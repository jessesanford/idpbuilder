#!/usr/bin/env python3
"""
Software Factory State Machine Comprehensive Validator (No External Dependencies)
Validates both structural integrity and logical consistency.
"""

import json
import sys
import os
from pathlib import Path
from typing import Dict, List, Set, Tuple, Optional, Any
from collections import defaultdict
from datetime import datetime
import re

# ANSI color codes for terminal output
class Colors:
    HEADER = '\033[95m'
    OKBLUE = '\033[94m'
    OKCYAN = '\033[96m'
    OKGREEN = '\033[92m'
    WARNING = '\033[93m'
    FAIL = '\033[91m'
    ENDC = '\033[0m'
    BOLD = '\033[1m'
    UNDERLINE = '\033[4m'

def colored(text: str, color: str) -> str:
    """Apply color to text for terminal output."""
    return f"{color}{text}{Colors.ENDC}"

def print_header(text: str):
    """Print a formatted header."""
    print(f"\n{colored('=' * 80, Colors.HEADER)}")
    print(f"{colored(text, Colors.HEADER + Colors.BOLD)}")
    print(f"{colored('=' * 80, Colors.HEADER)}")

def print_section(text: str):
    """Print a formatted section header."""
    print(f"\n{colored(text, Colors.OKBLUE + Colors.BOLD)}")
    print(colored("-" * len(text), Colors.OKBLUE))

def print_success(text: str):
    """Print success message."""
    print(colored(f"✅ {text}", Colors.OKGREEN))

def print_warning(text: str):
    """Print warning message."""
    print(colored(f"⚠️  {text}", Colors.WARNING))

def print_error(text: str):
    """Print error message."""
    print(colored(f"❌ {text}", Colors.FAIL))

def print_info(text: str):
    """Print info message."""
    print(colored(f"ℹ️  {text}", Colors.OKCYAN))

