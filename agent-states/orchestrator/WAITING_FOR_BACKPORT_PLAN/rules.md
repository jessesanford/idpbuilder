# Orchestrator - WAITING_FOR_BACKPORT_PLAN State Rules

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

**YOU HAVE ENTERED WAITING_FOR_BACKPORT_PLAN STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_WAITING_FOR_BACKPORT_PLAN
echo "$(date +%s) - Rules read and acknowledged for WAITING_FOR_BACKPORT_PLAN" > .state_rules_read_orchestrator_WAITING_FOR_BACKPORT_PLAN
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

## 📋 PRIMARY DIRECTIVES FOR WAITING_FOR_BACKPORT_PLAN STATE

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
   - Summary: Waiting states must check progress regularly

## 🎯 STATE OBJECTIVES - WAITING FOR CODE REVIEWER'S BACKPORT PLAN

In the WAITING_FOR_BACKPORT_PLAN state, the ORCHESTRATOR is responsible for:

1. **Monitoring Code Reviewer Progress**
   - Check for code-reviewer-state.yaml updates
   - Look for BACKPORT-PLAN.md creation
   - Monitor for any BLOCKED status
   - Track time elapsed

2. **Verifying Plan Completeness**
   - Once plan appears, verify it exists
   - Check plan covers all needed fixes
   - Ensure plan has clear instructions
   - Validate plan structure

3. **Preparing for Next State**
   - When plan is complete, prepare transition
   - Update state to SPAWN_SW_ENGINEER_BACKPORT_FIXES
   - Document plan location
   - Stop and wait for continuation

## 📝 REQUIRED ACTIONS

### Step 1: Initial Status Check
```bash
# Check current status of Code Reviewer
cd /efforts/integration-testing

echo "📊 Checking Code Reviewer status..."

# Check for state file
if [ -f "code-reviewer-state.yaml" ]; then
    echo "Code Reviewer state file found:"
    cat code-reviewer-state.yaml
    
    # Extract current state
    REVIEWER_STATE=$(grep "current_state:" code-reviewer-state.yaml | awk '{print $2}')
    echo "Code Reviewer is in state: $REVIEWER_STATE"
else
    echo "⏳ Code Reviewer state file not yet created"
    echo "Code Reviewer may still be initializing"
fi

# Check for plan file
if [ -f "BACKPORT-PLAN.md" ]; then
    echo "✅ BACKPORT-PLAN.md exists!"
    echo "Plan size: $(wc -l BACKPORT-PLAN.md | awk '{print $1}') lines"
else
    echo "⏳ BACKPORT-PLAN.md not yet created"
fi
```

### Step 2: Monitor Progress Loop
```bash
# Monitor Code Reviewer progress
MONITORING_INTERVAL=30  # Check every 30 seconds
MAX_WAIT_TIME=600      # Maximum 10 minutes
ELAPSED_TIME=0

while [ $ELAPSED_TIME -lt $MAX_WAIT_TIME ]; do
    echo "⏰ Monitoring check at $(date '+%H:%M:%S')"
    
    # Check if plan is complete
    if [ -f "BACKPORT-PLAN.md" ]; then
        # Check if Code Reviewer marked complete
        if [ -f "code-reviewer-state.yaml" ]; then
            STATE=$(grep "current_state:" code-reviewer-state.yaml | awk '{print $2}')
            
            if [ "$STATE" = "BACKPORT_PLAN_COMPLETE" ]; then
                echo "✅ Code Reviewer has completed the backport plan!"
                break
            else
                echo "📝 Plan exists but Code Reviewer still working (state: $STATE)"
            fi
        fi
    else
        echo "⏳ Still waiting for BACKPORT-PLAN.md..."
    fi
    
    # Check for blocked status
    if grep -q "BLOCKED" code-reviewer-state.yaml 2>/dev/null; then
        echo "⚠️ WARNING: Code Reviewer reports BLOCKED status!"
        echo "Check for issues that need resolution"
    fi
    
    # Wait before next check
    echo "Next check in ${MONITORING_INTERVAL} seconds..."
    sleep $MONITORING_INTERVAL
    
    ELAPSED_TIME=$((ELAPSED_TIME + MONITORING_INTERVAL))
    echo "Total wait time: ${ELAPSED_TIME} seconds"
done

# Check if we timed out
if [ $ELAPSED_TIME -ge $MAX_WAIT_TIME ]; then
    echo "⚠️ TIMEOUT: Code Reviewer has not completed plan in ${MAX_WAIT_TIME} seconds"
    echo "Manual intervention may be required"
fi
```

