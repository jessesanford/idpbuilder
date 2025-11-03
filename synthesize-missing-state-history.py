#!/usr/bin/env python3
"""
Synthesize missing state_history entries from git commits
Reconstructs gaps between BUILD_VALIDATION (14:38:02Z) and current REVIEW_PHASE_INTEGRATION
"""

import json
import sys
from datetime import datetime

# Missing state history entries reconstructed from git commits
MISSING_ENTRIES = [
    {
        "from_state": "BUILD_VALIDATION",
        "to_state": "SETUP_PHASE_INFRASTRUCTURE",
        "timestamp": "2025-11-03T15:18:28Z",
        "validated_by": "state-manager",
        "reason": "Wave 2.3 build validation passed successfully. Final artifact built and verified per R323. Ready for phase-level integration.",
        "phase": 2,
        "wave": 3,
        "build_status": "PASSED",
        "rule": "R323",
        "source": "RECONSTRUCTED from git commit 46c7696cdee8cc09124e736d964dc0d8347cf410"
    },
    {
        "from_state": "SETUP_PHASE_INFRASTRUCTURE",
        "to_state": "START_PHASE_ITERATION",
        "timestamp": "2025-11-03T15:39:07Z",
        "validated_by": "state-manager",
        "reason": "Phase 2 integration container initialized. All 3 waves complete. Ready to begin phase integration iteration 1.",
        "phase": 2,
        "iteration": 1,
        "waves_completed": 3,
        "integration_branch": "phase-2-integration",
        "source": "RECONSTRUCTED from git commit 6a7437002d467262eddf8ce492a35703bc65c00c"
    },
    {
        "from_state": "START_PHASE_ITERATION",
        "to_state": "INTEGRATE_PHASE_WAVES",
        "timestamp": "2025-11-03T15:51:13Z",
        "validated_by": "state-manager",
        "reason": "Phase iteration 1 started. Iteration counter incremented, backport counter reset. Ready to integrate 3 waves (2.1, 2.2, 2.3).",
        "phase": 2,
        "iteration": 1,
        "backport_attempts": 0,
        "waves_to_integrate": [2.1, 2.2, 2.3],
        "source": "RECONSTRUCTED from git commit 5c601db8fe0e250be7b99c8b015522dcabcaabfa"
    },
    {
        "from_state": "INTEGRATE_PHASE_WAVES",
        "to_state": "REVIEW_PHASE_INTEGRATION",
        "timestamp": "2025-11-03T16:17:20Z",
        "validated_by": "state-manager",
        "reason": "Phase 2 wave integration complete. Waves 2.1, 2.2, 2.3 merged into phase-2-integration branch. Build PASSED, tests PASSED_WITH_KNOWN_ISSUES. Ready for phase integration review.",
        "phase": 2,
        "waves_integrated": 3,
        "conflicts": 0,
        "build_status": "PASSED",
        "test_status": "PASSED_WITH_KNOWN_ISSUES",
        "integration_report": ".software-factory/phase2/integration/INTEGRATE_PHASE_WAVES-REPORT--20251103-160832.md",
        "orchestrator_proposal": "REVIEW_PHASE_INTEGRATION",
        "proposal_accepted": True,
        "source": "RECONSTRUCTED from git commit b3857f287e2cecee07005be0091bfc765d052c34"
    }
]

def synthesize_state_history(state_file_path):
    """Add missing state history entries to orchestrator-state-v3.json"""

    print(f"📖 Reading state file: {state_file_path}")
    with open(state_file_path, 'r') as f:
        state_data = json.load(f)

    # Validate last entry matches expected
    last_entry = state_data['state_machine']['state_history'][-1]
    expected_last = {
        'from_state': 'REVIEW_WAVE_ARCHITECTURE',
        'to_state': 'BUILD_VALIDATION',
        'timestamp': '2025-11-03T14:38:02Z'
    }

    print(f"\n✓ Validating last existing entry...")
    print(f"  From: {last_entry['from_state']} → {last_entry['to_state']}")
    print(f"  Time: {last_entry['timestamp']}")

    if (last_entry['from_state'] != expected_last['from_state'] or
        last_entry['to_state'] != expected_last['to_state'] or
        last_entry['timestamp'] != expected_last['timestamp']):
        print(f"\n❌ ERROR: Last entry doesn't match expected!")
        print(f"   Expected: {expected_last['from_state']} → {expected_last['to_state']} @ {expected_last['timestamp']}")
        return False

    # Add missing entries
    print(f"\n✓ Adding {len(MISSING_ENTRIES)} missing state history entries...")
    for i, entry in enumerate(MISSING_ENTRIES, 1):
        print(f"\n  Entry {i}:")
        print(f"    {entry['from_state']} → {entry['to_state']}")
        print(f"    Timestamp: {entry['timestamp']}")
        print(f"    Source: {entry['source'][:60]}...")
        state_data['state_machine']['state_history'].append(entry)

    # Update last_transition_time to match newest entry
    newest_timestamp = MISSING_ENTRIES[-1]['timestamp']
    state_data['state_machine']['last_transition_time'] = newest_timestamp

    print(f"\n✓ Updated last_transition_time: {newest_timestamp}")

    # Verify current state is correct
    current_state = state_data['state_machine']['current_state']
    expected_current_state = 'REVIEW_PHASE_INTEGRATION'

    if current_state != expected_current_state:
        print(f"\n❌ WARNING: Current state is {current_state}, expected {expected_current_state}")
    else:
        print(f"✓ Current state verified: {current_state}")

    # Write updated state file
    print(f"\n💾 Writing updated state file...")
    with open(state_file_path, 'w') as f:
        json.dump(state_data, f, indent=2)

    print(f"\n✅ State history synthesis complete!")
    print(f"   Total entries: {len(state_data['state_machine']['state_history'])}")
    print(f"   Added: {len(MISSING_ENTRIES)} entries")

    return True

if __name__ == '__main__':
    if len(sys.argv) != 2:
        print("Usage: synthesize-missing-state-history.py <state-file-path>")
        sys.exit(1)

    state_file = sys.argv[1]
    success = synthesize_state_history(state_file)
    sys.exit(0 if success else 1)
