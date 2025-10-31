#!/bin/bash

# Install strict pre-commit hooks for Software Factory 2.0 planning repositories
# This script installs the hooks that prevent code contamination in planning repos

set -euo pipefail

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
BOLD='\033[1m'
NC='\033[0m'

# Function to print messages
print_error() {
    echo -e "${RED}${BOLD}❌ ERROR:${NC} ${RED}$1${NC}" >&2
}

print_success() {
    echo -e "${GREEN}${BOLD}✅ PROJECT_DONE:${NC} ${GREEN}$1${NC}"
}

print_info() {
    echo -e "${BLUE}${BOLD}ℹ️  INFO:${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}${BOLD}⚠️  WARNING:${NC} ${YELLOW}$1${NC}"
}

# Main installation function
install_hooks() {
    local target_repo="${1:-$(pwd)}"

    # Check if target is a git repository
    if [ ! -d "$target_repo/.git" ]; then
        print_error "Target directory is not a git repository: $target_repo"
        exit 1
    fi

    print_info "Installing strict planning repository hooks to: $target_repo"

    # Get template directory
    local template_dir="$(dirname "$(dirname "$(realpath "$0")")")"

    # Create hooks directory if it doesn't exist
    mkdir -p "$target_repo/.git/hooks"

    # Copy master pre-commit hook
    print_info "Installing master pre-commit hook..."
    cp "$template_dir/tools/git-commit-hooks/master-pre-commit.sh" \
       "$target_repo/.git/hooks/pre-commit"
    chmod +x "$target_repo/.git/hooks/pre-commit"

    # Copy hook files to repository (for hook to find them)
    print_info "Installing hook validation files..."
    mkdir -p "$target_repo/tools/git-commit-hooks"
    cp -r "$template_dir/tools/git-commit-hooks/"* \
          "$target_repo/tools/git-commit-hooks/"

    # Create marker file to identify this as a planning repo
    echo "This is a Software Factory 2.0 planning repository" > "$target_repo/.planning-repo"

    print_success "Hooks installed successfully!"

    echo ""
    echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${GREEN}${BOLD}PLANNING REPOSITORY PROTECTION ACTIVATED${NC}"
    echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    echo "The following branches are now allowed:"
    echo -e "  ${GREEN}• software-factory-2.0${NC} - For all Software Factory planning work"
    echo -e "  ${GREEN}• main${NC} - Only for initial setup (before PROJECT-IMPLEMENTATION-PLAN.md exists)"
    echo ""
    echo -e "${YELLOW}ALL OTHER BRANCHES WILL BE BLOCKED!${NC}"
    echo ""
    echo "Implementation code must go in:"
    echo "  • \$CLAUDE_PROJECT_DIR/efforts/phaseX/waveY/**/*"
    echo "  • \$CLAUDE_PROJECT_DIR/efforts/**/integration"
    echo ""
    echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

# Run installation
if [ $# -gt 0 ]; then
    install_hooks "$1"
else
    install_hooks
fi