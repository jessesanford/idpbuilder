# 🔴🔴🔴 SUPREME RULE R351: Cascade Execution Protocol

## Criticality: SUPREME LAW
**Violation = -100% AUTOMATIC FAILURE**

## Description
This rule defines the EXACT protocol for executing cascade operations when fixes are applied. It distinguishes between rebasing (for efforts) and re-integration (for integrations), ensures proper execution order, handles both sequential and parallel operations where safe, and guarantees no dependency is missed.

## 🔴🔴🔴 THE CASCADE EXECUTION LAW 🔴🔴🔴

**CASCADES MUST BE EXECUTED IN EXACT DEPENDENCY ORDER WITH ZERO TOLERANCE FOR SHORTCUTS!**

### The Problem This Solves
```
❌ BROKEN EXECUTION (What fails without protocol):
1. Random execution order
2. Rebasing after integration (backwards!)
3. Parallel execution of dependent operations
4. Missing operations in the chain
5. PROJECT CORRUPTION

✅ CORRECT EXECUTION (With protocol):
1. Dependency-ordered execution
2. Rebases complete before integrations
3. Parallel only where safe
4. Every operation tracked and verified
5. PROJECT INTEGRITY MAINTAINED
```

## Core Requirements

### 1. CASCADE EXECUTION STATE MACHINE

The CASCADE_REINTEGRATION state MUST follow this execution flow:

```yaml
CASCADE_REINTEGRATION:
  entry_actions:
    - Load pending cascade chains
    - Validate dependency graph
    - Create execution plan
  
  execution_loop:
    1. SELECT_NEXT_OPERATION
    2. VALIDATE_PREREQUISITES
    3. EXECUTE_OPERATION
    4. VERIFY_PROJECT_DONE
    5. UPDATE_CASCADE_STATUS
    6. LOOP or EXIT
  
  exit_conditions:
    - All operations completed
    - All dependencies satisfied
    - All integrations fresh
```

### 2. OPERATION TYPES AND EXECUTION

#### REBASE Operations (for efforts)
```bash
execute_rebase_operation() {
    local TARGET=$1
    local BASE=$2
    local CASCADE_ID=$3
    local OP_INDEX=$4
    
    echo "🔄 R351: Executing REBASE: $TARGET onto $BASE"
    
    # Pre-execution validation
    if ! validate_rebase_prerequisites "$TARGET" "$BASE"; then
        echo "❌ R351: Prerequisites not met for rebase"
        return 1
    fi
    
    # Execute rebase
    cd "/efforts/$(dirname $TARGET)"
    git checkout "$TARGET"
    
    # Store current head for rollback
    ORIGINAL_HEAD=$(git rev-parse HEAD)
    
    # Attempt rebase
    if git rebase "$BASE"; then
        echo "✅ Rebase successful"
        
        # Push with force-with-lease for safety
        if git push --force-with-lease origin "$TARGET"; then
            update_cascade_operation_status "$CASCADE_ID" "$OP_INDEX" "completed"
            return 0
        else
            echo "❌ Push failed, rolling back"
            git reset --hard "$ORIGINAL_HEAD"
            update_cascade_operation_status "$CASCADE_ID" "$OP_INDEX" "failed"
            return 1
        fi
    else
        echo "❌ Rebase failed, aborting"
        git rebase --abort
        update_cascade_operation_status "$CASCADE_ID" "$OP_INDEX" "failed"
        return 1
    fi
}
```

