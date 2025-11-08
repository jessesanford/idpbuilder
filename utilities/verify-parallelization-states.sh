#!/bin/bash

# verify-parallelization-states.sh
# Verifies that the new parallelization analysis states are properly configured

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

echo "════════════════════════════════════════════════════════════════"
echo -e "${BOLD}PARALLELIZATION STATES VERIFICATION${NC}"
echo "Checking new mandatory analysis states configuration"
echo "════════════════════════════════════════════════════════════════"
echo ""

# Set base directory
BASE_DIR="/workspaces/software-factory-2.0-template"

# Counters
TOTAL_CHECKS=0
PASSED_CHECKS=0
WARNINGS=0

# Function to check if a file exists
check_file() {
    local FILE="$1"
    local DESC="$2"
    
    ((TOTAL_CHECKS++))
    
    if [ -f "$FILE" ]; then
        echo -e "${GREEN}✅ $DESC exists${NC}"
        ((PASSED_CHECKS++))
        return 0
    else
        echo -e "${RED}❌ $DESC MISSING: $FILE${NC}"
        return 1
    fi
}

# Function to check file content
check_content() {
    local FILE="$1"
    local PATTERN="$2"
    local DESC="$3"
    
    ((TOTAL_CHECKS++))
    
    if [ -f "$FILE" ]; then
        if grep -q "$PATTERN" "$FILE" 2>/dev/null; then
            echo -e "${GREEN}✅ $DESC found in $FILE${NC}"
            ((PASSED_CHECKS++))
            return 0
        else
            echo -e "${RED}❌ $DESC NOT FOUND in $FILE${NC}"
            return 1
        fi
    else
        echo -e "${RED}❌ File doesn't exist: $FILE${NC}"
        return 1
    fi
}

# Function to check state machine
check_state_machine() {
    local STATE="$1"
    local DESC="$2"
    
    ((TOTAL_CHECKS++))
    
    if grep -q "$STATE" "$BASE_DIR/software-factory-3.0-state-machine.json" 2>/dev/null; then
        echo -e "${GREEN}✅ $DESC in state machine${NC}"
        ((PASSED_CHECKS++))
        return 0
    else
        echo -e "${RED}❌ $DESC NOT in state machine${NC}"
        return 1
    fi
}

echo -e "${BLUE}1. Checking state directories...${NC}"
echo "────────────────────────────────────────────"

# Check ANALYZE_CODE_REVIEWER_PARALLELIZATION state
check_file "$BASE_DIR/agent-states/software-factory/orchestrator/ANALYZE_CODE_REVIEWER_PARALLELIZATION/rules.md" \
    "ANALYZE_CODE_REVIEWER_PARALLELIZATION rules"

check_file "$BASE_DIR/agent-states/software-factory/orchestrator/ANALYZE_CODE_REVIEWER_PARALLELIZATION/grading.md" \
    "ANALYZE_CODE_REVIEWER_PARALLELIZATION grading"

check_file "$BASE_DIR/agent-states/software-factory/orchestrator/ANALYZE_CODE_REVIEWER_PARALLELIZATION/checkpoint.md" \
    "ANALYZE_CODE_REVIEWER_PARALLELIZATION checkpoint"

# Check ANALYZE_IMPLEMENTATION_PARALLELIZATION state
check_file "$BASE_DIR/agent-states/software-factory/orchestrator/ANALYZE_IMPLEMENTATION_PARALLELIZATION/rules.md" \
    "ANALYZE_IMPLEMENTATION_PARALLELIZATION rules"

check_file "$BASE_DIR/agent-states/software-factory/orchestrator/ANALYZE_IMPLEMENTATION_PARALLELIZATION/grading.md" \
    "ANALYZE_IMPLEMENTATION_PARALLELIZATION grading"

check_file "$BASE_DIR/agent-states/software-factory/orchestrator/ANALYZE_IMPLEMENTATION_PARALLELIZATION/checkpoint.md" \
    "ANALYZE_IMPLEMENTATION_PARALLELIZATION checkpoint"

echo ""
echo -e "${BLUE}2. Checking state machine updates...${NC}"
echo "────────────────────────────────────────────"

check_state_machine "ANALYZE_CODE_REVIEWER_PARALLELIZATION" \
    "ANALYZE_CODE_REVIEWER_PARALLELIZATION state"

check_state_machine "ANALYZE_IMPLEMENTATION_PARALLELIZATION" \
    "ANALYZE_IMPLEMENTATION_PARALLELIZATION state"

# Check transitions
check_content "$BASE_DIR/software-factory-3.0-state-machine.json" \
    "CREATE_NEXT_INFRASTRUCTURE → ANALYZE_CODE_REVIEWER_PARALLELIZATION" \
    "Transition to ANALYZE_CODE_REVIEWER_PARALLELIZATION"

check_content "$BASE_DIR/software-factory-3.0-state-machine.json" \
    "WAITING_FOR_EFFORT_PLANS → ANALYZE_IMPLEMENTATION_PARALLELIZATION" \
    "Transition to ANALYZE_IMPLEMENTATION_PARALLELIZATION"

echo ""
echo -e "${BLUE}3. Checking rule enforcement...${NC}"
echo "────────────────────────────────────────────"

# Check R218 reference
check_content "$BASE_DIR/agent-states/software-factory/orchestrator/ANALYZE_CODE_REVIEWER_PARALLELIZATION/rules.md" \
    "R218" "R218 rule reference"

