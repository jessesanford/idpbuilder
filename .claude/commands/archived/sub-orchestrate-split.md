# /sub-orchestrate-split Command

## Purpose
Entry point for SPLIT_COORDINATION sub-orchestrator. Handles effort splitting when size limits are exceeded, coordinating sequential implementation of splits.

## Usage
```bash
claude -p "$CLAUDE_PROJECT_DIR" \
  --command "/sub-orchestrate-split" \
  --params "file=/tmp/params-uuid.json" \
  --headless \
  --state "sub-state-split-uuid.json"
```

## Parameter File Format
```json
{
  "sub_orchestrator_type": "SPLIT_COORDINATION",
  "unique_id": "uuid-v4",
  "master_state_file": "orchestrator-state.json",
  "input_parameters": {
    "original_effort": "E1.1",
    "split_plan": {
      "splits": ["E1.1a", "E1.1b", "E1.1c"],
      "dependencies": {
        "E1.1b": ["E1.1a"],
        "E1.1c": ["E1.1b"]
      },
      "sequential": true,
      "estimated_sizes": {
        "E1.1a": 650,
        "E1.1b": 680,
        "E1.1c": 450
      }
    },
    "size_limit": 700,
    "validation_after_each": true,
    "review_required": true
  },
  "output_location": "/tmp/sub-orch-uuid/output.json",
  "heartbeat_file": "/tmp/sub-orch-uuid/heartbeat.json",
  "checkpoint_file": "/tmp/sub-orch-uuid/checkpoint.json",
  "max_duration_seconds": 10800,
  "checkpoint_interval": 30
}
```

## Execution Flow

### 1. INITIALIZATION
```bash
# Parse split plan
PARAM_FILE="${1:-/tmp/params.json}"
PARAMS=$(cat "$PARAM_FILE")
SUB_ID=$(echo "$PARAMS" | jq -r '.unique_id')
ORIGINAL_EFFORT=$(echo "$PARAMS" | jq -r '.input_parameters.original_effort')
SPLITS=($(echo "$PARAMS" | jq -r '.input_parameters.split_plan.splits[]'))

# Set up split tracking
export SUB_ORCHESTRATOR_ID="$SUB_ID"
export SUB_TYPE="SPLIT_COORDINATION"
export TOTAL_SPLITS="${#SPLITS[@]}"
export COMPLETED_SPLITS=0

# Initialize state
initialize_split_tracking
write_initial_heartbeat
```

### 2. MAIN EXECUTION
```bash
# Split coordination state machine
execute_split_coordination() {
  case "$CURRENT_STATE" in
    "SPLIT_INIT")
      validate_split_plan
      prepare_split_infrastructure
      ;;
    "SPLIT_CREATE_BRANCHES")
      create_split_branches
      ;;
    "SPLIT_SPAWN_ENGINEER")
      spawn_engineer_for_current_split
      ;;
    "SPLIT_MONITOR_IMPLEMENTATION")
      monitor_split_progress
      ;;
    "SPLIT_MEASURE_SIZE")
      measure_split_size
      validate_within_limit
      ;;
    "SPLIT_SPAWN_REVIEWER")
      spawn_reviewer_for_split
      ;;
    "SPLIT_MONITOR_REVIEW")
      monitor_review_progress
      ;;
    "SPLIT_NEXT_SPLIT")
      advance_to_next_split
      ;;
    "SPLIT_MERGE_SPLITS")
      merge_all_splits
      ;;
    "SPLIT_COMPLETE")
      finalize_split_coordination
      ;;
  esac
}
```

### 3. SEQUENTIAL SPLIT HANDLING
```bash
handle_sequential_splits() {
  local SPLITS=("$@")

  for SPLIT in "${SPLITS[@]}"; do
    echo "Processing split: $SPLIT"

    # Create split branch
    create_split_branch "$SPLIT"

    # Spawn SW engineer
    spawn_engineer "$SPLIT"

    # Monitor until complete
    monitor_until_complete "$SPLIT"

    # Measure size
    SIZE=$(measure_branch_size "$SPLIT")

    # Validate size limit
    if [[ $SIZE -gt $SIZE_LIMIT ]]; then
      handle_oversized_split "$SPLIT" "$SIZE"
    fi

    # Spawn reviewer if required
    if [[ "$REVIEW_REQUIRED" == "true" ]]; then
      spawn_reviewer "$SPLIT"
      monitor_review "$SPLIT"
    fi

    # Update progress
    COMPLETED_SPLITS=$((COMPLETED_SPLITS + 1))
    update_progress $((COMPLETED_SPLITS * 100 / TOTAL_SPLITS))

    # Checkpoint
    checkpoint_split_completion "$SPLIT"
  done
}
```

### 4. DEPENDENCY MANAGEMENT
```bash
check_dependencies() {
  local SPLIT="$1"
  local DEPS=$(echo "$PARAMS" | jq -r ".input_parameters.split_plan.dependencies.$SPLIT[]" 2>/dev/null)

  if [[ -n "$DEPS" ]]; then
    for DEP in $DEPS; do
      if ! is_split_complete "$DEP"; then
        echo "Waiting for dependency: $DEP"
        return 1
      fi
    done
  fi

  return 0  # All dependencies satisfied
}
```

### 5. SIZE MEASUREMENT
```bash
measure_split_size() {
  local SPLIT="$1"
  local BRANCH="split/$SPLIT"

  # Use line counter tool
  if [[ -f "$CLAUDE_PROJECT_DIR/tools/line-counter.sh" ]]; then
    SIZE=$($CLAUDE_PROJECT_DIR/tools/line-counter.sh "$BRANCH")
  else
    # Fallback measurement
    SIZE=$(git diff main.."$BRANCH" --stat | tail -1 | awk '{print $4}')
  fi

  echo "Split $SPLIT size: $SIZE lines"

  # Record measurement
  record_split_measurement "$SPLIT" "$SIZE"

  echo "$SIZE"
}
```

