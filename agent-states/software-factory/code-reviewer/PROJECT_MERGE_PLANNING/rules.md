# Code Reviewer - PROJECT_MERGE_PLANNING State Rules

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
## PRIMARY PURPOSE
Create a detailed PROJECT MERGE PLAN document for integrating all phases into the final project.
DO NOT execute any merges - only plan them!

### 🔴🔴🔴 CRITICAL RULE: USE ORIGINAL BRANCHES ONLY! 🔴🔴🔴

**MERGE FROM ORIGINAL EFFORT BRANCHES ONLY!**
- ✅ CORRECT: phase1/wave1/effort1-foundation (original effort)
- ✅ CORRECT: phase2/wave1/effort2-core-feature (original effort)
- ❌ WRONG: phase1/integration-20250827 (phase integration)
- ❌ WRONG: project/integration-20250827 (project integration)

Phase integration branches are intermediate artifacts, not sources!

## State Context
You are creating a project-level merge plan to integrate ALL phases into the final deliverable. The orchestrator has set up the project integration infrastructure. Your role is to plan the merging of ALL phase efforts.

## PROJECT MERGE PLAN REQUIREMENTS

### 1. Inventory All Branches to Merge
```bash
#!/bin/bash
# Find all branches for project integration

echo "📊 Analyzing branches for Project integration..."

# Inventory all phase efforts
for phase in $(seq 1 5); do  # Adjust based on actual phase count
    echo "=== PHASE $phase EFFORTS ==="
    for wave in $(seq 1 4); do  # Adjust based on actual wave count
        echo "  Wave $wave:"
        for effort_dir in /efforts/phase${phase}/wave${wave}/*/; do
            [[ "$effort_dir" == *"integration-workspace"* ]] && continue

            effort=$(basename "$effort_dir")
            cd "$effort_dir"

            current_branch=$(git branch --show-current)

            # Include splits, exclude "too large" originals
            if [[ "$current_branch" == *"-split"* ]]; then
                echo "    ✅ $current_branch (split)"
            elif git log --oneline | grep -q "too large"; then
                echo "    ❌ $current_branch (too large, use splits)"
            else
                echo "    ✅ $current_branch"
            fi
        done
    done
done
```

### 2. Determine Project-Level Merge Order
```markdown
## Project Merge Order Analysis

### Merge Groups (IN ORDER):

#### Group 1: Phase 1 (Foundation)
All Phase 1 efforts across all waves
- Foundation layer must be established first

#### Group 2: Phase 2 (Core Features)
All Phase 2 efforts across all waves
- Depends on Phase 1 foundation

#### Group 3: Phase 3 (Advanced Features)
All Phase 3 efforts across all waves
- Depends on Phases 1-2

#### Group 4: Phase 4 (Integration & Polish)
All Phase 4 efforts across all waves
- Depends on Phases 1-3

#### Group 5: Phase 5 (Final Polish)
All Phase 5 efforts across all waves
- Depends on all previous phases

### Dependency Analysis:
- Phases are merged in sequence (1→2→3→4→5)
- Within a phase, waves are merged in order
- Within a wave, splits must maintain order (split1→split2→split3)
```

### 3. Create PROJECT-MERGE-PLAN.md

**Location:** `${INTEGRATE_WAVE_EFFORTS_DIR}/PROJECT-MERGE-PLAN.md`

