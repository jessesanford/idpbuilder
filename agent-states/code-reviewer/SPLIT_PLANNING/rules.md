# Code Reviewer - SPLIT_PLANNING State Rules

## State Context
You are planning how to split an oversized implementation (>800 lines) into compliant sub-efforts while maintaining functionality and architectural integrity.

---
### 🚨 RULE R056.0.0 - Split Plan Creation
**Source:** rule-library/RULE-REGISTRY.md#R056
**Criticality:** CRITICAL - Major impact on grading

SPLIT PLANNING PROTOCOL:
1. Analyze oversized implementation structure
2. Identify logical split boundaries
3. Design split strategy preserving functionality
4. Plan integration sequence between splits
5. Validate each split will be <800 lines
6. Create detailed split execution plan
---

## Size Violation Analysis

---
### 🚨🚨 RULE R007.0.0 - Size Limit Enforcement
**Source:** rule-library/RULE-REGISTRY.md#R007
**Criticality:** MANDATORY - Required for approval

MANDATORY SPLIT TRIGGERS:
1. Implementation >800 lines (measured by tmc-pr-line-counter)
2. Effort estimated >800 lines during planning
3. Implementation approaching 800 lines (>750)
4. Complex integration requiring size management
---

```python
def analyze_size_violation(implementation_data):
    """Analyze why implementation exceeded size limits"""
    
    size_analysis = {
        'current_size': implementation_data.get('actual_lines', 0),
        'size_limit': 800,
        'overage': implementation_data.get('actual_lines', 0) - 800,
        'overage_percentage': ((implementation_data.get('actual_lines', 0) - 800) / 800) * 100,
        'violation_severity': determine_violation_severity(implementation_data.get('actual_lines', 0))
    }
    
    # Analyze size distribution by component
    component_analysis = analyze_component_sizes(implementation_data)
    
    # Identify root causes of size violation
    root_causes = identify_size_violation_causes(implementation_data, component_analysis)
    
    # Assess split complexity
    split_complexity = assess_split_complexity(implementation_data, component_analysis)
    
    return {
        'size_analysis': size_analysis,
        'component_analysis': component_analysis,
        'root_causes': root_causes,
        'split_complexity': split_complexity,
        'split_required': size_analysis['current_size'] > 800,
        'split_urgency': determine_split_urgency(size_analysis['overage'])
    }

def determine_violation_severity(actual_lines):
    """Determine severity of size violation"""
    
    if actual_lines <= 800:
        return 'NONE'
    elif actual_lines <= 1000:
        return 'MODERATE'  # 0-25% overage
    elif actual_lines <= 1200:
        return 'HIGH'      # 25-50% overage  
    else:
        return 'SEVERE'    # >50% overage

def analyze_component_sizes(implementation_data):
    """Analyze size distribution across implementation components"""
    
    # Parse detailed size breakdown from line counter output
    size_breakdown = parse_line_counter_output(
        implementation_data.get('size_compliance', {}).get('raw_output', '')
    )
    
    component_analysis = {}
    
    for file_path, line_count in size_breakdown.items():
        component = classify_file_component(file_path)
        
        if component not in component_analysis:
            component_analysis[component] = {
                'files': [],
                'total_lines': 0,
                'percentage_of_total': 0
            }
        
        component_analysis[component]['files'].append({
            'file': file_path,
            'lines': line_count
        })
        component_analysis[component]['total_lines'] += line_count
    
    # Calculate percentages
    total_lines = implementation_data.get('actual_lines', 0)
    for component_data in component_analysis.values():
        component_data['percentage_of_total'] = (
            component_data['total_lines'] / total_lines * 100 if total_lines > 0 else 0
        )
    
    return component_analysis

def classify_file_component(file_path):
    """Classify file into logical component based on path"""
    
    path_lower = file_path.lower()
    
    if 'api' in path_lower or 'types' in path_lower:
        return 'API_TYPES'
    elif 'controller' in path_lower or 'reconcile' in path_lower:
        return 'CONTROLLERS'
    elif 'webhook' in path_lower:
        return 'WEBHOOKS'
    elif 'client' in path_lower:
        return 'CLIENTS'
    elif 'server' in path_lower or 'handler' in path_lower:
        return 'SERVER_COMPONENTS'
    elif 'config' in path_lower or 'settings' in path_lower:
        return 'CONFIGURATION'
    elif 'util' in path_lower or 'helper' in path_lower:
        return 'UTILITIES'
    elif 'test' in path_lower:
        return 'TESTS'
    else:
        return 'OTHER'
```

