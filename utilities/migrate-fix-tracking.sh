#!/bin/bash
# migrate-fix-tracking.sh - Migrate from legacy fix tracking to R406 format
#
# Purpose: Automate migration of old fix tracking fields to R406 Fix Cascade Tracking
# Usage: bash utilities/migrate-fix-tracking.sh [state-file]
#
# Part of Software Factory 2.0 - R406 Migration Tool

set -euo pipefail

# Configuration
STATE_FILE="${1:-${CLAUDE_PROJECT_DIR}/orchestrator-state-v3.json}"
BACKUP_DIR="${CLAUDE_PROJECT_DIR}/backups"
MIGRATION_LOG="${CLAUDE_PROJECT_DIR}/migration-r406-$(date +%Y%m%d-%H%M%S).log"

# Source cascade helpers
source "${CLAUDE_PROJECT_DIR}/utilities/cascade-helpers.sh"

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Logging
log() {
    echo -e "${BLUE}[$(date '+%Y-%m-%d %H:%M:%S')]${NC} $*" | tee -a "$MIGRATION_LOG"
}

log_success() {
    echo -e "${GREEN}[$(date '+%Y-%m-%d %H:%M:%S')] ✅ $*${NC}" | tee -a "$MIGRATION_LOG"
}

log_error() {
    echo -e "${RED}[$(date '+%Y-%m-%d %H:%M:%S')] ❌ $*${NC}" | tee -a "$MIGRATION_LOG"
}

log_warning() {
    echo -e "${YELLOW}[$(date '+%Y-%m-%d %H:%M:%S')] ⚠️ $*${NC}" | tee -a "$MIGRATION_LOG"
}

#======================================================================
# PRE-MIGRATION CHECKS
#======================================================================

pre_migration_checks() {
    log "Running pre-migration checks..."

    # Check state file exists
    if [[ ! -f "$STATE_FILE" ]]; then
        log_error "State file not found: $STATE_FILE"
        exit 1
    fi
    log_success "State file exists"

    # Check if state file is valid JSON
    if ! jq empty "$STATE_FILE" 2>/dev/null; then
        log_error "State file is not valid JSON"
        exit 1
    fi
    log_success "State file is valid JSON"

    # Check current state
    CURRENT_STATE=$(jq -r '.current_state' "$STATE_FILE")
    log "Current state: $CURRENT_STATE"

    # Warn if not in expected state
    if [[ "$CURRENT_STATE" != "ERROR_RECOVERY" ]] && [[ "$CURRENT_STATE" != "CASCADE_REINTEGRATION" ]]; then
        log_warning "Current state is $CURRENT_STATE (expected ERROR_RECOVERY or CASCADE_REINTEGRATION)"
        log_warning "Migration can proceed but verify this is correct"
    fi

    # Check if already migrated
    if jq -e '.bug_registry' "$STATE_FILE" >/dev/null 2>&1; then
        REGISTRY_COUNT=$(jq -r '.bug_registry | length' "$STATE_FILE")
        if [[ "$REGISTRY_COUNT" -gt 0 ]]; then
            log_warning "bug_registry already exists with $REGISTRY_COUNT bugs"
            log_warning "This may be a re-migration. Proceeding will merge old and new bugs."
            read -p "Continue? (y/n) " -n 1 -r
            echo
            if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                log "Migration cancelled by user"
                exit 0
            fi
        fi
    fi

    log_success "Pre-migration checks complete"
}

#======================================================================
# BACKUP
#======================================================================

create_backup() {
    log "Creating backup..."

    mkdir -p "$BACKUP_DIR"
    BACKUP_FILE="${BACKUP_DIR}/orchestrator-state.pre-r406-$(date +%Y%m%d-%H%M%S).json"

    cp "$STATE_FILE" "$BACKUP_FILE"
    log_success "Backup created: $BACKUP_FILE"

    # Verify backup
    if ! jq empty "$BACKUP_FILE" 2>/dev/null; then
        log_error "Backup file is corrupted!"
        exit 1
    fi
    log_success "Backup verified"

    echo "$BACKUP_FILE"
}

#======================================================================
# DATA EXTRACTION
#======================================================================

