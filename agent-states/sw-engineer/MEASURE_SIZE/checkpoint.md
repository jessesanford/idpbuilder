# SW Engineer - MEASURE_SIZE State Checkpoint

## When to Save State

Save checkpoint at these critical measurement milestones:

1. **Pre-Measurement Planning**
   - Size measurement triggered by threshold or warning
   - Measurement approach planned
   - Analysis scope determined

2. **Measurement Execution**
   - Primary size measurement completed
   - Tool compliance verified
   - Initial results obtained

3. **Analysis Completion**
   - Detailed breakdown analysis finished
   - Growth trend analysis completed
   - Optimization opportunities identified

4. **Decision Point**
   - Size-based transition decision made
   - Decision reasoning documented
   - Next actions planned

## Required Data to Preserve

```yaml
measure_size_checkpoint:
  # State identification
  state: "MEASURE_SIZE"
  effort_id: "effort2-controller"
  branch: "phase1/wave2/effort2-controller"
  working_dir: "/workspaces/efforts/phase1/wave2/effort2-controller"
  checkpoint_timestamp: "2025-08-23T16:45:30Z"
  
  # Measurement trigger context
  measurement_trigger:
    trigger_source: "IMPLEMENTATION"  # State that triggered measurement
    trigger_reason: "THRESHOLD_WARNING"  # THRESHOLD_WARNING, SCHEDULED, MANUAL
    trigger_threshold: 700  # Line threshold that triggered measurement
    previous_measurement:
      lines: 687
      timestamp: "2025-08-23T15:30:00Z"
      utilization_percentage: 85.9
    
  # Measurement execution
  measurement_execution:
    started_at: "2025-08-23T16:40:00Z"
    completed_at: "2025-08-23T16:45:00Z"
    duration_minutes: 5.0
    
    primary_measurement:
      tool_used: "/workspaces/kcp-shared-tools/tmc-pr-line-counter.sh"
      command_executed: "tmc-pr-line-counter.sh -c phase1/wave2/effort2-controller"
      command_output: "742 lines in effort (excluding generated code)"
      measurement_result: 742
      tool_compliant: true
      
    measurement_validation:
      repeat_measurements_performed: 3
      measurements: [742, 742, 741]
      variance: 1
      final_measurement: 742
      measurement_consistent: true
      
  # Size analysis results
  size_analysis:
    current_size_status:
      total_lines: 742
      size_limit: 800
      utilization_percentage: 92.75
      status: "DANGER"  # COMPLIANT, WARNING, DANGER, VIOLATION
      remaining_capacity: 58
      
    detailed_breakdown:
      tool_used: "/workspaces/kcp-shared-tools/tmc-pr-line-counter.sh -d"
      breakdown_by_directory:
        - path: "pkg/controllers"
          lines: 345
          percentage: 46.5
          
        - path: "pkg/api/v1"
          lines: 189
          percentage: 25.5
          
        - path: "pkg/webhooks"
          lines: 134
          percentage: 18.1
          
        - path: "test"
          lines: 74
          percentage: 10.0
          
      breakdown_by_file_type:
        - type: "*.go"
          lines: 668
          count: 8
          percentage: 90.0
          
        - type: "*_test.go"
          lines: 74
          count: 4
          percentage: 10.0
          
      largest_contributors:
        - path: "pkg/controllers/resource_controller.go"
          lines: 234
          
        - path: "pkg/api/v1/resource_types.go"
          lines: 145
          
        - path: "pkg/webhooks/admission.go"
          lines: 98
          
        - path: "pkg/controllers/status_manager.go"
          lines: 87
          
    optimization_opportunities:
      - type: "LARGE_FILE"
        description: "resource_controller.go at 234 lines could be split"
        estimated_savings: 50
        priority: "HIGH"
        
      - type: "SIMILAR_PATTERNS"
        description: "Status management code duplicated across files"
        estimated_savings: 35
        priority: "MEDIUM"
        
    test_coverage_analysis:
      test_lines: 74
      implementation_lines: 668
      test_ratio: 11.1  # percentage
      test_ratio_status: "LOW"  # Target should be 15-25%
      
  # Growth trend analysis
  growth_trend_analysis:
    historical_measurements:
      - timestamp: "2025-08-23T14:00:00Z"
        lines: 387
        source: "work_log"
        
      - timestamp: "2025-08-23T15:30:00Z"
        lines: 687
        source: "checkpoint"
        
      - timestamp: "2025-08-23T16:45:00Z"
        lines: 742
        source: "current_measurement"
        
    growth_metrics:
      timespan_hours: 2.75
      total_growth_lines: 355
      average_growth_rate: 129.1  # lines per hour
      recent_growth_rate: 73.3   # last 1.25 hours
      growth_trend: "DECELERATING"
      
    completion_projection:
      implementation_plan_file: "/workspaces/efforts/phase1/wave2/effort2-controller/IMPLEMENTATION-PLAN.md"
      current_completion_percentage: 78
      remaining_work_percentage: 22
      estimated_remaining_hours: 1.5
      
      projected_additional_lines: 110  # based on recent growth rate
      projected_final_size: 852
      projection_exceeds_limit: true
      
  # Decision analysis
  decision_analysis:
    analysis_started_at: "2025-08-23T16:45:30Z"
    analysis_completed_at: "2025-08-23T16:48:00Z"
    
    size_thresholds:
      current_lines: 742
      threshold_700: true   # exceeded
      threshold_750: false  # not exceeded
      threshold_800: false  # not exceeded
      violation: false      # not exceeded
      
    decision_factors:
      size_status: "DANGER"
      completion_percentage: 78
      projected_final_size: 852
      projection_exceeds_limit: true
      optimization_opportunities_exist: true
      recent_growth_rate_manageable: true
      
    decision_matrix_evaluation:
      continue_implementation:
        viable: false
        reason: "Projected to exceed size limit"
        risk_level: "HIGH"
        
      optimize_first:
        viable: true
        reason: "Optimization could save ~85 lines"
        estimated_savings: 85
        post_optimization_projected_size: 767
        risk_level: "MEDIUM"
        
      split_effort:
        viable: true
        reason: "Guaranteed to stay within limits"
        split_complexity: "MEDIUM"
        risk_level: "LOW"
        
    final_decision:
      next_state: "FIX_ISSUES"  # Optimize before continuing
      primary_reason: "Optimization opportunities can prevent size violation"
      confidence: 85
      fallback_plan: "SPLIT_WORK if optimization insufficient"
      
      decision_reasoning:
        - "Current size 742 lines in danger zone (92.75% of limit)"
        - "Projected final size 852 lines exceeds 800-line limit"
        - "Optimization opportunities identified with ~85 line savings potential"
        - "Post-optimization projected size ~767 lines within acceptable range"
        - "Risk mitigation: Split effort if optimization proves insufficient"
        
      immediate_actions:
        - "Transition to FIX_ISSUES state for optimization"
        - "Focus on resource_controller.go refactoring (234 lines → ~180 lines)"
        - "Extract common status management patterns"
        - "Re-measure size after optimization"
        - "Prepare SPLIT_WORK plan as contingency"
        
  # Documentation and reporting
  measurement_documentation:
    measurement_report_generated: true
    report_location: "/workspaces/efforts/phase1/wave2/effort2-controller/SIZE-MEASUREMENT-REPORT.md"
    report_timestamp: "2025-08-23T16:48:00Z"
    
    work_log_entry:
      entry_added: true
      entry_timestamp: "2025-08-23T16:48:00Z"
      entry_content: |
        ## [2025-08-23 16:45] SIZE MEASUREMENT - DANGER ZONE
        **Tool**: tmc-pr-line-counter.sh -c phase1/wave2/effort2-controller
        **Result**: 742/800 lines (92.75% utilization)
        **Status**: DANGER - Optimization required
        
        ### Analysis Results
        - Projected final size: 852 lines (EXCEEDS LIMIT)
        - Optimization potential: ~85 lines savings
        - Decision: FIX_ISSUES (optimize first)
        
        ### Next Actions
        - Refactor resource_controller.go (234 → ~180 lines)
        - Extract status management patterns
        - Re-measure after optimization
        
    size_history_updated: true
    
  # Risk assessment
  risk_assessment:
    current_risk_level: "HIGH"
    
    identified_risks:
      - risk: "SIZE_LIMIT_VIOLATION"
        probability: "HIGH"
        impact: "CRITICAL"
        mitigation: "Optimization and potential split"
        
      - risk: "OPTIMIZATION_INSUFFICIENT"
        probability: "MEDIUM"
        impact: "HIGH"
        mitigation: "Split effort if optimization saves <60 lines"
        
      - risk: "IMPLEMENTATION_DELAY"
        probability: "LOW"
        impact: "MEDIUM"
        mitigation: "Optimization should be quick (1-2 hours max)"
        
    contingency_plans:
      - trigger: "Optimization saves <60 lines"
        action: "Immediately transition to SPLIT_WORK"
        preparation_required: "Have split plan ready"
        
      - trigger: "Optimization takes >2 hours"
        action: "Abort optimization, proceed with split"
        reason: "Time efficiency vs. guaranteed compliance"
        
  # Performance metrics
  measurement_performance:
    analysis_efficiency: 95  # 5 minutes for comprehensive analysis
    tool_compliance: 100     # Used correct tool
    decision_quality: 90     # Sound reasoning and contingencies
    documentation_completeness: 100
    
  # Next session planning
  next_session:
    planned_state: "FIX_ISSUES"
    primary_focus: "Code optimization to prevent size violation"
    estimated_duration_hours: 1.5
    
    optimization_plan:
      - target: "pkg/controllers/resource_controller.go"
        current_lines: 234
        target_lines: 180
        approach: "Extract helper methods, split reconcile logic"
        estimated_savings: 54
        
      - target: "Status management patterns"
        current_lines: 87  # distributed across files
        target_lines: 52
        approach: "Create shared status manager utility"
        estimated_savings: 35
        
    success_criteria:
      - "Post-optimization size ≤767 lines"
      - "All tests still passing"
      - "No functionality reduction"
      - "Code quality maintained or improved"
      
    fallback_criteria:
      - "If savings <60 lines: immediate split"
      - "If optimization takes >2 hours: proceed to split"
      - "If tests break: revert and split"
```

