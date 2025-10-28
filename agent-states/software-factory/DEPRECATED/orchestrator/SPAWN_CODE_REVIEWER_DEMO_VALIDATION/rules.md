# Orchestrator State: SPAWN_CODE_REVIEWER_DEMO_VALIDATION

## Purpose
Spawn Code Reviewer to validate integration demos by RUNNING them and verifying
they execute successfully. This is the R291 GATE 4 enforcement mechanism.

## Entry Conditions
- Integration code review passed (APPROVED)
- Integration infrastructure exists
- Demo files should be present (created by integration agent)

## State Responsibilities

### 1. Verify Demo Infrastructure
Check that demos directory exists and demo scripts are present:

```bash
# Determine integration level and demo path
INTEGRATE_WAVE_EFFORTS_TYPE=$(jq -r '.current_integration_type' orchestrator-state-v3.json)
CURRENT_PHASE=$(jq -r '.current_phase' orchestrator-state-v3.json)
CURRENT_WAVE=$(jq -r '.current_wave' orchestrator-state-v3.json)

# Set demo paths based on level
if [ "$INTEGRATE_WAVE_EFFORTS_TYPE" = "wave" ]; then
    DEMO_DIR="demos/phase${CURRENT_PHASE}/wave${CURRENT_WAVE}/integration"
elif [ "$INTEGRATE_WAVE_EFFORTS_TYPE" = "phase" ]; then
    DEMO_DIR="demos/phase${CURRENT_PHASE}/integration"
elif [ "$INTEGRATE_WAVE_EFFORTS_TYPE" = "project" ]; then
    DEMO_DIR="demos/project"
else
    echo "🔴 Unknown integration type: $INTEGRATE_WAVE_EFFORTS_TYPE"
    exit 291
fi

# Verify demo infrastructure exists
if [ ! -d "$DEMO_DIR" ]; then
    echo "⚠️ WARNING: Demo directory missing: $DEMO_DIR"
    echo "This may indicate integration agent did not create demos per R291"
    echo "Code Reviewer will attempt to find and execute any demos present"
fi
```

### 2. Prepare Code Reviewer Spawn Instructions

Create clear instructions for the Code Reviewer to run demos:

```bash
# Create demo validation task file
TASK_FILE=".software-factory/demo-validation-task.json"

cat > "$TASK_FILE" << EOF
{
  "task_name": "Validate ${INTEGRATE_WAVE_EFFORTS_TYPE} integration demos",
  "task_description": "Run and verify integration demos execute successfully per R291 Gate 4",
  "integration_type": "${INTEGRATE_WAVE_EFFORTS_TYPE}",
  "demo_directory": "${DEMO_DIR}",
  "requirements": {
    "execute_all_demos": true,
    "capture_outputs": true,
    "verify_success": true,
    "create_report": true,
    "save_logs": true
  },
  "report_location": ".software-factory/phase${CURRENT_PHASE}/${INTEGRATE_WAVE_EFFORTS_TYPE}/demo-evaluation-report.md",
  "log_location": ".software-factory/phase${CURRENT_PHASE}/${INTEGRATE_WAVE_EFFORTS_TYPE}/demo-evaluation.log",
  "state_to_enter": "DEMO_VALIDATION",
  "spawned_at": "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
}
EOF

echo "✅ Created demo validation task file: $TASK_FILE"
```

### 3. Spawn Code Reviewer for Demo Validation

```bash
# Spawn Code Reviewer with demo validation instructions
echo "🚀 Spawning Code Reviewer for demo validation..."

# Code Reviewer will:
# 1. Navigate to demos directory
# 2. Execute all demo scripts
# 3. Capture outputs
# 4. Verify successful execution
# 5. Create demo-evaluation-report.md
# 6. Save logs to .software-factory/.../demo-evaluation.log

# SPAWN COMMAND (actual implementation depends on agent spawning mechanism)
spawn_code_reviewer \
    --state DEMO_VALIDATION \
    --task-file "$TASK_FILE" \
    --integration-type "$INTEGRATE_WAVE_EFFORTS_TYPE" \
    --demo-directory "$DEMO_DIR"
```

