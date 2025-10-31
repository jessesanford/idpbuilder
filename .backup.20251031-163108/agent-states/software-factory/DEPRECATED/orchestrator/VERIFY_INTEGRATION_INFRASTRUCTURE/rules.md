# ORCHESTRATOR STATE: VERIFY_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE

## 🚨🚨🚨 STATE PURPOSE [BLOCKING]
Validate integration infrastructure configuration against authoritative sources AFTER creation and BEFORE spawning integration agents. This is a MANDATORY gate that prevents catastrophic failures from wrong repository targeting.

## 🔴🔴🔴 SUPREME LAW ENFORCEMENT
This state enforces:
- **R507**: Mandatory Infrastructure Validation [BLOCKING]
- **R508**: Target Repository Enforcement [SUPREME LAW]
- **R308**: Incremental Branching Strategy [SUPREME LAW]

ANY violation requires immediate transition to ERROR_RECOVERY with CONTINUE-SOFTWARE-FACTORY=FALSE.

## 🚨🚨🚨 CRITICAL: WRONG REPOSITORY = CATASTROPHIC FAILURE 🚨🚨🚨

**THIS STATE EXISTS TO PREVENT THE #1 CATASTROPHIC FAILURE:**
- Integration infrastructure on PLANNING repository instead of TARGET repository
- This is a SUPREME LAW violation (R508) = -100% IMMEDIATE FAILURE
- Exit code 911 = CATASTROPHIC repository mismatch

## 🚨🚨🚨 ENTRY CONDITIONS [BLOCKING]
MUST have:
1. ✅ Just completed SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE
2. ✅ Integration infrastructure directories created
3. ✅ Integration branch created and pushed
4. ✅ Need to verify EVERYTHING before spawning agents

## 🚨🚨🚨 REQUIRED ACTIONS [MANDATORY]

### 1. LOAD INTEGRATE_WAVE_EFFORTS METADATA
```bash
echo "🔍 Loading integration metadata from orchestrator-state-v3.json..."

INTEGRATE_WAVE_EFFORTS_TYPE=$(jq -r '.integration_infrastructure.type' orchestrator-state-v3.json)
INTEGRATE_WAVE_EFFORTS_BRANCH=$(jq -r '.integration_infrastructure.branch' orchestrator-state-v3.json)
INTEGRATE_WAVE_EFFORTS_DIR=$(jq -r '.integration_infrastructure.directory' orchestrator-state-v3.json)
BASE_BRANCH=$(jq -r '.integration_infrastructure.base_branch' orchestrator-state-v3.json)
PHASE=$(jq -r '.integration_infrastructure.phase' orchestrator-state-v3.json)
WAVE=$(jq -r '.integration_infrastructure.wave' orchestrator-state-v3.json)

echo "Integration Type: $INTEGRATE_WAVE_EFFORTS_TYPE"
echo "Integration Branch: $INTEGRATE_WAVE_EFFORTS_BRANCH"
echo "Integration Directory: $INTEGRATE_WAVE_EFFORTS_DIR"
echo "Base Branch: $BASE_BRANCH"
```

### 2. VALIDATE TARGET REPOSITORY (R508 - SUPREME LAW)
```bash
echo "🔴🔴🔴 VALIDATING TARGET REPOSITORY (R508 SUPREME LAW)"

# Get configured target repository
CONFIGURED_REPO=$(yq '.repository_url' target-repo-config.yaml)
if [ -z "$CONFIGURED_REPO" ]; then
    echo "🔴🔴🔴 CATASTROPHIC: No target repository configured!"
    exit 911
fi

# Check integration directory
cd "$INTEGRATE_WAVE_EFFORTS_DIR"

# Get actual remote
ACTUAL_REMOTE=$(git remote get-url origin 2>/dev/null || git remote get-url target 2>/dev/null)

# SUPREME LAW CHECK
if [ "$ACTUAL_REMOTE" != "$CONFIGURED_REPO" ]; then
    echo "🔴🔴🔴 CATASTROPHIC FAILURE: INTEGRATE_WAVE_EFFORTS ON WRONG REPOSITORY!"
    echo "Expected: $CONFIGURED_REPO"
    echo "Actual: $ACTUAL_REMOTE"
    echo "THIS IS A SUPREME LAW VIOLATION (R508)"
    echo ""
    echo "❌ INTEGRATE_WAVE_EFFORTS INFRASTRUCTURE IS ON THE WRONG REPOSITORY!"
    echo "❌ THIS IS LIKELY THE PLANNING REPOSITORY!"
    echo "❌ ALL WORK WILL BE LOST IF WE CONTINUE!"
    echo ""
    echo "IMMEDIATE ACTIONS REQUIRED:"
    echo "1. STOP ALL WORK"
    echo "2. TRANSITION TO ERROR_RECOVERY"
    echo "3. DELETE INCORRECT INFRASTRUCTURE"
    echo "4. RECREATE ON TARGET REPOSITORY"
    echo ""
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
    exit 911
fi

echo "✅ Integration repository verified: $CONFIGURED_REPO"
```

