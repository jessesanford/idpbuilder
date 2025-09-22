# 🔴🔴🔴 CRITICAL: TARGET REPOSITORY CONFIGURATION 🔴🔴🔴

## THE MOST COMMON AND CATASTROPHIC MISTAKE IN SOFTWARE FACTORY 2.0

### THE PROBLEM THAT CAUSES 90% OF FAILURES

**Software Factory has TWO completely separate repositories:**

1. **SOFTWARE FACTORY REPOSITORY** (The Framework)
   - Location: `/workspaces/your-sf-instance/`
   - Contains: Rules, agents, state files, planning docs
   - Purpose: ORCHESTRATION AND MANAGEMENT ONLY
   - **NEVER WRITE CODE HERE!**

2. **TARGET REPOSITORY** (Your Actual Project)
   - Location: Gets cloned into `/workspaces/your-sf-instance/efforts/`
   - Contains: Your actual project code
   - Purpose: WHERE ALL CODE IMPLEMENTATION HAPPENS
   - **THIS IS CONFIGURED IN target-repo-config.yaml**

## 🚨 WITHOUT target-repo-config.yaml YOU WILL FAIL 🚨

The orchestrator CANNOT proceed without knowing:
- **WHAT** repository to clone
- **WHERE** to branch from
- **HOW** to name branches

## MANDATORY CONFIGURATION FILE

### Location
```
/workspaces/your-sf-instance/target-repo-config.yaml
```

### Complete Example Configuration
```yaml
# target-repo-config.yaml
# THIS FILE IS MANDATORY - WITHOUT IT, NOTHING WORKS!

# 1. THE TARGET REPOSITORY (Your actual project)
target_repository:
  # The URL of YOUR PROJECT repository (NOT the SF repo!)
  url: "https://github.com/your-org/your-actual-project.git"
  
  # The default base branch to create efforts from
  base_branch: "main"
  
  # Optional: Credentials method if private repo
  auth_method: "ssh"  # or "https" with token

# 2. BRANCH NAMING CONVENTIONS
branch_naming:
  # Optional project prefix for all branches
  project_prefix: "my-project"  # or null for no prefix
  
  # Format for effort branches
  # Variables: {prefix}, {phase}, {wave}, {effort_name}
  effort_format: "{prefix}/phase{phase}/wave{wave}/{effort_name}"
  
  # Format for wave integration branches
  integration_format: "{prefix}/phase{phase}/wave{wave}/integration"
  
  # Format for phase integration branches
  phase_integration_format: "{prefix}/phase{phase}/integration"

# 3. WORKSPACE CONFIGURATION
workspace:
  # Root directory for all efforts (relative to SF root)
  efforts_root: "efforts"
  
  # Path pattern for individual efforts
  effort_path: "phase{phase}/wave{wave}/{effort_name}"
  
  # Optional: Max parallel efforts
  max_parallel_efforts: 3
```

## COMMON MISTAKES AND THEIR CONSEQUENCES

### Mistake 1: No target-repo-config.yaml
```bash
# Orchestrator tries to proceed without config
cd /workspaces/my-sf
ls target-repo-config.yaml
# ls: cannot access 'target-repo-config.yaml': No such file or directory

# RESULT: Orchestrator FAILS immediately
# ERROR: "CRITICAL: target-repo-config.yaml not found!"
# GRADE: -100% (Cannot even start)
```

### Mistake 2: Using SF Repo URL as Target
```yaml
# WRONG - Points to Software Factory itself!
target_repository:
  url: "https://github.com/user/software-factory-instance.git"  # WRONG!
```

**CONSEQUENCE:** You'll modify the Software Factory framework instead of your project!

### Mistake 3: Creating Code in SF Repository
```bash
# WRONG - Creating code in SF repo
cd /workspaces/my-sf-instance
mkdir src
echo "package main" > src/main.go  # VIOLATION!
```

**CONSEQUENCE:** 
- Code in wrong repository
- Cannot be integrated
- Complete failure of implementation

### Mistake 4: Cloning SF Repo into Efforts
```bash
# WRONG - Cloning SF into efforts
cd /workspaces/my-sf-instance/efforts/phase1/wave1/my-effort
git clone https://github.com/user/software-factory-instance.git .
```

**CONSEQUENCE:** Infinite recursion of SF instances!

## CORRECT WORKFLOW

### Step 1: Create Configuration
```bash
cd /workspaces/your-sf-instance
cat > target-repo-config.yaml << 'EOF'
target_repository:
  url: "$PROJECT-TARGET-REPO"  # Set from setup-config.yaml or environment
  base_branch: "main"

branch_naming:
  project_prefix: "your-project-prefix"
  effort_format: "{prefix}/phase{phase}/wave{wave}/{effort_name}"

workspace:
  efforts_root: "efforts"
EOF
```

