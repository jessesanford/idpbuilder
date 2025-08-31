# SW Engineer - INIT State Rules

## State Context
You are initializing as a SW Engineer. Your FIRST priority is establishing directory isolation.

## 🔴🔴🔴 MANDATORY STARTUP SEQUENCE 🔴🔴🔴

**YOU MUST DO THESE IN ORDER:**

### Step 0: CAPTURE ORCHESTRATOR'S PROMPT
```bash
# Store the entire prompt you received for error reporting
ORCHESTRATOR_PROMPT="[Store the complete prompt/instructions you received from the orchestrator here]"
```

### Step 1: FIND AND NAVIGATE TO EFFORT DIRECTORY
```bash
echo "═══════════════════════════════════════════════════════"
echo "🚨 SW ENGINEER STARTUP - FINDING EFFORT DIRECTORY 🚨"
echo "═══════════════════════════════════════════════════════"
echo ""
echo "📍 INITIAL DIRECTORY: $(pwd)"
echo ""

# Try to extract directory from orchestrator's prompt
PROMPT_DIR=$(echo "$ORCHESTRATOR_PROMPT" | grep -oE "(TARGET_DIRECTORY|WORKING_DIR|Directory):[[:space:]]*/[^[:space:]]+" | head -1 | cut -d: -f2 | xargs)
if [ -n "$PROMPT_DIR" ]; then
    echo "📋 Orchestrator specified directory: $PROMPT_DIR"
fi

# Check if we're already in the right place (handle both old and new formats)
# Look for timestamped plans first
LATEST_PLAN=$(ls -t IMPLEMENTATION-PLAN-*.md 2>/dev/null | head -n1)

# Fallback to old format if no timestamped versions
if [ -z "$LATEST_PLAN" ] && [ -f "IMPLEMENTATION-PLAN.md" ]; then
    LATEST_PLAN="IMPLEMENTATION-PLAN.md"
    echo "⚠️ Using legacy plan format"
fi

if [ -n "$LATEST_PLAN" ]; then
    echo "✅ Already in effort directory with plan: $LATEST_PLAN"
    # Set global variable for other steps
    export IMPLEMENTATION_PLAN="$LATEST_PLAN"
else
    echo "⚠️ Not in effort directory, searching..."
    
    # Try to find our effort directory
    # First check if orchestrator gave us a hint in the current path
    if [[ "$(pwd)" == *"/efforts/phase"*/wave*/* ]]; then
        echo "📂 Already in an effort path, checking for plan..."
    else
        # Search for effort directories
        echo "🔍 Searching for effort directories with plans..."
        EFFORT_DIRS=$(find /workspaces -type d -path "*/efforts/phase*/wave*/*" -maxdepth 7 2>/dev/null)
        
        for dir in $EFFORT_DIRS; do
            # Check for timestamped plans first
            PLAN_COUNT=$(ls "$dir"/IMPLEMENTATION-PLAN*.md 2>/dev/null | wc -l)
            if [ $PLAN_COUNT -gt 0 ]; then
                echo "✅ Found effort directory with plan(s): $dir"
                cd "$dir"
                # Find the latest plan
                LATEST_PLAN=$(ls -t IMPLEMENTATION-PLAN-*.md 2>/dev/null | head -n1)
                if [ -z "$LATEST_PLAN" ] && [ -f "IMPLEMENTATION-PLAN.md" ]; then
                    LATEST_PLAN="IMPLEMENTATION-PLAN.md"
                fi
                export IMPLEMENTATION_PLAN="$LATEST_PLAN"
                break
            fi
        done
    fi
    
    # Final check with R254 error reporting
    if [ -z "$IMPLEMENTATION_PLAN" ]; then
        echo "❌ ENVIRONMENT ERROR: Cannot find effort directory with IMPLEMENTATION-PLAN*.md"
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
        echo "WHAT I EXPECTED:"
        echo "- Clear TARGET_DIRECTORY: /efforts/phase{X}/wave{Y}/{effort-name}"
        echo "- File: IMPLEMENTATION-PLAN-*.md (or legacy IMPLEMENTATION-PLAN.md) in that directory"
        echo ""
        echo "WHAT I FOUND:"
        echo "- Current directory: $(pwd)"
        echo "- Searched paths: /workspaces/*/efforts/phase*/wave*/*"
        echo "- Plans found: NONE"
        echo "- Directory from prompt: ${PROMPT_DIR:-NOT SPECIFIED}"
        echo ""
        echo "PLEASE TRY AGAIN WITH:"
        echo "1. Correct directory path for my effort"
        echo "2. Verification that infrastructure exists"
        echo "3. Clear TARGET_DIRECTORY specification"
        echo ""
        echo "GRADING VIOLATION: Cannot proceed without proper infrastructure (R208/R209)"
        exit 254
    fi
fi

echo "✅ Now in directory with plan $IMPLEMENTATION_PLAN: $(pwd)"
```

