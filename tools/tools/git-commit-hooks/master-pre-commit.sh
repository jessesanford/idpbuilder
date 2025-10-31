#!/bin/bash

# Master Pre-commit Hook for Software Factory 2.0/3.0
# This script combines multiple validation hooks based on repository type and SF version
# It automatically detects SF 2.0 vs SF 3.0 and applies appropriate validation

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

print_info() {
    echo -e "${BLUE}${BOLD}ℹ️  INFO:${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}${BOLD}⚠️  WARNING:${NC} ${YELLOW}$1${NC}"
}

# Function to detect Software Factory version
detect_sf_version() {
    # SF 3.0: Has v3 state files
    if [ -f "$GIT_ROOT/orchestrator-state-v3.json" ] || \
       [ -f "$GIT_ROOT/bug-tracking.json" ] || \
       [ -f "$GIT_ROOT/integration-containers.json" ]; then
        echo "3.0"
        return
    fi

    # SF 2.0: Has old state file
    if [ -f "$GIT_ROOT/orchestrator-state-v3.json" ]; then
        echo "2.0"
        return
    fi

    # Unknown/no SF
    echo "unknown"
}

# Function to detect repository type
detect_repo_type() {
    local remote_url="$(git config --get remote.origin.url 2>/dev/null || true)"
    local repo_basename="$(basename "$GIT_ROOT")"

    # TEMPLATE CHECK: Skip strict validation for template repository
    # Template repos need multiple branches for development/maintenance
    if [ -f "$GIT_ROOT/.template-repository" ]; then
        echo "template"
        return
    fi

    # PRIMARY CHECK: Planning repos have orchestrator state file at root
    # This is the DEFINITIVE marker of a planning repository
    # SF 2.0: orchestrator-state-v3.json
    # SF 3.0: orchestrator-state-v3.json
    if [ -f "$GIT_ROOT/orchestrator-state-v3.json" ] || [ -f "$GIT_ROOT/orchestrator-state-v3.json" ]; then
        echo "planning"
        return
    fi

    # Check if this is the planning repository by name patterns
    # Includes software-factory-template repos and common planning repos
    if [[ "$remote_url" == *"software-factory-template"* ]] || \
       [[ "$remote_url" == *"software-factory-2.0-template"* ]] || \
       [[ "$repo_basename" == "idpbuilder" ]] || \
       [[ "$repo_basename" == "idpbuilder-push" ]] || \
       [[ "$repo_basename" == "idpbuilder-push-oci" ]] || \
       [[ "$repo_basename" == *"-push" ]] || \
       [[ "$repo_basename" == *"-push-oci" ]] || \
       [[ "$repo_basename" == *"-planning" ]] || \
       [[ -f "$GIT_ROOT/PROJECT-IMPLEMENTATION-PLAN.md" && ! -f "$GIT_ROOT/target-repo-config.yaml" ]] || \
       [[ -f "$GIT_ROOT/.planning-repo" ]]; then
        echo "planning"
        return
    fi

    # Check if this is an effort repository
    if [[ "$GIT_ROOT" =~ efforts/ ]] || [ -f "$GIT_ROOT/target-repo-config.yaml" ]; then
        echo "effort"
        return
    fi

    # Check if we're in the efforts directory structure
    if [[ -d "../software-factory-template" ]]; then
        echo "effort"
        return
    fi

    # Default to general repository
    echo "general"
}

# Function to run a hook if it exists
run_hook() {
    local hook_path="$1"
    local hook_name="$2"

    if [ -f "$hook_path" ]; then
        print_info "Running $hook_name..."
        if bash "$hook_path"; then
            print_success "$hook_name passed"
            return 0
        else
            print_error "$hook_name failed"
            return 1
        fi
    else
        # Hook doesn't exist at this path, check if it's in the template
        local template_base="/home/vscode/software-factory-template"
        # Remove "tools/git-commit-hooks/" from hook_path if it starts with it
        local hook_file="${hook_path#tools/git-commit-hooks/}"
        if [ -f "$template_base/tools/git-commit-hooks/$hook_file" ]; then
            print_info "Running $hook_name from template..."
            if bash "$template_base/tools/git-commit-hooks/$hook_file"; then
                print_success "$hook_name passed"
                return 0
            else
                print_error "$hook_name failed"
                return 1
            fi
        fi
        return 0
    fi
}

