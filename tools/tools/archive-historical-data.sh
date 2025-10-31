#!/bin/bash
# SAFE STATE REDUCTION - Strategy 2: Historical Data Archival
# Moves completed/historical data to separate archive files
# WARNING: Only archive sections with ZERO active references!

set -e

STATE_FILE="${1:-orchestrator-state-v3.json}"
ARCHIVE_DIR="${2:-archives}"

if [[ ! -f "$STATE_FILE" ]]; then
    echo "ERROR: State file not found: $STATE_FILE"
    exit 1
fi

echo "SAFE STATE REDUCTION - Historical Data Archival"
echo "================================================"
echo ""
echo "⚠️  WARNING: This script only archives sections with ZERO active references!"
echo ""
echo "SAFE sections (verified zero references):"
echo "  - state_transition_log"
echo "  - phase_integration_results"
echo ""
echo "UNSAFE sections (have active references - DO NOT ARCHIVE):"
echo "  - spawn_timing (referenced in R003)"
echo "  - agents_spawned (referenced in R288, R313, etc.)"
echo "  - waves_completed (referenced in R288, R105, etc.)"
echo "  - review_results (referenced in R313)"
echo ""

read -p "Continue with archiving SAFE sections only? (y/N) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Aborted."
    exit 0
fi

# Create archive directory
mkdir -p "$ARCHIVE_DIR"

# Create backup
backup="${STATE_FILE}.backup-$(date +%Y%m%d-%H%M%S)"
cp "$STATE_FILE" "$backup"
echo "Backup created: $backup"
echo ""

# Calculate initial size
initial_size=$(jq -c '.' "$STATE_FILE" | wc -c)
initial_tokens=$((initial_size / 4))

echo "Archiving safe sections..."

# Archive state_transition_log
if jq -e '.state_transition_log' "$STATE_FILE" > /dev/null 2>&1; then
    jq '.state_transition_log' "$STATE_FILE" > "$ARCHIVE_DIR/state-transition-log.json"
    jq 'del(.state_transition_log)' "$STATE_FILE" > "${STATE_FILE}.tmp"
    mv "${STATE_FILE}.tmp" "$STATE_FILE"
    echo "✅ Archived: state_transition_log → $ARCHIVE_DIR/state-transition-log.json"
fi

# Archive phase_integration_results
if jq -e '.phase_integration_results' "$STATE_FILE" > /dev/null 2>&1; then
    jq '.phase_integration_results' "$STATE_FILE" > "$ARCHIVE_DIR/phase-integration-results.json"
    jq 'del(.phase_integration_results)' "$STATE_FILE" > "${STATE_FILE}.tmp"
    mv "${STATE_FILE}.tmp" "$STATE_FILE"
    echo "✅ Archived: phase_integration_results → $ARCHIVE_DIR/phase-integration-results.json"
fi

# Calculate final size
final_size=$(jq -c '.' "$STATE_FILE" | wc -c)
final_tokens=$((final_size / 4))
saved_chars=$((initial_size - final_size))
saved_tokens=$((initial_tokens - final_tokens))
percent_saved=$((saved_chars * 100 / initial_size))

echo ""
echo "✅ Archival complete!"
echo ""
echo "Results:"
echo "  Initial size: $initial_size chars (~$initial_tokens tokens)"
echo "  Final size:   $final_size chars (~$final_tokens tokens)"
echo "  Saved:        $saved_chars chars (~$saved_tokens tokens)"
echo "  Reduction:    ${percent_saved}%"
echo ""
echo "Archives created in: $ARCHIVE_DIR/"
echo ""

# Verify integrity
echo "Verifying state file integrity..."
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
echo "  1. Test thoroughly before committing"
echo "  2. Commit archives: git add $ARCHIVE_DIR/"
echo "  3. Commit state: git add $STATE_FILE"
echo "  4. Create commit: git commit -m 'perf: archive historical data (~${saved_tokens} tokens saved)'"
echo ""
echo "To restore: bash tools/merge-archives.sh"
