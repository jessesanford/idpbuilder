# 🔴🔴🔴 RULE R513 - ORCHESTRATOR LAYER TRACKING PROTOCOL (SUPREME LAW) 🔴🔴🔴

## Metadata
- **Rule ID**: R513
- **Title**: Orchestrator Layer Tracking Protocol
- **Criticality**: 🔴🔴🔴 SUPREME LAW
- **Category**: State Management / Cascade Operations
- **Version**: 1.0
- **Date**: 2025-10-05
- **Status**: PRODUCTION
- **Penalty**: -100% for layer confusion violations
- **Related Rules**: R512 (Trunk-Based Integration), R308 (Incremental Branching), R327 (Cascade)

## The Problem

**ORCHESTRATOR LAYER CONFUSION:**

During cascade operations, the orchestrator can lose track of:
1. **WHICH integration layer** it's currently working on (wave/phase/project)
2. **HOW DEEP** the fix cascade is (original bug vs fix-of-fix vs fix-of-fix-of-fix)
3. **WHAT BASE BRANCH** should be used for the current layer
4. **WHERE FIXES** should be applied (which effort branches)

This leads to catastrophic failures:
- Wrong base branch selected
- Fixes applied to wrong integration level
- Cascade loops (fixing same layer repeatedly)
- Lost track of cascade progress

## The Rule

**THE ORCHESTRATOR MUST MAINTAIN EXPLICIT LAYER TRACKING METADATA!**

All cascade and integration operations require:
1. **Layer identification** - Explicit metadata for current layer
2. **Depth tracking** - How many cascade iterations deep
3. **Base branch validation** - Verify base matches layer
4. **Source validation** - Verify sources match layer per R512

## 🔴🔴🔴 REQUIRED STATE MACHINE METADATA 🔴🔴🔴

### Layer Tracking Structure

Add to `orchestrator-state-v3.json`:

```json
{
  "layer_tracking": {
    "current_integration_layer": {
      "type": "wave|phase|project",
      "identifier": "phase1-wave2-integration",
      "phase_number": 1,
      "wave_number": 2,
      "base_branch": "phase1-wave1-integration",
      "last_updated": "2025-10-05T16:42:00Z"
    },
    "cascade_layer_depth": {
      "current_depth": 1,
      "max_depth_seen": 1,
      "layer_history": [
        {
          "depth": 1,
          "trigger": "integration_build_failure",
          "started_at": "2025-10-05T16:00:00Z",
          "integration_layer": "wave",
          "integration_id": "phase1-wave2-integration"
        }
      ]
    },
    "integration_source_tracking": {
      "expected_source_type": "effort",
      "expected_source_pattern": "phase1/wave2/*",
      "actual_sources": [
        "phase1/wave2/effort1",
        "phase1/wave2/effort2"
      ],
      "validation_status": "valid",
      "last_validated": "2025-10-05T16:42:00Z"
    }
  }
}
```

### Field Definitions

#### current_integration_layer
- **type**: "wave", "phase", or "project" - WHICH layer
- **identifier**: Full branch name (phase1-wave2-integration)
- **phase_number**: Integer phase number
- **wave_number**: Integer wave number (null for phase/project)
- **base_branch**: Expected R308 incremental base
- **last_updated**: ISO timestamp

#### cascade_layer_depth
- **current_depth**: Integer (1 = first cascade, 2 = fix-of-fix, etc.)
- **max_depth_seen**: Highest depth reached (prevent infinite loops)
- **layer_history**: Array of all cascade layers with triggers

#### integration_source_tracking
- **expected_source_type**: "effort" (per R512 - NEVER "integration")
- **expected_source_pattern**: Regex pattern for valid source branches
- **actual_sources**: List of branches being integrated
- **validation_status**: "valid", "invalid", "needs_check"
- **last_validated**: ISO timestamp

## 🔴🔴🔴 MANDATORY TRACKING OPERATIONS 🔴🔴🔴

### 1. Layer Initialization (All Integration States)

