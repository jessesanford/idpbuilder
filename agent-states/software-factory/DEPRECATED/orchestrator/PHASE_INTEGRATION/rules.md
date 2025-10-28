# Orchestrator - INTEGRATE_PHASE_WAVES State Rules


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

**YOU HAVE ENTERED INTEGRATE_PHASE_WAVES STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_INTEGRATE_PHASE_WAVES
echo "$(date +%s) - Rules read and acknowledged for INTEGRATE_PHASE_WAVES" > .state_rules_read_orchestrator_INTEGRATE_PHASE_WAVES
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY INTEGRATE_PHASE_WAVES WORK UNTIL RULES ARE READ:
- ❌ Start merge wave branches
- ❌ Start create phase branch
- ❌ Start integrate wave work
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
   ❌ WRONG: "I acknowledge R151, R208, R053..."
   (No Read tool calls detected = AUTOMATIC FAILURE)
   ```

2. **Bulk Acknowledgment**:
   ```
   ❌ WRONG: "I acknowledge all INTEGRATE_PHASE_WAVES rules"
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
   ❌ WRONG: "I know R208 requires CD before spawn..."
   (Must READ from file, not recall from memory)
   ```

5. **Skipping Rules in PRIMARY DIRECTIVES**:
   ```
   ❌ WRONG: Reading only some rules from the list
   (ALL rules in PRIMARY DIRECTIVES are MANDATORY)
   ```

### ✅ CORRECT PATTERN FOR INTEGRATE_PHASE_WAVES:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/[first-rule-file].md
2. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/[second-rule-file].md  
4. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. "Ready to execute INTEGRATE_PHASE_WAVES work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY INTEGRATE_PHASE_WAVES work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ You have stated readiness to execute INTEGRATE_PHASE_WAVES work
1. **READ** every rule file listed in PRIMARY DIRECTIVES below
2. **ACKNOWLEDGE** each rule individually with number and description
3. **ONLY THEN** proceed with INTEGRATE_PHASE_WAVES work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY INTEGRATE_PHASE_WAVES work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

**THE SYSTEM IS MONITORING_SWE_PROGRESS YOUR READ TOOL CALLS!**

## ⚠️⚠️⚠️ MANDATORY RULE READING AND ACKNOWLEDGMENT ⚠️⚠️⚠️

**YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. YOUR READ TOOL CALLS ARE BEING MONITORED.**

### ❌ ANTI-PATTERNS THAT WILL CAUSE FAILURE:
1. Fake acknowledgment without reading
2. Bulk acknowledgment
3. Reading from memory

### ✅ CORRECT PATTERN:
1. READ each rule file
2. Acknowledge individually with rule number and description

## 📋 PRIMARY DIRECTIVES FOR INTEGRATE_PHASE_WAVES STATE

## 🔴🔴🔴 MANDATORY MEASUREMENT RULE - R304 🔴🔴🔴

**ABSOLUTE REQUIREMENTS:**
- ✅ MUST use `$CLAUDE_PROJECT_DIR/tools/line-counter.sh` for ALL line counting
- ❌ NEVER use `wc -l` or any manual counting method
- ❌ NEVER count lines any other way - this is a -100% automatic failure
- ✅ MUST specify both -b (base branch) and -c (current branch) parameters
- ✅ Base branch MUST be phase integration branch (NOT "main")

**Failure to use the line counter tool = AUTOMATIC -100% GRADE**

### 🔴🔴🔴 R301 - Integration Branch Current Tracking (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R301-integration-branch-current-tracking.md`
**Criticality**: SUPREME LAW - Only ONE current integration allowed
**Summary**: Track current vs deprecated integrations, prevent wrong branch usage

### 🚨🚨🚨 R006 - Orchestrator NEVER Writes Code [BLOCKING]
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`
**Criticality**: BLOCKING - Any code operation = -100% IMMEDIATE FAILURE
**Summary**: Orchestrator coordinates but NEVER implements or fixes code

### 🚨🚨🚨 R329 - Orchestrator NEVER Performs Git Merges [BLOCKING]
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R329-orchestrator-never-performs-merges.md`
**Criticality**: BLOCKING - Any merge operation = -100% IMMEDIATE FAILURE
**Summary**: Orchestrator MUST spawn Integration Agent for ALL merges - NO EXCEPTIONS

