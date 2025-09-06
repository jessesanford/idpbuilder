# Code Reviewer - PERFORM_CODE_REVIEW State Rules

## State Context
You are performing a comprehensive code review of an implementation that has already been measured and confirmed to be within size limits (≤800 lines).

---

### 🚨🚨🚨 RULE R320 - No Stub Implementations (CRITICAL BLOCKER) 🚨🚨🚨
**Source:** rule-library/R320-no-stub-implementations.md
**Criticality:** BLOCKING - Any stub = FAILED REVIEW

**MANDATORY STUB DETECTION PROTOCOL:**
1. Search for ALL "not implemented" patterns
2. Check for TODO in function bodies  
3. Verify each function has ACTUAL logic
4. Any stub found = CRITICAL BLOCKER
5. Stub implementations = IMMEDIATE REJECTION

**Common stub patterns to detect:**
- `return fmt.Errorf("not implemented")`
- `panic("TODO")` or `panic("unimplemented")`
- `raise NotImplementedError`
- Empty function bodies with just return
- `throw new Error("Not implemented")`

**GRADING PENALTIES:**
- **-50%**: Passing ANY stub implementation
- **-30%**: Classifying stub as "minor issue"
- **-40%**: Marking stub code as "properly implemented"

---

### ℹ️ RULE R108.0.0 - CODE_REVIEW Rules
**Source:** rule-library/RULE-REGISTRY.md#R108
**Criticality:** INFO - Best practice

CODE REVIEW PROTOCOL:
1. **CHECK FOR STUBS FIRST (R320)** - Any stub = FAILED REVIEW
2. Validate implementation against plan
3. Verify test coverage requirements
4. Validate KCP/Kubernetes patterns
5. Check multi-tenancy implementation
6. Assess security and performance
7. Provide detailed feedback

---

## Review Focus Areas

### 1. Stub Implementation Detection (R320) - HIGHEST PRIORITY

```python
def detect_stub_implementations(effort_dir):
    """Detect stub implementations per R320 requirements
    
    ANY STUB FOUND = CRITICAL BLOCKER = FAILED REVIEW
    """
    
    stubs_found = []
    
    # Define patterns for different languages
    stub_patterns = {
        'go': [
            r'return.*fmt\.Errorf\("not.*implemented',
            r'return.*errors\.New\("not.*implemented',
            r'return.*errors\.New\("TODO"',
            r'panic\("TODO"\)',
            r'panic\("unimplemented"\)',
            r'panic\("not.*implemented"\)',
        ],
        'python': [
            r'raise NotImplementedError',
            r'return\s+"TODO"',
            r'pass\s+#.*TODO',
        ],
        'javascript': [
            r'throw new Error\("Not implemented',
            r'throw new Error\("TODO',
            r'Promise\.reject\("TODO',
            r'console\.warn\("Not implemented',
        ],
        'typescript': [
            r'throw new Error\("Not implemented',
            r'throw new Error\("TODO',
            r'Promise\.reject\("TODO',
        ]
    }
    
    # Search exhaustively for stubs
    # Check empty function bodies
    # Verify actual implementation logic exists
    
    return {
        'stubs_found': len(stubs_found) > 0,
        'stub_count': len(stubs_found),
        'stub_locations': stubs_found,
        'review_result': 'FAILED' if stubs_found else 'PASSED',
        'critical_blocker': len(stubs_found) > 0
    }
```

### 2. Test Coverage Validation

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
    
    # Check for critical coverage gaps
    unit_coverage = coverage_results['unit_test_coverage'].get('percentage', 0)
    if unit_coverage < 90:
        return {
            'critical_issue': 'INSUFFICIENT_UNIT_COVERAGE',
            'blocking': True
        }
    
    return coverage_results
