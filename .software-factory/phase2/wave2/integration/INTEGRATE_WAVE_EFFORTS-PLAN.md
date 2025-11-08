# Integration Plan - Wave 2.2 (Iteration 7)

**Date**: 2025-11-02 23:36:35 UTC
**Phase**: 2
**Wave**: 2.2
**Iteration**: 7
**Target Branch**: idpbuilder-oci-push/phase2/wave2/integration
**Base Branch**: idpbuilder-oci-push/phase2/wave1/integration

## Context

This is a **re-integration** (Iteration 7) after fixing BUG-019 (R359 code deletion violation).

### Previous Integration Issues (Iteration 6)
- Integration completed but contained R359 violation
- Bug fixes applied to effort branches
- Fresh integration required from clean Wave 2.1 base

### R300 Verification
- Effort branches contain fix commits
- Effort 2.2.1: 37b5e68 "fix: remove out-of-scope stub files to fix size violation"
- Effort 2.2.1: aa20b98 "fix(root): pass viper instance to NewPushCommand"
- All fixes are in effort branches (R300 compliance verified)

## Branches to Integrate (Sequential per R308)

### 1. Effort 2.2.1 - Registry Override & Viper Integration
- **Branch**: origin/idpbuilder-oci-push/phase2/wave2/effort-1-registry-override-viper
- **Latest Commit**: 37b5e68 "fix: remove out-of-scope stub files to fix size violation"
- **Dependencies**: Wave 2.1 (Push Command Core, Progress Reporter)
- **Features**:
  - Registry override configuration via Viper
  - Environment variable support for registry settings
  - Bug fixes applied post-review

### 2. Effort 2.2.2 - Environment Variable Support
- **Branch**: origin/idpbuilder-oci-push/phase2/wave2/effort-2-env-variable-support
- **Latest Commit**: e08ef95 "todo: orchestrator - INTEGRATE_WAVE_EFFORTS complete"
- **Dependencies**: Effort 2.2.1 (requires Viper integration)
- **Features**:
  - Comprehensive environment variable configuration
  - Integration with Viper-based registry override
  - Full test coverage

## Merge Strategy

### Step 1: Reset Integration Branch
- Hard reset to origin/idpbuilder-oci-push/phase2/wave1/integration
- Clean slate for fresh integration
- Preserve iteration 6 work in git history (not deleted)

### Step 2: Sequential Merging (R308)
1. Merge effort-1 (registry-override-viper) using --no-ff
2. Validate build after effort-1
3. Merge effort-2 (env-variable-support) using --no-ff
4. Validate build after effort-2

### Step 3: Comprehensive Testing (R265)
- Run full test suite
- Verify all tests pass
- Check test coverage meets wave targets
- Validate binary functionality

### Step 4: Version Consistency (R381)
- Verify go.mod consistency across merged branches
- No version updates during integration
- All efforts must have matching library versions

## Conflict Resolution Strategy

If conflicts occur:
- Resolve per R262 guidelines (conflict resolution only)
- Document all conflict resolutions in work log
- NO new code creation (R361)
- NO bug fixes during integration (R266)

## Expected Outcome

- Fully integrated Wave 2.2 with both efforts
- All tests passing
- No R359 violations
- Clean build
- Ready for architect review

## Risk Mitigation

- Original effort branches remain unmodified (R262)
- No cherry-picking (preserves history)
- No pre-commit bypass (R506)
- All changes tracked in work log (R264)

## Success Criteria

- Build completes successfully
- All tests pass (100% of existing tests)
- No new bugs introduced
- Integration report complete with all required sections
- Work log is replayable
