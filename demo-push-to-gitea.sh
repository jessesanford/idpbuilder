#!/bin/bash
################################################################################
# R291 INTEGRATION DEMO - idpbuilder push to Gitea Registry
################################################################################
# This demo proves that the idpbuilder push command successfully pushes OCI
# images to a real Gitea registry, satisfying R291 Gate 4 requirements.
################################################################################

set -e  # Exit on error
set -o pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
REGISTRY_URL="gitea.cnoe.localtest.me:8443"
REGISTRY_USER="giteaAdmin"
REGISTRY_PASSWORD='#@ko7q#J53roJt$'"'"'sG,%38xl^s<W3l=KYPENE1s}'
TEST_IMAGE="alpine:latest"
TARGET_REPO="giteaadmin/r291-demo"
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
TARGET_IMAGE="${REGISTRY_URL}/${TARGET_REPO}:${TIMESTAMP}"

echo -e "${BLUE}=========================================="
echo "R291 INTEGRATION DEMO"
echo "idpbuilder push to Gitea Registry"
echo -e "==========================================${NC}"
echo ""
echo -e "${YELLOW}Demo Configuration:${NC}"
echo "  Registry: ${REGISTRY_URL}"
echo "  User: ${REGISTRY_USER}"
echo "  Source Image: ${TEST_IMAGE}"
echo "  Target: ${TARGET_IMAGE}"
echo ""

################################################################################
# Step 1: Verify Prerequisites
################################################################################
echo -e "${BLUE}[Step 1/6] Verifying Prerequisites${NC}"
echo "----------------------------------------"

# Check binary exists
if [ ! -f "./idpbuilder" ]; then
    echo -e "${RED}✗ idpbuilder binary not found!${NC}"
    echo "Building from source..."
    make build
    if [ ! -f "./idpbuilder" ]; then
        echo -e "${RED}✗ Build failed!${NC}"
        exit 1
    fi
fi
echo -e "${GREEN}✓ idpbuilder binary exists${NC}"

# Check registry accessibility
if ! curl -k -s https://${REGISTRY_URL}/v2/ | grep -q "UNAUTHORIZED"; then
    echo -e "${RED}✗ Registry not accessible at https://${REGISTRY_URL}/v2/${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Registry is accessible${NC}"

# Check Docker is available
if ! command -v docker &> /dev/null; then
    echo -e "${RED}✗ Docker not found!${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Docker is available${NC}"
echo ""

################################################################################
# Step 2: Prepare Test Image
################################################################################
echo -e "${BLUE}[Step 2/6] Preparing Test Image${NC}"
echo "----------------------------------------"

# Pull a small test image
echo "Pulling ${TEST_IMAGE}..."
docker pull ${TEST_IMAGE} >/dev/null 2>&1
echo -e "${GREEN}✓ Test image ready${NC}"

# Tag for target registry
echo "Tagging as: ${TARGET_IMAGE}"
docker tag ${TEST_IMAGE} ${TARGET_IMAGE}
echo -e "${GREEN}✓ Image tagged for Gitea registry${NC}"
echo ""

################################################################################
# Step 3: Execute idpbuilder push Command
################################################################################
echo -e "${BLUE}[Step 3/6] Executing idpbuilder push${NC}"
echo "----------------------------------------"
echo "Command: ./idpbuilder push ${TARGET_IMAGE} --username <user> --password <hidden> --insecure --verbose"
echo ""

# Execute the push
if ./idpbuilder push ${TARGET_IMAGE} \
    --username "${REGISTRY_USER}" \
    --password "${REGISTRY_PASSWORD}" \
    --insecure \
    --verbose; then
    echo ""
    echo -e "${GREEN}✓ Push command completed successfully${NC}"
else
    echo -e "${RED}✗ Push command failed!${NC}"
    exit 1
fi
echo ""

################################################################################
# Step 4: Verify Image in Registry (OCI API)
################################################################################
echo -e "${BLUE}[Step 4/6] Verifying via OCI Registry API${NC}"
echo "----------------------------------------"

# Check catalog
echo "Checking registry catalog..."
CATALOG_RESPONSE=$(curl -k -s -u "${REGISTRY_USER}:${REGISTRY_PASSWORD}" \
    "https://${REGISTRY_URL}/v2/_catalog")
echo "Catalog response: ${CATALOG_RESPONSE}"

if echo "${CATALOG_RESPONSE}" | grep -q "${TARGET_REPO}"; then
    echo -e "${GREEN}✓ Repository found in catalog${NC}"
else
    echo -e "${YELLOW}⚠ Repository not in catalog (may be normal for Gitea)${NC}"
fi

# Check tags
echo ""
echo "Checking image tags..."
TAGS_RESPONSE=$(curl -k -s -u "${REGISTRY_USER}:${REGISTRY_PASSWORD}" \
    "https://${REGISTRY_URL}/v2/${TARGET_REPO}/tags/list" || echo '{"error":"endpoint may not be supported"}')
echo "Tags response: ${TAGS_RESPONSE}"

