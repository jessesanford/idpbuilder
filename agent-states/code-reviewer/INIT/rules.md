# Code Reviewer - INIT State Rules

## State Context
You are initializing as a Code Reviewer. Your FIRST priority is establishing directory isolation and understanding your task.

## 🔴🔴🔴 MANDATORY STARTUP SEQUENCE 🔴🔴🔴

**YOU MUST DO THESE IN ORDER:**

### Step 0: CAPTURE ORCHESTRATOR'S PROMPT
```bash
# Store the entire prompt you received for error reporting
ORCHESTRATOR_PROMPT="[Store the complete prompt/instructions you received from the orchestrator here]"
```

### Step 1: DETERMINE YOUR TASK TYPE
```bash
echo "═══════════════════════════════════════════════════════"
echo "🚨 CODE REVIEWER STARTUP - DETERMINING TASK TYPE 🚨"
echo "═══════════════════════════════════════════════════════"
echo ""
echo "📍 INITIAL DIRECTORY: $(pwd)"
echo ""

# Determine what type of review/planning we're doing
if echo "$ORCHESTRATOR_PROMPT" | grep -qi "create.*implementation.*plan\|effort.*planning\|plan.*creation"; then
    TASK_TYPE="EFFORT_PLANNING"
    echo "📋 Task Type: Creating Implementation Plan"
    EXPECTED_FILES="work-log.md"
elif echo "$ORCHESTRATOR_PROMPT" | grep -qi "review.*code\|validate.*implementation\|check.*compliance"; then
    TASK_TYPE="CODE_REVIEW"
    echo "📋 Task Type: Reviewing Implementation"
    EXPECTED_FILES="IMPLEMENTATION-PLAN.md work-log.md"
elif echo "$ORCHESTRATOR_PROMPT" | grep -qi "split.*plan\|plan.*split\|too.*large"; then
    TASK_TYPE="SPLIT_PLANNING"
    echo "📋 Task Type: Planning Split Strategy"
    EXPECTED_FILES="IMPLEMENTATION-PLAN.md work-log.md"
elif echo "$ORCHESTRATOR_PROMPT" | grep -qi "wave.*plan\|phase.*plan"; then
    TASK_TYPE="WAVE_PLANNING"
    echo "📋 Task Type: Creating Wave/Phase Plan"
    EXPECTED_FILES=""
else
    TASK_TYPE="UNKNOWN"
    echo "⚠️ Unable to determine task type from prompt"
fi
```

