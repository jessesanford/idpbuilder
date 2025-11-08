# R406 - Fix Cascade Tracking Protocol

## Rule Metadata
- **ID**: R406.0.0
- **Criticality**: 🚨🚨🚨 BLOCKING - FIX CASCADE TRACKING
- **Category**: Fix Management, State Integrity
- **Scope**: Orchestrator, ERROR_RECOVERY, CASCADE_REINTEGRATION
- **Enforcement**: Schema validation, state file integrity
- **Penalty**: -100% for lost fixes, -50% for inconsistent tracking

## Purpose

Provides a BULLETPROOF tracking system for fix cascades across multiple integration layers (wave → phase → project). Ensures NO bugs are lost and ALL fixes are tracked from detection through integration, even across multiple cascade iterations.

## Problem Statement

**The Challenge:**
When integration fails and bugs are found, the system enters a complex cascade:
1. Project integration fails → bugs detected
2. Fixes applied to effort branches → project integration STALE
3. Cascade to Phase 2 → Phase re-integrated → MORE bugs found!
4. Cascade to Wave 2 → Wave re-integrated → MORE bugs found!
5. Now we have 3 layers of bugs to track
6. Fix Wave bugs → re-integrate → cascade back up
7. Fix Phase bugs → re-integrate → cascade back up
8. Fix Project bugs → re-integrate → finally done

**Without proper tracking:**
- Easy to lose track of which bugs belong to which integration
- Unclear what layer of cascade we're on
- Hard to know when cascade is complete
- Fixes can be forgotten during multiple iterations
- No audit trail of what was fixed when

## Data Structure Overview

R406 defines FOUR interconnected tracking structures in `orchestrator-state-v3.json`:

### 1. fix_cascade_state (Overall Orchestration)
Tracks the overall cascade state, layers, and progress.

### 2. bug_registry (Single Source of Truth)
Every bug tracked from detection to integration.

### 3. integration_fix_states (Per Integration)
Tracks bugs and fix status for each integration (wave/phase/project).

### 4. effort_fix_states (Per Effort)
Tracks which bugs are assigned to each effort and their fix progress.

## Design Principles

✅ **UNIFIED**: One structure works for wave/phase/project
✅ **LAYERED**: Track cascade layers (which integration triggered which fixes)
✅ **COMPLETE**: Every bug, every effort, every integration state
✅ **CHECKPOINTED**: Clear status at each step
✅ **DEPENDENCY-AWARE**: Know what needs re-integration after fixes
✅ **ITERATIVE**: Support multiple rounds of fixes
✅ **QUERYABLE**: Easy to answer "what bugs are pending for wave 2?"

## Complete Schema Reference

### 1. fix_cascade_state

```json
{
  "fix_cascade_state": {
    "active": true,
    "cascade_id": "cascade-20251004-160000",
    "triggered_by_integration": {
      "type": "project",
      "phase": 1,
      "wave": null,
      "integration_name": "project_integration"
    },
    "cascade_origin": {
      "integration_name": "phase1_wave2",
      "detected_at": "2025-10-04T15:50:00Z"
    },
    "cascade_chain": [
      {
        "layer": 1,
        "integration_name": "phase1_wave1",
        "type": "wave",
        "status": "complete"
      },
      {
        "layer": 2,
        "integration_name": "phase1_wave2",
        "type": "wave",
        "status": "fixing"
      },
      {
        "layer": 3,
        "integration_name": "phase1_integration",
        "type": "phase",
        "status": "pending"
      }
    ],
    "current_layer": 2,
    "total_layers": 3,
    "status": "fixing",
    "created_at": "2025-10-04T16:00:00Z",
    "updated_at": "2025-10-04T16:30:00Z",
    "validation": {
      "total_bugs_detected": 5,
      "total_bugs_fixed": 2,
      "total_bugs_pending": 3,
      "checksum": "md5(sorted_bug_ids)",
      "last_validated": "2025-10-04T16:30:00Z"
    }
  }
}
```

**Status Values:**
- `detecting`: Scanning for bugs
- `planning`: Creating fix plans
- `fixing`: Fixes in progress
- `reintegrating`: Recreating integrations
- `cascading_up`: Moving to next layer
- `complete`: All layers done
- `aborted`: Manual abort

**Layer Status Values:**
- `pending`: Not started yet
- `detecting`: Finding bugs
- `fixing`: Applying fixes
- `reintegrating`: Recreating integration
- `complete`: Layer finished
- `failed`: Layer failed

### 2. bug_registry

