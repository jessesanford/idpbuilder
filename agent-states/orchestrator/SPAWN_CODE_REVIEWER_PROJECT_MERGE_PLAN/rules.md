# Orchestrator - SPAWN_CODE_REVIEWER_PROJECT_MERGE_PLAN State Rules

## 🛑🛑🛑 R313 MANDATORY STOP AFTER SPAWNING AGENTS 🛑🛑🛑

**SUPREME LAW - VIOLATION = -100% IMMEDIATE FAILURE**

You MUST STOP IMMEDIATELY after spawning the Code Reviewer to preserve context:
- Record what was spawned in state file
- Save TODOs and commit state changes
- EXIT with clear continuation instructions
- This prevents agent responses from overflowing context

## State Context

**Purpose:**
Spawn Code Reviewer to create a comprehensive merge plan for integrating all phase branches into the project integration branch.

## Primary Actions

1. **Spawn Code Reviewer** with PROJECT_MERGE_PLANNING directive
2. **Provide Context**:
   - List of all phase integration branches
   - Target project integration branch
   - Dependency order requirements (R270)
3. **Update State File** with spawned agent details
4. **STOP per R313** - Exit immediately after spawning

## Valid State Transitions

- **SUCCESS** → WAITING_FOR_PROJECT_MERGE_PLAN (agent spawned)
- **FAILURE** → ERROR_RECOVERY (spawn failed)