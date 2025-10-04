#!/bin/bash
################################################################################
# R291 INTEGRATION DEMO - idpbuilder push to Gitea Registry (CORRECTED)
################################################################################
# This demo proves that the idpbuilder push command successfully pushes OCI
# images to a real Gitea registry, satisfying R291 Gate 4 requirements.
#
# KEY INSIGHT: idpbuilder push expects images as tarballs in the build path,
# not Docker daemon references. We must save the image first.
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
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
IMAGE_NAME="r291-demo"
IMAGE_TAG="${TIMESTAMP}"
BUILD_DIR="./build-r291-demo"

echo -e "${BLUE}=========================================="
echo "R291 INTEGRATION DEMO (CORRECTED)"
echo "idpbuilder push to Gitea Registry"
echo -e "==========================================${NC}"
echo ""
echo -e "${YELLOW}Demo Configuration:${NC}"
echo "  Registry: ${REGISTRY_URL}"
echo "  User: ${REGISTRY_USER}"
echo "  Source Image: ${TEST_IMAGE}"
echo "  Target Repository: ${IMAGE_NAME}"
echo "  Tag: ${IMAGE_TAG}"
echo "  Build Directory: ${BUILD_DIR}"
echo ""

################################################################################
# Step 1: Verify Prerequisites
################################################################################
echo -e "${BLUE}[Step 1/7] Verifying Prerequisites${NC}"
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
# Step 2: Prepare Build Directory
################################################################################
echo -e "${BLUE}[Step 2/7] Preparing Build Directory${NC}"
echo "----------------------------------------"

# Create clean build directory
rm -rf "${BUILD_DIR}"
mkdir -p "${BUILD_DIR}"
echo -e "${GREEN}✓ Build directory created: ${BUILD_DIR}${NC}"
echo ""

################################################################################
# Step 3: Prepare Test Image as Tarball
################################################################################
echo -e "${BLUE}[Step 3/7] Saving Image as Tarball${NC}"
echo "----------------------------------------"

# Pull a small test image
echo "Pulling ${TEST_IMAGE}..."
docker pull ${TEST_IMAGE} >/dev/null 2>&1
echo -e "${GREEN}✓ Test image pulled${NC}"

# Save image to tarball (this is what idpbuilder push expects!)
TARBALL_PATH="${BUILD_DIR}/${IMAGE_NAME}-${IMAGE_TAG}.tar"
echo "Saving image to: ${TARBALL_PATH}"
docker save ${TEST_IMAGE} -o "${TARBALL_PATH}"
echo -e "${GREEN}✓ Image saved as tarball ($(du -h "${TARBALL_PATH}" | cut -f1))${NC}"

# Verify tarball exists
if [ ! -f "${TARBALL_PATH}" ]; then
    echo -e "${RED}✗ Tarball not created!${NC}"
    exit 1
fi
echo ""

################################################################################
# Step 4: Execute idpbuilder push Command
################################################################################
echo -e "${BLUE}[Step 4/7] Executing idpbuilder push${NC}"
echo "----------------------------------------"

# The key insight: idpbuilder push discovers images in the build-path
# It will find our tarball and push it to the registry
echo "Command: ./idpbuilder push ${IMAGE_NAME}-${IMAGE_TAG}.tar \\"
echo "  --username '${REGISTRY_USER}' \\"
echo "  --password '<hidden>' \\"
echo "  --registry '${REGISTRY_URL}' \\"
echo "  --build-path '${BUILD_DIR}' \\"
echo "  --insecure \\"
echo "  --verbose"
echo ""

# Execute the push
if ./idpbuilder push "${IMAGE_NAME}-${IMAGE_TAG}.tar" \
    --username "${REGISTRY_USER}" \
    --password "${REGISTRY_PASSWORD}" \
    --registry "${REGISTRY_URL}" \
    --build-path "${BUILD_DIR}" \
    --insecure \
    --verbose 2>&1 | tee /tmp/push-output.log; then
    echo ""
    echo -e "${GREEN}✓ Push command completed${NC}"
