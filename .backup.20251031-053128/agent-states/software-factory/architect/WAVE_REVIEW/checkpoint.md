# REVIEW_WAVE_ARCHITECTURE State Checkpoints

## When to Save Wave Review State

### Critical Checkpoint Triggers

1. **Wave Review Initiation** - When entering REVIEW_WAVE_ARCHITECTURE state
   - Save: Wave scope, effort list, assessment baseline
   - Reason: Establish review boundary and effort inventory

2. **Effort Assessment Completion** - After completing each effort evaluation
   - Save: Individual effort grades, issues found, evidence artifacts
   - Reason: Preserve individual assessments for wave score calculation

3. **Integration Testing Complete** - After wave integration testing finishes
   - Save: Integration test results, merge status, performance benchmarks
   - Reason: Critical data for integration stability scoring

4. **Pre-Decision State** - Before making PROCEED/CHANGES_REQUIRED/STOP decision
   - Save: Complete wave assessment, calculated scores, evidence summary
   - Reason: Support decision audit and potential reconsideration

5. **Decision Made** - After determining final wave review outcome
   - Save: Decision rationale, action items, next steps, approval record
   - Reason: Permanent record for wave integration approval/rejection

## Required Decision Data to Preserve

### Wave Assessment Context
```yaml
wave_review_context:
  wave_id: "phase-1-wave-2"
  wave_scope: "Webhook Framework and Configuration Management"
  efforts_in_wave:
    - effort_id: "webhook-core"
      branch: "phase1/wave2/effort-webhook-core"
      line_count: 650
      status: "ASSESSED"
    - effort_id: "webhook-validation"  
      branch: "phase1/wave2/effort-webhook-validation"
      line_count: 720
      status: "ASSESSED"
    - effort_id: "config-management"
      branch: "phase1/wave2/effort-config-management"
      line_count: 580
      status: "ASSESSED"
  
  wave_integration_branch: "phase1/wave2/integration-main"
  review_start: "2025-08-23T09:30:00Z"
  reviewer_id: "architect-agent"
  parent_phase: "phase-1"
```

### Individual Effort Assessments
```yaml
effort_assessments:
  webhook-core:
    size_compliance:
      line_count: 650
      measurement_tool: "line-counter.sh"
      status: "PASS"
      evidence: "line-count-webhook-core-20250823.txt"
    
    kcp_patterns:
      score: 95
      issues:
        - "Minor: Controller indexing could be more efficient"
      violations: []
      evidence: "kcp-pattern-audit-webhook-core.yaml"
    
    code_quality:
      test_coverage: 92
      documentation: "Complete work log"
      technical_debt: "Low"
    
    integration_readiness:
      merge_status: "clean"
      dependencies_resolved: true
      api_compatibility: "maintained"

  webhook-validation:
    size_compliance:
      line_count: 720
      measurement_tool: "line-counter.sh"
      status: "PASS"
      evidence: "line-count-webhook-validation-20250823.txt"
    
    kcp_patterns:
      score: 88
      issues:
        - "LogicalCluster field missing in ValidationRule CRD"
      violations: ["Missing LogicalCluster in CRD schema"]
      evidence: "kcp-pattern-audit-webhook-validation.yaml"
    
    code_quality:
      test_coverage: 89
      documentation: "Work log incomplete - missing validation decisions"
      technical_debt: "Medium"
    
    integration_readiness:
      merge_status: "conflicts_resolved"
      dependencies_resolved: true
      api_compatibility: "breaking_change_documented"

  config-management:
    size_compliance:
      line_count: 580
      measurement_tool: "line-counter.sh"
      status: "PASS"
      evidence: "line-count-config-management-20250823.txt"
    
    kcp_patterns:
      score: 98
      issues: []
      violations: []
      evidence: "kcp-pattern-audit-config-management.yaml"
    
    code_quality:
      test_coverage: 95
      documentation: "Excellent work log with all decisions documented"
      technical_debt: "None"
    
    integration_readiness:
      merge_status: "clean"
      dependencies_resolved: true
      api_compatibility: "maintained"
```

