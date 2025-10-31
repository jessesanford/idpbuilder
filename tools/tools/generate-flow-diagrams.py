#!/usr/bin/env python3
"""
Flow Diagram Generator for Software Factory 2.0 State Machine
Generates professional flowchart-style PNG diagrams of the state machine
"""

import json
import os
from pathlib import Path
from graphviz import Digraph
import argparse
from typing import Dict, List, Set, Tuple

# Color scheme for different state types
STATE_COLORS = {
    'INIT': '#90EE90',  # Light green
    'PLANNING': '#87CEEB',  # Sky blue
    'IMPLEMENTATION': '#FFD700',  # Gold
    'REVIEW': '#DDA0DD',  # Plum
    'ERROR': '#FF6B6B',  # Light red
    'PROJECT_DONE': '#98FB98',  # Pale green
    'INTEGRATE_WAVE_EFFORTS': '#FFA500',  # Orange
    'MONITOR': '#B0C4DE',  # Light steel blue
    'RECOVERY': '#F0E68C',  # Khaki
    'SPLIT': '#FFB6C1',  # Light pink
    'SPAWN': '#D8BFD8',  # Thistle
    'WAVE': '#20B2AA',  # Light sea green
    'PHASE': '#4682B4',  # Steel blue
    'MEASURE': '#DEB887',  # Burlywood
    'FIX': '#F4A460',  # Sandy brown
    'TEST': '#9ACD32',  # Yellow green
    'WAITING': '#E6E6FA',  # Lavender
    'COMPLETE': '#98FB98',  # Pale green
    'DEFAULT': '#E0E0E0'  # Light gray
}

# Shape mapping for different state categories
STATE_SHAPES = {
    'INIT': 'ellipse',
    'ERROR': 'octagon',
    'PROJECT_DONE': 'doublecircle',
    'COMPLETE': 'doublecircle',
    'decision': 'diamond',
    'default': 'box'
}

