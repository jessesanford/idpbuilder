# 🔴🔴🔴 SUPREME RULE R348: Cascade State Transitions

## Criticality: SUPREME LAW
**Violation = -100% AUTOMATIC FAILURE**

## Description
This rule mandates the use of CASCADE_REINTEGRATION state whenever stale integrations are detected. The CASCADE_REINTEGRATION state is a PERSISTENT COORDINATOR that maintains control until ALL cascade recreations are complete, supporting MULTIPLE OVERLAPPING CASCADE CHAINS per R352. This ensures R327 compliance is absolute and unstoppable.

## 🔴🔴🔴 THE CASCADE STATE LAW 🔴🔴🔴

**ANY DETECTION OF STALE INTEGRATE_WAVE_EFFORTSS MUST IMMEDIATELY TRANSITION TO CASCADE_REINTEGRATION!**

### The Problem This Solves
```
❌ BROKEN FLOW (What was happening):
1. Fixes applied to effort branches
2. Wave integration marked stale
3. Orchestrator tries to continue without recreating phase/project
4. Phase integration contains old broken code
5. Project integration is unbuildable

✅ CORRECT FLOW (Enforced by CASCADE_REINTEGRATION):
1. Fixes applied to effort branches
2. Stale integrations detected
3. IMMEDIATE transition to CASCADE_REINTEGRATION
4. CASCADE_REINTEGRATION blocks all work
5. Forces recreation: wave → phase → project
6. Only when ALL fresh, allows continuation
```

## Core Requirements

### 1. MANDATORY TRANSITION TRIGGERS

These conditions MUST trigger transition to CASCADE_REINTEGRATION:

```bash
# From MONITORING_INTEGRATE_WAVE_EFFORTS
if [[ -n "$PENDING_CASCADES" ]] || [[ -n "$STALE_INTEGRATE_WAVE_EFFORTSS" ]]; then
    UPDATE_STATE="CASCADE_REINTEGRATION"
fi

# From WAVE_COMPLETE
if [[ -n "$FIXES_APPLIED" ]]; then
    NEXT_STATE="CASCADE_REINTEGRATION"
fi

# From COMPLETE_PHASE
if [[ -n "$PENDING_CASCADES" ]] || [[ "$PHASE_IS_STALE" == "true" ]]; then
    transition_to "CASCADE_REINTEGRATION"
fi
```

### 2. CASCADE_REINTEGRATION STATE BEHAVIOR

```yaml
CASCADE_REINTEGRATION:
  type: PERSISTENT_COORDINATOR
  characteristics:
    - Cannot be skipped
    - Maintains control after EVERY operation (R352)
    - Supports multiple overlapping cascade chains (R352)
    - Returns to CASCADE_REINTEGRATION after all integration states
    - Uses R350 dependency graph for cascade calculation
    - Uses R351 execution protocol for cascade execution
    - Uses R352 overlapping protocol for multiple chains
    - Processes cascades in dependency order
    - Handles both rebases (efforts) and recreations (integrations)
  
  entry_actions:
    - Load dependency graph (R350)
    - Calculate cascade chain (R350)
    - Create execution plan (R351)
    - Begin cascade execution (R351)
  
  allowed_transitions:
    - CASCADE_REINTEGRATION → INTEGRATE_WAVE_EFFORTS (recreate wave)
    - CASCADE_REINTEGRATION → INTEGRATE_PHASE_WAVES (recreate phase - SF 3.0)
    - CASCADE_REINTEGRATION → PROJECT_INTEGRATE_WAVE_EFFORTS (recreate project)
    - CASCADE_REINTEGRATION → CASCADE_REINTEGRATION (more cascades)
    - CASCADE_REINTEGRATION → REVIEW_WAVE_INTEGRATION (all done)

  return_transitions_per_R352:
    - INTEGRATE_WAVE_EFFORTS → CASCADE_REINTEGRATION (when cascade_mode=true)
    - INTEGRATE_PHASE_WAVES → CASCADE_REINTEGRATION (when cascade_mode=true - SF 3.0)
    - PROJECT_INTEGRATE_WAVE_EFFORTS → CASCADE_REINTEGRATION (when cascade_mode=true)
    - MONITORING_INTEGRATE_WAVE_EFFORTS → CASCADE_REINTEGRATION (when cascade_mode=true)
    - All integration states → CASCADE_REINTEGRATION (when cascade_mode=true)
  
  forbidden_transitions:
    - CASCADE_REINTEGRATION → WAVE_COMPLETE (❌ CANNOT skip)
    - CASCADE_REINTEGRATION → COMPLETE_PHASE (❌ CANNOT skip)
    - CASCADE_REINTEGRATION → PROJECT_DONE (❌ CANNOT skip)
    - CASCADE_REINTEGRATION → Any other state (❌ BLOCKED)
  
  execution_flow:
    1. Calculate complete cascade chain using R350 dependency graph
    2. Create execution plan with R351 protocol
    3. Execute rebases for dependent efforts
    4. Execute recreations for stale integrations
    5. Verify all dependencies satisfied
    6. Exit only when cascade fully complete
```

