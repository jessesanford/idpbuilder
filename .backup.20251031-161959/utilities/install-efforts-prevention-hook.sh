#!/bin/bash

# Install/Update git pre-commit hook to prevent efforts/ directory commits
# This ensures implementation work goes to the target repository, not the planning repo
#
# Purpose: Prevents agents from accidentally committing implementation work
#          to the Software Factory planning repository. All effort work must
#          be pushed to the target repository defined in target-repo-config.yaml
#
# Usage: ./utilities/install-efforts-prevention-hook.sh

set -euo pipefail

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
BOLD='\033[1m'
NC='\033[0m' # No Color

print_error() {
    echo -e "${RED}${BOLD}❌ ERROR:${NC} ${RED}$1${NC}" >&2
}

print_success() {
    echo -e "${GREEN}${BOLD}✅ PROJECT_DONE:${NC} ${GREEN}$1${NC}"
}

print_info() {
    echo -e "${BLUE}${BOLD}ℹ️  INFO:${NC} $1"
}

# Get the git root directory
GIT_ROOT="$(git rev-parse --show-toplevel 2>/dev/null)" || {
    print_error "Not in a git repository!"
    exit 1
}

cd "$GIT_ROOT"

# Check if pre-commit hook exists
HOOK_FILE=".git/hooks/pre-commit"
BACKUP_FILE=".git/hooks/pre-commit.backup-$(date +%Y%m%d-%H%M%S)"

if [ -f "$HOOK_FILE" ]; then
    print_info "Existing pre-commit hook found. Creating backup..."
    cp "$HOOK_FILE" "$BACKUP_FILE"
    print_success "Backup created: $BACKUP_FILE"

    # Check if the efforts check already exists
    if grep -q "efforts/ directory files - CRITICAL: This is the planning repo!" "$HOOK_FILE"; then
        print_info "Efforts directory check already exists in pre-commit hook"
        exit 0
    fi

    # Add the efforts check to existing hook
    print_info "Adding efforts directory check to existing pre-commit hook..."

    # Find the line before "Exit with appropriate code" and insert our check
    sed -i '/# Exit with appropriate code/i\
    # Check for efforts\/ directory files - CRITICAL: This is the planning repo!\
    # This check prevents agents from accidentally committing implementation work\
    # to the planning repository. All effort work must go to the target repository\
    # specified in target-repo-config.yaml\
    efforts_files=$(git diff --cached --name-only | grep "^efforts\/" || true)\
    if [ -n "$efforts_files" ]; then\
        echo ""\
        print_error "WRONG! THIS IS THE PLANNING REPO. YOU ARE NOT ALLOWED TO COMMIT EFFORT WORK TO THE PLANNING REPO. REMOVE THESE COMMITS AND INSTEAD PUSH THEM TO THE TARGET REPO DENOTED IN THE \\$CLAUDE_PROJECT_DIR\/target-repo-config.yaml"\
        echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"\
        echo ""\
        print_info "The following files from efforts\/ directory were detected:"\
        echo "$efforts_files" | while read -r file; do\
            echo "  - $file"\
        done\
        echo ""\
        print_info "This is the Software Factory planning repository."\
        print_info "Implementation work must be committed to the target repository."\
        print_info "Check target-repo-config.yaml for the correct repository."\
        echo ""\
        print_info "To remove these files from staging:"\
        echo -e "  ${BOLD}git reset HEAD efforts\/${NC}"\
        echo ""\
        validation_failed=true\
    fi\
' "$HOOK_FILE"

else
    print_info "Creating new pre-commit hook with efforts directory check..."

    cat > "$HOOK_FILE" << 'EOF'
#!/bin/bash

# Pre-commit hook for Software Factory 2.0
# Prevents commits of efforts/ directory files to the planning repository

set -euo pipefail

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
BOLD='\033[1m'
NC='\033[0m' # No Color

print_error() {
    echo -e "${RED}${BOLD}❌ ERROR:${NC} ${RED}$1${NC}" >&2
}

print_info() {
    echo -e "${BLUE}${BOLD}ℹ️  INFO:${NC} $1"
}

# Check for efforts/ directory files - CRITICAL: This is the planning repo!
# This check prevents agents from accidentally committing implementation work
# to the planning repository. All effort work must go to the target repository
# specified in target-repo-config.yaml
efforts_files=$(git diff --cached --name-only | grep "^efforts/" || true)
if [ -n "$efforts_files" ]; then
    echo ""
    print_error "WRONG! THIS IS THE PLANNING REPO. YOU ARE NOT ALLOWED TO COMMIT EFFORT WORK TO THE PLANNING REPO. REMOVE THESE COMMITS AND INSTEAD PUSH THEM TO THE TARGET REPO DENOTED IN THE \$CLAUDE_PROJECT_DIR/target-repo-config.yaml"
    echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    print_info "The following files from efforts/ directory were detected:"
    echo "$efforts_files" | while read -r file; do
        echo "  - $file"
    done
    echo ""
    print_info "This is the Software Factory planning repository."
    print_info "Implementation work must be committed to the target repository."
    print_info "Check target-repo-config.yaml for the correct repository."
    echo ""
    print_info "To remove these files from staging:"
    echo -e "  ${BOLD}git reset HEAD efforts/${NC}"
    echo ""
    exit 1
fi

exit 0
EOF
fi

# Make the hook executable
chmod +x "$HOOK_FILE"
print_success "Pre-commit hook has been updated with efforts directory prevention!"

# Test if .gitignore includes efforts/
if ! grep -q "^efforts/" .gitignore 2>/dev/null; then
    print_info "Adding efforts/ to .gitignore..."
    echo "efforts/" >> .gitignore
    print_success "Added efforts/ to .gitignore"
else
    print_info "efforts/ already in .gitignore"
fi

print_success "Installation complete!"
print_info "The pre-commit hook will now prevent any files in efforts/ from being committed to this repository."
print_info "Implementation work should go to: $(grep -A1 "repository:" target-repo-config.yaml 2>/dev/null | grep "url:" | cut -d'"' -f2 || echo "[Check target-repo-config.yaml]")"