#!/bin/bash

# 🏭 SOFTWARE FACTORY 2.0 - EFFORT RESTORATION VERIFICATION
# ═══════════════════════════════════════════════════════════════════════════
# Purpose: Verify that restored efforts match the orchestrator state file
# Usage: ./verify-restored-efforts.sh [orchestrator-state-v3.json]
# 
# This script validates:
# 1. All efforts from state file exist on disk
# 2. Efforts are on correct branches
# 3. Directory structure matches expectations
# 4. No orphaned efforts exist
# ═══════════════════════════════════════════════════════════════════════════

set -euo pipefail

# ANSI Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m'

# Configuration
STATE_FILE="${1:-orchestrator-state-v3.json}"
EFFORTS_ROOT="efforts"
VERIFICATION_LOG="effort-verification.log"
ISSUES_FOUND=()
WARNINGS_FOUND=()

# Counters
EXPECTED_EFFORTS=0
FOUND_EFFORTS=0
CORRECT_BRANCH=0
WRONG_BRANCH=0
MISSING_EFFORTS=0
ORPHANED_EFFORTS=0

# Logging function
log() {
    local level="$1"
    shift
    local message="$*"
    echo -e "${level}${message}${NC}" | tee -a "$VERIFICATION_LOG"
}

# Print header
print_header() {
    echo "═══════════════════════════════════════════════════════════════════════════" | tee "$VERIFICATION_LOG"
    echo "🏭 SOFTWARE FACTORY 2.0 - EFFORT VERIFICATION" | tee -a "$VERIFICATION_LOG"
    echo "═══════════════════════════════════════════════════════════════════════════" | tee -a "$VERIFICATION_LOG"
    echo "Timestamp: $(date '+%Y-%m-%d %H:%M:%S %Z')" | tee -a "$VERIFICATION_LOG"
    echo "State File: ${STATE_FILE}" | tee -a "$VERIFICATION_LOG"
    echo "═══════════════════════════════════════════════════════════════════════════" | tee -a "$VERIFICATION_LOG"
}

# Check prerequisites
check_prerequisites() {
    log "${CYAN}" "\n📋 Checking prerequisites..."
    
    if ! command -v jq &> /dev/null; then
        log "${RED}" "❌ ERROR: jq is not installed"
        exit 1
    fi
    
    if [ ! -f "$STATE_FILE" ]; then
        log "${RED}" "❌ ERROR: State file '$STATE_FILE' not found"
        exit 1
    fi
    
    if [ ! -d "$EFFORTS_ROOT" ]; then
        log "${RED}" "❌ ERROR: Efforts directory '$EFFORTS_ROOT' not found"
        log "${YELLOW}" "  Run restore-all-efforts.sh first"
        exit 1
    fi
    
    log "${GREEN}" "✅ Prerequisites satisfied"
}

# Verify single effort
verify_effort() {
    local phase="$1"
    local wave="$2"
    local name="$3"
    local expected_branch="$4"
    local effort_type="${5:-standard}"
    
    EXPECTED_EFFORTS=$((EXPECTED_EFFORTS + 1))
    
    local dir_path="${EFFORTS_ROOT}/phase${phase}/wave${wave}/${name}"
    
    log "${CYAN}" "\n🔍 Verifying: ${name}"
    log "${BLUE}" "  Expected path: ${dir_path}"
    log "${BLUE}" "  Expected branch: ${expected_branch}"
    log "${BLUE}" "  Type: ${effort_type}"
    
    if [ ! -d "$dir_path" ]; then
        log "${RED}" "  ❌ Directory not found!"
        ISSUES_FOUND+=("${name}: Directory missing at ${dir_path}")
        MISSING_EFFORTS=$((MISSING_EFFORTS + 1))
        return 1
    fi
    
    FOUND_EFFORTS=$((FOUND_EFFORTS + 1))
    
    if [ ! -d "${dir_path}/.git" ]; then
        log "${RED}" "  ❌ Not a git repository!"
        ISSUES_FOUND+=("${name}: Not a git repository")
        return 1
    fi
    
    # Check branch
    cd "$dir_path"
    local current_branch=$(git branch --show-current 2>/dev/null || echo "unknown")
    
    if [ "$current_branch" == "$expected_branch" ]; then
        log "${GREEN}" "  ✅ Correct branch: ${current_branch}"
        CORRECT_BRANCH=$((CORRECT_BRANCH + 1))
    else
        log "${RED}" "  ❌ Wrong branch! Current: ${current_branch}, Expected: ${expected_branch}"
        ISSUES_FOUND+=("${name}: Wrong branch (${current_branch} != ${expected_branch})")
        WRONG_BRANCH=$((WRONG_BRANCH + 1))
    fi
    
    # Check for uncommitted changes
    local changes=$(git status --porcelain 2>/dev/null | wc -l)
    if [ "$changes" -gt 0 ]; then
        log "${YELLOW}" "  ⚠️  Has ${changes} uncommitted changes"
        WARNINGS_FOUND+=("${name}: ${changes} uncommitted changes")
    fi
    
    # Check remote tracking
    local remote=$(git remote -v 2>/dev/null | head -1)
    if [ -z "$remote" ]; then
        log "${YELLOW}" "  ⚠️  No remote configured"
        WARNINGS_FOUND+=("${name}: No remote repository")
    else
        log "${BLUE}" "  📡 Remote: $(echo "$remote" | awk '{print $2}')"
    fi
    
    cd - > /dev/null
    
    return 0
}

