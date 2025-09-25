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
7. BUILD SUCCESS with ALL fixes
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
    ' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
}
```

### 4. NEW FIX DETECTION DURING CASCADE

**CASCADE_REINTEGRATION must check for new fixes after EVERY operation:**

```bash
check_for_new_fixes() {
    echo "🔍 Checking for new fixes that arrived during cascade operation..."
    
    # Get timestamp of last cascade operation
    LAST_OP_TIME=$(jq -r '.cascade_coordination.last_operation_completed_at' orchestrator-state.json)
    
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
    ' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
    
    # Check for new fixes
    check_for_new_fixes
    
    # Return to CASCADE_REINTEGRATION for next operation
    echo "🔄 Returning control to CASCADE_REINTEGRATION"
    jq '.current_state = "CASCADE_REINTEGRATION"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
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
    ' orchestrator-state.json)
    
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
    ' orchestrator-state.json)
    
    if [[ -n "$PENDING_FIXES" ]]; then
        echo "❌ Pending fixes without cascade chains: $PENDING_FIXES"
        return 1
    fi
    
    # 3. Check all integrations are fresh
    STALE_INTEGRATIONS=$(jq -r '
        .stale_integration_tracking.stale_integrations[]? |
        select(.recreation_completed != true) |
        .integration_id
    ' orchestrator-state.json)
    
    if [[ -n "$STALE_INTEGRATIONS" ]]; then
        echo "❌ Stale integrations remain: $STALE_INTEGRATIONS"
        return 1
    fi
    
    # 4. Check project integration is fresh
    PROJECT_STALE=$(jq -r '.project_integration.is_stale // false' orchestrator-state.json)
    if [[ "$PROJECT_STALE" == "true" ]]; then
        echo "❌ Project integration is still stale"
        return 1
    fi
    
    # 5. Final check for any new fixes
    check_for_new_fixes
    NEW_CHAINS=$(jq '.cascade_coordination.active_cascade_chains | 
                     map(select(.status == "pending")) | length' orchestrator-state.json)
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
# In every integration state (INTEGRATION, PHASE_INTEGRATION, PROJECT_INTEGRATION, etc.)
check_cascade_mode_return() {
    CASCADE_MODE=$(jq -r '.cascade_coordination.cascade_mode // false' orchestrator-state.json)
    
    if [[ "$CASCADE_MODE" == "true" ]]; then
        echo "🔄 CASCADE MODE ACTIVE - Returning to CASCADE_REINTEGRATION"
        jq '.current_state = "CASCADE_REINTEGRATION"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
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
    ' orchestrator-state.json)
    
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