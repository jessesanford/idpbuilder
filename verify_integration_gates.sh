#!/bin/bash
echo "🔴🔴🔴 R291 MANDATORY GATES CHECK 🔴🔴🔴"
echo "Started: $(date '+%Y-%m-%d %H:%M:%S')"
echo ""

FAILED=false
PASSED_COUNT=0
FAILED_COUNT=0
EFFORTS_TESTED=""

# Phase 2 Wave 1 efforts (we only have these)
echo "=== PHASE 2 WAVE 1 DEMO VERIFICATION ==="
echo ""

# Check each effort branch for demo files
for effort in image-builder gitea-client gitea-client-split-001 gitea-client-split-002; do
    echo "🎬 [GATE] Checking demos in: $effort"
    
    # Check if branch exists
    if git show-ref --verify --quiet refs/remotes/integration/idpbuilder-oci-build-push/phase2/wave1/$effort; then
        echo "  ✅ Branch exists: idpbuilder-oci-build-push/phase2/wave1/$effort"
        
        # Check for demo files in branch
        if git ls-tree -r integration/idpbuilder-oci-build-push/phase2/wave1/$effort | grep -q "demo-features.sh"; then
            echo "  ✅ demo-features.sh found"
            ((PASSED_COUNT++))
            EFFORTS_TESTED="$EFFORTS_TESTED\n  ✅ $effort: DEMO IMPLEMENTED"
        else
            echo "  🔴 demo-features.sh missing!"
            ((FAILED_COUNT++))
            FAILED=true
            EFFORTS_TESTED="$EFFORTS_TESTED\n  ❌ $effort: NO DEMO FOUND"
        fi
        
        # Check for DEMO.md documentation
        if git ls-tree -r integration/idpbuilder-oci-build-push/phase2/wave1/$effort | grep -q "DEMO.md"; then
            echo "  ✅ DEMO.md documentation found"
        else
            echo "  ⚠️  DEMO.md documentation missing"
        fi
        
        # Check for test-data
        if git ls-tree -r integration/idpbuilder-oci-build-push/phase2/wave1/$effort | grep -q "test-data"; then
            echo "  ✅ test-data directory found"
        else
            echo "  ⚠️  test-data directory missing"
        fi
    else
        echo "  🔴 Branch not found!"
        ((FAILED_COUNT++))
        FAILED=true
        EFFORTS_TESTED="$EFFORTS_TESTED\n  ❌ $effort: BRANCH MISSING"
    fi
    echo ""
done

echo "======================================="
echo "📊 VERIFICATION SUMMARY"
echo "======================================="
echo -e "Efforts tested:$EFFORTS_TESTED"
echo ""
echo "✅ Passed: $PASSED_COUNT efforts"
echo "❌ Failed: $FAILED_COUNT efforts"
echo ""

if [ "$FAILED" = true ]; then
    echo "🔴🔴🔴 DEMO GATES FAILED - MUST ENTER ERROR_RECOVERY 🔴🔴🔴"
    exit 1
fi

echo "✅✅✅ ALL DEMO GATES PASSED ✅✅✅"
exit 0