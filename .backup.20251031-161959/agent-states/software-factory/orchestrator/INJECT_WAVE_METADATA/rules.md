# Orchestrator - INJECT_WAVE_METADATA State Rules


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

**YOU HAVE ENTERED INJECT_WAVE_METADATA STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_INJECT_WAVE_METADATA
echo "$(date +%s) - Rules read and acknowledged for INJECT_WAVE_METADATA" > .state_rules_read_orchestrator_INJECT_WAVE_METADATA
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY INJECT_WAVE_METADATA WORK UNTIL RULES ARE READ:
- ❌ Start inject wave metadata
- ❌ Start update tracking files
- ❌ Start configure wave settings
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
   ❌ WRONG: "I acknowledge all INJECT_WAVE_METADATA rules"
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

### ✅ CORRECT PATTERN FOR INJECT_WAVE_METADATA:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/[first-rule-file].md
2. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/[second-rule-file].md  
4. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. "Ready to execute INJECT_WAVE_METADATA work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY INJECT_WAVE_METADATA work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ You have stated readiness to execute INJECT_WAVE_METADATA work
1. **READ** every rule file listed in PRIMARY DIRECTIVES below
2. **ACKNOWLEDGE** each rule individually with number and description
3. **ONLY THEN** proceed with INJECT_WAVE_METADATA work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY INJECT_WAVE_METADATA work before reading and acknowledging rules:**
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

## 📋 PRIMARY DIRECTIVES FOR INJECT_WAVE_METADATA STATE

### 🚨🚨🚨 R213 - Wave and Effort Metadata Injection
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R213-wave-and-effort-metadata-protocol.md`
**Criticality**: BLOCKING - Must inject metadata before spawning
**Summary**: Inject parallelization metadata into wave implementation plans

### 🔴🔴🔴 R234 - Mandatory State Traversal (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R234-mandatory-state-traversal-supreme-law.md`
**Criticality**: SUPREME LAW - Violation = -100% automatic failure
**Summary**: Must traverse all states in sequence, no skipping allowed

### 🔴🔴🔴 R288 - State File Update and Commit Protocol (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
**Criticality**: SUPREME LAW - Update on every transition
**Summary**: Update orchestrator-state-v3.json on all state changes

### 🚨🚨🚨 R288 - State File Update and Commit Protocol
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
**Criticality**: BLOCKING - Push within 60 seconds
**Summary**: Commit and push state file immediately after updates

### 🔴🔴🔴 R232 - TodoWrite Pending Items Override (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R232-todowrite-pending-items-override.md`
**Criticality**: SUPREME LAW - Pending items are COMMANDS
**Summary**: Any pending TODO items must be executed immediately

## 🚨 INJECT_WAVE_METADATA IS A VERB - START INJECTING METADATA IMMEDIATELY! 🚨

### IMMEDIATE ACTIONS UPON ENTERING INJECT_WAVE_METADATA

**THE MOMENT YOU ENTER THIS STATE, YOU MUST:**
1. Open wave implementation plan for editing NOW
2. Insert R213 parallelization metadata immediately
3. Check TodoWrite for pending items and process them
4. Save and validate metadata without delay

**FORBIDDEN - AUTOMATIC FAILURE:**
- ❌ "STATE TRANSITION COMPLETE: Now in INJECT_WAVE_METADATA" [stops]
- ❌ "Successfully entered INJECT_WAVE_METADATA state" [waits]
- ❌ "Ready to start injecting metadata" [pauses]
- ❌ "I'm in INJECT_WAVE_METADATA state" [does nothing]

**REQUIRED - IMMEDIATE ACTION:**
- ✅ "Entering INJECT_WAVE_METADATA, opening wave implementation plan for editing NOW..."
- ✅ "START INJECTING METADATA per R213, inserting parallelization metadata immediately..."
- ✅ "INJECT_WAVE_METADATA: Saving and validating metadata without delay..."

## 🔴🔴🔴 MANDATORY MEASUREMENT RULE - R304 🔴🔴🔴

**ABSOLUTE REQUIREMENTS:**
- ✅ MUST use `$CLAUDE_PROJECT_DIR/tools/line-counter.sh` for ALL line counting
- ❌ NEVER use `wc -l` or any manual counting method
- ❌ NEVER count lines any other way - this is a -100% automatic failure
- ✅ MUST specify both -b (base branch) and -c (current branch) parameters
- ✅ Base branch MUST be phase integration branch (NOT "main")

**Failure to use the line counter tool = AUTOMATIC -100% GRADE**

## State Context
This state is responsible for injecting R213 parallelization metadata into wave implementation plans before any agent spawning occurs.

## R213 Metadata Injection Protocol

**CRITICAL**: This metadata MUST be injected BEFORE spawning Code Reviewers or SW Engineers!

**🚨 MANDATORY UUID SUFFIX REQUIREMENT 🚨**

