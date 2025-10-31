# R602: Test Fixture Creation Standards

**Criticality**: WARNING
**Category**: Testing / Quality Assurance
**Applies To**: All test developers, Software Factory Manager
**Related Rules**: R601 (Runtime Test Isolation Protocol - BLOCKING), R600 (Checklist Execution Protocol)

---

## Purpose

Define standards and best practices for creating test fixtures in the Software Factory runtime test suite to ensure:
1. Maintainability and readability of test fixtures
2. Consistency across all test files
3. Ease of updating fixtures when framework changes
4. Prevention of hardcoded values and isolation violations
5. Clear separation of test data from test logic

**Why WARNING (not BLOCKING)?**: While fixture quality impacts maintainability and developer productivity, violations won't cause test isolation failures (R601 handles that). Poor fixture patterns increase technical debt but don't break concurrent test execution.

---

## Problem Statement

Test fixtures (JSON state files) can be created in multiple ways:
- External fixture files with runtime substitution (preferred)
- Template-based generation with token replacement (acceptable)
- Inline jq heredocs with variable interpolation (acceptable for simple cases)
- Complex inline jq heredocs with hardcoded values (anti-pattern - violates R601)

Without standards, test fixtures become:
- **Hard to maintain**: Complex inline jq makes fixtures difficult to read/update
- **Error-prone**: Hardcoded values slip through (violates R601)
- **Inconsistent**: Different tests use wildly different fixture creation approaches
- **Fragile**: Framework variable changes require updates in dozens of locations

**This rule provides guidance on fixture creation patterns to maximize maintainability.**

---

## Fixture Creation Patterns

### 🥇 PREFERRED: External Fixtures with Runtime Substitution

**Best for**: Complex state files (>50 lines), frequently reused fixtures, fixtures with multiple variable substitutions

**Pattern**:
```bash
# Store fixture as template with placeholders
tests/fixtures/templates/orchestrator-state-v3.template.json:
{
  "project_prefix": "{{PROJECT_PREFIX}}",
  "target_repo_url": "{{TARGET_REPO}}",
  "test_workspace": "{{TEST_WORKSPACE}}",
  ...
}

# Test file uses substitution function
cat tests/fixtures/templates/orchestrator-state-v3.template.json \
  | sed "s|{{PROJECT_PREFIX}}|$PROJECT_PREFIX|g" \
  | sed "s|{{TARGET_REPO}}|$TARGET_REPO|g" \
  | sed "s|{{TEST_WORKSPACE}}|$TEST_WORKSPACE|g" \
  > "$TEST_WORKSPACE/orchestrator-state-v3.json"
```

**Advantages**:
- ✅ Fixture stored as real JSON (syntax validation, IDE support)
- ✅ Easy to review fixture structure (git diff shows structure, not jq code)
- ✅ Single source of truth for fixture template
- ✅ Variables clearly marked with {{TOKEN}} syntax
- ✅ Reusable across multiple tests
- ✅ Easy to update when framework changes

**Disadvantages**:
- Requires separate template file maintenance
- Substitution logic in test file

**When to use**:
- Fixtures >50 lines
- Fixtures used by multiple tests
- Complex state files with many fields
- Fixtures updated frequently

---

### ✅ ACCEPTABLE: Template System with Batch Generation

**Best for**: Multiple similar fixtures, standardized fixture families, framework-wide fixture generation

**Pattern**:
```bash
# Framework provides template generation functions
# See: tests/fixtures/templates/generate-fixtures.sh

generate_orchestrator_state_fixture() {
  local project_prefix="$1"
  local target_repo="$2"
  local output_file="$3"

  local template="tests/fixtures/templates/orchestrator-state-v3.template.json"
  substitute_framework_variables "$template" "$project_prefix" "$target_repo" > "$output_file"
}

# Test uses generation function
generate_orchestrator_state_fixture \
  "$PROJECT_PREFIX" \
  "$TARGET_REPO" \
  "$TEST_WORKSPACE/orchestrator-state-v3.json"
```

