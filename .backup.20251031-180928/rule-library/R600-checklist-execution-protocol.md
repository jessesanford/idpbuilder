# 🔴🔴🔴 RULE R600 - CHECKLIST EXECUTION PROTOCOL [SUPREME LAW]

**Status**: ACTIVE
**Criticality**: SUPREME LAW - BLOCKING
**Enforcement**: Exit code 600, -100% for unchecked items without DoD validation

## Purpose

Defines mandatory protocol for executing `/home/vscode/software-factory-template/docs/SF-3.0-COMPREHENSIVE-EXECUTION-CHECKLIST.md` systematically, ensuring EVERY item is verified before checking boxes and ALL work is tracked.

## Supreme Law Status

This rule is a SUPREME LAW because:
- Checklist execution determines SF 3.0 success/failure
- Unchecked items = incomplete implementation
- Checking boxes without DoD validation = FALSE PROGRESS
- Missing items discovered late = cascading failures

**NEVER check a box without completing DoD validation = -100% CATASTROPHIC FAILURE**

## Core Protocol

### 1. Read Checklist Systematically

**ALWAYS read from top to bottom, in order:**

```bash
# Start at current position
READ: /home/vscode/software-factory-template/docs/SF-3.0-COMPREHENSIVE-EXECUTION-CHECKLIST.md

# Find first unchecked item
PATTERN: "- [ ]" (unchecked checkbox)

# Record item details:
- Task description
- DoD (Definition of Done)
- Validation command/check
- Owner (which agent/role)
```

**NEVER skip items** - if not applicable, mark `N/A` with justification in worklog.

### 2. Determine Current Position

**Before ANY work, find where you are:**

```bash
# 1. Check execution state
READ: /home/vscode/software-factory-template/sf3-implementation/execution-state.json

# Extract:
current_week=$(jq -r '.current_week' execution-state.json)
current_phase=$(jq -r '.current_phase' execution-state.json)
next_item_index=$(jq -r '.next_item_index' execution-state.json)

# 2. Find next unchecked item in checklist
grep -n "^- \[ \]" SF-3.0-COMPREHENSIVE-EXECUTION-CHECKLIST.md | head -n 1

# 3. Cross-reference with execution-state.json
# If mismatch, execution-state.json is authoritative

# 4. Report position BEFORE starting work
echo "📍 CURRENT POSITION:"
echo "   Week: $current_week"
echo "   Phase: $current_phase"
echo "   Next Item: #$next_item_index"
```

**Position determination is MANDATORY before ANY work.**

### 3. Verify DoD Before Checking Box

**THE ABSOLUTE RULE: NEVER check box without DoD verification**

```bash
# For EVERY item:

# 1. Read DoD criteria
DoD=$(grep -A 1 "DoD:" checklist | head -1)

# 2. Read validation command
VALIDATION=$(grep -A 1 "Validation:" checklist | head -1)

# 3. Execute validation command
eval "$VALIDATION" > validation_result.txt 2>&1
VALIDATION_EXIT_CODE=$?

# 4. Verify DoD is met
if [ $VALIDATION_EXIT_CODE -eq 0 ]; then
    # DoD MET - safe to check box
    echo "✅ DoD VERIFIED: $DoD"

    # Update checklist (change [ ] to [x])
    sed -i "s/^- \[ \] ${TASK}/- [x] ${TASK}/" SF-3.0-COMPREHENSIVE-EXECUTION-CHECKLIST.md

    # Add timestamp and validator
    echo "   (✅ $(date '+%Y-%m-%d %H:%M UTC') - @$(whoami))"
else
    # DoD NOT MET - STOP
    echo "❌ DoD VALIDATION FAILED: $DoD"
    echo "   Exit code: $VALIDATION_EXIT_CODE"
    echo "   Output: $(cat validation_result.txt)"

    # Update execution-state.json
    jq '.blocked_items += 1' execution-state.json > tmp.json && mv tmp.json execution-state.json

    # STOP IMMEDIATELY
    exit 600  # R600 violation - DoD not met
fi
```

**ENFORCEMENT:**
- Checking box without validation = -100% IMMEDIATE FAILURE
- Validation failed but continued = -100% IMMEDIATE FAILURE
- Wrong validation command used = -50%

### 4. Update execution-state.json After Each Item

**After EVERY completed item:**

```bash
# Update execution state atomically
jq --arg timestamp "$(date -u '+%Y-%m-%dT%H:%M:%SZ')" \
   --arg item_desc "$TASK_DESCRIPTION" \
   '.last_updated = $timestamp |
    .completed_items += 1 |
    .next_item_index += 1 |
    .last_completed_item = $item_desc |
    .completion_percentage = ((.completed_items / .total_items) * 100 | floor)' \
   execution-state.json > tmp.json && mv tmp.json execution-state.json

# Commit immediately (R288 compliance)
git add execution-state.json
git add SF-3.0-COMPREHENSIVE-EXECUTION-CHECKLIST.md
git commit -m "checklist: completed item #${NEXT_ITEM_INDEX} - ${TASK_DESCRIPTION} [R600]"
git push
```

