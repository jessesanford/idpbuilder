# 🚨🚨🚨 RULE R329 - Orchestrator NEVER Performs Git Merges

**Criticality:** BLOCKING - Automatic termination  
**Grading Impact:** IMMEDIATE FAILURE (0% grade)  
**Enforcement:** ZERO TOLERANCE - Single violation = termination
**Detection:** Automatic monitoring of git merge commands

## Rule Statement

The Orchestrator is a COORDINATOR ONLY and MUST NEVER perform git merge operations directly. ALL merge operations MUST be delegated to the Integration Agent, regardless of complexity. This extends R006 (Orchestrator never writes code) to include merge operations, which are a form of code manipulation.

## 🔴🔴🔴 ABSOLUTE PROHIBITIONS 🔴🔴🔴

The Orchestrator MUST NEVER:
- ❌ Execute `git merge` commands directly
- ❌ Perform `git rebase` operations
- ❌ Use `git cherry-pick` to apply commits
- ❌ Resolve merge conflicts manually
- ❌ Apply patches with `git apply`
- ❌ Use `git pull` with merge strategy
- ❌ Execute integration scripts that perform merges
- ❌ Combine branches in any way
- ❌ Edit files to resolve conflicts
- ❌ Make "simple" merges (NO EXCEPTIONS)

## ⚠️⚠️⚠️ CRITICAL CLARIFICATION ⚠️⚠️⚠️

**THERE ARE NO "SIMPLE" MERGES FOR ORCHESTRATOR**
- Single branch merge = ❌ FORBIDDEN (spawn Integration Agent)
- "No conflicts expected" = ❌ FORBIDDEN (spawn Integration Agent)  
- "Quick test merge" = ❌ FORBIDDEN (spawn Integration Agent)
- "Linear history merge" = ❌ FORBIDDEN (spawn Integration Agent)
- ALL merges = ❌ FORBIDDEN (ALWAYS spawn Integration Agent)

**CREATING INFRASTRUCTURE IS NOT MERGING**
- Creating integration workspace = ✅ Infrastructure
- Creating integration branch = ✅ Infrastructure
- Cloning repository = ✅ Infrastructure
- Merging branches = ❌ IMPLEMENTATION (FORBIDDEN)
- Combining code = ❌ IMPLEMENTATION (FORBIDDEN)

## What Orchestrator CAN Do

The Orchestrator MAY:
- ✅ Create integration workspace directories
- ✅ Clone repositories for integration
- ✅ Create empty integration branches
- ✅ Push empty branches to establish tracking
- ✅ Spawn Code Reviewer to create merge plans
- ✅ Spawn Integration Agent to execute merges
- ✅ Monitor integration progress via reports
- ✅ Coordinate fix distribution if issues found

## Required Delegation Pattern

### CORRECT: Always Spawn Integration Agent

```bash
# Step 1: Create integration infrastructure
create_integration_infrastructure() {
    echo "📋 Setting up integration workspace..."
    
    # Create workspace directory (✅ ALLOWED)
    INTEGRATION_DIR="/efforts/phase${PHASE}/wave${WAVE}/integration-workspace"
    mkdir -p "$INTEGRATION_DIR"
    
    # Clone repository (✅ ALLOWED)
    git clone "$TARGET_REPO_URL" "$INTEGRATION_DIR"
    cd "$INTEGRATION_DIR"
    
    # Create metadata directory (✅ ALLOWED)
    mkdir -p .software-factory
    
    # Create integration branch (✅ ALLOWED)
    git checkout -b "wave-${WAVE}-integration"
    git push -u origin "wave-${WAVE}-integration"
    
    # Document infrastructure (✅ ALLOWED)
    cat > INTEGRATION-SETUP.md << EOF
    Integration workspace created at: $INTEGRATION_DIR
    Integration branch: wave-${WAVE}-integration
    Ready for Integration Agent
    EOF
}

# Step 2: Spawn Code Reviewer for merge plan
spawn_code_reviewer_for_plan() {
    echo "🚀 Spawning Code Reviewer for merge plan..."
    
    Task: Create merge plan for Wave ${WAVE}
    Agent: code-reviewer
    State: WAVE_MERGE_PLANNING
    Output: WAVE-MERGE-PLAN.md
    
    # Transition to WAITING_FOR_MERGE_PLAN
}

# Step 3: Spawn Integration Agent to execute
spawn_integration_agent() {
    echo "🚀 Spawning Integration Agent to execute merges..."
    
    Task: Execute merges per WAVE-MERGE-PLAN.md
    Agent: integration-agent
    Workspace: $INTEGRATION_DIR
    Plan: WAVE-MERGE-PLAN.md
    
    # Transition to MONITORING_INTEGRATION
}
```

## Detection Mechanisms

