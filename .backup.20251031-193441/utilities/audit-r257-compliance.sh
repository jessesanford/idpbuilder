#!/bin/bash

# R257 Compliance Audit Script
# Verifies that all phase assessment requirements are properly enforced

set -uo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo "=================================================="
echo "       R257 COMPLIANCE AUDIT REPORT"
echo "   Mandatory Phase Assessment Report Rule"
echo "=================================================="
echo ""

VIOLATIONS=0
WARNINGS=0
PASSED=0

# Function to check file for R257 reference
check_r257_reference() {
    local file="$1"
    local context="$2"
    
    if grep -q "R257" "$file" 2>/dev/null; then
        echo -e "  ${GREEN}✅ $context references R257${NC}"
        ((PASSED++))
        return 0
    else
        echo -e "  ${RED}❌ $context MISSING R257 reference${NC}"
        ((VIOLATIONS++))
        return 1
    fi
}

# Function to check for mandatory file creation language
check_mandatory_language() {
    local file="$1"
    local context="$2"
    
    if grep -qi "must create.*assessment.*report\|mandatory.*report.*file\|PHASE.*ASSESSMENT.*REPORT\.md" "$file" 2>/dev/null; then
        echo -e "  ${GREEN}✅ $context has mandatory report requirement${NC}"
        ((PASSED++))
        return 0
    else
        echo -e "  ${YELLOW}⚠️ $context may lack clear mandatory language${NC}"
        ((WARNINGS++))
        return 1
    fi
}

echo "1. RULE LIBRARY CHECKS"
echo "----------------------"

# Check R257 exists
if [ -f "rule-library/R257-mandatory-phase-assessment-report.md" ]; then
    echo -e "${GREEN}✅ R257 rule file exists${NC}"
    ((PASSED++))
    
    # Check rule has proper criticality
    if grep -q "BLOCKING" rule-library/R257-mandatory-phase-assessment-report.md; then
        echo -e "${GREEN}✅ R257 marked as BLOCKING criticality${NC}"
        ((PASSED++))
    else
        echo -e "${RED}❌ R257 not marked as BLOCKING${NC}"
        ((VIOLATIONS++))
    fi
else
    echo -e "${RED}❌ R257 rule file NOT FOUND${NC}"
    ((VIOLATIONS++))
fi

# Check rule registry
if grep -q "R257.*Mandatory Phase Assessment Report" rule-library/RULE-REGISTRY.md 2>/dev/null; then
    echo -e "${GREEN}✅ R257 registered in RULE-REGISTRY.md${NC}"
    ((PASSED++))
else
    echo -e "${RED}❌ R257 NOT in RULE-REGISTRY.md${NC}"
    ((VIOLATIONS++))
fi

echo ""
echo "2. ARCHITECT STATE CHECKS"
echo "-------------------------"

# Check architect PHASE_ASSESSMENT rules
ARCH_PHASE_RULES="agent-states/architect/PHASE_ASSESSMENT/rules.md"
if [ -f "$ARCH_PHASE_RULES" ]; then
    check_r257_reference "$ARCH_PHASE_RULES" "Architect PHASE_ASSESSMENT rules"
    check_mandatory_language "$ARCH_PHASE_RULES" "Architect PHASE_ASSESSMENT"
    
    # Check for specific file name requirement
    if grep -q "PHASE-{N}-ASSESSMENT-REPORT.md" "$ARCH_PHASE_RULES"; then
        echo -e "  ${GREEN}✅ Specifies exact report file name${NC}"
        ((PASSED++))
    else
        echo -e "  ${RED}❌ Missing exact file name specification${NC}"
        ((VIOLATIONS++))
    fi
else
    echo -e "${RED}❌ Architect PHASE_ASSESSMENT rules NOT FOUND${NC}"
    ((VIOLATIONS++))
fi

# Check architect grading
ARCH_GRADING="agent-states/architect/PHASE_ASSESSMENT/grading.md"
if [ -f "$ARCH_GRADING" ]; then
    check_r257_reference "$ARCH_GRADING" "Architect grading.md"
else
    echo -e "${YELLOW}⚠️ Architect grading.md not checked${NC}"
    ((WARNINGS++))
fi

echo ""
echo "3. ORCHESTRATOR STATE CHECKS"
echo "----------------------------"

# Check WAITING_FOR_PHASE_ASSESSMENT
ORCH_WAITING="agent-states/software-factory/orchestrator/WAITING_FOR_PHASE_ASSESSMENT/rules.md"
if [ -f "$ORCH_WAITING" ]; then
    check_r257_reference "$ORCH_WAITING" "Orchestrator WAITING_FOR_PHASE_ASSESSMENT"
    
    # Check for verification logic
    if grep -q "verify.*assessment.*report\|! -f.*ASSESSMENT-REPORT" "$ORCH_WAITING"; then
        echo -e "  ${GREEN}✅ Has report verification logic${NC}"
        ((PASSED++))
    else
        echo -e "  ${RED}❌ Missing report verification logic${NC}"
        ((VIOLATIONS++))
    fi