```json
{
  "bug_registry": [
    {
      "bug_id": "BUG-cascade-20251004-160000-001",
      "cascade_id": "cascade-20251004-160000",
      "cascade_layer": 2,
      "detected_in_integration": {
        "name": "phase1_wave2",
        "type": "wave",
        "phase": 1,
        "wave": 2
      },
      "detected_at": "2025-10-04T15:50:00Z",
      "detected_by": "R291 Build Gate",
      "severity": "CRITICAL",
      "category": "build_failure",
      "description": "PushCmd redeclared in multiple packages",
      "error_details": {
        "error_message": "PushCmd redeclared in this block",
        "file_path": "cmd/push.go",
        "line_number": 45,
        "stack_trace": "..."
      },
      "affected_efforts": ["E1.2.1-command-structure"],
      "primary_effort": "E1.2.1-command-structure",
      "requires_coordinated_fix": false,
      "fix_status": "fixed",
      "fix_attempts": [
        {
          "attempt_number": 1,
          "started_at": "2025-10-04T16:00:00Z",
          "completed_at": "2025-10-04T16:15:00Z",
          "commits": ["abc123def"],
          "outcome": "successful",
          "verified": true,
          "notes": "Removed duplicate PushCmd declaration from cmd package"
        }
      ],
      "current_attempt": 1,
      "integration_status": "not_integrated",
      "integrated_at": null,
      "integration_commit": null,
      "blocked_by": [],
      "blocks": [],
      "resolution_notes": "Fixed by consolidating command declarations in single file",
      "created_at": "2025-10-04T15:50:00Z",
      "updated_at": "2025-10-04T16:15:00Z"
    }
  ]
}
```

**Severity Values:**
- `CRITICAL`: Blocks all progress (build failure, crash)
- `HIGH`: Major functionality broken
- `MEDIUM`: Partial functionality affected
- `LOW`: Minor issues, edge cases

**Category Values:**
- `build_failure`: Code doesn't compile
- `test_failure`: Tests failing
- `lint_error`: Linting issues
- `runtime_error`: Crashes, panics
- `integration_conflict`: Merge conflicts, interface mismatches
- `other`: Other issues

**Fix Status Values:**
- `pending`: Not started
- `in_progress`: Being fixed
- `fixed`: Fix complete
- `verified`: Fix tested and verified
- `integrated`: Fix in integration branch
- `blocked`: Waiting on other fixes
- `abandoned`: Fix cancelled

**Integration Status Values:**
- `not_integrated`: Fix not in any integration
- `integrated`: Fix in integration branch
- `failed_integration`: Integration attempt failed

### 3. integration_fix_states

```json
{
  "integration_fix_states": {
    "phase1_wave2": {
      "integration_type": "wave",
      "phase": 1,
      "wave": 2,
      "integration_name": "phase1_wave2",
      "cascade_layer": 2,
      "cascade_id": "cascade-20251004-160000",
      "bugs_detected": [
        "BUG-cascade-20251004-160000-001",
        "BUG-cascade-20251004-160000-002"
      ],
      "bugs_by_status": {
        "pending": 0,
        "in_progress": 0,
        "fixed": 2,
        "verified": 2,
        "integrated": 0
      },
      "bugs_total_count": 2,
      "status": "ready_for_reintegration",
      "dependent_integrations": [
        "phase1_integration",
        "project_integration"
      ],
      "requires_reintegration": true,
      "reintegration_attempts": [
        {
          "attempt_number": 1,
          "started_at": "2025-10-04T16:20:00Z",
          "completed_at": null,
          "outcome": "in_progress",
          "new_bugs_found": []
        }
      ],
      "reintegration_complete": false,
      "created_at": "2025-10-04T15:50:00Z",
      "updated_at": "2025-10-04T16:20:00Z"
    }
  }
}
```

**Status Values:**
- `detecting`: Finding bugs in this integration
- `planning`: Creating fix plan
- `fixing`: Fixes being applied to source efforts
- `ready_for_reintegration`: All fixes complete, ready to recreate integration
- `reintegrating`: Integration being recreated
- `complete`: Integration recreated and verified
- `failed`: Integration recreation failed

### 4. effort_fix_states

```json
{
  "effort_fix_states": {
    "E1.2.1-command-structure": {
      "effort_id": "E1.2.1-command-structure",
      "phase": 1,
      "wave": 2,
      "bugs_assigned": [
        "BUG-cascade-20251004-160000-001"
      ],
      "bugs_by_status": {
        "pending": 0,
        "in_progress": 0,
        "fixed": 1,
        "verified": 1
      },
      "fixes_in_progress": [],
      "fixes_complete": [
        "BUG-cascade-20251004-160000-001"
      ],
      "ready_for_integration": true,
      "last_fix_commit": "abc123def",
      "fix_branch": "E1.2.1-command-structure",
      "created_at": "2025-10-04T16:00:00Z",
      "updated_at": "2025-10-04T16:15:00Z"
    }
  }
}
```

## Common Query Patterns

