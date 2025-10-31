# ANALYZE_IMPLEMENTATION_PARALLELIZATION State Rules

## 🔴🔴🔴 CRITICAL: THIS STATE READS WAVE PLANS, NOT EFFORT PLANS 🔴🔴🔴

**CONTEXT**: This state occurs AFTER effort plans are created but BEFORE implementation begins. The parallelization strategy COMES FROM the wave plan that was created by the architect.

## State Context

**Current Phase**: Read from `orchestrator-state-v3.json`
**Current Wave**: Read from `orchestrator-state-v3.json`
**Entry From**: WAITING_FOR_EFFORT_PLANS
**Exit To**: SPAWN_SW_ENGINEERS

## 🚨🚨🚨 PRIMARY OBJECTIVE 🚨🚨🚨

**EXTRACT and VALIDATE the parallelization strategy FROM THE WAVE PLAN to determine HOW to spawn SW Engineer agents.**

You are NOT deciding parallelization - you are IMPLEMENTING the architect's decision!

## Required Inputs

### 1. READ Wave Plan (PRIMARY SOURCE)
```bash
PHASE=$(jq -r '.current_phase' orchestrator-state-v3.json)
WAVE=$(jq -r '.current_wave' orchestrator-state-v3.json)

# Primary source - the wave implementation plan with timestamp pattern (updated path structure)
WAVE_PLAN=$(ls -t "$CLAUDE_PROJECT_DIR/phase-plans/phase${PHASE}/wave${WAVE}/WAVE-${PHASE}-${WAVE}-PLAN--"*.md 2>/dev/null | head -1)
if [ -n "$WAVE_PLAN" ] && [ -f "$WAVE_PLAN" ]; then
    echo "📖 Reading wave plan: $WAVE_PLAN"
    # READ: $WAVE_PLAN
else
    echo "❌ FATAL: Wave plan not found for Phase ${PHASE} Wave ${WAVE}"
    echo "Expected location: phase-plans/phase${PHASE}/wave${WAVE}/WAVE-${PHASE}-${WAVE}-PLAN--*.md"
    exit 1
fi
```

### 2. EXTRACT Parallelization Strategy from Wave Plan
The wave plan contains a "🚀 Parallelization Strategy" section with:
- **Execution Groups**: Sequential and parallel groupings
- **Parallelization Rules**: Specific dependencies and ordering
- **Orchestrator Spawning Strategy**: Explicit spawn commands

**Example Wave Plan Section:**
```yaml
group_1_sequential:  # MUST COMPLETE FIRST
  - effort_1: Contracts & Interfaces (blocking)

group_2_sequential:  # MUST COMPLETE SECOND
  - effort_2: Shared Libraries (blocking)

group_3_parallel:    # CAN RUN SIMULTANEOUSLY
  - effort_3: Feature A (independent)
  - effort_4: Feature B (independent)
  - effort_5: Feature C (independent)
```

### 3. READ Effort Plans (SECONDARY - for validation only)
```bash
# Helper function for R383/R343 compliant metadata paths
sf_metadata_path() {
    local file_type="$1"  # IMPLEMENTATION-PLAN, CODE-REVIEW-REPORT, etc.
    local phase="$2"
    local wave="$3"
    local effort="$4"
    local timestamp="${5:-$(date +%Y%m%d-%H%M%S)}"

    echo ".software-factory/phase${phase}/wave${wave}/${effort}/${file_type}--${timestamp}.md"
}

# Read effort plans to VALIDATE they exist and match wave plan
for effort_dir in efforts/phase${PHASE}/wave${WAVE}/*/; do
    if [ -d "$effort_dir" ]; then
        EFFORT_NAME=$(basename "$effort_dir")
        # Look for the most recent implementation plan
        IMPL_PLAN=$(ls -t "${effort_dir}.software-factory/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}/IMPLEMENTATION-PLAN--"*.md 2>/dev/null | head -1)
        if [ -f "$IMPL_PLAN" ]; then
            echo "✅ Found effort plan: $EFFORT_NAME at $IMPL_PLAN"
            # Quick validation - ensure effort is mentioned in wave plan
        fi
    fi
done
```

