# 🔴🔴🔴 SUPREME RULE R352: Overlapping Cascade Protocol

## Criticality: SUPREME LAW
**Violation = -100% AUTOMATIC FAILURE**

## Description
This rule mandates support for OVERLAPPING CASCADE CHAINS - multiple cascade operations running simultaneously from different trigger points. When fixes arrive at different levels while cascades are already in progress, the system MUST handle all cascades correctly, ensuring NO FIX IS EVER LOST. CASCADE_REINTEGRATION acts as a persistent coordinator that maintains control throughout ALL overlapping cascades.

## 🔴🔴🔴 THE OVERLAPPING CASCADE LAW 🔴🔴🔴

**CASCADE_REINTEGRATION IS A PERSISTENT COORDINATOR THAT NEVER RELEASES CONTROL UNTIL ALL FIXES REACH PROJECT LEVEL!**

### The Problem This Solves
```
❌ BROKEN FLOW (Sequential-only cascades):
1. Fix in wave1 triggers cascade
2. System locks into single cascade chain
3. Fix arrives in wave2 during cascade
4. New fix is IGNORED or LOST
5. Project integration missing critical fixes
6. BUILD FAILURE

✅ CORRECT FLOW (Overlapping cascades):
1. Fix in wave1 triggers cascade chain #1
2. CASCADE_REINTEGRATION maintains control
3. Fix arrives in wave2, triggers cascade chain #2
4. Both chains tracked independently
5. Chains merge when they converge
6. ALL fixes reach project integration
7. BUILD PROJECT_DONE with ALL fixes
```

## Core Requirements

### 1. CASCADE_REINTEGRATION PERSISTENCE

**CASCADE_REINTEGRATION maintains control after EVERY operation:**

```yaml
CASCADE_REINTEGRATION:
  persistent_control:
    - After wave effort rebases → Return to CASCADE_REINTEGRATION
    - After wave reintegration → Return to CASCADE_REINTEGRATION
    - After phase reintegration → Return to CASCADE_REINTEGRATION
    - After project reintegration → Return to CASCADE_REINTEGRATION
    - After any integration state → Return to CASCADE_REINTEGRATION
  
  control_release_conditions:
    ALL_OF:
      - All cascade chains status = "completed"
      - No pending fixes in any effort
      - No stale integrations remain
      - Project integration is fresh
      - No new fixes detected during last operation
```

### 2. MULTIPLE CASCADE CHAIN TRACKING

**State must support multiple active cascade chains:**

```json
{
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
            "type": "rebase",
            "target": "phase1/wave1/effort2",
            "base": "phase1/wave1/effort1",
            "status": "completed"
          },
          {
            "type": "recreate",
            "target": "phase1-wave1-integration",
            "status": "pending"
          }
        ],
        "started_at": "2025-01-14T10:00:00Z",
        "last_operation_at": "2025-01-14T10:15:00Z"
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
    "pending_fixes": {
      "phase1/wave1/effort1": {
        "fix_ids": ["fix_001", "fix_002"],
        "applied_at": "2025-01-14T10:00:00Z",
        "cascade_chain": "cascade_001"
      },
      "phase2/wave1/effort2": {
        "fix_ids": ["fix_003"],
        "applied_at": "2025-01-14T11:00:00Z",
        "cascade_chain": "cascade_002"
      }
    },
    "cascade_complete_when": {
      "all_chains_complete": true,
      "no_pending_fixes": true,
      "project_integration_fresh": true,
      "no_new_fixes_detected": true
    }
  }
}
```

### 3. CASCADE CHAIN MERGING

**When cascade chains converge, they MUST be merged:**

