# 🔴🔴🔴 SUPREME RULE R350: Complete Cascade Dependency Graph

## Criticality: SUPREME LAW
**Violation = -100% AUTOMATIC FAILURE**

## Description
This rule mandates the tracking and management of ALL dependency relationships in the Software Factory, including effort-to-effort dependencies within waves, effort dependencies on previous wave integrations, and the complete cascade chain that must be executed when any fix is applied at any level.

## 🔴🔴🔴 THE DEPENDENCY GRAPH LAW 🔴🔴🔴

**EVERY DEPENDENCY MUST BE TRACKED! EVERY CASCADE MUST BE COMPLETE!**

### The Problem This Solves
```
❌ BROKEN FLOW (What happens without dependency tracking):
1. effort1 gets a fix
2. effort2 (based on effort1) is NOT rebased
3. wave integration gets stale effort2
4. Next wave efforts based on broken integration
5. CASCADE FAILURE - entire project corrupted

✅ CORRECT FLOW (With dependency graph):
1. effort1 gets a fix
2. Dependency graph shows effort2 depends on effort1
3. effort2 automatically rebased on fixed effort1
4. All dependent efforts cascade-rebased
5. Wave re-integrated with all updated efforts
6. Next wave efforts rebase on fresh integration
7. CASCADE PROJECT_DONE - all dependencies satisfied
```

## Core Requirements

### 1. COMPLETE DEPENDENCY TRACKING STRUCTURE

The orchestrator-state-v3.json MUST maintain a complete dependency graph:

```json
"dependency_graph": {
  "effort_dependencies": [
    {
      "effort": "phase1-wave1-effort2",
      "depends_on": "phase1-wave1-effort1",
      "dependency_type": "sequential",
      "requires_rebase_on_change": true
    },
    {
      "effort": "phase1-wave2-effort1",
      "depends_on": "phase1-wave1-integration",
      "dependency_type": "wave_base",
      "requires_rebase_on_change": true
    }
  ],
  "integration_dependencies": [
    {
      "integration": "phase1-wave1-integration",
      "depends_on": ["phase1-wave1-effort1", "phase1-wave1-effort2"],
      "dependency_type": "source_efforts",
      "requires_recreation_on_change": true
    },
    {
      "integration": "phase1-integration",
      "depends_on": ["phase1-wave1-integration", "phase1-wave2-integration"],
      "dependency_type": "wave_integrations",
      "requires_recreation_on_change": true
    }
  ],
  "cascade_chains": []
}
```

### 2. DEPENDENCY TYPES AND RULES

#### Effort-to-Effort Dependencies (Within Wave)
```bash
# Sequential dependencies within a wave
# effort2 depends on effort1, effort3 depends on effort2, etc.
track_effort_dependency() {
    local EFFORT=$1
    local DEPENDS_ON=$2
    local TYPE=$3  # sequential|parallel|independent
    
    jq --arg effort "$EFFORT" \
       --arg depends "$DEPENDS_ON" \
       --arg type "$TYPE" \
       '.dependency_graph.effort_dependencies += [{
          "effort": $effort,
          "depends_on": $depends,
          "dependency_type": $type,
          "requires_rebase_on_change": true,
          "tracked_at": now | todate
       }]' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
}
```

#### Effort-to-Integration Dependencies (Cross-Wave)
```bash
# Efforts in new wave depend on previous wave's integration
track_wave_base_dependency() {
    local PHASE=$1
    local WAVE=$2
    
    if [[ $WAVE -gt 1 ]]; then
        local PREV_WAVE=$((WAVE - 1))
        local BASE="phase${PHASE}-wave${PREV_WAVE}-integration"
        
        # ALL efforts in this wave depend on previous wave integration
        for effort in $(list_wave_efforts $PHASE $WAVE); do
            track_effort_dependency "$effort" "$BASE" "wave_base"
        done
    fi
}
```

#### Integration Dependencies
```bash
# Integrations depend on their source branches
track_integration_dependencies() {
    local INTEGRATE_WAVE_EFFORTS=$1
    shift
    local SOURCES=("$@")
    
    local DEPS_JSON=$(printf '%s\n' "${SOURCES[@]}" | jq -R . | jq -s .)
    
    jq --arg integration "$INTEGRATE_WAVE_EFFORTS" \
       --argjson deps "$DEPS_JSON" \
       '.dependency_graph.integration_dependencies += [{
          "integration": $integration,
          "depends_on": $deps,
          "dependency_type": "source_efforts",
          "requires_recreation_on_change": true
       }]' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
}
```

### 3. CASCADE CHAIN CALCULATION

