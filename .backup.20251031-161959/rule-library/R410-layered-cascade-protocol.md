# R410 - Layered Cascade Protocol

## Rule Metadata
- **ID**: R410.0.0
- **Criticality**: 🔴🔴🔴 SUPREME LAW
- **Category**: Fix Cascade Management
- **Enforcement**: Automatic layer creation, recursive solution
- **Grading Impact**: -100% for treating cascade build failures as blockers

## Summary
Build failures during CASCADE re-integration are EXPECTED and trigger automatic creation of NEW cascade layers. This is a RECURSIVE solution where each layer fixes bugs discovered by the previous layer's integration attempts.

## Problem Statement

When the orchestrator performs CASCADE_REINTEGRATION (per R327), it:
1. Deletes stale integration branches
2. Recreates them with fresh merges
3. Runs builds to validate integration

**CRITICAL ISSUE**: Build failures at step 3 reveal **NEW cross-effort bugs** that weren't visible in individual efforts. The current system treats these as blockers requiring manual intervention.

**WRONG APPROACH**:
```bash
# Build failed during phase1-wave2-integration recreation
echo "❌ Build failed - need manual approval to fix"
echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
# User must run /continue-orchestrating to proceed
```

**CORRECT APPROACH (R410)**:
```bash
# Build failed during phase1-wave2-integration recreation
echo "📋 Build failure discovered NEW bugs (cross-effort issues)"
echo "🔄 Starting CASCADE LAYER 2 to fix upstream branches"
# Document bugs, pause layer 1, start layer 2
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # AUTOMATION CONTINUES!
```

## The Layered Cascade Concept

### Layer Definition
A **cascade layer** is a complete fix-and-cascade cycle:
1. Bugs discovered (from integration build failures)
2. Fixes applied to effort branches
3. Cascade re-integration performed (R327)
4. Integration validated

### Layer Relationships
```
CASCADE LAYER 1 (original bugs from upstream merge)
├─ Fix effort branches with bugs
├─ Cascade to wave integration → BUILD FAILS! (discovers layer 2 bugs)
└─ PAUSE layer 1, start layer 2

CASCADE LAYER 2 (bugs discovered during layer 1 re-integration)
├─ Fix effort branches with NEW bugs
├─ Cascade to wave integration → BUILD FAILS! (discovers layer 3 bugs)
└─ PAUSE layer 2, start layer 3

CASCADE LAYER 3 (bugs discovered during layer 2 re-integration)
├─ Fix effort branches with NEW bugs
├─ Cascade to wave integration → BUILD PROJECT_DONE! ✅
└─ Complete layer 3, RESUME layer 2

RESUME LAYER 2 (layer 3 fixes now integrated)
├─ Re-attempt wave integration → BUILD PROJECT_DONE! ✅
└─ Complete layer 2, RESUME layer 1

RESUME LAYER 1 (layer 2+3 fixes now integrated)
├─ Re-attempt wave integration → BUILD PROJECT_DONE! ✅
└─ Complete layer 1, CASCADE COMPLETE!
```

### Why Layers Are Necessary

**Root Cause**: Cross-effort integration bugs are **NOT visible** until efforts are merged together.

**Example**:
- Effort 1 implements `getUserData()` returning `Promise<User>`
- Effort 2 implements `displayUser(user: User)` expecting sync data
- Both build INDIVIDUALLY (no errors)
- Integration FAILS because Effort 2 doesn't await the promise
- **This bug is ONLY visible during integration!**

## Mandatory Layer Tracking

### Schema Addition: `cascade_layers[]`

**Add to orchestrator-state.schema.json:**

