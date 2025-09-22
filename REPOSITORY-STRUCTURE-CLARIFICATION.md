# 🔴🔴🔴 CRITICAL: REPOSITORY STRUCTURE CLARIFICATION 🔴🔴🔴

**THIS DOCUMENT PREVENTS ORCHESTRATOR CONFUSION ABOUT WHERE THINGS ARE**

## The Confusion Problem

The orchestrator has been getting confused about:
1. Looking for effort branches in Software Factory repo (they don't exist there!)
2. Not finding integration workspaces (looking in wrong places)
3. Merging stale branches (not checking where fresh branches are)
4. Not understanding which repository contains what

## Repository Structure - THE TRUTH

### 1. SOFTWARE FACTORY REPOSITORY
```
Location: /home/vscode/software-factory-template/
Purpose:  Contains Software Factory code, rules, and orchestration
Branches: 
  - main (SF template code)
  - software-factory-2.0 (SF 2.0 implementation)
  
NEVER CONTAINS:
  ❌ Effort branches (phase1/wave1/effort-name)
  ❌ Integration branches (wave1-integration, phase1-integration)
  ❌ Project implementation code
  
ALWAYS CONTAINS:
  ✅ orchestrator-state.json
  ✅ rule-library/
  ✅ agent-states/
  ✅ .claude/
```

### 2. TARGET REPOSITORY CLONES (EFFORTS)
```
Location Pattern: efforts/phase*/wave*/[effort-name]/
Purpose:         Contains actual project implementation code
Example:         efforts/phase1/wave1/image-builder/
Branches:        
  - Effort branches (e.g., idpbuilder-oci/phase1/wave1/image-builder)
  - Split branches (e.g., phase1/wave1/image-builder-split-001)
  
Each effort directory is a SEPARATE CLONE of the target repository!
These are WHERE THE CODE LIVES!
```

### 3. INTEGRATION WORKSPACES
```
Location Pattern: efforts/*/integration-workspace/
Purpose:         Isolated workspaces for merging branches
Examples:
  - efforts/phase1/wave1/integration-workspace/ (wave integration)
  - efforts/phase1/integration-workspace/ (phase integration)
  - efforts/project/integration-workspace/ (project integration)
  
Branches:
  - Integration branches (wave1-integration, phase1-integration, project-integration)
  
Contains:
  - MERGE-PLAN.md files
  - INTEGRATION-REPORT.md files
  - Merged code from efforts
```

## Navigation Guide for Orchestrator

### To Find Effort Code:
```bash
# CORRECT - Look in effort directories
cd efforts/phase1/wave1/image-builder/
git branch --show-current  # Shows effort branch

# WRONG - Looking in SF repo
cd /home/vscode/software-factory-template/
git branch -r | grep image-builder  # WON'T FIND IT!
```

### To Find Integration Workspaces:
```bash
# CORRECT - Check state file for location
WORKSPACE=$(jq -r '.current_wave_integration.workspace' orchestrator-state.json)
cd "$WORKSPACE"

# ALTERNATIVE - Search for integration workspaces
find efforts -type d -name "integration-workspace"

# WRONG - Looking in SF directory
ls -la integration-workspace/  # DOESN'T EXIST HERE!
```

### To Find Merge Plans:
```bash
# CORRECT - Look in integration workspace
cd efforts/phase1/wave1/integration-workspace/
cat WAVE-MERGE-PLAN.md

# WRONG - Looking in SF directory
cat WAVE-MERGE-PLAN.md  # NOT HERE!
```

### To Check Branch Freshness:
```bash
# CORRECT - Compare timestamps from state file
INTEGRATION_TIME=$(jq -r '.current_wave_integration.created_at' orchestrator-state.json)
EFFORT_TIME=$(jq -r '.efforts_completed[0].last_updated_at' orchestrator-state.json)

# WRONG - Assuming integration is always fresh
git merge origin/wave1-integration  # MIGHT BE STALE!
```

## State File Tracking Requirements

The orchestrator-state.json MUST track:

```json
{
  "repository_contexts": {
    "software_factory_repo": {
      "path": "/home/vscode/software-factory-template",
      "purpose": "SF orchestration and rules"
    },
    "target_repo_clones": {
      "pattern": "efforts/*/*/",
      "purpose": "Project implementation code"
    }
  },
  "integration_workspaces": {
    "wave1": {
      "workspace": "/full/path/to/efforts/phase1/wave1/integration-workspace",
      "merge_plan": "/full/path/to/efforts/phase1/wave1/integration-workspace/WAVE-MERGE-PLAN.md"
    }
  },
  "efforts_completed": [{
    "workspace": "/full/path/to/efforts/phase1/wave1/effort-name",
    "branch": "project/phase1/wave1/effort-name",
    "last_updated_at": "2025-08-25T09:00:00Z"
  }]
}
```

## Common Orchestrator Mistakes

### Mistake 1: Wrong Repository
```bash
# WRONG
cd /home/vscode/software-factory-template
git checkout phase1/wave1/image-builder  # DOESN'T EXIST!

# CORRECT
cd efforts/phase1/wave1/image-builder/
git checkout phase1/wave1/image-builder  # EXISTS IN TARGET REPO CLONE
```

### Mistake 2: Can't Find Merge Plan
```bash
# WRONG
ls -la *.md  # Looking in SF directory

# CORRECT
cd $(jq -r '.integration_workspaces.wave1.workspace' orchestrator-state.json)
ls -la *MERGE-PLAN.md
```

### Mistake 3: Merging Stale Branches
```bash
# WRONG
git merge origin/project-integration  # Could be outdated!

# CORRECT
# First check if efforts have newer commits
if check_integration_freshness; then
    git merge origin/project-integration
else
    echo "Must use fresh effort branches instead!"
fi
```

## Quick Reference Commands

### List All Effort Branches
```bash
for effort in efforts/*/*/; do
    [[ -d "$effort/.git" ]] || continue
    echo "Effort: $effort"
    (cd "$effort" && git branch --show-current)
done
```

### Find All Integration Workspaces
```bash
find efforts -type d -name "integration-workspace" -exec echo {} \;
```

### Check All Merge Plans
```bash
find efforts -name "*MERGE-PLAN.md" -exec echo "Found: {}" \;
```

### Verify Repository Types
```bash
# Check if directory is SF or target repo
check_repo_type() {
    local dir="$1"
    cd "$dir"
    REMOTE=$(git remote get-url origin)
    if [[ "$REMOTE" == *"software-factory"* ]]; then
        echo "Software Factory repository"
    else
        echo "Target project repository"
    fi
}
```

## Summary Rules

1. **Effort branches** → In effort directories (target repo clones)
2. **Integration branches** → In integration workspaces
3. **Merge plans** → In integration workspaces
4. **State files** → In SF repository root
5. **Rules** → In SF repository rule-library/

**NEVER look for project branches in the Software Factory repository!**
**ALWAYS check the state file for workspace locations!**
**ALWAYS verify branch freshness before merging!**