```bash
# BEFORE any integration operation
initialize_layer_tracking() {
    local integration_type="$1"  # wave, phase, or project
    local integration_id="$2"     # phase1-wave2-integration
    local base_branch="$3"        # phase1-wave1-integration

    echo "🔴 R513: Initializing layer tracking for $integration_type"

    # Extract phase/wave numbers
    PHASE_NUM=$(echo "$integration_id" | grep -oP 'phase\K[0-9]+')
    WAVE_NUM=$(echo "$integration_id" | grep -oP 'wave\K[0-9]+' || echo "null")

    # Set expected source pattern per R512
    case "$integration_type" in
        wave)
            EXPECTED_PATTERN="phase${PHASE_NUM}/wave${WAVE_NUM}/*"
            ;;
        phase)
            EXPECTED_PATTERN="phase${PHASE_NUM}/*"
            ;;
        project)
            EXPECTED_PATTERN="phase*/*"
            ;;
    esac

    # Initialize tracking
    jq --arg type "$integration_type" \
       --arg id "$integration_id" \
       --argjson phase "$PHASE_NUM" \
       --argjson wave "$WAVE_NUM" \
       --arg base "$base_branch" \
       --arg pattern "$EXPECTED_PATTERN" \
       --arg ts "$(date -Iseconds)" '
        .layer_tracking.current_integration_layer = {
            "type": $type,
            "identifier": $id,
            "phase_number": $phase,
            "wave_number": $wave,
            "base_branch": $base,
            "last_updated": $ts
        } |
        .layer_tracking.integration_source_tracking = {
            "expected_source_type": "effort",
            "expected_source_pattern": $pattern,
            "actual_sources": [],
            "validation_status": "needs_check",
            "last_validated": null
        }' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json

    echo "✅ R513: Layer tracking initialized for $integration_type integration"
}
```

### 2. Cascade Depth Tracking (CASCADE_REINTEGRATION State)

```bash
# When entering CASCADE_REINTEGRATION
increment_cascade_depth() {
    local trigger="$1"  # integration_build_failure, test_failure, etc.
    local integration_layer="$2"  # wave, phase, project
    local integration_id="$3"  # branch name

    echo "🔴 R513: Incrementing cascade depth"

    # Get current depth
    CURRENT_DEPTH=$(jq -r '.layer_tracking.cascade_layer_depth.current_depth // 0' orchestrator-state-v3.json)
    NEW_DEPTH=$((CURRENT_DEPTH + 1))

    # Safety check: Prevent infinite cascades
    if [[ "$NEW_DEPTH" -gt 10 ]]; then
        echo "❌ R513 VIOLATION: Cascade depth exceeds maximum (10 layers)!"
        echo "This indicates an infinite loop or severe architectural problem"
        return 1
    fi

    # Record new layer
    jq --argjson depth "$NEW_DEPTH" \
       --arg trigger "$trigger" \
       --arg layer "$integration_layer" \
       --arg id "$integration_id" \
       --arg ts "$(date -Iseconds)" '
        .layer_tracking.cascade_layer_depth.current_depth = $depth |
        .layer_tracking.cascade_layer_depth.max_depth_seen =
            (if .layer_tracking.cascade_layer_depth.max_depth_seen > $depth
             then .layer_tracking.cascade_layer_depth.max_depth_seen
             else $depth end) |
        .layer_tracking.cascade_layer_depth.layer_history += [{
            "depth": $depth,
            "trigger": $trigger,
            "started_at": $ts,
            "integration_layer": $layer,
            "integration_id": $id
        }]' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json

    echo "✅ R513: Now at cascade depth $NEW_DEPTH"
}

# When completing a cascade layer
decrement_cascade_depth() {
    local completed_integration_id="$1"

    echo "🔴 R513: Decrementing cascade depth"

    CURRENT_DEPTH=$(jq -r '.layer_tracking.cascade_layer_depth.current_depth' orchestrator-state-v3.json)
    NEW_DEPTH=$((CURRENT_DEPTH - 1))

    if [[ "$NEW_DEPTH" -lt 0 ]]; then
        echo "⚠️ R513 WARNING: Cascade depth would go negative, resetting to 0"
        NEW_DEPTH=0
    fi

    jq --argjson depth "$NEW_DEPTH" \
       --arg id "$completed_integration_id" \
       --arg ts "$(date -Iseconds)" '
        .layer_tracking.cascade_layer_depth.current_depth = $depth |
        (.layer_tracking.cascade_layer_depth.layer_history[] |
         select(.integration_id == $id)) |= . + {"completed_at": $ts}
    ' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json

    echo "✅ R513: Cascade depth reduced to $NEW_DEPTH"
}
```

### 3. Source Validation (All Integration States)

