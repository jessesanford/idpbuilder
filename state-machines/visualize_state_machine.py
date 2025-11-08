#!/usr/bin/env python3
"""
Software Factory 2.0 - State Machine Visualization Tool
Generates a comprehensive graph visualization of the entire state machine
"""

import json
import plotly.graph_objects as go
import plotly.io as pio
import networkx as nx
from collections import defaultdict
import numpy as np

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
        elif 'infrastructure' in state_lower or 'setup' in state_lower:
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
                        if to_state in states:
                            # Add edge with agent label
                            G.add_edge(from_state, to_state, agent=agent)

    return G

def create_layout(G):
    """Create a hierarchical layout for the graph"""
    # Try different layouts and pick the best one

    # Use a combination of layouts for better visualization
    # Start with hierarchical layout for main flow
    try:
        # Try to identify layers based on state types
        layers = defaultdict(list)

        for node in G.nodes():
            if 'INIT' in node:
                layers[0].append(node)
            elif 'SPAWN' in node:
                layers[1].append(node)
            elif 'WAITING' in node:
                layers[2].append(node)
            elif 'MONITOR' in node:
                layers[3].append(node)
            elif 'REVIEW' in node:
                layers[4].append(node)
            elif 'INTEGRATE_WAVE_EFFORTS' in node:
                layers[5].append(node)
            elif 'COMPLETE' in node or 'PROJECT_DONE' in node:
                layers[6].append(node)
            elif 'ERROR' in node or 'RECOVERY' in node:
                layers[7].append(node)
            else:
                layers[3].append(node)  # Default middle layer

        # Create positions manually for better control
        pos = {}
        max_width = max(len(nodes) for nodes in layers.values()) if layers else 1

        for layer_num, nodes in sorted(layers.items()):
            y = -layer_num * 2  # Vertical spacing
            x_spacing = 20.0 / max(len(nodes), 1)  # Horizontal spacing
            x_start = -(len(nodes) - 1) * x_spacing / 2

            for i, node in enumerate(sorted(nodes)):
                x = x_start + i * x_spacing
                pos[node] = (x, y)

        # Add any missing nodes
        missing = set(G.nodes()) - set(pos.keys())
        if missing:
            # Use spring layout for remaining nodes
            subgraph = G.subgraph(missing)
            sub_pos = nx.spring_layout(subgraph, k=2, iterations=50)
            # Offset these positions
            for node, (x, y) in sub_pos.items():
                pos[node] = (x * 10 - 15, y * 10 - 5)

        return pos

    except Exception as e:
        print(f"Layout error: {e}, falling back to spring layout")
        # Fallback to spring layout with high spacing
        return nx.spring_layout(G, k=3, iterations=100, scale=20)

