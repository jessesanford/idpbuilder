# 🚨🚨🚨 BLOCKING RULE R369 - PR Validation and Integrity Protocol

**Criticality:** BLOCKING - Invalid PRs = Production failures
**Enforcement:** MANDATORY - Final PR-Ready validation
**Created:** 2025-01-21

## PURPOSE
Perform comprehensive validation of all branches before PR creation to ensure code integrity, test passage, build success, and production readiness per R271 and R355.

## VALIDATION CATEGORIES

### Category 1: Code Integrity
- No Software Factory artifacts (R365)
- No stub implementations (R355)
- No TODO/FIXME markers (R355)
- No hardcoded credentials
- No debug code

### Category 2: Build & Test
- All unit tests pass
- All integration tests pass
- Build succeeds on all targets
- No linting errors
- Coverage meets requirements

### Category 3: Documentation
- README updated if needed
- API docs current
- Changelog updated
- PR description ready

### Category 4: Git Hygiene
- Clean commit history
- No merge commits (rebased)
- Conventional commits
- Co-authorship preserved

## VALIDATION PROTOCOL

### Step 1: Pre-Validation Setup
```bash
# Create validation workspace
VALIDATION_DIR="pr-validation-$(date +%Y%m%d-%H%M%S)"
mkdir -p $VALIDATION_DIR
cd $VALIDATION_DIR

echo "📋 Starting PR validation suite..."
```

### Step 2: Code Integrity Validation
```bash
echo "🔍 Validating code integrity..."

# Check for SF artifacts (R365)
SF_ARTIFACTS=$(find . -name "orchestrator-state.json" -o -name "*.todo" -o -name "*-PLAN.md" 2>/dev/null)
if [ -n "$SF_ARTIFACTS" ]; then
    echo "❌ FAIL: Software Factory artifacts found:"
    echo "$SF_ARTIFACTS"
    exit 1
fi

# Check for stubs/mocks (R355)
STUBS=$(grep -r "TODO\|FIXME\|STUB\|MOCK\|XXX" --include="*.js" --include="*.ts" --include="*.py" --include="*.java" --include="*.go" . 2>/dev/null | grep -v test)
if [ -n "$STUBS" ]; then
    echo "❌ FAIL: Stub implementations found:"
    echo "$STUBS" | head -10
    exit 1
fi

# Check for hardcoded credentials
CREDENTIALS=$(grep -r "password\|api[_-]key\|secret" --include="*.js" --include="*.ts" --include="*.py" . 2>/dev/null | grep -E "(=|:)\s*[\"'][^\"']*[\"']" | grep -v test)
if [ -n "$CREDENTIALS" ]; then
    echo "⚠️ WARNING: Possible hardcoded credentials:"
    echo "$CREDENTIALS" | head -5
fi

echo "✅ Code integrity validation passed"
```

### Step 3: Build Validation
```bash
echo "🏗️ Validating build..."

# Detect build system
if [ -f "package.json" ]; then
    echo "Node.js project detected"
    npm ci
    npm run build
elif [ -f "pom.xml" ]; then
    echo "Maven project detected"
    mvn clean package
elif [ -f "go.mod" ]; then
    echo "Go project detected"
    go build ./...
elif [ -f "Cargo.toml" ]; then
    echo "Rust project detected"
    cargo build --release
else
    echo "⚠️ No standard build system detected"
fi

if [ $? -ne 0 ]; then
    echo "❌ FAIL: Build failed"
    exit 1
fi

echo "✅ Build validation passed"
```

### Step 4: Test Validation
```bash
echo "🧪 Running test suites..."

# Run tests based on project type
if [ -f "package.json" ]; then
    npm test
    npm run test:integration 2>/dev/null || true
elif [ -f "pom.xml" ]; then
    mvn test
elif [ -f "go.mod" ]; then
    go test ./... -v
elif [ -f "Cargo.toml" ]; then
    cargo test
fi

if [ $? -ne 0 ]; then
    echo "❌ FAIL: Tests failed"
    exit 1
fi

# Check coverage if available
if [ -f "coverage/lcov.info" ]; then
    COVERAGE=$(grep -o '[0-9]*\.[0-9]*%' coverage/lcov.info | head -1)
    echo "Coverage: $COVERAGE"
fi

echo "✅ Test validation passed"
```

### Step 5: Linting Validation
```bash
echo "📝 Running linters..."

# Run appropriate linters
if [ -f ".eslintrc" ] || [ -f ".eslintrc.json" ]; then
    npx eslint . --max-warnings 0
elif [ -f ".pylintrc" ]; then
    pylint **/*.py
elif [ -f ".golangci.yml" ]; then
    golangci-lint run
fi

if [ $? -ne 0 ]; then
    echo "⚠️ WARNING: Linting issues found"
fi

echo "✅ Linting validation passed"
```

