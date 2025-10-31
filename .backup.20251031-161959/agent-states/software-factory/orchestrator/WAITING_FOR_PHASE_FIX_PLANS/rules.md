# WAITING_FOR_PHASE_FIX_PLANS State Rules


## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

## 📋 PRIMARY DIRECTIVES FOR WAITING_FOR_PHASE_FIX_PLANS STATE

### Core Mandatory Rules (ALL orchestrator states must have these):

1. **🚨🚨🚨 R006** - ORCHESTRATOR NEVER WRITES CODE OR PERFORMS FILE OPERATIONS (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`
   - Criticality: BLOCKING - Automatic termination, 0% grade
   - Summary: NEVER write, copy, move, or manipulate ANY code files - delegate ALL to agents

2. **🔴🔴🔴 R287** - TODO PERSISTENCE COMPREHENSIVE (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-persistence-comprehensive.md`
   - Criticality: SUPREME - -20% to -100% penalty for violations
   - Summary: MUST save TODOs within 30s after write, every 10 messages, before transitions

3. **🔴🔴🔴 R288** - STATE FILE UPDATE REQUIREMENTS (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
   - Criticality: SUPREME - State updates required for all transitions
   - Summary: MUST update orchestrator-state-v3.json before EVERY state transition

4. **🔴🔴🔴 R322 Part A** - Mandatory Stop After Spawn States
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-before-state-transitions.md`
   - Criticality: SUPREME LAW - Must stop after spawning
   - Summary: ALL spawn states require STOP after spawning agents

### State-Specific Rules:

5. **🔴🔴🔴 R233** - Immediate Action On State Entry
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R233-all-states-immediate-action.md`
   - Criticality: SUPREME LAW - Must act immediately on entering state
   - Summary: WAITING states require active monitoring, not passive waiting

## State Purpose
Actively monitor Code Reviewer creating fix plans for phase-level integration failures. These are complex fixes that require careful planning before implementation.

## Critical Rules

### 🔴🔴🔴 RULE R322: MANDATORY STOP BEFORE STATE TRANSITION (SUPREME LAW)
- **STOP** and save state before ANY transition
- **READ** orchestrator-state-v3.json to verify current state
- **VALIDATE** next state exists in software-factory-3.0-state-machine.json
- **VIOLATION = IMMEDIATE FAILURE**

### 🔴🔴🔴 RULE R233: IMMEDIATE ACTION REQUIRED (SUPREME LAW)
- **NO PASSIVE WAITING** - Must actively check for completion
- **IMMEDIATE ACTION** - Start checking within first response
- **CONTINUOUS MONITORING_SWE_PROGRESS** - Check every 30-60 seconds
- **States are VERBS** - "WAITING" means "ACTIVELY CHECKING"

### 🔴🔴🔴 RULE R321: IMMEDIATE BACKPORT REQUIRED (SUPREME LAW)
- Phase fixes MUST be applied to source branches
- Integration branches are READ-ONLY for code
- Fix plans must target effort/wave branches
- Deferred backporting is DEPRECATED

### 🚨🚨🚨 RULE R290: STATE RULE VERIFICATION (BLOCKING)
- **MUST** verify this rules file exists and is loaded
- **MUST** acknowledge all rules before proceeding
- **MUST** validate state transitions against state machine

### 🚨🚨🚨 RULE R232: MONITOR STATE REQUIREMENTS (BLOCKING)
- **MUST** check TodoWrite for pending items BEFORE transition
- **MUST** process ALL pending items immediately
- **NO** "I will..." statements - only "I am..." with action
- **VIOLATION = AUTOMATIC FAILURE**

### ⚠️⚠️⚠️ RULE R282: PHASE FIX PLAN REQUIREMENTS (WARNING)
- Fix plans MUST identify root cause of phase failure
- Plans MUST specify which effort branches need fixes
- Plans MUST maintain phase architectural integrity
- Plans MUST be created as PHASE-FIX-PLAN-*.md files

### ⚠️⚠️⚠️ RULE R287: TODO PERSISTENCE (WARNING)
- **MUST** save TODOs every 10 messages or 15 minutes
- **MUST** save before state transition
- **MUST** commit and push TODO state

## Required Actions

1. **Initial Check (IMMEDIATE)**
   ```bash
   # Verify Code Reviewer was spawned for fix planning
   grep "spawned_agents" orchestrator-state-v3.json | grep -i "fix"
   
   # Check phase integration directory
   ls -la phase-*-integration/
   
   # Look for integration failure report
   cat phase-*/PHASE-INTEGRATE_WAVE_EFFORTS-REPORT.md | grep -i "status"
   ```

2. **Active Monitoring Loop**
   ```bash
   # Monitor for fix plan creation
   while true; do
     # Check for fix plan files
     FIX_PLANS=$(ls phase-*/PHASE-FIX-PLAN-*.md 2>/dev/null)
     
     if [ -n "$FIX_PLANS" ]; then
       echo "✓ Fix plans detected: $FIX_PLANS"
       break
     fi
     
     # Check Code Reviewer status
     if [ -f phase-*/fix-planning-status.yaml ]; then
       STATUS=$(grep "status:" phase-*/fix-planning-status.yaml | awk '{print $2}')
       echo "Fix planning status: $STATUS at $(date)"
       
       if [ "$STATUS" = "COMPLETED" ]; then
         echo "✓ Fix planning completed"
         break
       elif [ "$STATUS" = "BLOCKED" ]; then
         echo "✗ Fix planning blocked - escalating"
         # Transition to ERROR_RECOVERY
         break
       fi
     fi
     
     echo "Waiting for fix plans... checking again in 30s"
     sleep 30
   done
   ```

3. **Validate Fix Plans**
   ```bash
   # Validate each fix plan
   for PLAN in phase-*/PHASE-FIX-PLAN-*.md; do
     echo "Validating: $PLAN"
     
     # Check required sections
     grep -q "## Root Cause" "$PLAN" || echo "✗ Missing root cause"
     grep -q "## Affected Efforts" "$PLAN" || echo "✗ Missing effort list"
     grep -q "## Fix Strategy" "$PLAN" || echo "✗ Missing strategy"
     grep -q "## Implementation Steps" "$PLAN" || echo "✗ Missing steps"
     
     # Verify targets source branches (R321)
     if grep -q "integration" "$PLAN" | grep -v "after"; then
       echo "✗ ERROR: Fix plan targets integration branch (violates R321)"
     fi
   done
   ```

4. **Consolidate Fix Information**
   ```bash
   # Create summary of all fixes needed
   cat > phase-*/CONSOLIDATED-FIX-SUMMARY.md << EOF
   # Phase Fix Summary
   
   ## Fix Plans Created
   $(ls -1 phase-*/PHASE-FIX-PLAN-*.md)
   
   ## Affected Efforts
   $(grep "## Affected Efforts" phase-*/PHASE-FIX-PLAN-*.md -A 10 | grep "^-" | sort -u)
   
   ## Next Steps (R321 Compliance)
   1. Spawn SW Engineers to fix source branches
   2. Each fix in effort/wave branch, not integration
   3. Re-review after fixes
   4. Re-run phase integration with fixed sources
   EOF
   ```

5. **Check for Timeout**
   ```bash
   # Get spawn time
   SPAWN_TIME=$(grep "phase_fix_planning" orchestrator-state-v3.json -A 2 | grep "timestamp" | cut -d'"' -f2)
   
   # Calculate elapsed
   ELAPSED=$(($(date +%s) - $(date -d "$SPAWN_TIME" +%s)))
   
   # Timeout after 45 minutes (complex phase fixes)
   if [ $ELAPSED -gt 2700 ]; then
     echo "✗ Timeout waiting for phase fix plans"
     # Transition to ERROR_RECOVERY
   fi
   ```

6. **Update State When Complete**
   ```yaml
   current_state: WAITING_FOR_PHASE_FIX_PLANS
   phase_fixes:
     plans_created:
       - PHASE-FIX-PLAN-001.md
       - PHASE-FIX-PLAN-002.md
     affected_efforts:
       - effort-1
       - effort-3
     fix_strategy: immediate_backport
     plans_validated: true
   ```

## Transition Rules

### Valid Next States
- **ERROR_RECOVERY** - Fix plans ready, proceed to coordinate fixes
- **ERROR_RECOVERY** - Timeout or planning blocked (also goes here)

### Invalid Transitions
- ❌ Direct to SPAWN_SW_ENGINEERS (must go through ERROR_RECOVERY)
- ❌ Back to INTEGRATE_PHASE_WAVES without fixes
- ❌ Skipping fix implementation



## 🔴🔴🔴 STATE COMPLETION CHECKLIST 🔴🔴🔴

**Execute these steps IN ORDER to properly complete WAITING_FOR_PHASE_FIX_PLANS:**

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
PROPOSED_NEXT_STATE="NEXT_STATE"
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


## Common Violations to Avoid

1. **Passive waiting** - Violates R233, must actively monitor
2. **Not checking R321 compliance** - Plans target integration branch
3. **Missing timeout handling** - Waiting forever
4. **Not validating all plans** - Incomplete validation
5. **Forgetting TODO persistence** - Violates R287

## Phase Fix Complexity

Phase fixes are more complex than wave fixes because:
- Multiple waves may have interdependencies
- Architectural decisions may need revision
- Cross-wave conflicts require coordination
- Testing scope is much larger

This is why the state machine routes through ERROR_RECOVERY for coordination.

## Monitoring Pattern

```bash
# CORRECT: Active monitoring with comprehensive checks
echo "Monitoring phase fix planning at $(date)"
while true; do
  # Check for plans
  PLAN_COUNT=$(ls phase-*/PHASE-FIX-PLAN-*.md 2>/dev/null | wc -l)
  echo "Found $PLAN_COUNT fix plans so far"
  
  # Check status
  test -f phase-*/fix-planning-status.yaml && \
    grep "progress:" phase-*/fix-planning-status.yaml
  
  # Check for completion
  if [ $PLAN_COUNT -gt 0 ] && grep -q "COMPLETED" phase-*/fix-planning-status.yaml 2>/dev/null; then
    break
  fi
  
  sleep 30
done

# WRONG: Just waiting without checks
sleep 1800  # Wait 30 minutes
ls phase-*/PHASE-FIX-PLAN-*.md  # Check once at end
```

## Verification Commands

```bash
# Verify state entry
echo "Entered WAITING_FOR_PHASE_FIX_PLANS at $(date)"
echo "Phase integration failed - monitoring fix plan creation"

# Check failure context
grep "phase_integration" orchestrator-state-v3.json | grep "status"

# Active monitoring
for i in {1..90}; do  # 45 minute timeout
  if ls phase-*/PHASE-FIX-PLAN-*.md 2>/dev/null; then
    echo "✓ Fix plans found after $i checks"
    break
  fi
  echo "Check $i: No plans yet"
  sleep 30
done

# Validate R321 compliance
for PLAN in phase-*/PHASE-FIX-PLAN-*.md; do
  if grep -q "fix.*integration.*branch" "$PLAN"; then
    echo "✗ R321 VIOLATION: Plan suggests fixing integration branch!"
  fi
done
```

## References
- R232: rule-library/R232-enforcement-examples.md
- R233: rule-library/R233-all-states-immediate-action.md
- R282: rule-library/R282-phase-integration-protocol.md
- R287: rule-library/R287-todo-persistence-comprehensive.md
- R290: rule-library/R290-state-rule-reading-verification-supreme-law.md
- R321: rule-library/R321-immediate-backport-during-integration.md
- R322: rule-library/R322-mandatory-stop-before-state-transitions.md
