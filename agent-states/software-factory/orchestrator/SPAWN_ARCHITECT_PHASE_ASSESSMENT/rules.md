# Orchestrator - SPAWN_ARCHITECT_PHASE_ASSESSMENT State Rules

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
## 🔴🔴🔴 CRITICAL: STATE FILES PROPOSE, STATE MANAGER DECIDES 🔴🔴🔴

**R288 SUPREME LAW - YOU NEVER CALL update_state DIRECTLY!**

This state file uses the **BOOKEND PATTERN** (R600):
1. **START**: Set `PROPOSED_NEXT_STATE` and `TRANSITION_REASON` variables
2. **WORK**: Execute state-specific logic
3. **END**: Follow 10-step completion checklist to exit properly

**NEVER CALL update_state() DIRECTLY - IT IS PROHIBITED!**
- ❌ `update_state "NEXT_STATE" "reason"` = SYSTEM VIOLATION
- ✅ `PROPOSED_NEXT_STATE="WAITING_FOR_PHASE_ASSESSMENT"` = CORRECT
- ✅ `TRANSITION_REASON="reason"` = CORRECT

The State Manager (`run-software-factory.sh`) handles ALL state transitions.

**See:**
- `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md` (SUPREME LAW)
- `$CLAUDE_PROJECT_DIR/rule-library/R600-checklist-execution-protocol.md`

---

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


## 🛑🛑🛑 R322 CONTEXT PRESERVATION CHECKPOINT 🛑🛑🛑

**CRITICAL: STOP INFERENCE BUT SET TRUE FLAG FOR NORMAL OPS**

### AFTER SPAWNING ARCHITECT YOU MUST:
1. ✅ Complete spawn operation
2. ✅ Update orchestrator-state-v3.json
3. ✅ Save TODOs per R287
4. ✅ Commit and push changes
5. ✅ Output CONTINUE-SOFTWARE-FACTORY=TRUE
6. ✅ Exit with code 0

### KEY UNDERSTANDING:
- **Stop inference**: YES (preserve context)
- **Set TRUE flag**: YES (normal operation)
- **Automation restarts**: Automatically

### CORRECT PATTERN:
```bash
echo "✅ Spawned architect for phase assessment"
echo "📊 Architect will evaluate phase completion"

# Propose state transition (R288)
PROPOSED_NEXT_STATE="WAITING_FOR_PHASE_ASSESSMENT"
TRANSITION_REASON="Architect spawned for phase assessment"
save_todos "SPAWNED_ARCHITECT"

# Commit changes
git add -A && git commit -m "state: spawned architect"
git push

# Follow completion checklist for proper R600 bookend exit
```

**See STATE COMPLETION CHECKLIST below for complete exit protocol.**

**This is NORMAL OPERATION - not an error!**

---

## 🔴🔴🔴 R322 MANDATORY: STOP INFERENCE AFTER SPAWNING 🔴🔴🔴

**CRITICAL DISTINCTION: STOP INFERENCE ≠ SET FALSE FLAG**

After spawning the architect in this state, you MUST:
1. Record what was spawned in state file
2. Save TODOs per R287
3. Commit and push state changes
4. Output CONTINUE-SOFTWARE-FACTORY=TRUE (normal operation!)
5. EXIT immediately with code 0 (preserve context)

**KEY UNDERSTANDING:**
- Stop inference: YES (required for context preservation)
- Set flag FALSE: NO! (spawning is normal, set TRUE)
- External automation will auto-restart when it sees TRUE

**VIOLATION = AUTOMATIC -100% FAILURE**

