# Orchestrator - PHASE_COMPLETE State Rules

## 🛑🛑🛑 R322 MANDATORY STOP BEFORE STATE TRANSITIONS 🛑🛑🛑

**SUPREME LAW - VIOLATION = -100% IMMEDIATE FAILURE**

### YOU MUST STOP AFTER:
1. ✅ Completing all TODOs for this state
2. ✅ Updating orchestrator-state.yaml with new state
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

**YOU HAVE ENTERED PHASE_COMPLETE STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_PHASE_COMPLETE
echo "$(date +%s) - Rules read and acknowledged for PHASE_COMPLETE" > .state_rules_read_orchestrator_PHASE_COMPLETE
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY PHASE_COMPLETE WORK UNTIL RULES ARE READ:
- ❌ Start finalize phase work
- ❌ Start update documentation
- ❌ Start prepare next phase
- ❌ Update state files
- ❌ Continue to next state
- ❌ Think about what to do in this state

### ✅ YOU MUST IMMEDIATELY:

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

### ❌ ANTI-PATTERNS THAT WILL CAUSE FAILURE:

1. **Fake Acknowledgment Without Reading**:
   ```
   ❌ WRONG: "I acknowledge R151, R208, R053..."
   (No Read tool calls detected = AUTOMATIC FAILURE)
   ```

2. **Bulk Acknowledgment**:
   ```
   ❌ WRONG: "I acknowledge all PHASE_COMPLETE rules"
   (YOU Must READ AND ACKNOWLEDGE EACH rule individually)
   ```

3. **Silent Reading**:
   ```
   ❌ WRONG: [Reads rules but doesn't acknowledge]
   "Now I've read the rules, let me start work..."
   (MUST explicitly acknowledge EACH rule)
   ```

4. **Reading From Memory**:
   ```
   ❌ WRONG: "I know R208 requires CD before spawn..."
   (Must READ from file, not recall from memory)
   ```

5. **Skipping Rules in PRIMARY DIRECTIVES**:
   ```
   ❌ WRONG: Reading only some rules from the list
   (ALL rules in PRIMARY DIRECTIVES are MANDATORY)
   ```

### ✅ CORRECT PATTERN FOR PHASE_COMPLETE:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/[first-rule-file].md
2. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/[second-rule-file].md  
4. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. "Ready to execute PHASE_COMPLETE work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY PHASE_COMPLETE work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ You have stated readiness to execute PHASE_COMPLETE work
1. **READ** every rule file listed in PRIMARY DIRECTIVES below
2. **ACKNOWLEDGE** each rule individually with number and description
3. **ONLY THEN** proceed with PHASE_COMPLETE work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY PHASE_COMPLETE work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

**THE SYSTEM IS MONITORING YOUR READ TOOL CALLS!**

## ⚠️⚠️⚠️ MANDATORY RULE READING AND ACKNOWLEDGMENT ⚠️⚠️⚠️

**YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. YOUR READ TOOL CALLS ARE BEING MONITORED.**

### ❌ ANTI-PATTERNS THAT WILL CAUSE FAILURE:
1. Fake acknowledgment without reading
2. Bulk acknowledgment
3. Reading from memory

### ✅ CORRECT PATTERN:
1. READ each rule file
2. Acknowledge individually with rule number and description

## 📋 PRIMARY DIRECTIVES FOR PHASE_COMPLETE STATE

### 🚨🚨🚨 R257 - Mandatory Phase Assessment Report
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R257-mandatory-phase-assessment-report.md`
**Criticality**: BLOCKING - Required for phase completion
**Summary**: Architect must create assessment report before phase completion

### 🔴🔴🔴 R288 - State File Update and Commit Protocol (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
**Criticality**: SUPREME LAW - Update on every transition
**Summary**: Update orchestrator-state.yaml with phase completion metrics

### 🚨🚨🚨 R288 - State File Update and Commit Protocol
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
**Criticality**: BLOCKING - Push within 60 seconds
**Summary**: Commit and push state file immediately after updates

### 🚨🚨🚨 R035 - Phase Completion Testing
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R035-phase-completion-testing.md`
**Criticality**: BLOCKING - All tests must pass
**Summary**: Run full test suite at phase level before completion

### 🚨🚨🚨 R040 - Documentation Requirements
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R040-documentation-requirements.md`
**Criticality**: BLOCKING - Complete documentation required
**Summary**: Generate phase completion report and documentation

### 🔴🔴🔴 R322 - Mandatory Stop Before State Transitions (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-before-state-transitions.md`
**Criticality**: SUPREME LAW - Continue until SUCCESS or next phase
**Summary**: Keep working through phase completion tasks

