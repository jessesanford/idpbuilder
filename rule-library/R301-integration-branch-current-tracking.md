# 🔴🔴🔴 RULE R301: Integration Branch Current Tracking (SUPREME LAW)

## Summary
**ONLY ONE integration branch per level (wave/phase) can be "current" at any time. Using a deprecated integration branch = -100% AUTOMATIC FAILURE.**

## Criticality
**🔴🔴🔴 SUPREME LAW** - Violation results in automatic project failure

## Description

### The Problem This Solves
Multiple integration branches are created during error recovery cycles at BOTH wave and phase levels, causing confusion about which branch is "current" for assessment. This rule ELIMINATES that confusion by enforcing single "current" pointers for both wave and phase integrations.

### Integration Hierarchy
1. **Wave Integrations**: Created after each wave completes (WAVE_COMPLETE state)
2. **Phase Integrations**: Created after all waves in a phase complete (PHASE_INTEGRATION state)

### Mandatory State File Structure

```yaml
# REQUIRED: Current WAVE integration pointer (one per wave)
current_wave_integration:
  phase: 2
  wave: 1
  branch: "phase2-wave1-integration-20250901-143000"
  status: "active"  # MUST be "active"
  created_at: "2025-09-01T14:30:00Z"
  type: "initial"  # "initial" or "post_fixes"
  
# REQUIRED: Current PHASE integration pointer (one per phase)
current_phase_integration:
  phase: 2
  branch: "phase2-integration-20250901-180000"
  status: "active"  # MUST be "active"
  created_at: "2025-09-01T18:00:00Z"
  type: "post_fixes"  # "initial" or "post_fixes"
  
# REQUIRED: Deprecated wave integrations tracking
deprecated_wave_integrations:
  - phase: 2
    wave: 1
    branch: "phase2-wave1-integration-20250901-120000"
    status: "deprecated"  # MUST be "deprecated"
    deprecated_at: "2025-09-01T14:30:00Z"
    reason: "build failures, new attempt needed"
    
# REQUIRED: Deprecated phase integrations tracking
deprecated_phase_integrations:
  - phase: 2
    branch: "phase2-integration-20250901-160000"
    status: "deprecated"  # MUST be "deprecated"
    deprecated_at: "2025-09-01T18:00:00Z"
    reason: "test failures, new attempt needed"
```

## Enforcement Rules

### 1. Creating New WAVE Integration Branch
When creating a wave integration branch (in WAVE_COMPLETE state):
```bash
# MANDATORY: Deprecate existing wave integration for this wave
EXISTING_WAVE=$(yq ".current_wave_integration | select(.phase == $PHASE and .wave == $WAVE)" orchestrator-state.yaml)
if [ ! -z "$EXISTING_WAVE" ]; then
    yq -i ".deprecated_wave_integrations += (.current_wave_integration | select(.phase == $PHASE and .wave == $WAVE))" orchestrator-state.yaml
fi

# Set the NEW current wave integration
yq -i ".current_wave_integration = {
  \"phase\": $PHASE,
  \"wave\": $WAVE,
  \"branch\": \"$NEW_WAVE_BRANCH\",
  \"status\": \"active\",
  \"created_at\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
  \"type\": \"$INTEGRATION_TYPE\"
}" orchestrator-state.yaml
```

### 2. Creating New PHASE Integration Branch
When creating a phase integration branch (in PHASE_INTEGRATION state):
```bash
# MANDATORY: Deprecate existing phase integration for this phase
EXISTING_PHASE=$(yq ".current_phase_integration | select(.phase == $PHASE)" orchestrator-state.yaml)
if [ ! -z "$EXISTING_PHASE" ]; then
    yq -i ".deprecated_phase_integrations += (.current_phase_integration | select(.phase == $PHASE))" orchestrator-state.yaml
fi

# Set the NEW current phase integration
yq -i ".current_phase_integration = {
  \"phase\": $PHASE,
  \"branch\": \"$NEW_PHASE_BRANCH\",
  \"status\": \"active\",
  \"created_at\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
  \"type\": \"$INTEGRATION_TYPE\"
}" orchestrator-state.yaml
```

### 3. Validation Before ANY Integration Work

