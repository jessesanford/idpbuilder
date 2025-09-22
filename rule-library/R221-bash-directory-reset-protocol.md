# 🔴🔴🔴 R221: BASH DIRECTORY RESET PROTOCOL - SUPREME LAW #1 🔴🔴🔴

**Category:** SUPREME LAW  
**Agents:** ALL AGENTS - NO EXCEPTIONS  
**Criticality:** ABSOLUTE - VIOLATION = -100% GRADE (AUTOMATIC FAILURE)  
**Priority:** SUPREME LAW #1 - THE HIGHEST PRIORITY RULE IN THE ENTIRE SYSTEM

## 🚨🚨🚨 THIS IS SUPREME LAW #1 - SUPERSEDES ALL OTHER RULES 🚨🚨🚨

### THE FUNDAMENTAL TRUTH

**BASH RESETS TO HOME DIRECTORY AFTER EVERY COMMAND!**

This is not a bug, it's how Bash tool works. Every new Bash invocation starts fresh.

### 🔴 THE ABSOLUTE LAW

**YOU MUST CD TO YOUR WORKING DIRECTORY IN EVERY SINGLE BASH COMMAND!**

```bash
# ❌❌❌ THIS WILL FAIL CATASTROPHICALLY:
Bash: cd /path/to/effort
Bash: git add .  # WRONG! You're in HOME, not effort directory!
Bash: git commit -m "fix"  # WRONG! Commits random home files!

# ✅✅✅ THIS IS THE ONLY CORRECT WAY:
EFFORT_DIR="/path/to/effort"  # Set this variable ONCE
Bash: cd $EFFORT_DIR && git add .
Bash: cd $EFFORT_DIR && git commit -m "fix"
Bash: cd $EFFORT_DIR && git push
```

### 🔴 ENFORCEMENT FOR ALL AGENTS

#### SW-ENGINEER
```bash
# Set your effort directory once
EFFORT_DIR="/workspace/efforts/phase1/wave1/auth-implementation"

# Then use it in EVERY command:
Bash: cd $EFFORT_DIR && go test ./...
Bash: cd $EFFORT_DIR && git add pkg/
Bash: cd $EFFORT_DIR && git commit -m "feat: implement auth"
```

#### CODE-REVIEWER
```bash
# Set your review directory once  
REVIEW_DIR="/workspace/efforts/phase1/wave1/auth-implementation"

# Then use it in EVERY command:
Bash: cd $REVIEW_DIR && $CLAUDE_PROJECT_DIR/tools/line-counter.sh
Bash: cd $REVIEW_DIR && git diff main...HEAD
```

#### ORCHESTRATOR
```bash
# When spawning agents - CD FIRST per R208:
TARGET_DIR="/workspace/efforts/phase1/wave1/auth-implementation"
Bash: cd $TARGET_DIR && claude_spawn sw-engineer
```

#### ARCHITECT
```bash
# Set assessment directory once
ASSESS_DIR="/workspace/wave-reviews/phase1-wave1"

# Use in every command:
Bash: cd $ASSESS_DIR && ls -la
Bash: cd $ASSESS_DIR && cat WAVE-ASSESSMENT.md
```

### 🔴 CRITICAL VIOLATIONS TO AVOID

```bash
# ❌ ABSOLUTE FAILURE #1: Assuming directory persists
Bash: cd /some/path
Bash: ls  # You're NOT in /some/path anymore!

# ❌ ABSOLUTE FAILURE #2: Multi-step operations without CD
Bash: mkdir pkg/auth
Bash: touch pkg/auth/handler.go  # Creates in wrong location!

# ❌ ABSOLUTE FAILURE #3: Git operations without CD  
Bash: git add .  # Adds random home directory files!
Bash: git commit  # Commits wrong files!
Bash: git push  # Pushes contaminated commits!
```

### 🔴 STATE-SPECIFIC APPLICATIONS

#### IMPLEMENTATION State
```bash
EFFORT_DIR="/workspace/efforts/phase1/wave1/my-effort"
Bash: cd $EFFORT_DIR && mkdir -p pkg/feature
Bash: cd $EFFORT_DIR && touch pkg/feature/implementation.go
```

