# Rule R197: One Agent Per Effort - No Reuse Policy

## Rule Statement
The orchestrator MUST spawn a NEW agent instance for EVERY effort. Agents MUST NOT be reused across efforts. Each agent gets exactly ONE effort directory, completes its work, and terminates. No agent may work on multiple efforts.

## Criticality Level
**BLOCKING** - Agent reuse causes context pollution and state confusion

## Enforcement Mechanism
- **Technical**: Agent ID tracking and effort assignment validation
- **Behavioral**: Agents refuse work outside their assigned effort
- **Grading**: -40% for agent reuse, -30% for multi-effort agents

## Core Principles

### 1. One-to-One Mapping
```
1 Effort = 1 Agent Instance = 1 Working Directory = 1 Task
```

### 2. Agent Lifecycle
```
SPAWN → ASSIGN_EFFORT → WORK → COMPLETE → TERMINATE
         ↓
    NEVER REASSIGN
```

### 3. No Context Carry-Over
- Each agent starts fresh with no knowledge of other efforts
- No shared state between agent instances
- Clean context for every effort

## Detailed Requirements

### ORCHESTRATOR: Spawning Protocol

```bash
# ❌❌❌ WRONG - Reusing agents across efforts
spawn_agent_once() {
    Task: software-engineer
    echo "Work on effort 1 in efforts/phase1/wave1/api-types"
    echo "Now work on effort 2 in efforts/phase1/wave1/controllers"  # WRONG!
    echo "Now work on effort 3 in efforts/phase1/wave1/webhooks"    # WRONG!
}

# ✅✅✅ CORRECT - New agent for each effort
spawn_agents_correctly() {
    # Effort 1: API Types
    prepare_effort_for_agent 1 1 "api-types"
    Task: software-engineer
    Working directory: efforts/phase1/wave1/api-types
    Your ONLY task: Implement API types per plan
    When complete, you terminate. You work on NO other efforts.
    
    # Effort 2: Controllers (NEW AGENT)
    prepare_effort_for_agent 1 1 "controllers"
    Task: software-engineer  # NEW INSTANCE
    Working directory: efforts/phase1/wave1/controllers
    Your ONLY task: Implement controllers per plan
    When complete, you terminate. You work on NO other efforts.
    
    # Effort 3: Webhooks (NEW AGENT)
    prepare_effort_for_agent 1 1 "webhooks"
    Task: software-engineer  # NEW INSTANCE
    Working directory: efforts/phase1/wave1/webhooks
    Your ONLY task: Implement webhooks per plan
    When complete, you terminate. You work on NO other efforts.
}
```

### Parallel Spawning (When Dependencies Allow)

```bash
# Check implementation plan for parallelizable efforts
identify_parallel_efforts() {
    # From implementation plan:
    # - api-types: No dependencies (can start immediately)
    # - controllers: No dependencies (can start immediately)  
    # - webhooks: Depends on api-types (must wait)
    
    # Spawn parallel efforts (SEPARATE AGENTS)
    echo "Spawning parallel agents for independent efforts..."
    
    # Agent 1
    prepare_effort_for_agent 1 1 "api-types"
    Task: software-engineer
    Working directory: efforts/phase1/wave1/api-types
    You work ONLY on api-types. One effort only.
    
    # Agent 2 (DIFFERENT INSTANCE, PARALLEL)
    prepare_effort_for_agent 1 1 "controllers"
    Task: software-engineer
    Working directory: efforts/phase1/wave1/controllers
    You work ONLY on controllers. One effort only.
    
    # Wait for api-types completion, then:
    # Agent 3 (NEW INSTANCE)
    prepare_effort_for_agent 1 1 "webhooks"
    Task: software-engineer
    Working directory: efforts/phase1/wave1/webhooks
    You work ONLY on webhooks. One effort only.
}
```

### AGENT: Single Effort Verification

