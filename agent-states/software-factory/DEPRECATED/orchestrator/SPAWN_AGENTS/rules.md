# SPAWN_SW_ENGINEERS State Rules

## 🔴🔴🔴 ABSOLUTE REQUIREMENT: STATE MANAGER CONSULTATION 🔴🔴🔴

**SUPREME LAW - R517 - UNIVERSAL STATE MANAGER CONSULTATION LAW**

**BEFORE EXITING THIS STATE, YOU MUST:**

1. **SPAWN STATE MANAGER FOR SHUTDOWN_CONSULTATION** (MANDATORY - NO EXCEPTIONS)
2. **NEVER UPDATE STATE FILES DIRECTLY** (orchestrator-state-v3.json, bug-tracking.json, etc.)
3. **NEVER COMMIT STATE FILES YOURSELF** (State Manager does this atomically)

**FAILURE TO CONSULT STATE MANAGER = IMMEDIATE SYSTEM HALT (-100% GRADE)**

### Enforcement Mechanism

If you attempt to exit this state without spawning State Manager:
- ❌ Pre-commit hooks will REJECT your commit
- ❌ Validation tools will FAIL the build
- ❌ Grading system will assign -100% penalty
- ❌ System will transition to ERROR_RECOVERY

### Required Pattern (COPY THIS EXACTLY)

```bash
# At end of state work, BEFORE any state file updates:

echo "🔴 MANDATORY: Spawning State Manager for SHUTDOWN_CONSULTATION"

# Spawn State Manager (REQUIRED - NOT OPTIONAL)
# Task: state-manager
# State: SHUTDOWN_CONSULTATION
# Current State: [YOUR_CURRENT_STATE]
# Proposed Next State: [YOUR_PROPOSED_NEXT_STATE]
# Work Summary: [SUMMARY_OF_WORK_COMPLETED]

# State Manager will:
# 1. Validate proposed transition against state machine
# 2. Update all 4 state files atomically
# 3. Commit with [R288] tag
# 4. Return REQUIRED next state (may differ from proposal)

# Wait for State Manager response
# Follow State Manager's directive (REQUIRED next state)
# DO NOT proceed until State Manager confirms
```

**YOU MUST NEVER:**
- ❌ Update orchestrator-state-v3.json yourself
- ❌ Update bug-tracking.json yourself
- ❌ Update integration-containers.json yourself
- ❌ Use `jq` to modify state files
- ❌ Use `sed/awk` to modify state files
- ❌ Set `validated_by: "orchestrator"` (must be "state-manager")
- ❌ Skip State Manager consultation "just this once"
- ❌ Think "I'll validate it manually"

**ONLY State Manager may update state files. This is NON-NEGOTIABLE.**

See: `rule-library/R517-universal-state-manager-consultation-law.md`
## State Purpose
Spawn appropriate agents (SW Engineers, Code Reviewers, Architects) based on current work requirements, ensuring proper parallelization and workspace isolation.

## Entry Conditions
- From ANALYZE_IMPLEMENTATION_PARALLELIZATION when ready to spawn SW Engineers
- From decision states when specific agent types needed
- From any state requiring agent deployment

## Mandatory Validations

### 1. Determine Agent Requirements
```bash
# Get spawn context from state file
SPAWN_TYPE=$(jq -r '.spawn_context.type' orchestrator-state-v3.json)
SPAWN_COUNT=$(jq -r '.spawn_context.count' orchestrator-state-v3.json)
SPAWN_PURPOSE=$(jq -r '.spawn_context.purpose' orchestrator-state-v3.json)

echo "Agent Spawn Request:"
echo "  Type: ${SPAWN_TYPE}"
echo "  Count: ${SPAWN_COUNT}"
echo "  Purpose: ${SPAWN_PURPOSE}"

# Validate spawn context
if [[ "$SPAWN_TYPE" == "null" ]] || [[ -z "$SPAWN_TYPE" ]]; then
    echo "ERROR: No spawn type specified"
    update_state "ERROR_RECOVERY" "Missing spawn context"
    exit 1
fi
```

