# SOFTWARE FACTORY 2.0 - STATE MACHINE VISUALIZATION TOOLS

## Overview

These tools provide ASCII art visualization of the Software Factory 2.0 state machine, helping developers understand:
- Current state position
- State transitions and flow
- Parallel execution paths
- Integration and fix cascade flows

## Tools Available

### 1. visualize-state-machine.py

Provides a detailed view of the current state machine position with context.

```bash
# Full visualization with current state highlighted
python3 tools/visualize-state-machine.py

# Compact view for quick status checks
python3 tools/visualize-state-machine.py --compact

# Without colors (for logging/documentation)
python3 tools/visualize-state-machine.py --no-color
```

**Features:**
- Shows current state with highlighting
- Displays state path from INIT
- Lists possible next transitions
- Color-coded state types (spawn, waiting, action, etc.)
- Shows required actions and entry conditions
- Includes timestamp and phase/wave info

### 2. visualize-state-flow.py

Provides flow diagrams showing the complete state machine architecture.

```bash
# Complete flow diagram
python3 tools/visualize-state-flow.py

# Specific flow sections:
python3 tools/visualize-state-flow.py --wave        # Wave execution flow
python3 tools/visualize-state-flow.py --integration # Integration flow
python3 tools/visualize-state-flow.py --fix         # Fix cascade flow
python3 tools/visualize-state-flow.py --split       # Split handling flow
```

**Features:**
- Wave execution flow with parallelization
- Integration flow (wave, phase, project)
- Fix cascade flow with loops
- Split handling for large efforts
- Key rules and constraints reminder

## Output Examples

### Current State View
```
================================================================================
SOFTWARE FACTORY 2.0 - STATE MACHINE VISUALIZATION
================================================================================

Current State: MONITORING
Previous State: SPAWN_SW_ENGINEERS
Phase: 2, Wave: 1

--------------------------------------------------------------------------------
State Flow:
--------------------------------------------------------------------------------

[INIT]
  │
  ▼
[WAVE_START]
  │
  ▼
[CREATE_NEXT_INFRASTRUCTURE]
  │
  ▼
[SPAWN_SW_ENGINEERS]
  │
  ▼
╔══════════════════════════════════════╗
║             MONITORING               ║  ◄── YOU ARE HERE
╚══════════════════════════════════════╝

Possible Next States:
  → SPAWN_CODE_REVIEWERS_EFFORT_REVIEW
  → ERROR_RECOVERY
  → MONITORING
```

### Wave Flow Diagram
```
WAVE EXECUTION FLOW
======================================================================

┌─────────────┐
│  WAVE_START │
└──────┬──────┘
       │
       ▼
┌─────────────────────────────────┐
│ CREATE_NEXT_INFRASTRUCTURE      │
└──────┬──────────────────────────┘
       │
       ▼
┌─────────────────────────────────────────┐
│ SPAWN_SW_ENGINEERS                             │
│ ┌─────────┐ ┌─────────┐ ┌─────────┐    │
│ │ SWE-1   │ │ SWE-2   │ │ SWE-N   │    │ (Parallel)
│ └─────────┘ └─────────┘ └─────────┘    │
└──────┬──────────────────────────────────┘
       │
       ▼
┌─────────────────────────────────┐
│ MONITORING                       │
└─────────────────────────────────┘
```

## Color Legend

When colors are enabled:
- **GREEN**: Spawn states (agent creation)
- **YELLOW**: Waiting states (monitoring/waiting)
- **CYAN**: Action states (active operations)
- **RED**: Terminal states (endpoints)
- **MAGENTA**: Initial states (starting points)
- **BLUE**: Checkpoint states (user review required)

## Integration with Other Tools

These visualization tools work with:
- `orchestrator-state-v3.json` - Current state tracking
- `state-machines/software-factory-3.0-state-machine.json` - State definitions
- `tools/claude_orchestrator_runner.py` - State execution

## Use Cases

1. **Debugging State Transitions**
   - Verify current position in state machine
   - Check valid next transitions
   - Understand why certain transitions are blocked

2. **Understanding Parallelization**
   - See which agents can run in parallel
   - Understand timing requirements
   - Visualize wave execution patterns

3. **Documentation**
   - Generate diagrams for documentation
   - Explain state machine to new developers
   - Create training materials

4. **Monitoring Progress**
   - Quick status checks during execution
   - Understanding current phase/wave
   - Tracking sub-state machines

## Requirements

- Python 3.6+
- No external dependencies (uses only standard library)
- Terminal with UTF-8 support for box drawing characters
- Optional: Terminal with color support for enhanced visualization

## Troubleshooting

### No State Found
If the visualizer can't find the current state:
1. Check that `orchestrator-state-v3.json` exists
2. Verify the state name matches state machine definition
3. Use the example file if needed

### Missing State Machine
If the state machine can't be loaded:
1. Check `state-machines/software-factory-3.0-state-machine.json` exists
2. Verify JSON is valid (no syntax errors)
3. Check file permissions

### Character Display Issues
If box characters don't display correctly:
1. Ensure terminal supports UTF-8
2. Try `--no-color` option
3. Check terminal font supports box drawing characters

## Contributing

When adding new states or transitions:
1. Update the state machine JSON
2. Test visualization tools still work
3. Add any new state types to color coding
4. Update this README if needed