#### Validate Current WAVE Integration
```bash
validate_current_wave_integration() {
    local PHASE=$1
    local WAVE=$2
    local BRANCH_TO_USE=$3
    
    # Get the CURRENT wave integration branch
    CURRENT=$(yq ".current_wave_integration | select(.phase == $PHASE and .wave == $WAVE).branch" orchestrator-state.yaml)
    CURRENT_STATUS=$(yq ".current_wave_integration | select(.phase == $PHASE and .wave == $WAVE).status" orchestrator-state.yaml)
    
    # FATAL if not using current
    if [[ "$BRANCH_TO_USE" != "$CURRENT" ]]; then
        echo "🔴🔴🔴 SUPREME LAW VIOLATION: R301 - Using deprecated wave integration branch!"
        echo "  Current: $CURRENT (status: $CURRENT_STATUS)"
        echo "  Attempted: $BRANCH_TO_USE"
        echo "  PENALTY: -100% AUTOMATIC FAILURE"
        exit 1
    fi
    
    # FATAL if current is not active
    if [[ "$CURRENT_STATUS" != "active" ]]; then
        echo "🔴🔴🔴 FATAL: Current wave integration is not active!"
        exit 1
    fi
    
    echo "✅ Using current wave integration: $CURRENT"
}
```

#### Validate Current PHASE Integration
```bash
validate_current_phase_integration() {
    local PHASE=$1
    local BRANCH_TO_USE=$2
    
    # Get the CURRENT phase integration branch
    CURRENT=$(yq ".current_phase_integration | select(.phase == $PHASE).branch" orchestrator-state.yaml)
    CURRENT_STATUS=$(yq ".current_phase_integration | select(.phase == $PHASE).status" orchestrator-state.yaml)
    
    # FATAL if not using current
    if [[ "$BRANCH_TO_USE" != "$CURRENT" ]]; then
        echo "🔴🔴🔴 SUPREME LAW VIOLATION: R301 - Using deprecated phase integration branch!"
        echo "  Current: $CURRENT (status: $CURRENT_STATUS)"
        echo "  Attempted: $BRANCH_TO_USE"
        echo "  PENALTY: -100% AUTOMATIC FAILURE"
        exit 1
    fi
    
    # FATAL if current is not active
    if [[ "$CURRENT_STATUS" != "active" ]]; then
        echo "🔴🔴🔴 FATAL: Current phase integration is not active!"
        exit 1
    fi
    
    echo "✅ Using current phase integration: $CURRENT"
}
```

### 4. Architect MUST Use Current Integration

#### For WAVE Reviews (SPAWN_ARCHITECT_WAVE_REVIEW):
```bash
# MANDATORY: Get ONLY the current wave integration
WAVE_BRANCH=$(yq ".current_wave_integration | select(.phase == $PHASE and .wave == $WAVE).branch" orchestrator-state.yaml)

if [ -z "$WAVE_BRANCH" ]; then
    echo "❌ No current wave integration branch for phase $PHASE wave $WAVE!"
    exit 1
fi

# Validate it's active
STATUS=$(yq ".current_wave_integration | select(.phase == $PHASE and .wave == $WAVE).status" orchestrator-state.yaml)
if [ "$STATUS" != "active" ]; then
    echo "🔴🔴🔴 Current wave integration is not active!"
    exit 1
fi

echo "✅ Using current wave integration: $WAVE_BRANCH"
```

#### For PHASE Assessments (SPAWN_ARCHITECT_PHASE_ASSESSMENT):
```bash
# MANDATORY: Get ONLY the current phase integration
PHASE_BRANCH=$(yq ".current_phase_integration | select(.phase == $PHASE).branch" orchestrator-state.yaml)

if [ -z "$PHASE_BRANCH" ]; then
    echo "❌ No current phase integration branch for phase $PHASE!"
    exit 1
fi

# Validate it's active
STATUS=$(yq ".current_phase_integration | select(.phase == $PHASE).status" orchestrator-state.yaml)
if [ "$STATUS" != "active" ]; then
    echo "🔴🔴🔴 Current phase integration is not active!"
    exit 1
fi

echo "✅ Using current phase integration: $PHASE_BRANCH"
```

## Migration for Existing Systems

