#!/bin/bash
# Update split-001 review status
yq -i '.split_tracking."gitea-client".splits[0].review_status = "APPROVED"' orchestrator-state.json
yq -i '.split_tracking."gitea-client".splits[0].review_complete_time = "2025-09-08T06:15:00Z"' orchestrator-state.json
yq -i '.split_tracking."gitea-client".splits[0].actual_lines = 629' orchestrator-state.json
yq -i '.split_tracking."gitea-client".splits[0].review_report = "CODE-REVIEW-REPORT-SPLIT001-20250908-061509.md"' orchestrator-state.json
yq -i '.split_tracking."gitea-client".splits[0].measurement_corrected = true' orchestrator-state.json
yq -i '.split_tracking."gitea-client".splits[0].note = "Split-001 APPROVED - 629 lines when measured correctly from Phase1/integration"' orchestrator-state.json

# Update current state (R324 CRITICAL - MUST UPDATE BEFORE STOPPING!)
yq -i '.state_machine.previous_state = "SPAWN_CODE_REVIEWERS_FOR_REVIEW"' orchestrator-state.json
yq -i '.state_machine.current_state = "MONITOR_REVIEWS"' orchestrator-state.json
yq -i '.state_machine.transition_time = "2025-09-08T06:16:00Z"' orchestrator-state.json
yq -i '.state_machine.transition_reason = "E2.1.2 Split-001 review complete and approved"' orchestrator-state.json

# Add to next actions
yq -i '.next_actions = ["Monitor for Split-001 integration readiness", "Prepare for Split-002 creation", "Check if E2.1.1 fixes are complete"]' orchestrator-state.json

echo "✅ State file updated successfully"
