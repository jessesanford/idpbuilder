# Orchestrator - MONITORING_EFFORT_FIXES State Rules

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

**YOU HAVE ENTERED MONITORING_EFFORT_FIXES STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_MONITORING_EFFORT_FIXES
echo "$(date +%s) - Rules read and acknowledged for MONITORING_EFFORT_FIXES" > .state_rules_read_orchestrator_MONITORING_EFFORT_FIXES
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

## 📋 PRIMARY DIRECTIVES FOR MONITORING_EFFORT_FIXES STATE

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
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-checkpoints.md`
   - Criticality: SUPREME LAW - Must stop at checkpoints
   - Summary: STOP after state completion, await /continue-software-factory

6. **🔴🔴🔴 R355** - Production Ready Code Enforcement
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R355-production-ready-code-enforcement-supreme-law.md`
   - Criticality: SUPREME LAW - No stubs, mocks, or placeholders
   - Summary: ALL fixes must be production-ready code, no TODO markers

## 🎯 STATE OBJECTIVES - MONITORING SW ENGINEER BUG FIX PROGRESS

In the MONITORING_EFFORT_FIXES state, the ORCHESTRATOR is responsible for:

1. **Active Fix Progress Monitoring**
   - Check each SW Engineer's fix progress regularly
   - Track completion status for each effort requiring fixes
   - Identify any BLOCKED engineers
   - Monitor for timeout conditions

2. **Fix Aggregation**
   - Count completed fixes
   - Track remaining work
   - Calculate completion percentage
   - Estimate time to completion

3. **Issue Detection and Response**
   - Detect blocked engineers
   - Identify failed fix attempts
   - Document issues for resolution
   - Determine if ERROR_RECOVERY needed

4. **Completion Verification**
   - Verify all fixes complete
   - Check all bugs resolved
   - Confirm builds/tests pass
   - Prepare for re-review cycle

## 📝 REQUIRED ACTIONS

### Step 1: Initial Status Assessment
```bash
# Get initial status of all SW Engineers working on fixes
echo "📊 Initial SW Engineer Fix Status Check..."
echo "=================================="

# Get list of efforts requiring fixes from state file
EFFORTS_WITH_FIXES=$(jq -r '.spawned_agents.sw_engineers[]? | select(.state != "COMPLETE" and .task == "FIX_ISSUES") | .effort_id' "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json")

if [ -z "$EFFORTS_WITH_FIXES" ]; then
    echo "⚠️ No active fix efforts found in state file"
    echo "Checking for FIX-IN-PROGRESS markers as fallback..."

    # Fallback: scan for FIX-IN-PROGRESS markers
    find efforts -name "FIX-IN-PROGRESS.marker" 2>/dev/null | while read marker; do
        EFFORT_DIR=$(dirname "$marker")
        echo "Found fix in progress: $EFFORT_DIR"
    done
fi

# Check each effort being fixed
for EFFORT in $EFFORTS_WITH_FIXES; do
    echo ""
    echo "Checking ${EFFORT}..."

    # Determine effort directory
    EFFORT_DIR=$(find efforts -type d -name "$EFFORT" 2>/dev/null | head -1)

    if [ -z "$EFFORT_DIR" ]; then
        echo "  Status: ⚠️ EFFORT DIRECTORY NOT FOUND"
        continue
    fi

    # Check for fix progress markers
    if [ -f "$EFFORT_DIR/FIX-COMPLETE.marker" ]; then
        echo "  Status: ✅ FIXES COMPLETE"

        # Check for fix summary (R343: in .software-factory with timestamp)
        FIX_SUMMARY=$(find "$EFFORT_DIR/.software-factory" -name "FIX-SUMMARY--*.md" 2>/dev/null | head -1)
        if [ -n "$FIX_SUMMARY" ]; then
            echo "  Summary: $(basename "$FIX_SUMMARY")"

            # Extract fix count
            FIXES_APPLIED=$(grep -c "^- \[x\]" "$FIX_SUMMARY" 2>/dev/null || echo "0")
            echo "  Fixes Applied: $FIXES_APPLIED"
        else
            echo "  Summary: ⚠️ MISSING (R343 violation!)"
        fi

        # Check build/test status
        if [ -f "$EFFORT_DIR/BUILD-PASS.marker" ]; then
            echo "  Build: ✅ PASS"
        else
            echo "  Build: ⚠️ Unknown or failed"
        fi

        if [ -f "$EFFORT_DIR/TESTS-PASS.marker" ]; then
            echo "  Tests: ✅ PASS"
        else
            echo "  Tests: ⚠️ Unknown or failed"
        fi

    elif [ -f "$EFFORT_DIR/FIX-BLOCKED.marker" ]; then
        echo "  Status: ❌ BLOCKED"

        # Get block reason
        if [ -f "$EFFORT_DIR/FIX-BLOCKED.marker" ]; then
            BLOCK_REASON=$(cat "$EFFORT_DIR/FIX-BLOCKED.marker")
            echo "  Reason: $BLOCK_REASON"
        fi

    elif [ -f "$EFFORT_DIR/FIX-IN-PROGRESS.marker" ]; then
        echo "  Status: ⏳ IN PROGRESS"

        # Check how long fixes have been running
        STARTED=$(stat -c %Y "$EFFORT_DIR/FIX-IN-PROGRESS.marker" 2>/dev/null || echo "0")
        CURRENT=$(date +%s)
        ELAPSED=$((CURRENT - STARTED))
        echo "  Elapsed: ${ELAPSED}s"

        # Check for FIX-INSTRUCTIONS to see total fixes needed
        FIX_INSTRUCTIONS=$(find "$EFFORT_DIR/.software-factory" -name "FIX-INSTRUCTIONS--*.md" 2>/dev/null | head -1)
        if [ -n "$FIX_INSTRUCTIONS" ]; then
            TOTAL_FIXES=$(grep -c "^- \[ \]" "$FIX_INSTRUCTIONS" 2>/dev/null || echo "0")
            echo "  Fixes Required: $TOTAL_FIXES"
        fi

    else
        echo "  Status: ⏳ NOT STARTED (no markers found)"
    fi

    # Check git status for fix branch
    if [ -d "$EFFORT_DIR/.git" ]; then
        cd "$EFFORT_DIR"
        CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null)
        echo "  Branch: $CURRENT_BRANCH"

        # Check for uncommitted fixes
        UNCOMMITTED=$(git status --porcelain 2>/dev/null | wc -l)
        if [ "$UNCOMMITTED" -gt 0 ]; then
            echo "  ⚠️ Warning: $UNCOMMITTED uncommitted changes"
        fi

        # Check recent fix commits
        FIX_COMMITS=$(git log --oneline --grep="fix:" -5 2>/dev/null | wc -l)
        if [ "$FIX_COMMITS" -gt 0 ]; then
            echo "  Fix commits: $FIX_COMMITS (last 5 commits)"
        fi

        cd - > /dev/null
    fi
done
```

