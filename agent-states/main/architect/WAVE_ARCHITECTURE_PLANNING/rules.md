# Architect - WAVE_ARCHITECTURE_PLANNING State Rules

## State Context
This is the WAVE_ARCHITECTURE_PLANNING state for the architect.

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
