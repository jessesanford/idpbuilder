# Phase 1 Wave 2 Integration - Wave Merge Plan

**Created**: 2025-10-05 17:19:21 UTC
**Agent**: Code Reviewer (WAVE_MERGE_PLANNING state)
**Compliance**: R269 (original branches only), R270 (dependency-based order)

## Overview

This document provides the authoritative merge plan for integrating all Phase 1 Wave 2 efforts into the wave integration branch.

**Integration Branch**: `idpbuilder-push-oci/phase1-wave2-integration`
**Base Branch**: `idpbuilder-push-oci/phase1-wave1-integration`
**Integration Directory**: `/home/vscode/workspaces/idpbuilder-push-oci/efforts/phase1/wave2/integration-workspace`
**Total Branches to Merge**: 6 branches

## R269 Compliance: Branch Selection

Per R269, we use ONLY original effort branches (not integration branches):
- ✅ Original unsplit efforts: Include directly
- ✅ Split effort branches: Include ALL splits, EXCLUDE parent
- ❌ Integration branches: NEVER merge integration branches

### Branches to Merge:
1. ✅ `idpbuilder-push-oci/phase1/wave2/command-structure` (E1.2.1 - no splits)
2. ❌ `idpbuilder-push-oci/phase1/wave2/registry-authentication` (EXCLUDE - was split)
3. ✅ `idpbuilder-push-oci/phase1/wave2/registry-authentication-split-001` (E1.2.2 part 1)
4. ✅ `idpbuilder-push-oci/phase1/wave2/registry-authentication-split-002` (E1.2.2 part 2)
5. ❌ `idpbuilder-push-oci/phase1/wave2/image-push-operations` (EXCLUDE - was split)
6. ✅ `idpbuilder-push-oci/phase1/wave2/image-push-operations-split-001` (E1.2.3 part 1)
7. ✅ `idpbuilder-push-oci/phase1/wave2/image-push-operations-split-002` (E1.2.3 part 2)
8. ✅ `idpbuilder-push-oci/phase1/wave2/image-push-operations-split-003` (E1.2.3 part 3)

**Final Count**: 6 branches to merge (1 original + 5 splits)

## Effort Summary

### E1.2.1: Command Structure (Original - No Splits)
**Branch**: `idpbuilder-push-oci/phase1/wave2/command-structure`
**Status**: Complete, reviewed
**Content**:
- Push command framework with Cobra
- Command flags and validation
- Input validation logic
- Package: `pkg/cmd/push/`

### E1.2.2: Registry Authentication (Split into 2)
**Original Branch**: ~~`registry-authentication`~~ (EXCLUDED per R269)
**Split Branches**:

#### Split 001: Authentication Basics
**Branch**: `idpbuilder-push-oci/phase1/wave2/registry-authentication-split-001`
**Status**: Complete, reviewed
**Content**:
- Authentication interfaces
- Credential handling
- Insecure registry support
- Auth error types
- Partial retry logic
- Packages: `pkg/push/auth/`, `pkg/push/errors/`, `pkg/push/retry/` (partial)

#### Split 002: Retry Logic Complete
**Branch**: `idpbuilder-push-oci/phase1/wave2/registry-authentication-split-002`
**Status**: Complete, reviewed
**Content**:
- Complete retry logic implementation
- Retry logic tests
- Backoff strategies
- Package: `pkg/push/retry/` (complete with tests)

### E1.2.3: Image Push Operations (Split into 3)
**Original Branch**: ~~`image-push-operations`~~ (EXCLUDED per R269)
**Split Branches**:

#### Split 001: Core Operations
**Branch**: `idpbuilder-push-oci/phase1/wave2/image-push-operations-split-001`
**Status**: Complete, reviewed
**Content**:
- Image discovery logic
- Logging utilities
- Core operations
- Progress tracking
- Pusher implementation
- Package: `pkg/push/`

