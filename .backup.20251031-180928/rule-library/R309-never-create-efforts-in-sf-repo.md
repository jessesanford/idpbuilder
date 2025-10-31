# 🔴🔴🔴 RULE R309 - NEVER Create Efforts in SF Repo (PARAMOUNT LAW)

## Rule Definition
**Criticality:** PARAMOUNT - Automatic -100% failure for violation
**Category:** Infrastructure
**Applies To:** orchestrator, all infrastructure setup
**Created:** 2025-09-02 (to stop critical pollution issue)

## 🔴🔴🔴 THE FUNDAMENTAL SEPARATION 🔴🔴🔴

### TWO COMPLETELY DIFFERENT REPOSITORIES:

1. **SOFTWARE FACTORY REPO** (Where you are now)
   - Path: `/home/vscode/software-factory-template/` or similar
   - Contains: Rules, agents, state files, orchestration
   - Has: `.claude/`, `rule-library/`, `orchestrator-state-v3.json`
   - Purpose: PLANNING AND ORCHESTRATION ONLY
   - **NEVER CREATE CODE HERE!**

2. **TARGET REPOSITORY** (What you clone for efforts)
   - URL: From `target-repo-config.yaml`
   - Contains: ACTUAL PROJECT CODE
   - Gets cloned to: `efforts/phaseX/waveY/effort-name/`
   - Purpose: WHERE ALL IMPLEMENTATION HAPPENS
   - **THIS IS WHERE EFFORT BRANCHES GO!**

## 🔴🔴🔴 ABSOLUTE PROHIBITIONS 🔴🔴🔴

### NEVER DO THIS (AUTOMATIC -100% FAILURE):

```bash
# ❌❌❌ FATAL ERROR - Creating branch in SF repo
cd /home/vscode/software-factory-template
git checkout -b phase1/wave1/effort-name  # DESTROYS EVERYTHING!

# ❌❌❌ FATAL ERROR - Using SF repo as target
TARGET_REPO_URL="https://github.com/user/software-factory-template"
git clone "$TARGET_REPO_URL" efforts/...  # WRONG REPO!

# ❌❌❌ FATAL ERROR - Creating code in SF repo
cd /home/vscode/software-factory-template
mkdir pkg/api  # NO! This is orchestration repo!
vim main.go    # NO! Code goes in TARGET repo clones!
```

### ALWAYS DO THIS (ONLY CORRECT WAY):

```bash
# ✅ CORRECT - Clone TARGET repo for efforts
TARGET_URL=$(yq '.target_repository.url' target-repo-config.yaml)
# Verify it's NOT the SF repo
if [[ "$TARGET_URL" == *"software-factory"* ]]; then
    echo "🔴 FATAL: Target is SF template! Fix config!"
    exit 1
fi
git clone "$TARGET_URL" efforts/phase1/wave1/effort-name

# ✅ CORRECT - Create branches in TARGET clone
cd efforts/phase1/wave1/effort-name
git checkout -b phase1/wave1/effort-name  # In TARGET clone!

# ✅ CORRECT - Implementation in TARGET clone
cd efforts/phase1/wave1/effort-name
mkdir pkg/api  # In TARGET repo clone
vim main.go    # In TARGET repo clone
```

## 🚨 DETECTION AND ENFORCEMENT

### Pre-Branch Validation (MANDATORY):
```bash
validate_not_in_sf_repo() {
    # Check 1: SF marker files
    if [ -f ".claude/CLAUDE.md" ] || [ -f "rule-library/RULE-REGISTRY.md" ]; then
        echo "🔴🔴🔴 FATAL ERROR: YOU ARE IN SF REPO!"
        echo "NEVER CREATE EFFORT BRANCHES HERE!"
        echo "Clone TARGET repo to efforts/ instead!"
        exit 309
    fi
    
    # Check 2: Remote URL
    local REMOTE=$(git remote get-url origin 2>/dev/null)
    if [[ "$REMOTE" == *"software-factory"* ]]; then
        echo "🔴🔴🔴 FATAL: This is Software Factory repo!"
        echo "Efforts go in TARGET repository clones!"
        exit 309
    fi
    
    # Check 3: Directory path
    if [[ "$PWD" == *"software-factory-template"* ]] && [[ "$PWD" != *"/efforts/"* ]]; then
        echo "🔴🔴🔴 FATAL: In SF root, not effort directory!"
        exit 309
    fi
}

# CALL THIS BEFORE EVERY git checkout -b
```

### Continuous Monitoring:
```bash
monitor_sf_repo_pollution() {
    cd $CLAUDE_PROJECT_DIR
    
    # Check for effort branches
    POLLUTION=$(git branch -a | grep -E "effort|wave|split" | grep -v main)
    
    if [[ -n "$POLLUTION" ]]; then
        echo "🔴🔴🔴 CRITICAL POLLUTION DETECTED!"
        echo "Effort branches found in SF repo:"
        echo "$POLLUTION"
        echo ""
        echo "THESE MUST BE DELETED IMMEDIATELY!"
        echo "Efforts belong in TARGET repo clones!"
        echo ""
        echo "To fix:"
        echo "1. Delete these branches"
        echo "2. Clone TARGET repo to efforts/"
        echo "3. Create branches there instead"
        exit 309
    fi
}
```

