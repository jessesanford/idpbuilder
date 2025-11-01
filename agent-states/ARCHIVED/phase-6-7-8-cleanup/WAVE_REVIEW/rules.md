# Orchestrator - REVIEW_WAVE_ARCHITECTURE State Rules

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

## 🟢🟢🟢 REVIEW_WAVE_ARCHITECTURE - NORMAL CONTINUATION AFTER ARCHITECT REVIEW 🟢🟢🟢

**Per R322: Wave reviews are NORMAL OPERATIONS that continue automatically!**

### IMPORTANT CLARIFICATION:
- **MONITORING_INTEGRATE_WAVE_EFFORTS → REVIEW_WAVE_ARCHITECTURE** requires a stop (per R322 Part D)
- **REVIEW_WAVE_ARCHITECTURE → next state** also requires stop with TRUE flag for normal operations

### MANDATORY STOP PROTOCOL:
1. ✅ Process architect review results
2. ✅ Determine appropriate next state based on review
3. ✅ Update orchestrator-state-v3.json with new state
4. ✅ Commit and push the state file
5. ✅ Output CONTINUE-SOFTWARE-FACTORY=TRUE (for normal operations)
6. ✅ **STOP PROCESSING** (mandatory for context preservation)

### NORMAL TRANSITIONS (ALL CONTINUE):
- ✅ REVIEW_WAVE_ARCHITECTURE → WAVE_START (approved, continue to next wave)
- ✅ REVIEW_WAVE_ARCHITECTURE → CASCADE_REINTEGRATION (fixes needed, but automated)
- ✅ REVIEW_WAVE_ARCHITECTURE → COMPLETE_PHASE (last wave of phase reviewed)

### EXCEPTIONAL CASES (RARE - CONTINUE-SOFTWARE-FACTORY=FALSE):
- ❌ Architect marked as OFF_TRACK (major architectural issues)
- ❌ CRITICAL blocking issues that require manual intervention
- ❌ State corruption or invalid transitions

### AUTOMATION CONTINUITY:
```bash
# After processing architect review:
if [ "$REVIEW_STATUS" == "APPROVED" ]; then
    echo "✅ Architect review passed, continuing to next wave"
    NEXT_STATE="WAVE_START"
elif [ "$REVIEW_STATUS" == "NEEDS_FIXES" ]; then
    echo "⚠️ Fixes needed, but continuing automatically"
    NEXT_STATE="CASCADE_REINTEGRATION"
elif [ "$REVIEW_STATUS" == "OFF_TRACK" ]; then
    echo "❌ OFF_TRACK - manual intervention required"
    NEXT_STATE="ERROR_RECOVERY"
    CONTINUE_FLAG="FALSE"  # Only stop for OFF_TRACK
else
    CONTINUE_FLAG="TRUE"  # Default: keep going
fi

echo "CONTINUE-SOFTWARE-FACTORY=$CONTINUE_FLAG"
```

**WAVE REVIEWS ARE NORMAL - KEEP THE FACTORY RUNNING!**

---


## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED REVIEW_WAVE_ARCHITECTURE STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_REVIEW_WAVE_ARCHITECTURE
echo "$(date +%s) - Rules read and acknowledged for REVIEW_WAVE_ARCHITECTURE" > .state_rules_read_orchestrator_REVIEW_WAVE_ARCHITECTURE
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY REVIEW_WAVE_ARCHITECTURE WORK UNTIL RULES ARE READ:
- ❌ Start spawn architecture reviewer
- ❌ Start request wave assessment
- ❌ Start collect review feedback
- ❌ Update state files
- ❌ Continue to next state
- ❌ Think about what to do in this state

### ✅ YOU MUST IMMEDIATELY:

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**
### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

## 📋 PRIMARY DIRECTIVES FOR REVIEW_WAVE_ARCHITECTURE STATE

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