## Recovery Protocol

### Context Recovery After Interruption

```python
def recover_measure_size_state(checkpoint_data):
    """Recover size measurement state from checkpoint"""
    
    print("🔄 RECOVERING MEASURE_SIZE STATE")
    
    effort_info = checkpoint_data.get('effort_id', 'unknown')
    measurement = checkpoint_data.get('measurement_execution', {})
    decision = checkpoint_data.get('decision_analysis', {})
    
    print(f"Effort: {effort_info}")
    print(f"Measurement: {measurement.get('measurement_result', 0)}/800 lines")
    print(f"Status: {checkpoint_data.get('size_analysis', {}).get('current_size_status', {}).get('status', 'UNKNOWN')}")
    print(f"Decision: {decision.get('final_decision', {}).get('next_state', 'UNDECIDED')}")
    
    # Verify measurement is still valid
    measurement_verification = verify_measurement_still_valid(checkpoint_data)
    
    # Check for changes since measurement
    changes_detected = detect_changes_since_measurement(checkpoint_data)
    
    # Determine recovery actions
    recovery_actions = determine_measure_size_recovery_actions(
        checkpoint_data, measurement_verification, changes_detected
    )
    
    return {
        'effort_id': effort_info,
        'measured_lines': measurement.get('measurement_result', 0),
        'measurement_status': checkpoint_data.get('size_analysis', {}).get('current_size_status', {}).get('status', 'UNKNOWN'),
        'decision_made': decision.get('final_decision', {}).get('next_state'),
        'measurement_verification': measurement_verification,
        'changes_since_measurement': changes_detected,
        'recovery_actions': recovery_actions,
        'recovery_needed': len(recovery_actions) > 0
    }

def verify_measurement_still_valid(checkpoint_data):
    """Verify that the size measurement is still accurate"""
    
    measurement_time = datetime.fromisoformat(checkpoint_data['measurement_execution']['completed_at'])
    current_time = datetime.now()
    time_elapsed = (current_time - measurement_time).total_seconds() / 60  # minutes
    
    verification_results = {
        'still_valid': True,
        'time_elapsed_minutes': time_elapsed,
        'issues_detected': []
    }
    
    # If measurement is very recent (<30 minutes), likely still valid
    if time_elapsed < 30:
        return verification_results
    
    # For older measurements, need to re-verify
    branch = checkpoint_data.get('branch', '')
    working_dir = checkpoint_data.get('working_dir', '')
    
    if not os.path.exists(working_dir):
        verification_results['still_valid'] = False
        verification_results['issues_detected'].append('Working directory no longer exists')
        return verification_results
    
    try:
        # Re-measure current size
        result = subprocess.run([
            '/workspaces/kcp-shared-tools/tmc-pr-line-counter.sh',
            '-c', branch
        ], cwd=working_dir, capture_output=True, text=True)
        
        if result.returncode == 0:
            current_lines = int(result.stdout.strip().split()[-1])
            checkpoint_lines = checkpoint_data.get('measurement_execution', {}).get('measurement_result', 0)
            
            # Allow small variance (up to 10 lines)
            variance = abs(current_lines - checkpoint_lines)
            if variance > 10:
                verification_results['still_valid'] = False
                verification_results['issues_detected'].append(
                    f'Significant size change: {current_lines} vs checkpoint {checkpoint_lines} (variance: {variance})'
                )
            else:
                verification_results['current_measurement'] = current_lines
                verification_results['checkpoint_measurement'] = checkpoint_lines
                verification_results['variance'] = variance
        else:
            verification_results['still_valid'] = False
            verification_results['issues_detected'].append(f'Re-measurement failed: {result.stderr}')
            
    except Exception as e:
        verification_results['still_valid'] = False
        verification_results['issues_detected'].append(f'Verification error: {str(e)}')
    
    return verification_results

def detect_changes_since_measurement(checkpoint_data):
    """Detect changes since size measurement was performed"""
    
    measurement_time = datetime.fromisoformat(checkpoint_data['measurement_execution']['completed_at'])
    working_dir = checkpoint_data.get('working_dir', '')
    
    changes = {
        'new_commits': [],
        'modified_files': [],
        'build_changes': []
    }
    
    if not os.path.exists(working_dir):
        return changes
    
    try:
        # Check for commits since measurement
        git_log = subprocess.run([
            'git', 'log', '--since', measurement_time.isoformat(),
            '--pretty=format:%H|%s|%ai'
        ], cwd=working_dir, capture_output=True, text=True)
        
        if git_log.stdout.strip():
            for line in git_log.stdout.strip().split('\n'):
                parts = line.split('|')
                if len(parts) >= 3:
                    changes['new_commits'].append({
                        'sha': parts[0],
                        'message': parts[1],
                        'timestamp': parts[2]
                    })
        
        # Check for file modifications
        modified_files = find_files_modified_since(working_dir, measurement_time)
        changes['modified_files'] = modified_files
        
        # Check for build-related changes
        build_files = ['go.mod', 'go.sum', 'Makefile', 'Dockerfile']
        for build_file in build_files:
            build_path = os.path.join(working_dir, build_file)
            if os.path.exists(build_path):
                if os.path.getmtime(build_path) > measurement_time.timestamp():
                    changes['build_changes'].append({
                        'file': build_file,
                        'modified_at': datetime.fromtimestamp(os.path.getmtime(build_path)).isoformat()
                    })
    
    except Exception as e:
        changes['error'] = str(e)
    
    return changes

def determine_measure_size_recovery_actions(checkpoint, verification, changes):
    """Determine recovery actions needed for measure size state"""
    
    recovery_actions = []
    
    # Handle verification issues
    if not verification['still_valid']:
        for issue in verification['issues_detected']:
            if 'directory' in issue.lower():
                recovery_actions.append({
                    'type': 'RESTORE_WORKSPACE',
                    'description': issue,
                    'priority': 'CRITICAL'
                })
            elif 'size change' in issue.lower():
                recovery_actions.append({
                    'type': 'RE_MEASURE_SIZE',
                    'description': issue,
                    'priority': 'HIGH'
                })
            elif 'measurement failed' in issue.lower():
                recovery_actions.append({
                    'type': 'VERIFY_MEASUREMENT_TOOLS',
                    'description': issue,
                    'priority': 'HIGH'
                })
    
    # Handle detected changes
    if changes['new_commits']:
        recovery_actions.append({
            'type': 'VALIDATE_NEW_COMMITS',
            'description': f'Validate {len(changes["new_commits"])} commits since measurement',
            'priority': 'MEDIUM',
            'details': changes['new_commits']
        })
    
    if changes['modified_files']:
        recovery_actions.append({
            'type': 'RE_MEASURE_AFTER_CHANGES',
            'description': f'Re-measure size after {len(changes["modified_files"])} file modifications',
            'priority': 'HIGH',
            'details': changes['modified_files']
        })
    
    if changes['build_changes']:
        recovery_actions.append({
            'type': 'VERIFY_BUILD_IMPACT',
            'description': f'Verify build changes impact on size measurement',
            'priority': 'MEDIUM',
            'details': changes['build_changes']
        })
    
    # Check if decision is still valid
    decision = checkpoint.get('decision_analysis', {}).get('final_decision', {})
    if decision and verification['time_elapsed_minutes'] > 60:
        recovery_actions.append({
            'type': 'REVALIDATE_DECISION',
            'description': 'Measurement is over 1 hour old - revalidate decision',
            'priority': 'MEDIUM'
        })
    
    return recovery_actions
```

