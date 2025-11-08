#!/bin/bash

# add-schema-to-project.sh - Adds orchestrator-state.schema.json to existing projects
# This script ensures existing projects have the schema file for validation

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SF_DIR="$(dirname "$SCRIPT_DIR")"

echo -e "${CYAN}═══ Adding Schema File to Project ═══${NC}\n"

# Determine target directory
if [ -n "$1" ]; then
    TARGET_DIR="$1"
else
    TARGET_DIR="$(pwd)"
fi

# Verify target is a valid project directory
if [ ! -f "$TARGET_DIR/orchestrator-state-v3.json" ] && [ ! -f "$TARGET_DIR/.claude/CLAUDE.md" ]; then
    echo -e "${YELLOW}⚠️  Warning: Target doesn't appear to be a Software Factory project${NC}"
    echo "   Looking for orchestrator-state-v3.json or .claude/CLAUDE.md"
    echo ""
    read -p "Continue anyway? (y/N): " confirm
    if [[ ! "$confirm" =~ ^[Yy]$ ]]; then
        echo -e "${RED}Aborted${NC}"
        exit 1
    fi
fi

# Check if schema already exists
if [ -f "$TARGET_DIR/orchestrator-state.schema.json" ]; then
    echo -e "${GREEN}✓ Schema file already exists in project${NC}"
    echo "   Path: $TARGET_DIR/orchestrator-state.schema.json"
    
    # Check if it needs updating
    if [ -f "$SF_DIR/orchestrator-state.schema.json" ]; then
        if ! diff -q "$TARGET_DIR/orchestrator-state.schema.json" "$SF_DIR/orchestrator-state.schema.json" >/dev/null 2>&1; then
            echo -e "${YELLOW}⚠️  Schema file differs from template${NC}"
            read -p "Update to latest version? (y/N): " update
            if [[ "$update" =~ ^[Yy]$ ]]; then
                cp "$SF_DIR/orchestrator-state.schema.json" "$TARGET_DIR/"
                echo -e "${GREEN}✓ Schema file updated${NC}"
            fi
        else
            echo -e "${GREEN}✓ Schema file is up-to-date${NC}"
        fi
    fi
else
    # Copy schema from Software Factory template
    if [ -f "$SF_DIR/orchestrator-state.schema.json" ]; then
        echo -e "${CYAN}Copying schema from Software Factory template...${NC}"
        cp "$SF_DIR/orchestrator-state.schema.json" "$TARGET_DIR/"
        echo -e "${GREEN}✓ Schema file copied successfully${NC}"
        echo "   Path: $TARGET_DIR/orchestrator-state.schema.json"
    else
        echo -e "${RED}❌ Schema file not found in Software Factory template${NC}"
        echo "   Expected at: $SF_DIR/orchestrator-state.schema.json"
        exit 1
    fi
fi

# Also ensure validate-state.sh is up-to-date
if [ -f "$SF_DIR/tools/validate-state.sh" ]; then
    echo -e "\n${CYAN}Checking validate-state.sh tool...${NC}"
    
    mkdir -p "$TARGET_DIR/tools"
    
    if [ -f "$TARGET_DIR/tools/validate-state.sh" ]; then
        if ! diff -q "$TARGET_DIR/tools/validate-state.sh" "$SF_DIR/tools/validate-state.sh" >/dev/null 2>&1; then
            echo -e "${YELLOW}⚠️  validate-state.sh differs from template${NC}"
            read -p "Update to latest version? (y/N): " update
            if [[ "$update" =~ ^[Yy]$ ]]; then
                cp "$SF_DIR/tools/validate-state.sh" "$TARGET_DIR/tools/"
                chmod +x "$TARGET_DIR/tools/validate-state.sh"
                echo -e "${GREEN}✓ validate-state.sh updated${NC}"
            fi
        else
            echo -e "${GREEN}✓ validate-state.sh is up-to-date${NC}"
        fi
    else
        cp "$SF_DIR/tools/validate-state.sh" "$TARGET_DIR/tools/"
        chmod +x "$TARGET_DIR/tools/validate-state.sh"
        echo -e "${GREEN}✓ validate-state.sh copied${NC}"
    fi
fi

# Test validation
echo -e "\n${CYAN}Testing validation...${NC}"
if [ -f "$TARGET_DIR/orchestrator-state-v3.json" ]; then
    if [ -f "$TARGET_DIR/tools/validate-state.sh" ]; then
        if "$TARGET_DIR/tools/validate-state.sh" "$TARGET_DIR/orchestrator-state-v3.json" >/dev/null 2>&1; then
            echo -e "${GREEN}✓ Validation test passed${NC}"
        else
            echo -e "${YELLOW}⚠️  Validation test failed - running full validation for details:${NC}"
            echo ""
            "$TARGET_DIR/tools/validate-state.sh" "$TARGET_DIR/orchestrator-state-v3.json" || true
        fi
    else
        echo -e "${YELLOW}⚠️  validate-state.sh not found, skipping test${NC}"
    fi
else
    echo -e "${YELLOW}⚠️  No orchestrator-state-v3.json to validate${NC}"
fi

echo -e "\n${GREEN}═══ Schema Setup Complete ═══${NC}"
echo ""
echo "You can now validate your orchestrator state with:"
echo "  ${CYAN}tools/validate-state.sh${NC}"
echo ""
echo "The validation will automatically find the schema in:"
echo "  1. Your project directory"
echo "  2. \$SOFTWARE_FACTORY_DIR environment variable"
echo "  3. \$CLAUDE_PROJECT_DIR environment variable"
echo "  4. Known Software Factory locations"