class FlowDiagramGenerator:
    def __init__(self, state_machine_path: str, high_res: bool = True, direction: str = 'horizontal'):
        """Initialize with state machine JSON file

        Args:
            state_machine_path: Path to the state machine JSON file
            high_res: If True, generate high-resolution diagrams (300 DPI),
                     otherwise generate standard resolution (150 DPI)
            direction: Layout direction - 'horizontal' (LR) or 'vertical' (TB)
        """
        self.state_machine_path = Path(state_machine_path)
        self.output_dir = Path('state-machines/flow-diagrams')
        self.output_dir.mkdir(parents=True, exist_ok=True)
        self.high_res = high_res
        self.direction = direction

        # Set rankdir based on direction
        self.rankdir = 'LR' if direction == 'horizontal' else 'TB'

        # Set DPI and size multiplier based on resolution mode
        self.dpi = '300' if high_res else '150'
        self.size_mult = 2.0 if high_res else 1.0

        # Load state machine
        with open(self.state_machine_path, 'r') as f:
            self.state_machine = json.load(f)

        # Build complete states dictionary from agent definitions
        self.states = {}
        self.agents = self.state_machine.get('agents', {})

        # Extract all states from agent definitions
        for agent_name, agent_data in self.agents.items():
            agent_states = agent_data.get('states', {})
            if isinstance(agent_states, dict):
                for state_name, state_info in agent_states.items():
                    # Store state with agent association
                    self.states[state_name] = {
                        'agent': agent_name,
                        'type': state_info.get('type', ''),
                        'description': state_info.get('description', ''),
                        'transitions': state_info.get('valid_transitions', []),
                        'full_info': state_info
                    }

    def get_state_color(self, state_name: str) -> str:
        """Get color for a state based on its type"""
        state_upper = state_name.upper()

        # Check each keyword in order of priority
        for key in ['ERROR', 'PROJECT_DONE', 'COMPLETE', 'INIT', 'WAITING', 'PLANNING',
                   'IMPLEMENTATION', 'REVIEW', 'INTEGRATE_WAVE_EFFORTS', 'MONITOR', 'RECOVERY',
                   'SPLIT', 'SPAWN', 'WAVE', 'PHASE', 'MEASURE', 'FIX', 'TEST']:
            if key in state_upper:
                return STATE_COLORS.get(key, STATE_COLORS['DEFAULT'])

        return STATE_COLORS['DEFAULT']

    def get_state_shape(self, state_name: str) -> str:
        """Get shape for a state based on its type"""
        state_upper = state_name.upper()

        if 'INIT' in state_upper:
            return STATE_SHAPES['INIT']
        elif 'ERROR' in state_upper:
            return STATE_SHAPES['ERROR']
        elif 'PROJECT_DONE' in state_upper or 'COMPLETE' in state_upper:
            return STATE_SHAPES['COMPLETE']
        elif any(word in state_upper for word in ['DECISION', 'CHECK', 'VERIFY']):
            return STATE_SHAPES['decision']
        else:
            return STATE_SHAPES['default']

    def create_overall_flow(self):
        """Create the main system flow diagram"""
        print("Creating overall system flow diagram...")

        dot = Digraph('Software Factory State Machine',
                     comment='Software Factory 2.0 Overall Flow',
                     format='png')

        # Configure graph attributes for flowchart layout
        # Use dynamic DPI and sizes based on resolution mode and direction
        if self.direction == 'horizontal':
            # Wider for horizontal layout
            width = int(30 * self.size_mult)
            height = int(20 * self.size_mult)
        else:
            # Taller for vertical layout
            width = int(20 * self.size_mult)
            height = int(30 * self.size_mult)

        dot.attr(rankdir=self.rankdir, size=f'{width},{height}', dpi=self.dpi)
        ranksep = 1.5 if self.high_res else 1.0
        nodesep = 0.8 if self.high_res else 0.5
        dot.attr('graph', bgcolor='white', pad='0.5', ranksep=str(ranksep), nodesep=str(nodesep))
        # Scale font sizes based on resolution
        node_fontsize = '12' if self.high_res else '9'
        edge_fontsize = '11' if self.high_res else '8'
        penwidth = '2.0' if self.high_res else '1.5'
        dot.attr('node', fontname='Arial', fontsize=node_fontsize, style='filled', penwidth=penwidth)
        dot.attr('edge', fontname='Arial', fontsize=edge_fontsize, penwidth='1.5')

        # Track which states we've added to avoid duplicates
        added_states = set()

        # Group states by agent for subgraphs
        for agent_name, agent_data in self.agents.items():
            with dot.subgraph(name=f'cluster_{agent_name}') as subgraph:
                cluster_fontsize = '14' if self.high_res else '11'
                subgraph.attr(label=f'{agent_name.upper()} Agent',
                            style='rounded,filled',
                            fillcolor='#f0f0f0',
                            fontsize=cluster_fontsize,
                            fontname='Arial Bold')

                agent_states = agent_data.get('states', {})

                if isinstance(agent_states, dict):
                    for state_name in agent_states.keys():
                        if state_name not in added_states:
                            # Create node with appropriate shape and color
                            shape = self.get_state_shape(state_name)
                            color = self.get_state_color(state_name)

                            # Format label - shorter for overview
                            label = state_name.replace('_', '\\n')

                            subgraph.node(state_name,
                                        label=label,
                                        shape=shape,
                                        fillcolor=color,
                                        color='black')

                            added_states.add(state_name)

        # Add transitions between states
        for state_name, state_info in self.states.items():
            transitions = state_info.get('transitions', [])

            # Transitions might be a list of state names or a list of dicts
            for transition in transitions:
                if isinstance(transition, str):
                    next_state = transition
                    condition = ''
                elif isinstance(transition, dict):
                    next_state = transition.get('to_state', '')
                    condition = transition.get('condition', '')
                else:
                    continue

                if state_name in added_states and next_state in added_states:
                    # Different edge styles for different transition types
                    # Scale arrow and pen sizes based on resolution
                    arrow_size = '1.2' if self.high_res else '0.7'
                    arrow_size_std = '1.0' if self.high_res else '0.7'
                    pen_error = '2.0' if self.high_res else '1.0'
                    pen_success = '2.5' if self.high_res else '1.2'
                    pen_std = '1.5' if self.high_res else '1.0'

                    if 'ERROR' in next_state or 'error' in condition.lower():
                        dot.edge(state_name, next_state,
                               color='red',
                               style='dashed',
                               arrowsize=arrow_size,
                               penwidth=pen_error)
                    elif 'PROJECT_DONE' in next_state or 'COMPLETE' in next_state:
                        dot.edge(state_name, next_state,
                               color='green',
                               style='bold',
                               arrowsize=arrow_size,
                               penwidth=pen_success)
                    elif state_name == next_state:  # Self-loop
                        dot.edge(state_name, next_state,
                               color='blue',
                               style='dotted',
                               constraint='false',
                               arrowsize=arrow_size_std,
                               penwidth=pen_std)
                    else:
                        dot.edge(state_name, next_state,
                               color='black',
                               arrowsize=arrow_size_std,
                               penwidth=pen_std)

        # Save the diagram
        output_path = self.output_dir / 'overall-system-flow'
        dot.render(output_path, cleanup=True)
        print(f"✅ Created: {output_path}.png")

    def create_agent_flow(self, agent_name: str):
        """Create flow diagram for a specific agent"""
        print(f"Creating flow diagram for {agent_name} agent...")

        if agent_name not in self.agents:
            print(f"⚠️ Agent {agent_name} not found")
            return

        dot = Digraph(f'{agent_name}_flow',
                     comment=f'{agent_name.upper()} Agent Flow',
                     format='png')

        # Configure for high-resolution flowchart with direction support
        if self.direction == 'horizontal':
            width = int(32 * self.size_mult)
            height = int(24 * self.size_mult)
        else:
            width = int(24 * self.size_mult)
            height = int(32 * self.size_mult)

        dot.attr(rankdir=self.rankdir, size=f'{width},{height}', dpi=self.dpi)
        # Adjust spacing based on resolution
        ranksep = 1.0 if self.high_res else 0.7
        nodesep = 0.6 if self.high_res else 0.4
        fontsize = '18' if self.high_res else '14'
        dot.attr('graph', bgcolor='white', pad='0.5', ranksep=str(ranksep), nodesep=str(nodesep),
                label=f'{agent_name.upper()} Agent State Flow',
                fontsize=fontsize, fontname='Arial Bold', labelloc='t')
        node_fontsize = '13' if self.high_res else '10'
        edge_fontsize = '11' if self.high_res else '9'
        penwidth = '2.5' if self.high_res else '1.5'
        dot.attr('node', fontname='Arial', fontsize=node_fontsize, style='filled', penwidth=penwidth)
        dot.attr('edge', fontname='Arial', fontsize=edge_fontsize, penwidth='2.0')

        agent_data = self.agents[agent_name]
        agent_states = agent_data.get('states', {})

        # Track added states
        added_states = set()

        # Add all states for this agent
        if isinstance(agent_states, dict):
            for state_name, state_info in agent_states.items():
                shape = self.get_state_shape(state_name)
                color = self.get_state_color(state_name)

                # More detailed label for agent-specific view
                label = state_name.replace('_', '\\n')

                # Add type if available
                if state_info.get('type'):
                    state_type = state_info['type']
                    if len(state_type) <= 20:
                        label += f'\\n[{state_type}]'

                dot.node(state_name,
                       label=label,
                       shape=shape,
                       fillcolor=color,
                       color='black')

                added_states.add(state_name)

        # Add transitions between agent states
        for state_name in added_states:
            if state_name in agent_states:
                state_info = agent_states[state_name]
                transitions = state_info.get('valid_transitions', [])

                for transition in transitions:
                    if isinstance(transition, str):
                        next_state = transition
                    elif isinstance(transition, dict):
                        next_state = transition.get('to_state', '')
                    else:
                        continue

                    # Only include transitions within this agent
                    if next_state in added_states:
                        # Color-code edges
                        if 'ERROR' in next_state:
                            dot.edge(state_name, next_state,
                                   color='red',
                                   style='dashed')
                        elif state_name == next_state:  # Self-loop
                            dot.edge(state_name, next_state,
                                   color='blue',
                                   style='dotted',
                                   constraint='false')
                        elif 'PROJECT_DONE' in next_state or 'COMPLETE' in next_state:
                            dot.edge(state_name, next_state,
                                   color='green',
                                   style='bold')
                        else:
                            dot.edge(state_name, next_state,
                                   color='black')

        # Save the diagram
        output_path = self.output_dir / f'{agent_name}-flow'
        dot.render(output_path, cleanup=True)
        print(f"✅ Created: {output_path}.png")

    def create_phase_progression_flow(self):
        """Create a flow showing phase and wave progression"""
        print("Creating phase progression flow diagram...")

        dot = Digraph('phase_progression',
                     comment='Phase and Wave Progression Flow',
                     format='png')

        # Configure based on direction and resolution
        # Phase progression especially benefits from horizontal flow
        if self.direction == 'horizontal':
            width = int(32 * self.size_mult)
            height = int(20 * self.size_mult)
        else:
            width = int(20 * self.size_mult)
            height = int(32 * self.size_mult)

        dot.attr(rankdir=self.rankdir, size=f'{width},{height}', dpi=self.dpi)
        ranksep = 1.8 if self.high_res else 1.2
        nodesep = 0.9 if self.high_res else 0.6
        fontsize = '18' if self.high_res else '14'
        dot.attr('graph', bgcolor='white', pad='0.5', ranksep=str(ranksep), nodesep=str(nodesep),
                label='Phase and Wave Progression', fontsize=fontsize,
                fontname='Arial Bold', labelloc='t')
        node_fontsize = '13' if self.high_res else '10'
        edge_fontsize = '11' if self.high_res else '9'
        penwidth = '2.5' if self.high_res else '1.5'
        dot.attr('node', fontname='Arial', fontsize=node_fontsize, style='filled', penwidth=penwidth)
        dot.attr('edge', fontname='Arial', fontsize=edge_fontsize, penwidth=penwidth)

        # Find phase and wave related states
        phase_states = []
        wave_states = []
        planning_states = []

        for state_name in self.states:
            if 'PHASE' in state_name.upper():
                phase_states.append(state_name)
            elif 'WAVE' in state_name.upper():
                wave_states.append(state_name)
            elif 'PLANNING' in state_name.upper():
                planning_states.append(state_name)

        # Create phase nodes
        if phase_states:
            with dot.subgraph(name='cluster_phases') as phases:
                phases.attr(label='Phase States',
                           style='rounded,filled',
                           fillcolor='#e6f3ff',
                           fontsize='16',
                           fontname='Arial Bold')

                for state in phase_states:
                    phases.node(state,
                              label=state.replace('_', '\\n'),
                              shape='box3d',
                              fillcolor=STATE_COLORS['PHASE'],
                              color='black')

        # Create wave nodes
        if wave_states:
            with dot.subgraph(name='cluster_waves') as waves:
                waves.attr(label='Wave States',
                          style='rounded,filled',
                          fillcolor='#ffe6e6',
                          fontsize='16',
                          fontname='Arial Bold')

                for state in wave_states:
                    waves.node(state,
                             label=state.replace('_', '\\n'),
                             shape='box',
                             fillcolor=STATE_COLORS['WAVE'],
                             color='black')

        # Create planning nodes
        if planning_states[:5]:  # Limit to avoid clutter
            with dot.subgraph(name='cluster_planning') as planning:
                planning.attr(label='Planning States',
                             style='rounded,filled',
                             fillcolor='#f0f0e6',
                             fontsize='16',
                             fontname='Arial Bold')

                for state in planning_states[:5]:
                    planning.node(state,
                                label=state.replace('_', '\\n'),
                                shape='parallelogram',
                                fillcolor=STATE_COLORS['PLANNING'],
                                color='black')

        # Add some key transitions to show flow
        all_progression_states = phase_states + wave_states + planning_states[:5]

        for state_name in all_progression_states:
            if state_name in self.states:
                state_info = self.states[state_name]
                transitions = state_info.get('transitions', [])

                for transition in transitions:
                    if isinstance(transition, str):
                        next_state = transition
                    else:
                        continue

                    if next_state in all_progression_states:
                        if 'COMPLETE' in next_state:
                            dot.edge(state_name, next_state,
                                   color='green',
                                   style='bold')
                        else:
                            dot.edge(state_name, next_state,
                                   color='blue')

        # Save the diagram
        output_path = self.output_dir / 'phase-progression-flow'
        dot.render(output_path, cleanup=True)
        print(f"✅ Created: {output_path}.png")

    def create_error_recovery_flow(self):
        """Create a flow showing error states and recovery paths"""
        print("Creating error recovery flow diagram...")

        dot = Digraph('error_recovery',
                     comment='Error Recovery Flow',
                     format='png')

        # Configure based on direction and resolution
        if self.direction == 'horizontal':
            width = int(28 * self.size_mult)
            height = int(24 * self.size_mult)
        else:
            width = int(24 * self.size_mult)
            height = int(28 * self.size_mult)

        dot.attr(rankdir=self.rankdir, size=f'{width},{height}', dpi=self.dpi)
        ranksep = 1.2 if self.high_res else 0.8
        nodesep = 0.7 if self.high_res else 0.5
        fontsize = '18' if self.high_res else '14'
        dot.attr('graph', bgcolor='#fff5f5', pad='0.5', ranksep=str(ranksep), nodesep=str(nodesep),
                label='Error Recovery Flow', fontsize=fontsize,
                fontname='Arial Bold', labelloc='t')
        node_fontsize = '13' if self.high_res else '10'
        edge_fontsize = '11' if self.high_res else '9'
        penwidth = '2.5' if self.high_res else '1.5'
        dot.attr('node', fontname='Arial', fontsize=node_fontsize, style='filled', penwidth=penwidth)
        dot.attr('edge', fontname='Arial', fontsize=edge_fontsize, penwidth='2.0')

        # Find error and recovery states
        error_states = []
        recovery_states = []
        normal_states_with_error_transitions = set()

        for state_name in self.states:
            if 'ERROR' in state_name.upper():
                error_states.append(state_name)
            elif 'RECOVERY' in state_name.upper() or 'FIX' in state_name.upper():
                recovery_states.append(state_name)

        # Find states that can transition to error states
        for state_name, state_info in self.states.items():
            transitions = state_info.get('transitions', [])
            for transition in transitions:
                if isinstance(transition, str) and transition in error_states:
                    normal_states_with_error_transitions.add(state_name)

        # Limit normal states to avoid clutter
        normal_states_list = list(normal_states_with_error_transitions)[:10]

        # Add normal states that can error
        if normal_states_list:
            with dot.subgraph(name='cluster_normal') as normal:
                normal.attr(label='States That Can Error',
                           style='rounded,filled',
                           fillcolor='#f0f0f0',
                           fontsize='16',
                           fontname='Arial Bold')

                for state in normal_states_list:
                    if state not in error_states and state not in recovery_states:
                        normal.node(state,
                                  label=state.replace('_', '\\n'),
                                  shape='box',
                                  fillcolor=self.get_state_color(state),
                                  color='black')

        # Add error states
        if error_states:
            with dot.subgraph(name='cluster_errors') as errors:
                errors.attr(label='Error States',
                           style='rounded,filled',
                           fillcolor='#ffcccc',
                           fontsize='16',
                           fontname='Arial Bold')

                for state in error_states:
                    errors.node(state,
                              label=state.replace('_', '\\n'),
                              shape='octagon',
                              fillcolor=STATE_COLORS['ERROR'],
                              color='darkred',
                              penwidth='3')

        # Add recovery states
        if recovery_states:
            with dot.subgraph(name='cluster_recovery') as recovery:
                recovery.attr(label='Recovery States',
                             style='rounded,filled',
                             fillcolor='#ccffcc',
                             fontsize='16',
                             fontname='Arial Bold')

                for state in recovery_states:
                    recovery.node(state,
                                label=state.replace('_', '\\n'),
                                shape='hexagon',
                                fillcolor=STATE_COLORS['RECOVERY'],
                                color='darkgreen',
                                penwidth='2')

        # Add transitions
        all_error_related = normal_states_list + error_states + recovery_states

        for state_name in all_error_related:
            if state_name in self.states:
                state_info = self.states[state_name]
                transitions = state_info.get('transitions', [])

                for transition in transitions:
                    if isinstance(transition, str):
                        next_state = transition
                    else:
                        continue

                    if next_state in all_error_related:
                        # Color code based on transition type
                        if next_state in error_states:
                            dot.edge(state_name, next_state,
                                   label='ERROR',
                                   color='red',
                                   style='dashed',
                                   penwidth='2')
                        elif state_name in error_states and next_state in recovery_states:
                            dot.edge(state_name, next_state,
                                   label='RECOVER',
                                   color='orange',
                                   style='bold',
                                   penwidth='2')
                        elif state_name in recovery_states:
                            dot.edge(state_name, next_state,
                                   label='FIXED',
                                   color='green',
                                   style='bold',
                                   penwidth='2')

        # Save the diagram
        output_path = self.output_dir / 'error-recovery-flow'
        dot.render(output_path, cleanup=True)
        print(f"✅ Created: {output_path}.png")

    def create_split_flow(self):
        """Create a flow showing the split implementation process"""
        print("Creating split implementation flow diagram...")

        dot = Digraph('split_flow',
                     comment='Split Implementation Flow',
                     format='png')

        # Configure based on direction and resolution
        if self.direction == 'horizontal':
            width = int(28 * self.size_mult)
            height = int(24 * self.size_mult)
        else:
            width = int(24 * self.size_mult)
            height = int(28 * self.size_mult)

        dot.attr(rankdir=self.rankdir, size=f'{width},{height}', dpi=self.dpi)
        ranksep = 1.0 if self.high_res else 0.7
        nodesep = 0.6 if self.high_res else 0.4
        fontsize = '18' if self.high_res else '14'
        dot.attr('graph', bgcolor='white', pad='0.5', ranksep=str(ranksep), nodesep=str(nodesep),
                label='Split Implementation Flow', fontsize=fontsize,
                fontname='Arial Bold', labelloc='t')
        node_fontsize = '13' if self.high_res else '10'
        edge_fontsize = '11' if self.high_res else '9'
        penwidth = '2.5' if self.high_res else '1.5'
        dot.attr('node', fontname='Arial', fontsize=node_fontsize, style='filled', penwidth=penwidth)
        dot.attr('edge', fontname='Arial', fontsize=edge_fontsize, penwidth='2.0')

        # Find split-related states
        split_states = []
        measure_states = []
        infrastructure_states = []

        for state_name in self.states:
            if 'SPLIT' in state_name.upper():
                split_states.append(state_name)
            elif 'MEASURE' in state_name.upper():
                measure_states.append(state_name)
            elif 'INFRASTRUCTURE' in state_name.upper():
                infrastructure_states.append(state_name)

        # Add measurement states
        if measure_states:
            with dot.subgraph(name='cluster_measure') as measure:
                measure.attr(label='Size Measurement',
                            style='rounded,filled',
                            fillcolor='#fff0e6',
                            fontsize='16',
                            fontname='Arial Bold')

                for state in measure_states:
                    measure.node(state,
                               label=state.replace('_', '\\n'),
                               shape='diamond',
                               fillcolor=STATE_COLORS['MEASURE'],
                               color='black')

        # Add split states
        if split_states:
            with dot.subgraph(name='cluster_split') as split:
                split.attr(label='Split Implementation',
                          style='rounded,filled',
                          fillcolor='#ffe6f0',
                          fontsize='16',
                          fontname='Arial Bold')

                for state in split_states:
                    split.node(state,
                             label=state.replace('_', '\\n'),
                             shape='box',
                             fillcolor=STATE_COLORS['SPLIT'],
                             color='black')

        # Add infrastructure states
        if infrastructure_states[:5]:  # Limit to avoid clutter
            with dot.subgraph(name='cluster_infra') as infra:
                infra.attr(label='Infrastructure Setup',
                          style='rounded,filled',
                          fillcolor='#e6ffe6',
                          fontsize='16',
                          fontname='Arial Bold')

                for state in infrastructure_states[:5]:
                    infra.node(state,
                             label=state.replace('_', '\\n'),
                             shape='component',
                             fillcolor='#c0ffc0',
                             color='black')

        # Add transitions
        all_split_related = split_states + measure_states + infrastructure_states[:5]

        for state_name in all_split_related:
            if state_name in self.states:
                state_info = self.states[state_name]
                transitions = state_info.get('transitions', [])

                for transition in transitions:
                    if isinstance(transition, str):
                        next_state = transition
                    else:
                        continue

                    if next_state in all_split_related:
                        if 'EXCEEDED' in state_name or 'EXCEEDED' in next_state:
                            dot.edge(state_name, next_state,
                                   label='Size Exceeded',
                                   color='red',
                                   style='bold')
                        elif 'COMPLETE' in next_state:
                            dot.edge(state_name, next_state,
                                   label='Complete',
                                   color='green',
                                   style='bold')
                        else:
                            dot.edge(state_name, next_state,
                                   color='black')

        # Save the diagram
        output_path = self.output_dir / 'split-implementation-flow'
        dot.render(output_path, cleanup=True)
        print(f"✅ Created: {output_path}.png")

    def generate_all_diagrams(self):
        """Generate all flow diagrams"""
        resolution = "High Resolution (300 DPI)" if self.high_res else "Standard Resolution (150 DPI)"
        layout = self.direction.upper()
        print(f"🚀 Generating all flow diagrams - {resolution} - {layout} layout")
        print(f"Output directory: {self.output_dir}")
        print(f"Found {len(self.states)} states across {len(self.agents)} agents")
        print("-" * 50)

        # Create overall system flow
        self.create_overall_flow()

        # Create agent-specific flows
        for agent_name in self.agents:
            self.create_agent_flow(agent_name)

        # Create specialized flows
        self.create_phase_progression_flow()
        self.create_error_recovery_flow()
        self.create_split_flow()

        print("-" * 50)
        print(f"✅ All diagrams generated in: {self.output_dir}")
        print("\nGenerated diagrams:")
        for png_file in sorted(self.output_dir.glob("*.png")):
            size = png_file.stat().st_size
            print(f"  - {png_file.name} ({size:,} bytes)")

