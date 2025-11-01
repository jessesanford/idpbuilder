# ORCHESTRATOR STATE: VALIDATE_INFRASTRUCTURE

## 🔴🔴🔴 ABSOLUTE REQUIREMENT: STATE MANAGER CONSULTATION 🔴🔴🔴

**SUPREME LAW - R517 - UNIVERSAL STATE MANAGER CONSULTATION LAW**

**BEFORE EXITING THIS STATE, YOU MUST:**

1. **SPAWN STATE MANAGER FOR SHUTDOWN_CONSULTATION** (MANDATORY - NO EXCEPTIONS)
2. **NEVER UPDATE STATE FILES DIRECTLY** (orchestrator-state-v3.json, bug-tracking.json, etc.)
3. **NEVER COMMIT STATE FILES YOURSELF** (State Manager does this atomically)

**FAILURE TO CONSULT STATE MANAGER = IMMEDIATE SYSTEM HALT (-100% GRADE)**

### Enforcement Mechanism

If you attempt to exit this state without spawning State Manager:
- ❌ Pre-commit hooks will REJECT your commit
- ❌ Validation tools will FAIL the build
- ❌ Grading system will assign -100% penalty
- ❌ System will transition to ERROR_RECOVERY

### Required Pattern (COPY THIS EXACTLY)

```bash
# At end of state work, BEFORE any state file updates:

echo "🔴 MANDATORY: Spawning State Manager for SHUTDOWN_CONSULTATION"

# Spawn State Manager (REQUIRED - NOT OPTIONAL)
# Task: state-manager
# State: SHUTDOWN_CONSULTATION
# Current State: [YOUR_CURRENT_STATE]
# Proposed Next State: [YOUR_PROPOSED_NEXT_STATE]
# Work Summary: [SUMMARY_OF_WORK_COMPLETED]

# State Manager will:
# 1. Validate proposed transition against state machine
# 2. Update all 4 state files atomically
# 3. Commit with [R288] tag
# 4. Return REQUIRED next state (may differ from proposal)

# Wait for State Manager response
# Follow State Manager's directive (REQUIRED next state)
# DO NOT proceed until State Manager confirms
```

**YOU MUST NEVER:**
- ❌ Update orchestrator-state-v3.json yourself
- ❌ Update bug-tracking.json yourself
- ❌ Update integration-containers.json yourself
- ❌ Use `jq` to modify state files
- ❌ Use `sed/awk` to modify state files
- ❌ Set `validated_by: "orchestrator"` (must be "state-manager")
- ❌ Skip State Manager consultation "just this once"
- ❌ Think "I'll validate it manually"

**ONLY State Manager may update state files. This is NON-NEGOTIABLE.**

See: `rule-library/R517-universal-state-manager-consultation-law.md`
## 🚨🚨🚨 STATE PURPOSE [BLOCKING]
Validate ALL infrastructure configuration against authoritative sources BEFORE any implementation work begins. This is a MANDATORY gate that prevents catastrophic failures.

## 🔴🔴🔴 SUPREME LAW ENFORCEMENT
This state enforces:
- **R507**: Mandatory Infrastructure Validation [BLOCKING]
- **R508**: Target Repository Enforcement [SUPREME LAW]

ANY violation requires immediate transition to ERROR_RECOVERY.

## 🚨🚨🚨 ENTRY CONDITIONS [BLOCKING]
MUST have:
1. ✅ Just completed CREATE_NEXT_INFRASTRUCTURE
2. ✅ Infrastructure directories created
3. ✅ Need to verify configuration before spawning agents

## 🚨🚨🚨 REQUIRED ACTIONS [MANDATORY]

### 1. RUN VALIDATION SCRIPT
```bash
echo "🔍 Starting infrastructure validation (R507/R508)..."
bash $CLAUDE_PROJECT_DIR/utilities/validate-infrastructure.sh

VALIDATION_RESULT=$?

if [ $VALIDATION_RESULT -eq 0 ]; then
    echo "✅ Infrastructure validation PASSED"
    # Can proceed to next state
elif [ $VALIDATION_RESULT -eq 911 ]; then
    echo "🔴🔴🔴 CATASTROPHIC: Wrong repository configured (R508 violation)"
    # MUST transition to ERROR_RECOVERY
else
    echo "❌ Infrastructure validation FAILED"
    # MUST transition to ERROR_RECOVERY
fi
```

### 2. CHECK EACH INFRASTRUCTURE COMPONENT
For EVERY effort/split/integration directory:

#### A. VALIDATE EFFORT INFRASTRUCTURE
```bash
# For each effort in pre_planned_infrastructure
for effort_key in $(jq -r '.pre_planned_infrastructure.efforts | keys[]' orchestrator-state-v3.json); do
    EFFORT_DATA=$(jq -r ".pre_planned_infrastructure.efforts.\"$effort_key\"" orchestrator-state-v3.json)
    EFFORT_PATH=$(echo "$EFFORT_DATA" | jq -r '.full_path')
    EFFORT_BRANCH=$(echo "$EFFORT_DATA" | jq -r '.branch_name')
    EFFORT_TARGET_REPO=$(echo "$EFFORT_DATA" | jq -r '.target_repo_url')

    echo "Validating effort: $effort_key"

    if [ -d "$EFFORT_PATH" ]; then
        cd "$EFFORT_PATH"
        ACTUAL_REMOTE=$(git remote get-url origin 2>/dev/null || git remote get-url target)
        ACTUAL_BRANCH=$(git branch --show-current)

        if [ "$ACTUAL_REMOTE" != "$EFFORT_TARGET_REPO" ]; then
            echo "🔴🔴🔴 CATASTROPHIC: Wrong repository for effort $effort_key!"
            echo "Expected: $EFFORT_TARGET_REPO"
            echo "Actual: $ACTUAL_REMOTE"
            exit 911  # IMMEDIATE ERROR_RECOVERY
        fi

        if [ "$ACTUAL_BRANCH" != "$EFFORT_BRANCH" ]; then
            echo "❌ Wrong branch name for effort $effort_key!"
            echo "Expected: $EFFORT_BRANCH"
            echo "Actual: $ACTUAL_BRANCH"
            VALIDATION_FAILED=true
        fi
    fi
done
```

#### B. VALIDATE INTEGRATE_WAVE_EFFORTS INFRASTRUCTURE (R504 CRITICAL)
```bash
# Check wave integrations
for wave_key in $(jq -r '.pre_planned_infrastructure.integrations.wave_integrations | keys[]' orchestrator-state-v3.json); do
    WAVE_DATA=$(jq -r ".pre_planned_infrastructure.integrations.wave_integrations.\"$wave_key\"" orchestrator-state-v3.json)
    WAVE_DIR=$(echo "$WAVE_DATA" | jq -r '.directory')
    WAVE_BRANCH=$(echo "$WAVE_DATA" | jq -r '.branch_name')
    WAVE_TARGET_REPO=$(echo "$WAVE_DATA" | jq -r '.target_repo_url')
    WAVE_CREATED=$(echo "$WAVE_DATA" | jq -r '.created')

    echo "Validating wave integration: $wave_key"

    if [ "$WAVE_CREATED" == "true" ] && [ -d "$WAVE_DIR" ]; then
        cd "$WAVE_DIR"
        ACTUAL_REMOTE=$(git remote get-url origin 2>/dev/null || git remote get-url target)

        if [ "$ACTUAL_REMOTE" != "$WAVE_TARGET_REPO" ]; then
            echo "🔴🔴🔴 CATASTROPHIC: Wrong repository for integration $wave_key!"
            echo "Expected: $WAVE_TARGET_REPO"
            echo "Actual: $ACTUAL_REMOTE"
            exit 911  # IMMEDIATE ERROR_RECOVERY
        fi
    fi
done

# Check phase integrations
for phase_key in $(jq -r '.pre_planned_infrastructure.integrations.phase_integrations | keys[]' orchestrator-state-v3.json 2>/dev/null); do
    PHASE_DATA=$(jq -r ".pre_planned_infrastructure.integrations.phase_integrations.\"$phase_key\"" orchestrator-state-v3.json)
    PHASE_TARGET_REPO=$(echo "$PHASE_DATA" | jq -r '.target_repo_url')

    if [ "$PHASE_TARGET_REPO" != "null" ]; then
        CONFIGURED_REPO=$(yq '.repository_url' target-repo-config.yaml)
        if [ "$PHASE_TARGET_REPO" != "$CONFIGURED_REPO" ]; then
            echo "🔴🔴🔴 CATASTROPHIC: Phase integration pre-planned for wrong repository!"
            echo "Expected: $CONFIGURED_REPO"
            echo "Pre-planned: $PHASE_TARGET_REPO"
            exit 911
        fi
    fi
done
```

