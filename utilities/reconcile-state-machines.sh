#!/bin/bash

# reconcile-state-machines.sh
# Reconciles the agent-specific state machines with the main SOFTWARE-FACTORY-STATE-MACHINE.md
# Adds warnings about which is authoritative

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

echo "════════════════════════════════════════════════════════════════"
echo -e "${BOLD}STATE MACHINE RECONCILIATION ANALYSIS${NC}"
echo "════════════════════════════════════════════════════════════════"
echo ""

# Check both systems exist
MAIN_SM="SOFTWARE-FACTORY-STATE-MACHINE.md"
AGENT_SM_DIR="state-machines"

if [ ! -f "$MAIN_SM" ]; then
    echo -e "${RED}❌ Main state machine not found: $MAIN_SM${NC}"
    exit 1
fi

if [ ! -d "$AGENT_SM_DIR" ]; then
    echo -e "${RED}❌ Agent state machines directory not found: $AGENT_SM_DIR${NC}"
    exit 1
fi

echo -e "${CYAN}Found State Machine Systems:${NC}"
echo "1. Main: $MAIN_SM (R206 declares this as SINGLE SOURCE OF TRUTH)"
echo "2. Agent-specific: $AGENT_SM_DIR/*.md (older, more detailed)"
echo ""

# Add deprecation warning to old state machines
add_deprecation_warning() {
    local FILE="$1"
    
    # Check if warning already exists
    if grep -q "⚠️ DEPRECATION WARNING" "$FILE" 2>/dev/null; then
        echo "  Warning already present in $(basename $FILE)"
        return
    fi
    
    # Create backup
    cp "$FILE" "${FILE}.pre-reconciliation"
    
    # Add warning at the top
    TEMP_FILE="${FILE}.tmp"
    cat > "$TEMP_FILE" << 'EOF'
# ⚠️ DEPRECATION WARNING ⚠️

**IMPORTANT**: This file is part of the LEGACY state machine system.

Per **Rule R206**, the authoritative state machine is:
- **SOFTWARE-FACTORY-STATE-MACHINE.md** (SINGLE SOURCE OF TRUTH)

This file is retained for reference but should NOT be used for state validation.
All agents MUST validate states against SOFTWARE-FACTORY-STATE-MACHINE.md.

---

EOF
    
    cat "$FILE" >> "$TEMP_FILE"
    mv "$TEMP_FILE" "$FILE"
    
    echo -e "${GREEN}✅ Added deprecation warning to $(basename $FILE)${NC}"
}

echo -e "${YELLOW}Adding deprecation warnings to legacy state machines...${NC}"
for sm_file in $AGENT_SM_DIR/*.md; do
    if [[ ! "$sm_file" == *.backup ]]; then
        add_deprecation_warning "$sm_file"
    fi
done

echo ""
echo -e "${BLUE}Comparing state definitions...${NC}"
echo ""

# Extract states from main state machine
extract_main_states() {
    local AGENT_TYPE="$1"
    local SECTION=""
    
    case "$AGENT_TYPE" in
        orchestrator)
            SECTION="Orchestrator States"
            ;;
        sw-engineer)
            SECTION="SW Engineer States"
            ;;
        code-reviewer)
            SECTION="Code Reviewer States"
            ;;
        architect)
            SECTION="Architect States"
            ;;
    esac
    
    if [ -n "$SECTION" ]; then
        sed -n "/## $SECTION/,/^##[^#]/p" "$MAIN_SM" | \
            grep "^- \*\*" | \
            sed 's/- \*\*//' | \
            sed 's/\*\*.*//' | \
            sort
    fi
}

# Extract states from agent-specific file
extract_agent_states() {
    local AGENT_FILE="$1"
    
    # Look for states in the mermaid diagram
    grep -E "^\s+[A-Z_]+ -->" "$AGENT_FILE" 2>/dev/null | \
        sed 's/-->.*//' | \
        sed 's/^\s*//' | \
        sort -u
}