## 🔴 WHY THIS MATTERS

### Pollution Consequences:
1. **Corrupts Planning Repository** - SF repo becomes unusable
2. **Breaks State Machine** - Can't track actual vs plan
3. **Destroys Isolation** - Efforts contaminate each other
4. **Makes Integration Impossible** - Can't merge from wrong repo
5. **Violates Architecture** - Fundamental system design broken

### Correct Architecture:
```
software-factory-template/          # SF REPO (orchestration)
├── .claude/                       # Agent configs
├── rule-library/                  # Rules
├── orchestrator-state-v3.json        # State tracking
├── target-repo-config.yaml        # Points to TARGET
└── efforts/                       # Clone TARGET here
    └── phase1/
        └── wave1/
            ├── effort1/           # TARGET REPO CLONE
            │   ├── .git/          # Separate git repo
            │   ├── pkg/           # Actual code
            │   └── main.go        # Implementation
            └── effort2/           # ANOTHER TARGET CLONE
                ├── .git/          # Separate git repo
                └── pkg/           # Actual code
```

## 🚨 GRADING IMPACT

### Violations and Penalties:
- Creating ANY branch in SF repo: -100% IMMEDIATE FAILURE
- Cloning SF repo into efforts/: -100% IMMEDIATE FAILURE  
- Writing code in SF repo: -100% IMMEDIATE FAILURE
- Not validating before branch creation: -50%
- Missing pollution checks: -30%

### Required Validations:
- ✅ Before EVERY branch creation: validate_not_in_sf_repo()
- ✅ Before EVERY clone: verify TARGET_URL != SF repo
- ✅ After setup: monitor_sf_repo_pollution()
- ✅ In spawn commands: verify effort directory is TARGET clone

## 📋 ORCHESTRATOR CHECKLIST

Before creating ANY infrastructure:
- [ ] Read target-repo-config.yaml
- [ ] Verify TARGET_URL is NOT software-factory repo
- [ ] Understand SF repo = planning, TARGET repo = code
- [ ] Never run git checkout -b in SF root
- [ ] Always clone TARGET to efforts/
- [ ] Always cd to effort directory before branches
- [ ] Run validate_not_in_sf_repo() before branches
- [ ] Run monitor_sf_repo_pollution() after setup

## 🔴 ENFORCEMENT IN ALL STATES

### In CREATE_NEXT_INFRASTRUCTURE:
```bash
# FIRST thing - verify we understand the separation
echo "🔴 R309 CHECK: Verifying repository understanding..."
echo "SF Repo (planning): $(pwd)"
echo "Target Repo (code): $(yq '.target_repository.url' target-repo-config.yaml)"

if [[ "$(pwd)" == *"/efforts/"* ]]; then
    echo "❌ ERROR: Already in efforts! Return to SF root!"
    cd $CLAUDE_PROJECT_DIR
fi

# NEVER create branches here!
validate_not_in_sf_repo
```

### In SPAWN_SW_ENGINEERS:
```bash
# Verify agent will work in TARGET clone
EFFORT_DIR="efforts/phase${PHASE}/wave${WAVE}/${EFFORT}"
if [ ! -d "$EFFORT_DIR/.git" ]; then
    echo "❌ No git repo in effort directory!"
    exit 1
fi

cd "$EFFORT_DIR"
validate_not_in_sf_repo  # Should pass - we're in clone
```

## 🔴 RECOVERY PROCEDURE

If you've already polluted the SF repo:

```bash
# 1. Document the pollution
cd $CLAUDE_PROJECT_DIR
git branch -a | grep -E "effort|wave|split" > pollution.txt

# 2. Delete effort branches from SF repo
while read branch; do
    git branch -D "$branch" 2>/dev/null
    git push origin --delete "$branch" 2>/dev/null
done < pollution.txt

# 3. Clean up any effort code in SF repo
find . -path "./efforts" -prune -o -name "*.go" -print -delete
find . -path "./efforts" -prune -o -name "*.py" -print -delete

# 4. Recreate efforts in TARGET clones
TARGET_URL=$(yq '.target_repository.url' target-repo-config.yaml)
for effort in effort1 effort2; do
    git clone "$TARGET_URL" "efforts/phase1/wave1/$effort"
    cd "efforts/phase1/wave1/$effort"
    git checkout -b "phase1/wave1/$effort"
    cd $CLAUDE_PROJECT_DIR
done

echo "✅ Pollution cleaned, efforts in correct location"
```

---

**REMEMBER:** 
- SF REPO = PLANNING ONLY (where orchestrator-state-v3.json lives)
- TARGET REPO = CODE ONLY (what gets cloned to efforts/)
- NEVER MIX THEM UP!
- VIOLATION = -100% AUTOMATIC FAILURE!