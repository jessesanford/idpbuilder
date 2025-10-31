#!/bin/bash

# Install Software Factory 2.0 Git Hooks in Effort Working Copy
#
# This script installs pre-commit hooks in an effort, split, or integration working copy.
# It is called automatically during infrastructure creation (CREATE_NEXT_INFRASTRUCTURE state)
# and can be run manually to fix missing hooks or update existing hooks.
#
# Usage:
#   install-effort-hooks.sh [effort-directory]
#
# If no directory is specified, uses current working directory.
#
# The script installs:
#   - master-pre-commit.sh (main orchestrator hook)
#   - shared-hooks/* (validation modules)
#
# The master-pre-commit.sh automatically detects repository type and applies
# appropriate validation rules.

set -euo pipefail

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Function to print colored messages
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

# Determine effort directory
EFFORT_DIR="${1:-.}"

# Normalize path to absolute
if [[ "$EFFORT_DIR" != /* ]]; then
    EFFORT_DIR="$(cd "$EFFORT_DIR" && pwd)"
fi

print_info "Installing Software Factory hooks in: $EFFORT_DIR"

# Validate that this is a git repository
if [ ! -d "$EFFORT_DIR/.git" ]; then
    print_error "Not a git repository: $EFFORT_DIR"
    print_error "Expected .git directory not found"
    exit 1
fi

# Determine template directory
# Try multiple possible locations
TEMPLATE_DIR=""

# 1. Check environment variable
if [ -n "${CLAUDE_PROJECT_DIR:-}" ] && [ -d "$CLAUDE_PROJECT_DIR/tools/git-commit-hooks" ]; then
    TEMPLATE_DIR="$CLAUDE_PROJECT_DIR"
# 2. Check standard location
elif [ -d "/workspaces/software-factory-template/tools/git-commit-hooks" ]; then
    TEMPLATE_DIR="/workspaces/software-factory-template"
# 3. Check home directory
elif [ -d "$HOME/software-factory-template/tools/git-commit-hooks" ]; then
    TEMPLATE_DIR="$HOME/software-factory-template"
# 4. Try to find by traversing up from effort directory
elif [ -f "$EFFORT_DIR/../../../../../../tools/git-commit-hooks/master-pre-commit.sh" ]; then
    TEMPLATE_DIR="$(cd "$EFFORT_DIR/../../../../../../" && pwd)"
else
    print_error "Cannot locate Software Factory template directory"
    print_error "Searched locations:"
    print_error "  - \$CLAUDE_PROJECT_DIR/tools/git-commit-hooks"
    print_error "  - /workspaces/software-factory-template/tools/git-commit-hooks"
    print_error "  - \$HOME/software-factory-template/tools/git-commit-hooks"
    print_error "  - Relative paths from effort directory"
    exit 1
fi

print_info "Using template directory: $TEMPLATE_DIR"

# Validate that required hook files exist in template
if [ ! -f "$TEMPLATE_DIR/tools/git-commit-hooks/master-pre-commit.sh" ]; then
    print_error "master-pre-commit.sh not found in template"
    print_error "Expected: $TEMPLATE_DIR/tools/git-commit-hooks/master-pre-commit.sh"
    exit 1
fi

if [ ! -d "$TEMPLATE_DIR/tools/git-commit-hooks/shared-hooks" ]; then
    print_error "shared-hooks directory not found in template"
    print_error "Expected: $TEMPLATE_DIR/tools/git-commit-hooks/shared-hooks/"
    exit 1
fi

# Create hooks directory if it doesn't exist
mkdir -p "$EFFORT_DIR/.git/hooks"

# Install master pre-commit hook
print_info "Installing master-pre-commit.sh..."
cp "$TEMPLATE_DIR/tools/git-commit-hooks/master-pre-commit.sh" \
   "$EFFORT_DIR/.git/hooks/pre-commit"
chmod +x "$EFFORT_DIR/.git/hooks/pre-commit"
print_success "Installed pre-commit hook"

# Create hooks-shared directory for shared validation modules
# The master-pre-commit.sh looks for hooks in:
#   1. tools/git-commit-hooks/shared-hooks/ (if repo has tools dir)
#   2. .git/hooks-shared/ (fallback location)
#   3. $TEMPLATE_DIR/tools/git-commit-hooks/shared-hooks/ (absolute fallback)
#
# We install to .git/hooks-shared/ to ensure they're always available
print_info "Installing shared validation hooks..."
mkdir -p "$EFFORT_DIR/.git/hooks-shared"
cp -r "$TEMPLATE_DIR/tools/git-commit-hooks/shared-hooks/"* \
      "$EFFORT_DIR/.git/hooks-shared/"
print_success "Installed shared hooks"

# Also install planning-hooks and effort-hooks if they exist (for future use)
if [ -d "$TEMPLATE_DIR/tools/git-commit-hooks/planning-hooks" ]; then
    mkdir -p "$EFFORT_DIR/.git/planning-hooks"
    cp -r "$TEMPLATE_DIR/tools/git-commit-hooks/planning-hooks/"* \
          "$EFFORT_DIR/.git/planning-hooks/" 2>/dev/null || true
fi

if [ -d "$TEMPLATE_DIR/tools/git-commit-hooks/effort-hooks" ]; then
    mkdir -p "$EFFORT_DIR/.git/effort-hooks"
    cp -r "$TEMPLATE_DIR/tools/git-commit-hooks/effort-hooks/"* \
          "$EFFORT_DIR/.git/effort-hooks/" 2>/dev/null || true
fi

# Verify installation
echo ""
echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
print_success "Software Factory hooks installed successfully!"
echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo "Installed components:"
echo -e "  ${GREEN}✓${NC} Pre-commit hook: $EFFORT_DIR/.git/hooks/pre-commit"
echo -e "  ${GREEN}✓${NC} Shared hooks: $EFFORT_DIR/.git/hooks-shared/"
echo ""
echo "Enforced rules:"
echo -e "  ${GREEN}•${NC} R383 - Metadata placement validation (SUPREME LAW)"
echo -e "  ${GREEN}•${NC} R343 - Metadata directory standardization"
echo -e "  ${GREEN}•${NC} R506 - Absolute prohibition on pre-commit bypass"
echo -e "  ${GREEN}•${NC} R500 - Branch HEAD tracking synchronization"
echo -e "  ${GREEN}•${NC} Orchestrator state validation (if orchestrator-state-v3.json exists)"
echo ""
echo -e "${YELLOW}⚠️  IMPORTANT:${NC} ${BOLD}NEVER use --no-verify when committing!${NC}"
echo -e "   Using --no-verify violates R506 and causes system-wide corruption."
echo ""
echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# Test the hook installation
print_info "Testing hook installation..."
cd "$EFFORT_DIR"

# Try to run the hook with --help or just test it exists
if bash .git/hooks/pre-commit --version 2>/dev/null || [ $? -eq 0 ]; then
    print_success "Hook is executable and ready"
elif [ -x .git/hooks/pre-commit ]; then
    print_success "Hook is executable (cannot test execution without staged files)"
else
    print_warning "Hook installed but may not be executable"
fi

exit 0
