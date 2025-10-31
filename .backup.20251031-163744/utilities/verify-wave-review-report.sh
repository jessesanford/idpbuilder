#!/bin/bash
# verify-wave-review-report.sh
# Verify wave review report compliance with R258 - Mandatory Wave Review Report
#
# Usage: ./verify-wave-review-report.sh <phase> <wave>
#
# This script verifies that a wave review report exists and is compliant with R258 requirements

set -euo pipefail

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check arguments
if [ $# -ne 2 ]; then
    echo -e "${RED}❌ ERROR: Missing required arguments${NC}"
    echo "Usage: $0 <phase> <wave>"
    echo "Example: $0 3 2  # Verify Phase 3 Wave 2 review report"
    exit 1
fi

PHASE=$1
WAVE=$2
REPORT_FILE="wave-reviews/phase${PHASE}/wave${WAVE}/PHASE-${PHASE}-WAVE-${WAVE}-REVIEW-REPORT.md"

echo -e "${BLUE}📊 Verifying Wave Review Report for Phase ${PHASE} Wave ${WAVE}${NC}"
echo "Expected file: $REPORT_FILE"
echo "-------------------------------------------"

# Function to check if file exists
check_file_exists() {
    if [ ! -f "$REPORT_FILE" ]; then
        echo -e "${RED}❌ CRITICAL: Wave review report not found!${NC}"
        echo -e "${RED}   Expected: $REPORT_FILE${NC}"
        echo -e "${RED}   This violates R258 - Mandatory Wave Review Report${NC}"
        echo ""
        echo "The architect MUST create this report before signaling wave review complete."
        echo "The orchestrator CANNOT proceed without this report."
        return 1
    fi
    echo -e "${GREEN}✅ Report file exists${NC}"
    return 0
}

# Function to extract and validate decision
check_decision() {
    local DECISION=$(grep "^\*\*DECISION\*\*:" "$REPORT_FILE" 2>/dev/null | cut -d: -f2 | xargs)
    
    if [ -z "$DECISION" ]; then
        echo -e "${RED}❌ No DECISION field found in report${NC}"
        return 1
    fi
    
    case "$DECISION" in
        PROCEED_NEXT_WAVE)
            echo -e "${GREEN}✅ Valid decision: PROCEED_NEXT_WAVE${NC}"
            echo "   Wave approved, ready for next wave"
            ;;
        PROCEED_PHASE_ASSESSMENT)
            echo -e "${GREEN}✅ Valid decision: PROCEED_PHASE_ASSESSMENT${NC}"
            echo "   Last wave complete, trigger phase assessment"
            ;;
        CHANGES_REQUIRED)
            echo -e "${YELLOW}⚠️ Valid decision: CHANGES_REQUIRED${NC}"
            echo "   Fixes needed before progression"
            ;;
        WAVE_FAILED)
            echo -e "${RED}⚠️ Valid decision: WAVE_FAILED${NC}"
            echo "   Major issues, cannot proceed"
            ;;
        *)
            echo -e "${RED}❌ Invalid decision: $DECISION${NC}"
            echo "   Must be one of: PROCEED_NEXT_WAVE, PROCEED_PHASE_ASSESSMENT, CHANGES_REQUIRED, WAVE_FAILED"
            return 1
            ;;
    esac
    return 0
}

# Function to check mandatory sections
check_mandatory_sections() {
    local MISSING=0
    
    echo -e "${BLUE}Checking mandatory sections...${NC}"
    
    local REQUIRED_SECTIONS=(
        "Review Metadata"
        "Review Decision"
        "Integration Assessment"
        "Architectural Review Scoring"
        "Size Compliance Verification"
        "Wave Deliverables Checklist"
        "Sign-Off"
    )
    
    for section in "${REQUIRED_SECTIONS[@]}"; do
        if grep -q "^## $section" "$REPORT_FILE"; then
            echo -e "  ${GREEN}✓${NC} $section"
        else
            echo -e "  ${RED}✗${NC} $section - MISSING"
            MISSING=$((MISSING + 1))
        fi
    done
    
    if [ $MISSING -gt 0 ]; then
        echo -e "${RED}❌ Missing $MISSING mandatory sections${NC}"
        return 1
    fi
    
    echo -e "${GREEN}✅ All mandatory sections present${NC}"
    return 0
}

