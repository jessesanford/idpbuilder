#!/bin/bash
#
# Software Factory Line Counter
# 
# Purpose: Accurately counts hand-written implementation lines in feature branches
# Usage: ./line-counter.sh -c <branch-name> [-b <base-branch>] [-d]
# 
# This script properly excludes generated files while counting
# hand-written implementation code only.
# 
# Configure the patterns below for your specific project.
#

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Configuration - CUSTOMIZE FOR YOUR PROJECT
MAX_LINES_WARNING=700
MAX_LINES_ERROR=800

# Function to print colored output
print_colored() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# Function to print section headers
print_header() {
    echo
    print_colored "$BLUE" "=========================================="
    print_colored "$BLUE" "$1"
    print_colored "$BLUE" "=========================================="
}

# Function to check if a file is generated
# CUSTOMIZE: Add your project's generated file patterns
is_generated_file() {
    local file=$1
    
    # Common generated file patterns - CUSTOMIZE FOR YOUR PROJECT
    
    # Go generated files
    if [[ "$file" =~ zz_generated ]]; then
        return 0
    fi
    if [[ "$file" =~ _generated\.go$ ]]; then
        return 0
    fi
    if [[ "$file" =~ \.pb\.go$ ]]; then  # Protocol buffers
        return 0
    fi
    
    # Python generated files
    if [[ "$file" =~ __pycache__ ]]; then
        return 0
    fi
    if [[ "$file" =~ \.pyc$ ]]; then
        return 0
    fi
    
    # JavaScript/TypeScript generated files
    if [[ "$file" =~ /dist/ ]]; then
        return 0
    fi
    if [[ "$file" =~ /build/ ]]; then
        return 0
    fi
    if [[ "$file" =~ \.min\.js$ ]]; then
        return 0
    fi
    if [[ "$file" =~ \.d\.ts$ ]]; then
        return 0
    fi
    
    # Common build artifacts
    if [[ "$file" =~ /target/ ]]; then  # Maven/Rust
        return 0
    fi
    if [[ "$file" =~ /out/ ]]; then
        return 0
    fi
    
    # Vendor/dependencies
    if [[ "$file" =~ /vendor/ ]]; then
        return 0
    fi
    if [[ "$file" =~ /node_modules/ ]]; then
        return 0
    fi
    if [[ "$file" =~ /.venv/ ]]; then
        return 0
    fi
    
    # Test data/fixtures
    if [[ "$file" =~ /testdata/ ]]; then
        return 0
    fi
    if [[ "$file" =~ /fixtures/ ]]; then
        return 0
    fi
    if [[ "$file" =~ /mocks/ ]]; then
        return 0
    fi
    
    return 1
}

# Function to check if a file is documentation
is_doc_file() {
    local file=$1
    
    if [[ "$file" =~ \.md$ ]]; then
        return 0
    fi
    if [[ "$file" =~ \.rst$ ]]; then
        return 0
    fi
    if [[ "$file" =~ \.txt$ ]]; then
        return 0
    fi
    if [[ "$file" =~ \.adoc$ ]]; then
        return 0
    fi
    if [[ "$file" =~ /docs/ ]]; then
        return 0
    fi
    if [[ "$file" =~ README ]]; then
        return 0
    fi
    if [[ "$file" =~ LICENSE ]]; then
        return 0
    fi
    
    return 1
}

# Function to check if a file is configuration
is_config_file() {
    local file=$1
    
    if [[ "$file" =~ \.yaml$ ]] || [[ "$file" =~ \.yml$ ]]; then
        # Allow application code YAML but not config
        if [[ "$file" =~ config/ ]] || [[ "$file" =~ \.config\. ]]; then
            return 0
        fi
    fi
    
    if [[ "$file" =~ \.json$ ]]; then
        if [[ "$file" =~ package\.json$ ]] || [[ "$file" =~ tsconfig\.json$ ]]; then
            return 0
        fi
    fi
    
    if [[ "$file" =~ \.toml$ ]]; then
        return 0
    fi
    
    if [[ "$file" =~ \.ini$ ]]; then
        return 0
    fi
    
    if [[ "$file" =~ Makefile ]]; then
        return 0
    fi
    
    if [[ "$file" =~ Dockerfile ]]; then
        return 0
    fi
    
    if [[ "$file" =~ \.gitignore$ ]]; then
        return 0
    fi
    
    return 1
}

# Parse command line arguments
COMPARE_BRANCH=""
BASE_BRANCH="main"
DETAILED=false

while getopts "c:b:dh" opt; do
    case ${opt} in
        c )
            COMPARE_BRANCH=$OPTARG
            ;;
        b )
            BASE_BRANCH=$OPTARG
            ;;
        d )
            DETAILED=true
            ;;
        h )
            echo "Usage: $0 -c <branch-name> [-b <base-branch>] [-d]"
            echo "  -c: Branch to compare (required)"
            echo "  -b: Base branch (default: main)"
            echo "  -d: Show detailed breakdown by file"
            echo "  -h: Show this help message"
            exit 0
            ;;
        \? )
            echo "Invalid option: $OPTARG" 1>&2
            exit 1
            ;;
    esac
done

if [ -z "$COMPARE_BRANCH" ]; then
    print_colored "$RED" "Error: Branch name required. Use -c <branch-name>"
    exit 1
