# ⚠️⚠️⚠️ R178: Effort Directory Structure ⚠️⚠️⚠️

**Category:** Utility and Hook Management  
**Agents:** orchestrator (creates), sw-engineer (uses), code-reviewer (validates)  
**Criticality:** WARNING - Wrong structure causes workflow failures  
**Related:** R176 (Workspace Isolation), R209 (Effort Directory Isolation)

## MANDATORY DIRECTORY STRUCTURE

### 1. STANDARD EFFORT STRUCTURE

Every effort MUST follow this exact structure:

```
efforts/
└── phase{N}/
    └── wave{N}/
        └── {effort-name}/
            ├── .git/                    # Git repository (sparse)
            ├── pkg/                     # Code goes here
            │   └── {module}/            # Module organization
            │       ├── {feature}.go     # Implementation
            │       └── {feature}_test.go # Tests
            ├── IMPLEMENTATION-PLAN.md   # From code-reviewer
            ├── REVIEW-FEEDBACK.md       # Review results
            ├── SPLIT-INVENTORY.md       # If split required
            └── effort-metadata.yaml     # Effort tracking
```

### 2. METADATA FILE FORMAT

```yaml
# effort-metadata.yaml
effort:
  name: auth-base
  phase: 1
  wave: 1
  branch: phase1/wave1/auth-base
  
status:
  state: IMPLEMENTATION  # Current state
  lines_written: 0
  last_measured: null
  review_status: PENDING
  
agent:
  type: sw-engineer
  spawned_at: 2024-01-20T10:00:00Z
  pid: 12345
  
dependencies:
  - none  # or list effort names
  
size_tracking:
  current: 0
  limit: 800
  measurements:
    - timestamp: 2024-01-20T10:30:00Z
      lines: 245
    - timestamp: 2024-01-20T11:00:00Z
      lines: 487
```

### 3. SPLIT EFFORT STRUCTURE

When splitting is required:

```
efforts/
└── phase1/
    └── wave1/
        └── auth-base/
            ├── pkg/               # Original code
            ├── splits/            # Split organization
            │   ├── split-001/
            │   │   ├── pkg/       # Split 1 code
            │   │   └── split-metadata.yaml
            │   ├── split-002/
            │   │   ├── pkg/       # Split 2 code
            │   │   └── split-metadata.yaml
            │   └── SPLIT-PLAN.md
            └── SPLIT-INVENTORY.md # Tracks all splits
```

### 4. CREATION PROTOCOL (Orchestrator)

```bash
create_effort_structure() {
    local phase="$1"
    local wave="$2"
    local effort_name="$3"
    local base_dir="efforts/phase${phase}/wave${wave}/${effort_name}"
    
    echo "📁 Creating effort structure: $base_dir"
    
    # Create directory hierarchy
    mkdir -p "$base_dir/pkg"
    
    # Initialize git (sparse checkout)
    cd "$base_dir"
    git init
    git remote add origin "$TARGET_REPO_URL"
    git sparse-checkout init --cone
    git sparse-checkout set pkg/
    
    # Create branch
    git checkout -b "phase${phase}/wave${wave}/${effort_name}"
    
    # Create metadata
    cat > effort-metadata.yaml << EOF
effort:
  name: $effort_name
  phase: $phase
  wave: $wave
  branch: phase${phase}/wave${wave}/${effort_name}
  
status:
  state: INIT
  lines_written: 0
  last_measured: null
  review_status: PENDING
EOF
    
    echo "✅ Effort structure created"
}
```

### 5. USAGE RULES (SW-Engineer)

```bash
# ALWAYS verify structure before work
verify_effort_structure() {
    local effort_dir="$1"
    
    # Check required directories
    [ -d "$effort_dir/pkg" ] || {
        echo "❌ Missing pkg/ directory"
        return 178
    }
    
    # Check metadata
    [ -f "$effort_dir/effort-metadata.yaml" ] || {
        echo "❌ Missing effort-metadata.yaml"
        return 178
    }
    
    # Check git
    [ -d "$effort_dir/.git" ] || {
        echo "❌ Not a git repository"
        return 178
    }
    
    echo "✅ Effort structure valid"
}
```

