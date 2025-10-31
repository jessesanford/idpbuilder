# R601: Runtime Test Isolation Protocol

**Criticality**: BLOCKING
**Category**: Testing Standards
**Scope**: All Runtime Tests
**Status**: Active (SF 3.0+)

## Purpose

Ensure all runtime tests achieve complete isolation through UUID-based naming and parameterized configuration, preventing conflicts between concurrent test runs and enabling proper test cleanup.

## The Problem

Runtime tests that hardcode project names, repository URLs, or workspace paths create catastrophic failures:

1. **Concurrent Execution Conflicts**: Multiple test runs overwrite each other's workspaces
2. **Remote Branch Collisions**: Tests push to same branch names, corrupting test state
3. **State File Pollution**: Tests share state files, causing cross-contamination
4. **Cleanup Failures**: Cannot identify which files belong to which test run
5. **False Test Failures**: Tests fail due to interference, not actual bugs

**Real-World Impact**: Test 02 hardcoded "hello-world-fullstack" in 20 locations (2025-10-14), breaking all concurrent test execution.

## Rule Requirements

### MANDATORY Variable Usage (BLOCKING)

**ALL runtime tests MUST use framework-provided variables:**

```bash
# ✅ CORRECT - Framework variables
PROJECT_PREFIX="fastapi-hello-sf3-test-${TEST_UUID}"  # Unique per run
TARGET_REPO="https://github.com/jessesanford/fastapi-hello.git"
PLANNING_REPO="https://github.com/jessesanford/software-factory-v3-test-planning-repo"
TEST_WORKSPACE="/tmp/${PROJECT_PREFIX}"
TEST_UUID="$(uuidgen | tr '[:upper:]' '[:lower:]')"

# ✅ CORRECT - Using variables in jq
jq -n --arg project_prefix "$PROJECT_PREFIX" \
      --arg target_repo "$TARGET_REPO" \
      '{
        project_prefix: $project_prefix,
        target_repo_url: $target_repo
      }'

# ❌ FORBIDDEN - Hardcoded values
"project_prefix": "hello-world-fullstack"         # -100% FAILURE
"target_repo_url": "https://github.com/test/..."  # -100% FAILURE
TEST_WORKSPACE="/tmp/sf3-test-$$"                 # -100% FAILURE (PID conflicts)
```

### Required Framework Variables

**Variable Table**:

| Variable | Purpose | Example | Usage |
|----------|---------|---------|-------|
| `$PROJECT_PREFIX` | Unique project name | `fastapi-hello-sf3-test-abc123` | Branch names, state files, directories |
| `$TEST_UUID` | Test run UUID | `abc123de-4567-89ab-cdef-0123456789ab` | Workspace isolation, uniqueness |
| `$TARGET_REPO` | Implementation repository | `https://github.com/jessesanford/fastapi-hello.git` | Effort branches, code pushes |
| `$PLANNING_REPO` | Planning repository | `https://github.com/.../planning-repo` | Planning branch pushes |
| `$TEST_WORKSPACE` | Test workspace path | `/tmp/fastapi-hello-sf3-test-abc123` | All test file operations |

**Derivation Rules**:
```bash
# Framework generates in this order:
TEST_UUID="$(uuidgen | tr '[:upper:]' '[:lower:]')"
PROJECT_PREFIX="fastapi-hello-sf3-test-${TEST_UUID}"
TEST_WORKSPACE="/tmp/${PROJECT_PREFIX}"
```

### FORBIDDEN Practices (BLOCKING)

**❌ These practices result in -100% CATASTROPHIC FAILURE:**

1. **Hardcoded Project Names**:
   ```bash
   # FORBIDDEN
   "project_prefix": "hello-world-fullstack"
   "project_prefix": "my-test-project"
   ```

2. **Hardcoded Repository URLs**:
   ```bash
   # FORBIDDEN
   "target_repo_url": "https://github.com/test/hello-world-fullstack.git"
   ```

3. **PID-Based Workspace Paths**:
   ```bash
   # FORBIDDEN - PID reuse causes conflicts
   TEST_WORKSPACE="/tmp/sf3-test-$$"
   TEST_WORKSPACE="/tmp/test-$(date +%s)"
   ```

4. **Static Branch Names**:
   ```bash
   # FORBIDDEN - Concurrent tests collide
   "branch_name": "software-factory-3.0"
   "branch_name": "test-branch"
   ```

5. **Manual jq Heredocs Without Variables**:
   ```bash
   # FORBIDDEN - Hardcoded values embedded
   jq -n '{
     "project_prefix": "hello-world-fullstack",
     "target_repo": "https://..."
   }'
   ```

### Variable Substitution Patterns

**Pattern 1: jq with --arg (PREFERRED)**
```bash
# ✅ CORRECT - Variables passed via --arg
jq -n \
  --arg project_prefix "$PROJECT_PREFIX" \
  --arg target_repo "$TARGET_REPO" \
  --arg test_uuid "$TEST_UUID" \
  '{
    project_prefix: $project_prefix,
    target_repo_url: $target_repo,
    test_metadata: {
      test_uuid: $test_uuid
    }
  }'
```

