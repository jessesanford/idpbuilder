#!/bin/bash
# Software Factory 3.0 Initialization Script
# Creates a new SF 3.0 project from template
#
# Usage:
#   bash utilities/init-software-factory.sh --project-name my-project [--non-interactive]
#
# This script:
# 1. Creates 4 SF 3.0 state files from examples
# 2. Installs pre-commit hooks (R506 enforcement)
# 3. Creates README.md from template
# 4. Initializes git repository (if needed)

set -euo pipefail

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Script configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
PROJECT_NAME=""
NON_INTERACTIVE=0
LICENSE_INFO="See LICENSE file"

# Logging functions
log_info() {
    echo -e "${BLUE}${BOLD}ℹ INFO:${NC} ${BLUE}$*${NC}"
}

log_success() {
    echo -e "${GREEN}${BOLD}✅ PROJECT_DONE:${NC} ${GREEN}$*${NC}"
}

log_error() {
    echo -e "${RED}${BOLD}❌ ERROR:${NC} ${RED}$*${NC}" >&2
}

log_warning() {
    echo -e "${YELLOW}${BOLD}⚠ WARNING:${NC} ${YELLOW}$*${NC}"
}

# Usage function
usage() {
    cat <<EOF
Software Factory 3.0 Initialization Script

Usage:
    $0 --project-name <name> [OPTIONS]

Required Arguments:
    --project-name <name>     Name of the project to create

Optional Arguments:
    --non-interactive         Run without user prompts (for testing)
    --license <info>          License information (default: "See LICENSE file")
    --help                    Show this help message

Examples:
    # Interactive mode
    bash utilities/init-software-factory.sh --project-name my-web-app

    # Non-interactive mode (for CI/testing)
    bash utilities/init-software-factory.sh --project-name test-app --non-interactive

EOF
    exit 1
}

# Parse command line arguments
parse_args() {
    if [ $# -eq 0 ]; then
        usage
    fi

    while [[ $# -gt 0 ]]; do
        case $1 in
            --project-name)
                PROJECT_NAME="$2"
                shift 2
                ;;
            --non-interactive)
                NON_INTERACTIVE=1
                shift
                ;;
            --license)
                LICENSE_INFO="$2"
                shift 2
                ;;
            --help)
                usage
                ;;
            *)
                log_error "Unknown option: $1"
                usage
                ;;
        esac
    done

    if [ -z "$PROJECT_NAME" ]; then
        log_error "Project name is required"
        usage
    fi
}

# Validate project name
validate_project_name() {
    if [[ ! "$PROJECT_NAME" =~ ^[a-zA-Z0-9_-]+$ ]]; then
        log_error "Project name must contain only letters, numbers, hyphens, and underscores"
        exit 1
    fi
}

# Check if git repository exists
check_git_repo() {
    if [ ! -d "$PROJECT_ROOT/.git" ]; then
        log_warning "Not a git repository. Initializing git..."
        cd "$PROJECT_ROOT"
        git init
        log_success "Git repository initialized"
    else
        log_info "Git repository already initialized"
    fi
}

