#!/bin/bash

# 🏭 SOFTWARE FACTORY 2.0 - COMPREHENSIVE EFFORT RESTORATION UTILITY
# ═══════════════════════════════════════════════════════════════════════════
# Purpose: Restore all effort directories from remote branches in the TARGET
#          repository (separate from the planning repository)
# 
# IMPORTANT: Software Factory 2.0 uses TWO repositories:
#   1. Planning Repository: Contains orchestrator state, rules, agents, templates
#      (e.g., software-factory-template or your planning repo)
#   2. Target Repository: Contains actual code implementations on effort branches
#      (e.g., idpbuilder-oci-build-push or your project repo)
#
# This script restores efforts from the TARGET repository branches.
#
# Usage: ./restore-all-efforts.sh [orchestrator-state.json] [target-repo-url]
# 
# Example: ./restore-all-efforts.sh orchestrator-state.json https://github.com/org/project.git
#
# This script:
# 1. Reads orchestrator-state.json to identify all efforts
# 2. Creates proper directory structure (efforts/phase#/wave#/effort-name)
# 3. Clones/checks out each effort from its remote branch in the TARGET repo
# 4. Handles split efforts with proper suffixes
# 5. Manages both local worktrees and remote branches
# 6. Provides comprehensive progress reporting and error handling
# ═══════════════════════════════════════════════════════════════════════════

set -euo pipefail

# ANSI Color codes for better output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
MAGENTA='\033[0;35m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Default values
STATE_FILE="${1:-orchestrator-state.json}"
TARGET_REPO="${2:-}"
EFFORTS_ROOT="efforts"
RESTORE_LOG="effort-restoration.log"
FAILED_EFFORTS=()
SUCCESSFUL_EFFORTS=()
SKIPPED_EFFORTS=()

# Timestamp function
timestamp() {
    date '+%Y-%m-%d %H:%M:%S %Z'
}

# Logging function
log() {
    local level="$1"
    shift
    local message="$*"
    echo -e "${level}${message}${NC}" | tee -a "$RESTORE_LOG"
}

# Print header
print_header() {
    echo "═══════════════════════════════════════════════════════════════════════════" | tee -a "$RESTORE_LOG"
    echo "🏭 SOFTWARE FACTORY 2.0 - EFFORT RESTORATION UTILITY" | tee -a "$RESTORE_LOG"
    echo "═══════════════════════════════════════════════════════════════════════════" | tee -a "$RESTORE_LOG"
    echo "Timestamp: $(timestamp)" | tee -a "$RESTORE_LOG"
    echo "═══════════════════════════════════════════════════════════════════════════" | tee -a "$RESTORE_LOG"
}