#### C. VALIDATE REMOTE REPOSITORY CONSISTENCY (R508 - SUPREME LAW)
```bash
CONFIGURED_REPO=$(yq '.repository_url' target-repo-config.yaml)

# Verify ALL pre-planned infrastructure points to correct target repo
WRONG_REPOS=$(jq -r '
    [.pre_planned_infrastructure.efforts[].target_repo_url,
     .pre_planned_infrastructure.integrations.wave_integrations[].target_repo_url // empty,
     .pre_planned_infrastructure.integrations.phase_integrations[].target_repo_url // empty]
    | map(select(. != null and . != "'"$CONFIGURED_REPO"'"))
    | unique[]' orchestrator-state-v3.json 2>/dev/null)

if [ -n "$WRONG_REPOS" ]; then
    echo "🔴🔴🔴 CATASTROPHIC: Pre-planned infrastructure has wrong repository URLs!"
    echo "Expected: $CONFIGURED_REPO"
    echo "Found wrong URLs: $WRONG_REPOS"
    exit 911
fi
```

### 3. UPDATE STATE FILE
```bash
# Record validation results
jq '.infrastructure_validation = {
  "timestamp": "'$(date -u +"%Y-%m-%dT%H:%M:%SZ")'",
  "status": "validated",
  "repositories_verified": true,
  "branches_verified": true,
  "paths_verified": true
}' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
```

## 🚨🚨🚨 EXIT CONDITIONS [MANDATORY]
ONE of:
1. ✅ All validation checks PASSED → Continue to next state
2. ❌ Validation FAILED → Must transition to ERROR_RECOVERY
3. 🔴 CATASTROPHIC failure (exit code 911) → IMMEDIATE ERROR_RECOVERY

## 🚨🚨🚨 STATE TRANSITIONS [BLOCKING]

### PROJECT_DONE PATH
If validation passes:
- Single effort detected → `SPAWN_CODE_REVIEWERS_EFFORT_PLANNING` (R356 optimization - skip parallelization analysis)
- Multiple efforts → `ANALYZE_CODE_REVIEWER_PARALLELIZATION` (normal flow)

### FAILURE PATH
If validation fails:
- ANY failure → `ERROR_RECOVERY`
- Exit code 911 → `ERROR_RECOVERY` with CATASTROPHIC flag

## 🚨🚨 COMMON ERRORS TO AVOID

1. **NEVER skip validation**
```bash
# ❌ WRONG - Missing validation
CREATE_NEXT_INFRASTRUCTURE
transition_to SPAWN_SW_ENGINEERS  # NO VALIDATION!
```

2. **NEVER ignore validation failures**
```bash
# ❌ WRONG - Ignoring failure
bash validate-infrastructure.sh || true  # NEVER DO THIS
transition_to SPAWN_SW_ENGINEERS  # Proceeding despite failure!
```

3. **NEVER hardcode repository URLs**
```bash
# ❌ WRONG - Hardcoded URL
git remote add origin "https://github.com/some/repo.git"

# ✅ CORRECT - From config
REPO=$(yq '.repository_url' target-repo-config.yaml)
git remote add origin "$REPO"
```

## 🚨🚨🚨 GRADING PENALTIES
- Skipping validation: **-50%**
- Proceeding after validation failure: **-75%**
- Wrong repository (R508 violation): **-100% IMMEDIATE FAILURE**
- Not transitioning to ERROR_RECOVERY on failure: **-50%**

## 📋 STATE OUTPUT REQUIREMENTS
At completion, MUST output:
```
=== VALIDATE_INFRASTRUCTURE Complete ===
Validation Status: [PASSED/FAILED]
Repositories Verified: [YES/NO]
Branches Verified: [YES/NO]
Paths Verified: [YES/NO]
Next State: [state_name]
===
```

## 🚨 MANDATORY STOP (R313)
After validation (success OR failure), you MUST:
1. Update orchestrator-state-v3.json
2. Output the state completion summary
3. **STOP IMMEDIATELY**
4. Do NOT continue to next state automatically

The continuation command will handle the transition based on validation results.

## Automation Flag

```bash
# After infrastructure validation
if [ "$VALIDATION_PASSED" = true ]; then
    echo "✅ Infrastructure validation passed"
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # Continue to next state
else
    echo "❌ Infrastructure validation failed - transitioning to error recovery"
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # Continue per R405 - system will handle recovery
fi
```

---

**REMEMBER**: This state is a CRITICAL SAFETY GATE. NEVER bypass or skip validation. The entire Software Factory depends on correct infrastructure configuration!
## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

**MUST output as LAST line before state completion:**
- `CONTINUE-SOFTWARE-FACTORY=TRUE` - State completed, automation can continue
- `CONTINUE-SOFTWARE-FACTORY=FALSE` - Error/block, manual intervention required

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-CONTINUATION-FLAG-MASTER-GUIDE.md**
