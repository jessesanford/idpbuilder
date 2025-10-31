# ⚠️⚠️⚠️ DEPRECATED - CONSOLIDATED INTO R290 ⚠️⚠️⚠️

**This rule has been consolidated into R290: State Rule Reading and Verification - Supreme Law**
**See: [R290](R290-state-rule-reading-verification-supreme-law.md)**
**Deprecated: 2025-08-30**

---

# 🔴🔴🔴 RULE R236: MANDATORY STATE RULE READING - SUPREME LAW #3 🔴🔴🔴

## CRITICALITY: SUPREME
**Status**: ENFORCED  
**Penalty**: -100% (COMPLETE FAILURE)  
**Created**: 2025-01-29  
**Supersedes**: Implicit behavior in R231, R203  
**Enforces**: STATE-AWARE operation  

## THE ABSOLUTE LAW

**STATE RULES MUST BE READ BEFORE STATE ACTIONS - NO EXCEPTIONS**

This SUPREME LAW establishes the inviolable sequence:
1. **FIRST**: Enter state (state machine transition)
2. **SECOND**: Read state-specific rules COMPLETELY
3. **THIRD**: Execute state actions per those rules

## 🚨🚨🚨 MANDATORY ENFORCEMENT 🚨🚨🚨

### THE IRON SEQUENCE

```markdown
STATE TRANSITION PROTOCOL:
1. ✅ State machine authorizes transition
2. ✅ Agent enters new state
3. 🔴 MANDATORY: Read state rules BEFORE ANY ACTION
4. ✅ Execute state-specific behavior per rules
5. ✅ Continue until next transition
```

### VIOLATIONS = IMMEDIATE FAILURE

**ANY of these = -100% penalty:**
- ❌ Taking ANY action before reading state rules
- ❌ Skipping state rule reading
- ❌ Partial state rule reading
- ❌ Reading wrong state's rules
- ❌ Assuming rules from memory

## 📋 IMPLEMENTATION REQUIREMENTS

### For All Agents

```bash
# MANDATORY on state entry
enter_state() {
    local new_state="$1"
    
    # Step 1: Transition authorized by state machine
    echo "Entering state: $new_state"
    
    # Step 2: 🔴 MANDATORY - Read state rules FIRST
    local rules_file="agent-states/${AGENT_TYPE}/${new_state}/rules.md"
    if [[ ! -f "$rules_file" ]]; then
        echo "❌ FATAL: No rules for state $new_state"
        exit 1
    fi
    
    echo "📖 Reading state rules from $rules_file..."
    cat "$rules_file"
    
    # Step 3: NOW execute state actions
    echo "✅ Rules loaded, executing state actions..."
    execute_state_actions "$new_state"
}
```

### Clarification of R231

This rule CLARIFIES R231's "immediately continue working":
- "Immediately" means AFTER reading state rules
- The sequence is: transition → read rules → work
- NO work before rules are read

### Audit Requirements

```bash
# Verify compliance
audit_state_rule_reading() {
    echo "🔍 Auditing R236 compliance..."
    
    # Check if rules were read before actions
    local last_transition=$(grep "Entering state:" agent.log | tail -1)
    local rules_read=$(grep "Reading state rules" agent.log | tail -1)
    local first_action=$(grep "Executing action" agent.log | tail -1)
    
    # Timestamps must show: transition < rules < action
    # Any deviation = FAILURE
}
```

## 🔴 RELATIONSHIP TO OTHER RULES

### Supersedes
- **Implicit behaviors**: Any assumption that work starts immediately on transition
- **Old interpretations**: That R231 means skip reading to work faster

### Works With
- **R231**: Clarifies "immediately" means after rules
- **R203**: Extends state-aware startup to ALL transitions
- **R234**: State traversal includes rule reading at each state
- **State Machine**: Rules are part of state definition

### Penalty Interaction
- **R231 violation**: -30% to -50%
- **R236 violation**: -100% (COMPLETE FAILURE)
- **Both violated**: -100% (R236 supersedes)

## 📊 COMPLIANCE METRICS

Track and report:
- States entered: ___
- Rules read: ___
- Compliance rate: MUST BE 100%
- Violations: MUST BE 0

## ⚠️⚠️⚠️ CRITICAL UNDERSTANDING ⚠️⚠️⚠️

**This is not optional optimization - it's MANDATORY PROTOCOL**

Every state has specific rules that:
- Define exact behavior for that state
- Override general agent behavior
- Contain critical constraints
- Specify exact outputs required

**Skipping rules = Operating blind = GUARANTEED FAILURE**

## 🚨 ENFORCEMENT EXAMPLES

### ✅ CORRECT: Read Then Act
```
[TRANSITION] INIT → PLANNING
[READ] agent-states/orchestrator/PLANNING/rules.md
[ACK] "Rules loaded for PLANNING state"
[ACTION] Begin planning per state rules
```

### ❌ WRONG: Act Without Reading
```
[TRANSITION] INIT → PLANNING
[ACTION] Start planning immediately  ← VIOLATION! -100%
```

### ❌ WRONG: Partial Reading
```
[TRANSITION] INIT → PLANNING
[READ] Skim first few lines  ← VIOLATION! -100%
[ACTION] Miss critical rules at bottom
```

## 🔴 FINAL WORD

**NO AGENT MAY TAKE STATE ACTIONS WITHOUT READING STATE RULES FIRST**

This is SUPREME LAW #3 - absolute, inviolable, and mandatory.
Violation = Complete System Failure = -100% penalty = DO NOT PASS GO

---
*Rule R236: Because state rules define state behavior - ignore them at your peril*