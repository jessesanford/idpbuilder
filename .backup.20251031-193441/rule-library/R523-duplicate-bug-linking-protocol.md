# 🚨🚨🚨 RULE R523: Duplicate Bug Linking Protocol

## Criticality: BLOCKING
**Violation Penalty: -50% per unlinked duplicate, -100% for transitive duplicates**

## Description
This rule defines how to handle duplicate bugs that WERE created (despite R522 prevention). It provides the protocol for marking bugs as duplicates, linking them to canonical bugs, and preventing transitive duplicate chains.

## The Problem This Solves

### Without R523 (CHAOS):
```
BUG-001: "Duplicate PushCmd" (created first)
BUG-006: "Duplicate PushCmd" (created accidentally)
BUG-007: "Duplicate PushCmd" (created again)

Result: 3 bugs, no linkage
        Agents work on all 3 separately
        Fixes waste effort
        No single source of truth
```

### With R523 (ORGANIZED):
```
BUG-007: is_duplicate=false, duplicates=["BUG-001", "BUG-006"]  (CANONICAL)
BUG-001: is_duplicate=true, duplicate_of="BUG-007"
BUG-006: is_duplicate=true, duplicate_of="BUG-007"

Result: Clear hierarchy
        All work redirected to BUG-007
        Single fix covers all
        Audit trail maintained
```

## 🔗 MARKING A BUG AS DUPLICATE

### Step 1: Identify Canonical Bug
```bash
# The canonical bug is the one that should be worked on
# Criteria: earliest created, most complete info, or most recent fix attempt
identify_canonical_bug() {
    local bug_a="$1"
    local bug_b="$2"

    # Get creation timestamps
    CREATED_A=$(jq -r --arg bug "$bug_a" '.bug_registry[] | select(.bug_id == $bug) | .created_at' orchestrator-state-v3.json)
    CREATED_B=$(jq -r --arg bug "$bug_b" '.bug_registry[] | select(.bug_id == $bug) | .created_at' orchestrator-state-v3.json)

    # Earliest created becomes canonical
    if [[ "$CREATED_A" < "$CREATED_B" ]]; then
        echo "$bug_a"  # A is canonical
    else
        echo "$bug_b"  # B is canonical
    fi
}
```

### Step 2: Validate No Transitive Duplicates
```bash
# CRITICAL: Prevent duplicate-of-duplicate chains
# BUG-A cannot be duplicate of BUG-B if BUG-B is itself a duplicate
validate_no_transitive_duplicates() {
    local duplicate_bug="$1"
    local canonical_bug="$2"

    echo "🔍 R523: Validating no transitive duplicates..."

    # Check if proposed canonical is itself a duplicate
    IS_DUPLICATE=$(jq -r --arg bug "$canonical_bug" '
        .bug_registry[] |
        select(.bug_id == $bug) |
        .is_duplicate // false
    ' orchestrator-state-v3.json)

    if [[ "$IS_DUPLICATE" == "true" ]]; then
        # Find the TRUE canonical bug
        TRUE_CANONICAL=$(jq -r --arg bug "$canonical_bug" '
            .bug_registry[] |
            select(.bug_id == $bug) |
            .duplicate_of
        ' orchestrator-state-v3.json)

        echo "⚠️ R523: $canonical_bug is itself a duplicate of $TRUE_CANONICAL"
        echo "   Using $TRUE_CANONICAL as canonical instead"
        echo "$TRUE_CANONICAL"
        return 0
    fi

    echo "✅ R523: $canonical_bug is canonical (not a duplicate)"
    echo "$canonical_bug"
}
```

### Step 3: Link Bugs Atomically
```bash
# Update BOTH bugs in a single atomic operation
link_duplicate_bugs() {
    local duplicate_bug="$1"
    local canonical_bug="$2"

    # Validate no transitive chain
    FINAL_CANONICAL=$(validate_no_transitive_duplicates "$duplicate_bug" "$canonical_bug")

    echo "🔗 R523: Linking $duplicate_bug as duplicate of $FINAL_CANONICAL"

    # Atomic update: modify both bugs
    jq --arg dup "$duplicate_bug" \
       --arg canon "$FINAL_CANONICAL" \
       '
       # Update duplicate bug: set is_duplicate=true, duplicate_of=canonical
       (.bug_registry[] | select(.bug_id == $dup)) |= {
           is_duplicate: true,
           duplicate_of: $canon,
           .
       } |
       # Update canonical bug: add to duplicates array
       (.bug_registry[] | select(.bug_id == $canon)) |= {
           duplicates: (.duplicates // [] | . + [$dup] | unique),
           .
       }
       ' orchestrator-state-v3.json > /tmp/state.json

    mv /tmp/state.json orchestrator-state-v3.json

    echo "✅ R523: Duplicate linked successfully"
    echo "   $duplicate_bug → $FINAL_CANONICAL"
}
```

