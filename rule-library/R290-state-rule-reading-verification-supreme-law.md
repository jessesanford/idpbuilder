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
   markers/state-verification/state_rules_read_[agent]_[STATE]-[TIMESTAMP]
   ```
   (Note: Markers organized in `markers/state-verification/` with timestamps)

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

    # Step 2: Create marker directory if needed
    mkdir -p markers/state-verification

    # Step 3: 🔴 MANDATORY - Read state rules FIRST
    local rules_file="agent-states/${agent_type}/${new_state}/rules.md"
    if [[ ! -f "$rules_file" ]]; then
        echo "❌ FATAL: No rules for state $new_state"
        exit 290
    fi

    echo "📖 READING STATE RULES FOR ${new_state}..."
    cat "$rules_file"

    # Step 4: Create verification marker with timestamp
    local timestamp=$(date +%Y%m%d-%H%M%S)
    local marker_file="markers/state-verification/state_rules_read_${agent_type}_${new_state}-${timestamp}"
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

    # Check for marker (most recent one)
    local marker_pattern="markers/state-verification/state_rules_read_${agent_type}_${state}-*"
    local marker_file=$(ls $marker_pattern 2>/dev/null | tail -1)

    if [[ -z "$marker_file" ]]; then
        echo "🔴🔴🔴 FATAL ERROR: R290 VIOLATION DETECTED! 🔴🔴🔴"
        echo "State work attempted in ${state} WITHOUT reading rules!"
        echo "Missing verification marker: $marker_pattern"
        echo "AUTOMATIC FAILURE: -100% penalty"
        exit 290
    fi

    # Check for backward compatibility (old location)
    local old_marker=".state_rules_read_${agent_type}_${state}"
    if [[ -f "$old_marker" ]]; then
        echo "⚠️  Found marker in old location, migrating..."
        local timestamp=$(date +%Y%m%d-%H%M%S)
        mv "$old_marker" "markers/state-verification/state_rules_read_${agent_type}_${state}-${timestamp}"
    fi
}
```

### ✅ CORRECT: Read, Verify, Then Act
```
[TRANSITION] INIT → PLANNING
[READ] agent-states/software-factory/orchestrator/PLANNING/rules.md
[MARKER] markers/state-verification/state_rules_read_orchestrator_PLANNING-20251102-152030 created
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
## Software Factory 3.0 Integration

**State Tracking**: In SF 3.0, state transitions are tracked in `orchestrator-state-v3.json`:
```json
{
  "state_machine": {
    "current_state": "CURRENT_STATE_NAME",
    "previous_state": "PREVIOUS_STATE_NAME",
    "state_history": [...]
  }
}
```

**Compliance**: This rule applies to SF 3.0 state machine with appropriate state name mappings per R516 naming conventions.

**Reference**: See `docs/SOFTWARE-FACTORY-3.0-ARCHITECTURE.md` Part 2 for state machine design.

## State Manager Coordination (SF 3.0)

State Manager provides current state context during startup consultation:
- **STARTUP_CONSULTATION** reads `orchestrator-state-v3.json` → `.state_machine.current_state`
- **directive_report** includes current state name for Orchestrator
- **Orchestrator** uses this to load correct state-specific rules per R203

The bookend pattern ensures agents ALWAYS have correct state context before beginning work.

State Manager does NOT read state-specific rules itself (it only has STARTUP/SHUTDOWN consultation states).

See: R203 (state-aware startup), `.claude/agents/state-manager.md`