```json
{
  "cascade_layers": {
    "type": "array",
    "description": "Tracks multiple cascade layers (recursive fix-cascade cycles)",
    "items": {
      "type": "object",
      "required": ["layer_id", "started_at", "trigger", "status"],
      "properties": {
        "layer_id": {
          "type": "integer",
          "description": "Sequential layer number (1, 2, 3...)"
        },
        "started_at": {
          "type": "string",
          "format": "date-time",
          "description": "When this layer started"
        },
        "trigger": {
          "type": "string",
          "description": "What triggered this layer",
          "examples": [
            "original_upstream_bugs",
            "build_failures_from_layer_1_reintegration",
            "build_failures_from_layer_2_reintegration"
          ]
        },
        "bugs": {
          "type": "array",
          "description": "Bugs discovered at this layer",
          "items": {
            "type": "string",
            "description": "Bug ID from bug_registry (e.g., BUG-001)"
          }
        },
        "status": {
          "type": "string",
          "enum": [
            "fixing_upstream",
            "reintegrating",
            "paused_for_layer_N",
            "completed",
            "failed"
          ],
          "description": "Current status of this layer"
        },
        "progress": {
          "type": "string",
          "description": "Human-readable progress for this layer",
          "examples": [
            "phase1_wave2_integration: build failed, needs fixes",
            "effort fixes complete, attempting re-integration",
            "all integrations passed, layer complete"
          ]
        },
        "completed_at": {
          "type": "string",
          "format": "date-time",
          "description": "When this layer completed (if status=completed)"
        }
      }
    }
  }
}
```

### Example Layer Tracking

```json
{
  "cascade_layers": [
    {
      "layer_id": 1,
      "started_at": "2025-10-05T04:00:00Z",
      "trigger": "original_upstream_bugs",
      "bugs": ["BUG-001", "BUG-002", "BUG-003"],
      "status": "paused_for_layer_2",
      "progress": "phase1_wave2_integration: build failed with 3 errors, started layer 2"
    },
    {
      "layer_id": 2,
      "started_at": "2025-10-05T05:10:00Z",
      "trigger": "build_failures_from_layer_1_reintegration",
      "bugs": ["BUG-004", "BUG-005"],
      "status": "fixing_upstream",
      "progress": "applying fixes to effort-3 and effort-5"
    }
  ]
}
```

## Decisive Decision Tree for CASCADE_REINTEGRATION

**When integration build FAILS during cascade:**

```
Integration build failed during CASCADE_REINTEGRATION?
├─ YES → Are these NEW bugs (not in current layer)?
│        ├─ YES → START NEW CASCADE LAYER
│        │        ├─ Document new bugs in bug_registry
│        │        ├─ Create layer N+1 metadata
│        │        ├─ Pause current layer (status = "paused_for_layer_N+1")
│        │        ├─ Transition to ERROR_RECOVERY (to start fixes)
│        │        └─ CONTINUE-SOFTWARE-FACTORY=TRUE! (automation!)
│        │
│        └─ NO → SAME bugs we're already fixing?
│                 ├─ ERROR: Infinite loop detected!
│                 ├─ Fixes didn't actually fix bugs
│                 └─ CONTINUE-SOFTWARE-FACTORY=FALSE (genuine block)
│
└─ NO → Build PROJECT_DONE?
         └─ YES → Complete current layer, RESUME previous layer
                  └─ CONTINUE-SOFTWARE-FACTORY=TRUE! (automation!)
```

## Implementation Protocol

### 1. Detect Build Failure During CASCADE

**In CASCADE_REINTEGRATION state:**

```bash
# After attempting integration build
if ! build_integration; then
    echo "🔍 Build failed during cascade re-integration"

    # Extract build errors
    BUILD_ERRORS=$(extract_build_errors build.log)

    # Check if these are NEW bugs
    if are_bugs_new "$BUILD_ERRORS"; then
        echo "📋 NEW bugs discovered (cross-effort issues)"
        echo "🔄 Triggering R410: Layered Cascade Protocol"

        # Start new layer (automated decision!)
        start_new_cascade_layer "$BUILD_ERRORS"
    else
        echo "❌ SAME bugs as current layer - fixes ineffective"
        echo "🚨 Infinite loop detected - MANUAL INTERVENTION REQUIRED"
        echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
        exit 1
    fi
fi
```

