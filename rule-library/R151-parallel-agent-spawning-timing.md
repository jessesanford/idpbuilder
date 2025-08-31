# Rule R151: Parallel Agent Spawning Timing and Acknowledgment

## Rule Statement
When spawning multiple agents that CAN be parallelized according to the plan:
1. **CHECK** parallelization headers in the appropriate plan document
2. **ACKNOWLEDGE** the plan creator's parallelization decision  
3. **DETERMINE** the correct working directory for EACH agent (R208)
4. **CD** to the appropriate directory BEFORE spawning (R208)
5. **SPAWN** all parallel agents in ONE message with <5s average delta
6. **CITE** the specific plan location that authorized parallelization

## 🚨 SPAWN_AGENTS IS A VERB - IMMEDIATE ACTION REQUIRED 🚨

### IMMEDIATE ACTIONS UPON ENTERING SPAWN_AGENTS STATE

**THE MOMENT YOU ENTER THIS STATE, YOU MUST:**
1. Start spawning SW Engineers for ready efforts NOW
2. Check implementation plans exist before spawning
3. Check TodoWrite for pending items and process them immediately
4. Follow parallelization plan from state file without delay

**FORBIDDEN - AUTOMATIC FAILURE:**
- ❌ "STATE TRANSITION COMPLETE: Now in SPAWN_AGENTS" [stops]
- ❌ "Successfully entered SPAWN_AGENTS state" [waits]
- ❌ "Ready to start spawning agents" [pauses]
- ❌ "I'm in SPAWN_AGENTS state" [does nothing]

**REQUIRED - IMMEDIATE ACTION:**
- ✅ "Entering SPAWN_AGENTS, spawning SW Engineers for ready efforts NOW..."
- ✅ "SPAWN_AGENTS: Immediately checking and spawning per implementation plans..."
- ✅ "SPAWN_AGENTS: Following parallelization plan from state file, spawning now..."

The SPAWN_AGENTS state name is a VERB - it means actively spawning agents RIGHT NOW, not preparing to spawn or thinking about spawning.

## Criticality Level
**CRITICAL** - Worth 50% of orchestrator grade (Timing violation = -50%)

## Enforcement Mechanism
- **Technical**: All parallel spawns MUST be in single message
- **Behavioral**: MUST acknowledge plan creator's decision
- **Grading**: -50% if average spawn delta >5s, -25% if no acknowledgment
- **Audit**: Check spawn timestamps and acknowledgment messages

## Core Principle

```
Plan Says Parallel → Spawn Together in ONE Message
Plan Says Sequential → Spawn One at a Time
ALWAYS Acknowledge the Plan Creator's Decision
```

## Context-Aware Header Locations

### When Spawning Code Reviewers for Effort Planning
```bash
# CHECK: Wave Implementation Plan parallelization headers
PLAN_FILE="phase-plans/PHASE-${PHASE}-WAVE-${WAVE}-IMPLEMENTATION-PLAN.md"

# LOOK FOR: Effort parallelization metadata
grep "Can Parallelize:" "$PLAN_FILE"
grep "Parallel With:" "$PLAN_FILE"

# ACKNOWLEDGE:
echo "📋 R151: Acknowledging parallelization decision from Wave Plan creator"
echo "   Plan location: $PLAN_FILE"
echo "   Decision: Efforts [X,Y,Z] marked as 'Can Parallelize: Yes'"
echo "   Spawning these agents IN PARALLEL as directed by plan"
```

### When Spawning SW Engineers for Implementation
```bash
# CHECK: Effort Implementation Plans in effort directories
for effort in $EFFORTS; do
    PLAN_FILE="efforts/phase${PHASE}/wave${WAVE}/${effort}/IMPLEMENTATION-PLAN.md"
    
    # LOOK FOR: Parallelization context section
    grep "Can Parallelize:" "$PLAN_FILE"
    grep "Parallel With:" "$PLAN_FILE"
done

# ACKNOWLEDGE:
echo "📋 R151: Acknowledging parallelization decision from Code Reviewer plans"
echo "   These efforts can run in parallel per their IMPLEMENTATION-PLAN.md files"
echo "   Spawning SW Engineers TOGETHER as authorized by Code Reviewer"
```