**MANDATORY after EVERY item - no batching, no deferrals.**

### 5. Commit Progress Frequently (R288 Compliance)

**Commit frequency rules:**

```bash
# Trigger commits on ANY of:
# 1. Every 5 items completed
if [ $((completed_items % 5)) -eq 0 ]; then
    COMMIT_NOW=true
fi

# 2. Every 30 minutes
LAST_COMMIT_TIME=$(git log -1 --format=%ct)
NOW=$(date +%s)
if [ $((NOW - LAST_COMMIT_TIME)) -gt 1800 ]; then
    COMMIT_NOW=true
fi

# 3. Before any state transition
if [ "$TRANSITIONING_STATE" = "true" ]; then
    COMMIT_NOW=true
fi

# 4. After critical items (marked with 🔴)
if echo "$TASK" | grep -q "🔴"; then
    COMMIT_NOW=true
fi

if [ "$COMMIT_NOW" = "true" ]; then
    git add -A
    git commit -m "checklist: checkpoint after $completed_items items [R600+R288]"
    git push
fi
```

**NEVER lose progress - commit early, commit often.**

### 6. Escalate/Stop When DoD Cannot Be Met

**When validation fails:**

```bash
# 1. Document blocker immediately
BLOCKER_ID="BLOCKER-$(date '+%Y%m%d-%H%M%S')"

cat >> sf3-implementation/BLOCKERS.md <<EOF

## $BLOCKER_ID - Item #${NEXT_ITEM_INDEX} Blocked

**Date**: $(date '+%Y-%m-%d %H:%M UTC')
**Item**: $TASK_DESCRIPTION
**DoD**: $DOD
**Validation**: $VALIDATION
**Failure Output**:
\`\`\`
$(cat validation_result.txt)
\`\`\`

**Root Cause**: [TO BE INVESTIGATED]
**Resolution Plan**: [PENDING]
**ETA**: [UNKNOWN]
EOF

# 2. Update execution-state.json
jq --arg blocker_id "$BLOCKER_ID" \
   --arg item_desc "$TASK_DESCRIPTION" \
   '.blocked_items += 1 |
    .blockers += [{
        blocker_id: $blocker_id,
        item: $item_desc,
        timestamp: now | strftime("%Y-%m-%dT%H:%M:%SZ"),
        status: "OPEN"
    }]' \
   execution-state.json > tmp.json && mv tmp.json execution-state.json

# 3. Add blocker note to checklist
sed -i "s/^- \[ \] ${TASK}/- [ ] ${TASK} (⚠️ BLOCKED: $BLOCKER_ID - ETA PENDING)/" \
    SF-3.0-COMPREHENSIVE-EXECUTION-CHECKLIST.md

# 4. Commit blocker state
git add -A
git commit -m "checklist: item #${NEXT_ITEM_INDEX} blocked - $BLOCKER_ID [R600]"
git push

# 5. Set continuation flag to FALSE
echo "CONTINUE-SOFTWARE-FACTORY=FALSE  # BLOCKER: $BLOCKER_ID"

# 6. Exit with blocker code
exit 600
```

**NEVER proceed past a blocker - always escalate.**

## Critical Items (🔴) Protocol

**Items marked with 🔴 require validation by another agent (not self-validation):**

**MODERN APPROACH (SF 3.0 - Automated Agent Delegation):**

```bash
# 1. Critical items require peer validation - spawn the appropriate agent
if echo "$TASK" | grep -q "🔴"; then
    echo "🔴 CRITICAL ITEM - Requires peer validation"

    # 2. Determine which agent to spawn based on validation requirement
    if echo "$VALIDATION" | grep -qi "code reviewer"; then
        echo "✅ Spawning Code Reviewer agent for validation"
        spawn_code_reviewer_agent "$TASK" "$DOD" "$VALIDATION"
        # Agent will perform validation and report results
    elif echo "$VALIDATION" | grep -qi "architect"; then
        echo "✅ Spawning Architect agent for validation"
        spawn_architect_agent "$TASK" "$DOD" "$VALIDATION"
    elif echo "$VALIDATION" | grep -qi "software engineer\|sw-engineer"; then
        echo "✅ Spawning Software Engineer agent for validation"
        spawn_sw_engineer_agent "$TASK" "$DOD" "$VALIDATION"
    elif echo "$VALIDATION" | grep -qi "user must\|human\|approval required"; then
        # Requires human decision - cannot automate
        echo "⚠️  REQUIRES HUMAN DECISION - Cannot delegate to agent"
        echo "CONTINUE-SOFTWARE-FACTORY=FALSE  # Human input required"
        exit 0
    else
        # Critical but unclear who validates - default to spawning code reviewer
        echo "⚠️  Critical item without clear validator - spawning Code Reviewer"
        spawn_code_reviewer_agent "$TASK" "$DOD" "$VALIDATION"
    fi

    # 3. After agent completes, check the box if validation passed
    # (Agent will report PASS/FAIL in its output)
fi
```