#### RECREATE Operations (for integrations)
```bash
execute_recreate_operation() {
    local INTEGRATE_WAVE_EFFORTS=$1
    local CASCADE_ID=$2
    local OP_INDEX=$3
    
    echo "🔄 R351: Executing RECREATE: $INTEGRATE_WAVE_EFFORTS"
    
    # Determine integration level
    if [[ "$INTEGRATE_WAVE_EFFORTS" =~ wave.*integration ]]; then
        LEVEL="wave"
    elif [[ "$INTEGRATE_WAVE_EFFORTS" =~ phase.*integration ]]; then
        LEVEL="phase"
    elif [[ "$INTEGRATE_WAVE_EFFORTS" =~ project.*integration ]]; then
        LEVEL="project"
    fi
    
    # Delete old integration
    echo "🗑️ Deleting stale $INTEGRATE_WAVE_EFFORTS"
    git push origin --delete "$INTEGRATE_WAVE_EFFORTS" 2>/dev/null || true
    rm -rf "/efforts/*/$(basename $INTEGRATE_WAVE_EFFORTS)"
    
    # Mark integration as stale in state
    mark_integration_stale "$INTEGRATE_WAVE_EFFORTS"
    
    # Trigger re-integration through state machine
    case $LEVEL in
        wave)
            transition_to_state "INTEGRATE_WAVE_EFFORTS"
            ;;
        phase)
            transition_to_state "INTEGRATE_PHASE_WAVES"  # SF 3.0
            ;;
        project)
            transition_to_state "PROJECT_INTEGRATE_WAVE_EFFORTS"
            ;;
    esac
    
    # Wait for integration completion
    wait_for_integration_completion "$INTEGRATE_WAVE_EFFORTS" || {
        update_cascade_operation_status "$CASCADE_ID" "$OP_INDEX" "failed"
        return 1
    }
    
    update_cascade_operation_status "$CASCADE_ID" "$OP_INDEX" "completed"
    return 0
}
```

### 3. CASCADE EXECUTION PLAN

Before execution, create a detailed plan:

```bash
create_cascade_execution_plan() {
    local CASCADE_ID=$1
    
    echo "📋 R351: Creating cascade execution plan"
    
    # Get cascade chain
    CASCADE_CHAIN=$(jq -r --arg id "$CASCADE_ID" '
        .dependency_graph.cascade_chains[] |
        select(.cascade_id == $id) |
        .operations
    ' orchestrator-state-v3.json)
    
    # Group operations by parallelizability
    EXECUTION_PLAN='{"sequential_groups": []}'
    CURRENT_GROUP='{"parallel_operations": []}'
    PREVIOUS_TARGETS=""
    
    echo "$CASCADE_CHAIN" | jq -c '.[]' | while read -r op; do
        TARGET=$(echo "$op" | jq -r '.target')
        TYPE=$(echo "$op" | jq -r '.type')
        
        # Check if this can be parallel with current group
        CAN_PARALLEL=$(check_parallelizability "$TARGET" "$PREVIOUS_TARGETS")
        
        if [[ "$CAN_PARALLEL" == "true" ]]; then
            # Add to current parallel group
            CURRENT_GROUP=$(echo "$CURRENT_GROUP" | jq --argjson op "$op" \
                '.parallel_operations += [$op]')
        else
            # Start new group
            if [[ $(echo "$CURRENT_GROUP" | jq '.parallel_operations | length') -gt 0 ]]; then
                EXECUTION_PLAN=$(echo "$EXECUTION_PLAN" | jq --argjson group "$CURRENT_GROUP" \
                    '.sequential_groups += [$group]')
            fi
            CURRENT_GROUP=$(echo '{"parallel_operations": []}' | jq --argjson op "$op" \
                '.parallel_operations += [$op]')
        fi
        
        PREVIOUS_TARGETS="$PREVIOUS_TARGETS $TARGET"
    done
    
    # Add final group
    if [[ $(echo "$CURRENT_GROUP" | jq '.parallel_operations | length') -gt 0 ]]; then
        EXECUTION_PLAN=$(echo "$EXECUTION_PLAN" | jq --argjson group "$CURRENT_GROUP" \
            '.sequential_groups += [$group]')
    fi
    
    # Save execution plan
    jq --arg id "$CASCADE_ID" \
       --argjson plan "$EXECUTION_PLAN" \
       '.dependency_graph.cascade_chains |= map(
           if .cascade_id == $id then
               . + {"execution_plan": $plan}
           else . end
       )' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
    
    echo "✅ Execution plan created with $(echo "$EXECUTION_PLAN" | jq '.sequential_groups | length') sequential groups"
}
```

