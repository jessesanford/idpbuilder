#!/bin/bash
# Comprehensive state name validation against SF 3.0 state machine
# Created: 2025-10-18
# Purpose: Find all invalid state name references across codebase

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
STATE_MACHINE="$PROJECT_ROOT/state-machines/software-factory-3.0-state-machine.json"

echo "================================================================"
echo "STATE NAME VALIDATION REPORT"
echo "================================================================"
echo "Timestamp: $(date)"
echo "State Machine: $STATE_MACHINE"
echo ""

# Extract valid state names
VALID_STATES=$(jq -r '.states | keys[]' "$STATE_MACHINE" | sort)
VALID_COUNT=$(echo "$VALID_STATES" | wc -l)

echo "✅ Valid states in SF 3.0 state machine: $VALID_COUNT"
echo ""

# Create temp file with valid states
VALID_STATES_FILE=$(mktemp)
echo "$VALID_STATES" > "$VALID_STATES_FILE"

# Find all state-like patterns in codebase
echo "🔍 Scanning codebase for state name references..."
echo ""

# Patterns to search for state names
STATE_PATTERNS=(
    'current_state.*:.*"([A-Z_]+)"'
    'to_state.*:.*"([A-Z_]+)"'
    'from_state.*:.*"([A-Z_]+)"'
    'transition_to\s+"([A-Z_]+)"'
    'state-manager\s+([A-Z_]+)'
)

# Search in specific file types
SEARCH_PATHS=(
    "agent-states/software-factory/orchestrator"
    ".claude/agents"
    ".claude/commands"
    "tests/fixtures"
    "rule-library"
)

# Collect all potential state names
POTENTIAL_STATES=$(mktemp)

for path in "${SEARCH_PATHS[@]}"; do
    if [ -d "$PROJECT_ROOT/$path" ]; then
        # Find state-like capitalized words (3+ uppercase letters with underscores)
        grep -roh '\b[A-Z][A-Z_]\{10,\}\b' "$PROJECT_ROOT/$path" 2>/dev/null | sort -u >> "$POTENTIAL_STATES" || true
    fi
done

# Remove duplicates
sort -u "$POTENTIAL_STATES" -o "$POTENTIAL_STATES"

# Check each potential state name
INVALID_STATES=()
INVALID_COUNT=0

while IFS= read -r state_name; do
    # Skip common non-state patterns
    if [[ "$state_name" =~ ^(CLAUDE_PROJECT_DIR|SOFTWARE_FACTORY|CONTINUE_SOFTWARE_FACTORY)$ ]]; then
        continue
    fi

    # Check if state is valid
    if ! grep -q "^${state_name}$" "$VALID_STATES_FILE"; then
        INVALID_STATES+=("$state_name")
        ((INVALID_COUNT++))
    fi
done < "$POTENTIAL_STATES"

# Report findings
if [ $INVALID_COUNT -eq 0 ]; then
    echo "✅ PROJECT_DONE: No invalid state names found!"
    echo ""
else
    echo "❌ FOUND $INVALID_COUNT INVALID STATE NAMES:"
    echo ""

    for invalid_state in "${INVALID_STATES[@]}"; do
        echo "  ❌ $invalid_state"

        # Find where this invalid state is referenced
        echo "     Referenced in:"
        grep -rl "$invalid_state" "$PROJECT_ROOT/agent-states/software-factory/orchestrator" \
            "$PROJECT_ROOT/.claude" \
            "$PROJECT_ROOT/tests/fixtures" 2>/dev/null | head -5 | while read -r file; do
            rel_path="${file#$PROJECT_ROOT/}"
            count=$(grep -c "$invalid_state" "$file" 2>/dev/null || echo "0")
            echo "       - $rel_path ($count occurrences)"
        done
        echo ""
    done
fi

# Validate state directories match state machine
echo "================================================================"
echo "STATE DIRECTORY VALIDATION"
echo "================================================================"
echo ""

ORCHESTRATOR_STATE_DIRS=$(find "$PROJECT_ROOT/agent-states/software-factory/orchestrator" \
    -maxdepth 1 -type d -not -name "orchestrator" -not -name "DEPRECATED" \
    | xargs -n1 basename | sort)

MISSING_DIRS=()
EXTRA_DIRS=()

# Check for states in machine without directories
while IFS= read -r state; do
    if [ ! -d "$PROJECT_ROOT/agent-states/software-factory/orchestrator/$state" ]; then
        MISSING_DIRS+=("$state")
    fi
done < "$VALID_STATES_FILE"

# Check for directories not in machine
while IFS= read -r dir; do
    if ! grep -q "^${dir}$" "$VALID_STATES_FILE"; then
        EXTRA_DIRS+=("$dir")
    fi
done <<< "$ORCHESTRATOR_STATE_DIRS"

if [ ${#MISSING_DIRS[@]} -gt 0 ]; then
    echo "⚠️  MISSING STATE DIRECTORIES (${#MISSING_DIRS[@]}):"
    for state in "${MISSING_DIRS[@]}"; do
        echo "  - $state"
    done
    echo ""
fi

if [ ${#EXTRA_DIRS[@]} -gt 0 ]; then
    echo "❌ EXTRA STATE DIRECTORIES NOT IN MACHINE (${#EXTRA_DIRS[@]}):"
    for dir in "${EXTRA_DIRS[@]}"; do
        echo "  - $dir"
    done
    echo ""
fi

if [ ${#MISSING_DIRS[@]} -eq 0 ] && [ ${#EXTRA_DIRS[@]} -eq 0 ]; then
    echo "✅ All state directories match state machine!"
    echo ""
fi

# Cleanup
rm -f "$VALID_STATES_FILE" "$POTENTIAL_STATES"

# Final summary
echo "================================================================"
echo "VALIDATION SUMMARY"
echo "================================================================"
echo "Valid states: $VALID_COUNT"
echo "Invalid state references: $INVALID_COUNT"
echo "Missing state directories: ${#MISSING_DIRS[@]}"
echo "Extra state directories: ${#EXTRA_DIRS[@]}"
echo ""

if [ $INVALID_COUNT -eq 0 ] && [ ${#EXTRA_DIRS[@]} -eq 0 ]; then
    echo "✅ VALIDATION PASSED - All state names valid!"
    exit 0
else
    echo "❌ VALIDATION FAILED - Fix invalid state names and directories!"
    exit 1
fi
