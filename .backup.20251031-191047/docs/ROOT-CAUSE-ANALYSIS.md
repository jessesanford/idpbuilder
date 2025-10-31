# ROOT CAUSE ANALYSIS: Wave 2 Infrastructure Creation Bypass

**Date**: 2025-10-29
**Issue**: CREATE_NEXT_INFRASTRUCTURE and VALIDATE_INFRASTRUCTURE states skipped for Wave 2
**Impact**: 3 out of 4 SW Engineers blocked due to missing Git branch infrastructure
**Severity**: CRITICAL - Complete wave execution failure

---

## EXECUTIVE SUMMARY

The Software Factory 3.0 state machine allowed an **INVALID TRANSITION** from `WAITING_FOR_EFFORT_PLANS` directly to `SPAWN_SW_ENGINEERS`, bypassing the mandatory infrastructure creation states (`CREATE_NEXT_INFRASTRUCTURE` and `VALIDATE_INFRASTRUCTURE`). This caused 75% of Wave 2 SW Engineers to fail immediately upon spawn.

**Root Cause Category**: **State Machine Definition Error**

The state machine JSON incorrectly lists `SPAWN_SW_ENGINEERS` as an allowed transition from `WAITING_FOR_EFFORT_PLANS` without any guard conditions, despite the state rules clearly mandating that `ANALYZE_IMPLEMENTATION_PARALLELIZATION` must come first.

---

## INVESTIGATION FINDINGS

### 1. ACTUAL vs EXPECTED STATE FLOW

#### What SHOULD Have Happened (Wave 2 - New Infrastructure Needed):
```
INJECT_WAVE_METADATA
  ↓
ANALYZE_CODE_REVIEWER_PARALLELIZATION
  ↓
SPAWN_CODE_REVIEWERS_EFFORT_PLANNING
  ↓
WAITING_FOR_EFFORT_PLANS
  ↓
ANALYZE_IMPLEMENTATION_PARALLELIZATION ← MANDATORY per state rules
  ↓
CREATE_NEXT_INFRASTRUCTURE              ← Create Wave 2 branches!
  ↓
VALIDATE_INFRASTRUCTURE                 ← Validate branches exist!
  ↓
SPAWN_SW_ENGINEERS                      ← NOW spawn agents
```

#### What ACTUALLY Happened (Wave 2):
```
INJECT_WAVE_METADATA (06:33:12)
  ↓
ANALYZE_CODE_REVIEWER_PARALLELIZATION (06:33:12)
  ↓
SPAWN_CODE_REVIEWERS_EFFORT_PLANNING (06:53:40)
  ↓
WAITING_FOR_EFFORT_PLANS (06:53:40)
  ↓
[MISSING: ANALYZE_IMPLEMENTATION_PARALLELIZATION]
  ↓
[MISSING: CREATE_NEXT_INFRASTRUCTURE]
  ↓
[MISSING: VALIDATE_INFRASTRUCTURE]
  ↓
SPAWN_SW_ENGINEERS (07:02:00)           ← Spawned WITHOUT infrastructure!
```

**Result**: Wave 2 directories exist but Wave 2 Git branches DO NOT exist.

### 2. STATE MACHINE DEFINITION FLAW

**File**: `/home/vscode/workspaces/idpbuilder-oci-push-planning/state-machines/software-factory-3.0-state-machine.json`

**Line 1624-1653**: `WAITING_FOR_EFFORT_PLANS` state definition

```json
"WAITING_FOR_EFFORT_PLANS": {
  "description": "Monitor Code Reviewer creating effort-level implementation plans with R340 quality validation",
  "agent": "orchestrator",
  "checkpoint": false,
  "iteration_level": "wave",
  "iteration_type": "MONITORING",
  "allowed_transitions": [
    "ANALYZE_IMPLEMENTATION_PARALLELIZATION",
    "SPAWN_SW_ENGINEERS",                          ← ❌ INVALID WITHOUT GUARDS!
    "ERROR_RECOVERY"
  ],
  "actions": [
    "Apply R356 optimization: single effort → SPAWN_SW_ENGINEERS, multiple → ANALYZE_IMPLEMENTATION_PARALLELIZATION"
  ],
  "guards": {
    "SPAWN_SW_ENGINEERS": "effort_count == 1 (R356 optimization)",
    "ANALYZE_IMPLEMENTATION_PARALLELIZATION": "effort_count > 1 (parallelization analysis required)"
  }
}
```

