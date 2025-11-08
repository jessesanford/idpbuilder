# Orchestrator - MONITORING_EFFORT_REVIEWS State Rules

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

## 🛑🛑🛑 R322 MANDATORY STOP BEFORE STATE TRANSITIONS 🛑🛑🛑

**SUPREME LAW - VIOLATION = -100% IMMEDIATE FAILURE**

### YOU MUST STOP AFTER:
1. ✅ Completing all TODOs for this state
2. ✅ Updating orchestrator-state-v3.json with new state
3. ✅ Committing and pushing the state file
4. ✅ Providing work summary

### YOU MUST NOT:
- ❌ Continue to the next state automatically
- ❌ Start work for the new state
- ❌ Spawn agents for the new state
- ❌ Assume permission to continue

### STOP PROTOCOL:
```markdown
## 🛑 STATE TRANSITION CHECKPOINT: CURRENT_STATE → NEXT_STATE

### ✅ Current State Work Completed:
- [List completed work]

### 📊 Current Status:
- Current State: CURRENT_STATE
- Next State: NEXT_STATE
- TODOs Completed: X/Y
- State Files: Updated and committed ✅

### ⏸️ STOPPED - Awaiting User Continuation
Ready to transition to NEXT_STATE. Please use /continue-software-factory.
```

**STOP MEANS STOP - Exit and wait for /continue-software-factory**

---

## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED MONITORING_EFFORT_REVIEWS STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
mkdir -p markers/state-verification && touch "markers/state-verification/state_rules_read_orchestrator_MONITORING_EFFORT_REVIEWS-$(date +%Y%m%d-%H%M%S)"
echo "$(date +%s) - Rules read and acknowledged for MONITORING_EFFORT_REVIEWS" > "markers/state-verification/state_rules_read_orchestrator_MONITORING_EFFORT_REVIEWS-$(date +%Y%m%d-%H%M%S)"
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

## 📋 PRIMARY DIRECTIVES FOR MONITORING_EFFORT_REVIEWS STATE

### Core Mandatory Rules (ALL orchestrator states must have these):

