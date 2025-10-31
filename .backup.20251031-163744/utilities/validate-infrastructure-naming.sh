#!/bin/bash

# 🚨🚨🚨 R504 Pre-Infrastructure Planning - Infrastructure Naming/Pathing Validator
# Validates all pre-calculated infrastructure against naming and pathing rules

set -euo pipefail

CLAUDE_PROJECT_DIR="${CLAUDE_PROJECT_DIR:-/home/vscode/software-factory-template}"
STATE_FILE="${1:-$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Validation results tracking
VALIDATION_PASSED=true
ERRORS=()
WARNINGS=()

# Function to log messages
log() {
    echo -e "${GREEN}[VALIDATOR]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
    ERRORS+=("$1")
    VALIDATION_PASSED=false
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1" >&2
    WARNINGS+=("$1")
}

success() {
    echo -e "${GREEN}[✅ PASS]${NC} $1"
}

# Check if state file exists and has pre_planned_infrastructure
check_state_file() {
    if [[ ! -f "$STATE_FILE" ]]; then
        error "State file not found: $STATE_FILE"
        exit 504
    fi

    local has_pre_planned=$(yq '.pre_planned_infrastructure' "$STATE_FILE")
    if [[ "$has_pre_planned" == "null" ]]; then
        error "No pre_planned_infrastructure found in state file. Run pre-calculate-infrastructure.sh first!"
        exit 504
    fi
}

# R313: Validate effort directory structure
validate_r313_paths() {
    log "Validating R313 - Effort directory structure..."

    local effort_count=0
    local paths_valid=true

    # Get all effort paths
    while IFS= read -r effort_id; do
        local path=$(yq ".pre_planned_infrastructure.efforts.\"$effort_id\".full_path" "$STATE_FILE")
        ((effort_count++))

        # R313 Pattern: $CLAUDE_PROJECT_DIR/efforts/phase{N}/wave{N}/{effort-name}/
        if [[ ! "$path" =~ ^.*/efforts/phase[0-9]+/wave[0-9]+/[a-z0-9-]+/?$ ]]; then
            error "R313 violation - Invalid path format for $effort_id: $path"
            error "  Expected: .../efforts/phase{N}/wave{N}/{effort-name}/"
            paths_valid=false
        else
            success "R313 compliant path for $effort_id"
        fi

        # Check for absolute path
        if [[ ! "$path" =~ ^/ ]]; then
            error "Path must be absolute for $effort_id: $path"
            paths_valid=false
        fi

        # Check for forbidden characters
        if [[ "$path" =~ [[:upper:]] ]]; then
            error "Path contains uppercase characters for $effort_id: $path"
            paths_valid=false
        fi

        if [[ "$path" =~ [[:space:]] ]]; then
            error "Path contains spaces for $effort_id: $path"
            paths_valid=false
        fi
    done < <(yq '.pre_planned_infrastructure.efforts | keys | .[]' "$STATE_FILE")

    if [[ $paths_valid == true ]]; then
        success "All $effort_count effort paths comply with R313"
    else
        error "R313 path validation failed"
    fi
}

# R327: Validate cascade branching conventions
validate_r327_branches() {
    log "Validating R327 - Cascade branching conventions..."

    local branch_count=0
    local branches_valid=true

    # Get all effort branches
    while IFS= read -r effort_id; do
        local branch=$(yq ".pre_planned_infrastructure.efforts.\"$effort_id\".branch_name" "$STATE_FILE")
        ((branch_count++))

        # R327 Pattern: {project-prefix}/phase{N}/wave{N}/{effort-name}
        if [[ ! "$branch" =~ ^[a-z0-9-]+/phase[0-9]+/wave[0-9]+/[a-z0-9-]+$ ]]; then
            error "R327 violation - Invalid branch format for $effort_id: $branch"
            error "  Expected: {project-prefix}/phase{N}/wave{N}/{effort-name}"
            branches_valid=false
        else
            success "R327 compliant branch for $effort_id"
        fi

        # Check for forbidden characters
        if [[ "$branch" =~ [[:upper:]] ]]; then
            error "Branch contains uppercase characters for $effort_id: $branch"
            branches_valid=false
        fi

        if [[ "$branch" =~ [[:space:]] ]]; then
            error "Branch contains spaces for $effort_id: $branch"
            branches_valid=false
        fi

        # Validate split pattern
        local split_pattern=$(yq ".pre_planned_infrastructure.efforts.\"$effort_id\".split_pattern" "$STATE_FILE")
        if [[ ! "$split_pattern" =~ ^[a-z0-9-]+--split-$ ]]; then
            warning "Non-standard split pattern for $effort_id: $split_pattern"
        fi
    done < <(yq '.pre_planned_infrastructure.efforts | keys | .[]' "$STATE_FILE")

    # Validate integration branches
    while IFS= read -r integration_id; do
        local branch=$(yq ".pre_planned_infrastructure.integrations.\"$integration_id\".branch_name" "$STATE_FILE")
        ((branch_count++))

        # Integration pattern: {project-prefix}/phase{N}/wave{N}/integration
        # OR: {project-prefix}/phase{N}/integration
        # OR: {project-prefix}/project/integration
        if [[ ! "$branch" =~ ^[a-z0-9-]+/(phase[0-9]+/(wave[0-9]+/)?|project/)integration$ ]]; then
            error "Invalid integration branch format for $integration_id: $branch"
            branches_valid=false
        else
            success "Valid integration branch for $integration_id"
        fi
    done < <(yq '.pre_planned_infrastructure.integrations | keys | .[]' "$STATE_FILE" 2>/dev/null || true)

    if [[ $branches_valid == true ]]; then
        success "All $branch_count branches comply with R327"
    else
        error "R327 branch validation failed"
    fi
}