extract_old_bugs() {
    log "Extracting bugs from old format..."

    # Try to extract from various old field formats
    local bugs_found=0

    # Format 1: upstream_bugs_wave2 (array of strings)
    if jq -e '.upstream_bugs_wave2' "$STATE_FILE" >/dev/null 2>&1; then
        local count=$(jq -r '.upstream_bugs_wave2 | length' "$STATE_FILE")
        log "Found $count bugs in upstream_bugs_wave2"
        bugs_found=$((bugs_found + count))
    fi

    # Format 2: project_fixes_in_progress (array of objects)
    if jq -e '.project_fixes_in_progress' "$STATE_FILE" >/dev/null 2>&1; then
        local count=$(jq -r '.project_fixes_in_progress | length' "$STATE_FILE")
        log "Found $count bugs in project_fixes_in_progress"
        bugs_found=$((bugs_found + count))
    fi

    # Format 3: efforts_needing_fixes (array of effort names)
    if jq -e '.efforts_needing_fixes' "$STATE_FILE" >/dev/null 2>&1; then
        local count=$(jq -r '.efforts_needing_fixes | length' "$STATE_FILE")
        log "Found $count efforts needing fixes in efforts_needing_fixes"
    fi

    # Format 4: fix_cascade.bugs (ad-hoc structure)
    if jq -e '.fix_cascade.bugs' "$STATE_FILE" >/dev/null 2>&1; then
        local count=$(jq -r '.fix_cascade.bugs | length' "$STATE_FILE")
        log "Found $count bugs in fix_cascade.bugs"
        bugs_found=$((bugs_found + count))
    fi

    if [[ $bugs_found -eq 0 ]]; then
        log_warning "No bugs found in old format - nothing to migrate"
        log_warning "If you expected bugs, check the old field names"
        return 1
    fi

    log_success "Found $bugs_found total bugs to migrate"
    echo "$bugs_found"
}

#======================================================================
# MIGRATION LOGIC
#======================================================================

migrate_bugs() {
    local integration_name="$1"
    local integration_type="$2"
    local cascade_id="$3"

    log "Migrating bugs to bug_registry..."

    # Initialize counters
    local migrated_count=0
    local failed_count=0

    # Migrate from upstream_bugs_wave2
    if jq -e '.upstream_bugs_wave2' "$STATE_FILE" >/dev/null 2>&1; then
        log "Migrating from upstream_bugs_wave2..."

        while IFS= read -r bug_description; do
            # Skip empty or null
            if [[ -z "$bug_description" ]] || [[ "$bug_description" == "null" ]]; then
                continue
            fi

            # Parse bug description to extract info
            # Format: "Description" or "Category: Description" or "Severity - Description"
            local severity="HIGH"
            local category="other"
            local description="$bug_description"

            # Try to detect severity
            if echo "$bug_description" | grep -iq "critical\|crash\|build fail"; then
                severity="CRITICAL"
            elif echo "$bug_description" | grep -iq "minor\|low"; then
                severity="LOW"
            fi

            # Try to detect category
            if echo "$bug_description" | grep -iq "build"; then
                category="build_failure"
            elif echo "$bug_description" | grep -iq "test"; then
                category="test_failure"
            elif echo "$bug_description" | grep -iq "lint"; then
                category="lint_error"
            elif echo "$bug_description" | grep -iq "runtime\|crash\|panic"; then
                category="runtime_error"
            elif echo "$bug_description" | grep -iq "conflict\|merge"; then
                category="integration_conflict"
            fi

            # Extract effort from efforts_needing_fixes if available
            local effort=""
            if jq -e '.efforts_needing_fixes[0]' "$STATE_FILE" >/dev/null 2>&1; then
                effort=$(jq -r '.efforts_needing_fixes[0]' "$STATE_FILE")
            else
                # Default effort name
                effort="UNKNOWN-EFFORT"
            fi

            # Register bug
            log "Registering bug: $description (severity: $severity, category: $category)"

            local bug_id=$(cascade_register_bug \
                "$cascade_id" \
                "$integration_name" \
                "$severity" \
                "$category" \
                "$description" \
                "$effort")

            if [[ -n "$bug_id" ]]; then
                log_success "Registered bug: $bug_id"
                ((migrated_count++))
            else
                log_error "Failed to register bug: $description"
                ((failed_count++))
            fi

        done < <(jq -r '.upstream_bugs_wave2[]?' "$STATE_FILE")
    fi

    # Migrate from project_fixes_in_progress (structured format)
    if jq -e '.project_fixes_in_progress' "$STATE_FILE" >/dev/null 2>&1; then
        log "Migrating from project_fixes_in_progress..."

        while IFS= read -r bug_json; do
            # Extract structured bug data
            local description=$(echo "$bug_json" | jq -r '.description // .error // "No description"')
            local severity=$(echo "$bug_json" | jq -r '.severity // "HIGH"')
            local category=$(echo "$bug_json" | jq -r '.category // "other"')
            local effort=$(echo "$bug_json" | jq -r '.effort // .affected_effort // "UNKNOWN-EFFORT"')
            local fix_status=$(echo "$bug_json" | jq -r '.status // "pending"')

            # Register bug
            log "Registering structured bug: $description"

            local bug_id=$(cascade_register_bug \
                "$cascade_id" \
                "$integration_name" \
                "$severity" \
                "$category" \
                "$description" \
                "$effort")

            if [[ -n "$bug_id" ]]; then
                # If bug was already in progress, update status
                if [[ "$fix_status" == "in_progress" ]] || [[ "$fix_status" == "fixing" ]]; then
                    cascade_start_fix "$bug_id"
                    log "Started fix for migrated bug: $bug_id"
                elif [[ "$fix_status" == "fixed" ]] || [[ "$fix_status" == "complete" ]]; then
                    cascade_start_fix "$bug_id"
                    cascade_complete_fix "$bug_id" "MIGRATED" "true" "Migrated from old format as fixed"
                    log "Marked migrated bug as fixed: $bug_id"
                fi

                log_success "Registered structured bug: $bug_id"
                ((migrated_count++))
            else
                log_error "Failed to register bug: $description"
                ((failed_count++))
            fi

        done < <(jq -c '.project_fixes_in_progress[]?' "$STATE_FILE")
    fi

    # Migrate from fix_cascade.bugs
    if jq -e '.fix_cascade.bugs' "$STATE_FILE" >/dev/null 2>&1; then
        log "Migrating from fix_cascade.bugs..."

        while IFS= read -r bug_json; do
            local description=$(echo "$bug_json" | jq -r '.description // "No description"')
            local severity=$(echo "$bug_json" | jq -r '.severity // "HIGH"')
            local category=$(echo "$bug_json" | jq -r '.category // "other"')
            local effort=$(echo "$bug_json" | jq -r '.effort // "UNKNOWN-EFFORT"')

            local bug_id=$(cascade_register_bug \
                "$cascade_id" \
                "$integration_name" \
                "$severity" \
                "$category" \
                "$description" \
                "$effort")

            if [[ -n "$bug_id" ]]; then
                log_success "Registered cascade bug: $bug_id"
                ((migrated_count++))
            else
                log_error "Failed to register bug: $description"
                ((failed_count++))
            fi

        done < <(jq -c '.fix_cascade.bugs[]?' "$STATE_FILE")
    fi

    log_success "Migration complete: $migrated_count bugs migrated, $failed_count failed"

    if [[ $failed_count -gt 0 ]]; then
        log_error "Some bugs failed to migrate - manual review required"
        return 1
    fi

    return 0
}

