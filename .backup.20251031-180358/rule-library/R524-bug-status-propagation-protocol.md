# 🚨🚨🚨 RULE R524: Bug Status Propagation Protocol

## Criticality: BLOCKING
**Violation Penalty: -75% for manual duplicate fixes, -100% for inconsistent states**

## Description
This rule defines the automatic status propagation that occurs when a canonical bug is fixed: ALL duplicate bugs must be automatically marked as fixed with the same resolution. This prevents wasted effort fixing duplicates manually and ensures consistent bug state across the registry.

## The Problem This Solves

### Without R524 (WASTEFUL):
```
BUG-007 (canonical) is fixed → marked FIXED
BUG-001 (duplicate) → still OPEN (agents waste time fixing it)
BUG-006 (duplicate) → still OPEN (agents waste time fixing it)

Result: Same bug "fixed" 3 separate times
        Triple the effort
        Inconsistent bug registry
```

### With R524 (EFFICIENT):
```
BUG-007 (canonical) is fixed → marked FIXED
  ↓ AUTOMATIC STATUS PROPAGATION
BUG-001 (duplicate) → automatically marked FIXED
BUG-006 (duplicate) → automatically marked FIXED

Result: One fix, status propagates to all
        No wasted effort
        Consistent bug registry
```

## 🔄 AUTOMATIC STATUS PROPAGATION

### Trigger: Canonical Bug Fixed
```bash
# When canonical bug is marked as fixed, propagate status to all duplicates
mark_bug_fixed_with_propagation() {
    local canonical_bug="$1"
    local resolution_notes="$2"
    local fixed_by="$3"

    echo "✅ R524: Marking $canonical_bug as FIXED"

    # Get all duplicate bugs
    DUPLICATES=$(jq -r --arg canon "$canonical_bug" '
        .bug_registry[] |
        select(.bug_id == $canon) |
        .duplicates[]? // empty
    ' orchestrator-state-v3.json)

    # Mark canonical bug as fixed
    jq --arg bug "$canonical_bug" \
       --arg notes "$resolution_notes" \
       --arg fixed "$fixed_by" \
       '
       (.bug_registry[] | select(.bug_id == $bug)) |= {
           fix_status: "FIXED",
           resolution_notes: $notes,
           fixed_by: $fixed,
           updated_at: (now | todate),
           .
       }
       ' orchestrator-state-v3.json > /tmp/state.json

    mv /tmp/state.json orchestrator-state-v3.json

    echo "✅ Canonical bug $canonical_bug marked FIXED"

    # PROPAGATE: Mark all duplicates as fixed
    if [[ -n "$DUPLICATES" ]]; then
        echo "🔄 R524: Propagating fix status to duplicates..."

        for dup in $DUPLICATES; do
            propagate_fix_to_duplicate "$dup" "$canonical_bug" "$resolution_notes" "$fixed_by"
        done

        echo "✅ R524: Status propagated to all $(echo "$DUPLICATES" | wc -w) duplicates"
    else
        echo "ℹ️ No duplicates to propagate to"
    fi
}
```

### Propagate Fix Status to Duplicate Bug
```bash
# Mark a duplicate bug as fixed (because its canonical was fixed)
propagate_fix_to_duplicate() {
    local duplicate_bug="$1"
    local canonical_bug="$2"
    local canonical_resolution="$3"
    local canonical_fixed_by="$4"

    echo "  🔄 Propagating fix status to $duplicate_bug..."

    jq --arg dup "$duplicate_bug" \
       --arg canon "$canonical_bug" \
       --arg notes "$canonical_resolution" \
       --arg fixed "$canonical_fixed_by" \
       '
       (.bug_registry[] | select(.bug_id == $dup)) |= {
           fix_status: "FIXED_AS_DUPLICATE",
           resolution_notes: "Fixed as duplicate of \($canon): \($notes)",
           fixed_by: $fixed,
           fixed_via_propagation: true,
           propagation_source: $canon,
           updated_at: (now | todate),
           .
       }
       ' orchestrator-state-v3.json > /tmp/state.json

    mv /tmp/state.json orchestrator-state-v3.json

    echo "  ✅ $duplicate_bug marked FIXED_AS_DUPLICATE"
}
```

## 📋 FIX STATUS TYPES

### Status Values for Bug Registry
```
OPEN                 - Bug discovered, not yet fixed
IN_PROGRESS          - Fix effort underway
FIXED                - Canonical bug fixed directly
FIXED_AS_DUPLICATE   - Duplicate bug fixed via status propagation from canonical
DUPLICATE            - Bug marked as duplicate (not yet fixed)
```

