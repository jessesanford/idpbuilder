# SOFTWARE ENGINEER FIX IMPLEMENTATION TASK

🔴🔴🔴 CRITICAL STATE INFORMATION (R295):
YOU ARE IN STATE: FIX_ISSUES
This means you should: Fix the R320 violations identified in CODE-REVIEW-REPORT-20250908-041718.md
🔴🔴🔴

📋 YOUR INSTRUCTIONS (R295):
FOLLOW ONLY: CODE-REVIEW-REPORT-20250908-041718.md
LOCATION: In your effort directory (efforts/phase2/wave1/image-builder)
IGNORE: Any files named *-COMPLETED-*.md (these are from previous fix cycles)

⚠️⚠️⚠️ IMPORTANT:
- The review found R320 violations (stub implementations)
- You MUST fix the stub methods: ListImages(), RemoveImage(), TagImage()
- Either fully implement them OR remove them entirely from the interface
- DO NOT leave any stub implementations (R320 is ZERO TOLERANCE)
⚠️⚠️⚠️

🎯 CONTEXT:
- EFFORT: image-builder (E2.1.1)
- WAVE: 1
- PHASE: 2
- PREVIOUS WORK: Implementation complete at 615 lines, but has R320 violations
- YOUR TASK: Fix ALL R320 violations listed in the code review report

## Critical Information
- **Working Directory**: efforts/phase2/wave1/image-builder
- **Branch**: idpbuilder-oci-build-push/phase2/wave1/image-builder
- **Fix Plan**: CODE-REVIEW-REPORT-20250908-041718.md

## Required Actions

1. **Read the code review report**:
   - File: CODE-REVIEW-REPORT-20250908-041718.md in your effort directory
   - Focus on Section: "CRITICAL BLOCKER: R320 Violation - Stub Implementations Found"
   - The stub methods are: ListImages(), RemoveImage(), TagImage()

2. **Fix the R320 violations**:
   - Option A: Fully implement these methods with real functionality
   - Option B: Remove these methods entirely from the interface and struct
   - DO NOT leave them as stubs or "not implemented"
   - R320 has ZERO TOLERANCE for stub implementations

3. **Verify fixes**:
   - Ensure no methods return "not implemented" errors
   - Ensure no TODO or stub comments remain
   - Run tests to confirm everything still works
   - Verify build succeeds

4. **Update test coverage**:
   - Current coverage: 47.9%
   - Required: 80%
   - Add tests for any new implementations
   - Focus on tar creation and storage edge cases

5. **Complete the fix**:
   - Archive the review report when complete (per R294)
   - Create FIX_COMPLETE.flag with summary
   - Commit all changes with clear message

## Success Criteria
- ALL stub implementations removed (R320 compliance)
- Test coverage increased toward 80%
- Build passes successfully  
- Tests pass
- CODE-REVIEW-REPORT archived as COMPLETED
- FIX_COMPLETE.flag created

## Example Fix Commit Message
```
fix: remove stub implementations to comply with R320

- Removed ListImages(), RemoveImage(), TagImage() stub methods
- These will be implemented in separate future efforts
- Ensures R320 compliance (zero tolerance for stubs)
- Added additional test coverage for tar and storage edge cases
```