### 3. THE FATAL FLAW

**The problem**: The `guards` section defines when each transition should be used, BUT:

1. **Guards are DOCUMENTATION ONLY** - They are not enforced by the state machine!
2. The state machine allows BOTH transitions without actual validation
3. R356 optimization is ONLY for implementation planning REUSE, not infrastructure!

**What R356 Actually Says**:
- Single effort → Skip parallelization ANALYSIS (not infrastructure!)
- Multiple efforts → Perform parallelization analysis

**What the orchestrator did**:
- Saw `SPAWN_SW_ENGINEERS` in `allowed_transitions`
- Assumed infrastructure already existed (Wave 1 pattern)
- Skipped infrastructure creation entirely

### 4. WAVE 1 vs WAVE 2 COMPARISON

#### Wave 1 (Worked Correctly):
```
State History Lines 386-421:
WAITING_FOR_WAVE_IMPLEMENTATION_PLAN (02:45:00)
  ↓
CREATE_NEXT_INFRASTRUCTURE (02:51:30)     ← Created 4 effort branches
  ↓
VALIDATE_INFRASTRUCTURE (03:08:45)        ← Validated branches exist
  ↓
[LOOP: Created efforts 1.1.1, 1.1.2, 1.1.3, 1.1.4]
  ↓
VALIDATE_INFRASTRUCTURE (03:45:01)
  ↓
SPAWN_CODE_REVIEWERS_EFFORT_PLANNING (03:45:01)
```

**Key Difference**: Wave 1 went through a different path because it came from `WAITING_FOR_WAVE_IMPLEMENTATION_PLAN`, which has different `allowed_transitions`:

```json
"WAITING_FOR_WAVE_IMPLEMENTATION_PLAN": {
  "allowed_transitions": [
    "CREATE_NEXT_INFRASTRUCTURE",  ← Only valid transition
    "ERROR_RECOVERY"
  ]
}
```

#### Wave 2 (Failed):
```
State History Lines 1204-1226:
WAITING_FOR_EFFORT_PLANS (06:53:40)
  ↓
[Infrastructure creation completely skipped!]
  ↓
SPAWN_SW_ENGINEERS (07:02:00)
```

**Why Different**:
- Wave 1 used **Code Reviewer for wave implementation planning** (old pattern)
- Wave 2 used **Code Reviewer for effort-level planning** (new SF 3.0 pattern)
- Wave 2 path has `SPAWN_SW_ENGINEERS` shortcut that bypasses infrastructure

### 5. INFRASTRUCTURE STATUS VALIDATION

#### pre_planned_infrastructure Analysis:

**File**: `orchestrator-state-v3.json` lines 1407-1502

```json
"pre_planned_infrastructure": {
  "validated": true,
  "efforts": {
    "phase1_wave1_effort-1-docker-interface": { "created": true },
    "phase1_wave1_effort-2-registry-interface": { "created": true },
    "phase1_wave1_effort-3-auth-tls-interfaces": { "created": true },
    "phase1_wave1_effort-4-command-structure": { "created": true }
    // ❌ NO WAVE 2 EFFORTS PRESENT!
  }
}
```

**Missing**:
- `phase1_wave2_effort-1-docker-client`
- `phase1_wave2_effort-2-registry-client`
- `phase1_wave2_effort-3-auth`
- `phase1_wave2_effort-4-tls`

**Conclusion**: Infrastructure metadata was NEVER injected for Wave 2!

---

## ROOT CAUSE IDENTIFICATION

### Primary Root Cause: **STATE MACHINE DEFINITION ERROR**

