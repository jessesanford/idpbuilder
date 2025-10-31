# Rule Delimiter and Criticality System

## Combined Approach: Clean Delimiters + Visual Hierarchy

### The Format Template

```markdown
---
### [EMOJI] RULE R### - Rule Title
**Source:** rule-library/RULE-REGISTRY.md#R###
**Criticality:** [LEVEL] - [Consequence]

[Rule content here]
---
```

## Criticality Levels with Delimiters

### LEVEL 1: BLOCKING (🚨🚨🚨)
```markdown
---
### 🚨🚨🚨 RULE R001 - Pre-Flight Checks
**Source:** rule-library/RULE-REGISTRY.md#R001
**Criticality:** BLOCKING - Failure = Immediate Stop (exit 1)

- Working directory MUST be correct or `exit 1`
- Git branch MUST match pattern or `exit 1`
- Required files MUST exist or `exit 1`
---
```

### LEVEL 2: MANDATORY (🚨🚨)
```markdown
---
### 🚨🚨 RULE R037 - KCP Multi-Tenancy
**Source:** rule-library/RULE-REGISTRY.md#R037
**Criticality:** MANDATORY - Required for approval

Every resource MUST include:
- `LogicalCluster` field in spec
- Workspace isolation logic
- RBAC boundary enforcement
---
```

### LEVEL 3: CRITICAL (🚨)
```markdown
---
### 🚨 RULE R002 - Size Limits
**Source:** rule-library/RULE-REGISTRY.md#R002
**Criticality:** CRITICAL - Exceeding = automatic split required

Maximum 800 lines per effort (measured by line-counter.sh)
---
```

### LEVEL 4: IMPORTANT (⚠️)
```markdown
---
### ⚠️ RULE R016 - TODO Management
**Source:** rule-library/RULE-REGISTRY.md#R016
**Criticality:** IMPORTANT - Save TODOs every 10-15 messages
---
```

### LEVEL 5: INFORMATIONAL (ℹ️)
```markdown
---
### ℹ️ RULE R099 - Naming Conventions
**Source:** rule-library/RULE-REGISTRY.md#R099
**Criticality:** INFO - Use kebab-case for branches
---
```

## Why This Works

1. **Clear Boundaries**: `---` delimiters are minimal but clear
2. **Scannable Headers**: Emoji count = importance level
3. **Consistent Structure**: Every rule has same 3-line header
4. **Machine Parseable**: Can grep for `^---$` to find rules
5. **Human Readable**: Clean markdown formatting

## Parsing Rules

To find all rules in a file:
```bash
# Find all rule starts
grep -n "^### .* RULE R" file.md

# Find all BLOCKING rules
grep -n "^### 🚨🚨🚨 RULE" file.md

# Extract rule numbers
grep "^### .* RULE R" file.md | sed 's/.*RULE R\([0-9]*\).*/R\1/'
```

## Migration Examples

### OLD (Verbose Box):
```
┌─────────────────────────────────────────────────────────────────┐
│ RULE R001.0.0 - Pre-Flight Checks                              │
│ Source: rule-library/RULE-REGISTRY.md#R001                     │
├─────────────────────────────────────────────────────────────────┤
│ EVERY AGENT MUST COMPLETE THESE CHECKS BEFORE ANY WORK        │
│ Failure = Immediate Stop (exit 1)                             │
└─────────────────────────────────────────────────────────────────┘
```

### NEW (Clean + Hierarchical):
```markdown
---
### 🚨🚨🚨 RULE R001 - Pre-Flight Checks
**Source:** rule-library/RULE-REGISTRY.md#R001
**Criticality:** BLOCKING - Failure = Immediate Stop (exit 1)

EVERY AGENT MUST COMPLETE THESE CHECKS BEFORE ANY WORK
---
```

## File Update Priority

