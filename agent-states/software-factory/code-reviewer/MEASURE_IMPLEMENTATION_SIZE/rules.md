# Code Reviewer - MEASURE_IMPLEMENTATION_SIZE State Rules

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

### ⚠️⚠️⚠️ CRITICAL: LINE-COUNTER.SH AUTO-DETECTS BASE - NO PARAMETERS! ⚠️⚠️⚠️

**🔴🔴🔴 TOOL UPDATE: AUTO-DETECTION IS NOW MANDATORY! 🔴🔴🔴**

### THE TOOL IS SMART - IT KNOWS THE CORRECT BASE:
1. **ALWAYS use ${PROJECT_ROOT}/tools/line-counter.sh** - NO EXCEPTIONS
2. **NO PARAMETERS NEEDED** - Tool auto-detects EVERYTHING!
3. **NEVER specify -b parameter** - That's OLD/WRONG syntax!
4. **NEVER do manual counting** - AUTOMATIC FAILURE (-100%)
5. **Tool shows detected base** - Look for "🎯 Detected base:" in output

### HOW IT WORKS:
- **For efforts**: Detects phase/wave pattern, uses correct integration branch
- **For splits**: Knows to measure split-001 vs original, split-002 vs split-001
- **For integrations**: Knows to measure against main
- **Shows its work**: Output includes "🎯 Detected base: [branch]"

### CORRECT USAGE (UPDATED FOR AUTO-DETECTION):
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

# STEP 3: Find project root (where orchestrator-state-v3.json lives)
PROJECT_ROOT=$(pwd); while [ "$PROJECT_ROOT" != "/" ]; do 
    [ -f "$PROJECT_ROOT/orchestrator-state-v3.json" ] && break; 
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT"); 
done
echo "Project root: $PROJECT_ROOT"

# STEP 4: RUN THE TOOL - NO PARAMETERS AT ALL!
$PROJECT_ROOT/tools/line-counter.sh
# That's it! The tool does EVERYTHING automatically!

# Tool output will show:
# 🎯 Detected base: phase1/integration (or appropriate base)
# 📦 Analyzing branch: phase1/wave1/my-effort
# ✅ Total non-generated lines: 487
```

### 🔴🔴🔴 CRITICAL: Just Let The Tool Auto-Detect! 🔴🔴🔴

**THE FATAL MISTAKES TO AVOID:**
```bash
# ❌❌❌ WRONG - Trying to specify base manually:
cd efforts/phase2/wave1/go-containerregistry-image-builder
./line-counter.sh -b main  # WRONG! No -b parameter!

# ❌❌❌ WRONG - Using git diff with wrong base:
git diff main --stat  # WRONG! This will count ALL code since main!

# ❌❌❌ WRONG - Manual counting:
find . -name "*.go" | xargs wc -l  # WRONG! Manual counting forbidden!

# ✅✅✅ CORRECT - Just run the tool, no parameters:
cd efforts/phase2/wave1/go-containerregistry-image-builder
../../tools/line-counter.sh  # CORRECT! Tool auto-detects everything!
```

### THE TOOL HANDLES SPLITS AUTOMATICALLY:
```bash
# 🔴🔴🔴 YOU DON'T NEED TO FIGURE OUT THE BASE! 🔴🔴🔴

# The tool AUTOMATICALLY detects split patterns and uses correct base:
# - split-001: Measures against the original effort branch
# - split-002: Measures against split-001
# - split-003: Measures against split-002
# etc.

# Just run:
$PROJECT_ROOT/tools/line-counter.sh

# Tool output for splits will show:
# 🎯 Detected base: phase1/wave1/my-effort (for split-001)
# 🎯 Detected base: phase1/wave1/my-effort--split-001 (for split-002)
# etc.

# YOU DON'T NEED TO CALCULATE THIS - THE TOOL DOES IT!
```

## Size Compliance Decision Tree

```python
def determine_size_action(line_count):
    """Determine action based on measured size"""
    
    if line_count <= 800:
        return {
            'compliant': True,
            'action': 'PROCEED_TO_REVIEW',
            'next_state': 'CODE_REVIEW',
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
- **Next State**: [CODE_REVIEW / CREATE_SPLIT_INVENTORY]

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
- **SIZE_COMPLIANT** → CODE_REVIEW (≤800 lines, proceed with review)
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

