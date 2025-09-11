# Rule R060: Test Implementation Protocol

## Rule Statement
All code MUST have corresponding tests before being marked complete. Tests must cover core functionality, edge cases, and error conditions. Test coverage must be verified before integration, and all tests must pass before phase completion.

## Criticality Level
**BLOCKING** - Untested code cannot be integrated or deployed

## Enforcement Mechanism
- **Technical**: Test execution gates at integration points
- **Behavioral**: Tests required for effort completion
- **Grading**: -40% for missing tests, -100% for integration without tests

## Core Principle

```
Test Implementation = Unit Tests → Integration Tests → Validation → All Pass
NEVER mark code complete without tests
NEVER integrate failing tests
ALWAYS test edge cases and errors
Tests are NOT optional
```

## Detailed Requirements

### Test Coverage Requirements

1. **Unit Tests**
   ```python
   # Every function/method needs tests
   def test_function_normal_case():
       """Test normal operation"""
       assert function(valid_input) == expected_output
   
   def test_function_edge_case():
       """Test boundary conditions"""
       assert function(edge_input) == edge_output
   
   def test_function_error_case():
       """Test error handling"""
       with pytest.raises(ExpectedException):
           function(invalid_input)
   ```

2. **Integration Tests**
   ```python
   # Test component interactions
   def test_component_integration():
       """Test components work together"""
       component_a = ComponentA()
       component_b = ComponentB()
       result = component_a.process(component_b.output())
       assert result.is_valid()
   ```

3. **End-to-End Tests**
   ```python
   # Test complete workflows
   def test_full_workflow():
       """Test entire feature flow"""
       setup_test_environment()
       execute_workflow()
       verify_results()
       cleanup_test_environment()
   ```

### Test Implementation Stages

1. **During Development**
   - Write tests alongside code
   - Run tests after each change
   - Fix failures immediately

2. **Before Commit**
   ```bash
   # Run all tests
   npm test           # or
   pytest             # or
   cargo test         # or
   go test ./...
   
   # Verify coverage
   npm run coverage   # or equivalent
   ```

3. **Before Integration**
   ```bash
   # Full test suite
   ./run-all-tests.sh
   
   # Coverage report
   coverage report --fail-under=80
   ```

### Test Standards

1. **Descriptive Names**: Test names describe what they test
2. **Isolation**: Tests don't depend on each other
3. **Repeatability**: Tests produce same results every run
4. **Speed**: Unit tests complete in <1 second
5. **Coverage**: Minimum 80% code coverage

### Test Organization

```
tests/
├── unit/
│   ├── test_component_a.py
│   └── test_component_b.py
├── integration/
│   └── test_workflow.py
└── e2e/
    └── test_full_system.py
```

### Failure Handling

When tests fail:
1. **STOP** implementation immediately
2. **FIX** the failing test or code
3. **VERIFY** fix doesn't break other tests
4. **DOCUMENT** the issue and resolution
5. **CONTINUE** only when all tests pass

## Relationship to Other Rules
- **R035**: Phase completion testing
- **R265**: Integration testing requirements
- **R272**: Integration testing branch
- **R273**: Runtime-specific validation

## Implementation Notes
- Tests must be committed with the code they test
- Test files must follow project naming conventions
- Mock external dependencies in unit tests
- Use test fixtures for common setup
- Document complex test scenarios