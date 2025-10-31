# 🔴🔴🔴 RULE R550: Formalized Just-In-Time Planning (SUPREME LAW)

**Criticality:** SUPREME LAW - BLOCKING
**Status:** Active
**Created:** 2025-10-31
**Penalty:** -100% for skipping planning

---

## THE PROBLEM THIS SOLVES

### Before R550:
- Just-in-time planning was documented in R504 but not enforced by state machine
- State Manager had no way to guarantee planning happened before infrastructure
- COMPLETE_PHASE could transition to SETUP_PHASE_INFRASTRUCTURE even without planning
- ERROR_RECOVERY couldn't transition to planning states to fix missing plans
- Result: Cascading failures when phases/waves started without planning

### After R550:
- ✅ Mandatory sequences enforce planning → infrastructure flow
- ✅ Guard conditions prevent infrastructure setup without planning
- ✅ State Manager can choose planning states with confidence
- ✅ ERROR_RECOVERY can transition to planning for recovery
- ✅ System GUARANTEES planning happens before execution

---

## THE FORMALIZED SYSTEM

### 1. MANDATORY SEQUENCES

#### Phase Transition Sequence (New):
```
COMPLETE_PHASE
    ↓ (mandatory if: more_phases AND no_planning)
SPAWN_ARCHITECT_PHASE_PLANNING
    ↓
WAITING_FOR_PHASE_ARCHITECTURE
    ↓
SPAWN_CODE_REVIEWER_PHASE_TEST_PLANNING
    ↓
WAITING_FOR_PHASE_TEST_PLAN
    ↓
CREATE_PHASE_INTEGRATION_BRANCH_EARLY
    ↓
SPAWN_CODE_REVIEWER_PHASE_IMPL
    ↓
WAITING_FOR_PHASE_IMPLEMENTATION_PLAN
    ↓
WAVE_START
```

**Enforcement:** BLOCKING - Cannot skip any state
**Trigger:** `more_phases_in_project == true AND next_phase_planning_missing == true`
**Exit:** Only to ERROR_RECOVERY allowed

#### Wave Transition Sequence (New):
```
COMPLETE_WAVE
    ↓ (mandatory if: more_waves AND no_planning)
SPAWN_ARCHITECT_WAVE_PLANNING
    ↓
WAITING_FOR_WAVE_ARCHITECTURE
    ↓
SPAWN_CODE_REVIEWER_WAVE_TEST_PLANNING
    ↓
WAITING_FOR_WAVE_TEST_PLAN
    ↓
ANALYZE_IMPLEMENTATION_PARALLELIZATION
    ↓
WAITING_FOR_IMPLEMENTATION_ANALYSIS
```

**Enforcement:** BLOCKING - Cannot skip any state
**Trigger:** `more_waves_in_phase == true AND next_wave_planning_missing == true`
**Exit:** Only to ERROR_RECOVERY allowed

---

### 2. GUARD CONDITIONS

#### COMPLETE_PHASE Guards:
```json
{
  "guards": {
    "SPAWN_ARCHITECT_PHASE_PLANNING": "more_phases_in_project == true AND next_phase_planning_missing == true",
    "SETUP_PHASE_INFRASTRUCTURE": "more_phases_in_project == true AND next_phase_planning_exists == true",
    "SETUP_PROJECT_INFRASTRUCTURE": "all_phases_complete == true"
  }
}
```

**Logic:**
- If more phases exist AND next phase has no planning → MUST do planning first
- If more phases exist AND next phase has planning → Can setup infrastructure
- If all phases done → Setup project infrastructure

#### COMPLETE_WAVE Guards:
```json
{
  "guards": {
    "SPAWN_ARCHITECT_WAVE_PLANNING": "more_waves_in_phase == true AND next_wave_planning_missing == true",
    "SETUP_WAVE_INFRASTRUCTURE": "more_waves_in_phase == true AND next_wave_planning_exists == true",
    "SETUP_PHASE_INFRASTRUCTURE": "all_waves_complete_in_phase == true"
  }
}
```

**Logic:**
- If more waves exist AND next wave has no planning → MUST do planning first
- If more waves exist AND next wave has planning → Can setup infrastructure
- If all waves done → Setup phase infrastructure

