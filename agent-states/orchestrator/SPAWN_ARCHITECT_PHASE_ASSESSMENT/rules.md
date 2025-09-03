# Orchestrator - SPAWN_ARCHITECT_PHASE_ASSESSMENT State Rules

## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED SPAWN_ARCHITECT_PHASE_ASSESSMENT STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_SPAWN_ARCHITECT_PHASE_ASSESSMENT
echo "$(date +%s) - Rules read and acknowledged for SPAWN_ARCHITECT_PHASE_ASSESSMENT" > .state_rules_read_orchestrator_SPAWN_ARCHITECT_PHASE_ASSESSMENT
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY SPAWN_ARCHITECT_PHASE_ASSESSMENT WORK UNTIL RULES ARE READ:
- ❌ Start spawn architect agent
- ❌ Start request phase assessment
- ❌ Start evaluate phase completion
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
   ❌ WRONG: "I acknowledge all SPAWN_ARCHITECT_PHASE_ASSESSMENT rules"
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

### ✅ CORRECT PATTERN FOR SPAWN_ARCHITECT_PHASE_ASSESSMENT:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/[first-rule-file].md
2. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/[second-rule-file].md  
4. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. "Ready to execute SPAWN_ARCHITECT_PHASE_ASSESSMENT work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY SPAWN_ARCHITECT_PHASE_ASSESSMENT work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ You have stated readiness to execute SPAWN_ARCHITECT_PHASE_ASSESSMENT work
1. **READ** every rule file listed in PRIMARY DIRECTIVES below
2. **ACKNOWLEDGE** each rule individually with number and description
3. **ONLY THEN** proceed with SPAWN_ARCHITECT_PHASE_ASSESSMENT work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY SPAWN_ARCHITECT_PHASE_ASSESSMENT work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

**THE SYSTEM IS MONITORING YOUR READ TOOL CALLS!**

## 📋 PRIMARY DIRECTIVES FOR SPAWN_ARCHITECT_PHASE_ASSESSMENT STATE

### 🔴🔴🔴 R301 - Integration Branch Current Tracking (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R301-integration-branch-current-tracking.md`
**Criticality**: SUPREME LAW - Only ONE current integration allowed
**Summary**: MUST use current_integration.branch, NEVER deprecated branches

### 🚨🚨🚨 R257 - Mandatory Phase Assessment Report
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R257-mandatory-phase-assessment-report.md`
**Criticality**: BLOCKING - Phase cannot complete without report
**Summary**: Architect MUST create assessment report file

### 🚨🚨🚨 R285 - Mandatory Phase Integration Before Assessment
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R285-mandatory-phase-integration-before-assessment.md`
**Criticality**: BLOCKING - Integration must precede assessment
**Summary**: Assess integrated work, not individual efforts

### 🔴🔴🔴 R233 - All States Require Immediate Action
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R233-all-states-immediate-action.md`
**Criticality**: CRITICAL - States are verbs
**Summary**: SPAWN_ARCHITECT_PHASE_ASSESSMENT means SPAWN NOW

## 🚨 SPAWN_ARCHITECT_PHASE_ASSESSMENT IS A VERB - SPAWN IMMEDIATELY! 🚨

### IMMEDIATE ACTIONS UPON ENTERING SPAWN_ARCHITECT_PHASE_ASSESSMENT

**THE MOMENT YOU ENTER THIS STATE, YOU MUST:**
1. Spawn Architect for PHASE-LEVEL assessment NOW
2. Provide all wave integration branches for the phase
3. Include phase completion metrics
4. Request comprehensive phase validation

**FORBIDDEN - AUTOMATIC FAILURE:**
- ❌ "STATE TRANSITION COMPLETE: Now in SPAWN_ARCHITECT_PHASE_ASSESSMENT" [stops]
- ❌ "Successfully entered SPAWN_ARCHITECT_PHASE_ASSESSMENT state" [waits]
- ❌ "Ready to spawn architect" [pauses]
- ❌ "I'm in SPAWN_ARCHITECT_PHASE_ASSESSMENT state" [does nothing]

**REQUIRED - IMMEDIATE ACTION:**
- ✅ "Entering SPAWN_ARCHITECT_PHASE_ASSESSMENT, spawning architect NOW..."
- ✅ "SPAWN ARCHITECT for complete Phase $PHASE assessment..."
- ✅ "Providing all wave integrations for phase validation..."

## State Context

This state is entered in two scenarios:
1. **Initial Assessment**: After the last wave of a phase completes successfully
2. **Reassessment**: After PHASE_INTEGRATION creates integration branch with fixes (R259)

The architect must assess the ENTIRE PHASE before allowing transition to SUCCESS.

## 🔴🔴🔴 CRITICAL: Phase Completion Gate 🔴🔴🔴

**This state is the MANDATORY gateway to phase completion:**
- NO phase can be marked complete without architect assessment
- NO transition to SUCCESS without passing this gate
- NO shortcuts or bypassing allowed

## Primary Purpose

The SPAWN_ARCHITECT_PHASE_ASSESSMENT state is for:
1. Spawning architect to assess complete phase work
2. Providing all phase deliverables for review
3. Ensuring phase-level architectural integrity
4. Validating readiness for phase completion

## Phase Assessment Scope

The architect must review:
- **All Wave Integrations**: Every wave branch merged correctly
- **Phase Architecture**: Overall design coherence
- **Feature Completeness**: All planned features implemented
- **API Stability**: APIs ready for external use
- **Test Coverage**: Phase-level test adequacy
- **Documentation**: Phase documentation complete
- **Performance**: Phase-level performance validation
- **Security**: Security requirements met

## Spawning the Architect

```bash
# R301 MANDATORY: Get ONLY the current phase integration branch
CURRENT_BRANCH=$(yq '.current_phase_integration | select(.phase == env(PHASE)).branch' orchestrator-state.yaml)
CURRENT_STATUS=$(yq '.current_phase_integration | select(.phase == env(PHASE)).status' orchestrator-state.yaml)
INTEGRATION_TYPE=$(yq '.current_phase_integration | select(.phase == env(PHASE)).type' orchestrator-state.yaml)

