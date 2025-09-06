# Orchestrator - MONITORING_BACKPORT_PROGRESS State Rules

## 🛑🛑🛑 R322 MANDATORY STOP BEFORE STATE TRANSITIONS 🛑🛑🛑

**SUPREME LAW - VIOLATION = -100% IMMEDIATE FAILURE**

### YOU MUST STOP AFTER:
1. ✅ Completing all TODOs for this state
2. ✅ Updating orchestrator-state.yaml with new state
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
Ready to transition to NEXT_STATE. Please use /continue-orchestrating.
```

**STOP MEANS STOP - Exit and wait for /continue-orchestrating**

---

## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED MONITORING_BACKPORT_PROGRESS STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_MONITORING_BACKPORT_PROGRESS
echo "$(date +%s) - Rules read and acknowledged for MONITORING_BACKPORT_PROGRESS" > .state_rules_read_orchestrator_MONITORING_BACKPORT_PROGRESS
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

## 📋 PRIMARY DIRECTIVES FOR MONITORING_BACKPORT_PROGRESS STATE

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
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-requirements.md`
   - Criticality: SUPREME - State updates required for all transitions
   - Summary: MUST update orchestrator-state.yaml before EVERY state transition
   - **CRITICAL**: Commit and push state changes immediately

4. **🚨🚨🚨 R304** - MANDATORY LINE COUNTER TOOL ENFORCEMENT (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R304-mandatory-line-counter-enforcement.md`
   - Criticality: BLOCKING - Manual counting = AUTOMATIC -100% FAILURE
   - Summary: MUST use tools/line-counter.sh for ALL line counting
   - **CRITICAL**: NEVER use wc -l or manual counting

### State-Specific Rules:

5. **⚠️⚠️⚠️ R237** - Waiting State Monitoring Requirements
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R237-waiting-state-monitoring.md`
   - Criticality: WARNING - Must actively monitor
   - Summary: Must check progress regularly, not passively wait

6. **🔴🔴🔴 R321** - Immediate Backport During Integration Protocol
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R321-immediate-backport-during-integration.md`
   - Criticality: SUPREME LAW - Immediate backporting required
   - Summary: ALL fixes must go to source branches IMMEDIATELY

## 🎯 STATE OBJECTIVES - MONITORING SW ENGINEER BACKPORT PROGRESS

In the MONITORING_BACKPORT_PROGRESS state, the ORCHESTRATOR is responsible for:

1. **Active Progress Monitoring**
   - Check each SW Engineer's state file regularly
   - Track completion status for each effort
   - Identify any BLOCKED engineers
   - Monitor for timeout conditions

2. **Progress Aggregation**
   - Count completed backports
   - Track remaining work
   - Calculate completion percentage
   - Estimate time to completion

3. **Issue Detection and Response**
   - Detect blocked engineers
   - Identify failed backports
   - Document issues for resolution
   - Determine if ERROR_RECOVERY needed

4. **Completion Verification**
   - Verify all backports complete
   - Check all branches updated
   - Confirm builds pass
   - Prepare for next state transition

## 📝 REQUIRED ACTIONS

### Step 1: Initial Status Assessment
```bash
# Get initial status of all SW Engineers
cd /efforts

echo "📊 Initial SW Engineer Status Check..."
echo "=================================="

# List all efforts being backported
EFFORTS_IN_PROGRESS=("effort-1" "effort-2" "effort-3")  # From spawn records

for EFFORT in "${EFFORTS_IN_PROGRESS[@]}"; do
    echo ""
    echo "Checking ${EFFORT}..."
    
    STATE_FILE="/efforts/${EFFORT}/sw-engineer-state.yaml"
    
    if [ -f "$STATE_FILE" ]; then
        CURRENT_STATE=$(grep "current_state:" "$STATE_FILE" | awk '{print $2}')
        echo "  State: $CURRENT_STATE"
        
        # Check for completion
        if [ "$CURRENT_STATE" = "BACKPORT_COMPLETE" ]; then
            echo "  Status: ✅ COMPLETE"
        elif [ "$CURRENT_STATE" = "BLOCKED" ]; then
            echo "  Status: ❌ BLOCKED"
            # Get block reason if available
            grep "block_reason:" "$STATE_FILE" 2>/dev/null
        else
            echo "  Status: ⏳ IN PROGRESS"
        fi
    else
        echo "  Status: ⏳ NOT STARTED (no state file yet)"
    fi
    
    # Check if branch was updated
    if [ -d "${EFFORT}/.git" ]; then
        cd "$EFFORT"
        LAST_COMMIT=$(git log -1 --format="%s" 2>/dev/null)
        echo "  Last commit: $LAST_COMMIT"
        cd ..
    fi
done
```

### Step 2: Continuous Monitoring Loop
```bash
# Monitor until all complete or timeout
MONITOR_INTERVAL=30  # Check every 30 seconds
MAX_MONITOR_TIME=1200  # Maximum 20 minutes
ELAPSED_TIME=0
ALL_COMPLETE=false

while [ $ELAPSED_TIME -lt $MAX_MONITOR_TIME ] && [ "$ALL_COMPLETE" = false ]; do
    echo ""
    echo "⏰ Monitor check at $(date '+%H:%M:%S') - Elapsed: ${ELAPSED_TIME}s"
    echo "------------------------------------------------"
    
    COMPLETED_COUNT=0
    BLOCKED_COUNT=0
    IN_PROGRESS_COUNT=0
    TOTAL_COUNT=${#EFFORTS_IN_PROGRESS[@]}
    
    # Check each effort
    for EFFORT in "${EFFORTS_IN_PROGRESS[@]}"; do
        STATE_FILE="/efforts/${EFFORT}/sw-engineer-state.yaml"
        
        if [ -f "$STATE_FILE" ]; then
            CURRENT_STATE=$(grep "current_state:" "$STATE_FILE" | awk '{print $2}')
            
            case "$CURRENT_STATE" in
                "BACKPORT_COMPLETE")
                    ((COMPLETED_COUNT++))
                    echo "✅ ${EFFORT}: Complete"
                    ;;
                "BLOCKED")
                    ((BLOCKED_COUNT++))
                    echo "❌ ${EFFORT}: BLOCKED"
                    BLOCK_REASON=$(grep "block_reason:" "$STATE_FILE" | cut -d: -f2-)
                    echo "   Reason: $BLOCK_REASON"
                    ;;
                *)
                    ((IN_PROGRESS_COUNT++))
                    echo "⏳ ${EFFORT}: In Progress ($CURRENT_STATE)"
                    ;;
            esac
        else
            ((IN_PROGRESS_COUNT++))
            echo "⏳ ${EFFORT}: Initializing..."
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
        echo "🎉 All backports complete!"
        ALL_COMPLETE=true
        break
    fi
    
    # Check if any blocked
    if [ $BLOCKED_COUNT -gt 0 ]; then
        echo "⚠️ WARNING: ${BLOCKED_COUNT} engineers blocked - may need intervention"
    fi
    
    # Wait before next check
    sleep $MONITOR_INTERVAL
    ELAPSED_TIME=$((ELAPSED_TIME + MONITOR_INTERVAL))
