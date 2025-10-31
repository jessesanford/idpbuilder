# 🚨🚨🚨 RULE R522: Duplicate Bug Detection Protocol

## Criticality: BLOCKING
**Violation Penalty: -25% per duplicate bug created**

## Description
This rule mandates that ALL agents MUST check for duplicate bugs before creating new bug entries. It provides the protocol for detecting duplicates, updating existing bugs with new occurrences, and tracking bugs across multiple integrations/efforts.

## The Problem This Solves

### Without R522 (WASTEFUL):
```
Integration 1: Creates BUG-007 (Duplicate PushCmd)
Integration 2: Creates BUG-012 (Duplicate PushCmd) ← DUPLICATE!
Integration 3: Creates BUG-019 (Duplicate PushCmd) ← DUPLICATE!

Result: 3 bug entries for same issue
        3 separate fix efforts
        Confused tracking
```

### With R522 (EFFICIENT):
```
Integration 1: Creates BUG-007 (Duplicate PushCmd)
Integration 2: Updates BUG-007 with occurrence #2
Integration 3: Updates BUG-007 with occurrence #3

Result: 1 bug entry with full history
        1 fix effort covers all occurrences
        Clear tracking
```

## 🔍 MANDATORY DUPLICATE CHECK

### Before Creating ANY Bug
```bash
# MANDATORY: Run BEFORE creating bug entry
check_for_duplicate_bug() {
    local bug_summary="$1"  # e.g., "Duplicate PushCmd definition"
    local bug_type="$2"      # e.g., "COMPILATION"

    echo "🔍 R522: Checking for duplicate bugs..."

    # Search in orchestrator-state-v3.json bug registry
    BUG_REGISTRY=$(jq -r '.bug_registry // {}' orchestrator-state-v3.json)

    # Search for matching bug by summary or error signature
    DUPLICATE=$(jq -r --arg summary "$bug_summary" '
        .bug_registry | to_entries[] |
        select(.value.summary == $summary or .value.error_signature == $summary) |
        .key
    ' orchestrator-state-v3.json)

    if [ -n "$DUPLICATE" ]; then
        echo "✅ R522: Found duplicate bug: $DUPLICATE"
        return 0  # Duplicate found
    else
        echo "ℹ️ R522: No duplicate found - this is a new bug"
        return 1  # New bug
    fi
}
```

### Update Existing Bug (If Duplicate)
```bash
# If duplicate found, UPDATE existing bug instead of creating new
update_duplicate_bug() {
    local bug_id="$1"
    local new_occurrence="$2"

    echo "📝 R522: Updating existing bug $bug_id with new occurrence"

    # Add new occurrence to affected_integrations array
    jq --arg bug_id "$bug_id" \
       --arg integration "$INTEGRATE_WAVE_EFFORTS_ID" \
       --arg timestamp "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
       '.bug_registry[$bug_id].affected_integrations += [{
           integration: $integration,
           discovered_at: $timestamp
       }]' orchestrator-state-v3.json > /tmp/state.json

    mv /tmp/state.json orchestrator-state-v3.json

    # Update occurrence count
    jq --arg bug_id "$bug_id" \
       '.bug_registry[$bug_id].occurrence_count += 1' \
       orchestrator-state-v3.json > /tmp/state.json

    mv /tmp/state.json orchestrator-state-v3.json

    echo "✅ Bug $bug_id updated with occurrence #$(jq -r --arg bug_id "$bug_id" '.bug_registry[$bug_id].occurrence_count' orchestrator-state-v3.json)"
}
```

