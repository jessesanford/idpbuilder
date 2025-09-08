# Orchestrator - PRODUCTION_READY_VALIDATION State Rules

## 🛑🛑🛑 R322 MANDATORY STOP BEFORE STATE TRANSITIONS 🛑🛑🛑

**SUPREME LAW - VIOLATION = -100% IMMEDIATE FAILURE**

### YOU MUST STOP AFTER:
1. ✅ Completing all TODOs for this state
2. ✅ Updating orchestrator-state.json with new state
3. ✅ Committing and pushing the state file  
4. ✅ Providing work summary

### YOU MUST NOT:
- ❌ Continue to the next state automatically
- ❌ Start work for the new state
- ❌ Spawn agents for the new state
- ❌ Assume permission to continue

### STOP PROTOCOL:
```markdown
## 🛑 STATE TRANSITION CHECKPOINT: CURRENT_STATE → NEXT_STATE

### ✅ Current State Work Completed:
- [List completed work]

### 📊 Current Status:
- Current State: CURRENT_STATE
- Next State: NEXT_STATE
- TODOs Completed: X/Y
- State Files: Updated and committed ✅

### ⏸️ STOPPED - Awaiting User Continuation
Ready to transition to NEXT_STATE. Please use /continue-orchestrating.
```

**STOP MEANS STOP - Exit and wait for /continue-orchestrating**

---


## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED PRODUCTION_READY_VALIDATION STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_PRODUCTION_READY_VALIDATION
echo "$(date +%s) - Rules read and acknowledged for PRODUCTION_READY_VALIDATION" > .state_rules_read_orchestrator_PRODUCTION_READY_VALIDATION
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY VALIDATION WORK UNTIL RULES ARE READ:
- ❌ Start running tests
- ❌ Start checking dependencies
- ❌ Update state files
- ❌ Continue to next state
- ❌ Think about what to do in this state

### ✅ YOU MUST IMMEDIATELY:

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

## 📋 PRIMARY DIRECTIVES

### 🚨🚨🚨 RULE R006 - Orchestrator NEVER Writes Code [BLOCKING]
**Source:** `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`
**Criticality:** BLOCKING - Any code operation = -100% IMMEDIATE FAILURE

READ THIS RULE FILE: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`

### 🚨🚨🚨 RULE R319 - Orchestrator NEVER Measures or Assesses Code [BLOCKING]
**Source:** `$CLAUDE_PROJECT_DIR/rule-library/R319-orchestrator-never-measures-code.md`
**Criticality:** BLOCKING - Any technical assessment = -100% IMMEDIATE FAILURE

READ THIS RULE FILE: `$CLAUDE_PROJECT_DIR/rule-library/R319-orchestrator-never-measures-code.md`

### 🚨🚨🚨 RULE R323 - Mandatory Final Artifact Build [BLOCKING]
**Source:** `$CLAUDE_PROJECT_DIR/rule-library/R323-mandatory-final-artifact-build.md`
**Criticality:** BLOCKING - No artifact = -50% to -100% FAILURE

READ THIS RULE FILE: `$CLAUDE_PROJECT_DIR/rule-library/R323-mandatory-final-artifact-build.md`

**⚠️ R006, R319 & R323 WARNING FOR PRODUCTION_READY_VALIDATION STATE:**
- DO NOT run tests yourself - spawn Code Reviewer!
- DO NOT execute builds yourself - spawn Code Reviewer!
- DO NOT validate dependencies yourself - spawn Code Reviewer!
- DO NOT fix test failures yourself - spawn SW Engineers!
- DO NOT modify code to make tests pass - spawn SW Engineers!
- **🔴 R323: Code Reviewer MUST verify final artifact still exists**
- **🔴 R323: Code Reviewer MUST test the actual built artifact**
- You ONLY coordinate validation - NEVER execute it yourself

### Production Ready Validation Requirements

**This state validates that the integrated code is ready for production deployment.**

Key validation areas:
- All tests must pass
- Dependencies must be secure and up-to-date
- Build must complete successfully
- **🔴 R323: Final artifact must exist and be tested**
- **🔴 R323: Artifact path must be documented**
- Documentation must be present

## 🎯 STATE OBJECTIVES

In the PRODUCTION_READY_VALIDATION state, you are the COORDINATOR:

1. **Spawning Code Reviewer for Test Validation**
   - Spawn Code Reviewer to run all test suites
   - Provide clear instructions on what to validate
   - Specify test types to execute
   - Request comprehensive test report

2. **Spawning Code Reviewer for Dependency Validation**
   - Spawn Code Reviewer to audit dependencies
   - Request security vulnerability checks
   - Ask for version compatibility verification
   - Ensure dependency resolution validation

3. **Monitoring Validation Progress**
   - Check if Code Reviewer has completed
   - Read validation reports
   - Track any issues identified
   - Coordinate fixes if needed

4. **Managing State Transitions**
   - Determine next state based on validation
   - Update state file with validation status
   - Track production readiness
   - Coordinate any required fixes

## 📝 REQUIRED ACTIONS

### Step 1: Spawn Code Reviewer for Test Validation
```bash
# ✅ CORRECT: Delegate test validation to Code Reviewer
echo "🧪 Production readiness validation needed"
echo "🚀 Spawning Code Reviewer to run test suites..."

