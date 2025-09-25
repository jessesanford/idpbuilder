# 🔴🔴🔴 MERGE ORDERING WITH SPLITS - COMPLETE GUIDE 🔴🔴🔴

## SUPREME LAW: Dependencies Work at EFFORT Level, Not Split Level

When an effort is split into multiple branches, ALL splits must be merged before any dependent efforts can be merged. This is NON-NEGOTIABLE.

## Core Principles

### 1. Dependencies are at EFFORT Level
- If Effort E2 depends on Effort E1, it depends on **ALL of E1**
- ALL E1 splits must merge before E2 can merge
- You cannot cherry-pick just part of a dependency

### 2. Splits Merge Sequentially  
- Splits are created sequentially (each based on previous)
- Splits MUST merge in the same sequential order
- Split-001 → Split-002 → Split-003 (never out of order)

### 3. Integration is Progressive
- Each integration builds on the previous
- No parallel merging of splits from same effort
- Dependencies wait for complete prerequisites

## 🚨🚨🚨 CRITICAL MERGE ORDER RULES 🚨🚨🚨

### Rule 1: Complete All Splits Before Dependents

**CORRECT Pattern:**
```
E1 has 3 splits, E2 depends on E1:
1. Merge E1-split-001
2. Merge E1-split-002  
3. Merge E1-split-003
4. NOW merge E2 (has complete E1)
```

**WRONG Pattern:**
```
❌ FORBIDDEN:
1. Merge E1-split-001
2. Merge E2 (WRONG - missing E1-split-002 and 003!)
3. Merge E1-split-002
4. Merge E1-split-003
```

### Rule 2: Sequential Split Merging

Since splits are created sequentially (R302):
- split-002 is based on split-001
- split-003 is based on split-002

Therefore, they MUST merge in order:
```bash
# Correct merge sequence
git merge origin/E1-split-001 --no-ff
git merge origin/E1-split-002 --no-ff  # Contains split-001 changes
git merge origin/E1-split-003 --no-ff  # Contains split-001 and 002 changes
```

### Rule 3: Dependency Chain Resolution

When multiple efforts have dependencies and splits:
```
E1 (3 splits) ← E2 depends ← E3 depends
E4 (2 splits) ← E5 depends

Merge Order:
1. E1-split-001
2. E1-split-002
3. E1-split-003  
4. E2 (now has complete E1)
5. E3 (now has complete E2)
6. E4-split-001
7. E4-split-002
8. E5 (now has complete E4)
```

## Scenario Examples

### Scenario 1: Simple Split with Dependency

**Setup:**
- E1.1.1: Authentication (1200 lines → 3 splits)
- E1.1.2: User Profile (depends on E1.1.1)
- E1.1.3: Settings (independent)

**Correct Merge Order:**
```yaml
merge_sequence:
  - step: 1
    branch: phase1/wave1/E1.1.1-authentication-split-001
    type: split
    effort: E1.1.1
    description: Core auth functions
    
  - step: 2
    branch: phase1/wave1/E1.1.1-authentication-split-002
    type: split
    effort: E1.1.1
    description: Session management
    
  - step: 3
    branch: phase1/wave1/E1.1.1-authentication-split-003
    type: split
    effort: E1.1.1
    description: OAuth providers
    
  - step: 4
    branch: phase1/wave1/E1.1.2-user-profile
    type: effort
    depends_on: [E1.1.1]  # Now complete!
    description: User profile using auth
    
  - step: 5
    branch: phase1/wave1/E1.1.3-settings
    type: effort
    depends_on: []  # Independent
    description: Application settings
```

### Scenario 2: Multiple Splits and Dependencies

**Setup:**
- E2.1.1: Database Layer (1500 lines → 3 splits)
- E2.1.2: API Endpoints (900 lines → 2 splits, depends on E2.1.1)
- E2.1.3: Frontend (depends on E2.1.2)

**Correct Merge Order:**
```yaml
merge_sequence:
  # First: Complete Database Layer
  - step: 1
    branch: phase2/wave1/E2.1.1-database-split-001
    description: Schema and migrations
    
  - step: 2
    branch: phase2/wave1/E2.1.1-database-split-002
    description: CRUD operations
    
  - step: 3
    branch: phase2/wave1/E2.1.1-database-split-003
    description: Complex queries
    
  # Second: Complete API (depends on COMPLETE database)
  - step: 4
    branch: phase2/wave1/E2.1.2-api-split-001
    description: Basic endpoints
    
  - step: 5
    branch: phase2/wave1/E2.1.2-api-split-002
    description: Advanced endpoints
    
  # Third: Frontend (depends on COMPLETE API)
  - step: 6
    branch: phase2/wave1/E2.1.3-frontend
    description: UI using complete API
```

