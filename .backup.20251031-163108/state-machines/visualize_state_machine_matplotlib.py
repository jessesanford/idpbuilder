#!/usr/bin/env python3
"""
Software Factory 2.0 - State Machine Visualization Tool (Matplotlib version)
Generates a comprehensive graph visualization of the entire state machine
"""

import json
import matplotlib.pyplot as plt
import matplotlib.patches as mpatches
import networkx as nx
from collections import defaultdict
import numpy as np

# Set matplotlib to use a non-GUI backend
import matplotlib
matplotlib.use('Agg')

def load_state_machine(filepath):
    """Load the state machine JSON file"""
    with open(filepath, 'r') as f:
        return json.load(f)

def categorize_states(states):
    """Categorize states by type for coloring"""
    categories = {
        'init': [],
        'spawn': [],
        'waiting': [],
        'monitoring': [],
        'integration': [],
        'error': [],
        'terminal': [],
        'review': [],
        'infrastructure': [],
        'other': []
    }

    for state in states:
        state_lower = state.lower()
        if 'init' in state_lower:
            categories['init'].append(state)
        elif 'spawn' in state_lower:
            categories['spawn'].append(state)
        elif 'waiting' in state_lower:
            categories['waiting'].append(state)
        elif 'monitor' in state_lower:
            categories['monitoring'].append(state)
        elif 'integration' in state_lower:
            categories['integration'].append(state)
        elif 'error' in state_lower or 'recovery' in state_lower:
            categories['error'].append(state)
        elif state in ['PROJECT_DONE', 'ERROR_RECOVERY', 'COMPLETE_PHASE', 'WAVE_COMPLETE', 'PROJECT_COMPLETE']:
            categories['terminal'].append(state)
        elif 'review' in state_lower:
            categories['review'].append(state)
        elif 'infrastructure' in state_lower or 'setup' in state_lower or 'create' in state_lower and 'infrastructure' in state_lower:
            categories['infrastructure'].append(state)
        else:
            categories['other'].append(state)

    return categories

def get_color_for_category(category):
    """Return color for each category"""
    color_map = {
        'init': '#4CAF50',        # Green
        'spawn': '#2196F3',        # Blue
        'waiting': '#FFC107',      # Amber
        'monitoring': '#FF9800',   # Orange
        'integration': '#9C27B0',  # Purple
        'error': '#F44336',        # Red
        'terminal': '#795548',     # Brown
        'review': '#00BCD4',       # Cyan
        'infrastructure': '#607D8B', # Blue Grey
        'other': '#9E9E9E'         # Grey
    }
    return color_map.get(category, '#9E9E9E')

def build_graph(data):
    """Build NetworkX graph from state machine data"""
    G = nx.DiGraph()

    # Add all states as nodes
    states = data['states']
    for state in states:
        G.add_node(state)

    # Add edges from transition matrix
    # Focus on orchestrator as the main agent
    if 'orchestrator' in data['transition_matrix']:
        transitions = data['transition_matrix']['orchestrator']
        for from_state, to_states in transitions.items():
            if from_state in states:  # Ensure state exists
                for to_state in to_states:
                    if to_state in states:  # Ensure destination exists
                        G.add_edge(from_state, to_state)

    # Add transitions from other agents if they connect to orchestrator states
    for agent in ['sw-engineer', 'code-reviewer', 'architect']:
        if agent in data['transition_matrix']:
            agent_transitions = data['transition_matrix'][agent]
            for from_state, to_states in agent_transitions.items():
                if from_state in states:
                    for to_state in to_states:
                        if to_state in states and not G.has_edge(from_state, to_state):
                            # Add edge with agent label
                            G.add_edge(from_state, to_state, agent=agent)

    return G

