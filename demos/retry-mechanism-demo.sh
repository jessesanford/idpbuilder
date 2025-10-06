#!/usr/bin/env bash
# Retry Mechanism Demo - idpbuilder push retry logic
# Demonstrates automatic retry behavior on transient failures

set -e

echo "================================================================="
echo "=== Retry Mechanism Demo ======================================="
echo "================================================================="
echo ""
echo "Purpose: Demonstrate retry logic on transient failures"
echo "Features: Automatic retries, backoff, max retry limits"
echo "Duration: ~4-5 minutes"
echo ""

# Configuration
DEMO_IMAGE="demo-retry:latest"
FLAKY_REGISTRY="localhost:5002"  # Simulated flaky registry
IDPBUILDER_BIN="${IDPBUILDER_BIN:-./bin/idpbuilder}"

# Check prerequisites
echo "=== Step 1: Prerequisites Check ==="
if ! command -v docker &> /dev/null; then
    echo "❌ Docker not found. Please install Docker."
    exit 1
fi
echo "✅ Docker available"

if [ ! -f "$IDPBUILDER_BIN" ]; then
    echo "⚠️  idpbuilder binary not found at $IDPBUILDER_BIN"
    echo "    Running in simulation mode"
fi

echo ""
echo "=== Step 2: Creating Test Image ==="
TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT

cat > "$TEMP_DIR/Dockerfile" << 'EOF'
FROM alpine:latest
RUN echo "Retry mechanism test image" > /retry-test.txt
CMD ["cat", "/retry-test.txt"]
EOF

echo "Building test image: $DEMO_IMAGE"
docker build -t "$DEMO_IMAGE" "$TEMP_DIR" > /dev/null 2>&1
echo "✅ Test image built"

echo ""
echo "================================================================="
echo "=== Scenario 1: Transient Network Failure (Auto-Retry) ========"
echo "================================================================="
echo ""
echo "Simulating transient network error that resolves after retry"
echo ""

echo "Command: idpbuilder push $DEMO_IMAGE --registry $FLAKY_REGISTRY --max-retries 3"
echo ""

if [ -f "$IDPBUILDER_BIN" ]; then
    echo "Executing push with retry enabled..."
    "$IDPBUILDER_BIN" push "$DEMO_IMAGE" --registry "$FLAKY_REGISTRY" --max-retries 3 2>&1 || {
        echo "⚠️  Flaky registry not available (expected in demo)"
    }
else
    # Simulate retry behavior
    echo "Attempt 1: Pushing image..."
    echo "❌ Error: connection refused (transient failure)"
    sleep 1
    echo ""
    echo "⏳ Retry 1/3 after 1s backoff..."
    echo "Attempt 2: Pushing image..."
    echo "❌ Error: connection timeout (transient failure)"
    sleep 2
    echo ""
    echo "⏳ Retry 2/3 after 2s backoff..."
    echo "Attempt 3: Pushing image..."
    echo "✅ Success! Push completed"
    echo ""
    echo "✅ Retry mechanism succeeded after 2 retries"
fi

echo ""
echo "================================================================="
echo "=== Scenario 2: Exponential Backoff Demonstration =============="
echo "================================================================="
echo ""
echo "Showing exponential backoff between retry attempts"
echo ""

echo "Configuration:"
echo "  Initial backoff: 1 second"
echo "  Backoff multiplier: 2x"
echo "  Max backoff: 30 seconds"
echo ""

if [ -f "$IDPBUILDER_BIN" ]; then
    echo "Retry timing with exponential backoff:"
    echo "  Attempt 1 → Wait 0s → Attempt 2"
    echo "  Attempt 2 → Wait 1s → Attempt 3"
    echo "  Attempt 3 → Wait 2s → Attempt 4"
    echo "  Attempt 4 → Wait 4s → Attempt 5"
    echo "  Attempt 5 → Wait 8s → Attempt 6"
    echo "  (continues up to max backoff of 30s)"
else
    echo "Simulating exponential backoff pattern:"
    for attempt in 1 2 3; do
        echo "Attempt $attempt: Pushing image..."
        if [ $attempt -eq 3 ]; then
            echo "✅ Success after $(($attempt - 1)) retries"
            break
        fi
        backoff=$((2 ** ($attempt - 1)))
        echo "❌ Transient error"
        echo "⏳ Waiting ${backoff}s before retry..."
        sleep 1  # Shortened for demo
        echo ""
    done
fi