## 🔴🔴🔴 Parallelization Decision Matrix 🔴🔴🔴

### ✅ PARALLEL Execution Criteria
Efforts can run in PARALLEL when ALL of these are true:
- **No code dependencies**: Neither effort needs the other's code
- **No shared modules**: They work on independent components
- **No naming conflicts**: No risk of duplicate interfaces/types
- **Same base branch**: All can branch from the same cascade base
- **Wave plan specifies**: Listed in same `parallel` group

### ❌ SEQUENTIAL Execution Criteria
Efforts MUST run SEQUENTIALLY when ANY of these are true:
- **Code dependency**: Effort N+1 needs code from Effort N
- **Shared module work**: Both modify the same component
- **Risk of conflicts**: Overlapping functionality
- **Different base branches**: Per R308 cascade branching
- **Wave plan specifies**: Listed in `sequential` group

## Required Analysis Steps

### Step 1: Parse Wave Plan Parallelization Section
```bash
# Extract parallelization strategy from wave plan
extract_parallelization_strategy() {
    local WAVE_PLAN=$1

    # Find the Parallelization Strategy section
    grep -A 50 "## 🚀 Parallelization Strategy" "$WAVE_PLAN" | \
    grep -E "(group_[0-9]+_(sequential|parallel):|^\s+-\s+effort_)"
}
```

### Step 2: Build Execution Groups
```markdown
Based on wave plan, create execution groups:

**Group 1 (Sequential - Blocking)**:
- effort-1: Must complete first
- Reason: Creates interfaces others depend on

**Group 2 (Parallel - Can run simultaneously)**:
- effort-2: Independent feature A
- effort-3: Independent feature B
- effort-4: Independent feature C
- Reason: No interdependencies, work on separate modules

**Group 3 (Sequential - Depends on Group 2)**:
- effort-5: Integration effort
- Reason: Needs all parallel efforts complete
```

### Step 3: Determine Cascade Base Branches (R308)
```bash
# For PARALLEL efforts - all use same base
WAVE_BASE=$(determine_effort_base_branch $PHASE $WAVE)

# For SEQUENTIAL efforts - chain from each other
determine_sequential_bases() {
    local GROUP=$1
    local PREV_EFFORT=""

    for effort in $(get_efforts_in_group $GROUP); do
        if [ -z "$PREV_EFFORT" ]; then
            # First in sequence uses wave base
            BASE="$WAVE_BASE"
        else
            # Subsequent efforts use previous effort as base
            BASE="phase${PHASE}/wave${WAVE}/${PREV_EFFORT}"
        fi

        echo "Effort $effort will use base: $BASE"
        PREV_EFFORT=$effort
    done
}
```

## 📝 Output: Parallelization Plan

### Save to orchestrator-state-v3.json
```json
{
  "current_wave_parallelization": {
    "strategy": "MIXED",  // PARALLEL | SEQUENTIAL | MIXED
    "source": "phase-plans/PHASE-2-WAVE-1-PLAN.md",
    "execution_groups": [
      {
        "group_number": 1,
        "type": "sequential",
        "efforts": ["effort-1"],
        "blocking": true,
        "reason": "Creates contracts and interfaces others depend on"
      },
      {
        "group_number": 2,
        "type": "parallel",
        "efforts": ["effort-2", "effort-3", "effort-4"],
        "blocking": false,
        "reason": "Independent features with no interdependencies"
      },
      {
        "group_number": 3,
        "type": "sequential",
        "efforts": ["effort-5"],
        "blocking": false,
        "reason": "Integration effort requiring all features complete"
      }
    ],
    "cascade_bases": {
      "effort-1": "phase1-wave3-integration",
      "effort-2": "phase1-wave3-integration",
      "effort-3": "phase1-wave3-integration",
      "effort-4": "phase1-wave3-integration",
      "effort-5": "phase2/wave1/effort-5"
    },
    "spawn_plan": {
      "step_1": {
        "spawn": ["effort-1"],
        "wait_for": [],
        "type": "sequential"
      },
      "step_2": {
        "spawn": ["effort-2", "effort-3", "effort-4"],
        "wait_for": ["effort-1"],
        "type": "parallel"
      },
      "step_3": {
        "spawn": ["effort-5"],
        "wait_for": ["effort-2", "effort-3", "effort-4"],
        "type": "sequential"
      }
    },
    "analyzed_at": "2024-12-19T10:30:00Z",
    "analyzed_by": "orchestrator"
  }
}
```