cd /efforts/integration-testing

# Update state to show spawning Code Reviewer
yq -i '.current_state = "SPAWN_CODE_REVIEWER_PROD_VALIDATION"' orchestrator-state.json
yq -i '.spawn_in_progress.agent = "code-reviewer"' orchestrator-state.json
yq -i '.spawn_in_progress.purpose = "production_validation"' orchestrator-state.json

# Spawn Code Reviewer with validation task
Task: subagent_type="code-reviewer" \
      state="VALIDATION" \
      prompt="Perform production readiness validation. Run all test suites (unit, integration, e2e), validate dependencies, check security, verify deployment readiness. Create PRODUCTION-VALIDATION-REPORT.md with comprehensive findings." \
      workspace="/efforts/integration-testing" \
      description="Production readiness validation"

echo "⏳ Code Reviewer spawned for production validation"
echo "📋 Waiting for PRODUCTION-VALIDATION-REPORT.md"

# ❌❌❌ FORBIDDEN - Orchestrator CANNOT run tests!
# npm test  # VIOLATION OF R006/R319 = IMMEDIATE FAILURE!
    npm run test:coverage || echo "No coverage script"
elif [ -f "go.mod" ]; then
    # Go project
    go test ./... -v
    go test ./... -cover
elif [ -f "Cargo.toml" ]; then
    # Rust project
    cargo test
    cargo test --all-features
elif [ -f "pom.xml" ]; then
    # Java/Maven project
    mvn test
    mvn verify
elif [ -f "build.gradle" ] || [ -f "build.gradle.kts" ]; then
    # Java/Gradle project
    ./gradlew test
    ./gradlew check
elif [ -f "requirements.txt" ] || [ -f "setup.py" ]; then
    # Python project
    pip install -r requirements.txt || pip install -e .
    pytest || python -m pytest
    pytest --cov || echo "No coverage configured"
else
    echo "⚠️ No standard test framework detected"
    echo "Looking for test files..."
    find . -name "*test*" -type f | head -20
fi

# Save test results
echo "Test execution completed at $(date)" > TEST-RESULTS.txt
```

### Step 2: Validate Dependencies
```bash
# Check for known vulnerability scanners
if command -v npm &> /dev/null; then
    npm audit || true
    npm list --depth=0 > DEPENDENCIES.txt
elif command -v go &> /dev/null; then
    go list -m all > DEPENDENCIES.txt
    # Check for vulnerable dependencies
    go list -json -m all | nancy sleuth || echo "Nancy not installed"
elif command -v cargo &> /dev/null; then
    cargo tree > DEPENDENCIES.txt
    cargo audit || echo "cargo-audit not installed"
elif command -v mvn &> /dev/null; then
    mvn dependency:tree > DEPENDENCIES.txt
    mvn dependency-check:check || echo "OWASP check not configured"
elif command -v pip &> /dev/null; then
    pip list > DEPENDENCIES.txt
    safety check || pip install safety && safety check || echo "Safety not available"
fi
```

### Step 3: Check Production Requirements
```bash
# Check for essential files
MISSING_FILES=""

# Configuration files
for config in "config.yaml" "config.json" ".env.example" "settings.ini"; do
    if [ ! -f "$config" ]; then
        MISSING_FILES="$MISSING_FILES\n- $config"
    fi
done

# Documentation
if [ ! -f "README.md" ]; then
    MISSING_FILES="$MISSING_FILES\n- README.md"
fi

# Deployment files
for deploy in "Dockerfile" "docker-compose.yml" "kubernetes.yaml" ".github/workflows/deploy.yml"; do
    if [ -f "$deploy" ]; then
        echo "✅ Found deployment file: $deploy"
        break
    fi
done

if [ -n "$MISSING_FILES" ]; then
    echo "⚠️ Missing files:$MISSING_FILES"
fi
```

### Step 4: Create Validation Report
```bash
cat > PRODUCTION-READY-VALIDATION-REPORT.md << 'EOF'
# Production Ready Validation Report
Date: $(date)
State: PRODUCTION_READY_VALIDATION

## Test Results
### Unit Tests
- Passed: [X/Y]
- Failed: [List any failures]
- Coverage: [XX%]

### Integration Tests
- Passed: [X/Y]
- Failed: [List any failures]

### End-to-End Tests
- Status: [PASS/FAIL/NOT_APPLICABLE]

