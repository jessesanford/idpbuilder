#!/bin/bash
# detect-sf-artifacts.sh - Detect Software Factory artifacts in branches

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Software Factory artifact patterns
SF_DIRS=(
  "todos"
  "efforts"
  "agent-states"
  "rule-library"
  "templates"
  "utilities"
  "phase-plans"
  "wave-plans"
  "protocols"
  ".claude/agents"
  ".claude/commands"
)

SF_FILE_PATTERNS=(
  "*-state.json"
  "*.todo"
  "CODE-REVIEW-REPORT.md"
  "SPLIT-PLAN*.md"
  "PROJECT-IMPLEMENTATION-PLAN.md"
  "EFFORT-IMPLEMENTATION-PLAN.md"
  "software-factory-3.0-state-machine.json"
  "RECOVERY-*.md"
  "CURRENT-TODO-STATE.md"
  "CRITICAL-*.md"
  "FINAL-*.md"
  "TODO-*.md"
)

# Function to scan a branch
scan_branch() {
  local branch=$1
  local artifact_count=0

  echo -e "${YELLOW}Scanning branch: $branch${NC}"
  git checkout $branch 2>/dev/null

  # Check for SF directories
  for dir in "${SF_DIRS[@]}"; do
    if [ -d "$dir" ]; then
      local file_count=$(find "$dir" -type f 2>/dev/null | wc -l)
      if [ $file_count -gt 0 ]; then
        echo -e "${RED}  ✗ Found: $dir/ ($file_count files)${NC}"
        artifact_count=$((artifact_count + file_count))
      fi
    fi
  done

  # Check for SF files
  for pattern in "${SF_FILE_PATTERNS[@]}"; do
    local files=$(ls $pattern 2>/dev/null || true)
    if [ ! -z "$files" ]; then
      local file_count=$(echo "$files" | wc -w)
      echo -e "${RED}  ✗ Found: $pattern ($file_count files)${NC}"
      artifact_count=$((artifact_count + file_count))
    fi
  done

  if [ $artifact_count -eq 0 ]; then
    echo -e "${GREEN}  ✓ Clean - no SF artifacts found${NC}"
  else
    echo -e "${RED}  Total artifacts: $artifact_count${NC}"
  fi

  return $artifact_count
}

# Main execution
main() {
  echo "========================================="
  echo "Software Factory Artifact Detection Tool"
  echo "========================================="
  echo

  # Get branch name from argument or scan all
  if [ $# -eq 1 ]; then
    BRANCHES=$1
  else
    # Get all branches except main
    BRANCHES=$(git branch -a | grep -v main | grep -E "(phase|wave|effort|feature)" | sed 's/^[* ]*//' | sed 's/remotes\/origin\///' | sort -u)
  fi

  local total_artifacts=0
  local contaminated_branches=0
  local clean_branches=0

  # Save current branch
  ORIGINAL_BRANCH=$(git branch --show-current)

  # Scan each branch
  for branch in $BRANCHES; do
    scan_branch $branch
    local branch_artifacts=$?

    if [ $branch_artifacts -gt 0 ]; then
      contaminated_branches=$((contaminated_branches + 1))
      total_artifacts=$((total_artifacts + branch_artifacts))
    else
      clean_branches=$((clean_branches + 1))
    fi
    echo
  done

  # Return to original branch
  git checkout $ORIGINAL_BRANCH 2>/dev/null

  # Summary
  echo "========================================="
  echo "SUMMARY"
  echo "========================================="
  echo -e "Contaminated branches: ${RED}$contaminated_branches${NC}"
  echo -e "Clean branches: ${GREEN}$clean_branches${NC}"
  echo -e "Total artifacts found: ${YELLOW}$total_artifacts${NC}"

  # Exit with error if artifacts found
  if [ $total_artifacts -gt 0 ]; then
    exit 1
  fi

  echo -e "${GREEN}All branches are clean!${NC}"
  exit 0
}

# Run main function
main "$@"