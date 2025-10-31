# Orchestrator - INTEGRATE_WAVE_EFFORTS_TESTING State Rules


## 🛑🛑🛑 R322 MANDATORY STOP BEFORE STATE TRANSITIONS 🛑🛑🛑

**SUPREME LAW - VIOLATION = -100% IMMEDIATE FAILURE**

### YOU MUST STOP AFTER:
1. ✅ Completing all TODOs for this state
2. ✅ Updating orchestrator-state-v3.json with new state
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


## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED INTEGRATE_WAVE_EFFORTS_TESTING STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_INTEGRATE_WAVE_EFFORTS_TESTING
echo "$(date +%s) - Rules read and acknowledged for INTEGRATE_WAVE_EFFORTS_TESTING" > .state_rules_read_orchestrator_INTEGRATE_WAVE_EFFORTS_TESTING
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY INTEGRATE_WAVE_EFFORTS TESTING WORK UNTIL RULES ARE READ:
- ❌ Start merging efforts
- ❌ Start running tests
- ❌ Update state files
- ❌ Continue to next state
- ❌ Think about what to do in this state

### ✅ YOU MUST IMMEDIATELY:

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

## 📋 PRIMARY DIRECTIVES

### 🚨🚨🚨 RULE R006 - Orchestrator NEVER Writes Code [BLOCKING]
**Source:** `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`
**Criticality:** BLOCKING - Any code operation = -100% IMMEDIATE FAILURE