# Verify completed efforts
verify_completed_efforts() {
    log "${BOLD}${CYAN}" "\n📋 Verifying COMPLETED efforts..."
    
    local count=$(jq '.efforts_completed | length' "$STATE_FILE")
    
    if [ "$count" -eq 0 ]; then
        log "${YELLOW}" "  No completed efforts in state file"
        return
    fi
    
    for i in $(seq 0 $((count - 1))); do
        local effort=$(jq -r ".efforts_completed[$i]" "$STATE_FILE")
        local phase=$(echo "$effort" | jq -r '.phase')
        local wave=$(echo "$effort" | jq -r '.wave')
        local name=$(echo "$effort" | jq -r '.name')
        local branch=$(echo "$effort" | jq -r '.branch')
        local status=$(echo "$effort" | jq -r '.status')
        
        if [ "$status" == "SPLIT_DEPRECATED" ]; then
            log "${YELLOW}" "  ⏭️  Skipping deprecated effort: ${name}"
            
            # Verify splits instead
            local splits=$(echo "$effort" | jq -r '.replacement_splits[]?' 2>/dev/null)
            if [ -n "$splits" ]; then
                local split_num=1
                for split_branch in $splits; do
                    local split_name="${name}-split-$(printf "%03d" $split_num)"
                    verify_effort "$phase" "$wave" "$split_name" "$split_branch" "split"
                    split_num=$((split_num + 1))
                done
            fi
        else
            verify_effort "$phase" "$wave" "$name" "$branch" "completed"
        fi
    done
}

# Verify in-progress efforts
verify_in_progress_efforts() {
    log "${BOLD}${CYAN}" "\n📋 Verifying IN-PROGRESS efforts..."
    
    local count=$(jq '.efforts_in_progress | length' "$STATE_FILE")
    
    if [ "$count" -eq 0 ]; then
        log "${YELLOW}" "  No in-progress efforts in state file"
        return
    fi
    
    for i in $(seq 0 $((count - 1))); do
        local effort=$(jq -r ".efforts_in_progress[$i]" "$STATE_FILE")
        local phase=$(echo "$effort" | jq -r '.phase')
        local wave=$(echo "$effort" | jq -r '.wave')
        local name=$(echo "$effort" | jq -r '.name')
        local branch=$(echo "$effort" | jq -r '.branch')
        
        verify_effort "$phase" "$wave" "$name" "$branch" "in-progress"
    done
}

