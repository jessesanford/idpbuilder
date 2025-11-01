# INTEGRATE_WAVE_EFFORTS_BUILD_VALIDATION State Rules

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
## State Purpose
Validate that the integrated code builds successfully - the first major quality gate.

## Entry Conditions
- All branches successfully merged
- No unresolved conflicts
- Integration branch is complete
- Build validation is required (per requirements)

## Required Actions

### 1. Prepare Build Environment
```bash
# Clean any previous build artifacts
make clean 2>/dev/null || true
rm -rf build/ dist/ target/ 2>/dev/null || true

# Ensure all dependencies are available
if [[ -f "requirements.txt" ]]; then
    pip install -r requirements.txt
elif [[ -f "package.json" ]]; then
    npm install
elif [[ -f "go.mod" ]]; then
    go mod download
fi
```

### 2. Set Build Configuration
```bash
# Set build flags for integration testing
export BUILD_TYPE="integration"
export ENABLE_TESTS="true"
export VERBOSE_OUTPUT="true"

# Record build start time
BUILD_START=$(date +%s)
```

### 3. Capture Build Environment
```json
{
  "build_validation": {
    "started_at": "2025-01-21T10:30:00Z",
    "attempt": 2,
    "environment": {
      "build_type": "integration",
      "compiler": "gcc 9.3.0",
      "flags": "-O2 -Wall -Werror"
    }
  }
}
```

### 4. Transition Decision
- Proceed to INTEGRATE_WAVE_EFFORTS_RUN_BUILD to execute build

## Critical Validation Rules

### Build Must Be Clean
```bash
# No warnings allowed in integration
export CFLAGS="-Wall -Werror"
export GOFLAGS="-buildvcs=false"

# Treat warnings as errors
export CI=true
export INTEGRATE_WAVE_EFFORTS_BUILD=true
```

### Build Must Be Complete
- All modules must build
- All dependencies must resolve
- All artifacts must be generated

### Build Must Be Reproducible
```bash
# Clear ccache and other caches
ccache -C 2>/dev/null || true

# Use consistent timestamps
export SOURCE_DATE_EPOCH=$(date +%s)
```

## State Tracking Updates
```json
{
  "validation_status": {
    "merge": "COMPLETE",
    "build": "IN_PROGRESS",
    "unit_tests": "PENDING",
    "functional_tests": "PENDING"
  },
  "build_config": {
    "clean_build": true,
    "warnings_as_errors": true,
    "parallel_jobs": 4
  }
}
```

## Quality Gates

### Gate 1: Pre-Build Checks
- ✅ All merges complete
- ✅ No conflict markers in code
- ✅ Dependencies available
- ✅ Build system ready

### Gate 2: Build Configuration
- ✅ Appropriate flags set
- ✅ Error handling configured
- ✅ Output capture ready
- ✅ Timeout configured

### Gate 3: Environment Validation
- ✅ Correct branch checked out
- ✅ Clean workspace
- ✅ Tools available
- ✅ Disk space sufficient

## 🔴🔴🔴 R410: BUILD FAILURES DURING CASCADE (LAYERED CASCADE RECOGNITION) 🔴🔴🔴

### CRITICAL: Build Failures Are NOT Blockers During CASCADE!

**If this integration is part of CASCADE_REINTEGRATION:**

Build failures are **EXPECTED** and trigger **automatic new cascade layer creation** per R410.

```bash
# Check if we're in cascade mode
IN_CASCADE=$(jq -r '.cascade_mode // false' "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json")

if [[ "$IN_CASCADE" == "true" ]]; then
    echo "🔄 R410: Build failure during CASCADE mode"
    echo "This is EXPECTED (discovers cross-effort bugs)"
    echo "Orchestrator will automatically start new cascade layer"
    echo ""
    echo "BUILD FAILURE IS NOT A BLOCKER - IT'S A NEW BUG DISCOVERY!"
    echo ""

    # Report failure to orchestrator (will trigger R410 layered cascade)
    # Orchestrator handles decision-making per R410
    # NO manual intervention needed!

    exit 1  # Signal failure to orchestrator (triggers R410)
fi
```

**What Happens Next (R410 Automation):**
1. Orchestrator detects build failure
2. R410: Are these NEW bugs? (hash comparison with existing layers)
3. If YES → Start CASCADE LAYER N+1 (automated!)
4. If NO → Infinite loop detected (genuine block)

