# 🔴🔴🔴 R327 CASCADE RE-INTEGRATION PLAN 🔴🔴🔴

## CRITICAL: SUPREME LAW ENFORCEMENT

**Date**: 2025-09-14
**Violation Detected**: Stale wave integrations being used after effort fixes
**Required Action**: CASCADE DELETE AND RECREATE per R327

## Timeline of Violation

1. **Wave 1 Integration Created**: 2025-09-12 03:24:01
2. **Wave 2 Integration Created**: 2025-09-13 14:54:02
3. **Fixes Applied to Efforts**: 2025-09-13 to 2025-09-14
4. **Violation**: Attempted Phase integration without recreating waves
5. **Correction**: NOW - Full CASCADE re-integration

## Phase 1 Wave 1 Re-Integration Requirements

### Effort Branches to Integrate (WITH FIXES):
1. **kind-cert-extraction** ✅ CONTAINS FIXES
   - Branch: `idpbuilder-oci-build-push/phase1/wave1/kind-cert-extraction`
   - Fixed: Docker API imports, format strings
   - Lines: 841 (over limit, was split)

2. **registry-tls-trust**
   - Branch: `idpbuilder-oci-build-push/phase1/wave1/registry-tls-trust`
   - Lines: 700

3. **registry-auth-types-split-001**
   - Branch: `idpbuilder-oci-build-push/phase1/wave1/registry-auth-types-split-001`
   - Lines: 800

4. **registry-auth-types-split-002**
   - Branch: `idpbuilder-oci-build-push/phase1/wave1/registry-auth-types-split-002`
   - Lines: 800

### Wave 1 Integration Configuration:
- **New Branch Name**: `idpbuilder-oci-build-push/phase1/wave1-integration`
- **Base Branch**: `main`
- **Working Directory**: `worktrees/phase1-wave1-integration`
- **Merge Order**: Sequential as listed above
- **Critical**: Must include Docker API fixes from kind-cert-extraction

## Phase 1 Wave 2 Re-Integration Requirements

### Effort Branches to Integrate (WITH FIXES):
1. **cert-validation-split-001** ✅ CONTAINS FIXES
   - Branch: `idpbuilder-oci-build-push/phase1/wave2/cert-validation-split-001`
   - Fixed: TLSConfig duplicates
   - Lines: 779

2. **cert-validation-split-002** ✅ CONTAINS FIXES
   - Branch: `idpbuilder-oci-build-push/phase1/wave2/cert-validation-split-002`
   - Fixed: TLSConfig duplicates
   - Lines: 684

3. **cert-validation-split-003** ✅ CONTAINS FIXES
   - Branch: `idpbuilder-oci-build-push/phase1/wave2/cert-validation-split-003`
   - Fixed: TLSConfig duplicates
   - Lines: 790

4. **fallback-strategies**
   - Branch: `idpbuilder-oci-build-push/phase1/wave2/fallback-strategies`
   - Lines: 697

### Wave 2 Integration Configuration:
- **New Branch Name**: `idpbuilder-oci-build-push/phase1/wave2-integration`
- **Base Branch**: `idpbuilder-oci-build-push/phase1/wave1-integration` (R308 incremental)
- **Working Directory**: `worktrees/phase1-wave2-integration`
- **Merge Order**: Sequential as listed above
- **Critical**: Must include TLSConfig fixes from cert-validation splits

## R327 CASCADE Sequence

### Step 1: Wave 1 Re-Integration
1. Create fresh worktree from main
2. Create new wave1-integration branch
3. Merge all Wave 1 efforts IN ORDER
4. Validate build and tests pass
5. Push to remote

### Step 2: Wave 2 Re-Integration
1. Create fresh worktree from wave1-integration (R308)
2. Create new wave2-integration branch
3. Merge all Wave 2 efforts IN ORDER
4. Validate build and tests pass
5. Push to remote

### Step 3: Phase Integration (AFTER WAVES)
1. Delete old phase integration (already done)
2. Create fresh from wave2-integration (R308)
3. No additional merges needed (waves contain everything)
4. Validate and push

## Validation Requirements

### Per-Wave Validation:
- ✅ All effort branches contain fixes
- ✅ Integration builds successfully
- ✅ Tests pass
- ✅ No duplicate definitions
- ✅ Docker imports resolved
- ✅ TLSConfig conflicts resolved

### R327 Timestamp Validation:
```bash
# Must verify: Integration timestamp > All fix timestamps
INTEGRATION_TIME=$(git log -1 --format=%ct wave-integration)
for effort in efforts/*; do
    FIX_TIME=$(git log -1 --grep="fix:" --format=%ct $effort)
    [ "$INTEGRATION_TIME" -gt "$FIX_TIME" ] || FAIL
done
```

## Integration Agent Instructions

### For Wave 1 Integration Agent:
```
You are integrating Phase 1 Wave 1 with CASCADE fixes per R327.
Base: main
Efforts: kind-cert-extraction, registry-tls-trust, registry-auth-types-split-001, registry-auth-types-split-002
Critical: kind-cert-extraction contains Docker API fixes that MUST be included
Validate: Build and test after integration
```

### For Wave 2 Integration Agent:
```
You are integrating Phase 1 Wave 2 with CASCADE fixes per R327.
Base: phase1/wave1-integration (R308 incremental)
Efforts: cert-validation-split-001/002/003, fallback-strategies
Critical: cert-validation splits contain TLSConfig fixes that MUST be included
Validate: Build and test after integration
```

## Success Criteria

1. ✅ Both wave integrations recreated with fixes
2. ✅ Timestamps show integrations NEWER than fixes
3. ✅ Build succeeds with no errors
4. ✅ Tests pass
5. ✅ Phase integration can proceed cleanly

## Failure Consequences

Per R327:
- Using stale integrations = -100% AUTOMATIC FAILURE
- Not cascading = -100% AUTOMATIC FAILURE
- Partial recreation = -50% penalty

---

**R327 CASCADE IS SUPREME LAW - NO EXCEPTIONS**