### 4. OVERLAPPING CASCADE CHAIN MANAGEMENT (R352)

Manage multiple cascade chains running simultaneously:

```bash
manage_overlapping_cascades() {
    echo "🔄 R352: Managing overlapping cascade chains"
    
    # Get all active chains
    ACTIVE_CHAINS=$(jq -r '.cascade_coordination.active_cascade_chains[] | 
                           select(.status == "in_progress" or .status == "pending") | 
                           .chain_id' orchestrator-state-v3.json)
    
    echo "📋 Active cascade chains: $(echo $ACTIVE_CHAINS | wc -w)"
    
    # Process each chain
    for chain_id in $ACTIVE_CHAINS; do
        # Get next operation for this chain
        NEXT_OP=$(get_next_operation_for_chain "$chain_id")
        
        if [[ -n "$NEXT_OP" ]]; then
            # Check if operation conflicts with other chains
            if ! check_operation_conflict "$NEXT_OP" "$chain_id"; then
                execute_cascade_operation "$NEXT_OP" "$chain_id"
                
                # After operation, check for convergence
                check_for_chain_convergence "$NEXT_OP"
            else
                echo "⏸️ Deferring operation due to conflict"
            fi
        fi
    done
    
    # Check for new fixes (R352 requirement)
    check_for_new_fixes_and_create_chains
    
    # Return to CASCADE_REINTEGRATION
    return_to_cascade_reintegration
}

# Merge converging cascade chains
merge_cascade_chains() {
    local PRIMARY=$1
    local SECONDARY=$2
    local MERGE_POINT=$3
    
    echo "🔀 R352: Merging cascade chains at $MERGE_POINT"
    
    # Combine fix lists
    jq --arg p "$PRIMARY" --arg s "$SECONDARY" --arg mp "$MERGE_POINT" '
        (.cascade_coordination.active_cascade_chains[] | select(.chain_id == $p)) |= 
        . + {
            merged_with: ((.merged_with // []) + [$s]),
            merge_point: $mp,
            merged_at: now | todate
        } |
        (.cascade_coordination.active_cascade_chains[] | select(.chain_id == $s)) |= 
        . + {status: "merged_into", merged_into: $p}' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
}
```

### 5. PARALLEL VS SEQUENTIAL EXECUTION

Determine when operations can be parallel:

```bash
check_parallelizability() {
    local TARGET=$1
    local PREVIOUS_TARGETS=$2
    
    # Rules for parallelization:
    # 1. Different waves can be parallel if no dependencies
    # 2. Independent efforts within same wave can be parallel
    # 3. Recreations must be sequential (wave → phase → project)
    
    # Check if target depends on any previous targets
    for prev in $PREVIOUS_TARGETS; do
        DEPENDS=$(jq -r --arg target "$TARGET" --arg prev "$prev" '
            .dependency_graph.effort_dependencies[] |
            select(.effort == $target and .depends_on == $prev) |
            .effort
        ' orchestrator-state-v3.json)
        
        if [[ -n "$DEPENDS" ]]; then
            echo "false"  # Has dependency, must be sequential
            return
        fi
    done
    
    echo "true"  # No dependencies, can be parallel
}
```

### 5. CASCADE EXECUTION ENGINE

The main execution loop:

```bash
execute_cascade_chain() {
    local CASCADE_ID=$1
    
    echo "🚀 R351: Starting cascade execution for $CASCADE_ID"
    
    # Create execution plan
    create_cascade_execution_plan "$CASCADE_ID"
    
    # Get execution plan
    PLAN=$(jq -r --arg id "$CASCADE_ID" '
        .dependency_graph.cascade_chains[] |
        select(.cascade_id == $id) |
        .execution_plan
    ' orchestrator-state-v3.json)
    
    # Execute each sequential group
    GROUP_COUNT=$(echo "$PLAN" | jq '.sequential_groups | length')
    
    for ((g=0; g<$GROUP_COUNT; g++)); do
        echo "📦 Executing sequential group $((g+1))/$GROUP_COUNT"
        
        GROUP=$(echo "$PLAN" | jq ".sequential_groups[$g]")
        OPS=$(echo "$GROUP" | jq -c '.parallel_operations[]')
        
        # Execute operations in parallel within group
        PIDS=""
        while IFS= read -r op; do
            TYPE=$(echo "$op" | jq -r '.type')
            TARGET=$(echo "$op" | jq -r '.target')
            ON=$(echo "$op" | jq -r '.on // empty')
            
            # Execute in background for parallelization
            (
                case $TYPE in
                    rebase)
                        execute_rebase_operation "$TARGET" "$ON" "$CASCADE_ID" "$g"
                        ;;
                    recreate)
                        execute_recreate_operation "$TARGET" "$CASCADE_ID" "$g"
                        ;;
                esac
            ) &
            
            PIDS="$PIDS $!"
        done <<< "$OPS"
        
        # Wait for all operations in group to complete
        for pid in $PIDS; do
            wait $pid || {
                echo "❌ R351: Operation failed in group $((g+1))"
                mark_cascade_failed "$CASCADE_ID"
                return 1
            }
        done
        
        echo "✅ Sequential group $((g+1)) completed"
    done
    
    # Mark cascade as completed
    mark_cascade_completed "$CASCADE_ID"
    
    echo "🎉 R351: Cascade execution completed successfully!"
    return 0
}
```

### 6. VALIDATION AND RECOVERY

#### Pre-execution Validation
```bash
validate_cascade_prerequisites() {
    local CASCADE_ID=$1
    
    echo "🔍 R351: Validating cascade prerequisites"
    
    # Check 1: All source branches exist
    # Check 2: No conflicting operations in progress
    # Check 3: Dependency graph is complete
    # Check 4: State machine is in CASCADE_REINTEGRATION
    
    [[ "$(jq -r '.state_machine.current_state' orchestrator-state-v3.json)" == "CASCADE_REINTEGRATION" ]] || {
        echo "❌ Not in CASCADE_REINTEGRATION state"
        return 1
    }
    
    # More validation...
    return 0
}
```

#### Recovery from Failed Operations
```bash
recover_from_cascade_failure() {
    local CASCADE_ID=$1
    local FAILED_OP=$2
    
    echo "🔧 R351: Attempting cascade recovery"
    
    # Options:
    # 1. Retry the failed operation
    # 2. Skip if non-critical (NO! All cascade ops are critical!)
    # 3. Rollback entire cascade
    # 4. Manual intervention required
    
    # For now, mark for manual intervention
    jq --arg id "$CASCADE_ID" --arg op "$FAILED_OP" \
        '.cascade_recovery_required += [{
            "cascade_id": $id,
            "failed_operation": $op,
            "timestamp": now | todate,
            "action_required": "manual_intervention"
        }]' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
    
    echo "❌ R351: Manual intervention required for cascade $CASCADE_ID"
    transition_to_state "ERROR_RECOVERY"
}
```

### 7. CASCADE COMPLETION VERIFICATION

After execution, verify ALL objectives met:

```bash
verify_cascade_completion() {
    local CASCADE_ID=$1
    
    echo "🔍 R351: Verifying cascade completion"
    
    # Check 1: All operations marked completed
    PENDING=$(jq -r --arg id "$CASCADE_ID" '
        .dependency_graph.cascade_chains[] |
        select(.cascade_id == $id) |
        .operations[] |
        select(.status != "completed")
    ' orchestrator-state-v3.json)
    
    [[ -z "$PENDING" ]] || {
        echo "❌ Pending operations remain"
        return 1
    }
    
    # Check 2: All integrations fresh
    STALE=$(jq -r '
        .stale_integration_tracking.stale_integrations[] |
        select(.recreation_completed != true)
    ' orchestrator-state-v3.json)
    
    [[ -z "$STALE" ]] || {
        echo "❌ Stale integrations remain"
        return 1
    }
    
    # Check 3: All dependencies satisfied
    if ! validate_all_dependencies_satisfied; then
        echo "❌ Dependencies not satisfied"
        return 1
    fi
    
    echo "✅ R351: Cascade $CASCADE_ID fully completed!"
    return 0
}
```

