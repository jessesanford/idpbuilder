#!/bin/bash

# Standalone test for branch naming patterns

# Function to validate branch naming format (copied from enforce-state-validation.sh)
validate_branch_format() {
    local branch="$1"
    local effort_name="$2"
    local phase="$3"
    local wave="$4"

    # Determine branch type and validate accordingly
    local is_valid=false
    local branch_type=""

    # Check for different branch types based on R014 and other rules
    if [ -n "$PROJECT_PREFIX" ]; then
        # With project prefix
        if echo "$branch" | grep -qE "^${PROJECT_PREFIX}/phase[0-9]+/wave[0-9]+-integration$"; then
            # Wave integration branch: {prefix}/phase{X}/wave{Y}-integration
            branch_type="wave-integration"
            is_valid=true
        elif echo "$branch" | grep -qE "^${PROJECT_PREFIX}/phase[0-9]+-integration$"; then
            # Phase integration branch: {prefix}/phase{X}-integration
            branch_type="phase-integration"
            is_valid=true
        elif echo "$branch" | grep -qE "^${PROJECT_PREFIX}/integration$"; then
            # Project integration branch: {prefix}/integration
            branch_type="project-integration"
            is_valid=true
        elif echo "$branch" | grep -qE "^${PROJECT_PREFIX}/phase[0-9]+/wave[0-9]+/[-a-z0-9.-]+--split-[0-9]{3}$"; then
            # Split branch: {prefix}/phase{X}/wave{Y}/{effort-name}--split-{NNN}
            branch_type="split"
            is_valid=true
        elif echo "$branch" | grep -qE "^${PROJECT_PREFIX}/phase[0-9]+/wave[0-9]+/[-a-z0-9.-]+-fix$"; then
            # Fix branch: {prefix}/phase{X}/wave{Y}/{effort-name}-fix
            branch_type="fix"
            is_valid=true
        elif echo "$branch" | grep -qE "^${PROJECT_PREFIX}/phase[0-9]+/wave[0-9]+/[-a-z0-9.-]+$"; then
            # Standard effort branch: {prefix}/phase{X}/wave{Y}/{effort-name}
            branch_type="effort"
            is_valid=true
        fi
    else
        # Without project prefix
        if echo "$branch" | grep -qE "^phase[0-9]+/wave[0-9]+-integration$"; then
            # Wave integration branch: phase{X}/wave{Y}-integration
            branch_type="wave-integration"
            is_valid=true
        elif echo "$branch" | grep -qE "^phase[0-9]+-integration$"; then
            # Phase integration branch: phase{X}-integration
            branch_type="phase-integration"
            is_valid=true
        elif echo "$branch" | grep -qE "^integration$"; then
            # Project integration branch: integration
            branch_type="project-integration"
            is_valid=true
        elif echo "$branch" | grep -qE "^phase[0-9]+/wave[0-9]+/[-a-z0-9.-]+--split-[0-9]{3}$"; then
            # Split branch: phase{X}/wave{Y}/{effort-name}--split-{NNN}
            branch_type="split"
            is_valid=true
        elif echo "$branch" | grep -qE "^phase[0-9]+/wave[0-9]+/[-a-z0-9.-]+-fix$"; then
            # Fix branch: phase{X}/wave{Y}/{effort-name}-fix
            branch_type="fix"
            is_valid=true
        elif echo "$branch" | grep -qE "^phase[0-9]+/wave[0-9]+/[-a-z0-9.-]+$"; then
            # Standard effort branch: phase{X}/wave{Y}/{effort-name}
            branch_type="effort"
            is_valid=true
        fi
    fi

    if [ "$is_valid" = true ]; then
        # For effort/split/fix branches, verify phase/wave numbers match
        if [[ "$branch_type" =~ ^(effort|split|fix)$ ]]; then
            local branch_phase=$(echo "$branch" | sed 's/.*phase\([0-9]*\).*/\1/')
            local branch_wave=$(echo "$branch" | sed 's/.*wave\([0-9]*\).*/\1/')

            if [ "$branch_phase" == "$phase" ] && [ "$branch_wave" == "$wave" ]; then
                return 0
            else
                return 2  # Phase/wave mismatch
            fi
        else
            # For integration branches, just check if phase number matches (if applicable)
            if [[ "$branch_type" =~ ^(wave|phase)-integration$ ]]; then
                local branch_phase=$(echo "$branch" | sed 's/.*phase\([0-9]*\).*/\1/')
                if [ -n "$phase" ] && [ "$branch_phase" != "$phase" ]; then
                    return 2  # Phase mismatch
                fi
            fi
            return 0
        fi
    else
        return 1  # Format mismatch
    fi
}

