# 🚨🚨🚨 RULE R601 - WORKLOG MAINTENANCE PROTOCOL [BLOCKING]

**Status**: ACTIVE
**Criticality**: BLOCKING
**Enforcement**: Exit code 601, -50% for missing worklogs, -100% for incomplete audit trail

## Purpose

Mandates comprehensive worklog maintenance for ALL SF 3.0 implementation work, providing complete audit trail of execution, blockers, and progress.

## Blocking Status

This rule is BLOCKING because:
- Worklog is the ONLY reliable record of what happened
- Missing worklog = no accountability
- Incomplete worklog = lost knowledge on interruption
- No worklog = cannot resume work correctly

**Missing worklog entries = LOST WORK HISTORY = -50% minimum penalty**

## Worklog File

**Location**: `/home/vscode/software-factory-template/sf3-implementation/WORKLOG.md`

**Format**: Markdown with structured sections

**Created**: During initial setup (execution-state.json initialization)

## Worklog Format Specification

### Session Header (MANDATORY)

```markdown
## YYYY-MM-DD

### Session N - [SESSION_TYPE]
- **Time**: HH:MM:SS UTC
- **Agent**: [agent-name]
- **Phase**: [current_phase from execution-state.json]
- **Week**: [current_week from execution-state.json]
- **Starting Item**: #[next_item_index]
- **Work Planned**: [Brief description of session goals]
```

**SESSION_TYPE values:**
- `Initialization`: First session, setup work
- `Execution`: Normal checklist execution
- `Blocker Investigation`: Investigating/resolving blockers
- `Validation`: Peer validation of critical items
- `Recovery`: Resuming after interruption

### Item Completions (MANDATORY)

```markdown
#### Items Completed
- ✅ Item #123: [task description] (DoD verified - [validation command] passed)
- ✅ Item #124: [task description] (DoD verified - [validation command] passed)
- ⊘ Item #125: [task description] (N/A - [justification])
```

**Format rules:**
- ✅ for completed items with DoD validation
- ⊘ for N/A items with justification
- NEVER use ✅ without "(DoD verified - ...)" proof

### Items Blocked (IF ANY)

```markdown
#### Items Blocked
- ❌ Item #126: [task description]
  - **Blocker ID**: BLOCKER-YYYYMMDD-HHMMSS
  - **Reason**: [Root cause description]
  - **Impact**: [What this blocks]
  - **ETA**: [Expected resolution time or "UNKNOWN"]
```

### Next Steps (MANDATORY)

```markdown
#### Next Steps
1. [Specific next action for this work stream]
2. [Dependency or prerequisite needed]
3. [Escalation if blocked]

**Next Item**: #[next_item_index + 1] - [task description]
```

### Session Summary (MANDATORY)

```markdown
#### Session Summary
- **Duration**: [X] minutes
- **Items Completed**: [N]
- **Items Blocked**: [M]
- **Commits Made**: [K]
- **Files Changed**: [List key files]
- **Completion %**: [X.X]% ([completed]/[total])
```

## Complete Session Template

```bash
# Generate session worklog entry
cat >> sf3-implementation/WORKLOG.md <<EOF

## $(date '+%Y-%m-%d')

### Session $(jq -r '.sessions_count // 0' execution-state.json) - Execution
- **Time**: $(date '+%H:%M:%S UTC')
- **Agent**: ${AGENT_NAME}
- **Phase**: $(jq -r '.current_phase' execution-state.json)
- **Week**: $(jq -r '.current_week' execution-state.json)
- **Starting Item**: #$(jq -r '.next_item_index' execution-state.json)
- **Work Planned**: Execute items #$(jq -r '.next_item_index' execution-state.json) through #$(($(jq -r '.next_item_index' execution-state.json) + 5))

#### Items Completed
[POPULATED DURING SESSION]

#### Items Blocked
[IF ANY - OTHERWISE "None"]

#### Next Steps
[POPULATED AT SESSION END]

#### Session Summary
- **Duration**: [CALCULATED] minutes
- **Items Completed**: [COUNTED]
- **Items Blocked**: [COUNTED]
- **Commits Made**: [COUNTED]
- **Files Changed**: [LISTED]
- **Completion %**: [CALCULATED]%

EOF
```

## Session Start Protocol

**At the BEGINNING of EVERY work session:**

```bash
# 1. Calculate session number
SESSION_NUM=$(grep -c "^### Session" sf3-implementation/WORKLOG.md)
SESSION_NUM=$((SESSION_NUM + 1))

# 2. Record start time
SESSION_START=$(date +%s)

# 3. Get current state
CURRENT_PHASE=$(jq -r '.current_phase' sf3-implementation/execution-state.json)
CURRENT_WEEK=$(jq -r '.current_week' sf3-implementation/execution-state.json)
NEXT_ITEM=$(jq -r '.next_item_index' sf3-implementation/execution-state.json)

# 4. Write session header
cat >> sf3-implementation/WORKLOG.md <<EOF

## $(date '+%Y-%m-%d')

### Session $SESSION_NUM - Execution
- **Time**: $(date '+%H:%M:%S UTC')
- **Agent**: software-factory-manager
- **Phase**: $CURRENT_PHASE
- **Week**: $CURRENT_WEEK
- **Starting Item**: #$NEXT_ITEM
- **Work Planned**: Continue SF 3.0 implementation checklist

#### Items Completed
EOF

# 5. Commit session start
git add sf3-implementation/WORKLOG.md
git commit -m "worklog: session $SESSION_NUM started [R601]"
git push

# 6. Update execution-state.json
jq --arg session_num "$SESSION_NUM" \
   '.sessions_count = ($session_num | tonumber) |
    .current_session_start = now' \
   sf3-implementation/execution-state.json > tmp.json && mv tmp.json sf3-implementation/execution-state.json
```

