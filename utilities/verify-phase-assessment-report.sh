#!/bin/bash

# Utility function to verify phase assessment report per R257
# This function MUST be called by orchestrator in WAITING_FOR_PHASE_ASSESSMENT state

set -euo pipefail

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

verify_phase_assessment_report() {
    local PHASE="${1:-}"
    
    if [ -z "$PHASE" ]; then
        echo -e "${RED}❌ CRITICAL: Phase number not provided!${NC}"
        return 1
    fi
    
    local REPORT_FILE="phase-assessments/phase${PHASE}/PHASE-${PHASE}-ASSESSMENT-REPORT.md"
    
    echo "🔍 Verifying phase assessment report for Phase $PHASE..."
    echo "📄 Expected location: $REPORT_FILE"
    
    # Check if file exists
    if [ ! -f "$REPORT_FILE" ]; then
        echo -e "${RED}❌ CRITICAL: No phase assessment report found!${NC}"
        echo -e "${RED}❌ Expected: $REPORT_FILE${NC}"
        echo -e "${RED}❌ This violates R257 - Mandatory Phase Assessment Report${NC}"
        echo -e "${RED}❌ Cannot proceed without assessment report${NC}"
        return 1
    fi
    
    echo -e "${GREEN}✅ Report file exists${NC}"
    
    # Validate mandatory sections
    local MISSING_SECTIONS=()
    local REQUIRED_SECTIONS=(
        "Assessment Metadata"
        "Assessment Decision"
        "Scoring Summary"
        "Assessment Details"
        "Sign-Off"
    )
    
    for section in "${REQUIRED_SECTIONS[@]}"; do
        if ! grep -q "^## $section" "$REPORT_FILE"; then
            MISSING_SECTIONS+=("$section")
        fi
    done
    
    if [ ${#MISSING_SECTIONS[@]} -gt 0 ]; then
        echo -e "${RED}❌ Missing mandatory sections:${NC}"
        for section in "${MISSING_SECTIONS[@]}"; do
            echo -e "${RED}  - $section${NC}"
        done
        return 1
    fi
    
    echo -e "${GREEN}✅ All mandatory sections present${NC}"
    
    # Extract and validate decision
    local DECISION=$(grep "^\*\*DECISION\*\*:" "$REPORT_FILE" 2>/dev/null | cut -d: -f2 | xargs)
    
    if [ -z "$DECISION" ]; then
        echo -e "${RED}❌ Invalid assessment report - no decision found!${NC}"
        return 1
    fi
    
    # Validate decision value
    if [[ ! "$DECISION" =~ ^(PHASE_COMPLETE|NEEDS_WORK|PHASE_FAILED)$ ]]; then
        echo -e "${RED}❌ Invalid decision value: $DECISION${NC}"
        echo -e "${RED}   Must be one of: PHASE_COMPLETE, NEEDS_WORK, PHASE_FAILED${NC}"
        return 1
    fi
    
    echo -e "${GREEN}✅ Valid decision found: $DECISION${NC}"
    
    # Extract score if present
    local SCORE=$(grep "\*\*TOTAL SCORE\*\*" "$REPORT_FILE" 2>/dev/null | grep -oE "[0-9]+" | tail -1)
    
    if [ -n "$SCORE" ]; then
        echo "📊 Phase assessment score: $SCORE"
    fi
    
    # Verify architect sign-off
    if ! grep -q "^\*\*Architect Sign-Off\*\*:" "$REPORT_FILE"; then
        echo -e "${YELLOW}⚠️ Warning: No architect sign-off found${NC}"
    else
        local SIGNOFF_DATE=$(grep "^\*\*Architect Sign-Off\*\*:" "$REPORT_FILE" | cut -d: -f2-)
        echo -e "${GREEN}✅ Architect sign-off present: $SIGNOFF_DATE${NC}"
    fi
    
    # Check for report hash (integrity verification)
    if grep -q "^\*\*Report Hash\*\*:" "$REPORT_FILE"; then
        echo -e "${GREEN}✅ Report hash present for integrity verification${NC}"
    else
        echo -e "${YELLOW}⚠️ Warning: No report hash found${NC}"
    fi
    
    # Update state file with report information
    if command -v yq &> /dev/null && [ -f "orchestrator-state.json" ]; then
        echo "📝 Updating state file with report information..."
        yq -i ".phase_assessment.report_file = \"$REPORT_FILE\"" orchestrator-state.json
        yq -i ".phase_assessment.decision = \"$DECISION\"" orchestrator-state.json
        [ -n "$SCORE" ] && yq -i ".phase_assessment.score = $SCORE" orchestrator-state.json
        yq -i ".phase_assessment.verified_at = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" orchestrator-state.json
        echo -e "${GREEN}✅ State file updated${NC}"
    fi
    
    # Summary
    echo ""
    echo "═══════════════════════════════════════════════════════"
    echo -e "${GREEN}✅ PHASE ASSESSMENT REPORT VERIFICATION COMPLETE${NC}"
    echo "═══════════════════════════════════════════════════════"
    echo "  📄 Report: $REPORT_FILE"
    echo "  📊 Decision: $DECISION"
    [ -n "$SCORE" ] && echo "  📈 Score: $SCORE"
    echo "  ✅ All validations passed"
    echo "═══════════════════════════════════════════════════════"
    echo ""
    
    return 0
}

# Function to list all phase assessment reports
list_phase_assessment_reports() {
    echo "📋 Listing all phase assessment reports..."
    
    if [ ! -d "phase-assessments" ]; then
        echo "  No phase-assessments directory found"
        return 0
    fi
    
    local FOUND_REPORTS=0
    for phase_dir in phase-assessments/phase*/; do
        if [ -d "$phase_dir" ]; then
            local PHASE_NUM=$(basename "$phase_dir" | grep -oE "[0-9]+")
            local REPORT_FILE="${phase_dir}PHASE-${PHASE_NUM}-ASSESSMENT-REPORT.md"
            
            if [ -f "$REPORT_FILE" ]; then
                local DECISION=$(grep "^\*\*DECISION\*\*:" "$REPORT_FILE" 2>/dev/null | cut -d: -f2 | xargs)
                local DATE=$(grep "^\*\*Assessment Date\*\*:" "$REPORT_FILE" 2>/dev/null | cut -d: -f2 | xargs)
                echo "  ✅ Phase $PHASE_NUM: $DECISION (assessed: $DATE)"
                FOUND_REPORTS=$((FOUND_REPORTS + 1))
            else
                echo -e "  ${YELLOW}⚠️ Phase $PHASE_NUM: Directory exists but no report found${NC}"
            fi
        fi
    done
    
    if [ $FOUND_REPORTS -eq 0 ]; then
        echo "  No assessment reports found"
    else
        echo "  Total reports: $FOUND_REPORTS"
    fi
}

# Function to audit phase completions vs reports
audit_phase_completions() {
    echo "🔍 Auditing phase completions vs assessment reports..."
    
    if [ ! -f "orchestrator-state.json" ]; then
        echo -e "${YELLOW}⚠️ No orchestrator-state.json found${NC}"
        return 1
    fi
    
    local CURRENT_PHASE=$(yq '.current_phase' orchestrator-state.json)
    echo "  Current phase: $CURRENT_PHASE"
    
    # Check each completed phase for assessment report
    for phase in $(seq 1 $((CURRENT_PHASE - 1))); do
        local REPORT_FILE="phase-assessments/phase${phase}/PHASE-${phase}-ASSESSMENT-REPORT.md"
        
        if [ -f "$REPORT_FILE" ]; then
            echo -e "  ${GREEN}✅ Phase $phase: Has assessment report${NC}"
        else
            echo -e "  ${RED}❌ Phase $phase: MISSING assessment report (R257 violation!)${NC}"
        fi
    done
}

# Main execution
case "${1:-verify}" in
    verify)
        if [ -z "${2:-}" ]; then
            echo "Usage: $0 verify <phase_number>"
            exit 1
        fi
        verify_phase_assessment_report "$2"
        ;;
    list)
        list_phase_assessment_reports
        ;;
    audit)
        audit_phase_completions
        ;;
    *)
        echo "Usage: $0 {verify|list|audit} [phase_number]"
        echo "  verify <phase>  - Verify assessment report for specific phase"
        echo "  list            - List all phase assessment reports"
        echo "  audit           - Audit phase completions vs reports"
        exit 1
        ;;
esac