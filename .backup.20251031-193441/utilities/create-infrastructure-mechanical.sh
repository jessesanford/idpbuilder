#!/bin/bash

# 🚨🚨🚨 R504 Pre-Infrastructure Planning - Mechanical Infrastructure Creation
# This script ONLY executes pre-planned infrastructure - NO runtime decisions allowed!

set -euo pipefail

CLAUDE_PROJECT_DIR="${CLAUDE_PROJECT_DIR:-/home/vscode/software-factory-template}"
STATE_FILE="$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json"
EFFORT_ID="${1:-}"
TARGET_REPO="${2:-}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to log messages
log() {
    echo -e "${GREEN}[MECHANICAL]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
    exit 504
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1" >&2
}

# Usage function
usage() {
    echo "Usage: $0 <effort_id> [target_repo_url]"
    echo ""
    echo "Mechanically creates infrastructure from pre-planned configuration (R504)"
    echo ""
    echo "Arguments:"
    echo "  effort_id       - The effort ID from pre_planned_infrastructure (e.g., phase1_wave1_effort-name)"
    echo "  target_repo_url - Optional. The target repository URL to clone (defaults to env var TARGET_REPO)"
    echo ""
    echo "Examples:"
    echo "  $0 phase1_wave1_authentication-module"
    echo "  $0 phase1_wave1_database-layer https://github.com/user/repo.git"
    exit 1
}

# Check arguments
if [[ -z "$EFFORT_ID" ]]; then
    usage
fi

# Set target repo
if [[ -z "$TARGET_REPO" ]] && [[ -z "${TARGET_REPO_URL:-}" ]]; then
    error "No target repository specified. Set TARGET_REPO environment variable or provide as second argument."
fi
TARGET_REPO="${TARGET_REPO:-$TARGET_REPO_URL}"

# Check dependencies
check_dependencies() {
    if ! command -v yq &> /dev/null; then
        error "yq is required but not installed. Install with: ./utilities/install-requirements.sh"
    fi

    if ! command -v git &> /dev/null; then
        error "git is required but not installed"
    fi
}

# Validate pre-planned infrastructure exists
validate_pre_planned() {
    log "Validating pre-planned infrastructure (R504)..."

    # Check state file exists
    if [[ ! -f "$STATE_FILE" ]]; then
        error "State file not found: $STATE_FILE"
    fi

    # Check if pre_planned_infrastructure exists and is validated
    local validated=$(yq '.pre_planned_infrastructure.validated' "$STATE_FILE")
    if [[ "$validated" != "true" ]]; then
        error "Pre-planned infrastructure not validated! Run validate-infrastructure-naming.sh first"
    fi

    # Check if effort exists in pre-planned infrastructure
    local effort_exists=$(yq ".pre_planned_infrastructure.efforts.\"$EFFORT_ID\"" "$STATE_FILE")
    if [[ "$effort_exists" == "null" ]]; then
        error "Effort '$EFFORT_ID' not found in pre_planned_infrastructure"
    fi

    # Check if already created
    local already_created=$(yq ".pre_planned_infrastructure.efforts.\"$EFFORT_ID\".created" "$STATE_FILE")
    if [[ "$already_created" == "true" ]]; then
        warning "Infrastructure for '$EFFORT_ID' already created!"
        echo "  To recreate, set .created to false in orchestrator-state-v3.json"
        exit 0
    fi

    log "✅ Pre-planned infrastructure validated"
}

# Read pre-calculated configuration (NO DECISIONS!)
read_configuration() {
    log "📖 Reading pre-calculated configuration for: $EFFORT_ID"

    # Extract ALL values from pre_planned_infrastructure - NO CALCULATIONS!
    FULL_PATH=$(yq ".pre_planned_infrastructure.efforts.\"$EFFORT_ID\".full_path" "$STATE_FILE")
    BRANCH_NAME=$(yq ".pre_planned_infrastructure.efforts.\"$EFFORT_ID\".branch_name" "$STATE_FILE")
    REMOTE_BRANCH=$(yq ".pre_planned_infrastructure.efforts.\"$EFFORT_ID\".remote_branch" "$STATE_FILE")
    TARGET_REMOTE=$(yq ".pre_planned_infrastructure.efforts.\"$EFFORT_ID\".target_remote" "$STATE_FILE")
    PLANNING_REMOTE=$(yq ".pre_planned_infrastructure.efforts.\"$EFFORT_ID\".planning_remote" "$STATE_FILE")
    INTEGRATE_WAVE_EFFORTS_BRANCH=$(yq ".pre_planned_infrastructure.efforts.\"$EFFORT_ID\".integration_branch" "$STATE_FILE")
    PHASE=$(yq ".pre_planned_infrastructure.efforts.\"$EFFORT_ID\".phase" "$STATE_FILE")
    WAVE=$(yq ".pre_planned_infrastructure.efforts.\"$EFFORT_ID\".wave" "$STATE_FILE")

    # Display configuration
    echo ""
    echo "  📁 Path: $FULL_PATH"
    echo "  🌿 Branch: $BRANCH_NAME"
    echo "  🔗 Remote: $REMOTE_BRANCH"
    echo "  🎯 Target Remote: $TARGET_REMOTE"
    echo "  📝 Planning Remote: $PLANNING_REMOTE"
    echo "  🔄 Integration: $INTEGRATE_WAVE_EFFORTS_BRANCH"
    echo "  📊 Phase: $PHASE, Wave: $WAVE"
    echo ""

    log "✅ Configuration loaded - ready for mechanical execution"
}

