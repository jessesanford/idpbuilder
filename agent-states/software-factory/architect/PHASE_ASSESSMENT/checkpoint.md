# PHASE_ASSESSMENT State Checkpoints

## When to Save Assessment State

### Critical Checkpoint Triggers

1. **Assessment Initiation** - When entering PHASE_ASSESSMENT state
   - Save: Phase scope, components to assess, entry timestamp
   - Reason: Establish assessment baseline and scope

2. **Midpoint Assessment** - After completing 50% of evaluation criteria
   - Save: Partial grades, identified issues, remaining work
   - Reason: Enable recovery if assessment is interrupted

3. **Pre-Decision State** - Before making final ON_TRACK/NEEDS_CORRECTION/OFF_TRACK decision  
   - Save: Complete assessment data, grade calculations, evidence
   - Reason: Support decision audit and potential appeals

4. **Decision Made** - After determining final assessment outcome
   - Save: Final decision, rationale, action items, next steps
   - **MANDATORY**: Create `PHASE-{N}-ASSESSMENT-REPORT.md` per R257
   - Location: `phase-assessments/phase{N}/`
   - Reason: Permanent record for phase transition approval (R257 requirement)

5. **Issue Escalation** - When critical architectural problems discovered
   - Save: Issue details, severity assessment, escalation rationale
   - Reason: Document decision to halt progression

## Required Decision Data to Preserve

### Assessment Context Data
```yaml
assessment_context:
  phase_id: "phase-1"
  phase_scope: "TMC Base Infrastructure"  
  waves_assessed: ["wave-1", "wave-2", "wave-3"]
  assessment_start: "2025-08-23T14:30:00Z"
  assessor_id: "architect-agent"
  integration_branch: "phase1/integration-main"
  previous_phase_baseline: "phase0-final-state.yaml"
```

### Technical Assessment Data
```yaml
technical_assessment:
  kcp_compliance:
    score: 100
    status: "PASS"
    violations: []
    evidence_artifacts:
      - "kcp-pattern-audit-report.md"
      - "logical-cluster-field-verification.yaml"
  
  api_quality:
    score: 88
    status: "NEEDS_MINOR_CORRECTION"
    issues:
      - "Inconsistent field naming in Config CRD"
      - "Missing validation for webhook timeout"
    evidence_artifacts:
      - "api-compatibility-matrix.yaml"
      - "schema-validation-results.json"
  
  integration_stability:
    score: 95
    status: "PASS"
    merge_conflicts_resolved: 2
    test_results:
      unit_tests: "98% pass rate"
      integration_tests: "100% pass rate"
      e2e_tests: "95% pass rate"
    evidence_artifacts:
      - "integration-test-report.html"
      - "merge-conflict-resolution-log.md"
  
  performance_benchmarks:
    score: 92
    status: "PASS"
    metrics:
      api_response_time: "145ms (target: <150ms)"
      memory_usage: "720MB (target: <768MB)" 
      throughput: "850 req/s (target: >800 req/s)"
      startup_time: "42s (target: <45s)"
    evidence_artifacts:
      - "performance-benchmark-report.json"
      - "load-test-results.html"
  
  security_posture:
    score: 100
    status: "PASS"
    vulnerabilities:
      critical: 0
      high: 0
      medium: 1
      low: 3
    evidence_artifacts:
      - "security-scan-report.json"
      - "rbac-audit-results.yaml"
```

### Decision Outcome Data
```yaml
assessment_decision:
  final_grade: 93.2
  outcome: "NEEDS_CORRECTION"
  rationale: "High overall quality with minor API consistency issues requiring correction"
  
  action_items:
    - id: "AI-001"
      priority: "HIGH"
      description: "Standardize field naming in Config CRD"
      assigned_to: "orchestrator"
      estimated_effort: "2 hours"
    
    - id: "AI-002"
      priority: "MEDIUM"
      description: "Add timeout validation to webhook configuration"
      assigned_to: "orchestrator"
      estimated_effort: "4 hours"
  
  next_steps:
    - "Return to orchestrator with action items"
    - "Re-assess after corrections implemented"
    - "Approve phase transition if corrections meet standards"
  
  escalation_required: false
  approval_authority: "architect-agent"
  decision_timestamp: "2025-08-23T16:45:00Z"
```

