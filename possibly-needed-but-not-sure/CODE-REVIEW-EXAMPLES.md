# Code Review Examples

## Example 1: Successful Review (E1.1.1 - API Types)

### Review Input
- Branch: `/phase1/wave1/effort1-api-types-core`
- Phase: 1 (APIs only)
- Base: main

### Review Process

```bash
# Step 1: Measure size
$ /workspaces/kcp-shared-tools/tmc-pr-line-counter.sh -c /phase1/wave1/effort1-api-types-core
Lines changed: 456
```

```bash
# Step 2: Verify structure
$ git diff --name-only main.../phase1/wave1/effort1-api-types-core
apis/tmc/v1alpha1/types.go
apis/tmc/v1alpha1/doc.go
apis/tmc/v1alpha1/zz_generated.deepcopy.go
apis/workload/v1alpha1/types.go
apis/workload/v1alpha1/validation.go
# ✅ Only API files, no implementations
```

```bash
# Step 3: Check commits
$ git log --oneline main.../phase1/wave1/effort1-api-types-core
a1b2c3d feat(apis): add TMC configuration types
d4e5f6g feat(apis): add workload placement types
g7h8i9j feat(apis): add validation for placement
j1k2l3m chore: generate deepcopy functions
# ✅ Linear, atomic, well-described
```

```bash
# Step 4: Run tests
$ go test ./apis/... -coverprofile=coverage.out
PASS
coverage: 82.3% of statements
# ✅ Exceeds Phase 1 requirement of 80%
```

### Review Output

```markdown
# Code Review - /phase1/wave1/effort1-api-types-core

## Review Result: ✅ ACCEPTED

### Compliance Summary
- Size: 456 lines ✅ (well under 800 limit)
- Phase: APIs only (Phase 1 compliant) ✅
- Style: KCP compliant ✅
- Documentation: All exports documented ✅
- Testing: 82.3% coverage ✅ (exceeds 80% requirement)
- Build: Clean compilation ✅
- Commits: Linear history, atomic, proper messages ✅

### Strengths
- Clean separation of API types from implementation
- Comprehensive validation using CEL
- Well-structured godoc comments
- Good use of kubebuilder markers

### Architecture Compliance
- Uses standard k8s.io/apimachinery types ✅
- Follows KCP API conventions ✅
- No duplicate type definitions ✅

### Commit Quality
- 4 atomic commits telling clear story
- Proper feat/chore prefixes
- No merge commits

### PR Readiness
Ready for PR creation. Suggest PR title:
"feat(apis): introduce TMC v1alpha1 API types and CRDs"

### Merge Requirements
- Base: main ✅
- No conflicts detected ✅
- Can merge immediately after PR approval
```

---

## Example 2: Review Requiring Fixes (E2.2.1 - TMC Controller)

### Review Input
- Branch: `/phase2/wave2/effort1-tmc-controller`
- Phase: 2 (Controllers)
- Base: phase1-integration

### Review Process

```bash
# Step 1: Size check
$ /workspaces/kcp-shared-tools/tmc-pr-line-counter.sh -c /phase2/wave2/effort1-tmc-controller
Lines changed: 743
# ✅ Under limit
```

```go
// Step 2: Code inspection reveals issues
// File: pkg/controller/tmc/controller.go

func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    // ❌ ISSUE: Hardcoded value
    time.Sleep(30 * time.Second)
    
    // ❌ ISSUE: No error handling
    config := r.getConfig()
    
    // ❌ ISSUE: Function too long (127 lines)
    // ... lots of code ...
}

// ❌ ISSUE: Missing godoc
func processSync(target *v1alpha1.SyncTarget) error {
    // ...
}
```

```bash
# Step 3: Test coverage
$ go test ./pkg/controller/tmc/... -coverprofile=coverage.out
coverage: 71.2% of statements
# ❌ Below Phase 2 requirement of 85%
```

```bash
# Step 4: Commit history
$ git log --oneline phase1-integration.../phase2/wave2/effort1-tmc-controller
abc123 WIP: controller stuff
def456 Merge branch 'phase1-integration' into /phase2/wave2/effort1-tmc-controller
ghi789 fix: more stuff
# ❌ Poor commit messages, merge commit present
```

### Review Output