# Execute infrastructure creation mechanically
create_infrastructure() {
    log "🤖 EXECUTING PRE-PLANNED INFRASTRUCTURE (100% MECHANICAL)..."

    # Step 1: Create directory
    log "Creating directory: $FULL_PATH"
    mkdir -p "$FULL_PATH"

    # Step 2: Clone repository
    log "Cloning repository..."
    if [[ -d "$FULL_PATH/.git" ]]; then
        warning "Git repository already exists at $FULL_PATH - cleaning up"
        rm -rf "$FULL_PATH"
        mkdir -p "$FULL_PATH"
    fi

    git clone "$TARGET_REPO" "$FULL_PATH"
    cd "$FULL_PATH"

    # Step 3: Determine base branch (ONLY runtime decision allowed)
    local BASE_BRANCH=""
    local WAVE_NUM="${WAVE#wave}"  # Extract number from "wave1" -> "1"

    if [[ "$WAVE_NUM" == "1" ]]; then
        # First wave always bases from main
        BASE_BRANCH="origin/main"
        log "Using main as base (first wave)"
    else
        # Check if integration branch exists
        if git ls-remote --heads origin "$INTEGRATE_WAVE_EFFORTS_BRANCH" | grep -q "$INTEGRATE_WAVE_EFFORTS_BRANCH"; then
            BASE_BRANCH="origin/$INTEGRATE_WAVE_EFFORTS_BRANCH"
            log "Using integration branch as base: $INTEGRATE_WAVE_EFFORTS_BRANCH"
        else
            # Fall back to main if integration doesn't exist yet
            BASE_BRANCH="origin/main"
            warning "Integration branch doesn't exist yet, using main as base"
        fi
    fi

    # Step 4: Create and push branch
    log "Creating branch: $BRANCH_NAME from $BASE_BRANCH"
    git checkout -b "$BRANCH_NAME" "$BASE_BRANCH"

    log "Pushing branch to remote..."
    git push -u origin "$BRANCH_NAME"

    # Step 5: Lock git config (R312)
    log "🔒 Locking git config per R312..."
    chmod 444 .git/config

    # Step 6: Install branch validation hook if available
    if [[ -f "$CLAUDE_PROJECT_DIR/utilities/install-branch-validation-hook.sh" ]]; then
        log "Installing branch validation pre-commit hook..."
        bash "$CLAUDE_PROJECT_DIR/utilities/install-branch-validation-hook.sh" "$FULL_PATH" single false
        if [[ $? -eq 0 ]]; then
            log "✅ Branch validation hook installed"
        else
            warning "Branch validation hook installation failed"
        fi
    fi

    # Step 7: Create initial structure files
    log "Creating initial structure files..."
    echo "# $EFFORT_ID" > README.md
    echo "" >> README.md
    echo "Branch: $BRANCH_NAME" >> README.md
    echo "Created: $(date -u +%Y-%m-%dT%H:%M:%SZ)" >> README.md
    echo "Phase: $PHASE, Wave: $WAVE" >> README.md

    git add README.md
    git commit -m "chore: initial infrastructure for $EFFORT_ID [R504]"
    git push

    log "✅ Infrastructure created successfully!"
}

# Update state file
update_state() {
    log "📝 Updating orchestrator-state-v3.json..."

    # Mark as created in pre_planned_infrastructure
    yq -i ".pre_planned_infrastructure.efforts.\"$EFFORT_ID\".created = true" "$STATE_FILE"

    # Add timestamp
    local timestamp=$(date -u +%Y-%m-%dT%H:%M:%SZ)
    yq -i ".pre_planned_infrastructure.efforts.\"$EFFORT_ID\".created_at = \"$timestamp\"" "$STATE_FILE"

    # Update legacy tracking if it exists (for compatibility)
    local effort_name="${EFFORT_ID##*_}"  # Extract effort name from ID
    if yq ".effort_dependencies.\"$effort_name\"" "$STATE_FILE" | grep -qv "null"; then
        yq -i ".effort_dependencies.\"$effort_name\".infrastructure_created = true |
               .effort_dependencies.\"$effort_name\".branch = \"$BRANCH_NAME\" |
               .effort_dependencies.\"$effort_name\".status = \"ready\"" "$STATE_FILE"
    fi

    # Commit state changes
    cd "$CLAUDE_PROJECT_DIR"
    git add orchestrator-state-v3.json
    git commit -m "state: created infrastructure for $EFFORT_ID [R504]" || true
    git push || true

    log "✅ State file updated"
}

# Generate summary
generate_summary() {
    echo ""
    echo "========================================="
    echo "INFRASTRUCTURE CREATION COMPLETE"
    echo "========================================="
    echo ""
    echo "Effort ID: $EFFORT_ID"
    echo "Location: $FULL_PATH"
    echo "Branch: $BRANCH_NAME"
    echo "Remote: $REMOTE_BRANCH"
    echo ""
    echo "Next Steps:"
    echo "1. Spawn agents to work in this infrastructure"
    echo "2. Agents should cd to: $FULL_PATH"
    echo "3. Agents should verify branch: $BRANCH_NAME"
    echo ""
    echo "✅ Infrastructure is ready for use!"
    echo "========================================="
}

# Main execution
main() {
    log "Starting mechanical infrastructure creation (R504)..."

    # Check dependencies
    check_dependencies

    # Validate pre-planned infrastructure
    validate_pre_planned

    # Read configuration (NO DECISIONS!)
    read_configuration

    # Create infrastructure (MECHANICAL EXECUTION!)
    create_infrastructure

    # Update state
    update_state

    # Generate summary
    generate_summary

    log "Mechanical infrastructure creation complete!"
}

# Execute main function
main "$@"