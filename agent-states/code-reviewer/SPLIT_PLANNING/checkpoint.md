# Code Reviewer - SPLIT_PLANNING State Checkpoint

## When to Save State

Save checkpoint at these critical split planning points:

1. **Size Violation Analysis Complete**
   - Oversized implementation analyzed
   - Root causes of size violation identified
   - Component size distribution understood

2. **Split Strategy Design Complete**
   - Split boundaries identified and evaluated
   - Split options generated and compared
   - Optimal split strategy selected

3. **Split Plan Validation Complete**
   - All splits validated to be <800 lines
   - Functional integrity preservation confirmed
   - KCP compliance across splits verified

4. **Execution Plan Finalized**
   - Detailed execution sequence planned
   - Integration strategy defined
   - Risk mitigation strategies prepared

5. **Split Plan Approved**
   - Plan validation passed all checks
   - Ready for split execution
   - First split ready to be implemented

## Required Data to Preserve

```yaml
split_planning_checkpoint:
  # State identification
  state: "SPLIT_PLANNING"
  original_effort: "phase1-wave2-effort3-webhooks"
  phase: 1
  wave: 2
  checkpoint_timestamp: "2025-08-23T20:15:00Z"
  
  # Split planning context
  planning_session:
    started_at: "2025-08-23T19:45:00Z"
    duration_minutes: 30
    trigger_reason: "SIZE_VIOLATION"
    original_implementation_location: "/workspaces/efforts/phase1/wave2/effort3-webhooks"
    original_branch: "phase1/wave2/effort3-webhooks"
    size_violation_detected_at: "2025-08-23T19:30:00Z"
    
  # Size violation analysis
  size_violation_analysis:
    violation_details:
      original_size_lines: 1247
      size_limit: 800
      overage_lines: 447
      overage_percentage: 55.9
      violation_severity: "HIGH"
      
    size_measurement_data:
      tool_used: "line-counter.sh"
      measurement_timestamp: "2025-08-23T19:30:00Z"
      raw_output: |
        pkg/webhooks/admission/validator.go: 156 lines
        pkg/webhooks/admission/mutator.go: 128 lines
        pkg/webhooks/admission/types.go: 67 lines
        pkg/webhooks/server/server.go: 198 lines
        pkg/webhooks/server/routing.go: 134 lines
        pkg/webhooks/server/handler.go: 89 lines
        pkg/webhooks/config/configuration.go: 92 lines
        pkg/webhooks/config/loader.go: 56 lines
        pkg/webhooks/config/updater.go: 78 lines
        pkg/webhooks/utils/helpers.go: 43 lines
        pkg/webhooks/utils/validation.go: 67 lines
        config/webhooks/manifests.go: 38 lines
        test/webhooks/integration_test.go: 101 lines
        Total: 1247
    
    component_size_analysis:
      admission_logic:
        files: ["validator.go", "mutator.go", "types.go"]
        total_lines: 351
        percentage_of_total: 28.1
        
      server_infrastructure:
        files: ["server.go", "routing.go", "handler.go"]
        total_lines: 421
        percentage_of_total: 33.8
        
      configuration:
        files: ["configuration.go", "loader.go", "updater.go"]
        total_lines: 226
        percentage_of_total: 18.1
        
      utilities:
        files: ["helpers.go", "validation.go"]
        total_lines: 110
        percentage_of_total: 8.8
        
      configuration_manifests:
        files: ["manifests.go"]
        total_lines: 38
        percentage_of_total: 3.0
        
      integration_tests:
        files: ["integration_test.go"]
        total_lines: 101
        percentage_of_total: 8.1
    
    root_cause_analysis:
      primary_causes:
        - "Server infrastructure more complex than estimated"
        - "Configuration management expanded during implementation"
        - "Additional utility functions added for KCP integration"
        
      contributing_factors:
        - "APIExport integration required more code than planned"
        - "Multi-tenant routing logic more involved than expected"
        - "Error handling and logging expanded during development"
        
      lessons_learned:
        - "KCP integration complexity often underestimated"
        - "Server components tend to grow during implementation"
        - "Configuration management requires more infrastructure"
        
  # Split strategy design
  split_strategy_design:
    analysis_timestamp: "2025-08-23T19:55:00Z"
    
    # Identified split boundaries
    potential_boundaries:
      component_boundaries:
        - boundary: "ADMISSION_LOGIC_SEPARATION"
          description: "Separate admission validation and mutation logic"
          estimated_split_size: 351
          feasibility: "HIGH"
          
        - boundary: "SERVER_INFRASTRUCTURE_SEPARATION"
          description: "Separate core server and routing infrastructure"
          estimated_split_size: 421
          feasibility: "MEDIUM"
          
        - boundary: "CONFIGURATION_SEPARATION"
          description: "Separate configuration management"
          estimated_split_size: 226
          feasibility: "LOW"  # Too small for standalone split
          
      functional_boundaries:
        - boundary: "VALIDATION_VS_MUTATION"
          description: "Split validating and mutating admission webhooks"
          validation_size: 189
          mutation_size: 162
          feasibility: "MEDIUM"
          
        - boundary: "CORE_VS_ROUTING"
          description: "Split core server from routing logic"
          core_size: 287
          routing_size: 134
          feasibility: "HIGH"
          
      interface_boundaries:
        - boundary: "API_VS_IMPLEMENTATION"
          description: "Split interface definitions from implementations"
          api_size: 156
          implementation_size: 1091
          feasibility: "LOW"  # Implementation still oversized
    
    # Generated split options
    split_options:
      option_1_component_based:
        strategy: "COMPONENT_BASED"
        description: "Split along major component boundaries"
        splits:
          - id: "webhooks-admission"
            description: "Admission validation and mutation logic"
            estimated_lines: 351
            components: ["admission/validator.go", "admission/mutator.go", "admission/types.go"]
            
          - id: "webhooks-server"
            description: "Server infrastructure and routing"
            estimated_lines: 421  
            components: ["server/server.go", "server/routing.go", "server/handler.go"]
            
          - id: "webhooks-config-utils"
            description: "Configuration, utilities, and manifests"
            estimated_lines: 475  # 226 + 110 + 38 + 101 (includes tests)
            components: ["config/*", "utils/*", "manifests.go", "integration_test.go"]
            
        total_splits: 3
        max_split_size: 475
        compliant: true
        
      option_2_functional:
        strategy: "FUNCTIONAL"
        description: "Split along functional boundaries"
        splits:
          - id: "webhooks-validation"
            description: "Validating admission webhooks and infrastructure"
            estimated_lines: 467  # validator + server components
            components: ["admission/validator.go", "server/*", "config/*"]
            
          - id: "webhooks-mutation"  
            description: "Mutating admission webhooks and utilities"
            estimated_lines: 389  # mutator + types + utils + tests
            components: ["admission/mutator.go", "admission/types.go", "utils/*", "integration_test.go"]
            
          - id: "webhooks-config"
            description: "Configuration management and manifests"
            estimated_lines: 391  # All config + manifests + handler
            components: ["config/*", "manifests.go", "server/handler.go"]
            
        total_splits: 3
        max_split_size: 467
        compliant: true
        
      option_3_hybrid:
        strategy: "HYBRID" 
        description: "Combination of component and functional splits"
        splits:
          - id: "webhooks-core-server"
            description: "Core server infrastructure and configuration"
            estimated_lines: 384  # server + config (no routing)
            components: ["server/server.go", "server/handler.go", "config/*"]
            
          - id: "webhooks-admission-routing"
            description: "Admission logic with routing integration"
            estimated_lines: 485  # admission + routing + utils
            components: ["admission/*", "server/routing.go", "utils/*"]
            
          - id: "webhooks-integration-config"
            description: "Integration tests and deployment configuration"
            estimated_lines: 378  # tests + manifests + remaining
            components: ["integration_test.go", "manifests.go", "remaining files"]
            
        total_splits: 3
        max_split_size: 485
        compliant: true
    
    # Selected optimal strategy
    selected_strategy: "COMPONENT_BASED"  # option_1_component_based
    selection_rationale:
      - "Cleanest separation of concerns"
      - "Lowest integration complexity"
      - "Best preserves existing architecture"
      - "Easiest to implement and test"
      - "All splits well under 800-line limit"
      
    selection_confidence: 92
    
  # Detailed split plan (selected strategy)
  detailed_split_plan:
    strategy: "COMPONENT_BASED"
    total_splits: 3
    
    splits:
      split_1:
        id: "webhooks-admission"
        description: "Admission webhook validation and mutation logic"
        
        # Size and compliance
        estimated_lines: 351
        size_compliant: true
        size_margin: 449
        target_completion_size: 400  # Allow for some expansion
        
        # Components included
        components:
          - file: "pkg/webhooks/admission/validator.go"
            lines: 156
            purpose: "Validating admission webhook implementation"
            
          - file: "pkg/webhooks/admission/mutator.go"
            lines: 128
            purpose: "Mutating admission webhook implementation"
            
          - file: "pkg/webhooks/admission/types.go"
            lines: 67
            purpose: "Admission webhook types and interfaces"
        
        # Interfaces and dependencies
        interfaces_provided:
          - "AdmissionValidator interface"
          - "AdmissionMutator interface"
          - "WebhookHandler interface"
          
        interfaces_required:
          - "WebhookServer interface (from server split)"
          - "ConfigurationManager interface (from config split)"
          
        external_dependencies:
          - "KCP logical cluster client"
          - "Kubernetes admission types"
          - "APIExport client interfaces"
          
        internal_dependencies:
          - dependency_on: "webhooks-server"
            reason: "Needs server registration and request routing"
            interface: "WebhookServer.RegisterHandler()"
            
          - dependency_on: "webhooks-config-utils"
            reason: "Needs access to webhook configuration"
            interface: "ConfigurationManager.GetWebhookConfig()"
            
        # KCP compliance considerations
        kcp_compliance:
          multi_tenancy_support: "PRESERVED"
          logical_cluster_integration: "MAINTAINED"
          workspace_isolation: "ENFORCED"
          api_export_integration: "COMPATIBLE"
          
        # Testing strategy
        testing_strategy:
          unit_tests_location: "test/admission/"
          integration_test_approach: "Mock server interfaces"
          multi_tenant_test_scenarios: 4
          target_coverage: 90
          
      split_2:
        id: "webhooks-server"
        description: "Core webhook server infrastructure and routing"
        
        # Size and compliance  
        estimated_lines: 421
        size_compliant: true
        size_margin: 379
        target_completion_size: 450
        
        # Components included
        components:
          - file: "pkg/webhooks/server/server.go"
            lines: 198
            purpose: "Core webhook HTTP server implementation"
            
          - file: "pkg/webhooks/server/routing.go"
            lines: 134
            purpose: "Multi-tenant request routing logic"
            
          - file: "pkg/webhooks/server/handler.go"
            lines: 89
            purpose: "HTTP request handlers and middleware"
            
        # Interfaces and dependencies
        interfaces_provided:
          - "WebhookServer interface"
          - "TenantRouter interface"
          - "RequestHandler interface"
          
        interfaces_required:
          - "AdmissionHandler interface (from admission split)"
          - "ConfigurationManager interface (from config split)"
          
        external_dependencies:
          - "KCP APIExport client"
          - "HTTP server libraries"
          - "TLS certificate management"
          
        internal_dependencies:
          - dependency_on: "webhooks-admission" 
            reason: "Routes requests to admission handlers"
            interface: "AdmissionHandler.ProcessAdmissionRequest()"
            
          - dependency_on: "webhooks-config-utils"
            reason: "Server configuration and TLS setup"
            interface: "ConfigurationManager.GetServerConfig()"
            
        # KCP compliance considerations
        kcp_compliance:
          multi_tenancy_support: "CRITICAL_COMPONENT"
          logical_cluster_integration: "ROUTING_LAYER"
          workspace_isolation: "ENFORCED_AT_ENTRY"
          api_export_integration: "PRIMARY_INTEGRATION_POINT"
          
      split_3:
        id: "webhooks-config-utils"
        description: "Configuration management, utilities, and deployment"
        
        # Size and compliance
        estimated_lines: 475
        size_compliant: true
        size_margin: 325
        target_completion_size: 500
        
        # Components included
        components:
          - file: "pkg/webhooks/config/configuration.go"
            lines: 92
            purpose: "Webhook configuration structures"
            
          - file: "pkg/webhooks/config/loader.go" 
            lines: 56
            purpose: "Configuration loading and validation"
            
          - file: "pkg/webhooks/config/updater.go"
            lines: 78
            purpose: "Dynamic configuration updates"
            
          - file: "pkg/webhooks/utils/helpers.go"
            lines: 43
            purpose: "Common utility functions"
            
          - file: "pkg/webhooks/utils/validation.go"
            lines: 67
            purpose: "Validation helper functions"
            
          - file: "config/webhooks/manifests.go"
            lines: 38
            purpose: "Kubernetes manifest generation"
            
          - file: "test/webhooks/integration_test.go"
            lines: 101
            purpose: "Integration tests for complete webhook flow"
        
        # Interfaces and dependencies
        interfaces_provided:
          - "ConfigurationManager interface"
          - "ValidationHelper interface"
          - "ManifestGenerator interface"
          
        interfaces_required:
          - "No internal interface dependencies"
          
        external_dependencies:
          - "Kubernetes client libraries"
          - "YAML/JSON parsing libraries"
          - "File system access"
          
        internal_dependencies: []  # No dependencies on other splits
        
        # KCP compliance considerations
        kcp_compliance:
          multi_tenancy_support: "CONFIGURATION_SUPPORT"
          logical_cluster_integration: "CONFIG_AWARE"
          workspace_isolation: "SUPPORTS_ISOLATION"
          api_export_integration: "CONFIG_GENERATION"
    
  # Integration strategy
  integration_strategy:
    approach: "SEQUENTIAL_INTEGRATION_WITH_VALIDATION"
    integration_order: ["webhooks-config-utils", "webhooks-server", "webhooks-admission"]
    
    integration_sequence:
      phase_1_foundation:
        duration_hours: 3
        split_to_implement: "webhooks-config-utils"
        rationale: "No dependencies - can be implemented independently"
        deliverables:
          - "Configuration management working"
          - "Utility functions available"
          - "Manifest generation functional"
        validation_criteria:
          - "Configuration loading works correctly"
          - "Utility functions pass unit tests"
          - "Manifest generation produces valid YAML"
          
      phase_2_server:
        duration_hours: 5
        split_to_implement: "webhooks-server"
        depends_on: "phase_1_foundation"
        rationale: "Depends on config split for server configuration"
        deliverables:
          - "HTTP server infrastructure working"
          - "Multi-tenant routing functional"
          - "Integration with config split validated"
        validation_criteria:
          - "Server starts successfully"
          - "Request routing works correctly"
          - "Multi-tenant isolation verified"
          
      phase_3_admission:
        duration_hours: 4
        split_to_implement: "webhooks-admission" 
        depends_on: "phase_2_server"
        rationale: "Depends on server split for request handling"
        deliverables:
          - "Admission webhook logic working"
          - "Integration with server split validated"
          - "End-to-end webhook flow functional"
        validation_criteria:
          - "Admission webhooks process requests correctly"
          - "Validation and mutation logic works"
          - "Multi-tenant admission processing verified"
          
      phase_4_integration:
        duration_hours: 2
        rationale: "Final integration validation and optimization"
        deliverables:
          - "All splits integrated and tested together"
          - "Performance targets validated"
          - "Complete functionality verified"
        validation_criteria:
          - "Original functionality fully preserved"
          - "Performance within acceptable ranges"
          - "All multi-tenant scenarios working"
    
    total_estimated_duration: 14  # hours
    
    # Interface integration points
    integration_points:
      config_to_server:
        interface: "ConfigurationManager → WebhookServer"
        integration_method: "Dependency injection during server initialization"
        validation_approach: "Unit tests with mock configurations"
        
      server_to_admission:
        interface: "WebhookServer → AdmissionHandler"
        integration_method: "Handler registration and request routing"
        validation_approach: "Integration tests with real admission requests"
        
      config_to_admission:
        interface: "ConfigurationManager → AdmissionHandler"
        integration_method: "Configuration access during admission processing"  
        validation_approach: "Configuration change tests"
    
    # Integration testing strategy
    integration_testing:
      test_levels:
        - level: "Interface contract tests"
          description: "Validate interface contracts between splits"
          scope: "Each integration point"
          
        - level: "Component integration tests"
          description: "Test integration between two splits"
          scope: "Pairwise split combinations"
          
        - level: "End-to-end integration tests"
          description: "Test complete webhook flow across all splits"
          scope: "Full admission webhook scenarios"
          
        - level: "Multi-tenant integration tests"
          description: "Validate multi-tenancy across split boundaries"
          scope: "Cross-workspace scenarios"
  
  # Risk assessment and mitigation
  risk_assessment:
    split_planning_risks:
      high_risks:
        - risk: "Interface integration complexity underestimated"
          probability: "MEDIUM"
          impact: "HIGH"
          mitigation:
            - "Design interfaces early and validate with prototypes"
            - "Create comprehensive interface contract tests"
            - "Plan for interface evolution and versioning"
            
        - risk: "Multi-tenant behavior broken across splits"
          probability: "MEDIUM"
          impact: "CRITICAL"
          mitigation:
            - "Preserve logical cluster context in all interfaces"
            - "Test multi-tenancy at each integration point"
            - "Validate workspace isolation across split boundaries"
            
        - risk: "Performance degradation from split overhead"
          probability: "LOW"
          impact: "HIGH"
          mitigation:
            - "Measure performance at each integration milestone"
            - "Optimize critical paths if performance issues detected"
            - "Consider inlining if overhead becomes significant"
    
      medium_risks:
        - risk: "Configuration management complexity"
          probability: "MEDIUM"
          impact: "MEDIUM"
          mitigation:
            - "Centralize configuration access patterns"
            - "Use consistent configuration interfaces across splits"
            - "Test configuration changes with all splits"
            
        - risk: "Testing coordination across splits"
          probability: "MEDIUM"
          impact: "MEDIUM"
          mitigation:
            - "Establish shared testing infrastructure early"
            - "Create integration test coordination strategy"
            - "Plan for cross-split test data management"
    
    risk_mitigation_confidence: 85
    
  # Validation results
  split_plan_validation:
    validation_timestamp: "2025-08-23T20:10:00Z"
    
    size_compliance_validation:
      all_splits_compliant: true
      compliance_details:
        - split: "webhooks-admission"
          estimated_lines: 351
          compliant: true
          margin: 449
          
        - split: "webhooks-server"
          estimated_lines: 421
          compliant: true
          margin: 379
          
        - split: "webhooks-config-utils"
          estimated_lines: 475
          compliant: true
          margin: 325
      
      overall_compliance: "PASSED"
      
    functional_integrity_validation:
      data_flow_preserved: 94
      api_contracts_maintained: 89
      error_handling_complete: 87
      interface_design_quality: 91
      dependency_management: 86
      
      overall_integrity_score: 89.4
      integrity_status: "GOOD"
      
    kcp_compliance_validation:
      multi_tenancy_preserved: 93
      logical_cluster_context_flow: 91
      workspace_isolation_maintained: 89
      api_export_integration_intact: 87
      syncer_compatibility_preserved: 85
      
      overall_kcp_score: 89.0
      kcp_status: "GOOD"
      
    integration_feasibility_validation:
      interface_complexity: "MEDIUM"
      dependency_chain_length: "SHORT"
      integration_points_count: "REASONABLE"
      circular_dependencies: "NONE"
      
      feasibility_score: 82
      feasibility_status: "GOOD"
      
    execution_viability_validation:
      plan_completeness: 91
      timeline_realism: 87
      resource_requirements: "REASONABLE"
      risk_mitigation_quality: 85
      
      viability_score: 88
      viability_status: "GOOD"
      
    overall_validation_result:
      validation_score: 87.8
      validation_status: "GOOD"
      plan_approved: true
      blocking_issues: []
      recommendations:
        - "Monitor interface integration complexity closely"
        - "Validate multi-tenant behavior at each integration point"
        - "Consider performance impact of split architecture"
        
  # Next actions and readiness
  execution_readiness:
    ready_for_execution: true
    next_split_to_implement: "webhooks-config-utils"
    
    immediate_next_actions:
      - action: "CREATE_SPLIT_WORKING_DIRECTORIES"
        priority: "HIGH"
        estimated_duration: "15 minutes"
        description: "Set up working directories for all three splits"
        
      - action: "GENERATE_SPLIT_IMPLEMENTATION_PLANS"
        priority: "HIGH"
        estimated_duration: "45 minutes"
        description: "Create detailed implementation plans for each split"
        
      - action: "SPAWN_SW_ENGINEER_FOR_FIRST_SPLIT"
        priority: "HIGH"
        estimated_duration: "5 minutes"
        description: "Request orchestrator to spawn SW Engineer for config-utils split"
        
    success_criteria:
      - "All three splits implemented and tested individually"
      - "Integration between splits working correctly"
      - "Original webhook functionality fully preserved"
      - "Performance targets maintained"
      - "Multi-tenant behavior verified across all splits"
      
  # Split plan quality metrics
  plan_quality_assessment:
    size_compliance_score: 100
    functional_integrity_score: 89
    integration_feasibility_score: 82
    kcp_preservation_score: 89
    execution_viability_score: 88
    
    overall_plan_quality: 89.6
    plan_quality_grade: "GOOD"
    ready_for_execution: true
    confidence_level: 87
```

