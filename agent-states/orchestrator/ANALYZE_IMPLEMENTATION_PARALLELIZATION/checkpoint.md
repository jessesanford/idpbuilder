# ANALYZE_IMPLEMENTATION_PARALLELIZATION - Checkpoint Template

## State Entry
```yaml
checkpoint:
  state: "ANALYZE_IMPLEMENTATION_PARALLELIZATION"
  timestamp: "{ISO-8601}"
  phase: {PHASE}
  wave: {WAVE}
  reason: "Analyzing implementation plans for SW Engineer parallelization"
```

## Required Inputs
```yaml
inputs:
  effort_directories: []  # List of effort directories
  implementation_plans: []  # Paths to IMPLEMENTATION-PLAN.md files
  code_reviewer_parallelization: {}  # From previous analysis
  efforts_ready_count: 0  # Number of completed effort plans
```

## Analysis Progress
```yaml
analysis:
  plans_read: {}  # effort_name: read_status
  metadata_extracted: {}  # effort_name: metadata
  blocking_implementations: []
  parallel_implementations: []
  dependencies_verified: false
  consistency_checked: false
```

## SW Engineer Parallelization Plan
```yaml
sw_engineer_parallelization_plan:
  wave: {WAVE}
  phase: {PHASE}
  analysis_timestamp: "{ISO-8601}"
  
  blocking_implementations:
    - effort_id: ""
      name: ""
      implementation_plan: ""
      can_parallelize: false
      reason: ""
      dependencies: []
      
  parallel_groups:
    group_1:
      can_start_after: []
      efforts:
        - effort_id: ""
          name: ""
          implementation_plan: ""
          
  spawn_sequence:
    - step: 1
      action: ""
      agent_type: "sw-engineer"
      efforts: []
      wait_for_completion: true/false
      expected_duration: ""
```

## Consistency Verification
```yaml
consistency:
  wave_plan_match: true/false
  metadata_preserved: true/false
  conflicts_found: []
  resolution_applied: ""
```

## Validation Checklist
```yaml
validation:
  all_plans_exist: false
  all_plans_read: false
  blocking_identified: false
  parallel_groups_defined: false
  spawn_sequence_complete: false
  consistency_verified: false
  state_file_updated: false
  acknowledgment_output: false
  directories_verified: false
```

## State Exit
```yaml
exit:
  next_state: "SPAWN_AGENTS"
  parallelization_strategy: "defined"
  blocking_count: 0
  parallel_count: 0
  total_spawn_steps: 0
  validation_passed: true/false
```

## Recovery Information
```yaml
recovery:
  last_effort_analyzed: ""
  error_message: ""
  retry_possible: true/false
  fallback_state: "WAITING_FOR_EFFORT_PLANS"
```