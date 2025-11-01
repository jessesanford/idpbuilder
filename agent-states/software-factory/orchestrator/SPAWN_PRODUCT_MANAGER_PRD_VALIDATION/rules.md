# Orchestrator - SPAWN_PRODUCT_MANAGER_PRD_VALIDATION State Rules

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
## 🔴🔴🔴 MANDATORY: R322 STOP + R405 CONTINUATION FLAG 🔴🔴🔴

**CRITICAL FOR SPAWN STATES - READ THIS FIRST!**

### 🔴 CRITICAL DISTINCTION: TWO INDEPENDENT DECISIONS 🔴

#### Decision 1: Should Agent Stop? (R322 - Context Preservation)
**YES - ALWAYS stop after spawning for context preservation**

- **Purpose**: Prevent context overflow between states
- **Action**: `exit 0` to end conversation
- **User Experience**: User sees "/continue-orchestrating" as next step
- **This is NORMAL!** Not an error!

#### Decision 2: Should Factory Continue? (R405 - Automation Control)
**YES - ALWAYS emit TRUE for normal spawning operations**

- **Purpose**: Tell automation whether it CAN restart
- **Action**: `echo "CONTINUE-SOFTWARE-FACTORY=TRUE CHECKPOINT=R322"` (LAST output before exit)
- **Automation**: Framework will auto-restart orchestrator
- **This is NORMAL!** Designed behavior!

### ✅ REQUIRED PATTERN FOR ALL SPAWN STATES

```bash
# 1. Complete spawning work
echo "✅ Spawned Product Manager for PRD validation"

# 2. Spawn State Manager for state transition
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "SPAWN_PRODUCT_MANAGER_PRD_VALIDATION" \
  --proposed-next-state "WAITING_FOR_PRD_CREATION" \
  --transition-reason "Product Manager spawned for PRD validation"

# 3. Save TODOs per R287
save_todos "SPAWNED_PRODUCT_MANAGER_PRD_VALIDATION"

# 4. R322: Stop conversation (context preservation)
echo "🛑 R322: Stopping after spawn for context preservation"

# 5. R405: CONTINUATION FLAG - MUST BE TRUE CHECKPOINT=R322 FOR SPAWNING!
echo "CONTINUE-SOFTWARE-FACTORY=TRUE CHECKPOINT=R322"

# 6. Exit to end conversation
exit 0
```

**Enhanced Format**: The `CHECKPOINT=R322` context tells the test framework this is a normal R322 checkpoint, enabling automatic continuation.

---

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

---

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

After acknowledging state rules, create verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_SPAWN_PRODUCT_MANAGER_PRD_VALIDATION
echo "$(date +%s) - Rules read and acknowledged for SPAWN_PRODUCT_MANAGER_PRD_VALIDATION" > .state_rules_read_orchestrator_SPAWN_PRODUCT_MANAGER_PRD_VALIDATION
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

---

## 📋 PRIMARY DIRECTIVES FOR SPAWN_PRODUCT_MANAGER_PRD_VALIDATION STATE

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
   - Summary: MUST update orchestrator-state-v3.json before EVERY state transition (via State Manager in SF 3.0)