fi

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    print_colored "$RED" "Error: Not in a git repository"
    exit 1
fi

# Check if branches exist
if ! git rev-parse "$BASE_BRANCH" > /dev/null 2>&1; then
    print_colored "$RED" "Error: Base branch '$BASE_BRANCH' does not exist"
    exit 1
fi

if ! git rev-parse "$COMPARE_BRANCH" > /dev/null 2>&1; then
    print_colored "$RED" "Error: Compare branch '$COMPARE_BRANCH' does not exist"
    exit 1
fi

print_header "Software Factory Line Counter"
echo "Comparing: $COMPARE_BRANCH"
echo "Against:   $BASE_BRANCH"
echo

# Get the list of changed files
CHANGED_FILES=$(git diff --name-only "$BASE_BRANCH...$COMPARE_BRANCH" | sort)

if [ -z "$CHANGED_FILES" ]; then
    print_colored "$YELLOW" "No files changed between $BASE_BRANCH and $COMPARE_BRANCH"
    exit 0
fi

# Process files
TOTAL_LINES=0
IMPL_LINES=0
GENERATED_LINES=0
DOC_LINES=0
CONFIG_LINES=0
TEST_LINES=0

declare -a IMPL_FILES
declare -a IMPL_LINES_PER_FILE

print_header "Analyzing Files"

while IFS= read -r file; do
    # Skip if file was deleted
    if ! git show "$COMPARE_BRANCH:$file" > /dev/null 2>&1; then
        continue
    fi
    
    # Count lines in the file
    if [ -f "$file" ]; then
        lines=$(wc -l < "$file" 2>/dev/null || echo 0)
    else
        lines=$(git show "$COMPARE_BRANCH:$file" 2>/dev/null | wc -l || echo 0)
    fi
    
    # Categorize the file
    if is_generated_file "$file"; then
        GENERATED_LINES=$((GENERATED_LINES + lines))
        if [ "$DETAILED" = true ]; then
            print_colored "$PURPLE" "  [GENERATED] $file: $lines lines"
        fi
    elif is_doc_file "$file"; then
        DOC_LINES=$((DOC_LINES + lines))
        if [ "$DETAILED" = true ]; then
            print_colored "$CYAN" "  [DOC] $file: $lines lines"
        fi
    elif is_config_file "$file"; then
        CONFIG_LINES=$((CONFIG_LINES + lines))
        if [ "$DETAILED" = true ]; then
            print_colored "$YELLOW" "  [CONFIG] $file: $lines lines"
        fi
    elif [[ "$file" =~ test ]] || [[ "$file" =~ spec ]]; then
        TEST_LINES=$((TEST_LINES + lines))
        IMPL_LINES=$((IMPL_LINES + lines))  # Tests count toward implementation
        IMPL_FILES+=("$file")
        IMPL_LINES_PER_FILE+=("$lines")
        if [ "$DETAILED" = true ]; then
            print_colored "$GREEN" "  [TEST] $file: $lines lines"
        fi
    else
        IMPL_LINES=$((IMPL_LINES + lines))
        IMPL_FILES+=("$file")
        IMPL_LINES_PER_FILE+=("$lines")
        if [ "$DETAILED" = true ]; then
            print_colored "$GREEN" "  [IMPL] $file: $lines lines"
        fi
    fi
    
    TOTAL_LINES=$((TOTAL_LINES + lines))
done <<< "$CHANGED_FILES"

print_header "Summary"

# Print summary with color coding
echo "Implementation: $([ $IMPL_LINES -le $MAX_LINES_WARNING ] && print_colored "$GREEN" "$IMPL_LINES lines" || [ $IMPL_LINES -le $MAX_LINES_ERROR ] && print_colored "$YELLOW" "$IMPL_LINES lines" || print_colored "$RED" "$IMPL_LINES lines")"
echo "├── Source:     $((IMPL_LINES - TEST_LINES)) lines"
echo "└── Tests:      $TEST_LINES lines"
echo
echo "Excluded:"
echo "├── Generated:  $GENERATED_LINES lines"
echo "├── Docs:       $DOC_LINES lines"
echo "└── Config:     $CONFIG_LINES lines"
echo
echo "Total Changed:  $TOTAL_LINES lines"

# Status indicator
echo
if [ $IMPL_LINES -le $MAX_LINES_WARNING ]; then
    print_colored "$GREEN" "✓ SIZE OK: Implementation is under $MAX_LINES_WARNING lines"
    exit 0
elif [ $IMPL_LINES -le $MAX_LINES_ERROR ]; then
    print_colored "$YELLOW" "⚠ WARNING: Implementation approaching limit ($IMPL_LINES/$MAX_LINES_ERROR lines)"
    exit 0
else
    print_colored "$RED" "✗ ERROR: Implementation exceeds $MAX_LINES_ERROR line limit!"
    print_colored "$RED" "  Action Required: Split this effort into smaller parts"
    
    if [ "$DETAILED" = true ]; then
        echo
        print_colored "$RED" "Files contributing to implementation:"
        for i in "${!IMPL_FILES[@]}"; do
            print_colored "$RED" "  ${IMPL_FILES[$i]}: ${IMPL_LINES_PER_FILE[$i]} lines"
        done
    fi
    
    exit 1
fi