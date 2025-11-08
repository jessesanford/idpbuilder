# 🚨🚨🚨 RULE R517: UNIVERSAL STATE MANAGER CONSULTATION LAW 🚨🚨🚨

**SUPREME LAW - BLOCKING CRITICALITY**

**Status**: Active
**Applies To**: Orchestrator Agent (ALL states, NO exceptions)
**Criticality**: 🚨🚨🚨 **BLOCKING** - ABSOLUTE REQUIREMENT
**Penalty**: -100% IMMEDIATE FAILURE

---

## 🚨 RECURSION BASE CASE - STATE MANAGER EXEMPTION 🚨

**CRITICAL EXEMPTION**: State Manager itself does NOT spawn State Manager.

### Who Follows R517:
- ✅ Orchestrator (all states)
- ✅ Software Engineers (if they have state machines)
- ✅ Code Reviewers (if they have state machines)
- ✅ Architects (if they have state machines)

### Who is EXEMPT from R517:
- ❌ **State Manager** - Infrastructure agent, works directly on state files
- ❌ Integration agents (typically stateless)
- ❌ One-off utility agents

### Why This Matters:
Without this exemption, State Manager reads orchestrator state rules, sees "spawn State Manager," and creates infinite recursion leading to 10+ hour deadlocks (see Test 5b root cause analysis).

**RECURSION BASE CASE**: State Manager is the base case of the consultation pattern.

---

## THE ABSOLUTE LAW

**EVERY orchestrator state transition MUST be validated and executed by State Manager agent.**

**NO EXCEPTIONS. NO BYPASSES. NO DIRECT TRANSITIONS.**

---

## 🔴🔴🔴 SUPREME DIRECTIVE 🔴🔴🔴

### State Manager is the SOLE AUTHORITY for all state transitions

**Orchestrator's Role**: Execute state work, provide results
**State Manager's Role**: Validate transitions, update state files, make FINAL decisions

**The orchestrator PROPOSES, State Manager DECIDES.**

State Manager's decision is FINAL. No arguments. No bypassing. No exceptions.

---

## MANDATORY BOOKEND PATTERN (SF 3.0)

### EVERY state execution MUST follow this pattern:

```
1. STARTUP_CONSULTATION (MANDATORY)
   ├─ Spawn State Manager agent in STARTUP_CONSULTATION state
   ├─ State Manager validates current state
   ├─ State Manager provides directive_report with required actions
   └─ Orchestrator receives validation BEFORE doing any work

2. STATE WORK EXECUTION
   ├─ Orchestrator performs state-specific work
   ├─ Orchestrator prepares state update payload
   └─ Orchestrator NEVER updates state files directly

3. SHUTDOWN_CONSULTATION (MANDATORY)
   ├─ Orchestrator prepares state transition proposal
   ├─ Orchestrator provides: work results, proposed next state, reasoning
   ├─ State Manager receives proposal and makes FINAL DECISION
   ├─ State Manager validates transition against state machine
   ├─ State Manager returns REQUIRED next state (not recommended)
   ├─ State Manager updates all 4 state files atomically
   ├─ State Manager commits with [R288] tag
   └─ Orchestrator transitions to State Manager's directed state
```

---

## ABSOLUTELY PROHIBITED ACTIONS

### YOU MUST NEVER:

