# Code-reviewer - PHASE_IMPLEMENTATION_PLANNING State Rules

## State Context
This is the PHASE_IMPLEMENTATION_PLANNING state for the code-reviewer within SF 3.0.

## SF 3.0 Phase Planning Context

In this state, the Code Reviewer creates phase implementation plans:
- Reads phase architecture and test plans from orchestrator-state-v3.json `metadata_locations` per R340
- Creates comprehensive phase implementation plan with wave breakdown
- Reports plan location for orchestrator to record in `metadata_locations.phase_implementation_plans`
- Phase plan guides orchestrator's wave-by-wave execution strategy
- All planning artifacts stored with atomic state updates per R288

## Acknowledgment Required
Thank you for reading the rules file for the PHASE_IMPLEMENTATION_PLANNING state.

**IMPORTANT**: Please report that you have successfully read the PHASE_IMPLEMENTATION_PLANNING rules file.

Say: "✅ Successfully read PHASE_IMPLEMENTATION_PLANNING rules for code-reviewer"

## Critical Rules Referenced

This state enforces the following critical rules:

### R355 - Production Ready Code Enforcement (SUPREME LAW)
- **File**: `$CLAUDE_PROJECT_DIR/rule-library/R355-production-ready-code-enforcement-supreme-law.md`
- **Classification**: SUPREME LAW
- **Penalty**: -100% for production code violations

**Production Readiness Requirements for Plans**:
1. Plans must NOT allow TODO comments without valid future implementation
2. Plans must NOT include stub code that remains unfixed
3. Plans must require configuration-based values (no hardcoding)
4. Plans must ensure all functionality is actually implemented
5. Plans must verify implementation completeness

**TODO Planning Criteria**:
- If plan defers functionality, MUST specify exact future effort
- Must identify exact file and line for TODO placement
- Must document removal criteria in plan
- Vague deferrals = Planning failure

### R332 - Mandatory Bug Filing Protocol (SUPREME LAW)
- **File**: `$CLAUDE_PROJECT_DIR/rule-library/R332-mandatory-bug-filing-protocol.md`
- **Classification**: SUPREME LAW
- **Integration**: Ensures all implementation issues are tracked

**Planning Phase Requirements**:
1. Plan must include bug tracking system setup
2. Plan must define bug filing protocol
3. Plan must prevent "pre-existing" bug excuses
4. Plan must ensure all issues documented

**See**: R332 for complete bug filing protocol.

## 🔴🔴🔴 MANDATORY MEASUREMENT RULE - R304 🔴🔴🔴

### R304: Mandatory Line Counter Tool Enforcement
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R304-mandatory-line-counter-enforcement.md`
**Criticality**: BLOCKING - Manual counting = AUTOMATIC -100% FAILURE

**ABSOLUTE REQUIREMENTS:**
- ✅ MUST use `$CLAUDE_PROJECT_DIR/tools/line-counter.sh` for ALL line counting
- ❌ NEVER use `wc -l` or any manual counting method
- ❌ NEVER count lines any other way - this is a -100% automatic failure
- ✅ MUST specify both -b (base branch) and -c (current branch) parameters
- ✅ Base branch MUST be phase integration branch (NOT "main")

**Failure to use the line counter tool = AUTOMATIC -100% GRADE**

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


## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

### YOU MUST OUTPUT THE CONTINUATION FLAG AS YOUR LAST ACTION

**EVERY STATE COMPLETION MUST END WITH EXACTLY ONE OF:**
```bash
# If state completed successfully and factory should continue:
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"

# If error/block/manual review needed:
echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
```

### CRITICAL REQUIREMENTS:
1. **ABSOLUTE LAST OUTPUT**: This MUST be the very last line of output before state completion
2. **NO TEXT AFTER**: No explanations, summaries, or any text after the flag
3. **EXACTLY AS SHOWN**: Use exact format - no variations like CONTINUE-ORCHESTRATING
4. **ALWAYS REQUIRED**: Every single state must output this flag
5. **GREPPABLE**: Must be on its own line for automation parsing

### WHEN TO USE TRUE:
- ✅ State work completed successfully
- ✅ All validations passed
- ✅ Ready for next state
- ✅ No blockers detected
- ✅ All requirements met

### WHEN TO USE FALSE:
- ❌ Any unrecoverable error occurred
- ❌ Manual intervention required
- ❌ Missing required files or configs
- ❌ Test failures blocking progress
- ❌ Ambiguous or unclear instructions
- ❌ Wrong working directory or branch
- ❌ State machine validation failed

### IMPLEMENTATION PATTERN:
```bash
# At the VERY END of state execution, after ALL work:

# Determine success/failure
if [[ "$STATE_COMPLETED_PROJECT_DONEFULLY" == "true" ]]; then
    echo "✅ State work completed successfully"
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
else
    echo "❌ State encountered issues requiring intervention"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
fi

# DO NOT OUTPUT ANYTHING AFTER THE FLAG!
```

### GRADING IMPACT:
- **Missing flag**: -100% AUTOMATIC FAILURE
- **Text after flag**: -50% penalty
- **Wrong format**: -100% AUTOMATIC FAILURE
- **Multiple flags**: -50% penalty

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

