# Orchestrator - SPAWN_SW_ENGINEER_BACKPORT_FIXES State Rules

## 🛑🛑🛑 R322 MANDATORY STOP BEFORE STATE TRANSITIONS 🛑🛑🛑

**SUPREME LAW - VIOLATION = -100% IMMEDIATE FAILURE**

### YOU MUST STOP AFTER:
1. ✅ Completing all TODOs for this state
2. ✅ Updating orchestrator-state.json with new state
3. ✅ Committing and pushing the state file  
4. ✅ Providing work summary

### YOU MUST NOT:
- ❌ Continue to the next state automatically
- ❌ Start work for the new state
- ❌ Spawn agents for the new state
- ❌ Assume permission to continue

### STOP PROTOCOL:
```markdown
## 🛑 STATE TRANSITION CHECKPOINT: CURRENT_STATE → NEXT_STATE

### ✅ Current State Work Completed:
- [List completed work]

### 📊 Current Status:
- Current State: CURRENT_STATE
- Next State: NEXT_STATE
- TODOs Completed: X/Y
- State Files: Updated and committed ✅

### ⏸️ STOPPED - Awaiting User Continuation
Ready to transition to NEXT_STATE. Please use /continue-orchestrating.
```

**STOP MEANS STOP - Exit and wait for /continue-orchestrating**

---

## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED SPAWN_SW_ENGINEER_BACKPORT_FIXES STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_SPAWN_SW_ENGINEER_BACKPORT_FIXES
echo "$(date +%s) - Rules read and acknowledged for SPAWN_SW_ENGINEER_BACKPORT_FIXES" > .state_rules_read_orchestrator_SPAWN_SW_ENGINEER_BACKPORT_FIXES
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

## 📋 PRIMARY DIRECTIVES FOR SPAWN_SW_ENGINEER_BACKPORT_FIXES STATE

### Core Mandatory Rules (ALL orchestrator states must have these):