**Integration agent does NOT make layer decisions** - report failure and let orchestrator handle it per R410.

## Error Handling

### Pre-Build Failures
- Missing dependencies → Attempt to install
- Dirty workspace → Clean and retry
- Wrong branch → Checkout correct branch

### Build Preparation Failures
- Cannot clean → Log warning, continue
- Cannot set environment → INTEGRATE_WAVE_EFFORTS_ERROR
- Missing tools → INTEGRATE_WAVE_EFFORTS_ABORT

### Build Failures During CASCADE (R410)
- **IN CASCADE MODE**: Exit with failure code → Orchestrator starts new layer
- **NOT IN CASCADE**: Follow normal error recovery
- Integration agent reports, orchestrator decides (per R410)

## Logging Requirements
```bash
echo "[INTEGRATE_WAVE_EFFORTS_BUILD_VALIDATION] Preparing build validation"
echo "[INTEGRATE_WAVE_EFFORTS_BUILD_VALIDATION] Attempt: ${ATTEMPT_NUMBER}"
echo "[INTEGRATE_WAVE_EFFORTS_BUILD_VALIDATION] Build type: ${BUILD_TYPE}"
echo "[INTEGRATE_WAVE_EFFORTS_BUILD_VALIDATION] Merged branches: ${MERGED_COUNT}"

# Log build configuration
echo "Build Configuration:"
echo "  - Clean build: YES"
echo "  - Warnings as errors: YES"
echo "  - Parallel jobs: ${PARALLEL_JOBS:-4}"
echo "  - Timeout: ${BUILD_TIMEOUT:-30m}"
```

## Metrics to Track
- Build preparation time
- Configuration complexity
- Environment setup duration
- Cache hit rates

## Common Issues

### Issue 1: Dependency Conflicts
**Problem**: Merged branches have conflicting dependencies
**Solution**: Resolve in source branches, re-attempt integration

### Issue 2: Build Tool Version Mismatch
**Problem**: Different branches expect different tool versions
**Solution**: Use project-standard versions, fix in source branches

### Issue 3: Missing Build Configuration
**Problem**: No Makefile, package.json, or build config
**Solution**: Check if build is actually required, or fix in source

## 🔴🔴🔴 R291 + R331 MANDATORY DEMO CREATION (SUPREME LAW - BLOCKING) 🔴🔴🔴

**INTEGRATE_WAVE_EFFORTS AGENT MUST CREATE DEMO SCRIPTS AFTER PROJECT_DONEFUL BUILD**

**🚨 CRITICAL: All demos MUST comply with R331 (Demo Validation Protocol) - NO SIMULATION ALLOWED! 🚨**

### Demo Creation Requirements (R291 + R331)

After build completes successfully, integration agent MUST create demo scripts and documentation per the integration plan:

```bash
# 🔴🔴🔴 R504 + R330: USE PRE-PLANNED DEMO INFORMATION - NO RUNTIME DECISIONS! 🔴🔴🔴

# Determine integration level and load pre-planned demo info
INTEGRATE_WAVE_EFFORTS_TYPE=$(jq -r '.integration_type' "$INTEGRATE_WAVE_EFFORTS_DIR/integration-state.json")
CURRENT_PHASE=$(jq -r '.current_phase' orchestrator-state-v3.json 2>/dev/null || echo "1")
CURRENT_WAVE=$(jq -r '.current_wave' orchestrator-state-v3.json 2>/dev/null || echo "1")

# READ PRE-PLANNED DEMO INFORMATION from pre_planned_infrastructure (R504 ENFORCEMENT)
if [ "$INTEGRATE_WAVE_EFFORTS_TYPE" = "wave" ]; then
    INTEGRATE_WAVE_EFFORTS_KEY="phase${CURRENT_PHASE}_wave${CURRENT_WAVE}"
    DEMO_SCRIPT_FILE=$(jq -r ".pre_planned_infrastructure.integrations.wave_integrations.\"${INTEGRATE_WAVE_EFFORTS_KEY}\".demo_script_file // empty" orchestrator-state-v3.json)
    DEMO_DESCRIPTION=$(jq -r ".pre_planned_infrastructure.integrations.wave_integrations.\"${INTEGRATE_WAVE_EFFORTS_KEY}\".demo_description // empty" orchestrator-state-v3.json)
    DEMO_PLAN_FILE=$(jq -r ".pre_planned_infrastructure.integrations.wave_integrations.\"${INTEGRATE_WAVE_EFFORTS_KEY}\".demo_plan_file // empty" orchestrator-state-v3.json)
elif [ "$INTEGRATE_WAVE_EFFORTS_TYPE" = "phase" ]; then
    INTEGRATE_WAVE_EFFORTS_KEY="phase${CURRENT_PHASE}"
    DEMO_SCRIPT_FILE=$(jq -r ".pre_planned_infrastructure.integrations.phase_integrations.\"${INTEGRATE_WAVE_EFFORTS_KEY}\".demo_script_file // empty" orchestrator-state-v3.json)
    DEMO_DESCRIPTION=$(jq -r ".pre_planned_infrastructure.integrations.phase_integrations.\"${INTEGRATE_WAVE_EFFORTS_KEY}\".demo_description // empty" orchestrator-state-v3.json)
    DEMO_PLAN_FILE=$(jq -r ".pre_planned_infrastructure.integrations.phase_integrations.\"${INTEGRATE_WAVE_EFFORTS_KEY}\".demo_plan_file // empty" orchestrator-state-v3.json)
elif [ "$INTEGRATE_WAVE_EFFORTS_TYPE" = "project" ]; then
    DEMO_SCRIPT_FILE=$(jq -r ".pre_planned_infrastructure.integrations.project_integration.demo_script_file // empty" orchestrator-state-v3.json)
    DEMO_DESCRIPTION=$(jq -r ".pre_planned_infrastructure.integrations.project_integration.demo_description // empty" orchestrator-state-v3.json)
    DEMO_PLAN_FILE=$(jq -r ".pre_planned_infrastructure.integrations.project_integration.demo_plan_file // empty" orchestrator-state-v3.json)
fi

# VALIDATE pre-planned demo information exists (R504 + R330 ENFORCEMENT)
if [ -z "$DEMO_SCRIPT_FILE" ] || [ "$DEMO_SCRIPT_FILE" = "null" ]; then
    echo "❌ FATAL: No pre-planned demo_script_file for ${INTEGRATE_WAVE_EFFORTS_TYPE} integration!"
    echo "  R504 VIOLATION: Demo paths must be pre-planned in pre_planned_infrastructure"
    echo "  R330 VIOLATION: Demos must be planned before integration"
    exit 1
fi

echo "✅ Using pre-planned demo information from orchestrator-state-v3.json:"
echo "  Demo Script: $DEMO_SCRIPT_FILE"
echo "  Description: $DEMO_DESCRIPTION"
echo "  Plan File: $DEMO_PLAN_FILE"

# Extract demo directory from pre-planned script path
DEMO_DIR=$(dirname "$DEMO_SCRIPT_FILE")
DEMO_SCRIPT_NAME=$(basename "$DEMO_SCRIPT_FILE")

echo "📝 Creating demos in pre-planned location: $DEMO_DIR"
mkdir -p "$DEMO_DIR"
cd "$DEMO_DIR"
```

### Create Main Demo Script