### Wave Integration Assessment
```yaml
wave_integration_assessment:
  merge_compatibility:
    status: "PROJECT_DONE"
    conflicts_found: 2
    conflicts_resolved: 2
    merge_time: "2025-08-23T13:45:00Z"
    evidence: "wave2-integration-merge-log.txt"
  
  integration_testing:
    test_suite: "wave2-integration-tests"
    total_tests: 42
    passed: 40
    failed: 2
    skipped: 0
    pass_rate: 95.2
    failed_tests:
      - "TestWebhookValidationIntegration"
      - "TestConfigWebhookAuth" 
    evidence: "wave2-integration-test-results.xml"
  
  performance_testing:
    baseline_metrics:
      api_response_time: "98ms"
      memory_usage: "485MB"
      cpu_usage: "1.8 cores"
      throughput: "1050 req/s"
    
    wave_metrics:
      api_response_time: "112ms"  # +14ms increase
      memory_usage: "545MB"       # +60MB increase
      cpu_usage: "2.1 cores"      # +0.3 cores increase
      throughput: "980 req/s"     # -70 req/s decrease
    
    performance_impact:
      response_time_change: "+14.3%"
      memory_change: "+12.4%"
      cpu_change: "+16.7%"
      throughput_change: "-6.7%"
    
    assessment: "WARNING - Response time and CPU usage above 10% threshold"
    evidence: "wave2-performance-benchmark-results.json"
```

### Wave Grade Calculation
```yaml
wave_grading:
  size_compliance:
    score: 100
    rationale: "All 3 efforts under 800 line limit"
    evidence: ["line-count-webhook-core-20250823.txt", "line-count-webhook-validation-20250823.txt", "line-count-config-management-20250823.txt"]
  
  kcp_pattern_consistency:
    individual_scores: [95, 88, 98]
    average_score: 93.7
    rationale: "One LogicalCluster field missing in webhook-validation"
    evidence: ["kcp-pattern-audit-webhook-core.yaml", "kcp-pattern-audit-webhook-validation.yaml", "kcp-pattern-audit-config-management.yaml"]
  
  integration_stability:
    score: 87
    rationale: "Clean merges but 2 integration test failures"
    components:
      merge_compatibility: 95  # -5 for initial conflicts
      integration_tests: 90    # -10 for 2 failed tests
      cross_effort_apis: 85    # -15 for auth integration issue
    evidence: ["wave2-integration-merge-log.txt", "wave2-integration-test-results.xml"]
  
  api_coherence:
    score: 88
    rationale: "Minor naming inconsistencies, one breaking change documented"
    components:
      naming_consistency: 92
      interface_compatibility: 82  # webhook-validation breaking change
      schema_alignment: 90
    evidence: ["api-coherence-analysis-wave2.yaml"]
  
  performance_impact:
    score: 75
    rationale: "Response time and CPU usage exceed 10% threshold but under 20%"
    components:
      response_time: 65  # +14.3% over threshold
      memory_usage: 88   # +12.4% acceptable
      cpu_usage: 67      # +16.7% over threshold  
      throughput: 80     # -6.7% acceptable
    evidence: ["wave2-performance-benchmark-results.json"]
  
  documentation_quality:
    score: 85
    rationale: "One effort has incomplete work log"
    components:
      work_log_completeness: 67  # webhook-validation incomplete
      implementation_notes: 90
      integration_guide: 100
    evidence: ["work-log-review-summary.md"]
  
  final_calculation:
    formula: "(100×0.25) + (93.7×0.25) + (87×0.20) + (88×0.15) + (75×0.10) + (85×0.05)"
    result: "25 + 23.42 + 17.4 + 13.2 + 7.5 + 4.25 = 90.77%"
    final_score: 90.8
```

### Decision Outcome Record
```yaml
wave_review_decision:
  final_score: 90.8
  decision: "CHANGES_REQUIRED"
  rationale: "Good wave quality but performance regression and missing LogicalCluster field require correction"
  
  required_changes:
    - change_id: "WR-001"
      priority: "HIGH"
      description: "Add LogicalCluster field to ValidationRule CRD in webhook-validation effort"
      effort_affected: "webhook-validation"
      estimated_effort: "2 hours"
      verification: "Re-run KCP pattern audit"
    
    - change_id: "WR-002"  
      priority: "HIGH"
      description: "Optimize webhook validation logic causing 14% response time increase"
      effort_affected: "webhook-core"
      estimated_effort: "6 hours"
      verification: "Re-run performance benchmarks"
    
    - change_id: "WR-003"
      priority: "MEDIUM"
      description: "Fix 2 failing integration tests: TestWebhookValidationIntegration, TestConfigWebhookAuth"
      effort_affected: "webhook-validation"
      estimated_effort: "4 hours"
      verification: "Re-run integration test suite"
    
    - change_id: "WR-004"
      priority: "LOW"
      description: "Complete work log documentation for webhook-validation effort"
      effort_affected: "webhook-validation"
      estimated_effort: "1 hour"
      verification: "Work log review"
  
  acceptance_criteria:
    - "KCP pattern score >95% (add LogicalCluster field)"
    - "Performance impact <10% (optimize response time)"
    - "Integration test pass rate >98% (fix failing tests)"
    - "All work logs complete"
  
  resubmission_timeline: "3 days maximum"
  next_review_trigger: "All required changes implemented and verified"
  
  approver: "architect-agent"
  decision_timestamp: "2025-08-23T15:30:00Z"
  decision_valid_until: "2025-08-30T15:30:00Z"
```