**Session start MUST be logged before ANY work.**

## Item Completion Logging

**After EVERY item completion:**

```bash
# 1. Extract item details
ITEM_NUM=$next_item_index
ITEM_TASK="[Task description from checklist]"
DOD_VALIDATION="[Validation command that was run]"

# 2. Append to worklog (under "Items Completed" section)
# Find the last "#### Items Completed" line and append after it
cat >> sf3-implementation/WORKLOG.md <<EOF
- ✅ Item #${ITEM_NUM}: ${ITEM_TASK} (DoD verified - ${DOD_VALIDATION} passed)
EOF

# 3. Commit item completion
git add sf3-implementation/WORKLOG.md
git commit -m "worklog: item #${ITEM_NUM} completed [R601]"
git push
```

**MANDATORY log for every completion - real-time, not batched.**

## Blocker Logging

**When item is blocked:**

```bash
# 1. Generate blocker ID
BLOCKER_ID="BLOCKER-$(date '+%Y%m%d-%H%M%S')"

# 2. Document blocker in worklog
cat >> sf3-implementation/WORKLOG.md <<EOF

#### Items Blocked
- ❌ Item #${ITEM_NUM}: ${ITEM_TASK}
  - **Blocker ID**: $BLOCKER_ID
  - **Reason**: [Root cause - extracted from validation failure]
  - **Impact**: Prevents items #${ITEM_NUM} through #[DEPENDENT_ITEMS]
  - **ETA**: [Investigation required / Known fix available / Escalated]
EOF

# 3. Also create detailed blocker file
cat >> sf3-implementation/BLOCKERS.md <<EOF

## $BLOCKER_ID

**Date**: $(date '+%Y-%m-%d %H:%M UTC')
**Item**: #${ITEM_NUM} - ${ITEM_TASK}
**DoD**: ${DOD_REQUIREMENT}
**Validation Failed**: ${DOD_VALIDATION}

### Failure Output
\`\`\`
[Actual command output that showed failure]
\`\`\`

### Root Cause Analysis
[TO BE COMPLETED]

### Resolution Plan
[TO BE COMPLETED]

### Status
- [ ] Root cause identified
- [ ] Resolution plan created
- [ ] Fix implemented
- [ ] Validation passes
- [ ] Item unblocked

**Resolution Date**: [PENDING]
**Resolved By**: [PENDING]

EOF

# 4. Commit blocker documentation
git add sf3-implementation/WORKLOG.md sf3-implementation/BLOCKERS.md
git commit -m "worklog: item #${ITEM_NUM} blocked - $BLOCKER_ID [R601]"
git push
```

**Blockers MUST be logged immediately - no delays.**

## Session End Protocol

**At the END of EVERY work session:**

```bash
# 1. Calculate session duration
SESSION_END=$(date +%s)
DURATION=$(( (SESSION_END - SESSION_START) / 60 ))

# 2. Count session accomplishments
ITEMS_COMPLETED=$(grep -c "^- ✅" sf3-implementation/WORKLOG.md | tail -1)
ITEMS_BLOCKED=$(grep -c "^- ❌" sf3-implementation/WORKLOG.md | tail -1)
COMMITS_MADE=$(git rev-list --count HEAD --since="@${SESSION_START}")

# 3. Get updated completion percentage
COMPLETION_PCT=$(jq -r '.completion_percentage' sf3-implementation/execution-state.json)

# 4. Determine next steps
NEXT_ITEM_IDX=$(jq -r '.next_item_index' sf3-implementation/execution-state.json)
NEXT_TASK=$(sed -n "${NEXT_ITEM_IDX}p" SF-3.0-COMPREHENSIVE-EXECUTION-CHECKLIST.md)

# 5. Append session summary to worklog
cat >> sf3-implementation/WORKLOG.md <<EOF

#### Next Steps
1. Continue with Item #${NEXT_ITEM_IDX}: ${NEXT_TASK}
2. Resolve any blockers identified in this session
3. Maintain commit frequency per R600

**Next Item**: #${NEXT_ITEM_IDX}

#### Session Summary
- **Duration**: ${DURATION} minutes
- **Items Completed**: ${ITEMS_COMPLETED}
- **Items Blocked**: ${ITEMS_BLOCKED}
- **Commits Made**: ${COMMITS_MADE}
- **Files Changed**: [LIST FROM GIT]
- **Completion %**: ${COMPLETION_PCT}%

---

EOF

# 6. Commit session end
git add sf3-implementation/WORKLOG.md
git commit -m "worklog: session $SESSION_NUM complete - ${ITEMS_COMPLETED} items [R601]"
git push
```