class StateMachineValidator:
    def __init__(self, state_machine_path: str):
        self.state_machine_path = Path(state_machine_path)
        self.state_machine = None
        self.errors = []
        self.warnings = []
        self.info = []

    def load_file(self) -> bool:
        """Load the state machine file."""
        try:
            with open(self.state_machine_path, 'r') as f:
                self.state_machine = json.load(f)
            print_success(f"Loaded state machine from {self.state_machine_path}")
            return True
        except FileNotFoundError:
            print_error(f"File not found: {self.state_machine_path}")
            return False
        except json.JSONDecodeError as e:
            print_error(f"Invalid JSON: {e.msg} at line {e.lineno}, column {e.colno}")
            return False
        except Exception as e:
            print_error(f"Error loading file: {str(e)}")
            return False

    def validate_structure(self) -> bool:
        """Validate basic structure of the state machine."""
        print_section("Structural Validation")

        valid = True

        # Check required top-level keys
        required_keys = ['metadata', 'agents', 'transition_matrix']
        missing_keys = []

        for key in required_keys:
            if key not in self.state_machine:
                missing_keys.append(key)
                print_error(f"Missing required top-level key: '{key}'")
                valid = False

        if not missing_keys:
            print_success("All required top-level keys present")

        # Check metadata structure
        if 'metadata' in self.state_machine:
            metadata = self.state_machine['metadata']
            meta_required = ['version', 'description', 'source', 'last_updated']
            for key in meta_required:
                if key not in metadata:
                    print_error(f"Missing required metadata field: '{key}'")
                    valid = False

        # Check agents structure
        if 'agents' in self.state_machine:
            required_agents = ['orchestrator', 'sw-engineer', 'code-reviewer', 'architect']
            agents = self.state_machine['agents']
            for agent in required_agents:
                if agent not in agents:
                    print_error(f"Missing required agent: '{agent}'")
                    valid = False

        return valid

    def get_all_states(self) -> Dict[str, Set[str]]:
        """Get all states organized by agent."""
        states = {}
        for agent_name, agent_def in self.state_machine.get('agents', {}).items():
            states[agent_name] = set(agent_def.get('states', {}).keys())
        return states

    def validate_transition_matrix(self) -> bool:
        """Validate that transition matrix references only existing states."""
        print_section("Transition Matrix Validation")

        valid = True
        all_states = self.get_all_states()
        transition_matrix = self.state_machine.get('transition_matrix', {})

        for agent_name, transitions in transition_matrix.items():
            if agent_name not in all_states:
                print_error(f"Agent '{agent_name}' in transition matrix not found in agents definition")
                valid = False
                continue

            agent_states = all_states[agent_name]

            # Check source states
            for source_state, target_states in transitions.items():
                if source_state not in agent_states:
                    print_error(f"Source state '{source_state}' for agent '{agent_name}' not defined")
                    valid = False

                # Check target states
                for target_state in target_states:
                    if target_state not in agent_states:
                        print_error(f"Target state '{target_state}' for agent '{agent_name}' not defined")
                        valid = False

        if valid:
            print_success("All transition matrix entries reference valid states")

        return valid

    def find_orphaned_states(self) -> Dict[str, List[str]]:
        """Find states with no incoming transitions (except initial states)."""
        print_section("Orphaned States Check")

        orphaned = {}
        transition_matrix = self.state_machine.get('transition_matrix', {})

        for agent_name, agent_def in self.state_machine.get('agents', {}).items():
            agent_transitions = transition_matrix.get(agent_name, {})
            states = agent_def.get('states', {})

            # Get all states that are targets of transitions
            target_states = set()
            for source, targets in agent_transitions.items():
                target_states.update(targets)

            # Find orphaned states (no incoming transitions)
            agent_orphaned = []
            for state_name, state_def in states.items():
                state_type = state_def.get('type', '')
                # Initial states are expected to have no incoming transitions
                if state_type != 'initial' and state_name not in target_states:
                    agent_orphaned.append(state_name)

            if agent_orphaned:
                orphaned[agent_name] = agent_orphaned

        if orphaned:
            print_warning("Found orphaned states (no incoming transitions):")
            for agent, states in orphaned.items():
                print_warning(f"  {agent}: {', '.join(states)}")
        else:
            print_success("No orphaned states found")

        return orphaned

    def find_dead_end_states(self) -> Dict[str, List[str]]:
        """Find non-terminal states with no outgoing transitions."""
        print_section("Dead-End States Check")

        dead_ends = {}
        transition_matrix = self.state_machine.get('transition_matrix', {})

        for agent_name, agent_def in self.state_machine.get('agents', {}).items():
            agent_transitions = transition_matrix.get(agent_name, {})
            states = agent_def.get('states', {})

            agent_dead_ends = []
            for state_name, state_def in states.items():
                state_type = state_def.get('type', '')
                # Terminal states are expected to have no outgoing transitions
                if state_type != 'terminal':
                    # Check if state has outgoing transitions
                    if state_name not in agent_transitions or not agent_transitions[state_name]:
                        agent_dead_ends.append(state_name)

            if agent_dead_ends:
                dead_ends[agent_name] = agent_dead_ends

        if dead_ends:
            print_error("Found dead-end states (non-terminal with no outgoing transitions):")
            for agent, states in dead_ends.items():
                print_error(f"  {agent}: {', '.join(states)}")
        else:
            print_success("No dead-end states found")

        return dead_ends

    def validate_terminal_states(self) -> bool:
        """Validate that terminal states are properly configured."""
        print_section("Terminal States Validation")

        valid = True
        transition_matrix = self.state_machine.get('transition_matrix', {})

        for agent_name, agent_def in self.state_machine.get('agents', {}).items():
            agent_transitions = transition_matrix.get(agent_name, {})
            states = agent_def.get('states', {})

            for state_name, state_def in states.items():
                state_type = state_def.get('type', '')
                valid_transitions = state_def.get('valid_transitions', [])

                if state_type == 'terminal':
                    # Terminal states should have empty valid_transitions
                    if valid_transitions:
                        print_error(f"Terminal state '{state_name}' in '{agent_name}' has non-empty valid_transitions: {valid_transitions}")
                        valid = False

                    # Terminal states should not appear in transition matrix or have empty list
                    if state_name in agent_transitions and agent_transitions[state_name]:
                        print_error(f"Terminal state '{state_name}' in '{agent_name}' has outgoing transitions in matrix")
                        valid = False

        if valid:
            print_success("All terminal states properly configured")

        return valid

    def validate_state_consistency(self) -> bool:
        """Check consistency between state definitions and transition matrix."""
        print_section("State Consistency Validation")

        valid = True
        transition_matrix = self.state_machine.get('transition_matrix', {})

        for agent_name, agent_def in self.state_machine.get('agents', {}).items():
            agent_transitions = transition_matrix.get(agent_name, {})
            states = agent_def.get('states', {})

            for state_name, state_def in states.items():
                valid_transitions = state_def.get('valid_transitions', [])
                matrix_transitions = agent_transitions.get(state_name, [])

                # Check if valid_transitions matches transition matrix
                valid_set = set(valid_transitions)
                matrix_set = set(matrix_transitions)

                if valid_set != matrix_set:
                    print_warning(f"Mismatch for state '{state_name}' in '{agent_name}':")

                    in_valid_not_matrix = valid_set - matrix_set
                    if in_valid_not_matrix:
                        print_warning(f"  In valid_transitions but not in matrix: {in_valid_not_matrix}")

                    in_matrix_not_valid = matrix_set - valid_set
                    if in_matrix_not_valid:
                        print_warning(f"  In matrix but not in valid_transitions: {in_matrix_not_valid}")

                    valid = False

        if valid:
            print_success("State definitions and transition matrix are consistent")

        return valid

    def generate_statistics(self):
        """Generate and display statistics about the state machine."""
        print_section("State Machine Statistics")

        stats = {
            'total_agents': len(self.state_machine.get('agents', {})),
            'total_states': 0,
            'state_types': defaultdict(int),
            'states_per_agent': {},
            'transitions_per_agent': {},
            'rules_referenced': len(self.state_machine.get('metadata', {}).get('rules_referenced', []))
        }

        for agent_name, agent_def in self.state_machine.get('agents', {}).items():
            states = agent_def.get('states', {})
            stats['states_per_agent'][agent_name] = len(states)
            stats['total_states'] += len(states)

            # Count state types
            for state_def in states.values():
                state_type = state_def.get('type', 'unknown')
                stats['state_types'][state_type] += 1

            # Count transitions
            transitions = self.state_machine.get('transition_matrix', {}).get(agent_name, {})
            total_transitions = sum(len(targets) for targets in transitions.values())
            stats['transitions_per_agent'][agent_name] = total_transitions

        # Display statistics
        print_info(f"Total Agents: {stats['total_agents']}")
        print_info(f"Total States: {stats['total_states']}")
        print_info(f"Total Rules Referenced: {stats['rules_referenced']}")

        print("\n  States per Agent:")
        for agent, count in stats['states_per_agent'].items():
            print(f"    {agent}: {count} states")

        print("\n  State Types Distribution:")
        for state_type, count in sorted(stats['state_types'].items()):
            print(f"    {state_type}: {count}")

        print("\n  Transitions per Agent:")
        for agent, count in stats['transitions_per_agent'].items():
            print(f"    {agent}: {count} transitions")

    def run_all_validations(self) -> bool:
        """Run all validation checks."""
        print_header("SOFTWARE FACTORY STATE MACHINE VALIDATOR")
        print_info(f"Validating: {self.state_machine_path}")
        print_info(f"Timestamp: {datetime.now().isoformat()}")

        if not self.load_file():
            return False

        results = []

        # Run all validations
        results.append(("Basic Structure", self.validate_structure()))
        results.append(("Transition Matrix", self.validate_transition_matrix()))
        results.append(("State Consistency", self.validate_state_consistency()))
        results.append(("Terminal States", self.validate_terminal_states()))

        # Check for orphaned and dead-end states
        orphaned = self.find_orphaned_states()
        dead_ends = self.find_dead_end_states()

        results.append(("Orphaned States", len(orphaned) == 0))
        results.append(("Dead-End States", len(dead_ends) == 0))

        # Generate statistics
        self.generate_statistics()

        # Summary
        print_header("VALIDATION SUMMARY")

        passed = sum(1 for _, result in results if result)
        failed = len(results) - passed

        print(f"\n  Total Checks: {len(results)}")
        print(colored(f"  Passed: {passed}", Colors.OKGREEN))
        if failed > 0:
            print(colored(f"  Failed: {failed}", Colors.FAIL))

        print("\n  Results by Check:")
        for check_name, result in results:
            status = "✅ PASS" if result else "❌ FAIL"
            color = Colors.OKGREEN if result else Colors.FAIL
            print(colored(f"    {check_name:.<35} {status}", color))

        overall_valid = all(result for _, result in results)

        print()
        if overall_valid:
            print_success("🎉 STATE MACHINE VALIDATION PROJECT_DONEFUL!")
        else:
            print_error("⚠️  STATE MACHINE VALIDATION FAILED!")
            print_info("Review the errors and warnings above for details.")

        return overall_valid

def main():
    """Main entry point."""
    # Determine paths
    script_dir = Path(__file__).parent
    project_root = script_dir.parent

    state_machine_path = project_root / "software-factory-3.0-state-machine.json"

    # Allow command-line override
    if len(sys.argv) > 1:
        state_machine_path = Path(sys.argv[1])

    # Create validator and run
    validator = StateMachineValidator(str(state_machine_path))
    success = validator.run_all_validations()

    # Exit with appropriate code
    sys.exit(0 if success else 1)

if __name__ == "__main__":
    main()