### Step 2: VERIFY CORRECT EFFORT DIRECTORY
```bash
# Extract the required directory from the plan
WORKING_DIR=$(grep "**WORKING_DIRECTORY**:" "$IMPLEMENTATION_PLAN" | cut -d: -f2- | xargs)

if [ -n "$WORKING_DIR" ]; then
    echo "📋 Plan specifies directory: $WORKING_DIR"
    
    # Navigate if we're not there yet
    if [ "$(pwd)" != "$WORKING_DIR" ]; then
        echo "📂 Navigating to correct effort directory..."
        if [ -d "$WORKING_DIR" ]; then
            cd "$WORKING_DIR"
            echo "✅ Successfully navigated to: $(pwd)"
        else
            echo "❌ ENVIRONMENT ERROR: Specified directory doesn't exist"
            echo ""
            echo "🔴 ORCHESTRATOR, YOUR PROMPT LED ME TO A NON-EXISTENT DIRECTORY!"
            echo ""
            echo "YOUR PROMPT SPECIFIED:"
            echo "════════════════════════════════════════"
            echo "$ORCHESTRATOR_PROMPT"
            echo "════════════════════════════════════════"
            echo ""
            echo "EXTRACTED DIRECTORY: $WORKING_DIR"
            echo "THIS DIRECTORY DOES NOT EXIST!"
            echo ""
            echo "CURRENT SITUATION:"
            echo "- Attempted directory: $WORKING_DIR"
            echo "- Directory exists: NO ❌"
            echo "- Current location: $(pwd)"
            echo ""
            echo "ORCHESTRATOR, PLEASE:"
            echo "1. Run SETUP_EFFORT_INFRASTRUCTURE for this effort"
            echo "2. Verify the directory was actually created"
            echo "3. Check if you have the correct path"
            echo "4. Re-spawn me with the correct, existing directory"
            echo ""
            echo "GRADING VIOLATION: R208 - Wrong directory"
            exit 208
        fi
    else
        echo "✅ Already in correct directory"
    fi
else
    echo "⚠️ No WORKING_DIRECTORY in plan, using current: $(pwd)"
fi
```

### Step 3: EXTRACT AND VERIFY R209 METADATA
```bash
# Extract metadata
WORKING_DIR=$(grep "**WORKING_DIRECTORY**:" "$IMPLEMENTATION_PLAN" | cut -d: -f2- | xargs)
BRANCH=$(grep "**BRANCH**:" "$IMPLEMENTATION_PLAN" | cut -d: -f2- | xargs)
EFFORT_NAME=$(grep "**EFFORT_NAME**:" "$IMPLEMENTATION_PLAN" | cut -d: -f2- | xargs)
PHASE=$(grep "**PHASE**:" "$IMPLEMENTATION_PLAN" | cut -d: -f2- | xargs)
WAVE=$(grep "**WAVE**:" "$IMPLEMENTATION_PLAN" | cut -d: -f2- | xargs)

if [ -z "$WORKING_DIR" ]; then
    echo "❌ FATAL: No R209 metadata in $IMPLEMENTATION_PLAN!"
    echo "Orchestrator failed to inject directory metadata!"
    exit 1
fi

echo "📋 EFFORT METADATA:"
echo "   Effort: $EFFORT_NAME"
echo "   Phase:  $PHASE"
echo "   Wave:   $WAVE"
echo "   Branch: $BRANCH"
echo "   Required Dir: $WORKING_DIR"
```

