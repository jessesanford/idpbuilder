#!/usr/bin/env python3
"""
Interactive State Machine Viewer for Software Factory 3.0
Two views: Full zoomed-out view and focused current state view
Press 'v' to switch views, 'q' to quit
"""

import json
import os
import sys
import termios
import tty
import shutil
from pathlib import Path
from datetime import datetime
from typing import Dict, List, Optional, Tuple

# ANSI color codes
class Colors:
    RESET = '\033[0m'
    BOLD = '\033[1m'
    DIM = '\033[2m'

    # Foreground colors
    BLACK = '\033[30m'
    RED = '\033[31m'
    GREEN = '\033[32m'
    YELLOW = '\033[33m'
    BLUE = '\033[34m'
    MAGENTA = '\033[35m'
    CYAN = '\033[36m'
    WHITE = '\033[37m'

    # Background colors
    BG_BLACK = '\033[40m'
    BG_RED = '\033[41m'
    BG_GREEN = '\033[42m'
    BG_YELLOW = '\033[43m'
    BG_BLUE = '\033[44m'
    BG_MAGENTA = '\033[45m'
    BG_CYAN = '\033[46m'
    BG_WHITE = '\033[47m'

    # Bright colors
    BRIGHT_RED = '\033[91m'
    BRIGHT_GREEN = '\033[92m'
    BRIGHT_YELLOW = '\033[93m'
    BRIGHT_BLUE = '\033[94m'
    BRIGHT_MAGENTA = '\033[95m'
    BRIGHT_CYAN = '\033[96m'

