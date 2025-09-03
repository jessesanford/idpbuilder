# ORCHESTRATOR STATE MACHINE - ASCII DIAGRAM

## Complete State Flow with Integration Feedback Loop

```
┌─────────────┐
│    START    │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│    INIT     │ (Load config, validate environment)
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  PLANNING   │ (Create/load project plan)
└──────┬──────┘
       │
       ▼
┌──────────────────────┐
│  PHASE LOOP START    │
└──────────┬───────────┘
           │
           ▼
    ┌─────────────┐
    │ PHASE_INIT  │ (Start new phase)
    └──────┬──────┘
           │
           ▼
    ┌──────────────────────┐
    │   WAVE LOOP START    │
    └──────────┬───────────┘
               │
               ▼
        ┌─────────────┐
        │  WAVE_INIT  │ (Start new wave)
        └──────┬──────┘
               │
               ▼
        ┌──────────────────────────┐
        │ SETUP_EFFORT_INFRASTRUCTURE│ (Create branches, worktrees)
        └──────────┬────────────────┘
                   │
                   ▼
        ┌─────────────┐
        │ SPAWN_AGENTS│ (Deploy SWEs, Reviewers)
        └──────┬──────┘
               │
               ▼
        ┌─────────────┐
        │   MONITOR   │ (Track progress, check sizes)
        └──────┬──────┘
               │
               ├─[Issues Found]──────┐
               │                     ▼
               │              ┌──────────────┐
               │              │ERROR_RECOVERY│
               │              └──────┬───────┘
               │                     │
               │◄────────────────────┘
               │
               ├─[Size Violation]────┐
               │                     ▼
               │              ┌──────────────┐
               │              │HANDLE_SPLITS │
               │              └──────┬───────┘
               │                     │
               │◄────────────────────┘
               │
               └─[All Complete]──────┐
                                     ▼
        ┌──────────────┐
        │ WAVE_COMPLETE│ (Merge to wave branch)
        └──────┬───────┘
               │
               ▼
        ┌──────────────┐
        │ WAVE_REVIEW  │ (Architect reviews wave)
        └──────┬───────┘
               │
               ├─[More Waves]────────┐
               │                     │
               │◄────────────────────┘
               │
               └─[Waves Done]────────┐
                                     ▼
    ┌───────────────┐
    │PHASE_COMPLETE │ (Merge to phase branch)
    └──────┬────────┘
           │
           ▼
    ┌───────────────┐
    │ PHASE_REVIEW  │ (Architect phase assessment)
    └──────┬────────┘
           │
           ├─[ON_TRACK]──────────────┐
           │                         ▼
           │                  ┌──────────────┐
           │                  │ NEXT_PHASE   │
           │                  └──────┬───────┘
           │                         │
           │◄────────────────────────┘
           │
           ├─[NEEDS_ADJUSTMENT]──────┐
           │                         ▼
           │                  ┌──────────────────────┐
           │                  │ REVIEW_PLAN_DEVIATIONS│
           │                  └──────────┬───────────┘
           │                             │
           │                             ▼
           │                  ┌──────────────────────┐
           │                  │ ADJUSTMENT_PLANNING  │
           │                  └──────────┬───────────┘
           │                             │
           │                             ▼
           │                  ┌──────────────────────┐
           │                  │ ARCHITECT_FEEDBACK   │
           │                  └──────────┬───────────┘
           │                             │
           │◄────────────────────────────┘
           │
           └─[OFF_TRACK]─────────────┐
                                     ▼
           ┌──────────────────────────────────┐
           │        INTEGRATION_REVIEW        │
           └─────────────┬────────────────────┘
                         │
                         ▼
           ┌──────────────────────────────────┐
           │    ANALYZE_PROJECT_DEVIATIONS    │
           └─────────────┬────────────────────┘
                         │
                         ▼
           ┌──────────────────────────────────┐
           │     CREATE_RECOVERY_STRATEGY     │
           └─────────────┬────────────────────┘
                         │
                         ▼
           ┌──────────────────────────────────┐
           │      ARCHITECT_INTEGRATION       │
           └─────────────┬────────────────────┘
                         │
                         ├─[Recovery Approved]─┐
                         │                     ▼
                         │              ┌──────────────┐
                         │              │EXECUTE_RECOVERY│
                         │              └──────┬───────┘
                         │                     │
                         │◄────────────────────┘
                         │
                         └─[Project Reset]─────┐
                                               ▼
                                        ┌──────────────┐
                                        │ RESET_PROJECT│
                                        └──────┬───────┘
                                               │
                                               ▼
                                        ┌──────────────┐
                                        │   PLANNING   │
                                        └──────────────┘

## Integration Feedback States (Detailed)

### INTEGRATION_REVIEW
```
┌────────────────────────────────────────┐
│         INTEGRATION_REVIEW             │
│                                        │
│ Triggers:                              │
│ - Phase marked OFF_TRACK               │
│ - Critical architecture violations     │
│ - Accumulated technical debt > threshold│
│                                        │
│ Actions:                               │
│ 1. Stop all active work                │
│ 2. Merge outstanding branches          │
│ 3. Create integration checkpoint       │
│ 4. Spawn architect for deep review     │
└────────────────────────────────────────┘
```

### ANALYZE_PROJECT_DEVIATIONS
```
┌────────────────────────────────────────┐
│    ANALYZE_PROJECT_DEVIATIONS         │
│                                        │
│ Actions:                               │
│ 1. Compare actual vs planned:         │
│    - Architecture alignment            │
│    - Feature completeness              │
│    - Code quality metrics              │
│    - Test coverage                     │
│ 2. Identify root causes                │
│ 3. Quantify technical debt             │
│ 4. Document deviation patterns         │
│                                        │
│ Outputs:                               │
│ - DEVIATION_ANALYSIS.md                │
│ - Technical debt inventory             │
│ - Risk assessment matrix               │
└────────────────────────────────────────┘
```

### CREATE_RECOVERY_STRATEGY
```
┌────────────────────────────────────────┐
│     CREATE_RECOVERY_STRATEGY          │
│                                        │
│ Actions:                               │
│ 1. Generate recovery options:         │
│    - Minimal: Fix critical only       │
│    - Standard: Fix + refactor         │
│    - Complete: Reset to checkpoint    │
│ 2. Estimate effort for each           │
│ 3. Risk/benefit analysis              │
│ 4. Create detailed recovery plan      │
│                                        │
│ Outputs:                               │
│ - RECOVERY_STRATEGY.md                 │
│ - Effort estimates                     │
│ - Success criteria                     │
└────────────────────────────────────────┘
```

### ARCHITECT_INTEGRATION
```
┌────────────────────────────────────────┐
│      ARCHITECT_INTEGRATION            │
│                                        │
│ Decision Points:                      │
│ ┌──────────────────────────┐         │
│ │ Recovery feasible?        │         │
│ │ - Within time budget?     │         │
│ │ - Addresses root causes?  │         │
│ │ - Sustainable solution?   │         │
│ └──────────────────────────┘         │
│           │                            │
│           ├─[YES]→ EXECUTE_RECOVERY   │
│           │                            │
│           └─[NO]→ RESET_PROJECT       │
│                                        │
│ Approval Criteria:                     │
│ - Clear success metrics                │
│ - Defined completion criteria          │
│ - Risk mitigation plan                 │
└────────────────────────────────────────┘
```

### EXECUTE_RECOVERY
```
┌────────────────────────────────────────┐
│         EXECUTE_RECOVERY               │
│                                        │
│ Sub-states:                            │
│ 1. RECOVERY_PLANNING                   │
│ 2. RECOVERY_IMPLEMENTATION             │
│ 3. RECOVERY_VALIDATION                 │
│ 4. RECOVERY_INTEGRATION                │
│                                        │
│ Success → Return to PHASE_INIT        │
│ Failure → INTEGRATION_REVIEW (retry)   │
└────────────────────────────────────────┘
```

## State Transition Rules

### Normal Flow
- INIT → PLANNING → PHASE_INIT → WAVE_INIT → ... → COMPLETE

### Error Recovery Flows
- MONITOR → ERROR_RECOVERY → MONITOR (retry)
- MONITOR → HANDLE_SPLITS → MONITOR (continue)

### Integration Feedback Flows
- PHASE_REVIEW → [ON_TRACK] → NEXT_PHASE
- PHASE_REVIEW → [NEEDS_ADJUSTMENT] → ADJUSTMENT_PLANNING → PHASE_INIT
- PHASE_REVIEW → [OFF_TRACK] → INTEGRATION_REVIEW → ...

### Critical Decision Points
1. **Phase Review Assessment**
   - ON_TRACK: Continue to next phase
   - NEEDS_ADJUSTMENT: Minor corrections, continue
   - OFF_TRACK: Major intervention required

2. **Integration Review Decision**
   - EXECUTE_RECOVERY: Attempt to fix issues
   - RESET_PROJECT: Start fresh with new plan

3. **Recovery Validation**
   - SUCCESS: Return to normal flow
   - FAILURE: Re-enter integration review

## State Machine Rules

### R206: State Machine Validation
- All transitions must be validated
- State file must be updated atomically
- Previous state must be valid source
- Target state must be reachable

### R281: Complete State Initialization
- All state fields must be populated
- No partial state transitions allowed
- State consistency checks required

### Integration States Requirements
- Must preserve all work before entering
- Must create checkpoint branches
- Must document all decisions
- Must track recovery metrics

## Error Handling States

### ERROR_RECOVERY
- Handles: Build failures, test failures, merge conflicts
- Max retries: 3
- Escalation: INTEGRATION_REVIEW after max retries

### HANDLE_SPLITS
- Triggers: Size violations (>700 lines soft, >800 hard)
- Creates: Split plan, sequential execution
- Returns: To MONITOR after completion

## Terminal States

### COMPLETE
- All phases complete
- All integration tests passing
- Project delivered

### ABANDONED
- Project cancelled
- Unrecoverable failure
- Resource exhaustion

## Monitoring Requirements

Throughout all states:
- TODO persistence (R187-R190)
- Size compliance checks
- State file updates
- Progress tracking
- Error detection

## State File Example

```yaml
current_state: MONITOR
previous_state: SPAWN_AGENTS
current_phase: 1
current_wave: 2
phase_name: "Core Implementation"
wave_name: "API Development"
efforts_in_progress:
  - name: "User Authentication"
    branch: "wave-1-2-auth"
    agent: "sw-engineer-1"
    status: "implementing"
efforts_completed:
  - name: "Database Schema"
    branch: "wave-1-1-schema"
    lines: 450
    review_status: "approved"
integration_status:
  last_review: "2025-01-20T14:30:00Z"
  health: "ON_TRACK"
  technical_debt: 12
  recovery_attempts: 0
```

## Critical Transitions

### Entering Integration Review
```
PRE-CONDITIONS:
- Phase marked OFF_TRACK
- OR technical_debt > threshold
- OR critical_violations > 0

ACTIONS:
1. Save all work
2. Create integration branch
3. Update state file
4. Spawn architect
5. Begin analysis
```

### Exiting Integration Review
```
POST-CONDITIONS:
- Recovery plan approved
- OR project reset completed
- State file updated
- All agents notified
- Checkpoints created
```

## Usage Notes

1. **State Validation**: Always validate transitions with R206
2. **Atomic Updates**: State changes must be atomic
3. **Persistence**: Save state before/after transitions
4. **Recovery**: Always create checkpoints before risky operations
5. **Monitoring**: Check state health continuously

This ASCII diagram represents the complete orchestrator state machine including the new integration feedback loop for handling OFF_TRACK scenarios and major project corrections.