**Advantages**:
- ✅ Centralized fixture generation logic
- ✅ Consistent variable substitution across all tests
- ✅ Easy to batch-update all fixtures when framework changes
- ✅ Reduces code duplication in test files
- ✅ Framework provides validation/testing of fixture quality

**Disadvantages**:
- Requires learning fixture generation API
- Additional abstraction layer

**When to use**:
- Creating standardized fixtures for framework
- Multiple tests need same fixture with different variables
- Framework-wide fixture generation needed

---

### ⚠️ ACCEPTABLE (with caution): Simple Inline jq Updates

**Best for**: Small fixture modifications (<10 lines), single-field updates, test-specific tweaks

**Pattern**:
```bash
# Start with base fixture, modify specific fields
cat tests/fixtures/orchestrator-state-base.json \
  | jq --arg prefix "$PROJECT_PREFIX" \
       --arg repo "$TARGET_REPO" \
       '.project_prefix = $prefix |
        .pre_planned_infrastructure.target_repo = $repo' \
  > "$TEST_WORKSPACE/orchestrator-state-v3.json"
```

**Advantages**:
- ✅ Quick for small modifications
- ✅ Uses framework variables correctly (jq --arg)
- ✅ Based on validated base fixture
- ✅ Self-contained in test file

**Disadvantages**:
- ⚠️ Harder to review fixture structure (need to trace jq pipeline)
- ⚠️ Easy to make mistakes with jq syntax
- ⚠️ Doesn't scale well beyond 3-4 field updates

**When to use**:
- Modifying 1-3 fields in existing fixture
- Test-specific field overrides
- Quick prototyping/debugging

**Requirements**:
- MUST use `jq --arg` for variable substitution (NEVER string interpolation in jq filter)
- MUST start with validated base fixture
- MUST keep jq pipeline <10 lines
- MUST NOT create fixtures from scratch with inline jq

---

### ❌ ANTI-PATTERN: Complex Inline jq Heredocs

**What NOT to do**:
```bash
# ❌ BAD: Complex inline fixture creation with jq heredoc
cat <<'EOF' | jq --arg prefix "$PROJECT_PREFIX" ... > state.json
{
  "version": "3.0",
  "project_prefix": "'"$PROJECT_PREFIX"'",  # ❌ String interpolation
  "state_machine": {
    "current_state": "INIT",
    "previous_state": null,
    ...
  },
  "efforts": [
    {
      "effort_id": "E1",
      "effort_name": "hello-world-fullstack/feature-login",  # ❌ Hardcoded!
      ...
    }
  ]
}
EOF
```

**Why this is bad**:
- ❌ 50+ line fixtures embedded in test files
- ❌ Mixes string interpolation and jq --arg (confusing, error-prone)
- ❌ Easy to miss hardcoded values (violates R601)
- ❌ Impossible to validate JSON syntax before runtime
- ❌ Git diffs show jq code changes, not fixture structure changes
- ❌ Duplication across test files
- ❌ Framework changes require updates in every test

**Penalty**: -10% code quality, fixture refactoring required during code review

---

## Framework Variable Usage (MANDATORY - R601)

**All fixture creation patterns MUST use framework-provided variables:**

| Variable | Purpose | Example |
|----------|---------|---------|
| `$PROJECT_PREFIX` | Unique project identifier | `fastapi-hello-sf3-test-a1b2c3d4` |
| `$TARGET_REPO` | Target repository URL | `https://github.com/jessesanford/fastapi-hello.git` |
| `$TEST_WORKSPACE` | Test workspace path | `/tmp/fastapi-hello-sf3-test-a1b2c3d4` |
| `$TEST_UUID` | Unique test instance ID | `a1b2c3d4-e5f6-7890-abcd-ef1234567890` |

