# SW Engineer - INIT_CREATE_NEW_REPO State Rules

## Purpose
Initialize a new repository for a greenfield project.

## Entry Criteria
- Repository decision made: new_repository
- Project name provided
- No upstream codebase

## Required Actions

### 1. Create Full Directory Structure
```bash
# Create main project directory
mkdir -p $CLAUDE_PROJECT_DIR/efforts/$PROJECT_PREFIX/project

# Create subdirectories based on language
cd $CLAUDE_PROJECT_DIR/efforts/$PROJECT_PREFIX/project

# For most languages
mkdir -p src
mkdir -p tests
mkdir -p docs
mkdir -p configs
mkdir -p scripts
mkdir -p .github/workflows

# Language-specific additions
# For Go:
mkdir -p cmd pkg internal

# For Python:
mkdir -p src/${PROJECT_PREFIX}

# For Node.js:
mkdir -p src lib bin
```

### 2. Initialize Git Repository
```bash
cd $CLAUDE_PROJECT_DIR/efforts/$PROJECT_PREFIX/project
git init
git config user.name "SF2 Developer"
git config user.email "dev@sf2.local"
```

### 3. Create Initial Files

#### .gitignore
```bash
cat > .gitignore << 'EOF'
# Build artifacts
build/
dist/
*.egg-info/
__pycache__/
node_modules/
target/
*.o
*.so
*.dylib

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# Environment
.env
.env.local
venv/
.venv/

# Logs
*.log
logs/

# OS
.DS_Store
Thumbs.db

# Testing
coverage/
.coverage
*.cover
.pytest_cache/
EOF
```

#### README.md
```bash
cat > README.md << EOF
# $PROJECT_NAME

$PROJECT_IDEA

## Status
Project initialized by Software Factory 2.0
Phase 1 implementation pending

## Structure
- src/ - Source code
- tests/ - Test files
- docs/ - Documentation
- configs/ - Configuration files

## Development
See IMPLEMENTATION-PLAN.md for development phases
EOF
```

#### Initial Source File (language-specific)
Create a minimal valid source file:
- Go: src/main.go with package main
- Python: src/__init__.py
- JavaScript: src/index.js
- Java: src/Main.java

### 4. Make Initial Commit
```bash
git add -A
git commit -m "Initial commit: SF2 project structure

- Created directory structure
- Added .gitignore
- Added README.md
- Initialized for $LANGUAGE development"
```

### 5. Create Development Branch
```bash
# Create main branch
git branch -M main

# Create development branch
git checkout -b development

# Create Phase 1 branch
git checkout -b phase-1-wave-1
```

### 6. Update State File
Record in init-state-${PROJECT_PREFIX}.json:
```json
"repository_setup": {
  "type": "new_repository",
  "project_path": "efforts/$PROJECT_PREFIX/project",
  "git_initialized": true,
  "initial_commit": "[hash]",
  "main_branch": "main",
  "current_branch": "phase-1-wave-1",
  "setup_complete": true
}
```

## Exit Criteria
- Git repository initialized
- Initial commit made
- Basic file structure created
- Development branches set up
- State file updated

## Transition
**MANDATORY**: → INIT_GENERATE_CONFIGS

## Language-Specific Initialization

### Go Project
```go
// src/main.go
package main

import "fmt"

func main() {
    fmt.Println("SF2 Project Initialized")
}
```

### Python Project
```python
# src/__init__.py
"""SF2 Project - $PROJECT_NAME"""
__version__ = "0.1.0"
```

### Node.js Project
```javascript
// src/index.js
console.log('SF2 Project Initialized');
```

## Validation Checks
- [ ] .git directory exists
- [ ] Initial commit successful
- [ ] README.md created
- [ ] .gitignore appropriate for language
- [ ] Source file validates for language
- [ ] Branches created correctly