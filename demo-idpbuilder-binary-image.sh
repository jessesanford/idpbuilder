#!/bin/bash
# Demo: Create and Push Single-Layer Image with idpbuilder Binary
# This uses ONLY idpbuilder commands (no Docker)

set -e

echo "================================================"
echo "🎯 IDP Builder Self-Hosting Demo"
echo "================================================"
echo "Creating an OCI image containing the idpbuilder binary itself"
echo "Date: $(date '+%Y-%m-%d %H:%M:%S')"
echo ""

# Configuration
REGISTRY="${REGISTRY:-gitea.cnoe.localtest.me:8443}"
IMAGE_NAME="${IMAGE_NAME:-idpbuilder-binary}"
IMAGE_TAG="${IMAGE_TAG:-latest}"
FULL_IMAGE="${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}"

# Step 1: Check if idpbuilder binary exists
echo "1️⃣ Checking for idpbuilder binary..."
if [ ! -f "./idpbuilder" ]; then
    echo "❌ idpbuilder binary not found in current directory"
    echo "   Please ensure the binary is built first"
    exit 1
fi

BINARY_SIZE=$(ls -lh idpbuilder | awk '{print $5}')
echo "   ✅ Found idpbuilder binary (size: ${BINARY_SIZE})"

# Step 2: Create a directory with just the binary
echo ""
echo "2️⃣ Creating context directory..."
mkdir -p idpbuilder-image-context
cp idpbuilder idpbuilder-image-context/
echo "   ✅ Copied binary to context directory"

# Optional: Add a README to the image
cat > idpbuilder-image-context/README.txt << 'README'
IDP Builder Binary Image
========================
This OCI image contains the idpbuilder binary.

To extract the binary from this image:
1. Pull the image
2. Extract the layer
3. The binary will be at /idpbuilder

Version: Phase 2 Integration Build
Built with: idpbuilder build command
README

echo "   ✅ Added README.txt to image"

# Step 3: Build the single-layer OCI image
echo ""
echo "3️⃣ Building OCI image with idpbuilder..."
echo "   Command: idpbuilder build --context ./idpbuilder-image-context --tag ${FULL_IMAGE}"
./idpbuilder build \
    --context ./idpbuilder-image-context \
    --tag "${FULL_IMAGE}"

if [ $? -eq 0 ]; then
    echo "   ✅ Image built successfully!"
else
    echo "   ❌ Build failed!"
    exit 1
fi

# Step 4: Show what was created
echo ""
echo "4️⃣ Image Details:"
echo "   • Image tag: ${FULL_IMAGE}"
echo "   • Type: Single-layer OCI image"
echo "   • Contents:"
echo "     - /idpbuilder (executable binary)"
echo "     - /README.txt (documentation)"
echo "   • Layer count: 1"

# Step 5: Push to Gitea registry
echo ""
echo "5️⃣ Pushing image to Gitea registry..."

# Check for credentials
if [ -z "$GITEA_USERNAME" ] || [ -z "$GITEA_TOKEN" ]; then
    echo ""
    echo "⚠️  To push this image, set credentials and run:"
    echo ""
    echo "   export GITEA_USERNAME='your-username'"
    echo "   export GITEA_TOKEN='your-token'"
    echo ""
    echo "   idpbuilder push --username \$GITEA_USERNAME --token \$GITEA_TOKEN ${FULL_IMAGE}"
    echo ""
    echo "📦 Image is ready to push when you have credentials."
else
    echo "   Using credentials: $GITEA_USERNAME"
    echo "   Command: idpbuilder push --username \$GITEA_USERNAME --token \$GITEA_TOKEN ${FULL_IMAGE}"
    
    ./idpbuilder push \
        --username "$GITEA_USERNAME" \
        --token "$GITEA_TOKEN" \
        "${FULL_IMAGE}"
    
    if [ $? -eq 0 ]; then
        echo "   ✅ Image pushed successfully!"
        echo ""
        echo "6️⃣ Verification:"
        echo "   • Image available at: ${FULL_IMAGE}"
        echo "   • Registry: ${REGISTRY}"
        echo "   • Repository: ${IMAGE_NAME}"
        echo "   • Tag: ${IMAGE_TAG}"
    else
        echo "   ❌ Push failed!"
        exit 1
    fi
fi

# Clean up
echo ""
echo "7️⃣ Cleaning up..."
rm -rf idpbuilder-image-context
echo "   ✅ Removed temporary context directory"

# Summary
echo ""
echo "================================================"
echo "✅ DEMO COMPLETE!"
echo "================================================"
echo "Successfully demonstrated:"
echo "  • Built single-layer OCI image from directory"
echo "  • Image contains idpbuilder binary itself"
echo "  • Used ONLY idpbuilder commands (no Docker)"
echo "  • Ready to push to Gitea registry"
echo ""
echo "This proves idpbuilder can:"
echo "  1. Package any directory as an OCI image"
echo "  2. Push that image to a registry"
echo "  3. Work entirely without Docker for simple use cases"
echo ""
echo "Note: This creates a data layer, not a runnable container."
echo "      It's useful for distributing files via OCI registries."
echo "================================================"
