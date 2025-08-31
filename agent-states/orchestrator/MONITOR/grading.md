# Orchestrator - MONITOR State Grading

## Critical Performance Metrics

---
### ℹ️ RULE  - 
**Source:** 
**Criticality:** INFO - Best practice

PRIMARY METRIC: Agent Coordination Efficiency
Measurement: Response time to agent status changes
Target: <2 minutes to detect and act on status changes
Grade: PASS/FAIL (binary)
Weight: 40% of overall orchestrator grade
Consequence: FAIL = Loss of control over implementation
---

## Grading Rubric

| Metric | Excellent | Good | Acceptable | FAIL |
|--------|-----------|------|------------|------|
| Detection Time | <1 min | 1-2 min | 2-3 min | >3 min |
| Intervention Speed | <2 min | 2-5 min | 5-10 min | >10 min |
| Dependency Coordination | 100% accuracy | 95% accuracy | 90% accuracy | <90% |
| Progress Prediction | <10% variance | 10-20% variance | 20-30% variance | >30% |
| Resource Utilization | >90% optimal | 80-90% | 70-80% | <70% |

## Real-Time Scoring

```python
class MonitorStateGrader:
    def __init__(self):
        self.monitoring_sessions = []
        self.detection_times = []
        self.intervention_times = []
        
    def grade_monitoring_session(self, session_data):
        """Grade a complete monitoring session"""
        
        # Critical: Agent status detection
        detection_grade = self.calculate_detection_performance(session_data)
        
        # Intervention response time
        intervention_grade = self.calculate_intervention_performance(session_data)
        
        # Dependency coordination accuracy
        dependency_grade = self.evaluate_dependency_coordination(session_data)
        
        # Progress prediction accuracy
        prediction_grade = self.evaluate_progress_predictions(session_data)
        
        # Resource utilization optimization
        resource_grade = self.evaluate_resource_utilization(session_data)
        
        overall = self.calculate_overall_grade(
            detection_grade, intervention_grade, dependency_grade,
            prediction_grade, resource_grade
        )
        
        return {
            'detection': detection_grade,
            'intervention': intervention_grade,
            'dependency': dependency_grade,
            'prediction': prediction_grade,
            'resource': resource_grade,
            'overall': overall,
            'timestamp': datetime.now().isoformat()
        }
    
    def calculate_detection_performance(self, session):
        """Calculate agent status change detection performance"""
        
        status_changes = session.get('status_changes', [])
        detection_times = []
        
        for change in status_changes:
            actual_change_time = datetime.fromisoformat(change['occurred_at'])
            detection_time = datetime.fromisoformat(change['detected_at'])
            
            delay_seconds = (detection_time - actual_change_time).total_seconds()
            detection_times.append(delay_seconds)
        
        if not detection_times:
            return {
                'avg_detection_seconds': 0,
                'max_detection_seconds': 0,
                'grade': 'EXCELLENT',
                'score': 100
            }
        
        avg_delay = sum(detection_times) / len(detection_times)
        max_delay = max(detection_times)
        
        # Grade based on average detection time
        if avg_delay < 60:  # <1 minute
            grade = 'EXCELLENT'
            score = 100
        elif avg_delay < 120:  # <2 minutes
            grade = 'GOOD'
            score = 90
        elif avg_delay < 180:  # <3 minutes
            grade = 'PASS'
            score = 75
        else:
            grade = 'FAIL'
            score = 0
        
        # Penalty for any detection >5 minutes
        if max_delay > 300:
            grade = 'FAIL'
            score = 0
        
        return {
            'avg_detection_seconds': avg_delay,
            'max_detection_seconds': max_delay,
            'detection_count': len(detection_times),
            'grade': grade,
            'score': score
        }
    
    def calculate_intervention_performance(self, session):
        """Calculate intervention response time performance"""
        
        interventions = session.get('interventions', [])
        response_times = []
        
        for intervention in interventions:
            if 'triggered_at' in intervention and 'action_taken_at' in intervention:
                trigger_time = datetime.fromisoformat(intervention['triggered_at'])
                action_time = datetime.fromisoformat(intervention['action_taken_at'])
                
                response_seconds = (action_time - trigger_time).total_seconds()
                response_times.append({
                    'seconds': response_seconds,
                    'type': intervention['type'],
                    'urgency': intervention.get('urgency', 'MEDIUM')
                })
        
        if not response_times:
            return {
                'avg_response_seconds': 0,
                'interventions_count': 0,
                'grade': 'EXCELLENT',
                'score': 100
            }
        
        # Weight by urgency
        weighted_scores = []
        for resp in response_times:
            urgency_multiplier = {
                'IMMEDIATE': 1.0,  # No tolerance for delay
                'HIGH': 1.2,
                'MEDIUM': 1.5,
                'LOW': 2.0
            }.get(resp['urgency'], 1.0)
            
            target_seconds = {
                'IMMEDIATE': 120,  # 2 minutes
                'HIGH': 300,       # 5 minutes
                'MEDIUM': 600,     # 10 minutes
                'LOW': 1200        # 20 minutes
            }.get(resp['urgency'], 300)
            
            if resp['seconds'] <= target_seconds:
                weighted_scores.append(100)
            elif resp['seconds'] <= target_seconds * 1.5:
                weighted_scores.append(80)
            elif resp['seconds'] <= target_seconds * 2:
                weighted_scores.append(60)
            else:
                weighted_scores.append(0)
        
        avg_score = sum(weighted_scores) / len(weighted_scores)
        avg_response = sum(r['seconds'] for r in response_times) / len(response_times)
        
        if avg_score >= 95:
            grade = 'EXCELLENT'
        elif avg_score >= 85:
            grade = 'GOOD'
        elif avg_score >= 70:
            grade = 'PASS'
        else:
            grade = 'FAIL'
        
        return {
            'avg_response_seconds': avg_response,
            'interventions_count': len(response_times),
            'weighted_score': avg_score,
            'grade': grade,
            'score': avg_score
        }
    
    def evaluate_dependency_coordination(self, session):
        """Evaluate accuracy of dependency management"""
        
        dependency_events = session.get('dependency_events', [])
        coordination_accuracy = []
        
        for event in dependency_events:
            expected_action = event.get('expected_action')
            actual_action = event.get('actual_action')
            timing_accuracy = event.get('timing_accuracy', 100)
            
            if expected_action == actual_action:
                coordination_accuracy.append(timing_accuracy)
            else:
                coordination_accuracy.append(0)
        
        if not coordination_accuracy:
            return {
                'coordination_events': 0,
                'accuracy_percentage': 100,
                'grade': 'EXCELLENT',
                'score': 100
            }
        
        avg_accuracy = sum(coordination_accuracy) / len(coordination_accuracy)
        
        if avg_accuracy >= 95:
            grade = 'EXCELLENT'
        elif avg_accuracy >= 85:
            grade = 'GOOD'
        elif avg_accuracy >= 75:
            grade = 'PASS'
        else:
            grade = 'FAIL'
        
        return {
            'coordination_events': len(coordination_accuracy),
            'accuracy_percentage': avg_accuracy,
            'grade': grade,
            'score': avg_accuracy
        }
    
    def evaluate_progress_predictions(self, session):
        """Evaluate accuracy of progress and timeline predictions"""
        
        predictions = session.get('progress_predictions', [])
        prediction_accuracies = []
        
        for prediction in predictions:
            predicted_completion = prediction.get('predicted_completion_time')
            actual_completion = prediction.get('actual_completion_time')
            
            if predicted_completion and actual_completion:
                pred_time = datetime.fromisoformat(predicted_completion)
                actual_time = datetime.fromisoformat(actual_completion)
                
                difference_hours = abs((actual_time - pred_time).total_seconds() / 3600)
                original_estimate_hours = prediction.get('original_estimate_hours', 4)
                
                # Calculate accuracy as percentage
                accuracy = max(0, 100 - (difference_hours / original_estimate_hours * 100))
                prediction_accuracies.append(accuracy)
        
        if not prediction_accuracies:
            return {
                'predictions_count': 0,
                'avg_accuracy': 100,
                'grade': 'EXCELLENT',
                'score': 100
            }
        
        avg_accuracy = sum(prediction_accuracies) / len(prediction_accuracies)
        
        if avg_accuracy >= 90:
            grade = 'EXCELLENT'
        elif avg_accuracy >= 80:
            grade = 'GOOD'
        elif avg_accuracy >= 70:
            grade = 'PASS'
        else:
            grade = 'FAIL'
        
        return {
            'predictions_count': len(prediction_accuracies),
            'avg_accuracy': avg_accuracy,
            'variance_hours': self.calculate_prediction_variance(predictions),
            'grade': grade,
            'score': avg_accuracy
        }
    
    def evaluate_resource_utilization(self, session):
        """Evaluate how well resources (agents) were utilized"""
        
        agent_utilization = session.get('agent_utilization', {})
        
        utilization_scores = []
        for agent, data in agent_utilization.items():
            active_time = data.get('active_time_percentage', 0)
            parallel_efficiency = data.get('parallel_efficiency', 100)
            
            # Optimal utilization is high active time with good parallel coordination
            utilization_score = (active_time * 0.7) + (parallel_efficiency * 0.3)
            utilization_scores.append(utilization_score)
        
        if not utilization_scores:
            return {
                'agents_monitored': 0,
                'avg_utilization': 100,
                'grade': 'EXCELLENT',
                'score': 100
            }
        
        avg_utilization = sum(utilization_scores) / len(utilization_scores)
        
        if avg_utilization >= 90:
            grade = 'EXCELLENT'
        elif avg_utilization >= 80:
            grade = 'GOOD'
        elif avg_utilization >= 70:
            grade = 'PASS'
        else:
            grade = 'FAIL'
        
        return {
            'agents_monitored': len(utilization_scores),
            'avg_utilization': avg_utilization,
            'grade': grade,
            'score': avg_utilization
        }
    
    def calculate_overall_grade(self, detection, intervention, dependency, prediction, resource):
        """Calculate weighted overall grade"""
        
        # Detection: 40% (critical)
        # Intervention: 25%
        # Dependency: 20%
        # Prediction: 10%
        # Resource: 5%
        
        weighted_score = (
            detection['score'] * 0.40 +
            intervention['score'] * 0.25 +
            dependency['score'] * 0.20 +
            prediction['score'] * 0.10 +
            resource['score'] * 0.05
        )
        
        # Critical failure overrides
        if detection['grade'] == 'FAIL' or intervention['grade'] == 'FAIL':
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
            'critical_failure': detection['grade'] == 'FAIL' or intervention['grade'] == 'FAIL'
        }
```

