# SW-Engineer - SPLIT_IMPLEMENTATION State Rules

## 🔴🔴🔴 CRITICAL: SPLIT IMPLEMENTATION REQUIREMENTS 🔴🔴🔴

**YOU ARE IMPLEMENTING A SPLIT OF A TOO-LARGE BRANCH - SPECIAL PROTOCOLS APPLY!**

## State Context
You are in SPLIT_IMPLEMENTATION state because an effort exceeded size limits and must be split into smaller branches.

## 🔴🔴🔴 CRITICAL: Sequential Split Branching 🔴🔴🔴

**SPLITS ARE CREATED SEQUENTIALLY - EACH BASED ON THE PREVIOUS!**

### The Sequential Chain:
```
Original Effort (1200 lines from phase-integration) - TOO LARGE!
    ↓
Split-001 (400 lines from phase-integration)
    ↓
Split-002 (400 lines from split-001) ← NOT from phase-integration!
    ↓
Split-003 (400 lines from split-002) ← NOT from phase-integration!
```

### Why This Matters for You:
1. **Line Counting**: Your split will measure ONLY what you add
2. **Dependencies**: You can use code from previous splits
3. **No Duplicatio**: Don't re-implement what previous splits did
4. **Clean Merging**: Your work builds progressively

### How to Verify:
```bash
# Check what branch you're based on
git log --oneline -1 --format="%B" | grep "from branch:"
# Should show previous split for split-002 and later
```

## Primary Responsibilities

### 1. Verify Split Infrastructure (R204)
```bash
# MANDATORY: Verify orchestrator created split infrastructure
if [[ ! -f "SPLIT-PLAN-*.md" ]]; then
    echo "❌ FATAL: No split plan found!"
    echo "Orchestrator must create split infrastructure first per R204"
    exit 1
fi

# Verify you're in a split directory
if [[ $(pwd) != *"--split-"* ]]; then
    echo "❌ FATAL: Not in a split directory!"
    exit 1
fi
```

### 2. Sequential Split Implementation (R202, R205)
- **R202**: Single SW engineer implements ALL splits for an effort
- **R205**: Splits must be implemented SEQUENTIALLY, not in parallel
- Complete split-001, then split-002, then split-003, etc.

### 3. Size Compliance (R007, R200)
- Each split MUST stay under 800 lines (HARD LIMIT)
- Use line-counter.sh with BRANCH NAMES (not directory names!)
- Remember: Splits are in SEPARATE git repositories
- Stop immediately if approaching limit

### 4. 🚨🚨🚨 UPDATE SPLIT TRACKING (R302) 🚨🚨🚨

**MANDATORY: Report split details for comprehensive tracking:**
```bash
# After completing each split, report to orchestrator:
SPLIT_REPORT="Split $CURRENT_SPLIT of $TOTAL_SPLITS complete
Branch: $SPLIT_BRANCH
Lines: $(../../tools/line-counter.sh -b $BASE_BRANCH -c $(git branch --show-current) | grep Total | awk '{print $NF}')
Description: [What this split contains]
Status: COMPLETED"

echo "$SPLIT_REPORT" > /tmp/split-${EFFORT_NAME}-${CURRENT_SPLIT}.status
```

### 5. 🚨🚨🚨 NOTIFY ORCHESTRATOR WHEN ALL SPLITS COMPLETE (R296, R302) 🚨🚨🚨

**CRITICAL: After completing the FINAL split for an effort:**

```bash
# When you complete the LAST split (e.g., split-003 of 3 total)
if [ "$CURRENT_SPLIT" == "$TOTAL_SPLITS" ]; then
    echo "═══════════════════════════════════════════════════════════════"
    echo "✅ ALL SPLITS COMPLETE FOR EFFORT: $EFFORT_NAME"
    echo "═══════════════════════════════════════════════════════════════"
    
    # Create completion marker for orchestrator
    COMPLETION_FILE="/tmp/splits-complete-${EFFORT_NAME}.marker"
    cat > "$COMPLETION_FILE" << EOF
EFFORT: $EFFORT_NAME
TOTAL_SPLITS: $TOTAL_SPLITS
COMPLETED_AT: $(date -u +%Y-%m-%dT%H:%M:%SZ)
SPLITS_BRANCHES:
$(for i in $(seq 1 $TOTAL_SPLITS); do
    echo "  - ${EFFORT_NAME}--split-$(printf "%03d" $i)"
done)
STATUS: READY_FOR_DEPRECATION
ORIGINAL_BRANCH: $ORIGINAL_BRANCH
ACTION_REQUIRED: Orchestrator must mark original branch as deprecated per R296
EOF
    
    echo "📋 Completion marker created at: $COMPLETION_FILE"
    echo ""
    echo "🚨 ORCHESTRATOR ACTION REQUIRED:"
    echo "  1. Mark original branch as DEPRECATED per R296"
    echo "  2. Update state file with SPLIT_DEPRECATED status"
    echo "  3. List replacement splits in state file"
    echo ""
    echo "Original branch to deprecate: $ORIGINAL_BRANCH"
    echo "Replacement splits to use for integration:"
    for i in $(seq 1 $TOTAL_SPLITS); do
        echo "  - ${EFFORT_NAME}--split-$(printf "%03d" $i)"
    done
fi
```

