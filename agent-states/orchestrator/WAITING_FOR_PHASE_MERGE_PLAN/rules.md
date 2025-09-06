# WAITING_FOR_PHASE_MERGE_PLAN State Rules

## State Purpose
Actively monitor Code Reviewer creating phase merge plan for integrating all wave branches into the phase integration branch. This is critical for phase-level integration.

## Critical Rules

### 🔴🔴🔴 RULE R322: MANDATORY STOP BEFORE STATE TRANSITION (SUPREME LAW)
- **STOP** and save state before ANY transition
- **READ** orchestrator-state.yaml to verify current state
- **VALIDATE** next state exists in SOFTWARE-FACTORY-STATE-MACHINE.md
- **VIOLATION = IMMEDIATE FAILURE**

### 🔴🔴🔴 RULE R233: IMMEDIATE ACTION REQUIRED (SUPREME LAW)
- **NO PASSIVE WAITING** - Must actively check for completion
- **IMMEDIATE ACTION** - Start checking within first response
- **CONTINUOUS MONITORING** - Check every 30-60 seconds
- **States are VERBS** - "WAITING" means "ACTIVELY CHECKING"

### 🔴🔴🔴 RULE R285: MANDATORY PHASE INTEGRATION (SUPREME LAW)
- Phase integration MUST happen before phase assessment
- All waves MUST be integrated into phase branch
- Integration MUST follow sequential wave order
- Failed integration triggers IMMEDIATE_BACKPORT per R321

### 🚨🚨🚨 RULE R290: STATE RULE VERIFICATION (BLOCKING)
- **MUST** verify this rules file exists and is loaded
- **MUST** acknowledge all rules before proceeding
- **MUST** validate state transitions against state machine

### 🚨🚨🚨 RULE R232: MONITOR STATE REQUIREMENTS (BLOCKING)
- **MUST** check TodoWrite for pending items BEFORE transition
- **MUST** process ALL pending items immediately
- **NO** "I will..." statements - only "I am..." with action
- **VIOLATION = AUTOMATIC FAILURE**

### ⚠️⚠️⚠️ RULE R269: PHASE MERGE PLAN REQUIREMENTS (WARNING)
- Plan MUST be created as PHASE-MERGE-PLAN.md
- Plan MUST list all wave branches in sequential order
- Plan MUST specify integration strategy
- Plan MUST identify cross-wave dependencies

### ⚠️⚠️⚠️ RULE R287: TODO PERSISTENCE (WARNING)
- **MUST** save TODOs every 10 messages or 15 minutes
- **MUST** save before state transition
- **MUST** commit and push TODO state

## Required Actions

1. **Initial Check (IMMEDIATE)**
   ```bash
   # Verify Code Reviewer was spawned
   grep "SPAWN_CODE_REVIEWER_PHASE_MERGE_PLAN" orchestrator-state.yaml
   grep "spawned_agents" orchestrator-state.yaml | tail -5
   
   # Check phase integration directory exists
   ls -la phase-*-integration/
   
   # Verify all waves completed
   grep "waves_completed" orchestrator-state.yaml
   ```

2. **Active Monitoring Loop**
   ```bash
   # Monitor for phase merge plan
   while true; do
     # Check for plan file
     if [ -f phase-*/PHASE-MERGE-PLAN.md ]; then
       echo "✓ Phase merge plan detected at $(date)"
       break
     fi
     
     # Check Code Reviewer status
     if [ -f phase-*/merge-planning-status.yaml ]; then
       STATUS=$(grep "status:" phase-*/merge-planning-status.yaml | awk '{print $2}')
       PROGRESS=$(grep "progress:" phase-*/merge-planning-status.yaml | cut -d: -f2-)
       
       echo "Merge planning status: $STATUS"
       echo "Progress: $PROGRESS"
       
       if [ "$STATUS" = "COMPLETED" ]; then
         echo "✓ Merge planning completed"
         break
       elif [ "$STATUS" = "BLOCKED" ]; then
         echo "✗ Merge planning blocked"
         # Log blocker details
         grep "blocker:" phase-*/merge-planning-status.yaml
         # Transition to ERROR_RECOVERY
         break
       fi
     fi
     
     echo "Waiting for phase merge plan... checking again in 30s"
     sleep 30
   done
   ```

3. **Validate Phase Merge Plan**
   ```bash
   PLAN_FILE="phase-*/PHASE-MERGE-PLAN.md"
   
   if [ -f $PLAN_FILE ]; then
     echo "Validating phase merge plan..."
     
     # Check required sections
     grep -q "## Wave Merge Order" "$PLAN_FILE" || echo "✗ Missing wave order"
     grep -q "## Integration Strategy" "$PLAN_FILE" || echo "✗ Missing strategy"
     grep -q "## Wave Branches" "$PLAN_FILE" || echo "✗ Missing wave list"
     grep -q "## Dependencies" "$PLAN_FILE" || echo "✗ Missing dependencies"
     grep -q "## Conflict Analysis" "$PLAN_FILE" || echo "✗ Missing conflict analysis"
     
     # Verify all waves included
     WAVE_COUNT=$(grep -c "wave-.*-integration" "$PLAN_FILE")
     EXPECTED_WAVES=$(grep "total_waves:" orchestrator-state.yaml | awk '{print $2}')
     
     if [ "$WAVE_COUNT" -ne "$EXPECTED_WAVES" ]; then
       echo "✗ Plan missing waves: found $WAVE_COUNT, expected $EXPECTED_WAVES"
     else
       echo "✓ All $WAVE_COUNT waves included in plan"
     fi
     
     # Verify sequential order (R285)
     grep "wave-" "$PLAN_FILE" | grep -n "integration" | while read line; do
       echo "  $line"
     done
   fi
   ```

