# Split Infrastructure Just-In-Time Creation Fix

## Executive Summary

Fixed a critical issue where the orchestrator was creating ALL split infrastructure at once, violating the sequential split principle. The solution implements just-in-time creation where each split's infrastructure is created only after the previous split is complete.

## The Problem

### Previous Behavior (WRONG)
The orchestrator was creating all split infrastructure simultaneously:
```
1. Code Reviewer creates split plans (001, 002, 003)
2. Orchestrator creates ALL infrastructure at once:
   - split-001 based on main
   - split-002 based on main (WRONG!)
   - split-003 based on main (WRONG!)
3. When measuring split-002, it appeared to have 800+ lines
   (because it included split-001's work in the diff)
```

### Root Cause
- R204 specified creating "ALL split infrastructure BEFORE spawning SW engineering agents"
- This made sequential branching impossible since split-001 hadn't been implemented yet
- Split-002 couldn't be based on split-001 if split-001 had no code

## The Solution

### New Behavior (CORRECT)
Implemented just-in-time infrastructure creation:
```
1. Code Reviewer creates split plans (001, 002, 003)
2. Orchestrator creates ONLY split-001 infrastructure
   - Based on incremental base (per R308)
3. SW Engineer implements split-001
4. After split-001 complete, orchestrator creates split-002
   - Based on split-001 branch (now has code!)
5. SW Engineer implements split-002
6. After split-002 complete, orchestrator creates split-003
   - Based on split-002 branch
7. Continue pattern for all splits
```

## Implementation Details

### 1. New State: CREATE_NEXT_SPLIT_INFRASTRUCTURE
Added a new orchestrator state that:
- Creates infrastructure for EXACTLY ONE split
- Determines correct base branch (first split vs subsequent)
- Updates split tracking in state file
- Transitions to SPAWN_AGENTS for implementation

### 2. State Machine Updates
**Added transitions:**
```
MONITOR → CREATE_NEXT_SPLIT_INFRASTRUCTURE (when next split needed)
CREATE_NEXT_SPLIT_INFRASTRUCTURE → SPAWN_AGENTS (infrastructure ready)
```

**Flow for splits:**
```
1. MONITOR detects split plans exist
2. Transition to CREATE_NEXT_SPLIT_INFRASTRUCTURE
3. Create infrastructure for split-001 only
4. Transition to SPAWN_AGENTS
5. Spawn SW Engineer for split-001
6. Return to MONITOR
7. When split-001 complete, detect next split needed
8. Transition to CREATE_NEXT_SPLIT_INFRASTRUCTURE
9. Create infrastructure for split-002 (based on split-001)
10. Continue pattern...
```

### 3. Updated Rules

#### R204 - Orchestrator Split Infrastructure
- Changed from "create ALL infrastructure" to "create infrastructure JUST-IN-TIME"
- Emphasized sequential creation pattern
- Added tracking of current split in progress

#### MONITOR State Rules
- Added detection for when next split infrastructure needed
- Transition logic to CREATE_NEXT_SPLIT_INFRASTRUCTURE
- Split tracking updates

## Benefits

### 1. Correct Line Counting
Each split now shows ONLY its own additions:
- split-001: 400 lines (from main)
- split-002: 380 lines (incremental from split-001)
- split-003: 420 lines (incremental from split-002)

### 2. True Sequential Dependencies
- split-002 can use code from split-001
- split-003 can use code from split-001 and split-002
- No merge conflicts between splits

### 3. Clean Integration
Sequential merge order is natural:
```
1. Merge split-003 → split-002
2. Merge split-002 → split-001
3. Merge split-001 → main
```

## Verification Steps

### 1. Check Split Creation
```bash
# Verify only one split directory exists initially
ls -la /efforts/phase1/wave1/ | grep SPLIT
# Should show only auth-SPLIT-001

# After split-001 completes, check again
ls -la /efforts/phase1/wave1/ | grep SPLIT
# Should now show auth-SPLIT-001 and auth-SPLIT-002
```

### 2. Verify Branch Dependencies
```bash
# Check split-002 is based on split-001
cd /efforts/phase1/wave1/auth-SPLIT-002
git log --oneline | grep split-001
# Should show commits from split-001
```

### 3. Verify Line Counts
```bash
# Each split should show reasonable line count (tool auto-detects base)
cd /efforts/phase1/wave1/auth-SPLIT-002
../../tools/line-counter.sh
# Tool will show: 🎯 Detected base: split-001
# Should show ~400 lines, not 800+
```

## Files Modified

1. **SOFTWARE-FACTORY-STATE-MACHINE.md**
   - Added CREATE_NEXT_SPLIT_INFRASTRUCTURE state
   - Updated state transitions
   - Modified Split Infrastructure Flow section

2. **rule-library/R204-orchestrator-split-infrastructure.md**
   - Changed to just-in-time creation model
   - Updated code examples
   - Emphasized sequential dependency

3. **agent-states/orchestrator/MONITOR/rules.md**
   - Added split infrastructure detection logic
   - Transition to CREATE_NEXT_SPLIT_INFRASTRUCTURE

4. **agent-states/orchestrator/CREATE_NEXT_SPLIT_INFRASTRUCTURE/rules.md**
   - New file with complete state rules
   - Implementation functions
   - Tracking updates

## Impact on Grading

This fix ensures compliance with:
- **R308**: Incremental branching strategy (splits build on each other)
- **R306**: Merge ordering with splits (sequential dependencies)
- **R302**: Comprehensive split tracking (one at a time)
- **R204**: Split infrastructure creation (just-in-time)

**Penalty Avoided**: -100% for violating sequential split principles

## Migration Guide

For existing efforts with splits already created:
1. No action needed if splits are complete
2. If splits in progress, may need to recreate incorrectly based splits
3. Check with `git merge-base` to verify correct branching

## Summary

The just-in-time split infrastructure creation ensures:
- ✅ Each split is based on the previous split
- ✅ Line counts are accurate (not cumulative)
- ✅ Sequential dependencies work correctly
- ✅ Clean merge strategy without conflicts
- ✅ Compliance with all split-related rules

This is a CRITICAL fix that prevents split measurement errors and ensures proper sequential implementation of split efforts.