# /sub-orchestrate-integration Command

## Purpose
Entry point for INTEGRATE_WAVE_EFFORTS sub-orchestrator. Handles branch integration independently from main orchestrator for wave, phase, or project level merges.

## Usage
```bash
claude -p "$CLAUDE_PROJECT_DIR" \
  --command "/sub-orchestrate-integration" \
  --params "file=/tmp/params-uuid.json" \
  --headless \
  --state "sub-state-integration-uuid.json"
```

## Parameter File Format
```json
{
  "sub_orchestrator_type": "INTEGRATE_WAVE_EFFORTS",
  "unique_id": "uuid-v4",
  "master_state_file": "orchestrator-state-v3.json",
  "input_parameters": {
    "integration_type": "WAVE|PHASE|PROJECT",
    "branches_to_merge": ["effort-1", "effort-2", "effort-3"],
    "target_branch": "wave-1-integration",
    "validation_level": "BASIC|FULL|COMPREHENSIVE",
    "conflict_resolution": "MANUAL|AUTO_THEIRS|AUTO_OURS",
    "merge_strategy": "SEQUENTIAL|OCTOPUS",
    "build_required": true,
    "test_required": true
  },
  "output_location": "/tmp/sub-orch-uuid/output.json",
  "heartbeat_file": "/tmp/sub-orch-uuid/heartbeat.json",
  "checkpoint_file": "/tmp/sub-orch-uuid/checkpoint.json",
  "max_duration_seconds": 7200,
  "checkpoint_interval": 30
}
```

## Execution Flow

### 1. INITIALIZATION
```bash
# Parse parameters
PARAM_FILE="${1:-/tmp/params.json}"
PARAMS=$(cat "$PARAM_FILE")
SUB_ID=$(echo "$PARAMS" | jq -r '.unique_id')
INTEGRATE_WAVE_EFFORTS_TYPE=$(echo "$PARAMS" | jq -r '.input_parameters.integration_type')

# Set up environment
export SUB_ORCHESTRATOR_ID="$SUB_ID"
export SUB_TYPE="INTEGRATE_WAVE_EFFORTS"
export INTEGRATE_WAVE_EFFORTS_LEVEL="$INTEGRATE_WAVE_EFFORTS_TYPE"

# Initialize tracking
initialize_integration_tracking
write_initial_heartbeat
```

### 2. MAIN EXECUTION
```bash
# Integration state machine
execute_integration() {
  case "$CURRENT_STATE" in
    "INTEGRATE_WAVE_EFFORTS_INIT")
      validate_branches_exist
      check_branch_status
      ;;
    "INTEGRATE_WAVE_EFFORTS_CREATE_TARGET")
      create_integration_branch
      ;;
    "INTEGRATE_WAVE_EFFORTS_MERGE_BRANCHES")
      perform_sequential_merges
      ;;
    "INTEGRATE_WAVE_EFFORTS_RESOLVE_CONFLICTS")
      handle_merge_conflicts
      ;;
    "INTEGRATE_WAVE_EFFORTS_BUILD")
      run_build_validation
      ;;
    "INTEGRATE_WAVE_EFFORTS_TEST")
      run_test_suites
      ;;
    "INTEGRATE_WAVE_EFFORTS_VALIDATE")
      validate_integration_complete
      ;;
    "INTEGRATE_WAVE_EFFORTS_COMPLETE")
      finalize_integration
      ;;
  esac
}
```

### 3. MERGE STRATEGY IMPLEMENTATION

#### Sequential Merge
```bash
sequential_merge() {
  local BRANCHES=("$@")
  local MERGED=()

  for BRANCH in "${BRANCHES[@]}"; do
    echo "Merging $BRANCH into $TARGET_BRANCH"

    if git merge --no-ff "$BRANCH"; then
      MERGED+=("$BRANCH")
      write_checkpoint "MERGE_PROJECT_DONE" "$BRANCH"
    else
      handle_merge_conflict "$BRANCH"
    fi

    update_progress $((${#MERGED[@]} * 100 / ${#BRANCHES[@]}))
  done
}
```

#### Octopus Merge
```bash
octopus_merge() {
  local BRANCHES=("$@")

  git merge --no-ff "${BRANCHES[@]}" || {
    # Octopus failed, fall back to sequential
    echo "Octopus merge failed, using sequential"
    sequential_merge "${BRANCHES[@]}"
  }
}
```

### 4. CONFLICT RESOLUTION
```bash
handle_conflicts() {
  local CONFLICT_STRATEGY="$1"

  case "$CONFLICT_STRATEGY" in
    "AUTO_THEIRS")
      git checkout --theirs .
      git add -A
      git commit --no-edit
      ;;
    "AUTO_OURS")
      git checkout --ours .
      git add -A
      git commit --no-edit
      ;;
    "MANUAL")
      # Record conflicts for manual resolution
      record_conflicts_for_escalation
      set_status "NEEDS_MANUAL_RESOLUTION"
      ;;
  esac
}
```

### 5. VALIDATION LEVELS

