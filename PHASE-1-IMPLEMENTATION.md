# PHASE 1: FOUNDATION & COMMAND STRUCTURE IMPLEMENTATION PLAN

## Overview
**Phase**: 1
**Title**: Foundation & Command Structure
**TDD Requirement**: Tests are already written in `phase-tests/phase-1/`
**Total Estimated Size**: 800 lines
**Duration**: 2-3 days
**Branch Prefix**: idpbuilderpush/phase1

### Phase Objectives
- Establish the `idpbuilder push` command foundation
- Implement command registration with Cobra
- Define and implement all CLI flags (--username, --password, --insecure)
- Implement basic input validation and error handling
- Ensure TDD compliance by implementing against existing tests

### Success Criteria
✅ All tests in `phase-tests/phase-1/tests/phase1/` pass
✅ Command properly integrated with idpbuilder CLI
✅ Flags parse correctly with validation
✅ Help text is comprehensive
✅ 80%+ test coverage maintained
✅ Each effort is an atomic, mergeable PR

## Wave Breakdown

### Wave 1: Command Structure & Skeleton
**Wave ID**: 1.1
**Theme**: Core command infrastructure
**Total Lines**: ~500
**Can Parallelize Between Efforts**: No (sequential dependencies)

### Wave 2: Flag Implementation & Validation
**Wave ID**: 1.2
**Theme**: Input handling and validation
**Total Lines**: ~300
**Can Parallelize Between Efforts**: Yes

## Effort Details

### Wave 1.1: Command Structure & Skeleton

#### Effort 1.1.1: Command Registration
**Branch**: `idpbuilderpush/phase1/wave1/command-registration`
**Can Parallelize**: No
**Parallel With**: None (first effort)
**Size Estimate**: 150 lines
**Dependencies**: None

**Files to Create/Modify**:
- `pkg/cmd/push/root.go` - Basic command structure
- `pkg/cmd/root.go` - Register push command with main CLI

**Implementation Requirements**:
1. Create push command using Cobra
2. Register with root command
3. Define basic command metadata (Use, Short, Long descriptions)
4. Implement empty Run function that returns nil
5. Make tests in `command_test.go` pass

**Test Coverage**:
- `TestPushCommandExists` must pass
- `TestPushCommandRegistered` must pass

**Atomic PR Requirements**:
- Command is registered but does nothing (safe to merge)
- No functionality exposed yet
- Tests validate structure only

#### Effort 1.1.2: Flag Definitions
**Branch**: `idpbuilderpush/phase1/wave1/flag-definitions`
**Can Parallelize**: No
**Parallel With**: None (depends on 1.1.1)
**Size Estimate**: 200 lines
**Dependencies**: Effort 1.1.1

**Files to Create/Modify**:
- `pkg/cmd/push/root.go` - Add flag definitions
- `pkg/cmd/push/flags.go` - Flag helper functions

**Implementation Requirements**:
1. Define --username/-u flag (string, default empty)
2. Define --password/-p flag (string, default empty)
3. Define --insecure/-k flag (bool, default false)
4. Implement flag initialization in init() function
5. Create flag accessor functions

**Test Coverage**:
- All tests in `flags_test.go` must pass
- `TestUsernameFlag`, `TestPasswordFlag`, `TestInsecureFlag`
- `TestFlagShorthands`, `TestFlagDefaults`

**Atomic PR Requirements**:
- Flags are defined but not used
- No behavior changes
- Backward compatible

#### Effort 1.1.3: Help Text & Documentation
**Branch**: `idpbuilderpush/phase1/wave1/help-documentation`
**Can Parallelize**: No
**Parallel With**: None (depends on 1.1.2)
**Size Estimate**: 150 lines
**Dependencies**: Efforts 1.1.1, 1.1.2

**Files to Create/Modify**:
- `pkg/cmd/push/root.go` - Complete help text
- `pkg/cmd/push/examples.go` - Usage examples

**Implementation Requirements**:
1. Write comprehensive Long description
2. Add usage examples
3. Document all flags in help text
4. Ensure proper formatting
5. Add argument placeholders (IMAGE REGISTRY)

**Test Coverage**:
- All tests in `help_test.go` must pass
- `TestHelpTextContent`, `TestUsageString`, `TestLongDescription`
- `TestExampleSection`, `TestFlagDescriptions`

**Atomic PR Requirements**:
- Only documentation changes
- No functional impact
- Improves user experience

### Wave 2: Flag Implementation & Validation

#### Effort 1.2.1: Argument Validation
**Branch**: `idpbuilderpush/phase1/wave2/argument-validation`
**Can Parallelize**: Yes
**Parallel With**: Effort 1.2.2
**Size Estimate**: 150 lines
**Dependencies**: Wave 1 completion

**Files to Create/Modify**:
- `pkg/cmd/push/validation.go` - Validation functions
- `pkg/cmd/push/root.go` - Add validation to PreRunE

