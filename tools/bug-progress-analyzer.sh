#!/bin/bash

# Bug Progress Analyzer Tool
# Implements R615 (Progress-Based Iteration Limits) and R616 (Bug Lifecycle Tracking)
#
# Purpose: Analyze bug fix progress across iterations to determine if system is making
#          actual progress (bugs fixed) or just churning (same bugs repeated).
#
# Usage:
#   bash tools/bug-progress-analyzer.sh analyze_progress SCOPE CURRENT_ITER PREVIOUS_ITER
#   bash tools/bug-progress-analyzer.sh count_closed_bugs SCOPE CURRENT_ITER
#   bash tools/bug-progress-analyzer.sh detect_reopened_bugs SCOPE CURRENT_ITER PREVIOUS_ITER
#   bash tools/bug-progress-analyzer.sh calculate_progress_score SCOPE CURRENT_ITER PREVIOUS_ITER
#   bash tools/bug-progress-analyzer.sh should_continue_or_escalate SCOPE CURRENT_ITER

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get project directory
PROJECT_DIR="${CLAUDE_PROJECT_DIR:-$(pwd)}"
STATE_FILE="${PROJECT_DIR}/orchestrator-state-v3.json"

# ============================================================================
# FUNCTION: get_bug_tracking_file
# Get the bug-tracking.json file for a given scope
# ============================================================================
get_bug_tracking_file() {
    local scope="$1"

    case "$scope" in
        WAVE)
            # Bug tracking in wave integration workspace
            wave_num=$(jq -r '.project_progression.current_wave.wave_number' "$STATE_FILE")
            echo "${PROJECT_DIR}/integration-workspaces/wave${wave_num}/bug-tracking.json"
            ;;
        PHASE)
            # Bug tracking in phase integration workspace
            phase_num=$(jq -r '.project_progression.current_phase.phase_number' "$STATE_FILE")
            echo "${PROJECT_DIR}/integration-workspaces/phase${phase_num}/bug-tracking.json"
            ;;
        PROJECT)
            # Bug tracking in project integration workspace
            echo "${PROJECT_DIR}/integration-workspaces/project/bug-tracking.json"
            ;;
        *)
            echo "ERROR: Invalid scope: $scope" >&2
            exit 1
            ;;
    esac
}

# ============================================================================
# FUNCTION: count_closed_bugs
# Count how many bugs were closed in a specific iteration
# ============================================================================
count_closed_bugs() {
    local scope="$1"
    local iteration="$2"

    local bug_file
    bug_file=$(get_bug_tracking_file "$scope")

    if [ ! -f "$bug_file" ]; then
        echo "0"
        return
    fi

    # Count bugs with closed_iteration = $iteration
    jq -r --argjson iter "$iteration" \
        '.bugs[] | select(.lifecycle.closed_iteration == $iter) | .bug_id' \
        "$bug_file" | wc -l
}

# ============================================================================
# FUNCTION: count_open_bugs
# Count how many bugs are open in a specific iteration
# ============================================================================
count_open_bugs() {
    local scope="$1"
    local iteration="$2"

    local bug_file
    bug_file=$(get_bug_tracking_file "$scope")

    if [ ! -f "$bug_file" ]; then
        echo "0"
        return
    fi

    # Count bugs that are OPEN or REOPENED in this iteration
    # (first_found_iteration <= iteration AND (not closed OR closed_iteration < iteration))
    jq -r --argjson iter "$iteration" '
        .bugs[] |
        select(
            .lifecycle.first_found_iteration <= $iter and
            (
                .lifecycle.closed_iteration == null or
                .lifecycle.closed_iteration < $iter
            )
        ) | .bug_id
    ' "$bug_file" | wc -l
}

