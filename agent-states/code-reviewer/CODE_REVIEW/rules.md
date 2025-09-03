# Code Reviewer - CODE_REVIEW State Rules

## State Context
You are reviewing a completed implementation, validating code quality, architecture compliance, and test coverage before approval.

---
### ℹ️ RULE R108.0.0 - CODE_REVIEW Rules
**Source:** rule-library/RULE-REGISTRY.md#R108
**Criticality:** INFO - Best practice

CODE REVIEW PROTOCOL:
1. Validate implementation against plan
2. Check size compliance using line counter
3. Verify test coverage requirements
4. Validate KCP/Kubernetes patterns
5. Check multi-tenancy implementation
6. Assess security and performance
7. Provide detailed feedback
---

## 🔴🔴🔴 MANDATORY LINE COUNTING REQUIREMENTS 🔴🔴🔴

### ⚠️⚠️⚠️ CRITICAL: YOU MUST USE LINE-COUNTER.SH TOOL WITH CORRECT BASE BRANCH ⚠️⚠️⚠️

**VIOLATION = -100% IMMEDIATE FAILURE**

### MANDATORY STEPS:
1. **ALWAYS use ${PROJECT_ROOT}/tools/line-counter.sh** - NO EXCEPTIONS
2. **ALWAYS specify -b parameter with CORRECT BASE BRANCH**
3. **NEVER do manual counting** - AUTOMATIC FAILURE (-100%)
4. **NEVER use "main" as base for efforts** - AUTOMATIC FAILURE (-100%)
5. **NEVER count test/doc files separately** - tool handles this

### CORRECT USAGE:
```bash
# STEP 1: Navigate to effort directory (IT'S A SEPARATE GIT REPO!)
cd /path/to/effort/directory
pwd  # Confirm location
ls -la .git  # MUST exist - this is the effort's own git repository!

# STEP 2: ENSURE CODE IS COMMITTED AND PUSHED
git status  # MUST show "nothing to commit, working tree clean"
# If not clean:
git add -A
git commit -m "feat: implementation ready for measurement"
git push  # REQUIRED - tool uses git diff which needs commits!

# STEP 3: Get ACTUAL BRANCH NAMES (not directory names!)
CURRENT_BRANCH=$(git branch --show-current)
echo "Current branch: $CURRENT_BRANCH"  # e.g., phase2-wave1-gcr-image-builder

# STEP 4: Find the BASE branch IN THIS REPOSITORY
git branch -a | grep -E "integration|main"
# Use what exists, e.g., phase2/integration or main
BASE_BRANCH="phase2/integration"  # NOT always "main"!

# STEP 5: Find project root (where orchestrator-state.yaml lives)
PROJECT_ROOT=$(pwd); while [ "$PROJECT_ROOT" != "/" ]; do 
    [ -f "$PROJECT_ROOT/orchestrator-state.yaml" ] && break; 
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT"); 
done

# STEP 6: Run the tool WITH BRANCH NAMES FROM THIS REPO
$PROJECT_ROOT/tools/line-counter.sh -b $BASE_BRANCH -c $CURRENT_BRANCH
```

### 🔴🔴🔴 CRITICAL: Directory Names vs Branch Names 🔴🔴🔴

**THE FATAL MISTAKE:**
```bash
# ❌❌❌ WRONG - Using directory name as branch:
cd /workspaces/project
./tools/line-counter.sh -b main -c go-containerregistry-image-builder
# "go-containerregistry-image-builder" is a DIRECTORY, not a branch!
# Git can't find this "branch" = measurement fails!

# ✅✅✅ RIGHT - Using actual branch name:
cd efforts/phase2/wave1/go-containerregistry-image-builder
BRANCH=$(git branch --show-current)  # Gets: phase2-wave1-gcr-image-builder
../../tools/line-counter.sh -b phase2/integration -c "$BRANCH"
```

