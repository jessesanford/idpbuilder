# ⚠️⚠️⚠️ RULE R316 - Orchestrator Commit Restrictions

**Criticality:** WARNING  
**Grading Impact:** -30% for violations  
**Enforcement:** Commit content validation
**Applies To:** Orchestrator when committing changes

## Rule Statement

The Orchestrator may ONLY commit orchestration-related files. ANY commit containing code files MUST be rejected. The Orchestrator must verify commit contents before pushing.

## Allowed Commit Contents

### ✅ ORCHESTRATOR MAY COMMIT:
- `orchestrator-state.json`
- `*.md` files (documentation/planning)
- `todos/*.todo` files
- `.gitkeep` files
- `YAML/JSON` configuration (non-code)
- Directory structure changes (empty)
- Planning documents
- State tracking files

### ❌ ORCHESTRATOR MUST NEVER COMMIT:
- Source code files (any language)
- Test files
- Implementation files
- Build artifacts
- Compiled binaries
- Package files
- Script files (unless orchestration scripts)

## Validation Protocol

### REQUIRED: Pre-Commit Validation
```bash
# MUST run before EVERY orchestrator commit
validate_orchestrator_commit() {
    echo "🔍 Validating commit contents for R316 compliance..."
    
    # Check staged files
    local staged_files=$(git diff --cached --name-only)
    
    for file in $staged_files; do
        # Check for code files
        if [[ "$file" =~ \.(go|py|js|ts|java|cpp|c|rs|rb|php|jsx|tsx|vue|swift|kt|scala)$ ]]; then
            echo "🚨 R316 VIOLATION: Code file in commit: $file"
            echo "❌ Orchestrator cannot commit code files!"
            git reset HEAD "$file"
            return 1
        fi
        
        # Check for test files
        if [[ "$file" =~ _(test|spec)\. ]]; then
            echo "🚨 R316 VIOLATION: Test file in commit: $file"
            echo "❌ Orchestrator cannot commit test files!"
            git reset HEAD "$file"
            return 1
        fi
    done
    
    echo "✅ Commit contents validated - no code files"
    return 0
}
```

## Required Commit Pattern

### EVERY Orchestrator Commit MUST:
```bash
# 1. Stage only allowed files
git add orchestrator-state.json todos/*.todo *.md

# 2. Validate before commit
validate_orchestrator_commit || {
    echo "❌ Commit blocked by R316"
    exit 1
}

# 3. Commit with clear prefix
git commit -m "orchestrator: update state to MONITOR" \
           -m "No code files included per R316"

# 4. Push immediately
git push
```

## Common Violations

### ❌ VIOLATION: Committing After Code Copy
```bash
# Orchestrator copies files (already a violation)
cp *.go efforts/split2/

# Then commits everything (double violation)
git add -A
git commit -m "Set up split2"  # VIOLATION: Contains code
```

### ❌ VIOLATION: Blanket Add Commands
```bash
# Dangerous commands that may capture code
git add .  # MAY include code files
git add -A  # MAY include code files
git add *  # MAY include code files

# MUST use specific adds instead
git add orchestrator-state.json todos/*.todo
```

## Correct Patterns

### ✅ GOOD: Specific File Commits
```bash
# Only add orchestration files
git add orchestrator-state.json
git add todos/orchestrator-*.todo
git add efforts/phase1/PLAN.md

# Validate and commit
validate_orchestrator_commit && \
git commit -m "orchestrator: transition to SPAWN_AGENTS state"
```

### ✅ GOOD: State-Only Updates
```bash
# Update state file
jq '.current_state = "MONITOR"' orchestrator-state.json

# Commit ONLY the state file
git add orchestrator-state.json
git commit -m "orchestrator: update state to MONITOR"
git push
```

## Enforcement Hooks

### Git Pre-Commit Hook (Recommended)
```bash
#!/bin/bash
# .git/hooks/pre-commit

# Detect if orchestrator is committing
if [[ "$GIT_AUTHOR_NAME" == *"orchestrator"* ]]; then
    echo "🔍 R316: Validating orchestrator commit..."
    
    # Check for code files
    git diff --cached --name-only | while read file; do
        if [[ "$file" =~ \.(go|py|js|ts|java|cpp|c|rs|rb|php)$ ]]; then
            echo "🚨 R316 VIOLATION: Orchestrator cannot commit $file"
            exit 1
        fi
    done
fi
```

## Recovery from Violations

If orchestrator accidentally stages code files:
```bash
# 1. Reset all staged files
git reset HEAD

# 2. Add ONLY allowed files
git add orchestrator-state.json todos/*.todo

# 3. Validate before proceeding
validate_orchestrator_commit

# 4. Commit with compliance note
git commit -m "orchestrator: state update (R316 compliant)"
```

## Grading Impact

- Committing code files: -30% per incident
- Pushing code changes: -50% (escalated)
- Repeated violations: Progressive penalties

## Self-Audit Command

```bash
# Orchestrator should run periodically
audit_my_commits() {
    echo "📊 Auditing last 10 orchestrator commits for R316..."
    
    git log --author="orchestrator" --name-only -10 | \
    grep -E '\.(go|py|js|ts|java|cpp|c|rs|rb|php)$' && {
        echo "🚨 WARNING: Found code files in recent commits!"
        return 1
    } || {
        echo "✅ All recent commits are R316 compliant"
        return 0
    }
}
```

---
**REMEMBER:** Orchestrator commits = orchestration only. NO CODE EVER.