# ============================================================================
# FUNCTION: detect_reopened_bugs
# Detect bugs that were closed in previous iteration but are open again
# ============================================================================
detect_reopened_bugs() {
    local scope="$1"
    local current_iteration="$2"
    local previous_iteration="$3"

    local bug_file
    bug_file=$(get_bug_tracking_file "$scope")

    if [ ! -f "$bug_file" ]; then
        echo "[]"
        return
    fi

    # Find bugs closed in previous iteration
    local closed_in_previous
    closed_in_previous=$(jq -r --argjson prev "$previous_iteration" \
        '.bugs[] | select(.lifecycle.closed_iteration == $prev) | .bug_id' \
        "$bug_file" | sort)

    # Find bugs open in current iteration
    local open_in_current
    open_in_current=$(jq -r --argjson curr "$current_iteration" \
        '.bugs[] |
         select(
             .lifecycle.first_found_iteration <= $curr and
             (
                 .lifecycle.closed_iteration == null or
                 .lifecycle.closed_iteration < $curr
             )
         ) | .bug_id' \
        "$bug_file" | sort)

    # Find intersection (reopened bugs)
    local reopened
    if [ -n "$closed_in_previous" ] && [ -n "$open_in_current" ]; then
        reopened=$(comm -12 <(echo "$closed_in_previous") <(echo "$open_in_current"))
    else
        reopened=""
    fi

    # Return as JSON array
    if [ -n "$reopened" ]; then
        echo "$reopened" | jq -R . | jq -s .
    else
        echo "[]"
    fi
}

# ============================================================================
# FUNCTION: calculate_progress_category
# Determine progress category based on bugs closed and net change
# ============================================================================
calculate_progress_category() {
    local bugs_closed="$1"
    local net_change="$2"

    if [ "$bugs_closed" -eq 0 ]; then
        if [ "$net_change" -eq 0 ]; then
            echo "STALL"
        else
            echo "REGRESSION"
        fi
    else
        if [ "$net_change" -le 0 ]; then
            echo "PURE_PROGRESS"
        elif [ "$net_change" -le 3 ]; then
            echo "DISCOVERY_PROGRESS"
        else
            echo "SLOW_PROGRESS"
        fi
    fi
}

# ============================================================================
# FUNCTION: analyze_progress
# Complete progress analysis between two iterations
# ============================================================================
analyze_progress() {
    local scope="$1"
    local current_iteration="$2"
    local previous_iteration="$3"

    echo "📊 Analyzing bug progress for $scope (iteration $previous_iteration → $current_iteration)..." >&2

    # Count bugs closed in current iteration
    local bugs_closed
    bugs_closed=$(count_closed_bugs "$scope" "$current_iteration")
    echo "  Bugs closed: $bugs_closed" >&2

    # Count open bugs in each iteration
    local prev_open curr_open
    prev_open=$(count_open_bugs "$scope" "$previous_iteration")
    curr_open=$(count_open_bugs "$scope" "$current_iteration")
    echo "  Open bugs: $prev_open → $curr_open" >&2

    # Calculate net change
    local net_change=$((curr_open - prev_open))
    echo "  Net change: $net_change" >&2

    # Detect reopened bugs
    local reopened_bugs reopened_count
    reopened_bugs=$(detect_reopened_bugs "$scope" "$current_iteration" "$previous_iteration")
    reopened_count=$(echo "$reopened_bugs" | jq 'length')
    if [ "$reopened_count" -gt 0 ]; then
        echo "  ${YELLOW}⚠️  Reopened bugs: $reopened_count${NC}" >&2
        echo "$reopened_bugs" | jq -r '.[]' | while read bug_id; do
            echo "    - $bug_id" >&2
        done
    fi

    # Calculate progress category
    local progress_category
    progress_category=$(calculate_progress_category "$bugs_closed" "$net_change")

    case "$progress_category" in
        PURE_PROGRESS)
            echo "  ${GREEN}✅ PURE_PROGRESS: Closing bugs faster than finding them${NC}" >&2
            ;;
        DISCOVERY_PROGRESS)
            echo "  ${GREEN}✅ DISCOVERY_PROGRESS: Fixing bugs while discovering hidden issues${NC}" >&2
            ;;
        SLOW_PROGRESS)
            echo "  ${YELLOW}⚠️  SLOW_PROGRESS: Fixing bugs but finding many new ones${NC}" >&2
            ;;
        STALL)
            echo "  ${RED}❌ STALL: No bugs fixed, no movement${NC}" >&2
            ;;
        REGRESSION)
            echo "  ${RED}❌ REGRESSION: No bugs fixed, more bugs found${NC}" >&2
            ;;
    esac

    # Output JSON result
    jq -n \
        --arg category "$progress_category" \
        --argjson closed "$bugs_closed" \
        --argjson prev_open "$prev_open" \
        --argjson curr_open "$curr_open" \
        --argjson net_change "$net_change" \
        --argjson reopened "$reopened_count" \
        --argjson reopened_bugs "$reopened_bugs" \
        '{
            progress_category: $category,
            bugs_closed: $closed,
            bugs_open_previous: $prev_open,
            bugs_open_current: $curr_open,
            net_change: $net_change,
            bugs_reopened: $reopened,
            reopened_bug_ids: $reopened_bugs
        }'
}

