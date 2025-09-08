# Rule R203: State-Aware Agent Startup and Rule Acknowledgment

## Rule Statement
Every agent MUST perform a state-aware startup sequence that includes: (1) loading core agent configuration, (2) determining current state from context, (3) loading state-specific rules, and (4) acknowledging BOTH core and state rules. Agents that fail to load and acknowledge state-specific rules must terminate immediately.

## Criticality Level
**BLOCKING** - Agents without proper state awareness waste context and cause failures

## Enforcement Mechanism
- **Technical**: Agents detect state from context files/instructions
- **Behavioral**: Agents must acknowledge both core and state rules
- **Grading**: -40% for not loading state-specific rules

## Core Principle

```
Agent Startup = Core Identity + State Detection + State Rules + Full Acknowledgment
Every agent loads ONLY what it needs for its current state
```

## Detailed Requirements

### MANDATORY STARTUP SEQUENCE

```bash
# ═══════════════════════════════════════════════════════════════
# UNIVERSAL AGENT STARTUP SEQUENCE (ALL AGENTS MUST FOLLOW)
# ═══════════════════════════════════════════════════════════════

agent_startup_sequence() {
    echo "═══════════════════════════════════════════════════════════════"
    echo "🚀 AGENT STARTUP SEQUENCE INITIATED"
    echo "TIMESTAMP: $(date '+%Y-%m-%d %H:%M:%S %Z')"
    echo "═══════════════════════════════════════════════════════════════"
    
    # STEP 1: IDENTIFY AGENT TYPE
    AGENT_TYPE="UNKNOWN"
    if [[ "$0" == *"orchestrator"* ]] || [[ "$PROMPT" == *"orchestrator"* ]]; then
        AGENT_TYPE="orchestrator"
    elif [[ "$0" == *"sw-engineer"* ]] || [[ "$PROMPT" == *"sw-engineer"* ]]; then
        AGENT_TYPE="sw-engineer"
    elif [[ "$0" == *"code-reviewer"* ]] || [[ "$PROMPT" == *"code-reviewer"* ]]; then
        AGENT_TYPE="code-reviewer"
    elif [[ "$0" == *"architect"* ]] || [[ "$PROMPT" == *"architect"* ]]; then
        AGENT_TYPE="architect"
    fi
    
    echo "AGENT TYPE IDENTIFIED: $AGENT_TYPE"
    
    # Get project root (Software Factory 2.0 root)
    local SF_ROOT="${SF_ROOT:-/workspaces/software-factory-2.0-template}"
    
    # STEP 2: LOAD CORE CONFIGURATION
    echo "═══════════════════════════════════════════════════════════════"
    echo "📚 LOADING CORE AGENT CONFIGURATION"
    echo "═══════════════════════════════════════════════════════════════"
    CORE_CONFIG="${SF_ROOT}/.claude/agents/${AGENT_TYPE}.md"
    echo "MUST READ WITH READ TOOL: $CORE_CONFIG"
    # 🚨 USE READ TOOL: ${SF_ROOT}/.claude/agents/${AGENT_TYPE}.md
    # DO NOT PROCEED WITHOUT ACTUALLY READING THIS FILE!
    
    # STEP 3: LOAD TODO PERSISTENCE RULES
    echo "═══════════════════════════════════════════════════════════════"
    echo "📋 LOADING TODO PERSISTENCE RULE (R287)"
    echo "═══════════════════════════════════════════════════════════════"
    echo "MUST READ WITH READ TOOL:"
    echo "  - ${SF_ROOT}/rule-library/R287-todo-persistence-comprehensive.md"
    # 🚨 USE READ TOOL FOR EACH TODO RULE FILE!
    # THESE ARE BLOCKING CRITICALITY - MUST BE READ!
    
    # STEP 4: DETERMINE CURRENT STATE
    echo "═══════════════════════════════════════════════════════════════"
    echo "🔍 DETERMINING CURRENT STATE"
    echo "═══════════════════════════════════════════════════════════════"
    CURRENT_STATE=$(determine_agent_state)
    echo "CURRENT STATE: $CURRENT_STATE"
    
    # STEP 5: LOAD STATE-SPECIFIC RULES
    echo "═══════════════════════════════════════════════════════════════"
    echo "📋 LOADING STATE-SPECIFIC RULES"
    echo "═══════════════════════════════════════════════════════════════"
    STATE_RULES_PATH="${SF_ROOT}/agent-states/${AGENT_TYPE}/${CURRENT_STATE}/rules.md"
    
    if [ -f "$STATE_RULES_PATH" ]; then
        echo "✅ Found state rules: $STATE_RULES_PATH"
        echo "MUST READ WITH READ TOOL: $STATE_RULES_PATH"
        # 🚨 USE READ TOOL: ${SF_ROOT}/agent-states/${AGENT_TYPE}/${CURRENT_STATE}/rules.md
        # DO NOT PROCEED WITHOUT ACTUALLY READING THIS FILE!
    else
        echo "❌ FATAL: State rules not found at $STATE_RULES_PATH"
        echo "Cannot proceed without state-specific rules!"
        exit 1
    fi
    
    # STEP 6: ACKNOWLEDGE ALL RULES
    acknowledge_all_rules
}
```

