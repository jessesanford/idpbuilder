# 🔴🔴🔴 RULE R512 - TRUNK-BASED DEVELOPMENT INTEGRATE_WAVE_EFFORTS MODEL (SUPREME LAW) 🔴🔴🔴

## Metadata
- **Rule ID**: R512
- **Title**: Trunk-Based Development Integration Model
- **Criticality**: 🔴🔴🔴 SUPREME LAW
- **Category**: Integration Architecture
- **Version**: 1.0
- **Date**: 2025-10-05
- **Status**: PRODUCTION
- **Penalty**: -100% for violations
- **Related Rules**: R308 (Incremental Branching), R410 (Planned)

## The Rule

**THE INTEGRATE_WAVE_EFFORTS MODEL: We integrate EFFORT BRANCHES, not INTEGRATE_WAVE_EFFORTS BRANCHES!**

All wave, phase, and project integrations follow a trunk-based development model where:
1. **EFFORT BRANCHES** are integrated (the source code changes)
2. **INTEGRATE_WAVE_EFFORTS BRANCHES** are TARGET destinations (not sources!)
3. Integration branches build INCREMENTALLY on previous integrations (R308)

## 🔴🔴🔴 CRITICAL UNDERSTANDING 🔴🔴🔴

### WHAT WE INTEGRATE (Source Branches)

**ALL integration levels integrate EFFORT BRANCHES:**

#### Wave Integration (Sequential Rebuild - Wave-Scoped)
```bash
# Wave 2 Integration
BASE: phase1/wave2/effort1 (FIRST effort of wave 2)
INTEGRATE:
  - phase1/wave2/effort2  # EFFORT branch (sequential merge)
  - phase1/wave2/effort3  # EFFORT branch (sequential merge)
RESULT: phase1-wave2-integration

# Wave-scoped sequential rebuild:
# - Base: First effort of THIS wave
# - Merge: Remaining efforts of THIS wave
# - Purpose: Test wave-internal sequential mergeability
# - NOT integrating previous wave integration!
# - Creates REAL merges for testing
```

#### Phase Integration (Sequential Rebuild - Phase-Scoped)
```bash
# Phase 1 Integration
BASE: phase1/wave1/effort1 (FIRST effort of phase)
INTEGRATE:
  - phase1/wave1/effort2  # Sequential merge
  - phase1/wave1/effort3  # Sequential merge
  - phase1/wave2/effort1  # REAL MERGE (not in base via cascade)
  - phase1/wave2/effort2  # REAL MERGE
  - phase1/wave3/effort1  # REAL MERGE
  # ... etc
RESULT: phase1-integration

# Phase-scoped sequential rebuild:
# - Base: First effort of PHASE (not last wave integration!)
# - Merge: ALL efforts from ALL waves in sequence
# - Purpose: Test cross-wave sequential mergeability
# - NOT integrating wave-integration branches!
# - Creates REAL merges testing cross-wave compatibility
```

#### Project Integration (Sequential Rebuild - Project-Scoped)
```bash
# Project Integration
BASE: main (the trunk!)
INTEGRATE:
  - phase1/wave1/effort1  # FIRST effort of entire project
  - phase1/wave1/effort2  # Sequential merge
  - phase1/wave2/effort1  # Sequential merge
  # ... through ALL phases
  - phase2/wave1/effort1  # Sequential merge
  - phase2/wave2/effort1  # Sequential merge
  # ... all efforts from ALL phases
RESULT: project-integration

# Project-scoped sequential rebuild:
# - Base: MAIN (the trunk!)
# - Merge: ALL efforts from ALL phases in sequence
# - Purpose: Test COMPLETE trunk-based sequential mergeability
# - NOT integrating phase-integration branches!
# - Creates REAL merges testing if every effort can merge to main
```

## Why This Model?

### Simulates Trunk-Based Development PRs

**This model simulates how PRs work in real trunk-based development:**

