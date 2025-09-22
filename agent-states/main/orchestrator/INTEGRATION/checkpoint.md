# Orchestrator - INTEGRATION State Checkpoint

## When to Save State

Save checkpoint at these critical points:

1. **Before Starting Integration**
   - All wave efforts validated and ready
   - Integration branch created
   - Conflict resolution strategy planned

2. **During Conflict Resolution**
   - Each conflict detected and analyzed
   - Resolution approach determined
   - Progress on multi-conflict scenarios

3. **After Each Effort Merge**
   - Individual effort successfully integrated
   - Cumulative test results
   - Size validation results

4. **At Integration Completion**
   - Final integration branch ready
   - All tests passing
   - Ready for architect review or next wave

## Required Data to Preserve

```yaml
integration_checkpoint:
  # State identification
  state: "INTEGRATION"
  phase: 1
  wave: 2
  checkpoint_timestamp: "2025-08-23T15:30:45Z"
  
  # Integration progress
  integration_branch: "phase1/wave2-integration"
  efforts_included:
    - name: "effort1-api-types"
      status: "MERGED"
      merge_commit: "abc123"
      conflicts: 0
    - name: "effort2-controller"
      status: "MERGING"
      conflicts: 2
      resolution_status: "IN_PROGRESS"
    - name: "effort3-webhooks"
      status: "PENDING"
      dependencies: ["effort2-controller"]
  
  # Conflict tracking
  conflicts_detected:
    - effort: "effort2-controller"
      files: ["pkg/controller/tenant.go", "api/types.go"]
      type: "OVERLAPPING_CHANGES"
      resolution_strategy: "SPAWN_REVIEWER"
      started_at: "2025-08-23T15:25:00Z"
      
  # Validation results
  size_validations:
    - effort: "effort1-api-types"
      lines: 245
      compliant: true
    - current_integration_total: 487
      limit: 800
      remaining_capacity: 313
      
  # Test results
  test_results:
    unit_tests:
      status: "PASSING"
      coverage: 87.5
    integration_tests:
      status: "RUNNING"
      progress: "60%"
    build:
      status: "SUCCESS"
      
  # Grading data
  performance_metrics:
    integration_started: "2025-08-23T15:20:00Z"
    conflicts_detected: 2
    resolution_time_so_far: "5min 30s"
    success_rate_current: 50.0  # 1 of 2 completed successfully
    
  # Next actions
  pending_actions:
    - type: "RESOLVE_CONFLICT"
      target: "effort2-controller"
      assigned_to: "code-reviewer"
      priority: "HIGH"
    - type: "MERGE_EFFORT"
      target: "effort3-webhooks"
      depends_on: "effort2-controller"
      priority: "MEDIUM"
    - type: "FINAL_VALIDATION"
      trigger: "ALL_EFFORTS_MERGED"
      priority: "HIGH"
```

## Recovery Protocol

### Context Recovery After Interruption

```python
def recover_integration_state(checkpoint_data):
    """Recover integration state from checkpoint"""
    
    print("🔄 RECOVERING INTEGRATION STATE")
    
    # Verify integration branch exists
    branch = checkpoint_data['integration_branch']
    verify_branch_exists(branch)
    
    # Check current status vs checkpoint
    current_efforts = analyze_current_integration_status(branch)
    checkpoint_efforts = checkpoint_data['efforts_included']
    
    # Identify what needs to be resumed
    pending_work = []
    
    for effort in checkpoint_efforts:
        current_status = get_current_effort_status(effort['name'])
        checkpoint_status = effort['status']
        
        if checkpoint_status == 'MERGING' and current_status != 'MERGED':
            pending_work.append({
                'type': 'RESUME_MERGE',
                'effort': effort['name'],
                'conflicts': effort.get('conflicts', 0)
            })
        elif checkpoint_status == 'PENDING':
            pending_work.append({
                'type': 'START_MERGE',
                'effort': effort['name']
            })
    
    # Resume conflict resolution
    active_conflicts = checkpoint_data.get('conflicts_detected', [])
    for conflict in active_conflicts:
        if conflict.get('resolution_status') == 'IN_PROGRESS':
            pending_work.append({
                'type': 'RESUME_CONFLICT_RESOLUTION',
                'conflict': conflict
            })
    
    return {
        'integration_branch': branch,
        'pending_work': pending_work,
        'current_progress': calculate_integration_progress(checkpoint_data),
        'next_action': determine_next_action(pending_work)
    }

def calculate_integration_progress(checkpoint):
    """Calculate integration completion percentage"""
    efforts = checkpoint['efforts_included']
    total_efforts = len(efforts)
    
    completed = sum(1 for e in efforts if e['status'] == 'MERGED')
    in_progress = sum(1 for e in efforts if e['status'] == 'MERGING')
    
    # Weight: completed=100%, in_progress=50%, pending=0%
    progress = ((completed * 100) + (in_progress * 50)) / (total_efforts * 100)
    
    return {
        'percentage': progress * 100,
        'completed': completed,
        'in_progress': in_progress,
        'pending': total_efforts - completed - in_progress
    }
```

