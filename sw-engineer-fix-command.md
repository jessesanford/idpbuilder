# SOFTWARE ENGINEER FIX IMPLEMENTATION TASK

🔴🔴🔴 CRITICAL STATE INFORMATION (R295):
YOU ARE IN STATE: FIX_ISSUES
This means you should: Fix the issues identified in CODE-REVIEW-REPORT.md
🔴🔴🔴

📋 YOUR INSTRUCTIONS (R295):
FOLLOW ONLY: CODE-REVIEW-REPORT.md
LOCATION: In your effort directory
IGNORE: Any files named *-COMPLETED-*.md (these are archived from previous fix cycles)

⚠️⚠️⚠️ IMPORTANT:
- SPLIT-PLAN-003-COMPLETED-*.md = old, already done
- R509-CASCADE-VIOLATION-REPORT-COMPLETED-*.md = old, already done
- ONLY follow CODE-REVIEW-REPORT.md
⚠️⚠️⚠️

🎯 CONTEXT:
- EFFORT: E1.2.3-image-push-operations-split-003
- WAVE: 2
- PHASE: 1
- PREVIOUS WORK: Initial implementation revealed cascade violation
- YOUR TASK: Properly implement split-003 after split-002 is fixed

## Critical Information
- **Working Directory**: efforts/phase1/wave2/E1.2.3-image-push-operations-split-003
- **Branch**: phase1/wave2/image-push-operations-split-003
- **Fix Plan**: CODE-REVIEW-REPORT.md (already in your directory)

## Required Actions (from CODE-REVIEW-REPORT.md)

### CRITICAL: Fix CASCADE Violations
1. **Wait for split-002 completion**:
   - Split-003 MUST be based on split-002, not split-001
   - Verify split-002 is properly implemented first

2. **Rebase onto split-002**:
   - After split-002 is fixed, rebase split-003 onto it
   - This establishes proper cascade dependency

3. **Implement operations.go**:
   - Should be ~390 lines for operations.go
   - This is the main orchestration component

4. **Add Tests**:
   - Create operations_test.go with unit tests
   - Test the orchestration logic thoroughly

5. **Clean up TODO comments**:
   - Address or remove TODO comments found in the code

## Success Criteria
- Based on split-002 (after it's fixed)
- Contains ONLY operations.go (the orchestrator)
- Under 800 lines total
- Tests added with good coverage
- TODO comments addressed
- Cascade dependency properly established

## Completion
- Archive CODE-REVIEW-REPORT.md when complete per R294
- Create FIX_COMPLETE.flag with summary
- Commit all changes with clear message
