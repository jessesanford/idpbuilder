# 🚨🚨🚨 BLOCKING RULE R406: Orchestrator State Schema Validation

## Rule Definition
ALL orchestrator state files MUST pass JSON schema validation before being saved or committed. Invalid state files cause cascading failures across the entire Software Factory system.

## Enforcement Points

### 1. Pre-Save Validation (MANDATORY)
```bash
# Before ANY write to orchestrator-state-v3.json
tools/validate-state.sh orchestrator-state-v3.json || {
    echo "❌ BLOCKED: State file validation failed"
    exit 1
}
```

### 2. Pre-Commit Hook (MANDATORY)
```bash
# In .git/hooks/pre-commit
if git diff --cached --name-only | grep -q "orchestrator-state-v3.json"; then
    tools/validate-state.sh orchestrator-state-v3.json || {
        echo "❌ COMMIT BLOCKED: State file validation failed"
        exit 1
    }
fi
```

### 3. State Transition Validation (CRITICAL)
```bash
# Before ANY state transition
validate_state_transition() {
    local new_state="$1"
    local current_state=$(jq -r '.state_machine.current_state' orchestrator-state-v3.json)

    # Validate against state machine
    grep -q "$current_state.*->.*$new_state" software-factory-3.0-state-machine.json || {
        echo "❌ INVALID TRANSITION: $current_state -> $new_state"
        return 1
    }

    # Validate schema after update
    tools/validate-state.sh orchestrator-state-v3.json || return 1
}
```

## Schema Location
- **Primary Schema**: `orchestrator-state.schema.json`
- **Sub-State Schemas**:
  - `integration-state.schema.json`
  - `fix-cascade-state.schema.json`
  - `pr-ready-state.schema.json`

## Required Fields (NEVER OPTIONAL)
```json
{
    "current_phase": "integer ≥ 1",
    "current_wave": "integer ≥ 0",
    "current_state": "valid state from state machine",
    "previous_state": "valid state or null",
    "transition_time": "ISO 8601 timestamp",
    "phases_planned": "integer ≥ 1",
    "waves_per_phase": "array of integers",
    "efforts_completed": "array",
    "efforts_in_progress": "array",
    "efforts_pending": "array",
    "project_info": "object with name, description, start_date"
}
```

## Validation Tools

### 1. Schema Validator
```bash
# Primary validation tool
tools/validate-state.sh [state-file]

# Returns:
# 0 - Valid
# 1 - Invalid (with detailed error messages)
```

### 2. Python Validator (Comprehensive)
```python
# For detailed validation with jsonschema
python3 utilities/validate-schema-comprehensive.py
```

### 3. Quick Validation Function
```bash
validate_state() {
    local state_file="${1:-orchestrator-state-v3.json}"

    # Check file exists
    [ -f "$state_file" ] || {
        echo "❌ State file not found: $state_file"
        return 1
    }

    # Validate JSON syntax
    jq empty "$state_file" 2>/dev/null || {
        echo "❌ Invalid JSON syntax"
        return 1
    }

    # Validate against schema
    tools/validate-state.sh "$state_file"
}
```

## Auto-Fix Capabilities

### 1. Add Missing Required Fields
```bash
fix_missing_fields() {
    local state_file="orchestrator-state-v3.json"

    # Add missing required fields with defaults
    jq '. + {
        "project_info": (.project_info // {
            "name": "unnamed",
            "description": "No description",
            "start_date": (now | strftime("%Y-%m-%d"))
        }),
        "efforts_completed": (.efforts_completed // []),
        "efforts_in_progress": (.efforts_in_progress // []),
        "efforts_pending": (.efforts_pending // [])
    }' "$state_file" > "$state_file.tmp" && mv "$state_file.tmp" "$state_file"
}
```

### 2. Fix Invalid State Values
```bash
fix_invalid_state() {
    local current_state=$(jq -r '.state_machine.current_state' orchestrator-state-v3.json)

    # Check if state exists in state machine
    if ! grep -q "STATE: $current_state" software-factory-3.0-state-machine.json; then
        echo "⚠️  Invalid state '$current_state', resetting to INIT"
        jq '.state_machine.current_state = "INIT" | .state_machine.previous_state = null' \
            orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
    fi
}
```

## Integration with State Machine

