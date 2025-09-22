# PHASE TEST PLANNING INSTRUCTIONS (TDD)

## Your Mission
Create comprehensive tests for Phase 1 BEFORE implementation begins. This enforces Test-Driven Development (TDD) at the phase level.

## Context
- **Phase**: 1
- **Architecture Document**: /home/vscode/workspaces/idpbuilder-push/phase-1-architecture.md
- **Implementation Plan**: /home/vscode/workspaces/idpbuilder-push/IMPLEMENTATION-PLAN.md
- **Project**: idpbuilder-push (adding push command to idpbuilder for OCI registry support)

## Required Deliverables

### 1. PHASE-TEST-PLAN.md
Document containing:
- Overview of all test scenarios for Phase 1
- Test categorization (command tests, validation tests, error tests)
- Coverage mapping to architecture promises
- Expected behaviors and assertions
- Test execution order and dependencies

### 2. PHASE-TEST-HARNESS.sh
Executable shell script that:
- Runs all Phase 1 tests in sequence
- Reports pass/fail for each test category
- Generates summary report
- Can run even before implementation (will fail initially)
- Supports both individual and batch test execution

### 3. tests/phase1/*.test.go
Actual test files containing:
- Command interface tests (skeleton testing)
- Flag validation tests
- Error handling tests
- Help text validation
- All tests MUST fail initially (no implementation yet)
- Tests use standard Go testing framework

### 4. PHASE-DEMO-PLAN.md
Demo scenarios document containing:
- Interactive demo scripts for Phase 1 functionality
- Step-by-step execution guides
- Expected outputs for each demo
- Integration points with tests
- User experience validation scenarios

## Test Requirements

### Technical Requirements
- Tests must be FUNCTIONAL/BEHAVIORAL (not unit tests)
- Tests must validate ALL promised capabilities from architecture
- Tests must FAIL initially (no implementation yet)
- Tests must be executable and automated
- Tests must include both success and failure scenarios
- Tests must use idpbuilder's existing test patterns

### Coverage Requirements
From Phase 1 Architecture, you must test:
1. Command registration and availability
2. All command-line flags functionality
3. Help text accuracy and completeness
4. Error messages and validation
5. Command lifecycle (init, validate, execute)
6. Integration with existing idpbuilder commands

### Test Categories
1. **Command Tests**: Verify command exists and is callable
2. **Flag Tests**: Validate all flags work as specified
3. **Validation Tests**: Test input validation and error cases
4. **Help Tests**: Ensure documentation is complete
5. **Integration Tests**: Verify command fits within idpbuilder

## Demo Integration (R330/R291 Consolidation)
- Demo scenarios must be designed alongside tests
- Each test category should have corresponding demo scenarios
- Demos should showcase both success and error paths
- Integration points must be clearly documented
- User workflows should be validated through demos

## Implementation Guidelines

### File Structure
```
phase-tests/phase-1/
├── PHASE-TEST-PLAN.md
├── PHASE-TEST-HARNESS.sh
├── PHASE-DEMO-PLAN.md
└── tests/
    └── phase1/
        ├── command_test.go
        ├── flags_test.go
        ├── validation_test.go
        ├── help_test.go
        └── integration_test.go
```

### Test Naming Convention
- Test files: `*_test.go`
- Test functions: `Test<Feature><Scenario>`
- Example: `TestPushCommandExists`, `TestPushFlagsValidation`

### Expected Initial State
All tests should fail with clear messages like:
- "Push command not implemented"
- "Flag --registry not yet available"
- "Validation logic pending implementation"

## Success Criteria
- ✅ All architectural promises have corresponding tests
- ✅ Test harness is executable (even if tests fail)
- ✅ Demo scenarios are clearly defined
- ✅ Tests are ready for implementation teams to target
- ✅ Clear failure messages guide implementation
- ✅ Tests follow Go and idpbuilder conventions

## Notes
- Review the Phase 1 architecture thoroughly before creating tests
- Focus on WHAT to test, not HOW it's implemented
- Tests define the contract that implementation must fulfill
- These tests will guide SW Engineers during implementation
- Demo scenarios will be used for stakeholder presentations