#### BASIC Validation
```bash
basic_validation() {
  # Syntax checks
  check_syntax_all_files

  # Basic linting
  run_linters

  # Smoke tests
  run_smoke_tests
}
```

#### FULL Validation
```bash
full_validation() {
  # Everything from BASIC plus:
  basic_validation

  # Full build
  run_full_build

  # Unit tests
  run_unit_tests

  # Integration tests
  run_integration_tests
}
```

#### COMPREHENSIVE Validation
```bash
comprehensive_validation() {
  # Everything from FULL plus:
  full_validation

  # Performance tests
  run_performance_tests

  # Security scans
  run_security_scans

  # E2E tests
  run_e2e_tests
}
```

### 6. OUTPUT GENERATION
```bash
write_integration_output() {
  cat > "$OUTPUT_FILE" <<EOF
{
  "status": "$STATUS",
  "integration_type": "$INTEGRATE_WAVE_EFFORTS_TYPE",
  "integration_branch": "$TARGET_BRANCH",
  "merge_results": {
    $(generate_merge_results)
  },
  "conflicts_encountered": $CONFLICTS,
  "conflicts_resolved": $RESOLVED,
  "validation_results": {
    "build_status": "$BUILD_STATUS",
    "test_status": "$TEST_STATUS",
    "validation_level": "$VALIDATION_LEVEL"
  },
  "duration_seconds": $DURATION,
  "next_action": "$NEXT_ACTION"
}
EOF
}
```

## State Machine States
- `INTEGRATE_WAVE_EFFORTS_INIT` - Validate prerequisites
- `INTEGRATE_WAVE_EFFORTS_CREATE_TARGET` - Create integration branch
- `INTEGRATE_WAVE_EFFORTS_MERGE_BRANCHES` - Perform merges
- `INTEGRATE_WAVE_EFFORTS_RESOLVE_CONFLICTS` - Handle conflicts
- `INTEGRATE_WAVE_EFFORTS_BUILD` - Build validation
- `INTEGRATE_WAVE_EFFORTS_TEST` - Test execution
- `INTEGRATE_WAVE_EFFORTS_VALIDATE` - Final validation
- `INTEGRATE_WAVE_EFFORTS_ROLLBACK` - Rollback on failure
- `INTEGRATE_WAVE_EFFORTS_COMPLETE` - Finalize

## Heartbeat Updates
```bash
update_heartbeat() {
  cat > "$HEARTBEAT_FILE" <<EOF
{
  "pid": $$,
  "status": "$STATUS",
  "progress_percentage": $PROGRESS,
  "current_state": "$CURRENT_STATE",
  "last_update": "$(date -Iseconds)",
  "branches_merged": ${#MERGED[@]},
  "branches_total": ${#BRANCHES[@]},
  "conflicts": $CONFLICTS,
  "build_status": "$BUILD_STATUS",
  "test_status": "$TEST_STATUS"
}
EOF
}
```

## Recovery Support
```bash
# Checkpoint at each successful merge
checkpoint_merge() {
  local BRANCH="$1"
  cat > "$CHECKPOINT_FILE" <<EOF
{
  "state": "INTEGRATE_WAVE_EFFORTS_MERGE_BRANCHES",
  "data": {
    "merged_branches": $(json_array "${MERGED[@]}"),
    "current_branch": "$BRANCH",
    "target_branch": "$TARGET_BRANCH"
  },
  "can_resume": true
}
EOF
}
```

## Error Handling
```bash
handle_integration_error() {
  local ERROR_TYPE="$1"

  case "$ERROR_TYPE" in
    "MERGE_CONFLICT")
      if [[ "$CONFLICT_RESOLUTION" == "MANUAL" ]]; then
        escalate_for_manual_resolution
      else
        attempt_auto_resolution
      fi
      ;;
    "BUILD_FAILURE")
      capture_build_logs
      determine_if_retryable
      ;;
    "TEST_FAILURE")
      capture_test_results
      identify_failing_tests
      ;;
  esac
}
```

## Rules Applied
- R377: Communication Protocol
- R378: Lifecycle Management
- R379: Monitoring and Recovery
- R361: Integration Testing Requirements
- R362: Architecture Validation

## Success Criteria
- All branches successfully merged
- No unresolved conflicts
- Build passes (if required)
- Tests pass (if required)
- Integration branch created
- Output file written

## Common Issues
1. Merge conflicts between branches
2. Build failures after integration
3. Test failures in integrated code
4. Timeout on large merges
5. Git history complications

## Example Launch
```bash
# Prepare parameters
cat > /tmp/integration-params.json <<'EOF'
{
  "sub_orchestrator_type": "INTEGRATE_WAVE_EFFORTS",
  "unique_id": "integ-wave1-uuid",
  "input_parameters": {
    "integration_type": "WAVE",
    "branches_to_merge": ["effort-1", "effort-2"],
    "target_branch": "wave-1-integration",
    "validation_level": "FULL",
    "conflict_resolution": "MANUAL"
  }
}
EOF

# Launch sub-orchestrator
/sub-orchestrate-integration /tmp/integration-params.json
```