**Session end MUST be logged even if interrupted.**

## Commit Protocol

**Worklog commits ALWAYS accompany checklist commits:**

```bash
# Combined commit pattern (from R600)
git add SF-3.0-COMPREHENSIVE-EXECUTION-CHECKLIST.md
git add sf3-implementation/execution-state.json
git add sf3-implementation/WORKLOG.md

git commit -m "checklist: completed item #${ITEM_NUM} - ${TASK_DESC}
worklog: session $SESSION_NUM item completion
[R600+R601]"

git push
```

**Worklog is NEVER committed separately - always with state updates.**

## Progress Reporting Format

**Daily/Weekly Summary (OPTIONAL but recommended):**

```markdown
## Week N Summary - YYYY-MM-DD

### Achievements
- ✅ [X] items completed this week
- ✅ [Key milestone reached]
- ✅ [Phase/week completion]

### Challenges
- ❌ [N] items blocked
- ⚠️ [Challenge description]

### Metrics
- **Weekly Completion Rate**: [X] items/day
- **Total Progress**: [XX]%
- **Blockers Resolved**: [N]
- **Average Session Duration**: [X] minutes

### Outlook
- **Next Week Goal**: Complete items #[START] through #[END]
- **Expected Challenges**: [Describe]
- **Dependencies**: [List]
```

## Machine-Readable Format

**Worklog MUST be parseable:**

```bash
# Extract session count
grep -c "^### Session" sf3-implementation/WORKLOG.md

# Extract completed items
grep "^- ✅" sf3-implementation/WORKLOG.md | wc -l

# Extract blocked items
grep "^- ❌" sf3-implementation/WORKLOG.md | wc -l

# Extract session durations
grep "Duration:" sf3-implementation/WORKLOG.md | awk '{print $3}'

# Calculate average session duration
grep "Duration:" sf3-implementation/WORKLOG.md | awk '{sum+=$3; count++} END {print sum/count}'

# Get current completion rate
tail -20 sf3-implementation/WORKLOG.md | grep "Completion %" | tail -1
```

**Structured format enables automated analysis.**

## Recovery from Interruption

**Using worklog to resume:**

```bash
# 1. Read last session
LAST_SESSION=$(grep -A 20 "^### Session" sf3-implementation/WORKLOG.md | tail -20)

echo "📖 LAST SESSION SUMMARY:"
echo "$LAST_SESSION"

# 2. Check for blockers
OPEN_BLOCKERS=$(grep "^- ❌" sf3-implementation/BLOCKERS.md | grep -v "RESOLVED")

if [ -n "$OPEN_BLOCKERS" ]; then
    echo ""
    echo "⚠️  OPEN BLOCKERS FROM PREVIOUS SESSION:"
    echo "$OPEN_BLOCKERS"
    echo ""
    echo "Resolve blockers before continuing."
fi

# 3. Display next steps from last session
NEXT_STEPS=$(grep -A 5 "^#### Next Steps" sf3-implementation/WORKLOG.md | tail -5)

echo ""
echo "🎯 NEXT STEPS (from last session):"
echo "$NEXT_STEPS"

# 4. Ready to resume
echo ""
echo "✅ Ready to resume execution"
```

**Worklog provides complete context for resumption.**

## Enforcement & Penalties

### CATASTROPHIC FAILURES (-100%)
- No worklog for entire work stream
- Claiming work done without worklog proof
- Checking items without worklog entry

### CRITICAL FAILURES (-50% to -75%)
- Missing session headers
- Incomplete session summaries
- Blockers not logged in worklog

### STANDARD FAILURES (-20% to -40%)
- Missing "Next Steps" sections
- Inaccurate duration reporting
- Missing blocker details
- Session not committed

### MINOR VIOLATIONS (-10% to -15%)
- Formatting inconsistencies
- Missing file change lists
- Incomplete metrics

## Success Criteria

Worklog is complete when:

✅ EVERY work session has entry
✅ EVERY completed item logged
✅ EVERY blocker documented
✅ ALL sessions have start/end
✅ Machine-readable format maintained
✅ Complete audit trail from start to finish

## Related Rules

- **R600**: Checklist Execution Protocol (worklog is part of execution)
- **R288**: State File Update and Commit (worklog commits)
- **R405**: Automation Continuation Flag (session completion)

## File References

- **Worklog**: `/home/vscode/software-factory-template/sf3-implementation/WORKLOG.md`
- **Blockers**: `/home/vscode/software-factory-template/sf3-implementation/BLOCKERS.md`
- **Execution State**: `/home/vscode/software-factory-template/sf3-implementation/execution-state.json`
- **Checklist**: `/home/vscode/software-factory-template/docs/SF-3.0-COMPREHENSIVE-EXECUTION-CHECKLIST.md`

---

**REMEMBER**: The worklog is your memory. Without it, you have no proof of what was done, no record of blockers, and no way to resume. Maintain it religiously.