```bash
# BEFORE spawning integration agent
validate_integration_sources() {
    local SOURCE_BRANCHES=("$@")

    echo "🔴 R513: Validating integration sources per R512"

    # Get expected metadata
    EXPECTED_TYPE=$(jq -r '.layer_tracking.integration_source_tracking.expected_source_type' orchestrator-state-v3.json)
    EXPECTED_PATTERN=$(jq -r '.layer_tracking.integration_source_tracking.expected_source_pattern' orchestrator-state-v3.json)
    INTEGRATE_WAVE_EFFORTS_TYPE=$(jq -r '.layer_tracking.current_integration_layer.type' orchestrator-state-v3.json)

    # Validate per R512
    ALL_VALID=true
    for source in "${SOURCE_BRANCHES[@]}"; do
        # R512 RULE: Sources MUST be effort branches, NOT integrations
        if [[ "$source" =~ -integration$ ]]; then
            echo "❌ R512/R513 VIOLATION: Cannot integrate integration branch '$source'!"
            echo "At layer: $INTEGRATE_WAVE_EFFORTS_TYPE"
            echo "Expected: effort branches matching '$EXPECTED_PATTERN'"
            ALL_VALID=false
        fi

        # Validate pattern match
        if ! [[ "$source" =~ $EXPECTED_PATTERN ]]; then
            echo "⚠️ R513 WARNING: Source '$source' doesn't match expected pattern '$EXPECTED_PATTERN'"
        fi
    done

    # Update tracking
    if [[ "$ALL_VALID" == "true" ]]; then
        SOURCES_JSON=$(printf '%s\n' "${SOURCE_BRANCHES[@]}" | jq -R . | jq -s .)
        jq --argjson sources "$SOURCES_JSON" \
           --arg ts "$(date -Iseconds)" '
            .layer_tracking.integration_source_tracking.actual_sources = $sources |
            .layer_tracking.integration_source_tracking.validation_status = "valid" |
            .layer_tracking.integration_source_tracking.last_validated = $ts
        ' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json

        echo "✅ R513: All sources validated successfully"
        return 0
    else
        jq '.layer_tracking.integration_source_tracking.validation_status = "invalid"' \
           orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
        echo "❌ R513: Source validation FAILED!"
        return 1
    fi
}
```

### 4. Base Branch Validation (All Integration States)

```bash
# BEFORE creating infrastructure
validate_base_branch() {
    local actual_base="$1"

    echo "🔴 R513: Validating base branch selection"

    # Get expected base from layer tracking
    EXPECTED_BASE=$(jq -r '.layer_tracking.current_integration_layer.base_branch' orchestrator-state-v3.json)
    INTEGRATE_WAVE_EFFORTS_TYPE=$(jq -r '.layer_tracking.current_integration_layer.type' orchestrator-state-v3.json)

    if [[ "$actual_base" != "$EXPECTED_BASE" ]]; then
        echo "❌ R513 VIOLATION: Base branch mismatch!"
        echo "Expected (per R308/R513): $EXPECTED_BASE"
        echo "Actual: $actual_base"
        echo "Integration type: $INTEGRATE_WAVE_EFFORTS_TYPE"
        return 1
    fi

    echo "✅ R513: Base branch validated: $actual_base"
    return 0
}
```

## 🔴🔴🔴 INTEGRATE_WAVE_EFFORTS WITH EXISTING STATES 🔴🔴🔴

### CASCADE_REINTEGRATION State

Add to state entry:
```bash
# FIRST THING in CASCADE_REINTEGRATION
# Determine current layer from cascade context
INTEGRATE_WAVE_EFFORTS_TYPE=$(jq -r '.cascade_coordination.cascade_progress.current_operation_type' orchestrator-state-v3.json)

# Extract integration ID being recreated
TARGET_INTEGRATE_WAVE_EFFORTS=$(jq -r '.cascade_coordination.cascade_progress.current_operation_target' orchestrator-state-v3.json)

# Parse layer type from integration ID
if [[ "$TARGET_INTEGRATE_WAVE_EFFORTS" =~ wave[0-9]+-integration ]]; then
    LAYER_TYPE="wave"
elif [[ "$TARGET_INTEGRATE_WAVE_EFFORTS" =~ phase[0-9]+-integration ]]; then
    LAYER_TYPE="phase"
elif [[ "$TARGET_INTEGRATE_WAVE_EFFORTS" == "project-integration" ]]; then
    LAYER_TYPE="project"
fi

# Initialize or verify layer tracking
initialize_layer_tracking "$LAYER_TYPE" "$TARGET_INTEGRATE_WAVE_EFFORTS" "$BASE_BRANCH"

# Track cascade depth
increment_cascade_depth "integration_build_failure" "$LAYER_TYPE" "$TARGET_INTEGRATE_WAVE_EFFORTS"
```