## 🔍 POST-CREATION DUPLICATE DETECTION

### Audit All Bugs for Duplicates
```bash
# Run after integration or periodically to find missed duplicates
audit_for_duplicates() {
    echo "📊 R523: Post-creation duplicate audit"

    # Get all bugs
    ALL_BUGS=$(jq -r '.bug_registry[].bug_id' orchestrator-state-v3.json)

    # Compare each bug with others
    for bug_a in $ALL_BUGS; do
        # Skip if already marked as duplicate
        IS_DUP=$(jq -r --arg bug "$bug_a" '.bug_registry[] | select(.bug_id == $bug) | .is_duplicate // false' orchestrator-state-v3.json)
        [[ "$IS_DUP" == "true" ]] && continue

        # Get bug A's description
        DESC_A=$(jq -r --arg bug "$bug_a" '.bug_registry[] | select(.bug_id == $bug) | .description' orchestrator-state-v3.json)

        for bug_b in $ALL_BUGS; do
            [[ "$bug_a" == "$bug_b" ]] && continue

            # Skip if B already marked as duplicate
            IS_DUP_B=$(jq -r --arg bug "$bug_b" '.bug_registry[] | select(.bug_id == $bug) | .is_duplicate // false' orchestrator-state-v3.json)
            [[ "$IS_DUP_B" == "true" ]] && continue

            # Get bug B's description
            DESC_B=$(jq -r --arg bug "$bug_b" '.bug_registry[] | select(.bug_id == $bug) | .description' orchestrator-state-v3.json)

            # Check if descriptions match (exact or high similarity)
            if [[ "$DESC_A" == "$DESC_B" ]]; then
                echo "🔴 R523: Found duplicate bugs!"
                echo "   $bug_a: $DESC_A"
                echo "   $bug_b: $DESC_B"

                # Determine canonical and link
                CANONICAL=$(identify_canonical_bug "$bug_a" "$bug_b")
                if [[ "$CANONICAL" == "$bug_a" ]]; then
                    link_duplicate_bugs "$bug_b" "$bug_a"
                else
                    link_duplicate_bugs "$bug_a" "$bug_b"
                fi
            fi
        done
    done

    echo "✅ R523: Duplicate audit complete"
}
```

## 🔄 MANUAL DUPLICATE MARKING

### Command: Mark Bug as Duplicate
```bash
# For when humans/orchestrator identify duplicates manually
mark_as_duplicate() {
    local duplicate_bug="$1"
    local canonical_bug="$2"

    echo "📝 R523: Manual duplicate marking requested"
    echo "   Duplicate: $duplicate_bug"
    echo "   Canonical: $canonical_bug"

    # Validate both bugs exist
    BUG_A_EXISTS=$(jq -r --arg bug "$duplicate_bug" '.bug_registry[] | select(.bug_id == $bug) | .bug_id' orchestrator-state-v3.json)
    BUG_B_EXISTS=$(jq -r --arg bug "$canonical_bug" '.bug_registry[] | select(.bug_id == $bug) | .bug_id' orchestrator-state-v3.json)

    if [[ -z "$BUG_A_EXISTS" ]] || [[ -z "$BUG_B_EXISTS" ]]; then
        echo "❌ R523: One or both bugs don't exist!"
        return 1
    fi

    # Link them
    link_duplicate_bugs "$duplicate_bug" "$canonical_bug"

    echo "✅ R523: Manual duplicate marking complete"
}
```

## 📊 BUG REGISTRY STRUCTURE (WITH DUPLICATES)