ALL effort branch names MUST include the `project_prefix` suffix from `orchestrator-state-v3.json` to ensure test isolation and prevent branch collisions in concurrent environments.

**Extraction of project_prefix:**
```bash
# Extract project_prefix (includes UUID suffix) from state file
PROJECT_PREFIX=$(jq -r '.pre_planned_infrastructure.project_prefix // .current_project.project_id' orchestrator-state-v3.json)
echo "Project prefix with UUID: $PROJECT_PREFIX"
# Example: fastapi-hello-sf3-test-6a37fb6a-cf44-40fb-b02b-11a2ca5591af
```

**Branch Name Construction:**
```bash
# WRONG - Missing UUID suffix
branch: "phase3/wave1/sync-engine-foundation"
branch: "effort/phase1-wave1-fastapi-app-structure"

# CORRECT - Includes project_prefix UUID suffix
branch: "phase3/wave1/sync-engine-foundation-${PROJECT_PREFIX}"
branch: "effort/phase1-wave1-fastapi-app-structure-${PROJECT_PREFIX}"
```

```bash
# Example metadata to inject into wave plan:
EFFORT_METADATA:
  effort_id: "E3.1.1"
  name: "sync-engine-foundation"
  can_parallelize: false
  blocks: ["E3.1.2", "E3.1.3", "E3.1.4", "E3.1.5"]
  dependencies: []
  estimated_lines: 600
  assigned_to: "sw-engineer-1"
  working_directory: "/efforts/phase3/wave1/sync-engine-foundation"
  branch: "phase3/wave1/sync-engine-foundation-${PROJECT_PREFIX}"
  # ☝️ MANDATORY: Branch name MUST include ${PROJECT_PREFIX} for isolation
```

## Critical Requirements

1. **READ** the wave implementation plan with Read tool
2. **IDENTIFY** all efforts in the wave
3. **INJECT** parallelization metadata for EACH effort
4. **SAVE** the updated plan
5. **VERIFY** metadata is present before proceeding
6. **TRANSITION** to next state in mandatory sequence

## State Transitions

- **FROM**: Previous state in mandatory sequence
- **TO**: Next state per R234 mandatory traversal
- **CANNOT SKIP**: This state is part of mandatory sequence

## Validation Before Transition

**🚨 MANDATORY VALIDATION: UUID Suffix Compliance 🚨**

Before transitioning from INJECT_WAVE_METADATA, you MUST validate that all effort branch names include the project_prefix UUID suffix. This prevents isolation violations that waste API costs.

```bash
validate_metadata_injection() {
    echo "🔍 Validating R213 metadata injection..."

    # Extract project_prefix for validation
    PROJECT_PREFIX=$(jq -r '.pre_planned_infrastructure.project_prefix // .current_project.project_id' orchestrator-state-v3.json)

    if [ -z "$PROJECT_PREFIX" ]; then
        echo "❌ FATAL: Cannot validate - project_prefix not found in state file!"
        exit 213
    fi

    echo "Validating against project_prefix: $PROJECT_PREFIX"

    # Check each effort has metadata
    for effort in efforts/phase${PHASE}/wave${WAVE}/*/; do
        IMPL_PLAN=$(ls -t "${effort}.software-factory/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}/IMPLEMENTATION-PLAN--"*.md 2>/dev/null | head -1)

        # Validate parallelization metadata exists
        if ! grep -q "can_parallelize:" "$IMPL_PLAN"; then
            echo "❌ FATAL: Missing parallelization metadata in ${effort}"
            exit 213
        fi

        # CRITICAL: Validate branch name includes project_prefix UUID
        BRANCH_NAME=$(grep -E "^\s*branch:" "$IMPL_PLAN" | head -1 | cut -d: -f2- | xargs)
        if [ -n "$BRANCH_NAME" ]; then
            if ! echo "$BRANCH_NAME" | grep -q "$PROJECT_PREFIX"; then
                echo "❌ FATAL: Effort branch missing UUID suffix!"
                echo "  Effort: ${effort}"
                echo "  Branch: $BRANCH_NAME"
                echo "  Expected to contain: $PROJECT_PREFIX"
                echo "  This will cause isolation validation failures!"
                exit 213
            fi
            echo "✅ Branch has UUID suffix: $BRANCH_NAME"
        fi
    done

    echo "✅ All efforts have R213 metadata with proper UUID suffixes"
}
```

## Next Steps

After successfully injecting metadata:
1. Update orchestrator-state-v3.json (R288)
2. Commit and push changes (R288)
3. Check TodoWrite for pending items (R232)
4. Transition to next mandatory state (R234)

## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-orchestrating
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**



## 🔴🔴🔴 STATE COMPLETION CHECKLIST 🔴🔴🔴

**Execute these steps IN ORDER to properly complete INJECT_WAVE_METADATA:**

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
- Missing Step 2: No next state = stuck forever
- Missing Step 3: No State Manager spawn = state machine broken (R288 violation, -100%)
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

