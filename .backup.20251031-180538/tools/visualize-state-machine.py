#!/usr/bin/env python3
"""
SOFTWARE FACTORY 2.0 - STATE MACHINE VISUALIZER
Generates ASCII art visualization of the current state machine position
"""

import json
import os
import sys
from datetime import datetime
from typing import Dict, List, Optional, Tuple, Set
from pathlib import Path

# ANSI color codes for terminal output
class Colors:
    RESET = '\033[0m'
    BOLD = '\033[1m'
    RED = '\033[91m'
    GREEN = '\033[92m'
    YELLOW = '\033[93m'
    BLUE = '\033[94m'
    MAGENTA = '\033[95m'
    CYAN = '\033[96m'
    WHITE = '\033[97m'
    GRAY = '\033[90m'

    # Background colors
    BG_RED = '\033[41m'
    BG_GREEN = '\033[42m'
    BG_YELLOW = '\033[43m'
    BG_BLUE = '\033[44m'
    BG_MAGENTA = '\033[45m'
    BG_CYAN = '\033[46m'

def load_json_file(filepath: str) -> Dict:
    """Load JSON file with error handling"""
    try:
        with open(filepath, 'r') as f:
            return json.load(f)
    except FileNotFoundError:
        print(f"{Colors.RED}Error: File not found: {filepath}{Colors.RESET}")
        return {}
    except json.JSONDecodeError as e:
        print(f"{Colors.RED}Error parsing JSON in {filepath}: {e}{Colors.RESET}")
        return {}

def get_project_root() -> Path:
    """Get the project root directory"""
    # Try multiple possible locations
    for possible_root in [
        Path.cwd(),
        Path.home() / "software-factory-template",
        Path("/home/vscode/software-factory-template"),
        Path("/workspaces/software-factory-2.0-template")
    ]:
        if (possible_root / "orchestrator-state-v3.json").exists() or \
           (possible_root / "orchestrator-state-demo.json").exists():
            return possible_root
    return Path.cwd()

