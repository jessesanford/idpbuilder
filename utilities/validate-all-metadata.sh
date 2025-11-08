#!/bin/bash

# validate-all-metadata.sh
# Comprehensive validation of R209 (Effort), R212 (Phase), and R213 (Wave) metadata
# Ensures orchestrator is properly acting as master of all directory structures

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
MAGENTA='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Counters
TOTAL_CHECKS=0
PASSED_CHECKS=0
FAILED_CHECKS=0
WARNINGS=0

echo "════════════════════════════════════════════════════════════════"
echo "🔍 COMPREHENSIVE METADATA VALIDATION"
echo "Checking R209 (Effort), R212 (Phase), R213 (Wave) Compliance"
echo "════════════════════════════════════════════════════════════════"
echo ""
echo "Orchestrator MUST be the master of ALL directory structures!"
echo ""

# Function to check if orchestrator is source
check_orchestrator_source() {
    local FILE="$1"
    local RULE="$2"
    
    if [ ! -f "$FILE" ]; then
        return 1
    fi
    
    if grep -q "METADATA_SOURCE.*ORCHESTRATOR" "$FILE" 2>/dev/null; then
        echo -e "${GREEN}✅ $RULE: Orchestrator is source${NC}"
        return 0
    elif grep -q "Added by Orchestrator" "$FILE" 2>/dev/null; then
        echo -e "${GREEN}✅ $RULE: Added by Orchestrator${NC}"
        return 0
    else
        echo -e "${YELLOW}⚠️  $RULE: Source not clearly marked${NC}"
        ((WARNINGS++))
        return 1
    fi
}

# Function to validate phase metadata (R212)
validate_phase_metadata() {
    local PHASE="$1"
    echo ""
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${CYAN}📁 PHASE ${PHASE} - R212 Validation${NC}"
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    
    local PHASE_ARCH="phase-plans/PHASE-${PHASE}-ARCHITECTURE-PLAN.md"
    local PHASE_IMPL="phase-plans/PHASE-${PHASE}-IMPLEMENTATION-PLAN.md"
    local CHECKS=0
    local FAILS=0
    
    for plan in "$PHASE_ARCH" "$PHASE_IMPL"; do
        if [ -f "$plan" ]; then
            echo "Checking: $(basename $plan)"
            ((TOTAL_CHECKS++))
            ((CHECKS++))
            
            # Check for R212 metadata
            if grep -q "PHASE INFRASTRUCTURE METADATA" "$plan"; then
                echo -e "${GREEN}✅ Has R212 phase metadata${NC}"
                ((PASSED_CHECKS++))
                
                # Verify it's from orchestrator
                check_orchestrator_source "$plan" "R212"
                
                # Extract and verify paths
                PHASE_ROOT=$(grep "**PHASE_ROOT**:" "$plan" 2>/dev/null | cut -d: -f2- | xargs || true)
                EFFORTS_ROOT=$(grep "**EFFORTS_ROOT**:" "$plan" 2>/dev/null | cut -d: -f2- | xargs || true)
                
                if [ -n "$PHASE_ROOT" ]; then
                    echo "  Phase root: $PHASE_ROOT"
                    [ -d "$PHASE_ROOT" ] && echo -e "${GREEN}  ✓ Directory exists${NC}" || echo -e "${YELLOW}  ⚠ Directory not created yet${NC}"
                fi
                
                if [ -n "$EFFORTS_ROOT" ]; then
                    echo "  Efforts root: $EFFORTS_ROOT"
                    [ -d "$EFFORTS_ROOT" ] && echo -e "${GREEN}  ✓ Directory exists${NC}" || echo -e "${YELLOW}  ⚠ Directory not created yet${NC}"
                fi
            else
                echo -e "${RED}❌ Missing R212 phase metadata${NC}"
                ((FAILED_CHECKS++))
                ((FAILS++))
            fi
        else
            echo -e "${YELLOW}⚠️  $(basename $plan) not created yet${NC}"
        fi
    done
    
    if [ $CHECKS -gt 0 ]; then
        if [ $FAILS -eq 0 ]; then
            echo -e "${GREEN}✅ Phase ${PHASE}: R212 compliant${NC}"
        else
            echo -e "${RED}❌ Phase ${PHASE}: R212 violations found${NC}"
        fi
    fi
}

