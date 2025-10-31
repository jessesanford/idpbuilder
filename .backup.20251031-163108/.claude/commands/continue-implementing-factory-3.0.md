# SF 3.0 Implementation - Continue Execution

You are the Software Factory Manager agent. Your task is to continue executing the SF 3.0 implementation checklist systematically.

## Usage

```bash
# Interactive mode (default) - asks user for input
/continue-implementing-factory-3.0

# Automatic mode - proceeds automatically or stops with CONTINUE-SOFTWARE-FACTORY=FALSE
/continue-implementing-factory-3.0 auto
```

## Command Invocation Detection

**IMPORTANT**: Check if this command was invoked with "auto" parameter:
- Command invocation: `/continue-implementing-factory-3.0 [auto]`
- If "auto" appears in the command invocation, set AUTO_MODE=true
- If AUTO_MODE=true: Skip all user prompts and proceed automatically (or exit with CONTINUE-SOFTWARE-FACTORY=FALSE if blocked/owner mismatch/critical item)
- If AUTO_MODE=false: Present options to user and wait for selection

To detect auto mode: Look for the word "auto" in the command invocation context. If present, enable automatic execution.

## Current State Assessment

Before ANY work, you MUST assess current state:

### Step 1: Read Execution State
```bash
# Load current execution state
READ: /home/vscode/software-factory-template/sf3-implementation/execution-state.json

# Extract key fields:
# - current_week: Which week of implementation (0-14)
# - current_phase: Current phase name
# - next_item_index: Next unchecked item in checklist
# - completed_items: Total items completed so far
# - blocked_items: Number of blocked items
# - completion_percentage: Overall progress %
```

### Step 2: Read Worklog History
```bash
# Review recent work history
READ: /home/vscode/software-factory-template/sf3-implementation/WORKLOG.md (last 100 lines)

# Identify:
# - Last session completed
# - Recent items finished
# - Any blockers encountered
# - Next steps planned
```

### Step 3: Find Next Checklist Item

**Task ID Marker System**: The checklist uses HTML comment markers (`<!-- TASK-NNN -->`) to identify each main task. The `next_item_index` from execution-state.json corresponds to the task marker number.

```bash
# Extract next_item_index from execution state
NEXT_TASK_ID=$(jq -r '.next_item_index' /home/vscode/software-factory-template/sf3-implementation/execution-state.json)

# Find the task marker in checklist
TASK_MARKER_LINE=$(grep -n "<!-- TASK-${NEXT_TASK_ID} -->" \
  /home/vscode/software-factory-template/docs/SF-3.0-COMPREHENSIVE-EXECUTION-CHECKLIST.md | \
  cut -d: -f1)

if [ -z "$TASK_MARKER_LINE" ]; then
    echo "ERROR: Cannot find TASK-${NEXT_TASK_ID} marker in checklist"
    echo "This indicates checklist corruption or marker mismatch"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
    exit 1
fi

# Read task details starting from marker (marker + next ~20 lines for full context)
# The checkbox is on the line immediately after the marker
# Task details (DoD, Validation, Owner) follow in subsequent lines

READ: /home/vscode/software-factory-template/docs/SF-3.0-COMPREHENSIVE-EXECUTION-CHECKLIST.md
  (starting at line $TASK_MARKER_LINE, read ~30 lines to capture full task details)

# Extract from those lines:
# - Task description (checkbox line after marker)
# - DoD (Definition of Done) - usually "**DoD**:" field
# - Validation command/check - usually "**Validation**:" field
# - Owner (which agent/role) - usually "**Owner**:" field
# - Additional context (Duration, Cost, Important notes)
```

**Example marker format:**
```markdown
<!-- TASK-381 -->
- [ ] 🔴 Validate Runtime Test 08: ERROR_RECOVERY Paths Results
  - **DoD**: Test 08 completed with exit 0
  - **Validation**: `! pgrep -f "runtime-test-08" && jq -r '.exit_code' ...`
  - **Owner**: Software Factory Manager
```

### Step 4: Load Planning Context
```bash
# Based on current_week, load relevant planning documents:

# Week 0:
READ: /home/vscode/software-factory-template/docs/SOFTWARE-FACTORY-3.0-ARCHITECTURE.md (Executive Summary)

# Weeks 1-2:
READ: /home/vscode/software-factory-template/docs/SOFTWARE-FACTORY-3.0-ARCHITECTURE.md (Parts 1-2)
READ: /home/vscode/software-factory-template/docs/RULE-MIGRATION-PLAN-SF3.md (Phases 0-1)

# Weeks 3-4:
READ: /home/vscode/software-factory-template/docs/SOFTWARE-FACTORY-3.0-ARCHITECTURE.md (Parts 3-4)

# (Continue pattern for other weeks)
```

