# ORCHESTRATOR CONTEXT-SWITCHING TESTING PROTOCOL

## Critical Testing Principle

**REMEMBER**: The orchestrator has ZERO MEMORY between states. Each test must validate that a completely fresh orchestrator can function with only:
1. The minimal bootstrap rules
2. The orchestrator-state.json file
3. The state-specific rules for its current state

## Test Environment Setup

### Prerequisites
```bash
# 1. Create test environment
mkdir -p /tmp/orchestrator-testing
cd /tmp/orchestrator-testing

# 2. Create mock orchestrator-state.json files for different states
cat > state-init.yaml << 'EOF'
current_state: INIT
current_phase: 1
current_wave: 1
EOF

cat > state-spawn.yaml << 'EOF'
current_state: SPAWN_AGENTS
current_phase: 1
current_wave: 1
efforts_planned:
  - effort1-auth
  - effort2-database
EOF

cat > state-monitor.yaml << 'EOF'
current_state: MONITOR_IMPLEMENTATION
current_phase: 1
current_wave: 1
efforts_in_progress:
  - effort1-auth
  - effort2-database
EOF

# 3. Create test TODO file for recovery testing
cat > test-todos.todo << 'EOF'
- [ ] Complete effort implementation
- [ ] Run size checks
- [ ] Update work logs
EOF
```

## Test Suite 1: Cold Start Testing

### Test 1.1: Fresh Start with No Prior Context
```markdown
**Scenario**: Orchestrator starts for the first time
**Initial State**: No orchestrator-state.json exists
**Expected Behavior**:

1. Start orchestrator
2. Read minimal bootstrap (7 rules)
3. Detect no state file
4. Enter INIT state
5. Read INIT state rules
6. Create complete state file (R281)
7. Validate all phases/waves/efforts included

**Validation Checklist**:
- [ ] Only bootstrap rules loaded initially
- [ ] INIT rules loaded after state determination
- [ ] State file created with ALL items from plan
- [ ] No missing rules for initialization
```

### Test 1.2: Resume from Existing State
```markdown
**Scenario**: Orchestrator resumes from SPAWN_AGENTS
**Initial State**: orchestrator-state.json shows SPAWN_AGENTS
**Expected Behavior**:

1. Start orchestrator (no prior memory)
2. Read minimal bootstrap
3. Read orchestrator-state.json
4. Determine state = SPAWN_AGENTS
5. Read SPAWN_AGENTS/rules.md
6. Have all rules needed for spawning (R208, R151, etc.)
7. Execute spawn operations

**Validation Checklist**:
- [ ] Did NOT load INIT rules
- [ ] Did NOT load MONITOR rules
- [ ] DID load all SPAWN rules
- [ ] Can execute spawning correctly
```

## Test Suite 2: State Transition Testing

### Test 2.1: Normal State Transition
```markdown
**Scenario**: Transition from SPAWN_AGENTS to MONITOR
**Process**:

PART A - Complete SPAWN_AGENTS:
1. Orchestrator in SPAWN_AGENTS state
2. Complete spawning work
3. Update orchestrator-state.json to MONITOR
4. Save TODOs (R287)
5. Commit changes (R288)
6. STOP per R322
7. Orchestrator exits completely

PART B - Start fresh in MONITOR:
1. User runs /continue-orchestrating
2. NEW orchestrator starts (zero memory)
3. Reads minimal bootstrap
4. Reads state file (now shows MONITOR)
5. Reads MONITOR/rules.md
6. Has monitoring rules (R319, R006)
7. Does NOT have spawning rules

**Validation Checklist**:
- [ ] Clean stop after state update
- [ ] New instance has no memory of spawn
- [ ] Correct rules loaded for MONITOR
- [ ] No spawn rules in memory
```