1. **In Real Development**:
   - Developer creates feature branch from main
   - Developer submits PR to main
   - CI runs tests against main + feature
   - If approved, feature merges to main
   - Next developer starts from new main state

2. **In Software Factory**:
   - SW Engineer creates effort branch from wave base
   - Code Reviewer validates effort
   - Integration applies effort to wave-integration (simulates merge to main)
   - Next wave starts from wave-integration (simulates starting from new main)

### Benefits

1. **Testing Realism**: Each effort is tested against the "trunk" (integration branch)
2. **Conflict Detection**: Conflicts appear when integrating efforts, not when merging integrations
3. **Incremental Validation**: Each layer builds on validated previous work
4. **Clear Ownership**: Effort branches own code, integration branches are build artifacts
5. **Fix Propagation**: Fixes go to effort branches (source), then re-integrate

## 🚨🚨🚨 WHAT THE PREVENTION-RECOMMENDATIONS.md GOT WRONG 🚨🚨🚨

### INVALID Recommendations (DO NOT IMPLEMENT)

The PREVENTION-RECOMMENDATIONS.md file contains **FUNDAMENTALLY FLAWED** validation logic:

#### Lines 73-101 - Phase Integration Validation  - **COMPLETELY WRONG**
```bash
# ❌❌❌ WRONG - DO NOT IMPLEMENT THIS!
validate_phase_integration() {
    # Sources MUST be wave integration branches only
    for source in "${SOURCE_BRANCHES[@]}"; do
        if [[ ! "$source" =~ -wave[0-9]+-integration$ ]]; then
            echo "❌ Cannot merge non-wave-integration branch"
            return 1
        fi
    done
}
```

**WHY THIS IS WRONG:**
- Phase integration integrates EFFORT branches, NOT wave-integration branches!
- This validation would reject the correct integration model
- This contradicts trunk-based development principles

#### Lines 105-128 - Project Integration Validation - **COMPLETELY WRONG**
```bash
# ❌❌❌ WRONG - DO NOT IMPLEMENT THIS!
validate_project_integration() {
    # Sources MUST be phase integration branches only
    for source in "${SOURCE_BRANCHES[@]}"; do
        if [[ ! "$source" =~ phase[0-9]+-integration$ ]]; then
            echo "❌ Cannot merge '$source' into project integration!"
            return 1
        fi
    done
}
```

**WHY THIS IS WRONG:**
- Project integration integrates ALL EFFORT branches across all phases!
- Phase-integration branches are TARGETS, not sources
- This validation would prevent the entire system from working

### VALID Concepts (Need Correction)

These concepts are good but need the logic fixed:

1. **Pre-merge validation** - YES, but validate EFFORT branches, not integration branches
2. **Integration level tracking** - YES, track which integration layer we're operating in
3. **Base branch validation** - YES, verify correct R308 incremental base
4. **Workspace validation** - YES, ensure orchestrator knows its location

## 🔴🔴🔴 CORRECT VALIDATION LOGIC 🔴🔴🔴

### Wave Integration Validation
```bash
validate_wave_integration() {
    local TARGET_BRANCH=$1
    local SOURCE_BRANCHES=("${@:2}")

    # Target must be wave integration branch
    if [[ ! "$TARGET_BRANCH" =~ -wave[0-9]+-integration$ ]]; then
        echo "❌ R512 VIOLATION: Target '$TARGET_BRANCH' is not a wave integration branch!"
        return 1
    fi

    # Sources MUST be effort branches from this wave
    local PHASE=$(echo "$TARGET_BRANCH" | grep -oP 'phase\K[0-9]+')
    local WAVE=$(echo "$TARGET_BRANCH" | grep -oP 'wave\K[0-9]+')

    for source in "${SOURCE_BRANCHES[@]}"; do
        # Must NOT be an integration branch
        if [[ "$source" =~ -integration$ ]]; then
            echo "❌ R512 VIOLATION: Cannot merge integration branch '$source' into wave integration!"
            echo "Wave integration merges EFFORT branches only"
            return 1
        fi

        # Should be from correct phase/wave
        if [[ ! "$source" =~ phase${PHASE}/wave${WAVE}/ ]]; then
            echo "⚠️ WARNING: Effort '$source' not from phase$PHASE/wave$WAVE"
        fi
    done

    echo "✅ R512: Wave integration validation passed"
    return 0
}
```

