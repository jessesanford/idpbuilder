# R421 - Flexible Bug Tracking Architecture

## Rule Metadata
- **ID**: R421.0.0
- **Criticality**: 🚨🚨🚨 BLOCKING - FLEXIBLE BUG TRACKING
- **Category**: Bug Management, State Integrity, Universal Tracking
- **Scope**: All Agents, All States
- **Enforcement**: Schema validation, state file integrity
- **Penalty**: -100% for lost bugs, -50% for incorrect source tracking

## Purpose

Extends R406 Fix Cascade Tracking Protocol to support bugs from MULTIPLE SOURCES while maintaining efficiency for cascade bugs (which remain the most common case). Enables tracking bugs from implementation, code review, user feedback, documentation, and testing phases.

## Relationship to R406

**R406 remains the foundation** for cascade bug tracking:
- All cascade logic UNCHANGED
- Layered cascade protocol INTACT (R410)
- CASCADE_REINTEGRATION workflow PRESERVED
- Automatic status reporting MAINTAINED

**R421 adds flexibility**:
- Non-cascade bugs supported
- Simpler bug creation for non-integration issues
- Unified query interface
- Single source of truth for ALL bugs

## The Gap R421 Fills

### Before R421 (CASCADE-ONLY)

```json
{
  "bug_registry": [
    {
      "bug_id": "BUG-cascade-20251005-051000-001",
      "cascade_id": "cascade-20251005-051000",  // REQUIRED
      "detected_in_integration": {...}           // REQUIRED
      // Every bug must be cascade-related
    }
  ]
}
```

**Problems:**
- ❌ Can't track bugs found during initial implementation
- ❌ Can't track bugs found during code review (before integration)
- ❌ Can't track user feedback bugs (post-deployment)
- ❌ Can't track documentation bugs
- ❌ Forces all bugs through cascade workflow even when inappropriate

### After R421 (FLEXIBLE)

```json
{
  "bug_registry": [
    {
      "bug_id": "BUG-cascade-20251005-051000-001",
      "bug_source": "cascade",
      "cascade_id": "cascade-20251005-051000",  // Non-null for cascade
      "detected_in_integration": {...}
      // CASCADE BUG - Full cascade metadata
    },
    {
      "bug_id": "BUG-20251005-001",
      "bug_source": "implementation",
      "cascade_id": null,                        // Null for non-cascade
      "detected_in_integration": null,
      "description": "Found null pointer during coding"
      // SIMPLE BUG - No cascade overhead
    },
    {
      "bug_id": "BUG-20251005-002",
      "bug_source": "user_feedback",
      "cascade_id": null,
      "description": "User reports 500 error"
      // USER FEEDBACK - Tracked independently
    }
  ]
}
```

**Benefits:**
- ✅ Track bugs from ANY source
- ✅ Cascade bugs keep full metadata
- ✅ Non-cascade bugs stay simple
- ✅ Unified registry for all bugs
- ✅ Flexible intake processes

## Bug Sources

### 1. CASCADE (Most Common - 80%+ of bugs)

**When**: Integration build/test failures during CASCADE_REINTEGRATION

**Characteristics:**
- Cross-effort bugs invisible until integration
- Requires cascade metadata (cascade_id, layer, integration)
- Follows R406 + R410 protocols
- Automatic reintegration workflow

**Example:**
```bash
# Build failed during wave integration
utilities/create-bug.sh \
    --source cascade \
    --cascade-id cascade-20251005-051000 \
    --layer 1 \
    --integration phase1_wave2 \
    --integration-type wave \
    --phase 1 \
    --wave 2 \
    --severity CRITICAL \
    --category build_failure \
    --description "PushCmd redeclared in multiple packages" \
    --efforts "E1.2.1,E1.2.2"
```

**Bug ID Format**: `BUG-cascade-YYYYMMDD-HHMMSS-NNN`

### 2. IMPLEMENTATION (Found During Coding)

