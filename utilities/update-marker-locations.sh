#!/bin/bash

# Script to update all marker file locations to new .software-factory/ structure

set -e

echo "📦 Updating marker file locations to new .software-factory/ structure..."

# Function to update state verification markers in a file
update_state_markers() {
    local file="$1"
    local agent_name=$(basename $(dirname $(dirname "$file")))
    local state_name=$(basename $(dirname "$file"))
    
    echo "  📝 Processing: $file"
    
    # Create backup
    cp "$file" "$file.backup.marker-update" 2>/dev/null || true
    
    # Update the old marker pattern to new structure
    # Old: touch .state_rules_read_orchestrator_MONITOR
    # New: mkdir -p .software-factory/markers && touch .software-factory/markers/state_rules_read_orchestrator_MONITOR-$(date +%Y%m%d-%H%M%S)
    
    # Replace simple touch commands
    sed -i 's|touch \.state_rules_read_\([^"]*\)$|mkdir -p .software-factory/markers \&\& touch .software-factory/markers/state_rules_read_\1-$(date +%Y%m%d-%H%M%S)|g' "$file"
    
    # Replace echo commands with date into marker files
    sed -i 's|echo "\$(date +%s)[^"]*" > \.state_rules_read_\([^"]*\)$|mkdir -p .software-factory/markers \&\& echo "$(date +%s) - Rules read and acknowledged for \1" > .software-factory/markers/state_rules_read_\1-$(date +%Y%m%d-%H%M%S)|g' "$file"
    
    # Update references to marker files in verification checks
    sed -i 's|\[[ -f \.state_rules_read_\([^]]*\) \]|\[ -n "$(ls .software-factory/markers/state_rules_read_\1-* 2>/dev/null)" \]|g' "$file"
    
    # Update references in text/documentation
    sed -i 's|`\.state_rules_read_\([^`]*\)`|`.software-factory/markers/state_rules_read_\1-TIMESTAMP`|g' "$file"
    sed -i 's|(\.state_rules_read_\([^)]*\))|(.software-factory/markers/state_rules_read_\1-TIMESTAMP)|g' "$file"
    
    # Check if file was modified
    if ! diff -q "$file" "$file.backup.marker-update" >/dev/null 2>&1; then
        echo "    ✅ Updated marker references"
        rm "$file.backup.marker-update"
    else
        echo "    ⏭️  No marker references found"
        rm "$file.backup.marker-update"
    fi
}

# Function to update plan/report file locations
update_plan_report_markers() {
    local file="$1"
    
    echo "  📝 Processing plan/report references: $file"
    
    # Create backup
    cp "$file" "$file.backup.marker-update" 2>/dev/null || true
    
    # Update IMPLEMENTATION-PLAN.md references
    sed -i 's|IMPLEMENTATION-PLAN\.md|.software-factory/phase\${PHASE}/wave\${WAVE}/\${EFFORT_NAME}/plans/IMPLEMENTATION-PLAN-$(date +%Y%m%d-%H%M%S).md|g' "$file"
    
    # Update CODE-REVIEW-REPORT.md references  
    sed -i 's|CODE-REVIEW-REPORT\.md|.software-factory/phase\${PHASE}/wave\${WAVE}/\${EFFORT_NAME}/reports/CODE-REVIEW-REPORT-$(date +%Y%m%d-%H%M%S).md|g' "$file"
    
    # Update SPLIT-PLAN.md references
    sed -i 's|SPLIT-PLAN\.md|.software-factory/phase\${PHASE}/wave\${WAVE}/\${EFFORT_NAME}/plans/SPLIT-PLAN-$(date +%Y%m%d-%H%M%S).md|g' "$file"
    
    # Update FIX-PLAN.md references
    sed -i 's|FIX-PLAN\.md|.software-factory/phase\${PHASE}/wave\${WAVE}/\${EFFORT_NAME}/plans/FIX-PLAN-$(date +%Y%m%d-%H%M%S).md|g' "$file"
    
    # Check if file was modified
    if ! diff -q "$file" "$file.backup.marker-update" >/dev/null 2>&1; then
        echo "    ✅ Updated plan/report references"
        rm "$file.backup.marker-update"
    else
        echo "    ⏭️  No plan/report references found"
        rm "$file.backup.marker-update"
    fi
}

# Process all state rule files
echo "🔍 Finding all state rule files..."
STATE_FILES=$(find /home/vscode/software-factory-template/agent-states -name "rules.md" -type f | grep -v backup)

echo "📊 Found $(echo "$STATE_FILES" | wc -l) state rule files to process"

for file in $STATE_FILES; do
    update_state_markers "$file"
    update_plan_report_markers "$file"
done

# Also update agent configuration files
echo ""
echo "🔍 Updating agent configuration files..."
AGENT_FILES="/home/vscode/software-factory-template/.claude/agents/*.md"

for file in $AGENT_FILES; do
    if [ -f "$file" ]; then
        update_state_markers "$file"
        update_plan_report_markers "$file"
    fi
done

# Update utility scripts that check for markers
echo ""
echo "🔍 Updating utility scripts..."
UTIL_FILES=$(find /home/vscode/software-factory-template/utilities -name "*.sh" -type f)

for file in $UTIL_FILES; do
    update_state_markers "$file"
done

# Update rule library files that might reference markers
echo ""
echo "🔍 Updating rule library files..."
RULE_FILES=$(find /home/vscode/software-factory-template/rule-library -name "*.md" -type f | grep -v backup)

for file in $RULE_FILES; do
    update_plan_report_markers "$file"
done

echo ""
echo "✅ Marker location updates complete!"
echo ""
echo "📋 Summary of changes:"
echo "  - State verification markers: .state_rules_read_* → .software-factory/markers/state_rules_read_*-TIMESTAMP"
echo "  - Plan files: [NAME].md → .software-factory/phase*/wave*/{effort}/plans/[NAME]-TIMESTAMP.md"
echo "  - Report files: [NAME].md → .software-factory/phase*/wave*/{effort}/reports/[NAME]-TIMESTAMP.md"
echo ""
echo "⚠️  Note: Review the changes and test before committing!"