# 🚨🚨🚨 RULE R291: Integration Demo Requirement

## Classification
- **Category**: Integration Management
- **Criticality Level**: 🚨🚨🚨 BLOCKING
- **Enforcement**: MANDATORY for all integrations
- **Penalty**: -50% to -75% for violations

## The Rule

**EVERY integration at EVERY level (Wave, Phase, Project) MUST produce a working build, automated test harness, and demonstrable functionality before marking integration as complete.**

## 🔴🔴🔴 SUPREME GATE: BUILD/TEST/DEMO MUST PASS OR ERROR_RECOVERY 🔴🔴🔴

**ABSOLUTE REQUIREMENT - NO EXCEPTIONS:**
Every integration at EVERY level (effort/wave/phase/project) MUST:
1. ✅ **BUILD SUCCESSFULLY** - Compilation completes without errors
2. ✅ **PASS ALL TESTS** - Unit, integration, and E2E tests pass
3. ✅ **PRODUCE WORKING OUTPUT** - Binary/package/dist created
4. ✅ **DEMO FUNCTIONALITY** - Features actually work

**🚨🚨🚨 MANDATORY ERROR_RECOVERY TRIGGERS 🚨🚨🚨**

**If ANY of these fail, you MUST:**
- 🔴 **IMMEDIATELY STOP** - No proceeding whatsoever
- 🔴 **TRANSITION TO ERROR_RECOVERY** - This is MANDATORY, not optional
- 🔴 **DOCUMENT FAILURE** - Record exact error in state file
- 🔴 **INITIATE FIX PROTOCOL** - Follow R300 for fixes in effort branches

**Specific Failure → ERROR_RECOVERY Mappings:**
- Build failure (make/npm build/cargo build returns non-zero) → **ERROR_RECOVERY**
- Test failure (ANY test fails) → **ERROR_RECOVERY**
- No output produced (missing dist/build/target) → **ERROR_RECOVERY**
- Demo script fails (exit code != 0) → **ERROR_RECOVERY**
- Feature doesn't work as expected → **ERROR_RECOVERY**

**VIOLATION = -100% AUTOMATIC FAILURE**

Marking integration complete without passing build/test/demo = **IMMEDIATE DISQUALIFICATION**

**THIS IS AN ABSOLUTE GATE!** The code must build, test, and run successfully or you CANNOT proceed!

## Requirements

### 1. 🏗️ MANDATORY BUILD VERIFICATION WITH ERROR_RECOVERY TRIGGER

**NO INTEGRATION IS COMPLETE WITHOUT A WORKING BUILD:**

```bash
# MANDATORY VERIFICATION - MUST RUN FOR EVERY INTEGRATION
verify_integration_gates() {
    echo "🔴🔴🔴 R291 MANDATORY GATES CHECK 🔴🔴🔴"
    local FAILED=false
    local FAILURE_REASON=""
    
    # 1. BUILD GATE - MUST PASS OR ERROR_RECOVERY
    echo "🏗️ [GATE 1] Build Verification..."
    rm -rf dist/ build/ out/ target/
    
    if make build 2>&1 | tee build.log || \
       npm run build 2>&1 | tee build.log || \
       cargo build 2>&1 | tee build.log; then
        echo "✅ BUILD GATE: PASSED"
    else
        echo "🔴 BUILD GATE: FAILED - MUST ENTER ERROR_RECOVERY"
        FAILED=true
        FAILURE_REASON="Build compilation failed"
    fi
    
    # 2. ARTIFACT GATE - MUST EXIST OR ERROR_RECOVERY
    echo "📦 [GATE 2] Artifact Verification..."
    if [ -d "dist" ] || [ -d "build" ] || [ -d "out" ] || [ -d "target" ]; then
        echo "✅ ARTIFACT GATE: PASSED"
        ls -la dist/ build/ out/ target/ 2>/dev/null | tee artifacts.log
    else
        echo "🔴 ARTIFACT GATE: FAILED - MUST ENTER ERROR_RECOVERY"
        FAILED=true
        FAILURE_REASON="No build artifacts produced"
    fi
    
    # 3. TEST GATE - ALL MUST PASS OR ERROR_RECOVERY
    echo "🧪 [GATE 3] Test Verification..."
    if make test 2>&1 | tee test.log || \
       npm test 2>&1 | tee test.log || \
       cargo test 2>&1 | tee test.log || \
       pytest 2>&1 | tee test.log; then
        echo "✅ TEST GATE: PASSED"
    else
        echo "🔴 TEST GATE: FAILED - MUST ENTER ERROR_RECOVERY"
        FAILED=true
        FAILURE_REASON="Tests failed"
    fi
    
    # 4. DEMO GATE - MUST WORK OR ERROR_RECOVERY
    echo "🎬 [GATE 4] Demo Verification..."
    if [ -f "./demo-features.sh" ] && ./demo-features.sh; then
        echo "✅ DEMO GATE: PASSED"
    else
        echo "🔴 DEMO GATE: FAILED - MUST ENTER ERROR_RECOVERY"
        FAILED=true
        FAILURE_REASON="Demo script failed or missing"
    fi
    
    # FINAL VERDICT - ERROR_RECOVERY IF ANY GATE FAILED
    if [ "$FAILED" = true ]; then
        echo "🔴🔴🔴 INTEGRATION GATES FAILED 🔴🔴🔴"
        echo "FAILURE REASON: $FAILURE_REASON"
        echo "MANDATORY ACTION: Transition to ERROR_RECOVERY state"
        
        # Update state file to ERROR_RECOVERY
        yq eval ".current_state = \"ERROR_RECOVERY\"" -i orchestrator-state.json
        yq eval ".error_recovery.trigger = \"R291_BUILD_TEST_GATE_FAILURE\"" -i orchestrator-state.json
        yq eval ".error_recovery.reason = \"$FAILURE_REASON\"" -i orchestrator-state.json
        yq eval ".error_recovery.timestamp = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" -i orchestrator-state.json
        
        git add orchestrator-state.json
        git commit -m "error: R291 gate failure - entering ERROR_RECOVERY: $FAILURE_REASON"
        git push
        
        exit 1  # STOP IMMEDIATELY
    fi
    
    echo "✅✅✅ ALL GATES PASSED - Integration may proceed ✅✅✅"
    return 0
}
```