**Location**: `state-machines/software-factory-3.0-state-machine.json` line 1632

**Issue**: `WAITING_FOR_EFFORT_PLANS.allowed_transitions` includes `SPAWN_SW_ENGINEERS` without proper guard enforcement

**Why This Happened**:

1. **R356 Misapplication**: R356 optimization was designed for **Wave 1** where:
   - Infrastructure already created by `CREATE_NEXT_INFRASTRUCTURE` loop
   - Effort plans created AFTER infrastructure exists
   - Single effort can skip parallelization ANALYSIS

2. **SF 3.0 Pattern Change**: Wave 2+ uses **effort-level planning BEFORE infrastructure**:
   - Code Reviewers create effort plans FIRST
   - Infrastructure should be created AFTER planning
   - R356 no longer applies because infrastructure doesn't exist yet!

3. **Guard Condition Not Enforced**: The state machine treats `guards` as documentation, not enforcement

### Secondary Root Cause: **MISSING INFRASTRUCTURE INJECTION**

**Location**: Between `WAITING_FOR_EFFORT_PLANS` and infrastructure creation

**Issue**: No mechanism to ensure Wave 2 efforts are injected into `pre_planned_infrastructure`

**Expected**: After effort plans are created, infrastructure metadata should be injected

**Actual**: `pre_planned_infrastructure.efforts` only contains Wave 1

### Contributing Factor: **STATE RULES vs STATE MACHINE MISMATCH**

**State Rules File**: `agent-states/software-factory/orchestrator/WAITING_FOR_EFFORT_PLANS/rules.md` lines 265-284

```markdown
## 🔴🔴🔴 SUPREME LAW R234 - STAY IN SEQUENCE 🔴🔴🔴

YOUR POSITION IN THE MANDATORY SEQUENCE:
CREATE_NEXT_INFRASTRUCTURE (✓ completed)
    ↓
ANALYZE_CODE_REVIEWER_PARALLELIZATION (✓ completed)
    ↓
SPAWN_CODE_REVIEWERS_EFFORT_PLANNING (✓ completed)
    ↓
WAITING_FOR_EFFORT_PLANS (👈 YOU ARE HERE)
    ↓ (MUST GO HERE NEXT)
ANALYZE_IMPLEMENTATION_PARALLELIZATION  ← MANDATORY!
    ↓
SPAWN_SW_ENGINEERS

FORBIDDEN: Skipping analysis to go directly to SPAWN_SW_ENGINEERS = -100%
```

**The Contradiction**:
- **State rules SAY**: MUST go to `ANALYZE_IMPLEMENTATION_PARALLELIZATION`
- **State machine ALLOWS**: Direct transition to `SPAWN_SW_ENGINEERS`
- **Orchestrator behavior**: Follows state machine, not state rules

---

## ANALYSIS: WHY ORCHESTRATOR MADE THIS CHOICE

### Orchestrator's Reasoning (Incorrect):

1. Saw `SPAWN_SW_ENGINEERS` in `allowed_transitions` for `WAITING_FOR_EFFORT_PLANS`
2. Checked `guards`: "effort_count == 1 → SPAWN_SW_ENGINEERS" (R356 optimization)
3. **Misunderstood R356**: Thought "skip parallelization" meant "skip ALL intermediate states"
4. **Assumed infrastructure exists**: Pattern matching from Wave 1 where infrastructure was pre-created
5. Transitioned directly to `SPAWN_SW_ENGINEERS`

### What Orchestrator SHOULD Have Done:

1. Recognize: Wave 2 needs NEW infrastructure (not Wave 1 infrastructure!)
2. Check `pre_planned_infrastructure.efforts` for Wave 2 entries → NOT FOUND
3. Conclude: MUST create infrastructure before spawning
4. Transition to `ANALYZE_IMPLEMENTATION_PARALLELIZATION` (validates infrastructure needs)
5. Then to `CREATE_NEXT_INFRASTRUCTURE` → `VALIDATE_INFRASTRUCTURE` → `SPAWN_SW_ENGINEERS`