4. **Check Wave Readiness**
   ```bash
   # Verify all wave branches exist and are pushed
   for WAVE_DIR in wave-*-integration/; do
     if [ -d "$WAVE_DIR" ]; then
       cd "$WAVE_DIR"
       BRANCH=$(git branch --show-current)
       echo "Wave branch: $BRANCH"
       
       # Check remote exists
       if git ls-remote origin "$BRANCH" > /dev/null 2>&1; then
         echo "  ✓ Remote branch exists"
       else
         echo "  ✗ Remote branch missing!"
       fi
       
       # Check for integration report
       if [ -f WAVE-INTEGRATION-REPORT.md ]; then
         echo "  ✓ Wave integration complete"
       else
         echo "  ⚠ No integration report"
       fi
       cd ..
     fi
   done
   ```

5. **Check for Timeout**
   ```bash
   # Get spawn timestamp
   SPAWN_TIME=$(grep "SPAWN_CODE_REVIEWER_PHASE_MERGE_PLAN" orchestrator-state.yaml -A 10 | \
                grep "timestamp" | tail -1 | cut -d'"' -f2)
   
   # Calculate elapsed time
   CURRENT=$(date +%s)
   SPAWN_EPOCH=$(date -d "$SPAWN_TIME" +%s)
   ELAPSED=$((CURRENT - SPAWN_EPOCH))
   
   # Timeout after 30 minutes
   if [ $ELAPSED -gt 1800 ]; then
     echo "✗ Timeout: Phase merge planning took > 30 minutes"
     # Transition to ERROR_RECOVERY
   fi
   ```

6. **Update State When Complete**
   ```yaml
   current_state: WAITING_FOR_PHASE_MERGE_PLAN
   phase_integration:
     merge_plan: phase-X-integration/PHASE-MERGE-PLAN.md
     merge_plan_created: YYYY-MM-DD HH:MM:SS
     wave_count: X
     plan_validated: true
     ready_for_integration: true
   ```

## Transition Rules

### Valid Next States
- **SPAWN_INTEGRATION_AGENT_PHASE** - Plan created and validated
- **ERROR_RECOVERY** - Timeout or Code Reviewer blocked

### Invalid Transitions
- ❌ Direct to MONITORING_PHASE_INTEGRATION (must spawn first)
- ❌ Back to PHASE_INTEGRATION (already in that flow)
- ❌ Skipping to phase assessment (integration required first)

## Common Violations to Avoid

1. **Passive waiting** - Violates R233, must actively check
2. **Not validating wave count** - Plan missing waves
3. **Ignoring sequential order** - Waves out of order (R285)
4. **Missing dependency check** - Cross-wave conflicts
5. **Not checking remote branches** - Integration will fail

## Phase vs Wave Merge Plans

Key differences:
- **Scope**: Phase plans merge waves, wave plans merge efforts
- **Complexity**: Phase has cross-wave dependencies
- **Scale**: Phase typically 3-5 waves vs 5-10 efforts
- **Risk**: Phase failures affect entire project phase

## Monitoring Pattern

```bash
# CORRECT: Comprehensive active monitoring
echo "Starting phase merge plan monitoring at $(date)"
CHECKS=0
MAX_CHECKS=60  # 30 minutes

while [ $CHECKS -lt $MAX_CHECKS ]; do
  CHECKS=$((CHECKS + 1))
  echo "Check #$CHECKS at $(date)"
  
  # Check for plan
  if [ -f phase-*/PHASE-MERGE-PLAN.md ]; then
    echo "✓ Plan found!"
    # Validate immediately
    grep -c "wave-" phase-*/PHASE-MERGE-PLAN.md
    break
  fi
  
  # Check status
  if [ -f phase-*/merge-planning-status.yaml ]; then
    cat phase-*/merge-planning-status.yaml
  fi
  
  sleep 30
done

if [ $CHECKS -eq $MAX_CHECKS ]; then
  echo "✗ Timeout reached"
fi

# WRONG: Minimal checking
echo "Waiting for plan..."
sleep 900  # Wait 15 minutes
ls phase-*/PHASE-MERGE-PLAN.md || echo "No plan"
```

## Verification Commands

```bash
# Verify state entry
echo "===================="
echo "WAITING_FOR_PHASE_MERGE_PLAN"
echo "Entered at: $(date)"
echo "===================="

# Check context
echo "Current phase: $(grep current_phase orchestrator-state.yaml | awk '{print $2}')"
echo "Waves completed: $(grep -c "status: COMPLETE" orchestrator-state.yaml)"

# Monitor with timeout
timeout 1800 bash -c 'while [ ! -f phase-*/PHASE-MERGE-PLAN.md ]; do 
  echo "Checking... $(date +%H:%M:%S)"
  sleep 30
done && echo "✓ Plan created"' || echo "✗ Timeout"

# Final validation
if [ -f phase-*/PHASE-MERGE-PLAN.md ]; then
  echo "Plan size: $(wc -l phase-*/PHASE-MERGE-PLAN.md)"
  echo "Waves in plan: $(grep -c "wave-" phase-*/PHASE-MERGE-PLAN.md)"
fi
```

## References
- R232: rule-library/R232-monitor-state-requirements.md
- R233: rule-library/R233-all-states-immediate-action.md
- R269: rule-library/R269-merge-plan-requirements.md
- R285: rule-library/R285-mandatory-phase-integration.md
- R287: rule-library/R287-todo-persistence-comprehensive.md
- R290: rule-library/R290-state-rule-verification.md
- R322: rule-library/R322-mandatory-stop-before-transition.md