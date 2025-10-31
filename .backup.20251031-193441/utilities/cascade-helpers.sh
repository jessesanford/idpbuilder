#!/bin/bash
# cascade-helpers.sh - Helper functions for R406 Fix Cascade Tracking
#
# Purpose: Provides easy-to-use functions for managing fix cascade state
# Usage: source utilities/cascade-helpers.sh
#
# Part of Software Factory 2.0 - Rule R406

set -euo pipefail

# Default state file location
STATE_FILE="${STATE_FILE:-${CLAUDE_PROJECT_DIR}/orchestrator-state-v3.json}"

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

#======================================================================
# CASCADE INITIALIZATION
#======================================================================

# Initialize a new fix cascade
# Usage: cascade_init "phase1_wave2" "wave"
cascade_init() {
    local integration_name="$1"
    local integration_type="$2"
    local state_file="${3:-$STATE_FILE}"

    echo -e "${BLUE}🔄 Initializing new fix cascade...${NC}"

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
        } |
        .bug_registry = (.bug_registry // []) |
        .integration_fix_states = (.integration_fix_states // {}) |
        .effort_fix_states = (.effort_fix_states // {})' \
        "$state_file" > "${state_file}.tmp"
    mv "${state_file}.tmp" "$state_file"

    echo -e "${GREEN}✅ Cascade initialized: $cascade_id${NC}"
    echo "$cascade_id"
}

# Get current cascade ID
cascade_get_id() {
    local state_file="${1:-$STATE_FILE}"
    jq -r '.fix_cascade_state.cascade_id // "none"' "$state_file"
}

# Check if cascade is active
cascade_is_active() {
    local state_file="${1:-$STATE_FILE}"
    local active=$(jq -r '.fix_cascade_state.active // false' "$state_file")
    [[ "$active" == "true" ]]
}

#======================================================================
# BUG REGISTRATION
#======================================================================

# Register a new bug in the cascade
# Usage: cascade_register_bug CASCADE_ID INTEGRATE_WAVE_EFFORTS SEVERITY CATEGORY DESCRIPTION EFFORT1 [EFFORT2 ...]
cascade_register_bug() {
    local cascade_id="$1"
    local integration="$2"
    local severity="$3"
    local category="$4"
    local description="$5"
    shift 5
    local efforts=("$@")
    local state_file="${STATE_FILE}"

    echo -e "${BLUE}📝 Registering new bug...${NC}"

    # Generate bug ID
    local bug_count=$(jq -r '.bug_registry | length' "$state_file")
    local bug_id=$(printf "BUG-%s-%03d" "$cascade_id" $((bug_count + 1)))

    # Convert efforts array to JSON
    local efforts_json=$(printf '%s\n' "${efforts[@]}" | jq -R . | jq -s .)

    # Add to bug registry
    jq --arg bid "$bug_id" \
       --arg cid "$cascade_id" \
       --arg integration "$integration" \
       --arg severity "$severity" \
       --arg category "$category" \
       --arg description "$description" \
       --argjson efforts "$efforts_json" \
       '.bug_registry += [{
          bug_id: $bid,
          cascade_id: $cid,
          cascade_layer: 1,
          detected_in_integration: {
            name: $integration,
            type: "wave"
          },
          detected_at: (now | todate),
          detected_by: "Manual",
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
        }]' "$state_file" > "${state_file}.tmp"
    mv "${state_file}.tmp" "$state_file"

    echo -e "${GREEN}✅ Bug registered: $bug_id${NC}"
    echo "$bug_id"
}

#======================================================================
# BUG STATUS MANAGEMENT
#======================================================================

# Update bug status
# Usage: cascade_update_bug_status BUG_ID NEW_STATUS
cascade_update_bug_status() {
    local bug_id="$1"
    local new_status="$2"
    local state_file="${3:-$STATE_FILE}"

    echo -e "${BLUE}📝 Updating bug $bug_id status to: $new_status${NC}"

    jq --arg bid "$bug_id" \
       --arg status "$new_status" \
       '(.bug_registry[] | select(.bug_id == $bid) | .fix_status) = $status |
        (.bug_registry[] | select(.bug_id == $bid) | .updated_at) = (now | todate)' \
        "$state_file" > "${state_file}.tmp"
    mv "${state_file}.tmp" "$state_file"

    echo -e "${GREEN}✅ Bug status updated${NC}"
}

