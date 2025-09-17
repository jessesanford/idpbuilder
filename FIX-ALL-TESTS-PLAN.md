# Fix All Tests Plan - Phase 1 Wave 3 Upstream Fixes

**Date**: 2025-09-17
**Priority**: CRITICAL - User directive to make ALL tests pass
**Line Limit**: NONE - Focus on making everything work

## Directive
Fix ALL test failures and make the tests pass completely. Don't worry about line count limits.

## Test Failures to Fix

### 1. pkg/util/git_repository_test.go - Missing Functions
**Required implementations**:
```go
// All these functions need to be implemented in pkg/util/git.go or git_repository.go

func CloneRemoteRepoToDir(url string, dir string) error
func CloneRemoteRepoToMemory(url string) (*git.Repository, error)
func CopyTreeToTree(src, dst *object.Tree) error
func ReadWorktreeFile(worktree *git.Worktree, path string) ([]byte, error)
func GetWorktreeYamlFiles(worktree *git.Worktree, dir string) ([]string, error)
```

### 2. pkg/cmd/get/secrets_test.go - Missing Constants and Functions
**Add to pkg/cmd/get/secrets.go**:
```go
const argoCDInitialAdminSecretName = "argocd-initial-admin-secret"
const giteaAdminSecretName = "gitea-admin-secret"

var packages = []string{"argocd", "gitea"}

func printPackageSecrets(pkg string, secrets []v1.Secret)
func printAllPackageSecrets(secrets map[string][]v1.Secret)
```

### 3. pkg/kind/cluster_test.go - Fix API Signatures
**Update NewCluster function signature**:
- Current test expects: `NewCluster(name string)`
- Add method: `func (c *Cluster) getConfig() (*rest.Config, error)`

### 4. pkg/controllers/localbuild/argo_test.go - Add Missing Module
**Create pkg/util/fs package**:
```go
// pkg/util/fs/fs.go
package fs

func FileExists(path string) bool
func CreateDir(path string) error
func ReadFile(path string) ([]byte, error)
func WriteFile(path string, data []byte) error
```

## Implementation Steps

1. Navigate to upstream-fixes directory
2. Implement ALL missing functions with full functionality
3. Add all missing constants and variables
4. Fix all API signature mismatches
5. Create missing packages
6. Test each package individually
7. Run full test suite

## Success Criteria
- `go test ./pkg/...` passes completely
- No compilation errors
- No test failures
- All functions properly implemented (not stubs)

## Note
This is a complete fix - not minimal stubs. Make everything work properly.