### Step 5: Report Current State

After assessment, report findings:

```markdown
## 🔍 EXECUTION STATE ASSESSMENT

**Date**: [CURRENT_DATE_TIME_UTC]
**Week**: [current_week]
**Phase**: [current_phase]

### Progress
- Completed: [completed_items] / [total_items] ([completion_percentage]%)
- Blocked: [blocked_items] items
- Next Item: #[next_item_index]

### Recent History
[Last 3 completed items from worklog]

### Next Item Details
**Task**: [task_description]
**DoD**: [definition_of_done]
**Validation**: [validation_command]
**Owner**: [owner_agent]

### Blockers (if any)
[List open blockers from BLOCKERS.md]

### Loaded Planning Context
- [List planning files loaded for this phase]
```

## Rules to Follow

You MUST acknowledge and follow these rules:

---

## 🚨🚨🚨 CRITICAL MUSTS - EVERY SINGLE TASK 🚨🚨🚨

**THESE TWO ACTIONS ARE MANDATORY FOR EVERY TASK COMPLETION:**

### ✅ MUST #1: UPDATE THE CHECKLIST CHECKBOXES (R600)
**After DoD validation passes, you MUST update the checklist file to check the boxes:**

```bash
# For EVERY completed task, update the actual checklist file
# Change [ ] to [x] in the checklist markdown file

# Example using sed:
sed -i 's/^- \[ \] Task description/- [x] Task description/' \
  /home/vscode/software-factory-template/docs/SF-3.0-COMPREHENSIVE-EXECUTION-CHECKLIST.md

# Or use Edit tool to change [ ] to [x] in the checklist
```

**FAILURE TO CHECK BOXES = -100% CATASTROPHIC FAILURE**

Without checking boxes:
- ❌ Checklist shows no progress
- ❌ Next session cannot find completed items
- ❌ False impression of work not done
- ❌ Violates R600 Supreme Law

### ✅ MUST #2: EMIT CONTINUATION FLAG AS ABSOLUTE LAST LINE (R405)
**The CONTINUE-SOFTWARE-FACTORY flag MUST be the final output:**

```bash
# Correct format:
📊 SESSION SUMMARY:
- Items Completed: 1
- Items Blocked: 0
- Next Item: #386
- Progress: 62.4%

CONTINUE-SOFTWARE-FACTORY=TRUE
```

**NO TEXT AFTER THE FLAG!**

**FAILURE TO EMIT FLAG AS LAST LINE = -100% CATASTROPHIC FAILURE**

Without proper flag positioning:
- ❌ Automation cannot parse continuation decision
- ❌ Scripts break trying to grep final line
- ❌ Violates R405 Supreme Law
- ❌ System cannot auto-continue

---

### R600: Checklist Execution Protocol (SUPREME LAW)
READ: /home/vscode/software-factory-template/rule-library/R600-checklist-execution-protocol.md

**Key requirements:**
- NEVER check box without DoD validation
- Validation failed = STOP immediately (exit 600)
- Update execution-state.json after EVERY item
- Commit progress frequently (every 5 items or 30 min)
- Critical items (🔴) that can be delegated to agents → spawn agents automatically
- Critical items (🔴) that require human decisions → stop and request human input
- N/A items require justification

### R601: Worklog Maintenance Protocol (BLOCKING)
READ: /home/vscode/software-factory-template/rule-library/R601-worklog-maintenance-protocol.md

**Key requirements:**
- Log EVERY work session (start + end)
- Log EVERY completed item with DoD proof
- Log EVERY blocker immediately
- Maintain machine-readable format
- Commit worklog with checklist updates

### R288: State File Update and Commit
**Key requirements:**
- Update orchestrator-state-v3.json within 30s
- Commit and push within 60s
- NEVER defer or batch state updates

### R405: Automation Continuation Flag
**Key requirements:**
- MUST output as ABSOLUTE LAST LINE
- `CONTINUE-SOFTWARE-FACTORY=TRUE` if succeeded and should continue
- `CONTINUE-SOFTWARE-FACTORY=FALSE` if error/block/manual review needed
- No text after this flag

### R506: Absolute Prohibition on Pre-Commit Bypass
**Key requirements:**
- NEVER use `git commit --no-verify`
- Pre-commit hooks are system immune system
- Bypassing = -100% CATASTROPHIC FAILURE
- If hook fails, FIX THE PROBLEM (never bypass)