# Start a fix attempt
# Usage: cascade_start_fix BUG_ID
cascade_start_fix() {
    local bug_id="$1"
    local state_file="${2:-$STATE_FILE}"

    echo -e "${BLUE}🔧 Starting fix for bug: $bug_id${NC}"

    local current_attempt=$(jq -r ".bug_registry[] | select(.bug_id == \"$bug_id\") | .current_attempt" "$state_file")
    local next_attempt=$((current_attempt + 1))

    jq --arg bid "$bug_id" \
       --argjson attempt "$next_attempt" \
       '(.bug_registry[] | select(.bug_id == $bid) | .fix_status) = "in_progress" |
        (.bug_registry[] | select(.bug_id == $bid) | .current_attempt) = $attempt |
        (.bug_registry[] | select(.bug_id == $bid) | .fix_attempts) += [{
          attempt_number: $attempt,
          started_at: (now | todate),
          completed_at: null,
          commits: [],
          outcome: "in_progress",
          verified: false,
          notes: ""
        }] |
        (.bug_registry[] | select(.bug_id == $bid) | .updated_at) = (now | todate)' \
        "$state_file" > "${state_file}.tmp"
    mv "${state_file}.tmp" "$state_file"

    echo -e "${GREEN}✅ Fix attempt $next_attempt started${NC}"
}

# Complete a fix attempt
# Usage: cascade_complete_fix BUG_ID COMMIT_SHA PROJECT_DONE NOTES
cascade_complete_fix() {
    local bug_id="$1"
    local commit_sha="$2"
    local success="${3:-true}"
    local notes="${4:-Fix completed}"
    local state_file="${5:-$STATE_FILE}"

    echo -e "${BLUE}🔧 Completing fix for bug: $bug_id${NC}"

    local outcome="successful"
    local new_status="fixed"
    [[ "$success" != "true" ]] && outcome="failed" && new_status="pending"

    jq --arg bid "$bug_id" \
       --arg commit "$commit_sha" \
       --arg outcome "$outcome" \
       --arg status "$new_status" \
       --arg notes "$notes" \
       '(.bug_registry[] | select(.bug_id == $bid) | .fix_status) = $status |
        (.bug_registry[] | select(.bug_id == $bid) | .fix_attempts[-1].completed_at) = (now | todate) |
        (.bug_registry[] | select(.bug_id == $bid) | .fix_attempts[-1].commits) += [$commit] |
        (.bug_registry[] | select(.bug_id == $bid) | .fix_attempts[-1].outcome) = $outcome |
        (.bug_registry[] | select(.bug_id == $bid) | .fix_attempts[-1].notes) = $notes |
        (.bug_registry[] | select(.bug_id == $bid) | .updated_at) = (now | todate)' \
        "$state_file" > "${state_file}.tmp"
    mv "${state_file}.tmp" "$state_file"

    if [[ "$success" == "true" ]]; then
        echo -e "${GREEN}✅ Fix completed successfully${NC}"
    else
        echo -e "${YELLOW}⚠️ Fix attempt failed${NC}"
    fi
}

#======================================================================
# QUERY FUNCTIONS
#======================================================================

# Get all pending bugs
cascade_get_pending_bugs() {
    local state_file="${1:-$STATE_FILE}"
    jq -r '.bug_registry[] | select(.fix_status == "pending") | .bug_id' "$state_file"
}

# Get bugs for specific integration
cascade_get_integration_bugs() {
    local integration="$1"
    local state_file="${2:-$STATE_FILE}"
    jq -r ".bug_registry[] | select(.detected_in_integration.name == \"$integration\") | .bug_id" "$state_file"
}

# Get bugs assigned to effort
cascade_get_effort_bugs() {
    local effort="$1"
    local state_file="${2:-$STATE_FILE}"
    jq -r ".bug_registry[] | select(.affected_efforts[] == \"$effort\") | .bug_id" "$state_file"
}

