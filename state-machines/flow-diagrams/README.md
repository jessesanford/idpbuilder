# Software Factory 2.0 - Flow Diagrams

This directory contains auto-generated flow diagrams of the Software Factory 2.0 state machine.

## Generated Diagrams

### Overall System Flow
- **File**: `overall-system-flow.png`
- **Description**: Complete system overview showing all agents and their states grouped by agent, with transitions between states

### Agent-Specific Flows
- **orchestrator-flow.png**: Orchestrator agent state flow
- **sw-engineer-flow.png**: Software Engineer agent state flow
- **code-reviewer-flow.png**: Code Reviewer agent state flow
- **architect-flow.png**: Architect agent state flow

### Specialized Views
- **phase-progression-flow.png**: Shows phase and wave progression through the system
- **error-recovery-flow.png**: Highlights error states and recovery paths
- **split-implementation-flow.png**: Details the split implementation process for large efforts

## How to Regenerate

To regenerate all diagrams:
```bash
python3 tools/generate-flow-diagrams.py
```

To generate specific diagrams:
```bash
# Generate only orchestrator flow
python3 tools/generate-flow-diagrams.py --agent orchestrator

# Generate only error recovery flow
python3 tools/generate-flow-diagrams.py --type error
```

## Visual Legend

### Colors
- **Light Green**: INIT states
- **Sky Blue**: PLANNING states
- **Gold**: IMPLEMENTATION states
- **Plum**: REVIEW states
- **Light Red**: ERROR states
- **Pale Green**: PROJECT_DONE/COMPLETE states
- **Orange**: INTEGRATE_WAVE_EFFORTS states
- **Light Pink**: SPLIT states
- **Light Sea Green**: WAVE states
- **Steel Blue**: PHASE states

### Shapes
- **Ellipse**: INIT states
- **Octagon**: ERROR states
- **Double Circle**: PROJECT_DONE/COMPLETE states
- **Diamond**: Decision/Check states
- **Box**: Default states

### Edge Styles
- **Red Dashed**: Error transitions
- **Green Bold**: Success transitions
- **Blue Dotted**: Self-loops
- **Black**: Normal transitions

## Requirements

- Python 3.x
- graphviz package: `pip install graphviz`
- Graphviz system package: `sudo apt-get install graphviz`

## Notes

The diagrams are generated from `state-machines/software-factory-3.0-state-machine.json` which contains the complete state machine definition for all agents in the Software Factory 2.0 system.