### 6. SPLIT MERGING
```bash
merge_splits_to_original() {
  local ORIGINAL="$1"
  shift
  local SPLITS=("$@")

  # Create merge branch
  git checkout -b "merged/$ORIGINAL" main

  # Merge each split
  for SPLIT in "${SPLITS[@]}"; do
    echo "Merging split: $SPLIT"

    if ! git merge --no-ff "split/$SPLIT"; then
      handle_merge_conflict "$SPLIT"
    fi

    checkpoint_merge "$SPLIT"
  done

  # Validate merged result
  validate_merged_splits "$ORIGINAL"
}
```

### 7. OUTPUT GENERATION
```bash
write_split_output() {
  cat > "$OUTPUT_FILE" <<EOF
{
  "status": "$STATUS",
  "original_effort": "$ORIGINAL_EFFORT",
  "splits_completed": $(json_array "${COMPLETED_SPLITS[@]}"),
  "line_counts": {
    $(generate_line_counts)
  },
  "all_within_limit": $ALL_WITHIN_LIMIT,
  "review_results": {
    $(generate_review_results)
  },
  "merge_status": "$MERGE_STATUS",
  "duration_seconds": $DURATION,
  "next_action": "$NEXT_ACTION"
}
EOF
}
```

## State Machine States
- `SPLIT_INIT` - Initialize split coordination
- `SPLIT_CREATE_BRANCHES` - Create split branches
- `SPLIT_SPAWN_ENGINEER` - Spawn SW engineer
- `SPLIT_MONITOR_IMPLEMENTATION` - Monitor progress
- `SPLIT_MEASURE_SIZE` - Measure and validate
- `SPLIT_SPAWN_REVIEWER` - Spawn code reviewer
- `SPLIT_MONITOR_REVIEW` - Monitor review
- `SPLIT_HANDLE_ISSUES` - Handle review issues
- `SPLIT_NEXT_SPLIT` - Move to next split
- `SPLIT_MERGE_SPLITS` - Merge all splits
- `SPLIT_COMPLETE` - Finalize

## Heartbeat Updates
```bash
update_split_heartbeat() {
  cat > "$HEARTBEAT_FILE" <<EOF
{
  "pid": $$,
  "status": "$STATUS",
  "progress_percentage": $PROGRESS,
  "current_state": "$CURRENT_STATE",
  "last_update": "$(date -Iseconds)",
  "current_split": "$CURRENT_SPLIT",
  "splits_completed": $COMPLETED_SPLITS,
  "splits_total": $TOTAL_SPLITS,
  "current_split_size": $CURRENT_SIZE,
  "within_limit": $WITHIN_LIMIT
}
EOF
}
```

## Recovery Support
```bash
# Checkpoint after each split
checkpoint_split() {
  local SPLIT="$1"
  local STATUS="$2"

  cat > "$CHECKPOINT_FILE" <<EOF
{
  "state": "SPLIT_COORDINATION",
  "data": {
    "completed_splits": $(json_array "${COMPLETED[@]}"),
    "current_split": "$SPLIT",
    "split_status": "$STATUS",
    "measurements": $(json_object "$MEASUREMENTS")
  },
  "can_resume": true
}
EOF
}

# Resume from checkpoint
resume_split_coordination() {
  local CHECKPOINT=$(cat "$CHECKPOINT_FILE")
  local LAST_SPLIT=$(echo "$CHECKPOINT" | jq -r '.data.current_split')

  # Find position in split array
  local INDEX=$(find_split_index "$LAST_SPLIT")

  # Resume from next split
  resume_from_index $((INDEX + 1))
}
```

## Error Handling
```bash
handle_split_error() {
  local ERROR_TYPE="$1"
  local SPLIT="$2"

  case "$ERROR_TYPE" in
    "SIZE_EXCEEDED")
      # Split is still too large
      create_sub_splits "$SPLIT"
      ;;
    "IMPLEMENTATION_FAILED")
      # SW engineer failed
      retry_split_implementation "$SPLIT"
      ;;
    "REVIEW_FAILED")
      # Review found issues
      spawn_fix_engineer "$SPLIT"
      ;;
    "MERGE_CONFLICT")
      # Splits conflict with each other
      resolve_split_conflicts
      ;;
  esac
}
```

## Rules Applied
- R377: Communication Protocol
- R378: Lifecycle Management
- R379: Monitoring and Recovery
- R220: Size Limit Enforcement
- R221: Continuous Delivery Requirements

## Success Criteria
- All splits completed successfully
- Each split within size limit
- Reviews passed (if required)
- Splits merged successfully
- Output file written
- Clean exit

## Common Issues
1. Split still exceeds size limit
2. Dependencies between splits complex
3. Merge conflicts between splits
4. Review finding issues late
5. Resource exhaustion with many splits

## Example Launch
```bash
# Prepare parameters
cat > /tmp/split-params.json <<'EOF'
{
  "sub_orchestrator_type": "SPLIT_COORDINATION",
  "unique_id": "split-e11-uuid",
  "input_parameters": {
    "original_effort": "E1.1",
    "split_plan": {
      "splits": ["E1.1a", "E1.1b"],
      "sequential": true
    },
    "size_limit": 700
  }
}
EOF

# Launch sub-orchestrator
/sub-orchestrate-split /tmp/split-params.json
```