### Step 2: Continuous Monitoring Loop (R233 Compliance)
```bash
# Monitor until all fixes complete or timeout
MONITOR_INTERVAL=30  # Check every 30 seconds (R233)
MAX_MONITOR_TIME=7200  # Maximum 2 hours for fixes
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
    TOTAL_COUNT=$(echo "$EFFORTS_WITH_FIXES" | wc -w)

    # Check each effort fix
    for EFFORT in $EFFORTS_WITH_FIXES; do
        EFFORT_DIR=$(find efforts -type d -name "$EFFORT" 2>/dev/null | head -1)

        if [ -z "$EFFORT_DIR" ]; then
            echo "⚠️ ${EFFORT}: Directory not found"
            continue
        fi

        if [ -f "$EFFORT_DIR/FIX-COMPLETE.marker" ]; then
            ((COMPLETED_COUNT++))
            echo "✅ ${EFFORT}: Fixes Complete"

            # Verify fix artifacts exist
            FIX_SUMMARY=$(find "$EFFORT_DIR/.software-factory" -name "FIX-SUMMARY--*.md" 2>/dev/null | head -1)
            if [ -n "$FIX_SUMMARY" ]; then
                FIXES=$(grep -c "^- \[x\]" "$FIX_SUMMARY" 2>/dev/null || echo "0")
                echo "   Summary: $(basename "$FIX_SUMMARY") ($FIXES fixes)"
            else
                echo "   ⚠️ WARNING: Fixes complete but no summary found!"
            fi

            # Check build/test status
            BUILD_STATUS="Unknown"
            TEST_STATUS="Unknown"
            [ -f "$EFFORT_DIR/BUILD-PASS.marker" ] && BUILD_STATUS="✅ PASS"
            [ -f "$EFFORT_DIR/BUILD-FAIL.marker" ] && BUILD_STATUS="❌ FAIL"
            [ -f "$EFFORT_DIR/TESTS-PASS.marker" ] && TEST_STATUS="✅ PASS"
            [ -f "$EFFORT_DIR/TESTS-FAIL.marker" ] && TEST_STATUS="❌ FAIL"

            echo "   Build: $BUILD_STATUS | Tests: $TEST_STATUS"

        elif [ -f "$EFFORT_DIR/FIX-BLOCKED.marker" ]; then
            ((BLOCKED_COUNT++))
            echo "❌ ${EFFORT}: BLOCKED"
            BLOCK_REASON=$(head -1 "$EFFORT_DIR/FIX-BLOCKED.marker" 2>/dev/null || echo "Unknown")
            echo "   Reason: $BLOCK_REASON"

        else
            ((IN_PROGRESS_COUNT++))
            echo "⏳ ${EFFORT}: In Progress"

            # Check for signs of progress
            if [ -f "$EFFORT_DIR/FIX-IN-PROGRESS.marker" ]; then
                STARTED=$(stat -c %Y "$EFFORT_DIR/FIX-IN-PROGRESS.marker" 2>/dev/null || echo "0")
                CURRENT=$(date +%s)
                FIX_ELAPSED=$((CURRENT - STARTED))

                # Warn if fixes taking too long
                if [ $FIX_ELAPSED -gt 3600 ]; then  # >60 minutes
                    echo "   ⚠️ WARNING: Fixes running for ${FIX_ELAPSED}s (>60min)"
                fi

                # Check progress via git commits
                cd "$EFFORT_DIR" 2>/dev/null
                if [ -d .git ]; then
                    RECENT_COMMITS=$(git log --oneline --since="5 minutes ago" 2>/dev/null | wc -l)
                    if [ "$RECENT_COMMITS" -gt 0 ]; then
                        echo "   Progress: $RECENT_COMMITS commits in last 5 min"
                    fi
                fi
                cd - > /dev/null 2>&1
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
        echo "🎉 All fixes complete!"
        ALL_COMPLETE=true
        break
    fi

    # Check if any blocked
    if [ $BLOCKED_COUNT -gt 0 ]; then
        echo "⚠️ WARNING: ${BLOCKED_COUNT} engineers blocked - may need intervention"
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
        echo "❌ CRITICAL: Fixes blocked - ERROR_RECOVERY needed"
        NEXT_STATE="ERROR_RECOVERY"
        TRANSITION_REASON="SW Engineers blocked during bug fixes"
    else
        echo "⚠️ TIMEOUT: Not all fixes completed in time"
        NEXT_STATE="ERROR_RECOVERY"
        TRANSITION_REASON="Fix timeout after ${MAX_MONITOR_TIME}s"
    fi
else
    echo "✅ All fixes completed successfully"
fi
```

