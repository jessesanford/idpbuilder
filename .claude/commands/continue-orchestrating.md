# Continue Orchestrating Command

## Purpose
This command resumes or initiates the software factory orchestration process, coordinating multiple agents through phases, waves, and efforts to deliver complete software implementation.

## Usage
```
/continue-orchestrating
```

## What This Command Does

1. **Spawns the Orchestrator**: Activates @agent-orchestrator-task-master
2. **Loads Current State**: Reads orchestrator-state.yaml to understand progress
3. **Resumes State Machine**: Continues from the last known state
4. **Coordinates Agents**: Spawns appropriate agents for current tasks
5. **Enforces Quality Gates**: Ensures all checks and reviews are completed

## Prerequisites

### Required Directory Structure
```
/home/vscode/workspaces/idpbuilder/
├── core/
│   └── SOFTWARE-FACTORY-STATE-MACHINE.md
├── protocols/
│   └── [protocol files]
├── phase-plans/
│   └── [phase-specific implementation plans]
├── orchestrator-state.yaml
├── PROJECT-IMPLEMENTATION-PLAN.md
├── efforts/
├── tools/
└── todos/
```

### Required Files
- core/SOFTWARE-FACTORY-STATE-MACHINE.md
- orchestrator-state.yaml
- PROJECT-IMPLEMENTATION-PLAN.md
- phase-plans/PHASE{X}-SPECIFIC-IMPL-PLAN.md files
- protocols/[various protocol files]
- .claude/agents/[agent configurations]

## Orchestrator Workflow

### Initial Invocation
```markdown
Task @agent-orchestrator-task-master:

STARTUP REQUIREMENTS:
1. Print timestamp and acknowledgment
2. Verify environment
3. Read core/SOFTWARE-FACTORY-STATE-MACHINE.md

INITIAL TASK:
1. CHECK: Does orchestrator-state.yaml exist?
   - YES: Load state and resume from current position
   - NO: Initialize from PROJECT-IMPLEMENTATION-PLAN.md
2. IDENTIFY: Current state in state machine
3. EXECUTE: Appropriate state transition
4. CONTINUE: Until reaching terminal state (SUCCESS or STOP)
```

### State Machine Execution

The orchestrator will:
1. **Determine Current State**: From orchestrator-state.yaml
2. **Execute State Logic**: Per SOFTWARE-FACTORY-STATE-MACHINE.md
3. **Spawn Agents**: As needed for current state
4. **Update State File**: Record all progress
5. **Continue Transitions**: Until terminal state

## Agent Spawning Templates

### Code Reviewer - Planning
```markdown
Task @agent-code-reviewer:

PURPOSE: Create implementation plan for Effort E{X}.{Y}.{Z}

CONTEXT:
- Working directory: /workspaces/efforts/phase{X}/wave{Y}/effort{Z}
- Requirements: [From phase plan]

DELIVERABLES:
- IMPLEMENTATION-PLAN.md
- work-log.md
```

### SW Engineer - Implementation
```markdown
Task @agent-sw-engineer:

PURPOSE: Implement Effort E{X}.{Y}.{Z}

INSTRUCTIONS:
1. READ: IMPLEMENTATION-PLAN.md
2. IMPLEMENT: Per plan
3. MEASURE: Size continuously
4. TEST: Per requirements
5. COMMIT: When complete

SIZE LIMIT: {configured_limit} lines
```

### Architect - Review
```markdown
Task @agent-architect-reviewer:

PURPOSE: Review Wave {X}.{Y} completion

ASSESS:
- Architectural consistency
- Integration readiness
- Technical debt

DECISION: PROCEED / CHANGES_REQUIRED / STOP
```

## Quality Gates

### Enforced Automatically
- **Size Limits**: No effort exceeds configured limit
- **Review Completion**: Every effort reviewed
- **Test Coverage**: Meets project requirements
- **Architectural Approval**: At wave and phase boundaries

### Recovery Procedures
- **Size Violation**: Automatic split protocol
- **Review Failure**: Fix and re-review cycle
- **Architecture Issues**: Addendum and correction

## State File Management

### orchestrator-state.yaml
```yaml
current_phase: X
current_wave: Y
current_state: "STATE_NAME"

efforts_completed:
  - [list of completed efforts]

efforts_in_progress:
  - [current active efforts]

efforts_pending:
  - [queued efforts]

reviews:
  - [review outcomes]
```

### Update Frequency
- After every state transition
- When effort status changes
- After review completions
- Before spawning agents

## Troubleshooting

### Common Issues

| Issue | Resolution |
|-------|------------|
| "State file not found" | Initialize with PROJECT-IMPLEMENTATION-PLAN.md |
| "Stuck in state" | Check efforts_in_progress for blockers |
| "Agent confusion" | Verify working directories and branches |
| "Size violation" | Implement split protocol immediately |

### Recovery Commands
```bash
# Check current state
cat orchestrator-state.yaml

# Verify working directories
ls -la /workspaces/efforts/

# Check for blocking reviews
grep "in_progress" orchestrator-state.yaml
```

## Best Practices

1. **Commit State Frequently**: Preserve progress
2. **Monitor Agent Output**: Catch issues early
3. **Review State Transitions**: Ensure compliance
4. **Document Deviations**: In state file notes
5. **Clear Blockers Quickly**: Don't let efforts stall

## Terminal States

### SUCCESS
- All phases complete
- All tests passing
- Feature coverage 100%
- Final review accepted

### STOP
- OFF_TRACK assessment
- Critical architecture violation
- Unrecoverable error
- Manual intervention required

## Notes
- Orchestrator NEVER writes code
- All splits are sequential
- Reviews are mandatory
- Gates cannot be skipped
- State machine is law