## Recovery Protocol for Wave Reviews

### State Recovery Procedure

1. **Detect Incomplete Wave Review**
   ```bash
   # Check for wave review checkpoint files
   find /checkpoints -name "wave-review-*.yaml" -mtime -7
   
   # Verify review completion status
   grep -q "review_complete: true" latest_wave_checkpoint.yaml
   ```

2. **Load Wave Review Context**
   ```yaml
   # From checkpoint file
   recovery_context:
     wave_id: "phase-1-wave-2"
     efforts_assessed: ["webhook-core", "webhook-validation"]
     efforts_remaining: ["config-management"]
     completed_phases:
       - "individual_effort_assessment"
       - "size_compliance_verification"
     remaining_phases:
       - "integration_testing"
       - "performance_benchmarking"
       - "final_grading"
   ```

3. **Resume Wave Review Process**
   - Continue from last completed assessment phase
   - Validate existing evidence artifacts are accessible
   - Re-run time-sensitive tests (performance benchmarks)
   - Maintain assessment methodology consistency

4. **Validate Recovery Integrity**
   ```yaml
   recovery_validation:
     checkpoint_integrity: "verified"
     evidence_artifacts: "accessible"
     assessment_consistency: "maintained"
     timeline_impact: "documented"
   ```

### Wave Review Interruption Handling

#### Planned Interruptions
- **Before Pause**: Save complete state with current assessment progress
- **During Pause**: Preserve all test results and evidence artifacts
- **After Resume**: Validate checkpoint and continue from last phase

#### Unplanned Interruptions
- **Detection**: Incomplete checkpoint without completion marker
- **Recovery**: Load most recent valid checkpoint state
- **Continuation**: Resume from last completed assessment phase
- **Validation**: Cross-check partial results for consistency

### Checkpoint File Management

#### Checkpoint File Naming Convention
```
wave-review-{wave-id}-{timestamp}-{status}.yaml

Examples:
- wave-review-phase1-wave2-20250823-093000-initiated.yaml
- wave-review-phase1-wave2-20250823-123000-efforts-assessed.yaml
- wave-review-phase1-wave2-20250823-143000-integration-tested.yaml
- wave-review-phase1-wave2-20250823-153000-complete.yaml
```

#### Checkpoint File Locations
- **Active Reviews**: `/workspaces/architect-checkpoints/wave-reviews/active/`
- **Completed Reviews**: `/workspaces/architect-checkpoints/wave-reviews/completed/`
- **Evidence Archive**: `/workspaces/architect-checkpoints/wave-reviews/evidence/`

#### Checkpoint Retention Policy
- **Active**: Keep until wave review complete
- **Completed**: Keep for 60 days after wave integration
- **Evidence**: Keep for 6 months for audit purposes
- **Archive**: Long-term storage for project retrospectives

### Recovery Validation Checklist

Before resuming wave review after interruption:

- [ ] Checkpoint file integrity verified
- [ ] All evidence artifacts accessible and valid
- [ ] Wave integration branch state unchanged
- [ ] Individual effort branches still accessible
- [ ] Assessment methodology documentation current
- [ ] Testing tools and environments available
- [ ] Performance baseline data still valid
- [ ] Previous individual effort assessments consistent

### Emergency Recovery Procedures

If checkpoint data is corrupted or lost:

1. **Reconstruct Wave Review State**
   - Examine orchestrator state for wave context
   - Re-identify efforts in wave scope
   - Recreate evidence artifacts from source data
   - Document review restart rationale

2. **Fresh Wave Review Protocol**
   - Begin new assessment with full methodology
   - Document previous attempt in review record
   - Apply consistent standards across all efforts
   - Preserve any valid prior evidence artifacts

3. **Escalation Triggers**
   - Multiple recovery failures indicate systemic issues
   - Evidence artifacts permanently corrupted or lost
   - Wave integration branch significantly altered
   - Assessment methodology changed during interruption
   - Performance baseline no longer valid