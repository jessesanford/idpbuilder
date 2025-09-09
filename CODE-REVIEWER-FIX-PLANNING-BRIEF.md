# Brief for Code Reviewer Fix Planning
Date: 2025-09-09T20:59:00Z
From: orchestrator/ANALYZE_BUILD_FAILURES
To: code-reviewer

## Your Task
Create detailed fix plans for 3 pre-existing test failures in the idpbuilder codebase that are blocking the test suite.

## Context
These are NOT failures introduced by our Phase 1 or Phase 2 implementation. They are pre-existing issues in the original idpbuilder codebase that need to be fixed to achieve clean test execution.

## Input Documents
1. BUILD-ERROR-ANALYSIS.md - Categorized errors with root cause analysis
2. ERROR-TO-EFFORT-MAP.md - Error locations and mapping
3. efforts/project/integration-workspace/test-output-verbose.txt - Raw test output with full error details
4. efforts/project/integration-workspace/PRODUCTION-VALIDATION-REPORT.md - Complete validation report

## Required Outputs
Create a single comprehensive fix plan document:
- **FIX-PLAN-PREEXISTING-TESTS.md**

The fix plan must include for each of the 3 errors:
1. **Error 1: Docker API Type Issue**
   - File: pkg/kind/cluster_test.go:232
   - Current problematic code (exact snippet)
   - Fixed code (exact replacement)
   - Explanation of the fix
   - How to verify the fix works

2. **Error 2: Format String Issue**
   - File: pkg/util/git_repository_test.go:102
   - Current problematic code (exact snippet)
   - Fixed code (exact replacement)
   - Explanation of the fix
   - How to verify the fix works

3. **Error 3: Missing etcd Binary**
   - File: pkg/controllers/custompackage/controller_test.go
   - Analysis of why etcd is missing
   - Solution approach (download binary or update test setup)
   - Exact commands or code changes needed
   - How to verify the fix works

## Priority Order
1. Fix Docker API type issue (compilation error)
2. Fix format string issue (compilation error)
3. Fix missing etcd binary (test execution error)

## Success Criteria
- All 3 errors have implementable fix plans
- Each fix includes exact code snippets (before/after)
- Plans are clear enough for SW Engineers to implement without ambiguity
- Each fix is independently testable
- All fixes can be applied to the project-integration branch

## Working Directory
You should work in: /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/project/integration-workspace

## Important Notes
- These fixes will be applied to the project-integration branch
- No backporting is needed since these are pre-existing issues
- The production binary already builds successfully - we're fixing tests only
- Focus on minimal, targeted fixes that resolve the specific errors

## After Creating Fix Plan
Return control to the orchestrator who will spawn SW Engineers to implement the fixes.