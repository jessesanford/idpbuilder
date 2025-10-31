# 🔴🔴🔴 RULE R208: ORCHESTRATOR WORKING DIRECTORY SPAWN PROTOCOL - SUPREME LAW 🔴🔴🔴

## ABSOLUTE SUPREMACY DECLARATION

**THIS RULE IS SUPREME LAW #2 OF THE SOFTWARE FACTORY**
- **NO OTHER RULE CAN OVERRIDE THIS (except R234)**
- **NOT EFFICIENCY CONCERNS**
- **NOT TIME CONSTRAINTS**
- **NOT "CONTINUOUS OPERATION"**
- **NOT ANYTHING**

## VIOLATION = IMMEDIATE CATASTROPHIC FAILURE

**Penalty for ANY violation:**
- **-100% GRADE (AUTOMATIC FAIL)**
- **IMMEDIATE TERMINATION**
- **NO RECOVERY POSSIBLE**
- **NO EXCUSES ACCEPTED**

## THE SUPREME LAW OF SPAWNING

### YOU MUST ALWAYS CD BEFORE SPAWN - NO EXCEPTIONS EVER

**Category:** SUPREME LAW  
**Agent:** Orchestrator ONLY  
**Criticality:** 🔴🔴🔴 ABSOLUTE - VIOLATION = INSTANT FAILURE 🔴🔴🔴  
**Priority:** SUPREME - This rule CANNOT be overridden by ANY other rule

## 🚨🚨🚨 MISSION CRITICAL RULE 🚨🚨🚨

**THE ONLY RELIABLE WAY TO SPAWN AGENTS IN THE CORRECT DIRECTORY IS FOR THE ORCHESTRATOR TO CHANGE TO THAT DIRECTORY FIRST!**

**SPAWNING WITHOUT CD'ING FIRST = -100% GRADE = AUTOMATIC FAILURE**

## The Problem This Solves

Agents inherit their working directory from their spawner. If the orchestrator spawns an agent while in the wrong directory, the agent starts in the wrong location and all subsequent work fails. This is the #1 cause of effort failures.

## 🔴🔴🔴 MANDATORY ENFORCEMENT - NO EXCEPTIONS 🔴🔴🔴

### EVERY SINGLE SPAWN MUST FOLLOW THIS EXACT SEQUENCE:

1. **DETERMINE** target directory for the agent
2. **CD** to that directory (MANDATORY - NO EXCEPTIONS)
3. **VERIFY** pwd output shows correct directory
4. **SPAWN** the agent (ONLY after successful CD)
5. **RETURN** to orchestrator directory

### FORBIDDEN ACTIONS (INSTANT FAILURE):
- ❌ Spawning without CD'ing first = **-100% GRADE**
- ❌ Assuming agent will CD itself = **-100% GRADE**
- ❌ Using --working-directory flag instead of CD = **-100% GRADE**
- ❌ Spawning from wrong directory "for efficiency" = **-100% GRADE**
- ❌ Skipping pwd verification = **-100% GRADE**

### THE ONLY ACCEPTABLE PATTERN:
```bash
# MANDATORY BEFORE EVERY SPAWN - NO EXCEPTIONS!
cd /path/to/effort/directory || exit 1
pwd  # MUST verify correct directory
/usr/bin/env bash -c 'task spawn agent-name "instructions"'
cd /original/directory  # Return after spawn
```

## The Protocol (MANDATORY)

### Step 1: THINK About Target Directory
```bash
# BEFORE ANY SPAWN, ASK:
# 1. What type of agent am I spawning?
# 2. Where does this agent need to work?
# 3. Does that directory exist?
# 4. Am I currently in that directory?

determine_agent_working_directory() {
    local AGENT_TYPE="$1"
    local EFFORT_NAME="$2"
    local SPLIT_NUM="${3:-}"
    
    case "$AGENT_TYPE" in
        "code-reviewer")
            if [ -n "$SPLIT_NUM" ]; then
                echo "efforts/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}/split-${SPLIT_NUM}"
            else
                echo "efforts/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}"
            fi
            ;;
        "sw-engineer")
            if [ -n "$SPLIT_NUM" ]; then
                echo "efforts/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}/split-${SPLIT_NUM}"
            else
                echo "efforts/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}"
            fi
            ;;
        "architect")
            echo "."  # Architect works from project root
            ;;
        *)
            echo "efforts/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}"
            ;;
    esac
}
```