### Step 3: Verify Fix Quality (R355 Production-Ready Code)
```bash
# For completed fixes, verify quality and completeness
echo ""
echo "🔍 Verifying Fix Quality (R355 Compliance)..."
echo "================================"

VERIFICATION_PASSED=true
EFFORTS_READY_FOR_REREVIEW=()
EFFORTS_WITH_ISSUES=()

for EFFORT in $EFFORTS_WITH_FIXES; do
    EFFORT_DIR=$(find efforts -type d -name "$EFFORT" 2>/dev/null | head -1)

    if [ -z "$EFFORT_DIR" ]; then
        continue
    fi

    if [ -f "$EFFORT_DIR/FIX-COMPLETE.marker" ]; then
        echo ""
        echo "Verifying ${EFFORT}..."

        # Check for fix summary (R343: in .software-factory with timestamp)
        FIX_SUMMARY=$(find "$EFFORT_DIR/.software-factory" -name "FIX-SUMMARY--*.md" 2>/dev/null | head -1)

        if [ -z "$FIX_SUMMARY" ]; then
            echo "  ❌ FAIL: No fix summary found (R343 violation!)"
            VERIFICATION_PASSED=false
            EFFORTS_WITH_ISSUES+=("$EFFORT")
            continue
        fi

        echo "  ✅ Fix summary: $(basename "$FIX_SUMMARY")"

        # Count fixes applied
        FIXES_COMPLETED=$(grep -c "^- \[x\]" "$FIX_SUMMARY" 2>/dev/null || echo "0")
        FIXES_PENDING=$(grep -c "^- \[ \]" "$FIX_SUMMARY" 2>/dev/null || echo "0")

        echo "  Fixes completed: $FIXES_COMPLETED"
        if [ "$FIXES_PENDING" -gt 0 ]; then
            echo "  ⚠️ WARNING: $FIXES_PENDING fixes still pending!"
            VERIFICATION_PASSED=false
            EFFORTS_WITH_ISSUES+=("$EFFORT")
            continue
        fi

        cd "$EFFORT_DIR"

        # R355 Check: No TODO/FIXME/HACK markers in code
        echo "  Checking R355 compliance (no TODO/FIXME)..."
        TODO_COUNT=$(grep -r "TODO\|FIXME\|HACK" --include="*.go" --include="*.py" --include="*.js" pkg/ 2>/dev/null | grep -v "test" | wc -l)
        if [ "$TODO_COUNT" -gt 0 ]; then
            echo "  ❌ FAIL: Found $TODO_COUNT TODO/FIXME markers (R355 violation!)"
            VERIFICATION_PASSED=false
            EFFORTS_WITH_ISSUES+=("$EFFORT")
        else
            echo "  ✅ No TODO/FIXME markers found"
        fi

        # R355 Check: No stub implementations
        echo "  Checking for stub implementations..."
        STUB_COUNT=$(grep -r "not implemented\|stub\|mock" --include="*.go" --include="*.py" --include="*.js" pkg/ 2>/dev/null | grep -v "test" | wc -l)
        if [ "$STUB_COUNT" -gt 0 ]; then
            echo "  ❌ FAIL: Found $STUB_COUNT stub/mock implementations (R355 violation!)"
            VERIFICATION_PASSED=false
            EFFORTS_WITH_ISSUES+=("$EFFORT")
        else
            echo "  ✅ No stubs found"
        fi

        # Check build status
        if [ -f "BUILD-PASS.marker" ]; then
            echo "  ✅ Build: PASS"
        else
            echo "  ❌ Build: FAIL or unknown"
            VERIFICATION_PASSED=false
            EFFORTS_WITH_ISSUES+=("$EFFORT")
        fi

        # Check test status
        if [ -f "TESTS-PASS.marker" ]; then
            echo "  ✅ Tests: PASS"
        else
            echo "  ❌ Tests: FAIL or unknown"
            VERIFICATION_PASSED=false
            EFFORTS_WITH_ISSUES+=("$EFFORT")
        fi

        # Check if changes committed
        UNCOMMITTED=$(git status --porcelain 2>/dev/null | wc -l)
        if [ "$UNCOMMITTED" -gt 0 ]; then
            echo "  ⚠️ WARNING: $UNCOMMITTED uncommitted changes"
            VERIFICATION_PASSED=false
            EFFORTS_WITH_ISSUES+=("$EFFORT")
        else
            echo "  ✅ All changes committed"
        fi

        # If all checks passed, ready for re-review
        if [ "$TODO_COUNT" -eq 0 ] && [ "$STUB_COUNT" -eq 0 ] && \
           [ -f "BUILD-PASS.marker" ] && [ -f "TESTS-PASS.marker" ] && \
           [ "$UNCOMMITTED" -eq 0 ]; then
            echo "  ✅ READY FOR RE-REVIEW"
            EFFORTS_READY_FOR_REREVIEW+=("$EFFORT")
        fi

        cd - > /dev/null
    fi
done

if [ "$VERIFICATION_PASSED" = true ]; then
    echo ""
    echo "✅ All fix verifications PASSED"
else
    echo ""
    echo "⚠️ Some verifications failed - issues need resolution"
fi

# Determine next state based on results
if [ ${#EFFORTS_READY_FOR_REREVIEW[@]} -eq $TOTAL_COUNT ]; then
    echo ""
    echo "✅ All fixes verified and ready for re-review"
    echo "Next action: SPAWN_CODE_REVIEWERS_FOR_REREVIEW"
    NEXT_STATE="SPAWN_CODE_REVIEWERS_FOR_REREVIEW"
    TRANSITION_REASON="All fixes complete and verified, ready for re-review"
elif [ ${#EFFORTS_WITH_ISSUES[@]} -gt 0 ]; then
    echo ""
    echo "⚠️ ${#EFFORTS_WITH_ISSUES[@]} efforts have verification issues"
    echo "Next action: ERROR_RECOVERY for issue resolution"
    NEXT_STATE="ERROR_RECOVERY"
    TRANSITION_REASON="${#EFFORTS_WITH_ISSUES[@]} efforts failed verification"
fi
```

