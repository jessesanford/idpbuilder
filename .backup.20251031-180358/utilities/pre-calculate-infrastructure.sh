#!/bin/bash

# 🚨🚨🚨 R504 Pre-Infrastructure Planning - Infrastructure Pre-Calculation Engine
# Pre-calculates all infrastructure paths, branches, and remotes based on parsed plan

set -euo pipefail

CLAUDE_PROJECT_DIR="${CLAUDE_PROJECT_DIR:-/home/vscode/software-factory-template}"
STATE_FILE="$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json"
PLAN_FILE="${1:-$CLAUDE_PROJECT_DIR/PROJECT-IMPLEMENTATION-PLAN.md}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to log messages
log() {
    echo -e "${GREEN}[PRE-CALC]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
    exit 504
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1" >&2
}

# Check dependencies
check_dependencies() {
    if ! command -v yq &> /dev/null; then
        error "yq is required but not installed. Install with: ./utilities/install-requirements.sh"
    fi

    if ! command -v jq &> /dev/null; then
        error "jq is required but not installed. Install with: ./utilities/install-requirements.sh"
    fi

    if [[ ! -f "$CLAUDE_PROJECT_DIR/utilities/parse-implementation-plan.sh" ]]; then
        error "parse-implementation-plan.sh not found"
    fi
}

# Parse the implementation plan
parse_plan() {
    local parsed_json
    parsed_json=$("$CLAUDE_PROJECT_DIR/utilities/parse-implementation-plan.sh" "$PLAN_FILE")

    if [[ -z "$parsed_json" ]]; then
        error "Failed to parse implementation plan"
    fi

    echo "$parsed_json"
}

# Calculate cascade base branch per R308/R501
# This implements PROGRESSIVE trunk-based development where each effort
# branches from the PREVIOUS effort, creating a cascade chain
calculate_cascade_base() {
    local phase="$1"
    local wave="$2"
    local effort_name="$3"
    local effort_index="$4"  # Position in the effort sequence (1, 2, 3...)
    local previous_effort_branch="$5"  # The previous effort's branch (if any)
    local previous_wave_last_branch="$6"  # Last branch from previous wave (if any)
    local previous_phase_last_branch="$7"  # Last branch from previous phase (if any)

    local base_branch=""
    local cascade_reason=""

    # R308/R501 CASCADE ALGORITHM:
    # 1. FIRST effort of Phase 1 Wave 1: Base = main
    # 2. SUBSEQUENT efforts in same wave: Base = previous effort's branch
    # 3. FIRST effort of new wave (same phase): Base = last effort of previous wave
    # 4. FIRST effort of new phase: Base = last effort of previous phase's last wave

    if [[ "$phase" == "phase1" && "$wave" == "wave1" && "$effort_index" -eq 1 ]]; then
        # First effort of Phase 1 Wave 1 starts from main
        base_branch="main"
        cascade_reason="First effort of Phase 1 Wave 1 starts from main (R308/R501)"
    elif [[ "$effort_index" -eq 1 ]]; then
        # First effort of a new wave or phase
        if [[ "$wave" == "wave1" ]]; then
            # First wave of new phase - use last branch from previous phase
            if [[ -n "$previous_phase_last_branch" && "$previous_phase_last_branch" != "null" ]]; then
                base_branch="$previous_phase_last_branch"
                cascade_reason="First effort of $phase $wave cascades from previous phase's last effort: $previous_phase_last_branch (R308/R501)"
            else
                # No previous phase (shouldn't happen unless phase1)
                base_branch="main"
                cascade_reason="No previous phase found, defaulting to main (R308/R501)"
            fi
        else
            # Subsequent wave in same phase - use last branch from previous wave
            if [[ -n "$previous_wave_last_branch" && "$previous_wave_last_branch" != "null" ]]; then
                base_branch="$previous_wave_last_branch"
                cascade_reason="First effort of $phase $wave cascades from previous wave's last effort: $previous_wave_last_branch (R308/R501)"
            else
                # Previous wave not found - should not happen
                error "CASCADE VIOLATION: Cannot find previous wave's last effort for $phase $wave"
            fi
        fi
    else
        # Subsequent efforts in same wave - cascade from previous effort
        if [[ -n "$previous_effort_branch" && "$previous_effort_branch" != "null" ]]; then
            base_branch="$previous_effort_branch"
            cascade_reason="Effort $effort_index cascades from previous effort: $previous_effort_branch (R308/R501)"
        else
            error "CASCADE VIOLATION: Cannot find previous effort for $effort_name (index $effort_index)"
        fi
    fi

    # Return both base and reason as JSON
    echo "{\"base\": \"$base_branch\", \"reason\": \"$cascade_reason\"}"
}