# Function to validate wave metadata (R213)
validate_wave_metadata() {
    local PHASE="$1"
    local WAVE="$2"
    echo ""
    echo -e "${MAGENTA}────────────────────────────────────────────────────────${NC}"
    echo -e "${MAGENTA}🌊 Phase ${PHASE} Wave ${WAVE} - R213 Validation${NC}"
    echo -e "${MAGENTA}────────────────────────────────────────────────────────${NC}"
    
    local WAVE_ARCH="phase-plans/PHASE-${PHASE}-WAVE-${WAVE}-ARCHITECTURE-PLAN.md"
    local WAVE_IMPL="phase-plans/PHASE-${PHASE}-WAVE-${WAVE}-IMPLEMENTATION-PLAN.md"
    local CHECKS=0
    local FAILS=0
    
    for plan in "$WAVE_ARCH" "$WAVE_IMPL"; do
        if [ -f "$plan" ]; then
            echo "Checking: $(basename $plan)"
            ((TOTAL_CHECKS++))
            ((CHECKS++))
            
            # Check for R213 metadata
            if grep -q "WAVE INFRASTRUCTURE METADATA" "$plan"; then
                echo -e "${GREEN}✅ Has R213 wave metadata${NC}"
                ((PASSED_CHECKS++))
                
                # CRITICAL: Check if orchestrator is the source
                if grep -q "METADATA_SOURCE.*ORCHESTRATOR" "$plan"; then
                    echo -e "${GREEN}✅ ORCHESTRATOR is source (correct!)${NC}"
                else
                    echo -e "${RED}❌ ORCHESTRATOR NOT marked as source!${NC}"
                    echo -e "${RED}   This violates R213 - orchestrator MUST be master!${NC}"
                    ((FAILED_CHECKS++))
                    ((FAILS++))
                fi
                
                # Extract and verify paths
                WAVE_ROOT=$(grep "**WAVE_ROOT**:" "$plan" 2>/dev/null | cut -d: -f2- | xargs || true)
                EFFORT_COUNT=$(grep "**EFFORT_COUNT**:" "$plan" 2>/dev/null | cut -d: -f2- | xargs || true)
                
                if [ -n "$WAVE_ROOT" ]; then
                    echo "  Wave root: $WAVE_ROOT"
                    if [ -d "$WAVE_ROOT" ]; then
                        echo -e "${GREEN}  ✓ Directory exists${NC}"
                        
                        # Count actual efforts
                        ACTUAL_EFFORTS=$(ls -d ${WAVE_ROOT}/effort-* 2>/dev/null | wc -l || echo "0")
                        echo "  Metadata says: $EFFORT_COUNT efforts"
                        echo "  Actually found: $ACTUAL_EFFORTS efforts"
                        
                        if [ "$ACTUAL_EFFORTS" = "$EFFORT_COUNT" ]; then
                            echo -e "${GREEN}  ✓ Effort count matches${NC}"
                        else
                            echo -e "${YELLOW}  ⚠ Effort count mismatch${NC}"
                            ((WARNINGS++))
                        fi
                    else
                        echo -e "${YELLOW}  ⚠ Directory not created yet${NC}"
                    fi
                fi
            else
                echo -e "${RED}❌ Missing R213 wave metadata${NC}"
                ((FAILED_CHECKS++))
                ((FAILS++))
            fi
        fi
    done
    
    if [ $CHECKS -gt 0 ]; then
        if [ $FAILS -eq 0 ]; then
            echo -e "${GREEN}✅ Wave ${WAVE}: R213 compliant${NC}"
        else
            echo -e "${RED}❌ Wave ${WAVE}: R213 violations found${NC}"
        fi
    fi
}