### 2. Validate Parallelization Rules (R151)
```bash
# Check parallelization constraints
validate_parallelization() {
    local agent_type="$1"
    local count="$2"

    # R151: Same agent type can be parallelized
    if [[ "$count" -gt 1 ]]; then
        echo "Parallelization check for ${count} ${agent_type} agents..."

        # Check if work can be parallelized
        case "$agent_type" in
            "sw-engineer")
                # Check if efforts are independent
                INDEPENDENT=$(jq -r '.spawn_context.independent' orchestrator-state-v3.json)
                if [[ "$INDEPENDENT" != "true" ]]; then
                    echo "WARNING: Efforts not independent - using sequential deployment"
                    return 1
                fi
                ;;

            "code-reviewer")
                # Code reviewers can parallelize different efforts
                echo "Code reviewers can work in parallel on different efforts"
                ;;

            "architect")
                # Typically single architect per phase
                echo "WARNING: Multiple architects unusual - verify requirement"
                ;;
        esac
    fi

    echo "✅ Parallelization validated for ${count} ${agent_type} agents"
    return 0
}

validate_parallelization "$SPAWN_TYPE" "$SPAWN_COUNT"
ALLOW_PARALLEL=$?
```

### 3. Prepare Workspace Isolation (R356)
```bash
# Each agent needs isolated workspace
prepare_workspaces() {
    local agent_type="$1"
    local count="$2"

    echo "Preparing ${count} isolated workspaces..."

    for i in $(seq 1 "$count"); do
        # Define workspace for this agent
        case "$SPAWN_PURPOSE" in
            *"effort"*)
                EFFORT_ID=$(jq -r ".spawn_context.efforts[$((i-1))].id" orchestrator-state-v3.json)
                WORKSPACE="efforts/${EFFORT_ID}"
                ;;
            *"review"*)
                REVIEW_ID="review-$(date +%Y%m%d-%H%M%S)-${i}"
                WORKSPACE="reviews/${REVIEW_ID}"
                ;;
            *)
                WORKSPACE="${agent_type}-${i}-$(date +%Y%m%d-%H%M%S)"
                ;;
        esac

        echo "  Agent ${i}: ${WORKSPACE}"

        # Record workspace assignment
        jq --arg idx "$((i-1))" \
           --arg ws "$WORKSPACE" \
           '.spawn_context.workspaces[$idx | tonumber] = $ws' \
           orchestrator-state-v3.json > /tmp/state.json
        mv /tmp/state.json orchestrator-state-v3.json
    done
}

prepare_workspaces "$SPAWN_TYPE" "$SPAWN_COUNT"
```

## State Actions

### 1. Create Agent Spawn Commands
```bash
# Generate spawn commands for agents
create_spawn_commands() {
    local agent_type="$1"
    local count="$2"
    local parallel="$3"

    echo "Generating spawn commands..."

    # Create spawn script
    cat > /tmp/spawn_agents.sh <<'SPAWN_SCRIPT'
#!/bin/bash
# Agent Spawn Script
# Generated: $(date '+%Y-%m-%d %H:%M:%S %Z')

# R151: Record spawn timestamp for parallelization verification
SPAWN_TIMESTAMP=$(date +%s)
echo "SPAWN INITIATED: ${SPAWN_TIMESTAMP}"

SPAWN_SCRIPT

    for i in $(seq 1 "$count"); do
        WORKSPACE=$(jq -r ".spawn_context.workspaces[$((i-1))]" orchestrator-state-v3.json)
        EFFORT_ID=$(jq -r ".spawn_context.efforts[$((i-1))].id" orchestrator-state-v3.json)

        cat >> /tmp/spawn_agents.sh <<SPAWN_CMD

# Agent ${i} of ${count}
echo "Spawning ${agent_type} #${i}..."
echo "  Workspace: ${WORKSPACE}"
echo "  Effort: ${EFFORT_ID}"
echo "  Timestamp: \$(date '+%Y-%m-%d %H:%M:%S %Z')"

# Spawn command would go here
# /spawn-${agent_type} --workspace "${WORKSPACE}" --effort "${EFFORT_ID}"

SPAWN_CMD

        if [[ "$parallel" != "true" ]]; then
            echo "echo 'Waiting for agent ${i} to complete...'" >> /tmp/spawn_agents.sh
            echo "# Wait for sequential execution" >> /tmp/spawn_agents.sh
        fi
    done

    cat >> /tmp/spawn_agents.sh <<'SPAWN_END'

# Verify timing for parallel spawns (R151)
if [[ "$count" -gt 1 ]] && [[ "$parallel" == "true" ]]; then
    END_TIMESTAMP=$(date +%s)
    SPAWN_DURATION=$((END_TIMESTAMP - SPAWN_TIMESTAMP))

    if [[ $SPAWN_DURATION -gt 5 ]]; then
        echo "⚠️ WARNING: Parallel spawn took ${SPAWN_DURATION}s (>5s limit per R151)"
    else
        echo "✅ Parallel spawn completed in ${SPAWN_DURATION}s"
    fi
fi

echo "All agents spawned successfully"
SPAWN_END

    chmod +x /tmp/spawn_agents.sh
    echo "Spawn script created: /tmp/spawn_agents.sh"
}

create_spawn_commands "$SPAWN_TYPE" "$SPAWN_COUNT" "$ALLOW_PARALLEL"
```

