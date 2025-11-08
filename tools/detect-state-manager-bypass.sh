#!/bin/bash
# tools/detect-state-manager-bypass.sh
# Detects State Manager bypass violations in git history
# Part of R517 enforcement - Supreme Law

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}🔍 State Manager Bypass Detection Tool${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
echo ""
echo "Scanning git history for State Manager bypass violations..."
echo "Rule: R517 - Universal State Manager Consultation Law"
echo ""

cd "$PROJECT_ROOT"

VIOLATIONS_FOUND=0
COMMITS_SCANNED=0

# Function to check a single commit
check_commit() {
    local commit_hash="$1"
    local commit_msg
    commit_msg=$(git log -1 --format='%s' "$commit_hash")

    ((COMMITS_SCANNED++))

    # Check if commit modifies state files
    local state_files_modified
    state_files_modified=$(git show --name-only --format="" "$commit_hash" | grep -E 'orchestrator-state.*\.json$|bug-tracking\.json$|integration-containers\.json$|fix-cascade-state\.json$' || true)

    if [ -z "$state_files_modified" ]; then
        return 0  # No state files modified, skip
    fi

    echo -e "${YELLOW}📋 Checking commit: $commit_hash${NC}"
    echo "   Message: $commit_msg"
    echo "   Modified state files:"
    echo "$state_files_modified" | sed 's/^/     - /'

    # VIOLATION CHECK 1: Commit message lacks State Manager reference
    if ! echo "$commit_msg" | grep -qi "state-manager\|State Manager\|\[R288\]"; then
        echo -e "${RED}   ❌ VIOLATION: Commit message lacks State Manager reference${NC}"
        ((VIOLATIONS_FOUND++))
    fi

    # VIOLATION CHECK 2: Check validated_by field in diff
    for file in $state_files_modified; do
        if [ -f "$file" ]; then
            # Check current validated_by
            local validated_by
            validated_by=$(jq -r '.state_machine.state_history[0].validated_by // "missing"' "$file" 2>/dev/null || echo "error")

            if [ "$validated_by" != "state-manager" ] && [ "$validated_by" != "missing" ]; then
                echo -e "${RED}   ❌ VIOLATION: $file has validated_by=$validated_by (not 'state-manager')${NC}"
                ((VIOLATIONS_FOUND++))
            fi
        fi
    done

    # VIOLATION CHECK 3: Direct jq/sed usage in commit
    local forbidden_patterns
    forbidden_patterns=$(git show "$commit_hash" | grep -E '^\+.*jq.*orchestrator-state|^\+.*sed.*orchestrator-state|^\+.*yq.*orchestrator-state' || true)

    if [ -n "$forbidden_patterns" ]; then
        echo -e "${RED}   ❌ VIOLATION: Direct state file manipulation detected:${NC}"
        echo "$forbidden_patterns" | sed 's/^/     /'
        ((VIOLATIONS_FOUND++))
    fi

    # VIOLATION CHECK 4: Check if state_history entry exists for transition
    if [ -f "orchestrator-state-v3.json" ]; then
        local current_state
        current_state=$(jq -r '.state_machine.current_state' orchestrator-state-v3.json 2>/dev/null || echo "unknown")
        local previous_state
        previous_state=$(jq -r '.state_machine.previous_state' orchestrator-state-v3.json 2>/dev/null || echo "unknown")

        # Check if this transition is in state_history
        local history_entry
        history_entry=$(jq --arg from "$previous_state" --arg to "$current_state" '.state_machine.state_history[] | select(.from_state == $from and .to_state == $to) | .validated_by' orchestrator-state-v3.json 2>/dev/null || echo "")

        if [ -z "$history_entry" ] && [ "$current_state" != "unknown" ]; then
            echo -e "${RED}   ❌ VIOLATION: No state_history entry for transition $previous_state → $current_state${NC}"
            ((VIOLATIONS_FOUND++))
        fi
    fi

    echo ""
}

# Scan last N commits (default 100, or specify as argument)
NUM_COMMITS="${1:-100}"

echo "Scanning last $NUM_COMMITS commits..."
echo ""

# Get list of commits
commits=$(git log -"$NUM_COMMITS" --format='%H')

for commit in $commits; do
    check_commit "$commit"
done

echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}📊 Scan Results${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
echo ""
echo "Commits scanned: $COMMITS_SCANNED"
echo "Violations found: $VIOLATIONS_FOUND"
echo ""

if [ "$VIOLATIONS_FOUND" -eq 0 ]; then
    echo -e "${GREEN}✅ NO VIOLATIONS DETECTED${NC}"
    echo "All state transitions appear to use State Manager correctly."
    exit 0
else
    echo -e "${RED}❌ VIOLATIONS DETECTED!${NC}"
    echo ""
    echo "State Manager bypass violations were found in git history."
    echo "This indicates system corruption risk."
    echo ""
    echo "Recommended actions:"
    echo "  1. Review flagged commits"
    echo "  2. Verify state file consistency"
    echo "  3. Check state_history completeness"
    echo "  4. Consider rollback if recent"
    echo ""
    echo "See: rule-library/R517-universal-state-manager-consultation-law.md"
    exit 1
fi