### When Spawning for Split Efforts
```bash
# CHECK: Split Summary for parallelization guidance
SPLIT_SUMMARY="efforts/phase${PHASE}/wave${WAVE}/${effort}/SPLIT-SUMMARY.md"

# LOOK FOR: Split execution strategy
grep "Execution Order:" "$SPLIT_SUMMARY"
grep "Can Parallelize Splits:" "$SPLIT_SUMMARY"

# ACKNOWLEDGE:
echo "📋 R151: Acknowledging split strategy from Code Reviewer"
echo "   Plan location: $SPLIT_SUMMARY"
echo "   Decision: Splits must run SEQUENTIALLY (standard protocol)"
echo "   Spawning splits one at a time as required"
```

## Detailed Requirements

### 🚨 MANDATORY: Acknowledge Plan Creator's Decision

**BEFORE SPAWNING, OUTPUT:**
```
📋 R151: Reading parallelization decision from [PLAN_TYPE] at:
   [full/path/to/plan.md]
   
✅ Plan Creator's Decision:
   - Efforts that CAN parallelize: [list]
   - Efforts that CANNOT parallelize: [list]
   - Rationale: [Dependencies/Independence noted in plan]
   
🚀 R151: Spawning agents according to plan creator's parallelization strategy
```

### 🎯 CRITICAL: Think About AND VERIFY Directories (R208)

**🚨 MANDATORY PWD VERIFICATION PROTOCOL USING BASH TOOL 🚨**

The orchestrator MUST:
1. **USE THE BASH TOOL** to CD to each effort directory
2. **USE THE BASH TOOL** to output `pwd` 
3. **SEE THE ACTUAL PATH** in Bash tool output
4. **ONLY THEN** spawn the agent

⚠️ JUST ECHOING COMMANDS IS NOT ENOUGH - MUST USE BASH TOOL! ⚠️

**BEFORE ANY PARALLEL SPAWN, YOU MUST:**
```bash
# 1. DETERMINE the working directory for EACH agent
echo "🗂️ R208: THINKING about where each agent should work..."
for effort in $PARALLEL_EFFORTS; do
    EFFORT_DIR="efforts/phase${PHASE}/wave${WAVE}/${effort}"
    echo "   Effort ${effort} needs directory: ${EFFORT_DIR}"
    
    # VERIFY directory exists or create it
    if [ ! -d "$EFFORT_DIR" ]; then
        echo "   Creating directory: $EFFORT_DIR"
        mkdir -p "$EFFORT_DIR"
    fi
done

# 2. PLAN the spawn sequence with MANDATORY pwd verification
echo "📋 R151 + R208: Planning parallel spawn with VERIFIED directories:"
echo "   Will CD to each directory"
echo "   Will OUTPUT pwd to PROVE correct location"
echo "   Will spawn ONLY after verification"
echo "   All spawns in ONE message for <5s delta"
```

**REQUIRED SEQUENCE FOR EACH SPAWN (USING BASH TOOL!):**

```
Step 1: Announce intention
   echo "🗂️ R208: Need to CD to: efforts/phase1/wave1/oci-stack-types"

Step 2: USE BASH TOOL to CD and verify
   Bash: cd efforts/phase1/wave1/oci-stack-types && pwd
   
Step 3: Bash tool outputs (YOU MUST SEE THIS):
   /home/vscode/workspaces/idpbuilder-oci-mgmt/efforts/phase1/wave1/oci-stack-types
   
Step 4: Confirm correct directory
   echo "✅ Confirmed in correct directory via Bash tool output"
   
Step 5: NOW spawn the agent
   Task: sw-engineer
   Working directory: efforts/phase1/wave1/oci-stack-types
```

⚠️ IF YOU DON'T SEE THE BASH TOOL OUTPUT WITH PWD, YOU'RE DOING IT WRONG! ⚠️

