# ERROR RECOVERY PLAN - Phase 1 pkg/certs Duplicate Declarations

## 🔴 CRITICAL: THIS PLAN MUST BE FOLLOWED EXACTLY

**Date Created:** 2025-09-14T18:05:00Z
**Issue:** Build failure due to duplicate type declarations in pkg/certs
**Root Cause:** Phase 1 latent bug exposed by Phase 2 integration

## The Problem

### Duplicate Declarations Found:
1. **TLSConfig struct** - Declared in BOTH:
   - `pkg/certs/types.go:55`
   - `pkg/certs/utilities.go:130` (DUPLICATE - MUST REMOVE)

2. **DefaultTLSConfig function** - Declared in BOTH:
   - `pkg/certs/types.go:150` (approx)
   - `pkg/certs/utilities.go:139` (DUPLICATE - MUST REMOVE)

### Impact:
- Phase 2 Wave 1 integration cannot build
- Registry package cannot import certs package
- All Phase 2 functionality blocked

## Fix Sequence (R327 CASCADE MANDATORY)

### Step 1: Identify Phase 1 Source Branch
- Find which Phase 1 effort introduced pkg/certs
- This is likely in Phase 1 Wave 2 (certificate infrastructure)
- Check efforts: cert-manager, secure-build, or similar

### Step 2: Fix in Source Branch
1. Checkout Phase 1 source branch containing pkg/certs
2. Edit `pkg/certs/utilities.go`:
   - REMOVE lines 130-138 (TLSConfig struct)
   - REMOVE lines 139-147 (DefaultTLSConfig function)
3. Verify build: `go build ./pkg/certs`
4. Run tests: `go test ./pkg/certs/...`
5. Commit: "fix: remove duplicate declarations in pkg/certs/utilities.go"
6. Push to source branch

### Step 3: R327 CASCADE - Re-integrate EVERYTHING

**⚠️ CRITICAL: Each integration must be DELETED and RECREATED from scratch!**

#### Phase 1 Re-integrations:
1. **Phase 1 Wave 1 Re-integration**
   - Delete old integration branch
   - Re-integrate all Wave 1 efforts
   - Verify build passes

2. **Phase 1 Wave 2 Re-integration**
   - Delete old integration branch
   - Re-integrate all Wave 2 efforts (including fixed branch)
   - Verify build passes

3. **Phase 1 Phase-Level Integration**
   - Delete `idpbuilder-oci-build-push/phase1-integration`
   - Re-integrate Wave 1 + Wave 2
   - Verify build passes
   - This becomes new Phase 1 baseline

#### Phase 2 Rebases and Re-integrations:
4. **Rebase ALL Phase 2 Efforts**
   - image-builder → rebase onto new Phase 1 integration
   - gitea-client-split-001 → rebase onto new Phase 1 integration
   - gitea-client-split-002 → rebase onto split-001

5. **Phase 2 Wave 1 Re-integration**
   - Delete old integration branch
   - Re-integrate all rebased Phase 2 efforts
   - Verify build passes
   - Run all tests
   - Verify demos work

## State Machine Progression

```
Current: ERROR_RECOVERY
    ↓
SPAWN_CODE_REVIEWER_FIX_PLAN (create fix plan for Phase 1)
    ↓
WAITING_FOR_FIX_PLANS
    ↓
SPAWN_ENGINEERS_FOR_FIXES (fix Phase 1 source)
    ↓
MONITOR_FIXES
    ↓
[Phase 1 Wave 1 Re-integration]
    ↓
[Phase 1 Wave 2 Re-integration]
    ↓
[Phase 1 Phase Integration]
    ↓
[Phase 2 Rebases]
    ↓
[Phase 2 Wave 1 Re-integration]
    ↓
WAVE_COMPLETE (finally!)
```

## Success Criteria

✅ All duplicate declarations removed
✅ Phase 1 builds successfully
✅ Phase 1 integration branch updated
✅ Phase 2 efforts rebased onto fixed Phase 1
✅ Phase 2 Wave 1 integration builds successfully
✅ All tests pass
✅ All demos functional

## Tracking

The orchestrator-state.json has been updated with:
- `error_recovery_context` - Full details of the issue
- `r327_cascade_required` - Complete sequence of re-integrations
- `next_state_after_recovery` - Where to continue

## DO NOT SKIP ANY STEPS!

Every single re-integration is mandatory per R327. Skipping any step will result in:
- Inconsistent branches
- Build failures reappearing
- Failed R327 compliance
- Grading failure

---
**REMEMBER:** When we restart, check this plan and execute it completely!