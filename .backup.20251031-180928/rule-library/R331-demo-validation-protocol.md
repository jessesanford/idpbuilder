# 🚨🚨🚨 RULE R331: Demo Validation Protocol

## Classification
- **Category**: Quality Assurance
- **Criticality Level**: 🚨🚨🚨 BLOCKING
- **Enforcement**: MANDATORY for all demos
- **Penalty**: -100% for simulation violations, -50% to -75% for other violations
- **Related Rules**: R291 (Integration Demo Requirement), R330 (Demo Planning Requirements)

## The Rule

**ALL demos (wave, phase, project integrations) MUST be REAL working implementations that execute actual code and verify external side effects. NO SIMULATION ALLOWED.**

Demos that use non-existent flags, skip actual execution, or cannot fail when implementation is broken are **FALSE POSITIVES** and result in **IMMEDIATE FAILURE (-100%)**.

## 🔴🔴🔴 ABSOLUTE PROHIBITIONS 🔴🔴🔴

### PROHIBITED: Simulation and Fake Validation

**The following are STRICTLY FORBIDDEN and result in -100% penalty:**

```bash
# ❌ FORBIDDEN: Simulated output without actual execution
echo "✅ Feature working!"  # WITHOUT actually running the feature

# ❌ FORBIDDEN: Hardcoded success without validation
DEMO_RESULT="PASSED"  # Without actually testing anything

# ❌ FORBIDDEN: Non-existent flags
my-binary --demo-mode  # Flag doesn't exist in --help output

# ❌ FORBIDDEN: Stubbed implementation with TODO comments
function demo_feature() {
    # TODO: Implement actual feature
    echo "Feature works!"  # LIES!
}

# ❌ FORBIDDEN: Demos that cannot fail
demo_critical_feature() {
    echo "✅ All tests passed"  # Always returns 0
    return 0
}

# ❌ FORBIDDEN: Mock responses without external verification
curl http://fake-endpoint  # Not checking if endpoint exists
echo '{"status": "success"}'  # Faking the response

# ❌ FORBIDDEN: Internal state only (no external verification)
MY_VAR="deployed"  # Just setting a variable
echo "✅ Deployment successful"  # No actual deployment verified
```

### REQUIRED: Real Execution and External Verification

**All demos MUST:**

```bash
# ✅ REQUIRED: Execute actual implementation
my-binary feature-command 2>&1 | tee execution.log
if [ ${PIPESTATUS[0]} -ne 0 ]; then
    echo "❌ DEMO FAILED: Feature execution returned error"
    exit 1
fi

# ✅ REQUIRED: Verify external side effects
if ! docker images | grep -q "myimage:latest"; then
    echo "❌ DEMO FAILED: Expected image not in registry"
    exit 1
fi

# ✅ REQUIRED: Use only flags that exist in --help
my-binary --help | grep -q -- "--actual-flag" || {
    echo "❌ Flag --actual-flag not supported"
    exit 1
}
my-binary --actual-flag value  # Now use the verified flag

# ✅ REQUIRED: Verify implementation has no TODOs in execution path
grep -r "TODO\|FIXME\|XXX" src/critical_feature.go && {
    echo "❌ DEMO FAILED: Implementation incomplete (TODO found)"
    exit 1
}

# ✅ REQUIRED: Test that demo can fail
# If demo always succeeds even when implementation broken, it's a FALSE POSITIVE
```

## Requirements

### 1. 🔍 Pre-Demo Implementation Completeness Scan (MANDATORY)

**BEFORE running any demo, scan for incomplete implementation:**

