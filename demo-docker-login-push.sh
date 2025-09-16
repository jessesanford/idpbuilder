#!/bin/bash
# Workaround Demo: Use Docker login + Docker push
# Since idpbuilder push has a bug with authentication

set -e

echo "================================================"
echo "🔧 WORKAROUND: Docker Login + Docker Push"
echo "================================================"
echo "Note: idpbuilder push has a bug - credentials not passed to registry"
echo ""

# Configuration
REGISTRY="${REGISTRY:-gitea.cnoe.localtest.me:8443}"
IMAGE_NAME="${IMAGE_NAME:-demo/idpbuilder-binary}"
IMAGE_TAG="${IMAGE_TAG:-v1.0.0}"
FULL_IMAGE="${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}"

# Check prerequisites
if [ -z "$GITEA_USERNAME" ] || [ -z "$GITEA_TOKEN" ]; then
    echo "⚠️  Please set credentials first:"
    echo ""
    echo "export GITEA_USERNAME='giteaadmin'"
    echo "export GITEA_TOKEN='your-token'"
    echo ""
    exit 1
fi

# Step 1: Build image with Docker
echo "1️⃣ Building image..."
mkdir -p docker-build-context
cp idpbuilder docker-build-context/

cat > docker-build-context/Dockerfile << 'DOCKERFILE'
FROM scratch
COPY idpbuilder /idpbuilder
LABEL description="IDP Builder Binary"
DOCKERFILE

docker build -t "${FULL_IMAGE}" docker-build-context/
echo "   ✅ Image built"

# Step 2: Login to Docker registry
echo ""
echo "2️⃣ Logging in to Gitea registry..."
echo "${GITEA_TOKEN}" | docker login "${REGISTRY}" \
    --username "${GITEA_USERNAME}" \
    --password-stdin
echo "   ✅ Login successful"

# Step 3: Push with Docker
echo ""
echo "3️⃣ Pushing with Docker..."
docker push "${FULL_IMAGE}"
echo "   ✅ Push successful!"

# Step 4: Verify
echo ""
echo "4️⃣ Verifying in registry..."
curl -k -u "${GITEA_USERNAME}:${GITEA_TOKEN}" \
    "https://${REGISTRY}/v2/${IMAGE_NAME}/tags/list" 2>/dev/null | \
    python3 -m json.tool || echo "   (Manual verification needed)"

# Cleanup
rm -rf docker-build-context
docker logout "${REGISTRY}" >/dev/null 2>&1

echo ""
echo "================================================"
echo "✅ SUCCESS using Docker login + push"
echo "================================================"
echo "Image available at: ${FULL_IMAGE}"
echo ""
echo "⚠️  BUG IDENTIFIED IN PHASE 2:"
echo "   idpbuilder push accepts --username/--token"
echo "   but doesn't actually use them in registry.Push()"
echo "   See: pkg/registry/push.go lines 93-96"
echo "================================================"