### Query 1: What bugs are pending for wave 2?
```bash
jq -r '.bug_registry[] |
  select(.detected_in_integration.name == "phase1_wave2" and
         .fix_status == "pending")' \
  orchestrator-state-v3.json
```

### Query 2: Are all bugs fixed for this integration?
```bash
integration="phase1_wave2"
pending=$(jq -r ".integration_fix_states.\"$integration\".bugs_by_status.pending" orchestrator-state-v3.json)
in_progress=$(jq -r ".integration_fix_states.\"$integration\".bugs_by_status.in_progress" orchestrator-state-v3.json)

if [[ "$pending" -eq 0 ]] && [[ "$in_progress" -eq 0 ]]; then
  echo "✅ All bugs fixed for $integration"
else
  echo "⏳ $pending pending, $in_progress in progress"
fi
```

### Query 3: What efforts need fixes?
```bash
jq -r '.effort_fix_states |
  to_entries[] |
  select(.value.bugs_by_status.pending > 0 or
         .value.bugs_by_status.in_progress > 0) |
  .key' \
  orchestrator-state-v3.json
```

### Query 4: What's blocking this integration from reintegration?
```bash
integration="phase1_wave2"
jq -r ".bug_registry[] |
  select(.detected_in_integration.name == \"$integration\" and
         .fix_status != \"verified\")" \
  orchestrator-state-v3.json
```

### Query 5: What's the overall cascade status?
```bash
jq -r '.fix_cascade_state |
  "Status: \(.status)\nLayer: \(.current_layer)/\(.total_layers)\n" +
  (.cascade_chain[] | "  Layer \(.layer): \(.integration_name) - \(.status)")' \
  orchestrator-state-v3.json
```

### Query 6: Validation - Are all bugs accounted for?
```bash
total_in_registry=$(jq -r '.bug_registry | length' orchestrator-state-v3.json)
total_in_validation=$(jq -r '.fix_cascade_state.validation.total_bugs_detected' orchestrator-state-v3.json)

if [[ "$total_in_registry" -eq "$total_in_validation" ]]; then
  echo "✅ All bugs accounted for: $total_in_registry"
else
  echo "❌ MISMATCH: Registry=$total_in_registry, Validation=$total_in_validation"
fi
```

## State Transition Protocols

### Initializing a Cascade

When ERROR_RECOVERY detects fixes are needed:

```bash
# 1. Create cascade ID
CASCADE_ID="cascade-$(date +%Y%m%d-%H%M%S)"

# 2. Initialize cascade state
jq --arg cid "$CASCADE_ID" \
   --arg integration "$INTEGRATE_WAVE_EFFORTS_NAME" \
   --arg type "$INTEGRATE_WAVE_EFFORTS_TYPE" \
   '.fix_cascade_state = {
      active: true,
      cascade_id: $cid,
      triggered_by_integration: {
        type: $type,
        integration_name: $integration
      },
      cascade_origin: {
        integration_name: $integration,
        detected_at: (now | todate)
      },
      cascade_chain: [],
      current_layer: 0,
      total_layers: 0,
      status: "detecting",
      created_at: (now | todate),
      updated_at: (now | todate),
      validation: {
        total_bugs_detected: 0,
        total_bugs_fixed: 0,
        total_bugs_pending: 0,
        checksum: "",
        last_validated: (now | todate)
      }
    }' orchestrator-state-v3.json > tmp.json
mv tmp.json orchestrator-state-v3.json
```

### Registering a Bug

When a bug is detected:

```bash
# Generate bug ID
BUG_NUMBER=$(printf "%03d" $(($(jq '.bug_registry | length' orchestrator-state-v3.json) + 1)))
BUG_ID="BUG-${CASCADE_ID}-${BUG_NUMBER}"

# Add to bug registry
jq --arg bid "$BUG_ID" \
   --arg cid "$CASCADE_ID" \
   --argjson layer "$CASCADE_LAYER" \
   --arg integration "$INTEGRATE_WAVE_EFFORTS_NAME" \
   --arg severity "$SEVERITY" \
   --arg category "$CATEGORY" \
   --arg description "$DESCRIPTION" \
   --argjson efforts "$(echo $EFFORTS | jq -R 'split(" ")')" \
   '.bug_registry += [{
      bug_id: $bid,
      cascade_id: $cid,
      cascade_layer: $layer,
      detected_in_integration: {
        name: $integration,
        type: "wave"
      },
      detected_at: (now | todate),
      detected_by: "R291 Build Gate",
      severity: $severity,
      category: $category,
      description: $description,
      affected_efforts: $efforts,
      primary_effort: $efforts[0],
      requires_coordinated_fix: false,
      fix_status: "pending",
      fix_attempts: [],
      current_attempt: 0,
      integration_status: "not_integrated",
      integrated_at: null,
      integration_commit: null,
      blocked_by: [],
      blocks: [],
      resolution_notes: "",
      created_at: (now | todate),
      updated_at: (now | todate)
    }]' orchestrator-state-v3.json > tmp.json
mv tmp.json orchestrator-state-v3.json
```