### 3. VALIDATE BRANCH NAME AND TRACKING
```bash
echo "🔍 Validating branch configuration..."

CURRENT_BRANCH=$(git branch --show-current)

if [ "$CURRENT_BRANCH" != "$INTEGRATE_WAVE_EFFORTS_BRANCH" ]; then
    echo "❌ Wrong branch!"
    echo "Expected: $INTEGRATE_WAVE_EFFORTS_BRANCH"
    echo "Actual: $CURRENT_BRANCH"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
    exit 1
fi

# Check remote tracking
TRACKING=$(git rev-parse --abbrev-ref --symbolic-full-name @{u} 2>/dev/null)
EXPECTED_TRACKING="origin/$INTEGRATE_WAVE_EFFORTS_BRANCH"

if [ "$TRACKING" != "$EXPECTED_TRACKING" ]; then
    echo "❌ Branch tracking incorrect!"
    echo "Expected: $EXPECTED_TRACKING"
    echo "Actual: $TRACKING"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
    exit 1
fi

echo "✅ Branch correctly configured: $CURRENT_BRANCH → $TRACKING"
```

### 4. VALIDATE INCREMENTAL BASE (R308 - SUPREME LAW)
```bash
echo "🔴🔴🔴 VALIDATING INCREMENTAL BASE (R308 SUPREME LAW)"

# Critical R308 validation for Phase 2+ Wave 1
if [[ "$INTEGRATE_WAVE_EFFORTS_TYPE" == "wave" && $PHASE -gt 1 && $WAVE -eq 1 ]]; then
    if [ "$BASE_BRANCH" == "main" ]; then
        echo "🔴🔴🔴 R308 VIOLATION DETECTED!"
        echo "Phase $PHASE Wave 1 integration CANNOT use main as base!"
        echo "Must use: phase$((PHASE-1))-integration"
        echo "Actual base: $BASE_BRANCH"
        echo ""
        echo "THIS IS A SUPREME LAW VIOLATION!"
        echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
        exit 308
    fi

    # Verify correct incremental base
    EXPECTED_BASE="phase$((PHASE-1))-integration"
    if [ "$BASE_BRANCH" != "$EXPECTED_BASE" ]; then
        echo "❌ Wrong incremental base!"
        echo "Expected: $EXPECTED_BASE"
        echo "Actual: $BASE_BRANCH"
        echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
        exit 308
    fi
fi

echo "✅ R308 Incremental base verified: $BASE_BRANCH"
```

### 5. VALIDATE DIRECTORY STRUCTURE
```bash
echo "🔍 Validating directory structure..."

if [ ! -d "$INTEGRATE_WAVE_EFFORTS_DIR" ]; then
    echo "❌ Integration directory doesn't exist!"
    echo "Expected: $INTEGRATE_WAVE_EFFORTS_DIR"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
    exit 1
fi

if [ ! -d "$INTEGRATE_WAVE_EFFORTS_DIR/.git" ]; then
    echo "❌ Integration directory is not a git repository!"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
    exit 1
fi

if [ ! -d "$INTEGRATE_WAVE_EFFORTS_DIR/.software-factory" ]; then
    echo "⚠️ Warning: .software-factory metadata directory missing"
fi

echo "✅ Directory structure validated"
```

### 6. RUN COMPREHENSIVE VALIDATION SCRIPT
```bash
echo "🔍 Running comprehensive infrastructure validation..."

bash $CLAUDE_PROJECT_DIR/utilities/validate-infrastructure.sh

VALIDATION_RESULT=$?

if [ $VALIDATION_RESULT -eq 911 ]; then
    echo "🔴🔴🔴 CATASTROPHIC: Wrong repository configured (R508 violation)"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
    exit 911
elif [ $VALIDATION_RESULT -ne 0 ]; then
    echo "❌ Infrastructure validation FAILED"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
    exit 1
fi

echo "✅ Comprehensive validation PASSED"
```

### 7. UPDATE STATE FILE WITH VALIDATION RESULTS
```bash
cd $CLAUDE_PROJECT_DIR

# Record validation success
jq '.integration_infrastructure.validation = {
  "timestamp": "'$(date -u +"%Y-%m-%dT%H:%M:%SZ")'",
  "status": "validated",
  "repository_verified": true,
  "branch_verified": true,
  "incremental_base_verified": true,
  "directory_verified": true,
  "r508_compliant": true,
  "r308_compliant": true
}' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json

echo "✅ Validation results recorded"
```

