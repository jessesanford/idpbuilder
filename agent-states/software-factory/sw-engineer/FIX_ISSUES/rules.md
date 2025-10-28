# SW Engineer - FIX_ISSUES State Rules

## State Context
You are focused on resolving specific issues, bugs, or optimization requirements that are blocking progress.

## 🔴🔴🔴 SUPREME LAW R355: PRODUCTION READY FIXES ONLY 🔴🔴🔴

### FIXES MUST NOT INTRODUCE:
- ❌ **Hardcoded Credentials** - No quick-fix passwords
- ❌ **Stub Implementations** - No "fix later" placeholders
- ❌ **Mock/Fake Objects** - Real fixes only
- ❌ **Static Values** - Keep everything configurable
- ❌ **TODO/FIXME Comments** - Complete the fix fully

### VERIFY FIXES ARE PRODUCTION READY:
```bash
echo "🔴 R355: VERIFYING FIX IS PRODUCTION READY"
cd $EFFORT_DIR
# Check fix doesn't introduce violations
VIOLATIONS=0
grep -r "password.*=.*['\"]" --exclude-dir=test --include="*.go" --include="*.py" --include="*.js" && VIOLATIONS=1
grep -r "stub\|mock\|fake" --exclude-dir=test --include="*.go" --include="*.py" --include="*.js" && VIOLATIONS=1
grep -r "TODO\|FIXME\|TEMPORARY" --exclude-dir=test --include="*.go" --include="*.py" --include="*.js" && VIOLATIONS=1
grep -r "not.*implemented" --exclude-dir=test --include="*.go" --include="*.py" --include="*.js" && VIOLATIONS=1

if [ $VIOLATIONS -eq 1 ]; then
    echo "🚨 R355 VIOLATION: Fix contains non-production code!"
    echo "Even fixes must be production ready!"
    exit 355
fi
echo "✅ R355: Fix is production ready"
```

## 🔴🔴🔴 CRITICAL: R340 Plan Location Tracking 🔴🔴🔴

**RULE R340: Planning File Metadata Tracking (BLOCKING)**
- Review reports are tracked in orchestrator-state-v3.json
- Read review report location from .effort_repo_files.review_reports
- NEVER search for reports using `find` or `ls` commands
- The orchestrator tracks ALL planning files in the state file
- Violation of R340 = Integration delays and failures

### R340: Reading Review Reports from State

```bash
# CORRECT: Read review report location from state (R340 compliant)
EFFORT_NAME=$(basename $(pwd))
ORCHESTRATOR_STATE_PATH="$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json"

REVIEW_REPORT=$(jq -r ".effort_repo_files.review_reports.\"${EFFORT_NAME}\".file_path // null" \
  "$ORCHESTRATOR_STATE_PATH")

if [ -z "$REVIEW_REPORT" ] || [ "$REVIEW_REPORT" == "null" ]; then
  echo "⚠️  No review report tracked for effort: $EFFORT_NAME"
  echo "Orchestrator should have recorded review report location per R340"
else
  echo "📋 R340: Using tracked review report: $REVIEW_REPORT"

  # Verify file exists
  if [ ! -f "$REVIEW_REPORT" ]; then
    echo "❌ Tracked review report missing: $REVIEW_REPORT"
    exit 1
  fi

  # Read review feedback
  echo "Reading review feedback..."
  cat "$REVIEW_REPORT"
fi

# ❌ WRONG: Searching for review reports (R340 VIOLATION)
# NEVER DO THIS:
# LATEST_REPORT=$(ls -t .software-factory/phase*/wave*/*/CODE-REVIEW-REPORT--*.md | head -n1)
```

## 🔴🔴🔴 CRITICAL: FIX APPLICATION LOCATION 🔴🔴🔴

---
### 🔴🔴🔴 RULE R300 - Comprehensive Fix Management Protocol (SUPREME LAW)
**Source:** rule-library/R300-comprehensive-fix-management-protocol.md
**Criticality:** SUPREME LAW - Violation = -100% AUTOMATIC FAILURE

ALL FIXES MUST BE APPLIED TO EFFORT BRANCHES, NEVER TO INTEGRATE_WAVE_EFFORTS BRANCHES!

MANDATORY VERIFICATION:
```bash
# FIRST THING IN FIX_ISSUES STATE:
cd /efforts/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}
CURRENT_BRANCH=$(git branch --show-current)
if [[ ! "$CURRENT_BRANCH" =~ ^effort- ]]; then
    echo "🔴 CRITICAL ERROR: Not on effort branch! Violates R300!"
    exit 1
fi
echo "✅ Confirmed on effort branch: $CURRENT_BRANCH"
```

KEY POINTS:
- Fixes in integration branches are LOST when branches recreated
- This causes infinite loops of fixing the same issues
- ALWAYS work in effort directory and effort branch
- ALWAYS push to effort branch remote
---

---
### 🚨 RULE R339 - Fix Grace Period Protocol
**Source:** rule-library/R339-fix-grace-period-protocol.md
**Criticality:** CRITICAL - Prevents unnecessary splits during fixes

