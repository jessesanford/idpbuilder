# Code Reviewer - EFFORT_PLANNING State Checkpoint

## When to Save State

Save checkpoint at these critical effort planning points:

1. **Requirements Analysis Complete**
   - Effort scope fully analyzed
   - Dependencies identified and validated
   - KCP-specific requirements understood

2. **Architecture Design Complete**
   - Technical architecture designed
   - Component interfaces defined
   - KCP patterns identified and planned

3. **Implementation Plan Created**
   - Detailed implementation sequence defined
   - File structure planned
   - Size estimates calculated

4. **Test Strategy Finalized**
   - Test coverage strategy defined
   - Multi-tenancy test scenarios planned
   - Performance test requirements set

5. **Plan Validation Complete**
   - Implementation plan validated
   - Size compliance strategy confirmed
   - Ready to spawn Software Engineer

## Required Data to Preserve

```yaml
effort_planning_checkpoint:
  # State identification
  state: "EFFORT_PLANNING"
  effort_id: "phase1-wave2-effort3-webhooks"
  phase: 1
  wave: 2
  checkpoint_timestamp: "2025-08-23T18:15:00Z"
  
  # Planning context
  planning_session:
    started_at: "2025-08-23T17:45:00Z"
    duration_minutes: 30
    requirements_source: "/workspaces/efforts/phase1/wave2/effort3-requirements.yaml"
    assigned_sw_engineer: "TBD"
    
  # Requirements analysis
  requirements_analysis:
    effort_objectives:
      - "Implement admission webhooks for resource validation"
      - "Add mutating webhooks for resource defaulting"
      - "Integrate with KCP multi-tenant architecture"
      - "Ensure proper workspace isolation"
    
    deliverables:
      - "Validating admission webhook server"
      - "Mutating admission webhook server"
      - "Webhook configuration manifests"
      - "Multi-tenant webhook routing logic"
      - "Comprehensive test suite (90% coverage)"
    
    dependencies:
      - effort: "effort1-api-types"
        reason: "Webhook validation depends on core API types"
        status: "COMPLETE"
      - effort: "effort2-controller"
        reason: "Webhook server integrates with controller patterns"
        status: "COMPLETE"
      - external: "KCP APIExport integration"
        reason: "Webhook routing requires APIExport awareness"
        status: "AVAILABLE"
        
    constraints:
      - type: "SIZE_LIMIT"
        description: "Implementation must stay under 800 lines"
        impact: "HIGH"
        mitigation: "Plan splits if necessary"
      - type: "MULTI_TENANCY"
        description: "Must support logical cluster isolation"
        impact: "CRITICAL"
        mitigation: "Design tenant-aware routing from start"
      - type: "PERFORMANCE"
        description: "Webhook latency <100ms"
        impact: "HIGH"
        mitigation: "Include performance tests in plan"
    
    kcp_specific_requirements:
      logical_cluster_impact: "HIGH"
      multi_tenancy_support: "CRITICAL"
      api_export_integration: "REQUIRED"
      syncer_compatibility: "REQUIRED"
      workspace_isolation: "CRITICAL"
      
  # Technical architecture
  architecture_design:
    component_structure:
      webhook_server:
        description: "Core webhook HTTP server with KCP routing"
        estimated_lines: 200
        key_interfaces:
          - "WebhookServer interface"
          - "LogicalClusterRouter interface"
          - "WebhookHandler interface"
        
      admission_controllers:
        description: "Validating and mutating admission logic"
        estimated_lines: 270
        key_components:
          - "ValidatingAdmissionController"
          - "MutatingAdmissionController"
          - "ResourceValidator interface"
          
      routing_logic:
        description: "KCP-aware request routing and tenant isolation"
        estimated_lines: 150
        key_components:
          - "TenantRouter"
          - "WorkspaceIsolator"
          - "APIExportIntegrator"
          
      configuration:
        description: "Webhook configuration and management"
        estimated_lines: 100
        key_components:
          - "WebhookConfiguration"
          - "ConfigurationLoader"
          - "DynamicConfigUpdater"
    
    total_estimated_size: 720  # Within size limits
    
    kcp_pattern_integration:
      multi_tenant_design:
        approach: "Logical cluster aware webhook routing"
        isolation_method: "Workspace-scoped webhook handlers"
        tenant_identification: "Extract from request context"
        
      api_export_integration:
        approach: "Webhook registration via APIExports"
        discovery_method: "APIExport watcher integration"
        routing_strategy: "Export-aware handler selection"
        
      syncer_compatibility:
        sync_behavior: "Webhook configs synchronized to physical clusters"
        conflict_resolution: "Logical cluster priority-based"
        performance_impact: "Minimal - config-only sync"
  
  # Implementation plan
  implementation_plan:
    file_structure:
      pkg/webhooks/:
        admission/:
          - validator.go        # 120 lines
          - mutator.go         # 100 lines
          - types.go           # 50 lines
        server/:
          - server.go          # 150 lines
          - routing.go         # 100 lines
          - handler.go         # 50 lines
        config/:
          - configuration.go   # 70 lines
          - loader.go         # 30 lines
          - updater.go        # 40 lines
      
      config/webhooks/:
        - validating-webhook.yaml    # 15 lines
        - mutating-webhook.yaml      # 15 lines
        - rbac.yaml                  # 20 lines
    
    implementation_sequence:
      phase_1_foundation:
        duration_hours: 4
        deliverables:
          - "Basic webhook server framework"
          - "KCP logical cluster context extraction"
          - "Request routing foundation"
        validation: "Server starts and responds to health checks"
        
      phase_2_admission:
        duration_hours: 6
        depends_on: "phase_1_foundation"
        deliverables:
          - "Validating admission webhook implementation"
          - "Resource validation logic"
          - "Error handling and logging"
        validation: "Validation webhooks accept/reject resources correctly"
        
      phase_3_mutation:
        duration_hours: 5
        depends_on: "phase_2_admission"
        deliverables:
          - "Mutating admission webhook implementation"
          - "Resource defaulting logic"
          - "Tenant-aware mutations"
        validation: "Mutation webhooks properly modify resources"
        
      phase_4_integration:
        duration_hours: 3
        depends_on: "phase_3_mutation"
        deliverables:
          - "APIExport integration"
          - "Multi-tenant routing verification"
          - "Performance optimization"
        validation: "Full multi-tenant webhook workflow works"
    
    total_estimated_duration: 18  # hours
    
  # Testing strategy
  testing_strategy:
    unit_tests:
      target_coverage: 90
      estimated_test_lines: 400
      
      test_packages:
        pkg/webhooks/admission:
          coverage_target: 95
          key_tests:
            - "Validation logic correctness"
            - "Mutation logic correctness"
            - "Error handling scenarios"
            
        pkg/webhooks/server:
          coverage_target: 90
          key_tests:
            - "Server lifecycle management"
            - "Request routing logic"
            - "Handler registration/deregistration"
            
        pkg/webhooks/config:
          coverage_target: 85
          key_tests:
            - "Configuration loading"
            - "Dynamic configuration updates"
            - "Configuration validation"
    
    integration_tests:
      estimated_test_lines: 300
      
      test_scenarios:
        single_tenant:
          - "Basic admission webhook flow"
          - "Resource validation and rejection"
          - "Resource mutation and defaulting"
          
        multi_tenant:
          - "Workspace isolation verification"
          - "Cross-tenant request blocking"
          - "Tenant-specific webhook routing"
          
        api_export_integration:
          - "Webhook discovery via APIExports"
          - "Export-based handler routing"
          - "Export lifecycle webhook updates"
          
        error_scenarios:
          - "Webhook server unavailability"
          - "Invalid webhook configuration"
          - "Network timeout handling"
    
    performance_tests:
      estimated_test_lines: 150
      
      performance_targets:
        - metric: "Webhook response latency"
          target: "<100ms p95"
          test: "Load test with realistic request patterns"
          
        - metric: "Memory usage"
          target: "<50MB steady state"
          test: "Extended operation monitoring"
          
        - metric: "CPU usage"
          target: "<10% under normal load"
          test: "Sustained webhook request processing"
    
    total_test_lines_estimate: 850
    test_to_implementation_ratio: 1.18  # Good ratio
    
  # Size management strategy
  size_management:
    estimated_implementation_lines: 720
    estimated_test_lines: 850
    total_estimated_lines: 1570
    
    size_compliance_status: "COMPLIANT"  # Implementation only: 720 < 800
    split_required: false
    
    monitoring_strategy:
      checkpoints:
        - at_lines: 200
          action: "Verify on track for size target"
        - at_lines: 400
          action: "Detailed size review and projection"
        - at_lines: 600
          action: "Final size validation before completion"
    
    contingency_plan:
      if_approaching_limit:
        threshold: 750
        actions:
          - "Immediate size assessment"
          - "Identify refactoring opportunities"
          - "Consider moving complex logic to utilities"
      
      if_exceeds_limit:
        threshold: 800
        actions:
          - "STOP implementation immediately"
          - "Request split planning from orchestrator"
          - "Document current progress for split coordination"
  
  # Risk assessment
  risk_assessment:
    high_risks:
      - risk: "KCP logical cluster routing complexity"
        probability: "MEDIUM"
        impact: "HIGH"
        mitigation: 
          - "Study existing KCP routing patterns"
          - "Create extensive integration tests"
          - "Implement comprehensive logging"
          
      - risk: "Multi-tenant isolation violations"
        probability: "MEDIUM"
        impact: "CRITICAL"
        mitigation:
          - "Design isolation-first architecture"
          - "Create tenant isolation test scenarios"
          - "Regular security review checkpoints"
          
      - risk: "Webhook performance under load"
        probability: "LOW"
        impact: "HIGH"
        mitigation:
          - "Include performance tests from start"
          - "Use profiling during development"
          - "Design for horizontal scaling"
    
    medium_risks:
      - risk: "APIExport integration complexity"
        probability: "MEDIUM"
        impact: "MEDIUM"
        mitigation:
          - "Study APIExport patterns in existing code"
          - "Create focused integration tests"
          
      - risk: "Configuration management complexity"
        probability: "LOW"
        impact: "MEDIUM"
        mitigation:
          - "Use established configuration patterns"
          - "Implement configuration validation"
    
    risk_mitigation_readiness: "HIGH"
    
  # Validation checkpoints
  validation_checkpoints:
    architectural_review:
      criteria:
        - "KCP patterns properly integrated"
        - "Multi-tenancy design verified"
        - "Component interfaces well-defined"
      required_before: "Implementation start"
      
    mid_implementation_review:
      criteria:
        - "Size targets on track"
        - "Test coverage progressing appropriately"
        - "Integration points working correctly"
      required_at: "50% implementation complete"
      
    pre_completion_review:
      criteria:
        - "All functionality implemented"
        - "Test coverage targets met"
        - "Performance targets verified"
      required_before: "Mark effort complete"
  
  # Next actions
  planned_next_actions:
    - action: "VALIDATE_PLAN_WITH_ORCHESTRATOR"
      priority: "HIGH"
      estimated_duration: "15 minutes"
      success_criteria: "Orchestrator approves plan and spawns SW Engineer"
      
    - action: "CREATE_EFFORT_WORKING_DIRECTORY"
      priority: "HIGH"
      depends_on: "VALIDATE_PLAN_WITH_ORCHESTRATOR"
      estimated_duration: "10 minutes"
      
    - action: "GENERATE_WORK_LOG_TEMPLATE"
      priority: "MEDIUM"
      depends_on: "CREATE_EFFORT_WORKING_DIRECTORY" 
      estimated_duration: "5 minutes"
      
    - action: "SPAWN_SOFTWARE_ENGINEER"
      priority: "HIGH"
      depends_on: "GENERATE_WORK_LOG_TEMPLATE"
      orchestrator_action: true
      
  # Plan quality assessment
  plan_quality_metrics:
    completeness_score: 95
    kcp_compliance_score: 92
    risk_coverage_score: 88
    size_management_score: 90
    test_strategy_score: 94
    overall_plan_quality: "EXCELLENT"
    ready_for_implementation: true
```

