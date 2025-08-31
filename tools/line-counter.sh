#!/bin/bash
# Line counter for Software Factory 2.0
# Excludes generated code, tests, and documentation
# Properly handles parameters and auto-detects branches

set -euo pipefail

# Initialize variables
BRANCH=""
BASE=""
DETAILED=false
VERBOSE=false
HELP=false

# Function to show usage
show_usage() {
    cat << EOF
Usage: $0 [OPTIONS]

Software Factory 2.0 Line Counter
Counts non-generated lines in git branches, excluding test files and docs.

OPTIONS:
    -c, --current BRANCH    Branch to measure (default: current branch)
    -b, --base BASE        Base branch to compare against (default: auto-detect or main)
    -d, --detailed         Show detailed file breakdown
    -v, --verbose          Show verbose output including what's excluded
    -h, --help            Show this help message

EXAMPLES:
    $0                     # Measure current branch against its base
    $0 -c feature-x        # Measure feature-x branch
    $0 -c feature-x -b main # Measure feature-x against main
    $0 -d                  # Show detailed breakdown

AUTO-DETECTION:
    - If no branch specified, uses current branch
    - If no base specified, tries to detect from commit messages
    - Falls back to 'main' or 'master' if available

EXCLUDED PATTERNS:
    - Generated code: *.pb.go, *_generated.go, zz_generated*, *.gen.go
    - Vendored code: vendor/*
    - Documentation: *.md
    - Test files: *_test.go
    - CRD YAML: *.crd.yaml, *.crd.yml
    - Build artifacts: bin/*, dist/*, build/*
EOF
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -c|--current)
            BRANCH="$2"
            shift 2
            ;;
        -b|--base)
            BASE="$2"
            shift 2
            ;;
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
        *)
            # Handle old positional parameter format for backwards compatibility
            if [ -z "$BRANCH" ]; then
                BRANCH="$1"
            elif [ -z "$BASE" ]; then
                BASE="$1"
            elif [ "$1" = "-d" ]; then
                DETAILED=true
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
        echo "Please specify branch with -c option"
        exit 1
    fi
    [ "$VERBOSE" = true ] && echo "Auto-detected branch: $BRANCH"
fi

# Function to detect base branch from commit messages
detect_base_branch() {
    local branch="$1"
    
    # Try to find base from commit messages
    local base_from_commit=$(git log "$branch" --format=%B -n 20 2>/dev/null | \
        grep -i "from branch:" | head -1 | \
        sed 's/.*from branch:[ ]*//' | \
        awk '{print $1}' | tr -d ' ')
    
    if [ -n "$base_from_commit" ]; then
        # Check if the base branch exists
        if git rev-parse --verify "$base_from_commit" >/dev/null 2>&1; then
            echo "$base_from_commit"
            return 0
        fi
    fi
    
    # Try to find merge-base with common branches
    for candidate in main master develop development; do
        if git rev-parse --verify "$candidate" >/dev/null 2>&1; then
            echo "$candidate"
            return 0
        fi
    done
    
    # Default to main
    echo "main"
}

# Auto-detect base branch if not specified
if [ -z "$BASE" ]; then
    BASE=$(detect_base_branch "$BRANCH")
    [ "$VERBOSE" = true ] && echo "Auto-detected base: $BASE"
fi

# Verify branches exist
if ! git rev-parse --verify "$BRANCH" >/dev/null 2>&1; then
    echo "Error: Branch '$BRANCH' does not exist"
    echo "Available branches:"
    git branch -a | head -10
    exit 1
fi

if ! git rev-parse --verify "$BASE" >/dev/null 2>&1; then
    echo "Warning: Base branch '$BASE' does not exist, trying origin/$BASE"
    if git rev-parse --verify "origin/$BASE" >/dev/null 2>&1; then
        BASE="origin/$BASE"
    else
        echo "Error: Neither '$BASE' nor 'origin/$BASE' exists"
        echo "Available branches:"
        git branch -a | grep -E "(main|master|develop)" | head -5
        exit 1
    fi
fi

# Show what we're measuring
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 Measuring: $BRANCH"
echo "📍 Against base: $BASE"
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
    echo "No changes found between $BASE and $BRANCH"
    echo "Total non-generated lines: 0"
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
echo "✅ Total non-generated lines: $TOTAL_LINES"

# Check if over limit
if [ "$TOTAL_LINES" -gt 800 ]; then
    echo ""
    echo "⚠️  WARNING: Branch exceeds 800 line limit!"
    echo "   Consider splitting into multiple efforts."
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
    exit 2  # Warning exit code
else
    exit 0  # Success
fi