def main():
    """Main entry point"""
    parser = argparse.ArgumentParser(
        description='Generate flow diagrams for Software Factory 2.0 State Machine'
    )
    parser.add_argument(
        '--state-machine',
        default='state-machines/software-factory-3.0-state-machine.json',
        help='Path to state machine JSON file'
    )
    parser.add_argument(
        '--agent',
        help='Generate diagram for specific agent only'
    )
    parser.add_argument(
        '--type',
        choices=['overall', 'phase', 'error', 'split', 'all'],
        default='all',
        help='Type of diagram to generate'
    )
    parser.add_argument(
        '--resolution',
        choices=['high', 'standard', 'both'],
        default='high',
        help='Resolution mode: high (300 DPI), standard (150 DPI), or both'
    )
    parser.add_argument(
        '--direction', '--layout',
        choices=['horizontal', 'vertical', 'LR', 'TB'],
        default='horizontal',
        help='Layout direction: horizontal (left-to-right) or vertical (top-to-bottom). Default is horizontal.'
    )

    args = parser.parse_args()

    # Normalize direction parameter (LR -> horizontal, TB -> vertical)
    direction = args.direction
    if direction == 'LR':
        direction = 'horizontal'
    elif direction == 'TB':
        direction = 'vertical'

    # Check if state machine file exists
    if not os.path.exists(args.state_machine):
        print(f"❌ State machine file not found: {args.state_machine}")
        return 1

    # Helper function to generate diagrams
    def generate_diagrams(gen):
        if args.agent:
            gen.create_agent_flow(args.agent)
        elif args.type == 'overall':
            gen.create_overall_flow()
        elif args.type == 'phase':
            gen.create_phase_progression_flow()
        elif args.type == 'error':
            gen.create_error_recovery_flow()
        elif args.type == 'split':
            gen.create_split_flow()
        else:  # 'all'
            gen.generate_all_diagrams()

    # Handle resolution modes
    if args.resolution == 'both':
        # Generate both high and standard resolution diagrams

        # First, generate high-resolution diagrams
        print(f"🎨 Generating HIGH RESOLUTION diagrams (300 DPI) - {direction.upper()} layout...")
        print("=" * 60)
        generator_high = FlowDiagramGenerator(args.state_machine, high_res=True, direction=direction)
        original_dir = generator_high.output_dir
        generator_high.output_dir = original_dir / 'high-res'
        generator_high.output_dir.mkdir(parents=True, exist_ok=True)
        generate_diagrams(generator_high)

        # Then, generate standard resolution diagrams
        print(f"\n🎨 Generating STANDARD RESOLUTION diagrams (150 DPI) - {direction.upper()} layout...")
        print("=" * 60)
        generator_std = FlowDiagramGenerator(args.state_machine, high_res=False, direction=direction)
        generator_std.output_dir = original_dir / 'standard-res'
        generator_std.output_dir.mkdir(parents=True, exist_ok=True)
        generate_diagrams(generator_std)

        print("\n📊 Generated both high and standard resolution diagrams!")
    else:
        # Single resolution mode
        high_res = args.resolution == 'high'
        res_text = "HIGH RESOLUTION (300 DPI)" if high_res else "STANDARD RESOLUTION (150 DPI)"
        print(f"🎨 Generating {res_text} diagrams - {direction.upper()} layout...")
        generator = FlowDiagramGenerator(args.state_machine, high_res=high_res, direction=direction)
        generate_diagrams(generator)

    return 0

if __name__ == '__main__':
    exit(main())