### STATE DETERMINATION FUNCTIONS

```bash
# SW ENGINEER State Detection
determine_sw_engineer_state() {
    if [ -f "SPLIT-INVENTORY.md" ]; then
        echo "SPLIT_WORK"
    elif [ -f "REVIEW-FEEDBACK.md" ]; then
        echo "FIX_ISSUES"
    elif [ -f "IMPLEMENTATION-PLAN.md" ]; then
        if [ -f "work-log.md" ] && grep -q "Implementation complete" work-log.md; then
            echo "TEST_IMPLEMENTATION"
        else
            echo "IMPLEMENTATION"
        fi
    else
        echo "INIT"
    fi
}

# CODE REVIEWER State Detection
determine_code_reviewer_state() {
    # Check instructions or context
    if grep -q "plan.*split" <<< "$INSTRUCTIONS"; then
        echo "SPLIT_PLANNING"
    elif grep -q "create.*plan" <<< "$INSTRUCTIONS"; then
        echo "PLANNING"
    elif grep -q "review.*code" <<< "$INSTRUCTIONS"; then
        echo "CODE_REVIEW"
    elif grep -q "validate" <<< "$INSTRUCTIONS"; then
        echo "VALIDATION"
    else
        echo "INIT"
    fi
}

# ORCHESTRATOR State Detection
determine_orchestrator_state() {
    if [ -f "orchestrator-state.json" ]; then
        # Extract current_state from YAML
        CURRENT_STATE=$(grep "current_state:" orchestrator-state.json | awk '{print $2}')
        echo "$CURRENT_STATE"
    else
        echo "INIT"
    fi
}

# ARCHITECT State Detection
determine_architect_state() {
    if grep -q "wave.*review" <<< "$INSTRUCTIONS"; then
        echo "WAVE_REVIEW"
    elif grep -q "phase.*review" <<< "$INSTRUCTIONS"; then
        echo "PHASE_REVIEW"
    elif grep -q "integration" <<< "$INSTRUCTIONS"; then
        echo "INTEGRATION_REVIEW"
    else
        echo "INIT"
    fi
}

# Master state determination
determine_agent_state() {
    case "$AGENT_TYPE" in
        "sw-engineer")
            determine_sw_engineer_state
            ;;
        "code-reviewer")
            determine_code_reviewer_state
            ;;
        "orchestrator")
            determine_orchestrator_state
            ;;
        "architect")
            determine_architect_state
            ;;
        *)
            echo "INIT"
            ;;
    esac
}
```

### MANDATORY ACKNOWLEDGMENT FORMAT