# Validate current integration is active (R301)
if [ "$CURRENT_STATUS" != "active" ]; then
    echo "🔴🔴🔴 FATAL: Current integration is not active!"
    echo "  Phase: $PHASE"
    echo "  Branch: $CURRENT_BRANCH"
    echo "  Status: $CURRENT_STATUS"
    exit 1
fi

# Determine assessment type based on integration type
if [ "$INTEGRATION_TYPE" = "post_fixes" ]; then
    ASSESSMENT_TYPE="REASSESSMENT after fixes"
    PHASE_BRANCH="$CURRENT_BRANCH"
    ORIGINAL_REPORT="phase-assessments/phase${PHASE}/PHASE-${PHASE}-ASSESSMENT-REPORT.md"
    IS_REASSESSMENT="true"
else
    ASSESSMENT_TYPE="INITIAL ASSESSMENT"
    PHASE_BRANCH="$CURRENT_BRANCH"
    ORIGINAL_REPORT=""
    IS_REASSESSMENT="false"
fi

echo "✅ Using current integration per R301: $PHASE_BRANCH"

# Prepare phase assessment request
PHASE_CONTEXT="Complete Phase $PHASE ${ASSESSMENT_TYPE}"

# Gather all wave branches
WAVE_BRANCHES=$(list_all_wave_integration_branches "$PHASE")

# Key information for architect
ASSESSMENT_PROMPT="Perform COMPLETE PHASE $PHASE ${ASSESSMENT_TYPE}.

🚨🚨🚨 MANDATORY REQUIREMENT [R257]: You MUST create a phase assessment report file! 🚨🚨🚨

1. CREATE REPORT FILE: phase-assessments/phase${PHASE}/PHASE-${PHASE}-ASSESSMENT-REPORT.md
2. Include ALL mandatory sections per R257
3. Document your DECISION in the report
4. Commit and push the report
5. ONLY THEN signal assessment complete

$([ -n "$ORIGINAL_REPORT" ] && echo "
📋 THIS IS A REASSESSMENT:
- Original report with issues: $ORIGINAL_REPORT  
- Fixes have been integrated in: $PHASE_BRANCH
- Verify all Priority 1 issues are addressed
- Check if original concerns are resolved
")

Integration Branches:
$WAVE_BRANCHES

Phase Branch: $PHASE_BRANCH

Assess:
1. Feature completeness against phase plan
2. Architectural integrity across all waves
3. API readiness and stability
4. Test coverage adequacy at phase level
5. Documentation completeness
6. Performance and security validation
7. Readiness for production use

Provide decision IN YOUR REPORT FILE: 
- PHASE_COMPLETE (ready for SUCCESS)
- NEEDS_WORK (specific fixes required)
- PHASE_FAILED (major issues, cannot complete)

This assessment gates phase completion - be thorough.
Your assessment is INVALID without the report file per R257."

# Spawn architect reviewer
Task: subagent_type="architect-reviewer" \
      prompt="$ASSESSMENT_PROMPT" \
      description="Complete Phase $PHASE Assessment (Gates SUCCESS)"
```

## State File Updates

```bash
# Record phase assessment request
yq -i ".phase_assessment.phase = $PHASE" orchestrator-state.yaml
yq -i ".phase_assessment.requested_at = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" orchestrator-state.yaml
yq -i ".phase_assessment.expected_report = \"phase-assessments/phase${PHASE}/PHASE-${PHASE}-ASSESSMENT-REPORT.md\"" orchestrator-state.yaml
yq -i ".phase_assessment.type = \"COMPLETE_PHASE\"" orchestrator-state.yaml
yq -i ".phase_assessment.phase_branch = \"$PHASE_BRANCH\"" orchestrator-state.yaml
yq -i ".phase_assessment.wave_count = $WAVE_COUNT" orchestrator-state.yaml
yq -i ".phase_assessment.status = \"PENDING\"" orchestrator-state.yaml

# Immediate transition to waiting
transition_to "WAITING_FOR_PHASE_ASSESSMENT"
```

## Success Criteria for Phase Assessment

The phase assessment passes when:
- [ ] All phase features implemented
- [ ] Architectural consistency verified
- [ ] APIs stable and documented
- [ ] Test coverage meets requirements
- [ ] Performance benchmarks met
- [ ] Security requirements satisfied
- [ ] Documentation complete
- [ ] No blocking issues remain

## State Transitions

From SPAWN_ARCHITECT_PHASE_ASSESSMENT:
- **Always** → WAITING_FOR_PHASE_ASSESSMENT (immediate after spawn)

## Multi-Phase Considerations

If this is NOT the final phase:
- Architect may recommend starting next phase planning
- PHASE_COMPLETE state will handle transition to next phase
- Document lessons learned for next phase

## Required Actions

1. Update state file with assessment request
2. Spawn architect with comprehensive context
3. Provide all integration branches
4. Include phase metrics and achievements
5. Transition to WAITING_FOR_PHASE_ASSESSMENT

## Grading Impact

- Spawning architect promptly: +10 points
- Providing complete context: +10 points
- Including all branches: +10 points
- Proper state file updates: +10 points
- Missing phase assessment before SUCCESS: -100 points (CRITICAL FAILURE)