### 6. FILE PLACEMENT RULES

```bash
# ✅ CORRECT - In effort's pkg directory
efforts/phase1/wave1/auth-base/pkg/auth/handler.go
efforts/phase1/wave1/auth-base/pkg/auth/handler_test.go

# ❌ WRONG - Outside pkg directory
efforts/phase1/wave1/auth-base/handler.go
efforts/phase1/wave1/auth-base/src/handler.go

# ❌ WRONG - In root pkg
/workspace/pkg/auth/handler.go

# ❌ WRONG - In another effort
efforts/phase1/wave1/user-service/pkg/auth/handler.go
```

### 7. DOCUMENTATION FILES

Each effort should maintain:

```markdown
# IMPLEMENTATION-PLAN.md
Created by code-reviewer, contains:
- Detailed implementation steps
- File creation order
- Testing requirements

# REVIEW-FEEDBACK.md
Created after review, contains:
- Issues found
- Required fixes
- Suggestions

# SPLIT-INVENTORY.md (if applicable)
Tracks split progress:
- Split boundaries
- Completion status
- Dependencies between splits
```

### 8. STATE-BASED FILES

Different files appear based on state:

```yaml
INIT:
  - effort-metadata.yaml
  
IMPLEMENTATION:
  - IMPLEMENTATION-PLAN.md
  - pkg/**/*.go
  
CODE_REVIEW:
  - REVIEW-FEEDBACK.md
  
SPLIT_REQUIRED:
  - SPLIT-PLAN.md
  - splits/
  
COMPLETE:
  - COMPLETION-REPORT.md
```

### 9. COMMON STRUCTURE VIOLATIONS

**AVOID THESE:**
- ❌ Creating code outside pkg/ directory
- ❌ Missing effort-metadata.yaml
- ❌ Wrong directory naming (must match pattern)
- ❌ Mixing efforts in same directory
- ❌ Using absolute paths instead of relative
- ❌ Creating files in planning repository

### 10. VALIDATION SCRIPT

```bash
#!/bin/bash
validate_all_effort_structures() {
    local base_dir="efforts"
    
    # Find all efforts
    for effort_dir in $(find $base_dir -type d -name "phase*" | \
                       xargs -I {} find {} -maxdepth 3 -mindepth 3 -type d); do
        
        echo "Checking: $effort_dir"
        
        # Validate structure
        if [ ! -d "$effort_dir/pkg" ]; then
            echo "  ❌ Missing pkg/ directory"
            continue
        fi
        
        if [ ! -f "$effort_dir/effort-metadata.yaml" ]; then
            echo "  ❌ Missing metadata file"
            continue
        fi
        
        echo "  ✅ Valid structure"
    done
}
```

## GRADING IMPACT

```yaml
structure_violations:
  wrong_directory_structure: -15%
  missing_metadata_file: -10%
  code_outside_pkg: -20%
  no_git_repository: -15%
  contaminated_structure: -30%
```

## INTEGRATE_WAVE_EFFORTS WITH OTHER RULES

- **R176**: Workspace isolation requirement
- **R209**: Effort directory isolation protocol
- **R181**: Orchestrator workspace setup responsibility
- **R204**: Orchestrator split infrastructure creation

## SUMMARY

**R178 Core Mandate: Every effort follows exact directory structure!**

- Standard structure for all efforts
- Metadata tracking required
- Code only in pkg/ directory
- Proper git initialization
- Documentation in effort root

---
**Created**: Standardize effort organization
**Purpose**: Enable consistent workflow and tooling
**Enforcement**: WARNING - Structure violations break automation