### Querying Fixed Bugs
```bash
# Get all truly fixed bugs (canonical fixes)
get_canonical_fixes() {
    jq -r '.bug_registry[] | select(.fix_status == "FIXED") | .bug_id' orchestrator-state-v3.json
}

# Get all bugs fixed via status propagation
get_propagated_fixes() {
    jq -r '.bug_registry[] | select(.fix_status == "FIXED_AS_DUPLICATE") | .bug_id' orchestrator-state-v3.json
}

# Get ALL fixed bugs (both canonical and duplicates)
get_all_fixed_bugs() {
    jq -r '.bug_registry[] | select(.fix_status == "FIXED" or .fix_status == "FIXED_AS_DUPLICATE") | .bug_id' orchestrator-state-v3.json
}
```

## 🔍 REVERSE PROPAGATION (REOPENING)

### If Canonical Bug Reopened
```bash
# If canonical bug needs to be reopened, propagate status to duplicates
reopen_bug_with_propagation() {
    local canonical_bug="$1"
    local reopen_reason="$2"

    echo "⚠️ R524: Reopening $canonical_bug"

    # Get duplicates
    DUPLICATES=$(jq -r --arg canon "$canonical_bug" '
        .bug_registry[] |
        select(.bug_id == $canon) |
        .duplicates[]? // empty
    ' orchestrator-state-v3.json)

    # Reopen canonical
    jq --arg bug "$canonical_bug" \
       --arg reason "$reopen_reason" \
       '
       (.bug_registry[] | select(.bug_id == $bug)) |= {
           fix_status: "OPEN",
           resolution_notes: "Reopened: \($reason)",
           fixed_by: null,
           updated_at: (now | todate),
           .
       }
       ' orchestrator-state-v3.json > /tmp/state.json

    mv /tmp/state.json orchestrator-state-v3.json

    # PROPAGATE: Reopen all duplicates
    if [[ -n "$DUPLICATES" ]]; then
        echo "🔄 R524: Propagating reopen status to duplicates..."

        for dup in $DUPLICATES; do
            jq --arg dup "$dup" \
               --arg canon "$canonical_bug" \
               --arg reason "$reopen_reason" \
               '
               (.bug_registry[] | select(.bug_id == $dup)) |= {
                   fix_status: "DUPLICATE",
                   resolution_notes: "Reopened as duplicate of \($canon): \($reason)",
                   fixed_by: null,
                   fixed_via_propagation: null,
                   propagation_source: null,
                   updated_at: (now | todate),
                   .
               }
               ' orchestrator-state-v3.json > /tmp/state.json

            mv /tmp/state.json orchestrator-state-v3.json

            echo "  ✅ $dup reopened"
        done
    fi
}
```

## 🎯 AGENT-SPECIFIC RESPONSIBILITIES

### SW Engineer (After Fixing Bug)
```bash
# When SW engineer completes bug fix, check if it's canonical
sw_engineer_bug_fix_completion() {
    local bug_id="$1"
    local fix_commit="$2"

    echo "✅ SW Engineer: Bug $bug_id fixed in commit $fix_commit"

    # Check if this is a canonical bug with duplicates
    HAS_DUPLICATES=$(jq -r --arg bug "$bug_id" '
        .bug_registry[] |
        select(.bug_id == $bug) |
        .duplicates | length
    ' orchestrator-state-v3.json)

    if [[ "$HAS_DUPLICATES" -gt 0 ]]; then
        echo "🔄 R524: This is a canonical bug with $HAS_DUPLICATES duplicates"
        echo "   Will trigger status propagation when marked as FIXED"

        # Mark as fixed with status propagation
        mark_bug_fixed_with_propagation "$bug_id" "Fixed in commit $fix_commit" "sw-engineer"
    else
        # Check if it's a duplicate (shouldn't be working on it!)
        IS_DUPLICATE=$(jq -r --arg bug "$bug_id" '
            .bug_registry[] |
            select(.bug_id == $bug) |
            .is_duplicate // false
        ' orchestrator-state-v3.json)

        if [[ "$IS_DUPLICATE" == "true" ]]; then
            echo "❌ R525 VIOLATION: Fixing a duplicate bug directly!"
            echo "   Should work on canonical bug instead"
            return 1
        fi

        # Normal bug (no duplicates)
        mark_bug_fixed_with_propagation "$bug_id" "Fixed in commit $fix_commit" "sw-engineer"
    fi
}
```

