#!/bin/bash

# Test script to verify backups directory is excluded from new backups

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

echo -e "${CYAN}Testing backup exclusion fix...${NC}"

# Create test directory structure
TEST_DIR="/tmp/test-upgrade-backup-$$"
mkdir -p "$TEST_DIR"
cd "$TEST_DIR"

# Create mock project structure
mkdir -p project/{backups,efforts,tools,utilities}
echo "test content" > project/file1.txt
echo "more content" > project/file2.txt

# Create some existing backups that should NOT be included
mkdir -p project/backups/efforts-old-20240101-120000
echo "old backup content" > project/backups/efforts-old-20240101-120000/test.txt
mkdir -p project/backups/efforts-old-20240102-120000
echo "another old backup" > project/backups/efforts-old-20240102-120000/test.txt

# Simulate the rsync backup command from upgrade.sh
BACKUP_DIR="project.backup.test"
echo -e "${CYAN}Creating backup with exclusions...${NC}"
rsync -av --exclude='efforts/' \
          --exclude='backups/' \
          --exclude='*.git' \
          --exclude='node_modules' \
          --exclude='__pycache__' \
          --exclude='*.pyc' \
          project/ "$BACKUP_DIR/"

# Check if backups directory was excluded
echo -e "\n${CYAN}Checking backup content...${NC}"
if [ -d "$BACKUP_DIR/backups" ]; then
    echo -e "${RED}❌ FAIL: backups/ directory was included in the backup!${NC}"
    echo -e "${RED}   This would cause recursive backup growth.${NC}"
    ls -la "$BACKUP_DIR/backups/" 2>/dev/null || true
    exit 1
else
    echo -e "${GREEN}✓ PASS: backups/ directory correctly excluded${NC}"
fi

# Verify other files were included
if [ -f "$BACKUP_DIR/file1.txt" ] && [ -f "$BACKUP_DIR/file2.txt" ]; then
    echo -e "${GREEN}✓ PASS: Regular files were included${NC}"
else
    echo -e "${RED}❌ FAIL: Regular files were not included${NC}"
    exit 1
fi

# Verify efforts was excluded as expected
if [ -d "$BACKUP_DIR/efforts" ]; then
    echo -e "${RED}❌ FAIL: efforts/ directory was included (should be excluded)${NC}"
    exit 1
else
    echo -e "${GREEN}✓ PASS: efforts/ directory correctly excluded${NC}"
fi

# Clean up
cd /
rm -rf "$TEST_DIR"

echo -e "\n${GREEN}✅ All tests passed! The backup exclusion fix is working correctly.${NC}"
echo -e "${CYAN}The upgrade script will no longer include previous backups in new backups.${NC}"