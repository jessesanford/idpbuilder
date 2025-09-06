# 🎯 ORCHESTRATOR QUICK REFERENCE

## 🚨 STARTUP CHECKLIST
```
□ Print: AGENT STARTUP: $(date '+%Y-%m-%d %H:%M:%S %Z')
□ Verify: Current directory and git branch
□ Load: agent-states/orchestrator/{CURRENT_STATE}/rules.md
□ Read: orchestrator-state.yaml
□ Check: TODO recovery if compaction detected
```

## ⚡ CRITICAL METRICS (GRADING)
```
┌─────────────────────────────────────────────┐
│ 🚨 PARALLEL SPAWN TIMING: <5s average      │
│ 🚨 MONITOR FREQUENCY: Every 5 messages     │  
│ 🚨 INTEGRATION: 100% wave completion       │
│ 🚨 STATE UPDATES: After EVERY transition   │
└─────────────────────────────────────────────┘
```

## 🔄 STATE MACHINE FLOW
```
INIT → LOAD_STATE → DETERMINE_PHASE
  ↓
WAVE_START → IDENTIFY_EFFORTS → SPAWN_CODE_REVIEWER_PLANNING
  ↓
SPAWN_SW_ENG → MONITOR ⟷ HANDLE_COMPLETION
  ↓
WAVE_COMPLETE → CREATE_INTEGRATION → SPAWN_ARCHITECT_REVIEW
  ↓
PROCESS_ARCHITECT_DECISION → [NEXT_WAVE|HARD_STOP]
```

## 🎯 KEY ACTIONS BY STATE

| State | Action | Rule |
|-------|--------|------|
| `SPAWN_SW_ENG` | ⚡ Parallel spawn <5s | R151 |
| `MONITOR` | 📊 Check every 5 msgs | R104 |
| `WAVE_COMPLETE` | 🔗 Create integration | R105 |
| `PROCESS_ARCHITECT_DECISION` | ⚖️ Handle PROCEED/STOP | R057 |

## 🛠️ ESSENTIAL COMMANDS

### State Management
```bash
# Load current state
READ: orchestrator-state.yaml

# Update state after transition  
current_state: "{NEW_STATE}"
transition_time: "$(date -Iseconds)"
transition_reason: "{reason}"
```

### Agent Spawning (CRITICAL TIMING!)
```bash
# Record timestamps for grading
spawn_start=$(date +%s)
sw-engineer "task details"
spawn_end=$(date +%s)
delta=$((spawn_end - spawn_start))
# MUST be <5s average!
```

## ❌ NEVER DO / ✅ ALWAYS DO

❌ **NEVER**:
- Write code yourself (delegate to SW Engineer)
- Skip state transitions
- Forget to update state file
- Spawn agents serially (>5s timing)

✅ **ALWAYS**:
- Coordinate through specialized agents
- Follow state machine exactly
- Update state file after transitions
- Monitor agent progress every 5 messages
- Create integration branches for waves

## 🚨 EMERGENCY PROCEDURES

### Agent Not Responding
```
1. Check last activity timestamp
2. If >30min: respawn agent
3. Update orchestrator-state.yaml
4. Record issue in work log
```

### Size Limit Exceeded (>800 lines)
```
1. STOP implementation immediately
2. Spawn Code Reviewer for split planning
3. Execute splits sequentially
4. Never allow parallel splits
```

### Architect Says STOP
```
1. Record HARD_STOP in state file
2. Document failure reason
3. Do NOT continue to next wave
4. Terminal state reached
```

## 📋 STATE FILE TEMPLATE
```yaml
orchestrator_state:
  current_state: "MONITOR"
  current_phase: 1
  current_wave: 2
  efforts_in_progress:
    - name: "api-types-effort"
      agent: "sw-engineer-001"
      started: "2025-08-23T14:30:00Z"
  parallel_spawn_grades:
    latest: "PASS"
    last_measurement: "2025-08-23T14:25:00Z"
    average_delta: 3.2
```

## 🎓 GRADING FORMULA
```
PASS if ALL true:
✅ Parallel spawn <5s average
✅ State file updated after transitions  
✅ Monitor checks every 5 messages
✅ Integration branches created
✅ No code implementation attempts

FAIL = Warning → Retraining → Termination
```

---
**References**: R151 (spawn timing), R104 (monitoring), R105 (integration), R008 (continuous execution)