### Size Measurement Re-validation

```python
def revalidate_size_measurement(checkpoint_data):
    """Re-validate size measurement after recovery"""
    
    print("🔍 RE-VALIDATING SIZE MEASUREMENT")
    
    working_dir = checkpoint_data.get('working_dir', '')
    branch = checkpoint_data.get('branch', '')
    
    # Re-execute measurement
    current_measurement = execute_size_measurement(branch, working_dir)
    
    # Compare with checkpoint measurement
    checkpoint_measurement = checkpoint_data.get('measurement_execution', {})
    comparison = compare_measurements(checkpoint_measurement, current_measurement)
    
    # Re-evaluate decision if significant changes
    decision_revalidation = None
    if not comparison['consistent']:
        decision_revalidation = re_evaluate_size_decision(
            current_measurement, checkpoint_data.get('decision_analysis', {})
        )
    
    return {
        'revalidation_timestamp': datetime.now().isoformat(),
        'current_measurement': current_measurement,
        'checkpoint_comparison': comparison,
        'decision_revalidation': decision_revalidation,
        'action_required': 'UPDATE_DECISION' if decision_revalidation else 'CONTINUE_AS_PLANNED'
    }

def execute_size_measurement(branch, working_dir):
    """Execute fresh size measurement"""
    
    try:
        result = subprocess.run([
            '/workspaces/kcp-shared-tools/tmc-pr-line-counter.sh',
            '-c', branch
        ], cwd=working_dir, capture_output=True, text=True, check=True)
        
        lines = int(result.stdout.strip().split()[-1])
        
        return {
            'success': True,
            'measurement_result': lines,
            'tool_output': result.stdout.strip(),
            'utilization_percentage': (lines / 800) * 100,
            'status': determine_size_status(lines)
        }
        
    except subprocess.CalledProcessError as e:
        return {
            'success': False,
            'error': f'Measurement failed: {e.stderr}',
            'measurement_result': 0
        }
    except Exception as e:
        return {
            'success': False,
            'error': f'Unexpected error: {str(e)}',
            'measurement_result': 0
        }

def compare_measurements(checkpoint_measurement, current_measurement):
    """Compare checkpoint measurement with current measurement"""
    
    if not current_measurement.get('success', False):
        return {
            'consistent': False,
            'error': 'Current measurement failed',
            'variance': 'unknown'
        }
    
    checkpoint_lines = checkpoint_measurement.get('measurement_result', 0)
    current_lines = current_measurement.get('measurement_result', 0)
    
    variance = abs(current_lines - checkpoint_lines)
    percentage_change = (variance / checkpoint_lines) * 100 if checkpoint_lines > 0 else 0
    
    # Consider consistent if variance is small
    consistent = variance <= 10 or percentage_change <= 2.0
    
    return {
        'consistent': consistent,
        'checkpoint_lines': checkpoint_lines,
        'current_lines': current_lines,
        'variance': variance,
        'percentage_change': percentage_change,
        'significance': 'MINOR' if consistent else ('MAJOR' if variance > 50 else 'MODERATE')
    }

def determine_size_status(lines):
    """Determine size compliance status"""
    
    if lines > 800:
        return 'VIOLATION'
    elif lines > 750:
        return 'DANGER'
    elif lines > 700:
        return 'WARNING'
    else:
        return 'COMPLIANT'
```

