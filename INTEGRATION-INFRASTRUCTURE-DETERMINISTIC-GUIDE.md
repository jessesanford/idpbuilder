# 🔴🔴🔴 INTEGRATION INFRASTRUCTURE DETERMINISTIC GUIDE 🔴🔴🔴

**THIS DOCUMENT PROVIDES EXACT, UNAMBIGUOUS SPECIFICATIONS FOR ALL INTEGRATION INFRASTRUCTURE**

## 📋 EXECUTIVE SUMMARY

The orchestrator was creating incorrect integration infrastructure paths because rules were ambiguous. This guide provides **EXACT, DETERMINISTIC** specifications that leave **NO ROOM FOR INTERPRETATION**.

## 🔴 THE PROBLEM

**AMBIGUOUS SPECIFICATION (OLD):**
```
Create integration workspace at: /efforts/phase1/wave1/integration-workspace
```

**RESULT:** Orchestrator didn't know if:
- Target repo should be cloned AS integration-workspace
- Target repo should be cloned INTO integration-workspace  
- What to do on re-integration attempts

This created the bug: target repo being cloned with a `/repo` subdirectory inside integration-workspace

## ✅ THE SOLUTION: DETERMINISTIC SPECIFICATIONS

### 1. EXACT DIRECTORY STRUCTURE (NO AMBIGUITY)

```bash
# DETERMINISTIC STRUCTURE FOR ALL INTEGRATIONS:
${SF_INSTANCE_DIR}/                           # Software Factory instance root
└── efforts/                                  # All implementation work
    ├── phase1/
    │   ├── wave1/
    │   │   ├── effort-foo/                  # Effort workspace (target repo clone)
    │   │   │   ├── .software-factory/      # Effort metadata
    │   │   │   └── src/                    # Implementation code
    │   │   ├── effort-bar/                  # Effort workspace (target repo clone)
    │   │   │   ├── .software-factory/      # Effort metadata
    │   │   │   └── src/                    # Implementation code
    │   │   └── integration-workspace/       # Target repo cloned HERE directly
    │   │       ├── .git/
    │   │       ├── .software-factory/      # Integration metadata
    │   │       │   ├── WAVE-MERGE-PLAN.md
    │   │       │   └── INTEGRATION-REPORT.md
    │   │       ├── src/                    # Merged code
    │   │       └── ... (project files)
    │   └── phase-integration-workspace/     # Phase integration (target repo clone)
    │       ├── .software-factory/          # Phase integration metadata
    │       └── ... (project files)
    └── project-integration-workspace/       # Project integration (target repo clone)
        ├── .software-factory/              # Project integration metadata
        └── ... (project files)
```

**KEY SPECIFICATIONS:**
- Target repository is ALWAYS cloned directly as `integration-workspace`
- ALL metadata goes in `.software-factory/` subdirectory
- NEVER create a `/repo/` subdirectory - this pattern is deprecated

### 2. DETERMINISTIC CLONE COMMANDS

```bash
# ✅ CORRECT (DETERMINISTIC):
INTEGRATION_WORKSPACE="${SF_INSTANCE_DIR}/efforts/phase2/wave1/integration-workspace"

# Clone target repo directly as integration-workspace
git clone --branch "$BASE_BRANCH" "$TARGET_REPO_URL" "$INTEGRATION_WORKSPACE"
cd "$INTEGRATION_WORKSPACE"

# Create metadata directory
mkdir -p .software-factory

# ❌ WRONG (OLD PATTERN):
mkdir -p "$INTEGRATION_WORKSPACE"
git clone "$TARGET_REPO_URL" "${INTEGRATION_WORKSPACE}/repo"  # Creates unnecessary subdirectory!
```

### 3. RE-INTEGRATION HANDLING (DETERMINISTIC)

When integration fails and must be redone:

```bash
# DETERMINISTIC RE-INTEGRATION PROTOCOL:
handle_reintegration() {
    local INTEGRATION_WORKSPACE="$1"
    
    # Step 1: Archive old integration (numbered suffix)
    if [ -d "$INTEGRATION_WORKSPACE" ]; then
        ARCHIVE_NUM=1
        while [ -d "${INTEGRATION_WORKSPACE}-archived-${ARCHIVE_NUM}" ]; do
            ARCHIVE_NUM=$((ARCHIVE_NUM + 1))
        done
        mv "$INTEGRATION_WORKSPACE" "${INTEGRATION_WORKSPACE}-archived-${ARCHIVE_NUM}"
        echo "📦 Archived to: ${INTEGRATION_WORKSPACE}-archived-${ARCHIVE_NUM}"
    fi
    
    # Step 2: Clone target repo directly as workspace
    git clone --branch "$BASE_BRANCH" "$TARGET_REPO_URL" "$INTEGRATION_WORKSPACE"
    
    # Step 3: Create metadata directory
    cd "$INTEGRATION_WORKSPACE"
    mkdir -p .software-factory
    
    # Step 4: Use SAME branch name (force-push)
    git checkout -b "$INTEGRATION_BRANCH"
    git push --force-with-lease -u origin "$INTEGRATION_BRANCH"
}
```

**ARCHIVE NAMING:**
- First re-integration: `integration-workspace-archived-1`
- Second re-integration: `integration-workspace-archived-2`
- Nth re-integration: `integration-workspace-archived-N`

### 4. BRANCH NAMING (DETERMINISTIC)

**Integration branches use FIXED names (no versioning):**
- Wave integration: `phase{N}-wave{M}-integration`
- Phase integration: `phase{N}-integration`
- Project integration: `project-integration`

**On re-integration:**
- Use SAME branch name
- Force-push updates
- NEVER create versioned branches (e.g., phase1-wave1-integration-v2)

### 5. BASE BRANCH DETERMINATION (R308 COMPLIANT)

```bash
# DETERMINISTIC BASE BRANCH LOGIC:
determine_integration_base() {
    local TYPE=$1 PHASE=$2 WAVE=$3
    
    case "$TYPE" in
        "wave")
            if [[ $PHASE -eq 1 && $WAVE -eq 1 ]]; then
                echo "main"
            elif [[ $WAVE -eq 1 ]]; then
                echo "phase$((PHASE-1))-integration"  # NOT main!
            else
                echo "phase${PHASE}-wave$((WAVE-1))-integration"
            fi
            ;;
        "phase")
            LAST_WAVE=$(get_last_wave_number $PHASE)
            echo "phase${PHASE}-wave${LAST_WAVE}-integration"
            ;;
        "project")
            LAST_PHASE=$(get_last_phase_number)
            echo "phase${LAST_PHASE}-integration"
            ;;
    esac
}
```

## 📊 COMPLETE EXAMPLES

### Example 1: Phase 2, Wave 1 Integration (First Attempt)

```bash
# Starting state
SF_INSTANCE_DIR="/home/vscode/software-factory-template"
PHASE=2
WAVE=1
TYPE="wave"

# Determine paths (DETERMINISTIC)
INTEGRATION_WORKSPACE="${SF_INSTANCE_DIR}/efforts/phase2/wave1/integration-workspace"

# Determine base branch (R308)
BASE_BRANCH="phase1-integration"  # NOT main for Phase 2!

# Create infrastructure
git clone --branch "$BASE_BRANCH" "$TARGET_REPO_URL" "$INTEGRATION_WORKSPACE"
cd "$INTEGRATION_WORKSPACE"
mkdir -p .software-factory
git checkout -b "phase2-wave1-integration"
git push -u origin "phase2-wave1-integration"

# Result structure:
# /home/vscode/software-factory-template/efforts/phase2/wave1/integration-workspace/
```

### Example 2: Phase 2, Wave 1 Integration (Re-integration After Fixes)

