# Unit Test Execution & Fixes Implementation Plan

## Effort Metadata
- **Effort**: E2.1.1 - Unit Test Execution & Fixes
- **Phase**: 2, Wave: 1
- **Branch**: `idpbuilder-push-oci/phase2/wave1/unit-test-execution`
- **Base Branch**: `idpbuilder-push-oci/phase1-integration`
- **Estimated Size**: 450 lines
- **Implementation Time**: 3-4 hours
- **Theme**: TEST_EXECUTION_AND_FIXES

## 1. SCOPE DEFINITION

### Files TO BE MODIFIED (In Priority Order)

#### Test Compilation Fixes (150 lines)
1. `pkg/push/pusher_test.go` - Add missing SetError method to mockProgressReporter
2. `pkg/cmd/push/push_test.go` - Fix undefined variable references (username, password, insecureTLS)
3. `pkg/push/operations_test.go` - Verify and fix any compilation issues
4. `pkg/push/discovery_test.go` - Verify and fix any compilation issues

#### Test Coverage Improvements (200 lines)
5. `pkg/push/progress_test.go` (NEW) - Add comprehensive tests for ProgressReporter
6. `pkg/push/auth/auth_test.go` (NEW) - Add tests for authentication components
7. `pkg/push/errors/errors_test.go` (NEW) - Add error handling tests
8. `pkg/push/logging_test.go` (NEW) - Add logging tests

#### Fix Implementation Issues Found by Tests (100 lines)
9. `pkg/push/pusher.go` - Fix any bugs found during testing
10. `pkg/push/operations.go` - Fix any bugs found during testing
11. `pkg/cmd/push/root.go` - Fix any command initialization issues
12. `pkg/push/progress.go` - Fix any progress reporting issues

### Test Files TO BE VERIFIED (No Modifications)
- `pkg/k8s/deserialize_test.go` - Run and verify passes
- `pkg/k8s/util_test.go` - Run and verify passes
- `pkg/kind/config_test.go` - Run and verify passes
- `pkg/kind/cluster_test.go` - Run and verify passes
- `pkg/cmd/get/secrets_test.go` - Run and verify passes
- `pkg/cmd/helpers/validation_test.go` - Run and verify passes
- `pkg/controllers/*_test.go` - Run and verify all controller tests
- `pkg/util/*_test.go` - Run and verify all utility tests
- `pkg/tls/config_test.go` - Run and verify passes

## 2. OUT OF SCOPE

### EXPLICITLY NOT TOUCHING
- ❌ Main application logic changes (unless fixing test-discovered bugs)
- ❌ Adding new features or functionality
- ❌ Refactoring existing code (unless required for test fixes)
- ❌ Documentation updates (separate effort)
- ❌ Integration tests (Phase 2, Wave 2)
- ❌ End-to-end tests (Phase 2, Wave 3)
- ❌ Performance optimizations
- ❌ Dependency updates
- ❌ CI/CD pipeline changes
- ❌ Files outside the `pkg/` directory (except test utilities)

## 3. THEME DECLARATION

**SINGLE THEME**: TEST_EXECUTION_AND_FIXES

This effort focuses exclusively on:
1. Making all existing unit tests compile and pass
2. Fixing implementation bugs discovered by tests
3. Adding missing unit tests to achieve >80% coverage
4. Documenting edge cases found during testing

All changes serve this single theme of test execution and fixing.

## 4. IMPLEMENTATION STEPS

### Step 1: Fix Test Compilation Issues (50 lines)
```go
// File: pkg/push/pusher_test.go
// Add the missing SetError method to mockProgressReporter
func (m *mockProgressReporter) SetError(digest string, err error) {
    // Implementation for tracking errors in tests
}

// File: pkg/cmd/push/push_test.go
// Fix undefined variables by using command flags properly
// Access flags via cmd.Flags() instead of module-level variables
```

### Step 2: Run Existing Tests and Document Failures (0 lines - analysis only)
```bash
# Run all tests with coverage
go test -v ./pkg/... -cover -coverprofile=coverage.out 2>&1 | tee test-results.txt

# Generate coverage report
go tool cover -html=coverage.out -o coverage.html

# Document all failures and their causes
```

### Step 3: Fix Implementation Bugs Found by Tests (100 lines)
Based on test failures, fix:
- Nil pointer dereferences
- Incorrect error handling
- Missing validations
- Race conditions
- Resource leaks

### Step 4: Add Missing Unit Tests for push Package (150 lines)
```go
// pkg/push/progress_test.go - Test ConsoleProgressReporter
// - Test StartImage/FinishImage lifecycle
// - Test concurrent layer updates
// - Test error reporting
// - Test progress calculation

// pkg/push/auth/auth_test.go - Test authentication
// - Test credential extraction
// - Test various auth types
// - Test auth failures

// pkg/push/errors/errors_test.go - Test error handling
// - Test error wrapping
// - Test error types
// - Test retry logic for transient errors

// pkg/push/logging_test.go - Test logging
// - Test log levels
// - Test structured logging
// - Test sensitive data redaction
```

