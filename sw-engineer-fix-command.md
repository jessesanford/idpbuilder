# SOFTWARE ENGINEER FIX IMPLEMENTATION TASK

## Your State
You are in state: FIX_INTEGRATION_ISSUES

## Critical Information
- **Working Directory**: efforts/phase2/wave1/gitea-registry-client
- **Branch**: idpbuilder-oci-mvp/phase2/wave1/gitea-registry-client
- **Fix Plan Location**: Check FIX_PLAN_LOCATION.txt for the fix plan file

## Required Actions

### 1. Read the fix plan
- The file path is in FIX_PLAN_LOCATION.txt
- This contains specific compilation fixes needed
- Follow ALL instructions exactly

### 2. Implement fixes
The key fixes needed are:
- Add missing import: `github.com/containers/image/v5/copy`
- Remove unused imports (crypto/tls, net/http, build package)
- Replace `image.Copy` with `copy.Image` (2 occurrences)
- Replace `image.Options` with `copy.Options` (2 occurrences)
- Fix docker.Reference type handling

### 3. Verify fixes
- Run: `go build ./pkg/registry/...`
- Ensure no compilation errors
- Run tests if available: `go test ./pkg/registry/...`

### 4. Update status
- Remove FIX_REQUIRED.flag when complete
- Create FIX_COMPLETE.flag with summary
- Commit all changes with clear message

### 5. Important
- Do NOT create new features
- ONLY fix the compilation issues specified
- Follow the fix plan exactly
- Stay in your assigned directory

## Success Criteria
- All compilation errors resolved
- Build passes successfully
- Tests pass (if applicable)
- FIX_COMPLETE.flag created
- Changes committed and pushed
