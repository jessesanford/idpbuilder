#!/bin/bash

# Update state verification markers to new .software-factory/ structure
# More conservative approach - only update the marker creation/checking patterns

set -e

echo "📦 Updating state verification markers to new structure..."
echo ""

# Counter for changes
CHANGES_MADE=0

# Function to update a single file
update_file() {
    local file="$1"
    local agent_name=$(basename $(dirname $(dirname "$file")))
    local state_name=$(basename $(dirname "$file"))
    
    # Skip backup files
    if [[ "$file" == *".backup"* ]]; then
        return
    fi
    
    echo "  Checking: $file"
    
    # Create temp file for changes
    local temp_file="${file}.temp"
    cp "$file" "$temp_file"
    
    # Pattern 1: Update touch commands for state markers
    # Old: touch .state_rules_read_orchestrator_MONITOR
    # New: mkdir -p .software-factory/markers && touch ".software-factory/markers/state_rules_read_orchestrator_MONITOR-$(date +%Y%m%d-%H%M%S)"
    
    if grep -q "touch \.state_rules_read_" "$temp_file"; then
        sed -i 's|^touch \.state_rules_read_\(.*\)$|mkdir -p .software-factory/markers \&\& touch ".software-factory/markers/state_rules_read_\1-$(date +%Y%m%d-%H%M%S)"|g' "$temp_file"
        echo "    ✅ Updated touch commands"
        ((CHANGES_MADE++))
    fi
    
    # Pattern 2: Update echo commands that write to state markers
    # Old: echo "$(date +%s) - Rules read..." > .state_rules_read_orchestrator_MONITOR
    # New: mkdir -p .software-factory/markers && echo "$(date +%s) - Rules read..." > ".software-factory/markers/state_rules_read_orchestrator_MONITOR-$(date +%Y%m%d-%H%M%S)"
    
    if grep -q 'echo.*> \.state_rules_read_' "$temp_file"; then
        sed -i 's|> \.state_rules_read_\(.*\)$|> ".software-factory/markers/state_rules_read_\1-$(date +%Y%m%d-%H%M%S)"|g' "$temp_file"
        
        # Also add mkdir before echo commands if not already present
        sed -i '/echo.*> "\.software-factory\/markers\/state_rules_read/i mkdir -p .software-factory/markers' "$temp_file"
        
        echo "    ✅ Updated echo commands"
        ((CHANGES_MADE++))
    fi
    
    # Pattern 3: Update references in backticks and parentheses (documentation)
    if grep -q '`\.state_rules_read_' "$temp_file"; then
        sed -i 's|`\.state_rules_read_\([^`]*\)`|`.software-factory/markers/state_rules_read_\1-TIMESTAMP`|g' "$temp_file"
        echo "    ✅ Updated documentation references"
        ((CHANGES_MADE++))
    fi
    
    if grep -q '(\.state_rules_read_' "$temp_file"; then
        sed -i 's|(\.state_rules_read_\([^)]*\))|(.software-factory/markers/state_rules_read_\1-TIMESTAMP)|g' "$temp_file"
        echo "    ✅ Updated parenthetical references"
        ((CHANGES_MADE++))
    fi
    
    # Only replace file if changes were made
    if ! diff -q "$file" "$temp_file" >/dev/null 2>&1; then
        mv "$temp_file" "$file"
    else
        rm "$temp_file"
    fi
}

# Process all state rule files
echo "🔍 Processing state rule files..."
find /home/vscode/software-factory-template/agent-states -name "rules.md" -type f | while read -r file; do
    update_file "$file"
done

echo ""
echo "🔍 Processing agent configuration files..."
for file in /home/vscode/software-factory-template/.claude/agents/*.md; do
    if [ -f "$file" ]; then
        update_file "$file"
    fi
done

echo ""
echo "✅ State verification marker updates complete!"
echo "📊 Total files with changes: $CHANGES_MADE"
echo ""
echo "📋 New structure:"
echo "  .software-factory/"
echo "  └── markers/"
echo "      ├── state_rules_read_orchestrator_[STATE]-[TIMESTAMP]"
echo "      ├── state_rules_read_sw-engineer_[STATE]-[TIMESTAMP]"
echo "      ├── state_rules_read_code-reviewer_[STATE]-[TIMESTAMP]"
echo "      └── state_rules_read_architect_[STATE]-[TIMESTAMP]"