## 🚨🚨🚨 EXIT CONDITIONS [MANDATORY]
ONE of:
1. ✅ All validation checks PASSED → Continue to SPAWN_CODE_REVIEWER_MERGE_PLAN
2. ❌ Validation FAILED → Must transition to ERROR_RECOVERY
3. 🔴 CATASTROPHIC failure (exit code 911) → IMMEDIATE ERROR_RECOVERY with FALSE flag

## 🚨🚨🚨 STATE TRANSITIONS [BLOCKING]

### PROJECT_DONE PATH
If validation passes:
- Wave integration → `SPAWN_CODE_REVIEWER_MERGE_PLAN`
- Phase integration → `SPAWN_CODE_REVIEWER_PHASE_MERGE_PLAN`
- Project integration → `SPAWN_CODE_REVIEWER_PROJECT_MERGE_PLAN`

### FAILURE PATH
If validation fails:
- ANY failure → `ERROR_RECOVERY`
- Exit code 911 → `ERROR_RECOVERY` with CONTINUE-SOFTWARE-FACTORY=FALSE

## 🚨🚨 CRITICAL ERRORS TO PREVENT

### 1. **NEVER skip repository validation**
```bash
# ❌ CATASTROPHIC - Missing repository check
SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE
transition_to SPAWN_CODE_REVIEWER_MERGE_PLAN  # NO VALIDATION!
```

### 2. **NEVER ignore wrong repository**
```bash
# ❌ CATASTROPHIC - Ignoring wrong repo
if [ "$ACTUAL_REMOTE" != "$CONFIGURED_REPO" ]; then
    echo "Warning: Different repository"  # NO! THIS IS CATASTROPHIC!
    # Continuing anyway...  # NEVER DO THIS!
fi
```

### 3. **NEVER assume repository is correct**
```bash
# ❌ WRONG - Assuming without checking
echo "Integration infrastructure created, proceeding..."
# No actual validation!
```

## 🚨🚨🚨 GRADING PENALTIES
- Wrong repository (R508 violation): **-100% IMMEDIATE FAILURE**
- Wrong incremental base (R308 violation): **-100% IMMEDIATE FAILURE**
- Skipping validation: **-75%**
- Proceeding after validation failure: **-100%**
- Not transitioning to ERROR_RECOVERY on failure: **-75%**
- Not using CONTINUE-SOFTWARE-FACTORY=FALSE for catastrophic failures: **-50%**

## 📋 STATE OUTPUT REQUIREMENTS
At completion, MUST output:
```
=== VERIFY_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE Complete ===
Validation Status: [PASSED/FAILED/CATASTROPHIC]
Target Repository: [VERIFIED/WRONG]
Incremental Base: [VERIFIED/WRONG]
Branch Configuration: [VERIFIED/WRONG]
Directory Structure: [VERIFIED/MISSING]
R508 Compliance: [YES/NO]
R308 Compliance: [YES/NO]
Next State: [state_name]
===
```

## 🚨 MANDATORY STOP (R313)
After validation (success OR failure), you MUST:
1. Update orchestrator-state-v3.json with validation results
2. Output the state completion summary
3. **STOP IMMEDIATELY**
4. Do NOT continue to next state automatically

## 🔴🔴🔴 AUTOMATION FLAG (R405) 🔴🔴🔴

```bash
# After integration infrastructure validation
if [ "$VALIDATION_PASSED" = true ]; then
    echo "✅ Integration infrastructure validation passed"
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # Continue to merge plan
else
    echo "❌ Integration infrastructure validation failed"
    if [ "$CATASTROPHIC" = true ]; then
        echo "🔴🔴🔴 CATASTROPHIC FAILURE - WRONG REPOSITORY!"
        echo "CONTINUE-SOFTWARE-FACTORY=FALSE"  # R405: Legitimate use of FALSE
    else
        echo "Transitioning to error recovery"
        echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # System will handle recovery
    fi
fi
```

**CRITICAL**: Wrong repository (R508 violation) is one of the LEGITIMATE cases for CONTINUE-SOFTWARE-FACTORY=FALSE per R405, as it's an unrecoverable error requiring human intervention.

---

**REMEMBER**: This state is the LAST LINE OF DEFENSE against catastrophic repository targeting failures. NEVER skip or weaken these validations. The entire Software Factory depends on infrastructure being on the CORRECT repository!
## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

**MUST output as LAST line before state completion:**
- `CONTINUE-SOFTWARE-FACTORY=TRUE` - State completed, automation can continue
- `CONTINUE-SOFTWARE-FACTORY=FALSE` - Error/block, manual intervention required

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-CONTINUATION-FLAG-MASTER-GUIDE.md**