```bash
#!/bin/bash
# Pre-demo validation scan

echo "🔍 R331 Pre-Demo Validation: Implementation Completeness Scan"

FAILED=false

# 1. Check for TODO comments in implementation (not tests/docs)
echo "📋 Scanning for TODO/FIXME/XXX in implementation..."
TODO_FILES=$(find src/ pkg/ internal/ -type f \( -name "*.go" -o -name "*.ts" -o -name "*.py" -o -name "*.rs" \) \
    -exec grep -l "TODO\|FIXME\|XXX\|STUB\|NOT IMPLEMENTED" {} \; 2>/dev/null)

if [ -n "$TODO_FILES" ]; then
    echo "❌ CRITICAL: Incomplete implementation detected!"
    echo "Files with TODOs:"
    echo "$TODO_FILES"
    echo ""
    echo "🚨 R331 VIOLATION: Cannot demo incomplete implementation!"
    echo "PENALTY: -100% (implementation not ready)"
    FAILED=true
fi

# 2. Check for stub/mock functions in production code
echo "📋 Scanning for stub/mock functions..."
STUB_FUNCTIONS=$(grep -rn "func.*Stub\|func.*Mock\|function.*stub\|def.*mock" src/ pkg/ internal/ 2>/dev/null)

if [ -n "$STUB_FUNCTIONS" ]; then
    echo "❌ WARNING: Stub/Mock functions in production code!"
    echo "$STUB_FUNCTIONS"
    echo ""
    echo "⚠️ R331 WARNING: Stubs should be in test code, not production"
fi

# 3. Verify build succeeds
echo "📋 Verifying build..."
if ! make build 2>&1 | tee /tmp/build.log; then
    echo "❌ CRITICAL: Build failed!"
    echo "🚨 R331 VIOLATION: Cannot demo code that doesn't build!"
    echo "PENALTY: -100% (must build successfully)"
    FAILED=true
fi

if [ "$FAILED" = true ]; then
    echo ""
    echo "🔴🔴🔴 R331 PRE-DEMO VALIDATION FAILED 🔴🔴🔴"
    echo "REQUIRED ACTION: Complete implementation before creating demo"
    exit 1
fi

echo "✅ R331 Pre-Demo Validation: PASSED"
echo "Implementation is ready for demo creation"
```

**Trigger this scan:**
- Before creating demo script
- Before running demo validation
- In CI/CD before demo execution

**Penalty for skipping:** -100% (allows incomplete implementation to pass)

### 2. 🎯 Flag Validation Against --help Output (MANDATORY)

**Every flag/option used MUST exist in the actual binary:**

```bash
#!/bin/bash
# Flag validation protocol

BINARY="my-application"

# Get actual help output
HELP_OUTPUT=$($BINARY --help 2>&1)

# Validate each flag before use
validate_flag() {
    local flag="$1"

    if echo "$HELP_OUTPUT" | grep -q -- "$flag"; then
        echo "✅ Flag validated: $flag"
        return 0
    else
        echo "❌ INVALID FLAG: $flag not found in --help output"
        echo ""
        echo "Available flags:"
        echo "$HELP_OUTPUT"
        echo ""
        echo "🚨 R331 VIOLATION: Using non-existent flag!"
        echo "PENALTY: -50% (false positive - flag doesn't exist)"
        exit 1
    fi
}

# Before using any flag, validate it exists
validate_flag "--output"
validate_flag "--format"
validate_flag "--registry"

# NOW use the validated flags
$BINARY --output result.txt --format json --registry localhost:5000
```

**Flags must be validated because:**
- Prevents demos from "working" with imaginary features
- Ensures demo tests actual implementation
- Catches copy-paste errors from other projects
- Proves feature actually exists

**Penalty for fake flags:** -50% per flag (false positive detection)

### 3. 🌍 External Side Effect Verification (MANDATORY)

**Demos MUST verify EXTERNAL state changes, not just internal variables:**