## Recovery Protocol

### Context Recovery After Interruption

```python
def recover_effort_planning_state(checkpoint_data):
    """Recover effort planning state from checkpoint"""
    
    print("🔄 RECOVERING EFFORT_PLANNING STATE")
    
    effort_info = checkpoint_data['effort_id']
    plan_status = checkpoint_data.get('plan_quality_metrics', {})
    
    print(f"Effort: {effort_info}")
    print(f"Plan Quality: {plan_status.get('overall_plan_quality', 'UNKNOWN')}")
    print(f"Ready for Implementation: {plan_status.get('ready_for_implementation', False)}")
    
    # Verify planning artifacts still exist
    artifacts_verification = verify_planning_artifacts(checkpoint_data)
    
    # Check if requirements have changed since planning
    requirements_changes = detect_requirements_changes(checkpoint_data)
    
    # Determine recovery actions needed
    recovery_actions = determine_planning_recovery_actions(
        checkpoint_data, artifacts_verification, requirements_changes
    )
    
    return {
        'effort_id': effort_info,
        'plan_quality': plan_status.get('overall_plan_quality'),
        'artifacts_verification': artifacts_verification,
        'requirements_changes': requirements_changes,
        'recovery_actions': recovery_actions,
        'planning_complete': plan_status.get('ready_for_implementation', False),
        'recovery_needed': len(recovery_actions) > 0
    }

def verify_planning_artifacts(checkpoint_data):
    """Verify all planning artifacts are still valid"""
    
    verification_results = {
        'artifacts_present': True,
        'artifacts_checked': [],
        'missing_artifacts': [],
        'outdated_artifacts': []
    }
    
    # Check implementation plan file
    plan_file = checkpoint_data.get('planning_session', {}).get('implementation_plan_file')
    if plan_file:
        if os.path.exists(plan_file):
            file_time = os.path.getmtime(plan_file)
            checkpoint_time = datetime.fromisoformat(checkpoint_data['checkpoint_timestamp']).timestamp()
            
            if file_time >= checkpoint_time:
                verification_results['artifacts_checked'].append(f"Plan file current: {plan_file}")
            else:
                verification_results['outdated_artifacts'].append(f"Plan file outdated: {plan_file}")
        else:
            verification_results['artifacts_present'] = False
            verification_results['missing_artifacts'].append(f"Plan file missing: {plan_file}")
    
    # Check requirements source
    req_source = checkpoint_data.get('planning_session', {}).get('requirements_source')
    if req_source and not os.path.exists(req_source):
        verification_results['artifacts_present'] = False
        verification_results['missing_artifacts'].append(f"Requirements source missing: {req_source}")
    else:
        verification_results['artifacts_checked'].append(f"Requirements source present: {req_source}")
    
    return verification_results

def detect_requirements_changes(checkpoint_data):
    """Detect if requirements have changed since planning"""
    
    checkpoint_time = datetime.fromisoformat(checkpoint_data['checkpoint_timestamp'])
    req_source = checkpoint_data.get('planning_session', {}).get('requirements_source')
    
    changes = {
        'requirements_modified': False,
        'new_dependencies': [],
        'changed_constraints': [],
        'scope_changes': []
    }
    
    if req_source and os.path.exists(req_source):
        req_mod_time = datetime.fromtimestamp(os.path.getmtime(req_source))
        
        if req_mod_time > checkpoint_time:
            changes['requirements_modified'] = True
            
            # Load current requirements and compare
            try:
                with open(req_source, 'r') as f:
                    current_requirements = yaml.safe_load(f)
                
                original_requirements = checkpoint_data['requirements_analysis']
                
                # Compare key sections
                changes.update(compare_requirements_sections(
                    original_requirements, current_requirements
                ))
                
            except Exception as e:
                changes['comparison_error'] = str(e)
    
    return changes

def determine_planning_recovery_actions(checkpoint, artifacts, changes):
    """Determine what actions are needed to recover planning state"""
    
    recovery_actions = []
    
    # Handle missing artifacts
    if not artifacts['artifacts_present']:
        for missing in artifacts['missing_artifacts']:
            recovery_actions.append({
                'type': 'RECREATE_ARTIFACT',
                'description': f"Recreate missing artifact: {missing}",
                'priority': 'CRITICAL'
            })
    
    # Handle outdated artifacts
    if artifacts['outdated_artifacts']:
        recovery_actions.append({
            'type': 'UPDATE_ARTIFACTS',
            'description': f"Update outdated artifacts: {artifacts['outdated_artifacts']}",
            'priority': 'HIGH'
        })
    
    # Handle requirements changes
    if changes['requirements_modified']:
        recovery_actions.append({
            'type': 'REVALIDATE_PLAN',
            'description': "Requirements changed - revalidate implementation plan",
            'priority': 'HIGH',
            'details': changes
        })
    
    # Check if plan needs updates based on changes
    if changes.get('scope_changes'):
        recovery_actions.append({
            'type': 'REVISE_PLAN_SCOPE',
            'description': f"Revise plan scope due to changes: {changes['scope_changes']}",
            'priority': 'HIGH'
        })
    
    if changes.get('new_dependencies'):
        recovery_actions.append({
            'type': 'UPDATE_DEPENDENCIES',
            'description': f"Update plan for new dependencies: {changes['new_dependencies']}",
            'priority': 'MEDIUM'
        })
    
    # Check if ready to proceed or need to complete planning
    plan_quality = checkpoint.get('plan_quality_metrics', {})
    if not plan_quality.get('ready_for_implementation', False):
        recovery_actions.append({
            'type': 'COMPLETE_PLANNING',
            'description': "Planning not complete - finish plan validation and approval",
            'priority': 'HIGH'
        })
    
    return recovery_actions
```