# Check for orphaned efforts
check_orphaned_efforts() {
    log "${BOLD}${CYAN}" "\n🔍 Checking for orphaned efforts..."
    
    # Build list of expected effort paths from state file
    local expected_paths_file=$(mktemp)
    
    # Add completed efforts
    jq -r '.efforts_completed[] | 
        "\(.phase)/\(.wave)/\(.name)"' "$STATE_FILE" >> "$expected_paths_file" 2>/dev/null || true
    
    # Add completed effort splits
    jq -r '.efforts_completed[] | 
        select(.status == "SPLIT_DEPRECATED") | 
        .replacement_splits[]? as $split | 
        "\(.phase)/\(.wave)/\(.name)-split-" + ($split | split("/")[-1] | split("-")[-1])' "$STATE_FILE" >> "$expected_paths_file" 2>/dev/null || true
    
    # Add in-progress efforts
    jq -r '.efforts_in_progress[] | 
        "\(.phase)/\(.wave)/\(.name)"' "$STATE_FILE" >> "$expected_paths_file" 2>/dev/null || true
    
    # Add integration workspaces
    jq -r '.integration_branches[] | 
        "\(.phase)/\(.wave)/integration-workspace"' "$STATE_FILE" >> "$expected_paths_file" 2>/dev/null || true
    
    # Find all actual effort directories
    local actual_paths_file=$(mktemp)
    find "$EFFORTS_ROOT" -type d -name ".git" | while read -r git_dir; do
        local effort_dir=$(dirname "$git_dir")
        local rel_path=${effort_dir#${EFFORTS_ROOT}/}
        # Convert to phase/wave/name format
        echo "$rel_path" | sed 's/^phase//; s/\/wave/\//' >> "$actual_paths_file"
    done
    
    # Find orphaned efforts (in filesystem but not in state)
    while read -r actual_path; do
        if ! grep -q "^${actual_path}$" "$expected_paths_file"; then
            log "${YELLOW}" "  ⚠️  Orphaned effort found: ${EFFORTS_ROOT}/${actual_path}"
            WARNINGS_FOUND+=("Orphaned: ${actual_path}")
            ORPHANED_EFFORTS=$((ORPHANED_EFFORTS + 1))
        fi
    done < "$actual_paths_file"
    
    if [ "$ORPHANED_EFFORTS" -eq 0 ]; then
        log "${GREEN}" "  ✅ No orphaned efforts found"
    else
        log "${YELLOW}" "  ⚠️  Found ${ORPHANED_EFFORTS} orphaned efforts"
    fi
    
    # Cleanup temp files
    rm -f "$expected_paths_file" "$actual_paths_file"
}

# Verify directory structure
verify_directory_structure() {
    log "${BOLD}${CYAN}" "\n📂 Verifying directory structure..."
    
    # Check phases from state
    local phases=$(jq -r '[.efforts_completed[].phase, .efforts_in_progress[].phase] | unique | .[]' "$STATE_FILE" 2>/dev/null | sort -u)
    
    for phase in $phases; do
        local phase_dir="${EFFORTS_ROOT}/phase${phase}"
        if [ -d "$phase_dir" ]; then
            log "${GREEN}" "  ✅ Phase ${phase} directory exists"
            
            # Check waves
            local waves=$(jq -r "[.efforts_completed[] | select(.phase == ${phase}) | .wave, 
                                  .efforts_in_progress[] | select(.phase == ${phase}) | .wave] | 
                                  unique | .[]" "$STATE_FILE" 2>/dev/null | sort -u)
            
            for wave in $waves; do
                local wave_dir="${phase_dir}/wave${wave}"
                if [ -d "$wave_dir" ]; then
                    log "${GREEN}" "    ✅ Wave ${wave} directory exists"
                else
                    log "${RED}" "    ❌ Wave ${wave} directory missing"
                    ISSUES_FOUND+=("Missing directory: ${wave_dir}")
                fi
            done
        else
            log "${RED}" "  ❌ Phase ${phase} directory missing"
            ISSUES_FOUND+=("Missing directory: ${phase_dir}")
        fi
    done
}

# Generate verification report
generate_report() {
    log "${BOLD}${CYAN}" "\n═══════════════════════════════════════════════════════════════════════════"
    log "${BOLD}${CYAN}" "📊 VERIFICATION REPORT"
    log "${BOLD}${CYAN}" "═══════════════════════════════════════════════════════════════════════════"
    
    log "${CYAN}" "\n📈 Statistics:"
    log "${BLUE}" "  Expected efforts: ${EXPECTED_EFFORTS}"
    log "${BLUE}" "  Found efforts: ${FOUND_EFFORTS}"
    log "${GREEN}" "  ✅ Correct branch: ${CORRECT_BRANCH}"
    log "${RED}" "  ❌ Wrong branch: ${WRONG_BRANCH}"
    log "${RED}" "  ❌ Missing efforts: ${MISSING_EFFORTS}"
    log "${YELLOW}" "  ⚠️  Orphaned efforts: ${ORPHANED_EFFORTS}"
    
    if [ ${#ISSUES_FOUND[@]} -gt 0 ]; then
        log "${RED}" "\n❌ Critical Issues (${#ISSUES_FOUND[@]}):"
        for issue in "${ISSUES_FOUND[@]}"; do
            log "${RED}" "  • ${issue}"
        done
    else
        log "${GREEN}" "\n✅ No critical issues found!"
    fi
    
    if [ ${#WARNINGS_FOUND[@]} -gt 0 ]; then
        log "${YELLOW}" "\n⚠️  Warnings (${#WARNINGS_FOUND[@]}):"
        for warning in "${WARNINGS_FOUND[@]}"; do
            log "${YELLOW}" "  • ${warning}"
        done
    fi
    
    # Overall status
    log "${CYAN}" "\n📋 Overall Status:"
    if [ ${#ISSUES_FOUND[@]} -eq 0 ] && [ "$MISSING_EFFORTS" -eq 0 ] && [ "$WRONG_BRANCH" -eq 0 ]; then
        log "${GREEN}" "  ✅ VERIFICATION PASSED - All efforts correctly restored!"
        local exit_code=0
    else
        log "${RED}" "  ❌ VERIFICATION FAILED - Issues need attention"
        local exit_code=1
    fi
    
    log "${CYAN}" "\n📝 Full verification log: ${VERIFICATION_LOG}"
    log "${CYAN}" "⏰ Completed at: $(date '+%Y-%m-%d %H:%M:%S %Z')"
    
    # Create JSON report
    local json_report="effort-verification-report.json"
    {
        echo "{"
        echo "  \"verification_timestamp\": \"$(date -Iseconds)\","
        echo "  \"state_file\": \"${STATE_FILE}\","
        echo "  \"statistics\": {"
        echo "    \"expected_efforts\": ${EXPECTED_EFFORTS},"
        echo "    \"found_efforts\": ${FOUND_EFFORTS},"
        echo "    \"correct_branch\": ${CORRECT_BRANCH},"
        echo "    \"wrong_branch\": ${WRONG_BRANCH},"
        echo "    \"missing_efforts\": ${MISSING_EFFORTS},"
        echo "    \"orphaned_efforts\": ${ORPHANED_EFFORTS}"
        echo "  },"
        echo "  \"issues_count\": ${#ISSUES_FOUND[@]},"
        echo "  \"warnings_count\": ${#WARNINGS_FOUND[@]},"
        echo "  \"verification_passed\": $([ $exit_code -eq 0 ] && echo "true" || echo "false")"
        echo "}"
    } > "$json_report"
    
    log "${GREEN}" "📊 JSON report saved to: ${json_report}"
    
    return $exit_code
}

# Main execution
main() {
    print_header
    check_prerequisites
    
    log "${BOLD}${CYAN}" "\n🔍 Starting verification of restored efforts..."
    
    verify_completed_efforts
    verify_in_progress_efforts
    verify_directory_structure
    check_orphaned_efforts
    
    generate_report
    exit $?
}

# Help function
show_help() {
    cat << EOF
═══════════════════════════════════════════════════════════════════════════
🏭 SOFTWARE FACTORY 2.0 - EFFORT VERIFICATION UTILITY
═══════════════════════════════════════════════════════════════════════════

USAGE:
    $0 [orchestrator-state-v3.json]

DESCRIPTION:
    Verifies that restored efforts match the orchestrator state file.
    Checks directory structure, branch correctness, and identifies issues.

ARGUMENTS:
    orchestrator-state-v3.json  Path to state file (default: orchestrator-state-v3.json)

VERIFICATIONS:
    • All efforts from state file exist on disk
    • Each effort is on the correct git branch
    • Directory structure matches phase/wave organization
    • No orphaned efforts (on disk but not in state)
    • Git repository health checks
    • Remote tracking configuration

EXAMPLES:
    # Basic verification
    $0

    # Verify with specific state file
    $0 my-state.json

OUTPUT FILES:
    • effort-verification.log - Detailed verification log
    • effort-verification-report.json - JSON verification summary

EXIT CODES:
    0 - All verifications passed
    1 - One or more verification failures

═══════════════════════════════════════════════════════════════════════════
EOF
}

# Parse arguments
if [ "$#" -gt 0 ] && [[ "$1" == "-h" || "$1" == "--help" || "$1" == "help" ]]; then
    show_help
    exit 0
fi

# Run main
main "$@"