# ============================================================================
# FUNCTION: calculate_progress_score
# Calculate numeric progress score (for trending)
# ============================================================================
calculate_progress_score() {
    local scope="$1"
    local current_iteration="$2"
    local previous_iteration="$3"

    local analysis
    analysis=$(analyze_progress "$scope" "$current_iteration" "$previous_iteration")

    local bugs_closed net_change bugs_reopened
    bugs_closed=$(echo "$analysis" | jq -r '.bugs_closed')
    net_change=$(echo "$analysis" | jq -r '.net_change')
    bugs_reopened=$(echo "$analysis" | jq -r '.bugs_reopened')

    # Progress score calculation:
    # +10 points per bug closed
    # -5 points per net new bug
    # -20 points per reopened bug
    local score=$((bugs_closed * 10 - net_change * 5 - bugs_reopened * 20))

    echo "$score"
}

# ============================================================================
# FUNCTION: should_continue_or_escalate
# Determine if iteration should continue or escalate to ERROR_RECOVERY
# Implements R615 two-tiered limits
# ============================================================================
should_continue_or_escalate() {
    local scope="$1"
    local current_iteration="$2"

    echo "🔍 R615: Checking iteration continuation criteria..." >&2

    # Get previous iteration
    local previous_iteration=$((current_iteration - 1))

    if [ "$previous_iteration" -lt 1 ]; then
        echo "  First iteration - automatically continue" >&2
        echo '{"decision": "CONTINUE", "reason": "First iteration"}'
        return 0
    fi

    # Analyze progress
    local analysis
    analysis=$(analyze_progress "$scope" "$current_iteration" "$previous_iteration")

    local progress_category bugs_reopened
    progress_category=$(echo "$analysis" | jq -r '.progress_category')
    bugs_reopened=$(echo "$analysis" | jq -r '.bugs_reopened')

    # Get current stall counter from state file
    local stall_count
    stall_count=$(jq -r ".project_progression.current_${scope,,}.bugs.progress_stalls // 0" "$STATE_FILE")

    # Update stall counter based on progress
    local new_stall_count
    case "$progress_category" in
        PURE_PROGRESS|DISCOVERY_PROGRESS)
            new_stall_count=0
            echo "  ✅ Progress detected - resetting stall counter" >&2
            ;;
        SLOW_PROGRESS)
            new_stall_count="$stall_count"
            echo "  ⚠️  Slow progress - maintaining stall counter at $stall_count" >&2
            ;;
        STALL|REGRESSION)
            new_stall_count=$((stall_count + 1))
            echo "  ❌ No progress - incrementing stall counter to $new_stall_count" >&2
            ;;
    esac

    # Check for flapping
    if [ "$bugs_reopened" -gt 0 ]; then
        echo "  🚨 WARNING: $bugs_reopened bugs reopened - treating as no-progress" >&2
        new_stall_count=$((new_stall_count + 1))
    fi

    # Check Tier 1: No-progress limit (5 stalls)
    if [ "$new_stall_count" -ge 5 ]; then
        echo "" >&2
        echo "${RED}❌❌❌ R615 NO-PROGRESS LIMIT EXCEEDED ❌❌❌${NC}" >&2
        echo "Stall counter: $new_stall_count / 5" >&2
        echo "No bugs fixed for 5 consecutive iterations" >&2
        echo "ESCALATING TO ERROR_RECOVERY" >&2

        jq -n \
            --arg decision "ERROR_RECOVERY" \
            --arg reason "R615: No-progress limit exceeded ($new_stall_count stalls)" \
            --argjson stalls "$new_stall_count" \
            --argjson analysis "$analysis" \
            '{
                decision: $decision,
                reason: $reason,
                stall_count: $stalls,
                limit_type: "NO_PROGRESS_LIMIT",
                progress_analysis: $analysis
            }'
        return 1
    fi

    # Check Tier 2: Some-progress limit (10 iterations)
    if [ "$current_iteration" -ge 10 ]; then
        echo "" >&2
        echo "${RED}❌❌❌ R615 SOME-PROGRESS LIMIT EXCEEDED ❌❌❌${NC}" >&2
        echo "Total iterations: $current_iteration / 10" >&2
        echo "Maximum iterations reached even with progress" >&2
        echo "ESCALATING TO ERROR_RECOVERY for replanning" >&2

        jq -n \
            --arg decision "ERROR_RECOVERY" \
            --arg reason "R615: Some-progress limit exceeded ($current_iteration iterations)" \
            --argjson iter "$current_iteration" \
            --argjson stalls "$new_stall_count" \
            --argjson analysis "$analysis" \
            '{
                decision: $decision,
                reason: $reason,
                total_iterations: $iter,
                stall_count: $stalls,
                limit_type: "SOME_PROGRESS_LIMIT",
                progress_analysis: $analysis
            }'
        return 1
    fi

    # Continue iteration
    echo "" >&2
    echo "${GREEN}✅ R615: Iteration continuation approved${NC}" >&2
    echo "  Stall counter: $new_stall_count / 5" >&2
    echo "  Total iterations: $current_iteration / 10" >&2

    jq -n \
        --arg decision "CONTINUE" \
        --arg reason "Within iteration limits, can continue" \
        --argjson iter "$current_iteration" \
        --argjson stalls "$new_stall_count" \
        --argjson analysis "$analysis" \
        '{
            decision: $decision,
            reason: $reason,
            total_iterations: $iter,
            stall_count: $stalls,
            progress_analysis: $analysis
        }'
    return 0
}