## State Persistence

Save size measurement checkpoint with complete analysis context:

```bash
# Primary checkpoint location
CHECKPOINT_DIR="/workspaces/software-factory-2.0-template/checkpoints/active"
EFFORT_ID="effort2-controller"
CHECKPOINT_FILE="$CHECKPOINT_DIR/sw-engineer-measure-size-${EFFORT_ID}-$(date +%Y%m%d-%H%M%S).yaml"

# Measurement-specific backup (critical for decision validation)
BACKUP_DIR="/workspaces/software-factory-2.0-template/checkpoints/size-measurements"
mkdir -p "$BACKUP_DIR"
BACKUP_FILE="$BACKUP_DIR/size-measurement-${EFFORT_ID}-latest.yaml"

# Size analysis archive
ARCHIVE_DIR="/workspaces/software-factory-2.0-template/checkpoints/size-analysis"
mkdir -p "$ARCHIVE_DIR"
ANALYSIS_FILE="$ARCHIVE_DIR/size-analysis-${EFFORT_ID}-$(date +%Y%m%d-%H%M%S).yaml"

# Save checkpoint and analysis
cp "$CHECKPOINT_FILE" "$BACKUP_FILE"
cp "$CHECKPOINT_FILE" "$ANALYSIS_FILE"

# Generate measurement report
REPORT_FILE="/workspaces/efforts/${EFFORT_PATH}/SIZE-MEASUREMENT-REPORT.md"
generate_size_measurement_report "$CHECKPOINT_FILE" > "$REPORT_FILE"

# Update work log with measurement
echo "- [$(date '+%Y-%m-%d %H:%M')] SIZE MEASUREMENT: ${MEASURED_LINES}/800 lines (${UTILIZATION}%) - ${STATUS}" >> work-log.md

# Commit measurement results
git add .
git commit -m "measure: size analysis ${MEASURED_LINES}/800 lines - transition to ${NEXT_STATE}"
git push
```

