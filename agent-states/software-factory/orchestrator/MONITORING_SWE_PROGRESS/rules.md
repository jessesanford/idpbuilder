# Orchestrator - MONITORING_SWE_PROGRESS State Rules

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

**YOU HAVE ENTERED MONITORING_SWE_PROGRESS STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
mkdir -p markers/state-verification && touch "markers/state-verification/state_rules_read_orchestrator_MONITORING_SWE_PROGRESS-$(date +%Y%m%d-%H%M%S)"
echo "$(date +%s) - Rules read and acknowledged for MONITORING_SWE_PROGRESS" > "markers/state-verification/state_rules_read_orchestrator_MONITORING_SWE_PROGRESS-$(date +%Y%m%d-%H%M%S)"
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

## 📋 PRIMARY DIRECTIVES FOR MONITORING_SWE_PROGRESS STATE

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

## 🎯 STATE OBJECTIVES - MONITORING SW ENGINEER IMPLEMENTATION PROGRESS

In the MONITORING_SWE_PROGRESS state, the ORCHESTRATOR is responsible for:

1. **Active Implementation Progress Monitoring**
   - Check each SW Engineer's implementation progress regularly
   - Track completion status for each effort
   - Identify any BLOCKED engineers
   - Monitor for timeout conditions

2. **Progress Aggregation**
   - Count completed implementations
   - Track remaining work
   - Calculate completion percentage
   - Estimate time to completion

3. **Issue Detection and Response**
   - Detect blocked engineers
   - Identify failed implementations
   - Document issues for resolution
   - Determine if ERROR_RECOVERY needed

4. **Completion Verification**
   - Verify all implementations complete
   - Check all IMPLEMENTATION-COMPLETE markers
   - Confirm code committed and pushed
   - Prepare for next state transition (typically code review)

## 📝 REQUIRED ACTIONS

### Step 1: Initial Status Assessment
```bash
# Get initial status of all SW Engineers
echo "📊 Initial SW Engineer Status Check..."
echo "=================================="

# Get list of active SW engineer sessions from state file
SW_ENGINEERS=$(jq -r '.spawned_agents.sw_engineers[]? | select(.state != "COMPLETE") | .effort_id' "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json")

if [ -z "$SW_ENGINEERS" ]; then
    echo "⚠️ No active SW engineers found in state file"
    echo "Checking for IMPLEMENTATION-IN-PROGRESS markers as fallback..."

    # Fallback: scan for implementation markers
    find efforts -name "IMPLEMENTATION-IN-PROGRESS.marker" 2>/dev/null | while read marker; do
        EFFORT_DIR=$(dirname "$marker")
        echo "Found implementation in progress: $EFFORT_DIR"
    done
fi

# Check each effort being implemented
for EFFORT in $SW_ENGINEERS; do
    echo ""
    echo "Checking ${EFFORT}..."

    # Determine effort directory
    EFFORT_DIR=$(find efforts -type d -name "$EFFORT" 2>/dev/null | head -1)

    if [ -z "$EFFORT_DIR" ]; then
        echo "  Status: ⚠️ EFFORT DIRECTORY NOT FOUND"
        continue
    fi

    # Check for implementation progress markers
    if [ -f "$EFFORT_DIR/IMPLEMENTATION-COMPLETE.marker" ]; then
        echo "  Status: ✅ IMPLEMENTATION COMPLETE"

        # Check for work log (R343: in .software-factory with timestamp)
        WORK_LOG=$(find "$EFFORT_DIR/.software-factory" -name "work-log--*.md" 2>/dev/null | head -1)
        if [ -n "$WORK_LOG" ]; then
            echo "  Work Log: $(basename "$WORK_LOG")"

            # Extract implementation summary
            LINES_ADDED=$(grep "Lines added:" "$WORK_LOG" 2>/dev/null | tail -1 | awk '{print $3}')
            if [ -n "$LINES_ADDED" ]; then
                echo "  Lines Added: $LINES_ADDED"
            fi
        else
            echo "  Work Log: ⚠️ MISSING (R343 violation!)"
        fi

    elif [ -f "$EFFORT_DIR/IMPLEMENTATION-BLOCKED.marker" ]; then
        echo "  Status: ❌ BLOCKED"

        # Get block reason
        if [ -f "$EFFORT_DIR/IMPLEMENTATION-BLOCKED.marker" ]; then
            BLOCK_REASON=$(cat "$EFFORT_DIR/IMPLEMENTATION-BLOCKED.marker")
            echo "  Reason: $BLOCK_REASON"
        fi

    elif [ -f "$EFFORT_DIR/IMPLEMENTATION-IN-PROGRESS.marker" ]; then
        echo "  Status: ⏳ IN PROGRESS"

        # Check how long implementation has been running
        STARTED=$(stat -c %Y "$EFFORT_DIR/IMPLEMENTATION-IN-PROGRESS.marker" 2>/dev/null || echo "0")
        CURRENT=$(date +%s)
        ELAPSED=$((CURRENT - STARTED))
        echo "  Elapsed: ${ELAPSED}s"

    else
        echo "  Status: ⏳ NOT STARTED (no markers found)"
    fi

    # Check git status for implementation branch
    if [ -d "$EFFORT_DIR/.git" ]; then
        cd "$EFFORT_DIR"
        CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null)
        echo "  Branch: $CURRENT_BRANCH"

        # Check for uncommitted implementation code
        UNCOMMITTED=$(git status --porcelain 2>/dev/null | wc -l)
        if [ "$UNCOMMITTED" -gt 0 ]; then
            echo "  ⚠️ Warning: $UNCOMMITTED uncommitted files"
        fi

        # Check recent implementation commits
        IMPL_COMMITS=$(git log --oneline --grep="feat:\|impl:" -10 2>/dev/null | wc -l)
        if [ "$IMPL_COMMITS" -gt 0 ]; then
            echo "  Implementation commits: $IMPL_COMMITS (last 10 commits)"
        fi

        cd - > /dev/null
    fi
done
```

