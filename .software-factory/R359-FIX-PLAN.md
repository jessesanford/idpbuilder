# Fix Plan for E1.2.2-split-001 (R359 SUPREME LAW VIOLATION)

## Issue Summary
Split-001 catastrophically violated R359 by DELETING 380 lines of existing retry package code instead of partitioning the 900 lines of NEW work that was added in the parent effort.

## Root Cause
**FUNDAMENTAL MISUNDERSTANDING**: The SW Engineer incorrectly believed that splits should contain ONLY the portion of code for that split, deleting everything else. This is completely wrong!

**What Actually Happened:**
- Parent effort (E1.2.2-registry-authentication) added 900 lines of NEW code
- Split-001 should have contained ~500 lines of those NEW additions
- Instead, Split-001 DELETED the existing retry package (380 lines)
- This is a SUPREME LAW violation that could destroy the codebase

## The Correct Understanding of Splits

### What "Split" REALLY Means
When an effort exceeds 800 lines, we split the NEW WORK into pieces:
- Split-001: First ~500 lines of NEW code
- Split-002: Next ~400 lines of NEW code
- Each split ADDS to the codebase, never deletes

### Visual Example
```
BEFORE EFFORT (main branch): 10,000 lines existing code
AFTER EFFORT (full work):     10,900 lines (added 900 NEW lines)

CORRECT SPLITTING:
Split-001: 10,500 lines (10,000 existing + 500 NEW)
Split-002: 10,400 lines (10,000 existing + 400 NEW)
When merged: 10,900 lines total

WRONG (what happened):
Split-001: 9,620 lines (DELETED 380 lines!)
This is CATASTROPHIC!
```

## Correct Split Strategy for E1.2.2

The 900 lines of NEW authentication work should be divided as:

### Split-001 (Target: ~500 lines of NEW code)
Should contain:
- Core authentication types and interfaces
- Basic credential management structures
- Essential authentication methods
- Core error types
- Basic validation logic

### Split-002 (Target: ~400 lines of NEW code)
Should contain:
- Retry mechanism with backoff
- Advanced authentication features
- Additional helper methods
- Extended error handling
- Integration utilities

## Fix Instructions

### Step 1: Understand What Went Wrong
1. Read this entire document carefully
2. Read R359 rule at `rule-library/R359-code-deletion-prohibition.md`
3. Understand: Splits partition NEW work, NEVER delete existing code

### Step 2: Restore All Deleted Code
```bash
# From the split-001 directory
git checkout phase1/wave2/registry-authentication -- pkg/push/retry/
# This restores the ENTIRE retry package that was deleted
```

### Step 3: Identify the NEW Code to Keep
Review what was added in the parent effort and identify which ~500 lines belong in split-001:
- Focus on core authentication logic
- Keep foundational types and methods
- Save advanced features for split-002

### Step 4: Verify No Deletions
```bash
# This MUST show ONLY additions, NO deletions
git diff phase1/wave2/registry-authentication..HEAD --stat
# If you see any deletions, you're doing it wrong!
```

### Step 5: Measure Correctly
```bash
# Use the official line counter
PROJECT_ROOT=$(pwd)
while [ "$PROJECT_ROOT" != "/" ]; do
    [ -f "$PROJECT_ROOT/orchestrator-state.json" ] && break
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
done
$PROJECT_ROOT/tools/line-counter.sh
# Should show ~500 lines of NEW code
```

## Files to Restore (CRITICAL)
These files were DELETED and MUST be restored:
- `pkg/push/retry/` - ENTIRE package (380 lines)
  - All retry logic
  - Backoff algorithms
  - Error handling for retries

## Files to Keep (NEW work for split-001)
Focus on core authentication from the 900 NEW lines:
- Basic auth types and interfaces (~150 lines)
- Credential structures (~100 lines)
- Core authentication methods (~200 lines)
- Essential error types (~50 lines)

## Verification Steps

### 1. Check for Deletions (MUST BE ZERO)
```bash
deleted=$(git diff --numstat phase1/wave2/registry-authentication..HEAD | awk '{sum+=$2} END {print sum}')
if [ "$deleted" -gt 0 ]; then
    echo "FAIL: Still deleting $deleted lines!"
    exit 359
fi
```

### 2. Verify Size is Correct
```bash
# Should show ~500 lines of NEW code, not total repository size
$PROJECT_ROOT/tools/line-counter.sh
```

### 3. Ensure Compilation
```bash
go build ./...
go test ./...
```

## Critical Learning Points

1. **The 800-line limit applies ONLY to NEW code you ADD**
2. **NEVER delete existing code to fit the limit**
3. **Splits partition your NEW additions into manageable pieces**
4. **Each split builds ON TOP of existing code, not instead of it**
5. **The repository grows with each effort - that's EXPECTED**

## Severity Notice
This is a SUPREME LAW violation (R359). The penalty for this violation is -1000% and immediate termination. This fix plan MUST be followed exactly to correct the catastrophic misunderstanding.

## Expected Outcome
After fixing:
- Split-001 will ADD ~500 lines to the existing codebase
- NO deletions of any existing code
- The retry package will be fully restored
- The split will contain the first portion of NEW authentication work

Remember: You're adding a new floor to a building, not demolishing floors to make room!