# 🔴🔴🔴 RULE R363: Sequential Direct Mergeability to Main (SUPREME)

## Classification
- **Category**: CD/Integration/Core
- **Criticality**: SUPREME (Fundamental Architecture Principle)
- **Penalty**: -100% for violating merge strategy
- **Introduced**: v2.0.0

## THE FUNDAMENTAL PRINCIPLE

**EVERY effort and split branch MUST be directly mergeable to main in SEQUENCE.**

This is the FOUNDATION of our entire branching strategy:
1. **Efforts merge to main DIRECTLY** - No intermediary branches
2. **Efforts merge in SEQUENCE** - E1 → main, then E2 → main, then E3 → main
3. **Integration branches are for TESTING ONLY** - Never merge to main
4. **Each merge stands alone** - No batch merging, no integration merging

## The Absolute Architecture

```
main (production-ready at all times)
  ↑ Direct merge (after testing)
  E1.1.1 (effort branch)
  ↑ Direct merge (after E1.1.1)
  E1.1.2 (effort branch)
  ↑ Direct merge (after E1.1.2)
  E1.1.3-split-001
  ↑ Direct merge (after split-001)
  E1.1.3-split-002
  ↑ Direct merge (after split-002)
  E1.2.1 (effort branch)

Integration branches (TESTING ONLY):
  - wave1-integration (tests E1.1.1 + E1.1.2 work together)
  - phase1-integration (tests all phase1 efforts work together)
  - NEVER merge these to main
  - ONLY for validation
```

## What This Means

### ✅ CORRECT: Sequential Direct Merging
```bash
# Each effort merges to main directly, in sequence
git checkout main
git merge E1.1.1  # First effort merges
git push

git checkout main
git pull
git merge E1.1.2  # Second effort merges (after E1.1.1)
git push

git checkout main
git pull
git merge E1.1.3-split-001  # First split merges
git push

git checkout main
git pull
git merge E1.1.3-split-002  # Second split merges
git push
```

### ❌ WRONG: Merging Integration Branches
```bash
# NEVER DO THIS
git checkout main
git merge wave1-integration  # NO! Integration branches never merge to main
git merge phase1-integration # NO! Testing branches stay separate
```

### ❌ WRONG: Batch Merging
```bash
# NEVER DO THIS
git checkout main
git merge E1.1.1 E1.1.2 E1.1.3  # NO! Must be sequential, one at a time
```

## Why This Architecture

### 1. True Continuous Deployment
- Each effort can deploy independently
- No waiting for "integration cycles"
- Immediate value delivery
- Rollback is simple (revert one PR)

### 2. Clean History
- Linear progression of features
- Clear attribution (who did what)
- Easy bisection for bugs
- Simple revert strategy

### 3. Independent Validation
- Each effort proves it works alone
- Integration testing proves they work together
- But integration branches are throwaway
- Source of truth is always main + effort branches

### 4. Flexibility
- Can merge efforts with months between them
- Can skip efforts if not ready
- Can reorder if dependencies change
- Can pause and resume anytime

## Integration Branches: Testing Theater

Integration branches are **TEST ENVIRONMENTS**, not merge targets:

```bash
# Wave Integration (TESTING ONLY)
git checkout -b wave1-integration-test main
git merge E1.1.1
git merge E1.1.2
git merge E1.1.3
# Run tests, verify interfaces work
# Build, test, demo
# Then DELETE this branch - it served its purpose

# The ACTUAL merging happens sequentially to main:
git checkout main
git merge E1.1.1  # Direct to main
git merge E1.1.2  # Direct to main
git merge E1.1.3  # Direct to main
```

## Validation Requirements

### Before EVERY Effort Merge
```bash
validate_ready_for_main() {
    local effort_branch="$1"

    # 1. Must build in isolation
    git checkout main
    git checkout -b test-isolation
    git merge "$effort_branch"
    make build || return 1

    # 2. Must pass all tests
    make test || return 1

    # 3. Must not depend on unmerged work
    check_dependencies "$effort_branch" || return 1

    # 4. Must be next in sequence
    verify_merge_sequence "$effort_branch" || return 1

    echo "✅ $effort_branch ready for main"
}
```

