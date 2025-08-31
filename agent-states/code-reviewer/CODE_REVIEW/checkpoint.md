# Code Reviewer - CODE_REVIEW State Checkpoint

## When to Save State

Save checkpoint at these critical code review points:

1. **Initial Review Assessment**
   - Implementation received for review
   - Initial size compliance check completed
   - Review scope and criteria established

2. **Size Compliance Validation**
   - Size measured using tmc-pr-line-counter.sh
   - Compliance decision made
   - Split determination if size exceeded

3. **Technical Review Complete**
   - Code quality assessment finished
   - Architecture compliance validated
   - KCP pattern compliance checked

4. **Test Coverage Validation**
   - Test coverage measured and assessed
   - Coverage requirements validation complete
   - Multi-tenant test scenarios reviewed

5. **Final Review Decision**
   - All review areas completed
   - Final approval/rejection decision made
   - Feedback and recommendations prepared

## Required Data to Preserve

```yaml
code_review_checkpoint:
  # State identification
  state: "CODE_REVIEW"
  effort_id: "phase1-wave2-effort3-webhooks"
  phase: 1
  wave: 2
  checkpoint_timestamp: "2025-08-23T19:30:00Z"
  
  # Review context
  review_session:
    started_at: "2025-08-23T18:45:00Z"
    duration_minutes: 45
    implementation_location: "/workspaces/efforts/phase1/wave2/effort3-webhooks"
    implementation_branch: "phase1/wave2/effort3-webhooks"
    original_plan_location: "/workspaces/efforts/phase1/wave2/effort3-webhooks/IMPLEMENTATION-PLAN.md"
    assigned_sw_engineer: "sw-engineer-agent-001"
    
  # Implementation metadata
  implementation_info:
    completed_at: "2025-08-23T18:30:00Z"
    implementation_duration_hours: 8.5
    commits_count: 23
    files_changed: 15
    sw_engineer_reported_complete: true
    work_log_location: "/workspaces/efforts/phase1/wave2/effort3-webhooks/work-log.md"
    
  # Size compliance validation (CRITICAL)
  size_compliance:
    validation_timestamp: "2025-08-23T18:50:00Z"
    tool_used: "tmc-pr-line-counter.sh"
    command_executed: "/workspaces/kcp-shared-tools/tmc-pr-line-counter.sh -c phase1/wave2/effort3-webhooks"
    raw_output: |
      pkg/webhooks/admission/validator.go: 127 lines
      pkg/webhooks/admission/mutator.go: 104 lines
      pkg/webhooks/admission/types.go: 52 lines
      pkg/webhooks/server/server.go: 156 lines
      pkg/webhooks/server/routing.go: 98 lines
      pkg/webhooks/server/handler.go: 47 lines
      pkg/webhooks/config/configuration.go: 73 lines
      pkg/webhooks/config/loader.go: 34 lines
      Total: 691
    
    actual_lines: 691
    size_limit: 800
    compliant: true
    margin: 109
    size_risk_level: "LOW"
    
    # Detailed breakdown by component
    component_sizes:
      admission_logic: 283  # validator + mutator + types
      server_framework: 301  # server + routing + handler
      configuration: 107   # configuration + loader
      
    size_validation_passed: true
    split_required: false
    
  # Test coverage validation
  test_coverage_validation:
    validation_timestamp: "2025-08-23T19:00:00Z"
    
    unit_test_coverage:
      tool_used: "go test -coverprofile"
      overall_coverage: 92.3
      coverage_by_package:
        pkg/webhooks/admission: 94.7
        pkg/webhooks/server: 91.2
        pkg/webhooks/config: 89.8
      meets_90_percent_requirement: true
      
      detailed_coverage:
        total_lines: 691
        covered_lines: 638
        uncovered_lines: 53
        test_files: 8
        test_functions: 47
        
    integration_test_coverage:
      integration_tests_present: true
      test_scenarios_count: 12
      api_endpoints_covered: 6
      error_scenarios_covered: 8
      
      key_scenarios:
        - "Validating webhook accepts valid resources"
        - "Validating webhook rejects invalid resources"
        - "Mutating webhook applies default values"
        - "Multi-tenant webhook routing validation"
        - "Webhook server startup and health checks"
        - "Configuration loading and validation"
        
    multi_tenant_test_coverage:
      multi_tenant_scenarios: 6
      logical_cluster_tests: 4
      workspace_isolation_tests: 3
      cross_workspace_security_tests: 2
      
      scenarios_details:
        - scenario: "Webhook routing by logical cluster"
          test_file: "test/webhooks/routing_test.go"
          coverage: "ADEQUATE"
        - scenario: "Workspace isolation enforcement"
          test_file: "test/webhooks/isolation_test.go" 
          coverage: "GOOD"
        - scenario: "Cross-workspace access blocking"
          test_file: "test/webhooks/security_test.go"
          coverage: "ADEQUATE"
          
    performance_test_coverage:
      performance_tests_present: true
      latency_tests: 3
      load_tests: 2
      memory_usage_tests: 1
      
      performance_targets_validated:
        webhook_latency: "<100ms - PASS"
        memory_usage: "<50MB - PASS" 
        cpu_usage: "<10% - PASS"
        
    overall_test_assessment:
      coverage_score: 94
      meets_requirements: true
      test_quality: "HIGH"
      multi_tenant_coverage: "ADEQUATE"
      
  # KCP/Kubernetes pattern compliance
  kcp_compliance_validation:
    validation_timestamp: "2025-08-23T19:10:00Z"
    
    multi_tenancy_compliance:
      logical_cluster_awareness: "IMPLEMENTED"
      workspace_isolation: "PROPERLY_ENFORCED"
      tenant_routing: "CORRECTLY_IMPLEMENTED"
      
      findings:
        - "WebhookServer properly extracts logical cluster from request context"
        - "TenantRouter correctly isolates webhook handling by workspace"
        - "No cross-tenant data leakage detected in implementation"
        
      compliance_score: 95
      
    api_export_integration:
      api_export_aware: "IMPLEMENTED"
      export_discovery: "CORRECTLY_INTEGRATED"
      export_lifecycle: "PROPERLY_HANDLED"
      
      findings:
        - "WebhookConfiguration registers via APIExports"
        - "Webhook discovery uses APIExport watcher"
        - "Export lifecycle changes trigger webhook reconfiguration"
        
      compliance_score: 90
      
    syncer_compatibility:
      sync_behavior: "COMPATIBLE"
      conflict_resolution: "PROPERLY_HANDLED"
      performance_impact: "MINIMAL"
      
      findings:
        - "Webhook configurations sync correctly to physical clusters"
        - "No syncer conflicts detected in webhook resources"
        - "Minimal impact on syncer performance"
        
      compliance_score: 88
      
    rbac_implementation:
      workspace_scoped_permissions: "IMPLEMENTED"
      least_privilege: "FOLLOWED"
      resource_access_control: "PROPER"
      
      findings:
        - "Webhook server uses workspace-scoped ServiceAccount"
        - "Permissions limited to required resources only"
        - "No cluster-wide permissions requested unnecessarily"
        
      compliance_score: 92
      
    overall_kcp_compliance:
      average_compliance: 91.3
      critical_issues: 0
      high_issues: 1  # Minor syncer optimization opportunity
      medium_issues: 2
      compliance_level: "HIGH"
      
  # Architecture review
  architecture_review:
    validation_timestamp: "2025-08-23T19:15:00Z"
    
    plan_adherence:
      implementation_follows_plan: true
      planned_components_implemented: 6
      additional_components_added: 1  # Extra utility component
      missing_planned_components: 0
      
      component_analysis:
        webhook_server:
          planned: "Core webhook HTTP server with KCP routing"
          implemented: "WebhookServer with KCP-aware request handling"
          adherence: "EXCELLENT"
          
        admission_controllers:
          planned: "Validating and mutating admission logic"
          implemented: "ValidatingController and MutatingController with tenant awareness"
          adherence: "EXCELLENT"
          
        routing_logic:
          planned: "KCP-aware request routing and tenant isolation"
          implemented: "TenantRouter with workspace isolation"
          adherence: "GOOD"
          
        configuration:
          planned: "Webhook configuration and management"
          implemented: "ConfigurationLoader with dynamic updates"
          adherence: "GOOD"
          
      plan_adherence_score: 89
      
    design_pattern_usage:
      controller_pattern: "PROPERLY_IMPLEMENTED"
      factory_pattern: "USED_APPROPRIATELY"
      strategy_pattern: "USED_FOR_ROUTING"
      observer_pattern: "USED_FOR_CONFIGURATION"
      
      pattern_quality_score: 92
      
    interface_compliance:
      planned_interfaces_implemented: 8
      interface_segregation: "GOOD"
      dependency_injection: "PROPERLY_USED"
      
      key_interfaces:
        - "WebhookServer interface - IMPLEMENTED"
        - "AdmissionController interface - IMPLEMENTED" 
        - "TenantRouter interface - IMPLEMENTED"
        - "ConfigurationLoader interface - IMPLEMENTED"
        
      interface_quality_score: 94
      
    component_structure:
      separation_of_concerns: "GOOD"
      cohesion: "HIGH"
      coupling: "LOW"
      modularity: "EXCELLENT"
      
      structure_quality_score: 90
      
    overall_architecture_score: 91
    
  # Security review
  security_review:
    validation_timestamp: "2025-08-23T19:20:00Z"
    
    input_validation:
      user_input_validated: true
      request_sanitization: "IMPLEMENTED"
      parameter_validation: "COMPREHENSIVE"
      
      findings:
        - "All webhook requests validated before processing"
        - "Resource specifications sanitized appropriately"
        - "Configuration parameters validated on load"
        
      input_validation_score: 94
      
    workspace_isolation_security:
      tenant_data_isolation: "PROPERLY_ENFORCED"
      cross_tenant_access_prevention: "IMPLEMENTED"
      workspace_boundary_respect: "CORRECT"
      
      findings:
        - "No cross-workspace data access detected"
        - "Workspace context properly validated"
        - "Tenant boundaries respected in all operations"
        
      isolation_security_score: 96
      
    rbac_security:
      least_privilege_followed: true
      permission_scope_appropriate: true
      service_account_security: "PROPER"
      
      findings:
        - "ServiceAccount permissions minimal and scoped"
        - "No unnecessary cluster-wide permissions"
        - "Role bindings properly workspace-scoped"
        
      rbac_security_score: 93
      
    secret_management:
      no_hardcoded_secrets: true
      environment_variable_usage: "APPROPRIATE"
      secret_rotation_support: "IMPLEMENTED"
      
      findings:
        - "No hardcoded credentials detected"
        - "TLS certificates loaded from environment"
        - "Webhook secrets properly managed"
        
      secret_management_score: 95
      
    error_handling_security:
      no_information_disclosure: true
      error_messages_safe: true
      logging_security: "APPROPRIATE"
      
      findings:
        - "Error messages don't leak sensitive information"
        - "Stack traces not exposed to clients"
        - "Logging excludes sensitive data"
        
      error_handling_score: 92
      
    overall_security_score: 94
    critical_security_issues: 0
    high_security_issues: 0
    medium_security_issues: 1  # Suggest additional rate limiting
    
  # Performance review
  performance_review:
    validation_timestamp: "2025-08-23T19:25:00Z"
    
    latency_performance:
      webhook_response_time: "78ms avg"
      meets_100ms_requirement: true
      p95_latency: "94ms"
      p99_latency: "127ms"
      
      performance_grade: "EXCELLENT"
      
    resource_usage:
      memory_usage_steady_state: "38MB"
      memory_usage_under_load: "45MB" 
      meets_50mb_requirement: true
      
      cpu_usage_idle: "2%"
      cpu_usage_normal_load: "8%"
      meets_10_percent_requirement: true
      
      resource_grade: "EXCELLENT"
      
    scalability:
      concurrent_request_handling: "GOOD"
      horizontal_scaling_ready: true
      resource_scaling_behavior: "LINEAR"
      
      scalability_grade: "GOOD"
      
    overall_performance_score: 93
    performance_issues: 0
    
  # Final review decision
  review_decision:
    decision_timestamp: "2025-08-23T19:30:00Z"
    
    decision_factors:
      size_compliant: true
      test_coverage_adequate: true
      kcp_compliance_high: true
      architecture_sound: true
      security_acceptable: true
      performance_meets_requirements: true
      
    blocking_issues: []
    
    non_blocking_issues:
      - type: "OPTIMIZATION"
        severity: "LOW"
        description: "Consider adding rate limiting for webhook requests"
        recommendation: "Add configurable rate limiting to webhook server"
        
      - type: "ENHANCEMENT"
        severity: "LOW"
        description: "Additional error handling tests could improve coverage"
        recommendation: "Add more edge case error scenarios to test suite"
        
    warnings:
      - "Performance tests could cover more concurrent request scenarios"
      - "Consider adding metrics collection for webhook operations"
      
    recommendations:
      - "Add comprehensive logging for debugging multi-tenant webhook routing"
      - "Consider implementing webhook request caching for frequently accessed configurations"
      - "Document webhook configuration patterns for other teams"
      
    final_decision:
      result: "APPROVED"
      can_proceed: true
      confidence_level: 95
      
    approval_conditions: []
    
    next_steps:
      - "Mark effort as complete in orchestrator state"
      - "Update work log with review results"
      - "Notify orchestrator of approval for integration planning"
      
  # Review quality metrics
  review_quality_assessment:
    review_thoroughness: 96
    issue_identification_effectiveness: 94
    decision_appropriateness: 95
    recommendation_quality: 89
    
    areas_reviewed:
      - "Size compliance (mandatory tool usage)"
      - "Test coverage (unit, integration, multi-tenant)"
      - "KCP pattern compliance (multi-tenancy, APIExport, syncer)"
      - "Architecture adherence to implementation plan"
      - "Security compliance (isolation, RBAC, secrets)"
      - "Performance validation (latency, resources, scalability)"
      
    review_completeness: "COMPREHENSIVE"
    reviewer_confidence: "HIGH"
    
  # Post-review tracking
  post_review_tracking:
    sw_engineer_feedback_requested: true
    orchestrator_notification_sent: true
    
    follow_up_actions:
      - action: "Monitor implementation integration into wave"
        priority: "MEDIUM"
        estimated_timeline: "Next wave integration"
        
      - action: "Validate recommendations are considered in future efforts"
        priority: "LOW"
        estimated_timeline: "Next similar effort"
        
    success_metrics_to_track:
      - "No critical issues discovered in integration testing"
      - "Performance targets maintained in production"
      - "Multi-tenancy works correctly in full system"
```

