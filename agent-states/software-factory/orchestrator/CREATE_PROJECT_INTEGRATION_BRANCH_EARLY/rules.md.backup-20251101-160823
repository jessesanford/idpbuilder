# ORCHESTRATOR STATE: CREATE_PROJECT_INTEGRATION_BRANCH_EARLY


## 🚨 State Manager Bookend Pattern (MANDATORY)

**BEFORE this state**:
- State Manager validated transition via STARTUP_CONSULTATION
- You are here because State Manager directed you here
- orchestrator-state-v3.json shows validated_by: "state-manager"

**DURING this state**:
- Perform state-specific work
- NEVER call update_state directly
- Prepare results for State Manager
- Propose next state (don't decide!)

**AFTER this state**:
- Spawn State Manager SHUTDOWN_CONSULTATION
- Provide results and proposed next state
- State Manager validates and decides actual next state
- Transition to State Manager's required_next_state

**CRITICAL**: The orchestrator PROPOSES, the State Manager DECIDES!

---

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

## 📋 PRIMARY DIRECTIVES FOR CREATE_PROJECT_INTEGRATION_BRANCH_EARLY STATE

### Core Mandatory Rules (ALL orchestrator states must have these):

1. **🚨🚨🚨 R006** - ORCHESTRATOR NEVER WRITES CODE OR PERFORMS FILE OPERATIONS (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`
   - Criticality: BLOCKING - Automatic termination, 0% grade
   - Summary: NEVER write, copy, move, or manipulate ANY code files - delegate ALL to agents

2. **🔴🔴🔴 R287** - TODO PERSISTENCE COMPREHENSIVE (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-persistence-comprehensive.md`
   - Criticality: SUPREME - -20% to -100% penalty for violations
   - Summary: MUST save TODOs within 30s after write, every 10 messages, before transitions

3. **🔴🔴🔴 R288** - STATE FILE UPDATE REQUIREMENTS (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
   - Criticality: SUPREME - State updates required for all transitions
   - Summary: MUST update orchestrator-state-v3.json before EVERY state transition

4. **🔴🔴🔴 R322 Part A** - Mandatory Stop After Spawn States
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-before-state-transitions.md`
   - Criticality: SUPREME LAW - Must stop after spawning
   - Summary: ALL spawn states require STOP after spawning agents

### State-Specific Rules:

5. **🔴🔴🔴 R308** - Incremental Branching Strategy
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R308-incremental-branching-strategy.md`
   - Criticality: SUPREME LAW - Phase 2 NEVER from main!
   - Summary: Integration branches must follow incremental strategy

6. **🔴🔴🔴 R009** - Mandatory Wave/Phase Integration Protocol
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R009-integration-branch-creation.md`
   - Criticality: SUPREME LAW - Wave N+1 REQUIRES Wave N integration
   - Summary: Enforces trunk-based development

## Overview
This state creates the project-integration branch immediately after test planning, storing tests in their permanent location per R342.

## Entry Criteria
- Transitioned from WAITING_FOR_PROJECT_TEST_PLAN
- PROJECT-TEST-PLAN.md exists
- Project tests have been created

## State Responsibilities

### 1. Create Project Integration Workspace
```bash
# R104: Proper workspace isolation
PROJECT_INT_DIR="/efforts/project/integration-workspace"
mkdir -p "$PROJECT_INT_DIR"
cd "$PROJECT_INT_DIR"
echo "📁 Created project integration workspace"
```

### 2. Clone Target Repository
```bash
# Get target repository from config
TARGET_REPO_URL=$(yq '.pre_planned_infrastructure.target_repo_url' "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json")
if [ -z "$TARGET_REPO_URL" ] || [ "$TARGET_REPO_URL" = "null" ]; then
    echo "❌ No target repository configured!"
    exit 1
fi

git clone "$TARGET_REPO_URL" target-repo
cd target-repo
echo "✅ Cloned target repository"
```

### 3. Create Project Integration Branch
```bash
# Create branch from main
git checkout -b project-integration
echo "🌿 Created project-integration branch"
```

### 4. Store Project Tests (R342 Core Requirement)
```bash
# Copy tests to integration branch
mkdir -p tests/project
if [ -d "$CLAUDE_PROJECT_DIR/project-tests" ]; then
    cp -r "$CLAUDE_PROJECT_DIR/project-tests/"* tests/project/
    echo "📋 Copied project tests to branch"
fi

# Copy test plan and harness
cp "$CLAUDE_PROJECT_DIR/PROJECT-TEST-PLAN.md" .
cp "$CLAUDE_PROJECT_DIR/PROJECT-TEST-HARNESS.sh" tests/project/ 2>/dev/null || true

# Verify tests exist
ls -la tests/project/
```

### 5. Commit and Push Tests
```bash
# Add and commit tests
git add tests/ PROJECT-TEST-PLAN.md
git commit -m "test: add project-level tests (R341 TDD, R342 early branch)"
git push -u origin project-integration

echo "✅ R342 COMPLIANT: Tests stored in project-integration branch"
```

### 6. Update State File with Branch Info
```bash
# Track integration branch creation
yq -i '.project_integration.branch = "project-integration"' orchestrator-state-v3.json
yq -i '.project_integration.created_at = "'$(date -Iseconds)'"' orchestrator-state-v3.json
yq -i '.project_integration.has_tests = true' orchestrator-state-v3.json
yq -i '.project_integration.workspace = "'$PROJECT_INT_DIR'"' orchestrator-state-v3.json
```

### 7. Prepare for Phase 1 Start
```bash
PROPOSED_NEXT_STATE="SPAWN_ARCHITECT_PHASE_PLANNING"
TRANSITION_REASON="Project integration branch created with tests per R342, ready for Phase 1"
echo "📍 Will propose transition to: $PROPOSED_NEXT_STATE"
# NOTE: Actual transition via State Manager in completion checklist
```

## Exit Criteria
- project-integration branch created
- Tests committed to branch
- Branch pushed to remote
- State transitions to INIT for Phase 1 start

## Success Metrics
- ✅ Integration branch exists
- ✅ Tests stored in tests/project/
- ✅ Branch pushed to remote
- ✅ R342 compliance achieved
- ✅ Ready for Phase 1



## 🔴🔴🔴 STATE COMPLETION CHECKLIST 🔴🔴🔴

**Execute these steps IN ORDER to properly complete CREATE_PROJECT_INTEGRATION_BRANCH_EARLY:**

### ✅ Step 1: Complete State-Specific Work
**Refer to "Primary Actions" and "Mandatory Validations" sections above for state-specific tasks.**

Once all state work is complete, proceed to mandatory exit protocol:

---

### ✅ Step 2: Set Proposed Next State and Transition Reason
```bash
# Based on state work results, set variables for State Manager
PROPOSED_NEXT_STATE="[DETERMINE_FROM_STATE_LOGIC]"
TRANSITION_REASON="[REASON_FOR_TRANSITION]"
echo "Proposed next state: $PROPOSED_NEXT_STATE"
echo "Transition reason: $TRANSITION_REASON"
```

---

### ✅ Step 3: Spawn State Manager for State Transition (R288 - SUPREME LAW)
```bash
# State Manager handles ALL state file updates atomically
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "CURRENT_STATE_NAME" \
  --proposed-next-state "$PROPOSED_NEXT_STATE" \
  --transition-reason "$TRANSITION_REASON"

# State Manager will:
# 1. Validate the transition against state machine
# 2. Update all 4 state tracking locations atomically:
#    - orchestrator-state-v3.json
#    - orchestrator-state-demo.json
#    - .cascade-state-backup.json
#    - .orchestrator-state-v3.json
# 3. Commit and push all changes
# 4. Return control

echo "✅ State Manager completed transition"
```

---

### ✅ Step 4: Save TODOs (R287 - SUPREME LAW)
```bash
# Save TODO state before exit (R287 trigger)
save_todos "STATE_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - state complete [R287]"; then
    echo "❌ ERROR: Failed to commit TODO files"
    echo "This is non-fatal but TODOs may be lost in compaction"
    echo "Proceeding with state execution..."
    # Don't exit - TODO commit failure is not fatal
fi

git push || echo "⚠️ WARNING: TODO push failed - committed locally"
echo "✅ TODOs saved and committed"
git push
echo "✅ TODOs saved and committed"
```

---

### ✅ Step 5: Output Continuation Flag (R405 - SUPREME LAW) ⚠️ MANDATORY
```bash
# Output continuation flag as LAST action (R405)
# Use TRUE for normal completion, FALSE only for catastrophic errors

echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"
```

**⚠️ THIS MUST BE THE ABSOLUTE LAST LINE OF OUTPUT BEFORE EXIT! ⚠️**

---

### ✅ Step 6: Stop Processing (R322 - SUPREME LAW)
```bash
# Stop for context preservation (R322)
echo "🛑 Stopping for context preservation - use /continue-orchestrating to resume"
exit 0
```

---

## 🚨 CHECKLIST ENFORCEMENT 🚨

**Skipping ANY step in this checklist = FAILURE:**
- Missing Step 2: No proposed next state = State Manager can't proceed
- Missing Step 3: No State Manager consultation = bypassing bookend pattern (R288 violation, -100%)
- Missing Step 4: No TODO save = work lost (R287 violation, -20% to -100%)
- Missing Step 5: No continuation flag = automation stops (R405 violation, -100%)
- Missing Step 6: No exit = R322 violation (-100%)

**ALL 6 STEPS ARE MANDATORY - NO EXCEPTIONS**
**NOTE**: State file validation and commits now handled by State Manager (SF 3.0 pattern)

## ✅ EXIT CHECKLIST (R405 - Continuation Flag Protocol) 🔴🔴🔴

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

### THE PATTERN AT R322 CHECKPOINTS (SF 3.0)

```bash
# 1. Complete state work
echo "✅ State work complete"

# 2. Set proposed next state
PROPOSED_NEXT_STATE="SPAWN_ARCHITECT_PHASE_PLANNING"
TRANSITION_REASON="State work complete"

# 3. Spawn State Manager for state transition
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "CURRENT_STATE" \
  --proposed-next-state "$PROPOSED_NEXT_STATE" \
  --transition-reason "$TRANSITION_REASON"
# State Manager updates all 4 state files atomically

# 4. Save TODOs
save_todos "R322_CHECKPOINT"

# 5. Factory continues (operational status)
echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"

# 6. Agent stops (technical requirement)
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
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**


## Related Rules
- R342: Early integration branch creation (SUPREME LAW)
- R341: TDD - tests before implementation
- R104: Target repository isolation
- R308: Incremental branching strategy
