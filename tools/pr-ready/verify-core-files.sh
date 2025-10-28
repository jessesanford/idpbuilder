#!/bin/bash
# verify-core-files.sh - Verify core application files are preserved

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Core files that must never be deleted
CORE_FILES_PATTERNS=(
  "main.*"           # main.go, main.py, main.rs, main.js, etc.
  "Makefile"
  "README*"
  "LICENSE*"
  "package.json"
  "package-lock.json"
  "go.mod"
  "go.sum"
  "Cargo.toml"
  "Cargo.lock"
  "requirements.txt"
  "setup.py"
  "pyproject.toml"
  ".gitignore"
  "Dockerfile"
  "docker-compose.yml"
)

# Core directories that should exist
CORE_DIRS=(
  "src"
  "pkg"
  "lib"
  "app"
  "api"
  "cmd"
  "internal"
  "public"
  "static"
  "assets"
)

# Function to check core files in a branch
check_branch() {
  local branch=$1
  local missing_files=()
  local found_files=()

  echo -e "${YELLOW}Checking branch: $branch${NC}"
  git checkout $branch 2>/dev/null

  # Check for core files
  for pattern in "${CORE_FILES_PATTERNS[@]}"; do
    local files=$(ls $pattern 2>/dev/null || true)
    if [ -z "$files" ]; then
      # Only report if file exists in main
      git checkout main 2>/dev/null
      if [ -n "$(ls $pattern 2>/dev/null || true)" ]; then
        missing_files+=("$pattern")
      fi
      git checkout $branch 2>/dev/null
    else
      found_files+=("$files")
    fi
  done

  # Check for core directories
  for dir in "${CORE_DIRS[@]}"; do
    if [ -d "$dir" ]; then
      echo -e "${GREEN}  ✓ Directory exists: $dir/${NC}"
    fi
  done

  # Report missing files
  if [ ${#missing_files[@]} -gt 0 ]; then
    echo -e "${RED}  ✗ CRITICAL: Missing core files:${NC}"
    for file in "${missing_files[@]}"; do
      echo -e "${RED}    - $file${NC}"
    done
    return 1
  fi

  # Check for mass deletions
  local deletions=$(git diff main --numstat | awk '$2 > 1000 {print $3, "-"$2" lines"}')
  if [ ! -z "$deletions" ]; then
    echo -e "${RED}  ⚠ WARNING: Large deletions detected:${NC}"
    echo "$deletions" | while read line; do
      echo -e "${YELLOW}    $line${NC}"
    done
  fi

  echo -e "${GREEN}  ✓ Core files intact${NC}"
  return 0
}

# Function to detect destructive changes
detect_destructive_changes() {
  local branch=$1

  git checkout $branch 2>/dev/null

  # Check for massive deletions (>10,000 lines)
  local total_deletions=$(git diff main --shortstat | grep -oE '[0-9]+ deletion' | grep -oE '[0-9]+' || echo 0)

  if [ "$total_deletions" -gt 10000 ]; then
    echo -e "${RED}⚠ CRITICAL: Branch has $total_deletions deletions!${NC}"
    echo -e "${RED}  This branch may have accidentally deleted core functionality${NC}"

    # Show what was deleted
    echo -e "${YELLOW}  Top deleted files:${NC}"
    git diff main --numstat | sort -rn -k2 | head -10 | while read added deleted file; do
      if [ "$deleted" -gt 100 ]; then
        echo -e "${RED}    $file: -$deleted lines${NC}"
      fi
    done

    return 1
  fi

  return 0
}

# Main execution
main() {
  echo "========================================="
  echo "Core File Integrity Verification Tool"
  echo "========================================="
  echo

  # Get branch name from argument or check all
  if [ $# -eq 1 ]; then
    BRANCHES=$1
  else
    # Get all branches except main
    BRANCHES=$(git branch -a | grep -v main | grep -E "(phase|wave|effort|feature)" | sed 's/^[* ]*//' | sed 's/remotes\/origin\///' | sort -u)
  fi

  local failed_branches=()
  local destructive_branches=()

  # Save current branch
  ORIGINAL_BRANCH=$(git branch --show-current)

  # Check each branch
  for branch in $BRANCHES; do
    if ! check_branch $branch; then
      failed_branches+=("$branch")
    fi

    if ! detect_destructive_changes $branch; then
      destructive_branches+=("$branch")
    fi
    echo
  done

  # Return to original branch
  git checkout $ORIGINAL_BRANCH 2>/dev/null

  # Summary
  echo "========================================="
  echo "SUMMARY"
  echo "========================================="

  if [ ${#failed_branches[@]} -gt 0 ]; then
    echo -e "${RED}Branches with missing core files:${NC}"
    for branch in "${failed_branches[@]}"; do
      echo -e "${RED}  - $branch${NC}"
    done
  fi

  if [ ${#destructive_branches[@]} -gt 0 ]; then
    echo -e "${RED}Branches with destructive changes:${NC}"
    for branch in "${destructive_branches[@]}"; do
      echo -e "${RED}  - $branch${NC}"
    done
  fi

  if [ ${#failed_branches[@]} -eq 0 ] && [ ${#destructive_branches[@]} -eq 0 ]; then
    echo -e "${GREEN}All branches passed integrity check!${NC}"
    exit 0
  fi

  exit 1
}

# Run main function
main "$@"