## Recovery Protocol

### Context Recovery After Interruption

```python
def recover_code_review_state(checkpoint_data):
    """Recover code review state from checkpoint"""
    
    print("🔄 RECOVERING CODE_REVIEW STATE")
    
    effort_info = checkpoint_data['effort_id']
    review_status = checkpoint_data.get('review_decision', {})
    
    print(f"Effort: {effort_info}")
    print(f"Review Decision: {review_status.get('result', 'IN_PROGRESS')}")
    print(f"Can Proceed: {review_status.get('can_proceed', False)}")
    
    # Verify implementation still exists and unchanged
    implementation_verification = verify_implementation_unchanged(checkpoint_data)
    
    # Check if any review findings need revalidation
    revalidation_needed = assess_revalidation_needs(checkpoint_data)
    
    # Determine recovery actions needed
    recovery_actions = determine_review_recovery_actions(
        checkpoint_data, implementation_verification, revalidation_needed
    )
    
    return {
        'effort_id': effort_info,
        'review_decision': review_status.get('result', 'IN_PROGRESS'),
        'implementation_verification': implementation_verification,
        'revalidation_needed': revalidation_needed,
        'recovery_actions': recovery_actions,
        'review_complete': review_status.get('result') in ['APPROVED', 'CHANGES_REQUIRED'],
        'recovery_needed': len(recovery_actions) > 0
    }

def verify_implementation_unchanged(checkpoint_data):
    """Verify implementation hasn't changed since review"""
    
    implementation_info = checkpoint_data['implementation_info']
    review_branch = checkpoint_data['review_session']['implementation_branch']
    
    verification_results = {
        'implementation_unchanged': True,
        'branch_exists': True,
        'commit_count_stable': True,
        'files_unchanged': True,
        'issues_detected': []
    }
    
    try:
        # Check if branch still exists
        branch_check = subprocess.run([
            'git', 'show-ref', '--verify', f'refs/heads/{review_branch}'
        ], capture_output=True, text=True)
        
        if branch_check.returncode != 0:
            verification_results['implementation_unchanged'] = False
            verification_results['branch_exists'] = False
            verification_results['issues_detected'].append(
                f"Implementation branch missing: {review_branch}"
            )
            return verification_results
        
        # Check commit count
        current_commits = subprocess.run([
            'git', 'rev-list', '--count', review_branch
        ], capture_output=True, text=True, check=True)
        
        current_commit_count = int(current_commits.stdout.strip())
        original_commit_count = implementation_info.get('commits_count', 0)
        
        if current_commit_count != original_commit_count:
            verification_results['implementation_unchanged'] = False
            verification_results['commit_count_stable'] = False
            verification_results['issues_detected'].append(
                f"Commit count changed: {original_commit_count} -> {current_commit_count}"
            )
        
        # Check for uncommitted changes
        status_check = subprocess.run([
            'git', 'status', '--porcelain'
        ], capture_output=True, text=True, check=True)
        
        if status_check.stdout.strip():
            verification_results['implementation_unchanged'] = False
            verification_results['files_unchanged'] = False
            verification_results['issues_detected'].append(
                "Uncommitted changes detected in working directory"
            )
        
        # Verify size hasn't changed (critical for size compliance)
        current_size = measure_current_implementation_size(review_branch)
        original_size = checkpoint_data['size_compliance']['actual_lines']
        
        if current_size != original_size:
            verification_results['implementation_unchanged'] = False
            verification_results['issues_detected'].append(
                f"Implementation size changed: {original_size} -> {current_size} lines"
            )
    
    except Exception as e:
        verification_results['implementation_unchanged'] = False
        verification_results['issues_detected'].append(
            f"Verification error: {str(e)}"
        )
    
    return verification_results

def assess_revalidation_needs(checkpoint_data):
    """Assess what aspects of review need revalidation"""
    
    checkpoint_time = datetime.fromisoformat(checkpoint_data['checkpoint_timestamp'])
    current_time = datetime.now()
    time_since_review = (current_time - checkpoint_time).total_seconds() / 3600  # hours
    
    revalidation_needs = {
        'size_compliance': False,
        'test_coverage': False, 
        'kcp_compliance': False,
        'security_review': False,
        'performance_review': False,
        'reasons': []
    }
    
    # Size compliance - revalidate if implementation changed
    if not verify_implementation_unchanged(checkpoint_data)['implementation_unchanged']:
        revalidation_needs['size_compliance'] = True
        revalidation_needs['reasons'].append("Implementation changed - size may have changed")
    
    # Test coverage - revalidate if tests were modified
    test_files_changed = check_test_files_modified(checkpoint_data, checkpoint_time)
    if test_files_changed:
        revalidation_needs['test_coverage'] = True
        revalidation_needs['reasons'].append("Test files modified since review")
    
    # KCP compliance - revalidate if core patterns changed
    kcp_files_changed = check_kcp_pattern_files_modified(checkpoint_data, checkpoint_time)
    if kcp_files_changed:
        revalidation_needs['kcp_compliance'] = True
        revalidation_needs['reasons'].append("KCP pattern implementation files changed")
    
    # Performance - revalidate if review was more than 24 hours ago
    if time_since_review > 24:
        revalidation_needs['performance_review'] = True
        revalidation_needs['reasons'].append("Performance review over 24 hours old")
    
    # Security - revalidate if security-sensitive files changed
    security_files_changed = check_security_files_modified(checkpoint_data, checkpoint_time)
    if security_files_changed:
        revalidation_needs['security_review'] = True
        revalidation_needs['reasons'].append("Security-sensitive files modified")
    
    return revalidation_needs

def determine_review_recovery_actions(checkpoint, verification, revalidation):
    """Determine what actions are needed to recover review state"""
    
    recovery_actions = []
    
    # Handle implementation verification issues
    if not verification['implementation_unchanged']:
        for issue in verification['issues_detected']:
            if 'branch missing' in issue.lower():
                recovery_actions.append({
                    'type': 'CRITICAL_ERROR',
                    'description': issue,
                    'priority': 'CRITICAL',
                    'action': 'Cannot continue review - implementation branch missing'
                })
            elif 'size changed' in issue.lower():
                recovery_actions.append({
                    'type': 'REVALIDATE_SIZE',
                    'description': issue,
                    'priority': 'CRITICAL',
                    'action': 'Re-measure size compliance immediately'
                })
            elif 'commit count changed' in issue.lower() or 'uncommitted changes' in issue.lower():
                recovery_actions.append({
                    'type': 'INVESTIGATE_CHANGES',
                    'description': issue,
                    'priority': 'HIGH',
                    'action': 'Investigate and validate implementation changes'
                })
    
    # Handle revalidation needs
    if revalidation['size_compliance']:
        recovery_actions.append({
            'type': 'REVALIDATE_SIZE_COMPLIANCE',
            'description': 'Size compliance needs revalidation',
            'priority': 'CRITICAL',
            'action': 'Re-run tmc-pr-line-counter.sh and validate results'
        })
    
    if revalidation['test_coverage']:
        recovery_actions.append({
            'type': 'REVALIDATE_TEST_COVERAGE',
            'description': 'Test coverage needs revalidation',
            'priority': 'HIGH',
            'action': 'Re-run test coverage analysis and validation'
        })
    
    if revalidation['kcp_compliance']:
        recovery_actions.append({
            'type': 'REVALIDATE_KCP_PATTERNS',
            'description': 'KCP compliance needs revalidation',
            'priority': 'HIGH',
            'action': 'Re-review KCP pattern implementation'
        })
    
    # Check if review is complete and ready to proceed
    review_decision = checkpoint.get('review_decision', {})
    if review_decision.get('result') in ['APPROVED', 'CHANGES_REQUIRED']:
        if not recovery_actions:  # No critical issues
            recovery_actions.append({
                'type': 'PROCEED_WITH_DECISION',
                'description': f"Review complete with decision: {review_decision.get('result')}",
                'priority': 'NORMAL',
                'action': 'Proceed with review decision'
            })
    else:
        recovery_actions.append({
            'type': 'COMPLETE_REVIEW',
            'description': 'Review not complete - finish review process',
            'priority': 'HIGH',
            'action': 'Complete remaining review steps'
        })
    
    return recovery_actions
```