FIX GRACE PERIOD FOR RE-INTEGRATE_WAVE_EFFORTS FIXES:
- Original effort + fix < 900 lines = NO SPLIT REQUIRED (100 line grace)
- Original effort + fix >= 900 lines = SPLIT REQUIRED

WHEN APPLYING FIXES:
```bash
# Check if fix will trigger split with grace period
check_fix_size_with_grace() {
    # Get current effort size
    CURRENT_SIZE=$($CLAUDE_PROJECT_DIR/tools/line-counter.sh -c $(git branch --show-current))
    echo "📊 Current effort size: ${CURRENT_SIZE} lines"
    
    # Estimate fix size (be conservative)
    echo "📝 Estimating fix impact..."
    # Your fix estimate here
    ESTIMATED_FIX_SIZE=70  # Example
    
    TOTAL_WITH_FIX=$((CURRENT_SIZE + ESTIMATED_FIX_SIZE))
    echo "📊 Projected total: ${TOTAL_WITH_FIX} lines"
    
    if [ "$TOTAL_WITH_FIX" -lt 900 ]; then
        echo "✅ Within grace period (${TOTAL_WITH_FIX} < 900)"
        echo "✅ No split required - proceed with fix"
    else
        echo "⚠️ Would exceed grace period (${TOTAL_WITH_FIX} >= 900)"
        echo "🔀 Split will be required after fix"
        echo "💡 Consider minimizing fix to stay under 900"
    fi
}

# Run this check BEFORE starting fix implementation
check_fix_size_with_grace
```

KEY BENEFITS:
- Prevents cascade disruption for minor overages
- Allows fixes up to 900 lines total
- Grace only applies ONCE per effort
- Does NOT apply to feature additions
---

---
### ℹ️ RULE R109.0.0 - FIX_ISSUES Rules
**Source:** rule-library/RULE-REGISTRY.md#R109
**Criticality:** INFO - Best practice

ISSUE RESOLUTION PROTOCOL:
1. Clearly identify and document the issue
2. Analyze root cause before implementing fix
3. Design minimal, focused solution
4. Implement fix with proper testing
5. Verify issue resolution and measure impact
6. Document fix and lessons learned
---

## Issue Classification and Prioritization

### Reading Issues from Bug Tracking (SF 3.0)

In SF 3.0, issues may be tracked centrally in bug-tracking.json:

```bash
# Read issues from bug tracking system
ORCHESTRATOR_STATE_PATH="$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json"
BUG_TRACKING_PATH="$CLAUDE_PROJECT_DIR/bug-tracking.json"

if [ -f "$BUG_TRACKING_PATH" ]; then
    # Get bugs assigned to current effort
    EFFORT_NAME=$(basename $(pwd))
    EFFORT_BUGS=$(jq -r ".bugs[] | select(.effort_name == \"$EFFORT_NAME\" and .status != \"RESOLVED\")" \
      "$BUG_TRACKING_PATH")

    if [ -n "$EFFORT_BUGS" ]; then
        echo "📋 Found tracked bugs for this effort in bug-tracking.json"
        echo "$EFFORT_BUGS" | jq -r '.id + ": " + .title'
    fi
fi
```

---
### 🚨🚨🚨 RULE R019.0.0 - Error Recovery
**Source:** rule-library/RULE-REGISTRY.md#R019
**Criticality:** BLOCKING - Failure = Immediate Stop (exit 1)

ISSUE PRIORITY CLASSIFICATION:

CRITICAL (Fix immediately):
- Size limit violations (>800 lines)
- Build failures preventing compilation
- Test failures blocking code review
- Security vulnerabilities

HIGH (Fix within session):
- Performance regressions
- API compatibility issues
- Missing error handling
- Code quality violations

MEDIUM (Address before completion):
- Code optimization opportunities
- Maintainability improvements
- Documentation gaps
- Test coverage improvements

LOW (Nice to have):
- Code style improvements
- Minor refactoring opportunities
- Performance optimizations
---