## Monitoring Performance Tracking

```yaml
# Update orchestrator-state.yaml
grading:
  MONITOR:
    current_session:
      started_at: "2025-08-23T14:30:00Z"
      agents_monitored: 4
      status_changes_detected: 12
      interventions_triggered: 2
      avg_detection_time_seconds: 45
      avg_intervention_time_seconds: 180
      
    latest:
      timestamp: "2025-08-23T16:45:30Z"
      session_duration_hours: 2.25
      detection_grade: "EXCELLENT"
      intervention_grade: "GOOD"
      dependency_grade: "EXCELLENT"
      prediction_grade: "GOOD"
      resource_grade: "GOOD"
      overall: "GOOD"
    
    history:
      - {timestamp: "...", duration_hours: 1.8, grade: "EXCELLENT", agents: 3}
      - {timestamp: "...", duration_hours: 2.1, grade: "GOOD", agents: 4}
    
    cumulative:
      sessions: 15
      excellent: 8
      good: 5
      pass: 2
      fail: 0
      avg_detection_seconds: 52.3
      avg_intervention_seconds: 204.1
      avg_session_duration: 1.9
```

## Warning Triggers

---
### 🚨 RULE  - 
**Source:** 
**Criticality:** CRITICAL - Major impact on grading

MONITORING PERFORMANCE WARNINGS
Detection Time >2 minutes:
⚠️ WARNING: Slow agent status detection
⚠️ Average: {time}s (target: <120s)
⚠️ Review monitoring frequency and alerting

