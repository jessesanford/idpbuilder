#!/bin/bash
# Migrate all orchestrator-state-v3.json (v2) references to orchestrator-state-v3.json (SF 3.0)

set -e

echo "🔄 Migrating v2 state file references to v3..."

# Backup before migration
echo "📦 Creating backup..."
git stash push -m "Pre-v3-migration backup $(date +%Y%m%d-%H%M%S)" 2>/dev/null || echo "  (No changes to stash)"

# Count references before
BEFORE_COUNT=$(grep -r "orchestrator-state\.json" . --exclude-dir=.git --exclude-dir=node_modules 2>/dev/null | wc -l)
echo "📊 Found $BEFORE_COUNT v2 references to migrate"

# Perform migration (exclude .git and node_modules)
echo "🔧 Replacing references..."
find . -type f \( -name "*.md" -o -name "*.sh" -o -name "*.py" -o -name "*.json" \) \
  -not -path "./.git/*" \
  -not -path "*/node_modules/*" \
  -not -path "*/docs/STATE-FILE-V2-TO-V3-MIGRATION.md" \
  -exec sed -i 's/orchestrator-state\.json/orchestrator-state-v3.json/g' {} +

# Count references after
AFTER_COUNT=$(grep -r "orchestrator-state\.json" . --exclude-dir=.git --exclude-dir=node_modules 2>/dev/null | wc -l || echo "0")
echo "📊 Remaining v2 references: $AFTER_COUNT"

if [ "$AFTER_COUNT" -eq "0" ]; then
  echo "✅ Migration complete! All v2 references updated to v3"
else
  echo "⚠️  Warning: $AFTER_COUNT references still remain"
  echo ""
  echo "Remaining references:"
  grep -r "orchestrator-state\.json" . --exclude-dir=.git --exclude-dir=node_modules 2>/dev/null
fi

echo ""
echo "Next steps:"
echo "1. Review changes: git diff"
echo "2. Commit: git add -A && git commit -m 'fix: Migrate all orchestrator-state-v3.json (v2) to orchestrator-state-v3.json (SF 3.0) [R600]'"
echo "3. Verify: grep -r 'orchestrator-state\.json' . --exclude-dir=.git"