### Step 4: ESTABLISH UNREMOVABLE DIRECTORY LOCK
```bash
# Verify we're in the correct directory
CURRENT_DIR=$(pwd)
if [ "$CURRENT_DIR" != "$WORKING_DIR" ]; then
    echo "❌❌❌ DIRECTORY MISMATCH!"
    echo "   Current:  $CURRENT_DIR"
    echo "   Required: $WORKING_DIR"
    echo ""
    echo "ATTEMPTING TO NAVIGATE TO CORRECT DIRECTORY..."
    
    if [ -d "$WORKING_DIR" ]; then
        cd "$WORKING_DIR"
        echo "✅ Navigated to: $(pwd)"
    else
        echo "❌ FATAL: Target directory doesn't exist!"
        exit 1
    fi
fi

# SET UNREMOVABLE ENVIRONMENT VARIABLES
export EFFORT_ISOLATION_DIR="$(pwd)"
export EFFORT_NAME="$EFFORT_NAME"
export EFFORT_PHASE="$PHASE"
export EFFORT_WAVE="$WAVE"

# MAKE THEM READONLY (CANNOT BE CHANGED)
readonly EFFORT_ISOLATION_DIR
readonly EFFORT_NAME
readonly EFFORT_PHASE
readonly EFFORT_WAVE

echo ""
echo "🔐🔐🔐 DIRECTORY LOCK ESTABLISHED 🔐🔐🔐"
echo "   EFFORT_ISOLATION_DIR = $EFFORT_ISOLATION_DIR (READONLY)"
echo "   You CANNOT leave this directory!"
echo ""
```

### Step 5: CREATE R209 ACKNOWLEDGMENT FILE
```bash
# Create audit trail
cat > .r209-acknowledged << EOF
R209 Directory Isolation Acknowledged
======================================
Timestamp: $(date '+%Y-%m-%d %H:%M:%S %Z')
SW Engineer: $(whoami)
Effort: $EFFORT_NAME
Phase: $PHASE
Wave: $WAVE
Directory: $(pwd)
Branch: $(git branch --show-current)
Environment Variables: LOCKED (readonly)

I acknowledge that:
1. ALL work happens in this directory
2. I cannot cd out of this boundary
3. All code goes in pkg/ subdirectory
4. This is my isolated workspace
EOF

echo "✅ Created .r209-acknowledged audit file"
```

### Step 6: FINAL ACKNOWLEDGMENT
```bash
echo ""
echo "═══════════════════════════════════════════════════════"
echo "✅✅✅ SW ENGINEER INITIALIZATION COMPLETE ✅✅✅"
echo "═══════════════════════════════════════════════════════"
echo ""
echo "📍 WORKING IN: $(pwd)"
echo "🔒 LOCKED TO: $EFFORT_ISOLATION_DIR"
echo "📝 EFFORT: $EFFORT_NAME (Phase $PHASE, Wave $WAVE)"
echo ""
echo "NEXT STEPS:"
echo "1. Read $IMPLEMENTATION_PLAN completely"
echo "2. Update work-log.md with startup entry"
echo "3. Begin implementation in pkg/ directory"
echo "4. NEVER attempt to leave this directory"
echo "═══════════════════════════════════════════════════════"
```

## State Transition

After completing ALL initialization steps:
1. Verify .r209-acknowledged exists
2. Confirm readonly environment variables are set
3. Transition to IMPLEMENTATION state

## Critical Reminders

- **NEVER** skip the directory verification
- **ALWAYS** create the acknowledgment file
- **NEVER** attempt to unset readonly variables
- **ALWAYS** echo your current directory before starting work
- If spawned in wrong directory, **EXIT IMMEDIATELY**