# Function to run SF 3.0 validation
run_sf3_validation() {
    print_info "Validating SF 3.0 state files..."

    local status=0
    local validator="$GIT_ROOT/tools/validate-state-file.sh"

    # Check if validator exists
    if [ ! -f "$validator" ]; then
        print_error "SF 3.0 validator not found: $validator"
        print_warning "Pre-commit hook expected tools/validate-state-file.sh"
        return 1
    fi

    # Define SF 3.0 state files to validate
    local state_files=(
        "orchestrator-state-v3.json"
        "bug-tracking.json"
        "integration-containers.json"
        "fix-cascade-state.json"
    )

    # Validate each file that exists and is staged for commit
    for file in "${state_files[@]}"; do
        if [ -f "$GIT_ROOT/$file" ]; then
            # Check if file is staged for commit
            if git diff --cached --name-only | grep -q "^${file}$"; then
                if bash "$validator" "$GIT_ROOT/$file"; then
                    print_success "$file validation passed"
                else
                    print_error "$file validation failed"
                    status=1
                fi
            fi
        fi
    done

    # R550: Plan Path Consistency Validation
    if [ -f "$GIT_ROOT/tools/validate-R550-compliance.sh" ]; then
        print_info "Running R550 plan path consistency validation..."
        if bash "$GIT_ROOT/tools/validate-R550-compliance.sh"; then
            print_success "R550 plan path consistency validation passed"
        else
            print_error "R550 plan path consistency validation failed"
            print_warning "See: rule-library/R550-plan-path-consistency-and-discovery.md"
            status=1
        fi
    fi

    return $status
}

# Main function
main() {
    local validation_failed=false
    local repo_type=$(detect_repo_type)
    local sf_version=$(detect_sf_version)

    echo ""
    echo -e "${BOLD}===========================================${NC}"
    echo -e "${BOLD}Software Factory Pre-Commit Validation${NC}"
    echo -e "${BOLD}===========================================${NC}"
    print_info "Software Factory Version: $sf_version"
    print_info "Repository type: $repo_type"
    echo -e "${BOLD}===========================================${NC}"

    # SF 3.0: Use new validation
    if [ "$sf_version" = "3.0" ]; then
        if ! run_sf3_validation; then
            validation_failed=true
        fi

        # Final status for SF 3.0
        echo ""
        echo -e "${BOLD}===========================================${NC}"
        if [ "$validation_failed" = true ]; then
            print_error "SF 3.0 validation failed!"
            echo -e "${RED}${BOLD}⚠️  WARNING - R506 VIOLATION:${NC}"
            echo -e "${RED}${BOLD}Bypassing pre-commit hooks will violate R506 and cause system corruption!${NC}"
            echo ""
            exit 1
        else
            print_success "All SF 3.0 state file validations passed"
            echo -e "${BOLD}===========================================${NC}"
            print_success "All pre-commit validations passed!"
            echo -e "${GREEN}Proceeding with commit...${NC}"
            echo ""
            exit 0
        fi
    fi

    # SF 2.0: Use legacy validation (existing code continues below)
    echo ""
    echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

    # Shared hooks (run for all repository types)
    # R383 SUPREME LAW: Metadata placement validation (FIRST - prevents merge conflicts!)
    if ! run_hook "tools/git-commit-hooks/shared-hooks/metadata-placement-validation.hook" \
                  "R383 Metadata Placement Validation"; then
        validation_failed=true
    fi

    if ! run_hook "tools/git-commit-hooks/shared-hooks/orchestrator-state-validation.hook" \
                  "Orchestrator State Validation"; then
        validation_failed=true
    fi

    # Repository-specific hooks
    case "$repo_type" in
        template)
            # Template repository - allow multi-branch development
            print_info "Template repository - allowing multiple branches for development"
            print_info "Branch validation skipped - this is the Software Factory codebase itself"
            ;;

        planning)
            # Planning repository hooks
            if ! run_hook "tools/git-commit-hooks/planning-hooks/branch-name-validation.hook" \
                          "Planning Branch Name Validation"; then
                validation_failed=true
            fi

            if ! run_hook "tools/git-commit-hooks/planning-hooks/efforts-protection.hook" \
                          "Efforts Directory Protection"; then
                validation_failed=true
            fi
            ;;

        effort)
            # Effort repository hooks
            if ! run_hook "tools/git-commit-hooks/effort-hooks/branch-name-validation.hook" \
                          "Branch Name Validation"; then
                validation_failed=true
            fi

            if ! run_hook "tools/git-commit-hooks/effort-hooks/r251-main-branch-protection.hook" \
                          "R251 Main Branch Protection"; then
                validation_failed=true
            fi
            ;;

        general)
            # General repository - run minimal validation
            print_info "General repository - minimal validation applied"
            ;;
    esac

    echo ""
    echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

    # Final status
    if [ "$validation_failed" = true ]; then
        print_error "Pre-commit validation failed!"
        echo ""
        exit 1
    else
        print_success "All pre-commit validations passed!"
        echo ""
        exit 0
    fi
}

# Run main function
main "$@"