### BASE BRANCH IDENTIFICATION:
```bash
# Check orchestrator-state.yaml for current_phase_integration.branch
cat orchestrator-state.yaml | grep "current_phase_integration" -A 5
# OR use the phase integration branch naming pattern:
# phase[N]/integration or phase[N]/integration-[timestamp]
```

### FORBIDDEN ACTIONS:
- ❌ Manual line counting (wc -l, etc.)
- ❌ Using "main" as base for effort measurements
- ❌ Counting test files separately
- ❌ Counting documentation files
- ❌ Using old tool locations (/workspaces/kcp-shared-tools/)

## Size Compliance Validation

---
### 🚨🚨 RULE R007.0.0 - Size Limit Enforcement
**Source:** rule-library/RULE-REGISTRY.md#R007
**Criticality:** MANDATORY - Required for approval

MANDATORY SIZE VALIDATION:
1. Use ONLY ${PROJECT_ROOT}/tools/line-counter.sh
2. Specify correct base branch (-b parameter)
3. Exclude generated code (tool handles automatically)
4. Implementation must be ≤800 lines
5. If >800 lines, IMMEDIATE split required
6. Document exact command and output in review
---

```python
def validate_effort_size_compliance(effort_branch, base_branch):
    """Validate effort size using mandatory line counter"""
    
    # CRITICAL: Identify correct base branch first!
    if base_branch == "main":
        raise ValueError("❌ FATAL: Cannot use 'main' as base for effort measurement!")
    
    try:
        # Find project root first
        project_root = find_project_root()
        line_counter = f"{project_root}/tools/line-counter.sh"
        
        # CRITICAL: Use only the approved line counter with CORRECT BASE
        result = subprocess.run([
            line_counter,
            '-b', base_branch,  # MUST be phase integration branch
            '-c', effort_branch
        ], capture_output=True, text=True, check=True)
        
        # Parse line count from output
        output_lines = result.stdout.strip().split('\n')
        line_count = int(output_lines[-1].split()[-1])
        
        compliance_result = {
            'compliant': line_count <= 800,
            'actual_lines': line_count,
            'limit': 800,
            'margin': 800 - line_count,
            'tool_used': './tools/line-counter.sh',
            'base_branch': base_branch,
            'command': f'./tools/line-counter.sh -b {base_branch} -c {effort_branch}',
            'raw_output': result.stdout.strip()
        }
        
        # Critical failure if over limit
        if not compliance_result['compliant']:
            compliance_result['critical_failure'] = True
            compliance_result['required_action'] = 'IMMEDIATE_SPLIT_REQUIRED'
            compliance_result['split_urgency'] = 'CRITICAL'
        
        return compliance_result
        
    except subprocess.CalledProcessError as e:
        return {
            'compliant': False,
            'error': f"Line counter failed: {e}",
            'critical_failure': True,
            'required_action': 'INVESTIGATE_SIZE_CHECK_FAILURE'
        }
    except Exception as e:
        return {
            'compliant': False, 
            'error': f"Size validation error: {e}",
            'critical_failure': True,
            'required_action': 'INVESTIGATE_SIZE_CHECK_FAILURE'
        }

def document_size_measurement(size_result, review_context):
    """Document size measurement results in review"""
    
    review_context['size_compliance'] = {
        'measured_at': datetime.now().isoformat(),
        'tool_used': size_result.get('tool_used', './tools/line-counter.sh'),
        'base_branch': size_result.get('base_branch', 'UNKNOWN'),
        'command_executed': size_result.get('command', 'UNKNOWN'),
        'actual_lines': size_result.get('actual_lines', 0),
        'limit': size_result.get('limit', 800),
        'compliant': size_result.get('compliant', False),
        'raw_measurement': size_result.get('raw_output', '')
    }
    
    if not size_result.get('compliant', False):
        review_context['critical_issues'] = review_context.get('critical_issues', [])
        review_context['critical_issues'].append({
            'type': 'SIZE_VIOLATION',
            'severity': 'CRITICAL',
            'description': f"Implementation {size_result.get('actual_lines', 0)} lines exceeds 800-line limit",
            'required_action': 'Split effort before approval',
            'blocking': True
        })
    
    return review_context
```

