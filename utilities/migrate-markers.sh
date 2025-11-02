#!/bin/bash
# migrate-markers.sh - Migrate old marker files to new structure
# Part of marker file reorganization (R290 update)

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$PROJECT_ROOT"

echo "🔄 Migrating marker files to new structure..."
echo "   Project: $PROJECT_ROOT"
echo ""

# Create new directories
echo "📁 Creating marker directories..."
mkdir -p markers/state-verification
mkdir -p markers/state-backups
mkdir -p markers/coordination
echo "   ✅ Created: markers/state-verification/"
echo "   ✅ Created: markers/state-backups/"
echo "   ✅ Created: markers/coordination/"
echo ""

# Migrate state verification markers
echo "🔄 Migrating state verification markers..."
marker_count=0
for marker in .state_rules_read_* .state_rules_*; do
    [ -f "$marker" ] || continue
    base=$(basename "$marker")
    # Remove leading dot
    base_no_dot="${base#.}"
    new_name="markers/state-verification/${base_no_dot}-$(date +%Y%m%d-%H%M%S)"
    mv "$marker" "$new_name"
    echo "   ✅ Migrated: $marker → $new_name"
    marker_count=$((marker_count + 1))
done

if [ $marker_count -eq 0 ]; then
    echo "   ℹ️  No state verification markers found to migrate"
else
    echo "   ✅ Migrated $marker_count state verification markers"
fi
echo ""

# Migrate state backups
echo "🔄 Migrating state backups..."
if [ -d .state-backup ]; then
    backup_count=$(find .state-backup -mindepth 1 -maxdepth 1 -type d | wc -l)
    if [ "$backup_count" -gt 0 ]; then
        mv .state-backup/* markers/state-backups/ 2>/dev/null || true
        rmdir .state-backup 2>/dev/null || true
        echo "   ✅ Migrated: .state-backup/ → markers/state-backups/ ($backup_count backups)"
    else
        rmdir .state-backup 2>/dev/null || true
        echo "   ℹ️  .state-backup/ was empty, removed"
    fi
else
    echo "   ℹ️  No .state-backup/ directory found"
fi
echo ""

# Summary
echo "✅ Migration complete!"
echo ""
echo "📊 New structure:"
echo "   markers/"
echo "   ├── state-verification/  ($(ls -1 markers/state-verification/ 2>/dev/null | wc -l) files)"
echo "   ├── state-backups/       ($(ls -1d markers/state-backups/*/ 2>/dev/null | wc -l) backups)"
echo "   └── coordination/        (reserved for future use)"
echo ""
echo "💡 Old marker patterns are still supported for backward compatibility"
echo "   but will be automatically migrated on detection."
