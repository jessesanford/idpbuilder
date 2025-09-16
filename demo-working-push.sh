#!/bin/bash
# Working Demo: Build with Docker, Push with idpbuilder
# This actually works with the current implementation

set -e

echo "================================================"
echo "✅ WORKING DEMO: Docker Build + idpbuilder Push"
echo "================================================"
echo "This demo uses the working combination:"
echo "  • Docker for building (creates image in daemon)"
echo "  • idpbuilder for pushing (with auth & certs)"
echo ""

# Configuration
REGISTRY="${REGISTRY:-gitea.cnoe.localtest.me:8443}"
IMAGE_NAME="${IMAGE_NAME:-demo/idpbuilder-binary}"
IMAGE_TAG="${IMAGE_TAG:-v1.0.0}"
FULL_IMAGE="${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}"

# Check prerequisites
if [ ! -f "./idpbuilder" ]; then
    echo "❌ idpbuilder binary not found"
    exit 1
fi

if [ -z "$GITEA_USERNAME" ] || [ -z "$GITEA_TOKEN" ]; then
    echo "⚠️  Please set credentials first:"
    echo ""
    echo "export GITEA_USERNAME='giteaadmin'"
    echo "export GITEA_TOKEN='your-token'"
    echo ""
    exit 1
fi

# Step 1: Create content directory
echo "1️⃣ Preparing content..."
mkdir -p docker-build-context
cp idpbuilder docker-build-context/

# Add metadata
cat > docker-build-context/info.txt << INFO
IDP Builder Binary
Version: Phase 2 Integration
Date: $(date '+%Y-%m-%d')
Size: $(ls -lh idpbuilder | awk '{print $5}')
INFO

echo "   ✅ Content ready"

# Step 2: Create Dockerfile
echo ""
echo "2️⃣ Creating Dockerfile..."
cat > docker-build-context/Dockerfile << 'DOCKERFILE'
FROM scratch
COPY idpbuilder /idpbuilder
COPY info.txt /info.txt
LABEL description="IDP Builder Binary Image"
LABEL version="phase2-integration"
DOCKERFILE
echo "   ✅ Dockerfile created"

# Step 3: Build with Docker
echo ""
echo "3️⃣ Building image with Docker..."
echo "   Command: docker build -t ${FULL_IMAGE} docker-build-context/"
docker build -t "${FULL_IMAGE}" docker-build-context/
echo "   ✅ Image built and stored in Docker daemon"

# Step 4: Verify image exists in Docker
echo ""
echo "4️⃣ Verifying image in Docker daemon..."
docker images | grep "${IMAGE_NAME}" | head -1
echo "   ✅ Image confirmed in Docker"

# Step 5: Push with idpbuilder
echo ""
echo "5️⃣ Pushing with idpbuilder..."
echo "   Registry: ${REGISTRY}"
echo "   Image: ${IMAGE_NAME}:${IMAGE_TAG}"
echo "   User: ${GITEA_USERNAME}"
echo ""
echo "   Command: idpbuilder push --username \$GITEA_USERNAME --token \$GITEA_TOKEN ${FULL_IMAGE}"

./idpbuilder push \
    --username "${GITEA_USERNAME}" \
    --token "${GITEA_TOKEN}" \
    "${FULL_IMAGE}"

if [ $? -eq 0 ]; then
    echo ""
    echo "   ✅ Push successful!"
    echo ""
    echo "6️⃣ Image is now available at:"
    echo "   ${FULL_IMAGE}"
    echo ""
    echo "   To pull this image:"
    echo "   docker pull ${FULL_IMAGE}"
else
    echo "   ❌ Push failed"
    exit 1
fi

# Cleanup
echo ""
echo "7️⃣ Cleaning up..."
rm -rf docker-build-context
echo "   ✅ Cleanup complete"

echo ""
echo "================================================"
echo "🎉 SUCCESS!"
echo "================================================"
echo "The image has been pushed to Gitea registry!"
echo ""
echo "What this demonstrated:"
echo "  ✅ Docker build creates images in daemon"
echo "  ✅ idpbuilder push works with Docker images"
echo "  ✅ Authentication with --username/--token"
echo "  ✅ Certificate handling (automatic)"
echo "  ✅ Full integration with Gitea registry"
echo "================================================"
