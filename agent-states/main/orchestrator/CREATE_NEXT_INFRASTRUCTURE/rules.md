# Orchestrator - CREATE_NEXT_INFRASTRUCTURE State Rules

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

## 🔴🔴🔴 R360 - JUST-IN-TIME INFRASTRUCTURE CREATION 🔴🔴🔴

**NEW APPROACH - CREATE INFRASTRUCTURE ONLY WHEN READY TO USE**

### This State Handles BOTH:
1. **Regular Efforts**: Created just before implementation
2. **Split Efforts**: Created just before split implementation

### Infrastructure Creation Timing:
```
Create infrastructure ONLY when:
- All dependencies are complete ✓
- Effort/split is next to be implemented ✓
- About to spawn agents for implementation ✓
```

### Dependency-Based Base Branch Selection:
```bash
# For efforts with dependencies
if effort.depends_on exists:
    base_branch = completed_effort(depends_on).branch
else:
    base_branch = previous_wave_integration

# For splits (always chain sequentially)
if split_number > 1:
    base_branch = previous_split.branch
else:
    base_branch = effort_base_or_wave_integration
```

---

## PRIMARY DIRECTIVES

You MUST read and acknowledge these rules:
- R006: Orchestrator cannot write code
- R287: TODO persistence requirements
- R288: State file update protocol
- R322: Mandatory stop before state transitions
- R360: Just-in-time infrastructure creation (NEW)
- R308: Incremental branching strategy
- R312: Git config immutability
- R340: Planning file metadata tracking
- R302: Comprehensive split tracking

## State Context

CREATE_NEXT_INFRASTRUCTURE = You are ACTIVELY creating infrastructure for the NEXT effort or split that needs it!

This could be:
- A regular effort at wave start
- A dependent effort after its dependency completes
- A split after previous split completes
- An independent effort that can run in parallel

## 🚨🚨🚨 IMMEDIATE ACTIONS UPON ENTERING STATE 🚨🚨🚨

**THE INSTANT YOU ENTER THIS STATE:**

### 1. Determine What Needs Infrastructure

```bash
echo "🔧 DETERMINING NEXT INFRASTRUCTURE TO CREATE..."

# Check for efforts needing infrastructure
next_effort=$(jq -r '.effort_dependencies | to_entries[] |
  select(.value.infrastructure_created == false and
         (.value.depends_on == null or
          .effort_dependencies[.value.depends_on].status == "complete")) |
  .key' orchestrator-state.json | head -1)

# Check for splits needing infrastructure
next_split=$(jq -r '.split_tracking | to_entries[] |
  select(.value.current_split < .value.total_splits and
         .value.status == "in_progress") |
  .key' orchestrator-state.json | head -1)

if [ -n "$next_effort" ]; then
    echo "Creating infrastructure for effort: $next_effort"
    infrastructure_type="effort"
    infrastructure_target="$next_effort"
elif [ -n "$next_split" ]; then
    echo "Creating infrastructure for split: $next_split"
    infrastructure_type="split"
    infrastructure_target="$next_split"
else
    echo "No infrastructure needed - all complete!"
    # Transition to appropriate completion state
fi
```

### 2. Determine Base Branch

```bash
if [ "$infrastructure_type" = "effort" ]; then
    # Check if effort has dependency
    depends_on=$(jq -r ".effort_dependencies.$infrastructure_target.depends_on" orchestrator-state.json)

    if [ "$depends_on" != "null" ]; then
        # Use completed dependency as base
        BASE_BRANCH=$(jq -r ".effort_dependencies.$depends_on.branch" orchestrator-state.json)
        echo "Using dependency branch as base: $BASE_BRANCH"
    else
        # Use previous wave integration
        BASE_BRANCH=$(jq -r '.current_integration_base' orchestrator-state.json)
        echo "Using wave integration as base: $BASE_BRANCH"
    fi
elif [ "$infrastructure_type" = "split" ]; then
    # Splits always chain sequentially
    current_split=$(jq -r ".split_tracking.$infrastructure_target.current_split" orchestrator-state.json)

    if [ "$current_split" -eq 0 ]; then
        # First split uses effort base or wave integration
        BASE_BRANCH=$(jq -r '.current_integration_base' orchestrator-state.json)
    else
        # Subsequent splits chain from previous
        prev_split=$((current_split))
        BASE_BRANCH="${infrastructure_target}-split-$(printf "%03d" $prev_split)"
    fi
fi
```