```bash
# Previous integration failed, must redo

# Archive old workspace
mv "$INTEGRATION_WORKSPACE" "${INTEGRATION_WORKSPACE}-archived-1"

# Create fresh infrastructure (SAME paths)
git clone --branch "$BASE_BRANCH" "$TARGET_REPO_URL" "$INTEGRATION_WORKSPACE"
cd "$INTEGRATION_WORKSPACE"
mkdir -p .software-factory
git checkout -b "phase2-wave1-integration"
git push --force-with-lease -u origin "phase2-wave1-integration"  # Force-push

# Result structure:
# /home/vscode/software-factory-template/efforts/phase2/wave1/integration-workspace-archived-1/  (old)
# /home/vscode/software-factory-template/efforts/phase2/wave1/integration-workspace/             (new)
```

## 🚨 COMMON MISTAKES TO AVOID

### ❌ MISTAKE 1: Creating /repo subdirectory
```bash
# WRONG - Creates unnecessary subdirectory
mkdir -p integration-workspace
git clone "$REPO" "integration-workspace/repo"
```

### ❌ MISTAKE 2: Not archiving on re-integration
```bash
# WRONG - Overwrites previous attempt
rm -rf integration-workspace
mkdir integration-workspace
```

### ❌ MISTAKE 3: Creating versioned branches
```bash
# WRONG - Creates branch proliferation
git checkout -b "phase1-wave1-integration-v2"
```

### ❌ MISTAKE 4: Not creating .software-factory directory
```bash
# WRONG - Missing metadata directory
git clone "$REPO" "integration-workspace"
# Forgot to create .software-factory/!
```

## ✅ VALIDATION CHECKLIST

After creating integration infrastructure, verify:

1. **Directory Structure:**
   ```bash
   [ -d "${INTEGRATION_WORKSPACE}/.git" ] || echo "ERROR: Wrong structure!"
   [ -d "${INTEGRATION_WORKSPACE}/.software-factory" ] || echo "ERROR: Missing metadata directory!"
   ```

2. **No Nested Efforts:**
   ```bash
   [ ! -d "${INTEGRATION_WORKSPACE}/efforts" ] || echo "ERROR: Nested structure!"
   ```

3. **Correct Base Branch:**
   ```bash
   cd "${INTEGRATION_WORKSPACE}"
   git merge-base HEAD "origin/${BASE_BRANCH}" || echo "ERROR: Wrong base!"
   ```

4. **Archive Exists (if re-integration):**
   ```bash
   ls -d ${INTEGRATION_WORKSPACE}-archived-* 2>/dev/null || echo "First attempt"
   ```

## 📋 RULES UPDATED

The following rules have been updated with deterministic specifications:

1. **R250**: Integration Isolation Requirement
   - Added exact directory structure
   - Added re-integration archiving protocol
   - Specified repo/ subdirectory requirement

2. **R308**: Incremental Branching Strategy
   - Added re-integration handling
   - Specified force-push protocol
   - Clarified archive naming
   - Removed /repo subdirectory pattern

3. **SETUP_INTEGRATION_INFRASTRUCTURE** state rules
   - Updated with deterministic paths
   - Added archiving logic
   - Specified exact clone commands
   - Added .software-factory metadata directory

## 🔴 ENFORCEMENT

**ANY DEVIATION FROM THESE SPECIFICATIONS = -100% AUTOMATIC FAILURE**

The orchestrator MUST follow these exact specifications with NO interpretation or optimization.

---

**Remember**: DETERMINISTIC means EXACTLY THE SAME EVERY TIME. No decisions, no variations, no "close enough"!

## 🔴 METADATA ORGANIZATION

**ALL metadata MUST go in `.software-factory/` directory:**
- Implementation plans
- Review reports
- Integration reports
- Merge plans
- Work logs

**NEVER put metadata in the root directory of the workspace!**