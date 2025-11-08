#!/bin/bash

# Installation script for Software Factory 2.0 Branch Name Validation Hook
# This script installs the branch validation pre-commit hook into git repositories
# It can be used by orchestrators when setting up effort infrastructure

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

# Function to get the script directory
get_script_dir() {
    local source="${BASH_SOURCE[0]}"
    while [ -h "$source" ]; do
        local dir="$(cd -P "$(dirname "$source")" && pwd)"
        source="$(readlink "$source")"
        [[ $source != /* ]] && source="$dir/$source"
    done
    echo "$(cd -P "$(dirname "$source")" && pwd)"
}

# Function to install the hook in a git repository
install_hook() {
    local repo_path="${1:-.}"
    local force="${2:-false}"

    # Get absolute path
    repo_path="$(cd "$repo_path" && pwd)"

    # Check if it's a git repository
    if [ ! -d "$repo_path/.git" ]; then
        print_error "Not a git repository: $repo_path"
        return 1
    fi

    print_info "Installing branch validation hook in: $repo_path"

    # Get the hooks directory
    local hooks_dir="$repo_path/.git/hooks"

    # Create hooks directory if it doesn't exist
    mkdir -p "$hooks_dir"

    # Get the path to the validation hook script
    local script_dir=$(get_script_dir)

    # Try new location first
    local validation_hook="$script_dir/../tools/git-commit-hooks/effort-hooks/branch-name-validation.hook"

    # Fall back to old location if new one doesn't exist
    if [ ! -f "$validation_hook" ]; then
        validation_hook="$script_dir/branch-name-validation-hook.sh"
    fi

    # Also check the template directory
    if [ ! -f "$validation_hook" ]; then
        validation_hook="/home/vscode/software-factory-template/tools/git-commit-hooks/effort-hooks/branch-name-validation.hook"
    fi

    if [ ! -f "$validation_hook" ]; then
        print_error "Branch validation hook script not found"
        print_info "Looked in:"
        print_info "  - $script_dir/../tools/git-commit-hooks/effort-hooks/"
        print_info "  - $script_dir/"
        print_info "  - /home/vscode/software-factory-template/tools/git-commit-hooks/effort-hooks/"
        return 1
    fi

    # Check if there's an existing pre-commit hook
    local pre_commit_hook="$hooks_dir/pre-commit"
    local backup_created=false

    if [ -f "$pre_commit_hook" ]; then
        # Read the existing hook to see if it already includes our validation
        if grep -q "Branch Name Validation" "$pre_commit_hook" 2>/dev/null; then
            print_info "Branch validation already installed in pre-commit hook"
            if [ "$force" != "true" ]; then
                return 0
            fi
            print_info "Force flag set, reinstalling..."
        fi

        # Back up the existing hook
        local backup_file="${pre_commit_hook}.backup.$(date +%Y%m%d-%H%M%S)"
        cp "$pre_commit_hook" "$backup_file"
        backup_created=true
        print_info "Backed up existing pre-commit hook to: ${backup_file#$repo_path/}"
    fi

    # Create or update the pre-commit hook
    if [ -f "$pre_commit_hook" ] && [ "$backup_created" = true ]; then
        # Merge with existing hook
        print_info "Merging with existing pre-commit hook..."

        # Create a new hook that calls both
        cat > "${pre_commit_hook}.new" << 'EOF'
#!/bin/bash

# Combined pre-commit hook for Software Factory 2.0
# This hook runs multiple validation checks before allowing commits

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

print_info() {
    echo -e "${BLUE}${BOLD}ℹ️  INFO:${NC} $1"
}

# Track overall validation status
validation_failed=false

EOF

        # Add branch validation first
        echo "" >> "${pre_commit_hook}.new"
        echo "# ============================================" >> "${pre_commit_hook}.new"
        echo "# Branch Name Validation" >> "${pre_commit_hook}.new"
        echo "# ============================================" >> "${pre_commit_hook}.new"
        echo "" >> "${pre_commit_hook}.new"

        # Include the branch validation code (without shebang)
        tail -n +2 "$validation_hook" | sed 's/^main "$@"$/# main function will be called later/' >> "${pre_commit_hook}.new"

        echo "" >> "${pre_commit_hook}.new"
        echo "# Run branch validation" >> "${pre_commit_hook}.new"
        echo "if ! main; then" >> "${pre_commit_hook}.new"
        echo "    validation_failed=true" >> "${pre_commit_hook}.new"
        echo "fi" >> "${pre_commit_hook}.new"
        echo "" >> "${pre_commit_hook}.new"

        # Add the existing hook content (without shebang)
        echo "# ============================================" >> "${pre_commit_hook}.new"
        echo "# Original Pre-commit Hook" >> "${pre_commit_hook}.new"
        echo "# ============================================" >> "${pre_commit_hook}.new"
        echo "" >> "${pre_commit_hook}.new"

        # Extract the main logic from the existing hook (skip shebang and initial setup)
        tail -n +2 "$pre_commit_hook" | grep -v "^set -euo pipefail" | grep -v "^# Color codes" >> "${pre_commit_hook}.new" || true

        echo "" >> "${pre_commit_hook}.new"
        echo "# Exit with appropriate code" >> "${pre_commit_hook}.new"
        echo 'if [ "$validation_failed" = true ]; then' >> "${pre_commit_hook}.new"
        echo "    exit 1" >> "${pre_commit_hook}.new"
        echo "fi" >> "${pre_commit_hook}.new"
        echo "" >> "${pre_commit_hook}.new"
        echo "exit 0" >> "${pre_commit_hook}.new"

        # Replace the old hook with the new one
        mv "${pre_commit_hook}.new" "$pre_commit_hook"
    else
        # Simply copy the validation hook as the pre-commit hook
        print_info "Installing fresh pre-commit hook..."
        cp "$validation_hook" "$pre_commit_hook"
    fi

    # Make the hook executable
    chmod +x "$pre_commit_hook"

    print_success "Branch validation hook installed successfully!"

    # Test the hook
    print_info "Testing branch validation..."
    current_branch="$(cd "$repo_path" && git branch --show-current)"
    if [ -n "$current_branch" ]; then
        print_info "Current branch: $current_branch"

        # Run validation in a subshell to test
        if (cd "$repo_path" && SKIP_BRANCH_VALIDATION=false bash "$pre_commit_hook" 2>/dev/null); then
            print_success "Current branch passes validation"
        else
            print_warning "Current branch may not follow naming conventions (this is expected for main/master)"
        fi
    fi

    return 0
}

# Function to install in multiple repositories
install_in_multiple_repos() {
    local base_path="${1:-.}"
    local pattern="${2:-efforts/phase*/wave*/*}"

    print_info "Installing branch validation hooks in multiple repositories..."
    print_info "Base path: $base_path"
    print_info "Pattern: $pattern"

    local count=0
    local failed=0

    # Find all git repositories matching the pattern
    while IFS= read -r -d '' git_dir; do
        local repo_dir="${git_dir%/.git}"
        echo ""
        if install_hook "$repo_dir" false; then
            ((count++))
        else
            ((failed++))
        fi
    done < <(find "$base_path" -type d -path "$pattern/.git" -prune -print0 2>/dev/null)

    echo ""
    echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    print_info "Installation summary:"
    print_success "Successfully installed in $count repositories"
    if [ $failed -gt 0 ]; then
        print_warning "Failed to install in $failed repositories"
    fi
}

# Main function
main() {
    local target="${1:-.}"
    local mode="${2:-single}"
    local force="${3:-false}"

    case "$mode" in
        single)
            install_hook "$target" "$force"
            ;;
        multiple)
            install_in_multiple_repos "$target" "efforts/phase*/wave*/*"
            ;;
        all)
            # Install in main repository
            echo -e "${BOLD}Installing in main repository...${NC}"
            install_hook "$target" "$force"
            echo ""
            # Install in all effort repositories
            echo -e "${BOLD}Installing in effort repositories...${NC}"
            install_in_multiple_repos "$target" "efforts/phase*/wave*/*"
            ;;
        *)
            print_error "Unknown mode: $mode"
            print_info "Usage: $0 [target_path] [single|multiple|all] [force]"
            exit 1
            ;;
    esac
}

# Show usage if --help is provided
if [ "${1:-}" = "--help" ] || [ "${1:-}" = "-h" ]; then
    echo "Software Factory 2.0 - Branch Validation Hook Installer"
    echo ""
    echo "Usage: $0 [target_path] [mode] [force]"
    echo ""
    echo "Arguments:"
    echo "  target_path  Path to the repository or base directory (default: current directory)"
    echo "  mode         Installation mode: single, multiple, or all (default: single)"
    echo "               - single: Install in a single repository"
    echo "               - multiple: Install in all effort repositories"
    echo "               - all: Install in main repo and all effort repositories"
    echo "  force        Set to 'true' to force reinstallation (default: false)"
    echo ""
    echo "Examples:"
    echo "  $0                                    # Install in current directory"
    echo "  $0 /path/to/repo                      # Install in specific repository"
    echo "  $0 /path/to/project multiple          # Install in all effort repos"
    echo "  $0 /path/to/project all true          # Force install everywhere"
    echo ""
    exit 0
fi

# Run main function
main "$@"