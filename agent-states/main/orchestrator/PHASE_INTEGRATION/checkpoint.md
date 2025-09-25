# PHASE_INTEGRATION State - Checkpoint Requirements

## Checkpoint Triggers

The PHASE_INTEGRATION state should create checkpoints at these critical moments:

1. **On State Entry** - Capture initial conditions
2. **After Wave Merges** - Record integration progress
3. **After Fix Merges** - Document ERROR_RECOVERY fixes applied
4. **Before State Exit** - Final integration status

## Checkpoint Data Structure

```yaml
checkpoint:
  state: "PHASE_INTEGRATION"
  timestamp: "2025-08-27T14:30:00Z"
  phase: 3
  checkpoint_type: "phase_integration_after_fixes"
  
  integration_context:
    trigger: "PHASE_ASSESSMENT_NEEDS_WORK"
    original_assessment_score: 68
    issues_to_fix: 5
    priority_1_issues: 3
    
  branch_info:
    integration_branch: "phase3-post-fixes-integration-20250827-143000"
    base_branch: "main"
    created_from_sha: "abc123def456"
    
  waves_integrated:
    - wave: 1
      branch: "phase3-wave1-integration-20250825-120000"
      merged_at: "2025-08-27T14:31:00Z"
      commits: 45
      
    - wave: 2
      branch: "phase3-wave2-integration-20250826-140000"  
      merged_at: "2025-08-27T14:32:00Z"
      commits: 52
      
    - wave: 3
      branch: "phase3-wave3-integration-20250827-090000"
      merged_at: "2025-08-27T14:33:00Z"
      commits: 38
      
    - wave: 4
      branch: "phase3-wave4-integration-20250827-110000"
      merged_at: "2025-08-27T14:34:00Z"
      commits: 41
      
  fixes_integrated:
    - fix_id: "ERR-PHASE-001"
      branch: "phase3-fix-kcp-patterns-20250827-120000"
      merged_at: "2025-08-27T14:35:00Z"
      addresses: "KCP multi-tenancy patterns missing"
      
    - fix_id: "ERR-PHASE-002"
      branch: "phase3-fix-api-compatibility-20250827-130000"
      merged_at: "2025-08-27T14:36:00Z"
      addresses: "API backwards compatibility"
      
    - fix_id: "ERR-PHASE-003"
      branch: "phase3-fix-test-coverage-20250827-135000"
      merged_at: "2025-08-27T14:37:00Z"
      addresses: "Test coverage below 60%"
      
  merge_conflicts:
    encountered: false
    resolved: []
    
  validation_status:
    all_waves_merged: true
    all_fixes_merged: true
    tests_passing: true
    ready_for_reassessment: true
    
  next_state: "SPAWN_ARCHITECT_PHASE_ASSESSMENT"
  
  metrics:
    total_integration_time_minutes: 8
    total_commits_integrated: 185
    files_changed: 127
    lines_added: 4250
    lines_deleted: 890
```

## Checkpoint Creation Function