## Recovery Protocol

### Context Recovery After Interruption

```python
def recover_split_planning_state(checkpoint_data):
    """Recover split planning state from checkpoint"""
    
    print("🔄 RECOVERING SPLIT_PLANNING STATE")
    
    original_effort = checkpoint_data['original_effort']
    plan_status = checkpoint_data.get('execution_readiness', {})
    
    print(f"Original Effort: {original_effort}")
    print(f"Ready for Execution: {plan_status.get('ready_for_execution', False)}")
    print(f"Next Split: {plan_status.get('next_split_to_implement', 'UNKNOWN')}")
    
    # Verify original implementation still exists and is unchanged
    original_verification = verify_original_implementation_unchanged(checkpoint_data)
    
    # Check if split plan is still valid
    plan_validity = assess_split_plan_validity(checkpoint_data)
    
    # Determine recovery actions needed
    recovery_actions = determine_split_planning_recovery_actions(
        checkpoint_data, original_verification, plan_validity
    )
    
    return {
        'original_effort': original_effort,
        'plan_ready': plan_status.get('ready_for_execution', False),
        'original_verification': original_verification,
        'plan_validity': plan_validity,
        'recovery_actions': recovery_actions,
        'planning_complete': plan_status.get('ready_for_execution', False),
        'recovery_needed': len(recovery_actions) > 0
    }

def verify_original_implementation_unchanged(checkpoint_data):
    """Verify the original implementation hasn't changed since split planning"""
    
    original_branch = checkpoint_data['planning_session']['original_branch']
    original_size = checkpoint_data['size_violation_analysis']['violation_details']['original_size_lines']
    
    verification_results = {
        'implementation_unchanged': True,
        'size_still_matches': True,
        'branch_accessible': True,
        'issues_detected': []
    }
    
    try:
        # Check if original branch still exists
        branch_check = subprocess.run([
            'git', 'show-ref', '--verify', f'refs/heads/{original_branch}'
        ], capture_output=True, text=True)
        
        if branch_check.returncode != 0:
            verification_results['implementation_unchanged'] = False
            verification_results['branch_accessible'] = False
            verification_results['issues_detected'].append(
                f"Original branch not accessible: {original_branch}"
            )
            return verification_results
        
        # Re-measure size of original implementation
        current_size = measure_implementation_size(original_branch)
        
        if current_size != original_size:
            verification_results['implementation_unchanged'] = False
            verification_results['size_still_matches'] = False
            verification_results['issues_detected'].append(
                f"Original implementation size changed: {original_size} -> {current_size} lines"
            )
        
        # Check for recent commits on original branch
        recent_commits = check_recent_commits_since(
            original_branch, 
            checkpoint_data['checkpoint_timestamp']
        )
        
        if recent_commits:
            verification_results['implementation_unchanged'] = False
            verification_results['issues_detected'].append(
                f"New commits on original branch: {len(recent_commits)} commits"
            )
    
    except Exception as e:
        verification_results['implementation_unchanged'] = False
        verification_results['issues_detected'].append(
            f"Verification error: {str(e)}"
        )
    
    return verification_results

def assess_split_plan_validity(checkpoint_data):
    """Assess if the split plan is still valid and optimal"""
    
    split_plan = checkpoint_data['detailed_split_plan']
    validation_results = checkpoint_data['split_plan_validation']
    
    validity_assessment = {
        'plan_still_valid': True,
        'size_estimates_valid': True,
        'integration_strategy_valid': True,
        'kcp_compliance_maintained': True,
        'validity_issues': []
    }
    
    # Re-validate split size estimates
    for split in split_plan['splits']:
        split_id = split['id']
        estimated_lines = split['estimated_lines']
        
        # Check if estimate is still reasonable based on component analysis
        if estimated_lines > 750:  # Getting close to limit
            validity_assessment['validity_issues'].append(
                f"Split {split_id} estimate {estimated_lines} approaching size limit"
            )
    
    # Check if integration complexity has changed
    integration_complexity = validation_results.get('integration_feasibility_validation', {}).get('feasibility_score', 0)
    if integration_complexity < 70:
        validity_assessment['plan_still_valid'] = False
        validity_assessment['integration_strategy_valid'] = False
        validity_assessment['validity_issues'].append(
            f"Integration feasibility too low: {integration_complexity}%"
        )
    
    # Check KCP compliance scores
    kcp_score = validation_results.get('kcp_compliance_validation', {}).get('overall_kcp_score', 0)
    if kcp_score < 85:
        validity_assessment['plan_still_valid'] = False
        validity_assessment['kcp_compliance_maintained'] = False
        validity_assessment['validity_issues'].append(
            f"KCP compliance too low: {kcp_score}%"
        )
    
    return validity_assessment

def determine_split_planning_recovery_actions(checkpoint, original_verification, plan_validity):
    """Determine what actions are needed to recover split planning state"""
    
    recovery_actions = []
    
    # Handle original implementation verification issues
    if not original_verification['implementation_unchanged']:
        for issue in original_verification['issues_detected']:
            if 'branch not accessible' in issue.lower():
                recovery_actions.append({
                    'type': 'CRITICAL_ERROR',
                    'description': issue,
                    'priority': 'CRITICAL',
                    'action': 'Cannot continue - original implementation not accessible'
                })
            elif 'size changed' in issue.lower():
                recovery_actions.append({
                    'type': 'REVALIDATE_SPLIT_PLAN',
                    'description': issue,
                    'priority': 'HIGH',
                    'action': 'Re-analyze size violation and adjust split plan'
                })
            elif 'new commits' in issue.lower():
                recovery_actions.append({
                    'type': 'ASSESS_IMPLEMENTATION_CHANGES',
                    'description': issue,
                    'priority': 'HIGH',
                    'action': 'Review changes and update split plan if necessary'
                })
    
    # Handle plan validity issues
    if not plan_validity['plan_still_valid']:
        for issue in plan_validity['validity_issues']:
            if 'approaching size limit' in issue.lower():
                recovery_actions.append({
                    'type': 'ADJUST_SPLIT_SIZES',
                    'description': issue,
                    'priority': 'MEDIUM',
                    'action': 'Review and adjust split size estimates'
                })
            elif 'integration feasibility' in issue.lower():
                recovery_actions.append({
                    'type': 'REVISE_INTEGRATION_STRATEGY',
                    'description': issue,
                    'priority': 'HIGH',
                    'action': 'Simplify integration approach or revise split boundaries'
                })
            elif 'kcp compliance' in issue.lower():
                recovery_actions.append({
                    'type': 'FIX_KCP_COMPLIANCE',
                    'description': issue,
                    'priority': 'HIGH',
                    'action': 'Revise split plan to better preserve KCP patterns'
                })
    
    # Check if ready to proceed or need to complete planning
    execution_readiness = checkpoint.get('execution_readiness', {})
    if not execution_readiness.get('ready_for_execution', False):
        recovery_actions.append({
            'type': 'COMPLETE_PLANNING',
            'description': 'Split planning not complete - finish validation and approval',
            'priority': 'HIGH',
            'action': 'Complete remaining planning steps'
        })
    else:
        if not recovery_actions:  # No critical issues
            recovery_actions.append({
                'type': 'PROCEED_WITH_EXECUTION',
                'description': f"Split plan ready - start with {execution_readiness.get('next_split_to_implement', 'first split')}",
                'priority': 'NORMAL',
                'action': 'Begin split execution sequence'
            })
    
    return recovery_actions
```

