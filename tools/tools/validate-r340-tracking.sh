#!/bin/bash
# R340 Metadata Tracking Validation Script
# Validates that orchestrator-state-v3.json complies with R340 requirements

set -e

STATE_FILE="${1:-orchestrator-state-v3.json}"

if [ ! -f "$STATE_FILE" ]; then
  echo "Usage: $0 <state-file>"
  echo "State file not found: $STATE_FILE"
  exit 1
fi

echo "🔍 Validating R340 metadata tracking in: $STATE_FILE"
echo ""

# Check if R340 fields exist
echo "Checking R340 field existence..."
MISSING_FIELDS=0

if ! jq -e '.planning_repo_files' "$STATE_FILE" >/dev/null 2>&1; then
  echo "  ❌ MISSING: .planning_repo_files field"
  MISSING_FIELDS=$((MISSING_FIELDS + 1))
else
  echo "  ✅ .planning_repo_files exists"
fi

if ! jq -e '.effort_repo_files' "$STATE_FILE" >/dev/null 2>&1; then
  echo "  ❌ MISSING: .effort_repo_files field"
  MISSING_FIELDS=$((MISSING_FIELDS + 1))
else
  echo "  ✅ .effort_repo_files exists"
fi

if [ $MISSING_FIELDS -gt 0 ]; then
  echo ""
  echo "🚨 R340 VIOLATION: Missing required fields"
  exit 1
fi

echo ""

# Check tracked files exist (only from R340 fields)
echo "Checking tracked file existence..."

# Only check file_path values from planning_repo_files and effort_repo_files
PLANNING_FILES=$(jq -r '.planning_repo_files | .. | objects | select(has("file_path")) | .file_path' "$STATE_FILE" 2>/dev/null || echo "")
EFFORT_FILES=$(jq -r '.effort_repo_files | .. | objects | select(has("file_path")) | .file_path' "$STATE_FILE" 2>/dev/null || echo "")

TRACKED_FILES=$(echo -e "$PLANNING_FILES\n$EFFORT_FILES" | grep -v '^$' || echo "")
MISSING_COUNT=0
CHECKED_COUNT=0

if [ -n "$TRACKED_FILES" ]; then
  while IFS= read -r path; do
    CHECKED_COUNT=$((CHECKED_COUNT + 1))
    if [ ! -f "$path" ]; then
      echo "  ⚠️  File missing: $path"
      MISSING_COUNT=$((MISSING_COUNT + 1))
    fi
  done <<< "$TRACKED_FILES"

  if [ $MISSING_COUNT -eq 0 ]; then
    echo "  ✅ All $CHECKED_COUNT tracked files exist"
  else
    echo "  ⚠️  $MISSING_COUNT of $CHECKED_COUNT tracked files are missing (may be deprecated)"
  fi
else
  echo "  ℹ️  No files tracked yet (R340 fields are empty - this is OK for new projects)"
fi

echo ""

# Check R383 compliance (paths under .software-factory/ or planning/)
echo "Checking R383 compliance (file locations)..."

# Only check paths from R340 fields
PLANNING_VIOLATIONS=$(jq -r '.planning_repo_files | .. | objects | select(has("file_path")) | .file_path | select((startswith(".software-factory/") or startswith("planning/")) | not)' "$STATE_FILE" 2>/dev/null || echo "")
EFFORT_VIOLATIONS=$(jq -r '.effort_repo_files | .. | objects | select(has("file_path")) | .file_path | select((startswith(".software-factory/") or startswith("planning/")) | not)' "$STATE_FILE" 2>/dev/null || echo "")

VIOLATIONS=$(echo -e "$PLANNING_VIOLATIONS\n$EFFORT_VIOLATIONS" | grep -v '^$' || echo "")

if [ -n "$VIOLATIONS" ]; then
  echo "  ❌ R383 VIOLATIONS: Files outside .software-factory/ or planning/:"
  echo "$VIOLATIONS" | sed 's/^/    /'
  echo ""
  echo "🚨 R383 VIOLATION: Metadata files in wrong locations"
  exit 1
else
  echo "  ✅ All tracked files comply with R383 (or no files tracked yet)"
fi

echo ""
echo "✅ R340 metadata tracking validation PASSED"
echo ""
exit 0