create_integration_states() {
    log "Creating integration_fix_states..."

    # Get all unique integration names from bug_registry
    local integrations=$(jq -r '.bug_registry[].detected_in_integration.name' "$STATE_FILE" | sort -u)

    for integration in $integrations; do
        if [[ -z "$integration" ]] || [[ "$integration" == "null" ]]; then
            continue
        fi

        log "Creating state for integration: $integration"

        # Get bugs for this integration
        local bug_ids=$(jq -r ".bug_registry[] | select(.detected_in_integration.name == \"$integration\") | .bug_id" "$STATE_FILE" | jq -R . | jq -s .)
        local bug_count=$(echo "$bug_ids" | jq 'length')

        # Count by status
        local pending=$(jq -r "[.bug_registry[] | select(.detected_in_integration.name == \"$integration\" and .fix_status == \"pending\")] | length" "$STATE_FILE")
        local in_progress=$(jq -r "[.bug_registry[] | select(.detected_in_integration.name == \"$integration\" and .fix_status == \"in_progress\")] | length" "$STATE_FILE")
        local fixed=$(jq -r "[.bug_registry[] | select(.detected_in_integration.name == \"$integration\" and .fix_status == \"fixed\")] | length" "$STATE_FILE")
        local verified=$(jq -r "[.bug_registry[] | select(.detected_in_integration.name == \"$integration\" and .fix_status == \"verified\")] | length" "$STATE_FILE")

        # Determine integration type and phase/wave
        local int_type="wave"
        local phase=1
        local wave=1

        if [[ "$integration" =~ phase([0-9]+)_wave([0-9]+) ]]; then
            phase="${BASH_REMATCH[1]}"
            wave="${BASH_REMATCH[2]}"
            int_type="wave"
        elif [[ "$integration" =~ phase([0-9]+)_integration ]]; then
            phase="${BASH_REMATCH[1]}"
            wave=null
            int_type="phase"
        elif [[ "$integration" =~ project_integration ]]; then
            phase=null
            wave=null
            int_type="project"
        fi

        # Create integration fix state
        jq --arg integration "$integration" \
           --arg type "$int_type" \
           --argjson phase "$phase" \
           --argjson wave "$wave" \
           --argjson bugs "$bug_ids" \
           --argjson total "$bug_count" \
           --argjson pending "$pending" \
           --argjson in_progress "$in_progress" \
           --argjson fixed "$fixed" \
           --argjson verified "$verified" \
           ".integration_fix_states[\$integration] = {
               integration_type: \$type,
               phase: \$phase,
               wave: \$wave,
               integration_name: \$integration,
               cascade_layer: 1,
               cascade_id: .fix_cascade_state.cascade_id,
               bugs_detected: \$bugs,
               bugs_by_status: {
                   pending: \$pending,
                   in_progress: \$in_progress,
                   fixed: \$fixed,
                   verified: \$verified,
                   integrated: 0
               },
               bugs_total_count: \$total,
               status: \"fixing\",
               dependent_integrations: [],
               requires_reintegration: true,
               reintegration_attempts: [],
               reintegration_complete: false,
               created_at: (now | todate),
               updated_at: (now | todate)
           }" "$STATE_FILE" > "${STATE_FILE}.tmp"
        mv "${STATE_FILE}.tmp" "$STATE_FILE"

        log_success "Created integration_fix_state for $integration ($bug_count bugs)"
    done
}

