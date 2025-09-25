# PHASE 1 TEST PLAN - COMMAND SKELETON & FOUNDATION
**Project**: idpbuilder-push
**Methodology**: Test-Driven Development (TDD)
**Generated**: 2025-09-22
**Test Author**: Code Reviewer Agent

## Overview

This test plan defines comprehensive test scenarios for Phase 1 of the idpbuilder push command implementation. All tests MUST be written BEFORE implementation begins, following strict TDD methodology. These tests will initially FAIL with clear messages to guide implementation.

## Test Categories

### 1. Command Registration Tests
Tests to verify the push command is properly registered with idpbuilder's Cobra command tree.

### 2. Flag Definition Tests
Tests to verify all command flags are properly defined and accessible.

### 3. Argument Validation Tests
Tests to verify command arguments are properly validated.

### 4. Help Text Tests
Tests to verify help and usage documentation is complete and accurate.

### 5. Error Handling Tests
Tests to verify appropriate error messages for various failure scenarios.

### 6. Integration Tests
Tests to verify the command integrates properly with the existing idpbuilder CLI.

## Coverage Mapping to Architecture Promises

| Architecture Promise | Test Category | Test File | Test Functions |
|---------------------|---------------|-----------|----------------|
| Command registration with Cobra | Command Registration | command_test.go | TestPushCommandExists, TestPushCommandRegistered |
| Flag parsing (--username, --password, --insecure) | Flag Definition | flags_test.go | TestUsernameFlag, TestPasswordFlag, TestInsecureFlag |
| Required arguments (IMAGE, REGISTRY) | Argument Validation | validation_test.go | TestRequiresTwoArguments, TestValidatesImagePath, TestValidatesRegistryURL |
| Help text and documentation | Help Text | help_test.go | TestHelpTextContent, TestUsageString, TestLongDescription |
| Error messages and validation | Error Handling | validation_test.go | TestInvalidArgumentCount, TestInvalidImagePath, TestInvalidRegistryURL |
| Integration with idpbuilder | Integration | integration_test.go | TestCommandInRootHelp, TestCommandExecution, TestFlagPrecedence |

## Detailed Test Scenarios

### Command Registration Tests (command_test.go)

#### Test: Push Command Exists
```go
func TestPushCommandExists(t *testing.T) {
    // Test that the push command can be instantiated
    // Expected: Command object created successfully
    // Initial state: FAIL - "push command not implemented"
}
```

#### Test: Push Command Registered with Root
```go
func TestPushCommandRegistered(t *testing.T) {
    // Test that push command appears in root command's subcommands
    // Expected: Push command found in subcommand list
    // Initial state: FAIL - "push command not registered"
}
```

### Flag Definition Tests (flags_test.go)

#### Test: Username Flag
```go
func TestUsernameFlag(t *testing.T) {
    // Test --username/-u flag is defined
    // Test accepts string value
    // Test default is empty string
    // Initial state: FAIL - "username flag not defined"
}
```

#### Test: Password Flag
```go
func TestPasswordFlag(t *testing.T) {
    // Test --password/-p flag is defined
    // Test accepts string value
    // Test default is empty string
    // Initial state: FAIL - "password flag not defined"
}
```

#### Test: Insecure Flag
```go
func TestInsecureFlag(t *testing.T) {
    // Test --insecure/-k flag is defined
    // Test accepts boolean value
    // Test default is false
    // Initial state: FAIL - "insecure flag not defined"
}
```

#### Test: Flag Shorthand Support
```go
func TestFlagShorthands(t *testing.T) {
    // Test -u works for --username
    // Test -p works for --password
    // Test -k works for --insecure
    // Initial state: FAIL - "shorthand flags not implemented"
}
```

### Argument Validation Tests (validation_test.go)

