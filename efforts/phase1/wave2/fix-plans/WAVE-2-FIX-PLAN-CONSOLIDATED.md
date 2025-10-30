# Wave 2 Fix Plan - Phase 1 Wave 2

## Executive Summary

**Created**: 2025-10-30T00:07:30Z
**Orchestrator State**: CREATE_WAVE_FIX_PLAN
**Phase**: 1, **Wave**: 2

### Overview
Code Reviewers have identified 2 efforts requiring fixes before Wave 2 integration can proceed:
- **Effort 1.2.2** (registry-client): CRITICAL R320 violation - stub implementations
- **Effort 1.2.3** (auth): MINOR R383 violation - metadata file placement

Both fixes are low-risk and can be executed in parallel.

### Fix Plan Summary

| Effort ID | Effort Name | Severity | Est. Time | Risk | Blocking |
|-----------|-------------|----------|-----------|------|----------|
| 1.2.2 | registry-client | CRITICAL | 45m | LOW | YES |
| 1.2.3 | auth | MINOR | 15m | VERY_LOW | YES |

**Total estimated time**: 60 minutes (can execute in parallel)

---

## Effort 1.2.2: registry-client (CRITICAL)

### Issue
**R320 Violation**: Stub implementations with `panic()` in production code

### Affected Files
- `pkg/auth/interface.go:44` - NewBasicAuthProvider()
- `pkg/tls/interface.go:41` - NewConfigProvider()

### Impact
- Runtime panics if functions called
- Violates R320 Supreme Law (no stubs in production)
- BLOCKS integration until fixed

### Fix Instructions
**Detailed plan**: `efforts/phase1/wave2/fix-plans/FIX-PLAN-effort-2-registry-client--20251029-233955.md`

**Summary**:
1. Remove stub implementations (NewBasicAuthProvider, NewConfigProvider)
2. Keep interface definitions only
3. Ensure all 28 registry tests still pass
4. Run R355 production readiness scan

**Expected Outcome**:
- `pkg/auth/interface.go` contains ONLY interface definition
- `pkg/tls/interface.go` contains ONLY interface definition
- No `panic` statements in production code
- All tests pass

---

## Effort 1.2.3: auth (MINOR)

### Issue
**R383 Violation**: Metadata file in wrong location

### Affected Files
- `./IMPLEMENTATION-COMPLETE.marker` (should be in `.software-factory/`)

### Impact
- Merge conflicts during integration
- Clutters working tree
- Violates R383 metadata organization

### Fix Instructions
**Detailed plan**: `efforts/phase1/wave2/fix-plans/FIX-PLAN-effort-3-auth--20251029-233955.md`

**Summary**:
1. Move `IMPLEMENTATION-COMPLETE.marker` to `.software-factory/`
2. Add timestamp to marker filename
3. Verify all tests still pass (94.1% coverage)
4. Confirm R383 compliance

**Expected Outcome**:
- No metadata files in root directory
- Marker file in `.software-factory/` with timestamp
- All tests pass (94.1% coverage maintained)

---

## Orchestrator Execution Plan

### Step 1: Distribute Fix Plans to Efforts
```bash
# Copy fix plan to each effort directory
cp efforts/phase1/wave2/fix-plans/FIX-PLAN-effort-2-registry-client--20251029-233955.md \
   efforts/phase1/wave2/effort-2-registry-client/FIX-INSTRUCTIONS.md

cp efforts/phase1/wave2/fix-plans/FIX-PLAN-effort-3-auth--20251029-233955.md \
   efforts/phase1/wave2/effort-3-auth/FIX-INSTRUCTIONS.md
```

### Step 2: Update bug-tracking.json
- Record fix assignments for each effort
- Link bugs to fix plans
- Track fix status

### Step 3: Transition to FIX_WAVE_UPSTREAM_BUGS
- Spawn State Manager for shutdown consultation
- Propose next state: `FIX_WAVE_UPSTREAM_BUGS`
- State Manager validates and transitions

### Step 4: SW Engineer Spawning (from FIX_WAVE_UPSTREAM_BUGS state)
**Parallel execution allowed**: YES

Spawn 2 SW Engineers in parallel per R151:
```
Agent 1: Fix effort-2-registry-client (45m)
Agent 2: Fix effort-3-auth (15m)
```

Both agents work in isolated effort directories, no dependencies.

---

## Post-Fix Validation

### After Fixes Complete
1. Code Reviewer re-reviews both efforts
2. Verify all tests pass
3. Run R355 scan on effort-2
4. Verify R383 compliance on effort-3

### Integration Readiness Criteria
- ✅ All fix plans executed successfully
- ✅ All tests passing
- ✅ R320 compliance verified (no stubs)
- ✅ R383 compliance verified (metadata placed correctly)
- ✅ Code quality maintained/improved

---

## Risk Assessment

**Overall Risk**: LOW

**Effort 2 (registry-client)**:
- Risk: LOW (simple removal of stub functions)
- Mitigation: Comprehensive test suite (28 tests)
- Rollback: Keep original branch as backup

**Effort 3 (auth)**:
- Risk: VERY_LOW (metadata move only, no code changes)
- Mitigation: No code touched, tests unchanged
- Rollback: Simple file move reversal

---

## Success Metrics

**Pre-Fix**:
- Effort 2: 90% code quality, 0% production readiness (stubs)
- Effort 3: 100% code quality, 95% production readiness (metadata)

**Post-Fix Expected**:
- Effort 2: 95% code quality, 100% production readiness
- Effort 3: 100% code quality, 100% production readiness

---

## Notes

1. Both efforts have excellent code quality foundations
2. Fixes are straightforward and low-risk
3. No architectural changes required
4. No test changes required
5. Parallel execution safe - no dependencies between fixes
6. Effort 1 (docker-client) already approved - can integrate first
7. Effort 4 (tls) may need similar review for stub patterns

---

## State Machine Path

```
CREATE_WAVE_FIX_PLAN (current)
  → FIX_WAVE_UPSTREAM_BUGS (spawn SW Engineers with fix instructions)
  → MONITORING_EFFORT_FIXES (monitor fix progress)
  → SPAWN_CODE_REVIEWERS_EFFORT_REVIEW (re-review fixed code)
  → MONITORING_EFFORT_REVIEWS (verify fixes approved)
  → INTEGRATE_WAVE_EFFORTS (if all approved)
```

---

**Plan created per**: R313 (bug tracking), R321 (upstream fixes)
**Complies with**: R006 (orchestrator delegation), R383 (metadata organization)
**Execution**: Sequential state progression with parallel SW Engineer spawning