done

# Check final status
if [ "$ALL_COMPLETE" = false ]; then
    if [ $BLOCKED_COUNT -gt 0 ]; then
        echo "❌ CRITICAL: Backports blocked - ERROR_RECOVERY needed"
    else
        echo "⚠️ TIMEOUT: Not all backports completed in time"
    fi
fi
```

### Step 3: Verify Backport Quality
```bash
# For completed backports, verify quality
echo ""
echo "🔍 Verifying Backport Quality..."
echo "================================"

VERIFICATION_PASSED=true

for EFFORT in "${EFFORTS_IN_PROGRESS[@]}"; do
    STATE_FILE="/efforts/${EFFORT}/sw-engineer-state.yaml"
    
    if [ -f "$STATE_FILE" ]; then
        CURRENT_STATE=$(grep "current_state:" "$STATE_FILE" | awk '{print $2}')
        
        if [ "$CURRENT_STATE" = "BACKPORT_COMPLETE" ]; then
            echo ""
            echo "Verifying ${EFFORT}..."
            
            cd "/efforts/${EFFORT}"
            
            # Check if commits were made
            RECENT_COMMITS=$(git log --oneline -5 --grep="backport\|fix\|integration" 2>/dev/null | wc -l)
            if [ $RECENT_COMMITS -gt 0 ]; then
                echo "  ✅ Found ${RECENT_COMMITS} backport commits"
            else
                echo "  ⚠️ No backport commits found"
                VERIFICATION_PASSED=false
            fi
            
            # Check if branch is ahead of origin
            AHEAD_COUNT=$(git rev-list --count origin/HEAD..HEAD 2>/dev/null || echo "0")
            if [ "$AHEAD_COUNT" = "0" ]; then
                echo "  ⚠️ Branch not pushed or no new commits"
            else
                echo "  ✅ Branch has ${AHEAD_COUNT} new commits"
            fi
            
            # Check build status if recorded
            if grep -q "build_status: PASS" "$STATE_FILE"; then
                echo "  ✅ Build passed"
            else
                echo "  ⚠️ Build status unknown or failed"
            fi
            
            # Check test status if recorded
            if grep -q "test_status: PASS" "$STATE_FILE"; then
                echo "  ✅ Tests passed"
            else
                echo "  ⚠️ Test status unknown or failed"
            fi
            
            cd - > /dev/null
        fi
    fi
