# Orchestrator - PROJECT_FIX_PLANNING State Rules

## 🛑🛑🛑 R322 MANDATORY STOP BEFORE STATE TRANSITIONS 🛑🛑🛑

**SUPREME LAW - VIOLATION = -100% IMMEDIATE FAILURE**

### YOU MUST STOP AFTER:
1. ✅ Completing all TODOs for this state
2. ✅ Updating orchestrator-state.json with new state
3. ✅ Committing and pushing the state file  
4. ✅ Providing work summary

### YOU MUST NOT:
- ❌ Continue to the next state automatically
- ❌ Start work for the new state
- ❌ Spawn agents for the new state
- ❌ Assume permission to continue

### STOP PROTOCOL:
```markdown
## 🛑 STATE TRANSITION CHECKPOINT: CURRENT_STATE → NEXT_STATE

### ✅ Current State Work Completed:
- [List completed work]

### 📊 Current Status:
- Current State: CURRENT_STATE
- Next State: NEXT_STATE
- TODOs Completed: X/Y
- State Files: Updated and committed ✅

### ⏸️ STOPPED - Awaiting User Continuation
Ready to transition to NEXT_STATE. Please use /continue-orchestrating.
```

**STOP MEANS STOP - Exit and wait for /continue-orchestrating**

---

## State Context

**Purpose:**
Create fix plans for bugs found and documented during project integration per R266.

## Primary Actions

1. **Analyze Bug Documentation**:
   - Read PROJECT-INTEGRATION-REPORT.md
   - Extract all bugs from "UPSTREAM BUGS IDENTIFIED" section
   - Categorize by severity and source branch

2. **Create Fix Plans**:
   - For each bug, determine fix approach
   - Identify which phase/wave/effort branch contains the bug
   - Create PROJECT-FIX-PLAN.md with detailed instructions

3. **Prepare for SW Engineer Spawn**:
   - Determine parallelization strategy for fixes
   - Group fixes by dependency
   - Prepare spawn instructions for each engineer

## Fix Plan Format

```markdown
# PROJECT INTEGRATION FIX PLAN

## Bug Summary
- Total Bugs Found: X
- Critical: Y
- High: Z
- Medium: A
- Low: B

## Fix Strategy

### Fix Group 1: [Independent Fixes - Can Parallelize]
#### Bug #1: [Title from report]
- **Source Branch**: phase-1-wave-2-effort-3
- **Fix Location**: src/api/handler.go:234-237
- **Fix Instructions**:
  ```go
  // Replace line 235:
  userData[5] = request.Extra  // BUG
  // With:
  if len(userData) > 5 {
      userData[5] = request.Extra
  }
  ```
- **Assigned To**: sw-engineer-1

#### Bug #2: [Title]
[Same format...]

### Fix Group 2: [Sequential Fixes - Dependencies]
[Bugs that must be fixed in order]

## Spawn Instructions

### Parallel Spawn:
- SW Engineer 1: Fix bugs #1, #3, #5 in respective branches
- SW Engineer 2: Fix bugs #2, #4 in respective branches

### Sequential Work:
- After parallel fixes complete, spawn engineer for bug #6
```

## Valid State Transitions

- **PLANS_CREATED** → SPAWN_SW_ENGINEER_PROJECT_FIXES (ready to spawn engineers)
- **NO_BUGS** → ERROR_RECOVERY (should not be in this state if no bugs)
- **PLANNING_FAILED** → ERROR_RECOVERY (unable to create fix plans)

## Critical Requirements

1. **NEVER fix bugs directly** - Only create plans
2. **Each bug must have clear fix instructions**
3. **Respect R321** - Fixes go to source branches, not integration branch
4. **Document which engineer gets which fix**
5. **Consider dependencies between fixes**

## Grading Impact

- **+20%** for comprehensive fix plans
- **+10%** for optimal parallelization strategy
- **-50%** if attempting to fix bugs directly
- **-30%** if fix plans are incomplete or unclear