```markdown
# Code Review - /phase2/wave2/effort1-tmc-controller

## Review Result: ❌ NEEDS FIXES

### Size Analysis
- Lines: 743 ✅ (under 800 limit)

### Critical Issues (MUST FIX)

#### 1. Hardcoded Values
**Location**: `pkg/controller/tmc/controller.go:45`
**Issue**: Hardcoded sleep duration
```go
time.Sleep(30 * time.Second)  // ❌ Hardcoded
```
**Fix Required**:
```go
const DefaultReconcileDelay = 30 * time.Second
time.Sleep(DefaultReconcileDelay)
```

#### 2. Missing Error Handling
**Location**: `pkg/controller/tmc/controller.go:67`
**Issue**: Ignoring potential error from getConfig()
**Fix Required**:
```go
config, err := r.getConfig()
if err != nil {
    return ctrl.Result{}, fmt.Errorf("failed to get config: %w", err)
}
```

#### 3. Function Too Long
**Location**: `pkg/controller/tmc/controller.go:Reconcile`
**Issue**: Function is 127 lines (limit 50)
**Fix Required**: Decompose into smaller functions:
- `validateRequest()`
- `fetchResources()`  
- `computeDesiredState()`
- `applyChanges()`

### Quality Issues

#### Documentation
- [ ] Missing godoc for exported function `ProcessSync` at line 234
- [ ] Missing godoc for exported type `ReconcilerOptions` at line 23

#### Testing
- [ ] Coverage 71.2% - MUST reach 85% for Phase 2
- [ ] Missing test cases for error conditions
- [ ] No tests for concurrent reconciliation

### Commit History Issues

#### Poor Commit Messages
- ❌ "WIP: controller stuff" - Must use proper format
- ❌ "fix: more stuff" - Not descriptive

#### Merge Commit Present
- ❌ Commit def456 is a merge commit
- **Fix**: Rebase onto phase1-integration:
```bash
git rebase phase1-integration
```

### Required Actions Before Acceptance

1. **Fix all hardcoded values** - Use constants
2. **Add error handling** - Check all error returns
3. **Decompose large function** - Max 50 lines per function
4. **Add missing godoc** - All exported symbols
5. **Increase test coverage to 85%** - Add error case tests
6. **Fix commit history**:
   ```bash
   git rebase -i phase1-integration
   # Reword commits with proper messages
   # Remove merge commit
   ```

### Suggested Commit Messages After Rebase
```
feat(controller): implement TMC reconciliation logic
test(controller): add reconciliation error cases  
docs(controller): add godoc for exported types
```

### Re-review Required
Once fixes are complete, request re-review for acceptance.
```

---

## Example 3: Review Requiring Split (E3.2.1 - Resource Sync)

### Review Input
- Branch: `/phase3/wave2/effort1-resource-sync`
- Phase: 3 (Syncer)
- Lines: 1,234 (over limit)

### Review Output

```markdown
# Code Review - /phase3/wave2/effort1-resource-sync

## Review Result: 🔄 NEEDS SPLIT

### Size Analysis
- Current: 1,234 lines
- Limit: 800 lines  
- Overage: 434 lines

### Code Quality Assessment
- Code quality: Good ✅
- Tests present: Yes ✅
- Style compliant: Yes ✅

### Exception Evaluation
**Exception Denied** - Code can be logically separated without breaking functionality.

### Split Plan

#### Part 1: Resource Sync APIs and Types (~300 lines)
**Branch**: `/phase3/wave2/effort1-resource-sync-part1`
**Files**:
- `apis/sync/v1alpha1/resource_types.go`
- `apis/sync/v1alpha1/resource_validation.go`
- `pkg/syncer/resource/interfaces.go`

**Cherry-pick commits**:
```bash
git cherry-pick abc123  # Resource type definitions
git cherry-pick def456  # Validation logic
```

#### Part 2: Resource Sync Implementation (~650 lines)
**Branch**: `/phase3/wave2/effort1-resource-sync-part2`
**Dependencies**: Requires part1
**Files**:
- `pkg/syncer/resource/reconciler.go`
- `pkg/syncer/resource/transformer.go`
- `pkg/syncer/resource/cache.go`

**Cherry-pick commits**:
```bash
git cherry-pick ghi789  # Core reconciler
git cherry-pick jkl012  # Transformation logic
git cherry-pick mno345  # Caching layer
```

#### Part 3: Resource Sync Tests (~284 lines)
**Branch**: `/phase3/wave2/effort1-resource-sync-part3`
**Dependencies**: Requires part1 and part2
**Files**:
- `pkg/syncer/resource/reconciler_test.go`
- `pkg/syncer/resource/transformer_test.go`
- `test/e2e/resource_sync_test.go`

**Cherry-pick commits**:
```bash
git cherry-pick pqr678  # Unit tests
git cherry-pick stu901  # E2E tests
```

### Implementation Instructions

1. **Rename current branch**:
```bash
git branch -m /phase3/wave2/effort1-resource-sync \
            /phase3/wave2/effort1-resource-sync-to-be-split
