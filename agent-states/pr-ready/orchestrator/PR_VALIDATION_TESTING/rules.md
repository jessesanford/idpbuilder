# PR_VALIDATION_TESTING State Rules

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
   find . -name "orchestrator-state.json" -o -name "*.todo" 2>/dev/null

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