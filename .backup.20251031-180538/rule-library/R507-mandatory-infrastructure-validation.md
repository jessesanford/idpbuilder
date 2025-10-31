# 🚨🚨🚨 RULE R507 - Mandatory Infrastructure Validation [BLOCKING]

## Category: Infrastructure Safety

## Priority: SUPREME

## Description
ALL infrastructure MUST be validated before ANY implementation work begins. This is a hard gate that prevents catastrophic failures from misconfigured infrastructure.

## Requirements

### MANDATORY VALIDATION POINTS
1. **Remote Repository Validation**
   - MUST verify remote matches target-repo-config.yaml
   - MUST check both 'origin' and 'target' remotes
   - ANY mismatch = IMMEDIATE ERROR_RECOVERY

2. **Branch Name Validation**
   - MUST match pre-planned names in orchestrator-state-v3.json
   - MUST be on correct branch before ANY work
   - Wrong branch = IMMEDIATE ERROR_RECOVERY

3. **Directory Path Validation**
   - MUST exist at pre-planned path
   - MUST have correct permissions
   - Missing directory = IMMEDIATE ERROR_RECOVERY

### ENFORCEMENT
```bash
# Validation is MANDATORY after infrastructure creation
CREATE_NEXT_INFRASTRUCTURE → VALIDATE_INFRASTRUCTURE (REQUIRED)
CREATE_NEXT_INFRASTRUCTURE → VALIDATE_INFRASTRUCTURE (REQUIRED)
```

### VALIDATION SCRIPT
```bash
# Use the official validation script
bash $CLAUDE_PROJECT_DIR/utilities/validate-infrastructure.sh || {
    echo "❌ INFRASTRUCTURE VALIDATION FAILED"
    transition_to ERROR_RECOVERY
}
```

## Penalties
- Skipping validation: -50% IMMEDIATE
- Proceeding with invalid infrastructure: -100% FAILURE
- Not transitioning to ERROR_RECOVERY on failure: -75%

## Examples

### CORRECT: Proper Validation
```bash
# After creating infrastructure
echo "✅ Infrastructure created, validating..."
bash $CLAUDE_PROJECT_DIR/utilities/validate-infrastructure.sh
if [ $? -eq 0 ]; then
    echo "✅ Infrastructure validated"
    transition_to SPAWN_SW_ENGINEERS
else
    echo "❌ Validation failed"
    transition_to ERROR_RECOVERY
fi
```

### WRONG: Skipping Validation
```bash
# ❌ NEVER DO THIS
create_infrastructure
transition_to SPAWN_SW_ENGINEERS  # Missing validation!
```

## Related Rules
- R508: Target Repository Enforcement
- R360: Just-in-Time Infrastructure
- R504: Pre-Infrastructure Planning Protocol

## Source
Created to prevent catastrophic infrastructure failures where code goes to wrong repositories.

## Metadata
- Created: 2025-09-27
- Criticality: BLOCKING
- Enforcement: AUTOMATIC via state machine