---

### 3. INFRASTRUCTURE STATE REQUIREMENTS

#### SETUP_PHASE_INFRASTRUCTURE:
```json
{
  "requires": {
    "conditions": [
      "Phase code review clean (bugs_found == 0)",
      "Architect approved phase",
      "All phase tests passing",
      "phase_planning_exists OR phase_is_first_phase"
    ]
  }
}
```

**New requirement:** Cannot setup phase infrastructure unless:
- Planning documents exist for the phase, OR
- This is the first phase (planned during project initialization)

#### SETUP_WAVE_INFRASTRUCTURE:
```json
{
  "requires": {
    "conditions": [
      "Code review clean (bugs_found == 0)",
      "Architect approved",
      "All tests passing",
      "wave_planning_exists OR wave_is_first_wave_of_first_phase"
    ]
  }
}
```

**New requirement:** Cannot setup wave infrastructure unless:
- Planning documents exist for the wave, OR
- This is first wave of first phase (planned during project initialization)

---

### 4. ERROR RECOVERY ENHANCEMENTS

#### ERROR_RECOVERY Transitions (Updated):
```json
{
  "allowed_transitions": [
    "SETUP_WAVE_INFRASTRUCTURE",
    "START_WAVE_ITERATION",
    "SETUP_PHASE_INFRASTRUCTURE",
    "START_PHASE_ITERATION",
    "SETUP_PROJECT_INFRASTRUCTURE",
    "START_PROJECT_ITERATION",
    "PROJECT_DONE",
    "SPAWN_ARCHITECT_PHASE_PLANNING",  // ← NEW
    "SPAWN_ARCHITECT_WAVE_PLANNING"    // ← NEW
  ]
}
```

**Now ERROR_RECOVERY can:**
- Transition to planning states to fix missing plans
- Recover from "no planning" errors properly
- Break out of stuck loops by doing planning first

---

## HOW STATE MANAGER USES THIS

### Decision Algorithm (Enhanced):

```python
def choose_next_state(current_state, work_results):
    # 1. Check if in mandatory sequence
    if in_mandatory_sequence(current_state):
        return next_state_in_sequence()

    # 2. Evaluate all allowed transitions
    allowed = get_allowed_transitions(current_state)

    # 3. For each transition, check guards
    valid_transitions = []
    for next_state in allowed:
        guard = get_guard_condition(current_state, next_state)
        if evaluate_guard(guard, work_results):
            valid_transitions.append(next_state)

    # 4. Prioritize planning over infrastructure
    if "SPAWN_ARCHITECT_PHASE_PLANNING" in valid_transitions:
        if not phase_planning_exists(next_phase):
            return "SPAWN_ARCHITECT_PHASE_PLANNING"  # Planning first!

    if "SPAWN_ARCHITECT_WAVE_PLANNING" in valid_transitions:
        if not wave_planning_exists(next_wave):
            return "SPAWN_ARCHITECT_WAVE_PLANNING"  # Planning first!

    # 5. Check requires conditions
    for next_state in valid_transitions:
        if meets_requirements(next_state, work_results):
            return next_state

    # 6. If no valid transition, go to ERROR_RECOVERY
    return "ERROR_RECOVERY"
```

**Key behaviors:**
- Mandatory sequences take precedence
- Planning states prioritized over infrastructure
- Guard conditions prevent wrong transitions
- Requirements checked before allowing state
- ERROR_RECOVERY as fallback

---

## PRACTICAL EXAMPLES

### Example 1: Phase 1 Complete → Phase 2 Start

**Scenario:** Phase 1 integration done, Phase 2 has no planning

```
Current: COMPLETE_PHASE (Phase 1)
Orchestrator proposes: SETUP_PHASE_INFRASTRUCTURE (Phase 2)

State Manager evaluation:
1. Check guards for SETUP_PHASE_INFRASTRUCTURE
   → Guard: "next_phase_planning_exists == true"
   → Evaluate: Phase 2 has no planning documents
   → Result: FALSE - Guard blocks transition

2. Check guards for SPAWN_ARCHITECT_PHASE_PLANNING
   → Guard: "next_phase_planning_missing == true"
   → Evaluate: Phase 2 has no planning documents
   → Result: TRUE - Guard allows transition

3. Check mandatory sequence trigger
   → Trigger: "more_phases AND next_phase_planning_missing"
   → Result: TRUE - Enters mandatory sequence

Decision: SPAWN_ARCHITECT_PHASE_PLANNING (mandatory)
Mandatory Sequence: phase_transition_to_next_phase activated
```

