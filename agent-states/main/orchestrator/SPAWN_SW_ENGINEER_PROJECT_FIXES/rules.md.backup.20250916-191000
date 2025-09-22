# Orchestrator - SPAWN_SW_ENGINEER_PROJECT_FIXES State Rules

## 🔴🔴🔴 R313 MANDATORY STOP AFTER SPAWNING AGENTS 🔴🔴🔴

**CONTEXT PRESERVATION LAW - PREVENTS RULE LOSS!**

The orchestrator MUST STOP IMMEDIATELY after spawning to preserve context:
- Record what was spawned in state file
- Save TODOs and commit state changes
- EXIT with clear continuation instructions
- This prevents agent responses from overflowing context

**STOP MEANS STOP - Exit after spawning, wait for /continue-orchestrating**

---

## State Context

**Purpose:**
Spawn SW Engineers to fix bugs identified during project integration per the fix plan created by Code Reviewer in WAITING_FOR_PROJECT_FIX_PLANS state.

## Primary Actions

1. **Read Fix Plan**:
   - Load PROJECT-FIX-PLAN.md
   - Identify spawn strategy (parallel vs sequential)
   - Extract fix assignments per engineer

2. **Navigate to Source Branches**:
   - Each bug must be fixed in its original source branch
   - Engineers work in effort directories, not integration branch
   - Respect R321: No fixes in integration branches

3. **Spawn Engineers**:
   ```bash
   # For each engineer in parallel group:
   cd /efforts/phase-X-wave-Y-effort-Z
   git checkout effort-branch-name
   
   # Copy fix instructions to effort directory
   cp PROJECT-FIX-PLAN.md ./FIX-INSTRUCTIONS.md
   
   # Spawn engineer with specific fix task
   echo "Spawning SW Engineer to fix bug #N in $(pwd)"
   echo "@agent-sw-engineer Please fix the following bug per instructions in FIX-INSTRUCTIONS.md"
   echo "Bug #N: [Title]"
   echo "After fixing, commit with message: 'fix: [description] (project integration bug #N)'"
   ```

4. **Update State File**:
   ```json
   {
     "project_fixes_in_progress": [
       {
         "bug_id": 1,
         "engineer": "sw-engineer-1",
         "branch": "phase-1-wave-2-effort-3",
         "status": "spawned"
       }
     ]
   }
   ```

## Spawn Protocol

### Parallel Spawning:
```markdown
## 🚀 SPAWNING SW ENGINEERS FOR PROJECT FIXES

### Parallel Group (Independent Fixes):
- **Engineer 1**: Fixing bug #1 in phase-1-wave-2-effort-3
- **Engineer 2**: Fixing bug #2 in phase-1-wave-1-effort-2
- **Engineer 3**: Fixing bug #4 in phase-2-wave-1-effort-1

All engineers spawned with timestamp: [TIMESTAMP]
```

### Sequential Notice:
```markdown
### Sequential Fixes Pending:
After parallel fixes complete, will spawn engineer for:
- Bug #6 (depends on bug #1 fix)
- Bug #7 (depends on bug #6 fix)
```

## Valid State Transitions

- **SPAWNED** → MONITORING_PROJECT_FIXES (engineers working on fixes)
- **SPAWN_FAILED** → ERROR_RECOVERY (unable to spawn engineers)
- **NO_FIX_PLAN** → ERROR_RECOVERY (missing PROJECT-FIX-PLAN.md)

## Critical Requirements

1. **MUST spawn in source branches** - Never in integration branch
2. **MUST stop after spawning** - R313 enforcement
3. **MUST record all spawns** - Track who's fixing what
4. **MUST provide fix instructions** - Each engineer needs clear direction
5. **MUST respect dependencies** - Don't spawn sequential fixes prematurely

## R313 Enforcement

After spawning:
1. Update orchestrator-state.json with spawn details
2. Save TODOs with spawn information
3. Commit and push state changes
4. Display continuation message
5. **EXIT IMMEDIATELY** - Do not wait for responses

## Grading Impact

- **-100%** if continuing after spawn (R313 violation)
- **-50%** if spawning in integration branch instead of source
- **-30%** if missing spawn tracking in state file
- **+15%** for optimal parallel spawning
- **+10%** for clear fix instructions to each engineer