# Check prerequisites
check_prerequisites() {
    log "${CYAN}" "\n📋 Checking prerequisites..."
    
    # Check if jq is installed
    if ! command -v jq &> /dev/null; then
        log "${RED}" "❌ ERROR: jq is not installed. Please install jq first."
        log "${YELLOW}" "  Install with: apt-get install jq (Debian/Ubuntu) or brew install jq (macOS)"
        exit 1
    fi
    
    # Check if state file exists
    if [ ! -f "$STATE_FILE" ]; then
        log "${RED}" "❌ ERROR: State file '$STATE_FILE' not found!"
        log "${YELLOW}" "  Please provide path to orchestrator-state.json"
        exit 1
    fi
    
    # Check if git is available
    if ! command -v git &> /dev/null; then
        log "${RED}" "❌ ERROR: git is not installed!"
        exit 1
    fi
    
    # Try to detect target repository if not provided
    if [ -z "$TARGET_REPO" ]; then
        log "${YELLOW}" "⚠️  No target repository URL provided. Attempting to detect from state file..."
        
        # First, try to get project_repository field from state file (new format)
        local project_repo=$(jq -r '.project_info.project_repository // ""' "$STATE_FILE" 2>/dev/null)
        if [ -n "$project_repo" ] && [ "$project_repo" != "null" ]; then
            TARGET_REPO="$project_repo"
            log "${GREEN}" "✅ Detected target repository from state file: $TARGET_REPO"
        else
            # Try legacy repository field
            project_repo=$(jq -r '.project_info.repository // ""' "$STATE_FILE" 2>/dev/null)
            if [ -n "$project_repo" ] && [ "$project_repo" != "null" ]; then
                TARGET_REPO="$project_repo"
                log "${GREEN}" "✅ Detected target repository from state file (legacy field): $TARGET_REPO"
            else
                # Try to extract from branch names in state file
                local sample_branch=$(jq -r '.efforts_completed[0].branch // .efforts_in_progress[0].branch // ""' "$STATE_FILE" 2>/dev/null)
                if [ -n "$sample_branch" ] && [ "$sample_branch" != "null" ]; then
                    # Extract project name from branch (e.g., idpbuilder-oci-build-push/phase1/wave1/effort)
                    local project_name=$(echo "$sample_branch" | cut -d'/' -f1)
                    
                    # Try to construct from current git remote
                    if git remote -v 2>/dev/null | grep -q origin; then
                        local base_url=$(git remote get-url origin 2>/dev/null | sed 's|/[^/]*\.git$||')
                        if [ -n "$base_url" ] && [ -n "$project_name" ]; then
                            TARGET_REPO="${base_url}/${project_name}.git"
                            log "${GREEN}" "✅ Detected target repository: $TARGET_REPO"
                        fi
                    fi
                fi
            fi
        fi
        
        if [ -z "$TARGET_REPO" ]; then
            log "${RED}" "❌ ERROR: Could not detect target repository URL"
            log "${YELLOW}" "  Please provide as second argument: $0 $STATE_FILE <target-repo-url>"
            log "${YELLOW}" "  Note: This should be your PROJECT repository, not the planning repository"
            exit 1
        fi
    fi
    
    log "${GREEN}" "✅ All prerequisites met"
}

# Create directory structure
create_directory_structure() {
    local phase="$1"
    local wave="$2"
    local effort="$3"
    
    local dir_path="${EFFORTS_ROOT}/phase${phase}/wave${wave}/${effort}"
    
    if [ ! -d "$(dirname "$dir_path")" ]; then
        mkdir -p "$(dirname "$dir_path")"
        log "${BLUE}" "  📁 Created directory structure: $(dirname "$dir_path")"
    fi
    
    echo "$dir_path"
}