## KCP/Kubernetes Pattern Validation

---
### ℹ️ RULE R037.0.0 - Pattern Compliance
**Source:** rule-library/RULE-REGISTRY.md#R037
**Criticality:** INFO - Best practice

KCP PATTERN CHECKLIST:
✅ Multi-tenancy: Logical cluster awareness
✅ APIExport: Proper integration patterns
✅ Virtual Workspace: Compliance with VW model
✅ Syncer: Compatible with syncer patterns
✅ RBAC: Workspace-scoped permissions
✅ Resource Quotas: Tenant isolation enforcement
---

```go
// Code Review Checklist: Multi-tenant Controller Pattern
func reviewMultiTenantController(implementation string) []ReviewIssue {
    issues := []ReviewIssue{}
    
    // 1. Check logical cluster awareness
    if !containsPattern(implementation, "logicalcluster.Name") {
        issues = append(issues, ReviewIssue{
            Type: "MISSING_LOGICAL_CLUSTER",
            Severity: "CRITICAL",
            Description: "Controller must be logical cluster aware",
            Example: "Add LogicalCluster logicalcluster.Name field",
            Pattern: "Multi-tenancy requirement"
        })
    }
    
    // 2. Check workspace isolation
    if !containsPattern(implementation, "workspace.*access.*check") {
        issues = append(issues, ReviewIssue{
            Type: "MISSING_WORKSPACE_ISOLATION", 
            Severity: "CRITICAL",
            Description: "Must validate workspace access before operations",
            Example: "Implement hasWorkspaceAccess() validation",
            Pattern: "Security requirement"
        })
    }
    
    // 3. Check APIExport integration
    if containsPattern(implementation, "apiexport") && 
       !containsPattern(implementation, "APIExportClient") {
        issues = append(issues, ReviewIssue{
            Type: "INCOMPLETE_APIEXPORT_INTEGRATION",
            Severity: "HIGH",
            Description: "APIExport usage requires proper client integration",
            Example: "Add APIExportClient field and initialization",
            Pattern: "APIExport compliance"
        })
    }
    
    // 4. Check error handling for multi-tenancy
    if !containsPattern(implementation, "unauthorized.*workspace") {
        issues = append(issues, ReviewIssue{
            Type: "MISSING_TENANT_ERROR_HANDLING",
            Severity: "MEDIUM", 
            Description: "Should handle unauthorized workspace access gracefully",
            Example: "Return nil error for unauthorized workspaces (silent skip)",
            Pattern: "Multi-tenant error handling"
        })
    }
    
    return issues
}

// Code Review Checklist: KCP API Types
func reviewKCPAPITypes(implementation string) []ReviewIssue {
    issues := []ReviewIssue{}
    
    // Check for proper KCP annotations
    requiredAnnotations := []string{
        "kcp.io/cluster",
        "kcp.io/workspace", 
        "apis.kcp.io/binding"
    }
    
    for _, annotation := range requiredAnnotations {
        if containsAPIUsage(implementation) && 
           !containsPattern(implementation, annotation) {
            issues = append(issues, ReviewIssue{
                Type: "MISSING_KCP_ANNOTATION",
                Severity: "HIGH",
                Description: fmt.Sprintf("API types should consider %s annotation", annotation),
                Example: fmt.Sprintf("Add %s annotation handling", annotation),
                Pattern: "KCP API compliance"
            })
        }
    }
    
    return issues
}
```

## Test Coverage Validation

---
### 🚨🚨 RULE R032.0.0 - Test Coverage Requirements
**Source:** rule-library/RULE-REGISTRY.md#R032
**Criticality:** MANDATORY - Required for approval

MANDATORY COVERAGE VALIDATION:
- Unit Tests: 90% line coverage minimum
- Integration Tests: All API endpoints covered
- Multi-tenant Tests: Cross-workspace scenarios tested
- Error Cases: All error paths validated
- Performance: Resource usage within limits
---