```bash
# Create integration demo script using PRE-PLANNED NAME and based on PRE-PLANNED DEMO PLAN
# Read demo plan if it exists to get detailed steps
DEMO_STEPS="# TODO: Implement demo steps from demo plan: $DEMO_PLAN_FILE"
if [ -f "$DEMO_PLAN_FILE" ]; then
    echo "✅ Found pre-planned demo plan: $DEMO_PLAN_FILE"
    DEMO_STEPS=$(grep -A 100 "Demo Steps\|Demo Scenarios\|How to Run" "$DEMO_PLAN_FILE" | head -50 || echo "# TODO: Extract steps from $DEMO_PLAN_FILE")
fi

# 🚨 R331 REQUIREMENT: Perform pre-demo implementation scan
echo "🔍 R331: Pre-demo implementation completeness scan..."
if find src/ pkg/ internal/ -type f \( -name "*.go" -o -name "*.ts" -o -name "*.py" -o -name "*.rs" \) \
    -exec grep -l "TODO\|FIXME\|XXX\|STUB" {} \; 2>/dev/null | grep -q .; then
    echo "❌ R331 VIOLATION: TODO/FIXME/STUB found in implementation!"
    echo "PENALTY: -100% (implementation incomplete)"
    echo "Cannot create demo for incomplete implementation"
    exit 331
fi
echo "✅ R331: Implementation completeness check passed"

# Create demo script with PRE-PLANNED filename (R331 COMPLIANT)
cat > "$DEMO_SCRIPT_NAME" << EOF
#!/bin/bash
set -e  # R331 requirement: Exit on any error

# Integration Demo Script
# Auto-generated by Integration Agent per R291 + R331
# Description: ${DEMO_DESCRIPTION}
# Based on plan: ${DEMO_PLAN_FILE}

echo "🎬 Starting Integration Demo"
echo "=============================="
echo "Description: ${DEMO_DESCRIPTION}"
echo ""

# R331 REQUIREMENT: Pre-demo validation
echo "🔍 R331: Pre-demo validation..."
if find src/ pkg/ internal/ -type f -exec grep -l "TODO\|FIXME" {} \; 2>/dev/null | grep -q .; then
    echo "❌ Implementation incomplete (TODO found)"
    exit 1
fi

# Demo steps from pre-planned demo plan
${DEMO_STEPS}

# Example structure (customize based on integration plan):
# 1. Setup environment
# 2. Run integrated features (R331: real execution, not simulation)
# 3. Verify outputs (R331: external state verification required)
# 4. Cleanup

# R331 REQUIREMENT: External verification example
# TODO: Customize based on your integration
# Example: Verify image in registry, file exists, database record, etc.
# if ! docker images | grep -q "expected-image"; then
#     echo "❌ DEMO FAILED: Expected artifact not found"
#     exit 1
# fi

echo ""
echo "✅ Demo completed successfully"
exit 0
EOF

chmod +x "$DEMO_SCRIPT_NAME"
echo "✅ Created ${DEMO_SCRIPT_NAME} using pre-planned name from orchestrator-state-v3.json"
echo "✅ Demo is R331 compliant (includes validation, error handling, external checks)"
```

### Create Demo Documentation

```bash
# Create demo documentation
cat > DEMO.md << EOF
# ${INTEGRATE_WAVE_EFFORTS_TYPE^} Integration Demo

## Integration Summary
- Integration Type: ${INTEGRATE_WAVE_EFFORTS_TYPE}
- Build Status: ✅ PASSING
- Tests Status: $([ -f "test-results.log" ] && echo "✅ PASSING" || echo "⏳ PENDING")
- Created: $(date -Iseconds)

## Features Demonstrated
<!-- List the features being demonstrated -->
1. [Feature 1]: Brief description
2. [Feature 2]: Brief description
3. [Feature 3]: Brief description

## How to Run Demo

\`\`\`bash
# Execute the demo
./integration-demo.sh

# Expected output:
# - Demo starts successfully
# - Features execute without errors
# - Demo completes with exit code 0
\`\`\`

## Demo Success Criteria
- Demo script exits with code 0
- No errors during execution
- Features demonstrate integration working correctly

## Evidence
- Build log: ../../build.log
- Test results: ../../test-results.log
- Demo output: (captured when demo runs)

---
**Generated by Integration Agent**
**Per R291: Integration Demo Requirement**
EOF

echo "✅ Created DEMO.md"
```

### Commit Demos to Repository

```bash
# Return to integration root
cd "$OLDPWD"

# Add demos to git
git add "$DEMO_DIR/"

# Commit demos with integration code
git commit -m "feat: add ${INTEGRATE_WAVE_EFFORTS_TYPE} integration demos per R291

Created demo scripts and documentation in ${DEMO_DIR}
- integration-demo.sh: Main demo script
- DEMO.md: Demo documentation and instructions

Demos created after successful build validation.
Required by R291 for all integrations.

Demo directory: ${DEMO_DIR}"

git push

echo "✅ Demos committed and pushed"
```

### Demo Creation Validation

```bash
# Verify demos were created successfully
if [ ! -f "$DEMO_DIR/integration-demo.sh" ]; then
    echo "🔴 FAILED: Demo script not created!"
    exit 291
fi

if [ ! -f "$DEMO_DIR/DEMO.md" ]; then
    echo "🔴 FAILED: Demo documentation not created!"
    exit 291
fi

if [ ! -x "$DEMO_DIR/integration-demo.sh" ]; then
    echo "🔴 FAILED: Demo script not executable!"
    exit 291
fi

echo "✅ Demo creation validation passed"
echo "   Demos created in: $DEMO_DIR"
echo "   Demo script: integration-demo.sh"
echo "   Demo documentation: DEMO.md"
```