### Step 4: Create Completion Report
```bash
# Create comprehensive completion report
cat > "$CLAUDE_PROJECT_DIR/MONITORING-FIXES-REPORT.md" << EOF
# Bug Fix Monitoring - Final Report

## Monitoring Summary
- Start Time: $(date -d @$(($(date +%s) - ELAPSED_TIME)) '+%Y-%m-%d %H:%M:%S')
- End Time: $(date '+%Y-%m-%d %H:%M:%S')
- Total Duration: ${ELAPSED_TIME} seconds
- Monitor Checks: $((ELAPSED_TIME / MONITOR_INTERVAL))

## Fix Status

### Completed Successfully: ${COMPLETED_COUNT}/${TOTAL_COUNT}
$(for effort in $EFFORTS_WITH_FIXES; do
    EFFORT_DIR=$(find efforts -type d -name "$effort" 2>/dev/null | head -1)
    if [ -f "$EFFORT_DIR/FIX-COMPLETE.marker" ]; then
        FIX_SUMMARY=$(find "$EFFORT_DIR/.software-factory" -name "FIX-SUMMARY--*.md" 2>/dev/null | head -1)
        FIXES=$(grep -c "^- \[x\]" "$FIX_SUMMARY" 2>/dev/null || echo "0")
        echo "- $effort: ✅ Complete ($FIXES fixes applied)"
    fi
done)

### In Progress: ${IN_PROGRESS_COUNT}/${TOTAL_COUNT}
$(for effort in $EFFORTS_WITH_FIXES; do
    EFFORT_DIR=$(find efforts -type d -name "$effort" 2>/dev/null | head -1)
    if [ -f "$EFFORT_DIR/FIX-IN-PROGRESS.marker" ] && [ ! -f "$EFFORT_DIR/FIX-COMPLETE.marker" ]; then
        echo "- $effort: ⏳ Still in progress"
    fi
done)

### Blocked: ${BLOCKED_COUNT}/${TOTAL_COUNT}
$(for effort in $EFFORTS_WITH_FIXES; do
    EFFORT_DIR=$(find efforts -type d -name "$effort" 2>/dev/null | head -1)
    if [ -f "$EFFORT_DIR/FIX-BLOCKED.marker" ]; then
        REASON=$(head -1 "$EFFORT_DIR/FIX-BLOCKED.marker" 2>/dev/null || echo "Unknown")
        echo "- $effort: ❌ Blocked - $REASON"
    fi
done)

## Verification Results (R355 Production-Ready Code)

### Verified and Ready for Re-Review: ${#EFFORTS_READY_FOR_REREVIEW[@]}
$(for effort in "${EFFORTS_READY_FOR_REREVIEW[@]}"; do
    echo "- $effort: ✅ All checks passed"
done)

### Efforts with Verification Issues: ${#EFFORTS_WITH_ISSUES[@]}
$(for effort in "${EFFORTS_WITH_ISSUES[@]}"; do
    echo "- $effort: ⚠️ Verification issues found"
done)

## Quality Checks Performed
- R355 compliance (no TODO/FIXME markers): $([ "$TODO_COUNT" -eq 0 ] && echo "✅ PASS" || echo "⚠️ ISSUES")
- No stub implementations: $([ "$STUB_COUNT" -eq 0 ] && echo "✅ PASS" || echo "⚠️ ISSUES")
- Build status: $([ -f "BUILD-PASS.marker" ] && echo "✅ PASS" || echo "⚠️ FAIL")
- Test status: $([ -f "TESTS-PASS.marker" ] && echo "✅ PASS" || echo "⚠️ FAIL")
- All changes committed: $([ "$UNCOMMITTED" -eq 0 ] && echo "✅ YES" || echo "⚠️ NO")

## Next State Recommendation
**Proposed Next State:** $NEXT_STATE
**Reason:** $TRANSITION_REASON

## Issues Encountered
$(if [ $BLOCKED_COUNT -gt 0 ]; then
    echo "- $BLOCKED_COUNT engineers blocked during fixes"
fi)
$(if [ "$ALL_COMPLETE" = false ] && [ $BLOCKED_COUNT -eq 0 ]; then
    echo "- Timeout after ${MAX_MONITOR_TIME}s"
fi)
$(if [ ${#EFFORTS_WITH_ISSUES[@]} -gt 0 ]; then
    echo "- ${#EFFORTS_WITH_ISSUES[@]} efforts failed verification"
fi)

## R233 Active Monitoring Compliance
- Monitor interval: ${MONITOR_INTERVAL}s
- Total checks performed: $((ELAPSED_TIME / MONITOR_INTERVAL))
- Active monitoring: ✅ COMPLIANT
EOF

echo "✅ Monitoring report created: MONITORING-FIXES-REPORT.md"
```