# Clone or update effort
clone_or_update_effort() {
    local branch="$1"
    local dir_path="$2"
    local effort_name="$3"
    local base_branch="${4:-main}"
    
    log "${CYAN}" "\n🔄 Processing effort: ${effort_name}"
    log "${BLUE}" "  Branch: ${branch}"
    log "${BLUE}" "  Path: ${dir_path}"
    log "${BLUE}" "  Base: ${base_branch}"
    
    # Check if directory already exists
    if [ -d "$dir_path" ]; then
        log "${YELLOW}" "  ⚠️  Directory already exists. Checking status..."
        
        cd "$dir_path"
        
        # Check if it's a git repository
        if [ -d .git ]; then
            local current_branch=$(git branch --show-current 2>/dev/null || echo "unknown")
            local has_changes=$(git status --porcelain 2>/dev/null | wc -l)
            
            if [ "$has_changes" -gt 0 ]; then
                log "${YELLOW}" "  ⚠️  Uncommitted changes detected! Skipping to preserve local work."
                SKIPPED_EFFORTS+=("${effort_name} (uncommitted changes)")
                cd - > /dev/null
                return 1
            fi
            
            if [ "$current_branch" == "$branch" ]; then
                log "${GREEN}" "  ✅ Already on correct branch. Pulling latest changes..."
                if git pull origin "$branch" 2>&1 | tee -a "$RESTORE_LOG"; then
                    log "${GREEN}" "  ✅ Successfully updated ${effort_name}"
                    SUCCESSFUL_EFFORTS+=("${effort_name} (updated)")
                else
                    log "${YELLOW}" "  ⚠️  Pull failed, but continuing..."
                fi
            else
                log "${YELLOW}" "  ⚠️  On different branch ($current_branch). Switching to $branch..."
                if git checkout "$branch" 2>&1 | tee -a "$RESTORE_LOG"; then
                    git pull origin "$branch" 2>&1 | tee -a "$RESTORE_LOG"
                    log "${GREEN}" "  ✅ Successfully switched and updated ${effort_name}"
                    SUCCESSFUL_EFFORTS+=("${effort_name} (switched)")
                else
                    log "${RED}" "  ❌ Failed to switch branch"
                    FAILED_EFFORTS+=("${effort_name} (branch switch failed)")
                fi
            fi
        else
            log "${YELLOW}" "  ⚠️  Directory exists but is not a git repository. Backing up and re-cloning..."
            cd ..
            mv "${effort_name}" "${effort_name}.backup.$(date +%s)"
            clone_fresh_effort "$branch" "$dir_path" "$effort_name"
        fi
        
        cd - > /dev/null 2>&1 || true
    else
        clone_fresh_effort "$branch" "$dir_path" "$effort_name"
    fi
}

# Clone fresh effort
clone_fresh_effort() {
    local branch="$1"
    local dir_path="$2"
    local effort_name="$3"
    
    log "${BLUE}" "  📥 Cloning fresh effort..."
    
    # First check if branch exists in remote
    local branch_exists=$(git ls-remote --heads "$TARGET_REPO" "$branch" 2>/dev/null | wc -l)
    
    if [ "$branch_exists" -eq 0 ]; then
        log "${YELLOW}" "  ⚠️  Branch '$branch' does not exist in remote repository"
        log "${YELLOW}" "     This effort may not have been pushed yet or was deleted"
        
        # Check if it's a local worktree that needs to be pushed
        local local_worktree="/home/vscode/workspaces/$(echo "$branch" | cut -d'/' -f1)"
        if [ -d "$local_worktree" ]; then
            log "${BLUE}" "     Found local worktree at: $local_worktree"
            log "${YELLOW}" "     You may need to push this branch from the local worktree"
        fi
        
        SKIPPED_EFFORTS+=("${effort_name} (branch not in remote)")
        return 1
    fi
    
    # Try to clone the specific branch
    if git clone -b "$branch" --single-branch "$TARGET_REPO" "$dir_path" 2>&1 | tee -a "$RESTORE_LOG"; then
        log "${GREEN}" "  ✅ Successfully cloned ${effort_name}"
        SUCCESSFUL_EFFORTS+=("${effort_name} (cloned)")
        return 0
    else
        log "${YELLOW}" "  ⚠️  Direct branch clone failed. Trying alternative approach..."
        
        # Try cloning main and then checking out the branch
        if git clone "$TARGET_REPO" "$dir_path" 2>&1 | tee -a "$RESTORE_LOG"; then
            cd "$dir_path"
            
            # Fetch all branches
            git fetch --all 2>&1 | tee -a "$RESTORE_LOG"
            
            # Try to checkout the branch
            if git checkout -b "$branch" "origin/$branch" 2>&1 | tee -a "$RESTORE_LOG"; then
                log "${GREEN}" "  ✅ Successfully cloned and checked out ${effort_name}"
                SUCCESSFUL_EFFORTS+=("${effort_name} (cloned-alt)")
            else
                log "${RED}" "  ❌ Branch '$branch' not found in remote"
                FAILED_EFFORTS+=("${effort_name} (branch not found)")
            fi
            
            cd - > /dev/null
        else
            log "${RED}" "  ❌ Failed to clone repository"
            FAILED_EFFORTS+=("${effort_name} (clone failed)")
            return 1
        fi
    fi
}