### Plan Validation Recovery

```python
def revalidate_implementation_plan(checkpoint_data):
    """Re-validate implementation plan after recovery"""
    
    print("🔍 RE-VALIDATING IMPLEMENTATION PLAN")
    
    plan_data = checkpoint_data.get('implementation_plan', {})
    architecture = checkpoint_data.get('architecture_design', {})
    
    # Re-run plan validation
    validation_results = run_plan_validation(plan_data, architecture)
    
    # Compare with checkpoint validation
    checkpoint_quality = checkpoint_data.get('plan_quality_metrics', {})
    validation_comparison = compare_plan_validations(checkpoint_quality, validation_results)
    
    # Check size estimates still valid
    size_validation = revalidate_size_estimates(plan_data)
    
    # Check KCP compliance still valid
    kcp_validation = revalidate_kcp_compliance(architecture)
    
    still_valid = (
        validation_results['overall_quality'] in ['EXCELLENT', 'GOOD', 'PASS'] and
        validation_comparison['consistent'] and
        size_validation['still_compliant'] and
        kcp_validation['still_compliant']
    )
    
    return {
        'validation_timestamp': datetime.now().isoformat(),
        'still_valid': still_valid,
        'plan_validation': validation_results,
        'comparison_with_checkpoint': validation_comparison,
        'size_validation': size_validation,
        'kcp_validation': kcp_validation,
        'action_required': 'NONE' if still_valid else 'REVISE_PLAN'
    }
```