### Split Plan Revalidation

```python
def revalidate_split_plan_after_recovery(checkpoint_data):
    """Re-validate split plan after recovery"""
    
    print("🔍 RE-VALIDATING SPLIT PLAN")
    
    original_branch = checkpoint_data['planning_session']['original_branch']
    
    # Re-measure original implementation size
    current_size_validation = revalidate_original_size(original_branch)
    
    # Re-validate split plan if size changed
    if current_size_validation['size_changed']:
        split_plan_revalidation = revalidate_split_strategy(
            checkpoint_data, current_size_validation
        )
    else:
        split_plan_revalidation = {
            'revalidation_needed': False,
            'plan_still_valid': True,
            'original_validation_current': True
        }
    
    # Compare with original validation results
    original_validation = checkpoint_data['split_plan_validation']
    validation_comparison = compare_split_validations(
        original_validation, split_plan_revalidation
    )
    
    # Determine if split plan is still viable
    plan_still_viable = (
        not current_size_validation['size_changed'] or
        (split_plan_revalidation.get('plan_still_valid', False))
    ) and validation_comparison.get('consistent', True)
    
    return {
        'revalidation_timestamp': datetime.now().isoformat(),
        'size_revalidation': current_size_validation,
        'plan_revalidation': split_plan_revalidation,
        'validation_comparison': validation_comparison,
        'plan_still_viable': plan_still_viable,
        'action_required': 'NONE' if plan_still_viable else 'REVISE_SPLIT_PLAN'
    }

def revalidate_original_size(original_branch):
    """Re-validate size of original implementation"""
    
    try:
        # Re-measure using same tool as original
        result = subprocess.run([
            '$PROJECT_ROOT/tools/line-counter.sh',
            '-c', original_branch
        ], capture_output=True, text=True, check=True)
        
        # Parse current size
        output_lines = result.stdout.strip().split('\n')
        current_size = int(output_lines[-1].split()[-1])
        
        return {
            'current_size': current_size,
            'size_changed': False,  # Will be set by caller when comparing
            'measurement_successful': True,
            'raw_output': result.stdout.strip()
        }
        
    except Exception as e:
        return {
            'current_size': 0,
            'size_changed': True,
            'measurement_successful': False,
            'error': str(e)
        }
```