```bash
#!/bin/bash
# External verification protocol

demo_container_push() {
    local IMAGE="myapp:latest"
    local REGISTRY="localhost:5000"

    # ❌ WRONG: Internal state only
    # PUSHED=true
    # echo "✅ Image pushed"

    # ✅ CORRECT: Verify external registry state
    echo "📦 Pushing image to registry..."
    docker push "${REGISTRY}/${IMAGE}" 2>&1 | tee push.log

    if [ ${PIPESTATUS[0]} -ne 0 ]; then
        echo "❌ DEMO FAILED: Push command failed"
        cat push.log
        exit 1
    fi

    echo "🔍 Verifying image in external registry..."

    # Method 1: Query registry API
    if curl -s "http://${REGISTRY}/v2/${IMAGE}/tags/list" | grep -q "latest"; then
        echo "✅ VERIFIED: Image exists in registry (API check)"
    else
        echo "❌ DEMO FAILED: Image not found in registry API"
        exit 1
    fi

    # Method 2: Pull from registry to confirm
    docker rmi "${REGISTRY}/${IMAGE}" 2>/dev/null  # Remove local copy
    if docker pull "${REGISTRY}/${IMAGE}"; then
        echo "✅ VERIFIED: Image pullable from registry"
    else
        echo "❌ DEMO FAILED: Cannot pull image from registry"
        exit 1
    fi

    echo "✅ External verification complete: Image confirmed in registry"
}

demo_file_creation() {
    local OUTPUT_FILE="/tmp/demo-output.txt"

    # ❌ WRONG: Just echo success
    # echo "✅ File created"

    # ✅ CORRECT: Verify file exists with expected content
    echo "📝 Creating output file..."
    my-binary generate --output "$OUTPUT_FILE"

    if [ ! -f "$OUTPUT_FILE" ]; then
        echo "❌ DEMO FAILED: Output file not created"
        exit 1
    fi

    if [ ! -s "$OUTPUT_FILE" ]; then
        echo "❌ DEMO FAILED: Output file is empty"
        exit 1
    fi

    echo "✅ VERIFIED: File exists and has content"
    echo "Content preview:"
    head -5 "$OUTPUT_FILE"
}

demo_database_write() {
    local DB="testdb"
    local TABLE="users"

    # ❌ WRONG: Assume write worked
    # echo "INSERT INTO users..." | sqlite3 $DB
    # echo "✅ User created"

    # ✅ CORRECT: Verify record exists in database
    echo "💾 Creating user record..."
    echo "INSERT INTO users (name, email) VALUES ('Test', 'test@example.com');" | sqlite3 "$DB"

    echo "🔍 Verifying record in database..."
    RECORD_COUNT=$(sqlite3 "$DB" "SELECT COUNT(*) FROM users WHERE email='test@example.com';")

    if [ "$RECORD_COUNT" -eq 1 ]; then
        echo "✅ VERIFIED: Record exists in database"
    else
        echo "❌ DEMO FAILED: Record not found (count: $RECORD_COUNT)"
        exit 1
    fi
}
```

**External verification examples:**
- Container registry: Query API or pull image
- Files: Check existence, size, content
- Database: Query for records
- Network: Verify endpoints respond
- Process: Check if running in `ps`
- Logs: Verify log entries exist

**Penalty for no external verification:** -75% (cannot prove feature works)

### 4. 🧪 Demo Failure Testing (MANDATORY)

**Demos MUST be capable of failing when implementation is broken:**

```bash
#!/bin/bash
# Demo failure testing protocol

echo "🧪 R331 Requirement: Verify demo can detect failures"

# Test 1: Intentionally break implementation
test_demo_can_fail() {
    echo "Test 1: Breaking implementation to verify demo fails..."

    # Save original binary
    cp bin/myapp bin/myapp.original

    # Replace with broken version
    cat > bin/myapp << 'EOF'
#!/bin/bash
echo "ERROR: Feature not implemented"
exit 1
EOF
    chmod +x bin/myapp

    # Run demo - should FAIL
    if ./demo.sh > /dev/null 2>&1; then
        echo "❌ CRITICAL: Demo passed with broken implementation!"
        echo "🚨 R331 VIOLATION: Demo is a FALSE POSITIVE"
        echo "PENALTY: -100% (demo cannot detect failures)"

        # Restore original
        mv bin/myapp.original bin/myapp
        exit 1
    else
        echo "✅ GOOD: Demo correctly failed with broken implementation"
    fi

    # Restore original
    mv bin/myapp.original bin/myapp

    # Verify demo passes with real implementation
    if ./demo.sh; then
        echo "✅ GOOD: Demo passes with working implementation"
    else
        echo "❌ Demo fails even with working implementation"
        exit 1
    fi
}

# Test 2: Verify error detection
test_error_detection() {
    echo "Test 2: Verify demo detects errors..."

    # Demo should check exit codes
    if grep -q "exit \$?\|return \$?" demo.sh; then
        echo "❌ WARNING: Demo may ignore errors (uses \$?)"
    fi

    # Demo should use set -e or check PIPESTATUS
    if ! grep -q "set -e\|PIPESTATUS" demo.sh; then
        echo "⚠️ WARNING: Demo should use 'set -e' or check PIPESTATUS"
    fi

    # Demo should have explicit error checks
    ERROR_CHECKS=$(grep -c "if.*; then" demo.sh)
    if [ "$ERROR_CHECKS" -lt 3 ]; then
        echo "⚠️ WARNING: Demo has few error checks ($ERROR_CHECKS found)"
    fi
}

test_demo_can_fail
test_error_detection

echo "✅ Demo failure testing complete"
```