5. **🔴🔴🔴 R009** - Wave Integration Protocol
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R009-integration-branch-creation.md`
   - Criticality: SUPREME LAW - Wave integration requirements
   - Summary: Wave management and integration protocols


### ❌ ANTI-PATTERNS THAT WILL CAUSE FAILURE:

1. **Fake Acknowledgment Without Reading**:
   ```
   ❌ WRONG: "I acknowledge R151, R208, R053..."
   (No Read tool calls detected = AUTOMATIC FAILURE)
   ```

2. **Bulk Acknowledgment**:
   ```
   ❌ WRONG: "I acknowledge all REVIEW_WAVE_ARCHITECTURE rules"
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

### ✅ CORRECT PATTERN FOR REVIEW_WAVE_ARCHITECTURE:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/[first-rule-file].md
2. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/[second-rule-file].md  
4. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. "Ready to execute REVIEW_WAVE_ARCHITECTURE work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY REVIEW_WAVE_ARCHITECTURE work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ You have stated readiness to execute REVIEW_WAVE_ARCHITECTURE work
1. **READ** every rule file listed in PRIMARY DIRECTIVES below
2. **ACKNOWLEDGE** each rule individually with number and description
3. **ONLY THEN** proceed with REVIEW_WAVE_ARCHITECTURE work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY REVIEW_WAVE_ARCHITECTURE work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

**THE SYSTEM IS MONITORING_SWE_PROGRESS YOUR READ TOOL CALLS!**

## 🚨 REVIEW_WAVE_ARCHITECTURE IS A VERB - START REVIEW PROCESS IMMEDIATELY! 🚨

### IMMEDIATE ACTIONS UPON ENTERING REVIEW_WAVE_ARCHITECTURE

**THE MOMENT YOU ENTER THIS STATE, YOU MUST:**
1. Spawn Architect for wave review NOW
2. Provide integration branch for review
3. Check TodoWrite for pending items and process them
4. Include all effort completion status

**FORBIDDEN - AUTOMATIC FAILURE:**
- ❌ "STATE TRANSITION COMPLETE: Now in REVIEW_WAVE_ARCHITECTURE" [stops]
- ❌ "Successfully entered REVIEW_WAVE_ARCHITECTURE state" [waits]
- ❌ "Ready to start review process" [pauses]
- ❌ "I'm in REVIEW_WAVE_ARCHITECTURE state" [does nothing]

**REQUIRED - IMMEDIATE ACTION:**
- ✅ "Entering REVIEW_WAVE_ARCHITECTURE, Spawn Architect for wave review NOW..."
- ✅ "START REVIEW PROCESS, provide integration branch for review..."
- ✅ "REVIEW_WAVE_ARCHITECTURE: Include all effort completion status..."

## 🔴🔴🔴 MANDATORY MEASUREMENT RULE - R304 🔴🔴🔴

**ABSOLUTE REQUIREMENTS:**
- ✅ MUST use `$CLAUDE_PROJECT_DIR/tools/line-counter.sh` for ALL line counting
- ❌ NEVER use `wc -l` or any manual counting method
- ❌ NEVER count lines any other way - this is a -100% automatic failure
- ✅ MUST specify both -b (base branch) and -c (current branch) parameters
- ✅ Base branch MUST be phase integration branch (NOT "main")

**Failure to use the line counter tool = AUTOMATIC -100% GRADE**

## State Context
You are requesting an architect review of a completed wave's integration. The architect will assess the integrated code for architectural compliance, patterns, and readiness to proceed.

## 🔴🔴🔴 CRITICAL: MANDATORY STATE FILE UPDATE (R288) 🔴🔴🔴

### IMMEDIATELY upon entering REVIEW_WAVE_ARCHITECTURE state:

```bash
# 1. Update state machine
update_orchestrator_state "REVIEW_WAVE_ARCHITECTURE" "Requesting architect review of wave integration"

# 2. Record review request in state file
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
git commit -m "state: REVIEW_WAVE_ARCHITECTURE → $NEXT_STATE - REVIEW_WAVE_ARCHITECTURE complete [R288]"
git push
echo "✅ State file committed and pushed"
```

---

### ✅ Step 6: Save TODOs (R287 - SUPREME LAW)
```bash
# Save TODO state before transition (R287 trigger)
save_todos "REVIEW_WAVE_ARCHITECTURE_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo
git commit -m "todo: orchestrator - REVIEW_WAVE_ARCHITECTURE complete [R287]"
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