## State Persistence

```bash
# Save effort planning checkpoint
CHECKPOINT_DIR="/workspaces/software-factory-2.0-template/checkpoints/active"
EFFORT_ID="phase${PHASE}-wave${WAVE}-effort${EFFORT_NUM}"
CHECKPOINT_FILE="$CHECKPOINT_DIR/code-reviewer-effort-planning-${EFFORT_ID}-$(date +%Y%m%d-%H%M%S).yaml"

# Planning artifacts backup
PLANNING_DIR="/workspaces/efforts/phase${PHASE}/wave${WAVE}/effort${EFFORT_NUM}"
mkdir -p "$PLANNING_DIR"

# Implementation plan backup
PLAN_FILE="$PLANNING_DIR/IMPLEMENTATION-PLAN.md"
WORK_LOG_TEMPLATE="$PLANNING_DIR/work-log-template.md"

# Save all planning artifacts
cp "$CHECKPOINT_FILE" "$PLANNING_DIR/planning-checkpoint.yaml"
generate_implementation_plan_document "$CHECKPOINT_FILE" > "$PLAN_FILE"
generate_work_log_template "$CHECKPOINT_FILE" > "$WORK_LOG_TEMPLATE"

# Commit planning artifacts
git add checkpoints/ efforts/
git commit -m "checkpoint: EFFORT_PLANNING complete for ${EFFORT_ID} - ready for implementation"
git push
```

## Critical Recovery Points