**Build Requirements:**
- Must compile/build successfully
- Must produce verifiable artifacts
- Must capture build logs
- Must be runnable/executable
- Build failures = integration incomplete

### 2. 🧪 MANDATORY TEST HARNESS

**EVERY INTEGRATION MUST HAVE AN AUTOMATED TEST HARNESS:**

```bash
# Template for test-harness.sh
cat > test-harness.sh << 'EOF'
#!/bin/bash
# Integration Test Harness
echo "🧪 Starting Integration Test Suite"
echo "=================================="

FAILED=0

# Unit tests (REQUIRED)
echo "📦 Running unit tests..."
if npm test 2>&1 | tee unit-tests.log; then
    echo "✅ Unit tests passed"
else
    echo "❌ Unit tests failed"
    ((FAILED++))
fi

# Integration tests (REQUIRED)
echo "🔗 Running integration tests..."
if npm run test:integration 2>&1 | tee integration-tests.log; then
    echo "✅ Integration tests passed"
else
    echo "❌ Integration tests failed"
    ((FAILED++))
fi

# Feature verification (REQUIRED)
echo "🎯 Verifying new features..."
if ./verify-features.sh; then
    echo "✅ Features verified"
else
    echo "❌ Feature verification failed"
    ((FAILED++))
fi

echo "=================================="
if [ $FAILED -eq 0 ]; then
    echo "✅ ALL TESTS PASSED!"
    exit 0
else
    echo "❌ $FAILED test suites failed!"
    exit 1
fi
EOF

chmod +x test-harness.sh
```

**Test Harness Requirements:**
- Must be automated and repeatable
- Must test integrated functionality
- Must clearly show pass/fail status
- Must capture test logs
- Must verify new features work

### 3. 🎬 MANDATORY DEMO

**EVERY INTEGRATION MUST DEMONSTRATE WORKING FUNCTIONALITY:**

```bash
# Demo documentation template with timestamp
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
DEMO_FILE="INTEGRATION-DEMO-${TIMESTAMP}.md"

cat > "$DEMO_FILE" << 'EOF'
# Integration Demo

## Build Status
- Build: ✅ PASSING
- Tests: ✅ ALL PASSING
- Integration: ✅ COMPLETE
- Created: [timestamp]

## Features Demonstrated
1. [Feature 1]: Working implementation with evidence
2. [Feature 2]: Integration verified with test results
3. [Feature 3]: Functionality demonstrated

## How to Run Demo
```bash
# Start application
npm start

# Run demo script
./demo-features.sh

# Verify outputs
curl http://localhost:3000/api/new-feature
```

## Evidence
- Build log: build.log
- Test results: test-results.log
- Screenshots: demos/integration/
- Demo script: demo-features.sh
EOF
```

**Demo Requirements:**
- Must create demo documentation
- Must show actual functionality working
- Must provide reproduction steps
- Must capture evidence (logs/screenshots)
- Must prove integration delivers value

### 4. Level-Specific Requirements

#### Wave Integration Demo
```bash
# Wave level requirements
- Demonstrate wave-specific features
- Show integration of all efforts in wave
- Verify no regression in previous features
- Create WAVE-DEMO.md
- Create wave-test-harness.sh
```