done

if [ "$VERIFICATION_PASSED" = true ]; then
    echo ""
    echo "✅ All backport verifications PASSED"
else
    echo ""
    echo "⚠️ Some verifications failed - review needed"
fi
```

### Step 4: Create Completion Report
```bash
# Create comprehensive completion report
cd /efforts/integration-testing

cat > BACKPORT-COMPLETION-REPORT.md << 'EOF'
# Backport Progress Monitoring - Final Report

## Monitoring Summary
- Start Time: [timestamp]
- End Time: [timestamp]
- Total Duration: ${ELAPSED_TIME} seconds
- Monitor Checks: $((ELAPSED_TIME / MONITOR_INTERVAL))

## Backport Status
### Completed Successfully: ${COMPLETED_COUNT}
$(for EFFORT in completed_efforts; do
    echo "- ${EFFORT}: ✅ Complete"
done)

### In Progress: ${IN_PROGRESS_COUNT}
$(for EFFORT in progress_efforts; do
    echo "- ${EFFORT}: ⏳ Still working"
done)

### Blocked: ${BLOCKED_COUNT}
$(for EFFORT in blocked_efforts; do
    echo "- ${EFFORT}: ❌ Blocked"
done)

## Verification Results
- All builds passing: YES/NO
- All tests passing: YES/NO
- All branches pushed: YES/NO
- Ready for PR creation: YES/NO

## Issues Encountered
[List any issues, blocks, or failures]

## Next State Recommendation
$(if [ "$ALL_COMPLETE" = true ] && [ "$VERIFICATION_PASSED" = true ]; then
    echo "✅ Proceed to PR_PLAN_CREATION"
else
    echo "⚠️ ERROR_RECOVERY may be needed"
fi)

## Branch Update Summary
[List commits added to each branch]
EOF

echo "✅ Completion report created"
```

### Step 5: Prepare State Transition
```bash
# Update state based on results
cd $CLAUDE_PROJECT_DIR

if [ "$ALL_COMPLETE" = true ] && [ "$VERIFICATION_PASSED" = true ]; then
    # Success - all backports complete
    cat > orchestrator-state.yaml << 'EOF'
