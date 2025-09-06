# Integration Agent - TESTING State Rules

## State Definition
The TESTING state validates the integrated code through build and test execution.

## Required Actions

### 1. Build Attempt
```bash
# Attempt build - DO NOT FIX if it fails
make build || BUILD_STATUS="FAILED"
echo "Build Status: $BUILD_STATUS"
```

### 2. Test Execution
```bash
# Run tests - DO NOT FIX failures
make test || TEST_STATUS="FAILED"
go test ./... -v > test-output.txt 2>&1
```

### 3. Coverage Analysis (if available)
```bash
go test ./... -cover -coverprofile=coverage.out
go tool cover -func=coverage.out
```

## CRITICAL RULE
- R266 - Upstream Bug Documentation
  - **NEVER FIX** any bugs found
  - **ONLY DOCUMENT** issues
  - Mark all bugs as "NOT FIXED - upstream"

## Documentation Requirements
Document in INTEGRATION-REPORT.md:
- Build results (success/failure)
- Test results (passed/failed count)
- Failed test details
- Upstream bugs found
- Recommendations (but don't implement)

## Transition Rules
- Can transition to: REPORTING
- Always transitions regardless of test results
- Must document all findings

## Success Criteria
- Build attempted and documented
- Tests executed and results captured
- All bugs documented (not fixed)
- Results added to INTEGRATION-REPORT.md

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
