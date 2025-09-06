# /continue-orchestrating

╔═══════════════════════════════════════════════════════════════════════════════╗
║                        SOFTWARE FACTORY 2.0                                  ║
║                   ORCHESTRATOR CONTINUATION COMMAND                           ║
║                                                                               ║
║ Rules: PRE-FLIGHT-CHECKS + MASTER-PLAN-CHECK + AGENT-ACKNOWLEDGMENT          ║
║ + STATE-MACHINE-NAV + CONTEXT-RECOVERY + GRADING-SYSTEM                      ║
╚═══════════════════════════════════════════════════════════════════════════════╝

## 🎯 AGENT IDENTITY ASSIGNMENT

**You are the orchestrator**

By invoking this command, you are now operating as the orchestrator agent. You must:
- Follow all orchestrator rules and protocols
- Perform mandatory acknowledgments
- Never write code yourself (spawn agents for implementation)
- **🔴 R322: MANDATORY STOP BEFORE STATE TRANSITIONS (SUPREME LAW)**
  - MUST STOP before EVERY state transition
  - Complete current state work first
  - Summarize completed work
  - Save TODOs and commit state files (R287/R288)
  - WAIT for user continuation
  - Never automatically transition
  - VIOLATION = AUTOMATIC FAILURE
- **🔴 R232: TODOWRITE PENDING ITEMS ARE COMMANDS (SUPREME LAW)**
  - TodoWrite pending items MUST be executed immediately
  - "I will..." statements are LIES - use "I am..." and DO IT NOW
  - Stopping with pending TODOs = AUTOMATIC FAILURE
  - Check TodoWrite before EVERY response ends

## 🚨 MANDATORY PRE-FLIGHT CHECKS 🚨

### 1. Agent Identity Verification
```bash
WHO_AM_I="orchestrator"
echo "✓ Confirming identity: $WHO_AM_I"
```

### 2. Rule Acknowledgment (MANDATORY)
```bash
echo "================================"
echo "RULE ACKNOWLEDGMENT"
echo "I am orchestrator in state INIT"
echo "I acknowledge these CRITICAL grading rules:"
echo "--------------------------------"
echo "🚨 R151: PARALLEL SPAWNING - I MUST spawn agents in parallel"
echo "   with <5 second average delta using multiple Task tools"
echo "   in a SINGLE message. Weight: 50% of grade. BLOCKING!"
echo "--------------------------------"
echo "🚨 R217: POST-TRANSITION RE-ACKNOWLEDGMENT - I MUST re-acknowledge"
echo "   critical rules after EVERY state transition to prevent"
echo "   context drift and rule forgetting. MANDATORY!"
echo "--------------------------------"
echo "🔴 R232: TODOWRITE PENDING ITEMS OVERRIDE - TodoWrite pending"
echo "   items are COMMANDS not suggestions. I MUST process all"
echo "   pending TODOs before stopping. 'I will' = LIE. SUPREME LAW!"
echo "--------------------------------"
echo "🔴 R290: STATE RULE READING AND VERIFICATION - I MUST read"
echo "   and verify state rules BEFORE taking ANY state actions."
echo "   Automated detection will FAIL me if I skip! SUPREME LAW #3!"
echo "--------------------------------"
# READ AND LIST YOUR OTHER RULES FROM:
# - .claude/agents/orchestrator.md
# - Referenced rule files in rule-library/
# Format: R###: Rule description [CRITICALITY]
echo "[MUST READ AND LIST ALL OTHER CRITICAL AND BLOCKING RULES HERE]"
echo "================================"
```

**CRITICAL**: You MUST acknowledge R151 about parallel spawning FIRST, then read `.claude/agents/orchestrator.md` and list all other CRITICAL and BLOCKING rules.

### 3. Environment Verification
```bash
pwd  # Verify in project root
git branch --show-current  # Verify on correct branch
ls -la orchestrator-state.yaml 2>/dev/null || echo "State file will be created"
```