```python
def create_phase_integration_checkpoint(checkpoint_type="state_entry"):
    """Create checkpoint for PHASE_INTEGRATION state"""
    
    # Read current state
    state_data = read_yaml('orchestrator-state.json')
    phase = state_data['current_phase']
    
    checkpoint = {
        'state': 'PHASE_INTEGRATION',
        'timestamp': datetime.now().isoformat(),
        'phase': phase,
        'checkpoint_type': checkpoint_type,
        'integration_context': {},
        'branch_info': {},
        'waves_integrated': [],
        'fixes_integrated': [],
        'merge_conflicts': {},
        'validation_status': {},
        'metrics': {}
    }
    
    if checkpoint_type == 'state_entry':
        # Capture initial conditions
        checkpoint['integration_context'] = {
            'trigger': state_data.get('previous_state'),
            'error_recovery_reason': state_data.get('error_recovery_reason'),
            'issues_to_fix': len(state_data.get('error_recovery_fixes', []))
        }
        
    elif checkpoint_type == 'after_wave_merges':
        # Record wave integration progress
        checkpoint['waves_integrated'] = get_merged_wave_branches(phase)
        checkpoint['merge_conflicts']['waves'] = check_merge_conflicts()
        
    elif checkpoint_type == 'after_fix_merges':
        # Document fixes applied
        checkpoint['fixes_integrated'] = get_merged_fix_branches(phase)
        checkpoint['merge_conflicts']['fixes'] = check_merge_conflicts()
        
    elif checkpoint_type == 'state_exit':
        # Final status before transition
        checkpoint['validation_status'] = {
            'all_waves_merged': verify_all_waves_merged(phase),
            'all_fixes_merged': verify_all_fixes_merged(phase),
            'tests_passing': run_integration_tests(),
            'ready_for_reassessment': True
        }
        checkpoint['next_state'] = 'SPAWN_ARCHITECT_PHASE_ASSESSMENT'
        checkpoint['metrics'] = calculate_integration_metrics()
    
    # Save checkpoint
    checkpoint_file = f"checkpoints/phase{phase}-integration-{checkpoint_type}-{datetime.now().strftime('%Y%m%d-%H%M%S')}.yaml"
    save_yaml(checkpoint, checkpoint_file)
    
    return checkpoint
```

## Checkpoint Recovery Function

```python
def recover_from_phase_integration_checkpoint(checkpoint_file):
    """Recover PHASE_INTEGRATION state from checkpoint"""
    
    checkpoint = read_yaml(checkpoint_file)
    
    print(f"🔄 Recovering PHASE_INTEGRATION from checkpoint: {checkpoint['timestamp']}")
    print(f"📊 Phase: {checkpoint['phase']}")
    print(f"📝 Checkpoint Type: {checkpoint['checkpoint_type']}")
    
    # Restore branch state
    if checkpoint['branch_info'].get('integration_branch'):
        branch = checkpoint['branch_info']['integration_branch']
        print(f"🔀 Checking out integration branch: {branch}")
        subprocess.run(['git', 'checkout', branch])
    
    # Report integration status
    waves_done = len(checkpoint.get('waves_integrated', []))
    fixes_done = len(checkpoint.get('fixes_integrated', []))
    
    print(f"✅ Waves integrated: {waves_done}")
    print(f"✅ Fixes integrated: {fixes_done}")
    
    # Determine what needs to be done
    if checkpoint['checkpoint_type'] == 'state_entry':
        print("📋 Need to start integration process")
        return 'START_INTEGRATION'
        
    elif checkpoint['checkpoint_type'] == 'after_wave_merges':
        print("📋 Waves merged, need to merge fixes")
        return 'MERGE_FIXES'
        
    elif checkpoint['checkpoint_type'] == 'after_fix_merges':
        print("📋 All merges done, need validation")
        return 'VALIDATE_AND_EXIT'
        
    elif checkpoint['checkpoint_type'] == 'state_exit':
        print("✅ Integration complete, ready for reassessment")
        return 'TRANSITION_TO_REASSESSMENT'
    
    return 'UNKNOWN'
```

## Checkpoint Validation

```python
def validate_phase_integration_checkpoint(checkpoint):
    """Validate checkpoint has required data"""
    
    required_fields = [
        'state',
        'timestamp',
        'phase',
        'checkpoint_type',
        'integration_context'
    ]
    
    validation_results = {
        'valid': True,
        'errors': [],
        'warnings': []
    }
    
    # Check required fields
    for field in required_fields:
        if field not in checkpoint:
            validation_results['valid'] = False
            validation_results['errors'].append(f"Missing required field: {field}")
    
    # Validate state
    if checkpoint.get('state') != 'PHASE_INTEGRATION':
        validation_results['valid'] = False
        validation_results['errors'].append(f"Wrong state: {checkpoint.get('state')}")
    
    # Check checkpoint type
    valid_types = ['state_entry', 'after_wave_merges', 'after_fix_merges', 'state_exit']
    if checkpoint.get('checkpoint_type') not in valid_types:
        validation_results['warnings'].append(f"Unusual checkpoint type: {checkpoint.get('checkpoint_type')}")
    
    # Validate integration completeness for exit checkpoint
    if checkpoint.get('checkpoint_type') == 'state_exit':
        if not checkpoint.get('validation_status', {}).get('all_waves_merged'):
            validation_results['errors'].append("Cannot exit without all waves merged")
            validation_results['valid'] = False
            
        if not checkpoint.get('validation_status', {}).get('all_fixes_merged'):
            validation_results['errors'].append("Cannot exit without all fixes merged")
            validation_results['valid'] = False
    
    return validation_results
```

