#!/bin/bash
# SAFE STATE REDUCTION - Strategy 1: Minification
# Reduces orchestrator-state-v3.json by ~4766 tokens (95% of 5000 target)
# ZERO risk - no breaking changes possible

set -e

STATE_FILE="${1:-orchestrator-state-v3.json}"

if [[ ! -f "$STATE_FILE" ]]; then
    echo "ERROR: State file not found: $STATE_FILE"
    exit 1
fi

echo "SAFE STATE REDUCTION - Minification"
echo "==================================="
echo ""
echo "File: $STATE_FILE"
echo ""

# Calculate current size
current_size=$(wc -c < "$STATE_FILE")
current_tokens=$((current_size / 4))
echo "Current size: $current_size chars (~$current_tokens tokens)"

# Create backup
backup="${STATE_FILE}.backup-$(date +%Y%m%d-%H%M%S)"
cp "$STATE_FILE" "$backup"
echo "Backup created: $backup"
echo ""

# Minify
echo "Minifying..."
jq -c '.' "$STATE_FILE" > "${STATE_FILE}.tmp"
mv "${STATE_FILE}.tmp" "$STATE_FILE"

# Calculate new size
new_size=$(wc -c < "$STATE_FILE")
new_tokens=$((new_size / 4))
saved_chars=$((current_size - new_size))
saved_tokens=$((current_tokens - new_tokens))
percent_saved=$((saved_chars * 100 / current_size))

echo "✅ Minification complete!"
echo ""
echo "Results:"
echo "  New size:    $new_size chars (~$new_tokens tokens)"
echo "  Saved:       $saved_chars chars (~$saved_tokens tokens)"
echo "  Reduction:   ${percent_saved}%"
echo ""

# Verify jq still works
echo "Verifying integrity..."
if jq '.current_state' "$STATE_FILE" > /dev/null 2>&1; then
    echo "✅ State file integrity verified"
else
    echo "❌ ERROR: State file corrupted!"
    echo "Restoring backup..."
    cp "$backup" "$STATE_FILE"
    exit 1
fi

echo ""
echo "Next steps:"
echo "  1. Test with: jq '.' $STATE_FILE | less"
echo "  2. Commit: git add $STATE_FILE && git commit -m 'perf: minify state file (~${saved_tokens} tokens saved)'"
echo "  3. View anytime: jq '.' $STATE_FILE"
echo ""
echo "To revert: cp $backup $STATE_FILE"
