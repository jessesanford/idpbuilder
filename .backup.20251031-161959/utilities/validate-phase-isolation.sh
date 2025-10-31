#!/bin/bash

# validate-phase-isolation.sh
# Validates R212 phase directory isolation protocol compliance
# Checks phase plans for required metadata and validates agent acknowledgments

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Counters
TOTAL_PHASES=0
VALID_PHASES=0
INVALID_PHASES=0
WARNINGS=0

echo "════════════════════════════════════════════════════════════════"
echo "🔍 R212: Phase Directory Isolation Validation"
echo "════════════════════════════════════════════════════════════════"

# Function to validate phase plan metadata
validate_phase_plan() {
    local PLAN="$1"
    local PLAN_NAME=$(basename "$PLAN")
    local ERRORS=0
    
    echo ""
    echo "Checking: $PLAN_NAME"
    echo "─────────────────────────────────────────────────────────"
    
    # Check if plan exists
    if [ ! -f "$PLAN" ]; then
        echo -e "${YELLOW}⚠️  Plan not found yet${NC}"
        return 0
    fi
    
    # Required metadata fields per R212
    local REQUIRED_FIELDS=(
        "PHASE INFRASTRUCTURE METADATA"
        "PHASE_NUMBER"
        "WORKING_DIRECTORY"
        "PHASE_ROOT"
        "EFFORTS_ROOT"
        "PHASE_PLANS_DIR"
        "INTEGRATE_WAVE_EFFORTS_BRANCH"
        "CRITICAL PHASE ISOLATION RULES"
    )
    
    # Check each required field
    for field in "${REQUIRED_FIELDS[@]}"; do
        if grep -q "$field" "$PLAN"; then
            echo -e "${GREEN}✅ Found: $field${NC}"
        else
            echo -e "${RED}❌ Missing: $field${NC}"
            ((ERRORS++))
        fi
    done
    
    # Extract and validate metadata values
    if grep -q "PHASE_NUMBER" "$PLAN"; then
        PHASE_NUM=$(grep "**PHASE_NUMBER**:" "$PLAN" 2>/dev/null | cut -d: -f2- | xargs || true)
        PHASE_ROOT=$(grep "**PHASE_ROOT**:" "$PLAN" 2>/dev/null | cut -d: -f2- | xargs || true)
        EFFORTS_ROOT=$(grep "**EFFORTS_ROOT**:" "$PLAN" 2>/dev/null | cut -d: -f2- | xargs || true)
        
        if [ -n "$PHASE_NUM" ]; then
            echo -e "${GREEN}✅ Phase number: ${PHASE_NUM}${NC}"
            
            # Check if directories match expected structure
            if [[ "$PHASE_ROOT" == *"/phase${PHASE_NUM}" ]]; then
                echo -e "${GREEN}✅ Phase root correctly formatted${NC}"
            else
                echo -e "${RED}❌ Phase root mismatch!${NC}"
                echo "   Expected pattern: */phase${PHASE_NUM}"
                echo "   Found: $PHASE_ROOT"
                ((ERRORS++))
            fi
            
            if [[ "$EFFORTS_ROOT" == *"/efforts/phase${PHASE_NUM}" ]]; then
                echo -e "${GREEN}✅ Efforts root correctly formatted${NC}"
            else
                echo -e "${RED}❌ Efforts root mismatch!${NC}"
                echo "   Expected pattern: */efforts/phase${PHASE_NUM}"
                echo "   Found: $EFFORTS_ROOT"
                ((ERRORS++))
            fi
        fi
    fi
    
    return $ERRORS
}

# Function to check phase directory structure
validate_phase_structure() {
    local PHASE="$1"
    local ERRORS=0
    
    echo ""
    echo "Validating Phase ${PHASE} directory structure..."
    
    # Check phase plans directory
    if [ -d "phase-plans" ]; then
        echo -e "${GREEN}✅ phase-plans/ directory exists${NC}"
        
        # Count phase plans
        local ARCH_PLAN="phase-plans/PHASE-${PHASE}-ARCHITECTURE-PLAN.md"
        local IMPL_PLAN="phase-plans/PHASE-${PHASE}-IMPLEMENTATION-PLAN.md"
        
        [ -f "$ARCH_PLAN" ] && echo -e "${GREEN}✅ Architecture plan exists${NC}" || echo -e "${YELLOW}⚠️  Architecture plan not created yet${NC}"
        [ -f "$IMPL_PLAN" ] && echo -e "${GREEN}✅ Implementation plan exists${NC}" || echo -e "${YELLOW}⚠️  Implementation plan not created yet${NC}"
    else
        echo -e "${RED}❌ phase-plans/ directory missing${NC}"
        ((ERRORS++))
    fi
    
    # Check phase root directory
    if [ -d "phase${PHASE}" ]; then
        echo -e "${GREEN}✅ phase${PHASE}/ directory exists${NC}"
    else
        echo -e "${YELLOW}⚠️  phase${PHASE}/ directory not created yet${NC}"
    fi
    
    # Check efforts directory
    if [ -d "efforts/phase${PHASE}" ]; then
        echo -e "${GREEN}✅ efforts/phase${PHASE}/ directory exists${NC}"
        
        # Count waves
        local WAVE_COUNT=$(find "efforts/phase${PHASE}" -type d -name "wave*" 2>/dev/null | wc -l)
        if [ $WAVE_COUNT -gt 0 ]; then
            echo -e "${GREEN}✅ Found ${WAVE_COUNT} wave directories${NC}"
        fi
    else
        echo -e "${YELLOW}⚠️  efforts/phase${PHASE}/ directory not created yet${NC}"
    fi
    
    return $ERRORS
}

