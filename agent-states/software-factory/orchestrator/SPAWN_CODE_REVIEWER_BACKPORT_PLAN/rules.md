# Orchestrator - SPAWN_CODE_REVIEWER_BACKPORT_PLAN State Rules

# PRIMARY DIRECTIVES

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


# 🔴🔴🔴 MANDATORY: R322 STOP + R405 CONTINUATION FLAG 🔴🔴🔴

**CRITICAL FOR SPAWN STATES - READ THIS FIRST OR FAIL TEST 2!**

## 🚨 THE PATTERN THAT FAILED TEST 2 🚨

**WHAT HAPPENED IN TEST 2:**
- Orchestrator spawned Code Reviewers ✅ (correct)
- Orchestrator stopped per R322 ✅ (correct)
- Orchestrator **DID NOT emit `CONTINUE-SOFTWARE-FACTORY=TRUE`** ❌ (WRONG!)
- Test framework saw no continuation flag → stopped automation
- Test 2 FAILED at iteration 8

**ROOT CAUSE:** Confusion between R322 "stop" and R405 continuation flag

## 🔴 CRITICAL DISTINCTION: TWO INDEPENDENT DECISIONS 🔴

### Decision 1: Should Agent Stop? (R322 - Context Preservation)
**YES - ALWAYS stop after spawning for context preservation**

- **Purpose**: Prevent context overflow between states
- **Action**: `exit 0` to end conversation
- **User Experience**: User sees "/continue-orchestrating" as next step
- **This is NORMAL!** Not an error!

### Decision 2: Should Factory Continue? (R405 - Automation Control)
**YES - ALWAYS emit TRUE for normal spawning operations**

- **Purpose**: Tell automation whether it CAN restart
- **Action**: `echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"` (LAST output before exit)
- **Automation**: Framework will auto-restart orchestrator
- **This is NORMAL!** Designed behavior!

## ✅ REQUIRED PATTERN FOR ALL SPAWN STATES

```bash
# 1. Complete spawning work
echo "✅ Spawned [agent type] for [purpose]"

# 2. Update state file per R324/R288
update_state "[NEXT_STATE]"
commit_state_files_per_r288()

# 3. Save TODOs per R287
save_todos "SPAWNED_[AGENT_TYPE]"

# 4. R322: Stop conversation (context preservation)
echo "🛑 R322: Stopping after spawn for context preservation"

# 5. R405: CONTINUATION FLAG - MUST BE TRUE FOR SPAWNING!
echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"

# 6. Exit to end conversation
exit 0
```

## ❌ WRONG PATTERN (CAUSES TEST FAILURES)

```bash
# ❌ THIS KILLS AUTOMATION - DO NOT DO THIS!
echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=MANUAL_INTERVENTION_REQUIRED"
exit 0

# Result: Test framework stops, Test 2 fails at iteration 8
```

## 🎯 WHY TRUE IS CORRECT FOR SPAWNING

**Spawning is NORMAL operation:**
- ✅ System knows next state (from state machine)
- ✅ Automation can continue (designed workflow)
- ✅ No manual intervention needed
- ✅ Context preservation ≠ error condition

**The orchestrator stopping (`exit 0`) is for:**
- Preserving context between conversation turns
- Allowing state file commits
- Creating clean state boundaries

**The TRUE flag indicates:**
- Automation CAN restart the conversation
- System knows what to do next (check state file)
- Normal operation is proceeding

## 🔴 WHEN TO USE FALSE (NOT FOR SPAWNING!)

**FALSE should ONLY be used for catastrophic failures:**
- ❌ State file corrupted beyond parsing
- ❌ Critical infrastructure destroyed
- ❌ Unrecoverable system errors
- ❌ **NEVER for normal spawning operations!**

## 📋 SPAWN STATE CHECKLIST

**Before exiting this spawn state, verify:**
1. [ ] All agents spawned successfully
2. [ ] State file updated to next state per R324
3. [ ] State files committed per R288
4. [ ] TODOs saved per R287
5. [ ] R322 stop message displayed
6. [ ] **CONTINUE-SOFTWARE-FACTORY=TRUE emitted** ← Critical!
7. [ ] Exited with `exit 0`

**Missing step 6 = Test 2 failure = -100% grade**

---


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

**YOU HAVE ENTERED SPAWN_CODE_REVIEWER_BACKPORT_PLAN STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_SPAWN_CODE_REVIEWER_BACKPORT_PLAN
echo "$(date +%s) - Rules read and acknowledged for SPAWN_CODE_REVIEWER_BACKPORT_PLAN" > .state_rules_read_orchestrator_SPAWN_CODE_REVIEWER_BACKPORT_PLAN
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

## 🎯 STATE OBJECTIVES - SPAWN CODE REVIEWER FOR BACKPORT PLANNING