## Dependency Validation
### Security Scan
- Vulnerabilities Found: [Count]
- Critical: [Count]
- High: [Count]
- Medium: [Count]
- Low: [Count]

### Dependency List
- Total Dependencies: [Count]
- Direct: [Count]
- Transitive: [Count]

## Production Readiness Checklist
- [ ] All tests passing
- [ ] No critical vulnerabilities
- [ ] Configuration files present
- [ ] Documentation complete
- [ ] Deployment scripts validated
- [ ] Environment variables documented
- [ ] Logging configured
- [ ] Error handling implemented
- [ ] Performance acceptable

## Issues Found
[List any issues that need fixing]

## Recommendation
[PROCEED_TO_BUILD / FIX_REQUIRED]

## Next Steps
[Transition to BUILD_VALIDATION or FIX_BUILD_ISSUES]
EOF
```

## ⚠️ CRITICAL REQUIREMENTS

### 🔴 R006 VIOLATION DETECTION 
If you find yourself about to:
- Fix failing tests by editing test code
- Modify source code to make tests pass
- Update dependencies directly in code
- Apply patches to fix issues
- Make excuses like "just a quick fix for the test"

**STOP IMMEDIATELY - You are violating R006!**
Document failures and spawn SW Engineers to fix them!

### Test Coverage Thresholds
- Unit test coverage: Minimum 60% (configurable)
- Critical paths: Must have integration tests
- API endpoints: Must have tests
- Database operations: Must have tests

### Security Requirements
- NO critical vulnerabilities allowed
- High vulnerabilities require justification
- All dependencies must be current stable versions
- No deprecated dependencies

### Documentation Requirements
- README.md with setup instructions
- API documentation if applicable
- Configuration documentation
- Deployment instructions

## 🚫 FORBIDDEN ACTIONS

1. **NEVER edit any code files yourself** - R006 VIOLATION = -100%
2. **NEVER fix test failures yourself** - R006 VIOLATION = -100%
3. **NEVER modify test code to make it pass** - R006 VIOLATION = -100%
4. **NEVER update dependencies in code** - R006 VIOLATION = -100%
5. **NEVER skip tests even if they're slow**
6. **NEVER ignore security vulnerabilities**
7. **NEVER proceed with failing tests**
8. **NEVER fake test results**

## ✅ SUCCESS CRITERIA

Before transitioning to BUILD_VALIDATION:
- [ ] All test suites executed
- [ ] Test results documented
- [ ] No critical failures
- [ ] Dependencies validated
- [ ] No critical security issues
- [ ] Production checklist reviewed
- [ ] Validation report created

## 🔄 STATE TRANSITIONS

### Success Path:
```
PRODUCTION_READY_VALIDATION → BUILD_VALIDATION
```
- All tests pass
- No critical issues
- Ready for build validation

### Error Path:
```
PRODUCTION_READY_VALIDATION → FIX_BUILD_ISSUES
```
- Test failures found
- Security vulnerabilities detected
- Missing critical components

## 📊 VERIFICATION CHECKLIST

Before leaving this state:
```bash
# Verify test results captured
ls -la TEST-RESULTS.txt

# Verify dependency audit completed
ls -la DEPENDENCIES.txt

# Verify validation report exists
ls -la PRODUCTION-READY-VALIDATION-REPORT.md

# Check for any uncommitted changes
git status

# Commit validation results
git add -A
git commit -m "validation: production readiness validation complete"
git push
```

## 🎓 GRADING CRITERIA

You will be evaluated on:
1. **Test Execution** (30%)
   - All test suites run
   - Results properly captured
   
2. **Security Validation** (30%)
   - Vulnerability scanning performed
   - Results documented
   
3. **Production Checklist** (20%)
   - All items reviewed
   - Missing items identified
   
4. **Documentation** (20%)
   - Complete validation report
   - Clear recommendations

## 💡 TIPS FOR SUCCESS

1. **Run Tests in Order**: Unit → Integration → E2E
2. **Don't Ignore Warnings**: They become errors in production
3. **Check Logs**: Test output contains valuable information
4. **Be Thorough**: Better to find issues now than in production

## 🚨 COMMON PITFALLS TO AVOID

1. **Assuming Tests Pass**: Always run them
2. **Ignoring Flaky Tests**: They indicate real problems
3. **Skipping Security Scans**: Critical for production
4. **Incomplete Reports**: Future states need details

## 🔧 TROUBLESHOOTING GUIDE

### If Tests Won't Run:
1. Check test framework installation
2. Verify working directory
3. Check for test configuration files
4. Look for test scripts in package.json/Makefile

### If Dependencies Can't Be Validated:
1. Install security scanning tools
2. Use online vulnerability databases
3. Check package manager documentation
4. Document manual review if automated tools unavailable

Remember: This state ensures the software is truly production-ready!
## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-orchestrating
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**