## Split Strategy Design

```python
def design_split_strategy(size_analysis, architectural_constraints):
    """Design strategy for splitting oversized implementation"""
    
    # Identify potential split boundaries
    split_boundaries = identify_split_boundaries(
        size_analysis['component_analysis'],
        architectural_constraints
    )
    
    # Design split options
    split_options = generate_split_options(split_boundaries, size_analysis)
    
    # Evaluate split options
    evaluated_options = evaluate_split_options(split_options, architectural_constraints)
    
    # Select optimal split strategy
    optimal_strategy = select_optimal_split_strategy(evaluated_options)
    
    return {
        'split_boundaries': split_boundaries,
        'split_options': split_options,
        'evaluated_options': evaluated_options,
        'recommended_strategy': optimal_strategy,
        'split_count': len(optimal_strategy['splits']),
        'strategy_confidence': optimal_strategy['confidence_score']
    }

def identify_split_boundaries(component_analysis, constraints):
    """Identify logical boundaries where implementation can be split"""
    
    boundaries = []
    
    # Component-based boundaries (natural splits)
    for component, data in component_analysis.items():
        if data['total_lines'] >= 200:  # Minimum viable split size
            boundaries.append({
                'type': 'COMPONENT_BOUNDARY',
                'component': component,
                'size': data['total_lines'],
                'files': data['files'],
                'split_feasibility': assess_component_split_feasibility(component, data, constraints)
            })
    
    # Functional boundaries (within components)
    functional_boundaries = identify_functional_boundaries(component_analysis, constraints)
    boundaries.extend(functional_boundaries)
    
    # Interface boundaries (API/service boundaries)
    interface_boundaries = identify_interface_boundaries(component_analysis, constraints)
    boundaries.extend(interface_boundaries)
    
    # Layer boundaries (presentation, business, data)
    layer_boundaries = identify_layer_boundaries(component_analysis, constraints)
    boundaries.extend(layer_boundaries)
    
    return boundaries

def generate_split_options(boundaries, size_analysis):
    """Generate different split configuration options"""
    
    total_size = size_analysis['size_analysis']['current_size']
    target_split_size = 700  # Target size per split (buffer below 800)
    min_splits = math.ceil(total_size / target_split_size)
    
    split_options = []
    
    # Option 1: Component-based split
    component_split = design_component_based_split(boundaries, target_split_size)
    if component_split:
        split_options.append({
            'strategy': 'COMPONENT_BASED',
            'description': 'Split along component boundaries',
            'splits': component_split,
            'split_count': len(component_split)
        })
    
    # Option 2: Functional split
    functional_split = design_functional_split(boundaries, target_split_size)
    if functional_split:
        split_options.append({
            'strategy': 'FUNCTIONAL',
            'description': 'Split along functional boundaries', 
            'splits': functional_split,
            'split_count': len(functional_split)
        })
    
    # Option 3: Layered split
    layered_split = design_layered_split(boundaries, target_split_size)
    if layered_split:
        split_options.append({
            'strategy': 'LAYERED',
            'description': 'Split along architectural layers',
            'splits': layered_split,
            'split_count': len(layered_split)
        })
    
    # Option 4: Hybrid split (combination approach)
    hybrid_split = design_hybrid_split(boundaries, target_split_size, min_splits)
    if hybrid_split:
        split_options.append({
            'strategy': 'HYBRID',
            'description': 'Combination of component and functional splits',
            'splits': hybrid_split,
            'split_count': len(hybrid_split)
        })
    
    return split_options
```

## KCP-Specific Split Considerations

---
### ℹ️ RULE R037.0.0 - Pattern Compliance
**Source:** rule-library/RULE-REGISTRY.md#R037
**Criticality:** INFO - Best practice