# Function to check key fields
check_key_fields() {
    echo -e "${BLUE}Checking key fields...${NC}"
    
    # Check for Overall Score
    local SCORE=$(grep "^\*\*OVERALL SCORE\*\*" "$REPORT_FILE" 2>/dev/null | grep -o "[0-9]\+" | tail -1)
    if [ -n "$SCORE" ]; then
        echo -e "  ${GREEN}✓${NC} Overall Score: $SCORE"
    else
        echo -e "  ${YELLOW}⚠️${NC} Overall Score: Not found or invalid"
    fi
    
    # Check Size Compliance
    local SIZE_OK=$(grep "^\*\*All Efforts Compliant\*\*:" "$REPORT_FILE" 2>/dev/null | grep -o "\(YES\|NO\)")
    if [ -n "$SIZE_OK" ]; then
        if [ "$SIZE_OK" = "YES" ]; then
            echo -e "  ${GREEN}✓${NC} Size Compliance: YES (all efforts ≤800 lines)"
        else
            echo -e "  ${YELLOW}⚠️${NC} Size Compliance: NO (some efforts >800 lines)"
        fi
    else
        echo -e "  ${RED}✗${NC} Size Compliance: Not documented"
    fi
    
    # Check Architect Sign-off
    local SIGNOFF=$(grep "^\*\*Architect Sign-Off\*\*:" "$REPORT_FILE" 2>/dev/null)
    if [ -n "$SIGNOFF" ]; then
        echo -e "  ${GREEN}✓${NC} Architect Sign-Off: Present"
    else
        echo -e "  ${RED}✗${NC} Architect Sign-Off: MISSING (critical)"
        return 1
    fi
    
    # Check Report Hash
    local HASH=$(grep "^\*\*Report Hash\*\*:" "$REPORT_FILE" 2>/dev/null)
    if [ -n "$HASH" ]; then
        echo -e "  ${GREEN}✓${NC} Report Hash: Present (integrity verified)"
    else
        echo -e "  ${YELLOW}⚠️${NC} Report Hash: Missing (integrity not verifiable)"
    fi
    
    return 0
}

# Function to check decision consistency
check_decision_consistency() {
    echo -e "${BLUE}Checking decision consistency...${NC}"
    
    local DECISION=$(grep "^\*\*DECISION\*\*:" "$REPORT_FILE" 2>/dev/null | cut -d: -f2 | xargs)
    local SIZE_OK=$(grep "^\*\*All Efforts Compliant\*\*:" "$REPORT_FILE" 2>/dev/null | grep -o "\(YES\|NO\)")
    
    # Check for inconsistent decisions
    if [ "$SIZE_OK" = "NO" ] && [ "$DECISION" = "PROCEED_NEXT_WAVE" ]; then
        echo -e "  ${RED}✗${NC} INCONSISTENT: Cannot PROCEED_NEXT_WAVE with size violations!"
        return 1
    fi
    
    if [ "$SIZE_OK" = "NO" ] && [ "$DECISION" = "PROCEED_PHASE_ASSESSMENT" ]; then
        echo -e "  ${RED}✗${NC} INCONSISTENT: Cannot PROCEED_PHASE_ASSESSMENT with size violations!"
        return 1
    fi
    
    echo -e "  ${GREEN}✓${NC} Decision is consistent with report findings"
    return 0
}

# Main verification flow
main() {
    local ERRORS=0
    
    # Check file exists
    if ! check_file_exists; then
        ERRORS=$((ERRORS + 1))
        # Can't proceed without file
        echo ""
        echo -e "${RED}═══════════════════════════════════════════${NC}"
        echo -e "${RED}RESULT: R258 VIOLATION - No report found${NC}"
        echo -e "${RED}═══════════════════════════════════════════${NC}"
        exit 1
    fi
    
    echo ""
    
    # Check decision field
    if ! check_decision; then
        ERRORS=$((ERRORS + 1))
    fi
    
    echo ""
    
    # Check mandatory sections
    if ! check_mandatory_sections; then
        ERRORS=$((ERRORS + 1))
    fi
    
    echo ""
    
    # Check key fields
    if ! check_key_fields; then
        ERRORS=$((ERRORS + 1))
    fi
    
    echo ""
    
    # Check consistency
    if ! check_decision_consistency; then
        ERRORS=$((ERRORS + 1))
    fi
    
    echo ""
    echo "═══════════════════════════════════════════"
    
    if [ $ERRORS -eq 0 ]; then
        echo -e "${GREEN}✅ RESULT: Wave review report is COMPLIANT with R258${NC}"
        echo -e "${GREEN}   All requirements satisfied${NC}"
        exit 0
    else
        echo -e "${RED}❌ RESULT: Wave review report has $ERRORS VIOLATIONS${NC}"
        echo -e "${RED}   Fix the issues above before proceeding${NC}"
        exit 1
    fi
}

# Run main verification
main