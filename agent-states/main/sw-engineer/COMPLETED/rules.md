# Sw-engineer - COMPLETED State Rules

## State Context
This is the COMPLETED state for the sw-engineer.

## Acknowledgment Required
Thank you for reading the rules file for the COMPLETED state.

**IMPORTANT**: Please report that you have successfully read the COMPLETED rules file.

Say: "✅ Successfully read COMPLETED rules for sw-engineer"

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
No additional state-specific rules are defined for this state at this time.

## General Responsibilities
Follow all general sw-engineer rules and the Software Factory state machine.

## Next Steps
Proceed with the standard workflow for the COMPLETED state as defined in the state machine.
