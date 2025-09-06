# 🔴🔴🔴 RULE R313 - MANDATORY STOP AFTER SPAWNING AGENTS (SUPREME LAW)

**Criticality:** SUPREME LAW - CONTEXT PRESERVATION  
**Grading Impact:** -100% IMMEDIATE FAILURE for violation  
**Enforcement:** MANDATORY - Every spawn state must exit after spawning

## SUPREME LAW STATEMENT

**THE ORCHESTRATOR MUST STOP IMMEDIATELY AFTER SPAWNING ANY AGENTS TO PRESERVE CONTEXT AND PREVENT RULE FORGETTING.**

## 🚨🚨🚨 CRITICAL CONTEXT OVERFLOW PROBLEM 🚨🚨🚨

### THE FATAL FLAW WE'RE FIXING:
```markdown
❌ OLD BROKEN PATTERN (CAUSES CONTEXT OVERFLOW):
1. Orchestrator spawns agents
2. Orchestrator continues running
3. Agents complete and report back (thousands of lines)
4. Reports accumulate in orchestrator's context
5. Rules get pushed out of context window
6. Orchestrator forgets critical rules
7. Orchestrator becomes confused and ineffective
8. SOFTWARE FACTORY FAILS

✅ NEW REQUIRED PATTERN (PRESERVES CONTEXT):
1. Orchestrator spawns agents
2. Orchestrator records what was spawned in state file
3. Orchestrator STOPS IMMEDIATELY
4. Human restarts with fresh context
5. Orchestrator reads state file to understand situation
6. Orchestrator processes agent results with full rule awareness
7. SOFTWARE FACTORY SUCCEEDS
```

## 🔴🔴🔴 MANDATORY STOP STATES 🔴🔴🔴

### THE ORCHESTRATOR MUST STOP AFTER THESE STATES:

```markdown
SPAWN STATES REQUIRING IMMEDIATE STOP:
✅ SPAWN_AGENTS - After spawning SWE agents
✅ SPAWN_CODE_REVIEWERS_EFFORT_PLANNING - After spawning effort planners
✅ SPAWN_CODE_REVIEWERS_FOR_REVIEW - After spawning code reviewers
✅ SPAWN_ENGINEERS_FOR_FIXES - After spawning fix engineers
✅ SPAWN_INTEGRATION_AGENT - After spawning integration agent
✅ SPAWN_ARCHITECT_PHASE_PLANNING - After spawning architect for planning
✅ SPAWN_ARCHITECT_PHASE_ASSESSMENT - After spawning architect for assessment
✅ SPAWN_ARCHITECT_WAVE_PLANNING - After spawning architect for wave planning
✅ SPAWN_CODE_REVIEWER_FIX_PLAN - After spawning fix planner
✅ SPAWN_CODE_REVIEWER_PHASE_FIX_PLAN - After spawning phase fix planner
✅ SPAWN_CODE_REVIEWER_MERGE_PLAN - After spawning merge planner

ANY STATE THAT SPAWNS AN AGENT MUST STOP!
```

## 🛡️ CONTEXT PRESERVATION PROTOCOL

### BEFORE STOPPING (MANDATORY RECORDING):
```bash
# Record what was spawned in state file
update_state "agents_spawned" "[
  {agent: 'sw-engineer-1', effort: 'EFFORT_001', status: 'spawned'},
  {agent: 'sw-engineer-2', effort: 'EFFORT_002', status: 'spawned'}
]"

# Record continuation point
update_state "continuation_state" "MONITOR"
update_state "continuation_message" "Spawned 2 SWE agents for parallel efforts"

# Commit and push state
git add orchestrator-state.yaml
git commit -m "state: spawned agents, stopping per R313"
git push

# Create continuation command
echo "To continue after agents complete: claude --continue"
```

### THE STOP MESSAGE (REQUIRED FORMAT):
```bash
echo "🛑 STOPPING PER R313 - CONTEXT PRESERVATION"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Agents spawned:"
echo "  • sw-engineer-1 → EFFORT_001"
echo "  • sw-engineer-2 → EFFORT_002"
echo ""
echo "Next state: MONITOR"
echo "To continue: claude --continue"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
exit 0  # MANDATORY EXIT
```

## 🔴🔴🔴 WHY THIS IS CRITICAL 🔴🔴🔴

### CONTEXT WINDOW MATHEMATICS:
```
Typical context window: ~200,000 tokens
Agent response size: 5,000-20,000 tokens each
Multiple agents: 3-5 agents common

RESULT:
- 3 agents × 15,000 tokens = 45,000 tokens
- 5 agents × 20,000 tokens = 100,000 tokens
- This pushes out critical rules from context
- Orchestrator loses its instructions
- CATASTROPHIC FAILURE
```

### WHAT GETS LOST WHEN CONTEXT OVERFLOWS:
1. **Critical Rules** - R287 (TODO persistence), R304 (line counting), etc.
2. **State Machine Knowledge** - Forgets valid transitions
3. **Grading Criteria** - Forgets what it's being graded on
4. **Current Objectives** - Loses track of goals
5. **Agent Coordination** - Forgets who's doing what

## ✅ CORRECT IMPLEMENTATION PATTERN