### Example: Linked Duplicates
```json
{
  "bug_registry": [
    {
      "bug_id": "BUG-007",
      "description": "Duplicate PushCmd definition in pkg/cmd/push.go",
      "is_duplicate": false,
      "duplicate_of": null,
      "duplicates": ["BUG-001", "BUG-006"],
      "fix_status": "IN_PROGRESS",
      "fix_attempts": [
        {
          "effort_id": "E1.2.3-fix-bug-007",
          "status": "IN_PROGRESS"
        }
      ]
    },
    {
      "bug_id": "BUG-001",
      "description": "Duplicate PushCmd definition in pkg/cmd/push.go",
      "is_duplicate": true,
      "duplicate_of": "BUG-007",
      "duplicates": [],
      "fix_status": "DUPLICATE",
      "resolution_notes": "Marked as duplicate of BUG-007"
    },
    {
      "bug_id": "BUG-006",
      "description": "Duplicate PushCmd definition in pkg/cmd/push.go",
      "is_duplicate": true,
      "duplicate_of": "BUG-007",
      "duplicates": [],
      "fix_status": "DUPLICATE",
      "resolution_notes": "Marked as duplicate of BUG-007"
    }
  ]
}
```

## 🎯 AGENT-SPECIFIC RESPONSIBILITIES

### Integration Agent
```bash
# During integration, check for post-creation duplicates
integration_duplicate_check() {
    # After creating any new bugs, audit for duplicates
    echo "🔍 R523: Post-integration duplicate check"

    # Get bugs created in THIS integration
    NEW_BUGS=$(jq -r --arg integration "$INTEGRATE_WAVE_EFFORTS_ID" '
        .bug_registry[] |
        select(.detected_in_integration.name == $integration) |
        .bug_id
    ' orchestrator-state-v3.json)

    # Check each new bug against existing bugs
    for new_bug in $NEW_BUGS; do
        NEW_DESC=$(jq -r --arg bug "$new_bug" '.bug_registry[] | select(.bug_id == $bug) | .description' orchestrator-state-v3.json)

        # Search for duplicates (R522 may have been bypassed)
        EXISTING_MATCH=$(jq -r --arg desc "$NEW_DESC" --arg new "$new_bug" '
            .bug_registry[] |
            select(.description == $desc and .bug_id != $new and .is_duplicate == false) |
            .bug_id
        ' orchestrator-state-v3.json | head -1)

        if [[ -n "$EXISTING_MATCH" ]]; then
            echo "⚠️ R523: Found duplicate after creation!"
            echo "   New: $new_bug"
            echo "   Existing: $EXISTING_MATCH"

            # Link them (existing becomes canonical)
            link_duplicate_bugs "$new_bug" "$EXISTING_MATCH"
        fi
    done
}
```

### Orchestrator
```bash
# During MONITOR state, check for duplicates
orchestrator_duplicate_monitoring() {
    # Run duplicate audit periodically
    echo "📊 R523: Orchestrator duplicate monitoring"

    # Count duplicates
    DUPLICATE_COUNT=$(jq -r '.bug_registry[] | select(.is_duplicate == true) | .bug_id' orchestrator-state-v3.json | wc -l)

    if [[ $DUPLICATE_COUNT -gt 0 ]]; then
        echo "⚠️ Found $DUPLICATE_COUNT duplicate bugs"

        # List them
        jq -r '.bug_registry[] | select(.is_duplicate == true) | "\(.bug_id) → \(.duplicate_of)"' orchestrator-state-v3.json

        # Ensure all work is on canonical bugs (R525)
        validate_no_work_on_duplicates
    fi
}
```

## 🚨 ENFORCEMENT AND VALIDATION

