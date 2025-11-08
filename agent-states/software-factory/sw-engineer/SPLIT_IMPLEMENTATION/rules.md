# SW-Engineer - SPLIT_IMPLEMENTATION State Rules

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
## 🔴🔴🔴 SUPREME LAW R355: PRODUCTION READY CODE ONLY 🔴🔴🔴

### EVEN IN SPLITS - ZERO TOLERANCE FOR:
- ❌ **Hardcoded Credentials** - No passwords/tokens in code
- ❌ **Stub Implementations** - Every split must be functional
- ❌ **Mock/Fake Objects** - Real implementations only
- ❌ **Static Values** - Everything configurable
- ❌ **TODO/FIXME Comments** - No incomplete work

### MANDATORY SPLIT CODE CHECK:
```bash
echo "🔴 R355: SPLIT PRODUCTION CODE CHECK"
cd $SPLIT_DIR
# ANY violation = STOP IMMEDIATELY
grep -r "password.*=.*['\"]" --exclude-dir=test --include="*.go" --include="*.py" --include="*.js" && exit 355
grep -r "stub\|mock\|fake" --exclude-dir=test --include="*.go" --include="*.py" --include="*.js" && exit 355
grep -r "TODO\|FIXME" --exclude-dir=test --include="*.go" --include="*.py" --include="*.js" && exit 355
grep -r "not.*implemented" --exclude-dir=test --include="*.go" --include="*.py" --include="*.js" && exit 355
echo "✅ R355: Split production code verified"
```

## 🔴🔴🔴 CRITICAL: SPLIT IMPLEMENTATION REQUIREMENTS 🔴🔴🔴

**YOU ARE IMPLEMENTING A SPLIT OF A TOO-LARGE BRANCH - SPECIAL PROTOCOLS APPLY!**

## Helper Functions

```bash
# Generate metadata file path with R383/R343 compliance
sf_metadata_path() {
    local file_type="$1"  # IMPLEMENTATION-PLAN, CODE-REVIEW-REPORT, etc.
    local phase="$2"
    local wave="$3"
    local effort="$4"
    local timestamp="${5:-$(date +%Y%m%d-%H%M%S)}"

    echo ".software-factory/phase${phase}/wave${wave}/${effort}/${file_type}--${timestamp}.md"
}
```

## State Context
You are in SPLIT_IMPLEMENTATION state because an effort exceeded size limits and must be split into smaller branches.

## 🔴🔴🔴 CRITICAL: R340 Plan Location Tracking 🔴🔴🔴

**RULE R340: Planning File Metadata Tracking (BLOCKING)**
- You MUST read split plan locations from orchestrator-state-v3.json
- NEVER search for plans using `find` or `ls` commands
- The orchestrator tracks ALL split planning files in the state file
- Violation of R340 = Integration delays and failures

## 🔴🔴🔴 PARAMOUNT: Repository Separation (R251 & R309) 🔴🔴🔴

### R251: Universal Repository Separation Law
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R251-REPOSITORY-SEPARATION-LAW.md`
**Criticality**: PARAMOUNT - Automatic -100% failure for violation
**CRITICAL**: Split directories are TARGET repo clones, NOT SF repo!

### R309: Never Create Efforts in SF Repo  
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R309-never-create-efforts-in-sf-repo.md`
**Criticality**: PARAMOUNT - Automatic -100% failure for violation
**CRITICAL**: Each split gets its own TARGET repo clone under /efforts/

**VERIFY SPLIT IS IN TARGET REPO:**
```bash
echo "🔴 R251/R309: Verifying split is in TARGET repo..."
if [ -f "orchestrator-state-v3.json" ] || [ -f ".claude/CLAUDE.md" ]; then
    echo "🔴🔴🔴 FATAL: Split in Software Factory repo!"
    echo "SPLITS MUST BE IN TARGET REPO CLONES!"
    exit 309
fi

if [[ "$(pwd)" != *"/efforts/"* ]] || [[ "$(pwd)" != *"-SPLIT-"* ]]; then
    echo "🔴 FATAL: Not in a proper split directory!"
    echo "Should be: /efforts/phaseX/waveY/effort-SPLIT-XXX/"
    exit 251
fi

echo "✅ Confirmed: Split in TARGET repo clone under /efforts/"
```

## 🔴🔴🔴 SUPREME LAW R359: NEVER DELETE CODE FOR SIZE LIMITS 🔴🔴🔴

