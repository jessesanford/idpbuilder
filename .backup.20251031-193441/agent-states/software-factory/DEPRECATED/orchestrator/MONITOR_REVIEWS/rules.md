# MONITORING_EFFORT_REVIEWS State Rules

## State Context

**Current Phase**: Read from `orchestrator-state-v3.json`
**Current Wave**: Read from `orchestrator-state-v3.json`
**Entry From**: SPAWN_CODE_REVIEWERS_EFFORT_REVIEW
**Exit To**: SPAWN_SW_ENGINEERS, WAVE_COMPLETE, MONITORING_EFFORT_REVIEWS (recursive)

## 🚨🚨🚨 PRIMARY OBJECTIVE 🚨🚨🚨

**MONITOR progress of Code Reviewer agents, checking for:**
- Review completion
- Issues found requiring fixes
- Clean approvals
- Size violations discovered

This is an ACTIVE MONITORING_SWE_PROGRESS state - the orchestrator checks review progress and manages the review-fix cycle.

## Required Inputs

### 1. List of Active Reviews
```bash
PHASE=$(jq -r '.current_phase' orchestrator-state-v3.json)
WAVE=$(jq -r '.current_wave' orchestrator-state-v3.json)

# Get efforts currently under review
REVIEWS_IN_PROGRESS=$(jq -r '.reviews_in_progress[]?' orchestrator-state-v3.json 2>/dev/null)

if [ -z "$REVIEWS_IN_PROGRESS" ]; then
    echo "⚠️ No reviews in progress - may have all completed"
fi

echo "👀 Monitoring reviews for:"
echo "$REVIEWS_IN_PROGRESS"
```

## 🔴🔴🔴 REVIEW MONITORING_SWE_PROGRESS PROTOCOL 🔴🔴🔴

### Monitoring Cycle (Every 5 Messages)

```bash
echo "🔍 CODE REVIEW MONITORING_SWE_PROGRESS CYCLE - $(date -Iseconds)"
echo "═══════════════════════════════════════════════════════"

COMPLETED_REVIEWS=""
REVIEWS_WITH_ISSUES=""
APPROVED_EFFORTS=""
SIZE_VIOLATIONS=""

for effort in $REVIEWS_IN_PROGRESS; do
    echo ""
    echo "📊 Checking review: $effort"

    EFFORT_DIR="efforts/phase${PHASE}/wave${WAVE}/${effort}"

    # Check 1: Look for review report
    REVIEW_REPORT=$(ls -t "$EFFORT_DIR/.software-factory/phase${PHASE}/wave${WAVE}/${effort}/CODE-REVIEW-REPORT--"*.md 2>/dev/null | head -1)

    if [ -z "$REVIEW_REPORT" ] || [ ! -f "$REVIEW_REPORT" ]; then
        echo "  ⏳ Review not yet complete (no report found)"
        continue
    fi

    echo "  ✅ Review report found: $REVIEW_REPORT"
    COMPLETED_REVIEWS="$COMPLETED_REVIEWS $effort"

    # Check 2: Parse review outcome
    # Look for approval indicator
    if grep -qi "APPROVED\|NO ISSUES\|LOOKS GOOD\|READY.*MERGE" "$REVIEW_REPORT"; then
        echo "  ✅ Review APPROVED - no issues found"
        APPROVED_EFFORTS="$APPROVED_EFFORTS $effort"
        continue
    fi

    # Check 3: Look for blocking issues
    BLOCKING_COUNT=$(grep -c "BLOCKING\|\\*\\*BLOCKING\\*\\*\|MUST FIX" "$REVIEW_REPORT" || echo "0")
    HIGH_COUNT=$(grep -c "HIGH\|\\*\\*HIGH\\*\\*" "$REVIEW_REPORT" || echo "0")
    MEDIUM_COUNT=$(grep -c "MEDIUM\|\\*\\*MEDIUM\\*\\*" "$REVIEW_REPORT" || echo "0")
    LOW_COUNT=$(grep -c "LOW\|\\*\\*LOW\\*\\*" "$REVIEW_REPORT" || echo "0")

    TOTAL_ISSUES=$((BLOCKING_COUNT + HIGH_COUNT + MEDIUM_COUNT + LOW_COUNT))

    if [ "$TOTAL_ISSUES" -gt 0 ]; then
        echo "  ⚠️ Issues found:"
        echo "    BLOCKING: $BLOCKING_COUNT"
        echo "    HIGH: $HIGH_COUNT"
        echo "    MEDIUM: $MEDIUM_COUNT"
        echo "    LOW: $LOW_COUNT"

        REVIEWS_WITH_ISSUES="$REVIEWS_WITH_ISSUES $effort"
    fi

    # Check 4: Look for size violations
    if grep -qi "SIZE.*VIOLATION\|EXCEEDS.*800\|TOO LARGE\|>.*800.*LINE" "$REVIEW_REPORT"; then
        echo "  🚨 SIZE VIOLATION detected in review"
        SIZE_VIOLATIONS="$SIZE_VIOLATIONS $effort"
    fi

    # Check 5: Look for fix instructions
    FIX_INSTRUCTIONS=$(ls -t "$EFFORT_DIR/.software-factory/phase${PHASE}/wave${WAVE}/${effort}/CODE-REVIEW-INSTRUCTIONS--"*.md 2>/dev/null | head -1)

    if [ -f "$FIX_INSTRUCTIONS" ]; then
        echo "  📝 Fix instructions found: $FIX_INSTRUCTIONS"
    else
        echo "  ⚠️ Issues found but no fix instructions - review may be incomplete"
    fi
done

echo ""
echo "═══════════════════════════════════════════════════════"
```