## ⚠️ CRITICAL REQUIREMENTS

### Active Monitoring (R233)
- Check progress every 30 seconds
- Cannot passively wait
- Document each monitoring cycle
- Check every 5 messages for stalls

### R355 Production-Ready Code Enforcement
- ALL fixes must be production-ready
- No TODO/FIXME markers allowed
- No stub implementations
- Builds and tests must pass

### No Direct Intervention
- If engineer is blocked, document it
- Do NOT try to fix bugs yourself
- Transition to ERROR_RECOVERY if needed

### Complete Verification
- All engineers must complete
- All fix summaries must exist (R343)
- All R355 checks must pass
- Builds and tests must pass

## 🚫 FORBIDDEN ACTIONS

1. **Writing fix code yourself** - R006 violation (orchestrator never writes code)
2. **Passive waiting without checks** - R233 violation
3. **Continuing with incomplete fixes** - Will cause re-review failures
4. **Modifying SW Engineer work** - They own their fix branches
5. **Skipping R355 verification** - Production code requirement

## ✅ PROJECT_DONE CRITERIA

For successful transition to next state:
- [ ] All SW Engineers report FIX_COMPLETE
- [ ] No engineers in BLOCKED state
- [ ] All fix summaries created (R343 compliant)
- [ ] R355 checks pass (no TODO/stubs)
- [ ] Builds and tests pass
- [ ] Verification checks pass
- [ ] Completion report created
- [ ] Next state determined from results

