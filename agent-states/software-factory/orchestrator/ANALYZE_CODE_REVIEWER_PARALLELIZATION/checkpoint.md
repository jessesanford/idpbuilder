# ANALYZE_CODE_REVIEWER_PARALLELIZATION - Checkpoint Template

## State Entry
```yaml
checkpoint:
  state: "ANALYZE_CODE_REVIEWER_PARALLELIZATION"
  timestamp: "{ISO-8601}"
  phase: {PHASE}
  wave: {WAVE}
  reason: "Analyzing parallelization before spawning Code Reviewers"
```

## Required Inputs
```yaml
inputs:
  wave_plan_location: "phase-plans/PHASE-{PHASE}-WAVE-{WAVE}-IMPLEMENTATION-PLAN.md"
  current_efforts: []  # List of efforts in wave
  infrastructure_ready: true  # From CREATE_NEXT_INFRASTRUCTURE
```

## Analysis Progress
```yaml
analysis:
  wave_plan_read: false  # Set true after Read tool
  blocking_efforts_identified: []
  parallel_efforts_identified: []
  dependencies_mapped: {}
  spawn_sequence_created: false
```

## Parallelization Plan
```yaml
code_reviewer_parallelization_plan:
  wave: {WAVE}
  phase: {PHASE}
  analysis_timestamp: "{ISO-8601}"
  
  blocking_efforts:
    - effort_id: ""
      name: ""
      reason: ""
      can_parallelize: false
      
  parallel_groups:
    group_1:
      can_start_after: []
      efforts: []
      
  spawn_sequence:
    - step: 1
      action: ""
      efforts: []
      wait_for_completion: true/false
```

## Validation Checklist
```yaml
validation:
  wave_plan_exists: false
  wave_plan_read: false
  blocking_efforts_found: false
  parallel_groups_defined: false
  spawn_sequence_complete: false
  state_file_updated: false
  acknowledgment_output: false
```

## State Exit
```yaml
exit:
  next_state: "SPAWN_CODE_REVIEWERS_EFFORT_PLANNING"
  parallelization_strategy: "defined"
  blocking_count: 0
  parallel_count: 0
  total_spawn_steps: 0
  validation_passed: true/false
```

## Recovery Information
```yaml
recovery:
  last_action: ""
  error_message: ""
  retry_possible: true/false
  fallback_state: "CREATE_NEXT_INFRASTRUCTURE"
```