**PENALTY: IMMEDIATE TERMINATION (-1000%)**

### WHAT "SPLITTING" ACTUALLY MEANS:
✅ **CORRECT**: Break your NEW work into 800-line pieces
❌ **WRONG**: Delete everything except your 800-line portion

### EXAMPLE OF THE CATASTROPHIC MISTAKE:
```bash
# ❌❌❌ THIS IS WHAT CAUSED THE DISASTER:
# Original effort had 10,000 lines (existing) + 2,000 lines (new)
# Agent thought: "I need to make this fit in 800 lines"
git rm -rf pkg/build/
git rm -rf pkg/cmd/
git rm -rf pkg/controllers/
git rm main.go
git rm LICENSE
# Result: DELETED 9,552 LINES OF APPROVED CODE!

# ✅✅✅ THIS IS WHAT YOU MUST DO:
# Keep ALL existing code, split only YOUR NEW work
# Split 1: Add 800 lines (repo now 10,800 lines total)
# Split 2: Add 800 lines (repo now 11,600 lines total)
# Split 3: Add 400 lines (repo now 12,000 lines total)
```

### MANDATORY CHECK BEFORE ANY SPLIT COMMIT:
```bash
echo "🔴 R359: Checking for code deletion..."
deleted_lines=$(git diff --numstat main..HEAD | awk '{sum+=$2} END {print sum}')
if [ "$deleted_lines" -gt 100 ]; then
    echo "🔴🔴🔴 R359 VIOLATION: Attempting to delete $deleted_lines lines!"
    echo "NEVER DELETE CODE TO MEET SIZE LIMITS!"
    exit 359
fi
echo "✅ R359: No excessive deletions detected"
```

## 🔴🔴🔴 CRITICAL: STRICT SCOPE ADHERENCE (R310 & R314) 🔴🔴🔴

**IMPLEMENT EXACTLY WHAT'S SPECIFIED - NO MORE, NO LESS!**

### THE WORST SPLIT VIOLATION IN HISTORY:
```
🚨 2667% VIOLATION - CATASTROPHIC FAILURE 🚨
Plan said: "Implement 3 files in pkg/builder"
SW did: Implemented 80 files across 26 packages
Result: Complete project failure, total re-implementation required
```

### ACTUAL 3.4X FAILURE FROM TRANSCRIPT:
```
🚨 3.4X OVERRUN - Split-003 Failure 🚨
Plan said: "selected command files" (~650 lines)
SW thought: "I'll implement all the important commands"
Actual: ALL commands + utilities + tests = 2,215 lines
Cause: "Complete the feature" mindset instead of "stay within budget"
```

### YOUR MANDATORY MINDSET CHANGE:

❌ **WRONG MINDSET (Causes 3.4x overruns):**
- "I should complete this feature properly"
- "This seems incomplete, let me add what's missing"
- "Selected probably means the important ones"
- "While I'm here, I'll add this related thing"

✅ **CORRECT MINDSET (Stays within budget):**
- "I will implement EXACTLY what's listed, nothing more"
- "Incomplete is intentional - other splits handle the rest"
- "Selected means ONLY the ones explicitly named"
- "If it's not listed, I won't add it, period"

### YOUR IMPLEMENTATION RULES:
- Read the "DO NOT IMPLEMENT" section FIRST - it's your primary guide
- Count the EXACT number of files/functions/methods before starting
- STOP at the specified boundaries, even mid-implementation
- NEVER add "helpful" extras or "while I'm here" improvements
- Measure progress every 100 lines to catch overruns early

## 🔴🔴🔴 CRITICAL: Split Infrastructure Requirements 🔴🔴🔴

**EACH SPLIT MUST HAVE ITS OWN DIRECTORY AND CLONE!**

### Split Directory Structure:
```
efforts/phase1/wave1/
├── original-effort/                    # DEPRECATED (too large)
├── original-effort-SPLIT-001/          # Split 1 (separate clone)
├── original-effort-SPLIT-002/          # Split 2 (separate clone)
└── original-effort-SPLIT-003/          # Split 3 (separate clone)
```

### Key Requirements:
1. **Separate Directory**: Each split in its own `-SPLIT-00Z` directory
2. **Separate Clone**: Each split directory has own git repository
3. **Clear Naming**: Directory names include `-SPLIT-` suffix
4. **Sequential Branches**: Each split branch based on previous

