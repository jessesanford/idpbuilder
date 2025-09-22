# ORCHESTRATOR STATE: SPAWN_CODE_REVIEWER_PROJECT_TEST_PLANNING

# PRIMARY DIRECTIVES

You MUST read and acknowledge these rules:

1. **R006** - Orchestrator cannot write code (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`

2. **R355** - Code Reviewer Test Planning
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R355-code-reviewer-test-planning.md`

4. **R287** - TODO Persistence Requirements (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-persistence-comprehensive.md`

5. **R288** - State File Update Requirements (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-requirements.md`

6. **R304** - Mandatory Line Counter Usage (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R304-mandatory-line-counter-usage.md`

7. **R322** - Mandatory Stop After Spawn States (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-after-spawn.md`

8. **R324** - State Transition Validation (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R324-state-transition-validation.md`


## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

## Overview
This state spawns a Code Reviewer to create project-level tests from the master architecture, implementing TDD at the project level per R341.

## Entry Criteria
- Master architecture exists (PROJECT-ARCHITECTURE.md or MASTER-ARCHITECTURE.md)
- Transitioned from WAITING_FOR_MASTER_ARCHITECTURE
- No existing PROJECT-TEST-PLAN.md

## State Responsibilities

### 1. Verify Architecture Exists
```bash
# Ensure we have architecture to work from
if [ ! -f "PROJECT-ARCHITECTURE.md" ] && [ ! -f "MASTER-ARCHITECTURE.md" ]; then
    echo "❌ Cannot create tests without master architecture!"
    exit 341
fi
```

### 2. Spawn Code Reviewer for Test Planning
```bash
/spawn-agent code-reviewer \
    --state PROJECT_TEST_PLANNING \
    --task "Create project-level tests from master architecture" \
    --tdd "Tests must be created BEFORE Phase 1 implementation per R341"
```

### 3. Record Spawn with R151 Timing
```bash
# Track spawn time for R151 compliance
SPAWN_TIME=$(date +%s)
yq -i '.spawned_agents.project_test_reviewer.id = "'$AGENT_ID'"' orchestrator-state.json
yq -i '.spawned_agents.project_test_reviewer.timestamp = "'$SPAWN_TIME'"' orchestrator-state.json
yq -i '.spawned_agents.project_test_reviewer.state = "PROJECT_TEST_PLANNING"' orchestrator-state.json
```

### 4. Update State and Stop (R313)
```bash
update_state "WAITING_FOR_PROJECT_TEST_PLAN"
save_todos "Spawned Code Reviewer for project test planning"
git add orchestrator-state.json todos/
git commit -m "state: spawned code reviewer for project tests (R341 TDD)"
git push
echo "🛑 Stopping per R313 after spawn"
```

## Exit Criteria
- Code Reviewer spawned for PROJECT_TEST_PLANNING
- State file updated to WAITING_FOR_PROJECT_TEST_PLAN
- R313: MUST stop after spawning

## Success Metrics
- ✅ Code Reviewer spawned within 5s (R151)
- ✅ State updated before stop (R324)
- ✅ R341 TDD compliance tracked
- ✅ TODOs saved (R287)

## Related Rules
- R341: TDD - tests before implementation
- R342: Early branch creation for test storage
- R151: Parallel agent timing
- R313: Mandatory stop after spawn
