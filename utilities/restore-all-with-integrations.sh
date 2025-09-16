#!/bin/bash

# 🏭 SOFTWARE FACTORY 2.0 - ENHANCED RESTORATION WITH ALL INTEGRATIONS
# ═══════════════════════════════════════════════════════════════════════════
# Purpose: Restore all efforts AND integration workspaces (wave, phase, project)
#          from remote branches in the TARGET repository
#
# This enhanced version includes:
# - All regular efforts and splits
# - Wave-level integrations
# - Phase-level integrations
# - Project-level integration
# ═══════════════════════════════════════════════════════════════════════════

set -euo pipefail

# ANSI Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Default values
STATE_FILE="${1:-orchestrator-state.json}"
TARGET_REPO="${2:-https://github.com/jessesanford/idpbuilder.git}"
EFFORTS_ROOT="efforts"
RESTORE_LOG="restoration-with-integrations.log"

# Logging function
log() {
    local level="$1"
    shift
    local message="$*"
    echo -e "${level}${message}${NC}" | tee -a "$RESTORE_LOG"
}

# Clone or update a directory
clone_or_update() {
    local branch="$1"
    local dir_path="$2"
    local name="$3"

    log "${CYAN}" "  Processing: ${name}"
    log "${BLUE}" "  Branch: ${branch}"
    log "${BLUE}" "  Path: ${dir_path}"

    if [ -d "$dir_path/.git" ]; then
        log "${YELLOW}" "  Directory exists, updating..."
        cd "$dir_path"
        git fetch --all
        git checkout "$branch" || git checkout -b "$branch" "origin/$branch"
        git pull origin "$branch" || true
        cd - > /dev/null
        log "${GREEN}" "  ✅ Updated ${name}"
    else
        log "${BLUE}" "  Cloning fresh..."
        mkdir -p "$(dirname "$dir_path")"
        if git clone -b "$branch" --single-branch "$TARGET_REPO" "$dir_path" 2>&1 | tee -a "$RESTORE_LOG"; then
            log "${GREEN}" "  ✅ Cloned ${name}"
        else
            log "${RED}" "  ❌ Failed to clone ${name}"
        fi
    fi
}

# Process phase-level integrations
process_phase_integrations() {
    log "${BOLD}${CYAN}" "\n📋 Processing PHASE-LEVEL integrations..."

    # Check for Phase 1 integration
    local phase1_branch=$(jq -r '.phase1_integration_infrastructure.branch // empty' "$STATE_FILE")
    if [ -n "$phase1_branch" ]; then
        local dir_path="${EFFORTS_ROOT}/phase1/phase-integration-workspace"
        clone_or_update "$phase1_branch" "$dir_path" "Phase 1 Integration"
    else
        # Try alternative location
        phase1_branch="idpbuilder-oci-build-push/phase1-integration"
        if git ls-remote --heads "$TARGET_REPO" "$phase1_branch" 2>/dev/null | grep -q "$phase1_branch"; then
            local dir_path="${EFFORTS_ROOT}/phase1/phase-integration-workspace"
            clone_or_update "$phase1_branch" "$dir_path" "Phase 1 Integration"
        fi
    fi

    # Check for Phase 2 integration (multiple possible branches)
    local phase2_branches=(
        "idpbuilder-oci-build-push/phase2-integration-20250916-033720"
        "idpbuilder-oci-build-push/phase2-integration"
    )

    for branch in "${phase2_branches[@]}"; do
        if git ls-remote --heads "$TARGET_REPO" "$branch" 2>/dev/null | grep -q "$branch"; then
            local dir_path="${EFFORTS_ROOT}/phase2/phase-integration-workspace"
            clone_or_update "$branch" "$dir_path" "Phase 2 Integration"
            break
        fi
    done
}

# Process project-level integration
process_project_integration() {
    log "${BOLD}${CYAN}" "\n📋 Processing PROJECT-LEVEL integration..."

    # Check state file for project integration
    local project_branch=$(jq -r '.project_integration_infrastructure.branch // empty' "$STATE_FILE")
    local project_workspace=$(jq -r '.project_integration_infrastructure.working_copy // empty' "$STATE_FILE")

    if [ -n "$project_branch" ]; then
        # Use the workspace path from state if available, otherwise default
        if [ -n "$project_workspace" ]; then
            local dir_path="$project_workspace"
        else
            local dir_path="${EFFORTS_ROOT}/project/project-integration-workspace"
        fi

        clone_or_update "$project_branch" "$dir_path" "Project Integration"
    else
        # Try known project integration branch
        project_branch="idpbuilder-oci-build-push/project-integration-20250916-152718"
        if git ls-remote --heads "$TARGET_REPO" "$project_branch" 2>/dev/null | grep -q "$project_branch"; then
            local dir_path="${EFFORTS_ROOT}/project/project-integration-workspace"
            clone_or_update "$project_branch" "$dir_path" "Project Integration"
        fi
    fi
}

# Main execution
main() {
    echo "═══════════════════════════════════════════════════════════════════════════" | tee "$RESTORE_LOG"
    echo "🏭 SOFTWARE FACTORY 2.0 - ENHANCED RESTORATION UTILITY" | tee -a "$RESTORE_LOG"
    echo "═══════════════════════════════════════════════════════════════════════════" | tee -a "$RESTORE_LOG"
    echo "Timestamp: $(date '+%Y-%m-%d %H:%M:%S %Z')" | tee -a "$RESTORE_LOG"
    echo "Target Repository: $TARGET_REPO" | tee -a "$RESTORE_LOG"
    echo "═══════════════════════════════════════════════════════════════════════════" | tee -a "$RESTORE_LOG"

    # First run the original restore script for all efforts
    if [ -f "utilities/restore-all-efforts.sh" ]; then
        log "${CYAN}" "\n🔄 Running original effort restoration..."
        bash utilities/restore-all-efforts.sh "$STATE_FILE" "$TARGET_REPO"
    fi

    # Then add phase-level integrations
    process_phase_integrations

    # Finally add project-level integration
    process_project_integration

    # Summary
    log "${BOLD}${GREEN}" "\n✅ RESTORATION COMPLETE!"
    log "${CYAN}" "\n📂 Final structure:"

    # Show phase integrations
    if [ -d "${EFFORTS_ROOT}/phase1/phase-integration-workspace" ]; then
        log "${GREEN}" "  ✅ Phase 1 integration restored"
    fi
    if [ -d "${EFFORTS_ROOT}/phase2/phase-integration-workspace" ]; then
        log "${GREEN}" "  ✅ Phase 2 integration restored"
    fi
    if [ -d "${EFFORTS_ROOT}/project/project-integration-workspace" ]; then
        log "${GREEN}" "  ✅ Project integration restored"
    fi

    # Show tree if available
    if command -v tree &> /dev/null; then
        tree -L 3 "$EFFORTS_ROOT" 2>/dev/null | head -50
    else
        find "$EFFORTS_ROOT" -type d -name "*integration*" | sort
    fi

    log "${CYAN}" "\nRestore log saved to: $RESTORE_LOG"
}

# Run main
main