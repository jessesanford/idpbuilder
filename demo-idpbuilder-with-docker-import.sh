#!/bin/bash
# Demo: Build with idpbuilder, import to Docker, then push
# This shows the current workflow limitation

set -e

echo "================================================"
echo "🔧 IDP Builder Demo with Docker Import Workaround"
echo "================================================"
echo "Date: $(date '+%Y-%m-%d %H:%M:%S')"
echo ""

# Configuration
REGISTRY="${REGISTRY:-gitea.cnoe.localtest.me:8443}"
IMAGE_NAME="${IMAGE_NAME:-idpbuilder-binary}"
IMAGE_TAG="${IMAGE_TAG:-latest}"
LOCAL_TAG="${IMAGE_NAME}:${IMAGE_TAG}"
FULL_IMAGE="${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}"

echo "📝 Note: Current limitation in Phase 2 implementation:"
echo "   - idpbuilder build creates images in temp storage"
echo "   - idpbuilder push expects images in Docker daemon"
echo "   - Workaround: Import to Docker first"
echo ""

# Step 1: Create context
echo "1️⃣ Creating context directory..."
mkdir -p idpbuilder-image-context
cp idpbuilder idpbuilder-image-context/
echo "   ✅ Context ready"

# Step 2: Build with idpbuilder (creates temp image)
echo ""
echo "2️⃣ Building OCI image with idpbuilder..."
./idpbuilder build \
    --context ./idpbuilder-image-context \
    --tag "${LOCAL_TAG}"
echo "   ✅ Image built (in temp storage)"

# Step 3: The missing link - need to get image into Docker
echo ""
echo "3️⃣ Current gap in implementation:"
echo "   ❌ idpbuilder build output is not in Docker daemon"
echo "   ❌ idpbuilder push looks for image in Docker daemon"
echo "   📋 This is the disconnect in Phase 2 implementation"

# Step 4: Workaround using Docker
echo ""
echo "4️⃣ Workaround: Use Docker to create the image..."
cat > idpbuilder-image-context/Dockerfile << 'DOCKERFILE'
FROM scratch
COPY idpbuilder /idpbuilder
DOCKERFILE

docker build -t "${FULL_IMAGE}" idpbuilder-image-context/
echo "   ✅ Image now in Docker daemon"

# Step 5: Now push works
echo ""
echo "5️⃣ Now idpbuilder push will work..."
if [ -z "$GITEA_USERNAME" ] || [ -z "$GITEA_TOKEN" ]; then
    echo "   ⚠️  Set GITEA_USERNAME and GITEA_TOKEN to push"
    echo ""
    echo "   export GITEA_USERNAME='your-username'"
    echo "   export GITEA_TOKEN='your-token'"
    echo "   ./idpbuilder push --username \$GITEA_USERNAME --token \$GITEA_TOKEN ${FULL_IMAGE}"
else
    echo "   Pushing with idpbuilder..."
    ./idpbuilder push \
        --username "$GITEA_USERNAME" \
        --token "$GITEA_TOKEN" \
        "${FULL_IMAGE}"
    echo "   ✅ Push complete!"
fi

# Cleanup
rm -rf idpbuilder-image-context

echo ""
echo "================================================"
echo "📊 Summary of Phase 2 Implementation Gap:"
echo "================================================"
echo "✅ idpbuilder build: Creates OCI images (but in temp storage)"
echo "✅ idpbuilder push: Pushes to registry (from Docker daemon)"
echo "❌ Missing: Bridge between build output and Docker daemon"
echo ""
echo "The build command needs to either:"
echo "  1. Save images directly to Docker daemon, OR"
echo "  2. Save to a format that push can read, OR"
echo "  3. Provide an 'import' command to bridge the gap"
echo "================================================"