# Process completed efforts
process_completed_efforts() {
    log "${BOLD}${CYAN}" "\n📋 Processing COMPLETED efforts..."
    
    local completed_count=$(jq '.efforts_completed | length' "$STATE_FILE")
    
    if [ "$completed_count" -eq 0 ]; then
        log "${YELLOW}" "  No completed efforts found"
        return
    fi
    
    log "${BLUE}" "  Found ${completed_count} completed efforts"
    
    for i in $(seq 0 $((completed_count - 1))); do
        # Extract values directly from state file using proper jq paths
        local phase=$(jq -r ".efforts_completed[$i].phase" "$STATE_FILE")
        local wave=$(jq -r ".efforts_completed[$i].wave" "$STATE_FILE")
        local effort_id=$(jq -r ".efforts_completed[$i].effort_id" "$STATE_FILE")
        local branch=$(jq -r ".efforts_completed[$i].branch" "$STATE_FILE")
        local base_branch=$(jq -r ".efforts_completed[$i].base_branch // \"main\"" "$STATE_FILE")
        local status=$(jq -r ".efforts_completed[$i].status" "$STATE_FILE")
        
        # Extract effort name from effort_id (e.g., E2.1.1-image-builder -> image-builder)
        local name="${effort_id#*-}"  # Remove everything before and including the first dash
        if [ "$name" == "$effort_id" ] || [ -z "$name" ] || [ "$name" == "null" ]; then
            # If no dash found or empty, try extracting from branch name
            name=$(echo "$branch" | sed 's|.*/||')  # Get last part after /
        fi
        
        # Handle split efforts
        if [ "$status" == "SPLIT_DEPRECATED" ]; then
            log "${YELLOW}" "  ⚠️  Effort ${name} was split. Processing splits..."
            
            # Get replacement splits directly from state file
            local splits=$(jq -r ".efforts_completed[$i].replacement_splits[]?" "$STATE_FILE" 2>/dev/null)
            if [ -n "$splits" ]; then
                local split_num=1
                for split_branch in $splits; do
                    local split_name="${name}-split-$(printf "%03d" $split_num)"
                    local dir_path=$(create_directory_structure "$phase" "$wave" "$split_name")
                    clone_or_update_effort "$split_branch" "$dir_path" "$split_name" "$base_branch"
                    split_num=$((split_num + 1))
                done
            fi
        else
            local dir_path=$(create_directory_structure "$phase" "$wave" "$name")
            clone_or_update_effort "$branch" "$dir_path" "$name" "$base_branch"
        fi
    done
}

# Process in-progress efforts
process_in_progress_efforts() {
    log "${BOLD}${CYAN}" "\n📋 Processing IN-PROGRESS efforts..."
    
    local progress_count=$(jq '.efforts_in_progress | length' "$STATE_FILE")
    
    if [ "$progress_count" -eq 0 ]; then
        log "${YELLOW}" "  No in-progress efforts found"
        return
    fi
    
    log "${BLUE}" "  Found ${progress_count} in-progress efforts"
    
    for i in $(seq 0 $((progress_count - 1))); do
        # Extract values directly from state file using proper jq paths
        local phase=$(jq -r ".efforts_in_progress[$i].phase" "$STATE_FILE")
        local wave=$(jq -r ".efforts_in_progress[$i].wave" "$STATE_FILE")
        local effort_id=$(jq -r ".efforts_in_progress[$i].effort_id" "$STATE_FILE")
        local branch=$(jq -r ".efforts_in_progress[$i].branch" "$STATE_FILE")
        local base_branch=$(jq -r ".efforts_in_progress[$i].base_branch // \"main\"" "$STATE_FILE")
        
        # Extract effort name from effort_id (e.g., E2.1.1-image-builder -> image-builder)
        local name="${effort_id#*-}"  # Remove everything before and including the first dash
        if [ "$name" == "$effort_id" ] || [ -z "$name" ] || [ "$name" == "null" ]; then
            # If no dash found or empty, try extracting from branch name
            name=$(echo "$branch" | sed 's|.*/||')  # Get last part after /
        fi
        
        local dir_path=$(create_directory_structure "$phase" "$wave" "$name")
        clone_or_update_effort "$branch" "$dir_path" "$name" "$base_branch"
    done
}