### 2. Start New Cascade Layer (R410 Automation)

**Function: `start_new_cascade_layer()`**

```bash
start_new_cascade_layer() {
    local build_errors="$1"

    # Get current layer ID
    CURRENT_LAYER=$(jq -r '.cascade_layers | length' orchestrator-state-v3.json)
    NEW_LAYER=$((CURRENT_LAYER + 1))

    echo "🆕 Starting CASCADE LAYER $NEW_LAYER"

    # 1. DOCUMENT NEW BUGS in bug_registry
    document_bugs_in_registry "$build_errors" "$NEW_LAYER"

    # 2. CREATE LAYER METADATA
    jq --arg layer "$NEW_LAYER" \
       --arg trigger "build_failures_from_layer_${CURRENT_LAYER}_reintegration" \
       --arg started "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
       '.cascade_layers += [{
           layer_id: ($layer | tonumber),
           started_at: $started,
           trigger: $trigger,
           bugs: [],
           status: "fixing_upstream",
           progress: "documenting bugs and starting fixes"
       }]' orchestrator-state-v3.json > /tmp/state.json

    mv /tmp/state.json orchestrator-state-v3.json

    # 3. PAUSE CURRENT LAYER
    if [[ $CURRENT_LAYER -gt 0 ]]; then
        jq --argjson idx $((CURRENT_LAYER - 1)) \
           --arg pause_status "paused_for_layer_${NEW_LAYER}" \
           --arg progress "integration build failed, started layer ${NEW_LAYER} for fixes" \
           '.cascade_layers[$idx].status = $pause_status |
            .cascade_layers[$idx].progress = $progress' \
           orchestrator-state-v3.json > /tmp/state.json

        mv /tmp/state.json orchestrator-state-v3.json
    fi

    # 4. TRANSITION TO ERROR_RECOVERY to start fixes
    echo "🔄 Transitioning to ERROR_RECOVERY to apply layer $NEW_LAYER fixes"
    update_state "ERROR_RECOVERY" "Starting cascade layer $NEW_LAYER fixes"

    # 5. SET CONTINUATION FLAG = TRUE!
    # This is AUTOMATED workflow - NOT a blocker!
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
}
```

### 3. ERROR_RECOVERY: Layer-Aware Fix Application

**Enhancement to ERROR_RECOVERY rules:**

```bash
# In ERROR_RECOVERY state
ACTIVE_LAYER=$(jq -r '.cascade_layers | map(select(.status == "fixing_upstream")) | .[0].layer_id // 0' orchestrator-state-v3.json)

if [[ $ACTIVE_LAYER -gt 0 ]]; then
    echo "🔧 R410: Applying fixes for CASCADE LAYER $ACTIVE_LAYER"

    # Get bugs for this layer
    LAYER_BUGS=$(jq -r ".cascade_layers[] | select(.layer_id == $ACTIVE_LAYER) | .bugs[]" orchestrator-state-v3.json)

    echo "📋 Fixing bugs: $LAYER_BUGS"

    # Apply fixes to upstream effort branches
    # (existing fix application logic)

    # After fixes complete, update layer status
    jq --argjson layer "$ACTIVE_LAYER" \
       '.cascade_layers[] |=
        if .layer_id == $layer then
            .status = "reintegrating" |
            .progress = "fixes applied, starting re-integration"
        else . end' \
       orchestrator-state-v3.json > /tmp/state.json

    mv /tmp/state.json orchestrator-state-v3.json

    # Transition BACK to CASCADE_REINTEGRATION
    # Layer $ACTIVE_LAYER will now re-attempt integration
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # Still automated!
fi
```

### 4. Complete Layer and Resume Previous

**When layer's integration succeeds:**