### Step 2: Save Current Directory
```bash
# ALWAYS save where you are before changing
ORCHESTRATOR_DIR=$(pwd)
echo "📍 Orchestrator currently in: $ORCHESTRATOR_DIR"
```

### Step 3: Create Target Directory (If Needed)
```bash
TARGET_DIR=$(determine_agent_working_directory "$AGENT_TYPE" "$EFFORT_NAME" "$SPLIT_NUM")

if [ ! -d "$TARGET_DIR" ]; then
    echo "📁 Creating target directory: $TARGET_DIR"
    mkdir -p "$TARGET_DIR"
    
    # For effort directories, also create standard structure
    if [[ "$TARGET_DIR" == *"effort-"* ]]; then
        mkdir -p "$TARGET_DIR"/{src,tests,docs}
        touch "$TARGET_DIR"/work-log.md
    fi
fi
```

### Step 4: CHANGE TO TARGET DIRECTORY
```bash
# THIS IS THE CRITICAL STEP!
echo "🔄 Changing to target directory: $TARGET_DIR"
cd "$TARGET_DIR" || {
    echo "❌ FATAL: Cannot change to $TARGET_DIR"
    exit 1
}

# VERIFY we're in the right place
echo "✅ Now in: $(pwd)"
```

### Step 5: Spawn the Agent
```bash
# NOW spawn the agent - it will inherit our current directory
echo "🚀 Spawning $AGENT_TYPE in $(pwd)"

task spawn "$AGENT_TYPE" 
    --working-directory "$(pwd)" 
    --effort "$EFFORT_NAME" 
    --instructions "$INSTRUCTIONS"

# The agent is now running in the correct directory!
```

### Step 6: Return to Orchestrator Directory
```bash
# ALWAYS return to where you were
cd "$ORCHESTRATOR_DIR"
echo "📍 Returned to orchestrator directory: $(pwd)"
```

## Complete Spawn Function

```bash
spawn_agent_in_correct_directory() {
    local AGENT_TYPE="$1"
    local EFFORT_NAME="$2"
    local INSTRUCTIONS="$3"
    local SPLIT_NUM="${4:-}"
    
    echo "═══════════════════════════════════════════════"
    echo "🎯 R208: SPAWN DIRECTORY PROTOCOL"
    echo "═══════════════════════════════════════════════"
    
    # Step 1: Save current location
    local ORCHESTRATOR_DIR=$(pwd)
    echo "📍 Orchestrator starting from: $ORCHESTRATOR_DIR"
    
    # Step 2: Determine target directory
    local TARGET_DIR=$(determine_agent_working_directory "$AGENT_TYPE" "$EFFORT_NAME" "$SPLIT_NUM")
    echo "🎯 Agent needs to be in: $TARGET_DIR"
    
    # Step 3: Create if needed
    if [ ! -d "$TARGET_DIR" ]; then
        echo "📁 Creating: $TARGET_DIR"
        mkdir -p "$TARGET_DIR"
        
        # Create standard effort structure
        if [[ "$TARGET_DIR" == *"effort-"* ]]; then
            mkdir -p "$TARGET_DIR"/{src,tests,docs}
            touch "$TARGET_DIR"/work-log.md
            echo "📝 Created standard effort structure"
        fi
    fi
    
    # Step 4: CRITICAL - Change to target directory
    echo "🔄 CHANGING TO TARGET DIRECTORY..."
    cd "$TARGET_DIR" || {
        echo "❌ FATAL: Cannot change to $TARGET_DIR"
        return 1
    }
    
    # Step 5: Verify location
    echo "✅ Now in: $(pwd)"
    echo "📂 Contents:"
    ls -la
    
    # Step 6: Spawn agent (inherits our directory)
    echo "🚀 Spawning $AGENT_TYPE..."
    echo "   Agent will start in: $(pwd)"
    
    # ACTUAL SPAWN COMMAND HERE
    task spawn "$AGENT_TYPE" "$INSTRUCTIONS"
    
    # Step 7: Return to orchestrator directory
    echo "🔄 Returning to orchestrator directory..."
    cd "$ORCHESTRATOR_DIR"
    echo "📍 Back in: $(pwd)"
    
    echo "✅ R208 Protocol Complete"
    echo "═══════════════════════════════════════════════"
}
```