## 🚨 Validation Requirements

### 1. Verify All Efforts Have Plans
```bash
# Every effort in wave plan must have an implementation plan
for effort in $(extract_efforts_from_wave_plan); do
    # Check for any implementation plan with timestamp
    IMPL_PLAN=$(ls -t "efforts/phase${PHASE}/wave${WAVE}/${effort}/.software-factory/phase${PHASE}/wave${WAVE}/${effort}/IMPLEMENTATION-PLAN--"*.md 2>/dev/null | head -1)
    if [ -z "$IMPL_PLAN" ] || [ ! -f "$IMPL_PLAN" ]; then
        echo "❌ FATAL: No implementation plan for $effort"
        echo "  Expected in: efforts/phase${PHASE}/wave${WAVE}/${effort}/.software-factory/phase${PHASE}/wave${WAVE}/${effort}/"
        exit 1
    fi
done
```

### 2. Verify Parallelization Logic
```bash
# Check parallel efforts don't have dependencies
for parallel_group in $(get_parallel_groups); do
    for effort in $(get_efforts_in_group $parallel_group); do
        # Find the most recent implementation plan
        IMPL_PLAN=$(ls -t "efforts/phase${PHASE}/wave${WAVE}/${effort}/.software-factory/phase${PHASE}/wave${WAVE}/${effort}/IMPLEMENTATION-PLAN--"*.md 2>/dev/null | head -1)
        # Check effort plan doesn't list dependencies on other parallel efforts
        if [ -f "$IMPL_PLAN" ] && grep -q "Dependencies:.*effort" "$IMPL_PLAN"; then
            echo "⚠️ WARNING: Parallel effort $effort may have dependencies"
        fi
    done
done
```

### 3. R356 Single-Effort Optimization
```bash
# If only one effort, skip complex analysis
EFFORT_COUNT=$(ls -d efforts/phase${PHASE}/wave${WAVE}/*/ 2>/dev/null | wc -l)
if [ "$EFFORT_COUNT" -eq 1 ]; then
    echo "📊 Single effort detected - applying R356 optimization"
    echo "Proceeding directly to SPAWN_SW_ENGINEERS"
    # Save minimal parallelization plan
    jq '.current_wave_parallelization.strategy = "SINGLE" |
        .current_wave_parallelization.r356_applied = true' \
        orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
fi
```

## 🔴 Common Mistakes to Avoid

### ❌ WRONG: Reading Individual Effort Plans for Strategy
```bash
# WRONG - Looking in wrong location without proper structure!
for effort in efforts/phase${PHASE}/wave${WAVE}/*/; do
    IMPL_PLAN="${effort}IMPLEMENTATION-PLAN.md"  # WRONG - Missing .software-factory and timestamp
    # This violates R383 and R343!
done
```

### ✅ CORRECT: Reading Wave Plan for Strategy
```bash
# CORRECT - Wave plan exists and contains parallelization strategy
WAVE_PLAN="phase-plans/PHASE-${PHASE}-WAVE-${WAVE}-PLAN.md"
# Extract parallelization from architect's wave plan

# For effort plans (when they exist), use proper R383/R343 structure:
for effort_dir in efforts/phase${PHASE}/wave${WAVE}/*/; do
    EFFORT_NAME=$(basename "$effort_dir")
    # Find most recent implementation plan with proper structure
    IMPL_PLAN=$(ls -t "${effort_dir}.software-factory/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}/IMPLEMENTATION-PLAN--"*.md 2>/dev/null | head -1)
    if [ -f "$IMPL_PLAN" ]; then
        echo "Found compliant plan: $IMPL_PLAN"
    fi
done
```

