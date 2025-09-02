# Code Reviewer - VALIDATION State Rules

## State Context
You are performing comprehensive validation of implementations, ensuring they meet all quality gates, compliance requirements, and integration standards before final approval.

---
### ℹ️ RULE R055.0.0 - Review Decision Making
**Source:** rule-library/RULE-REGISTRY.md#R055
**Criticality:** INFO - Best practice

VALIDATION PROTOCOL:
1. Execute comprehensive validation suite
2. Verify all quality gates pass
3. Validate integration readiness
4. Confirm KCP/Kubernetes compliance
5. Assess performance and security
6. Make final approval/rejection decision
7. Provide actionable feedback
---

## Quality Gate Validation

---
### 🚨🚨 RULE R035.0.0 - Phase Completion Testing
**Source:** rule-library/RULE-REGISTRY.md#R035
**Criticality:** MANDATORY - Required for approval

MANDATORY QUALITY GATES:
✅ Size Compliance: ≤800 lines (line-counter.sh)
✅ Test Coverage: ≥90% unit tests, integration tests present
✅ KCP Compliance: Multi-tenancy, APIExport, workspace
✅ Security: Input validation, RBAC, secret management
✅ Performance: Latency, resource usage within limits
✅ Documentation: Complete work logs, API documentation
---