### Integration Testing (Never Merges to Main)
```bash
test_wave_integration() {
    local wave="$1"

    # Create TEMPORARY test branch
    git checkout -b "wave${wave}-test-$(date +%s)" main

    # Merge all wave efforts for testing
    for effort in $(get_wave_efforts "$wave"); do
        git merge "$effort"
    done

    # Test that they work together
    make build
    make test
    ./demo-features.sh

    # Record results
    echo "Wave $wave integration: TESTED ✅" >> test-results.md

    # DELETE the test branch - it's not needed anymore
    git checkout main
    git branch -D "wave${wave}-test-*"

    # Now merge efforts to main SEQUENTIALLY
    for effort in $(get_wave_efforts "$wave"); do
        git checkout main
        git merge "$effort"
        git push
        sleep 5  # Let CI/CD process each one
    done
}
```

## Common Misunderstandings

### Misunderstanding 1: "Integration branches make merging easier"
**WRONG!** Integration branches are for TESTING only. They prove things work together but should NEVER merge to main. Each effort merges independently.

### Misunderstanding 2: "We should batch merge for efficiency"
**WRONG!** Sequential merging ensures each change is validated independently. Batch merging hides problems and makes rollback impossible.

### Misunderstanding 3: "Integration branches are the path to main"
**WRONG!** The path to main is:
- Effort branch → main (directly, sequentially)
- NOT: Effort → Integration → main

### Misunderstanding 4: "We need integration branches in the merge chain"
**WRONG!** Integration branches are disposable test environments. The merge chain is always: effort → main, effort → main, effort → main

## Enforcement

### Orchestrator Rules
```bash
# When managing merges to main
merge_effort_to_main() {
    local effort="$1"

    # MUST merge directly to main
    git checkout main
    git pull origin main

    # MUST verify it's next in sequence
    verify_next_in_sequence "$effort" || {
        echo "❌ $effort not next in sequence!"
        return 1
    }

    # MUST merge only this effort
    git merge "$effort" --no-ff

    # MUST push immediately
    git push origin main

    # MUST wait before next merge
    sleep 10  # Let CI/CD process
}
```

### Integration Agent Rules
```bash
# Integration agents create test branches ONLY
create_integration_test() {
    local level="$1"  # wave/phase

    # Create TEMPORARY branch for testing
    TEST_BRANCH="${level}-integration-test-$(date +%s)"
    git checkout -b "$TEST_BRANCH" main

    # Merge efforts for testing
    merge_efforts_for_testing

    # Run comprehensive tests
    run_integration_tests

    # Document results
    create_test_report

    # DELETE the test branch
    git checkout main
    git branch -D "$TEST_BRANCH"

    echo "✅ Integration tested. Efforts ready for sequential merge to main"
    echo "❌ NOT merging integration branch to main (R363)"
}
```

## The Golden Rules

1. **Every effort merges to main DIRECTLY**
2. **Merges happen in SEQUENCE, not batches**
3. **Integration branches are for TESTING, not merging**
4. **Each merge must stand alone**
5. **No effort depends on integration branches**

## Penalties

- Merging integration branch to main: **-100% IMMEDIATE FAILURE**
- Batch merging efforts: **-50%**
- Wrong merge sequence: **-40%**
- Creating dependencies on integration branches: **-50%**
- Not deleting test integration branches: **-20%**

## Success Metrics

- 100% of efforts merge directly to main
- 0 integration branches in main's history
- Sequential merge order maintained
- Each merge passes CI/CD independently
- Clean linear history in main

## Related Rules

- R307: Independent Branch Mergeability
- R270: No Integration Branches as Sources
- R306: Merge Ordering with Splits
- R364: Integration Testing Only Branches (companion rule)

---

**REMEMBER**: Integration branches are testing theaters where we verify our code plays well together. But the show goes to Broadway (main) one act at a time, in order, each standing on its own merit.