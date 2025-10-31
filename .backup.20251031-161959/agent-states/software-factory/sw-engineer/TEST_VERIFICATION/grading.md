# SW Engineer - TEST_WRITING State Grading

## Critical Performance Metrics

---
### ℹ️ RULE  - 
**Source:** 
**Criticality:** INFO - Best practice

PRIMARY METRIC: Test Coverage Achievement
Measurement: Coverage increase per test writing session
Target: +8-12% coverage increase per hour
Grade: Based on coverage goals achieved
Weight: 50% of overall test writing grade
Consequence: Low coverage achievement = failed testing
---

## Grading Rubric

| Metric | Excellent | Good | Acceptable | FAIL |
|--------|-----------|------|------------|------|
| Coverage Increase | +12%+ per hour | +8-12% per hour | +5-8% per hour | <5% per hour |
| Test Quality Score | 90-100 | 80-90 | 70-80 | <70 |
| Test Case Efficiency | <15 lines/case | 15-20 lines/case | 20-25 lines/case | >25 lines/case |
| Test Failure Rate | 0% failing | <5% failing | 5-10% failing | >10% failing |
| Size Management | No size issues | Minor impact | Moderate impact | Size violation |
| Test Completeness | 100% scenarios | 95% scenarios | 85% scenarios | <85% scenarios |

## Real-Time Scoring

