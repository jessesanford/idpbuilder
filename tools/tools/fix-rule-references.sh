#!/bin/bash

# Fix incorrect rule references in state files
# This script will identify and fix all mismatched rule references

echo "🔍 Starting rule reference validation and fix process..."
echo "=================================================="

# Define the working directory
WORK_DIR="/home/vscode/software-factory-template"
cd "$WORK_DIR"

# Track fixes
FIX_COUNT=0
FILES_FIXED=()

# Function to check and fix a specific rule reference
fix_rule_reference() {
    local wrong_ref="$1"
    local correct_ref="$2"
    local description="$3"

    echo -n "Checking for incorrect reference: $wrong_ref -> $correct_ref ... "

    # Find files with the incorrect reference
    files_with_error=$(grep -rl "$wrong_ref" agent-states/ 2>/dev/null || true)

    if [ -z "$files_with_error" ]; then
        echo "✅ No incorrect references found"
        return
    fi

    echo "❌ Found incorrect references"

    for file in $files_with_error; do
        echo "  Fixing: $file"
        sed -i "s|$wrong_ref|$correct_ref|g" "$file"
        FIX_COUNT=$((FIX_COUNT + 1))
        FILES_FIXED+=("$file")
    done
}

# Known incorrect references to fix
echo ""
echo "📝 Fixing known incorrect rule references..."
echo "--------------------------------------------"

# R288 - state file update protocol
fix_rule_reference \
    "R288-state-file-update-requirements.md" \
    "R288-state-file-update-and-commit-protocol.md" \
    "R288 state file update protocol"

# R322 - mandatory stop before transitions
fix_rule_reference \
    "R322-mandatory-stop-before-transitions.md" \
    "R322-mandatory-stop-before-state-transitions.md" \
    "R322 mandatory stop before state transitions"

# Now let's check for any other potentially incorrect references
echo ""
echo "🔍 Scanning for all rule references in state files..."
echo "-----------------------------------------------------"

# Extract all rule references from state files
echo "Extracting all rule references..."
all_rule_refs=$(grep -rho "R[0-9]\{3\}[^[:space:]]*\.md" agent-states/ 2>/dev/null | sort -u || true)

echo ""
echo "📊 Validating all rule references against rule library..."
echo "--------------------------------------------------------"

INVALID_REFS=()

for ref in $all_rule_refs; do
    # Check if this rule file exists in rule-library
    if [ ! -f "rule-library/$ref" ]; then
        echo "❌ Invalid reference found: $ref"
        INVALID_REFS+=("$ref")

        # Try to find the correct filename
        rule_number=$(echo "$ref" | grep -o "R[0-9]\{3\}")
        correct_file=$(ls rule-library/${rule_number}*.md 2>/dev/null | head -1 || true)

        if [ -n "$correct_file" ]; then
            correct_ref=$(basename "$correct_file")
            echo "  ✅ Found correct file: $correct_ref"

            # Fix this reference
            fix_rule_reference "$ref" "$correct_ref" "Auto-detected correction for $rule_number"
        else
            echo "  ⚠️  No matching file found for $rule_number"
        fi
    fi
done

# Now verify all fixes
echo ""
echo "🔍 Verifying all fixes..."
echo "------------------------"

# Re-scan for invalid references
remaining_invalid=$(grep -rho "R[0-9]\{3\}[^[:space:]]*\.md" agent-states/ 2>/dev/null | while read ref; do
    if [ ! -f "rule-library/$ref" ]; then
        echo "$ref"
    fi
done | sort -u)

if [ -z "$remaining_invalid" ]; then
    echo "✅ All rule references are now valid!"
else
    echo "⚠️ Some invalid references remain:"
    echo "$remaining_invalid"
fi

# Summary
echo ""
echo "📊 Summary"
echo "========="
echo "Total fixes applied: $FIX_COUNT"
echo "Files modified: ${#FILES_FIXED[@]}"

if [ ${#FILES_FIXED[@]} -gt 0 ]; then
    echo ""
    echo "Modified files:"
    printf '%s\n' "${FILES_FIXED[@]}" | sort -u | head -20

    if [ ${#FILES_FIXED[@]} -gt 20 ]; then
        echo "... and $((${#FILES_FIXED[@]} - 20)) more files"
    fi
fi

echo ""
echo "✅ Script completed!"