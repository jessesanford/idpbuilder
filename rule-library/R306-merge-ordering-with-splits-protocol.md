# 🔴🔴🔴 SUPREME RULE R306: Merge Ordering with Splits Protocol

## Rule Definition
When efforts have splits, ALL splits of an effort MUST be merged before ANY dependent efforts can be merged. Dependencies are at the EFFORT level, not the SPLIT level. This ensures complete functionality is available to dependent code.

## Criticality: 🔴🔴🔴 SUPREME
Violating merge order with splits causes cascade failures, broken dependencies, and corrupted integrations that are nearly impossible to recover from.

## Core Principle
**Dependencies work at EFFORT level, not SPLIT level!**

If Effort E2 depends on Effort E1:
- E2 depends on **ALL** of E1 (complete functionality)
- ALL E1 splits must merge before E2
- You cannot cherry-pick partial dependencies

## Requirements

### 1. Split Merge Sequencing
Splits created sequentially (per R302) MUST merge sequentially:
```bash
# CORRECT: Sequential merge order
git merge E1-split-001  # First split
git merge E1-split-002  # Based on split-001
git merge E1-split-003  # Based on split-002

# WRONG: Out of order
git merge E1-split-002  # ERROR: split-001 not merged!
```

### 2. Dependency Completion Before Dependent Merges
```yaml
# Example dependency chain
efforts:
  E1:
    splits: [split-001, split-002, split-003]
  E2:
    depends_on: [E1]  # Depends on COMPLETE E1

# CORRECT merge order:
1. E1-split-001
2. E1-split-002
3. E1-split-003  # E1 now complete
4. E2            # Can merge now

# WRONG merge order:
1. E1-split-001
2. E2            # ERROR: E1 incomplete!
3. E1-split-002
4. E1-split-003
```

### 3. Integration Agent Validation
The Integration Agent MUST validate before each merge:
```bash
validate_merge_readiness() {
    local branch="$1"
    local effort=$(echo "$branch" | sed 's/-split-[0-9]*//')
    
    # Get dependencies
    DEPS=$(yq ".efforts.\"$effort\".dependencies[]" orchestrator-state.yaml)
    
    for dep in $DEPS; do
        # Check split tracking
        SPLIT_COUNT=$(yq ".split_tracking.\"$dep\".split_count // 0" orchestrator-state.yaml)
        
        if [ "$SPLIT_COUNT" -gt 0 ]; then
            # Verify ALL splits merged
            for i in $(seq 1 $SPLIT_COUNT); do
                SPLIT_STATUS=$(yq ".split_tracking.\"$dep\".split_branches[$((i-1))].status" orchestrator-state.yaml)
                if [ "$SPLIT_STATUS" != "INTEGRATED" ]; then
                    echo "❌ BLOCKED: Cannot merge $branch"
                    echo "   Dependency $dep split $i not integrated"
                    return 1
                fi
            done
        fi
    done
    
    return 0
}
```

### 4. Code Reviewer Merge Plan Requirements
Merge plans MUST respect split ordering:
```markdown
## Branches to Merge (IN ORDER)

### Group 1: E1 Authentication (has splits)
1. phase1/wave1/E1-auth-split-001 (423 lines)
2. phase1/wave1/E1-auth-split-002 (389 lines)
3. phase1/wave1/E1-auth-split-003 (401 lines)

### Group 2: E2 User Profile (depends on E1)
4. phase1/wave1/E2-user-profile (567 lines)
   - Dependencies: [E1] ✓ (all splits above)

### Group 3: E3 Settings (independent)
5. phase1/wave1/E3-settings (234 lines)
   - Dependencies: none
```

### 5. Orchestrator Monitoring
During MONITORING_INTEGRATION, detect violations:
```bash
monitor_merge_order() {
    WORK_LOG="$INTEGRATION_DIR/work-log.md"
    
    # Track merged branches
    MERGED=$(grep "MERGED:" "$WORK_LOG" | awk '{print $2}')
    
    for branch in $MERGED; do
        # Check if dependencies were complete when merged
        if ! validate_dependencies_at_merge_time "$branch"; then
            echo "🔴 MERGE ORDER VIOLATION DETECTED!"
            echo "STOP_SIGNAL" > "$INTEGRATION_DIR/STOP"
            transition_to "ERROR_RECOVERY"
        fi
    done
}
```

