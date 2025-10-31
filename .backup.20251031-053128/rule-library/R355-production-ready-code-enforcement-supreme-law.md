# 🔴🔴🔴 SUPREME LAW R355: PRODUCTION READY CODE ENFORCEMENT 🔴🔴🔴

## 🚨🚨🚨 BLOCKING - ZERO TOLERANCE - AUTOMATIC FAILURE 🚨🚨🚨

**THIS IS SUPREME LAW #5 - SUPERSEDES ALL OTHER CODE QUALITY RULES**

**VIOLATION = -100% GRADE = AUTOMATIC PROJECT FAILURE**

## ABSOLUTE PROHIBITION - NO EXCEPTIONS

### ❌ FORBIDDEN IN ALL PRODUCTION CODE:

1. **STUBS** - Any unimplemented functionality
2. **MOCKS** - Any mock implementations (except test files)
3. **HARDCODED CREDENTIALS** - Any embedded passwords/tokens
4. **STATIC VALUES** - Any non-configurable constants
5. **PLACEHOLDER CODE** - Any temporary implementations
6. **TODO/FIXME** - Any incomplete work markers
7. **FAKE DATA** - Any hardcoded test data
8. **DUMMY IMPLEMENTATIONS** - Any non-functional code

## 🔴🔴🔴 MANDATORY DETECTION PROTOCOL 🔴🔴🔴

### SW ENGINEERS MUST CHECK BEFORE COMMIT:
```bash
# CRITICAL: Run ALL these checks BEFORE committing
cd $EFFORT_DIR

# Check for hardcoded credentials (ZERO TOLERANCE)
grep -r "password.*=.*['\"]" --exclude-dir=test --exclude-dir=.git --include="*.go" --include="*.py" --include="*.js" --include="*.ts" --include="*.java"
grep -r "username.*=.*['\"]" --exclude-dir=test --exclude-dir=.git --include="*.go" --include="*.py" --include="*.js" --include="*.ts" --include="*.java"
grep -r "token.*=.*['\"]" --exclude-dir=test --exclude-dir=.git --include="*.go" --include="*.py" --include="*.js" --include="*.ts" --include="*.java"
grep -r "secret.*=.*['\"]" --exclude-dir=test --exclude-dir=.git --include="*.go" --include="*.py" --include="*.js" --include="*.ts" --include="*.java"
grep -r "apikey.*=.*['\"]" --exclude-dir=test --exclude-dir=.git --include="*.go" --include="*.py" --include="*.js" --include="*.ts" --include="*.java"

# Check for stubs/mocks outside tests (AUTOMATIC FAILURE)
grep -r "stub\|mock\|fake\|dummy" --exclude-dir=test --exclude-dir=tests --exclude-dir=.git --include="*.go" --include="*.py" --include="*.js" --include="*.ts" --include="*.java"

# Check for incomplete work markers (NOT ALLOWED)
grep -r "TODO\|FIXME\|HACK\|XXX\|INCOMPLETE\|TEMPORARY" --exclude-dir=test --exclude-dir=.git --include="*.go" --include="*.py" --include="*.js" --include="*.ts" --include="*.java"

# Check for not implemented patterns (CRITICAL BLOCKER)
grep -r "not.*implemented\|unimplemented\|NotImplementedError" --exclude-dir=test --exclude-dir=.git --include="*.go" --include="*.py" --include="*.js" --include="*.ts" --include="*.java"

# Check for panic/throw placeholders (FORBIDDEN)
grep -r "panic.*TODO\|panic.*unimplemented\|throw.*TODO" --exclude-dir=test --exclude-dir=.git --include="*.go" --include="*.py" --include="*.js" --include="*.ts" --include="*.java"

# Check for hardcoded URLs/endpoints (MUST BE CONFIGURABLE)
grep -r "http://localhost\|127.0.0.1\|192.168" --exclude-dir=test --exclude-dir=.git --include="*.go" --include="*.py" --include="*.js" --include="*.ts" --include="*.java"
```

