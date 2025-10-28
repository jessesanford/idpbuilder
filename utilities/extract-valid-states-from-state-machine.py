#!/usr/bin/env python3
"""
Extract all valid states from state machine and output as JSON schema enums.

This ensures the schema is ALWAYS in sync with the state machine.
"""

import json
import sys
from pathlib import Path

def extract_orchestrator_states(state_machine_path):
    """Extract all valid orchestrator states from state machine."""
    with open(state_machine_path, 'r') as f:
        data = json.load(f)

    # Extract from transition_matrix
    if 'transition_matrix' not in data:
        print("ERROR: No transition_matrix in state machine", file=sys.stderr)
        sys.exit(1)

    if 'orchestrator' not in data['transition_matrix']:
        print("ERROR: No orchestrator in transition_matrix", file=sys.stderr)
        sys.exit(1)

    # Get all states (keys of the transition matrix)
    orchestrator_states = sorted(data['transition_matrix']['orchestrator'].keys())

    return orchestrator_states

def extract_all_agent_states(state_machine_path):
    """Extract all valid states from all agents."""
    with open(state_machine_path, 'r') as f:
        data = json.load(f)

    all_states = set()

    # Extract from transition_matrix
    if 'transition_matrix' in data:
        for agent, states in data['transition_matrix'].items():
            all_states.update(states.keys())

    return sorted(all_states)

def main():
    # Paths
    repo_root = Path(__file__).parent.parent
    state_machine_path = repo_root / 'state-machines' / 'software-factory-3.0-state-machine.json'

    if not state_machine_path.exists():
        print(f"ERROR: State machine not found at {state_machine_path}", file=sys.stderr)
        sys.exit(1)

    # Extract states
    orchestrator_states = extract_orchestrator_states(state_machine_path)
    all_states = extract_all_agent_states(state_machine_path)

    # Output as JSON
    output = {
        "orchestrator_states": orchestrator_states,
        "all_agent_states": all_states,
        "count": {
            "orchestrator": len(orchestrator_states),
            "total": len(all_states)
        }
    }

    print(json.dumps(output, indent=2))

if __name__ == '__main__':
    main()