### Step 2: FIND AND NAVIGATE TO WORKING DIRECTORY
```bash
echo ""
echo "🔍 FINDING WORKING DIRECTORY"
echo ""

# Extract directory from orchestrator's prompt
PROMPT_DIR=$(echo "$ORCHESTRATOR_PROMPT" | grep -oE "(TARGET_DIRECTORY|WORKING_DIR|Directory):[[:space:]]*/[^[:space:]]+" | head -1 | cut -d: -f2 | xargs)
if [ -n "$PROMPT_DIR" ]; then
    echo "📋 Orchestrator specified directory: $PROMPT_DIR"
    TARGET_DIR="$PROMPT_DIR"
elif echo "$ORCHESTRATOR_PROMPT" | grep -oE "efforts/phase[0-9]+/wave[0-9]+/[^[:space:]]+" | head -1; then
    # Extract from path mentioned in prompt
    TARGET_DIR="/$(echo "$ORCHESTRATOR_PROMPT" | grep -oE "efforts/phase[0-9]+/wave[0-9]+/[^[:space:]]+" | head -1)"
    echo "📋 Extracted directory from prompt: $TARGET_DIR"
else
    echo "⚠️ No directory specified in prompt, will search..."
    TARGET_DIR=""
fi

# Try to navigate to specified directory first
if [ -n "$TARGET_DIR" ]; then
    if [ -d "$TARGET_DIR" ]; then
        cd "$TARGET_DIR"
        echo "✅ Successfully navigated to: $(pwd)"
    else
        echo "❌ Specified directory doesn't exist: $TARGET_DIR"
        TARGET_DIR=""
    fi
fi

# If we couldn't navigate from prompt, search for appropriate directory
if [ -z "$TARGET_DIR" ] || [ ! -d "$(pwd)" ]; then
    echo "🔍 Searching for appropriate directory..."
    
    # For effort planning, look for directories with work-log.md but no IMPLEMENTATION-PLAN.md
    if [ "$TASK_TYPE" = "EFFORT_PLANNING" ]; then
        echo "Looking for effort directory ready for planning..."
        EFFORT_DIRS=$(find /workspaces -type d -path "*/efforts/phase*/wave*/*" -maxdepth 7 2>/dev/null)
        
        for dir in $EFFORT_DIRS; do
            if [ -f "$dir/work-log.md" ] && [ ! -f "$dir/IMPLEMENTATION-PLAN.md" ]; then
                echo "✅ Found effort directory ready for planning: $dir"
                cd "$dir"
                TARGET_DIR="$dir"
                break
            fi
        done
    
    # For code review, look for directories with IMPLEMENTATION-PLAN.md
    elif [ "$TASK_TYPE" = "CODE_REVIEW" ] || [ "$TASK_TYPE" = "SPLIT_PLANNING" ]; then
        echo "Looking for effort directory with implementation to review..."
        EFFORT_DIRS=$(find /workspaces -type d -path "*/efforts/phase*/wave*/*" -maxdepth 7 2>/dev/null)
        
        for dir in $EFFORT_DIRS; do
            if [ -f "$dir/IMPLEMENTATION-PLAN.md" ]; then
                echo "✅ Found effort directory with implementation: $dir"
                cd "$dir"
                TARGET_DIR="$dir"
                break
            fi
        done
    fi
fi

# Final verification
echo ""
echo "📂 DIRECTORY VERIFICATION"
echo "Current directory: $(pwd)"

# Check if we're in an appropriate directory
if [[ "$(pwd)" == *"/efforts/phase"*/wave*/* ]] || [[ "$(pwd)" == *"/phase-plans"* ]]; then
    echo "✅ In a valid working directory"
    
    # Verify expected files based on task type
    if [ -n "$EXPECTED_FILES" ]; then
        for file in $EXPECTED_FILES; do
            if [ -f "$file" ]; then
                echo "✅ Found expected file: $file"
            else
                echo "⚠️ Missing expected file: $file"
            fi
        done
    fi
else
    echo "❌ ENVIRONMENT ERROR: Not in a valid effort or planning directory"
    echo ""
    echo "🔴 ORCHESTRATOR, YOU GAVE ME THE WRONG PROMPT!"
    echo ""
    echo "THIS IS THE PROMPT YOU GAVE:"
    echo "════════════════════════════════════════"
    echo "$ORCHESTRATOR_PROMPT"
    echo "════════════════════════════════════════"
    echo ""
    echo "I FAILED TO FIND MY WORKING DIRECTORY BASED ON THIS PROMPT."
    echo ""
    echo "TASK TYPE DETECTED: $TASK_TYPE"
    echo ""
    echo "WHAT I EXPECTED:"
    echo "- Clear TARGET_DIRECTORY or WORKING_DIR specification"
    echo "- Directory path: /efforts/phase{X}/wave{Y}/{effort-name}"
    echo "- Or phase plan directory: /phase-plans/"
    echo ""
    echo "WHAT I FOUND:"
    echo "- Current directory: $(pwd)"
    echo "- Directory from prompt: ${PROMPT_DIR:-NOT SPECIFIED}"
    echo "- Searched paths: /workspaces/*/efforts/phase*/wave*/*"
    echo ""
    echo "PLEASE TRY AGAIN WITH:"
    echo "1. Explicit TARGET_DIRECTORY: /path/to/effort"
    echo "2. Verification that infrastructure exists"
    echo "3. Clear task specification (planning vs review)"
    echo ""
    echo "GRADING VIOLATION: Cannot proceed without proper infrastructure (R208/R209)"
    exit 254
fi
```

### Step 3: VERIFY BRANCH AND GIT STATE
```bash
echo ""
echo "🌿 VERIFYING GIT STATE"
echo ""

# Get current branch
CURRENT_BRANCH=$(git branch --show-current 2>/dev/null || echo "NONE")
echo "Current branch: $CURRENT_BRANCH"

# Extract expected branch from prompt or files
if [ -f "IMPLEMENTATION-PLAN.md" ]; then
    EXPECTED_BRANCH=$(grep "**BRANCH**:" IMPLEMENTATION-PLAN.md | cut -d: -f2- | xargs)
elif [ -f "work-log.md" ]; then
    EXPECTED_BRANCH=$(grep "Branch:" work-log.md | head -1 | cut -d: -f2- | xargs)
else
    EXPECTED_BRANCH=$(echo "$ORCHESTRATOR_PROMPT" | grep -oE "BRANCH:[[:space:]]*[^[:space:]]+" | cut -d: -f2 | xargs)
fi

if [ -n "$EXPECTED_BRANCH" ]; then
    echo "Expected branch: $EXPECTED_BRANCH"
    
    if [ "$CURRENT_BRANCH" != "$EXPECTED_BRANCH" ]; then
        echo "⚠️ Branch mismatch!"
        echo "Attempting to checkout expected branch..."
        
        if git checkout "$EXPECTED_BRANCH" 2>/dev/null; then
            echo "✅ Successfully switched to: $EXPECTED_BRANCH"
        else
            echo "❌ ENVIRONMENT ERROR: Cannot checkout expected branch"
            echo ""
            echo "🔴 ORCHESTRATOR, BRANCH CONFIGURATION ERROR!"
            echo ""
            echo "YOUR PROMPT:"
            echo "════════════════════════════════════════"
            echo "$ORCHESTRATOR_PROMPT"
            echo "════════════════════════════════════════"
            echo ""
            echo "BRANCH SITUATION:"
            echo "- Current: $CURRENT_BRANCH"
            echo "- Expected: $EXPECTED_BRANCH"
            echo "- Branch exists: NO"
            echo ""
            echo "ORCHESTRATOR, PLEASE:"
            echo "1. Run SETUP_EFFORT_INFRASTRUCTURE to create branches"
            echo "2. Verify branch was pushed to remote"
            echo "3. Re-spawn me after infrastructure is ready"
            echo ""
            echo "GRADING VIOLATION: R193 - Branch infrastructure missing"
            exit 193
        fi
    else
        echo "✅ On correct branch"
    fi
fi

# Check remote tracking
if git branch -vv | grep -q "$CURRENT_BRANCH.*\[origin/"; then
    echo "✅ Branch has remote tracking"
else
    echo "⚠️ No remote tracking configured"
fi
```