# Show cascade status
cascade_status() {
    local state_file="${1:-$STATE_FILE}"

    echo -e "${BLUE}📊 CASCADE STATUS${NC}"
    echo "════════════════════════════════════════"

    local active=$(jq -r '.fix_cascade_state.active // false' "$state_file")

    if [[ "$active" != "true" ]]; then
        echo -e "${YELLOW}No active cascade${NC}"
        return 0
    fi

    local cascade_id=$(jq -r '.fix_cascade_state.cascade_id' "$state_file")
    local status=$(jq -r '.fix_cascade_state.status' "$state_file")
    local current_layer=$(jq -r '.fix_cascade_state.current_layer // 0' "$state_file")
    local total_layers=$(jq -r '.fix_cascade_state.total_layers // 0' "$state_file")

    echo "Cascade ID: $cascade_id"
    echo "Status: $status"
    echo "Layer: $current_layer / $total_layers"
    echo ""

    # Bug counts
    local total=$(jq -r '.bug_registry | length' "$state_file")
    local pending=$(jq -r '[.bug_registry[] | select(.fix_status == "pending")] | length' "$state_file")
    local in_progress=$(jq -r '[.bug_registry[] | select(.fix_status == "in_progress")] | length' "$state_file")
    local fixed=$(jq -r '[.bug_registry[] | select(.fix_status == "fixed")] | length' "$state_file")
    local verified=$(jq -r '[.bug_registry[] | select(.fix_status == "verified")] | length' "$state_file")

    echo "Bugs:"
    echo "  Total: $total"
    echo "  Pending: $pending"
    echo "  In Progress: $in_progress"
    echo "  Fixed: $fixed"
    echo "  Verified: $verified"
    echo ""

    # Cascade chain
    echo "Cascade Chain:"
    jq -r '.fix_cascade_state.cascade_chain[]? | "  Layer \(.layer): \(.integration_name) [\(.type)] - \(.status)"' "$state_file"
}

# List all bugs with details
cascade_list_bugs() {
    local state_file="${1:-$STATE_FILE}"

    echo -e "${BLUE}🐛 BUG REGISTRY${NC}"
    echo "════════════════════════════════════════"

    jq -r '.bug_registry[] |
           "\(.bug_id) [\(.severity)] - \(.fix_status)\n" +
           "  Desc: \(.description)\n" +
           "  Efforts: \(.affected_efforts | join(", "))\n" +
           "  Detected: \(.detected_at)\n"' \
        "$state_file"
}

#======================================================================
# VALIDATION
#======================================================================

# Validate cascade state
cascade_validate() {
    local state_file="${1:-$STATE_FILE}"

    echo -e "${BLUE}🔍 Validating cascade state...${NC}"

    # Count bugs in registry
    local registry_count=$(jq -r '.bug_registry | length' "$state_file")

    # Count bugs in validation
    local validation_count=$(jq -r '.fix_cascade_state.validation.total_bugs_detected // 0' "$state_file")

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

    local errors=0

    if [[ "$registry_count" -ne "$total_by_status" ]]; then
        echo -e "${RED}❌ Registry count ($registry_count) != Status sum ($total_by_status)${NC}"
        ((errors++))
    fi

    if [[ "$validation_count" -ne 0 ]] && [[ "$registry_count" -ne "$validation_count" ]]; then
        echo -e "${RED}❌ Registry count ($registry_count) != Validation count ($validation_count)${NC}"
        ((errors++))
    fi

    if [[ $errors -eq 0 ]]; then
        echo -e "${GREEN}✅ Validation PASSED: All bugs accounted for${NC}"
        return 0
    else
        echo -e "${RED}❌ Validation FAILED: $errors error(s) detected${NC}"
        return 1
    fi
}

# Generate checksum of bug IDs
cascade_generate_checksum() {
    local state_file="${1:-$STATE_FILE}"

    # Get sorted bug IDs
    local bug_ids=$(jq -r '.bug_registry[].bug_id' "$state_file" | sort | tr '\n' ',')

    # Generate MD5
    echo -n "$bug_ids" | md5sum | cut -d' ' -f1
}

# Update validation metadata
cascade_update_validation() {
    local state_file="${1:-$STATE_FILE}"

    echo -e "${BLUE}📝 Updating validation metadata...${NC}"

    local total=$(jq -r '.bug_registry | length' "$state_file")
    local fixed=$(jq -r '[.bug_registry[] | select(.fix_status == "fixed" or .fix_status == "verified" or .fix_status == "integrated")] | length' "$state_file")
    local pending=$(jq -r '[.bug_registry[] | select(.fix_status == "pending" or .fix_status == "in_progress")] | length' "$state_file")
    local checksum=$(cascade_generate_checksum "$state_file")

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
        }' "$state_file" > "${state_file}.tmp"
    mv "${state_file}.tmp" "$state_file"

    echo -e "${GREEN}✅ Validation metadata updated${NC}"
}