### Step 5: Add Missing Unit Tests for cmd Package (50 lines)
```go
// Fix and enhance pkg/cmd/push/push_test.go
// - Test command initialization
// - Test flag validation
// - Test argument parsing
// - Test error scenarios
```

### Step 6: Achieve Coverage Target (100 lines)
```bash
# Identify packages below 80% coverage
go test ./pkg/... -cover | grep -E "coverage: [0-7][0-9]?\.[0-9]%"

# Add targeted tests for uncovered code paths
# Focus on error conditions and edge cases
```

### Step 7: Document Edge Cases
Create EDGE-CASES.md documenting:
- Boundary conditions found
- Error scenarios tested
- Performance considerations
- Security implications

## 5. TESTING REQUIREMENTS

### Unit Test Coverage
- **Target**: >80% overall package coverage
- **Critical Packages**:
  - `pkg/push/*`: Must have >85% coverage
  - `pkg/cmd/push/*`: Must have >80% coverage
- **Test Types**:
  - Positive path tests
  - Error condition tests
  - Boundary condition tests
  - Concurrent operation tests

### Test Execution Requirements
- All tests must pass: `go test ./pkg/... -race`
- No race conditions detected
- No resource leaks
- Tests must be deterministic (no flaky tests)
- Tests must complete in <30 seconds total

### Test Documentation
- Each test must have clear description
- Complex test scenarios must have comments
- Test data must be minimal and focused
- Mock objects must be properly documented

## 6. SUCCESS CRITERIA

### Must Have (Blocking)
- ✅ All existing tests compile without errors
- ✅ All tests pass (0 failures)
- ✅ Overall code coverage >80%
- ✅ Critical packages coverage >85%
- ✅ No race conditions detected
- ✅ All discovered bugs fixed
- ✅ Edge cases documented

### Should Have (Important)
- ✅ Test execution time <30 seconds
- ✅ Clear test failure messages
- ✅ Consistent test patterns across packages
- ✅ No test interdependencies

### Nice to Have (Optional)
- ⭕ Coverage badge updated
- ⭕ Test performance benchmarks
- ⭕ Test complexity metrics

## 7. SIZE MANAGEMENT

### Measurement Strategy
- Use `${PROJECT_ROOT}/tools/line-counter.sh` for all measurements
- Measure after each step completion
- Target: Stay under 500 lines total

### Size Breakdown
- Test fixes: 50 lines
- Bug fixes: 100 lines
- New tests: 300 lines
- **Total Estimated**: 450 lines

### Split Threshold
- Warning at 400 lines
- Stop and split at 450 lines
- Each test file should be <150 lines

## 8. RISK MITIGATION

### Identified Risks
1. **Disk Space Issues**: Tests failing due to "no space left on device"
   - Mitigation: Clean up test artifacts between runs
   - Use t.TempDir() for test files

2. **Compilation Errors**: Undefined variables in tests
   - Mitigation: Fix immediately before running tests
   - Verify all imports are correct

3. **Flaky Tests**: Non-deterministic test failures
   - Mitigation: Use fixed seeds for random operations
   - Mock time-dependent operations

4. **Coverage Gaps**: Hard-to-test code paths
   - Mitigation: Refactor for testability if needed
   - Use test doubles for external dependencies

## 9. VALIDATION CHECKLIST

Before marking complete:
- [ ] All tests compile successfully
- [ ] `go test ./pkg/...` passes with no failures
- [ ] `go test -race ./pkg/...` shows no race conditions
- [ ] Coverage report shows >80% overall
- [ ] Push package has >85% coverage
- [ ] All test files follow naming convention (*_test.go)
- [ ] No hardcoded test data or credentials
- [ ] Edge cases documented in EDGE-CASES.md
- [ ] Test execution time <30 seconds
- [ ] All TODOs from test failures addressed

## 10. NOTES

### Current Issues Found
1. `mockProgressReporter` missing `SetError` method
2. `push_test.go` references undefined module-level variables
3. Several packages have 0% coverage
4. Some tests failing due to disk space issues

### Dependencies
- This effort depends on Phase 1 implementation being complete
- Must be based on `idpbuilder-push-oci/phase1-integration` branch
- No external dependencies should be added

### Follow-up Work
- Integration tests (Phase 2, Wave 2)
- E2E tests (Phase 2, Wave 3)
- Performance testing (Phase 3)
- Security testing (Phase 3)

---
Generated: 2025-10-02 17:12:00 UTC