## Complex Scenarios

### Scenario 1: Diamond Dependencies with Splits
```
     E1 (2 splits)
       /        \
    E2            E3 (3 splits)
       \        /
         E4

Correct Order:
1. E1-split-001
2. E1-split-002      # E1 complete
3. E2                # Needs complete E1
4. E3-split-001
5. E3-split-002
6. E3-split-003      # E3 complete
7. E4                # Needs E2 and complete E3
```

### Scenario 2: Chain Dependencies with Multiple Splits
```
E1 (3 splits) → E2 → E3 (2 splits) → E4

Correct Order:
1. E1-split-001
2. E1-split-002
3. E1-split-003      # E1 complete
4. E2                # Needs complete E1
5. E3-split-001
6. E3-split-002      # E3 complete
7. E4                # Needs complete E3
```

## Error Detection and Recovery

### Detection Points
1. **Code Reviewer**: When creating merge plan
2. **Integration Agent**: Before each merge
3. **Orchestrator**: During monitoring
4. **Architect**: During review

### Error Messages
```
❌ CRITICAL: Merge Order Violation
Cannot merge: phase1/wave1/E2-user-profile
Reason: Dependency E1 incomplete
Missing splits:
  - E1-auth-split-002 (not merged)
  - E1-auth-split-003 (not merged)
Action: Merge missing splits first
```

### Recovery Process
1. **STOP** current integration immediately
2. **IDENTIFY** missing splits
3. **RESET** to last valid state
4. **REORDER** merge plan
5. **RESTART** integration with correct order

## Validation Functions

### Pre-Integration Validation
```bash
validate_merge_plan_ordering() {
    local plan_file="$1"
    local errors=0
    
    # Extract merge sequence
    SEQUENCE=$(grep "^[0-9]\." "$plan_file" | sed 's/^[0-9]\. //')
    
    # Track what's been merged
    declare -A MERGED
    
    for branch in $SEQUENCE; do
        EFFORT=$(echo "$branch" | sed 's/-split-[0-9]*//')
        
        # Check dependencies
        DEPS=$(yq ".efforts.\"$EFFORT\".dependencies[]" orchestrator-state.yaml)
        
        for dep in $DEPS; do
            # Check if dependency is complete
            SPLIT_COUNT=$(yq ".split_tracking.\"$dep\".split_count // 0" orchestrator-state.yaml)
            
            if [ "$SPLIT_COUNT" -gt 0 ]; then
                # Check all splits are in MERGED
                for i in $(seq 1 $SPLIT_COUNT); do
                    SPLIT_KEY="${dep}-split-$(printf "%03d" $i)"
                    if [ "${MERGED[$SPLIT_KEY]}" != "1" ]; then
                        echo "❌ ERROR: $branch scheduled before $SPLIT_KEY"
                        ((errors++))
                    fi
                done
            else
                if [ "${MERGED[$dep]}" != "1" ]; then
                    echo "❌ ERROR: $branch scheduled before $dep"
                    ((errors++))
                fi
            fi
        done
        
        # Mark this branch as merged
        MERGED["$branch"]=1
    done
    
    return $errors
}
```

## Related Rules
- R302: Comprehensive Split Tracking Protocol
- R204: Orchestrator Split Infrastructure  
- R262: Merge Operation Protocols
- R269: Code Reviewer Merge Plan No Execution
- R297: Architect Split Detection Protocol

## Penalties
- Merging dependent before splits complete: -50%
- Merging splits out of order: -40%
- Not detecting violations during monitoring: -30%
- Creating incorrect merge plans: -25%
- Not validating before merge: -20%

## Key Takeaway
**NEVER** merge a dependent effort until ALL splits of its dependencies are merged. Dependencies are satisfied at the EFFORT level only when the COMPLETE effort (all splits) is integrated.

---
*Rule Type*: Protocol
*Agents*: Integration Agent, Code Reviewer, Orchestrator
*Enforcement*: Pre-merge validation, monitoring, and review