# Check R151 reference
check_content "$BASE_DIR/agent-states/software-factory/orchestrator/ANALYZE_CODE_REVIEWER_PARALLELIZATION/rules.md" \
    "R151" "R151 rule reference"

# Check mandatory gate enforcement
check_content "$BASE_DIR/agent-states/software-factory/orchestrator/ANALYZE_CODE_REVIEWER_PARALLELIZATION/rules.md" \
    "MANDATORY GATE" "Mandatory gate enforcement"

check_content "$BASE_DIR/agent-states/software-factory/orchestrator/ANALYZE_IMPLEMENTATION_PARALLELIZATION/rules.md" \
    "MANDATORY GATE" "Mandatory gate enforcement"

echo ""
echo -e "${BLUE}4. Checking updated spawn states...${NC}"
echo "────────────────────────────────────────────"

# Check SPAWN_CODE_REVIEWERS_EFFORT_PLANNING updates
check_content "$BASE_DIR/agent-states/software-factory/orchestrator/SPAWN_CODE_REVIEWERS_EFFORT_PLANNING/rules.md" \
    "ANALYZE_CODE_REVIEWER_PARALLELIZATION" \
    "Reference to ANALYZE_CODE_REVIEWER_PARALLELIZATION"

# Check SPAWN_SW_ENGINEERS updates
check_content "$BASE_DIR/agent-states/software-factory/orchestrator/SPAWN_SW_ENGINEERS/rules.md" \
    "ANALYZE_IMPLEMENTATION_PARALLELIZATION" \
    "Reference to ANALYZE_IMPLEMENTATION_PARALLELIZATION"

echo ""
echo -e "${BLUE}5. Checking state directory map...${NC}"
echo "────────────────────────────────────────────"

check_content "$BASE_DIR/agent-states/software-factory/orchestrator/STATE-DIRECTORY-MAP.md" \
    "ANALYZE_CODE_REVIEWER_PARALLELIZATION" \
    "ANALYZE_CODE_REVIEWER_PARALLELIZATION in map"

check_content "$BASE_DIR/agent-states/software-factory/orchestrator/STATE-DIRECTORY-MAP.md" \
    "ANALYZE_IMPLEMENTATION_PARALLELIZATION" \
    "ANALYZE_IMPLEMENTATION_PARALLELIZATION in map"

echo ""
echo -e "${BLUE}6. Checking for critical enforcement...${NC}"
echo "────────────────────────────────────────────"

# Check for blocking requirements
check_content "$BASE_DIR/agent-states/software-factory/orchestrator/ANALYZE_CODE_REVIEWER_PARALLELIZATION/rules.md" \
    "Cannot proceed without compliance" \
    "Blocking enforcement"

check_content "$BASE_DIR/agent-states/software-factory/orchestrator/ANALYZE_IMPLEMENTATION_PARALLELIZATION/rules.md" \
    "Cannot proceed without compliance" \
    "Blocking enforcement"

# Check for grading impact
check_content "$BASE_DIR/agent-states/software-factory/orchestrator/ANALYZE_CODE_REVIEWER_PARALLELIZATION/grading.md" \
    "Skipping this state: ..-50%.." \
    "Critical grading penalty"

echo ""
echo -e "${BLUE}7. Checking utility scripts...${NC}"
echo "────────────────────────────────────────────"

check_file "$BASE_DIR/utilities/upgrade-parallelization-states.sh" \
    "Upgrade script for parallelization states"

if [ -f "$BASE_DIR/utilities/upgrade-parallelization-states.sh" ]; then
    if [ -x "$BASE_DIR/utilities/upgrade-parallelization-states.sh" ]; then
        echo -e "${GREEN}✅ Upgrade script is executable${NC}"
        ((PASSED_CHECKS++))
        ((TOTAL_CHECKS++))
    else
        echo -e "${YELLOW}⚠️  Upgrade script is not executable${NC}"
        ((WARNINGS++))
        ((TOTAL_CHECKS++))
    fi
fi

echo ""
echo "════════════════════════════════════════════════════════════════"
echo -e "${BOLD}VERIFICATION SUMMARY${NC}"
echo "────────────────────────────────────────────"
echo -e "Total Checks: ${TOTAL_CHECKS}"
echo -e "Passed: ${GREEN}${PASSED_CHECKS}${NC}"
echo -e "Failed: ${RED}$((TOTAL_CHECKS - PASSED_CHECKS))${NC}"
if [ $WARNINGS -gt 0 ]; then
    echo -e "Warnings: ${YELLOW}${WARNINGS}${NC}"
fi
echo ""

if [ $PASSED_CHECKS -eq $TOTAL_CHECKS ]; then
    echo -e "${GREEN}${BOLD}✅ ALL CHECKS PASSED!${NC}"
    echo ""
    echo "The parallelization analysis states are properly configured."
    echo "The orchestrator will now be FORCED to analyze parallelization"
    echo "before spawning agents, preventing R151 violations."
else
    echo -e "${RED}${BOLD}❌ SOME CHECKS FAILED${NC}"
    echo ""
    echo "Please review the failed checks above and ensure all"
    echo "parallelization state files are properly configured."
    echo ""
    echo -e "${YELLOW}To fix missing components, you may need to:${NC}"
    echo "  1. Re-run the creation script"
    echo "  2. Manually create missing files"
    echo "  3. Update the state machine"
fi

echo "════════════════════════════════════════════════════════════════"

# Exit with appropriate code
if [ $PASSED_CHECKS -eq $TOTAL_CHECKS ]; then
    exit 0
else
    exit 1
fi