**A demo that always succeeds is worthless because:**
- Cannot detect regressions
- Hides implementation problems
- Gives false confidence
- Wastes everyone's time

**Penalty for demos that cannot fail:** -100% (false positive)

### 5. 📋 Demo Validation Checklist (MANDATORY)

**Every demo MUST pass this checklist before being accepted:**

```bash
#!/bin/bash
# R331 Demo Validation Checklist

validate_demo() {
    local demo_script="$1"
    local VIOLATIONS=0

    echo "🔍 R331 Demo Validation Checklist"
    echo "=================================="
    echo "Demo: $demo_script"
    echo ""

    # Check 1: No TODO in implementation
    echo "✓ Check 1: Implementation completeness..."
    if find src/ pkg/ internal/ -type f -exec grep -l "TODO\|FIXME" {} \; 2>/dev/null | grep -q .; then
        echo "  ❌ FAILED: TODO comments in implementation"
        ((VIOLATIONS++))
    else
        echo "  ✅ PASSED: No TODOs in implementation"
    fi

    # Check 2: All flags exist in --help
    echo "✓ Check 2: Flag validation..."
    DEMO_FLAGS=$(grep -o '\-\-[a-z-]*' "$demo_script" | sort -u)
    BINARY=$(grep -m1 '^[^#]*\(./\|bin/\)[a-z-]*' "$demo_script" | awk '{print $1}')

    if [ -n "$BINARY" ] && [ -x "$BINARY" ]; then
        HELP_OUTPUT=$($BINARY --help 2>&1)
        FLAG_MISSING=false

        for flag in $DEMO_FLAGS; do
            if ! echo "$HELP_OUTPUT" | grep -q -- "$flag"; then
                echo "  ❌ FAILED: Flag $flag not in --help"
                FLAG_MISSING=true
                ((VIOLATIONS++))
            fi
        done

        if [ "$FLAG_MISSING" = false ]; then
            echo "  ✅ PASSED: All flags exist in --help"
        fi
    else
        echo "  ⚠️ SKIPPED: Binary not found or not executable"
    fi

    # Check 3: External verification present
    echo "✓ Check 3: External verification..."
    if grep -q "docker images\|docker pull\|curl.*registry\|sqlite3\|ls -la\|cat.*file" "$demo_script"; then
        echo "  ✅ PASSED: External verification detected"
    else
        echo "  ❌ FAILED: No external verification found"
        ((VIOLATIONS++))
    fi

    # Check 4: Error handling present
    echo "✓ Check 4: Error handling..."
    if grep -q "set -e\|PIPESTATUS\|exit 1" "$demo_script"; then
        echo "  ✅ PASSED: Error handling detected"
    else
        echo "  ❌ FAILED: No error handling found"
        ((VIOLATIONS++))
    fi

    # Check 5: No simulation patterns
    echo "✓ Check 5: No simulation..."
    SIM_PATTERNS="echo.*✅.*without actual\|RESULT=.*PASSED.*without check\|Mock.*Success"
    if grep -qi "$SIM_PATTERNS" "$demo_script"; then
        echo "  ❌ FAILED: Simulation patterns detected"
        ((VIOLATIONS++))
    else
        echo "  ✅ PASSED: No simulation patterns"
    fi

    # Check 6: Demo can fail (test mode)
    echo "✓ Check 6: Demo failure capability..."
    if [ -x "$demo_script" ]; then
        # This would need implementation-specific breaking
        echo "  ℹ️ MANUAL: Verify demo fails when implementation broken"
    fi

    echo ""
    echo "=================================="
    if [ $VIOLATIONS -eq 0 ]; then
        echo "✅ R331 VALIDATION: PASSED"
        echo "Demo meets all requirements"
        return 0
    else
        echo "❌ R331 VALIDATION: FAILED"
        echo "Violations: $VIOLATIONS"
        echo "PENALTY: -100% if simulation, -50% to -75% for other violations"
        return 1
    fi
}

# Usage
validate_demo "demo-features.sh"
```