**IMPORTANT:** Demos are created WITH the integration code, not separately.
They are committed and pushed as part of the integration branch.

## 🔴🔴🔴 R291 MANDATORY DEMO GATE CHECK (SUPREME LAW - BLOCKING) 🔴🔴🔴

**THIS IS A BLOCKING GATE - INTEGRATE_WAVE_EFFORTS CANNOT PROCEED WITHOUT DEMOS**

### DEMO VALIDATION MUST HAPPEN AFTER BUILD

```bash
# After build completes successfully, MANDATORY R291 gate check
verify_integration_demos() {
    echo "🔴🔴🔴 R291 MANDATORY DEMO GATE CHECK 🔴🔴🔴"
    local DEMO_FAILED=false
    local FAILURE_DETAILS=""

    # 1. Check for wave-level demo script
    if [ ! -f "demo-wave-integration.sh" ] && [ ! -f "demo-integration.sh" ]; then
        echo "🔴 R291 VIOLATION: Missing wave integration demo script!"
        DEMO_FAILED=true
        FAILURE_DETAILS="${FAILURE_DETAILS}\n- Missing demo script (demo-wave-integration.sh or demo-integration.sh)"
    fi

    # 2. Check for demo documentation
    if [ ! -f "WAVE-DEMO.md" ] && [ ! -f "DEMO.md" ] && [ ! -f "INTEGRATE_WAVE_EFFORTS-DEMO.md" ]; then
        echo "🔴 R291 VIOLATION: Missing demo documentation!"
        DEMO_FAILED=true
        FAILURE_DETAILS="${FAILURE_DETAILS}\n- Missing demo documentation (WAVE-DEMO.md, DEMO.md, or INTEGRATE_WAVE_EFFORTS-DEMO.md)"
    fi

    # 3. Test demo execution
    DEMO_SCRIPT=""
    [ -f "demo-wave-integration.sh" ] && DEMO_SCRIPT="demo-wave-integration.sh"
    [ -f "demo-integration.sh" ] && DEMO_SCRIPT="demo-integration.sh"

    if [ -n "$DEMO_SCRIPT" ]; then
        echo "🧪 Testing demo execution: $DEMO_SCRIPT"
        chmod +x "$DEMO_SCRIPT" 2>/dev/null || true

        if ! bash "$DEMO_SCRIPT"; then
            echo "🔴 R291 VIOLATION: Demo script fails to execute!"
            DEMO_FAILED=true
            FAILURE_DETAILS="${FAILURE_DETAILS}\n- Demo script execution failed"
        else
            echo "✅ Demo script executes successfully"
        fi
    fi

    # 4. BLOCKING: Fail integration if demos missing/broken
    if [ "$DEMO_FAILED" = true ]; then
        echo ""
        echo "🔴🔴🔴 R291 GATE CHECK: FAILED 🔴🔴🔴"
        echo "Integration CANNOT proceed without working demos"
        echo ""
        echo "MISSING/BROKEN ARTIFACTS:"
        echo -e "$FAILURE_DETAILS"
        echo ""
        echo "REQUIRED ACTION: Create demo artifacts before integration"
        echo ""

        # Document failure in DEMO-STATUS.md
        cat > DEMO-STATUS.md << EOF
# R291 Demo Gate Check - FAILED

Status: ❌ BLOCKING FAILURE
Timestamp: $(date -Iseconds)

## Missing/Broken Artifacts:
${FAILURE_DETAILS}

## Required Actions:
1. Create demo-wave-integration.sh script (or demo-integration.sh)
2. Create WAVE-DEMO.md documentation (or DEMO.md/INTEGRATE_WAVE_EFFORTS-DEMO.md)
3. Verify demo executes successfully
4. Re-run integration

## Next Step:
**MANDATORY**: Orchestrator MUST transition to ERROR_RECOVERY state.
This is a BLOCKING gate per R291 - integration cannot proceed.

See: rule-library/R291-integration-demo-requirement.md
EOF

        # Commit demo status for visibility
        git add DEMO-STATUS.md
        git commit -m "error: R291 demo gate check FAILED - blocking integration" || true
        git push || true

        echo ""
        echo "📝 DEMO-STATUS.md created documenting R291 violation"
        echo "🔴 MUST transition to ERROR_RECOVERY - this is MANDATORY per R291"
        echo ""

        # Exit with R291 code
        exit 291
    fi

    echo ""
    echo "✅✅✅ R291 GATE CHECK: PASSED ✅✅✅"
    echo "Demo artifacts present and functional"
    echo ""
}

# MANDATORY: Run R291 gate check BEFORE marking integration complete
verify_integration_demos
```

