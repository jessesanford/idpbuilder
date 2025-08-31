# SW-Engineer - SPLIT_IMPLEMENTATION State Rules

## 🔴🔴🔴 CRITICAL: SPLIT IMPLEMENTATION REQUIREMENTS 🔴🔴🔴

**YOU ARE IMPLEMENTING A SPLIT OF A TOO-LARGE BRANCH - SPECIAL PROTOCOLS APPLY!**

## State Context
You are in SPLIT_IMPLEMENTATION state because an effort exceeded size limits and must be split into smaller branches.

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
- Use line-counter.sh to measure ONLY your changes
- Stop immediately if approaching limit

### 4. 🚨🚨🚨 NOTIFY ORCHESTRATOR WHEN ALL SPLITS COMPLETE (R296) 🚨🚨🚨

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
# Use line-counter.sh to measure
$CLAUDE_PROJECT_DIR/tools/line-counter.sh

# Verify under limit
LINES=$(git diff --numstat origin/main | awk '{sum+=$1+$2} END {print sum}')
if [ "$LINES" -gt 800 ]; then
    echo "❌ CRITICAL: Split exceeds 800 lines!"
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