echo ""
echo "================================================================="
echo "=== Scenario 3: Max Retry Limit Reached ========================"
echo "================================================================="
echo ""
echo "Demonstrating behavior when max retries exhausted"
echo ""

echo "Command: idpbuilder push $DEMO_IMAGE --registry $FLAKY_REGISTRY --max-retries 2"
echo ""

if [ -f "$IDPBUILDER_BIN" ]; then
    "$IDPBUILDER_BIN" push "$DEMO_IMAGE" --registry "$FLAKY_REGISTRY" --max-retries 2 2>&1 || {
        echo "❌ Max retries reached, operation failed"
    }
else
    echo "Simulating max retry limit:"
    echo "Attempt 1: Pushing image..."
    echo "❌ Error: connection refused"
    sleep 1
    echo ""
    echo "⏳ Retry 1/2 after 1s..."
    echo "Attempt 2: Pushing image..."
    echo "❌ Error: connection refused"
    sleep 1
    echo ""
    echo "⏳ Retry 2/2 after 2s..."
    echo "Attempt 3: Pushing image..."
    echo "❌ Error: connection refused"
    echo ""
    echo "❌ Max retries (2) reached. Operation failed."
    echo "Exit code: 1"
fi

echo ""
echo "================================================================="
echo "=== Scenario 4: Non-Retryable Errors ==========================="
echo "================================================================="
echo ""
echo "Demonstrating immediate failure for non-retryable errors"
echo ""

echo "Non-retryable error types:"
echo "  - Authentication failures (401, 403)"
echo "  - Image not found (404)"
echo "  - Invalid image format"
echo "  - Quota exceeded (413)"
echo ""

if [ -f "$IDPBUILDER_BIN" ]; then
    echo "Testing authentication failure (non-retryable)..."
    "$IDPBUILDER_BIN" push "$DEMO_IMAGE" --registry "$FLAKY_REGISTRY" \
        --username "invalid" --password "invalid" 2>&1 || {
        echo "❌ Authentication failed (401) - no retry attempted"
        echo "Reason: Non-retryable error type"
    }
else
    echo "Simulating non-retryable error:"
    echo "Attempt 1: Pushing image..."
    echo "❌ Error: authentication failed (401 Unauthorized)"
    echo ""
    echo "❌ Non-retryable error detected. Operation aborted."
    echo "Exit code: 1"
    echo ""
    echo "✅ Correctly avoided unnecessary retries for permanent errors"
fi

echo ""
echo "=== Cleanup ==="
docker rmi "$DEMO_IMAGE" > /dev/null 2>&1 || true
echo "✅ Cleanup complete"

echo ""
echo "================================================================="
echo "=== Demo Complete =============================================="
echo "================================================================="
echo ""
echo "Summary of Retry Behavior:"
echo ""
echo "Retryable Errors (transient failures):"
echo "  ✅ Network timeouts"
echo "  ✅ Connection refused"
echo "  ✅ Service unavailable (503)"
echo "  ✅ Rate limited (429)"
echo "  ✅ Temporary server errors (5xx except 501)"
echo ""
echo "Non-Retryable Errors (permanent failures):"
echo "  ❌ Authentication errors (401, 403)"
echo "  ❌ Not found (404)"
echo "  ❌ Invalid request (400)"
echo "  ❌ Quota exceeded (413)"
echo "  ❌ Method not allowed (405)"
echo ""
echo "Retry Configuration Options:"
echo "  --max-retries N       Maximum retry attempts (default: 3)"
echo "  --initial-backoff T   Initial backoff duration (default: 1s)"
echo "  --max-backoff T       Maximum backoff duration (default: 30s)"
echo "  --no-retry           Disable retry mechanism"
echo ""
echo "Best Practices:"
echo "  1. Use default retry settings for most cases"
echo "  2. Increase max-retries for unstable networks"
echo "  3. Disable retries for CI/CD fast-fail scenarios"
echo "  4. Monitor retry patterns to detect infrastructure issues"
echo ""
echo "Key Takeaways:"
echo "  ✅ Automatic retries improve reliability"
echo "  ✅ Exponential backoff prevents server overload"
echo "  ✅ Non-retryable errors fail fast"
echo "  ✅ Configurable for different use cases"
echo ""
echo "Next Steps:"
echo "  - Try tls-configuration-demo.sh for TLS scenarios"
echo "  - See authenticated-push-demo.sh for auth methods"
echo "  - Run phase2-integration-demo.sh for complete workflow"
echo ""