READ THIS RULE FILE: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`

### 🚨🚨🚨 RULE R329 - Orchestrator NEVER Performs Git Merges [BLOCKING]
**Source:** `$CLAUDE_PROJECT_DIR/rule-library/R329-orchestrator-never-performs-merges.md`
**Criticality:** BLOCKING - Any merge operation = -100% IMMEDIATE FAILURE

READ THIS RULE FILE: `$CLAUDE_PROJECT_DIR/rule-library/R329-orchestrator-never-performs-merges.md`

**⚠️ R006 + R329 WARNING FOR INTEGRATE_WAVE_EFFORTS_TESTING STATE:**
- DO NOT execute git merge commands yourself! (R329)
- DO NOT resolve merge conflicts yourself! (R006 + R329)
- DO NOT edit code to fix integration issues! (R006)
- DO NOT apply patches or fixes directly! (R006 + R329)
- MUST spawn Integration Agent for ALL merges (R329)
- Document all issues for appropriate agents to resolve
- You only coordinate - NEVER execute merges or modify code

### 🚨🚨🚨 RULE R271 - Mandatory Production Ready Validation [BLOCKING]
**MUST validate production readiness** | Source: rule-library/R271-mandatory-production-ready-validation.md

READ THIS RULE FILE: `$CLAUDE_PROJECT_DIR/rule-library/R271-mandatory-production-ready-validation.md`

This rule defines the mandatory production ready validation process required before any code can be considered complete.

### 🚨🚨🚨 RULE R273 - Runtime Specific Validation [BLOCKING]
**MUST validate runtime specific requirements** | Source: rule-library/R273-runtime-specific-validation.md

READ THIS RULE FILE: `$CLAUDE_PROJECT_DIR/rule-library/R273-runtime-specific-validation.md`

This rule defines runtime-specific validation requirements based on the technology stack being used.

### 🚨🚨🚨 RULE R280 - Main Branch Protection [BLOCKING]
**SOFTWARE FACTORY NEVER MERGES TO MAIN** | Source: rule-library/R280-main-branch-protection.md

READ THIS RULE FILE: `$CLAUDE_PROJECT_DIR/rule-library/R280-main-branch-protection.md`

Software Factory creates MASTER-PR-PLAN.md for humans to execute PRs. We NEVER push to main ourselves.

### 🚨🚨🚨 RULE R328 - Integration Freshness Validation [BLOCKING]
**MUST verify integration branch freshness before merging** | Source: rule-library/R328-integration-freshness-validation.md

READ THIS RULE FILE: `$CLAUDE_PROJECT_DIR/rule-library/R328-integration-freshness-validation.md`

Integration branches become stale when fixes are applied after creation. ALWAYS check timestamps!

## 🎯 STATE OBJECTIVES

In the INTEGRATE_WAVE_EFFORTS_TESTING state, you are responsible for:

1. **Coordinating Integration via Integration Agent (R329 MANDATORY)**
   - Identify all effort directories and branches
   - Spawn Code Reviewer to create merge plan
   - **SPAWN INTEGRATE_WAVE_EFFORTS AGENT TO EXECUTE ALL MERGES**
   - Monitor Integration Agent progress via reports
   - NEVER execute merges yourself (R329 VIOLATION)

2. **Verifying Integration Success (via Agent Reports)**
   - Review Integration Agent's INTEGRATE_WAVE_EFFORTS-REPORT.md
   - **🔴 R291 MANDATORY: Validate demos before approving integration** ← NEW
   - Check for conflicts documented by Integration Agent
   - Spawn Code Reviewer for build/test validation if needed
   - Document any issues for next states

## 🔴🔴🔴 R291 DEMO VALIDATION (MANDATORY BEFORE APPROVAL) 🔴🔴🔴

**BEFORE approving ANY integration (wave or phase), you MUST validate demo requirements:**

```bash
# MANDATORY R291 validation after reading integration report
validate_integration_demos() {
    local INTEGRATE_WAVE_EFFORTS_REPORT="$1"  # Path to Integration Agent's report

    echo "🔴🔴🔴 R291 MANDATORY DEMO VALIDATION 🔴🔴🔴"

    if [ ! -f "$INTEGRATE_WAVE_EFFORTS_REPORT" ]; then
        echo "❌ Cannot find integration report: $INTEGRATE_WAVE_EFFORTS_REPORT"
        echo "CANNOT validate R291 without integration report"
        return 1
    fi

    # Check if integration report mentions R291 gate check
    if grep -q "R291.*FAILED\|Demo Gate.*FAILED\|DEMO.*FAILED" "$INTEGRATE_WAVE_EFFORTS_REPORT"; then
        echo ""
        echo "🔴🔴🔴 ORCHESTRATOR DETECTED: R291 DEMO GATE FAILURE 🔴🔴🔴"
        echo "Integration report shows missing/broken demos"
        echo ""
        echo "Integration Agent reported demo failures - CANNOT approve integration"
        echo ""
        echo "MANDATORY ACTIONS:"
        echo "1. ❌ DO NOT approve this integration"
        echo "2. ❌ DO NOT transition to next wave/phase"
        echo "3. ✅ MUST transition to ERROR_RECOVERY state"
        echo "4. ✅ MUST create fix plan for demo creation"
        echo ""

        # Update state to ERROR_RECOVERY
        yq -i '.state_machine.current_state = "ERROR_RECOVERY"' "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json"
        yq -i '.error_recovery.triggered_by = "R291_DEMO_GATE_FAILURE"' "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json"
        yq -i '.error_recovery.error_type = "MISSING_DEMO_REQUIREMENTS"' "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json"
        yq -i '.error_recovery.timestamp = "'"$(date -Iseconds)"'"' "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json"
        yq -i '.error_recovery.requires = "Demo artifacts creation"' "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json"

        cd "$CLAUDE_PROJECT_DIR"
        git add orchestrator-state-v3.json
        git commit -m "error: R291 demo gate failure detected - entering ERROR_RECOVERY" || true
        git push || true

        echo "🔴 STATE UPDATED: Transitioned to ERROR_RECOVERY per R291"
        echo "❌ Integration BLOCKED until demos created"
        exit 291
    fi

    # Check for positive demo confirmation
    if grep -q "R291.*PASSED\|Demo Gate.*PASSED\|✅.*demo\|✅.*Demo" "$INTEGRATE_WAVE_EFFORTS_REPORT"; then
        echo "✅ R291 validation PASSED - Integration report confirms demos present"
        return 0
    fi

    # If no demo gate check mentioned at all
    echo "⚠️  WARNING: Integration report does not mention R291 demo gate check"
    echo "⚠️  This may indicate the Integration Agent skipped R291 validation"
    echo ""
    echo "VERIFYING demo artifacts directly..."

    # Determine integration workspace
    INTEGRATE_WAVE_EFFORTS_DIR=$(dirname "$INTEGRATE_WAVE_EFFORTS_REPORT")

    # Check for demo files
    DEMO_FILES_FOUND=false
    if [ -f "$INTEGRATE_WAVE_EFFORTS_DIR/demo-wave-integration.sh" ] || \
       [ -f "$INTEGRATE_WAVE_EFFORTS_DIR/demo-integration.sh" ] || \
       [ -f "$INTEGRATE_WAVE_EFFORTS_DIR/WAVE-DEMO.md" ] || \
       [ -f "$INTEGRATE_WAVE_EFFORTS_DIR/DEMO.md" ] || \
       [ -f "$INTEGRATE_WAVE_EFFORTS_DIR/INTEGRATE_WAVE_EFFORTS-DEMO.md" ]; then
        DEMO_FILES_FOUND=true
    fi

    if [ "$DEMO_FILES_FOUND" = false ]; then
        echo "❌ R291 VIOLATION: No demo artifacts found in integration workspace!"
        echo "   Checked: $INTEGRATE_WAVE_EFFORTS_DIR"
        echo ""
        echo "🔴 MANDATORY ERROR_RECOVERY TRANSITION"
        echo "Integration CANNOT be approved without demos"

        # Transition to ERROR_RECOVERY
        yq -i '.state_machine.current_state = "ERROR_RECOVERY"' "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json"
        yq -i '.error_recovery.triggered_by = "R291_MISSING_DEMOS"' "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json"
        yq -i '.error_recovery.timestamp = "'"$(date -Iseconds)"'"' "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json"

        cd "$CLAUDE_PROJECT_DIR"
        git add orchestrator-state-v3.json
        git commit -m "error: R291 missing demos detected - entering ERROR_RECOVERY" || true
        git push || true

        exit 291
    fi

    echo "✅ R291 validation: Demo artifacts found in integration workspace"
    return 0
}