### Updating Bug Status

When a bug fix starts:

```bash
jq --arg bid "$BUG_ID" \
   '(.bug_registry[] | select(.bug_id == $bid) | .fix_status) = "in_progress" |
    (.bug_registry[] | select(.bug_id == $bid) | .current_attempt) = 1 |
    (.bug_registry[] | select(.bug_id == $bid) | .fix_attempts) += [{
      attempt_number: 1,
      started_at: (now | todate),
      completed_at: null,
      commits: [],
      outcome: "in_progress",
      verified: false,
      notes: ""
    }] |
    (.bug_registry[] | select(.bug_id == $bid) | .updated_at) = (now | todate)' \
   orchestrator-state-v3.json > tmp.json
mv tmp.json orchestrator-state-v3.json
```

When a bug fix completes:

```bash
jq --arg bid "$BUG_ID" \
   --arg commit "$COMMIT_SHA" \
   --arg notes "$FIX_NOTES" \
   '(.bug_registry[] | select(.bug_id == $bid) | .fix_status) = "fixed" |
    (.bug_registry[] | select(.bug_id == $bid) | .fix_attempts[-1].completed_at) = (now | todate) |
    (.bug_registry[] | select(.bug_id == $bid) | .fix_attempts[-1].commits) += [$commit] |
    (.bug_registry[] | select(.bug_id == $bid) | .fix_attempts[-1].outcome) = "successful" |
    (.bug_registry[] | select(.bug_id == $bid) | .fix_attempts[-1].notes) = $notes |
    (.bug_registry[] | select(.bug_id == $bid) | .updated_at) = (now | todate)' \
   orchestrator-state-v3.json > tmp.json
mv tmp.json orchestrator-state-v3.json
```

### Updating Integration Fix State

Track bugs per integration:

```bash
jq --arg integration "$INTEGRATE_WAVE_EFFORTS_NAME" \
   --arg bid "$BUG_ID" \
   '.integration_fix_states[$integration].bugs_detected += [$bid] |
    .integration_fix_states[$integration].bugs_total_count += 1 |
    .integration_fix_states[$integration].bugs_by_status.pending += 1 |
    .integration_fix_states[$integration].updated_at = (now | todate)' \
   orchestrator-state-v3.json > tmp.json
mv tmp.json orchestrator-state-v3.json
```

Update counts when bug status changes:

```bash
# Moving from pending to in_progress
jq --arg integration "$INTEGRATE_WAVE_EFFORTS_NAME" \
   '.integration_fix_states[$integration].bugs_by_status.pending -= 1 |
    .integration_fix_states[$integration].bugs_by_status.in_progress += 1 |
    .integration_fix_states[$integration].updated_at = (now | todate)' \
   orchestrator-state-v3.json > tmp.json
mv tmp.json orchestrator-state-v3.json
```

### Updating Effort Fix State

Assign bug to effort:

```bash
jq --arg effort "$EFFORT_ID" \
   --arg bid "$BUG_ID" \
   '.effort_fix_states[$effort].bugs_assigned += [$bid] |
    .effort_fix_states[$effort].bugs_by_status.pending += 1 |
    .effort_fix_states[$effort].updated_at = (now | todate)' \
   orchestrator-state-v3.json > tmp.json
mv tmp.json orchestrator-state-v3.json
```

Mark fix in progress:

```bash
jq --arg effort "$EFFORT_ID" \
   --arg bid "$BUG_ID" \
   --arg agent "$AGENT_NAME" \
   '.effort_fix_states[$effort].bugs_by_status.pending -= 1 |
    .effort_fix_states[$effort].bugs_by_status.in_progress += 1 |
    .effort_fix_states[$effort].fixes_in_progress += [{
      bug_id: $bid,
      started_at: (now | todate),
      assigned_to_agent: $agent,
      current_attempt: 1
    }] |
    .effort_fix_states[$effort].updated_at = (now | todate)' \
   orchestrator-state-v3.json > tmp.json
mv tmp.json orchestrator-state-v3.json
```

Mark fix complete:

```bash
jq --arg effort "$EFFORT_ID" \
   --arg bid "$BUG_ID" \
   --arg commit "$COMMIT_SHA" \
   '.effort_fix_states[$effort].bugs_by_status.in_progress -= 1 |
    .effort_fix_states[$effort].bugs_by_status.fixed += 1 |
    .effort_fix_states[$effort].fixes_in_progress = [.effort_fix_states[$effort].fixes_in_progress[] | select(.bug_id != $bid)] |
    .effort_fix_states[$effort].fixes_complete += [$bid] |
    .effort_fix_states[$effort].last_fix_commit = $commit |
    .effort_fix_states[$effort].updated_at = (now | todate)' \
   orchestrator-state-v3.json > tmp.json
mv tmp.json orchestrator-state-v3.json
```

