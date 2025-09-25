# Orchestrator - SPLIT_INIT State Rules
## SPLITTING Sub-State Machine Entry Point

## State Context
You have entered the SPLITTING sub-state machine to handle an oversized effort implementation. This is a specialized workflow to split the effort into compliant chunks while maintaining quality and functionality.

## Primary Directives

### 1. Initialize Split Context
```bash
# Create split tracking structure
init_split_context() {
    local EFFORT_NAME="$1"
    local ORIGINAL_SIZE="$2"
    local EFFORT_BRANCH="$3"

    echo "📊 Initializing split operation for: $EFFORT_NAME"
    echo "📏 Original size: $ORIGINAL_SIZE lines (limit: 800)"

    # Load effort details from main state
    EFFORT_DIR=$(jq -r --arg effort "$EFFORT_NAME" \
        '.efforts_in_progress[] | select(.effort_id == $effort) | .directory' \
        orchestrator-state.json)

    # Initialize split state file
    jq --arg effort "$EFFORT_NAME" \
       --arg size "$ORIGINAL_SIZE" \
       --arg branch "$EFFORT_BRANCH" \
       --arg dir "$EFFORT_DIR" \
       '.original_effort = {
          "effort_id": $effort,
          "size": ($size | tonumber),
          "branch": $branch,
          "directory": $dir
       } | .current_state = "SPLIT_INIT"' \
       splitting-${EFFORT_NAME}-state.json > tmp.json && \
       mv tmp.json splitting-${EFFORT_NAME}-state.json
}
```

### 2. Validate Split Prerequisites
- Verify effort exists and is accessible
- Confirm size violation is real (not measurement error)
- Check for existing split attempts
- Ensure no active work on the effort

### 3. Prepare for Analysis
- Record original effort state
- Set up tracking for split operations
- Initialize quality gate tracking
- Prepare for Code Reviewer spawn

## State Transitions

### Valid Transitions
- **SPLIT_INIT → SPLIT_ANALYSIS**: Always proceed to analysis after initialization

### Invalid Transitions
- ❌ Cannot skip to SPLIT_PLANNING
- ❌ Cannot return to main state without completion
- ❌ Cannot proceed without valid effort context

## Required Actions

### 1. Load Effort Context
```bash
# Must load from actual effort directory
cd /efforts/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}
git checkout ${EFFORT_BRANCH}
git pull origin ${EFFORT_BRANCH}
```

### 2. Verify Size Violation
```bash
# MUST use line-counter.sh (R304)
LINE_COUNT=$($CLAUDE_PROJECT_DIR/tools/line-counter.sh | grep "Total" | awk '{print $NF}')

if [ "$LINE_COUNT" -le 800 ]; then
    echo "⚠️ WARNING: Size is actually compliant ($LINE_COUNT ≤ 800)"
    echo "🔄 Returning to main state - no split needed"
    # Exit sub-state machine
fi
```

### 3. Check for Previous Splits
```bash
# Look for existing split branches
git branch -r | grep "${EFFORT_NAME}--split-"
if [ $? -eq 0 ]; then
    echo "⚠️ WARNING: Existing splits detected"
    # Load split history and continue from last point
fi
```

### 4. Spawn Code Reviewer for Analysis
```bash
# Transition to SPLIT_ANALYSIS and spawn Code Reviewer
update_split_state "SPLIT_ANALYSIS"
spawn_code_reviewer_for_split_analysis "$EFFORT_NAME" "$LINE_COUNT"
```

## Quality Gates

### Entry Gate
- ✅ Valid effort with actual size violation
- ✅ No active implementation work
- ✅ Clean git state (no uncommitted changes)

### Exit Gate
- ✅ Split context initialized
- ✅ Size violation confirmed
- ✅ Ready for analysis

## State File Updates

### Required Updates
```json
{
  "current_state": "SPLIT_INIT",
  "initialized_at": "timestamp",
  "original_effort": {
    "effort_id": "effort-name",
    "size": 1250,
    "branch": "phase1-wave2-effort2",
    "directory": "/efforts/phase1/wave2/effort2"
  },
  "validation": {
    "size_verified": true,
    "effort_accessible": true,
    "prerequisites_met": true
  }
}
```

## Error Handling

### Common Issues
1. **Effort Not Found**: Verify correct phase/wave/effort
2. **No Size Violation**: Return to main state
3. **Git Issues**: Clean and retry
4. **Previous Splits**: Load and continue

### Recovery Strategy
```bash
handle_split_init_error() {
    local ERROR_TYPE="$1"

    case "$ERROR_TYPE" in
        "EFFORT_NOT_FOUND")
            echo "❌ Cannot find effort - aborting split"
            transition_to_split_abort
            ;;
        "NO_VIOLATION")
            echo "✅ No split needed - returning to main"
            complete_sub_state_success
            ;;
        "GIT_ERROR")
            echo "⚠️ Git issues - attempting recovery"
            clean_git_state && retry_split_init
            ;;
    esac
}
```

## Rules Compliance

### R304 - Mandatory Line Counter Usage
- MUST use `$CLAUDE_PROJECT_DIR/tools/line-counter.sh`
- NEVER use wc -l or manual counting
- Size verification is critical

### R296 - Split Deprecation
- Original effort will be marked SPLIT_DEPRECATED
- Must preserve original for history
- Splits become new source of truth

### R204 - Infrastructure Creation
- Orchestrator will create split infrastructure
- Just-in-time creation per split
- Sequential dependencies maintained

## Success Criteria

Before transitioning to SPLIT_ANALYSIS:
- [ ] Split state file created and initialized
- [ ] Original effort context recorded
- [ ] Size violation confirmed (>800 lines)
- [ ] No previous incomplete splits
- [ ] Ready to spawn Code Reviewer

## Next State: SPLIT_ANALYSIS

Upon successful initialization:
1. Update state to SPLIT_ANALYSIS
2. Spawn Code Reviewer with effort context
3. Wait for analysis completion
4. Code Reviewer will analyze and create split plan

---

**Remember**: This is a sub-state machine. You are operating within the SPLITTING context, not the main orchestration flow. Focus solely on split operations.