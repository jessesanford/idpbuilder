# 🚨🚨🚨 RULE R540: STATE FILE SCHEMA COMPLIANCE 🚨🚨🚨

**BLOCKING CRITICALITY - SYSTEM INTEGRITY PROTECTION**

**Status**: Active
**Applies To**: State Manager, Orchestrator, All agents updating state files
**Criticality**: 🚨🚨🚨 **BLOCKING** - ABSOLUTE REQUIREMENT
**Penalty**: -100% + System corruption

---

## THE ABSOLUTE LAW

**ALL state file modifications MUST comply with JSON schema definitions.**

**NO state file updates without schema validation.**

**NO manual jq/yq commands that bypass validation.**

---

## 🔴🔴🔴 SUPREME DIRECTIVE 🔴🔴🔴

### JSON Schema is the SOURCE OF TRUTH

**Schema Location**: `/schemas/orchestrator-state-v3.schema.json`

**Every state file MUST**:
1. Validate against its JSON schema
2. Contain all REQUIRED fields
3. Use correct data types for all fields
4. Respect minimum/maximum constraints
5. Follow enum restrictions

**Schema violations = State corruption = System failure**

---

## CRITICAL SCHEMA REQUIREMENTS

### orchestrator-state-v3.json Schema Rules:

#### 1. `project_progression.current_phase` Requirements:

**REQUIRED FIELDS** (per schema lines 325-334):
```json
{
  "phase_number": <integer>,         // REQUIRED, minimum: 1
  "name": <string>,                  // REQUIRED
  "iteration": <integer>,            // REQUIRED, minimum: 0, default: 0
  "total_waves_in_phase": <integer>, // REQUIRED, minimum: 1
  "waves_completed": <integer>       // REQUIRED, minimum: 0, default: 0
}
```

**CRITICAL FIELD: `total_waves_in_phase`**
- **Type**: integer (not string, not array)
- **Constraint**: minimum 1
- **Purpose**: "CRITICAL for state machine guard evaluation at COMPLETE_WAVE"
- **Used By**: State machine guards to determine if all waves complete
- **Impact**: Missing this field = guard evaluation IMPOSSIBLE

**CRITICAL FIELD: `waves_completed`**
- **Type**: integer (not array, not string)
- **Constraint**: minimum 0
- **Purpose**: "Number of waves completed in this phase - used with total_waves_in_phase for guard evaluation"
- **Format**: COUNT of waves, not array of wave IDs
- **Impact**: Wrong type = guard evaluation FAILURE

#### 2. State Machine Requirements:

**REQUIRED FIELDS**:
```json
{
  "current_state": <string>,        // REQUIRED, must be valid state name
  "previous_state": <string>,       // REQUIRED
  "state_history": <array>,         // REQUIRED, must contain transition objects
  "last_transition_timestamp": <string> // REQUIRED, ISO 8601 format
}
```

#### 3. State History Entry Requirements:

**REQUIRED FIELDS** (per schema lines 43-53):
```json
{
  "from_state": <string>,     // REQUIRED
  "to_state": <string>,       // REQUIRED
  "timestamp": <string>,      // REQUIRED, ISO 8601 format
  "validated_by": <string>    // REQUIRED, MUST be "state-manager"
}
```

---

## COMMON SCHEMA VIOLATIONS

### Violation #1: Array Instead of Integer

**WRONG**:
```json
"waves_completed": ["wave_1_1", "wave_1_2"]  // Array of IDs
```

**CORRECT**:
```json
"waves_completed": 2  // Integer count
```

**Why This Matters**:
- Guards evaluate: `waves_completed == total_waves_in_phase`
- Array == integer comparison fails
- State machine cannot determine if all waves are complete
- System forced into ERROR_RECOVERY

### Violation #2: Missing Required Field

**WRONG**:
```json
"current_phase": {
  "phase_number": 1,
  "name": "Foundation",
  "iteration": 0,
  "waves_completed": 2
  // Missing: total_waves_in_phase
}
```