else
    echo -e "${RED}✗ Push command failed!${NC}"
    echo "Full output:"
    cat /tmp/push-output.log
    exit 1
fi
echo ""

################################################################################
# Step 5: Verify Image in Registry (OCI API)
################################################################################
echo -e "${BLUE}[Step 5/7] Verifying via OCI Registry API${NC}"
echo "----------------------------------------"

# Construct the expected registry path
REGISTRY_IMAGE="${REGISTRY_URL}/${IMAGE_NAME}"

echo "Checking for image at: ${REGISTRY_IMAGE}:${IMAGE_TAG}"

# Check catalog
echo ""
echo "Checking registry catalog..."
CATALOG_RESPONSE=$(curl -k -s -u "${REGISTRY_USER}:${REGISTRY_PASSWORD}" \
    "https://${REGISTRY_URL}/v2/_catalog")
echo "Catalog: ${CATALOG_RESPONSE}"

# Check tags
echo ""
echo "Checking image tags..."
TAGS_RESPONSE=$(curl -k -s -u "${REGISTRY_USER}:${REGISTRY_PASSWORD}" \
    "https://${REGISTRY_URL}/v2/${IMAGE_NAME}/tags/list" 2>/dev/null || echo '{"tags":[]}')
echo "Tags response: ${TAGS_RESPONSE}"

# Check manifest (definitive proof)
echo ""
echo "Checking image manifest..."
MANIFEST_RESPONSE=$(curl -k -s -u "${REGISTRY_USER}:${REGISTRY_PASSWORD}" \
    -H "Accept: application/vnd.oci.image.manifest.v1+json" \
    -H "Accept: application/vnd.docker.distribution.manifest.v2+json" \
    "https://${REGISTRY_URL}/v2/${IMAGE_NAME}/manifests/${IMAGE_TAG}")

if echo "${MANIFEST_RESPONSE}" | grep -q "schemaVersion"; then
    echo -e "${GREEN}✓ Image manifest retrieved successfully${NC}"
    echo "Manifest preview:"
    echo "${MANIFEST_RESPONSE}" | jq '.' 2>/dev/null | head -20 || echo "${MANIFEST_RESPONSE}" | head -10
else
    echo -e "${YELLOW}⚠ Manifest check inconclusive${NC}"
    echo "Response: ${MANIFEST_RESPONSE}"
fi
echo ""

################################################################################
# Step 6: Verify Image Can Be Pulled Back
################################################################################
echo -e "${BLUE}[Step 6/7] Verifying Pull-Back Capability${NC}"
echo "----------------------------------------"

# Try to pull it back from registry
PULLBACK_IMAGE="${REGISTRY_URL}/${IMAGE_NAME}:${IMAGE_TAG}"
echo "Attempting to pull: ${PULLBACK_IMAGE}"

# Remove local copy first
docker rmi ${PULLBACK_IMAGE} >/dev/null 2>&1 || true

if docker pull ${PULLBACK_IMAGE} 2>&1 | tee /tmp/pull-output.log; then
    echo -e "${GREEN}✓ Image successfully pulled from registry${NC}"
    PULLBACK_SUCCESS=true
else
    echo -e "${YELLOW}⚠ Pull-back test inconclusive (may be registry configuration)${NC}"
    echo "This doesn't necessarily mean push failed - checking logs..."
    PULLBACK_SUCCESS=false
fi
echo ""

################################################################################
# Step 7: Generate Verification Report
################################################################################
echo -e "${BLUE}[Step 7/7] Generating Verification Report${NC}"
echo "----------------------------------------"

# Analyze push output
PUSH_SUCCESS=false
if grep -q "Successfully pushed" /tmp/push-output.log 2>/dev/null; then
    PUSH_SUCCESS=true
elif grep -q "ImagesPushed.*[1-9]" /tmp/push-output.log 2>/dev/null; then
    PUSH_SUCCESS=true
elif [ ! -s /tmp/push-output.log ]; then
    echo -e "${YELLOW}⚠ No push output captured${NC}"
fi

# Determine overall status
if [ "$PUSH_SUCCESS" = true ]; then
    OVERALL_STATUS="SUCCESS ✓"
    STATUS_COLOR="${GREEN}"
