# Orchestrator - MONITOR State Checkpoint

## When to Save State

Save checkpoint at these critical monitoring points:

1. **Monitoring Session Start**
   - All agents identified and baseline established
   - Expected timelines and dependencies mapped
   - Monitoring schedule and thresholds configured

2. **Status Change Detection**
   - Agent status changes detected
   - Progress milestones reached
   - Dependency conditions met

3. **Intervention Triggers**
   - Issues detected requiring intervention
   - Agent spawn decisions made
   - Critical errors identified

4. **Wave Progress Milestones**
   - 25%, 50%, 75% completion points
   - Major dependency transitions
   - Timeline adjustments

## Required Data to Preserve

```yaml
monitor_checkpoint:
  # State identification
  state: "MONITOR"
  phase: 1
  wave: 2
  checkpoint_timestamp: "2025-08-23T15:45:30Z"
  
  # Monitoring session context
  session_info:
    started_at: "2025-08-23T14:30:00Z"
    session_id: "monitor-phase1-wave2-001"
    expected_duration_hours: 4
    last_full_check: "2025-08-23T15:40:00Z"
    check_interval_minutes: 10
    
  # Agent monitoring state
  monitored_agents:
    sw-engineer-effort1:
      agent_id: "sw-engineer-effort1"
      status: "ACTIVE"
      working_dir: "/workspaces/efforts/phase1/wave2/effort1-api-types"
      branch: "phase1/wave2/effort1-api-types"
      started_at: "2025-08-23T14:35:00Z"
      expected_completion: "2025-08-23T16:35:00Z"
      
      # Progress tracking
      progress:
        percentage: 75
        last_milestone: "API types defined"
        next_milestone: "Validation logic complete"
        milestones_completed: 3
        milestones_total: 4
        
      # Timeline status
      timeline:
        utilization: 0.65  # 65% of allocated time used
        grade: "ON_TIME"
        projected_completion: "2025-08-23T16:20:00Z"
        
      # Size monitoring
      size:
        current_lines: 445
        limit: 800
        compliant: true
        last_check: "2025-08-23T15:30:00Z"
        
      # Health indicators
      health:
        responsive: true
        last_response: "2025-08-23T15:44:00Z"
        response_time_avg: 12.3
        grade: "GOOD"
        
    code-reviewer-effort2:
      agent_id: "code-reviewer-effort2"
      status: "ACTIVE"
      working_dir: "/workspaces/efforts/phase1/wave2/effort2-controller"
      branch: "phase1/wave2/effort2-controller"
      started_at: "2025-08-23T14:40:00Z"
      expected_completion: "2025-08-23T15:40:00Z"
      
      progress:
        percentage: 90
        last_milestone: "Review complete"
        next_milestone: "Final approval"
        
      timeline:
        utilization: 0.95
        grade: "APPROACHING_DEADLINE"
        
      size:
        reviewed_lines: 756
        compliant: true
        
      health:
        responsive: true
        grade: "GOOD"
        
  # Dependency tracking
  dependency_status:
    effort1_to_effort2:
      prerequisite: "effort1-api-types"
      dependent: "effort2-controller"
      status: "SATISFIED"  # effort1 completed
      satisfied_at: "2025-08-23T15:30:00Z"
      
    effort1_to_effort3:
      prerequisite: "effort1-api-types"
      dependent: "effort3-webhooks"
      status: "READY_TO_START"
      can_start_at: "2025-08-23T15:30:00Z"
      
  # Intervention history
  interventions:
    - intervention_id: "INT-001"
      triggered_at: "2025-08-23T15:15:00Z"
      agent: "sw-engineer-effort1"
      issue: "PROGRESS_SLIGHTLY_BEHIND"
      action_taken: "STATUS_CHECK_REQUESTED"
      resolution: "RESOLVED"
      resolved_at: "2025-08-23T15:20:00Z"
      
  # Performance metrics
  performance:
    detection_times:
      - {change: "effort1_milestone", detected_in_seconds: 45}
      - {change: "effort2_status", detected_in_seconds: 32}
    
    intervention_times:
      - {issue: "progress_behind", response_seconds: 180}
      
    overall_session_health: "GOOD"
    
  # Pending actions
  pending_actions:
    - action: "START_EFFORT3"
      trigger: "effort1_completion_confirmed"
      ready: true
      scheduled_for: "2025-08-23T15:50:00Z"
      
    - action: "SPAWN_CODE_REVIEWER_EFFORT1"
      trigger: "effort1_implementation_complete"
      ready: false
      estimated_trigger: "2025-08-23T16:20:00Z"
      
  # Wave-level status
  wave_status:
    overall_progress: 68
    timeline_status: "ON_TRACK"
    completion_projection: "2025-08-23T17:30:00Z"
    risk_factors:
      - "effort2 approaching deadline"
      - "effort3 not yet started"
    mitigation_plans:
      - "Monitor effort2 closely for deadline extension"
      - "Start effort3 as soon as effort1 confirmed complete"
```

