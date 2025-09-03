# 🔴🔴🔴 RULE R220: ATOMIC PR DESIGN REQUIREMENT (SUPREME LAW) 🔴🔴🔴

**Version**: 220.0.0
**Status**: ACTIVE  
**Criticality**: 🔴🔴🔴 SUPREME LAW - ABSOLUTE REQUIREMENT
**Supersedes**: None
**Related Rules**: R307 (PARAMOUNT independent mergeability), R221 (CD requirement), R277 (continuous build), R280 (main branch protection)

## ABSOLUTE REQUIREMENT

**EVERY SINGLE EFFORT AT EVERY LEVEL MUST BE DESIGNED FOR ATOMIC PR MERGEABILITY**

This is a SUPREME LAW for trunk-based development and continuous delivery. No exceptions.

## Core Requirements

### 1. Atomic PR Mergeability (See R307 for PARAMOUNT details)
- **MANDATORY**: Each effort = ONE atomic PR to upstream main
- **MANDATORY**: Must merge independently without breaking build
- **MANDATORY**: Must not introduce bugs or regressions
- **MANDATORY**: Must maintain working binary at all times
- **MANDATORY**: Must be mergeable MONTHS/YEARS later (R307)
- **MANDATORY**: Feature flags for ALL incomplete features
- **FORBIDDEN**: Multi-effort PRs that cannot be split
- **FORBIDDEN**: PRs that break the build when merged alone
- **FORBIDDEN**: Dependencies on unmerged work

### 2. Required Implementation Techniques

#### Feature Flags (MANDATORY for incomplete features)
```yaml
feature_flags:
  - name: "FEATURE_PAYMENT_ENABLED"
    default: false
    purpose: "Hide incomplete payment processing"
  - name: "FEATURE_NEW_UI_ENABLED"
    default: false
    purpose: "Gradual UI rollout"
```

#### Interfaces and Abstractions (MANDATORY for gradual integration)
```yaml
interfaces:
  - name: "IPaymentProcessor"
    purpose: "Define contract for payment implementations"
  - name: "IDataStore"
    purpose: "Allow swapping data store implementations"
```

#### Stubs and Placeholders (MANDATORY for dependencies)
```yaml
stubs:
  - name: "MockPaymentService"
    replaces: "RealPaymentService"
    until: "Payment service effort complete"
```

#### Backward Compatibility (MANDATORY always)
- Never break existing functionality
- Maintain old APIs until migration complete
- Support gradual feature activation

## Planning Requirements

### Architect Responsibilities
When designing architecture at ANY level (phase/wave/effort):
1. **MUST** split work into independently mergeable efforts
2. **MUST** design interfaces that allow gradual implementation
3. **MUST** plan feature flags for incomplete functionality
4. **MUST** ensure each effort can PR to main alone
5. **MUST** document merge order and dependencies
6. **MUST** specify how build remains working between PRs

### Code Reviewer Responsibilities
When creating implementation plans:
1. **MUST** ensure each effort = one atomic PR
2. **MUST** include feature flag implementation details
3. **MUST** specify stub implementations needed
4. **MUST** define interface contracts early
5. **MUST** plan for backward compatibility
6. **MUST** verify tests pass behind flags

### Software Engineer Responsibilities
When implementing:
1. **MUST** implement feature flags as specified
2. **MUST** create stubs for missing dependencies
3. **MUST** maintain backward compatibility
4. **MUST** ensure PR merges cleanly to main
5. **MUST** verify build remains working

## Verification Requirements

### In Planning Templates
```yaml
atomic_pr_verification:
  can_merge_independently: true  # MUST be true
  feature_flags_defined: true     # MUST be true if incomplete
  interfaces_defined: true        # MUST be true if gradual
  stubs_implemented: true         # MUST be true if dependencies
  build_remains_working: true     # MUST ALWAYS be true
  backward_compatible: true       # MUST ALWAYS be true
```

### In Code Reviews
```yaml
atomic_pr_checklist:
  - "✅ PR can merge to main independently"
  - "✅ Build passes with just this PR"
  - "✅ No regressions introduced"
  - "✅ Feature flags hide incomplete work"
  - "✅ Interfaces allow future extension"
  - "✅ Tests pass in isolation"
```

## Example Implementation

### Good Example: E-commerce Platform
```markdown
Phase 1, Wave 1 - Each effort is ONE atomic PR:

Effort 1: User Authentication (PR #1 → main)
- Complete, working feature
- Can merge alone, build passes
- No dependencies on other efforts

Effort 2: Product Catalog (PR #2 → main)
- Behind CATALOG_ENABLED flag
- Stub data for initial implementation
- Can merge without authentication
- Build remains working

Effort 3: Shopping Cart (PR #3 → main)
- Behind CART_ENABLED flag
- Uses IProductCatalog interface
- Mock checkout for now
- Build still works

Each PR merges directly to main!
Build never breaks!
Features activate when complete!
```

### Bad Example: Monolithic PRs
```markdown
❌ WRONG: Single PR with all of Phase 1
❌ WRONG: PR that needs other PRs to work
❌ WRONG: PR that breaks build until next PR
❌ WRONG: No feature flags for incomplete work
```

## Enforcement

### At Planning Time
- Architect MUST design for atomic PRs
- Code Reviewer MUST verify atomicity in plans
- Plans missing atomic design = REJECTED

### At Implementation Time
- Each effort produces exactly ONE PR
- PR must pass ALL tests independently
- PR must merge cleanly to main

### At Review Time
- Verify PR can stand alone
- Verify feature flags present
- Verify build remains working

## Grading Impact

**Violation = -100% IMMEDIATE FAILURE**

Specific violations:
- PR cannot merge independently: -50%
- Build breaks when PR merged: -100%
- No feature flags for incomplete: -30%
- Multiple efforts in one PR: -40%
- Dependencies not stubbed: -30%

## State Machine Integration

This rule applies in ALL planning and implementation states:
- PHASE_ARCHITECTURE_PLANNING
- WAVE_ARCHITECTURE_PLANNING
- PHASE_IMPLEMENTATION_PLANNING
- WAVE_IMPLEMENTATION_PLANNING
- EFFORT_PLAN_CREATION
- IMPLEMENTATION
- SPLIT_IMPLEMENTATION

## Relationship to Continuous Delivery

This rule enables TRUE continuous delivery by ensuring:
1. **Trunk-based development**: All work goes to main
2. **Small, frequent merges**: Each effort = one PR
3. **Always deployable**: Build never breaks
4. **Gradual rollout**: Feature flags control activation
5. **Fast feedback**: Problems detected immediately

## REMEMBER

**EVERY EFFORT = ONE ATOMIC PR TO MAIN**
**BUILD MUST NEVER BREAK**
**FEATURE FLAGS HIDE INCOMPLETE WORK**
**THIS IS NON-NEGOTIABLE**

---

*Failure to follow this rule results in immediate project failure and -100% grade.*