### Create New Bug (If No Duplicate)
```bash
# Only called if NO duplicate found
create_new_bug() {
    local bug_summary="$1"
    local bug_details="$2"

    echo "🆕 R522: Creating new bug entry (no duplicate found)"

    # Generate unique bug ID
    BUG_COUNT=$(jq -r '.bug_registry | length' orchestrator-state-v3.json)
    BUG_ID="BUG-$(printf "%03d" $((BUG_COUNT + 1)))"

    # Create bug entry
    jq --arg bug_id "$BUG_ID" \
       --arg summary "$bug_summary" \
       --arg integration "$INTEGRATE_WAVE_EFFORTS_ID" \
       --arg timestamp "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
       --argjson details "$bug_details" \
       '.bug_registry[$bug_id] = {
           summary: $summary,
           error_signature: $summary,
           type: $details.type,
           severity: $details.severity,
           created_at: $timestamp,
           occurrence_count: 1,
           affected_integrations: [{
               integration: $integration,
               discovered_at: $timestamp
           }],
           status: "OPEN",
           fix_efforts: []
       }' orchestrator-state-v3.json > /tmp/state.json

    mv /tmp/state.json orchestrator-state-v3.json

    echo "✅ Created new bug: $BUG_ID"
    echo "$BUG_ID"
}
```

## 🔑 BUG SIGNATURE MATCHING

### Exact Match (Highest Confidence)
```bash
# Match by exact error message
match_by_error_message() {
    local error_msg="$1"

    jq -r --arg error "$error_msg" '
        .bug_registry | to_entries[] |
        select(.value.error_signature == $error) |
        .key
    ' orchestrator-state-v3.json
}
```

### Fuzzy Match (Medium Confidence)
```bash
# Match by similar summary (requires manual verification)
match_by_summary_similarity() {
    local bug_summary="$1"

    # Extract key words from summary
    KEY_WORDS=$(echo "$bug_summary" | tr ' ' '\n' | grep -E '^[A-Z][a-z]+' | head -3)

    # Search for bugs with similar key words
    for bug in $(jq -r '.bug_registry | keys[]' orchestrator-state-v3.json); do
        BUG_SUMMARY=$(jq -r --arg bug "$bug" '.bug_registry[$bug].summary' orchestrator-state-v3.json)

        # Check if key words appear in existing summary
        MATCH_COUNT=0
        for word in $KEY_WORDS; do
            if echo "$BUG_SUMMARY" | grep -qi "$word"; then
                ((MATCH_COUNT++))
            fi
        done

        # If 2+ key words match, flag as potential duplicate
        if [ $MATCH_COUNT -ge 2 ]; then
            echo "⚠️ Potential duplicate: $bug (similarity: $MATCH_COUNT/$(($(echo $KEY_WORDS | wc -w))))"
            echo "   Existing: $BUG_SUMMARY"
            echo "   New: $bug_summary"
            return 0
        fi
    done

    return 1
}
```

### File/Line Match (Location-Based)
```bash
# Match by file and approximate line number
match_by_location() {
    local file="$1"
    local line="$2"

    # Allow +/- 5 lines for line number matching
    LINE_MIN=$((line - 5))
    LINE_MAX=$((line + 5))

    jq -r --arg file "$file" \
          --arg line_min "$LINE_MIN" \
          --arg line_max "$LINE_MAX" '
        .bug_registry | to_entries[] |
        select(.value.location.file == $file and
               (.value.location.line >= ($line_min | tonumber)) and
               (.value.location.line <= ($line_max | tonumber))) |
        .key
    ' orchestrator-state-v3.json
}
```

## 📊 BUG REGISTRY STRUCTURE

### In orchestrator-state-v3.json
```json
{
  "bug_registry": {
    "BUG-007": {
      "summary": "Duplicate PushCmd definition",
      "error_signature": "redefinition of 'PushCmd' at pkg/cmd/push.go:45",
      "type": "COMPILATION",
      "severity": "CRITICAL",
      "created_at": "2025-10-06T10:30:00Z",
      "occurrence_count": 3,
      "affected_integrations": [
        {
          "integration": "phase1-wave1-integration",
          "discovered_at": "2025-10-06T10:30:00Z"
        },
        {
          "integration": "phase1-wave2-integration",
          "discovered_at": "2025-10-06T14:15:00Z"
        },
        {
          "integration": "phase1-wave3-integration",
          "discovered_at": "2025-10-06T16:45:00Z"
        }
      ],
      "status": "OPEN",
      "fix_efforts": [
        {
          "effort_id": "E1.2.3-fix-bug-007",
          "branch": "phase1-wave1-fix-bug-007",
          "status": "IN_PROGRESS"
        }
      ]
    }
  }
}
```