```bash
merge_cascade_chains() {
    local CHAIN1=$1
    local CHAIN2=$2
    local MERGE_POINT=$3
    
    echo "🔀 Merging cascade chains at $MERGE_POINT"
    
    # Create merged chain
    jq --arg c1 "$CHAIN1" --arg c2 "$CHAIN2" --arg mp "$MERGE_POINT" '
        .cascade_coordination.active_cascade_chains |= map(
            if .chain_id == $c1 then
                . + {
                    merged_with: [$c2],
                    merge_point: $mp,
                    includes_fixes: (
                        .trigger.fix_ids + 
                        (.cascade_coordination.active_cascade_chains[] | 
                         select(.chain_id == $c2) | .trigger.fix_ids)
                    )
                }
            elif .chain_id == $c2 then
                . + {status: "merged_into", merged_into: $c1}
            else . end
        )
    ' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
}
```

### 4. NEW FIX DETECTION DURING CASCADE

**CASCADE_REINTEGRATION must check for new fixes after EVERY operation:**

```bash
check_for_new_fixes() {
    echo "🔍 Checking for new fixes that arrived during cascade operation..."
    
    # Get timestamp of last cascade operation
    LAST_OP_TIME=$(jq -r '.cascade_coordination.last_operation_completed_at' orchestrator-state-v3.json)
    
    # Check git logs for new fixes since last operation
    for effort_dir in /efforts/phase*/wave*/effort-*; do
        if [[ -d "$effort_dir" ]]; then
            cd "$effort_dir"
            branch=$(git branch --show-current)
            
            # Check for new commits since last operation
            NEW_FIXES=$(git log --since="$LAST_OP_TIME" --oneline --grep="fix:" | wc -l)
            
            if [[ "$NEW_FIXES" -gt 0 ]]; then
                echo "⚠️ Found $NEW_FIXES new fixes in $branch"
                
                # Create new cascade chain for these fixes
                create_new_cascade_chain "$branch" "$NEW_FIXES"
            fi
        fi
    done
}

# Called after EVERY cascade operation completes
after_cascade_operation() {
    local OP_TYPE=$1
    local OP_TARGET=$2
    
    echo "✅ Completed $OP_TYPE on $OP_TARGET"
    
    # Update timestamp
    jq --arg ts "$(date -Iseconds)" '
        .cascade_coordination.last_operation_completed_at = $ts
    ' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
    
    # Check for new fixes
    check_for_new_fixes
    
    # Return to CASCADE_REINTEGRATION for next operation
    echo "🔄 Returning control to CASCADE_REINTEGRATION"
    jq '.state_machine.current_state = "CASCADE_REINTEGRATION"' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
}
```

### 5. CASCADE COMPLETION LOGIC

**CASCADE_REINTEGRATION only exits when ALL conditions are met:**

```bash
can_exit_cascade_reintegration() {
    echo "🔍 Checking if all cascade operations are complete..."
    
    # 1. Check all cascade chains are complete
    ACTIVE_CHAINS=$(jq -r '
        .cascade_coordination.active_cascade_chains[]? |
        select(.status != "completed" and .status != "merged_into") |
        .chain_id
    ' orchestrator-state-v3.json)
    
    if [[ -n "$ACTIVE_CHAINS" ]]; then
        echo "❌ Active cascade chains remain: $ACTIVE_CHAINS"
        return 1
    fi
    
    # 2. Check no pending fixes
    PENDING_FIXES=$(jq -r '
        .cascade_coordination.pending_fixes | 
        to_entries[] | 
        select(.value.cascade_chain == null) |
        .key
    ' orchestrator-state-v3.json)
    
    if [[ -n "$PENDING_FIXES" ]]; then
        echo "❌ Pending fixes without cascade chains: $PENDING_FIXES"
        return 1
    fi
    
    # 3. Check all integrations are fresh
    STALE_INTEGRATE_WAVE_EFFORTSS=$(jq -r '
        .stale_integration_tracking.stale_integrations[]? |
        select(.recreation_completed != true) |
        .integration_id
    ' orchestrator-state-v3.json)
    
    if [[ -n "$STALE_INTEGRATE_WAVE_EFFORTSS" ]]; then
        echo "❌ Stale integrations remain: $STALE_INTEGRATE_WAVE_EFFORTSS"
        return 1
    fi
    
    # 4. Check project integration is fresh
    PROJECT_STALE=$(jq -r '.project_integration.is_stale // false' orchestrator-state-v3.json)
    if [[ "$PROJECT_STALE" == "true" ]]; then
        echo "❌ Project integration is still stale"
        return 1
    fi
    
    # 5. Final check for any new fixes
    check_for_new_fixes
    NEW_CHAINS=$(jq '.cascade_coordination.active_cascade_chains | 
                     map(select(.status == "pending")) | length' orchestrator-state-v3.json)
    if [[ "$NEW_CHAINS" -gt 0 ]]; then
        echo "❌ New cascade chains were just created"
        return 1
    fi
    
    echo "✅ All cascade operations complete - can exit CASCADE_REINTEGRATION"
    return 0
}
```

