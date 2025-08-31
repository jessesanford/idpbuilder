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