## 🔄 INTEGRATE_WAVE_EFFORTS WITH FIX CASCADE

### Linking Bugs to Fix Efforts
```bash
# When fix effort is created, link it to bug
link_bug_to_fix_effort() {
    local bug_id="$1"
    local effort_id="$2"
    local branch="$3"

    jq --arg bug_id "$bug_id" \
       --arg effort_id "$effort_id" \
       --arg branch "$branch" \
       '.bug_registry[$bug_id].fix_efforts += [{
           effort_id: $effort_id,
           branch: $branch,
           status: "IN_PROGRESS",
           started_at: (now | todate)
       }]' orchestrator-state-v3.json > /tmp/state.json

    mv /tmp/state.json orchestrator-state-v3.json

    echo "✅ Linked $effort_id to $bug_id"
}
```

### Closing Bug After Fix
```bash
# When fix is verified, close the bug
close_fixed_bug() {
    local bug_id="$1"

    jq --arg bug_id "$bug_id" \
       '.bug_registry[$bug_id].status = "FIXED" |
        .bug_registry[$bug_id].fixed_at = (now | todate)' \
       orchestrator-state-v3.json > /tmp/state.json

    mv /tmp/state.json orchestrator-state-v3.json

    echo "✅ Closed $bug_id as FIXED"
}
```

## 🚨 ENFORCEMENT AND VALIDATION

### Pre-Bug-Creation Check (MANDATORY)
```bash
# MUST be called before creating ANY bug
validate_no_duplicate_r522() {
    local bug_summary="$1"

    echo "🔍 R522: MANDATORY duplicate check"

    # Try all matching strategies
    DUPLICATE=$(match_by_error_message "$bug_summary" || \
                match_by_summary_similarity "$bug_summary" || \
                echo "")

    if [ -n "$DUPLICATE" ]; then
        echo "❌ R522 VIOLATION: Duplicate bug detected!"
        echo "   Existing bug: $DUPLICATE"
        echo "   Attempted summary: $bug_summary"
        echo ""
        echo "   MUST update existing bug, not create new one!"
        return 1
    fi

    echo "✅ R522: No duplicate found - safe to create new bug"
    return 0
}
```

### Post-Integration Audit
```bash
# Check if any duplicate bugs were created
audit_duplicate_bugs() {
    echo "📊 R522: Duplicate bug audit"

    # Extract all bug summaries
    jq -r '.bug_registry | to_entries[] | .value.summary' orchestrator-state-v3.json | sort > /tmp/bug_summaries.txt

    # Check for duplicate summaries
    DUPLICATES=$(uniq -d /tmp/bug_summaries.txt)

    if [ -n "$DUPLICATES" ]; then
        echo "❌ R522 VIOLATION: Duplicate bugs detected!"
        echo "$DUPLICATES"
        return 1
    fi

    echo "✅ No duplicate bugs found"
    return 0
}
```

## 📈 BUG OCCURRENCE TRACKING

### High-Occurrence Bugs (Priority)
```bash
# Identify bugs that occur frequently
identify_high_occurrence_bugs() {
    echo "📊 R522: High-occurrence bug analysis"

    # Bugs with 3+ occurrences are high priority
    HIGH_OCCURRENCE=$(jq -r '
        .bug_registry | to_entries[] |
        select(.value.occurrence_count >= 3) |
        "\(.key): \(.value.occurrence_count) occurrences - \(.value.summary)"
    ' orchestrator-state-v3.json)

    if [ -n "$HIGH_OCCURRENCE" ]; then
        echo "⚠️ HIGH-OCCURRENCE BUGS (Fix these FIRST!):"
        echo "$HIGH_OCCURRENCE"
    fi
}
```