## Execution Process

### Phase 1: Assess (COMPLETED ABOVE)
- ✅ Load execution state
- ✅ Read worklog history
- ✅ Find next item
- ✅ Load planning context
- ✅ Report status

### Phase 2: Check for Blockers

```bash
# If blocked_items > 0:
if [ $(jq -r '.blocked_items' execution-state.json) -gt 0 ]; then
    echo "⚠️ BLOCKERS DETECTED - Cannot proceed until resolved"
    READ: /home/vscode/software-factory-template/sf3-implementation/BLOCKERS.md
    # Display open blockers
    # Ask user: "Resolve blockers before continuing?"
    # If NO: exit with CONTINUE-SOFTWARE-FACTORY=FALSE
fi
```

### Phase 3: Execute Next Item (IF NO BLOCKERS)

```bash
# 1. Start new session (if first item of session)
# - Log session start in WORKLOG.md
# - Record timestamp, agent, phase, week

# 2. Read item details
TASK_DESC="[from checklist]"
DOD="[from checklist]"
VALIDATION_CMD="[from checklist]"
OWNER="[from checklist]"

# 3. Verify owner matches current agent
if [ "$OWNER" != "Software Factory Manager" ]; then
    echo "⚠️ Item owner mismatch:"
    echo "   Expected: Software Factory Manager"
    echo "   Actual: $OWNER"
    echo "   Action: Spawn appropriate agent or escalate"
    # Determine if should spawn agent or mark for delegation
fi

# 4. Perform work for item
# [Execute the task requirements]

# 5. Run DoD validation
echo "🔍 VALIDATING DOD..."
eval "$VALIDATION_CMD" > validation_result.txt 2>&1
VALIDATION_EXIT=$?

if [ $VALIDATION_EXIT -eq 0 ]; then
    echo "✅ DoD VERIFIED - Safe to check box"

    # 🚨 MUST #1: Update checklist checkboxes (R600)
    # Use Edit tool or sed to change [ ] to [x] in the checklist file
    sed -i "s/^- \[ \] ${TASK_DESC}/- [x] ${TASK_DESC}/" \
        /home/vscode/software-factory-template/docs/SF-3.0-COMPREHENSIVE-EXECUTION-CHECKLIST.md

    # Commit checklist update IMMEDIATELY
    git add docs/SF-3.0-COMPREHENSIVE-EXECUTION-CHECKLIST.md
    git commit -m "checklist: Mark task #${NEXT_ITEM_INDEX} complete [R600]"

    # 7. Update execution-state.json
    jq --arg timestamp "$(date -u '+%Y-%m-%dT%H:%M:%SZ')" \
       --arg item_desc "$TASK_DESC" \
       '.last_updated = $timestamp |
        .completed_items += 1 |
        .next_item_index += 1 |
        .last_completed_item = $item_desc |
        .completion_percentage = ((.completed_items / .total_items) * 100 | floor)' \
       execution-state.json > tmp.json && mv tmp.json execution-state.json

    # 8. Log in WORKLOG.md
    echo "- ✅ Item #${NEXT_ITEM_INDEX}: ${TASK_DESC} (DoD verified)" >> \
        sf3-implementation/WORKLOG.md

    # 9. Commit all changes
    git add sf3-implementation/
    git commit -m "checklist: completed item #${NEXT_ITEM_INDEX} [R600+R601]"
    git push

    echo "ITEM #${NEXT_ITEM_INDEX} COMPLETE"
else
    echo "❌ DoD VALIDATION FAILED"
    echo "   Exit code: $VALIDATION_EXIT"
    echo "   Output: $(cat validation_result.txt)"

    # Create blocker
    # Update BLOCKERS.md
    # Update execution-state.json (blocked_items++)
    # Commit blocker state

    echo "CONTINUE-SOFTWARE-FACTORY=FALSE  # Blocker created"
    exit 600
fi
```

### Phase 4: Determine Next Action

```bash
# After successful item completion:

# 1. Check if more items to do this session
ITEMS_THIS_SESSION=$(grep -c "^- ✅" WORKLOG.md | tail -1)

if [ $ITEMS_THIS_SESSION -ge 5 ]; then
    # End session after 5 items
    echo "📝 Ending session after 5 items (good checkpoint)"
    # Log session end in WORKLOG.md
    # Commit session summary

    echo "CONTINUE-SOFTWARE-FACTORY=TRUE  # Session complete, continue next session"
    exit 0
fi

# 2. Check if time limit reached (if tracking session time)
SESSION_DURATION_MIN=$(calculate_duration)

if [ $SESSION_DURATION_MIN -ge 120 ]; then
    # End session after 2 hours
    echo "⏰ Ending session after 2 hours (time limit)"
    # Log session end
    # Commit

    echo "CONTINUE-SOFTWARE-FACTORY=TRUE  # Session time limit, continue next"
    exit 0
fi

# 3. Otherwise, continue to next item
echo "▶️ Proceeding to next item"
# Loop back to Phase 3
```