### 3. CASCADE TRACKING REQUIREMENTS

The state MUST maintain cascade tracking:

```json
"cascade_coordination": {
  "cascade_mode": true,
  "persistent_coordination": true,
  "active_cascade_chains": [
    {
      "chain_id": "cascade_001",
      "trigger": {
        "type": "fix_applied",
        "location": "phase1/wave1/effort1",
        "timestamp": "2025-01-14T10:00:00Z",
        "fix_ids": ["fix_001", "fix_002"]
      },
      "status": "in_progress",
      "operations": [
        {
          "type": "recreate",
          "target": "phase1-wave1-integration",
          "status": "completed"
        }
      ],
      "started_at": "2025-01-14T10:00:00Z"
    },
    {
      "chain_id": "cascade_002",
      "trigger": {
        "type": "fix_applied",
        "location": "phase2/wave1/effort2",
        "timestamp": "2025-01-14T11:00:00Z",
        "fix_ids": ["fix_003"]
      },
      "status": "pending",
      "operations": [],
      "started_at": "2025-01-14T11:00:00Z"
    }
  ],
  "cascade_complete_when": {
    "all_chains_complete": true,
    "no_pending_fixes": true,
    "project_integration_fresh": true,
    "no_new_fixes_detected": true
  }
},
"stale_integration_tracking": {
  "staleness_cascade": [...],
  "stale_integrations": [...]
}
```

### 4. CASCADE COMPLETION VALIDATION

Before exiting CASCADE_REINTEGRATION:

```bash
# ALL of these must be true (R352 compliant):
validate_cascade_complete() {
    # 1. No active cascade chains (R352)
    ACTIVE_CHAINS=$(jq -r '
        .cascade_coordination.active_cascade_chains[]? |
        select(.status != "completed" and .status != "merged_into") |
        .chain_id' orchestrator-state-v3.json)
    [[ -z "$ACTIVE_CHAINS" ]] || return 1
    
    # 2. No pending fixes without cascade chains (R352)
    PENDING_FIXES=$(jq -r '
        .cascade_coordination.pending_fixes | 
        to_entries[] | 
        select(.value.cascade_chain == null) |
        .key' orchestrator-state-v3.json)
    [[ -z "$PENDING_FIXES" ]] || return 1
    
    # 3. No stale integrations
    STALE=$(jq '.stale_integration_tracking.stale_integrations[]? | 
                select(.recreation_completed == false)' orchestrator-state-v3.json)
    [[ -z "$STALE" ]] || return 1
    
    # 4. All integrations fresh
    for level in wave phase project; do
        IS_STALE=$(jq -r ".current_${level}_integration.is_stale // false" orchestrator-state-v3.json)
        [[ "$IS_STALE" == "false" ]] || return 1
    done
    
    # 5. Check for new fixes that arrived during cascade (R352)
    check_for_new_fixes
    NEW_CHAINS=$(jq '.cascade_coordination.active_cascade_chains | 
                     map(select(.status == "pending")) | length' orchestrator-state-v3.json)
    [[ "$NEW_CHAINS" -eq 0 ]] || return 1
    
    return 0
}
```

## State Machine Integration

### Required Updates to State Machine

```yaml
# Add CASCADE_REINTEGRATION to valid states
VALID_STATES:
  - CASCADE_REINTEGRATION  # Trap state for cascade enforcement

# Add transitions
TRANSITIONS:
  MONITORING_INTEGRATE_WAVE_EFFORTS:
    - to: CASCADE_REINTEGRATION
      when: stale_integrations_detected
      
  WAVE_COMPLETE:
    - to: CASCADE_REINTEGRATION
      when: fixes_applied_during_wave
      
  COMPLETE_PHASE:
    - to: CASCADE_REINTEGRATION
      when: pending_cascades_exist
      
  CASCADE_REINTEGRATION:
    - to: INTEGRATE_WAVE_EFFORTS
      when: recreating_wave_integration
    - to: INTEGRATE_PHASE_WAVES
      when: recreating_phase_integration
    - to: PROJECT_INTEGRATE_WAVE_EFFORTS
      when: recreating_project_integration
    - to: CASCADE_REINTEGRATION
      when: more_cascades_pending
    - to: REVIEW_WAVE_INTEGRATION
      when: all_cascades_complete
```