**Note:** The `$PROJECT-TARGET-REPO` variable should be defined in:
- `setup-config.yaml` - For initial setup
- `target-repo-config.yaml` - For runtime configuration
- Environment variable - For CI/CD or automation

### Step 2: Orchestrator Validates
```bash
# Orchestrator on startup
if [ ! -f "target-repo-config.yaml" ]; then
    echo "CRITICAL: Cannot proceed without target config!"
    exit 191
fi

TARGET_URL=$(yq '.target_repository.url' target-repo-config.yaml)
echo "Target repository: $TARGET_URL"
```

### Step 3: Infrastructure Setup Uses Config
```bash
# In SETUP_EFFORT_INFRASTRUCTURE state
TARGET_REPO_URL=$(yq '.target_repository.url' target-repo-config.yaml)

# Clone TARGET repo (not SF repo!)
git clone --single-branch --branch main "$TARGET_REPO_URL" \
    /workspaces/my-sf/efforts/phase1/wave1/api-effort
```

### Step 4: Agents Work in Clones
```bash
# SW Engineer works in cloned target repo
cd /workspaces/my-sf/efforts/phase1/wave1/api-effort
# This is a clone of YOUR PROJECT
echo "package api" > pkg/api/types.go  # CORRECT!
git add .
git commit -m "feat: implement API types"
git push
```

## VALIDATION CHECKLIST

### For Orchestrators
- [ ] Check target-repo-config.yaml exists before ANY work
- [ ] Validate target URL is NOT the SF repository
- [ ] Load configuration before SETUP_EFFORT_INFRASTRUCTURE
- [ ] Pass configuration to all spawned agents
- [ ] Use target URL for all cloning operations

### For SW Engineers
- [ ] Verify working directory is under /efforts/
- [ ] Confirm repository is target project (not SF)
- [ ] Check git remote points to target URL
- [ ] Never create code outside effort directory

### For Code Reviewers
- [ ] Verify code location is in effort clones
- [ ] Check branches exist in target repository
- [ ] Ensure reviews happen on target project code

## ERROR MESSAGES AND SOLUTIONS

### "target-repo-config.yaml not found"
**Solution:** Create the configuration file with your project's repository URL

### "Target repository URL same as SF repository"
**Solution:** Update URL to point to your actual project, not the SF instance

### "No target_repository.url in config"
**Solution:** Add the target_repository.url field to your config

### "Cannot clone repository"
**Solution:** Verify URL is correct and accessible, check authentication

## ARCHITECTURE DIAGRAM

```
/workspaces/
├── my-software-factory/           # SOFTWARE FACTORY INSTANCE
│   ├── target-repo-config.yaml    # POINTS TO TARGET REPO ←──────┐
│   ├── rule-library/              # Framework rules              │
│   ├── .claude/agents/            # Agent configurations         │
│   ├── orchestrator-state.json    # State management             │
│   └── efforts/                   # CLONED TARGET REPOS          │
│       ├── phase1/                                               │
│       │   ├── wave1/                                            │
│       │   │   ├── api-effort/    # CLONE of target ────────────┤
│       │   │   │   ├── .git/      # Points to TARGET repo       │
│       │   │   │   ├── pkg/       # TARGET project code         │
│       │   │   │   └── tests/     # TARGET project tests        │
│       │   │   └── auth-effort/   # CLONE of target ────────────┤
│       │   │       └── ...        # Points to TARGET repo       │
│       │   └── wave2/                                            │
│       │       └── ...                                           │
│       └── phase2/                                               │
│           └── ...                                               │
│                                                                 │
└── target-project/                # SEPARATE TARGET REPO ←───────┘
    ├── .git/                      # Original project
    ├── pkg/                       # Original code
    └── ...                        # Gets cloned into efforts
```

## RULE REFERENCES

- **R191**: Target Repository Configuration (rule-library/R191-target-repo-config.md)
- **R192**: Repository Separation (rule-library/R192-repo-separation.md)
- **R193**: Effort Clone Protocol (rule-library/R193-effort-clone-protocol.md)
- **R251**: Repository Separation Law (rule-library/R251-REPOSITORY-SEPARATION-LAW.md)

## FINAL WARNINGS

1. **The #1 cause of Software Factory failures is working in the wrong repository**
2. **Without target-repo-config.yaml, the orchestrator CANNOT function**
3. **The target repository URL must NEVER be the SF repository itself**
4. **All code must be written in effort clones, not the SF instance**
5. **This is not optional - it's the foundation of the entire system**

---

**Remember:** The Software Factory is a MANAGEMENT FRAMEWORK. Your actual project is configured in target-repo-config.yaml and cloned into efforts/ for implementation.