See: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-before-state-transitions.md`

---

# Orchestrator - SPAWN_ARCHITECT_PHASE_ASSESSMENT State Rules

## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED SPAWN_ARCHITECT_PHASE_ASSESSMENT STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
mkdir -p markers/state-verification && touch "markers/state-verification/state_rules_read_orchestrator_SPAWN_ARCHITECT_PHASE_ASSESSMENT-$(date +%Y%m%d-%H%M%S)"
echo "$(date +%s) - Rules read and acknowledged for SPAWN_ARCHITECT_PHASE_ASSESSMENT" > "markers/state-verification/state_rules_read_orchestrator_SPAWN_ARCHITECT_PHASE_ASSESSMENT-$(date +%Y%m%d-%H%M%S)"
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY SPAWN_ARCHITECT_PHASE_ASSESSMENT WORK UNTIL RULES ARE READ:
- ❌ Start spawn architect agent
- ❌ Start request phase assessment
- ❌ Start evaluate phase completion
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
   ❌ WRONG: "I acknowledge all SPAWN_ARCHITECT_PHASE_ASSESSMENT rules"
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

### ✅ CORRECT PATTERN FOR SPAWN_ARCHITECT_PHASE_ASSESSMENT:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/[first-rule-file].md
2. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/[second-rule-file].md  
4. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. "Ready to execute SPAWN_ARCHITECT_PHASE_ASSESSMENT work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY SPAWN_ARCHITECT_PHASE_ASSESSMENT work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ You have stated readiness to execute SPAWN_ARCHITECT_PHASE_ASSESSMENT work
1. **READ** every rule file listed in PRIMARY DIRECTIVES below
2. **ACKNOWLEDGE** each rule individually with number and description
3. **ONLY THEN** proceed with SPAWN_ARCHITECT_PHASE_ASSESSMENT work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY SPAWN_ARCHITECT_PHASE_ASSESSMENT work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

**THE SYSTEM IS MONITORING_SWE_PROGRESS YOUR READ TOOL CALLS!**

## 📋 PRIMARY DIRECTIVES FOR SPAWN_ARCHITECT_PHASE_ASSESSMENT STATE

### 🔴🔴🔴 R301 - Integration Branch Current Tracking (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R301-integration-branch-current-tracking.md`
**Criticality**: SUPREME LAW - Only ONE current integration allowed
**Summary**: MUST use current_integration.branch, NEVER deprecated branches

### 🚨🚨🚨 R257 - Mandatory Phase Assessment Report
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R257-mandatory-phase-assessment-report.md`
**Criticality**: BLOCKING - Phase cannot complete without report
**Summary**: Architect MUST create assessment report file

### 🚨🚨🚨 R285 - Mandatory Phase Integration Before Assessment
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R285-mandatory-phase-integration-before-assessment.md`
**Criticality**: BLOCKING - Integration must precede assessment
**Summary**: Assess integrated work, not individual efforts

### 🔴🔴🔴 R233 - All States Require Immediate Action
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R233-all-states-immediate-action.md`
**Criticality**: CRITICAL - States are verbs
**Summary**: SPAWN_ARCHITECT_PHASE_ASSESSMENT means SPAWN NOW

### 🔴🔴🔴 R510 - STATE EXECUTION CHECKLIST COMPLIANCE (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R510-state-execution-checklist-compliance.md`
**Criticality**: SUPREME LAW - Checklist compliance required
**Summary**: MUST complete and acknowledge every checklist item before transition

---

## 🔴🔴🔴 MANDATORY EXECUTION CHECKLIST

**RULE**: R510 requires every state to have explicit execution checklist
**ENFORCEMENT**: All BLOCKING items must complete before transition
**ACKNOWLEDGMENT**: Must output "✅ CHECKLIST[n]: [description] [proof]" for each item

### BLOCKING REQUIREMENTS (Cannot proceed without)

- [ ] 1. Spawn Architect for phase-level assessment
  - Agent: `agent-architect`
  - State: `PHASE_ASSESSMENT`
  - Phase: Current phase number
  - Focus: Complete phase validation across all waves
  - Data: All wave integration branches for the phase
  - Metrics: Phase completion metrics and quality indicators
  - Validation: Agent ID returned with timestamp
  - **BLOCKING**: Must spawn before any transition

### STANDARD EXECUTION TASKS (Required)

- [ ] 2. Update state file with architect spawn metadata
  - Field: `architect_review_states.phase_assessments` append new entry
  - Include: architect_id, phase, spawned_at, wave_integrations
  - Expected: State file contains architect spawn record

- [ ] 3. Prepare phase assessment context for architect
  - Collect: All wave integration branch names
  - Collect: Phase completion metrics
  - Collect: Quality indicators from wave reviews
  - Expected: Architect has complete phase context