In the SPAWN_CODE_REVIEWER_BACKPORT_PLAN state, the ORCHESTRATOR is responsible for:

1. **Preparing Fix Documentation for Code Reviewer**
   - Gather all integration failure reports
   - Collect all fixes that were made during integration
   - Identify which effort branches are affected
   - Document what fixes need to be backported where

2. **Creating Code Reviewer Assignment**
   - Prepare clear instructions for Code Reviewer
   - Specify they must create BACKPORT-PLAN.md
   - Include all fix information and affected branches
   - Set clear expectations for plan structure

3. **Spawning Code Reviewer Agent**
   - Spawn Code Reviewer with BACKPORT_PLAN_CREATION state
   - Provide path to fix documentation
   - Monitor Code Reviewer state file creation
   - Track Code Reviewer progress

4. **Stopping After Spawn (R322 Part A)**
   - Update orchestrator-state-v3.json to WAITING_FOR_BACKPORT_PLAN
   - Commit and push state changes
   - STOP and wait for user continuation
   - Do NOT continue to next state automatically

## 📝 REQUIRED ACTIONS

### Step 1: Prepare Fix Documentation
```bash
# Create comprehensive fix documentation for Code Reviewer
cd /efforts/integration-testing

# Create fix manifest for Code Reviewer
cat > FIX-MANIFEST-FOR-BACKPORT.md << 'EOF'
# Fix Manifest for Backport Planning

## Integration Fixes Applied
[Document all fixes made during integration]

## Affected Effort Branches
[List each effort and what fixes it needs]

## Source Files Modified
[List all files that were changed]

## Test Failures Resolved
[Document what tests were failing and how they were fixed]

## Build Issues Fixed
[Document build problems and resolutions]

## Required Backports
Code Reviewer must analyze these fixes and create a detailed
BACKPORT-PLAN.md that specifies:
1. Which fixes go to which effort branches
2. Order of application
3. Verification requirements
4. Success criteria
EOF

echo "✅ Fix documentation prepared for Code Reviewer"
```

### Step 2: Create Code Reviewer Assignment
```bash
# Create clear assignment for Code Reviewer
cat > CODE-REVIEWER-BACKPORT-ASSIGNMENT.md << 'EOF'
# Code Reviewer Assignment: Create Backport Plan

## Your Task
You are being spawned to create a comprehensive BACKPORT-PLAN.md that will guide
SW Engineers in applying integration fixes back to source effort branches.

## Input Documentation
- FIX-MANIFEST-FOR-BACKPORT.md - Contains all fixes made during integration
- Integration test results and failure reports
- Current state of effort branches

## Required Output: BACKPORT-PLAN.md
Create a detailed plan that includes:

### For Each Effort Branch:
1. **Branch Name**: The exact branch to update
2. **Working Directory**: /efforts/[effort-name]
3. **Fixes Required**: Specific fixes from integration
4. **Files to Modify**: Exact file paths and changes
5. **Verification Steps**: How to verify fixes work
6. **Dependencies**: Any order requirements

### Plan Structure:
- Group fixes by effort branch
- Specify exact commands or cherry-picks needed
- Include validation criteria
- Document any risks or special considerations

## Success Criteria
- Every integration fix is mapped to source branches
- Clear instructions for SW Engineers
- No ambiguity about what goes where
- Verification steps for each backport

## Working Directory
/efforts/integration-testing

## State File
Create: /efforts/integration-testing/code-reviewer-state.yaml
Update current_state to BACKPORT_PLAN_COMPLETE when done
EOF

echo "✅ Code Reviewer assignment created"
```

### Step 3: Spawn Code Reviewer Agent
```bash
echo "🚀 Spawning Code Reviewer for backport planning..."

# Log the spawn action
echo "$(date): Spawning Code Reviewer for BACKPORT_PLAN_CREATION" >> SPAWN-LOG.md

# The actual spawn would be done through the Claude interface
cat > /tmp/spawn-code-reviewer-command.md << 'EOF'
@agent-code-reviewer

## BACKPORT PLAN CREATION ASSIGNMENT

You are being spawned to create a comprehensive backport plan for integration fixes.

### Your Immediate Tasks:
1. Read your assignment at: /efforts/integration-testing/CODE-REVIEWER-BACKPORT-ASSIGNMENT.md
2. Read the fix manifest at: /efforts/integration-testing/FIX-MANIFEST-FOR-BACKPORT.md
3. Analyze what fixes need to go to which effort branches
4. Create detailed BACKPORT-PLAN.md with clear instructions for SW Engineers
5. Update your state file when complete

### Working Directory: 
/efforts/integration-testing

### State Transition:
- Initial State: INIT
- Target State: BACKPORT_PLAN_CREATION
- Final State: BACKPORT_PLAN_COMPLETE

### Critical Requirements:
- Map EVERY integration fix to its source branches
- Provide EXACT instructions for SW Engineers
- Include verification steps
- No ambiguity allowed

Start immediately upon spawn.
EOF

echo "✅ Code Reviewer spawn command prepared"
```