# Check for path conflicts (duplicates)
check_path_conflicts() {
    log "Checking for path conflicts..."

    local conflicts_found=false
    local paths=()

    # Collect all paths
    while IFS= read -r effort_id; do
        local path=$(yq ".pre_planned_infrastructure.efforts.\"$effort_id\".full_path" "$STATE_FILE")
        paths+=("$path|$effort_id")
    done < <(yq '.pre_planned_infrastructure.efforts | keys | .[]' "$STATE_FILE")

    # Check for duplicates
    local sorted_paths=($(printf '%s\n' "${paths[@]}" | sort))
    local prev=""

    for item in "${sorted_paths[@]}"; do
        local path="${item%%|*}"
        local effort_id="${item##*|}"

        if [[ "$path" == "$prev" ]]; then
            error "Path conflict detected: $path used by multiple efforts"
            conflicts_found=true
        fi
        prev="$path"
    done

    if [[ $conflicts_found == false ]]; then
        success "No path conflicts detected"
    fi
}

# Check for branch conflicts (duplicates)
check_branch_conflicts() {
    log "Checking for branch conflicts..."

    local conflicts_found=false
    local branches=()

    # Collect all branches
    while IFS= read -r effort_id; do
        local branch=$(yq ".pre_planned_infrastructure.efforts.\"$effort_id\".branch_name" "$STATE_FILE")
        branches+=("$branch|$effort_id")
    done < <(yq '.pre_planned_infrastructure.efforts | keys | .[]' "$STATE_FILE")

    # Check for duplicates
    local sorted_branches=($(printf '%s\n' "${branches[@]}" | sort))
    local prev=""

    for item in "${sorted_branches[@]}"; do
        local branch="${item%%|*}"
        local effort_id="${item##*|}"

        if [[ "$branch" == "$prev" ]]; then
            error "Branch conflict detected: $branch used by multiple efforts"
            conflicts_found=true
        fi
        prev="$branch"
    done

    if [[ $conflicts_found == false ]]; then
        success "No branch conflicts detected"
    fi
}

# Validate project prefix consistency
validate_project_prefix() {
    log "Validating project prefix consistency..."

    local project_prefix=$(yq '.pre_planned_infrastructure.project_prefix' "$STATE_FILE")

    if [[ -z "$project_prefix" ]] || [[ "$project_prefix" == "null" ]]; then
        error "No project prefix defined in pre_planned_infrastructure"
        return 1
    fi

    # Check if all branches start with project prefix
    local prefix_valid=true

    while IFS= read -r effort_id; do
        local branch=$(yq ".pre_planned_infrastructure.efforts.\"$effort_id\".branch_name" "$STATE_FILE")

        if [[ ! "$branch" =~ ^${project_prefix}/ ]]; then
            error "Branch does not start with project prefix '$project_prefix': $branch"
            prefix_valid=false
        fi
    done < <(yq '.pre_planned_infrastructure.efforts | keys | .[]' "$STATE_FILE")

    if [[ $prefix_valid == true ]]; then
        success "All branches use consistent project prefix: $project_prefix"
    fi
}