### CODE REVIEWERS MUST VERIFY:
```bash
# Run COMPLETE production readiness scan
cd $EFFORT_DIR

# Create violation report
echo "=== R355 PRODUCTION READINESS SCAN ===" > production-violations.txt
echo "Timestamp: $(date)" >> production-violations.txt
echo "Directory: $(pwd)" >> production-violations.txt
echo "" >> production-violations.txt

# Check each violation category
for pattern in "password.*=.*['\"]" "username.*=.*['\"]" "token.*=.*['\"]" "stub\|mock\|fake" "TODO\|FIXME" "not.*implemented"; do
    echo "Checking: $pattern" >> production-violations.txt
    grep -r "$pattern" --exclude-dir=test --exclude-dir=.git --include="*.go" --include="*.py" --include="*.js" --include="*.ts" >> production-violations.txt 2>/dev/null || echo "  ✓ None found" >> production-violations.txt
    echo "" >> production-violations.txt
done

# ANY violations = FAILED REVIEW
if grep -q "password\|username\|token\|stub\|mock\|TODO\|FIXME\|not.*implemented" production-violations.txt; then
    echo "🚨🚨🚨 R355 VIOLATION: NON-PRODUCTION CODE DETECTED!" >> production-violations.txt
    echo "REVIEW STATUS: FAILED - CRITICAL BLOCKERS FOUND" >> production-violations.txt
    cat production-violations.txt
    exit 355
fi
```

## 🔴🔴🔴 GRADING PENALTIES 🔴🔴🔴

### SW ENGINEERS:
- **-100%**: ANY hardcoded credential in code
- **-100%**: ANY stub/mock in production code
- **-75%**: Static values instead of configuration
- **-50%**: TODO/FIXME markers in committed code
- **-100%**: Arguing these are acceptable

### CODE REVIEWERS:
- **-100%**: Passing code with hardcoded credentials
- **-100%**: Passing code with stubs/mocks
- **-75%**: Missing static value violations
- **-50%**: Classifying as "minor issues"
- **-100%**: Not running detection protocol

### ORCHESTRATOR:
- **-100%**: Not enforcing R355 on all agents
- **-100%**: Allowing merge with violations
- **-100%**: Not tracking R355 compliance

## 🚨🚨🚨 VIOLATION EXAMPLES - AUTOMATIC FAILURES 🚨🚨🚨

### ❌ HARDCODED CREDENTIALS (SECURITY BREACH)
```go
// ❌ AUTOMATIC FAILURE - SECURITY VIOLATION
password := "admin123"
username := "testuser"
apiKey := "sk-1234567890abcdef"
dbConn := "postgres://admin:password@localhost/db"
```

### ❌ STUB IMPLEMENTATIONS (PRODUCTION BROKEN)
```go
// ❌ AUTOMATIC FAILURE - NON-FUNCTIONAL CODE
func GetUserData() (*User, error) {
    // TODO: implement this
    return nil, fmt.Errorf("not implemented")
}
```

### ❌ MOCK IN PRODUCTION (FAKE FUNCTIONALITY)
```go
// ❌ AUTOMATIC FAILURE - MOCK OUTSIDE TESTS
type MockUserService struct{}

func NewUserService() *MockUserService {
    return &MockUserService{}
}
```

### ❌ STATIC VALUES (NOT CONFIGURABLE)
```go
// ❌ AUTOMATIC FAILURE - HARDCODED CONFIG
const API_URL = "http://localhost:8080"
const MAX_RETRIES = 3
const TIMEOUT = 30
```

### ❌ PLACEHOLDER CODE (INCOMPLETE WORK)
```python
# ❌ AUTOMATIC FAILURE - TEMPORARY CODE
def process_payment(amount):
    # FIXME: Add real payment processing
    print("Payment would be processed here")
    return True  # HACK: Always succeed for now
```

## ✅ REQUIRED PATTERNS - PRODUCTION READY

### ✅ CONFIGURATION-BASED CREDENTIALS
```go
// ✅ CORRECT - From environment/config
password := os.Getenv("DB_PASSWORD")
username := config.GetString("auth.username")
apiKey := viper.GetString("api.key")
dbConn := fmt.Sprintf("postgres://%s:%s@%s/%s",
    os.Getenv("DB_USER"),
    os.Getenv("DB_PASS"),
    os.Getenv("DB_HOST"),
    os.Getenv("DB_NAME"))
```

