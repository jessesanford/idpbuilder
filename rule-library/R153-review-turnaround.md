# 🚨 RULE R153 - Review Turnaround and Quality Metrics

**Criticality:** CRITICAL - Performance metric  
**Grading Impact:** 40% of Code Reviewer grade  
**Enforcement:** PER REVIEW - Tracked continuously

## Rule Statement

Code Reviewers MUST achieve >80% first-try success rate and complete reviews within 30 minutes.

## Review Performance Metrics

### Key Performance Indicators
| Metric | Minimum | Target | Excellent |
|--------|---------|--------|-----------|
| First-Try Success | 80% | 90% | 95%+ |
| Review Completion Time | 30 min | 20 min | 15 min |
| Critical Issues Caught | 100% | 100% | 100% |
| False Positives | <10% | <5% | <2% |
| Plan Quality Score | 80% | 90% | 95%+ |

## First-Try Success Rate

### Definition
Plans created by Code Reviewer that result in successful implementation without requiring rework.

### Calculation
```python
def calculate_first_try_success(reviewer):
    total_plans = reviewer.implementation_plans_created
    successful_first_try = 0
    
    for plan in reviewer.plans:
        if plan.implementation_succeeded_without_rework:
            successful_first_try += 1
    
    success_rate = successful_first_try / total_plans
    return {
        'rate': success_rate * 100,
        'grade': min(success_rate / 0.8, 1.0) * 40  # 40% of grade
    }
```

### What Constitutes Success
- ✅ Implementation completes without blockers
- ✅ No missing dependencies identified during implementation
- ✅ Size estimate accurate (within 20%)
- ✅ All requirements implementable as specified
- ✅ Test requirements achievable

### What Constitutes Failure
- ❌ SW Engineer blocked by unclear requirements
- ❌ Major dependency missed in plan
- ❌ Size estimate off by >50%
- ❌ Requirements technically impossible
- ❌ Plan requires significant revision

## Review Turnaround Time

### Time Limits
```bash
review_time_tracker() {
    local start_time=$1
    local review_type=$2
    
    case "$review_type" in
        "effort_plan")
            MAX_TIME=1800  # 30 minutes
            ;;
        "code_review")
            MAX_TIME=1800  # 30 minutes
            ;;
        "split_plan")
            MAX_TIME=900   # 15 minutes
            ;;
    esac
    
    local current=$(date +%s)
    local elapsed=$((current - start_time))
    
    if [ $elapsed -gt $MAX_TIME ]; then
        echo "⚠️ VIOLATION: Review exceeded time limit!"
        echo "Limit: $MAX_TIME seconds"
        echo "Actual: $elapsed seconds"
    fi
}
```

## Review Quality Standards

### Critical Issues (Must Catch 100%)
```bash
# Security vulnerabilities
- SQL injection risks
- Authentication bypasses
- Exposed secrets/credentials
- Unsafe type assertions

# Architectural violations
- Wrong patterns used
- Circular dependencies
- Breaking changes to APIs
- Incorrect service boundaries

# Size violations
- Effort exceeding 800 lines
- Incorrect measurement method
- Missing split planning
```

### Review Checklist
```markdown
## Code Review Checklist

### Critical (Automatic Fail if Missed)
- [ ] No security vulnerabilities
- [ ] Size under 800 lines (verified with line-counter.sh)
- [ ] No breaking changes without version bump
- [ ] All tests passing

### Important (Points Deduction if Missed)
- [ ] Follows project patterns
- [ ] Adequate test coverage (meets phase minimum)
- [ ] Error handling implemented
- [ ] Documentation updated

### Nice to Have (No Penalty)
- [ ] Performance optimizations
- [ ] Code elegance improvements
- [ ] Additional test cases
```

## Plan Creation Quality

### High-Quality Plan Characteristics
```yaml
implementation_plan:
  clarity:
    requirements: "Crystal clear, no ambiguity"
    success_criteria: "Measurable and specific"
    dependencies: "All identified upfront"
    
  accuracy:
    size_estimate: "Within 20% of actual"
    time_estimate: "Within 30% of actual"
    complexity: "Correctly assessed"
    
  completeness:
    files_to_create: "All listed with purposes"
    test_requirements: "Specific scenarios listed"
    integration_points: "All identified"
    
  structure:
    logical_flow: "Tasks in correct order"
    parallelizable: "Opportunities identified"
    checkpoints: "Clear validation points"
```