# Function to validate effort metadata (R209)
validate_effort_metadata() {
    local PHASE="$1"
    local WAVE="$2"
    local EFFORT_DIR="$3"
    local EFFORT_NAME=$(basename "$EFFORT_DIR")
    
    if [ ! -d "$EFFORT_DIR" ]; then
        return 0
    fi
    
    echo "  Checking effort: $EFFORT_NAME"
    ((TOTAL_CHECKS++))
    
    local IMPL_PLAN="$EFFORT_DIR/IMPLEMENTATION-PLAN.md"
    
    if [ -f "$IMPL_PLAN" ]; then
        if grep -q "EFFORT INFRASTRUCTURE METADATA" "$IMPL_PLAN"; then
            echo -e "${GREEN}    ✅ Has R209 effort metadata${NC}"
            ((PASSED_CHECKS++))
            
            # Check consistency with wave metadata
            EFFORT_WORKING_DIR=$(grep "**WORKING_DIRECTORY**:" "$IMPL_PLAN" 2>/dev/null | cut -d: -f2- | xargs || true)
            
            if [[ "$EFFORT_WORKING_DIR" == *"efforts/phase${PHASE}/wave${WAVE}"* ]]; then
                echo -e "${GREEN}    ✓ Path consistent with wave structure${NC}"
            else
                echo -e "${RED}    ❌ Path inconsistent with wave!${NC}"
                ((FAILED_CHECKS++))
            fi
        else
            echo -e "${RED}    ❌ Missing R209 effort metadata${NC}"
            ((FAILED_CHECKS++))
        fi
    else
        echo -e "${YELLOW}    ⚠ No implementation plan yet${NC}"
    fi
}

# Function to check metadata hierarchy consistency
check_hierarchy_consistency() {
    echo ""
    echo -e "${BLUE}════════════════════════════════════════════════════════${NC}"
    echo -e "${BLUE}🔗 METADATA HIERARCHY CONSISTENCY CHECK${NC}"
    echo -e "${BLUE}════════════════════════════════════════════════════════${NC}"
    
    echo ""
    echo "Checking that all metadata forms a consistent hierarchy:"
    echo "  ORCHESTRATOR (Master)"
    echo "    ├── R212: Phase Metadata"
    echo "    │   └── R213: Wave Metadata"
    echo "    │       └── R209: Effort Metadata"
    echo ""
    
    local CONSISTENT=true
    
    # For each phase with metadata
    for phase_plan in phase-plans/PHASE-*-ARCHITECTURE-PLAN.md; do
        if [ -f "$phase_plan" ] && grep -q "PHASE INFRASTRUCTURE METADATA" "$phase_plan"; then
            PHASE_NUM=$(grep "**PHASE_NUMBER**:" "$phase_plan" | cut -d: -f2- | xargs)
            PHASE_EFFORTS_ROOT=$(grep "**EFFORTS_ROOT**:" "$phase_plan" | cut -d: -f2- | xargs)
            
            echo "Phase $PHASE_NUM metadata found"
            
            # Check wave metadata consistency
            for wave_plan in phase-plans/PHASE-${PHASE_NUM}-WAVE-*-ARCHITECTURE-PLAN.md; do
                if [ -f "$wave_plan" ] && grep -q "WAVE INFRASTRUCTURE METADATA" "$wave_plan"; then
                    WAVE_NUM=$(basename "$wave_plan" | sed 's/.*WAVE-\([0-9]*\).*/\1/')
                    WAVE_ROOT=$(grep "**WAVE_ROOT**:" "$wave_plan" | cut -d: -f2- | xargs)
                    
                    # Check wave root is under phase efforts root
                    if [[ "$WAVE_ROOT" == "$PHASE_EFFORTS_ROOT/wave${WAVE_NUM}" ]]; then
                        echo -e "${GREEN}  ✓ Wave $WAVE_NUM path consistent with phase${NC}"
                    else
                        echo -e "${RED}  ❌ Wave $WAVE_NUM path INCONSISTENT!${NC}"
                        echo "     Expected: $PHASE_EFFORTS_ROOT/wave${WAVE_NUM}"
                        echo "     Found: $WAVE_ROOT"
                        CONSISTENT=false
                    fi
                fi
            done
        fi
    done
    
    if $CONSISTENT; then
        echo -e "${GREEN}✅ Metadata hierarchy is CONSISTENT${NC}"
    else
        echo -e "${RED}❌ Metadata hierarchy has INCONSISTENCIES${NC}"
        echo -e "${RED}   Orchestrator must ensure all paths align!${NC}"
    fi
}