## Recovery Protocol

### Context Recovery After Interruption

```python
def recover_monitor_state(checkpoint_data):
    """Recover monitoring state from checkpoint"""
    
    print("🔄 RECOVERING MONITOR STATE")
    
    session_info = checkpoint_data['session_info']
    monitored_agents = checkpoint_data['monitored_agents']
    
    print(f"Resuming monitoring session: {session_info['session_id']}")
    print(f"Session started: {session_info['started_at']}")
    print(f"Agents being monitored: {len(monitored_agents)}")
    
    # Verify all agents are still accessible
    agent_verification = verify_agents_accessible(monitored_agents)
    
    # Check current status vs checkpoint status
    status_comparison = compare_current_vs_checkpoint_status(
        monitored_agents, checkpoint_data['checkpoint_timestamp']
    )
    
    # Determine what monitoring activities need to resume
    resume_actions = determine_monitoring_resume_actions(
        checkpoint_data, agent_verification, status_comparison
    )
    
    # Calculate elapsed time since checkpoint
    checkpoint_time = datetime.fromisoformat(checkpoint_data['checkpoint_timestamp'])
    elapsed_since_checkpoint = datetime.now() - checkpoint_time
    
    return {
        'session_id': session_info['session_id'],
        'agents_verified': agent_verification,
        'status_changes_detected': status_comparison['changes'],
        'resume_actions': resume_actions,
        'elapsed_since_checkpoint': elapsed_since_checkpoint.total_seconds() / 60,
        'monitoring_continuation': 'READY'
    }

def verify_agents_accessible(monitored_agents):
    """Verify all monitored agents are still accessible"""
    
    verification_results = {}
    
    for agent_id, agent_data in monitored_agents.items():
        try:
            # Check if working directory exists
            working_dir = agent_data['working_dir']
            wd_exists = os.path.exists(working_dir)
            
            # Check if branch exists
            branch = agent_data['branch']
            branch_exists = check_branch_exists(branch)
            
            # Try to ping agent (if communication protocol exists)
            agent_responsive = ping_agent(agent_id, timeout=10)
            
            verification_results[agent_id] = {
                'accessible': wd_exists and branch_exists,
                'working_dir_exists': wd_exists,
                'branch_exists': branch_exists,
                'agent_responsive': agent_responsive,
                'status': 'VERIFIED' if (wd_exists and branch_exists) else 'ISSUES_DETECTED'
            }
            
        except Exception as e:
            verification_results[agent_id] = {
                'accessible': False,
                'status': 'ERROR',
                'error': str(e)
            }
    
    return verification_results

def compare_current_vs_checkpoint_status(monitored_agents, checkpoint_timestamp):
    """Compare current agent status with checkpoint status"""
    
    status_changes = []
    
    for agent_id, checkpoint_data in monitored_agents.items():
        try:
            current_status = get_current_agent_status(agent_id)
            checkpoint_status = checkpoint_data['status']
            
            if current_status != checkpoint_status:
                status_changes.append({
                    'agent': agent_id,
                    'checkpoint_status': checkpoint_status,
                    'current_status': current_status,
                    'change_type': determine_change_significance(
                        checkpoint_status, current_status
                    )
                })
                
            # Check progress changes
            current_progress = get_current_progress_percentage(agent_id)
            checkpoint_progress = checkpoint_data['progress']['percentage']
            
            progress_change = current_progress - checkpoint_progress
            if abs(progress_change) > 5:  # >5% change
                status_changes.append({
                    'agent': agent_id,
                    'type': 'PROGRESS_CHANGE',
                    'checkpoint_progress': checkpoint_progress,
                    'current_progress': current_progress,
                    'change': progress_change
                })
                
        except Exception as e:
            status_changes.append({
                'agent': agent_id,
                'type': 'STATUS_CHECK_ERROR',
                'error': str(e)
            })
    
    return {
        'changes': status_changes,
        'agents_with_changes': len(status_changes),
        'requires_intervention': any(
            change.get('change_type') == 'CRITICAL' for change in status_changes
        )
    }

def determine_monitoring_resume_actions(checkpoint, verification, comparison):
    """Determine what actions need to be taken to resume monitoring"""
    
    resume_actions = []
    
    # Handle agent verification issues
    for agent_id, verify_result in verification.items():
        if verify_result['status'] == 'ERROR':
            resume_actions.append({
                'type': 'INVESTIGATE_AGENT_ERROR',
                'agent': agent_id,
                'priority': 'HIGH',
                'details': verify_result['error']
            })
        elif not verify_result['accessible']:
            resume_actions.append({
                'type': 'RESTORE_AGENT_ACCESS',
                'agent': agent_id,
                'priority': 'HIGH',
                'missing': [k for k, v in verify_result.items() 
                           if k.endswith('_exists') and not v]
            })
    
    # Handle status changes detected
    for change in comparison['changes']:
        if change.get('change_type') == 'CRITICAL':
            resume_actions.append({
                'type': 'IMMEDIATE_INTERVENTION',
                'agent': change['agent'],
                'priority': 'CRITICAL',
                'change': change
            })
        elif change.get('type') == 'PROGRESS_CHANGE':
            if change['change'] < -10:  # Progress went backwards significantly
                resume_actions.append({
                    'type': 'INVESTIGATE_PROGRESS_REGRESSION',
                    'agent': change['agent'],
                    'priority': 'HIGH'
                })
    
    # Resume pending actions from checkpoint
    pending = checkpoint.get('pending_actions', [])
    for action in pending:
        if action.get('ready'):
            resume_actions.append({
                'type': 'RESUME_PENDING_ACTION',
                'action': action,
                'priority': 'MEDIUM'
            })
    
    # Re-establish monitoring schedule
    resume_actions.append({
        'type': 'RESTART_MONITORING_CYCLE',
        'interval': checkpoint['session_info']['check_interval_minutes'],
        'priority': 'HIGH'
    })
    
    return resume_actions
```