```python
def classify_and_prioritize_issues(issue_list):
    """Classify and prioritize issues for systematic resolution"""
    
    classified_issues = {
        'critical': [],
        'high': [],
        'medium': [],
        'low': []
    }
    
    for issue in issue_list:
        priority = determine_issue_priority(issue)
        classified_issues[priority].append({
            'issue': issue,
            'priority': priority,
            'estimated_effort': estimate_fix_effort(issue),
            'impact_analysis': analyze_fix_impact(issue),
            'dependencies': identify_fix_dependencies(issue)
        })
    
    # Sort within each priority level by effort and impact
    for priority_level in classified_issues:
        classified_issues[priority_level].sort(
            key=lambda x: (x['estimated_effort'], -x['impact_analysis']['positive_impact_score'])
        )
    
    return classified_issues

def determine_issue_priority(issue):
    """Determine issue priority based on type and impact"""
    
    issue_type = issue.get('type', '').lower()
    impact = issue.get('impact', 'medium').lower()
    blocking = issue.get('blocks_progress', False)
    
    # Critical issues
    if any(keyword in issue_type for keyword in ['size_violation', 'build_failure', 'security']):
        return 'critical'
    
    if blocking and impact == 'high':
        return 'critical'
    
    # High priority issues
    if any(keyword in issue_type for keyword in ['test_failure', 'performance', 'api_break']):
        return 'high'
    
    if blocking and impact == 'medium':
        return 'high'
    
    # Medium priority issues  
    if any(keyword in issue_type for keyword in ['optimization', 'maintainability', 'coverage']):
        return 'medium'
    
    if impact == 'medium' and not blocking:
        return 'medium'
    
    # Low priority (everything else)
    return 'low'

def estimate_fix_effort(issue):
    """Estimate effort required to fix the issue (in hours)"""
    
    issue_type = issue.get('type', '').lower()
    complexity = issue.get('complexity', 'medium').lower()
    scope = issue.get('scope', 'single_file').lower()
    
    base_effort = {
        'syntax_error': 0.25,
        'test_failure': 0.5,
        'build_failure': 1.0,
        'performance': 1.5,
        'optimization': 2.0,
        'refactoring': 3.0,
        'architecture': 4.0
    }.get(issue_type, 1.0)
    
    complexity_multiplier = {
        'simple': 0.5,
        'medium': 1.0,
        'complex': 2.0,
        'very_complex': 3.0
    }.get(complexity, 1.0)
    
    scope_multiplier = {
        'single_line': 0.5,
        'single_function': 1.0,
        'single_file': 1.5,
        'multiple_files': 2.0,
        'package_wide': 3.0
    }.get(scope, 1.5)
    
    return base_effort * complexity_multiplier * scope_multiplier
```

## Size-Related Issue Resolution

---
### ℹ️ RULE R007.0.0 - Size Limit Enforcement
**Source:** rule-library/RULE-REGISTRY.md#R007
**Criticality:** INFO - Best practice

SIZE OPTIMIZATION STRATEGIES:

CODE REDUCTION TECHNIQUES:
1. Extract helper functions to reduce duplication
2. Use table-driven patterns for repetitive logic
3. Optimize verbose error handling
4. Consolidate similar functions
5. Remove unnecessary comments and whitespace

REFACTORING APPROACHES:
- Split large functions into smaller, focused ones
- Create utility packages for common operations
- Use interfaces to reduce concrete implementations
- Leverage Go standard library instead of custom code

MEASUREMENT VERIFICATION:
- Use line-counter.sh per R304 to verify reductions
- Measure after each optimization
- Ensure functionality is preserved
- Maintain or improve test coverage
---

```go
// Example of code size optimization techniques

// BEFORE: Verbose, repetitive code (45 lines)
func (r *ResourceController) updateResourceStatus(ctx context.Context, resource *myapi.Resource, phase myapi.ResourcePhase, message string) error {
    if resource.Status.Phase != phase {
        resource.Status.Phase = phase
        resource.Status.LastTransitionTime = metav1.Now()
    }
    
    condition := metav1.Condition{
        Type:               myapi.ConditionTypeReady,
        Status:             metav1.ConditionTrue,
        LastTransitionTime: metav1.Now(),
        Reason:             "StatusUpdated", 
        Message:            message,
    }
    
    if phase != myapi.ResourcePhaseReady {
        condition.Status = metav1.ConditionFalse
        condition.Reason = "NotReady"
    }
    
    meta.SetStatusCondition(&resource.Status.Conditions, condition)
    
    if err := r.Status().Update(ctx, resource); err != nil {
        r.Log.Error(err, "Failed to update resource status", "resource", resource.Name, "phase", phase)
        return fmt.Errorf("failed to update resource status: %w", err)
    }
    
    r.Log.Info("Resource status updated", "resource", resource.Name, "phase", phase, "message", message)
    return nil
}

func (r *ResourceController) updateResourceStatusError(ctx context.Context, resource *myapi.Resource, err error) error {
    resource.Status.Phase = myapi.ResourcePhaseError
    resource.Status.LastTransitionTime = metav1.Now()
    
    condition := metav1.Condition{
        Type:               myapi.ConditionTypeReady,
        Status:             metav1.ConditionFalse,
        LastTransitionTime: metav1.Now(),
        Reason:             "ReconcileError",
        Message:            err.Error(),
    }
    
    meta.SetStatusCondition(&resource.Status.Conditions, condition)
    
    if updateErr := r.Status().Update(ctx, resource); updateErr != nil {
        r.Log.Error(updateErr, "Failed to update resource error status", "resource", resource.Name)
        return fmt.Errorf("failed to update error status: %w", updateErr)
    }
    
    return nil
}

// AFTER: Optimized, DRY code (25 lines) - 44% reduction
func (r *ResourceController) updateResourceStatus(ctx context.Context, resource *myapi.Resource, phase myapi.ResourcePhase, message string) error {
    return r.setResourceStatus(ctx, resource, phase, myapi.ConditionTypeReady, 
        r.getConditionStatus(phase), r.getConditionReason(phase), message)
}

func (r *ResourceController) updateResourceStatusError(ctx context.Context, resource *myapi.Resource, err error) error {
    return r.setResourceStatus(ctx, resource, myapi.ResourcePhaseError, myapi.ConditionTypeReady,
        metav1.ConditionFalse, "ReconcileError", err.Error())
}

func (r *ResourceController) setResourceStatus(ctx context.Context, resource *myapi.Resource, phase myapi.ResourcePhase, 
    condType string, condStatus metav1.ConditionStatus, reason, message string) error {
    
    if resource.Status.Phase != phase {
        resource.Status.Phase = phase
        resource.Status.LastTransitionTime = metav1.Now()
    }
    
    meta.SetStatusCondition(&resource.Status.Conditions, metav1.Condition{
        Type: condType, Status: condStatus, LastTransitionTime: metav1.Now(),
        Reason: reason, Message: message,
    })
    
    if err := r.Status().Update(ctx, resource); err != nil {
        r.Log.Error(err, "Failed to update status", "resource", resource.Name, "phase", phase)
        return fmt.Errorf("failed to update resource status: %w", err)
    }
    
    r.Log.Info("Resource status updated", "resource", resource.Name, "phase", phase)
    return nil
}

func (r *ResourceController) getConditionStatus(phase myapi.ResourcePhase) metav1.ConditionStatus {
    if phase == myapi.ResourcePhaseReady { return metav1.ConditionTrue }
    return metav1.ConditionFalse
}

func (r *ResourceController) getConditionReason(phase myapi.ResourcePhase) string {
    if phase == myapi.ResourcePhaseReady { return "StatusUpdated" }
    return "NotReady"
}
```

