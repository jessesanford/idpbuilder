#!/bin/bash

# detect-orphaned-states.sh - Detect and report orphaned states in the state machine
# Part of R289: Orphaned State Detection and Prevention

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
STATE_MACHINE="$PROJECT_ROOT/state-machines/software-factory-3.0-state-machine.json"

# Terminal states that are expected to have no outgoing transitions
TERMINAL_STATES="PROJECT_DONE|ERROR_RECOVERY|COMPLETED|BLOCKED|DECISION"

echo "========================================="
echo "   ORPHANED STATE DETECTION UTILITY"
echo "========================================="
echo ""
echo "Checking: $STATE_MACHINE"
echo ""

# Check if state machine file exists
if [ ! -f "$STATE_MACHINE" ]; then
    echo "❌ ERROR: State machine file not found!"
    echo "Expected at: $STATE_MACHINE"
    exit 1
fi

# Create temp directory for analysis
TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT

# Extract all defined states
echo "Step 1: Extracting all defined states..."
grep '^- \*\*[A-Z_]*\*\*' "$STATE_MACHINE" | \
    sed 's/- \*\*//' | \
    sed 's/\*\*.*//' | \
    sort -u > "$TEMP_DIR/all_states.txt"

TOTAL_STATES=$(wc -l < "$TEMP_DIR/all_states.txt")
echo "  Found $TOTAL_STATES total states"

# Extract states with transitions
echo ""
echo "Step 2: Finding states with transitions..."
grep -E '→' "$STATE_MACHINE" | \
    sed 's/→/\n/g' | \
    grep -E '^[A-Z_]+' | \
    sed 's/[^A-Z_].*//' | \
    sort -u > "$TEMP_DIR/states_with_transitions.txt"

STATES_WITH_TRANSITIONS=$(wc -l < "$TEMP_DIR/states_with_transitions.txt")
echo "  Found $STATES_WITH_TRANSITIONS states with transitions"

# Find orphaned states (excluding terminal states)
echo ""
echo "Step 3: Identifying orphaned states..."
comm -23 "$TEMP_DIR/all_states.txt" "$TEMP_DIR/states_with_transitions.txt" > "$TEMP_DIR/orphaned_candidates.txt"

# Filter out terminal states
grep -v -E "^($TERMINAL_STATES)$" "$TEMP_DIR/orphaned_candidates.txt" > "$TEMP_DIR/orphaned_states.txt" 2>/dev/null || true

ORPHANED_COUNT=$(wc -l < "$TEMP_DIR/orphaned_states.txt")

echo ""
echo "========================================="
echo "              RESULTS"
echo "========================================="
echo ""

if [ "$ORPHANED_COUNT" -eq 0 ]; then
    echo "✅ PROJECT_DONE: No orphaned states detected!"
    echo ""
    echo "Summary:"
    echo "  - Total states: $TOTAL_STATES"
    echo "  - States with transitions: $STATES_WITH_TRANSITIONS"
    echo "  - Terminal states (expected): $(comm -12 "$TEMP_DIR/orphaned_candidates.txt" <(echo -e "PROJECT_DONE\nERROR_RECOVERY\nCOMPLETED\nBLOCKED\nDECISION") | wc -l)"
    echo "  - Orphaned states: 0"
else
    echo "⚠️ WARNING: Found $ORPHANED_COUNT orphaned state(s)!"
    echo ""
    echo "Orphaned States (non-terminal with no transitions):"
    echo "------------------------------------------------"
    while IFS= read -r state; do
        echo "  ❌ $state"
        
        # Check for state directory
        for agent in orchestrator sw-engineer code-reviewer architect integration; do
            if [ -d "$PROJECT_ROOT/agent-states/$agent/$state" ]; then
                echo "     └─ Directory exists: agent-states/$agent/$state/"
            fi
        done
        
        # Check for any references in the codebase
        ref_count=$(grep -r "$state" "$PROJECT_ROOT" --include="*.md" --include="*.yaml" 2>/dev/null | wc -l || echo "0")
        echo "     └─ Found $ref_count reference(s) in codebase"
        echo ""
    done < "$TEMP_DIR/orphaned_states.txt"
    
    echo "Recommended Actions:"
    echo "-------------------"
    echo "1. Remove orphaned states from software-factory-3.0-state-machine.json"
    echo "2. Archive state directories with .DEPRECATED-$(date +%Y%m%d) suffix"
    echo "3. Update any references to use valid states"
    echo "4. Create migration guide if states were previously used"
    echo ""
    echo "See R289-orphaned-state-detection.md for detailed procedures"
fi

echo ""
echo "========================================="
echo "         TRANSITION ANALYSIS"
echo "========================================="
echo ""

# Check for states with only incoming transitions (dead ends)
echo "Checking for dead-end states (non-terminal with no outgoing transitions)..."
while IFS= read -r state; do
    # Skip if it's a terminal state
    if echo "$state" | grep -qE "^($TERMINAL_STATES)$"; then
        continue
    fi
    
    # Check if state has outgoing transitions
    if ! grep -qE "^$state →" "$STATE_MACHINE"; then
        # Check if it has incoming transitions
        if grep -qE "→ $state" "$STATE_MACHINE"; then
            echo "  ⚠️ $state - has incoming but no outgoing transitions (dead end)"
        fi
    fi
done < "$TEMP_DIR/all_states.txt"

echo ""
echo "Checking for unreachable states (no incoming transitions)..."
while IFS= read -r state; do
    # Skip INIT (entry point) and terminal states for outgoing
    if [ "$state" = "INIT" ]; then
        continue
    fi
    
    # Check if state has incoming transitions
    if ! grep -qE "→ $state" "$STATE_MACHINE"; then
        # Check if it has outgoing transitions
        if grep -qE "^$state →" "$STATE_MACHINE"; then
            echo "  ⚠️ $state - has outgoing but no incoming transitions (unreachable)"
        fi
    fi
done < "$TEMP_DIR/all_states.txt"

echo ""
echo "========================================="
echo "     STATE DIRECTORY CONSISTENCY"
echo "========================================="
echo ""

# Check for state directories without corresponding states
echo "Checking for orphaned state directories..."
for agent_dir in "$PROJECT_ROOT/agent-states"/*; do
    if [ ! -d "$agent_dir" ]; then
        continue
    fi
    
    agent=$(basename "$agent_dir")
    
    for state_dir in "$agent_dir"/*; do
        if [ ! -d "$state_dir" ]; then
            continue
        fi
        
        state=$(basename "$state_dir")
        
        # Skip deprecated directories
        if echo "$state" | grep -qE "DEPRECATED|backup"; then
            continue
        fi
        
        # Check if state exists in state machine
        if ! grep -qE "^\- \*\*$state\*\*" "$STATE_MACHINE"; then
            echo "  ⚠️ Directory without state: agent-states/$agent/$state/"
        fi
    done
done

echo ""
echo "========================================="
echo "           VALIDATION COMPLETE"
echo "========================================="

# Exit with error if orphaned states found
if [ "$ORPHANED_COUNT" -gt 0 ]; then
    exit 1
fi

exit 0