## 🔄 STATE TRANSITIONS

### Success Path:
```
MONITORING_EFFORT_FIXES → SPAWN_CODE_REVIEWERS_FOR_REREVIEW
```
- All fixes complete
- R355 verification passed
- Ready for re-review cycle

### Error Paths:
```
MONITORING_EFFORT_FIXES → ERROR_RECOVERY
```
- Engineers blocked
- Timeout occurred
- Verification failed
- Builds/tests failed

## 🎓 GRADING CRITERIA

You will be evaluated on:
1. **Active Monitoring** (30%)
   - Regular progress checks
   - R233 compliance

2. **R355 Verification** (30%)
   - Production-ready code checks
   - No TODO/stub detection
   - Build/test validation

3. **Issue Detection** (20%)
   - Identifying blocks quickly
   - Catching failures

4. **Documentation** (20%)
   - Clear progress tracking
   - Complete final report

## 💡 TIPS FOR PROJECT_DONE

1. **Check frequently** - Every 30 seconds
2. **Verify R355** - No shortcuts allowed
3. **Document everything** - Show active monitoring
4. **Detect issues early** - Don't wait for timeout

Remember: You're the PROJECT MANAGER - monitor, verify, enforce quality!

## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-software-factory
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**

## 🔴🔴🔴 STATE COMPLETION CHECKLIST 🔴🔴🔴

**Execute these steps IN ORDER to properly complete MONITORING_EFFORT_FIXES:**

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
  --current-state "MONITORING_EFFORT_FIXES" \
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
save_todos "MONITORING_EFFORT_FIXES_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - monitoring effort fixes complete [R287]"; then
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
PROPOSED_NEXT_STATE="SPAWN_CODE_REVIEWERS_EFFORT_REVIEW"
TRANSITION_REASON="State work complete"

# 3. Spawn State Manager for state transition
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "MONITORING_EFFORT_FIXES" \
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
- ✅ Fixes complete, re-review needed

**FALSE (0.1%):**
- ❌ CATASTROPHIC unrecoverable error
- ❌ Data corruption spreading
- ❌ Critical security violation
- ❌ NOT for R322 checkpoints
- ❌ NOT for user review needs

### 🚨 MONITORING_EFFORT_FIXES STATE PATTERN - NORMAL TRANSITIONS 🚨

**Monitoring states transition to next actions automatically:**
```bash
# After monitoring completes
echo "✅ Monitoring complete, all engineers finished fixes"

# Determine next action from results
if all_fixes_verified; then
    transition_to "SPAWN_CODE_REVIEWERS_FOR_REREVIEW"
elif verification_failed; then
    transition_to "ERROR_RECOVERY"
fi

# R322 checkpoint (if required by this transition)
echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"  # NORMAL operation!
exit 0  # If R322 checkpoint
```

**Why TRUE is correct:**
- Monitoring results drive automatic actions
- System knows what to do based on results
- **Fixes verified = Re-review needed = NORMAL!**

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**
