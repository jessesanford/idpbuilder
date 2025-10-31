#!/bin/bash

# Comprehensive orchestrator state rule reference audit script
# Purpose: Verify all rule references in orchestrator state files have referential integrity

STATES_DIR="/home/vscode/software-factory-template/agent-states/orchestrator"
RULE_LIBRARY="/home/vscode/software-factory-template/rule-library"
REPORT_FILE="/home/vscode/software-factory-template/orchestrator-audit-report.md"

# Initialize counters
TOTAL_STATES=0
TOTAL_REFERENCES=0
MISMATCHES_FOUND=0
FIXES_MADE=0

# Start report
cat > "$REPORT_FILE" << 'EOF'
# ORCHESTRATOR STATE RULE REFERENCE AUDIT REPORT
Generated: $(date '+%Y-%m-%d %H:%M:%S %Z')

## Summary Statistics
EOF

echo "🔍 Starting comprehensive orchestrator state rule audit..."
echo ""

# Track issues for later fixing
declare -A ISSUES

# Function to extract rule references from a file
extract_rule_references() {
    local file="$1"
    # Look for patterns like R###, focusing on rule references
    grep -oE 'R[0-9]{3}[^0-9]' "$file" 2>/dev/null | sed 's/[^R0-9]//g' | sort -u
}

# Function to check if a rule file exists
check_rule_exists() {
    local rule_num="$1"
    local expected_path="$RULE_LIBRARY/${rule_num}-*.md"
    
    # Use ls to check for the file (handles wildcards)
    if ls $expected_path >/dev/null 2>&1; then
        echo "$(ls $expected_path | head -1)"
        return 0
    else
        return 1
    fi
}

# Audit each state
for state_dir in "$STATES_DIR"/*/; do
    state_name=$(basename "$state_dir")
    
    # Skip non-directories
    [ ! -d "$state_dir" ] && continue
    
    rules_file="$state_dir/rules.md"
    
    # Check if rules.md exists
    if [ ! -f "$rules_file" ]; then
        echo "⚠️  State $state_name has no rules.md file"
        continue
    fi
    
    TOTAL_STATES=$((TOTAL_STATES + 1))
    echo "📂 Checking state: $state_name"
    
    # Extract rule references
    references=$(extract_rule_references "$rules_file")
    
    if [ -z "$references" ]; then
        echo "   No rule references found"
        continue
    fi
    
    # Check each reference
    for rule_ref in $references; do
        TOTAL_REFERENCES=$((TOTAL_REFERENCES + 1))
        
        if actual_path=$(check_rule_exists "$rule_ref"); then
            echo "   ✅ $rule_ref -> $(basename "$actual_path")"
        else
            echo "   ❌ $rule_ref -> FILE NOT FOUND!"
            MISMATCHES_FOUND=$((MISMATCHES_FOUND + 1))
            ISSUES["$state_name:$rule_ref"]="missing"
        fi
    done
    echo ""
done

# Complete the report
cat >> "$REPORT_FILE" << EOF

- **Total States Checked:** $TOTAL_STATES
- **Total Rule References:** $TOTAL_REFERENCES
- **Mismatches Found:** $MISMATCHES_FOUND
- **Fixes Applied:** $FIXES_MADE

## Detailed Findings

EOF

# Add detailed findings
if [ $MISMATCHES_FOUND -gt 0 ]; then
    echo "### Issues Found:" >> "$REPORT_FILE"
    for issue in "${!ISSUES[@]}"; do
        echo "- $issue: ${ISSUES[$issue]}" >> "$REPORT_FILE"
    done
else
    echo "✅ **All rule references have valid referential integrity!**" >> "$REPORT_FILE"
fi

echo ""
echo "═══════════════════════════════════════════════════"
echo "📊 AUDIT SUMMARY:"
echo "   Total States Checked: $TOTAL_STATES"
echo "   Total Rule References: $TOTAL_REFERENCES"
echo "   Mismatches Found: $MISMATCHES_FOUND"
echo "═══════════════════════════════════════════════════"
echo ""
echo "📝 Full report saved to: $REPORT_FILE"