## Test-Related Issue Resolution

```python
def resolve_test_related_issues(test_issues):
    """Resolve test-related issues systematically"""
    
    print("🧪 RESOLVING TEST-RELATED ISSUES")
    
    resolution_plan = {
        'failing_tests': [],
        'coverage_issues': [],
        'performance_issues': [],
        'maintenance_issues': []
    }
    
    # Categorize test issues
    for issue in test_issues:
        category = categorize_test_issue(issue)
        resolution_plan[category].append(issue)
    
    # Resolve in priority order
    results = {}
    
    # 1. Fix failing tests first (blocks everything)
    if resolution_plan['failing_tests']:
        results['failing_tests'] = resolve_failing_tests(resolution_plan['failing_tests'])
    
    # 2. Address coverage issues  
    if resolution_plan['coverage_issues']:
        results['coverage_issues'] = resolve_coverage_issues(resolution_plan['coverage_issues'])
    
    # 3. Fix performance issues
    if resolution_plan['performance_issues']:
        results['performance_issues'] = resolve_test_performance_issues(resolution_plan['performance_issues'])
    
    # 4. Maintenance improvements
    if resolution_plan['maintenance_issues']:
        results['maintenance_issues'] = resolve_test_maintenance_issues(resolution_plan['maintenance_issues'])
    
    return results

def resolve_failing_tests(failing_tests):
    """Resolve failing test issues"""
    
    resolution_results = []
    
    for test_issue in failing_tests:
        test_name = test_issue.get('test_name', '')
        failure_reason = test_issue.get('failure_reason', '')
        
        print(f"🔧 Fixing failing test: {test_name}")
        print(f"   Failure reason: {failure_reason}")
        
        # Analyze failure type and apply appropriate fix strategy
        if 'timeout' in failure_reason.lower():
            fix_result = fix_timeout_test(test_issue)
        elif 'assertion' in failure_reason.lower() or 'expected' in failure_reason.lower():
            fix_result = fix_assertion_test(test_issue)
        elif 'panic' in failure_reason.lower() or 'nil pointer' in failure_reason.lower():
            fix_result = fix_panic_test(test_issue)
        elif 'race' in failure_reason.lower():
            fix_result = fix_race_condition_test(test_issue)
        else:
            fix_result = fix_generic_test_failure(test_issue)
        
        resolution_results.append({
            'test_name': test_name,
            'original_issue': test_issue,
            'fix_applied': fix_result,
            'resolution_status': fix_result.get('success', False)
        })
        
        # Verify fix by running the specific test
        if fix_result.get('success', False):
            verification = verify_test_fix(test_name)
            resolution_results[-1]['verification'] = verification
    
    return resolution_results

def fix_timeout_test(test_issue):
    """Fix test timeout issues"""
    
    test_file = test_issue.get('test_file', '')
    test_name = test_issue.get('test_name', '')
    
    # Common timeout fixes
    fixes_applied = []
    
    # 1. Increase context timeout
    if 'context' in test_issue.get('failure_output', '').lower():
        fixes_applied.append('increased_context_timeout')
        # Implementation: modify context timeout in test
    
    # 2. Add proper wait conditions
    if 'eventually' not in test_issue.get('test_code', '').lower():
        fixes_applied.append('added_eventually_conditions')
        # Implementation: add gomega Eventually conditions
    
    # 3. Mock slow operations
    if 'http' in test_issue.get('failure_output', '').lower():
        fixes_applied.append('mocked_slow_operations')
        # Implementation: mock HTTP calls
    
    return {
        'success': len(fixes_applied) > 0,
        'fixes_applied': fixes_applied,
        'estimated_improvement': 'Test should now complete within timeout'
    }

def fix_assertion_test(test_issue):
    """Fix assertion-related test failures"""
    
    failure_output = test_issue.get('failure_output', '')
    
    fixes_applied = []
    
    # Parse expected vs actual values
    expected_actual = parse_assertion_failure(failure_output)
    
    if expected_actual:
        # Determine if it's a logic issue or assertion issue
        if looks_like_assertion_issue(expected_actual):
            fixes_applied.append('corrected_assertion_values')
            # Implementation: update expected values in test
        else:
            fixes_applied.append('identified_logic_bug')
            # Implementation: need to fix actual implementation
    
    # Check for timing-related assertion issues
    if 'eventually' not in failure_output.lower() and is_timing_sensitive(test_issue):
        fixes_applied.append('added_eventually_assertions')
        # Implementation: wrap assertions in Eventually
    
    return {
        'success': len(fixes_applied) > 0,
        'fixes_applied': fixes_applied,
        'requires_implementation_fix': 'identified_logic_bug' in fixes_applied
    }
```