### Timing Requirements for Parallel Spawns WITH Directory Protocol

```bash
# CRITICAL: All parallel spawns in ONE message WITH proper directories!
parallel_spawn_agents_with_directories() {
    local PHASE="$1"
    local WAVE="$2"
    local AGENTS="$3"
    local PLAN_LOCATION="$4"
    
    echo "═══════════════════════════════════════════════════════════════"
    echo "🚀 R151 + R208: PARALLEL SPAWN WITH DIRECTORY PROTOCOL"
    echo "═══════════════════════════════════════════════════════════════"
    
    # Record start time
    START_TIME=$(date +%s%N)
    
    echo "📋 R151: Acknowledging parallelization from:"
    echo "   $PLAN_LOCATION"
    echo ""
    echo "🗂️ R208: Setting up directories for parallel agents:"
    
    # R208: Show directory plan BEFORE spawning
    for effort in $AGENTS; do
        echo "   ${effort} → efforts/phase${PHASE}/wave${WAVE}/${effort}/"
    done
    
    echo ""
    echo "🚨 R151 + R208 CRITICAL: Spawning ALL agents in ONE message..."
    echo "   Each in its CORRECT directory per R208!"
    
    # ALL SPAWNS IN SINGLE MESSAGE - WITH DIRECTORY VERIFICATION!
    # MANDATORY: Show pwd before EACH spawn to prove correct directory
    for effort in $AGENTS; do
        echo "🗂️ CD'ing to: efforts/phase${PHASE}/wave${WAVE}/${effort}"
        echo "cd efforts/phase${PHASE}/wave${WAVE}/${effort}"
        echo "📍 VERIFYING pwd: $(pwd)  # MUST OUTPUT: efforts/phase${PHASE}/wave${WAVE}/${effort}"
        echo "✅ Confirmed in correct directory, spawning agent..."
        echo "Task: code-reviewer"
        echo "  Purpose: Create implementation plan"
        echo "  Effort: ${effort}"
        echo "  Working directory verified: $(pwd)"
    done
    
    # Calculate timing
    END_TIME=$(date +%s%N)
    DELTA=$((($END_TIME - $START_TIME) / 1000000))  # Convert to ms
    NUM_AGENTS=$(echo "$AGENTS" | wc -w)
    AVG_DELTA=$(($DELTA / $NUM_AGENTS))
    
    echo "✅ R151 Timing: ${NUM_AGENTS} agents spawned"
    echo "   Total time: ${DELTA}ms"
    echo "   Average delta: ${AVG_DELTA}ms"
    echo "✅ R208 Compliance: Each agent in correct directory"
    
    if [ $AVG_DELTA -lt 5000 ]; then
        echo "   Grade: PASS (< 5000ms + correct directories)"
    else
        echo "   Grade: FAIL (timing or directory violation!)"
    fi
}
```

## Integration with R218 and R208

R151 works in conjunction with other critical spawning rules:

### R151 + R218 + R208 Combined Protocol
- **R218**: Determines WHAT to spawn in parallel (reads wave plan)
- **R208**: Determines WHERE to spawn each agent (correct directory)
- **R151**: Ensures HOW to spawn in parallel (timing + acknowledgment)