### ❌ WRONG: Deciding Parallelization Yourself
```bash
# WRONG - Orchestrator doesn't decide, it implements!
if [[ "$effort" == "contracts" ]]; then
    echo "I think this should be sequential"
fi
```

### ✅ CORRECT: Implementing Architect's Decision
```bash
# CORRECT - Extract what architect already decided
STRATEGY=$(grep -A 20 "Parallelization Strategy" "$WAVE_PLAN")
# Implement exactly what the wave plan specifies
```

## 📊 Example Scenarios

### Scenario 1: All Parallel Wave
```yaml
Wave Plan Says:
  group_1_parallel:
    - effort-1: UI Theme
    - effort-2: Documentation
    - effort-3: Performance

Analysis Result:
  strategy: PARALLEL
  All efforts branch from: phase1-wave2-integration
  Spawn: All three engineers simultaneously
```

### Scenario 2: Sequential Dependencies
```yaml
Wave Plan Says:
  group_1_sequential:
    - effort-1: API Foundation
  group_2_sequential:
    - effort-2: API Endpoints (needs foundation)
  group_3_sequential:
    - effort-3: API Client (needs endpoints)

Analysis Result:
  strategy: SEQUENTIAL
  effort-1 branches from: phase1-wave2-integration
  effort-2 branches from: phase2/wave1/effort-1
  effort-3 branches from: phase2/wave1/effort-2
  Spawn: One at a time, waiting for completion
```

### Scenario 3: Mixed Strategy
```yaml
Wave Plan Says:
  group_1_sequential:
    - effort-1: Contracts (blocking)
  group_2_parallel:
    - effort-2: Feature A
    - effort-3: Feature B
  group_3_sequential:
    - effort-4: Integration

Analysis Result:
  strategy: MIXED
  effort-1: Sequential first, branches from integration
  efforts-2,3: Parallel after 1, both branch from integration
  effort-4: Sequential after 2&3, branches from effort-3
```

## Exit Criteria

Before transitioning to SPAWN_SW_ENGINEERS:
- ✅ Wave plan parallelization extracted
- ✅ Execution groups defined
- ✅ Cascade bases determined per R308
- ✅ Spawn plan saved to state
- ✅ All effort plans validated to exist
- ✅ State file updated with parallelization plan

## State Transition

```bash
# Update state to SPAWN_SW_ENGINEERS
jq '.state_machine.current_state = "SPAWN_SW_ENGINEERS" |
    .state_transition_log += [{
        "from": "ANALYZE_IMPLEMENTATION_PARALLELIZATION",
        "to": "SPAWN_SW_ENGINEERS",
        "timestamp": "'$(date -Iseconds)'",
        "reason": "Parallelization analysis complete"
    }]' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json

echo "✅ Parallelization analysis complete"
echo "📋 Strategy: $(jq -r '.current_wave_parallelization.strategy' orchestrator-state-v3.json)"
echo "🚀 Ready to spawn agents according to plan"
```

## Integration with Other Rules

- **R308**: Cascade branching determines base branches
- **R356**: Single-effort optimization
- **R053**: Parallelization decision criteria
- **R219**: Dependency-aware effort planning
- **R501**: Infrastructure planning requirements

## Automation Flag

```bash
# After successful analysis and strategy creation:
echo "✅ Parallelization strategy determined"
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # Continue to SPAWN_SW_ENGINEERS
```

---

**REMEMBER**: You are IMPLEMENTING the architect's parallelization strategy from the wave plan, NOT deciding it yourself! The wave plan is the source of truth.
## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

**MUST output as LAST line before state completion:**
- `CONTINUE-SOFTWARE-FACTORY=TRUE` - State completed, automation can continue
- `CONTINUE-SOFTWARE-FACTORY=FALSE` - Error/block, manual intervention required

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-CONTINUATION-FLAG-MASTER-GUIDE.md**