#### Phase Integration Demo
```bash
# Phase level requirements
- Demonstrate all waves integrated
- Show phase-level functionality complete
- Run comprehensive test suite
- Create PHASE-${N}-DEMO.md
- Create phase${N}-test-harness.sh
- Verify against phase plan deliverables
```

#### Project Integration Demo
```bash
# Project level requirements
- Demonstrate complete project functionality
- Show all phases working together
- Run full E2E test suite
- Create PROJECT-DEMO.md
- Create project-test-harness.sh
- Include performance metrics
- Include security scan results
```

## 🔧 FAILED DEMO FIX PROTOCOL

**When integration demo fails, follow this EXACT process:**

### 1. Code Reviewer Creates Fix Instructions
When demo fails (build errors, test failures, runtime issues):
```bash
# Code Reviewer analyzes failures
analyze_demo_failure() {
    echo "🔍 Analyzing demo failure..."
    
    # Check build logs
    grep -i "error\|fail" build.log > build-errors.txt
    
    # Check test logs
    grep -i "fail\|error" test-results.log > test-failures.txt
    
    # Identify affected efforts
    echo "Affected efforts:" > affected-efforts.txt
    # Trace errors back to source efforts
}

# Create timestamped fix instructions
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
FIX_FILE="FIX-INSTRUCTIONS-${TIMESTAMP}.md"

cat > "$FIX_FILE" << 'EOF'
# INTEGRATION DEMO FIX INSTRUCTIONS

## Demo Failure Summary
- Build Status: ❌ FAILED
- Test Status: ❌ FAILED
- Error Type: [compilation/runtime/test]
- Created: [timestamp]

## Root Cause Analysis
1. [Error 1]: Located in effort-X, file Y, line Z
2. [Error 2]: Located in effort-Y, file A, line B

## Required Fixes

### Effort: [effort-name-1]
**Branch**: feature/effort-name-1
**Issues**:
- [ ] Fix compilation error in src/module.ts line 45
- [ ] Update test expectations in tests/module.test.ts
- [ ] Add missing dependency to package.json

### Effort: [effort-name-2]
**Branch**: feature/effort-name-2
**Issues**:
- [ ] Fix API integration error
- [ ] Update configuration for new endpoint

## Verification Steps
1. Fix issues in effort branches
2. Run local build and tests
3. Push fixes to effort branches
4. Re-attempt integration
EOF
```

### 2. Orchestrator Receives Instructions
```bash
# Orchestrator reads fix instructions
read_fix_instructions() {
    # Find the latest fix instructions file
    LATEST_FIX=$(ls -t FIX-INSTRUCTIONS-*.md 2>/dev/null | head -n1)
    
    # Fallback to old format if needed
    if [ -z "$LATEST_FIX" ] && [ -f "FIX-INSTRUCTIONS.md" ]; then
        LATEST_FIX="FIX-INSTRUCTIONS.md"
        echo "⚠️ Using legacy fix instructions format"
    fi
    
    if [ -f "$LATEST_FIX" ]; then
        echo "📋 Processing fix instructions: $LATEST_FIX"
        
        # Extract affected efforts
        grep "^### Effort:" "$LATEST_FIX" | cut -d: -f2
        
        # Spawn SW Engineers for each effort
        for effort in $(get_affected_efforts); do
            spawn_sw_engineer_for_fixes "$effort"
        done
    fi
}
```

### 3. SW Engineers Fix in EFFORT BRANCHES

**⚠️⚠️⚠️ CRITICAL: ALL FIXES MUST BE IN EFFORT BRANCHES ⚠️⚠️⚠️**

```bash
# SW Engineer fixes in EFFORT branch
fix_in_effort_branch() {
    local effort_name="$1"
    local effort_branch="feature/${effort_name}"
    
    # ✅ CORRECT: Fix in effort branch
    git checkout "$effort_branch"
    
    # ❌ WRONG: Never fix in integration branch!
    # git checkout integration-wave-1  # NEVER DO THIS!
    
    # Apply fixes
    implement_fixes_from_instructions
    
    # Test locally
    npm test
    npm run build
    
    # Commit and push to effort branch
    git add -A
    git commit -m "fix: resolve integration demo failures for $effort_name"
    git push origin "$effort_branch"
}
```

**Why effort branch fixes are MANDATORY:**
- ✅ Ensures source branches are correct
- ✅ Makes future integrations work
- ✅ Prevents drift between effort and integration branches
- ✅ Maintains clean git history
- ✅ Allows proper PR reviews
- ✅ Enables rollback if needed

