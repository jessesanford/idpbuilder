#!/bin/bash
# Pre-commit hook for Software Factory 2.0 - Planning Repository
# This hook is specifically for the main planning/orchestration repository
# It combines all planning-specific validations

set -euo pipefail

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Get the git root directory
GIT_ROOT="$(git rev-parse --show-toplevel)"
cd "$GIT_ROOT"

# Function to print colored messages
print_error() {
    echo -e "${RED}${BOLD}❌ ERROR:${NC} ${RED}$1${NC}" >&2
}

print_success() {
    echo -e "${GREEN}${BOLD}✅ PROJECT_DONE:${NC} ${GREEN}$1${NC}"
}

print_warning() {
    echo -e "${YELLOW}${BOLD}⚠ WARNING:${NC} ${YELLOW}$1${NC}"
}

print_info() {
    echo -e "${BLUE}${BOLD}ℹ INFO:${NC} ${BLUE}$1${NC}"
}

# Function to run a hook and capture its exit code
run_hook() {
    local hook_path="$1"
    local hook_name="$(basename "$hook_path" .hook)"

    if [ -f "$hook_path" ]; then
        print_info "Running $hook_name validation..."
        if "$hook_path"; then
            print_success "$hook_name validation passed"
            return 0
        else
            print_error "$hook_name validation failed"
            return 1
        fi
    else
        print_warning "Hook not found: $hook_path"
        return 0  # Don't fail if hook is missing
    fi
}

# Main execution
echo -e "${BOLD}${BLUE}===========================================${NC}"
echo -e "${BOLD}${BLUE}Software Factory 2.0 - Planning Repository${NC}"
echo -e "${BOLD}${BLUE}Pre-Commit Validation${NC}"
echo -e "${BOLD}${BLUE}===========================================${NC}"

# Track overall status
overall_status=0

# Determine hook locations
HOOK_DIR="$GIT_ROOT/tools/git-commit-hooks"
if [ ! -d "$HOOK_DIR" ]; then
    # Fallback to template directory structure
    HOOK_DIR="$(dirname "$(readlink -f "$0")")/../git-commit-hooks"
fi

# 1. Run R383 metadata placement validation (SUPREME LAW - shared hook)
if ! run_hook "$HOOK_DIR/shared-hooks/metadata-placement-validation.hook"; then
    overall_status=1
fi

# 2. Run orchestrator state validation (shared hook)
if ! run_hook "$HOOK_DIR/shared-hooks/orchestrator-state-validation.hook"; then
    overall_status=1
fi

# 3. Run branch name validation for planning repos (flexible branch names)
if ! run_hook "$HOOK_DIR/planning-hooks/branch-name-validation.hook"; then
    overall_status=1
fi

# 4. Run efforts directory protection
if ! run_hook "$HOOK_DIR/planning-hooks/efforts-protection.hook"; then
    overall_status=1
fi

# Final status report
echo -e "${BOLD}${BLUE}===========================================${NC}"
if [ $overall_status -eq 0 ]; then
    print_success "All pre-commit validations passed!"
    echo -e "${GREEN}Proceeding with commit...${NC}"
else
    print_error "Pre-commit validation failed!"
    echo -e "${RED}Please fix the issues above and try again.${NC}"
    echo -e "${YELLOW}To bypass (NOT RECOMMENDED), use: git commit --no-verify${NC}"
    echo -e "${RED}${BOLD}WARNING: Bypassing will violate R506 and cause system corruption!${NC}"
fi

exit $overall_status