### State Transition Validator
```bash
# Validate state transitions against state machine
validate_transition() {
    local from_state="$1"
    local to_state="$2"
    local state_machine="software-factory-3.0-state-machine.json"

    # Extract valid transitions
    grep -A 20 "^## STATE: $from_state" "$state_machine" | \
    grep "^### TRANSITIONS TO:" -A 10 | \
    grep -q "- $to_state" || {
        echo "❌ Invalid transition: $from_state -> $to_state"
        return 1
    }
}
```

## Common Validation Errors

### 1. Missing Required Fields
```
❌ State file validation failed!
  Field: (root)
  Error: 'project_info' is a required property
```
**Fix**: Add the missing field with appropriate default values

### 2. Invalid State Value
```
❌ State file validation failed!
  Field: current_state
  Error: 'INVALID_STATE' is not one of ['INIT', 'PLANNING', ...]
```
**Fix**: Update to a valid state from the state machine

### 3. Type Mismatch
```
❌ State file validation failed!
  Field: current_phase
  Error: 'one' is not of type 'integer'
```
**Fix**: Ensure field types match schema definitions

## Enforcement in Agents

### Orchestrator Agent
```bash
# At startup
validate_state || {
    echo "❌ FATAL: Invalid state file - cannot proceed"
    exit 1
}

# Before state transitions
update_state() {
    local new_state="$1"
    local temp_file="orchestrator-state-v3.json.tmp"

    # Update state in temp file
    jq ".state_machine.current_state = \"$new_state\"" orchestrator-state-v3.json > "$temp_file"

    # Validate before committing
    tools/validate-state.sh "$temp_file" || {
        rm "$temp_file"
        return 1
    }

    mv "$temp_file" orchestrator-state-v3.json
}
```

## Schema Evolution Protocol

### Adding New Fields
1. Add field to schema as optional first
2. Update all state files to include field
3. After migration period, make field required if needed

### Removing Fields
1. Mark field as deprecated in schema
2. Stop writing to field in new code
3. After migration period, remove from schema

### Schema Version Tracking
```json
{
    "$schema": "http://json-schema.org/draft-07/schema#",
    "version": "2.0.0",
    "last_updated": "2025-01-23",
    ...
}
```

## Monitoring and Alerts

### Validation Failures Tracking
```bash
# Log all validation failures
log_validation_failure() {
    local timestamp=$(date -Iseconds)
    local error="$1"

    echo "[$timestamp] VALIDATION FAILURE: $error" >> validation-errors.log

    # Alert if too many failures
    local failure_count=$(grep "$(date +%Y-%m-%d)" validation-errors.log | wc -l)
    if [ "$failure_count" -gt 10 ]; then
        echo "⚠️  HIGH VALIDATION FAILURE RATE: $failure_count today"
    fi
}
```

## Recovery from Invalid State

### Emergency State Reset
```bash
# When state file is corrupted beyond repair
reset_to_last_valid() {
    local backup_dir="state-backups"

    # Find last valid backup
    for backup in $(ls -t "$backup_dir"/orchestrator-state-*.json); do
        if tools/validate-state.sh "$backup" 2>/dev/null; then
            echo "✅ Restoring from: $backup"
            cp "$backup" orchestrator-state-v3.json
            return 0
        fi
    done

    echo "❌ No valid backup found - manual intervention required"
    return 1
}
```

## Penalties for Violations
- **Writing invalid state**: System halt, -50% reliability score
- **Committing invalid state**: Repository corruption, -75% score
- **Ignoring validation errors**: Cascade failures, -100% score

## Related Rules
- R206: State machine validation
- R287: TODO state persistence
- R337: Base branch tracking
- R340: Planning file locations
- R344: Metadata location tracking

---

**REMEMBER**: The orchestrator state file is the MEMORY of the Software Factory. Invalid state = lost context = project failure. ALWAYS validate before save!
## State Manager Coordination (SF 3.0)

**THIS RULE IS IMPLEMENTED BY STATE MANAGER**

State Manager performs schema validation during shutdown consultation:
1. Validates `orchestrator-state-v3.json` against `schemas/orchestrator-state-v3.schema.json`
2. Validates `bug-tracking.json` against `schemas/bug-tracking.schema.json`
3. Validates `integration-containers.json` against `schemas/integration-containers.schema.json`
4. Validates `fix-cascade-state.json` (if exists) against `schemas/fix-cascade-state.schema.json`

Uses `tools/validate-state-file.sh` for all validations. Rolls back ALL changes if ANY file fails validation (atomicity guarantee).

See: `tools/validate-state-file.sh`, `tools/atomic-state-update.sh`, R506 (pre-commit hook also validates)
