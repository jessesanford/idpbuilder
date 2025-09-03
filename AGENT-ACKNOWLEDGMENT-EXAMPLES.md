# Agent Acknowledgment Examples

## Purpose

Every agent MUST acknowledge their identity and rules at startup. This document shows the exact output expected from each agent type.

## 🎯 Orchestrator Agent

**Agent ID:** `orchestrator`  
**Common States:** `INIT`, `WAVE_START`, `WAVE_COMPLETE`, `SPAWN_AGENTS`, `MONITOR`

### Expected Output:
```
================================
RULE ACKNOWLEDGMENT
I am orchestrator in state INIT
I acknowledge these rules:
--------------------------------
[AGENT MUST READ AND LIST THEIR OWN RULES HERE]
[Include all CRITICAL and BLOCKING rules from:
 - .claude/agents/orchestrator.md
 - Referenced rule files
 Format: R###: Rule description [CRITICALITY]]
================================
```

## ⚙️ Software Engineer Agent

**Agent ID:** `sw-engineer`  
**Common States:** `IMPLEMENTATION`, `MEASURE_SIZE`, `FIX_REVIEW_ISSUES`, `SPLIT_IMPLEMENTATION`

### Expected Output:
```
================================
RULE ACKNOWLEDGMENT
I am sw-engineer in state IMPLEMENTATION
I acknowledge these rules:
--------------------------------
[AGENT MUST READ AND LIST THEIR OWN RULES HERE]
[Include all CRITICAL and BLOCKING rules from:
 - .claude/agents/sw-engineer.md
 - Referenced rule files
 Format: R###: Rule description [CRITICALITY]]
================================
```

## 🔍 Code Reviewer Agent

**Agent ID:** `@agent-code-reviewer`  
**Common States:** `EFFORT_PLAN_CREATION`, `CODE_REVIEW`, `CREATE_SPLIT_PLAN`, `SPLIT_REVIEW`

### Expected Output:
```
================================
RULE ACKNOWLEDGMENT
I am @agent-code-reviewer in state CODE_REVIEW
I acknowledge these rules:
--------------------------------
[AGENT MUST READ AND LIST THEIR OWN RULES HERE]
[Include all CRITICAL and BLOCKING rules from:
 - .claude/agents/code-reviewer.md
 - Referenced rule files
 Format: R###: Rule description [CRITICALITY]]
================================
```

## 🏗️ Architect Agent

**Agent ID:** `@agent-architect`  
**Common States:** `PHASE_ASSESSMENT`, `WAVE_REVIEW`, `INTEGRATION_REVIEW`

### Expected Output:
```
================================
RULE ACKNOWLEDGMENT
I am @agent-architect in state WAVE_REVIEW
I acknowledge these rules:
--------------------------------
[AGENT MUST READ AND LIST THEIR OWN RULES HERE]
[Include all CRITICAL and BLOCKING rules from:
 - .claude/agents/architect.md
 - Referenced rule files
 Format: R###: Rule description [CRITICALITY]]
================================
```

## Key Requirements

### 1. Identity Declaration
Every agent MUST declare:
- Their correct agent ID (`@agent-{name}`)
- Their current state (from state machine)

### 2. Rule Acknowledgment
Every agent MUST acknowledge:
- Their role-specific rules (CRITICAL level)
- Universal TODO persistence rules (R187-R190)
- Target repository rules (R191-R195 as applicable)

### 3. Format Consistency
- Use the exact format shown above
- Include rule numbers and criticality levels
- Use proper spacing and separators

## Common Mistakes to Avoid

❌ **Wrong Agent ID:**
```
I am @agent-sw-engineer  # WRONG - should be sw-engineer
```

❌ **Missing State:**
```
I am orchestrator  # WRONG - missing "in state {STATE}"
```

❌ **Incomplete Rule List:**
```
I acknowledge these rules:
R006, R008, R009  # WRONG - missing descriptions and criticality
```

✅ **Correct Format:**
```
I am sw-engineer in state IMPLEMENTATION
I acknowledge these rules:
--------------------------------
R054: Implementation plan creation [CRITICAL]
```

## Verification

The orchestrator should verify that spawned agents:
1. Print their acknowledgment within first 3 messages
2. Use the correct agent ID
3. Include all required rules
4. Specify their current state

## Grading Impact

Failure to properly acknowledge:
- Wrong agent ID: -10% (identity confusion)
- Missing rules: -5% per missing rule
- No acknowledgment: -20% (non-compliance)
- Wrong format: -5% (process violation)