# Process split tracking
process_split_tracking() {
    log "${BOLD}${CYAN}" "\n📋 Processing SPLIT TRACKING efforts..."
    
    local has_splits=$(jq '.split_tracking | length' "$STATE_FILE")
    
    if [ "$has_splits" -eq 0 ]; then
        log "${YELLOW}" "  No split tracking found"
        return
    fi
    
    # Get all effort names from split_tracking
    local effort_names=$(jq -r '.split_tracking | keys[]' "$STATE_FILE")
    
    for effort_name in $effort_names; do
        log "${BLUE}" "  Processing splits for: ${effort_name}"
        
        # Check if this effort has splits
        local split_count=$(jq ".split_tracking[\"${effort_name}\"].splits | length" "$STATE_FILE" 2>/dev/null)
        
        if [ -n "$split_count" ] && [ "$split_count" -gt 0 ]; then
            # Process each split directly from state file
            for j in $(seq 0 $((split_count - 1))); do
                local split_branch=$(jq -r ".split_tracking[\"${effort_name}\"].splits[$j].branch" "$STATE_FILE")
                local split_status=$(jq -r ".split_tracking[\"${effort_name}\"].splits[$j].status" "$STATE_FILE")
                local split_id=$(jq -r ".split_tracking[\"${effort_name}\"].splits[$j].split_id" "$STATE_FILE")
                local base_branch=$(jq -r ".split_tracking[\"${effort_name}\"].splits[$j].base_branch // \"main\"" "$STATE_FILE")
                
                # Extract phase and wave from branch name
                if [[ "$split_branch" =~ phase([0-9]+)/wave([0-9]+)/ ]]; then
                    local phase="${BASH_REMATCH[1]}"
                    local wave="${BASH_REMATCH[2]}"
                    local split_name="${effort_name}-split-${split_id}"
                    
                    local dir_path=$(create_directory_structure "$phase" "$wave" "$split_name")
                    
                    if [ "$split_status" != "DEPRECATED" ]; then
                        clone_or_update_effort "$split_branch" "$dir_path" "$split_name" "$base_branch"
                    else
                        log "${YELLOW}" "    ⏭️  Skipping deprecated split: ${split_name}"
                        SKIPPED_EFFORTS+=("${split_name} (deprecated)")
                    fi
                fi
            done
        fi
    done
}

# Process integration workspaces
process_integration_workspaces() {
    log "${BOLD}${CYAN}" "\n📋 Processing INTEGRATION workspaces..."
    
    local has_integrations=$(jq '.integration_branches | length' "$STATE_FILE")
    
    if [ "$has_integrations" -eq 0 ]; then
        log "${YELLOW}" "  No integration branches found"
        return
    fi
    
    for i in $(seq 0 $((has_integrations - 1))); do
        # Extract values directly from state file using proper jq paths
        local branch=$(jq -r ".integration_branches[$i].branch" "$STATE_FILE")
        local phase=$(jq -r ".integration_branches[$i].phase" "$STATE_FILE")
        local wave=$(jq -r ".integration_branches[$i].wave" "$STATE_FILE")
        
        local dir_path="${EFFORTS_ROOT}/phase${phase}/wave${wave}/integration-workspace"
        
        log "${CYAN}" "\n🔄 Processing integration: phase${phase}/wave${wave}"
        log "${BLUE}" "  Branch: ${branch}"
        log "${BLUE}" "  Path: ${dir_path}"
        
        mkdir -p "$(dirname "$dir_path")"
        clone_or_update_effort "$branch" "$dir_path" "integration-phase${phase}-wave${wave}" "main"
    done
}