**LEGACY APPROACH (Manual Validation Requests):**

```bash
# For non-automated environments, create validation request
if echo "$TASK" | grep -q "🔴"; then
    cat > validation-request-${NEXT_ITEM_INDEX}.md <<EOF
# Validation Request - Item #${NEXT_ITEM_INDEX}
**Task**: $TASK_DESCRIPTION
**DoD**: $DOD
**Validation**: $VALIDATION
**Validator Required**: [Code Reviewer / Architect / etc.]
EOF

    echo "CONTINUE-SOFTWARE-FACTORY=FALSE  # Awaiting peer validation"
    exit 0
fi
```

**KEY PRINCIPLE**: Critical items are NEVER self-validated, but CAN be delegated to appropriate agents automatically.

## Worklog Requirement

**EVERY work session MUST be logged:**

```bash
# At start of session
cat >> sf3-implementation/WORKLOG.md <<EOF

### Session $(date '+%Y-%m-%d %H:%M UTC') - ${AGENT_NAME}
- **Phase**: $CURRENT_PHASE
- **Week**: $CURRENT_WEEK
- **Starting Item**: #${NEXT_ITEM_INDEX}
- **Work Planned**: [DESCRIBE]
EOF

# During session (as items complete)
cat >> sf3-implementation/WORKLOG.md <<EOF
- ✅ Item #${ITEM_INDEX}: $TASK_DESCRIPTION (DoD verified)
EOF

# At end of session
cat >> sf3-implementation/WORKLOG.md <<EOF
- **Items Completed**: $SESSION_COMPLETED_COUNT
- **Items Blocked**: $SESSION_BLOCKED_COUNT
- **Next Steps**: [DESCRIBE]
- **Duration**: ${SESSION_DURATION} minutes

EOF

# Commit worklog with checklist updates
git add sf3-implementation/WORKLOG.md
git commit -m "worklog: session complete - $SESSION_COMPLETED_COUNT items [R600+R601]"
git push
```

**Worklog is the audit trail - MANDATORY for every session.**

## Continuation/Resumption Protocol

**When resuming work after any interruption:**

```bash
# 1. Load execution state
READ: /home/vscode/software-factory-template/sf3-implementation/execution-state.json

current_week=$(jq -r '.current_week' execution-state.json)
current_phase=$(jq -r '.current_phase' execution-state.json)
next_item_index=$(jq -r '.next_item_index' execution-state.json)
completed_items=$(jq -r '.completed_items' execution-state.json)
blocked_items=$(jq -r '.blocked_items' execution-state.json)

# 2. Load worklog to understand recent history
tail -50 sf3-implementation/WORKLOG.md

# 3. Locate next unchecked item in checklist
NEXT_TASK=$(sed -n "${next_item_index}p" SF-3.0-COMPREHENSIVE-EXECUTION-CHECKLIST.md)

# 4. Report resumption state
echo "🔄 RESUMING CHECKLIST EXECUTION"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Week: $current_week | Phase: $current_phase"
echo "Progress: $completed_items / $total_items (${completion_percentage}%)"
echo "Blocked: $blocked_items items"
echo ""
echo "Next Item (#${next_item_index}):"
echo "$NEXT_TASK"
echo ""
echo "Recent History:"
tail -10 sf3-implementation/WORKLOG.md
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 5. Ask user if should proceed
# (In automation context, check for blockers)
if [ $blocked_items -gt 0 ]; then
    echo "⚠️  WARNING: $blocked_items items currently blocked"
    echo "   Review blockers in: sf3-implementation/BLOCKERS.md"
    echo ""
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE  # Blockers must be resolved first"
    exit 600
fi

# 6. Proceed with next item
echo "▶️  PROCEEDING with Item #${next_item_index}"
```

**Resumption ALWAYS starts with state assessment - NEVER guess position.**

## N/A Item Protocol

**When item is not applicable:**

