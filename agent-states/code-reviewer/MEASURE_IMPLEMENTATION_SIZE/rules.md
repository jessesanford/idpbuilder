# Code Reviewer - MEASURE_IMPLEMENTATION_SIZE State Rules

## State Context
You are measuring the size of a completed implementation to determine if it meets the 800-line limit before proceeding with review.

---

## 🔴🔴🔴 MANDATORY LINE COUNTING REQUIREMENTS 🔴🔴🔴

### 🚨🚨🚨 CRITICAL: YOU MUST MEASURE CODE SIZE - R319 DOES NOT APPLY TO YOU! 🚨🚨🚨

**ATTENTION CODE REVIEWER - READ THIS CAREFULLY:**

**YOU ARE A CODE REVIEWER, NOT AN ORCHESTRATOR!**
- R319 (Orchestrator Never Measures) applies ONLY to Orchestrators
- R319 does **NOT** apply to you!
- R006 (Orchestrator Never Writes/Measures) does **NOT** apply to you!

**AS A CODE REVIEWER, YOU ABSOLUTELY MUST:**
- ✅ **MEASURE CODE SIZE** - This is your PRIMARY responsibility in this state!
- ✅ **USE line-counter.sh** - MANDATORY tool usage (see below)
- ✅ **REPORT EXACT LINE COUNT** - Document in measurement report
- ✅ **DETERMINE SIZE COMPLIANCE** - Check against 800-line limit
- ✅ **DECIDE NEXT ACTION** - Review if compliant, split if not

**FAILURE TO MEASURE = -100% IMMEDIATE FAILURE**

### ⚠️⚠️⚠️ CRITICAL: USE LINE-COUNTER.SH WITH AUTOMATIC BASE DETECTION ⚠️⚠️⚠️

**VIOLATION = -100% IMMEDIATE FAILURE**

### MANDATORY STEPS:
1. **ALWAYS use ${PROJECT_ROOT}/tools/line-counter.sh** - NO EXCEPTIONS
2. **Tool automatically detects correct base branch** - NO -b parameter needed!
3. **NEVER do manual counting** - AUTOMATIC FAILURE (-100%)
4. **Base detection prevents wrong base errors** - No more main vs integration mistakes!
5. **NEVER count test/doc files separately** - tool handles this

### CORRECT USAGE:
```bash
# STEP 1: Navigate to effort directory (IT'S A SEPARATE GIT REPO!)
cd /path/to/effort/directory
pwd  # Confirm location
ls -la .git  # MUST exist - this is the effort's own git repository!

# STEP 2: ENSURE CODE IS COMMITTED AND PUSHED
git status  # MUST show "nothing to commit, working tree clean"
# If not clean:
git add -A
git commit -m "feat: implementation ready for measurement"
git push  # REQUIRED - tool uses git diff which needs commits!

# STEP 3: Get ACTUAL BRANCH NAMES (not directory names!)
CURRENT_BRANCH=$(git branch --show-current)
echo "Current branch: $CURRENT_BRANCH"  # e.g., phase2-wave1-gcr-image-builder

# STEP 4: Find project root (where orchestrator-state.yaml lives)
PROJECT_ROOT=$(pwd); while [ "$PROJECT_ROOT" != "/" ]; do 
    [ -f "$PROJECT_ROOT/orchestrator-state.yaml" ] && break; 
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT"); 
done

# STEP 5: Run the tool - it auto-detects the base!
$PROJECT_ROOT/tools/line-counter.sh $CURRENT_BRANCH
# Or even simpler - let it auto-detect current branch:
$PROJECT_ROOT/tools/line-counter.sh
```

### 🔴🔴🔴 CRITICAL: Directory Names vs Branch Names 🔴🔴🔴

**THE FATAL MISTAKE:**
```bash
# ❌❌❌ WRONG - Using directory name as branch:
cd efforts/phase2/wave1/go-containerregistry-image-builder
./line-counter.sh go-containerregistry-image-builder  # WRONG!

# ✅✅✅ CORRECT - Using actual git branch:
cd efforts/phase2/wave1/go-containerregistry-image-builder
BRANCH=$(git branch --show-current)  # Gets: phase2-wave1-gcr-image-builder
../../tools/line-counter.sh "$BRANCH"  # CORRECT!
```