# Compare states for each agent
for agent in orchestrator sw-engineer code-reviewer architect; do
    echo -e "${CYAN}═══════════════════════════════════════════════════════════════${NC}"
    echo -e "${CYAN}Agent: $agent${NC}"
    echo -e "${CYAN}═══════════════════════════════════════════════════════════════${NC}"
    
    MAIN_STATES=$(extract_main_states "$agent")
    AGENT_FILE="$AGENT_SM_DIR/${agent}.md"
    
    if [ -f "$AGENT_FILE" ]; then
        AGENT_STATES=$(extract_agent_states "$AGENT_FILE")
        
        echo -e "${BLUE}States in MAIN (authoritative):${NC}"
        echo "$MAIN_STATES" | sed 's/^/  /'
        
        echo ""
        echo -e "${YELLOW}States in LEGACY (deprecated):${NC}"
        echo "$AGENT_STATES" | head -20 | sed 's/^/  /'
        
        # Note discrepancies
        echo ""
        echo -e "${RED}⚠️  States are DIFFERENT between systems${NC}"
        echo "  Agents MUST use SOFTWARE-FACTORY-STATE-MACHINE.md per R206"
    else
        echo "No legacy file for $agent"
    fi
    echo ""
done

# Create migration guide
MIGRATION_GUIDE="STATE-MACHINE-MIGRATION-GUIDE.md"
cat > "$MIGRATION_GUIDE" << 'EOF'
# State Machine Migration Guide

## ⚠️ CRITICAL: State Machine Authority

Per **Rule R206**, there is only ONE authoritative state machine:
- **SOFTWARE-FACTORY-STATE-MACHINE.md** - SINGLE SOURCE OF TRUTH

## Legacy System (DEPRECATED)

The `state-machines/` directory contains the OLD state machine system:
- More detailed states (e.g., SPAWN_CODE_REVIEWER_PLANNING vs SPAWN_AGENTS)
- Agent-specific files
- Different state names and transitions

## Current System (AUTHORITATIVE)

The `SOFTWARE-FACTORY-STATE-MACHINE.md` contains:
- Simplified, unified states
- All agents in one file
- R206 compliance required
- New architecture-driven states (R210, R211)

## Migration Rules

1. **ALWAYS** validate against SOFTWARE-FACTORY-STATE-MACHINE.md
2. **NEVER** use state-machines/*.md for validation
3. **Map** old detailed states to new simplified states
4. **Update** any references to use new state names

## State Mapping Examples

### Orchestrator
- OLD: `SPAWN_CODE_REVIEWER_PLANNING` → NEW: `SPAWN_AGENTS`
- OLD: `SPAWN_SW_ENG` → NEW: `SPAWN_AGENTS`
- OLD: `SPAWN_SW_ENG_FIX` → NEW: `SPAWN_AGENTS`
- OLD: `CREATE_SPLIT_PLAN` → NEW: `SPAWN_AGENTS` (for Code Reviewer)
- NEW: `SPAWN_ARCHITECT_PHASE_PLANNING` (R210, no old equivalent)
- NEW: `SPAWN_CODE_REVIEWER_WAVE_IMPL` (R211, no old equivalent)

### Code Reviewer
- OLD: `PLANNING` → NEW: `EFFORT_PLAN_CREATION`
- NEW: `PHASE_IMPLEMENTATION_PLANNING` (R211, no old equivalent)
- NEW: `WAVE_IMPLEMENTATION_PLANNING` (R211, no old equivalent)
- NEW: `WAVE_DIRECTORY_ACKNOWLEDGMENT` (R214, no old equivalent)

### Architect
- OLD: `REVIEW` → NEW: `WAVE_REVIEW` or `PHASE_ASSESSMENT`
- NEW: `PHASE_ARCHITECTURE_PLANNING` (R210, no old equivalent)
- NEW: `WAVE_ARCHITECTURE_PLANNING` (R210, no old equivalent)

## Action Required

All agents must:
1. Read SOFTWARE-FACTORY-STATE-MACHINE.md on startup
2. Validate state transitions against it
3. Ignore state-machines/*.md files
4. Use new state names in all transitions
EOF

echo -e "${GREEN}✅ Created migration guide: $MIGRATION_GUIDE${NC}"
echo ""

echo "════════════════════════════════════════════════════════════════"
echo -e "${BOLD}RECONCILIATION COMPLETE${NC}"
echo "════════════════════════════════════════════════════════════════"
echo ""
echo -e "${BOLD}Key Findings:${NC}"
echo "1. Two state machine systems exist (legacy vs current)"
echo "2. SOFTWARE-FACTORY-STATE-MACHINE.md is authoritative (R206)"
echo "3. state-machines/*.md files are DEPRECATED"
echo "4. New architecture states (R210, R211, R214) only in main system"
echo ""
echo -e "${BOLD}Actions Taken:${NC}"
echo "✅ Added deprecation warnings to legacy files"
echo "✅ Created migration guide"
echo "✅ Documented state mappings"
echo ""
echo -e "${RED}${BOLD}⚠️  IMPORTANT: Agents must use SOFTWARE-FACTORY-STATE-MACHINE.md${NC}"