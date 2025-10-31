#!/usr/bin/env python3
"""
Synchronize schema state enums with state machine.

This ensures the schema ALWAYS matches the state machine exactly.
CRITICAL: This prevents invalid states from passing validation.
"""

import json
import sys
from pathlib import Path

def extract_orchestrator_states(state_machine_path):
    """Extract all valid orchestrator states from state machine."""
    with open(state_machine_path, 'r') as f:
        data = json.load(f)

    if 'transition_matrix' not in data or 'orchestrator' not in data['transition_matrix']:
        print("ERROR: Cannot find orchestrator states in state machine", file=sys.stderr)
        sys.exit(1)

    return sorted(data['transition_matrix']['orchestrator'].keys())

def extract_all_agent_states(state_machine_path):
    """Extract all valid states from all agents (for previous_state validation)."""
    with open(state_machine_path, 'r') as f:
        data = json.load(f)

    all_states = set()

    if 'transition_matrix' in data:
        for agent, states in data['transition_matrix'].items():
            all_states.update(states.keys())

    return sorted(all_states)

def update_schema_states(schema_path, orchestrator_states, all_states):
    """Update schema with correct state enums."""
    with open(schema_path, 'r') as f:
        schema = json.load(f)

    changes_made = False

    # Update current_state enum (orchestrator states only)
    if 'properties' in schema and 'current_state' in schema['properties']:
        current = schema['properties']['current_state'].get('enum', [])
        if set(current) != set(orchestrator_states):
            schema['properties']['current_state']['enum'] = orchestrator_states
            changes_made = True
            print(f"✓ Updated current_state enum: {len(orchestrator_states)} states")

            # Report missing/extra states
            missing = set(orchestrator_states) - set(current)
            extra = set(current) - set(orchestrator_states)
            if missing:
                print(f"  Added missing states: {', '.join(sorted(missing))}")
            if extra:
                print(f"  Removed invalid states: {', '.join(sorted(extra))}")

    # Update previous_state enum (all agent states + null)
    if 'properties' in schema and 'previous_state' in schema['properties']:
        current = schema['properties']['previous_state'].get('enum', [])
        # Filter out null to compare
        current_non_null = [s for s in current if s is not None]
        target_with_null = [None] + all_states

        if set(current_non_null) != set(all_states):
            schema['properties']['previous_state']['enum'] = target_with_null
            changes_made = True
            print(f"✓ Updated previous_state enum: {len(all_states)} states + null")

            # Report changes
            missing = set(all_states) - set(current_non_null)
            extra = set(current_non_null) - set(all_states)
            if missing:
                print(f"  Added missing states: {', '.join(sorted(missing))}")
            if extra:
                print(f"  Removed invalid states: {', '.join(sorted(extra))}")

    # Update next_state enum if it exists
    if 'properties' in schema and 'next_state' in schema['properties']:
        current = schema['properties']['next_state'].get('enum', [])
        current_non_null = [s for s in current if s is not None and s != ""]
        target_with_null = [None, ""] + all_states  # Can be null or empty string

        if set(current_non_null) != set(all_states):
            schema['properties']['next_state']['enum'] = target_with_null
            changes_made = True
            print(f"✓ Updated next_state enum: {len(all_states)} states + null + empty")

    return schema, changes_made

def main():
    # Paths
    repo_root = Path(__file__).parent.parent
    state_machine_path = repo_root / 'state-machines' / 'software-factory-3.0-state-machine.json'
    schema_path = repo_root / 'orchestrator-state.schema.json'

    if not state_machine_path.exists():
        print(f"ERROR: State machine not found at {state_machine_path}", file=sys.stderr)
        sys.exit(1)

    if not schema_path.exists():
        print(f"ERROR: Schema not found at {schema_path}", file=sys.stderr)
        sys.exit(1)

    print("🔍 Extracting states from state machine...")
    orchestrator_states = extract_orchestrator_states(state_machine_path)
    all_states = extract_all_agent_states(state_machine_path)

    print(f"📊 Found {len(orchestrator_states)} orchestrator states")
    print(f"📊 Found {len(all_states)} total agent states\n")

    print("🔄 Updating schema...")
    schema, changes_made = update_schema_states(schema_path, orchestrator_states, all_states)

    if changes_made:
        print("\n💾 Writing updated schema...")
        with open(schema_path, 'w') as f:
            json.dump(schema, f, indent=2)
        print(f"✅ Schema updated successfully at {schema_path}")
        sys.exit(0)
    else:
        print("✅ Schema already up-to-date (no changes needed)")
        sys.exit(0)

if __name__ == '__main__':
    main()
