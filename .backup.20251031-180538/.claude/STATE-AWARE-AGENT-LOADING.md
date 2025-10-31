# State-Aware Agent Loading Pattern

## 🚨 CRITICAL: Agents Must Load State-Specific Rules

### The Problem We're Solving
Agent definition files have become too large (30-40KB each) because they contain rules for ALL possible states. This wastes context and causes confusion.

### The Solution: State-Aware Loading

## Pattern for Orchestrator When Spawning Agents

```bash
# ❌ OLD WAY - No state awareness
Task: @agent-software-engineer
Working directory: efforts/phase1/wave1/api-types
Implement the effort per IMPLEMENTATION-PLAN.md

# ✅ NEW WAY - State-aware spawning
Task: @agent-software-engineer
Working directory: efforts/phase1/wave1/api-types
Current state: IMPLEMENTATION

CRITICAL STARTUP SEQUENCE:
1. Load your core identity from .claude/agents/sw-engineer.md
2. You are in state: IMPLEMENTATION
3. IMMEDIATELY READ: agent-states/sw-engineer/IMPLEMENTATION/rules.md
4. Follow state-specific rules for implementation
5. On state transition, load new state rules
```

## Pattern for Agents on Startup

```bash
# Agent startup sequence
startup_sequence() {
    echo "═══════════════════════════════════════════════════════"
    echo "AGENT STARTUP: $(date '+%Y-%m-%d %H:%M:%S %Z')"
    echo "═══════════════════════════════════════════════════════"
    
    # Step 1: Identify my type
    AGENT_TYPE="sw-engineer"  # or code-reviewer, orchestrator, architect
    
    # Step 2: Determine my state
    if [ -f "IMPLEMENTATION-PLAN.md" ] && [ ! -f "SPLIT-INVENTORY.md" ]; then
        CURRENT_STATE="IMPLEMENTATION"
    elif [ -f "SPLIT-INVENTORY.md" ]; then
        CURRENT_STATE="SPLIT_WORK"
    elif [ -f "REVIEW-FEEDBACK.md" ]; then
        CURRENT_STATE="FIX_ISSUES"
    else
        CURRENT_STATE="INIT"
    fi
    
    echo "Agent Type: $AGENT_TYPE"
    echo "Current State: $CURRENT_STATE"
    
    # Step 3: Load state-specific rules
    STATE_RULES="agent-states/$AGENT_TYPE/$CURRENT_STATE/rules.md"
    echo "Loading state rules from: $STATE_RULES"
    
    # READ: $STATE_RULES
}
```

## State Transition Pattern

```bash
# When transitioning between states
transition_state() {
    OLD_STATE="$CURRENT_STATE"
    NEW_STATE="$1"
    
    echo "═══════════════════════════════════════════════════════"
    echo "STATE TRANSITION: $OLD_STATE → $NEW_STATE"
    echo "═══════════════════════════════════════════════════════"
    
    # Save checkpoint from old state
    if [ -f "agent-states/$AGENT_TYPE/$OLD_STATE/checkpoint.md" ]; then
        echo "Saving checkpoint per $OLD_STATE requirements..."
        # Save work according to checkpoint.md
    fi
    
    # Load new state rules
    CURRENT_STATE="$NEW_STATE"
    STATE_RULES="agent-states/$AGENT_TYPE/$NEW_STATE/rules.md"
    echo "Loading new state rules from: $STATE_RULES"
    
    # READ: $STATE_RULES
}
```

## Example State-Specific Rule Files

### `agent-states/sw-engineer/IMPLEMENTATION/rules.md`
```markdown
# SW Engineer - IMPLEMENTATION State Rules

## Rules Active in This State
- R106: Implementation efficiency
- R152: Code quality standards
- R176: Workspace isolation
- R198: Line counter usage (every 200 lines)

## What to Do in This State
1. Read IMPLEMENTATION-PLAN.md
2. Implement code in pkg/ directory
3. Check size every 200 lines
4. Write tests alongside code

## Exit Conditions
- Size approaching 700 lines → MEASURE_SIZE
- Implementation complete → TEST_IMPLEMENTATION
- Review feedback received → FIX_ISSUES
```

### `agent-states/sw-engineer/SPLIT_WORK/rules.md`
```markdown
# SW Engineer - SPLIT_WORK State Rules

## Rules Active in This State
- R202: Single agent handles ALL splits sequentially
- R007: Each split must be <800 lines
- R200: Measure only changeset

## What to Do in This State
1. Read SPLIT-INVENTORY.md for all splits
2. Implement splits SEQUENTIALLY (never parallel)
3. Complete each split before starting next
4. Commit and push after each split

## Exit Conditions
- All splits complete → COMPLETE
- Split exceeds size → ERROR_STOP
```

## Benefits of State-Aware Loading

1. **Reduced Context Usage**: Load only ~10KB instead of 40KB
2. **Clearer Instructions**: State-specific guidance
3. **Fewer Conflicts**: Rules don't overlap between states
4. **Easier Debugging**: Know exactly which rules apply
5. **Better Performance**: Faster agent startup

## Migration Path

### Phase 1: Add State Awareness (NOW)
- Orchestrator includes state in spawn instructions
- Agents check for state indicators on startup

### Phase 2: Move Rules to States (NEXT)
- Gradually move rules from agent definitions to state files
- Keep only universal rules in agent definitions

### Phase 3: Enforce State Loading (FUTURE)
- Agents refuse to proceed without loading state rules
- Validation that correct state rules were loaded

## Quick Reference

| Agent | State Indicator Files | State to Load |
|-------|----------------------|---------------|
| SW Engineer | IMPLEMENTATION-PLAN.md | IMPLEMENTATION |
| SW Engineer | SPLIT-INVENTORY.md | SPLIT_WORK |
| SW Engineer | REVIEW-FEEDBACK.md | FIX_ISSUES |
| Code Reviewer | Request: "create plan" | PLANNING |
| Code Reviewer | Request: "review code" | CODE_REVIEW |
| Code Reviewer | Effort >800 lines | SPLIT_PLANNING |
| Orchestrator | orchestrator-state-v3.json | Per current_state field |
| Architect | Request: "wave review" | REVIEW_WAVE_ARCHITECTURE |
| Architect | Request: "phase review" | PHASE_REVIEW |

## Remember

**The Goal**: Agents should load ONLY the rules they need for their current state, not ALL rules for ALL states.

**The Pattern**: Core Identity (small) + Current State Rules (specific) = Optimal Context Usage