### Step 4: Update State and STOP (R322 Part A Enforcement)
```bash
# Update orchestrator state to waiting
cd $CLAUDE_PROJECT_DIR

# Update state file
cat > orchestrator-state-v3.json << 'EOF'
current_state: WAITING_FOR_BACKPORT_PLAN
previous_state: SPAWN_CODE_REVIEWER_BACKPORT_PLAN
agents_spawned:
  - agent: code-reviewer
    purpose: Create backport plan for integration fixes
    state: BACKPORT_PLAN_CREATION
    timestamp: $(date +%s)
backport_status: PLANNING
integration_branch: integration-testing
waiting_for:
  - Code Reviewer to complete BACKPORT-PLAN.md
  - Plan to map all fixes to source branches
EOF

# Commit state change
git add orchestrator-state-v3.json
git commit -m "state: transition to WAITING_FOR_BACKPORT_PLAN after spawning Code Reviewer"
git push

echo "✅ State updated and committed"
echo "🛑 STOPPING per R322 Part A - Must stop after spawn state"
```

## ⚠️ CRITICAL REQUIREMENTS

### Clear Separation of Responsibilities
- **Orchestrator**: ONLY coordinates and spawns
- **Code Reviewer**: Creates the backport plan
- **SW Engineers**: Will implement the fixes (next state)

### No Direct Code Work
- Orchestrator MUST NOT analyze fixes directly
- Orchestrator MUST NOT create the plan itself
- Orchestrator ONLY prepares documentation and spawns

### R322 Part A Compliance
- MUST stop after spawning Code Reviewer
- MUST update state to WAITING_FOR_BACKPORT_PLAN
- MUST NOT continue automatically
- Wait for /continue-orchestrating

## 🚫 FORBIDDEN ACTIONS

1. **Creating the backport plan yourself** - Must be done by Code Reviewer
2. **Analyzing code changes directly** - Code Reviewer's responsibility
3. **Continuing past spawn without stopping** - R322 Part A violation
4. **Spawning SW Engineers in this state** - That's the next state
5. **Making any code edits** - R006 violation

## ✅ PROJECT_DONE CRITERIA

Before transitioning to WAITING_FOR_BACKPORT_PLAN:
- [ ] Fix documentation prepared for Code Reviewer
- [ ] Clear assignment created with expectations
- [ ] Code Reviewer spawned with proper state
- [ ] Orchestrator state updated to WAITING_FOR_BACKPORT_PLAN
- [ ] State changes committed and pushed
- [ ] STOPPED per R322 Part A requirements

## 🔄 STATE TRANSITIONS

### Success Path:
```
SPAWN_CODE_REVIEWER_BACKPORT_PLAN → WAITING_FOR_BACKPORT_PLAN
```
- Code Reviewer spawned successfully
- Waiting for plan completion
- Will continue after user command

### Next States After Waiting:
```
WAITING_FOR_BACKPORT_PLAN → SPAWN_SW_ENGINEER_BACKPORT_FIXES
```
- Once plan is ready
- Spawn SW Engineers to implement

## 🎓 GRADING CRITERIA

You will be evaluated on:
1. **Separation of Concerns** (40%)
   - Code Reviewer does planning
   - Orchestrator only coordinates
   
2. **R322 Part A Compliance** (30%)
   - Proper stop after spawn
   - State update before stop
   
3. **Documentation Quality** (20%)
   - Clear instructions for Code Reviewer
   - Complete fix information provided
   
4. **State Management** (10%)
   - Proper state transitions
   - State file updates

## 💡 TIPS FOR PROJECT_DONE

1. **Let Code Reviewer analyze** - Don't try to understand fixes yourself
2. **Provide complete information** - Give Code Reviewer everything needed
3. **Stop means stop** - R322 Part A is absolute
4. **Clear expectations** - Tell Code Reviewer exactly what output you need

Remember: This state is about DELEGATION, not doing the work yourself!

## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-orchestrating
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**

### 🔴🔴🔴 MANDATORY VALIDATION REQUIREMENT 🔴🔴🔴

**Per R288 and R324**: ALL state file updates MUST be validated before commit:

```bash
# After ANY update to orchestrator-state-v3.json:
"$CLAUDE_PROJECT_DIR/tools/validate-state.sh" orchestrator-state-v3.json || {
    echo "❌ State file validation failed!"
    exit 288
}
```