# ============================================================================
# MAIN
# ============================================================================
main() {
    local command="${1:-}"

    if [ -z "$command" ]; then
        echo "Usage: $0 <command> [args...]"
        echo ""
        echo "Commands:"
        echo "  analyze_progress SCOPE CURRENT_ITER PREVIOUS_ITER"
        echo "  count_closed_bugs SCOPE ITERATION"
        echo "  count_open_bugs SCOPE ITERATION"
        echo "  detect_reopened_bugs SCOPE CURRENT_ITER PREVIOUS_ITER"
        echo "  calculate_progress_score SCOPE CURRENT_ITER PREVIOUS_ITER"
        echo "  should_continue_or_escalate SCOPE CURRENT_ITER"
        echo ""
        echo "Scopes: WAVE, PHASE, PROJECT"
        exit 1
    fi

    case "$command" in
        analyze_progress)
            analyze_progress "${2:-}" "${3:-}" "${4:-}"
            ;;
        count_closed_bugs)
            count_closed_bugs "${2:-}" "${3:-}"
            ;;
        count_open_bugs)
            count_open_bugs "${2:-}" "${3:-}"
            ;;
        detect_reopened_bugs)
            detect_reopened_bugs "${2:-}" "${3:-}" "${4:-}"
            ;;
        calculate_progress_score)
            calculate_progress_score "${2:-}" "${3:-}" "${4:-}"
            ;;
        should_continue_or_escalate)
            should_continue_or_escalate "${2:-}" "${3:-}"
            ;;
        *)
            echo "ERROR: Unknown command: $command" >&2
            exit 1
            ;;
    esac
}

# Run main if executed directly
if [ "${BASH_SOURCE[0]}" = "$0" ]; then
    main "$@"
fi