### ✅ COMPLETE IMPLEMENTATIONS
```go
// ✅ CORRECT - Full implementation
func GetUserData(userID string) (*User, error) {
    if userID == "" {
        return nil, errors.New("user ID required")
    }

    user, err := db.QueryUser(userID)
    if err != nil {
        return nil, fmt.Errorf("query user: %w", err)
    }

    return user, nil
}
```

### ✅ REAL SERVICES (NOT MOCKS)
```go
// ✅ CORRECT - Real service implementation
type UserService struct {
    db     Database
    cache  Cache
    logger Logger
}

func NewUserService(db Database, cache Cache, logger Logger) *UserService {
    return &UserService{
        db:     db,
        cache:  cache,
        logger: logger,
    }
}
```

### ✅ CONFIGURABLE VALUES
```go
// ✅ CORRECT - Configuration-driven
var (
    apiURL     = config.GetString("api.url")
    maxRetries = config.GetInt("api.max_retries")
    timeout    = config.GetDuration("api.timeout")
)
```

## 🔴🔴🔴 ENFORCEMENT HIERARCHY 🔴🔴🔴

1. **R355 OVERRIDES ALL OTHER RULES**
2. **NO EXCEPTIONS FOR:**
   - "Just getting it working first"
   - "Will fix in next PR"
   - "Size limit constraints"
   - "Time pressure"
   - "POC/Demo code"

3. **APPLIES TO ALL CODE:**
   - Main implementation
   - CLI tools
   - Scripts
   - Configuration loaders
   - Helper functions
   - ALL production code

## INTEGRATE_WAVE_EFFORTS WITH OTHER RULES

- **SUPERSEDES R320**: R355 includes and extends R320's stub detection
- **OVERRIDES R007**: Size limits NEVER excuse production violations
- **MANDATORY FOR R307**: Independent mergeability requires production-ready code
- **REQUIRED BY R323**: Final artifacts must be production-ready
- **TRIGGERS R332**: ALL violations MUST trigger mandatory bug filing (see below)

## 🚨🚨🚨 INTEGRATE_WAVE_EFFORTS WITH R332 - MANDATORY BUG FILING 🚨🚨🚨

**ALL R355 violations found in code reviews MUST trigger R332 Bug Filing Protocol.**

### NO EXCEPTIONS - NO DEFERRALS

**EXCEPTION: TODO Comments with Valid Plan Evidence (R332 TODO Acceptance Criteria)**

TODOs are ONLY acceptable if Code Reviewer provides ALL evidence per R332:
1. **Plan File Path** (exact)
2. **Line Numbers** (exact)
3. **Effort ID** (specific)
4. **Future Phase/Wave** (verified)
5. **Grep Output** (showing plan text)
6. **Phase/Wave Comparison** (current vs planned)
7. **Decision with Evidence** (not vague reasoning)

**See R332 TODO Acceptance Criteria for complete verification protocol.**

**PROHIBITED RESPONSES (R332 Violation = -100%):**
- ❌ "TODO comment acceptable for Phase X" (without evidence)
- ❌ "Pre-existing stub from earlier effort"
- ❌ "Well-documented future enhancement" (vague)
- ❌ "Out of scope for this effort"
- ❌ "Standard pattern across codebase"
- ❌ "Will be addressed later" (no proof)

**REQUIRED RESPONSE (R332 Compliance):**
```markdown
## R355 Violation Found: [Type]

**File**: [path/to/file.go:line]
**Violation**: [TODO/Stub/Mock/Hardcoded/etc.]
**Severity**: [CRITICAL/HIGH/MEDIUM/LOW]

**R332 Actions Taken**:
1. ✅ Filed BUG-YYYYMMDD-HHMMSS-[description].md
2. ✅ Checked fix cascade status
3. ✅ Created fix cascade entry (if needed)
4. ✅ Assigned to [effort or ERROR_RECOVERY]
5. 🔴 BLOCKING review until addressed

**Bug File**: bugs/filed/BUG-YYYYMMDD-HHMMSS-[description].md
**Fix Cascade**: [CASCADE_ID]
**Status**: [FILED/IN_PROGRESS/BLOCKED]

**Review Status**: BLOCKED (exit 332)
```

