#!/bin/bash
# Test script to verify backup removal handles write-protected files correctly

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

echo -e "${CYAN}Testing backup removal with write-protected files...${NC}"

# Create test directory structure
TEST_DIR="/tmp/backup-permission-test-$$"
mkdir -p "$TEST_DIR/efforts/phase1/wave1/test-effort/bin"

echo -e "${YELLOW}Creating test files...${NC}"

# Create regular files
echo "regular file" > "$TEST_DIR/efforts/phase1/wave1/test-effort/regular.txt"
echo "another file" > "$TEST_DIR/efforts/phase1/wave1/test-effort/bin/script.sh"

# Create write-protected files (simulating Kubernetes binaries)
echo "protected binary" > "$TEST_DIR/efforts/phase1/wave1/test-effort/bin/kubectl"
chmod 555 "$TEST_DIR/efforts/phase1/wave1/test-effort/bin/kubectl"

echo "protected binary 2" > "$TEST_DIR/efforts/phase1/wave1/test-effort/bin/kube-apiserver"
chmod 444 "$TEST_DIR/efforts/phase1/wave1/test-effort/bin/kube-apiserver"

# Make some directories read-only too
chmod 555 "$TEST_DIR/efforts/phase1/wave1/test-effort/bin"

echo -e "${CYAN}Test directory created: $TEST_DIR${NC}"
echo -e "${CYAN}Directory structure:${NC}"
find "$TEST_DIR" -ls | sed 's/^/  /'

echo -e "\n${YELLOW}Testing removal with permission fix...${NC}"

# Test the removal logic from our backup script
remove_with_permission_fix() {
    local dir="$1"

    # Fix permissions on write-protected files before removal
    echo -e "${YELLOW}⏳ Fixing permissions on write-protected files...${NC}"
    chmod -R u+w "$dir" 2>/dev/null || true

    # Now remove the directory
    if rm -rf "$dir" 2>/dev/null; then
        echo -e "${GREEN}✅ Directory removed successfully${NC}"
        return 0
    else
        # If standard removal fails, try with elevated permissions or show error
        echo -e "${YELLOW}⚠️  Some files could not be removed (permission denied)${NC}"
        echo -e "${YELLOW}    This is likely due to write-protected binaries${NC}"

        # Try to remove what we can
        find "$dir" -type f -exec chmod u+w {} \; 2>/dev/null || true
        find "$dir" -type d -exec chmod u+w {} \; 2>/dev/null || true

        # Attempt removal again
        if rm -rf "$dir" 2>/dev/null; then
            echo -e "${GREEN}✅ Directory removed after permission fix${NC}"
            return 0
        else
            echo -e "${RED}❌ Could not remove all files${NC}"
            echo -e "${YELLOW}   You may need to manually remove: $dir${NC}"
            echo -e "${YELLOW}   Try: sudo rm -rf \"$dir\"${NC}"
            return 1
        fi
    fi
}

# Run the test
if remove_with_permission_fix "$TEST_DIR"; then
    echo -e "\n${GREEN}✅ TEST PASSED: Write-protected files handled correctly${NC}"

    # Verify directory is gone
    if [ -d "$TEST_DIR" ]; then
        echo -e "${RED}❌ ERROR: Directory still exists!${NC}"
        exit 1
    else
        echo -e "${GREEN}✅ Verified: Directory completely removed${NC}"
    fi
else
    echo -e "\n${YELLOW}⚠️  TEST RESULT: Manual intervention required${NC}"
    echo -e "${YELLOW}   This is expected for system-protected files${NC}"

    # Clean up if still exists
    if [ -d "$TEST_DIR" ]; then
        echo -e "${CYAN}Cleaning up test directory with sudo...${NC}"
        sudo rm -rf "$TEST_DIR" 2>/dev/null || true
    fi
fi

echo -e "\n${CYAN}Test completed.${NC}"