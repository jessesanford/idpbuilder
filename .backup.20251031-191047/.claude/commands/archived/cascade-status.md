---
tags: ["project", "gitignored"]
description: Display comprehensive R406 fix cascade status report
name: /cascade-status
---

# Cascade Status Command

Display comprehensive status of active R406 fix cascades.

## Purpose

Shows detailed progress report for any active fix cascade, including:
- Bug registry status
- Per-integration breakdown
- Per-effort status
- Next actions
- Validation checksums

## Usage

```bash
/cascade-status
```

## When to Use

- Check progress of fix cascade
- See which bugs are pending/in_progress/fixed
- Understand what needs to happen next
- Verify no bugs are lost
- Debug stuck/blocked fixes

## Implementation

```bash
#!/bin/bash
# Display R406 cascade status report

cd "$CLAUDE_PROJECT_DIR" || exit 1

# Check if cascade is active
CASCADE_ACTIVE=$(jq -r '.fix_cascade_state.active // false' orchestrator-state-v3.json 2>/dev/null)

if [ "$CASCADE_ACTIVE" != "true" ]; then
    echo "No active fix cascade detected."
    echo ""
    echo "Cascade state: inactive"
    exit 0
fi

# Source and run the status report
if [ -f "utilities/cascade-status-report.sh" ]; then
    echo "═══════════════════════════════════════════════════════════"
    echo "📊 R406 FIX CASCADE STATUS"
    echo "═══════════════════════════════════════════════════════════"
    echo ""

    source utilities/cascade-status-report.sh
    cascade_status_report

    echo ""
    echo "═══════════════════════════════════════════════════════════"
else
    echo "ERROR: cascade-status-report.sh not found!"
    echo "Expected at: utilities/cascade-status-report.sh"
    exit 1
fi
```

## Example Output

```
═══════════════════════════════════════════════════════════
📊 R406 FIX CASCADE STATUS
═══════════════════════════════════════════════════════════

Cascade Overview:
  Cascade ID: cascade-20251004-160000
  Status: fixing
  Layer: 1 / 2
  Triggered by: phase1_wave2
  Started: 2025-10-04T16:00:00Z

Bug Registry Summary:
  Total Bugs: 5
  ✅ Verified: 2
  🔧 In Progress: 1
  ⏳ Pending: 2

Integration Status:
  phase1_wave1:
    Total Bugs: 2
    ✅ Fixed: 1
    ⏳ Pending: 1
    ⚠️ Requires Reintegration: Yes

Effort Status:
  ✅ E1.1.3: 1 bug fixed
  🔧 E1.2.1-command-structure: 1 bug in progress
  ⏳ E1.1.2-split-002: 1 bug pending

Next Actions:
  1. Complete fix for BUG-001 (in progress)
  2. Start fix for BUG-002 (pending)
  3. Reintegrate Wave 1 after both bugs fixed

Validation:
  ✅ Bug count valid: 5 bugs accounted for
  ✅ Checksum valid: no lost bugs
  Last validated: 2025-10-04T16:30:00Z
═══════════════════════════════════════════════════════════
```

## Related

- **R406**: Fix Cascade Tracking Protocol
- **utilities/cascade-status-report.sh**: Report implementation
- **State Rules**: Automatic reporting at end of states
- **R327**: Mandatory Re-Integration After Fixes

## Notes

- This command is available anytime during cascade operations
- Orchestrator states also print this automatically at end of execution
- Use when you need on-demand status check
- No side effects - read-only operation
