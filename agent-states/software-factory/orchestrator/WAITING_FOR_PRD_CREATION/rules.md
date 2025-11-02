# Orchestrator - WAITING_FOR_PRD_CREATION State Rules

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
## 🟢 DECISION STATE - CHECK FLAG AND ROUTE TO CORRECT NEXT STATE 🟢

**CRITICAL: This is a DECISION state that reads the CONTINUE-SOFTWARE-FACTORY flag from Product Manager and routes accordingly**

### STATE PURPOSE:
- **Read** Product Manager's output and CONTINUE flag
- **Check** if PRD generation was complete (TRUE) or incomplete (FALSE)
- **Route** to correct next state based on flag
- **Stop** with appropriate continuation flag

### DECISION LOGIC:
```bash
# Check Product Manager's CONTINUE-SOFTWARE-FACTORY flag
PM_CONTINUE_FLAG=$(grep "CONTINUE-SOFTWARE-FACTORY" "$PM_OUTPUT" | cut -d= -f2 | cut -d' ' -f1)

if [ "$PM_CONTINUE_FLAG" = "TRUE" ]; then
    # PRD description was complete - PM generated full PRD
    NEXT_STATE="SPAWN_ARCHITECT_MASTER_PLANNING"
    CONTINUE_FLAG="TRUE"  # Automation can proceed
elif [ "$PM_CONTINUE_FLAG" = "FALSE" ]; then
    # PRD description was incomplete - needs human input
    NEXT_STATE="WAITING_FOR_PRD_VALIDATION"
    CONTINUE_FLAG="FALSE"  # Automation must stop - human needed
else
    # Error case
    NEXT_STATE="ERROR_RECOVERY"
    CONTINUE_FLAG="FALSE"
fi
```

---

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

---

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

After acknowledging state rules, create verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
mkdir -p markers/state-verification && touch "markers/state-verification/state_rules_read_orchestrator_WAITING_FOR_PRD_CREATION-$(date +%Y%m%d-%H%M%S)"
echo "$(date +%s) - Rules read and acknowledged for WAITING_FOR_PRD_CREATION" > "markers/state-verification/state_rules_read_orchestrator_WAITING_FOR_PRD_CREATION-$(date +%Y%m%d-%H%M%S)"
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

---

## 📋 PRIMARY DIRECTIVES FOR WAITING_FOR_PRD_CREATION STATE

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
   - Summary: MUST stop after completing state work

### State-Specific Rules:

5. **🔴🔴🔴 R405** - Automation Continuation Flag (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md`
   - Criticality: SUPREME - -100% penalty for missing
   - Summary: MUST output CONTINUE-SOFTWARE-FACTORY flag based on Product Manager's result
   - Usage: TRUE if PM completed PRD, FALSE if PM needs human input

---

## State Manager Bookend Pattern (MANDATORY - SF 3.0)

**BEFORE this state**:
- State Manager validated transition to this state via STARTUP_CONSULTATION
- You are here because State Manager directed you here
- Product Manager was spawned and has completed PRD_CREATION

**DURING this state**:
- Check Product Manager's output and continuation flag
- Determine routing based on flag
- NEVER update state files directly
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

## 🎯 STATE OBJECTIVES - CHECK PRD CREATION RESULT AND ROUTE

In the WAITING_FOR_PRD_CREATION state, the ORCHESTRATOR is responsible for:

1. **Verify PRD File Exists**
   - Check that prd/${PROJECT_NAME}-prd.md was created
   - Verify file is not empty
   - Record PRD creation in state file

2. **Read Product Manager's Continuation Flag**
   - Check CONTINUE-SOFTWARE-FACTORY flag from PM output
   - TRUE = Description was complete, PRD is finished
   - FALSE = Description was incomplete, needs human input

3. **Make Routing Decision**
   - If TRUE: Route to SPAWN_ARCHITECT_MASTER_PLANNING (PRD complete, proceed)
   - If FALSE: Route to WAITING_FOR_PRD_VALIDATION (needs human completion)
   - If error: Route to ERROR_RECOVERY

4. **Update State File**
   - Set prd_exists = true
   - Record PRD file path
   - Update next state based on routing decision

5. **Stop and Emit Appropriate Flag**
   - Stop per R322
   - Emit CONTINUE-SOFTWARE-FACTORY=TRUE if routing to architecture
   - Emit CONTINUE-SOFTWARE-FACTORY=FALSE if routing to human validation

---

## 📝 IMMEDIATE ACTIONS UPON ENTERING STATE

### Step 1: Verify PRD File Was Created

```bash
echo "🔍 Verifying PRD file was created..."

STATE_FILE="$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json"
PROJECT_NAME=$(jq -r '.project_name' "$STATE_FILE")
PRD_FILE="$CLAUDE_PROJECT_DIR/prd/${PROJECT_NAME}-prd.md"

if [ ! -f "$PRD_FILE" ]; then
    echo "❌ CRITICAL: PRD file not found at $PRD_FILE"
    echo "Product Manager did not create PRD file"
    PROPOSED_NEXT_STATE="ERROR_RECOVERY"
    TRANSITION_REASON="PRD file not created by Product Manager"
    ERROR_OCCURRED="true"
    exit 1
fi

# Check file is not empty
if [ ! -s "$PRD_FILE" ]; then
    echo "❌ CRITICAL: PRD file exists but is empty"
    PROPOSED_NEXT_STATE="ERROR_RECOVERY"
    TRANSITION_REASON="PRD file is empty"
    ERROR_OCCURRED="true"
    exit 1
fi

echo "✅ PRD file verified: $PRD_FILE"
echo "   File size: $(wc -c < "$PRD_FILE") bytes"
```

### Step 2: Read Product Manager's Continuation Flag

```bash
echo "📊 Reading Product Manager's continuation flag..."

# Find Product Manager's output/log
# Location depends on spawn method - check common locations
PM_OUTPUT_FILE=$(find "$CLAUDE_PROJECT_DIR" -name "*product-manager*output*" -o -name "*PRD_CREATION*log*" 2>/dev/null | head -1)

if [ -z "$PM_OUTPUT_FILE" ]; then
    echo "⚠️ WARNING: Cannot find Product Manager output file"
    echo "Checking PRD file for [NEEDS INPUT] markers as fallback..."

    # Fallback: Check if PRD has incomplete markers
    if grep -q "\[NEEDS INPUT\]" "$PRD_FILE"; then
        echo "Found [NEEDS INPUT] markers in PRD - incomplete"
        PM_CONTINUE_FLAG="FALSE"
    else
        echo "No [NEEDS INPUT] markers in PRD - assuming complete"
        PM_CONTINUE_FLAG="TRUE"
    fi
else
    echo "✅ Found Product Manager output: $PM_OUTPUT_FILE"

    # Extract continuation flag
    PM_CONTINUE_FLAG=$(grep "CONTINUE-SOFTWARE-FACTORY" "$PM_OUTPUT_FILE" | tail -1 | cut -d= -f2 | cut -d' ' -f1)

    if [ -z "$PM_CONTINUE_FLAG" ]; then
        echo "⚠️ WARNING: No continuation flag found in output"
        echo "Defaulting to FALSE (safe - requires human review)"
        PM_CONTINUE_FLAG="FALSE"
    fi
fi

echo "✅ Product Manager continuation flag: $PM_CONTINUE_FLAG"
```

### Step 3: Make Routing Decision

```bash
echo "🔀 Determining routing based on continuation flag..."

if [ "$PM_CONTINUE_FLAG" = "TRUE" ]; then
    echo "✅ PRD creation COMPLETE - description was sufficient"
    echo "   Routing to: SPAWN_ARCHITECT_MASTER_PLANNING"
    echo "   Reason: PRD finished, ready for architecture planning"

    PROPOSED_NEXT_STATE="SPAWN_ARCHITECT_MASTER_PLANNING"
    TRANSITION_REASON="PRD creation complete - proceeding to architecture planning"
    ORCHESTRATOR_CONTINUE_FLAG="TRUE"  # Automation can continue

elif [ "$PM_CONTINUE_FLAG" = "FALSE" ]; then
    echo "⚠️ PRD creation INCOMPLETE - description insufficient"
    echo "   Routing to: WAITING_FOR_PRD_VALIDATION"
    echo "   Reason: PRD has [NEEDS INPUT] gaps requiring human completion"

    PROPOSED_NEXT_STATE="WAITING_FOR_PRD_VALIDATION"
    TRANSITION_REASON="PRD incomplete - requires human input for gaps"
    ORCHESTRATOR_CONTINUE_FLAG="FALSE"  # Automation must stop - human needed

    # Check for validation report
    VALIDATION_REPORT="$CLAUDE_PROJECT_DIR/prd/PRD-VALIDATION-REPORT.md"
    if [ -f "$VALIDATION_REPORT" ]; then
        echo "📋 Validation report available: $VALIDATION_REPORT"
        echo "   User should review this for gap details"
    fi

else
    echo "❌ ERROR: Invalid or missing continuation flag"
    echo "   Flag value: '$PM_CONTINUE_FLAG'"

    PROPOSED_NEXT_STATE="ERROR_RECOVERY"
    TRANSITION_REASON="Invalid continuation flag from Product Manager"
    ORCHESTRATOR_CONTINUE_FLAG="FALSE"
fi

echo "✅ Routing decision made:"
echo "   Next state: $PROPOSED_NEXT_STATE"
echo "   Continue flag: $ORCHESTRATOR_CONTINUE_FLAG"
```

---

## 🔴🔴🔴 STATE COMPLETION CHECKLIST 🔴🔴🔴

**Execute these steps IN ORDER to properly complete WAITING_FOR_PRD_CREATION:**

### ✅ Step 1: Complete State-Specific Work
**Refer to "IMMEDIATE ACTIONS UPON ENTERING STATE" section above.**

Once all state work is complete (PRD verified, flag read, routing determined), proceed to mandatory exit protocol:

---

### ✅ Step 2: Set Proposed Next State and Transition Reason
```bash
# Already determined in Step 3 of immediate actions
echo "Proposed next state: $PROPOSED_NEXT_STATE"
echo "Transition reason: $TRANSITION_REASON"
echo "Continuation flag: $ORCHESTRATOR_CONTINUE_FLAG"
```

---

### ✅ Step 3: Spawn State Manager for SHUTDOWN_CONSULTATION
```bash
# State Manager validates transition and updates state files (SF 3.0 Pattern)
echo "🔄 Spawning State Manager for SHUTDOWN_CONSULTATION..."

# Prepare work results summary
WORK_RESULTS=$(cat <<EOF
{
  "state_completed": "WAITING_FOR_PRD_CREATION",
  "work_accomplished": [
    "Verified PRD file created: $PRD_FILE",
    "Read Product Manager continuation flag: $PM_CONTINUE_FLAG",
    "Made routing decision based on flag",
    "Updated prd_exists field in state"
  ],
  "prd_details": {
    "file_path": "$PRD_FILE",
    "file_size": $(wc -c < "$PRD_FILE"),
    "creation_complete": $([ "$PM_CONTINUE_FLAG" = "TRUE" ] && echo "true" || echo "false"),
    "needs_human_input": $([ "$PM_CONTINUE_FLAG" = "FALSE" ] && echo "true" || echo "false")
  },
  "proposed_next_state": "$PROPOSED_NEXT_STATE",
  "transition_reason": "$TRANSITION_REASON",
  "orchestrator_continue_flag": "$ORCHESTRATOR_CONTINUE_FLAG"
}
EOF
)

# Spawn State Manager
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "WAITING_FOR_PRD_CREATION" \
  --proposed-next-state "$PROPOSED_NEXT_STATE" \
  --transition-reason "$TRANSITION_REASON" \
  --work-results "$WORK_RESULTS"

# State Manager will:
# 1. Validate PROPOSED_NEXT_STATE exists and transition is valid
# 2. Update all 4 state files atomically (R288)
# 3. Set prd_exists = true in state file
# 4. Commit and push state files
# 5. Return REQUIRED_NEXT_STATE (usually same as proposed unless invalid)

echo "✅ State Manager consultation complete"
echo "✅ State files updated by State Manager"
```

---

### ✅ Step 4: Save TODOs (R287 - SUPREME LAW)
```bash
# Save TODO state before transition (R287 trigger)
save_todos "WAITING_FOR_PRD_CREATION_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - WAITING_FOR_PRD_CREATION complete [R287]"; then
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
# Use the ORCHESTRATOR_CONTINUE_FLAG determined in Step 3 of immediate actions

if [ "$ORCHESTRATOR_CONTINUE_FLAG" = "TRUE" ]; then
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=PRD_COMPLETE"
else
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=HUMAN_INPUT_REQUIRED"
fi
```

**⚠️ THIS MUST BE THE ABSOLUTE LAST LINE OF OUTPUT BEFORE EXIT! ⚠️**

**CRITICAL DISTINCTION FOR THIS STATE**:
- If routing to SPAWN_ARCHITECT_MASTER_PLANNING → Use TRUE (automation continues)
- If routing to WAITING_FOR_PRD_VALIDATION → Use FALSE (human input required)

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
- Missing Step 1: No PRD verification = might route incorrectly
- Missing Step 2: No routing decision = stuck forever
- Missing Step 3: No State Manager consultation = bypassing bookend pattern (R288 violation, -100%)
- Missing Step 4: No TODO save = work lost (R287 violation, -20% to -100%)
- Missing Step 5: No continuation flag = automation doesn't know how to proceed (R405 violation, -100%)
- Missing Step 6: No exit = R322 violation (-100%)

**ALL 6 STEPS ARE MANDATORY - NO EXCEPTIONS**
**NOTE**: State file validation and commits now handled by State Manager (SF 3.0 pattern)

---

## STATE TRANSITIONS

### Valid Next States (from state-machines/software-factory-3.0-state-machine.json):

After WAITING_FOR_PRD_CREATION, the valid next states are:

1. **SPAWN_ARCHITECT_MASTER_PLANNING** (complete path)
   - Guard: PRD creation complete (PM flag = TRUE)
   - Condition: PRD file exists with no [NEEDS INPUT] markers
   - Continuation flag: TRUE (automation proceeds)
   - **This path means: PRD fully generated, ready for architecture**

2. **WAITING_FOR_PRD_VALIDATION** (incomplete path)
   - Guard: PRD creation incomplete (PM flag = FALSE)
   - Condition: PRD file has [NEEDS INPUT] markers
   - Continuation flag: FALSE (human input required)
   - **This path means: PRD has gaps, needs human to fill them**

3. **ERROR_RECOVERY** (error path)
   - Guard: PRD file missing or invalid flag
   - Condition: Critical errors detected
   - Continuation flag: FALSE
   - **Only use for errors, not for incomplete PRD**

### Decision Matrix

| PM Continue Flag | PRD Has [NEEDS INPUT] | Next State | Orchestrator Continue Flag |
|---|---|---|---|
| TRUE | No | SPAWN_ARCHITECT_MASTER_PLANNING | TRUE |
| FALSE | Yes | WAITING_FOR_PRD_VALIDATION | FALSE |
| Missing/Invalid | N/A | ERROR_RECOVERY | FALSE |

### Invalid Transitions

❌ **FORBIDDEN**:
- WAITING_FOR_PRD_CREATION → SPAWN_PRODUCT_MANAGER_PRD_VALIDATION - Skips waiting for human input
- WAITING_FOR_PRD_CREATION → SPAWN_ARCHITECT_PHASE_PLANNING - Skips master architecture
- WAITING_FOR_PRD_CREATION → WAVE_START - Skips all planning
- WAITING_FOR_PRD_CREATION → Any state not in allowed_transitions

---

## ✅ EXIT CHECKLIST (R405 - Continuation Flag Protocol) 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

### 🚨 CRITICAL DISTINCTION: AGENT STOPS ≠ FACTORY STOPS 🚨

**TWO INDEPENDENT DECISIONS - DO NOT CONFUSE THEM:**

#### 1. Should Agent Stop Work? (R322 Technical Requirement)
- Agent completes current state (checks PRD and routes)
- Agent saves TODOs and State Manager updates state
- Agent exits with `exit 0` (preserves context)
- User runs /continue-orchestrating to resume
- **This is NORMAL at checkpoints**

#### 2. Should Factory Continue? (R405 Operational Status)
- **Depends on routing decision!**
- TRUE if PRD complete (going to architecture) = Automation can proceed
- FALSE if PRD incomplete (going to human validation) = Human must act
- **This state can use EITHER flag depending on the situation**

### THE PATTERN FOR THIS STATE (SF 3.0)

```bash
# 1. Complete state work
echo "✅ PRD verified, routing decision made"

# 2. Set proposed next state based on PM flag
if [ "$PM_CONTINUE_FLAG" = "TRUE" ]; then
    PROPOSED_NEXT_STATE="SPAWN_ARCHITECT_MASTER_PLANNING"
    ORCHESTRATOR_FLAG="TRUE"
else
    PROPOSED_NEXT_STATE="WAITING_FOR_PRD_VALIDATION"
    ORCHESTRATOR_FLAG="FALSE"
fi

# 3. Spawn State Manager for state transition
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "WAITING_FOR_PRD_CREATION" \
  --proposed-next-state "$PROPOSED_NEXT_STATE" \
  --transition-reason "[REASON]"

# 4. Save TODOs
save_todos "PRD_CREATION_ROUTED"

# 5. Factory continues or stops (operational status - DEPENDS ON ROUTING!)
echo "CONTINUE-SOFTWARE-FACTORY=$ORCHESTRATOR_FLAG REASON=[REASON]"

# 6. Agent stops (technical requirement)
exit 0
```

**Agent ALWAYS stops, but factory continuation DEPENDS on PRD completeness**

### WHEN TO USE EACH FLAG VALUE (FOR THIS STATE)

**TRUE (when PRD is complete):**
- ✅ PRD generation finished (no [NEEDS INPUT] markers)
- ✅ Routing to SPAWN_ARCHITECT_MASTER_PLANNING
- ✅ Product Manager returned TRUE
- ✅ Automation can continue to architecture planning

**FALSE (when PRD needs human input):**
- ✅ PRD has [NEEDS INPUT] gaps
- ✅ Routing to WAITING_FOR_PRD_VALIDATION
- ✅ Product Manager returned FALSE
- ✅ Human must complete PRD before continuing

**FALSE (error cases):**
- ❌ PRD file missing
- ❌ Invalid continuation flag from PM
- ❌ Critical errors detected

**See**: `$CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md`

---

## Common Mistakes to Avoid

1. ❌ **Always using TRUE regardless of PM flag** - Breaks human-in-loop for incomplete PRDs
2. ❌ **Always using FALSE for this state** - Prevents automation from proceeding with complete PRDs
3. ❌ **Not checking for PRD file existence** - Might route to architecture without PRD
4. ❌ **Ignoring [NEEDS INPUT] markers** - Bypasses human validation when needed
5. ❌ **Not creating R290 verification marker** - Automatic -100% penalty
6. ❌ **Updating state files directly** - Violates SF 3.0 bookend pattern

---

**END OF WAITING_FOR_PRD_CREATION STATE RULES**