## Performance Issue Resolution

```go
// Example of performance issue resolution

// BEFORE: Inefficient reconciliation logic
func (r *ResourceController) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    // Inefficient: Multiple separate API calls
    var resource myapi.Resource
    if err := r.Get(ctx, req.NamespacedName, &resource); err != nil {
        return ctrl.Result{}, err
    }
    
    // Inefficient: Checking every possible child resource individually
    pods := &corev1.PodList{}
    if err := r.List(ctx, pods, client.InNamespace(resource.Namespace)); err != nil {
        return ctrl.Result{}, err
    }
    
    services := &corev1.ServiceList{}  
    if err := r.List(ctx, services, client.InNamespace(resource.Namespace)); err != nil {
        return ctrl.Result{}, err
    }
    
    deployments := &appsv1.DeploymentList{}
    if err := r.List(ctx, deployments, client.InNamespace(resource.Namespace)); err != nil {
        return ctrl.Result{}, err
    }
    
    // Process each resource type separately...
    for _, pod := range pods.Items {
        // Processing logic...
    }
    
    return ctrl.Result{}, nil
}

// AFTER: Optimized reconciliation logic  
func (r *ResourceController) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    var resource myapi.Resource
    if err := r.Get(ctx, req.NamespacedName, &resource); err != nil {
        if apierrors.IsNotFound(err) {
            return ctrl.Result{}, nil // Resource deleted
        }
        return ctrl.Result{}, err
    }
    
    // Optimized: Single call to get all owned resources using label selector
    ownedResources, err := r.getOwnedResources(ctx, &resource)
    if err != nil {
        return ctrl.Result{}, err
    }
    
    // Optimized: Process all resources in a single pass
    reconResult, err := r.reconcileOwnedResources(ctx, &resource, ownedResources)
    if err != nil {
        return ctrl.Result{}, err
    }
    
    return reconResult, nil
}

func (r *ResourceController) getOwnedResources(ctx context.Context, resource *myapi.Resource) (*OwnedResources, error) {
    // Optimized: Use label selector to get all owned resources at once
    selector := labels.SelectorFromSet(map[string]string{
        "app.kubernetes.io/managed-by": "resource-controller",
        "app.kubernetes.io/name":       resource.Name,
    })
    
    var allObjects unstructured.UnstructuredList
    allObjects.SetGroupVersionKind(schema.GroupVersionKind{
        Group:   "",
        Version: "v1", 
        Kind:    "List",
    })
    
    err := r.List(ctx, &allObjects, 
        client.InNamespace(resource.Namespace),
        client.MatchingLabelsSelector{Selector: selector})
    
    if err != nil {
        return nil, err
    }
    
    // Optimized: Categorize resources efficiently
    return r.categorizeOwnedResources(allObjects.Items), nil
}
```

## Issue Resolution Documentation

---
### ℹ️ RULE R018.0.0 - Progress Reporting
**Source:** rule-library/RULE-REGISTRY.md#R018
**Criticality:** INFO - Best practice

ISSUE RESOLUTION DOCUMENTATION:

REQUIRED DOCUMENTATION:
1. Issue description and root cause analysis
2. Solution approach and alternatives considered
3. Implementation details and code changes
4. Testing performed to verify the fix
5. Impact assessment and side effects
6. Lessons learned and prevention strategies

DOCUMENTATION LOCATION:
- Update work-log.md with fix details
- Create FIXES.md if multiple issues resolved
- Update implementation plan with changes
- Document in commit messages
---

```markdown
# Issue Resolution Documentation Template

## Issue: [Brief Description]
**Date**: 2025-08-23 17:45:00
**Priority**: CRITICAL/HIGH/MEDIUM/LOW
**Type**: SIZE_VIOLATION/TEST_FAILURE/PERFORMANCE/etc.
**Status**: RESOLVED/IN_PROGRESS/BLOCKED

### Problem Description
Clear description of the issue:
- What was observed/failing
- When it was detected
- Impact on development progress
- Blocking factors

### Root Cause Analysis
Deep dive into the underlying cause:
- Why the issue occurred
- Contributing factors
- How it was diagnosed
- Tools used for analysis

### Solution Approach
Explanation of the chosen solution:
- Alternatives considered
- Why this approach was selected
- Trade-offs and compromises made
- Expected impact and benefits

### Implementation Details
Specific changes made:
- Files modified
- Functions changed
- New code added
- Code removed/refactored
- Configuration changes

### Testing and Verification
How the fix was validated:
- Tests run to verify fix
- Regression testing performed
- Performance impact measured
- Edge cases validated

### Results and Impact
Outcome of the fix:
- Issue resolution confirmed
- Performance improvements
- Side effects observed
- Remaining work if any

### Lessons Learned
Key takeaways:
- What could have prevented this issue
- Process improvements identified
- Knowledge gained
- Prevention strategies for future

### Metrics
Quantitative measurements:
- Lines changed: +25, -18 
- Size impact: Reduced by 45 lines
- Performance improvement: 15% faster
- Test coverage: Maintained at 87%
```