```python
def validate_test_coverage(effort_dir):
    """Validate test coverage meets requirements"""
    
    coverage_results = {
        'unit_test_coverage': measure_unit_test_coverage(effort_dir),
        'integration_test_coverage': assess_integration_tests(effort_dir),
        'multi_tenant_test_coverage': assess_multi_tenant_tests(effort_dir),
        'error_case_coverage': assess_error_case_coverage(effort_dir),
        'performance_test_coverage': assess_performance_tests(effort_dir)
    }
    
    # Calculate overall coverage score
    coverage_score = calculate_coverage_score(coverage_results)
    
    # Determine compliance
    compliance = {
        'meets_requirements': coverage_score >= 90,
        'coverage_score': coverage_score,
        'detailed_results': coverage_results,
        'coverage_gaps': identify_coverage_gaps(coverage_results),
        'critical_issues': []
    }
    
    # Check for critical coverage gaps
    unit_coverage = coverage_results['unit_test_coverage'].get('percentage', 0)
    if unit_coverage < 90:
        compliance['critical_issues'].append({
            'type': 'INSUFFICIENT_UNIT_COVERAGE',
            'severity': 'CRITICAL',
            'description': f"Unit test coverage {unit_coverage}% < 90% required",
            'blocking': True
        })
    
    multi_tenant_tests = coverage_results['multi_tenant_test_coverage'].get('scenarios_covered', 0)
    if multi_tenant_tests == 0:
        compliance['critical_issues'].append({
            'type': 'NO_MULTI_TENANT_TESTS',
            'severity': 'CRITICAL', 
            'description': "No multi-tenant test scenarios found",
            'blocking': True
        })
    
    return compliance

def measure_unit_test_coverage(effort_dir):
    """Measure unit test coverage using go test"""
    
    try:
        # Run tests with coverage
        result = subprocess.run([
            'go', 'test', '-coverprofile=coverage.out', './...'
        ], cwd=effort_dir, capture_output=True, text=True)
        
        # Parse coverage percentage  
        coverage_result = subprocess.run([
            'go', 'tool', 'cover', '-func=coverage.out'
        ], cwd=effort_dir, capture_output=True, text=True)
        
        # Extract overall coverage percentage
        coverage_lines = coverage_result.stdout.strip().split('\n')
        total_line = [line for line in coverage_lines if 'total:' in line]
        
        if total_line:
            percentage_str = total_line[0].split()[-1].rstrip('%')
            percentage = float(percentage_str)
        else:
            percentage = 0.0
        
        return {
            'percentage': percentage,
            'detailed_output': coverage_result.stdout,
            'test_output': result.stdout,
            'success': result.returncode == 0
        }
        
    except Exception as e:
        return {
            'percentage': 0.0,
            'error': str(e),
            'success': False
        }

def assess_multi_tenant_tests(effort_dir):
    """Assess multi-tenant test scenario coverage"""
    
    test_files = glob.glob(f"{effort_dir}/**/*_test.go", recursive=True)
    
    multi_tenant_indicators = [
        'logical.*cluster', 'workspace.*isolation', 'multi.*tenant',
        'tenant.*specific', 'cross.*workspace', 'workspace.*access'
    ]
    
    scenarios_found = []
    total_tests = 0
    
    for test_file in test_files:
        try:
            with open(test_file, 'r') as f:
                content = f.read()
                
            # Count test functions
            test_functions = re.findall(r'func Test\w+', content)
            total_tests += len(test_functions)
            
            # Check for multi-tenant test patterns
            for indicator in multi_tenant_indicators:
                if re.search(indicator, content, re.IGNORECASE):
                    scenarios_found.append({
                        'file': test_file,
                        'pattern': indicator,
                        'context': extract_test_context(content, indicator)
                    })
                    
        except Exception as e:
            continue
    
    return {
        'scenarios_covered': len(scenarios_found),
        'total_tests': total_tests,
        'multi_tenant_ratio': len(scenarios_found) / max(total_tests, 1),
        'scenarios_details': scenarios_found,
        'adequate_coverage': len(scenarios_found) >= 3  # Minimum 3 scenarios
    }
```