```python
class TestWritingGrader:
    def __init__(self):
        self.test_sessions = []
        self.coverage_history = []
        
    def grade_test_writing_session(self, session_data):
        """Grade a test writing work session"""
        
        # Coverage achievement grading (primary metric)
        coverage_grade = self.calculate_coverage_achievement_grade(session_data)
        
        # Test quality assessment  
        quality_grade = self.calculate_test_quality_grade(session_data)
        
        # Test development efficiency
        efficiency_grade = self.calculate_test_efficiency_grade(session_data)
        
        # Test completeness and scenario coverage
        completeness_grade = self.calculate_test_completeness_grade(session_data)
        
        # Size management during test development
        size_management_grade = self.calculate_test_size_management_grade(session_data)
        
        # Test maintainability and structure
        maintainability_grade = self.calculate_test_maintainability_grade(session_data)
        
        overall = self.calculate_overall_test_writing_grade(
            coverage_grade, quality_grade, efficiency_grade,
            completeness_grade, size_management_grade, maintainability_grade
        )
        
        return {
            'coverage_achievement': coverage_grade,
            'test_quality': quality_grade,
            'test_efficiency': efficiency_grade,
            'test_completeness': completeness_grade,
            'size_management': size_management_grade,
            'maintainability': maintainability_grade,
            'overall': overall,
            'timestamp': datetime.now().isoformat()
        }
    
    def calculate_coverage_achievement_grade(self, session):
        """Calculate grade based on test coverage achievement"""
        
        duration_hours = session.get('duration_hours', 1)
        coverage_before = session.get('coverage_before_percentage', 0)
        coverage_after = session.get('coverage_after_percentage', 0)
        coverage_target = session.get('coverage_target_percentage', 85)
        
        coverage_increase = coverage_after - coverage_before
        coverage_rate = coverage_increase / duration_hours if duration_hours > 0 else 0
        
        # Assess coverage progress toward target
        progress_toward_target = min(100, (coverage_after / coverage_target) * 100) if coverage_target > 0 else 100
        
        # Base scoring on coverage increase rate
        if coverage_rate >= 12:
            rate_score = 100
            rate_grade = 'EXCELLENT'
        elif coverage_rate >= 8:
            rate_score = 80 + (coverage_rate - 8) * 5  # 80-100 range
            rate_grade = 'GOOD'
        elif coverage_rate >= 5:
            rate_score = 60 + (coverage_rate - 5) * 6.67  # 60-80 range  
            rate_grade = 'ACCEPTABLE'
        else:
            rate_score = max(0, coverage_rate * 12)  # 0-60 range
            rate_grade = 'FAIL'
        
        # Adjust for target achievement
        if coverage_after >= coverage_target:
            target_bonus = 10  # Bonus for meeting target
        elif progress_toward_target >= 90:
            target_bonus = 5   # Near target
        else:
            target_bonus = 0
        
        final_score = min(100, rate_score + target_bonus)
        
        # Override grade if target is achieved
        if coverage_after >= coverage_target and rate_grade != 'FAIL':
            if final_score >= 95:
                final_grade = 'EXCELLENT'
            elif final_score >= 85:
                final_grade = 'GOOD'
            else:
                final_grade = 'ACCEPTABLE'
        else:
            final_grade = rate_grade
        
        return {
            'coverage_before': coverage_before,
            'coverage_after': coverage_after,
            'coverage_increase': coverage_increase,
            'coverage_rate_per_hour': coverage_rate,
            'coverage_target': coverage_target,
            'progress_toward_target': progress_toward_target,
            'rate_score': rate_score,
            'target_bonus': target_bonus,
            'grade': final_grade,
            'score': final_score
        }
    
    def calculate_test_quality_grade(self, session):
        """Calculate grade based on test code quality"""
        
        quality_metrics = session.get('test_quality_metrics', {})
        
        # Test structure quality (0-25 points)
        structure_score = self.assess_test_structure_quality(quality_metrics)
        
        # Assertion quality (0-25 points)  
        assertion_score = self.assess_assertion_quality(quality_metrics)
        
        # Test coverage depth (0-25 points)
        depth_score = self.assess_coverage_depth_quality(quality_metrics)
        
        # Test maintainability (0-25 points)
        maintainability_score = self.assess_test_maintainability_quality(quality_metrics)
        
        total_score = structure_score + assertion_score + depth_score + maintainability_score
        
        if total_score >= 90:
            grade = 'EXCELLENT'
        elif total_score >= 80:
            grade = 'GOOD'
        elif total_score >= 70:
            grade = 'ACCEPTABLE'
        else:
            grade = 'FAIL'
        
        return {
            'structure_score': structure_score,
            'assertion_score': assertion_score,
            'depth_score': depth_score,
            'maintainability_score': maintainability_score,
            'total_score': total_score,
            'grade': grade,
            'score': total_score
        }
    
    def calculate_test_efficiency_grade(self, session):
        """Calculate test development efficiency grade"""
        
        duration_hours = session.get('duration_hours', 1)
        test_cases_written = session.get('test_cases_written', 0)
        test_lines_written = session.get('test_lines_written', 0)
        
        # Calculate efficiency metrics
        test_cases_per_hour = test_cases_written / duration_hours if duration_hours > 0 else 0
        lines_per_test_case = test_lines_written / test_cases_written if test_cases_written > 0 else 0
        lines_per_hour = test_lines_written / duration_hours if duration_hours > 0 else 0
        
        # Grade test case productivity (0-40 points)
        if test_cases_per_hour >= 8:  # 8+ test cases per hour
            case_productivity_score = 40
        elif test_cases_per_hour >= 6:
            case_productivity_score = 30 + (test_cases_per_hour - 6) * 5
        elif test_cases_per_hour >= 4:
            case_productivity_score = 20 + (test_cases_per_hour - 4) * 5
        else:
            case_productivity_score = max(0, test_cases_per_hour * 5)
        
        # Grade test case conciseness (0-30 points)
        if lines_per_test_case <= 15:  # Ideal: concise tests
            conciseness_score = 30
        elif lines_per_test_case <= 20:
            conciseness_score = 25 + (20 - lines_per_test_case)
        elif lines_per_test_case <= 25:
            conciseness_score = 20 + (25 - lines_per_test_case) * 0.8
        else:
            conciseness_score = max(0, 20 - (lines_per_test_case - 25) * 2)
        
        # Grade overall test writing speed (0-30 points)
        if lines_per_hour >= 120:  # Fast test development
            speed_score = 30
        elif lines_per_hour >= 80:
            speed_score = 20 + (lines_per_hour - 80) * 0.25
        elif lines_per_hour >= 50:
            speed_score = 10 + (lines_per_hour - 50) * 0.33
        else:
            speed_score = max(0, lines_per_hour * 0.2)
        
        total_score = case_productivity_score + conciseness_score + speed_score
        
        if total_score >= 90:
            grade = 'EXCELLENT'
        elif total_score >= 75:
            grade = 'GOOD'
        elif total_score >= 60:
            grade = 'ACCEPTABLE'
        else:
            grade = 'FAIL'
        
        return {
            'test_cases_per_hour': test_cases_per_hour,
            'lines_per_test_case': lines_per_test_case,
            'lines_per_hour': lines_per_hour,
            'case_productivity_score': case_productivity_score,
            'conciseness_score': conciseness_score,
            'speed_score': speed_score,
            'grade': grade,
            'score': total_score
        }
    
    def calculate_test_completeness_grade(self, session):
        """Calculate test scenario completeness grade"""
        
        planned_scenarios = session.get('planned_test_scenarios', [])
        completed_scenarios = session.get('completed_test_scenarios', [])
        edge_cases_covered = session.get('edge_cases_covered', 0)
        error_paths_tested = session.get('error_paths_tested', 0)
        
        # Scenario completion rate (0-40 points)
        if planned_scenarios:
            completion_rate = len(completed_scenarios) / len(planned_scenarios)
            completion_score = completion_rate * 40
        else:
            completion_score = 30  # Default if no explicit planning
        
        # Edge case coverage (0-30 points)
        expected_edge_cases = session.get('expected_edge_cases', 0)
        if expected_edge_cases > 0:
            edge_case_rate = edge_cases_covered / expected_edge_cases
            edge_case_score = min(30, edge_case_rate * 30)
        else:
            edge_case_score = 20 if edge_cases_covered > 0 else 10
        
        # Error path coverage (0-30 points)
        expected_error_paths = session.get('expected_error_paths', 0)
        if expected_error_paths > 0:
            error_path_rate = error_paths_tested / expected_error_paths
            error_path_score = min(30, error_path_rate * 30)
        else:
            error_path_score = 20 if error_paths_tested > 0 else 10
        
        total_score = completion_score + edge_case_score + error_path_score
        
        if total_score >= 90:
            grade = 'EXCELLENT'
        elif total_score >= 75:
            grade = 'GOOD'
        elif total_score >= 60:
            grade = 'ACCEPTABLE'
        else:
            grade = 'FAIL'
        
        return {
            'planned_scenarios': len(planned_scenarios),
            'completed_scenarios': len(completed_scenarios),
            'completion_rate': len(completed_scenarios) / len(planned_scenarios) if planned_scenarios else 0,
            'edge_cases_covered': edge_cases_covered,
            'error_paths_tested': error_paths_tested,
            'completion_score': completion_score,
            'edge_case_score': edge_case_score,
            'error_path_score': error_path_score,
            'grade': grade,
            'score': total_score
        }
    
    def calculate_test_size_management_grade(self, session):
        """Calculate grade for managing test code size impact"""
        
        test_lines_added = session.get('test_lines_written', 0)
        total_size_after = session.get('total_effort_size_after', 0)
        size_limit = session.get('size_limit', 800)
        implementation_lines = session.get('implementation_lines', 0)
        
        # Size compliance (0-50 points)
        size_utilization = (total_size_after / size_limit) * 100
        
        if total_size_after <= size_limit:
            if size_utilization <= 85:  # Well within limits
                size_compliance_score = 50
            elif size_utilization <= 95:  # Close to limits
                size_compliance_score = 40 + (95 - size_utilization)
            else:  # Very close to limits
                size_compliance_score = 30 + (100 - size_utilization) * 2
        else:
            size_compliance_score = 0  # Size violation
        
        # Test-to-implementation ratio (0-30 points)
        if implementation_lines > 0:
            test_ratio = (test_lines_added / implementation_lines) * 100
            
            if 15 <= test_ratio <= 25:  # Ideal ratio
                ratio_score = 30
            elif 10 <= test_ratio <= 35:  # Acceptable ratio
                ratio_score = 25
            elif 5 <= test_ratio <= 45:  # Suboptimal ratio
                ratio_score = 15
            else:  # Poor ratio
                ratio_score = 5
        else:
            ratio_score = 15  # Default when no implementation measured
        
        # Test efficiency (0-20 points)
        if test_lines_added > 0 and session.get('coverage_increase', 0) > 0:
            lines_per_coverage_percent = test_lines_added / session.get('coverage_increase', 1)
            
            if lines_per_coverage_percent <= 8:  # Very efficient
                efficiency_score = 20
            elif lines_per_coverage_percent <= 12:  # Good efficiency
                efficiency_score = 15 + (12 - lines_per_coverage_percent) * 1.25
            elif lines_per_coverage_percent <= 20:  # Acceptable efficiency
                efficiency_score = 10 + (20 - lines_per_coverage_percent) * 0.625
            else:  # Poor efficiency
                efficiency_score = max(0, 10 - (lines_per_coverage_percent - 20) * 0.5)
        else:
            efficiency_score = 10  # Default
        
        total_score = size_compliance_score + ratio_score + efficiency_score
        
        # Critical failure conditions
        if total_size_after > size_limit:
            grade = 'FAIL'
        elif total_score >= 85:
            grade = 'EXCELLENT'
        elif total_score >= 70:
            grade = 'GOOD'
        elif total_score >= 55:
            grade = 'ACCEPTABLE'
        else:
            grade = 'FAIL'
        
        return {
            'total_size_after': total_size_after,
            'size_limit': size_limit,
            'size_utilization': size_utilization,
            'test_lines_added': test_lines_added,
            'test_ratio': (test_lines_added / implementation_lines * 100) if implementation_lines > 0 else 0,
            'size_compliance_score': size_compliance_score,
            'ratio_score': ratio_score,
            'efficiency_score': efficiency_score,
            'grade': grade,
            'score': total_score
        }
    
    def calculate_test_maintainability_grade(self, session):
        """Calculate test code maintainability grade"""
        
        test_files = session.get('test_files_created', [])
        test_quality_metrics = session.get('test_quality_metrics', {})
        
        # Code organization (0-35 points)
        organization_score = self.assess_test_organization(test_files, test_quality_metrics)
        
        # Naming and clarity (0-25 points)
        clarity_score = self.assess_test_clarity(test_quality_metrics)
        
        # Test independence (0-25 points)
        independence_score = self.assess_test_independence(test_quality_metrics)
        
        # Documentation and comments (0-15 points)
        documentation_score = self.assess_test_documentation(test_quality_metrics)
        
        total_score = organization_score + clarity_score + independence_score + documentation_score
        
        if total_score >= 90:
            grade = 'EXCELLENT'
        elif total_score >= 80:
            grade = 'GOOD'
        elif total_score >= 65:
            grade = 'ACCEPTABLE'
        else:
            grade = 'FAIL'
        
        return {
            'organization_score': organization_score,
            'clarity_score': clarity_score,
            'independence_score': independence_score,
            'documentation_score': documentation_score,
            'grade': grade,
            'score': total_score
        }
    
    def calculate_overall_test_writing_grade(self, coverage, quality, efficiency, completeness, size_mgmt, maintainability):
        """Calculate weighted overall test writing grade"""
        
        # Weighted scoring:
        # Coverage Achievement: 50% (primary goal of test writing)
        # Test Quality: 20% (important for long-term maintainability)
        # Size Management: 15% (critical constraint)
        # Test Completeness: 8% (thoroughness)
        # Efficiency: 5% (development speed)
        # Maintainability: 2% (code quality)
        
        weighted_score = (
            coverage['score'] * 0.50 +
            quality['score'] * 0.20 +
            size_mgmt['score'] * 0.15 +
            completeness['score'] * 0.08 +
            efficiency['score'] * 0.05 +
            maintainability['score'] * 0.02
        )
        
        # Critical failure conditions
        critical_failures = []
        if size_mgmt['grade'] == 'FAIL':
            critical_failures.append('Size limit exceeded during test development')
        if coverage['score'] < 40:
            critical_failures.append('Insufficient coverage progress')
        if quality['score'] < 50:
            critical_failures.append('Test quality unacceptable')
        
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
    
    def assess_test_structure_quality(self, metrics):
        """Assess test structure and organization quality (0-25 points)"""
        
        score = 25
        
        # Check for table-driven tests
        table_driven_tests = metrics.get('table_driven_tests', 0)
        total_test_functions = metrics.get('total_test_functions', 1)
        table_driven_ratio = table_driven_tests / total_test_functions
        
        if table_driven_ratio >= 0.7:  # Most tests are table-driven
            score += 0  # Already at max
        elif table_driven_ratio >= 0.4:
            score -= 5
        else:
            score -= 10
        
        # Check for proper setup/teardown
        if not metrics.get('proper_setup_teardown', True):
            score -= 8
        
        # Check for test helper functions
        if metrics.get('helper_functions_count', 0) == 0:
            score -= 5  # No helpers might indicate code duplication
        
        # Check for excessive test function size
        avg_test_function_size = metrics.get('avg_test_function_lines', 20)
        if avg_test_function_size > 50:
            score -= 7
        elif avg_test_function_size > 30:
            score -= 3
        
        return max(0, score)
    
    def assess_assertion_quality(self, metrics):
        """Assess assertion quality and specificity (0-25 points)"""
        
        score = 25
        
        # Check assertion-to-test ratio
        total_assertions = metrics.get('total_assertions', 0)
        total_test_cases = metrics.get('total_test_cases', 1)
        assertion_ratio = total_assertions / total_test_cases
        
        if assertion_ratio >= 3:  # Good assertion coverage
            score += 0  # Already at max
        elif assertion_ratio >= 2:
            score -= 3
        elif assertion_ratio >= 1:
            score -= 7
        else:
            score -= 15  # Very few assertions
        
        # Check for specific vs. generic assertions
        specific_assertions = metrics.get('specific_assertions', 0)
        generic_assertions = metrics.get('generic_assertions', 0)
        
        if specific_assertions + generic_assertions > 0:
            specificity_ratio = specific_assertions / (specific_assertions + generic_assertions)
            if specificity_ratio < 0.6:
                score -= 5  # Too many generic assertions
        
        # Check for error assertion patterns
        error_assertions = metrics.get('error_assertions', 0)
        expected_error_tests = metrics.get('error_path_tests', 0)
        
        if expected_error_tests > 0 and error_assertions == 0:
            score -= 5  # Error tests without proper error assertions
        
        return max(0, score)
```

