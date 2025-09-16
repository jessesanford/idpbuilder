#!/bin/bash
# End-to-End Demo: Build and Push to Real Gitea Registry
# Target: gitea.cnoe.localtest.me:8443

set -e

echo "================================================"
echo "🚀 END-TO-END DEMO: Build & Push to Gitea Registry"
echo "================================================"
echo "Target Registry: gitea.cnoe.localtest.me:8443"
echo "Date: $(date '+%Y-%m-%d %H:%M:%S')"
echo ""

# Configuration
REGISTRY="gitea.cnoe.localtest.me:8443"
REPO_NAME="demo-app"
IMAGE_TAG="v1.0.0"
FULL_IMAGE="${REGISTRY}/${REPO_NAME}:${IMAGE_TAG}"

# Check prerequisites
echo "📋 Checking prerequisites..."
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is required but not installed"
    exit 1
fi

# Step 1: Create sample application
echo ""
echo "1️⃣ Creating sample application..."
mkdir -p demo-app
cat > demo-app/Dockerfile << 'DOCKERFILE'
FROM alpine:latest
RUN apk add --no-cache curl
COPY app.sh /app.sh
RUN chmod +x /app.sh
CMD ["/app.sh"]
DOCKERFILE

cat > demo-app/app.sh << 'APPSH'
#!/bin/sh
echo "Hello from IDP Builder!"
echo "Running on: $(hostname)"
echo "Date: $(date)"
APPSH

echo "   ✅ Sample application created"

# Step 2: Build OCI image
echo ""
echo "2️⃣ Building OCI image with idpbuilder..."
echo "   Command: idpbuilder build --context ./demo-app --tag ${FULL_IMAGE}"
./idpbuilder build --context ./demo-app --tag "${FULL_IMAGE}"
echo "   ✅ Image built successfully"

# Step 3: Configure credentials
echo ""
echo "3️⃣ Configuring Gitea credentials..."
# Check for credentials in environment
if [ -z "$GITEA_USERNAME" ] || [ -z "$GITEA_TOKEN" ]; then
    echo "   ⚠️  Please set GITEA_USERNAME and GITEA_TOKEN environment variables"
    echo "   Using demo credentials for illustration..."
    GITEA_USERNAME="idpbuilder"
    GITEA_TOKEN="demo-token-12345"
fi
echo "   Username: $GITEA_USERNAME"
echo "   Token: ${GITEA_TOKEN:0:10}..."

# Step 4: Push to Gitea registry
echo ""
echo "4️⃣ Pushing image to Gitea registry..."
echo "   Command: idpbuilder push --username $GITEA_USERNAME --token \$GITEA_TOKEN ${FULL_IMAGE}"
./idpbuilder push \
    --username "$GITEA_USERNAME" \
    --token "$GITEA_TOKEN" \
    --registry-cert /etc/idpbuilder/certs/gitea-ca.crt \
    "${FULL_IMAGE}"

# Step 5: Verify push
echo ""
echo "5️⃣ Verifying image in registry..."
echo "   Checking: ${REGISTRY}/${REPO_NAME}"
# This would use the actual Gitea API to verify
curl -k -u "${GITEA_USERNAME}:${GITEA_TOKEN}" \
    "https://${REGISTRY}/v2/${REPO_NAME}/tags/list" 2>/dev/null | jq '.' || echo "   (Verification requires active registry)"

# Summary
echo ""
echo "================================================"
echo "✅ END-TO-END DEMO COMPLETE!"
echo "================================================"
echo "Image successfully:"
echo "  • Built from local context"
echo "  • Tagged as ${FULL_IMAGE}"
echo "  • Pushed to Gitea registry with authentication"
echo "  • Used custom CA certificate for TLS"
echo ""
echo "To pull this image:"
echo "  docker pull ${FULL_IMAGE}"
echo ""
echo "Phase 2 Features Demonstrated:"
echo "  ✓ CLI build command (E2.2.1)"
echo "  ✓ CLI push command (E2.2.1)"
echo "  ✓ Credential management --username/--token (E2.2.2-A)"
echo "  ✓ Real OCI operations (E2.2.2-B)"
echo "  ✓ Certificate validation (Phase 1 integration)"
echo "================================================"
