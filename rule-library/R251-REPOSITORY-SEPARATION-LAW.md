# 🔴🔴🔴 RULE R251 - UNIVERSAL REPOSITORY SEPARATION LAW 🔴🔴🔴

## ⚠️⚠️⚠️ THIS APPLIES TO ALL AGENTS - NO EXCEPTIONS ⚠️⚠️⚠️

## SECTION 1: THE ABSOLUTE LAW OF REPOSITORY SEPARATION

### TWO DISTINCT REPOSITORIES EXIST:

1. **SOFTWARE FACTORY INSTANCE** (Planning & Orchestration)
   - Location: The root directory with orchestrator-state-v3.json
   - Purpose: Planning, configuration, state management, agent coordination
   - Contains: .claude/, orchestrator-state-v3.json, target-repo-config.yaml, phase-plans/
   - NEVER contains: Source code, tests, application logic

2. **TARGET REPOSITORY** (Actual Software Development)
   - Location: ALWAYS under `/efforts/` directory (multiple clones)
   - Purpose: Actual software development, code implementation, testing
   - Contains: pkg/, cmd/, test/, go.mod, Makefile, actual source code
   - NEVER contains: orchestrator-state-v3.json, phase plans, agent configs

## SECTION 2: MANDATORY PRE-FLIGHT CHECK FOR ALL AGENTS

```bash
# EVERY AGENT MUST ACKNOWLEDGE THIS ON STARTUP:
echo "🔴 R251: REPOSITORY SEPARATION ACKNOWLEDGMENT"
echo "════════════════════════════════════════════════════════"
echo "I understand:"
echo "1. SOFTWARE FACTORY = Planning only (NO CODE)"
echo "2. TARGET REPO = Code only (under /efforts/)"
echo "3. I will NEVER write code in the SF instance"
echo "4. I will ONLY write code in target repo clones"
echo "════════════════════════════════════════════════════════"

# Verify current location
if [ -f "orchestrator-state-v3.json" ]; then
    echo "📍 Currently in: SOFTWARE FACTORY INSTANCE"
    echo "⚠️  This is for PLANNING ONLY - NO CODE HERE!"
elif [ -f "go.mod" ] || [ -f "Makefile" ]; then
    echo "📍 Currently in: TARGET REPOSITORY CLONE"
    echo "✅ This is where code development happens"
else
    echo "❓ Location unclear - checking..."
    pwd
fi
```

## SECTION 3: INTEGRATE_WAVE_EFFORTS TESTING PROTOCOL (CRITICAL UPDATE)

### The Complete Software Lifecycle:

1. **DEVELOPMENT PHASE**: Code in `/efforts/` (isolated)
   - Each effort in its own target repo clone
   - Complete separation from SF instance
   
2. **INTEGRATE_WAVE_EFFORTS TESTING PHASE**: Validate in `integration-testing` branch
   - Create integration-testing branch from main's HEAD
   - Merge all efforts sequentially
   - Prove everything builds and runs
   - **NEVER touch main branch directly**
   
3. **PR PLAN PHASE**: Generate MASTER-PR-PLAN.md
   - Document exact PR sequence
   - Provide conflict resolutions
   - Give humans executable instructions
   
4. **HUMAN REVIEW PHASE**: PRs to main (humans only)
   - Humans create PRs from effort branches
   - Each PR reviewed and approved
   - Main branch modified ONLY through PRs

### Critical Addition: The Integration Path

```
/efforts/phase1/wave1/effort1 ↘
/efforts/phase1/wave1/effort2 → integration-testing → MASTER-PR-PLAN.md
/efforts/phase1/wave2/effort1 ↗                            ↓
                                                     Human PRs to main
```

### Main Branch Is SACRED (R280 SUPREME LAW):

```bash
# ❌ ABSOLUTELY FORBIDDEN:
git checkout main
git merge anything  # NEVER!
git push origin main  # NEVER!

# ✅ REQUIRED APPROACH:
git checkout main
git pull origin main
git checkout -b integration-testing-$(date +%Y%m%d-%H%M%S)
# Do ALL integration work in integration-testing branch
# Generate PR plan for humans
```

### Integration Testing Implementation:

```bash
integrate_for_validation() {
    # Start from clean main
    git checkout main
    git pull origin main
    
    # Create integration testing branch
    INTEGRATE_WAVE_EFFORTS_BRANCH="integration-testing-$(date +%Y%m%d-%H%M%S)"
    git checkout -b "$INTEGRATE_WAVE_EFFORTS_BRANCH"
    
    # Merge all efforts in order
    for effort in $(cat effort-order.txt); do
        echo "Integrating $effort..."
        git merge --no-ff "origin/$effort" || {
            # Document conflicts for PR plan
            echo "$effort: merge conflict" >> conflicts.log
        }
    done
    
    # Validate it works
    make build || go build .
    make test || go test ./...
    
    # Generate PR plan
    generate_master_pr_plan > MASTER-PR-PLAN.md
    
    echo "✅ Integration validated in $INTEGRATE_WAVE_EFFORTS_BRANCH"
    echo "✅ Main branch untouched (per R280)"
    echo "✅ PR plan ready for humans"
}
```

