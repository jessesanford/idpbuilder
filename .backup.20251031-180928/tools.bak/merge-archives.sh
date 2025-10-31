#!/bin/bash
# Restore archived historical data back to state file

set -e

STATE_FILE="${1:-orchestrator-state-v3.json}"
ARCHIVE_DIR="${2:-archives}"

if [[ ! -f "$STATE_FILE" ]]; then
    echo "ERROR: State file not found: $STATE_FILE"
    exit 1
fi

if [[ ! -d "$ARCHIVE_DIR" ]]; then
    echo "ERROR: Archive directory not found: $ARCHIVE_DIR"
    exit 1
fi

echo "Merging archives back to state file..."
echo ""

# Create backup
backup="${STATE_FILE}.backup-$(date +%Y%m%d-%H%M%S)"
cp "$STATE_FILE" "$backup"
echo "Backup created: $backup"
echo ""

# Merge state_transition_log
if [[ -f "$ARCHIVE_DIR/state-transition-log.json" ]]; then
    jq --slurpfile archive "$ARCHIVE_DIR/state-transition-log.json" \
       '.state_transition_log = $archive[0]' \
       "$STATE_FILE" > "${STATE_FILE}.tmp"
    mv "${STATE_FILE}.tmp" "$STATE_FILE"
    echo "✅ Restored: state_transition_log"
fi

# Merge phase_integration_results
if [[ -f "$ARCHIVE_DIR/phase-integration-results.json" ]]; then
    jq --slurpfile archive "$ARCHIVE_DIR/phase-integration-results.json" \
       '.phase_integration_results = $archive[0]' \
       "$STATE_FILE" > "${STATE_FILE}.tmp"
    mv "${STATE_FILE}.tmp" "$STATE_FILE"
    echo "✅ Restored: phase_integration_results"
fi

# Calculate size
final_size=$(wc -c < "$STATE_FILE")
final_tokens=$((final_size / 4))

echo ""
echo "✅ Merge complete!"
echo "New size: $final_size chars (~$final_tokens tokens)"
echo ""
echo "To commit: git add $STATE_FILE && git commit -m 'revert: restore archived historical data'"