# Function to check acknowledgment files
check_acknowledgments() {
    echo ""
    echo "Checking R212 acknowledgments..."
    
    if [ -f ".r212-phase-acknowledged" ]; then
        echo -e "${GREEN}✅ R212 acknowledgment file found${NC}"
        echo "Recent acknowledgments:"
        tail -10 .r212-phase-acknowledged | sed 's/^/  /'
    else
        echo -e "${YELLOW}⚠️  No R212 acknowledgment file yet${NC}"
        echo "  Agents haven't acknowledged phase directories yet"
    fi
}

# Main validation loop
echo ""
echo "Scanning for phases..."
echo "════════════════════════════════════════════════════════════════"

# Create phase-plans directory if it doesn't exist (for validation)
mkdir -p phase-plans 2>/dev/null || true

# Check phases 1-5 (standard phases)
for PHASE in {1..5}; do
    ((TOTAL_PHASES++))
    
    echo ""
    echo "════════════════════════════════════════════════════════════════"
    echo "📁 Phase ${PHASE}"
    echo "════════════════════════════════════════════════════════════════"
    
    # Validate architecture plan
    ARCH_ERRORS=0
    if validate_phase_plan "phase-plans/PHASE-${PHASE}-ARCHITECTURE-PLAN.md"; then
        :  # Success
    else
        ARCH_ERRORS=$?
    fi
    
    # Validate implementation plan
    IMPL_ERRORS=0
    if validate_phase_plan "phase-plans/PHASE-${PHASE}-IMPLEMENTATION-PLAN.md"; then
        :  # Success
    else
        IMPL_ERRORS=$?
    fi
    
    # Validate directory structure
    STRUCT_ERRORS=0
    if validate_phase_structure "$PHASE"; then
        :  # Success
    else
        STRUCT_ERRORS=$?
    fi
    
    # Tally results
    TOTAL_ERRORS=$((ARCH_ERRORS + IMPL_ERRORS + STRUCT_ERRORS))
    
    if [ $TOTAL_ERRORS -eq 0 ]; then
        if [ -f "phase-plans/PHASE-${PHASE}-ARCHITECTURE-PLAN.md" ] || [ -f "phase-plans/PHASE-${PHASE}-IMPLEMENTATION-PLAN.md" ]; then
            echo ""
            echo -e "${GREEN}✅ PHASE VALID: All R212 requirements met${NC}"
            ((VALID_PHASES++))
        else
            echo ""
            echo -e "${YELLOW}⚠️  PHASE NOT STARTED: No plans created yet${NC}"
        fi
    else
        echo ""
        echo -e "${RED}❌ PHASE INVALID: $TOTAL_ERRORS errors found${NC}"
        ((INVALID_PHASES++))
    fi
done

# Check acknowledgments
check_acknowledgments

# Check for wave plans with metadata
echo ""
echo "Checking wave plans for R212 metadata..."
for wave_plan in phase-plans/PHASE-*-WAVE-*.md; do
    if [ -f "$wave_plan" ]; then
        if grep -q "PHASE INFRASTRUCTURE METADATA" "$wave_plan"; then
            echo -e "${GREEN}✅ $(basename $wave_plan) has R212 metadata${NC}"
        else
            echo -e "${YELLOW}⚠️  $(basename $wave_plan) missing R212 metadata${NC}"
            ((WARNINGS++))
        fi
    fi
done

# Summary report
echo ""
echo "════════════════════════════════════════════════════════════════"
echo "📊 R212 Phase Isolation Validation Summary"
echo "════════════════════════════════════════════════════════════════"
echo "Total Phases Checked: $TOTAL_PHASES"
echo -e "${GREEN}Valid Phases: $VALID_PHASES${NC}"
echo -e "${RED}Invalid Phases: $INVALID_PHASES${NC}"
echo -e "${YELLOW}Warnings: $WARNINGS${NC}"
echo ""

if [ $INVALID_PHASES -eq 0 ] && [ $WARNINGS -eq 0 ]; then
    echo -e "${GREEN}🎉 PROJECT_DONE: All phases comply with R212 isolation protocol!${NC}"
    exit 0
elif [ $INVALID_PHASES -eq 0 ]; then
    echo -e "${YELLOW}⚠️  MOSTLY COMPLIANT: No errors but $WARNINGS warnings found${NC}"
    echo ""
    echo "Recommended Actions:"
    echo "1. Orchestrator should inject R212 metadata into wave plans"
    echo "2. Agents should acknowledge phase directories when reading plans"
    exit 0
else
    echo -e "${RED}❌ VALIDATION FAILED: $INVALID_PHASES phases violate R212${NC}"
    echo ""
    echo "Required Actions:"
    echo "1. Orchestrator must inject metadata into all phase plans"
    echo "2. Agents must acknowledge phase directories"
    echo "3. All phase work must stay in designated directories"
    echo "4. No cross-phase contamination allowed"
    echo ""
    echo "Run orchestrator with R212 enforcement to fix these issues."
    exit 1
fi