### 4. Update Orchestrator State

Record the code reviewer spawn and transition to waiting state:

```bash
# Update state to reflect demo validation in progress
jq ".state_machine.current_state = \"WAITING_FOR_DEMO_VALIDATION\"" -i orchestrator-state-v3.json
jq ".waiting_for = \"demo_validation\"" -i orchestrator-state-v3.json
jq ".demo_validation = {
  \"integration_type\": \"$INTEGRATE_WAVE_EFFORTS_TYPE\",
  \"demo_directory\": \"$DEMO_DIR\",
  \"validation_status\": \"in_progress\",
  \"code_reviewer_spawned_at\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
  \"evaluation_report\": \".software-factory/phase${CURRENT_PHASE}/${INTEGRATE_WAVE_EFFORTS_TYPE}/demo-evaluation-report.md\",
  \"evaluation_log\": \".software-factory/phase${CURRENT_PHASE}/${INTEGRATE_WAVE_EFFORTS_TYPE}/demo-evaluation.log\"
}" -i orchestrator-state-v3.json

# Commit state update
git add orchestrator-state-v3.json "$TASK_FILE"
git commit -m "state: spawn code reviewer for demo validation (R291 Gate 4)"
git push
```

## Exit Conditions

### ✅ CONTINUE-SOFTWARE-FACTORY=TRUE
Conditions for successful completion:
- Code Reviewer spawned successfully
- Demo infrastructure verified (or noted as missing)
- State updated to WAITING_FOR_DEMO_VALIDATION
- Task file created with clear instructions

**Next State**: WAITING_FOR_DEMO_VALIDATION

### 🔴 CONTINUE-SOFTWARE-FACTORY=FALSE
Conditions requiring manual intervention:
- Demo infrastructure completely missing (no recovery possible)
- Cannot spawn Code Reviewer (infrastructure broken)
- State file corruption detected
- Integration type unknown/invalid

**Next State**: ERROR_RECOVERY (manual intervention required)

## Next States
- **WAITING_FOR_DEMO_VALIDATION** (normal path)
- **ERROR_RECOVERY** (if infrastructure fundamentally broken)

## R291 Compliance

This state enforces R291 GATE 4 (Demo Verification). It CANNOT be skipped.

**Per R291 line 44:**
> "Marking integration complete without passing build/test/demo = IMMEDIATE DISQUALIFICATION"

**Penalty for skipping this state**: -100% (AUTOMATIC FAILURE)

**State machine enforcement**:
- Direct transitions from code review to completion are PROHIBITED
- Demo validation is a MANDATORY step
- Cannot bypass via state manipulation

## State Machine Reference

This state is part of the R291 enforcement chain:

```
WAITING_FOR_REVIEW_WAVE_INTEGRATION (approved)
  ↓ (REQUIRED)
SPAWN_CODE_REVIEWER_DEMO_VALIDATION  ← YOU ARE HERE
  ↓ (REQUIRED)
WAITING_FOR_DEMO_VALIDATION
  ↓ (based on results)
  ├─ Demos passed → WAVE/PHASE/PROJECT_COMPLETE
  └─ Demos failed → ERROR_RECOVERY (MANDATORY)
```

## Critical Rules Referenced
- **R291**: Integration Demo Requirement (BLOCKING)
- **R322**: Mandatory Orchestrator Checkpoints (SUPREME LAW)
- **R263**: Integration Documentation Requirements
- **R265**: Integration Testing Requirements

## Remember

**Demos are not optional.** They are PROOF that the integration actually works.

Per R291:
- Build + Tests + Demo MUST all pass
- No exceptions
- No shortcuts
- No "we'll do it later"

**This state ensures that promise is kept.**

CONTINUE-SOFTWARE-FACTORY=TRUE

## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

**MUST output as LAST line before state completion:**
- `CONTINUE-SOFTWARE-FACTORY=TRUE` - State completed, automation can continue
- `CONTINUE-SOFTWARE-FACTORY=FALSE` - Error/block, manual intervention required

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-CONTINUATION-FLAG-MASTER-GUIDE.md**
