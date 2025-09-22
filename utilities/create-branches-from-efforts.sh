#!/bin/bash
# Creates git branches for effort directories that don't have branches yet

set -e

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
BOLD='\033[1m'
RESET='\033[0m'

# Print colored output
log() {
    local level="$1"
    shift
    echo -e "${level}$*${RESET}"
}

# Find all effort directories
find_effort_directories() {
    find efforts -type d -name "phase*" | while read phase_dir; do
        find "$phase_dir" -type d -name "wave*" | while read wave_dir; do
            find "$wave_dir" -maxdepth 1 -type d ! -name "wave*" ! -name "integration-workspace" | while read effort_dir; do
                if [ "$effort_dir" != "$wave_dir" ]; then
                    echo "$effort_dir"
                fi
            done
        done
    done
}

log "${BOLD}${CYAN}" "🏭 Creating branches for effort directories"
log "${BLUE}" "Repository: $(pwd)"

# Get current branch
ORIGINAL_BRANCH=$(git branch --show-current)
log "${BLUE}" "Original branch: $ORIGINAL_BRANCH"

# Find all effort directories
EFFORT_DIRS=$(find_effort_directories)

if [ -z "$EFFORT_DIRS" ]; then
    log "${YELLOW}" "No effort directories found"
    exit 0
fi

log "${CYAN}" "\nFound effort directories:"
echo "$EFFORT_DIRS" | while read dir; do
    log "${BLUE}" "  • $dir"
done

# Extract project name from current branch or use default
PROJECT_NAME=$(echo "$ORIGINAL_BRANCH" | cut -d'/' -f1)
if [ "$PROJECT_NAME" == "$ORIGINAL_BRANCH" ]; then
    # No slash in branch name, use repository name
    PROJECT_NAME=$(basename $(git rev-parse --show-toplevel))
fi

log "${CYAN}" "\nProject name: $PROJECT_NAME"

# Process each effort directory
echo "$EFFORT_DIRS" | while read effort_dir; do
    if [ -z "$effort_dir" ]; then
        continue
    fi
    
    # Extract phase, wave, and effort name
    # Example: efforts/phase2/wave1/image-builder -> phase2, wave1, image-builder
    phase=$(echo "$effort_dir" | cut -d'/' -f2)
    wave=$(echo "$effort_dir" | cut -d'/' -f3)
    effort=$(basename "$effort_dir")
    
    # Construct branch name
    branch_name="${PROJECT_NAME}/${phase}/${wave}/${effort}"
    
    log "${CYAN}" "\nProcessing: $effort_dir"
    log "${BLUE}" "  Branch name: $branch_name"
    
    # Check if branch already exists
    if git show-ref --verify --quiet "refs/heads/$branch_name"; then
        log "${YELLOW}" "  ⚠️  Branch already exists"
        continue
    fi
    
    # Create the branch
    log "${GREEN}" "  ✅ Creating branch: $branch_name"
    
    # Create branch from current branch
    git checkout -b "$branch_name" "$ORIGINAL_BRANCH" 2>/dev/null || {
        log "${RED}" "  ❌ Failed to create branch"
        continue
    }
    
    # Add the effort directory contents
    if [ -n "$(ls -A "$effort_dir")" ]; then
        git add "$effort_dir"
        
        # Check if there are changes to commit
        if ! git diff --cached --quiet; then
            git commit -m "feat($effort): add $phase $wave implementation

Effort: $effort
Phase: $phase
Wave: $wave
Directory: $effort_dir"
            log "${GREEN}" "  ✅ Committed effort files"
        else
            log "${YELLOW}" "  ℹ️  No new files to commit"
        fi
    else
        log "${YELLOW}" "  ℹ️  Directory is empty"
    fi
done

# Return to original branch
log "${CYAN}" "\nReturning to original branch: $ORIGINAL_BRANCH"
git checkout "$ORIGINAL_BRANCH"

log "${BOLD}${GREEN}" "\n✅ Branch creation complete!"
log "${CYAN}" "Run './utilities/push-local-efforts.sh' to push these branches to remote"