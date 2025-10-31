# SPAWN_CODE_REVIEWER_INTEGRATE_WAVE_EFFORTS_FIX_PLAN State Rules

# PRIMARY DIRECTIVES

You MUST read and acknowledge these rules:

1. **R006** - Orchestrator cannot write code (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`

2. **R256** - Fix Planning Protocol
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R256-mandatory-phase-assessment-gate.md`

4. **R287** - TODO Persistence Requirements (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-persistence-comprehensive.md`

5. **R288** - State File Update Requirements (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`

6. **R304** - Mandatory Line Counter Usage (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R304-mandatory-line-counter-enforcement.md`

7. **R322** - Mandatory Stop After Spawn States (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-before-state-transitions.md`

8. **R324** - State Transition Validation (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R324-mandatory-line-counter-auto-detection.md`


## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

## State Purpose
Spawn Code Reviewer to create fix plans for integration issues identified during code review. This state is reached when integration code review fails.

## Entry Criteria
- **From**: WAITING_FOR_REVIEW_WAVE_INTEGRATION
- **Condition**: Integration code review failed with issues identified
- **Required**: 
  - REVIEW_WAVE_INTEGRATION_REPORT.md exists with FAIL status
  - Integration issues documented
  - Fix plans needed

## State Actions

### 1. IMMEDIATE: Spawn Code Reviewer for Fix Planning
```bash
# Spawn Code Reviewer to create integration fix plans
/spawn agent-code-reviewer CREATE_INTEGRATE_WAVE_EFFORTS_FIX_PLAN \
  --wave "W${current_wave}" \
  --branch "${wave_integration_branch}" \
  --issues-file "REVIEW_WAVE_INTEGRATION_REPORT.md" \
  --focus "integration-fixes"
```

### 2. Code Reviewer Responsibilities
The spawned Code Reviewer will:
- Analyze integration code review failures
- Create detailed fix plans for each issue
- Determine which source branches need fixes
- Identify backport requirements per R321
- Create fix distribution strategy
- Generate INTEGRATE_WAVE_EFFORTS_FIX_PLAN.md

### 3. Update State File
```json
{
  "current_state": "SPAWN_CODE_REVIEWER_INTEGRATE_WAVE_EFFORTS_FIX_PLAN",
  "phase": "integration_fixes",
  "fix_planning": {
    "reviewer": "agent-code-reviewer",
    "issues_from": "REVIEW_WAVE_INTEGRATION_REPORT.md",
    "target_branch": "${wave_integration_branch}",
    "started_at": "timestamp"
  }
}
```

## Exit Criteria

### Success Path → WAITING_FOR_INTEGRATE_WAVE_EFFORTS_FIX_PLANS
- Code Reviewer spawned successfully
- State file updated
- Transition to waiting for fix plans

### Failure Scenarios
- **Spawn Failure** → ERROR_RECOVERY
- **Missing Review Report** → ERROR_RECOVERY

## Rules Enforced
- R233: Immediate action upon state entry
- R313: Stop after spawning agent
- R321: Fixes must go to source branches
- R300: Comprehensive fix management

## Expected Output
The Code Reviewer will create:
- `INTEGRATE_WAVE_EFFORTS_FIX_PLAN.md` with:
  - Issue categorization
  - Fix assignments by effort
  - Backport requirements
  - Fix sequence and dependencies
  - Validation criteria

## Transition Rules
- **ALWAYS** → WAITING_FOR_INTEGRATE_WAVE_EFFORTS_FIX_PLANS (after spawn)


## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**


### 🚨 CRITICAL DISTINCTION: AGENT STOPS ≠ FACTORY STOPS 🚨

**TWO INDEPENDENT DECISIONS - DO NOT CONFUSE THEM:**

#### 1. Should Agent Stop Work? (R322 Technical Requirement)
- Agent completes current state
- Agent saves TODOs and commits state
- Agent exits with `exit 0` (preserves context)
- User runs /continue-orchestrating to resume
- **This is NORMAL at checkpoints**

#### 2. Should Factory Continue? (R405 Operational Status)
- Even though agent stopped, can automation proceed?
- TRUE = Healthy completion, automation can continue
- FALSE = Catastrophic failure, must halt everything
- **R322 checkpoints = TRUE (99.9% of cases)**

### THE PATTERN AT R322 CHECKPOINTS

```bash
# 1. Complete state work
echo "✅ State work complete"

# 2. Update state file
jq '.state_machine.current_state = "NEXT_STATE"' orchestrator-state-v3.json > tmp && mv tmp orchestrator-state-v3.json

# 3. Save TODOs
save_todos "R322_CHECKPOINT"

# 4. Factory continues (operational status)
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"

# 5. Agent stops (technical requirement)
exit 0
```

**Both happen together! Agent stops AND factory continues!**

### WHEN TO USE EACH FLAG VALUE

**TRUE (99.9%):**
- ✅ R322 checkpoint reached
- ✅ State work completed successfully
- ✅ Ready for /continue-orchestrating
- ✅ Waiting for user to continue (NORMAL)
- ✅ Plan ready for review (agent done, factory proceeds)

**FALSE (0.1%):**
- ❌ CATASTROPHIC unrecoverable error
- ❌ Data corruption spreading
- ❌ Critical security violation
- ❌ NOT for R322 checkpoints
- ❌ NOT for user review needs
### 🚨 SPAWN STATE PATTERN - R322 + R405 USAGE 🚨

**Spawning operations require R322 stop for context preservation:**
```bash
# After spawning agent(s)
echo "✅ Spawned agents for work"

# R322 checkpoint (context preservation)
echo "🛑 R322: Stopping after spawn for context preservation"

# Flag? → MUST BE TRUE (normal operation!)
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"

# Stop inference
exit 0
```

**Why TRUE is correct:**
- Spawning is NORMAL operation
- System knows next state
- Automation can continue
- **Context preservation ≠ manual intervention needed!**
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

