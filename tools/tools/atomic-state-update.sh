#!/usr/bin/env bash
# Atomic State Update Tool for SF 3.0
# Purpose: Update all 4 state files atomically (backup → update → validate → commit → push)
# Rule: R288 Multi-File Atomic Update Protocol

set -euo pipefail

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# State files (all 4 must be updated atomically)
STATE_FILES=(
    "orchestrator-state-v3.json"
    "bug-tracking.json"
    "integration-containers.json"
    "fix-cascade-state.json"
)

BACKUP_DIR=".state-backup"
BACKUP_TIMESTAMP=$(date -u '+%Y%m%d-%H%M%S')

# Usage message
usage() {
    cat << EOF
Usage: $0 [OPTIONS]

Atomically update all 4 SF 3.0 state files with validation and git commit.

Options:
    --test              Run in test mode (validation only, no commit)
    --help              Show this help message

The tool performs 6 steps in sequence:
    1. backup_state_files       - Create timestamped backups
    1.5 validate_state_transition - Validate state machine transition (P0 fix #2)
    2. update_state_files       - Apply updates to state files
    3. validate_state_files     - Validate all files against schemas
    4. commit_state_files       - Atomic git commit (all 4 files)
    5. push_state_files         - Push to remote

On validation failure, automatic rollback restores from backup.

Example:
    $0                  # Normal execution
    $0 --test           # Test mode (validation only)

EOF
    exit 0
}

# Step 1: Backup all state files
backup_state_files() {
    echo -e "${BLUE}[STEP 1/5] Backing up state files...${NC}"

    # Create backup directory if it doesn't exist
    mkdir -p "$BACKUP_DIR"

    local backup_path="$BACKUP_DIR/$BACKUP_TIMESTAMP"
    mkdir -p "$backup_path"

    local backed_up=0
    for state_file in "${STATE_FILES[@]}"; do
        if [ -f "$state_file" ]; then
            cp "$state_file" "$backup_path/$state_file"
            echo -e "  ✓ Backed up: $state_file → $backup_path/$state_file"
            ((backed_up++))
        else
            echo -e "  ${YELLOW}⚠ File not found (skipping): $state_file${NC}"
        fi
    done

    if [ $backed_up -eq 0 ]; then
        echo -e "${RED}✗ ERROR: No state files found to backup${NC}"
        return 1
    fi

    echo -e "${GREEN}✓ Backup complete: $backed_up files saved to $backup_path${NC}"
    return 0
}

# Step 1.5: Validate state transition (Test 1.5 P0 fix #2)
validate_state_transition() {
    echo -e "${BLUE}[STEP 1.5/6] Validating state transition...${NC}"

    local state_file="orchestrator-state-v3.json"
    local validation_script="utilities/validate-state-transition.sh"

    # Skip if orchestrator state file doesn't exist
    if [ ! -f "$state_file" ]; then
        echo -e "  ${YELLOW}ℹ Skipping: $state_file not found (may be initial setup)${NC}"
        echo -e "${GREEN}✓ State transition validation skipped${NC}"
        return 0
    fi

    # Skip if validation script doesn't exist
    if [ ! -f "$validation_script" ]; then
        echo -e "  ${YELLOW}⚠ WARNING: Validation script not found: $validation_script${NC}"
        echo -e "  ${YELLOW}ℹ Skipping state transition validation${NC}"
        return 0
    fi

    # Extract previous and current state
    local previous_state=$(jq -r '.state_machine.previous_state // "INIT"' "$state_file")
    local current_state=$(jq -r '.state_machine.current_state // "UNKNOWN"' "$state_file")

    if [ "$current_state" = "UNKNOWN" ] || [ "$current_state" = "null" ]; then
        echo -e "  ${YELLOW}⚠ WARNING: Cannot determine current state from $state_file${NC}"
        echo -e "  ${YELLOW}ℹ Skipping state transition validation${NC}"
        return 0
    fi

    echo -e "  Transition: ${BLUE}$previous_state${NC} → ${BLUE}$current_state${NC}"

    # Run validation
    if bash "$validation_script" "$previous_state" "$current_state" > /tmp/state-validation.log 2>&1; then
        echo -e "  ${GREEN}✓ State transition valid${NC}"
        echo -e "${GREEN}✓ State transition validation passed${NC}"
        return 0
    else
        echo -e "  ${RED}✗ INVALID STATE TRANSITION${NC}"
        echo -e ""
        cat /tmp/state-validation.log | sed 's/^/    /'
        echo -e ""
        echo -e "${RED}✗ State transition validation failed${NC}"
        echo -e "${RED}   This prevents invalid states like CREATE_NEXT_INFRASTRUCTURE from entering git history${NC}"
        return 1
    fi
}

# Step 2: Update state files (placeholder - actual updates applied externally)
update_state_files() {
    echo -e "${BLUE}[STEP 2/6] Updating state files...${NC}"

    # This function is a hook point for external updates
    # In practice, the caller modifies the state files before calling this script
    # OR passes update operations as parameters

    echo -e "  ${YELLOW}ℹ State files should be updated by caller before running this script${NC}"
    echo -e "  ${YELLOW}ℹ This step validates that files have been modified${NC}"

    # Check if any files were modified since backup
    local modified=0
    for state_file in "${STATE_FILES[@]}"; do
        if [ -f "$state_file" ]; then
            if git diff --quiet "$state_file" 2>/dev/null; then
                echo -e "  - No changes: $state_file"
            else
                echo -e "  ✓ Modified: $state_file"
                ((modified++))
            fi
        fi
    done

    if [ $modified -eq 0 ]; then
        echo -e "${YELLOW}⚠ WARNING: No state files were modified${NC}"
    fi

    echo -e "${GREEN}✓ Update step complete (caller responsibility)${NC}"
    return 0
}

# Step 3: Validate all state files against schemas
validate_state_files() {
    echo -e "${BLUE}[STEP 3/6] Validating state files against schemas...${NC}"

    local validation_script="tools/validate-state-file.sh"

    if [ ! -x "$validation_script" ]; then
        echo -e "${RED}✗ ERROR: Validation script not found or not executable: $validation_script${NC}"
        return 1
    fi

    local validated=0
    local failed=0

    for state_file in "${STATE_FILES[@]}"; do
        if [ -f "$state_file" ]; then
            echo -e "  Validating: $state_file"
            if bash "$validation_script" "$state_file" > /dev/null 2>&1; then
                echo -e "  ${GREEN}✓ Valid: $state_file${NC}"
                ((validated++))
            else
                echo -e "  ${RED}✗ INVALID: $state_file${NC}"
                bash "$validation_script" "$state_file" 2>&1 | sed 's/^/    /'
                ((failed++))
            fi
        fi
    done

    if [ $failed -gt 0 ]; then
        echo -e "${RED}✗ VALIDATION FAILED: $failed file(s) invalid${NC}"
        return 1
    fi

    if [ $validated -eq 0 ]; then
        echo -e "${YELLOW}⚠ WARNING: No state files found to validate${NC}"
        return 1
    fi

    echo -e "${GREEN}✓ Validation complete: $validated file(s) valid${NC}"
    return 0
}

# Step 4: Commit all state files atomically
commit_state_files() {
    echo -e "${BLUE}[STEP 4/6] Committing state files atomically...${NC}"

    # Add all state files to staging
    local staged=0
    for state_file in "${STATE_FILES[@]}"; do
        if [ -f "$state_file" ]; then
            git add "$state_file"
            echo -e "  ✓ Staged: $state_file"
            ((staged++))
        fi
    done

    if [ $staged -eq 0 ]; then
        echo -e "${YELLOW}⚠ No state files to commit${NC}"
        return 0
    fi

    # Create atomic commit message
    local commit_msg="state: Atomic update of $staged state file(s) [R288]

Updated state files:
$(for f in "${STATE_FILES[@]}"; do [ -f "$f" ] && echo "- $f"; done)

Timestamp: $BACKUP_TIMESTAMP
Backup: $BACKUP_DIR/$BACKUP_TIMESTAMP

🤖 Generated with Claude Code - Software Factory 3.0
"

    # Commit with R288 tag
    if git commit -m "$commit_msg"; then
        local commit_hash=$(git rev-parse --short HEAD)
        echo -e "${GREEN}✓ Commit successful: $commit_hash${NC}"
        return 0
    else
        echo -e "${RED}✗ Commit failed${NC}"
        return 1
    fi
}

# Step 5: Push to remote
push_state_files() {
    echo -e "${BLUE}[STEP 5/6] Pushing to remote...${NC}"

    local current_branch=$(git branch --show-current)

    if [ -z "$current_branch" ]; then
        echo -e "${RED}✗ ERROR: Not on a branch${NC}"
        return 1
    fi

    echo -e "  Branch: $current_branch"

    if git push; then
        echo -e "${GREEN}✓ Push successful to remote${NC}"
        return 0
    else
        echo -e "${RED}✗ Push failed${NC}"
        echo -e "${YELLOW}  Commit is local only. Run 'git push' manually.${NC}"
        return 1
    fi
}

# Rollback function - restore from backup
rollback_state_files() {
    echo -e "${YELLOW}[ROLLBACK] Restoring state files from backup...${NC}"

    local backup_path="$BACKUP_DIR/$BACKUP_TIMESTAMP"

    if [ ! -d "$backup_path" ]; then
        echo -e "${RED}✗ ERROR: Backup not found: $backup_path${NC}"
        return 1
    fi

    local restored=0
    for state_file in "${STATE_FILES[@]}"; do
        if [ -f "$backup_path/$state_file" ]; then
            cp "$backup_path/$state_file" "$state_file"
            echo -e "  ✓ Restored: $state_file"
            ((restored++))
        fi
    done

    echo -e "${GREEN}✓ Rollback complete: $restored file(s) restored${NC}"
    return 0
}

# Main execution
main() {
    local test_mode=false

    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --test)
                test_mode=true
                shift
                ;;
            --help)
                usage
                ;;
            *)
                echo -e "${RED}Unknown option: $1${NC}"
                usage
                ;;
        esac
    done

    echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"
    echo -e "${BLUE}  Atomic State Update Tool - SF 3.0${NC}"
    echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"

    if [ "$test_mode" = true ]; then
        echo -e "${YELLOW}[TEST MODE] Validation only, no commit/push${NC}"
    fi

    echo ""

    # Execute 6-step sequence
    if ! backup_state_files; then
        echo -e "${RED}✗ FAILED at step 1 (backup)${NC}"
        exit 1
    fi

    echo ""

    if ! validate_state_transition; then
        echo -e "${RED}✗ FAILED at step 1.5 (state transition validation)${NC}"
        rollback_state_files
        exit 2
    fi

    echo ""

    if ! update_state_files; then
        echo -e "${RED}✗ FAILED at step 2 (update)${NC}"
        rollback_state_files
        exit 3
    fi

    echo ""

    if ! validate_state_files; then
        echo -e "${RED}✗ FAILED at step 3 (validation)${NC}"
        rollback_state_files
        exit 4
    fi

    echo ""

    if [ "$test_mode" = true ]; then
        echo -e "${GREEN}✓ TEST MODE: Validation successful, skipping commit/push${NC}"
        exit 0
    fi

    if ! commit_state_files; then
        echo -e "${RED}✗ FAILED at step 4 (commit)${NC}"
        rollback_state_files
        exit 5
    fi

    echo ""

    if ! push_state_files; then
        echo -e "${YELLOW}⚠ PARTIAL PROJECT_DONE: Committed locally but push failed${NC}"
        exit 6
    fi

    echo ""
    echo -e "${GREEN}═══════════════════════════════════════════════════════${NC}"
    echo -e "${GREEN}  ✓ ATOMIC UPDATE COMPLETE${NC}"
    echo -e "${GREEN}═══════════════════════════════════════════════════════${NC}"

    exit 0
}

# Run main
main "$@"