### Scenario 3: Diamond Dependencies with Splits

**Setup:**
```
     E1 (2 splits)
       /        \
    E2            E3 (3 splits)
       \        /
         E4
```

**Correct Merge Order:**
```yaml
merge_sequence:
  # Complete E1 first (common dependency)
  - step: 1
    branch: E1-split-001
  - step: 2
    branch: E1-split-002
    
  # E2 and E3 can go in either order, but splits must be complete
  - step: 3
    branch: E2  # Single branch, depends on complete E1
    
  - step: 4
    branch: E3-split-001
  - step: 5
    branch: E3-split-002
  - step: 6
    branch: E3-split-003
    
  # E4 needs BOTH E2 and complete E3
  - step: 7
    branch: E4
```

## Validation Functions

### 1. Validate Merge Readiness
```bash
validate_merge_readiness() {
    local EFFORT=$1
    local STATE_FILE="orchestrator-state.json"
    
    # Get dependencies for this effort
    DEPENDENCIES=$(yq ".efforts.\"$EFFORT\".dependencies[]" "$STATE_FILE" 2>/dev/null)
    
    for DEP in $DEPENDENCIES; do
        echo "Checking dependency: $DEP"
        
        # Check if dependency has splits
        SPLIT_COUNT=$(yq ".split_tracking.\"$DEP\".split_count // 0" "$STATE_FILE")
        
        if [[ $SPLIT_COUNT -gt 0 ]]; then
            echo "  $DEP has $SPLIT_COUNT splits"
            
            # Verify ALL splits are merged
            for i in $(seq 1 $SPLIT_COUNT); do
                SPLIT_NUM=$(printf "%03d" $i)
                SPLIT_BRANCH="${DEP}-split-${SPLIT_NUM}"
                SPLIT_STATUS=$(yq ".split_tracking.\"$DEP\".split_branches[] | select(.branch == \"*$SPLIT_BRANCH*\") | .status" "$STATE_FILE")
                
                if [[ "$SPLIT_STATUS" != "INTEGRATED" ]]; then
                    echo "❌ ERROR: Cannot merge $EFFORT"
                    echo "   Missing split: $SPLIT_BRANCH from dependency $DEP"
                    echo "   Status: $SPLIT_STATUS (needs INTEGRATED)"
                    return 1
                fi
            done
            echo "  ✅ All $SPLIT_COUNT splits of $DEP are integrated"
        else
            # Check if non-split dependency is merged
            STATUS=$(yq ".efforts_completed.\"$DEP\".integration_status" "$STATE_FILE")
            if [[ "$STATUS" != "INTEGRATED" ]]; then
                echo "❌ ERROR: Cannot merge $EFFORT"
                echo "   Dependency $DEP not integrated (status: $STATUS)"
                return 1
            fi
            echo "  ✅ $DEP is integrated"
        fi
    done
    
    echo "✅ $EFFORT is ready to merge (all dependencies complete)"
    return 0
}
```

### 2. Generate Correct Merge Order
```bash
generate_merge_order() {
    local WAVE=$1
    local STATE_FILE="orchestrator-state.json"
    local MERGE_ORDER=()
    
    echo "Generating merge order for Wave $WAVE..."
    
    # Get all efforts for this wave
    EFFORTS=$(yq ".waves.wave${WAVE}.efforts[]" "$STATE_FILE")
    
    # First pass: Collect all branches (splits or originals)
    for effort in $EFFORTS; do
        SPLIT_COUNT=$(yq ".split_tracking.\"$effort\".split_count // 0" "$STATE_FILE")
        
        if [[ $SPLIT_COUNT -gt 0 ]]; then
            # Add all splits in sequence
            for i in $(seq 1 $SPLIT_COUNT); do
                SPLIT_NUM=$(printf "%03d" $i)
                SPLIT_BRANCH=$(yq ".split_tracking.\"$effort\".split_branches[$((i-1))].branch" "$STATE_FILE")
                MERGE_ORDER+=("$SPLIT_BRANCH|$effort|split-$SPLIT_NUM")
            done
        else
            # Add original branch
            BRANCH=$(yq ".efforts_completed.\"$effort\".branch" "$STATE_FILE")
            MERGE_ORDER+=("$BRANCH|$effort|original")
        fi
    done
    
    # Sort by dependencies (topological sort)
    # This is simplified - real implementation would do full topological sort
    
    echo "Merge Order:"
    for entry in "${MERGE_ORDER[@]}"; do
        IFS='|' read -r branch effort type <<< "$entry"
        echo "  - $branch ($effort - $type)"
    done
}
```