### Decision Logic

```bash
# Priority 1: Handle Size Violations (HARD STOP)
if [ -n "$SIZE_VIOLATIONS" ]; then
    echo "🚨 SIZE VIOLATIONS FOUND IN REVIEW"
    echo "Efforts with size issues: $SIZE_VIOLATIONS"
    echo ""
    echo "REQUIRED ACTION:"
    echo "1. These efforts CANNOT be merged"
    echo "2. Must create SPLIT PLANS"
    echo "3. Re-implement as multiple smaller efforts"
    echo ""
    # Transition to split planning state
    exit 1
fi

# Priority 2: Handle Completed Reviews with Issues
if [ -n "$REVIEWS_WITH_ISSUES" ]; then
    echo "⚠️ REVIEWS COMPLETED WITH ISSUES"
    echo "Need fixes for: $REVIEWS_WITH_ISSUES"
    echo ""
    echo "REQUIRED ACTION:"
    echo "1. Spawn SW Engineers to fix issues"
    echo "2. Provide CODE-REVIEW-INSTRUCTIONS to engineers"
    echo "3. After fixes, re-review"
    echo ""
    # Transition to spawn engineers for fixes
fi

# Priority 3: Handle Approved Efforts
if [ -n "$APPROVED_EFFORTS" ]; then
    echo "✅ REVIEWS APPROVED"
    echo "Clean implementations: $APPROVED_EFFORTS"
    echo ""
    echo "These efforts are ready for:"
    echo "- Wave completion (if all efforts done)"
    echo "- Wave integration"
    echo ""
fi

# Priority 4: Continue Monitoring
STILL_REVIEWING=$(echo "$REVIEWS_IN_PROGRESS" | wc -w)
COMPLETED_COUNT=$(echo "$COMPLETED_REVIEWS" | wc -w)
REMAINING=$((STILL_REVIEWING - COMPLETED_COUNT))

if [ "$REMAINING" -gt 0 ]; then
    echo "⏳ Still monitoring $REMAINING reviews in progress"
    echo "Next check in 5 messages or 3 minutes"
else
    echo "✅ All reviews complete"
fi
```

## State Transition Logic

### Case 1: Reviews Find Issues - Need Fixes
```bash
# Extract efforts needing fixes
EFFORTS_NEEDING_FIXES=$(echo "$REVIEWS_WITH_ISSUES" | tr ' ' '\n' | grep -v "^$")

jq --argjson needs_fixes "$(echo $EFFORTS_NEEDING_FIXES | jq -R -s -c 'split(" ") | map(select(length > 0))')" \
   --arg timestamp "$(date -Iseconds)" \
   '.efforts_needing_fixes = $needs_fixes |
    .reviews_in_progress = null |
    .state_machine.current_state = "SPAWN_SW_ENGINEERS" |
    .state_machine.previous_state = "MONITORING_EFFORT_REVIEWS" |
    .state_transition_log += [{
        "from": "MONITORING_EFFORT_REVIEWS",
        "to": "SPAWN_SW_ENGINEERS",
        "timestamp": $timestamp,
        "reason": "Reviews completed with issues requiring fixes"
    }]' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
```

### Case 2: All Reviews Approved - Wave Complete
```bash
# All efforts approved, move to wave completion
jq --argjson approved "$(echo $APPROVED_EFFORTS | jq -R -s -c 'split(" ") | map(select(length > 0))')" \
   --arg timestamp "$(date -Iseconds)" \
   '.reviews_in_progress = null |
    .efforts_completed += $approved |
    .state_machine.current_state = "WAVE_COMPLETE" |
    .state_machine.previous_state = "MONITORING_EFFORT_REVIEWS" |
    .state_transition_log += [{
        "from": "MONITORING_EFFORT_REVIEWS",
        "to": "WAVE_COMPLETE",
        "timestamp": $timestamp,
        "reason": "All reviews approved, wave ready for completion"
    }]' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
```

