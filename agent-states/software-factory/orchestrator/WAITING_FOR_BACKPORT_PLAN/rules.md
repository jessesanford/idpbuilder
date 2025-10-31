# Orchestrator - WAITING_FOR_BACKPORT_PLAN State Rules


## ✅ BACKPORTS ARE NORMAL OPERATIONS - CONTINUE AUTOMATICALLY

**IMPORTANT CLARIFICATION PER R322 AND R405:**

Backports are NORMAL software development operations that should proceed automatically:
- ✅ NO user review required for backport plans
- ✅ Continue automatically to SPAWN_SW_ENGINEER_BACKPORT_FIXES
- ✅ Set CONTINUE-SOFTWARE-FACTORY=TRUE when plan is ready
- ✅ This is routine fix propagation, not an exceptional situation

### When transitioning from WAITING_FOR_BACKPORT_PLAN → SPAWN_SW_ENGINEER_BACKPORT_FIXES:
```markdown
## ✅ Backport Plan Ready - Proceeding Automatically

### Plan Details:
- Location: BACKPORT-PLAN.md
- Branches to update: [List branches]
- Fixes to apply: [Summary]

### Mandatory Stop Protocol:
- Current State: WAITING_FOR_BACKPORT_PLAN ✅
- Next State: SPAWN_SW_ENGINEER_BACKPORT_FIXES
- Action: Stop with TRUE flag after updating state
- STOP PROCESSING (mandatory for context preservation)

CONTINUE-SOFTWARE-FACTORY=TRUE
```

**Backports are NORMAL - Stop with TRUE flag for automatic restart!**

See: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-before-state-transitions.md`
See: `$CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md`

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
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
   - Criticality: SUPREME - State updates required for all transitions
   - Summary: MUST update orchestrator-state-v3.json before EVERY state transition
   - **CRITICAL**: Commit and push state changes immediately

   - **CRITICAL**: NEVER use wc -l or manual counting

### State-Specific Rules:

5. **🔴🔴🔴 R233** - All States Require Immediate Action (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R233-all-states-immediate-action.md`
   - Criticality: SUPREME - Violation = AUTOMATIC FAILURE
   - Summary: WAITING states must ACTIVELY poll/monitor - not passively wait
   - **CRITICAL**: "Waiting" means ACTIVELY CHECKING (poll every few seconds, report what you're checking, set timeouts)

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
    
    cat > orchestrator-state-v3.json << 'EOF'
current_state: SPAWN_SW_ENGINEER_BACKPORT_FIXES
previous_state: WAITING_FOR_BACKPORT_PLAN
backport_status: PLAN_READY
backport_plan: /efforts/integration-testing/BACKPORT-PLAN.md
plan_validation: PASSED
next_action: Spawn SW Engineers to implement backport fixes
code_reviewer_status: COMPLETE
wait_duration: ${ELAPSED_TIME}
EOF
    
    git add orchestrator-state-v3.json
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

## ✅ PROJECT_DONE CRITERIA

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

## 💡 TIPS FOR PROJECT_DONE

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

### 🔴🔴🔴 MANDATORY VALIDATION REQUIREMENT 🔴🔴🔴

**Per R288 and R324**: ALL state file updates MUST be validated before commit:

```bash
# After ANY update to orchestrator-state-v3.json:
"$CLAUDE_PROJECT_DIR/tools/validate-state.sh" orchestrator-state-v3.json || {
    echo "❌ State file validation failed!"
    exit 288
}
```

**Use helper functions for automatic validation:**
```bash
# Source the helper functions
source "$CLAUDE_PROJECT_DIR/utilities/state-file-update-functions.sh"

# Use safe functions that include validation:
safe_state_transition "NEW_STATE" "reason"
safe_update_field "field_name" "value"
```



## 🔴🔴🔴 STATE COMPLETION CHECKLIST 🔴🔴🔴

**Execute these steps IN ORDER to properly complete WAITING_FOR_BACKPORT_PLAN:**

### ✅ Step 1: Complete State-Specific Work
**Refer to "Primary Actions" and "Mandatory Validations" sections above for state-specific tasks.**

Once all state work is complete, proceed to mandatory exit protocol:

---

### ✅ Step 2: Set Proposed Next State and Transition Reason
```bash
# Based on state work results, set variables for State Manager
PROPOSED_NEXT_STATE="[DETERMINE_FROM_STATE_LOGIC]"
TRANSITION_REASON="[REASON_FOR_TRANSITION]"
echo "Proposed next state: $PROPOSED_NEXT_STATE"
echo "Transition reason: $TRANSITION_REASON"
```

---

### ✅ Step 3: Spawn State Manager for State Transition (R288 - SUPREME LAW)
```bash
# State Manager handles ALL state file updates atomically
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "CURRENT_STATE_NAME" \
  --proposed-next-state "$PROPOSED_NEXT_STATE" \
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
save_todos "STATE_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - state complete [R287]"; then
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
echo "🛑 Stopping for context preservation - use /continue-orchestrating to resume"
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
- User runs /continue-orchestrating to resume
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
PROPOSED_NEXT_STATE="SPAWN_SW_ENGINEER_BACKPORT_FIXES"
TRANSITION_REASON="State work complete"

# 3. Spawn State Manager for state transition
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "CURRENT_STATE" \
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
- ✅ Ready for /continue-orchestrating
- ✅ Waiting for user to continue (NORMAL)
- ✅ Plan ready for review (agent done, factory proceeds)

**FALSE (0.1%):**
- ❌ CATASTROPHIC unrecoverable error
- ❌ Data corruption spreading
- ❌ Critical security violation
- ❌ NOT for R322 checkpoints
- ❌ NOT for user review needs
### 🚨 WAITING STATE PATTERN - CRITICAL UNDERSTANDING 🚨

**This is a WAITING state. Common source of incorrect FALSE usage!**

**WRONG interpretation:**
> "R322 mandates stop before transition"
> "State work is complete (validation done)"
> "User needs to invoke /continue-orchestrating"
> "Therefore I must set CONTINUE-SOFTWARE-FACTORY=FALSE"

**CORRECT interpretation:**
> "R322 checkpoint is NORMAL procedure for context preservation"
> "State work completed successfully = NORMAL outcome"
> "Waiting for /continue is DESIGNED user experience"
> "System KNOWS next step from state file"
> "NO manual intervention required, just normal continuation"
> "Therefore set CONTINUE-SOFTWARE-FACTORY=TRUE"

**The key distinction:**
- **Stopping inference** (`exit 0`) = Context management (ALWAYS at R322 points)
- **Continuation flag** = Can automation proceed? (TRUE unless catastrophic failure)

**ONLY use FALSE if:**
- ❌ The thing we're waiting for completely disappeared (agents crashed with no recovery)
- ❌ Results arrived but are completely corrupted/unreadable
- ❌ State file corruption prevents determining what to wait for
- ❌ System deadlock with no automated resolution
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

