#!/bin/bash

# 🚨🚨🚨 DEPRECATED RULE VALIDATION SCRIPT 🚨🚨🚨
# This script ensures no deprecated rules are referenced in agent states or configs
# Must be run before any commit that modifies agent-states/ or .claude/

set -e

echo "🔍 Validating for deprecated rule references..."
echo "========================================="

# Define deprecated rules that have been consolidated
DEPRECATED_RULES=(
    "R236"  # Consolidated into R290
    "R237"  # Consolidated into R290
    "R252"  # Consolidated into R288
    "R253"  # Consolidated into R288
    "R187"  # Consolidated into R287
    "R188"  # Consolidated into R287
    "R189"  # Consolidated into R287
    "R190"  # Consolidated into R287
)

# Track validation results
VALIDATION_PASSED=true
ERRORS_FOUND=()

# Function to check for deprecated rule references
check_deprecated_references() {
    local dir="$1"
    local dir_name="$2"
    
    echo "📂 Checking $dir_name..."
    
    for rule in "${DEPRECATED_RULES[@]}"; do
        # Search for the rule in non-backup files
        if grep -r "\b${rule}\b" "$dir" --include="*.md" --exclude="*.backup.*" --exclude="RULE-REGISTRY.md" 2>/dev/null | grep -v "DEPRECATED" | grep -v "deprecated" > /dev/null; then
            echo "  ❌ ERROR: Deprecated rule $rule still referenced in $dir_name!"
            
            # Show where it's referenced
            echo "     Found in:"
            grep -r "\b${rule}\b" "$dir" --include="*.md" --exclude="*.backup.*" --exclude="RULE-REGISTRY.md" 2>/dev/null | grep -v "DEPRECATED" | grep -v "deprecated" | while read -r line; do
                echo "       - ${line%%:*}"
            done | sort -u
            
            VALIDATION_PASSED=false
            ERRORS_FOUND+=("$rule in $dir_name")
        fi
    done
}

# Check agent-states directory
if [ -d "agent-states" ]; then
    check_deprecated_references "agent-states" "agent-states"
fi

# Check .claude directory (agent configs and commands)
if [ -d ".claude" ]; then
    check_deprecated_references ".claude" ".claude configs"
fi

# Check orchestrator state files
if [ -f "orchestrator-state.json" ] || [ -f "orchestrator-state.json.example" ]; then
    echo "📂 Checking orchestrator state files..."
    for rule in "${DEPRECATED_RULES[@]}"; do
        if grep "\b${rule}\b" orchestrator-state.json* 2>/dev/null > /dev/null; then
            echo "  ❌ ERROR: Deprecated rule $rule found in orchestrator state file!"
            VALIDATION_PASSED=false
            ERRORS_FOUND+=("$rule in orchestrator-state.json")
        fi
    done
fi

echo "========================================="

# Report results
if $VALIDATION_PASSED; then
    echo "✅ VALIDATION PASSED: No deprecated rule references found!"
    echo ""
    echo "📋 Rules checked as deprecated:"
    for rule in "${DEPRECATED_RULES[@]}"; do
        case "$rule" in
            R236|R237)
                echo "  - $rule → Use R290 instead"
                ;;
            R252|R253)
                echo "  - $rule → Use R288 instead"
                ;;
            R187|R188|R189|R190)
                echo "  - $rule → Use R287 instead"
                ;;
        esac
    done
    exit 0
else
    echo "❌ VALIDATION FAILED: Deprecated rules still referenced!"
    echo ""
    echo "📋 Summary of issues found:"
    for error in "${ERRORS_FOUND[@]}"; do
        echo "  - $error"
    done
    echo ""
    echo "🔧 HOW TO FIX:"
    echo "1. Replace deprecated rule references with their consolidated versions:"
    echo "   - R236, R237 → R290 (State Rule Reading and Verification)"
    echo "   - R252, R253 → R288 (State File Update and Commit Protocol)"
    echo "   - R187, R188, R189, R190 → R287 (Comprehensive TODO Persistence)"
    echo ""
    echo "2. Update the references in the affected files"
    echo "3. Run this validation again to verify"
    echo ""
    echo "⚠️ IMPORTANT: Never reference deprecated rules in new code!"
    exit 1
fi