## Common Quality Issues

### VIOLATION: Vague Requirements
```markdown
# ❌ BAD: Unclear plan
## Requirements
- Implement user stuff
- Add some validation
- Make it work
```

### VIOLATION: Missing Dependencies
```markdown
# ❌ BAD: Dependencies not identified
## Implementation
- Create UserService
# Missing: Requires AuthService, Database, Config
```

### VIOLATION: Wrong Size Estimate
```markdown
# ❌ BAD: Severe underestimate
## Size Estimate: 200 lines
# Reality: 850 lines (325% off)
```

## Correct Patterns

### GOOD: Clear, Complete Plan
```markdown
## Implementation Plan for User Authentication

### Requirements
1. JWT-based authentication with refresh tokens
2. Role-based access control (admin, user, guest)
3. Session management with Redis
4. Password requirements: 8+ chars, 1 upper, 1 number

### Dependencies
- External: Redis client, JWT library, bcrypt
- Internal: ConfigService, LoggerService
- Database: users table with indexes

### Files to Create
1. `auth/jwt.go` (150 lines) - JWT token management
2. `auth/session.go` (100 lines) - Session handling
3. `auth/rbac.go` (120 lines) - Role checking
4. `auth/password.go` (80 lines) - Password validation
5. Tests: (200 lines total)

### Total Estimate: 650 lines ✅
```

## Review Speed Techniques

### 1. Pattern Recognition
```bash
# Quick pattern checks
check_common_patterns() {
    grep -r "panic\|Fatal" . && echo "⚠️ Unhandled panic found"
    grep -r "TODO\|FIXME" . && echo "⚠️ Unfinished work found"
    grep -r "password.*=.*\"" . && echo "🚨 Hardcoded password!"
}
```

### 2. Automated Checks First
```bash
# Run automated tools before manual review
quick_review() {
    make lint || echo "⚠️ Linting issues"
    make test || echo "⚠️ Test failures"
    line-counter.sh -c branch || echo "⚠️ Check size"
    # Then focus manual review on logic
}
```

## Grading Formula

```python
def grade_code_reviewer(reviewer):
    # 40% - First-try success rate
    success_grade = min(reviewer.first_try_rate / 80, 1.0) * 40
    
    # 30% - Review accuracy (no missed critical issues)
    accuracy_grade = (1.0 - reviewer.missed_critical_rate) * 30
    
    # 20% - Turnaround time
    time_grade = min(30 / reviewer.avg_review_minutes, 1.0) * 20
    
    # 10% - False positive rate
    false_positive_grade = max(0, 1.0 - reviewer.false_positive_rate) * 10
    
    total_grade = success_grade + accuracy_grade + time_grade + false_positive_grade
    
    # Automatic failures
    if reviewer.missed_security_issue:
        return 0  # Immediate failure
    
    if reviewer.approved_oversized_effort:
        return 0  # Immediate failure
        
    return total_grade
```

## Daily Performance Report

```yaml
reviewer_performance:
  date: "2025-08-26"
  reviewer: "code-reviewer-1"
  
  reviews_completed: 5
  
  metrics:
    first_try_success:
      successful: 4
      total: 5
      rate: 80%
      status: "MEETS_MINIMUM"
    
    turnaround:
      average: 22
      fastest: 15
      slowest: 28
      status: "GOOD"
    
    accuracy:
      critical_caught: 2
      critical_missed: 0
      false_positives: 1
      status: "EXCELLENT"
  
  grade_calculation:
    first_try: 32  # (80/80) * 40 = 40
    accuracy: 30    # No misses = full 30
    speed: 18       # (30/22) * 20, capped
    false_pos: 9    # Small penalty
    total: 89       # B+ grade
    
  status: "PASS"
  trend: "stable"
```

---
**Remember:** Quality plans prevent rework. Fast, accurate reviews keep projects moving.