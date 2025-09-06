# 🔴🔴🔴 R208 SUPREME LAW ENFORCEMENT - CD BEFORE SPAWN 🔴🔴🔴

## THIS FILE ENFORCES SUPREME LAW #2 - NO EXCEPTIONS EVER

### VIOLATION = -100% GRADE = AUTOMATIC FAILURE

## 🚨🚨🚨 MANDATORY ENFORCEMENT FUNCTIONS 🚨🚨🚨

Add these functions to your orchestrator startup and use them ALWAYS:

```bash
# 🔴🔴🔴 R208 SUPREME LAW ENFORCEMENT FUNCTIONS 🔴🔴🔴

# MANDATORY: Load this at orchestrator startup
enforce_r208_supreme_law() {
    echo "🔴🔴🔴 R208 SUPREME LAW ENFORCEMENT ACTIVATED 🔴🔴🔴"
    echo "SPAWNING WITHOUT CD = -100% GRADE = AUTOMATIC FAILURE"
    
    # Create enforcement wrapper
    export R208_ENFORCED=true
    export ORCHESTRATOR_HOME=$(pwd)
}

# MANDATORY: Use for EVERY spawn operation
r208_spawn_agent() {
    local AGENT_TYPE="$1"
    local TARGET_DIR="$2"
    local INSTRUCTIONS="$3"
    
    echo "═══════════════════════════════════════════════"
    echo "🔴🔴🔴 R208 SUPREME LAW: CD BEFORE SPAWN 🔴🔴🔴"
    echo "═══════════════════════════════════════════════"
    
    # Step 1: Verify enforcement is active
    if [[ "$R208_ENFORCED" != "true" ]]; then
        echo "❌ R208 ENFORCEMENT NOT ACTIVE!"
        echo "❌ GRADE: -100% (AUTOMATIC FAILURE)"
        exit 1
    fi
    
    # Step 2: Save current directory
    local ORIGINAL_DIR=$(pwd)
    echo "📍 Orchestrator currently in: $ORIGINAL_DIR"
    
    # Step 3: CD to target directory (MANDATORY)
    echo "🔴 R208: CD'ing to: $TARGET_DIR"
    cd "$TARGET_DIR" || {
        echo "❌ R208 VIOLATION: Failed to CD to $TARGET_DIR"
        echo "❌ GRADE: -100% (AUTOMATIC FAILURE)"
        exit 1
    }
    
    # Step 4: Verify we're in correct directory (MANDATORY)
    local ACTUAL_DIR=$(pwd)
    echo "📍 R208 PWD VERIFICATION: $ACTUAL_DIR"
    
    if [[ "$ACTUAL_DIR" != *"$TARGET_DIR"* ]]; then
        echo "❌ R208 VIOLATION: Not in expected directory!"
        echo "❌ Expected: $TARGET_DIR"
        echo "❌ Actual: $ACTUAL_DIR"
        echo "❌ GRADE: -100% (AUTOMATIC FAILURE)"
        exit 1
    fi
    
    echo "✅ R208: Confirmed in correct directory"
    
    # Step 5: Spawn the agent (inherits our directory)
    echo "🚀 Spawning $AGENT_TYPE (will start in: $ACTUAL_DIR)"
    /usr/bin/env bash -c "task spawn $AGENT_TYPE '$INSTRUCTIONS'"
    
    # Step 6: Return to original directory (MANDATORY)
    cd "$ORIGINAL_DIR"
    echo "📍 R208: Returned to orchestrator directory: $(pwd)"
    
    echo "✅ R208 SUPREME LAW COMPLIANCE: VERIFIED"
    echo "═══════════════════════════════════════════════"
}

# MANDATORY: Validate spawn attempt
r208_validate_spawn() {
    local TARGET_DIR="$1"
    
    echo "🔍 R208 Pre-Spawn Validation"
    
    # Check if we can access target
    if [[ ! -d "$TARGET_DIR" ]]; then
        echo "⚠️ Target directory doesn't exist: $TARGET_DIR"
        echo "📁 Creating directory..."
        mkdir -p "$TARGET_DIR" || {
            echo "❌ R208: Cannot create $TARGET_DIR"
            echo "❌ GRADE: -100% (AUTOMATIC FAILURE)"
            return 1
        }
    fi
    
    # Test CD capability
    local ORIGINAL=$(pwd)
    cd "$TARGET_DIR" 2>/dev/null || {
        echo "❌ R208: Cannot CD to $TARGET_DIR"
        echo "❌ GRADE: -100% (AUTOMATIC FAILURE)"
        return 1
    }
    cd "$ORIGINAL"
    
    echo "✅ R208 Pre-validation passed"
    return 0
}
```

## 🔴 USAGE EXAMPLES - THE ONLY ACCEPTABLE WAY