### Phase Integration Validation
```bash
validate_phase_integration() {
    local TARGET_BRANCH=$1
    local SOURCE_BRANCHES=("${@:2}")

    # Target must be phase integration branch (not wave)
    if [[ "$TARGET_BRANCH" =~ -wave[0-9]+-integration$ ]]; then
        echo "❌ R512 VIOLATION: Using wave integration branch for phase integration!"
        return 1
    fi

    if [[ ! "$TARGET_BRANCH" =~ ^phase[0-9]+-integration$ ]]; then
        echo "❌ R512 VIOLATION: Target '$TARGET_BRANCH' is not a phase integration branch!"
        return 1
    fi

    # Sources MUST be effort branches from ALL waves in this phase
    local PHASE=$(echo "$TARGET_BRANCH" | grep -oP 'phase\K[0-9]+')

    for source in "${SOURCE_BRANCHES[@]}"; do
        # Must NOT be an integration branch
        if [[ "$source" =~ -integration$ ]]; then
            echo "❌ R512 VIOLATION: Cannot merge integration branch '$source' into phase integration!"
            echo "Phase integration merges ALL EFFORT branches from all waves"
            return 1
        fi

        # Should be from correct phase
        if [[ ! "$source" =~ phase${PHASE}/ ]]; then
            echo "⚠️ WARNING: Effort '$source' not from phase$PHASE"
        fi
    done

    echo "✅ R512: Phase integration validation passed"
    return 0
}
```

### Project Integration Validation
```bash
validate_project_integration() {
    local TARGET_BRANCH=$1
    local SOURCE_BRANCHES=("${@:2}")

    # Target must be project integration
    if [[ "$TARGET_BRANCH" != "project-integration" ]]; then
        echo "❌ R512 VIOLATION: Target '$TARGET_BRANCH' is not project integration!"
        return 1
    fi

    # Sources MUST be effort branches from ALL phases and waves
    for source in "${SOURCE_BRANCHES[@]}"; do
        # Must NOT be an integration branch
        if [[ "$source" =~ -integration$ ]]; then
            echo "❌ R512 VIOLATION: Cannot merge integration branch '$source' into project integration!"
            echo "Project integration merges ALL EFFORT branches from all phases/waves"
            return 1
        fi

        # Should match phase/wave/effort pattern
        if [[ ! "$source" =~ phase[0-9]+/wave[0-9]+/ ]]; then
            echo "⚠️ WARNING: Effort '$source' doesn't match expected pattern"
        fi
    done

    echo "✅ R512: Project integration validation passed"
    return 0
}
```

## Integration Agent Requirements

### Wave Integration (Sequential Rebuild)
```bash
# Integration agent receives:
INTEGRATE_WAVE_EFFORTS_TYPE="wave"
TARGET_BRANCH="phase1-wave2-integration"
BASE_BRANCH="phase1/wave2/effort1"  # FIRST effort of THIS wave!
SOURCE_BRANCHES=(
    "phase1/wave2/effort2"  # Remaining efforts only
    "phase1/wave2/effort3"
)

# Agent validates:
# 1. All sources are effort branches (NOT integrations)
# 2. Target is wave-integration branch
# 3. Base is FIRST effort of THIS wave (wave-scoped)
# 4. Sources are remaining efforts of THIS wave only
```