## 🔴🔴🔴 CRITICAL: Sequential Split Branching 🔴🔴🔴

**SPLITS ARE CREATED SEQUENTIALLY - EACH BASED ON THE PREVIOUS!**

### The Sequential Chain:
```
Original Effort (1200 lines from phase-integration) - TOO LARGE!
    ↓
Split-001 (400 lines from phase-integration)
    ↓
Split-002 (400 lines from split-001) ← NOT from phase-integration!
    ↓
Split-003 (400 lines from split-002) ← NOT from phase-integration!
```

### Why This Matters for You:
1. **Line Counting**: Your split will measure ONLY what you add
2. **Dependencies**: You can use code from previous splits
3. **No Duplicatio**: Don't re-implement what previous splits did
4. **Clean Merging**: Your work builds progressively

### How to Verify:
```bash
# Check what branch you're based on
git log --oneline -1 --format="%B" | grep "from branch:"
# Should show previous split for split-002 and later
```

## 🔴🔴🔴 CRITICAL: SPLIT SIZE MEASUREMENT 🔴🔴🔴

**SPLITS MUST BE MEASURED AGAINST THEIR IMMEDIATE PREDECESSOR!**

### The Critical Measurement Rule (AUTO-DETECTED!):
```bash
# ✅ NEW TOOL - Automatically detects correct base:
# For split-001:
./tools/line-counter.sh split-001
# Tool output: 🎯 Detected base: <original-branch>

# For split-002:
./tools/line-counter.sh split-002
# Tool output: 🎯 Detected base: split-001

# For split-003:
./tools/line-counter.sh split-003
# Tool output: 🎯 Detected base: split-002

# THE TOOL PREVENTS MEASUREMENT ERRORS!
# No more 5,584 vs 280 line mistakes!
```

### Why This Matters:
- Your split should be 400-600 lines MAX
- Measuring against wrong base shows cumulative size
- Reviewers will reject valid splits if measured wrong
- You could be blamed for exceeding size limits unfairly

## Primary Responsibilities