### R355 Violation Categories → R332 Bug Filing

| R355 Violation | Bug Severity | R332 Action |
|----------------|--------------|-------------|
| Hardcoded credentials | CRITICAL | File bug, block review, ERROR_RECOVERY |
| Stub/Mock in production | CRITICAL | File bug, block review, immediate fix |
| TODO/FIXME comments | HIGH | File bug, create fix cascade |
| Static values | MEDIUM | File bug, track for next effort |
| Placeholder code | HIGH | File bug, block review |
| Not implemented | CRITICAL | File bug, block review, ERROR_RECOVERY |

### Code Reviewer R355 + R332 Protocol

```bash
# 1. Run R355 production readiness scan
cd $EFFORT_DIR

# 2. For EACH violation found, file R332 bug
for violation in $(detect_r355_violations); do
    BUG_ID="BUG-$(date +%Y%m%d-%H%M%S)-${violation}"

    # File bug per R332
    cat > bugs/filed/${BUG_ID}.md << EOF
# Bug: R355 Violation - ${violation}
**Severity**: [CRITICAL/HIGH/MEDIUM]
**File**: ${violation_file}
**Violation Type**: R355 Production Readiness
...
EOF

    # Check/create fix cascade
    if ! jq -e ".fix_cascades[] | select(.bugs[] | contains(\"${BUG_ID}\"))" orchestrator-state-v3.json; then
        # Create cascade per R332
        create_fix_cascade "${BUG_ID}"
    fi

    # Block if critical
    if [[ "${severity}" == "CRITICAL" ]]; then
        echo "🔴 BLOCKING review - R355 CRITICAL violation"
        exit 332  # R332 violation (untracked bug)
    fi
done

# 3. Only pass if ALL violations filed and tracked
verify_all_bugs_in_cascades || exit 332
```

### Why R332 Eliminates R355 Deferrals

**Before R332:**
```
Phase 1: "TODO acceptable for Phase 2" ❌
Phase 2: "TODO pre-existing" ❌
Result: TODO never fixed ❌
```

**After R332:**
```
Phase 1: File BUG-001, create fix cascade ✅
Phase 2: Check BUG-001 status, verify fix ✅
Result: Bug MUST be fixed or review blocked ✅
```

**R332 ensures R355 violations cannot hide as "pre-existing" or "deferred to future phase".**

## ROOT CAUSE ANALYSIS

Engineers create non-production code because:
1. **Unclear expectations** - Now crystal clear with R355
2. **Time pressure** - Not an excuse for security violations
3. **Incremental development** - Use feature flags, not stubs
4. **Size constraints** - Split properly, don't leave stubs
5. **"TODO later" mindset** - NO! Production ready NOW

## RECOVERY WHEN VIOLATIONS FOUND

1. **IMMEDIATE STOP** - No further work until fixed
2. **COMPLETE REWRITE** - No patching over violations
3. **SECURITY REVIEW** - For any credential violations
4. **FULL RESCAN** - Verify all violations removed
5. **DOCUMENT** - Explain how violations occurred

## ESCALATION PROTOCOL

When R355 violations detected:
1. **SW Engineer**: EXIT immediately with code 355
2. **Code Reviewer**: FAIL review with R355 citation
3. **Orchestrator**: BLOCK all merges, initiate recovery
4. **Architect**: Assess systemic issues if repeated

## SUMMARY

**PRODUCTION READY CODE IS NON-NEGOTIABLE**

The Software Factory produces PRODUCTION SOFTWARE that:
- **WORKS** in production environments
- **SECURES** sensitive data properly
- **CONFIGURES** from environment
- **IMPLEMENTS** all functionality
- **DELIVERS** real value

**ANY VIOLATION = PROJECT FAILURE**

This is SUPREME LAW - no exceptions, no excuses, no negotiations.