**NEVER fix directly in integration branch because:**
- ❌ Creates divergence from source branches
- ❌ Makes future merges conflict
- ❌ Hides problems in effort code
- ❌ Breaks traceability
- ❌ Violates CD principles

### 4. Re-attempt Integration

```bash
# After fixes are pushed to effort branches
retry_integration() {
    echo "🔄 Re-attempting integration with fixed code..."
    
    # Create fresh integration branch
    git checkout main
    git pull origin main
    git checkout -b integration-wave-X-retry-$(date +%s)
    
    # Merge fixed effort branches
    for effort_branch in $(get_effort_branches); do
        echo "Merging fixed $effort_branch..."
        git merge "origin/$effort_branch" --no-ff
    done
    
    # Build and test again
    npm install
    npm run build | tee build.log
    ./test-harness.sh | tee test-results.log
    
    # Demo must now pass!
    if ./demo-features.sh; then
        echo "✅ Demo passes! Integration can proceed"
    else
        echo "❌ Demo still failing - repeat fix protocol"
    fi
}
```

### 5. Fix Protocol State Machine

```mermaid
graph TD
    A[Demo Fails] --> B[Code Reviewer Analyzes]
    B --> C[Create FIX-INSTRUCTIONS.md]
    C --> D[Orchestrator Reads Instructions]
    D --> E[Spawn SW Engineers]
    E --> F[Fix in EFFORT Branches]
    F --> G[Push Fixes to Effort Branches]
    G --> H[Re-attempt Integration]
    H --> I{Demo Passes?}
    I -->|Yes| J[Integration Complete]
    I -->|No| B
```

## Implementation Process

### Step 1: Build Verification
```bash
# Clean and build
make clean && make build
# or
npm run build:clean && npm run build:prod
```

### Step 2: Test Harness Creation
```bash
# Create and run test harness
./create-test-harness.sh [wave|phase|project]
./test-harness.sh
```

### Step 3: Demo Creation
```bash
# Create demo artifacts
./create-demo.sh [wave|phase|project]
# Run demo
./demo-features.sh
```

### Step 4: Verification
```bash
# Verify all requirements met
verify_integration_complete() {
    [ -f "build.log" ] || { echo "❌ Missing build log"; return 1; }
    [ -f "test-harness.sh" ] || { echo "❌ Missing test harness"; return 1; }
    [ -f "*-DEMO.md" ] || { echo "❌ Missing demo documentation"; return 1; }
    [ -d "dist" ] || [ -d "build" ] || { echo "❌ Missing build artifacts"; return 1; }
    
    echo "✅ All integration requirements met"
    return 0
}
```

## Failure Conditions

### Critical Failures (Immediate Stop)
- 🚨 No build artifacts = FAIL
- 🚨 Build doesn't compile = FAIL
- 🚨 No test harness = FAIL
- 🚨 Tests not passing = FAIL
- 🚨 No demo created = FAIL

### Grading Penalties
- Missing build verification: **-25%**
- Missing test harness: **-25%**
- Missing demo: **-25%**
- Tests failing but ignored: **-50%**
- Build broken but claimed complete: **-75%**

## Success Criteria

Before marking ANY integration complete:
- ✅ Build compiles and runs successfully
- ✅ Build artifacts verified to exist
- ✅ Test harness created and executed
- ✅ All tests passing (unit + integration)
- ✅ Demo documentation created
- ✅ Demo script functional
- ✅ Features verified working
- ✅ Evidence captured (logs, screenshots)

## Examples

### ✅ CORRECT: Complete integration
```bash
# 1. Build verification
npm run build:prod | tee build.log
ls -la dist/

# 2. Test harness
./create-test-harness.sh wave
./test-harness.sh

# 3. Demo creation
./create-demo.sh wave
./demo-wave-features.sh

# 4. Verification
./verify-integration-complete.sh
```

### ❌ WRONG: Incomplete integration
```bash
# Just merging branches without verification
git merge feature-branch
git push
echo "Integration complete"  # NO BUILD, NO TESTS, NO DEMO!
```

## Related Rules
- R034: Integration Requirements
- R282: Phase Integration Protocol
- R283: Project Integration Protocol
- R265: Integration Testing Requirements
- R263: Integration Documentation Requirements

## Enforcement

This rule is enforced at:
1. **Wave Integration** - Every wave must demo
2. **Phase Integration** - Every phase must demo
3. **Project Integration** - Final project must demo
4. **PR Reviews** - No merge without demo evidence
5. **State Transitions** - Cannot proceed without demo

## Remember

**"If it doesn't build, it doesn't work"**
**"If it doesn't test, it's not verified"**
**"If it doesn't demo, it's not complete"**

Every integration MUST prove the code actually works through building, testing, and demonstration. No exceptions!