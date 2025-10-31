# MONITORING_SWE_PROGRESS State Rules

## State Context

**Current Phase**: Read from `orchestrator-state-v3.json`
**Current Wave**: Read from `orchestrator-state-v3.json`
**Entry From**: SPAWN_SW_ENGINEERS
**Exit To**: SPAWN_CODE_REVIEWERS_EFFORT_REVIEW, MONITORING_SWE_PROGRESS (recursive)

## 🚨🚨🚨 PRIMARY OBJECTIVE 🚨🚨🚨

**MONITOR progress of Software Engineer agents implementing efforts, checking for:**
- Implementation completion
- Size violations (>800 lines = HARD STOP)
- Errors or blocks
- Need for intervention

This is an ACTIVE MONITORING_SWE_PROGRESS state - the orchestrator periodically checks on agent progress and takes action when needed.

## Required Inputs

### 1. List of Active Implementations
```bash
PHASE=$(jq -r '.current_phase' orchestrator-state-v3.json)
WAVE=$(jq -r '.current_wave' orchestrator-state-v3.json)

# Get efforts currently being implemented
EFFORTS_IN_PROGRESS=$(jq -r '.efforts_in_progress[]?' orchestrator-state-v3.json 2>/dev/null)

if [ -z "$EFFORTS_IN_PROGRESS" ]; then
    echo "⚠️ No efforts in progress - may have all completed"
    # Check if we should transition to review spawning
fi

echo "👀 Monitoring implementation of:"
echo "$EFFORTS_IN_PROGRESS"
```

### 2. Expected Agent Working Directories
```bash
for effort in $EFFORTS_IN_PROGRESS; do
    EFFORT_DIR="efforts/phase${PHASE}/wave${WAVE}/${effort}"
    echo "  - $effort: $EFFORT_DIR"
done
```

## 🔴🔴🔴 MONITORING_SWE_PROGRESS PROTOCOL 🔴🔴🔴

### Monitoring Cycle (Every 5 Messages)

```bash
echo "🔍 IMPLEMENTATION MONITORING_SWE_PROGRESS CYCLE - $(date -Iseconds)"
echo "═══════════════════════════════════════════════════════"

COMPLETED_EFFORTS=""
BLOCKED_EFFORTS=""
SIZE_VIOLATIONS=""
EFFORTS_NEEDING_REVIEW=""

for effort in $EFFORTS_IN_PROGRESS; do
    echo ""
    echo "📊 Checking: $effort"

    EFFORT_DIR="efforts/phase${PHASE}/wave${WAVE}/${effort}"
    EFFORT_BRANCH="phase${PHASE}/wave${WAVE}/${effort}"

    # Check 1: Read work log for status
    WORKLOG=$(ls -t "$EFFORT_DIR/.software-factory/phase${PHASE}/wave${WAVE}/${effort}/work-log--"*.log 2>/dev/null | head -1)

    if [ -f "$WORKLOG" ]; then
        LAST_STATUS=$(grep "^## Status" "$WORKLOG" -A 1 | tail -1)
        echo "  Last status: $LAST_STATUS"

        # Check for completion indicators
        if echo "$LAST_STATUS" | grep -qi "COMPLETE\|READY.*REVIEW"; then
            echo "  ✅ Effort appears complete"
            COMPLETED_EFFORTS="$COMPLETED_EFFORTS $effort"
            EFFORTS_NEEDING_REVIEW="$EFFORTS_NEEDING_REVIEW $effort"
            continue
        fi

        # Check for error indicators
        if echo "$LAST_STATUS" | grep -qi "ERROR\|BLOCKED\|FAILED\|STUCK"; then
            echo "  ⚠️ Effort may be blocked"
            BLOCKED_EFFORTS="$BLOCKED_EFFORTS $effort"
        fi
    else
        echo "  ⚠️ No work log found"
    fi

    # Check 2: Measure size using line counter
    echo "  📏 Measuring size..."

    # Get base branch from state or infrastructure
    BASE_BRANCH=$(jq -r ".infrastructure_created.efforts.\"${effort}\".base_branch" orchestrator-state-v3.json)

    if [ -z "$BASE_BRANCH" ] || [ "$BASE_BRANCH" = "null" ]; then
        echo "  ⚠️ Cannot determine base branch, skipping size check"
        continue
    fi

    # Run line counter
    cd "$EFFORT_DIR"
    SIZE_RESULT=$("$CLAUDE_PROJECT_DIR/tools/line-counter.sh" -b "$BASE_BRANCH" -c "$EFFORT_BRANCH" 2>&1)
    TOTAL_LINES=$(echo "$SIZE_RESULT" | grep "Total lines changed" | grep -oE "[0-9]+" | head -1)

    if [ -n "$TOTAL_LINES" ]; then
        echo "  Current size: $TOTAL_LINES lines"

        # Hard limit check (R220)
        if [ "$TOTAL_LINES" -gt 800 ]; then
            echo "  🚨 HARD LIMIT VIOLATION: $TOTAL_LINES > 800 lines"
            SIZE_VIOLATIONS="$SIZE_VIOLATIONS $effort"
        elif [ "$TOTAL_LINES" -gt 700 ]; then
            echo "  ⚠️ Approaching limit: $TOTAL_LINES > 700 lines (soft limit)"
        else
            echo "  ✅ Size within limits"
        fi
    else
        echo "  ⚠️ Could not measure size"
    fi

    cd "$CLAUDE_PROJECT_DIR"

    # Check 3: Git status
    cd "$EFFORT_DIR"
    UNCOMMITTED=$(git status --porcelain 2>/dev/null | wc -l)
    if [ "$UNCOMMITTED" -gt 0 ]; then
        echo "  ⚠️ Uncommitted changes: $UNCOMMITTED files"
    else
        echo "  ✅ All changes committed"
    fi
    cd "$CLAUDE_PROJECT_DIR"

    # Check 4: Look for recent activity
    cd "$EFFORT_DIR"
    LAST_COMMIT=$(git log -1 --format="%cr" "$EFFORT_BRANCH" 2>/dev/null)
    echo "  Last commit: $LAST_COMMIT"
    cd "$CLAUDE_PROJECT_DIR"
done

echo ""
echo "═══════════════════════════════════════════════════════"
```

