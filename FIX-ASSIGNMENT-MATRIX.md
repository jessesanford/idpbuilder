# Fix Assignment Matrix
Date: 2025-09-09
State: COORDINATE_BUILD_FIXES
Orchestrator: Active Coordination

## Assignment Overview
| SW Engineer | Error | Fix Type | Priority | Dependencies | Files |
|-------------|-------|----------|----------|--------------|-------|
| SWE-1 | Error 1 | Docker API Compilation | HIGH | None | pkg/kind/cluster_test.go |
| SWE-2 | Error 2 | Format String Compilation | HIGH | None | pkg/util/git_repository_test.go |
| SWE-3 | Error 3 | Test Infrastructure Setup | HIGH | None | pkg/controllers/custompackage/controller_test.go |
| SWE-4 | Error 4 | Nil Pointer Fix | MEDIUM | Error 3 | pkg/controllers/custompackage/controller_test.go |

## Parallelization Strategy
### Parallel Group 1 (Independent fixes) - R151 CRITICAL
- SWE-1: Fix Docker API type issue in pkg/kind/cluster_test.go
- SWE-2: Fix format string issue in pkg/util/git_repository_test.go
- SWE-3: Setup test infrastructure (etcd binaries)

### Sequential Group (Dependent fixes)
- After SWE-3 completes → SWE-4 starts (nil pointer fix)

## Timing Requirements (R151)
For parallel spawns:
- Maximum spawn delay: 5 seconds
- All parallel agents must emit timestamps on startup
- Spawn all 3 parallel agents in ONE message

## Success Criteria per Engineer
### SWE-1
- [ ] Import github.com/docker/docker/api/types/container added
- [ ] ContainerList signature updated to use container.ListOptions
- [ ] pkg/kind tests compile successfully
- [ ] Tests pass without Docker API errors

### SWE-2
- [ ] Format string issue fixed in git_repository_test.go line 102
- [ ] Changed to: t.Fatalf("failed to clone repository: %v", err)
- [ ] pkg/util tests compile successfully
- [ ] Tests pass without format errors

### SWE-3
- [ ] Test binaries downloaded or existing cluster configured
- [ ] envtest can start successfully
- [ ] No "missing etcd" errors
- [ ] Controller tests can initialize environment

### SWE-4
- [ ] Proper error handling added after testEnv.Start()
- [ ] Nil checks implemented for cfg and k8sClient
- [ ] Test cleanup function added
- [ ] TestReconcileCustomPkgAppSet passes without panic

## Backport Requirements (R321)
Since these are fixes to the main project integration branch, backports needed to:
- project-integration branch (already there)
- No upstream effort branches affected (fixes are in existing idpbuilder code)

## Verification Command
```bash
cd efforts/project/integration-workspace
go test ./pkg/kind ./pkg/util ./pkg/controllers/custompackage -v
```