## Detection Points

### 1. During Monitoring States
```bash
# In MONITORING_INTEGRATE_WAVE_EFFORTS, MONITORING_EFFORT_FIXES, etc.
./utilities/stale-integration-manager.sh check || {
    echo "🔴 Stale integrations detected!"
    transition_to "CASCADE_REINTEGRATION"
}
```

### 2. Before Major Transitions
```bash
# Before WAVE_COMPLETE → INTEGRATE_WAVE_EFFORTS
# Before COMPLETE_PHASE → INIT/PROJECT_DONE
check_for_cascades() {
    NEEDS_CASCADE=$(./utilities/stale-integration-manager.sh check > /dev/null; echo $?)
    if [[ "$NEEDS_CASCADE" -ne 0 ]]; then
        return 1  # Block transition, need CASCADE_REINTEGRATION
    fi
    return 0
}
```

### 3. After Fix Application
```bash
# When fixes are tracked
./utilities/stale-integration-manager.sh track-fix "$FIX_ID" ...
# This automatically creates cascade requirements
```

## Common Violations (ALL RESULT IN AUTOMATIC FAILURE)

### ❌ VIOLATION 1: Skipping CASCADE_REINTEGRATION
```bash
# WRONG:
if stale_detected; then
    recreate_wave_integration  # Just recreating wave
    continue_to_next_state
fi
```

### ✅ CORRECTION 1: Mandatory CASCADE_REINTEGRATION
```bash
# RIGHT:
if stale_detected; then
    transition_to "CASCADE_REINTEGRATION"  # Let state handle ALL cascades
fi
```

### ❌ VIOLATION 2: Partial Cascade Handling
```bash
# WRONG:
# In CASCADE_REINTEGRATION
recreate_wave_integration
transition_to "REVIEW_WAVE_INTEGRATION"  # Skipping phase/project!
```

### ✅ CORRECTION 2: Complete Cascade
```bash
# RIGHT:
# CASCADE_REINTEGRATION handles ALL levels
while cascades_pending; do
    process_next_cascade
    transition_to_appropriate_integration_state
done
```

### ❌ VIOLATION 3: Ignoring Cascade Detection
```bash
# WRONG:
# In WAVE_COMPLETE
transition_to "INTEGRATE_WAVE_EFFORTS"  # Not checking for fixes!
```

### ✅ CORRECTION 3: Always Check
```bash
# RIGHT:
if fixes_applied || cascades_pending; then
    transition_to "CASCADE_REINTEGRATION"
else
    transition_to "INTEGRATE_WAVE_EFFORTS"
fi
```

## Real-World Example: Test 5 Execution

### Test 5: Cross-Container Cascade (Actual Execution)

Test 5 demonstrates the CASCADE_REINTEGRATION persistent coordinator pattern with actual state transitions from production test execution.

**Scenario:** Wave-level bugs require fixes that cascade to phase integration

**Complete State Sequence:**

```
1. FIX_WAVE_UPSTREAM_BUGS → CASCADE_REINTEGRATION
   Reason: "Bug fixes complete in Wave 1 efforts, triggering cascade reintegration"

2. CASCADE_REINTEGRATION → INTEGRATE_WAVE_EFFORTS
   Reason: "Beginning bottom-up reintegration - wave_1_1 integration after fixes"

3. MONITORING_INTEGRATE_WAVE_EFFORTS → CASCADE_REINTEGRATION
   Reason: "wave_1_1 reintegration complete, returning to CASCADE_REINTEGRATION to process wave_1_2"
   ☝️ R352 PERSISTENT COORDINATOR - Returns after each operation

4. CASCADE_REINTEGRATION → INTEGRATE_WAVE_EFFORTS
   Reason: "Beginning bottom-up reintegration - wave_1_2 integration"

5. INTEGRATE_WAVE_EFFORTS → CASCADE_REINTEGRATION
   Reason: "wave_1_2 reintegration complete, processing completion"
   ☝️ R352 PERSISTENT COORDINATOR - Returns again

6. CASCADE_REINTEGRATION → INTEGRATE_PHASE_WAVES
   Reason: "Wave-level reintegrations complete, proceeding to phase_1 reintegration"
```

**Key Observations:**

1. **Persistent Coordinator Pattern (R352):**
   - CASCADE_REINTEGRATION appeared 6 times (3 entries, 3 exits)
   - Pattern: Entry → Exit → Return → Exit → Return → Exit
   - Maintains control throughout entire cascade operation