1. **Agent Configs** (4 files) - Highest visibility
   - `.claude/agents/orchestrator.md`
   - `.claude/agents/sw-engineer.md`
   - `.claude/agents/code-reviewer.md`
   - `.claude/agents/architect.md`

2. **Rule Library** - Core enforcement rules
   - `rule-library/R002-agent-acknowledgment.md`
   - `rule-library/R003-performance-grading.md`
   - `rule-library/R006-orchestrator-never-writes-code.md`
   - `rule-library/R007-size-limit-compliance.md`
   - `rule-library/R008-monitoring-frequency.md`
   - `rule-library/R009-integration-branch-creation.md`
   - `rule-library/R152-implementation-speed.md`
   - `rule-library/R153-review-turnaround.md`
   - `rule-library/R158-pattern-compliance-rate.md`

3. **State Rules** (~45 files) - Context-specific
   - `agent-states/orchestrator/*/rules.md`
   - `agent-states/sw-engineer/*/rules.md`
   - `agent-states/code-reviewer/*/rules.md`
   - `agent-states/architect/*/rules.md`

4. **State Machines** (4 files)
   - `state-machines/orchestrator.md`
   - `state-machines/sw-engineer.md`
   - `state-machines/code-reviewer.md`
   - `state-machines/architect.md`

## Criticality Assignment Guide

When converting a rule, ask:

1. **Does failure stop the agent?** → 🚨🚨🚨 BLOCKING
2. **Is it required for approval?** → 🚨🚨 MANDATORY  
3. **Does it affect grading significantly?** → 🚨 CRITICAL
4. **Does it disrupt workflow?** → ⚠️ IMPORTANT
5. **Is it just guidance?** → ℹ️ INFO

## Bulk Conversion Script

```bash
#!/bin/bash
# convert-rules.sh - Convert box delimiters to new format

# Find all files with old box delimiters
files=$(grep -r "┌─────" --include="*.md" . | cut -d: -f1 | sort -u)

for file in $files; do
    echo "Converting $file..."
    
    # Replace box tops with ---
    sed -i 's/┌─.*┐/---/' "$file"
    
    # Replace box middles with header format
    sed -i 's/│ RULE R\([0-9.]*\) - \(.*\) *│/### RULE R\1 - \2/' "$file"
    
    # Replace source lines
    sed -i 's/│ Source: \(.*\) *│/**Source:** \1/' "$file"
    
    # Replace box bottoms with ---
    sed -i 's/└─.*┘/---/' "$file"
    
    # Remove box side characters
    sed -i 's/^│ //' "$file"
    sed -i 's/ *│$//' "$file"
done
```

## Validation Checklist

After conversion, each rule should:
- [ ] Start with `---`
- [ ] Have emoji matching criticality
- [ ] Include rule number and title
- [ ] Show source reference
- [ ] State criticality level
- [ ] End with `---`
- [ ] Be scannable in 2 seconds

## Example: Complete Agent Config Section

```markdown
# 🎯 SOFTWARE FACTORY 2.0 - ORCHESTRATOR AGENT

## 🚨🚨🚨 MANDATORY PRE-FLIGHT CHECKS

---
### 🚨🚨🚨 RULE R001 - Pre-Flight Checks
**Source:** rule-library/RULE-REGISTRY.md#R001
**Criticality:** BLOCKING - Failure = Immediate Stop (exit 1)

[Pre-flight check code here]
---

---
### 🚨🚨 RULE R010 - Wrong Location Handling
**Source:** rule-library/RULE-REGISTRY.md#R010
**Criticality:** MANDATORY - Working in wrong location = GRADING FAILURE

- NEVER attempt to cd or checkout to "fix"
- Report error and stop immediately
---

---
### 🚨 RULE R006 - Orchestrator Never Codes
**Source:** rule-library/RULE-REGISTRY.md#R006
**Criticality:** CRITICAL - Writing code = automatic FAIL

Orchestrator coordinates ONLY, never implements
---
```

This provides clean visual separation while maintaining hierarchy!