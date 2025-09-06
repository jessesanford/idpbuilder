#!/bin/bash
# validate-state-completeness.sh - Validates state file has ALL items from plan
# Required by R281 - SUPREME LAW #7

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check arguments
if [ $# -ne 2 ]; then
    echo "Usage: $0 <implementation-plan.md> <orchestrator-state.yaml>"
    echo "Example: $0 PROJECT-IMPLEMENTATION-PLAN.md orchestrator-state.yaml"
    exit 1
fi

PLAN_FILE="$1"
STATE_FILE="$2"

# Verify files exist
if [ ! -f "$PLAN_FILE" ]; then
    echo -e "${RED}❌ Implementation plan not found: $PLAN_FILE${NC}"
    exit 1
fi

if [ ! -f "$STATE_FILE" ]; then
    echo -e "${RED}❌ State file not found: $STATE_FILE${NC}"
    exit 1
fi

echo "🔍 Validating state file completeness per R281..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Count phases in plan (look for Phase patterns)
PLAN_PHASES=$(grep -E "^#{1,3}\s+(Phase|PHASE)\s+[0-9]+" "$PLAN_FILE" | wc -l)
if [ $PLAN_PHASES -eq 0 ]; then
    # Try alternative patterns
    PLAN_PHASES=$(grep -E "^#{1,3}\s+.*Phase\s+[0-9]+" "$PLAN_FILE" | wc -l)
fi

# Count waves in plan (look for Wave patterns)
PLAN_WAVES=$(grep -E "^#{1,4}\s+(Wave|WAVE)\s+[0-9]+" "$PLAN_FILE" | wc -l)
if [ $PLAN_WAVES -eq 0 ]; then
    # Try alternative patterns
    PLAN_WAVES=$(grep -E "^#{1,4}\s+.*Wave\s+[0-9]+" "$PLAN_FILE" | wc -l)
fi

# Count efforts in plan (look for Effort patterns or numbered lists)
PLAN_EFFORTS=$(grep -E "^#{1,5}\s+(Effort|EFFORT)\s+[0-9]+" "$PLAN_FILE" | wc -l)
if [ $PLAN_EFFORTS -eq 0 ]; then
    # Try counting numbered implementation items
    PLAN_EFFORTS=$(grep -E "^\s*[0-9]+\.\s+\[" "$PLAN_FILE" | wc -l)
fi

# Count items in state file using yq
if ! command -v yq &> /dev/null; then
    echo -e "${YELLOW}⚠️ yq not found, using grep fallback${NC}"
    STATE_PHASES=$(grep -c "phase_id:" "$STATE_FILE" || true)
    STATE_WAVES=$(grep -c "wave_id:" "$STATE_FILE" || true)
    STATE_EFFORTS=$(grep -c "effort_id:" "$STATE_FILE" || true)
else
    STATE_PHASES=$(yq '.phases | length' "$STATE_FILE" 2>/dev/null || echo 0)
    STATE_WAVES=$(yq '.phases[].waves | length' "$STATE_FILE" 2>/dev/null | awk '{sum+=$1} END {print sum}')
    STATE_EFFORTS=$(yq '.phases[].waves[].efforts | length' "$STATE_FILE" 2>/dev/null | awk '{sum+=$1} END {print sum}')
fi

# Display results
echo "📊 PHASE COUNT VALIDATION:"
echo "   Plan has: $PLAN_PHASES phases"
echo "   State has: $STATE_PHASES phases"
if [ "$PLAN_PHASES" -eq "$STATE_PHASES" ] && [ "$PLAN_PHASES" -gt 0 ]; then
    echo -e "   ${GREEN}✅ Phase count matches!${NC}"
    PHASE_OK=1
else
    echo -e "   ${RED}❌ MISMATCH: Missing $(($PLAN_PHASES - $STATE_PHASES)) phases${NC}"
    PHASE_OK=0
fi

echo ""
echo "📊 WAVE COUNT VALIDATION:"
echo "   Plan has: $PLAN_WAVES waves"
echo "   State has: $STATE_WAVES waves"
if [ "$PLAN_WAVES" -eq "$STATE_WAVES" ] && [ "$PLAN_WAVES" -gt 0 ]; then
    echo -e "   ${GREEN}✅ Wave count matches!${NC}"
    WAVE_OK=1
else
    echo -e "   ${RED}❌ MISMATCH: Missing $(($PLAN_WAVES - $STATE_WAVES)) waves${NC}"
    WAVE_OK=0
fi

echo ""
echo "📊 EFFORT COUNT VALIDATION:"
echo "   Plan has: $PLAN_EFFORTS efforts"
echo "   State has: $STATE_EFFORTS efforts"
if [ "$PLAN_EFFORTS" -eq "$STATE_EFFORTS" ] && [ "$PLAN_EFFORTS" -gt 0 ]; then
    echo -e "   ${GREEN}✅ Effort count matches!${NC}"
    EFFORT_OK=1
else
    echo -e "   ${RED}❌ MISMATCH: Missing $(($PLAN_EFFORTS - $STATE_EFFORTS)) efforts${NC}"
    EFFORT_OK=0
fi

# Check for required fields in state file
echo ""
echo "📋 STRUCTURE VALIDATION:"
STRUCTURE_OK=1

# Check required top-level fields
for field in "current_phase" "current_wave" "current_state" "phases"; do
    if grep -q "^$field:" "$STATE_FILE"; then
        echo -e "   ${GREEN}✅ Found required field: $field${NC}"
    else
        echo -e "   ${RED}❌ Missing required field: $field${NC}"
        STRUCTURE_OK=0
    fi
done

# Final verdict
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
if [ $PHASE_OK -eq 1 ] && [ $WAVE_OK -eq 1 ] && [ $EFFORT_OK -eq 1 ] && [ $STRUCTURE_OK -eq 1 ]; then
    echo -e "${GREEN}🎉 VALIDATION PASSED! State file is COMPLETE per R281${NC}"
    echo "All phases, waves, and efforts from the plan are present."
    exit 0
else
    echo -e "${RED}🚨 VALIDATION FAILED! State file is INCOMPLETE${NC}"
    echo -e "${RED}R281 VIOLATION: -100% PENALTY${NC}"
    echo ""
    echo "Required actions:"
    [ $PHASE_OK -eq 0 ] && echo "  • Add missing phases from implementation plan"
    [ $WAVE_OK -eq 0 ] && echo "  • Add missing waves from implementation plan"
    [ $EFFORT_OK -eq 0 ] && echo "  • Add missing efforts from implementation plan"
    [ $STRUCTURE_OK -eq 0 ] && echo "  • Fix missing required fields"
    echo ""
    echo "Use templates/initial-state-template.yaml as reference"
    exit 1
fi