current_state: PR_PLAN_CREATION
previous_state: MONITORING_BACKPORT_PROGRESS
backport_status: COMPLETE
backports_completed: ${COMPLETED_COUNT}/${TOTAL_COUNT}
verification: PASSED
all_branches_updated: true
ready_for_pr: true
monitoring_duration: ${ELAPSED_TIME}
next_action: Create PR plan for updated branches
EOF
    
    git add orchestrator-state.yaml
    git commit -m "state: all backports complete - transition to PR_PLAN_CREATION"
    
elif [ $BLOCKED_COUNT -gt 0 ]; then
    # Blocked - need error recovery
    cat > orchestrator-state.yaml << 'EOF'
current_state: ERROR_RECOVERY
previous_state: MONITORING_BACKPORT_PROGRESS
backport_status: BLOCKED
backports_completed: ${COMPLETED_COUNT}/${TOTAL_COUNT}
backports_blocked: ${BLOCKED_COUNT}/${TOTAL_COUNT}
issues: "SW Engineers blocked during backport implementation"
requires_intervention: true
EOF
    
    git add orchestrator-state.yaml
    git commit -m "state: backports blocked - transition to ERROR_RECOVERY"
    
else
    # Timeout or other issue
    cat > orchestrator-state.yaml << 'EOF'
current_state: ERROR_RECOVERY
previous_state: MONITORING_BACKPORT_PROGRESS
backport_status: TIMEOUT
backports_completed: ${COMPLETED_COUNT}/${TOTAL_COUNT}
backports_in_progress: ${IN_PROGRESS_COUNT}/${TOTAL_COUNT}
timeout_after: ${MAX_MONITOR_TIME}
EOF
    
    git add orchestrator-state.yaml
    git commit -m "state: backport timeout - transition to ERROR_RECOVERY"
fi

git push
echo "✅ State updated and pushed"
```

## ⚠️ CRITICAL REQUIREMENTS

### Active Monitoring (R237)
- Check progress every 30-60 seconds
- Cannot passively wait
- Document each monitoring cycle

### No Direct Intervention
- If engineer is blocked, document it
- Do NOT try to fix issues yourself
- Transition to ERROR_RECOVERY if needed

### Complete Verification
- All engineers must complete
- All branches must be updated
- Verification must pass for success

## 🚫 FORBIDDEN ACTIONS

1. **Fixing blocked backports yourself** - R006 violation
2. **Passive waiting without checks** - R237 violation
3. **Continuing with incomplete backports** - Will cause PR failures
4. **Modifying SW Engineer work** - They own their branches
5. **Skipping verification** - Must ensure quality

## ✅ SUCCESS CRITERIA

For successful transition to PR_PLAN_CREATION:
- [ ] All SW Engineers report BACKPORT_COMPLETE
- [ ] No engineers in BLOCKED state
- [ ] All branches have new commits
- [ ] Verification checks pass
- [ ] Completion report created
- [ ] State updated appropriately

## 🔄 STATE TRANSITIONS

### Success Path:
```
MONITORING_BACKPORT_PROGRESS → PR_PLAN_CREATION
```
- All backports complete
- Verification passed
- Ready for PRs

### Error Paths:
```
MONITORING_BACKPORT_PROGRESS → ERROR_RECOVERY
```
- Engineers blocked
- Timeout occurred
- Verification failed

## 🎓 GRADING CRITERIA

You will be evaluated on:
1. **Active Monitoring** (35%)
   - Regular progress checks
   - R237 compliance
   
2. **Issue Detection** (25%)
   - Identifying blocks quickly
   - Catching failures
   
3. **Verification Quality** (25%)
   - Thorough verification
   - Quality assurance
   
4. **Documentation** (15%)
   - Clear progress tracking
   - Complete final report

## 💡 TIPS FOR SUCCESS

1. **Check frequently** - Every 30 seconds is good
2. **Document everything** - Show active monitoring
3. **Detect issues early** - Don't wait for timeout
4. **Verify thoroughly** - Ensure quality backports

Remember: You're the PROJECT MANAGER - monitor, verify, report!

## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-orchestrating
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**