### 2. Record Agent Deployments
```bash
# Update state with deployed agents
for i in $(seq 1 "$SPAWN_COUNT"); do
    WORKSPACE=$(jq -r ".spawn_context.workspaces[$((i-1))]" orchestrator-state-v3.json)
    EFFORT_ID=$(jq -r ".spawn_context.efforts[$((i-1))].id" orchestrator-state-v3.json)

    # Add to active agents list
    jq --arg type "$SPAWN_TYPE" \
       --arg ws "$WORKSPACE" \
       --arg effort "$EFFORT_ID" \
       --arg ts "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
       '.active_agents += [{
           type: $type,
           workspace: $ws,
           effort: $effort,
           spawned_at: $ts,
           status: "spawned"
       }]' orchestrator-state-v3.json > /tmp/state.json

    mv /tmp/state.json orchestrator-state-v3.json
done

echo "Recorded ${SPAWN_COUNT} agent deployments"
```

### 3. Determine Next State
```bash
# Based on spawn purpose, determine monitoring state
case "$SPAWN_PURPOSE" in
    *"implementation"*)
        NEXT_STATE="MONITORING_SWE_PROGRESS"
        ;;
    *"review"*)
        NEXT_STATE="MONITORING_EFFORT_REVIEWS"
        ;;
    *"fix"*)
        NEXT_STATE="MONITORING_EFFORT_FIXES"
        ;;
    *"planning"*)
        # Planning agents report back differently
        case "$SPAWN_TYPE" in
            "architect")
                NEXT_STATE="WAITING_FOR_PHASE_PLANS"
                ;;
            "code-reviewer")
                NEXT_STATE="WAITING_FOR_EFFORT_PLANS"
                ;;
            *)
                NEXT_STATE="MONITORING_PLANNING"
                ;;
        esac
        ;;
    *)
        NEXT_STATE="MONITORING_SWE_PROGRESS"
        ;;
esac

echo "Next state after spawn: ${NEXT_STATE}"
```

## Exit Conditions

### Success Criteria
- Agent type and count determined
- Parallelization validated per R151
- Workspaces prepared
- Spawn commands generated
- Agents recorded in state

### State Transitions
- **MONITORING_SWE_PROGRESS**: After spawning SW Engineers
- **MONITORING_EFFORT_REVIEWS**: After spawning Code Reviewers
- **MONITORING_EFFORT_FIXES**: After spawning for fixes
- **WAITING_FOR_[X]_PLANS**: After spawning planners
- **ERROR_RECOVERY**: If spawn fails

### State Update Requirements
```bash
# Update state to monitoring
update_state() {
    local next_state="$1"
    local notes="${2:-Agent spawn complete}"

    jq --arg state "$next_state" \
       --arg notes "$notes" \
       --arg timestamp "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
       '.state_machine.current_state = $state |
        .spawn_context.completed = true |
        .last_transition = {
            from: "SPAWN_SW_ENGINEERS",
            to: $state,
            timestamp: $timestamp,
            notes: $notes
        }' orchestrator-state-v3.json > /tmp/state.json

    mv /tmp/state.json orchestrator-state-v3.json
    echo "State updated to: ${next_state}"
}

# Transition to monitoring
update_state "${NEXT_STATE}" "Spawned ${SPAWN_COUNT} ${SPAWN_TYPE} agents"
```