def create_hierarchical_layout(G, states):
    """Create a hierarchical layout for better visualization"""
    # Group states into layers
    layers = {
        0: [],  # Init states
        1: [],  # Spawn states
        2: [],  # Waiting states
        3: [],  # Monitor/Other states
        4: [],  # Review states
        5: [],  # Integration states
        6: [],  # Terminal states
        7: []   # Error states
    }

    for state in states:
        state_lower = state.lower()
        if 'init' in state_lower:
            layers[0].append(state)
        elif 'spawn' in state_lower:
            layers[1].append(state)
        elif 'waiting' in state_lower:
            layers[2].append(state)
        elif 'review' in state_lower:
            layers[4].append(state)
        elif 'integration' in state_lower:
            layers[5].append(state)
        elif state in ['PROJECT_DONE', 'ERROR_RECOVERY', 'COMPLETE_PHASE', 'WAVE_COMPLETE']:
            layers[6].append(state)
        elif 'error' in state_lower or 'recovery' in state_lower:
            layers[7].append(state)
        else:
            layers[3].append(state)  # Monitor and other states

    # Create positions
    pos = {}
    for layer_idx, nodes in layers.items():
        if not nodes:
            continue

        # Vertical position based on layer
        y = 10 - layer_idx * 1.5

        # Horizontal positions spread across the layer
        n_nodes = len(nodes)
        if n_nodes == 1:
            x_positions = [0]
        else:
            x_positions = np.linspace(-10, 10, n_nodes)

        for i, node in enumerate(sorted(nodes)):
            if node in G.nodes():
                pos[node] = (x_positions[i], y)

    return pos

def create_full_visualization(G, data):
    """Create the complete state machine visualization"""

    # Create a large figure for the complete graph
    fig, ax = plt.subplots(figsize=(40, 30))

    # Categorize states
    categories = categorize_states(data['states'])

    # Create hierarchical layout
    pos = create_hierarchical_layout(G, data['states'])

    # If some nodes are missing positions, use spring layout for them
    missing_nodes = [n for n in G.nodes() if n not in pos]
    if missing_nodes:
        subgraph = G.subgraph(missing_nodes)
        sub_pos = nx.spring_layout(subgraph, k=2, iterations=50)
        for node, (x, y) in sub_pos.items():
            pos[node] = (x * 5, y * 5 - 12)  # Place below main hierarchy

    # Draw edges
    nx.draw_networkx_edges(G, pos, edge_color='gray', arrows=True,
                           arrowsize=10, alpha=0.5, ax=ax,
                           arrowstyle='->', node_size=3000,
                           connectionstyle='arc3,rad=0.1')

    # Draw nodes by category with different colors
    for category, nodes in categories.items():
        node_list = [n for n in nodes if n in G.nodes() and n in pos]
        if node_list:
            nx.draw_networkx_nodes(G, pos, nodelist=node_list,
                                  node_color=get_color_for_category(category),
                                  node_size=3000, ax=ax, alpha=0.9,
                                  label=category.title())

    # Draw labels with smaller font
    labels = {}
    for node in G.nodes():
        if len(node) > 20:
            # Truncate long names
            labels[node] = node[:17] + '...'
        else:
            labels[node] = node

    nx.draw_networkx_labels(G, pos, labels, font_size=6, ax=ax)

    # Add legend
    legend_patches = []
    for category in categories.keys():
        if categories[category]:  # Only add to legend if category has nodes
            patch = mpatches.Patch(color=get_color_for_category(category),
                                  label=f'{category.title()} ({len(categories[category])})')
            legend_patches.append(patch)

    ax.legend(handles=legend_patches, loc='upper left', fontsize=10,
             bbox_to_anchor=(1.02, 1), borderaxespad=0)

    # Set title and remove axes
    ax.set_title('Software Factory 2.0 - Complete State Machine', fontsize=24, fontweight='bold')
    ax.axis('off')

    # Adjust layout to prevent label cutoff
    plt.tight_layout()

    return fig

