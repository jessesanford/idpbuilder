# SW Engineer - IMPLEMENTATION State Grading

## Critical Performance Metrics

---
### ℹ️ RULE  - 
**Source:** 
**Criticality:** INFO - Best practice

PRIMARY METRIC: Implementation Velocity
Measurement: Lines of quality code per hour
Target: 80-120 lines/hour (including tests)
Grade: Based on quality-adjusted velocity
Weight: 40% of overall implementation grade
Consequence: Low velocity = effort timeline risk
---

## Grading Rubric

| Metric | Excellent | Good | Acceptable | FAIL |
|--------|-----------|------|------------|------|
| Lines/Hour | 100-120 | 80-100 | 60-80 | <60 |
| Test Coverage | >90% | 85-90% | 80-85% | <80% |
| Size Management | Always <700 lines | <750 lines | <800 lines | >800 lines |
| Plan Adherence | 100% on track | 95% on track | 90% on track | <90% |
| Code Quality | 0 issues | 1-2 minor | 3-5 minor | Major issues |
| Commit Frequency | Every 1-2 hours | Every 2-3 hours | Every 4 hours | >4 hours |

## Real-Time Scoring

```python
class ImplementationGrader:
    def __init__(self):
        self.implementation_sessions = []
        self.velocity_history = []
        
    def grade_implementation_session(self, session_data):
        """Grade an implementation work session"""
        
        # Core velocity measurement
        velocity_grade = self.calculate_velocity_grade(session_data)
        
        # Test coverage achievement
        coverage_grade = self.calculate_coverage_grade(session_data)
        
        # Size management effectiveness
        size_management_grade = self.calculate_size_management_grade(session_data)
        
        # Implementation plan adherence
        plan_adherence_grade = self.calculate_plan_adherence_grade(session_data)
        
        # Code quality assessment
        quality_grade = self.calculate_code_quality_grade(session_data)
        
        # Commit and progress discipline
        discipline_grade = self.calculate_discipline_grade(session_data)
        
        overall = self.calculate_overall_implementation_grade(
            velocity_grade, coverage_grade, size_management_grade,
            plan_adherence_grade, quality_grade, discipline_grade
        )
        
        return {
            'velocity': velocity_grade,
            'test_coverage': coverage_grade,
            'size_management': size_management_grade,
            'plan_adherence': plan_adherence_grade,
            'code_quality': quality_grade,
            'discipline': discipline_grade,
            'overall': overall,
            'timestamp': datetime.now().isoformat()
        }
    
    def calculate_velocity_grade(self, session):
        """Calculate implementation velocity grade"""
        
        duration_hours = session.get('duration_hours', 1)
        lines_added = session.get('lines_added', 0)
        lines_tested = session.get('test_lines_added', 0)
        
        # Quality-adjusted lines (implementation + test lines)
        quality_lines = lines_added + (lines_tested * 0.7)  # Tests count as 70% toward velocity
        velocity = quality_lines / duration_hours if duration_hours > 0 else 0
        
        # Adjust for code quality issues
        quality_multiplier = self.get_quality_multiplier(session)
        adjusted_velocity = velocity * quality_multiplier
        
        if adjusted_velocity >= 100:
            grade = 'EXCELLENT'
            score = min(100, 80 + (adjusted_velocity - 100) / 5)  # Bonus for high velocity
        elif adjusted_velocity >= 80:
            grade = 'GOOD'
            score = 70 + (adjusted_velocity - 80) / 2
        elif adjusted_velocity >= 60:
            grade = 'ACCEPTABLE'
            score = 60 + (adjusted_velocity - 60) / 2
        else:
            grade = 'FAIL'
            score = max(0, adjusted_velocity)
        
        return {
            'raw_velocity': velocity,
            'quality_adjusted_velocity': adjusted_velocity,
            'quality_multiplier': quality_multiplier,
            'grade': grade,
            'score': score,
            'target_range': '80-120 lines/hour'
        }
    
    def get_quality_multiplier(self, session):
        """Calculate quality multiplier based on code quality indicators"""
        
        multiplier = 1.0
        
        # Reduce for linting issues
        linting_issues = session.get('linting_issues', 0)
        if linting_issues > 0:
            multiplier *= max(0.7, 1.0 - (linting_issues * 0.05))
        
        # Reduce for failing tests
        failing_tests = session.get('failing_tests', 0)
        if failing_tests > 0:
            multiplier *= max(0.5, 1.0 - (failing_tests * 0.1))
        
        # Reduce for build failures
        if session.get('build_failures', 0) > 0:
            multiplier *= 0.6
        
        # Bonus for clean, well-structured code
        if session.get('refactoring_performed', False):
            multiplier *= 1.1
        
        return max(0.3, min(1.2, multiplier))  # Keep within reasonable bounds
    
    def calculate_coverage_grade(self, session):
        """Calculate test coverage achievement grade"""
        
        current_coverage = session.get('current_coverage_percentage', 0)
        target_coverage = session.get('target_coverage_percentage', 85)
        coverage_delta = session.get('coverage_increase', 0)
        
        # Base score on current coverage
        if current_coverage >= 90:
            base_score = 100
            grade = 'EXCELLENT'
        elif current_coverage >= 85:
            base_score = 90
            grade = 'GOOD'
        elif current_coverage >= 80:
            base_score = 75
            grade = 'ACCEPTABLE'
        else:
            base_score = max(0, current_coverage)
            grade = 'FAIL'
        
        # Adjust for improvement trend
        improvement_bonus = min(10, coverage_delta * 2)  # Up to 10 points bonus
        final_score = min(100, base_score + improvement_bonus)
        
        return {
            'current_coverage': current_coverage,
            'target_coverage': target_coverage,
            'coverage_delta': coverage_delta,
            'base_score': base_score,
            'improvement_bonus': improvement_bonus,
            'grade': grade,
            'score': final_score
        }
    
    def calculate_size_management_grade(self, session):
        """Calculate size management effectiveness grade"""
        
        current_lines = session.get('current_total_lines', 0)
        lines_limit = session.get('lines_limit', 800)
        measurement_frequency = session.get('size_checks_performed', 0)
        session_duration = session.get('duration_hours', 1)
        
        # Calculate current size status
        size_utilization = (current_lines / lines_limit) * 100
        
        # Grade based on size utilization
        if current_lines > lines_limit:
            grade = 'FAIL'
            score = 0
        elif size_utilization > 93.75:  # >750 lines
            grade = 'FAIL'
            score = 30  # Some credit for being close but over warning threshold
        elif size_utilization > 87.5:  # >700 lines
            grade = 'ACCEPTABLE'
            score = 60 + (10 * (1 - (size_utilization - 87.5) / 6.25))  # 60-70 range
        elif size_utilization > 75:  # >600 lines
            grade = 'GOOD'
            score = 80 + (10 * (1 - (size_utilization - 75) / 12.5))  # 80-90 range
        else:
            grade = 'EXCELLENT'
            score = 95 + (5 * (1 - size_utilization / 75))  # 95-100 range
        
        # Adjust for measurement discipline
        expected_measurements = max(1, session_duration * 0.5)  # Every 2 hours minimum
        measurement_discipline = min(1.0, measurement_frequency / expected_measurements)
        
        final_score = score * (0.8 + 0.2 * measurement_discipline)  # Up to 20% penalty for poor measurement
        
        return {
            'current_lines': current_lines,
            'lines_limit': lines_limit,
            'size_utilization': size_utilization,
            'measurement_frequency': measurement_frequency,
            'expected_measurements': expected_measurements,
            'measurement_discipline': measurement_discipline,
            'grade': grade,
            'score': final_score
        }
    
    def calculate_plan_adherence_grade(self, session):
        """Calculate implementation plan adherence grade"""
        
        planned_tasks = session.get('planned_tasks_count', 1)
        completed_tasks = session.get('completed_tasks_count', 0)
        out_of_order_tasks = session.get('out_of_order_completions', 0)
        scope_deviations = session.get('scope_deviations', 0)
        
        # Base completion percentage
        completion_percentage = (completed_tasks / planned_tasks) * 100 if planned_tasks > 0 else 0
        
        # Calculate adherence penalties
        order_penalty = out_of_order_tasks * 5  # 5% penalty per out-of-order task
        scope_penalty = scope_deviations * 10   # 10% penalty per scope deviation
        
        adherence_score = max(0, completion_percentage - order_penalty - scope_penalty)
        
        if adherence_score >= 95 and scope_deviations == 0:
            grade = 'EXCELLENT'
            score = 100
        elif adherence_score >= 90 and scope_deviations <= 1:
            grade = 'GOOD'
            score = 85 + (adherence_score - 90)
        elif adherence_score >= 80 and scope_deviations <= 2:
            grade = 'ACCEPTABLE'
            score = 70 + (adherence_score - 80) / 2
        else:
            grade = 'FAIL'
            score = max(0, adherence_score / 2)
        
        return {
            'completion_percentage': completion_percentage,
            'out_of_order_tasks': out_of_order_tasks,
            'scope_deviations': scope_deviations,
            'order_penalty': order_penalty,
            'scope_penalty': scope_penalty,
            'adherence_score': adherence_score,
            'grade': grade,
            'score': score
        }
    
    def calculate_code_quality_grade(self, session):
        """Calculate code quality grade"""
        
        quality_metrics = {
            'linting_issues': session.get('linting_issues', 0),
            'cyclomatic_complexity': session.get('avg_cyclomatic_complexity', 5),
            'test_failures': session.get('test_failures', 0),
            'build_failures': session.get('build_failures', 0),
            'code_smells': session.get('code_smells_detected', 0)
        }
        
        # Start with perfect score
        score = 100
        grade_points = []
        
        # Deduct for linting issues
        if quality_metrics['linting_issues'] > 0:
            penalty = min(30, quality_metrics['linting_issues'] * 3)
            score -= penalty
            grade_points.append(f"Linting issues: -{penalty}")
        
        # Deduct for complexity
        if quality_metrics['cyclomatic_complexity'] > 10:
            penalty = min(20, (quality_metrics['cyclomatic_complexity'] - 10) * 2)
            score -= penalty
            grade_points.append(f"High complexity: -{penalty}")
        
        # Deduct for test failures
        if quality_metrics['test_failures'] > 0:
            penalty = min(25, quality_metrics['test_failures'] * 5)
            score -= penalty
            grade_points.append(f"Test failures: -{penalty}")
        
        # Deduct for build failures
        if quality_metrics['build_failures'] > 0:
            penalty = 20
            score -= penalty
            grade_points.append(f"Build failures: -{penalty}")
        
        # Deduct for code smells
        if quality_metrics['code_smells'] > 0:
            penalty = min(15, quality_metrics['code_smells'] * 2)
            score -= penalty
            grade_points.append(f"Code smells: -{penalty}")
        
        score = max(0, score)
        
        if score >= 95:
            grade = 'EXCELLENT'
        elif score >= 85:
            grade = 'GOOD'
        elif score >= 70:
            grade = 'ACCEPTABLE'
        else:
            grade = 'FAIL'
        
        return {
            'quality_metrics': quality_metrics,
            'deductions': grade_points,
            'grade': grade,
            'score': score
        }
    
    def calculate_discipline_grade(self, session):
        """Calculate development discipline grade"""
        
        duration_hours = session.get('duration_hours', 1)
        commits_made = session.get('commits_made', 0)
        work_log_updates = session.get('work_log_updates', 0)
        size_checks = session.get('size_checks_performed', 0)
        
        discipline_score = 100
        
        # Commit frequency discipline
        expected_commits = max(1, duration_hours / 2)  # Every 2 hours
        if commits_made < expected_commits:
            commit_penalty = min(20, (expected_commits - commits_made) * 10)
            discipline_score -= commit_penalty
        
        # Work log discipline
        expected_log_updates = max(1, duration_hours)  # Every hour
        if work_log_updates < expected_log_updates:
            log_penalty = min(15, (expected_log_updates - work_log_updates) * 7.5)
            discipline_score -= log_penalty
        
        # Size monitoring discipline
        expected_size_checks = max(1, duration_hours / 2)  # Every 2 hours
        if size_checks < expected_size_checks:
            size_penalty = min(25, (expected_size_checks - size_checks) * 12.5)
            discipline_score -= size_penalty
        
        discipline_score = max(0, discipline_score)
        
        if discipline_score >= 95:
            grade = 'EXCELLENT'
        elif discipline_score >= 85:
            grade = 'GOOD'
        elif discipline_score >= 70:
            grade = 'ACCEPTABLE'
        else:
            grade = 'FAIL'
        
        return {
            'commits_made': commits_made,
            'expected_commits': expected_commits,
            'work_log_updates': work_log_updates,
            'expected_log_updates': expected_log_updates,
            'size_checks': size_checks,
            'expected_size_checks': expected_size_checks,
            'discipline_score': discipline_score,
            'grade': grade,
            'score': discipline_score
        }
    
    def calculate_overall_implementation_grade(self, velocity, coverage, size_mgmt, plan, quality, discipline):
        """Calculate weighted overall implementation grade"""
        
        # Weighted scoring
        # Velocity: 30% (most important for delivery)
        # Quality: 25% (code must be maintainable)
        # Size Management: 20% (critical constraint)
        # Test Coverage: 15% (quality assurance)
        # Plan Adherence: 7% (process compliance)
        # Discipline: 3% (professional practices)
        
        weighted_score = (
            velocity['score'] * 0.30 +
            quality['score'] * 0.25 +
            size_mgmt['score'] * 0.20 +
            coverage['score'] * 0.15 +
            plan['score'] * 0.07 +
            discipline['score'] * 0.03
        )
        
        # Critical failure conditions
        critical_failures = []
        if size_mgmt['grade'] == 'FAIL':
            critical_failures.append('Size limit exceeded')
        if quality['score'] < 50:
            critical_failures.append('Code quality unacceptable')
        if coverage['score'] < 60:
            critical_failures.append('Test coverage insufficient')
        
        # Determine final grade
        if critical_failures:
            overall_grade = 'FAIL'
        elif weighted_score >= 90:
            overall_grade = 'EXCELLENT'
        elif weighted_score >= 80:
            overall_grade = 'GOOD'
        elif weighted_score >= 70:
            overall_grade = 'PASS'
        else:
            overall_grade = 'FAIL'
        
        return {
            'weighted_score': weighted_score,
            'grade': overall_grade,
            'critical_failures': critical_failures,
            'has_critical_failures': len(critical_failures) > 0
        }
```

