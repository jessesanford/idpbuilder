# Orchestrator - INIT_ERROR_RECOVERY State Rules

## Purpose
Handle initialization failures and guide recovery to successful completion.

## Entry Criteria
- Validation failed OR
- Error encountered during initialization OR
- Missing required information OR
- User requested correction

## Required Actions

### 1. Identify Error Category

**Load State File and Determine Issue**:
```bash
# Read init-state-${PROJECT_PREFIX}.json
# Check validation.issues_found array
# Determine error category
```

**Error Categories**:
1. **MISSING_INFO**: Required fields empty/missing
2. **INVALID_CONFIG**: Configuration syntax/format errors
3. **REPO_FAILURE**: Git operations failed
4. **FILE_ERROR**: Cannot write/read files
5. **NETWORK_ERROR**: Cannot reach remote repositories
6. **VALIDATION_FAIL**: Content doesn't meet requirements

### 2. Determine Recovery Strategy

#### For MISSING_INFO
```
Recovery Path: Return to INIT_REQUIREMENTS_GATHERING
Actions:
1. Identify specific missing fields
2. Prepare targeted questions
3. Spawn architect with specific focus
4. Document what fields need collection
```

#### For INVALID_CONFIG
```
Recovery Path: Return to INIT_GENERATE_CONFIGS
Actions:
1. Identify which config has issues
2. Parse error messages
3. Spawn code reviewer to regenerate
4. Validate YAML syntax
```

#### For REPO_FAILURE
```
Recovery Path: Return to INIT_REPO_DECISION
Actions:
1. Determine if URL is invalid
2. Check network connectivity
3. Verify credentials if needed
4. Re-attempt repository setup
```

#### For FILE_ERROR
```
Recovery Path: Manual intervention needed
Actions:
1. Check file permissions
2. Verify directory exists
3. Request user to create directories
4. Retry file operations
```

#### For NETWORK_ERROR
```
Recovery Path: Retry with backoff
Actions:
1. Wait 5 seconds
2. Retry operation
3. If fails 3x, request manual intervention
4. Allow offline mode if possible
```

#### For VALIDATION_FAIL
```
Recovery Path: Depends on specific failure
Actions:
1. Parse validation report
2. Determine which component failed
3. Route to appropriate state for fixes
4. Document required corrections
```

### 3. Implement Recovery

#### Interactive Recovery
```markdown
## Initialization Error Detected

**Issue**: [Specific error description]
**Category**: [Error category]

**Required Action**: [What user needs to do]

Would you like me to:
1. Try to fix this automatically
2. Guide you through manual correction
3. Skip this step (if optional)
4. Abort initialization

Please respond with 1, 2, 3, or 4:
```

#### Automated Recovery
For recoverable errors:
1. Log the error
2. Apply fix automatically
3. Document what was fixed
4. Continue to next state

### 4. Update State File

Record recovery attempt:
```json
"error_recovery": {
  "error_type": "[category]",
  "error_details": "[specific issue]",
  "recovery_strategy": "[chosen approach]",
  "recovery_attempts": [count],
  "timestamp": "[ISO_TIME]"
}
```

### 5. Retry Limits

**Maximum Attempts**:
- Network errors: 3 retries with exponential backoff
- File errors: 2 retries
- Validation errors: No automatic retry
- User input errors: Unlimited (with user consent)

**After Max Retries**:
```markdown
## Initialization Failed

Despite multiple attempts, initialization could not complete.

**Final Error**: [Error details]
**Attempts Made**: [Number]

**Manual Resolution Required**:
1. [Step to fix issue]
2. [Step to fix issue]

Once resolved, run: /init-software-factory --resume

**Debug Information Saved**: init-error-[timestamp].log
```

## Common Recovery Scenarios

### Scenario 1: Invalid Git URL
**Error**: "fatal: repository 'invalid-url' not found"
**Recovery**:
1. Return to INIT_REQUIREMENTS_GATHERING
2. Ask: "The git URL appears invalid. Please provide the correct repository URL:"
3. Validate format before accepting
4. Update requirements and continue

### Scenario 2: Missing Technology Info
**Error**: "technology.primary_language is required but empty"
**Recovery**:
1. Return to INIT_REQUIREMENTS_GATHERING
2. Ask: "What programming language will this project use?"
3. Update requirements
4. Regenerate configs

### Scenario 3: Directory Permission Denied
**Error**: "Permission denied: cannot create directory"
**Recovery**:
1. Request user action:
   ```
   Please run: sudo chown -R $USER $CLAUDE_PROJECT_DIR
   Then press Enter to continue
   ```
2. Wait for confirmation
3. Retry operation

### Scenario 4: Network Timeout
**Error**: "Failed to connect to github.com"
**Recovery**:
1. Wait 5 seconds
2. Retry with longer timeout (30s)
3. If still fails, offer offline mode
4. Mark repository setup as "pending"

## Exit Criteria
- Error identified and documented
- Recovery strategy selected
- Either recovered OR clear manual steps provided
- State file updated with recovery information

## Transitions
**SUCCESSFUL RECOVERY**:
- → INIT_REQUIREMENTS_GATHERING (for missing info)
- → INIT_REPO_DECISION (for repo issues)
- → INIT_GENERATE_CONFIGS (for config issues)
- → INIT_VALIDATION (to revalidate after fixes)

**FAILED RECOVERY**:
- → TERMINAL_ERROR (with manual resolution steps)

## Best Practices
- Always explain the error in user-friendly terms
- Provide specific, actionable recovery steps
- Save state before attempting recovery
- Log all recovery attempts for debugging
- Offer manual override options

## Recovery State Persistence

Save recovery state for resume:
```json
{
  "recovery_checkpoint": {
    "last_successful_state": "[state]",
    "pending_operations": [list],
    "user_inputs_collected": {map},
    "can_resume": true/false
  }
}
```

This allows `/init-software-factory --resume` to continue from last good state.