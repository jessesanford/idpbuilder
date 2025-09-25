# SOFTWARE FACTORY 2.0: BRANCHING STRATEGY CLARIFICATION

## THE FUNDAMENTAL REQUIREMENT

**Every effort and split branch must be directly mergeable to main in sequence.**

This document clarifies the CORE branching and merging strategy that is FUNDAMENTAL to the Software Factory's architecture.

## 🔴🔴🔴 THE SUPREME PRINCIPLES 🔴🔴🔴

### Principle 1: Sequential Direct Mergeability (R363)
- **EVERY** effort merges to main DIRECTLY
- **NO** intermediary branches in the merge path
- **SEQUENTIAL** merging, not batch merging
- **INDEPENDENT** - each effort can merge alone

### Principle 2: Integration Testing Only (R364)
- Integration branches are **TESTING ENVIRONMENTS**
- They **NEVER** merge to main
- They validate interfaces and interactions
- They get **DELETED** after testing

### Principle 3: Independent Branch Mergeability (R307)
- Each effort must work in isolation
- Months can pass between merges
- Feature flags protect incomplete features
- No breaking changes ever

## THE CORRECT BRANCHING FLOW

### 1. Development Flow
```
main (always stable, always deployable)
  |
  ├── E1.1.1 (effort branch - developed independently)
  ├── E1.1.2 (effort branch - developed independently)
  ├── E1.1.3-split-001 (split branch)
  ├── E1.1.3-split-002 (split branch)
  └── E1.2.1 (effort branch - developed independently)
```

### 2. Testing Flow (Integration Branches)
```
Wave Integration Test:
  main
    |
    ├── wave1-integration-test-20250120 (TEMPORARY)
         ├── merge E1.1.1 (for testing)
         ├── merge E1.1.2 (for testing)
         └── merge E1.1.3 (for testing)

         RUN TESTS → PASS → DELETE BRANCH

Phase Integration Test:
  main
    |
    └── phase1-integration-test-20250120 (TEMPORARY)
         ├── merge all wave1 efforts
         ├── merge all wave2 efforts
         └── merge all wave3 efforts

         RUN TESTS → PASS → DELETE BRANCH
```

### 3. Merging Flow (Sequential to Main)
```
After Testing Passes:

Step 1: git checkout main && git merge E1.1.1
        ↓ (wait for CI/CD)
Step 2: git checkout main && git merge E1.1.2
        ↓ (wait for CI/CD)
Step 3: git checkout main && git merge E1.1.3-split-001
        ↓ (wait for CI/CD)
Step 4: git checkout main && git merge E1.1.3-split-002
        ↓ (wait for CI/CD)
Step 5: git checkout main && git merge E1.2.1

Each merge is INDEPENDENT and DIRECT to main
```

## WHAT THIS MEANS IN PRACTICE

### For Software Engineers
1. **Design for independence** - Your effort must work alone
2. **Use feature flags** - Protect incomplete features
3. **Never depend on unmerged work** - Each effort stands alone
4. **Test in isolation** - Ensure your effort works by itself

### For Integration Agents
1. **Create test branches only** - Never permanent branches
2. **Test thoroughly** - Build, test suites, demos
3. **Delete after testing** - Clean up test branches
4. **Never merge to main** - That's not your job

### For Orchestrators
1. **Coordinate sequential merging** - One effort at a time
2. **Wait between merges** - Let CI/CD process each one
3. **Never batch merge** - Each effort merges independently
4. **Track merge order** - Maintain proper sequence

### For Code Reviewers
1. **Verify independence** - Each effort must work alone
2. **Check feature flags** - Incomplete features must be flagged
3. **Create merge plans** - Define the sequence
4. **Never execute merges** - Just plan them

## COMMON MISUNDERSTANDINGS CORRECTED

### ❌ WRONG: "Integration branches merge to main"
**✅ CORRECT:** Integration branches are for testing only. They validate that efforts work together but are NEVER merged to main. Each effort merges directly to main.

### ❌ WRONG: "We merge all efforts at once for efficiency"
**✅ CORRECT:** Efforts merge sequentially, one at a time. This ensures each change is validated independently and can be rolled back if needed.

### ❌ WRONG: "Integration branches are the path to production"
**✅ CORRECT:** The path to production is: Effort → Main (directly). Integration branches are temporary testing environments.

### ❌ WRONG: "Wave integration contains all wave changes for main"
**✅ CORRECT:** Wave integration tests that wave efforts work together. After testing passes, each effort merges to main individually.

## TESTING VS MERGING

### Testing (What Integration Branches Do)
```bash
# Create temporary test environment
git checkout -b wave1-test-$(date +%s) main

# Merge all efforts for testing
git merge E1.1.1 E1.1.2 E1.1.3

# Run comprehensive tests
make build && make test && ./demo.sh

# If passes, delete test branch
git checkout main && git branch -D wave1-test-*

# Test proved they work together
```

### Merging (What Actually Goes to Main)
```bash
# After testing passes, merge sequentially
git checkout main && git merge E1.1.1 && git push
# Wait for CI/CD...

git checkout main && git merge E1.1.2 && git push
# Wait for CI/CD...

git checkout main && git merge E1.1.3 && git push
# Each effort proven independently
```

## THE BENEFITS OF THIS APPROACH

### 1. Clean History
- Linear progression of features
- Clear attribution of changes
- Easy to understand what changed when
- Simple bisection for debugging

### 2. Safe Rollbacks
- Can revert individual efforts
- No cascade failures
- Minimal impact radius
- Quick recovery

### 3. Continuous Deployment
- Each effort can deploy immediately
- No waiting for "integration windows"
- Faster value delivery
- Reduced deployment risk

### 4. True Independence
- Efforts developed in parallel
- No blocking dependencies
- Flexible merge timing
- Resilient to delays

## VALIDATION CHECKLIST

Before any effort merges to main:
- [ ] Builds successfully in isolation
- [ ] All tests pass independently
- [ ] Feature flags protect incomplete features
- [ ] No dependencies on unmerged work
- [ ] Integration testing completed (but branch deleted)
- [ ] Ready for sequential merge to main

## ENFORCEMENT RULES

The following rules enforce this strategy:
- **R363**: Sequential Direct Mergeability (SUPREME)
- **R364**: Integration Testing Only Branches (SUPREME)
- **R307**: Independent Branch Mergeability (PARAMOUNT)
- **R270**: No Integration Branches as Sources (SUPREME)
- **R306**: Merge Ordering with Splits Protocol
- **R291**: Integration Demo Requirement

## THE GOLDEN RULE

**"Test together, merge separately"**

Integration branches prove our code works together, but each effort must still earn its place in main through independent merit.

## SUMMARY

1. **Efforts merge directly to main** - No intermediaries
2. **Merges happen sequentially** - One at a time
3. **Integration branches test only** - Never merge to main
4. **Each effort stands alone** - True independence
5. **Clean, linear history** - Easy to understand and manage

This is not just a preference - it's the FOUNDATION of our continuous deployment architecture. Every agent, every process, every workflow depends on this fundamental branching strategy.

---

**Created**: 2025-01-20
**Authority**: Software Factory Architecture Team
**Enforcement**: MANDATORY - No exceptions