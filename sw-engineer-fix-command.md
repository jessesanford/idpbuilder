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
- SPLIT-PLAN-002-COMPLETED-*.md = old, already done
- ONLY follow CODE-REVIEW-REPORT.md
⚠️⚠️⚠️

🎯 CONTEXT:
- EFFORT: E1.2.3-image-push-operations-split-002
- WAVE: 2
- PHASE: 1
- PREVIOUS WORK: Initial implementation complete, but includes wrong files
- YOUR TASK: Fix boundary violations - rebuild split correctly

## Critical Information
- **Working Directory**: efforts/phase1/wave2/E1.2.3-image-push-operations-split-002
- **Branch**: phase1/wave2/image-push-operations-split-002
- **Fix Plan**: CODE-REVIEW-REPORT.md (already in your directory)

## Required Actions (from CODE-REVIEW-REPORT.md)

### CRITICAL: Fix Split Boundary Violations
1. **Reset to correct base**:
   - Start from split-001 branch (NOT main)
   - Split-002 should ONLY have discovery.go and pusher.go
   - Remove logging.go, progress.go, and operations.go

2. **Correct Implementation**:
   - Total should be ~689 lines (326 for discovery.go + 363 for pusher.go)
   - Do NOT include files from other splits

3. **Add Tests**:
   - Create discovery_test.go with unit tests
   - Create pusher_test.go with unit tests
   - Achieve at least 80% coverage

4. **Verify Size**:
   - Run line counter to ensure under 800 lines
   - Measure from split-001 as base, not main

## Success Criteria
- Split contains ONLY discovery.go and pusher.go
- Based on split-001 branch
- Under 800 lines when measured correctly
- Tests added with good coverage
- All fixes committed and pushed

## Completion
- Archive CODE-REVIEW-REPORT.md when complete per R294
- Create FIX_COMPLETE.flag with summary
- Commit all changes with clear message
