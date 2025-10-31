#!/bin/bash

# Unified Git Hook Installation Script for Software Factory 2.0
# This script installs the appropriate pre-commit hooks based on repository type
# It uses the new organized hook structure in tools/git-commit-hooks/

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

# Function to detect repository type
detect_repo_type() {
    local repo_path="${1:-.}"
    repo_path="$(cd "$repo_path" && pwd)"

    local remote_url="$(cd "$repo_path" && git config --get remote.origin.url 2>/dev/null || true)"

    # Check if this is the planning repository
    if [[ "$remote_url" == *"software-factory-template"* ]] || [[ "$remote_url" == *"software-factory-2.0-template"* ]]; then
        echo "planning"
        return
    fi

    # Check if this is an effort repository
    if [[ "$repo_path" =~ efforts/ ]] || [ -f "$repo_path/target-repo-config.yaml" ]; then
        echo "effort"
        return
    fi

    # Check if we're in the efforts directory structure
    if [[ -d "$repo_path/../software-factory-template" ]]; then
        echo "effort"
        return
    fi

    # Default to general repository
    echo "general"
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

# Function to install the master hook
install_master_hook() {
    local repo_path="${1:-.}"
    local force="${2:-false}"

    # Get absolute path
    repo_path="$(cd "$repo_path" && pwd)"

    # Check if it's a git repository
    if [ ! -d "$repo_path/.git" ]; then
        print_error "Not a git repository: $repo_path"
        return 1
    fi

    print_info "Installing hooks in: $repo_path"

    # Get the hooks directory
    local hooks_dir="$repo_path/.git/hooks"
    mkdir -p "$hooks_dir"

    # Backup existing pre-commit hook if it exists
    local pre_commit_hook="$hooks_dir/pre-commit"
    if [ -f "$pre_commit_hook" ]; then
        local backup_file="${pre_commit_hook}.backup.$(date +%Y%m%d-%H%M%S)"
        cp "$pre_commit_hook" "$backup_file"
        print_info "Backed up existing pre-commit hook to: ${backup_file#$repo_path/}"
    fi

    # Find the master hook script
    local script_dir=$(get_script_dir)
    local master_hook=""

    # Try multiple locations
    if [ -f "$script_dir/../tools/git-commit-hooks/master-pre-commit.sh" ]; then
        master_hook="$script_dir/../tools/git-commit-hooks/master-pre-commit.sh"
    elif [ -f "/home/vscode/software-factory-template/tools/git-commit-hooks/master-pre-commit.sh" ]; then
        master_hook="/home/vscode/software-factory-template/tools/git-commit-hooks/master-pre-commit.sh"
    else
        print_error "Master pre-commit hook not found"
        return 1
    fi

    # Copy master hook
    cp "$master_hook" "$pre_commit_hook"
    chmod +x "$pre_commit_hook"

    # Also copy the hooks directory structure if not present
    local hooks_base=""
    if [ -d "$script_dir/../tools/git-commit-hooks" ]; then
        hooks_base="$script_dir/../tools/git-commit-hooks"
    elif [ -d "/home/vscode/software-factory-template/tools/git-commit-hooks" ]; then
        hooks_base="/home/vscode/software-factory-template/tools/git-commit-hooks"
    fi

    if [ -n "$hooks_base" ] && [ ! -d "$repo_path/tools/git-commit-hooks" ]; then
        print_info "Copying hook library to repository..."
        mkdir -p "$repo_path/tools/git-commit-hooks"
        cp -r "$hooks_base/"* "$repo_path/tools/git-commit-hooks/"
        print_success "Hook library copied to tools/git-commit-hooks/"
    fi

    # Detect repository type and provide guidance
    local repo_type=$(detect_repo_type "$repo_path")
    print_success "Pre-commit hook installed successfully!"
    print_info "Repository type detected: $repo_type"

    case "$repo_type" in
        planning)
            print_info "Hooks enabled:"
            print_info "  ✅ Efforts directory protection"
            print_info "  ✅ Orchestrator state validation"
            print_info "  ✅ State machine validation"
            ;;
        effort)
            print_info "Hooks enabled:"
            print_info "  ✅ Branch name validation"
            print_info "  ✅ R251 main branch protection"
            print_info "  ✅ Orchestrator state validation"
            ;;
        general)
            print_info "Hooks enabled:"
            print_info "  ✅ Orchestrator state validation (if applicable)"
            ;;
    esac

    return 0
}

# Function to install in multiple repositories
install_in_multiple_repos() {
    local base_path="${1:-.}"
    local pattern="${2:-efforts/phase*/wave*/*}"

    print_info "Installing hooks in multiple repositories..."
    print_info "Base path: $base_path"
    print_info "Pattern: $pattern"

    local count=0
    local failed=0

    # Find all git repositories matching the pattern
    while IFS= read -r -d '' git_dir; do
        local repo_dir="${git_dir%/.git}"
        echo ""
        if install_master_hook "$repo_dir" false; then
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

    echo ""
    echo -e "${BOLD}Software Factory 2.0 - Git Hook Installation${NC}"
    echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""

    case "$mode" in
        single)
            install_master_hook "$target" "$force"
            ;;
        effort)
            install_in_multiple_repos "$target" "efforts/phase*/wave*/*"
            ;;
        all)
            # Install in main repository
            echo -e "${BOLD}Installing in main repository...${NC}"
            install_master_hook "$target" "$force"
            echo ""
            # Install in all effort repositories
            echo -e "${BOLD}Installing in effort repositories...${NC}"
            install_in_multiple_repos "$target" "efforts/phase*/wave*/*"
            ;;
        *)
            print_error "Unknown mode: $mode"
            print_info "Usage: $0 [target_path] [single|effort|all] [force]"
            exit 1
            ;;
    esac

    echo ""
    echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    print_success "Installation complete!"
    echo ""
}

# Show usage if --help is provided
if [ "${1:-}" = "--help" ] || [ "${1:-}" = "-h" ]; then
    echo "Software Factory 2.0 - Git Hook Installation"
    echo ""
    echo "Usage: $0 [target_path] [mode] [force]"
    echo ""
    echo "Arguments:"
    echo "  target_path  Path to the repository or base directory (default: current directory)"
    echo "  mode         Installation mode (default: single)"
    echo "               - single: Install in a single repository"
    echo "               - effort: Install in all effort repositories"
    echo "               - all: Install in main repo and all effort repositories"
    echo "  force        Set to 'true' to force reinstallation (default: false)"
    echo ""
    echo "The script will automatically detect repository type and install appropriate hooks:"
    echo "  • Planning repos: Efforts protection + state validation"
    echo "  • Effort repos: Branch validation + R251 protection + state validation"
    echo "  • General repos: State validation only"
    echo ""
    echo "Examples:"
    echo "  $0                          # Install in current directory"
    echo "  $0 /path/to/repo            # Install in specific repository"
    echo "  $0 /path/to/project effort  # Install in all effort repos"
    echo "  $0 /path/to/project all     # Install everywhere"
    echo ""
    exit 0
fi

# Run main function
main "$@"