else
    echo -e "${RED}❌ WAITING_FOR_PHASE_ASSESSMENT rules NOT FOUND${NC}"
    ((VIOLATIONS++))
fi

# Check SPAWN_ARCHITECT_PHASE_ASSESSMENT
ORCH_SPAWN="agent-states/software-factory/orchestrator/SPAWN_ARCHITECT_PHASE_ASSESSMENT/rules.md"
if [ -f "$ORCH_SPAWN" ]; then
    if grep -q "R257\|mandatory.*report\|MUST create.*assessment.*report" "$ORCH_SPAWN"; then
        echo -e "${GREEN}✅ SPAWN_ARCHITECT instructs about report requirement${NC}"
        ((PASSED++))
    else
        echo -e "${YELLOW}⚠️ SPAWN_ARCHITECT may not emphasize report requirement${NC}"
        ((WARNINGS++))
    fi
else
    echo -e "${YELLOW}⚠️ SPAWN_ARCHITECT_PHASE_ASSESSMENT not checked${NC}"
    ((WARNINGS++))
fi

# Check COMPLETE_PHASE
ORCH_COMPLETE="agent-states/software-factory/orchestrator/COMPLETE_PHASE/rules.md"
if [ -f "$ORCH_COMPLETE" ]; then
    if grep -q "assessment.*report\|R257" "$ORCH_COMPLETE"; then
        echo -e "${GREEN}✅ COMPLETE_PHASE references assessment report${NC}"
        ((PASSED++))
    else
        echo -e "${YELLOW}⚠️ COMPLETE_PHASE may not reference report${NC}"
        ((WARNINGS++))
    fi
fi

echo ""
echo "4. UTILITY SCRIPTS"
echo "------------------"

# Check for verification script
if [ -f "utilities/verify-phase-assessment-report.sh" ]; then
    echo -e "${GREEN}✅ Verification utility script exists${NC}"
    ((PASSED++))
    
    if [ -x "utilities/verify-phase-assessment-report.sh" ]; then
        echo -e "${GREEN}✅ Verification script is executable${NC}"
        ((PASSED++))
    else
        echo -e "${YELLOW}⚠️ Verification script not executable${NC}"
        ((WARNINGS++))
    fi
else
    echo -e "${RED}❌ Verification utility script NOT FOUND${NC}"
    ((VIOLATIONS++))
fi

echo ""
echo "5. ENFORCEMENT VERIFICATION"
echo "---------------------------"

# Check if any state can transition to PROJECT_DONE without COMPLETE_PHASE
echo "Checking for illegal PROJECT_DONE transitions..."
ILLEGAL_TRANSITIONS=$(grep -r "transition_to.*PROJECT_DONE\|next_state.*=.*PROJECT_DONE" agent-states/software-factory/orchestrator/ 2>/dev/null | grep -v "COMPLETE_PHASE" | grep -v "rules.md.backup" | grep -v "Never transition\|NO transition\|not.*PROJECT_DONE" || true)
if [ -z "$ILLEGAL_TRANSITIONS" ]; then
    echo -e "${GREEN}✅ No illegal PROJECT_DONE transitions found${NC}"
    ((PASSED++))
else
    echo -e "${RED}❌ Found illegal PROJECT_DONE transitions:${NC}"
    echo "$ILLEGAL_TRANSITIONS"
    ((VIOLATIONS++))
fi

echo ""
echo "6. SAMPLE REPORT STRUCTURE"
echo "--------------------------"

# Check if R257 includes complete report template
if grep -q "Assessment Metadata\|Assessment Decision\|Scoring Summary\|Sign-Off" rule-library/R257-mandatory-phase-assessment-report.md 2>/dev/null; then
    echo -e "${GREEN}✅ R257 includes complete report template${NC}"
    ((PASSED++))
else
    echo -e "${YELLOW}⚠️ R257 may lack complete template${NC}"
    ((WARNINGS++))
fi

echo ""
echo "=================================================="
echo "            AUDIT SUMMARY"
echo "=================================================="
echo -e "Passed Checks:     ${GREEN}$PASSED${NC}"
echo -e "Warnings:          ${YELLOW}$WARNINGS${NC}"
echo -e "Violations:        ${RED}$VIOLATIONS${NC}"
echo ""

if [ $VIOLATIONS -eq 0 ]; then
    if [ $WARNINGS -eq 0 ]; then
        echo -e "${GREEN}🎉 FULL COMPLIANCE: R257 is properly enforced!${NC}"
        exit 0
    else
        echo -e "${YELLOW}⚠️ PARTIAL COMPLIANCE: R257 enforced with minor issues${NC}"
        exit 0
    fi
else
    echo -e "${RED}❌ NON-COMPLIANT: R257 enforcement has critical gaps!${NC}"
    echo ""
    echo "Required fixes:"
    echo "1. Ensure all architect PHASE_ASSESSMENT states require report creation"
    echo "2. Ensure orchestrator verifies report existence before proceeding"
    echo "3. Update all relevant states to reference R257"
    echo "4. Add verification logic to prevent bypass"
    exit 1
fi