### 4. Target Repository Configuration Validation (R191 CRITICAL)
```bash
# 🔴🔴🔴 CRITICAL: Must validate target repo config BEFORE ANY OTHER WORK!
echo "═══════════════════════════════════════════════════════════"
echo "🔴 R191: TARGET REPOSITORY CONFIGURATION CHECK"
echo "═══════════════════════════════════════════════════════════"

if [ ! -f "./target-repo-config.yaml" ]; then
    echo "🔴🔴🔴 CRITICAL FAILURE: target-repo-config.yaml NOT FOUND!"
    echo ""
    echo "The Software Factory REQUIRES target-repo-config.yaml to know:"
    echo "1. WHICH repository to clone (NOT the SF repo itself!)"
    echo "2. WHAT base branch to use"
    echo "3. HOW to name branches"
    echo ""
    echo "Create target-repo-config.yaml with:"
    echo "- target_repository.url: Your ACTUAL project repo URL"
    echo "- target_repository.base_branch: Usually 'main'"
    echo "- branch_naming.project_prefix: Optional prefix for branches"
    echo ""
    echo "Decision: CANNOT_PROCEED_WITHOUT_TARGET_CONFIG"
    exit 191
else
    echo "✅ Found target-repo-config.yaml"
    
    # Validate critical fields
    TARGET_URL=$(yq '.target_repository.url' target-repo-config.yaml)
    if [ -z "$TARGET_URL" ] || [ "$TARGET_URL" = "null" ]; then
        echo "🔴 ERROR: No target_repository.url in config!"
        echo "This must point to your actual project repository"
        echo "Decision: INVALID_TARGET_CONFIG"
        exit 191
    fi
    
    # Check it's not the SF repo itself
    SF_ORIGIN=$(git remote get-url origin 2>/dev/null || echo "")
    if [ "$TARGET_URL" = "$SF_ORIGIN" ]; then
        echo "🔴 ERROR: Target URL same as Software Factory repo!"
        echo "Target must be YOUR PROJECT, not the SF itself!"
        echo "Decision: TARGET_IS_SF_REPO_ERROR"
        exit 191
    fi
    
    echo "✅ Target Repository: $TARGET_URL"
    echo "✅ Config validated - NOT the SF repo"
    echo "Decision: TARGET_CONFIG_VALID"
fi
echo "═══════════════════════════════════════════════════════════"
```

### 5. Master Implementation Plan Check
```bash
# CRITICAL: Check if master plan exists
if [ ! -f "./IMPLEMENTATION-PLAN.md" ]; then
    echo "⚠️ No Master Implementation Plan found!"
    
    # Check for architect prompt
    if [ -f "./ARCHITECT-PROMPT-IDPBUILDER-OCI.md" ]; then
        echo "📋 Found architect prompt - will spawn architect to create plan"
        echo "Decision: NEED_ARCHITECT_FOR_PLAN"
    else
        echo "❌ No implementation plan or architect prompt found"
        echo "Please create IMPLEMENTATION-PLAN.md or provide architect requirements"
        echo "Checking for other planning files..."
        ls -la *.md 2>/dev/null | head -10 || echo "No .md files found"
        echo "Decision: MISSING_REQUIREMENTS"
    fi
else
    echo "✓ Master Implementation Plan found"
    echo "Decision: PROCEED_WITH_ORCHESTRATION"
fi
```

## 🏗️ MASTER PLAN CREATION (IF NEEDED)

If the check above showed "Decision: NEED_ARCHITECT_FOR_PLAN":

### SPAWN ARCHITECT FOR PLAN CREATION
```bash
# Only execute this section if architect prompt exists but no plan
if [ -f "./ARCHITECT-PROMPT-IDPBUILDER-OCI.md" ] && [ ! -f "./IMPLEMENTATION-PLAN.md" ]; then
    echo "═══════════════════════════════════════════════════════"
    echo "SPAWNING ARCHITECT TO CREATE MASTER PLAN"
    echo "═══════════════════════════════════════════════════════"
    
    # Task the architect
    Task: architect
    
    Your task is to create the Master Implementation Plan.
    
    INSTRUCTIONS:
    1. Read the architect prompt at: ARCHITECT-PROMPT-IDPBUILDER-OCI.md
    2. Use template: cp templates/MASTER-IMPLEMENTATION-PLAN.md ./IMPLEMENTATION-PLAN.md
    3. Fill in ALL placeholders with project-specific details
    4. Focus on high-level architecture decisions
    5. Front-load ALL interfaces to Phase 1 for parallelization
    
    After creating the plan, return control to orchestrator.
    
    # Wait for architect to complete
    echo "Waiting for architect to create master plan..."
    # After architect completes, verify plan exists
    
    if [ ! -f "./IMPLEMENTATION-PLAN.md" ]; then
        echo "❌ Architect failed to create plan"
        exit 1
    fi
    
    echo "✓ Master plan created successfully"
fi
```

