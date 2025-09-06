# SPAWN_INTEGRATION_AGENT_PHASE State Rules

## State Purpose
Spawn Integration Agent (SW Engineer in integration mode) to execute the phase merge plan, integrating all wave branches into the phase integration branch.

## Critical Rules

### 🔴🔴🔴 RULE R322: MANDATORY STOP BEFORE STATE TRANSITION (SUPREME LAW)
- **STOP** and save state before ANY transition
- **READ** orchestrator-state.yaml to verify current state
- **VALIDATE** next state exists in SOFTWARE-FACTORY-STATE-MACHINE.md
- **VIOLATION = IMMEDIATE FAILURE**

### 🔴🔴🔴 RULE R313: MANDATORY STOP AFTER SPAWNING AGENTS (SUPREME LAW)
- **MUST STOP IMMEDIATELY** after spawning Integration Agent
- **RECORD** what was spawned in state file
- **SAVE** TODOs and commit state changes
- **EXIT** with clear continuation instructions
- **VIOLATION = CONTEXT LOSS AND RULE FORGETTING**

### 🚨🚨🚨 RULE R290: STATE RULE VERIFICATION (BLOCKING)
- **MUST** verify this rules file exists and is loaded
- **MUST** acknowledge all rules before proceeding
- **MUST** validate state transitions against state machine

### 🚨🚨🚨 RULE R208: SPAWN DIRECTORY PROTOCOL (BLOCKING)
- **MUST** spawn Integration Agent in phase integration directory
- **MUST** provide PHASE-MERGE-PLAN.md location
- **MUST** ensure agent has phase integration branch checked out
- **MUST** verify all wave branches are accessible

### 🚨🚨🚨 RULE R285: PHASE INTEGRATION REQUIREMENTS (BLOCKING)
- Integration MUST follow PHASE-MERGE-PLAN.md exactly
- Waves MUST be merged in sequential order
- Each merge MUST be tested before next merge
- Failed merges trigger IMMEDIATE_BACKPORT_REQUIRED per R321

### ⚠️⚠️⚠️ RULE R321: IMMEDIATE BACKPORT PROTOCOL (WARNING)
- ANY merge conflict or failure requires immediate fix in source branch
- Integration branches are READ-ONLY for code changes
- Fixes MUST go to effort branches first, then re-merge

### ⚠️⚠️⚠️ RULE R287: TODO PERSISTENCE (WARNING)
- **MUST** save TODOs before spawning
- **MUST** save TODOs after spawn completes
- **MUST** commit and push TODO state

## Required Actions

1. **Verify Phase Merge Plan Exists**
   ```bash
   # Check for merge plan
   ls -la phase-*-integration/PHASE-MERGE-PLAN.md
   
   # Verify plan contents
   cat phase-*-integration/PHASE-MERGE-PLAN.md
   ```

2. **Verify Phase Integration Infrastructure**
   ```bash
   # Check phase integration directory
   ls -la phase-*-integration/
   
   # Verify phase integration branch
   cd phase-*-integration/
   git branch --show-current  # Should be phase-X-integration
   git pull origin phase-X-integration
   ```

3. **Verify Wave Branches Ready**
   ```bash
   # List all wave branches for phase
   git ls-remote origin | grep "wave-${PHASE_NUM}-"
   
   # Verify all waves marked complete
   grep "waves_completed" orchestrator-state.yaml
   ```

4. **Spawn Integration Agent**
   ```bash
   # Spawn SW Engineer as Integration Agent
   cd phase-X-integration/
   
   /spawn @agent-sw-engineer PHASE_INTEGRATION_EXECUTION \
     --merge-plan "PHASE-MERGE-PLAN.md" \
     --target-branch "phase-X-integration" \
     --wave-branches "wave-X-1,wave-X-2,..." \
     --output "PHASE-INTEGRATION-REPORT.md"
   ```

5. **Update State File**
   ```yaml
   current_state: SPAWN_INTEGRATION_AGENT_PHASE
   spawned_agents:
     - agent: sw-engineer
       role: integration-agent
       directory: phase-X-integration
       task: phase_integration_execution
       merge_plan: PHASE-MERGE-PLAN.md
       timestamp: YYYY-MM-DD HH:MM:SS
   phase_integration:
     status: IN_PROGRESS
     agent_spawned: true
   ```

6. **Save and Exit (R313 MANDATORY)**
   ```bash
   # Save TODOs
   save_todos "SPAWN_INTEGRATION_AGENT_PHASE"
   
   # Commit state
   git add orchestrator-state.yaml todos/*.todo
   git commit -m "state: spawned Integration Agent for phase integration"
   git push
   
   # EXIT IMMEDIATELY
   echo "Integration Agent spawned for phase merges. Use /continue orchestrator to resume."
   exit 0
   ```

## Transition Rules

### Valid Next States
- **MONITORING_PHASE_INTEGRATION** - After spawning (MANDATORY per R313)

### Invalid Transitions
- ❌ Any state other than MONITORING_PHASE_INTEGRATION
- ❌ Continuing work after spawn (violates R313)
- ❌ Attempting integration yourself (orchestrator doesn't merge)

## Common Violations to Avoid

1. **Not stopping after spawn** - Violates R313, causes context loss
2. **Missing PHASE-MERGE-PLAN.md** - Integration has no guidance
3. **Wrong directory for spawn** - Integration happens in wrong place
4. **Not verifying wave completeness** - Merging incomplete work
5. **Forgetting R321 immediate backport** - Creating fix commits in integration branch

## Verification Commands

```bash
# Verify state entry
echo "Entered SPAWN_INTEGRATION_AGENT_PHASE at $(date)"

# Verify merge plan exists
test -f phase-*/PHASE-MERGE-PLAN.md && echo "✓ Merge plan found" || echo "✗ Missing merge plan"

# Verify phase infrastructure
ls -la phase-*-integration/
git branch -r | grep "phase-.*-integration"

# Verify all waves complete
grep -c "status: COMPLETE" orchestrator-state.yaml

# After spawn, verify stop
echo "STOPPING per R313 - Integration Agent spawned for phase merges"
```

## Integration Failure Handling

If integration fails (detected in MONITORING_PHASE_INTEGRATION):
1. Transition to **IMMEDIATE_BACKPORT_REQUIRED** (R321)
2. DO NOT attempt fixes in integration branch
3. Spawn engineers to fix source branches
4. Re-run entire phase integration after fixes

## References
- R313: rule-library/R313-mandatory-stop-after-spawn.md
- R208: rule-library/R208-spawn-directory-protocol.md
- R285: rule-library/R285-phase-integration-requirements.md
- R321: rule-library/R321-immediate-backport-during-integration.md
- R287: rule-library/R287-todo-persistence-comprehensive.md
- R290: rule-library/R290-state-rule-verification.md
- R322: rule-library/R322-mandatory-stop-before-transition.md