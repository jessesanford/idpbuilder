# SW Engineer - MEASURE_SIZE State Grading

## Critical Performance Metrics

---
### ℹ️ RULE  - 
**Source:** 
**Criticality:** INFO - Best practice

PRIMARY METRIC: Measurement Accuracy
Measurement: Correct tool usage and accurate results
Target: 100% compliance with line-counter.sh per R304
Grade: PASS/FAIL (binary)
Weight: 60% of overall measure size grade
Consequence: Wrong tool usage = FAIL
---

## Grading Rubric

| Metric | Excellent | Good | Acceptable | FAIL |
|--------|-----------|------|------------|------|
| Tool Compliance | 100% correct | 100% correct | 100% correct | Wrong tool used |
| Decision Accuracy | 100% correct | 90% correct | 80% correct | <80% |
| Analysis Depth | Comprehensive | Detailed | Basic | Insufficient |
| Speed | <5 minutes | 5-10 minutes | 10-15 minutes | >15 minutes |
| Documentation | Complete | Good | Adequate | Poor |

## Real-Time Scoring

```python
class MeasureSizeGrader:
    def __init__(self):
        self.measurement_sessions = []
        
    def grade_measure_size_session(self, session_data):
        """Grade a size measurement and analysis session"""
        
        # Critical: Tool compliance (mandatory line-counter.sh usage per R304)
        tool_compliance_grade = self.calculate_tool_compliance_grade(session_data)
        
        # Decision accuracy based on measurement results
        decision_accuracy_grade = self.calculate_decision_accuracy_grade(session_data)
        
        # Analysis comprehensiveness
        analysis_depth_grade = self.calculate_analysis_depth_grade(session_data)
        
        # Measurement completion speed
        speed_grade = self.calculate_speed_grade(session_data)
        
        # Documentation quality
        documentation_grade = self.calculate_documentation_grade(session_data)
        
        overall = self.calculate_overall_measure_size_grade(
            tool_compliance_grade, decision_accuracy_grade, analysis_depth_grade,
            speed_grade, documentation_grade
        )
        
        return {
            'tool_compliance': tool_compliance_grade,
            'decision_accuracy': decision_accuracy_grade,
            'analysis_depth': analysis_depth_grade,
            'speed': speed_grade,
            'documentation': documentation_grade,
            'overall': overall,
            'timestamp': datetime.now().isoformat()
        }
    
    def calculate_tool_compliance_grade(self, session):
        """Calculate tool compliance grade - CRITICAL for software factory"""
        
        tool_used = session.get('measurement_tool', '')
        command_used = session.get('measurement_command', '')
        
        # Check if correct tool was used
        correct_tool = 'line-counter.sh' in tool_used or 'line-counter.sh' in command_used
        
        if not correct_tool:
            # Check for wrong tools
            wrong_tools = ['wc -l', 'find', 'cloc', 'sloccount', 'tokei']
            tool_violation = any(wrong_tool in tool_used.lower() or wrong_tool in command_used.lower() 
                               for wrong_tool in wrong_tools)
            
            return {
                'correct_tool_used': False,
                'tool_violation_detected': tool_violation,
                'tool_used': tool_used,
                'grade': 'FAIL',
                'score': 0,
                'critical_failure': True,
                'failure_reason': f'Wrong tool used: {tool_used}' if tool_used else 'No measurement tool detected'
            }
        
        # Verify command format correctness
        command_format_score = 100
        
        if '-c' not in command_used:
            command_format_score -= 20  # Missing branch specification
        
        if '-d' not in command_used and session.get('detailed_analysis_performed', False):
            command_format_score -= 10  # Could have used detailed flag
        
        # Check for generated code exclusion awareness
        if session.get('generated_code_excluded', True):  # Should be default behavior
            exclusion_bonus = 5
        else:
            exclusion_bonus = -15  # Penalty for including generated code
        
        final_score = min(100, command_format_score + exclusion_bonus)
        
        if final_score >= 95:
            grade = 'EXCELLENT'
        elif final_score >= 85:
            grade = 'GOOD'
        elif final_score >= 70:
            grade = 'ACCEPTABLE'
        else:
            grade = 'FAIL'
        
        return {
            'correct_tool_used': True,
            'tool_used': tool_used,
            'command_format_score': command_format_score,
            'exclusion_bonus': exclusion_bonus,
            'grade': grade,
            'score': final_score,
            'critical_failure': grade == 'FAIL'
        }
    
    def calculate_decision_accuracy_grade(self, session):
        """Calculate accuracy of size-based state transition decision"""
        
        current_lines = session.get('measured_lines', 0)
        decision_made = session.get('next_state_decision', '')
        decision_reasoning = session.get('decision_reasoning', [])
        
        # Determine correct decision based on size thresholds
        if current_lines > 800:
            correct_decision = 'SPLIT_WORK'
            urgency_level = 'CRITICAL'
        elif current_lines > 750:
            # Complex decision - depends on completion percentage and projections
            completion_pct = session.get('implementation_completion_percentage', 50)
            projected_size = session.get('projected_completion_size', current_lines)
            
            if completion_pct > 90 and projected_size <= 800:
                correct_decision = 'IMPLEMENTATION'
            else:
                correct_decision = 'SPLIT_WORK'
            urgency_level = 'HIGH'
        elif current_lines > 700:
            correct_decision = 'IMPLEMENTATION'
            urgency_level = 'MEDIUM'
        else:
            correct_decision = 'IMPLEMENTATION'
            urgency_level = 'LOW'
        
        # Evaluate decision accuracy
        decision_correct = (decision_made == correct_decision)
        
        # Evaluate reasoning quality
        reasoning_quality = self.evaluate_decision_reasoning(
            decision_reasoning, current_lines, urgency_level
        )
        
        # Calculate score
        if decision_correct and reasoning_quality >= 80:
            grade = 'EXCELLENT'
            score = 100
        elif decision_correct and reasoning_quality >= 60:
            grade = 'GOOD'
            score = 90
        elif decision_correct:
            grade = 'ACCEPTABLE'
            score = 75
        elif self.is_reasonable_alternative_decision(decision_made, correct_decision, current_lines):
            grade = 'ACCEPTABLE'
            score = 70
        else:
            grade = 'FAIL'
            score = max(0, reasoning_quality / 2)
        
        return {
            'measured_lines': current_lines,
            'correct_decision': correct_decision,
            'actual_decision': decision_made,
            'decision_correct': decision_correct,
            'reasoning_quality': reasoning_quality,
            'urgency_level': urgency_level,
            'grade': grade,
            'score': score
        }
    
    def calculate_analysis_depth_grade(self, session):
        """Calculate comprehensiveness of size analysis performed"""
        
        analysis_components = {
            'basic_measurement': session.get('basic_measurement_performed', False),
            'detailed_breakdown': session.get('detailed_breakdown_performed', False),
            'trend_analysis': session.get('trend_analysis_performed', False),
            'optimization_analysis': session.get('optimization_opportunities_identified', False),
            'projection_analysis': session.get('completion_projection_performed', False),
            'test_ratio_analysis': session.get('test_ratio_calculated', False)
        }
        
        components_completed = sum(analysis_components.values())
        total_components = len(analysis_components)
        completion_percentage = (components_completed / total_components) * 100
        
        # Quality assessment of each completed component
        component_quality_scores = []
        
        if analysis_components['basic_measurement']:
            quality = self.assess_basic_measurement_quality(session)
            component_quality_scores.append(quality)
        
        if analysis_components['detailed_breakdown']:
            quality = self.assess_breakdown_analysis_quality(session)
            component_quality_scores.append(quality)
        
        if analysis_components['trend_analysis']:
            quality = self.assess_trend_analysis_quality(session)
            component_quality_scores.append(quality)
        
        # Calculate weighted score
        avg_component_quality = sum(component_quality_scores) / len(component_quality_scores) if component_quality_scores else 0
        depth_score = (completion_percentage * 0.6) + (avg_component_quality * 0.4)
        
        if depth_score >= 90:
            grade = 'EXCELLENT'
        elif depth_score >= 80:
            grade = 'GOOD'
        elif depth_score >= 65:
            grade = 'ACCEPTABLE'
        else:
            grade = 'FAIL'
        
        return {
            'components_completed': components_completed,
            'total_components': total_components,
            'completion_percentage': completion_percentage,
            'avg_component_quality': avg_component_quality,
            'component_details': analysis_components,
            'grade': grade,
            'score': depth_score
        }
    
    def calculate_speed_grade(self, session):
        """Calculate measurement and analysis completion speed"""
        
        start_time = session.get('analysis_start_time')
        end_time = session.get('analysis_end_time')
        
        if not start_time or not end_time:
            return {
                'duration_minutes': 0,
                'grade': 'ACCEPTABLE',
                'score': 75  # Default score when timing not available
            }
        
        start_dt = datetime.fromisoformat(start_time)
        end_dt = datetime.fromisoformat(end_time)
        duration_minutes = (end_dt - start_dt).total_seconds() / 60
        
        # Speed benchmarks based on analysis complexity
        analysis_complexity = session.get('analysis_complexity', 'MEDIUM')  # LOW, MEDIUM, HIGH
        
        if analysis_complexity == 'LOW':  # Basic measurement only
            target_time = 3
            good_time = 5
            acceptable_time = 8
        elif analysis_complexity == 'HIGH':  # Comprehensive analysis with trends
            target_time = 8
            good_time = 12
            acceptable_time = 18
        else:  # MEDIUM - standard analysis
            target_time = 5
            good_time = 8
            acceptable_time = 12
        
        if duration_minutes <= target_time:
            grade = 'EXCELLENT'
            score = 100
        elif duration_minutes <= good_time:
            grade = 'GOOD'
            score = 85 + (15 * (good_time - duration_minutes) / (good_time - target_time))
        elif duration_minutes <= acceptable_time:
            grade = 'ACCEPTABLE'
            score = 70 + (15 * (acceptable_time - duration_minutes) / (acceptable_time - good_time))
        else:
            grade = 'FAIL'
            score = max(0, 70 - (duration_minutes - acceptable_time) * 2)
        
        return {
            'duration_minutes': duration_minutes,
            'target_time': target_time,
            'analysis_complexity': analysis_complexity,
            'grade': grade,
            'score': score
        }
    
    def calculate_documentation_grade(self, session):
        """Calculate quality of size measurement documentation"""
        
        documentation_elements = {
            'measurement_command_logged': bool(session.get('measurement_command')),
            'results_documented': bool(session.get('measurement_results')),
            'decision_reasoning_recorded': bool(session.get('decision_reasoning')),
            'next_actions_specified': bool(session.get('next_actions')),
            'work_log_updated': session.get('work_log_updated', False),
            'size_trend_documented': session.get('size_trend_documented', False)
        }
        
        elements_completed = sum(documentation_elements.values())
        total_elements = len(documentation_elements)
        completeness_percentage = (elements_completed / total_elements) * 100
        
        # Quality assessment of documentation content
        documentation_quality = 100
        
        # Check command documentation quality
        if documentation_elements['measurement_command_logged']:
            command = session.get('measurement_command', '')
            if not command or 'line-counter.sh' not in command:
                documentation_quality -= 20
        
        # Check results documentation quality
        if documentation_elements['results_documented']:
            results = session.get('measurement_results', '')
            if not results or not any(char.isdigit() for char in results):
                documentation_quality -= 15
        
        # Check reasoning quality
        if documentation_elements['decision_reasoning_recorded']:
            reasoning = session.get('decision_reasoning', [])
            if not reasoning or len(reasoning) < 2:
                documentation_quality -= 15
        
        # Calculate final score
        final_score = (completeness_percentage * 0.7) + (documentation_quality * 0.3)
        
        if final_score >= 90:
            grade = 'EXCELLENT'
        elif final_score >= 80:
            grade = 'GOOD'
        elif final_score >= 65:
            grade = 'ACCEPTABLE'
        else:
            grade = 'FAIL'
        
        return {
            'elements_completed': elements_completed,
            'total_elements': total_elements,
            'completeness_percentage': completeness_percentage,
            'documentation_quality': documentation_quality,
            'element_details': documentation_elements,
            'grade': grade,
            'score': final_score
        }
    
    def calculate_overall_measure_size_grade(self, tool_compliance, decision_accuracy, analysis_depth, speed, documentation):
        """Calculate weighted overall measure size grade"""
        
        # Weighted scoring:
        # Tool Compliance: 60% (critical - must use correct tool)
        # Decision Accuracy: 20% (important for state transitions)
        # Analysis Depth: 12% (comprehensive analysis valuable)
        # Documentation: 5% (important for traceability)
        # Speed: 3% (efficiency matters but not critical)
        
        weighted_score = (
            tool_compliance['score'] * 0.60 +
            decision_accuracy['score'] * 0.20 +
            analysis_depth['score'] * 0.12 +
            documentation['score'] * 0.05 +
            speed['score'] * 0.03
        )
        
        # Critical failure conditions override everything
        critical_failures = []
        if tool_compliance.get('critical_failure', False):
            critical_failures.append('Wrong measurement tool used')
        if decision_accuracy['score'] < 50:
            critical_failures.append('Poor decision accuracy')
        
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
    
    def evaluate_decision_reasoning(self, reasoning_list, current_lines, urgency_level):
        """Evaluate quality of decision reasoning"""
        
        if not reasoning_list:
            return 0
        
        quality_score = 0
        max_score = 100
        
        # Check for size-specific reasoning
        size_mentioned = any('lines' in reason.lower() or 'size' in reason.lower() 
                           for reason in reasoning_list)
        if size_mentioned:
            quality_score += 25
        
        # Check for threshold awareness
        threshold_mentioned = any(str(threshold) in reason 
                                for reason in reasoning_list
                                for threshold in [700, 750, 800])
        if threshold_mentioned:
            quality_score += 25
        
        # Check for urgency awareness
        urgency_mentioned = any(urgency_level.lower() in reason.lower() 
                              for reason in reasoning_list)
        if urgency_mentioned:
            quality_score += 20
        
        # Check for projection consideration
        projection_mentioned = any('project' in reason.lower() or 'complet' in reason.lower()
                                 for reason in reasoning_list)
        if projection_mentioned:
            quality_score += 15
        
        # Check for implementation plan awareness
        plan_mentioned = any('plan' in reason.lower() or 'implement' in reason.lower()
                           for reason in reasoning_list)
        if plan_mentioned:
            quality_score += 10
        
        # Length and detail bonus (up to 5 points)
        total_reasoning_length = sum(len(reason) for reason in reasoning_list)
        if total_reasoning_length > 200:
            quality_score += 5
        elif total_reasoning_length > 100:
            quality_score += 3
        elif total_reasoning_length > 50:
            quality_score += 1
        
        return min(max_score, quality_score)
    
    def is_reasonable_alternative_decision(self, actual_decision, correct_decision, current_lines):
        """Check if actual decision is a reasonable alternative to the correct one"""
        
        # Allow some flexibility in the 700-750 range
        if 700 < current_lines <= 750:
            reasonable_alternatives = {
                'IMPLEMENTATION': ['SPLIT_WORK', 'FIX_ISSUES'],  # Could optimize first
                'SPLIT_WORK': ['IMPLEMENTATION'],  # Could try to finish
                'FIX_ISSUES': ['IMPLEMENTATION', 'SPLIT_WORK']  # Could optimize then decide
            }
            return actual_decision in reasonable_alternatives.get(correct_decision, [])
        
        return False
    
    def assess_basic_measurement_quality(self, session):
        """Assess quality of basic measurement execution"""
        
        quality_score = 100
        
        # Check if tool path is correct
        tool_used = session.get('measurement_tool', '')
        if 'line-counter.sh' not in tool_used:
            quality_score -= 20
        
        # Check if branch was specified
        command = session.get('measurement_command', '')
        if '-c' not in command:
            quality_score -= 15
        
        # Check if results are numeric and reasonable
        try:
            lines = int(session.get('measured_lines', 0))
            if lines <= 0 or lines > 10000:  # Sanity check
                quality_score -= 25
        except (ValueError, TypeError):
            quality_score -= 30
        
        return max(0, quality_score)
    
    def assess_breakdown_analysis_quality(self, session):
        """Assess quality of detailed breakdown analysis"""
        
        breakdown_data = session.get('size_breakdown', {})
        
        if not breakdown_data:
            return 0
        
        quality_score = 100
        
        # Check for directory-level breakdown
        if not breakdown_data.get('breakdown_by_directory'):
            quality_score -= 30
        
        # Check for file type analysis
        if not breakdown_data.get('breakdown_by_file_type'):
            quality_score -= 20
        
        # Check for largest contributors identification
        if not breakdown_data.get('largest_contributors'):
            quality_score -= 20
        
        # Check for optimization opportunities
        if not breakdown_data.get('optimization_opportunities'):
            quality_score -= 15
        
        # Check for test ratio calculation
        if not breakdown_data.get('test_to_implementation_ratio'):
            quality_score -= 15
        
        return max(0, quality_score)
    
    def assess_trend_analysis_quality(self, session):
        """Assess quality of growth trend analysis"""
        
        trend_data = session.get('trend_analysis', {})
        
        if not trend_data:
            return 0
        
        quality_score = 100
        
        # Check for historical data usage
        measurements_count = trend_data.get('measurements_count', 0)
        if measurements_count < 3:
            quality_score -= 25
        
        # Check for growth rate calculation
        if not trend_data.get('average_growth_rate'):
            quality_score -= 20
        
        # Check for trend direction assessment
        if not trend_data.get('growth_trend'):
            quality_score -= 20
        
        # Check for completion projection
        if not trend_data.get('projected_completion_size'):
            quality_score -= 20
        
        # Bonus for identifying concerning trends
        if trend_data.get('growth_trend') in ['ACCELERATING'] and session.get('trend_warnings_noted', False):
            quality_score += 10
        
        return max(0, min(100, quality_score))
```

