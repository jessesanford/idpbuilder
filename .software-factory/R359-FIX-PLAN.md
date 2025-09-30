# Fix Plan for E1.2.3-split-001 (R359 SUPREME LAW VIOLATION)

## Issue Summary
Split-001 catastrophically violated R359 by DELETING 1085 lines of existing code (discovery.go, operations.go, pusher.go) instead of partitioning the 1706 lines of NEW work that was added in the parent effort.

## Root Cause
**CRITICAL MISUNDERSTANDING**: The SW Engineer believed splits should be self-contained units with ONLY their portion of code. This led to DELETING over 1000 lines of existing, approved code - a catastrophic violation that could destroy the project.

**What Actually Happened:**
- Parent effort (E1.2.3-image-push-operations) added 1706 lines of NEW code
- Split-001 should have contained ~550 lines of those NEW additions
- Instead, Split-001 DELETED three entire existing files (1085 lines)
- This is the WORST kind of R359 violation - deleting core functionality

## The Correct Understanding of Splits

### Critical Concept: Additive Development
```
EXISTING CODEBASE: Like a building with 10 floors
YOUR EFFORT:       Adds 3 new floors (1706 lines)
SPLITTING:         Build floor 11 first (550 lines)
                   Then floor 12 (550 lines)
                   Then floor 13 (606 lines)

NEVER:             Demolish floors 7-9 to make room!
```

### What This Means in Practice
```
BEFORE EFFORT (main branch): 10,000 lines existing code
AFTER FULL EFFORT:           11,706 lines (added 1706 NEW lines)

CORRECT SPLITTING:
Split-001: 10,550 lines (10,000 existing + 550 NEW)
Split-002: 10,550 lines (10,000 existing + 550 NEW)
Split-003: 10,606 lines (10,000 existing + 606 NEW)
When all merged: 11,706 lines total

WRONG (what happened):
Split-001: 8,915 lines (DELETED 1085 lines of existing code!)
This would DESTROY the project!
```

## Correct Split Strategy for E1.2.3

The 1706 lines of NEW push operations work should be divided as:

### Split-001 (Target: ~550 lines of NEW code)
Should contain:
- Core push command structure
- Basic configuration types
- Essential validation logic
- Core error types
- Fundamental push interfaces

### Split-002 (Target: ~550 lines of NEW code)
Should contain:
- Image discovery mechanisms
- Registry interaction logic
- Authentication integration
- Progress tracking basics

### Split-003 (Target: ~606 lines of NEW code)
Should contain:
- Advanced operations
- Retry and error recovery
- Detailed logging
- Performance optimizations
- Integration utilities

## Fix Instructions

### Step 1: Understand the Violation
1. Read this entire document twice
2. Study R359 at `rule-library/R359-code-deletion-prohibition.md`
3. Understand: You DELETED 1085 lines that should NEVER have been touched

### Step 2: Restore ALL Deleted Files
```bash
# From the split-001 directory
# Restore the THREE deleted files
git checkout phase1/wave2/image-push-operations -- pkg/push/discovery.go
git checkout phase1/wave2/image-push-operations -- pkg/push/operations.go
git checkout phase1/wave2/image-push-operations -- pkg/push/pusher.go

# Verify restoration
ls -la pkg/push/
# Should see discovery.go, operations.go, pusher.go restored
```

### Step 3: Identify Correct Split-001 Content
From the 1706 lines of NEW code added in parent effort, select ~550 lines:
- Focus on foundational push command code
- Keep basic structures and types
- Save complex operations for later splits

### Step 4: Verification (CRITICAL)
```bash
# This command MUST show ZERO deletions
git diff --numstat phase1/wave2/image-push-operations..HEAD
# Look at second column - MUST be 0 for all files

# If ANY deletions exist:
echo "FATAL: Still deleting code! R359 VIOLATION!"
exit 359
```

### Step 5: Correct Measurement
```bash
# Navigate to split directory and measure
PROJECT_ROOT=$(pwd)
while [ "$PROJECT_ROOT" != "/" ]; do
    [ -f "$PROJECT_ROOT/orchestrator-state.json" ] && break
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
done
$PROJECT_ROOT/tools/line-counter.sh
# Should show ~550 lines of NEW additions
```

## Files to Restore (CATASTROPHIC DELETIONS)
These files were COMPLETELY DELETED and MUST be restored immediately:

### 1. `pkg/push/discovery.go` (~350 lines)
- Image discovery logic
- Registry scanning
- Manifest detection

### 2. `pkg/push/operations.go` (~400 lines)
- Core push operations
- Transfer logic
- Progress tracking

### 3. `pkg/push/pusher.go` (~335 lines)
- Push orchestration
- Error handling
- State management

**TOTAL DELETED: 1085 lines of EXISTING, WORKING CODE**

## Files to Keep (NEW work for split-001)
From the 1706 NEW lines, keep ~550 lines focusing on:
- Push command cobra setup (~150 lines)
- Configuration structures (~100 lines)
- Basic validation (~100 lines)
- Core interfaces (~100 lines)
- Essential types (~100 lines)

## Verification Steps

### 1. Absolute Deletion Check
```bash
#!/bin/bash
deleted=$(git diff --numstat phase1/wave2/image-push-operations..HEAD | awk '{sum+=$2} END {print sum}')
if [ "$deleted" -gt 0 ]; then
    echo "🔴🔴🔴 FATAL R359 VIOLATION!"
    echo "You are STILL deleting $deleted lines!"
    echo "This is a SUPREME LAW violation!"
    exit 359
fi
echo "✅ No deletions detected"
```

### 2. File Existence Check
```bash
# All these files MUST exist
for file in discovery.go operations.go pusher.go; do
    if [ ! -f "pkg/push/$file" ]; then
        echo "🔴 MISSING: pkg/push/$file was deleted!"
        exit 359
    fi
done
echo "✅ All files restored"
```

### 3. Size Verification
```bash
$PROJECT_ROOT/tools/line-counter.sh
# Must show ~550 lines of NEW code
# NOT the total repository size
```

### 4. Build Test
```bash
go build ./...
go test ./...
# Must compile and pass tests
```

## Critical Learning Points

### The Devastating Impact
1. You deleted **1085 lines** of working code
2. Three entire files were removed
3. This would have broken ALL push functionality
4. The project would be unusable

### The Correct Mental Model
Think of splits like this:
- You're adding chapters to a book
- The book already has 100 chapters
- You need to add 3 new chapters (too long for one commit)
- So you add chapter 101 first, then 102, then 103
- You NEVER tear out chapters 70-80 to make room!

### Remember
- **Splits partition NEW work into smaller pieces**
- **Each split is additive to the existing codebase**
- **NEVER delete existing code for size reasons**
- **The repository GROWS with each effort - this is normal**

## Severity Notice
This is a SUPREME LAW R359 violation with the HIGHEST severity:
- **Lines Deleted**: 1085 (catastrophic)
- **Files Deleted**: 3 entire files
- **Penalty**: -1000% and immediate termination
- **Impact**: Would destroy the project

This fix plan MUST be followed EXACTLY. There is ZERO tolerance for deletions.

## Expected Outcome After Fix
- ALL 1085 lines of deleted code restored
- Split-001 contains ~550 lines of NEW push command code
- ZERO deletions from the parent branch
- Full compilation and test success
- Proper additive development

## Final Warning
Deleting existing code to meet size limits is like burning down rooms in your house to make it "smaller." It's not just wrong - it's destructive. NEVER do this again. Splits ADD functionality in pieces, they don't DELETE to make room.