### 🚨🚨🚨 R285 - Mandatory Phase Integration Before Assessment  
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R285-mandatory-phase-integration-before-assessment.md`
**Criticality**: BLOCKING - Must integrate before assessment
**Summary**: Phase integration required in normal flow (from REVIEW_WAVE_ARCHITECTURE)

### 🚨🚨🚨 R259 - Mandatory Phase Integration After Fixes
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R259-mandatory-phase-integration-after-fixes.md`
**Criticality**: BLOCKING - Must create integration branch after fixes
**Summary**: Create phase-level integration after ERROR_RECOVERY fixes

### 🚨🚨🚨 R257 - Mandatory Phase Assessment Report
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R257-mandatory-phase-assessment-report.md`
**Criticality**: BLOCKING - Required for phase completion
**Summary**: Verify all assessment issues are addressed

### 🚨🚨🚨 R014 - Branch Naming Convention
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R014-branch-naming-convention.md`
**Criticality**: BLOCKING - Mandatory project prefix for all branches
**Summary**: Use project prefix for phase integration branches

### 🚨🚨🚨 R296 - Deprecated Branch Marking Protocol
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R296-deprecated-branch-marking-protocol.md`
**Criticality**: BLOCKING - Prevents integration of wrong branches
**Summary**: Check for and prevent integration of deprecated split branches

### 🚨🚨🚨 R271 - Mandatory Production-Ready Validation
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R271-mandatory-production-ready-validation.md`
**Criticality**: BLOCKING - Full checkouts required
**Summary**: Phase integration must use full repository clones

### 🚨🚨🚨 R269 - Code Reviewer Merge Plan No Execution
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R269-code-reviewer-merge-plan-no-execution.md`
**Criticality**: BLOCKING - Code Reviewer only plans
**Summary**: Code Reviewer creates plan, Integration Agent executes

### 🚨🚨🚨 R260 - Integration Agent Core Requirements
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R260-integration-agent-core-requirements.md`
**Criticality**: BLOCKING - Integration Agent requirements
**Summary**: Integration Agent must acknowledge INTEGRATE_WAVE_EFFORTS_DIR

### 🔴🔴🔴 R321 - Immediate Backport During Integration Protocol
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R321-immediate-backport-during-integration.md`
**Criticality**: SUPREME LAW - Immediate backporting required
**Summary**: ANY fix during integration MUST be immediately backported to source branches before continuing

### 🚨🚨🚨 R280 - Main Branch Protection Protocol [BLOCKING]
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R280-main-branch-protection.md`
**Criticality**: BLOCKING - Direct commits to main/master are forbidden
**Summary**: All changes must go through PR process with proper reviews

### 🚨🚨🚨 R307 - Branch Mergeability Check [BLOCKING]
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R307-independent-branch-mergeability.md`
**Criticality**: BLOCKING - Must verify branches are mergeable before attempting
**Summary**: Check for conflicts and mergeability before integration operations

### 🔴🔴🔴 R233 - All States Require Immediate Action (CRITICAL)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R233-all-states-immediate-action.md`
**Criticality**: CRITICAL - States are verbs
**Summary**: INTEGRATE_PHASE_WAVES means START INTEGRATING NOW

### 🔴🔴🔴 R288 - State File Update and Commit Protocol (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
**Criticality**: SUPREME LAW - Update on every transition
**Summary**: Update orchestrator-state-v3.json with integration details

### 🚨🚨🚨 R288 - State File Update and Commit Protocol
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
**Criticality**: BLOCKING - Push within 60 seconds
**Summary**: Commit and push state file immediately

## 🚨 INTEGRATE_PHASE_WAVES IS A VERB - COORDINATE PHASE INTEGRATE_WAVE_EFFORTS NOW! 🚨

### IMMEDIATE ACTIONS UPON ENTERING INTEGRATE_PHASE_WAVES

