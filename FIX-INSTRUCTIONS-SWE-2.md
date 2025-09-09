# Fix Instructions for SW Engineer 2 - Format String Issue
Date: 2025-09-09
Assigned by: orchestrator/COORDINATE_BUILD_FIXES
Priority: HIGH

## Your Assignment
Fix format string compilation error in pkg/util/git_repository_test.go

## Working Directory
```bash
cd /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/project/integration-workspace
```

## Fix Plan Reference
See: FIX-PLAN-BUILD-FAILURES.md Error 2 (lines 102-145)

## Specific Task
Fix the non-constant format string error in t.Fatalf() call at line 102

## Code Change Required

**File**: pkg/util/git_repository_test.go  
**Location**: Line 102  

**BEFORE** (Current incorrect code):
```go
if err != nil {
	t.Fatalf(err.Error())
}
```

**AFTER** (Corrected code):
```go
if err != nil {
	t.Fatalf("failed to clone repository: %v", err)
}
```

## Testing Requirements
```bash
# Step 1: Verify compilation
go test -c ./pkg/util
# Should compile without "non-constant format string" error

# Step 2: Run tests
go test ./pkg/util -v
# Should pass or at least not fail with compilation errors

# Step 3: Create completion marker
touch FIX-COMPLETE-SWE-2.marker
echo "Format string fix completed at $(date)" > FIX-COMPLETE-SWE-2.marker
```

## Success Indicators
- ✅ No compilation errors for pkg/util
- ✅ Tests compile successfully
- ✅ t.Fatalf uses constant format string
- ✅ FIX-COMPLETE-SWE-2.marker created

## Important Notes
- Simple one-line fix
- Go's testing framework requires format strings to be constants
- Independent fix - can be done in parallel with other fixes
- This is fixing existing idpbuilder test code, not our Phase 1/2 implementations