```markdown
# Project Integration Merge Plan

**Generated:** $(date -u +%Y-%m-%dT%H:%M:%SZ)
**Code Reviewer:** code-reviewer
**State:** PROJECT_MERGE_PLANNING
**Type:** Complete Project Integration

## Target Integration Branch
- **Branch Name:** project/integration-${timestamp}
- **Base:** main at ${commit}
- **Purpose:** Complete project integration
- **Location:** /efforts/project-integration-workspace

## 🎬 Demo Execution Plan (R330/R291 Compliance)

### Demo Requirements Overview
Per R330 and R291, ALL integrations MUST demonstrate working functionality.
Project-level demos validate that ALL phases work together as a complete, production-ready system.

### Demo Execution Sequence
1. **After Each Phase Group Merge:**
   - Run phase-specific demo script if exists
   - Capture output in `demo-results/phase-X-demo.log`
   - Continue even if individual demo fails (document for review)

2. **After All Merges Complete:**
   - Run integrated project demo: `./project-demo.sh`
   - Verify all project features work together
   - Verify end-to-end user scenarios
   - Capture evidence in `demo-results/project-integration-demo.log`

3. **Demo Validation Gates (R291):**
   - ✅ BUILD GATE: Code must compile
   - ✅ TEST GATE: All tests must pass
   - ✅ DEMO GATE: Demo scripts must execute
   - ✅ ARTIFACT GATE: Build outputs must exist
   - ✅ E2E GATE: End-to-end scenarios work

### Demo Files Expected
Based on phase integration demos (R330/R291), these demos should exist:
- [ ] phase1-integration/demo-features.sh (from phase 1 integration)
- [ ] phase2-integration/demo-features.sh (from phase 2 integration)
- [ ] phase3-integration/demo-features.sh (from phase 3 integration)
- [ ] phase4-integration/demo-features.sh (from phase 4 integration)
- [ ] phase5-integration/demo-features.sh (from phase 5 integration)
- [ ] PROJECT-DEMO.md (project integration demo documentation)
- [ ] project-demo.sh (integrated project demo script)
- [ ] E2E-SCENARIOS.md (end-to-end scenario documentation)
- [ ] e2e-test-suite.sh (end-to-end automated tests)

### Demo Failure Protocol
If ANY demo fails during integration:
1. Document failure in PROJECT-INTEGRATE_WAVE_EFFORTS-REPORT.md
2. Mark Demo Status: FAILED
3. This will trigger ERROR_RECOVERY per R291
4. Fixes must be made in effort branches (R300)

## Project Scope Summary
- **Total Phases:** 5
- **Total Waves:** 20 (across all phases)
- **Total Efforts:** 60+ (including splits)
- **Target:** Production-ready deliverable

## Branches to Merge (IN STRICT ORDER)

### Phase 1 Efforts (Foundation)

[List all Phase 1 efforts following the same pattern as PHASE_MERGE_PLANNING]

### Phase 2 Efforts (Core Features)

[Continue for all phases...]

## Excluded Branches (DO NOT MERGE)
These branches are superseded or intermediate:
- phase*/wave*/integration-* (intermediate wave integrations)
- phase*/integration-* (intermediate phase integrations)
- All "too large" original branches (use splits instead)

## Merge Strategy
1. **Phase-by-Phase:** Complete each phase before moving to next
2. **Waves in Sequence:** Complete all waves within a phase sequentially
3. **Splits in Order:** Maintain split order within efforts
4. **Testing:** Project-level tests after all merges
5. **E2E Validation:** Full end-to-end scenario testing

## Expected Integration Challenges
1. **Large Scope:** 60+ branches is significant complexity
2. **Cross-Phase Dependencies:** Features spanning multiple phases
3. **API Evolution:** APIs evolved across phases, ensure compatibility
4. **Performance:** Combined system must meet performance requirements
5. **Security:** All security requirements must be validated

## Project-Level Validation
```bash
# After ALL merges complete:

# 1. Run complete project test suite
make test-all

# 2. Check combined size
$PROJECT_ROOT/tools/line-counter.sh -c $(git branch --show-current)

# 3. Run integrated project demo (R291 requirement)
echo "🎬 Running integrated project demo..."
if [ -f "./project-demo.sh" ]; then
    bash ./project-demo.sh > demo-results/project-integration-demo.log 2>&1
    DEMO_STATUS=$?
    echo "Project demo status: $DEMO_STATUS"
else
    echo "🔴 CRITICAL: No project-demo.sh found - R291 VIOLATION!"
    exit 291
fi

# 4. Verify all phase demos executed (R291)
for phase in $(seq 1 5); do
    echo "Checking phase $phase demo..."
    if [ -f "phase${phase}-integration/demo-features.sh" ]; then
        bash "phase${phase}-integration/demo-features.sh" > "demo-results/phase${phase}-demo.log" 2>&1
        echo "Phase $phase demo exit code: $?"
    else
        echo "⚠️ WARNING: Phase $phase demo missing"
    fi
done

# 5. Run end-to-end test suite (R291 requirement)
echo "🎬 Running end-to-end test scenarios..."
if [ -f "./e2e-test-suite.sh" ]; then
    bash ./e2e-test-suite.sh > demo-results/e2e-tests.log 2>&1
    E2E_STATUS=$?
    echo "E2E test status: $E2E_STATUS"
    if [ $E2E_STATUS -ne 0 ]; then
        echo "🔴 E2E tests failed - R291 GATE FAILURE"
        exit 291
    fi
else
    echo "🔴 CRITICAL: No e2e-test-suite.sh found - R291 VIOLATION!"
    exit 291
fi

# 6. Performance benchmarks
make benchmark-all

# 7. Security scan
make security-scan

# 8. Production readiness checks
make production-ready

