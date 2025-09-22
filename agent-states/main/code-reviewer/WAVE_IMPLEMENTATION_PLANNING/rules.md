# Code-reviewer - WAVE_IMPLEMENTATION_PLANNING State Rules

## State Context
This is the WAVE_IMPLEMENTATION_PLANNING state for the code-reviewer.

## Acknowledgment Required
Thank you for reading the rules file for the WAVE_IMPLEMENTATION_PLANNING state.

**IMPORTANT**: Please report that you have successfully read the WAVE_IMPLEMENTATION_PLANNING rules file.

Say: "✅ Successfully read WAVE_IMPLEMENTATION_PLANNING rules for code-reviewer"

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

When creating wave implementation plans, you MUST ensure EVERY effort can be implemented as an atomic PR:

1. **Wave-Level Atomic PR Planning**
   - Each effort in wave = ONE independent PR
   - PRs can merge in any order (unless documented)
   - No cross-effort dependencies that block merging

2. **Feature Flag Coordination**
   - Wave-level feature flags for incomplete features
   - Effort-level flags for granular control
   - Document flag hierarchy and dependencies
   - Plan coordinated activation strategy

3. **Interface Implementation Sequencing**
   - Define which effort implements interface
   - Plan stub usage until real implementation
   - Ensure each PR maintains contracts
   - Document replacement sequence

4. **Parallel vs Sequential Efforts**
   - Identify which efforts can PR in parallel
   - Document any required sequencing
   - Ensure parallel efforts don't conflict
   - Plan merge conflict resolution

5. **Testing Independence**
   - Each effort PR must test independently
   - No test coupling between efforts
   - Feature flag permutation testing
   - Integration tests per PR

### Wave Implementation Plan Must Include

```yaml
wave_implementation_atomic_design:
  parallel_pr_efforts:
    - effort_1: "User auth - can merge anytime"
    - effort_3: "Logging - can merge anytime"
  sequential_pr_efforts:
    - effort_2: "Requires effort_1 interface"
  feature_flags:
    wave_flag: "WAVE_1_ENABLED"
    effort_flags:
      - "EFFORT_1_AUTH_ENABLED"
      - "EFFORT_2_PROFILE_ENABLED"
  stub_plan:
    - stub: "MockAuthService"
      replaced_by: "effort_1"
      used_by: ["effort_2", "effort_4"]
  merge_strategy:
    order_independent: true
    conflict_resolution: "Rebase on main"
  test_independence:
    each_pr_isolated: true
    flag_permutations: ["all off", "each on", "all on"]
```

**VIOLATION = -100% IMMEDIATE FAILURE**

## General Responsibilities
Follow all general code-reviewer rules and the Software Factory state machine.

## Next Steps
Proceed with the standard workflow for the WAVE_IMPLEMENTATION_PLANNING state as defined in the state machine.
