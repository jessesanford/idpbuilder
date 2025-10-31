#!/usr/bin/env python3
"""
SOFTWARE FACTORY 3.0 - STATE MACHINE VISUALIZER
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
        Path("/workspaces/software-factory-3.0-template")
    ]:
        if (possible_root / "orchestrator-state-v3.json").exists() or \
           (possible_root / "orchestrator-state-demo.json").exists() or \
           (possible_root / ".state-backup").exists():
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

        # Try backup if main state file doesn't exist
        if not self.state_file.exists():
            backup_dir = self.project_root / ".state-backup"
            if backup_dir.exists():
                # Get most recent backup
                backups = sorted(backup_dir.glob("*/orchestrator-state-v3.json"))
                if backups:
                    self.state_file = backups[-1]

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

        # Extract state history if available (SF 3.0 format)
        if "state_machine" in self.current_state and "state_history" in self.current_state["state_machine"]:
            self.state_history = self.current_state["state_machine"]["state_history"]
        elif "state_history" in self.current_state:
            self.state_history = self.current_state["state_history"]

        return True

    def get_all_states(self) -> Dict[str, Dict]:
        """Extract all states from the SF 3.0 state machine"""
        states = {}

        # SF 3.0: states are at top level .states
        if "states" in self.state_machine:
            return self.state_machine["states"]

        # Legacy fallback for SF 2.0 structure
        if "agents" in self.state_machine and "orchestrator" in self.state_machine["agents"]:
            orch_states = self.state_machine["agents"]["orchestrator"].get("states", {})
            for state_name, state_data in orch_states.items():
                states[state_name] = state_data

        return states

    def get_current_state_info(self) -> Dict:
        """Extract current state info from SF 3.0 structure"""
        # SF 3.0 format: nested under state_machine
        if "state_machine" in self.current_state:
            return {
                "current_state": self.current_state["state_machine"].get("current_state", "UNKNOWN"),
                "previous_state": self.current_state["state_machine"].get("previous_state", "N/A"),
                "sub_state_machine": self.current_state["state_machine"].get("sub_state_machine", {})
            }

        # Legacy format
        return {
            "current_state": self.current_state.get("current_state", "UNKNOWN"),
            "previous_state": self.current_state.get("previous_state", "N/A"),
            "sub_state_machine": self.current_state.get("sub_state_machine", {})
        }

    def get_progression_info(self) -> Dict:
        """Extract project progression info (phase/wave/efforts)"""
        if "project_progression" not in self.current_state:
            return {}

        prog = self.current_state["project_progression"]
        return {
            "project": prog.get("current_project", {}),
            "phase": prog.get("current_phase", {}),
            "wave": prog.get("current_wave", {}),
            "iteration_tracking": prog.get("iteration_tracking", {})
        }

    def draw_progress_bar(self, current: int, maximum: int, width: int = 20) -> str:
        """Draw a progress bar with color coding"""
        if maximum == 0:
            return f"[{'░' * width}] 0/0"

        percentage = (current / maximum) * 100
        filled = int((current / maximum) * width)

        # Color code based on percentage
        if percentage < 60:
            color = Colors.GREEN
            fill_char = '█'
        elif percentage < 80:
            color = Colors.YELLOW
            fill_char = '▓'
        else:
            color = Colors.RED
            fill_char = '▓'

        bar = f"{color}{fill_char * filled}{Colors.GRAY}{'░' * (width - filled)}{Colors.RESET}"
        return f"[{bar}] {current}/{maximum}"

    def draw_iteration_warning(self, current: int, maximum: int, name: str) -> Optional[str]:
        """Generate warning if approaching iteration limit"""
        if maximum == 0:
            return None

        percentage = (current / maximum) * 100
        remaining = maximum - current

        if percentage >= 90:
            return f"{Colors.RED}🔴 CRITICAL:{Colors.RESET} {name} at {current}/{maximum} ({percentage:.0f}%) - near exhaustion! ({remaining} remaining)"
        elif percentage >= 80:
            return f"{Colors.YELLOW}⚠️  WARNING:{Colors.RESET} {name} at {current}/{maximum} ({percentage:.0f}%) - approaching limit ({remaining} remaining)"

        return None

    def find_state_path(self, target_state: str) -> List[str]:
        """Find the path to reach the current state"""
        states = self.get_all_states()

        # Build a graph of state transitions (SF 3.0: allowed_transitions)
        graph = {}
        for state_name, state_data in states.items():
            transitions = state_data.get("allowed_transitions", state_data.get("valid_transitions", []))
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
        output.append(f"{Colors.BOLD}{Colors.CYAN}SOFTWARE FACTORY 3.0 - STATE MACHINE VISUALIZATION{Colors.RESET}")
        output.append("=" * 80)
        output.append("")

        # Current state info
        state_info = self.get_current_state_info()
        current = state_info["current_state"]
        previous = state_info["previous_state"]

        output.append(f"{Colors.BOLD}Current State:{Colors.RESET} {Colors.YELLOW}{current}{Colors.RESET}")
        output.append(f"{Colors.BOLD}Previous State:{Colors.RESET} {previous}")

        # SF 3.0: Project progression info
        prog_info = self.get_progression_info()
        if prog_info:
            output.append("")
            output.append(f"{Colors.BOLD}{Colors.CYAN}Project Progression:{Colors.RESET}")

            # Project info
            if prog_info.get("project"):
                project = prog_info["project"]
                output.append(f"├─ Project: {project.get('name', 'N/A')}")
                if project.get("iteration") is not None:
                    bar = self.draw_progress_bar(project.get("iteration", 0), project.get("max_iterations", 10))
                    output.append(f"│  └─ Iterations: {bar}")

            # Phase info
            if prog_info.get("phase"):
                phase = prog_info["phase"]
                # Handle both list and int formats for waves_completed
                waves_completed = phase.get("waves_completed", [])
                waves_done = waves_completed if isinstance(waves_completed, int) else len(waves_completed)
                output.append(f"├─ Phase {phase.get('phase_number', '?')}: {phase.get('name', 'N/A')}")
                output.append(f"│  ├─ Status: {phase.get('status', 'UNKNOWN')}")
                output.append(f"│  ├─ Waves Completed: {waves_done}")
                if phase.get("iteration") is not None:
                    bar = self.draw_progress_bar(phase.get("iteration", 0), phase.get("max_iterations", 10))
                    output.append(f"│  └─ Iterations: {bar}")

            # Wave info
            if prog_info.get("wave"):
                wave = prog_info["wave"]
                # Handle both list and int formats for efforts_completed
                efforts_completed = wave.get("efforts_completed", [])
                efforts_done = efforts_completed if isinstance(efforts_completed, int) else len(efforts_completed)
                output.append(f"└─ Wave {wave.get('wave_number', '?')}: {wave.get('name', 'N/A')}")
                output.append(f"   ├─ Status: {wave.get('status', 'UNKNOWN')}")
                output.append(f"   ├─ Efforts Completed: {efforts_done}")
                if wave.get("iteration") is not None:
                    bar = self.draw_progress_bar(wave.get("iteration", 0), wave.get("max_iterations", 10))
                    output.append(f"   └─ Iterations: {bar}")

            # Container iteration tracking
            if prog_info.get("iteration_tracking"):
                it_track = prog_info["iteration_tracking"]
                output.append("")
                output.append(f"{Colors.BOLD}Active Container:{Colors.RESET} {it_track.get('active_container_level', 'N/A')} (ID: {it_track.get('active_container_id', 'N/A')})")

        # Iteration warnings
        warnings = []
        if prog_info:
            for container_name, container in [("Project", prog_info.get("project", {})),
                                              ("Phase", prog_info.get("phase", {})),
                                              ("Wave", prog_info.get("wave", {}))]:
                if container.get("iteration") is not None:
                    warning = self.draw_iteration_warning(
                        container.get("iteration", 0),
                        container.get("max_iterations", 10),
                        container_name
                    )
                    if warning:
                        warnings.append(warning)

        if warnings:
            output.append("")
            output.append(f"{Colors.BOLD}⚠️  ITERATION WARNINGS:{Colors.RESET}")
            for warning in warnings:
                output.append(f"  {warning}")

        # Sub-state machine info if active
        sub_state = state_info.get("sub_state_machine") or {}
        if sub_state and sub_state.get("active"):
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

        # Show possible next states (SF 3.0: allowed_transitions)
        if current in states_dict:
            next_states = states_dict[current].get("allowed_transitions", states_dict[current].get("valid_transitions", []))
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
            state_data = states_dict[current]
            description = state_data.get("description", "No description available")
            agent = state_data.get("agent", "unknown")
            checkpoint = state_data.get("checkpoint", False)
            iteration_level = state_data.get("iteration_level", "N/A")

            output.append(f"Agent: {Colors.CYAN}{agent}{Colors.RESET}")
            output.append(f"Description: {description}")
            output.append(f"Checkpoint: {Colors.GREEN if checkpoint else Colors.GRAY}{checkpoint}{Colors.RESET}")
            output.append(f"Iteration Level: {Colors.YELLOW}{iteration_level}{Colors.RESET}")

            # Show required actions
            actions = state_data.get("actions", [])
            if actions:
                output.append("")
                output.append(f"{Colors.BOLD}Required Actions:{Colors.RESET}")
                for action in actions:
                    output.append(f"  • {action}")

            # Show entry conditions
            requires = state_data.get("requires", {})
            conditions = requires.get("conditions", [])
            if conditions:
                output.append("")
                output.append(f"{Colors.BOLD}Entry Conditions:{Colors.RESET}")
                for condition in conditions:
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
        output.append(f"{Colors.GREEN}[███░░░░░░░]{Colors.RESET}            - Progress bar (< 60% = green)")
        output.append(f"{Colors.YELLOW}[▓▓▓▓▓▓░░░░]{Colors.RESET}            - Progress bar (60-80% = yellow)")
        output.append(f"{Colors.RED}[▓▓▓▓▓▓▓▓▓░]{Colors.RESET}            - Progress bar (> 80% = red)")

        # Footer
        output.append("")
        output.append("=" * 80)
        output.append(f"Generated: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
        output.append(f"State File: {self.state_file}")
        output.append("=" * 80)

        return "\n".join(output)

    def generate_compact_view(self) -> str:
        """Generate a compact view for quick reference"""
        output = []

        state_info = self.get_current_state_info()
        current = state_info["current_state"]
        states_dict = self.get_all_states()

        output.append(f"{Colors.BOLD}{Colors.CYAN}╔{'═' * 50}╗{Colors.RESET}")
        output.append(f"{Colors.BOLD}{Colors.CYAN}║{'SOFTWARE FACTORY 3.0 STATE MACHINE'.center(50)}║{Colors.RESET}")
        output.append(f"{Colors.BOLD}{Colors.CYAN}╚{'═' * 50}╝{Colors.RESET}")
        output.append("")

        # Show recent history
        if self.state_history:
            output.append(f"{Colors.BOLD}Recent Path:{Colors.RESET}")
            recent = self.state_history[-5:] if len(self.state_history) > 5 else self.state_history
            for state in recent:
                if isinstance(state, dict):
                    state_name = state.get("to_state", state.get("state", "?"))
                else:
                    state_name = state
                output.append(f"  {state_name}")
            output.append("  ▼")

        # Current state
        output.append(f"{Colors.GREEN}{Colors.BOLD}[{current}] ◄── CURRENT{Colors.RESET}")

        # Progression info (compact)
        prog_info = self.get_progression_info()
        if prog_info:
            if prog_info.get("phase"):
                phase = prog_info["phase"]
                output.append(f"  Phase {phase.get('phase_number', '?')}: {phase.get('name', 'N/A')}")
            if prog_info.get("wave"):
                wave = prog_info["wave"]
                output.append(f"  Wave {wave.get('wave_number', '?')}: {wave.get('name', 'N/A')}")

        # Next states
        if current in states_dict:
            next_states = states_dict[current].get("allowed_transitions", states_dict[current].get("valid_transitions", []))
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

    parser = argparse.ArgumentParser(description="Visualize Software Factory 3.0 State Machine")
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