### Phase Integration (Sequential Rebuild)
```bash
# Integration agent receives:
INTEGRATE_WAVE_EFFORTS_TYPE="phase"
TARGET_BRANCH="phase1-integration"
BASE_BRANCH="phase1/wave1/effort1"  # FIRST effort of PHASE!
SOURCE_BRANCHES=(
    # ALL efforts from ALL waves in the phase (in sequence)
    "phase1/wave1/effort2"
    "phase1/wave1/effort3"
    "phase1/wave2/effort1"
    "phase1/wave2/effort2"
    "phase1/wave3/effort1"
)

# Agent validates:
# 1. All sources are effort branches (NOT integrations)
# 2. Target is phase-integration branch
# 3. Base is FIRST effort of PHASE (phase-scoped)
# 4. Sources are ALL remaining efforts from ALL waves in sequence
```

### Project Integration
```bash
# Integration agent receives:
INTEGRATE_WAVE_EFFORTS_TYPE="project"
TARGET_BRANCH="project-integration"
BASE_BRANCH="main"  # The trunk!
SOURCE_BRANCHES=(
    # ALL efforts from ALL phases and waves (including first!)
    "phase1/wave1/effort1"  # First effort of entire project
    "phase1/wave1/effort2"
    # ... through all phases
    "phase2/wave1/effort1"
    "phase2/wave2/effort1"
    # ... all efforts
)

# Agent validates:
# 1. All sources are effort branches (NOT integrations)
# 2. Target is project-integration branch
# 3. Base is MAIN (the trunk!)
# 4. Sources span all phases and waves in sequential order
```

## Common Misunderstandings

### ❌ WRONG: "Phase integration merges wave-integration branches"
**NO!** Phase integration applies ALL effort branches from all waves to the phase-integration target.

The wave-integration branches are CHECKPOINTS showing what the code looked like after each wave. They are NOT merged together.

### ❌ WRONG: "Project integration merges phase-integration branches"
**NO!** Project integration applies ALL effort branches from all phases to the project-integration target.

The phase-integration branches are CHECKPOINTS showing what the code looked like after each phase. They are NOT merged together.

### ✅ CORRECT: "Integration branches are build artifacts"
YES! Think of integration branches like compiled binaries - they are the RESULT of integrating effort branches (source code), not inputs to further integrations.

### ✅ CORRECT: "Each integration level re-applies all relevant efforts"
YES! This ensures:
- Testing against actual integrated state
- Conflict detection at appropriate level
- Clear lineage of changes
- Fix propagation works correctly

## Grading Impact

- **Implementing PREVENTION-RECOMMENDATIONS.md validation as-is**: -100% (breaks entire system)
- **Merging integration branches together**: -100% (violates trunk-based model)
- **Not understanding effort-vs-integration distinction**: -75%
- **Correct implementation of R512 validation**: +100% (system works correctly)

## Integration with Other Rules

- **R308**: Integration branches use incremental bases (wave2 from wave1, phase2-wave1 from phase1)
- **R321**: Fixes go to effort branches, then re-integrate (not to integration branches)
- **R410** (Planned): Will validate integration level and source branch types
- **R327**: Cascade re-integration recreates integration branches from fixed effort branches

## Acknowledgment Required

Before performing ANY integration work, agents must acknowledge:

```
I acknowledge R512 - Trunk-Based Development Integration Model:
- Wave integrations merge EFFORT branches (not wave-integrations)
- Phase integrations merge ALL EFFORT branches from all waves (not wave-integrations)
- Project integrations merge ALL EFFORT branches from all phases/waves (not phase-integrations)
- Integration branches are TARGETS and CHECKPOINTS, not sources for further integration
- The PREVENTION-RECOMMENDATIONS.md file contains INCORRECT validation logic that must NOT be implemented
```

## Summary

**REMEMBER THE MANTRA:**
```
Efforts are source code
Integrations are build artifacts
You don't merge binaries together
You rebuild from source!
```

Every integration level rebuilds from ALL relevant effort branches (source code), not from previous integration branches (build artifacts).

**THIS IS FUNDAMENTAL TO THE ENTIRE SYSTEM!**