## Velocity Tracking

```python
class VelocityTracker:
    def __init__(self):
        self.session_history = []
        
    def track_implementation_session(self, start_time, end_time, work_data):
        """Track a complete implementation session"""
        
        duration_hours = (end_time - start_time).total_seconds() / 3600
        
        session_metrics = {
            'session_id': f"impl_{start_time.strftime('%Y%m%d_%H%M%S')}",
            'start_time': start_time.isoformat(),
            'end_time': end_time.isoformat(),
            'duration_hours': duration_hours,
            
            # Code metrics
            'lines_added': work_data.get('lines_added', 0),
            'lines_modified': work_data.get('lines_modified', 0),
            'lines_deleted': work_data.get('lines_deleted', 0),
            'net_lines': work_data.get('lines_added', 0) - work_data.get('lines_deleted', 0),
            
            # Test metrics
            'test_lines_added': work_data.get('test_lines_added', 0),
            'test_cases_added': work_data.get('test_cases_added', 0),
            'coverage_before': work_data.get('coverage_before', 0),
            'coverage_after': work_data.get('coverage_after', 0),
            
            # Quality metrics
            'linting_issues': work_data.get('linting_issues', 0),
            'build_failures': work_data.get('build_failures', 0),
            'test_failures': work_data.get('test_failures', 0),
            
            # Process metrics
            'commits_made': work_data.get('commits_made', 0),
            'size_checks': work_data.get('size_checks', 0),
            'work_log_updates': work_data.get('work_log_updates', 0)
        }
        
        # Calculate derived metrics
        session_metrics['velocity_lines_per_hour'] = session_metrics['net_lines'] / duration_hours if duration_hours > 0 else 0
        session_metrics['test_velocity'] = session_metrics['test_lines_added'] / duration_hours if duration_hours > 0 else 0
        session_metrics['quality_score'] = self.calculate_quality_score(session_metrics)
        
        self.session_history.append(session_metrics)
        return session_metrics
    
    def calculate_quality_score(self, metrics):
        """Calculate quality score for a session"""
        
        base_score = 100
        
        # Penalize quality issues
        base_score -= min(20, metrics['linting_issues'] * 2)
        base_score -= min(30, metrics['build_failures'] * 15)
        base_score -= min(25, metrics['test_failures'] * 5)
        
        # Reward good practices
        if metrics['test_lines_added'] > 0:
            base_score += 5
        if metrics['commits_made'] >= metrics['duration_hours'] / 2:
            base_score += 3
        
        return max(0, min(100, base_score))
    
    def get_velocity_trends(self, lookback_sessions=5):
        """Analyze velocity trends over recent sessions"""
        
        if len(self.session_history) < 2:
            return {'trend': 'INSUFFICIENT_DATA', 'sessions_analyzed': len(self.session_history)}
        
        recent_sessions = self.session_history[-lookback_sessions:]
        velocities = [s['velocity_lines_per_hour'] for s in recent_sessions]
        quality_scores = [s['quality_score'] for s in recent_sessions]
        
        # Calculate trends
        velocity_trend = self.calculate_trend(velocities)
        quality_trend = self.calculate_trend(quality_scores)
        
        # Overall assessment
        if velocity_trend['direction'] == 'IMPROVING' and quality_trend['direction'] != 'DECLINING':
            overall_trend = 'POSITIVE'
        elif velocity_trend['direction'] == 'DECLINING' or quality_trend['direction'] == 'DECLINING':
            overall_trend = 'CONCERNING'
        else:
            overall_trend = 'STABLE'
        
        return {
            'overall_trend': overall_trend,
            'velocity_trend': velocity_trend,
            'quality_trend': quality_trend,
            'avg_velocity': sum(velocities) / len(velocities),
            'avg_quality': sum(quality_scores) / len(quality_scores),
            'sessions_analyzed': len(recent_sessions)
        }
    
    def calculate_trend(self, values):
        """Calculate trend direction from a series of values"""
        
        if len(values) < 3:
            return {'direction': 'INSUFFICIENT_DATA', 'slope': 0}
        
        # Simple linear regression slope
        n = len(values)
        x_sum = sum(range(n))
        y_sum = sum(values)
        xy_sum = sum(i * values[i] for i in range(n))
        x2_sum = sum(i * i for i in range(n))
        
        slope = (n * xy_sum - x_sum * y_sum) / (n * x2_sum - x_sum * x_sum)
        
        if slope > 0.1:
            direction = 'IMPROVING'
        elif slope < -0.1:
            direction = 'DECLINING'
        else:
            direction = 'STABLE'
        
        return {'direction': direction, 'slope': slope}
```