**FORBIDDEN** (violates R601 - BLOCKING):
- ❌ Hardcoded project names: `"hello-world-fullstack"`
- ❌ Hardcoded repository URLs: `"https://github.com/test/hello-world-fullstack.git"`
- ❌ Hardcoded workspace paths: `"/tmp/sf3-test-12345"`

**See R601 for complete variable usage requirements and -100% penalties for violations.**

---

## Fixture Organization Standards

### Directory Structure
```
tests/fixtures/
├── templates/                    # Preferred: External fixture templates
│   ├── orchestrator-state-v3.template.json
│   ├── bug-tracking.template.json
│   ├── integration-containers.template.json
│   └── generate-fixtures.sh     # Template generation functions
├── base/                         # Base fixtures for jq updates
│   ├── orchestrator-state-base.json
│   └── minimal-state.json
└── test-01/                      # Test-specific fixtures (if needed)
    └── custom-state.json
```

### Naming Conventions
- **Templates**: `<filename>.template.json` (use {{TOKEN}} placeholders)
- **Base fixtures**: `<filename>-base.json` (validated, minimal fixture)
- **Test-specific**: `tests/fixtures/test-NN/<filename>.json`

---

## Variable Substitution Methods

### Method 1: sed (Simple, Fast)
```bash
cat template.json \
  | sed "s|{{PROJECT_PREFIX}}|$PROJECT_PREFIX|g" \
  | sed "s|{{TARGET_REPO}}|$TARGET_REPO|g" \
  > output.json
```

**Use when**: Simple token replacement, no complex JSON manipulation

### Method 2: jq --arg (Safe, Robust)
```bash
jq --arg prefix "$PROJECT_PREFIX" \
   --arg repo "$TARGET_REPO" \
   '.project_prefix = $prefix |
    .pre_planned_infrastructure.target_repo = $repo' \
  base.json > output.json
```

**Use when**: Modifying specific JSON fields, need JSON validation

### Method 3: Framework Functions (Recommended)
```bash
source tests/fixtures/templates/generate-fixtures.sh
generate_orchestrator_state_fixture "$PROJECT_PREFIX" "$TARGET_REPO" "output.json"
```

**Use when**: Creating standardized fixtures, following framework patterns

---

## Validation Requirements

### Pre-Commit Validation

All fixture templates MUST pass validation before commit:

```bash
# R601 validation (BLOCKING)
bash tools/validate-test-fixtures.sh staged
# Checks: No hardcoded project names/repos/paths

# JSON syntax validation
for template in tests/fixtures/templates/*.template.json; do
  # Replace tokens with dummy values, validate JSON
  cat "$template" \
    | sed 's/{{[^}]*}}/dummy/g' \
    | jq '.' > /dev/null || echo "INVALID: $template"
done
```

### Runtime Validation

Tests SHOULD validate generated fixtures:

```bash
# After fixture creation
jq '.' "$TEST_WORKSPACE/orchestrator-state-v3.json" > /dev/null || {
  echo "ERROR: Invalid JSON in generated fixture"
  exit 1
}

# Validate required fields present
jq -e '.project_prefix' "$TEST_WORKSPACE/orchestrator-state-v3.json" > /dev/null || {
  echo "ERROR: Missing project_prefix in fixture"
  exit 1
}
```

---

## Migration Guide

### Migrating from Inline jq to External Templates

**Step 1**: Extract inline fixture to template file
```bash
# OLD: Inline jq heredoc in test file
cat <<'EOF' | jq ... > state.json
{ "project_prefix": "...", ... }
EOF

# NEW: Extract to template
cat > tests/fixtures/templates/my-fixture.template.json <<'EOF'
{ "project_prefix": "{{PROJECT_PREFIX}}", ... }
EOF
```

**Step 2**: Update test to use template
```bash
# Replace inline jq with template substitution
cat tests/fixtures/templates/my-fixture.template.json \
  | sed "s|{{PROJECT_PREFIX}}|$PROJECT_PREFIX|g" \
  > "$TEST_WORKSPACE/state.json"
```

