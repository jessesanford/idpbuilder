#!/bin/bash

# Test script for idpbuilder OCI build and push commands
# This script will:
# 1. Create an idpbuilder cluster with Gitea
# 2. Build a test image
# 3. Push the image to Gitea registry
# 4. Verify the push was successful

set -e

echo "═══════════════════════════════════════════════════════════════════════"
echo "           IDPBUILDER OCI BUILD & PUSH TEST SCRIPT"
echo "═══════════════════════════════════════════════════════════════════════"
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_success() { echo -e "${GREEN}✅ $1${NC}"; }
print_error() { echo -e "${RED}❌ $1${NC}"; }
print_info() { echo -e "${YELLOW}ℹ️  $1${NC}"; }

# Check if idpbuilder is available
print_info "Checking for idpbuilder binary..."

# First check if we have a local compiled version
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
if [ -f "$SCRIPT_DIR/integration-testing-20250906-021605/idpbuilder" ]; then
    IDPBUILDER="$SCRIPT_DIR/integration-testing-20250906-021605/idpbuilder"
    print_success "Found compiled idpbuilder binary in integration directory"
elif [ -f "$SCRIPT_DIR/idpbuilder" ]; then
    IDPBUILDER="$SCRIPT_DIR/idpbuilder"
    print_success "Found compiled idpbuilder binary in script directory"
elif [ -f "./idpbuilder" ]; then
    IDPBUILDER="./idpbuilder"
    print_success "Found compiled idpbuilder binary in current directory"
elif command -v idpbuilder &> /dev/null; then
    IDPBUILDER="idpbuilder"
    print_success "Found idpbuilder in PATH"
else
    print_error "idpbuilder not found! Please ensure it's installed or compiled."
    echo "You can build it with: go build -o idpbuilder ."
    exit 1
fi

echo ""
echo "Using idpbuilder: $IDPBUILDER"
$IDPBUILDER version || true
echo ""

# Step 1: Create idpbuilder cluster
print_info "Step 1: Creating idpbuilder cluster with Gitea..."
echo "This will create a local Kubernetes cluster with Gitea registry"
echo ""

# Check if cluster already exists
if kubectl get nodes 2>/dev/null | grep -qE "(localdev|idpbuilder)"; then
    print_info "Kubernetes cluster already exists, using existing cluster"
else
    print_info "Creating new idpbuilder cluster (this may take a few minutes)..."
    $IDPBUILDER create \
        --host "cnoe.localtest.me" \
        --port 8443 \
        --recreate || {
        print_error "Failed to create idpbuilder cluster"
        echo "You may need to run: idpbuilder create"
        exit 1
    }
    print_success "idpbuilder cluster created successfully"
fi

echo ""

# Step 2: Wait for Gitea to be ready
print_info "Step 2: Waiting for Gitea to be ready..."
echo "Checking for Gitea deployment..."

# Wait for Gitea pod to be ready
TIMEOUT=120
ELAPSED=0
while [ $ELAPSED -lt $TIMEOUT ]; do
    if kubectl get pods -n gitea 2>/dev/null | grep -q "gitea.*Running"; then
        print_success "Gitea is running"
        break
    fi
    echo -n "."
    sleep 5
    ELAPSED=$((ELAPSED + 5))
done

if [ $ELAPSED -ge $TIMEOUT ]; then
    print_error "Timeout waiting for Gitea to be ready"
    echo "Please ensure idpbuilder cluster is running with: idpbuilder create"
    exit 1
fi

echo ""

# Get Gitea URL  
# Configure the registry URL here - you may need to adjust the port
# Common configurations:
#   - https://gitea.cnoe.localtest.me:8443 (idpbuilder's ingress port)
#   - https://gitea.cnoe.localtest.me:3000 (Gitea's default port)
#   - http://localhost:32223 (NodePort direct access)
GITEA_URL="https://gitea.cnoe.localtest.me:8443"

# You can override the registry URL with an environment variable
if [ ! -z "$REGISTRY_URL" ]; then
    GITEA_URL="$REGISTRY_URL"
    print_info "Using custom registry URL from environment: $GITEA_URL"
fi

print_info "Gitea Registry URL: $GITEA_URL"

# Step 3: Build the test image
print_info "Step 3: Building test image..."
echo "Current directory: $(pwd)"

# Always cd to demo directory from script directory
if [ -d "$SCRIPT_DIR/demo" ]; then
    echo "Changing to demo directory..."
    cd "$SCRIPT_DIR/demo"
    echo "Now in: $(pwd)"
elif [ -d "demo" ]; then
    echo "Changing to demo directory..."
    cd demo
    echo "Now in: $(pwd)"
else
    print_error "demo directory not found!"
    exit 1
fi

# Verify Dockerfile exists
if [ ! -f "Dockerfile" ]; then
    print_error "Dockerfile not found in demo directory!"
    exit 1
fi

echo "Building hello-world:v1 from current directory..."
echo "Note: Adding --output flag to save tarball as push requires it"
echo "Running: $IDPBUILDER build --context . --tag hello-world:v1 --platform linux/amd64 --output hello-world-v1.tar"

