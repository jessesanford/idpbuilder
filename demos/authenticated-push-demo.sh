#!/usr/bin/env bash
# Authenticated Push Demo - idpbuilder push with authentication
# Demonstrates various authentication methods for registry access

set -e

echo "================================================================="
echo "=== Authenticated Push Demo ===================================="
echo "================================================================="
echo ""
echo "Purpose: Demonstrate authenticated push with credentials"
echo "Features: Multiple auth methods (env vars, config, stdin)"
echo "Duration: ~3-4 minutes"
echo ""

# Configuration
DEMO_IMAGE="demo-auth-push:latest"
TEST_REGISTRY="localhost:5001"  # Separate registry for auth testing
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
    echo "    Running in demonstration mode"
fi

echo ""
echo "=== Step 2: Creating Test Image ==="
TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT

cat > "$TEMP_DIR/Dockerfile" << 'EOF'
FROM alpine:latest
RUN echo "Authenticated push test image" > /auth-test.txt
CMD ["cat", "/auth-test.txt"]
EOF

echo "Building test image: $DEMO_IMAGE"
docker build -t "$DEMO_IMAGE" "$TEMP_DIR" > /dev/null 2>&1
echo "✅ Test image built"

echo ""
echo "================================================================="
echo "=== Method 1: Environment Variable Authentication =============="
echo "================================================================="
echo ""
echo "Using REGISTRY_USERNAME and REGISTRY_PASSWORD environment variables"
echo ""

# Demonstrate env var method
export REGISTRY_USERNAME="demo-user"
export REGISTRY_PASSWORD="demo-pass"

echo "Environment variables set:"
echo "  REGISTRY_USERNAME=demo-user"
echo "  REGISTRY_PASSWORD=***hidden***"
echo ""
echo "Command: idpbuilder push $DEMO_IMAGE --registry $TEST_REGISTRY"
echo ""

if [ -f "$IDPBUILDER_BIN" ]; then
    echo "Executing push with env var auth..."
    "$IDPBUILDER_BIN" push "$DEMO_IMAGE" --registry "$TEST_REGISTRY" 2>&1 || {
        echo "⚠️  Auth registry not available (expected in demo)"
    }
else
    echo "Simulating: Push would use credentials from environment"
    echo "✅ Environment variable authentication demonstrated"
fi

echo ""
echo "================================================================="
echo "=== Method 2: Docker Config File Authentication ================"
echo "================================================================="
echo ""
echo "Using existing Docker credentials from ~/.docker/config.json"
echo ""

echo "Command: idpbuilder push $DEMO_IMAGE --registry $TEST_REGISTRY"
echo "         (credentials read from Docker config automatically)"
echo ""

if [ -f "$IDPBUILDER_BIN" ]; then
    # Unset env vars to test config file method
    unset REGISTRY_USERNAME
    unset REGISTRY_PASSWORD

    echo "Executing push with Docker config auth..."
    "$IDPBUILDER_BIN" push "$DEMO_IMAGE" --registry "$TEST_REGISTRY" 2>&1 || {
        echo "⚠️  Auth registry not available (expected in demo)"
    }
else
    echo "Simulating: Push would use Docker config.json credentials"
    echo "✅ Docker config authentication demonstrated"
fi

echo ""
echo "================================================================="
echo "=== Method 3: Stdin Authentication (CI/CD friendly) ============"
echo "================================================================="
echo ""
echo "Passing credentials via stdin (secure for CI/CD pipelines)"
echo ""

echo "Command: echo 'username:password' | idpbuilder push $DEMO_IMAGE --registry $TEST_REGISTRY --auth-stdin"
echo ""

if [ -f "$IDPBUILDER_BIN" ]; then
    echo "demo-user:demo-pass" | "$IDPBUILDER_BIN" push "$DEMO_IMAGE" --registry "$TEST_REGISTRY" --auth-stdin 2>&1 || {
        echo "⚠️  Auth registry not available (expected in demo)"
    }
else
    echo "Simulating: Push would receive credentials from stdin"
    echo "✅ Stdin authentication demonstrated"
fi

echo ""
echo "================================================================="
echo "=== Method 4: Explicit Credentials via Flags ==================="
echo "================================================================="
echo ""
echo "Using --username and --password flags (not recommended for security)"
echo ""

echo "Command: idpbuilder push $DEMO_IMAGE --registry $TEST_REGISTRY \\"
echo "           --username demo-user --password ***hidden***"
echo ""

if [ -f "$IDPBUILDER_BIN" ]; then
    "$IDPBUILDER_BIN" push "$DEMO_IMAGE" --registry "$TEST_REGISTRY" \
        --username "demo-user" --password "demo-pass" 2>&1 || {
        echo "⚠️  Auth registry not available (expected in demo)"
    }
else
    echo "Simulating: Push would use credentials from flags"
    echo "⚠️  Warning: Command-line credentials visible in process list"
    echo "✅ Flag-based authentication demonstrated (use with caution)"
fi

echo ""
echo "=== Cleanup ==="
docker rmi "$DEMO_IMAGE" > /dev/null 2>&1 || true
unset REGISTRY_USERNAME REGISTRY_PASSWORD
echo "✅ Cleanup complete"

echo ""
echo "================================================================="
echo "=== Demo Complete =============================================="
echo "================================================================="
echo ""
echo "Summary of Authentication Methods:"
echo ""
echo "1. Environment Variables (REGISTRY_USERNAME, REGISTRY_PASSWORD)"
echo "   ✅ Good for: Local development, simple scripts"
echo "   ⚠️  Warning: Variables may leak in logs"
echo ""
echo "2. Docker Config File (~/.docker/config.json)"
echo "   ✅ Good for: Using existing Docker credentials"
echo "   ✅ Automatic, no additional configuration needed"
echo ""
echo "3. Stdin Authentication (--auth-stdin)"
echo "   ✅ Good for: CI/CD pipelines, secure automation"
echo "   ✅ Best practice: Credentials not in env or flags"
echo ""
echo "4. Command-line Flags (--username, --password)"
echo "   ⚠️  Warning: Visible in process list"
echo "   ⚠️  Use only when other methods unavailable"
echo ""
echo "Recommended Priority:"
echo "  1. Docker config (automatic)"
echo "  2. Stdin (CI/CD secure)"
echo "  3. Environment variables (simple scripts)"
echo "  4. Flags (last resort only)"
echo ""
echo "Next Steps:"
echo "  - See retry-mechanism-demo.sh for failure handling"
echo "  - Try tls-configuration-demo.sh for TLS scenarios"
echo "  - Run phase2-integration-demo.sh for complete workflow"
echo ""