# 9. Capture demo evidence for review (R291)
ls -la demo-results/
echo "✅ All demos executed and captured"
```

## Risk Mitigation
- **Very High Risk:** Large number of branches, complex dependencies
- **Mitigation 1:** Test after each phase group
- **Mitigation 2:** Save state between phase groups
- **Mitigation 3:** Comprehensive logging in work-log.md
- **Mitigation 4:** E2E testing at multiple levels

## Integration Agent Instructions
1. This is a PROJECT-level integration (largest scope)
2. Execute merges in EXACT order specified
3. Test after each PHASE GROUP (not just each merge)
4. Document all conflicts with detailed resolution notes
5. Verify ALL demos pass (R291 is NON-NEGOTIABLE)
6. Create comprehensive PROJECT-INTEGRATE_WAVE_EFFORTS-REPORT.md
7. Include production readiness assessment

## Success Criteria
- All branches merged successfully
- No test failures after integration
- All demos pass (R291 compliance)
- E2E scenarios pass
- Performance benchmarks pass
- Security scan passes
- Production readiness validated
- Ready for final release
```

## Validation Before Completion

```bash
#!/bin/bash
# Validate project merge plan completeness

validate_project_merge_plan() {
    local plan_file="$1"

    echo "🔍 Validating Project Merge Plan..."

    # Check phase coverage
    for phase in $(seq 1 5); do
        if ! grep -q "Phase $phase Efforts" "$plan_file"; then
            echo "❌ Missing Phase $phase in merge plan"
            return 1
        fi
    done

    # Check demo section exists (R291/R330 requirement)
    if ! grep -q "Demo Execution Plan" "$plan_file"; then
        echo "❌ CRITICAL: Missing Demo Execution Plan section (R291/R330 VIOLATION)"
        return 1
    fi

    # Check demo validation gates
    if ! grep -q "Demo Validation Gates" "$plan_file"; then
        echo "❌ CRITICAL: Missing Demo Validation Gates (R291 VIOLATION)"
        return 1
    fi

    # Check E2E requirements
    if ! grep -q "E2E" "$plan_file"; then
        echo "❌ CRITICAL: Missing E2E requirements (R291 project requirement)"
        return 1
    fi

    # Verify no integration branches as sources
    if grep -E "phase[0-9]/integration\|project/integration" "$plan_file" | grep "git merge"; then
        echo "❌ ERROR: Integration branches used as sources!"
        return 1
    fi

    # Count total merges planned
    merge_count=$(grep -c "git merge origin/" "$plan_file")
    echo "📊 Total branches to merge: $merge_count"

    if [[ $merge_count -lt 20 ]]; then
        echo "⚠️ Warning: Fewer merges than expected for a full project"
    fi

    echo "✅ Project merge plan validation passed"
    return 0
}

# Run validation
validate_project_merge_plan "PROJECT-MERGE-PLAN.md"
```

## State Transitions

From PROJECT_MERGE_PLANNING state:
- **PLAN_COMPLETE** → Return to orchestrator
- **VALIDATION_FAILED** → Fix and re-validate

## Critical Success Criteria

1. ✅ PROJECT-MERGE-PLAN.md created in project integration directory
2. ✅ All phase efforts catalogued (excluding "too large" originals)
3. ✅ Phase-by-phase merge order established
4. ✅ NO phase/wave integration branches used as sources
5. ✅ **Demo Execution Plan included (R291/R330 MANDATORY)**
6. ✅ **Demo Validation Gates specified (R291 MANDATORY)**
7. ✅ **E2E testing requirements included (R291 project requirement)**
8. ✅ Project-level validation steps included
9. ✅ Clear instructions for Integration Agent

## Common Mistakes to Avoid

1. **Merging from phase integration branches**
   - ❌ WRONG: Use phase1/integration as source
   - ✅ RIGHT: Use original effort branches only

2. **Missing demo planning**
   - ❌ WRONG: No demo section (R291/R330 VIOLATION)
   - ✅ RIGHT: Comprehensive demo plan with all gates

3. **Skipping E2E validation**
   - ❌ WRONG: Only unit/integration tests
   - ✅ RIGHT: Full E2E scenario validation

4. **Wrong phase order**
   - ❌ WRONG: Random phase ordering
   - ✅ RIGHT: Sequential phase progression

5. **Executing merges**
   - ❌ WRONG: Running git merge commands
   - ✅ RIGHT: Only documenting merge plan

---
### ⚠️⚠️⚠️ RULE R261 - Code Reviewer Merge Plan No Execution
**Source:** rule-library/R261-integration-planning-requirements.md
**Criticality:** WARNING - Violation = Role confusion