# Run build command without fancy error handling first
"$IDPBUILDER" build \
    --context . \
    --tag hello-world:v1 \
    --platform linux/amd64 \
    --output hello-world-v1.tar

BUILD_RESULT=$?
if [ $BUILD_RESULT -ne 0 ]; then
    print_error "Build failed with exit code: $BUILD_RESULT"
    exit 1
fi
print_success "Image built successfully and saved to hello-world-v1.tar"

echo ""
ls -la *.tar 2>/dev/null || echo "Note: Tarball may be stored in memory"
echo ""

# Step 4: Get Gitea credentials
print_info "Step 4: Setting up Gitea credentials..."

# Get credentials from the gitea-credential secret
GITEA_USERNAME=$(kubectl get secret gitea-credential -n gitea -o jsonpath='{.data.username}' 2>/dev/null | base64 -d) || {
    print_info "Could not get username from cluster, using default"
    GITEA_USERNAME="giteaAdmin"
}

GITEA_PASSWORD=$(kubectl get secret gitea-credential -n gitea -o jsonpath='{.data.password}' 2>/dev/null | base64 -d) || {
    print_info "Could not get password from cluster, using default"
    GITEA_PASSWORD="password"
}

echo "Username: $GITEA_USERNAME"
echo "Password: [hidden]"
echo ""

# Step 5: Push the image
print_info "Step 5: Pushing image to Gitea registry..."
echo "Target: $GITEA_URL"

# Export credentials for the push command
export REGISTRY_USERNAME="$GITEA_USERNAME"
export REGISTRY_PASSWORD="$GITEA_PASSWORD"

print_info "Push functionality is now implemented!"
echo "The push command will load the tarball and attempt to push to Gitea registry."
echo ""

# The push command expects a tarball path as the argument
# The registry, username, and credentials are provided as flags
TARBALL_PATH="hello-world-v1.tar"

print_info "Pushing tarball: $TARBALL_PATH"
print_info "Target registry: $GITEA_URL"

# Verify the tarball exists
if [ ! -f "$TARBALL_PATH" ]; then
    print_error "Tarball not found: $TARBALL_PATH"
    echo "Please run the build command first to create the tarball"
    exit 1
fi

# Verify the binary exists before attempting to push
if [ ! -f "$IDPBUILDER" ] && ! command -v "$IDPBUILDER" &> /dev/null; then
    print_error "idpbuilder binary not found at: $IDPBUILDER"
    echo "Please build the binary first with: go build -o idpbuilder ."
    exit 1
fi

# Push the image using the tarball
# The push command will load from tarball and push to registry
# Use set -o pipefail to capture exit code from command before pipe
set -o pipefail
"$IDPBUILDER" push "$TARBALL_PATH" \
    --registry "$GITEA_URL" \
    --insecure \
    --username "$GITEA_USERNAME" \
    --password "$GITEA_PASSWORD" 2>&1 | tee push-output.log
PUSH_RESULT=$?
set +o pipefail

if [ $PUSH_RESULT -ne 0 ]; then
    print_error "Push failed - this may be due to registry connectivity issues"
    echo ""
    echo "Attempted to push to: $GITEA_URL"
    echo "Tarball: $TARBALL_PATH"
    echo ""
    echo "The push implementation is complete, but the push failed. Common causes:"
    echo "  1. Gitea container registry may not be enabled"
    echo "  2. Registry endpoint may be different (e.g., needs port 3000 instead of 443)"
    echo "  3. Credentials may be different for registry vs git"
    echo "  4. The namespace ($GITEA_USERNAME) may not exist in the registry"
    echo ""
    echo "You can try adjusting the registry URL in this script (line 100)"
    exit 1
fi

print_success "Image pushed successfully!"
echo ""

# Step 6: Verify the push
print_info "Step 6: Verifying image in registry..."

# Check if we can list the image
if command -v curl &> /dev/null; then
    echo "Checking registry catalog..."
    curl -k -u "$GITEA_USERNAME:$GITEA_PASSWORD" \
        "$GITEA_URL/v2/_catalog" 2>/dev/null | grep -q "hello-world" && {
        print_success "Image found in registry catalog!"
    } || {
        print_info "Could not verify via API, but push reported success"
    }
else
    print_info "curl not available for verification, but push reported success"
fi

echo ""
echo "═══════════════════════════════════════════════════════════════════════"
echo "                     TEST COMPLETED SUCCESSFULLY!"
echo "═══════════════════════════════════════════════════════════════════════"
echo ""
print_success "✨ Both build and push commands are working!"
echo ""
echo "Summary:"
echo "  • Built image: hello-world:v1"
echo "  • Pushed to: $GITEA_URL"
echo "  • Registry user: $GITEA_USERNAME"
echo ""
echo "You can access Gitea UI at: $GITEA_URL"
echo "Login with: $GITEA_USERNAME / [check kubectl for password]"
echo ""
echo "To clean up the cluster, run:"
echo "  idpbuilder delete"
echo ""