**CRITICAL COMBINATION:**
```bash
# WRONG - Spawning without considering directory
spawn_parallel_agents "effort1 effort2 effort3"  # Where are they?

# RIGHT - Full compliance with R151 + R208 + R218
echo "📖 R218: Reading wave plan to determine parallelization..."
READ: phase-plans/PHASE-1-WAVE-1-IMPLEMENTATION-PLAN.md

echo "📋 R151: Acknowledging plan creator's parallel decision"
echo "🗂️ R208: Determining correct directories for each agent:"
echo "   Effort 1: efforts/phase1/wave1/effort1/"
echo "   Effort 2: efforts/phase1/wave1/effort2/"
echo "   Effort 3: efforts/phase1/wave1/effort3/"

echo "🚀 R151 + R208: Spawning with VERIFIED directories in ONE message:"

# MANDATORY: CD and VERIFY directory for effort1
echo "📍 CD'ing to effort1 directory..."
cd efforts/phase1/wave1/effort1
echo "📍 pwd verification: $(pwd)"  # MUST SHOW: efforts/phase1/wave1/effort1
echo "✅ Confirmed in effort1 directory, spawning..."
Task: code-reviewer "Create implementation plan for effort1"

# MANDATORY: CD and VERIFY directory for effort2
echo "📍 CD'ing to effort2 directory..."
cd efforts/phase1/wave1/effort2
echo "📍 pwd verification: $(pwd)"  # MUST SHOW: efforts/phase1/wave1/effort2
echo "✅ Confirmed in effort2 directory, spawning..."
Task: code-reviewer "Create implementation plan for effort2"

# MANDATORY: CD and VERIFY directory for effort3
echo "📍 CD'ing to effort3 directory..."
cd efforts/phase1/wave1/effort3
echo "📍 pwd verification: $(pwd)"  # MUST SHOW: efforts/phase1/wave1/effort3
echo "✅ Confirmed in effort3 directory, spawning..."
Task: code-reviewer "Create implementation plan for effort3"

echo "✅ All agents spawned in VERIFIED directories with <5s delta"
```

### Combined Compliance Example
```bash
# R218: Read wave plan to determine parallelization
echo "📖 R218: Using Read tool to read wave plan..."
READ: phase-plans/PHASE-1-WAVE-1-IMPLEMENTATION-PLAN.md

echo "📖 R218: I have READ the parallelization headers from wave plan located at:"
echo "   phase-plans/PHASE-1-WAVE-1-IMPLEMENTATION-PLAN.md"

# R151: Acknowledge the plan creator's decision
echo "📋 R151: Acknowledging Wave Plan creator's parallelization strategy:"
echo "   Efforts 3,4,5 marked 'Can Parallelize: Yes' by Code Reviewer"
echo "   Efforts 1,2 marked 'Can Parallelize: No' (blocking)"

# R218 + R151: Spawn according to plan with proper timing
echo "✅ Spawning Code Reviewers as allowed by the parallelization headers in wave plan"
echo "🚀 R151: Executing parallel spawn for efforts 3,4,5 in ONE message..."

# Spawn all parallel agents together
Task: code-reviewer "Effort 3"
Task: code-reviewer "Effort 4" 
Task: code-reviewer "Effort 5"
```

## Grading Checklist

**R151 Compliance (50% of grade):**
- [ ] Did orchestrator check parallelization headers? (-10% if no)
- [ ] Did orchestrator acknowledge plan creator's decision? (-15% if no)
- [ ] Did orchestrator cite specific plan location? (-10% if no)
- [ ] Did orchestrator CD to correct directory for EACH agent? (R208) (-20% if no)
- [ ] **Did orchestrator OUTPUT pwd to VERIFY directory?** (-25% if no)
- [ ] Were parallel agents spawned in ONE message? (-25% if no)
- [ ] Was average spawn delta <5s? (-50% if no)

**Directory Violations (R208 Integration):**
- [ ] Spawning in wrong directory = MISSION CRITICAL FAILURE
- [ ] Not outputting pwd verification = -25% (CANNOT VERIFY COMPLIANCE)
- [ ] Not checking effort directories before spawn = -20%
- [ ] Spawning all agents in same directory = -30%
- [ ] pwd shows wrong directory but spawn continues = -50%

## Common Violations to Avoid

### ❌ Spawning Parallel Agents Sequentially
```bash
# WRONG - Violates R151 timing requirement
spawn_agent "Effort 3"
wait 2
spawn_agent "Effort 4"  # Should be in same message!
wait 2
spawn_agent "Effort 5"  # Should be in same message!
# Average delta: >5s = FAIL
```

### ❌ Not Acknowledging Plan Creator
```bash
# WRONG - No acknowledgment of parallelization decision
spawn_parallel_agents "3 4 5"  # Where did this come from?
```

