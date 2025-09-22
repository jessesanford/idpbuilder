# Effort 1.1.1: Write Command Tests - Implementation Plan

## CRITICAL EFFORT METADATA (FROM WAVE PLAN)
**Branch**: `idpbuilderpush/phase1/wave1/command-tests`
**Can Parallelize**: No
**Parallel With**: None
**Size Estimate**: 150 lines
**Dependencies**: None (First effort - foundational)

## Overview
- **Effort**: Write comprehensive TDD RED phase tests for push command
- **Phase**: 1, Wave: 1.1
- **Estimated Size**: 150 lines
- **Implementation Time**: 1-2 hours
- **TDD Phase**: RED (tests must FAIL initially)

## Library Version Requirements (R381)
**CRITICAL**: ALL versions below are IMMUTABLE per R381

### Locked Dependencies (DO NOT UPDATE)
- `github.com/spf13/cobra` v1.8.0 (LOCKED - already in go.mod)
- `github.com/stretchr/testify` v1.9.0 (LOCKED - already in go.mod)
- `github.com/spf13/pflag` v1.0.5 (transitive via cobra)

These versions are already established in the project. DO NOT suggest or implement any version updates.

## TDD Context
This is the RED phase of Test-Driven Development. All tests written here MUST:
1. FAIL when first executed (no implementation exists)
2. Define the expected behavior clearly
3. Cover all command registration and flag scenarios
4. Be comprehensive enough to drive the GREEN phase implementation

## File Structure
- `cmd/push/root_test.go`: Command registration and flag tests (150 lines)
  - Test package declaration and imports
  - Command registration tests
  - Flag validation tests
  - Help text tests
  - Environment variable tests
  - Default value tests

## Detailed Test Implementation Plan

### Test Structure and Organization

#### Package and Imports (15 lines)
```go
package push

import (
    "testing"
    "os"
    "bytes"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/spf13/cobra"
    "github.com/spf13/pflag"
)
```

#### Test 1: TestPushCommandRegistration (20 lines)
- Verify command exists and is properly registered
- Check command Use, Short, and Long descriptions
- Verify command can be executed (will fail initially)
- Assert command has expected name "push"

#### Test 2: TestPushCommandFlags (30 lines)
- Test all required flags are registered:
  - `--username` (string)
  - `--password` (string)
  - `--namespace` (string)
  - `--dir` (string)
  - `--insecure` (bool)
  - `--plain-http` (bool)
- Verify flag types and default values
- Check flag descriptions are present

#### Test 3: TestPushCommandArgValidation (20 lines)
- Test with no arguments (should fail)
- Test with valid registry URL argument
- Test with multiple arguments (should fail)
- Test with invalid URL format

#### Test 4: TestPushCommandHelp (15 lines)
- Capture help output
- Verify help contains usage examples
- Check for flag descriptions
- Verify command description present

#### Test 5: TestPushCommandFlagShorthands (15 lines)
- Verify shorthand flags work:
  - `-u` for `--username`
  - `-p` for `--password`
  - `-n` for `--namespace`
  - `-d` for `--dir`

#### Test 6: TestPushCommandEnvVariables (20 lines)
- Test IDPBUILDER_USERNAME env var
- Test IDPBUILDER_PASSWORD env var
- Test IDPBUILDER_NAMESPACE env var
- Verify env vars are read when flags not provided

#### Test 7: TestPushCommandDefaults (15 lines)
- Verify default namespace is "idpbuilder"
- Verify default dir is "."
- Verify insecure defaults to false
- Verify plain-http defaults to false

### PushConfig Struct Definition
The tests will reference a PushConfig struct that should contain:
```go
type PushConfig struct {
    RegistryURL string
    Username    string
    Password    string
    Namespace   string
    Dir         string
    Insecure    bool
    PlainHTTP   bool
}
```

## Implementation Steps
1. **Create test file structure**: Set up cmd/push/root_test.go
2. **Add package declaration**: Define package push
3. **Import dependencies**: Add all required testing imports
4. **Write failing tests**: Implement all 7 test functions
5. **Define PushConfig**: Include struct definition in tests
6. **Run tests to verify RED**: Ensure all tests fail as expected
7. **Document failures**: Note specific failure messages for GREEN phase

## Size Management
- **Estimated Lines**: 150 lines
- **Measurement Tool**: ${PROJECT_ROOT}/tools/line-counter.sh
- **Check Frequency**: After completing each test function
- **Split Threshold**: N/A (well under 800 line limit)

## Expected Test Output (RED Phase)
All tests should fail with errors like:
- "undefined: pushCmd"
- "undefined: PushConfig"
- "command not found"
- "flags not registered"

This is EXPECTED and CORRECT for the RED phase.

## Test Requirements
- **Coverage Goal**: Define behavior for 100% of command surface
- **Test Independence**: Each test must be runnable in isolation
- **Clear Assertions**: Each assertion should have a clear failure message
- **No Implementation**: Tests only - no production code

## Pattern Compliance
- **Go Testing Patterns**: Table-driven tests where appropriate
- **Cobra Patterns**: Standard Cobra command testing approach
- **Error Messages**: Use descriptive assertion messages
- **Test Naming**: Follow Go convention Test[Function][Scenario]

## Out of Scope (DO NOT IMPLEMENT)
- L Actual command implementation (root.go)
- L Registry client creation
- L Authentication logic
- L OCI push operations
- L Integration with registry
- L Credential retrieval from secrets
- L Complex error handling
- L Network operations

## Success Criteria
-  All 7 test functions written and failing
-  Total test code d 150 lines
-  Tests cover all command aspects
-  Clear test names and assertions
-  No production code written
-  Tests define complete expected behavior

## Next Steps (for Effort 1.1.2)
The SW Engineer implementing the GREEN phase will:
1. Read these failing tests
2. Implement minimal code to make tests pass
3. Follow the behavior defined by tests exactly
4. Not add functionality beyond test requirements