```bash
# In CASCADE_REINTEGRATION after successful build
ACTIVE_LAYER=$(jq -r '.cascade_layers | map(select(.status == "reintegrating")) | .[0].layer_id // 0' orchestrator-state-v3.json)

if [[ $ACTIVE_LAYER -gt 0 ]]; then
    echo "✅ CASCADE LAYER $ACTIVE_LAYER: Integration successful!"

    # Mark layer complete
    jq --argjson layer "$ACTIVE_LAYER" \
       --arg completed "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
       '.cascade_layers[] |=
        if .layer_id == $layer then
            .status = "completed" |
            .completed_at = $completed |
            .progress = "integration successful, layer complete"
        else . end' \
       orchestrator-state-v3.json > /tmp/state.json

    mv /tmp/state.json orchestrator-state-v3.json

    # Check if previous layer exists and is paused
    PREV_LAYER=$((ACTIVE_LAYER - 1))
    if [[ $PREV_LAYER -gt 0 ]]; then
        PREV_STATUS=$(jq -r ".cascade_layers[] | select(.layer_id == $PREV_LAYER) | .status" orchestrator-state-v3.json)

        if [[ "$PREV_STATUS" == "paused_for_layer_${ACTIVE_LAYER}" ]]; then
            echo "🔄 Resuming CASCADE LAYER $PREV_LAYER"

            # Resume previous layer
            jq --argjson layer "$PREV_LAYER" \
               '.cascade_layers[] |=
                if .layer_id == $layer then
                    .status = "reintegrating" |
                    .progress = "resuming after layer completion"
                else . end' \
               orchestrator-state-v3.json > /tmp/state.json

            mv /tmp/state.json orchestrator-state-v3.json

            # Continue cascade for previous layer
            echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
        fi
    else
        echo "🎉 All cascade layers complete!"
        echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
    fi
fi
```

## Prevention of Infinite Loops

### Layer Count Limit
```bash
# Maximum 10 cascade layers before manual review
MAX_CASCADE_LAYERS=10

LAYER_COUNT=$(jq -r '.cascade_layers | length' orchestrator-state-v3.json)
if [[ $LAYER_COUNT -ge $MAX_CASCADE_LAYERS ]]; then
    echo "🚨 R410: Exceeded maximum cascade layers ($MAX_CASCADE_LAYERS)"
    echo "🚨 This indicates systematic design issues"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
    exit 410
fi
```

### Bug Deduplication
```bash
# Ensure new bugs are actually NEW
are_bugs_new() {
    local new_errors="$1"

    # Hash the error messages
    NEW_HASH=$(echo "$new_errors" | sha256sum | cut -d' ' -f1)

    # Check all previous layers
    for layer_bugs in $(jq -r '.cascade_layers[].bugs[]' orchestrator-state-v3.json); do
        EXISTING_BUG=$(jq -r ".bug_registry.bugs[] | select(.id == \"$layer_bugs\") | .error_hash" orchestrator-state-v3.json)

        if [[ "$NEW_HASH" == "$EXISTING_BUG" ]]; then
            echo "❌ Bug hash matches existing bug - NOT new!"
            return 1  # NOT new
        fi
    done

    return 0  # New bugs
}
```

## Automation Continuation Flags

### ALWAYS Use TRUE for New Layers

**R410 RULE: Starting a new cascade layer is AUTOMATED WORKFLOW**

```bash
# ✅ CORRECT: New layer discovery
echo "📋 Build failure discovered 3 new cross-effort bugs"
echo "🔄 Starting CASCADE LAYER 2 to fix upstream branches"
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # Automation continues!

# ❌ WRONG: Treating new layer as blocker
echo "📋 Build failure discovered 3 new bugs"
echo "⚠️ Need manual approval to start cascade layer 2"
echo "CONTINUE-SOFTWARE-FACTORY=FALSE"  # Defeats automation!
```

### Only Use FALSE for Genuine Blocks