## Size Optimization Systematic Approach

```python
def systematic_size_optimization(effort_path):
    """Perform systematic size optimization to meet limits"""
    
    print("📉 SYSTEMATIC SIZE OPTIMIZATION")
    
    # 1. Measure current size
    current_size = measure_effort_size(effort_path)
    print(f"Current size: {current_size} lines")
    
    if current_size <= 800:
        print("✅ Size already compliant")
        return {'success': True, 'reduction_needed': 0}
    
    reduction_needed = current_size - 800
    print(f"Reduction needed: {reduction_needed} lines")
    
    # 2. Analyze optimization opportunities
    optimization_opportunities = analyze_optimization_opportunities(effort_path)
    
    # 3. Apply optimizations in order of impact/effort ratio
    optimization_plan = create_optimization_plan(optimization_opportunities, reduction_needed)
    
    results = {
        'initial_size': current_size,
        'reduction_target': reduction_needed,
        'optimizations_applied': [],
        'final_size': current_size,
        'success': False
    }
    
    for optimization in optimization_plan:
        print(f"🔧 Applying: {optimization['description']}")
        
        opt_result = apply_optimization(effort_path, optimization)
        
        if opt_result['success']:
            new_size = measure_effort_size(effort_path)
            actual_reduction = results['final_size'] - new_size
            
            results['optimizations_applied'].append({
                'optimization': optimization,
                'result': opt_result,
                'size_before': results['final_size'],
                'size_after': new_size,
                'actual_reduction': actual_reduction
            })
            
            results['final_size'] = new_size
            
            print(f"   Reduced by {actual_reduction} lines (now {new_size} lines)")
            
            # Check if we've achieved compliance
            if new_size <= 800:
                results['success'] = True
                print("✅ Size optimization successful!")
                break
        else:
            print(f"   ❌ Optimization failed: {opt_result.get('error', 'Unknown error')}")
    
    # 4. Verify functionality is preserved
    if results['success']:
        print("🧪 Verifying functionality after optimization...")
        verification = verify_functionality_preserved(effort_path)
        results['functionality_verification'] = verification
        
        if not verification['all_tests_passing']:
            print("⚠️ Tests failing after optimization - may need to revert some changes")
            results['success'] = False
    
    return results

def analyze_optimization_opportunities(effort_path):
    """Analyze code for optimization opportunities"""
    
    opportunities = []
    
    # Find Go files to analyze
    go_files = find_go_files(effort_path)
    
    for go_file in go_files:
        file_opportunities = analyze_file_for_optimizations(go_file)
        opportunities.extend(file_opportunities)
    
    # Sort by impact/effort ratio
    opportunities.sort(key=lambda x: x['impact_effort_ratio'], reverse=True)
    
    return opportunities

def analyze_file_for_optimizations(file_path):
    """Analyze individual Go file for optimization opportunities"""
    
    with open(file_path, 'r') as f:
        content = f.read()
    
    opportunities = []
    
    # 1. Find duplicate code blocks
    duplicates = find_duplicate_code_blocks(content)
    for duplicate in duplicates:
        opportunities.append({
            'type': 'EXTRACT_DUPLICATE_CODE',
            'file': file_path,
            'description': f'Extract {duplicate["lines"]} duplicate lines into helper function',
            'estimated_reduction': duplicate['lines'] * (duplicate['occurrences'] - 1),
            'effort_hours': 0.5,
            'impact_effort_ratio': duplicate['lines'] * duplicate['occurrences'] / 0.5,
            'details': duplicate
        })
    
    # 2. Find long functions that can be split
    long_functions = find_long_functions(content)
    for func in long_functions:
        if func['lines'] > 50:
            opportunities.append({
                'type': 'SPLIT_LONG_FUNCTION',
                'file': file_path,
                'description': f'Split {func["name"]} function ({func["lines"]} lines)',
                'estimated_reduction': func['lines'] * 0.1,  # Modest reduction through better organization
                'effort_hours': 1.0,
                'impact_effort_ratio': func['lines'] * 0.1 / 1.0,
                'details': func
            })
    
    # 3. Find verbose error handling
    verbose_errors = find_verbose_error_handling(content)
    for error_pattern in verbose_errors:
        opportunities.append({
            'type': 'OPTIMIZE_ERROR_HANDLING',
            'file': file_path,
            'description': f'Optimize verbose error handling pattern',
            'estimated_reduction': error_pattern['lines_saved'],
            'effort_hours': 0.3,
            'impact_effort_ratio': error_pattern['lines_saved'] / 0.3,
            'details': error_pattern
        })
    
    # 4. Find opportunities to use standard library
    stdlib_opportunities = find_stdlib_replacement_opportunities(content)
    for stdlib_opp in stdlib_opportunities:
        opportunities.append({
            'type': 'USE_STDLIB',
            'file': file_path,
            'description': f'Replace custom code with standard library',
            'estimated_reduction': stdlib_opp['lines_saved'],
            'effort_hours': 0.25,
            'impact_effort_ratio': stdlib_opp['lines_saved'] / 0.25,
            'details': stdlib_opp
        })
    
    return opportunities
```