### Test 2.2: Mandatory Sequence Transition
```markdown
**Scenario**: R234 Mandatory sequence enforcement
**States**: SETUP_EFFORT_INFRASTRUCTURE → ANALYZE_CODE_REVIEWER_PARALLELIZATION

PART A - In SETUP_EFFORT_INFRASTRUCTURE:
1. Has R234 in state rules
2. Knows next state MUST be ANALYZE_CODE_REVIEWER_PARALLELIZATION
3. Updates state file correctly
4. Stops

PART B - Fresh start:
1. Reads new state
2. Loads ANALYZE rules
3. Has R234 for next transition
4. Continues sequence

**Validation Checklist**:
- [ ] R234 present in sequence states
- [ ] Sequence enforced correctly
- [ ] No skipping allowed
```

## Test Suite 3: Recovery Testing

### Test 3.1: Recovery from Compaction
```markdown
**Scenario**: Context lost mid-state, TODOs saved
**Initial State**: In MONITOR, context compacted
**Recovery Process**:

1. Orchestrator starts fresh (no memory)
2. Reads minimal bootstrap (includes R287)
3. Detects TODO file exists
4. Loads TODOs into TodoWrite tool
5. Reads orchestrator-state.json (MONITOR)
6. Reads MONITOR/rules.md
7. Continues monitoring work

**Validation Checklist**:
- [ ] TODO detection works
- [ ] TODOs properly loaded
- [ ] Correct state rules loaded
- [ ] Work continues correctly
```

### Test 3.2: Recovery with Missing State Rules
```markdown
**Scenario**: State rules file deleted/corrupted
**Initial State**: orchestrator-state.json shows SPAWN_AGENTS
**Expected Behavior**:

1. Read bootstrap
2. Read state file
3. Attempt to read SPAWN_AGENTS/rules.md
4. File missing/corrupted
5. Transition to ERROR_RECOVERY
6. Report specific error
7. Request manual intervention

**Validation Checklist**:
- [ ] Error detected properly
- [ ] Clear error message
- [ ] Safe transition to ERROR_RECOVERY
- [ ] No undefined behavior
```

## Test Suite 4: Rule Loading Verification

### Test 4.1: Bootstrap Minimal Load
```markdown
**Test**: Verify ONLY essential rules in bootstrap
**Process**:

1. Start orchestrator
2. After reading orchestrator.md
3. Before reading state rules
4. Check loaded rules

**Must Have**:
- R283 (Complete file reading)
- R290 (State rule reading)
- R203 (State-aware startup)
- R206 (State validation)
- R288 (State file updates)
- R287 (TODO persistence)
- R322 (Stop before transitions)

**Must NOT Have**:
- R208 (CD before spawn)
- R151 (Parallel timing)
- R281 (State initialization)
- Other state-specific rules
```

### Test 4.2: State-Specific Load Verification
```markdown
**Test**: Each state has complete rules
**Process**: For each state directory

1. Simulate entering that state
2. Load only bootstrap + state rules
3. Verify can execute all state operations
4. Check for missing rules
5. Check for unnecessary rules

**Example - SPAWN_AGENTS must have**:
- From bootstrap: R203, R206, R288, R287, R290, R283, R322
- From state: R208, R151, R052, R197, R295, R309, R221
- NOT have: R281, R319, R307, etc.
```

## Test Suite 5: Edge Cases

### Test 5.1: Corrupted State File
```markdown
**Scenario**: orchestrator-state.json corrupted
**Expected**:
1. Detect corruption
2. Attempt recovery
3. If unrecoverable, transition to ERROR_RECOVERY
4. Clear error reporting
```

### Test 5.2: Invalid State in State File
```markdown
**Scenario**: State file contains "INVALID_STATE"
**Expected**:
1. Read state file
2. Validate against SOFTWARE-FACTORY-STATE-MACHINE.md
3. Detect invalid state
4. Report error
5. Transition to ERROR_RECOVERY
```

### Test 5.3: Rapid Context Switching
```markdown
**Scenario**: Multiple quick transitions
**Process**:
1. INIT → PLANNING (stop, restart)
2. PLANNING → SETUP (stop, restart)
3. SETUP → ANALYZE (stop, restart)
4. Each restart has correct rules
5. No rule accumulation
6. No missing rules
```

## Test Suite 6: Performance Testing

### Test 6.1: Startup Time Comparison
```markdown
**Measure**:
- Time to load full orchestrator.md (878 lines)
- Time to load minimal orchestrator.md (200 lines)
- Time to load state rules
- Total startup time difference

**Target**: <20% increase acceptable
```

