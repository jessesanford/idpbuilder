# 🔴🔴🔴 RULE R281: MANDATORY COMPLETE STATE FILE INITIALIZATION

## Criticality: SUPREME LAW #7 - BLOCKING
**Penalty**: -100% for incomplete initial state files

## Description
The orchestrator MUST create ALL 4 COMPLETE state files during INIT:

### SF 3.0: 4-File State Initialization
Software Factory 3.0 requires creating **4 separate state files** during initialization:

1. **orchestrator-state-v3.json** - Complete project state including:
   - ALL phases from the implementation plan
   - ALL waves from the implementation plan
   - ALL efforts from the implementation plan
   - State machine tracking
   - Project progression data

2. **bug-tracking.json** - Bug tracking system (initially empty):
   - Empty bugs array
   - Metadata and statistics initialized to zero
   - Ready to receive bugs from code reviews

3. **integration-containers.json** - Iteration container system (initially empty):
   - Empty active_integrations array
   - Ready for wave/phase/project integration tracking
   - Convergence metrics initialized

4. **fix-cascade-state.json** - Fix cascade tracking (created on demand):
   - Only created when cross-container bugs detected
   - Not required during initial INIT state
   - Will be created atomically when first cascade triggered

## Requirements

### 1. MANDATORY COMPLETENESS
The initial state file MUST contain:
```yaml
current_phase: 1
current_wave: 1
current_state: PLANNING
phases:
  - phase_id: 1
    phase_name: "[FROM PLAN]"
    status: pending
    waves:
      - wave_id: 1
        wave_name: "[FROM PLAN]"
        status: pending
        efforts:
          - effort_id: 1
            effort_name: "[FROM PLAN]"
            assigned_agent: null
            branch_name: null
            status: pending
            metrics:
              planned_lines: 0
              actual_lines: 0
              review_cycles: 0
          # ... ALL efforts from plan
      # ... ALL waves from plan
  # ... ALL phases from plan
```

### 2. VALIDATION REQUIREMENTS
Before transitioning from INIT:
- Count phases in implementation plan
- Count waves per phase in plan
- Count efforts per wave in plan
- Verify state file has EXACT same counts
- FAIL if any mismatch

### 3. STRUCTURE REQUIREMENTS
Every phase MUST have:
- phase_id (sequential)
- phase_name (from plan)
- status: "pending"
- waves: [] (populated with ALL waves)

Every wave MUST have:
- wave_id (sequential within phase)
- wave_name (from plan)
- status: "pending"
- efforts: [] (populated with ALL efforts)

Every effort MUST have:
- effort_id (sequential within wave)
- effort_name (from plan)
- assigned_agent: null
- branch_name: null
- status: "pending"
- metrics: (with zeros)

### 4. VERIFICATION SCRIPT (SF 3.0)
Use validation tools to verify all 4 state files:

```bash
# Validate orchestrator-state-v3.json structure
bash tools/validate-state-file.sh orchestrator-state-v3.json

# Validate bug-tracking.json structure
bash tools/validate-state-file.sh bug-tracking.json

# Validate integration-containers.json structure
bash tools/validate-state-file.sh integration-containers.json

# Verify completeness against implementation plan
bash utilities/validate-state-completeness.sh \
  PROJECT-IMPLEMENTATION-PLAN.md \
  orchestrator-state-v3.json
```

### 5. INITIALIZATION TEMPLATES
Use the provided example files as templates:
- `orchestrator-state-v3.json.example` - Copy and customize
- `bug-tracking.json.example` - Copy as-is (starts empty)
- `integration-containers.json.example` - Copy as-is (starts empty)
- Schemas in `schemas/` directory define required structure

## Violations (SF 3.0)
- Missing ANY of the 4 required state files: -100% IMMEDIATE FAILURE
- Missing any phase from plan in orchestrator-state-v3.json: -100% IMMEDIATE FAILURE
- Missing any wave from plan in orchestrator-state-v3.json: -100% IMMEDIATE FAILURE
- Missing any effort from plan in orchestrator-state-v3.json: -100% IMMEDIATE FAILURE
- Invalid JSON in any state file: -100% IMMEDIATE FAILURE
- Incomplete placeholder structure: -50%
- Wrong initial values: -25%

## Implementation Notes (SF 3.0)
1. Parse implementation plan completely
2. Extract ALL phases, waves, efforts
3. Build orchestrator-state-v3.json with complete project structure
4. Create bug-tracking.json (empty bugs array, zero statistics)
5. Create integration-containers.json (empty active_integrations array)
6. Validate ALL files against schemas
7. Use atomic-state-update.sh to commit all 4 files together
8. Validate completeness before proceeding

## State Manager Coordination (SF 3.0)

State Manager handles initialization of all 4 SF 3.0 state files during project setup:

1. **orchestrator-state-v3.json** - State machine tracking + project progression
2. **bug-tracking.json** - Bugs array for all discovered issues
3. **integration-containers.json** - Iteration tracking for wave/phase/project convergence
4. **fix-cascade-state.json** - Created on-demand during cross-container bug cascades

Initialization protocol:
- Uses `tools/atomic-state-update.sh` for atomic creation and validation
- Validates against schemas in `schemas/*.schema.json`
- Commits all 4 files together (or none if any fail validation)
- Installs pre-commit hooks (R506) to enforce validation on all future commits

See: `.claude/agents/state-manager.md`, `schemas/orchestrator-state-v3.schema.json`

## Related Rules
- R288: State file updates (ongoing) - State Manager implements this
- R206: State machine compliance
- R234: State traversal requirements

## Enforcement (SF 3.0)
- Orchestrator INIT state MUST create all 4 state files
- No transitions allowed until ALL files validated
- All validation scripts MUST pass
- Manual verification recommended
- All 4 files MUST be committed together using atomic-state-update.sh

## Example Validation Output (SF 3.0)
```
✅ orchestrator-state-v3.json: Valid JSON, schema compliant
✅ bug-tracking.json: Valid JSON, schema compliant
✅ integration-containers.json: Valid JSON, schema compliant
✅ Plan has 3 phases - State file has 3 phases
✅ Phase 1 has 2 waves - State file matches
✅ Wave 1 has 5 efforts - State file matches
✅ All structure validated
✅ All 4 state files created and committed atomically
VALIDATION PASSED
```

---
*Rule R281 - Software Factory 3.0*
*Criticality: SUPREME LAW #7*
*Enforcement: IMMEDIATE*
*4-File Requirement: orchestrator-state-v3.json, bug-tracking.json, integration-containers.json, fix-cascade-state.json (on-demand)*