### Validation and Checksums

Run validation to ensure no bugs lost:

```bash
#!/bin/bash
# validate-cascade.sh

validate_cascade() {
    local state_file="${1:-orchestrator-state-v3.json}"

    echo "🔍 Validating cascade state..."

    # Count bugs in registry
    local registry_count=$(jq -r '.bug_registry | length' "$state_file")

    # Count bugs in validation
    local validation_count=$(jq -r '.fix_cascade_state.validation.total_bugs_detected' "$state_file")

    # Count by status
    local pending=$(jq -r '[.bug_registry[] | select(.fix_status == "pending")] | length' "$state_file")
    local in_progress=$(jq -r '[.bug_registry[] | select(.fix_status == "in_progress")] | length' "$state_file")
    local fixed=$(jq -r '[.bug_registry[] | select(.fix_status == "fixed")] | length' "$state_file")
    local verified=$(jq -r '[.bug_registry[] | select(.fix_status == "verified")] | length' "$state_file")
    local integrated=$(jq -r '[.bug_registry[] | select(.fix_status == "integrated")] | length' "$state_file")

    echo "📊 Bug Registry Statistics:"
    echo "  Total bugs: $registry_count"
    echo "  Pending: $pending"
    echo "  In Progress: $in_progress"
    echo "  Fixed: $fixed"
    echo "  Verified: $verified"
    echo "  Integrated: $integrated"

    # Validate counts match
    local total_by_status=$((pending + in_progress + fixed + verified + integrated))

    if [[ "$registry_count" -eq "$validation_count" ]] && \
       [[ "$registry_count" -eq "$total_by_status" ]]; then
        echo "✅ Validation PASSED: All bugs accounted for"
        return 0
    else
        echo "❌ Validation FAILED:"
        [[ "$registry_count" -ne "$validation_count" ]] && \
            echo "  - Registry count ($registry_count) != Validation count ($validation_count)"
        [[ "$registry_count" -ne "$total_by_status" ]] && \
            echo "  - Registry count ($registry_count) != Status sum ($total_by_status)"
        return 1
    fi
}

# Generate checksum
generate_checksum() {
    local state_file="${1:-orchestrator-state-v3.json}"

    # Get sorted bug IDs
    local bug_ids=$(jq -r '.bug_registry[].bug_id | sort | @csv' "$state_file")

    # Generate MD5
    local checksum=$(echo "$bug_ids" | md5sum | cut -d' ' -f1)

    echo "$checksum"
}

# Update validation metadata
update_validation() {
    local state_file="${1:-orchestrator-state-v3.json}"

    local total=$(jq -r '.bug_registry | length' "$state_file")
    local fixed=$(jq -r '[.bug_registry[] | select(.fix_status == "fixed" or .fix_status == "verified" or .fix_status == "integrated")] | length' "$state_file")
    local pending=$(jq -r '[.bug_registry[] | select(.fix_status == "pending" or .fix_status == "in_progress")] | length' "$state_file")
    local checksum=$(generate_checksum "$state_file")

    jq --argjson total "$total" \
       --argjson fixed "$fixed" \
       --argjson pending "$pending" \
       --arg checksum "$checksum" \
       '.fix_cascade_state.validation = {
          total_bugs_detected: $total,
          total_bugs_fixed: $fixed,
          total_bugs_pending: $pending,
          checksum: $checksum,
          last_validated: (now | todate)
        }' "$state_file" > tmp.json
    mv tmp.json "$state_file"

    echo "✅ Validation metadata updated"
}

validate_cascade
update_validation
```

## Example Scenarios

### Scenario 1: Wave Integration Fails

```
1. Wave 2 integration built → R291 detects 2 build failures
2. Initialize cascade:
   - cascade_id: "cascade-20251004-160000"
   - triggered_by: phase1_wave2
   - bugs: BUG-001, BUG-002

3. Register bugs in bug_registry:
   - BUG-cascade-20251004-160000-001 (affects E1.2.1)
   - BUG-cascade-20251004-160000-002 (affects E1.2.2)

4. Create integration_fix_states["phase1_wave2"]:
   - bugs_detected: [BUG-001, BUG-002]
   - status: "fixing"

5. Create effort_fix_states:
   - E1.2.1: bugs_assigned: [BUG-001]
   - E1.2.2: bugs_assigned: [BUG-002]

6. Spawn SW Engineers to fix bugs
7. Update bug_registry as fixes complete
8. When all bugs fixed → integration ready for reintegration
9. CASCADE_REINTEGRATION recreates wave integration
10. Mark cascade complete
```