### Test 6.2: Memory Usage
```markdown
**Measure**:
- Memory with all rules loaded
- Memory with minimal + state rules
- Reduction percentage

**Target**: >50% reduction expected
```

## Automated Test Script

```bash
#!/bin/bash
# orchestrator-migration-test.sh

echo "=== ORCHESTRATOR MIGRATION TEST SUITE ==="

# Test 1: Cold Start
test_cold_start() {
    echo "TEST 1: Cold Start"
    rm -f orchestrator-state.json
    
    # Simulate orchestrator start
    # Should load bootstrap, detect no state, enter INIT
    # Verify only INIT rules loaded
}

# Test 2: State Resume
test_state_resume() {
    echo "TEST 2: State Resume"
    cp state-spawn.yaml orchestrator-state.json
    
    # Simulate orchestrator start
    # Should load bootstrap, read state, load SPAWN rules
    # Verify has spawn rules, not others
}

# Test 3: Transition
test_transition() {
    echo "TEST 3: State Transition"
    
    # Part A: Complete current state
    cp state-spawn.yaml orchestrator-state.json
    # Simulate work completion
    # Update to MONITOR
    # Stop
    
    # Part B: Fresh start
    # Should have MONITOR rules only
}

# Test 4: Recovery
test_recovery() {
    echo "TEST 4: Recovery from Compaction"
    cp state-monitor.yaml orchestrator-state.json
    cp test-todos.todo todos/orchestrator-MONITOR-*.todo
    
    # Simulate fresh start
    # Should detect and load TODOs
    # Continue with MONITOR rules
}

# Run all tests
test_cold_start
test_state_resume
test_transition
test_recovery

echo "=== TEST SUITE COMPLETE ==="
```

## Manual Testing Checklist

### Phase 1: Pre-Migration Baseline
- [ ] Document current startup time
- [ ] Document current memory usage
- [ ] Document current rule loading
- [ ] Save as baseline metrics

### Phase 2: Post-Migration Validation
For each orchestrator state (all 67):
- [ ] Can enter state from fresh start
- [ ] Has all needed rules
- [ ] No missing functionality
- [ ] Can transition to valid next states
- [ ] Proper error handling

### Phase 3: Workflow Testing
Complete workflows:
- [ ] Full phase implementation
- [ ] Wave with splits
- [ ] Integration with conflicts
- [ ] Error recovery
- [ ] Parallel operations

### Phase 4: Stress Testing
- [ ] Rapid state transitions
- [ ] Large state files
- [ ] Many TODOs
- [ ] Corrupted files
- [ ] Missing dependencies

## Success Criteria

### Functional Success
- ✅ All states accessible
- ✅ All transitions work
- ✅ No missing rules
- ✅ Recovery works
- ✅ Error handling works

### Performance Success
- ✅ Startup time acceptable
- ✅ Memory usage reduced
- ✅ No performance degradation

### Maintainability Success
- ✅ Clearer rule organization
- ✅ Easier to debug
- ✅ Simpler to modify
- ✅ Better documentation

## Rollback Testing

### Rollback Trigger Conditions
1. Any state unreachable
2. Missing critical rules
3. Performance degradation >50%
4. Recovery failures
5. Transition failures

### Rollback Process
1. Stop all orchestrator instances
2. Restore original orchestrator.md
3. Restore original state rules
4. Verify system functional
5. Document failure reason

## Test Report Template

```markdown
# Orchestrator Migration Test Report

**Date**: [DATE]
**Tester**: [NAME]
**Version**: [pre/post migration]

## Test Results Summary
- Total Tests: X
- Passed: Y
- Failed: Z
- Skipped: W

## Detailed Results
[Test-by-test results]

## Performance Metrics
- Startup Time: Xms → Yms
- Memory Usage: XMB → YMB
- Rule Load Time: Xms → Yms

## Issues Found
[List any issues]

## Recommendation
[ ] Proceed with migration
[ ] Fix issues and retest
[ ] Rollback required
```

This comprehensive testing protocol ensures the migration maintains all functionality while improving efficiency through context-aware rule loading.