### INTEGRATE_WAVE_EFFORTS, INTEGRATE_PHASE_WAVES, PROJECT_INTEGRATE_WAVE_EFFORTS States

Add to state entry:
```bash
# Initialize layer tracking
initialize_layer_tracking "$INTEGRATE_WAVE_EFFORTS_TYPE" "$TARGET_BRANCH" "$BASE_BRANCH"

# Before spawning integration agent, validate sources
validate_integration_sources "${SOURCE_BRANCHES[@]}" || {
    echo "❌ R512/R513: Source validation failed!"
    exit 512
}

# Validate base branch
validate_base_branch "$BASE_BRANCH" || {
    echo "❌ R513: Base branch validation failed!"
    exit 513
}
```

### SPAWN_CODE_REVIEWER_FIX_PLAN State

Add metadata for fix plan context:
```bash
# Record which layer fixes are for
CURRENT_LAYER=$(jq -r '.layer_tracking.current_integration_layer.type' orchestrator-state-v3.json)
CURRENT_LAYER_ID=$(jq -r '.layer_tracking.current_integration_layer.identifier' orchestrator-state-v3.json)

# Pass to code reviewer
jq --arg layer "$CURRENT_LAYER" \
   --arg id "$CURRENT_LAYER_ID" '
    .fix_plan_context = {
        "integration_layer": $layer,
        "integration_id": $id,
        "cascade_depth": .layer_tracking.cascade_layer_depth.current_depth
    }
' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
```

## Grading Impact

- **Implementing R513 layer tracking**: REQUIRED for cascade operations
- **Missing layer tracking metadata**: -100% (orchestrator loses position)
- **Wrong base branch due to layer confusion**: -100% (cascade corruption)
- **Integrating wrong source type**: -100% (R512 violation)
- **Infinite cascade loop (depth >10)**: -100% (system failure)
- **Correct layer tracking throughout cascade**: +100% (system works correctly)

## Integration with Other Rules

- **R512**: Validates source branches are efforts (not integrations)
- **R308**: Validates base branch follows incremental strategy
- **R327**: Ensures cascade recreations use correct layers
- **R410**: Tracks cascade depth for layered cascade protocol
- **R406**: Bug registry can reference cascade layer metadata

## Acknowledgment Required

Before performing ANY cascade or integration work, agents must acknowledge:

```
I acknowledge R513 - Orchestrator Layer Tracking Protocol:
- I will initialize layer_tracking metadata for all integration operations
- I will track cascade depth to prevent infinite loops
- I will validate source branches match expected layer per R512
- I will validate base branches match expected layer per R308
- I will maintain explicit metadata so orchestrator always knows its position
- Layer confusion = -100% automatic failure
```

## Summary

**REMEMBER: THE ORCHESTRATOR MUST ALWAYS KNOW WHERE IT IS!**

```
Without layer tracking:
  "Where am I?" - Unknown
  "What should I integrate?" - Guess
  "What's my base branch?" - Hope for best
  Result: CASCADE CORRUPTION

With R513 layer tracking:
  "Where am I?" - phase1-wave2-integration (wave layer, depth 1)
  "What should I integrate?" - phase1/wave2/* effort branches (R512)
  "What's my base branch?" - phase1-wave1-integration (R308)
  Result: CASCADE WORKS PERFECTLY
```

Every integration and cascade operation MUST maintain this metadata!

**THIS IS FUNDAMENTAL TO PREVENTING LAYER CONFUSION!**

## State Manager Coordination (SF 3.0)

State Manager tracks orchestrator layer (wave/phase/project) through iteration containers:
- **orchestrator-state-v3.json** `.state_machine.iteration_container` field
- **Current layer** determines which iteration container is active
- **integration-containers.json** tracks active integrations per layer
- **State transitions** update layer during COMPLETE_WAVE/COMPLETE_PHASE/COMPLETE_PROJECT

Layer tracking is atomic - all 4 files updated together when changing layers.

See: `integration-containers.json`, iteration container states (SETUP_WAVE_INFRASTRUCTURE, etc.)