**Pattern 2: Template Substitution**
```bash
# ✅ CORRECT - Template + sed substitution
sed -e "s/{{PROJECT_PREFIX}}/$PROJECT_PREFIX/g" \
    -e "s/{{TARGET_REPO}}/$TARGET_REPO/g" \
    -e "s/{{TEST_UUID}}/$TEST_UUID/g" \
    tests/fixtures/templates/orchestrator-state.template.json
```

**Pattern 3: Framework Functions**
```bash
# ✅ CORRECT - Framework provides substitution functions
init_orchestrator_state "$PROJECT_PREFIX" "$TARGET_REPO"
generate_fixture_from_template "orchestrator-state" "$PROJECT_PREFIX"
```

**Anti-Pattern: Inline Heredoc**
```bash
# ❌ FORBIDDEN - Hardcoded values in heredoc
cat > state.json <<EOF
{
  "project_prefix": "hello-world-fullstack",
  "target_repo": "https://github.com/test/repo.git"
}
EOF
```

## Validation Requirements

### Pre-Commit Validation

**Tool**: `tools/validate-test-fixtures.sh`

**Checks**:
1. No hardcoded "hello-world-fullstack" strings
2. No hardcoded "software-factory-3.0" branch names
3. All project names use `$PROJECT_PREFIX` variable
4. All repository URLs use `$TARGET_REPO` or `$PLANNING_REPO`
5. All workspace paths use `$TEST_WORKSPACE`

**Execution**:
```bash
# Validate all test files
bash tools/validate-test-fixtures.sh all

# Validate staged files only
bash tools/validate-test-fixtures.sh staged

# Exit codes:
# 0 = All checks passed
# 1 = Hardcoded values found (BLOCKS COMMIT)
```

### Runtime Validation

**Function**: `validate_test_isolation()` (in runtime-test-framework.sh)

**Checks**:
1. State file `project_prefix` matches `$PROJECT_PREFIX`
2. State file `target_repo_url` matches `$TARGET_REPO`
3. All branch names contain `$PROJECT_PREFIX`
4. Workspace path starts with `/tmp/${PROJECT_PREFIX}`

**Auto-Fail Behavior**:
```bash
# If validation fails, test IMMEDIATELY exits with code 3
validate_test_isolation || {
    echo "❌ FATAL: Test isolation broken - hardcoded values detected"
    exit 3  # Exit code 3 = isolation violation
}
```

## Test Framework Integration

### Framework Setup

**Location**: `tests/runtime-test-framework.sh`

**Required Setup**:
```bash
# Framework MUST provide these variables
export TEST_UUID="$(uuidgen | tr '[:upper:]' '[:lower:]')"
export PROJECT_PREFIX="fastapi-hello-sf3-test-${TEST_UUID}"
export TARGET_REPO="https://github.com/jessesanford/fastapi-hello.git"
export PLANNING_REPO="https://github.com/jessesanford/software-factory-v3-test-planning-repo"
export TEST_WORKSPACE="/tmp/${PROJECT_PREFIX}"

# Create unique workspace
mkdir -p "$TEST_WORKSPACE"
cd "$TEST_WORKSPACE"

# Validate isolation before test starts
validate_test_isolation
```

### Test Cleanup

**Required Cleanup**:
```bash
# Cleanup function MUST use $TEST_WORKSPACE
cleanup_test() {
    # Remove workspace using variable
    rm -rf "$TEST_WORKSPACE"

    # Delete remote branches using $PROJECT_PREFIX
    git push origin --delete "$PROJECT_PREFIX" 2>/dev/null || true

    # Clean up planning repo branch
    cd "$PLANNING_REPO_CLONE"
    git push origin --delete "$PROJECT_PREFIX" 2>/dev/null || true
}

# Register cleanup trap
trap cleanup_test EXIT INT TERM
```

## Concurrent Test Execution

### Isolation Guarantees

**With R601 Compliance**:
- ✅ Multiple test runs can execute simultaneously
- ✅ Each test gets unique workspace (`/tmp/fastapi-hello-sf3-test-{UUID1}`, `/tmp/...{UUID2}`)
- ✅ Each test gets unique remote branches (`$PROJECT_PREFIX/effort-1`, etc.)
- ✅ No state file conflicts (separate workspaces)
- ✅ No cleanup conflicts (UUID identifies ownership)

**Without R601 (Hardcoded Values)**:
- ❌ Tests overwrite each other's workspaces
- ❌ Tests push to same remote branches (state corruption)
- ❌ Tests share state files (cross-contamination)
- ❌ Cleanup removes wrong test's files
- ❌ Random test failures due to race conditions

### Concurrent Execution Test

**Validation**:
```bash
# Run two test instances simultaneously
bash tests/runtime-test-02-wave-start-to-effort-creation.sh &
PID1=$!
bash tests/runtime-test-02-wave-start-to-effort-creation.sh &
PID2=$!

# Both should complete successfully
wait $PID1 && wait $PID2 || echo "❌ Concurrent execution failed"

# Verify unique workspaces exist
ls -d /tmp/fastapi-hello-sf3-test-* | wc -l  # Should be >= 2
```