## Architecture Review

```python
def review_architecture_compliance(implementation_dir, original_plan):
    """Review implementation against architectural plan"""
    
    review_results = {
        'architecture_compliance': assess_architecture_adherence(implementation_dir, original_plan),
        'design_pattern_usage': validate_design_patterns(implementation_dir),
        'interface_compliance': validate_interface_implementation(implementation_dir, original_plan),
        'component_structure': validate_component_structure(implementation_dir, original_plan)
    }
    
    # Calculate compliance score
    compliance_score = calculate_architecture_compliance_score(review_results)
    
    return {
        'compliance_score': compliance_score,
        'detailed_results': review_results,
        'architecture_issues': identify_architecture_issues(review_results),
        'recommendations': generate_architecture_recommendations(review_results)
    }

def assess_architecture_adherence(implementation_dir, plan):
    """Assess how well implementation follows planned architecture"""
    
    planned_components = plan.get('architecture_design', {}).get('component_structure', {})
    implemented_structure = analyze_actual_structure(implementation_dir)
    
    adherence_results = {}
    
    for component_name, component_plan in planned_components.items():
        actual_impl = implemented_structure.get(component_name, {})
        
        adherence_results[component_name] = {
            'planned_interfaces': component_plan.get('key_interfaces', []),
            'implemented_interfaces': actual_impl.get('interfaces', []),
            'interface_match': calculate_interface_match(
                component_plan.get('key_interfaces', []),
                actual_impl.get('interfaces', [])
            ),
            'size_adherence': assess_size_adherence(
                component_plan.get('estimated_lines', 0),
                actual_impl.get('actual_lines', 0)
            ),
            'structure_match': assess_structure_match(component_plan, actual_impl)
        }
    
    return adherence_results
```

## Security and Performance Review

---
### ℹ️ RULE R038.0.0 - Security Review
**Source:** rule-library/RULE-REGISTRY.md#R038
**Criticality:** INFO - Best practice

SECURITY CHECKLIST:
✅ Input validation on all external data
✅ Workspace isolation properly enforced
✅ RBAC permissions correctly implemented
✅ No hardcoded credentials or secrets
✅ Error messages don't leak sensitive information
✅ Resource access properly authorized
---