### Review Revalidation

```python
def revalidate_critical_review_areas(checkpoint_data):
    """Re-validate critical review areas after recovery"""
    
    print("🔍 RE-VALIDATING CRITICAL REVIEW AREAS")
    
    effort_branch = checkpoint_data['review_session']['implementation_branch']
    
    # Always revalidate size compliance (most critical)
    size_revalidation = revalidate_size_compliance(effort_branch)
    
    # Revalidate test coverage if needed
    coverage_revalidation = revalidate_test_coverage(checkpoint_data)
    
    # Compare with original review results
    original_size = checkpoint_data['size_compliance']
    size_consistent = compare_size_validations(original_size, size_revalidation)
    
    original_coverage = checkpoint_data['test_coverage_validation']
    coverage_consistent = compare_coverage_validations(original_coverage, coverage_revalidation)
    
    # Determine if review results are still valid
    review_still_valid = (
        size_revalidation['compliant'] and
        size_consistent['consistent'] and
        coverage_consistent['consistent']
    )
    
    return {
        'revalidation_timestamp': datetime.now().isoformat(),
        'size_revalidation': size_revalidation,
        'coverage_revalidation': coverage_revalidation,
        'size_consistent': size_consistent,
        'coverage_consistent': coverage_consistent,
        'review_still_valid': review_still_valid,
        'action_required': 'NONE' if review_still_valid else 'COMPLETE_REVALIDATION'
    }

def revalidate_size_compliance(effort_branch):
    """Re-validate size compliance using mandatory tool"""
    
    try:
        # Use the same tool as original validation
        result = subprocess.run([
            '/workspaces/kcp-shared-tools/tmc-pr-line-counter.sh',
            '-c', effort_branch
        ], capture_output=True, text=True, check=True)
        
        # Parse line count
        output_lines = result.stdout.strip().split('\n')
        line_count = int(output_lines[-1].split()[-1])
        
        return {
            'compliant': line_count <= 800,
            'actual_lines': line_count,
            'limit': 800,
            'margin': 800 - line_count,
            'tool_used': 'tmc-pr-line-counter.sh',
            'raw_output': result.stdout.strip(),
            'revalidation_successful': True
        }
        
    except Exception as e:
        return {
            'compliant': False,
            'revalidation_successful': False,
            'error': str(e),
            'critical_failure': True
        }

def compare_size_validations(original, revalidation):
    """Compare original and revalidated size measurements"""
    
    original_lines = original.get('actual_lines', 0)
    revalidated_lines = revalidation.get('actual_lines', 0)
    
    consistent = original_lines == revalidated_lines
    
    return {
        'consistent': consistent,
        'original_lines': original_lines,
        'revalidated_lines': revalidated_lines,
        'difference': abs(original_lines - revalidated_lines),
        'compliance_changed': original.get('compliant') != revalidation.get('compliant')
    }
```