## Performance Dashboard

```python
def generate_implementation_dashboard(session_data):
    """Generate real-time implementation performance dashboard"""
    
    grader = ImplementationGrader()
    tracker = VelocityTracker()
    
    # Grade current session
    current_grade = grader.grade_implementation_session(session_data)
    
    # Get velocity trends
    velocity_trends = tracker.get_velocity_trends()
    
    dashboard = {
        'current_session': current_grade,
        'velocity_trends': velocity_trends,
        'health_status': assess_implementation_health(current_grade, velocity_trends),
        'recommendations': generate_performance_recommendations(current_grade, velocity_trends)
    }
    
    print("📊 IMPLEMENTATION PERFORMANCE DASHBOARD")
    print(f"Overall Grade: {current_grade['overall']['grade']} ({current_grade['overall']['weighted_score']:.1f}/100)")
    print(f"Velocity: {current_grade['velocity']['quality_adjusted_velocity']:.1f} lines/hour (target: 80-120)")
    print(f"Size Status: {current_grade['size_management']['current_lines']}/{current_grade['size_management']['lines_limit']} lines ({current_grade['size_management']['size_utilization']:.1f}%)")
    print(f"Test Coverage: {current_grade['test_coverage']['current_coverage']:.1f}%")
    print(f"Code Quality: {current_grade['code_quality']['grade']} ({current_grade['code_quality']['score']}/100)")
    
    if current_grade['overall']['critical_failures']:
        print("❌ CRITICAL FAILURES:")
        for failure in current_grade['overall']['critical_failures']:
            print(f"  - {failure}")
    
    if velocity_trends['overall_trend'] == 'CONCERNING':
        print("⚠️ VELOCITY CONCERN: Performance trending downward")
    
    return dashboard

def assess_implementation_health(current_grade, trends):
    """Assess overall implementation health"""
    
    health_indicators = {
        'performance': current_grade['overall']['grade'],
        'velocity_trend': trends['overall_trend'],
        'critical_issues': len(current_grade['overall']['critical_failures']) > 0,
        'size_risk': current_grade['size_management']['size_utilization'] > 90
    }
    
    if health_indicators['critical_issues'] or health_indicators['size_risk']:
        status = 'CRITICAL'
    elif health_indicators['velocity_trend'] == 'CONCERNING':
        status = 'WARNING'
    elif health_indicators['performance'] in ['EXCELLENT', 'GOOD']:
        status = 'HEALTHY'
    else:
        status = 'NEEDS_IMPROVEMENT'
    
    return {
        'status': status,
        'indicators': health_indicators
    }
```