```go
// Security Review Patterns
func reviewSecurityPatterns(implementation string) []SecurityIssue {
    issues := []SecurityIssue{}
    
    // Check for input validation
    if containsUserInput(implementation) && !containsValidation(implementation) {
        issues = append(issues, SecurityIssue{
            Type: "MISSING_INPUT_VALIDATION",
            Severity: "HIGH",
            Description: "User input not validated",
            Recommendation: "Add input validation before processing",
            CWE: "CWE-20"
        })
    }
    
    // Check for hardcoded secrets
    secretPatterns := []string{
        `password\s*=\s*"[^"]*"`,
        `token\s*=\s*"[^"]*"`,
        `key\s*=\s*"[^"]*"`,
        `secret\s*=\s*"[^"]*"`
    }
    
    for _, pattern := range secretPatterns {
        if matched := regexp.MustCompile(pattern); matched.MatchString(implementation) {
            issues = append(issues, SecurityIssue{
                Type: "HARDCODED_SECRET",
                Severity: "CRITICAL",
                Description: "Hardcoded secret detected",
                Recommendation: "Use environment variables or secret management",
                CWE: "CWE-798"
            })
        }
    }
    
    return issues
}
```

## Review Decision Framework

```python
def make_review_decision(review_data):
    """Make final review decision based on all validation results"""
    
    # Critical blocking issues
    blocking_issues = []
    
    # Size compliance (CRITICAL)
    size_result = review_data.get('size_compliance', {})
    if not size_result.get('compliant', False):
        blocking_issues.append({
            'type': 'SIZE_VIOLATION',
            'description': f"Size {size_result.get('actual_lines', 0)} > 800 lines",
            'action_required': 'SPLIT_EFFORT'
        })
    
    # Test coverage (CRITICAL)
    coverage_result = review_data.get('test_coverage', {})
    if not coverage_result.get('meets_requirements', False):
        blocking_issues.append({
            'type': 'INSUFFICIENT_COVERAGE',
            'description': f"Coverage {coverage_result.get('coverage_score', 0)}% < 90%",
            'action_required': 'IMPROVE_TESTS'
        })
    
    # KCP compliance (CRITICAL)  
    kcp_result = review_data.get('kcp_compliance', {})
    if kcp_result.get('critical_issues', []):
        blocking_issues.append({
            'type': 'KCP_COMPLIANCE_FAILURE',
            'description': "Critical KCP pattern violations",
            'action_required': 'FIX_PATTERNS'
        })
    
    # Security issues (CRITICAL)
    security_result = review_data.get('security_review', {})
    critical_security = [issue for issue in security_result.get('issues', []) 
                        if issue.get('severity') == 'CRITICAL']
    if critical_security:
        blocking_issues.append({
            'type': 'CRITICAL_SECURITY_ISSUES',
            'description': f"{len(critical_security)} critical security issues",
            'action_required': 'FIX_SECURITY'
        })
    
    # Make decision
    if blocking_issues:
        decision = {
            'result': 'CHANGES_REQUIRED',
            'blocking_issues': blocking_issues,
            'can_proceed': False,
            'required_actions': [issue['action_required'] for issue in blocking_issues]
        }
    else:
        # Check for non-blocking issues
        warnings = collect_review_warnings(review_data)
        
        if len(warnings) == 0:
            decision_result = 'APPROVED'
        elif len(warnings) <= 3:
            decision_result = 'APPROVED_WITH_WARNINGS'
        else:
            decision_result = 'CHANGES_RECOMMENDED'
        
        decision = {
            'result': decision_result,
            'blocking_issues': [],
            'warnings': warnings,
            'can_proceed': True,
            'recommendations': generate_recommendations(review_data)
        }
    
    return decision
```

## Review Documentation

```yaml
# Code Review Report Template
code_review_report:
  effort_id: "phase1-wave2-effort3-webhooks"
  reviewed_at: "2025-08-23T19:30:00Z"
  reviewer: "code-reviewer-agent"
  
  size_compliance:
    tool_used: "line-counter.sh"
    measured_lines: 687
    limit: 800
    compliant: true
    margin: 113
    
  test_coverage:
    unit_test_coverage: 92.3
    integration_tests: 8
    multi_tenant_scenarios: 5
    performance_tests: 3
    overall_score: 94
    meets_requirements: true
    
  kcp_compliance:
    multi_tenancy_score: 95
    api_export_integration: 90
    workspace_isolation: 92
    syncer_compatibility: 88
    overall_compliance: 91
    
  architecture_review:
    plan_adherence: 89
    design_patterns: 92
    interface_compliance: 94
    component_structure: 87
    
  security_review:
    input_validation: "PASS"
    workspace_isolation: "PASS" 
    rbac_implementation: "PASS"
    secret_management: "PASS"
    critical_issues: 0
    
  performance_review:
    resource_usage: "WITHIN_LIMITS"
    response_times: "ACCEPTABLE"
    scalability: "GOOD"
    
  final_decision:
    result: "APPROVED"
    can_proceed: true
    blocking_issues: []
    warnings: 
      - "Consider adding more error handling tests"
      - "Performance tests could cover more edge cases"
    recommendations:
      - "Add logging for debugging multi-tenant scenarios"
      - "Consider caching for frequently accessed configurations"
```

## State Transitions

From CODE_REVIEW state:
- **APPROVED** → SPAWN_AGENTS (Spawn next effort or integration)
- **CHANGES_REQUIRED** → SPAWN_AGENTS (Spawn SW Engineer for fixes)
- **SIZE_VIOLATION** → SPLIT_PLANNING (Plan effort split)
- **CRITICAL_ISSUES** → ERROR_RECOVERY (Address blocking problems)
