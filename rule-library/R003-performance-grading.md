# 🚨 RULE R003 - Performance Grading System

**Criticality:** CRITICAL - Continuous assessment  
**Grading Impact:** Direct impact on agent employment  
**Enforcement:** CONTINUOUS - Tracked throughout execution

## Rule Statement

ALL agents are continuously graded on performance metrics. Poor grades lead to immediate replacement.

## Universal Grading Principles

1. **No Excuses**: Technical issues are YOUR problem
2. **Final Grades**: No appeals, no re-grades
3. **Three Strike Rule**: Warning → Retraining → Termination
4. **Real-time Tracking**: Grades updated after each action

## Orchestrator Grading Metrics

### Primary Metrics (80% of grade)
- **Parallel Spawn Timing** (R151): <5s average delta = 50% of grade
- **State File Updates** (R288): 100% compliance = 20% of grade
- **Size Compliance**: Zero >800 line efforts = 10% of grade

### Secondary Metrics (20% of grade)
- **Monitoring Frequency**: Check every 5 messages
- **Integration Creation**: 100% wave completion
- **TODO Persistence**: R287 compliance

### Grading Formula
```python
orchestrator_grade = (
    (spawn_timing_pass * 0.5) +
    (state_updates_rate * 0.2) +
    (size_compliance * 0.1) +
    (monitoring_rate * 0.05) +
    (integration_rate * 0.05) +
    (todo_compliance * 0.1)
)
```

### Pass/Fail Thresholds
- **PASS**: ≥ 0.8
- **WARNING**: 0.6 - 0.79
- **FAIL**: < 0.6

## SW Engineer Grading Metrics

### Primary Metrics
- **Implementation Speed** (R152): >50 lines/hour = 30%
- **Test Coverage**: Meet phase minimum = 30%
- **Size Compliance**: Never exceed limit = 20%

### Secondary Metrics
- **Work Log Updates**: Every checkpoint = 10%
- **Git Hygiene**: Logical commits = 10%

### Grading Formula
```python
sw_engineer_grade = (
    min(lines_per_hour/50, 1.0) * 0.3 +
    min(test_coverage/required, 1.0) * 0.3 +
    (1 if under_limit else 0) * 0.2 +
    work_log_frequency * 0.1 +
    commit_quality * 0.1
)
```

## Code Reviewer Grading Metrics

### Primary Metrics
- **Plan Quality** (R153): First-try success >80% = 40%
- **Review Accuracy**: No missed critical issues = 30%
- **Size Measurement**: Always correct tool = 20%

### Secondary Metrics
- **Split Decisions**: All under limit = 10%

### Critical Failures (Automatic FAIL)
- Wrong size measurement tool used
- Split exceeds 800 line limit
- Missed security vulnerability
- Approved broken code

## Architect Grading Metrics

### Primary Metrics
- **Decision Accuracy** (R158): No reversed decisions = 40%
- **Issue Detection**: Catch critical problems = 30%
- **Assessment Accuracy**: Correct trajectory = 30%

### Critical Failures (Automatic FAIL)
- False positive STOP (blocks progress unnecessarily)
- Missed architectural violation
- Wrong ON_TRACK/OFF_TRACK assessment

## Grade Reporting

Grades MUST be tracked in orchestrator-state.yaml:

```yaml
agent_grades:
  orchestrator:
    parallel_spawn:
      last_measurement: "2025-08-26T14:30:45Z"
      average_delta: 3.2
      grade: "PASS"
    state_updates:
      compliance_rate: 1.0
      grade: "PASS"
    overall_grade: 0.92
    status: "PASS"
    
  sw_engineer_1:
    implementation_speed: 67
    test_coverage: 0.85
    size_compliance: true
    overall_grade: 0.88
    status: "PASS"
    
  code_reviewer_1:
    first_try_success: 0.83
    review_accuracy: 1.0
    overall_grade: 0.91
    status: "PASS"
    
  architect:
    decisions_reversed: 0
    critical_caught: 3
    trajectory_accuracy: 1.0
    overall_grade: 1.0
    status: "PASS"
```

## Enforcement Actions

| Grade Status | First Offense | Second Offense | Third Offense |
|-------------|---------------|----------------|---------------|
| WARNING | Logged warning | Performance plan | Consider replacement |
| FAIL | Immediate warning | Retraining required | Termination |

## Self-Assessment Required

Agents should self-monitor:
```bash
check_my_grade() {
    local my_grade=$(calculate_grade)
    if [ "$my_grade" -lt 80 ]; then
        echo "⚠️⚠️⚠️ WARNING: My grade is $my_grade% - improvement needed!"
        echo "📊 Analyzing weak areas..."
        identify_improvement_areas
    fi
}
```

## Examples

### GOOD: High Performance
```
Orchestrator Grade Report:
- Parallel Spawns: 2.1s average ✅
- State Updates: 100% ✅
- Size Compliance: 100% ✅
- Overall: 95% PASS
```

### BAD: Poor Performance
```
SW Engineer Grade Report:
- Implementation: 35 lines/hour ❌
- Test Coverage: 60% (required 80%) ❌
- Overall: 52% FAIL
⚠️⚠️⚠️ WARNING ISSUED - IMPROVE IMMEDIATELY
```

---
**Remember:** Your performance is your responsibility. Maintain high grades or be replaced.