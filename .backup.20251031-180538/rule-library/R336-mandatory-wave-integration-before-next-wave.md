# 🔴🔴🔴 SUPREME LAW R336: Mandatory Wave Integration Before Next Wave

## Classification
- **Category**: Core Development Flow / Integration Requirements
- **Criticality Level**: 🔴🔴🔴 SUPREME LAW
- **Enforcement**: MANDATORY - Trunk-based development foundation
- **Penalty**: -100% for skipping wave integration

## The Rule

**EVERY wave MUST be fully integrated and its integration branch created BEFORE the next wave can start. The next wave's efforts MUST use the previous wave's integration branch as their base per R308. This creates TRUE trunk-based development.**

**CRITICAL CLARIFICATION:**
- Single-effort waves STILL require integration (the integration might be trivial but is MANDATORY)
- Single-wave phases STILL require phase integration (may duplicate wave integration but is MANDATORY)
- The number of efforts/waves does NOT affect integration requirements
- Integration and demos are ALWAYS required at wave/phase/project levels

### 🔴 THE MANDATORY INTEGRATE_WAVE_EFFORTS FLOW 🔴

```
Wave N efforts complete
    ↓
WAVE_COMPLETE (verify all reviews passed)
    ↓
INTEGRATE_WAVE_EFFORTS (setup infrastructure)
    ↓
SPAWN_CODE_REVIEWER_MERGE_PLAN (create merge plan)
    ↓
SPAWN_INTEGRATION_AGENT (execute merges)
    ↓
MONITORING_INTEGRATE_WAVE_EFFORTS (verify success)
    ↓
REVIEW_WAVE_ARCHITECTURE (architect reviews integrated wave)
    ↓
WAVE_START (Wave N+1 MUST use wave-N-integration as base!)
```

**NO WAVE MAY START WITHOUT PREVIOUS WAVE'S INTEGRATE_WAVE_EFFORTS BRANCH!**

## Why This Is CRITICAL

### Without Wave Integration (BROKEN):
```
main
  ├─→ P1W1 efforts (from main)
  ├─→ P1W2 efforts (ALSO from main - missing W1 work!)
  └─→ P1W3 efforts (ALSO from main - missing W1 & W2 work!)
```
**Result**: Massive conflicts, broken dependencies, integration nightmare

### With Mandatory Wave Integration (CORRECT):
```
main
  └─→ P1W1 efforts → wave1-integration
                        └─→ P1W2 efforts → wave2-integration
                                              └─→ P1W3 efforts
```
**Result**: Incremental integration, conflicts detected early, true CI/CD

## Mandatory Implementation

### 1. Wave Completion MUST Lead to Integration

**FORBIDDEN State Transitions:**
- ❌ `WAVE_COMPLETE` → `WAVE_START` (skips integration!)
- ❌ `WAVE_COMPLETE` → `PLANNING` (skips integration!)
- ❌ `REVIEW_WAVE_ARCHITECTURE` → `WAVE_START` without verifying integration exists

**REQUIRED State Transitions:**
- ✅ `WAVE_COMPLETE` → `INTEGRATE_WAVE_EFFORTS`
- ✅ `INTEGRATE_WAVE_EFFORTS` → `SPAWN_CODE_REVIEWER_MERGE_PLAN`
- ✅ `MONITORING_INTEGRATE_WAVE_EFFORTS` → `REVIEW_WAVE_ARCHITECTURE`
- ✅ `REVIEW_WAVE_ARCHITECTURE` → `WAVE_START` (only after integration verified)

### 2. Integration Branch Creation is MANDATORY

```bash
# After WAVE_COMPLETE, orchestrator MUST:
create_wave_integration() {
    local PHASE=$1
    local WAVE=$2
    local INTEGRATE_WAVE_EFFORTS_BRANCH="phase${PHASE}-wave${WAVE}-integration"
    
    echo "🔴 R336: Creating MANDATORY wave integration branch"
    
    # This branch becomes the base for next wave per R308!
    echo "📌 This will be the base for Wave $((WAVE + 1))"
    
    # Record in state file
    jq --arg branch "$INTEGRATE_WAVE_EFFORTS_BRANCH" \
       '.wave_integrations.phase'$PHASE'.wave'$WAVE' = {
          "branch": $branch,
          "created_at": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'",
          "is_base_for_next_wave": true
       }' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
}
```

### 3. Next Wave MUST Use Integration as Base (R308 + R336)