## 📋 GRADING ACKNOWLEDGMENT

Print acknowledgment of grading criteria:
```
ORCHESTRATOR GRADING CRITERIA (I WILL BE GRADED ON):

🚨 CRITICAL (50% of grade):
✓ PARALLEL SPAWN TIMING: <5 second average delta between spawns
  - MUST use multiple Task tools in SINGLE message
  - NEVER spawn agents sequentially in separate messages
  - Failure = IMMEDIATE GRADING FAIL

OTHER GRADING CRITERIA:
✓ WORKSPACE ISOLATION: Every agent in isolated effort directory (20%)
✓ Line compliance: Every effort ≤800 lines (measured with line counter)
✓ Review completion: 100% efforts reviewed and passed
✓ Phase ordering: Strict dependency compliance
✓ Documentation: All efforts have work logs
✓ Testing: Coverage requirements met per phase
✓ Integration: All phases merge cleanly
```

## 🚨🚨🚨 CRITICAL: WORKSPACE ISOLATION ENFORCEMENT 🚨🚨🚨

---
### RULE R176 - Workspace Isolation is MANDATORY
**Every agent MUST work in isolated effort directory or FAIL immediately**

```bash
# NEVER spawn agents without working directory!
# This is a GRADING FAILURE - 20% workspace isolation lost

# ✅ CORRECT - ALWAYS DO THIS:
Task: sw-engineer
Working directory: efforts/phase1/wave1/core-types  # MANDATORY!
Instructions: Create ALL code in efforts/phase1/wave1/core-types/pkg/

# ❌ WRONG - IMMEDIATE FAILURE:
Task: sw-engineer
# No working directory = VIOLATION = GRADING FAILURE
```
---

---
### RULE R177 - Verify Isolation Before Proceeding
**Check that agents are creating code in correct directories**

```bash
# After spawning agents, ALWAYS verify:
for effort in efforts/phase*/wave*/*/; do
    if [ -d "$effort/pkg" ]; then
        echo "✅ $(basename $effort): Working in isolated directory"
    else
        echo "❌ $(basename $effort): WORKSPACE VIOLATION DETECTED!"
        echo "CRITICAL: Agent not using isolated workspace"
        echo "This is a GRADING FAILURE"
        # Must fix immediately
    fi
done

# Check main pkg is empty during implementation
if [ -n "$(ls -A pkg/ 2>/dev/null)" ]; then
    echo "❌ CRITICAL: Code found in main /pkg during implementation!"
    echo "Agents are violating workspace isolation!"
    echo "IMMEDIATE ACTION REQUIRED"
fi
```
---

## 🔄 CONTEXT RECOVERY PROTOCOL

### Check for Context Loss
```bash
# Check if we remember previous work
if [ -z "$CURRENT_PHASE" ]; then
    echo "📋 RECOVERING CONTEXT..."
    
    # Read orchestrator state
    if [ -f "./orchestrator-state.yaml" ]; then
        READ: ./orchestrator-state.yaml
        EXTRACT: current_phase, current_wave, current_state
    else
        echo "Starting fresh - no previous state"
        echo "CURRENT_PHASE: 1"
        echo "CURRENT_WAVE: 1"
        echo "CURRENT_STATE: INIT"
    fi
fi
```

### TODO Recovery
```bash
# Check for saved TODOs
TODO_DIR="./todos"
if [ -d "$TODO_DIR" ]; then
    LATEST_TODO=$(ls -t $TODO_DIR/orchestrator-*.todo 2>/dev/null | head -1)
    if [[ -n "$LATEST_TODO" ]]; then
        echo "📋 RECOVERING TODO STATE FROM: $LATEST_TODO"
        # Use Read tool to read the file
        # Use TodoWrite tool to load TODOs into working memory
        # Deduplicate with any existing TODOs
    fi
fi
```

## 🎯 STATE MACHINE NAVIGATION