**Use helper functions for automatic validation:**
```bash
# Source the helper functions
source "$CLAUDE_PROJECT_DIR/utilities/state-file-update-functions.sh"

# Use safe functions that include validation:
safe_state_transition "NEW_STATE" "reason"
safe_update_field "field_name" "value"
```



## 🔴🔴🔴 STATE COMPLETION CHECKLIST 🔴🔴🔴

**Execute these steps IN ORDER to properly complete SPAWN_CODE_REVIEWER_BACKPORT_PLAN:**

### ✅ Step 1: Complete State-Specific Work
**Refer to "Primary Actions" and "Mandatory Validations" sections above for state-specific tasks.**

Once all state work is complete, proceed to mandatory exit protocol:

---

### ✅ Step 2: Set Proposed Next State and Transition Reason
```bash
# Based on state work results, determine proposed next state
PROPOSED_NEXT_STATE="[DETERMINE_FROM_STATE_LOGIC]"
TRANSITION_REASON="SPAWN_CODE_REVIEWER_BACKPORT_PLAN complete - [describe what was accomplished]"
echo "Proposed next state: $PROPOSED_NEXT_STATE"
echo "Transition reason: $TRANSITION_REASON"
```

---

### ✅ Step 3: Update State File (R288 - SUPREME LAW)
```bash
# State Manager validates transition and updates state files (SF 3.0 Pattern)
echo "🔄 Spawning State Manager for SHUTDOWN_CONSULTATION..."

# Prepare work results summary
WORK_RESULTS=$(cat <<EOF
{
  "state_completed": "SPAWN_CODE_REVIEWER_BACKPORT_PLAN",
  "work_accomplished": [
    "[List accomplishments from state work]"
  ],
  "proposed_next_state": "$PROPOSED_NEXT_STATE",
  "transition_reason": "$TRANSITION_REASON"
}
EOF
)

# Spawn State Manager
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "SPAWN_CODE_REVIEWER_BACKPORT_PLAN" \
  --proposed-next-state "$PROPOSED_NEXT_STATE" \
  --transition-reason "$TRANSITION_REASON" \
  --work-results "$WORK_RESULTS"

# State Manager will:
# 1. Validate PROPOSED_NEXT_STATE exists and transition is valid
# 2. Update all 4 state files atomically (R288)
# 3. Commit and push state files
# 4. Return REQUIRED_NEXT_STATE (usually same as proposed unless invalid)

echo "✅ State Manager consultation complete"
echo "✅ State files updated by State Manager"
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

if ! git commit -m "state: SPAWN_CODE_REVIEWER_BACKPORT_PLAN → $NEXT_STATE - SPAWN_CODE_REVIEWER_BACKPORT_PLAN complete [R288]"; then
    echo "❌ CRITICAL: Git commit failed - likely schema validation error"
    echo "State: SPAWN_CODE_REVIEWER_BACKPORT_PLAN"
    echo "Attempted transition from: SPAWN_CODE_REVIEWER_BACKPORT_PLAN"
    echo ""
    echo "Common causes:"
    echo "  - Schema validation failure (check pre-commit hook output above)"
    echo "  - Missing required fields in JSON files"
    echo "  - Invalid JSON syntax"
    echo ""
    echo "🛑 Cannot proceed - manual intervention required"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=SCHEMA_VALIDATION"
    exit 1
fi

git push || echo "⚠️ WARNING: Push failed - committed locally"
echo "✅ State file committed and pushed"
git push
echo "✅ State file committed and pushed"
```

---

### ✅ Step 6: Save TODOs (R287 - SUPREME LAW)
```bash
# Save TODO state before transition (R287 trigger)
save_todos "SPAWN_CODE_REVIEWER_BACKPORT_PLAN_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - SPAWN_CODE_REVIEWER_BACKPORT_PLAN complete [R287]"; then
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

### ✅ Step 7: Output Continuation Flag (R405 - SUPREME LAW) ⚠️ MANDATORY
```bash
# Output continuation flag as LAST action (R405)
# Use TRUE for normal completion, FALSE only for catastrophic errors

echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"
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
PROPOSED_NEXT_STATE="WAITING_FOR_BACKPORT_PLAN"
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
### 🚨 SPAWN STATE PATTERN - R322 + R405 USAGE 🚨

**Spawning operations require R322 stop for context preservation:**
```bash
# After spawning agent(s)
echo "✅ Spawned agents for work"

# R322 checkpoint (context preservation)
echo "🛑 R322: Stopping after spawn for context preservation"

# Flag? → MUST BE TRUE (normal operation!)
echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"

# Stop inference
exit 0
```

**Why TRUE is correct:**
- Spawning is NORMAL operation
- System knows next state
- Automation can continue
- **Context preservation ≠ manual intervention needed!**
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

