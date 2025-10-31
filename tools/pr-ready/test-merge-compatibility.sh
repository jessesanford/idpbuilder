#!/bin/bash
# test-merge-compatibility.sh - Test if branches can merge cleanly

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to test merge for a branch
test_merge() {
  local branch=$1
  local target=${2:-main}
  local conflicts=()

  echo -e "${YELLOW}Testing merge: $branch → $target${NC}"

  # Create test branch
  git checkout $target 2>/dev/null
  git checkout -b test-merge-$branch 2>/dev/null

  # Try merge
  if git merge --no-commit --no-ff origin/$branch 2>/dev/null; then
    echo -e "${GREEN}  ✓ Merge successful - no conflicts${NC}"

    # Check if build still works
    if [ -f "Makefile" ]; then
      echo -e "${BLUE}  Testing build...${NC}"
      if make build > /dev/null 2>&1; then
        echo -e "${GREEN}  ✓ Build successful${NC}"
      else
        echo -e "${RED}  ✗ Build failed after merge${NC}"
      fi
    fi

    git merge --abort 2>/dev/null || true
  else
    # Get conflict list
    local conflict_files=$(git status --porcelain | grep "^UU\|^AA" | awk '{print $2}')
    echo -e "${RED}  ✗ Merge has conflicts:${NC}"
    echo "$conflict_files" | while read file; do
      echo -e "${RED}    - $file${NC}"
      conflicts+=("$file")
    done

    git merge --abort 2>/dev/null || true
  fi

  # Clean up test branch
  git checkout main 2>/dev/null
  git branch -D test-merge-$branch 2>/dev/null || true

  return ${#conflicts[@]}
}

# Function to test sequential merge order
test_sequential_merges() {
  local branches=("$@")
  local merged_branch="main"

  echo -e "${BLUE}Testing sequential merge order...${NC}"
  echo

  # Create test integration branch
  git checkout main 2>/dev/null
  git checkout -b test-integration 2>/dev/null

  for branch in "${branches[@]}"; do
    echo -e "${YELLOW}Merging: $branch${NC}"

    if git merge --no-edit origin/$branch 2>/dev/null; then
      echo -e "${GREEN}  ✓ Merged successfully${NC}"
      merged_branch=$branch
    else
      # Document conflicts
      echo -e "${RED}  ✗ Conflicts when merging $branch after $merged_branch${NC}"
      local conflict_files=$(git status --porcelain | grep "^UU\|^AA" | awk '{print $2}')
      echo "$conflict_files" | while read file; do
        echo -e "${RED}    Conflict in: $file${NC}"
      done

      # Try automatic resolution for known patterns
      echo -e "${YELLOW}  Attempting automatic resolution...${NC}"

      # Common resolution: accept incoming for add/add conflicts
      for file in $conflict_files; do
        if [ -f "$file" ]; then
          git checkout --theirs "$file"
          git add "$file"
        fi
      done

      if git commit -m "Merge $branch with auto-resolved conflicts" 2>/dev/null; then
        echo -e "${GREEN}  ✓ Auto-resolution successful${NC}"
      else
        echo -e "${RED}  ✗ Auto-resolution failed${NC}"
        git merge --abort 2>/dev/null || true
        break
      fi
    fi
    echo
  done

  # Test final build
  echo -e "${BLUE}Testing final build...${NC}"
  if [ -f "Makefile" ]; then
    if make build > /dev/null 2>&1; then
      echo -e "${GREEN}✓ Final build successful${NC}"
    else
      echo -e "${RED}✗ Final build failed${NC}"
    fi
  fi

  # Clean up
  git checkout main 2>/dev/null
  git branch -D test-integration 2>/dev/null || true
}

# Function to generate conflict resolution guide
generate_resolution_guide() {
  local branches=("$@")

  echo "# Merge Conflict Resolution Guide" > PR-MERGE-GUIDE.md
  echo "" >> PR-MERGE-GUIDE.md
  echo "## Merge Order" >> PR-MERGE-GUIDE.md
  echo "" >> PR-MERGE-GUIDE.md

  local order=1
  for branch in "${branches[@]}"; do
    echo "$order. **$branch**" >> PR-MERGE-GUIDE.md

    # Test merge and document any conflicts
    git checkout main 2>/dev/null
    git checkout -b test-guide-$branch 2>/dev/null

    if ! git merge --no-commit --no-ff origin/$branch 2>/dev/null; then
      echo "   - Expected conflicts in:" >> PR-MERGE-GUIDE.md
      git status --porcelain | grep "^UU\|^AA" | awk '{print $2}' | while read file; do
        echo "     - \`$file\`" >> PR-MERGE-GUIDE.md
      done
      echo "   - Resolution: \`git checkout --theirs $file\` (accept incoming)" >> PR-MERGE-GUIDE.md
      git merge --abort 2>/dev/null || true
    else
      echo "   - No conflicts expected" >> PR-MERGE-GUIDE.md
      git merge --abort 2>/dev/null || true
    fi

    git checkout main 2>/dev/null
    git branch -D test-guide-$branch 2>/dev/null || true

    echo "" >> PR-MERGE-GUIDE.md
    order=$((order + 1))
  done

  echo -e "${GREEN}✓ Generated PR-MERGE-GUIDE.md${NC}"
}

# Main execution
main() {
  echo "========================================="
  echo "Merge Compatibility Testing Tool"
  echo "========================================="
  echo

  # Get branches from argument or detect
  if [ $# -gt 0 ]; then
    BRANCHES=("$@")
  else
    # Get branches in dependency order
    BRANCHES=($(git branch -a | grep -v main | grep -E "(phase|wave|effort|feature)" | sed 's/^[* ]*//' | sed 's/remotes\/origin\///' | sort))
  fi

  # Save current branch
  ORIGINAL_BRANCH=$(git branch --show-current)

  # Test individual merges
  echo -e "${BLUE}=== Individual Merge Tests ===${NC}"
  echo
  for branch in "${BRANCHES[@]}"; do
    test_merge $branch
    echo
  done

  # Test sequential merges
  echo -e "${BLUE}=== Sequential Merge Test ===${NC}"
  echo
  test_sequential_merges "${BRANCHES[@]}"

  # Generate resolution guide
  echo
  echo -e "${BLUE}=== Generating Resolution Guide ===${NC}"
  generate_resolution_guide "${BRANCHES[@]}"

  # Return to original branch
  git checkout $ORIGINAL_BRANCH 2>/dev/null

  echo
  echo "========================================="
  echo -e "${GREEN}Testing complete!${NC}"
  echo "See PR-MERGE-GUIDE.md for resolution instructions"
  echo "========================================="
}

# Run main function
main "$@"