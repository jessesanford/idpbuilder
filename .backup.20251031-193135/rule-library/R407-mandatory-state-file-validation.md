# 🚨🚨🚨 RULE R407 - Mandatory State File Validation [BLOCKING]

**Category:** State Management / System Integrity
**Impact:** BLOCKING - Invalid state = immediate halt
**Priority:** SUPREME
**Created:** 2025-01-23
**Updated:** 2025-01-23

## Summary
The orchestrator-state-v3.json file is the single source of truth for the entire Software Factory 2.0 system. ALL agents MUST validate this file at critical points to prevent corruption, ensure consistency, and maintain system integrity. Any validation failure must halt the agent immediately.

## The Problem This Rule Solves
- State file corruption can break the entire orchestration pipeline
- Malformed JSON can cause agent crashes and data loss
- Invalid state transitions can violate the state machine
- Missing required fields can cause undefined behavior
- Timestamp inconsistencies can break synchronization
- Branch name mismatches can corrupt the git repository

## Required Implementation

### 1. MANDATORY Validation Points
Agents MUST validate orchestrator-state-v3.json at these points:
```bash
# BEFORE reading state (detect existing issues)
$CLAUDE_PROJECT_DIR/tools/enforce-state-validation.sh orchestrator-state-v3.json || {
    echo "❌ CRITICAL: State file invalid before read!"
    save_todos "STATE_INVALID_BEFORE_READ"
    exit 1
}

# BEFORE any modification
$CLAUDE_PROJECT_DIR/tools/enforce-state-validation.sh orchestrator-state-v3.json || {
    echo "❌ CRITICAL: Cannot modify invalid state file!"
    save_todos "STATE_INVALID_BEFORE_MODIFY"
    exit 1
}

# AFTER any modification
$CLAUDE_PROJECT_DIR/tools/enforce-state-validation.sh orchestrator-state-v3.json || {
    echo "❌ CRITICAL: State file invalid after modification!"
    save_todos "STATE_INVALID_AFTER_MODIFY"
    exit 1
}

# BEFORE state transitions
$CLAUDE_PROJECT_DIR/tools/enforce-state-validation.sh orchestrator-state-v3.json || {
    echo "❌ CRITICAL: Cannot transition from invalid state!"
    save_todos "STATE_INVALID_BEFORE_TRANSITION"
    exit 1
}

# AFTER state transitions
$CLAUDE_PROJECT_DIR/tools/enforce-state-validation.sh orchestrator-state-v3.json || {
    echo "❌ CRITICAL: State transition resulted in invalid state!"
    save_todos "STATE_INVALID_AFTER_TRANSITION"
    exit 1
}

# BEFORE spawning new agents
$CLAUDE_PROJECT_DIR/tools/enforce-state-validation.sh orchestrator-state-v3.json || {
    echo "❌ CRITICAL: Cannot spawn agents with invalid state!"
    save_todos "STATE_INVALID_BEFORE_SPAWN"
    exit 1
}

# AFTER completing any effort/wave/phase
$CLAUDE_PROJECT_DIR/tools/enforce-state-validation.sh orchestrator-state-v3.json || {
    echo "❌ CRITICAL: Completion resulted in invalid state!"
    save_todos "STATE_INVALID_AFTER_COMPLETION"
    exit 1
}
```

### 2. Validation Requirements
The validation script MUST check:

#### JSON Syntax
- Valid JSON structure
- No trailing commas
- Proper quote usage
- Correct bracket/brace matching

#### Required Fields
```json
{
    "current_state": "string (required)",
    "current_phase": "integer (required)",
    "current_wave": "integer (required)",
    "phases": "array (required)",
    "waves": "array (required)",
    "efforts_completed": "array (required)",
    "efforts_in_progress": "array (required)",
    "integration_branches": "object (required)",
    "last_updated": "ISO 8601 timestamp (required)",
    "version": "string (required)",
    "metadata": {
        "created_at": "ISO 8601 timestamp (required)",
        "created_by": "string (required)",
        "last_validated": "ISO 8601 timestamp (required)",
        "schema_version": "string (required)"
    }
}
```

#### State Machine Validity
- current_state MUST exist in software-factory-3.0-state-machine.json
- State transitions MUST be valid according to state machine
- No orphaned or undefined states

#### Effort Metadata
Each effort must have:
- Unique ID
- Valid status (pending/in_progress/completed/failed)
- Valid timestamps
- Proper branch naming
- Line count if completed

#### Branch Naming Convention
- Must follow pattern: `phase-X-wave-Y-effort-Z` or variants
- No spaces or special characters except hyphens
- Must match git branch name rules

#### Timestamp Consistency
- All timestamps must be valid ISO 8601
- last_updated must be most recent
- Timestamps must be in chronological order
- No future timestamps

### 3. Recovery Protocol
When validation fails:

```bash
# 1. STOP immediately - do not attempt automatic fixes
echo "🚨🚨🚨 STATE FILE VALIDATION FAILED"

# 2. Report validation errors in detail
$CLAUDE_PROJECT_DIR/tools/enforce-state-validation.sh --verbose orchestrator-state-v3.json

# 3. Save current TODOs before stopping
save_todos "STATE_VALIDATION_FAILURE"

# 4. Request orchestrator intervention
echo "❌ Requesting orchestrator intervention for state repair"
echo "DO NOT attempt manual fixes - could make corruption worse"

# 5. Exit with specific error code
exit 127  # State validation failure
```