### Decision Logic

```bash
# Priority 1: Handle Size Violations (HARD STOP)
if [ -n "$SIZE_VIOLATIONS" ]; then
    echo "🚨 SIZE VIOLATIONS DETECTED"
    echo "Efforts exceeding 800 lines: $SIZE_VIOLATIONS"
    echo ""
    echo "REQUIRED ACTION:"
    echo "1. STOP these implementations immediately"
    echo "2. Spawn Code Reviewer to create SPLIT PLAN"
    echo "3. Cannot proceed until split"
    echo ""
    # Orchestrator should transition to SPAWN_CODE_REVIEWER_SPLIT_PLAN state
    exit 1
fi

# Priority 2: Handle Completed Efforts
if [ -n "$COMPLETED_EFFORTS" ]; then
    echo "✅ COMPLETED EFFORTS DETECTED"
    echo "Ready for review: $COMPLETED_EFFORTS"
    echo ""
    echo "REQUIRED ACTION:"
    echo "1. Spawn Code Reviewers for these efforts"
    echo "2. Keep monitoring other in-progress efforts"
    echo ""
    # Should trigger code reviewer spawning
fi

# Priority 3: Handle Blocked Efforts
if [ -n "$BLOCKED_EFFORTS" ]; then
    echo "⚠️ BLOCKED EFFORTS DETECTED"
    echo "May need intervention: $BLOCKED_EFFORTS"
    echo ""
    echo "SUGGESTED ACTION:"
    echo "1. Investigate work logs"
    echo "2. Check for error messages"
    echo "3. May need to respawn or provide guidance"
    echo ""
fi

# Priority 4: Continue Monitoring
STILL_IN_PROGRESS=$(echo "$EFFORTS_IN_PROGRESS" | wc -w)
COMPLETED_COUNT=$(echo "$COMPLETED_EFFORTS" | wc -w)
REMAINING=$((STILL_IN_PROGRESS - COMPLETED_COUNT))

if [ "$REMAINING" -gt 0 ]; then
    echo "⏳ Still monitoring $REMAINING efforts in progress"
    echo "Next check in 5 messages or 3 minutes"
else
    echo "✅ All efforts complete - transitioning to reviews"
fi
```

## State Transition Logic

### Case 1: Some Efforts Complete, Others In Progress
```bash
# Spawn reviewers for completed efforts
# Stay in MONITORING_SWE_PROGRESS for remaining efforts
jq --argjson completed_list "$(echo $COMPLETED_EFFORTS | jq -R -s -c 'split(" ") | map(select(length > 0))')" \
   '.efforts_ready_for_review = $completed_list |
    .state_machine.current_state = "SPAWN_CODE_REVIEWERS_EFFORT_REVIEW"' \
   orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
```

### Case 2: All Efforts Complete
```bash
# Spawn reviewers for all efforts
# Transition to MONITORING_EFFORT_REVIEWS
jq --argjson all_efforts "$(echo $EFFORTS_IN_PROGRESS | jq -R -s -c 'split(" ") | map(select(length > 0))')" \
   '.efforts_in_progress = null |
    .efforts_ready_for_review = $all_efforts |
    .state_machine.current_state = "SPAWN_CODE_REVIEWERS_EFFORT_REVIEW"' \
   orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
```

### Case 3: Size Violations Detected
```bash
# Stop implementations, create split plans
jq --argjson violations "$(echo $SIZE_VIOLATIONS | jq -R -s -c 'split(" ") | map(select(length > 0))')" \
   '.size_violations = $violations |
    .state_machine.current_state = "SPAWN_CODE_REVIEWER_SPLIT_PLAN"' \
   orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
```