```bash
# In CREATE_NEXT_INFRASTRUCTURE state:
determine_effort_base_branch() {
    local PHASE=$1
    local WAVE=$2
    
    # First wave of first phase: from main
    if [[ $PHASE -eq 1 && $WAVE -eq 1 ]]; then
        echo "main"
        return
    fi
    
    # First wave of new phase: from previous phase integration
    if [[ $WAVE -eq 1 ]]; then
        PREV_PHASE=$((PHASE - 1))
        BASE="phase${PREV_PHASE}-integration"
    else
        # SUBSEQUENT WAVES: FROM PREVIOUS WAVE INTEGRATE_WAVE_EFFORTS (R336!)
        PREV_WAVE=$((WAVE - 1))
        BASE="phase${PHASE}-wave${PREV_WAVE}-integration"
    fi
    
    # VERIFY INTEGRATE_WAVE_EFFORTS BRANCH EXISTS (R336 ENFORCEMENT)
    if ! git ls-remote --heads origin "$BASE" > /dev/null 2>&1; then
        echo "🔴🔴🔴 R336 VIOLATION: Wave $PREV_WAVE integration not found!"
        echo "Cannot start Wave $WAVE without Wave $PREV_WAVE integration!"
        echo "Required branch: $BASE"
        exit 336  # R336 violation
    fi
    
    echo "$BASE"
}
```

### 4. Architect Must Review Integrated Wave

```bash
# In REVIEW_WAVE_ARCHITECTURE state:
verify_integration_before_review() {
    local PHASE=$1
    local WAVE=$2
    local INTEGRATE_WAVE_EFFORTS_BRANCH="phase${PHASE}-wave${WAVE}-integration"
    
    echo "🔍 R336: Verifying wave integration exists for review"
    
    if ! git ls-remote --heads origin "$INTEGRATE_WAVE_EFFORTS_BRANCH" > /dev/null 2>&1; then
        echo "❌ R336 VIOLATION: No integration branch for review!"
        echo "Cannot review wave without integration!"
        exit 336
    fi
    
    echo "✅ Integration branch verified: $INTEGRATE_WAVE_EFFORTS_BRANCH"
    echo "📋 Architect will review the INTEGRATED wave"
}
```

## Integration with State Machine

### Required State Machine Updates:

1. **Remove Invalid Transitions:**
   - DELETE: `WAVE_COMPLETE → WAVE_START`
   - DELETE: `WAVE_COMPLETE → PLANNING`
   - DELETE: `REVIEW_WAVE_ARCHITECTURE → WAVE_START` (without integration check)

2. **Add Validation States:**
   - ADD: Integration verification before `WAVE_START`
   - ADD: Base branch validation in `CREATE_NEXT_INFRASTRUCTURE`

3. **Update Decision Logic:**
   ```yaml
   # In REVIEW_WAVE_ARCHITECTURE state:
   if decision == "PROCEED_NEXT_WAVE":
       # MUST verify integration exists first!
       verify_wave_integration_exists(phase, wave)
       transition_to("WAVE_START")  # Next wave uses integration as base
   ```

## Verification Protocol

### Pre-Wave Start Verification
```bash
# MANDATORY before starting any wave > 1
echo "🔍 R336: Pre-wave verification..."

PHASE=$(jq '.current_phase' orchestrator-state-v3.json)
WAVE=$(jq '.current_wave' orchestrator-state-v3.json)

if [[ $WAVE -gt 1 ]]; then
    PREV_WAVE=$((WAVE - 1))
    REQUIRED_BASE="phase${PHASE}-wave${PREV_WAVE}-integration"
    
    if ! git ls-remote --heads origin "$REQUIRED_BASE" > /dev/null 2>&1; then
        echo "🔴🔴🔴 R336 VIOLATION DETECTED!"
        echo "Wave $WAVE cannot start without Wave $PREV_WAVE integration!"
        echo "Missing: $REQUIRED_BASE"
        echo "IMMEDIATE STOP REQUIRED"
        exit 336
    fi
    
    echo "✅ R336 verified: Integration branch exists"
fi
```

### Post-Integration Verification
```bash
# After integration completes
echo "📋 R336: Recording integration for next wave base"

jq --arg branch "$INTEGRATE_WAVE_EFFORTS_BRANCH" \
   '.last_completed_integration = $branch' orchestrator-state-v3.json > tmp.json
mv tmp.json orchestrator-state-v3.json

echo "✅ Wave $WAVE integration ready as base for Wave $((WAVE + 1))"
```