def create_simplified_visualization(G, data):
    """Create a simplified visualization showing main flow"""

    # Define critical states for simplified view
    critical_states = [
        'INIT', 'WAVE_START', 'CREATE_NEXT_INFRASTRUCTURE',
        'SPAWN_SW_ENGINEERS', 'SPAWN_CODE_REVIEWERS_EFFORT_PLANNING',
        'WAITING_FOR_EFFORT_PLANS', 'SPAWN_CODE_REVIEWERS_EFFORT_REVIEW',
        'MONITORING_SWE_PROGRESS', 'MONITORING_EFFORT_REVIEWS',
        'WAVE_COMPLETE', 'INTEGRATE_WAVE_EFFORTS', 'COMPLETE_PHASE',
        'PROJECT_DONE', 'ERROR_RECOVERY', 'CASCADE_REINTEGRATION'
    ]

    # Find these states in the graph
    subgraph_nodes = set()
    for state in critical_states:
        if state in G.nodes():
            subgraph_nodes.add(state)
            # Add immediate neighbors for context
            for neighbor in G.successors(state):
                if any(keyword in neighbor.lower() for keyword in ['spawn', 'waiting', 'monitor', 'complete', 'error']):
                    subgraph_nodes.add(neighbor)

    # Create subgraph
    subgraph = G.subgraph(subgraph_nodes)

    # Create figure
    fig, ax = plt.subplots(figsize=(20, 15))

    # Use hierarchical layout
    pos = nx.spring_layout(subgraph, k=3, iterations=100)

    # Categorize states for coloring
    categories = categorize_states(list(subgraph_nodes))

    # Draw edges with arrows
    nx.draw_networkx_edges(subgraph, pos, edge_color='gray', arrows=True,
                          arrowsize=15, alpha=0.6, ax=ax,
                          arrowstyle='->', node_size=5000,
                          connectionstyle='arc3,rad=0.1', width=2)

    # Draw nodes by category
    for category, nodes in categories.items():
        node_list = [n for n in nodes if n in subgraph.nodes()]
        if node_list:
            nx.draw_networkx_nodes(subgraph, pos, nodelist=node_list,
                                  node_color=get_color_for_category(category),
                                  node_size=5000, ax=ax, alpha=0.9,
                                  label=category.title())

    # Draw labels
    labels = {}
    for node in subgraph.nodes():
        if len(node) > 25:
            labels[node] = node[:22] + '...'
        else:
            labels[node] = node

    nx.draw_networkx_labels(subgraph, pos, labels, font_size=9, ax=ax)

    # Add legend
    legend_patches = []
    for category in categories.keys():
        if categories[category]:
            patch = mpatches.Patch(color=get_color_for_category(category),
                                  label=category.title())
            legend_patches.append(patch)

    ax.legend(handles=legend_patches, loc='upper right', fontsize=11)

    # Set title and remove axes
    ax.set_title('Software Factory 2.0 - Simplified Main Flow', fontsize=20, fontweight='bold')
    ax.axis('off')

    plt.tight_layout()

    return fig

def create_agent_specific_view(G, data, agent_name):
    """Create a view specific to an agent's states"""

    # Get agent-specific transitions
    if agent_name not in data['transition_matrix']:
        return None

    agent_transitions = data['transition_matrix'][agent_name]

    # Build agent-specific graph
    agent_graph = nx.DiGraph()

    for from_state, to_states in agent_transitions.items():
        if from_state in data['states']:
            agent_graph.add_node(from_state)
            for to_state in to_states:
                if to_state in data['states']:
                    agent_graph.add_edge(from_state, to_state)

    if not agent_graph.nodes():
        return None

    # Create figure
    fig, ax = plt.subplots(figsize=(16, 12))

    # Layout
    pos = nx.spring_layout(agent_graph, k=2, iterations=100)

    # Categorize states
    categories = categorize_states(list(agent_graph.nodes()))

    # Draw edges
    nx.draw_networkx_edges(agent_graph, pos, edge_color='darkgray', arrows=True,
                          arrowsize=12, alpha=0.7, ax=ax,
                          arrowstyle='->', node_size=4000, width=1.5)

    # Draw nodes by category
    for category, nodes in categories.items():
        node_list = [n for n in nodes if n in agent_graph.nodes()]
        if node_list:
            nx.draw_networkx_nodes(agent_graph, pos, nodelist=node_list,
                                  node_color=get_color_for_category(category),
                                  node_size=4000, ax=ax, alpha=0.9)

    # Draw labels
    nx.draw_networkx_labels(agent_graph, pos, font_size=10, ax=ax)

    # Set title and remove axes
    ax.set_title(f'Software Factory 2.0 - {agent_name.title()} Agent States',
                fontsize=18, fontweight='bold')
    ax.axis('off')

    plt.tight_layout()

    return fig

def generate_statistics(G, data):
    """Generate statistics about the state machine"""
    stats = {
        'total_states': len(data['states']),
        'total_transitions': G.number_of_edges(),
        'agents': len(data.get('agents', {})),
        'max_in_degree': 0,
        'max_out_degree': 0,
        'most_connected_state': None,
        'terminal_states': [],
        'entry_points': [],
        'cycles_detected': False
    }

    # Find states with most connections
    max_connections = 0
    for node in G.nodes():
        in_deg = G.in_degree(node)
        out_deg = G.out_degree(node)
        total = in_deg + out_deg

        if in_deg > stats['max_in_degree']:
            stats['max_in_degree'] = in_deg
        if out_deg > stats['max_out_degree']:
            stats['max_out_degree'] = out_deg
        if total > max_connections:
            max_connections = total
            stats['most_connected_state'] = node

    # Find terminal states (no outgoing edges)
    for node in G.nodes():
        if G.out_degree(node) == 0:
            stats['terminal_states'].append(node)

    # Find entry points (no incoming edges or INIT states)
    for node in G.nodes():
        if G.in_degree(node) == 0 or 'INIT' in node:
            stats['entry_points'].append(node)

    # Check for cycles
    try:
        stats['cycles_detected'] = not nx.is_directed_acyclic_graph(G)
    except:
        stats['cycles_detected'] = None

    return stats