KCP SPLIT REQUIREMENTS:
1. Preserve logical cluster isolation across splits
2. Maintain APIExport integration integrity
3. Ensure multi-tenancy works across split boundaries
4. Preserve syncer compatibility in all splits
5. Maintain workspace isolation mechanisms
---

```go
// KCP-specific split validation
func validateKCPSplitCompliance(splitPlan SplitPlan) []KCPSplitIssue {
    issues := []KCPSplitIssue{}
    
    // 1. Validate logical cluster context preservation
    for _, split := range splitPlan.Splits {
        if !validateLogicalClusterContextPreservation(split) {
            issues = append(issues, KCPSplitIssue{
                Type: "LOGICAL_CLUSTER_CONTEXT_BROKEN",
                Split: split.ID,
                Severity: "CRITICAL",
                Description: "Split breaks logical cluster context flow",
                Recommendation: "Ensure logical cluster context passes between splits"
            })
        }
    }
    
    // 2. Validate APIExport integration integrity
    apiExportIssues := validateAPIExportIntegration(splitPlan)
    issues = append(issues, apiExportIssues...)
    
    // 3. Validate multi-tenant isolation preserved
    isolationIssues := validateMultiTenantIsolation(splitPlan)
    issues = append(issues, isolationIssues...)
    
    // 4. Validate syncer compatibility maintained
    syncerIssues := validateSyncerCompatibility(splitPlan)
    issues = append(issues, syncerIssues...)
    
    return issues
}

func validateAPIExportIntegration(splitPlan SplitPlan) []KCPSplitIssue {
    issues := []KCPSplitIssue{}
    
    // Check if APIExport components are properly distributed
    apiExportComponents := identifyAPIExportComponents(splitPlan)
    
    if len(apiExportComponents) > 1 {
        // APIExport logic split across multiple components
        
        // Validate integration points
        for i, component := range apiExportComponents {
            for j, otherComponent := range apiExportComponents {
                if i != j {
                    integration := validateAPIExportIntegration(component, otherComponent)
                    if !integration.Valid {
                        issues = append(issues, KCPSplitIssue{
                            Type: "APIEXPORT_INTEGRATION_BROKEN",
                            Splits: []string{component.SplitID, otherComponent.SplitID},
                            Severity: "HIGH",
                            Description: fmt.Sprintf("APIExport integration broken between %s and %s", 
                                component.SplitID, otherComponent.SplitID),
                            Recommendation: "Add integration interfaces or merge APIExport logic"
                        })
                    }
                }
            }
        }
    }
    
    return issues
}
```

## Split Validation Framework

```python
def validate_split_plan(split_strategy):
    """Comprehensively validate split plan before execution"""
    
    validation_results = {
        'size_compliance': validate_split_sizes(split_strategy),
        'functional_integrity': validate_functional_integrity(split_strategy),
        'architectural_soundness': validate_architectural_soundness(split_strategy),
        'kcp_compliance': validate_kcp_split_compliance(split_strategy),
        'integration_feasibility': validate_integration_feasibility(split_strategy),
        'test_coverage_preservation': validate_test_coverage_preservation(split_strategy)
    }
    
    # Calculate overall validation score
    validation_score = calculate_split_validation_score(validation_results)
    
    # Identify blocking issues
    blocking_issues = identify_blocking_split_issues(validation_results)
    
    return {
        'validation_results': validation_results,
        'validation_score': validation_score,
        'blocking_issues': blocking_issues,
        'plan_viable': len(blocking_issues) == 0,
        'confidence_level': calculate_split_confidence(validation_results)
    }

def validate_split_sizes(split_strategy):
    """Validate each split meets size requirements"""
    
    size_validation = {
        'all_splits_compliant': True,
        'split_sizes': [],
        'oversized_splits': [],
        'undersized_splits': []
    }
    
    for split in split_strategy['splits']:
        estimated_size = split.get('estimated_lines', 0)
        
        size_validation['split_sizes'].append({
            'split_id': split['id'],
            'estimated_lines': estimated_size,
            'compliant': estimated_size <= 800,
            'margin': 800 - estimated_size
        })
        
        if estimated_size > 800:
            size_validation['all_splits_compliant'] = False
            size_validation['oversized_splits'].append({
                'split_id': split['id'],
                'size': estimated_size,
                'overage': estimated_size - 800
            })
        elif estimated_size < 100:  # Too small to be viable
            size_validation['undersized_splits'].append({
                'split_id': split['id'],
                'size': estimated_size,
                'recommendation': 'Consider merging with another split'
            })
    
    return size_validation

def validate_functional_integrity(split_strategy):
    """Validate splits preserve functional integrity"""
    
    integrity_checks = {
        'data_flow_preserved': validate_data_flow_preservation(split_strategy),
        'api_contracts_maintained': validate_api_contracts(split_strategy),
        'error_handling_complete': validate_error_handling_completeness(split_strategy),
        'configuration_accessible': validate_configuration_access(split_strategy)
    }
    
    # Check for functional dependencies between splits
    dependency_issues = analyze_split_dependencies(split_strategy)
    
    return {
        'integrity_checks': integrity_checks,
        'dependency_issues': dependency_issues,
        'functional_integrity_score': calculate_functional_integrity_score(integrity_checks),
        'critical_issues': identify_critical_functional_issues(integrity_checks, dependency_issues)
    }
```

