# Orchestrator - SPAWN_CODE_REVIEWER_PROJECT_VALIDATION State Rules

## 🛑🛑🛑 R313 MANDATORY STOP AFTER SPAWNING AGENTS 🛑🛑🛑

**SUPREME LAW - VIOLATION = -100% IMMEDIATE FAILURE**

You MUST STOP IMMEDIATELY after spawning the Code Reviewer.

## State Context

**Purpose:**
Spawn Code Reviewer to perform comprehensive validation of the integrated project, ensuring all phases work together correctly.

## Primary Actions

1. **Spawn Code Reviewer** with PROJECT_VALIDATION directive
2. **Provide Context**:
   - Project integration branch location
   - List of all integrated phases
   - Validation requirements per R271-R275
3. **Request Validation**:
   - Functional correctness across phases
   - No inter-phase conflicts
   - Production readiness
4. **Update State File** and **STOP per R313**

## Valid State Transitions

- **SUCCESS** → WAITING_FOR_PROJECT_VALIDATION (reviewer spawned)
- **FAILURE** → ERROR_RECOVERY (spawn failed)