```python
def execute_validation_suite(implementation_data):
    """Execute comprehensive validation suite"""
    
    validation_suite = {
        'size_compliance': validate_size_compliance_gate(implementation_data),
        'test_coverage': validate_test_coverage_gate(implementation_data),
        'kcp_compliance': validate_kcp_compliance_gate(implementation_data),
        'security_compliance': validate_security_gate(implementation_data),
        'performance_compliance': validate_performance_gate(implementation_data),
        'documentation_compliance': validate_documentation_gate(implementation_data),
        'integration_readiness': validate_integration_readiness_gate(implementation_data)
    }
    
    # Calculate overall validation score
    validation_results = calculate_validation_results(validation_suite)
    
    # Identify critical failures
    critical_failures = identify_critical_validation_failures(validation_suite)
    
    # Make final decision
    final_decision = make_validation_decision(validation_results, critical_failures)
    
    return {
        'validation_suite_results': validation_suite,
        'validation_score': validation_results['overall_score'],
        'critical_failures': critical_failures,
        'final_decision': final_decision,
        'validation_timestamp': datetime.now().isoformat()
    }

def validate_size_compliance_gate(implementation_data):
    """Validate size compliance quality gate (CRITICAL)"""
    
    branch = implementation_data.get('branch')
    if not branch:
        return {
            'gate_passed': False,
            'error': 'No branch specified for size measurement',
            'critical_failure': True
        }
    
    try:
        # MANDATORY: Use only the approved line counter
        result = subprocess.run([
            '$PROJECT_ROOT/tools/line-counter.sh',
            '-c', branch
        ], capture_output=True, text=True, check=True)
        
        # Parse line count
        output_lines = result.stdout.strip().split('\n')
        line_count = int(output_lines[-1].split()[-1])
        
        gate_result = {
            'gate_passed': line_count <= 800,
            'actual_lines': line_count,
            'limit': 800,
            'margin': 800 - line_count,
            'tool_used': 'line-counter.sh',
            'raw_output': result.stdout.strip(),
            'measurement_timestamp': datetime.now().isoformat()
        }
        
        if not gate_result['gate_passed']:
            gate_result['critical_failure'] = True
            gate_result['failure_reason'] = f"Size {line_count} exceeds 800-line limit"
            gate_result['required_action'] = 'SPLIT_REQUIRED'
        
        return gate_result
        
    except Exception as e:
        return {
            'gate_passed': False,
            'critical_failure': True,
            'error': f"Size measurement failed: {e}",
            'required_action': 'INVESTIGATE_SIZE_MEASUREMENT_FAILURE'
        }

def validate_test_coverage_gate(implementation_data):
    """Validate test coverage quality gate"""
    
    working_dir = implementation_data.get('working_dir')
    if not working_dir or not os.path.exists(working_dir):
        return {
            'gate_passed': False,
            'error': 'Working directory not accessible',
            'critical_failure': True
        }
    
    coverage_validation = {
        'unit_test_coverage': measure_and_validate_unit_coverage(working_dir),
        'integration_test_presence': validate_integration_tests_present(working_dir),
        'multi_tenant_test_coverage': validate_multi_tenant_test_coverage(working_dir),
        'error_case_coverage': validate_error_case_coverage(working_dir)
    }
    
    # Calculate overall test gate status
    unit_coverage_ok = coverage_validation['unit_test_coverage'].get('meets_requirement', False)
    integration_ok = coverage_validation['integration_test_presence'].get('adequate', False)
    multi_tenant_ok = coverage_validation['multi_tenant_test_coverage'].get('adequate', False)
    
    gate_passed = unit_coverage_ok and integration_ok and multi_tenant_ok
    
    gate_result = {
        'gate_passed': gate_passed,
        'coverage_details': coverage_validation,
        'overall_coverage_score': calculate_overall_coverage_score(coverage_validation)
    }
    
    if not gate_passed:
        gate_result['critical_failure'] = True
        gate_result['failure_reasons'] = []
        
        if not unit_coverage_ok:
            gate_result['failure_reasons'].append(
                f"Unit test coverage {coverage_validation['unit_test_coverage'].get('percentage', 0)}% < 90% required"
            )
        if not integration_ok:
            gate_result['failure_reasons'].append("Insufficient integration test coverage")
        if not multi_tenant_ok:
            gate_result['failure_reasons'].append("Inadequate multi-tenant test scenarios")
        
        gate_result['required_action'] = 'IMPROVE_TEST_COVERAGE'
    
    return gate_result

def validate_kcp_compliance_gate(implementation_data):
    """Validate KCP/Kubernetes compliance quality gate"""
    
    working_dir = implementation_data.get('working_dir')
    
    kcp_compliance_checks = {
        'multi_tenancy_implementation': validate_multi_tenancy_implementation(working_dir),
        'logical_cluster_integration': validate_logical_cluster_integration(working_dir),
        'api_export_compliance': validate_api_export_compliance(working_dir),
        'workspace_isolation': validate_workspace_isolation(working_dir),
        'syncer_compatibility': validate_syncer_compatibility(working_dir),
        'rbac_implementation': validate_rbac_implementation(working_dir)
    }
    
    # Calculate KCP compliance score
    compliance_scores = [
        check.get('compliance_score', 0) for check in kcp_compliance_checks.values()
        if isinstance(check, dict) and 'compliance_score' in check
    ]
    
    overall_kcp_score = sum(compliance_scores) / len(compliance_scores) if compliance_scores else 0
    
    # KCP compliance requirements (stricter than planning)
    gate_passed = overall_kcp_score >= 90
    
    gate_result = {
        'gate_passed': gate_passed,
        'overall_kcp_score': overall_kcp_score,
        'compliance_checks': kcp_compliance_checks,
        'kcp_compliance_level': determine_kcp_compliance_level(overall_kcp_score)
    }
    
    if not gate_passed:
        gate_result['critical_failure'] = True
        gate_result['failure_reason'] = f"KCP compliance {overall_kcp_score:.1f}% < 90% required"
        gate_result['required_action'] = 'IMPROVE_KCP_COMPLIANCE'
        
        # Identify specific KCP issues
        gate_result['critical_kcp_issues'] = []
        for check_name, check_result in kcp_compliance_checks.items():
            if isinstance(check_result, dict) and check_result.get('compliance_score', 0) < 85:
                gate_result['critical_kcp_issues'].append({
                    'area': check_name,
                    'score': check_result.get('compliance_score', 0),
                    'issues': check_result.get('issues', [])
                })
    
    return gate_result
```