class StateVisualizer:
    def __init__(self):
        self.project_root = get_project_root()
        self.state_file = self.project_root / "orchestrator-state-v3.json"
        self.state_machine_file = self.project_root / "state-machines" / "software-factory-3.0-state-machine.json"

        # Fallback to root if state-machines doesn't exist
        if not self.state_machine_file.exists():
            self.state_machine_file = self.project_root / "software-factory-3.0-state-machine.json"

        self.current_state = {}
        self.state_machine = {}
        self.state_history = []

    def load_data(self) -> bool:
        """Load state files"""
        # Load current state
        self.current_state = load_json_file(str(self.state_file))
        if not self.current_state:
            print(f"{Colors.YELLOW}Warning: Could not load current state, using example{Colors.RESET}")
            self.state_file = self.project_root / "orchestrator-state-demo.json"
            self.current_state = load_json_file(str(self.state_file))

        # Load state machine definition
        self.state_machine = load_json_file(str(self.state_machine_file))
        if not self.state_machine:
            print(f"{Colors.RED}Error: Could not load state machine definition{Colors.RESET}")
            return False

        # Extract state history if available
        if "state_history" in self.current_state:
            self.state_history = self.current_state["state_history"]

        return True

    def get_all_states(self) -> Dict[str, Dict]:
        """Extract all states from the state machine"""
        states = {}

        # Get orchestrator states
        if "agents" in self.state_machine and "orchestrator" in self.state_machine["agents"]:
            orch_states = self.state_machine["agents"]["orchestrator"].get("states", {})
            for state_name, state_data in orch_states.items():
                states[state_name] = state_data

        # Also check for states in sequences
        if "state_sequences" in self.state_machine:
            for sequence in self.state_machine["state_sequences"].values():
                if isinstance(sequence, list):
                    for state in sequence:
                        if state not in states:
                            states[state] = {"type": "sequence", "description": "State from sequence"}

        return states

    def find_state_path(self, target_state: str) -> List[str]:
        """Find the path to reach the current state"""
        states = self.get_all_states()

        # Build a graph of state transitions
        graph = {}
        for state_name, state_data in states.items():
            transitions = state_data.get("valid_transitions", [])
            graph[state_name] = transitions

        # Try to find a path from INIT to target
        def find_path_recursive(current: str, target: str, visited: Set[str], path: List[str]) -> Optional[List[str]]:
            if current == target:
                return path + [current]

            if current in visited:
                return None

            visited.add(current)

            for next_state in graph.get(current, []):
                result = find_path_recursive(next_state, target, visited.copy(), path + [current])
                if result:
                    return result

            return None

        # Try to find path from INIT
        path = find_path_recursive("INIT", target_state, set(), [])
        if path:
            return path

        # If no path from INIT, just return the state itself
        return [target_state]

    def draw_state_box(self, state: str, is_current: bool = False, width: int = 40) -> List[str]:
        """Draw a state box"""
        lines = []

        # Adjust width to fit state name
        state_len = len(state) + 4
        width = max(width, state_len)

        if is_current:
            # Current state with highlighting
            top = f"{Colors.GREEN}╔{'═' * (width - 2)}╗{Colors.RESET}"
            middle = f"{Colors.GREEN}║{Colors.BOLD}{Colors.YELLOW} {state.center(width - 4)} {Colors.RESET}{Colors.GREEN}║{Colors.RESET}"
            bottom = f"{Colors.GREEN}╚{'═' * (width - 2)}╝{Colors.RESET}"
            marker = f"{Colors.RED}{Colors.BOLD}◄── YOU ARE HERE{Colors.RESET}"

            lines.append(top)
            lines.append(middle + "  " + marker)
            lines.append(bottom)
        else:
            # Normal state
            lines.append(f"[{state}]")

        return lines

    def draw_arrow(self, label: str = "") -> str:
        """Draw a downward arrow"""
        if label:
            return f"  │ {Colors.GRAY}{label}{Colors.RESET}"
        return "  │"

    def generate_visualization(self) -> str:
        """Generate the full ASCII visualization"""
        output = []

        # Header
        output.append("=" * 80)
        output.append(f"{Colors.BOLD}{Colors.CYAN}SOFTWARE FACTORY 2.0 - STATE MACHINE VISUALIZATION{Colors.RESET}")
        output.append("=" * 80)
        output.append("")

        # Current state info
        current = self.current_state.get("current_state", "UNKNOWN")
        previous = self.current_state.get("previous_state", "N/A")
        phase = self.current_state.get("current_phase", "N/A")
        wave = self.current_state.get("current_wave", "N/A")

        output.append(f"{Colors.BOLD}Current State:{Colors.RESET} {Colors.YELLOW}{current}{Colors.RESET}")
        output.append(f"{Colors.BOLD}Previous State:{Colors.RESET} {previous}")
        output.append(f"{Colors.BOLD}Phase:{Colors.RESET} {phase}, {Colors.BOLD}Wave:{Colors.RESET} {wave}")

        # Sub-state machine info if active
        sub_state = self.current_state.get("sub_state_machine", {})
        if sub_state.get("active"):
            output.append("")
            output.append(f"{Colors.MAGENTA}● SUB-STATE MACHINE ACTIVE:{Colors.RESET}")
            output.append(f"  Type: {sub_state.get('type', 'N/A')}")
            output.append(f"  State: {sub_state.get('current_state', 'N/A')}")
            output.append(f"  Return to: {sub_state.get('return_state', 'N/A')}")

        output.append("")
        output.append("-" * 80)
        output.append(f"{Colors.BOLD}State Flow:{Colors.RESET}")
        output.append("-" * 80)
        output.append("")

        # Get the path to current state
        path = self.find_state_path(current)
        states_dict = self.get_all_states()

        # Draw the state flow
        for i, state in enumerate(path):
            is_current = (state == current)

            # Draw the state
            state_lines = self.draw_state_box(state, is_current)
            for line in state_lines:
                output.append(line)

            # Draw arrow to next state (if not last)
            if i < len(path) - 1:
                output.append(self.draw_arrow())
                output.append("  ▼")

        # Show possible next states
        if current in states_dict:
            next_states = states_dict[current].get("valid_transitions", [])
            if next_states:
                output.append("")
                output.append(self.draw_arrow())
                output.append("  ▼")
                output.append("")
                output.append(f"{Colors.BOLD}Possible Next States:{Colors.RESET}")
                for next_state in next_states[:3]:  # Show max 3 next states
                    output.append(f"  → {Colors.CYAN}{next_state}{Colors.RESET}")
                if len(next_states) > 3:
                    output.append(f"  → {Colors.GRAY}... and {len(next_states) - 3} more{Colors.RESET}")

        # Show state type and description
        output.append("")
        output.append("-" * 80)
        output.append(f"{Colors.BOLD}Current State Details:{Colors.RESET}")
        output.append("-" * 80)

        if current in states_dict:
            state_info = states_dict[current]
            state_type = state_info.get("type", "unknown")
            description = state_info.get("description", "No description available")

            # Color code by type
            type_color = {
                "spawn": Colors.GREEN,
                "waiting": Colors.YELLOW,
                "action": Colors.CYAN,
                "terminal": Colors.RED,
                "initial": Colors.MAGENTA,
                "checkpoint": Colors.BLUE
            }.get(state_type, Colors.WHITE)

            output.append(f"Type: {type_color}{state_type.upper()}{Colors.RESET}")
            output.append(f"Description: {description}")

            # Show required actions
            actions = state_info.get("required_actions", [])
            if actions:
                output.append("")
                output.append(f"{Colors.BOLD}Required Actions:{Colors.RESET}")
                for action in actions:
                    output.append(f"  • {action}")

            # Show entry conditions
            entry_conditions = state_info.get("entry_conditions", [])
            if entry_conditions:
                output.append("")
                output.append(f"{Colors.BOLD}Entry Conditions:{Colors.RESET}")
                for condition in entry_conditions:
                    output.append(f"  ✓ {condition}")

        # Legend
        output.append("")
        output.append("-" * 80)
        output.append(f"{Colors.BOLD}Legend:{Colors.RESET}")
        output.append("-" * 80)
        output.append(f"[STATE]                    - Normal state")
        output.append(f"{Colors.GREEN}╔════════════╗{Colors.RESET}")
        output.append(f"{Colors.GREEN}║{Colors.YELLOW} CURRENT ST {Colors.GREEN}║{Colors.RESET} {Colors.RED}◄── YOU ARE HERE{Colors.RESET} - Current state")
        output.append(f"{Colors.GREEN}╚════════════╝{Colors.RESET}")
        output.append(f"──▶                        - State transition")
        output.append(f"{Colors.GREEN}GREEN{Colors.RESET}                      - Spawn states")
        output.append(f"{Colors.YELLOW}YELLOW{Colors.RESET}                     - Waiting states")
        output.append(f"{Colors.CYAN}CYAN{Colors.RESET}                       - Action states")
        output.append(f"{Colors.RED}RED{Colors.RESET}                        - Terminal states")
        output.append(f"{Colors.MAGENTA}MAGENTA{Colors.RESET}                    - Initial states")
        output.append(f"{Colors.BLUE}BLUE{Colors.RESET}                       - Checkpoint states")

        # Footer
        output.append("")
        output.append("=" * 80)
        output.append(f"Generated: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
        output.append("=" * 80)

        return "\n".join(output)

    def generate_compact_view(self) -> str:
        """Generate a compact view for quick reference"""
        output = []

        current = self.current_state.get("current_state", "UNKNOWN")
        states_dict = self.get_all_states()

        output.append(f"{Colors.BOLD}{Colors.CYAN}╔{'═' * 50}╗{Colors.RESET}")
        output.append(f"{Colors.BOLD}{Colors.CYAN}║{'SOFTWARE FACTORY STATE MACHINE'.center(50)}║{Colors.RESET}")
        output.append(f"{Colors.BOLD}{Colors.CYAN}╚{'═' * 50}╝{Colors.RESET}")
        output.append("")

        # Show recent history
        if self.state_history:
            output.append(f"{Colors.BOLD}Recent Path:{Colors.RESET}")
            recent = self.state_history[-5:] if len(self.state_history) > 5 else self.state_history
            for state in recent:
                if isinstance(state, dict):
                    state_name = state.get("state", "?")
                else:
                    state_name = state
                output.append(f"  {state_name}")
            output.append("  ▼")

        # Current state
        output.append(f"{Colors.GREEN}{Colors.BOLD}[{current}] ◄── CURRENT{Colors.RESET}")

        # Next states
        if current in states_dict:
            next_states = states_dict[current].get("valid_transitions", [])
            if next_states:
                output.append("  ▼")
                output.append(f"{Colors.BOLD}Next:{Colors.RESET}")
                for next_state in next_states[:3]:
                    output.append(f"  → {next_state}")

        return "\n".join(output)

    def run(self, compact: bool = False):
        """Main execution"""
        if not self.load_data():
            return 1

        if compact:
            print(self.generate_compact_view())
        else:
            print(self.generate_visualization())

        return 0

def main():
    """Main entry point"""
    import argparse

    parser = argparse.ArgumentParser(description="Visualize Software Factory State Machine")
    parser.add_argument("-c", "--compact", action="store_true",
                       help="Show compact view")
    parser.add_argument("--no-color", action="store_true",
                       help="Disable colored output")

    args = parser.parse_args()

    if args.no_color:
        # Disable colors
        for attr in dir(Colors):
            if not attr.startswith("__"):
                setattr(Colors, attr, "")

    visualizer = StateVisualizer()
    return visualizer.run(compact=args.compact)

if __name__ == "__main__":
    sys.exit(main())