## Usage Examples

### Example 1: Spawning Code Reviewer for Planning
```bash
# Orchestrator is in project root
pwd  # /workspaces/project

# Spawn code reviewer for effort planning
spawn_agent_in_correct_directory 
    "code-reviewer" 
    "effort-api-types" 
    "Create implementation plan for API types"

# Code reviewer starts in: /workspaces/project/efforts/phase1/wave1/effort-api-types
# Orchestrator returns to: /workspaces/project
```

### Example 2: Spawning SW Engineer for Split
```bash
# Orchestrator in project root
pwd  # /workspaces/project

# Spawn SW engineer for split-002
spawn_agent_in_correct_directory 
    "sw-engineer" 
    "effort-api-types" 
    "Implement split-002 per plan" 
    "002"

# SW engineer starts in: /workspaces/project/efforts/phase1/wave1/effort-api-types/split-002
# Orchestrator returns to: /workspaces/project
```

### Example 3: Spawning Architect for Review
```bash
# Orchestrator anywhere
pwd  # /workspaces/project/efforts/phase1/wave2/effort-controllers

# Spawn architect (needs project root)
spawn_agent_in_correct_directory 
    "architect" 
    "wave2" 
    "Review wave 2 completion"

# Architect starts in: /workspaces/project (root)
# Orchestrator returns to: /workspaces/project/efforts/phase1/wave2/effort-controllers
```

## Validation Checks

### Pre-Spawn Checklist
```bash
pre_spawn_validation() {
    echo "🔍 R208 Pre-Spawn Validation"
    
    # 1. Know where you are
    echo "Current directory: $(pwd)"
    
    # 2. Know where agent needs to be
    echo "Target directory: $TARGET_DIR"
    
    # 3. Verify you can change there
    if [ -d "$TARGET_DIR" ] || mkdir -p "$TARGET_DIR" 2>/dev/null; then
        echo "✅ Target directory accessible"
    else
        echo "❌ Cannot access/create target directory!"
        return 1
    fi
    
    # 4. Test the round trip
    local ORIGINAL=$(pwd)
    cd "$TARGET_DIR" && cd "$ORIGINAL" || {
        echo "❌ Cannot perform directory round trip!"
        return 1
    }
    
    echo "✅ Pre-spawn validation passed"
}
```

## Common Failures (AND HOW TO PREVENT THEM)

### ❌ WRONG: Spawning without changing directory
```bash
# DON'T DO THIS!
task spawn sw-engineer "implement effort"
# Agent starts in whatever random directory orchestrator is in!
```

### ❌ WRONG: Assuming agent will cd itself
```bash
# DON'T DO THIS!
task spawn sw-engineer "cd to effort directory and implement"
# Agent might cd, but preflight checks run in wrong directory first!
```

### ❌ WRONG: Not creating directory first
```bash
# DON'T DO THIS!
cd efforts/phase1/wave1/effort-new/split-001  # Fails if doesn't exist!
```

### ✅ CORRECT: Always use the protocol
```bash
# ALWAYS DO THIS!
spawn_agent_in_correct_directory 
    "sw-engineer" 
    "effort-new" 
    "implement" 
    "001"
```

## Integration Points