#### Test: Requires Exactly Two Arguments
```go
func TestRequiresTwoArguments(t *testing.T) {
    // Test with 0 arguments - should fail
    // Test with 1 argument - should fail
    // Test with 2 arguments - should pass
    // Test with 3+ arguments - should fail
    // Initial state: FAIL - "argument validation not implemented"
}
```

#### Test: Validates Image Path
```go
func TestValidatesImagePath(t *testing.T) {
    // Test with valid file path
    // Test with valid directory path
    // Test with invalid path (path traversal)
    // Test with non-existent path
    // Initial state: FAIL - "image path validation not implemented"
}
```

#### Test: Validates Registry URL
```go
func TestValidatesRegistryURL(t *testing.T) {
    // Test with valid HTTPS URL
    // Test with valid HTTP URL (should warn)
    // Test with invalid URL format
    // Test with missing protocol
    // Initial state: FAIL - "registry URL validation not implemented"
}
```

### Help Text Tests (help_test.go)

#### Test: Help Text Content
```go
func TestHelpTextContent(t *testing.T) {
    // Test short description is present
    // Test long description is present
    // Test usage string is correct
    // Test all flags are documented
    // Initial state: FAIL - "help text not defined"
}
```

#### Test: Usage String Format
```go
func TestUsageString(t *testing.T) {
    // Test usage shows: "push IMAGE REGISTRY"
    // Test positional arguments are documented
    // Initial state: FAIL - "usage string not defined"
}
```

#### Test: Examples in Help
```go
func TestHelpExamples(t *testing.T) {
    // Test at least one example is present
    // Test examples are properly formatted
    // Initial state: FAIL - "examples not provided"
}
```

### Error Handling Tests (validation_test.go)

#### Test: Invalid Argument Count Errors
```go
func TestInvalidArgumentCount(t *testing.T) {
    // Test error message for no arguments
    // Test error message for one argument
    // Test error message for too many arguments
    // Initial state: FAIL - "error handling not implemented"
}
```

#### Test: Invalid Image Path Errors
```go
func TestInvalidImagePath(t *testing.T) {
    // Test error for non-existent file
    // Test error for invalid path characters
    // Test error for path traversal attempts
    // Initial state: FAIL - "image validation errors not implemented"
}
```

#### Test: Invalid Registry URL Errors
```go
func TestInvalidRegistryURL(t *testing.T) {
    // Test error for malformed URL
    // Test error for unsupported protocol
    // Test error for missing host
    // Initial state: FAIL - "registry validation errors not implemented"
}
```

#### Test: Missing Credential Handling
```go
func TestMissingCredentials(t *testing.T) {
    // Test behavior when no credentials provided
    // Test informative error message
    // Initial state: FAIL - "credential validation not implemented"
}
```

### Integration Tests (integration_test.go)

#### Test: Command Appears in Root Help
```go
func TestCommandInRootHelp(t *testing.T) {
    // Execute: idpbuilder --help
    // Verify: "push" appears in available commands
    // Initial state: FAIL - "push command not in help"
}
```

#### Test: Command Execution Flow
```go
func TestCommandExecution(t *testing.T) {
    // Test PreRun hooks execute
    // Test Run function executes
    // Test PostRun hooks execute
    // Initial state: FAIL - "execution flow not implemented"
}
```

#### Test: Flag Precedence
```go
func TestFlagPrecedence(t *testing.T) {
    // Test CLI flag overrides environment variable
    // Test environment variable used if no CLI flag
    // Test default used if neither provided
    // Initial state: FAIL - "flag precedence not implemented"
}
```

#### Test: Environment Variable Support
```go
func TestEnvironmentVariables(t *testing.T) {
    // Test IDPBUILDER_REGISTRY_USERNAME
    // Test IDPBUILDER_REGISTRY_PASSWORD
    // Test IDPBUILDER_INSECURE
    // Initial state: FAIL - "environment variable support not implemented"
}
```

## Test Execution Order

