# Orchestrator - BUILD_VALIDATION State Rules


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

**YOU HAVE ENTERED BUILD_VALIDATION STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_BUILD_VALIDATION
echo "$(date +%s) - Rules read and acknowledged for BUILD_VALIDATION" > .state_rules_read_orchestrator_BUILD_VALIDATION
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY VALIDATION WORK UNTIL RULES ARE READ:
- ❌ Start running tests
- ❌ Start checking dependencies
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

### 🚨🚨🚨 RULE R319 - Orchestrator NEVER Measures or Assesses Code [BLOCKING]
**Source:** `$CLAUDE_PROJECT_DIR/rule-library/R319-orchestrator-never-measures-code.md`
**Criticality:** BLOCKING - Any technical assessment = -100% IMMEDIATE FAILURE

READ THIS RULE FILE: `$CLAUDE_PROJECT_DIR/rule-library/R319-orchestrator-never-measures-code.md`

### 🚨🚨🚨 RULE R323 - Mandatory Final Artifact Build [BLOCKING]
**Source:** `$CLAUDE_PROJECT_DIR/rule-library/R323-mandatory-final-artifact-build.md`
**Criticality:** BLOCKING - No artifact = -50% to -100% FAILURE

READ THIS RULE FILE: `$CLAUDE_PROJECT_DIR/rule-library/R323-mandatory-final-artifact-build.md`

**⚠️ R006, R319 & R323 WARNING FOR BUILD_VALIDATION STATE:**
- DO NOT run tests yourself - spawn Code Reviewer!
- DO NOT execute builds yourself - spawn Code Reviewer!
- DO NOT validate dependencies yourself - spawn Code Reviewer!
- DO NOT fix test failures yourself - spawn SW Engineers!
- DO NOT modify code to make tests pass - spawn SW Engineers!
- **🔴 R323: Code Reviewer MUST verify final artifact still exists**
- **🔴 R323: Code Reviewer MUST test the actual built artifact**
- You ONLY coordinate validation - NEVER execute it yourself

### Production Ready Validation Requirements

**This state validates that the integrated code is ready for production deployment.**

Key validation areas:
- All tests must pass
- Dependencies must be secure and up-to-date
- Build must complete successfully
- **🔴 R323: Final artifact must exist and be tested**
- **🔴 R323: Artifact path must be documented**
- Documentation must be present

## 🎯 STATE OBJECTIVES

In the BUILD_VALIDATION state, you are the COORDINATOR:

1. **Spawning Code Reviewer for Test Validation**
   - Spawn Code Reviewer to run all test suites
   - Provide clear instructions on what to validate
   - Specify test types to execute
   - Request comprehensive test report

2. **Spawning Code Reviewer for Dependency Validation**
   - Spawn Code Reviewer to audit dependencies
   - Request security vulnerability checks
   - Ask for version compatibility verification
   - Ensure dependency resolution validation

3. **Monitoring Validation Progress**
   - Check if Code Reviewer has completed
   - Read validation reports
   - Track any issues identified
   - Coordinate fixes if needed

4. **Managing State Transitions**
   - Determine next state based on validation
   - Update state file with validation status
   - Track production readiness
   - Coordinate any required fixes

## 📝 REQUIRED ACTIONS

### Step 1: Spawn Code Reviewer for Test Validation
```bash
# ✅ CORRECT: Delegate test validation to Code Reviewer
echo "🧪 Production readiness validation needed"
echo "🚀 Spawning Code Reviewer to run test suites..."

cd /efforts/integration-testing

# Update state to show spawning Code Reviewer
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
git commit -m "state: BUILD_VALIDATION → $NEXT_STATE - BUILD_VALIDATION complete [R288]"
git push
echo "✅ State file committed and pushed"
```

---

### ✅ Step 6: Save TODOs (R287 - SUPREME LAW)
```bash
# Save TODO state before transition (R287 trigger)
save_todos "BUILD_VALIDATION_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo
git commit -m "todo: orchestrator - BUILD_VALIDATION complete [R287]"
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

