# Code-reviewer - PHASE_IMPLEMENTATION_PLANNING State Rules

## State Context
This is the PHASE_IMPLEMENTATION_PLANNING state for the code-reviewer.

## Acknowledgment Required
Thank you for reading the rules file for the PHASE_IMPLEMENTATION_PLANNING state.

**IMPORTANT**: Please report that you have successfully read the PHASE_IMPLEMENTATION_PLANNING rules file.

Say: "✅ Successfully read PHASE_IMPLEMENTATION_PLANNING rules for code-reviewer"

## State-Specific Rules

### 🔴🔴🔴 ATOMIC PR IMPLEMENTATION REQUIREMENTS (R220 - SUPREME LAW) 🔴🔴🔴

When creating phase implementation plans, you MUST ensure EVERY effort can be implemented as an atomic PR:

1. **Each Effort = One Atomic PR**
   - Implementation plan must result in ONE PR per effort
   - PR must be mergeable to main independently
   - No multi-effort PRs allowed

2. **Feature Flag Implementation Details**
   - Specify exact flag names and locations
   - Document flag initialization and defaults
   - Plan testing with flags on/off
   - Include flag cleanup strategy

3. **Stub Implementation Planning**
   - Identify all external dependencies
   - Plan mock/stub implementations
   - Document when stubs get replaced
   - Ensure stubs maintain interface contracts

4. **Interface Contract Definition**
   - Define all interfaces before implementation
   - Document expected behavior
   - Plan for future extensions
   - Ensure backward compatibility

5. **Testing Strategy for Atomic PRs**
   - Each PR must have complete test coverage
   - Tests must pass with feature flags off
   - Integration tests for gradual activation
   - No test dependencies on other PRs

### Implementation Plan Must Include

```yaml
phase_implementation_atomic_design:
  effort_pr_mapping: "1 effort = 1 PR"
  feature_flag_implementation:
    - flag: "PHASE_1_FEATURES"
      location: "config/features.yaml"
      default: false
      testing: "Test with flag on and off"
  stub_implementations:
    - name: "MockPaymentGateway"
      implements: "IPaymentGateway"
      replacement_effort: "effort_5"
  interface_definitions:
    - interface: "IUserService"
      methods: ["authenticate", "authorize"]
      implementation_efforts: ["effort_1", "effort_2"]
  pr_testing_strategy:
    isolated_tests: true
    flag_coverage: true
    backward_compatible: true
```

**VIOLATION = -100% IMMEDIATE FAILURE**

## General Responsibilities
Follow all general code-reviewer rules and the Software Factory state machine.

## Next Steps
Proceed with the standard workflow for the PHASE_IMPLEMENTATION_PLANNING state as defined in the state machine.