### Phase 1 Test Sequence
1. Command Registration Tests (foundation)
2. Flag Definition Tests (configuration)
3. Argument Validation Tests (input handling)
4. Help Text Tests (documentation)
5. Error Handling Tests (robustness)
6. Integration Tests (system level)

### Dependencies Between Tests
- Flag tests depend on command registration
- Validation tests depend on flag definitions
- Integration tests depend on all previous tests

## Expected Initial State

### All Tests Must Initially Fail
When tests are first run (before implementation), every test should fail with a clear, actionable message:

```
FAIL: TestPushCommandExists
Error: push command not implemented
Hint: Create cmd/push/root.go with NewPushCommand() function

FAIL: TestUsernameFlag
Error: username flag not defined
Hint: Add flag definition in push command init()

FAIL: TestRequiresTwoArguments
Error: argument validation not implemented
Hint: Set Args: cobra.ExactArgs(2) in command definition
```

## Test Coverage Requirements

### Phase 1 Coverage Targets
- Command Registration: 100%
- Flag Handling: 100%
- Argument Validation: 100%
- Help Text: 100%
- Error Paths: 100%
- Overall Phase 1: 100%

### Critical Path Coverage
All critical paths must have 100% coverage:
- Command initialization
- Flag parsing
- Argument validation
- Error message generation

## Test Data and Fixtures

### Sample Test Data
```go
// Valid test inputs
validImagePath := "./test-image.tar"
validRegistryURL := "https://gitea.cnoe.localtest.me"
validUsername := "testuser"
validPassword := "testpass"

// Invalid test inputs
invalidImagePath := "../../../etc/passwd"
invalidRegistryURL := "not-a-url"
emptyString := ""
```

### Mock Command Structure
```go
// Minimal command structure for testing
type mockRootCommand struct {
    subcommands []*cobra.Command
}

type mockPushCommand struct {
    args []string
    flags map[string]interface{}
}
```

## Success Criteria

### Test Plan Success
- ✅ All Phase 1 architectural promises have tests
- ✅ Tests are executable before implementation
- ✅ Tests fail with helpful messages
- ✅ Test structure follows Go conventions
- ✅ Tests use standard testing libraries

### Implementation Guidance Success
- ✅ Each failing test provides clear error message
- ✅ Error messages include implementation hints
- ✅ Tests define expected behavior clearly
- ✅ Test names indicate what feature to implement

## Test Maintenance

### As Implementation Progresses
1. Tests should gradually turn from RED to GREEN
2. No test should be modified to pass
3. Only implementation code should change
4. All tests must remain at 100% coverage

### After Phase 1 Complete
1. All tests should pass
2. No tests should be skipped
3. Coverage report should show 100%
4. Tests become regression suite

## Integration with Test Harness

The test harness (`PHASE-TEST-HARNESS.sh`) will execute all tests in the correct order and report results:

```bash
#!/bin/bash
# Phase 1 Test Harness

echo "Running Phase 1 Tests..."

# Run each test category
go test ./tests/phase1/command_test.go -v
go test ./tests/phase1/flags_test.go -v
go test ./tests/phase1/validation_test.go -v
go test ./tests/phase1/help_test.go -v
go test ./tests/phase1/integration_test.go -v

# Generate coverage report
go test ./tests/phase1/... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

echo "Test execution complete. See coverage.html for details."
```

## Notes for Implementation Teams

### TDD Workflow Reminder
1. **RED**: Run tests, see failures
2. **GREEN**: Write minimal code to pass
3. **REFACTOR**: Improve code quality

### Do NOT:
- Skip writing tests first
- Modify tests to make them pass
- Write more code than needed to pass tests
- Ignore failing tests

### DO:
- Run tests frequently
- Commit tests before implementation
- Keep tests simple and focused
- Use test failures to guide implementation

---

*This test plan ensures Phase 1 implementation follows strict TDD methodology. All tests must be written and failing before any implementation code is created.*