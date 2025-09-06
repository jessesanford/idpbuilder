# Orchestrator - INIT State Rules (OPTIMIZED)

## State Entry Requirements
**Prerequisites**: Bootstrap rules already loaded (R203, R006, R319, R322, R288)

## INIT-Specific Rules (3 FILES)

### 1. 🚨🚨🚨 R281 - Complete State File Initialization (SUPREME LAW #7)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R281-initial-state-file-creation.md`
**INIT Requirement**: Create complete orchestrator-state.yaml with ALL phases/waves/efforts
**Validation**: Must match implementation plan exactly or -100% failure

### 2. 🚨🚨🚨 R191 - Target Repository Configuration (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R191-target-repo-config.md`  
**INIT Requirement**: Validate target-repo-config.yaml exists and is NOT the SF repo
**Validation**: Missing or invalid config = cannot proceed

### 3. 🚨🚨🚨 R192 - Repository Separation Protocol (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R192-repo-separation.md`
**INIT Requirement**: Verify Software Factory and target repo are completely separate
**Validation**: Any mixing = immediate termination

## INIT State Actions

### Phase 1: Configuration Validation
```bash
# Validate target repository config (R191)
if [ ! -f "./target-repo-config.yaml" ]; then
    echo "FATAL: No target-repo-config.yaml"
    exit 191
fi

# Verify separation (R192)
TARGET_URL=$(yq '.target_repository.url' target-repo-config.yaml)
SF_ORIGIN=$(git remote get-url origin)
if [ "$TARGET_URL" = "$SF_ORIGIN" ]; then
    echo "FATAL: Target is SF repo!"
    exit 192
fi
```

### Phase 2: State File Creation (R281)
```bash
# Parse implementation plan
# Extract ALL phases, waves, efforts
# Create complete state file structure
# Validate completeness before proceeding
```

### Phase 3: Transition Preparation
```yaml
Next State: PLANNING or WAITING_FOR_IMPLEMENTATION_PLAN
Trigger: State file created and plan exists
Stop Required: YES (R322)
```

## Verification Marker
```bash
touch .state_rules_read_orchestrator_INIT
echo "$(date +%s) - INIT rules acknowledged" > .state_rules_read_orchestrator_INIT
```

## State Exit Protocol (R322)
1. Complete state file creation
2. Update current_state to next state
3. Commit orchestrator-state.yaml
4. STOP and summarize work
5. Wait for /continue-orchestrating

---
*INIT State Rules - Optimized Version*
*Dependencies: Bootstrap rules must be loaded first*