## Warning Triggers

---
### 🚨 RULE  - 
**Source:** 
**Criticality:** CRITICAL - Major impact on grading

IMPLEMENTATION STATE WARNINGS
Velocity Below 60 lines/hour:
⚠️ WARNING: Implementation velocity too low
⚠️ Review complexity, tooling, or approach

Size >750 lines:
🚨 CRITICAL: Immediate size limit danger
🚨 Stop implementation and transition to MEASURE_SIZE

Test Coverage <80%:
⚠️⚠️ WARNING: Test coverage insufficient
⚠️⚠️ Focus on test writing before continuing

>4 hours without commits:
⚠️ WARNING: Commit discipline issues
⚠️ May indicate blocked or unfocused work
---

## Performance State Tracking

```yaml
# Update orchestrator-state-v3.json
grading:
  SW_ENGINEER:
    IMPLEMENTATION:
      latest:
        timestamp: "2025-08-23T16:30:45Z"
        session_duration_hours: 2.5
        velocity_lines_per_hour: 94.2
        test_coverage: 88.5
        size_utilization: 67.3
        overall: "GOOD"
        
      cumulative:
        sessions_completed: 8
        total_hours: 18.5
        total_lines_delivered: 1847
        avg_velocity: 89.7
        excellent: 3
        good: 4
        acceptable: 1
        fail: 0