This rule MUST be referenced and applied at:

1. **SPAWN_SW_ENGINEERS state** - Initial effort spawns
2. **CREATE_SPLIT_PLAN state** - Spawning code reviewer for splits
3. **SPLIT_IMPLEMENTATION state** - Spawning SW engineers for splits
4. **REVIEW_WAVE_ARCHITECTURE state** - Spawning architect
5. **FIX_ISSUES state** - Re-spawning for fixes
6. **ANY spawn operation** - No exceptions!

## 🔴🔴🔴 SUPREME LAW ENFORCEMENT 🔴🔴🔴

```bash
# THIS ENFORCEMENT IS ABSOLUTE - NO OVERRIDES POSSIBLE
supreme_law_R208_enforcement() {
    echo "🔴🔴🔴 R208 SUPREME LAW ENFORCEMENT ACTIVE 🔴🔴🔴"
    echo "SPAWNING WITHOUT CD = -100% GRADE = AUTOMATIC FAILURE"
    echo "Current directory: $(pwd)"
    
    # BLOCK ALL DIRECT SPAWNS - NO EXCEPTIONS
    task() {
        if [ "$1" = "spawn" ]; then
            echo "🔴🔴🔴 SUPREME LAW R208 VIOLATION DETECTED! 🔴🔴🔴"
            echo "❌ CATASTROPHIC FAILURE: Direct spawn without CD!"
            echo "❌ GRADE: -100% (AUTOMATIC FAIL)"
            echo "❌ You MUST use spawn_agent_in_correct_directory!"
            echo "🔴🔴🔴 TERMINATING IMMEDIATELY 🔴🔴🔴"
            exit 1  # IMMEDIATE TERMINATION
        fi
        command task "$@"
    }
    
    # VALIDATE EVERY SPAWN ATTEMPT
    validate_spawn_directory() {
        local current_dir=$(pwd)
        local expected_dir="$1"
        
        if [ "$current_dir" != "$expected_dir" ]; then
            echo "🔴🔴🔴 R208 SUPREME LAW VIOLATION! 🔴🔴🔴"
            echo "❌ Current: $current_dir"
            echo "❌ Expected: $expected_dir"
            echo "❌ GRADE: -100% (AUTOMATIC FAIL)"
            exit 1
        fi
    }
}

# ACTIVATE ON ORCHESTRATOR STARTUP - MANDATORY
supreme_law_R208_enforcement
```

## 🔴 ABSOLUTE REQUIREMENTS - NO NEGOTIATION 🔴

1. **CD FIRST** - Always change directory before spawn
2. **VERIFY PWD** - Always verify you're in the right place
3. **SPAWN SECOND** - Only spawn after successful CD
4. **NO SHORTCUTS** - No efficiency bypasses allowed
5. **NO EXCEPTIONS** - This applies to EVERY spawn, ALWAYS

## Summary

**R208 IS SUPREME LAW #2 - IT CANNOT BE VIOLATED FOR ANY REASON**

1. **THINK** about where the agent needs to be
2. **CD** to that directory (MANDATORY)
3. **VERIFY** with pwd (MANDATORY)
4. **SPAWN** the agent (inherits directory)
5. **RETURN** to your original directory

**SPAWNING WITHOUT CD'ING FIRST = -100% GRADE = AUTOMATIC FAILURE**

This is not optional. This is not a suggestion. This is SUPREME LAW.
## Software Factory 3.0 Integration

**State Tracking**: In SF 3.0, state transitions are tracked in `orchestrator-state-v3.json`:
```json
{
  "state_machine": {
    "current_state": "CURRENT_STATE_NAME",
    "previous_state": "PREVIOUS_STATE_NAME",
    "state_history": [...]
  }
}
```

**Compliance**: This rule applies to SF 3.0 state machine with appropriate state name mappings per R516 naming conventions.

**Reference**: See `docs/SOFTWARE-FACTORY-3.0-ARCHITECTURE.md` Part 2 for state machine design.