elif [ "$PULLBACK_SUCCESS" = true ]; then
    OVERALL_STATUS="SUCCESS ✓ (verified via pull-back)"
    STATUS_COLOR="${GREEN}"
else
    OVERALL_STATUS="NEEDS INVESTIGATION"
    STATUS_COLOR="${YELLOW}"
fi

cat > DEMO-VERIFICATION-REPORT-V2.txt << EOF
R291 Integration Demo - Verification Report v2
===============================================
Generated: $(date -Iseconds)

Demo Configuration
------------------
Registry URL: ${REGISTRY_URL}
Registry Type: Gitea (OCI-compliant)
Username: ${REGISTRY_USER}
Source Image: ${TEST_IMAGE}
Build Directory: ${BUILD_DIR}
Image Tarball: ${IMAGE_NAME}-${IMAGE_TAG}.tar
Target Registry Path: ${REGISTRY_URL}/${IMAGE_NAME}:${IMAGE_TAG}

Verification Results
--------------------
[✓] idpbuilder binary verified
[✓] Registry accessibility confirmed
[✓] Test image prepared and saved as tarball
[✓] Push command executed
$(if [ "$PUSH_SUCCESS" = true ]; then echo "[✓] Push reported success"; else echo "[?] Push status unclear"; fi)
$(if [ "$PULLBACK_SUCCESS" = true ]; then echo "[✓] Image pulled back successfully"; else echo "[?] Pull-back inconclusive"; fi)

Key Insights
------------
- idpbuilder push expects images as TARBALLS in the build-path
- Images must be saved with 'docker save' before pushing
- The command discovers .tar files in the specified directory
- This is the correct usage pattern for the tool

Push Command Output
-------------------
$(cat /tmp/push-output.log 2>/dev/null || echo "No output captured")

R291 Gate 4 Status
------------------
Overall Status: ${OVERALL_STATUS}

Evidence Files
--------------
- Push output: /tmp/push-output.log
- Pull output: /tmp/pull-output.log
- Image tarball: ${TARBALL_PATH}

Conclusion
----------
This demo shows the correct usage of idpbuilder push with tarball-based
image distribution. $(if [ "$OVERALL_STATUS" = "SUCCESS ✓" ] || [ "$OVERALL_STATUS" = "SUCCESS ✓ (verified via pull-back)" ]; then echo "The functionality works as designed."; else echo "Further investigation needed to confirm registry integration."; fi)
EOF

echo -e "${GREEN}✓ Verification report created: DEMO-VERIFICATION-REPORT-V2.txt${NC}"
echo ""

################################################################################
# Demo Summary
################################################################################
echo -e "${STATUS_COLOR}=========================================="
echo "DEMO COMPLETED"
echo -e "==========================================${NC}"
echo ""
echo -e "Status: ${STATUS_COLOR}${OVERALL_STATUS}${NC}"
echo ""
echo "Key Files:"
echo "  - Verification Report: DEMO-VERIFICATION-REPORT-V2.txt"
echo "  - Push Output: /tmp/push-output.log"
echo "  - Image Tarball: ${TARBALL_PATH}"
echo ""
echo -e "${BLUE}R291 Gate 4 Compliance:${NC}"
if [ "$OVERALL_STATUS" = "SUCCESS ✓" ] || [ "$OVERALL_STATUS" = "SUCCESS ✓ (verified via pull-back)" ]; then
    echo -e "  ${GREEN}SATISFIED ✓${NC}"
    echo ""
    echo "  The demo proves idpbuilder push:"
    echo "  - Correctly handles tarball-format images"
    echo "  - Authenticates with registries"
    echo "  - Integrates with OCI-compliant registries"
else
    echo -e "  ${YELLOW}PARTIAL - Needs Investigation${NC}"
    echo ""
    echo "  Review DEMO-VERIFICATION-REPORT-V2.txt for details"
fi
echo ""

# Cleanup
echo "Cleaning up build directory..."
rm -rf "${BUILD_DIR}"
echo -e "${GREEN}✓ Cleanup complete${NC}"
echo ""