### Cross-Integration Bug Patterns
```bash
# Detect bugs that appear across multiple integrations
detect_cross_integration_bugs() {
    echo "📊 R522: Cross-integration bug pattern analysis"

    # Bugs affecting 2+ integrations indicate systemic issues
    CROSS_INTEGRATE_WAVE_EFFORTS=$(jq -r '
        .bug_registry | to_entries[] |
        select(.value.affected_integrations | length >= 2) |
        "\(.key): affects \(.value.affected_integrations | length) integrations"
    ' orchestrator-state-v3.json)

    if [ -n "$CROSS_INTEGRATE_WAVE_EFFORTS" ]; then
        echo "🔴 CROSS-INTEGRATE_WAVE_EFFORTS BUGS (Systemic issues!):"
        echo "$CROSS_INTEGRATE_WAVE_EFFORTS"
    fi
}
```

## 🎯 AGENT-SPECIFIC RESPONSIBILITIES

### Integration Agent
```bash
# During integration, ALWAYS check for duplicates
integration_agent_bug_protocol() {
    local bug_summary="$1"

    # 1. MANDATORY: Check for duplicate
    if check_for_duplicate_bug "$bug_summary"; then
        # Duplicate found - update it
        DUPLICATE_ID=$(match_by_error_message "$bug_summary")
        update_duplicate_bug "$DUPLICATE_ID" "$INTEGRATE_WAVE_EFFORTS_ID"
        echo "Updated existing bug: $DUPLICATE_ID"
    else
        # No duplicate - create new (but only after R521 known fix check!)
        # First check R521: Is this a known fix?
        if search_known_fixes "$bug_summary"; then
            apply_known_fix "$bug_summary"  # R521 protocol
        else
            # Truly new - create bug entry
            BUG_ID=$(create_new_bug "$bug_summary" "$bug_details")
            echo "Created new bug: $BUG_ID"
        fi
    fi
}
```

### Orchestrator
```bash
# Monitor bug registry for patterns
orchestrator_bug_monitoring() {
    # Check for high-occurrence bugs
    identify_high_occurrence_bugs

    # Check for cross-integration bugs
    detect_cross_integration_bugs

    # Prioritize fixes based on occurrence count
    prioritize_bug_fixes
}
```

## Common Scenarios

### Scenario 1: Same Bug in Multiple Integrations
```bash
# Integration 1
check_for_duplicate_bug "Duplicate PushCmd"  # Not found
create_new_bug "Duplicate PushCmd"  # Creates BUG-007

# Integration 2 (CASCADE RETRY)
check_for_duplicate_bug "Duplicate PushCmd"  # Found: BUG-007
update_duplicate_bug "BUG-007" "phase1-wave2-integration"  # Update, don't create

# Integration 3
check_for_duplicate_bug "Duplicate PushCmd"  # Found: BUG-007
update_duplicate_bug "BUG-007" "phase1-wave3-integration"  # Update again
```

### Scenario 2: Similar But Different Bugs
```bash
# Bug 1: "Undefined variable 'client'" in file A
create_new_bug "Undefined variable 'client' in registry.go"  # BUG-008

# Bug 2: "Undefined variable 'client'" in file B (different location!)
check_for_duplicate_bug "Undefined variable 'client' in gitea.go"  # NOT EXACT MATCH
match_by_summary_similarity "..."  # Potential duplicate flagged
# Manual review determines: Different root cause, create new bug
create_new_bug "Undefined variable 'client' in gitea.go"  # BUG-009
```

## Grading Impact

### COMPLIANCE BONUS (+15%)
- Correctly detecting and updating duplicate bugs
- Maintaining clean bug registry
- Accurate occurrence tracking

### VIOLATIONS
- Creating duplicate bug: **-25% per duplicate**
- Not checking for duplicates: **-50%**
- Corrupted bug registry: **-100%**

## Related Rules
- R521: Integration Known Fixes Protocol (check before creating bug)
- R266: Upstream Bug Documentation (how to document)
- R300: Comprehensive Fix Management (fix cascade protocol)

## Remember

**"One bug, one entry - update, don't duplicate"**
**"Check R522 before creating ANY bug"**
**"High occurrence = high priority"**
**"Bug registry is the single source of truth"**

Duplicate bugs waste time, create confusion, and lead to redundant fix efforts. The bug registry should be a clean, deduplicated source of truth for all issues discovered during integrations.

## Date Added
2025-10-06

## Changelog
- 2025-10-06: Initial creation based on orchestrator analysis identifying duplicate bug creation pattern