### Step 4: ESTABLISH R209 ISOLATION
```bash
echo ""
echo "🔒 ESTABLISHING R209 DIRECTORY ISOLATION"
echo ""

# Set effort isolation directory
export EFFORT_ISOLATION_DIR="$(pwd)"
export readonly EFFORT_ISOLATION_DIR

echo "EFFORT_ISOLATION_DIR set to: $EFFORT_ISOLATION_DIR"
echo "✅ Directory isolation established"

# Verify we can access necessary files
if [ "$TASK_TYPE" = "EFFORT_PLANNING" ]; then
    # Check for wave plan
    WAVE_PLAN_LOCATIONS=(
        "../../PHASE-*-WAVE-*-IMPLEMENTATION-PLAN.md"
        "../../../phase-plans/PHASE-*-WAVE-*-IMPLEMENTATION-PLAN.md"
        "../../../../phase-plans/PHASE-*-WAVE-*-IMPLEMENTATION-PLAN.md"
    )
    
    WAVE_PLAN_FOUND=false
    for pattern in "${WAVE_PLAN_LOCATIONS[@]}"; do
        if ls $pattern 2>/dev/null | head -1; then
            WAVE_PLAN_FOUND=true
            echo "✅ Found wave plan at: $(ls $pattern 2>/dev/null | head -1)"
            break
        fi
    done
    
    if [ "$WAVE_PLAN_FOUND" = false ]; then
        echo "⚠️ Cannot find wave implementation plan"
        echo "Will need to request from orchestrator if needed"
    fi
fi
```

### Step 5: ACKNOWLEDGE STARTUP REQUIREMENTS
```bash
echo ""
echo "═══════════════════════════════════════════════════════"
echo "📋 CODE REVIEWER STARTUP COMPLETE"
echo "═══════════════════════════════════════════════════════"
echo ""
echo "AGENT STARTUP: $(date '+%Y-%m-%d %H:%M:%S %Z')"
echo ""
echo "STARTUP VERIFICATION:"
echo "  Task Type: $TASK_TYPE"
echo "  Directory: $EFFORT_ISOLATION_DIR"
echo "  Branch: $(git branch --show-current)"
echo "  Remote: $(git remote -v | head -1)"
echo ""
echo "ACKNOWLEDGING RULES:"
echo "  ✅ R208 - Directory spawn protocol"
echo "  ✅ R209 - Effort directory isolation"
echo "  ✅ R254 - Error reporting protocol"
echo "  ✅ R203 - State-aware agent startup"
echo "  ✅ R214 - Wave directory acknowledgment"
echo ""
echo "READY TO PROCEED WITH: $TASK_TYPE"
echo "═══════════════════════════════════════════════════════"
```

## Critical Error Conditions

If ANY of these occur, STOP and report to orchestrator:

1. **Cannot determine task type** - Need clearer instructions
2. **Cannot find working directory** - Infrastructure not ready
3. **Expected files missing** - May be in wrong state
4. **Cannot checkout branch** - Git infrastructure issue
5. **No remote tracking** - Push may have failed

## Task Type Specific Requirements

### EFFORT_PLANNING
- Must have work-log.md
- Must NOT have IMPLEMENTATION-PLAN.md yet
- Need access to wave plan

### CODE_REVIEW
- Must have IMPLEMENTATION-PLAN.md
- Must have work-log.md
- Code should be implemented

### SPLIT_PLANNING
- Must have IMPLEMENTATION-PLAN.md showing >800 lines
- Need measurement results

### WAVE_PLANNING
- May be in phase-plans directory
- Need phase requirements

## State Transition

After successful initialization:
1. Verify all startup checks passed
2. Confirm task type identified
3. Transition to appropriate next state:
   - EFFORT_PLANNING → EFFORT_PLAN_CREATION
   - CODE_REVIEW → CODE_REVIEW
   - SPLIT_PLANNING → CREATE_SPLIT_PLAN
   - WAVE_PLANNING → WAVE_IMPLEMENTATION_PLANNING

## Remember

- **ALWAYS** capture orchestrator prompt first
- **ALWAYS** verify directory before proceeding
- **ALWAYS** check git branch and remote
- **NEVER** proceed if infrastructure missing
- **NEVER** guess at task requirements
- **ALWAYS** report detailed errors with prompt echo-back