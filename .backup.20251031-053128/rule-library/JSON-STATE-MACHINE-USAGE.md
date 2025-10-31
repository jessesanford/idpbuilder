# JSON State Machine Usage Guide

## File Locations
- Main: `$CLAUDE_PROJECT_DIR/state-machines/software-factory-3.0-state-machine.json`
- Sub-machines: `$CLAUDE_PROJECT_DIR/state-machines/[agent].json`

## Common Operations

### Check if state exists
```bash
jq -e '.agents.orchestrator.states.INIT' software-factory-3.0-state-machine.json
```

### Get valid transitions
```bash
jq -r '.transition_matrix.orchestrator.INIT[]' software-factory-3.0-state-machine.json
```

### List all states
```bash
jq -r '.agents.orchestrator.states | keys[]' software-factory-3.0-state-machine.json
```

### Validate transition
```bash
STATE="INIT"
TARGET="WAVE_START"
jq -e ".transition_matrix.orchestrator.$STATE | index(\"$TARGET\")" software-factory-3.0-state-machine.json
```

## Migration Notes
Replace grep patterns with jq queries when working with the state machine.
