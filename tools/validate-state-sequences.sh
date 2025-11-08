#!/bin/bash
# tools/validate-state-sequences.sh
# Validates that state rules files document correct sequences from state machine
# Part of R517 enforcement - prevents documentation drift

set -euo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}🔍 State Sequence Documentation Validator${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
echo ""

cd "$PROJECT_ROOT"

STATE_MACHINE="state-machines/software-factory-3.0-state-machine.json"

if [ ! -f "$STATE_MACHINE" ]; then
    echo -e "${RED}❌ ERROR: State machine not found: $STATE_MACHINE${NC}"
    exit 1
fi

echo "State Machine: $STATE_MACHINE"
echo "Agent States Dir: agent-states/software-factory"
echo ""

MISMATCHES=0
FILES_CHECKED=0
STATES_WITHOUT_FILES=0

# Get list of all states in state machine
all_states=$(jq -r '.states | keys[]' "$STATE_MACHINE")

echo -e "${BLUE}Validating state sequence documentation...${NC}"
echo ""

for state in $all_states; do
    # Find corresponding rules file
    rules_file=$(find agent-states/software-factory -type f -path "*/$state/rules.md" 2>/dev/null | head -1)

    if [ -z "$rules_file" ]; then
        ((STATES_WITHOUT_FILES++))
        echo -e "${YELLOW}⚠️  State $state has no rules file (may be deprecated)${NC}"
        continue
    fi

    ((FILES_CHECKED++))

    # Get allowed transitions from state machine
    allowed_transitions=$(jq -r --arg state "$state" '.states[$state].allowed_transitions[]' "$STATE_MACHINE" 2>/dev/null || echo "")

    if [ -z "$allowed_transitions" ]; then
        echo -e "${YELLOW}⚠️  State $state has no allowed_transitions in state machine${NC}"
        continue
    fi

    # Get primary next state (first non-ERROR_RECOVERY transition)
    primary_next=$(echo "$allowed_transitions" | grep -v "ERROR_RECOVERY" | head -1 || echo "")

    if [ -z "$primary_next" ]; then
        # Only ERROR_RECOVERY available
        primary_next="ERROR_RECOVERY"
    fi

    # Check if rules file documents this next state
    if ! grep -q "$primary_next" "$rules_file" 2>/dev/null; then
        echo -e "${RED}❌ MISMATCH: $state${NC}"
        echo "   Rules file: $rules_file"
        echo "   Expected next state: $primary_next"
        echo "   But rules file does NOT mention this state!"
        echo ""
        ((MISMATCHES++))
        continue
    fi

    # Check for sequence diagram section
    if grep -q "YOUR POSITION IN THE MANDATORY SEQUENCE" "$rules_file" 2>/dev/null || \
       grep -q "POSITION IN SEQUENCE" "$rules_file" 2>/dev/null || \
       grep -q "State Sequence" "$rules_file" 2>/dev/null; then

        # Extract the documented next state from sequence diagram
        # Look for patterns like "↓ (MUST GO HERE NEXT)" or "→ NEXT_STATE"
        documented_next=$(grep -A 2 "👈 YOU ARE HERE" "$rules_file" 2>/dev/null | grep -oE '[A-Z_]{3,}' | grep -v "YOU\|ARE\|HERE\|MUST\|GO\|NEXT" | head -1 || echo "")

        if [ -n "$documented_next" ] && [ "$documented_next" != "$primary_next" ]; then
            echo -e "${RED}❌ SEQUENCE MISMATCH: $state${NC}"
            echo "   Rules file: $rules_file"
            echo "   State machine says next: $primary_next"
            echo "   Rules file documents: $documented_next"
            echo "   These MUST match!"
            echo ""
            ((MISMATCHES++))
        else
            echo -e "${GREEN}✅ $state → $primary_next (documented correctly)${NC}"
        fi
    else
        echo -e "${YELLOW}⚠️  $state: No sequence diagram found in rules file${NC}"
    fi
done

echo ""
echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}📊 Validation Results${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
echo ""
echo "Files checked: $FILES_CHECKED"
echo "Mismatches found: $MISMATCHES"
echo "States without rules files: $STATES_WITHOUT_FILES"
echo ""

if [ "$MISMATCHES" -eq 0 ]; then
    echo -e "${GREEN}✅ ALL SEQUENCES VALID${NC}"
    echo "State rules documentation matches state machine definitions."
    exit 0
else
    echo -e "${RED}❌ MISMATCHES DETECTED!${NC}"
    echo ""
    echo "State sequence documentation does NOT match state machine."
    echo "This can cause agents to skip mandatory states!"
    echo ""
    echo "Required actions:"
    echo "  1. Update rules files to match state machine"
    echo "  2. Verify state machine is correct"
    echo "  3. Check for documentation drift"
    echo ""
    echo "See: rule-library/R234-mandatory-state-traversal-supreme-law.md"
    exit 1
fi