```bash
# 1. Document justification
JUSTIFICATION="[Provide detailed reason why this item is N/A]"

# 2. Update checklist with N/A marking
sed -i "s/^- \[ \] ${TASK}/- [N\/A] ${TASK} (Justification: $JUSTIFICATION - @$(whoami) $(date '+%Y-%m-%d'))/" \
    SF-3.0-COMPREHENSIVE-EXECUTION-CHECKLIST.md

# 3. Log in worklog
cat >> sf3-implementation/WORKLOG.md <<EOF
- ⊘ Item #${NEXT_ITEM_INDEX}: $TASK (N/A - $JUSTIFICATION)
EOF

# 4. Update execution state (counts as completed for progress %)
jq '.completed_items += 1 |
    .next_item_index += 1 |
    .na_items += 1' \
   execution-state.json > tmp.json && mv tmp.json execution-state.json

# 5. Commit
git add -A
git commit -m "checklist: item #${NEXT_ITEM_INDEX} marked N/A - $JUSTIFICATION [R600]"
git push
```

**N/A items MUST have justification - no silent skips.**

## Integration with Planning Files

**Reference planning files for context:**

```bash
# Before starting phase/week work, load planning files
case "$current_phase" in
    "WEEK_1-2")
        READ: /home/vscode/software-factory-template/docs/SOFTWARE-FACTORY-3.0-ARCHITECTURE.md
        READ: /home/vscode/software-factory-template/docs/RULE-MIGRATION-PLAN-SF3.md
        ;;
    "WEEK_3-4")
        # State Manager implementation context
        READ: /home/vscode/software-factory-template/docs/SOFTWARE-FACTORY-3.0-ARCHITECTURE.md (Part 5)
        ;;
    # ... etc for each phase
esac

# Extract relevant requirements for current items
echo "📖 PLANNING CONTEXT FOR CURRENT WORK:"
[Display relevant excerpts]
```

**Planning files provide the "why" - checklist provides the "what".**

## Validation Script Integration

**Use validation commands exactly as specified:**

```bash
# Example from checklist:
# - [ ] Verify SF 3.0 architecture document is complete
#   - **DoD**: docs/SOFTWARE-FACTORY-3.0-ARCHITECTURE.md contains all 10 parts (INIT through Part 9.5)
#   - **Validation**: grep -c "^## Part" docs/SOFTWARE-FACTORY-3.0-ARCHITECTURE.md
#   - **Owner**: Orchestrator

# Execute validation EXACTLY as written
VALIDATION_CMD='grep -c "^## Part" /home/vscode/software-factory-template/docs/SOFTWARE-FACTORY-3.0-ARCHITECTURE.md'
EXPECTED_RESULT=11  # (includes Part 3.5, 7.5, 9.5)

ACTUAL_RESULT=$(eval "$VALIDATION_CMD")

if [ "$ACTUAL_RESULT" -eq "$EXPECTED_RESULT" ]; then
    echo "✅ Validation PASSED: $ACTUAL_RESULT parts found (expected $EXPECTED_RESULT)"
    CHECK_BOX=true
else
    echo "❌ Validation FAILED: $ACTUAL_RESULT parts found (expected $EXPECTED_RESULT)"
    CHECK_BOX=false
    exit 600
fi
```

**Validation commands are EXACT SPECIFICATIONS - do not modify.**

## Enforcement & Penalties

### CATASTROPHIC FAILURES (-100%)
- Checking box without DoD verification
- Validation failed but box checked anyway
- Skipping items without N/A justification
- Lost progress (no commits for >1 hour)

### CRITICAL FAILURES (-50% to -75%)
- Wrong validation command used
- Self-validation of critical items (🔴)
- Proceeding past blockers
- Missing worklog entries

### STANDARD FAILURES (-20% to -40%)
- Incorrect execution-state.json updates
- Missing blocker documentation
- Commit frequency violations
- Inaccurate progress reporting

## Success Criteria

Checklist execution is successful when:

✅ ALL checkboxes checked (or marked N/A with justification)
✅ ALL DoD criteria validated and verified
✅ execution-state.json shows 100% completion
✅ Zero open blockers
✅ Complete worklog for all sessions
✅ All files committed and pushed

## Related Rules

- **R288**: State File Update and Commit Protocol (commit frequency)
- **R405**: Automation Continuation Flag (CONTINUE-SOFTWARE-FACTORY output)
- **R601**: Worklog Maintenance Protocol (logging requirements)
- **R510**: State Execution Checklist Compliance (general checklist rules)

## File References

- **Checklist**: `/home/vscode/software-factory-template/docs/SF-3.0-COMPREHENSIVE-EXECUTION-CHECKLIST.md`
- **Execution State**: `/home/vscode/software-factory-template/sf3-implementation/execution-state.json`
- **Worklog**: `/home/vscode/software-factory-template/sf3-implementation/WORKLOG.md`
- **Blockers**: `/home/vscode/software-factory-template/sf3-implementation/BLOCKERS.md`

---

**FINAL REMINDER**: This is a SUPREME LAW. Checking boxes without DoD validation = -100% CATASTROPHIC FAILURE. The checklist is your contract - honor every line.