create_effort_states() {
    log "Creating effort_fix_states..."

    # Get all unique efforts from bug_registry
    local efforts=$(jq -r '.bug_registry[].affected_efforts[]?' "$STATE_FILE" | sort -u)

    for effort in $efforts; do
        if [[ -z "$effort" ]] || [[ "$effort" == "null" ]] || [[ "$effort" == "UNKNOWN-EFFORT" ]]; then
            continue
        fi

        log "Creating state for effort: $effort"

        # Get bugs for this effort
        local bug_ids=$(jq -r ".bug_registry[] | select(.affected_efforts[]? == \"$effort\") | .bug_id" "$STATE_FILE" | jq -R . | jq -s .)

        # Count by status
        local pending=$(jq -r "[.bug_registry[] | select(.affected_efforts[]? == \"$effort\" and .fix_status == \"pending\")] | length" "$STATE_FILE")
        local in_progress=$(jq -r "[.bug_registry[] | select(.affected_efforts[]? == \"$effort\" and .fix_status == \"in_progress\")] | length" "$STATE_FILE")
        local fixed=$(jq -r "[.bug_registry[] | select(.affected_efforts[]? == \"$effort\" and .fix_status == \"fixed\")] | length" "$STATE_FILE")
        local verified=$(jq -r "[.bug_registry[] | select(.affected_efforts[]? == \"$effort\" and .fix_status == \"verified\")] | length" "$STATE_FILE")

        # Extract phase/wave from effort name (format: E1.2.1-name or phase1-wave2-effort1)
        local phase=1
        local wave=1

        if [[ "$effort" =~ ^E([0-9]+)\.([0-9]+) ]]; then
            phase="${BASH_REMATCH[1]}"
            wave="${BASH_REMATCH[2]}"
        elif [[ "$effort" =~ phase([0-9]+)-wave([0-9]+) ]]; then
            phase="${BASH_REMATCH[1]}"
            wave="${BASH_REMATCH[2]}"
        fi

        # Create effort fix state
        jq --arg effort "$effort" \
           --argjson phase "$phase" \
           --argjson wave "$wave" \
           --argjson bugs "$bug_ids" \
           --argjson pending "$pending" \
           --argjson in_progress "$in_progress" \
           --argjson fixed "$fixed" \
           --argjson verified "$verified" \
           ".effort_fix_states[\$effort] = {
               effort_id: \$effort,
               phase: \$phase,
               wave: \$wave,
               bugs_assigned: \$bugs,
               bugs_by_status: {
                   pending: \$pending,
                   in_progress: \$in_progress,
                   fixed: \$fixed,
                   verified: \$verified
               },
               fixes_in_progress: [],
               fixes_complete: [],
               ready_for_integration: false,
               last_fix_commit: null,
               fix_branch: \$effort,
               created_at: (now | todate),
               updated_at: (now | todate)
           }" "$STATE_FILE" > "${STATE_FILE}.tmp"
        mv "${STATE_FILE}.tmp" "$STATE_FILE"

        log_success "Created effort_fix_state for $effort ($(echo "$bug_ids" | jq 'length') bugs)"
    done
}