## Associated Rules
- **R290**: State rule reading verification (SUPREME LAW)
- **R151**: Parallelization timing requirements (<5s)
- **R356**: Workspace isolation requirements
- **R233**: Single operation per state (SUPREME LAW)
- **R322**: Agent deployment protocols

## Prohibitions
- ❌ Spawn without workspace isolation
- ❌ Mix agent types in single spawn
- ❌ Exceed 5s timing for parallel spawns

## Exit Conditions and Continuation Flag

### ✅ CONTINUE-SOFTWARE-FACTORY=TRUE (Standard Operation - DEFAULT)

**Use TRUE when:**
- ✅ All required agents spawned successfully
- ✅ Workspace isolation verified for each agent
- ✅ Parallel spawns completed within 5s (R151)
- ✅ State file updated with spawn tracking
- ✅ Following designed workflow
- ✅ Ready to transition to monitoring state

**THIS IS NORMAL WORKFLOW.** Spawning agents for implementation, review, or planning
is the DESIGNED PROCESS. This is automation working as intended. Transitioning to
monitoring states after spawning is EXPECTED BEHAVIOR.

### 🔴 CONTINUE-SOFTWARE-FACTORY=FALSE (Exceptional Conditions ONLY)

**Use FALSE ONLY when:**
- ❌ Cannot create working copies for agents
- ❌ Workspace isolation failures (R356 violations)
- ❌ State machine corruption detected
- ❌ Required infrastructure files missing
- ❌ Cannot determine which agents to spawn
- ❌ Unrecoverable error prevents agent creation

**DO NOT set FALSE because:**
- ❌ Spawning agents (this is NORMAL workflow!)
- ❌ Transitioning to monitoring (EXPECTED process)
- ❌ R322 requires stop (stop ≠ FALSE flag!)
- ❌ "User might want to review" (only if exceptional)
- ❌ Multiple agents spawned (parallelization is NORMAL!)

### Critical Distinction: R322 Stop vs Continuation Flag

**R322 "stop"** = End conversation turn (`exit 0`)
- Purpose: Context preservation
- Always required after spawning agents

**Continuation flag** = Can system auto-restart?
- TRUE = Normal operations (default)
- FALSE = Exceptional/error conditions (rare)

**Correct pattern for this state:**
```bash
# Spawn all agents
# Update state file
exit 0  # R322 stop
```

**Last line before exit:**
```bash
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # Normal operation
```

**Grading Impact:**
- Using FALSE for normal spawning: -20% per violation
- Breaking automation flow: -30%

## R313 Enforcement - MANDATORY STOP (Context Preservation)

```bash
# This is the ABSOLUTE LAST thing that happens in this state
echo ""
echo "🛑 R313 ENFORCEMENT: STOPPING INFERENCE NOW (to preserve context)"
echo "The orchestrator MUST stop inference to prevent context overflow."
echo "SW Engineers have been spawned for implementation."
echo ""
echo "⚠️ IMPORTANT: This is a NORMAL stop for context preservation, not an error!"
echo "Next state will be: MONITORING_SWE_PROGRESS"
echo "System will automatically continue when ready."
```

## Automation Flag

```bash
# After successful spawn and state transition:
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # TRUE because spawning agents is NORMAL operation!
# The system stops inference but sets TRUE to allow automatic restart
```
- ❌ Spawn without recording in state
- ❌ Continue without monitoring spawned agents

## Notes
- Spawn timing is critical for R151 compliance
- Each agent must have isolated workspace
- Sequential spawning for dependent work
- Parallel spawning for independent work
- Always transition to appropriate monitoring state
## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

**MUST output as LAST line before state completion:**
- `CONTINUE-SOFTWARE-FACTORY=TRUE` - State completed, automation can continue
- `CONTINUE-SOFTWARE-FACTORY=FALSE` - Error/block, manual intervention required

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-CONTINUATION-FLAG-MASTER-GUIDE.md**