**When**: SW Engineer discovers bug while implementing feature

**Characteristics:**
- Found before any integration attempt
- Caught during development
- Can be fixed immediately or queued
- No cascade overhead needed

**Example:**
```bash
# Engineer discovers issue during coding
utilities/create-bug.sh \
    --source implementation \
    --severity HIGH \
    --category runtime_error \
    --description "Null pointer in auth.Validate()" \
    --efforts "E1.2.1" \
    --detected-by "SW Engineer" \
    --file-path "pkg/auth/validator.go" \
    --line-number 42
```

**Bug ID Format**: `BUG-YYYYMMDD-NNN`

### 3. REVIEW (Found During Code Review)

**When**: Code Reviewer finds issues during review

**Characteristics:**
- Logic errors, design issues
- Found before integration
- Part of review feedback
- May require rework

**Example:**
```bash
# Code reviewer finds issue
utilities/create-bug.sh \
    --source review \
    --severity MEDIUM \
    --category other \
    --description "Inconsistent error handling in auth package" \
    --efforts "E1.3.1" \
    --detected-by "Code Reviewer"
```

**Bug ID Format**: `BUG-YYYYMMDD-NNN`

### 4. USER_FEEDBACK (Post-Deployment Reports)

**When**: User reports issue with delivered software

**Characteristics:**
- Discovered in production/staging
- User-reported symptoms
- May affect multiple efforts
- Requires investigation to identify root cause

**Example:**
```bash
# User reports issue
utilities/create-bug.sh \
    --source user_feedback \
    --severity HIGH \
    --category runtime_error \
    --description "Push command fails with 500 error for files >10MB" \
    --efforts "E1.2.3" \
    --detected-by "User Report #42 (GitHub Issue)"
```

**Bug ID Format**: `BUG-YYYYMMDD-NNN`

### 5. DOCUMENTATION (Docs Issues)

**When**: Errors, typos, outdated information in documentation

**Characteristics:**
- Non-code bugs
- Important for usability
- Usually lower severity
- Can be batched

**Example:**
```bash
# Documentation error found
utilities/create-bug.sh \
    --source documentation \
    --severity LOW \
    --category other \
    --description "README example uses deprecated API" \
    --efforts "E2.1.1" \
    --detected-by "Documentation Review"
```

**Bug ID Format**: `BUG-YYYYMMDD-NNN`

### 6. TEST (Found During Testing)

**When**: Test suite execution reveals issues

**Characteristics:**
- Found during test writing or execution
- Before integration (unit/integration tests)
- Clear reproduction steps
- Test failure logs available

**Example:**
```bash
# Test failure found
utilities/create-bug.sh \
    --source test \
    --severity CRITICAL \
    --category test_failure \
    --description "TestPushLargeFile fails with timeout" \
    --efforts "E1.2.3" \
    --detected-by "Test Suite"
```

**Bug ID Format**: `BUG-YYYYMMDD-NNN`

### 7. OTHER (Miscellaneous)

**When**: Any other bug source not covered above

**Bug ID Format**: `BUG-YYYYMMDD-NNN`

## Schema Changes

### Updated bug_registry Structure

```json
{
  "bug_registry": {
    "type": "array",
    "description": "R421 - Single source of truth for ALL bugs from ANY source",
    "items": {
      "required": [
        "bug_id",
        "bug_source",        // NEW: Source type
        "detected_at",
        "severity",
        "description",
        "affected_efforts",
        "fix_status"
      ],
      "properties": {
        "bug_id": {
          "pattern": "^BUG-(cascade-\\d{8}-\\d{6}-\\d{3}|\\d{8}-\\d{3})$"
          // UPDATED: Supports both cascade and non-cascade formats
        },
        "bug_source": {
          "enum": ["cascade", "implementation", "review", "user_feedback", "documentation", "test", "other"]
          // NEW: Bug source type
        },
        "cascade_id": {
          "type": ["string", "null"]
          // UPDATED: Nullable (null for non-cascade bugs)
        },
        "cascade_layer": {
          "type": ["integer", "null"]
          // UPDATED: Nullable (null for non-cascade bugs)
        },
        "detected_in_integration": {
          "type": ["object", "null"]
          // UPDATED: Nullable (null for non-cascade bugs)
        }
        // All other fields remain the same
      }
    }
  }
}
```

