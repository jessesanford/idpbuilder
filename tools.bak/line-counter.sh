#!/bin/bash
# Line counter for Software Factory 2.0
# Automatically determines correct base branch - no manual specification needed
# Excludes generated code, tests, and documentation

set -euo pipefail

# Initialize variables
BRANCH=""
DETAILED=false
VERBOSE=false
HELP=false
PROJECT_PREFIX=""
PREFIX_SOURCE="none"

# Function to show usage
show_usage() {
    cat << EOF
Usage: $0 [OPTIONS] [BRANCH]

Software Factory 2.0 Line Counter - Smart Base Detection
Automatically determines the correct base branch and counts non-generated lines.

OPTIONS:
    -d, --detailed         Show detailed file breakdown
    -v, --verbose          Show verbose output including what's excluded
    -h, --help            Show this help message

ARGUMENTS:
    BRANCH                 Branch to measure (default: current branch)

EXAMPLES:
    $0                     # Measure current branch against auto-detected base
    $0 feature-branch      # Measure feature-branch against its base
    $0 -d                  # Show detailed breakdown for current branch
    $0 phase2/wave1/api-refactor -d  # Measure specific effort with details
    $0 my-project/phase2/wave1/api-refactor  # Measure with project prefix

AUTO-DETECTION LOGIC:
    The tool automatically determines the correct base branch.
    
    PROJECT PREFIX DETECTION:
      1. Reads target-repo-config.yaml from:
         - Current directory
         - \$CLAUDE_PROJECT_DIR (if set)
         - Parent directories up to orchestrator root
      2. Falls back to pattern detection if config not found
      3. Shows prefix source in output for transparency
    
    Supports optional project prefixes (e.g., my-project/phase1/wave1/effort):
    
    For split branches (--split-NNN):
      - First split (--split-001): base is the SAME BASE as the oversized effort (NOT the effort itself!)
        Example: If effort was based on phase1/integration, split-001 uses phase1/integration too
      - Later splits (--split-002+): base is the previous split
    
    For effort branches (phase*/wave*/effort-name):
      - Phase 1 efforts: base is main/master
      - Later phase efforts: base is previous phase integration
      - Within wave: base is previous wave integration if exists
    
    For integration branches (final merges back to main):
      - Wave integration (phaseN/waveM-integration): base is main
      - Phase integration (phaseN-integration): base is main
      Note: Integration branches with timestamps are supported
    
    Fallback: Uses git merge-base with main/master

EXCLUDED PATTERNS:
    - Generated code: *.pb.go, *_generated.go, zz_generated*, *.gen.go
    - Vendored code: vendor/*
    - Documentation: *.md
    - Test files: *_test.go
    - CRD YAML: *.crd.yaml, *.crd.yml
    - Build artifacts: bin/*, dist/*, build/*

BASE BRANCH DETECTION:
    No need to specify base branch - automatically determined from branch name!
EOF
}

# Function to find and read project prefix from target-repo-config.yaml
find_project_prefix() {
    local config_file=""
    local prefix=""
    
    # Method 1: Check current directory
    if [ -f "target-repo-config.yaml" ]; then
        config_file="target-repo-config.yaml"
        PREFIX_SOURCE="current directory"
    # Method 2: Check CLAUDE_PROJECT_DIR if set
    elif [ -n "${CLAUDE_PROJECT_DIR:-}" ] && [ -f "$CLAUDE_PROJECT_DIR/target-repo-config.yaml" ]; then
        config_file="$CLAUDE_PROJECT_DIR/target-repo-config.yaml"
        PREFIX_SOURCE="CLAUDE_PROJECT_DIR"
    else
        # Method 3: Search up from current directory to find orchestrator root
        local search_dir="$(pwd)"
        while [ "$search_dir" != "/" ]; do
            # Check for orchestrator-state.json as a marker of the orchestrator root
            if [ -f "$search_dir/orchestrator-state.json" ] || [ -f "$search_dir/orchestrator-state.json.example" ]; then
                if [ -f "$search_dir/target-repo-config.yaml" ]; then
                    config_file="$search_dir/target-repo-config.yaml"
                    PREFIX_SOURCE="orchestrator root ($search_dir)"
                    break
                fi
            fi
            # Also check for the config file directly
            if [ -f "$search_dir/target-repo-config.yaml" ]; then
                config_file="$search_dir/target-repo-config.yaml"
                PREFIX_SOURCE="parent directory ($search_dir)"
                break
            fi
            search_dir="$(dirname "$search_dir")"
        done
    fi
    
    # If config file found, try to read the prefix
    if [ -n "$config_file" ]; then
        # Check if yq is available
        if command -v yq &> /dev/null; then
            prefix=$(jq '.branch_naming.project_prefix // ""' "$config_file" 2>/dev/null || echo "")
            if [ -n "$prefix" ]; then
                PROJECT_PREFIX="$prefix"
                [ "$VERBOSE" = true ] && echo "✓ Found project prefix '$PROJECT_PREFIX' from config in $PREFIX_SOURCE"
                return 0
            else
                [ "$VERBOSE" = true ] && echo "ℹ Config found in $PREFIX_SOURCE but project_prefix is empty"
                PREFIX_SOURCE="config (empty)"
                return 0
            fi
        else
            # Try grep as fallback if yq not available
            local prefix_line=$(grep "project_prefix:" "$config_file" 2>/dev/null || echo "")
            if [ -n "$prefix_line" ]; then
                # Extract value after colon, handling both quoted and unquoted values
                local value="${prefix_line#*: }"
                # Remove leading/trailing spaces and quotes
                value="${value#"${value%%[![:space:]]*}"}"  # Remove leading whitespace
                value="${value%"${value##*[![:space:]]}"}"  # Remove trailing whitespace
                value="${value#\"}"  # Remove leading quote
                value="${value%\"}"  # Remove trailing quote
                
                if [ -n "$value" ]; then
                    PROJECT_PREFIX="$value"
                    [ "$VERBOSE" = true ] && echo "✓ Found project prefix '$PROJECT_PREFIX' from config in $PREFIX_SOURCE (using grep fallback)"
                else
                    [ "$VERBOSE" = true ] && echo "ℹ Config found in $PREFIX_SOURCE but project_prefix is empty (using grep fallback)"
                    PREFIX_SOURCE="config (empty)"
                fi
                return 0
            fi
            [ "$VERBOSE" = true ] && echo "⚠ Config found but unable to parse project_prefix"
        fi
    else
        [ "$VERBOSE" = true ] && echo "ℹ No target-repo-config.yaml found, will use pattern detection"
    fi
    
    PREFIX_SOURCE="pattern detection"
    return 1
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -d|--detailed)
            DETAILED=true
            shift
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -h|--help)
            show_usage
            exit 0
            ;;
        -*)
            echo "Error: Unknown option $1"
            echo "Use -h or --help for usage"
            exit 1
            ;;
        *)
            # This is the branch name
            if [ -z "$BRANCH" ]; then
                BRANCH="$1"
            else
                echo "Error: Multiple branch names specified"
                echo "Use -h or --help for usage"
                exit 1
            fi
            shift
            ;;
    esac
done

# Auto-detect current branch if not specified
if [ -z "$BRANCH" ]; then
    BRANCH=$(git branch --show-current 2>/dev/null || echo "")
    if [ -z "$BRANCH" ]; then
        echo "Error: Not in a git repository or cannot detect current branch"
        echo "Please specify branch as argument"
        exit 1
    fi
    [ "$VERBOSE" = true ] && echo "Using current branch: $BRANCH"
fi

# Verify branch exists
if ! git rev-parse --verify "$BRANCH" >/dev/null 2>&1; then
    echo "Error: Branch '$BRANCH' does not exist"
    echo ""
    echo "Available local branches:"
    git branch | head -10
    echo ""
    echo "Available remote branches matching pattern:"
    git branch -r | grep -E "(phase[0-9]+|split)" | head -10 || echo "  (none found)"
    echo ""
    echo "Tip: If the branch exists remotely, try:"
    echo "  git fetch origin"
    echo "  git checkout -b '$BRANCH' origin/'$BRANCH'"
    exit 1
fi

# Try to find project prefix from configuration
find_project_prefix

# Function to check if branch exists (locally or remote)
branch_exists() {
    local branch="$1"
    git rev-parse --verify "$branch" >/dev/null 2>&1 || \
    git rev-parse --verify "origin/$branch" >/dev/null 2>&1
}

# Function to find branch with possible timestamp suffix
find_branch_with_suffix() {
    local base_pattern="$1"
    local found_branch=""
    
    # First try exact match
    if branch_exists "$base_pattern"; then
        echo "$base_pattern"
        return 0
    fi
    
    # Escape special regex characters in the pattern (especially slashes)
    local escaped_pattern=$(echo "$base_pattern" | sed 's/[[\.*^$(){}?+|]/\\&/g')
    
    # Try to find branch with timestamp suffix (format: -YYYYMMDD-HHMMSS)
    # Search both local and remote branches
    found_branch=$(git branch -a | grep -E "${escaped_pattern}(-[0-9]{8}-[0-9]{6})?$" | head -1 | sed 's/^[* ]*//' | sed 's|^remotes/||')
    
    if [ -n "$found_branch" ]; then
        # Clean up the branch name (remove origin/ if present for consistency)
        found_branch="${found_branch#origin/}"
        [ "$VERBOSE" = true ] && echo "  Found branch with suffix: $found_branch" >&2
        echo "$found_branch"
        return 0
    fi
    
    return 1
}