### Step 2: Continuous Monitoring Loop (R233 Compliance)
```bash
# Monitor until all complete or timeout
MONITOR_INTERVAL=30  # Check every 30 seconds (R233)
MAX_MONITOR_TIME=14400  # Maximum 4 hours for implementation
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
    TOTAL_COUNT=$(echo "$SW_ENGINEERS" | wc -w)

    # Check each effort implementation
    for EFFORT in $SW_ENGINEERS; do
        EFFORT_DIR=$(find efforts -type d -name "$EFFORT" 2>/dev/null | head -1)

        if [ -z "$EFFORT_DIR" ]; then
            echo "⚠️ ${EFFORT}: Directory not found"
            continue
        fi

        if [ -f "$EFFORT_DIR/IMPLEMENTATION-COMPLETE.marker" ]; then
            ((COMPLETED_COUNT++))
            echo "✅ ${EFFORT}: Implementation Complete"

            # Verify implementation artifacts exist
            WORK_LOG=$(find "$EFFORT_DIR/.software-factory" -name "work-log--*.md" 2>/dev/null | head -1)
            if [ -n "$WORK_LOG" ]; then
                echo "   Work Log: $(basename "$WORK_LOG")"
            else
                echo "   ⚠️ WARNING: Implementation complete but no work log!"
            fi

        elif [ -f "$EFFORT_DIR/IMPLEMENTATION-BLOCKED.marker" ]; then
            ((BLOCKED_COUNT++))
            echo "❌ ${EFFORT}: BLOCKED"
            BLOCK_REASON=$(head -1 "$EFFORT_DIR/IMPLEMENTATION-BLOCKED.marker" 2>/dev/null || echo "Unknown")
            echo "   Reason: $BLOCK_REASON"

        else
            ((IN_PROGRESS_COUNT++))
            echo "⏳ ${EFFORT}: In Progress"

            # Check for signs of progress
            if [ -f "$EFFORT_DIR/IMPLEMENTATION-IN-PROGRESS.marker" ]; then
                STARTED=$(stat -c %Y "$EFFORT_DIR/IMPLEMENTATION-IN-PROGRESS.marker" 2>/dev/null || echo "0")
                CURRENT=$(date +%s)
                IMPL_ELAPSED=$((CURRENT - STARTED))

                # Warn if implementation taking too long
                if [ $IMPL_ELAPSED -gt 7200 ]; then  # >2 hours
                    echo "   ⚠️ WARNING: Implementation running for ${IMPL_ELAPSED}s (>2hrs)"
                fi

                # Check progress via git commits
                cd "$EFFORT_DIR" 2>/dev/null
                if [ -d .git ]; then
                    RECENT_COMMITS=$(git log --oneline --since="10 minutes ago" 2>/dev/null | wc -l)
                    if [ "$RECENT_COMMITS" -gt 0 ]; then
                        echo "   Progress: $RECENT_COMMITS commits in last 10 min"
                    else
                        echo "   ⚠️ No recent commits (possible stall?)"
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
        echo "🎉 All implementations complete!"
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
        echo "❌ CRITICAL: Implementations blocked - ERROR_RECOVERY needed"
        NEXT_STATE="ERROR_RECOVERY"
        TRANSITION_REASON="SW Engineers blocked during implementation"
    else
        echo "⚠️ TIMEOUT: Not all implementations completed in time"
        NEXT_STATE="ERROR_RECOVERY"
        TRANSITION_REASON="Implementation timeout after ${MAX_MONITOR_TIME}s"
    fi
else
    echo "✅ All implementations completed successfully"
fi
```