Code Reviewer creates merge plans ONLY. NEVER executes merges. That's the Integration Agent's job.
---

---
### 🔴🔴🔴 RULE R262 - No Integration Branches as Sources
**Source:** rule-library/R262-merge-operation-protocols.md
**Criticality:** SUPREME - Violation = Recursive integration chaos

CRITICAL: Only original effort branches in merge plans. Integration branches are TARGETS not SOURCES.
---

---
### 🔴🔴🔴 RULE R291 - Integration Demo Requirement
**Source:** rule-library/R291-integration-demo-requirement.md
**Criticality:** BLOCKING - Violation = -50% to -75%

**PROJECT-LEVEL DEMOS ARE MANDATORY** and include:
- Phase-level demo validation
- Project-level integrated demo
- End-to-end scenario testing
- Production readiness validation

**See R291 for complete requirements**
---

---
### 🔴🔴🔴 RULE R330 - Demo Planning Requirements
**Source:** rule-library/R330-demo-planning-requirements.md
**Criticality:** BLOCKING - Violation = -25% to -50%

**EVERY merge plan MUST include explicit demo requirements:**
- Demo objectives
- Demo scenarios
- Demo validation gates
- Expected demo files

**See R330 for complete requirements**
---

### 🔴🔴🔴 RULE R340: Report Merge Plan Location to Orchestrator 🔴🔴🔴

**File**: `$CLAUDE_PROJECT_DIR/rule-library/R340-planning-file-metadata-tracking.md`
**Criticality**: BLOCKING - All planning files must be tracked

**AFTER CREATING PROJECT-MERGE-PLAN.md, YOU MUST REPORT:**

```markdown
## 📋 PLANNING FILE CREATED

**Type**: merge_plan (project)
**Path**: /efforts/project-integration-workspace/PROJECT-MERGE-PLAN.md
**Target Branch**: project/integration
**Phases Count**: {N}
**Created At**: {ISO-8601-timestamp}

ORCHESTRATOR: Please update planning_files.merge_plans.project in state file per R340
```

**R340 AUTOMATION: Update Orchestrator State Directly**

After creating the project merge plan, Code Reviewer MUST update orchestrator-state-v3.json with R340 metadata:

```bash
# R340: Record project merge plan location in orchestrator state
PROJECT_ID="project"
TIMESTAMP=$(date -u +%Y-%m-%dT%H:%M:%SZ)

# Get absolute path to merge plan
MERGE_PLAN_ABS=$(realpath "$MERGE_PLAN")

# Update orchestrator state with R340 metadata
jq --arg project_id "$PROJECT_ID" \
   --arg path "$MERGE_PLAN_ABS" \
   --arg timestamp "$TIMESTAMP" \
   --arg creator "code-reviewer" \
   '
   # Initialize planning_repo_files if not exists
   .planning_repo_files //= {} |
   .planning_repo_files.merge_plans //= {} |
   .planning_repo_files.merge_plans.project //= {} |

   # Record project merge plan metadata
   .planning_repo_files.merge_plans.project[$project_id] = {
     "file_path": $path,
     "created_at": $timestamp,
     "created_by": $creator,
     "status": "active",
     "deprecated": false
   }
   ' orchestrator-state-v3.json > orchestrator-state-v3.json.tmp

mv orchestrator-state-v3.json.tmp orchestrator-state-v3.json

echo "✅ R340: Project merge plan tracked in orchestrator-state-v3.json"
echo "   Project ID: $PROJECT_ID"
echo "   Path: $MERGE_PLAN_ABS"

# Commit the state update
git add orchestrator-state-v3.json
git commit -m "feat: Track project merge plan per R340

R340 metadata recorded in planning_repo_files.merge_plans.project

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>"

git push
```

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
3. **EXACTLY AS SHOWN**: Use exact format - no variations
4. **ALWAYS REQUIRED**: Every single state must output this flag
5. **GREPPABLE**: Must be on its own line for automation parsing

### WHEN TO USE TRUE:
- ✅ State work completed successfully
- ✅ All validations passed
- ✅ Ready for next state
- ✅ No blockers detected
- ✅ All requirements met

### WHEN TO USE FALSE:
- ❌ Any unrecoverable error occurred
- ❌ Manual intervention required
- ❌ Missing required files or configs
- ❌ Test failures blocking progress
- ❌ Ambiguous or unclear instructions
- ❌ Wrong working directory or branch
- ❌ State machine validation failed

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**