### Backward Compatibility

**Existing cascade bugs fully compatible:**
- All R406 queries work unchanged
- R410 layered cascade logic intact
- CASCADE_REINTEGRATION flow unchanged

**Migration not required:**
- Old bugs implicitly have `bug_source: "cascade"`
- Schema validation ensures consistency
- New bugs use explicit source field

## Bug Creation Workflows

### CASCADE BUG (R406 Protocol)

**Trigger**: Integration build/test failure during CASCADE_REINTEGRATION

**Process**:
```bash
# 1. Detect build failure
if ! make build; then
    # 2. Initialize cascade (if not active)
    CASCADE_ID="${CASCADE_ID:-cascade-$(date +%Y%m%d-%H%M%S)}"

    # 3. Create cascade bug
    BUG_ID=$(utilities/create-bug.sh \
        --source cascade \
        --cascade-id "$CASCADE_ID" \
        --layer "$CASCADE_LAYER" \
        --integration "$INTEGRATE_WAVE_EFFORTS_NAME" \
        --integration-type "$INTEGRATE_WAVE_EFFORTS_TYPE" \
        --phase "$PHASE" \
        --wave "$WAVE" \
        --severity CRITICAL \
        --category build_failure \
        --description "$BUILD_ERROR" \
        --efforts "$AFFECTED_EFFORTS")

    # 4. Add to cascade tracking
    # (R406 cascade_layers, integration_fix_states, effort_fix_states)

    # 5. Spawn engineers to fix
    # (R406 protocol)
fi
```

### NON-CASCADE BUG (New R421 Protocol)

**Trigger**: Bug found outside integration context

**Process**:
```bash
# 1. Discover bug (during implementation, review, user report, etc.)
# 2. Create simple bug entry
BUG_ID=$(utilities/create-bug.sh \
    --source "$SOURCE" \
    --severity "$SEVERITY" \
    --description "$DESCRIPTION" \
    --efforts "$AFFECTED_EFFORTS")

# 3. Handle based on context:
#    a) Fix immediately (implementation/review bugs)
#    b) Queue for next sprint (user feedback)
#    c) Add to backlog (documentation/low priority)

# 4. Update bug status as fix progresses
jq --arg bid "$BUG_ID" \
   '(.bug_registry[] | select(.bug_id == $bid) | .fix_status) = "in_progress"' \
   orchestrator-state-v3.json > tmp.json
mv tmp.json orchestrator-state-v3.json
```

## Query Patterns

### Query All Bugs (Any Source)

```bash
# Get all pending bugs regardless of source
jq '.bug_registry[] | select(.fix_status == "pending")' \
   orchestrator-state-v3.json
```

### Query CASCADE Bugs (R406 Queries Unchanged)

```bash
# Get bugs for specific cascade
jq --arg cid "cascade-20251005-051000" \
   '.bug_registry[] | select(.cascade_id == $cid)' \
   orchestrator-state-v3.json

# Get bugs for cascade layer
jq --argjson layer 2 \
   '.bug_registry[] | select(.cascade_layer == $layer)' \
   orchestrator-state-v3.json

# Get bugs for integration
jq --arg integration "phase1_wave2" \
   '.bug_registry[] | select(.detected_in_integration.name == $integration)' \
   orchestrator-state-v3.json
```

### Query Non-CASCADE Bugs (New)

