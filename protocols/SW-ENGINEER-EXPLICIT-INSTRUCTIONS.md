# Software Engineering Agent Explicit Instructions

## MANDATORY STARTUP PROCEDURE

**🚨 CRITICAL: BEFORE DOING ANY WORK 🚨**

Every SW Engineering agent MUST follow the startup requirements at:
`/workspaces/[project]/protocols/SW-ENGINEER-STARTUP-REQUIREMENTS.md`

This includes:
1. Print startup timestamp
2. Print all instruction file paths you're using
3. Verify working directory, branch, and remote tracking
4. Confirm task understanding

## Critical Instructions for ALL Development Agents

### MANDATORY FOR EVERY EFFORT

#### 1. Git Commands - EXACT SYNTAX
```bash
# NEVER use pseudo-code. Use these EXACT commands:

# Starting an effort:
git checkout <base-branch>
git checkout -b <effort-branch-name>

# Cherry-picking all commits from a branch:
for commit in $(git rev-list --reverse origin/<source-branch>); do
    git cherry-pick $commit || {
        echo "Conflict on commit $commit"
        # If conflict, prefer source branch version:
        git checkout --theirs .
        git add -A
        git cherry-pick --continue
    }
done

# Cherry-picking specific commits:
git cherry-pick <commit-hash>

# Handling conflicts:
# Option 1: Keep their version (from source branch)
git checkout --theirs <file-path>
# Option 2: Keep our version (current branch)  
git checkout --ours <file-path>
# Option 3: Manual merge
git mergetool
```

#### 2. Validation After EVERY Step
```bash
# After cherry-picking:
git status  # Should show no conflicts
git log --oneline -5  # Verify commits applied

# After code changes:
[BUILD_COMMAND] ./...  # Must compile (e.g., go build, npm run build, mvn compile)
[TEST_COMMAND] ./...   # Must pass (e.g., go test, npm test, mvn test)
[LINT_COMMAND] ./...   # Must have no errors (e.g., golangci-lint, eslint, checkstyle)

# After generating code:
make generate       # or equivalent code generation command
git diff           # Check what was generated
make manifests     # or equivalent manifest generation
make verify        # or equivalent verification command
```

#### 3. Testing Requirements - MUST HAVE
Every effort MUST include tests that prove the requirement is met:

```[LANGUAGE]
// MINIMUM test structure for EVERY feature
// Adapt syntax for your language (Go example shown)

func TestFeatureName(t *testing.T) {
    // Arrange
    // ... setup
    
    // Act
    // ... execute feature
    
    // Assert
    require.NoError(t, err, "Feature must work without errors")
    assert.Equal(t, expected, actual, "Feature must produce expected result")
}

// MUST test error cases
func TestFeatureNameErrors(t *testing.T) {
    // Test invalid input
    // Test edge cases
    // Test concurrent access
}
```

#### 4. File Structure Verification
```bash
# After implementing, verify expected files exist:
ls -la [PACKAGE_DIR]/<feature>/  # Should show expected files
tree [PACKAGE_DIR]/<feature>/    # Should match expected structure

# Verify no duplicate implementations:
find . -name "*.[EXT]" -type f | xargs grep -l "func [FUNCTION_NAME]" | wc -l
# Should be exactly 1 for unique implementations
```

---

## Phase-Specific Explicit Instructions

### PHASE 1 ENHANCEMENTS

#### E1.1.1 - API Types Core - EXPLICIT
```bash
# EXACT commands to run:
git checkout main
git checkout -b /phase1/wave1/effort1-api-types-core

# Cherry-pick from reference implementation
git cherry-pick [COMMIT_HASH]  # API types implementation

# Create EXACT directory structure:
mkdir -p [API_DIR]/[VERSION]
mkdir -p [MODELS_DIR]
mkdir -p [TYPES_DIR]

# Copy base files:
cat > [API_DIR]/[VERSION]/types.[EXT] << 'EOF'
[LANGUAGE_SPECIFIC_CONTENT]
// Base API type definitions
[TYPE_DEFINITIONS]
EOF

# Generate code (adapt for your language/framework):
[CODE_GEN_COMMAND] object:headerFile=hack/boilerplate.[EXT].txt paths=./[API_DIR]/...

# Verify generation worked:
ls [API_DIR]/[VERSION]/[GENERATED_FILES] || exit 1

# Run tests:
[TEST_COMMAND] ./[API_DIR]/... -v
```

#### E1.1.2 - Model Types - EXPLICIT
```bash
# EXACT sequence:
git checkout /phase1/wave1/effort1-api-types-core
git checkout -b /phase1/wave1/effort2-model-types

# Get the model implementation:
git cherry-pick [COMMIT_HASH]  # Model API types

# Update imports (adapt for your language):
sed -i 's|[OLD_PACKAGE_PATH]|[NEW_PACKAGE_PATH]|g' [API_DIR]/[VERSION]/*.[EXT]

# Add validation:
cat > [API_DIR]/[VERSION]/[MODEL]_validation.[EXT] << 'EOF'
[LANGUAGE_SPECIFIC_CONTENT]
// Validation functions
[VALIDATION_LOGIC]
EOF

# Test the validation:
[TEST_COMMAND] ./[API_DIR]/[VERSION]/... -v -run TestValidate
```