class StateMachineViewer:
    def __init__(self, manual_width=None):
        self.current_view = 'journey'  # 'journey', 'focused', or 'overview'
        self.orchestrator_state = None
        self.state_machine = None
        self.current_state = None
        self.previous_state = None
        self.all_states = []
        self.transitions = {}
        self.scroll_offset = 0  # For scrolling in full view
        self.max_lines = 40  # Max lines to show at once
        self.visited_states = set()  # States we've actually been to
        self.state_journey = []  # Ordered journey through states
        self.sub_state_history = []  # Sub-state machine executions

        # State name mapping for known inconsistencies
        self.state_name_mapping = {
            'CREATE_NEXT_SPLIT_INFRASTRUCTURE': 'CREATE_NEXT_INFRASTRUCTURE',
            # Add more mappings as discovered
        }
        self.unmapped_states = set()  # Track states not in state machine

        # Terminal dimensions
        self.manual_width = manual_width  # Manual override if provided
        self.term_width = 80  # Default, will be updated
        self.term_height = 40  # Default, will be updated
        self.update_terminal_size()

    def update_terminal_size(self):
        """Update terminal dimensions with manual override support"""
        # Priority order:
        # 1. Manual width set via constructor (command-line argument)
        # 2. COLUMNS environment variable
        # 3. Auto-detected terminal size
        # 4. Default fallback

        if self.manual_width is not None:
            # Manual override takes highest priority
            self.term_width = self.manual_width
            # Still try to get height automatically
            try:
                _, self.term_height = shutil.get_terminal_size((80, 40))
            except:
                self.term_height = 40
        else:
            # Check for COLUMNS environment variable
            env_columns = os.environ.get('COLUMNS')
            if env_columns and env_columns.isdigit():
                self.term_width = int(env_columns)
                # Still try to get height automatically
                try:
                    _, self.term_height = shutil.get_terminal_size((80, 40))
                except:
                    self.term_height = 40
            else:
                # Auto-detect both width and height
                try:
                    self.term_width, self.term_height = shutil.get_terminal_size((80, 40))
                except:
                    self.term_width, self.term_height = 80, 40

    def center_text(self, text: str, width: int = None) -> str:
        """Center text within the given width (defaults to terminal width)"""
        import re
        width = width or self.term_width

        # Remove ANSI escape sequences to calculate actual visible length
        ansi_escape = re.compile(r'\x1B(?:[@-Z\\-_]|\[[0-?]*[ -/]*[@-~])')
        visible_text = ansi_escape.sub('', text)
        visible_len = len(visible_text)

        # Calculate padding needed
        if visible_len >= width:
            return text

        padding = (width - visible_len) // 2
        remainder = (width - visible_len) % 2

        # Return centered text with proper padding
        return ' ' * padding + text + ' ' * (padding + remainder)

    def draw_line(self, char: str = '─', width: int = None) -> str:
        """Draw a horizontal line with the given character"""
        width = width or self.term_width
        return char * width

    # --- Lightweight programmatic API for Textual wrapper ---
    def set_view(self, view: str) -> None:
        if view in ('journey', 'focused', 'overview'):
            self.current_view = view

    def cycle_view(self) -> None:
        if self.current_view == 'journey':
            self.current_view = 'focused'
        elif self.current_view == 'focused':
            self.current_view = 'overview'
        else:
            self.current_view = 'journey'

    def get_lines(self, view: Optional[str] = None) -> List[str]:
        """Return lines for the requested view (defaults to current)."""
        v = view or self.current_view
        if v == 'focused':
            # Focused view returns a single string; split to lines
            return (self.render_focused_view() or '').splitlines()
        if v == 'overview':
            return self.render_overview_view()
        # default journey
        return self.render_journey_view()

    def normalize_state_name(self, state: str) -> str:
        """Normalize state names to handle known inconsistencies"""
        if state in self.state_name_mapping:
            mapped = self.state_name_mapping[state]
            if state not in self.unmapped_states:
                self.unmapped_states.add(state)
                print(f"{Colors.YELLOW}Warning: Mapped state '{state}' to '{mapped}'{Colors.RESET}")
            return mapped
        return state

    def handle_key(self, key: str) -> None:
        """Update internal view/scroll based on a simple key string.
        Keys supported: 'v', 'UP', 'DOWN', 'PAGEUP', 'PAGEDOWN', 'HOME', 'END', 'r'.
        """
        if not key:
            return
        k = key.upper()
        if k == 'V':
            self.cycle_view()
            self.scroll_offset = 0
            return
        if k == 'R':
            # reload data best-effort
            try:
                self.load_data()
            except Exception:
                pass
            return
        if k in ('UP', 'DOWN', 'PAGEUP', 'PAGEDOWN', 'HOME', 'END'):
            if self.current_view in ('journey', 'overview'):
                if k == 'UP':
                    self.scroll_offset = max(0, self.scroll_offset - 1)
                elif k == 'DOWN':
                    self.scroll_offset += 1
                elif k == 'PAGEUP':
                    self.scroll_offset = max(0, self.scroll_offset - max(1, self.max_lines - 3))
                elif k == 'PAGEDOWN':
                    self.scroll_offset += max(1, self.max_lines - 3)
                elif k == 'HOME':
                    self.scroll_offset = 0
                elif k == 'END':
                    # end: move down a big chunk; wrapper will clamp
                    self.scroll_offset += 10_000
            return

    def load_data(self):
        """Load orchestrator state and state machine data"""
        # Find the orchestrator state file
        state_file = None
        for path in [
            Path.cwd() / "orchestrator-state-v3.json",
            Path.home() / "software-factory-template" / "orchestrator-state-v3.json",
            Path("/home/vscode/software-factory-template/orchestrator-state-v3.json"),
            Path("/home/vscode/workspaces") / "idpbuilder-push" / "orchestrator-state-v3.json"
        ]:
            if path.exists():
                state_file = path
                break

        if not state_file:
            print(f"{Colors.RED}Error: Could not find orchestrator-state-v3.json{Colors.RESET}")
            return False

        # Find the state machine file
        machine_file = None
        for path in [
            Path.cwd() / "state-machines" / "software-factory-3.0-state-machine.json",
            Path.home() / "software-factory-template" / "state-machines" / "software-factory-3.0-state-machine.json",
            Path("/home/vscode/software-factory-template/state-machines/software-factory-3.0-state-machine.json"),
            Path("/home/vscode/workspaces") / "idpbuilder-push" / "state-machines" / "software-factory-3.0-state-machine.json"
        ]:
            if path.exists():
                machine_file = path
                break

        if not machine_file:
            print(f"{Colors.RED}Error: Could not find software-factory-3.0-state-machine.json{Colors.RESET}")
            return False

        # Load the files
        with open(state_file, 'r') as f:
            self.orchestrator_state = json.load(f)

        # Load fix-cascade-state.json if it exists (SF 3.0 separate file structure)
        self.fix_cascade_state = None
        cascade_file = None
        for path in [
            Path.cwd() / "fix-cascade-state.json",
            state_file.parent / "fix-cascade-state.json"
        ]:
            if path.exists():
                cascade_file = path
                break

        if cascade_file:
            try:
                with open(cascade_file, 'r') as f:
                    self.fix_cascade_state = json.load(f)
            except Exception as e:
                # Non-fatal - cascade state is optional
                pass

        # Load integration-containers.json if it exists (for iteration tracking)
        self.integration_containers = None
        integration_file = None
        for path in [
            Path.cwd() / "integration-containers.json",
            state_file.parent / "integration-containers.json"
        ]:
            if path.exists():
                integration_file = path
                break

        if integration_file:
            try:
                with open(integration_file, 'r') as f:
                    self.integration_containers = json.load(f)
            except Exception as e:
                # Non-fatal - integration containers is optional
                pass

        with open(machine_file, 'r') as f:
            self.state_machine = json.load(f)

        # Extract current state info and normalize names
        # SF 3.0: nested under state_machine, SF 2.0: at top level
        if 'state_machine' in self.orchestrator_state:
            # SF 3.0 format
            raw_current = self.orchestrator_state['state_machine'].get('current_state', 'UNKNOWN')
            raw_previous = self.orchestrator_state['state_machine'].get('previous_state', 'UNKNOWN')
        else:
            # Legacy SF 2.0 format
            raw_current = self.orchestrator_state.get('current_state', 'UNKNOWN')
            raw_previous = self.orchestrator_state.get('previous_state', 'UNKNOWN')

        # Normalize state names to handle known inconsistencies
        self.current_state = self.normalize_state_name(raw_current)
        self.previous_state = self.normalize_state_name(raw_previous)

        # Extract all states and transitions
        if 'states' in self.state_machine:
            self.all_states = self.state_machine['states']

        # Get transitions for orchestrator
        if 'transition_matrix' in self.state_machine:
            if 'orchestrator' in self.state_machine['transition_matrix']:
                self.transitions = self.state_machine['transition_matrix']['orchestrator']

        # Build journey from orchestrator state history
        self.build_state_journey()

        return True

    def get_fix_cascade_count(self) -> int:
        """Get the count of fix cascades from fix-cascade-state.json"""
        if self.fix_cascade_state:
            # Try to get from statistics first (SF 3.0 structure)
            if 'statistics' in self.fix_cascade_state:
                stats = self.fix_cascade_state['statistics']
                # Could be total, active, or completed depending on what user wants to see
                # Return total cascades (both active + completed)
                return stats.get('total_cascades', 0)
            # Fallback: check if there's a cascade object (active cascade)
            elif 'cascade' in self.fix_cascade_state:
                # There's an active cascade
                return 1

        # Legacy fallback: check orchestrator state for fix_cascade_summary
        if 'fix_cascade_summary' in self.orchestrator_state:
            return self.orchestrator_state['fix_cascade_summary'].get('completed_fix_cascades', 0)

        return 0

    def get_iteration_container_counts(self) -> dict:
        """Get counts of iteration containers from integration-containers.json"""
        counts = {
            'total_containers': 0,
            'containers_with_iterations': 0,
            'total_iterations': 0,
            'max_iteration': 0
        }

        if not self.integration_containers:
            return counts

        # Check wave integrations
        for wave in self.integration_containers.get('wave_integrations', []):
            counts['total_containers'] += 1
            iteration = wave.get('iteration', 0)
            if iteration > 0:
                counts['containers_with_iterations'] += 1
                counts['total_iterations'] += iteration
                counts['max_iteration'] = max(counts['max_iteration'], iteration)

        # Check phase integrations
        for phase in self.integration_containers.get('phase_integrations', []):
            counts['total_containers'] += 1
            iteration = phase.get('iteration', 0)
            if iteration > 0:
                counts['containers_with_iterations'] += 1
                counts['total_iterations'] += iteration
                counts['max_iteration'] = max(counts['max_iteration'], iteration)

        # Check project integration if it has iteration tracking
        project = self.integration_containers.get('project_integration', {})
        if 'iteration' in project:
            counts['total_containers'] += 1
            iteration = project.get('iteration', 0)
            if iteration > 0:
                counts['containers_with_iterations'] += 1
                counts['total_iterations'] += iteration
                counts['max_iteration'] = max(counts['max_iteration'], iteration)

        return counts

    def build_state_journey(self):
        """Build the journey of states we've visited"""
        self.visited_states = set()
        self.state_journey = []

        # Extract sub-state history
        if 'sub_state_history' in self.orchestrator_state:
            self.sub_state_history = self.orchestrator_state['sub_state_history']

        # Always start with INIT if it exists
        if 'INIT' in self.all_states:
            self.state_journey.append('INIT')
            self.visited_states.add('INIT')

        # 1. Check for explicit state history (if exists)
        if 'state_history' in self.orchestrator_state:
            for entry in self.orchestrator_state['state_history']:
                if isinstance(entry, str):
                    # Simple string state
                    self.visited_states.add(entry)
                    if entry not in self.state_journey:
                        self.state_journey.append(entry)
                elif isinstance(entry, dict) and 'state' in entry:
                    # Dict with state field
                    state = entry['state']
                    self.visited_states.add(state)
                    if state not in self.state_journey:
                        self.state_journey.append(state)

        # 2. Look at sub-state machine history (FIX_CASCADE, etc.)
        if 'sub_state_history' in self.orchestrator_state:
            for sub_state in self.orchestrator_state['sub_state_history']:
                if isinstance(sub_state, dict):
                    if sub_state.get('type') == 'FIX_CASCADE':
                        # Add actual fix cascade states from the history
                        if 'states' in sub_state:
                            for state in sub_state['states']:
                                if state.startswith('FIX_CASCADE'):
                                    self.visited_states.add(state)
                        else:
                            # Default fix cascade states if no specific ones listed
                            self.visited_states.update([
                                'FIX_CASCADE_INIT',
                                'FIX_CASCADE_ANALYSIS',
                                'FIX_CASCADE_PLANNING',
                                'FIX_CASCADE_EXECUTION',
                                'FIX_CASCADE_VALIDATION'
                            ])
                        # Mark in journey - but place it at the right position
                        # Look for where this fix cascade occurred in the timeline
                        if '[FIX_CASCADE_LOOP]' not in self.state_journey:
                            # Insert after INIT if we have INIT, otherwise append
                            if 'INIT' in self.state_journey:
                                init_idx = self.state_journey.index('INIT')
                                self.state_journey.insert(init_idx + 1, '[FIX_CASCADE_LOOP]')
                            else:
                                self.state_journey.append('[FIX_CASCADE_LOOP]')

        # 3. Check fix cascade summary
        if 'fix_cascade_summary' in self.orchestrator_state:
            summary = self.orchestrator_state['fix_cascade_summary']
            if summary.get('completed_fix_cascades', 0) > 0:
                self.visited_states.add('FIX_CASCADE')
                self.visited_states.add('CASCADE_REINTEGRATION')

        # 4. Look at completed efforts and their states
        if 'efforts_completed' in self.orchestrator_state:
            if self.orchestrator_state['efforts_completed']:
                # These are states we must have gone through
                implementation_states = [
                    'PLANNING',
                    'CREATE_NEXT_INFRASTRUCTURE',
                    'SPAWN_SW_ENGINEERS',
                    'MONITORING',
                    'EFFORT_COMPLETE',
                    'WAVE_COMPLETE'
                ]
                for state in implementation_states:
                    # Add to visited states for counting
                    self.visited_states.add(state)
                    # ALSO add to journey so they show up in the display!
                    if state not in self.state_journey:
                        self.state_journey.append(state)

        # 5. Check integration branches (suggests integration states visited)
        if 'integration_branches' in self.orchestrator_state:
            if self.orchestrator_state['integration_branches']:
                integration_states = ['INTEGRATE_WAVE_EFFORTS', 'SPAWN_INTEGRATION_AGENT']
                for state in integration_states:
                    self.visited_states.add(state)
                    if state not in self.state_journey:
                        self.state_journey.append(state)

        # 6. Check for violations and fixes
        if 'violations' in self.orchestrator_state:
            if self.orchestrator_state['violations']:
                self.visited_states.update([
                    'SPAWN_CODE_REVIEWER_FIX_PLAN',
                    'WAITING_FOR_FIX_PLANS',
                    'SPAWN_SW_ENGINEERS',
                    'MONITORING_EFFORT_FIXES'
                ])

        # 7. Add previous state
        if self.previous_state and self.previous_state != 'UNKNOWN':
            self.visited_states.add(self.previous_state)
            if self.previous_state not in self.state_journey:
                self.state_journey.append(self.previous_state)

        # 8. Add current state
        if self.current_state and self.current_state != 'UNKNOWN':
            self.visited_states.add(self.current_state)
            if self.current_state not in self.state_journey:
                self.state_journey.append(self.current_state)

    def get_state_color(self, state: str) -> str:
        """Get color for a state based on its status"""
        if state == self.current_state:
            return Colors.BRIGHT_YELLOW + Colors.BOLD
        elif state in self.visited_states:
            return Colors.GREEN  # Visited states in green
        elif self.is_next_state(state):
            return Colors.BRIGHT_CYAN
        else:
            return Colors.DIM + Colors.WHITE

    def is_next_state(self, state: str) -> bool:
        """Check if a state is a valid next state from current"""
        # Normalize both the state and current state for comparison
        normalized_current = self.normalize_state_name(self.current_state)
        normalized_state = self.normalize_state_name(state)

        if normalized_current in self.transitions:
            next_states = self.transitions[normalized_current]
            if isinstance(next_states, list):
                return normalized_state in next_states
            elif isinstance(next_states, str):
                return normalized_state == next_states
        return False

    def build_graph_structure(self):
        """Build a graph structure from the transitions"""
        graph = {}

        # Build adjacency list
        for from_state, to_states in self.transitions.items():
            if from_state not in graph:
                graph[from_state] = []

            if isinstance(to_states, list):
                graph[from_state].extend(to_states)
            elif isinstance(to_states, str):
                graph[from_state].append(to_states)

        return graph

    def render_state_node(self, state: str, width: int = 25) -> List[str]:
        """Render a single state as a box node"""
        color = self.get_state_color(state)

        # Truncate state name if too long
        display_name = state[:width-4] + "..." if len(state) > width-1 else state

        if state == self.current_state:
            # Special rendering for current state
            return [
                f"{color}╔{'═' * width}╗{Colors.RESET}",
                f"{color}║ ▶ {display_name:^{width-4}} ◀ ║{Colors.RESET}",
                f"{color}╚{'═' * width}╝{Colors.RESET}"
            ]
        else:
            # Regular state box
            return [
                f"{color}┌{'─' * width}┐{Colors.RESET}",
                f"{color}│ {display_name:^{width-2}} │{Colors.RESET}",
                f"{color}└{'─' * width}┘{Colors.RESET}"
            ]

    def render_journey_view(self):
        """Render the journey map view - where we've been"""
        self.update_terminal_size()  # Refresh terminal size
        lines = []

        # Header with proper spacing
        header_line = self.draw_line('═')
        lines.append(f"{Colors.BRIGHT_MAGENTA}{header_line}{Colors.RESET}")
        header_text = 'SOFTWARE FACTORY 3.0 - JOURNEY MAP'
        centered_header = self.center_text(header_text)
        lines.append(f"{Colors.BRIGHT_MAGENTA}{centered_header}{Colors.RESET}")
        lines.append(f"{Colors.BRIGHT_MAGENTA}{header_line}{Colors.RESET}")
        lines.append("")  # Add extra spacing after header
        lines.append("")  # Additional spacing to prevent overlap

        # Current state info - handle SF 3.0 project_progression structure
        if 'project_progression' in self.orchestrator_state:
            # SF 3.0 format
            prog = self.orchestrator_state['project_progression']
            phase_info = prog.get('current_phase', {})
            wave_info = prog.get('current_wave', {})

            phase = f"Phase {phase_info.get('phase_number', '?')}" if phase_info else 'N/A'
            wave = f"Wave {wave_info.get('wave_number', '?')}" if wave_info else 'N/A'
        else:
            # Legacy format
            phase = self.orchestrator_state.get('current_phase', 'N/A')
            wave = self.orchestrator_state.get('current_wave', 'N/A')

        # Get fix cascade and iteration metrics
        fix_cascades = self.get_fix_cascade_count()
        iteration_counts = self.get_iteration_container_counts()

        info_line1 = f"Phase: {Colors.BRIGHT_GREEN}{phase}{Colors.RESET}  "
        info_line1 += f"Wave: {Colors.BRIGHT_GREEN}{wave}{Colors.RESET}  "
        info_line1 += f"Current: {Colors.BRIGHT_YELLOW}{self.current_state}{Colors.RESET}"

        info_line2 = f"States Visited: {Colors.GREEN}{len(self.visited_states)}{Colors.RESET}  "
        info_line2 += f"Fix Cascades: {Colors.BRIGHT_MAGENTA}{fix_cascades}{Colors.RESET}  "

        # Add iteration container info if any exist
        if iteration_counts['total_containers'] > 0:
            info_line2 += f"Iterations: {Colors.BRIGHT_CYAN}{iteration_counts['containers_with_iterations']}/{iteration_counts['total_containers']} containers{Colors.RESET}"

        # Center the info lines if there's extra space
        if self.term_width > 100:
            lines.append(self.center_text(info_line1))
            lines.append(self.center_text(info_line2))
        else:
            lines.append(info_line1)
            lines.append(info_line2)

        lines.append("")  # Extra spacing after info lines
        lines.append(f"{Colors.DIM}{self.draw_line('─')}{Colors.RESET}")
        lines.append("")  # Spacing before journey title
        lines.append(f"{Colors.BOLD}Your Journey So Far:{Colors.RESET}")
        lines.append("")

        # Build the journey path
        if self.state_journey:
            # Calculate box width based on terminal width
            box_width = min(max(35, self.term_width // 3), self.term_width - 10)
            center_offset = max(0, (self.term_width - box_width - 2) // 2)
            arrow_offset = center_offset + box_width // 2

            for i, state in enumerate(self.state_journey):
                # Check if this is the fix cascade marker
                if state == '[FIX_CASCADE_LOOP]':
                    # Show fix cascade as a special sub-journey
                    cascade_width = min(25, box_width - 10)
                    cascade_offset = center_offset + 10
                    lines.append(f"{Colors.BRIGHT_MAGENTA}{' ' * cascade_offset}╔{'═' * cascade_width}╗{Colors.RESET}")
                    lines.append(f"{Colors.BRIGHT_MAGENTA}{' ' * cascade_offset}║{'FIX CASCADE ENTERED':^{cascade_width}}║{Colors.RESET}")
                    lines.append(f"{Colors.BRIGHT_MAGENTA}{' ' * cascade_offset}╚{'═' * cascade_width}╝{Colors.RESET}")
                    lines.append(f"{Colors.BRIGHT_MAGENTA}{' ' * (cascade_offset + cascade_width // 2)}│{Colors.RESET}")

                    # Show fix cascade states if we have them
                    fix_states = ['FIX_CASCADE_INIT', 'FIX_CASCADE_ANALYSIS',
                                  'FIX_CASCADE_PLANNING', 'FIX_CASCADE_EXECUTION',
                                  'FIX_CASCADE_VALIDATION']
                    for fix_state in fix_states:
                        if fix_state in self.visited_states:
                            lines.append(f"{Colors.MAGENTA}{' ' * cascade_offset}┌{'─' * cascade_width}┐{Colors.RESET}")
                            lines.append(f"{Colors.MAGENTA}{' ' * cascade_offset}│{fix_state:^{cascade_width}}│{Colors.RESET}")
                            lines.append(f"{Colors.MAGENTA}{' ' * cascade_offset}└{'─' * cascade_width}┘{Colors.RESET}")
                            lines.append(f"{Colors.MAGENTA}{' ' * (cascade_offset + cascade_width // 2)}│{Colors.RESET}")

                    lines.append(f"{Colors.BRIGHT_MAGENTA}{' ' * cascade_offset}╔{'═' * cascade_width}╗{Colors.RESET}")
                    lines.append(f"{Colors.BRIGHT_MAGENTA}{' ' * cascade_offset}║{'FIX CASCADE COMPLETE':^{cascade_width}}║{Colors.RESET}")
                    lines.append(f"{Colors.BRIGHT_MAGENTA}{' ' * cascade_offset}╚{'═' * cascade_width}╝{Colors.RESET}")

                    # Continue with arrow if not last
                    if i < len(self.state_journey) - 1:
                        lines.append(f"{Colors.GREEN}{' ' * arrow_offset}│{Colors.RESET}")
                        lines.append(f"{Colors.GREEN}{' ' * arrow_offset}▼{Colors.RESET}")

                    continue

                # Determine if this is a key state
                is_current = (state == self.current_state)
                is_previous = (state == self.previous_state)

                # Render the state node
                if is_current:
                    # Current state - emphasized
                    lines.append(f"{Colors.BRIGHT_YELLOW}{' ' * center_offset}╔{'═' * box_width}╗{Colors.RESET}")
                    you_are_here = "▶ YOU ARE HERE ◀"
                    padding = (box_width - len(you_are_here)) // 2
                    lines.append(f"{Colors.BRIGHT_YELLOW}{' ' * center_offset}║{Colors.BRIGHT_RED}{' ' * padding}{you_are_here}{' ' * (box_width - len(you_are_here) - padding)}{Colors.BRIGHT_YELLOW}║{Colors.RESET}")
                    lines.append(f"{Colors.BRIGHT_YELLOW}{' ' * center_offset}║{state:^{box_width}}║{Colors.RESET}")
                    lines.append(f"{Colors.BRIGHT_YELLOW}{' ' * center_offset}╚{'═' * box_width}╝{Colors.RESET}")
                else:
                    # Previously visited state
                    lines.append(f"{Colors.GREEN}{' ' * center_offset}┌{'─' * box_width}┐{Colors.RESET}")
                    lines.append(f"{Colors.GREEN}{' ' * center_offset}│{state:^{box_width}}│{Colors.RESET}")
                    lines.append(f"{Colors.GREEN}{' ' * center_offset}└{'─' * box_width}┘{Colors.RESET}")

                # Add arrow to next state (if not last)
                if i < len(self.state_journey) - 1:
                    lines.append(f"{Colors.GREEN}{' ' * arrow_offset}│{Colors.RESET}")
                    lines.append(f"{Colors.GREEN}{' ' * arrow_offset}▼{Colors.RESET}")

        else:
            # Fallback if no journey built - show current and previous
            box_width = min(max(35, self.term_width // 3), self.term_width - 10)
            center_offset = max(0, (self.term_width - box_width - 2) // 2)
            arrow_offset = center_offset + box_width // 2

            if self.previous_state and self.previous_state != 'UNKNOWN':
                lines.append(f"{Colors.GREEN}{' ' * center_offset}┌{'─' * box_width}┐{Colors.RESET}")
                lines.append(f"{Colors.GREEN}{' ' * center_offset}│{self.previous_state:^{box_width}}│{Colors.RESET}")
                lines.append(f"{Colors.GREEN}{' ' * center_offset}└{'─' * box_width}┘{Colors.RESET}")
                lines.append(f"{Colors.GREEN}{' ' * arrow_offset}│{Colors.RESET}")
                lines.append(f"{Colors.GREEN}{' ' * arrow_offset}▼{Colors.RESET}")

            lines.append(f"{Colors.BRIGHT_YELLOW}{' ' * center_offset}╔{'═' * box_width}╗{Colors.RESET}")
            you_are_here = "▶ YOU ARE HERE ◀"
            padding = (box_width - len(you_are_here)) // 2
            lines.append(f"{Colors.BRIGHT_YELLOW}{' ' * center_offset}║{Colors.BRIGHT_RED}{' ' * padding}{you_are_here}{' ' * (box_width - len(you_are_here) - padding)}{Colors.BRIGHT_YELLOW}║{Colors.RESET}")
            lines.append(f"{Colors.BRIGHT_YELLOW}{' ' * center_offset}║{self.current_state:^{box_width}}║{Colors.RESET}")
            lines.append(f"{Colors.BRIGHT_YELLOW}{' ' * center_offset}╚{'═' * box_width}╝{Colors.RESET}")

        # Show next possible states (limited)
        box_width = min(max(35, self.term_width // 3), self.term_width - 10)
        center_offset = max(0, (self.term_width - box_width - 2) // 2)
        arrow_offset = center_offset + box_width // 2

        lines.append("")
        lines.append(f"{Colors.BRIGHT_CYAN}{' ' * arrow_offset}│{Colors.RESET}")
        lines.append(f"{Colors.BRIGHT_CYAN}{' ' * arrow_offset}▼{Colors.RESET}")
        lines.append("")

        next_states = []
        normalized_current = self.normalize_state_name(self.current_state)
        if normalized_current in self.transitions:
            next_val = self.transitions[normalized_current]
            if isinstance(next_val, list):
                # Show more next states on wider terminals
                max_next_states = 3 if self.term_width < 120 else 5
                next_states = next_val[:max_next_states]
            elif isinstance(next_val, str):
                next_states = [next_val]

        if next_states:
            lines.append(f"{Colors.BRIGHT_CYAN}{self.center_text('Next Possible States:')}{Colors.RESET}")
            lines.append("")

            for i, state in enumerate(next_states):
                if i > 0:
                    lines.append("")  # Spacing between states
                lines.append(f"{Colors.BRIGHT_CYAN}{' ' * center_offset}┌{'─' * box_width}┐{Colors.RESET}")
                lines.append(f"{Colors.BRIGHT_CYAN}{' ' * center_offset}│{state:^{box_width}}│{Colors.RESET}")
                lines.append(f"{Colors.BRIGHT_CYAN}{' ' * center_offset}└{'─' * box_width}┘{Colors.RESET}")

            if len(self.transitions.get(normalized_current, [])) > len(next_states):
                remaining = len(self.transitions[normalized_current]) - len(next_states)
                lines.append(f"\n{Colors.DIM}{self.center_text(f'... and {remaining} more possible states')}{Colors.RESET}")
        else:
            lines.append(f"{Colors.DIM}{self.center_text('(No valid next states from current position)')}{Colors.RESET}")

        # Legend
        lines.append("")
        lines.append(f"{Colors.DIM}{self.draw_line('─')}{Colors.RESET}")
        lines.append(f"{Colors.BOLD}Legend:{Colors.RESET}")
        lines.append(f"  {Colors.BRIGHT_YELLOW}╔══════╗{Colors.RESET}")
        lines.append(f"  {Colors.BRIGHT_YELLOW}║ HERE ║{Colors.RESET} = Current State (You Are Here)")
        lines.append(f"  {Colors.BRIGHT_YELLOW}╚══════╝{Colors.RESET}")
        lines.append(f"  {Colors.GREEN}┌──────┐{Colors.RESET}")
        lines.append(f"  {Colors.GREEN}│      │{Colors.RESET} = Previously Visited State")
        lines.append(f"  {Colors.GREEN}└──────┘{Colors.RESET}")
        lines.append(f"  {Colors.BRIGHT_CYAN}┌──────┐{Colors.RESET}")
        lines.append(f"  {Colors.BRIGHT_CYAN}│      │{Colors.RESET} = Next Possible State")
        lines.append(f"  {Colors.BRIGHT_CYAN}└──────┘{Colors.RESET}")
        lines.append("")
        lines.append(f"{Colors.DIM}Press 'v' to cycle views (journey/focused/overview), 'q' to quit, ↑/↓ to scroll{Colors.RESET}")

        return lines  # Return as list for scrolling

    def render_overview_view(self):
        """Render complete state machine overview with all states"""
        self.update_terminal_size()  # Refresh terminal size
        lines = []

        # Header with proper spacing
        lines.append(f"{Colors.BRIGHT_BLUE}{self.draw_line('═')}{Colors.RESET}")
        lines.append(f"{Colors.BRIGHT_BLUE}{self.center_text('SOFTWARE FACTORY 3.0 - COMPLETE STATE MACHINE OVERVIEW')}{Colors.RESET}")
        lines.append(f"{Colors.BRIGHT_BLUE}{self.draw_line('═')}{Colors.RESET}")
        lines.append("")  # Add extra spacing after header
        lines.append("")  # Additional spacing to prevent overlap

        # Statistics
        total_states = len(self.all_states)
        visited_count = len(self.visited_states)
        coverage = int((visited_count / total_states * 100)) if total_states > 0 else 0

        stats_line = f"Total States: {Colors.BRIGHT_CYAN}{total_states}{Colors.RESET}  "
        stats_line += f"Visited: {Colors.GREEN}{visited_count}{Colors.RESET}  "
        stats_line += f"Coverage: {Colors.BRIGHT_YELLOW}{coverage}%{Colors.RESET}  "
        stats_line += f"Current: {Colors.BRIGHT_YELLOW}{self.current_state}{Colors.RESET}"

        if self.term_width > 100:
            lines.append(self.center_text(stats_line))
        else:
            lines.append(stats_line)

        # Show fix cascade and iteration info if present
        fix_cascades = self.get_fix_cascade_count()
        iteration_counts = self.get_iteration_container_counts()

        metrics_added = False
        if fix_cascades > 0 or iteration_counts['total_containers'] > 0:
            metrics_line = ""
            if fix_cascades > 0:
                metrics_line += f"Fix Cascades: {Colors.BRIGHT_MAGENTA}{fix_cascades}{Colors.RESET}  "
                metrics_added = True

            if iteration_counts['total_containers'] > 0:
                metrics_line += f"Iteration Containers: {Colors.BRIGHT_CYAN}{iteration_counts['containers_with_iterations']}/{iteration_counts['total_containers']}{Colors.RESET}  "
                metrics_line += f"Max Iteration: {Colors.BRIGHT_CYAN}{iteration_counts['max_iteration']}{Colors.RESET}"
                metrics_added = True

            if metrics_added:
                if self.term_width > 100:
                    lines.append(self.center_text(metrics_line))
                else:
                    lines.append(metrics_line)

        lines.append("")  # Extra spacing after statistics
        lines.append(f"{Colors.DIM}{self.draw_line('─')}{Colors.RESET}")
        lines.append("")  # Additional spacing before content

        # Group states by category for organization
        categories = {
            'INIT & PLANNING': [],
            'SPAWNING AGENTS': [],
            'WAITING STATES': [],
            'MONITORING': [],
            'INTEGRATE_WAVE_EFFORTS': [],
            'FIXES & CASCADE': [],
            'REVIEW & VALIDATION': [],
            'PHASE & PROJECT': [],
            'BUILD & DEPLOY': [],
            'OTHER': []
        }

        # Categorize all states
        for state in sorted(self.all_states):
            categorized = False

            if 'INIT' in state or 'PLANNING' in state:
                categories['INIT & PLANNING'].append(state)
            elif 'SPAWN' in state:
                categories['SPAWNING AGENTS'].append(state)
            elif 'WAITING' in state or 'WAIT' in state:
                categories['WAITING STATES'].append(state)
            elif 'MONITORING' in state or 'MONITOR' in state:
                categories['MONITORING'].append(state)
            elif 'INTEGRATE_WAVE_EFFORTS' in state or 'INTEGRATE' in state:
                categories['INTEGRATE_WAVE_EFFORTS'].append(state)
            elif 'FIX' in state or 'CASCADE' in state or 'BACKPORT' in state:
                categories['FIXES & CASCADE'].append(state)
            elif 'REVIEW' in state or 'VALIDATION' in state or 'TEST' in state:
                categories['REVIEW & VALIDATION'].append(state)
            elif 'PHASE' in state or 'PROJECT' in state:
                categories['PHASE & PROJECT'].append(state)
            elif 'BUILD' in state or 'DEPLOY' in state or 'PR' in state:
                categories['BUILD & DEPLOY'].append(state)
            else:
                categories['OTHER'].append(state)

        # Display categories
        for category, states in categories.items():
            if not states:
                continue

            lines.append(f"{Colors.BOLD}{Colors.BRIGHT_CYAN}{category} ({len(states)} states):{Colors.RESET}")
            lines.append(f"{Colors.DIM}{self.draw_line('─', min(50, self.term_width))}{Colors.RESET}")

            # Calculate number of columns based on terminal width
            if self.term_width < 80:
                num_cols = 1
                col_width = self.term_width - 4
            elif self.term_width < 120:
                num_cols = 2
                col_width = (self.term_width - 6) // 2
            elif self.term_width < 180:
                num_cols = 3
                col_width = (self.term_width - 8) // 3
            else:
                num_cols = 4
                col_width = (self.term_width - 10) // 4

            # Display states in dynamic columns
            for i in range(0, len(states), num_cols):
                line = ""
                for j in range(num_cols):
                    if i + j < len(states):
                        state = states[i + j]

                        # Color based on status
                        if state == self.current_state:
                            color = Colors.BRIGHT_YELLOW + Colors.BOLD
                            prefix = "▶ "
                            suffix = " ◀"
                        elif state in self.visited_states:
                            color = Colors.GREEN
                            prefix = "✓ "
                            suffix = ""
                        elif self.is_next_state(state):
                            color = Colors.BRIGHT_CYAN
                            prefix = "→ "
                            suffix = ""
                        else:
                            color = Colors.DIM
                            prefix = "  "
                            suffix = ""

                        # Truncate state name if needed
                        max_state_len = col_width - len(prefix) - len(suffix) - 2
                        display_state = state[:max_state_len-3] + "..." if len(state) > max_state_len else state
                        formatted = f"{color}{prefix}{display_state}{suffix}{Colors.RESET}"
                        line += f"  {formatted:<{col_width}}"

                lines.append(line.rstrip())
            lines.append("")

        # Show sub-state machines if any
        if self.sub_state_history:
            lines.append(f"{Colors.BRIGHT_MAGENTA}{self.draw_line('─')}{Colors.RESET}")
            lines.append(f"{Colors.BRIGHT_MAGENTA}Sub-State Machine History:{Colors.RESET}")
            for sub in self.sub_state_history:
                if isinstance(sub, dict):
                    type_name = sub.get('type', 'UNKNOWN')
                    result = sub.get('result', 'UNKNOWN')
                    lines.append(f"  • {type_name}: {result}")

        # Legend
        lines.append("")
        lines.append(f"{Colors.DIM}{self.draw_line('─')}{Colors.RESET}")
        lines.append(f"{Colors.BOLD}Legend:{Colors.RESET}")
        lines.append(f"  {Colors.BRIGHT_YELLOW}▶ STATE ◀{Colors.RESET} = Current State")
        lines.append(f"  {Colors.GREEN}✓ STATE{Colors.RESET} = Previously Visited")
        lines.append(f"  {Colors.BRIGHT_CYAN}→ STATE{Colors.RESET} = Possible Next State")
        lines.append(f"  {Colors.DIM}  STATE{Colors.RESET} = Unvisited State")
        lines.append("")
        lines.append(f"{Colors.DIM}Press 'v' to cycle views, 'q' to quit, ↑/↓ to scroll{Colors.RESET}")

        return lines  # Return as list for scrolling

    def render_focused_view(self):
        """Render the focused view showing n-1, current, n+1 states"""
        self.update_terminal_size()  # Refresh terminal size
        lines = []

        # Calculate box widths based on terminal size
        header_width = min(max(60, self.term_width - 20), self.term_width)
        prev_box_width = min(max(40, self.term_width // 2), self.term_width - 10)
        current_box_width = min(max(50, self.term_width * 2 // 3), self.term_width - 10)
        next_box_width = min(max(50, self.term_width * 2 // 3), self.term_width - 10)

        # Center offsets
        header_offset = max(0, (self.term_width - header_width) // 2)
        prev_offset = max(0, (self.term_width - prev_box_width - 2) // 2)
        current_offset = max(0, (self.term_width - current_box_width - 2) // 2)
        next_offset = max(0, (self.term_width - next_box_width - 2) // 2)

        # Header with proper spacing
        lines.append(f"{Colors.BRIGHT_MAGENTA}{' ' * header_offset}{'═' * header_width}{Colors.RESET}")
        header_text = 'SOFTWARE FACTORY 3.0 - FOCUSED VIEW'
        # Center text within the header box, not the full terminal width
        padding = (header_width - len(header_text)) // 2
        lines.append(f"{Colors.BRIGHT_MAGENTA}{' ' * header_offset}{' ' * padding}{header_text}{' ' * (header_width - len(header_text) - padding)}{Colors.RESET}")
        lines.append(f"{Colors.BRIGHT_MAGENTA}{' ' * header_offset}{'═' * header_width}{Colors.RESET}")
        lines.append("")  # Add extra spacing after header
        lines.append("")  # Additional spacing to prevent overlap

        # Previous State (n-1)
        lines.append(f"{Colors.GREEN}{' ' * prev_offset}┌{'─' * prev_box_width}┐{Colors.RESET}")
        lines.append(f"{Colors.GREEN}{' ' * prev_offset}│  {'Previous State (n-1):':<{prev_box_width-4}}  │{Colors.RESET}")
        lines.append(f"{Colors.GREEN}{' ' * prev_offset}│  {self.previous_state:<{prev_box_width-4}}  │{Colors.RESET}")
        lines.append(f"{Colors.GREEN}{' ' * prev_offset}└{'─' * prev_box_width}┘{Colors.RESET}")
        arrow_pos = prev_offset + prev_box_width // 2
        lines.append(f"{Colors.GREEN}{' ' * arrow_pos}│{Colors.RESET}")
        lines.append(f"{Colors.GREEN}{' ' * arrow_pos}▼{Colors.RESET}")
        lines.append("")

        # Current State with emphasis
        lines.append(f"{Colors.BRIGHT_YELLOW}{' ' * current_offset}╔{'═' * current_box_width}╗{Colors.RESET}")
        you_are_here = "YOU ARE HERE"
        padding = (current_box_width - len(you_are_here)) // 2
        lines.append(f"{Colors.BRIGHT_YELLOW}{' ' * current_offset}║{Colors.BRIGHT_RED}{' ' * padding}{you_are_here}{' ' * (current_box_width - len(you_are_here) - padding)}{Colors.BRIGHT_YELLOW}║{Colors.RESET}")
        lines.append(f"{Colors.BRIGHT_YELLOW}{' ' * current_offset}║{' ' * current_box_width}║{Colors.RESET}")
        lines.append(f"{Colors.BRIGHT_YELLOW}{' ' * current_offset}║  {'Current State:':<{current_box_width-4}}  ║{Colors.RESET}")
        lines.append(f"{Colors.BRIGHT_YELLOW}{' ' * current_offset}║  {Colors.BOLD}{self.current_state:^{current_box_width-4}}{Colors.RESET}{Colors.BRIGHT_YELLOW}  ║{Colors.RESET}")
        lines.append(f"{Colors.BRIGHT_YELLOW}{' ' * current_offset}║{' ' * current_box_width}║{Colors.RESET}")

        # Add metadata about current state
        phase = self.orchestrator_state.get('current_phase', 'N/A')
        wave = self.orchestrator_state.get('current_wave', 'N/A')
        meta_line = f"Phase: {phase:<10} Wave: {wave:<10}"
        padding_needed = current_box_width - len(meta_line) - 2
        lines.append(f"{Colors.BRIGHT_YELLOW}{' ' * current_offset}║  {meta_line}{' ' * padding_needed}║{Colors.RESET}")

        # Transition time if available
        transition_time = self.orchestrator_state.get('transition_time', '')
        if transition_time:
            try:
                dt = datetime.fromisoformat(transition_time.replace('Z', '+00:00'))
                time_str = dt.strftime('%Y-%m-%d %H:%M:%S')
                time_line = f"Since: {time_str}"
                padding_needed = current_box_width - len(time_line) - 2
                lines.append(f"{Colors.BRIGHT_YELLOW}{' ' * current_offset}║  {time_line}{' ' * padding_needed}║{Colors.RESET}")
            except:
                pass

        lines.append(f"{Colors.BRIGHT_YELLOW}{' ' * current_offset}╚{'═' * current_box_width}╝{Colors.RESET}")
        arrow_pos = current_offset + current_box_width // 2
        lines.append(f"{Colors.BRIGHT_CYAN}{' ' * arrow_pos}│{Colors.RESET}")
        lines.append(f"{Colors.BRIGHT_CYAN}{' ' * arrow_pos}▼{Colors.RESET}")
        lines.append("")

        # Next States (n+1)
        next_states = []
        normalized_current = self.normalize_state_name(self.current_state)
        if normalized_current in self.transitions:
            next_val = self.transitions[normalized_current]
            if isinstance(next_val, list):
                next_states = next_val
            elif isinstance(next_val, str):
                next_states = [next_val]

        if next_states:
            lines.append(f"{Colors.BRIGHT_CYAN}{' ' * next_offset}┌{'─' * next_box_width}┐{Colors.RESET}")
            lines.append(f"{Colors.BRIGHT_CYAN}{' ' * next_offset}│  {'Possible Next States (n+1):':<{next_box_width-4}}  │{Colors.RESET}")

            # Show more states on wider terminals
            max_states = 5 if self.term_width < 120 else 8
            for i, state in enumerate(next_states[:max_states]):
                state_line = f"{i+1}. {state}"
                padding_needed = next_box_width - len(state_line) - 4
                lines.append(f"{Colors.BRIGHT_CYAN}{' ' * next_offset}│  {state_line:<{next_box_width-4}}  │{Colors.RESET}")

            if len(next_states) > max_states:
                more_line = f"... and {len(next_states)-max_states} more"
                lines.append(f"{Colors.BRIGHT_CYAN}{' ' * next_offset}│  {more_line:<{next_box_width-4}}  │{Colors.RESET}")

            lines.append(f"{Colors.BRIGHT_CYAN}{' ' * next_offset}└{'─' * next_box_width}┘{Colors.RESET}")
        else:
            lines.append(f"{Colors.DIM}{self.center_text('No valid transitions from current state')}{Colors.RESET}")

        # Additional info
        lines.append("")
        lines.append(f"{Colors.DIM}{self.draw_line('─')}{Colors.RESET}")

        # Transition reason if available
        reason = self.orchestrator_state.get('transition_reason', '')
        if reason:
            lines.append(f"{Colors.BOLD}Last Transition Reason:{Colors.RESET}")
            # Word wrap the reason
            words = reason.split()
            current_line = []
            for word in words:
                if len(' '.join(current_line + [word])) > 55:
                    lines.append(f"  {' '.join(current_line)}")
                    current_line = [word]
                else:
                    current_line.append(word)
            if current_line:
                lines.append(f"  {' '.join(current_line)}")

        lines.append("")
        lines.append(f"{Colors.DIM}Press 'v' to switch to full view, 'q' to quit{Colors.RESET}")

        return '\n'.join(lines)

    def clear_screen(self):
        """Clear the terminal screen"""
        os.system('clear' if os.name == 'posix' else 'cls')

    def get_key(self):
        """Get a single keypress, handle arrow keys"""
        fd = sys.stdin.fileno()
        old_settings = termios.tcgetattr(fd)
        try:
            tty.setraw(sys.stdin.fileno())
            key = sys.stdin.read(1)
            # Check for escape sequences (arrow keys)
            if key == '\x1b':
                key += sys.stdin.read(2)
                if key == '\x1b[A':
                    return 'UP'
                elif key == '\x1b[B':
                    return 'DOWN'
                elif key == '\x1b[C':
                    return 'RIGHT'
                elif key == '\x1b[D':
                    return 'LEFT'
        finally:
            termios.tcsetattr(fd, termios.TCSADRAIN, old_settings)
        return key

    def run(self):
        """Main run loop"""
        if not self.load_data():
            return

        try:
            while True:
                self.clear_screen()
                self.update_terminal_size()  # Refresh terminal size on each iteration

                # Render current view
                if self.current_view in ['journey', 'overview']:
                    # Get the view lines
                    if self.current_view == 'journey':
                        all_lines = self.render_journey_view()
                    else:
                        all_lines = self.render_overview_view()

                    # Use the instance's terminal size
                    term_width, term_height = self.term_width, self.term_height

                    # Calculate visible window
                    visible_height = term_height - 3  # Leave room for status

                    # Adjust scroll offset
                    max_scroll = max(0, len(all_lines) - visible_height)
                    self.scroll_offset = min(self.scroll_offset, max_scroll)
                    self.scroll_offset = max(0, self.scroll_offset)

                    # Display visible portion
                    visible_lines = all_lines[self.scroll_offset:self.scroll_offset + visible_height]
                    for line in visible_lines:
                        print(line)

                    # Show scroll indicator if needed
                    if len(all_lines) > visible_height:
                        scroll_pct = int((self.scroll_offset / max_scroll * 100)) if max_scroll > 0 else 0
                        print(f"\n{Colors.DIM}[Line {self.scroll_offset + 1}-{min(self.scroll_offset + visible_height, len(all_lines))} of {len(all_lines)}] [{scroll_pct}%]{Colors.RESET}")
                else:
                    print(self.render_focused_view())

                # Get user input
                key = self.get_key()

                if key.lower() == 'q':
                    break
                elif key.lower() == 'v':
                    # Cycle through views: journey -> focused -> overview -> journey
                    if self.current_view == 'journey':
                        self.current_view = 'focused'
                    elif self.current_view == 'focused':
                        self.current_view = 'overview'
                    else:
                        self.current_view = 'journey'
                    self.scroll_offset = 0
                elif key == 'UP':
                    # Scroll up
                    if self.current_view in ['journey', 'overview']:
                        self.scroll_offset = max(0, self.scroll_offset - 1)
                elif key == 'DOWN':
                    # Scroll down
                    if self.current_view in ['journey', 'overview']:
                        self.scroll_offset += 1
                elif key.lower() == 'r':
                    # Reload data
                    self.load_data()

        except KeyboardInterrupt:
            pass
        finally:
            self.clear_screen()
            print(f"{Colors.GREEN}Thanks for using State Machine Viewer!{Colors.RESET}")

def main():
    import sys
    import argparse as _argparse
    p = _argparse.ArgumentParser(description='Software Factory 3.0 State Machine Viewer')
    p.add_argument('--non-interactive', action='store_true',
                   help='Non-interactive mode (just print and exit)')
    p.add_argument('--view', choices=['journey','focused','overview'], default='journey',
                   help='Initial view mode')
    p.add_argument('--width', '--columns', type=int, dest='width',
                   help='Manually set terminal width (overrides auto-detection and COLUMNS env var)')
    args, _unknown = p.parse_known_args()

    # Create viewer with manual width if provided
    viewer = StateMachineViewer(manual_width=args.width)

    if args.non_interactive or '--no-interactive' in sys.argv or not sys.stdout.isatty():
        if viewer.load_data():
            if args.view == 'focused':
                print(viewer.render_focused_view())
            elif args.view == 'overview':
                for line in viewer.render_overview_view():
                    print(line)
            else:
                for line in viewer.render_journey_view():
                    print(line)
        return
    # Interactive mode
    viewer.set_view(args.view)
    viewer.run()

if __name__ == "__main__":
    main()