## Integration with Other Rules

### With R291 (Integration Demo Requirement)

R291 requires demos exist and pass. R331 defines HOW to validate demos are real:

```bash
# R291: Demo must exist and pass
if [ ! -f "demo.sh" ]; then
    echo "❌ R291 VIOLATION: No demo created"
    exit 291
fi

# R331: Demo must be REAL (not simulated)
if ! bash /path/to/r331-validation.sh demo.sh; then
    echo "❌ R331 VIOLATION: Demo is simulated or invalid"
    echo "PENALTY: -100%"
    exit 331
fi

# Both R291 and R331 satisfied
./demo.sh
if [ $? -eq 0 ]; then
    echo "✅ Demo passed R291 and R331 requirements"
fi
```

### With R330 (Demo Planning Requirements)

R330 requires demo planning. R331 requires planned demos are implemented correctly:

```bash
# R330: Demo must be planned
if ! grep -q "Demo Requirements" IMPLEMENTATION-PLAN.md; then
    echo "❌ R330 VIOLATION: Demo not planned"
    exit 330
fi

# R331: Planned demo must be real, not simulated
validate_demo_implementation() {
    # Check demo matches plan
    PLANNED_SCENARIOS=$(grep -A 20 "Demo Scenarios" IMPLEMENTATION-PLAN.md | grep "^#### Scenario")

    # Verify each scenario has real implementation
    for scenario in $PLANNED_SCENARIOS; do
        if ! demo_implements_scenario "$scenario"; then
            echo "❌ R331 VIOLATION: Scenario not actually implemented"
            exit 331
        fi
    done
}
```

## Enforcement

### Integration Agent (Demo Creation)

```bash
# When creating demos during integration
create_integration_demo() {
    echo "📝 Creating integration demo per R330 plan..."

    # Create demo script
    cat > demo.sh << 'EOF'
#!/bin/bash
set -e  # Fail on any error (R331 requirement)

# R331 Requirement: Pre-demo validation
echo "🔍 R331: Pre-demo implementation scan..."
bash /path/to/pre-demo-scan.sh || exit 1

# R331 Requirement: Validate flags before use
echo "🔍 R331: Validating command flags..."
BINARY="./bin/myapp"
$BINARY --help | grep -q -- "--registry" || {
    echo "❌ Flag --registry not supported"
    exit 1
}

# R331 Requirement: Execute real implementation
echo "🚀 Executing actual feature..."
$BINARY push --registry localhost:5000 myimage:latest 2>&1 | tee execution.log

if [ ${PIPESTATUS[0]} -ne 0 ]; then
    echo "❌ DEMO FAILED: Command returned error"
    exit 1
fi

# R331 Requirement: Verify external state
echo "🔍 Verifying external registry..."
if ! curl -s http://localhost:5000/v2/myimage/tags/list | grep -q "latest"; then
    echo "❌ DEMO FAILED: Image not in registry"
    exit 1
fi

echo "✅ Demo passed: Feature working with external verification"
EOF

    chmod +x demo.sh

    # R331: Validate demo before committing
    bash /path/to/r331-validation.sh demo.sh || {
        echo "❌ R331 validation failed - demo not compliant"
        exit 1
    }

    echo "✅ Demo created and R331 validated"
}
```

### Code Reviewer (Demo Validation State)

