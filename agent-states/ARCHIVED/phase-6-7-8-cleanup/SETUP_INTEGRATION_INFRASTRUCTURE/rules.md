# Orchestrator - SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE State Rules

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

## 🔴🔴🔴 CRITICAL: INTEGRATE_WAVE_EFFORTS BASE BRANCH DETERMINATION (R308) 🔴🔴🔴

**VIOLATION = -100% AUTOMATIC FAILURE**

### YOU MUST DETERMINE THE CORRECT BASE BRANCH (SEQUENTIAL REBUILD MODEL):
- **Wave Integration**: Base on the FIRST EFFORT of THIS wave (R009 Sequential Rebuild)
- **Phase Integration**: Base on the FIRST EFFORT of THIS phase (R282 Sequential Rebuild)
- **Project Integration**: Base on MAIN - the trunk! (R283 Sequential Rebuild)

### NOTE: R308 IS ABOUT EFFORT BRANCHING, NOT INTEGRATE_WAVE_EFFORTS BASES!
**R308 defines how EFFORTS cascade** (development flow):
```
Efforts within waves cascade from each other:
P1W1: main → effort1 → effort2 → effort3
P1W2: effort3 → effort4 → effort5
P2W1: effort5 → effort6 → effort7
```

**Integration bases** are determined by R009/R282/R283 Sequential Rebuild:
```
Wave integration: base = first effort of WAVE
Phase integration: base = first effort of PHASE
Project integration: base = main
```

### 🔴 CRITICAL EXAMPLE:
```bash
# Phase 2 Integration (after completing all Phase 2 waves)
# WRONG - AUTOMATIC FAILURE:
BASE_BRANCH="main"  # ❌ WRONG! Only for project integration!
BASE_BRANCH="phase1-integration"  # ❌ WRONG! Not from previous phase!
BASE_BRANCH="phase2-wave3-integration"  # ❌ WRONG! Wave integrations are testing checkpoints (R364)!

# CORRECT (R282 Sequential Rebuild):
BASE_BRANCH="phase2/wave1/auth-system"  # ✅ First effort of THIS phase!
```

**Acknowledge: "I understand integration bases MUST follow Sequential Rebuild model (R009/R282/R283), NOT R308 (which is for effort cascading)"**

## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE
echo "$(date +%s) - Rules read and acknowledged for SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE" > .state_rules_read_orchestrator_SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE WORK UNTIL RULES ARE READ:
- ❌ Start creating integration workspace
- ❌ Start creating integration branch
- ❌ Start cloning repository
- ❌ Start pushing to remote
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
   ❌ WRONG: "I acknowledge R308, R250, R034..."
   (No Read tool calls detected = AUTOMATIC FAILURE)
   ```

2. **Bulk Acknowledgment**:
   ```
   ❌ WRONG: "I acknowledge all SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE rules"
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
   ❌ WRONG: "I know R308 requires incremental bases..."
   (Must READ from file, not recall from memory)
   ```

5. **Skipping Rules in PRIMARY DIRECTIVES**:
   ```
   ❌ WRONG: Reading only some rules from the list
   (ALL rules in PRIMARY DIRECTIVES are MANDATORY)
   ```

### ✅ CORRECT PATTERN FOR SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/[first-rule-file].md
2. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/[second-rule-file].md  
4. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. "Ready to execute SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ You have stated readiness to execute SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE work
1. **READ** every rule file listed in PRIMARY DIRECTIVES below
2. **ACKNOWLEDGE** each rule individually with number and description
3. **ONLY THEN** proceed with SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

**THE SYSTEM IS MONITORING_SWE_PROGRESS YOUR READ TOOL CALLS!**

## 📋 PRIMARY DIRECTIVES FOR SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE

**YOU MUST READ EACH RULE LISTED HERE. YOUR READ TOOL CALLS ARE BEING MONITORED.**

### 🔴🔴🔴 R308 - INCREMENTAL BRANCHING STRATEGY (SUPREME LAW!)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R308-incremental-branching-strategy.md`
**Criticality**: 🔴🔴🔴 SUPREME LAW - Integration branches MUST be incremental!
**Summary**: Wave N+1 integration based on Wave N integration, Phase N+1 based on Phase N
**CRITICAL**: Phase 2 Wave 1 integration MUST use phase1-integration, NOT main!

### 🚨🚨🚨 R250 - Integration Isolation Requirement
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R250-integration-isolation-requirement.md`
**Criticality**: BLOCKING - Integration must use separate target clone
**Summary**: Integration must happen under /efforts/ directory structure

### 🚨🚨🚨 R034 - Integration Requirements
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R034-integration-requirements.md`
**Criticality**: BLOCKING - Required for wave approval
**Summary**: Complete integration protocol with testing and validation

### 🚨🚨🚨 R014 - Branch Naming Convention
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R014-branch-naming-convention.md`
**Criticality**: BLOCKING - Mandatory project prefix for all branches
**Summary**: Use project prefix for all integration branches

### 🚨🚨🚨 R271 - Mandatory Production-Ready Validation
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R271-mandatory-production-ready-validation.md`
**Criticality**: BLOCKING - Full checkouts required for integration
**Summary**: Integration must use full repository clones, no sparse checkouts