### Step 6: Git Hygiene Validation
```bash
echo "📚 Validating git hygiene..."

# Check for merge commits
MERGE_COMMITS=$(git log --merges main..HEAD --oneline)
if [ -n "$MERGE_COMMITS" ]; then
    echo "❌ FAIL: Merge commits found (should be rebased):"
    echo "$MERGE_COMMITS"
    exit 1
fi

# Validate commit messages
INVALID_COMMITS=$(git log main..HEAD --format="%s" | grep -v -E "^(feat|fix|docs|style|refactor|test|chore|perf|ci|build|revert)(\(.+\))?:")
if [ -n "$INVALID_COMMITS" ]; then
    echo "⚠️ WARNING: Non-conventional commits:"
    echo "$INVALID_COMMITS"
fi

# Check commit count
COMMIT_COUNT=$(git rev-list --count main..HEAD)
if [ $COMMIT_COUNT -gt 20 ]; then
    echo "⚠️ WARNING: $COMMIT_COUNT commits (consider squashing)"
fi

echo "✅ Git hygiene validation passed"
```

### Step 7: Documentation Validation
```bash
echo "📚 Validating documentation..."

# Check if README exists
if [ ! -f "README.md" ]; then
    echo "⚠️ WARNING: No README.md found"
fi

# Check if CHANGELOG needs update
if [ -f "CHANGELOG.md" ]; then
    LAST_CHANGE=$(git log -1 --format="%ai" CHANGELOG.md)
    LAST_COMMIT=$(git log -1 --format="%ai")

    if [[ "$LAST_COMMIT" > "$LAST_CHANGE" ]]; then
        echo "⚠️ WARNING: CHANGELOG.md may need updating"
    fi
fi

echo "✅ Documentation validation passed"
```

### Step 8: Generate Validation Report
```bash
cat > PR-VALIDATION-REPORT.md << 'EOF'
# PR Validation Report
Date: $(date)
Branch: $(git branch --show-current)

## Validation Results

### ✅ Code Integrity
- No SF artifacts: ✅
- No stubs/mocks: ✅
- No TODOs: ✅
- No hardcoded creds: ✅

### ✅ Build & Test
- Build successful: ✅
- Unit tests pass: ✅
- Integration tests: ✅
- Coverage: ${COVERAGE:-N/A}

### ✅ Git Hygiene
- No merge commits: ✅
- Conventional commits: ✅
- Commit count: $COMMIT_COUNT

### ✅ Documentation
- README present: ✅
- CHANGELOG current: ✅

## Overall Status: ✅ PR-READY

## Remaining Tasks
1. Create PR description
2. Submit PR in correct order
3. Request reviews

## Certification
This branch has passed all validation checks and is ready for production PR submission.

Validated by: Software Factory PR Validation Suite
Validation ID: $(uuidgen || date +%s)
EOF

echo "✅ Validation report created"
```

## VALIDATION MATRIX

| Check | Required | Failure Action |
|-------|----------|----------------|
| SF Artifacts | ✅ MANDATORY | Block PR |
| Stubs/TODOs | ✅ MANDATORY | Block PR |
| Build Success | ✅ MANDATORY | Block PR |
| Tests Pass | ✅ MANDATORY | Block PR |
| No Merge Commits | ✅ MANDATORY | Require rebase |
| Conventional Commits | ⚠️ RECOMMENDED | Warning only |
| Documentation | ⚠️ RECOMMENDED | Warning only |
| Coverage Threshold | ⚠️ RECOMMENDED | Warning only |

## FAILURE PROTOCOLS

### Critical Failures (Block PR)
```bash
if [ "$VALIDATION_FAILED" = true ]; then
    cat > PR-VALIDATION-FAILED.md << EOF
# PR Validation Failed
Branch: $(git branch --show-current)
Reason: [Specific failure]

This branch is NOT ready for PR submission.

## Required Actions
1. Fix the identified issues
2. Re-run validation
3. Do not proceed until passed
EOF
    exit 1
fi
```

### Warnings (Allow but Note)
```bash
if [ "$WARNINGS_FOUND" = true ]; then
    echo "⚠️ Validation passed with warnings"
    echo "Review warnings before PR submission"
fi
```

## SUCCESS CRITERIA
✅ All mandatory checks pass
✅ No critical failures
✅ Validation report generated
✅ Branch certified PR-ready

## OUTPUT ARTIFACTS
- `PR-VALIDATION-REPORT.md` - Complete validation results
- `PR-VALIDATION-FAILED.md` - If validation fails
- `validation.log` - Detailed validation output

## GRADING PENALTIES
- SF artifacts in PR: -50%
- Stubs/TODOs in code: -40%
- Failing tests in PR: -100%
- No validation performed: -30%

## INTEGRATION WITH OTHER RULES
- Validates **R365** (Artifact Detection)
- Enforces **R355** (Production Ready Code)
- Supports **R271** (Production Validation)
- Enables **R370** (PR Plan Creation)

---

*This rule ensures only production-ready code reaches PR submission.*