## Split Execution Planning

```yaml
# Split Execution Plan Template
split_execution_plan:
  original_effort: "phase1-wave2-effort3-webhooks"
  split_trigger: "Size violation: 1247 lines > 800 limit"
  split_strategy: "HYBRID"
  
  # Split configuration
  splits:
    - split_id: "webhooks-core-server"
      description: "Core webhook server and routing infrastructure"
      estimated_lines: 387
      
      components:
        - "pkg/webhooks/server/"
        - "pkg/webhooks/routing/"
        - "pkg/webhooks/config/"
        
      dependencies:
        - external: "KCP APIExport client"
          reason: "Server registration with APIExports"
        - internal: "webhook-admission-logic"
          reason: "Webhook handlers for admission processing"
          
      interfaces_provided:
        - "WebhookServer interface"
        - "TenantRouter interface"
        - "ConfigurationManager interface"
        
      interfaces_required:
        - "AdmissionHandler interface (from admission split)"
        
    - split_id: "webhooks-admission-logic"
      description: "Admission webhook validation and mutation logic"
      estimated_lines: 423
      
      components:
        - "pkg/webhooks/admission/"
        - "pkg/webhooks/validation/"
        - "pkg/webhooks/mutation/"
        
      dependencies:
        - external: "KCP logical cluster client"
          reason: "Multi-tenant admission processing"
        - internal: "webhooks-core-server"
          reason: "Server registration and request routing"
          
      interfaces_provided:
        - "AdmissionHandler interface"
        - "ResourceValidator interface"
        - "ResourceMutator interface"
        
      interfaces_required:
        - "WebhookServer interface (from server split)"
        
    - split_id: "webhooks-integration-tests"
      description: "Integration tests and multi-tenant scenarios"
      estimated_lines: 437
      
      components:
        - "test/webhooks/integration/"
        - "test/webhooks/multitenant/"
        - "test/webhooks/performance/"
        
      dependencies:
        - internal: "webhooks-core-server"
          reason: "Test server functionality"
        - internal: "webhooks-admission-logic"
          reason: "Test admission processing"
          
  # Execution sequence
  execution_sequence:
    phase_1_foundation:
      duration_hours: 6
      splits_to_create: ["webhooks-core-server"]
      deliverables:
        - "Split 1: Core server working directory"
        - "Split 1: Implementation plan"
        - "Split 1: Interfaces defined and documented"
      validation_criteria:
        - "Server split compiles independently"
        - "Interface contracts well-defined"
        - "Size compliance maintained"
        
    phase_2_admission:
      duration_hours: 5
      depends_on: "phase_1_foundation"
      splits_to_create: ["webhooks-admission-logic"]
      deliverables:
        - "Split 2: Admission logic working directory"
        - "Split 2: Implementation plan"
        - "Interface integration with Split 1"
      validation_criteria:
        - "Admission split compiles independently"
        - "Integration interfaces work correctly"
        - "Multi-tenancy preserved"
        
    phase_3_integration:
      duration_hours: 4
      depends_on: "phase_2_admission"
      splits_to_create: ["webhooks-integration-tests"]
      deliverables:
        - "Split 3: Test infrastructure"
        - "Integration between all splits validated"
        - "End-to-end functionality verified"
      validation_criteria:
        - "All splits integrate correctly"
        - "Original functionality preserved"
        - "Test coverage maintained"
        
    phase_4_validation:
      duration_hours: 2
      depends_on: "phase_3_integration"
      deliverables:
        - "Complete split integration validation"
        - "Performance targets verified"
        - "KCP compliance across all splits"
      validation_criteria:
        - "Integration performance acceptable"
        - "No functionality regressions"
        - "Multi-tenant isolation working"
  
  total_estimated_duration: 17  # hours
  
  # Integration strategy
  integration_strategy:
    approach: "Sequential integration with validation gates"
    
    integration_points:
      - between: ["webhooks-core-server", "webhooks-admission-logic"]
        mechanism: "Interface-based dependency injection"
        validation: "Unit tests for interface contracts"
        
      - between: ["webhooks-core-server", "webhooks-integration-tests"]
        mechanism: "Test harness setup and teardown"
        validation: "Integration test execution"
        
      - between: ["webhooks-admission-logic", "webhooks-integration-tests"]
        mechanism: "Test scenario execution"
        validation: "Multi-tenant test scenarios"
    
    validation_gates:
      - gate: "Split 1 Complete"
        criteria: ["Server functionality working", "Interfaces documented", "Size compliant"]
        
      - gate: "Split 2 Complete"
        criteria: ["Admission logic working", "Integration working", "Size compliant"]
        
      - gate: "All Splits Integrated"
        criteria: ["End-to-end functionality", "Performance targets met", "KCP compliance"]
  
  # Risk mitigation
  risk_mitigation:
    high_risks:
      - risk: "Interface integration complexity"
        mitigation: 
          - "Define interfaces clearly before implementation"
          - "Create integration tests early"
          - "Validate interfaces with mock implementations"
          
      - risk: "Multi-tenant functionality broken across splits"
        mitigation:
          - "Preserve logical cluster context in all interfaces"
          - "Test multi-tenant scenarios with each split integration"
          - "Validate workspace isolation at integration points"
          
      - risk: "Performance degradation from split overhead"
        mitigation:
          - "Measure performance at each integration point"
          - "Optimize interface calls if needed"
          - "Consider inlining critical paths if performance issues"
    
    medium_risks:
      - risk: "Configuration management complexity"
        mitigation:
          - "Centralize configuration in server split"
          - "Use dependency injection for configuration access"
          - "Test configuration changes across all splits"
```

