#!/bin/bash
# R505: Pre-Planned Infrastructure Synchronization Validation Script
# Purpose: Validate that pre_planned_infrastructure is synchronized and correct

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Get project directory
PROJECT_DIR="${CLAUDE_PROJECT_DIR:-$(pwd)}"
cd "$PROJECT_DIR"

echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}R505: Pre-Planned Infrastructure Synchronization Check${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"

# Check if orchestrator-state-v3.json exists
if [ ! -f "orchestrator-state-v3.json" ]; then
    echo -e "${RED}❌ ERROR: orchestrator-state-v3.json not found${NC}"
    exit 1
fi

# Check if target-repo-config.yaml exists
if [ ! -f "target-repo-config.yaml" ]; then
    echo -e "${RED}❌ ERROR: target-repo-config.yaml not found${NC}"
    exit 1
fi

# Get current phase and wave
CURRENT_PHASE=$(jq -r '.current_phase' orchestrator-state-v3.json)
CURRENT_WAVE=$(jq -r '.current_wave' orchestrator-state-v3.json)
TARGET_REPO=$(yq '.repository' target-repo-config.yaml)
PROJECT_PREFIX=$(jq -r '.project_prefix // "project"' orchestrator-state-v3.json)

echo -e "${BLUE}Current Phase: ${NC}$CURRENT_PHASE"
echo -e "${BLUE}Current Wave: ${NC}$CURRENT_WAVE"
echo -e "${BLUE}Target Repository: ${NC}$TARGET_REPO"
echo -e "${BLUE}Project Prefix: ${NC}$PROJECT_PREFIX"
echo ""

# Function to check infrastructure entry
check_infrastructure_entry() {
    local EFFORT_ID=$1
    local ENTRY=$(jq ".pre_planned_infrastructure.efforts.\"$EFFORT_ID\"" orchestrator-state-v3.json)

    if [ "$ENTRY" == "null" ]; then
        echo -e "${RED}❌ Missing: $EFFORT_ID${NC}"
        return 1
    fi

    # Check required fields
    local FULL_PATH=$(echo "$ENTRY" | jq -r '.full_path')
    local BRANCH_NAME=$(echo "$ENTRY" | jq -r '.branch_name')
    local BASE_BRANCH=$(echo "$ENTRY" | jq -r '.base_branch')
    local TARGET_REPOSITORY=$(echo "$ENTRY" | jq -r '.target_repository')

    local ERRORS=0

    # Validate path
    if [ -z "$FULL_PATH" ] || [ "$FULL_PATH" == "null" ]; then
        echo -e "${RED}  ❌ Missing full_path for $EFFORT_ID${NC}"
        ((ERRORS++))
    fi

    # Validate branch name
    if [ -z "$BRANCH_NAME" ] || [ "$BRANCH_NAME" == "null" ]; then
        echo -e "${RED}  ❌ Missing branch_name for $EFFORT_ID${NC}"
        ((ERRORS++))
    fi

    # Validate base branch
    if [ -z "$BASE_BRANCH" ] || [ "$BASE_BRANCH" == "null" ]; then
        echo -e "${RED}  ❌ Missing base_branch for $EFFORT_ID${NC}"
        ((ERRORS++))
    fi

    # Validate target repository
    if [ "$TARGET_REPOSITORY" != "$TARGET_REPO" ]; then
        echo -e "${RED}  ❌ Wrong repository for $EFFORT_ID: $TARGET_REPOSITORY (expected: $TARGET_REPO)${NC}"
        ((ERRORS++))
    fi

    if [ $ERRORS -eq 0 ]; then
        echo -e "${GREEN}✅ Valid: $EFFORT_ID${NC}"
        echo -e "   Path: $FULL_PATH"
        echo -e "   Branch: $BRANCH_NAME"
        echo -e "   Base: $BASE_BRANCH"
        return 0
    else
        return 1
    fi
}

# Function to check cascade pattern (R308)
check_cascade_pattern() {
    local PHASE=$1
    local WAVE=$2

    echo -e "\n${BLUE}Checking R308 Cascade Pattern for Wave ${PHASE}.${WAVE}...${NC}"

    local EFFORTS=$(jq -r ".pre_planned_infrastructure.efforts | keys | .[] | select(. | test(\"phase${PHASE}_wave${WAVE}_\"))" orchestrator-state-v3.json | sort)

    if [ -z "$EFFORTS" ]; then
        echo -e "${YELLOW}⚠️  No efforts found for Wave ${PHASE}.${WAVE}${NC}"
        return 1
    fi

    local PREV_BRANCH=""
    local EFFORT_INDEX=0

    for EFFORT_ID in $EFFORTS; do
        ((EFFORT_INDEX++))

        local BASE_BRANCH=$(jq -r ".pre_planned_infrastructure.efforts.\"$EFFORT_ID\".base_branch" orchestrator-state-v3.json)
        local BRANCH_NAME=$(jq -r ".pre_planned_infrastructure.efforts.\"$EFFORT_ID\".branch_name" orchestrator-state-v3.json)

        # Check cascade pattern
        if [[ $PHASE -eq 1 && $WAVE -eq 1 && $EFFORT_INDEX -eq 1 ]]; then
            # First effort of P1W1 should be from main
            if [ "$BASE_BRANCH" != "main" ]; then
                echo -e "${RED}  ❌ First effort of P1W1 should base on main, got: $BASE_BRANCH${NC}"
            else
                echo -e "${GREEN}  ✅ $EFFORT_ID cascades from main (correct for first effort)${NC}"
            fi
        elif [[ $EFFORT_INDEX -eq 1 ]]; then
            # First effort of wave should cascade from previous wave/phase
            echo -e "${BLUE}  ℹ️  $EFFORT_ID (first of wave) bases on: $BASE_BRANCH${NC}"
        else
            # Subsequent efforts should cascade from previous effort
            if [ "$BASE_BRANCH" != "$PREV_BRANCH" ]; then
                echo -e "${YELLOW}  ⚠️  $EFFORT_ID may not follow cascade (base: $BASE_BRANCH, expected: $PREV_BRANCH)${NC}"
            else
                echo -e "${GREEN}  ✅ $EFFORT_ID cascades from $BASE_BRANCH${NC}"
            fi
        fi

        PREV_BRANCH="$BRANCH_NAME"
    done
}

# Function to check synchronization timestamps
check_sync_timestamps() {
    echo -e "\n${BLUE}Checking Synchronization Timestamps...${NC}"

    local LAST_PHASE_SYNC=$(jq -r '.pre_planned_infrastructure.last_phase_sync // "never"' orchestrator-state-v3.json)
    local LAST_WAVE_SYNC=$(jq -r '.pre_planned_infrastructure.last_wave_sync // "never"' orchestrator-state-v3.json)
    local VALIDATED=$(jq -r '.pre_planned_infrastructure.validated // false' orchestrator-state-v3.json)

    echo -e "Last Phase Sync: $LAST_PHASE_SYNC"
    echo -e "Last Wave Sync: $LAST_WAVE_SYNC"
    echo -e "Validated: $VALIDATED"

    if [ "$LAST_PHASE_SYNC" == "never" ] && [ "$LAST_WAVE_SYNC" == "never" ]; then
        echo -e "${YELLOW}⚠️  Infrastructure has never been synchronized${NC}"
    fi

    if [ "$VALIDATED" == "false" ]; then
        echo -e "${YELLOW}⚠️  Infrastructure needs validation${NC}"
    fi
}

# Main validation
echo -e "\n${BLUE}Validating Current Wave Infrastructure...${NC}"

# Get all efforts for current wave
WAVE_EFFORTS=$(jq -r ".pre_planned_infrastructure.efforts | keys | .[] | select(. | test(\"phase${CURRENT_PHASE}_wave${CURRENT_WAVE}_\"))" orchestrator-state-v3.json)

if [ -z "$WAVE_EFFORTS" ]; then
    echo -e "${RED}❌ No infrastructure entries found for current wave!${NC}"
    echo -e "${RED}   Run R505 synchronization after phase/wave planning${NC}"
    exit 505
fi

# Check each effort
TOTAL_EFFORTS=0
VALID_EFFORTS=0

for EFFORT_ID in $WAVE_EFFORTS; do
    ((TOTAL_EFFORTS++))
    if check_infrastructure_entry "$EFFORT_ID"; then
        ((VALID_EFFORTS++))
    fi
done

echo -e "\n${BLUE}Summary:${NC}"
echo -e "Total Efforts: $TOTAL_EFFORTS"
echo -e "Valid Efforts: $VALID_EFFORTS"

# Check cascade pattern
check_cascade_pattern "$CURRENT_PHASE" "$CURRENT_WAVE"

# Check timestamps
check_sync_timestamps

# Final result
echo -e "\n${BLUE}═══════════════════════════════════════════════════════${NC}"

if [ $VALID_EFFORTS -eq $TOTAL_EFFORTS ] && [ $TOTAL_EFFORTS -gt 0 ]; then
    echo -e "${GREEN}✅ PASS: Infrastructure synchronization is valid${NC}"

    # Mark as validated
    jq '.pre_planned_infrastructure.validated = true |
        .pre_planned_infrastructure.validation_timestamp = (now | todate)' \
       orchestrator-state-v3.json > orchestrator-state.tmp && \
       mv orchestrator-state.tmp orchestrator-state-v3.json

    echo -e "${GREEN}Infrastructure marked as validated${NC}"
    exit 0
else
    echo -e "${RED}❌ FAIL: Infrastructure synchronization has issues${NC}"
    echo -e "${RED}Please run synchronization or fix the issues above${NC}"
    exit 505
fi