## KCP-Specific Validation Rules

```go
// KCP Multi-tenancy Validation
func validateMultiTenancyImplementation(workingDir string) ValidationResult {
    result := ValidationResult{
        ComplianceScore: 100,
        Issues: []string{},
    }
    
    // 1. Check for logical cluster awareness
    if !hasLogicalClusterIntegration(workingDir) {
        result.ComplianceScore -= 30
        result.Issues = append(result.Issues, 
            "Missing logical cluster integration - controllers must be cluster-aware")
    }
    
    // 2. Check for workspace isolation
    if !hasWorkspaceIsolation(workingDir) {
        result.ComplianceScore -= 25
        result.Issues = append(result.Issues,
            "Missing workspace isolation - must validate workspace access")
    }
    
    // 3. Check for tenant-scoped operations
    if !hasTenantScopedOperations(workingDir) {
        result.ComplianceScore -= 20
        result.Issues = append(result.Issues,
            "Operations not properly tenant-scoped - risk of cross-tenant access")
    }
    
    // 4. Check for multi-tenant error handling
    if !hasMultiTenantErrorHandling(workingDir) {
        result.ComplianceScore -= 15
        result.Issues = append(result.Issues,
            "Missing multi-tenant error handling - should silently skip unauthorized workspaces")
    }
    
    // 5. Check for tenant context propagation
    if !hasTenantContextPropagation(workingDir) {
        result.ComplianceScore -= 10
        result.Issues = append(result.Issues,
            "Tenant context not properly propagated through call stack")
    }
    
    return result
}

// APIExport Integration Validation
func validateAPIExportCompliance(workingDir string) ValidationResult {
    result := ValidationResult{
        ComplianceScore: 100,
        Issues: []string{},
    }
    
    // Check for APIExport client usage
    apiExportUsage := analyzeAPIExportUsage(workingDir)
    
    if apiExportUsage.Used && !apiExportUsage.ProperlyIntegrated {
        result.ComplianceScore -= 40
        result.Issues = append(result.Issues,
            "APIExport used but not properly integrated - missing client initialization or lifecycle management")
    }
    
    // Check for export discovery patterns
    if apiExportUsage.Used && !hasExportDiscovery(workingDir) {
        result.ComplianceScore -= 30
        result.Issues = append(result.Issues,
            "Missing APIExport discovery - should watch for export changes")
    }
    
    // Check for export lifecycle handling
    if apiExportUsage.Used && !hasExportLifecycleHandling(workingDir) {
        result.ComplianceScore -= 30
        result.Issues = append(result.Issues,
            "Missing APIExport lifecycle handling - should handle export creation/deletion")
    }
    
    return result
}
```

## Performance Validation