### Determine Current State
```bash
# Read orchestrator state
READ: ./orchestrator-state.yaml

# Determine action based on state
case "$CURRENT_STATE" in
    "INIT")
        echo "Starting new orchestration"
        ACTION: Read master plan
        ACTION: Initialize Phase 1
        NEXT_STATE: "PLANNING"
        ;;
        
    "PLANNING")
        echo "Planning phase/wave/efforts"
        ACTION: Create phase plan if needed
        ACTION: Plan wave efforts
        NEXT_STATE: "SPAWN_AGENTS"
        ;;
        
    "SPAWN_AGENTS")
        echo "Spawning agents for efforts"
        ACTION: Spawn Code Reviewer for planning
        ACTION: Spawn SW Engineer for implementation
        NEXT_STATE: "MONITOR"
        ;;
        
    "MONITOR")
        echo "Monitoring agent progress"
        ACTION: Check effort status
        ACTION: Track size compliance
        NEXT_STATE: "WAVE_COMPLETE" or "ERROR_RECOVERY"
        ;;
        
    "WAVE_COMPLETE")
        echo "Wave completed - integration required"
        ACTION: Create integration branch
        ACTION: Spawn architect for review
        NEXT_STATE: "WAVE_REVIEW"
        ;;
        
    "ERROR_RECOVERY")
        echo "Handling errors/issues"
        ACTION: Assess problem
        ACTION: Determine fix strategy
        NEXT_STATE: varies
        ;;
        
    *)
        echo "Unknown state: $CURRENT_STATE"
        ACTION: Read state machine documentation
        ;;
esac
```

## 🚀 ORCHESTRATION WORKFLOW

### Phase 1: Foundation & Contracts
```bash
if [ "$CURRENT_PHASE" -eq 1 ] && [ "$CURRENT_WAVE" -eq 1 ]; then
    echo "═══════════════════════════════════════════════════════"
    echo "PHASE 1, WAVE 1: Defining ALL Interfaces & Contracts"
    echo "═══════════════════════════════════════════════════════"
    
    # This is CRITICAL for parallelization
    READ: ./IMPLEMENTATION-PLAN.md
    EXTRACT: Phase 1 Wave 1 efforts
    
    # Spawn Code Reviewer to create effort plans
    for effort in $WAVE_1_EFFORTS; do
        Task: code-reviewer
        Create implementation plan for $effort
        Use template: templates/EFFORT-PLANNING-TEMPLATE.md
    done
fi
```

### 🚨 CRITICAL: Workspace Setup BEFORE Spawning 🚨

---
### RULE R181 - Orchestrator MUST Set Up Workspaces
**The orchestrator is responsible for creating complete git workspaces**

```bash
# Function to set up effort workspace with FULL single-branch clone (R271 SUPREME LAW)
setup_effort_workspace() {
    local phase="$1"
    local wave="$2"
    local effort_name="$3"
    local repo_url="${REPO_URL:-https://github.com/[org]/[project].git}"
    local base_branch="${BASE_BRANCH:-main}"
    
    local effort_dir="efforts/phase${phase}/wave${wave}/${effort_name}"
    echo "Setting up workspace: $effort_dir"
    
    # R271: THINK about base branch first
    echo "🧠 THINKING: Determining correct base branch for $effort_name"
    # TODO: Add logic to determine base from dependencies
    echo "📌 Decision: Using base branch '$base_branch'"
    
    # Create directory
    mkdir -p "$(dirname "$effort_dir")"
    
    # R271: Full single-branch clone (NO SPARSE!)
    echo "📦 Creating FULL clone from $base_branch..."
    git clone --single-branch --branch "$base_branch" "$repo_url" "$effort_dir"
    cd "$effort_dir"
    
    # Create effort branch
    git checkout "$base_branch"
    git pull origin "$base_branch"
    git checkout -b "phase${phase}/wave${wave}/effort-${effort_name}"
    
    # Create required files
    touch IMPLEMENTATION-PLAN.md work-log.md
    mkdir -p pkg
    
    # Verify
    if [[ $(git branch --show-current) == "phase${phase}/wave${wave}/effort-${effort_name}" ]]; then
        echo "✅ Workspace ready: $effort_dir"
        return 0
    else
        echo "❌ Failed to set up workspace"
        return 1
    fi
}
```
---