---

## COMPLETE FAILURE CHAIN

```
1. State Machine Definition Flaw
   └─> WAITING_FOR_EFFORT_PLANS allows SPAWN_SW_ENGINEERS
       └─> R356 guard documented but not enforced
           └─> Orchestrator follows state machine literal definition

2. Missing Infrastructure Metadata
   └─> pre_planned_infrastructure only has Wave 1 efforts
       └─> Wave 2 efforts never injected after planning
           └─> No validation check before spawning

3. No Pre-Spawn Validation
   └─> SPAWN_SW_ENGINEERS state does not verify infrastructure exists
       └─> Agents spawned into non-existent branches
           └─> 3 out of 4 agents blocked immediately

4. R151 Timing Violation
   └─> Sequential spawn due to failures (not true parallel)
       └─> 10m 24s delta between spawns
           └─> -50% grading penalty
```

---

## CORRECT SEQUENCE FOR WAVE 2

### What SHOULD Happen After Effort Planning:

```
Step 1: WAITING_FOR_EFFORT_PLANS
  - Code Reviewers complete effort plans
  - Plans tracked in orchestrator-state-v3.json
  - Transition trigger: All plans created

Step 2: ANALYZE_IMPLEMENTATION_PARALLELIZATION (MANDATORY)
  - Read effort plans to extract R213 metadata
  - Determine parallelization strategy
  - Identify infrastructure requirements
  - **CRITICAL**: Detect that Wave 2 infrastructure doesn't exist

Step 3: CREATE_NEXT_INFRASTRUCTURE
  - Inject Wave 2 efforts into pre_planned_infrastructure
  - Create Git branches for each effort:
    - idpbuilder-oci-push/phase1/wave2/effort-1-docker-client
    - idpbuilder-oci-push/phase1/wave2/effort-2-registry-client
    - idpbuilder-oci-push/phase1/wave2/effort-3-auth
    - idpbuilder-oci-push/phase1/wave2/effort-4-tls
  - Push branches to remote with -u flag
  - Update pre_planned_infrastructure.efforts[*].created = true

Step 4: VALIDATE_INFRASTRUCTURE
  - Verify all branches exist on remote
  - Verify upstream tracking configured
  - Verify directory structure correct
  - Mark pre_planned_infrastructure.efforts[*].validated = true

Step 5: SPAWN_SW_ENGINEERS
  - NOW spawn agents into validated infrastructure
  - All 4 agents can proceed immediately
  - R151 timing compliance achieved
```

---

## FIXES REQUIRED

### Fix #1: State Machine JSON Correction (PRIMARY FIX)

**File**: `state-machines/software-factory-3.0-state-machine.json`

**Change**: Remove `SPAWN_SW_ENGINEERS` from `WAITING_FOR_EFFORT_PLANS.allowed_transitions`

**Before** (lines 1630-1633):
```json
"allowed_transitions": [
  "ANALYZE_IMPLEMENTATION_PARALLELIZATION",
  "SPAWN_SW_ENGINEERS",    ← REMOVE THIS!
  "ERROR_RECOVERY"
]
```

**After**:
```json
"allowed_transitions": [
  "ANALYZE_IMPLEMENTATION_PARALLELIZATION",
  "ERROR_RECOVERY"
]
```

**Rationale**:
- R356 optimization applies to parallelization ANALYSIS, not infrastructure creation
- ALL waves must go through infrastructure creation after effort planning
- State machine should enforce mandatory sequence

### Fix #2: Add Infrastructure Validation to SPAWN_SW_ENGINEERS

**File**: `state-machines/software-factory-3.0-state-machine.json`

**Change**: Add `requires.conditions` to SPAWN_SW_ENGINEERS state

**Before** (around line 1251):
```json
"SPAWN_SW_ENGINEERS": {
  "description": "Spawn SW Engineer agents to implement features per effort plans",
  "requires": {
    "conditions": [
      "Code review analysis complete",
      "Wave efforts identified"
    ]
  }
}
```