def create_plotly_graph(G, pos, categories):
    """Create interactive Plotly graph"""

    # Create edge traces
    edge_traces = []

    for edge in G.edges():
        x0, y0 = pos[edge[0]]
        x1, y1 = pos[edge[1]]

        # Create arrow annotation for directed edges
        edge_trace = go.Scatter(
            x=[x0, x1, None],
            y=[y0, y1, None],
            mode='lines',
            line=dict(width=0.5, color='#888'),
            hoverinfo='none',
            showlegend=False
        )
        edge_traces.append(edge_trace)

    # Create node traces by category
    node_traces = []

    for category, nodes in categories.items():
        if not nodes:
            continue

        node_x = []
        node_y = []
        node_text = []

        for node in nodes:
            if node in pos:
                x, y = pos[node]
                node_x.append(x)
                node_y.append(y)

                # Create hover text
                in_edges = list(G.in_edges(node))
                out_edges = list(G.out_edges(node))
                hover_text = f"<b>{node}</b><br>"
                hover_text += f"Category: {category}<br>"
                hover_text += f"In-transitions: {len(in_edges)}<br>"
                hover_text += f"Out-transitions: {len(out_edges)}<br>"

                if in_edges:
                    from_states = [e[0] for e in in_edges[:3]]
                    hover_text += f"From: {', '.join(from_states)}"
                    if len(in_edges) > 3:
                        hover_text += f" (+{len(in_edges)-3} more)"
                    hover_text += "<br>"

                if out_edges:
                    to_states = [e[1] for e in out_edges[:3]]
                    hover_text += f"To: {', '.join(to_states)}"
                    if len(out_edges) > 3:
                        hover_text += f" (+{len(out_edges)-3} more)"

                node_text.append(hover_text)

        if node_x:  # Only create trace if there are nodes
            node_trace = go.Scatter(
                x=node_x,
                y=node_y,
                mode='markers+text',
                name=category.title(),
                text=[n for n in nodes if n in pos],
                textposition='top center',
                textfont=dict(size=8),
                hovertext=node_text,
                hoverinfo='text',
                marker=dict(
                    color=get_color_for_category(category),
                    size=12,
                    line=dict(width=2, color='white')
                )
            )
            node_traces.append(node_trace)

    # Create figure
    fig = go.Figure(
        data=edge_traces + node_traces,
        layout=go.Layout(
            title=dict(
                text='Software Factory 2.0 - Complete State Machine',
                font=dict(size=20)
            ),
            showlegend=True,
            hovermode='closest',
            margin=dict(b=20, l=5, r=5, t=40),
            width=2400,
            height=1800,
            xaxis=dict(showgrid=False, zeroline=False, showticklabels=False),
            yaxis=dict(showgrid=False, zeroline=False, showticklabels=False),
            plot_bgcolor='white',
            legend=dict(
                orientation="v",
                yanchor="top",
                y=1,
                xanchor="left",
                x=1.02,
                bgcolor="rgba(255,255,255,0.9)",
                bordercolor="black",
                borderwidth=1
            )
        )
    )

    # Add arrows for directed edges
    annotations = []
    for edge in G.edges():
        x0, y0 = pos[edge[0]]
        x1, y1 = pos[edge[1]]

        # Calculate arrow position (slightly before the target)
        dx = x1 - x0
        dy = y1 - y0
        norm = np.sqrt(dx**2 + dy**2)
        if norm > 0:
            dx /= norm
            dy /= norm

            # Position arrow head
            arrow_x = x1 - dx * 0.3
            arrow_y = y1 - dy * 0.3

            annotations.append(
                dict(
                    ax=x0 + dx * 0.3,
                    ay=y0 + dy * 0.3,
                    x=arrow_x,
                    y=arrow_y,
                    xref='x',
                    yref='y',
                    axref='x',
                    ayref='y',
                    showarrow=True,
                    arrowhead=2,
                    arrowsize=0.8,
                    arrowwidth=0.5,
                    arrowcolor='#888',
                    opacity=0.6
                )
            )

    # Limit annotations to prevent clutter
    if len(annotations) <= 200:  # Only show arrows if not too many
        fig.update_layout(annotations=annotations[:200])

    return fig

