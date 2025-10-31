#!/bin/bash
# compress-orchestrator-state.sh
# Applies Phase 1 optimizations to orchestrator-state-v3.json
# Reduces file size by ~15-23% with zero semantic loss

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
STATE_FILE="${1:-$PROJECT_ROOT/orchestrator-state-v3.json}"

if [[ ! -f "$STATE_FILE" ]]; then
    echo "ERROR: State file not found: $STATE_FILE"
    exit 1
fi

echo "Compressing orchestrator-state-v3.json..."
echo "Input file: $STATE_FILE"
echo "Original size: $(wc -c < "$STATE_FILE") bytes"

# Create backup
BACKUP_FILE="${STATE_FILE}.backup-$(date +%Y%m%d-%H%M%S)"
cp "$STATE_FILE" "$BACKUP_FILE"
echo "Backup created: $BACKUP_FILE"

# Apply compression using jq
jq '
  # Phase 1: Top-level key compression
  with_entries(
    if .key == "effort_planning_in_progress" then .key = "eff_plan_prog"
    elif .key == "effort_plans_completed" then .key = "eff_plan_done"
    elif .key == "effort_plans_in_progress" then .key = "eff_plan_wip"
    elif .key == "efforts_in_progress" then .key = "eff_wip"
    elif .key == "efforts_needing_fixes" then .key = "eff_fix_need"
    elif .key == "phase_integration_infrastructure" then .key = "ph_integ_infra"
    elif .key == "phase_integration_results" then .key = "ph_integ_res"
    elif .key == "code_reviewer_parallelization_plan" then .key = "cr_para_plan"
    elif .key == "code_reviewer_spawn_info" then .key = "cr_spawn"
    elif .key == "state_transition_log" then .key = "state_log"
    elif .key == "transition_history" then .key = "trans_hist"
    elif .key == "pre_planned_infrastructure" then .key = "preplan_infra"
    elif .key == "upstream_bugs_wave2" then .key = "bugs_w2"
    elif .key == "phase2_wave1_parallelization_analysis" then .key = "ph2w1_para"
    elif .key == "phase_1_integration_status" then .key = "ph1_integ_st"
    elif .key == "r291_gate_status_wave2" then .key = "r291_w2"
    elif .key == "split_parallelization_plan" then .key = "split_para"
    elif .key == "stale_integration_tracking" then .key = "stale_integ"
    elif .key == "project_integration_code_review" then .key = "proj_integ_cr"
    elif .key == "project_integration_fix_planning" then .key = "proj_integ_fix"
    elif .key == "project_integration_results" then .key = "proj_integ_res"
    else .
    end
  ) |

  # Phase 2: Recursive nested key compression for common patterns
  walk(
    if type == "object" then
      with_entries(
        if .key == "detected_in_integration" then .key = "det_in_integ"
        elif .key == "detected_at" then .key = "det_at"
        elif .key == "detected_by" then .key = "det_by"
        elif .key == "error_details" then .key = "err_det"
        elif .key == "error_message" then .key = "err_msg"
        elif .key == "affected_efforts" then .key = "aff_eff"
        elif .key == "primary_effort" then .key = "pri_eff"
        elif .key == "requires_coordinated_fix" then .key = "coord_fix"
        elif .key == "fix_attempts" then .key = "fix_att"
        elif .key == "attempt_number" then .key = "att_num"
        elif .key == "integration_status" then .key = "integ_st"
        elif .key == "integrated_at" then .key = "integ_at"
        elif .key == "integration_commit" then .key = "integ_cmt"
        elif .key == "resolution_notes" then .key = "res_notes"
        elif .key == "total_splits" then .key = "tot_sp"
        elif .key == "current_split" then .key = "cur_sp"
        elif .key == "effort_name" then .key = "eff_nm"
        elif .key == "base_branch" then .key = "base_br"
        elif .key == "infrastructure_created" then .key = "infra_cr"
        elif .key == "implementation_complete" then .key = "impl_done"
        elif .key == "actual_lines" then .key = "act_ln"
        elif .key == "target_lines" then .key = "tgt_ln"
        elif .key == "r359_fix_complete" then .key = "r359_done"
        elif .key == "split_reason" then .key = "sp_why"
        elif .key == "split_plan_location" then .key = "sp_plan_loc"
        else .
        end
      )
    else .
    end
  ) |

  # Phase 3: Internal enum compression
  walk(
    if type == "string" then
      if . == "COMPLETED" then "C"
      elif . == "PENDING" then "P"
      elif . == "IN_PROGRESS" then "W"
      elif . == "FAILED" then "F"
      elif . == "PROJECT_DONEFUL" then "S"
      elif . == "NOT_STARTED" then "N"
      elif . == "not_integrated" then "ni"
      elif . == "integrated" then "i"
      elif . == "integration_pending" then "ip"
      elif . == "stale" then "st"
      elif . == "CRITICAL" then "CR"
      elif . == "HIGH" then "H"
      elif . == "MEDIUM" then "M"
      elif . == "LOW" then "L"
      else .
      end
    else .
    end
  )
' "$STATE_FILE" > "${STATE_FILE}.tmp"

# Replace original with compressed version
mv "${STATE_FILE}.tmp" "$STATE_FILE"

NEW_SIZE=$(wc -c < "$STATE_FILE")
OLD_SIZE=$(wc -c < "$BACKUP_FILE")
SAVINGS=$((OLD_SIZE - NEW_SIZE))
PERCENT=$((SAVINGS * 100 / OLD_SIZE))

echo ""
echo "Compression complete!"
echo "New size: $NEW_SIZE bytes"
echo "Savings: $SAVINGS bytes ($PERCENT% reduction)"
echo ""
echo "Backup preserved at: $BACKUP_FILE"