#======================================================================
# VALIDATION
#======================================================================

validate_migration() {
    log "Validating migration..."

    # Run cascade validation
    if cascade_validate; then
        log_success "Cascade validation passed"
    else
        log_error "Cascade validation failed!"
        return 1
    fi

    # Update validation metadata
    cascade_update_validation

    # Verify integration states
    local integration_count=$(jq -r '.integration_fix_states | length' "$STATE_FILE")
    log "Integration states created: $integration_count"

    # Verify effort states
    local effort_count=$(jq -r '.effort_fix_states | length' "$STATE_FILE")
    log "Effort states created: $effort_count"

    # Verify bug registry
    local bug_count=$(jq -r '.bug_registry | length' "$STATE_FILE")
    log "Bugs in registry: $bug_count"

    log_success "Migration validation complete"
}

#======================================================================
# MAIN MIGRATION FLOW
#======================================================================

main() {
    echo -e "${BLUE}╔════════════════════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║   R406 FIX TRACKING MIGRATION TOOL                     ║${NC}"
    echo -e "${BLUE}║   Legacy Format → Fix Cascade Tracking Protocol       ║${NC}"
    echo -e "${BLUE}╚════════════════════════════════════════════════════════╝${NC}"
    echo

    log "Starting migration process..."
    log "State file: $STATE_FILE"
    log "Migration log: $MIGRATION_LOG"

    # Step 1: Pre-migration checks
    pre_migration_checks

    # Step 2: Create backup
    BACKUP_FILE=$(create_backup)

    # Step 3: Extract old bug data
    if ! OLD_BUG_COUNT=$(extract_old_bugs); then
        log_warning "No bugs to migrate - exiting"
        exit 0
    fi

    # Step 4: Determine integration context
    CURRENT_PHASE=$(jq -r '.current_phase' "$STATE_FILE")
    CURRENT_WAVE=$(jq -r '.current_wave' "$STATE_FILE")
    INTEGRATE_WAVE_EFFORTS_NAME="phase${CURRENT_PHASE}_wave${CURRENT_WAVE}"
    INTEGRATE_WAVE_EFFORTS_TYPE="wave"

    log "Migration context: $INTEGRATE_WAVE_EFFORTS_NAME ($INTEGRATE_WAVE_EFFORTS_TYPE)"

    # Step 5: Initialize R406 cascade state
    log "Initializing R406 cascade state..."
    CASCADE_ID=$(cascade_init "$INTEGRATE_WAVE_EFFORTS_NAME" "$INTEGRATE_WAVE_EFFORTS_TYPE" "$STATE_FILE")
    log_success "Cascade initialized: $CASCADE_ID"

    # Step 6: Migrate bugs
    if ! migrate_bugs "$INTEGRATE_WAVE_EFFORTS_NAME" "$INTEGRATE_WAVE_EFFORTS_TYPE" "$CASCADE_ID"; then
        log_error "Bug migration failed!"
        log_error "Backup available at: $BACKUP_FILE"
        exit 1
    fi

    # Step 7: Create integration fix states
    create_integration_states

    # Step 8: Create effort fix states
    create_effort_states

    # Step 9: Validate migration
    if ! validate_migration; then
        log_error "Migration validation failed!"
        log_error "Backup available at: $BACKUP_FILE"
        exit 1
    fi

    # Step 10: Summary
    NEW_BUG_COUNT=$(jq -r '.bug_registry | length' "$STATE_FILE")

    echo
    echo -e "${GREEN}╔════════════════════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║   MIGRATION COMPLETE                                   ║${NC}"
    echo -e "${GREEN}╚════════════════════════════════════════════════════════╝${NC}"
    echo
    log_success "Migrated $NEW_BUG_COUNT bugs from legacy format to R406"
    log_success "Backup saved: $BACKUP_FILE"
    log_success "Migration log: $MIGRATION_LOG"
    echo
    echo -e "${BLUE}Next steps:${NC}"
    echo "  1. Review the migration log: cat $MIGRATION_LOG"
    echo "  2. Validate state: source utilities/cascade-helpers.sh && cascade_status"
    echo "  3. Commit changes: git add orchestrator-state-v3.json && git commit -m 'migrate: R406 fix cascade tracking'"
    echo "  4. Continue with fixes using new system"
}

# Run main if executed directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
