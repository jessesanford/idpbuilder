#!/usr/bin/env python3
"""
SOFTWARE FACTORY 2.0 - ENHANCED STATE FLOW VISUALIZER
Shows complex state transitions and parallel execution paths
"""

import json
import sys
from pathlib import Path
from typing import Dict, List, Set, Optional
from collections import defaultdict

class FlowVisualizer:
    def __init__(self):
        self.project_root = self._find_project_root()
        self.state_machine = {}
        self.current_state = {}

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
        # Try to load state machine
        state_machine_paths = [
            self.project_root / "state-machines" / "software-factory-3.0-state-machine.json",
            self.project_root / "software-factory-3.0-state-machine.json"
        ]

        for path in state_machine_paths:
            if path.exists():
                try:
                    with open(path, 'r') as f:
                        self.state_machine = json.load(f)
                        break
                except Exception as e:
                    print(f"Error loading {path}: {e}")

        # Load current state
        state_paths = [
            self.project_root / "orchestrator-state-v3.json",
            self.project_root / "orchestrator-state-demo.json"
        ]

        for path in state_paths:
            if path.exists():
                try:
                    with open(path, 'r') as f:
                        self.current_state = json.load(f)
                        break
                except Exception as e:
                    print(f"Error loading {path}: {e}")

        return bool(self.state_machine)

    def generate_wave_flow(self) -> str:
        """Generate wave execution flow diagram"""
        output = []

        output.append("WAVE EXECUTION FLOW")
        output.append("=" * 70)
        output.append("")
        output.append("┌─────────────┐")
        output.append("│  WAVE_START │")
        output.append("└──────┬──────┘")
        output.append("       │")
        output.append("       ▼")
        output.append("┌─────────────────────────────────┐")
        output.append("│ SPAWN_ARCHITECT_WAVE_PLANNING    │ (if needed)")
        output.append("└──────┬──────────────────────────┘")
        output.append("       │")
        output.append("       ▼")
        output.append("┌─────────────────────────────────┐")
        output.append("│ CREATE_NEXT_INFRASTRUCTURE      │")
        output.append("└──────┬──────────────────────────┘")
        output.append("       │")
        output.append("       ▼")
        output.append("┌─────────────────────────────────┐")
        output.append("│ ANALYZE_CODE_REVIEWER_PARALLEL   │")
        output.append("└──────┬──────────────────────────┘")
        output.append("       │")
        output.append("    ┌──┴──┐")
        output.append("    │ 1?  │ ─────► Single effort optimization (R356)")
        output.append("    └──┬──┘")
        output.append("       │ Multiple")
        output.append("       ▼")
        output.append("┌─────────────────────────────────────────┐")
        output.append("│ SPAWN_CODE_REVIEWERS_EFFORT_PLANNING     │")
        output.append("│ ┌─────────┐ ┌─────────┐ ┌─────────┐    │")
        output.append("│ │ CR-1    │ │ CR-2    │ │ CR-N    │    │ (Parallel)")
        output.append("│ └─────────┘ └─────────┘ └─────────┘    │")
        output.append("└──────┬──────────────────────────────────┘")
        output.append("       │")
        output.append("       ▼")
        output.append("┌─────────────────────────────────┐")
        output.append("│ WAITING_FOR_EFFORT_PLANS         │")
        output.append("└──────┬──────────────────────────┘")
        output.append("       │")
        output.append("       ▼")
        output.append("┌─────────────────────────────────┐")
        output.append("│ ANALYZE_IMPLEMENTATION_PARALLEL  │")
        output.append("└──────┬──────────────────────────┘")
        output.append("       │")
        output.append("       ▼")
        output.append("┌─────────────────────────────────────────┐")
        output.append("│ SPAWN_SW_ENGINEERS                             │")
        output.append("│ ┌─────────┐ ┌─────────┐ ┌─────────┐    │")
        output.append("│ │ SWE-1   │ │ SWE-2   │ │ SWE-N   │    │ (Parallel)")
        output.append("│ └─────────┘ └─────────┘ └─────────┘    │")
        output.append("└──────┬──────────────────────────────────┘")
        output.append("       │")
        output.append("       ▼")
        output.append("┌─────────────────────────────────┐")
        output.append("│ MONITORING                       │")
        output.append("└──────┬──────────────────────────┘")
        output.append("       │")
        output.append("       ▼")
        output.append("┌─────────────────────────────────┐")
        output.append("│ SPAWN_CODE_REVIEWERS_EFFORT_REVIEW  │")
        output.append("└──────┬──────────────────────────┘")
        output.append("       │")
        output.append("    ┌──┴──┐")
        output.append("    │Issues│ YES ──► FIX_CASCADE_INIT")
        output.append("    └──┬──┘")
        output.append("       │ NO")
        output.append("       ▼")
        output.append("┌─────────────────────────────────┐")
        output.append("│ WAVE_COMPLETE                    │")
        output.append("└─────────────────────────────────┘")

        return "\n".join(output)

    def generate_integration_flow(self) -> str:
        """Generate integration flow diagram"""
        output = []

        output.append("INTEGRATE_WAVE_EFFORTS FLOW")
        output.append("=" * 70)
        output.append("")
        output.append("                    WAVE INTEGRATE_WAVE_EFFORTS")
        output.append("┌─────────────────┐         ┌──────────────────────┐")
        output.append("│  WAVE_COMPLETE  │ ──────► │ CREATE_WAVE_INT_BR   │")
        output.append("└─────────────────┘         └──────────┬───────────┘")
        output.append("                                       │")
        output.append("                                       ▼")
        output.append("                            ┌──────────────────────┐")
        output.append("                            │ SPAWN_INTEGRATE_WAVE_EFFORTS     │")
        output.append("                            │      AGENT            │")
        output.append("                            └──────────┬───────────┘")
        output.append("                                       │")
        output.append("                                       ▼")
        output.append("                            ┌──────────────────────┐")
        output.append("                            │ WAITING_FOR_MERGE    │")
        output.append("                            └──────────┬───────────┘")
        output.append("                                       │")
        output.append("                                    ┌──┴──┐")
        output.append("                                    │Fixes│ YES ──► IMMEDIATE_BACKPORT")
        output.append("                                    └──┬──┘")
        output.append("                                       │ NO")
        output.append("                                       ▼")
        output.append("                            ┌──────────────────────┐")
        output.append("                            │ ARCHITECT_REVIEW_WAVE_ARCHITECTURE│")
        output.append("                            └──────────────────────┘")
        output.append("")
        output.append("                    PHASE INTEGRATE_WAVE_EFFORTS")
        output.append("┌─────────────────┐         ┌──────────────────────┐")
        output.append("│ COMPLETE_PHASE  │ ──────► │ CREATE_PHASE_INT_BR  │")
        output.append("└─────────────────┘         └──────────┬───────────┘")
        output.append("                                       │")
        output.append("                                       ▼")
        output.append("                            ┌──────────────────────┐")
        output.append("                            │ SPAWN_INTEGRATE_WAVE_EFFORTS     │")
        output.append("                            │    AGENT_PHASE        │")
        output.append("                            └──────────┬───────────┘")
        output.append("                                       │")
        output.append("                                       ▼")
        output.append("                            ┌──────────────────────┐")
        output.append("                            │ ARCHITECT_PHASE_REVIEW│")
        output.append("                            └──────────────────────┘")

        return "\n".join(output)

    def generate_fix_cascade_flow(self) -> str:
        """Generate fix cascade flow"""
        output = []

        output.append("FIX CASCADE FLOW")
        output.append("=" * 70)
        output.append("")
        output.append("┌─────────────────────────┐")
        output.append("│ Issues Detected         │")
        output.append("└────────┬────────────────┘")
        output.append("         │")
        output.append("         ▼")
        output.append("┌─────────────────────────┐")
        output.append("│ FIX_CASCADE_INIT        │")
        output.append("└────────┬────────────────┘")
        output.append("         │")
        output.append("         ▼")
        output.append("┌─────────────────────────┐")
        output.append("│ SPAWN_CODE_REVIEWER     │")
        output.append("│    FIX_PLANNING         │")
        output.append("└────────┬────────────────┘")
        output.append("         │")
        output.append("         ▼")
        output.append("┌─────────────────────────┐")
        output.append("│ WAITING_FOR_FIX_PLANS   │")
        output.append("└────────┬────────────────┘")
        output.append("         │")
        output.append("         ▼")
        output.append("┌─────────────────────────┐")
        output.append("│ CREATE_WAVE_FIX_PLAN    │ (checkpoint)")
        output.append("└────────┬────────────────┘")
        output.append("         │")
        output.append("         ▼")
        output.append("┌─────────────────────────┐")
        output.append("│ SPAWN_ENGINEERS_FOR_FIX │")
        output.append("└────────┬────────────────┘")
        output.append("         │")
        output.append("         ▼")
        output.append("┌─────────────────────────┐")
        output.append("│ MONITORING_EFFORT_FIXES │")
        output.append("└────────┬────────────────┘")
        output.append("         │")
        output.append("         ▼")
        output.append("┌─────────────────────────┐")
        output.append("│ SPAWN_CODE_REVIEWERS    │")
        output.append("│    FOR_FIX_REVIEW       │")
        output.append("└────────┬────────────────┘")
        output.append("         │")
        output.append("         ├─► Issues? ──► (Loop back)")
        output.append("         │")
        output.append("         ▼")
        output.append("┌─────────────────────────┐")
        output.append("│ FIX_CASCADE_COMPLETE    │")
        output.append("└─────────────────────────┘")

        return "\n".join(output)

    def generate_split_flow(self) -> str:
        """Generate split handling flow"""
        output = []

        output.append("SPLIT HANDLING FLOW")
        output.append("=" * 70)
        output.append("")
        output.append("┌─────────────────────────┐")
        output.append("│ Size > 700 lines        │")
        output.append("└────────┬────────────────┘")
        output.append("         │")
        output.append("         ▼")
        output.append("┌─────────────────────────┐")
        output.append("│ EFFORT_SPLIT_REQUIRED   │")
        output.append("└────────┬────────────────┘")
        output.append("         │")
        output.append("         ▼")
        output.append("┌─────────────────────────┐")
        output.append("│ CREATE_SPLIT_PLAN       │")
        output.append("└────────┬────────────────┘")
        output.append("         │")
        output.append("         ▼")
        output.append("┌─────────────────────────────────────┐")
        output.append("│ SEQUENTIAL SPLIT EXECUTION           │")
        output.append("│                                      │")
        output.append("│  Split-1 ──► Review ──► Issues? ──┐ │")
        output.append("│     ▼                             │ │")
        output.append("│  Split-2 ──► Review ──► Issues? ──┤ │")
        output.append("│     ▼                             │ │")
        output.append("│  Split-N ──► Review ──► Issues? ──┘ │")
        output.append("│                                      │")
        output.append("└──────────────┬───────────────────────┘")
        output.append("               │")
        output.append("               ▼")
        output.append("┌─────────────────────────┐")
        output.append("│ All Splits Complete     │")
        output.append("└─────────────────────────┘")

        return "\n".join(output)

    def generate_complete_flow(self) -> str:
        """Generate complete flow diagram"""
        output = []

        # Header
        output.append("=" * 80)
        output.append("SOFTWARE FACTORY 2.0 - COMPLETE STATE FLOW DIAGRAM")
        output.append("=" * 80)
        output.append("")

        # Current state marker
        current = self.current_state.get("current_state", "UNKNOWN")
        phase = self.current_state.get("current_phase", "N/A")
        wave = self.current_state.get("current_wave", "N/A")

        output.append(f"CURRENT POSITION: {current}")
        output.append(f"Phase: {phase}, Wave: {wave}")
        output.append("")
        output.append("-" * 80)
        output.append("")

        # Show different flow sections
        output.append(self.generate_wave_flow())
        output.append("")
        output.append("-" * 80)
        output.append("")
        output.append(self.generate_integration_flow())
        output.append("")
        output.append("-" * 80)
        output.append("")
        output.append(self.generate_fix_cascade_flow())
        output.append("")
        output.append("-" * 80)
        output.append("")
        output.append(self.generate_split_flow())

        # Add rules reminder
        output.append("")
        output.append("-" * 80)
        output.append("KEY RULES & CONSTRAINTS")
        output.append("-" * 80)
        output.append("• R313: STOP after spawning agents (context preservation)")
        output.append("• R322: Checkpoint before critical transitions")
        output.append("• R234: No state skipping - mandatory traversal")
        output.append("• R356: Single-effort optimization allowed")
        output.append("• R336: Wave integration before next wave")
        output.append("• R321: Immediate backport during integration")
        output.append("")
        output.append("PARALLELIZATION RULES:")
        output.append("• Same agent type in parallel: ALLOWED")
        output.append("• Different agent types: FORBIDDEN")
        output.append("• Timing delta must be <5s between parallel spawns")
        output.append("")
        output.append("-" * 80)

        return "\n".join(output)

    def run(self):
        """Main execution"""
        if not self.load_data():
            print("Error: Could not load state machine data")
            return 1

        print(self.generate_complete_flow())
        return 0

def main():
    """Main entry point"""
    import argparse

    parser = argparse.ArgumentParser(description="Visualize State Machine Flow")
    parser.add_argument("-w", "--wave", action="store_true",
                       help="Show only wave flow")
    parser.add_argument("-i", "--integration", action="store_true",
                       help="Show only integration flow")
    parser.add_argument("-f", "--fix", action="store_true",
                       help="Show only fix cascade flow")
    parser.add_argument("-s", "--split", action="store_true",
                       help="Show only split flow")

    args = parser.parse_args()

    visualizer = FlowVisualizer()

    if not visualizer.load_data():
        print("Error loading data")
        return 1

    # Show specific flow if requested
    if args.wave:
        print(visualizer.generate_wave_flow())
    elif args.integration:
        print(visualizer.generate_integration_flow())
    elif args.fix:
        print(visualizer.generate_fix_cascade_flow())
    elif args.split:
        print(visualizer.generate_split_flow())
    else:
        # Show complete flow
        return visualizer.run()

    return 0

if __name__ == "__main__":
    sys.exit(main())