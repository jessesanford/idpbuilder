# Phase 2 Wave 1 Integration Merge Plan

**Generated:** 2025-09-14T18:52:00Z
**Code Reviewer:** code-reviewer
**State:** WAVE_MERGE_PLANNING

## Target Integration Branch
- **Branch Name:** idpbuilder-oci-build-push/phase2/wave1/integration
- **Base:** idpbuilder-oci-build-push/phase1/integration
- **Location:** /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase2/wave1/integration-workspace/repo

## 🎬 Demo Execution Plan (R330/R291 Compliance)

### Demo Requirements Overview
Per R330 and R291, ALL integrations MUST demonstrate working functionality.

### Demo Execution Sequence
1. **After Each Effort Merge:**
   - Run effort-specific demo script if exists
   - Capture output in `demo-results/effort-X-demo.log`
   - Continue even if individual demo fails (document for review)

2. **After All Merges Complete:**
   - Run integrated wave demo: `./wave-demo.sh`
   - Verify all effort features work together
   - Capture evidence in `demo-results/wave-integration-demo.log`

3. **Demo Validation Gates (R291):**
   - ✅ BUILD GATE: Code must compile
   - ✅ TEST GATE: All tests must pass
   - ✅ DEMO GATE: Demo scripts must execute
   - ✅ ARTIFACT GATE: Build outputs must exist

### Demo Files Expected
Based on effort plans (R330), these demos should exist:
- [ ] image-builder/demo-features.sh (from effort plan)
- [ ] gitea-client-split-001/demo-features.sh (from effort plan)
- [ ] gitea-client-split-002/demo-features.sh (from effort plan)
- [ ] WAVE-DEMO.md (integration demo documentation)
- [ ] wave-demo.sh (integrated demo script)

### Demo Failure Protocol
If ANY demo fails during integration:
1. Document failure in INTEGRATION_REPORT.md
2. Mark Demo Status: FAILED
3. This will trigger ERROR_RECOVERY per R291
4. Fixes must be made in effort branches (R292)

## Branches to Merge (IN ORDER)

### 1. idpbuilder-oci-build-push/phase2/wave1/image-builder
- **Type:** Original effort branch (with duplicate fix applied)
- **Base:** idpbuilder-oci-build-push/phase1/integration
- **Size:** ~600 lines (estimated)
- **Dependencies:** None (foundational image builder)
- **Conflicts Expected:** None (new functionality)
- **Special Notes:** Contains emergency fix for duplicate TLSConfig struct
- **Merge Command:**
  ```bash
  git fetch origin idpbuilder-oci-build-push/phase2/wave1/image-builder
  git merge origin/idpbuilder-oci-build-push/phase2/wave1/image-builder --no-ff \
    -m "Integrate E2.1.1-image-builder into Phase 2 Wave 1"
  ```

### 2. idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-001
- **Type:** Split branch (1 of 2)
- **Base:** idpbuilder-oci-build-push/phase1/integration
- **Size:** ~400 lines (first part of gitea client)
- **Dependencies:** May depend on image-builder types
- **Conflicts Expected:** Possible in go.mod/go.sum
- **Special Notes:** First split of gitea-client effort
- **Merge Command:**
  ```bash
  git fetch origin idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-001
  git merge origin/idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-001 --no-ff \
    -m "Integrate E2.1.2-gitea-client-split-001 into Phase 2 Wave 1"
  ```

### 3. idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-002
- **Type:** Split branch (2 of 2)
- **Base:** Follows gitea-client-split-001
- **Size:** ~400 lines (completion of gitea client)
- **Dependencies:** Must merge after split-001
- **Conflicts Expected:** None (sequential splits)
- **Special Notes:** Final split of gitea-client effort
- **Merge Command:**
  ```bash
  git fetch origin idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-002
  git merge origin/idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-002 --no-ff \
    -m "Integrate E2.1.2-gitea-client-split-002 into Phase 2 Wave 1"
  ```

## Excluded Branches (Too Large)
These original branches were split and should NOT be merged:
- idpbuilder-oci-build-push/phase2/wave1/gitea-client (original, exceeded limit - replaced by splits)