### 1. Verify Split Infrastructure (R204) and File Placement (R326)
```bash
# MANDATORY: Verify orchestrator created split infrastructure
echo "🔍 Verifying split infrastructure..."

# 🔴🔴🔴 CRITICAL PRE-WORK VALIDATION: DETECT NESTED STRUCTURES! 🔴🔴🔴
CURRENT_DIR=$(pwd)
echo "📍 Current directory: $CURRENT_DIR"

# FATAL CHECK: Detect nested effort paths (the catastrophic bug!)
if [[ "$CURRENT_DIR" == *"/efforts/"*"/efforts/"* ]]; then
    echo "🔴🔴🔴 CATASTROPHIC ERROR: NESTED EFFORT STRUCTURE DETECTED!"
    echo "Path contains duplicate '/efforts/' segments!"
    echo "Current: $CURRENT_DIR"
    echo ""
    echo "This is the WRONG infrastructure pattern:"
    echo "  efforts/phase2/wave1/gitea-client/efforts/phase2/wave1/gitea-client-split-001"
    echo ""
    echo "CORRECT pattern should be:"
    echo "  efforts/phase2/wave1/gitea-client-split-001"
    echo ""
    echo "SW ENGINEER REFUSING TO WORK IN CORRUPTED INFRASTRUCTURE!"
    echo "Orchestrator must fix infrastructure setup!"
    exit 1
fi

# Check 1: We're in a properly named split directory
if [[ "$CURRENT_DIR" != *"-split-"* ]] && [[ "$CURRENT_DIR" != *"-SPLIT-"* ]]; then
    echo "❌ FATAL: Not in a split directory!"
    echo "   Current: $CURRENT_DIR"
    echo "   Expected: path containing '-split-XXX' or '-SPLIT-XXX'"
    echo "   Orchestrator must create split directories first per R204"
    exit 1
fi

# 🔴🔴🔴 R326: CRITICAL CHECK - NO SPLIT SUBDIRECTORIES! 🔴🔴🔴
echo "🔴 R326: Checking for illegal split subdirectories..."
if [ -d "split-"* ] 2>/dev/null; then
    echo "🔴🔴🔴 FATAL: Split subdirectory detected!"
    ls -d split-*
    echo "DELETE THESE IMMEDIATELY! Files go in standard directories!"
    exit 326
fi

echo "✅ R326: No split subdirectories - CORRECT"
echo "📁 Files will go in standard directories: pkg/, cmd/, tests/, etc."

# Extract split identifier from directory name
# Could be effort-name-split-001 or effort-name-SPLIT-001
SPLIT_ID=$(basename "$CURRENT_DIR")
echo "📋 Split identifier: $SPLIT_ID"

# Check 2: R340 - Read split plan location from orchestrator-state-v3.json
echo "🔍 R340: Reading split plan location from orchestrator-state-v3.json..."

STATE_FILE="/workspaces/software-factory-template/orchestrator-state-v3.json"
if [ ! -f "$STATE_FILE" ]; then
    STATE_FILE="/home/vscode/software-factory-template/orchestrator-state-v3.json"
fi

if [ ! -f "$STATE_FILE" ]; then
    echo "❌ CRITICAL: Cannot find orchestrator-state-v3.json!"
    exit 1
fi

# Extract split plan path from state file
if command -v jq &> /dev/null; then
    SPLIT_PLAN=$(jq -r ".effort_repo_files.split_plans[\"${SPLIT_ID}\"].file_path" "$STATE_FILE")
elif command -v yq &> /dev/null; then
    SPLIT_PLAN=$(yq ".effort_repo_files.split_plans[\"${SPLIT_ID}\"].file_path" "$STATE_FILE")
else
    echo "❌ CRITICAL: Neither jq nor yq available for parsing state file!"
    exit 1
fi

if [ "$SPLIT_PLAN" = "null" ] || [ -z "$SPLIT_PLAN" ]; then
    echo "❌ R340 VIOLATION: No split plan tracked for '$SPLIT_ID' in orchestrator-state-v3.json!"
    echo ""
    echo "🔴 ORCHESTRATOR ERROR DETECTED!"
    echo "The orchestrator failed to track the split planning file in state."
    echo "This violates R340: Planning File Metadata Tracking"
    echo ""
    echo "Expected entry in orchestrator-state-v3.json:"
    echo '  "planning_files": {'
    echo '    "split_plans": {'
    echo "      \"$SPLIT_ID\": {"
    echo '        "file_path": "/path/to/.software-factory/SPLIT-PLAN--${TIMESTAMP}.md"'
    echo '      }'
    echo '    }'
    echo '  }'
    echo ""
    echo "ORCHESTRATOR MUST:"
    echo "1. Ensure Code Reviewer reports split plan creation"
    echo "2. Update orchestrator-state-v3.json with split plan metadata"
    echo "3. Commit the state file update"
    echo "4. Re-spawn SW Engineer with proper tracking"
    exit 340
fi

# Verify the split plan actually exists
if [ ! -f "$SPLIT_PLAN" ]; then
    echo "❌ CRITICAL: Tracked split plan does not exist at: $SPLIT_PLAN"
    echo "State file references a non-existent plan!"
    exit 1
fi

echo "✅ R340: Found tracked split plan at: $SPLIT_PLAN"

# Check 3: This is a separate git repository
if [[ ! -d ".git" ]]; then
    echo "❌ FATAL: This is not a git repository!"
    echo "   Each split must have its own clone per R204"
    exit 1
fi

# Check 4: We're on a split branch (canonical format: --split- lowercase)
CURRENT_BRANCH=$(git branch --show-current)
if [[ "$CURRENT_BRANCH" != *"--split-${SPLIT_NUM}"* ]]; then
    echo "❌ FATAL: Not on correct split branch!"
    echo "   Current branch: $CURRENT_BRANCH"
    echo "   Expected: branch containing '--split-${SPLIT_NUM}' (double hyphen, lowercase)"
    echo "   Per R204: Branch format should be {original-branch}--split-{number}"
    exit 1
fi

echo "✅ Split infrastructure verified:"
echo "   Directory: $CURRENT_DIR"
echo "   Split number: $SPLIT_NUM"
echo "   Branch: $CURRENT_BRANCH"
echo "   Plan file: $SPLIT_PLAN"
```

### 2. Sequential Split Implementation (R202, R205)
- **R202**: Single SW engineer implements ALL splits for an effort
- **R205**: Splits must be implemented SEQUENTIALLY, not in parallel
- Complete split-001, then split-002, then split-003, etc.