## Health Monitoring

```python
def monitor_measure_size_health():
    """Monitor size measurement process health"""
    
    health_indicators = {
        'tool_compliance': check_measurement_tool_compliance(),
        'measurement_accuracy': assess_measurement_accuracy(),
        'decision_quality': evaluate_decision_quality(),
        'analysis_depth': measure_analysis_comprehensiveness()
    }
    
    overall_health = calculate_measure_size_health(health_indicators)
    
    if overall_health['status'] != 'HEALTHY':
        print(f"⚠️ SIZE MEASUREMENT HEALTH: {overall_health['status']}")
        for concern in overall_health['concerns']:
            print(f"  - {concern}")
    
    return overall_health

def check_measurement_tool_compliance():
    """Check for compliance with mandatory measurement tool"""
    
    # Check recent measurement history for tool compliance
    recent_measurements = get_recent_size_measurements()
    
    if not recent_measurements:
        return {'status': 'NO_DATA', 'compliance_rate': 0}
    
    compliant_count = sum(
        1 for measurement in recent_measurements
        if 'tmc-pr-line-counter.sh' in measurement.get('tool_used', '')
    )
    
    compliance_rate = (compliant_count / len(recent_measurements)) * 100
    
    if compliance_rate == 100:
        status = 'EXCELLENT'
    elif compliance_rate >= 90:
        status = 'GOOD'
    elif compliance_rate >= 80:
        status = 'ACCEPTABLE'
    else:
        status = 'POOR'
    
    return {
        'status': status,
        'compliance_rate': compliance_rate,
        'total_measurements': len(recent_measurements),
        'compliant_measurements': compliant_count
    }
```

## Critical Recovery Points