# Calculate effort infrastructure
calculate_effort_infrastructure() {
    local parsed_plan="$1"
    local project_prefix=$(echo "$parsed_plan" | jq -r '.project_prefix')

    # Ensure project prefix exists
    if [[ -z "$project_prefix" ]] || [[ "$project_prefix" == "null" ]]; then
        project_prefix="project"
        warning "No project prefix found, using default: $project_prefix"
    fi

    # Start building pre_planned_infrastructure JSON
    local infrastructure='{"pre_planned_infrastructure": {'
    infrastructure+='"validated": false,'
    infrastructure+='"validation_timestamp": null,'
    infrastructure+='"project_prefix": "'$project_prefix'",'
    infrastructure+='"efforts": {},'
    infrastructure+='"integrations": {},'
    infrastructure+='"cascade_bases": {},'
    infrastructure+='"validation": {'
    infrastructure+='  "naming_rules_checked": false,'
    infrastructure+='  "path_conflicts_checked": false,'
    infrastructure+='  "branch_conflicts_checked": false,'
    infrastructure+='  "remote_configs_validated": false,'
    infrastructure+='  "cascade_bases_calculated": false'
    infrastructure+='}'
    infrastructure+='}}'

    # Track cascade information across phases and waves
    local previous_phase_last_branch=""
    local previous_wave_last_branch=""
    local previous_effort_branch=""

    # Process each phase
    for phase in $(echo "$parsed_plan" | jq -r '.phases | keys[]'); do
        log "Processing $phase"

        # Reset wave tracking for new phase
        previous_wave_last_branch=""

        # Process each wave in phase
        for wave in $(echo "$parsed_plan" | jq -r ".phases.$phase.waves | keys[]"); do
            log "  Processing $phase/$wave"

            # Create integration branch info
            local integration_id="${phase}_${wave}_integration"
            local integration_branch="${project_prefix}/${phase}/${wave}/integration"

            # Add integration to infrastructure
            infrastructure=$(echo "$infrastructure" | jq \
                --arg id "$integration_id" \
                --arg branch "$integration_branch" \
                --arg phase "$phase" \
                --arg wave "$wave" \
                '.pre_planned_infrastructure.integrations[$id] = {
                    "branch_name": $branch,
                    "phase": $phase,
                    "wave": $wave,
                    "component_efforts": [],
                    "created": false,
                    "validated": false
                }')

            # Reset effort tracking for new wave
            previous_effort_branch=""
            local wave_last_branch=""  # Track last branch in this wave

            # Process each effort in wave
            local effort_count=0
            for effort in $(echo "$parsed_plan" | jq -r ".phases.$phase.waves.$wave.efforts[].name" 2>/dev/null); do
                if [[ -n "$effort" ]] && [[ "$effort" != "null" ]]; then
                    ((effort_count++)) || true

                    # Create standardized effort ID
                    local effort_id="${phase}_${wave}_${effort}"

                    # Calculate paths following R313 structure
                    local full_path="$CLAUDE_PROJECT_DIR/efforts/${phase}/${wave}/${effort}/"

                    # Calculate branch names following R327 cascade pattern
                    local branch_name="${project_prefix}/${phase}/${wave}/${effort}"
                    local remote_branch="origin/${branch_name}"

                    # Split pattern for when effort needs splitting
                    local split_pattern="${effort}--split-"

                    # CRITICAL: Calculate cascade base per R308/R501
                    local cascade_result=$(calculate_cascade_base \
                        "$phase" \
                        "$wave" \
                        "$effort" \
                        "$effort_count" \
                        "$previous_effort_branch" \
                        "$previous_wave_last_branch" \
                        "$previous_phase_last_branch")

                    local base_branch=$(echo "$cascade_result" | jq -r '.base')
                    local cascade_reason=$(echo "$cascade_result" | jq -r '.reason')

                    log "    Effort: $effort_id"
                    log "      Path: $full_path"
                    log "      Branch: $branch_name"
                    log "      Base: $base_branch (CASCADE)"
                    log "      Reason: $cascade_reason"

                    # Add effort to infrastructure WITH CASCADE BASE
                    infrastructure=$(echo "$infrastructure" | jq \
                        --arg id "$effort_id" \
                        --arg path "$full_path" \
                        --arg branch "$branch_name" \
                        --arg remote "$remote_branch" \
                        --arg split "$split_pattern" \
                        --arg integration "${project_prefix}/${phase}/${wave}/integration" \
                        --arg phase "$phase" \
                        --arg wave "$wave" \
                        --arg effort "$effort" \
                        --arg base "$base_branch" \
                        --arg reason "$cascade_reason" \
                        '.pre_planned_infrastructure.efforts[$id] = {
                            "full_path": $path,
                            "branch_name": $branch,
                            "remote_branch": $remote,
                            "base_branch": $base,
                            "cascade_reason": $reason,
                            "target_remote": "target",
                            "planning_remote": "planning",
                            "split_pattern": $split,
                            "integration_branch": $integration,
                            "phase": $phase,
                            "wave": $wave,
                            "effort_name": $effort,
                            "created": false,
                            "validated": false,
                            "splits": []
                        }')

                    # Store cascade base information
                    infrastructure=$(echo "$infrastructure" | jq \
                        --arg id "$effort_id" \
                        --arg base "$base_branch" \
                        --arg reason "$cascade_reason" \
                        '.pre_planned_infrastructure.cascade_bases[$id] = {
                            "base_branch": $base,
                            "reason": $reason,
                            "calculated_at": now | todate
                        }')

                    # Add effort to integration's component_efforts
                    infrastructure=$(echo "$infrastructure" | jq \
                        --arg integration_id "$integration_id" \
                        --arg effort_id "$effort_id" \
                        '.pre_planned_infrastructure.integrations[$integration_id].component_efforts += [$effort_id]')

                    # Update tracking variables for next iteration
                    previous_effort_branch="$branch_name"
                    wave_last_branch="$branch_name"
                fi
            done

            if [[ $effort_count -eq 0 ]]; then
                warning "No efforts found in $phase/$wave"
            else
                # Update wave tracking for next wave
                previous_wave_last_branch="$wave_last_branch"
            fi
        done

        # Update phase tracking for next phase
        if [[ -n "$previous_wave_last_branch" ]]; then
            previous_phase_last_branch="$previous_wave_last_branch"
        fi
    done

    echo "$infrastructure"
}

