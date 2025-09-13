# Phase 2 Wave 1 Rebase Plan

## Situation Analysis
Phase 2 Wave 1 efforts were created and implemented before Phase 1 integration was finalized. Now they need to be rebased onto the fresh Phase 1 integration branch to ensure they have all the Phase 1 foundation code.

## Current State
- **Phase 1 Integration**: `idpbuilder-oci-build-push/phase1/integration` (COMPLETE)
- **Phase 2 Wave 1 Efforts**:
  - gitea-client (main effort)
  - gitea-client-split-001
  - gitea-client-split-002

## Rebase Requirements

### 1. gitea-client
- **Current Branch**: `idpbuilder-oci-build-push/phase2/wave1/gitea-client`
- **Working Directory**: `efforts/phase2/wave1/gitea-client`
- **Current Base**: Unknown (likely main)
- **Target Base**: `idpbuilder-oci-build-push/phase1/integration`
- **Strategy**: Interactive rebase with conflict resolution

### 2. gitea-client-split-001
- **Current Branch**: `idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-001`
- **Working Directory**: `efforts/phase2/wave1/gitea-client-split-001`
- **Dependency**: Must be rebased after gitea-client
- **Target Base**: Rebased gitea-client

### 3. gitea-client-split-002
- **Current Branch**: `idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-002`
- **Working Directory**: `efforts/phase2/wave1/gitea-client-split-002`
- **Dependency**: Must be rebased after split-001
- **Target Base**: Rebased split-001

## Execution Strategy

### Phase 1: Prepare
1. Fetch latest Phase 1 integration branch
2. Create backup branches for all efforts
3. Document current commit SHAs

### Phase 2: Rebase Main Effort
1. Checkout gitea-client branch
2. Rebase onto `idpbuilder-oci-build-push/phase1/integration`
3. Resolve any conflicts (expected in import paths, API calls)
4. Test compilation
5. Push with force-with-lease

### Phase 3: Rebase Splits Sequentially
1. Rebase split-001 onto rebased gitea-client
2. Resolve conflicts
3. Rebase split-002 onto rebased split-001
4. Resolve conflicts
5. Push both splits

### Phase 4: Verification
1. Verify all branches compile
2. Run demos (R291 compliance)
3. Check line counts remain within limits
4. Update state file

## Expected Conflicts
- Import paths (Phase 1 packages may have moved)
- API signatures (Phase 1 APIs may have evolved)
- Certificate handling (Phase 1 implemented cert management)
- Registry types (Phase 1 defined registry interfaces)

## Success Criteria
- [ ] All three branches successfully rebased
- [ ] Code compiles on all branches
- [ ] Demos run successfully
- [ ] Line counts remain within limits
- [ ] No functionality lost in rebase
- [ ] State file updated with completion

## Next Steps After Rebase
1. Continue with Phase 2 Wave 1 review
2. Integrate Wave 1 efforts
3. Proceed to Phase 2 Wave 2