### For Wave Integration Branches:
```bash
# Migration script for wave integrations
PHASE=$(yq '.current_phase' orchestrator-state.yaml)
WAVE=$(yq '.current_wave' orchestrator-state.yaml)

# Find existing wave integration branches
WAVE_BRANCHES=$(yq ".integration_branches[] | select(.branch | test(\"phase$PHASE.*wave$WAVE.*integration\"))" orchestrator-state.yaml)

if [ ! -z "$WAVE_BRANCHES" ]; then
    # Get the most recent wave integration
    LATEST_WAVE_BRANCH=$(echo "$WAVE_BRANCHES" | yq '.branch' | tail -1)
    
    # Set as current wave integration
    yq -i ".current_wave_integration = {
      \"phase\": $PHASE,
      \"wave\": $WAVE,
      \"branch\": \"$LATEST_WAVE_BRANCH\",
      \"status\": \"active\",
      \"created_at\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
      \"type\": \"migrated\"
    }" orchestrator-state.yaml
fi
```

### For Phase Integration Branches:
```bash
# Migration script for phase integrations
PHASE=$(yq '.current_phase' orchestrator-state.yaml)

# Find existing phase integration branches (NOT wave integrations)
PHASE_BRANCHES=$(yq ".phase_integration_branches[] | select(.phase == $PHASE)" orchestrator-state.yaml)

if [ ! -z "$PHASE_BRANCHES" ]; then
    # Get the most recent phase integration
    LATEST_PHASE_BRANCH=$(echo "$PHASE_BRANCHES" | yq '.branch' | tail -1)
    
    # Set as current phase integration
    yq -i ".current_phase_integration = {
      \"phase\": $PHASE,
      \"branch\": \"$LATEST_PHASE_BRANCH\",
      \"status\": \"active\",
      \"created_at\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
      \"type\": \"migrated\"
    }" orchestrator-state.yaml
    
    # Move others to deprecated
    yq -i ".deprecated_phase_integrations = (.phase_integration_branches[] | 
           select(.branch != \"$LATEST_PHASE_BRANCH\"))" orchestrator-state.yaml
    
    # Remove old structure
    yq -i 'del(.phase_integration_branches)' orchestrator-state.yaml
fi
```

## Violations and Penalties

### 🔴🔴🔴 AUTOMATIC FAILURE (-100%)
- Using ANY wave integration branch other than `current_wave_integration.branch`
- Using ANY phase integration branch other than `current_phase_integration.branch`
- Architect reviewing a deprecated wave integration branch
- Architect assessing a deprecated phase integration branch
- Creating new integration without deprecating old ones
- Having multiple "active" integrations for same wave or phase

### 🚨🚨🚨 BLOCKING VIOLATIONS (-50%)
- Missing `current_wave_integration` field when wave integration exists
- Missing `current_phase_integration` field when phase integration exists
- Missing `deprecated_wave_integrations` tracking
- Missing `deprecated_phase_integrations` tracking
- Not validating current before use
- Mixing wave and phase integration branches

## Example Scenarios

### Scenario 1: Initial Wave Integration
```yaml
# Before: No wave integration exists
current_wave_integration: null

# Create first wave integration
current_wave_integration:
  phase: 2
  wave: 1
  branch: "phase2-wave1-integration-20250901-100000"
  status: "active"
  created_at: "2025-09-01T10:00:00Z"
  type: "initial"
```

### Scenario 2: Wave Integration Retry (After Failure)
```yaml
# Before: Failed wave integration
current_wave_integration:
  phase: 2
  wave: 1
  branch: "phase2-wave1-integration-20250901-100000"
  status: "active"
  
# After: New retry attempt
current_wave_integration:
  phase: 2
  wave: 1
  branch: "phase2-wave1-integration-20250901-143000"
  status: "active"
  created_at: "2025-09-01T14:30:00Z"
  type: "post_fixes"
  
deprecated_wave_integrations:
  - phase: 2
    wave: 1
    branch: "phase2-wave1-integration-20250901-100000"
    status: "deprecated"
    deprecated_at: "2025-09-01T14:30:00Z"
    reason: "test failures, retry needed"
```

### Scenario 3: Phase Integration After All Waves
```yaml
# State after all waves complete
current_wave_integration:
  phase: 2
  wave: 3  # Last wave
  branch: "phase2-wave3-integration-20250901-160000"
  status: "active"
  
# Create phase integration
current_phase_integration:
  phase: 2
  branch: "phase2-integration-20250901-180000"
  status: "active"
  created_at: "2025-09-01T18:00:00Z"
  type: "initial"
```

