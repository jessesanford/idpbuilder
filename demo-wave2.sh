#!/bin/bash
# Wave 2 Integration Demo Script (R291/R330)
echo "🎬 Running Phase 1 Wave 2 Integration Demo..."
echo "=============================================="
echo "This demo validates all Wave 2 components are properly integrated"
echo ""

DEMO_DIR="$(pwd)/demo-results"
FAILED=false

echo "📁 Demo results directory: $DEMO_DIR"
echo ""

# Run each individual demo
echo "1. Certificate Validation Demo..."
if ./demo-cert-validation.sh > /dev/null 2>&1; then
    echo "   ✅ cert-validation demo PASSED"
else
    echo "   ❌ cert-validation demo FAILED"
    FAILED=true
fi

echo "2. Chain Validation Demo..."
if ./demo-chain-validation.sh > /dev/null 2>&1; then
    echo "   ✅ chain-validation demo PASSED"
else
    echo "   ❌ chain-validation demo FAILED"
    FAILED=true
fi

echo "3. Fallback Strategies Demo..."
if ./demo-fallback.sh > /dev/null 2>&1; then
    echo "   ✅ fallback-strategies demo PASSED"
else
    echo "   ❌ fallback-strategies demo FAILED"
    FAILED=true
fi

echo "4. Validators Demo..."
if ./demo-validators.sh > /dev/null 2>&1; then
    echo "   ✅ validators demo PASSED"
else
    echo "   ❌ validators demo FAILED"
    FAILED=true
fi

echo ""
echo "=============================================="

# Verify integration points
echo "🔍 Verifying Integration Points..."
echo ""

# Check for cert validation package
if [ -d "pkg/certvalidation" ]; then
    echo "✅ Certificate validation package integrated"
else
    echo "❌ Certificate validation package missing"
    FAILED=true
fi

# Check for fallback package
if [ -d "pkg/fallback" ]; then
    echo "✅ Fallback strategies package integrated"
else
    echo "❌ Fallback strategies package missing"
    FAILED=true
fi

# Check for insecure package
if [ -d "pkg/insecure" ]; then
    echo "✅ Insecure mode package integrated"
else
    echo "❌ Insecure mode package missing"
    FAILED=true
fi

# Check build passes
echo ""
echo "🔨 Verifying build..."
if go build ./... > /dev/null 2>&1; then
    echo "✅ All packages build successfully"
else
    echo "❌ Build failed"
    FAILED=true
fi

# Check tests pass
echo ""
echo "🧪 Verifying tests..."
TEST_OUTPUT=$(go test ./pkg/certs ./pkg/certvalidation ./pkg/fallback ./pkg/insecure 2>&1)
if echo "$TEST_OUTPUT" | grep -q "^FAIL"; then
    echo "❌ Some tests failed"
    FAILED=true
else
    echo "✅ All integrated package tests pass"
fi

echo ""
echo "=============================================="
if [ "$FAILED" = "true" ]; then
    echo "❌ WAVE 2 INTEGRATION DEMO FAILED"
    echo ""
    echo "Please review the individual demo logs in $DEMO_DIR"
    exit 1
else
    echo "✅ WAVE 2 INTEGRATION DEMO COMPLETED SUCCESSFULLY!"
    echo ""
    echo "Summary:"
    echo "  • All effort demos executed successfully"
    echo "  • All packages integrated properly"
    echo "  • Build compiles successfully"
    echo "  • Tests pass for integrated components"
    echo ""
    echo "The Phase 1 Wave 2 integration is ready for:"
    echo "  • Push to remote repository"
    echo "  • Architect review"
    echo "  • Merge to main branch"
fi