### 3. Size Compliance (R007, R200)
- Each split MUST stay under 800 lines (HARD LIMIT)
- Use line-counter.sh with BRANCH NAMES (not directory names!)
- Remember: Splits are in SEPARATE git repositories
- Stop immediately if approaching limit

### 4. 🚨🚨🚨 UPDATE SPLIT TRACKING (R302) 🚨🚨🚨

**MANDATORY: Report split details for comprehensive tracking:**
```bash
# After completing each split, report to orchestrator:
SPLIT_REPORT="Split $CURRENT_SPLIT of $TOTAL_SPLITS complete
Branch: $SPLIT_BRANCH
Lines: $(../../tools/line-counter.sh | grep Total | awk '{print $NF}')
Description: [What this split contains]
Status: COMPLETED"

echo "$SPLIT_REPORT" > /tmp/split-${EFFORT_NAME}-${CURRENT_SPLIT}.status
```

### 5. 🚨🚨🚨 NOTIFY ORCHESTRATOR WHEN ALL SPLITS COMPLETE (R296, R302) 🚨🚨🚨

**CRITICAL: After completing the FINAL split for an effort:**

```bash
# When you complete the LAST split (e.g., split-003 of 3 total)
if [ "$CURRENT_SPLIT" == "$TOTAL_SPLITS" ]; then
    echo "═══════════════════════════════════════════════════════════════"
    echo "✅ ALL SPLITS COMPLETE FOR EFFORT: $EFFORT_NAME"
    echo "═══════════════════════════════════════════════════════════════"
    
    # Create completion marker for orchestrator
    COMPLETION_FILE="/tmp/splits-complete-${EFFORT_NAME}.marker"
    cat > "$COMPLETION_FILE" << EOF
EFFORT: $EFFORT_NAME
TOTAL_SPLITS: $TOTAL_SPLITS
COMPLETED_AT: $(date -u +%Y-%m-%dT%H:%M:%SZ)
SPLITS_BRANCHES:
$(for i in $(seq 1 $TOTAL_SPLITS); do
    echo "  - ${EFFORT_NAME}--split-$(printf "%03d" $i)"
done)
STATUS: READY_FOR_DEPRECATION
ORIGINAL_BRANCH: $ORIGINAL_BRANCH
ACTION_REQUIRED: Orchestrator must mark original branch as deprecated per R296
EOF
    
    echo "📋 Completion marker created at: $COMPLETION_FILE"
    echo ""
    echo "🚨 ORCHESTRATOR ACTION REQUIRED:"
    echo "  1. Mark original branch as DEPRECATED per R296"
    echo "  2. Update state file with SPLIT_DEPRECATED status"
    echo "  3. List replacement splits in state file"
    echo ""
    echo "Original branch to deprecate: $ORIGINAL_BRANCH"
    echo "Replacement splits to use for integration:"
    for i in $(seq 1 $TOTAL_SPLITS); do
        echo "  - ${EFFORT_NAME}--split-$(printf "%03d" $i)"
    done
fi
```

## Required Actions for Each Split

### 1. Read Split Plan (R340 Compliant)
```bash
# R340: Read split plan location from orchestrator-state-v3.json
# We already have SPLIT_PLAN from the infrastructure verification step above
# which reads it from orchestrator-state-v3.json per R340 requirements

if [ -z "$SPLIT_PLAN" ]; then
    echo "❌ ERROR: No split plan tracked in orchestrator-state-v3.json!"
    echo "This violates R340 - orchestrator must track all planning files"
    exit 340
fi

echo "Reading split plan (from R340 tracking): $SPLIT_PLAN"
cat "$SPLIT_PLAN"
```

### 2. Create Work Log with Proper Naming (R301)

**MANDATORY: Create timestamped work log to prevent collisions:**
```bash
# R301: Create properly named work log
WORKLOG_NAME="worklog-${EFFORT_NAME}-SPLIT-${SPLIT_NUM}-$(date +%Y%m%d-%H%M%S).txt"
echo "Creating work log: $WORKLOG_NAME"

cat > "$WORKLOG_NAME" << EOF
# Work Log - ${EFFORT_NAME} Split ${SPLIT_NUM}
Date: $(date)
Branch: $(git branch --show-current)
Split Plan: ${SPLIT_PLAN}

## Scope from Plan:
$(grep -A 5 "EXACTLY" "${SPLIT_PLAN}" | head -10)

## Implementation Log:
[Record your work here]
EOF

echo "✅ Work log created: $WORKLOG_NAME"
```

