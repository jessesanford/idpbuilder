# WAITING_FOR_MERGE_PLAN State Rules

## State Purpose
Actively monitor Code Reviewer creating wave merge plan. Check for completion and validate the merge plan meets requirements.

## Critical Rules

### 🛑🛑🛑 RULE R322: MANDATORY CHECKPOINT BEFORE SPAWN_INTEGRATION_AGENT (SUPREME LAW) 🛑🛑🛑

**THIS IS A CRITICAL R322 CHECKPOINT STATE!**

When transitioning from WAITING_FOR_MERGE_PLAN → SPAWN_INTEGRATION_AGENT:
- **MUST STOP** to allow user review of WAVE-MERGE-PLAN.md
- **MUST UPDATE** state file to SPAWN_INTEGRATION_AGENT before stopping
- **MUST DISPLAY** checkpoint message with plan location
- **MUST EXIT** cleanly to preserve context
- **VIOLATION = -100% IMMEDIATE FAILURE**

See: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-before-state-transitions.md`

### 🔴🔴🔴 RULE R233: IMMEDIATE ACTION REQUIRED (SUPREME LAW)
- **NO PASSIVE WAITING** - Must actively check for completion
- **IMMEDIATE ACTION** - Start checking within first response
- **CONTINUOUS MONITORING** - Check every 30-60 seconds
- **States are VERBS** - "WAITING" means "ACTIVELY CHECKING"

### 🚨🚨🚨 RULE R290: STATE RULE VERIFICATION (BLOCKING)
- **MUST** verify this rules file exists and is loaded
- **MUST** acknowledge all rules before proceeding
- **MUST** validate state transitions against state machine

### 🚨🚨🚨 RULE R232: MONITOR STATE REQUIREMENTS (BLOCKING)
- **MUST** check TodoWrite for pending items BEFORE transition
- **MUST** process ALL pending items immediately
- **NO** "I will..." statements - only "I am..." with action
- **VIOLATION = AUTOMATIC FAILURE**

### ⚠️⚠️⚠️ RULE R269: MERGE PLAN VALIDATION (WARNING)
- Merge plan MUST exist as WAVE-MERGE-PLAN.md
- Plan MUST list all effort branches in order
- Plan MUST specify merge strategy
- Plan MUST identify potential conflicts

### ⚠️⚠️⚠️ RULE R287: TODO PERSISTENCE (WARNING)
- **MUST** save TODOs every 10 messages or 15 minutes
- **MUST** save before state transition
- **MUST** commit and push TODO state

## Required Actions

1. **Initial Check (IMMEDIATE)**
   ```bash
   # Check if Code Reviewer was spawned
   grep "spawned_agents" orchestrator-state.json
   
   # Verify wave integration directory
   ls -la wave-*-integration/
   ```

2. **Active Monitoring Loop**
   ```bash
   # Get integration workspace location
   PHASE=$(jq '.current_phase' orchestrator-state.json)
   WAVE=$(jq '.current_wave' orchestrator-state.json)
   INTEGRATION_DIR="/efforts/phase${PHASE}/wave${WAVE}/integration-workspace"
   
   # Check for merge plan existence in the integration workspace
   while true; do
     if [ -f "${INTEGRATION_DIR}/WAVE-MERGE-PLAN.md" ]; then
       echo "✓ Merge plan detected at $(date)"
       echo "📍 Location: ${INTEGRATION_DIR}/WAVE-MERGE-PLAN.md"
       break
     fi
     
     # Check agent status (if status file exists in integration workspace)
     if [ -f "${INTEGRATION_DIR}/status.yaml" ]; then
       STATUS=$(grep "status:" "${INTEGRATION_DIR}/status.yaml" | awk '{print $2}')
       echo "Code Reviewer status: $STATUS"
       
       if [ "$STATUS" = "COMPLETED" ]; then
         echo "✓ Code Reviewer completed"
         break
       elif [ "$STATUS" = "BLOCKED" ]; then
         echo "✗ Code Reviewer blocked - need intervention"
         # Transition to ERROR_RECOVERY
         break
       fi
     fi
     
     echo "Waiting for merge plan in ${INTEGRATION_DIR}... checking again in 30s"
     sleep 30
   done
   ```

3. **Validate Merge Plan Contents**
   ```bash
   # Once plan exists, validate it
   PLAN_FILE="${INTEGRATION_DIR}/WAVE-MERGE-PLAN.md"
   
   # Check plan has required sections
   grep -q "## Merge Order" "$PLAN_FILE" || echo "✗ Missing merge order"
   grep -q "## Effort Branches" "$PLAN_FILE" || echo "✗ Missing effort list"
   grep -q "## Merge Strategy" "$PLAN_FILE" || echo "✗ Missing strategy"
   grep -q "## Potential Conflicts" "$PLAN_FILE" || echo "✗ Missing conflict analysis"
   
   # Count effort branches in plan
   EFFORT_COUNT=$(grep -c "effort-" "$PLAN_FILE")
   echo "Plan includes $EFFORT_COUNT efforts"
   ```

4. **Check for Timeout**
   ```bash
   # Get spawn time
   SPAWN_TIME=$(grep "spawned_agents" orchestrator-state.json -A 5 | grep "timestamp" | tail -1 | cut -d'"' -f2)
   
   # Calculate elapsed time
   ELAPSED=$(($(date +%s) - $(date -d "$SPAWN_TIME" +%s)))
   
   # Timeout after 30 minutes
   if [ $ELAPSED -gt 1800 ]; then
     echo "✗ Timeout waiting for merge plan"
     # Transition to ERROR_RECOVERY
   fi
   ```

5. **Update State When Complete**
   ```yaml
   current_state: WAITING_FOR_MERGE_PLAN
   wave_integration:
     merge_plan: wave-X-integration/WAVE-MERGE-PLAN.md
     merge_plan_created: YYYY-MM-DD HH:MM:SS
     plan_validation: PASSED
   ```

## Transition Rules

### 🛑🛑🛑 R322 MANDATORY CHECKPOINT - DO NOT SKIP! 🛑🛑🛑

**THIS TRANSITION REQUIRES A MANDATORY USER CHECKPOINT!**

When merge plan is ready and validated:

```bash
# STEP 1: Validate merge plan exists
if [ ! -f "${INTEGRATION_DIR}/WAVE-MERGE-PLAN.md" ]; then
    echo "❌ Cannot transition - no merge plan found!"
    exit 1
