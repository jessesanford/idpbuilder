# 📊 GRADING METRICS QUICK REFERENCE

## 🚨 GRADING SYSTEM OVERVIEW

```
┌─────────────────────────────────────────────────────────────────┐
│                    AGENT PERFORMANCE GRADING                   │
│ Poor performance = Replacement | Grades are final & binding   │
└─────────────────────────────────────────────────────────────────┘
```

### Strike System:
| Grade | First Offense | Second Offense | Third Offense |
|-------|---------------|----------------|---------------|
| FAIL  | Warning       | Retraining     | **Termination** |

## 🎯 ORCHESTRATOR GRADING METRICS

### Primary Metrics:
```
┌─────────────────────────────────────────────┐
│ 🚨 PARALLEL SPAWN TIMING (R151.0.0)        │
│ Target: <5 seconds average delta            │
│ Weight: CRITICAL (can cause immediate FAIL) │
│ Formula: avg(spawn_timestamps_delta)        │
└─────────────────────────────────────────────┘
```

### Grading Formula:
```python
def grade_orchestrator_performance():
    # Critical Metrics (Must Pass)
    parallel_spawn_grade = "PASS" if avg_spawn_delta < 5.0 else "FAIL"
    state_file_updates = "PASS" if updates_after_transitions == 100 else "FAIL"  
    agent_monitoring = "PASS" if check_frequency <= 5_messages else "FAIL"
    gate_enforcement = "PASS" if compliance_rate == 100 else "FAIL"
    
    # Any FAIL in critical metrics = Overall FAIL
    if any([parallel_spawn_grade, state_file_updates, agent_monitoring, gate_enforcement] == "FAIL"):
        return "FAIL"
    
    # Additional Performance Metrics
    completion_rate = completed_efforts / total_efforts
    integration_success = successful_integrations / total_integrations
    
    overall_score = (completion_rate * 0.3 + 
                    integration_success * 0.3 +
                    efficiency_score * 0.4)
    
    return "PASS" if overall_score >= 0.85 else "FAIL"
```

### Measurement Commands:
```bash
# Track spawn timing (CRITICAL!)
start_time=$(date +%s.%N)
# ... spawn agents in parallel ...
end_time=$(date +%s.%N)
delta=$(echo "$end_time - $start_time" | bc)
echo "Spawn delta: ${delta}s (target: <5s)"
```

## ⚡ SW ENGINEER GRADING METRICS

### Primary Metrics:
```
┌─────────────────────────────────────────────┐
│ 🚨 IMPLEMENTATION SPEED (R152.0.0)          │
│ Target: >50 lines/hour                      │
│ Weight: 30% of total score                  │
│ Measurement: (total_lines / hours_worked)   │
└─────────────────────────────────────────────┘
```

### Grading Formula:
```python
def grade_sw_engineer_performance():
    # Core Metrics  
    lines_per_hour = total_lines_implemented / hours_worked
    test_coverage_ratio = actual_coverage / required_coverage
    size_compliance = 1.0 if max_branch_size <= 800 else 0.0
    work_log_frequency = updates_per_checkpoint / expected_updates
    commit_quality = logical_commits / total_commits
    
    # Weighted Score
    score = (
        (lines_per_hour / 50.0) * 0.3 +      # Speed component
        test_coverage_ratio * 0.3 +           # Quality component  
        size_compliance * 0.2 +               # Compliance component
        work_log_frequency * 0.1 +            # Process component
        commit_quality * 0.1                  # Hygiene component
    )
    
    # Critical Failures (Instant FAIL)
    if max_branch_size > 800:
        return "FAIL"  # Size violation
    if actual_coverage < required_coverage:
        return "FAIL"  # Coverage violation
        
    return "PASS" if score >= 0.8 else "FAIL"
```