### Example 1: Spawning SW Engineer
```bash
# MANDATORY: Always use r208_spawn_agent
r208_spawn_agent \
    "sw-engineer" \
    "efforts/phase1/wave1/effort-api-types" \
    "Implement API types per IMPLEMENTATION-PLAN.md"
```

### Example 2: Spawning Code Reviewer for Split
```bash
# MANDATORY: Always use r208_spawn_agent
r208_spawn_agent \
    "code-reviewer" \
    "efforts/phase1/wave1/effort-api-types/split-001" \
    "Review split-001 implementation"
```

### Example 3: Spawning Integration Agent
```bash
# MANDATORY: Always use r208_spawn_agent
r208_spawn_agent \
    "integration" \
    "efforts/phase1/wave1/integration-workspace" \
    "Execute WAVE-MERGE-PLAN.md"
```

### Example 4: Parallel Spawning (ALL IN ONE MESSAGE)
```bash
# For parallel agents, call r208_spawn_agent multiple times in ONE message:

# Agent 1
r208_spawn_agent \
    "sw-engineer" \
    "efforts/phase1/wave1/effort-api" \
    "Implement API effort"

# Agent 2 (same message!)
r208_spawn_agent \
    "sw-engineer" \
    "efforts/phase1/wave1/effort-auth" \
    "Implement auth effort"

# Agent 3 (same message!)
r208_spawn_agent \
    "sw-engineer" \
    "efforts/phase1/wave1/effort-db" \
    "Implement database effort"
```

## 🚨 FORBIDDEN PATTERNS - INSTANT FAILURE

### ❌ NEVER DO THIS: Direct spawn without CD
```bash
# THIS IS -100% FAILURE!
task spawn sw-engineer "implement something"
```

### ❌ NEVER DO THIS: Assume agent will CD
```bash
# THIS IS -100% FAILURE!
task spawn sw-engineer "cd to effort directory and implement"
```

### ❌ NEVER DO THIS: Use --working-directory flag
```bash
# THIS IS -100% FAILURE!
task spawn sw-engineer --working-directory /some/path "implement"
```

### ❌ NEVER DO THIS: Skip pwd verification
```bash
# THIS IS -100% FAILURE!
cd /some/directory
task spawn sw-engineer "implement"  # No pwd verification!
```

## 🔴 INTEGRATION WITH STATE MACHINE

### States That MUST Use R208 Enforcement:
1. **SPAWN_AGENTS** - SW Engineers
2. **SPAWN_CODE_REVIEWERS_EFFORT_PLANNING** - Code Reviewers for efforts
3. **SPAWN_CODE_REVIEWER_MERGE_PLAN** - Code Reviewer for integration
4. **SPAWN_CODE_REVIEWER_PHASE_IMPL** - Code Reviewer for phase plan
5. **SPAWN_CODE_REVIEWER_WAVE_IMPL** - Code Reviewer for wave plan
6. **SPAWN_INTEGRATION_AGENT** - Integration Agent
7. **SPAWN_ARCHITECT_**** - Any Architect spawns
8. **ANY OTHER SPAWN** - No exceptions!

## 🔴 MANDATORY STARTUP SEQUENCE

```bash
# EVERY orchestrator startup MUST include:
echo "🏭 ORCHESTRATOR STARTUP"
echo "📋 Loading R208 SUPREME LAW enforcement..."

# Load enforcement
enforce_r208_supreme_law

# Verify enforcement is active
if [[ "$R208_ENFORCED" != "true" ]]; then
    echo "❌ FATAL: R208 enforcement not active!"
    echo "❌ GRADE: -100% (AUTOMATIC FAILURE)"
    exit 1
fi

echo "✅ R208 SUPREME LAW enforcement active"
echo "🔴 Remember: SPAWN WITHOUT CD = -100% FAIL"
```

## 🔴 GRADING IMPACT

**R208 Violations Result In:**
- -100% Grade (AUTOMATIC FAIL)
- Immediate Termination
- No Recovery Possible
- No Partial Credit
- No Excuses Accepted

**Common Failure Scenarios:**
1. Spawning from wrong directory = -100%
2. Skipping CD step = -100%
3. Not verifying pwd = -100%
4. Using shortcuts = -100%
5. ANY bypass attempt = -100%

## 🔴 REMEMBER: THIS IS SUPREME LAW #2

- **R234** (State Traversal) is SUPREME LAW #1
- **R208** (CD Before Spawn) is SUPREME LAW #2
- **NO OTHER RULE** can override these (except R234 can override R208)
- **NO EFFICIENCY CONCERNS** can bypass this
- **NO TIME PRESSURE** can skip this
- **NO "CONTINUOUS OPERATION"** can avoid this

**SPAWN WITHOUT CD = AUTOMATIC FAILURE**

---

**THIS LAW IS ABSOLUTE. THIS LAW IS FINAL. THIS LAW CANNOT BE NEGOTIATED.**