#!/usr/bin/env python3
"""
SOFTWARE FACTORY 2.0 - STATE POSITION VALIDATOR
Validates current state against state machine definition
"""

import json
import sys
from pathlib import Path
from typing import Dict, List, Optional, Tuple

class StateValidator:
    def __init__(self):
        self.project_root = self._find_project_root()
        self.state_machine = {}
        self.current_state = {}
        self.errors = []
        self.warnings = []

    def _find_project_root(self) -> Path:
        """Find project root directory"""
        for possible_root in [
            Path.cwd(),
            Path("/home/vscode/software-factory-template"),
            Path("/workspaces/software-factory-2.0-template")
        ]:
            if (possible_root / "orchestrator-state-v3.json").exists() or \
               (possible_root / "orchestrator-state-demo.json").exists():
                return possible_root
        return Path.cwd()

    def load_data(self) -> bool:
        """Load state machine and current state"""
        # Load state machine
        state_machine_paths = [
            self.project_root / "state-machines" / "software-factory-3.0-state-machine.json",
            self.project_root / "software-factory-3.0-state-machine.json"
        ]

        for path in state_machine_paths:
            if path.exists():
                try:
                    with open(path, 'r') as f:
                        self.state_machine = json.load(f)
                        print(f"✓ Loaded state machine from: {path}")
                        break
                except Exception as e:
                    self.errors.append(f"Failed to load state machine: {e}")
                    return False

        if not self.state_machine:
            self.errors.append("No state machine definition found")
            return False

        # Load current state
        state_path = self.project_root / "orchestrator-state-v3.json"
        if not state_path.exists():
            self.warnings.append("No orchestrator-state-v3.json found, using example")
            state_path = self.project_root / "orchestrator-state-demo.json"

        if state_path.exists():
            try:
                with open(state_path, 'r') as f:
                    self.current_state = json.load(f)
                    print(f"✓ Loaded state from: {state_path}")
            except Exception as e:
                self.errors.append(f"Failed to load current state: {e}")
                return False
        else:
            self.errors.append("No state file found")
            return False

        return True

    def get_all_states(self) -> Dict[str, Dict]:
        """Extract all valid states from state machine"""
        states = {}

        # Get orchestrator states
        if "agents" in self.state_machine:
            if "orchestrator" in self.state_machine["agents"]:
                orch_states = self.state_machine["agents"]["orchestrator"].get("states", {})
                for state_name, state_data in orch_states.items():
                    states[state_name] = state_data

        # Check for states in sequences
        if "state_sequences" in self.state_machine:
            for seq_name, sequence in self.state_machine["state_sequences"].items():
                if isinstance(sequence, list):
                    for state in sequence:
                        if state not in states:
                            states[state] = {
                                "type": "sequence",
                                "description": f"State from {seq_name} sequence"
                            }

        return states

    def validate_current_state(self) -> bool:
        """Validate the current state"""
        current = self.current_state.get("current_state", "UNKNOWN")
        previous = self.current_state.get("previous_state")

        states = self.get_all_states()

        # Check if current state exists
        if current not in states:
            self.errors.append(f"Current state '{current}' is not in state machine definition!")
            return False

        print(f"✓ Current state '{current}' is valid")

        # Check if transition from previous state was valid
        if previous and previous in states:
            valid_transitions = states[previous].get("valid_transitions", [])
            if valid_transitions and current not in valid_transitions:
                self.warnings.append(
                    f"Transition from '{previous}' to '{current}' is not in valid_transitions list"
                )
                print(f"  Valid transitions from {previous}: {', '.join(valid_transitions)}")

        return True

    def validate_sub_state_machine(self) -> bool:
        """Validate sub-state machine configuration"""
        sub_state = self.current_state.get("sub_state_machine", {})

        if not sub_state.get("active"):
            print("✓ No sub-state machine active")
            return True

        sub_type = sub_state.get("type")
        sub_current = sub_state.get("current_state")
        return_state = sub_state.get("return_state")

        # Validate sub-state type
        valid_sub_types = ["FIX_CASCADE", "PR_READY", "INITIALIZATION", "INTEGRATE_WAVE_EFFORTS", "SPLITTING"]
        if sub_type not in valid_sub_types:
            self.warnings.append(f"Sub-state machine type '{sub_type}' may not be valid")
        else:
            print(f"✓ Sub-state machine type '{sub_type}' is valid")

        # Check return state exists
        states = self.get_all_states()
        if return_state and return_state not in states:
            self.warnings.append(f"Sub-state return state '{return_state}' not found in main state machine")

        return True

    def validate_phase_wave(self) -> bool:
        """Validate phase and wave numbers"""
        phase = self.current_state.get("current_phase")
        wave = self.current_state.get("current_wave")

        if phase is not None:
            if not isinstance(phase, int) or phase < 0:
                self.errors.append(f"Invalid phase number: {phase}")
                return False
            print(f"✓ Phase {phase} is valid")

        if wave is not None:
            if not isinstance(wave, int) or wave < 0:
                self.errors.append(f"Invalid wave number: {wave}")
                return False
            print(f"✓ Wave {wave} is valid")

        return True

    def validate_efforts(self) -> bool:
        """Validate effort tracking"""
        efforts_planned = self.current_state.get("efforts_planned", [])
        efforts_in_progress = self.current_state.get("efforts_in_progress", [])
        efforts_completed = self.current_state.get("efforts_completed", [])

        # Check for duplicates
        all_efforts = []
        for effort in efforts_planned + efforts_in_progress + efforts_completed:
            if isinstance(effort, dict):
                effort_id = effort.get("id", effort.get("name"))
            else:
                effort_id = effort

            if effort_id in all_efforts:
                self.warnings.append(f"Duplicate effort found: {effort_id}")
            all_efforts.append(effort_id)

        print(f"✓ Effort tracking: {len(efforts_planned)} planned, "
              f"{len(efforts_in_progress)} in progress, {len(efforts_completed)} completed")

        return True

    def check_required_fields(self) -> bool:
        """Check for required fields in state file"""
        required_fields = [
            "current_state",
            "current_phase",
            "current_wave"
        ]

        missing = []
        for field in required_fields:
            if field not in self.current_state:
                missing.append(field)

        if missing:
            self.errors.append(f"Missing required fields: {', '.join(missing)}")
            return False

        print(f"✓ All required fields present")
        return True

    def check_state_rules(self) -> bool:
        """Check state-specific rules"""
        current = self.current_state.get("current_state")
        states = self.get_all_states()

        if current not in states:
            return False

        state_info = states[current]
        state_type = state_info.get("type", "")

        # Check spawn states (R313)
        if state_type == "spawn" or "SPAWN" in current:
            print(f"⚠ State '{current}' is a spawn state - R313 requires STOP after spawn")

        # Check checkpoint states (R322)
        checkpoint_states = [
            "WAITING_FOR_MERGE_PLAN",
            "WAITING_FOR_PHASE_MERGE_PLAN",
            "WAITING_FOR_PROJECT_MERGE_PLAN",
            "WAITING_FOR_FIX_PLANS",
            "WAITING_FOR_BACKPORT_PLAN"
        ]

        if current in checkpoint_states:
            print(f"⚠ State '{current}' requires user checkpoint before transition (R322)")

        return True

    def generate_report(self) -> str:
        """Generate validation report"""
        lines = []
        lines.append("=" * 70)
        lines.append("SOFTWARE FACTORY STATE VALIDATION REPORT")
        lines.append("=" * 70)
        lines.append("")

        # Current position
        current = self.current_state.get("current_state", "UNKNOWN")
        phase = self.current_state.get("current_phase", "?")
        wave = self.current_state.get("current_wave", "?")

        lines.append(f"Current Position: {current}")
        lines.append(f"Phase: {phase}, Wave: {wave}")
        lines.append("")

        # Validation results
        if not self.errors:
            lines.append("✅ VALIDATION PASSED")
        else:
            lines.append("❌ VALIDATION FAILED")

        lines.append("")

        # Errors
        if self.errors:
            lines.append("ERRORS:")
            for error in self.errors:
                lines.append(f"  ❌ {error}")
            lines.append("")

        # Warnings
        if self.warnings:
            lines.append("WARNINGS:")
            for warning in self.warnings:
                lines.append(f"  ⚠️  {warning}")
            lines.append("")

        # State details
        states = self.get_all_states()
        if current in states:
            state_info = states[current]
            lines.append("Current State Details:")
            lines.append(f"  Type: {state_info.get('type', 'unknown')}")
            lines.append(f"  Description: {state_info.get('description', 'N/A')}")

            valid_transitions = state_info.get("valid_transitions", [])
            if valid_transitions:
                lines.append(f"  Valid Next States: {', '.join(valid_transitions[:5])}")
                if len(valid_transitions) > 5:
                    lines.append(f"                     ... and {len(valid_transitions) - 5} more")

        lines.append("")
        lines.append("-" * 70)

        return "\n".join(lines)

    def run(self) -> int:
        """Run validation"""
        print("Starting state validation...")
        print()

        if not self.load_data():
            print("Failed to load data files")
            return 1

        # Run all validations
        self.check_required_fields()
        self.validate_current_state()
        self.validate_sub_state_machine()
        self.validate_phase_wave()
        self.validate_efforts()
        self.check_state_rules()

        # Generate report
        print()
        print(self.generate_report())

        return 0 if not self.errors else 1

def main():
    """Main entry point"""
    import argparse

    parser = argparse.ArgumentParser(description="Validate Software Factory state position")
    parser.add_argument("-q", "--quiet", action="store_true",
                       help="Only show errors and warnings")
    parser.add_argument("-j", "--json", action="store_true",
                       help="Output results as JSON")

    args = parser.parse_args()

    validator = StateValidator()

    if args.json:
        # JSON output mode
        result = {
            "valid": not validator.errors,
            "errors": validator.errors,
            "warnings": validator.warnings,
            "current_state": validator.current_state.get("current_state"),
            "phase": validator.current_state.get("current_phase"),
            "wave": validator.current_state.get("current_wave")
        }
        print(json.dumps(result, indent=2))
        return 0 if result["valid"] else 1

    return validator.run()

if __name__ == "__main__":
    sys.exit(main())