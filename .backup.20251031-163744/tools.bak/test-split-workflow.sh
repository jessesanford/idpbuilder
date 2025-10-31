#!/bin/bash
# Test script to demonstrate the split planning workflow

set -euo pipefail

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}Testing Split Plan Workflow${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"
echo ""

# Step 1: Simulate Code Reviewer creating split plans
echo -e "${YELLOW}Step 1: Code Reviewer Creates Split Plans${NC}"
echo "In the too-large branch directory:"
echo ""
cat << 'EOF'
# Code Reviewer would execute:
cd /efforts/phase1/wave1/api-types  # Too-large branch

# Create SPLIT-INVENTORY.md
cat > SPLIT-INVENTORY.md << 'INVENTORY'
# Split Inventory for api-types

## Overview
api-types effort exceeded 800 lines (actual: 1247) and requires splitting.

## Split Structure
| Split # | Name | Description | Est. Lines | Status |
|---------|------|-------------|------------|--------|
| 001 | core-types | Core API type definitions | 420 | Planned |
| 002 | validators | Validation logic | 380 | Planned |
| 003 | converters | Type converters and helpers | 400 | Planned |
INVENTORY

# Create individual split plans
for num in 001 002 003; do
    cat > SPLIT-PLAN-${num}.md << 'PLAN'
# Split Plan ${num}

## Split Metadata
- Split Number: ${num}
- Parent Effort: api-types
- Target Size: <800 lines

## Implementation Scope
[Details here...]
PLAN
done

# Commit and push
git add SPLIT-*.md
git commit -m "feat: add split plans for oversized api-types effort"
git push
EOF

echo ""
echo -e "${GREEN}✓ Split plans created and pushed to too-large branch${NC}"
echo ""

# Step 2: Show what Orchestrator would do
echo -e "${YELLOW}Step 2: Orchestrator Creates Split Infrastructure${NC}"
echo "The orchestrator would execute:"
echo ""
cat << 'EOF'
# Orchestrator reads split inventory
cd /efforts/phase1/wave1/api-types
git pull  # Get latest split plans

# Read number of splits
TOTAL_SPLITS=$(grep -c "^| [0-9]" SPLIT-INVENTORY.md)
echo "Found $TOTAL_SPLITS splits to create"

# For each split, create infrastructure
for split_num in $(seq 1 $TOTAL_SPLITS); do
    SPLIT_NAME=$(printf "%03d" $split_num)
    SPLIT_DIR="/efforts/phase1/wave1/api-types--split-${SPLIT_NAME}"
    
    echo "Creating split-${SPLIT_NAME} infrastructure..."
    
    # Create directory
    mkdir -p "$SPLIT_DIR"
    
    # Clone and setup branch
    git clone --sparse [repo] "$SPLIT_DIR"
    cd "$SPLIT_DIR"
    git checkout -b "tmc-workspace/phase1/wave1/api-types--split-${SPLIT_NAME}"
    
    # Copy ONLY the specific split plan from too-large branch
    cp ../api-types/SPLIT-PLAN-${SPLIT_NAME}.md .
    
    # Commit initial setup
    git add -A
    git commit -m "chore: initialize split-${SPLIT_NAME}"
    git push -u origin [branch]
done
EOF

echo ""
echo -e "${GREEN}✓ All split directories created with plans${NC}"
echo ""

# Step 3: Show final structure
echo -e "${YELLOW}Step 3: Final Directory Structure${NC}"
cat << 'EOF'
efforts/phase1/wave1/
├── api-types/                    # Too-large branch (to be abandoned)
│   ├── SPLIT-INVENTORY.md        # Created by Code Reviewer (stays here)
│   ├── SPLIT-PLAN-001.md         # Created by Code Reviewer
│   ├── SPLIT-PLAN-002.md         # Created by Code Reviewer
│   └── SPLIT-PLAN-003.md         # Created by Code Reviewer
│
├── api-types--split-001/         # Created by Orchestrator
│   └── SPLIT-PLAN-001.md         # Copied from too-large (only this file)
│
├── api-types--split-002/         # Created by Orchestrator
│   └── SPLIT-PLAN-002.md         # Copied from too-large (only this file)
│
└── api-types--split-003/         # Created by Orchestrator
    └── SPLIT-PLAN-003.md         # Copied from too-large (only this file)
EOF

echo ""
echo -e "${YELLOW}Step 4: SW Engineer Implementation${NC}"
echo "Each SW Engineer receives their split directory with:"
echo "  • Their specific SPLIT-PLAN-XXX.md (only file needed)"
echo "  • Clean workspace ready for implementation"
echo "  • No conflicting files that would cause merge issues"
echo ""

echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}✅ Workflow Test Complete${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"
echo ""
echo "Key Benefits of This Approach:"
echo "  1. Split plans are version-controlled in too-large branch"
echo "  2. Clear indication that branch was split (has SPLIT-*.md files)"
echo "  3. SW Engineers get clear, specific instructions"
echo "  4. Orchestrator can verify plans exist before creating infrastructure"
echo "  5. Process is auditable and reversible"