Intervention Delay >5 minutes:
⚠️⚠️ WARNING: Slow intervention response
⚠️⚠️ Critical issues may escalate during delays

Multiple Detection Failures:
❌ CRITICAL: Monitoring system unreliable
❌ Agent coordination compromised
❌ Manual oversight required
---

## Performance Optimization

```python
def optimize_monitoring_performance():
    """Guidelines for excellent monitoring grades"""
    
    optimization_strategies = {
        'detection_optimization': [
            'Implement real-time agent status webhooks',
            'Use multiple monitoring channels (status files, API, heartbeat)',
            'Set up automated alerting for status changes',
            'Cache recent status to detect changes quickly'
        ],
        
        'intervention_optimization': [
            'Pre-define intervention procedures for common issues',
            'Maintain ready-to-spawn agent contexts',
            'Use parallel intervention execution where possible',
            'Implement automatic retry for transient failures'
        ],
        
        'dependency_optimization': [
            'Model dependencies explicitly in state files',
            'Use event-driven dependency notifications',
            'Implement predictive dependency analysis',
            'Optimize agent start timing for maximum parallelization'
        ],
        
        'prediction_optimization': [
            'Use historical data to improve timeline estimates',
            'Factor in agent-specific performance patterns',
            'Account for dependency delays in predictions',
            'Update predictions based on real-time progress'
        ]
    }
    
    return optimization_strategies
```

