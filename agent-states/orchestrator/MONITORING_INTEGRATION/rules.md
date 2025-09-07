# Orchestrator - MONITORING_INTEGRATION State Rules

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

# ORCHESTRATOR STATE: MONITORING_INTEGRATION

## 🔴🔴🔴 SUPREME DIRECTIVE: INTEGRATION FEEDBACK ENFORCEMENT 🔴🔴🔴

**YOU MUST CHECK FOR INTEGRATION REPORTS AND ACT ON FAILURES!**

## State Overview

In MONITORING_INTEGRATION, you are monitoring the Integration Agent's progress and MUST check for integration reports to determine next state.

## Required Actions

### 1. Monitor Integration Agent Progress

#### 🔴🔴🔴 CRITICAL: Monitor Split-Aware Merge Ordering 🔴🔴🔴

**When monitoring integration with splits, verify correct merge order!**

```bash
# Check if Integration Agent is still running
INTEGRATION_PID=$(pgrep -f "integration-agent" || echo "")
if [ -n "$INTEGRATION_PID" ]; then
    echo "Integration Agent still running (PID: $INTEGRATION_PID)"
    
    # Monitor for split ordering violations
    WORK_LOG="efforts/phase${PHASE}/wave${WAVE}/integration-workspace/work-log.md"
    if [ -f "$WORK_LOG" ]; then
        # Check for out-of-order split merges
        check_split_merge_order() {
            local log_file="$1"
            
            # Extract merged branches in order
            MERGED_BRANCHES=$(grep "MERGED:" "$log_file" | awk '{print $2}')
            
            for branch in $MERGED_BRANCHES; do
                # If this is a dependent effort, check its dependencies
                EFFORT=$(echo "$branch" | sed 's/-split-[0-9]*//')
                DEPENDENCIES=$(yq ".efforts.\"$EFFORT\".dependencies[]" orchestrator-state.yaml 2>/dev/null)
                
                for dep in $DEPENDENCIES; do
                    # Check if dependency has splits
                    SPLIT_COUNT=$(yq ".split_tracking.\"$dep\".split_count // 0" orchestrator-state.yaml)
                    if [ "$SPLIT_COUNT" -gt 0 ]; then
                        # Verify ALL splits of dependency are merged
                        for i in $(seq 1 $SPLIT_COUNT); do
                            SPLIT_PATTERN="${dep}-split-$(printf "%03d" $i)"
                            if ! grep -q "MERGED:.*$SPLIT_PATTERN" "$log_file"; then
                                echo "🔴 CRITICAL: Merge order violation detected!"
                                echo "   $branch merged but dependency $SPLIT_PATTERN not merged!"
                                echo "   This violates R302 split tracking protocol!"
                                
                                # Signal Integration Agent to STOP
                                echo "STOP: MERGE_ORDER_VIOLATION" > efforts/phase${PHASE}/wave${WAVE}/integration-workspace/STOP_SIGNAL
                                return 1
                            fi
                        done
                    fi
                done
            done
            
            echo "✅ Merge order correct so far"
        }
        
        check_split_merge_order "$WORK_LOG"
    fi
    
    # Stay in MONITORING_INTEGRATION
    sleep 5
    continue
fi
```

### 2. 🔴🔴🔴 CHECK FOR INTEGRATION REPORT AND ENFORCE BUILD/TEST GATES 🔴🔴🔴

**Per R291 SUPREME GATE: Build/Test/Demo MUST ALL PASS or transition to ERROR_RECOVERY!**