### EXIT REQUIREMENTS (Must complete before transition)

- [ ] 4. Update state file to WAITING_FOR_PHASE_ASSESSMENT per R288
  - Field: `current_state`
  - Value: `"WAITING_FOR_PHASE_ASSESSMENT"`
  - Also update: `previous_state`, `transition_time`, `transition_reason`
jq '.phase_assessment.requested_at = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
jq '.phase_assessment.expected_report = \"phase-assessments/phase${PHASE}/PHASE-${PHASE}-ASSESSMENT-REPORT.md\"' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
jq '.phase_assessment.type = \"COMPLETE_PHASE\"' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
jq '.phase_assessment.phase_branch = \"$PHASE_BRANCH\"' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
jq '.phase_assessment.wave_count = $WAVE_COUNT' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
jq '.phase_assessment.status = \"PENDING\"' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json

# Immediate transition to waiting
transition_to "WAITING_FOR_PHASE_ASSESSMENT"
```

## Success Criteria for Phase Assessment

The phase assessment passes when:
- [ ] All phase features implemented
- [ ] Architectural consistency verified
- [ ] APIs stable and documented
- [ ] Test coverage meets requirements
- [ ] Performance benchmarks met
- [ ] Security requirements satisfied
- [ ] Documentation complete
- [ ] No blocking issues remain

## State Transitions

From SPAWN_ARCHITECT_PHASE_ASSESSMENT:
- **Always** → WAITING_FOR_PHASE_ASSESSMENT (immediate after spawn)

## Multi-Phase Considerations

If this is NOT the final phase:
- Architect may recommend starting next phase planning
- COMPLETE_PHASE state will handle transition to next phase
- Document lessons learned for next phase

## Required Actions

1. Update state file with assessment request
2. Spawn architect with comprehensive context
3. Provide all integration branches
4. Include phase metrics and achievements
5. Transition to WAITING_FOR_PHASE_ASSESSMENT

## Grading Impact

- Spawning architect promptly: +10 points
- Providing complete context: +10 points
- Including all branches: +10 points
- Proper state file updates: +10 points
- Missing phase assessment before PROJECT_DONE: -100 points (CRITICAL FAILURE)

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

**Execute these steps IN ORDER to properly complete SPAWN_ARCHITECT_PHASE_ASSESSMENT:**

### ✅ Step 1: Complete State-Specific Work
**Refer to "Primary Actions" and "Mandatory Validations" sections above for state-specific tasks.**

**CRITICAL**: During state work, set these variables (DO NOT call update_state):
```bash
PROPOSED_NEXT_STATE="[state determined by logic]"
TRANSITION_REASON="[reason for transition]"
```

Once all state work is complete, proceed to mandatory exit protocol:

---

### ✅ Step 2: Verify Proposed Transition Set
```bash
# Verify PROPOSED_NEXT_STATE was set during state work
if [ -z "$PROPOSED_NEXT_STATE" ]; then
    echo "❌ CRITICAL: PROPOSED_NEXT_STATE not set!"
    echo "State work must set PROPOSED_NEXT_STATE and TRANSITION_REASON"
    exit 288
fi

echo "✅ Proposed transition: SPAWN_ARCHITECT_PHASE_ASSESSMENT → $PROPOSED_NEXT_STATE"
echo "   Reason: $TRANSITION_REASON"
```

---

### ✅ Step 3: Save Proposed Transition to File (R288)
```bash
# Write proposed transition to file for State Manager
cat > "$CLAUDE_PROJECT_DIR/.proposed-transition" <<EOF
PROPOSED_NEXT_STATE="$PROPOSED_NEXT_STATE"
TRANSITION_REASON="$TRANSITION_REASON"
CURRENT_STATE="SPAWN_ARCHITECT_PHASE_ASSESSMENT"
TIMESTAMP="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
EOF