# Main validation
echo "Starting comprehensive validation..."

# Check phase-plans directory
if [ ! -d "phase-plans" ]; then
    echo -e "${YELLOW}Creating phase-plans directory for validation...${NC}"
    mkdir -p phase-plans
fi

# Validate all phases (1-5)
for PHASE in {1..5}; do
    validate_phase_metadata "$PHASE"
    
    # Check for waves in this phase
    for WAVE in {1..6}; do
        if ls phase-plans/PHASE-${PHASE}-WAVE-${WAVE}-*.md 2>/dev/null | head -1 > /dev/null; then
            validate_wave_metadata "$PHASE" "$WAVE"
            
            # Check efforts in this wave
            WAVE_ROOT="efforts/phase${PHASE}/wave${WAVE}"
            if [ -d "$WAVE_ROOT" ]; then
                echo "  Checking efforts in $WAVE_ROOT..."
                for effort_dir in ${WAVE_ROOT}/effort-*; do
                    if [ -d "$effort_dir" ]; then
                        validate_effort_metadata "$PHASE" "$WAVE" "$effort_dir"
                    fi
                done
            fi
        fi
    done
done

# Check hierarchy consistency
check_hierarchy_consistency

# Check for acknowledgment files
echo ""
echo -e "${BLUE}════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}📋 ACKNOWLEDGMENT FILES${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════${NC}"

[ -f ".r209-acknowledged" ] && echo -e "${GREEN}✅ R209 acknowledgments found${NC}" || echo -e "${YELLOW}⚠️  No R209 acknowledgments${NC}"
[ -f ".r212-phase-acknowledged" ] && echo -e "${GREEN}✅ R212 acknowledgments found${NC}" || echo -e "${YELLOW}⚠️  No R212 acknowledgments${NC}"

# Final summary
echo ""
echo "════════════════════════════════════════════════════════════════"
echo "📊 FINAL VALIDATION SUMMARY"
echo "════════════════════════════════════════════════════════════════"
echo "Total Checks: $TOTAL_CHECKS"
echo -e "${GREEN}Passed: $PASSED_CHECKS${NC}"
echo -e "${RED}Failed: $FAILED_CHECKS${NC}"
echo -e "${YELLOW}Warnings: $WARNINGS${NC}"
echo ""

if [ $FAILED_CHECKS -eq 0 ] && [ $WARNINGS -eq 0 ]; then
    echo -e "${GREEN}🎉 PERFECT COMPLIANCE!${NC}"
    echo "The orchestrator is properly acting as master of all directory structures!"
    echo "All metadata (R209, R212, R213) is correctly injected and consistent!"
    exit 0
elif [ $FAILED_CHECKS -eq 0 ]; then
    echo -e "${YELLOW}⚠️  MOSTLY COMPLIANT with $WARNINGS warnings${NC}"
    echo ""
    echo "Recommendations:"
    echo "1. Ensure all metadata explicitly marks ORCHESTRATOR as source"
    echo "2. Create missing directories to match metadata"
    echo "3. Keep effort counts synchronized"
    exit 0
else
    echo -e "${RED}❌ COMPLIANCE FAILURES DETECTED!${NC}"
    echo ""
    echo "CRITICAL REQUIREMENTS:"
    echo "1. Orchestrator MUST inject metadata at ALL levels (phase, wave, effort)"
    echo "2. All metadata MUST mark ORCHESTRATOR as source"
    echo "3. Directory structures MUST match metadata exactly"
    echo "4. Agents MUST follow orchestrator-defined paths"
    echo ""
    echo "The orchestrator is NOT properly acting as master of directory structures!"
    exit 1
fi