### IN SPAWN STATES:
```bash
# SPAWN_AGENTS state implementation
spawn_agents_for_efforts() {
    echo "🚀 Spawning agents for parallel efforts..."
    
    # Spawn the agents
    for effort in "${efforts[@]}"; do
        spawn_sw_engineer "$effort"
        record_spawned_agent "$effort"
    done
    
    # UPDATE STATE FILE WITH SPAWN INFO
    update_state "agents_spawned" "$spawned_list"
    update_state "continuation_state" "MONITOR"
    
    # COMMIT STATE CHANGES
    git add orchestrator-state.yaml
    git commit -m "state: spawned ${#efforts[@]} agents, stopping per R313"
    git push
    
    # MANDATORY STOP WITH CLEAR MESSAGE
    echo "
🛑 STOPPING PER R313 - CONTEXT PRESERVATION
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Spawned ${#efforts[@]} SWE agents for parallel work
State saved to: orchestrator-state.yaml
Next state: MONITOR

To continue after agents complete:
  claude --continue

This stop preserves context and prevents rule loss.
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
"
    exit 0  # MANDATORY EXIT - NOT OPTIONAL!
}
```

### ON RESTART (CONTINUATION):
```bash
# When restarted with --continue
continue_after_spawn() {
    echo "📖 Loading state after agent spawn..."
    
    # Read what was spawned
    agents=$(yq '.agents_spawned' orchestrator-state.yaml)
    next_state=$(yq '.continuation_state' orchestrator-state.yaml)
    
    echo "✅ Found spawned agents:"
    echo "$agents"
    
    echo "🔄 Transitioning to $next_state..."
    transition_to_state "$next_state"
    
    # Now we have fresh context to process agent results!
}
```

## ❌ VIOLATIONS TO AVOID

### VIOLATION: Continuing After Spawn
```bash
# ❌ CATASTROPHIC VIOLATION - CAUSES CONTEXT OVERFLOW
spawn_agents_in_parallel
echo "Now monitoring agents..."  # NO! STOP HERE!
monitor_agent_progress  # VIOLATION - Will lose context!
```

### VIOLATION: Waiting for Agents in Same Session
```bash
# ❌ CATASTROPHIC VIOLATION - ACCUMULATES RESPONSES
spawn_agents
while agents_still_running; do
    check_status  # NO! This accumulates responses!
    sleep 30
done
# By now, context is destroyed
```

### VIOLATION: Processing Results Without Stop
```bash
# ❌ CATASTROPHIC VIOLATION - NO CONTEXT RESET
spawn_code_reviewers
wait_for_reviews
process_review_results  # Context overflow! Rules lost!
```

## 🚨 ENFORCEMENT MECHANISM

### Self-Check in Spawn States:
```bash
after_spawning_agents() {
    echo "🔴 R313 ENFORCEMENT CHECK"
    echo "Have I spawned agents? YES"
    echo "Must I stop now? YES - R313 MANDATES IT"
    echo "Recording spawn info and stopping..."
    
    # CANNOT PROCEED WITHOUT STOPPING
    save_state_and_exit
}
```

### Grading Enforcement:
```yaml
violations:
  continued_after_spawn: -100%  # IMMEDIATE FAILURE
  waited_in_same_session: -100%  # IMMEDIATE FAILURE
  processed_without_stop: -100%  # IMMEDIATE FAILURE
  
excellence:
  clean_stops_after_spawns: +10%
  proper_continuation_handling: +10%
  zero_context_overflows: MAXIMUM_GRADE
```

## 🔄 THE NEW WORKFLOW PATTERN

### SPAWN → STOP → RESTART → PROCESS
```
1. ORCHESTRATOR SPAWNS:
   - Records what was spawned
   - Saves state
   - STOPS with clear message

2. AGENTS WORK:
   - Complete their tasks
   - Report results

3. HUMAN RESTARTS:
   - Fresh context loaded
   - State file read
   - Knows what was spawned

4. ORCHESTRATOR PROCESSES:
   - Has full rule awareness
   - Can properly handle results
   - Continues to next phase
```

## 📋 INTEGRATION WITH OTHER RULES

### SUPERSEDES AND REPLACES:
- **R021** (Orchestrator Never Stops) - DEPRECATED for spawn states
- **R231** (Continuous Operation) - MODIFIED to exclude spawn states

### WORKS WITH:
- **R287** - TODO persistence before stopping
- **R288** - State file updates before stopping
- **R206** - State machine validation
- **R290** - State rule reading after restart

## 🔴🔴🔴 CRITICAL REMINDERS 🔴🔴🔴

1. **EVERY spawn state MUST stop**
2. **NO exceptions for "quick" agent tasks**
3. **NO waiting for agents in same session**
4. **ALWAYS record what was spawned**
5. **ALWAYS provide continuation instructions**
6. **NEVER process agent results without context reset**

## 📢 THE NEW ORCHESTRATOR MANTRAS

### Repeat These Constantly:
1. **"Spawn and stop. Context is precious."**
2. **"Agent responses overflow context. I must stop."**
3. **"Fresh context preserves rules."**
4. **"Stop after spawn is not failure, it's wisdom."**
5. **"The system requires stops to maintain integrity."**

## 🎯 IMPLEMENTATION CHECKLIST

When implementing spawn states:
- [ ] Spawn the required agents
- [ ] Record spawn details in state file
- [ ] Update continuation_state field
- [ ] Commit and push state changes
- [ ] Display clear stop message
- [ ] Include continuation command
- [ ] EXIT with code 0
- [ ] DO NOT continue in same session
- [ ] DO NOT wait for agents
- [ ] DO NOT process results yet

## SUMMARY

**R313 establishes a new critical pattern:**
- Spawning agents requires immediate stop
- Stops preserve context and prevent rule loss
- Context overflow is the enemy of orchestration
- Fresh contexts maintain Software Factory integrity
- This rule SUPERSEDES R021 and R231 for spawn states

**REMEMBER**: An orchestrator that continues after spawning will lose its rules and fail. ALWAYS STOP AFTER SPAWN!