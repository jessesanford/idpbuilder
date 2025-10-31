# INTEGRATE_WAVE_EFFORTS AGENT - EXECUTE P2W1 CASCADE INTEGRATE_WAVE_EFFORTS

You are @agent-integration in INIT state tasked with executing Phase 2 Wave 1 integration as part of CASCADE Operation #5.

## CRITICAL CASCADE CONTEXT
- This is CASCADE Op #5 RERUN: Recreating P2W1 integration after P1W1 rebuild
- The builder effort has been rebased onto the new phase1-wave1-integration
- Infrastructure already created in: `/home/vscode/workspaces/this-is-not-the-target-repo-this-is-for-orchestrator-planning-only/efforts/phase2/wave1/wave1-integration-cascade-20250919-233312`
- Branch: `idpbuilder-oci-build-push/phase2-wave1-integration-cascade-20250919-233415`

## YOUR TASK

1. **First emit timestamp**: echo "🕐 Integration Agent P2W1 CASCADE: $(date '+%Y-%m-%d %H:%M:%S')"

2. **Navigate to workspace**:
   ```bash
   cd /home/vscode/workspaces/this-is-not-the-target-repo-this-is-for-orchestrator-planning-only/efforts/phase2/wave1/wave1-integration-cascade-20250919-233312
   ```

3. **Read the merge plan**: PHASE2-WAVE1-CASCADE-MERGE-PLAN.md

4. **Verify current state**:
   ```bash
   git status
   git branch --show-current
   git log --oneline -5
   ```

5. **Clean reset to base**:
   ```bash
   git fetch origin
   git reset --hard origin/idpbuilder-oci-build-push/phase1-wave1-integration
   git clean -fd
   ```

6. **Execute merges per plan**:
   - Merge the builder effort branch
   - Use the merge commands from the CASCADE merge plan
   - Verify each merge completes cleanly

7. **After merge**:
   - Run tests: `go test ./...`
   - Check build: `go build ./...`
   - If issues, document them

8. **Measure implementation size**:
   ```bash
   /home/vscode/software-factory-template/tools/line-counter.sh \
     --base origin/idpbuilder-oci-build-push/phase1-wave1-integration \
     --branch HEAD
   ```

9. **Push integration branch**:
   ```bash
   git push origin HEAD:idpbuilder-oci-build-push/phase2-wave1-integration-cascade-20250919-233415 --force-with-lease
   ```

10. **Report completion** with:
    - Integration status (PROJECT_DONE/FAILED)
    - Total lines merged
    - Test results
    - Build status
    - Ready for CASCADE Op #6

## IMPORTANT
- This is a CASCADE RERUN fixing integration after P1W1 was rebuilt
- The base is the NEW phase1-wave1-integration (not the original)
- If any conflicts, resolve favoring incoming changes (the rebased effort)
- Document any issues for the orchestrator

## PROJECT_DONE CRITERIA
- ✅ Builder effort merged cleanly
- ✅ Code compiles
- ✅ Tests pass
- ✅ Integration branch pushed
- ✅ Size within limits (<800 lines)
- ✅ Ready for CASCADE to continue with Op #6