### MANDATORY ERROR_RECOVERY TRIGGER

**If R291 gate check fails:**
1. Integration agent MUST set status = "DEMO_REQUIREMENTS_MISSING"
2. Integration agent MUST exit with code 291
3. Orchestrator MUST detect R291 failure in integration report
4. Orchestrator MUST transition to ERROR_RECOVERY (NOT optional)
5. Integration CANNOT be approved until demos created

**Penalty for proceeding without demos: -100% IMMEDIATE FAILURE**

**Related Rules:**
- R330: Demo Planning Requirements (ensures demos are planned in efforts)
- R363: Sequential Direct Mergeability (integration is for testing only)
- R364: Integration Testing Only Branches (integration branches never merge)

## Success Criteria
✅ Build environment prepared
✅ All dependencies available
✅ Build configuration set
✅ **R291: Demo artifacts present and working** ← NEW MANDATORY
✅ Ready to execute build

## Next State
- If R291 gate passes → Continue to next integration step
- If R291 gate fails → EXIT 291 (triggers ERROR_RECOVERY in orchestrator)

## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

### YOU MUST OUTPUT THE CONTINUATION FLAG AS YOUR LAST ACTION

**EVERY STATE COMPLETION MUST END WITH EXACTLY ONE OF:**
```bash
# If state completed successfully and factory should continue:
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"

# If error/block/manual review needed:
echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
```

### CRITICAL REQUIREMENTS:
1. **ABSOLUTE LAST OUTPUT**: This MUST be the very last line of output before state completion
2. **NO TEXT AFTER**: No explanations, summaries, or any text after the flag
3. **EXACTLY AS SHOWN**: Use exact format - no variations like CONTINUE-ORCHESTRATING
4. **ALWAYS REQUIRED**: Every single state must output this flag
5. **GREPPABLE**: Must be on its own line for automation parsing

## 🎯 R405 DECISION GUIDE FOR INTEGRATE_WAVE_EFFORTS AGENTS

**DEFAULT: TRUE** - Integration agents should default to TRUE unless TRULY exceptional

### ✅ USE TRUE (Normal Integration Operations):

#### Scenario 1: Build Failures with Known Fixes (CASCADE)
```bash
# Build failed, but BUILD-FIX-SUMMARY.md has the fix
if detect_cascade_context && search_build_fix_summaries "$build_error"; then
    apply_known_fix "$build_error"  # R521 protocol
    echo "✅ Applied known fix - continuing"
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # ← CORRECT
fi
```

#### Scenario 2: Demo Creation Needed
```bash
# Demos don't exist yet - CREATE THEM per R291
if [ ! -f "$DEMO_SCRIPT" ]; then
    create_integration_demos  # R291 requirement
    echo "✅ Created demos - continuing"
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # ← CORRECT
fi
```

#### Scenario 3: Build/Test Failures (NEW BUGS - First Occurrence)
```bash
# Build failed with NEW bug (not in BUILD-FIX-SUMMARY)
if ! search_build_fix_summaries "$build_error"; then
    document_new_bug "$build_error"  # R266 protocol
    echo "🔴 New bug documented - orchestrator will spawn fix cascade"
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # ← CORRECT
    # Orchestrator handles cascade - integration agent continues to next state
fi
```

#### Scenario 4: Duplicate Bug Detection
```bash
# Same bug as before - update occurrence count
if check_for_duplicate_bug "$bug_signature"; then
    update_duplicate_bug "$bug_signature"  # R522 protocol
    echo "✅ Updated duplicate bug - continuing"
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # ← CORRECT
fi
```

#### Scenario 5: Integration Complete with Documented Issues
```bash
# Integration done, some bugs documented for upstream fix
if [ $NEW_BUGS_COUNT -gt 0 ]; then
    echo "✅ Integration complete, $NEW_BUGS_COUNT new bugs documented"
    echo "   Orchestrator will coordinate fixes per R300"
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # ← CORRECT
fi
```

### ❌ USE FALSE (Exceptional Situations ONLY):