fi

# STEP 2: Update state file to SPAWN_INTEGRATION_AGENT
echo "📝 R322: Updating state for checkpoint..."
jq '.current_state = "SPAWN_INTEGRATION_AGENT"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
jq '.previous_state = "WAITING_FOR_MERGE_PLAN"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
jq '.transition_time = "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
jq '.wave_integration.merge_plan = "integration-workspace/WAVE-MERGE-PLAN.md"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json

# STEP 3: Save TODOs per R287
save_todos "R322_CHECKPOINT_MERGE_PLAN_READY"

# STEP 4: Commit state changes
git add orchestrator-state.json todos/*.todo
git commit -m "state: R322 checkpoint - merge plan ready for user review"
git push

# STEP 5: Display R322 checkpoint message
cat << 'EOF'
## 🛑 R322 MANDATORY CHECKPOINT - USER REVIEW REQUIRED

### ✅ Merge Plan Created Successfully:
- Location: integration-workspace/WAVE-MERGE-PLAN.md
- Created by: Code Reviewer
- Status: Ready for review

### 📊 State Transition Ready:
- Current State: WAITING_FOR_MERGE_PLAN ✅
- Next State: SPAWN_INTEGRATION_AGENT (pending user approval)

### 📋 CRITICAL: Review Required Before Execution
The merge plan has been created but NOT executed. Please review:
1. Merge order and dependencies
2. Conflict resolution strategy
3. Branch selection (splits vs originals)

### ⏸️ STOPPED FOR USER REVIEW
DO NOT proceed without reviewing the merge plan!
To continue after review: /continue-orchestrating
EOF

# STEP 6: EXIT CLEANLY (R322 MANDATORY)
exit 0  # STOP HERE - DO NOT CONTINUE!
```

### ❌ R322 VIOLATION - AUTOMATIC FAILURE:
```bash
# ❌❌❌ NEVER DO THIS - IMMEDIATE -100% FAILURE!
echo "Merge plan ready, spawning integration agent..."
/spawn integration-agent  # R322 VIOLATION!
```

### Valid Next States
- **SPAWN_INTEGRATION_AGENT** - After R322 checkpoint and user review
- **ERROR_RECOVERY** - Timeout or Code Reviewer blocked

### Invalid Transitions
- ❌ Skipping to MONITORING_INTEGRATION without spawn
- ❌ Going back to spawn states
- ❌ Transitioning without merge plan
- ❌ Stopping without updating current_state (CAUSES LOOPS!)

## Common Violations to Avoid

1. **Passive waiting** - Violates R233, must actively check
2. **Not checking TodoWrite** - Violates R232 before transition
3. **Missing plan validation** - Proceeding with invalid plan
4. **Ignoring timeouts** - Waiting forever for blocked agent
5. **Not saving progress** - Violates R287 persistence

## Monitoring Pattern

```bash
# CORRECT: Active monitoring with status updates
echo "Checking for merge plan at $(date)"
while [ ! -f wave-*/WAVE-MERGE-PLAN.md ]; do
  echo "Plan not ready, checking again in 30s..."
  # Also check agent status
  test -f wave-*/status.yaml && cat wave-*/status.yaml
  sleep 30
done
echo "✓ Merge plan created"

# WRONG: Passive waiting
echo "Waiting for Code Reviewer to finish..."
sleep 600  # Just sleeping without checking
```

## Verification Commands

```bash
# Verify state entry
echo "Entered WAITING_FOR_MERGE_PLAN at $(date)"
echo "Starting active monitoring for merge plan"

# Check spawn record
grep "spawned_agents" orchestrator-state.json

# Monitor loop
for i in {1..60}; do
  test -f wave-*/WAVE-MERGE-PLAN.md && break
  echo "Check $i: Plan not ready"
  sleep 30
done

# Validate plan when found
if [ -f wave-*/WAVE-MERGE-PLAN.md ]; then
  echo "✓ Merge plan found"
  wc -l wave-*/WAVE-MERGE-PLAN.md
else
  echo "✗ Timeout - no merge plan after 30 minutes"
fi
```

## References
- R232: rule-library/R232-monitor-state-requirements.md
- R233: rule-library/R233-all-states-immediate-action.md
- R269: rule-library/R269-merge-plan-requirements.md
- R287: rule-library/R287-todo-persistence-comprehensive.md
- R290: rule-library/R290-state-rule-verification.md
- R322: rule-library/R322-mandatory-stop-before-transition.md