**Step 3**: Validate migration
```bash
# Ensure generated fixture matches old inline version
diff <(old_method) <(new_method) || echo "MIGRATION ERROR"
```

---

## Examples: Good vs Bad

### Example 1: Orchestrator State Fixture

**❌ BAD: Complex inline jq heredoc**
```bash
cat <<'EOF' | jq --arg prefix "$PROJECT_PREFIX" ... > state.json
{
  "version": "3.0",
  "project_prefix": "hello-world-fullstack",  # ❌ Hardcoded!
  "state_machine": { ... 50 more lines ... }
}
EOF
```

**✅ GOOD: External template with substitution**
```bash
cat tests/fixtures/templates/orchestrator-state-v3.template.json \
  | sed "s|{{PROJECT_PREFIX}}|$PROJECT_PREFIX|g" \
  > "$TEST_WORKSPACE/orchestrator-state-v3.json"
```

### Example 2: Simple Field Update

**❌ BAD: Recreate entire fixture inline**
```bash
cat <<'EOF' | jq ... > state.json
{ "version": "3.0", "project_prefix": "...", ... 100 lines ... }
EOF
```

**✅ GOOD: Update specific field from base**
```bash
jq --arg prefix "$PROJECT_PREFIX" \
   '.project_prefix = $prefix' \
  tests/fixtures/base/orchestrator-state-base.json \
  > "$TEST_WORKSPACE/orchestrator-state-v3.json"
```

---

## Grading Penalties

| Violation | Severity | Penalty | Fix Required |
|-----------|----------|---------|--------------|
| Complex inline jq (>50 lines) | WARNING | -10% | Refactor to external template |
| No framework variable usage | BLOCKING (R601) | -100% | Add variable substitution |
| Hardcoded project names/repos | BLOCKING (R601) | -100% | Use $PROJECT_PREFIX, $TARGET_REPO |
| Invalid JSON in template | WARNING | -5% | Fix syntax, add validation |
| Inconsistent patterns across tests | WARNING | -5% | Standardize fixture creation |

**Code review will flag:**
- Inline jq heredocs >20 lines (suggest external template)
- Missing framework variables (R601 violation)
- Duplicated fixture creation logic (suggest shared functions)

---

## Enforcement

### Pre-Commit Hooks
```bash
# tools/validate-test-fixtures.sh checks:
# 1. R601 compliance (no hardcoded values) - BLOCKING
# 2. JSON syntax in templates - WARNING
# 3. Fixture size (>100 lines inline = warning) - WARNING
```

### Code Review Checklist
- [ ] Fixtures use preferred pattern when applicable
- [ ] Framework variables used correctly (R601)
- [ ] No complex inline jq (>20 lines)
- [ ] Template files have .template.json extension
- [ ] Validation logic present for generated fixtures

### Runtime Tests
```bash
# validate_test_isolation() in runtime-test-framework.sh
# Auto-fails tests with R601 violations (exit code 3)
```

---

## Related Rules

- **R601 (BLOCKING)**: Runtime Test Isolation Protocol
  - Mandates framework variable usage
  - -100% penalty for hardcoded values
  - Pre-commit and runtime enforcement

- **R600**: Checklist Execution Protocol
  - Test execution standards
  - DoD validation requirements

---

## Summary

**PREFERRED**: External templates with {{TOKEN}} placeholders + sed/jq substitution
**ACCEPTABLE**: Template generation functions, simple jq updates from base fixtures
**AVOID**: Complex inline jq heredocs (>20 lines)
**FORBIDDEN**: Hardcoded values (violates R601 - BLOCKING)

**Key principle**: Fixtures are data, not code. Treat them as first-class artifacts with proper versioning, validation, and reusability.

---

**Rule Version**: 1.0
**Last Updated**: 2025-10-14
**Effective Date**: Immediate (applies to all new test development)
**Enforcement**: Pre-commit validation (WARNING), Code review guidance