### Monitoring State Validation

```python
def validate_monitoring_state_integrity(checkpoint):
    """Validate monitoring checkpoint data integrity"""
    
    validation_results = {
        'valid': True,
        'issues': [],
        'warnings': []
    }
    
    # Check required fields
    required_sections = ['session_info', 'monitored_agents', 'dependency_status']
    for section in required_sections:
        if section not in checkpoint:
            validation_results['issues'].append(f"Missing required section: {section}")
            validation_results['valid'] = False
    
    # Validate agent data completeness
    for agent_id, agent_data in checkpoint.get('monitored_agents', {}).items():
        required_agent_fields = ['status', 'working_dir', 'progress', 'timeline', 'health']
        for field in required_agent_fields:
            if field not in agent_data:
                validation_results['issues'].append(
                    f"Agent {agent_id} missing {field} data"
                )
    
    # Check timestamp consistency
    checkpoint_time = checkpoint.get('checkpoint_timestamp')
    if checkpoint_time:
        checkpoint_dt = datetime.fromisoformat(checkpoint_time)
        session_start = datetime.fromisoformat(
            checkpoint['session_info']['started_at']
        )
        
        if checkpoint_dt < session_start:
            validation_results['issues'].append(
                "Checkpoint timestamp before session start"
            )
    
    # Validate dependency consistency
    dependencies = checkpoint.get('dependency_status', {})
    monitored_agents = checkpoint.get('monitored_agents', {})
    
    for dep_id, dep_data in dependencies.items():
        prerequisite = dep_data.get('prerequisite')
        dependent = dep_data.get('dependent')
        
        # Check if referenced agents are being monitored
        if prerequisite and prerequisite not in monitored_agents:
            validation_results['warnings'].append(
                f"Dependency references unmonitored prerequisite: {prerequisite}"
            )
        if dependent and dependent not in monitored_agents:
            validation_results['warnings'].append(
                f"Dependency references unmonitored dependent: {dependent}"
            )
    
    return validation_results
```

## State Persistence

Save monitoring checkpoint to multiple locations for reliability:

