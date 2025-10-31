# Architect - WAVE_ARCHITECTURE_PLANNING State Rules

## State Context
This is the WAVE_ARCHITECTURE_PLANNING state for the architect within SF 3.0.

## SF 3.0 Wave Planning Context

In this state, the Architect creates wave-level architectural plans:
- Reads current phase and wave from `state_machine.current_state` in orchestrator-state-v3.json
- Reviews wave implementation requirements and effort breakdown from orchestrator-state-v3.json
- Creates architectural plans that orchestrator will record in `metadata_locations.wave_architecture_plans` per R340
- Plans are stored with timestamps and metadata for orchestrator tracking
- All plan locations reported back for atomic update to orchestrator-state-v3.json per R288

## Acknowledgment Required
Thank you for reading the rules file for the WAVE_ARCHITECTURE_PLANNING state.

**IMPORTANT**: Please report that you have successfully read the WAVE_ARCHITECTURE_PLANNING rules file.

Say: "✅ Successfully read WAVE_ARCHITECTURE_PLANNING rules for architect"

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

### 🔴🔴🔴 ATOMIC PR ARCHITECTURE REQUIREMENTS (R220 - SUPREME LAW) 🔴🔴🔴

When designing wave architecture, you MUST ensure EVERY effort is designed for atomic PR mergeability:

1. **Each Effort = One Atomic PR**
   - Design efforts to be independently PR-able to main
   - No effort should depend on another to merge
   - Each PR must pass all tests in isolation

2. **Wave-Level Feature Flags**
   - Plan flags to hide incomplete wave functionality
   - Document activation strategy for wave features
   - Ensure gradual feature rollout is possible

3. **Interface Design for Wave**
   - Define interfaces that efforts will implement
   - Ensure interfaces support incremental implementation
   - Plan stub implementations for missing components

4. **Merge Order Planning**
   - Document any required merge sequence
   - Identify efforts that can merge in parallel
   - Ensure no circular dependencies

5. **Build Verification Strategy**
   - Each effort PR must maintain working build
   - Plan how to test partial implementations
   - Document backward compatibility requirements

### Wave Architecture Plan Must Include

```yaml
wave_atomic_design:
  all_efforts_atomic: true
  parallel_mergeable_efforts: ["effort1", "effort3"]
  sequential_efforts: ["effort2 after effort1"]
  feature_flags:
    - flag: "WAVE_1_FEATURE_ENABLED"
      controls: "All wave 1 features"
  interface_contracts:
    - name: "IWaveService"
      implementers: ["effort1", "effort2"]
  stub_requirements:
    - stub: "MockDataService"
      until: "effort3 complete"
```

**VIOLATION = -100% IMMEDIATE FAILURE**

## General Responsibilities
Follow all general architect rules and the Software Factory state machine.

## Next Steps
Proceed with the standard workflow for the WAVE_ARCHITECTURE_PLANNING state as defined in the state machine.


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