## Required Actions for Each Split

### 1. Read Split Plan
```bash
SPLIT_PLAN=$(ls SPLIT-PLAN-*.md | head -1)
echo "Reading split plan: $SPLIT_PLAN"
cat "$SPLIT_PLAN"
```

### 2. Implement Only What's in Plan
- Follow the split plan EXACTLY
- Do NOT add extra features
- Do NOT refactor beyond the plan
- Stay within size limits

### 3. Measure After Implementation
```bash
# CRITICAL: Measurement depends on split number!
CURRENT_BRANCH=$(git branch --show-current)
SPLIT_NUM=$(echo "$CURRENT_BRANCH" | grep -o 'split-[0-9]*' | grep -o '[0-9]*')

# SEQUENTIAL BRANCHING: Each split measured against PREVIOUS split
if [ "$SPLIT_NUM" = "001" ]; then
    # First split: measure against original base (e.g., phase-integration)
    BASE_BRANCH="phase1-integration"  # Same base as original effort
else
    # Subsequent splits: measure against PREVIOUS split
    PREV_NUM=$(printf "%03d" $((10#$SPLIT_NUM - 1)))
    BASE_BRANCH="${CURRENT_BRANCH/split-${SPLIT_NUM}/split-${PREV_NUM}}"
fi

echo "📊 Measuring split-${SPLIT_NUM} against: $BASE_BRANCH"
echo "   (Each split measured against previous, NOT original base)"

# Use line-counter.sh with CORRECT base
$CLAUDE_PROJECT_DIR/tools/line-counter.sh -b $BASE_BRANCH -c $CURRENT_BRANCH

# Get the actual line count
LINES=$($CLAUDE_PROJECT_DIR/tools/line-counter.sh -b $BASE_BRANCH -c $CURRENT_BRANCH | grep Total | awk '{print $NF}')
if [ "$LINES" -gt 800 ]; then
    echo "❌ CRITICAL: Split exceeds 800 lines!"
    echo "   This split adds $LINES lines to $BASE_BRANCH"
    exit 1
fi
```

### 4. Commit and Push
```bash
# Commit with clear message
git add -A
git commit -m "feat: implement split-${SPLIT_NUM} - ${DESCRIPTION}"
git push
```

## State Transitions

### After Each Split:
- If more splits remain → Stay in SPLIT_IMPLEMENTATION
- If all splits complete → Transition to COMPLETE
- If size exceeded → Transition to ERROR

### Success Criteria:
- ✅ All splits implemented sequentially
- ✅ Each split under 800 lines
- ✅ All splits pushed to remote
- ✅ Orchestrator notified when complete (R296)

## Related Rules
- **R202**: Single agent per split effort
- **R204**: Split infrastructure creation
- **R205**: Sequential split navigation
- **R207**: Split boundary validation
- **R296**: Deprecated branch marking protocol
- **R007**: Size limit compliance
- **R200**: Measure only changeset

## Common Violations to Avoid

### ❌ Implementing Splits in Parallel
```bash
# WRONG - Opening multiple terminals
terminal1: implement split-001
terminal2: implement split-002  # VIOLATION!
```

### ❌ Exceeding Size Limits
```bash
# WRONG - Adding extra features
"While I'm here, let me also refactor..."  # NO!
```

### ❌ Not Notifying Orchestrator
```bash
# WRONG - Completing without notification
"Split 3 done, my work is complete" [exits]  # MUST NOTIFY!
```

### ✅ Correct Pattern
```bash
1. Implement split-001 → measure → commit → push
2. Implement split-002 → measure → commit → push  
3. Implement split-003 → measure → commit → push
4. Notify orchestrator that ALL splits are complete
5. Orchestrator marks original branch as deprecated
```

## Acknowledgment Required
After reading these rules, acknowledge:
"✅ Successfully read SPLIT_IMPLEMENTATION rules for sw-engineer. I understand I must notify the orchestrator when ALL splits are complete per R296."