### 3. Create the Infrastructure

```bash
echo "📁 Creating infrastructure for: $infrastructure_target"
echo "🔗 Base branch: $BASE_BRANCH"

# Determine directory path
if [ "$infrastructure_type" = "effort" ]; then
    INFRA_DIR="$CLAUDE_PROJECT_DIR/efforts/phase${CURRENT_PHASE}/wave${CURRENT_WAVE}/$infrastructure_target"
else
    INFRA_DIR="$CLAUDE_PROJECT_DIR/efforts/phase${CURRENT_PHASE}/wave${CURRENT_WAVE}/$infrastructure_target-split-$(printf "%03d" $((current_split + 1)))"
fi

# Clone repository
echo "📦 Cloning repository..."
git clone "$TARGET_REPO" "$INFRA_DIR"
cd "$INFRA_DIR"

# Create and push branch
BRANCH_NAME="$PROJECT_PREFIX/phase${CURRENT_PHASE}/wave${CURRENT_WAVE}/$infrastructure_target"
echo "🌿 Creating branch: $BRANCH_NAME from $BASE_BRANCH"
git checkout -b "$BRANCH_NAME" "$BASE_BRANCH"
git push -u origin "$BRANCH_NAME"

# Lock git config (R312)
echo "🔒 Locking git config per R312..."
chmod 444 .git/config

# Copy plan if it exists (R340)
if [ "$infrastructure_type" = "effort" ]; then
    plan_location=$(jq -r ".planning_files.effort_plans.$infrastructure_target" orchestrator-state.json)
else
    plan_location=$(jq -r ".planning_files.split_plans.$infrastructure_target" orchestrator-state.json)
fi

if [ -n "$plan_location" ] && [ "$plan_location" != "null" ]; then
    echo "📋 Plan will be used from: $plan_location"
    # SW Engineer will read the plan from tracked location
fi
```

### 4. Update Tracking

```bash
# Update orchestrator-state.json
if [ "$infrastructure_type" = "effort" ]; then
    # Update effort tracking
    jq --arg effort "$infrastructure_target" \
       --arg branch "$BRANCH_NAME" \
       --arg base "$BASE_BRANCH" \
       '.effort_dependencies[$effort].infrastructure_created = true |
        .effort_dependencies[$effort].base_branch = $base |
        .effort_dependencies[$effort].branch = $branch |
        .effort_dependencies[$effort].status = "ready"' \
       orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
else
    # Update split tracking
    jq --arg split "$infrastructure_target" \
       '.split_tracking[$split].current_split += 1' \
       orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
fi

# Commit state changes
git add orchestrator-state.json
git commit -m "state: Created infrastructure for $infrastructure_target [R360]"
git push
```

## State Exit Criteria

Infrastructure is successfully created when:
1. ✅ Repository cloned to correct location
2. ✅ Branch created from correct base (dependency or wave integration)
3. ✅ Branch pushed to remote
4. ✅ Git config locked (R312)
5. ✅ Tracking updated in orchestrator-state.json
6. ✅ State committed and pushed

## Next State Transitions

After successfully creating infrastructure:

- **If effort infrastructure created** → SPAWN_AGENTS
- **If split infrastructure created** → SPAWN_AGENTS
- **If no more infrastructure needed** → WAVE_COMPLETE or appropriate state

## Common Mistakes to Avoid

1. ❌ Creating all effort infrastructure at wave start
2. ❌ Creating infrastructure before dependencies complete
3. ❌ Using wrong base branch for dependent efforts
4. ❌ Forgetting to lock git config
5. ❌ Not updating tracking in state file

## Validation

Before transitioning to next state:
```bash
# Verify infrastructure was created
if [ ! -d "$INFRA_DIR/.git" ]; then
    echo "ERROR: Infrastructure directory not created!"
    exit 1
fi

# Verify correct branch
cd "$INFRA_DIR"
current_branch=$(git branch --show-current)
if [ "$current_branch" != "$BRANCH_NAME" ]; then
    echo "ERROR: Wrong branch!"
    exit 1
fi

# Verify git config is locked
if [ -w .git/config ]; then
    echo "ERROR: Git config not locked!"
    exit 1
fi
```

## Error Recovery

If infrastructure creation fails:
1. Check target repository access
2. Verify base branch exists
3. Ensure no duplicate branches
4. Clean up partial infrastructure
5. Update state to ERROR_RECOVERY