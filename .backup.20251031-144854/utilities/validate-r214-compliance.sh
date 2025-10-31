#!/bin/bash

# validate-r214-compliance.sh
# Validates that Code Reviewers acknowledge wave directories before creating effort plans
# Ensures effort plans are created in orchestrator-defined locations

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Counters
TOTAL_WAVES=0
COMPLIANT_WAVES=0
NON_COMPLIANT_WAVES=0
WARNINGS=0

echo "════════════════════════════════════════════════════════════════"
echo "🔍 R214: Code Reviewer Wave Directory Acknowledgment Validation"
echo "════════════════════════════════════════════════════════════════"
echo ""
echo "Checking that Code Reviewers use orchestrator-defined directories..."
echo ""

# Function to check R214 acknowledgment file
check_acknowledgment_file() {
    echo "Checking for R214 acknowledgment file..."
    
    if [ -f ".r214-wave-acknowledged" ]; then
        echo -e "${GREEN}✅ R214 acknowledgment file exists${NC}"
        echo "Recent acknowledgments:"
        tail -10 .r214-wave-acknowledged | while read line; do
            echo "  $line"
        done
        return 0
    else
        echo -e "${YELLOW}⚠️  No R214 acknowledgment file found${NC}"
        echo "  Code Reviewers haven't acknowledged any wave directories yet"
        return 1
    fi
}

# Function to validate wave compliance
validate_wave_r214_compliance() {
    local PHASE="$1"
    local WAVE="$2"
    local WAVE_IMPL_PLAN="phase-plans/PHASE-${PHASE}-WAVE-${WAVE}-IMPLEMENTATION-PLAN.md"
    
    echo ""
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${CYAN}Phase ${PHASE} Wave ${WAVE} - R214 Compliance Check${NC}"
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    
    ((TOTAL_WAVES++))
    
    # Check if wave implementation plan exists
    if [ ! -f "$WAVE_IMPL_PLAN" ]; then
        echo -e "${YELLOW}⚠️  Wave implementation plan not found${NC}"
        return 0
    fi
    
    # Check for R213 metadata (orchestrator-defined structure)
    if ! grep -q "WAVE INFRASTRUCTURE METADATA" "$WAVE_IMPL_PLAN"; then
        echo -e "${RED}❌ Wave plan missing R213 metadata!${NC}"
        echo "  Orchestrator must inject metadata first"
        ((NON_COMPLIANT_WAVES++))
        return 1
    fi
    
    # Verify orchestrator is the source
    METADATA_SOURCE=$(grep "**METADATA_SOURCE**:" "$WAVE_IMPL_PLAN" 2>/dev/null | cut -d: -f2- | xargs || true)
    if [ "$METADATA_SOURCE" != "ORCHESTRATOR" ]; then
        echo -e "${RED}❌ Metadata source is not ORCHESTRATOR!${NC}"
        echo "  Found: $METADATA_SOURCE"
        echo "  R214 requires ORCHESTRATOR as source"
        ((NON_COMPLIANT_WAVES++))
        return 1
    fi
    
    echo -e "${GREEN}✅ Wave has orchestrator metadata${NC}"
    
    # Extract wave structure
    WAVE_ROOT=$(grep "**WAVE_ROOT**:" "$WAVE_IMPL_PLAN" | cut -d: -f2- | xargs)
    EFFORT_COUNT=$(grep "**EFFORT_COUNT**:" "$WAVE_IMPL_PLAN" | cut -d: -f2- | xargs || echo "0")
    
    echo "  Wave root: $WAVE_ROOT"
    echo "  Expected efforts: $EFFORT_COUNT"
    
    # Check if wave directory exists
    if [ ! -d "$WAVE_ROOT" ]; then
        echo -e "${YELLOW}⚠️  Wave directory not created yet${NC}"
        return 0
    fi
    
    # Check effort plans for R214 compliance
    local COMPLIANT_EFFORTS=0
    local NON_COMPLIANT_EFFORTS=0
    
    for i in $(seq 1 $EFFORT_COUNT); do
        EFFORT_DIR="${WAVE_ROOT}/effort-${i}"
        EFFORT_PLAN="${EFFORT_DIR}/IMPLEMENTATION-PLAN.md"
        
        echo ""
        echo "  Checking effort ${i}..."
        
        if [ ! -d "$EFFORT_DIR" ]; then
            echo -e "${YELLOW}    ⚠ Directory not created yet${NC}"
            continue
        fi
        
        if [ ! -f "$EFFORT_PLAN" ]; then
            echo -e "${YELLOW}    ⚠ No implementation plan yet${NC}"
            continue
        fi
        
        # Check if effort plan is in correct location
        if [[ "$EFFORT_PLAN" == "${WAVE_ROOT}/effort-${i}/IMPLEMENTATION-PLAN.md" ]]; then
            echo -e "${GREEN}    ✓ Plan in correct location${NC}"
            
            # Check for R214 compliance note
            if grep -q "R214 Compliance" "$EFFORT_PLAN" 2>/dev/null; then
                echo -e "${GREEN}    ✓ Has R214 compliance note${NC}"
                ((COMPLIANT_EFFORTS++))
            else
                echo -e "${YELLOW}    ⚠ Missing R214 compliance note${NC}"
                echo "      Code Reviewer should add compliance note"
                ((WARNINGS++))
            fi
            
            # Check if location in plan matches actual location
            if grep -q "orchestrator-defined directory" "$EFFORT_PLAN" 2>/dev/null; then
                NOTED_DIR=$(grep "orchestrator-defined directory:" "$EFFORT_PLAN" | cut -d: -f2- | xargs || true)
                if [ -n "$NOTED_DIR" ]; then
                    if [[ "$NOTED_DIR" == "$EFFORT_DIR" ]]; then
                        echo -e "${GREEN}    ✓ Noted directory matches actual${NC}"
                    else
                        echo -e "${RED}    ❌ Directory mismatch!${NC}"
                        echo "      Noted: $NOTED_DIR"
                        echo "      Actual: $EFFORT_DIR"
                        ((NON_COMPLIANT_EFFORTS++))
                    fi
                fi
            fi
        else
            echo -e "${RED}    ❌ Plan in WRONG location!${NC}"
            echo "      Expected: ${WAVE_ROOT}/effort-${i}/IMPLEMENTATION-PLAN.md"
            echo "      Found: $EFFORT_PLAN"
            ((NON_COMPLIANT_EFFORTS++))
        fi
    done
    
    # Wave summary
    echo ""
    if [ $NON_COMPLIANT_EFFORTS -eq 0 ]; then
        echo -e "${GREEN}✅ Wave ${WAVE}: R214 COMPLIANT${NC}"
        echo "  All effort plans in orchestrator-defined locations"
        ((COMPLIANT_WAVES++))
    else
        echo -e "${RED}❌ Wave ${WAVE}: R214 VIOLATIONS FOUND${NC}"
        echo "  $NON_COMPLIANT_EFFORTS efforts have issues"
        ((NON_COMPLIANT_WAVES++))
    fi
}