**After**:
```json
"SPAWN_SW_ENGINEERS": {
  "description": "Spawn SW Engineer agents to implement features per effort plans",
  "requires": {
    "conditions": [
      "Code review analysis complete",
      "Wave efforts identified",
      "Infrastructure validated (pre_planned_infrastructure.efforts[*].validated == true)",
      "Git branches exist for all efforts",
      "Workspace directories created"
    ]
  }
}
```

### Fix #3: Update WAITING_FOR_EFFORT_PLANS State Rules

**File**: `agent-states/software-factory/orchestrator/WAITING_FOR_EFFORT_PLANS/rules.md`

**Change**: Clarify that R356 does NOT apply to infrastructure creation

**Add After Line 284**:
```markdown
## 🚨 CRITICAL: R356 DOES NOT SKIP INFRASTRUCTURE! 🚨

**R356 Optimization Scope**:
- ✅ APPLIES TO: Parallelization analysis complexity
  - Single effort → Simple spawn decision
  - Multiple efforts → Complex dependency analysis

- ❌ DOES NOT APPLY TO: Infrastructure creation
  - ALL waves need infrastructure created
  - ALL efforts need Git branches
  - Infrastructure is NEVER optional

**Correct Interpretation**:
- R356 means: "Skip complex parallelization analysis for single effort"
- R356 does NOT mean: "Skip infrastructure creation"
- R356 does NOT mean: "Skip ANALYZE_IMPLEMENTATION_PARALLELIZATION entirely"

**WAITING_FOR_EFFORT_PLANS MUST ALWAYS GO TO**:
- ANALYZE_IMPLEMENTATION_PARALLELIZATION (even for single effort!)
- Which then determines if infrastructure exists or needs creation
```

### Fix #4: Add Infrastructure Metadata Injection Point

**File**: `agent-states/software-factory/orchestrator/ANALYZE_IMPLEMENTATION_PARALLELIZATION/rules.md`

**Change**: Add mandatory check for infrastructure existence

**Add to "Mandatory Analysis Protocol" section**:
```markdown
### Step 0.5: Verify Infrastructure Exists or Plan Creation

```bash
echo "🔍 Checking if infrastructure exists for current wave efforts..."

CURRENT_PHASE=$(yq '.current_phase' orchestrator-state-v3.json)
CURRENT_WAVE=$(yq '.current_wave' orchestrator-state-v3.json)