```bash
# Get all user feedback bugs
jq '.bug_registry[] | select(.bug_source == "user_feedback")' \
   orchestrator-state-v3.json

# Get all implementation bugs
jq '.bug_registry[] | select(.bug_source == "implementation")' \
   orchestrator-state-v3.json

# Get all review bugs
jq '.bug_registry[] | select(.bug_source == "review")' \
   orchestrator-state-v3.json

# Get bugs for specific effort (any source)
jq --arg effort "E1.2.1" \
   '.bug_registry[] | select(.affected_efforts[] == $effort)' \
   orchestrator-state-v3.json
```

### Query By Source Distribution

```bash
# Count bugs by source
jq -r '[.bug_registry[] | .bug_source] |
        group_by(.) |
        map({source: .[0], count: length}) |
        .[]' \
   orchestrator-state-v3.json
```

## State Integration

### Where Non-Cascade Bugs Can Be Created

#### 1. IMPLEMENTATION State (SW Engineer)

**When**: Engineer discovers bug while coding

```bash
# In sw-engineer IMPLEMENTATION state
if bug_discovered_during_coding; then
    BUG_ID=$(utilities/create-bug.sh \
        --source implementation \
        --severity "$SEVERITY" \
        --description "$BUG_DESC" \
        --efforts "$CURRENT_EFFORT")

    # Fix immediately or queue
    if can_fix_immediately; then
        apply_fix
        update_bug_status "$BUG_ID" "fixed"
    else
        echo "⚠️ Bug $BUG_ID queued for separate fix cycle"
    fi
fi
```

#### 2. CODE_REVIEW State (Code Reviewer)

**When**: Reviewer finds issues during review

```bash
# In code-reviewer CODE_REVIEW state
if issues_found_in_review; then
    for issue in "${review_issues[@]}"; do
        BUG_ID=$(utilities/create-bug.sh \
            --source review \
            --severity "$ISSUE_SEVERITY" \
            --description "$ISSUE_DESC" \
            --efforts "$REVIEWED_EFFORT")

        # Add to review feedback
        add_to_review_report "$BUG_ID" "$ISSUE_DESC"
    done
fi
```

#### 3. Any State (User Feedback Intake)

**When**: User reports issue via GitHub, email, etc.

```bash
# User feedback intake process
# (Could be orchestrator or dedicated intake state)
process_user_feedback() {
    local github_issue="$1"

    # Parse user report
    DESCRIPTION=$(parse_issue_description "$github_issue")
    SEVERITY=$(determine_severity "$github_issue")
    AFFECTED_EFFORTS=$(analyze_affected_code "$github_issue")

    # Create bug
    BUG_ID=$(utilities/create-bug.sh \
        --source user_feedback \
        --severity "$SEVERITY" \
        --description "$DESCRIPTION" \
        --efforts "$AFFECTED_EFFORTS" \
        --detected-by "GitHub Issue #$github_issue")

    # Route for fixing
    if [[ "$SEVERITY" == "CRITICAL" ]]; then
        # Immediate fix cycle
        route_to_error_recovery "$BUG_ID"
    else
        # Queue for next sprint
        add_to_backlog "$BUG_ID"
    fi
}
```

#### 4. ERROR_RECOVERY State (Any Bug Source)

**When**: Handling bugs that need immediate attention

```bash
# ERROR_RECOVERY can handle ANY bug source
handle_bug_fix() {
    local bug_id="$1"

    # Get bug details
    BUG_SOURCE=$(jq -r --arg bid "$bug_id" \
        '.bug_registry[] | select(.bug_id == $bid) | .bug_source' \
        orchestrator-state-v3.json)

    case "$BUG_SOURCE" in
        cascade)
            # Use R406 cascade protocol
            handle_cascade_bug "$bug_id"
            ;;
        implementation|review|user_feedback|test|other)
            # Use simple fix protocol
            handle_simple_bug "$bug_id"
            ;;
    esac
}

handle_simple_bug() {
    local bug_id="$1"

    # Spawn engineer to fix
    spawn_engineer_for_fix "$bug_id"

    # Wait for fix
    monitor_fix_progress "$bug_id"

    # Verify fix
    verify_fix "$bug_id"

    # No cascade needed (not from integration)
}
```