### Key Principles:

1. **Development Isolation**: All code development in `/efforts/`
2. **Integration Testing**: Validate everything works in test branch
3. **Main Protection**: NEVER modify main directly (R280)
4. **Human Control**: Final integration through reviewed PRs
5. **Audit Trail**: Every change to main has human approval

## SECTION 4: AGENT-SPECIFIC IMPLICATIONS

### 🤖 SW ENGINEER:
- **MUST** work in `/efforts/phaseX/waveY/effort-name/` (target repo clone)
- **NEVER** create .go files in the SF instance
- **NEVER** run `go test` in the SF instance
- **ALWAYS** verify: `pwd | grep -q "/efforts/"` before coding

### 📋 CODE REVIEWER:
- **Reviews** happen in `/efforts/phaseX/waveY/effort-name/` (target repo)
- **Plans** can be created in effort directory
- **NEVER** review orchestrator-state-v3.json as if it's source code

### 🏗️ ARCHITECT:
- **Can** work from SF instance for high-level review
- **Must** check target repo clones for actual code review
- **Understands** the separation between planning and implementation

### 🎭 ORCHESTRATOR:
- **Lives** in the SF instance
- **Spawns** agents to work in target repo clones
- **NEVER** writes code anywhere
- **Creates** target repo clones under `/efforts/`

## SECTION 5: FORBIDDEN AND CORRECT ACTIONS

### 🚫 FORBIDDEN ACTIONS:

### In Software Factory Instance (where orchestrator-state-v3.json lives):
```bash
# ❌ NEVER DO THESE IN SF INSTANCE:
echo "package main" > main.go           # NO! Wrong repo!
go test ./...                           # NO! No Go code here!
make build                              # NO! No Makefile here!
git merge effort-branch                 # NO! Different repo!
```

### In Target Repository Clones (under /efforts/):
```bash
# ❌ NEVER DO THESE IN TARGET CLONES:
echo "current_state: PLANNING" > orchestrator-state-v3.json  # NO! Wrong repo!
vi .claude/agents/orchestrator.md                        # NO! Wrong repo!
vi phase-plans/phase1-plan.md                           # NO! Wrong repo!
```

## ✅ CORRECT ACTIONS:

### In Software Factory Instance:
```bash
# ✅ CORRECT in SF instance:
vi orchestrator-state-v3.json              # Yes, state management
vi phase-plans/phase1-wave1-plan.md    # Yes, planning
spawn_agent "sw-engineer"               # Yes, orchestration
```

### In Target Repository Clones:
```bash
# ✅ CORRECT in target clones:
vi pkg/controller/webhook.go           # Yes, source code
go test ./pkg/...                      # Yes, testing
make build                             # Yes, building
git commit -m "feat: add webhook"      # Yes, version control
```

## SECTION 6: DIRECTORY STRUCTURE CLARITY

```
software-factory-instance/              # 🎭 ORCHESTRATION HAPPENS HERE
├── .claude/                           # Agent configs
├── orchestrator-state-v3.json            # State tracking
├── target-repo-config.yaml            # Points to target repo
├── phase-plans/                       # Planning documents
├── utilities/                         # SF tools
└── efforts/                           # 🔽 TARGET REPO CLONES BELOW
    └── phase1/
        └── wave1/
            ├── core-api/              # 💻 Clone 1 of TARGET REPO
            │   ├── go.mod            # Real code here!
            │   ├── pkg/              # Real code here!
            │   └── Makefile          # Real build here!
            ├── webhooks/              # 💻 Clone 2 of TARGET REPO
            │   ├── go.mod            # Real code here!
            │   ├── pkg/              # Real code here!
            │   └── Makefile          # Real build here!
            └── controllers/           # 💻 Clone 3 of TARGET REPO
                ├── go.mod            # Real code here!
                ├── pkg/              # Real code here!
                └── Makefile          # Real build here!
```

## SECTION 7: THE GOLDEN RULES

1. **ONE** Software Factory instance (for orchestration)
2. **MANY** Target Repository clones (for development)
3. **NEVER** mix planning files with source code
4. **ALWAYS** develop in `/efforts/` subdirectories
5. **EFFORT ISOLATION**: Each effort gets its own target repo clone

## SECTION 8: VALIDATION QUESTIONS

Before ANY action, ask yourself:
1. Am I in the right repository for this task?
2. Should this file exist in a planning repo or code repo?
3. Is `orchestrator-state-v3.json` visible? Then NO CODE HERE!
4. Am I under `/efforts/`? Then CODE GOES HERE!

## SECTION 9: GRADING IMPACT

**AUTOMATIC FAILURE** if:
- Source code is written in the SF instance
- Planning files are created in target repo clones
- Integration happens in the SF instance
- Confusion between the two repositories

## REMEMBER:

**Software Factory ORCHESTRATES work on the Target Repository**
**Software Factory NEVER CONTAINS Target Repository code**
**Target Repository clones NEVER CONTAIN Software Factory configs**

This is the law. There are no exceptions.