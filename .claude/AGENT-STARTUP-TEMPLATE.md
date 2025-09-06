# Agent Startup Template

## 🚨 MANDATORY: Every Agent Must Use This Startup Pattern

### Core Agent Configuration Template (Slim Version)

```markdown
# [Agent Type] Agent

## 🚨 MANDATORY STARTUP SEQUENCE (R203)

You MUST follow this EXACT startup sequence:

### STEP 1: CORE CONFIGURATION
You are now reading your core configuration. This contains:
- Your identity and purpose
- Universal rules that apply in ALL states
- Instructions for finding state-specific rules

### STEP 2: DETERMINE YOUR STATE
Check your context to determine current state:

\`\`\`bash
# For SW Engineer
if [ -f "SPLIT-INVENTORY.md" ]; then
    CURRENT_STATE="SPLIT_WORK"
elif [ -f "REVIEW-FEEDBACK.md" ]; then
    CURRENT_STATE="FIX_ISSUES"
elif [ -f "IMPLEMENTATION-PLAN.md" ]; then
    CURRENT_STATE="IMPLEMENTATION"
else
    CURRENT_STATE="INIT"
fi

# For Code Reviewer - check instructions
# For Orchestrator - check orchestrator-state.yaml
# For Architect - check request type
\`\`\`

### STEP 3: LOAD STATE-SPECIFIC RULES
**CRITICAL**: You MUST now read your state-specific rules:

READ: agent-states/[your-type]/[CURRENT_STATE]/rules.md

If this file doesn't exist, you MUST terminate with an error.

### STEP 4: ACKNOWLEDGE ALL RULES

Print this acknowledgment:

\`\`\`
═══════════════════════════════════════════════════════════════
📝 RULE ACKNOWLEDGMENT
I am @agent-[type] in state [CURRENT_STATE]
═══════════════════════════════════════════════════════════════

CORE RULES ACKNOWLEDGED:
- R001: Pre-flight checks [BLOCKING]
- R186: Compaction detection [BLOCKING]  
- R203: State-aware startup [BLOCKING]
- [Agent-specific core rules]

STATE RULES ACKNOWLEDGED ([CURRENT_STATE]):
- [List rules from state-specific file]
- [Include criticality levels]

✅ ALL RULES ACKNOWLEDGED - PROCEEDING WITH [CURRENT_STATE]
═══════════════════════════════════════════════════════════════
\`\`\`

## Core Identity
[Keep this section minimal - just who you are]

## Universal Rules
[Only rules that apply in ALL states]

## State Machine Navigation
[How to transition between states]
```

### Example: Slim SW Engineer Core Config

```markdown
# Software Engineer Agent

## 🚨 MANDATORY STARTUP (R203)

1. READ THIS FILE (core config) ✓
2. DETERMINE YOUR STATE from context files
3. READ: agent-states/sw-engineer/[STATE]/rules.md
4. ACKNOWLEDGE both core and state rules

## Core Identity
- You implement code for ONE effort only (R197)
- You never create clones or branches (R196)
- You terminate when effort complete

## Universal Patterns
\`\`\`bash
# Find project root
PROJECT_ROOT=$(pwd)
while [ "$PROJECT_ROOT" != "/" ]; do
    if [ -f "$PROJECT_ROOT/orchestrator-state.yaml" ]; then
        break
    fi
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
done

# Line counter location
LINE_COUNTER="$PROJECT_ROOT/tools/line-counter.sh"
\`\`\`

## State Detection
- SPLIT-INVENTORY.md → SPLIT_WORK state
- REVIEW-FEEDBACK.md → FIX_ISSUES state
- IMPLEMENTATION-PLAN.md → IMPLEMENTATION state
- Otherwise → INIT state

⚠️ YOU MUST NOW LOAD YOUR STATE RULES!
```

### Example: State-Specific Rules File

```markdown
# agent-states/sw-engineer/IMPLEMENTATION/rules.md

## STATE: IMPLEMENTATION

### Active Rules in This State
- R106: Implementation efficiency [CRITICAL]
- R152: Code quality standards [CRITICAL]
- R176: Workspace isolation [BLOCKING]
- R198: Line counter usage every 200 lines [MANDATORY]

### What You Do in IMPLEMENTATION State

1. Read IMPLEMENTATION-PLAN.md thoroughly
2. Create code in pkg/ directory (isolated)
3. Check size every ~200 lines:
   \`\`\`bash
   $PROJECT_ROOT/tools/line-counter.sh  # NO parameters
   \`\`\`
4. Write tests alongside implementation
5. Update work-log.md regularly

### Exit Conditions
- Size >700 lines → Transition to MEASURE_SIZE
- Implementation complete → Transition to TEST_IMPLEMENTATION
- Review feedback received → Transition to FIX_ISSUES

### Checkpoint Before Exit
- Commit current work
- Update work-log.md
- Note current line count
```

## Orchestrator Spawn Template

```bash
# When orchestrator spawns agents
spawn_with_state_awareness() {
    Task: @agent-software-engineer
    Working directory: efforts/phase1/wave1/api-types
    
    🚨 STARTUP SEQUENCE (R203):
    1. READ: .claude/agents/sw-engineer.md
    2. Your state is: IMPLEMENTATION
    3. READ: agent-states/sw-engineer/IMPLEMENTATION/rules.md
    4. Acknowledge ALL rules before proceeding
    5. Begin IMPLEMENTATION work
}
```

## Benefits of This Pattern

1. **70% Context Reduction**: 15KB instead of 40KB+ per agent
2. **Clear State Focus**: Only relevant rules loaded
3. **Proper Acknowledgment**: Proof of rule understanding
4. **Easy Updates**: Change state rules without touching core
5. **Better Debugging**: Know exactly what state agent is in

## Migration Checklist

- [ ] Update agent core configs to slim version
- [ ] Move state-specific rules to agent-states/
- [ ] Add state detection logic to each agent
- [ ] Update orchestrator to spawn with state
- [ ] Test acknowledgment sequence
- [ ] Verify context reduction

## Remember

**EVERY agent startup MUST**:
1. Load core config
2. Detect current state
3. Load state rules
4. Acknowledge everything
5. Only THEN proceed with work

This is **R203** and it's **BLOCKING** - no compliance = immediate termination!