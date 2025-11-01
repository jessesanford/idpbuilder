# PR_VALIDATION_TESTING State Rules

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
## 🔴🔴🔴 STATE PURPOSE: Final Validation Before PR 🔴🔴🔴

### 🚨🚨🚨 RULE R369 - PR Validation and Integrity Protocol [BLOCKING]
**Source:** `$CLAUDE_PROJECT_DIR/rule-library/R369-pr-validation-and-integrity-protocol.md`
**Criticality:** BLOCKING - Invalid PRs = Production failures

This state implements R369 requirements for comprehensive validation.

### 🚨🚨🚨 RULE R271 - Mandatory Production-Ready Validation [BLOCKING]
**Source:** `$CLAUDE_PROJECT_DIR/rule-library/R271-mandatory-production-ready-validation.md`
**Criticality:** BLOCKING - Production readiness required

### 🔴🔴🔴 RULE R355 - Production Ready Code Enforcement [SUPREME LAW]
**Source:** `$CLAUDE_PROJECT_DIR/rule-library/R355-production-ready-code-enforcement-supreme-law.md`
**Criticality:** SUPREME LAW - Zero tolerance for non-production code

### MANDATORY ACTIONS (R233 + R369 + R271 + R355 COMPLIANT)
Upon entering this state, IMMEDIATELY:

1. **Run comprehensive validation suite**
   ```bash
   # Execute R369 validation protocol
   echo "🔍 Starting PR validation suite..."

   # Check for SF artifacts (R365)
   find . -name "orchestrator-state-v3.json" -o -name "*.todo" 2>/dev/null

   # Check for stubs/TODOs (R355)
   grep -r "TODO\|FIXME\|STUB" --include="*.js" --include="*.ts" --include="*.py" . | grep -v test

   # Run build validation
   npm run build || mvn package || go build || cargo build

   # Run all tests
   npm test || mvn test || go test ./... || cargo test
   ```

2. **Verify sequential mergeability**
   ```bash
   # Per R363 - verify branches can merge to main sequentially
   for BRANCH in $(cat PR-MERGE-ORDER.txt); do
     git merge --no-commit --no-ff $BRANCH || exit 1
   done
   git reset --hard main
   ```

3. **Generate validation report**
   ```bash
   cat > PR-VALIDATION-REPORT.md << EOF
   # PR Validation Report
   Date: $(date)

   ## R369 Compliance
   - No SF artifacts: ✅
   - No stubs/TODOs: ✅
   - Tests passing: ✅
   - Build successful: ✅

   ## R271 Compliance
   - Production ready: ✅
   - Deployable: ✅

   ## R355 Compliance
   - No non-production code: ✅

   ## R363 Compliance
   - Sequential mergeability: ✅

   ## Overall Status: ✅ PR-READY
   EOF
   ```

### EXIT CRITERIA
✅ All validation checks pass
✅ No SF artifacts present (R365)
✅ No stubs or TODOs (R355)
✅ All tests passing
✅ Build successful
✅ Sequential mergeability verified (R363)
✅ Validation report generated

### OUTPUT FILES
- `PR-VALIDATION-REPORT.md`
- `validation.log`
- `test-results.xml` (if applicable)

### PROHIBITED ACTIONS
❌ Do NOT proceed if ANY validation fails
❌ Do NOT skip any validation category
❌ Do NOT ignore warnings
❌ Do NOT modify code to pass validation

### FAILURE PROTOCOL
If validation fails:
1. Document failure in `PR-VALIDATION-FAILED.md`
2. Transition to `PR_FIX_VALIDATION_ISSUES`
3. Do NOT proceed to PR_PLAN_CREATION

### GRADING CRITERIA
- SF artifacts in PR: -50%
- Stubs/TODOs in code: -40%
- Failing tests: -100%
- No validation performed: -30%

## STATE TRANSITIONS

### Success Path:
```
PR_VALIDATION_TESTING → PR_PLAN_CREATION
```
All validations passed, ready for PR plan

### Failure Path:
```
PR_VALIDATION_TESTING → PR_FIX_VALIDATION_ISSUES
```
Validation failures need resolution

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
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

