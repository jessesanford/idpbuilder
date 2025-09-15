<<<<<<< HEAD
# SOFTWARE ENGINEER REBASE COMPLETION

🔴🔴🔴 CRITICAL STATE INFORMATION:
REBASE TASK COMPLETED SUCCESSFULLY
🔴🔴🔴

📋 REBASE FINAL STATUS:
TASK: Rebase image-builder branch onto latest phase1/integration
LOCATION: efforts/phase2/wave1/image-builder
STATUS: COMPLETED - 21 of 21 commits processed

✅ REBASE SUMMARY:
- Image-builder branch successfully rebased onto 2c39501
- All Phase 2 image-builder functionality preserved
- Certificate infrastructure conflicts resolved by accepting Phase 1 versions
- Ready for force push to origin

🎯 FINAL CONTEXT:
- EFFORT: image-builder (E2.1.1)
- WAVE: 1, PHASE: 2
- PREVIOUS BASE: Old phase1/integration commit 4f0e259
- NEW BASE: Latest phase1/integration commit 2c39501
- TASK STATUS: REBASE COMPLETE

## Final Status
- **Working Directory**: efforts/phase2/wave1/image-builder
- **Branch**: idpbuilder-oci-build-push/phase2/wave1/image-builder
- **Target**: origin/idpbuilder-oci-build-push/phase1/integration

## Final Actions Required
✅ All 21 commits rebased successfully
✅ All conflicts resolved
✅ Image-builder functionality preserved
✅ Ready for force push to origin

REBASE TASK: COMPLETE
=======
# SOFTWARE ENGINEER FIX IMPLEMENTATION TASK

🔴🔴🔴 CRITICAL STATE INFORMATION (R295):
YOU ARE IN STATE: FIX_INTEGRATION_ISSUES
This means you should: Fix integration issues found during integration testing
🔴🔴🔴

📋 YOUR INSTRUCTIONS (R295):
FOLLOW ONLY: INTEGRATION-REPORT.md
LOCATION: In your effort directory (already present)
IGNORE: Any files named *-COMPLETED-*.md (these are from previous fix cycles)

⚠️⚠️⚠️ IMPORTANT:
- SPLIT-PLAN-COMPLETED-*.md = old, already done
- CODE-REVIEW-REPORT-COMPLETED-*.md = old, already done
- ONLY follow INTEGRATION-REPORT.md
⚠️⚠️⚠️

🎯 CONTEXT:
- EFFORT: gitea-client-split-002
- WAVE: 1
- PHASE: 2
- PREVIOUS WORK: Implementation complete, integration testing revealed issues
- YOUR TASK: Fix ONLY the issues for your effort listed in INTEGRATION-REPORT.md

## Critical Information
- **Working Directory**: efforts/phase2/wave1/gitea-client-split-002
- **Branch**: phase2-wave1-gitea-client-split-002
- **Fix Plan**: INTEGRATION-REPORT.md (R293: Already in your directory)

## Required Actions

1. **Read the integration report**:
   - File: INTEGRATION-REPORT.md in your effort directory
   - Find the section for your effort
   - Follow ALL fix instructions for your effort

2. **Implement fixes (R300 compliance)**:
   - Make ALL fixes in your effort branch
   - NEVER modify the integration branch directly
   - Apply only the changes specified for your effort
   - Install any missing dependencies listed

3. **Archive completed plans (R294)**:
   - If you see any non-archived fix plans, archive them:
   - mv SPLIT-PLAN.md SPLIT-PLAN-COMPLETED-$(date +%Y%m%d-%H%M%S).md (if exists)
   - mv CODE-REVIEW-REPORT.md CODE-REVIEW-REPORT-COMPLETED-$(date +%Y%m%d-%H%M%S).md (if exists)

4. **Verify fixes**:
   - Run all verification commands from INTEGRATION-REPORT.md
   - Ensure build passes
   - Run tests to confirm fixes work

5. **Update status**:
   - Archive INTEGRATION-REPORT.md when complete (R294)
   - Create FIX_COMPLETE.flag with summary
   - Commit all changes with clear message

## Success Criteria
- All issues from INTEGRATION-REPORT.md resolved for your effort
- Build passes successfully
- Tests pass (if applicable)
- INTEGRATION-REPORT.md archived as COMPLETED
- FIX_COMPLETE.flag created
>>>>>>> gitea-split-002/idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-002