### Scenario 2: Cascading Failures (Multi-Layer)

```
1. Project integration fails → 1 bug found
   - Initialize cascade: cascade-20251004-160000
   - BUG-001 affects E1.1.1

2. Fix E1.1.1 → effort branch updated
3. Detect stale integrations:
   - phase1_wave1 contains E1.1.1 → STALE
   - phase1_integration contains wave1 → STALE
   - project_integration contains phase1 → STALE

4. Build cascade_chain:
   - Layer 1: phase1_wave1 (deepest)
   - Layer 2: phase1_integration
   - Layer 3: project_integration

5. Execute cascade layer by layer:
   - Recreate wave1 → NEW BUGS FOUND! (BUG-002, BUG-003)
   - Add to bug_registry
   - Fix new bugs
   - Recreate phase1 → clean
   - Recreate project → clean

6. Mark all layers complete
7. Exit cascade
```

### Scenario 3: Iterative Fixes

```
1. Bug BUG-001 assigned to E1.2.1
2. Attempt 1:
   - SW Engineer applies fix
   - Commit abc123
   - Reintegration → STILL FAILING
   - Mark attempt 1 as "failed"

3. Attempt 2:
   - SW Engineer applies different fix
   - Commit def456
   - Reintegration → PROJECT_DONE
   - Mark attempt 2 as "successful"

4. Bug history shows both attempts
5. Can analyze why attempt 1 failed
```

## Migration from Old Fields

### Old Fields to Deprecate

These fields should migrate to the new structure:

- `upstream_bugs_wave2` → `bug_registry` + `integration_fix_states`
- `project_fixes_in_progress` → `bug_registry` + `effort_fix_states`
- `efforts_needing_fixes` → derived from `effort_fix_states`
- `fix_cascade` (ad-hoc) → `fix_cascade_state`

### Migration Script

```bash
#!/bin/bash
# migrate-to-r406.sh

migrate_old_bugs() {
    local state_file="${1:-orchestrator-state-v3.json}"

    echo "🔄 Migrating old bug tracking to R406 format..."

    # Check if old fields exist
    if jq -e '.upstream_bugs_wave2' "$state_file" > /dev/null 2>&1; then
        echo "  Found upstream_bugs_wave2, migrating..."

        # Initialize new structures if not exist
        jq '.bug_registry = .bug_registry // [] |
            .integration_fix_states = .integration_fix_states // {} |
            .effort_fix_states = .effort_fix_states // {}' \
            "$state_file" > tmp.json
        mv tmp.json "$state_file"

        # Migrate each old bug (would need custom logic based on structure)
        echo "  ⚠️ Manual migration required for upstream_bugs_wave2"
        echo "  Please review and migrate to bug_registry format"
    fi

    # Similar for other old fields

    echo "✅ Migration check complete"
}

migrate_old_bugs
```

### Compatibility Layer

For backward compatibility during transition:

```bash
# Helper function to check both old and new formats
get_pending_bugs() {
    local integration="$1"
    local state_file="${2:-orchestrator-state-v3.json}"

    # Try new format first
    if jq -e ".integration_fix_states.\"$integration\"" "$state_file" > /dev/null 2>&1; then
        jq -r ".integration_fix_states.\"$integration\".bugs_detected[]" "$state_file"
    # Fall back to old format
    elif jq -e '.upstream_bugs_wave2' "$state_file" > /dev/null 2>&1; then
        jq -r '.upstream_bugs_wave2[]' "$state_file"
    fi
}
```

## Shell Helper Functions

```bash
# Source this file for helper functions
# source: rule-library/helpers/cascade-helpers.sh

# Initialize a new cascade
cascade_init() {
    local integration_name="$1"
    local integration_type="$2"

    local cascade_id="cascade-$(date +%Y%m%d-%H%M%S)"

    jq --arg cid "$cascade_id" \
       --arg integration "$integration_name" \
       --arg type "$integration_type" \
       '.fix_cascade_state = {
          active: true,
          cascade_id: $cid,
          triggered_by_integration: {
            type: $type,
            integration_name: $integration
          },
          status: "detecting",
          created_at: (now | todate),
          updated_at: (now | todate)
        }' orchestrator-state-v3.json > tmp.json
    mv tmp.json orchestrator-state-v3.json

    echo "$cascade_id"
}

# Register a new bug
cascade_register_bug() {
    local cascade_id="$1"
    local integration="$2"
    local severity="$3"
    local category="$4"
    local description="$5"
    shift 5
    local efforts=("$@")

    local bug_num=$(jq '.bug_registry | length' orchestrator-state-v3.json)
    local bug_id=$(printf "BUG-%s-%03d" "$cascade_id" $((bug_num + 1)))

    # Add to registry (implementation as shown above)

    echo "$bug_id"
}

# Get cascade status
cascade_status() {
    jq -r '.fix_cascade_state |
      "Cascade: \(.cascade_id // "none")\n" +
      "Status: \(.status // "inactive")\n" +
      "Layer: \(.current_layer // 0)/\(.total_layers // 0)"' \
      orchestrator-state-v3.json
}

# Check if cascade is complete
cascade_is_complete() {
    local pending=$(jq -r '[.bug_registry[] | select(.fix_status == "pending" or .fix_status == "in_progress")] | length' orchestrator-state-v3.json)

    [[ "$pending" -eq 0 ]]
}
```