## Utility: utilities/create-bug.sh

**Location**: `/home/vscode/software-factory-template/utilities/create-bug.sh`

**Purpose**: Universal bug creation supporting both cascade and non-cascade bugs

**Usage**:
```bash
# CASCADE BUG
utilities/create-bug.sh \
    --source cascade \
    --cascade-id CASCADE_ID \
    --layer LAYER \
    --integration NAME \
    --integration-type TYPE \
    --severity SEVERITY \
    --description DESC \
    --efforts EFFORTS

# NON-CASCADE BUG
utilities/create-bug.sh \
    --source SOURCE \
    --severity SEVERITY \
    --description DESC \
    --efforts EFFORTS
```

**See**: `utilities/create-bug.sh --help` for full options

## Validation and Integrity

### Schema Validation

**Required fields for ALL bugs:**
- `bug_id` (proper format for source type)
- `bug_source` (valid enum value)
- `detected_at`
- `severity`
- `description`
- `affected_efforts`
- `fix_status`

**Conditional requirements:**
- CASCADE bugs: `cascade_id`, `cascade_layer`, `detected_in_integration` must be non-null
- Non-cascade bugs: `cascade_id`, `cascade_layer`, `detected_in_integration` must be null

### Integrity Checks

```bash
# Validate cascade bugs have cascade metadata
validate_cascade_bugs() {
    local invalid=$(jq -r '[.bug_registry[] |
        select(.bug_source == "cascade" and
               (.cascade_id == null or
                .cascade_layer == null or
                .detected_in_integration == null))] |
        length' orchestrator-state-v3.json)

    if [[ "$invalid" -gt 0 ]]; then
        echo "❌ $invalid cascade bugs missing cascade metadata!"
        return 1
    fi
}

# Validate non-cascade bugs don't have cascade metadata
validate_noncascade_bugs() {
    local invalid=$(jq -r '[.bug_registry[] |
        select(.bug_source != "cascade" and
               (.cascade_id != null or
                .cascade_layer != null or
                .detected_in_integration != null))] |
        length' orchestrator-state-v3.json)

    if [[ "$invalid" -gt 0 ]]; then
        echo "❌ $invalid non-cascade bugs have cascade metadata!"
        return 1
    fi
}
```

## Reporting and Monitoring

### Universal Bug Report

```bash
# Report ALL bugs by source
bug_report_all_sources() {
    echo "📊 BUG REGISTRY REPORT"
    echo "===================="

    # Total bugs
    local total=$(jq '.bug_registry | length' orchestrator-state-v3.json)
    echo "Total Bugs: $total"
    echo ""

    # By source
    echo "By Source:"
    jq -r '[.bug_registry[] | .bug_source] |
           group_by(.) |
           map("  \(.[0]): \(length)") |
           .[]' orchestrator-state-v3.json
    echo ""

    # By status
    echo "By Status:"
    jq -r '[.bug_registry[] | .fix_status] |
           group_by(.) |
           map("  \(.[0]): \(length)") |
           .[]' orchestrator-state-v3.json
    echo ""

    # By severity
    echo "By Severity:"
    jq -r '[.bug_registry[] | .severity] |
           group_by(.) |
           map("  \(.[0]): \(length)") |
           .[]' orchestrator-state-v3.json
}
```

### Cascade Status Report (R406 - Unchanged)

**The R406 automatic cascade status reporting continues unchanged:**
```bash
# R406 cascade reports still work
if [ -f "orchestrator-state-v3.json" ] && \
   jq -e '.fix_cascade_state.active == true' orchestrator-state-v3.json > /dev/null 2>&1; then
    source utilities/cascade-status-report.sh
    cascade_status_report
fi
```

## Best Practices

### 1. Choose Correct Source