1. **🚨🚨🚨 R006** - ORCHESTRATOR NEVER WRITES CODE OR PERFORMS FILE OPERATIONS (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`
   - Criticality: BLOCKING - Automatic termination, 0% grade
   - Summary: NEVER write, copy, move, or manipulate ANY code files - delegate ALL to agents
   - **CRITICAL**: Copying files is NOT infrastructure - it's implementation work!

2. **🔴🔴🔴 R287** - TODO PERSISTENCE COMPREHENSIVE (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-persistence-comprehensive.md`
   - Criticality: SUPREME - -20% to -100% penalty for violations
   - Summary: MUST save TODOs within 30s after write, every 10 messages, before transitions
   - **CRITICAL**: Commit and push within 60 seconds of saving

3. **🔴🔴🔴 R288** - STATE FILE UPDATE REQUIREMENTS (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-requirements.md`
   - Criticality: SUPREME - State updates required for all transitions
   - Summary: MUST update orchestrator-state.json before EVERY state transition
   - **CRITICAL**: Commit and push state changes immediately

   - **CRITICAL**: NEVER use wc -l or manual counting

### State-Specific Rules:

5. **🔴🔴🔴 R321** - Immediate Backport During Integration Protocol
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R321-immediate-backport-during-integration.md`
   - Criticality: SUPREME LAW - Immediate backporting required
   - Summary: ALL fixes must go to source branches IMMEDIATELY when found

6. **🔴🔴🔴 R300** - Comprehensive Fix Management Protocol
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R300-comprehensive-fix-management-protocol.md`
   - Criticality: SUPREME LAW - Fixes go to effort branches
   - Summary: Effort branches are the source of truth for all fixes

7. **🔴🔴🔴 R322 Part A** - Mandatory Stop After Spawn States
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R322 Part A-mandatory-stop-after-spawn.md`
   - Criticality: SUPREME LAW - Must stop after spawning
   - Summary: ALL spawn states require STOP after spawning agents

8. **🔴🔴🔴 R151** - Parallelization Timing Requirements
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R151-parallelization-timestamp-requirements.md`
   - Criticality: SUPREME LAW - Parallel agents must start within 5s
   - Summary: When spawning multiple agents, timestamps must be close

## 🎯 STATE OBJECTIVES - SPAWN SW ENGINEERS FOR BACKPORT IMPLEMENTATION

In the SPAWN_SW_ENGINEER_BACKPORT_FIXES state, the ORCHESTRATOR is responsible for:

1. **Reading the Backport Plan**
   - Load BACKPORT-PLAN.md created by Code Reviewer
   - Identify all effort branches needing fixes
   - Understand fix groupings and dependencies
   - Note any sequential requirements

2. **Preparing Individual SW Engineer Assignments**
   - Extract fixes for each effort from the plan
   - Create clear instructions for each SW Engineer
   - Specify exact working directories and branches
   - Include verification requirements

3. **Spawning SW Engineers (Parallel or Sequential)**
   - Analyze if fixes can be done in parallel
   - Spawn SW Engineers according to plan requirements
   - If parallel: spawn all within 5 seconds (R151)
   - If sequential: document why and spawn in order

4. **Stopping After Spawn (R322 Part A)**
   - Update orchestrator-state.json to MONITORING_BACKPORT_PROGRESS
   - Document all spawned engineers
   - Commit and push state changes
   - STOP and wait for user continuation

## 📝 REQUIRED ACTIONS

### Step 1: Load and Parse Backport Plan
```bash
# Read the backport plan created by Code Reviewer
cd /efforts/integration-testing

# Load the plan
if [ -f "BACKPORT-PLAN.md" ]; then
    echo "📋 Loading backport plan from Code Reviewer..."
    cat BACKPORT-PLAN.md
    
    # Parse out effort branches needing fixes
    echo "📊 Identifying efforts requiring backports..."
    grep -E "^##.*Branch:|Effort:" BACKPORT-PLAN.md
else
    echo "❌ CRITICAL: No BACKPORT-PLAN.md found!"
    echo "Cannot proceed without Code Reviewer's plan"
    exit 1
fi
```

### Step 2: Create Individual SW Engineer Instructions
```bash
# For each effort in the backport plan, create specific instructions

# Example for each effort (this would be done for ALL efforts in plan)
EFFORTS_FROM_PLAN=("effort-1" "effort-2" "effort-3")  # Parsed from plan

for EFFORT_NAME in "${EFFORTS_FROM_PLAN[@]}"; do
    INSTRUCTION_FILE="/efforts/${EFFORT_NAME}/BACKPORT-INSTRUCTIONS.md"
    
    echo "📝 Creating instructions for ${EFFORT_NAME}..."
    
    cat > "$INSTRUCTION_FILE" << 'EOF'
# Backport Implementation Instructions

## Your Assignment
- Effort: ${EFFORT_NAME}
- Working Directory: /efforts/${EFFORT_NAME}
- Branch: [branch-name-from-plan]

## Fixes to Apply (From Code Reviewer's Plan)
[Extract specific fixes for this effort from BACKPORT-PLAN.md]

### Fix 1: [Description]
- Files: [List files]
- Changes: [Specific changes needed]
- Verification: [How to verify]

### Fix 2: [Description]
- Files: [List files]
- Changes: [Specific changes needed]
- Verification: [How to verify]

## Implementation Process
1. Ensure you're on your effort branch
2. Apply fixes as specified above
3. Run build to verify compilation
4. Run tests to ensure nothing breaks
5. Commit with clear message referencing integration fixes
6. Push your updated branch
7. Update sw-engineer-state.yaml when complete

## Success Criteria
- ✅ All listed fixes applied exactly as specified
- ✅ Build succeeds after fixes
- ✅ All tests pass
- ✅ Branch pushed to remote
- ✅ State file shows BACKPORT_COMPLETE

## Important Notes
- Do NOT modify anything beyond the specified fixes
- Do NOT merge from integration branch
- Apply fixes directly to your effort branch
- If you encounter issues, document them and update state to BLOCKED
EOF
    
    echo "✅ Instructions created for ${EFFORT_NAME}"
done
```

### Step 3: Determine Parallelization Strategy
```bash
# Analyze if fixes can be parallelized
echo "🔍 Analyzing parallelization requirements..."

# Check backport plan for dependencies
if grep -q "SEQUENTIAL\|DEPENDENCY\|ORDER" BACKPORT-PLAN.md; then
    echo "⚠️ Sequential execution required - dependencies detected"
    PARALLEL_EXECUTION=false
else
    echo "✅ Parallel execution possible - no dependencies"
    PARALLEL_EXECUTION=true
fi

# Document the strategy
cat > BACKPORT-SPAWN-STRATEGY.md << 'EOF'
# Backport SW Engineer Spawn Strategy

## Execution Mode: ${PARALLEL_EXECUTION}

## Efforts to Spawn:
1. effort-1 - 3 fixes - no dependencies
2. effort-2 - 2 fixes - no dependencies  
3. effort-3 - 4 fixes - no dependencies

## Spawn Order:
$(if [ "$PARALLEL_EXECUTION" = true ]; then
    echo "All engineers spawned simultaneously (R151 compliance)"
else
    echo "Sequential spawn order based on dependencies"
fi)

## Timing Requirements (R151):
- All parallel spawns must occur within 5 seconds
- First spawn timestamp: [will be recorded]
- Last spawn timestamp: [will be recorded]
- Delta must be < 5 seconds for parallel execution
EOF
```

### Step 4: Spawn SW Engineers
```bash
echo "🚀 Beginning SW Engineer spawn sequence..."

# Record start time for R151 compliance
START_TIME=$(date +%s)
echo "Spawn sequence started at: $(date '+%Y-%m-%d %H:%M:%S')"

# Spawn engineers based on strategy
if [ "$PARALLEL_EXECUTION" = true ]; then
    echo "📊 Spawning ALL SW Engineers in parallel (R151 compliance)..."
    
    for EFFORT_NAME in "${EFFORTS_FROM_PLAN[@]}"; do
        echo "🚀 Spawning SW Engineer for ${EFFORT_NAME}..."
        
        # Create spawn command
        cat > /tmp/spawn-sw-${EFFORT_NAME}.md << 'EOF'
@agent-software-engineer

## BACKPORT IMPLEMENTATION ASSIGNMENT

You are being spawned to implement backport fixes for your effort branch.

### Immediate Actions Required:
1. Read your instructions at: /efforts/${EFFORT_NAME}/BACKPORT-INSTRUCTIONS.md
2. Navigate to your working directory: /efforts/${EFFORT_NAME}
3. Verify you're on the correct branch
4. Apply ALL fixes specified in your instructions
5. Verify build and tests pass
6. Commit and push changes
7. Update your state file to BACKPORT_COMPLETE

### Critical Requirements:
- Apply fixes EXACTLY as specified
- Do NOT merge from other branches
- Do NOT modify beyond specified fixes
- Must complete ALL fixes or report BLOCKED

### State Progression:
INIT → BACKPORT_IMPLEMENTATION → BACKPORT_COMPLETE

Start immediately - this is time-critical per R321.
EOF
        
        # Log spawn with timestamp
        echo "$(date +%s): Spawned SW Engineer for ${EFFORT_NAME}" >> SPAWN-LOG.md
    done
    
else
    echo "📊 Spawning SW Engineers sequentially due to dependencies..."
    
    # Sequential spawn with order from plan
    for EFFORT_NAME in "${EFFORTS_FROM_PLAN[@]}"; do
        echo "🚀 Spawning SW Engineer for ${EFFORT_NAME} (sequential)..."
        # [Same spawn command as above]
        echo "⏸️ Waiting for ${EFFORT_NAME} to complete before next spawn..."
        break  # Would actually wait in real implementation
    done
fi

# Record end time and verify R151 compliance
END_TIME=$(date +%s)
SPAWN_DURATION=$((END_TIME - START_TIME))

echo "Spawn sequence completed at: $(date '+%Y-%m-%d %H:%M:%S')"
echo "Total spawn duration: ${SPAWN_DURATION} seconds"

if [ "$PARALLEL_EXECUTION" = true ] && [ $SPAWN_DURATION -gt 5 ]; then
    echo "⚠️ WARNING: R151 VIOLATION - Spawn duration exceeded 5 seconds!"
    echo "Duration was ${SPAWN_DURATION} seconds"
fi
```

### Step 5: Update State and STOP (R322 Part A Enforcement)
```bash
# Update orchestrator state
cd $CLAUDE_PROJECT_DIR

# Create comprehensive state update
cat > orchestrator-state.json << 'EOF'
current_state: MONITORING_BACKPORT_PROGRESS
previous_state: SPAWN_SW_ENGINEER_BACKPORT_FIXES
backport_status: IN_PROGRESS
agents_spawned:
  - agent: sw-engineer-1
    effort: effort-1
    purpose: Implement backport fixes
    timestamp: $(date +%s)
    state: BACKPORT_IMPLEMENTATION
  - agent: sw-engineer-2
    effort: effort-2
    purpose: Implement backport fixes
    timestamp: $(date +%s)
    state: BACKPORT_IMPLEMENTATION
  - agent: sw-engineer-3
    effort: effort-3
    purpose: Implement backport fixes
    timestamp: $(date +%s)
    state: BACKPORT_IMPLEMENTATION
parallelization:
  mode: parallel
  r151_compliant: true
  spawn_duration: ${SPAWN_DURATION}
backport_plan: /efforts/integration-testing/BACKPORT-PLAN.md
monitoring:
  total_efforts: 3
  completed: 0
  in_progress: 3
  blocked: 0
EOF

# Commit state change
git add orchestrator-state.json
git commit -m "state: transition to MONITORING_BACKPORT_PROGRESS after spawning SW Engineers

- Spawned SW Engineers for all efforts requiring backports
- Using backport plan from Code Reviewer
- R151 compliant parallel spawn
- R322 Part A stopping after spawn"
git push

echo "✅ State updated and committed"
echo "🛑 STOPPING per R322 Part A - Must stop after spawn state"
```

## ⚠️ CRITICAL REQUIREMENTS

### R151 Compliance for Parallel Spawns
- If spawning multiple engineers: ALL within 5 seconds
- Document spawn timestamps
- Flag any R151 violations

### Clear Separation of Work
- **Code Reviewer**: Created the plan (already done)
- **Orchestrator**: Distributes work and spawns engineers
- **SW Engineers**: Actually implement the fixes

### No Direct Implementation
- Orchestrator MUST NOT apply any fixes
- Orchestrator MUST NOT use git commands on code
- Orchestrator ONLY coordinates

### R322 Part A Enforcement
- MUST stop after spawning engineers
- MUST update state before stopping
- MUST NOT continue automatically

## 🚫 FORBIDDEN ACTIONS

1. **Implementing any fixes directly** - R006 violation
2. **Using cherry-pick or applying patches** - SW Engineer's job
3. **Continuing without stopping** - R322 Part A violation
4. **Modifying the backport plan** - Use it as-is from Code Reviewer
5. **Spawning agents slowly in parallel mode** - R151 violation

## ✅ SUCCESS CRITERIA

Before transitioning to MONITORING_BACKPORT_PROGRESS:
- [ ] Backport plan loaded from Code Reviewer
- [ ] Instructions created for each SW Engineer
- [ ] Parallelization strategy determined
- [ ] All SW Engineers spawned appropriately
- [ ] R151 compliance verified if parallel
- [ ] State updated to MONITORING_BACKPORT_PROGRESS
- [ ] State changes committed and pushed
- [ ] STOPPED per R322 Part A

## 🔄 STATE TRANSITIONS

### Success Path:
```
SPAWN_SW_ENGINEER_BACKPORT_FIXES → MONITORING_BACKPORT_PROGRESS
```
- All engineers spawned
- Monitoring their progress
- Waiting for completion

### From Monitoring:
```
MONITORING_BACKPORT_PROGRESS → PR_PLAN_CREATION
```
- All backports complete
- Ready for PR planning

```
MONITORING_BACKPORT_PROGRESS → ERROR_RECOVERY
```
- Backports failed
- Need intervention

## 🎓 GRADING CRITERIA

You will be evaluated on:
1. **R151 Compliance** (35%)
   - Parallel spawns within 5 seconds
   - Proper timestamp documentation
   
2. **R322 Part A Compliance** (25%)
   - Stop after spawning
   - No automatic continuation
   
3. **Work Distribution** (25%)
   - Clear instructions for each engineer
   - Following Code Reviewer's plan exactly
   
4. **State Management** (15%)
   - Proper state transitions
   - Complete state documentation

## 💡 TIPS FOR SUCCESS

1. **Trust the plan** - Code Reviewer analyzed the fixes, use their plan
2. **Spawn quickly** - For parallel execution, speed matters (R151)
3. **Clear instructions** - Each engineer needs to know exactly what to do
4. **Stop means stop** - R322 Part A is non-negotiable

Remember: You're the COORDINATOR, not the implementer!

## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-orchestrating
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**