def main():
    """Main function to generate visualizations"""

    print("Loading state machine data...")
    data = load_state_machine('software-factory-3.0-state-machine.json')

    print("Categorizing states...")
    categories = categorize_states(data['states'])

    print("\nState Categories:")
    for cat, states in categories.items():
        if states:
            print(f"  {cat.title()}: {len(states)} states")

    print("\nBuilding graph...")
    G = build_graph(data)

    print("Generating statistics...")
    stats = generate_statistics(G, data)

    print("\nState Machine Statistics:")
    print(f"  Total States: {stats['total_states']}")
    print(f"  Total Transitions: {stats['total_transitions']}")
    print(f"  Entry Points: {len(stats['entry_points'])}")
    print(f"  Terminal States: {len(stats['terminal_states'])}")
    print(f"  Most Connected: {stats['most_connected_state']}")
    print(f"  Max In-Degree: {stats['max_in_degree']}")
    print(f"  Max Out-Degree: {stats['max_out_degree']}")
    print(f"  Contains Cycles: {stats['cycles_detected']}")

    print("\nCreating visualizations...")

    # Create complete visualization
    print("  1. Creating complete graph...")
    fig_complete = create_full_visualization(G, data)
    fig_complete.savefig('software-factory-state-machine-complete.png',
                        dpi=150, bbox_inches='tight', facecolor='white')
    plt.close(fig_complete)

    # Create simplified visualization
    print("  2. Creating simplified graph...")
    fig_simple = create_simplified_visualization(G, data)
    fig_simple.savefig('software-factory-state-machine-simplified.png',
                      dpi=150, bbox_inches='tight', facecolor='white')
    plt.close(fig_simple)

    # Create agent-specific views
    agents = ['orchestrator', 'sw-engineer', 'code-reviewer', 'architect']
    for agent in agents:
        print(f"  3. Creating {agent} agent view...")
        fig_agent = create_agent_specific_view(G, data, agent)
        if fig_agent:
            fig_agent.savefig(f'software-factory-state-machine-{agent}.png',
                            dpi=150, bbox_inches='tight', facecolor='white')
            plt.close(fig_agent)

    # Save statistics to file
    print("\nSaving statistics...")
    with open('state-machine-statistics.txt', 'w') as f:
        f.write("Software Factory 2.0 - State Machine Statistics\n")
        f.write("=" * 50 + "\n\n")
        f.write(f"Total States: {stats['total_states']}\n")
        f.write(f"Total Transitions: {stats['total_transitions']}\n")
        f.write(f"Number of Agents: {stats.get('agents', 'N/A')}\n\n")

        f.write("Entry Points:\n")
        for ep in stats['entry_points']:
            f.write(f"  - {ep}\n")

        f.write("\nTerminal States:\n")
        for ts in stats['terminal_states']:
            f.write(f"  - {ts}\n")

        f.write(f"\nMost Connected State: {stats['most_connected_state']}\n")
        f.write(f"Maximum In-Degree: {stats['max_in_degree']}\n")
        f.write(f"Maximum Out-Degree: {stats['max_out_degree']}\n")
        f.write(f"Contains Cycles: {stats['cycles_detected']}\n\n")

        f.write("State Categories:\n")
        for cat, states in categories.items():
            if states:
                f.write(f"\n{cat.title()} ({len(states)} states):\n")
                for state in sorted(states):
                    f.write(f"  - {state}\n")

    print("\n✅ Visualization complete!")
    print("\nGenerated files:")
    print("  📊 software-factory-state-machine-complete.png - Full state machine")
    print("  📊 software-factory-state-machine-simplified.png - Simplified main flow")
    print("  📊 software-factory-state-machine-orchestrator.png - Orchestrator agent view")
    print("  📊 software-factory-state-machine-sw-engineer.png - SW Engineer agent view")
    print("  📊 software-factory-state-machine-code-reviewer.png - Code Reviewer agent view")
    print("  📊 software-factory-state-machine-architect.png - Architect agent view")
    print("  📄 state-machine-statistics.txt - Detailed statistics")

if __name__ == "__main__":
    main()