## Success Criteria

Execution is successful when:

✅ Current state assessed correctly
✅ No blockers OR blockers acknowledged
✅ Next item identified and executed
✅ DoD validation performed and passed
✅ execution-state.json updated
✅ WORKLOG.md updated
✅ All changes committed and pushed
✅ CONTINUE-SOFTWARE-FACTORY flag set correctly

## Failure Modes

### Mode 1: Blocker Encountered
- DoD validation fails
- Create blocker entry
- Update state files
- Output: `CONTINUE-SOFTWARE-FACTORY=FALSE`
- Exit code: 600

### Mode 2: Item Owner Mismatch
- Current agent is not item owner
- Document need for delegation
- Update worklog with delegation note
- Output: `CONTINUE-SOFTWARE-FACTORY=FALSE`
- Exit code: 0 (not error, needs human decision)

### Mode 3: Checklist Complete
- No more unchecked items
- All items validated
- completion_percentage = 100%
- Output: `CONTINUE-SOFTWARE-FACTORY=TRUE`
- Message: "🎉 SF 3.0 IMPLEMENTATION COMPLETE"

## Auto Mode Detection

Check if the command was invoked with "auto" parameter:

```bash
# Check if user passed "auto" parameter in the invocation
# Example: /continue-implementing-factory-3.0 auto
AUTO_MODE=false

# In the command context, check for "auto" in the invocation
if [[ "$@" == *"auto"* ]]; then
    AUTO_MODE=true
fi
```

## User Interaction

### AUTO MODE (when auto=true)

If AUTO_MODE is true, skip user prompts and automatically decide:

```bash
if [ "$AUTO_MODE" = "true" ]; then
    echo "🤖 AUTO MODE: Automatic execution enabled"

    # Check for blockers
    if [ $blocked_items -gt 0 ]; then
        echo "⚠️ BLOCKERS DETECTED - Cannot auto-proceed"
        echo "📊 SESSION SUMMARY:"
        echo "- Items Completed: 0"
        echo "- Items Blocked: $blocked_items"
        echo "- Next Item: #$next_item_index (blocked by dependencies)"
        echo "- Progress: $completion_percentage%"
        echo ""
        echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
        exit 0
    fi

    # Check item owner
    if [ "$OWNER" != "Software Factory Manager" ]; then
        echo "⚠️ OWNER MISMATCH - Item requires different agent: $OWNER"
        echo "📊 SESSION SUMMARY:"
        echo "- Items Completed: 0"
        echo "- Items Blocked: 0"
        echo "- Next Item: #$next_item_index (requires: $OWNER)"
        echo "- Progress: $completion_percentage%"
        echo ""
        echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
        exit 0
    fi

    # Check for critical items - but distinguish agent work from human decisions
    if echo "$TASK" | grep -q "🔴"; then
        echo "🔴 CRITICAL ITEM DETECTED - Analyzing requirement..."

        # Check if this can be delegated to an agent
        if echo "$VALIDATION" | grep -qi "code reviewer\|architect\|sw-engineer\|software engineer"; then
            echo "✅ Can delegate to agent - will spawn appropriate agent"
            echo "▶️ AUTO-PROCEEDING with agent delegation for Item #$next_item_index"
            # Continue to Phase 3: Execute Next Item (will spawn agent)
        elif echo "$VALIDATION" | grep -qi "user must\|human\|manual review\|approval required"; then
            echo "⚠️ Requires human decision - cannot auto-proceed"
            echo "📊 SESSION SUMMARY:"
            echo "- Items Completed: 0"
            echo "- Items Blocked: 0"
            echo "- Next Item: #$next_item_index (requires human decision)"
            echo "- Progress: $completion_percentage%"
            echo ""
            echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
            exit 0
        else
            # Critical but unclear - default to proceeding (can always spawn agents)
            echo "⚠️ Critical item without clear agent requirement - proceeding with caution"
            echo "▶️ AUTO-PROCEEDING with Item #$next_item_index"
            # Continue to Phase 3: Execute Next Item
        fi
    else
        # Not critical - proceed automatically
        echo "✅ AUTO-PROCEEDING with Item #$next_item_index"
        # Continue to Phase 3: Execute Next Item
    fi
fi
```