## 🚨 PHASE_COMPLETE IS A VERB - FINALIZE PHASE NOW! 🚨

### IMMEDIATE ACTIONS UPON ENTERING PHASE_COMPLETE

**THE MOMENT YOU ENTER THIS STATE, YOU MUST:**
1. Create final phase integration branch NOW
2. Document all phase achievements
3. Generate phase completion report
4. Determine if more phases exist

**FORBIDDEN - AUTOMATIC FAILURE:**
- ❌ "STATE TRANSITION COMPLETE: Phase is complete" [stops]
- ❌ "Successfully completed the phase" [celebrates early]
- ❌ "Phase done, waiting for next steps" [pauses]

**REQUIRED - IMMEDIATE ACTION:**
- ✅ "Creating phase integration branch..."
- ✅ "Documenting phase achievements..."
- ✅ "Generating completion report..."
- ✅ "Checking for additional phases..."

## 🔴🔴🔴 MANDATORY MEASUREMENT RULE - R304 🔴🔴🔴

### R304: Mandatory Line Counter Tool Enforcement
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R304-mandatory-line-counter-enforcement.md`
**Criticality**: BLOCKING - Manual counting = AUTOMATIC -100% FAILURE

**ABSOLUTE REQUIREMENTS:**
- ✅ MUST use `$CLAUDE_PROJECT_DIR/tools/line-counter.sh` for ALL line counting
- ❌ NEVER use `wc -l` or any manual counting method
- ❌ NEVER count lines any other way - this is a -100% automatic failure
- ✅ MUST specify both -b (base branch) and -c (current branch) parameters
- ✅ Base branch MUST be phase integration branch (NOT "main")

**Failure to use the line counter tool = AUTOMATIC -100% GRADE**

## State Context

This state is entered ONLY after architect approves the complete phase AND creates the mandatory assessment report (R257). It handles final phase-level tasks before SUCCESS or next phase.

## 🚨🚨🚨 PREREQUISITE: Verified Phase Assessment Report per R257 🚨🚨🚨

**SEE**: `$CLAUDE_PROJECT_DIR/rule-library/R257-mandatory-phase-assessment-report.md`

**CRITICAL**: You can ONLY be in this state if:
1. Architect completed phase assessment
2. Assessment report file exists at `phase-assessments/phase{N}/PHASE-{N}-ASSESSMENT-REPORT.md`
3. Report contains PHASE_COMPLETE decision
4. Report was verified in WAITING_FOR_PHASE_ASSESSMENT state

**The assessment report path MUST be in your state file:**
```yaml
phase_assessment:
  report_file: "phase-assessments/phase1/PHASE-1-ASSESSMENT-REPORT.md"
  decision: "PHASE_COMPLETE"
```

## 🔴🔴🔴 CRITICAL: Post-Assessment Actions 🔴🔴🔴

**This state performs MANDATORY phase finalization:**
- Create final phase integration branch
- Document all deliverables
- Update metrics and reporting
- Prepare for SUCCESS or next phase

## Primary Purpose

The PHASE_COMPLETE state is for:
1. Creating phase-level integration branch
2. Documenting phase achievements
3. Generating phase completion metrics
4. Determining single vs multi-phase flow
5. Transitioning to SUCCESS or next phase

## Phase Integration Tasks (FOLLOWS R104 AND R282)

### 🔴🔴🔴 SETUP PHASE INTEGRATION INFRASTRUCTURE (R104/R282) 🔴🔴🔴

**YOU MUST FOLLOW R104 AND R282 - Integration branches are created in TARGET repository!**

```bash
# 1. Read target repository configuration (R104 requirement)
TARGET_CONFIG="$CLAUDE_PROJECT_DIR/target-repo-config.yaml"
TARGET_REPO_URL=$(yq '.target_repository.url' "$TARGET_CONFIG")
TARGET_REPO_NAME=$(yq '.target_repository.name' "$TARGET_CONFIG")
DEFAULT_BRANCH=$(yq '.target_repository.default_branch // "main"' "$TARGET_CONFIG")

# 2. Verify target repository exists
if [ ! -f "$TARGET_CONFIG" ]; then
    echo "❌ CRITICAL: target-repo-config.yaml not found!"
    echo "Cannot proceed without target repository configuration"
    exit 1
fi

if [ "$TARGET_REPO_URL" = "null" ] || [ -z "$TARGET_REPO_URL" ]; then
    echo "❌ CRITICAL: No target repository URL configured!"
    echo "Please configure target_repository.url in target-repo-config.yaml"
    exit 1
fi