```bash
acknowledge_all_rules() {
    echo "═══════════════════════════════════════════════════════════════"
    echo "📝 RULE ACKNOWLEDGMENT"
    echo "═══════════════════════════════════════════════════════════════"
    echo "I am @agent-${AGENT_TYPE} in state ${CURRENT_STATE}"
    echo "═══════════════════════════════════════════════════════════════"
    
    echo "CORE RULES ACKNOWLEDGED:"
    echo "------------------------"
    # List core rules from agent definition
    echo "✅ R001: Pre-flight checks [BLOCKING]"
    echo "✅ R186: Compaction detection [BLOCKING]"
    echo "✅ R287: Comprehensive TODO persistence - Save/Commit/Recover [BLOCKING]"
    echo "✅ R203: State-aware startup [BLOCKING]"
    
    # Agent-specific core rules
    case "$AGENT_TYPE" in
        "orchestrator")
            echo "✅ R151: Parallel spawning <5s [CRITICAL]"
            echo "✅ Never write code [BLOCKING]"
            ;;
        "sw-engineer")
            echo "✅ R197: One agent per effort [BLOCKING]"
            echo "✅ R196: Never create clones [BLOCKING]"
            ;;
        "code-reviewer")
            echo "✅ R199: Single reviewer for splits [BLOCKING]"
            echo "✅ R200: Measure only changeset [BLOCKING]"
            ;;
        "architect")
            echo "✅ Architectural consistency [CRITICAL]"
            ;;
    esac
    
    echo ""
    echo "STATE-SPECIFIC RULES ACKNOWLEDGED (${CURRENT_STATE}):"
    echo "-----------------------------------------------------"
    # List rules from state-specific file
    acknowledge_state_rules
    
    echo "═══════════════════════════════════════════════════════════════"
    echo "✅ ALL RULES ACKNOWLEDGED - PROCEEDING WITH ${CURRENT_STATE}"
    echo "═══════════════════════════════════════════════════════════════"
}

# State-specific rule acknowledgment
acknowledge_state_rules() {
    case "$CURRENT_STATE" in
        "IMPLEMENTATION")
            echo "✅ R106: Implementation efficiency [CRITICAL]"
            echo "✅ R152: Code quality standards [CRITICAL]"
            echo "✅ R176: Workspace isolation [BLOCKING]"
            echo "✅ Check size every 200 lines [MANDATORY]"
            ;;
        "SPLIT_WORK")
            echo "✅ R202: Handle ALL splits sequentially [BLOCKING]"
            echo "✅ R007: Each split <800 lines [BLOCKING]"
            echo "✅ Complete each split before next [MANDATORY]"
            ;;
        "CODE_REVIEW")
            echo "✅ R007: Verify <800 lines [BLOCKING]"
            echo "✅ R153: Review effectiveness [CRITICAL]"
            echo "✅ Test coverage validation [MANDATORY]"
            ;;
        "SPAWN_AGENTS")
            echo "✅ R151: Parallel spawn <5s [CRITICAL]"
            echo "✅ R197: One agent per effort [BLOCKING]"
            echo "✅ Include state in spawn [MANDATORY]"
            ;;
        *)
            echo "✅ State-specific rules loaded from $STATE_RULES_PATH"
            ;;
    esac
}
```

### ORCHESTRATOR SPAWN PATTERN WITH STATE

```bash
# ✅ CORRECT - Spawn with state awareness
spawn_agent_with_state() {
    local agent_type=$1
    local working_dir=$2
    local state=$3
    
    echo "Spawning $agent_type in state $state"
    
    Task: @agent-$agent_type
    Working directory: $working_dir
    
    CRITICAL STARTUP INSTRUCTIONS:
    1. You are starting in state: $state
    2. READ: .claude/agents/$agent_type.md (core config)
    3. READ: agent-states/$agent_type/$state/rules.md (state rules)
    4. Acknowledge BOTH core and state rules
    5. Proceed with $state operations
}

# Example usage
spawn_agent_with_state "software-engineer" "efforts/phase1/wave1/api-types" "IMPLEMENTATION"
spawn_agent_with_state "code-reviewer" "efforts/phase1/wave1/api-types" "PLANNING"
```

### VERIFICATION CHECKLIST

```yaml
startup_verification:
  core_config_loaded:
    - file: ".claude/agents/{agent_type}.md"
    - status: "MUST be loaded"
    - acknowledgment: "REQUIRED"
  
  state_determined:
    - method: "Context analysis"
    - result: "Valid state name"
    - fallback: "INIT if uncertain"
  
  state_rules_loaded:
    - file: "agent-states/{agent_type}/{state}/rules.md"
    - status: "MUST be loaded"
    - acknowledgment: "REQUIRED"
  
  acknowledgment_complete:
    - core_rules: "Listed and acknowledged"
    - state_rules: "Listed and acknowledged"
    - format: "Proper acknowledgment block"
```

## State Transition Handling