**Decision Tree**:
```
Is bug from integration build/test failure?
├─ YES → source=cascade (use full R406 protocol)
└─ NO → Ask: Where was it found?
    ├─ During coding → source=implementation
    ├─ During review → source=review
    ├─ User reported → source=user_feedback
    ├─ In docs → source=documentation
    ├─ In tests → source=test
    └─ Other → source=other
```

### 2. Cascade Bugs Get Full Metadata

**Always provide for cascade bugs:**
- `cascade_id`
- `cascade_layer`
- `detected_in_integration` (name, type, phase, wave)

### 3. Non-Cascade Bugs Stay Simple

**Don't add cascade metadata to non-cascade bugs:**
- Keep `cascade_id`, `cascade_layer`, `detected_in_integration` as null
- Simpler = easier to track

### 4. Use Unified Registry

**All bugs in bug_registry:**
- Don't create separate tracking for non-cascade bugs
- Single source of truth
- Unified queries
- Consistent status tracking

### 5. Update Status Consistently

**Regardless of source, update status:**
```bash
# Status transitions work for all bugs
jq --arg bid "$BUG_ID" \
   '(.bug_registry[] | select(.bug_id == $bid) | .fix_status) = "fixed"' \
   orchestrator-state-v3.json > tmp.json
```

## Migration Strategy

### Existing Systems (R406 Only)

**No migration needed:**
1. Existing bugs remain valid (implicitly cascade source)
2. New field `bug_source` defaults to "cascade" for existing bugs
3. Continue using R406 for cascade bugs
4. Add R421 for new bug sources as needed

### Gradual Adoption

**Phase 1: Cascade Only (Current)**
- All bugs from integration failures
- R406 protocol exclusively

**Phase 2: Add Implementation/Review Bugs**
- SW Engineers report bugs during coding
- Code Reviewers report bugs during review
- Simple bug creation, immediate fix

**Phase 3: Add User Feedback**
- User reports tracked in bug_registry
- Prioritization and routing
- Integration with issue tracking

**Phase 4: Full Flexibility**
- All bug sources supported
- Unified tracking and reporting
- Comprehensive bug management

## Associated Rules

### MANDATORY DEPENDENCIES
- **R406**: Fix Cascade Tracking Protocol (cascade bugs foundation)
- **R410**: Layered Cascade Protocol (cascade bug handling)
- **R327**: Mandatory Re-Integration After Fixes (why cascades exist)
- **R407**: Mandatory State File Validation (schema enforcement)

### RELATED RULES
- **R348**: Cascade State Transitions
- **R350**: Complete Cascade Dependency Graph
- **R380**: Fix Registry Management

## Grading Penalties

- **Lost bugs** (bug not in registry): -100% IMMEDIATE FAILURE
- **Wrong bug source** (cascade marked as implementation, etc.): -50%
- **Missing cascade metadata** (cascade bug with null fields): -75%
- **Invalid cascade metadata** (non-cascade bug with cascade fields): -50%
- **Inconsistent status tracking**: -30%
- **Missing bug_source field**: -40%

## Summary

R421 extends R406 to support bugs from multiple sources:

✅ **CASCADE BUGS** (80%+ of bugs)
- Full R406 + R410 protocol
- Complete cascade metadata
- Layered reintegration support

✅ **NON-CASCADE BUGS** (20% of bugs)
- Implementation bugs
- Review bugs
- User feedback bugs
- Documentation bugs
- Test bugs

✅ **UNIFIED TRACKING**
- Single bug_registry
- Consistent status tracking
- Universal queries
- Flexible intake

✅ **BACKWARD COMPATIBLE**
- R406 queries unchanged
- Cascade logic intact
- No migration required

✅ **SIMPLE WHERE SIMPLE NEEDED**
- Cascade bugs: Complex metadata
- Other bugs: Minimal overhead

**Most bugs are still cascade bugs - R421 just adds flexibility for the rest!**