### PHASE 2 ENHANCEMENTS

#### E2.1.1 - Base Controller - CORRECTED
```bash
# USE THE RIGHT BRANCH (from correction note):
git checkout phase1-integration
git checkout -b /phase2/wave1/effort1-base-controller

# Get ALL commits from controller patterns implementation:
for commit in $(git log --oneline origin/feature/[CONTROLLER_BRANCH] | head -8 | awk '{print $1}' | tac); do
    echo "Cherry-picking $commit"
    git cherry-pick $commit || {
        # On conflict, keep the implementation version:
        git checkout --theirs [CONTROLLER_DIR]/
        git add -A
        git cherry-pick --continue
    }
done

# Verify we have expected patterns:
grep -r "[EXPECTED_PATTERN]" [CONTROLLER_DIR]/ || exit 1
grep -r "[ANOTHER_PATTERN]" [CONTROLLER_DIR]/ || exit 1
```

### PHASE 3 ENHANCEMENTS

#### E3.1.1 - Core Engine - CRITICAL
```bash
# MOST IMPORTANT - Use correct reference implementation:
git checkout phase2-integration
git checkout -b /phase3/wave1/effort1-core-engine

# Get ALL core engine commits:
BRANCH="origin/feature/[CORE_ENGINE_BRANCH]"
for commit in $(git rev-list --reverse $BRANCH); do
    echo "Applying $commit"
    git cherry-pick $commit || {
        echo "Handling conflict for $commit"
        # ALWAYS prefer reference implementation version:
        git checkout --theirs [CORE_DIR]/
        git add -A
        git cherry-pick --continue
    }
done

# VERIFY no duplicate implementations:
# Should NOT exist:
[ ! -d "[OLD_IMPLEMENTATION_DIR]" ] || { echo "ERROR: Old implementation location exists!"; exit 1; }
# Should exist:
[ -d "[NEW_IMPLEMENTATION_DIR]" ] || { echo "ERROR: New implementation location missing!"; exit 1; }

# Verify critical features are present:
grep -r "[CRITICAL_FEATURE]" [CORE_DIR]/ || { echo "ERROR: Missing critical feature!"; exit 1; }
```

### PHASE 4 ENHANCEMENTS

#### E4.1.1 - Bug Fix - CRITICAL
```bash
# This MUST fix the identified bug:
git checkout phase3-integration  
git checkout -b /phase4/wave1/effort1-bugfix

# Get the controller:
git cherry-pick $(git log --oneline origin/feature/[BUGFIX_BRANCH] | awk '{print $1}')

# Apply the bug fix:
cat > [PACKAGE_DIR]/[FEATURE]/bugfix.[EXT] << 'EOF'
[LANGUAGE_SPECIFIC_CONTENT]
// Bug fix implementation
[BUG_FIX_CODE]
EOF

# TEST the fix works:
cat > test/[FEATURE]/bugfix_test.[EXT] << 'EOF'
[LANGUAGE_SPECIFIC_CONTENT]
// Test that verifies the bug is fixed
[BUG_FIX_TESTS]
EOF

[TEST_COMMAND] ./test/[FEATURE]/... -v -run TestBugFixed
# If this fails, the bug is NOT fixed!
```

#### E4.4.1 - New Feature Plugin - NEW
```bash
# This is NEW - not in any existing branch:
git checkout phase4/wave3-integration
git checkout -b /phase4/wave4/effort1-new-feature-plugin

# Create the new feature implementation:
mkdir -p [PACKAGE_DIR]/[FEATURE]

# Implement the feature:
cat > [PACKAGE_DIR]/[FEATURE]/[FEATURE].[EXT] << 'EOF'
[LANGUAGE_SPECIFIC_CONTENT]
// New feature implementation
[FEATURE_CODE]
EOF

# Fix all imports and dependencies:
find [PACKAGE_DIR]/[FEATURE] -name "*.[EXT]" -type f -exec sed -i 's|[OLD_PACKAGE]|[NEW_PACKAGE]|g' {} \;

# Verify it compiles:
[BUILD_COMMAND] ./[PACKAGE_DIR]/[FEATURE]/...

# Create integration:
cat > [PACKAGE_DIR]/[FEATURE]/integration.[EXT] << 'EOF'
[LANGUAGE_SPECIFIC_CONTENT]
// Integration with main system
[INTEGRATION_CODE]
EOF
```

### PHASE 5 ENHANCEMENTS

