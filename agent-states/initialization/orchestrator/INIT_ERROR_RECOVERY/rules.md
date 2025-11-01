# Orchestrator - INIT_ERROR_RECOVERY State Rules

## 🔴🔴🔴 ABSOLUTE REQUIREMENT: STATE MANAGER CONSULTATION 🔴🔴🔴

**SUPREME LAW - R517 - UNIVERSAL STATE MANAGER CONSULTATION LAW**

**BEFORE EXITING THIS STATE, YOU MUST:**

1. **SPAWN STATE MANAGER FOR SHUTDOWN_CONSULTATION** (MANDATORY - NO EXCEPTIONS)
2. **NEVER UPDATE STATE FILES DIRECTLY** (orchestrator-state-v3.json, bug-tracking.json, etc.)
3. **NEVER COMMIT STATE FILES YOURSELF** (State Manager does this atomically)

**FAILURE TO CONSULT STATE MANAGER = IMMEDIATE SYSTEM HALT (-100% GRADE)**

### Enforcement Mechanism

If you attempt to exit this state without spawning State Manager:
- ❌ Pre-commit hooks will REJECT your commit
- ❌ Validation tools will FAIL the build
- ❌ Grading system will assign -100% penalty
- ❌ System will transition to ERROR_RECOVERY

### Required Pattern (COPY THIS EXACTLY)

```bash
# At end of state work, BEFORE any state file updates:

echo "🔴 MANDATORY: Spawning State Manager for SHUTDOWN_CONSULTATION"

# Spawn State Manager (REQUIRED - NOT OPTIONAL)
# Task: state-manager
# State: SHUTDOWN_CONSULTATION
# Current State: [YOUR_CURRENT_STATE]
# Proposed Next State: [YOUR_PROPOSED_NEXT_STATE]
# Work Summary: [SUMMARY_OF_WORK_COMPLETED]

# State Manager will:
# 1. Validate proposed transition against state machine
# 2. Update all 4 state files atomically
# 3. Commit with [R288] tag
# 4. Return REQUIRED next state (may differ from proposal)

# Wait for State Manager response
# Follow State Manager's directive (REQUIRED next state)
# DO NOT proceed until State Manager confirms
```

**YOU MUST NEVER:**
- ❌ Update orchestrator-state-v3.json yourself
- ❌ Update bug-tracking.json yourself
- ❌ Update integration-containers.json yourself
- ❌ Use `jq` to modify state files
- ❌ Use `sed/awk` to modify state files
- ❌ Set `validated_by: "orchestrator"` (must be "state-manager")
- ❌ Skip State Manager consultation "just this once"
- ❌ Think "I'll validate it manually"

**ONLY State Manager may update state files. This is NON-NEGOTIABLE.**

See: `rule-library/R517-universal-state-manager-consultation-law.md`
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
**PROJECT_DONEFUL RECOVERY**:
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

## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**


### 🚨 CRITICAL DISTINCTION: AGENT STOPS ≠ FACTORY STOPS 🚨

**TWO INDEPENDENT DECISIONS - DO NOT CONFUSE THEM:**

#### 1. Should Agent Stop Work? (R322 Technical Requirement)
- Agent completes current state
- Agent saves TODOs and commits state
- Agent exits with `exit 0` (preserves context)
- User runs /continue-orchestrating to resume
- **This is NORMAL at checkpoints**

#### 2. Should Factory Continue? (R405 Operational Status)
- Even though agent stopped, can automation proceed?
- TRUE = Healthy completion, automation can continue
- FALSE = Catastrophic failure, must halt everything
- **R322 checkpoints = TRUE (99.9% of cases)**

### THE PATTERN AT R322 CHECKPOINTS

```bash
# 1. Complete state work
echo "✅ State work complete"

# 2. Update state file
jq '.state_machine.current_state = "NEXT_STATE"' orchestrator-state-v3.json > tmp && mv tmp orchestrator-state-v3.json

# 3. Save TODOs
save_todos "R322_CHECKPOINT"

# 4. Factory continues (operational status)
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"

# 5. Agent stops (technical requirement)
exit 0
```

**Both happen together! Agent stops AND factory continues!**

### WHEN TO USE EACH FLAG VALUE

**TRUE (99.9%):**
- ✅ R322 checkpoint reached
- ✅ State work completed successfully
- ✅ Ready for /continue-orchestrating
- ✅ Waiting for user to continue (NORMAL)
- ✅ Plan ready for review (agent done, factory proceeds)

**FALSE (0.1%):**
- ❌ CATASTROPHIC unrecoverable error
- ❌ Data corruption spreading
- ❌ Critical security violation
- ❌ NOT for R322 checkpoints
- ❌ NOT for user review needs
### 🚨 ERROR_RECOVERY - SEQUENTIAL FIX WORK PATTERN 🚨

**MOST COMMON VIOLATION: Setting FALSE after each fix!**

**Sequential Fix Work:**
```bash
# Fix Bug 1
fix_bug_1()
commit_per_r288()

# R322 checkpoint (Bug 2-5 still need fixing)
echo "🛑 R322: Checkpoint after Bug 1 fixed"
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # TRUE! Continue to Bug 2
exit 0  # Stop inference
```

**Why TRUE is correct:**
- System knows current state: ERROR_RECOVERY
- System knows remaining work: Additional fixes
- System knows next action: Continue fixing
- **NO HUMAN INTERVENTION NEEDED!**

**ONLY use FALSE if:**
- ❌ Can't determine what to fix (state corruption)
- ❌ Same bug fails 4+ times (pattern of failure)
- ❌ Fix cascade recursion detected
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