echo "✅ Proposed transition saved to .proposed-transition"
```

---

### ✅ Step 4: Validate Current State File (R324)
```bash
# Validate state file integrity BEFORE proposing transition
"$CLAUDE_PROJECT_DIR/tools/validate-state.sh" orchestrator-state-v3.json || {
    echo "❌ State file validation failed!"
    exit 288
}
echo "✅ Current state file validated"
```

---

### ✅ Step 5: Commit Work Products (if any)
```bash
# Commit any work products created during this state
# (State file itself is committed by State Manager)
if [ -n "$(git status --porcelain | grep -v orchestrator-state-v3.json)" ]; then
    git add .
    git commit -m "work: SPAWN_ARCHITECT_PHASE_ASSESSMENT work products [R288]"
    git push
    echo "✅ Work products committed and pushed"
else
    echo "✅ No work products to commit"
fi
```

---

### ✅ Step 4: Save TODOs (R287 - SUPREME LAW)
```bash
# Save TODO state before transition (R287 trigger)
save_todos "SPAWN_ARCHITECT_PHASE_ASSESSMENT_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - SPAWN_ARCHITECT_PHASE_ASSESSMENT complete [R287]"; then
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

### ✅ Step 7: Output Transition Proposal (R288)
```bash
# Output the proposed transition for State Manager
echo "════════════════════════════════════════════════════════"
echo "PROPOSED STATE TRANSITION (R288):"
echo "  FROM: SPAWN_ARCHITECT_PHASE_ASSESSMENT"
echo "  TO:   $PROPOSED_NEXT_STATE"
echo "  REASON: $TRANSITION_REASON"
echo "════════════════════════════════════════════════════════"
```

---

### ✅ Step 8: Output Continuation Flag (R405 - SUPREME LAW)
```bash
# Output continuation flag (R405)
# TRUE = state work complete, transition proposed
# FALSE = error occurred, manual intervention needed

if [ "$PROPOSED_NEXT_STATE" = "ERROR_RECOVERY" ]; then
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=MANUAL_INTERVENTION_REQUIRED"
else
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"
fi
```

**⚠️ THIS MUST BE THE LAST OUTPUT BEFORE STOP! ⚠️**

---

### ✅ Step 9: State Manager Transition Checkpoint
```
🔄 STATE MANAGER TAKES CONTROL HERE (R288)

The State Manager (run-software-factory.sh) will:
1. Read .proposed-transition file
2. Validate proposed transition against state machine
3. Update orchestrator-state-v3.json with new state
4. Commit and push state file
5. Re-invoke orchestrator in new state (if CONTINUE-SOFTWARE-FACTORY=TRUE)

DO NOT PROCEED PAST THIS POINT - STATE MANAGER HANDLES TRANSITION
```

---

### ✅ Step 10: Stop Processing (R322 - SUPREME LAW)
```bash
# Stop for State Manager handoff (R322)
echo "🛑 State work complete - State Manager will handle transition"
exit 0
```

---

## 🚨 CHECKLIST ENFORCEMENT 🚨

**Skipping ANY step in this checklist = FAILURE:**
- Missing Step 1: State work incomplete
- Missing Step 2: No proposed transition = stuck forever (R288 violation, -100%)
- Missing Step 3: State Manager can't read proposal = broken automation (R288 violation, -100%)
- Missing Step 4: Invalid state = corruption (R324 violation)
- Missing Step 5: Work products lost (data loss)
- Missing Step 6: TODOs lost on compaction (R287 violation, -20% to -100%)
- Missing Step 7: No transition visibility (R288 violation)
- Missing Step 8: Automation doesn't know if it can continue (R405 violation, -100%)
- Missing Step 9: State Manager handoff failed (R288 violation, -100%)
- Missing Step 10: Context corruption (R322 violation, -100%)

- Missing Step 2: No proposed next state = State Manager can't proceed
- Missing Step 3: No State Manager consultation = bypassing bookend pattern (-100%)
- Missing Step 4: No TODO save = work lost (R287 violation, -20% to -100%)
- Missing Step 5: No continuation flag = automation stops (R405 violation, -100%)
- Missing Step 6: No exit = R322 violation (-100%)

**ALL 6 STEPS ARE MANDATORY - NO EXCEPTIONS**
**NOTE**: State file validation and commits now handled by State Manager (SF 3.0 pattern) - NO EXCEPTIONS**

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
PROPOSED_NEXT_STATE="WAITING_FOR_PHASE_ASSESSMENT"
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