### Test Coverage Requirements:
| Phase | Minimum Coverage | Critical Areas |
|-------|------------------|----------------|
| 1     | 70%             | Controllers (>90%) |
| 2     | 75%             | Webhooks (>90%) |
| 3     | 80%             | Integration (>95%) |

### Measurement Commands:
```bash
# Line count (MANDATORY TOOL!)
lines=$(line-counter.sh -c ${BRANCH} | grep "Total:" | awk '{print $2}')

# Test coverage
go test -coverprofile=coverage.out ./...
coverage=$(go tool cover -func=coverage.out | grep "total:" | awk '{print $3}' | sed 's/%//')
```

## 🔍 CODE REVIEWER GRADING METRICS

### Primary Metrics:
```
┌─────────────────────────────────────────────┐
│ 🚨 FIRST-TRY SUCCESS RATE (R153.0.0)        │
│ Target: >80% implementation success         │
│ Weight: CRITICAL (core competency)          │
│ Formula: successful_first_attempts / total  │
└─────────────────────────────────────────────┘
```

### Grading Formula:
```python
def grade_code_reviewer_performance():
    # Critical Metrics (Must Pass)
    first_try_success = successful_implementations / total_plans
    missed_critical_issues = count_missed_critical_issues()
    wrong_size_measurements = count_wrong_measurement_tool_usage()
    split_size_violations = count_splits_over_800_lines()
    
    # Instant FAIL conditions
    if first_try_success < 0.8:
        return "FAIL"
    if missed_critical_issues > 0:
        return "FAIL"  
    if wrong_size_measurements > 0:
        return "FAIL"
    if split_size_violations > 0:
        return "FAIL"
        
    # Additional Quality Metrics
    plan_completeness = complete_plans / total_plans
    review_thoroughness = issues_found / total_issues
    documentation_quality = clear_feedback_reports / total_reports
    
    score = (
        first_try_success * 0.4 +
        plan_completeness * 0.2 +
        review_thoroughness * 0.2 +
        documentation_quality * 0.2
    )
    
    return "PASS" if score >= 0.85 else "FAIL"
```

### Critical Compliance Checks:
```bash
# Size measurement (MUST use this tool)
if ! command -v line-counter.sh; then
    echo "🚨 GRADING FAIL: Wrong measurement tool"
fi

# Verify all splits under limit
for split in split-*; do
    lines=$(line-counter.sh -c $split | grep "Total:" | awk '{print $2}')
    if [ $lines -gt 800 ]; then
        echo "🚨 GRADING FAIL: Split $split exceeds 800 lines ($lines)"
    fi
done
```

## 🏗️ ARCHITECT GRADING METRICS

### Primary Metrics:
```
┌─────────────────────────────────────────────┐
│ 🚨 DECISION ACCURACY (R158.0.0)             │
│ Target: Zero reversed decisions              │
│ Weight: CRITICAL (architectural integrity)  │
│ Formula: stable_decisions / total_decisions │
└─────────────────────────────────────────────┘
```

### Grading Formula:
```python
def grade_architect_performance():
    # Critical Failure Conditions (Instant FAIL)
    false_positive_stops = count_unnecessary_stop_decisions()
    missed_critical_issues = count_missed_architectural_problems()
    wrong_trajectory_calls = count_incorrect_on_track_off_track()
    unclear_addendums = count_addendums_causing_next_wave_failure()
    
    # Any critical failure = FAIL
    critical_failures = [
        false_positive_stops,
        missed_critical_issues, 
        wrong_trajectory_calls,
        unclear_addendums
    ]
    
    if any(failure > 0 for failure in critical_failures):
        return "FAIL"
    
    # Performance Metrics
    decision_consistency = stable_decisions / total_decisions
    issue_detection_rate = critical_issues_found / total_critical_issues
    feature_assessment_accuracy = correct_trajectory_calls / total_assessments
    addendum_clarity_score = successful_next_waves / waves_with_addendums
    
    score = (
        decision_consistency * 0.3 +
        issue_detection_rate * 0.3 +
        feature_assessment_accuracy * 0.2 +
        addendum_clarity_score * 0.2
    )
    
    return "PASS" if score >= 0.9 else "FAIL"
```

