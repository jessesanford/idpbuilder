# Orchestrator - CREATE_INTEGRATE_WAVE_EFFORTS_TESTING State Rules

## 🔴🔴🔴 ABSOLUTE REQUIREMENT: STATE MANAGER CONSULTATION 🔴🔴🔴

**SUPREME LAW - R517 - UNIVERSAL STATE MANAGER CONSULTATION LAW**

**BEFORE EXITING THIS STATE, YOU MUST:**

1. **SPAWN STATE MANAGER FOR SHUTDOWN_CONSULTATION** (MANDATORY - NO EXCEPTIONS)
2. **NEVER UPDATE STATE FILES DIRECTLY** (orchestrator-state-v3.json, bug-tracking.json, etc.)
3. **NEVER COMMIT STATE FILES YOURSELF** (State Manager does this atomically)

**FAILURE TO CONSULT STATE MANAGER = IMMEDIATE SYSTEM HALT (-100% GRADE)**

### Enforcement Mechanism

If you attempt to exit this state without spawning State Manager:
- ❌ Pre-commit hooks will REJECT your commit
- ❌ Validation tools will FAIL the build
- ❌ Grading system will assign -100% penalty
- ❌ System will transition to ERROR_RECOVERY

### Required Pattern (COPY THIS EXACTLY)

```bash
# At end of state work, BEFORE any state file updates:

echo "🔴 MANDATORY: Spawning State Manager for SHUTDOWN_CONSULTATION"

# Spawn State Manager (REQUIRED - NOT OPTIONAL)
# Task: state-manager
# State: SHUTDOWN_CONSULTATION
# Current State: [YOUR_CURRENT_STATE]
# Proposed Next State: [YOUR_PROPOSED_NEXT_STATE]
# Work Summary: [SUMMARY_OF_WORK_COMPLETED]

# State Manager will:
# 1. Validate proposed transition against state machine
# 2. Update all 4 state files atomically
# 3. Commit with [R288] tag
# 4. Return REQUIRED next state (may differ from proposal)

# Wait for State Manager response
# Follow State Manager's directive (REQUIRED next state)
# DO NOT proceed until State Manager confirms
```

**YOU MUST NEVER:**
- ❌ Update orchestrator-state-v3.json yourself
- ❌ Update bug-tracking.json yourself
- ❌ Update integration-containers.json yourself
- ❌ Use `jq` to modify state files
- ❌ Use `sed/awk` to modify state files
- ❌ Set `validated_by: "orchestrator"` (must be "state-manager")
- ❌ Skip State Manager consultation "just this once"
- ❌ Think "I'll validate it manually"

**ONLY State Manager may update state files. This is NON-NEGOTIABLE.**

See: `rule-library/R517-universal-state-manager-consultation-law.md`

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

**YOU HAVE ENTERED CREATE_INTEGRATE_WAVE_EFFORTS_TESTING STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_CREATE_INTEGRATE_WAVE_EFFORTS_TESTING
echo "$(date +%s) - Rules read and acknowledged for CREATE_INTEGRATE_WAVE_EFFORTS_TESTING" > .state_rules_read_orchestrator_CREATE_INTEGRATE_WAVE_EFFORTS_TESTING
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY INTEGRATE_WAVE_EFFORTS TESTING SETUP WORK UNTIL RULES ARE READ:
- ❌ Start creating integration-testing branch
- ❌ Start cloning repositories
- ❌ Start setting up workspaces
- ❌ Update state files
- ❌ Continue to next state
- ❌ Think about what to do in this state

### ✅ YOU MUST IMMEDIATELY:

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

### ❌ ANTI-PATTERNS THAT WILL CAUSE FAILURE:

1. **Fake Acknowledgment Without Reading**:
   ```
   ❌ WRONG: "I acknowledge R272, R271, R273..."
   (No Read tool calls detected = AUTOMATIC FAILURE)
   ```