```bash
# Monitor for forbidden merge operations
detect_orchestrator_merge_violation() {
    local command="$1"
    
    # Check for any merge-related git commands
    if [[ "$command" =~ git[[:space:]]+(merge|rebase|cherry-pick|pull.*--rebase) ]]; then
        echo "🚨🚨🚨 CRITICAL VIOLATION: R329 🚨🚨🚨"
        echo "ORCHESTRATOR ATTEMPTED MERGE OPERATION!"
        echo "Command: $command"
        echo "CONSEQUENCE: IMMEDIATE TERMINATION"
        echo "CORRECT: Spawn Integration Agent to handle merges"
        exit 329
    fi
    
    # Check for git apply (patch application)
    if [[ "$command" =~ git[[:space:]]apply ]]; then
        echo "🚨🚨🚨 CRITICAL VIOLATION: R329 🚨🚨🚨"
        echo "ORCHESTRATOR ATTEMPTED PATCH APPLICATION!"
        echo "Command: $command"
        echo "CONSEQUENCE: IMMEDIATE TERMINATION"
        echo "CORRECT: Spawn Integration Agent to apply patches"
        exit 329
    fi
}

# Pre-execution check for orchestrator
before_git_command() {
    local git_cmd="$1"
    
    echo "🤔 Self-check: Is this a merge operation?"
    
    if [[ "$git_cmd" =~ (merge|rebase|cherry-pick|apply) ]]; then
        echo "🚫 STOP! That's a merge operation!"
        echo "✅ Correct action: Spawn Integration Agent"
        return 1
    fi
    
    echo "✅ Safe git operation - not a merge"
    return 0
}
```

## Common Violations

### 🔴 VIOLATION: Direct Merge in INTEGRATION State
```bash
# ❌ CATASTROPHIC VIOLATION - Orchestrator merging directly
cd /efforts/phase1/wave1/integration-workspace
git merge origin/effort1-branch
git merge origin/effort2-branch

# CONSEQUENCE: Immediate 0% grade, orchestrator terminated
# CORRECT: Spawn Integration Agent to perform ALL merges
```

### 🔴 VIOLATION: Direct Merge in INTEGRATION_TESTING State
```bash
# ❌ CATASTROPHIC VIOLATION - From actual transcript
cd /efforts/integration-testing
git checkout project-integration
git merge phase1-integration
git merge phase2-integration

# CONSEQUENCE: Immediate 0% grade
# WHY: Orchestrator performed merges instead of spawning Integration Agent
# RESULT: Stale branches merged, fixes lost
```

### 🔴 VIOLATION: Making Excuses for Direct Merges
```bash
# ❌ INVALID EXCUSE: "It's just a simple fast-forward merge"
# ❌ INVALID EXCUSE: "There are no conflicts to resolve"  
# ❌ INVALID EXCUSE: "Integration Agent is overkill for one branch"
# ❌ INVALID EXCUSE: "We're in a hurry, I'll just do it quickly"

# REALITY: ALL merges require Integration Agent
# REALITY: Integration Agent ensures proper tracking and documentation
# REALITY: Skipping Integration Agent = violated architecture
# CORRECT: ALWAYS spawn Integration Agent, NO EXCEPTIONS
```

## Correct Patterns

### GOOD: Integration State Entry
```bash
# Orchestrator in INTEGRATION state
echo "📋 Entering INTEGRATION state"
echo "🏗️ Creating integration infrastructure..."

# Create workspace (allowed)
mkdir -p /efforts/phase1/wave1/integration-workspace

# Spawn Code Reviewer for plan
echo "🚀 Spawning Code Reviewer for merge plan..."
/spawn code-reviewer WAVE_MERGE_PLANNING

# Later, spawn Integration Agent
echo "🚀 Spawning Integration Agent for execution..."
/spawn integration-agent EXECUTE_MERGES
```

### GOOD: Project Integration
```bash
# Orchestrator in PROJECT_INTEGRATION state
echo "📋 Setting up project integration infrastructure"

# Create project workspace (allowed)
mkdir -p /efforts/project/integration-workspace

# Clone and setup (allowed)
git clone $REPO_URL /efforts/project/integration-workspace
cd /efforts/project/integration-workspace
mkdir -p .software-factory
git checkout -b project-integration

# Spawn specialists for actual work
echo "🚀 Spawning Code Reviewer for project merge plan..."
/spawn code-reviewer PROJECT_MERGE_PLANNING

echo "🚀 Spawning Integration Agent for project merges..."
/spawn integration-agent PROJECT_INTEGRATION_EXECUTION
```

## Grading Consequences

| Violation | First Offense | Second Offense |
|-----------|--------------|----------------|
| Executed git merge | IMMEDIATE FAILURE | N/A - Terminated |
| Performed rebase | IMMEDIATE FAILURE | N/A - Terminated |
| Applied patches | IMMEDIATE FAILURE | N/A - Terminated |
| Resolved conflicts | IMMEDIATE FAILURE | N/A - Terminated |

## Relationship to Other Rules

- **Extends R006**: Orchestrator never writes code (merging is code manipulation)
- **Supersedes R268**: Makes Integration Agent MANDATORY (not optional)
- **Complements R269**: Code Reviewer plans, Integration Agent executes
- **Supports R260**: Integration Agent has clear responsibilities
- **Enables R321**: Proper backporting through correct agent separation

## Why This Rule Exists

The bug described in the user request occurred because:
1. Orchestrator performed merges directly in INTEGRATION_TESTING state
2. Orchestrator merged an OLD project-integration branch
3. Fixes that were in effort branches were NOT in the stale integration branch
4. Result: All fixes were lost in the final integration

If Integration Agent had been spawned:
1. Integration Agent would fetch LATEST branches
2. Integration Agent would merge fresh code with all fixes
3. No stale branch problem would occur
4. All fixes would be properly integrated

## Mantra for Orchestrator

```
I am a COORDINATOR, not a merger
I CREATE infrastructure, never combine code
I SPAWN Integration Agents, never merge myself
I MONITOR progress, never execute merges
Every merge requires an Integration Agent
No exceptions, no excuses, no direct merges
```

---
**REMEMBER:** One merge command = Career over. Integration Agent handles ALL merges. NO EXCEPTIONS.