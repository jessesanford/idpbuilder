# SW Engineer - INIT_SETUP_UPSTREAM_FORK State Rules

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