#### Scenario 1: Corrupted State Files
```bash
# orchestrator-state-v3.json is missing or corrupted
if [ ! -f "orchestrator-state-v3.json" ] || ! jq empty orchestrator-state-v3.json 2>/dev/null; then
    echo "❌ CRITICAL: orchestrator-state-v3.json corrupted"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"  # ← CORRECT - truly exceptional
fi
```

#### Scenario 2: Multiple Fix Cascades with Recurring Issues (4+)
```bash
# This is the 4th+ attempt and same bugs keep appearing
ATTEMPT_COUNT=$(ls -1 .software-factory/INTEGRATE_WAVE_EFFORTS-REPORT-*.md 2>/dev/null | wc -l)
if [ $ATTEMPT_COUNT -ge 4 ]; then
    echo "❌ EXCEPTIONAL: 4+ integration attempts with recurring failures"
    echo "   Pattern indicates fundamental architectural issue"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"  # ← CORRECT - need human review
fi
```

#### Scenario 3: Cannot Create Demos (R504 Pre-Planning Missing)
```bash
# R504 pre-planned demo info is missing from orchestrator-state-v3.json
if [ -z "$DEMO_SCRIPT_FILE" ] || [ "$DEMO_SCRIPT_FILE" = "null" ]; then
    echo "❌ FATAL: No pre-planned demo info in orchestrator-state-v3.json"
    echo "   R504 violation - demos must be pre-planned"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"  # ← CORRECT - planning failure
fi
```

#### Scenario 4: Wrong Working Directory/Branch
```bash
# Integration agent is in completely wrong location
if [ "$(git rev-parse --abbrev-ref HEAD)" != "$EXPECTED_INTEGRATE_WAVE_EFFORTS_BRANCH" ]; then
    echo "❌ FATAL: Wrong branch - expected $EXPECTED_INTEGRATE_WAVE_EFFORTS_BRANCH"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"  # ← CORRECT - env issue
fi
```

### ⚠️ COMMON MISTAKES TO AVOID:

```bash
# ❌ WRONG: Setting FALSE for normal build failures
echo "Build failed - need to fix upstream"
echo "CONTINUE-SOFTWARE-FACTORY=FALSE"  # ← WRONG!
# Should be TRUE - orchestrator handles fix cascade automatically

# ❌ WRONG: Setting FALSE when demos are missing
echo "Demos don't exist - can't continue"
echo "CONTINUE-SOFTWARE-FACTORY=FALSE"  # ← WRONG!
# Should CREATE demos per R291, then TRUE

# ❌ WRONG: Setting FALSE for known fix application
echo "Found known fix from BUILD-FIX-SUMMARY"
echo "But I'm not sure if I should apply it"
echo "CONTINUE-SOFTWARE-FACTORY=FALSE"  # ← WRONG!
# Should apply known fix per R521, then TRUE

# ❌ WRONG: Setting FALSE for first/second/third cascade
echo "This is the 2nd integration attempt and it still failed"
echo "CONTINUE-SOFTWARE-FACTORY=FALSE"  # ← WRONG!
# 1-3 cascades are NORMAL - only 4+ is exceptional
```

### 🔑 THE KEY PRINCIPLE:

**ASK**: "Can the system proceed automatically to handle this?"

- Build failures → Orchestrator spawns fix cascade → **TRUE**
- Missing demos → Integration agent creates them → **TRUE**
- Known fixes available → Integration agent applies them → **TRUE**
- New bugs discovered → Orchestrator coordinates fixes → **TRUE**
- Duplicate bugs → Integration agent updates registry → **TRUE**
- Corrupted state files → System cannot proceed → **FALSE**
- 4+ cascades with same issues → Need human analysis → **FALSE**
- R504 pre-planning missing → Planning failed → **FALSE**

**Default to TRUE unless the system is STUCK and CANNOT proceed!**

### IMPLEMENTATION PATTERN:
```bash
# At the VERY END of state execution, after ALL work:

# Determine success/failure
if [[ "$STATE_COMPLETED_PROJECT_DONEFULLY" == "true" ]]; then
    echo "✅ State work completed successfully"
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
else
    echo "❌ State encountered issues requiring intervention"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
fi

# DO NOT OUTPUT ANYTHING AFTER THE FLAG!
```

### GRADING IMPACT:
- **Missing flag**: -100% AUTOMATIC FAILURE
- **Text after flag**: -50% penalty
- **Wrong format**: -100% AUTOMATIC FAILURE
- **Multiple flags**: -50% penalty
- **Incorrect TRUE/FALSE decision**: -20% penalty

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

