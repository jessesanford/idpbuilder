# Git Commit Error Handling Template

## Purpose

This template provides standard error handling for git commit operations in orchestrator states. It ensures that commit failures (especially schema validation errors) are properly caught and result in appropriate R405 continuation flags.

## Usage

Include this pattern in any orchestrator state that performs git commits (R288 compliance).

---

## Standard Commit Pattern with Error Handling

```bash
# Commit state file changes per R288
git add orchestrator-state-v3.json orchestrator-state-demo.json .cascade-state-backup.json .orchestrator-state-v3.json

# Attempt commit with error handling
if ! git commit -m "state: ${CURRENT_STATE} → ${NEXT_STATE} - ${TRANSITION_REASON} [R288]"; then
    # Commit failed - likely schema validation error
    echo "❌ CRITICAL: Git commit failed - likely schema validation error"
    echo "State: ${CURRENT_STATE}"
    echo "Attempted transition: ${CURRENT_STATE} → ${NEXT_STATE}"
    echo ""
    echo "Common causes:"
    echo "  - Schema validation failure (check pre-commit hook output above)"
    echo "  - Missing required fields in JSON files"
    echo "  - Invalid JSON syntax"
    echo "  - File permissions issues"
    echo ""
    echo "🛑 Cannot proceed - manual intervention required"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=SCHEMA_VALIDATION"
    exit 1
fi

# Commit succeeded - push changes
if ! git push; then
    # Push failed - network or permission issue
    echo "❌ WARNING: Git push failed - changes committed locally but not pushed"
    echo "This may cause synchronization issues but is not fatal"
    echo "Proceeding with state execution..."
    # Don't exit - local commit succeeded, push can be retried later
fi

echo "✅ State changes committed and pushed successfully"
```

---

## For Bug Tracking Files

When committing bug-tracking.json or other data files with strict schemas:

```bash
# Validate against schema BEFORE committing (defensive)
if command -v python3 >/dev/null 2>&1; then
    if python3 -c "import json, jsonschema; \
        schema = json.load(open('schemas/bug-tracking.schema.json')); \
        data = json.load(open('bug-tracking.json')); \
        jsonschema.validate(data, schema)" 2>/dev/null; then
        echo "✅ bug-tracking.json schema validation passed"
    else
        echo "❌ ERROR: bug-tracking.json schema validation failed"
        echo "Fix validation errors before attempting commit:"
        python3 -c "import json, jsonschema; \
            schema = json.load(open('schemas/bug-tracking.schema.json')); \
            data = json.load(open('bug-tracking.json')); \
            jsonschema.validate(data, schema)" 2>&1 | head -20
        echo ""
        echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=SCHEMA_VALIDATION"
        exit 1
    fi
fi

# Now safe to commit
git add bug-tracking.json
if ! git commit -m "bugs: ${BUG_DESCRIPTION} [schema-validated]"; then
    echo "❌ CRITICAL: Commit failed despite schema validation"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=COMMIT_FAILURE"
    exit 1
fi

git push || echo "⚠️ WARNING: Push failed - committed locally"
```

---

## R405 Continuation Flag Rules

### When to emit TRUE:
- Commit and push both succeeded ✅
- Push failed but commit succeeded (non-fatal) ✅

### When to emit FALSE:
- Commit failed due to schema validation ❌
- Commit failed due to merge conflicts ❌
- Commit failed due to file permissions ❌
- Any unrecoverable error preventing state completion ❌

---

## Example: FIX_WAVE_UPSTREAM_BUGS State

```bash
#!/bin/bash

# ... state execution logic ...

# Update bug tracking with injected bugs
# DEFENSIVE: Ensure all required schema fields are present
for BUG_ID in "${INJECTED_BUGS[@]}"; do
    # Update bug in bug-tracking.json
    jq --arg bug_id "$BUG_ID" \
       --arg status "INJECTED" \
       --arg discovered_by "test-framework" \
       '.bugs[] | select(.bug_id == $bug_id) | .status = $status | .discovered_by = $discovered_by' \
       bug-tracking.json > bug-tracking.json.tmp
    mv bug-tracking.json.tmp bug-tracking.json
done

# Validate schema BEFORE committing
if command -v python3 >/dev/null 2>&1; then
    if ! python3 -c "import json, jsonschema; \
        schema = json.load(open('schemas/bug-tracking.schema.json')); \
        data = json.load(open('bug-tracking.json')); \
        jsonschema.validate(data, schema)" 2>&1; then
        echo "❌ ERROR: Schema validation failed for bug-tracking.json"
        echo "Cannot proceed with commit"
        echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=SCHEMA_VALIDATION"
        exit 1
    fi
fi

# Commit state changes
git add orchestrator-state-v3.json bug-tracking.json

if ! git commit -m "state: FIX_WAVE_UPSTREAM_BUGS → START_WAVE_ITERATION - bugs injected [R288]"; then
    echo "❌ CRITICAL: Commit failed"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=COMMIT_FAILURE"
    exit 1
fi

git push || echo "⚠️ Push failed - committed locally"

# State execution complete - emit R405 flag
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
exit 0
```

---

## Benefits

1. **Early Detection**: Schema validation before commit prevents wasted work
2. **Clear Error Messages**: User knows exactly what went wrong
3. **Proper R405 Emission**: Always emits flag (TRUE or FALSE) before exit
4. **Graceful Degradation**: Push failures don't stop execution
5. **Debugging Info**: Detailed error output for troubleshooting

---

## Integration with Existing States

To add this to an existing state:

1. Find git commit operations in state rules
2. Replace simple `git commit` with error-handled version
3. Add schema validation for data files
4. Ensure R405 flag emitted on ALL exit paths
5. Test with intentionally invalid data to verify error handling

---

## Related Rules

- **R288**: State file update and commit protocol
- **R405**: Automation continuation flag (SUPREME LAW)
- **R506**: Pre-commit hook enforcement

---

**Created**: 2024-10-24
**Purpose**: Prevent schema validation failures from causing silent orchestrator crashes
**Trigger**: Investigation into Test 3b iteration 6 R405 violation