if echo "${TAGS_RESPONSE}" | grep -q "${TIMESTAMP}"; then
    echo -e "${GREEN}✓ Image tag found in registry${NC}"
else
    echo -e "${YELLOW}⚠ Tag not found via API (checking manifest...)${NC}"
fi

# Check manifest (definitive proof)
echo ""
echo "Checking image manifest..."
MANIFEST_RESPONSE=$(curl -k -s -u "${REGISTRY_USER}:${REGISTRY_PASSWORD}" \
    -H "Accept: application/vnd.oci.image.manifest.v1+json" \
    "https://${REGISTRY_URL}/v2/${TARGET_REPO}/manifests/${TIMESTAMP}")

if echo "${MANIFEST_RESPONSE}" | grep -q "schemaVersion"; then
    echo -e "${GREEN}✓ Image manifest retrieved successfully${NC}"
    echo "Manifest preview:"
    echo "${MANIFEST_RESPONSE}" | jq -r '.config.digest' 2>/dev/null || echo "${MANIFEST_RESPONSE}" | head -3
else
    echo -e "${RED}✗ Could not retrieve manifest${NC}"
    echo "Response: ${MANIFEST_RESPONSE}"
    exit 1
fi
echo ""

################################################################################
# Step 5: Verify Image Can Be Pulled Back
################################################################################
echo -e "${BLUE}[Step 5/6] Verifying Pull-Back Capability${NC}"
echo "----------------------------------------"

# Remove local copy
echo "Removing local image..."
docker rmi ${TARGET_IMAGE} >/dev/null 2>&1 || true
echo -e "${GREEN}✓ Local image removed${NC}"

# Pull it back from registry
echo "Pulling image back from Gitea registry..."
if docker pull ${TARGET_IMAGE} 2>&1 | tee /tmp/pull-output.log; then
    echo -e "${GREEN}✓ Image successfully pulled from registry${NC}"
else
    echo -e "${RED}✗ Pull failed${NC}"
    cat /tmp/pull-output.log
    exit 1
fi
echo ""

################################################################################
# Step 6: Generate Verification Report
################################################################################
echo -e "${BLUE}[Step 6/6] Generating Verification Report${NC}"
echo "----------------------------------------"

# Get image details
IMAGE_DIGEST=$(docker inspect ${TARGET_IMAGE} --format='{{.RepoDigests}}' 2>/dev/null || echo "N/A")
IMAGE_SIZE=$(docker inspect ${TARGET_IMAGE} --format='{{.Size}}' 2>/dev/null || echo "N/A")

cat > DEMO-VERIFICATION-REPORT.txt << EOF
R291 Integration Demo - Verification Report
============================================
Generated: $(date -Iseconds)

Demo Configuration
------------------
Registry URL: ${REGISTRY_URL}
Registry Type: Gitea (OCI-compliant)
Username: ${REGISTRY_USER}
Test Image: ${TEST_IMAGE}
Target Image: ${TARGET_IMAGE}

Verification Results
--------------------
[✓] idpbuilder binary verified
[✓] Registry accessibility confirmed
[✓] Test image prepared
[✓] Push command executed successfully
[✓] Image manifest retrievable from registry
[✓] Image can be pulled back from registry

Image Details
-------------
Digest: ${IMAGE_DIGEST}
Size: ${IMAGE_SIZE} bytes
Push Timestamp: ${TIMESTAMP}

Proof of Success
----------------
1. Push Command: SUCCEEDED
2. OCI Manifest API: RESPONDED
3. Pull-back Test: SUCCEEDED

R291 Gate 4 Status
------------------
SATISFIED ✓

The idpbuilder push command successfully:
- Authenticated with Gitea registry
- Pushed OCI image layers
- Created verifiable manifest
- Enabled pull-back verification

Conclusion
----------
This demo conclusively proves that idpbuilder push functionality
integrates correctly with the Gitea OCI registry and satisfies
all R291 Gate 4 requirements for production readiness.
EOF

echo -e "${GREEN}✓ Verification report created: DEMO-VERIFICATION-REPORT.txt${NC}"
echo ""

################################################################################
# Demo Summary
################################################################################
echo -e "${GREEN}=========================================="
echo "DEMO COMPLETED SUCCESSFULLY"
echo -e "==========================================${NC}"
echo ""
echo -e "${GREEN}✓ All verification steps passed${NC}"
echo ""
echo "Evidence of Success:"
echo "  1. Image pushed to: ${TARGET_IMAGE}"
echo "  2. Manifest accessible via OCI API"
echo "  3. Image successfully pulled back"
echo "  4. Verification report: DEMO-VERIFICATION-REPORT.txt"
echo ""
echo -e "${BLUE}R291 Gate 4: ${GREEN}SATISFIED ✓${NC}"
echo ""
echo "Next Steps:"
echo "  - Review DEMO-VERIFICATION-REPORT.txt"
echo "  - Review demo-output.log for complete execution log"
echo "  - Mark R291 compliance in project status"
echo ""