#### E5.5.2 - Final Validation - MUST ALL PASS
```bash
# This script MUST exit 0 or the implementation is incomplete:
cat > test/final-validation.sh << 'EOF'
#!/bin/bash
set -e  # Exit on any error

echo "=== FINAL VALIDATION ==="

# 1. Build test
echo "Testing build..."
make build-all || { echo "FAIL: Build broken"; exit 1; }

# 2. Unit tests with coverage
echo "Testing units..."
[TEST_WITH_COVERAGE_COMMAND] || { echo "FAIL: Unit tests"; exit 1; }
COVERAGE=$([EXTRACT_COVERAGE_COMMAND])
if (( $(echo "$COVERAGE < 80" | bc -l) )); then
    echo "FAIL: Coverage $COVERAGE% < 80%"
    exit 1
fi

# 3. Integration tests
echo "Testing integration..."
make test-integration || { echo "FAIL: Integration tests"; exit 1; }

# 4. E2E tests
echo "Testing e2e..."
make test-e2e || { echo "FAIL: E2E tests"; exit 1; }

# 5. Bug fix verification
echo "Testing bug fix..."
[TEST_COMMAND] ./test/e2e -run TestBugFix -v || { echo "FAIL: Bug still present"; exit 1; }

# 6. Feature verification
echo "Verifying features..."
FEATURES=(
    "[FEATURE1_PATH]"
    "[FEATURE2_PATH]" 
    "[FEATURE3_PATH]"
    "[FEATURE4_PATH]"
)

for feature in "${FEATURES[@]}"; do
    if [ ! -e "$feature" ]; then
        echo "FAIL: Missing feature $feature"
        exit 1
    fi
    echo "✓ $feature"
done

# 7. No duplicate implementations
echo "Checking for duplicates..."
IMPL_COUNT=$(find . -name "*.[EXT]" -type f | xargs grep -l "[UNIQUE_PATTERN]" | wc -l)
if [ "$IMPL_COUNT" -ne 1 ]; then
    echo "FAIL: Found $IMPL_COUNT implementations (should be 1)"
    exit 1
fi

# 8. Lint check
echo "Checking lint..."
[LINT_COMMAND] ./... || { echo "FAIL: Lint errors"; exit 1; }

# 9. Generated code check
echo "Checking generated code..."
make generate
git diff --exit-code || { echo "FAIL: Generated code not committed"; exit 1; }

# 10. Manifest/Schema validation
echo "Checking manifests..."
make manifests
git diff --exit-code || { echo "FAIL: Manifests not up to date"; exit 1; }

echo "=== ALL VALIDATION PASSED ==="
echo "Ready to merge to main!"
EOF

chmod +x test/final-validation.sh
./test/final-validation.sh
```

---

## Agent Success Criteria

### For EVERY Effort, Agent MUST:

1. **Start with exact git commands** (not pseudo-code)
2. **Verify each cherry-pick succeeded** 
3. **Run tests after implementation**
4. **Verify no duplicates introduced**
5. **Check compilation: `[BUILD_COMMAND] ./...`**
6. **Check tests: `[TEST_COMMAND] ./...`**
7. **Check lint: `[LINT_COMMAND] ./...`**
8. **Commit with descriptive message**

### Red Flags - STOP if:
- Cherry-pick has unresolved conflicts
- Tests are failing
- Build is broken
- Duplicate implementations detected
- Coverage drops below 80%

### Rollback Procedure if Blocked:
```bash
# If effort fails:
git status  # Check state
git reset --hard HEAD  # Reset to clean state
git checkout <previous-branch>  # Go back
# Report issue to orchestrator with:
# - Exact error message
# - Conflicting files
# - Suggested resolution
```

---

## Orchestrator Validation Gates

Before marking any effort complete, orchestrator MUST verify:

```bash
# For effort branch:
git checkout <effort-branch>

# 1. Builds
[BUILD_COMMAND] ./... || EFFORT_FAILED

# 2. Tests pass
[TEST_COMMAND] ./... || EFFORT_FAILED  

# 3. No unresolved conflicts
git status | grep -q "conflict" && EFFORT_FAILED

# 4. Expected files exist
[ -f <expected-file> ] || EFFORT_FAILED

# 5. No regression
make test || EFFORT_FAILED
```

Only after ALL checks pass can effort be marked complete.

---

## Language/Framework Placeholders

Replace these placeholders with language-specific values:

- `[LANGUAGE]`: Programming language (Go, Java, TypeScript, Python, etc.)
- `[EXT]`: File extension (.go, .java, .ts, .py, etc.)
- `[BUILD_COMMAND]`: Build command (go build, npm run build, mvn compile, etc.)
- `[TEST_COMMAND]`: Test command (go test, npm test, mvn test, pytest, etc.)
- `[LINT_COMMAND]`: Lint command (golangci-lint, eslint, checkstyle, flake8, etc.)
- `[CODE_GEN_COMMAND]`: Code generation command (controller-gen, protoc, etc.)
- `[PACKAGE_DIR]`: Package/source directory (pkg, src, lib, etc.)
- `[API_DIR]`: API directory structure
- `[CONTROLLER_DIR]`: Controller/service directory
- `[CORE_DIR]`: Core implementation directory
- `[TEST_WITH_COVERAGE_COMMAND]`: Test with coverage command
- `[EXTRACT_COVERAGE_COMMAND]`: Command to extract coverage percentage
- `[UNIQUE_PATTERN]`: Pattern to identify unique implementations

This template provides a framework that can be adapted to any programming language or development stack while maintaining the critical validation and workflow patterns.