# Verify restoration
verify_restoration() {
    log "${BOLD}${CYAN}" "\n🔍 Verifying restoration..."
    
    local total_dirs=$(find "$EFFORTS_ROOT" -type d -name ".git" 2>/dev/null | wc -l)
    log "${BLUE}" "  Total git repositories restored: ${total_dirs}"
    
    # Check directory structure
    if [ -d "$EFFORTS_ROOT" ]; then
        log "${GREEN}" "  ✅ Efforts root directory exists"
        
        # Count phases
        local phase_count=$(find "$EFFORTS_ROOT" -maxdepth 1 -type d -name "phase*" | wc -l)
        log "${BLUE}" "  📁 Phases found: ${phase_count}"
        
        # List structure
        log "${CYAN}" "\n📂 Restored structure:"
        tree -L 3 "$EFFORTS_ROOT" 2>/dev/null || ls -la "$EFFORTS_ROOT"
    else
        log "${RED}" "  ❌ No efforts directory created"
    fi
}

# Generate summary report
generate_summary() {
    log "${BOLD}${CYAN}" "\n═══════════════════════════════════════════════════════════════════════════"
    log "${BOLD}${CYAN}" "📊 RESTORATION SUMMARY"
    log "${BOLD}${CYAN}" "═══════════════════════════════════════════════════════════════════════════"
    
    log "${GREEN}" "\n✅ Successful: ${#SUCCESSFUL_EFFORTS[@]} efforts"
    for effort in "${SUCCESSFUL_EFFORTS[@]}"; do
        log "${GREEN}" "    • ${effort}"
    done
    
    if [ ${#SKIPPED_EFFORTS[@]} -gt 0 ]; then
        log "${YELLOW}" "\n⏭️  Skipped: ${#SKIPPED_EFFORTS[@]} efforts"
        for effort in "${SKIPPED_EFFORTS[@]}"; do
            log "${YELLOW}" "    • ${effort}"
        done
    fi
    
    if [ ${#FAILED_EFFORTS[@]} -gt 0 ]; then
        log "${RED}" "\n❌ Failed: ${#FAILED_EFFORTS[@]} efforts"
        for effort in "${FAILED_EFFORTS[@]}"; do
            log "${RED}" "    • ${effort}"
        done
    fi
    
    log "${CYAN}" "\n📝 Full log saved to: ${RESTORE_LOG}"
    log "${CYAN}" "⏰ Completed at: $(timestamp)"
    
    # Create a restoration state file
    local restore_state="effort-restoration-state.json"
    {
        echo "{"
        echo "  \"restoration_timestamp\": \"$(timestamp)\","
        echo "  \"state_file_used\": \"$STATE_FILE\","
        echo "  \"target_repository\": \"$TARGET_REPO\","
        echo "  \"successful_count\": ${#SUCCESSFUL_EFFORTS[@]},"
        echo "  \"failed_count\": ${#FAILED_EFFORTS[@]},"
        echo "  \"skipped_count\": ${#SKIPPED_EFFORTS[@]},"
        echo "  \"successful_efforts\": ["
        local first=true
        for effort in "${SUCCESSFUL_EFFORTS[@]}"; do
            if [ "$first" = true ]; then
                echo -n "    \"$effort\""
                first=false
            else
                echo ","
                echo -n "    \"$effort\""
            fi
        done
        [ ${#SUCCESSFUL_EFFORTS[@]} -gt 0 ] && echo ""
        echo "  ],"
        echo "  \"failed_efforts\": ["
        first=true
        for effort in "${FAILED_EFFORTS[@]}"; do
            if [ "$first" = true ]; then
                echo -n "    \"$effort\""
                first=false
            else
                echo ","
                echo -n "    \"$effort\""
            fi
        done
        [ ${#FAILED_EFFORTS[@]} -gt 0 ] && echo ""
        echo "  ],"
        echo "  \"skipped_efforts\": ["
        first=true
        for effort in "${SKIPPED_EFFORTS[@]}"; do
            if [ "$first" = true ]; then
                echo -n "    \"$effort\""
                first=false
            else
                echo ","
                echo -n "    \"$effort\""
            fi
        done
        [ ${#SKIPPED_EFFORTS[@]} -gt 0 ] && echo ""
        echo "  ]"
        echo "}"
    } > "$restore_state"
    
    log "${GREEN}" "📊 Restoration state saved to: ${restore_state}"
}

# Main execution
main() {
    # Initialize log
    echo "Starting effort restoration at $(timestamp)" > "$RESTORE_LOG"
    
    print_header
    check_prerequisites
    
    log "${BOLD}${CYAN}" "\n🏭 Starting effort restoration from: ${STATE_FILE}"
    log "${BLUE}" "📍 Target repository: ${TARGET_REPO}"
    
    # Create efforts root if it doesn't exist
    if [ ! -d "$EFFORTS_ROOT" ]; then
        mkdir -p "$EFFORTS_ROOT"
        log "${GREEN}" "✅ Created efforts root directory"
    fi
    
    # Process different types of efforts
    process_completed_efforts
    process_in_progress_efforts
    process_split_tracking
    process_integration_workspaces
    
    # Verify and report
    verify_restoration
    generate_summary
    
    # Exit with appropriate code
    if [ ${#FAILED_EFFORTS[@]} -gt 0 ]; then
        exit 1
    else
        exit 0
    fi
}

# Help function
show_help() {
    cat << EOF
═══════════════════════════════════════════════════════════════════════════
🏭 SOFTWARE FACTORY 2.0 - EFFORT RESTORATION UTILITY
═══════════════════════════════════════════════════════════════════════════

USAGE:
    $0 [orchestrator-state.json] [target-repo-url]

DESCRIPTION:
    Restores all effort directories from their remote branches in the TARGET
    repository. Software Factory 2.0 uses a dual-repository architecture:
    
    • Planning Repository: Contains orchestrator state, rules, agents
      (e.g., software-factory-template)
    • Target Repository: Contains actual code on effort branches
      (e.g., idpbuilder-oci-build-push)
    
    This script restores efforts from the TARGET repository.

ARGUMENTS:
    orchestrator-state.json  Path to the state file (default: orchestrator-state.json)
    target-repo-url         URL of the TARGET repository containing efforts
                           (will attempt auto-detection if not provided)
                           NOTE: This is your PROJECT repo, not the planning repo

EXAMPLES:
    # Basic usage with defaults
    $0

    # Specify state file
    $0 my-state.json

    # Full specification
    $0 orchestrator-state.json https://github.com/org/project.git

FEATURES:
    • Reads orchestrator-state.json to identify all efforts
    • Creates proper directory structure (efforts/phase#/wave#/effort-name)
    • Clones/checks out each effort from its remote branch
    • Handles split efforts with proper suffixes (-split-001, -split-002)
    • Processes integration workspaces
    • Preserves local changes (skips dirs with uncommitted changes)
    • Comprehensive progress reporting
    • Error handling and recovery
    • Generates restoration state file for tracking

OUTPUT FILES:
    • effort-restoration.log - Full restoration log
    • effort-restoration-state.json - Summary of restoration results

EXIT CODES:
    0 - All efforts restored successfully (or skipped safely)
    1 - One or more efforts failed to restore

═══════════════════════════════════════════════════════════════════════════
EOF
}

# Parse command line arguments
if [ "$#" -gt 0 ] && [[ "$1" == "-h" || "$1" == "--help" || "$1" == "help" ]]; then
    show_help
    exit 0
fi

# Run main function
main "$@"