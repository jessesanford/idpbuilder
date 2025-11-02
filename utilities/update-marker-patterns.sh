#!/bin/bash
# update-marker-patterns.sh - Update marker file patterns in all state rule files
# Part of marker file reorganization (R290 update)

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$PROJECT_ROOT"

echo "🔄 Updating marker file patterns in state rule files..."
echo "   Project: $PROJECT_ROOT"
echo ""

# Find all rules.md files with old marker patterns
echo "📊 Scanning for files to update..."
files_with_markers=$(grep -rl "touch .state_rules_read_" agent-states/ --include="rules.md" 2>/dev/null || true)
file_count=$(echo "$files_with_markers" | grep -c . || echo "0")

if [ "$file_count" -eq 0 ]; then
    echo "   ℹ️  No files found with old marker patterns"
    echo "   ✅ All files already updated!"
    exit 0
fi

echo "   📋 Found $file_count files to update"
echo ""

# Update each file
update_count=0
for file in $files_with_markers; do
    echo "📝 Updating: $file"

    # Create backup
    cp "$file" "$file.bak"

    # Use perl for more reliable regex replacement
    # Replace: touch .state_rules_read_AGENT_STATE
    perl -i -pe 's/touch \.state_rules_read_([a-zA-Z_-]+)/mkdir -p markers\/state-verification \&\& touch "markers\/state-verification\/state_rules_read_$1-\$(date +%Y%m%d-%H%M%S)"/g' "$file"

    # Replace: echo "..." > .state_rules_read_AGENT_STATE
    perl -i -pe 's/> \.state_rules_read_([a-zA-Z_-]+)/> "markers\/state-verification\/state_rules_read_$1-\$(date +%Y%m%d-%H%M%S)"/g' "$file"

    # Replace: [ -f .state_rules_read_AGENT_STATE ]
    perl -i -pe 's/\[ -f \.state_rules_read_([a-zA-Z_-]+) \]/[ -n "\$(ls markers\/state-verification\/state_rules_read_$1-* 2>\/dev\/null | tail -1)" ]/g' "$file"

    # Check if file actually changed
    if ! diff -q "$file" "$file.bak" >/dev/null 2>&1; then
        echo "   ✅ Updated successfully"
        update_count=$((update_count + 1))
        rm "$file.bak"
    else
        echo "   ⚠️  No changes needed (pattern already correct)"
        mv "$file.bak" "$file"  # Restore original
    fi
done

echo ""
echo "✅ Update complete!"
echo "   📊 Files updated: $update_count / $file_count"
echo ""
echo "💡 Next steps:"
echo "   1. Review changes: git diff agent-states/"
echo "   2. Test marker creation in at least one state"
echo "   3. Commit changes: git add agent-states/ && git commit"