#### Split 002: Operation Tests
**Branch**: `idpbuilder-push-oci/phase1/wave2/image-push-operations-split-002`
**Status**: Complete, reviewed
**Content**:
- Discovery implementation and tests
- Pusher implementation and tests
- Package: `pkg/push/` (with overlapping files)

#### Split 003: Integration Operations
**Branch**: `idpbuilder-push-oci/phase1/wave2/image-push-operations-split-003`
**Status**: Complete, reviewed
**Content**:
- Complete operations implementation
- Operations tests
- Full integration
- Package: `pkg/push/` (with overlapping files)

## R270 Compliance: Dependency Analysis

Per R270, merge order is determined by analyzing branch bases and dependencies:

### Dependency Chain:
```
E1.2.1 (command-structure)
├── No dependencies (foundation)
└── Provides: Command framework, flags, validation

E1.2.2-split-001 (auth basics)
├── Depends on: E1.2.1 (uses command context)
└── Provides: Auth interfaces, credential handling

E1.2.2-split-002 (retry logic)
├── Depends on: E1.2.2-split-001 (extends retry package)
└── Provides: Complete retry implementation with tests

E1.2.3-split-001 (core operations)
├── Depends on: E1.2.1, E1.2.2-split-001 (uses auth + command)
└── Provides: Core push operations

E1.2.3-split-002 (operation tests)
├── Depends on: E1.2.3-split-001 (tests core operations)
└── Provides: Test coverage

E1.2.3-split-003 (integration)
├── Depends on: E1.2.3-split-001, E1.2.3-split-002
└── Provides: Complete integration
```

### Merge Order (Dependency-Based):
1. **E1.2.1** (command-structure) - Foundation, no dependencies
2. **E1.2.2-split-001** (auth basics) - Depends on E1.2.1
3. **E1.2.2-split-002** (retry complete) - Extends E1.2.2-split-001
4. **E1.2.3-split-001** (core ops) - Depends on E1.2.1 + E1.2.2
5. **E1.2.3-split-002** (ops tests) - Depends on E1.2.3-split-001
6. **E1.2.3-split-003** (integration) - Depends on all previous

## Expected Conflicts

### CRITICAL CONFLICT ZONES:

#### 1. E1.2.2-split-001 vs E1.2.2-split-002 (Retry Package Overlap)
**Files Affected**:
- `pkg/push/retry/backoff.go` - Both splits modify
- `pkg/push/retry/errors.go` - Both splits modify
- `pkg/push/retry/retry.go` - Both splits modify

**Analysis**:
- Split-001: Partial implementation (foundation)
- Split-002: Complete implementation + tests

**Resolution Strategy**:
- Use Split-002 version (more complete)
- Split-002 contains all of Split-001's code PLUS tests
- Safe to accept `--theirs` for all retry files

#### 2. E1.2.3 Splits (Major Overlaps)
**Files Affected**:
- `pkg/push/discovery.go` - All three splits
- `pkg/push/pusher.go` - All three splits
- `pkg/push/logging.go` - Split-001 and Split-003
- `pkg/push/operations.go` - Split-001 and Split-003
- `pkg/push/progress.go` - Split-001 and Split-003

**Analysis**:
- Split-001: Initial implementation
- Split-002: Adds tests, may have refinements
- Split-003: Most complete with integration

**Resolution Strategy**:
- Compare each file carefully
- Generally prefer later splits (more refined)
- Ensure test files from Split-002 are retained
- Manual merge may be required for some files

## Detailed Merge Sequence

### Prerequisites
```bash
# Navigate to integration workspace
cd /home/vscode/workspaces/idpbuilder-push-oci/efforts/phase1/wave2/integration-workspace

# Verify on integration branch
git branch --show-current
# Should show: idpbuilder-push-oci/phase1-wave2-integration

# Ensure clean state
git status
# Should show: nothing to commit, working tree clean

# Fetch latest from all remotes
git fetch origin
```