**Implementation Requirements**:
1. Validate exactly 2 arguments required (IMAGE, REGISTRY)
2. Validate image path exists and is accessible
3. Validate registry URL format
4. Return clear error messages
5. Implement PreRunE function

**Test Coverage**:
- Tests in `validation_test.go` must pass
- `TestRequiresTwoArguments`, `TestValidatesImagePath`
- `TestValidatesRegistryURL`, `TestInvalidArgumentCount`

**Feature Flag**:
```go
const pushCommandEnabled = false // Set to true when ready
```

**Atomic PR Requirements**:
- Validation logic is complete
- Feature flag prevents actual execution
- Tests verify validation behavior

#### Effort 1.2.2: Error Handling & Messages
**Branch**: `idpbuilderpush/phase1/wave2/error-handling`
**Can Parallelize**: Yes
**Parallel With**: Effort 1.2.1
**Size Estimate**: 150 lines
**Dependencies**: Wave 1 completion

**Files to Create/Modify**:
- `pkg/cmd/push/errors.go` - Error definitions
- `pkg/cmd/push/root.go` - Error handling in RunE

**Implementation Requirements**:
1. Define custom error types
2. Implement user-friendly error messages
3. Add error wrapping with context
4. Ensure consistent error format
5. Handle panic recovery

**Test Coverage**:
- Error handling tests in `validation_test.go`
- `TestInvalidImagePath`, `TestInvalidRegistryURL`
- `TestMissingArguments`, `TestErrorMessageFormat`

**Atomic PR Requirements**:
- Error handling is comprehensive
- No side effects
- Improves debugging

## Library Version Requirements (R381)
**CRITICAL**: ALL versions below are IMMUTABLE per R381

### Existing Dependencies (DO NOT UPDATE)
These dependencies are already in the idpbuilder project:
- `github.com/spf13/cobra v1.8.1` - CLI framework (LOCKED)
- `github.com/go-logr/logr v1.4.2` - Logging interface (LOCKED)
- Other existing dependencies in go.mod

### New Dependencies Allowed
For Phase 1, no new dependencies are required. The command structure uses only:
- Existing Cobra framework
- Go standard library

**WARNING**: Any version update requires user approval + full cascade per R382!

## Size Management Strategy
- **Measurement Tool**: `tools/line-counter.sh`
- **Measurement Frequency**: After each effort completion
- **Current Phase Total**: ~800 lines (well within limits)
- **Per-Effort Limits**: All efforts <200 lines (safe margin)

## Test Requirements
**Pre-existing Tests**: All tests for Phase 1 already exist in `phase-tests/phase-1/tests/phase1/`
- `command_test.go` - Command registration tests
- `flags_test.go` - Flag definition tests
- `help_test.go` - Documentation tests
- `validation_test.go` - Validation tests
- `integration_test.go` - Integration tests

**TDD Compliance**:
- ✅ Tests already written (RED phase complete)
- ⏳ Implementation makes tests pass (GREEN phase)
- ⏳ Refactor for quality after passing (REFACTOR phase)

## Integration Points
1. **Command Registration**: Add to existing `pkg/cmd/root.go`
2. **Import Path**: Use `github.com/cnoe-io/idpbuilder/pkg/cmd/push`
3. **Logging**: Reuse existing logr configuration
4. **Error Patterns**: Match existing idpbuilder error style

## Atomic PR Strategy (R220 Compliance)
Each effort creates one atomic, independently mergeable PR:

1. **Effort 1.1.1**: Command exists but does nothing
2. **Effort 1.1.2**: Flags defined but not processed
3. **Effort 1.1.3**: Help text only, no behavior
4. **Effort 1.2.1**: Validation with feature flag off
5. **Effort 1.2.2**: Error handling improvements

All PRs can be merged to main independently without breaking existing functionality.

## Wave Completion Criteria

### Wave 1 Complete When:
- [ ] Command appears in `idpbuilder --help`
- [ ] All flags are defined and show in help
- [ ] Help text is comprehensive
- [ ] Tests in phase-tests/phase-1 for Wave 1 pass

### Wave 2 Complete When:
- [ ] Argument validation works correctly
- [ ] Error messages are clear
- [ ] All Phase 1 tests pass
- [ ] 80%+ test coverage achieved

## Risk Mitigation
1. **Sequential Dependencies in Wave 1**: Cannot parallelize, must complete in order
2. **Test Compatibility**: Implementation must match existing test expectations exactly
3. **Integration Risk**: Each PR tested against main branch before merge

## Next Steps
1. Create Wave 1, Effort 1.1.1 working directory
2. Implement command registration to pass first tests
3. Measure with line-counter.sh after implementation
4. Submit PR when tests pass
5. Proceed to next effort

---
*Phase 1 enforces strict TDD with pre-written tests. Every implementation must make existing tests pass.*