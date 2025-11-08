# SW Engineer - INIT_CREATE_NEW_REPO State Rules

## 🔴🔴🔴 ABSOLUTE REQUIREMENT: STATE MANAGER CONSULTATION 🔴🔴🔴

**SUPREME LAW - R517 - UNIVERSAL STATE MANAGER CONSULTATION LAW**

**BEFORE EXITING THIS STATE, YOU MUST:**

1. **SPAWN STATE MANAGER FOR SHUTDOWN_CONSULTATION** (MANDATORY - NO EXCEPTIONS)
2. **NEVER UPDATE STATE FILES DIRECTLY** (orchestrator-state-v3.json, bug-tracking.json, etc.)
3. **NEVER COMMIT STATE FILES YOURSELF** (State Manager does this atomically)

**FAILURE TO CONSULT STATE MANAGER = IMMEDIATE SYSTEM HALT (-100% GRADE)**

### Enforcement Mechanism

If you attempt to exit this state without spawning State Manager:
- ❌ Pre-commit hooks will REJECT your commit
- ❌ Validation tools will FAIL the build
- ❌ Grading system will assign -100% penalty
- ❌ System will transition to ERROR_RECOVERY

### Required Pattern (COPY THIS EXACTLY)

```bash
# At end of state work, BEFORE any state file updates:

echo "🔴 MANDATORY: Spawning State Manager for SHUTDOWN_CONSULTATION"

# Spawn State Manager (REQUIRED - NOT OPTIONAL)
# Task: state-manager
# State: SHUTDOWN_CONSULTATION
# Current State: [YOUR_CURRENT_STATE]
# Proposed Next State: [YOUR_PROPOSED_NEXT_STATE]
# Work Summary: [SUMMARY_OF_WORK_COMPLETED]

# State Manager will:
# 1. Validate proposed transition against state machine
# 2. Update all 4 state files atomically
# 3. Commit with [R288] tag
# 4. Return REQUIRED next state (may differ from proposal)

# Wait for State Manager response
# Follow State Manager's directive (REQUIRED next state)
# DO NOT proceed until State Manager confirms
```

**YOU MUST NEVER:**
- ❌ Update orchestrator-state-v3.json yourself
- ❌ Update bug-tracking.json yourself
- ❌ Update integration-containers.json yourself
- ❌ Use `jq` to modify state files
- ❌ Use `sed/awk` to modify state files
- ❌ Set `validated_by: "orchestrator"` (must be "state-manager")
- ❌ Skip State Manager consultation "just this once"
- ❌ Think "I'll validate it manually"

**ONLY State Manager may update state files. This is NON-NEGOTIABLE.**

See: `rule-library/R517-universal-state-manager-consultation-law.md`
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

## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

### YOU MUST OUTPUT THE CONTINUATION FLAG AS YOUR LAST ACTION

**EVERY STATE COMPLETION MUST END WITH EXACTLY ONE OF:**
```bash
# If state completed successfully and factory should continue:
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"

# If error/block/manual review needed:
echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
```

### CRITICAL REQUIREMENTS:
1. **ABSOLUTE LAST OUTPUT**: This MUST be the very last line of output before state completion
2. **NO TEXT AFTER**: No explanations, summaries, or any text after the flag
3. **EXACTLY AS SHOWN**: Use exact format - no variations like CONTINUE-ORCHESTRATING
4. **ALWAYS REQUIRED**: Every single state must output this flag
5. **GREPPABLE**: Must be on its own line for automation parsing

### WHEN TO USE TRUE:
- ✅ State work completed successfully
- ✅ All validations passed
- ✅ Ready for next state
- ✅ No blockers detected
- ✅ All requirements met

### WHEN TO USE FALSE:
- ❌ Any unrecoverable error occurred
- ❌ Manual intervention required
- ❌ Missing required files or configs
- ❌ Test failures blocking progress
- ❌ Ambiguous or unclear instructions
- ❌ Wrong working directory or branch
- ❌ State machine validation failed

### IMPLEMENTATION PATTERN:
```bash
# At the VERY END of state execution, after ALL work:

# Determine success/failure
if [[ "$STATE_COMPLETED_PROJECT_DONEFULLY" == "true" ]]; then
    echo "✅ State work completed successfully"
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
else
    echo "❌ State encountered issues requiring intervention"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
fi

# DO NOT OUTPUT ANYTHING AFTER THE FLAG!
```

### GRADING IMPACT:
- **Missing flag**: -100% AUTOMATIC FAILURE
- **Text after flag**: -50% penalty
- **Wrong format**: -100% AUTOMATIC FAILURE
- **Multiple flags**: -50% penalty

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