```python
def validate_performance_gate(implementation_data):
    """Validate performance quality gate"""
    
    working_dir = implementation_data.get('working_dir')
    performance_requirements = implementation_data.get('performance_requirements', {})
    
    performance_validation = {
        'latency_performance': validate_latency_requirements(working_dir, performance_requirements),
        'resource_usage': validate_resource_usage_requirements(working_dir, performance_requirements),
        'scalability': validate_scalability_requirements(working_dir, performance_requirements),
        'throughput': validate_throughput_requirements(working_dir, performance_requirements)
    }
    
    # Check if performance tests exist and pass
    performance_tests_result = run_performance_tests(working_dir)
    
    gate_passed = all([
        performance_validation['latency_performance'].get('meets_requirements', False),
        performance_validation['resource_usage'].get('meets_requirements', False),
        performance_tests_result.get('all_passed', False)
    ])
    
    gate_result = {
        'gate_passed': gate_passed,
        'performance_validation': performance_validation,
        'performance_tests_result': performance_tests_result,
        'overall_performance_score': calculate_performance_score(performance_validation)
    }
    
    if not gate_passed:
        gate_result['failure_reasons'] = []
        
        if not performance_validation['latency_performance'].get('meets_requirements'):
            gate_result['failure_reasons'].append(
                f"Latency requirements not met: {performance_validation['latency_performance'].get('actual_latency')} > {performance_requirements.get('max_latency', 'unknown')}"
            )
        
        if not performance_validation['resource_usage'].get('meets_requirements'):
            gate_result['failure_reasons'].append(
                f"Resource usage requirements not met: {performance_validation['resource_usage'].get('issues', [])}"
            )
        
        if not performance_tests_result.get('all_passed'):
            gate_result['failure_reasons'].append(
                f"Performance tests failed: {performance_tests_result.get('failed_tests', [])}"
            )
        
        gate_result['required_action'] = 'OPTIMIZE_PERFORMANCE'
    
    return gate_result

def run_performance_tests(working_dir):
    """Run performance tests and analyze results"""
    
    try:
        # Look for performance test files
        perf_test_files = glob.glob(f"{working_dir}/**/*performance*test*.go", recursive=True)
        perf_test_files.extend(glob.glob(f"{working_dir}/**/*load*test*.go", recursive=True))
        perf_test_files.extend(glob.glob(f"{working_dir}/**/*bench*test*.go", recursive=True))
        
        if not perf_test_files:
            return {
                'tests_found': False,
                'all_passed': False,
                'message': 'No performance tests found'
            }
        
        # Run performance tests
        result = subprocess.run([
            'go', 'test', '-bench=.', '-benchmem', './...'
        ], cwd=working_dir, capture_output=True, text=True, timeout=300)
        
        # Parse benchmark results
        benchmark_results = parse_benchmark_results(result.stdout)
        
        # Analyze results against requirements
        performance_analysis = analyze_benchmark_performance(benchmark_results)
        
        return {
            'tests_found': True,
            'all_passed': result.returncode == 0 and performance_analysis['meets_requirements'],
            'benchmark_results': benchmark_results,
            'performance_analysis': performance_analysis,
            'raw_output': result.stdout
        }
        
    except subprocess.TimeoutExpired:
        return {
            'tests_found': True,
            'all_passed': False,
            'message': 'Performance tests timed out (>5 minutes)'
        }
    except Exception as e:
        return {
            'tests_found': False,
            'all_passed': False,
            'error': str(e)
        }
```

## Security Validation

---
### ℹ️ RULE R038.0.0 - Security Review
**Source:** rule-library/RULE-REGISTRY.md#R038
**Criticality:** INFO - Best practice

SECURITY VALIDATION CHECKLIST:
✅ Input Validation: All external inputs validated
✅ Authentication: Proper identity verification
✅ Authorization: RBAC correctly implemented
✅ Data Protection: Sensitive data properly handled
✅ Error Handling: No information leakage
✅ Logging: Security events properly logged
---