### ❌ Not Verifying Directory with pwd
```bash
# WRONG - No pwd verification, can't prove correct directory
cd efforts/phase1/wave1/effort1
Task: sw-engineer  # But are we really in effort1?

# WRONG - Just claiming without proving
echo "Spawning in correct directory"
Task: sw-engineer  # No pwd output = no proof!
```

### ✅ Correct Implementation
```bash
# RIGHT - Full compliance with pwd verification
echo "📋 R151: Code Reviewer's wave plan at phase-plans/PHASE-1-WAVE-1-IMPLEMENTATION-PLAN.md"
echo "        authorizes parallel execution of efforts 3,4,5"
echo "🚀 R151 + R208: Spawning 3 agents in ONE message with VERIFIED directories"

# All in one message WITH pwd verification for each
echo "🗂️ CD'ing to: efforts/phase1/wave1/effort3"
cd efforts/phase1/wave1/effort3
echo "📍 pwd verification: $(pwd)"  # Shows: /workspace/target-repo/efforts/phase1/wave1/effort3
Task: sw-engineer

echo "🗂️ CD'ing to: efforts/phase1/wave1/effort4"
cd efforts/phase1/wave1/effort4
echo "📍 pwd verification: $(pwd)"  # Shows: /workspace/target-repo/efforts/phase1/wave1/effort4
Task: sw-engineer

echo "🗂️ CD'ing to: efforts/phase1/wave1/effort5"
cd efforts/phase1/wave1/effort5
echo "📍 pwd verification: $(pwd)"  # Shows: /workspace/target-repo/efforts/phase1/wave1/effort5
Task: sw-engineer

echo "✅ R151: Average spawn delta: 500ms (PASS)"
echo "✅ R208: All directories verified with pwd (PASS)"
```

## Special Cases

### Sequential Requirements (No Parallelization)
When plan says "Can Parallelize: No":
```bash
echo "📋 R151: Acknowledging plan creator's SEQUENTIAL requirement"
echo "   Plan: $PLAN_FILE"
echo "   Decision: This effort blocks others - must complete first"
echo "   Spawning single agent as directed"

spawn_single_agent  # No parallel spawn needed
```

### Mixed Parallelization
When some can parallelize, some cannot:
```bash
# First: Sequential blocking efforts
echo "📋 R151: Effort 1 marked 'Can Parallelize: No' - spawning alone"
spawn_agent "Effort 1"
wait_for_completion

# Then: Parallel independent efforts
echo "📋 R151: Efforts 2,3,4 marked 'Can Parallelize: Yes' - spawning together"
spawn_parallel_agents "2 3 4"  # All in ONE message
```

## Relationship to Other Rules

- **R218**: Reads wave plan parallelization headers (determines WHAT to parallelize)
- **R208**: Orchestrator spawn directory protocol (determines WHERE to spawn)
- **R151**: Ensures proper spawn timing and acknowledgment (determines HOW to spawn)
- **R211**: Code Reviewers preserve parallelization metadata in plans
- **R176**: Workspace isolation requirements
- **R052**: General agent spawning protocol
- **R053**: Parallelization decision guidelines

**THE TRINITY OF PARALLEL SPAWNING:**
```
R218 (WHAT) + R208 (WHERE) + R151 (HOW) = Successful Parallel Spawn
```

## Summary

**R151 ensures efficient parallel spawning by:**
1. Checking the appropriate plan for parallelization headers
2. Acknowledging the plan creator's parallelization decisions
3. **THINKING about the correct directory for EACH agent (R208)**
4. **CD'ing to the appropriate directory BEFORE spawning (R208)**
5. Spawning all parallel agents in ONE message (<5s delta)
6. Citing specific plan locations that authorize parallelization
7. Respecting sequential requirements when parallelization is not allowed

**CRITICAL INTEGRATION:**
This rule works in tandem with:
- **R218** to determine WHAT to parallelize (read wave plan)
- **R208** to determine WHERE to spawn (correct directories)
- Together ensuring the orchestrator reads plans, respects parallelization decisions, AND spawns agents in their proper isolated workspaces.