### Step 1: Merge E1.2.1 - Command Structure
**Branch**: `origin/idpbuilder-push-oci/phase1/wave2/command-structure`
**Expected Conflicts**: None (new files only)

```bash
# Merge command structure
git merge origin/idpbuilder-push-oci/phase1/wave2/command-structure \
  --no-ff \
  -m "integrate: E1.2.1 command structure foundation

Creates pkg/cmd/push/ with command framework, flags, and validation.
No conflicts expected - new files only.

Phase: 1, Wave: 2, Effort: E1.2.1"

# Verify merge
git log --oneline -1
ls -la pkg/cmd/push/

# Test compilation
go build ./pkg/cmd/push/
```

**Validation**:
- [ ] Files created: `pkg/cmd/push/flags.go`, `push.go`, `validation.go`
- [ ] No conflicts occurred
- [ ] Code compiles successfully

### Step 2: Merge E1.2.2-split-001 - Auth Basics
**Branch**: `origin/idpbuilder-push-oci/phase1/wave2/registry-authentication-split-001`
**Expected Conflicts**: None (new files)

```bash
# Merge auth basics
git merge origin/idpbuilder-push-oci/phase1/wave2/registry-authentication-split-001 \
  --no-ff \
  -m "integrate: E1.2.2-split-001 authentication basics

Creates pkg/push/auth/, pkg/push/errors/, partial pkg/push/retry/.
Foundation for registry authentication.

Phase: 1, Wave: 2, Effort: E1.2.2, Split: 001"

# Verify merge
ls -la pkg/push/auth/
ls -la pkg/push/retry/

# Test compilation
go build ./pkg/push/...
```

**Validation**:
- [ ] Directories created: `pkg/push/auth/`, `pkg/push/errors/`, `pkg/push/retry/`
- [ ] Auth interfaces present
- [ ] Partial retry logic present
- [ ] No conflicts occurred

### Step 3: Merge E1.2.2-split-002 - Retry Complete (WITH CONFLICTS)
**Branch**: `origin/idpbuilder-push-oci/phase1/wave2/registry-authentication-split-002`
**Expected Conflicts**: YES - retry package files

```bash
# Attempt merge (conflicts expected)
git merge origin/idpbuilder-push-oci/phase1/wave2/registry-authentication-split-002 \
  --no-ff \
  -m "integrate: E1.2.2-split-002 retry logic complete

Completes pkg/push/retry/ with full implementation and tests.
Resolves conflicts by accepting split-002 version (more complete).

Phase: 1, Wave: 2, Effort: E1.2.2, Split: 002"
# This will pause for conflict resolution
```

**Conflict Resolution**:
```bash
# Check conflict status
git status

# Expected conflicts in:
# - pkg/push/retry/backoff.go
# - pkg/push/retry/errors.go
# - pkg/push/retry/retry.go

# Resolution: Accept split-002 version (includes tests + complete impl)
git checkout --theirs pkg/push/retry/backoff.go
git checkout --theirs pkg/push/retry/errors.go
git checkout --theirs pkg/push/retry/retry.go

# Stage resolved files
git add pkg/push/retry/

# Add new test files if any
git add pkg/push/retry/*_test.go

# Complete merge
git commit
```

**Validation**:
- [ ] All retry files use Split-002 version
- [ ] Test files present: `backoff_test.go`, `errors_test.go`, `retry_test.go`
- [ ] Tests pass: `go test ./pkg/push/retry/`
- [ ] No duplicate code

### Step 4: Merge E1.2.3-split-001 - Core Operations
**Branch**: `origin/idpbuilder-push-oci/phase1/wave2/image-push-operations-split-001`
**Expected Conflicts**: None (new files)

