# Agent Rule Refactoring Plan

## Problem Analysis

The agent definition files have become bloated:
- **orchestrator.md**: 43KB (1,000+ lines)
- **sw-engineer.md**: 33KB (800+ lines)  
- **code-reviewer.md**: 26KB (600+ lines)
- **architect.md**: 18KB (400+ lines)

These files contain many state-specific rules that should be in the `agent-states/` directories.

## Current State

### What's in Agent Definition Files (but shouldn't be)
1. **Implementation-specific rules** (should be in IMPLEMENTATION state)
2. **Split handling rules** (should be in SPLIT_WORK state)
3. **Measurement rules** (should be in MEASURE_SIZE state)
4. **Review-specific rules** (should be in CODE_REVIEW state)
5. **Detailed bash examples** (should be in relevant states)

### What SHOULD be in Agent Definition Files
1. **Core identity** (who the agent is)
2. **Universal pre-flight checks** (R001, R186 compaction detection)
3. **State machine awareness** (how to navigate states)
4. **Critical grading metrics** (R151 parallel spawning)
5. **Pointer to state-specific rules** (not the rules themselves)

## Refactoring Plan

### Phase 1: Agent Definition Slimming

#### SW Engineer (`sw-engineer.md`)
**Keep:**
- Core identity and termination rules (R197)
- Pre-flight checks (R001, R186)
- State machine navigation guide
- Project root finding pattern

**Move to States:**
- R200 (changeset measurement) → MEASURE_SIZE/rules.md
- R202 (split handling) → SPLIT_WORK/rules.md
- R198 (line counter usage) → MEASURE_SIZE/rules.md
- Workspace isolation details → IMPLEMENTATION/rules.md
- Git operations → COMMIT_PUSH/rules.md

#### Code Reviewer (`code-reviewer.md`)
**Keep:**
- Core identity
- Pre-flight checks
- Review decision framework
- State navigation

**Move to States:**
- Implementation plan creation → PLANNING/rules.md
- Size measurement details → CODE_REVIEW/rules.md
- Split planning (R199) → SPLIT_PLANNING/rules.md
- Test validation → VALIDATION/rules.md

#### Orchestrator (`orchestrator.md`)
**Keep:**
- Core identity (never writes code)
- Parallel spawning (R151)
- Pre-flight checks
- State machine orchestration

**Move to States:**
- Wave-specific rules → WAVE_START/rules.md
- Integration rules → INTEGRATION/rules.md
- Spawn patterns → SPAWN_AGENTS/rules.md
- Monitoring rules → MONITOR/rules.md

### Phase 2: State-Specific Enhancement

Each state directory should have:
```
agent-states/[agent]/[STATE]/
├── rules.md         # Rules specific to this state
├── checkpoint.md    # What to save before leaving state
├── grading.md       # How this state is graded
└── examples.md      # Detailed bash/code examples
```

### Phase 3: State Transition Guide

Create clear instructions for agents:

```markdown
## On State Entry
1. Read core agent definition (once at startup)
2. Identify current state from context/instructions
3. READ: agent-states/[my-type]/[CURRENT_STATE]/rules.md
4. Follow state-specific rules

## On State Exit
1. Complete checkpoint requirements
2. Save state per checkpoint.md
3. Transition to next state
4. Load new state rules
```

## Benefits

1. **Smaller agent files** - Easier to load and parse
2. **State-specific context** - Only load rules when needed
3. **Clearer separation** - Core identity vs. state behavior
4. **Easier updates** - Change state rules without touching agent definition
5. **Better performance** - Less context to process initially

## Implementation Priority

1. **HIGH**: Move split-handling rules to SPLIT_WORK states
2. **HIGH**: Move measurement rules to MEASURE_SIZE states
3. **MEDIUM**: Move implementation details to IMPLEMENTATION states
4. **LOW**: Move examples to state-specific example files

## Example: Slimmed SW Engineer

```markdown
# SW Engineer Agent Definition (SLIM VERSION)

## Core Identity
You are a software engineer who implements code.
You work on ONE effort only (R197).
You never create clones or branches (R196).

## Pre-Flight Checks
[R001 and R186 checks remain here]

## State Machine Navigation
When spawned, determine your state:
- If IMPLEMENTATION-PLAN.md exists → READ: agent-states/sw-engineer/IMPLEMENTATION/rules.md
- If SPLIT-INVENTORY.md exists → READ: agent-states/sw-engineer/SPLIT_WORK/rules.md
- If REVIEW-FEEDBACK.md exists → READ: agent-states/sw-engineer/FIX_ISSUES/rules.md

## Critical Universal Rules
- Find project root: [pattern stays here]
- Line counter: $PROJECT_ROOT/tools/line-counter.sh
- For detailed usage → See state-specific rules

## State Transition
Follow state machine, loading new rules at each transition.
```

## Metrics

**Before Refactoring:**
- Agent files: 120KB total
- Initial load: All rules for all states
- Context usage: High

**After Refactoring:**
- Agent files: ~40KB total (66% reduction)
- Initial load: Core identity + current state only
- Context usage: Optimized per state

## Next Steps

1. Start with SW Engineer as pilot
2. Move rules incrementally
3. Test state transitions
4. Apply pattern to other agents
5. Update orchestrator spawn instructions to include state