```python
def validate_security_gate(implementation_data):
    """Validate security quality gate"""
    
    working_dir = implementation_data.get('working_dir')
    
    security_validation = {
        'input_validation': validate_input_validation_security(working_dir),
        'authentication_security': validate_authentication_implementation(working_dir),
        'authorization_security': validate_authorization_implementation(working_dir),
        'data_protection': validate_data_protection_implementation(working_dir),
        'error_handling_security': validate_error_handling_security(working_dir),
        'logging_security': validate_logging_security(working_dir),
        'secret_management': validate_secret_management_security(working_dir)
    }
    
    # Calculate overall security score
    security_scores = [
        check.get('security_score', 0) for check in security_validation.values()
        if isinstance(check, dict) and 'security_score' in check
    ]
    
    overall_security_score = sum(security_scores) / len(security_scores) if security_scores else 0
    
    # Identify critical security issues
    critical_security_issues = []
    for check_name, check_result in security_validation.items():
        if isinstance(check_result, dict):
            issues = check_result.get('critical_issues', [])
            if issues:
                critical_security_issues.extend([
                    {'area': check_name, 'issue': issue} for issue in issues
                ])
    
    gate_passed = overall_security_score >= 85 and len(critical_security_issues) == 0
    
    gate_result = {
        'gate_passed': gate_passed,
        'overall_security_score': overall_security_score,
        'security_validation': security_validation,
        'critical_security_issues': critical_security_issues
    }
    
    if not gate_passed:
        gate_result['failure_reasons'] = []
        
        if overall_security_score < 85:
            gate_result['failure_reasons'].append(
                f"Overall security score {overall_security_score:.1f}% < 85% required"
            )
        
        if critical_security_issues:
            gate_result['failure_reasons'].append(
                f"{len(critical_security_issues)} critical security issues found"
            )
        
        gate_result['required_action'] = 'FIX_SECURITY_ISSUES'
    
    return gate_result

def validate_input_validation_security(working_dir):
    """Validate input validation security practices"""
    
    # Scan source files for input validation patterns
    source_files = glob.glob(f"{working_dir}/**/*.go", recursive=True)
    
    validation_analysis = {
        'input_sources_found': [],
        'validated_inputs': [],
        'unvalidated_inputs': [],
        'validation_patterns_used': [],
        'security_score': 100,
        'critical_issues': []
    }
    
    for file_path in source_files:
        try:
            with open(file_path, 'r') as f:
                content = f.read()
            
            # Look for input sources (HTTP requests, CLI args, config files)
            input_sources = identify_input_sources(content, file_path)
            validation_analysis['input_sources_found'].extend(input_sources)
            
            # Check for validation patterns
            validation_patterns = identify_validation_patterns(content, file_path)
            validation_analysis['validation_patterns_used'].extend(validation_patterns)
            
            # Match inputs with validations
            unvalidated = find_unvalidated_inputs(input_sources, validation_patterns)
            validation_analysis['unvalidated_inputs'].extend(unvalidated)
            
        except Exception:
            continue
    
    # Calculate security score based on validation coverage
    total_inputs = len(validation_analysis['input_sources_found'])
    unvalidated_count = len(validation_analysis['unvalidated_inputs'])
    
    if total_inputs > 0:
        validation_coverage = ((total_inputs - unvalidated_count) / total_inputs) * 100
        validation_analysis['security_score'] = validation_coverage
    
    # Identify critical issues
    if unvalidated_count > 0:
        validation_analysis['critical_issues'].extend([
            f"Unvalidated input: {input_info}" for input_info in validation_analysis['unvalidated_inputs']
        ])
    
    return validation_analysis
```

## Integration Readiness Validation