## Size Measurement Performance Tracking

```python
class SizeMeasurementTracker:
    def __init__(self):
        self.measurement_history = []
        
    def track_measurement_session(self, session_data):
        """Track a complete size measurement session"""
        
        session_metrics = {
            'session_id': f"measure_{datetime.now().strftime('%Y%m%d_%H%M%S')}",
            'timestamp': datetime.now().isoformat(),
            
            # Measurement data
            'measured_lines': session_data.get('measured_lines', 0),
            'measurement_tool': session_data.get('measurement_tool', ''),
            'tool_compliant': 'line-counter.sh' in session_data.get('measurement_tool', ''),
            
            # Analysis depth
            'analysis_components': session_data.get('analysis_components_completed', 0),
            'detailed_breakdown': session_data.get('detailed_breakdown_performed', False),
            'trend_analysis': session_data.get('trend_analysis_performed', False),
            'optimization_analysis': session_data.get('optimization_opportunities_identified', False),
            
            # Decision quality
            'decision_made': session_data.get('next_state_decision', ''),
            'decision_reasoning_count': len(session_data.get('decision_reasoning', [])),
            'immediate_actions_count': len(session_data.get('immediate_actions', [])),
            
            # Performance metrics
            'analysis_duration_minutes': session_data.get('analysis_duration_minutes', 0),
            'documentation_completeness': session_data.get('documentation_completeness_percentage', 0)
        }
        
        # Calculate derived metrics
        session_metrics['analysis_efficiency'] = self.calculate_analysis_efficiency(session_metrics)
        session_metrics['decision_comprehensiveness'] = self.calculate_decision_comprehensiveness(session_metrics)
        
        self.measurement_history.append(session_metrics)
        return session_metrics
    
    def calculate_analysis_efficiency(self, metrics):
        """Calculate analysis efficiency score"""
        
        components_completed = metrics['analysis_components']
        duration_minutes = metrics['analysis_duration_minutes']
        
        if duration_minutes <= 0:
            return 100  # Assume efficient if no timing data
        
        # Target: 1 component per 2 minutes for efficiency
        expected_duration = components_completed * 2
        
        if duration_minutes <= expected_duration:
            return 100
        else:
            efficiency = (expected_duration / duration_minutes) * 100
            return max(0, min(100, efficiency))
    
    def calculate_decision_comprehensiveness(self, metrics):
        """Calculate decision comprehensiveness score"""
        
        base_score = 0
        
        # Decision made (20 points)
        if metrics['decision_made']:
            base_score += 20
        
        # Reasoning provided (40 points max)
        reasoning_score = min(40, metrics['decision_reasoning_count'] * 10)
        base_score += reasoning_score
        
        # Immediate actions specified (30 points max)
        actions_score = min(30, metrics['immediate_actions_count'] * 7.5)
        base_score += actions_score
        
        # Analysis depth bonus (10 points max)
        if metrics['detailed_breakdown']:
            base_score += 5
        if metrics['trend_analysis']:
            base_score += 3
        if metrics['optimization_analysis']:
            base_score += 2
        
        return min(100, base_score)
    
    def get_measurement_performance_trends(self, lookback_sessions=10):
        """Analyze measurement performance trends"""
        
        if len(self.measurement_history) < 2:
            return {'trend': 'INSUFFICIENT_DATA', 'sessions_analyzed': len(self.measurement_history)}
        
        recent_sessions = self.measurement_history[-lookback_sessions:]
        
        # Calculate trend metrics
        tool_compliance_rate = sum(1 for s in recent_sessions if s['tool_compliant']) / len(recent_sessions) * 100
        avg_analysis_efficiency = sum(s['analysis_efficiency'] for s in recent_sessions) / len(recent_sessions)
        avg_decision_comprehensiveness = sum(s['decision_comprehensiveness'] for s in recent_sessions) / len(recent_sessions)
        avg_duration = sum(s['analysis_duration_minutes'] for s in recent_sessions) / len(recent_sessions)
        
        # Determine overall trend
        if tool_compliance_rate == 100 and avg_analysis_efficiency > 85 and avg_decision_comprehensiveness > 80:
            overall_trend = 'EXCELLENT'
        elif tool_compliance_rate >= 90 and avg_analysis_efficiency > 70 and avg_decision_comprehensiveness > 70:
            overall_trend = 'GOOD'
        elif tool_compliance_rate >= 80 and avg_analysis_efficiency > 60:
            overall_trend = 'ACCEPTABLE'
        else:
            overall_trend = 'NEEDS_IMPROVEMENT'
        
        return {
            'overall_trend': overall_trend,
            'tool_compliance_rate': tool_compliance_rate,
            'avg_analysis_efficiency': avg_analysis_efficiency,
            'avg_decision_comprehensiveness': avg_decision_comprehensiveness,
            'avg_duration_minutes': avg_duration,
            'sessions_analyzed': len(recent_sessions)
        }
```

