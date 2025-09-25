#!/bin/bash

# Software Factory 2.0 Universal Installer
# Can be run from anywhere to set up SF 2.0

set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

echo -e "${CYAN}${BOLD}"
cat << "EOF"
╔═══════════════════════════════════════════════════════════════════╗
║            Software Factory 2.0 Universal Installer              ║
╚═══════════════════════════════════════════════════════════════════╝
EOF
echo -e "${NC}"

# Find where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo -e "${CYAN}Software Factory 2.0 Template Location:${NC}"
echo -e "  ${BOLD}$SCRIPT_DIR${NC}\n"

echo -e "${YELLOW}What would you like to do?${NC}"
echo -e "  ${BOLD}1)${NC} Create new SF 2.0 project"
echo -e "  ${BOLD}2)${NC} Migrate SF 1.0 project (full migration)"
echo -e "  ${BOLD}3)${NC} Migrate SF 1.0 planning only (discard code)"
echo -e "  ${BOLD}4)${NC} Show documentation"

echo -ne "\n${CYAN}Select option (1-4): ${NC}"
read -r choice

case $choice in
    1)
        echo -e "\n${CYAN}Starting SF 2.0 project setup...${NC}"
        cd "$SCRIPT_DIR"
        ./setup.sh
        ;;
    2)
        echo -e "\n${CYAN}Starting full SF 1.0 migration...${NC}"
        cd "$SCRIPT_DIR"
        ./migrate-from-1.0.sh
        ;;
    3)
        echo -e "\n${CYAN}Starting planning-only migration...${NC}"
        cd "$SCRIPT_DIR"
        ./migrate-planning-only.sh
        ;;
    4)
        echo -e "\n${CYAN}Available Documentation:${NC}"
        echo -e "  • ${BOLD}README.md${NC} - Getting started guide"
        echo -e "  • ${BOLD}MIGRATION-GUIDE-1.0-TO-2.0.md${NC} - Migration guide"
        echo -e "  • ${BOLD}SF-1.0-VS-2.0-COMPARISON.md${NC} - Feature comparison"
        echo -e "  • ${BOLD}quick-reference/${NC} - Quick reference guides"
        
        echo -ne "\n${CYAN}View README? (y/n): ${NC}"
        read -r view
        if [ "$view" = "y" ]; then
            less "$SCRIPT_DIR/README.md"
        fi
        ;;
    *)
        echo -e "${YELLOW}Invalid option. Please run again and select 1-4.${NC}"
        exit 1
        ;;
esac

echo -e "\n${GREEN}Done! 🚀${NC}"