# Function to check for acknowledgment in wave plans
check_wave_acknowledgments() {
    echo ""
    echo -e "${BLUE}════════════════════════════════════════════════════════${NC}"
    echo -e "${BLUE}Checking Wave Acknowledgments in .r214 File${NC}"
    echo -e "${BLUE}════════════════════════════════════════════════════════${NC}"
    
    if [ ! -f ".r214-wave-acknowledged" ]; then
        echo -e "${YELLOW}No acknowledgments recorded yet${NC}"
        return
    fi
    
    # Parse acknowledgment file for waves
    while IFS= read -r line; do
        if [[ "$line" == *"Wave"*"Acknowledged"* ]]; then
            WAVE_NUM=$(echo "$line" | grep -oE "Wave [0-9]+" | cut -d' ' -f2)
            echo -e "${GREEN}✓ Wave $WAVE_NUM acknowledged by Code Reviewer${NC}"
        fi
        if [[ "$line" == *"wave root:"* ]]; then
            ACKED_ROOT=$(echo "$line" | cut -d: -f2- | xargs)
            echo "    Acknowledged root: $ACKED_ROOT"
        fi
    done < .r214-wave-acknowledged
}

# Main validation
echo "Starting R214 compliance validation..."

# Check acknowledgment file
check_acknowledgment_file

# Check each phase and wave
for PHASE in {1..5}; do
    for WAVE in {1..6}; do
        # Check if wave plan exists
        if ls phase-plans/PHASE-${PHASE}-WAVE-${WAVE}-*.md 2>/dev/null | head -1 > /dev/null; then
            validate_wave_r214_compliance "$PHASE" "$WAVE"
        fi
    done
done

# Check acknowledgments
check_wave_acknowledgments

# Summary report
echo ""
echo "════════════════════════════════════════════════════════════════"
echo "📊 R214 Compliance Summary"
echo "════════════════════════════════════════════════════════════════"
echo "Total Waves Checked: $TOTAL_WAVES"
echo -e "${GREEN}Compliant Waves: $COMPLIANT_WAVES${NC}"
echo -e "${RED}Non-Compliant Waves: $NON_COMPLIANT_WAVES${NC}"
echo -e "${YELLOW}Warnings: $WARNINGS${NC}"
echo ""

if [ $NON_COMPLIANT_WAVES -eq 0 ] && [ $WARNINGS -eq 0 ]; then
    echo -e "${GREEN}🎉 PERFECT R214 COMPLIANCE!${NC}"
    echo "All Code Reviewers are using orchestrator-defined directories!"
    echo "All effort plans are in correct locations!"
    exit 0
elif [ $NON_COMPLIANT_WAVES -eq 0 ]; then
    echo -e "${YELLOW}⚠️  MOSTLY COMPLIANT with $WARNINGS warnings${NC}"
    echo ""
    echo "Recommendations:"
    echo "1. Add R214 compliance notes to all effort plans"
    echo "2. Ensure all acknowledgments are recorded"
    echo "3. Verify directory paths in plans match actual locations"
    exit 0
else
    echo -e "${RED}❌ R214 COMPLIANCE FAILURES!${NC}"
    echo ""
    echo "CRITICAL REQUIREMENTS:"
    echo "1. Code Reviewers MUST acknowledge wave directories (R214)"
    echo "2. Must verify ORCHESTRATOR is metadata source"
    echo "3. Must create effort plans in orchestrator-defined locations"
    echo "4. Must add R214 compliance notes to plans"
    echo ""
    echo "Code Reviewers are NOT properly following orchestrator structure!"
    exit 1
fi