## Performance Dashboard

```python
def generate_measure_size_dashboard(session_data):
    """Generate real-time size measurement performance dashboard"""
    
    grader = MeasureSizeGrader()
    tracker = SizeMeasurementTracker()
    
    # Grade current session
    current_grade = grader.grade_measure_size_session(session_data)
    
    # Get performance trends
    performance_trends = tracker.get_measurement_performance_trends()
    
    dashboard = {
        'current_session': current_grade,
        'performance_trends': performance_trends,
        'size_status': assess_size_measurement_health(current_grade, session_data),
        'recommendations': generate_measurement_recommendations(current_grade, session_data)
    }
    
    print("📊 SIZE MEASUREMENT PERFORMANCE DASHBOARD")
    print(f"Overall Grade: {current_grade['overall']['grade']} ({current_grade['overall']['weighted_score']:.1f}/100)")
    print(f"Tool Compliance: {'✅' if current_grade['tool_compliance']['correct_tool_used'] else '❌'} {current_grade['tool_compliance']['grade']}")
    print(f"Decision Accuracy: {current_grade['decision_accuracy']['grade']} ({current_grade['decision_accuracy']['score']}/100)")
    print(f"Analysis Depth: {current_grade['analysis_depth']['grade']} ({current_grade['analysis_depth']['completion_percentage']:.1f}% complete)")
    print(f"Speed: {current_grade['speed']['duration_minutes']:.1f} minutes ({current_grade['speed']['grade']})")
    
    # Size status
    measured_lines = session_data.get('measured_lines', 0)
    utilization = (measured_lines / 800) * 100
    print(f"Size Status: {measured_lines}/800 lines ({utilization:.1f}%)")
    
    if current_grade['overall']['critical_failures']:
        print("❌ CRITICAL FAILURES:")
        for failure in current_grade['overall']['critical_failures']:
            print(f"  - {failure}")
    
    return dashboard

def assess_size_measurement_health(current_grade, session_data):
    """Assess overall size measurement health"""
    
    health_indicators = {
        'tool_compliance': current_grade['tool_compliance']['correct_tool_used'],
        'measurement_accuracy': current_grade['overall']['grade'] != 'FAIL',
        'size_status': session_data.get('measured_lines', 0) <= 800,
        'decision_quality': current_grade['decision_accuracy']['decision_correct']
    }
    
    if not health_indicators['tool_compliance']:
        status = 'CRITICAL'
    elif not health_indicators['size_status']:
        status = 'SIZE_VIOLATION'
    elif all(health_indicators.values()):
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

SIZE MEASUREMENT WARNINGS
Wrong Tool Used:
🚨 CRITICAL: Non-compliant measurement tool
🚨 Must use only line-counter.sh per R304

Size Limit Exceeded:
🚨 CRITICAL: >800 lines detected
🚨 Immediate split required

Analysis Duration >15 minutes:
⚠️ WARNING: Measurement taking too long
⚠️ May indicate analysis paralysis

Poor Decision Accuracy:
⚠️⚠️ WARNING: Incorrect state transition decisions
⚠️⚠️ Review decision criteria
---

## Performance State Tracking

```yaml
# Update orchestrator-state.yaml
grading:
  SW_ENGINEER:
    MEASURE_SIZE:
      latest:
        timestamp: "2025-08-23T16:45:30Z"
        measured_lines: 742
        tool_compliance: "PASS"
        decision_accuracy: 95
        analysis_depth: 88
        overall: "GOOD"
        
      cumulative:
        measurements_performed: 12
        tool_compliance_rate: 100
        avg_decision_accuracy: 91.2
        excellent: 4
        good: 6
        acceptable: 2
        fail: 0