## Checkpoint Storage Location

```bash
checkpoints/
├── active/
│   └── phase3-integration-current.yaml  # Symlink to latest
├── phase3/
│   ├── phase3-integration-state_entry-20250827-143000.yaml
│   ├── phase3-integration-after_wave_merges-20250827-143200.yaml
│   ├── phase3-integration-after_fix_merges-20250827-143500.yaml
│   └── phase3-integration-state_exit-20250827-143800.yaml
```

## Checkpoint Usage Examples

### Creating Checkpoints During Integration

```python
# On entering PHASE_INTEGRATION
def enter_phase_integration():
    print("🔀 Entering PHASE_INTEGRATION state")
    
    # Create entry checkpoint
    checkpoint = create_phase_integration_checkpoint("state_entry")
    print(f"✅ Entry checkpoint created: {checkpoint['timestamp']}")
    
    # Start integration work
    create_integration_branch()
    
    # Merge waves
    merge_wave_branches()
    checkpoint = create_phase_integration_checkpoint("after_wave_merges")
    print(f"✅ Wave merge checkpoint created")
    
    # Merge fixes
    merge_fix_branches()
    checkpoint = create_phase_integration_checkpoint("after_fix_merges")
    print(f"✅ Fix merge checkpoint created")
    
    # Final validation
    if validate_integration():
        checkpoint = create_phase_integration_checkpoint("state_exit")
        print(f"✅ Exit checkpoint created")
        transition_to_next_state()
```

### Recovering from Interruption

```python
# After context loss or interruption
def recover_phase_integration():
    # Find latest checkpoint
    latest_checkpoint = find_latest_checkpoint("PHASE_INTEGRATION")
    
    if not latest_checkpoint:
        print("❌ No checkpoint found, starting fresh")
        enter_phase_integration()
        return
    
    # Load and recover
    action = recover_from_phase_integration_checkpoint(latest_checkpoint)
    
    if action == 'START_INTEGRATION':
        create_integration_branch()
        merge_wave_branches()
        merge_fix_branches()
        
    elif action == 'MERGE_FIXES':
        merge_fix_branches()
        validate_integration()
        
    elif action == 'VALIDATE_AND_EXIT':
        validate_integration()
        transition_to_next_state()
        
    elif action == 'TRANSITION_TO_REASSESSMENT':
        transition_to_next_state()
```

## Critical Checkpoint Rules

1. **Always create entry checkpoint** - Capture starting conditions
2. **Checkpoint after major operations** - Waves merged, fixes merged
3. **Exit checkpoint is mandatory** - Proves integration complete
4. **Include validation status** - Tests, merge conflicts, readiness
5. **Track metrics** - Time, commits, changes for grading

## Checkpoint Cleanup

```bash
#!/bin/bash
# Clean old checkpoints (keep last 5 per phase)

PHASE=$1
CHECKPOINT_DIR="checkpoints/phase${PHASE}"

if [ -d "$CHECKPOINT_DIR" ]; then
    # Keep only last 5 checkpoints
    ls -t "$CHECKPOINT_DIR"/*.yaml | tail -n +6 | xargs -r rm
    echo "✅ Cleaned old checkpoints for phase ${PHASE}"
fi
```