### Case 3: Size Violations - Must Split
```bash
# Critical size violations, need split planning
jq --argjson violations "$(echo $SIZE_VIOLATIONS | jq -R -s -c 'split(" ") | map(select(length > 0))')" \
   --arg timestamp "$(date -Iseconds)" \
   '.size_violations = $violations |
    .reviews_in_progress = null |
    .state_machine.current_state = "SPAWN_CODE_REVIEWER_SPLIT_PLAN" |
    .state_machine.previous_state = "MONITORING_EFFORT_REVIEWS" |
    .state_transition_log += [{
        "from": "MONITORING_EFFORT_REVIEWS",
        "to": "SPAWN_CODE_REVIEWER_SPLIT_PLAN",
        "timestamp": $timestamp,
        "reason": "Size violations detected in code review"
    }]' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
```

### Case 4: Still In Progress
```bash
# Update monitoring timestamp, stay in same state
jq --arg timestamp "$(date -Iseconds)" \
   '.last_review_check = $timestamp |
    .review_monitoring_cycles = (.review_monitoring_cycles // 0) + 1' \
   orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
```

## Review-Fix Cycle Management

### Tracking Review Iterations
```bash
# Track how many review cycles each effort has gone through
# This helps detect if we're stuck in fix-review loop

for effort in $REVIEWS_WITH_ISSUES; do
    ITERATION=$(jq -r ".effort_review_iterations.\"${effort}\" // 0" orchestrator-state-v3.json)
    ITERATION=$((ITERATION + 1))

    jq --arg effort "$effort" \
       --argjson iteration "$ITERATION" \
       '.effort_review_iterations[$effort] = $iteration' \
       orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json

    echo "📊 $effort: Review iteration $ITERATION"

    if [ "$ITERATION" -gt 3 ]; then
        echo "⚠️ WARNING: $effort has been through $ITERATION review cycles"
        echo "  May indicate architectural issues or unclear requirements"
    fi
done
```

## Monitoring Frequency

### Active Monitoring
- Check every **5 messages** exchanged
- Check every **3 minutes** of wall clock time
- Immediate check on any reviewer completion signal
- Immediate check on any error signal

### What to Monitor
1. **Review reports**: CODE-REVIEW-REPORT files
2. **Fix instructions**: CODE-REVIEW-INSTRUCTIONS files
3. **Issue counts**: BLOCKING/HIGH/MEDIUM/LOW
4. **Size violations**: Line count checks
5. **Approval status**: Clean reviews vs issues found

## Integration with Rules

- **R220/R221**: Size Limits enforcement
- **R304**: Line Counter usage validation
- **R287**: TODO Persistence
- **R383/R343**: Metadata file standards
- **Grading**: Workflow Compliance (25%), Quality Assurance (20%)

## Exit Criteria

### Exit to SPAWN_SW_ENGINEERS when:
- ✅ Reviews complete with fixable issues
- ✅ Fix instructions exist for all issues
- ✅ No size violations blocking

### Exit to WAVE_COMPLETE when:
- ✅ All reviews approved
- ✅ No outstanding issues
- ✅ All code within size limits
- ✅ All tests passing

### Exit to SPAWN_CODE_REVIEWER_SPLIT_PLAN when:
- 🚨 Size violations detected
- 🚨 Must split before proceeding

### Stay in MONITORING_EFFORT_REVIEWS when:
- ⏳ Reviews still in progress
- ⏳ No completed reviews yet

## Common Issues

### Issue: Reviews Taking Too Long
**Detection**: Reviews in progress >30 minutes
**Resolution**: Check reviewer status, may need intervention

### Issue: Inconsistent Review Quality
**Detection**: Some reviews detailed, others superficial
**Resolution**: Provide better review guidelines, may respawn reviewers

### Issue: Review-Fix Loop
**Detection**: Same effort failing review multiple times
**Resolution**: Escalate to architect, may need requirements clarification

## Exit Conditions and Continuation Flag

### ✅ CONTINUE-SOFTWARE-FACTORY=TRUE (Standard Operation - DEFAULT)

**Use TRUE when:**
- ✅ Reviews are still in progress (monitoring continues)
- ✅ Reviews completed with issues (transition to fix workflow)
- ✅ Reviews approved (transition to wave complete)
- ✅ Size violations detected (transition to split planning)
- ✅ Any normal monitoring or transition operation