# Test the validate_branch_format function
test_branch() {
    local branch="$1"
    local effort="$2"
    local phase="$3"
    local wave="$4"
    local expected="$5"

    validate_branch_format "$branch" "$effort" "$phase" "$wave"
    local result=$?

    if [ "$result" -eq 0 ]; then
        if [ "$expected" = "PASS" ]; then
            echo "✅ PASS: $branch"
        else
            echo "❌ FAIL: $branch (should have failed)"
        fi
    else
        if [ "$expected" = "FAIL" ]; then
            echo "✅ PASS: $branch (correctly rejected)"
        else
            echo "❌ FAIL: $branch (should have passed, error code: $result)"
        fi
    fi
}

echo "========================================="
echo "BRANCH NAMING VALIDATION TESTS"
echo "========================================="
echo ""

# Test with project prefix
export PROJECT_PREFIX="idpbuilderpush"

echo "=== WITH PROJECT PREFIX: $PROJECT_PREFIX ==="
echo ""

echo "--- Standard Effort Branches ---"
test_branch "idpbuilderpush/phase1/wave1/command-tests" "command-tests" "1" "1" "PASS"
test_branch "idpbuilderpush/phase2/wave3/authentication" "authentication" "2" "3" "PASS"
test_branch "phase1/wave1/command-tests" "command-tests" "1" "1" "FAIL"  # Missing prefix
test_branch "idpbuilderpush/phase1/wave1/Command-Tests" "command-tests" "1" "1" "FAIL"  # Uppercase

echo ""
echo "--- Integration Branches ---"
test_branch "idpbuilderpush/phase1/wave1-integration" "integration" "1" "1" "PASS"
test_branch "idpbuilderpush/phase1-integration" "integration" "1" "" "PASS"
test_branch "idpbuilderpush/integration" "integration" "" "" "PASS"
test_branch "idpbuilderpush/phase1/wave1/integration" "integration" "1" "1" "FAIL"  # Wrong format
test_branch "phase1/wave1-integration" "integration" "1" "1" "FAIL"  # Missing prefix

echo ""
echo "--- Split Branches ---"
test_branch "idpbuilderpush/phase1/wave1/effort-1.1.1--split-001" "effort-1.1.1" "1" "1" "PASS"
test_branch "idpbuilderpush/phase1/wave1/effort-1.1.1--split-002" "effort-1.1.1" "1" "1" "PASS"
test_branch "idpbuilderpush/phase1/wave1/api-types--split-123" "api-types" "1" "1" "PASS"
test_branch "idpbuilderpush/phase1/wave1/effort-1.1.1-split-1" "effort-1.1.1" "1" "1" "FAIL"  # Wrong format
test_branch "idpbuilderpush/phase1/wave1/effort-1.1.1/split-001" "effort-1.1.1" "1" "1" "FAIL"  # Wrong format

echo ""
echo "--- Fix Branches ---"
test_branch "idpbuilderpush/phase1/wave1/command-tests-fix" "command-tests" "1" "1" "PASS"
test_branch "idpbuilderpush/phase2/wave2/api-gateway-fix" "api-gateway" "2" "2" "PASS"
test_branch "idpbuilderpush/phase1/wave1/command-tests--fix" "command-tests" "1" "1" "FAIL"  # Double dash

echo ""
echo "=== WITHOUT PROJECT PREFIX ==="
export PROJECT_PREFIX=""

echo ""
echo "--- Standard Effort Branches ---"
test_branch "phase1/wave1/command-tests" "command-tests" "1" "1" "PASS"
test_branch "phase2/wave3/authentication" "authentication" "2" "3" "PASS"
test_branch "feature/command-tests" "command-tests" "1" "1" "FAIL"  # Wrong format

echo ""
echo "--- Integration Branches ---"
test_branch "phase1/wave1-integration" "integration" "1" "1" "PASS"
test_branch "phase1-integration" "integration" "1" "" "PASS"
test_branch "integration" "integration" "" "" "PASS"
test_branch "phase1/wave1/integration" "integration" "1" "1" "FAIL"  # Wrong format

echo ""
echo "--- Split Branches ---"
test_branch "phase1/wave1/effort-1.1.1--split-001" "effort-1.1.1" "1" "1" "PASS"
test_branch "phase1/wave1/effort-1.1.1--split-999" "effort-1.1.1" "1" "1" "PASS"

echo ""
echo "--- Phase/Wave Mismatch Tests ---"
test_branch "phase1/wave1/command-tests" "command-tests" "2" "1" "FAIL"  # Phase mismatch
test_branch "phase1/wave1/command-tests" "command-tests" "1" "2" "FAIL"  # Wave mismatch

echo ""
echo "========================================="
echo "TEST COMPLETE"
echo "========================================="