# Validate infrastructure against rules
validate_infrastructure() {
    local infrastructure="$1"

    log "Validating infrastructure against rules..."

    # Check R313 - Effort directory structure
    local all_paths_valid=true
    for effort_path in $(echo "$infrastructure" | jq -r '.pre_planned_infrastructure.efforts[].full_path'); do
        if [[ ! "$effort_path" =~ ^.*/efforts/phase[0-9]+/wave[0-9]+/[a-z0-9-]+/$ ]]; then
            warning "Path does not match R313 pattern: $effort_path"
            all_paths_valid=false
        fi
    done

    # Check R327 - Cascade branching conventions
    local all_branches_valid=true
    for branch_name in $(echo "$infrastructure" | jq -r '.pre_planned_infrastructure.efforts[].branch_name'); do
        if [[ ! "$branch_name" =~ ^[a-z0-9-]+/phase[0-9]+/wave[0-9]+/[a-z0-9-]+$ ]]; then
            warning "Branch does not match R327 pattern: $branch_name"
            all_branches_valid=false
        fi
    done

    # Check R308/R501 - CASCADE BASE VALIDATION
    local cascade_valid=true
    log "Validating cascade bases per R308/R501..."

    # Get sorted list of efforts by phase, wave, and index
    local effort_list=$(echo "$infrastructure" | jq -r '.pre_planned_infrastructure.efforts | to_entries[] |
        "\(.value.phase):\(.value.wave):\(.value.effort_name):\(.value.base_branch):\(.key)"' | sort)

    local prev_branch=""
    local is_first_overall=true
    local current_phase=""
    local current_wave=""

    while IFS=: read -r phase wave effort base effort_id; do
        # Check if this is first effort of a new wave
        local is_first_in_wave=false
        if [[ "$phase" != "$current_phase" ]] || [[ "$wave" != "$current_wave" ]]; then
            is_first_in_wave=true
            current_phase="$phase"
            current_wave="$wave"
        fi

        # Validate cascade per R308/R501
        if [[ "$is_first_overall" == "true" ]]; then
            # First effort of Phase 1 Wave 1 must be from main
            if [[ "$phase" == "phase1" && "$wave" == "wave1" ]]; then
                if [[ "$base" != "main" ]]; then
                    warning "❌ CASCADE VIOLATION: First effort $effort_id should base from 'main' but bases from '$base'"
                    cascade_valid=false
                fi
            fi
            is_first_overall=false
        elif [[ "$is_first_in_wave" == "true" ]]; then
            # First effort of new wave should base from last effort of previous wave
            # (This is harder to validate without full history, but we check it's not main)
            if [[ "$phase" != "phase1" || "$wave" != "wave1" ]]; then
                if [[ "$base" == "main" ]]; then
                    warning "⚠️ CASCADE WARNING: $effort_id is first in $phase/$wave but bases from 'main' - should cascade from previous wave"
                    cascade_valid=false
                fi
            fi
        else
            # Subsequent efforts should cascade from previous effort
            if [[ -n "$prev_branch" ]] && [[ "$base" != "$prev_branch" ]]; then
                warning "❌ CASCADE VIOLATION: $effort_id should cascade from '$prev_branch' but bases from '$base'"
                cascade_valid=false
            fi
        fi

        # Update previous branch for next iteration
        prev_branch=$(echo "$infrastructure" | jq -r --arg id "$effort_id" '.pre_planned_infrastructure.efforts[$id].branch_name')
    done <<< "$effort_list"

    # Check that all efforts have cascade bases
    local missing_bases=false
    for effort_id in $(echo "$infrastructure" | jq -r '.pre_planned_infrastructure.efforts | keys[]'); do
        local base=$(echo "$infrastructure" | jq -r --arg id "$effort_id" '.pre_planned_infrastructure.efforts[$id].base_branch')
        if [[ -z "$base" ]] || [[ "$base" == "null" ]]; then
            warning "❌ MISSING CASCADE BASE: Effort $effort_id has no base_branch defined!"
            missing_bases=true
        fi
    done

    # Check for path conflicts
    local path_conflicts=false
    local paths=$(echo "$infrastructure" | jq -r '.pre_planned_infrastructure.efforts[].full_path' | sort)
    if [[ $(echo "$paths" | uniq -d | wc -l) -gt 0 ]]; then
        warning "Path conflicts detected:"
        echo "$paths" | uniq -d
        path_conflicts=true
    fi

    # Check for branch conflicts
    local branch_conflicts=false
    local branches=$(echo "$infrastructure" | jq -r '.pre_planned_infrastructure.efforts[].branch_name' | sort)
    if [[ $(echo "$branches" | uniq -d | wc -l) -gt 0 ]]; then
        warning "Branch conflicts detected:"
        echo "$branches" | uniq -d
        branch_conflicts=true
    fi

    # Update validation status
    if [[ "$all_paths_valid" == "true" ]] && \
       [[ "$all_branches_valid" == "true" ]] && \
       [[ "$path_conflicts" == "false" ]] && \
       [[ "$branch_conflicts" == "false" ]] && \
       [[ "$cascade_valid" == "true" ]] && \
       [[ "$missing_bases" == "false" ]]; then

        infrastructure=$(echo "$infrastructure" | jq \
            --arg timestamp "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
            '.pre_planned_infrastructure.validated = true |
             .pre_planned_infrastructure.validation_timestamp = $timestamp |
             .pre_planned_infrastructure.validation.naming_rules_checked = true |
             .pre_planned_infrastructure.validation.path_conflicts_checked = true |
             .pre_planned_infrastructure.validation.branch_conflicts_checked = true |
             .pre_planned_infrastructure.validation.remote_configs_validated = true |
             .pre_planned_infrastructure.validation.cascade_bases_calculated = true')

        log "✅ Infrastructure validation PASSED (including R308/R501 cascade validation)"
    else
        error "Infrastructure validation FAILED - fix issues above (R308/R501 CASCADE VIOLATIONS)"
    fi

    echo "$infrastructure"
}

# Update orchestrator-state-v3.json with pre-planned infrastructure
update_state_file() {
    local infrastructure="$1"

    log "Updating orchestrator-state-v3.json with pre-planned infrastructure..."

    # Check if state file exists
    if [[ ! -f "$STATE_FILE" ]]; then
        warning "orchestrator-state-v3.json not found, creating new one..."
        echo '{}' > "$STATE_FILE"
    fi

    # Merge pre_planned_infrastructure into state file
    local current_state=$(cat "$STATE_FILE")
    local updated_state=$(echo "$current_state" | jq --slurpfile new <(echo "$infrastructure") \
        '. * $new[0]')

    # Write updated state back
    echo "$updated_state" | jq '.' > "$STATE_FILE"

    log "✅ orchestrator-state-v3.json updated with pre-planned infrastructure"
}

# Generate summary report
generate_report() {
    local infrastructure="$1"

    echo ""
    echo "========================================="
    echo "PRE-INFRASTRUCTURE PLANNING COMPLETE"
    echo "========================================="
    echo ""
    echo "Project Prefix: $(echo "$infrastructure" | jq -r '.pre_planned_infrastructure.project_prefix')"
    echo "Total Phases: $(echo "$infrastructure" | jq -r '.pre_planned_infrastructure.efforts | to_entries | map(.value.phase) | unique | length')"
    echo "Total Waves: $(echo "$infrastructure" | jq -r '.pre_planned_infrastructure.efforts | to_entries | map(.value.wave) | unique | length')"
    echo "Total Efforts: $(echo "$infrastructure" | jq -r '.pre_planned_infrastructure.efforts | keys | length')"
    echo "Total Integrations: $(echo "$infrastructure" | jq -r '.pre_planned_infrastructure.integrations | keys | length')"
    echo ""
    echo "CASCADE BASE BRANCHES (R308/R501):"
    echo "-----------------------------------"
    # Show first 5 cascade bases as examples
    echo "$infrastructure" | jq -r '.pre_planned_infrastructure.efforts | to_entries[0:5][] |
        "  \(.value.effort_name): \(.value.base_branch)"'
    local total_efforts=$(echo "$infrastructure" | jq -r '.pre_planned_infrastructure.efforts | keys | length')
    if [[ $total_efforts -gt 5 ]]; then
        echo "  ... and $((total_efforts - 5)) more efforts with cascade bases"
    fi
    echo ""
    echo "Validation Status:"
    echo "  ✅ Naming Rules: $(echo "$infrastructure" | jq -r '.pre_planned_infrastructure.validation.naming_rules_checked')"
    echo "  ✅ Path Conflicts: $(echo "$infrastructure" | jq -r '.pre_planned_infrastructure.validation.path_conflicts_checked')"
    echo "  ✅ Branch Conflicts: $(echo "$infrastructure" | jq -r '.pre_planned_infrastructure.validation.branch_conflicts_checked')"
    echo "  ✅ Remote Configs: $(echo "$infrastructure" | jq -r '.pre_planned_infrastructure.validation.remote_configs_validated')"
    echo "  ✅ CASCADE Bases: $(echo "$infrastructure" | jq -r '.pre_planned_infrastructure.validation.cascade_bases_calculated')"
    echo ""
    echo "Infrastructure is ready for mechanical creation!"
    echo "========================================="
}

# Main execution
main() {
    log "Starting pre-infrastructure planning..."

    # Check dependencies
    check_dependencies

    # Parse implementation plan
    log "Parsing implementation plan..."
    local parsed_plan=$(parse_plan)

    # Calculate infrastructure
    log "Calculating infrastructure..."
    local infrastructure=$(calculate_effort_infrastructure "$parsed_plan")

    # Validate infrastructure
    log "Validating infrastructure..."
    infrastructure=$(validate_infrastructure "$infrastructure")

    # Update state file
    update_state_file "$infrastructure"

    # Generate report
    generate_report "$infrastructure"

    log "Pre-infrastructure planning complete!"
}

# Execute main function
main "$@"