# Create state files from examples
create_state_files() {
    log_info "Creating SF 3.0 state files..."

    cd "$PROJECT_ROOT"

    # 1. Create orchestrator-state-v3.json
    if [ "$NON_INTERACTIVE" -eq 1 ] || [ ! -f "orchestrator-state-v3.json.example" ]; then
        # Create minimal state file from scratch for non-interactive/test runs
        cat > orchestrator-state-v3.json <<EOF
{
  "state_machine": {
    "current_state": "INIT",
    "previous_state": "",
    "state_history": [],
    "last_transition_timestamp": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
    "transition_validated": true,
    "status": "ACTIVE",
    "sub_state_machine": {
      "active": false,
      "type": "",
      "state_file": "",
      "return_state": ""
    },
    "loop_detection": {
      "last_two_states": [],
      "ping_pong_count": 0,
      "same_state_count": 0,
      "error_recovery_entries": 0,
      "last_progress_timestamp": ""
    }
  },
  "project_progression": {
    "current_project": {
      "project_id": "project-$(python3 -c 'import uuid; print(str(uuid.uuid4()))')",
      "name": "$PROJECT_NAME",
      "iteration": 0,
      "max_iterations": 10,
      "status": "IN_PROGRESS"
    },
    "current_phase": {
      "phase_number": 1,
      "name": "Phase 1",
      "iteration": 0,
      "max_iterations": 10,
      "status": "IN_PROGRESS",
      "total_waves_in_phase": 1,
      "waves_completed": 0
    },
    "current_wave": {
      "wave_number": 1,
      "name": "Wave 1",
      "iteration": 0,
      "max_iterations": 10,
      "status": "IN_PROGRESS",
      "efforts_completed": []
    },
    "iteration_tracking": {
      "active_container_level": "project",
      "active_container_id": "project_1"
    }
  },
  "pre_planned_infrastructure": {
    "validated": false,
    "validation_timestamp": "",
    "project_prefix": "",
    "target_repo_url": "",
    "efforts": {}
  },
  "references": {
    "bug_tracking_file": "bug-tracking.json",
    "integration_containers_file": "integration-containers.json",
    "fix_cascade_state_file": "",
    "state_machine_definition": "state-machines/software-factory-3.0-state-machine.json",
    "schema_version": "3.0.0"
  },
  "metadata": {
    "version": "3.0.0",
    "created_at": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
    "last_updated": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
    "last_updated_by": "init-script"
  },
  "planning_artifacts": {
    "master_architecture_file": "planning/PROJECT-ARCHITECTURE.md",
    "master_architecture_status": "NOT_CREATED",
    "master_architecture_created_at": "",
    "master_architecture_created_by": "architect"
  },
  "planning_files": {
    "project": {
      "implementation_plan": null,
      "architecture_plan": "planning/project/PROJECT-ARCHITECTURE-PLAN.md",
      "test_plan": null
    },
    "phases": {}
  },
  "active_agents": [],
  "project_name": "$PROJECT_NAME"
}
EOF
        log_success "Created orchestrator-state-v3.json (minimal, schema-compliant)"
    else
        # Interactive mode: use the example file
        cp orchestrator-state-v3.json.example orchestrator-state-v3.json
        if command -v jq >/dev/null 2>&1; then
            jq --arg name "$PROJECT_NAME" '.project_name = $name' orchestrator-state-v3.json > orchestrator-state-v3.json.tmp
            mv orchestrator-state-v3.json.tmp orchestrator-state-v3.json
        fi
        log_success "Created orchestrator-state-v3.json from example"
    fi

    # 2. Create bug-tracking.json
    if [ -f "bug-tracking.json.example" ]; then
        cp bug-tracking.json.example bug-tracking.json
        log_success "Created bug-tracking.json"
    else
        cat > bug-tracking.json <<EOF
{
  "bugs": [],
  "metadata": {
    "last_updated": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
    "total_bugs": 0,
    "open_bugs": 0,
    "resolved_bugs": 0
  }
}
EOF
        log_success "Created bug-tracking.json (minimal)"
    fi

    # 3. Create integration-containers.json
    if [ -f "integration-containers.json.example" ]; then
        cp integration-containers.json.example integration-containers.json
        log_success "Created integration-containers.json"
    else
        cat > integration-containers.json <<EOF
{
  "active_integrations": [],
  "metadata": {
    "last_updated": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
    "total_containers": 0,
    "converged_containers": 0
  }
}
EOF
        log_success "Created integration-containers.json (minimal)"
    fi

    # 4. fix-cascade-state.json is created dynamically when needed
    log_info "fix-cascade-state.json will be created when needed"
}