def create_simplified_graph(G, pos, categories):
    """Create a simplified version focusing on main paths"""

    # Identify critical paths
    critical_states = ['INIT', 'WAVE_START', 'CREATE_NEXT_INFRASTRUCTURE',
                      'SPAWN_SW_ENGINEERS', 'MONITORING_SWE_PROGRESS', 'WAVE_COMPLETE',
                      'INTEGRATE_WAVE_EFFORTS', 'COMPLETE_PHASE', 'PROJECT_DONE', 'ERROR_RECOVERY']

    # Filter to show only critical states and their immediate connections
    subgraph_nodes = set()
    for state in critical_states:
        if state in G.nodes():
            subgraph_nodes.add(state)
            # Add immediate neighbors
            subgraph_nodes.update(G.predecessors(state))
            subgraph_nodes.update(G.successors(state))

    subgraph = G.subgraph(subgraph_nodes)

    # Create traces for simplified graph
    edge_traces = []

    for edge in subgraph.edges():
        if edge[0] in pos and edge[1] in pos:
            x0, y0 = pos[edge[0]]
            x1, y1 = pos[edge[1]]

            edge_trace = go.Scatter(
                x=[x0, x1, None],
                y=[y0, y1, None],
                mode='lines',
                line=dict(width=1.5, color='#888'),
                hoverinfo='none',
                showlegend=False
            )
            edge_traces.append(edge_trace)

    # Create node traces
    node_traces = []

    for category, nodes in categories.items():
        category_nodes = [n for n in nodes if n in subgraph_nodes and n in pos]
        if not category_nodes:
            continue

        node_x = []
        node_y = []
        node_text = []
        node_labels = []

        for node in category_nodes:
            x, y = pos[node]
            node_x.append(x)
            node_y.append(y)

            # Shorten long names for display
            label = node.replace('_', ' ')
            if len(label) > 20:
                label = label[:17] + '...'
            node_labels.append(label)

            # Create hover text
            hover_text = f"<b>{node}</b><br>"
            hover_text += f"Category: {category}<br>"
            hover_text += f"Connections: {subgraph.degree(node)}"
            node_text.append(hover_text)

        if node_x:
            node_trace = go.Scatter(
                x=node_x,
                y=node_y,
                mode='markers+text',
                name=category.title(),
                text=node_labels,
                textposition='top center',
                textfont=dict(size=10, color='black'),
                hovertext=node_text,
                hoverinfo='text',
                marker=dict(
                    color=get_color_for_category(category),
                    size=20,
                    line=dict(width=2, color='white')
                )
            )
            node_traces.append(node_trace)

    # Create figure for simplified view
    fig = go.Figure(
        data=edge_traces + node_traces,
        layout=go.Layout(
            title=dict(
                text='Software Factory 2.0 - Simplified State Flow',
                font=dict(size=20)
            ),
            showlegend=True,
            hovermode='closest',
            margin=dict(b=20, l=5, r=5, t=40),
            width=1600,
            height=1200,
            xaxis=dict(showgrid=False, zeroline=False, showticklabels=False),
            yaxis=dict(showgrid=False, zeroline=False, showticklabels=False),
            plot_bgcolor='white',
            legend=dict(
                orientation="v",
                yanchor="top",
                y=1,
                xanchor="left",
                x=1.02,
                bgcolor="rgba(255,255,255,0.9)",
                bordercolor="black",
                borderwidth=1
            )
        )
    )

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

    print("Creating layout...")
    pos = create_layout(G)

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

    print("\nCreating complete visualization...")
    fig_complete = create_plotly_graph(G, pos, categories)

    print("Creating simplified visualization...")
    fig_simple = create_simplified_graph(G, pos, categories)

    print("\nSaving visualizations...")

    # Save as static PNG images
    print("  - Saving complete graph as PNG...")
    pio.write_image(fig_complete, 'software-factory-state-machine-complete.png',
                   width=2400, height=1800, scale=2)

    print("  - Saving simplified graph as PNG...")
    pio.write_image(fig_simple, 'software-factory-state-machine-simplified.png',
                   width=1600, height=1200, scale=2)

    # Also save interactive HTML versions
    print("  - Saving interactive HTML versions...")
    fig_complete.write_html('software-factory-state-machine-complete.html')
    fig_simple.write_html('software-factory-state-machine-simplified.html')

    print("\n✅ Visualization complete!")
    print("\nGenerated files:")
    print("  📊 software-factory-state-machine-complete.png - Full state machine")
    print("  📊 software-factory-state-machine-simplified.png - Simplified main flow")
    print("  🌐 software-factory-state-machine-complete.html - Interactive full version")
    print("  🌐 software-factory-state-machine-simplified.html - Interactive simplified version")

    # Save statistics to file
    with open('state-machine-statistics.txt', 'w') as f:
        f.write("Software Factory 2.0 - State Machine Statistics\n")
        f.write("=" * 50 + "\n\n")
        f.write(f"Total States: {stats['total_states']}\n")
        f.write(f"Total Transitions: {stats['total_transitions']}\n")
        f.write(f"Number of Agents: {stats['agents']}\n\n")

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

    print("  📄 state-machine-statistics.txt - Detailed statistics")

if __name__ == "__main__":
    main()