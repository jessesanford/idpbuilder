# WAITING_FOR_PHASE_FIX_PLANS State Rules

## State Purpose
Actively monitor Code Reviewer creating fix plans for phase-level integration failures. These are complex fixes that require careful planning before implementation.

## Critical Rules

### 🔴🔴🔴 RULE R322: MANDATORY STOP BEFORE STATE TRANSITION (SUPREME LAW)
- **STOP** and save state before ANY transition
- **READ** orchestrator-state.json to verify current state
- **VALIDATE** next state exists in SOFTWARE-FACTORY-STATE-MACHINE.md
- **VIOLATION = IMMEDIATE FAILURE**

### 🔴🔴🔴 RULE R233: IMMEDIATE ACTION REQUIRED (SUPREME LAW)
- **NO PASSIVE WAITING** - Must actively check for completion
- **IMMEDIATE ACTION** - Start checking within first response
- **CONTINUOUS MONITORING** - Check every 30-60 seconds
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
   grep "spawned_agents" orchestrator-state.json | grep -i "fix"
   
   # Check phase integration directory
   ls -la phase-*-integration/
   
   # Look for integration failure report
   cat phase-*/PHASE-INTEGRATION-REPORT.md | grep -i "status"
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
   SPAWN_TIME=$(grep "phase_fix_planning" orchestrator-state.json -A 2 | grep "timestamp" | cut -d'"' -f2)
   
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
- ❌ Direct to SPAWN_ENGINEERS_FOR_FIXES (must go through ERROR_RECOVERY)
- ❌ Back to PHASE_INTEGRATION without fixes
- ❌ Skipping fix implementation

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
grep "phase_integration" orchestrator-state.json | grep "status"

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
- R232: rule-library/R232-monitor-state-requirements.md
- R233: rule-library/R233-all-states-immediate-action.md
- R282: rule-library/R282-phase-fix-requirements.md
- R287: rule-library/R287-todo-persistence-comprehensive.md
- R290: rule-library/R290-state-rule-verification.md
- R321: rule-library/R321-immediate-backport-during-integration.md
- R322: rule-library/R322-mandatory-stop-before-transition.md