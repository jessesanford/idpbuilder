# Software Factory 2.0 - State Machine Documentation

## Overview

This directory contains the state machine definition for the Software Factory 2.0 system.
The state machine controls the orchestrator's workflow through various states during project development.

## Primary File

- **software-factory-3.0-state-machine.json**: Complete state machine definition including:
  - State list
  - Transition matrix for each agent
  - Validation rules
  - Supreme laws and fundamental principles

## Project Completion Flow

### 🔴🔴🔴 CRITICAL: Multi-Phase Projects MUST Perform Project Integration (R283)

Multi-phase projects have a MANDATORY project-level integration phase before completion.
Direct transition from COMPLETE_PHASE to PROJECT_DONE is **PROHIBITED** for multi-phase projects.

### Multi-Phase Project Completion (MANDATORY Path)

```
COMPLETE_PHASE (final phase)
  ↓
PROJECT_INTEGRATE_WAVE_EFFORTS ← MANDATORY per R283
  ↓
SETUP_PROJECT_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE
  ↓
SPAWN_CODE_REVIEWER_PROJECT_MERGE_PLAN
  ↓
WAITING_FOR_PROJECT_MERGE_PLAN
  ↓
SPAWN_INTEGRATION_AGENT_PROJECT
  ↓
MONITORING_PROJECT_INTEGRATE_WAVE_EFFORTS
  ↓
PROJECT_REVIEW_WAVE_INTEGRATION
  ↓
WAITING_FOR_PROJECT_REVIEW_WAVE_INTEGRATION
  ↓ (validation/fix cycles as needed)
  ↓
PR_PLAN_CREATION
  ↓
PROJECT_DONE
```

**Penalty for skipping project integration:** -100% IMMEDIATE FAILURE

### Single-Phase Projects

#### Single-Phase, Single-Wave
For the simplest projects (one phase, one wave), may skip directly to PROJECT_DONE:
```
COMPLETE_PHASE → PROJECT_DONE
```

#### Single-Phase, Multi-Wave
Even single-phase projects with multiple waves should perform project integration
to create a comprehensive project-level demo:
```
COMPLETE_PHASE → PROJECT_INTEGRATE_WAVE_EFFORTS → ... → PROJECT_DONE
```

## Why Project Integration is Mandatory

1. **Verification**: Proves all phases actually work together as designed
2. **Integration Issues**: Catches phase interaction problems early
3. **Project Demo**: Creates comprehensive end-to-end demonstration (R291)
4. **Quality Gate**: Final validation before project completion
5. **The "Bow on the Project"**: Proper ending that demonstrates completeness

Without project integration, you have phases that individually work but **no proof**
the complete project functions as a unified system.

## R283 Compliance

Project integration is governed by:
- **R283**: Project Integration Protocol (BLOCKING - MANDATORY for completion)
- **R291**: Integration Demo Requirement (demos required at all integration levels)
- **R009**: Mandatory Wave/Phase Integration Protocol (general integration requirements)

See `rule-library/R283-project-integration-protocol.md` for complete details.

## Validation

To verify project integration requirement is properly enforced:
```bash
bash tools/validate-project-integration-required.sh
```

This script checks:
- State machine transitions are correct
- COMPLETE_PHASE state rules enforce R283
- PROJECT_INTEGRATE_WAVE_EFFORTS chain exists
- Only appropriate states can transition to PROJECT_DONE
- R283 rule documentation is complete

## State Transitions

### General Principles

1. **One Operation Per State**: Each state performs ONE atomic operation
2. **Mandatory Stops**: After state completion, agent MUST stop (R322)
3. **Continuation Commands**: Use `/continue-orchestrating` to proceed
4. **No Automatic Transitions**: State transitions are NEVER automatic

### Transition Matrix

The `transition_matrix` in the JSON defines valid transitions for each agent:
- **orchestrator**: Coordinates all workflows
- **sw-engineer**: Implements features
- **code-reviewer**: Reviews code and creates plans
- **architect**: Performs assessments

### Terminal States

- **PROJECT_DONE**: Project completed successfully
- **ERROR_RECOVERY**: Critical failure requiring manual intervention

Only these states have no outbound transitions.

## Common Patterns

### Spawning Pattern
```
SPAWN_SW_ENGINEERS → MONITORING_* → WAITING_FOR_*
```

### Review Pattern
```
SPAWN_CODE_REVIEWERS_* → MONITORING_REVIEWS → WAITING_FOR_*_PLANS
```

### Integration Pattern
```
*_INTEGRATE_WAVE_EFFORTS → SETUP_*_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE → SPAWN_INTEGRATION_AGENT_*
```

### Error Recovery
Most states can transition to `ERROR_RECOVERY` when problems are detected.

## References

- Rule R283: Project Integration Protocol
- Rule R291: Integration Demo Requirement
- Rule R322: Mandatory Stop Before State Transitions
- Rule R206: State Machine Validation Requirements

---

**Last Updated**: 2025-10-03
**Version**: 2.0