### 3. Verify Split Sequence Integrity
```bash
verify_split_sequence() {
    local EFFORT=$1
    local STATE_FILE="orchestrator-state.json"
    
    SPLITS=$(yq ".split_tracking.\"$EFFORT\".split_branches[].branch" "$STATE_FILE")
    
    PREV_SPLIT=""
    COUNT=0
    for split in $SPLITS; do
        ((COUNT++))
        echo "Checking split $COUNT: $split"
        
        if [[ $COUNT -eq 1 ]]; then
            # First split should be based on integration branch
            BASE=$(yq ".split_tracking.\"$EFFORT\".split_branches[0].base_branch" "$STATE_FILE")
            echo "  Base: $BASE"
        else
            # Subsequent splits should be based on previous split
            EXPECTED_BASE=$PREV_SPLIT
            ACTUAL_BASE=$(yq ".split_tracking.\"$EFFORT\".split_branches[$((COUNT-1))].base_branch" "$STATE_FILE")
            
            if [[ "$ACTUAL_BASE" != *"$EXPECTED_BASE"* ]]; then
                echo "❌ ERROR: Split sequence broken!"
                echo "   Expected base: $EXPECTED_BASE"
                echo "   Actual base: $ACTUAL_BASE"
                return 1
            fi
            echo "  ✅ Correctly based on: $EXPECTED_BASE"
        fi
        
        PREV_SPLIT=$(basename "$split")
    done
    
    echo "✅ Split sequence verified for $EFFORT"
}
```

## Common Mistakes and How to Avoid Them

### Mistake 1: Merging Dependent Before All Splits
```bash
# ❌ WRONG: E2 merged before E1 is complete
git merge E1-split-001
git merge E2  # ERROR: E1-split-002 and 003 not merged yet!

# ✅ CORRECT: Complete E1 first
git merge E1-split-001
git merge E1-split-002
git merge E1-split-003
git merge E2  # Now E1 is complete
```

### Mistake 2: Merging Splits Out of Order
```bash
# ❌ WRONG: Random split order
git merge E1-split-002  # ERROR: split-001 not merged yet!
git merge E1-split-001
git merge E1-split-003

# ✅ CORRECT: Sequential order
git merge E1-split-001
git merge E1-split-002
git merge E1-split-003
```

### Mistake 3: Treating Splits as Independent Efforts
```yaml
# ❌ WRONG: Treating splits as separate efforts
dependencies:
  E2: [E1-split-001]  # NO! Dependencies are on efforts, not splits

# ✅ CORRECT: Dependencies on complete efforts
dependencies:
  E2: [E1]  # Means ALL of E1, including all splits
```

### Mistake 4: Parallel Split Creation
```bash
# ❌ WRONG: Creating splits in parallel from same base
git checkout main
git checkout -b E1-split-001
git checkout main  # WRONG!
git checkout -b E1-split-002  # Now both splits have same base

# ✅ CORRECT: Sequential split creation (R302)
git checkout main
git checkout -b E1-split-001
# implement...
git checkout E1-split-001  # Base next split on previous!
git checkout -b E1-split-002
```

## Integration Agent Checklist

When executing merges, the Integration Agent MUST:

### Pre-Merge Validation
- [ ] Read split_tracking section for all efforts
- [ ] Identify which efforts have splits
- [ ] Verify split sequence integrity
- [ ] Check all dependencies are complete

### During Merge Execution
- [ ] Merge splits in sequential order (001, 002, 003...)
- [ ] Complete all splits before dependent efforts
- [ ] Never skip a split number
- [ ] Document each merge in work-log.md

