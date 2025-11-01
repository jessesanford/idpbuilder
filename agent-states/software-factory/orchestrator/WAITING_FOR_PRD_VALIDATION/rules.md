# Orchestrator - WAITING_FOR_PRD_VALIDATION State Rules

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
## 🟡 HUMAN-IN-LOOP STATE - USER MUST COMPLETE PRD GAPS 🟡

**CRITICAL: This is a HUMAN-IN-LOOP state where the agent MUST stop and wait for human action**

### STATE PURPOSE:
- **Notify** user that PRD has gaps needing human input
- **Display** gap report showing what needs to be filled
- **Wait** for human to edit PRD and fill [NEEDS INPUT] sections
- **Stop** with CONTINUE-SOFTWARE-FACTORY=FALSE (human action required)
- **Resume** when user runs /continue-orchestrating after editing PRD

### HUMAN ACTION REQUIRED:
```markdown
📋 PRD COMPLETION REQUIRED

Your PRD has sections marked [NEEDS INPUT] that require your expertise.

1. Open: prd/${PROJECT_NAME}-prd.md
2. Find all [NEEDS INPUT] markers
3. Replace markers with actual content
4. Save the file
5. Run: /continue-orchestrating

A validation report is available at: prd/PRD-VALIDATION-REPORT.md
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
touch .state_rules_read_orchestrator_WAITING_FOR_PRD_VALIDATION
echo "$(date +%s) - Rules read and acknowledged for WAITING_FOR_PRD_VALIDATION" > .state_rules_read_orchestrator_WAITING_FOR_PRD_VALIDATION
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

---

## 📋 PRIMARY DIRECTIVES FOR WAITING_FOR_PRD_VALIDATION STATE

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
   - Summary: MUST output CONTINUE-SOFTWARE-FACTORY=FALSE (human action required)
   - Usage: **ALWAYS FALSE for this state** - human must complete PRD gaps

6. **🟡🟡🟡 Human-in-Loop Protocol** (STATE-SPECIFIC)
   - Criticality: CRITICAL - User experience requirement
   - Summary: MUST display clear instructions to user about what to do
   - Must include: file path, what to edit, how to continue
   - Must be user-friendly and actionable

---

## State Manager Bookend Pattern (MANDATORY - SF 3.0)

**BEFORE this state**:
- State Manager validated transition to this state via STARTUP_CONSULTATION
- You are here because PRD has [NEEDS INPUT] gaps
- Previous state (WAITING_FOR_PRD_CREATION or SPAWN_PRODUCT_MANAGER_PRD_VALIDATION) determined PRD is incomplete

**DURING this state**:
- Display gap report to user
- Provide clear editing instructions
- Update state file to record waiting status
- NEVER update state files directly
- Prepare to stop and wait

**AFTER this state**:
- Spawn State Manager SHUTDOWN_CONSULTATION
- Set next state to SPAWN_PRODUCT_MANAGER_PRD_VALIDATION
- State Manager updates state files
- Stop with CONTINUE-SOFTWARE-FACTORY=FALSE
- User will edit PRD and run /continue-orchestrating

**PROHIBITED**:
- ❌ Calling update_state directly
- ❌ Updating orchestrator-state-v3.json directly
- ❌ Setting validated_by: "orchestrator"
- ❌ Bypassing State Manager consultation
- ❌ Using CONTINUE-SOFTWARE-FACTORY=TRUE (human action is required!)

---

## 🎯 STATE OBJECTIVES - NOTIFY USER AND WAIT FOR PRD COMPLETION

In the WAITING_FOR_PRD_VALIDATION state, the ORCHESTRATOR is responsible for:

1. **Display Gap Report**
   - Read prd/PRD-VALIDATION-REPORT.md
   - Show user what sections need input
   - Explain what information is missing

2. **Provide Clear Instructions**
   - Tell user exactly which file to edit
   - Explain how to find [NEEDS INPUT] markers
   - Describe what to fill in
   - Tell user how to resume after editing

3. **Update State File**
   - Record that system is waiting for human input
   - Set human_action_required = true
   - Set next_state = SPAWN_PRODUCT_MANAGER_PRD_VALIDATION (for after user completes)

4. **Stop and Wait**
   - Stop per R322
   - Emit CONTINUE-SOFTWARE-FACTORY=FALSE (human must act)
   - Exit and wait for user to run /continue-orchestrating

---

## 📝 IMMEDIATE ACTIONS UPON ENTERING STATE

### Step 1: Verify PRD and Gap Report Exist

```bash
echo "🔍 Verifying PRD and gap report exist..."

