# SW Engineer - INIT_SETUP_UPSTREAM_FORK State Rules

## Purpose
Set up repository structure for a project that works with an existing upstream codebase.

## Entry Criteria
- Repository decision made: upstream_fork
- Fork URL provided in requirements
- Target directory determined

## Required Actions

### 1. Create Directory Structure
```bash
mkdir -p $CLAUDE_PROJECT_DIR/efforts/$PROJECT_PREFIX/target-repo
mkdir -p $CLAUDE_PROJECT_DIR/efforts/$PROJECT_PREFIX/project
mkdir -p $CLAUDE_PROJECT_DIR/efforts/$PROJECT_PREFIX/configs
```

### 2. Clone Fork Repository
```bash
cd $CLAUDE_PROJECT_DIR/efforts/$PROJECT_PREFIX
git clone [fork_url] target-repo
cd target-repo

# Verify clone successful
if [ ! -d .git ]; then
    ERROR: Clone failed
fi
```

### 3. Configure Remotes
```bash
# Add upstream remote
git remote add upstream [upstream_url]
git remote -v

# Fetch upstream
git fetch upstream
git fetch origin
```

### 4. Set Up Branches
```bash
# Ensure on main/master
git checkout [main_branch]

# Create development branch
git checkout -b sf2-development

# Create feature branch for Phase 1
git checkout -b phase-1-wave-1
```

### 5. Initialize Project Directory
```bash
cd $CLAUDE_PROJECT_DIR/efforts/$PROJECT_PREFIX/project

# Create basic structure
mkdir -p src
mkdir -p tests
mkdir -p docs
mkdir -p configs

# Create initial .gitignore if needed
if [ ! -f .gitignore ]; then
    echo "# SF2 Generated" > .gitignore
    echo "*.tmp" >> .gitignore
    echo "*.log" >> .gitignore
fi
```

### 6. Verify Setup
```bash
# Check all directories exist
ls -la $CLAUDE_PROJECT_DIR/efforts/$PROJECT_PREFIX/

# Verify git configuration
cd $CLAUDE_PROJECT_DIR/efforts/$PROJECT_PREFIX/target-repo
git remote -v
git branch -a
```

### 7. Update State File
Record in init-state-${PROJECT_PREFIX}.json:
```json
"repository_setup": {
  "type": "upstream_fork",
  "target_repo_path": "efforts/$PROJECT_PREFIX/target-repo",
  "project_path": "efforts/$PROJECT_PREFIX/project",
  "fork_url": "[url]",
  "upstream_url": "[url]",
  "main_branch": "[branch]",
  "setup_complete": true
}
```

## Exit Criteria
- Fork cloned successfully
- Remotes configured (origin + upstream)
- Branches created
- Directory structure ready
- State file updated

## Transition
**MANDATORY**: → INIT_GENERATE_CONFIGS

## Error Handling

### Clone Failures
- Network timeout → Retry with backoff
- Auth required → Prompt for credentials
- Invalid URL → Return to requirements

### Remote Configuration
- Upstream doesn't exist → Document as new project
- Branch doesn't exist → Use default branch

## Validation Checks
- [ ] target-repo/.git exists
- [ ] upstream remote configured
- [ ] Can fetch from both remotes
- [ ] Development branches created
- [ ] Project directory initialized