```bash
# Merge core operations
git merge origin/idpbuilder-push-oci/phase1/wave2/image-push-operations-split-001 \
  --no-ff \
  -m "integrate: E1.2.3-split-001 core push operations

Creates pkg/push/ core implementation:
- discovery.go, logging.go, operations.go
- progress.go, pusher.go

Phase: 1, Wave: 2, Effort: E1.2.3, Split: 001"

# Verify merge
ls -la pkg/push/

# Test compilation
go build ./pkg/push/
```

**Validation**:
- [ ] Core files created: `discovery.go`, `logging.go`, `operations.go`, `progress.go`, `pusher.go`
- [ ] No conflicts occurred
- [ ] Code compiles successfully

### Step 5: Merge E1.2.3-split-002 - Operation Tests (WITH CONFLICTS)
**Branch**: `origin/idpbuilder-push-oci/phase1/wave2/image-push-operations-split-002`
**Expected Conflicts**: YES - `discovery.go`, `pusher.go`

```bash
# Attempt merge (conflicts expected)
git merge origin/idpbuilder-push-oci/phase1/wave2/image-push-operations-split-002 \
  --no-ff \
  -m "integrate: E1.2.3-split-002 operation tests

Adds test coverage for discovery and pusher.
Resolves conflicts preserving test files.

Phase: 1, Wave: 2, Effort: E1.2.3, Split: 002"
# This will pause for conflict resolution
```

**Conflict Resolution**:
```bash
# Check conflicts
git status

# Expected conflicts:
# - pkg/push/discovery.go
# - pkg/push/pusher.go

# Strategy: Compare versions, likely keep Split-001 base, add tests from Split-002
# Manual review recommended

# For each conflict file:
# 1. Check differences
git diff HEAD:pkg/push/discovery.go origin/idpbuilder-push-oci/phase1/wave2/image-push-operations-split-002:pkg/push/discovery.go

# 2. If Split-002 has enhancements, use it; otherwise keep ours
# Generally: Keep existing implementation
git checkout --ours pkg/push/discovery.go
git checkout --ours pkg/push/pusher.go

# 3. Ensure test files are added
git add pkg/push/discovery_test.go
git add pkg/push/pusher_test.go

# 4. Complete merge
git add pkg/push/
git commit
```

**Validation**:
- [ ] Test files present: `discovery_test.go`, `pusher_test.go`
- [ ] Implementation files conflict resolved
- [ ] Tests pass: `go test ./pkg/push/`
- [ ] No lost functionality

### Step 6: Merge E1.2.3-split-003 - Integration (MAJOR CONFLICTS)
**Branch**: `origin/idpbuilder-push-oci/phase1/wave2/image-push-operations-split-003`
**Expected Conflicts**: YES - Multiple files

```bash
# Attempt merge (major conflicts expected)
git merge origin/idpbuilder-push-oci/phase1/wave2/image-push-operations-split-003 \
  --no-ff \
  -m "integrate: E1.2.3-split-003 integration complete

Final integration of push operations with complete test coverage.
Manual conflict resolution for overlapping implementations.

Phase: 1, Wave: 2, Effort: E1.2.3, Split: 003"
# This will pause for conflict resolution
```

**Conflict Resolution**:
```bash
# Check all conflicts
git status

# Expected conflicts:
# - pkg/push/discovery.go
# - pkg/push/logging.go
# - pkg/push/operations.go
# - pkg/push/progress.go
# - pkg/push/pusher.go

# Strategy: Split-003 is most complete, but verify carefully
# Manual merge required - review each conflict

# For each file:
# 1. View conflict markers
cat pkg/push/discovery.go  # Look for <<<<<<< ======= >>>>>>>

# 2. Manually edit to combine best of both
# OR choose most complete version
# Generally: Split-003 is most refined

# 3. Recommended: Use Split-003 version as base
git checkout --theirs pkg/push/discovery.go
git checkout --theirs pkg/push/logging.go
git checkout --theirs pkg/push/operations.go
git checkout --theirs pkg/push/progress.go
git checkout --theirs pkg/push/pusher.go

# 4. Ensure operations_test.go is added
git add pkg/push/operations_test.go

# 5. Review all changes before committing
git diff --staged

# 6. Complete merge
git add pkg/push/
git commit
```

