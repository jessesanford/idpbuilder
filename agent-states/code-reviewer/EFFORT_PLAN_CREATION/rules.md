# Code-reviewer - EFFORT_PLAN_CREATION State Rules

## State Context
This is the EFFORT_PLAN_CREATION state for the code-reviewer.

## Acknowledgment Required
Thank you for reading the rules file for the EFFORT_PLAN_CREATION state.

**IMPORTANT**: Please report that you have successfully read the EFFORT_PLAN_CREATION rules file.

Say: "✅ Successfully read EFFORT_PLAN_CREATION rules for code-reviewer"

## State-Specific Rules

### 🔴🔴🔴 ATOMIC PR EFFORT REQUIREMENTS (R220 - SUPREME LAW) 🔴🔴🔴

When creating effort implementation plans, you MUST ensure the effort produces exactly ONE atomic PR:

1. **One Effort = One PR (ABSOLUTE)**
   - This effort must result in EXACTLY one PR to main
   - PR must merge independently of all other efforts
   - PR must not break the build when merged alone
   - NO EXCEPTIONS TO THIS RULE

2. **Feature Flags for This Effort**
   - Define specific flags for incomplete features
   - Document exact implementation location
   - Include flag initialization code
   - Plan tests with flag on/off
   - Specify cleanup conditions

3. **Stubs for Dependencies**
   - Identify what this effort depends on
   - Create stubs for missing dependencies
   - Ensure stubs match interface contracts
   - Document when stubs get replaced
   - Test with both stubs and real implementations

4. **Interface Implementation**
   - If defining interface: complete specification
   - If implementing interface: match contract exactly
   - Support both current and future use cases
   - Maintain backward compatibility
   - Document any assumptions

5. **PR Completeness Checklist**
   - All code for effort in ONE PR
   - All tests pass independently
   - Feature flags control activation
   - Documentation included
   - No dependencies on unmerged PRs

### Effort Plan MUST Include

```yaml
effort_atomic_pr_design:
  pr_summary: "Single PR implementing [specific feature]"
  can_merge_to_main_alone: true  # MUST be true
  
  feature_flags_needed:
    - flag: "EFFORT_X_FEATURE_Y"
      purpose: "Hide incomplete feature Y"
      default: false
      location: "config/features.yaml"
      activation: "When all components ready"
  
  stubs_required:
    - stub: "MockServiceZ"
      replaces: "ServiceZ (from effort_5)"
      interface: "IServiceZ"
      behavior: "Returns default success response"
  
  interfaces_to_implement:
    - interface: "IDataProcessor"
      methods: ["process", "validate"]
      implementation: "Complete in this PR"
  
  pr_verification:
    tests_pass_alone: true
    build_remains_working: true
    flags_tested_both_ways: true
    no_external_dependencies: true
    backward_compatible: true
  
  example_pr_structure:
    files_added:
      - "src/feature_x.go"
      - "src/feature_x_test.go"
      - "config/features.yaml"
      - "stubs/mock_service_z.go"
    tests_included:
      - "Unit tests with flag off"
      - "Unit tests with flag on"
      - "Integration test with stubs"
    documentation:
      - "README update"
      - "API documentation"
```

### CRITICAL VALIDATION

Before completing effort plan, verify:
- ✅ This effort = ONE atomic PR to main
- ✅ PR can merge without any other effort
- ✅ Build stays green when PR merges
- ✅ Feature flags hide incomplete work
- ✅ All dependencies stubbed/mocked
- ✅ Tests pass in complete isolation

**FAILURE TO ENSURE ATOMIC PR = -100% IMMEDIATE FAILURE**

## General Responsibilities
Follow all general code-reviewer rules and the Software Factory state machine.

## Next Steps
Proceed with the standard workflow for the EFFORT_PLAN_CREATION state as defined in the state machine.