When ANY fix is applied, calculate the COMPLETE cascade chain:

```bash
calculate_cascade_chain() {
    local FIXED_BRANCH=$1
    local CASCADE_ID=$(uuidgen)
    
    echo "🔍 R350: Calculating complete cascade chain for fix in $FIXED_BRANCH"
    
    # Initialize cascade chain
    local CHAIN='[]'
    
    # Step 1: Find all efforts that depend on this branch
    DEPENDENT_EFFORTS=$(jq -r --arg branch "$FIXED_BRANCH" '
        .dependency_graph.effort_dependencies[] |
        select(.depends_on == $branch) |
        .effort
    ' orchestrator-state-v3.json)
    
    # Add effort rebases to chain
    for effort in $DEPENDENT_EFFORTS; do
        CHAIN=$(echo "$CHAIN" | jq --arg effort "$effort" --arg on "$FIXED_BRANCH" '. += [{
            "type": "rebase",
            "target": $effort,
            "on": $on,
            "status": "pending"
        }]')
        
        # Recursively find efforts depending on this effort
        calculate_dependent_cascades "$effort" "$CHAIN"
    done
    
    # Step 2: Find integrations that need recreation
    AFFECTED_INTEGRATE_WAVE_EFFORTSS=$(jq -r --arg branch "$FIXED_BRANCH" '
        .dependency_graph.integration_dependencies[] |
        select(.depends_on | contains([$branch])) |
        .integration
    ' orchestrator-state-v3.json)
    
    # Add integration recreations to chain
    for integration in $AFFECTED_INTEGRATE_WAVE_EFFORTSS; do
        CHAIN=$(echo "$CHAIN" | jq --arg int "$integration" '. += [{
            "type": "recreate",
            "target": $int,
            "status": "pending"
        }]')
        
        # Find higher-level integrations affected
        calculate_integration_cascades "$integration" "$CHAIN"
    done
    
    # Step 3: Order the chain by dependency order
    ORDERED_CHAIN=$(order_cascade_operations "$CHAIN")
    
    # Save cascade chain
    jq --arg id "$CASCADE_ID" \
       --arg trigger "$FIXED_BRANCH" \
       --argjson chain "$ORDERED_CHAIN" \
       '.dependency_graph.cascade_chains += [{
          "cascade_id": $id,
          "trigger": $trigger,
          "triggered_at": now | todate,
          "operations": $chain,
          "status": "pending"
       }]' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
    
    echo "✅ R350: Cascade chain calculated with $(echo "$ORDERED_CHAIN" | jq length) operations"
    return 0
}
```

### 4. CASCADE EXECUTION ORDER

The cascade MUST be executed in DEPENDENCY ORDER:

```bash
order_cascade_operations() {
    local OPERATIONS=$1
    
    # Order rules:
    # 1. Rebases before recreations
    # 2. Lower-level before higher-level
    # 3. Dependencies before dependents
    
    echo "$OPERATIONS" | jq '
        sort_by(
            # First sort key: operation type (rebase=0, recreate=1)
            (if .type == "rebase" then 0 else 1 end),
            # Second sort key: level (effort=0, wave=1, phase=2, project=3)
            (if .target | contains("effort") then 0
             elif .target | contains("wave") then 1
             elif .target | contains("phase") then 2
             else 3 end),
            # Third sort key: wave number
            (.target | capture("wave(?<w>[0-9]+)").w // "0" | tonumber),
            # Fourth sort key: effort number
            (.target | capture("effort(?<e>[0-9]+)").e // "0" | tonumber)
        )
    '
}
```

### 5. DEPENDENCY VALIDATION GATES

Before ANY operation, validate dependencies are satisfied:

```bash
validate_dependencies_satisfied() {
    local BRANCH=$1
    
    echo "🔍 R350: Validating dependencies for $BRANCH"
    
    # Get dependencies for this branch
    DEPENDENCIES=$(jq -r --arg branch "$BRANCH" '
        .dependency_graph.effort_dependencies[] |
        select(.effort == $branch) |
        .depends_on
    ' orchestrator-state-v3.json)
    
    for dep in $DEPENDENCIES; do
        # Check if dependency is up-to-date
        DEP_TIME=$(git log -1 --format=%ct "$dep" 2>/dev/null || echo 0)
        BRANCH_BASE_TIME=$(git merge-base "$BRANCH" "$dep" | xargs git log -1 --format=%ct)
        
        if [[ $DEP_TIME -gt $BRANCH_BASE_TIME ]]; then
            echo "❌ R350 VIOLATION: $BRANCH not rebased on latest $dep"
            echo "   Dependency updated at: $(date -d @$DEP_TIME)"
            echo "   Branch based on: $(date -d @$BRANCH_BASE_TIME)"
            return 1
        fi
    done
    
    echo "✅ R350: All dependencies satisfied for $BRANCH"
    return 0
}
```