## State Transitions

From FIX_ISSUES state:
- **ISSUES_RESOLVED** → IMPLEMENTATION (Continue development)
- **SIZE_COMPLIANT** → IMPLEMENTATION (Size issues fixed)
- **SIZE_STILL_VIOLATION** → SPLIT_WORK (Optimization insufficient)
- **TESTS_FAILING** → TEST_WRITING (Need more comprehensive tests)
- **COMPLEX_ISSUES** → Continue FIX_ISSUES (More time needed)
- **BLOCKING_DEPENDENCIES** → MONITOR (Wait for external resolution)

## Issue Resolution Validation

```python
def validate_issue_resolution(issue_list, resolution_results):
    """Validate that issues have been properly resolved"""
    
    validation_results = {
        'issues_resolved': 0,
        'issues_remaining': 0,
        'new_issues_introduced': 0,
        'resolution_effectiveness': 0,
        'overall_success': False
    }
    
    for issue in issue_list:
        issue_id = issue.get('id', '')
        resolution = next((r for r in resolution_results if r.get('issue_id') == issue_id), None)
        
        if resolution and resolution.get('resolved', False):
            validation_results['issues_resolved'] += 1
            
            # Verify resolution didn't introduce new issues
            verification = verify_resolution_side_effects(issue, resolution)
            if verification['new_issues']:
                validation_results['new_issues_introduced'] += len(verification['new_issues'])
        else:
            validation_results['issues_remaining'] += 1
    
    total_issues = len(issue_list)
    if total_issues > 0:
        validation_results['resolution_effectiveness'] = (
            validation_results['issues_resolved'] / total_issues * 100
        )
    
    validation_results['overall_success'] = (
        validation_results['resolution_effectiveness'] >= 80 and
        validation_results['new_issues_introduced'] <= 1
    )
    
    return validation_results
```

## Fix Plan Archival Requirements (R294)

---
### 🚨🚨🚨 RULE R294 - Fix Plan Archival Protocol
**Source:** rule-library/R294-fix-plan-archival-protocol.md
**Criticality:** BLOCKING - Must archive completed fix plans

ARCHIVAL REQUIREMENTS:

When completing fixes from any fix plan, you MUST archive the plan:

1. **CODE-REVIEW-REPORT files (R340 Compliant)**:
   ```bash
   # R340: Review reports are tracked in orchestrator-state-v3.json
   # Read location from state instead of searching
   EFFORT_NAME=$(basename $(pwd))
   ORCHESTRATOR_STATE_PATH="$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json"

   REVIEW_REPORT=$(jq -r ".effort_repo_files.review_reports.\"${EFFORT_NAME}\".file_path // null" \
     "$ORCHESTRATOR_STATE_PATH")

   if [ -n "$REVIEW_REPORT" ] && [ "$REVIEW_REPORT" != "null" ] && [ -f "$REVIEW_REPORT" ]; then
       # Archive the review report
       ARCHIVED_NAME="${REVIEW_REPORT%.md}-COMPLETED-$(date +%Y%m%d-%H%M%S).md"
       mv "$REVIEW_REPORT" "$ARCHIVED_NAME"
       echo "✅ R340: Archived completed review report: $ARCHIVED_NAME"

       # Update state to mark as archived (orchestrator should handle this)
       echo "📋 R340: Orchestrator should update state to mark report as archived"
   fi

   # ❌ WRONG - R340 VIOLATION:
   # LATEST_REPORT=$(ls -t .software-factory/phase*/wave*/*/CODE-REVIEW-REPORT--*.md | head -n1)
   ```

2. **SPLIT-PLAN files (R340 Compliant)**:
   ```bash
   # R340: Split plans are tracked in orchestrator-state-v3.json
   # We don't search for them - orchestrator knows their location
   # If we need to archive a split plan, we should:
   # 1. Read its location from state file
   # 2. Archive it
   # 3. Report to orchestrator to update state
   
   echo "📋 Note: Split plan locations are tracked per R340"
   echo "Orchestrator must update state when plans are archived"
   ```

3. **INTEGRATE_WAVE_EFFORTS-REPORT.md**:
   ```bash
   # When integration fixes are complete
   mv INTEGRATE_WAVE_EFFORTS-REPORT.md INTEGRATE_WAVE_EFFORTS-REPORT-COMPLETED-$(date +%Y%m%d-%H%M%S).md
   echo "✅ Archived completed integration report"
   ```

ARCHIVAL VERIFICATION:
- Ensure no active fix plans remain after archival
- Check that archived files follow naming convention
- Commit the archival as part of fix completion