#======================================================================
# CASCADE COMPLETION
#======================================================================

# Check if cascade is complete
cascade_is_complete() {
    local state_file="${1:-$STATE_FILE}"

    local pending=$(jq -r '[.bug_registry[] | select(.fix_status == "pending" or .fix_status == "in_progress")] | length' "$state_file")

    [[ "$pending" -eq 0 ]]
}

# Mark cascade as complete
cascade_mark_complete() {
    local state_file="${1:-$STATE_FILE}"

    echo -e "${BLUE}🎉 Marking cascade as complete...${NC}"

    jq '.fix_cascade_state.active = false |
        .fix_cascade_state.status = "complete" |
        .fix_cascade_state.updated_at = (now | todate)' \
        "$state_file" > "${state_file}.tmp"
    mv "${state_file}.tmp" "$state_file"

    echo -e "${GREEN}✅ Cascade marked complete${NC}"
}

#======================================================================
# MIGRATION FROM LEGACY
#======================================================================

# Migrate from legacy fix tracking to R406 format
# Usage: cascade_migrate_from_legacy
cascade_migrate_from_legacy() {
    local state_file="${1:-$STATE_FILE}"

    echo -e "${BLUE}🔄 Migrating from legacy fix tracking format...${NC}"

    # Check if migration script exists
    local migration_script="${CLAUDE_PROJECT_DIR}/utilities/migrate-fix-tracking.sh"

    if [[ ! -f "$migration_script" ]]; then
        echo -e "${RED}❌ Migration script not found: $migration_script${NC}"
        return 1
    fi

    # Make sure it's executable
    chmod +x "$migration_script"

    # Run migration script
    echo -e "${BLUE}Running migration script...${NC}"
    bash "$migration_script" "$state_file"

    local result=$?

    if [[ $result -eq 0 ]]; then
        echo -e "${GREEN}✅ Migration completed successfully${NC}"

        # Show migration summary
        echo
        echo -e "${BLUE}Migration Summary:${NC}"
        cascade_status "$state_file"

        return 0
    else
        echo -e "${RED}❌ Migration failed - see migration log for details${NC}"
        return 1
    fi
}

#======================================================================
# MAIN EXECUTION (if called directly)
#======================================================================

# If script is executed directly (not sourced), show help
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    cat <<EOF
${BLUE}cascade-helpers.sh - R406 Fix Cascade Tracking Utilities${NC}

This script provides helper functions for managing fix cascades.

Usage: source utilities/cascade-helpers.sh

Available functions:
  ${GREEN}Initialization:${NC}
    cascade_init INTEGRATE_WAVE_EFFORTS_NAME TYPE      - Initialize new cascade
    cascade_get_id                           - Get current cascade ID
    cascade_is_active                        - Check if cascade active

  ${GREEN}Bug Management:${NC}
    cascade_register_bug ...                 - Register new bug
    cascade_update_bug_status BUG_ID STATUS - Update bug status
    cascade_start_fix BUG_ID                 - Start fix attempt
    cascade_complete_fix BUG_ID COMMIT ...   - Complete fix attempt

  ${GREEN}Queries:${NC}
    cascade_get_pending_bugs                 - List pending bugs
    cascade_get_integration_bugs INTEGRATE_WAVE_EFFORTS - Bugs for integration
    cascade_get_effort_bugs EFFORT           - Bugs for effort
    cascade_status                           - Show cascade status
    cascade_list_bugs                        - List all bugs

  ${GREEN}Validation:${NC}
    cascade_validate                         - Validate cascade state
    cascade_update_validation                - Update validation metadata
    cascade_is_complete                      - Check if cascade done
    cascade_mark_complete                    - Mark cascade complete

  ${GREEN}Migration:${NC}
    cascade_migrate_from_legacy              - Migrate from legacy fix tracking to R406

Example:
  $ source utilities/cascade-helpers.sh
  $ cascade_init "phase1_wave2" "wave"
  $ cascade_status
  $ cascade_validate

EOF
fi