### Duplicate Integrity Check
```bash
# Validate duplicate relationships are consistent
validate_duplicate_integrity() {
    echo "🔍 R523: Duplicate integrity validation"

    # Check 1: No transitive duplicates
    TRANSITIVE=$(jq -r '
        .bug_registry[] |
        select(.is_duplicate == true and .duplicate_of != null) as $dup |
        .bug_registry[] |
        select(.bug_id == $dup.duplicate_of and .is_duplicate == true) |
        "\($dup.bug_id) → \(.bug_id) → \(.duplicate_of)"
    ' orchestrator-state-v3.json)

    if [[ -n "$TRANSITIVE" ]]; then
        echo "❌ R523 VIOLATION: Transitive duplicate chain detected!"
        echo "$TRANSITIVE"
        return 1
    fi

    # Check 2: All duplicate_of bugs exist and are canonical
    INVALID_CANONICAL=$(jq -r '
        .bug_registry[] |
        select(.is_duplicate == true and .duplicate_of != null) as $dup |
        (.bug_registry[] | select(.bug_id == $dup.duplicate_of)) as $canon |
        select($canon == null or $canon.is_duplicate == true) |
        "\($dup.bug_id) → \($dup.duplicate_of) (INVALID CANONICAL)"
    ' orchestrator-state-v3.json)

    if [[ -n "$INVALID_CANONICAL" ]]; then
        echo "❌ R523 VIOLATION: Duplicate points to invalid canonical!"
        echo "$INVALID_CANONICAL"
        return 1
    fi

    # Check 3: Canonical bugs list all their duplicates
    MISSING_BACKREF=$(jq -r '
        .bug_registry[] |
        select(.is_duplicate == true and .duplicate_of != null) as $dup |
        .bug_registry[] |
        select(.bug_id == $dup.duplicate_of) as $canon |
        select($canon.duplicates == null or ($canon.duplicates | contains([$dup.bug_id]) | not)) |
        "\($dup.bug_id) not in \($canon.bug_id).duplicates[]"
    ' orchestrator-state-v3.json)

    if [[ -n "$MISSING_BACKREF" ]]; then
        echo "❌ R523 VIOLATION: Canonical missing duplicate in duplicates[]!"
        echo "$MISSING_BACKREF"
        return 1
    fi

    echo "✅ R523: Duplicate integrity validated"
}
```

## Common Scenarios

### Scenario 1: Two Bugs Created for Same Issue
```bash
# Integration creates BUG-001
# Later integration creates BUG-006 for same issue
# R522 was bypassed or didn't catch it

# Post-creation audit detects duplicate
audit_for_duplicates
# Output: Found duplicate bugs!
#         BUG-001: Duplicate PushCmd
#         BUG-006: Duplicate PushCmd

# Automatically link them (BUG-001 earlier, becomes canonical)
link_duplicate_bugs "BUG-006" "BUG-001"
# BUG-006 now marked as duplicate of BUG-001
```

### Scenario 2: Transitive Duplicate Prevented
```bash
# BUG-A is duplicate of BUG-B
# Later, someone tries to mark BUG-C as duplicate of BUG-A

mark_as_duplicate "BUG-C" "BUG-A"
# R523 validates and detects: BUG-A is itself a duplicate of BUG-B
# Automatically redirects: BUG-C marked as duplicate of BUG-B instead
# Prevents chain: BUG-C → BUG-A → BUG-B (WRONG)
# Creates: BUG-C → BUG-B, BUG-A → BUG-B (CORRECT)
```

### Scenario 3: Manual Duplicate Marking
```bash
# Human identifies duplicates that automated detection missed
mark_as_duplicate "BUG-012" "BUG-007"
# Validates both bugs exist
# Links BUG-012 → BUG-007
# Updates BUG-007.duplicates to include BUG-012
# Atomic operation ensures consistency
```

## Grading Impact

### COMPLIANCE BONUS (+25%)
- Correctly linking duplicate bugs
- Preventing transitive duplicates
- Maintaining duplicate integrity

### VIOLATIONS
- Transitive duplicate chain: **-100%** (CRITICAL FAILURE)
- Unlinked duplicate bugs: **-50% per unlinked**
- Invalid canonical reference: **-75%**
- Missing backref in duplicates[]: **-30%**

## Related Rules
- R522: Duplicate Bug Detection Protocol (prevention - check before creating)
- R524: Bug Status Propagation Protocol (fixing canonical propagates status to all duplicates)
- R525: Duplicate Work Redirection Rule (agents work on canonical only)

## Remember

**"Prevent duplicates with R522, link them with R523"**
**"No transitive duplicates - always point to canonical"**
**"Atomic linking - update both bugs or neither"**
**"Post-creation audit catches what prevention misses"**

Duplicate bugs happen despite prevention. R523 ensures they are properly linked, creating a clear hierarchy that prevents wasted work and maintains a single source of truth.

## Date Added
2025-10-06

## Changelog
- 2025-10-06: Initial creation to handle post-creation duplicate linking and prevent transitive duplicate chains
