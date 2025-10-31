#!/bin/bash
# validate-split-spawning.sh - Prevent parallel split agent spawning violations
# Enforces R202 - SINGLE agent for ALL splits, SEQUENTIAL execution only

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Locate the project root
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# State file location
STATE_FILE="${1:-$PROJECT_ROOT/orchestrator-state-v3.json}"

echo "=================================================="
echo "🔍 SPLIT SPAWNING VALIDATION (R202 Enforcement)"
echo "=================================================="
echo "State file: $STATE_FILE"
echo ""

# Check if state file exists
if [ ! -f "$STATE_FILE" ]; then
    echo -e "${RED}❌ ERROR: State file not found: $STATE_FILE${NC}"
    exit 1
fi

# Function to validate split spawning
validate_split_spawning() {
    local violations=0

    # Check efforts_in_progress for split violations
    echo "Checking efforts in progress for split violations..."

    # Extract all split-related efforts
    local split_efforts=$(jq -r '.efforts_in_progress | to_entries[] | select(.key | contains("split")) | .key' "$STATE_FILE" 2>/dev/null || true)

    if [ -n "$split_efforts" ]; then
        echo ""
        echo "Found split efforts in progress:"
        echo "$split_efforts"
        echo ""

        # Group by effort base name to detect parallel splits
        declare -A effort_splits

        while IFS= read -r split_effort; do
            # Extract base effort name (everything before -split- or _split_)
            local base_effort=$(echo "$split_effort" | sed -E 's/[-_]split[-_][0-9]+.*//')

            # Count splits for this base effort
            if [ -n "${effort_splits[$base_effort]:-}" ]; then
                effort_splits[$base_effort]=$((effort_splits[$base_effort] + 1))
            else
                effort_splits[$base_effort]=1
            fi
        done <<< "$split_efforts"

        # Check for violations
        for effort in "${!effort_splits[@]}"; do
            local split_count="${effort_splits[$effort]}"
            if [ "$split_count" -gt 1 ]; then
                echo -e "${RED}🔴🔴🔴 CRITICAL R202 VIOLATION DETECTED! 🔴🔴🔴${NC}"
                echo -e "${RED}Effort '$effort' has $split_count splits in parallel!${NC}"
                echo -e "${RED}This is FORBIDDEN - splits must be sequential!${NC}"
                violations=$((violations + 1))
            fi
        done
    fi

    # Check spawned_agents for parallel split agents
    echo "Checking spawned agents for parallel split execution..."

    local spawned_at_same_time=$(jq -r '
        .spawned_agents[]? |
        select(.agent_id | contains("split")) |
        "\(.agent_id)|\(.spawned_at)"
    ' "$STATE_FILE" 2>/dev/null || true)

    if [ -n "$spawned_at_same_time" ]; then
        echo ""
        echo "Split-related agents found:"

        # Group by timestamp to find parallel spawning
        declare -A timestamp_agents

        while IFS='|' read -r agent_id timestamp; do
            if [ -n "$timestamp" ]; then
                if [ -n "${timestamp_agents[$timestamp]:-}" ]; then
                    # Multiple agents at same timestamp - check if they're splits of same effort
                    local existing="${timestamp_agents[$timestamp]}"

                    # Extract base efforts
                    local base1=$(echo "$existing" | sed -E 's/.*[-_]([^-_]+)[-_]split[-_].*/\1/')
                    local base2=$(echo "$agent_id" | sed -E 's/.*[-_]([^-_]+)[-_]split[-_].*/\1/')

                    if [ "$base1" = "$base2" ]; then
                        echo -e "${RED}🔴🔴🔴 R202 VIOLATION: Parallel split agents detected! 🔴🔴🔴${NC}"
                        echo -e "${RED}Agents spawned at same time ($timestamp):${NC}"
                        echo -e "${RED}  - $existing${NC}"
                        echo -e "${RED}  - $agent_id${NC}"
                        violations=$((violations + 1))
                    fi
                else
                    timestamp_agents[$timestamp]="$agent_id"
                fi
            fi
        done <<< "$spawned_at_same_time"
    fi

    # Check split_tracking for proper sequential execution
    echo ""
    echo "Checking split_tracking for sequential execution..."

    local split_tracking=$(jq -r '.split_tracking | to_entries[]?' "$STATE_FILE" 2>/dev/null || true)

    if [ -n "$split_tracking" ]; then
        local effort_name=$(echo "$split_tracking" | jq -r '.key')
        local execution_mode=$(echo "$split_tracking" | jq -r '.value.execution_mode // "UNKNOWN"')
        local in_progress=$(echo "$split_tracking" | jq -r '.value.splits_in_progress | length')

        if [ "$execution_mode" = "PARALLEL" ]; then
            echo -e "${RED}🔴 R202 VIOLATION: Split execution mode is PARALLEL for $effort_name!${NC}"
            echo -e "${RED}Must be SEQUENTIAL per R202${NC}"
            violations=$((violations + 1))
        fi

        if [ "$in_progress" -gt 1 ]; then
            echo -e "${RED}🔴 R202 VIOLATION: Multiple splits in progress for $effort_name!${NC}"
            echo -e "${RED}Only ONE split at a time is allowed${NC}"
            violations=$((violations + 1))
        fi
    fi

    return $violations
}

# Function to check for recursive splits (R511)
check_recursive_splits() {
    echo ""
    echo "Checking for recursive splits (R511 enforcement)..."

    local recursive_found=0

    # Look for split-of-split patterns
    local nested_splits=$(jq -r '
        .efforts_in_progress |
        to_entries[] |
        select(.key | test("split.*split"; "i")) |
        .key
    ' "$STATE_FILE" 2>/dev/null || true)

    if [ -n "$nested_splits" ]; then
        echo -e "${RED}🔴🔴🔴 R511 VIOLATION: RECURSIVE SPLITS DETECTED! 🔴🔴🔴${NC}"
        echo -e "${RED}These are splits of splits (ABSOLUTELY FORBIDDEN):${NC}"
        echo "$nested_splits" | while read -r nested; do
            echo -e "${RED}  - $nested${NC}"
        done
        recursive_found=1
    fi

    return $recursive_found
}

# Main validation
echo "=================================================="
echo "VALIDATION RESULTS:"
echo "=================================================="

violations=0

# Run split spawning validation
if ! validate_split_spawning; then
    violations=$((violations + $?))
fi

# Run recursive split check
if ! check_recursive_splits; then
    violations=$((violations + 1))
fi

# Final report
echo ""
echo "=================================================="
if [ $violations -eq 0 ]; then
    echo -e "${GREEN}✅ VALIDATION PASSED${NC}"
    echo -e "${GREEN}No split spawning violations detected${NC}"
    echo -e "${GREEN}Compliant with R202 and R511${NC}"
else
    echo -e "${RED}❌ VALIDATION FAILED${NC}"
    echo -e "${RED}Found $violations violation(s)${NC}"
    echo ""
    echo -e "${RED}REQUIRED ACTIONS:${NC}"
    echo -e "${RED}1. STOP all split work immediately${NC}"
    echo -e "${RED}2. Kill parallel split agents${NC}"
    echo -e "${RED}3. Restart with SINGLE agent for ALL splits${NC}"
    echo -e "${RED}4. Ensure SEQUENTIAL execution only${NC}"
    echo ""
    echo -e "${RED}Remember: Parallel splits = -100% AUTOMATIC FAILURE${NC}"
    exit 1
fi

echo "=================================================="

# Provide guidance for correct implementation
echo ""
echo "📚 CORRECT SPLIT IMPLEMENTATION:"
echo "1. ONE Code Reviewer creates ALL split plans"
echo "2. Orchestrator creates infrastructure for split-001 ONLY"
echo "3. ONE SW Engineer implements split-001"
echo "4. Wait for split-001 review to PASS"
echo "5. THEN create split-002 infrastructure"
echo "6. SAME SW Engineer implements split-002"
echo "7. Continue sequentially for all splits"
echo ""
echo "See: rule-library/R202-single-agent-per-split.md"
echo "See: rule-library/R511-absolute-prohibition-recursive-splits.md"