# MANDATORY: Call this function after reading integration report
# validate_integration_demos "$INTEGRATE_WAVE_EFFORTS_WORKSPACE/INTEGRATE_WAVE_EFFORTS-REPORT.md"
```

### EXIT CRITERIA FOR INTEGRATE_WAVE_EFFORTS APPROVAL:

**Integration is ONLY approved when:**
- ✅ Integration Agent successfully merged all branches
- ✅ Build passed
- ✅ Tests passed
- ✅ **R291: Demos present and validated** ← NEW MANDATORY CHECK
- ✅ Integration report complete

**If R291 fails:**
- ❌ DO NOT approve integration
- ❌ DO NOT transition to next wave/phase
- ✅ MUST transition to ERROR_RECOVERY
- ✅ MUST create fix plan for demos

**Penalty for approving integration without demos: -100% IMMEDIATE FAILURE**

---

3. **Creating Orchestration Report**
   - Document which agents were spawned
   - Summarize Integration Agent's findings
   - Note any issues requiring attention
   - Track state transitions

## 📝 REQUIRED ACTIONS

### 🔴🔴🔴 CRITICAL: UNDERSTANDING REPOSITORY CONTEXTS 🔴🔴🔴

**BEFORE YOU START, YOU MUST UNDERSTAND WHERE THINGS ARE:**

1. **Software Factory Repository** (`/home/vscode/software-factory-template/`)
   - Contains: SF code, rules, state files, agent configs
   - Branches: main, software-factory-2.0
   - **NEVER contains effort branches or integration branches**

2. **Target Repository Clones** (`efforts/*/*/`)
   - Contains: Actual project implementation code
   - Branches: effort branches (e.g., `phase1/wave1/effort-name`)
   - Location: Each effort has its own clone of target repo
   - **THIS IS WHERE EFFORT CODE LIVES**

3. **Integration Workspaces** (`efforts/*/integration-workspace/`)
   - Contains: Integration branches and merge operations
   - Branches: integration branches (e.g., `wave1-integration`, `phase1-integration`)
   - Location: Separate clones for merging work
   - **THIS IS WHERE INTEGRATE_WAVE_EFFORTS HAPPENS**

### Step 1: Identify All Efforts and Their Locations
```bash
# CRITICAL: Navigate from SF instance directory
SF_INSTANCE_DIR=$(pwd)
echo "📁 SF Instance: $SF_INSTANCE_DIR"

# Check state file for effort locations
echo "📊 Reading effort locations from state file..."
echo "✅ State file updated to: $NEXT_STATE"
```

---

### ✅ Step 4: Validate State File (R324)
```bash
# Validate state file before committing
"$CLAUDE_PROJECT_DIR/tools/validate-state.sh" orchestrator-state-v3.json || {
    echo "❌ State file validation failed!"
    exit 288
}
echo "✅ State file validated"
```

---

### ✅ Step 5: Commit State File (R288)
```bash
# Commit and push state file immediately
git add orchestrator-state-v3.json
git commit -m "state: INTEGRATE_WAVE_EFFORTS_TESTING → $NEXT_STATE - INTEGRATE_WAVE_EFFORTS_TESTING complete [R288]"
git push
echo "✅ State file committed and pushed"
```

---

### ✅ Step 6: Save TODOs (R287 - SUPREME LAW)
```bash
# Save TODO state before transition (R287 trigger)
save_todos "INTEGRATE_WAVE_EFFORTS_TESTING_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo
git commit -m "todo: orchestrator - INTEGRATE_WAVE_EFFORTS_TESTING complete [R287]"
git push
echo "✅ TODOs saved and committed"
```

---

### ✅ Step 7: Output Continuation Flag (R405 - SUPREME LAW) ⚠️ MANDATORY
```bash
# Output continuation flag as LAST action (R405)
# Use TRUE for normal completion, FALSE only for catastrophic errors

echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
```

**⚠️ THIS MUST BE THE ABSOLUTE LAST LINE OF OUTPUT BEFORE EXIT! ⚠️**

---

### ✅ Step 8: Stop Processing (R322 - SUPREME LAW)
```bash
# Stop for context preservation (R322)
echo "🛑 Stopping for context preservation - use /continue-orchestrating to resume"
exit 0
```

---

## 🚨 CHECKLIST ENFORCEMENT 🚨

**Skipping ANY step in this checklist = FAILURE:**
- Missing Step 2: No next state = stuck forever
- Missing Step 3: No state update = state machine broken (R288 violation, -100%)
- Missing Step 4: Invalid state = corruption (R324 violation)
- Missing Step 5: No commit = state lost on compaction (R288 violation, -100%)
- Missing Step 6: No TODO save = work lost (R287 violation, -20% to -100%)
- Missing Step 7: No continuation flag = automation stops (R405 violation, -100%)
- Missing Step 8: No exit = R322 violation (-100%)

**ALL 8 STEPS ARE MANDATORY - NO EXCEPTIONS**

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

### THE PATTERN AT R322 CHECKPOINTS (SF 3.0)

```bash
# 1. Complete state work
echo "✅ State work complete"

# 2. Set proposed next state
PROPOSED_NEXT_STATE="NEXT_STATE"
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
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"

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