### Step 2.5: Cleanup Completed Agents (R610 - BLOCKING)
```bash
# R610: Agent Metadata Lifecycle Protocol - BLOCKING requirement
# R611: Active Agents Cleanup Protocol - WARNING requirement
#
# MUST cleanup completed SW-Engineer agents within 60 seconds of detection
# This prevents active_agents bloat and maintains state file performance

echo ""
echo "🧹 R610/R611: Cleaning up completed agents..."
echo "=============================================="

# Find completed agents in active_agents
COMPLETED_AGENTS=$(jq -r '
    .active_agents[] |
    select(.agent_type == "sw-engineer") |
    select(.state == "COMPLETE" or .state == "COMPLETED") |
    .agent_id
' orchestrator-state-v3.json)

if [ -z "$COMPLETED_AGENTS" ]; then
    echo "✅ R610: No completed agents to cleanup"
else
    COMPLETED_COUNT=$(echo "$COMPLETED_AGENTS" | wc -l)
    echo "📊 R610: Found $COMPLETED_COUNT completed SW-Engineer agent(s)"

    # Run cleanup utility (implements R610 + R612)
    if bash tools/cleanup-completed-agents.sh; then
        echo "✅ R610: Cleanup successful - agents moved to agents_history"
    else
        echo "❌ R610 VIOLATION: Cleanup failed!"
        echo "This is a BLOCKING violation - performance will degrade"
        # Don't fail state, but log the violation
    fi

    # Verify no completed agents remain (R611 validation)
    REMAINING=$(jq '
        [.active_agents[] |
         select(.agent_type == "sw-engineer") |
         select(.state == "COMPLETE" or .state == "COMPLETED")] |
        length
    ' orchestrator-state-v3.json)

    if [ "$REMAINING" -gt 0 ]; then
        echo "❌ R611 VIOLATION: $REMAINING completed agents still in active_agents"
    else
        echo "✅ R611: Active agents array clean (only active agents remain)"
    fi
fi

# R613: Monitor state file size
STATE_FILE_SIZE=$(wc -c < orchestrator-state-v3.json)
STATE_FILE_KB=$((STATE_FILE_SIZE / 1024))

echo "📊 R613: State file size: ${STATE_FILE_KB}KB"

if [ "$STATE_FILE_SIZE" -gt 1048576 ]; then
    echo "❌ R613 CRITICAL: State file >1MB (${STATE_FILE_KB}KB)"
    echo "Immediate cleanup required!"
elif [ "$STATE_FILE_SIZE" -gt 512000 ]; then
    echo "⚠️  R613 WARNING: State file >500KB (${STATE_FILE_KB}KB)"
    echo "Cleanup recommended"
else
    echo "✅ R613: State file size within targets"
fi

echo "✅ R610/R611/R613: Agent cleanup and size monitoring complete"
```

**R610/R611/R613 Integration Notes:**
- R610 defines automatic cleanup timing (within 60s of detection)
- R611 defines what "active" means (not COMPLETE/COMPLETED)
- R612 defines agents_history schema (used by cleanup utility)
- R613 monitors total state file size after cleanup
- This monitoring state is the PRIMARY cleanup point for SW-Engineer agents
- Cleanup happens automatically during monitoring loop
- Boundary states (COMPLETE_WAVE, COMPLETE_PHASE) provide safety net validation

---