```bash
# Primary location
CHECKPOINT_DIR="/workspaces/software-factory-2.0-template/checkpoints/active"
SESSION_ID=$(grep -o 'session_id: "[^"]*"' orchestrator-state.json | cut -d'"' -f2)
CHECKPOINT_FILE="$CHECKPOINT_DIR/orchestrator-monitor-${SESSION_ID}-$(date +%Y%m%d-%H%M%S).yaml"

# Backup location
BACKUP_DIR="/workspaces/software-factory-2.0-template/checkpoints/monitor-backup"
mkdir -p "$BACKUP_DIR"
BACKUP_FILE="$BACKUP_DIR/monitor-state-latest.yaml"

# Performance data location
METRICS_DIR="/workspaces/software-factory-2.0-template/checkpoints/metrics"
mkdir -p "$METRICS_DIR"
METRICS_FILE="$METRICS_DIR/monitor-performance-$(date +%Y%m%d).json"

# Save to all locations
cp "$CHECKPOINT_FILE" "$BACKUP_FILE"

# Extract and save performance metrics separately
extract_monitoring_metrics "$CHECKPOINT_FILE" > "$METRICS_FILE"

# Commit to git
git add checkpoints/
git commit -m "checkpoint: MONITOR state - ${SESSION_ID} wave progress"
git push
```

## Monitoring Health Checks

```python
def monitor_checkpoint_health():
    """Monitor the health of monitoring checkpoints themselves"""
    
    checkpoint_health = {
        'checkpoint_frequency': check_checkpoint_frequency(),
        'data_completeness': check_checkpoint_completeness(),
        'recovery_readiness': test_recovery_capability(),
        'storage_health': check_checkpoint_storage()
    }
    
    overall_health = calculate_checkpoint_health_grade(checkpoint_health)
    
    if overall_health['grade'] != 'PASS':
        print("⚠️ CHECKPOINT HEALTH ISSUE")
        for issue in overall_health['issues']:
            print(f"  - {issue}")
    
    return checkpoint_health

def check_checkpoint_frequency():
    """Check if checkpoints are being created at proper frequency"""
    
    checkpoint_dir = "/workspaces/software-factory-2.0-template/checkpoints/active"
    monitor_files = [f for f in os.listdir(checkpoint_dir) 
                    if f.startswith('orchestrator-monitor-')]
    
    if not monitor_files:
        return {'status': 'NO_CHECKPOINTS', 'grade': 'FAIL'}
    
    # Check most recent checkpoint age
    latest_file = max(monitor_files, key=lambda f: os.path.getmtime(
        os.path.join(checkpoint_dir, f)
    ))
    
    age_minutes = (time.time() - os.path.getmtime(
        os.path.join(checkpoint_dir, latest_file)
    )) / 60
    
    if age_minutes > 15:  # Should checkpoint every 10-15 minutes during monitoring
        return {'status': 'STALE_CHECKPOINTS', 'age_minutes': age_minutes, 'grade': 'FAIL'}
    else:
        return {'status': 'CURRENT', 'age_minutes': age_minutes, 'grade': 'PASS'}
```

## Critical Recovery Points

---
### 🚨 RULE  - 
**Source:** 
**Criticality:** CRITICAL - Major impact on grading

CRITICAL MONITORING RECOVERY SCENARIOS
1. Agent Communication Loss:
- Multiple agents unresponsive
- Re-establish communication channels
- Resume monitoring with status reconciliation

2. Monitoring Session Interrupted:
- Unexpected orchestrator termination
- Recover monitoring state from latest checkpoint
- Reconcile agent status changes during downtime

3. Dependency State Corruption:
- Dependency tracking inconsistent
- Rebuild dependency graph from agent statuses
- Validate dependency satisfaction

4. Performance Degradation:
- Monitoring response times increasing
- Switch to lightweight monitoring mode
- Maintain essential tracking only
---

## Cleanup and Maintenance

```python
def cleanup_monitoring_checkpoints():
    """Clean up old monitoring checkpoints"""
    
    checkpoint_dir = "/workspaces/software-factory-2.0-template/checkpoints/active"
    
    # Keep last 20 monitoring checkpoints per session
    monitor_files = [f for f in os.listdir(checkpoint_dir) 
                    if f.startswith('orchestrator-monitor-')]
    
    # Group by session
    sessions = {}
    for filename in monitor_files:
        session_id = extract_session_id_from_filename(filename)
        if session_id not in sessions:
            sessions[session_id] = []
        sessions[session_id].append(filename)
    
    # Keep last 20 per session, remove older
    for session_id, files in sessions.items():
        files.sort(key=lambda f: os.path.getmtime(os.path.join(checkpoint_dir, f)))
        
        for old_file in files[:-20]:  # Keep last 20
            file_age_hours = (time.time() - os.path.getmtime(
                os.path.join(checkpoint_dir, old_file)
            )) / 3600
            
            # Only remove if >6 hours old
            if file_age_hours > 6:
                os.remove(os.path.join(checkpoint_dir, old_file))
                print(f"Cleaned up old monitoring checkpoint: {old_file}")