```bash
# DEMO_VALIDATION state responsibilities
validate_integration_demos() {
    echo "🔍 Code Reviewer: Demo Validation (R331 Protocol)"

    DEMO_DIR="demos/phase1/wave1/integration"

    if [ ! -d "$DEMO_DIR" ]; then
        echo "❌ No demo directory found"
        exit 1
    fi

    # Find all demo scripts
    DEMO_SCRIPTS=$(find "$DEMO_DIR" -name "*.sh" -type f)

    if [ -z "$DEMO_SCRIPTS" ]; then
        echo "❌ No demo scripts found"
        exit 1
    fi

    # R331: Validate each demo
    for demo in $DEMO_SCRIPTS; do
        echo "📋 Validating: $demo"

        # Run R331 validation checklist
        if ! bash /path/to/r331-validation.sh "$demo"; then
            echo "❌ R331 VIOLATION: Demo failed validation"
            echo "Demo: $demo"
            exit 1
        fi

        # Execute demo
        echo "🚀 Executing demo..."
        if bash "$demo" 2>&1 | tee "${demo}.log"; then
            echo "✅ Demo passed: $demo"
        else
            echo "❌ Demo failed: $demo"
            echo "See log: ${demo}.log"
            exit 1
        fi
    done

    echo "✅ All demos validated and passed (R331 compliant)"

    # Create validation report
    cat > demo-evaluation-report.md << EOF
# Demo Validation Report

## Summary
- Demos Validated: $(echo "$DEMO_SCRIPTS" | wc -l)
- R331 Compliance: PASSED
- Execution Status: ALL PASSED

## R331 Validation Results
- ✅ No TODO in implementation
- ✅ All flags validated against --help
- ✅ External side effects verified
- ✅ Demos capable of failing
- ✅ No simulation detected

## Recommendation
Integration may proceed - all demos are REAL and working.
EOF
}
```

### Orchestrator (Demo Validation Enforcement)

```bash
# WAITING_FOR_DEMO_VALIDATION state
enforce_demo_validation() {
    echo "🔍 Orchestrator: Enforcing R291 + R331 requirements"

    # Wait for demo validation report
    REPORT="demo-evaluation-report.md"

    if [ ! -f "$REPORT" ]; then
        echo "❌ No demo validation report - cannot proceed"
        transition_to ERROR_RECOVERY "Missing demo validation"
        exit 1
    fi

    # Check R331 compliance
    if ! grep -q "R331 Compliance: PASSED" "$REPORT"; then
        echo "❌ R331 VIOLATION: Demos are not compliant"
        echo "Transitioning to ERROR_RECOVERY..."
        transition_to ERROR_RECOVERY "R331 demo validation failed"
        exit 1
    fi

    # Check demo execution
    if ! grep -q "Execution Status: ALL PASSED" "$REPORT"; then
        echo "❌ R291 GATE 4 FAILED: Demos did not pass"
        echo "Transitioning to ERROR_RECOVERY..."
        transition_to ERROR_RECOVERY "Demo execution failed"
        exit 1
    fi

    echo "✅ R291 + R331: All demos validated and passed"
    echo "Integration may proceed to completion"
}
```

## Failure Conditions

### Critical Failures (-100% IMMEDIATE)

- 🚨 Simulation detected (hardcoded success without execution)
- 🚨 TODO/FIXME in implementation execution path
- 🚨 Demo cannot fail (always returns success)
- 🚨 No external verification (internal state only)

### Major Violations (-75%)

- ⚠️ No external side effect verification
- ⚠️ Missing error handling
- ⚠️ Incomplete implementation with stubs

### Moderate Violations (-50%)

- ⚠️ Using non-existent flags
- ⚠️ Insufficient error detection
- ⚠️ Weak validation (echo instead of check)

## Success Criteria

Before accepting ANY demo:
- ✅ Pre-demo implementation scan passes (no TODOs)
- ✅ All flags validated against --help
- ✅ External side effects verified
- ✅ Demo can fail when implementation broken
- ✅ R331 validation checklist passes
- ✅ No simulation patterns detected
- ✅ Proper error handling implemented

## Examples

### ✅ CORRECT: Real Demo with Full Validation