### Step 3: Verify Implementation Quality
```bash
# For completed implementations, verify quality and completeness
echo ""
echo "🔍 Verifying Implementation Quality..."
echo "================================"

VERIFICATION_PASSED=true
IMPLEMENTATIONS_READY_FOR_REVIEW=()
IMPLEMENTATIONS_WITH_ISSUES=()

for EFFORT in $SW_ENGINEERS; do
    EFFORT_DIR=$(find efforts -type d -name "$EFFORT" 2>/dev/null | head -1)

    if [ -z "$EFFORT_DIR" ]; then
        continue
    fi

    if [ -f "$EFFORT_DIR/IMPLEMENTATION-COMPLETE.marker" ]; then
        echo ""
        echo "Verifying ${EFFORT}..."

        # Check for work log (R343: in .software-factory with timestamp)
        WORK_LOG=$(find "$EFFORT_DIR/.software-factory" -name "work-log--*.md" 2>/dev/null | head -1)

        if [ -z "$WORK_LOG" ]; then
            echo "  ❌ FAIL: No work log found (R343 violation!)"
            VERIFICATION_PASSED=false
            IMPLEMENTATIONS_WITH_ISSUES+=("$EFFORT")
            continue
        fi

        echo "  ✅ Work log: $(basename "$WORK_LOG")"

        cd "$EFFORT_DIR"

        # Check if all changes are committed
        UNCOMMITTED=$(git status --porcelain 2>/dev/null | wc -l)
        if [ "$UNCOMMITTED" -gt 0 ]; then
            echo "  ⚠️ WARNING: $UNCOMMITTED uncommitted files"
            VERIFICATION_PASSED=false
            IMPLEMENTATIONS_WITH_ISSUES+=("$EFFORT")
        else
            echo "  ✅ All changes committed"
        fi

        # Check if branch is pushed
        LOCAL_COMMIT=$(git rev-parse HEAD 2>/dev/null)
        REMOTE_COMMIT=$(git rev-parse @{u} 2>/dev/null || echo "no-remote")

        if [ "$LOCAL_COMMIT" = "$REMOTE_COMMIT" ]; then
            echo "  ✅ Branch is pushed"
        else
            echo "  ⚠️ WARNING: Local branch ahead of remote (not pushed)"
            VERIFICATION_PASSED=false
            IMPLEMENTATIONS_WITH_ISSUES+=("$EFFORT")
        fi

        # Check for implementation plan (should exist from planning)
        IMPL_PLAN=$(find ".software-factory" -name "IMPLEMENTATION-PLAN--*.md" 2>/dev/null | head -1)
        if [ -n "$IMPL_PLAN" ]; then
            echo "  ✅ Implementation plan: $(basename "$IMPL_PLAN")"
        else
            echo "  ⚠️ WARNING: No implementation plan found"
        fi

        # Check pkg/ directory exists (where code should be)
        if [ -d "pkg" ]; then
            CODE_FILES=$(find pkg -name "*.go" -o -name "*.py" -o -name "*.js" 2>/dev/null | wc -l)
            echo "  ✅ Implementation files: $CODE_FILES in pkg/"
        else
            echo "  ⚠️ WARNING: No pkg/ directory found"
            VERIFICATION_PASSED=false
            IMPLEMENTATIONS_WITH_ISSUES+=("$EFFORT")
        fi

        # If all checks passed, ready for review
        if [ "$UNCOMMITTED" -eq 0 ] && [ "$LOCAL_COMMIT" = "$REMOTE_COMMIT" ] && [ -d "pkg" ]; then
            echo "  ✅ READY FOR CODE REVIEW"
            IMPLEMENTATIONS_READY_FOR_REVIEW+=("$EFFORT")
        fi

        cd - > /dev/null
    fi
done

if [ "$VERIFICATION_PASSED" = true ]; then
    echo ""
    echo "✅ All implementation verifications PASSED"
else
    echo ""
    echo "⚠️ Some verifications failed - issues need resolution"
fi

# Determine next state based on results
if [ ${#IMPLEMENTATIONS_READY_FOR_REVIEW[@]} -eq $TOTAL_COUNT ]; then
    echo ""
    echo "✅ All implementations verified and ready for code review"
    echo "Next action: SPAWN_CODE_REVIEWERS_FOR_EFFORT_REVIEW"
    NEXT_STATE="SPAWN_CODE_REVIEWERS_FOR_EFFORT_REVIEW"
    TRANSITION_REASON="All implementations complete and verified, ready for code review"
elif [ ${#IMPLEMENTATIONS_WITH_ISSUES[@]} -gt 0 ]; then
    echo ""
    echo "⚠️ ${#IMPLEMENTATIONS_WITH_ISSUES[@]} implementations have verification issues"
    echo "Next action: ERROR_RECOVERY for issue resolution"
    NEXT_STATE="ERROR_RECOVERY"
    TRANSITION_REASON="${#IMPLEMENTATIONS_WITH_ISSUES[@]} implementations failed verification"
fi
```