1. **🚨🚨🚨 R006** - ORCHESTRATOR NEVER WRITES CODE OR PERFORMS FILE OPERATIONS (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`
   - Criticality: BLOCKING - Automatic termination, 0% grade
   - Summary: NEVER write, copy, move, or manipulate ANY code files - delegate ALL to agents
   - **CRITICAL**: Copying files is NOT infrastructure - it's implementation work!

2. **🔴🔴🔴 R287** - TODO PERSISTENCE COMPREHENSIVE (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-persistence-comprehensive.md`
   - Criticality: SUPREME - -20% to -100% penalty for violations
   - Summary: MUST save TODOs within 30s after write, every 10 messages, before transitions
   - **CRITICAL**: Commit and push within 60 seconds of saving

3. **🔴🔴🔴 R288** - STATE FILE UPDATE REQUIREMENTS (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
   - Criticality: SUPREME - State updates required for all transitions
   - Summary: MUST update orchestrator-state-v3.json before EVERY state transition
   - **CRITICAL**: Commit and push state changes immediately

### State-Specific Rules:

4. **⚠️⚠️⚠️ R233** - Active Monitoring Requirements
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R233-all-states-immediate-action.md`
   - Criticality: WARNING - Must actively monitor, not passively wait
   - Summary: Check progress every 5 messages, detect stalls/blockers

5. **🔴🔴🔴 R322** - Mandatory Checkpoints
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-before-state-transitions.md`
   - Criticality: SUPREME LAW - Must stop at checkpoints
   - Summary: STOP after state completion, await /continue-software-factory

## 🎯 STATE OBJECTIVES - MONITORING CODE REVIEWER EFFORT REVIEW PROGRESS

In the MONITORING_EFFORT_REVIEWS state, the ORCHESTRATOR is responsible for:

1. **Active Review Progress Monitoring**
   - Check each Code Reviewer's progress regularly
   - Track completion status for each effort review
   - Identify any BLOCKED reviewers
   - Monitor for timeout conditions

2. **Review Aggregation**
   - Count completed reviews
   - Track remaining work
   - Calculate completion percentage
   - Estimate time to completion

3. **Issue Detection and Response**
   - Detect blocked reviewers
   - Identify failed reviews
   - Document issues for resolution
   - Determine if ERROR_RECOVERY needed

4. **Completion Verification**
   - Verify all reviews complete
   - Check all review reports created
   - Confirm feedback documented
   - Prepare for next state transition

## 📝 REQUIRED ACTIONS

### Step 1: Initial Status Assessment
```bash
# Get initial status of all Code Reviewers
echo "📊 Initial Code Reviewer Status Check..."
echo "=================================="

# Get list of active code reviewer sessions from state file
CODE_REVIEWERS=$(jq -r '.spawned_agents.code_reviewers[]? | select(.state != "COMPLETE") | .effort_id' "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json")

if [ -z "$CODE_REVIEWERS" ]; then
    echo "⚠️ No active code reviewers found in state file"
    echo "Checking for review marker files as fallback..."

    # Fallback: scan for REVIEW-IN-PROGRESS markers
    find efforts -name "REVIEW-IN-PROGRESS.marker" 2>/dev/null | while read marker; do
        EFFORT_DIR=$(dirname "$marker")
        echo "Found review in progress: $EFFORT_DIR"
    done
fi

# Check each effort being reviewed
for EFFORT in $CODE_REVIEWERS; do
    echo ""
    echo "Checking ${EFFORT}..."

    # Determine effort directory
    EFFORT_DIR=$(find efforts -type d -name "$EFFORT" 2>/dev/null | head -1)

    if [ -z "$EFFORT_DIR" ]; then
        echo "  Status: ⚠️ EFFORT DIRECTORY NOT FOUND"
        continue
    fi

    # Check for review progress markers
    if [ -f "$EFFORT_DIR/REVIEW-COMPLETE.marker" ]; then
        echo "  Status: ✅ REVIEW COMPLETE"

        # Check for review report
        REVIEW_REPORT=$(find "$EFFORT_DIR/.software-factory" -name "CODE-REVIEW-REPORT--*.md" 2>/dev/null | head -1)
        if [ -n "$REVIEW_REPORT" ]; then
            echo "  Report: $(basename "$REVIEW_REPORT")"

            # Extract summary
            ISSUES_FOUND=$(grep -c "^###" "$REVIEW_REPORT" 2>/dev/null || echo "0")
            echo "  Issues Found: $ISSUES_FOUND"
        else
            echo "  Report: ⚠️ MISSING (R343 violation!)"
        fi

    elif [ -f "$EFFORT_DIR/REVIEW-BLOCKED.marker" ]; then
        echo "  Status: ❌ BLOCKED"

        # Get block reason
        if [ -f "$EFFORT_DIR/REVIEW-BLOCKED.marker" ]; then
            BLOCK_REASON=$(cat "$EFFORT_DIR/REVIEW-BLOCKED.marker")
            echo "  Reason: $BLOCK_REASON"
        fi

    elif [ -f "$EFFORT_DIR/REVIEW-IN-PROGRESS.marker" ]; then
        echo "  Status: ⏳ IN PROGRESS"

        # Check how long review has been running
        STARTED=$(stat -c %Y "$EFFORT_DIR/REVIEW-IN-PROGRESS.marker" 2>/dev/null || echo "0")
        CURRENT=$(date +%s)
        ELAPSED=$((CURRENT - STARTED))
        echo "  Elapsed: ${ELAPSED}s"

    else
        echo "  Status: ⏳ NOT STARTED (no markers found)"
    fi

    # Check git status for review branch
    if [ -d "$EFFORT_DIR/.git" ]; then
        cd "$EFFORT_DIR"
        CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null)
        echo "  Branch: $CURRENT_BRANCH"

        # Check for uncommitted review artifacts
        UNCOMMITTED=$(git status --porcelain 2>/dev/null | grep -c "CODE-REVIEW" || echo "0")
        if [ "$UNCOMMITTED" -gt 0 ]; then
            echo "  ⚠️ Warning: $UNCOMMITTED uncommitted review artifacts"
        fi
        cd - > /dev/null
    fi
done
```

### Step 2: Continuous Monitoring Loop (R233 Compliance)
```bash
# Monitor until all complete or timeout
MONITOR_INTERVAL=30  # Check every 30 seconds (R233)
MAX_MONITOR_TIME=3600  # Maximum 60 minutes for reviews
ELAPSED_TIME=0
ALL_COMPLETE=false
MESSAGE_COUNT=0

while [ $ELAPSED_TIME -lt $MAX_MONITOR_TIME ] && [ "$ALL_COMPLETE" = false ]; do
    echo ""
    echo "⏰ Monitor check $(date '+%H:%M:%S') - Elapsed: ${ELAPSED_TIME}s"
    echo "------------------------------------------------"

    COMPLETED_COUNT=0
    BLOCKED_COUNT=0
    IN_PROGRESS_COUNT=0
    TOTAL_COUNT=$(echo "$CODE_REVIEWERS" | wc -w)

    # Check each effort review
    for EFFORT in $CODE_REVIEWERS; do
        EFFORT_DIR=$(find efforts -type d -name "$EFFORT" 2>/dev/null | head -1)

        if [ -z "$EFFORT_DIR" ]; then
            echo "⚠️ ${EFFORT}: Directory not found"
            continue
        fi

        if [ -f "$EFFORT_DIR/REVIEW-COMPLETE.marker" ]; then
            ((COMPLETED_COUNT++))
            echo "✅ ${EFFORT}: Review Complete"

            # Verify review artifacts exist
            REVIEW_REPORT=$(find "$EFFORT_DIR/.software-factory" -name "CODE-REVIEW-REPORT--*.md" 2>/dev/null | head -1)
            if [ -n "$REVIEW_REPORT" ]; then
                ISSUES=$(grep -c "^###" "$REVIEW_REPORT" 2>/dev/null || echo "0")
                echo "   Report: $(basename "$REVIEW_REPORT") ($ISSUES issues)"
            else
                echo "   ⚠️ WARNING: Review complete but no report found!"
            fi

        elif [ -f "$EFFORT_DIR/REVIEW-BLOCKED.marker" ]; then
            ((BLOCKED_COUNT++))
            echo "❌ ${EFFORT}: BLOCKED"
            BLOCK_REASON=$(head -1 "$EFFORT_DIR/REVIEW-BLOCKED.marker" 2>/dev/null || echo "Unknown")
            echo "   Reason: $BLOCK_REASON"

        else
            ((IN_PROGRESS_COUNT++))
            echo "⏳ ${EFFORT}: In Progress"

            # Check for signs of progress
            if [ -f "$EFFORT_DIR/REVIEW-IN-PROGRESS.marker" ]; then
                STARTED=$(stat -c %Y "$EFFORT_DIR/REVIEW-IN-PROGRESS.marker" 2>/dev/null || echo "0")
                CURRENT=$(date +%s)
                REVIEW_ELAPSED=$((CURRENT - STARTED))

                # Warn if review taking too long
                if [ $REVIEW_ELAPSED -gt 1800 ]; then  # >30 minutes
                    echo "   ⚠️ WARNING: Review running for ${REVIEW_ELAPSED}s (>30min)"
                fi
            fi
        fi
    done

    # Summary
    echo ""
    echo "📊 Progress Summary:"
    echo "  Completed: ${COMPLETED_COUNT}/${TOTAL_COUNT}"
    echo "  In Progress: ${IN_PROGRESS_COUNT}/${TOTAL_COUNT}"
    echo "  Blocked: ${BLOCKED_COUNT}/${TOTAL_COUNT}"

    COMPLETION_PERCENT=$((COMPLETED_COUNT * 100 / TOTAL_COUNT))
    echo "  Completion: ${COMPLETION_PERCENT}%"

    # Check if all complete
    if [ $COMPLETED_COUNT -eq $TOTAL_COUNT ]; then
        echo "🎉 All reviews complete!"
        ALL_COMPLETE=true
        break
    fi

    # Check if any blocked
    if [ $BLOCKED_COUNT -gt 0 ]; then
        echo "⚠️ WARNING: ${BLOCKED_COUNT} reviewers blocked - may need intervention"
    fi

    # R233 compliance: Check every 5 messages
    ((MESSAGE_COUNT++))
    if [ $MESSAGE_COUNT -ge 5 ]; then
        echo "📝 R233: Active monitoring check #$((ELAPSED_TIME / MONITOR_INTERVAL))"
        MESSAGE_COUNT=0
    fi

    # Wait before next check
    sleep $MONITOR_INTERVAL
    ELAPSED_TIME=$((ELAPSED_TIME + MONITOR_INTERVAL))
done

# Check final status
if [ "$ALL_COMPLETE" = false ]; then
    if [ $BLOCKED_COUNT -gt 0 ]; then
        echo "❌ CRITICAL: Reviews blocked - ERROR_RECOVERY needed"
        NEXT_STATE="ERROR_RECOVERY"
        TRANSITION_REASON="Code reviewers blocked during effort reviews"
    else
        echo "⚠️ TIMEOUT: Not all reviews completed in time"
        NEXT_STATE="ERROR_RECOVERY"
        TRANSITION_REASON="Review timeout after ${MAX_MONITOR_TIME}s"
    fi
else
    echo "✅ All reviews completed successfully"
fi
```

### Step 2.5: Cleanup Completed Code Reviewer Agents (R610 - BLOCKING)
```bash
# R610: Agent Metadata Lifecycle Protocol - BLOCKING requirement
# Cleanup completed Code-Reviewer agents within 60 seconds

echo ""
echo "🧹 R610/R611: Cleaning up completed code reviewer agents..."
echo "============================================================"

# Find completed code-reviewer agents
COMPLETED_REVIEWERS=$(jq -r '
    .active_agents[] |
    select(.agent_type == "code-reviewer") |
    select(.state == "COMPLETE" or .state == "COMPLETED") |
    .agent_id
' orchestrator-state-v3.json)

if [ -z "$COMPLETED_REVIEWERS" ]; then
    echo "✅ R610: No completed reviewer agents to cleanup"
else
    REVIEWER_COUNT=$(echo "$COMPLETED_REVIEWERS" | wc -l)
    echo "📊 R610: Found $REVIEWER_COUNT completed code-reviewer agent(s)"

    # Run cleanup utility
    if bash tools/cleanup-completed-agents.sh; then
        echo "✅ R610: Code reviewer cleanup successful"
    else
        echo "❌ R610 VIOLATION: Cleanup failed!"
    fi
fi

echo "✅ R610/R611: Code reviewer cleanup complete"
```

---

### Step 3: Verify Review Quality
```bash
# For completed reviews, verify quality and completeness
echo ""
echo "🔍 Verifying Review Quality..."
echo "================================"

VERIFICATION_PASSED=true
REVIEWS_WITH_ISSUES=()
REVIEWS_NO_ISSUES=()

for EFFORT in $CODE_REVIEWERS; do
    EFFORT_DIR=$(find efforts -type d -name "$EFFORT" 2>/dev/null | head -1)

    if [ -z "$EFFORT_DIR" ]; then
        continue
    fi

    if [ -f "$EFFORT_DIR/REVIEW-COMPLETE.marker" ]; then
        echo ""
        echo "Verifying ${EFFORT}..."

        # Check for review report (R343: in .software-factory with timestamp)
        REVIEW_REPORT=$(find "$EFFORT_DIR/.software-factory" -name "CODE-REVIEW-REPORT--*.md" 2>/dev/null | head -1)

        if [ -z "$REVIEW_REPORT" ]; then
            echo "  ❌ FAIL: No review report found (R343 violation!)"
            VERIFICATION_PASSED=false
            continue
        fi

        echo "  ✅ Review report: $(basename "$REVIEW_REPORT")"

        # Count issues found
        ISSUE_COUNT=$(grep -c "^###" "$REVIEW_REPORT" 2>/dev/null || echo "0")
        echo "  Issues documented: $ISSUE_COUNT"

        if [ "$ISSUE_COUNT" -gt 0 ]; then
            REVIEWS_WITH_ISSUES+=("$EFFORT")

            # Check for FIX-INSTRUCTIONS (should be created for issues)
            FIX_INSTRUCTIONS=$(find "$EFFORT_DIR/.software-factory" -name "FIX-INSTRUCTIONS--*.md" 2>/dev/null | head -1)
            if [ -n "$FIX_INSTRUCTIONS" ]; then
                echo "  ✅ Fix instructions: $(basename "$FIX_INSTRUCTIONS")"
            else
                echo "  ⚠️ WARNING: Issues found but no FIX-INSTRUCTIONS created"
            fi
        else
            REVIEWS_NO_ISSUES+=("$EFFORT")
            echo "  ✅ No issues found - ready for integration"
        fi

        # Verify required sections in report
        REQUIRED_SECTIONS=("## Summary" "## Review Findings" "## Recommendations")
        for SECTION in "${REQUIRED_SECTIONS[@]}"; do
            if grep -q "$SECTION" "$REVIEW_REPORT"; then
                echo "  ✅ Section found: $SECTION"
            else
                echo "  ⚠️ WARNING: Missing section: $SECTION"
                VERIFICATION_PASSED=false
            fi
        done

        # Check if report was committed
        cd "$EFFORT_DIR"
        if git log -1 --pretty=format:"%s" | grep -q "review:"; then
            echo "  ✅ Review report committed to git"
        else
            echo "  ⚠️ WARNING: Review report not committed"
        fi
        cd - > /dev/null
    fi
done

if [ "$VERIFICATION_PASSED" = true ]; then
    echo ""
    echo "✅ All review verifications PASSED"
else
    echo ""
    echo "⚠️ Some verifications failed - review needed"
fi

# Determine next state based on results
if [ ${#REVIEWS_WITH_ISSUES[@]} -gt 0 ]; then
    echo ""
    echo "📋 Reviews with issues requiring fixes: ${#REVIEWS_WITH_ISSUES[@]}"
    echo "Next action: SPAWN_SW_ENGINEERS_FOR_FIXES"
    NEXT_STATE="SPAWN_SW_ENGINEERS_FOR_FIXES"
    TRANSITION_REASON="Reviews complete, ${#REVIEWS_WITH_ISSUES[@]} efforts need fixes"
elif [ ${#REVIEWS_NO_ISSUES[@]} -eq $TOTAL_COUNT ]; then
    echo ""
    echo "✅ All reviews passed with no issues"
    echo "Next action: Proceed to integration"
    NEXT_STATE="INTEGRATE_WAVE_EFFORTS"
    TRANSITION_REASON="All effort reviews passed, ready for integration"
fi
```

### Step 4: Create Completion Report
```bash
# Create comprehensive completion report
cat > "$CLAUDE_PROJECT_DIR/MONITORING-REVIEWS-REPORT.md" << EOF
# Code Review Monitoring - Final Report

## Monitoring Summary
- Start Time: $(date -d @$(($(date +%s) - ELAPSED_TIME)) '+%Y-%m-%d %H:%M:%S')
- End Time: $(date '+%Y-%m-%d %H:%M:%S')
- Total Duration: ${ELAPSED_TIME} seconds
- Monitor Checks: $((ELAPSED_TIME / MONITOR_INTERVAL))

## Review Status

### Completed Successfully: ${COMPLETED_COUNT}/${TOTAL_COUNT}
$(for effort in $CODE_REVIEWERS; do
    EFFORT_DIR=$(find efforts -type d -name "$effort" 2>/dev/null | head -1)
    if [ -f "$EFFORT_DIR/REVIEW-COMPLETE.marker" ]; then
        echo "- $effort: ✅ Complete"
    fi
done)

### In Progress: ${IN_PROGRESS_COUNT}/${TOTAL_COUNT}
$(for effort in $CODE_REVIEWERS; do
    EFFORT_DIR=$(find efforts -type d -name "$effort" 2>/dev/null | head -1)
    if [ -f "$EFFORT_DIR/REVIEW-IN-PROGRESS.marker" ] && [ ! -f "$EFFORT_DIR/REVIEW-COMPLETE.marker" ]; then
        echo "- $effort: ⏳ Still in progress"
    fi
done)

### Blocked: ${BLOCKED_COUNT}/${TOTAL_COUNT}
$(for effort in $CODE_REVIEWERS; do
    EFFORT_DIR=$(find efforts -type d -name "$effort" 2>/dev/null | head -1)
    if [ -f "$EFFORT_DIR/REVIEW-BLOCKED.marker" ]; then
        REASON=$(head -1 "$EFFORT_DIR/REVIEW-BLOCKED.marker" 2>/dev/null || echo "Unknown")
        echo "- $effort: ❌ Blocked - $REASON"
    fi
done)

## Review Results

### Reviews with Issues: ${#REVIEWS_WITH_ISSUES[@]}
$(for effort in "${REVIEWS_WITH_ISSUES[@]}"; do
    EFFORT_DIR=$(find efforts -type d -name "$effort" 2>/dev/null | head -1)
    REVIEW_REPORT=$(find "$EFFORT_DIR/.software-factory" -name "CODE-REVIEW-REPORT--*.md" 2>/dev/null | head -1)
    ISSUES=$(grep -c "^###" "$REVIEW_REPORT" 2>/dev/null || echo "0")
    echo "- $effort: $ISSUES issues found"
done)

### Reviews with No Issues: ${#REVIEWS_NO_ISSUES[@]}
$(for effort in "${REVIEWS_NO_ISSUES[@]}"; do
    echo "- $effort: ✅ No issues"
done)

## Verification Results
- All reviews complete: $([ "$ALL_COMPLETE" = true ] && echo "YES" || echo "NO")
- All reports generated: $([ "$VERIFICATION_PASSED" = true ] && echo "YES" || echo "NO")
- Reviews requiring fixes: ${#REVIEWS_WITH_ISSUES[@]}
- Reviews ready for integration: ${#REVIEWS_NO_ISSUES[@]}

## Next State Recommendation
**Proposed Next State:** $NEXT_STATE
**Reason:** $TRANSITION_REASON

## Issues Encountered
$(if [ $BLOCKED_COUNT -gt 0 ]; then
    echo "- $BLOCKED_COUNT reviewers blocked"
fi)
$(if [ "$ALL_COMPLETE" = false ] && [ $BLOCKED_COUNT -eq 0 ]; then
    echo "- Timeout after ${MAX_MONITOR_TIME}s"
fi)

## R233 Active Monitoring Compliance
- Monitor interval: ${MONITOR_INTERVAL}s
- Total checks performed: $((ELAPSED_TIME / MONITOR_INTERVAL))
- Active monitoring: ✅ COMPLIANT
EOF

echo "✅ Monitoring report created: MONITORING-REVIEWS-REPORT.md"
```

## ⚠️ CRITICAL REQUIREMENTS

### Active Monitoring (R233)
- Check progress every 30 seconds
- Cannot passively wait
- Document each monitoring cycle
- Check every 5 messages for stalls

### No Direct Intervention
- If reviewer is blocked, document it
- Do NOT try to fix review issues yourself
- Transition to ERROR_RECOVERY if needed

### Complete Verification
- All reviewers must complete
- All review reports must exist (R343)
- Verification must pass for success

## 🚫 FORBIDDEN ACTIONS

1. **Performing code reviews yourself** - R006 violation (orchestrator never writes code)
2. **Passive waiting without checks** - R233 violation
3. **Continuing with incomplete reviews** - Will cause integration failures
4. **Modifying Code Reviewer work** - They own their review artifacts
5. **Skipping verification** - Must ensure quality

## ✅ PROJECT_DONE CRITERIA

For successful transition to next state:
- [ ] All Code Reviewers report REVIEW_COMPLETE
- [ ] No reviewers in BLOCKED state
- [ ] All review reports created (R343 compliant)
- [ ] Verification checks pass
- [ ] Completion report created
- [ ] Next state determined from results

## 🔄 STATE TRANSITIONS

### Success Paths:

**If issues found:**
```
MONITORING_EFFORT_REVIEWS → SPAWN_SW_ENGINEERS_FOR_FIXES
```
- Reviews complete with issues
- Fix instructions created
- Engineers needed for fixes

**If no issues:**
```
MONITORING_EFFORT_REVIEWS → INTEGRATE_WAVE_EFFORTS
```
- All reviews passed
- No issues found
- Ready for integration

### Error Path:
```
MONITORING_EFFORT_REVIEWS → ERROR_RECOVERY
```
- Reviewers blocked
- Timeout occurred
- Verification failed

## 🎓 GRADING CRITERIA

You will be evaluated on:
1. **Active Monitoring** (35%)
   - Regular progress checks
   - R233 compliance

2. **Issue Detection** (25%)
   - Identifying blocks quickly
   - Catching failures

3. **Verification Quality** (25%)
   - Thorough verification
   - R343 artifact compliance

4. **Documentation** (15%)
   - Clear progress tracking
   - Complete final report

## 💡 TIPS FOR PROJECT_DONE

1. **Check frequently** - Every 30 seconds
2. **Document everything** - Show active monitoring
3. **Detect issues early** - Don't wait for timeout
4. **Verify thoroughly** - Ensure complete reviews with proper artifacts

Remember: You're the PROJECT MANAGER - monitor, verify, report!

## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-software-factory
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**

## 🔴🔴🔴 STATE COMPLETION CHECKLIST 🔴🔴🔴

**Execute these steps IN ORDER to properly complete MONITORING_EFFORT_REVIEWS:**

### ✅ Step 1: Complete State-Specific Work
**Refer to "Required Actions" sections above for state-specific tasks.**

Once all state work is complete, proceed to mandatory exit protocol:

---

### ✅ Step 2: Set Proposed Next State and Transition Reason
```bash
# Based on state work results, set variables for State Manager
# NEXT_STATE and TRANSITION_REASON already set in Step 3 above
echo "Proposed next state: $NEXT_STATE"
echo "Transition reason: $TRANSITION_REASON"
```

---

### ✅ Step 3: Spawn State Manager for State Transition (R288 - SUPREME LAW)
```bash
# State Manager handles ALL state file updates atomically
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "MONITORING_EFFORT_REVIEWS" \
  --proposed-next-state "$NEXT_STATE" \
  --transition-reason "$TRANSITION_REASON"

# State Manager will:
# 1. Validate the transition against state machine
# 2. Update all 4 state tracking locations atomically:
#    - orchestrator-state-v3.json
#    - orchestrator-state-demo.json
#    - .cascade-state-backup.json
#    - .orchestrator-state-v3.json
# 3. Commit and push all changes
# 4. Return control

echo "✅ State Manager completed transition"
```

---

### ✅ Step 4: Save TODOs (R287 - SUPREME LAW)
```bash
# Save TODO state before exit (R287 trigger)
save_todos "MONITORING_EFFORT_REVIEWS_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - monitoring effort reviews complete [R287]"; then
    echo "❌ ERROR: Failed to commit TODO files"
    echo "This is non-fatal but TODOs may be lost in compaction"
    echo "Proceeding with state execution..."
    # Don't exit - TODO commit failure is not fatal
fi

git push || echo "⚠️ WARNING: TODO push failed - committed locally"
echo "✅ TODOs saved and committed"
git push
echo "✅ TODOs saved and committed"
```

---

### ✅ Step 5: Output Continuation Flag (R405 - SUPREME LAW) ⚠️ MANDATORY
```bash
# Output continuation flag as LAST action (R405)
# Use TRUE for normal completion, FALSE only for catastrophic errors

echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"
```

**⚠️ THIS MUST BE THE ABSOLUTE LAST LINE OF OUTPUT BEFORE EXIT! ⚠️**

---

### ✅ Step 6: Stop Processing (R322 - SUPREME LAW)
```bash
# Stop for context preservation (R322)
echo "🛑 Stopping for context preservation - use /continue-software-factory to resume"
exit 0
```

---

## 🚨 CHECKLIST ENFORCEMENT 🚨

**Skipping ANY step in this checklist = FAILURE:**
- Missing Step 2: No next state = stuck forever
- Missing Step 3: No State Manager spawn = state machine broken (R288 violation, -100%)
- Missing Step 4: No TODO save = work lost (R287 violation, -20% to -100%)
- Missing Step 5: No continuation flag = automation stops (R405 violation, -100%)
- Missing Step 6: No exit = R322 violation (-100%)

**ALL 6 STEPS ARE MANDATORY - NO EXCEPTIONS**
**NOTE**: State file validation and commits now handled by State Manager (SF 3.0 pattern)


## ✅ EXIT CHECKLIST (R405 - Continuation Flag Protocol) 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

### 🚨 CRITICAL DISTINCTION: AGENT STOPS ≠ FACTORY STOPS 🚨

**TWO INDEPENDENT DECISIONS - DO NOT CONFUSE THEM:**

#### 1. Should Agent Stop Work? (R322 Technical Requirement)
- Agent completes current state
- Agent saves TODOs and commits state
- Agent exits with `exit 0` (preserves context)
- User runs /continue-software-factory to resume
- **This is NORMAL at checkpoints**

#### 2. Should Factory Continue? (R405 Operational Status)
- Even though agent stopped, can automation proceed?
- TRUE = Healthy completion, automation can continue
- FALSE = Catastrophic failure, must halt everything
- **R322 checkpoints = TRUE (99.9% of cases)**

### THE PATTERN AT R322 CHECKPOINTS (SF 3.0)

```bash
# 1. Complete state work
echo "✅ State work complete"

# 2. Set proposed next state
PROPOSED_NEXT_STATE="SPAWN_CODE_REVIEWER_FIX_PLAN"
TRANSITION_REASON="State work complete"

# 3. Spawn State Manager for state transition
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "MONITORING_EFFORT_REVIEWS" \
  --proposed-next-state "$PROPOSED_NEXT_STATE" \
  --transition-reason "$TRANSITION_REASON"
# State Manager updates all 4 state files atomically

# 4. Save TODOs
save_todos "R322_CHECKPOINT"

# 5. Factory continues (operational status)
echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"

# 6. Agent stops (technical requirement)
exit 0
```

**Both happen together! Agent stops AND factory continues!**

### WHEN TO USE EACH FLAG VALUE

**TRUE (99.9%):**
- ✅ R322 checkpoint reached
- ✅ State work completed successfully
- ✅ Ready for /continue-software-factory
- ✅ Waiting for user to continue (NORMAL)
- ✅ Reviews complete, next action determined

**FALSE (0.1%):**
- ❌ CATASTROPHIC unrecoverable error
- ❌ Data corruption spreading
- ❌ Critical security violation
- ❌ NOT for R322 checkpoints
- ❌ NOT for user review needs

### 🚨 MONITORING_EFFORT_REVIEWS STATE PATTERN - NORMAL TRANSITIONS 🚨

**Monitoring states transition to next actions automatically:**
```bash
# After monitoring completes
echo "✅ Monitoring complete, all reviewers finished work"

# Determine next action from results
if reviews_have_issues; then
    transition_to "SPAWN_SW_ENGINEERS_FOR_FIXES"
elif all_passed; then
    transition_to "INTEGRATE_WAVE_EFFORTS"
fi

# R322 checkpoint (if required by this transition)
echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"  # NORMAL operation!
exit 0  # If R322 checkpoint
```

**Why TRUE is correct:**
- Monitoring results drive automatic actions
- System knows what to do based on results
- **Review findings = Spawn fixes OR integrate = NORMAL!**

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**