2. **Bottom-Up Reintegration (R327):**
   - wave_1_1 reintegrated first
   - wave_1_2 reintegrated second
   - phase_1 reintegrated last
   - Project-level would follow (not shown in this test)

3. **R348 Compliance:**
   - Immediate transition to CASCADE_REINTEGRATION on stale detection
   - No direct integration-to-integration transitions
   - All operations routed through CASCADE_REINTEGRATION

**State History JSON (Actual from Test 5):**

```json
{
  "from_state": "CASCADE_REINTEGRATION",
  "to_state": "INTEGRATE_WAVE_EFFORTS",
  "timestamp": "2025-10-19T22:36:35Z",
  "validated_by": "state-manager",
  "reason": "CASCADE-001: Beginning bottom-up reintegration - wave_1_1 integration after fixes (R327/R352)"
}
```

**Test Location:** `tests/runtime-test-05-cross-container-cascade.sh`

**Validation Report:** See `TEST-5-CASCADE-REINTEGRATE_WAVE_EFFORTS-VALIDATION-REPORT.md` for complete analysis

---

## Grading Impact

### AUTOMATIC FAILURE (-100%)
- Detecting stale integrations but not transitioning to CASCADE_REINTEGRATION
- Exiting CASCADE_REINTEGRATION with pending cascades
- Skipping cascade levels
- Allowing phase/project completion with stale integrations

### MAJOR VIOLATIONS (-50%)
- Not checking for cascades at transition points
- Manual cascade handling instead of using CASCADE_REINTEGRATION
- Incomplete cascade tracking

### COMPLIANCE BONUS (+30%)
- Proper CASCADE_REINTEGRATION usage
- Complete cascade tracking
- Clean state transitions
- All integrations verified fresh

## Relationship to Other Rules

### Depends on
- R327: Mandatory Re-Integration After Fixes (defines cascade requirement)
- R350: Complete Cascade Dependency Graph (provides dependency tracking)
- R351: Cascade Execution Protocol (provides execution mechanism)
- R352: Overlapping Cascade Protocol (multiple cascade chains)
- R346: State Metadata Synchronization (tracking structure)

### Enables
- R291: Integration Demo Requirement (ensures fresh code for demos)
- R266: Project Bug Detection (fresh integrations contain fixes)
- R328: Integration Freshness Validation (validates cascade results)

### Works with
- R322: Mandatory Stop After State Transitions
- R206: State Machine Validation
- R337: Base Branch Single Source Truth

## Quick Reference

### Check if CASCADE_REINTEGRATION Needed
```bash
needs_cascade() {
    # Any of these trigger CASCADE_REINTEGRATION:
    [[ -n "$(jq '.stale_integration_tracking.staleness_cascade[]? | 
             select(.cascade_status != "completed")' orchestrator-state-v3.json)" ]] && return 0
    
    [[ -n "$(jq '.stale_integration_tracking.stale_integrations[]? | 
             select(.recreation_completed == false)' orchestrator-state-v3.json)" ]] && return 0
    
    [[ "$(jq -r '.current_wave_integration.is_stale // false' orchestrator-state-v3.json)" == "true" ]] && return 0
    
    return 1
}

if needs_cascade; then
    echo "🔴 CASCADE_REINTEGRATION required!"
    transition_to "CASCADE_REINTEGRATION"
fi
```

### Force CASCADE_REINTEGRATION
```bash
# Emergency cascade trigger
force_cascade() {
    echo "🔴 R348: Forcing CASCADE_REINTEGRATION"
    jq '.state_machine.current_state = "CASCADE_REINTEGRATION"' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
    jq '.transition_reason = "R348 enforcement - cascade required"' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
    git add orchestrator-state-v3.json
    git commit -m "state: force CASCADE_REINTEGRATION per R348"
    git push
}
```

## Remember

**"CASCADE_REINTEGRATION is a TRAP STATE - Once in, you MUST complete ALL cascades"**
**"Stale detection = Immediate CASCADE_REINTEGRATION"**
**"No shortcuts, no skips, no exceptions"**
**"The cascade chain is UNBREAKABLE"**

### 🔴🔴🔴 THE CASCADE STATE MANTRA 🔴🔴🔴
```
When staleness is detected,
CASCADE_REINTEGRATION is selected.
It traps you in its grasp so tight,
Until ALL integrations are made right.

Wave, then Phase, then Project too,
Each cascade must see you through.
No escape, no skip, no cheat,
Until the cascade is complete!
```

The goal: Make cascade enforcement AUTOMATIC and UNSTOPPABLE through the CASCADE_REINTEGRATION trap state.