**Validation**:
- [ ] All implementation files resolved
- [ ] All test files present
- [ ] No duplicate functions
- [ ] Full test suite passes: `go test ./pkg/push/`
- [ ] Integration test file present: `operations_test.go`

## Post-Merge Validation

### Compilation Validation
```bash
# Build all packages
go build ./...

# Should succeed with no errors
echo $?  # Should be 0
```

### Test Validation
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./pkg/cmd/push/...
go test -cover ./pkg/push/...

# Verify key packages
go test -v ./pkg/push/retry/
go test -v ./pkg/push/auth/
go test -v ./pkg/push/
```

### Integration Validation
```bash
# Verify command is accessible
go run main.go push --help

# Should display push command help
```

### Code Quality Checks
```bash
# Check for duplicate functions
go list -f '{{.Name}}' ./... | sort | uniq -d

# Should be empty or expected

# Lint code
golangci-lint run ./...
# OR
go vet ./...
```

### File Inventory Validation
```bash
# List all pkg/cmd/push files
ls -la pkg/cmd/push/

# Expected:
# - flags.go
# - push.go
# - validation.go
# - validation_new.go (if exists)

# List all pkg/push files
ls -la pkg/push/

# Expected:
# - auth/ (directory)
# - errors/ (directory)
# - retry/ (directory)
# - discovery.go
# - discovery_test.go
# - logging.go
# - operations.go
# - operations_test.go
# - progress.go
# - pusher.go
# - pusher_test.go
```

## Integration Success Criteria

### Must Have:
- [ ] All 6 branches successfully merged
- [ ] All conflicts resolved appropriately
- [ ] No duplicate code or functions
- [ ] All packages compile: `go build ./...`
- [ ] All tests pass: `go test ./...`
- [ ] Binary builds: `go build main.go`
- [ ] Command accessible: `go run main.go push --help`

### Quality Gates:
- [ ] Test coverage > 70% for pkg/push
- [ ] Test coverage > 80% for pkg/cmd/push
- [ ] No linting errors
- [ ] No TODO or FIXME markers in production code
- [ ] All error handling present

### Integration Completeness:
- [ ] Command structure present and working
- [ ] Authentication modules integrated
- [ ] Retry logic functional
- [ ] Push operations complete
- [ ] All features from individual efforts present
- [ ] No functionality lost during merges

## Risk Mitigation

### High-Risk Areas

#### 1. Overlapping Files in E1.2.3 Splits
**Risk**: Code loss or duplicate functions
**Mitigation**:
- Careful manual review of each conflicting file
- Use `git diff` extensively to compare versions
- Prefer later splits (more refined)
- Validate with tests after each resolution

#### 2. Test Coverage Gaps
**Risk**: Lost test coverage during conflict resolution
**Mitigation**:
- Track all `*_test.go` files explicitly
- Run coverage reports before/after
- Ensure all test files from all splits are present

#### 3. Dependency Breakage
**Risk**: Import cycles or missing dependencies
**Mitigation**:
- Build after each merge step
- Test inter-package dependencies
- Verify all imports resolve

### Recovery Procedures

#### If Merge Goes Wrong:
```bash
# Abort current merge
git merge --abort

# Reset to pre-merge state
git reset --hard HEAD

# Clean working directory
git clean -fd
```

#### Create Safety Checkpoints:
```bash
# Before risky merges (Step 3, 5, 6)
git branch backup-before-step-N