```bash
# MANDATORY: Check for integration report AND enforce gates
PHASE=$(yq '.current_phase' orchestrator-state.yaml)
WAVE=$(yq '.current_wave' orchestrator-state.yaml)
REPORT_FILE="efforts/phase${PHASE}/wave${WAVE}/integration-workspace/INTEGRATION_REPORT.md"
DEMO_STATUS_FILE="efforts/phase${PHASE}/wave${WAVE}/integration-workspace/DEMO_STATUS.md"

if [ ! -f "$REPORT_FILE" ]; then
    echo "🔴 CRITICAL: No integration report found at $REPORT_FILE"
    # NO REPORT = IMMEDIATE ERROR_RECOVERY
    UPDATE_STATE="ERROR_RECOVERY"
    UPDATE_REASON="No integration report generated - R291 gate cannot be verified"
else
    echo "✅ Found integration report, enforcing R291 gates..."
    
    # Extract status from report
    INTEGRATION_STATUS=$(grep "^Integration Status:" "$REPORT_FILE" | cut -d: -f2 | tr -d ' ')
    BUILD_STATUS=$(grep "^Build Status:" "$REPORT_FILE" | cut -d: -f2 | tr -d ' ')
    TEST_STATUS=$(grep "^Test Status:" "$REPORT_FILE" | cut -d: -f2 | tr -d ' ')
    DEMO_STATUS=$(grep "^Demo Status:" "$REPORT_FILE" | cut -d: -f2 | tr -d ' ' || echo "NOT_RUN")
    
    echo "🔍 Gate Status Check:"
    echo "  Build Status: $BUILD_STATUS"
    echo "  Test Status: $TEST_STATUS"
    echo "  Demo Status: $DEMO_STATUS"
    echo "  Integration Status: $INTEGRATION_STATUS"
    
    # 🔴🔴🔴 R291 SUPREME GATE ENFORCEMENT 🔴🔴🔴
    # ANY failure = MANDATORY ERROR_RECOVERY
    
    # BUILD GATE CHECK
    if [[ "$BUILD_STATUS" != "PASSING" ]] && [[ "$BUILD_STATUS" != "SUCCESS" ]]; then
        echo "🔴🔴🔴 BUILD GATE FAILED - MANDATORY ERROR_RECOVERY 🔴🔴🔴"
        echo "R291 VIOLATION: Build did not pass ($BUILD_STATUS)"
        UPDATE_STATE="ERROR_RECOVERY"
        UPDATE_REASON="R291 BUILD GATE FAILURE: $BUILD_STATUS - cannot proceed without successful build"
        
    # TEST GATE CHECK
    elif [[ "$TEST_STATUS" != "PASSING" ]] && [[ "$TEST_STATUS" != "SUCCESS" ]]; then
        echo "🔴🔴🔴 TEST GATE FAILED - MANDATORY ERROR_RECOVERY 🔴🔴🔴"
        echo "R291 VIOLATION: Tests did not pass ($TEST_STATUS)"
        UPDATE_STATE="ERROR_RECOVERY"
        UPDATE_REASON="R291 TEST GATE FAILURE: $TEST_STATUS - cannot proceed without passing tests"
        
    # DEMO GATE CHECK
    elif [[ "$DEMO_STATUS" != "PASSING" ]] && [[ "$DEMO_STATUS" != "SUCCESS" ]]; then
        echo "🔴🔴🔴 DEMO GATE FAILED - MANDATORY ERROR_RECOVERY 🔴🔴🔴"
        echo "R291 VIOLATION: Demo did not pass ($DEMO_STATUS)"
        UPDATE_STATE="ERROR_RECOVERY"
        UPDATE_REASON="R291 DEMO GATE FAILURE: $DEMO_STATUS - cannot proceed without working demo"
        
    # INTEGRATION STATUS CHECK
    elif [[ "$INTEGRATION_STATUS" != "SUCCESS" ]]; then
        echo "🔴 Integration failed - checking if it's fixable"
        # For integration failures that aren't build/test/demo, go to feedback review
        UPDATE_STATE="INTEGRATION_FEEDBACK_REVIEW"
        UPDATE_REASON="Integration issues detected - review needed (R300)"
        
    # ALL GATES PASSED
    else
        echo "✅✅✅ ALL R291 GATES PASSED ✅✅✅"
        echo "  ✅ Build: $BUILD_STATUS"
        echo "  ✅ Tests: $TEST_STATUS"
        echo "  ✅ Demo: $DEMO_STATUS"
        echo "  ✅ Integration: $INTEGRATION_STATUS"
        UPDATE_STATE="WAVE_REVIEW"
        UPDATE_REASON="All gates passed - integration successful"
    fi
fi

echo ""
echo "🎯 DECISION: Transitioning to $UPDATE_STATE"
echo "📝 REASON: $UPDATE_REASON"
```

### 3. Update State File
```bash
# Update orchestrator state
yq eval ".current_state = \"$UPDATE_STATE\"" -i orchestrator-state.yaml
yq eval ".state_transition_history += [{\"from\": \"MONITORING_INTEGRATION\", \"to\": \"$UPDATE_STATE\", \"timestamp\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\", \"reason\": \"$UPDATE_REASON\"}]" -i orchestrator-state.yaml

# Commit state change
git add orchestrator-state.yaml
git commit -m "state: MONITORING_INTEGRATION → $UPDATE_STATE - $UPDATE_REASON"
git push

# 🛑 STOP per R322 - State has been updated, now stop!
echo "🛑 Stopping before $UPDATE_STATE state (per R322)"
echo "State file updated to: $UPDATE_STATE"
echo "When restarted with /continue-orchestrating, will continue from $UPDATE_STATE"
# EXIT HERE - DO NOT CONTINUE
```

## Valid Transitions

Based on integration report analysis:

1. **SUCCESS Path**: `MONITORING_INTEGRATION` → `WAVE_REVIEW`
   - When: Integration, build, and tests all pass
   
2. **FAILURE Path**: `MONITORING_INTEGRATION` → `INTEGRATION_FEEDBACK_REVIEW`
   - When: Integration failed, build blocked, or tests fail
   
3. **ERROR Path**: `MONITORING_INTEGRATION` → `ERROR_RECOVERY`
   - When: No report found or unexpected status

## Grading Criteria

- ✅ **+20%**: Properly check for integration report
- ✅ **+20%**: Correctly parse report status fields
- ✅ **+20%**: Transition to INTEGRATION_FEEDBACK_REVIEW on failures
- ✅ **+20%**: Never ignore integration failures
- ✅ **+20%**: Update state file with proper reason

## Common Violations

- ❌ **-100%**: Ignoring integration failures and marking COMPLETE
- ❌ **-50%**: Not checking for integration report
- ❌ **-50%**: Transitioning to WAVE_REVIEW when integration failed
- ❌ **-30%**: Not parsing report status fields

## Related Rules

- R291: Integration Demo Requirement (CRITICAL - demo must pass)
- R300: Comprehensive Fix Management Protocol
- R238: Integration Report Evaluation Protocol (to be created)
- R260: Integration Agent Core Requirements
- R263: Integration Documentation Requirements
- R206: State Machine Transition Validation

## 🔴🔴🔴 MANDATORY MEASUREMENT RULE - R304 🔴🔴🔴

**ABSOLUTE REQUIREMENTS:**
- ✅ MUST use `$CLAUDE_PROJECT_DIR/tools/line-counter.sh` for ALL line counting
- ❌ NEVER use `wc -l` or any manual counting method
- ❌ NEVER count lines any other way - this is a -100% automatic failure
- ✅ MUST specify both -b (base branch) and -c (current branch) parameters
- ✅ Base branch MUST be phase integration branch (NOT "main")

**Failure to use the line counter tool = AUTOMATIC -100% GRADE**

## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-orchestrating
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**