4. **🔴🔴🔴 R322** - MANDATORY STOP BEFORE STATE TRANSITIONS (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-before-state-transitions.md`
   - Criticality: SUPREME - -100% penalty for violations
   - Summary: MUST stop after spawning agents, await /continue-orchestrating

### State-Specific Rules:

5. **🚨🚨🚨 R313** - Stop After Spawn (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R313-stop-after-agent-spawn.md`
   - Criticality: BLOCKING - Must stop after spawning
   - Summary: After spawning Product Manager, MUST stop and await continuation
   - Rationale: Prevents context/rule loss in spawned agent

6. **🔴🔴🔴 R405** - Automation Continuation Flag (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md`
   - Criticality: SUPREME - -100% penalty for missing
   - Summary: MUST output CONTINUE-SOFTWARE-FACTORY=TRUE CHECKPOINT=R322 before exit
   - Usage: TRUE for normal spawning (this is normal operation!)

---

## State Manager Bookend Pattern (MANDATORY - SF 3.0)

**BEFORE this state**:
- State Manager validated transition to this state via STARTUP_CONSULTATION
- You are here because State Manager directed you here
- User has completed PRD editing (from WAITING_FOR_PRD_VALIDATION)

**DURING this state**:
- Perform state-specific work (spawn Product Manager for validation)
- NEVER update state files directly
- NEVER call update_state function
- Prepare results for State Manager

**AFTER this state**:
- Spawn State Manager SHUTDOWN_CONSULTATION
- Provide results and proposed next state
- State Manager decides actual next state
- Transition to State Manager's required_next_state

**PROHIBITED**:
- ❌ Calling update_state directly
- ❌ Updating orchestrator-state-v3.json directly
- ❌ Setting validated_by: "orchestrator"
- ❌ Bypassing State Manager consultation

---

## 🎯 STATE OBJECTIVES - SPAWN PRODUCT MANAGER FOR PRD VALIDATION

In the SPAWN_PRODUCT_MANAGER_PRD_VALIDATION state, the ORCHESTRATOR is responsible for:

1. **Verify PRD File Exists and Was Edited**
   - Check that prd/${PROJECT_NAME}-prd.md exists
   - Verify file was modified (check timestamp or content)
   - Ensure user has filled [NEEDS INPUT] sections

2. **Spawn Product Manager Agent for Validation**
   - Spawn product-manager agent in PRD_VALIDATION state
   - Provide PRD file path to agent
   - Agent will:
     - Check for remaining [NEEDS INPUT] markers
     - Validate PRD completeness
     - Output CONTINUE=TRUE if complete
     - Output CONTINUE=FALSE and update gap report if still incomplete

3. **Record Spawn**
   - Update orchestrator-state-v3.json with spawn record
   - Record spawn timestamp
   - Record agent type and target state

4. **Stop Per R313**
   - Stop immediately after spawning
   - Emit CONTINUE-SOFTWARE-FACTORY=TRUE CHECKPOINT=R322
   - Exit with exit 0
   - Next state will be WAITING_FOR_PRD_CREATION (which will check PM's validation result)

---

## 📝 IMMEDIATE ACTIONS UPON ENTERING STATE

### Step 1: Verify PRD File Exists

```bash
echo "🔍 Verifying PRD file exists and was edited..."

STATE_FILE="$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json"
PROJECT_NAME=$(jq -r '.project_name' "$STATE_FILE")
PRD_FILE="$CLAUDE_PROJECT_DIR/prd/${PROJECT_NAME}-prd.md"

if [ ! -f "$PRD_FILE" ]; then
    echo "❌ CRITICAL: PRD file not found at $PRD_FILE"
    echo "User was supposed to complete PRD before continuing"
    PROPOSED_NEXT_STATE="ERROR_RECOVERY"
    TRANSITION_REASON="PRD file missing after user editing"
    ERROR_OCCURRED="true"
    exit 1
fi

# Check file is not empty
if [ ! -s "$PRD_FILE" ]; then
    echo "❌ CRITICAL: PRD file exists but is empty"
    PROPOSED_NEXT_STATE="ERROR_RECOVERY"
    TRANSITION_REASON="PRD file is empty after user editing"
    ERROR_OCCURRED="true"
    exit 1
fi

echo "✅ PRD file verified: $PRD_FILE"
echo "   File size: $(wc -c < "$PRD_FILE") bytes"
echo "   Last modified: $(stat -c %y "$PRD_FILE" 2>/dev/null || stat -f %Sm "$PRD_FILE")"
```

### Step 2: Quick Check for Remaining Gaps

```bash
echo "🔍 Quick check for remaining [NEEDS INPUT] markers..."

NEEDS_INPUT_COUNT=$(grep -c "\[NEEDS INPUT\]" "$PRD_FILE" || echo "0")

if [ "$NEEDS_INPUT_COUNT" -gt 0 ]; then
    echo "⚠️ NOTICE: Still found $NEEDS_INPUT_COUNT [NEEDS INPUT] markers"
    echo "   Product Manager will validate and may ask for more input"
    echo "   This is normal - validation will provide detailed feedback"
else
    echo "✅ No [NEEDS INPUT] markers found - appears complete"
    echo "   Product Manager will perform final validation"
fi
```

### Step 3: Spawn Product Manager Agent for Validation

```bash
echo "🚀 Spawning Product Manager agent for PRD validation..."

# Spawn product-manager agent in PRD_VALIDATION state
# Agent will:
# - Check for remaining [NEEDS INPUT] markers
# - Validate all required sections are present and complete
# - Check content quality and completeness
# - Output CONTINUE=TRUE if PRD passes validation
# - Output CONTINUE=FALSE and update gap report if validation fails

claude-code --agent product-manager \
    --state PRD_VALIDATION \
    --project-name "$PROJECT_NAME" \
    --prd-file "$PRD_FILE" \
    --output-dir "$CLAUDE_PROJECT_DIR/prd"

echo "✅ Product Manager agent spawned successfully"
echo "   Agent will validate PRD and output continuation flag"
```

### Step 4: Record Spawn in State File

**NOTE**: In SF 3.0, State Manager handles state file updates. We just prepare the information:

```bash
# Prepare spawn record for State Manager
SPAWN_TIMESTAMP=$(date -u +%Y-%m-%dT%H:%M:%SZ)

echo "📝 Spawn record prepared:"
echo "   Agent: product-manager"
echo "   State: PRD_VALIDATION"
echo "   Spawned at: $SPAWN_TIMESTAMP"
echo "   Project: $PROJECT_NAME"
echo "   PRD file: $PRD_FILE"

# State Manager will update state file during SHUTDOWN_CONSULTATION
```

---

## 🔴🔴🔴 STATE COMPLETION CHECKLIST 🔴🔴🔴

**Execute these steps IN ORDER to properly complete SPAWN_PRODUCT_MANAGER_PRD_VALIDATION:**

### ✅ Step 1: Complete State-Specific Work
**Refer to "IMMEDIATE ACTIONS UPON ENTERING STATE" section above.**

Once all state work is complete (Product Manager spawned for validation), proceed to mandatory exit protocol:

---

### ✅ Step 2: Set Proposed Next State and Transition Reason
```bash
# Product Manager has been spawned for validation
# Next state is WAITING_FOR_PRD_CREATION (which will check validation result)
PROPOSED_NEXT_STATE="WAITING_FOR_PRD_CREATION"
TRANSITION_REASON="Product Manager spawned for PRD validation - awaiting validation result"
echo "Proposed next state: $PROPOSED_NEXT_STATE"
echo "Transition reason: $TRANSITION_REASON"
```

**NOTE**: It might seem odd to go back to WAITING_FOR_PRD_CREATION, but this is correct!
- WAITING_FOR_PRD_CREATION checks PM's continuation flag
- If validation passes (PM outputs TRUE), it routes to SPAWN_ARCHITECT_MASTER_PLANNING
- If validation fails (PM outputs FALSE), it routes back to WAITING_FOR_PRD_VALIDATION
- This creates a validation loop until PRD is complete

---

### ✅ Step 3: Spawn State Manager for SHUTDOWN_CONSULTATION
```bash
# State Manager validates transition and updates state files (SF 3.0 Pattern)
echo "🔄 Spawning State Manager for SHUTDOWN_CONSULTATION..."

# Prepare work results summary
WORK_RESULTS=$(cat <<EOF
{
  "state_completed": "SPAWN_PRODUCT_MANAGER_PRD_VALIDATION",
  "work_accomplished": [
    "Verified PRD file exists and was edited: $PRD_FILE",
    "Checked for remaining [NEEDS INPUT] markers: $NEEDS_INPUT_COUNT found",
    "Spawned Product Manager agent for PRD validation",
    "Prepared spawn record"
  ],
  "spawn_details": {
    "agent": "product-manager",
    "agent_state": "PRD_VALIDATION",
    "spawned_at": "$SPAWN_TIMESTAMP",
    "project_name": "$PROJECT_NAME",
    "prd_file": "$PRD_FILE",
    "remaining_gaps": $NEEDS_INPUT_COUNT
  },
  "proposed_next_state": "$PROPOSED_NEXT_STATE",
  "transition_reason": "$TRANSITION_REASON"
}
EOF
)

# Spawn State Manager
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "SPAWN_PRODUCT_MANAGER_PRD_VALIDATION" \
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

### ✅ Step 4: Save TODOs (R287 - SUPREME LAW)
```bash
# Save TODO state before transition (R287 trigger)
save_todos "SPAWNED_PRODUCT_MANAGER_PRD_VALIDATION"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - SPAWN_PRODUCT_MANAGER_PRD_VALIDATION complete [R287]"; then
    echo "❌ ERROR: Failed to commit TODO files"
    echo "This is non-fatal but TODOs may be lost in compaction"
    echo "Proceeding with state execution..."
    # Don't exit - TODO commit failure is not fatal
fi

git push || echo "⚠️ WARNING: TODO push failed - committed locally"
echo "✅ TODOs saved and committed"
```

---

### ✅ Step 5: Output Continuation Flag (R405 - SUPREME LAW) ⚠️ MANDATORY
```bash
# Output continuation flag as LAST action (R405)
# Enhanced format with CHECKPOINT=R322 for spawn states
# This tells automation this is a normal R322 checkpoint, enabling auto-continue

echo "CONTINUE-SOFTWARE-FACTORY=TRUE CHECKPOINT=R322"
```

**⚠️ THIS MUST BE THE ABSOLUTE LAST LINE OF OUTPUT BEFORE EXIT! ⚠️**

**Enhanced Format**: The `CHECKPOINT=R322` context is **MANDATORY** for R322 checkpoints.
- Tells framework this is normal operation (not error)
- Enables automatic continuation in tests
- Makes intent explicit in logs

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
- Missing Step 1: No Product Manager spawned = validation doesn't happen
- Missing Step 2: No next state = stuck forever
- Missing Step 3: No State Manager consultation = bypassing bookend pattern (R288 violation, -100%)
- Missing Step 4: No TODO save = work lost (R287 violation, -20% to -100%)
- Missing Step 5: No continuation flag = automation stops (R405 violation, -100%)
- Missing Step 6: No exit = R322 violation (-100%)

**ALL 6 STEPS ARE MANDATORY - NO EXCEPTIONS**
**NOTE**: State file validation and commits now handled by State Manager (SF 3.0 pattern)

---

## STATE TRANSITIONS

### Valid Next States (from state-machines/software-factory-3.0-state-machine.json):

After SPAWN_PRODUCT_MANAGER_PRD_VALIDATION, the valid next states are:

1. **WAITING_FOR_PRD_CREATION** (normal path - validation result check)
   - Guard: Product Manager spawned successfully
   - Condition: Spawn record exists in orchestrator-state-v3.json
   - Action: Check PM's validation result and route accordingly
   - **This is the expected transition 99.9% of the time**
   - Note: Yes, going "back" to WAITING_FOR_PRD_CREATION is correct! It handles routing.

2. **WAITING_FOR_PRD_VALIDATION** (alternative direct path if validation fails)
   - Guard: Product Manager determined PRD still incomplete
   - Condition: Validation failed, loop back for more user input
   - **Less common - usually goes through WAITING_FOR_PRD_CREATION first**

3. **SPAWN_ARCHITECT_MASTER_PLANNING** (alternative direct path if validation passes)
   - Guard: Product Manager determined PRD complete
   - Condition: Validation passed, ready for architecture
   - **Less common - usually goes through WAITING_FOR_PRD_CREATION first**

4. **ERROR_RECOVERY** (error path)
   - Guard: Spawn failed or critical error
   - Condition: Spawn errors detected
   - **Only use for catastrophic failures**

### Routing Logic (via WAITING_FOR_PRD_CREATION)

The typical flow is:
```
SPAWN_PRODUCT_MANAGER_PRD_VALIDATION
  ↓
WAITING_FOR_PRD_CREATION (checks PM's flag)
  ↓
If PM flag = TRUE: → SPAWN_ARCHITECT_MASTER_PLANNING (PRD complete)
If PM flag = FALSE: → WAITING_FOR_PRD_VALIDATION (more user input needed)
```

This creates a validation loop:
```
User edits PRD
  ↓
SPAWN_PRODUCT_MANAGER_PRD_VALIDATION
  ↓
WAITING_FOR_PRD_CREATION
  ↓ (if incomplete)
WAITING_FOR_PRD_VALIDATION
  ↓
User edits PRD again
  ↓
(loop continues until PRD complete)
```

### Invalid Transitions

❌ **FORBIDDEN**:
- SPAWN_PRODUCT_MANAGER_PRD_VALIDATION → SPAWN_PRODUCT_MANAGER_PRD_CREATION - Going backwards
- SPAWN_PRODUCT_MANAGER_PRD_VALIDATION → WAVE_START - Skips all planning
- SPAWN_PRODUCT_MANAGER_PRD_VALIDATION → Any state not in allowed_transitions

---

## ✅ EXIT CHECKLIST (R405 - Continuation Flag Protocol) 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

### 🚨 CRITICAL DISTINCTION: AGENT STOPS ≠ FACTORY STOPS 🚨

**TWO INDEPENDENT DECISIONS - DO NOT CONFUSE THEM:**

#### 1. Should Agent Stop Work? (R322 Technical Requirement)
- Agent completes current state (spawns Product Manager for validation)
- Agent saves TODOs and State Manager updates state
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
echo "✅ Product Manager spawned for PRD validation"

# 2. Set proposed next state
PROPOSED_NEXT_STATE="WAITING_FOR_PRD_CREATION"
TRANSITION_REASON="Product Manager spawned, awaiting validation result"

# 3. Spawn State Manager for state transition
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "SPAWN_PRODUCT_MANAGER_PRD_VALIDATION" \
  --proposed-next-state "$PROPOSED_NEXT_STATE" \
  --transition-reason "$TRANSITION_REASON"
# State Manager updates all 4 state files atomically

# 4. Save TODOs
save_todos "R322_CHECKPOINT"

# 5. Factory continues (operational status)
echo "CONTINUE-SOFTWARE-FACTORY=TRUE CHECKPOINT=R322"

# 6. Agent stops (technical requirement)
exit 0
```

**Both happen together! Agent stops AND factory continues!**

### WHEN TO USE EACH FLAG VALUE

**TRUE (99.9% - Use for this state):**
- ✅ R322 checkpoint reached
- ✅ Product Manager spawned successfully
- ✅ Ready for /continue-orchestrating
- ✅ Waiting for user to continue (NORMAL)
- ✅ Normal spawning operation

**FALSE (0.1% - NOT for normal spawning):**
- ❌ CATASTROPHIC unrecoverable error
- ❌ Data corruption spreading
- ❌ Critical security violation
- ❌ NOT for R322 checkpoints
- ❌ NOT for normal spawn operations

**See**: `$CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md`

---

## Common Mistakes to Avoid

1. ❌ **Using CONTINUE-SOFTWARE-FACTORY=FALSE for normal spawning** - This kills automation!
2. ❌ **Not stopping after spawning** - Violates R313 and R322
3. ❌ **Updating state files directly** - Violates SF 3.0 bookend pattern
4. ❌ **Forgetting CHECKPOINT=R322 in continuation flag** - Breaks test automation
5. ❌ **Not creating R290 verification marker** - Automatic -100% penalty
6. ❌ **Missing mandatory rule readings** - Automatic failure
7. ❌ **Being confused by going "back" to WAITING_FOR_PRD_CREATION** - This is correct routing!

---

**END OF SPAWN_PRODUCT_MANAGER_PRD_VALIDATION STATE RULES**