### 6. STATE TRANSITION ENFORCEMENT

**ALL integration states MUST check cascade_mode:**

```bash
# In every integration state (INTEGRATE_WAVE_EFFORTS, INTEGRATE_PHASE_WAVES, PROJECT_INTEGRATE_WAVE_EFFORTS, etc.)
check_cascade_mode_return() {
    CASCADE_MODE=$(jq -r '.cascade_coordination.cascade_mode // false' orchestrator-state-v3.json)
    
    if [[ "$CASCADE_MODE" == "true" ]]; then
        echo "🔄 CASCADE MODE ACTIVE - Returning to CASCADE_REINTEGRATION"
        jq '.state_machine.current_state = "CASCADE_REINTEGRATION"' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
        return 0
    fi
    
    # Normal state transition if not in cascade mode
    return 1
}
```

## Examples

### Example 1: Simple Overlapping Cascade
```
10:00 - Fix applied to phase1/wave1/effort1
        CASCADE_REINTEGRATION starts cascade_001
        
10:05 - Rebasing phase1/wave1/effort2
        
10:10 - Fix applied to phase1/wave2/effort1
        CASCADE_REINTEGRATION creates cascade_002
        
10:15 - Both cascades proceed in parallel
        cascade_001: recreating wave1 integration
        cascade_002: rebasing wave2/effort2
        
10:30 - Cascades converge at phase1 integration
        Chains merged, all fixes included
        
10:45 - Project integration recreated with ALL fixes
        CASCADE_REINTEGRATION exits
```

### Example 2: Continuous Fix Arrival
```
14:00 - Fix in phase1/wave1/effort1 → cascade_001 starts
14:05 - Fix in phase1/wave1/effort2 → added to cascade_001
14:10 - Fix in phase2/wave1/effort1 → cascade_002 starts
14:15 - Fix in phase1/wave2/effort1 → cascade_003 starts
14:20 - All cascades tracked independently
14:45 - Project integration contains ALL fixes from ALL cascades
```

### Example 3: Recovery After Crash
```
09:00 - Multiple cascade chains active
09:15 - System crashes
09:20 - Recovery: CASCADE_REINTEGRATION resumes
        - Loads active cascade chains from state
        - Identifies incomplete operations
        - Resumes from last completed operation
        - Checks for new fixes during downtime
        - Continues until all cascades complete
```

---

## Production Example: Test 5 Persistent Coordinator

### Test 5 Demonstrates R352 Persistent Coordinator Pattern

**What Test 5 Validates:**
- CASCADE_REINTEGRATION returns after EVERY integration operation
- Persistent control maintained throughout cascade
- Multiple wave-level reintegrations coordinated correctly
- Only exits when ALL cascades complete

**Actual State Transition Pattern:**