```python
def validate_integration_readiness_gate(implementation_data):
    """Validate readiness for integration with other components"""
    
    integration_validation = {
        'interface_compliance': validate_interface_compliance(implementation_data),
        'dependency_satisfaction': validate_dependency_satisfaction(implementation_data),
        'api_stability': validate_api_stability(implementation_data),
        'backward_compatibility': validate_backward_compatibility(implementation_data),
        'configuration_completeness': validate_configuration_completeness(implementation_data),
        'deployment_readiness': validate_deployment_readiness(implementation_data)
    }
    
    # Calculate integration readiness score
    integration_scores = [
        check.get('readiness_score', 0) for check in integration_validation.values()
        if isinstance(check, dict) and 'readiness_score' in check
    ]
    
    overall_readiness_score = sum(integration_scores) / len(integration_scores) if integration_scores else 0
    
    gate_passed = overall_readiness_score >= 80
    
    gate_result = {
        'gate_passed': gate_passed,
        'overall_readiness_score': overall_readiness_score,
        'integration_validation': integration_validation
    }
    
    if not gate_passed:
        gate_result['failure_reason'] = f"Integration readiness {overall_readiness_score:.1f}% < 80% required"
        gate_result['required_action'] = 'IMPROVE_INTEGRATION_READINESS'
        
        # Identify specific readiness issues
        gate_result['readiness_issues'] = []
        for check_name, check_result in integration_validation.items():
            if isinstance(check_result, dict) and check_result.get('readiness_score', 0) < 75:
                gate_result['readiness_issues'].append({
                    'area': check_name,
                    'score': check_result.get('readiness_score', 0),
                    'issues': check_result.get('issues', [])
                })
    
    return gate_result

def make_validation_decision(validation_results, critical_failures):
    """Make final validation decision based on all quality gates"""
    
    decision_factors = {
        'all_gates_passed': len(critical_failures) == 0,
        'overall_score': validation_results['overall_score'],
        'critical_failure_count': len(critical_failures),
        'score_threshold_met': validation_results['overall_score'] >= 85
    }
    
    # Decision logic
    if decision_factors['all_gates_passed'] and decision_factors['score_threshold_met']:
        if validation_results['overall_score'] >= 95:
            decision = 'APPROVED_EXCELLENT'
        elif validation_results['overall_score'] >= 90:
            decision = 'APPROVED'
        else:
            decision = 'APPROVED_WITH_RECOMMENDATIONS'
    else:
        if decision_factors['critical_failure_count'] > 0:
            decision = 'REJECTED_CRITICAL_FAILURES'
        else:
            decision = 'REJECTED_INSUFFICIENT_QUALITY'
    
    return {
        'decision': decision,
        'decision_factors': decision_factors,
        'can_proceed': decision.startswith('APPROVED'),
        'requires_fixes': decision.startswith('REJECTED'),
        'decision_confidence': calculate_decision_confidence(validation_results, critical_failures),
        'decision_timestamp': datetime.now().isoformat()
    }
```

## Validation Reporting

```yaml
# Comprehensive Validation Report Template
validation_report:
  implementation_id: "phase1-wave2-effort3-webhooks"
  validated_at: "2025-08-23T21:00:00Z"
  validator: "code-reviewer-agent"
  
  quality_gates_results:
    size_compliance:
      gate_passed: true
      actual_lines: 687
      limit: 800
      margin: 113
      tool_used: "line-counter.sh"
      
    test_coverage:
      gate_passed: true
      unit_coverage: 92.3
      integration_tests: 8
      multi_tenant_scenarios: 5
      overall_coverage_score: 94
      
    kcp_compliance:
      gate_passed: true
      overall_kcp_score: 91.2
      multi_tenancy_score: 94
      api_export_integration: 89
      workspace_isolation: 92
      
    security:
      gate_passed: true
      overall_security_score: 88
      critical_security_issues: 0
      input_validation: "ADEQUATE"
      rbac_implementation: "PROPER"
      
    performance:
      gate_passed: true
      latency_requirements: "MET"
      resource_usage: "WITHIN_LIMITS"
      scalability: "GOOD"
      
    integration_readiness:
      gate_passed: true
      overall_readiness_score: 87
      interface_compliance: "GOOD"
      api_stability: "STABLE"
  
  overall_validation:
    all_gates_passed: true
    overall_score: 90.4
    critical_failures: 0
    
  final_decision:
    result: "APPROVED"
    can_proceed: true
    confidence_level: 94
    
    recommendations:
      - "Consider adding more error handling edge case tests"
      - "Monitor performance under sustained load in production"
      - "Add metrics collection for webhook operations"
      
  next_steps:
    - "Mark implementation as approved and ready for integration"
    - "Notify orchestrator of successful validation completion"
    - "Update work log with validation results and approval"
```

## State Transitions

From VALIDATION state:
- **APPROVED** → Complete (Implementation approved and ready)
- **APPROVED_WITH_RECOMMENDATIONS** → Complete (Approved with improvement suggestions)
- **REJECTED_CRITICAL_FAILURES** → SPAWN_AGENTS (Spawn SW Engineer for fixes)
- **REJECTED_INSUFFICIENT_QUALITY** → SPAWN_AGENTS (Spawn SW Engineer for improvements)