#### SPLIT_IMPLEMENTATION State
```bash
SPLIT_DIR="/workspace/efforts/phase1/wave1/my-effort/splits/split-001"
Bash: cd $SPLIT_DIR && go build ./...
```

#### MEASURE_SIZE State  
```bash
MEASURE_DIR="/workspace/efforts/phase1/wave1/my-effort"
Bash: cd $MEASURE_DIR && $CLAUDE_PROJECT_DIR/tools/line-counter.sh
```

#### FIX_ISSUES State
```bash
FIX_DIR="/workspace/efforts/phase1/wave1/my-effort"
Bash: cd $FIX_DIR && git add -u
Bash: cd $FIX_DIR && git commit -m "fix: address review feedback"
```

### 🔴 GRADING IMPACT

```yaml
r221_violations:
  git_operations_wrong_dir: -100%  # AUTOMATIC FAILURE
  creating_files_wrong_location: -100%  # AUTOMATIC FAILURE  
  assuming_directory_persists: -50%  # Major violation
  not_using_cd_prefix: -30%  # Per occurrence
  contaminating_home_directory: -100%  # AUTOMATIC FAILURE
```

### 🔴 VALIDATION SCRIPT

```bash
#!/bin/bash
# Every agent MUST run this validation

validate_r221_compliance() {
    # Check if in correct directory
    if [[ "$(pwd)" == "$HOME" ]]; then
        echo "🔴🔴🔴 R221 VIOLATION: In HOME directory!"
        echo "🔴🔴🔴 SUPREME LAW #1 VIOLATED!"
        exit 221
    fi
    
    # Set the working directory variable
    export WORKING_DIR="$(pwd)"
    echo "✅ R221: Working directory set to $WORKING_DIR"
    echo "✅ R221: Will use 'cd \$WORKING_DIR &&' in all commands"
}

# Run immediately
validate_r221_compliance
```

### 🔴 COMPOUND COMMAND RULES

When using compound commands, CD must be first:

```bash
# ✅ CORRECT: CD first in compound
Bash: cd $EFFORT_DIR && mkdir -p pkg && touch pkg/main.go && git add pkg/

# ❌ WRONG: CD in middle
Bash: mkdir -p pkg && cd $EFFORT_DIR && touch pkg/main.go

# ✅ CORRECT: Using subshell
Bash: (cd $EFFORT_DIR && go test ./... && git add -u && git commit -m "test: add tests")

# ✅ CORRECT: Using semicolons
Bash: cd $EFFORT_DIR; go build ./...; git add .; git commit -m "build: compile"
```

### 🔴 INTEGRATION WITH OTHER SUPREME LAWS

1. **Works with R208 (SUPREME LAW #2)**: Orchestrator must CD before spawn
2. **Works with R235 (SUPREME LAW #3)**: Pre-flight checks verify you're in right directory
3. **Works with R234 (SUPREME LAW #4)**: State transitions don't change this requirement

### 🔴 NO EXCEPTIONS CLAUSE

**THERE ARE NO EXCEPTIONS TO THIS RULE:**
- Not for "efficiency"
- Not for "I know what I'm doing"  
- Not for "continuous operation"
- Not for "the orchestrator said"
- Not for ANY reason

### 🔴 RECOVERY FROM VIOLATIONS

If you violate R221:
1. **STOP IMMEDIATELY**
2. **Assess contamination** (what files were affected)
3. **Clean up wrong files** (remove/revert as needed)
4. **Reset to correct directory**
5. **Resume with proper CD prefix**

## 🔴 SUMMARY

**R221 IS SUPREME LAW #1 - THE HIGHEST RULE**

**EVERY BASH COMMAND MUST START WITH: `cd $YOUR_DIR &&`**

**NO EXCEPTIONS - NO EXCUSES - NO ALTERNATIVES**

This prevents:
- Wrong directory operations
- File contamination
- Repository corruption  
- Git disasters
- Total system failure

**PENALTY FOR VIOLATION: -100% GRADE (AUTOMATIC FAILURE)**

---
**Created**: Emergency response to repeated directory disasters  
**Effective**: IMMEDIATELY - ALL AGENTS MUST COMPLY
**Enforcement**: ABSOLUTE - THIS IS SUPREME LAW #1