## Merge Strategy
1. **Merge Type:** --no-ff (preserve branch history)
2. **Conflict Resolution:** Favor newer implementation when conflicts occur
3. **Testing:** Run unit tests after each merge
4. **Validation:** Check total size after all merges

## Expected Conflicts
Based on branch analysis, conflicts are likely in:
- `go.mod/go.sum` - Multiple efforts adding dependencies
- Minimal other conflicts expected as Phase 2 adds new functionality

## Validation Steps
1. After each merge:
   ```bash
   # Test the code
   go test ./...

   # Run effort demo if exists (R330 compliance)
   EFFORT_NAME="[effort-name-from-merge]"
   if [ -f "${EFFORT_NAME}/demo-features.sh" ]; then
       echo "🎬 Running ${EFFORT_NAME} demo..."
       bash "${EFFORT_NAME}/demo-features.sh" > "demo-results/${EFFORT_NAME}-demo.log" 2>&1
       echo "Demo exit code: $?"
   fi
   ```
2. After all merges:
   ```bash
   # Run integration tests
   make test-integration || go test ./...

   # Check size compliance
   PROJECT_ROOT=$(pwd)
   while [ "$PROJECT_ROOT" != "/" ]; do
       [ -f "$PROJECT_ROOT/orchestrator-state.json" ] && break
       PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
   done
   $PROJECT_ROOT/tools/line-counter.sh

   # Run integrated wave demo (R291 requirement)
   echo "🎬 Running integrated wave demo..."
   if [ -f "./wave-demo.sh" ]; then
       bash ./wave-demo.sh > demo-results/wave-integration-demo.log 2>&1
       DEMO_STATUS=$?
       echo "Wave demo status: $DEMO_STATUS"
   else
       echo "⚠️ WARNING: No wave-demo.sh found - creating basic demo"
       # Integration agent should create basic demo if missing
   fi
   ```
3. Final validation:
   - Verify all effort features are present
   - Confirm no effort was missed
   - Check combined size is reasonable
   - ✅ Verify all demos executed (R291)
   - ✅ Capture demo evidence for review

## Risk Assessment
- **Low Risk:** Sequential splits should merge cleanly
- **Low Risk:** New Phase 2 functionality unlikely to conflict with Phase 1
- **Medium Risk:** Dependency management in go.mod
- **Mitigation:** Test after each merge to catch issues early

## Integration Agent Instructions
1. CD to integration directory before starting
2. Execute merges in the EXACT order specified
3. Run tests after EACH merge
4. Document any conflicts encountered
5. Create work-log.md with all operations
6. Generate INTEGRATION-REPORT.md when complete
7. DO NOT merge the original gitea-client branch (use splits only)

## Validation Checklist

After each merge:
- [ ] No uncommitted changes (`git status`)
- [ ] Build succeeds (`go build ./...`)
- [ ] Tests pass (`go test ./...`)
- [ ] No duplicate declarations
- [ ] Dependencies resolved (`go mod tidy`)

After all merges:
- [ ] All 3 branches integrated (image-builder + 2 gitea-client splits)
- [ ] Total size within expectations (~1,400 lines estimated)
- [ ] Integration tests pass
- [ ] No missing functionality
- [ ] Clean commit history
- [ ] Demo scripts executed and documented

## Notes for Integration Agent

1. **DO NOT** merge the parent 'gitea-client' branch - only splits
2. **ALWAYS** use `--no-ff` to preserve merge history
3. **RESOLVE** conflicts conservatively - when in doubt, accept both
4. **TEST** after each merge to catch issues early
5. **DOCUMENT** any unexpected conflicts in the integration report
6. **VERIFY** the duplicate TLSConfig fix is preserved from image-builder

## Expected Final State

After successful integration:
- Branch: `idpbuilder-oci-build-push/phase2/wave1/integration`
- Contains: All Phase 2 Wave 1 efforts (3 branches total: image-builder + 2 gitea-client splits)
- Size: Approximately 1,400 lines of new code
- Status: Ready for architect review and further integration

---

**Created**: 2025-09-14T18:52:00Z
**Created By**: Code Reviewer Agent (WAVE_MERGE_PLANNING state)
**Integration Target**: idpbuilder-oci-build-push/phase2/wave1/integration