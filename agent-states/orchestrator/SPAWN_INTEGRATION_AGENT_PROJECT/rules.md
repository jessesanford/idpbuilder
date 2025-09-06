# Orchestrator - SPAWN_INTEGRATION_AGENT_PROJECT State Rules

## 🛑🛑🛑 R313 MANDATORY STOP AFTER SPAWNING AGENTS 🛑🛑🛑

**SUPREME LAW - VIOLATION = -100% IMMEDIATE FAILURE**

You MUST STOP IMMEDIATELY after spawning the Integration Agent.

## State Context

**Purpose:**
Spawn Integration Agent to execute the project merge plan, merging all phase integration branches into the project integration branch.

## Primary Actions

1. **Spawn Integration Agent** with PROJECT_INTEGRATION directive
2. **Provide Resources**:
   - PROJECT-MERGE-PLAN.md location
   - Project integration workspace path
   - Phase branch details
3. **Update State File** with agent details
4. **STOP per R313**

## Valid State Transitions

- **SUCCESS** → MONITORING_PROJECT_INTEGRATION (agent spawned)
- **FAILURE** → ERROR_RECOVERY (spawn failed)