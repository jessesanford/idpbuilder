# 🔴🔴🔴 CRITICAL: Split Sequential Branching Clarification 🔴🔴🔴

## Executive Summary
Splits MUST be created SEQUENTIALLY, with each split based on the previous one. This is MANDATORY for correct line counting and clean integration.

## The Problem We're Solving

When an effort exceeds 800 lines and needs splitting, there are two possible branching strategies:

1. **WRONG: Parallel Branching** - All splits from same base
2. **CORRECT: Sequential Branching** - Each split from previous

The system REQUIRES sequential branching to function correctly.

## Visual Comparison

### ❌ WRONG: Parallel Branching (FORBIDDEN)
```
phase-integration (base branch)
    ├── split-001 (400 lines from phase-integration)
    ├── split-002 (400 lines from phase-integration) ← WRONG!
    └── split-003 (400 lines from phase-integration) ← WRONG!
```

**Problems with Parallel:**
- Split-002 can't use code from split-001
- Merge conflicts between splits
- Line counter shows cumulative totals:
  - split-001: 400 lines ✓
  - split-002: 800 lines ✗ (includes split-001!)
  - split-003: 1200 lines ✗ (includes all!)

### ✅ CORRECT: Sequential Branching (MANDATORY)
```
phase-integration (base branch)
    └── split-001 (400 lines from phase-integration)
            └── split-002 (400 lines from split-001)
                    └── split-003 (400 lines from split-002)
```

**Benefits of Sequential:**
- Each split measures ONLY its additions
- Later splits can use earlier code
- Clean progressive integration
- Line counter shows correct totals:
  - split-001: 400 lines ✓
  - split-002: 400 lines ✓
  - split-003: 400 lines ✓

## Real-World Example: Authentication System

### Scenario
Original effort "authentication" has 1200 lines - needs 3 splits.

### ❌ WRONG Implementation (Parallel)
```bash
# Developer creates all splits from phase-integration
git checkout phase-integration
git checkout -b project/phase1/wave1/authentication--split-001
# Implement core auth (400 lines)
./tools/line-counter.sh
# Shows: 400 lines ✓

git checkout phase-integration  # ← WRONG! Should checkout split-001
git checkout -b project/phase1/wave1/authentication--split-002
# Implement session management (400 lines)
# But can't use core auth from split-001!
# Have to duplicate some code
./tools/line-counter.sh
# Shows: 800 lines ✗ (includes duplicated code!)

git checkout phase-integration  # ← WRONG! Should checkout split-002
git checkout -b project/phase1/wave1/authentication--split-003
# Implement OAuth (400 lines)
# Can't use core auth OR session management!
# More duplication needed
./tools/line-counter.sh
# Shows: 1200 lines ✗ (all code!)
```

**Result**: Size violations, code duplication, merge conflicts

### ✅ CORRECT Implementation (Sequential)
```bash
# Split 1: Core authentication
git checkout phase-integration
git checkout -b project/phase1/wave1/authentication--split-001
# Implement core auth (400 lines)
./tools/line-counter.sh
# Shows: 400 lines ✓

# Split 2: Session management (based on split-001!)
git checkout project/phase1/wave1/authentication--split-001
git checkout -b project/phase1/wave1/authentication--split-002
# Implement session management (400 lines)
# CAN use core auth from split-001!
./tools/line-counter.sh
# Shows: 400 lines ✓ (only the additions)

# Split 3: OAuth (based on split-002!)
git checkout project/phase1/wave1/authentication--split-002
git checkout -b project/phase1/wave1/authentication--split-003
# Implement OAuth (400 lines)
# CAN use both core auth AND session management!
./tools/line-counter.sh
# Shows: 400 lines ✓ (only the additions)
```

**Result**: Clean separation, no duplication, perfect line counts

## Integration Sequence

### With Sequential Branching (CORRECT)
```bash
# Integration is simple - merge in order
git checkout phase-integration
git merge project/phase1/wave1/authentication--split-001  # Adds core auth
git merge project/phase1/wave1/authentication--split-002  # Adds sessions
git merge project/phase1/wave1/authentication--split-003  # Adds OAuth
# Result: Complete authentication system, no conflicts
```

### With Parallel Branching (WRONG)
```bash
# Integration has conflicts
git checkout phase-integration
git merge project/phase1/wave1/authentication--split-001  # OK
git merge project/phase1/wave1/authentication--split-002  # CONFLICTS!
# Both splits modified same base independently
git merge project/phase1/wave1/authentication--split-003  # MORE CONFLICTS!
# Result: Manual conflict resolution needed
```

## Validation Script

Add this to your workflow to verify sequential branching:

```bash
#!/bin/bash
# verify-sequential-splits.sh

verify_sequential_splits() {
    local effort_name="$1"
    local split_count="$2"
    
    echo "Verifying sequential branching for $effort_name..."
    
    for i in $(seq 2 $split_count); do
        CURRENT=$(printf "%03d" $i)
        PREVIOUS=$(printf "%03d" $((i-1)))
        
        CURRENT_BRANCH="${effort_name}--split-${CURRENT}"
        PREVIOUS_BRANCH="${effort_name}--split-${PREVIOUS}"
        
        # Check if current is based on previous
        if git merge-base --is-ancestor "$PREVIOUS_BRANCH" "$CURRENT_BRANCH"; then
            echo "✅ split-${CURRENT} correctly based on split-${PREVIOUS}"
        else
            echo "❌ ERROR: split-${CURRENT} NOT based on split-${PREVIOUS}!"
            echo "   This violates sequential branching requirement!"
            return 1
        fi
    done
    
    echo "✅ All splits follow sequential branching"
    return 0
}

# Example usage
verify_sequential_splits "project/phase1/wave1/authentication" 3
```

## Rules Updated

The following rules have been updated to enforce sequential branching:

1. **R302**: Comprehensive Split Tracking Protocol
   - Added `split_strategy: "SEQUENTIAL"` field
   - Added `base_branch` tracking for each split
   - Added verification functions

2. **R204**: Orchestrator Split Infrastructure
   - Changed from "optionally sequential" to "MANDATORY sequential"
   - Added explicit FORBIDDEN patterns
   - Updated implementation examples

3. **SW Engineer SPLIT_IMPLEMENTATION State**
   - Added sequential branching diagram
   - Updated measurement logic for correct base
   - Added verification steps

4. **Code Reviewer CREATE_SPLIT_PLAN State**
   - Added branching strategy section requirement
   - Clarified dependency descriptions
   - Added sequential pattern to plans

## Key Takeaways

1. **Sequential is MANDATORY** - Not optional, not recommended - REQUIRED
2. **Each split builds on previous** - Creates a chain of development
3. **Line counts are accurate** - Each split shows only its additions
4. **Integration is clean** - No conflicts between splits
5. **Dependencies work** - Later splits can use earlier code

## Enforcement

- **Orchestrator**: Creates splits with sequential branching
- **SW Engineer**: Verifies correct base before implementing
- **Code Reviewer**: Plans splits with sequential strategy
- **Line Counter**: Measures against correct base
- **Integration**: Applies splits in sequence

## Questions?

If you see any split implementation that doesn't follow sequential branching:
1. **STOP** immediately
2. **REPORT** the violation
3. **FIX** the branching before proceeding
4. **VERIFY** with the validation script

Remember: Sequential branching is not a suggestion - it's a REQUIREMENT for the system to function correctly!