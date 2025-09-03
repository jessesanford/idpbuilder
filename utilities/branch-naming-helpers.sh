#!/bin/bash

# Branch Naming Helper Functions
# These functions read from target-repo-config.yaml to generate consistent branch names
# Including support for project prefixes

# Get the Software Factory root directory
SF_ROOT="${SF_ROOT:-$(git rev-parse --show-toplevel 2>/dev/null || pwd)}"

# Function to get effort branch name
get_effort_branch_name() {
    local phase="$1"
    local wave="$2"
    local effort_name="$3"
    
    # Read format and prefix from config
    local format=$(yq '.branch_naming.effort_format' "$SF_ROOT/target-repo-config.yaml" 2>/dev/null || echo "phase{phase}/wave{wave}/{effort_name}")
    local project_prefix=$(yq '.branch_naming.project_prefix' "$SF_ROOT/target-repo-config.yaml" 2>/dev/null || echo "")
    
    # Add prefix with slash if it exists
    local prefix=""
    if [ -n "$project_prefix" ] && [ "$project_prefix" != "null" ]; then
        prefix="${project_prefix}/"
    fi
    
    # Replace variables
    local branch_name="${format}"
    branch_name="${branch_name//\{prefix\}/$prefix}"
    branch_name="${branch_name//\{phase\}/$phase}"
    branch_name="${branch_name//\{wave\}/$wave}"
    branch_name="${branch_name//\{effort_name\}/$effort_name}"
    
    echo "$branch_name"
}

# Function to get integration branch name
get_integration_branch_name() {
    local phase="$1"
    local wave="$2"
    
    # Read format and prefix from config
    local format=$(yq '.branch_naming.integration_format' "$SF_ROOT/target-repo-config.yaml" 2>/dev/null || echo "phase{phase}/wave{wave}/integration")
    local project_prefix=$(yq '.branch_naming.project_prefix' "$SF_ROOT/target-repo-config.yaml" 2>/dev/null || echo "")
    
    # Add prefix with slash if it exists
    local prefix=""
    if [ -n "$project_prefix" ] && [ "$project_prefix" != "null" ]; then
        prefix="${project_prefix}/"
    fi
    
    # Replace variables
    local branch_name="${format}"
    branch_name="${branch_name//\{prefix\}/$prefix}"
    branch_name="${branch_name//\{phase\}/$phase}"
    branch_name="${branch_name//\{wave\}/$wave}"
    
    echo "$branch_name"
}

# Function to get phase integration branch name
get_phase_integration_branch_name() {
    local phase="$1"
    
    # Read format and prefix from config
    local format=$(yq '.branch_naming.phase_integration_format' "$SF_ROOT/target-repo-config.yaml" 2>/dev/null || echo "phase{phase}/integration")
    local project_prefix=$(yq '.branch_naming.project_prefix' "$SF_ROOT/target-repo-config.yaml" 2>/dev/null || echo "")
    
    # Add prefix with slash if it exists
    local prefix=""
    if [ -n "$project_prefix" ] && [ "$project_prefix" != "null" ]; then
        prefix="${project_prefix}/"
    fi
    
    # Replace variables
    local branch_name="${format}"
    branch_name="${branch_name//\{prefix\}/$prefix}"
    branch_name="${branch_name//\{phase\}/$phase}"
    
    echo "$branch_name"
}

# Function to get split branch name
get_split_branch_name() {
    local original_branch="$1"  # The too-large branch (already includes prefix if configured)
    local split_number="$2"      # 001, 002, 003, etc.
    
    # Split branches are based on the original branch name
    # Format: {original_branch}--split-{number}
    echo "${original_branch}--split-${split_number}"
}

# Wrapper for get_split_branch_name that builds from components
get_split_branch_name_from_components() {
    local phase="$1"
    local wave="$2"
    local effort_name="$3"
    local split_number="$4"
    
    # First get the original branch name with prefix
    local original_branch=$(get_effort_branch_name "$phase" "$wave" "$effort_name")
    
    # Then create the split branch name
    echo "${original_branch}--split-${split_number}"
}

# Function to validate branch name pattern
validate_branch_name() {
    local branch="$1"
    local expected_type="$2"  # "effort", "integration", or "phase_integration"
    
    # Get the project prefix
    local project_prefix=$(yq '.branch_naming.project_prefix' "$SF_ROOT/target-repo-config.yaml" 2>/dev/null || echo "")
    
    # Build expected pattern
    local pattern=""
    if [ -n "$project_prefix" ] && [ "$project_prefix" != "null" ]; then
        pattern="^${project_prefix}/"
    else
        pattern="^"
    fi
    
    case "$expected_type" in
        "effort")
            pattern="${pattern}phase[0-9]+/wave[0-9]+/[a-z0-9-]+$"
            ;;
        "integration")
            pattern="${pattern}phase[0-9]+/wave[0-9]+/integration$"
            ;;
        "phase_integration")
            pattern="${pattern}phase[0-9]+/integration$"
            ;;
        *)
            echo "Unknown branch type: $expected_type" >&2
            return 1
            ;;
    esac
    
    if [[ "$branch" =~ $pattern ]]; then
        return 0
    else
        echo "Branch '$branch' does not match expected pattern: $pattern" >&2
        return 1
    fi
}

# Function to extract phase/wave/effort from branch name
parse_branch_name() {
    local branch="$1"
    
    # Remove project prefix if present
    local project_prefix=$(yq '.branch_naming.project_prefix' "$SF_ROOT/target-repo-config.yaml" 2>/dev/null || echo "")
    if [ -n "$project_prefix" ] && [ "$project_prefix" != "null" ]; then
        branch="${branch#${project_prefix}/}"
    fi
    
    # Extract components
    if [[ "$branch" =~ ^phase([0-9]+)/wave([0-9]+)/(.+)$ ]]; then
        echo "phase=${BASH_REMATCH[1]} wave=${BASH_REMATCH[2]} effort=${BASH_REMATCH[3]}"
    elif [[ "$branch" =~ ^phase([0-9]+)/integration$ ]]; then
        echo "phase=${BASH_REMATCH[1]} type=phase_integration"
    else
        echo "Invalid branch format: $branch" >&2
        return 1
    fi
}

# Export functions for use in other scripts
export -f get_effort_branch_name
export -f get_integration_branch_name
export -f get_phase_integration_branch_name
export -f get_split_branch_name
export -f get_split_branch_name_from_components
export -f validate_branch_name
export -f parse_branch_name