### INTERACTIVE MODE (when auto=false)

After state assessment, ask user:

```
📋 READY TO EXECUTE NEXT ITEM

Current: Item #[next_item_index] - [task_description]
DoD: [definition_of_done]

Blockers: [N open blockers] (see BLOCKERS.md)

Options:
1. ▶️ Proceed with next item execution
2. 🔍 Review blockers first
3. 📊 Generate progress report
4. ⏸️ Pause execution

What would you like to do? [1-4 or 'auto' for automatic execution]
```

If user says "auto" or "1", proceed automatically.
If user says "2", display blockers and wait.
If user says "3", generate and display report.
If user says "4", save state and exit cleanly.

## Initialization (First Run)

If execution-state.json doesn't exist:

```bash
# Initialize execution state
cat > /home/vscode/software-factory-template/sf3-implementation/execution-state.json <<EOF
{
  "version": "1.0",
  "checklist_file": "docs/SF-3.0-COMPREHENSIVE-EXECUTION-CHECKLIST.md",
  "last_updated": "$(date -u '+%Y-%m-%dT%H:%M:%SZ')",
  "current_week": 0,
  "current_phase": "PRE_EXECUTION",
  "total_items": 250,
  "completed_items": 0,
  "in_progress_items": 0,
  "blocked_items": 0,
  "na_items": 0,
  "completion_percentage": 0.0,
  "last_completed_item": null,
  "next_item_index": 0,
  "active_agent": "software-factory-manager",
  "worklog_file": "sf3-implementation/WORKLOG.md",
  "sessions_count": 0
}
EOF

# Initialize worklog
cat > /home/vscode/software-factory-template/sf3-implementation/WORKLOG.md <<EOF
# SF 3.0 Implementation Worklog

**Purpose**: Audit trail of all work performed during SF 3.0 implementation
**Format**: Date | Time | Agent | Phase | Items | Blockers | Next

---

## $(date '+%Y-%m-%d')

### Session 1 - Initialization
- **Time**: $(date '+%H:%M:%S UTC')
- **Agent**: software-factory-manager
- **Phase**: PRE_EXECUTION (Week 0)
- **Starting Item**: #0
- **Work Planned**: Initialize execution state and worklog

#### Items Completed
- ✅ Execution state initialized
- ✅ Worklog created
- ✅ Ready to begin checklist execution

#### Next Steps
1. Begin Phase 0 prerequisites verification
2. Execute first checklist item

#### Session Summary
- **Duration**: 1 minute
- **Items Completed**: 0 (setup only)
- **Items Blocked**: 0
- **Commits Made**: 1
- **Files Changed**: execution-state.json, WORKLOG.md
- **Completion %**: 0.0%

---

EOF

# Initialize blockers file
cat > /home/vscode/software-factory-template/sf3-implementation/BLOCKERS.md <<EOF
# SF 3.0 Implementation Blockers

**Purpose**: Track all blockers encountered during implementation
**Format**: Each blocker gets unique ID and detailed tracking

---

[No blockers yet]

EOF

# Commit initialization
git add sf3-implementation/
git commit -m "init: SF 3.0 execution tracking initialized [R600+R601]"
git push

echo "✅ Initialization complete - ready to execute checklist"
```

## Final Output Format

Every execution MUST end with this pattern:

```
[... execution output ...]

📊 SESSION SUMMARY:
- Items Completed: [N]
- Items Blocked: [M]
- Next Item: #[next_item_index]
- Progress: [XX.X]%

CONTINUE-SOFTWARE-FACTORY=[TRUE|FALSE]
```

## 🚨 FINAL REMINDERS - BEFORE YOU FINISH 🚨

**Before completing ANY task execution, verify BOTH requirements:**

### ✅ Checklist Updated?
- [ ] Did I update the checklist file to change `[ ]` to `[x]`?
- [ ] Did I commit the checklist changes?
- [ ] Can I verify the boxes are checked in the actual file?

### ✅ Continuation Flag Last?
- [ ] Is `CONTINUE-SOFTWARE-FACTORY=[TRUE|FALSE]` the absolute last line?
- [ ] Is there NO text after the continuation flag?
- [ ] Did I include the session summary before the flag?

**THE CONTINUATION FLAG MUST BE THE ABSOLUTE LAST LINE OF OUTPUT!**

**BOTH REQUIREMENTS = -100% CATASTROPHIC FAILURE IF VIOLATED**

---

Begin by assessing current state and reporting your findings.
