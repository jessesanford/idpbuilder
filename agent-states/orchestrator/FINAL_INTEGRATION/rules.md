# ORCHESTRATOR STATE: FINAL_INTEGRATION

## 🔴🔴🔴 R290 ENFORCEMENT: READ THESE RULES FIRST! 🔴🔴🔴

**SUPREME LAW #3 (R290): STATE RULES MUST BE READ BEFORE STATE ACTIONS**

## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED FINAL_INTEGRATION STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_FINAL_INTEGRATION
echo "$(date +%s) - Rules read and acknowledged for FINAL_INTEGRATION" > .state_rules_read_orchestrator_FINAL_INTEGRATION
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY FINAL_INTEGRATION WORK UNTIL RULES ARE READ:
- ❌ Start final merges
- ❌ Start integration branch creation
- ❌ Start target repository cloning
- ❌ Update state files
- ❌ Continue to next state
- ❌ Think about what to do in this state

### ✅ YOU MUST IMMEDIATELY:
1. **READ** every rule in this file completely
2. **ACKNOWLEDGE** understanding of the state requirements
3. **CREATE** the verification marker shown above
4. **ONLY THEN** proceed with FINAL_INTEGRATION work

## State Overview
The FINAL_INTEGRATION state is where all phases of a project are merged into the target project's main branch. This is the culmination of all development efforts.

## Entry Conditions
- All phases completed successfully
- All phase integration branches exist and tested
- Architect approval for project completion
- No outstanding critical issues

## PROJECT INTEGRATION INFRASTRUCTURE SETUP

### 1. Target Repository Identification
```bash
# Extract target project information from configuration
TARGET_REPO_URL=$(yq '.target_repository.url' orchestrator-state.yaml)
TARGET_REPO_NAME=$(yq '.target_repository.name' orchestrator-state.yaml)
PROJECT_PREFIX=$(yq '.project_prefix' orchestrator-state.yaml)

# Validate target repository is NOT the software-factory
if [[ "$TARGET_REPO_URL" == *"software-factory"* ]]; then
    echo "❌ CRITICAL: Target repository cannot be software-factory!"
    echo "Target URL: $TARGET_REPO_URL"
    exit 1
fi
```

### 2. Fresh Integration Workspace Creation
```bash
# Create isolated integration workspace
INTEGRATION_DIR="/efforts/integration"
rm -rf "$INTEGRATION_DIR"
mkdir -p "$INTEGRATION_DIR"
cd "$INTEGRATION_DIR"

# Clone target repository (NOT software-factory)
git clone "$TARGET_REPO_URL" "$TARGET_REPO_NAME"
cd "$TARGET_REPO_NAME"

# CRITICAL SAFETY CHECK - Verify correct repository
REMOTE_URL=$(git remote get-url origin)
if [[ "$REMOTE_URL" == *"software-factory"* ]]; then
    echo "❌ CRITICAL: Cloned orchestrator repository instead of target!"
    echo "Expected: Target project repository"
    echo "Got: $REMOTE_URL"
    exit 1
fi

echo "✅ Verified: Working in target repository"
echo "Repository: $REMOTE_URL"
```

### 3. Integration Branch Setup
```bash
# Create final integration branch
INTEGRATION_BRANCH="${PROJECT_PREFIX}-final-integration"
git checkout -b "$INTEGRATION_BRANCH"

# Verify isolation
if [[ $(pwd) != */efforts/integration/* ]]; then
    echo "❌ ERROR: Not in isolated integration workspace!"
    exit 1
fi
```

## Core Rules

### R283: Project Integration Protocol (CRITICAL)
- Must use fresh clone in `/efforts/integration/`
- Must verify target repository is NOT software-factory
- Must create project-prefixed integration branch
- Must merge all phase branches sequentially

### R206: State Machine Validation (BLOCKING)
- Validate transition from PHASE_INTEGRATION
- Ensure all phases marked complete
- Update state file after successful integration

### R287: TODO Persistence (BLOCKING)
- Save TODOs before starting integration
- Save after each phase merge
- Save before final push
- Commit all TODO updates

## Required Actions

1. **Setup Integration Infrastructure**
   - Clone target repository to `/efforts/integration/`
   - Verify NOT software-factory repository
   - Create final integration branch

2. **Sequential Phase Merging**
   - Merge each phase branch in order
   - Run tests after each merge
   - Document any conflicts

3. **Final Validation**
   - Run complete test suite
   - Generate integration report
   - Update orchestrator state

## Exit Conditions

### Success Criteria
- All phase branches merged successfully
- All tests passing
- Integration branch pushed to remote
- State updated to PROJECT_COMPLETE

### Failure Handling
- If merge conflicts: Document and transition to ERROR_RECOVERY
- If tests fail: Revert and investigate
- If repository wrong: IMMEDIATE STOP

## State Transitions

### On Success
- Next State: PROJECT_COMPLETE
- Update: Mark project as integrated
- Action: Generate final report

### On Failure
- Next State: ERROR_RECOVERY
- Update: Document integration issues
- Action: Spawn architect for assessment

## Critical Checks

```bash
# Continuous verification during integration
verify_integration_safety() {
    # Check we're in correct directory
    if [[ $(pwd) != */efforts/integration/* ]]; then
        echo "❌ CRITICAL: Wrong working directory!"
        return 1
    fi
    
    # Check we're not in software-factory
    if git remote get-url origin | grep -q "software-factory"; then
        echo "❌ CRITICAL: Working in wrong repository!"
        return 1
    fi
    
    # Check branch naming
    if ! git branch --show-current | grep -q "$PROJECT_PREFIX"; then
        echo "⚠️ WARNING: Branch missing project prefix"
    fi
    
    return 0
}
```

## Monitoring Requirements

- Log all merge operations
- Track test results after each merge
- Monitor repository state continuously
- Verify isolation maintained

## Documentation Requirements

Generate and commit:
- `FINAL-INTEGRATION-REPORT.md`
- Updated `PROJECT-STATUS.md`
- Merge conflict resolutions (if any)
- Test results summary