STATE_FILE="$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json"
PROJECT_NAME=$(jq -r '.project_name' "$STATE_FILE")
PRD_FILE="$CLAUDE_PROJECT_DIR/prd/${PROJECT_NAME}-prd.md"
GAP_REPORT="$CLAUDE_PROJECT_DIR/prd/PRD-VALIDATION-REPORT.md"

if [ ! -f "$PRD_FILE" ]; then
    echo "❌ CRITICAL: PRD file not found at $PRD_FILE"
    PROPOSED_NEXT_STATE="ERROR_RECOVERY"
    TRANSITION_REASON="PRD file missing"
    ERROR_OCCURRED="true"
    exit 1
fi

echo "✅ PRD file verified: $PRD_FILE"

# Gap report is optional but helpful
if [ -f "$GAP_REPORT" ]; then
    echo "✅ Gap report found: $GAP_REPORT"
else
    echo "⚠️ WARNING: Gap report not found (optional but recommended)"
fi
```

### Step 2: Check for [NEEDS INPUT] Markers

```bash
echo "🔍 Checking for [NEEDS INPUT] markers in PRD..."

NEEDS_INPUT_COUNT=$(grep -c "\[NEEDS INPUT\]" "$PRD_FILE" || echo "0")

if [ "$NEEDS_INPUT_COUNT" -eq 0 ]; then
    echo "⚠️ WARNING: No [NEEDS INPUT] markers found in PRD"
    echo "   This state should only be reached if PRD is incomplete"
    echo "   Either:"
    echo "   1. User already filled gaps (proceed to validation)"
    echo "   2. Routing error occurred"
    echo "   Proceeding to validation to check..."

    PROPOSED_NEXT_STATE="SPAWN_PRODUCT_MANAGER_PRD_VALIDATION"
    TRANSITION_REASON="No gaps found - proceeding to validation"
    ORCHESTRATOR_CONTINUE_FLAG="TRUE"  # Can proceed to validation
else
    echo "📊 Found $NEEDS_INPUT_COUNT [NEEDS INPUT] markers in PRD"
    echo "   User must complete these sections before continuing"

    PROPOSED_NEXT_STATE="SPAWN_PRODUCT_MANAGER_PRD_VALIDATION"
    TRANSITION_REASON="PRD gaps present - user must complete then validation will run"
    ORCHESTRATOR_CONTINUE_FLAG="FALSE"  # Human action required
fi
```

### Step 3: Display User Instructions

```bash
echo "📋 Displaying user instructions..."

cat <<EOF

═══════════════════════════════════════════════════════════════
🟡 HUMAN ACTION REQUIRED - PRD COMPLETION NEEDED 🟡
═══════════════════════════════════════════════════════════════

Your Product Requirements Document (PRD) has been generated but
contains sections that require your expert input.

📁 PRD FILE: $PRD_FILE
📊 GAPS FOUND: $NEEDS_INPUT_COUNT [NEEDS INPUT] markers

═══════════════════════════════════════════════════════════════
📝 INSTRUCTIONS:
═══════════════════════════════════════════════════════════════

1. Open the PRD file:
   - Location: $PRD_FILE
   - Use any text editor

2. Find [NEEDS INPUT] markers:
   - Search for "[NEEDS INPUT]" in the file
   - Each marker indicates a section needing completion

3. Replace markers with content:
   - Remove the [NEEDS INPUT] marker
   - Fill in the actual information
   - Be specific and detailed

4. Save the file:
   - Save your changes
   - Ensure file is valid Markdown

5. Resume the workflow:
   - Run: /continue-orchestrating
   - System will validate your changes

═══════════════════════════════════════════════════════════════
📋 GAP ANALYSIS REPORT:
═══════════════════════════════════════════════════════════════

EOF

if [ -f "$GAP_REPORT" ]; then
    echo "A detailed gap analysis is available:"
    echo "   $GAP_REPORT"
    echo ""
    echo "Preview of gaps:"
    head -50 "$GAP_REPORT"
    echo ""
    echo "(See full report at: $GAP_REPORT)"
else
    echo "Search for [NEEDS INPUT] markers in the PRD file."
    echo "Each marker shows what information is needed."
fi

cat <<EOF

═══════════════════════════════════════════════════════════════
⏸️  SYSTEM PAUSED - WAITING FOR YOUR INPUT
═══════════════════════════════════════════════════════════════

