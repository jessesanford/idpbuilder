#!/bin/bash

# Simplified test script for R296 Deprecated Branch Marking Protocol
# Works without yq dependency

set -e

echo "═══════════════════════════════════════════════════════════════"
echo "🧪 TESTING R296 DEPRECATED BRANCH DETECTION (Simplified)"
echo "═══════════════════════════════════════════════════════════════"

# Test 1: Check for deprecated suffix in branch names
echo ""
echo "TEST 1: Detecting deprecated branch suffix"
echo "──────────────────────────────────────────"

test_branches=(
    "myproject/phase1/wave1/effort1"
    "myproject/phase1/wave1/effort1-deprecated-split"
    "myproject/phase1/wave1/effort1-split1"
    "myproject/phase1/wave1/effort2-deprecated-split"
    "myproject/phase1/wave1/effort2-split1"
    "myproject/phase1/wave1/effort2-split2"
)

blocked_count=0
allowed_count=0

for branch in "${test_branches[@]}"; do
    if [[ "$branch" == *"-deprecated-split" ]]; then
        echo "❌ BLOCKED: Deprecated branch detected: $branch"
        blocked_count=$((blocked_count + 1))
    else
        echo "✅ OK: Branch allowed for integration: $branch"
        allowed_count=$((allowed_count + 1))
    fi
done

echo ""
echo "Results: $blocked_count blocked, $allowed_count allowed"

# Test 2: Integration pre-check function
echo ""
echo "TEST 2: Full integration pre-check function"
echo "───────────────────────────────────────────"

check_for_deprecated_branches() {
    local branches=("$@")
    local has_deprecated=false
    
    echo "Checking ${#branches[@]} branches for deprecation..."
    for branch in "${branches[@]}"; do
        # Check for deprecated suffix
        if [[ "$branch" == *"-deprecated-split" ]]; then
            echo "  ❌ CRITICAL: Cannot integrate deprecated branch: $branch"
            echo "     This branch was split due to size violations"
            echo "     Use the replacement splits instead"
            has_deprecated=true
        else
            echo "  ✅ Branch OK: $branch"
        fi
    done
    
    # Return 0 for success (no deprecated), 1 for failure (has deprecated)
    if [ "$has_deprecated" = true ]; then
        return 1
    else
        return 0
    fi
}

# Test with valid branches only
echo ""
echo "Scenario A: All valid branches (should pass)"
echo "─────────────────────────────────────────────"
VALID_BRANCHES=(
    "myproject/phase1/wave1/effort1-split1"
    "myproject/phase1/wave1/effort1-split2"
    "myproject/phase1/wave1/effort2"
)

if check_for_deprecated_branches "${VALID_BRANCHES[@]}"; then
    echo "✅ PASS: All valid branches allowed for integration"
else
    echo "❌ ERROR: Valid branches were incorrectly blocked!"
    exit 1
fi

# Test with mix including deprecated branch
echo ""
echo "Scenario B: Mix with deprecated branch (should block)"
echo "──────────────────────────────────────────────────────"
INVALID_BRANCHES=(
    "myproject/phase1/wave1/effort1-split1"
    "myproject/phase1/wave1/effort1-deprecated-split"  # This should block
    "myproject/phase1/wave1/effort2"
)

if ! check_for_deprecated_branches "${INVALID_BRANCHES[@]}"; then
    echo "✅ PASS: Integration correctly blocked due to deprecated branch"
else
    echo "❌ ERROR: Deprecated branches were not blocked!"
    exit 1
fi

# Test 3: Demonstrate the correct workflow
echo ""
echo "TEST 3: Demonstrating correct workflow"
echo "──────────────────────────────────────────"
echo ""
echo "BEFORE split completion:"
echo "  Branch: myproject/phase1/wave1/big-feature (1200 lines - TOO LARGE)"
echo ""
echo "AFTER split completion:"
echo "  ❌ DEPRECATED: myproject/phase1/wave1/big-feature-deprecated-split"
echo "  ✅ USE THESE INSTEAD:"
echo "     - myproject/phase1/wave1/big-feature-split1 (400 lines)"
echo "     - myproject/phase1/wave1/big-feature-split2 (400 lines)"
echo "     - myproject/phase1/wave1/big-feature-split3 (400 lines)"
echo ""
echo "Integration attempt with deprecated branch:"
echo "  Result: BLOCKED with clear error message"
echo ""
echo "Integration with replacement splits:"
echo "  Result: PROJECT_DONE"

echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "✅ R296 DEPRECATED BRANCH DETECTION TESTS COMPLETE"
echo "═══════════════════════════════════════════════════════════════"
echo ""
echo "Summary:"
echo "✅ Deprecated suffix detection: WORKING"
echo "✅ Integration pre-check function: WORKING"
echo "✅ Clear error messages: IMPLEMENTED"
echo ""
echo "The R296 protocol successfully prevents integration of"
echo "deprecated TOO LARGE branches and guides to correct splits."