### 🚨🚨🚨 R006 - Orchestrator NEVER Writes Code [BLOCKING]
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`
**Criticality**: BLOCKING - Any code operation = -100% IMMEDIATE FAILURE
**Summary**: Orchestrator coordinates but NEVER implements or fixes code

### 🚨🚨🚨 R329 - Orchestrator NEVER Performs Git Merges [BLOCKING]
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R329-orchestrator-never-performs-merges.md`
**Criticality**: BLOCKING - Any merge operation = -100% IMMEDIATE FAILURE
**Summary**: Orchestrator MUST spawn Integration Agent for ALL merges - NO EXCEPTIONS

### 🚨🚨🚨 R307 - Independent Branch Mergeability [BLOCKING]
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R307-independent-branch-mergeability.md`
**Criticality**: BLOCKING - Must verify branches are mergeable before attempting
**Summary**: Check for conflicts and mergeability before integration operations

### 🚨🚨🚨 R216 - Bash Execution Syntax Protocol (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R216-bash-execution-syntax.md`
**Criticality**: BLOCKING - Incorrect syntax causes failures
**Summary**: Use parentheses for subshells, proper variable syntax

### 🚨🚨🚨 R288 - State File Update and Commit Protocol
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
**Criticality**: SUPREME LAW - Update on every transition
**Summary**: Update orchestrator-state-v3.json with integration metadata

## 🚨 SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE IS A VERB - CREATE INFRASTRUCTURE NOW! 🚨

### 🔴🔴🔴 CRITICAL: YOU ARE ALREADY IN SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE STATE! 🔴🔴🔴

**If current_state = "SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE" in orchestrator-state-v3.json, you MUST:**
1. **IMMEDIATELY** start creating integration infrastructure
2. **NO ANNOUNCEMENTS** - just start working
3. **NO WAITING** - infrastructure creation begins NOW

### IMMEDIATE ACTIONS UPON ENTERING SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE

**THE MOMENT YOU SEE current_state: SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE, YOU MUST:**
1. Determine the integration type (wave/phase/project) NOW
2. Determine the correct incremental base branch per R308
3. Create integration working directory immediately
4. Clone repository with FULL checkout (R271)
5. Create integration branch following R308
6. Push integration branch to remote
7. Update state file with integration metadata
8. Transition to appropriate next state

**FORBIDDEN - AUTOMATIC FAILURE:**
- ❌ "STATE TRANSITION COMPLETE: Now in SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE" [stops]
- ❌ "Successfully entered SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE state" [waits]
- ❌ "Ready to setup integration infrastructure" [pauses]
- ❌ "I'm in integration infrastructure state" [does nothing]
- ❌ "Preparing to create integration workspace..." [delays]
- ❌ "I see we're in SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE state..." [announces]

**REQUIRED - IMMEDIATE ACTION:**
- ✅ "SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE: Determining integration type NOW..."
- ✅ "Creating wave integration infrastructure at /efforts/phase${X}/wave${Y}/integration-workspace..."
- ✅ "Applying R308: Phase 2 Wave 1 integration will use phase1-integration as base..."

## State Context
You are CREATING INFRASTRUCTURE for integration, not executing merges. Your responsibilities:
1. **DETERMINE** correct incremental base branch per R308
2. **CREATE** integration workspace directory
3. **CLONE** target repository with FULL checkout
4. **CREATE** integration branch with proper naming
5. **PUSH** integration branch to establish remote
6. **UPDATE** state file with integration metadata
7. **TRANSITION** to next appropriate state

**YOU MUST NEVER (R329 + R006 ENFORCEMENT):**
- ❌ Execute any git merges (R329 VIOLATION)
- ❌ Resolve any conflicts (R329 VIOLATION)
- ❌ Run any builds or tests (R006 VIOLATION)
- ❌ Write any code (R006 VIOLATION)

## 🔴🔴🔴 CRITICAL: SEQUENTIAL REBUILD BASE DETERMINATION (R009/R282/R283) 🔴🔴🔴

### Integration Base Branch Function
```bash
determine_integration_base_branch() {
    local INTEGRATE_WAVE_EFFORTS_TYPE="$1"  # wave, phase, or project
    local PHASE="$2"
    local WAVE="$3"

    echo "🔴 Sequential Rebuild: Determining base for $INTEGRATE_WAVE_EFFORTS_TYPE integration"

    case "$INTEGRATE_WAVE_EFFORTS_TYPE" in
        "wave")
            # R009 Sequential Rebuild: Wave integration based on FIRST effort of WAVE
            FIRST_EFFORT=$(echo "✅ State file updated to: $NEXT_STATE"
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
git commit -m "state: SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE → $NEXT_STATE - SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE complete [R288]"
git push
echo "✅ State file committed and pushed"
```

---

### ✅ Step 6: Save TODOs (R287 - SUPREME LAW)
```bash
# Save TODO state before transition (R287 trigger)
save_todos "SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo
git commit -m "todo: orchestrator - SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE complete [R287]"
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

