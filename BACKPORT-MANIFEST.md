# Backport Manifest
Date: 2025-09-09
State: COORDINATE_BUILD_FIXES
Integration Context: ACTIVE
Rule: R321 - Immediate Backport During Integration

## 🔴🔴🔴 R321 ENFORCEMENT 🔴🔴🔴
ALL fixes MUST be immediately backported to source branches!

## Fixes Requiring Backport

### Fix 1: Docker API Type Issue
- **File**: pkg/kind/cluster_test.go
- **Source Branch**: project-integration
- **SW Engineer**: SWE-1
- **Fix Status**: ASSIGNED
- **Changes to Backport**:
  - [ ] Add container package import
  - [ ] Update ContainerList signature
- **Backport Status**: PENDING
- **Note**: Fix in existing idpbuilder test code

### Fix 2: Format String Issue
- **File**: pkg/util/git_repository_test.go
- **Source Branch**: project-integration
- **SW Engineer**: SWE-2
- **Fix Status**: ASSIGNED
- **Changes to Backport**:
  - [ ] Fix t.Fatalf format string at line 102
- **Backport Status**: PENDING
- **Note**: Fix in existing idpbuilder test code

### Fix 3: Test Infrastructure Setup
- **Files**: scripts/download-test-binaries.sh, pkg/controllers/custompackage/controller_test.go
- **Source Branch**: project-integration
- **SW Engineer**: SWE-3
- **Fix Status**: ASSIGNED
- **Changes to Backport**:
  - [ ] Create binary download script
  - [ ] Update test configuration if needed
- **Backport Status**: PENDING
- **Note**: Infrastructure setup for tests

### Fix 4: Nil Pointer Dereference (Dependent on Fix 3)
- **File**: pkg/controllers/custompackage/controller_test.go
- **Source Branch**: project-integration
- **SW Engineer**: SWE-4 (to be spawned after SWE-3 completes)
- **Fix Status**: NOT_YET_ASSIGNED
- **Changes to Backport**:
  - [ ] Add error handling after testEnv.Start()
  - [ ] Add nil checks for cfg and k8sClient
  - [ ] Add test cleanup function
- **Backport Status**: PENDING
- **Note**: Depends on Fix 3 completion

## Backport Execution Plan
1. Each fix is applied in the project-integration branch
2. After verification, same fixes applied to main if needed
3. All changes tracked with commit references
4. Push to remote after each fix group

## Special Considerations
- All fixes are in existing idpbuilder test code (not our Phase 1/2 implementations)
- Fixes go directly to project-integration branch (already the integration point)
- No separate effort branches need backporting for these fixes
- Main branch update will happen via PR from project-integration

## Tracking
- [ ] All fixes documented
- [ ] Fix commits will be recorded as they complete
- [ ] Ready for BUILD_VALIDATION after all fixes
- [ ] PR to main will include all fixes

## Verification Command
```bash
# After all fixes complete
cd efforts/project/integration-workspace
go test ./pkg/kind ./pkg/util ./pkg/controllers/custompackage -v
```