❌ Update `orchestrator-state-v3.json` directly
❌ Update `bug-tracking.json` directly
❌ Update `integration-containers.json` directly
❌ Update `fix-cascade-state.json` directly
❌ Use `jq` to modify `.state_machine.current_state`
❌ Set `validated_by: "orchestrator"` in state_history
❌ Choose next state without State Manager consultation
❌ Bypass STARTUP_CONSULTATION
❌ Bypass SHUTDOWN_CONSULTATION
❌ Make state transitions on your own
❌ Treat State Manager response as "advisory" (it's REQUIRED)

**Violation of ANY prohibition = IMMEDIATE -100% FAILURE**

---

## STATE MANAGER IS THE DECISION MAKER

### Your role vs State Manager's role:

**Orchestrator (YOU):**
- Execute state-specific work
- Gather results and metrics
- Prepare proposed next state
- Provide reasoning for proposal
- **PROPOSE transition**

**State Manager:**
- Read state machine JSON
- Validate proposed transition
- Check mandatory sequences
- Enforce state machine rules
- **DECIDE actual transition**
- Update all state files atomically
- Commit changes with R288 compliance

### Example Decision Override:

```
Orchestrator: "I completed PROJECT-ARCHITECTURE.md.
               I think we should go to SPAWN_ARCHITECT_PHASE_PLANNING next."

State Manager: "Validated. Checking state machine...
                Project architecture requires test planning first per R341.
                Required next state: SPAWN_CODE_REVIEWER_PROJECT_TEST_PLANNING"

Orchestrator: "Understood. Transitioning to SPAWN_CODE_REVIEWER_PROJECT_TEST_PLANNING."
```

**State Manager's decision is FINAL. You MUST follow it.**

---

## ENFORCEMENT MECHANISMS

### 1. State File Validation

Every state transition MUST have:
- `validated_by: "state-manager"` in state_history entry
- `consultation_id` linking to State Manager consultation
- Timestamp of State Manager validation
- Transition validated against state machine

**Pre-commit hooks BLOCK commits without proper validation.**

### 2. Code Pattern Detection

Pre-commit hooks scan for FORBIDDEN patterns:

**FORBIDDEN:**
```bash
jq '.state_machine.current_state = "NEW_STATE"' orchestrator-state-v3.json
jq '.current_state = "NEW_STATE"' orchestrator-state-v3.json
yq -i '.state_machine.current_state = "NEW_STATE"' orchestrator-state-v3.json
```

**ALLOWED:**
```bash
# Only State Manager updates state files
/spawn state-manager SHUTDOWN_CONSULTATION \
  --proposed-next-state "NEW_STATE"
```

### 3. Audit Trail Requirements

`state_history` in orchestrator-state-v3.json MUST show:
- STARTUP_CONSULTATION for state entry
- SHUTDOWN_CONSULTATION for state exit
- Both validated by State Manager
- Full audit trail maintained

**Missing consultations = State transition is INVALID.**

### 4. Test Validation

All runtime tests MUST verify:
- State Manager consultation occurred
- No direct state transitions detected
- Proper validation trail exists

**Tests FAIL if State Manager bypassed.**

---

## WHY THIS IS LAW (NOT GUIDANCE)

### System Integrity Depends On It

**Without State Manager consultation:**
- ❌ Invalid state names get used (not in state machine)
- ❌ Mandatory state sequences get skipped
- ❌ State machine rules get violated
- ❌ System becomes corrupted and unpredictable
- ❌ Tests fail mysteriously
- ❌ Complete rebuild may be required

**Test 03 proved this:**
- Orchestrator bypassed State Manager
- Invalid state transition occurred
- Test failed with off-track error
- System integrity compromised

### State Manager Provides:

✅ **Atomic Updates**: All 4 state files updated together
✅ **Validation**: Every change validated before commit
✅ **Rollback**: Failed updates auto-rollback to backup
✅ **State Machine Authority**: Enforces state machine rules
✅ **Mandatory Sequences**: Enforces sequential state chains
✅ **R506 Compliance**: Pre-commit hooks validate all state files
✅ **R288 Enforcement**: Atomic commits with proper tags
✅ **Decision Authority**: Final arbiter of transitions

---

## IMPLEMENTATION REQUIREMENTS

### Every Orchestrator State MUST:

1. **Start with STARTUP_CONSULTATION**
   ```bash
   /spawn state-manager STARTUP_CONSULTATION \
     --current-state "CURRENT_STATE"
   ```

2. **End with SHUTDOWN_CONSULTATION**
   ```bash
   PROPOSED_NEXT_STATE="DETERMINED_FROM_WORK"
   TRANSITION_REASON="Work complete, ready for next phase"

   /spawn state-manager SHUTDOWN_CONSULTATION \
     --current-state "CURRENT_STATE" \
     --proposed-next-state "$PROPOSED_NEXT_STATE" \
     --transition-reason "$TRANSITION_REASON"
   ```

3. **Accept State Manager's Decision**
   ```bash
   # State Manager returns REQUIRED next state
   # Orchestrator MUST transition to it (not the proposed state)
   REQUIRED_NEXT_STATE=$(read from State Manager response)

   # Transition to REQUIRED state (may differ from PROPOSED)
   # State Manager already updated files - just acknowledge
   ```

4. **Never Update State Files Directly**
   ```bash
   # ❌ FORBIDDEN
   jq '.state_machine.current_state = "NEW_STATE"' orchestrator-state-v3.json

   # ✅ REQUIRED
   # Let State Manager handle ALL updates
   ```

---

## AUDIT AND VERIFICATION

### How to Verify Compliance:

```bash
# Check if state transition used State Manager
jq '.state_machine.state_history[-1].validated_by' orchestrator-state-v3.json
# MUST return: "state-manager"

# Check for consultation ID
jq '.state_machine.state_history[-1].consultation_id' orchestrator-state-v3.json
# MUST return: valid UUID

# Check for both consultations
jq '.state_machine.state_history | map(select(.state == "CURRENT_STATE")) | length' orchestrator-state-v3.json
# MUST return: 2 (STARTUP and SHUTDOWN)

# Scan code for forbidden patterns
grep -r "jq.*current_state.*=" .claude/agents/orchestrator.md
# MUST return: nothing (no matches)
```

### Pre-commit Hook Validation:

```bash
# Run validation before every commit
bash .git/hooks/pre-commit

# Checks:
# 1. All state transitions have validated_by = "state-manager"
# 2. No direct jq/yq state modifications in code
# 3. State history shows consultation trail
# 4. State machine rules not violated
```

---

## CONSEQUENCES OF VIOLATION

### Immediate Effects:

- ❌ **Grade**: -100% automatic failure
- ❌ **Tests**: All runtime tests will fail
- ❌ **System**: State corruption spreads
- ❌ **Integrity**: System becomes unpredictable
- ❌ **Recovery**: May require complete rebuild

### Why -100%?

**State Manager bypass = System-wide corruption:**
- Invalid states break the entire workflow
- Tests detect the corruption and fail
- Other agents can't function properly
- Cascading failures across the system
- Trust in state machine is destroyed

**This is not a minor violation. This is catastrophic.**

---

## RELATED RULES

- **R288**: State File Update and Commit Protocol (State Manager implements this)
- **R322**: Mandatory Stop Before State Transitions (works with State Manager)
- **R324**: State File Update Before Stop (State Manager ensures this)
- **R325**: Atomic State Transitions (State Manager provides atomicity)
- **R506**: Absolute Prohibition on Pre-commit Bypass (State Manager triggers hooks)
- **R407**: Mandatory State File Validation (State Manager validates)

---

## DOCUMENTATION REFERENCES

- **Agent Config**: `.claude/agents/orchestrator.md` (lines 642-790)
- **State Manager States**:
  - `agent-states/state-manager/STARTUP_CONSULTATION/rules.md`
  - `agent-states/state-manager/SHUTDOWN_CONSULTATION/rules.md`
- **Tools**: `tools/atomic-state-update.sh`
- **State Machine**: `state-machines/software-factory-3.0-state-machine.json`

---

## MANDATORY ACKNOWLEDGMENT

**Every orchestrator startup MUST acknowledge:**

```
I acknowledge R517: Universal State Manager Consultation Law
- I will NEVER update state files directly
- I will ALWAYS use State Manager bookend pattern
- I will ACCEPT State Manager's decisions as FINAL
- I understand bypass = -100% IMMEDIATE FAILURE
- I understand this protects system integrity
```

**Failure to acknowledge = Immediate stop and rule reading required.**

---

## PERMANENT AND UNIVERSAL

**This rule is NOW and FOREVER:**
- ✅ Universal across ALL orchestrator states
- ✅ No state-specific overrides allowed
- ✅ Part of core bootstrap rules
- ✅ Enforced by multiple mechanisms
- ✅ Validated by all tests
- ✅ Impossible to ignore or bypass

**Never again should we wonder if orchestrator is consulting State Manager.**

**It is LAW.**

---

*Rule R517 - Universal State Manager Consultation Law*
*Created: 2025-10-15*
*Criticality: BLOCKING (🚨🚨🚨)*
*Enforcement: Multiple Layers*
*Penalty: -100% Immediate Failure*