## State Persistence

```bash
# Save code review checkpoint
CHECKPOINT_DIR="/workspaces/software-factory-2.0-template/checkpoints/active"
EFFORT_ID="phase${PHASE}-wave${WAVE}-effort${EFFORT_NUM}"
CHECKPOINT_FILE="$CHECKPOINT_DIR/code-reviewer-code-review-${EFFORT_ID}-$(date +%Y%m%d-%H%M%S).yaml"

# Review artifacts backup
REVIEW_DIR="/workspaces/efforts/phase${PHASE}/wave${WAVE}/effort${EFFORT_NUM}"

# Review report backup
REVIEW_REPORT="$REVIEW_DIR/CODE-REVIEW-REPORT.md"
REVIEW_FEEDBACK="$REVIEW_DIR/REVIEW-FEEDBACK.md"

# Save all review artifacts
cp "$CHECKPOINT_FILE" "$REVIEW_DIR/review-checkpoint.yaml"
generate_review_report "$CHECKPOINT_FILE" > "$REVIEW_REPORT"
generate_review_feedback "$CHECKPOINT_FILE" > "$REVIEW_FEEDBACK"

# Commit review artifacts
git add checkpoints/ efforts/
git commit -m "checkpoint: CODE_REVIEW complete for ${EFFORT_ID} - ${REVIEW_RESULT}"
git push
```

## Critical Recovery Points

