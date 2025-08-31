# 🚀 ORCHESTRATOR WORKSPACE SETUP - QUICK REFERENCE

## 🚨🚨🚨 CRITICAL WORKSPACE RULES 🚨🚨🚨

### MANDATORY: Before Spawning ANY Agent

```bash
# R181: Orchestrator MUST set up complete workspace
# R182: Each effort gets sparse git clone  
# R183: Create proper base branch
# R184: Use phase/wave/effort-{name} naming
# R185: Verify before spawning agents
```

## ⚡ QUICK SETUP COMMANDS

### 1. Setup Single Effort Workspace
```bash
# Use the utility script (RECOMMENDED)
./utilities/setup-effort-workspace.sh 1 2 core-types https://github.com/idpbuilder/idpbuilder.git main

# Or manually:
EFFORT_DIR="efforts/phase1/wave2/core-types"
mkdir -p "$EFFORT_DIR" && cd "$EFFORT_DIR"
git clone --sparse --filter=blob:none https://github.com/idpbuilder/idpbuilder.git .
git sparse-checkout init --cone
git sparse-checkout set pkg/
git checkout main && git pull origin main
git checkout -b "phase1/wave2/effort-core-types"
touch IMPLEMENTATION-PLAN.md work-log.md
mkdir -p pkg
```

### 2. Spawn Agent with MANDATORY Working Directory
```bash
# ✅ CORRECT - ALWAYS DO THIS:
Task: sw-engineer
Working directory: efforts/phase1/wave2/core-types  # MANDATORY!
Instructions: Implement according to IMPLEMENTATION-PLAN.md

# ❌ WRONG - IMMEDIATE FAILURE:
Task: sw-engineer
# No working directory = GRADING FAILURE
```

### 3. Verify After Spawning
```bash
# Use verification script
./utilities/post-spawn-verify.sh 1 2

# Or quick manual check
for effort in efforts/phase*/wave*/*/; do
    if [ -d "$effort/pkg" ]; then
        echo "✅ $(basename $effort): Isolated workspace"
    else
        echo "❌ $(basename $effort): VIOLATION!"
    fi
done
```

## 📋 WORKSPACE STRUCTURE

```
efforts/
├── phase1/
│   ├── wave1/
│   │   ├── core-types/          ← Each effort
│   │   │   ├── .git/           ← Own git repo (sparse)
│   │   │   ├── IMPLEMENTATION-PLAN.md
│   │   │   ├── work-log.md
│   │   │   └── pkg/            ← ALL CODE HERE
│   │   │       └── oci/types/
│   │   ├── api-contracts/      ← Separate effort
│   │   │   ├── .git/
│   │   │   └── pkg/
│   │   └── controller-base/    ← Separate effort
│   │       ├── .git/
│   │       └── pkg/
```

## 🔴 GRADING FAILURES

### Automatic -20% Workspace Isolation:
- ❌ Agent working in main `/pkg/`
- ❌ No git repository in effort directory
- ❌ Wrong branch naming pattern
- ❌ No working directory specified when spawning

### How to Avoid:
1. ALWAYS use setup-effort-workspace.sh
2. ALWAYS specify working directory when spawning
3. ALWAYS verify with post-spawn-verify.sh
4. NEVER let agents cd to other directories

## 🎯 SPARSE CHECKOUT PATTERNS

### By Effort Type:
```bash
# API/Types efforts
git sparse-checkout set pkg/api pkg/apis pkg/types

# Controller efforts  
git sparse-checkout set pkg/controller pkg/controllers

# Webhook efforts
git sparse-checkout set pkg/webhook pkg/webhooks

# Test efforts
git sparse-checkout set test/ tests/ e2e/

# Default (most efforts)
git sparse-checkout set pkg/
```

## ⚠️ COMMON MISTAKES

### 1. Forgetting Working Directory
```bash
# ❌ WRONG
Task: sw-engineer
Instructions: Implement feature

# ✅ CORRECT
Task: sw-engineer
Working directory: efforts/phase1/wave2/feature
Instructions: Implement feature
```

### 2. Not Creating Git Workspace
```bash
# ❌ WRONG - Just mkdir
mkdir -p efforts/phase1/wave2/feature

# ✅ CORRECT - Full git setup
./utilities/setup-effort-workspace.sh 1 2 feature
```

### 3. Wrong Branch Name
```bash
# ❌ WRONG
git checkout -b phase1-wave2-feature
git checkout -b effort/feature
git checkout -b phase1/wave2/feature  # Missing "effort-"

# ✅ CORRECT
git checkout -b phase1/wave2/effort-feature
```

## 📊 VERIFICATION CHECKLIST

Before spawning agents:
- [ ] Effort directory created
- [ ] Sparse git clone completed
- [ ] Correct branch created (phase*/wave*/effort-*)
- [ ] IMPLEMENTATION-PLAN.md exists
- [ ] work-log.md exists
- [ ] pkg/ directory created
- [ ] Remote origin configured

After spawning agents:
- [ ] Agent acknowledged correct directory
- [ ] Agent on correct branch
- [ ] Code being created in effort/pkg/
- [ ] Main /pkg/ remains empty
- [ ] No workspace violations

## 🚀 PARALLEL SPAWN TEMPLATE

```bash
# Set up all workspaces FIRST
./utilities/setup-effort-workspace.sh 1 2 effort1 $REPO main
./utilities/setup-effort-workspace.sh 1 2 effort2 $REPO main
./utilities/setup-effort-workspace.sh 1 2 effort3 $REPO main

# Then spawn all agents in ONE message
Task: sw-engineer
Working directory: efforts/phase1/wave2/effort1
Instructions: Implement effort1

Task: sw-engineer  
Working directory: efforts/phase1/wave2/effort2
Instructions: Implement effort2

Task: sw-engineer
Working directory: efforts/phase1/wave2/effort3
Instructions: Implement effort3

# Verify all are working correctly
./utilities/post-spawn-verify.sh 1 2
```

## 📝 RULES REFERENCE

| Rule | Description | Criticality |
|------|-------------|-------------|
| R181 | Orchestrator Workspace Setup Responsibility | BLOCKING 🚨🚨🚨 |
| R182 | Sparse Clone Requirement | CRITICAL 🚨 |
| R183 | Branch Creation Protocol | CRITICAL 🚨 |
| R184 | Effort Branch Naming Scheme | MANDATORY 🚨🚨 |
| R185 | Workspace Verification Before Spawn | CRITICAL 🚨 |
| R176 | Workspace Isolation Requirement | BLOCKING 🚨🚨🚨 |
| R177 | Agent Working Directory Enforcement | CRITICAL 🚨 |

## 🔗 UTILITY SCRIPTS

- `utilities/setup-effort-workspace.sh` - Create complete effort workspace
- `utilities/post-spawn-verify.sh` - Verify workspace isolation
- `utilities/verify-workspace-isolation.sh` - Check single workspace

---

**Remember**: Workspace isolation is 20% of grading. One violation = automatic failure!