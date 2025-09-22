# 🚨 RULE R008 - Monitoring Frequency Requirements

**Criticality:** BLOCKING - Required for orchestrator operation  
**Grading Impact:** -5% per missed check, -20% for pattern  
**Enforcement:** CONTINUOUS - Every 5 messages

## Rule Statement

The Orchestrator MUST check agent progress EVERY 5 messages during active work periods.

## Monitoring Requirements

### Frequency Mandates
- **During Active Implementation**: Every 5 messages
- **During Reviews**: Every 3 messages  
- **During Integration**: Every 2 messages
- **While Waiting**: Every 10 messages maximum

### What to Monitor
```bash
perform_monitoring_check() {
    echo "📊 MONITORING CHECK (R008 Compliance)"
    echo "========================="
    
    # 1. Agent Status
    echo "🔍 Checking agent status..."
    check_agent_health
    
    # 2. Progress Metrics
    echo "📈 Checking progress..."
    check_effort_progress
    check_line_count
    check_test_coverage
    
    # 3. Blockers
    echo "🚧 Checking for blockers..."
    check_for_errors
    check_for_size_violations
    check_for_review_issues
    
    # 4. TODO State
    echo "📋 Checking TODO state..."
    verify_todo_saved
    
    echo "========================="
}
```

## Monitoring Triggers

### Message Count Based
```bash
MESSAGE_COUNT=0

after_each_message() {
    MESSAGE_COUNT=$((MESSAGE_COUNT + 1))
    
    # R008: Check every 5 messages
    if [ $((MESSAGE_COUNT % 5)) -eq 0 ]; then
        echo "⏰ R008: Monitoring check required (message #$MESSAGE_COUNT)"
        perform_monitoring_check
    fi
}
```

### Time Based (Backup)
```bash
LAST_MONITOR_TIME=$(date +%s)

check_monitoring_interval() {
    local current_time=$(date +%s)
    local elapsed=$((current_time - LAST_MONITOR_TIME))
    
    # If more than 10 minutes without a check
    if [ $elapsed -gt 600 ]; then
        echo "⚠️ R008 WARNING: ${elapsed}s since last monitor!"
        perform_monitoring_check
        LAST_MONITOR_TIME=$current_time
    fi
}
```

## Required Monitoring Actions

### 1. Agent Health Check
```bash
check_agent_health() {
    # For each active agent
    for agent in "${ACTIVE_AGENTS[@]}"; do
        echo "Checking $agent..."
        
        # Verify still working
        if ! agent_is_responsive "$agent"; then
            echo "❌ Agent $agent not responding!"
            handle_unresponsive_agent "$agent"
        fi
        
        # Check for errors
        if agent_has_errors "$agent"; then
            echo "⚠️ Agent $agent reported errors"
            review_agent_errors "$agent"
        fi
    done
}
```

### 2. Progress Verification
```bash
check_effort_progress() {
    local effort="$1"
    
    # Check line count
    local lines=$(line-counter.sh -c "$effort-branch" 2>/dev/null || echo 0)
    echo "📊 Current size: $lines/800 lines"
    
    if [ "$lines" -gt 700 ]; then
        echo "⚠️ WARNING: Approaching size limit!"
    fi
    
    # Check commits
    local commits=$(git rev-list --count HEAD...origin)
    if [ "$commits" -gt 5 ]; then
        echo "⚠️ WARNING: $commits unpushed commits"
    fi
}
```

### 3. State Consistency
```bash
verify_state_consistency() {
    # Load state file
    local current_state=$(grep "current_state:" orchestrator-state.json | awk '{print $2}')
    local efforts_in_progress=$(grep -c "in_progress" orchestrator-state.json)
    
    echo "📍 Current state: $current_state"
    echo "🔄 Active efforts: $efforts_in_progress"
    
    # Verify consistency
    if [[ "$current_state" == "IDLE" ]] && [ "$efforts_in_progress" -gt 0 ]; then
        echo "❌ INCONSISTENCY: IDLE but efforts active!"
        reconcile_state
    fi
}
```

## Monitoring Report Format

```yaml
monitoring_report:
  timestamp: "2025-08-26T14:00:00Z"
  message_count: 25
  checks_performed: 5
  
  agent_status:
    sw_engineer_1:
      status: "active"
      current_task: "implementing auth module"
      lines_written: 450
      time_elapsed: "1h 30m"
    
    code_reviewer_1:
      status: "waiting"
      waiting_for: "sw_engineer_1 completion"
      time_waiting: "10m"
  
  metrics:
    total_lines: 1250
    efforts_complete: 2
    efforts_in_progress: 1
    efforts_pending: 3
    
  issues_found:
    - type: "size_warning"
      effort: "E2.1.3"
      current: 720
      limit: 800
    
  actions_taken:
    - "Warned SW Engineer about size"
    - "Saved TODO state"
    - "Updated orchestrator-state.json"
```

## Common Violations

### VIOLATION: Skipping Checks
```bash
# ❌ WRONG: 10 messages without monitoring
Message 1: Spawn agent...
Message 2: Continue...
...
Message 10: Still no monitoring check
```

### VIOLATION: Superficial Checks
```bash
# ❌ WRONG: Check without verification
echo "Monitoring check done"  # No actual checks
```

### VIOLATION: Ignoring Issues
```bash
# ❌ WRONG: Finding issues but not acting
echo "Size violation detected: 850 lines"
# Continues without addressing
```

## Correct Patterns

### GOOD: Regular Monitoring
```bash
# Every 5 messages
Message 5: "📊 Monitoring check #1..."
Message 10: "📊 Monitoring check #2..."
Message 15: "📊 Monitoring check #3..."
```

### GOOD: Responsive Action
```bash
# Finding and addressing issues
echo "📊 Monitoring check..."
echo "⚠️ Size warning: 750 lines"
echo "🎬 ACTION: Preparing split plan"
echo "🚀 Spawning Code Reviewer for split"
```

### GOOD: Detailed Reporting
```bash
echo "📊 MONITORING REPORT"
echo "Agents: 2 active, 1 waiting"
echo "Progress: 450/800 lines"
echo "Issues: None"
echo "Next check: Message #30"
```

## Grading Impact

```python
def calculate_monitoring_grade(session):
    expected_checks = session.message_count // 5
    actual_checks = session.monitoring_checks_performed
    
    compliance_rate = actual_checks / expected_checks
    
    # Base grade
    if compliance_rate >= 0.95:
        grade = 100
    elif compliance_rate >= 0.80:
        grade = 90
    elif compliance_rate >= 0.60:
        grade = 75
    else:
        grade = 50
    
    # Penalties
    if session.missed_critical_issue:
        grade -= 20
    
    if session.longest_gap > 10:  # messages
        grade -= 10
    
    return grade
```

## Self-Monitoring

```bash
# Orchestrator self-check
monitor_myself() {
    echo "🪞 Self-monitoring check..."
    
    if [ $((MESSAGE_COUNT % 5)) -eq 0 ]; then
        echo "✅ R008: Monitoring on schedule"
    else
        local till_next=$((5 - (MESSAGE_COUNT % 5)))
        echo "📅 Next monitor in $till_next messages"
    fi
}
```

---
**Remember:** Monitoring prevents disasters. Check early, check often, act immediately.