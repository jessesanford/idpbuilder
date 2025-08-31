#!/bin/bash

# validate-spawn-directory.sh
# Validates R208 compliance - orchestrator spawn directory protocol

set -e

echo "=========================================="
echo "R208: Spawn Directory Validation"
echo "=========================================="
echo "MISSION CRITICAL: Agents MUST spawn in correct directories!"
echo ""

# Function to validate spawn setup
validate_spawn_setup() {
    local AGENT_TYPE="$1"
    local EFFORT_NAME="$2"
    local SPLIT_NUM="${3:-}"
    local PHASE="${4:-1}"
    local WAVE="${5:-1}"
    
    echo "Validating spawn setup for: $AGENT_TYPE"
    echo "  Effort: $EFFORT_NAME"
    echo "  Split: ${SPLIT_NUM:-none}"
    echo ""
    
    # Determine expected directory
    case "$AGENT_TYPE" in
        "code-reviewer"|"sw-engineer")
            if [ -n "$SPLIT_NUM" ]; then
                TARGET_DIR="efforts/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}/split-${SPLIT_NUM}"
            else
                TARGET_DIR="efforts/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}"
            fi
            ;;
        "architect")
            TARGET_DIR="."  # Architect works from root
            ;;
        "orchestrator")
            TARGET_DIR="."  # Orchestrator at root
            ;;
        *)
            echo "  ⚠️ Unknown agent type: $AGENT_TYPE"
            TARGET_DIR="efforts/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}"
            ;;
    esac
    
    echo "  Expected directory: $TARGET_DIR"
    
    # Check if directory exists
    if [ -d "$TARGET_DIR" ]; then
        echo "  ✅ Directory exists"
    else
        echo "  ❌ Directory does NOT exist!"
        echo "  📁 Creating: mkdir -p $TARGET_DIR"
        mkdir -p "$TARGET_DIR"
        echo "  ✅ Directory created"
    fi
    
    # Test we can cd to it
    local CURRENT_DIR=$(pwd)
    if cd "$TARGET_DIR" 2>/dev/null; then
        echo "  ✅ Can change to directory"
        echo "  📍 Agent would start in: $(pwd)"
        cd "$CURRENT_DIR"
        echo "  ✅ Returned to: $(pwd)"
    else
        echo "  ❌ FATAL: Cannot change to $TARGET_DIR!"
        return 1
    fi
    
    # Check for required files/structure
    if [[ "$TARGET_DIR" == *"effort-"* ]] && [ "$TARGET_DIR" != "." ]; then
        echo "  Checking effort structure..."
        if [ -f "$TARGET_DIR/work-log.md" ]; then
            echo "    ✅ work-log.md exists"
        else
            echo "    ⚠️ work-log.md missing - creating"
            touch "$TARGET_DIR/work-log.md"
        fi
        
        if [ -f "$TARGET_DIR/IMPLEMENTATION-PLAN.md" ]; then
            echo "    ✅ IMPLEMENTATION-PLAN.md exists"
        else
            echo "    ⚠️ IMPLEMENTATION-PLAN.md not yet created"
        fi
    fi
    
    echo "  ----------------------------------------"
    echo "  Summary: Spawn setup validated ✅"
    echo ""
}

# Function to simulate spawn with R208
simulate_spawn_with_R208() {
    local AGENT_TYPE="$1"
    local EFFORT_NAME="$2"
    local SPLIT_NUM="${3:-}"
    
    echo "🚀 Simulating R208 spawn protocol..."
    echo "=========================================="
    
    # Step 1: Save current directory
    local ORCHESTRATOR_DIR=$(pwd)
    echo "📍 Step 1: Orchestrator in: $ORCHESTRATOR_DIR"
    
    # Step 2: Determine target
    local TARGET_DIR
    if [ -n "$SPLIT_NUM" ]; then
        TARGET_DIR="efforts/phase1/wave1/${EFFORT_NAME}/split-${SPLIT_NUM}"
    else
        TARGET_DIR="efforts/phase1/wave1/${EFFORT_NAME}"
    fi
    [ "$AGENT_TYPE" = "architect" ] && TARGET_DIR="."
    
    echo "🎯 Step 2: Target directory: $TARGET_DIR"
    
    # Step 3: Create if needed
    if [ ! -d "$TARGET_DIR" ]; then
        echo "📁 Step 3: Creating directory..."
        mkdir -p "$TARGET_DIR"
    else
        echo "✅ Step 3: Directory exists"
    fi
    
    # Step 4: Change to target
    echo "🔄 Step 4: Changing to target directory..."
    cd "$TARGET_DIR" || {
        echo "❌ FATAL: Cannot cd to $TARGET_DIR"
        return 1
    }
    echo "✅ Now in: $(pwd)"
    echo "   Agent will spawn here!"
    
    # Step 5: Simulate spawn
    echo "🚀 Step 5: [SIMULATED] Spawning $AGENT_TYPE"
    echo "   Working directory: $(pwd)"
    
    # Step 6: Return
    echo "🔄 Step 6: Returning to orchestrator directory..."
    cd "$ORCHESTRATOR_DIR"
    echo "📍 Back in: $(pwd)"
    
    echo "=========================================="
    echo "✅ R208 Protocol simulation complete!"
    echo ""
}

# Main validation
echo "1. Testing spawn setup validation..."
echo "======================================"
validate_spawn_setup "sw-engineer" "effort-api-types" "" 1 1
validate_spawn_setup "sw-engineer" "effort-api-types" "001" 1 1
validate_spawn_setup "code-reviewer" "effort-controllers" "" 1 2
validate_spawn_setup "architect" "phase1-review" "" 1 1

echo ""
echo "2. Testing R208 spawn protocol..."
echo "======================================"

# Create temp test directory
TEST_DIR="/tmp/r208-test-$$"
mkdir -p "$TEST_DIR"
cd "$TEST_DIR"

simulate_spawn_with_R208 "sw-engineer" "effort-api-types"
simulate_spawn_with_R208 "sw-engineer" "effort-api-types" "002"
simulate_spawn_with_R208 "architect" "wave-review"

# Cleanup
cd /
rm -rf "$TEST_DIR"

echo ""
echo "=========================================="
echo "R208 Validation Complete"
echo "=========================================="
echo ""
echo "Key Takeaways:"
echo "1. ALWAYS change to target directory before spawning"
echo "2. Create directories if they don't exist"
echo "3. Return to orchestrator directory after spawn"
echo "4. Agents inherit the spawn directory"
echo "5. This is MISSION CRITICAL for success"