After completing the PRD, run:
  /continue-orchestrating

The system will then validate your changes and proceed if complete.

═══════════════════════════════════════════════════════════════

EOF
```

---

## 🔴🔴🔴 STATE COMPLETION CHECKLIST 🔴🔴🔴

**Execute these steps IN ORDER to properly complete WAITING_FOR_PRD_VALIDATION:**

### ✅ Step 1: Complete State-Specific Work
**Refer to "IMMEDIATE ACTIONS UPON ENTERING STATE" section above.**

Once all state work is complete (user notified, instructions displayed), proceed to mandatory exit protocol:

---

### ✅ Step 2: Set Proposed Next State and Transition Reason
```bash
# Next state is always SPAWN_PRODUCT_MANAGER_PRD_VALIDATION
# (User will run /continue-orchestrating after editing, which will spawn PM for validation)
PROPOSED_NEXT_STATE="SPAWN_PRODUCT_MANAGER_PRD_VALIDATION"
TRANSITION_REASON="Waiting for user to complete PRD gaps - will validate after user continues"
echo "Proposed next state: $PROPOSED_NEXT_STATE"
echo "Transition reason: $TRANSITION_REASON"
```

---

### ✅ Step 3: Spawn State Manager for SHUTDOWN_CONSULTATION
```bash
# State Manager validates transition and updates state files (SF 3.0 Pattern)
echo "🔄 Spawning State Manager for SHUTDOWN_CONSULTATION..."

# Prepare work results summary
WORK_RESULTS=$(cat <<EOF
{
  "state_completed": "WAITING_FOR_PRD_VALIDATION",
  "work_accomplished": [
    "Verified PRD file exists: $PRD_FILE",
    "Checked for [NEEDS INPUT] markers: $NEEDS_INPUT_COUNT found",
    "Displayed gap report and user instructions",
    "Set human_action_required flag"
  ],
  "prd_status": {
    "file_path": "$PRD_FILE",
    "needs_input_count": $NEEDS_INPUT_COUNT,
    "gap_report_available": $([ -f "$GAP_REPORT" ] && echo "true" || echo "false"),
    "human_action_required": true
  },
  "proposed_next_state": "$PROPOSED_NEXT_STATE",
  "transition_reason": "$TRANSITION_REASON"
}
EOF
)

# Spawn State Manager
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "WAITING_FOR_PRD_VALIDATION" \
  --proposed-next-state "$PROPOSED_NEXT_STATE" \
  --transition-reason "$TRANSITION_REASON" \
  --work-results "$WORK_RESULTS"

# State Manager will:
# 1. Validate PROPOSED_NEXT_STATE exists and transition is valid
# 2. Update all 4 state files atomically (R288)
# 3. Set human_action_required = true in state file
# 4. Commit and push state files
# 5. Return REQUIRED_NEXT_STATE (usually same as proposed unless invalid)

echo "✅ State Manager consultation complete"
echo "✅ State files updated by State Manager"
```

---

### ✅ Step 4: Save TODOs (R287 - SUPREME LAW)
```bash
# Save TODO state before transition (R287 trigger)
save_todos "WAITING_FOR_USER_PRD_COMPLETION"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - WAITING_FOR_PRD_VALIDATION (human action required) [R287]"; then
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
# **MUST BE FALSE** - human action is required

echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=HUMAN_INPUT_REQUIRED"
```

**⚠️ THIS MUST BE THE ABSOLUTE LAST LINE OF OUTPUT BEFORE EXIT! ⚠️**

**CRITICAL FOR THIS STATE**:
- **ALWAYS use FALSE** - This state exists specifically because human input is needed
- Automation CANNOT continue until user edits PRD and runs /continue-orchestrating
- TRUE would be wrong - system cannot proceed without human completing gaps

---

### ✅ Step 6: Stop Processing (R322 - SUPREME LAW)
```bash
# Stop for human action (R322)
echo "🛑 Stopping - waiting for user to complete PRD and run /continue-orchestrating"
exit 0
```

---

## 🚨 CHECKLIST ENFORCEMENT 🚨

**Skipping ANY step in this checklist = FAILURE:**
- Missing Step 1: User not notified = bad UX, user doesn't know what to do
- Missing Step 2: No next state = stuck forever when user continues
- Missing Step 3: No State Manager consultation = bypassing bookend pattern (R288 violation, -100%)
- Missing Step 4: No TODO save = work lost (R287 violation, -20% to -100%)
- Missing Step 5: No continuation flag OR using TRUE = automation continues without human input (R405 violation, -100%)
- Missing Step 6: No exit = R322 violation (-100%)

**ALL 6 STEPS ARE MANDATORY - NO EXCEPTIONS**
**NOTE**: State file validation and commits now handled by State Manager (SF 3.0 pattern)

---

## STATE TRANSITIONS

### Valid Next States (from state-machines/software-factory-3.0-state-machine.json):

After WAITING_FOR_PRD_VALIDATION, when user runs /continue-orchestrating, the ONLY valid next state is:

1. **SPAWN_PRODUCT_MANAGER_PRD_VALIDATION** (normal path)
   - Guard: User completed PRD editing
   - Condition: User ran /continue-orchestrating
   - Action: Spawn PM to validate user's changes
   - **This is the expected transition 100% of the time**

2. **ERROR_RECOVERY** (error path)
   - Guard: Critical system error
   - Condition: Unrecoverable failure
   - **Only use for catastrophic errors, not for incomplete PRD**

### Invalid Transitions

❌ **FORBIDDEN**:
- WAITING_FOR_PRD_VALIDATION → SPAWN_ARCHITECT_MASTER_PLANNING - Skips validation of user's edits
- WAITING_FOR_PRD_VALIDATION → WAITING_FOR_PRD_CREATION - Going backwards
- WAITING_FOR_PRD_VALIDATION → WAVE_START - Skips all planning
- WAITING_FOR_PRD_VALIDATION → Any state not in allowed_transitions

---

## ✅ EXIT CHECKLIST (R405 - Continuation Flag Protocol) 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

### 🚨 CRITICAL: THIS STATE ALWAYS USES FALSE 🚨

**FOR THIS STATE, THE FLAG IS SIMPLE:**

#### Should Agent Stop Work? (R322)
- YES - Agent completes state work (notifies user)
- Agent exits with `exit 0`

#### Should Factory Continue? (R405)
- **NO - ALWAYS FALSE for this state**
- Human must complete PRD gaps before automation can continue
- User will run /continue-orchestrating after editing
- **Using TRUE here would be WRONG - automation cannot proceed without human**

### THE PATTERN FOR THIS STATE (SF 3.0)

```bash
# 1. Complete state work
echo "✅ User notified about PRD gaps, instructions displayed"

# 2. Set proposed next state
PROPOSED_NEXT_STATE="SPAWN_PRODUCT_MANAGER_PRD_VALIDATION"
TRANSITION_REASON="User must complete PRD, then validation will run"

# 3. Spawn State Manager for state transition
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "WAITING_FOR_PRD_VALIDATION" \
  --proposed-next-state "$PROPOSED_NEXT_STATE" \
  --transition-reason "$TRANSITION_REASON"

# 4. Save TODOs
save_todos "WAITING_FOR_HUMAN_PRD_INPUT"

# 5. Factory CANNOT continue (operational status)
echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=HUMAN_INPUT_REQUIRED"

# 6. Agent stops (technical requirement)
exit 0
```

**Agent stops AND factory stops - both are correct for this state**

### WHY FALSE IS CORRECT (AND MANDATORY) FOR THIS STATE

**This state exists specifically for human-in-loop:**
- ❌ Automation CANNOT fill PRD gaps (requires domain knowledge)
- ❌ System CANNOT proceed without complete PRD
- ✅ Human MUST review and complete gaps
- ✅ User MUST run /continue-orchestrating to resume

**Using TRUE would be catastrophically wrong:**
- Would tell automation to continue despite incomplete PRD
- Would bypass human input requirement
- Would violate the entire purpose of this state
- Would potentially pass incomplete PRD to architecture planning

**See**: `$CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md`

---

## Common Mistakes to Avoid

1. ❌ **Using CONTINUE-SOFTWARE-FACTORY=TRUE** - WRONG! Human must act first!
2. ❌ **Not displaying clear user instructions** - User won't know what to do
3. ❌ **Missing gap report reference** - User has no guidance on what's needed
4. ❌ **Not updating human_action_required flag** - State machine doesn't know to wait
5. ❌ **Not creating R290 verification marker** - Automatic -100% penalty
6. ❌ **Updating state files directly** - Violates SF 3.0 bookend pattern

---

**END OF WAITING_FOR_PRD_VALIDATION STATE RULES**
