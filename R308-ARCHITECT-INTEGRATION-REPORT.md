# R308 Architect Integration Report

## Executive Summary

**Decision**: ✅ **R308 MUST be added to architect configuration**

The architect agent requires deep understanding of R308 (Incremental Branching Strategy) to properly design, review, and assess the incremental development flow that is fundamental to trunk-based development.

## Analysis Findings

### 1. Critical Need Identified

The architect has responsibilities that REQUIRE understanding incremental branching:

#### During PHASE_ARCHITECTURE_PLANNING:
- **Designs interfaces for gradual implementation** - Must know waves build incrementally
- **Plans feature layering** - Must understand Wave 2 sees Wave 1's work
- **Documents dependencies** - Must account for incremental base branches

#### During WAVE_ARCHITECTURE_PLANNING:
- **Plans merge strategies** - Must know efforts don't all merge to main
- **Designs incremental features** - Must understand cumulative development
- **Documents parallel vs sequential work** - Must consider incremental bases

#### During WAVE_REVIEW:
- **Verifies effort compliance** - Must check correct base branch usage
- **Validates integration quality** - Must confirm incremental building
- **Checks architectural consistency** - Must verify incremental interfaces

#### During PHASE_ASSESSMENT:
- **Validates development chain** - Must confirm all waves built incrementally
- **Assesses integration readiness** - Must verify no stale code bases
- **Reviews architectural integrity** - Must check incremental design worked

### 2. Gap Analysis Without R308

**Before R308 Integration:**
- ❌ Architect might design assuming all efforts branch from main
- ❌ Could miss violations where Wave 2 didn't build on Wave 1
- ❌ Might approve architectures incompatible with incremental flow
- ❌ Could overlook "big bang" integration anti-patterns

**After R308 Integration:**
- ✅ Architect designs explicitly for incremental building
- ✅ Reviews verify correct branching strategy usage
- ✅ Assessments validate complete incremental chain
- ✅ Architecture supports true trunk-based development

### 3. Relationship to Existing Rules

**R307 (Independent Mergeability)** + **R308 (Incremental Branching)**:
- R307 ensures efforts can merge independently
- R308 ensures they merge to the RIGHT base
- Together: Incrementally independent mergeability

**R220 (Atomic PRs)** + **R308 (Incremental Branching)**:
- R220 ensures each effort is one atomic PR
- R308 ensures PRs build on previous work
- Together: Atomic incremental development

## Implementation Completed

### 1. Main Architect Configuration Updated

**Location**: `.claude/agents/architect.md`

Added R308 as CORE TENANT immediately after R307, emphasizing:
- Wave-to-wave incremental building
- Phase-to-phase incremental progression
- Architecture supporting gradual enhancement
- No "big bang" integration patterns

### 2. State-Specific Rules Enhanced

#### PHASE_ARCHITECTURE_PLANNING:
- Added requirement to design for wave-to-wave building
- Emphasized incremental dependency planning
- Required documentation of incremental strategy
- Architecture must assume cumulative wave visibility

#### WAVE_REVIEW:
- Added R308 verification as first assessment rule
- Required checking correct base branch usage
- Verification of incremental commit inclusion
- -100% penalty for wrong base branch violations

#### PHASE_ASSESSMENT:
- Added complete incremental chain validation
- Required verification of wave-by-wave building
- Confirmed phase integration readiness for next phase
- -100% penalty for broken incremental chain

## Verification Steps

### Correct Incremental Architecture Example:
```yaml
phase_1_architecture:
  wave_1:
    base: main
    creates: "Core interfaces, base services"
  wave_2:
    base: phase1-wave1-integration
    extends: "Adds concrete implementations"
  wave_3:
    base: phase1-wave2-integration
    completes: "Advanced features on top"
```

### Architect Verification Checklist:
```bash
# During Planning
✅ Interfaces designed for incremental addition
✅ No wave assumes direct main access
✅ Features layer wave-by-wave

# During Review
✅ Each effort includes previous wave commits
✅ No stale main branch bases
✅ Integration branches used correctly

# During Assessment
✅ Complete incremental chain verified
✅ Phase integration ready for next phase
✅ No "big bang" merges detected
```

## Impact on Software Factory

### Positive Impacts:
1. **Better Architecture** - Designs explicitly support incremental development
2. **Earlier Detection** - Catch branching violations during reviews
3. **Smoother Integration** - Avoid conflicts from stale bases
4. **True CI/CD** - Enforce trunk-based development principles

### Grading Implications:
- **-100% for violating incremental flow** (R308 penalty)
- **-50% for not verifying base branches** (Review incompleteness)
- **-25% for not documenting incremental strategy** (Planning gap)

## Recommendations

### For Orchestrator:
1. Ensure architect receives phase/wave context in spawn commands
2. Verify architect reports include R308 compliance checks
3. Block progression if architect detects R308 violations

### For SW Engineers:
1. Must verify working from correct incremental base
2. Should refuse work if given wrong base branch
3. Must include previous wave work in their efforts

### For Code Reviewers:
1. Should verify effort has correct incremental base
2. Must check for conflicts with previous wave work
3. Should validate incremental development approach

## Conclusion

R308 integration into architect configuration is **COMPLETE and CRITICAL**. The architect now has full awareness of incremental branching requirements and will:

1. Design architectures that explicitly support incremental development
2. Verify correct branching during wave reviews
3. Validate complete incremental chains during phase assessments
4. Ensure true trunk-based development compliance

This change strengthens the Software Factory's adherence to modern CI/CD practices and prevents "big bang" integration failures.

## Change Summary

**Files Modified:**
- `.claude/agents/architect.md` - Added R308 as CORE TENANT
- `agent-states/architect/PHASE_ARCHITECTURE_PLANNING/rules.md` - Added incremental design requirements
- `agent-states/architect/WAVE_REVIEW/rules.md` - Added base branch verification
- `agent-states/architect/PHASE_ASSESSMENT/rules.md` - Added chain validation

**Commit**: `67c533d` - "update: architect awareness of incremental branching (R308)"

---
*Report generated by Software Factory Manager*
*Date: 2025-09-02*