### 4. Backup Requirements
Before ANY state modification:
```bash
# Create timestamped backup
cp orchestrator-state-v3.json "orchestrator-state.backup.$(date +%Y%m%d-%H%M%S).json"

# Validate backup was created
if [ ! -f "orchestrator-state.backup.*.json" ]; then
    echo "❌ Failed to create state backup"
    exit 1
fi
```

## Grading Penalties

| Violation | Penalty |
|-----------|---------|
| Missing validation before changes | -25% |
| Committing invalid state | -50% |
| Corrupting state file | -100% (IMMEDIATE FAIL) |
| Not halting on validation failure | -50% |
| Attempting automatic fixes | -25% |
| Missing backup before modification | -15% |

## Integration with Other Rules

### R206 (State Machine Validation)
R407 extends R206 by adding mandatory validation points and recovery procedures.

### R287 (TODO Persistence)
TODOs must be saved BEFORE any state validation failure exit.

### R203 (State-Aware Startup)
State validation is now part of mandatory startup sequence.

### R383 (Metadata Requirements)
Metadata fields are validated as part of state validation.

## Example Implementation

### In Orchestrator Agent
```bash
# Function to validate state with proper error handling
validate_state_file() {
    local context="$1"
    local state_file="${2:-orchestrator-state-v3.json}"

    echo "🔍 Validating state file: $context"

    if ! $CLAUDE_PROJECT_DIR/tools/enforce-state-validation.sh "$state_file"; then
        echo "❌❌❌ STATE VALIDATION FAILED: $context"
        echo "Validation errors:"
        $CLAUDE_PROJECT_DIR/tools/enforce-state-validation.sh --verbose "$state_file" 2>&1

        # Save TODOs before exit
        save_todos "STATE_INVALID_${context}"

        # Never try to fix automatically
        echo "⚠️ Manual intervention required - do not attempt automatic fixes"
        exit 127
    fi

    echo "✅ State validation passed: $context"
    return 0
}

# Use throughout agent lifecycle
validate_state_file "BEFORE_READ"
# ... read state ...

validate_state_file "BEFORE_MODIFY"
# ... modify state ...
validate_state_file "AFTER_MODIFY"

validate_state_file "BEFORE_TRANSITION"
# ... transition state ...
validate_state_file "AFTER_TRANSITION"
```

### In Software Engineer Agent
```bash
# Before reading effort configuration
validate_state_file "BEFORE_EFFORT_READ"

# After completing implementation
validate_state_file "AFTER_IMPLEMENTATION"
```

### In Code Reviewer Agent
```bash
# Before reading review targets
validate_state_file "BEFORE_REVIEW_READ"

# After completing review
validate_state_file "AFTER_REVIEW_COMPLETE"
```

## Common Validation Failures

### 1. Malformed JSON
```json
{
    "current_state": "PLANNING",
    "current_phase": 1,  // Trailing comma
}
```
**Fix:** Remove trailing comma

### 2. Invalid State
```json
{
    "current_state": "INVALID_STATE_NAME"
}
```
**Fix:** Use valid state from state machine

### 3. Missing Required Field
```json
{
    "current_state": "PLANNING"
    // Missing current_phase, current_wave, etc.
}
```
**Fix:** Add all required fields

### 4. Timestamp Issues
```json
{
    "last_updated": "2025-13-45"  // Invalid date
}
```
**Fix:** Use valid ISO 8601 timestamp

## Testing Validation

Test the validation with known-bad files:
```bash
# Test with invalid JSON
echo '{invalid json}' > test-state.json
$CLAUDE_PROJECT_DIR/tools/enforce-state-validation.sh test-state.json || echo "Correctly rejected invalid JSON"

# Test with missing fields
echo '{"current_state": "PLANNING"}' > test-state.json
$CLAUDE_PROJECT_DIR/tools/enforce-state-validation.sh test-state.json || echo "Correctly rejected missing fields"

# Test with invalid state
echo '{"current_state": "FAKE_STATE", "current_phase": 1, "current_wave": 1}' > test-state.json
$CLAUDE_PROJECT_DIR/tools/enforce-state-validation.sh test-state.json || echo "Correctly rejected invalid state"
```

## Related Documentation
- [R206: State Machine Validation](R206-state-machine-validation.md)
- [R287: TODO Persistence](R287-todo-persistence-comprehensive.md)
- [R203: State-Aware Startup](R203-state-aware-agent-startup.md)
- [R383: Metadata Requirements](R383-metadata-file-timestamp-requirements.md)
- [software-factory-3.0-state-machine.json](../state-machines/software-factory-3.0-state-machine.json)

## Enforcement Checklist
- [ ] Validation script created at tools/enforce-state-validation.sh
- [ ] All agents updated to call validation
- [ ] Backup mechanism implemented
- [ ] Recovery procedures documented
- [ ] Test cases verified
- [ ] Integration with R206, R287, R203, R383 confirmed
- [ ] Grading penalties documented
- [ ] Error codes standardized (127 for state validation failure)

## CRITICAL REMINDER
**NEVER ATTEMPT AUTOMATIC FIXES FOR STATE VALIDATION FAILURES**
Manual intervention by orchestrator is ALWAYS required to prevent cascading corruption.
## State Manager Coordination (SF 3.0)

State Manager enforces mandatory validation through:
- **Pre-commit hook** (R506): Validates on every commit attempt
- **Shutdown consultation**: Validates before atomic update
- **Startup consultation**: Validates consistency on read

Triple validation ensures corrupt state never persists:
1. Shutdown validation (before commit)
2. Pre-commit hook validation (during commit)
3. Startup validation (after checkout/pull)

See: `.git/hooks/pre-commit`, `tools/validate-state-file.sh`, R506
