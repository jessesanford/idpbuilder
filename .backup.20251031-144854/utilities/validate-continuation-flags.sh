#!/bin/bash
# Validation script for R405 Continuation Flag compliance
# Ensures all state files use correct CONTINUE-SOFTWARE-FACTORY=TRUE/FALSE format
# NEVER CONTINUE=TRUE

set -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$( cd "$SCRIPT_DIR/.." && pwd )"

echo "========================================"
echo "R405 Continuation Flag Validator"
echo "========================================"
echo ""

VIOLATIONS_FOUND=0
TOTAL_FILES_CHECKED=0

# Function to check a single file
check_file() {
    local file="$1"
    local violations=()

    ((TOTAL_FILES_CHECKED++))

    # Check for wrong pattern: CONTINUE=TRUE (without SOFTWARE-FACTORY)
    # Exclude lines that already have CONTINUE-SOFTWARE-FACTORY
    # Exclude documentation lines that are just explaining the pattern
    if grep -n "CONTINUE=TRUE" "$file" | \
       grep -v "CONTINUE-SOFTWARE-FACTORY=TRUE" | \
       grep -v "CONTINUE=TRUE/FALSE" | \
       grep -v "# or CONTINUE=TRUE" | \
       grep -v "Example.*CONTINUE=TRUE" | \
       grep -v "if.*CONTINUE.*=.*TRUE" | \
       grep -q .; then

        echo "❌ VIOLATION in $file:"
        echo "   Found CONTINUE=TRUE instead of CONTINUE-SOFTWARE-FACTORY=TRUE"
        echo ""
        echo "   Lines:"
        grep -n "CONTINUE=TRUE" "$file" | \
            grep -v "CONTINUE-SOFTWARE-FACTORY=TRUE" | \
            grep -v "CONTINUE=TRUE/FALSE" | \
            grep -v "# or CONTINUE=TRUE" | \
            grep -v "Example.*CONTINUE=TRUE" | \
            grep -v "if.*CONTINUE.*=.*TRUE"
        echo ""
        ((VIOLATIONS_FOUND++))
        return 1
    fi

    # Check that R405-related sections use the correct flag
    if grep -q "R405" "$file"; then
        # File mentions R405, verify it uses correct flag in examples
        if grep -B5 -A5 "R405" "$file" | \
           grep -E "echo.*CONTINUE.*TRUE" | \
           grep -v "CONTINUE-SOFTWARE-FACTORY=TRUE" | \
           grep -q .; then
            echo "⚠️  WARNING in $file:"
            echo "   R405 section may have incorrect flag format"
            echo "   Please manually review R405 examples"
            echo ""
        fi
    fi

    return 0
}

# Check all orchestrator state files
echo "Checking orchestrator state files..."
echo ""

STATE_FILES=$(find "$PROJECT_ROOT/agent-states/software-factory/orchestrator" -name "rules.md" 2>/dev/null || true)

if [ -z "$STATE_FILES" ]; then
    echo "⚠️  No orchestrator state files found"
else
    for file in $STATE_FILES; do
        check_file "$file"
    done
fi

# Check other agent states
echo "Checking other agent state files..."
echo ""

OTHER_STATE_FILES=$(find "$PROJECT_ROOT/agent-states" -name "rules.md" ! -path "*/orchestrator/*" 2>/dev/null || true)

if [ -n "$OTHER_STATE_FILES" ]; then
    for file in $OTHER_STATE_FILES; do
        check_file "$file"
    done
fi

# Check agent configs
echo "Checking agent configuration files..."
echo ""

AGENT_CONFIGS=$(find "$PROJECT_ROOT/.claude/agents" -name "*.md" 2>/dev/null || true)

if [ -n "$AGENT_CONFIGS" ]; then
    for file in $AGENT_CONFIGS; do
        check_file "$file"
    done
fi

# Report summary
echo "========================================"
echo "VALIDATION SUMMARY"
echo "========================================"
echo "Files checked: $TOTAL_FILES_CHECKED"
echo "Violations found: $VIOLATIONS_FOUND"
echo ""

if [ $VIOLATIONS_FOUND -eq 0 ]; then
    echo "✅ ALL FILES PASS R405 VALIDATION"
    echo ""
    echo "All files correctly use:"
    echo "  - CONTINUE-SOFTWARE-FACTORY=TRUE"
    echo "  - CONTINUE-SOFTWARE-FACTORY=FALSE"
    echo ""
    exit 0
else
    echo "❌ VALIDATION FAILED"
    echo ""
    echo "Found $VIOLATIONS_FOUND file(s) with incorrect continuation flag format"
    echo ""
    echo "REQUIRED FORMAT:"
    echo "  ✅ CONTINUE-SOFTWARE-FACTORY=TRUE"
    echo "  ✅ CONTINUE-SOFTWARE-FACTORY=FALSE"
    echo ""
    echo "WRONG FORMAT:"
    echo "  ❌ CONTINUE=TRUE"
    echo "  ❌ CONTINUE=FALSE"
    echo ""
    echo "See: $PROJECT_ROOT/rule-library/R405-automation-continuation-flag.md"
    echo ""
    exit 1
fi
