#!/bin/bash

# Validate all rule references in state files are correct
# This script checks that every rule reference points to an existing file

echo "🔍 RULE REFERENCE VALIDATION SCRIPT"
echo "===================================="
echo ""

cd /home/vscode/software-factory-template

# Initialize counters
TOTAL_REFS=0
VALID_REFS=0
INVALID_REFS=0
INVALID_FILES=()

echo "📊 Scanning all state files for rule references..."
echo "---------------------------------------------------"

# Get all .md files in agent-states directory (excluding backups)
STATE_FILES=$(find agent-states -name "*.md" -type f | grep -v ".backup" | sort)

for state_file in $STATE_FILES; do
    # Extract all rule references from this file
    refs=$(grep -ho "R[0-9]\{3\}[^[:space:]]*\.md" "$state_file" 2>/dev/null | sort -u)

    if [ -n "$refs" ]; then
        for ref in $refs; do
            TOTAL_REFS=$((TOTAL_REFS + 1))

            # Check if the referenced rule file exists
            if [ -f "rule-library/$ref" ]; then
                VALID_REFS=$((VALID_REFS + 1))
            else
                INVALID_REFS=$((INVALID_REFS + 1))
                echo "❌ Invalid reference in $state_file: $ref"
                INVALID_FILES+=("$state_file:$ref")
            fi
        done
    fi
done

echo ""
echo "📊 VALIDATION SUMMARY"
echo "===================="
echo "Total rule references checked: $TOTAL_REFS"
echo "Valid references: $VALID_REFS ✅"
echo "Invalid references: $INVALID_REFS ❌"
echo ""

if [ $INVALID_REFS -eq 0 ]; then
    echo "🎉 PROJECT_DONE: All rule references are valid!"
    exit 0
else
    echo "⚠️ WARNING: Found $INVALID_REFS invalid references"
    echo ""
    echo "Invalid references by file:"
    printf '%s\n' "${INVALID_FILES[@]}" | head -20

    if [ ${#INVALID_FILES[@]} -gt 20 ]; then
        echo "... and $((${#INVALID_FILES[@]} - 20)) more"
    fi

    exit 1
fi