# Function to get branch reference (local or origin/)
get_branch_ref() {
    local branch="$1"
    
    # First try the exact branch name
    if git rev-parse --verify "$branch" >/dev/null 2>&1; then
        echo "$branch"
    elif git rev-parse --verify "origin/$branch" >/dev/null 2>&1; then
        echo "origin/$branch"
    else
        # If not found, try to find with suffix
        local found_branch=$(find_branch_with_suffix "$branch")
        if [ -n "$found_branch" ]; then
            if git rev-parse --verify "$found_branch" >/dev/null 2>&1; then
                echo "$found_branch"
            elif git rev-parse --verify "origin/$found_branch" >/dev/null 2>&1; then
                echo "origin/$found_branch"
            fi
        fi
    fi
}

# Smart base branch detection based on Software Factory conventions
detect_base_branch() {
    local current="$1"
    local base=""
    
    [ "$VERBOSE" = true ] && echo "Analyzing branch pattern: $current" >&2
    
    # If we have a configured prefix, use it to parse more accurately
    local detected_prefix=""
    if [ -n "$PROJECT_PREFIX" ] && [ "$PREFIX_SOURCE" != "pattern detection" ]; then
        [ "$VERBOSE" = true ] && echo "  Using configured prefix: '$PROJECT_PREFIX' (from $PREFIX_SOURCE)" >&2
        # Strip the prefix from the branch name if it starts with it
        if [[ "$current" =~ ^${PROJECT_PREFIX}/(.*) ]]; then
            local branch_without_prefix="${BASH_REMATCH[1]}"
            detected_prefix="${PROJECT_PREFIX}/"
            [ "$VERBOSE" = true ] && echo "  Branch without prefix: $branch_without_prefix" >&2
        else
            local branch_without_prefix="$current"
        fi
    else
        # No configured prefix, will need to detect from patterns
        local branch_without_prefix="$current"
    fi
    
    # Pattern 1: Split branches (e.g., phase1/wave1/api--split-001 or phase1/wave1/api-split-001 or project/phase1/wave1/api--split-001)
    # Support both --split- (double dash) and -split- (single dash) formats for compatibility
    if [[ "$current" =~ (.*)(--|-split-)([0-9]+) ]]; then
        local effort_base="${BASH_REMATCH[1]}"
        local split_delimiter="${BASH_REMATCH[2]}"  # Either '--split-' or '-split-'
        local split_num="${BASH_REMATCH[3]}"
        
        # If we don't have a configured prefix, try to detect it from the pattern
        local project_prefix=""
        if [ -z "$detected_prefix" ] && [[ "$effort_base" =~ ^([^/]+/)phase[0-9]+/wave[0-9]+/ ]]; then
            project_prefix="${BASH_REMATCH[1]}"
            [ "$VERBOSE" = true ] && echo "  Detected project-prefixed split branch (pattern detection)" >&2
            [ "$VERBOSE" = true ] && echo "  Project prefix: ${project_prefix%/}" >&2
        elif [ -n "$detected_prefix" ]; then
            project_prefix="$detected_prefix"
        fi
        
        [ "$VERBOSE" = true ] && echo "  Detected split branch: split #$split_num of $effort_base" >&2
        
        if [ "$split_num" = "001" ] || [ "$split_num" = "1" ]; then
            # First split - CRITICAL: Use SAME BASE as the oversized effort, NOT the effort itself!
            # The oversized effort branch contains ALL the too-large code.
            # Split-001 must start clean from the same integration base.
            
            # Parse the effort base to determine what IT was based on
            if [[ "$effort_base" =~ ^(.*/)phase([0-9]+)/wave([0-9]+)/([^/]+)$ ]]; then
                local effort_prefix="${BASH_REMATCH[1]}"
                local effort_phase="${BASH_REMATCH[2]}"
                local effort_wave="${BASH_REMATCH[3]}"
                local effort_name="${BASH_REMATCH[4]}"
                
                # Determine what the oversized effort was based on (R308 incremental base)
                if [ "$effort_phase" = "1" ] && [ "$effort_wave" = "1" ]; then
                    # Phase 1, Wave 1 efforts are based on main
                    base=$(find_main_branch)
                else
                    # Other efforts use phase/wave integration as base
                    local integration_base="${effort_prefix}phase${effort_phase}/wave${effort_wave}/integration"
                    base=$(find_branch_with_suffix "$integration_base")
                    
                    if [ -z "$base" ]; then
                        # Try phase integration if wave integration doesn't exist
                        integration_base="${effort_prefix}phase${effort_phase}/integration"
                        base=$(find_branch_with_suffix "$integration_base")
                    fi
                    
                    if [ -z "$base" ]; then
                        # Fallback to main if no integration found
                        base=$(find_main_branch)
                    fi
                fi
                
                [ "$VERBOSE" = true ] && echo "  🔴 CRITICAL: First split uses SAME BASE as oversized effort" >&2
                [ "$VERBOSE" = true ] && echo "  Oversized effort '$effort_base' was based on: $base" >&2
                [ "$VERBOSE" = true ] && echo "  Split-001 will measure against: $base (NOT $effort_base)" >&2
            else
                # Fallback for non-standard naming - use the effort as base (old behavior)
                base="$effort_base"
                [ "$VERBOSE" = true ] && echo "  ⚠️ WARNING: Non-standard naming, using effort as base: $base" >&2
            fi
        else
            # Later splits - base is previous split (use same delimiter format as current branch)
            local prev_num=$((10#$split_num - 1))
            base="${effort_base}${split_delimiter}$(printf "%03d" $prev_num)"
            [ "$VERBOSE" = true ] && echo "  Split #$split_num - base is previous split: $base" >&2
        fi
        
    # Pattern 2: Phase/Wave/Effort branches (e.g., phase2/wave1/api-refactor or project/phase2/wave1/api-refactor)
    elif [[ "$current" =~ ^(([^/]+)/)?phase([0-9]+)/wave([0-9]+)/([^/]+)$ ]]; then
        local project_prefix=""
        local phase=""
        local wave=""
        local effort=""
        
        # Use configured prefix if available, otherwise detect from pattern
        if [ -n "$detected_prefix" ]; then
            # We already know the prefix, parse without it
            if [[ "$branch_without_prefix" =~ ^phase([0-9]+)/wave([0-9]+)/([^/]+)$ ]]; then
                phase="${BASH_REMATCH[1]}"
                wave="${BASH_REMATCH[2]}"
                effort="${BASH_REMATCH[3]}"
                project_prefix="${detected_prefix%/}"  # Remove trailing slash
            fi
        else
            # Pattern detection fallback
            project_prefix="${BASH_REMATCH[2]}"
            phase="${BASH_REMATCH[3]}"
            wave="${BASH_REMATCH[4]}"
            effort="${BASH_REMATCH[5]}"
            
            if [ -n "$project_prefix" ]; then
                [ "$VERBOSE" = true ] && echo "  Detected branch pattern: project-prefixed effort branch (pattern detection)" >&2
                [ "$VERBOSE" = true ] && echo "  Project prefix: $project_prefix" >&2
            fi
        fi
        
        [ "$VERBOSE" = true ] && echo "  Detected effort: Phase $phase, Wave $wave, Effort: $effort" >&2
        
        # For efforts in wave 2+, try previous wave's integration branch
        if [ "$wave" -gt 1 ]; then
            local prev_wave=$((wave - 1))
            local prev_wave_integration="phase${phase}/wave${prev_wave}-integration"
            local found_prev_wave_integration=""
            
            [ "$VERBOSE" = true ] && echo "  Looking for previous wave integration: $prev_wave_integration" >&2
            
            if [ -n "$project_prefix" ]; then
                # Try with prefix first (may have timestamp suffix)
                local prefixed_prev_integration="${project_prefix}/${prev_wave_integration}"
                found_prev_wave_integration=$(find_branch_with_suffix "$prefixed_prev_integration")
                if [ -z "$found_prev_wave_integration" ]; then
                    # Try without prefix as fallback
                    found_prev_wave_integration=$(find_branch_with_suffix "$prev_wave_integration")
                fi
            else
                # No prefix, just try the plain branch name
                found_prev_wave_integration=$(find_branch_with_suffix "$prev_wave_integration")
            fi
            
            if [ -n "$found_prev_wave_integration" ]; then
                base="$found_prev_wave_integration"
                [ "$VERBOSE" = true ] && echo "  Found previous wave integration branch: $base" >&2
            else
                [ "$VERBOSE" = true ] && echo "  Previous wave integration not found, will check other options" >&2
            fi
        fi
        
        # If we didn't find a base yet and phase > 1, check for previous phase integration
        if [ -z "$base" ] && [ "$phase" -gt 1 ]; then
            # For phase 2+, check for previous phase integration
            local prev_phase=$((phase - 1))
            local phase_integration="phase${prev_phase}-integration"
            local found_phase_integration=""
            
            if [ -n "$project_prefix" ]; then
                # Try with prefix first (may have timestamp suffix)
                local prefixed_phase="${project_prefix}/${phase_integration}"
                found_phase_integration=$(find_branch_with_suffix "$prefixed_phase")
                if [ -z "$found_phase_integration" ]; then
                    # Try without prefix as fallback
                    found_phase_integration=$(find_branch_with_suffix "$phase_integration")
                fi
            else
                # No prefix, just try the plain branch name
                found_phase_integration=$(find_branch_with_suffix "$phase_integration")
            fi
            
            if [ -n "$found_phase_integration" ]; then
                base="$found_phase_integration"
                [ "$VERBOSE" = true ] && echo "  Using previous phase integration: $base" >&2
            fi
        fi
        
        # Fallback to main/master for phase 1 or if no integration found
        if [ -z "$base" ]; then
            for candidate in main master; do
                if branch_exists "$candidate"; then
                    base="$candidate"
                    [ "$VERBOSE" = true ] && echo "  Using default base: $base" >&2
                    break
                fi
            done
        fi
        
    # Pattern 3: Integration branches (with optional prefix) - Handle both formats:
    # - phase2-wave2-integration (hyphen format)
    # - phase2/wave2-integration (slash format)
    elif [[ "$current" =~ ^(([^/]+)/)?phase([0-9]+)[-/]wave([0-9]+)-integration(-[0-9]{8}-[0-9]{6})?$ ]]; then
        local project_prefix=""
        local phase=""
        local wave=""
        
        # Use configured prefix if available, otherwise detect from pattern
        if [ -n "$detected_prefix" ]; then
            # We already know the prefix, parse without it - handle both formats
            if [[ "$branch_without_prefix" =~ ^phase([0-9]+)[-/]wave([0-9]+)-integration ]]; then
                phase="${BASH_REMATCH[1]}"
                wave="${BASH_REMATCH[2]}"
                project_prefix="${detected_prefix%/}"  # Remove trailing slash
            fi
        else
            # Pattern detection fallback
            project_prefix="${BASH_REMATCH[2]}"
            phase="${BASH_REMATCH[3]}"
            wave="${BASH_REMATCH[4]}"
            
            if [ -n "$project_prefix" ]; then
                [ "$VERBOSE" = true ] && echo "  Detected branch pattern: project-prefixed wave integration (pattern detection)" >&2
                [ "$VERBOSE" = true ] && echo "  Project prefix: $project_prefix" >&2
            fi
        fi
        
        [ "$VERBOSE" = true ] && echo "  Detected wave integration: Phase $phase, Wave $wave" >&2
        
        # Wave integration branches always compare against main for final integration
        # This is because wave integrations merge all work from the wave back to main
        for candidate in main master; do
            if branch_exists "$candidate"; then
                base="$candidate"
                [ "$VERBOSE" = true ] && echo "  Wave integration - using main branch: $base" >&2
                break
            fi
        done
        
    elif [[ "$current" =~ ^(([^/]+)/)?phase([0-9]+)-integration(-[0-9]{8}-[0-9]{6})?$ ]]; then
        local project_prefix=""
        local phase=""
        
        # Use configured prefix if available, otherwise detect from pattern
        if [ -n "$detected_prefix" ]; then
            # We already know the prefix, parse without it
            if [[ "$branch_without_prefix" =~ ^phase([0-9]+)-integration ]]; then
                phase="${BASH_REMATCH[1]}"
                project_prefix="${detected_prefix%/}"  # Remove trailing slash
            fi
        else
            # Pattern detection fallback
            project_prefix="${BASH_REMATCH[2]}"
            phase="${BASH_REMATCH[3]}"
            
            if [ -n "$project_prefix" ]; then
                [ "$VERBOSE" = true ] && echo "  Detected branch pattern: project-prefixed phase integration (pattern detection)" >&2
                [ "$VERBOSE" = true ] && echo "  Project prefix: $project_prefix" >&2
            fi
        fi
        
        [ "$VERBOSE" = true ] && echo "  Detected phase integration: Phase $phase" >&2
        
        # Phase integrations always compare against main
        # This is the final integration of a complete phase
        for candidate in main master; do
            if branch_exists "$candidate"; then
                base="$candidate"
                [ "$VERBOSE" = true ] && echo "  Phase integration - using main branch: $base" >&2
                break
            fi
        done
    fi
    
    # Fallback: Try git merge-base with common branches
    if [ -z "$base" ]; then
        [ "$VERBOSE" = true ] && echo "  No pattern match - trying merge-base detection" >&2
        
        for candidate in main master develop development; do
            if branch_exists "$candidate"; then
                # Check if there's a merge base
                local merge_base=$(git merge-base "$current" "$(get_branch_ref "$candidate")" 2>/dev/null || echo "")
                if [ -n "$merge_base" ]; then
                    base="$candidate"
                    [ "$VERBOSE" = true ] && echo "  Found merge-base with: $base" >&2
                    break
                fi
            fi
        done
    fi
    
    # Final fallback to main/master
    if [ -z "$base" ]; then
        for candidate in main master; do
            if branch_exists "$candidate"; then
                base="$candidate"
                [ "$VERBOSE" = true ] && echo "  Final fallback to: $base" >&2
                break
            fi
        done
    fi
    
    # Return the base branch (with origin/ prefix if needed)
    if [ -n "$base" ]; then
        echo "$(get_branch_ref "$base")"
    else
        echo ""
    fi
}

# Detect base branch
BASE=$(detect_base_branch "$BRANCH")

if [ -z "$BASE" ]; then
    echo "Error: Could not determine base branch for '$BRANCH'"
    echo ""
    echo "Debugging information:"
    echo "  Current branch: $BRANCH"
    
    # Check if it matches any known patterns
    if [[ "$BRANCH" =~ ^(([^/]+)/)?phase([0-9]+)/wave([0-9]+)/([^/]+)$ ]]; then
        echo "  ✓ Matches effort pattern"
        echo "    - Project prefix: ${BASH_REMATCH[2]:-none}"
        echo "    - Phase: ${BASH_REMATCH[3]}"
        echo "    - Wave: ${BASH_REMATCH[4]}"
        echo "    - Effort: ${BASH_REMATCH[5]}"
        echo ""
        echo "  ⚠ Pattern matched but base branch not found"
        echo "  Possible causes:"
        echo "    - Missing main/master branch"
        echo "    - No common ancestor with main/master"
    else
        echo "  ✗ Does not match expected patterns"
    fi
    
    echo ""
    echo "Branch naming conventions:"
    echo "  - Efforts: phaseN/waveM/effort-name"
    echo "  - Efforts with prefix: project-name/phaseN/waveM/effort-name"
    echo "  - Splits: phaseN/waveM/effort--split-NNN"
    echo "  - Splits with prefix: project-name/phaseN/waveM/effort--split-NNN"
    echo "  - Wave Integration: phaseN/waveM-integration or phaseN-waveM-integration"
    echo "  - Phase Integration: phaseN-integration"
    echo "  - With prefix: project-name/phaseN/waveM-integration"
    echo "  - With timestamp: any-pattern-YYYYMMDD-HHMMSS"
    echo ""
    echo "Run with -v flag for verbose pattern matching details."
    exit 1
fi

# Show what we're measuring
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 Line Counter - Software Factory 2.0"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📌 Analyzing branch: $BRANCH"
echo "🎯 Detected base:    $BASE"

# Show project prefix information if relevant
if [ -n "$PROJECT_PREFIX" ]; then
    echo "🏷️  Project prefix:  $PROJECT_PREFIX (from $PREFIX_SOURCE)"
elif [ "$PREFIX_SOURCE" = "config (empty)" ]; then
    echo "🏷️  Project prefix:  (none configured)"
elif [ "$PREFIX_SOURCE" = "pattern detection" ]; then
    # Only show if verbose and we tried to find config
    [ "$VERBOSE" = true ] && echo "🏷️  Project prefix:  (using pattern detection)"
fi

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Define exclusion patterns
EXCLUSIONS=(
    ':(exclude)*.pb.go'
    ':(exclude)*_generated.go'
    ':(exclude)zz_generated*'
    ':(exclude)*.gen.go'
    ':(exclude)vendor/*'
    ':(exclude)*.md'
    ':(exclude)*_test.go'
    ':(exclude)*.crd.yaml'
    ':(exclude)*.crd.yml'
    ':(exclude)bin/*'
    ':(exclude)dist/*'
    ':(exclude)build/*'
)

# Show exclusions in verbose mode
if [ "$VERBOSE" = true ]; then
    echo ""
    echo "Excluding patterns:"
    for pattern in "${EXCLUSIONS[@]}"; do
        echo "  - ${pattern#:(exclude)}"
    done
    echo ""
fi

# Get the diff stats
DIFF_OUTPUT=$(git diff --stat "$BASE..$BRANCH" -- "${EXCLUSIONS[@]}" 2>/dev/null || echo "")

if [ -z "$DIFF_OUTPUT" ]; then
    echo ""
    echo "ℹ️  No changes found between $BASE and $BRANCH"
    echo "   (This is expected for new branches or identical content)"
    echo ""
    echo "✅ Total non-generated lines: 0"
    exit 0
fi

# Extract total from last line
TOTAL_LINES=$(echo "$DIFF_OUTPUT" | tail -1 | grep -oE '[0-9]+ insertion' | grep -oE '[0-9]+' || echo "0")
TOTAL_DELETIONS=$(echo "$DIFF_OUTPUT" | tail -1 | grep -oE '[0-9]+ deletion' | grep -oE '[0-9]+' || echo "0")

# Calculate net lines
NET_LINES=$((TOTAL_LINES - TOTAL_DELETIONS))

# Display results
echo ""
echo "📈 Line Count Summary:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  Insertions:  +$TOTAL_LINES"
echo "  Deletions:   -$TOTAL_DELETIONS"
echo "  Net change:   $NET_LINES"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Size limit checks
if [ "$TOTAL_LINES" -gt 800 ]; then
    echo "🚨 HARD LIMIT VIOLATION: Branch exceeds 800 lines!"
    echo "   This branch MUST be split immediately."
    echo ""
    echo "✅ Total non-generated lines: $TOTAL_LINES"
elif [ "$TOTAL_LINES" -gt 700 ]; then
    echo "⚠️  WARNING: Branch exceeds 700 line soft limit!"
    echo "   Consider splitting into multiple efforts."
    echo ""
    echo "✅ Total non-generated lines: $TOTAL_LINES"
else
    echo "✅ Total non-generated lines: $TOTAL_LINES"
fi

# Show detailed breakdown if requested
if [ "$DETAILED" = true ]; then
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "📝 Detailed File Breakdown:"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    # Show all but the last line (summary)
    echo "$DIFF_OUTPUT" | head -n -1
fi

# Exit with appropriate code
if [ "$TOTAL_LINES" -gt 800 ]; then
    exit 2  # Hard limit violation
elif [ "$TOTAL_LINES" -gt 700 ]; then
    exit 1  # Soft limit warning
else
    exit 0  # Success
fi