```

### 3. KCP/Kubernetes Pattern Validation

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

### 4. Security Review

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

### 5. Architecture Compliance

Review implementation against architectural plan:
- Component structure matches design
- Interfaces properly implemented
- Design patterns correctly applied
- Dependencies appropriately managed

## Review Decision Framework

```python
def make_review_decision(review_data):
    """Make final review decision based on all validation results"""
    
    # Critical blocking issues
    blocking_issues = []
    
    # STUB IMPLEMENTATIONS (HIGHEST PRIORITY - R320)
    stub_result = review_data.get('stub_detection', {})
    if stub_result.get('stubs_found', False):
        blocking_issues.append({
            'type': 'STUB_IMPLEMENTATION_DETECTED',
            'severity': 'CRITICAL_BLOCKER',
            'description': f"Found {stub_result.get('stub_count', 0)} stub implementations",
            'action_required': 'COMPLETE_IMPLEMENTATION'
        })
    
    # Test coverage (CRITICAL)
    coverage_result = review_data.get('test_coverage', {})
    if not coverage_result.get('meets_requirements', False):
        blocking_issues.append({
            'type': 'INSUFFICIENT_COVERAGE',
            'description': f"Coverage {coverage_result.get('coverage_score', 0)}% < 90%",
            'action_required': 'IMPROVE_TESTS'
        })
    
    # Make decision
    if blocking_issues:
        return {
            'result': 'FAILED',
            'blocking_issues': blocking_issues,
            'can_proceed': False
        }
    else:
        warnings = collect_review_warnings(review_data)
        
        if len(warnings) == 0:
            return {'result': 'PASSED'}
        elif len(warnings) <= 3:
            return {'result': 'PASSED_WITH_WARNINGS'}
        else:
            return {'result': 'CHANGES_RECOMMENDED'}
```

## Review Documentation

Create CODE-REVIEW-REPORT.md:
```yaml
# Code Review Report
effort_id: "[effort-name]"
reviewed_at: "[timestamp]"
reviewer: "code-reviewer-agent"

## Pre-Review Verification
size_measurement:
  completed: true
  compliant: true
  lines: [number]
  
## Review Results

### 1. Stub Detection (R320)
stubs_found: [true/false]
stub_count: [number]
stub_locations: []
result: [PASSED/FAILED]

### 2. Test Coverage
unit_test_coverage: [percentage]
integration_tests: [count]
multi_tenant_scenarios: [count]
meets_requirements: [true/false]

### 3. KCP Compliance
multi_tenancy_score: [percentage]
api_export_integration: [percentage]
workspace_isolation: [percentage]
overall_compliance: [percentage]

### 4. Security Review
input_validation: [PASS/FAIL]
workspace_isolation: [PASS/FAIL]
secret_management: [PASS/FAIL]
critical_issues: [count]

### 5. Architecture Review
plan_adherence: [percentage]
design_patterns: [percentage]
interface_compliance: [percentage]

## Final Decision
REVIEW_STATUS: [PASSED/FAILED/PASSED_WITH_WARNINGS]
blocking_issues: []
warnings: []
recommendations: []

## Required Actions
[List any required fixes if FAILED]
```

## State Transitions

From PERFORM_CODE_REVIEW state:
- **REVIEW_PASSED** → COMPLETED (Implementation approved)
- **REVIEW_FAILED** → CREATE_FIX_PLAN (Issues need fixing)
- **CRITICAL_STUBS_FOUND** → CREATE_FIX_PLAN (R320 violation)

## Success Criteria
- ✅ Thoroughly checked for stub implementations
- ✅ Validated test coverage ≥90%
- ✅ Verified KCP patterns compliance
- ✅ Completed security assessment
- ✅ Created comprehensive review report
- ✅ Made clear pass/fail decision

## Failure Triggers
- ❌ Missing stub detection = R320 VIOLATION (-50%)
- ❌ Passing stub code = R320 VIOLATION (-50%)
- ❌ Incomplete review = -30% penalty
- ❌ Missing review report = -40% penalty