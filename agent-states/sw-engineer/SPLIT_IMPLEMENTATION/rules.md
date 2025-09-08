# SW-Engineer - SPLIT_IMPLEMENTATION State Rules

## 🔴🔴🔴 CRITICAL: SPLIT IMPLEMENTATION REQUIREMENTS 🔴🔴🔴

**YOU ARE IMPLEMENTING A SPLIT OF A TOO-LARGE BRANCH - SPECIAL PROTOCOLS APPLY!**

## State Context
You are in SPLIT_IMPLEMENTATION state because an effort exceeded size limits and must be split into smaller branches.

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
if [ -f "orchestrator-state.json" ] || [ -f ".claude/CLAUDE.md" ]; then
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

# Extract split number from directory name
SPLIT_NUM=$(echo "$CURRENT_DIR" | grep -o 'SPLIT-[0-9]*' | grep -o '[0-9]*')
if [ -z "$SPLIT_NUM" ]; then
    echo "❌ FATAL: Cannot determine split number from directory name!"
    exit 1
fi

# Check 2: Find the timestamped split plan for this split
# First check new .software-factory structure
if [ -d ".software-factory" ]; then
    SPLIT_PLAN=$(find .software-factory -name "SPLIT-PLAN-*.md" -type f | sort -r | head -1)
    if [ -n "$SPLIT_PLAN" ]; then
        echo "✅ Found split plan in .software-factory: $SPLIT_PLAN"
    fi
fi

# If not found in new location, check old location
if [ -z "$SPLIT_PLAN" ]; then
    # Look for the most recent split plan matching this split number
    SPLIT_PLAN=$(ls -t SPLIT-PLAN-*-split${SPLIT_NUM}-*.md 2>/dev/null | head -1)
    if [ -n "$SPLIT_PLAN" ]; then
        echo "⚠️ Found split plan in legacy location: $SPLIT_PLAN"
    fi
fi

if [ -z "$SPLIT_PLAN" ]; then
    # Fallback: Check for legacy numbered format (backwards compatibility)
    if [[ -f "SPLIT-PLAN-${SPLIT_NUM}.md" ]]; then
        echo "⚠️ WARNING: Using legacy split plan format (should be timestamped per R301)"
        SPLIT_PLAN="SPLIT-PLAN-${SPLIT_NUM}.md"
    else
        echo "❌ FATAL: No split plan found for split ${SPLIT_NUM}!"
        echo "   Searched in: .software-factory/ (new location)"
        echo "   Searched for: SPLIT-PLAN-*-split${SPLIT_NUM}-*.md (old location)"
        echo "   Also checked legacy: SPLIT-PLAN-${SPLIT_NUM}.md"
        echo "   Orchestrator must ensure split plans exist"
        exit 1
    fi
fi

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

### 1. Read Split Plan
```bash
# Find split plan - check new location first
if [ -d ".software-factory" ]; then
    SPLIT_PLAN=$(find .software-factory -name "SPLIT-PLAN-*.md" -type f | sort -r | head -1)
fi
# Fallback to old location
if [ -z "$SPLIT_PLAN" ]; then
    SPLIT_PLAN=$(ls SPLIT-PLAN-*.md 2>/dev/null | head -1)
fi

if [ -z "$SPLIT_PLAN" ]; then
    echo "❌ ERROR: No split plan found!"
    exit 1
fi

echo "Reading split plan: $SPLIT_PLAN"
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
FILE_COUNT=$(grep -c "^###.*File:" "${SPLIT_PLAN:-SPLIT-PLAN.md}" 2>/dev/null || echo 0)
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
grep -A 10 "DO NOT IMPLEMENT\|STOP BOUNDARIES" "${SPLIT_PLAN:-SPLIT-PLAN.md}"

# Count exact scope
FUNC_COUNT=$(grep -i "EXACTLY.*functions" "${SPLIT_PLAN:-SPLIT-PLAN.md}" | grep -o '[0-9]+' | head -1)
METHOD_COUNT=$(grep -i "EXACTLY.*methods" "${SPLIT_PLAN:-SPLIT-PLAN.md}" | grep -o '[0-9]+' | head -1)

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

**INCREMENTAL MONITORING (MANDATORY EVERY 100 LINES):**
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
PLANNED_FILES=$(cat .planned-file-count 2>/dev/null || grep -c "^###.*File:" "${SPLIT_PLAN:-SPLIT-PLAN.md}" 2>/dev/null)

if [ "$ACTUAL_FILES" -gt "$PLANNED_FILES" ]; then
    echo "❌ R314 VIOLATION: $ACTUAL_FILES files but plan specifies $PLANNED_FILES!"
    echo "This is a $(($ACTUAL_FILES * 100 / $PLANNED_FILES))% violation!"
    echo "Remove extra files before proceeding!"
    exit 314
fi

# THEN: Validate function scope (R310)
ACTUAL_FUNCS=$(grep -c "^func [A-Z]" *.go 2>/dev/null || echo 0)
PLANNED_FUNCS=$(grep -i "EXACTLY.*functions" "${SPLIT_PLAN:-SPLIT-PLAN.md}" | grep -o '[0-9]+' | head -1)

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

### 4. Commit and Push
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

# Commit with clear message
git add -A
git commit -m "feat: implement split-${SPLIT_NUM} - ${DESCRIPTION}"
git push
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