### Integration Agent (After Successful Integration)
```bash
# After integration succeeds, propagate all bug fix statuses
integration_propagate_bug_fixes() {
    echo "🔄 R524: Integration successful - propagating bug fix statuses"

    # Get all bugs fixed in this integration
    FIXED_BUGS=$(jq -r --arg integration "$INTEGRATE_WAVE_EFFORTS_ID" '
        .bug_registry[] |
        select(.fix_status == "FIXED" and .detected_in_integration.name == $integration) |
        .bug_id
    ' orchestrator-state-v3.json)

    for bug in $FIXED_BUGS; do
        # Check for duplicates and propagate status
        HAS_DUPLICATES=$(jq -r --arg bug "$bug" '
            .bug_registry[] |
            select(.bug_id == $bug) |
            .duplicates | length
        ' orchestrator-state-v3.json)

        if [[ "$HAS_DUPLICATES" -gt 0 ]]; then
            echo "  🔄 Propagating fix status for $bug to its duplicates"
            # Re-trigger propagation (idempotent)
            RESOLUTION=$(jq -r --arg bug "$bug" '.bug_registry[] | select(.bug_id == $bug) | .resolution_notes' orchestrator-state-v3.json)
            FIXED_BY=$(jq -r --arg bug "$bug" '.bug_registry[] | select(.bug_id == $bug) | .fixed_by' orchestrator-state-v3.json)
            mark_bug_fixed_with_propagation "$bug" "$RESOLUTION" "$FIXED_BY"
        fi
    done
}
```

### Orchestrator (Monitoring)
```bash
# Verify status propagation consistency
orchestrator_propagation_validation() {
    echo "📊 R524: Validating status propagation consistency"

    # Find canonical bugs that are FIXED
    FIXED_CANONICALS=$(jq -r '
        .bug_registry[] |
        select(.fix_status == "FIXED" and .is_duplicate == false and (.duplicates | length > 0)) |
        .bug_id
    ' orchestrator-state-v3.json)

    for canon in $FIXED_CANONICALS; do
        # Get its duplicates
        DUPLICATES=$(jq -r --arg canon "$canon" '
            .bug_registry[] |
            select(.bug_id == $canon) |
            .duplicates[]
        ' orchestrator-state-v3.json)

        # Verify all duplicates are FIXED_AS_DUPLICATE
        for dup in $DUPLICATES; do
            DUP_STATUS=$(jq -r --arg dup "$dup" '
                .bug_registry[] |
                select(.bug_id == $dup) |
                .fix_status
            ' orchestrator-state-v3.json)

            if [[ "$DUP_STATUS" != "FIXED_AS_DUPLICATE" ]]; then
                echo "❌ R524 VIOLATION: Canonical $canon is FIXED but duplicate $dup is $DUP_STATUS"
                echo "   Triggering status propagation repair..."

                # Repair propagation
                RESOLUTION=$(jq -r --arg canon "$canon" '.bug_registry[] | select(.bug_id == $canon) | .resolution_notes' orchestrator-state-v3.json)
                FIXED_BY=$(jq -r --arg canon "$canon" '.bug_registry[] | select(.bug_id == $canon) | .fixed_by' orchestrator-state-v3.json)
                propagate_fix_to_duplicate "$dup" "$canon" "$RESOLUTION" "$FIXED_BY"
            fi
        done
    done

    echo "✅ R524: Status propagation validation complete"
}
```

## 🚨 ENFORCEMENT AND VALIDATION

### Status Propagation Consistency Check
```bash
# Ensure all canonical fixes propagated properly
validate_propagation_consistency() {
    echo "🔍 R524: Status propagation consistency validation"

    # Rule 1: If canonical is FIXED, all duplicates must be FIXED_AS_DUPLICATE
    VIOLATION_1=$(jq -r '
        .bug_registry[] |
        select(.fix_status == "FIXED" and .is_duplicate == false) as $canon |
        $canon.duplicates[]? as $dup_id |
        .bug_registry[] |
        select(.bug_id == $dup_id and .fix_status != "FIXED_AS_DUPLICATE") |
        "VIOLATION: \($canon.bug_id) is FIXED but \(.bug_id) is \(.fix_status)"
    ' orchestrator-state-v3.json)

    if [[ -n "$VIOLATION_1" ]]; then
        echo "❌ R524 PROPAGATION VIOLATION:"
        echo "$VIOLATION_1"
        return 1
    fi

    # Rule 2: If duplicate is FIXED_AS_DUPLICATE, canonical must be FIXED
    VIOLATION_2=$(jq -r '
        .bug_registry[] |
        select(.fix_status == "FIXED_AS_DUPLICATE" and .propagation_source != null) as $dup |
        .bug_registry[] |
        select(.bug_id == $dup.propagation_source and .fix_status != "FIXED") |
        "VIOLATION: \($dup.bug_id) is FIXED_AS_DUPLICATE but \(.bug_id) is \(.fix_status)"
    ' orchestrator-state-v3.json)

    if [[ -n "$VIOLATION_2" ]]; then
        echo "❌ R524 PROPAGATION VIOLATION:"
        echo "$VIOLATION_2"
        return 1
    fi

    echo "✅ R524: Status propagation consistency validated"
}
```