2. **Bulk Acknowledgment**:
   ```
   ❌ WRONG: "I acknowledge all CREATE_INTEGRATE_WAVE_EFFORTS_TESTING rules"
   (YOU Must READ AND ACKNOWLEDGE EACH rule individually)
   ```

3. **Silent Reading**:
   ```
   ❌ WRONG: [Reads rules but doesn't acknowledge]
   "Now I've read the rules, let me start work..."
   (MUST explicitly acknowledge EACH rule)
   ```

4. **Reading From Memory**:
   ```
   ❌ WRONG: "I know R272 requires creating from main HEAD..."
   (Must READ from file, not recall from memory)
   ```

5. **Skipping Rules in PRIMARY DIRECTIVES**:
   ```
   ❌ WRONG: Reading only some rules from the list
   (ALL rules in PRIMARY DIRECTIVES are MANDATORY)
   ```

### ✅ CORRECT PATTERN FOR CREATE_INTEGRATE_WAVE_EFFORTS_TESTING:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/[first-rule-file].md
2. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/[second-rule-file].md  
4. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. "Ready to execute CREATE_INTEGRATE_WAVE_EFFORTS_TESTING work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY CREATE_INTEGRATE_WAVE_EFFORTS_TESTING work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ You have stated readiness to execute CREATE_INTEGRATE_WAVE_EFFORTS_TESTING work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY CREATE_INTEGRATE_WAVE_EFFORTS_TESTING work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

**THE SYSTEM IS MONITORING_SWE_PROGRESS YOUR READ TOOL CALLS!**

## 📋 PRIMARY DIRECTIVES FOR CREATE_INTEGRATE_WAVE_EFFORTS_TESTING STATE

### 🚨🚨🚨 R006 - Orchestrator NEVER Writes Code [BLOCKING]
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`
**Criticality**: BLOCKING - Any code operation = -100% IMMEDIATE FAILURE
**Summary**: Orchestrator is a COORDINATOR ONLY - never writes, edits, or modifies code

**⚠️ R006 WARNING FOR CREATE_INTEGRATE_WAVE_EFFORTS_TESTING STATE:**
- DO NOT create test files yourself!
- DO NOT write integration test code!
- DO NOT modify any source files!
- You only set up infrastructure - SW Engineers write ALL code

### 🚨🚨🚨 R272 - Integration Testing Branch Requirement
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R272-integration-testing-branch.md`
**Criticality**: BLOCKING - Must create from main HEAD
**Summary**: Create dedicated integration-testing branch from main's current HEAD

### 🚨🚨🚨 R271 - Mandatory Production-Ready Validation
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R271-mandatory-production-ready-validation.md`
**Criticality**: BLOCKING - Full checkouts required
**Summary**: Must use full repository clones, no sparse checkouts

### 🚨🚨🚨 R014 - Branch Naming Convention
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R014-branch-naming-convention.md`
**Criticality**: BLOCKING - Project prefix required
**Summary**: Use project prefix for integration-testing branch

### 🚨🚨🚨 R251 - Repository Separation Law
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R251-REPOSITORY-SEPARATION-LAW.md`
**Criticality**: BLOCKING - SF/target repo isolation
**Summary**: Integration testing happens in target repository only

### 🚨🚨🚨 R280 - Main Branch Protection
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R280-main-branch-protection.md`
**Criticality**: SUPREME LAW - Never modify main
**Summary**: Software Factory NEVER pushes to main branch

## 🚨 CREATE_INTEGRATE_WAVE_EFFORTS_TESTING IS A VERB - CREATE INFRASTRUCTURE NOW! 🚨

### IMMEDIATE ACTIONS UPON ENTERING CREATE_INTEGRATE_WAVE_EFFORTS_TESTING

**THE MOMENT YOU ENTER THIS STATE, YOU MUST:**
1. Checkout main branch and pull latest
2. Create timestamped integration-testing branch from main HEAD
3. Setup integration testing workspace
4. Document branch creation details
5. Transition to INTEGRATE_WAVE_EFFORTS_TESTING state

