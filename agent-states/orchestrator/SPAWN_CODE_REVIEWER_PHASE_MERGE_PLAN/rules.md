# SPAWN_CODE_REVIEWER_PHASE_MERGE_PLAN State Rules

## State Purpose
Spawn Code Reviewer agent to create a comprehensive merge plan for integrating all wave branches into the phase integration branch.

## Critical Rules

### 🔴🔴🔴 RULE R322: MANDATORY STOP BEFORE STATE TRANSITION (SUPREME LAW)
- **STOP** and save state before ANY transition
- **READ** orchestrator-state.json to verify current state
- **VALIDATE** next state exists in SOFTWARE-FACTORY-STATE-MACHINE.md
- **VIOLATION = IMMEDIATE FAILURE**

### 🔴🔴🔴 RULE R322 Part A: MANDATORY STOP AFTER SPAWNING AGENTS (SUPREME LAW)
- **MUST STOP IMMEDIATELY** after spawning Code Reviewer
- **RECORD** what was spawned in state file
- **SAVE** TODOs and commit state changes
- **EXIT** with clear continuation instructions
- **VIOLATION = CONTEXT LOSS AND RULE FORGETTING**

### 🚨🚨🚨 RULE R290: STATE RULE VERIFICATION (BLOCKING)
- **MUST** verify this rules file exists and is loaded
- **MUST** acknowledge all rules before proceeding
- **MUST** validate state transitions against state machine

### 🚨🚨🚨 RULE R208: SPAWN DIRECTORY PROTOCOL (BLOCKING)
- **MUST** spawn Code Reviewer in phase integration directory
- **MUST** provide phase integration branch information
- **MUST** pass list of wave branches to merge
- **MUST** ensure Code Reviewer has access to PHASE-INTEGRATION-PLAN.md

### ⚠️⚠️⚠️ RULE R269: MERGE PLAN REQUIREMENTS (WARNING)
- Code Reviewer MUST create PHASE-MERGE-PLAN.md
- Plan MUST list all wave branches in order
- Plan MUST specify merge strategy (sequential recommended)
- Plan MUST identify potential conflicts

### ⚠️⚠️⚠️ RULE R270: MERGE SEQUENCE VALIDATION (WARNING)
- Wave branches MUST be merged in sequential order
- Each merge MUST be tested before proceeding
- Failed merges MUST trigger immediate backport per R321

### ⚠️⚠️⚠️ RULE R287: TODO PERSISTENCE (WARNING)
- **MUST** save TODOs before spawning
- **MUST** save TODOs after spawn completes
- **MUST** commit and push TODO state

## Required Actions

1. **Setup Phase Integration Infrastructure**
   ```bash
   # Verify phase integration directory exists
   ls -la phase-*-integration/
   
   # Verify phase integration branch exists
   git ls-remote origin | grep "phase-.*-integration"
   ```

2. **Gather Wave Information**
   ```bash
   # List all wave branches for this phase
   grep "wave_.*_branches" orchestrator-state.json
   
   # Verify all waves are complete
   grep "waves_completed" orchestrator-state.json
   ```

3. **Spawn Code Reviewer**
   ```bash
   # Spawn in phase integration directory
   cd phase-X-integration/
   
   # Create spawn command with merge planning parameters
   /spawn @agent-code-reviewer PHASE_MERGE_PLANNING \
     --phase-branch "phase-X-integration" \
     --wave-branches "wave-X-1,wave-X-2,..." \
     --output "PHASE-MERGE-PLAN.md"
   ```

4. **Update State File**
   ```yaml
   current_state: SPAWN_CODE_REVIEWER_PHASE_MERGE_PLAN
   spawned_agents:
     - agent: code-reviewer
       directory: phase-X-integration
       task: phase_merge_planning
       timestamp: YYYY-MM-DD HH:MM:SS
   ```

5. **Save and Exit (R322 Part A MANDATORY)**
   ```bash
   # Save TODOs
   save_todos "SPAWN_CODE_REVIEWER_PHASE_MERGE_PLAN"
   
   # Commit state
   git add orchestrator-state.json todos/*.todo
   git commit -m "state: spawned Code Reviewer for phase merge planning"
   git push
   
   # EXIT IMMEDIATELY
   echo "Code Reviewer spawned for phase merge planning. Use /continue orchestrator to resume."
   exit 0
   ```

## Transition Rules

### Valid Next States
- **WAITING_FOR_PHASE_MERGE_PLAN** - After spawning (MANDATORY per R322 Part A)

### Invalid Transitions
- ❌ Any state other than WAITING_FOR_PHASE_MERGE_PLAN
- ❌ Continuing work after spawn (violates R322 Part A)
- ❌ Spawning multiple agents (violates R322 Part A)

## Common Violations to Avoid

1. **Not stopping after spawn** - Violates R322 Part A, causes context loss
2. **Forgetting to update state file** - Causes state machine corruption
3. **Not verifying phase infrastructure** - Causes spawn failures
4. **Missing wave branch information** - Results in incomplete merge plans
5. **Not saving TODOs** - Violates R287, loses progress

## Verification Commands

```bash
# Verify state entry
echo "Entered SPAWN_CODE_REVIEWER_PHASE_MERGE_PLAN at $(date)"

# Verify phase readiness
ls -la phase-*-integration/
git branch -r | grep "phase-.*-integration"

# Verify waves completed
grep "waves_completed" orchestrator-state.json

# After spawn, verify stop
echo "STOPPING per R322 Part A - Code Reviewer spawned for phase merge planning"
```

## References
- R322 Part A: rule-library/R322 Part A-mandatory-stop-after-spawn.md
- R208: rule-library/R208-spawn-directory-protocol.md
- R269: rule-library/R269-merge-plan-requirements.md
- R270: rule-library/R270-merge-sequence-validation.md
- R287: rule-library/R287-todo-persistence-comprehensive.md
- R290: rule-library/R290-state-rule-verification.md
- R322: rule-library/R322-mandatory-stop-before-transition.md

### 🔴🔴🔴 MANDATORY VALIDATION REQUIREMENT 🔴🔴🔴

**Per R288 and R324**: ALL state file updates MUST be validated before commit:

```bash
# After ANY update to orchestrator-state.json:
"$CLAUDE_PROJECT_DIR/tools/validate-state.sh" orchestrator-state.json || {
    echo "❌ State file validation failed!"
    exit 288
}
```

**Use helper functions for automatic validation:**
```bash
# Source the helper functions
source "$CLAUDE_PROJECT_DIR/utilities/state-file-update-functions.sh"

# Use safe functions that include validation:
safe_state_transition "NEW_STATE" "reason"
safe_update_field "field_name" "value"
```