### Scenario 4: Multiple Integration Levels Active
```yaml
# Valid state: Different integration levels can coexist
current_wave_integration:
  phase: 3
  wave: 1
  branch: "phase3-wave1-integration-20250902-090000"
  status: "active"
  
current_phase_integration:
  phase: 2  # Previous phase
  branch: "phase2-integration-20250901-180000"
  status: "active"
  
# This is VALID because they're different phases
```

## Integration with Other Rules

### Works With:
- **R105** - Wave Completion Protocol (creates wave integrations)
- **R258** - Mandatory Wave Review Report (uses current wave integration)
- **R259** - Mandatory Phase Integration After Fixes (creates new current phase integration)
- **R285** - Mandatory Phase Integration Before Assessment (creates initial current phase integration)
- **R257** - Phase Assessment Report (references current phase integration)
- **R300** - Fix Management Protocol (fixes go to efforts, not integration)
- **R250** - Integration Isolation (separate workspaces for wave/phase integrations)

### Supersedes:
- Old `phase_integration_branches` array structure
- Old `integration_branches` array structure
- Any ambiguous integration branch selection logic
- Manual wave/phase integration branch tracking

## Verification Commands

```bash
# Check current WAVE integration
yq '.current_wave_integration' orchestrator-state.yaml

# Check current PHASE integration
yq '.current_phase_integration' orchestrator-state.yaml

# Verify only one active wave integration
yq '.current_wave_integration | select(.status == "active")' orchestrator-state.yaml
# Must have status: "active"

# Verify only one active phase integration
yq '.current_phase_integration | select(.status == "active")' orchestrator-state.yaml
# Must have status: "active"

# List deprecated wave integrations
yq '.deprecated_wave_integrations[]' orchestrator-state.yaml

# List deprecated phase integrations
yq '.deprecated_phase_integrations[]' orchestrator-state.yaml

# Validate before architect wave review
PHASE=2
WAVE=1
WAVE_CURRENT=$(yq ".current_wave_integration | select(.phase == $PHASE and .wave == $WAVE).branch" orchestrator-state.yaml)
echo "Architect MUST review wave integration: $WAVE_CURRENT"

# Validate before architect phase assessment
PHASE=2
PHASE_CURRENT=$(yq ".current_phase_integration | select(.phase == $PHASE).branch" orchestrator-state.yaml)
echo "Architect MUST assess phase integration: $PHASE_CURRENT"
```

## Critical Implementation Notes

### Wave Integration Management:
1. **NEVER** allow multiple active wave integrations for the same wave
2. **ALWAYS** deprecate old wave integration before creating new
3. **VALIDATE** current wave integration before ANY wave review
4. **ARCHITECT** must ONLY review current_wave_integration.branch for waves
5. **WAVE_COMPLETE** state creates/updates current wave integration

### Phase Integration Management:
1. **NEVER** allow multiple active phase integrations for the same phase
2. **ALWAYS** deprecate old phase integration before creating new
3. **VALIDATE** current phase integration before ANY phase assessment
4. **ARCHITECT** must ONLY assess current_phase_integration.branch for phases
5. **PHASE_INTEGRATION** state creates/updates current phase integration

### Error Recovery:
1. **ERROR_RECOVERY** at wave level creates new current wave integration
2. **ERROR_RECOVERY** at phase level creates new current phase integration
3. **ALWAYS** deprecate failed integration attempts with clear reasons
4. **TRACK** all retry attempts in deprecated lists

## Grading Impact

- **Correct Usage**: +10% bonus for perfect integration tracking at both levels
- **Using Deprecated Wave Integration**: -100% AUTOMATIC FAILURE
- **Using Deprecated Phase Integration**: -100% AUTOMATIC FAILURE
- **Missing Wave Integration Validation**: -50% per occurrence
- **Missing Phase Integration Validation**: -50% per occurrence
- **Ambiguous Current** (multiple active for same wave/phase): -75%
- **Missing Tracking Fields**: -25% per missing field

---

**REMEMBER**: 
- There can be ONLY ONE current integration per wave
- There can be ONLY ONE current integration per phase
- Using the wrong one = PROJECT FAILURE!
- Wave and phase integrations are SEPARATE tracking systems