### 3. Implement Only What's in Plan (R310 & R314 ENFORCEMENT)

**🔴🔴🔴 MANDATORY PRE-IMPLEMENTATION CHECKLIST (NEW) 🔴🔴🔴**

```bash
# STEP 0: COMPLETE PRE-IMPLEMENTATION CHECKLIST (MANDATORY AFTER 2667% INCIDENT)
echo "📋 R314: Completing MANDATORY pre-implementation checklist..."
cp $CLAUDE_PROJECT_DIR/templates/SW-ENGINEER-PRE-IMPLEMENTATION-CHECKLIST.md ./pre-impl-checklist-$(date +%s).md

# You MUST fill this out completely before writing ANY code
# This prevents catastrophic violations like the 2667% incident
echo "⚠️ Fill out the checklist NOW - do not skip this!"
```

**BEFORE WRITING ANY CODE:**
```bash
# MANDATORY: Validate file count per R314
echo "🔴 R314: MANDATORY FILE COUNT VALIDATION"
FILE_COUNT=$(grep -c "^###.*File:" "${SPLIT_PLAN:-.software-factory/SPLIT-PLAN--*.md}" 2>/dev/null || echo 0)
if [ -z "$FILE_COUNT" ] || [ "$FILE_COUNT" -eq 0 ]; then
    echo "🚨 FATAL: No file count in split plan!"
    echo "Cannot proceed per R314"
    exit 314
fi

echo "📊 FILE COUNT CONTRACT:"
echo "  You will implement EXACTLY $FILE_COUNT files"
echo "  >200% = AUTOMATIC FAILURE (like the 2667% incident)"
echo "$FILE_COUNT" > .planned-file-count
touch .implementation-start-marker

# MANDATORY: Acknowledge scope boundaries
echo "╔═════════════════════════════════════════════════════════"
echo "🛑 R310 SCOPE BOUNDARIES - MANDATORY ACKNOWLEDGMENT"
echo "═════════════════════════════════════════════════════════"

# Extract DO NOT list
grep -A 10 "DO NOT IMPLEMENT\|STOP BOUNDARIES" "${SPLIT_PLAN:-.software-factory/SPLIT-PLAN--*.md}"

# Count exact scope
FUNC_COUNT=$(grep -i "EXACTLY.*functions" "${SPLIT_PLAN:-.software-factory/SPLIT-PLAN--*.md}" | grep -o '[0-9]+' | head -1)
METHOD_COUNT=$(grep -i "EXACTLY.*methods" "${SPLIT_PLAN:-.software-factory/SPLIT-PLAN--*.md}" | grep -o '[0-9]+' | head -1)

echo "📋 SCOPE FOR THIS SPLIT:"
echo "  Functions to implement: ${FUNC_COUNT:-0} (NO MORE)"
echo "  Methods to implement: ${METHOD_COUNT:-0} (NO MORE)"
echo "  ❌ I will NOT add validation unless explicitly listed"
echo "  ❌ I will NOT add Clone/Copy methods unless requested"
echo "  ❌ I will NOT write comprehensive tests"
echo "  ✅ I will STOP at the specified boundaries"
```

**ADHERENCE RULES:**
- Follow the split plan EXACTLY as written
- Do NOT add extra features (even if they seem necessary)
- Do NOT refactor beyond the plan
- Do NOT "complete" implementations
- If the plan seems insufficient, ASK rather than assume
- Stay within size limits

### 3. Validate Scope Adherence (R310 & R314) and Measure