## Test Coverage Tracking

```python
class TestCoverageTracker:
    def __init__(self):
        self.coverage_history = []
        
    def track_coverage_progress(self, session_data):
        """Track test coverage progress over time"""
        
        coverage_point = {
            'timestamp': datetime.now().isoformat(),
            'session_id': session_data.get('session_id', ''),
            'coverage_percentage': session_data.get('coverage_after_percentage', 0),
            'coverage_increase': session_data.get('coverage_increase', 0),
            'test_lines_added': session_data.get('test_lines_written', 0),
            'test_cases_added': session_data.get('test_cases_written', 0),
            'session_duration_hours': session_data.get('duration_hours', 0)
        }
        
        # Calculate derived metrics
        coverage_point['coverage_velocity'] = (
            coverage_point['coverage_increase'] / coverage_point['session_duration_hours']
            if coverage_point['session_duration_hours'] > 0 else 0
        )
        
        coverage_point['coverage_efficiency'] = (
            coverage_point['coverage_increase'] / coverage_point['test_lines_added'] * 100
            if coverage_point['test_lines_added'] > 0 else 0
        )
        
        self.coverage_history.append(coverage_point)
        return coverage_point
    
    def analyze_coverage_trends(self, lookback_sessions=5):
        """Analyze coverage improvement trends"""
        
        if len(self.coverage_history) < 2:
            return {'trend': 'INSUFFICIENT_DATA', 'sessions_analyzed': len(self.coverage_history)}
        
        recent_sessions = self.coverage_history[-lookback_sessions:]
        
        # Calculate trend metrics
        velocities = [s['coverage_velocity'] for s in recent_sessions]
        efficiencies = [s['coverage_efficiency'] for s in recent_sessions if s['coverage_efficiency'] > 0]
        
        avg_velocity = sum(velocities) / len(velocities) if velocities else 0
        avg_efficiency = sum(efficiencies) / len(efficiencies) if efficiencies else 0
        
        # Analyze velocity trend
        if len(velocities) >= 3:
            velocity_trend = self.calculate_trend_direction(velocities)
        else:
            velocity_trend = 'STABLE'
        
        # Overall assessment
        if avg_velocity >= 10 and velocity_trend == 'IMPROVING':
            overall_trend = 'EXCELLENT'
        elif avg_velocity >= 8 and velocity_trend != 'DECLINING':
            overall_trend = 'GOOD'
        elif avg_velocity >= 5:
            overall_trend = 'ACCEPTABLE'
        else:
            overall_trend = 'POOR'
        
        return {
            'overall_trend': overall_trend,
            'avg_velocity_per_hour': avg_velocity,
            'avg_efficiency_per_line': avg_efficiency,
            'velocity_trend': velocity_trend,
            'sessions_analyzed': len(recent_sessions)
        }
    
    def calculate_trend_direction(self, values):
        """Calculate trend direction from a series of values"""
        
        if len(values) < 3:
            return 'STABLE'
        
        # Simple slope calculation
        x_values = list(range(len(values)))
        n = len(values)
        
        x_mean = sum(x_values) / n
        y_mean = sum(values) / n
        
        numerator = sum((x_values[i] - x_mean) * (values[i] - y_mean) for i in range(n))
        denominator = sum((x_values[i] - x_mean) ** 2 for i in range(n))
        
        if denominator == 0:
            return 'STABLE'
        
        slope = numerator / denominator
        
        if slope > 0.5:
            return 'IMPROVING'
        elif slope < -0.5:
            return 'DECLINING'
        else:
            return 'STABLE'
    
    def predict_coverage_completion(self, target_coverage):
        """Predict when target coverage will be reached"""
        
        if len(self.coverage_history) < 2:
            return {'prediction': 'INSUFFICIENT_DATA'}
        
        current_coverage = self.coverage_history[-1]['coverage_percentage']
        
        if current_coverage >= target_coverage:
            return {
                'prediction': 'TARGET_ACHIEVED',
                'current_coverage': current_coverage,
                'target_coverage': target_coverage
            }
        
        # Calculate recent velocity
        recent_sessions = self.coverage_history[-3:] if len(self.coverage_history) >= 3 else self.coverage_history
        avg_velocity = sum(s['coverage_velocity'] for s in recent_sessions) / len(recent_sessions)
        
        if avg_velocity <= 0:
            return {
                'prediction': 'TARGET_UNREACHABLE',
                'reason': 'No coverage progress detected'
            }
        
        coverage_gap = target_coverage - current_coverage
        estimated_hours = coverage_gap / avg_velocity
        
        return {
            'prediction': 'ACHIEVABLE',
            'current_coverage': current_coverage,
            'target_coverage': target_coverage,
            'coverage_gap': coverage_gap,
            'estimated_hours': estimated_hours,
            'avg_velocity': avg_velocity
        }
```