### Step 3: Validate Plan Completeness
```bash
# Once plan exists, validate it
if [ -f "BACKPORT-PLAN.md" ]; then
    echo "🔍 Validating backport plan completeness..."
    
    # Check for required sections
    VALIDATION_PASSED=true
    
    # Check for effort branches
    if ! grep -q "## Effort\|## Branch\|effort-" BACKPORT-PLAN.md; then
        echo "❌ Plan missing effort branch specifications"
        VALIDATION_PASSED=false
    else
        echo "✅ Plan includes effort branches"
    fi
    
    # Check for fix descriptions
    if ! grep -q "## Fix\|### Fix\|Files to Modify" BACKPORT-PLAN.md; then
        echo "❌ Plan missing fix descriptions"
        VALIDATION_PASSED=false
    else
        echo "✅ Plan includes fix descriptions"
    fi
    
    # Check for verification steps
    if ! grep -q "Verification\|Test\|Validate" BACKPORT-PLAN.md; then
        echo "⚠️ Plan may be missing verification steps"
    else
        echo "✅ Plan includes verification steps"
    fi
    
    # Count efforts covered
    EFFORT_COUNT=$(grep -c "^## Effort\|^### Branch:" BACKPORT-PLAN.md)
    echo "📊 Plan covers ${EFFORT_COUNT} effort branches"
    
    if [ "$VALIDATION_PASSED" = true ]; then
        echo "✅ Backport plan validation PASSED"
    else
        echo "❌ Backport plan validation FAILED - may need Code Reviewer revision"
    fi
else
    echo "❌ No BACKPORT-PLAN.md found to validate"
fi
```

### Step 4: Prepare Transition When Complete
```bash
# When plan is ready and validated
if [ -f "BACKPORT-PLAN.md" ] && [ "$VALIDATION_PASSED" = true ]; then
    echo "✅ Ready to transition to SPAWN_SW_ENGINEER_BACKPORT_FIXES"
    
    # Document completion
    cat > WAITING-COMPLETION-REPORT.md << 'EOF'
# Waiting for Backport Plan - Completion Report

## Status: COMPLETE
- Code Reviewer has completed BACKPORT-PLAN.md
- Plan validation passed
- Ready to spawn SW Engineers

## Plan Summary
- Location: /efforts/integration-testing/BACKPORT-PLAN.md
- Efforts covered: ${EFFORT_COUNT}
- Plan validated: YES

## Next State
SPAWN_SW_ENGINEER_BACKPORT_FIXES - Will spawn SW Engineers to implement fixes

## Timing
- Wait started: [timestamp]
- Plan completed: [timestamp]
- Total wait time: ${ELAPSED_TIME} seconds
EOF
    
    # Update orchestrator state
    cd $CLAUDE_PROJECT_DIR
    
    cat > orchestrator-state.yaml << 'EOF'
current_state: SPAWN_SW_ENGINEER_BACKPORT_FIXES
previous_state: WAITING_FOR_BACKPORT_PLAN
backport_status: PLAN_READY
backport_plan: /efforts/integration-testing/BACKPORT-PLAN.md
plan_validation: PASSED
next_action: Spawn SW Engineers to implement backport fixes
code_reviewer_status: COMPLETE
wait_duration: ${ELAPSED_TIME}
EOF
    
    git add orchestrator-state.yaml
    git commit -m "state: transition to SPAWN_SW_ENGINEER_BACKPORT_FIXES - plan ready"
    git push
    
    echo "✅ State updated for transition"
fi
```

## ⚠️ CRITICAL REQUIREMENTS

### Active Monitoring Required
- Must check progress regularly (R237)
- Cannot just wait passively
- Document monitoring checks

### No Direct Work
- ONLY monitor and validate
- Do NOT create the plan yourself
- Do NOT modify Code Reviewer's plan

### Proper Validation
- Ensure plan covers all efforts
- Verify plan has actionable instructions
- Check for completeness before proceeding

## 🚫 FORBIDDEN ACTIONS

1. **Creating the plan yourself** - Code Reviewer's job
2. **Modifying the plan** - Use as-is from Code Reviewer
3. **Spawning SW Engineers prematurely** - Wait for complete plan
4. **Passive waiting without monitoring** - R237 violation
5. **Continuing without plan validation** - Must verify first

## ✅ SUCCESS CRITERIA

Before transitioning to SPAWN_SW_ENGINEER_BACKPORT_FIXES:
- [ ] Code Reviewer completed BACKPORT-PLAN.md
- [ ] Plan validated for completeness
- [ ] All effort branches covered in plan
- [ ] State updated to next state
- [ ] Monitoring documented

## 🔄 STATE TRANSITIONS

### Success Path:
```
WAITING_FOR_BACKPORT_PLAN → SPAWN_SW_ENGINEER_BACKPORT_FIXES
```
- Plan complete and validated
- Ready to spawn engineers

### Timeout/Error Path:
```
WAITING_FOR_BACKPORT_PLAN → ERROR_RECOVERY
```
- Code Reviewer blocked
- Plan incomplete after timeout
- Validation failures

## 🎓 GRADING CRITERIA

You will be evaluated on:
1. **Active Monitoring** (40%)
   - Regular progress checks
   - R237 compliance
   
2. **Validation Quality** (30%)
   - Thorough plan validation
   - Completeness verification
   
3. **State Management** (20%)
   - Proper transitions
   - Clear documentation
   
4. **Patience** (10%)
   - Not rushing
   - Allowing Code Reviewer time

## 💡 TIPS FOR SUCCESS

1. **Monitor actively** - Check every 30-60 seconds
2. **Validate thoroughly** - Ensure plan is complete
3. **Document monitoring** - Show you're checking
4. **Be patient** - Code Reviewer needs time to analyze

Remember: Good things come to those who WAIT (actively)!

## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-orchestrating
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**