## Common Violations (ALL RESULT IN AUTOMATIC FAILURE)

### ❌ VIOLATION 1: Out-of-Order Execution
```bash
# WRONG:
recreate_wave_integration
rebase_dependent_efforts  # Too late!
```

### ✅ CORRECTION 1: Dependency Order
```bash
# RIGHT:
execute_cascade_chain "$CASCADE_ID"  # Handles order automatically
```

### ❌ VIOLATION 2: Skipping Operations
```bash
# WRONG:
# Only executing some operations
for op in "effort1" "effort2"; do
    rebase_effort "$op"
done
# Missing effort3 and integrations!
```

### ✅ CORRECTION 2: Complete Execution
```bash
# RIGHT:
# Execute ALL operations in cascade
execute_cascade_chain "$CASCADE_ID"  # Executes EVERYTHING
```

### ❌ VIOLATION 3: Parallel Execution of Dependencies
```bash
# WRONG:
# Rebasing dependent efforts in parallel
rebase_effort "effort1" &
rebase_effort "effort2" &  # Depends on effort1!
wait
```

### ✅ CORRECTION 3: Respect Dependencies
```bash
# RIGHT:
# Sequential for dependencies, parallel where safe
execute_cascade_with_parallelization "$CASCADE_ID"
```

## Grading Impact

### AUTOMATIC FAILURE (-100%)
- Executing cascade operations out of order
- Skipping operations in cascade chain
- Parallel execution of dependent operations
- Not completing all cascade operations
- Exiting CASCADE_REINTEGRATION with incomplete cascade

### MAJOR VIOLATIONS (-50%)
- Manual cascade execution instead of automated
- Not creating execution plan
- Not verifying cascade completion

### COMPLIANCE BONUS (+30%)
- Perfect cascade execution order
- Optimal parallelization where safe
- Complete cascade verification
- Clean recovery from failures

## Relationship to Other Rules

### Depends on
- R350: Complete Cascade Dependency Graph (provides dependency data)
- R327: Mandatory Re-Integration After Fixes (defines cascade need)
- R348: Cascade State Transitions (provides CASCADE_REINTEGRATION)

### Enables
- R328: Integration Freshness Validation (cascades ensure freshness)
- R346: State Metadata Synchronization (cascade status tracking)

## Quick Reference

### Execute Complete Cascade
```bash
# From CASCADE_REINTEGRATION state
CASCADE_ID=$(get_pending_cascade_id)
execute_cascade_chain "$CASCADE_ID"
```

### Check Cascade Status
```bash
get_cascade_status "$CASCADE_ID"
```

### Verify Cascade Completion
```bash
verify_cascade_completion "$CASCADE_ID"
```

### Recover from Failure
```bash
recover_from_cascade_failure "$CASCADE_ID" "$FAILED_OP"
```

## Remember

**"EXECUTE IN ORDER, COMPLETE IN FULL"**
**"Dependencies first, integrations last"**
**"Parallel where safe, sequential where needed"**
**"No operation left behind"**

### 🔴🔴🔴 THE EXECUTION MANTRA 🔴🔴🔴
```
Rebase the efforts, one by one,
In order strict till all are done.
Then recreate integrations clean,
Wave, then phase, then project seen.

Parallel when paths don't cross,
Sequential when there'd be loss.
Track each step and verify,
Till all cascades satisfy!
```

The goal: PERFECT CASCADE EXECUTION - Every operation in order, every dependency respected, optimal parallelization, complete verification.