VIOLATIONS:
- ❌ Leaving fix plans unarchived after completion
- ❌ Wrong naming convention for archives
- ❌ Multiple active fix plans causing confusion
---

## 🔴🔴🔴 MANDATORY: Create FIX-COMPLETE Marker 🔴🔴🔴

### CRITICAL REQUIREMENT - FIXES ARE NOT COMPLETE WITHOUT THIS

**When all fixes are complete, you MUST create a completion marker:**

```bash
# MANDATORY when fixes are complete:
echo "🔴 Creating MANDATORY FIX-COMPLETE marker..."
touch FIX-COMPLETE.marker
cat > FIX-COMPLETE.marker << EOF
Completed at: $(date '+%Y-%m-%d %H:%M:%S %Z')
Effort: $(basename $(pwd))
Branch: $(git branch --show-current)
Fixes applied: [List number of fixes]
Tests passing: [Yes/No with count]
Build status: [Success/Failure]
Final commit: $(git log --oneline -1)
Status: FIXES COMPLETE
EOF

# MUST add, commit and push
git add FIX-COMPLETE.marker
git commit -m "marker: fixes complete - MANDATORY for orchestrator monitoring"
git push

echo "✅ FIX-COMPLETE.marker created and pushed"
echo "📋 Orchestrator can now detect fix completion"
```

**THIS IS NOT OPTIONAL:**
- ❌ WITHOUT marker = Orchestrator cannot detect fix completion
- ❌ WITHOUT marker = Fixes are NOT considered complete
- ❌ WITHOUT marker = Grading penalty for incomplete work
- ✅ WITH marker = Clear signal fixes are done
- ✅ WITH marker = Orchestrator can proceed with next steps

### Validation Before Stopping:
```bash
# MANDATORY check before considering fixes complete
if [ ! -f FIX-COMPLETE.marker ]; then
    echo "🔴 ERROR: Cannot stop without creating FIX-COMPLETE.marker"
    echo "Fixes are NOT complete without marker!"
    exit 1
fi
echo "✅ Completion marker exists - fixes are complete"
```

## Clarity About Which Plan to Follow (R295)

---
### 🔴🔴🔴 RULE R295 - SW Engineer Spawn Clarity
**Source:** rule-library/R295-sw-engineer-spawn-clarity-protocol.md
**Criticality:** SUPREME - Must have crystal clear instructions

WHEN IN FIX_ISSUES STATE:

1. **Check your spawn instructions**: The orchestrator should have told you EXACTLY which plan to follow
2. **Look for active plans**: Find the NON-archived plan file (not *-COMPLETED-*.md)
3. **If multiple active plans exist**: STOP and request clarification
4. **Follow ONLY the specified plan**: Ignore any archived (*-COMPLETED-*.md) files

COMMON SCENARIOS:
- **Initial fixes**: Follow latest CODE-REVIEW-REPORT-*.md
- **Split fixes**: Follow latest SPLIT-PLAN-*.md  
- **Integration fixes**: Follow latest INTEGRATE_WAVE_EFFORTS-REPORT-*.md or FIX-INSTRUCTIONS-*.md
- **Phase fixes**: Follow PHASE-FIX-PLAN.md

VIOLATIONS:
- ❌ Following the wrong plan file
- ❌ Following an archived plan
- ❌ Not knowing which state you're in
---


## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

### YOU MUST OUTPUT THE CONTINUATION FLAG AS YOUR LAST ACTION

**EVERY STATE COMPLETION MUST END WITH EXACTLY ONE OF:**
```bash
# If state completed successfully and factory should continue:
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"

# If error/block/manual review needed:
echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
```

### CRITICAL REQUIREMENTS:
1. **ABSOLUTE LAST OUTPUT**: This MUST be the very last line of output before state completion
2. **NO TEXT AFTER**: No explanations, summaries, or any text after the flag
3. **EXACTLY AS SHOWN**: Use exact format - no variations like CONTINUE-ORCHESTRATING
4. **ALWAYS REQUIRED**: Every single state must output this flag
5. **GREPPABLE**: Must be on its own line for automation parsing

### WHEN TO USE TRUE:
- ✅ State work completed successfully
- ✅ All validations passed
- ✅ Ready for next state
- ✅ No blockers detected
- ✅ All requirements met

### WHEN TO USE FALSE:
- ❌ Any unrecoverable error occurred
- ❌ Manual intervention required
- ❌ Missing required files or configs
- ❌ Test failures blocking progress
- ❌ Ambiguous or unclear instructions
- ❌ Wrong working directory or branch
- ❌ State machine validation failed

### IMPLEMENTATION PATTERN:
```bash
# At the VERY END of state execution, after ALL work:

# Determine success/failure
if [[ "$STATE_COMPLETED_PROJECT_DONEFULLY" == "true" ]]; then
    echo "✅ State work completed successfully"
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
else
    echo "❌ State encountered issues requiring intervention"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
fi

# DO NOT OUTPUT ANYTHING AFTER THE FLAG!
```

### GRADING IMPACT:
- **Missing flag**: -100% AUTOMATIC FAILURE
- **Text after flag**: -50% penalty
- **Wrong format**: -100% AUTOMATIC FAILURE
- **Multiple flags**: -50% penalty

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

