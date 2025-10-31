# ⚠️⚠️⚠️ RULE R265: Integration Testing Requirements ⚠️⚠️⚠️

## Rule Definition
**Criticality:** WARNING - Must attempt, document failures
**Category:** Agent-Specific
**Applies To:** integration-agent

## Testing Requirements

### 1. Build Verification
The integration agent MUST attempt:
- Build the integrated code
- Document build success/failure
- Capture build output
- **DO NOT FIX** build errors - only document

```bash
# Build attempt sequence
make build || BUILD_FAILED=true
if [[ "$BUILD_FAILED" == "true" ]]; then
    echo "⚠️ Build failed - documenting in INTEGRATE_WAVE_EFFORTS-REPORT.md"
    # Capture error output for report
fi
```

### 2. Test Execution
The integration agent MUST:
- Run all available tests
- Document test results
- Identify failing tests
- **DO NOT FIX** test failures - only document

```bash
# Test execution sequence
make test || TEST_FAILED=true
go test ./... -v > test-output.txt 2>&1
TEST_PASSED=$(grep -c "PASS:" test-output.txt)
TEST_FAILED=$(grep -c "FAIL:" test-output.txt)
```

### 3. Upstream Bug Documentation
**CRITICAL: DO NOT FIX UPSTREAM BUGS**

When bugs are found:
- Document the bug location
- Describe the issue
- Suggest potential fix
- Mark as "UPSTREAM BUG - NOT FIXED"
- Include in INTEGRATE_WAVE_EFFORTS-REPORT.md

### 4. Test Coverage Analysis
If coverage tools available:
```bash
# Coverage analysis
go test ./... -cover -coverprofile=coverage.out
go tool cover -func=coverage.out > coverage-summary.txt
COVERAGE=$(grep "total:" coverage-summary.txt | awk '{print $3}')
```

## Test Results Documentation

```markdown
## BUILD AND TEST RESULTS

### Build Attempt
Command: make build
Start Time: HH:MM:SS
End Time: HH:MM:SS
Status: PROJECT_DONE / FAILURE

#### Build Output (if failed):
```
error: undefined reference to 'functionX'
error: cannot find package "github.com/missing/package"
```

### Test Execution
Command: make test
Total Tests: 156
Passed: 148
Failed: 8
Skipped: 0
Coverage: 76.5%

#### Failed Tests:
1. TestAuthenticationFlow
   - File: auth/auth_test.go:45
   - Error: timeout waiting for response
   - Type: UPSTREAM BUG
   - Recommendation: Increase timeout to 30s

2. TestDatabaseConnection
   - File: db/connection_test.go:23  
   - Error: connection refused
   - Type: ENVIRONMENT ISSUE
   - Recommendation: Check database configuration

### Static Analysis (if available)
Command: golangci-lint run
Issues Found: 12
- 5 : unused variables (minor)
- 3 : error not checked (major)
- 4 : inefficient regex compilation (performance)
```

## Enforcement

```bash
# Verify testing was attempted
verify_test_execution() {
    local report="$1"
    
    # Check for test section
    if ! grep -q "Test Execution" "$report"; then
        echo "⚠️ WARNING: No test execution documented"
        return 1
    fi
    
    # Check for build attempt
    if ! grep -q "Build Attempt" "$report"; then
        echo "⚠️ WARNING: No build attempt documented"
        return 1
    fi
    
    # Verify no fixes were attempted
    if grep -q "Fixed\|Patched\|Corrected" "$report"; then
        echo "❌ ERROR: Integration agent should not fix bugs!"
        return 1
    fi
}
```

## What NOT to Do

```bash
# ❌ WRONG - Don't fix bugs
if test_fails; then
    vim src/broken_test.go  # NO! Document only
    git commit -m "fix: test"  # NO! Don't commit fixes
fi

# ✅ CORRECT - Document only  
if test_fails; then
    echo "Test failed: documenting in report"
    echo "- TestX failed at file:line" >> INTEGRATE_WAVE_EFFORTS-REPORT.md
    echo "  Recommendation: [suggestion]" >> INTEGRATE_WAVE_EFFORTS-REPORT.md
fi
```

## Grading Impact
- 10% - Build attempt and documentation
- 10% - Test execution and results
- 5% - Proper upstream bug documentation
- -20% PENALTY if bugs are fixed instead of documented

## Related Rules
- R263 - Integration Documentation Requirements
- R266 - Upstream Bug Documentation
- R267 - Integration Agent Grading Criteria