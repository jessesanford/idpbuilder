# Integration Agent - REPORTING State Rules

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
## State Definition
The REPORTING state completes all documentation and prepares final deliverables.

## Required Actions

### 1. Complete .software-factory/INTEGRATE_WAVE_EFFORTS-REPORT--${TIMESTAMP}.md
Must include ALL sections:
- Overview (branches integrated, statistics)
- Errors and Issues Found
- Compensating/Remediation Recommendations  
- Build and Test Results
- Upstream Bugs Identified
- Integration Verification Checklist
- Final State

### 2. Verify Work Log Completeness
```bash
# Ensure work-log is replayable
grep "^Command:" .software-factory/work-log--${TIMESTAMP}.log > replay.sh
bash -n replay.sh  # Verify syntax

# Count operations
OPERATION_COUNT=$(grep -c "^## Operation" .software-factory/work-log--${TIMESTAMP}.log)
echo "Total operations documented: $OPERATION_COUNT"
```

### 3. Report Integration Report Location (R340) and Update State Files

## SF 3.0 State File Updates

This state updates orchestrator-state-v3.json and integration-containers.json:
- Updates `state_machine.current_state` in orchestrator-state-v3.json to reflect reporting completion
- Updates integration container status in `integration-containers.json` with final merge results
- Records report locations in `metadata_locations` field for orchestrator access
- Documents convergence metrics and integration outcomes per SF 3.0 architecture
- All state updates are atomic per R288

```bash
# R340: Report integration report location to orchestrator
INTEGRATE_WAVE_EFFORTS_REPORT_PATH="$(pwd)/.software-factory/INTEGRATE_WAVE_EFFORTS-REPORT--${TIMESTAMP}.md"
WORK_LOG_PATH="$(pwd)/.software-factory/work-log--${TIMESTAMP}.log"

echo "📋 Integration Report: $INTEGRATE_WAVE_EFFORTS_REPORT_PATH"
echo "📋 Work Log: $WORK_LOG_PATH"
echo "Integration Type: $INTEGRATE_WAVE_EFFORTS_TYPE"
echo "R340: Created integration report at: $INTEGRATE_WAVE_EFFORTS_REPORT_PATH"

# Update orchestrator-state-v3.json with report location
yq -i ".metadata_locations.integration_reports += {\"$(date +%s)\": \"$INTEGRATE_WAVE_EFFORTS_REPORT_PATH\"}" "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json"
```

### 4. Commit Documentation
```bash
# Add all documentation to integration branch
git add .software-factory/INTEGRATE_WAVE_EFFORTS-PLAN--${TIMESTAMP}.md .software-factory/work-log--${TIMESTAMP}.log .software-factory/INTEGRATE_WAVE_EFFORTS-REPORT--${TIMESTAMP}.md
git commit -m "docs: complete integration documentation for [branch list]"
git push origin "$INTEGRATE_WAVE_EFFORTS_BRANCH"
```

## Documentation Quality Rules
- R263 - Integration Documentation Requirements
- R264 - Work Log Tracking Requirements

## Final Checklist
Before transition to COMPLETED:
- [ ] .software-factory/INTEGRATE_WAVE_EFFORTS-PLAN--${TIMESTAMP}.md exists and was followed
- [ ] .software-factory/work-log--${TIMESTAMP}.log is complete and replayable
- [ ] .software-factory/INTEGRATE_WAVE_EFFORTS-REPORT--${TIMESTAMP}.md has all sections
- [ ] No original branches were modified
- [ ] No cherry-picks were used
- [ ] All documentation committed and pushed

## Transition Rules
- Can transition to: COMPLETED
- Cannot transition if: Documentation incomplete
- Must have pushed integration branch

## Success Criteria
- All three documents complete
- Documentation committed to integration branch
- Integration branch pushed to remote
- Ready for external review

## 🔴🔴🔴 MANDATORY MEASUREMENT RULE - R304 🔴🔴🔴

### R304: Mandatory Line Counter Tool Enforcement
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R304-mandatory-line-counter-enforcement.md`
**Criticality**: BLOCKING - Manual counting = AUTOMATIC -100% FAILURE

**ABSOLUTE REQUIREMENTS:**
- ✅ MUST use `$CLAUDE_PROJECT_DIR/tools/line-counter.sh` for ALL line counting
- ❌ NEVER use `wc -l` or any manual counting method
- ❌ NEVER count lines any other way - this is a -100% automatic failure
- ✅ MUST specify both -b (base branch) and -c (current branch) parameters
- ✅ Base branch MUST be phase integration branch (NOT "main")

**Failure to use the line counter tool = AUTOMATIC -100% GRADE**


## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

### YOU MUST OUTPUT THE CONTINUATION FLAG AS YOUR LAST ACTION

**EVERY STATE COMPLETION MUST END WITH EXACTLY ONE OF:**
```bash
# If state completed successfully and factory should continue:
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"

# If error/block/manual review needed:
echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
```

### CRITICAL REQUIREMENTS:
1. **ABSOLUTE LAST OUTPUT**: This MUST be the very last line of output before state completion
2. **NO TEXT AFTER**: No explanations, summaries, or any text after the flag
3. **EXACTLY AS SHOWN**: Use exact format - no variations like CONTINUE-ORCHESTRATING
4. **ALWAYS REQUIRED**: Every single state must output this flag
5. **GREPPABLE**: Must be on its own line for automation parsing

### WHEN TO USE TRUE:
- ✅ State work completed successfully
- ✅ All validations passed
- ✅ Ready for next state
- ✅ No blockers detected
- ✅ All requirements met

### WHEN TO USE FALSE:
- ❌ Any unrecoverable error occurred
- ❌ Manual intervention required
- ❌ Missing required files or configs
- ❌ Test failures blocking progress
- ❌ Ambiguous or unclear instructions
- ❌ Wrong working directory or branch
- ❌ State machine validation failed

### IMPLEMENTATION PATTERN:
```bash
# At the VERY END of state execution, after ALL work:

# Determine success/failure
if [[ "$STATE_COMPLETED_PROJECT_DONEFULLY" == "true" ]]; then
    echo "✅ State work completed successfully"
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
else
    echo "❌ State encountered issues requiring intervention"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
fi

# DO NOT OUTPUT ANYTHING AFTER THE FLAG!
```

### GRADING IMPACT:
- **Missing flag**: -100% AUTOMATIC FAILURE
- **Text after flag**: -50% penalty
- **Wrong format**: -100% AUTOMATIC FAILURE
- **Multiple flags**: -50% penalty

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

