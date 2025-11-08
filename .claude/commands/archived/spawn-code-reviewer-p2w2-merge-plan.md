# Code Reviewer - P2W2 Merge Plan Creation (CASCADE Op#7)

You are @agent-code-reviewer tasked with creating the merge plan for Phase 2 Wave 2 integration as part of CASCADE Operation #7.

## CONTEXT
- Phase 1 is complete and integrated
- Phase 2 Wave 1 is complete and integrated
- Phase 2 Wave 2 efforts have been rebased onto P2W1 integration
- All R354 reviews for P2W2 rebases have passed
- Integration branch ready: `idpbuilder-oci-build-push/phase2-wave2-integration`

## YOUR IMMEDIATE TASKS

1. **Navigate to workspace**:
   ```bash
   cd /home/vscode/workspaces/this-is-not-the-target-repo-this-is-for-orchestrator-planning-only/efforts/phase2/wave2/integration-workspace
   ```

2. **Verify environment**:
   - Current branch: `idpbuilder-oci-build-push/phase2-wave2-integration`
   - Base: `idpbuilder-oci-build-push/phase2-wave1-integration`
   - Efforts to merge:
     * cli-commands
     * credential-management
     * image-operations

3. **Create merge plan**:
   - Fetch all remote branches
   - Analyze each effort for dependencies
   - Determine optimal merge order
   - Create WAVE-MERGE-PLAN.md

4. **Merge order should be**:
   1. cli-commands (foundation for CLI)
   2. credential-management (auth handling)
   3. image-operations (uses both above)

5. **Save merge plan** as `WAVE-MERGE-PLAN.md` in the workspace

## CRITICAL RULES
- Use standard merge plan format
- Include conflict resolution strategy
- Document any special merge considerations
- Report completion back to orchestrator

## EXPECTED OUTPUT
- WAVE-MERGE-PLAN.md created
- Plan includes all 3 efforts
- Proper merge sequence defined
- Ready for Integration Agent execution

Begin immediately by verifying your workspace and creating the merge plan.