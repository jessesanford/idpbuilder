# INTEGRATION AGENT - EXECUTE P1W2 CASCADE INTEGRATION

You are @agent-integration in INIT state tasked with executing Phase 1 Wave 2 integration as part of CASCADE Operation #2.

## CRITICAL CASCADE CONTEXT
- This is CASCADE Op #2: Recreating P1W2 integration after P1W1 rebuild
- The merge plan has been UPDATED to use single cert-validation branch (not broken splits)
- Both effort branches are rebased onto new P1W1 integration

## YOUR TASK

1. **First emit timestamp**: echo "🕐 Integration Agent P1W2 CASCADE: $(date '+%Y-%m-%d %H:%M:%S')"

2. **Navigate to workspace**:
   ```bash
   cd /home/vscode/workspaces/this-is-not-the-target-repo-this-is-for-orchestrator-planning-only/efforts/phase1/wave2/integration-workspace/repo
   ```

3. **Read the merge plan**: WAVE-MERGE-PLAN.md (already updated)

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
   - First: cert-validation (708 lines)
   - Second: fallback-strategies (660 lines)
   - Use the merge commands from the plan

7. **After each merge**:
   - Run tests: `go test ./...`
   - Check build: `go build ./...`
   - If issues, document them

8. **Final validation**:
   - Total lines should be ~1368
   - All tests passing
   - Code compiles cleanly

9. **Push integration branch**:
   ```bash
   git push origin HEAD:idpbuilder-oci-build-push/phase1-wave2-integration --force-with-lease
   ```

10. **Report completion** with:
    - Integration status
    - Total lines merged
    - Test results
    - Ready for CASCADE Op #3

## IMPORTANT
- This is fixing a CASCADE issue where broken splits caused integration failure
- Use the UPDATED merge plan that references single cert-validation branch
- If any conflicts, resolve favoring incoming changes (the rebased efforts)
- Document any issues for the orchestrator

## SUCCESS CRITERIA
- ✅ Both efforts merged cleanly
- ✅ Code compiles
- ✅ Tests pass
- ✅ Integration branch pushed
- ✅ Ready for CASCADE to continue