**INCREMENTAL MONITORING_SWE_PROGRESS (MANDATORY EVERY 100 LINES):**
```bash
# R314: Check file count every 100 lines to prevent 2667% violations
monitor_file_count_incremental() {
    PLANNED=$(cat .planned-file-count 2>/dev/null || echo 0)
    ACTUAL=$(find . -name "*.go" -newer .implementation-start-marker 2>/dev/null | wc -l)
    
    if [ "$PLANNED" -eq 0 ]; then
        echo "❌ No planned file count found - run validation first!"
        exit 314
    fi
    
    RATIO=$((ACTUAL * 100 / PLANNED))
    echo "📊 FILE COUNT MONITOR (Every 100 lines):"
    echo "  Planned: $PLANNED files"
    echo "  Current: $ACTUAL files" 
    echo "  Ratio: $RATIO%"
    
    if [ $RATIO -gt 200 ]; then
        echo "🚨🚨🚨 CATASTROPHIC FAILURE - EXCEEDING 2X FILE LIMIT!"
        echo "You're at $RATIO% - approaching the 2667% violation!"
        echo "STOP IMMEDIATELY!"
        exit 314
    elif [ $RATIO -gt 150 ]; then
        echo "⚠️ CRITICAL WARNING: File count at $RATIO% of plan!"
    fi
}

# Run this check frequently
LINES_ADDED=$(git diff --stat 2>/dev/null | tail -1 | awk '{print $4}')
if [ "${LINES_ADDED:-0}" -gt 100 ]; then
    monitor_file_count_incremental
fi
```

**SCOPE VALIDATION:**
```bash
# FIRST: Validate file count (R314 - HIGHEST PRIORITY)
ACTUAL_FILES=$(find . -name "*.go" -newer .implementation-start-marker 2>/dev/null | wc -l)
PLANNED_FILES=$(cat .planned-file-count 2>/dev/null || grep -c "^###.*File:" "${SPLIT_PLAN:-.software-factory/SPLIT-PLAN--*.md}" 2>/dev/null)

if [ "$ACTUAL_FILES" -gt "$PLANNED_FILES" ]; then
    echo "❌ R314 VIOLATION: $ACTUAL_FILES files but plan specifies $PLANNED_FILES!"
    echo "This is a $(($ACTUAL_FILES * 100 / $PLANNED_FILES))% violation!"
    echo "Remove extra files before proceeding!"
    exit 314
fi

# THEN: Validate function scope (R310)
ACTUAL_FUNCS=$(grep -c "^func [A-Z]" *.go 2>/dev/null || echo 0)
PLANNED_FUNCS=$(grep -i "EXACTLY.*functions" "${SPLIT_PLAN:-.software-factory/SPLIT-PLAN--*.md}" | grep -o '[0-9]+' | head -1)

if [ "$ACTUAL_FUNCS" -gt "${PLANNED_FUNCS:-999}" ]; then
    echo "❌ R310 VIOLATION: Implemented $ACTUAL_FUNCS functions but plan specifies $PLANNED_FUNCS!"
    echo "Remove extra functions before proceeding!"
    exit 1
fi

# THEN: Measure size
CURRENT_BRANCH=$(git branch --show-current)
# Extract split number from canonical format: --split-NNN
SPLIT_NUM=$(echo "$CURRENT_BRANCH" | grep -o '\-\-split-[0-9]*' | grep -o '[0-9]*')

# SEQUENTIAL BRANCHING: Each split measured against PREVIOUS split
if [ "$SPLIT_NUM" = "001" ]; then
    # First split: measure against original base (e.g., phase-integration)
    BASE_BRANCH="phase1-integration"  # Same base as original effort
else
    # Subsequent splits: measure against PREVIOUS split
    PREV_NUM=$(printf "%03d" $((10#$SPLIT_NUM - 1)))
    # Replace --split-NNN with --split-PPP in branch name
    BASE_BRANCH="${CURRENT_BRANCH/--split-${SPLIT_NUM}/--split-${PREV_NUM}}"
fi

echo "📊 Measuring split-${SPLIT_NUM} against: $BASE_BRANCH"
echo "   (Each split measured against previous, NOT original base)"

# Use line-counter.sh - it auto-detects the base!
$CLAUDE_PROJECT_DIR/tools/line-counter.sh $CURRENT_BRANCH

# Get the actual line count (tool auto-detects base)
LINES=$($CLAUDE_PROJECT_DIR/tools/line-counter.sh $CURRENT_BRANCH | grep Total | awk '{print $NF}')
if [ "$LINES" -gt 800 ]; then
    echo "❌ CRITICAL: Split exceeds 800 lines!"
    echo "   This split adds $LINES lines to $BASE_BRANCH"
    exit 1
fi
```