## Split Plan Validation

```python
def create_split_validation_checklist(split_plan):
    """Create validation checklist for split execution"""
    
    checklist = {
        'pre_split_validation': [
            'Original implementation size >800 lines confirmed',
            'Split strategy selected and validated',
            'Component boundaries clearly defined',
            'Interface contracts designed and documented',
            'Integration sequence planned',
            'Risk mitigation strategies prepared'
        ],
        
        'during_split_validation': [
            'Each split stays under 800 lines',
            'Interface contracts implemented correctly',
            'Integration points working as planned',
            'Test coverage preserved across splits',
            'KCP compliance maintained in all splits',
            'Multi-tenancy working across split boundaries'
        ],
        
        'post_split_validation': [
            'All splits integrate correctly',
            'Original functionality fully preserved',
            'Performance targets met',
            'Security model maintained',
            'Documentation updated for split architecture',
            'Future maintainability assessed'
        ]
    }
    
    return checklist
```

## State Transitions

From SPLIT_PLANNING state:
- **SPLIT_PLAN_APPROVED** → SPAWN_AGENTS (Spawn SW Engineer for first split)
- **PLAN_NEEDS_REVISION** → SPLIT_PLANNING (Revise split strategy)
- **INTEGRATION_TOO_COMPLEX** → ERROR_RECOVERY (Reconsider implementation approach)