# If needed to recover
git reset --hard backup-before-step-N
```

## Conflict Resolution Guidelines

### For Each Conflict File:

1. **Identify Nature of Overlap**:
   ```bash
   git diff --ours --theirs <file>
   ```

2. **Decision Matrix**:
   - Identical content → Keep one version (ours or theirs)
   - One extends other → Keep extended version
   - Different implementations → Manual merge required
   - Test files → Always include (both versions if different)

3. **Resolution Steps**:
   ```bash
   # View conflict
   cat <file>

   # Option A: Accept one side
   git checkout --ours <file>   # Keep current
   git checkout --theirs <file> # Use incoming

   # Option B: Manual edit
   nano <file>  # Edit conflict markers

   # Stage resolution
   git add <file>
   ```

4. **Validate Resolution**:
   ```bash
   # Compile affected package
   go build ./path/to/package

   # Run tests
   go test ./path/to/package
   ```

## Notes for Orchestrator

### Critical Instructions:

1. **DO NOT Merge in Parallel**: Sequential merges required due to dependencies
2. **DO NOT Auto-Resolve Conflicts**: Manual review required for Steps 3, 5, 6
3. **DO Validate After Each Step**: Build and test after each merge
4. **DO Create Backups**: Before Steps 3, 5, 6 (risky merges)

### Expected Timeline:
- Step 1 (E1.2.1): 5 minutes
- Step 2 (E1.2.2-001): 5 minutes
- Step 3 (E1.2.2-002): 15 minutes (conflicts)
- Step 4 (E1.2.3-001): 5 minutes
- Step 5 (E1.2.3-002): 20 minutes (conflicts)
- Step 6 (E1.2.3-003): 30 minutes (major conflicts)
- **Total**: ~80 minutes

### Success Indicators:
- ✅ Clean compilation after each step
- ✅ Tests passing after each step
- ✅ No uncommitted changes at end
- ✅ All expected files present
- ✅ Command works: `go run main.go push --help`

### Failure Indicators:
- ❌ Compilation errors after merge
- ❌ Test failures after merge
- ❌ Missing files from effort
- ❌ Duplicate function definitions
- ❌ Import cycle errors

## Appendix: Useful Commands

### Git Operations
```bash
# Preview merge without committing
git merge --no-commit --no-ff <branch>
git merge --abort  # Cancel if too many conflicts

# View differences between branches
git diff <branch1>..<branch2> -- <file>

# List files in a branch
git ls-tree -r <branch> --name-only

# Check merge conflicts in advance
git merge-tree $(git merge-base HEAD <branch>) HEAD <branch>
```

### Conflict Resolution
```bash
# Show conflict status
git status

# View conflict in file
git diff <file>

# Accept one side completely
git checkout --ours <file>    # Keep current version
git checkout --theirs <file>  # Use incoming version

# View both versions side by side
git diff --ours --theirs <file>
```

### Validation Commands
```bash
# Build check
go build ./...

# Test check
go test ./...

# Coverage check
go test -cover ./pkg/push/...

# Lint check
go vet ./...
```

### Recovery Commands
```bash
# Abort merge
git merge --abort

# Reset to previous state
git reset --hard HEAD~1

# Clean untracked files
git clean -fd

# Create backup branch
git checkout -b backup-$(date +%Y%m%d-%H%M%S)
```

## Document Control

**Version**: 2.0
**Created**: 2025-10-05 17:19:21 UTC
**Agent**: Code Reviewer (WAVE_MERGE_PLANNING)
**Status**: Ready for Execution
**Compliance**:
- ✅ R269: Uses only original effort branches (excludes split parents)
- ✅ R270: Dependency-based merge order
- ✅ R383: Metadata in .software-factory (not applicable - this is integration plan)

**Next Step**: Orchestrator executes merges following this plan sequentially

## R405 Compliance Note

This merge plan has been created in compliance with R405 automation requirements. The Integration Agent must:
- Execute these merges WITHOUT manual review
- Continue through the entire sequence unless a CRITICAL error occurs
- Report completion status with CONTINUE-SOFTWARE-FACTORY flag

CONTINUE-SOFTWARE-FACTORY=TRUE
