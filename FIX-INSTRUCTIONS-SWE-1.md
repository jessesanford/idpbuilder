# Fix Instructions for SW Engineer 1 - Docker API Type Issue
Date: 2025-09-09
Assigned by: orchestrator/COORDINATE_BUILD_FIXES
Priority: HIGH

## Your Assignment
Fix Docker API type compilation error in pkg/kind/cluster_test.go

## Working Directory
```bash
cd /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/project/integration-workspace
```

## Fix Plan Reference
See: FIX-PLAN-BUILD-FAILURES.md Error 1 (lines 23-100)

## Specific Tasks
1. Add container package import to pkg/kind/cluster_test.go
2. Update ContainerList method signature to use container.ListOptions
3. Verify compilation succeeds
4. Run tests to ensure fix is complete

## Code Changes Required

### Change 1: Add container import
**File**: pkg/kind/cluster_test.go  
**Location**: Import section (lines 3-17)  
**Action**: Add line after types import:
```go
"github.com/docker/docker/api/types/container"
```

### Change 2: Update ContainerList signature
**File**: pkg/kind/cluster_test.go  
**Location**: Line 232  
**BEFORE**: 
```go
func (m *DockerClientMock) ContainerList(ctx context.Context, listOptions types.ContainerListOptions) ([]types.Container, error) {
```
**AFTER**:
```go
func (m *DockerClientMock) ContainerList(ctx context.Context, listOptions container.ListOptions) ([]types.Container, error) {
```

## Testing Requirements
```bash
# Step 1: Verify compilation
go test -c ./pkg/kind
# Should compile without "undefined: types.ContainerListOptions" error

# Step 2: Run tests
go test ./pkg/kind -v
# Should pass or at least not fail with compilation errors

# Step 3: Create completion marker
touch FIX-COMPLETE-SWE-1.marker
echo "Docker API type fix completed at $(date)" > FIX-COMPLETE-SWE-1.marker
```

## Success Indicators
- ✅ No compilation errors for pkg/kind
- ✅ Tests compile successfully
- ✅ ContainerList method works with new type
- ✅ FIX-COMPLETE-SWE-1.marker created

## Important Notes
- This is fixing existing idpbuilder test code, not our Phase 1/2 implementations
- Docker SDK v28 reorganized types, requiring this update
- Independent fix - can be done in parallel with other fixes