## Real-Time Grade Calculation

```python
def calculate_monitoring_grade_realtime():
    """Calculate monitoring grade in real-time during session"""
    
    session_start = get_current_session_start()
    current_time = datetime.now()
    
    # Get current session metrics
    detection_metrics = get_detection_metrics_since(session_start)
    intervention_metrics = get_intervention_metrics_since(session_start)
    
    # Calculate running grades
    detection_grade = calculate_detection_grade(detection_metrics)
    intervention_grade = calculate_intervention_grade(intervention_metrics)
    
    # Project final grade based on current performance
    projected_grade = project_final_monitoring_grade(
        detection_grade, intervention_grade, current_time
    )
    
    print(f"📊 REAL-TIME MONITORING GRADE")
    print(f"Detection: {detection_grade['grade']} (avg: {detection_grade['avg_seconds']:.1f}s)")
    print(f"Intervention: {intervention_grade['grade']} (avg: {intervention_grade['avg_seconds']:.1f}s)")
    print(f"Projected Final: {projected_grade}")
    
    return {
        'detection': detection_grade,
        'intervention': intervention_grade,
        'projected': projected_grade,
        'session_duration': (current_time - session_start).total_seconds() / 3600
    }
```

## Dashboard Integration

```python
def generate_monitoring_grade_dashboard():
    """Generate monitoring performance dashboard"""
    
    current_grade = calculate_monitoring_grade_realtime()
    historical_trends = get_monitoring_grade_trends()
    
    dashboard = {
        'current_performance': current_grade,
        'trends': historical_trends,
        'recommendations': generate_performance_recommendations(current_grade),
        'alerts': check_performance_alerts(current_grade)
    }
    
    print("📊 MONITORING PERFORMANCE DASHBOARD")
    print(f"Current Grade: {current_grade['projected']}")
    print(f"Trend: {'📈' if historical_trends['improving'] else '📉'}")
    
    if dashboard['alerts']:
        print("⚠️ Performance Alerts:")
        for alert in dashboard['alerts']:
            print(f"  - {alert['message']}")
    
    return dashboard