# 3. Create phase integration workspace (R282 directory structure)
PHASE=$(yq '.current_phase' orchestrator-state.yaml)
PHASE_INTEGRATION_DIR="$CLAUDE_PROJECT_DIR/efforts/phase${PHASE}/integration-workspace"
echo "📂 Creating phase integration workspace at: $PHASE_INTEGRATION_DIR"

# Clean previous attempts if any
if [ -d "$PHASE_INTEGRATION_DIR" ]; then
    echo "🧹 Cleaning previous integration workspace..."
    rm -rf "$PHASE_INTEGRATION_DIR"
fi

mkdir -p "$(dirname "$PHASE_INTEGRATION_DIR")"

# 4. Clone TARGET repository (NOT software-factory!) - R104/R282 CRITICAL
echo "📦 Cloning TARGET repository for phase integration..."
echo "   Repository: $TARGET_REPO_URL"
echo "   Destination: $PHASE_INTEGRATION_DIR"

git clone "$TARGET_REPO_URL" "$PHASE_INTEGRATION_DIR"

if [ $? -ne 0 ]; then
    echo "❌ Failed to clone target repository!"
    echo "Please verify repository URL and access permissions"
    exit 1
fi

cd "$PHASE_INTEGRATION_DIR"

# 5. CRITICAL SAFETY CHECK - Verify we're in the TARGET repository (R282)
REMOTE_URL=$(git remote get-url origin)
if [[ "$REMOTE_URL" == *"software-factory"* ]] || [[ "$REMOTE_URL" == *"planning"* ]]; then
    echo "🔴🔴🔴 CRITICAL ERROR: Cloned wrong repository!"
    echo "Expected: Target project repository"
    echo "Got: $REMOTE_URL"
    echo "This violates R104 and R282!"
    exit 104
fi

echo "✅ Verified: Working in target repository"
echo "   Remote URL: $REMOTE_URL"

# 6. Create phase integration branch in TARGET repo (R104 naming)
git checkout "$DEFAULT_BRANCH"
git pull origin "$DEFAULT_BRANCH"
PHASE_BRANCH="phase-${PHASE}-integration"
echo "🌿 Creating phase integration branch: $PHASE_BRANCH"
git checkout -b "$PHASE_BRANCH"

# 7. Merge all wave integration branches (R282 sequential merging)
echo "📦 Starting wave integration merges..."
TOTAL_WAVES=$(yq ".phases[] | select(.phase_number == $PHASE) | .total_waves // 1" orchestrator-state.yaml)

for wave in $(seq 1 $TOTAL_WAVES); do
    WAVE_BRANCH="wave-${wave}-integration"
    echo "   Merging Wave $wave: $WAVE_BRANCH"
    
    # Fetch and merge
    git fetch origin "$WAVE_BRANCH"
    if git merge "origin/$WAVE_BRANCH" --no-ff -m "Integrate Wave $wave into Phase $PHASE"; then
        echo "   ✅ Wave $wave merged successfully"
    else
        echo "   ❌ Merge conflict in Wave $wave - manual resolution required"
        exit 1
    fi
    
    # Test after each wave (R282 requirement)
    echo "   🧪 Running tests after Wave $wave integration..."
    if command -v npm &> /dev/null && [ -f "package.json" ]; then
        npm test || echo "   ⚠️ npm tests need attention"
    fi
    if [ -f "Makefile" ] && grep -q "^test:" Makefile; then
        make test || echo "   ⚠️ make tests need attention"
    fi
done

# 8. Tag phase completion in TARGET repo
TAG_NAME="phase${PHASE}-complete-v1.0"
echo "🏷️ Creating phase completion tag: $TAG_NAME"
git tag -a "$TAG_NAME" -m "Phase $PHASE Complete - All waves integrated, Architect Approved"

# 9. Push branch and tags
echo "📤 Pushing phase integration branch and tags..."
git push -u origin "$PHASE_BRANCH"
git push origin --tags

echo "✅ Phase integration complete in TARGET repository"
echo "📍 Location: $PHASE_INTEGRATION_DIR"
echo "🌿 Branch: $PHASE_BRANCH"
echo "🏷️ Tag: $TAG_NAME"
```

**VIOLATION OF R104/R282 = AUTOMATIC FAILURE!**
- Never create integration branches in software-factory-template or planning repo
- Always clone target repository first
- Phase integration happens in target repo only
- Directory MUST be `efforts/phaseX/integration-workspace/`
- Repository verification is MANDATORY

## Phase Documentation

```bash
# Generate phase completion report
ASSESSMENT_REPORT=$(yq '.phase_assessment.report_file' orchestrator-state.yaml)
if [ -z "$ASSESSMENT_REPORT" ] || [ ! -f "$ASSESSMENT_REPORT" ]; then
    echo "❌ CRITICAL: No phase assessment report found in state!"
    echo "❌ This violates R257 - cannot complete phase without report"
    exit 1