### 4. Create Individual Split Completion Marker (MANDATORY)
```bash
# 🔴🔴🔴 MANDATORY: Create split completion marker 🔴🔴🔴
echo "📋 Creating MANDATORY split completion marker..."
SPLIT_MARKER="SPLIT-${SPLIT_NUM}-COMPLETE.marker"
cat > "$SPLIT_MARKER" << EOF
Split Number: ${SPLIT_NUM}
Completed at: $(date '+%Y-%m-%d %H:%M:%S %Z')
Effort: ${EFFORT_NAME}
Branch: $(git branch --show-current)
Total lines: $(./tools/line-counter.sh | grep Total | awk '{print $NF}') lines
Final commit: $(git log --oneline -1)
Status: SPLIT ${SPLIT_NUM} COMPLETE
EOF

echo "✅ Split completion marker created: $SPLIT_MARKER"
```

### 5. Commit and Push
```bash
# 🔴🔴🔴 R326: FINAL CHECK BEFORE COMMIT 🔴🔴🔴
echo "🔴 R326: Final check for split subdirectories..."
if find . -type d -name "split-*" | grep -q .; then
    echo "🔴🔴🔴 FATAL: Split subdirectories found!"
    find . -type d -name "split-*"
    echo "These cause CATASTROPHIC measurement errors!"
    echo "DELETE all split-XXX/ directories and move files to standard locations!"
    exit 326
fi

# Verify files are in correct locations
echo "✅ R326: Verifying files in standard directories..."
for dir in pkg cmd tests docs; do
    if [ -d "$dir" ]; then
        echo "  ✅ $dir/ contains:"
        ls -la "$dir" | head -5
    fi
done

# Commit with clear message INCLUDING THE MARKER
git add -A
git add "$SPLIT_MARKER"
git commit -m "feat: implement split-${SPLIT_NUM} - ${DESCRIPTION}

Includes mandatory SPLIT-${SPLIT_NUM}-COMPLETE.marker for orchestrator monitoring"
git push

# Validation check
if [ ! -f "$SPLIT_MARKER" ]; then
    echo "🔴 ERROR: Split marker not created! Work is NOT complete!"
    exit 1
fi
echo "✅ Split ${SPLIT_NUM} complete with marker"
```

## State Transitions

### After Each Split:
- If more splits remain → Stay in SPLIT_IMPLEMENTATION
- If all splits complete → Transition to COMPLETE
- If size exceeded → Transition to ERROR

### Success Criteria:
- ✅ All splits implemented sequentially
- ✅ Each split under 800 lines
- ✅ All splits pushed to remote
- ✅ Orchestrator notified when complete (R296)


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


## Related Rules
- **R310**: Split scope strict adherence protocol (MANDATORY)
- **R314**: Mandatory file count validation protocol (NEW - CRITICAL)
- **R202**: Single agent per split effort
- **R204**: Split infrastructure creation
- **R205**: Sequential split navigation
- **R207**: Split boundary validation
- **R296**: Deprecated branch marking protocol
- **R007**: Size limit compliance
- **R200**: Measure only changeset

## Common Violations to Avoid

### ❌ Implementing Splits in Parallel
```bash
# WRONG - Opening multiple terminals
terminal1: implement split-001
terminal2: implement split-002  # VIOLATION!
```

### ❌ Exceeding Size Limits
```bash
# WRONG - Adding extra features
"While I'm here, let me also refactor..."  # NO!
```

### ❌ Over-Engineering (R310 Violation)
```bash
# WRONG - Making it "complete"
Plan: "Implement WithImage, WithContext, WithPlatform"
Doing: Implementing 15 option functions for "completeness"  # VIOLATION!
```

### ❌ Ignoring DO NOT Instructions
```bash
# WRONG - Adding forbidden features
Plan: "DO NOT add validation"
Doing: "This needs validation to be production ready"  # VIOLATION!
```

### ❌ Not Notifying Orchestrator
```bash
# WRONG - Completing without notification
"Split 3 done, my work is complete" [exits]  # MUST NOTIFY!
```

### ✅ Correct Pattern
```bash
1. Implement split-001 → measure → commit → push
2. Implement split-002 → measure → commit → push  
3. Implement split-003 → measure → commit → push
4. Notify orchestrator that ALL splits are complete
5. Orchestrator marks original branch as deprecated
```

## Acknowledgment Required
After reading these rules, acknowledge:
"✅ Successfully read SPLIT_IMPLEMENTATION rules for sw-engineer. I understand:
- R310: I must implement EXACTLY what's specified, no more
- R296: I must notify the orchestrator when ALL splits are complete
- I will NOT add extra features or 'complete' implementations
- I will STOP at specified boundaries even if code seems incomplete"