### Audit Trail Verification
```bash
# Verify status propagation audit trail is complete
verify_propagation_audit_trail() {
    echo "🔍 R524: Verifying status propagation audit trail"

    # Check all FIXED_AS_DUPLICATE bugs have propagation metadata
    MISSING_TRAIL=$(jq -r '
        .bug_registry[] |
        select(.fix_status == "FIXED_AS_DUPLICATE") |
        select(.fixed_via_propagation != true or .propagation_source == null) |
        "\(.bug_id): Missing propagation audit trail"
    ' orchestrator-state-v3.json)

    if [[ -n "$MISSING_TRAIL" ]]; then
        echo "❌ R524 VIOLATION: Incomplete propagation audit trail"
        echo "$MISSING_TRAIL"
        return 1
    fi

    echo "✅ R524: Complete propagation audit trail verified"
}
```

## Common Scenarios

### Scenario 1: Canonical Bug Fixed
```bash
# BUG-007 (canonical) has duplicates: BUG-001, BUG-006
# SW Engineer fixes BUG-007

sw_engineer_bug_fix_completion "BUG-007" "abc123"
# Output: ✅ SW Engineer: Bug BUG-007 fixed in commit abc123
#         🔄 R524: This is a canonical bug with 2 duplicates
#         ✅ Canonical bug BUG-007 marked FIXED
#         🔄 R524: Propagating fix status to duplicates...
#           ✅ BUG-001 marked FIXED_AS_DUPLICATE
#           ✅ BUG-006 marked FIXED_AS_DUPLICATE
#         ✅ R524: Status propagated to all 2 duplicates

# Result: All 3 bugs marked as fixed
#         BUG-007: fix_status=FIXED
#         BUG-001: fix_status=FIXED_AS_DUPLICATE, propagation_source=BUG-007
#         BUG-006: fix_status=FIXED_AS_DUPLICATE, propagation_source=BUG-007
```

### Scenario 2: Integration Completes
```bash
# Integration succeeds after fixing several bugs
integration_propagate_bug_fixes
# Scans all bugs fixed in this integration
# Propagates each canonical bug's fix status to its duplicates
# Ensures no duplicate is left OPEN
```

### Scenario 3: Bug Needs Reopening
```bash
# BUG-007 was marked FIXED but fix was incomplete
reopen_bug_with_propagation "BUG-007" "Fix incomplete - still failing tests"
# Output: ⚠️ R524: Reopening BUG-007
#         🔄 R524: Propagating reopen status to duplicates...
#           ✅ BUG-001 reopened
#           ✅ BUG-006 reopened

# All bugs back to fixable state
```

## Grading Impact

### COMPLIANCE BONUS (+30%)
- Automatic status propagation on all canonical fixes
- Complete propagation audit trail
- Consistent fix states across duplicates

### VIOLATIONS
- Manual duplicate fix (not via status propagation): **-75%**
- Canonical FIXED but duplicate not: **-100%** (INCONSISTENT STATE)
- Missing propagation audit trail: **-50%**
- Duplicate FIXED_AS_DUPLICATE but canonical OPEN: **-100%** (CORRUPTED STATE)

## Related Rules
- R522: Duplicate Bug Detection Protocol (prevention)
- R523: Duplicate Bug Linking Protocol (linking duplicates)
- R525: Duplicate Work Redirection Rule (work on canonical only)

## Remember

**"One fix, status propagates to all - no manual duplicate fixes"**
**"FIXED propagates to FIXED_AS_DUPLICATE automatically"**
**"Reopening canonical propagates to all duplicates"**
**"Maintain propagation audit trail for every duplicate fix"**

When a canonical bug is fixed, R524 ensures that ALL its duplicates are automatically marked as fixed through status propagation. This prevents agents from wasting effort "fixing" bugs that are already resolved, and maintains perfect consistency across the bug registry.

**Note**: This rule handles **bug status propagation** (updating duplicate statuses). For **integration cascade** (rebuilding dependent integrations after fixes), see R406, R348, R350, R351. For terminology clarification, see R530.

## Date Added
2025-10-06

## Changelog
- 2025-10-06: Renamed from "Bug Resolution Cascade Protocol" to "Bug Status Propagation Protocol" to avoid confusion with integration cascade (R406). Updated all "cascade" references to "propagation" where referring to status updates. Added cross-reference to R530 for terminology disambiguation.
- 2025-10-06: Initial creation to automate duplicate bug status propagation and prevent wasteful duplicate fixing efforts