# Install pre-commit hooks
install_hooks() {
    log_info "Installing pre-commit hooks (R506 enforcement)..."

    cd "$PROJECT_ROOT"

    # Ensure hooks directory exists
    mkdir -p .git/hooks

    # Check if pre-commit hook template exists in tools/git-commit-hooks/
    if [ -f "tools/git-commit-hooks/master-pre-commit.sh" ]; then
        # Copy from tools template (the proper SF 2.0/3.0 hook)
        cp tools/git-commit-hooks/master-pre-commit.sh .git/hooks/pre-commit
        chmod +x .git/hooks/pre-commit
        log_success "Installed pre-commit hook from tools/git-commit-hooks/master-pre-commit.sh"
    elif [ -f "$SCRIPT_DIR/../.git/hooks/pre-commit" ]; then
        # We're initializing within the template repo itself
        # The hook is already installed, just verify it
        if [ -x ".git/hooks/pre-commit" ]; then
            log_success "Pre-commit hook already installed"
        else
            chmod +x .git/hooks/pre-commit
            log_success "Made pre-commit hook executable"
        fi
    else
        # Create a basic hook if neither exists
        log_warning "No pre-commit hook template found, creating basic version..."
        cat > .git/hooks/pre-commit <<'EOF'
#!/bin/bash
# Pre-commit hook for Software Factory 3.0
# Validates state files before commit

set -euo pipefail

GIT_ROOT="$(git rev-parse --show-toplevel)"
cd "$GIT_ROOT"

echo "🔍 Validating SF 3.0 state files..."

# Validate orchestrator-state-v3.json if it exists and is staged
if [ -f "orchestrator-state-v3.json" ] && git diff --cached --name-only | grep -q "orchestrator-state-v3.json"; then
    if ! jq . orchestrator-state-v3.json >/dev/null 2>&1; then
        echo "❌ ERROR: orchestrator-state-v3.json is invalid JSON"
        exit 1
    fi
    echo "✅ orchestrator-state-v3.json is valid"
fi

# Validate bug-tracking.json if it exists and is staged
if [ -f "bug-tracking.json" ] && git diff --cached --name-only | grep -q "bug-tracking.json"; then
    if ! jq . bug-tracking.json >/dev/null 2>&1; then
        echo "❌ ERROR: bug-tracking.json is invalid JSON"
        exit 1
    fi
    echo "✅ bug-tracking.json is valid"
fi

# Validate integration-containers.json if it exists and is staged
if [ -f "integration-containers.json" ] && git diff --cached --name-only | grep -q "integration-containers.json"; then
    if ! jq . integration-containers.json >/dev/null 2>&1; then
        echo "❌ ERROR: integration-containers.json is invalid JSON"
        exit 1
    fi
    echo "✅ integration-containers.json is valid"
fi

echo "✅ All validations passed"
exit 0
EOF
        chmod +x .git/hooks/pre-commit
        log_success "Created basic pre-commit hook"
    fi

    # Test the hook
    if bash .git/hooks/pre-commit >/dev/null 2>&1; then
        log_success "Pre-commit hook test passed"
    else
        log_warning "Pre-commit hook test returned non-zero (may be expected if no files staged)"
    fi
}

