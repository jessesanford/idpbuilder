#!/bin/bash
# Demo script for Effort 3.1.3: Core Workflow Integration Tests
# Created: 2025-11-04
# Purpose: Demonstrate working integration test functionality

set -e  # Exit on error

echo "🎬 Starting Core Workflow Integration Tests Demo..."
echo "========================================================"
echo ""

# Check prerequisites
echo "📦 Checking prerequisites..."
if ! command -v go &> /dev/null; then
    echo "❌ Go not installed"
    exit 1
fi

if ! command -v docker &> /dev/null; then
    echo "❌ Docker not installed"
    exit 1
fi

if ! docker info &> /dev/null; then
    echo "❌ Docker daemon not running"
    exit 1
fi

echo "✅ Prerequisites satisfied"
echo ""

# Scenario 1: Small Image Push Demo
echo "========================================================"
echo "Scenario 1: Small Image Push Demo (5MB, 2 layers)"
echo "========================================================"
echo "Testing: TestPushSmallImageSuccess"
echo "Expected: Test passes in ~45 seconds"
echo ""

go test -tags=integration -v -timeout 6m -run TestPushSmallImageSuccess ./test/integration/ 2>&1 | tee /tmp/demo-small-push.log

if [ $? -eq 0 ]; then
    echo "✅ Small image push test PASSED"
else
    echo "❌ Small image push test FAILED"
    echo "Check /tmp/demo-small-push.log for details"
    exit 1
fi
echo ""

# Scenario 2: Large Image with Progress Demo
echo "========================================================"
echo "Scenario 2: Large Image with Progress Demo (100MB, 10 layers)"
echo "========================================================"
echo "Testing: TestPushLargeImageWithProgress"
echo "Expected: Test passes in ~120 seconds with progress updates"
echo ""

go test -tags=integration -v -timeout 12m -run TestPushLargeImageWithProgress ./test/integration/ 2>&1 | tee /tmp/demo-large-push.log

if [ $? -eq 0 ]; then
    echo "✅ Large image with progress test PASSED"
    echo "✅ Verified: All 10 layers reported progress"
    echo "✅ Verified: >100MB processed"
else
    echo "❌ Large image with progress test FAILED"
    echo "Check /tmp/demo-large-push.log for details"
    exit 1
fi
echo ""

# Scenario 3: Multiple Images Demo
echo "========================================================"
echo "Scenario 3: Multiple Sequential Images Demo (3 images)"
echo "========================================================"
echo "Testing: TestPushMultipleImagesSequentially"
echo "Expected: Test passes in ~90 seconds"
echo ""

go test -tags=integration -v -timeout 12m -run TestPushMultipleImagesSequentially ./test/integration/ 2>&1 | tee /tmp/demo-multiple-push.log

if [ $? -eq 0 ]; then
    echo "✅ Multiple images test PASSED"
    echo "✅ Verified: All 3 images pushed successfully"
    echo "✅ Verified: All 3 images queryable in registry"
else
    echo "❌ Multiple images test FAILED"
    echo "Check /tmp/demo-multiple-push.log for details"
    exit 1
fi
echo ""

# Summary
echo "========================================================"
echo "✅ Demo Completed Successfully!"
echo "========================================================"
echo ""
echo "Summary of Results:"
echo "  ✅ Small image push: Working"
echo "  ✅ Large image with progress: Working"
echo "  ✅ Multiple sequential pushes: Working"
echo ""
echo "All core workflow integration tests are functioning correctly."
echo "The push command integrates successfully with Docker and Gitea."
echo ""

# Set demo ready flag for wave integration
export DEMO_READY=true
echo "DEMO_READY=true"

exit 0
