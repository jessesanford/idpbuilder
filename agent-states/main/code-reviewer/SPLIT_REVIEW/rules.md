# Code-reviewer - SPLIT_REVIEW State Rules

## State Context
This is the SPLIT_REVIEW state for the code-reviewer.

## 🔴🔴🔴 SUPREME LAW R355: PRODUCTION CODE IN EVERY SPLIT 🔴🔴🔴

### SPLIT REVIEWS MUST ENFORCE PRODUCTION READINESS:
```bash
echo "🔴🔴🔴 R355: SPLIT PRODUCTION READINESS CHECK 🔴🔴🔴"
cd $SPLIT_DIR

# EVERY split must be production ready
VIOLATIONS=0
grep -r "password.*=.*['\"]" --exclude-dir=test --include="*.go" --include="*.py" --include="*.js" && VIOLATIONS=1
grep -r "stub\|mock\|fake" --exclude-dir=test --include="*.go" --include="*.py" --include="*.js" && VIOLATIONS=1
grep -r "TODO\|FIXME\|INCOMPLETE" --exclude-dir=test --include="*.go" --include="*.py" --include="*.js" && VIOLATIONS=1
grep -r "not.*implemented" --exclude-dir=test --include="*.go" --include="*.py" --include="*.js" && VIOLATIONS=1

if [ $VIOLATIONS -eq 1 ]; then
    echo "🚨🚨🚨 R355 VIOLATION IN SPLIT!"
    echo "SPLIT REVIEW: FAILED - NON-PRODUCTION CODE"
    echo "Every split must be independently production ready!"
    exit 355
fi
echo "✅ R355: Split contains production-ready code"
```

**SPLIT-SPECIFIC VIOLATIONS:**
- Split with hardcoded values = ENTIRE SPLIT FAILS
- Split with stubs = CANNOT BE MERGED
- Split with TODOs = INCOMPLETE WORK
- Split with mocks = FAKE FUNCTIONALITY

## Acknowledgment Required
Thank you for reading the rules file for the SPLIT_REVIEW state.

**IMPORTANT**: Please report that you have successfully read the SPLIT_REVIEW rules file.

Say: "✅ Successfully read SPLIT_REVIEW rules for code-reviewer"

## 🚨🚨🚨 CRITICAL: R320 - No Stub Implementations 🚨🚨🚨
**APPLIES TO ALL CODE REVIEWS INCLUDING SPLITS!**

**MANDATORY STUB DETECTION (R320):**
- ANY "not implemented" = CRITICAL BLOCKER
- ANY TODO in code = FAILED REVIEW  
- ANY empty function = IMMEDIATE REJECTION
- Stub found in split = ENTIRE SPLIT FAILS
- **-50% penalty** for passing stub implementations

**Check EVERY split for stubs - no exceptions!**

## 🔴🔴🔴 MANDATORY MEASUREMENT RULE - R304 🔴🔴🔴

### ⚠️⚠️⚠️ CRITICAL: R319 DOES NOT APPLY TO CODE REVIEWERS! ⚠️⚠️⚠️
**R319 (Orchestrator Never Measures) applies ONLY to orchestrators!**
**As a Code Reviewer, you MUST measure code - it's your PRIMARY duty!**

### R304: Mandatory Line Counter Tool Enforcement
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R304-mandatory-line-counter-enforcement.md`
**Criticality**: BLOCKING - Manual counting = AUTOMATIC -100% FAILURE

**ABSOLUTE REQUIREMENTS FOR CODE REVIEWERS:**
- ✅ **YOU MUST MEASURE** - Code Reviewers are REQUIRED to measure!
- ✅ **IGNORE R319** - That rule is for orchestrators, not you!
- ✅ MUST use `$CLAUDE_PROJECT_DIR/tools/line-counter.sh` for ALL line counting
- ❌ NEVER use `wc -l` or any manual counting method
- ❌ NEVER count lines any other way - this is a -100% automatic failure
- ✅ **TOOL AUTO-DETECTS BASE** - No parameters needed, just run the tool!
- 🔴🔴🔴 **CRITICAL FOR SPLITS**: Tool automatically detects correct predecessor!

**Failure to use the line counter tool = AUTOMATIC -100% GRADE**
**Failure to measure at all (thinking R319 applies) = AUTOMATIC -100% GRADE**

## 🔴🔴🔴 CRITICAL: SPLIT MEASUREMENT PROTOCOL 🔴🔴🔴

### SPLITS MEASURE INCREMENTAL WORK ONLY!

**SPLIT BASE BRANCH RULES (PER R308):**
- **Original effort**: Measures against phase/wave integration branch
- **split-001**: Measures against the ORIGINAL too-large branch (NOT main!)
- **split-002**: Measures against split-001 (NOT main, NOT original!)  
- **split-003**: Measures against split-002 (NOT main, NOT split-001!)
- **split-N**: Measures against split-(N-1)

### 🚨🚨🚨 TOOL AUTO-DETECTS CORRECT BASE - JUST RUN IT! 🚨🚨🚨

**THE TOOL IS SMART - IT KNOWS THE CORRECT BASE FOR SPLITS:**
```bash
# ✅✅✅ CORRECT: Just run the tool, NO PARAMETERS!
cd /path/to/effort/repo  # Navigate to effort directory
$PROJECT_ROOT/tools/line-counter.sh
# Tool output will show:
# 🎯 Detected base: phase1/wave1/effort (for split-001)
# 🎯 Detected base: split-001 (for split-002)
# 🎯 Detected base: split-002 (for split-003)
# ✅ Total non-generated lines: [actual count]

