# Orchestrator State Directory Map

## Complete List of Valid Orchestrator States

Each state MUST have a corresponding directory with rules.md file.

### ✅ States with Rules Files:

1. **INIT** - `/agent-states/orchestrator/INIT/rules.md`
   - Initial startup and state detection

2. **PLANNING** - `/agent-states/orchestrator/PLANNING/rules.md`
   - Planning phases and waves

3. **SETUP_EFFORT_INFRASTRUCTURE** - `/agent-states/orchestrator/SETUP_EFFORT_INFRASTRUCTURE/rules.md`
   - Creating effort directories and clones

4. **ANALYZE_CODE_REVIEWER_PARALLELIZATION** - `/agent-states/orchestrator/ANALYZE_CODE_REVIEWER_PARALLELIZATION/rules.md`
   - MANDATORY - Analyzing wave plan to determine Code Reviewer spawn strategy

5. **SPAWN_CODE_REVIEWERS_EFFORT_PLANNING** - `/agent-states/orchestrator/SPAWN_CODE_REVIEWERS_EFFORT_PLANNING/rules.md`
   - Spawning code reviewers to create implementation plans

6. **WAITING_FOR_EFFORT_PLANS** - `/agent-states/orchestrator/WAITING_FOR_EFFORT_PLANS/rules.md`
   - Waiting for code reviewers to complete plans

7. **ANALYZE_IMPLEMENTATION_PARALLELIZATION** - `/agent-states/orchestrator/ANALYZE_IMPLEMENTATION_PARALLELIZATION/rules.md`
   - MANDATORY - Analyzing effort plans to determine SW Engineer spawn strategy

8. **SPAWN_AGENTS** - `/agent-states/orchestrator/SPAWN_AGENTS/rules.md`
   - Spawning SW engineers for implementation

9. **MONITOR** - `/agent-states/orchestrator/MONITOR/rules.md`
   - Monitoring agent progress

10. **WAVE_COMPLETE** - `/agent-states/orchestrator/WAVE_COMPLETE/rules.md`
    - All wave efforts completed and reviewed

11. **INTEGRATION** - `/agent-states/orchestrator/INTEGRATION/rules.md`
    - Creating integration branches in target repo

12. **WAVE_REVIEW** - `/agent-states/orchestrator/WAVE_REVIEW/rules.md`
    - Architect reviewing wave integration

13. **ERROR_RECOVERY** - `/agent-states/orchestrator/ERROR_RECOVERY/rules.md`
    - Handling errors and failures

14. **SPAWN_ARCHITECT_PHASE_ASSESSMENT** - `/agent-states/orchestrator/SPAWN_ARCHITECT_PHASE_ASSESSMENT/rules.md`
    - Request architect to assess complete phase (last wave only)

15. **WAITING_FOR_PHASE_ASSESSMENT** - `/agent-states/orchestrator/WAITING_FOR_PHASE_ASSESSMENT/rules.md`
    - Waiting for architect phase assessment decision

16. **PHASE_COMPLETE** - `/agent-states/orchestrator/PHASE_COMPLETE/rules.md`
    - Phase assessment passed, handling phase-level integration

17. **SUCCESS** - `/agent-states/orchestrator/SUCCESS/rules.md`
    - Project successfully completed (only after phase assessment)

18. **HARD_STOP** - `/agent-states/orchestrator/HARD_STOP/rules.md`
    - Critical failure requiring manual intervention

## State Transition Flow

```
INIT 
  ↓
PLANNING
  ↓
SETUP_EFFORT_INFRASTRUCTURE
  ↓
ANALYZE_CODE_REVIEWER_PARALLELIZATION  ← MANDATORY GATE
  ↓
SPAWN_CODE_REVIEWERS_EFFORT_PLANNING
  ↓
WAITING_FOR_EFFORT_PLANS
  ↓
ANALYZE_IMPLEMENTATION_PARALLELIZATION  ← MANDATORY GATE
  ↓
SPAWN_AGENTS
  ↓
MONITOR
  ↓
WAVE_COMPLETE
  ↓
INTEGRATION
  ↓
WAVE_REVIEW
  ↓
(Next wave PLANNING or SPAWN_ARCHITECT_PHASE_ASSESSMENT if last wave)
  ↓
SPAWN_ARCHITECT_PHASE_ASSESSMENT (last wave only)
  ↓
WAITING_FOR_PHASE_ASSESSMENT
  ↓
PHASE_COMPLETE
  ↓
SUCCESS
```

## Common State Transition Errors

1. **Missing State Directory**: If transitioning to a state without a rules.md file
   - Solution: Create the directory and rules.md file

2. **Invalid State Name**: Transitioning to a state not in the list
   - Solution: Check spelling and use exact state names

3. **Wrong Agent State**: Using a state from another agent type
   - Example: IMPLEMENTATION is for SW engineers, not orchestrator

## Verification Command

```bash
# Verify all state directories exist
for state in INIT PLANNING SETUP_EFFORT_INFRASTRUCTURE \
  ANALYZE_CODE_REVIEWER_PARALLELIZATION \
  SPAWN_CODE_REVIEWERS_EFFORT_PLANNING WAITING_FOR_EFFORT_PLANS \
  ANALYZE_IMPLEMENTATION_PARALLELIZATION \
  SPAWN_AGENTS MONITOR WAVE_COMPLETE INTEGRATION WAVE_REVIEW \
  SPAWN_ARCHITECT_PHASE_ASSESSMENT WAITING_FOR_PHASE_ASSESSMENT \
  PHASE_COMPLETE ERROR_RECOVERY SUCCESS HARD_STOP; do
  
  if [ -f "agent-states/orchestrator/$state/rules.md" ]; then
    echo "✅ $state"
  else
    echo "❌ $state - MISSING!"
  fi
done
```

## R217 Compliance

When transitioning to any state, the orchestrator MUST:
1. Update state file (R288)
2. Read the corresponding rules.md file (R217)
3. Acknowledge the rules before proceeding

Example:
```bash
# Transitioning to WAVE_REVIEW
update_orchestrator_state "WAVE_REVIEW" "Integration complete, requesting review"
# READ: agent-states/orchestrator/WAVE_REVIEW/rules.md
# Then proceed with WAVE_REVIEW work
```