**CORRECT**:
```json
"current_phase": {
  "phase_number": 1,
  "name": "Foundation",
  "iteration": 0,
  "waves_completed": 2,
  "total_waves_in_phase": 3  // Required for guard evaluation
}
```

### Violation #3: Wrong Data Type

**WRONG**:
```json
"phase_number": "1",           // String instead of integer
"iteration": "0",              // String instead of integer
"waves_completed": "2"         // String instead of integer
```

**CORRECT**:
```json
"phase_number": 1,             // Integer
"iteration": 0,                // Integer
"waves_completed": 2           // Integer
```

### Violation #4: Constraint Violations

**WRONG**:
```json
"waves_completed": -1,         // Negative (minimum is 0)
"total_waves_in_phase": 0,     // Zero (minimum is 1)
"phase_number": 0              // Zero (minimum is 1)
```

**CORRECT**:
```json
"waves_completed": 0,          // Zero or positive
"total_waves_in_phase": 1,     // One or more
"phase_number": 1              // One or more
```

---

## VALIDATION REQUIREMENTS

### Before ANY State File Update:

**MANDATORY VALIDATION STEPS**:
1. Read current state file
2. Apply proposed changes
3. Validate against JSON schema
4. Check all REQUIRED fields present
5. Verify all data types correct
6. Check all constraints satisfied
7. Only then commit changes

**Use Validation Tool**:
```bash
# Validate before updating
bash tools/validate-state-schema.sh orchestrator-state-v3.json

# If validation fails, DO NOT PROCEED
# Fix violations first
```

### State Manager Requirements:

**State Manager MUST**:
1. Validate schema before ANY state file update
2. Reject updates that violate schema
3. Provide clear error messages for violations
4. Never commit invalid state data
5. Rollback on validation failure

### Orchestrator Requirements:

**Orchestrator MUST**:
1. NEVER update state files directly
2. ALWAYS use State Manager for state updates
3. Provide complete, valid data to State Manager
4. Respect schema requirements in all proposals

---

## ENFORCEMENT MECHANISMS

### 1. Pre-commit Hook Validation

**Pre-commit hook MUST**:
```bash
# Check if state files changed
if git diff --cached --name-only | grep -E '(orchestrator-state|bug-tracking|integration-containers|fix-cascade-state)\.json'; then
    # Validate changed state files
    bash tools/validate-state-schema.sh || {
        echo "ERROR: State file schema validation failed"
        echo "Fix schema violations before committing"
        exit 1
    }
fi
```

### 2. State Manager Startup Validation

**State Manager MUST validate on startup**:
```bash
# Startup sequence
1. Read orchestrator-state-v3.json
2. Validate against schema
3. If invalid:
   - Log violations
   - Create backup
   - Refuse to proceed
   - Require manual fix
```

### 3. Automated Testing

**Test suite MUST include**:
- Schema validation tests for all state files
- Type checking tests
- Required field presence tests
- Constraint validation tests
- Guard evaluation tests with various data

---

## SCHEMA MIGRATION PROTOCOL

### When Schema Changes:

**MANDATORY STEPS**:
1. Update JSON schema file
2. Create migration script
3. Test migration on copies
4. Migrate all existing state files
5. Validate all migrated files
6. Document changes
7. Update agent instructions
8. Test end-to-end

**Migration Script Template**:
```bash
#!/usr/bin/env bash
# migrate-schema-vX-to-vY.sh

# Read state file
STATE=$(cat orchestrator-state-v3.json)

# Apply migrations
STATE=$(echo "$STATE" | jq '.field = new_value')

# Validate against new schema
ajv validate -s schemas/orchestrator-state-v3.schema.json -d <(echo "$STATE")

# Only write if valid
if [[ $? -eq 0 ]]; then
    echo "$STATE" > orchestrator-state-v3.json
else
    echo "Migration failed validation"
    exit 1
fi
```

---

## GUARD EVALUATION DEPENDENCIES

### Why Schema Compliance is CRITICAL:

**State machine guards depend on schema compliance**:

```javascript
// Example guard: all_waves_complete
function evaluateGuard(state) {
    // REQUIRES these fields to be integers
    const waves_completed = state.project_progression.current_phase.waves_completed;
    const total_waves = state.project_progression.current_phase.total_waves_in_phase;

    // If either is undefined or wrong type, guard FAILS
    if (typeof waves_completed !== 'number' || typeof total_waves !== 'number') {
        throw new Error('Guard evaluation impossible - schema violation');
    }

    return waves_completed === total_waves;
}
```

**Schema violations make guard evaluation IMPOSSIBLE**:
- Undefined fields = cannot evaluate
- Wrong types = comparison fails
- Missing constraints = logic breaks
- Result: State machine STUCK

---

## VALIDATION TOOLING

### Available Tools:

1. **Schema Validator**: `tools/validate-state-schema.sh`
   - Validates single or all state files
   - Checks against JSON schema
   - Reports known violations
   - Exit code 0 = valid, 1 = invalid

2. **AJV CLI**: `ajv validate -s schema.json -d data.json`
   - Industry-standard JSON schema validator
   - Detailed error messages
   - Strict mode available

3. **JQ Type Checking**: `jq 'type' file.json`
   - Quick type verification
   - Useful for debugging

---

## ERROR MESSAGES AND RECOVERY

### When Validation Fails:

**State Manager Response**:
```
ERROR: Schema validation failed for orchestrator-state-v3.json

Violations detected:
- Field 'total_waves_in_phase' is required but missing
- Field 'waves_completed' has type 'array' but should be 'integer'

Cannot proceed with state transition.
Current state: START_PHASE_ITERATION
Proposed state: INTEGRATE_PHASE_WAVES
Decision: ERROR_RECOVERY (schema violation)

Action required:
1. Fix schema violations in state file
2. Run: bash tools/validate-state-schema.sh
3. Retry state transition
```

### Recovery Steps:

1. **Identify violations** using validator
2. **Backup current state** before fixing
3. **Fix violations** manually or with migration script
4. **Validate fix** using validation tool
5. **Retry operation** that failed
6. **Test end-to-end** to ensure system operational

---

## RELATED RULES

- **R517**: Universal State Manager Consultation Law (State Manager handles updates)
- **R288**: State File Update and Commit Protocol (Atomic updates)
- **R506**: Absolute Prohibition on Pre-commit Bypass (Validation in hooks)
- **R407**: Mandatory State File Validation (General validation requirement)

---

## CONSEQUENCES OF VIOLATION

### Immediate Effects:

- ❌ **Grade**: -100% automatic failure
- ❌ **System**: State machine guard evaluation fails
- ❌ **Progression**: System stuck in ERROR_RECOVERY
- ❌ **Integrity**: State data corruption
- ❌ **Recovery**: Manual intervention required

### Cascading Failures:

1. Guards cannot evaluate (undefined/wrong types)
2. State transitions blocked
3. System forced to ERROR_RECOVERY
4. Manual fixes required
5. Development halted
6. Potential data loss

---

## MANDATORY ACKNOWLEDGMENT

**Every agent that updates state files MUST acknowledge**:

```
I acknowledge R540: State File Schema Compliance
- I will NEVER update state files without schema validation
- I will ALWAYS use validation tooling before commits
- I will RESPECT all REQUIRED fields and constraints
- I understand schema violations = system corruption
- I will FOLLOW schema migration protocols for changes
```

---

## PERMANENT AND UNIVERSAL

**This rule is NOW and FOREVER**:
- ✅ Universal across ALL state file updates
- ✅ No exceptions for "quick fixes"
- ✅ Enforced by multiple mechanisms
- ✅ Validated by automated testing
- ✅ Required for system integrity

**Schema compliance is not optional. It is LAW.**

---

*Rule R540 - State File Schema Compliance*
*Created: 2025-10-31*
*Criticality: BLOCKING (🚨🚨🚨)*
*Enforcement: Multiple Layers*
*Penalty: -100% + System Corruption*