### Step 4: Create Completion Report
```bash
# Create comprehensive completion report
cat > "$CLAUDE_PROJECT_DIR/MONITORING-IMPLEMENTATION-REPORT.md" << EOF
# Implementation Monitoring - Final Report

## Monitoring Summary
- Start Time: $(date -d @$(($(date +%s) - ELAPSED_TIME)) '+%Y-%m-%d %H:%M:%S')
- End Time: $(date '+%Y-%m-%d %H:%M:%S')
- Total Duration: ${ELAPSED_TIME} seconds ($(echo "scale=2; $ELAPSED_TIME / 3600" | bc) hours)
- Monitor Checks: $((ELAPSED_TIME / MONITOR_INTERVAL))

## Implementation Status

### Completed Successfully: ${COMPLETED_COUNT}/${TOTAL_COUNT}
$(for effort in $SW_ENGINEERS; do
    EFFORT_DIR=$(find efforts -type d -name "$effort" 2>/dev/null | head -1)
    if [ -f "$EFFORT_DIR/IMPLEMENTATION-COMPLETE.marker" ]; then
        echo "- $effort: ✅ Complete"
    fi
done)

### In Progress: ${IN_PROGRESS_COUNT}/${TOTAL_COUNT}
$(for effort in $SW_ENGINEERS; do
    EFFORT_DIR=$(find efforts -type d -name "$effort" 2>/dev/null | head -1)
    if [ -f "$EFFORT_DIR/IMPLEMENTATION-IN-PROGRESS.marker" ] && [ ! -f "$EFFORT_DIR/IMPLEMENTATION-COMPLETE.marker" ]; then
        echo "- $effort: ⏳ Still in progress"
    fi
done)

### Blocked: ${BLOCKED_COUNT}/${TOTAL_COUNT}
$(for effort in $SW_ENGINEERS; do
    EFFORT_DIR=$(find efforts -type d -name "$effort" 2>/dev/null | head -1)
    if [ -f "$EFFORT_DIR/IMPLEMENTATION-BLOCKED.marker" ]; then
        REASON=$(head -1 "$EFFORT_DIR/IMPLEMENTATION-BLOCKED.marker" 2>/dev/null || echo "Unknown")
        echo "- $effort: ❌ Blocked - $REASON"
    fi
done)

## Verification Results

### Ready for Code Review: ${#IMPLEMENTATIONS_READY_FOR_REVIEW[@]}
$(for effort in "${IMPLEMENTATIONS_READY_FOR_REVIEW[@]}"; do
    echo "- $effort: ✅ All checks passed"
done)

### Implementations with Issues: ${#IMPLEMENTATIONS_WITH_ISSUES[@]}
$(for effort in "${IMPLEMENTATIONS_WITH_ISSUES[@]}"; do
    echo "- $effort: ⚠️ Verification issues found"
done)

## Quality Checks Performed
- Work logs exist (R343): $([ -n "$WORK_LOG" ] && echo "✅ YES" || echo "⚠️ NO")
- All changes committed: $([ "$UNCOMMITTED" -eq 0 ] && echo "✅ YES" || echo "⚠️ NO")
- Branches pushed to remote: $([ "$LOCAL_COMMIT" = "$REMOTE_COMMIT" ] && echo "✅ YES" || echo "⚠️ NO")
- Implementation files in pkg/: $([ -d "pkg" ] && echo "✅ YES" || echo "⚠️ NO")

## Next State Recommendation
**Proposed Next State:** $NEXT_STATE
**Reason:** $TRANSITION_REASON

## Issues Encountered
$(if [ $BLOCKED_COUNT -gt 0 ]; then
    echo "- $BLOCKED_COUNT engineers blocked during implementation"
fi)
$(if [ "$ALL_COMPLETE" = false ] && [ $BLOCKED_COUNT -eq 0 ]; then
    echo "- Timeout after ${MAX_MONITOR_TIME}s"
fi)
$(if [ ${#IMPLEMENTATIONS_WITH_ISSUES[@]} -gt 0 ]; then
    echo "- ${#IMPLEMENTATIONS_WITH_ISSUES[@]} implementations failed verification"
fi)

## R233 Active Monitoring Compliance
- Monitor interval: ${MONITOR_INTERVAL}s
- Total checks performed: $((ELAPSED_TIME / MONITOR_INTERVAL))
- Active monitoring: ✅ COMPLIANT

## THIS STATE IS FROM SF 3.0 ARCHITECTURE DOC
**Reference:** Architecture Part 3.5, Line 377
**Proof:** SF 3.0 was designed for full implementation workflow, not just demos
EOF

echo "✅ Monitoring report created: MONITORING-IMPLEMENTATION-REPORT.md"
```