**THE MOMENT YOU ENTER THIS STATE, YOU MUST:**
1. Check if phase integration infrastructure exists NOW
2. If NO infrastructure: Transition to SETUP_INTEGRATE_PHASE_WAVES_INFRASTRUCTURE
3. If infrastructure EXISTS: Transition to SPAWN_CODE_REVIEWER_PHASE_MERGE_PLAN
4. Update state file with the appropriate next state
5. Stop per R322 for state transition

**FORBIDDEN - AUTOMATIC FAILURE:**
- ❌ "STATE TRANSITION COMPLETE: Now in INTEGRATE_PHASE_WAVES" [stops]
- ❌ "Successfully entered INTEGRATE_PHASE_WAVES state" [waits]
- ❌ "Ready to start phase integration" [pauses]
- ❌ "I'm in INTEGRATE_PHASE_WAVES state" [does nothing]
- ❌ Creating infrastructure yourself (INTEGRATE_PHASE_WAVES only coordinates!)
- ❌ Merging branches yourself (R329 violation!)

**REQUIRED - IMMEDIATE ACTION:**
- ✅ "INTEGRATE_PHASE_WAVES STATE: Checking for existing phase integration infrastructure..."
- ✅ "No infrastructure found, transitioning to SETUP_INTEGRATE_PHASE_WAVES_INFRASTRUCTURE..."
- ✅ "Infrastructure exists, transitioning to SPAWN_CODE_REVIEWER_PHASE_MERGE_PLAN..."

## Primary Purpose

**THIS STATE IS FOR COORDINATION ONLY!**

The INTEGRATE_PHASE_WAVES state is a decision point that:
1. **CHECKS** if phase integration infrastructure exists
2. **TRANSITIONS** to SETUP_INTEGRATE_PHASE_WAVES_INFRASTRUCTURE if no infrastructure
3. **TRANSITIONS** to SPAWN_CODE_REVIEWER_PHASE_MERGE_PLAN if infrastructure exists

**THIS STATE NEVER:**
- ❌ Creates phase integration workspace itself
- ❌ Sets up branches or directories itself
- ❌ Performs any actual integration work
- ❌ Executes git merges (R329 violation!)

You are the COORDINATOR of phase integration flow.

## 🔴 CRITICAL: Locating Effort Branches for Integration

### Effort Branch Locations (Per R193/R191)
All effort branches are located in specific directories with predictable patterns:

#### Directory Structure:
```bash
# Effort workspaces follow this pattern:
/efforts/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}/

# Example for Phase 2, Wave 1 with 3 efforts:
/efforts/phase2/wave1/auth-system/       # Contains effort branch
/efforts/phase2/wave1/user-management/   # Contains effort branch
/efforts/phase2/wave1/api-gateway/       # Contains effort branch
```

#### Branch Naming Convention:
```bash
# Effort branches follow naming from target-repo-config.yaml:
# Pattern: phase${PHASE}-wave${WAVE}-${EFFORT_NAME}
# Or with project prefix: ${PREFIX}/phase${PHASE}-wave${WAVE}-${EFFORT_NAME}

# Examples without prefix:
phase2-wave1-auth-system
phase2-wave1-user-management
phase2-wave1-api-gateway

# Examples with prefix (e.g., "tmc-workspace"):
tmc-workspace/phase2-wave1-auth-system
tmc-workspace/phase2-wave1-user-management
tmc-workspace/phase2-wave1-api-gateway
```

### Finding Efforts to Integrate

**MANDATORY: Before integration, locate all effort branches:**

```bash
#!/bin/bash
# Script to find all effort branches for current phase

PHASE=$(echo "✅ State file updated to: $NEXT_STATE"
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
git commit -m "state: INTEGRATE_PHASE_WAVES → $NEXT_STATE - INTEGRATE_PHASE_WAVES complete [R288]"
git push
echo "✅ State file committed and pushed"
```

---

### ✅ Step 6: Save TODOs (R287 - SUPREME LAW)
```bash
# Save TODO state before transition (R287 trigger)
save_todos "INTEGRATE_PHASE_WAVES_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo
git commit -m "todo: orchestrator - INTEGRATE_PHASE_WAVES complete [R287]"
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

