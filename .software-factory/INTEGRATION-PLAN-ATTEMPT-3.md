# Phase 1 Wave 2 Integration Plan - Attempt #3

**Date**: $(date -u +%Y-%m-%dT%H:%M:%SZ)
**Agent**: Integration Agent
**R520 Attempt**: 3 of 5
**Base Branch**: idpbuilder-push-oci/phase1-wave1-integration

## Integration Context

### Previous Attempts
- **Attempt #1**: Initial integration - discovered BUG-007 (PushCmd redeclared)
- **Attempt #2**: Applied BUG-007 fix - discovered BUG-010 (PushOptions struct mismatch)
- **Attempt #3 (Current)**: Check if BUG-010 fixed, apply known fixes

### Known Issues
- **BUG-007**: PushCmd redeclared (KNOWN FIX: delete duplicate push.go)
- **BUG-010**: PushOptions struct mismatch (Status: PENDING - check if fixed)

## Branches to Integrate (R306 Dependency Order)

1. **E1.2.1-command-structure**
   - Remote: effort-E1.2.1/idpbuilder-push-oci/phase1/wave2/command-structure
   - Creates: pkg/cmd/push/ command framework
   - Dependencies: None
   - Known Issue: Contains push.go that causes BUG-007

2. **E1.2.2-registry-authentication-split-001**
   - Remote: effort-E1.2.2-split-001/idpbuilder-push-oci/phase1/wave2/registry-authentication-split-001
   - Creates: Authentication basics
   - Dependencies: E1.2.1
   
3. **E1.2.2-registry-authentication-split-002**
   - Remote: effort-E1.2.2-split-002/idpbuilder-push-oci/phase1/wave2/registry-authentication-split-002
   - Creates: Retry logic
   - Dependencies: E1.2.2-split-001

4. **E1.2.3-image-push-operations-split-001**
   - Remote: effort-E1.2.3-split-001/idpbuilder-push-oci/phase1/wave2/image-push-operations-split-001
   - Creates: Core push operations
   - Dependencies: E1.2.2 (both splits)

5. **E1.2.3-image-push-operations-split-002**
   - Remote: effort-E1.2.3-split-002/idpbuilder-push-oci/phase1/wave2/image-push-operations-split-002
   - Creates: Discovery and pusher
   - Dependencies: E1.2.3-split-001
   - Recent Fixes: BUG-002 fix (commit 546bfb8)

6. **E1.2.3-image-push-operations-split-003**
   - Remote: effort-E1.2.3-split-003/idpbuilder-push-oci/phase1/wave2/image-push-operations-split-003
   - Creates: Operations and tests
   - Dependencies: E1.2.3-split-002
   - Recent Fixes: R291 fix (commit 6628dff)

## Merge Strategy

1. Reset to Wave 1 integration base
2. Create fresh attempt #3 branch
3. Merge all efforts in dependency order
4. After merging E1.2.1, check for BUG-007 (R521 protocol)
5. Build verification after each merge
6. Apply known fixes only (R521 - conflict resolution, not development)
7. Document new bugs (R522 - check for duplicates first)

## Expected Outcome

- All 6 branches merged successfully
- Build passes (BUG-007 and BUG-010 resolved or documented)
- Wave-level demo created (R291)
- Complete documentation (R343)
- R520 metadata updated

## R521 Known Fixes Available

**BUG-007 (PushCmd redeclared)**:
- Source: Wave 1 BUILD-FIX-SUMMARY.md
- Fix: Delete pkg/cmd/push/push.go after merging E1.2.1
- Classification: Conflict resolution (ALLOWED per R521)

**BUG-010 (PushOptions struct mismatch)**:
- Status: PENDING - check if upstream fix applied
- If NOT fixed: Document as new bug (R266 - do not fix)
- If FIXED: Verify build passes

## Success Criteria

- [ ] All merges complete without conflicts
- [ ] Build verification passes
- [ ] BUG-007 resolved via known fix
- [ ] BUG-010 status determined (fixed or documented)
- [ ] R291 demo created and functional
- [ ] R520 metadata updated with results
- [ ] CONTINUE-SOFTWARE-FACTORY=TRUE emitted