```bash
# When transitioning to a new state
handle_state_transition() {
    OLD_STATE="$CURRENT_STATE"
    NEW_STATE="$1"
    
    echo "═══════════════════════════════════════════════════════════════"
    echo "🔄 STATE TRANSITION: $OLD_STATE → $NEW_STATE"
    echo "═══════════════════════════════════════════════════════════════"
    
    # Save checkpoint from old state
    if [ -f "agent-states/$AGENT_TYPE/$OLD_STATE/checkpoint.md" ]; then
        echo "📝 Saving checkpoint per $OLD_STATE requirements"
        # READ: agent-states/$AGENT_TYPE/$OLD_STATE/checkpoint.md
        # Save required checkpoint data
    fi
    
    # Load new state rules
    CURRENT_STATE="$NEW_STATE"
    STATE_RULES_PATH="agent-states/$AGENT_TYPE/$NEW_STATE/rules.md"
    
    echo "📚 Loading new state rules: $STATE_RULES_PATH"
    # READ: $STATE_RULES_PATH
    
    # Acknowledge new state rules
    echo "═══════════════════════════════════════════════════════════════"
    echo "📝 NEW STATE ACKNOWLEDGMENT"
    echo "Transitioned to: $NEW_STATE"
    acknowledge_state_rules
    echo "═══════════════════════════════════════════════════════════════"
}
```

## 📚 EXPLICIT FILE PATHS FOR READ TOOL

**Every agent MUST use these EXACT paths with the Read tool at startup:**

```yaml
orchestrator_startup:
  core_config: ${SF_ROOT}/.claude/agents/orchestrator.md
  state_file: ${SF_ROOT}/orchestrator-state.json
  state_rules_pattern: ${SF_ROOT}/agent-states/orchestrator/${STATE}/rules.md
  examples:
    - ${SF_ROOT}/agent-states/orchestrator/INIT/rules.md
    - ${SF_ROOT}/agent-states/orchestrator/SPAWN_AGENTS/rules.md
    - ${SF_ROOT}/agent-states/orchestrator/WAVE_COMPLETE/rules.md

sw_engineer_startup:
  core_config: ${SF_ROOT}/.claude/agents/kcp-go-lang-sr-sw-eng.md
  state_rules_pattern: ${SF_ROOT}/agent-states/sw-engineer/${STATE}/rules.md
  examples:
    - ${SF_ROOT}/agent-states/sw-engineer/IMPLEMENTATION/rules.md
    - ${SF_ROOT}/agent-states/sw-engineer/FIX_ISSUES/rules.md
    - ${SF_ROOT}/agent-states/sw-engineer/SPLIT_WORK/rules.md

code_reviewer_startup:
  core_config: ${SF_ROOT}/.claude/agents/kcp-kubernetes-code-reviewer.md
  state_rules_pattern: ${SF_ROOT}/agent-states/code-reviewer/${STATE}/rules.md
  examples:
    - ${SF_ROOT}/agent-states/code-reviewer/PLANNING/rules.md
    - ${SF_ROOT}/agent-states/code-reviewer/CODE_REVIEW/rules.md
    - ${SF_ROOT}/agent-states/code-reviewer/SPLIT_PLANNING/rules.md

architect_startup:
  core_config: ${SF_ROOT}/.claude/agents/kcp-architect-reviewer.md
  state_rules_pattern: ${SF_ROOT}/agent-states/architect/${STATE}/rules.md
  examples:
    - ${SF_ROOT}/agent-states/architect/PHASE_PLANNING/rules.md
    - ${SF_ROOT}/agent-states/architect/WAVE_REVIEW/rules.md
```

**Where SF_ROOT is typically: /workspaces/software-factory-2.0-template**

## Common Violations to Avoid

### ❌ Loading All Rules
```bash
# WRONG - Loading everything wastes context
READ: .claude/agents/sw-engineer.md  # 40KB
READ: All state files  # Another 100KB
# Total: 140KB loaded, but only need 15KB!
```

### ❌ Not Loading State Rules
```bash
# WRONG - Missing state-specific guidance
READ: .claude/agents/sw-engineer.md
# Proceed without state rules = FAILURE
```

### ❌ No Acknowledgment
```bash
# WRONG - No proof rules were understood
# Just start working without acknowledgment
```

## Integration with Other Rules

- **R001**: Pre-flight checks (part of core)
- **R186**: Compaction detection (part of core)
- **R151**: Parallel spawning (orchestrator core)
- **R197**: One agent per effort (sw-engineer core)
- **R199**: Single reviewer for splits (code-reviewer core)
- **State-specific rules**: Loaded per current state

## Grading Impact

- **No state rules loaded**: -40% (Context waste)
- **No acknowledgment**: -30% (Compliance failure)
- **Wrong state detected**: -25% (Incorrect behavior)
- **All rules loaded**: -20% (Inefficiency)

## Summary

**Remember**:
- Load core config ALWAYS
- Detect state from context
- Load ONLY current state rules
- Acknowledge BOTH core and state rules
- Transition loads new state rules
- This reduces context usage by ~70%