### Post-Merge Validation
- [ ] Verify all planned merges completed
- [ ] Update integration_status for all efforts
- [ ] Run tests to verify integration
- [ ] Generate integration report

## Test Cases for Merge Ordering

### Test 1: Simple Split Dependency
```yaml
test_case: simple_split_dependency
setup:
  E1: 
    splits: [split-001, split-002]
  E2:
    depends_on: [E1]
    
valid_orders:
  - [E1-split-001, E1-split-002, E2]
  
invalid_orders:
  - [E1-split-001, E2, E1-split-002]  # E2 before E1 complete
  - [E1-split-002, E1-split-001, E2]  # Splits out of order
```

### Test 2: Multiple Dependencies with Splits
```yaml
test_case: multiple_dependencies
setup:
  E1:
    splits: [split-001, split-002, split-003]
  E2:
    depends_on: [E1]
  E3:
    splits: [split-001, split-002]
    depends_on: [E1]
  E4:
    depends_on: [E2, E3]
    
valid_orders:
  - [E1-s1, E1-s2, E1-s3, E2, E3-s1, E3-s2, E4]
  - [E1-s1, E1-s2, E1-s3, E3-s1, E3-s2, E2, E4]
  
invalid_orders:
  - [E1-s1, E2, E1-s2, E1-s3, E3-s1, E3-s2, E4]  # E2 too early
  - [E1-s1, E1-s2, E1-s3, E2, E3-s1, E4, E3-s2]  # E4 before E3 complete
```

### Test 3: Diamond with Splits
```yaml
test_case: diamond_dependency
setup:
  E1:
    splits: [split-001, split-002]
  E2:
    depends_on: [E1]
  E3:
    depends_on: [E1]
    splits: [split-001, split-002, split-003]
  E4:
    depends_on: [E2, E3]
    
valid_orders:
  - [E1-s1, E1-s2, E2, E3-s1, E3-s2, E3-s3, E4]
  - [E1-s1, E1-s2, E3-s1, E3-s2, E3-s3, E2, E4]
  
invalid_orders:
  - [E1-s1, E2, E1-s2, E3-s1, E3-s2, E3-s3, E4]  # E2 before E1 complete
  - [E1-s1, E1-s2, E2, E3-s1, E4, E3-s2, E3-s3]  # E4 before E3 complete
```

## Error Messages and Recovery

### ERROR: Incomplete Dependency
```
❌ CRITICAL: Cannot merge effort E2.1.2
Dependency E2.1.1 has 3 splits but only 2 are integrated:
  ✅ E2.1.1-split-001 (INTEGRATED)
  ✅ E2.1.1-split-002 (INTEGRATED)
  ❌ E2.1.1-split-003 (COMPLETED - not integrated)
  
Action Required:
1. Merge E2.1.1-split-003 first
2. Then retry merging E2.1.2
```

### ERROR: Out of Order Split
```
❌ CRITICAL: Cannot merge E1-split-002
Previous split not integrated:
  ❌ E1-split-001 (COMPLETED - not integrated)
  
Splits must be merged sequentially!

Action Required:
1. Merge E1-split-001 first
2. Then merge E1-split-002
```

### ERROR: Using Original Instead of Splits
```
❌ BLOCKED: Attempting to merge deprecated branch
Branch: phase1/wave1/E1.1.1-authentication
Status: SPLIT_DEPRECATED

This effort was split! Use these branches instead:
  - phase1/wave1/E1.1.1-authentication-split-001
  - phase1/wave1/E1.1.1-authentication-split-002
  - phase1/wave1/E1.1.1-authentication-split-003
```

## Related Rules and Documentation

- **R302**: Comprehensive Split Tracking Protocol
- **R204**: Orchestrator Split Infrastructure
- **R262**: Merge Operation Protocols
- **R269**: Code Reviewer Merge Plan No Execution
- **R297**: Architect Split Detection Protocol

## Summary

The key principle is simple but absolute:
**Dependencies work at the EFFORT level, not the SPLIT level.**

When an effort has splits:
1. ALL splits must be merged
2. In sequential order (001, 002, 003...)
3. BEFORE any dependent efforts

This ensures:
- Complete functionality is available to dependents
- No partial implementations cause failures
- Integration is predictable and reliable
- Split boundaries are respected

---
*Document Version*: 1.0
*Last Updated*: 2025-01-20
*Status*: AUTHORITATIVE