# Create README from template
create_readme() {
    log_info "Creating README.md from template..."

    cd "$PROJECT_ROOT"

    local template_file="templates/README.template.md"
    local readme_file="README.md"

    if [ -f "$template_file" ]; then
        # Replace template variables
        sed -e "s/{{PROJECT_NAME}}/$PROJECT_NAME/g" \
            -e "s/{{CREATION_DATE}}/$(date -u +%Y-%m-%d)/g" \
            -e "s/{{LICENSE_INFO}}/$LICENSE_INFO/g" \
            "$template_file" > "$readme_file"
        log_success "Created README.md from template"
    else
        log_warning "README template not found at $template_file, creating basic README..."
        cat > "$readme_file" <<EOF
# $PROJECT_NAME

**Software Factory 3.0 Project**

This project was initialized using Software Factory 3.0.

## Getting Started

1. Review your implementation plan (if created)
2. Run \`/init-software-factory\` to begin development
3. Monitor progress with \`/continue-software-factory\`

## State Files

- \`orchestrator-state-v3.json\` - Project state
- \`bug-tracking.json\` - Bug tracking
- \`integration-containers.json\` - Integration tracking

## Documentation

See the \`docs/\` directory for complete documentation.

---

Generated: $(date -u +%Y-%m-%d)
EOF
        log_success "Created basic README.md"
    fi
}

# Initial commit
create_initial_commit() {
    log_info "Creating initial commit..."

    cd "$PROJECT_ROOT"

    # Pre-validate state files before committing (catch issues early)
    if [ -x "tools/validate-state-file.sh" ]; then
        log_info "Pre-validating state files before commit..."

        if ! bash tools/validate-state-file.sh orchestrator-state-v3.json 2>&1 | tee /tmp/init-validation.log; then
            log_error "State file validation FAILED before commit"
            log_error "orchestrator-state-v3.json does not pass schema validation"
            log_error "See /tmp/init-validation.log for details"
            cat /tmp/init-validation.log
            log_error "Aborting initialization - will not commit invalid state"
            exit 1
        fi

        log_success "Pre-commit validation passed"
    fi

    # Stage the new files
    git add orchestrator-state-v3.json bug-tracking.json integration-containers.json README.md

    # Create commit and capture output for error detection
    log_info "Committing initial state files..."
    if git commit -m "init: $PROJECT_NAME - Software Factory 3.0 project initialized [R288]" 2>&1 | tee /tmp/init-commit-output.log; then
        log_success "Initial commit created"
    else
        log_error "Initial commit FAILED!"
        log_error "This is likely due to pre-commit hook validation failure"
        log_error "Output from commit attempt:"
        cat /tmp/init-commit-output.log
        log_error ""
        log_error "Common causes:"
        log_error "  1. State file validation failed (check orchestrator-state-v3.json)"
        log_error "  2. Pre-commit hook rejected the commit (check .git/hooks/pre-commit)"
        log_error "  3. Git configuration issues"
        log_error ""
        log_error "To debug, run:"
        log_error "  bash tools/validate-state-file.sh orchestrator-state-v3.json"
        exit 1
    fi

    # Verify the commit actually succeeded by checking if files are committed
    if ! git diff --quiet HEAD orchestrator-state-v3.json 2>/dev/null; then
        log_error "Commit verification FAILED - orchestrator-state-v3.json not committed"
        log_error "The commit appeared to succeed but files are not in git history"
        log_error "This indicates a serious pre-commit hook or git configuration issue"
        exit 1
    fi

    log_success "Commit verified - all state files committed successfully"
}

# Main initialization workflow
main() {
    echo ""
    echo "=========================================="
    echo "Software Factory 3.0 Initialization"
    echo "=========================================="
    echo ""

    parse_args "$@"
    validate_project_name

    log_info "Initializing project: $PROJECT_NAME"
    log_info "Project root: $PROJECT_ROOT"
    echo ""

    # Confirmation prompt (unless non-interactive)
    if [ $NON_INTERACTIVE -eq 0 ]; then
        read -p "Continue with initialization? (y/n) " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            log_info "Initialization cancelled"
            exit 0
        fi
    fi

    # Run initialization steps
    check_git_repo
    create_state_files
    install_hooks
    create_readme
    create_initial_commit

    echo ""
    echo "=========================================="
    log_success "Initialization Complete! 🎉"
    echo "=========================================="
    echo ""
    log_info "Next steps:"
    echo "  1. Review README.md for getting started guide"
    echo "  2. Create PROJECT-IMPLEMENTATION-PLAN.md (if not exists)"
    echo "  3. Run: /init-software-factory (in Claude Code)"
    echo ""
    log_info "State files created:"
    echo "  - orchestrator-state-v3.json"
    echo "  - bug-tracking.json"
    echo "  - integration-containers.json"
    echo "  - README.md"
    echo ""
    log_info "Pre-commit hook installed (R506 enforcement active)"
    echo ""
}

# Run main
main "$@"
