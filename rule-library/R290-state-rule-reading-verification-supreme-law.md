# 🔴🔴🔴 RULE R290: STATE RULE READING AND VERIFICATION - SUPREME LAW #3 🔴🔴🔴

## CRITICALITY: SUPREME
**Status**: ENFORCED  
**Penalty**: -100% (COMPLETE FAILURE)  
**Created**: 2025-08-30  
**Consolidates**: R236 (Mandatory State Rule Reading) + R237 (State Rule Verification Enforcement)  
**Enforces**: STATE-AWARE operation with automated verification  

## THE ABSOLUTE LAW

**STATE RULES MUST BE READ AND VERIFIED BEFORE STATE ACTIONS - NO EXCEPTIONS**

This SUPREME LAW establishes the inviolable sequence:
1. **FIRST**: Enter state (state machine transition)
2. **SECOND**: Read state-specific rules COMPLETELY with verification
3. **THIRD**: Execute state actions per those rules

## PART A: THE REQUIREMENT (from R236)

### 🚨🚨🚨 MANDATORY ENFORCEMENT 🚨🚨🚨

**THE IRON SEQUENCE:**
```markdown
STATE TRANSITION PROTOCOL:
1. ✅ State machine authorizes transition
2. ✅ Agent enters new state
3. 🔴 MANDATORY: Read state rules BEFORE ANY ACTION
4. ✅ Create verification evidence
5. ✅ Execute state-specific behavior per rules
6. ✅ Continue until next transition
```

**VIOLATIONS = IMMEDIATE FAILURE:**
- ❌ Taking ANY action before reading state rules
- ❌ Skipping state rule reading
- ❌ Partial state rule reading
- ❌ Reading wrong state's rules
- ❌ Assuming rules from memory

## PART B: THE ENFORCEMENT (from R237)

### 🚨🚨🚨 MANDATORY VERIFICATION PROTOCOL 🚨🚨🚨

**EVERY STATE TRANSITION MUST CREATE VERIFICATION EVIDENCE:**

1. **READ Operation Evidence**
   ```
   [TIMESTAMP] READ: agent-states/[agent]/[STATE]/rules.md
   ```

2. **Verification Marker File**
   ```
   .state_rules_read_[agent]_[STATE]
   ```

3. **Explicit Acknowledgment**
   ```
   ✅ STATE RULES READ AND ACKNOWLEDGED FOR [STATE]
   ```

4. **Work Authorization**
   ```
   ✅ Verification complete - authorized to proceed with [STATE] work
   ```

**AUTOMATIC DETECTION TRIGGERS FAILURE IF:**
- No READ tool call for state rules file
- No verification marker created
- Work attempted before marker exists
- Marker is stale (>60 seconds old in new context)
- Wrong state rules read (different state)

## PART C: IMPLEMENTATION PATTERNS

### Correct Implementation Pattern

```bash
# MANDATORY on state entry
enter_state() {
    local new_state="$1"
    local agent_type="$2"
    
    # Step 1: Transition authorized by state machine
    echo "Entering state: $new_state"
    
    # Step 2: Clear old verification marker
    local marker_file=".state_rules_read_${agent_type}_${new_state}"
    rm -f "$marker_file" 2>/dev/null
    
    # Step 3: 🔴 MANDATORY - Read state rules FIRST
    local rules_file="agent-states/${agent_type}/${new_state}/rules.md"
    if [[ ! -f "$rules_file" ]]; then
        echo "❌ FATAL: No rules for state $new_state"
        exit 290
    fi
    
    echo "📖 READING STATE RULES FOR ${new_state}..."
    cat "$rules_file"
    
    # Step 4: Create verification marker with timestamp
    echo "$(date +%s) - Rules read for ${new_state}" > "$marker_file"
    
    # Step 5: Explicit acknowledgment required
    echo "✅ STATE RULES READ AND ACKNOWLEDGED FOR ${new_state}"
    echo "📋 Verification marker created: $marker_file"
    
    # Step 6: NOW execute state actions
    echo "✅ Rules loaded and verified, executing state actions..."
    execute_state_actions "$new_state"
}

# MANDATORY check before ANY state work
check_rules_were_read() {
    local state="$1"
    local agent_type="$2"
    local marker_file=".state_rules_read_${agent_type}_${state}"
    
    if [[ ! -f "$marker_file" ]]; then
        echo "🔴🔴🔴 FATAL ERROR: R290 VIOLATION DETECTED! 🔴🔴🔴"
        echo "State work attempted in ${state} WITHOUT reading rules!"
        echo "Missing verification marker: $marker_file"
        echo "AUTOMATIC FAILURE: -100% penalty"
        exit 290
    fi
}
```

### ✅ CORRECT: Read, Verify, Then Act
```
[TRANSITION] INIT → PLANNING
[READ] agent-states/orchestrator/PLANNING/rules.md
[MARKER] .state_rules_read_orchestrator_PLANNING created
[ACK] "Rules loaded and verified for PLANNING state"
[ACTION] Begin planning per state rules
```

### ❌ WRONG: Act Without Reading or Verification
```
[TRANSITION] INIT → PLANNING
[ACTION] Start planning immediately  ← VIOLATION! -100%
```

## 📊 COMPLIANCE METRICS

Track and report:
- States entered: ___
- Rules read: ___
- Markers created: ___
- Compliance rate: MUST BE 100%
- Violations: MUST BE 0

## 🔴 RELATIONSHIP TO OTHER RULES

### Supersedes:
- **R236**: Consolidated here (Part A)
- **R237**: Consolidated here (Part B)
- Any informal rule reading patterns
- Optional verification approaches

### Works With:
- **R231**: Clarifies "immediately" means after rules and verification
- **R203**: Extends state-aware startup to ALL transitions
- **R234**: State traversal includes rule reading at each state
- **R206**: Validates transitions are legal
- **R288**: State file updates
- **State Machine**: Rules are part of state definition

### Penalty Interaction:
- **R231 violation**: -30% to -50%
- **R290 violation**: -100% (COMPLETE FAILURE)
- **Both violated**: -100% (R290 supersedes)

## ⚠️⚠️⚠️ CRITICAL UNDERSTANDING ⚠️⚠️⚠️

**This is not optional - it's MANDATORY PROTOCOL with AUTOMATED ENFORCEMENT**

Every state has specific rules that:
- Define exact behavior for that state
- Override general agent behavior
- Contain critical constraints
- Specify exact outputs required

The system will:
1. **DETECT** when rules aren't read (missing marker)
2. **BLOCK** any state work without verification
3. **FAIL** the entire session for violations
4. **GRADE** based on verification compliance

**Skipping rules = Operating blind = Automatic detection = GUARANTEED FAILURE**

## 🔴 FINAL WORD

**NO AGENT MAY TAKE STATE ACTIONS WITHOUT READING AND VERIFYING STATE RULES FIRST**

This is SUPREME LAW #3 - absolute, inviolable, and mandatory.
Violation = Complete System Failure = -100% penalty = DO NOT PASS GO

The days of skipping state rules are OVER. The system is watching.

---
*Rule R290: Because state rules define state behavior - read them, verify them, or fail*