## Visual Flow Diagram

```
Phase 1, Wave 1:
main ──→ W1 efforts ──→ wave1-integration ✅
                             │
                             ↓ (R336: MANDATORY BASE)
Phase 1, Wave 2:             │
wave1-integration ──→ W2 efforts ──→ wave2-integration ✅
                                          │
                                          ↓ (R336: MANDATORY BASE)
Phase 1, Wave 3:                          │
wave2-integration ──→ W3 efforts ──→ wave3-integration ✅
                                          │
                                          ↓
                                    phase1-integration
```

## Common Violations

### VIOLATION 1: Starting Wave Without Integration
```bash
# ❌ WRONG: Wave 2 starting without Wave 1 integration
current_state: WAVE_START
current_wave: 2
base_branch: main  # WRONG! Should be wave1-integration!
```

### VIOLATION 2: Skipping Integration State
```bash
# ❌ WRONG: Going directly from WAVE_COMPLETE to next wave
WAVE_COMPLETE → WAVE_START  # NO! Must go through INTEGRATE_WAVE_EFFORTS!
```

### VIOLATION 3: Wrong Base Branch
```bash
# ❌ WRONG: Using main instead of integration
git clone --branch main ...  # For Wave 2+
# ✅ CORRECT: Using previous wave integration
git clone --branch "phase1-wave1-integration" ...
```

## Grading Impact

- **Starting wave without previous integration**: -100% AUTOMATIC FAILURE
- **Using wrong base branch (main instead of integration)**: -100% FAILURE
- **Skipping INTEGRATE_WAVE_EFFORTS state**: -100% FAILURE
- **Not verifying integration exists**: -50% penalty
- **Integration branch not pushed**: -30% penalty
- **Not recording integration in state**: -25% penalty

## Stop Work Conditions

**IMMEDIATE STOP if:**
1. Wave > 1 starting without previous wave integration branch
2. Efforts using main as base when integration exists
3. WAVE_COMPLETE transitioning directly to WAVE_START
4. Integration branch missing when wave review requested
5. Next wave infrastructure created before integration complete

## Integration with Other Rules

- **R308**: Incremental branching (defines base branch logic)
- **R009**: Integration branch creation (defines how to integrate)
- **R285**: Phase integration before assessment (phase-level equivalent)
- **R234**: Mandatory state traversal (prevents skipping states)
- **R222**: Code review gate (ensures reviews before integration)

## Integration with Iteration Containers (SF 3.0)

**R336 wave integration creates iteration containers:**

```
┌──────────────────────────────────────────────┐
│  WAVE 1 INTEGRATE_WAVE_EFFORTS ITERATION CONTAINER      │
├──────────────────────────────────────────────┤
│  R336 mandates this integration exists       │
│  before Wave 2 can start.                    │
│                                              │
│  This container may iterate multiple times:  │
│  - Iteration 1: Merge → Review → Bugs       │
│  - Iteration 2: Fix → Merge → Review → Bugs │
│  - Iteration 3: Fix → Merge → Review → CLEAN│
│                                              │
│  Once CLEAN, wave1-integration branch        │
│  becomes the base for Wave 2 efforts         │
│  (per R308 incremental branching)            │
└──────────────────────────────────────────────┘
```

**Relationship to iteration containers:**
- Each wave integration IS an iteration container
- R336 blocks next wave until container completes (converges to clean)
- This ensures Wave 2 starts with stable, reviewed base
- Prevents cascading bugs across waves

**Container lifecycle:**
1. Wave completes → R336 mandates integration
2. Integration container created
3. Container iterates until clean (tracked in `integration-containers.json`)
4. Container completes → wave integration branch finalized
5. Next wave uses integration branch as base (R308)

**Tracked in**: `integration-containers.json` per SF 3.0 architecture
**Convergence expectation**: Each wave integration should converge (bugs → 0)

## Summary

This rule ensures TRUE trunk-based development where each wave builds incrementally on the integrated work of the previous wave. Integration is not optional - it's the foundation that enables the next wave to succeed.

**The chain of integration is sacred:**
- Wave 1 → wave1-integration (iteration container) → Wave 2 → wave2-integration → Wave 3
- Break the chain = Break the build = FAILURE
- Each integration IS an iteration container expecting convergence

**Remember**: Every wave stands on the shoulders of the integrated wave before it!