# 🔴🔴🔴 RULE R250 - INTEGRATION MUST USE SEPARATE TARGET REPO CLONE 🔴🔴🔴

## ⚠️⚠️⚠️ CRITICAL: NEVER INTEGRATE IN THE SF INSTANCE DIRECTORY! ⚠️⚠️⚠️

## THE ABSOLUTE LAW:

**INTEGRATION HAPPENS IN THE TARGET REPOSITORY, NOT THE SOFTWARE FACTORY!**

The Software Factory instance directory is for:
- Planning
- Configuration
- State management
- Agent coordination

The Software Factory instance directory is NOT for:
- ❌ Source code
- ❌ Integration branches
- ❌ Merging code
- ❌ Running tests

## MANDATORY INTEGRATION PROTOCOL:

```bash
perform_wave_integration() {
    local PHASE="$1"
    local WAVE="$2"
    
    # 0. SOURCE BRANCH NAMING HELPERS AND GET PREFIX
    SF_INSTANCE_DIR=$(pwd)
    source "$SF_INSTANCE_DIR/utilities/branch-naming-helpers.sh"
    PROJECT_PREFIX=$(jq '.branch_naming.project_prefix' "$SF_INSTANCE_DIR/target-repo-config.yaml")
    
    # 1. CREATE INTEGRATION DIRECTORY UNDER /efforts/ STRUCTURE
    WAVE_DIR="/efforts/phase${PHASE}/wave${WAVE}"
    INTEGRATION_DIR="${WAVE_DIR}/integration-workspace"
    echo "Creating integration workspace at: $INTEGRATION_DIR"
    mkdir -p "$INTEGRATION_DIR"
    
    # 2. FRESH CLONE OF TARGET REPO (NOT SF INSTANCE!)
    cd "$INTEGRATION_DIR"
    TARGET_REPO_URL=$(jq '.target_repository.url' "$SF_INSTANCE_DIR/target-repo-config.yaml")
    git clone "$TARGET_REPO_URL" .
    
    # 3. CREATE INTEGRATION BRANCH WITH PROPER PREFIX
    INTEGRATION_BRANCH=$(get_wave_integration_branch_name "$PHASE" "$WAVE")
    echo "Creating integration branch: $INTEGRATION_BRANCH"
    # Example with prefix: tmc-workspace/phase1/wave1-integration
    # Example without: phase1/wave1-integration
    git checkout -b "$INTEGRATION_BRANCH"
    
    # 4. FETCH ALL EFFORT BRANCHES FROM REMOTES
    for effort_dir in /efforts/phase${PHASE}/wave${WAVE}/*/; do
        # Skip the integration-workspace directory itself
        [[ "$effort_dir" == *"integration-workspace"* ]] && continue
        effort_name=$(basename "$effort_dir")
        
        # Get properly formatted effort branch with prefix
        EFFORT_BRANCH=$(get_effort_branch_name "$PHASE" "$WAVE" "$effort_name")
        
        echo "Fetching $effort_name (branch: $EFFORT_BRANCH)"
        git fetch origin "$EFFORT_BRANCH"
    done
    
    # 5. MERGE IN THE TARGET REPO CLONE
    for effort_dir in /efforts/phase${PHASE}/wave${WAVE}/*/; do
        # Skip the integration-workspace directory itself
        [[ "$effort_dir" == *"integration-workspace"* ]] && continue
        effort_name=$(basename "$effort_dir")
        
        # Get properly formatted effort branch with prefix
        EFFORT_BRANCH=$(get_effort_branch_name "$PHASE" "$WAVE" "$effort_name")
        
        echo "Merging $effort_name from $EFFORT_BRANCH"
        git merge "origin/$EFFORT_BRANCH" --no-ff -m "integrate: $effort_name into wave${WAVE}"
    done
    
    # 6. PUSH INTEGRATION BRANCH TO TARGET REPO
    git push -u origin "$INTEGRATION_BRANCH"
    
    # 7. RETURN TO SF INSTANCE
    cd "$SF_INSTANCE_DIR"
}
```

## DIRECTORY STRUCTURE:

```
software-factory-instance/          # THIS IS THE SF INSTANCE (orchestrator lives here)
├── orchestrator-state.json         # State management
├── target-repo-config.yaml         # Config pointing to TARGET repo
└── utilities/                      # Helper scripts

/efforts/                           # SEPARATE ROOT DIRECTORY for all code work
├── phase1/
│   ├── wave1/
│   │   ├── core-api/              # Clone of TARGET repo for effort
│   │   ├── webhooks/              # Clone of TARGET repo for effort
│   │   └── integration-workspace/ # Clone of TARGET repo for wave integration
│   ├── wave2/
│   │   ├── controllers/           # Clone of TARGET repo for effort
│   │   └── integration-workspace/ # Clone of TARGET repo for wave integration
│   └── phase-integration-workspace/ # Clone of TARGET repo for phase integration
└── phase2/
    └── wave1/
        ├── rbac/                  # Clone of TARGET repo for effort
        └── integration-workspace/ # Clone of TARGET repo for wave integration

TARGET-REPOSITORY/                  # THIS IS THE ACTUAL TARGET (separate repo!)
├── pkg/                           # Actual source code
├── cmd/                           # Actual applications
└── test/                          # Actual tests
```

## 🚫 FORBIDDEN ACTIONS IN SF INSTANCE:

```bash
# ❌ NEVER DO THIS IN THE SF INSTANCE DIRECTORY:
git checkout -b "integration"      # NO! Wrong repo!
git merge effort-branch            # NO! Wrong repo!
go test ./...                      # NO! Wrong repo!
make build                         # NO! Wrong repo!
```

## ✅ CORRECT INTEGRATION LOCATIONS:

```bash
# ✅ WAVE INTEGRATION - Under wave directory in /efforts/:
cd /efforts/phase1/wave1/integration-workspace     # Correct location
git checkout -b "tmc-workspace/phase1/wave1-integration"  # In TARGET repo
git merge effort-branches                          # In TARGET repo
go test ./...                                      # In TARGET repo

# ✅ PHASE INTEGRATION - Under phase directory in efforts/:
cd ${CLAUDE_PROJECT_DIR}/efforts/phase1/phase-integration-workspace     # Correct location
git checkout -b "tmc-workspace/phase1-integration" # In TARGET repo
git merge wave-integration-branches                # In TARGET repo
go test ./...                                      # In TARGET repo
```

## VALIDATION CHECKLIST:

Before starting integration, verify:
- [ ] You are NOT in the software-factory instance directory
- [ ] You have created a NEW clone of the TARGET repository
- [ ] Wave integration clone is in `${CLAUDE_PROJECT_DIR}/efforts/phaseX/waveY/integration-workspace/`
- [ ] Phase integration clone is in `${CLAUDE_PROJECT_DIR}/efforts/phaseX/phase-integration-workspace/`
- [ ] You are working with the TARGET repository URL from config
- [ ] You are NOT trying to merge in the SF instance

## ERROR MESSAGES TO WATCH FOR:

If you see these, you're in the WRONG place:
- "fatal: not a git repository" - You're in SF instance, not target
- "orchestrator-state.json" in git status - Wrong repo!
- ".claude/" directory visible - Wrong repo!
- "setup-config.yaml" visible - Wrong repo!

## THE GOLDEN RULE:

**The Software Factory instance orchestrates work on the target repository.**
**The Software Factory instance NEVER contains target repository code.**

## GRADING IMPACT:

**AUTOMATIC FAILURE** if integration is attempted in the SF instance directory.

This is equivalent to trying to merge code into your CI/CD system instead of your actual project!