**Orchestrator must:**
- Transition to SPAWN_ARCHITECT_PHASE_PLANNING
- Follow entire sequence through WAVE_START
- Cannot skip planning or jump ahead

---

### Example 2: Wave 1 Complete → Wave 2 Start

**Scenario:** Wave 1 integrated, Wave 2 needs planning

```
Current: COMPLETE_WAVE (Wave 1)
Orchestrator proposes: SETUP_WAVE_INFRASTRUCTURE (Wave 2)

State Manager evaluation:
1. Check guards for SETUP_WAVE_INFRASTRUCTURE
   → Guard: "next_wave_planning_exists == true"
   → Evaluate: Wave 2 has no planning documents
   → Result: FALSE - Guard blocks transition

2. Check guards for SPAWN_ARCHITECT_WAVE_PLANNING
   → Guard: "next_wave_planning_missing == true"
   → Evaluate: Wave 2 has no planning documents
   → Result: TRUE - Guard allows transition

3. Check mandatory sequence trigger
   → Trigger: "more_waves AND next_wave_planning_missing"
   → Result: TRUE - Enters mandatory sequence

Decision: SPAWN_ARCHITECT_WAVE_PLANNING (mandatory)
Mandatory Sequence: wave_transition_to_next_wave activated
```

---

### Example 3: ERROR_RECOVERY → Planning Recovery

**Scenario:** System in ERROR_RECOVERY, discovered Phase 2 has no planning

```
Current: ERROR_RECOVERY
Problem: Phase 2 was attempted without planning
Recovery needed: Create Phase 2 plan

State Manager evaluation:
1. Check allowed transitions from ERROR_RECOVERY
   → SPAWN_ARCHITECT_PHASE_PLANNING now in list (R550 added it!)

2. Orchestrator proposes: SPAWN_ARCHITECT_PHASE_PLANNING

3. State Manager validates:
   → Transition allowed: YES
   → Phase needs planning: YES
   → Can proceed: YES

Decision: SPAWN_ARCHITECT_PHASE_PLANNING
Result: Recovery to planning state successful
```

**Before R550:** This was IMPOSSIBLE - ERROR_RECOVERY couldn't transition to planning!
**After R550:** Clean recovery path available

---

## VALIDATION MECHANISMS

### 1. Pre-Transition Checks (State Manager):

```bash
validate_infrastructure_transition() {
    local next_state="$1"
    local target_phase="$2"
    local target_wave="$3"

    if [[ "$next_state" == "SETUP_PHASE_INFRASTRUCTURE" ]]; then
        # Check if phase planning exists
        if ! phase_planning_exists "$target_phase"; then
            echo "❌ BLOCKED: Phase $target_phase has no planning documents"
            echo "Required: SPAWN_ARCHITECT_PHASE_PLANNING first"
            return 1
        fi
    fi

    if [[ "$next_state" == "SETUP_WAVE_INFRASTRUCTURE" ]]; then
        # Check if wave planning exists
        if ! wave_planning_exists "$target_phase" "$target_wave"; then
            echo "❌ BLOCKED: Phase $target_phase Wave $target_wave has no planning"
            echo "Required: SPAWN_ARCHITECT_WAVE_PLANNING first"
            return 1
        fi
    fi

    return 0
}
```

### 2. Planning Existence Checks:

```bash
phase_planning_exists() {
    local phase="$1"

    # Check for required planning documents
    local planning_dir="planning/phase${phase}"

    [[ -f "$planning_dir/PHASE-ARCHITECTURE-PLAN.md" ]] || return 1
    [[ -f "$planning_dir/PHASE-TEST-PLAN.md" ]] || return 1

    # Check for efforts in pre_planned_infrastructure
    local effort_count=$(jq "[.pre_planned_infrastructure.efforts |
                             to_entries[] |
                             select(.key | startswith(\"phase${phase}_\"))] |
                             length" orchestrator-state-v3.json)

    [[ $effort_count -gt 0 ]] || return 1

    return 0
}

wave_planning_exists() {
    local phase="$1"
    local wave="$2"

    # Check for wave planning documents
    local planning_dir="planning/phase${phase}/wave${wave}"

    [[ -f "$planning_dir/WAVE-IMPLEMENTATION-PLAN.md" ]] || return 1
    [[ -f "$planning_dir/WAVE-TEST-PLAN.md" ]] || return 1

    return 0
}
```

---

## INTEGRATION WITH R504

### R504 (Pre-Infrastructure Planning):
- Defines WHAT planning must contain
- Requires pre-calculation of all infrastructure
- Specifies planning happens during planning states

### R550 (This Rule):
- Defines HOW planning is enforced by state machine
- GUARANTEES planning happens before infrastructure
- Makes R504 requirements MANDATORY through sequences

**Relationship:**
- R504 = Requirements (what must be planned)
- R550 = Enforcement (how system ensures it happens)

---

## GRADING IMPACT

### Compliance (100%):
- ✅ Planning states executed before infrastructure (40%)
- ✅ Mandatory sequences followed completely (30%)
- ✅ Guard conditions respected (20%)
- ✅ Planning documents exist before setup (10%)

### Violations (-100%):
- ❌ Skipping SPAWN_ARCHITECT_PHASE_PLANNING when required
- ❌ Skipping SPAWN_ARCHITECT_WAVE_PLANNING when required
- ❌ Setting up infrastructure without planning documents
- ❌ Bypassing mandatory sequences
- ❌ Ignoring guard conditions

---

## BENEFITS

### Before R550 (Broken):
```
COMPLETE_PHASE → SETUP_PHASE_INFRASTRUCTURE
                    ↓ (no planning!)
                 START_PHASE_ITERATION
                    ↓
                 ERROR_RECOVERY (detected problem too late)
```

### After R550 (Working):
```
COMPLETE_PHASE → SPAWN_ARCHITECT_PHASE_PLANNING (mandatory!)
                    ↓
                 WAITING_FOR_PHASE_ARCHITECTURE
                    ↓
                 [Complete planning sequence]
                    ↓
                 SETUP_PHASE_INFRASTRUCTURE (now has planning!)
                    ↓
                 START_PHASE_ITERATION (smooth execution)
```

---

## MIGRATION NOTES

### For Existing Projects:

If your project is stuck in ERROR_RECOVERY due to missing planning:

1. **Update state machine** with R550 changes
2. **Transition to planning**: ERROR_RECOVERY → SPAWN_ARCHITECT_PHASE_PLANNING
3. **Complete planning sequence**: Follow mandatory sequence through WAVE_START
4. **Continue normally**: System now has proper planning

### For New Projects:

- R550 is automatically enforced from INIT
- Just follow the state machine normally
- System GUARANTEES planning happens at right time

---

## RELATED RULES

- **R504**: Pre-Infrastructure Planning (requirements)
- **R360**: Just-In-Time Infrastructure Execution (timing)
- **R517**: State Manager Consultation Law (enforcement)
- **R206**: State Machine Validation (structure)
- **R234**: Mandatory State Traversal (sequences)

---

## SUMMARY

**R550 formalizes just-in-time planning by:**

1. ✅ Creating mandatory sequences for phase/wave transitions
2. ✅ Adding guard conditions that check for planning
3. ✅ Adding planning transitions to completion and recovery states
4. ✅ Requiring planning documents before infrastructure setup
5. ✅ Making State Manager enforce planning-before-infrastructure

**Result:** System GUARANTEES planning happens. No more cascading failures from missing planning.

**User confidence:** "I am confident that the planning will be chosen" - YES, now you can be!

---

*Rule R550 - Formalized Just-In-Time Planning*
*Created: 2025-10-31*
*Criticality: SUPREME LAW*
*Enforcement: Mandatory Sequences + Guard Conditions*
*Penalty: -100% for violations*
