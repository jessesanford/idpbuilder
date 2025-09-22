# CODE REVIEWER - UPDATE P1W2 MERGE PLAN

You are @agent-code-reviewer in WAVE_MERGE_PLANNING state tasked with UPDATING the existing Phase 1 Wave 2 merge plan.

## CRITICAL CONTEXT
- The existing merge plan at efforts/phase1/wave2/integration-workspace/repo/WAVE-MERGE-PLAN.md is OUTDATED
- It references cert-validation-split-001 and cert-validation-split-002 which are DELETED
- You must UPDATE it to use the single cert-validation branch instead

## YOUR TASK

1. **First emit timestamp**: echo "🕐 Code Reviewer P1W2 merge plan update: $(date '+%Y-%m-%d %H:%M:%S')"

2. **Navigate to workspace**:
   - cd /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase1/wave2/integration-workspace/repo

3. **Read existing merge plan** to understand structure

4. **Update the merge plan** to:
   - Remove ALL references to cert-validation-split-001 and cert-validation-split-002
   - Use single branch: idpbuilder-oci-build-push/phase1/wave2/cert-validation (708 lines)
   - Keep fallback-strategies as is

5. **Correct merge sequence**:
   - Effort 1: cert-validation (single branch, 708 lines)
   - Effort 2: fallback-strategies (660 lines)
   - Total: ~1368 lines

6. **Save updated plan** as WAVE-MERGE-PLAN.md (overwrite)

7. **Verify** the plan is correct

8. **Report completion** with summary of changes

## IMPORTANT
- Do NOT create a new plan from scratch
- UPDATE the existing one to fix the branch references
- Preserve all other structure and metadata
- This is fixing a CASCADE issue where old split references need updating

