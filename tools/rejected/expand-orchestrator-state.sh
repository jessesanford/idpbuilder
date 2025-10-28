#!/bin/bash
# expand-orchestrator-state.sh - Reverses compression for debugging

set -euo pipefail

STATE_FILE="${1:-orchestrator-state-v3.json}"

if [[ ! -f "$STATE_FILE" ]]; then
    echo "ERROR: State file not found: $STATE_FILE" >&2
    exit 1
fi

# Reverse all compression transformations
jq '
  with_entries(
    if .key == "eff_plan_prog" then .key = "effort_planning_in_progress"
    elif .key == "eff_plan_done" then .key = "effort_plans_completed"
    elif .key == "eff_plan_wip" then .key = "effort_plans_in_progress"
    elif .key == "eff_wip" then .key = "efforts_in_progress"
    elif .key == "eff_fix_need" then .key = "efforts_needing_fixes"
    elif .key == "ph_integ_infra" then .key = "phase_integration_infrastructure"
    elif .key == "ph_integ_res" then .key = "phase_integration_results"
    elif .key == "cr_para_plan" then .key = "code_reviewer_parallelization_plan"
    elif .key == "cr_spawn" then .key = "code_reviewer_spawn_info"
    elif .key == "state_log" then .key = "state_transition_log"
    elif .key == "trans_hist" then .key = "transition_history"
    elif .key == "preplan_infra" then .key = "pre_planned_infrastructure"
    elif .key == "bugs_w2" then .key = "upstream_bugs_wave2"
    elif .key == "ph2w1_para" then .key = "phase2_wave1_parallelization_analysis"
    elif .key == "ph1_integ_st" then .key = "phase_1_integration_status"
    elif .key == "r291_w2" then .key = "r291_gate_status_wave2"
    elif .key == "split_para" then .key = "split_parallelization_plan"
    elif .key == "stale_integ" then .key = "stale_integration_tracking"
    elif .key == "proj_integ_cr" then .key = "project_integration_code_review"
    elif .key == "proj_integ_fix" then .key = "project_integration_fix_planning"
    elif .key == "proj_integ_res" then .key = "project_integration_results"
    else .
    end
  ) |
  walk(
    if type == "object" then
      with_entries(
        if .key == "det_in_integ" then .key = "detected_in_integration"
        elif .key == "det_at" then .key = "detected_at"
        elif .key == "det_by" then .key = "detected_by"
        elif .key == "err_det" then .key = "error_details"
        elif .key == "err_msg" then .key = "error_message"
        elif .key == "aff_eff" then .key = "affected_efforts"
        elif .key == "pri_eff" then .key = "primary_effort"
        elif .key == "coord_fix" then .key = "requires_coordinated_fix"
        elif .key == "fix_att" then .key = "fix_attempts"
        elif .key == "att_num" then .key = "attempt_number"
        elif .key == "integ_st" then .key = "integration_status"
        elif .key == "integ_at" then .key = "integrated_at"
        elif .key == "integ_cmt" then .key = "integration_commit"
        elif .key == "res_notes" then .key = "resolution_notes"
        elif .key == "tot_sp" then .key = "total_splits"
        elif .key == "cur_sp" then .key = "current_split"
        elif .key == "eff_nm" then .key = "effort_name"
        elif .key == "base_br" then .key = "base_branch"
        elif .key == "infra_cr" then .key = "infrastructure_created"
        elif .key == "impl_done" then .key = "implementation_complete"
        elif .key == "act_ln" then .key = "actual_lines"
        elif .key == "tgt_ln" then .key = "target_lines"
        elif .key == "r359_done" then .key = "r359_fix_complete"
        elif .key == "sp_why" then .key = "split_reason"
        elif .key == "sp_plan_loc" then .key = "split_plan_location"
        else .
        end
      )
    else .
    end
  ) |
  walk(
    if type == "string" then
      if . == "C" then "COMPLETED"
      elif . == "P" then "PENDING"
      elif . == "W" then "IN_PROGRESS"
      elif . == "F" then "FAILED"
      elif . == "S" then "PROJECT_DONEFUL"
      elif . == "N" then "NOT_STARTED"
      elif . == "ni" then "not_integrated"
      elif . == "i" then "integrated"
      elif . == "ip" then "integration_pending"
      elif . == "st" then "stale"
      elif . == "CR" then "CRITICAL"
      elif . == "H" then "HIGH"
      elif . == "M" then "MEDIUM"
      elif . == "L" then "LOW"
      else .
      end
    else .
    end
  )
' "$STATE_FILE"