## Grading Penalties

### Violation Penalties

| Violation | Penalty | Rationale |
|-----------|---------|-----------|
| Hardcoded project name | -100% | Test isolation completely broken |
| Hardcoded repository URL | -100% | Tests push to wrong repository |
| PID-based workspace | -100% | PID reuse causes conflicts |
| Static branch name | -100% | Concurrent tests corrupt each other |
| Missing pre-commit validation | -50% | No enforcement mechanism |
| Missing runtime validation | -50% | No safety net during execution |
| Incomplete variable usage | -25% per instance | Partial isolation still fails |

### Compliance Verification

**100% Compliance Requirements**:
1. ✅ All 7 runtime tests use framework variables
2. ✅ Zero hardcoded project names in test files
3. ✅ Zero hardcoded repository URLs in test files
4. ✅ Pre-commit validation script exists and works
5. ✅ Runtime validation function integrated
6. ✅ Concurrent test execution succeeds
7. ✅ Documentation includes variable usage guide

## Related Rules

- **R600**: Checklist Execution Protocol - DoD validation requirements
- **R287**: TODO Persistence - State file management during tests
- **R506**: Absolute Prohibition on Pre-Commit Bypass - Validation enforcement
- **R602**: Test Fixture Creation Standards - Pattern guidance

## Migration Guide

### Converting Hardcoded Tests

**Step 1: Identify Hardcoded Values**
```bash
# Scan test file for violations
grep -n "hello-world-fullstack" tests/runtime-test-*.sh
grep -n '"target_repo_url": "https://' tests/runtime-test-*.sh
```

**Step 2: Replace with Variables**
```bash
# Before (WRONG):
jq -n '{
  "project_prefix": "hello-world-fullstack",
  "target_repo_url": "https://github.com/test/repo.git"
}'

# After (CORRECT):
jq -n \
  --arg project_prefix "$PROJECT_PREFIX" \
  --arg target_repo "$TARGET_REPO" \
  '{
    project_prefix: $project_prefix,
    target_repo_url: $target_repo
  }'
```

**Step 3: Validate**
```bash
# Run pre-commit validation
bash tools/validate-test-fixtures.sh all

# Run runtime validation
bash tests/runtime-test-XX-your-test.sh
```

## Examples

### Good Test (R601 Compliant)

```bash
#!/usr/bin/env bash
# Test 02: Wave Start → Effort Creation

# Source framework (provides variables)
source tests/runtime-test-framework.sh

# Framework Variables Available:
# - $PROJECT_PREFIX (unique per run)
# - $TARGET_REPO (implementation repo)
# - $TEST_WORKSPACE (unique workspace)

# Validate isolation
validate_test_isolation

# Create state using variables
jq -n \
  --arg project_prefix "$PROJECT_PREFIX" \
  --arg target_repo "$TARGET_REPO" \
  '{
    project_prefix: $project_prefix,
    efforts: [
      {
        branch_name: "\($project_prefix)/effort-1",
        target_repo_url: $target_repo,
        integration_branch: "\($project_prefix)/wave-1-integration"
      }
    ]
  }' > orchestrator-state-v3.json

# Verify isolation during test
grep -q "$PROJECT_PREFIX" orchestrator-state-v3.json || {
    echo "❌ FATAL: State file missing PROJECT_PREFIX"
    exit 3
}

echo "✅ Test complete - isolation verified"
```

### Bad Test (R601 Violation)

```bash
#!/usr/bin/env bash
# Test 02: Wave Start → Effort Creation (WRONG)

# ❌ VIOLATION: Hardcoded project name
jq -n '{
  "project_prefix": "hello-world-fullstack",
  "efforts": [
    {
      "branch_name": "hello-world-fullstack/effort-1",
      "target_repo_url": "https://github.com/test/repo.git"
    }
  ]
}' > orchestrator-state-v3.json

# Result: -100% CATASTROPHIC FAILURE
# - Concurrent tests will overwrite this file
# - Remote branches will collide
# - Test cleanup will fail
```

## Audit Trail

**Created**: 2025-10-14 19:50 UTC
**Reason**: Test 02 hardcoded "hello-world-fullstack" in 20 locations, breaking test isolation
**Reference**: `SOFTWARE-FACTORY-MANAGER-TEST-AUDIT-REPORT.md`
**Impact**: Prevents all future test isolation violations

**Enforcement Checklist**:
- [x] Pre-commit validation script created (tools/validate-test-fixtures.sh)
- [x] Runtime validation function created (validate_test_isolation)
- [x] All 7 tests audited for compliance (Test 02 fixed, others already compliant)
- [x] Documentation created (this rule)
- [x] Grading penalties defined (-100% for violations)

## Version History

- **v1.0.0** (2025-10-14): Initial creation for SF 3.0
  - Defined MANDATORY variable usage requirements
  - Specified FORBIDDEN practices with -100% penalties
  - Created validation mechanisms (pre-commit + runtime)
  - Documented migration guide and examples