**FORBIDDEN - AUTOMATIC FAILURE:**
- ❌ "STATE TRANSITION COMPLETE: Now in CREATE_INTEGRATE_WAVE_EFFORTS_TESTING" [stops]
- ❌ "Successfully entered CREATE_INTEGRATE_WAVE_EFFORTS_TESTING state" [waits]
- ❌ "Ready to create integration testing branch" [pauses]
- ❌ "I'm in CREATE_INTEGRATE_WAVE_EFFORTS_TESTING state" [does nothing]

**REQUIRED - IMMEDIATE ACTION:**
- ✅ "Entering CREATE_INTEGRATE_WAVE_EFFORTS_TESTING, creating branch from main HEAD NOW..."
- ✅ "START INTEGRATE_WAVE_EFFORTS TESTING SETUP, pulling latest main..."
- ✅ "CREATE_INTEGRATE_WAVE_EFFORTS_TESTING: Creating timestamped branch..."

## State Context
You are creating the final integration testing infrastructure where the PROJECT INTEGRATE_WAVE_EFFORTS branch will be merged to prove the entire project works. This is the ultimate validation before generating the MASTER-PR-PLAN for humans.

## 🔴🔴🔴 CRITICAL: INTEGRATE_WAVE_EFFORTS TESTING USES PROJECT INTEGRATE_WAVE_EFFORTS (R283) 🔴🔴🔴

**Prerequisites:**
- PROJECT_INTEGRATE_WAVE_EFFORTS has been completed (R283)
- All phases have been merged into project integration branch
- Code Reviewer has validated the project integration
- Ready for final integration testing

**This is NOT merging individual efforts - this is PROJECT-WIDE:**
1. **START** from main branch's current HEAD (not any integration branch)
2. **CREATE** timestamped integration-testing branch
3. **PREPARE** to merge the PROJECT INTEGRATE_WAVE_EFFORTS branch (which contains ALL phases)
4. **VALIDATE** entire project functionality
5. **PROVE** software is production-ready

**MERGE APPROACH (R283 compliant):**
- ✅ Merge the single project-integration branch (contains all phases)
- ❌ DO NOT merge individual effort branches
- ❌ DO NOT merge phase integration branches separately
- ❌ DO NOT merge wave integration branches

**ALWAYS BASE ON:**
- ✅ main branch HEAD (current state)
- ✅ Then merge project-integration branch into it

## Integration Testing Infrastructure Setup

### 🚨🚨🚨 R272 Compliance - Create From Main HEAD
**SEE**: `$CLAUDE_PROJECT_DIR/rule-library/R272-integration-testing-branch.md`

```bash
# 🔴 CRITICAL: Integration testing setup script
create_integration_testing_infrastructure() {
    echo "🏭 CREATE_INTEGRATE_WAVE_EFFORTS_TESTING: Starting infrastructure setup..."
    
    # 0. Save SF instance directory
    SF_INSTANCE_DIR=$(pwd)
    echo "📁 SF Instance: $SF_INSTANCE_DIR"
    
    # 1. Source branch naming helpers
    source "$SF_INSTANCE_DIR/utilities/branch-naming-helpers.sh"
    PROJECT_PREFIX=$(echo "✅ State file updated to: $NEXT_STATE"
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
git commit -m "state: CREATE_INTEGRATE_WAVE_EFFORTS_TESTING → $NEXT_STATE - CREATE_INTEGRATE_WAVE_EFFORTS_TESTING complete [R288]"
git push
echo "✅ State file committed and pushed"
```

---

### ✅ Step 6: Save TODOs (R287 - SUPREME LAW)
```bash
# Save TODO state before transition (R287 trigger)
save_todos "CREATE_INTEGRATE_WAVE_EFFORTS_TESTING_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo
git commit -m "todo: orchestrator - CREATE_INTEGRATE_WAVE_EFFORTS_TESTING complete [R287]"
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