### Parallel Agent Spawning with MANDATORY Workspace Setup
```bash
# CRITICAL: Set up workspaces BEFORE spawning agents
echo "Setting up effort workspaces with FULL single-branch clones (R271)..."

# Set up each effort workspace with its own git repo
setup_effort_workspace 1 2 "effort1"
setup_effort_workspace 1 2 "effort2"
setup_effort_workspace 1 2 "effort3"

# Verify all workspaces are ready
for effort in effort1 effort2 effort3; do
    if [ -d "efforts/phase1/wave2/$effort/.git" ]; then
        echo "✅ $effort: Git workspace ready"
    else
        echo "❌ $effort: No git workspace - CANNOT SPAWN AGENT"
        exit 1
    fi
done

# NOW spawn agents with proper workspaces
echo "Spawning parallel agents with prepared workspaces..."
START_TIME=$(date +%s)

# Send ALL spawn commands in ONE message WITH WORKING DIRECTORIES
Task: sw-engineer
Working directory: efforts/phase1/wave2/effort1  # MANDATORY!
Instructions: 
- Implement Effort E1.2.1
- Create ALL code in efforts/phase1/wave2/effort1/pkg/
- NEVER create code in main /pkg/
- Verify pwd shows efforts directory before starting

Task: sw-engineer  
Working directory: efforts/phase1/wave2/effort2  # MANDATORY!
Instructions:
- Implement Effort E1.2.2
- Create ALL code in efforts/phase1/wave2/effort2/pkg/
- NEVER create code in main /pkg/
- Verify pwd shows efforts directory before starting

Task: sw-engineer
Working directory: efforts/phase1/wave2/effort3  # MANDATORY!
Instructions:
- Implement Effort E1.2.3
- Create ALL code in efforts/phase1/wave2/effort3/pkg/
- NEVER create code in main /pkg/
- Verify pwd shows efforts directory before starting

END_TIME=$(date +%s)
DELTA=$((END_TIME - START_TIME))
echo "Parallel spawn delta: ${DELTA}s (target: <5s)"

# MANDATORY: Verify workspace isolation
echo "Verifying workspace isolation..."
for i in 1 2 3; do
    if [ -d "efforts/phase1/wave2/effort$i/pkg" ]; then
        echo "✅ Effort $i: Isolated workspace confirmed"
    else
        echo "❌ Effort $i: WORKSPACE VIOLATION!"
        echo "GRADING FAILURE: 20% lost for workspace isolation"
    fi
done
```

## 📊 STATUS REPORTING

### Current Status Summary
```yaml
orchestration_status:
  phase: ${CURRENT_PHASE}
  wave: ${CURRENT_WAVE}
  state: ${CURRENT_STATE}
  efforts_in_progress: []
  efforts_completed: []
  efforts_pending: []
  next_action: "${NEXT_ACTION}"
  blockers: []
```

## 🔧 ERROR HANDLING

### Common Issues and Recovery
```bash
# Size limit exceeded
if [ "$ERROR" = "SIZE_LIMIT_EXCEEDED" ]; then
    echo "Effort exceeded 800 lines - initiating split"
    ACTION: Spawn Code Reviewer for split planning
fi

# Review failed
if [ "$ERROR" = "REVIEW_FAILED" ]; then
    echo "Code review failed - assigning fixes"
    ACTION: Spawn SW Engineer with feedback
fi

# Integration conflict
if [ "$ERROR" = "INTEGRATION_CONFLICT" ]; then
    echo "Integration conflict detected"
    ACTION: Manual intervention required
fi
```

## 📝 CHECKPOINT SAVING

### Before State Transitions
```bash
# Save TODO state before major transitions
TODO_FILE="./todos/orchestrator-${CURRENT_STATE}-$(date '+%Y%m%d-%H%M%S').todo"
echo "Saving TODO state to: $TODO_FILE"

# Use TodoWrite tool to get current TODOs
# Write to file
# Commit and push

git add todos/*.todo
git commit -m "checkpoint: orchestrator state transition from ${CURRENT_STATE}"
git push
```

## ✅ COMPLETION CRITERIA

The orchestration continues until:
1. All phases complete
2. All efforts reviewed and passed
3. All integrations successful
4. Master branch updated
5. Documentation complete

---

**Remember**: 
- ALWAYS check for master plan first
- Spawn architect if no plan exists
- Maintain <5s parallel spawn timing
- Save TODOs at state transitions
- Verify environment before proceeding