fi

cat > "PHASE-${PHASE}-COMPLETION-REPORT.md" << EOF
# Phase $PHASE Completion Report

## Summary
- **Phase**: $PHASE
- **Waves Completed**: $WAVE_COUNT
- **Efforts Delivered**: $EFFORT_COUNT
- **Lines of Code**: $TOTAL_LINES
- **Test Coverage**: $COVERAGE%
- **Architect Assessment Report**: $ASSESSMENT_REPORT
- **Architect Approval**: $(date -u +%Y-%m-%dT%H:%M:%SZ)

## Achievements
$(list_phase_achievements)

## Delivered Features
$(list_delivered_features)

## Architecture Decisions
$(list_architecture_decisions)

## Metrics
- Code Review First-Try Success: $REVIEW_SUCCESS%
- Split Compliance Rate: $SPLIT_COMPLIANCE%
- Integration Success Rate: $INTEGRATION_SUCCESS%
- Average Effort Size: $AVG_EFFORT_SIZE lines

## Lessons Learned
$(compile_lessons_learned)

## Next Steps
$(determine_next_steps)
EOF

git add "PHASE-${PHASE}-COMPLETION-REPORT.md"
git commit -m "docs: Phase $PHASE completion report"
git push
```

## State File Updates

```bash
# Record phase completion
yq -i ".phases_completed[] = $PHASE" orchestrator-state.yaml
yq -i ".phase_metrics.phase$PHASE.completed_at = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" orchestrator-state.yaml
yq -i ".phase_metrics.phase$PHASE.waves = $WAVE_COUNT" orchestrator-state.yaml
yq -i ".phase_metrics.phase$PHASE.efforts = $EFFORT_COUNT" orchestrator-state.yaml
yq -i ".phase_metrics.phase$PHASE.integration_branch = \"$PHASE_BRANCH\"" orchestrator-state.yaml
```

## Multi-Phase Decision Logic

```bash
# Check if more phases exist
TOTAL_PHASES=$(get_total_phases_from_plan)
CURRENT_PHASE=$(yq '.current_phase' orchestrator-state.yaml)

if [ "$CURRENT_PHASE" -lt "$TOTAL_PHASES" ]; then
    # More phases to complete
    echo "Phase $CURRENT_PHASE complete. Preparing for Phase $((CURRENT_PHASE + 1))"
    
    # Update for next phase
    yq -i ".current_phase = $((CURRENT_PHASE + 1))" orchestrator-state.yaml
    yq -i ".current_wave = 1" orchestrator-state.yaml
    yq -i ".current_state = \"INIT\"" orchestrator-state.yaml
    
    # Transition to start next phase
    transition_to "INIT"  # Start next phase planning
else
    # All phases complete
    echo "All $TOTAL_PHASES phases completed successfully!"
    
    # Create final project summary
    create_project_completion_summary
    
    # Transition to terminal success
    transition_to "SUCCESS"  # Project complete!
fi
```

## Success Criteria

Phase completion tasks succeed when:
- [ ] Phase integration branch created
- [ ] All waves merged successfully
- [ ] Phase tagged appropriately
- [ ] Completion report generated
- [ ] Metrics updated in state file
- [ ] Next action determined

## State Transitions

From PHASE_COMPLETE:
- **Single phase project** → SUCCESS (project done)
- **Multi-phase, more phases** → INIT (next phase)
- **Multi-phase, last phase** → SUCCESS (project done)
- **Integration failure** → ERROR_RECOVERY

## Phase Completion Checklist

Before transitioning from PHASE_COMPLETE:
1. ✅ Phase integration branch created and pushed
2. ✅ All wave branches merged
3. ✅ Phase completion report written
4. ✅ Metrics documented
5. ✅ State file updated
6. ✅ Lessons learned captured
7. ✅ Next steps determined

## Common Tasks

1. **Create Pull Request**: For phase branch to main
2. **Update Documentation**: Phase-level docs
3. **Notify Stakeholders**: Phase complete
4. **Archive Branches**: Clean up wave branches
5. **Update Roadmap**: Mark phase done

## Required Actions

1. Create phase integration branch
2. Merge all wave integrations
3. Generate completion report
4. Update state file metrics
5. Determine if more phases exist
6. Transition to SUCCESS or INIT

## Grading Impact

- Phase integration successful: +20 points
- Complete documentation: +15 points
- Proper metrics capture: +10 points
- Clean state transitions: +10 points
- Skipping phase finalization: -50 points

## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-orchestrating
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**