### BASE BRANCH IDENTIFICATION FOR SPLITS:
```bash
# 🔴🔴🔴 SPLITS USE SEQUENTIAL BRANCHING - NOT INTEGRATION! 🔴🔴🔴
CURRENT_BRANCH=$(git branch --show-current)

if [[ "$CURRENT_BRANCH" =~ split-([0-9]{3}) ]]; then
    SPLIT_NUM="${BASH_REMATCH[1]}"
    SPLIT_NUM_INT=$((10#$SPLIT_NUM))
    
    if [[ $SPLIT_NUM_INT -eq 1 ]]; then
        # split-001 measures against original too-large branch
        BASE="<original-branch-name>"
    else
        # split-N measures against split-(N-1)
        PREV_NUM=$((SPLIT_NUM_INT - 1))
        BASE=$(printf "split-%03d" $PREV_NUM)
    fi
    echo "📊 Split detected! Measuring against: $BASE"
else
    # Regular effort - use phase integration
    BASE="phase${PHASE}/integration"
    echo "📊 Regular effort! Measuring against: $BASE"
fi
```

## Size Compliance Decision Tree

```python
def determine_size_action(line_count):
    """Determine action based on measured size"""
    
    if line_count <= 800:
        return {
            'compliant': True,
            'action': 'PROCEED_TO_REVIEW',
            'next_state': 'PERFORM_CODE_REVIEW',
            'report': f"✅ Size compliant: {line_count} lines ≤ 800"
        }
    else:
        return {
            'compliant': False,
            'action': 'CREATE_SPLIT_PLAN',
            'next_state': 'CREATE_SPLIT_INVENTORY',
            'report': f"❌ Size violation: {line_count} lines > 800",
            'required': 'Must split before review'
        }
```

## Measurement Report Documentation

Create SIZE-MEASUREMENT-REPORT.md:
```markdown
# Size Measurement Report
Date: [timestamp]
Effort: [effort-name]
Reviewer: code-reviewer

## Measurement Details
- Tool Used: line-counter.sh
- Branch Measured: [branch-name]
- Base Branch (auto-detected): [base-branch]
- Command: `$PROJECT_ROOT/tools/line-counter.sh [branch]`

## Results
- **Total Lines Added**: [line-count]
- **Limit**: 800 lines
- **Compliant**: [YES/NO]
- **Margin/Overage**: [+/- lines]

## Raw Output
```
[Paste complete tool output here]
```

## Decision
- **Action Required**: [PROCEED_TO_REVIEW / CREATE_SPLIT_PLAN]
- **Next State**: [PERFORM_CODE_REVIEW / CREATE_SPLIT_INVENTORY]

## Recommendations
[If size violation, recommend split approach]
```

## FORBIDDEN ACTIONS:
- ❌ Manual line counting (wc -l, etc.)
- ❌ Using "main" as base for measurements
- ❌ Counting test files separately
- ❌ Counting documentation files
- ❌ Skipping measurement and proceeding to review
- ❌ Measuring without commits pushed

## State Transitions

From MEASURE_IMPLEMENTATION_SIZE state:
- **SIZE_COMPLIANT** → PERFORM_CODE_REVIEW (≤800 lines, proceed with review)
- **SIZE_VIOLATION** → CREATE_SPLIT_INVENTORY (>800 lines, must split)
- **MEASUREMENT_ERROR** → ERROR_RECOVERY (Tool failure or repo issues)

## Success Criteria
- ✅ Used line-counter.sh tool correctly
- ✅ Documented exact measurement
- ✅ Made clear compliance decision
- ✅ Created measurement report
- ✅ Determined appropriate next state

## Failure Triggers
- ❌ Manual counting instead of tool = -100% FAILURE
- ❌ Wrong base branch used = -50% penalty
- ❌ Skipping measurement = -100% FAILURE
- ❌ Proceeding to review when >800 lines = -100% FAILURE