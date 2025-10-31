#!/bin/bash
# Restore pretty-printing to minified state file

set -e

STATE_FILE="${1:-orchestrator-state-v3.json}"

if [[ ! -f "$STATE_FILE" ]]; then
    echo "ERROR: State file not found: $STATE_FILE"
    exit 1
fi

echo "Restoring pretty-printing to: $STATE_FILE"

# Create backup
backup="${STATE_FILE}.backup-$(date +%Y%m%d-%H%M%S)"
cp "$STATE_FILE" "$backup"
echo "Backup created: $backup"

# Pretty-print
jq '.' "$STATE_FILE" > "${STATE_FILE}.tmp"
mv "${STATE_FILE}.tmp" "$STATE_FILE"

# Calculate sizes
current_size=$(wc -c < "$STATE_FILE")
current_tokens=$((current_size / 4))

echo "✅ Pretty-printing restored"
echo "New size: $current_size chars (~$current_tokens tokens)"
echo ""
echo "To commit: git add $STATE_FILE && git commit -m 'revert: restore pretty-printing to state file'"
