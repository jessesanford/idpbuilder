#!/bin/bash
# Software Factory 2.0 - State Machine Visualization Demo
# Shows different visualization modes and options

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/.."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

echo -e "${CYAN}${BOLD}════════════════════════════════════════════════════════════════${NC}"
echo -e "${CYAN}${BOLD}    SOFTWARE FACTORY 2.0 - STATE VISUALIZATION DEMO${NC}"
echo -e "${CYAN}${BOLD}════════════════════════════════════════════════════════════════${NC}"
echo

# Function to pause between demos
pause_demo() {
    echo
    echo -e "${YELLOW}Press Enter to continue...${NC}"
    read -r
    clear
}

# Demo 1: Current State Visualization
echo -e "${GREEN}${BOLD}Demo 1: Current State Visualization${NC}"
echo -e "${BLUE}Shows your current position in the state machine${NC}"
echo
python3 tools/visualize-state-machine.py
pause_demo

# Demo 2: Compact State View
echo -e "${GREEN}${BOLD}Demo 2: Compact State View${NC}"
echo -e "${BLUE}Quick status check - minimal output${NC}"
echo
python3 tools/visualize-state-machine.py --compact
pause_demo

# Demo 3: Wave Execution Flow
echo -e "${GREEN}${BOLD}Demo 3: Wave Execution Flow${NC}"
echo -e "${BLUE}Shows how a wave progresses through states${NC}"
echo
python3 tools/visualize-state-flow.py --wave
pause_demo

# Demo 4: Integration Flow
echo -e "${GREEN}${BOLD}Demo 4: Integration Flow${NC}"
echo -e "${BLUE}Shows wave and phase integration patterns${NC}"
echo
python3 tools/visualize-state-flow.py --integration
pause_demo

# Demo 5: Fix Cascade Flow
echo -e "${GREEN}${BOLD}Demo 5: Fix Cascade Flow${NC}"
echo -e "${BLUE}Shows how fixes are handled when issues are found${NC}"
echo
python3 tools/visualize-state-flow.py --fix
pause_demo

# Demo 6: Split Handling Flow
echo -e "${GREEN}${BOLD}Demo 6: Split Handling Flow${NC}"
echo -e "${BLUE}Shows how large efforts are split into smaller pieces${NC}"
echo
python3 tools/visualize-state-flow.py --split
pause_demo

# Demo 7: Complete Flow (No Color)
echo -e "${GREEN}${BOLD}Demo 7: Complete Flow Without Colors${NC}"
echo -e "${BLUE}Useful for documentation or logging${NC}"
echo
python3 tools/visualize-state-machine.py --no-color | head -40
echo -e "${YELLOW}... (truncated for demo)${NC}"
echo

# Summary
echo -e "${CYAN}${BOLD}════════════════════════════════════════════════════════════════${NC}"
echo -e "${CYAN}${BOLD}                    DEMO COMPLETE${NC}"
echo -e "${CYAN}${BOLD}════════════════════════════════════════════════════════════════${NC}"
echo
echo -e "${GREEN}Available Visualization Tools:${NC}"
echo "  • visualize-state-machine.py - Current state and context"
echo "  • visualize-state-flow.py    - Flow diagrams and patterns"
echo
echo -e "${GREEN}Common Use Cases:${NC}"
echo "  • Debug stuck states:    python3 tools/visualize-state-machine.py"
echo "  • Check wave progress:   python3 tools/visualize-state-flow.py --wave"
echo "  • Document workflows:    python3 tools/visualize-state-flow.py > flow.txt"
echo
echo -e "${GREEN}For more information:${NC}"
echo "  • See tools/README-VISUALIZATION.md"
echo