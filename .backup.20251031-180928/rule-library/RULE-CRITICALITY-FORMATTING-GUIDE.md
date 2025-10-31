# Rule Criticality Formatting Guide

## Purpose
Make rules instantly scannable by AI agents using visual hierarchy that matches cognitive importance. Less visual noise = better absorption.

## Criticality Levels & Formatting

### 🚨🚨🚨 LEVEL 1: BLOCKING CRITICAL
**When to use:** Rules that MUST be followed or work stops immediately
**Formatting:**
```markdown
# 🚨🚨🚨 RULE R001 - Pre-Flight Checks [BLOCKING]
**FAIL = EXIT IMMEDIATELY** | Source: [R001](rule-library/RULE-REGISTRY.md#R001)

- Working directory MUST be correct or `exit 1`
- Git branch MUST match pattern or `exit 1`
- Required files MUST exist or `exit 1`
```
**Characteristics:**
- H1 header (#)
- Three red lights (🚨🚨🚨)
- [BLOCKING] tag
- Bold failure consequence
- Concise bullet points

---

### 🚨🚨 LEVEL 2: MANDATORY
**When to use:** Rules that must be followed but don't cause immediate exit
**Formatting:**
```markdown
## 🚨🚨 RULE R037 - KCP Multi-Tenancy [MANDATORY]
**Required for approval** | [R037](rule-library/RULE-REGISTRY.md#R037)

Every resource MUST include:
- `LogicalCluster` field in spec
- Workspace isolation logic
- RBAC boundary enforcement
```
**Characteristics:**
- H2 header (##)
- Two red lights (🚨🚨)
- [MANDATORY] tag
- Bold requirement statement
- Clear requirements list

---

### 🚨 LEVEL 3: CRITICAL
**When to use:** Important rules that significantly impact quality/grading
**Formatting:**
```markdown
### 🚨 RULE R002 - Size Limits [CRITICAL]
*Exceeding = automatic split required* | [R002](rule-library/RULE-REGISTRY.md#R002)

Maximum 800 lines per effort (measured by line-counter.sh)
```
**Characteristics:**
- H3 header (###)
- Single red light (🚨)
- [CRITICAL] tag
- Italic consequence
- Single line explanation

---

### ⚠️ LEVEL 4: IMPORTANT
**When to use:** Rules that affect workflow but aren't blockers
**Formatting:**
```markdown
#### ⚠️ RULE R016 - TODO Management
Save TODOs every 10-15 messages | [R016](rule-library/RULE-REGISTRY.md#R016)
```
**Characteristics:**
- H4 header (####)
- Warning symbol (⚠️)
- No tag needed
- Inline consequence
- Very concise

---

### ℹ️ LEVEL 5: INFORMATIONAL
**When to use:** Guidelines, best practices, helpful context
**Formatting:**
```markdown
ℹ️ **R099 - Naming Conventions:** Use kebab-case for branches ([R099](rule-library/RULE-REGISTRY.md#R099))
```
**Characteristics:**
- No header
- Info symbol (ℹ️)
- Bold rule name
- Parenthetical link
- One line max

## Visual Hierarchy Principles

### 1. Inverse Pyramid of Detail
- BLOCKING: Most detail (agents MUST understand fully)
- MANDATORY: Good detail (agents need clarity)
- CRITICAL: Moderate detail (key point clear)
- IMPORTANT: Minimal detail (quick reminder)
- INFORMATIONAL: Least detail (FYI only)

### 2. Cognitive Load Management
```
BLOCKING:    ████████████ (100% attention required)
MANDATORY:   ████████░░░░ (75% attention required)
CRITICAL:    ██████░░░░░░ (50% attention required)
IMPORTANT:   ████░░░░░░░░ (25% attention required)
INFO:        ██░░░░░░░░░░ (10% attention required)
```

### 3. Scanning Pattern
Agents scan for:
1. Red lights first (🚨 = "pay attention")
2. Header size second (# = "this is big")
3. Bold/italic third (emphasis points)
4. Content last (actual requirements)

## Implementation Examples

### Before (Verbose Delimiters):
```
┌─────────────────────────────────────────────────────────────────┐
│ RULE R001.0.0 - Pre-Flight Checks                              │
│ Source: rule-library/RULE-REGISTRY.md#R001                     │
├─────────────────────────────────────────────────────────────────┤
│ EVERY AGENT MUST COMPLETE THESE CHECKS BEFORE ANY WORK        │
│ Failure = Immediate Stop (exit 1)                             │
└─────────────────────────────────────────────────────────────────┘
```
**Problems:** 
- 7 lines for 2 lines of content
- Box draws attention to itself, not content
- Hard to scan multiple rules quickly

### After (Criticality-Based):
```markdown
# 🚨🚨🚨 RULE R001 - Pre-Flight Checks [BLOCKING]
**FAIL = EXIT IMMEDIATELY** | [R001](rule-library/RULE-REGISTRY.md#R001)
- Working directory MUST be correct or `exit 1`
```
**Benefits:**
- 3 lines for same content
- Visual hierarchy matches importance
- Instantly scannable

## Quick Reference Card

| Level | Symbol | Header | Tag | Use When |
|-------|--------|--------|-----|----------|
| 1 | 🚨🚨🚨 | # | [BLOCKING] | Failure stops everything |
| 2 | 🚨🚨 | ## | [MANDATORY] | Required for approval |
| 3 | 🚨 | ### | [CRITICAL] | Major impact on grade |
| 4 | ⚠️ | #### | none | Affects workflow |
| 5 | ℹ️ | none | none | Good to know |

## Conversion Checklist

When converting rules from delimiter format:

1. **Assess Criticality**: What happens if this rule is violated?
   - Agent exits? → BLOCKING
   - Work rejected? → MANDATORY
   - Grade reduced? → CRITICAL
   - Workflow disrupted? → IMPORTANT
   - Nothing immediate? → INFORMATIONAL

2. **Apply Format**: Use the template for that level

3. **Trim Content**: Remove redundancy, keep essence

4. **Test Scannability**: Can you understand in 2 seconds?

## Cognitive Science Behind This

### Why It Works for AI Agents

1. **Pattern Matching**: LLMs are trained on markdown with headers
2. **Emoji Recognition**: Strong signal in training data for urgency
3. **Progressive Disclosure**: Most important info first
4. **Chunking**: Bite-sized rules easier to process
5. **Visual Distinctiveness**: Red lights trigger "important" pathways

### What Agents "See"

When Claude/Opus encounters:
- `# 🚨🚨🚨` → "CRITICAL INSTRUCTION INCOMING"
- `**FAIL = EXIT**` → "HARD CONSTRAINT"
- `exit 1` → "TERMINAL CONDITION"

This triggers stronger attention weights than verbose delimiters.

## Migration Priority

Update in this order:
1. Agent configuration files (highest visibility)
2. CRITICAL folder files (startup checks)
3. State-specific rules (contextual)
4. Reference documentation (lowest priority)

## Validation

A well-formatted rule should pass the "2-Second Test":
- Can you identify criticality in <0.5 seconds?
- Can you understand consequence in <1 second?
- Can you grasp requirement in <2 seconds?

If not, simplify further.

## Examples Gallery

### Perfect BLOCKING Rule:
```markdown
# 🚨🚨🚨 RULE R010 - Wrong Directory = STOP [BLOCKING]
**EXIT IMMEDIATELY if pwd != expected** | [R010](rule-library/RULE-REGISTRY.md#R010)
Never attempt to fix with `cd` - just exit 1
```

### Perfect MANDATORY Rule:
```markdown
## 🚨🚨 RULE R151 - Agent Spawn Speed [MANDATORY]
**Must average <5 seconds or FAIL** | [R151](rule-library/RULE-REGISTRY.md#R151)
Parallel spawning required for all multi-agent tasks
```

### Perfect CRITICAL Rule:
```markdown
### 🚨 RULE R002 - Line Count Limits [CRITICAL]
*>800 lines = immediate split* | [R002](rule-library/RULE-REGISTRY.md#R002)
```

### Perfect IMPORTANT Rule:
```markdown
#### ⚠️ RULE R173 - State Preservation
Save state before transitions | [R173](rule-library/RULE-REGISTRY.md#R173)
```

### Perfect INFO Rule:
```markdown
ℹ️ **R099:** Use semantic commit messages ([R099](rule-library/RULE-REGISTRY.md#R099))
```

## The Golden Rule

**If an agent misses a BLOCKING rule, the formatting failed.**

Always optimize for instant recognition over completeness. Agents can check the registry for details, but they must NEVER miss a critical requirement.