## State Persistence

```bash
# Save split planning checkpoint
CHECKPOINT_DIR="/workspaces/software-factory-2.0-template/checkpoints/active"
ORIGINAL_EFFORT="phase${PHASE}-wave${WAVE}-effort${EFFORT_NUM}"
CHECKPOINT_FILE="$CHECKPOINT_DIR/code-reviewer-split-planning-${ORIGINAL_EFFORT}-$(date +%Y%m%d-%H%M%S).yaml"

# Split plan artifacts backup
SPLITS_BASE_DIR="/workspaces/efforts/phase${PHASE}/wave${WAVE}"
ORIGINAL_EFFORT_DIR="$SPLITS_BASE_DIR/effort${EFFORT_NUM}"

# Split plan documentation
SPLIT_PLAN_DOC="$ORIGINAL_EFFORT_DIR/SPLIT-PLAN.md"
SPLIT_STRATEGY_DOC="$ORIGINAL_EFFORT_DIR/SPLIT-STRATEGY.md"

# Create split working directories (preparation for execution)
for SPLIT_ID in $(extract_split_ids_from_checkpoint "$CHECKPOINT_FILE"); do
    SPLIT_DIR="$SPLITS_BASE_DIR/${SPLIT_ID}"
    mkdir -p "$SPLIT_DIR"
    generate_split_implementation_plan "$CHECKPOINT_FILE" "$SPLIT_ID" > "$SPLIT_DIR/IMPLEMENTATION-PLAN.md"
    generate_split_work_log_template "$CHECKPOINT_FILE" "$SPLIT_ID" > "$SPLIT_DIR/work-log-template.md"
done

# Save all planning artifacts
cp "$CHECKPOINT_FILE" "$ORIGINAL_EFFORT_DIR/split-planning-checkpoint.yaml"
generate_split_plan_document "$CHECKPOINT_FILE" > "$SPLIT_PLAN_DOC"
generate_split_strategy_document "$CHECKPOINT_FILE" > "$SPLIT_STRATEGY_DOC"

# Commit planning artifacts
git add checkpoints/ efforts/
git commit -m "checkpoint: SPLIT_PLANNING complete for ${ORIGINAL_EFFORT} - ${SPLIT_COUNT} splits planned and ready"
git push
```

## Critical Recovery Points