# Validate remote configurations
validate_remote_configs() {
    log "Validating remote configurations..."

    local remotes_valid=true

    while IFS= read -r effort_id; do
        local target_remote=$(yq ".pre_planned_infrastructure.efforts.\"$effort_id\".target_remote" "$STATE_FILE")
        local planning_remote=$(yq ".pre_planned_infrastructure.efforts.\"$effort_id\".planning_remote" "$STATE_FILE")
        local remote_branch=$(yq ".pre_planned_infrastructure.efforts.\"$effort_id\".remote_branch" "$STATE_FILE")

        # Check required remotes
        if [[ "$target_remote" != "target" ]]; then
            warning "Non-standard target remote for $effort_id: $target_remote"
        fi

        if [[ "$planning_remote" != "planning" ]]; then
            warning "Non-standard planning remote for $effort_id: $planning_remote"
        fi

        # Check remote branch format
        if [[ ! "$remote_branch" =~ ^origin/ ]]; then
            error "Remote branch should start with 'origin/' for $effort_id: $remote_branch"
            remotes_valid=false
        fi
    done < <(yq '.pre_planned_infrastructure.efforts | keys | .[]' "$STATE_FILE")

    if [[ $remotes_valid == true ]]; then
        success "All remote configurations are valid"
    fi
}

# Update validation status in state file
update_validation_status() {
    local timestamp=$(date -u +%Y-%m-%dT%H:%M:%SZ)

    if [[ "$VALIDATION_PASSED" == "true" ]]; then
        yq -i ".pre_planned_infrastructure.validated = true |
               .pre_planned_infrastructure.validation_timestamp = \"$timestamp\" |
               .pre_planned_infrastructure.validation.naming_rules_checked = true |
               .pre_planned_infrastructure.validation.path_conflicts_checked = true |
               .pre_planned_infrastructure.validation.branch_conflicts_checked = true |
               .pre_planned_infrastructure.validation.remote_configs_validated = true |
               .pre_planned_infrastructure.validation.last_validated = \"$timestamp\"" "$STATE_FILE"

        log "✅ Validation status updated in state file"
    else
        yq -i ".pre_planned_infrastructure.validated = false |
               .pre_planned_infrastructure.validation_timestamp = null |
               .pre_planned_infrastructure.validation.last_validated = \"$timestamp\"" "$STATE_FILE"

        log "❌ Validation failed - state file marked as invalid"
    fi
}

# Generate validation report
generate_report() {
    echo ""
    echo "========================================="
    echo "INFRASTRUCTURE VALIDATION REPORT"
    echo "========================================="
    echo ""

    if [[ "$VALIDATION_PASSED" == "true" ]]; then
        echo -e "${GREEN}OVERALL STATUS: PASSED${NC}"
    else
        echo -e "${RED}OVERALL STATUS: FAILED${NC}"
    fi

    echo ""
    echo "Checks Performed:"
    echo "  • R313 - Effort directory structure"
    echo "  • R327 - Cascade branching conventions"
    echo "  • Path conflict detection"
    echo "  • Branch conflict detection"
    echo "  • Project prefix consistency"
    echo "  • Remote configuration validation"

    if [[ ${#ERRORS[@]} -gt 0 ]]; then
        echo ""
        echo -e "${RED}ERRORS (${#ERRORS[@]}):${NC}"
        for err in "${ERRORS[@]}"; do
            echo "  ❌ $err"
        done
    fi

    if [[ ${#WARNINGS[@]} -gt 0 ]]; then
        echo ""
        echo -e "${YELLOW}WARNINGS (${#WARNINGS[@]}):${NC}"
        for warn in "${WARNINGS[@]}"; do
            echo "  ⚠️ $warn"
        done
    fi

    echo ""
    echo "========================================="

    if [[ "$VALIDATION_PASSED" == "false" ]]; then
        echo ""
        echo -e "${RED}ACTION REQUIRED:${NC}"
        echo "  1. Fix the errors listed above"
        echo "  2. Re-run pre-calculate-infrastructure.sh"
        echo "  3. Re-run this validation"
        exit 504
    fi
}

# Main execution
main() {
    log "Starting infrastructure naming/pathing validation..."

    # Check state file
    check_state_file

    # Run validations
    validate_r313_paths
    validate_r327_branches
    check_path_conflicts
    check_branch_conflicts
    validate_project_prefix
    validate_remote_configs

    # Update state file
    update_validation_status

    # Generate report
    generate_report

    log "Validation complete!"
}

# Execute main function
main "$@"