### Critical Validation:
```bash
# Verify all efforts under size limit before PROCEED
for effort_branch in $(get_completed_efforts); do
    lines=$(line-counter.sh -c $effort_branch | grep "Total:" | awk '{print $2}')
    if [ $lines -gt 800 ]; then
        echo "🚨 ARCHITECTURAL FAIL: Cannot PROCEED with oversized effort"
        echo "Decision: STOP required"
    fi
done
```

## 📊 CROSS-AGENT METRICS

### Integration Success Rate:
```
Target: 100% wave integration success
Measurement: clean_merges / total_wave_integrations  
Owner: Orchestrator (coordination) + Architect (approval)
```

### Line Count Compliance:
```
Target: 100% efforts under 800 lines
Measurement: compliant_efforts / total_efforts
Tool: line-counter.sh (MANDATORY)
Owner: All agents (must verify)
```

### Pattern Compliance Rate:
```
Target: 100% KCP pattern adherence
Measurement: compliant_implementations / total_implementations
Owner: Code Reviewer (check) + Architect (validate)
```

## 📈 REAL-TIME GRADING DASHBOARD

### State File Grade Tracking:
```yaml
agent_grades:
  orchestrator:
    parallel_spawn:
      last_measurement: "2025-08-23T14:30:45Z"
      average_delta: 3.2
      grade: "PASS"
    monitoring_frequency:
      messages_between_checks: 4
      grade: "PASS"
    overall: "PASS"
    
  sw_engineer:
    implementation_speed: 67  # lines/hour
    test_coverage: 85        # percentage
    size_compliance: true    # under 800 lines
    work_log_updates: 0.9    # frequency score
    commit_quality: 0.85     # logical commit ratio
    grade: "PASS"
    
  code_reviewer:
    first_try_success: 0.83  # 83% success rate
    critical_misses: 0       # zero misses
    wrong_measurements: 0    # zero violations
    split_violations: 0      # zero oversized splits
    grade: "PASS"
    
  architect:
    decisions_reversed: 0    # zero reversals
    critical_caught: 3       # issues detected
    false_stops: 0          # zero false positives
    unclear_addendums: 0    # zero causing failures
    grade: "PASS"
```

## 🚨 GRADING ENFORCEMENT ACTIONS

### Warning Level (First FAIL):
```
1. Immediate performance review meeting
2. Specific improvement targets set
3. Enhanced monitoring activated
4. Retraining materials provided
5. Probationary period initiated
```

### Retraining Level (Second FAIL):
```
1. Intensive retraining program
2. Supervised execution required
3. Performance metrics tracked hourly
4. Competency re-certification needed
5. Final warning issued
```

### Termination Level (Third FAIL):
```
1. Agent immediately deactivated
2. All work suspended
3. Replacement agent spawned
4. Incident documented
5. Process improvements reviewed
```

## 🔧 GRADING MEASUREMENT TOOLS

### Line Counting (UNIVERSAL):
```bash
# ONLY acceptable tool
$PROJECT_ROOT/tools/line-counter.sh -c ${BRANCH}

# For detailed analysis
$PROJECT_ROOT/tools/line-counter.sh -c ${BRANCH} -d

# NEVER count manually or use other tools
```

### Performance Timing:
```bash
# High-precision timing for critical metrics
start=$(date +%s.%N)
# ... perform timed operation ...
end=$(date +%s.%N)
duration=$(echo "$end - $start" | bc)
```

### Coverage Analysis:
```bash
# Standard Go coverage
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep "total:" | awk '{print $3}'

# Detailed coverage by package
go tool cover -html=coverage.out -o coverage.html
```

---
**REMEMBER**: Grading is continuous, automated, and non-negotiable. Excellence is the only acceptable standard.