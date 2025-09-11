#!/bin/bash
echo "🔴🔴🔴 R291 MANDATORY GATES CHECK (LOCAL) 🔴🔴🔴"
echo "Started: $(date '+%Y-%m-%d %H:%M:%S')"
echo ""

FAILED=false
PASSED_COUNT=0
FAILED_COUNT=0
WARNING_COUNT=0
EFFORTS_TESTED=""

# Phase 2 Wave 1 efforts
echo "=== PHASE 2 WAVE 1 LOCAL DEMO VERIFICATION ==="
echo ""

# Check each effort for demo files
for effort in image-builder gitea-client gitea-client-split-001 gitea-client-split-002; do
    echo "🎬 [GATE] Checking demos in: efforts/phase2/wave1/$effort"
    
    EFFORT_DIR="efforts/phase2/wave1/$effort"
    
    # Check if effort directory exists
    if [ -d "$EFFORT_DIR" ]; then
        echo "  ✅ Effort directory exists"
        
        # Check for demo-features.sh (MANDATORY)
        if [ -f "$EFFORT_DIR/demo-features.sh" ]; then
            echo "  ✅ demo-features.sh found"
            
            # Verify it's executable
            if [ -x "$EFFORT_DIR/demo-features.sh" ]; then
                echo "  ✅ demo-features.sh is executable"
            else
                echo "  ⚠️  demo-features.sh is not executable"
                ((WARNING_COUNT++))
            fi
            
            # Check file size (should have content)
            DEMO_SIZE=$(wc -c < "$EFFORT_DIR/demo-features.sh")
            if [ $DEMO_SIZE -gt 100 ]; then
                echo "  ✅ demo-features.sh has content ($DEMO_SIZE bytes)"
                ((PASSED_COUNT++))
                EFFORTS_TESTED="$EFFORTS_TESTED\n  ✅ $effort: DEMO IMPLEMENTED"
            else
                echo "  🔴 demo-features.sh is too small ($DEMO_SIZE bytes)"
                ((FAILED_COUNT++))
                FAILED=true
                EFFORTS_TESTED="$EFFORTS_TESTED\n  ❌ $effort: DEMO STUB ONLY"
            fi
        else
            echo "  🔴 demo-features.sh missing!"
            ((FAILED_COUNT++))
            FAILED=true
            EFFORTS_TESTED="$EFFORTS_TESTED\n  ❌ $effort: NO DEMO FOUND"
        fi
        
        # Check for DEMO.md documentation (RECOMMENDED)
        if [ -f "$EFFORT_DIR/DEMO.md" ]; then
            echo "  ✅ DEMO.md documentation found"
            DEMO_MD_SIZE=$(wc -c < "$EFFORT_DIR/DEMO.md")
            echo "  ✅ DEMO.md size: $DEMO_MD_SIZE bytes"
        else
            echo "  ⚠️  DEMO.md documentation missing"
            ((WARNING_COUNT++))
        fi
        
        # Check for test-data directory (OPTIONAL)
        if [ -d "$EFFORT_DIR/test-data" ]; then
            echo "  ✅ test-data directory found"
            TEST_FILES=$(ls -1 "$EFFORT_DIR/test-data" 2>/dev/null | wc -l)
            echo "  ✅ test-data contains $TEST_FILES files"
        else
            echo "  ℹ️  test-data directory not present (optional)"
        fi
    else
        echo "  🔴 Effort directory not found!"
        ((FAILED_COUNT++))
        FAILED=true
        EFFORTS_TESTED="$EFFORTS_TESTED\n  ❌ $effort: DIRECTORY MISSING"
    fi
    echo ""
done

# Check for wave-level demo orchestration
echo "🎬 [GATE] Checking wave-level demo orchestration"
if [ -f "efforts/phase2/wave1/wave-demo.sh" ]; then
    echo "  ✅ wave-demo.sh found"
    if [ -x "efforts/phase2/wave1/wave-demo.sh" ]; then
        echo "  ✅ wave-demo.sh is executable"
    else
        echo "  ⚠️  wave-demo.sh is not executable"
        ((WARNING_COUNT++))
    fi
else
    echo "  ⚠️  wave-demo.sh missing (recommended for orchestration)"
    ((WARNING_COUNT++))
fi
echo ""

echo "======================================="
echo "📊 VERIFICATION SUMMARY"
echo "======================================="
echo -e "Efforts tested:$EFFORTS_TESTED"
echo ""
echo "✅ Passed: $PASSED_COUNT efforts"
echo "❌ Failed: $FAILED_COUNT efforts"
echo "⚠️  Warnings: $WARNING_COUNT items"
echo ""

if [ "$FAILED" = true ]; then
    echo "🔴🔴🔴 DEMO GATES FAILED - MUST ENTER ERROR_RECOVERY 🔴🔴🔴"
    echo ""
    echo "Required Actions:"
    echo "1. All efforts must have demo-features.sh with actual content"
    echo "2. Demo scripts must demonstrate real functionality"
    echo "3. Documentation (DEMO.md) should explain the demo"
    exit 1
fi

if [ $WARNING_COUNT -gt 0 ]; then
    echo "⚠️  Some minor issues detected but not blocking"
fi

echo "✅✅✅ ALL MANDATORY DEMO GATES PASSED ✅✅✅"
echo "Ready for RETROFIT_COMPLETE state"
exit 0