**ALL OF THESE ARE STANDARD CASES.** The review → fix → re-review cycle is **NORMAL SOFTWARE DEVELOPMENT WORKFLOW**. Code review finding issues is the system **WORKING CORRECTLY**, not failing!

### 🔴 CONTINUE-SOFTWARE-FACTORY=FALSE (Exceptional Conditions ONLY)

**Use FALSE ONLY when:**
- ❌ Review report files are corrupt/unreadable
- ❌ Critical state machine corruption detected
- ❌ Cannot determine next state (logic error)
- ❌ Infrastructure completely broken
- ❌ Unrecoverable error in monitoring logic

### ❌ DO NOT SET FALSE JUST BECAUSE:

**These are NOT reasons to use FALSE:**
- ❌ Code review found issues (this is **NORMAL**)
- ❌ Issues need to be fixed (this is **EXPECTED**)
- ❌ Multiple efforts have issues (this is **DESIGNED FOR**)
- ❌ Size violations detected (this is **HANDLED BY WORKFLOW**)
- ❌ Reviews taking time (this is **NORMAL**)
- ❌ Second or third review iteration (this is **EXPECTED**)
- ❌ "User might want to see this" (only if truly exceptional)

## Continuation Control

```bash
# Monitoring states continue in a loop until exit condition met
# Determine if we should continue monitoring or transition

# MANDATORY: Print cascade status if active (R406 auto-reporting)
if [ -f "orchestrator-state-v3.json" ] && \
   jq -e '.fix_cascade_state.active == true' orchestrator-state-v3.json > /dev/null 2>&1; then
    echo ""
    echo "═══════════════════════════════════════════════════════════"
    echo "📊 R406 FIX CASCADE STATUS (automatic report)"
    echo "═══════════════════════════════════════════════════════════"
    source utilities/cascade-status-report.sh
    cascade_status_report
    echo "═══════════════════════════════════════════════════════════"
    echo ""
fi

if [ "$REMAINING" -gt 0 ]; then
    echo "⏳ Still monitoring reviews..."
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # Standard: Continue monitoring loop
else
    # All reviews complete, determine next action
    if [ -n "$SIZE_VIOLATIONS" ]; then
        echo "🚨 Size violations require split planning"
        echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # Standard: System will handle splits
    elif [ -n "$REVIEWS_WITH_ISSUES" ]; then
        echo "⚠️ Issues found, need fixes"
        echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # Standard: System will spawn fixers (NORMAL!)
    elif [ -n "$APPROVED_EFFORTS" ]; then
        echo "✅ All approved, wave complete"
        echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # Standard: Continue to wave complete
    else
        echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # Standard: Continue monitoring
    fi
fi
```

## Grading Impact

**Correct flag usage:**
- ✅ Using TRUE for normal operations: No penalty
- ✅ Using TRUE when reviews find issues: No penalty (this is NORMAL!)
- ✅ Using FALSE for genuine errors: No penalty

**Incorrect flag usage:**
- ❌ Using FALSE when reviews find issues: **-20%** (misunderstanding normal workflow)
- ❌ Pattern of excessive FALSE usage: **-50%** (defeats automation purpose)
- ❌ Using FALSE "just to be safe": **-20%** (indicates lack of confidence in system)

## Common Scenarios

### Scenario 1: Review Found 5 Blocking Issues
**Status:** NORMAL
**Flag:** TRUE
**Reason:** Fix workflow will handle this automatically

### Scenario 2: Size Violation Detected
**Status:** NORMAL
**Flag:** TRUE
**Reason:** Split planning workflow will handle this automatically

### Scenario 3: Third Review Iteration Still Finding Issues
**Status:** NORMAL (may need escalation, but not a stop condition)
**Flag:** TRUE
**Reason:** Review iterations are expected, fix cycle continues

### Scenario 4: Review Report File Corrupt/Unreadable
**Status:** ERROR
**Flag:** FALSE
**Reason:** Cannot proceed without valid review data

### Scenario 5: All Reviews Approved
**Status:** NORMAL
**Flag:** TRUE
**Reason:** Standard transition to wave complete

---

**REMEMBER**: The review-fix cycle is normal and expected. Code review finding issues means the system is **WORKING CORRECTLY**! Only use FALSE for genuine system failures, not normal workflow events.

## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

**MUST output as LAST line before state completion:**
- `CONTINUE-SOFTWARE-FACTORY=TRUE` - State completed, automation can continue
- `CONTINUE-SOFTWARE-FACTORY=FALSE` - Error/block, manual intervention required

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-CONTINUATION-FLAG-MASTER-GUIDE.md**