git push origin --delete /phase3/wave2/effort1-resource-sync
git push -u origin /phase3/wave2/effort1-resource-sync-to-be-split
```

2. **Create split branches** (see EFFORT-SPLIT-CONTINUOUS-EXECUTION-PROTOCOL.md)

3. **Implement sequentially** - DO NOT parallelize

4. **Each split must**:
   - Build independently
   - Pass tests independently
   - Stay under 800 lines

### Expected Outcome
- 3 reviewable PRs under size limit
- Clear dependency chain
- Maintains full functionality when merged

### Next Steps
1. Execute split plan
2. Request review for each part
3. Merge in sequence: part1 → part2 → part3
```

---

## Example 4: Exception Grant (E4.1.1 - Cross-Workspace Controller)

### Review Input
- Branch: `/phase4/wave1/effort1-crossworkspace-ctrl`
- Lines: 923 (over limit)
- Contains critical bug fix

### Review Output

```markdown
# Code Review - /phase4/wave1/effort1-crossworkspace-ctrl

## Review Result: ✅ ACCEPTED WITH EXCEPTION

### Size Analysis
- Current: 923 lines
- Limit: 800 lines
- Overage: 123 lines

### Exception Evaluation

#### Exception Granted: YES

**Reason**: Cross-workspace placement bug fix requires atomic changes across controller, validator, and webhook components. Splitting would break the fix.

**Justification**:
The cross-workspace bug (contrib-tmc issue #4) involves a race condition between:
1. Workspace context propagation in controller
2. Validation webhook checks
3. Admission controller mutations

These three components share state through a synchronized cache that must be updated atomically. Splitting would either:
- Require duplicating the cache logic (400+ lines)
- Break the race condition fix
- Create an intermediate state where placement would fail

### Risk Mitigation

#### Required:
1. **Two additional reviewers** specializing in:
   - Multi-tenancy patterns
   - Kubernetes admission control
   
2. **Enhanced testing**:
   - Race condition tests with -race flag
   - Concurrent placement stress tests
   - Minimum 90% coverage (currently 91% ✅)

3. **Commit structure** for reviewability:
   ```
   Commit 1: Define cross-workspace types and interfaces (120 lines)
   Commit 2: Implement controller with context propagation (280 lines)
   Commit 3: Add validation webhook (185 lines)
   Commit 4: Implement admission mutations (180 lines)
   Commit 5: Add comprehensive tests (158 lines)
   ```

### Code Quality Assessment
- Architecture: Properly integrated ✅
- No duplicates: Verified ✅
- Documentation: Complete ✅
- Testing: 91% coverage ✅
- Bug fix verified: Cross-workspace test passes ✅

### Special Verification
```go
// Test specifically for the bug fix:
func TestCrossWorkspacePlacementBugFixed(t *testing.T) {
    // This used to fail in contrib-tmc
    err := controller.PlaceAcrossWorkspaces(ctx, placement)
    require.NoError(t, err, "Cross-workspace bug must be fixed")
}
// ✅ PASSES
```

### PR Requirements

1. **PR Description must include**:
   - Link to contrib-tmc issue #4
   - Explanation of why split would break fix
   - Performance impact assessment

2. **Review assignments**:
   - Primary: KCP multi-tenancy expert
   - Secondary: Admission control expert
   - Third: Original bug reporter (if available)

3. **Merge conditions**:
   - All three reviewers approve
   - Extended CI run passes (2x normal)
   - No performance regression

### Conclusion
Exception granted due to atomic bug fix requirement. Enhanced review process activated.
```

## Review Patterns Summary

These examples demonstrate:

1. **Successful Review** - Clean code meeting all requirements
2. **Fixes Required** - Specific issues with clear remediation
3. **Split Required** - Over-limit without exception justification  
4. **Exception Granted** - Over-limit but splitting breaks functionality

Each review provides:
- Specific line numbers/files
- Exact fix requirements
- Clear next steps
- Proper documentation