## Performance Dashboard

```python
def generate_test_writing_dashboard(session_data):
    """Generate real-time test writing performance dashboard"""
    
    grader = TestWritingGrader()
    tracker = TestCoverageTracker()
    
    # Grade current session
    current_grade = grader.grade_test_writing_session(session_data)
    
    # Get coverage trends
    coverage_trends = tracker.analyze_coverage_trends()
    
    # Get coverage prediction
    target_coverage = session_data.get('coverage_target_percentage', 85)
    coverage_prediction = tracker.predict_coverage_completion(target_coverage)
    
    dashboard = {
        'current_session': current_grade,
        'coverage_trends': coverage_trends,
        'coverage_prediction': coverage_prediction,
        'test_health': assess_test_writing_health(current_grade, session_data),
        'recommendations': generate_test_writing_recommendations(current_grade, coverage_trends)
    }
    
    print("📊 TEST WRITING PERFORMANCE DASHBOARD")
    print(f"Overall Grade: {current_grade['overall']['grade']} ({current_grade['overall']['weighted_score']:.1f}/100)")
    print(f"Coverage: {session_data.get('coverage_after_percentage', 0):.1f}% (target: {target_coverage}%)")
    print(f"Coverage Velocity: {current_grade['coverage_achievement']['coverage_rate_per_hour']:.1f}%/hour")
    print(f"Test Quality: {current_grade['test_quality']['grade']} ({current_grade['test_quality']['score']}/100)")
    print(f"Test Efficiency: {current_grade['test_efficiency']['lines_per_test_case']:.1f} lines/case")
    print(f"Size Impact: {session_data.get('total_effort_size_after', 0)}/800 lines")
    
    if coverage_prediction['prediction'] == 'ACHIEVABLE':
        print(f"Target ETA: {coverage_prediction['estimated_hours']:.1f} hours at current pace")
    elif coverage_prediction['prediction'] == 'TARGET_ACHIEVED':
        print("✅ Coverage target achieved!")
    
    if current_grade['overall']['critical_failures']:
        print("❌ CRITICAL ISSUES:")
        for failure in current_grade['overall']['critical_failures']:
            print(f"  - {failure}")
    
    return dashboard

def assess_test_writing_health(current_grade, session_data):
    """Assess overall test writing process health"""
    
    health_indicators = {
        'coverage_progress': current_grade['coverage_achievement']['grade'] != 'FAIL',
        'test_quality': current_grade['test_quality']['score'] >= 70,
        'size_compliance': current_grade['size_management']['grade'] != 'FAIL',
        'test_efficiency': current_grade['test_efficiency']['grade'] != 'FAIL'
    }
    
    if current_grade['overall']['has_critical_failures']:
        status = 'CRITICAL'
    elif all(health_indicators.values()):
        status = 'HEALTHY'
    elif sum(health_indicators.values()) >= 3:
        status = 'GOOD'
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

TEST WRITING WARNINGS
Coverage Velocity <5%/hour:
⚠️ WARNING: Slow coverage progress
⚠️ May indicate test complexity or approach issues

Size Limit Approached During Testing:
🚨 CRITICAL: Test code causing size violation
🚨 Optimize tests or split effort

Test Failure Rate >10%:
⚠️⚠️ WARNING: High test failure rate
⚠️⚠️ May indicate implementation issues

Test Quality Score <70:
⚠️ WARNING: Poor test quality detected
⚠️ Review test structure and assertions
---

## Performance State Tracking

```yaml
# Update orchestrator-state-v3.json
grading:
  SW_ENGINEER:
    TEST_WRITING:
      latest:
        timestamp: "2025-08-23T17:15:45Z"
        coverage_achieved: 87.3
        coverage_increase: 12.8
        test_quality_score: 89
        overall: "EXCELLENT"
        
      cumulative:
        test_sessions_completed: 6
        total_coverage_increase: 72.1
        avg_coverage_velocity: 9.2
        excellent: 2
        good: 3
        acceptable: 1
        fail: 0