```
Timestamp          | From State              | To State                 | Cascade Event
-------------------|-------------------------|--------------------------|-------------------
2025-10-15 20:45   | FIX_WAVE_UPSTREAM_BUGS  | CASCADE_REINTEGRATION   | ENTRY #1 (cascade start)
2025-10-19 22:36   | CASCADE_REINTEGRATION   | INTEGRATE_WAVE_EFFORTS  | EXIT #1 (process wave_1_1)
2025-10-19 23:30   | MONITORING_INTEGRATE_WAVE_EFFORTS  | CASCADE_REINTEGRATION   | RE-ENTRY #1 (persistence!)
2025-10-19 23:31   | CASCADE_REINTEGRATION   | INTEGRATE_WAVE_EFFORTS  | EXIT #2 (process wave_1_2)
2025-10-19 23:18   | INTEGRATE_WAVE_EFFORTS  | CASCADE_REINTEGRATION   | RE-ENTRY #2 (persistence!)
2025-10-19 23:30   | CASCADE_REINTEGRATION   | INTEGRATE_PHASE_WAVES   | EXIT #3 (phase level)
```

**Persistent Coordinator Validation:**

```bash
# Count CASCADE_REINTEGRATION appearances in Test 5
CASCADE_ENTRIES=$(jq '[.state_machine.state_history[] |
                       select(.to_state == "CASCADE_REINTEGRATION")] |
                       length' orchestrator-state-v3.json)

echo "CASCADE_REINTEGRATION entries: $CASCADE_ENTRIES"
# Output: 3 (initial + 2 returns)

CASCADE_EXITS=$(jq '[.state_machine.state_history[] |
                     select(.from_state == "CASCADE_REINTEGRATION")] |
                     length' orchestrator-state-v3.json)

echo "CASCADE_REINTEGRATION exits: $CASCADE_EXITS"
# Output: 3 (to wave_1_1, wave_1_2, phase_1)

# Total transitions: 6 (proves persistent coordinator)
```

**Why This Matters:**

1. **Visibility:** Each CASCADE_REINTEGRATION return is explicit in state history
2. **Auditability:** Can see exactly when coordinator resumed control
3. **Recovery:** If crash occurs during wave_1_2, can resume from last CASCADE_REINTEGRATION
4. **Proof:** State machine shows actual execution flow, not hidden internal logic

**Contrast with Meta-State Pattern:**

```
Meta-State Pattern (NOT USED):
CASCADE_REINTEGRATION
  ├─ (internal) Process wave_1_1
  ├─ (internal) Process wave_1_2
  └─ (internal) Process phase_1
EXIT CASCADE_REINTEGRATION

Problem: All coordination is hidden, only 1 state transition visible

Persistent Coordinator (R352 - USED):
CASCADE_REINTEGRATION → INTEGRATE_WAVE_EFFORTS → CASCADE_REINTEGRATION
  → INTEGRATE_WAVE_EFFORTS → CASCADE_REINTEGRATION → INTEGRATE_PHASE_WAVES

Benefit: Every operation is visible in state history
```

**Test Location:** `tests/runtime-test-05-cross-container-cascade.sh`

**Full Analysis:** `TEST-5-CASCADE-REINTEGRATE_WAVE_EFFORTS-VALIDATION-REPORT.md`

---

## Validation

### Pre-Exit Validation
```bash
validate_cascade_completion() {
    # Must pass ALL checks
    can_exit_cascade_reintegration || {
        echo "❌ CASCADE_REINTEGRATION cannot exit - conditions not met"
        return 1
    }
    
    # Verify project integration
    cd /efforts/project-integration
    git log --oneline -10
    
    # Count total fixes
    TOTAL_FIXES=$(jq '
        [.cascade_coordination.active_cascade_chains[].trigger.fix_ids[]] | 
        unique | length
    ' orchestrator-state-v3.json)
    
    echo "✅ CASCADE COMPLETE: $TOTAL_FIXES fixes successfully cascaded to project"
}
```

## Related Rules
- R327: Mandatory Re-Integration After Fixes (base cascade requirement)
- R348: Cascade State Transitions (state machine enforcement)
- R350: Complete Cascade Dependency Graph (dependency tracking)
- R351: Cascade Execution Protocol (execution order)

## Grading Impact
- Missing any fix in project integration: -100% AUTOMATIC FAILURE
- Not returning to CASCADE_REINTEGRATION: -50% per violation
- Losing cascade chain tracking: -30% per chain
- Exiting CASCADE_REINTEGRATION prematurely: -100% AUTOMATIC FAILURE