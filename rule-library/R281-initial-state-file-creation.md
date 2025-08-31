# 🔴🔴🔴 RULE R281: MANDATORY COMPLETE STATE FILE INITIALIZATION

## Criticality: SUPREME LAW #7 - BLOCKING
**Penalty**: -100% for incomplete initial state files

## Description
The orchestrator MUST create a COMPLETE state file during INIT that includes:
- ALL phases from the implementation plan
- ALL waves from the implementation plan  
- ALL efforts from the implementation plan
- Proper placeholder structure for every single item

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

### 4. VERIFICATION SCRIPT
Use `/home/vscode/software-factory-template/utilities/validate-state-completeness.sh` to verify:
```bash
bash utilities/validate-state-completeness.sh \
  PROJECT-IMPLEMENTATION-PLAN.md \
  orchestrator-state.yaml
```

## Violations
- Missing any phase from plan: -100% IMMEDIATE FAILURE
- Missing any wave from plan: -100% IMMEDIATE FAILURE  
- Missing any effort from plan: -100% IMMEDIATE FAILURE
- Incomplete placeholder structure: -50%
- Wrong initial values: -25%

## Implementation Notes
1. Parse implementation plan completely
2. Extract ALL phases, waves, efforts
3. Build complete YAML structure
4. Write atomically
5. Validate before proceeding

## Related Rules
- R288: State file updates (ongoing)
- R206: State machine compliance
- R234: State traversal requirements

## Enforcement
- Orchestrator INIT state MUST complete this
- No transitions allowed until verified
- Validation script MUST pass
- Manual verification recommended

## Example Validation Output
```
✅ Plan has 3 phases - State file has 3 phases
✅ Phase 1 has 2 waves - State file matches
✅ Wave 1 has 5 efforts - State file matches
✅ All structure validated
VALIDATION PASSED
```

---
*Rule R281 - Software Factory 2.0*
*Criticality: SUPREME LAW #7*
*Enforcement: IMMEDIATE*