## Enforcement

### Schema Validation
The schema enforces:
- Required fields for bugs and states
- Valid enum values for status fields
- Proper ID formats (cascade-YYYYMMDD-HHMMSS)
- Referential integrity (bug_ids must match pattern)

### State File Validation
R407 mandates validation at critical points:
- Before/after state transitions
- Before spawning agents
- After completing fixes

### Grading Penalties

- **Lost bugs** (bug not in registry): -100% IMMEDIATE FAILURE
- **Inconsistent counts** (registry != validation): -50%
- **Missing cascade state** (when cascade active): -75%
- **Wrong status transitions**: -30%
- **Missing timestamps**: -20%

## Associated Rules

### MANDATORY
- **R327**: Mandatory Re-Integration After Fixes (why cascades exist)
- **R348**: Cascade State Transitions (cascade flow)
- **R350**: Complete Cascade Dependency Graph
- **R407**: Mandatory State File Validation

### SUPPORTING
- **R351**: Cascade Execution Protocol
- **R352**: Overlapping Cascade Support
- **R353**: Cascade Focus Protocol
- **R354**: Post-Rebase Review Requirement
- **R380**: Fix Registry Management

## Best Practices

### 1. Always Update Counts
When changing bug status, update ALL related counts:
- bug_registry[].fix_status
- integration_fix_states[].bugs_by_status
- effort_fix_states[].bugs_by_status
- fix_cascade_state.validation counters

### 2. Run Validation Frequently
```bash
# After any bug status change
validate-cascade.sh
```

### 3. Use Atomic Updates
Update all related structures in single jq operation:
```bash
jq '# Update bug status
    (.bug_registry[] | select(.bug_id == $bid) | .fix_status) = "fixed" |
    # Update integration counts
    .integration_fix_states[$integration].bugs_by_status.in_progress -= 1 |
    .integration_fix_states[$integration].bugs_by_status.fixed += 1 |
    # Update effort counts
    .effort_fix_states[$effort].bugs_by_status.in_progress -= 1 |
    .effort_fix_states[$effort].bugs_by_status.fixed += 1' \
    orchestrator-state-v3.json > tmp.json
```

### 4. Maintain Audit Trail
Always populate:
- created_at / updated_at timestamps
- fix_attempts history
- resolution_notes

### 5. Validate Before Exit
Before leaving CASCADE_REINTEGRATION:
```bash
if ! validate-cascade.sh; then
    echo "❌ Cascade validation failed - cannot exit"
    exit 1
fi
```

## 🔴🔴🔴 MANDATORY STATUS REPORTING 🔴🔴🔴

### Automatic Cascade Status Reporting

**ALL orchestrator states MUST report cascade status automatically at end of execution when cascade is active.**

**Grading Penalty:** -15% for missing automatic reports during active cascade

### Implementation

Every orchestrator state that executes during an active cascade MUST include this code **BEFORE** the CONTINUE-SOFTWARE-FACTORY flag:

```bash
# MANDATORY: Print cascade status if active (R406 auto-reporting)
if [ -f "orchestrator-state-v3.json" ] && \
   jq -e '.fix_cascade_state.active == true' orchestrator-state-v3.json > /dev/null 2>&1; then
    echo ""
    echo "═══════════════════════════════════════════════════════════"
    echo "📊 R406 FIX CASCADE STATUS (automatic report)"
    echo "═══════════════════════════════════════════════════════════"
    source utilities/cascade-status-report.sh
    cascade_status_report
    echo "═══════════════════════════════════════════════════════════"
    echo ""
fi

# Then emit continuation flag
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # or FALSE
```

### Required States

The following states MUST include automatic cascade status reporting:

**MANDATORY (High Priority):**
- ERROR_RECOVERY (after fix work)
- CASCADE_REINTEGRATION (during cascade execution)
- SPAWN_SW_ENGINEERS (after spawning engineers)
- MONITORING_EFFORT_REVIEWS (while monitoring fix progress)