### Checkpoint Validation

```python
def validate_checkpoint_integrity(checkpoint):
    """Validate checkpoint data completeness"""
    
    required_fields = [
        'state', 'phase', 'wave', 'integration_branch',
        'efforts_included', 'test_results', 'performance_metrics'
    ]
    
    missing_fields = []
    for field in required_fields:
        if field not in checkpoint:
            missing_fields.append(field)
    
    if missing_fields:
        raise ValueError(f"Checkpoint missing required fields: {missing_fields}")
    
    # Validate effort statuses
    valid_statuses = ['PENDING', 'MERGING', 'MERGED', 'FAILED']
    for effort in checkpoint['efforts_included']:
        if effort['status'] not in valid_statuses:
            raise ValueError(f"Invalid effort status: {effort['status']}")
    
    # Validate branch exists
    if not branch_exists(checkpoint['integration_branch']):
        raise ValueError(f"Integration branch not found: {checkpoint['integration_branch']}")
    
    return True
```

## State Persistence

Save checkpoint to multiple locations for redundancy:

```bash
# Primary location
CHECKPOINT_DIR="/workspaces/software-factory-2.0-template/checkpoints/active"
CHECKPOINT_FILE="$CHECKPOINT_DIR/orchestrator-integration-$(date +%Y%m%d-%H%M%S).yaml"

# Backup location  
BACKUP_DIR="/workspaces/software-factory-2.0-template/checkpoints/backup"
BACKUP_FILE="$BACKUP_DIR/orchestrator-integration-latest.yaml"

# Save to both locations
cp "$CHECKPOINT_FILE" "$BACKUP_FILE"

# Commit to git for persistence
git add checkpoints/
git commit -m "checkpoint: orchestrator INTEGRATION state - wave ${WAVE} progress"
git push
```

## Cleanup Strategy

```python
def cleanup_old_checkpoints():
    """Clean up old checkpoint files"""
    
    checkpoint_dir = "/workspaces/software-factory-2.0-template/checkpoints/active"
    
    # Keep last 10 integration checkpoints
    integration_files = sorted([
        f for f in os.listdir(checkpoint_dir)
        if f.startswith('orchestrator-integration-')
    ], reverse=True)
    
    # Remove old files
    for old_file in integration_files[10:]:
        os.remove(os.path.join(checkpoint_dir, old_file))
        print(f"Cleaned up old checkpoint: {old_file}")
```

## Critical Recovery Points

---
### 🚨 RULE  - 
**Source:** 
**Criticality:** CRITICAL - Major impact on grading

CRITICAL RECOVERY SCENARIOS
1. Mid-Merge Interruption:
- Check git merge status
- Resume or abort current merge
- Validate branch consistency

2. Conflict Resolution Interrupted:
- Check spawned agent status
- Resume conflict resolution
- Re-validate resolution approach

3. Test Failure During Integration:
- Identify failing tests
- Determine if integration or effort issue
- Plan fix or rollback strategy

4. Size Limit Exceeded:
- Identify which effort caused overflow
- Plan integration split strategy
- Coordinate with architect for guidance
---

## Monitoring Integration Health

```python
def monitor_integration_health():
    """Continuous health monitoring during integration"""
    
    health_checks = {
        'branch_integrity': check_branch_integrity(),
        'size_compliance': validate_current_size(),
        'test_status': get_current_test_status(),
        'conflict_count': count_active_conflicts(),
        'performance_impact': measure_performance_impact()
    }
    
    # Save health data to checkpoint
    checkpoint_health_update(health_checks)
    
    # Alert on issues
    for check, status in health_checks.items():
        if status['grade'] == 'FAIL':
            print(f"⚠️ INTEGRATION HEALTH ISSUE: {check}")
            print(f"Details: {status['details']}")
    
    return health_checks