### Evidence Preservation
```yaml
evidence_preservation:
  assessment_artifacts:
    base_directory: "/assessments/phase-1-assessment-20250823"
    
  required_artifacts:
    - "phase-assessment-checklist.md"
    - "kcp-pattern-compliance-audit.yaml"  
    - "api-compatibility-analysis.json"
    - "integration-test-full-report.html"
    - "performance-benchmark-data.json"
    - "security-vulnerability-scan.json"
    - "architectural-decision-record.md"
    
  artifact_retention:
    duration: "2 years"
    backup_location: "s3://architecture-assessments/phase-1/"
    access_control: "architect-team-read-only"
```

## Recovery Protocol for Assessments

### State Recovery Procedure

1. **Detect Incomplete Assessment**
   ```bash
   # Check for checkpoint files
   find /checkpoints -name "phase-assessment-*.yaml" -mtime -1
   
   # Verify assessment completeness
   grep -q "assessment_complete: true" latest_checkpoint.yaml
   ```

2. **Load Assessment Context**
   ```yaml
   # From checkpoint file
   current_progress:
     completed_criteria: ["kcp_compliance", "api_quality"]
     remaining_criteria: ["integration_stability", "performance", "security"]
     partial_scores:
       kcp_compliance: 100
       api_quality: 88
     evidence_collected: 
       - "/evidence/kcp-audit-20250823.yaml"
       - "/evidence/api-analysis-20250823.json"
   ```

3. **Resume Assessment Process**
   - Continue from last completed criterion
   - Validate existing evidence artifacts still accessible
   - Re-run any time-sensitive measurements (performance benchmarks)
   - Preserve continuity of assessment methodology

4. **Validate Recovery Integrity**
   ```yaml
   recovery_validation:
     checkpoint_data_integrity: "verified"
     evidence_artifacts_accessible: "verified" 
     assessment_methodology_consistent: "verified"
     timeline_continuity: "acceptable_gap_documented"
   ```

### Assessment Interruption Handling

#### Planned Interruptions
- **Before Interruption**: Save complete state to checkpoint file
- **During Interruption**: Preserve all evidence artifacts
- **After Resume**: Validate checkpoint integrity and continue

#### Unplanned Interruptions  
- **Detection**: Missing completion marker in checkpoint
- **Recovery**: Load last valid checkpoint, verify evidence artifacts
- **Continuation**: Resume from last completed assessment criterion
- **Validation**: Cross-check partial results for consistency

### Checkpoint File Management

#### Checkpoint File Naming Convention
```
phase-assessment-{phase-id}-{timestamp}-{status}.yaml

Examples:
- phase-assessment-phase1-20250823-143000-initiated.yaml
- phase-assessment-phase1-20250823-153000-midpoint.yaml  
- phase-assessment-phase1-20250823-164500-complete.yaml
```

#### Checkpoint File Locations
- **Active Checkpoints**: `/workspaces/architect-checkpoints/active/`
- **Completed Assessments**: `/workspaces/architect-checkpoints/completed/`
- **Backup Location**: `/workspaces/architect-checkpoints/backup/`

#### Checkpoint Retention Policy
- **Active**: Keep until assessment complete
- **Completed**: Keep for 30 days after phase approval
- **Backup**: Keep for 1 year for audit purposes
- **Archive**: Long-term storage after project completion

### Recovery Validation Checklist

Before resuming assessment after interruption:

- [ ] Checkpoint file integrity verified
- [ ] All evidence artifacts accessible  
- [ ] Assessment timeline documented
- [ ] Partial results validated for consistency
- [ ] Assessment criteria methodology unchanged
- [ ] Required tools and access still available
- [ ] Phase integration branch state unchanged
- [ ] Previous assessment decisions still valid

### Emergency Recovery Procedures

If checkpoint data is corrupted or unavailable:

1. **Reconstruct Assessment State**
   - Review orchestrator state for phase context
   - Re-examine integration branch for current state
   - Recreate evidence artifacts from source data
   - Document assessment restart rationale

2. **Restart Assessment Protocol**
   - Begin fresh assessment with full documentation
   - Note previous assessment attempt in records  
   - Apply same standards and methodology
   - Preserve any valid prior evidence

3. **Escalation Triggers**
   - Multiple recovery failures indicate systemic issues
   - Evidence artifacts permanently lost
   - Assessment criteria changed during interruption
   - Phase state significantly altered during interruption