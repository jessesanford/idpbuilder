#!/bin/bash

# Orchestrator State Rules Cleanup Script
# This script will process all orchestrator state rule files to:
# 1. Remove inline rule duplications (content between --- delimiters)
# 2. Ensure proper rule references to rule-library
# 3. Standardize the structure

ORCHESTRATOR_DIR="/home/vscode/software-factory-template/agent-states/orchestrator"
RULE_LIBRARY="/home/vscode/software-factory-template/rule-library"
REPORT_FILE="/home/vscode/software-factory-template/ORCHESTRATOR-RULES-CLEANUP-REPORT.md"

echo "# Orchestrator State Rules Cleanup Report" > "$REPORT_FILE"
echo "Date: $(date '+%Y-%m-%d %H:%M:%S')" >> "$REPORT_FILE"
echo "" >> "$REPORT_FILE"
echo "## Summary" >> "$REPORT_FILE"
echo "" >> "$REPORT_FILE"

# Track metrics
TOTAL_FILES=0
FILES_MODIFIED=0
LINES_REMOVED=0
RULES_TO_CREATE=""

echo "🔍 Starting orchestrator state rules cleanup..."

# Process each state directory
for STATE_DIR in "$ORCHESTRATOR_DIR"/*; do
    if [ -d "$STATE_DIR" ]; then
        STATE_NAME=$(basename "$STATE_DIR")
        RULES_FILE="$STATE_DIR/rules.md"
        
        if [ -f "$RULES_FILE" ]; then
            echo "Processing: $STATE_NAME"
            TOTAL_FILES=$((TOTAL_FILES + 1))
            
            # Count lines with --- delimiters (indicates inline rules)
            DELIMITER_COUNT=$(grep -c "^---$" "$RULES_FILE" 2>/dev/null || echo 0)
            
            if [ "$DELIMITER_COUNT" -gt 0 ]; then
                echo "  ⚠️ Found $DELIMITER_COUNT inline rule blocks in $STATE_NAME"
                FILES_MODIFIED=$((FILES_MODIFIED + 1))
                
                # Extract rules referenced in inline blocks
                INLINE_RULES=$(sed -n '/^---$/,/^---$/p' "$RULES_FILE" | grep -o "RULE R[0-9]\+" | sed 's/RULE //' | sort -u)
                
                if [ -n "$INLINE_RULES" ]; then
                    echo "  📋 Inline rules found: $INLINE_RULES"
                    for RULE in $INLINE_RULES; do
                        # Check if rule exists in library
                        if ! find "$RULE_LIBRARY" -name "${RULE}*.md" | grep -q .; then
                            echo "    ❌ Rule $RULE not found in rule library"
                            RULES_TO_CREATE="$RULES_TO_CREATE $RULE"
                        fi
                    done
                fi
                
                # Count lines to be removed
                LINES_IN_BLOCKS=$(sed -n '/^---$/,/^---$/p' "$RULES_FILE" | wc -l)
                LINES_REMOVED=$((LINES_REMOVED + LINES_IN_BLOCKS))
                
                echo "### $STATE_NAME" >> "$REPORT_FILE"
                echo "- Inline rule blocks: $DELIMITER_COUNT" >> "$REPORT_FILE"
                echo "- Lines to remove: $LINES_IN_BLOCKS" >> "$REPORT_FILE"
                echo "- Rules referenced: $INLINE_RULES" >> "$REPORT_FILE"
                echo "" >> "$REPORT_FILE"
            fi
        fi
    fi
done

echo "" >> "$REPORT_FILE"
echo "## Metrics" >> "$REPORT_FILE"
echo "- Total files scanned: $TOTAL_FILES" >> "$REPORT_FILE"
echo "- Files with duplications: $FILES_MODIFIED" >> "$REPORT_FILE"
echo "- Total lines to remove: $LINES_REMOVED" >> "$REPORT_FILE"
echo "" >> "$REPORT_FILE"

if [ -n "$RULES_TO_CREATE" ]; then
    echo "## Rules Missing from Library" >> "$REPORT_FILE"
    echo "The following rules are referenced but don't exist in rule-library:" >> "$REPORT_FILE"
    for RULE in $(echo "$RULES_TO_CREATE" | tr ' ' '\n' | sort -u); do
        echo "- $RULE" >> "$REPORT_FILE"
    done
    echo "" >> "$REPORT_FILE"
fi

echo "## Files to Process" >> "$REPORT_FILE"
echo "" >> "$REPORT_FILE"

# List files that need processing
for STATE_DIR in "$ORCHESTRATOR_DIR"/*; do
    if [ -d "$STATE_DIR" ]; then
        STATE_NAME=$(basename "$STATE_DIR")
        RULES_FILE="$STATE_DIR/rules.md"
        
        if [ -f "$RULES_FILE" ]; then
            DELIMITER_COUNT=$(grep -c "^---$" "$RULES_FILE" 2>/dev/null || echo 0)
            if [ "$DELIMITER_COUNT" -gt 0 ]; then
                echo "- $STATE_NAME/rules.md ($DELIMITER_COUNT blocks)" >> "$REPORT_FILE"
            fi
        fi
    fi
done

echo "" >> "$REPORT_FILE"
echo "---" >> "$REPORT_FILE"
echo "Report generated at: $(date '+%Y-%m-%d %H:%M:%S')" >> "$REPORT_FILE"

echo "✅ Analysis complete. Report saved to: $REPORT_FILE"