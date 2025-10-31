# WAITING_FOR_IMPLEMENTATION_PLAN State Rules

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

## 📋 PRIMARY DIRECTIVES

### Core Mandatory Rules:
1. R006, R287, R288, R510, R405

### State-Specific Rules:
6. **🚨🚨🚨 R502** - Implementation Plan Quality Gates
7. **🚨🚨🚨 R213** - Effort Metadata Requirements
8. **🚨🚨🚨 R233** - Active Monitoring Protocol

## 🔴🔴🔴 MANDATORY EXECUTION CHECKLIST

### BLOCKING REQUIREMENTS

- [ ] 1. Wait for Code Reviewer to complete wave implementation plan
  - Expected: `planning/phase{N}/wave{M}/WAVE-IMPLEMENTATION-PLAN.md` (per R550)
  - Lookup: `jq -r '.planning_files.phases.phase{N}.waves.wave{M}.implementation_plan' orchestrator-state-v3.json`
  - Polling: Every 30 seconds (R233)
  - Timeout: 30 minutes
  - **NOTE**: This state ALWAYS executes - no conditional skipping

- [ ] 2. Validate wave implementation file exists
  - File: Use path from `planning_files` field in state file
  - Fallback: `planning/phase{N}/wave{M}/WAVE-IMPLEMENTATION-PLAN.md`
  - Minimum size: >2000 bytes (detailed plan)

- [ ] 3. Validate wave implementation quality (R502 + R213)
  - Check: Contains detailed effort definitions
  - Check: Each effort has R213 metadata block
  - Check: File paths specified for each effort
  - Check: Dependencies documented
  - Check: Effort IDs follow convention
  - **BLOCKING**: All R213 metadata must be present

- [ ] 4. Extract effort list for state file
  - Parse: Effort definitions and metadata
  - Store: In orchestrator-state-v3.json for effort creation
  - Format: Array of efforts with R213 metadata
  - **BLOCKING**: Effort list needed for CREATE_NEXT_INFRASTRUCTURE

### STANDARD EXECUTION TASKS

- [ ] 5. Record validation results in state file
- [ ] 6. Display validation summary
  - Show: Effort count, R213 compliance
  - Show: Next state

### EXIT REQUIREMENTS

- [ ] 7. Set proposed next state: `CREATE_NEXT_INFRASTRUCTURE` (typically)
- [ ] 8. Spawn State Manager for SHUTDOWN_CONSULTATION
- [ ] 9. Save TODOs per R287
- [ ] 10. Set CONTINUE-SOFTWARE-FACTORY=TRUE
- [ ] 11. Stop execution (exit 0)

## State Purpose

Monitor Code Reviewer creating wave implementation plan, validate it contains EXACT specifications with R213 metadata, and extract effort list for infrastructure creation.

**Fidelity Enforced:** EXACT (detailed efforts, R213 metadata, file paths)

**CRITICAL CHANGE**: This state now ALWAYS executes - removed conditional "skip if plan exists" logic. Every wave MUST create implementation plan just-in-time.

## Validation Criteria

**Must Have (EXACT fidelity):**
- Detailed effort definitions (e.g., "Effort 1.1: Create User Model - Implement User class with validation")
- R213 metadata for EACH effort:
  ```yaml
  effort_id: "1.1"
  estimated_lines: 150
  dependencies: []
  files_touched: ["src/models/user.py"]
  branch_name: "effort/1.1-user-model"
  ```
- Complete file paths
- Dependency graphs
- Task breakdowns

**Must NOT Have:**
- Wave-level descriptions only
- Missing R213 metadata
- Missing file paths
- General guidance without specifics

## Entry Criteria

- **From**: SPAWN_CODE_REVIEWER_WAVE_IMPL
- **Required**: Code Reviewer spawned for wave implementation planning
- **NO SKIPPING**: This state always executes (removed conditional logic per #84)

## Exit Criteria

**Success** → CREATE_NEXT_INFRASTRUCTURE
**Failure** → ERROR_RECOVERY (if validation fails or timeout)

## Rules Enforced

- R510, R502, R213, R233, R288, R287, R405, R006

## Additional Context

**MIGRATION NOTE (Phase 0 Item #84):**
- **OLD BEHAVIOR**: Orchestrator could skip wave implementation planning if plan already existed
- **NEW BEHAVIOR**: Wave implementation planning is MANDATORY - always create plan just-in-time
- **RATIONALE**: Progressive planning requires reviewing previous work before planning next wave

SF 3.0 Progressive Planning - Final Validation:
1. Phase Architecture (pseudocode) ✅
2. Phase Implementation (wave list) ✅
3. Wave Architecture (real code) ✅
4. **Wave Implementation (exact + R213)** ← THIS STATE validates

This is the most detailed planning artifact:
- Complete effort specifications
- R213 metadata for tracking and automation
- File-level granularity
- Dependency management
- Branch naming for parallel work

Common Validation Failures:
1. Missing R213 metadata → BLOCKING, must be fixed
2. Wave-level only → Wrong fidelity, request revision
3. No file paths → Incomplete, must add
4. Malformed effort IDs → Fix format
5. Missing dependencies → Add if applicable