```bash
# SW Engineer MUST verify single-effort assignment
verify_single_effort_only() {
    # Get current working directory
    EFFORT_DIR=$(pwd)
    
    # Extract effort name from path
    EFFORT_NAME=$(echo "$EFFORT_DIR" | grep -oP 'efforts/phase\d+/wave\d+/\K[^/]+')
    
    echo "═══════════════════════════════════════════════════════"
    echo "SINGLE EFFORT ASSIGNMENT CONFIRMATION"
    echo "═══════════════════════════════════════════════════════"
    echo "I am assigned to: $EFFORT_NAME"
    echo "I will work ONLY in: $EFFORT_DIR"
    echo "I will NOT work on any other efforts"
    echo "When complete, I will terminate"
    echo "═══════════════════════════════════════════════════════"
    
    # Lock to single directory
    export LOCKED_EFFORT_DIR="$EFFORT_DIR"
}

# Prevent directory changes to other efforts
prevent_effort_switching() {
    CURRENT_DIR=$(pwd)
    if [[ "$CURRENT_DIR" != "$LOCKED_EFFORT_DIR"* ]]; then
        echo "❌ FATAL: Attempted to leave assigned effort!"
        echo "Assigned: $LOCKED_EFFORT_DIR"
        echo "Attempted: $CURRENT_DIR"
        exit 1
    fi
}
```

## State Management

### Orchestrator State Tracking
```yaml
# orchestrator-state.yaml
efforts_in_progress:
  phase1_wave1_api_types:
    agent_id: "sw-eng-instance-001"  # Unique instance ID
    status: "in_progress"
    started: "2024-01-20T10:00:00Z"
    
  phase1_wave1_controllers:
    agent_id: "sw-eng-instance-002"  # Different instance
    status: "in_progress"
    started: "2024-01-20T10:00:00Z"

efforts_completed:
  phase1_wave1_api_types:
    agent_id: "sw-eng-instance-001"  # Same agent that started it
    completed: "2024-01-20T11:00:00Z"
    terminated: true  # Agent terminated after completion
```

### Agent Termination Protocol
```bash
# Agent completion and termination
complete_effort_and_terminate() {
    echo "Effort complete: $EFFORT_NAME"
    
    # Final commits
    git add -A
    git commit -m "feat(${EFFORT_NAME}): Complete implementation"
    git push
    
    # Update completion marker
    echo "$(date): Completed by $(whoami)" > .effort-complete
    
    echo "═══════════════════════════════════════════════════════"
    echo "AGENT TERMINATING"
    echo "Effort: $EFFORT_NAME - COMPLETE"
    echo "This agent instance will not be reused"
    echo "═══════════════════════════════════════════════════════"
    
    # Agent terminates here - orchestrator must spawn new one for next effort
}
```

## Common Violations to Avoid

### ❌ Sequential Reuse
```bash
# WRONG - Same agent doing multiple efforts in sequence
Task: software-engineer
for effort in api-types controllers webhooks; do
    cd efforts/phase1/wave1/$effort
    implement_effort $effort
done
```

### ❌ Context Carry-Over
```bash
# WRONG - Agent remembering previous effort
"Now that I finished api-types, let me work on controllers"
"Based on what I did in the previous effort..."
```

### ❌ Multi-Directory Work
```bash
# WRONG - Agent working across directories
cd ../api-types && cp types.go ../controllers/
cd ../webhooks && implement_based_on_controllers
```

## Correct Examples

### ✅ Proper Sequential Spawning
```bash
# Effort 1
prepare_effort_for_agent 1 1 "api-types"
Task: software-engineer  # Instance 1
# ... agent completes and terminates

# Effort 2 (after effort 1 completes)
prepare_effort_for_agent 1 1 "controllers"
Task: software-engineer  # Instance 2 (NEW)
# ... agent completes and terminates
```

### ✅ Proper Parallel Spawning
```bash
# Spawn multiple agents AT ONCE for parallel efforts
Task: software-engineer  # Instance 1
Working directory: efforts/phase1/wave1/api-types

Task: software-engineer  # Instance 2
Working directory: efforts/phase1/wave1/controllers

Task: software-engineer  # Instance 3
Working directory: efforts/phase1/wave1/validators
```

## Integration with Other Rules

- **R196**: Orchestrator creates workspace before spawning each new agent
- **R193**: Each agent gets its own effort clone
- **R052**: Agent spawning protocol for each instance
- **R053**: Parallelization of independent efforts with separate agents

## Grading Impact

- **Agent reuse across efforts**: -40% (Major violation)
- **Multi-effort agent**: -30% (Context pollution)
- **No effort isolation**: -25% (State leakage)
- **Missing termination**: -15% (Lifecycle violation)
- **Context carry-over**: -20% (Independence violation)

## Benefits of This Approach

1. **Clean Context**: Each effort gets fresh agent with no baggage
2. **Parallel Execution**: Independent efforts can run simultaneously
3. **Clear Boundaries**: No confusion about what agent owns what
4. **State Isolation**: No accidental dependencies between efforts
5. **Debugging Clarity**: Issues isolated to single agent/effort pair