**RECOMMENDED (Good Practice):**
- SPAWN_CODE_REVIEWER_FIX_PLAN (before fix planning)
- MONITORING_SWE_PROGRESS (if fixes in progress)
- Any custom states that handle fixes

### Report Utility

**Location:** `utilities/cascade-status-report.sh`

**Function:** `cascade_status_report()`

**Usage:**
```bash
# Source and call
source utilities/cascade-status-report.sh
cascade_status_report

# Or run directly
bash utilities/cascade-status-report.sh
```

### Report Contents

The status report MUST display:

1. **Cascade Overview**
   - Cascade ID
   - Status (detecting/fixing/reintegrating/complete)
   - Current layer / total layers
   - Triggered by which integration

2. **Bug Registry Summary**
   - Total bugs detected
   - Bugs by status (pending/in_progress/fixed/verified)

3. **Integration Status**
   - Per-integration bug breakdown
   - Which integrations need reintegration
   - Integration-specific bug counts

4. **Effort Status**
   - Per-effort bug assignments
   - Which efforts have pending/in-progress/fixed bugs
   - Ready for integration status

5. **Next Actions**
   - What needs to happen next
   - Blocking issues (if any)
   - Estimated progress

6. **Validation**
   - Bug count checksums
   - Integrity verification
   - Last validated timestamp

### User Benefits

**Why automatic reporting is critical:**

✅ **Transparency**: User always knows fix progress
✅ **No Lost Work**: All bugs tracked and visible
✅ **Clear Actions**: User sees what's next
✅ **Trust**: System shows it's working correctly
✅ **Debugging**: Easy to spot stuck/blocked bugs
✅ **Accountability**: Full audit trail

**Without automatic reporting:**
❌ User has to manually ask for status
❌ Easy to lose track of which bugs are fixed
❌ Unclear when cascade will complete
❌ No visibility into blocking issues
❌ Reduced confidence in automation

### Slash Command (Optional)

Users can also request status manually:

```bash
/cascade-status
```

See: `.claude/commands/cascade-status.md`

### Enforcement

**During cascade operations:**
- Every state transition = status report
- User sees progress without asking
- Clear feedback loop maintained

**Grading criteria:**
- Missing status reports: -15%
- Incomplete status info: -10%
- Status not user-friendly: -5%

### Example Output

```
═══════════════════════════════════════════════════════════
📊 R406 FIX CASCADE STATUS (automatic report)
═══════════════════════════════════════════════════════════

Cascade Overview:
  Cascade ID: cascade-20251004-160000
  Status: fixing
  Layer: 1 / 2
  Triggered by: phase1_wave2
  Started: 2025-10-04T16:00:00Z

Bug Registry Summary:
  Total Bugs: 5
  ✅ Verified: 2
  ✅ Fixed: 0
  🔧 In Progress: 1
  ⏳ Pending: 2
  ❌ Blocked: 0

Integration Status:
  phase1_wave1:
    Total Bugs: 2
    ✅ Fixed: 1
    ⏳ Pending: 1
    Status: fixing
    ⚠️ Requires Reintegration: Yes

Effort Status:
  ✅ E1.1.3: 1 bug fixed
  ⏳ E1.1.2-split-002: 1 bug pending
  🔧 E1.2.1-command-structure: 1 bug in progress

Next Actions:
  1. Complete fix for BUG-001 (in progress)
  2. Start fix for BUG-002 (pending)
  3. Start fix for BUG-003 (pending)
  4. Reintegrate Wave 1 after both bugs fixed

Validation:
  ✅ Bug count valid: 5 bugs accounted for
  ✅ Checksum valid: no lost bugs
  Last validated: 2025-10-04T16:30:00Z
═══════════════════════════════════════════════════════════
```

## Summary

R406 provides a BULLETPROOF tracking system for fix cascades that:

✅ Tracks every bug from detection to integration
✅ Works uniformly across wave/phase/project levels
✅ Supports multiple cascade layers
✅ Enables iterative fixes with history
✅ Provides checksums to prevent lost bugs
✅ Makes querying bug status trivial
✅ Integrates with R327 cascade requirements
✅ Enforced by schema validation
✅ **AUTOMATICALLY REPORTS STATUS** at end of every state

**THIS IS CRITICAL INFRASTRUCTURE - TREAT WITH CARE!**

## Terminology Note

R406 uses the term "cascade" to refer to **integration cascade** - the process of rebuilding dependent integrations after bug fixes (e.g., Wave 1 fix → rebuild Wave 2 → rebuild Phase 1).

For terminology clarification:
- **Integration cascade** (this rule): Rebuilding dependent integrations
- **Bug status propagation** (R524): Updating duplicate bug statuses
- **Layered cascade** (R410): Multiple integration cascades in succession

See R530 for complete cascade terminology disambiguation.