### 6. AUTOMATIC CASCADE TRIGGERING

Integrate with fix detection to automatically trigger cascades:

```bash
# Hook into fix application
on_fix_applied() {
    local FIXED_BRANCH=$1
    local FIX_COMMIT=$2
    
    echo "🔴 R350: Fix detected in $FIXED_BRANCH, triggering cascade analysis"
    
    # Calculate cascade chain
    calculate_cascade_chain "$FIXED_BRANCH"
    
    # Transition to CASCADE_REINTEGRATION
    jq '.state_machine.current_state = "CASCADE_REINTEGRATION" |
        .transition_reason = "R350: Fix applied, cascade required"' \
        orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
    
    echo "🔴 R350: Transitioned to CASCADE_REINTEGRATION for cascade execution"
}
```

## Common Violations (ALL RESULT IN AUTOMATIC FAILURE)

### ❌ VIOLATION 1: Missing Dependency Tracking
```bash
# WRONG:
# Creating effort2 without tracking dependency on effort1
git checkout -b phase1-wave1-effort2
```

### ✅ CORRECTION 1: Track All Dependencies
```bash
# RIGHT:
git checkout -b phase1-wave1-effort2
track_effort_dependency "phase1-wave1-effort2" "phase1-wave1-effort1" "sequential"
```

### ❌ VIOLATION 2: Incomplete Cascade Chain
```bash
# WRONG:
# Only rebasing immediate dependents
rebase_effort "effort2" "effort1"
# Missing effort3 which depends on effort2!
```

### ✅ CORRECTION 2: Complete Cascade
```bash
# RIGHT:
calculate_cascade_chain "effort1"
execute_cascade_chain  # Handles ALL dependencies
```

### ❌ VIOLATION 3: Wrong Cascade Order
```bash
# WRONG:
# Recreating integration before rebasing efforts
recreate_wave_integration
rebase_dependent_efforts  # Too late!
```

### ✅ CORRECTION 3: Dependency Order
```bash
# RIGHT:
# Follow calculated cascade order
execute_cascade_operations_in_order
```

## Grading Impact

### AUTOMATIC FAILURE (-100%)
- Not tracking effort-to-effort dependencies
- Missing dependencies in cascade chain
- Executing cascades out of order
- Leaving dependencies unsatisfied after cascade
- Not calculating cascade chains for fixes

### MAJOR VIOLATIONS (-50%)
- Incomplete dependency graph
- Manual cascade execution instead of automated
- Not validating dependencies before operations

### COMPLIANCE BONUS (+30%)
- Complete dependency graph maintained
- Automatic cascade calculation
- Perfect cascade execution order
- All dependencies validated and satisfied

## Relationship to Other Rules

### Depends on
- R327: Mandatory Re-Integration After Fixes (defines when cascades needed)
- R348: Cascade State Transitions (provides CASCADE_REINTEGRATION state)
- R308: Incremental Branching Strategy (defines base dependencies)

### Enables
- R351: Cascade Execution Protocol (uses dependency graph)
- R328: Integration Freshness Validation (validates using dependencies)

### Works with
- R346: State Metadata Synchronization
- R337: Base Branch Single Source Truth

## Quick Reference

### Track New Effort Dependency
```bash
track_effort_dependency "phase1-wave2-effort3" "phase1-wave2-effort2" "sequential"
```

### Calculate Cascade for Fix
```bash
calculate_cascade_chain "phase1-wave1-effort1"
```

### Validate Dependencies
```bash
validate_dependencies_satisfied "phase1-wave2-effort2"
```

### Execute Cascade
```bash
execute_cascade_in_dependency_order
```

## Remember

**"EVERY DEPENDENCY TRACKED, EVERY CASCADE COMPLETE"**
**"No hidden dependencies, no partial cascades"**
**"The graph is truth, the cascade is law"**
**"Dependency order is SACRED"**

### 🔴🔴🔴 THE DEPENDENCY MANTRA 🔴🔴🔴
```
Track each link in the chain,
From effort to effort, wave to main.
When fixes flow through the tree,
CASCADE COMPLETELY, dependency-free!

The graph reveals what must be done,
Each operation, one by one.
In perfect order, no shortcuts taken,
Until all branches are rewaken!
```

The goal: COMPLETE DEPENDENCY AWARENESS - Every relationship tracked, every cascade calculated, every operation ordered, no dependency left behind.