**Examples of genuine blocks:**
- Infinite loop detected (same bugs appearing in multiple layers)
- Exceeded maximum layer count (10+ layers = design problem)
- Cannot determine which bugs are new (bug registry corrupted)
- Cannot identify source efforts for bugs (dependency graph broken)

## Integration with Existing Rules

### R327 - Mandatory CASCADE Re-Integration
- R410 EXTENDS R327 by handling NEW bugs discovered during R327 execution
- R327 says "cascade when fixes applied"
- R410 says "and if cascade reveals MORE bugs, start NEW layer"

### R406 - Fix Cascade Tracking
- R410 uses R406's `bug_registry` for bug documentation
- R410 adds `cascade_layers[]` for layer tracking
- Both work together for complete cascade management

### R380 - Fix Registry Management
- R410 relies on R380 for unique bug IDs (BUG-001, BUG-002...)
- R410 adds "layer_id" context to bug records
- Bugs are associated with the layer that discovered them

## Grading Penalties

| Violation | Penalty | Description |
|-----------|---------|-------------|
| Treating cascade build failures as blockers | -100% | Using FALSE when should use TRUE |
| Not starting new layer for new bugs | -100% | Missing the automated protocol |
| Not tracking cascade layers | -50% | Cannot prove recursive solution |
| Not pausing current layer when starting new | -30% | Loses progress context |
| Not resuming previous layers after completion | -50% | Cascade never completes |
| Infinite loop (same bugs in multiple layers) | -100% | Fixes ineffective |

## Examples

### Example 1: Two-Layer Cascade

**Initial State:**
- Phase 1, Wave 2 complete
- Upstream merge introduces bugs
- Start CASCADE LAYER 1

**Layer 1 Execution:**
```
1. Apply fixes to effort-3, effort-5
2. Cascade to phase1-wave2-integration
3. Build phase1-wave2-integration → FAILS (3 new errors)
4. Start CASCADE LAYER 2 (pause layer 1)
```

**Layer 2 Execution:**
```
1. Apply fixes to effort-2, effort-4 (bugs from layer 1 integration)
2. Cascade to phase1-wave2-integration
3. Build phase1-wave2-integration → PROJECT_DONE!
4. Complete layer 2, resume layer 1
```

**Layer 1 Resume:**
```
1. Re-attempt phase1-wave2-integration (now has layer 2 fixes)
2. Build phase1-wave2-integration → PROJECT_DONE!
3. Complete layer 1
4. Cascade complete!
```

**All transitions used CONTINUE-SOFTWARE-FACTORY=TRUE!**

### Example 2: Three-Layer Cascade

**Scenario:** Complex cross-effort dependencies reveal bugs in stages

```
CASCADE LAYER 1 → Build fails → START LAYER 2
CASCADE LAYER 2 → Build fails → START LAYER 3
CASCADE LAYER 3 → Build PROJECT_DONE! → RESUME LAYER 2
CASCADE LAYER 2 → Build PROJECT_DONE! → RESUME LAYER 1
CASCADE LAYER 1 → Build PROJECT_DONE! → COMPLETE!
```

**Total user interventions: ZERO (fully automated!)**

## Summary

**R410 Makes the Orchestrator DECISIVE:**

1. **Build failure during cascade?** → Start new layer (automated!)
2. **New bugs discovered?** → Document and fix (automated!)
3. **Layer complete?** → Resume previous layer (automated!)
4. **All layers complete?** → Cascade done (automated!)

**The ONLY time to use FALSE:**
- Same bugs appearing multiple times (infinite loop)
- Exceeded maximum layers (design problem)
- Cannot determine what's new (system corruption)

**Otherwise: CONTINUE-SOFTWARE-FACTORY=TRUE** (this is designed workflow!)

## Related Rules
- R327: Mandatory CASCADE Re-Integration After Fixes
- R406: Fix Cascade Tracking Protocol
- R380: Fix Registry Management
- R405: Automation Continuation Flag
- R348: Cascade State Transitions