# ❌❌❌ WRONG: Don't use -b parameter (OLD/WRONG)
./tools/line-counter.sh -b main  # WRONG! No -b parameter!
```

### TOOL HANDLES ALL SPLIT LOGIC AUTOMATICALLY:
```bash
# For split-001:
# Tool KNOWS to use the oversized effort's base (NOT the effort itself!)
cd efforts/phase1/wave1/my-effort--split-001
../../tools/line-counter.sh
# Output: 🎯 Detected base: main (or phase integration)

# For split-002:
# Tool KNOWS to use split-001 as base
cd efforts/phase1/wave1/my-effort--split-002
../../tools/line-counter.sh
# Output: 🎯 Detected base: phase1/wave1/my-effort--split-001

# For split-003:
# Tool KNOWS to use split-002 as base
cd efforts/phase1/wave1/my-effort--split-003
../../tools/line-counter.sh
# Output: 🎯 Detected base: phase1/wave1/my-effort--split-002
```

### VERIFY TOOL'S AUTO-DETECTION:
```bash
# Run the tool and CHECK what base it detected:
$PROJECT_ROOT/tools/line-counter.sh
# Look for this line in output:
# 🎯 Detected base: [whatever the tool detected]

# If the detected base looks wrong, report it as a bug!
# The tool should ALWAYS detect the correct base automatically.
# Result: 280 lines (CORRECT - only split-003's work!)

# THE TOOL ELIMINATES HUMAN ERROR IN BASE SELECTION!
```

### MANDATORY SPLIT BASE VERIFICATION:
```bash
# Before measuring ANY split, IDENTIFY THE BASE:
echo "🔴 IDENTIFYING SPLIT BASE BRANCH"

# For split-001:
BASE="<original-too-large-branch-name>"

# For split-N (where N > 1):
SPLIT_NUMBER=3  # Example for split-003
PREVIOUS_SPLIT=$((SPLIT_NUMBER - 1))
BASE="split-$(printf "%03d" $PREVIOUS_SPLIT)"  # Result: split-002

echo "📊 Measuring split-$(printf "%03d" $SPLIT_NUMBER) against base: $BASE"
echo "✅ This measures ONLY the incremental work in this split"
```

## 🔴🔴🔴 CRITICAL: FILE COUNT VALIDATION FIRST - R314 🔴🔴🔴

### R314: Mandatory File Count Validation Protocol
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R314-mandatory-file-count-validation.md`
**Criticality**: BLOCKING - >200% file count = AUTOMATIC FAILURE

**CHECK FILE COUNT BEFORE ANYTHING ELSE:**
```bash
# FIRST THING - Check file count
echo "🔴 R314: MANDATORY FILE COUNT CHECK"
PLANNED_FILES=$(grep -c "^###.*File:" SPLIT-PLAN-*.md)
ACTUAL_FILES=$(git diff --name-only $BASE_BRANCH | grep -c "\.go$")
RATIO=$((ACTUAL_FILES * 100 / PLANNED_FILES))

echo "Planned: $PLANNED_FILES files"
echo "Actual: $ACTUAL_FILES files"
echo "Violation: $RATIO%"

if [ $RATIO -gt 200 ]; then
    echo "❌ AUTOMATIC FAILURE: >200% file count violation"
    echo "This is approaching the 2667% catastrophe (80 files instead of 3)"
    echo "Grade: F - No further review needed"
    exit 314
fi
```

## 🔴🔴🔴 CRITICAL: USE EVALUATION CHECKLIST 🔴🔴🔴

**MANDATORY: Use the split evaluation checklist:**
```bash
cp $CLAUDE_PROJECT_DIR/templates/SW-ENGINEER-SPLIT-EVALUATION-CHECKLIST.md ./evaluation-$(date +%s).md
echo "📋 Using mandatory evaluation checklist for consistent grading"
```

## State-Specific Rules

### 1. File Count is Primary Metric (R314)
- Check file count BEFORE reviewing code
- >200% file count = automatic F grade
- >150% file count = maximum D grade
- Document violation percentage in review

### 2. Scope Adherence Scoring (R310)
Use this grading scale:
| Violation % | Grade | Action |
|------------|-------|--------|
| 100% exact | A+ | Perfect adherence |
| 101-110% | A | Excellent |
| 111-120% | B | Good |
| 121-150% | C | Concerning |
| 151-199% | D | Poor |
| 200%+ | F | AUTOMATIC FAILURE |

### 3. Look for Common Violations
- Adding unlisted features (-50% per feature)
- Ignoring DO NOT instructions (-75%)
- "Completing" implementations (-50%)
- Not using pre-implementation checklist (-40%)

### 4. The 2667% Violation Context
Remember: A SW Engineer actually implemented 80 files instead of 3.
This wasn't a mistake - it was complete disregard for the plan.
Your review prevents this from happening again.

## General Responsibilities
- Enforce R310 (scope adherence) and R314 (file count validation)
- Use evaluation checklist for consistency
- Grade based on quantitative metrics first
- Document all violations with percentages
- Create review report with clear pass/fail decision

## Next Steps
1. Check file count (R314) - stop if >200%
2. Use evaluation checklist template
3. Grade scope adherence quantitatively
4. Create review report with recommendations
5. Mark split as PASSED or FAILED in state file