## ⚠️ CRITICAL REQUIREMENTS

### Active Monitoring (R233)
- Check progress every 30 seconds
- Cannot passively wait
- Document each monitoring cycle
- Check every 5 messages for stalls

### No Direct Intervention
- If engineer is blocked, document it
- Do NOT try to write implementation code yourself (R006)
- Transition to ERROR_RECOVERY if needed

### Complete Verification
- All engineers must complete
- All work logs must exist (R343)
- All code must be committed and pushed
- Verification must pass for success

## 🚫 FORBIDDEN ACTIONS

1. **Writing implementation code yourself** - R006 violation (orchestrator never writes code)
2. **Passive waiting without checks** - R233 violation
3. **Continuing with incomplete implementations** - Will cause review failures
4. **Modifying SW Engineer work** - They own their implementation branches
5. **Skipping verification** - Must ensure complete implementations

## ✅ PROJECT_DONE CRITERIA

For successful transition to next state:
- [ ] All SW Engineers report IMPLEMENTATION_COMPLETE
- [ ] No engineers in BLOCKED state
- [ ] All work logs created (R343 compliant)
- [ ] All code committed and pushed
- [ ] Verification checks pass
- [ ] Completion report created
- [ ] Next state determined from results

## 🔄 STATE TRANSITIONS

### Success Path:
```
MONITORING_SWE_PROGRESS → SPAWN_CODE_REVIEWERS_FOR_EFFORT_REVIEW
```
- All implementations complete
- Verification passed
- Ready for code review

### Error Paths:
```
MONITORING_SWE_PROGRESS → ERROR_RECOVERY
```
- Engineers blocked
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
4. **Verify thoroughly** - Ensure complete implementations with proper artifacts

Remember: You're the PROJECT MANAGER - monitor, verify, report!

## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-software-factory
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**

## 🔴🔴🔴 STATE COMPLETION CHECKLIST 🔴🔴🔴

**Execute these steps IN ORDER to properly complete MONITORING_SWE_PROGRESS:**

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
  --current-state "MONITORING_SWE_PROGRESS" \
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
save_todos "MONITORING_SWE_PROGRESS_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - monitoring SWE progress complete [R287]"; then
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
  --current-state "MONITORING_SWE_PROGRESS" \
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
- ✅ Implementations complete, code review needed

**FALSE (0.1%):**
- ❌ CATASTROPHIC unrecoverable error
- ❌ Data corruption spreading
- ❌ Critical security violation
- ❌ NOT for R322 checkpoints
- ❌ NOT for user review needs

### 🚨 MONITORING_SWE_PROGRESS STATE PATTERN - NORMAL TRANSITIONS 🚨

**Monitoring states transition to next actions automatically:**
```bash
# After monitoring completes
echo "✅ Monitoring complete, all engineers finished implementation"

# Determine next action from results
if all_implementations_complete; then
    transition_to "SPAWN_CODE_REVIEWERS_FOR_EFFORT_REVIEW"
elif blocked_or_failed; then
    transition_to "ERROR_RECOVERY"
fi

# R322 checkpoint (if required by this transition)
echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"  # NORMAL operation!
exit 0  # If R322 checkpoint
```

**Why TRUE is correct:**
- Monitoring results drive automatic actions
- System knows what to do based on results
- **Implementations complete = Code review needed = NORMAL!**

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**
