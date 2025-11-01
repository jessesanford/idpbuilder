# Code-reviewer - WAVE_IMPLEMENTATION_PLANNING State Rules

## 🔴🔴🔴 ABSOLUTE REQUIREMENT: STATE MANAGER CONSULTATION 🔴🔴🔴

**SUPREME LAW - R517 - UNIVERSAL STATE MANAGER CONSULTATION LAW**

**BEFORE EXITING THIS STATE, YOU MUST:**

1. **SPAWN STATE MANAGER FOR SHUTDOWN_CONSULTATION** (MANDATORY - NO EXCEPTIONS)
2. **NEVER UPDATE STATE FILES DIRECTLY** (orchestrator-state-v3.json, bug-tracking.json, etc.)
3. **NEVER COMMIT STATE FILES YOURSELF** (State Manager does this atomically)

**FAILURE TO CONSULT STATE MANAGER = IMMEDIATE SYSTEM HALT (-100% GRADE)**

### Enforcement Mechanism

If you attempt to exit this state without spawning State Manager:
- ❌ Pre-commit hooks will REJECT your commit
- ❌ Validation tools will FAIL the build
- ❌ Grading system will assign -100% penalty
- ❌ System will transition to ERROR_RECOVERY

### Required Pattern (COPY THIS EXACTLY)

```bash
# At end of state work, BEFORE any state file updates:

echo "🔴 MANDATORY: Spawning State Manager for SHUTDOWN_CONSULTATION"

# Spawn State Manager (REQUIRED - NOT OPTIONAL)
# Task: state-manager
# State: SHUTDOWN_CONSULTATION
# Current State: [YOUR_CURRENT_STATE]
# Proposed Next State: [YOUR_PROPOSED_NEXT_STATE]
# Work Summary: [SUMMARY_OF_WORK_COMPLETED]

# State Manager will:
# 1. Validate proposed transition against state machine
# 2. Update all 4 state files atomically
# 3. Commit with [R288] tag
# 4. Return REQUIRED next state (may differ from proposal)

# Wait for State Manager response
# Follow State Manager's directive (REQUIRED next state)
# DO NOT proceed until State Manager confirms
```

**YOU MUST NEVER:**
- ❌ Update orchestrator-state-v3.json yourself
- ❌ Update bug-tracking.json yourself
- ❌ Update integration-containers.json yourself
- ❌ Use `jq` to modify state files
- ❌ Use `sed/awk` to modify state files
- ❌ Set `validated_by: "orchestrator"` (must be "state-manager")
- ❌ Skip State Manager consultation "just this once"
- ❌ Think "I'll validate it manually"

**ONLY State Manager may update state files. This is NON-NEGOTIABLE.**

See: `rule-library/R517-universal-state-manager-consultation-law.md`
## State Context
This is the WAVE_IMPLEMENTATION_PLANNING state for the code-reviewer within SF 3.0.

## SF 3.0 Wave Planning Context

In this state, the Code Reviewer creates wave implementation plans:
- Reads wave architecture and test plans from orchestrator-state-v3.json `metadata_locations` per R340
- Creates wave implementation plan with effort breakdown and dependencies
- Reports plan location for orchestrator to record in `metadata_locations.wave_implementation_plans`
- Wave plan guides orchestrator's effort-level execution and parallelization decisions
- All planning artifacts stored with atomic state updates per R288

## Acknowledgment Required
Thank you for reading the rules file for the WAVE_IMPLEMENTATION_PLANNING rules file.

**IMPORTANT**: Please report that you have successfully read the WAVE_IMPLEMENTATION_PLANNING rules file.

Say: "✅ Successfully read WAVE_IMPLEMENTATION_PLANNING rules for code-reviewer"

## Critical Rules Referenced

This state enforces the following critical rules:

### R355 - Production Ready Code Enforcement (SUPREME LAW)
- **File**: `$CLAUDE_PROJECT_DIR/rule-library/R355-production-ready-code-enforcement-supreme-law.md`
- **Classification**: SUPREME LAW
- **Penalty**: -100% for production code violations

**Production Readiness Requirements for Wave Plans**:
1. Plans must NOT allow TODO comments without valid future implementation
2. Plans must NOT include stub code that remains unfixed in wave scope
3. Plans must require configuration-based values (no hardcoding)
4. Plans must ensure all wave functionality is actually implemented
5. Plans must verify wave implementation completeness

**Wave-Level TODO Planning**:
- If wave plan defers functionality beyond wave, MUST specify exact future wave/phase
- Must identify exact file and line for TODO placement
- Must document removal criteria in plan
- Vague deferrals = Planning failure

### R332 - Mandatory Bug Filing Protocol (SUPREME LAW)
- **File**: `$CLAUDE_PROJECT_DIR/rule-library/R332-mandatory-bug-filing-protocol.md`
- **Classification**: SUPREME LAW
- **Integration**: Ensures all wave implementation issues are tracked

**Wave Planning Requirements**:
1. Plan must include wave-level bug tracking
2. Plan must define bug filing protocol for wave
3. Plan must prevent "pre-existing" bug excuses
4. Plan must ensure all wave issues documented

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

