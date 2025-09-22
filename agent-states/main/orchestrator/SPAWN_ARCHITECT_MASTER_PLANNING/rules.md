# ORCHESTRATOR STATE: SPAWN_ARCHITECT_MASTER_PLANNING

# PRIMARY DIRECTIVES

You MUST read and acknowledge these rules:

1. **R006** - Orchestrator cannot write code (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`

2. **R308** - Architect Integration Role
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R308-architect-integration-role.md`

5. **R287** - TODO Persistence Requirements (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-persistence-comprehensive.md`

6. **R288** - State File Update Requirements (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-requirements.md`

7. **R304** - Mandatory Line Counter Usage (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R304-mandatory-line-counter-usage.md`

8. **R322** - Mandatory Stop After Spawn States (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-after-spawn.md`

9. **R324** - State Transition Validation (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R324-state-transition-validation.md`


## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

## Overview
This state spawns the Architect to create the master architecture that defines the entire project structure.

## Entry Criteria
- No existing PROJECT-ARCHITECTURE.md or MASTER-ARCHITECTURE.md
- Starting a brand new project
- Transitioned from INIT state

## State Responsibilities

### 1. Verify No Existing Master Plan
```bash
# Check for existing architecture
if [ -f "PROJECT-ARCHITECTURE.md" ] || [ -f "MASTER-ARCHITECTURE.md" ]; then
    echo "⚠️ Master architecture already exists"
    transition_to "SPAWN_CODE_REVIEWER_PROJECT_TEST_PLANNING"
else
    echo "✅ No master architecture found, proceeding with spawn"
fi
```

### 2. Spawn Architect Agent
```bash
/spawn-agent architect \
    --state MASTER_PLANNING \
    --task "Create master architecture for entire project"
```

### 3. Record Spawn in State File
- Update orchestrator-state.json with spawned agent
- Track timestamp per R151
- Record spawn ID

### 4. Transition to Waiting State
```bash
update_state "WAITING_FOR_MASTER_ARCHITECTURE"
save_todos "Spawned architect for master planning"
git add orchestrator-state.json todos/
git commit -m "state: spawned architect for master planning"
git push
```

## Exit Criteria
- Architect spawned successfully
- State file updated to WAITING_FOR_MASTER_ARCHITECTURE
- R313: MUST stop after spawning

## Success Metrics
- ✅ Architect spawned within 5s (R151)
- ✅ State file updated before stop (R324)
- ✅ TODOs saved (R287)

## Related Rules
- R151: Parallel agent timing
- R287: TODO persistence
- R313: Mandatory stop after spawn
- R324: State update before stop
- R341: TDD requirement