```bash
#!/bin/bash
set -e  # Exit on any error

echo "🎬 Wave 1 Integration Demo - OCI Image Push"

# R331: Pre-demo validation
echo "🔍 Step 1: Pre-demo implementation scan..."
if find internal/ -type f -name "*.go" -exec grep -l "TODO\|FIXME" {} \; | grep -q .; then
    echo "❌ Implementation incomplete (TODO found)"
    exit 1
fi

# R331: Flag validation
echo "🔍 Step 2: Validating flags..."
BINARY="./bin/ocictl"
HELP_OUTPUT=$($BINARY --help)

if ! echo "$HELP_OUTPUT" | grep -q -- "--registry"; then
    echo "❌ Flag --registry not supported"
    exit 1
fi

# R331: Real execution
echo "🚀 Step 3: Building OCI image..."
docker build -t testimage:latest . 2>&1 | tee build.log

if [ ${PIPESTATUS[0]} -ne 0 ]; then
    echo "❌ Build failed"
    exit 1
fi

echo "🚀 Step 4: Pushing to registry..."
$BINARY push --registry localhost:5000 testimage:latest 2>&1 | tee push.log

if [ ${PIPESTATUS[0]} -ne 0 ]; then
    echo "❌ Push failed"
    exit 1
fi

# R331: External verification
echo "🔍 Step 5: Verifying image in registry..."

# Check via API
REGISTRY_RESPONSE=$(curl -s http://localhost:5000/v2/testimage/tags/list)
if ! echo "$REGISTRY_RESPONSE" | grep -q "latest"; then
    echo "❌ Image not found in registry API"
    echo "Response: $REGISTRY_RESPONSE"
    exit 1
fi

# Check via pull
docker rmi localhost:5000/testimage:latest 2>/dev/null || true
if ! docker pull localhost:5000/testimage:latest; then
    echo "❌ Cannot pull image from registry"
    exit 1
fi

echo "✅ Demo PASSED: All verifications successful"
echo "  - Implementation complete (no TODOs)"
echo "  - Flags validated"
echo "  - Build successful"
echo "  - Push successful"
echo "  - External registry verified"
```

### ❌ WRONG: Simulated Demo (IMMEDIATE -100%)

```bash
#!/bin/bash

echo "🎬 Wave 1 Integration Demo - OCI Image Push"

# ❌ VIOLATION: No pre-demo scan
# ❌ VIOLATION: No flag validation

echo "🚀 Pushing to registry..."
# ❌ VIOLATION: Simulated without actual execution
echo "✅ Push successful"

# ❌ VIOLATION: No external verification
echo "✅ Image in registry"

# ❌ VIOLATION: Always succeeds, cannot fail
echo "✅ Demo PASSED"
exit 0  # ALWAYS returns success - FALSE POSITIVE!
```

### ❌ WRONG: Using Non-Existent Flags (-50%)

```bash
#!/bin/bash
set -e

echo "🎬 Demo with Non-Existent Flag"

# ❌ VIOLATION: Flag not validated
./bin/myapp --demo-mode --fake-flag data.txt

# Flag doesn't exist in --help, but demo doesn't check!
# This is a FALSE POSITIVE - pretending feature exists
```

### ❌ WRONG: No External Verification (-75%)

```bash
#!/bin/bash
set -e

echo "🎬 Demo without External Verification"

# Executes command
./bin/myapp push myimage:latest

# ❌ VIOLATION: Only checks exit code, not external state
if [ $? -eq 0 ]; then
    echo "✅ Push successful"
fi

# Never verifies image actually in registry!
# Could have failed internally but returned 0
```

## Related Rules

- **R291**: Integration Demo Requirement (enforces demos exist and pass)
- **R330**: Demo Planning Requirements (enforces demos are planned)
- **R007**: Size Limits (demos excluded from count)
- **R304**: Line Counter Usage (demos measured separately)

## Remember

**"A simulated demo is worse than no demo - it creates false confidence"**
**"Demos must prove features work, not pretend they work"**
**"If a demo can't fail, it can't validate"**
**"External verification is the only real verification"**

Demos exist to provide confidence that implementations actually work. Simulated demos destroy that confidence by creating false positives. R331 ensures every demo is a REAL validation of working code.

## Changelog

- **2025-10-06**: Created R331 to eliminate simulation and false positives in demos
- Addresses critical gap: demos were being accepted without actual validation
- Prevents TODO-filled implementations from claiming to be "complete"
- Ensures flags used in demos actually exist in binaries
- Requires external verification (not just internal state)
- Mandates demos must be capable of failing