### Case 4: Still In Progress
```bash
# Update monitoring timestamp, stay in same state
jq --arg timestamp "$(date -Iseconds)" \
   '.last_monitoring_check = $timestamp |
    .monitoring_cycles = (.monitoring_cycles // 0) + 1' \
   orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
```

## Monitoring Frequency

### Active Monitoring (R151 Adjacent)
- Check every **5 messages** exchanged with user
- Check every **3 minutes** of wall clock time
- Immediate check on any agent completion signal
- Immediate check on any error signal

### What to Monitor
1. **Work logs**: Status updates from SW Engineers
2. **Git status**: Committed vs uncommitted changes
3. **Branch size**: Line counts via line-counter.sh
4. **Recent activity**: Last commit time
5. **Error indicators**: Grep work logs for ERROR/FAILED

## Integration with Rules

- **R220/R221**: Size Limits and Continuous Delivery
- **R304**: Mandatory Line Counter Tool Enforcement
- **R287**: TODO Persistence (orchestrator must save frequently)
- **R208/R209**: Directory Isolation Monitoring
- **Grading**: Workflow Compliance (25%), Size Compliance (20%)

## Exit Criteria

### Exit to SPAWN_CODE_REVIEWERS_EFFORT_REVIEW when:
- ✅ At least one effort complete
- ✅ No size violations blocking
- ✅ Implementation work committed and pushed

### Stay in MONITORING_SWE_PROGRESS when:
- ⏳ Efforts still in progress
- ⏳ No completed efforts yet
- ⏳ No critical issues detected

### Exit to ERROR_RECOVERY when:
- ❌ Multiple efforts blocked
- ❌ Critical errors detected
- ❌ Agents not responding

## Common Issues

### Issue: Agents Working Too Long
**Detection**: Implementation >2 hours with no completion
**Resolution**: Check work logs, may need to provide guidance or respawn

### Issue: Size Creeping Up
**Detection**: Efforts approaching 700-800 line range
**Resolution**: Alert SW Engineers, prepare for potential split

### Issue: No Recent Activity
**Detection**: No commits in >30 minutes on active effort
**Resolution**: Check if agent is stuck, may need intervention

## Exit Conditions and Continuation Flag

### ✅ CONTINUE-SOFTWARE-FACTORY=TRUE (Standard Operation - DEFAULT)

**Use TRUE when:**
- ✅ Monitoring implementation progress (normal loop)
- ✅ Some efforts still in progress
- ✅ Efforts completing successfully
- ✅ Size violations detected (system will handle splits)
- ✅ Ready to spawn reviewers for completed work
- ✅ Transitioning to split planning or review spawning
- ✅ Following designed workflow

**THIS IS NORMAL WORKFLOW.** Monitoring implementation, detecting completions,
detecting size violations, and transitioning to appropriate next states is the
DESIGNED PROCESS. These are all EXPECTED behaviors of the monitoring loop.

### 🔴 CONTINUE-SOFTWARE-FACTORY=FALSE (Exceptional Conditions ONLY)

**Use FALSE ONLY when:**
- ❌ Cannot access state file or effort metadata
- ❌ All efforts show no progress for >1 hour
- ❌ Critical infrastructure corruption
- ❌ Cannot determine what to monitor
- ❌ State machine corruption

**DO NOT set FALSE because:**
- ❌ Size violations detected (system handles this!)
- ❌ Need to spawn reviewers (NORMAL transition!)
- ❌ Need split planning (EXPECTED workflow!)
- ❌ Still monitoring (NORMAL loop!)
- ❌ R322 requires stop (stop ≠ FALSE!)

**Correct pattern:** All monitoring outcomes use TRUE

## Continuation Control

```bash
# Monitoring states continue in a loop until exit condition met
# Check if we should continue monitoring or transition
if [ "$REMAINING" -gt 0 ]; then
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # Continue monitoring
else
    # All complete, transition based on results
    if [ -n "$SIZE_VIOLATIONS" ]; then
        echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # Continue per R405 - system will handle split planning
    elif [ -n "$COMPLETED_EFFORTS" ]; then
        echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # Continue per R405 - system will spawn reviewers
    else
        echo "CONTINUE-SOFTWARE-FACTORY=TRUE"   # Continue monitoring
    fi
fi
```

---

**REMEMBER**: Active monitoring is CRITICAL. The orchestrator must check frequently and react quickly to size violations or completion!

## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

**MUST output as LAST line before state completion:**
- `CONTINUE-SOFTWARE-FACTORY=TRUE` - State completed, automation can continue
- `CONTINUE-SOFTWARE-FACTORY=FALSE` - Error/block, manual intervention required

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-CONTINUATION-FLAG-MASTER-GUIDE.md**