# Check if pre_planned_infrastructure has entries for this wave
WAVE_EFFORTS=$(yq ".pre_planned_infrastructure.efforts | to_entries[] |
  select(.value.phase == \"phase${CURRENT_PHASE}\" and
         .value.wave == \"wave${CURRENT_WAVE}\") |
  .key" orchestrator-state-v3.json | wc -l)

if [ "$WAVE_EFFORTS" -eq 0 ]; then
    echo "❌ No infrastructure found for Phase $CURRENT_PHASE Wave $CURRENT_WAVE"
    echo "Infrastructure must be created before spawning SW Engineers"
    PROPOSED_NEXT_STATE="CREATE_NEXT_INFRASTRUCTURE"
    TRANSITION_REASON="Wave infrastructure needs creation"

    # Inject infrastructure metadata from effort plans
    echo "Injecting effort metadata into pre_planned_infrastructure..."
    # [Run infrastructure injection script]

elif [ "$(yq ".pre_planned_infrastructure.efforts.*.validated" orchestrator-state-v3.json | grep false)" ]; then
    echo "❌ Infrastructure exists but not validated"
    PROPOSED_NEXT_STATE="VALIDATE_INFRASTRUCTURE"
    TRANSITION_REASON="Infrastructure validation required"
else
    echo "✅ Infrastructure validated and ready"
    PROPOSED_NEXT_STATE="SPAWN_SW_ENGINEERS"
    TRANSITION_REASON="Infrastructure validated, ready to spawn"
fi
```
```

---

## IMPLEMENTATION STEPS FOR BOTH REPOS

### Repository 1: Current Project
**Path**: `/home/vscode/workspaces/idpbuilder-oci-push-planning/`

### Repository 2: Template
**Path**: `/home/vscode/software-factory-template/`

### Apply Fixes to BOTH:

1. **Update state machine JSON** (Fix #1)
2. **Add SPAWN_SW_ENGINEERS validation** (Fix #2)
3. **Update WAITING_FOR_EFFORT_PLANS rules** (Fix #3)
4. **Update ANALYZE_IMPLEMENTATION_PARALLELIZATION rules** (Fix #4)

---

## TESTING PROTOCOL

### Test Case 1: Wave 1 Regression Test
**Objective**: Ensure Wave 1 still works correctly

**Steps**:
1. Start new project from scratch
2. Complete Phase 1 Wave 1
3. Verify infrastructure created correctly
4. Verify SW Engineers spawn successfully

**Expected**: No changes to Wave 1 behavior

### Test Case 2: Wave 2 Success Test
**Objective**: Verify Wave 2 now creates infrastructure

**Steps**:
1. Complete Phase 1 Wave 1
2. Start Phase 1 Wave 2
3. Complete effort planning
4. Verify state machine goes:
   - WAITING_FOR_EFFORT_PLANS
   - → ANALYZE_IMPLEMENTATION_PARALLELIZATION
   - → CREATE_NEXT_INFRASTRUCTURE
   - → VALIDATE_INFRASTRUCTURE
   - → SPAWN_SW_ENGINEERS

**Expected**: All 4 SW Engineers spawn successfully with <5s delta

### Test Case 3: Guard Condition Enforcement
**Objective**: Ensure invalid transitions are rejected

**Steps**:
1. At WAITING_FOR_EFFORT_PLANS state
2. Attempt to transition directly to SPAWN_SW_ENGINEERS
3. Verify rejection by state machine

**Expected**: Transition rejected, ERROR_RECOVERY triggered

---

## PREVENTION MEASURES

### 1. State Machine Validation Tool
Create tool to validate:
- All `allowed_transitions` have matching state entries
- All `guards` are actually enforced
- All `requires.conditions` are validated before entry

### 2. Pre-Spawn Infrastructure Check
Add to SPAWN_SW_ENGINEERS state entry validation:
```bash
# Mandatory pre-spawn validation
for effort in $EFFORTS; do
    if ! git ls-remote --heads origin "$(get_branch_name $effort)" | grep -q .; then
        echo "BLOCKING: Infrastructure missing for $effort"
        exit 1
    fi
done
```

### 3. Documentation Sync
Add CI check to ensure:
- State rules match state machine `allowed_transitions`
- Guard conditions in state machine match state rules requirements
- No contradictions between state machine and state rules

---

## GRADING IMPACT ANALYSIS

### Current Failure:
- **Workspace Isolation**: -20% (3/4 agents blocked)
- **R151 Parallelization**: -50% (timing violation)
- **Workflow Compliance**: -25% (infrastructure incomplete)
- **Total**: ~17.5% / 100%

### After Fix:
- **Workspace Isolation**: +20% (all agents work correctly)
- **R151 Parallelization**: +50% (proper parallel spawn)
- **Workflow Compliance**: +25% (infrastructure complete)
- **Total**: ~100% / 100%

---

## CONCLUSION

The root cause is a **state machine definition error** where `WAITING_FOR_EFFORT_PLANS` incorrectly allows direct transition to `SPAWN_SW_ENGINEERS`, bypassing mandatory infrastructure creation states. This was compounded by:

1. Misapplication of R356 optimization
2. Missing infrastructure metadata for Wave 2
3. Lack of pre-spawn validation in SPAWN_SW_ENGINEERS
4. Contradiction between state rules and state machine

The fix is straightforward: Remove the invalid transition, add validation guards, and clarify documentation. This will prevent recurrence for all future waves and phases.

**Critical Takeaway**: Guard conditions in state machine JSON must be ENFORCED, not just documented. The state machine should prevent invalid transitions, not rely on orchestrator to interpret guard documentation correctly.
