# Orchestrator - SPAWN_INTEGRATION_AGENT_PROJECT State Rules

# PRIMARY DIRECTIVES

You MUST read and acknowledge these rules:

1. **R006** - Orchestrator cannot write code (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`

2. **R269** - WAVE Integration Merge Plan Protocol (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R269-WAVE-integration-merge-plan-protocol.md`

3. **R270** - PHASE Integration Merge Plan Protocol (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R270-PHASE-integration-merge-plan-protocol.md`

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

## 🛑🛑🛑 R322 MANDATORY STOP BEFORE STATE TRANSITIONS 🛑🛑🛑

**SUPREME LAW - VIOLATION = -100% IMMEDIATE FAILURE**

### YOU MUST STOP AFTER:
1. ✅ Completing all TODOs for this state
2. ✅ Updating orchestrator-state.json with new state
3. ✅ Committing and pushing the state file  
4. ✅ Providing work summary

### YOU MUST NOT:
- ❌ Continue to the next state automatically
- ❌ Start work for the new state
- ❌ Spawn agents for the new state
- ❌ Assume permission to continue

### STOP PROTOCOL:
```markdown
## 🛑 STATE TRANSITION CHECKPOINT: CURRENT_STATE → NEXT_STATE

### ✅ Current State Work Completed:
- [List completed work]

### 📊 Current Status:
- Current State: CURRENT_STATE
- Next State: NEXT_STATE
- TODOs Completed: X/Y
- State Files: Updated and committed ✅

### ⏸️ STOPPED - Awaiting User Continuation
Ready to transition to NEXT_STATE. Please use /continue-orchestrating.
```

**STOP MEANS STOP - Exit and wait for /continue-orchestrating**

---

## 🛑🛑🛑 R322 Part A MANDATORY STOP AFTER SPAWNING AGENTS 🛑🛑🛑

**SUPREME LAW - VIOLATION = -100% IMMEDIATE FAILURE**

You MUST STOP IMMEDIATELY after spawning the Integration Agent.

## State Context

**Purpose:**
Spawn Integration Agent to execute the project merge plan, merging all phase integration branches into the project integration branch.

## Primary Actions

1. **Spawn Integration Agent** with PROJECT_INTEGRATION directive
2. **Provide Resources**:
   - PROJECT-MERGE-PLAN.md location
   - Project integration workspace path
   - Phase branch details
3. **Update State File** with agent details
4. **STOP per R322 Part A**

## Valid State Transitions

- **SUCCESS** → MONITORING_PROJECT_INTEGRATION (agent spawned)
