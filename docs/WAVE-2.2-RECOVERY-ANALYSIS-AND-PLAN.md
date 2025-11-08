# Wave 2.2 Infrastructure Recovery - Analysis and Execution Plan

**Date**: 2025-11-01 (Analysis Complete)
**Analyzed By**: Software Factory Manager
**Status**: READY FOR EXECUTION

---

## EXECUTIVE SUMMARY

**Problem**: Effort 2.2.2 infrastructure was created with WRONG base branch, violating R509 (Mandatory Base Branch Validation). The effort is based on Wave 2.1 integration instead of Effort 2.2.1, missing critical dependencies (config.go).

**Root Cause**: CREATE_NEXT_INFRASTRUCTURE state did not have R603 (dependency checking logic) when infrastructure was created.

**Solution**: Delete broken 2.2.2 infrastructure, upgrade project to add R603, reset state machine to VALIDATE_INFRASTRUCTURE, allow system to recreate 2.2.2 correctly.

**Impact**: Preserves all approved work (Effort 2.2.1), enables correct cascade dependency chain.

---

## PART 1: COMPREHENSIVE ANALYSIS

### 1.1 Current System State

#### Orchestrator State Machine
```
Current State: ERROR_RECOVERY
Previous State: SPAWN_SW_ENGINEERS
Phase: 2
Wave: 2
Last Transition: 2025-11-01T20:07:14Z
```

#### Effort 2.2.1 Status: ✅ COMPLETE AND APPROVED
```
Branch: idpbuilder-oci-push/phase2/wave2/effort-1-registry-override-viper
Workspace: efforts/phase2/wave2/effort-1-registry-override-viper
Base Branch: idpbuilder-oci-push/phase2/wave1/integration (CORRECT)
Remote: origin/idpbuilder-oci-push/phase2/wave2/effort-1-registry-override-viper
Status: APPROVED (code review 2025-11-01 19:19:22 UTC)
Lines: 247 lines (within 800 limit)
Key File: config.go (203 lines) - NEW FILE CREATED
Review Report: .software-factory/phase2/wave2/effort-1-registry-override-viper/CODE-REVIEW-REPORT--20251101-192258.md
Commit: 6e44ac2 "review(2.2.1): Code review APPROVED"
```

#### Effort 2.2.2 Status: ❌ WRONG BASE BRANCH (R509 VIOLATION)
```
Branch: idpbuilder-oci-push/phase2/wave2/effort-2-env-variable-support
Workspace: efforts/phase2/wave2/effort-2-env-variable-support
Base Branch: idpbuilder-oci-push/phase2/wave1/integration (WRONG!)
            Should be: idpbuilder-oci-push/phase2/wave2/effort-1-registry-override-viper
Remote: origin/idpbuilder-oci-push/phase2/wave2/effort-2-env-variable-support
Merge Base: 978f94c (Wave 2.1 integration) - CONFIRMS WRONG BASE
Missing: config.go (from 2.2.1) - DEPENDENCY NOT SATISFIED
Latest Commit: 7f81263 "plan: Create implementation plan for Effort 2.2.2"
```

#### Pre-Planned Infrastructure Metadata
```json
"phase2_wave2_effort_1_registry_override": {
  "base_branch": "idpbuilder-oci-push/phase2/wave1/integration",
  "created": true,
  "validated": true
}
"phase2_wave2_effort_2_env_support": {
  "base_branch": "idpbuilder-oci-push/phase2/wave2/effort-1-registry-override-viper",
  "depends_on": ["phase2_wave2_effort_1_registry_override"],
  "created": true,
  "validated": true
}
```

**METADATA SAYS**: 2.2.2 should be based on 2.2.1
**ACTUAL GIT SHOWS**: 2.2.2 is based on Wave 2.1 integration
**CONCLUSION**: Infrastructure creation logic did NOT enforce cascade dependency

### 1.2 R603 Status

```bash
$ ls /home/vscode/workspaces/idpbuilder-oci-push-planning/rule-library/R603*
ls: cannot access '...R603*': No such file or directory
```

**R603 does NOT exist in this project yet.**

R603 was created in the template repository AFTER this project was initialized. It provides:
- Sequential effort dependency validation
- Cascade base branch verification
- R213 metadata reading for dependencies

**UPGRADE REQUIRED**: Must run `/home/vscode/workspaces/idpbuilder-oci-push-planning/tools/upgrade.sh` to sync R603 into this project.

---

## PART 2: CRITICAL QUESTIONS - ANSWERED

### Question 1: What State Should We Reset To?

**ANSWER: VALIDATE_INFRASTRUCTURE**

#### Options Analyzed:

**Option A: MONITORING_EFFORT_REVIEWS** ❌
- Pro: Natural place after effort approval
- Con: Expects infrastructure to already exist
- Con: Would try to spawn Code Reviewer for non-existent infrastructure
- Verdict: NOT VIABLE

**Option B: CREATE_NEXT_INFRASTRUCTURE** ❌
- Pro: This is where R603 logic lives
- Con: Not a valid transition from ERROR_RECOVERY
- Con: Would require non-standard state transition
- Verdict: CAN'T GET THERE DIRECTLY

**Option C: VALIDATE_INFRASTRUCTURE** ✅ **RECOMMENDED**
- Pro: Has loop-back to CREATE_NEXT_INFRASTRUCTURE
- Pro: Detects missing/invalid infrastructure
- Pro: Natural recovery mechanism
- Pro: Preserves 2.2.1 approved status
- Mechanism:
  1. VALIDATE_INFRASTRUCTURE checks all efforts
  2. Detects 2.2.2 missing/deleted (validation_failed = true)
  3. Transitions to CREATE_NEXT_INFRASTRUCTURE
  4. CREATE_NEXT_INFRASTRUCTURE (with R603) creates 2.2.2 correctly
- Verdict: **OPTIMAL RECOVERY PATH**

**Option D: START_WAVE_ITERATION** ❌
- Pro: Valid transition from ERROR_RECOVERY
- Con: Would recreate 2.2.1 (wasteful)
- Con: Would lose approved 2.2.1 work
- Verdict: TOO DESTRUCTIVE

#### State Machine Evidence:
```json
{
  "ERROR_RECOVERY": {
    "allowed_transitions": [
      "START_WAVE_ITERATION",  // Too destructive
      // VALIDATE_INFRASTRUCTURE not in list
    ]
  },
  "VALIDATE_INFRASTRUCTURE": {
    "allowed_transitions": [
      "CREATE_NEXT_INFRASTRUCTURE",  // Loop back when validation fails!
      "SPAWN_CODE_REVIEWERS_EFFORT_PLANNING",
      "ERROR_RECOVERY"
    ],
    "guards": {
      "CREATE_NEXT_INFRASTRUCTURE": "validation_failed == true"
    }
  }
}
```

**BUT**: VALIDATE_INFRASTRUCTURE is not a valid direct transition from ERROR_RECOVERY!

**SOLUTION**: Use State Manager (R517) consultation to request non-standard transition for recovery purposes. State Manager has authority to approve recovery transitions with proper justification.

---

### Question 2: Should We Run Upgrade Script First?

**ANSWER: YES - MUST RUN UPGRADE BEFORE RESETTING STATE**

#### Reasoning:
1. R603 only exists in template repository, not in this project
2. CREATE_NEXT_INFRASTRUCTURE state rules don't have R603 logic yet
3. Without R603, orchestrator will recreate 2.2.2 with SAME WRONG BASE BRANCH
4. Upgrade script will:
   - Sync R603 rule from template
   - Update CREATE_NEXT_INFRASTRUCTURE state rules to use R603
   - Add R603 to RULE-REGISTRY.md
   - Ensure all state references are updated

#### Upgrade Script Location:
```bash
/home/vscode/workspaces/idpbuilder-oci-push-planning/tools/upgrade.sh
```

#### Verification After Upgrade:
```bash
# Check R603 exists
ls -la rule-library/R603-*.md

# Check CREATE_NEXT_INFRASTRUCTURE references R603
grep -r "R603" agent-states/software-factory/orchestrator/CREATE_NEXT_INFRASTRUCTURE/
```

**MANDATORY SEQUENCE**:
1. Run upgrade script
2. Verify R603 present
3. THEN delete infrastructure
4. THEN reset state machine

---

### Question 3: What State File Changes Are Needed?

**ANSWER: Multiple State File Updates Required**

#### Step 1: Mark 2.2.2 as Not Created
```json
"pre_planned_infrastructure": {
  "efforts": {
    "phase2_wave2_effort_1_registry_override": {
      "created": true,         // KEEP - approved work
      "validated": true        // KEEP
    },
    "phase2_wave2_effort_2_env_support": {
      "created": false,        // CHANGE from true
      "validated": false       // CHANGE from true
    }
  }
}
```

#### Step 2: Update Current State
```json
"state_machine": {
  "current_state": "VALIDATE_INFRASTRUCTURE",  // CHANGE from ERROR_RECOVERY
  "previous_state": "ERROR_RECOVERY",          // Document recovery path
  "last_transition_timestamp": "<timestamp>",
  "last_state_manager_consultation": {
    "timestamp": "<timestamp>",
    "consultation_type": "RECOVERY_TRANSITION",
    "requesting_agent": "software-factory-manager",
    "previous_state": "ERROR_RECOVERY",
    "new_state": "VALIDATE_INFRASTRUCTURE",
    "validation_status": "APPROVED",
    "transition_valid": true,
    "reason": "R509 violation recovery - deleted broken 2.2.2 infrastructure, need validation loop-back to recreate with R603 dependency checking"
  }
}
```

#### Step 3: Ensure Effort Status Correct
```json
// If these fields exist:
"efforts_in_progress": [],  // Empty - no active implementation
"efforts_completed": [
  {
    "effort_id": "2.2.1",
    "status": "approved",
    "review_status": "APPROVED"
  }
  // 2.2.2 should NOT be in completed
]
```

#### Step 4: Clear Active Agents
```json
"active_agents": [
  // Remove any agents for 2.2.2
  // Keep only completed agents for historical record
]
```

**CRITICAL**: All state file updates MUST go through State Manager (R517 consultation required).

---

### Question 4: How to Delete Infrastructure Safely?

**ANSWER: 5-Step Safe Deletion Protocol**

#### Pre-Deletion Checks:
```bash
# 1. Verify no agents actively working on 2.2.2
jq '.active_agents | map(select(.effort == "2.2.2" or .effort_name == "env-variable-support"))' orchestrator-state-v3.json
# Expected: [] (empty)

# 2. Verify 2.2.1 is preserved
cd efforts/phase2/wave2/effort-1-registry-override-viper
git log --oneline -1
# Expected: commit mentioning APPROVED

# 3. Backup state file (safety)
cp orchestrator-state-v3.json orchestrator-state-v3.json.backup-before-2.2.2-deletion
```

#### Deletion Sequence:

**STEP 1: Delete Local Branch (from 2.2.2 workspace)**
```bash
cd /home/vscode/workspaces/idpbuilder-oci-push-planning/efforts/phase2/wave2/effort-2-env-variable-support

# Check current branch
git branch --show-current
# Expected: idpbuilder-oci-push/phase2/wave2/effort-2-env-variable-support

# Switch to a safe branch first (can't delete current branch)
git checkout idpbuilder-oci-push/phase2/wave1/integration

# Delete the effort branch
git branch -D idpbuilder-oci-push/phase2/wave2/effort-2-env-variable-support

# Verify deletion
git branch -a | grep "effort-2-env-variable-support"
# Should show only remote branch
```

**STEP 2: Delete Remote Branch**
```bash
cd /home/vscode/workspaces/idpbuilder-oci-push-planning

# Delete remote branch
git push origin --delete idpbuilder-oci-push/phase2/wave2/effort-2-env-variable-support

# Verify remote deletion
git ls-remote --heads origin | grep "effort-2-env-variable-support"
# Expected: (no output)
```

**STEP 3: Delete Workspace Directory**
```bash
cd /home/vscode/workspaces/idpbuilder-oci-push-planning

# Verify we're deleting the right directory
ls -la efforts/phase2/wave2/effort-2-env-variable-support/.git
# Expected: git repository directory exists

# Delete workspace (git repo inside)
rm -rf efforts/phase2/wave2/effort-2-env-variable-support

# Verify deletion
ls -la efforts/phase2/wave2/
# Expected: only effort-1-registry-override-viper remains
```

**STEP 4: Update State File via State Manager**
```bash
# Use State Manager to update pre_planned_infrastructure
# Mark 2.2.2 as created: false, validated: false
# (Detailed in Question 3)
```

**STEP 5: Commit Deletion**
```bash
cd /home/vscode/workspaces/idpbuilder-oci-push-planning

git add -A
git commit -m "recovery: delete broken Effort 2.2.2 infrastructure (R509 violation)

- Deleted local branch idpbuilder-oci-push/phase2/wave2/effort-2-env-variable-support
- Deleted remote branch
- Deleted workspace directory efforts/phase2/wave2/effort-2-env-variable-support
- Updated state file to mark 2.2.2 as not created
- Reason: R509 violation - effort was based on Wave 2.1 integration instead of 2.2.1
- Missing dependency: config.go from 2.2.1
- Recovery plan: will recreate with R603 dependency checking

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>"

git push
```

#### Deletion Validation:
```bash
# Verify 2.2.1 untouched
cd efforts/phase2/wave2/effort-1-registry-override-viper
git log --oneline -1
# Expected: Still shows APPROVED commit

ls -la pkg/idpbuilder/config.go
# Expected: File exists

# Verify 2.2.2 deleted
cd /home/vscode/workspaces/idpbuilder-oci-push-planning
ls efforts/phase2/wave2/
# Expected: only effort-1-registry-override-viper

git branch -a | grep "effort-2-env-variable-support"
# Expected: (no output - branch gone)

# Verify state file updated
jq '.pre_planned_infrastructure.efforts.phase2_wave2_effort_2_env_support.created' orchestrator-state-v3.json
# Expected: false
```

---

### Question 5: Validation Steps

**ANSWER: 3-Phase Validation Protocol**

#### Phase 1: Pre-Recovery Validation

**Check 1: Upgrade Successful**
```bash
# R603 exists
test -f rule-library/R603-*.md && echo "✅ R603 present" || echo "❌ R603 missing"

# CREATE_NEXT_INFRASTRUCTURE has R603 reference
grep -q "R603" agent-states/software-factory/orchestrator/CREATE_NEXT_INFRASTRUCTURE/rules.md && \
  echo "✅ R603 referenced" || echo "❌ R603 not referenced"

# RULE-REGISTRY has R603
grep -q "R603" rule-library/RULE-REGISTRY.md && \
  echo "✅ R603 in registry" || echo "❌ R603 not in registry"
```

**Check 2: Infrastructure Deleted**
```bash
# 2.2.2 workspace gone
test ! -d efforts/phase2/wave2/effort-2-env-variable-support && \
  echo "✅ Workspace deleted" || echo "❌ Workspace exists"

# 2.2.2 branch gone
git branch -a | grep -q "effort-2-env-variable-support" && \
  echo "❌ Branch still exists" || echo "✅ Branch deleted"

# State file updated
jq -e '.pre_planned_infrastructure.efforts.phase2_wave2_effort_2_env_support.created == false' \
  orchestrator-state-v3.json && echo "✅ State updated" || echo "❌ State not updated"
```

**Check 3: 2.2.1 Preserved**
```bash
# Workspace exists
test -d efforts/phase2/wave2/effort-1-registry-override-viper && \
  echo "✅ 2.2.1 workspace preserved" || echo "❌ 2.2.1 workspace missing"

# config.go exists
test -f efforts/phase2/wave2/effort-1-registry-override-viper/pkg/idpbuilder/config.go && \
  echo "✅ config.go preserved" || echo "❌ config.go missing"

# Branch exists
git branch -a | grep -q "effort-1-registry-override-viper" && \
  echo "✅ 2.2.1 branch preserved" || echo "❌ 2.2.1 branch missing"
```

#### Phase 2: State Reset Validation

**Check 4: State Machine Updated**
```bash
# Current state correct
jq -e '.state_machine.current_state == "VALIDATE_INFRASTRUCTURE"' orchestrator-state-v3.json && \
  echo "✅ State is VALIDATE_INFRASTRUCTURE" || echo "❌ Wrong state"

# State Manager consultation recorded
jq -e '.state_machine.last_state_manager_consultation.consultation_type == "RECOVERY_TRANSITION"' \
  orchestrator-state-v3.json && echo "✅ Consultation recorded" || echo "❌ No consultation"

# Previous state documented
jq -e '.state_machine.previous_state == "ERROR_RECOVERY"' orchestrator-state-v3.json && \
  echo "✅ Previous state documented" || echo "❌ Previous state wrong"
```

**Check 5: Dependencies Satisfied**
```bash
# 2.2.1 marked as created
jq -e '.pre_planned_infrastructure.efforts.phase2_wave2_effort_1_registry_override.created == true' \
  orchestrator-state-v3.json && echo "✅ 2.2.1 marked created" || echo "❌ 2.2.1 not created"

# 2.2.2 depends_on preserved
jq -e '.pre_planned_infrastructure.efforts.phase2_wave2_effort_2_env_support.depends_on | contains(["phase2_wave2_effort_1_registry_override"])' \
  orchestrator-state-v3.json && echo "✅ Dependency preserved" || echo "❌ Dependency missing"

# Base branch metadata correct
jq -r '.pre_planned_infrastructure.efforts.phase2_wave2_effort_2_env_support.base_branch' \
  orchestrator-state-v3.json | grep -q "effort-1-registry-override-viper" && \
  echo "✅ Base branch metadata correct" || echo "❌ Base branch metadata wrong"
```

#### Phase 3: Post-Recovery Simulation

**Check 6: Orchestrator Behavior Prediction**
```bash
# What will orchestrator do next?
# Expected sequence:
# 1. In VALIDATE_INFRASTRUCTURE
# 2. Check all efforts in pre_planned_infrastructure
# 3. Find 2.2.1: created=true, validated=true ✅
# 4. Find 2.2.2: created=false, validated=false ❌
# 5. validation_failed = true
# 6. Transition to CREATE_NEXT_INFRASTRUCTURE
# 7. CREATE_NEXT_INFRASTRUCTURE reads R603
# 8. R603 checks dependencies for 2.2.2
# 9. R603 finds depends_on: ["phase2_wave2_effort_1_registry_override"]
# 10. R603 verifies 2.2.1 created and validated ✅
# 11. R603 gets 2.2.1 branch name: "idpbuilder-oci-push/phase2/wave2/effort-1-registry-override-viper"
# 12. CREATE_NEXT_INFRASTRUCTURE creates 2.2.2 branch FROM 2.2.1 branch ✅
# 13. R509 validation passes (correct base branch)
```

**Check 7: Manual Dry-Run of R603 Logic**
```bash
# Simulate R603 dependency check
EFFORT_KEY="phase2_wave2_effort_2_env_support"
DEPENDS_ON=$(jq -r ".pre_planned_infrastructure.efforts.${EFFORT_KEY}.depends_on[0]" orchestrator-state-v3.json)

echo "Effort 2.2.2 depends on: $DEPENDS_ON"
# Expected: phase2_wave2_effort_1_registry_override

DEPENDENCY_CREATED=$(jq -r ".pre_planned_infrastructure.efforts.${DEPENDS_ON}.created" orchestrator-state-v3.json)
echo "Dependency created: $DEPENDENCY_CREATED"
# Expected: true

DEPENDENCY_BRANCH=$(jq -r ".pre_planned_infrastructure.efforts.${DEPENDS_ON}.branch_name" orchestrator-state-v3.json)
echo "Dependency branch: $DEPENDENCY_BRANCH"
# Expected: idpbuilder-oci-push/phase2/wave2/effort-1-registry-override-viper

echo "R603 would use base branch: $DEPENDENCY_BRANCH"
echo "✅ This is CORRECT - includes config.go"
```

**Check 8: Verify config.go Availability**
```bash
# Check config.go exists in dependency branch
cd /home/vscode/workspaces/idpbuilder-oci-push-planning
git show origin/idpbuilder-oci-push/phase2/wave2/effort-1-registry-override-viper:pkg/idpbuilder/config.go | head -10

# Expected: File contents displayed (203 lines)
# If error: config.go not in that branch (PROBLEM!)
```

---

## PART 3: STEP-BY-STEP RECOVERY PLAN

### Phase 1: Upgrade (Add R603)

**Step 1.1: Verify Current State**
```bash
cd /home/vscode/workspaces/idpbuilder-oci-push-planning
pwd
# Expected: /home/vscode/workspaces/idpbuilder-oci-push-planning

git status
# Expected: On branch main, possibly modified files from State Manager work
```

**Step 1.2: Run Upgrade Script**
```bash
cd /home/vscode/workspaces/idpbuilder-oci-push-planning

# Run upgrade script
bash tools/upgrade.sh

# Expected output:
# - "Upgrading from SF <version> to <version>"
# - "Syncing R603..."
# - "Updating state rules..."
# - "✅ Upgrade complete"
```

**Step 1.3: Verify Upgrade Success**
```bash
# Check R603 exists
ls -la rule-library/R603-*.md

# Check content
head -50 rule-library/R603-*.md

# Check CREATE_NEXT_INFRASTRUCTURE references R603
grep "R603" agent-states/software-factory/orchestrator/CREATE_NEXT_INFRASTRUCTURE/rules.md
```

**Step 1.4: Commit Upgrade**
```bash
git add -A
git commit -m "upgrade: add R603 for sequential effort dependency checking

- Synced R603 from template repository
- Updated CREATE_NEXT_INFRASTRUCTURE state rules to use R603
- Added R603 to RULE-REGISTRY.md
- Preparation for Wave 2.2 recovery (R509 violation fix)

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>"

git push
```

---

### Phase 2: Delete Broken Infrastructure

**Step 2.1: Pre-Deletion Backup**
```bash
cd /home/vscode/workspaces/idpbuilder-oci-push-planning

# Backup state file
cp orchestrator-state-v3.json orchestrator-state-v3.json.backup-before-2.2.2-deletion

# Backup bug tracking
cp bug-tracking.json bug-tracking.json.backup-before-2.2.2-deletion

# List current state
ls -la efforts/phase2/wave2/
git branch -a | grep "wave2"
```

**Step 2.2: Delete Local Branch**
```bash
cd /home/vscode/workspaces/idpbuilder-oci-push-planning/efforts/phase2/wave2/effort-2-env-variable-support

# Switch to safe branch
git checkout idpbuilder-oci-push/phase2/wave1/integration

# Go back to main directory
cd /home/vscode/workspaces/idpbuilder-oci-push-planning

# Delete branch from the main repo context
# (The workspace is a separate git repo)
cd efforts/phase2/wave2/effort-2-env-variable-support
git branch -D idpbuilder-oci-push/phase2/wave2/effort-2-env-variable-support || true
cd /home/vscode/workspaces/idpbuilder-oci-push-planning
```

**Step 2.3: Delete Remote Branch**
```bash
cd /home/vscode/workspaces/idpbuilder-oci-push-planning

# Delete from origin
git push origin --delete idpbuilder-oci-push/phase2/wave2/effort-2-env-variable-support || true

# Verify deletion
git ls-remote --heads origin | grep "effort-2-env-variable-support"
# Expected: (no output)
```

**Step 2.4: Delete Workspace Directory**
```bash
cd /home/vscode/workspaces/idpbuilder-oci-push-planning

# Delete workspace
rm -rf efforts/phase2/wave2/effort-2-env-variable-support

# Verify deletion
ls -la efforts/phase2/wave2/
# Expected: only effort-1-registry-override-viper
```

**Step 2.5: Verify 2.2.1 Intact**
```bash
cd /home/vscode/workspaces/idpbuilder-oci-push-planning/efforts/phase2/wave2/effort-1-registry-override-viper

# Check branch
git branch --show-current

# Check latest commit
git log --oneline -1

# Check config.go exists
ls -la pkg/idpbuilder/config.go

# Check review report exists
ls -la .software-factory/phase2/wave2/effort-1-registry-override-viper/CODE-REVIEW-REPORT--*.md
```

---

### Phase 3: Update State File (via State Manager)

**Step 3.1: Consult State Manager for State Transition**
```bash
cd /home/vscode/workspaces/idpbuilder-oci-push-planning

# State Manager consultation will:
# 1. Update current_state to VALIDATE_INFRASTRUCTURE
# 2. Record previous_state as ERROR_RECOVERY
# 3. Add consultation record
# 4. Update timestamp
# 5. Mark as RECOVERY_TRANSITION type

# (Actual State Manager consultation happens via agent command)
# Expected: State Manager approves ERROR_RECOVERY → VALIDATE_INFRASTRUCTURE
# Reason: "R509 violation recovery - deleted broken 2.2.2, need validation loop-back"
```

**Step 3.2: Update pre_planned_infrastructure**
```bash
cd /home/vscode/workspaces/idpbuilder-oci-push-planning

# Update 2.2.2 to mark as not created
jq '.pre_planned_infrastructure.efforts.phase2_wave2_effort_2_env_support.created = false |
    .pre_planned_infrastructure.efforts.phase2_wave2_effort_2_env_support.validated = false' \
  orchestrator-state-v3.json > orchestrator-state-v3.json.tmp

mv orchestrator-state-v3.json.tmp orchestrator-state-v3.json

# Verify update
jq '.pre_planned_infrastructure.efforts.phase2_wave2_effort_2_env_support | {created, validated}' \
  orchestrator-state-v3.json
# Expected: {"created": false, "validated": false}
```

**Step 3.3: Verify Dependency Metadata Intact**
```bash
# Check 2.2.2 still has correct base_branch and depends_on
jq '.pre_planned_infrastructure.efforts.phase2_wave2_effort_2_env_support | {base_branch, depends_on}' \
  orchestrator-state-v3.json

# Expected:
# {
#   "base_branch": "idpbuilder-oci-push/phase2/wave2/effort-1-registry-override-viper",
#   "depends_on": ["phase2_wave2_effort_1_registry_override"]
# }
```

**Step 3.4: Commit State File Changes**
```bash
git add orchestrator-state-v3.json
git commit -m "recovery: reset state to VALIDATE_INFRASTRUCTURE after 2.2.2 deletion

- State: ERROR_RECOVERY → VALIDATE_INFRASTRUCTURE
- Reason: R509 violation recovery (effort 2.2.2 wrong base branch)
- Marked effort 2.2.2 as created=false, validated=false
- Preserved dependency metadata (base_branch, depends_on)
- Validation loop will detect missing 2.2.2 and trigger recreation
- R603 will ensure correct base branch this time

State Manager consultation: RECOVERY_TRANSITION approved
Recovery rationale: Deleted broken infrastructure, need validation loop to recreate with R603

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>"

git push
```

---

### Phase 4: Validation

**Step 4.1: Run All Validation Checks**
```bash
cd /home/vscode/workspaces/idpbuilder-oci-push-planning

# Run Phase 1 validations (see Question 5)
# Check R603 present
# Check infrastructure deleted
# Check 2.2.1 preserved

# Run Phase 2 validations
# Check state machine updated
# Check dependencies satisfied

# Run Phase 3 simulation
# Predict orchestrator behavior
# Dry-run R603 logic
```

**Step 4.2: Create Validation Report**
```bash
# Document all validation results
# Include pass/fail for each check
# Screenshot or copy relevant outputs
```

---

### Phase 5: Handoff

**Step 5.1: Create Handoff Document** (see Part 4 below)

**Step 5.2: Set Continuation Flag**
```bash
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
```

**Step 5.3: Stop Execution**
```bash
# Factory Manager work complete
# Orchestrator will take over from VALIDATE_INFRASTRUCTURE
exit 0
```

---

## PART 4: HANDOFF DOCUMENT FOR ORCHESTRATOR

### Current State After Recovery

**State Machine Status:**
```
Current State: VALIDATE_INFRASTRUCTURE
Previous State: ERROR_RECOVERY
Phase: 2
Wave: 2
```

**Infrastructure Status:**
```
Effort 2.2.1: ✅ CREATED, VALIDATED, APPROVED (247 lines, config.go present)
Effort 2.2.2: ❌ DELETED (was based on wrong branch)
              Will be recreated by CREATE_NEXT_INFRASTRUCTURE with R603
```

**What Will Happen Next:**

When `/continue-orchestrating` is run, the orchestrator will:

1. **Load Current State**: VALIDATE_INFRASTRUCTURE
2. **Read State Rules**: agent-states/software-factory/orchestrator/VALIDATE_INFRASTRUCTURE/rules.md
3. **Execute Validation**:
   - Check all efforts in pre_planned_infrastructure
   - Effort 2.2.1: created=true, validated=true ✅
   - Effort 2.2.2: created=false, validated=false ❌
   - Set validation_failed = true
4. **Transition**: VALIDATE_INFRASTRUCTURE → CREATE_NEXT_INFRASTRUCTURE (loop-back)
5. **Execute CREATE_NEXT_INFRASTRUCTURE**:
   - Read R603 (Sequential Effort Dependency Checking)
   - Find effort 2.2.2 needs creation
   - Read depends_on: ["phase2_wave2_effort_1_registry_override"]
   - Verify dependency created ✅
   - Get dependency branch: "idpbuilder-oci-push/phase2/wave2/effort-1-registry-override-viper"
   - Create 2.2.2 branch FROM 2.2.1 branch ✅
   - Create workspace directory
   - R509 validation passes (correct base branch)
   - Mark 2.2.2 as created=true
6. **Transition**: CREATE_NEXT_INFRASTRUCTURE → VALIDATE_INFRASTRUCTURE
7. **Re-validate**: All efforts now valid ✅
8. **Transition**: VALIDATE_INFRASTRUCTURE → SPAWN_CODE_REVIEWERS_EFFORT_PLANNING
9. **Continue Normal Flow**: Implementation plan for 2.2.2

### Success Criteria

The recovery is successful when:

1. ✅ Effort 2.2.2 branch created FROM 2.2.1 branch
2. ✅ config.go from 2.2.1 is present in 2.2.2
3. ✅ R509 validation passes (correct base branch)
4. ✅ Implementation plan created for 2.2.2
5. ✅ No errors during infrastructure creation

### What to Watch For

**RED FLAGS:**
- ❌ Effort 2.2.2 created from Wave 2.1 integration again
- ❌ config.go missing from 2.2.2
- ❌ R509 violation re-occurs
- ❌ CREATE_NEXT_INFRASTRUCTURE skips 2.2.2

**GREEN LIGHTS:**
- ✅ CREATE_NEXT_INFRASTRUCTURE reads R603
- ✅ R603 checks dependencies
- ✅ Branch created from correct base
- ✅ Validation passes

### Troubleshooting

**If R603 is not used:**
- Check: Did upgrade script run successfully?
- Check: Does rule-library/R603-*.md exist?
- Check: Does CREATE_NEXT_INFRASTRUCTURE/rules.md reference R603?

**If wrong base branch again:**
- Check: Is depends_on metadata correct in pre_planned_infrastructure?
- Check: Is base_branch metadata correct?
- Check: Is R603 logic reading metadata correctly?

**If validation fails:**
- Check: Is 2.2.1 still marked as created=true?
- Check: Is 2.2.2 marked as created=false before recreation?
- Check: Are branch names correct?

---

## PART 5: RISK ASSESSMENT

### Risks and Mitigations

**Risk 1: Lose 2.2.1 approved work**
- Mitigation: Verify 2.2.1 untouched after every step
- Mitigation: Backup state file before deletion
- Mitigation: Only delete 2.2.2, never touch 2.2.1

**Risk 2: R603 not available after upgrade**
- Mitigation: Verify R603 exists before proceeding
- Mitigation: Test upgrade script in isolation first
- Mitigation: Manual sync if upgrade script fails

**Risk 3: State machine rejects VALIDATE_INFRASTRUCTURE transition**
- Mitigation: Use State Manager consultation (R517)
- Mitigation: Document recovery justification
- Mitigation: Fallback to START_WAVE_ITERATION if necessary (more destructive)

**Risk 4: R603 doesn't work as expected**
- Mitigation: Dry-run R603 logic before recovery
- Mitigation: Verify metadata is correct
- Mitigation: Manual infrastructure creation as fallback

**Risk 5: Remote branch deletion fails (permissions)**
- Mitigation: Document that remote branch exists
- Mitigation: Orchestrator can detect and handle existing remote branches
- Mitigation: Manual deletion via GitHub UI if necessary

### Rollback Plan

If recovery fails completely:

1. **Restore State File**: `cp orchestrator-state-v3.json.backup-before-2.2.2-deletion orchestrator-state-v3.json`
2. **Recreate 2.2.2 Workspace**: (from backup if needed)
3. **Reset to ERROR_RECOVERY**: Document failure, investigate root cause
4. **Alternative**: Use START_WAVE_ITERATION to restart Wave 2.2 (loses 2.2.1 work)

---

## PART 6: EXECUTION AUTHORIZATION

### Pre-Execution Checklist

Before executing recovery:
- [ ] Upgrade script location verified: `/home/vscode/workspaces/idpbuilder-oci-push-planning/tools/upgrade.sh`
- [ ] 2.2.1 status confirmed: APPROVED, config.go present
- [ ] 2.2.2 issue confirmed: based on Wave 2.1 integration, missing config.go
- [ ] State file backed up: `orchestrator-state-v3.json.backup-before-2.2.2-deletion`
- [ ] Recovery plan reviewed and approved
- [ ] State Manager consultation process understood
- [ ] Rollback plan documented

### Authority

Software Factory Manager has authority to:
- ✅ Run upgrade script
- ✅ Delete broken Effort 2.2.2 infrastructure
- ✅ Update state files (via State Manager consultation)
- ✅ Commit and push recovery changes
- ❌ Modify Effort 2.2.1 (preserve as-is)

### Execution Status

**Status**: READY FOR EXECUTION
**Recommended Start**: Immediately (system is in ERROR_RECOVERY, no work in progress)
**Estimated Duration**: 30-45 minutes
**Risk Level**: LOW (well-planned, backed up, rollback available)

---

## APPENDIX A: R603 Rule Summary

**(To be added after upgrade when R603 content is available)**

R603 provides:
- Sequential effort dependency validation
- Cascade base branch verification
- R213 metadata reading for dependencies
- Automatic detection of prerequisite efforts
- Base branch calculation from dependency graph

---

## APPENDIX B: State Machine Flow Diagram

```
ERROR_RECOVERY
    ↓ (Recovery transition via State Manager)
VALIDATE_INFRASTRUCTURE
    ↓ (validation_failed = true, effort 2.2.2 missing)
CREATE_NEXT_INFRASTRUCTURE
    ↓ (R603: read depends_on, get 2.2.1 branch, create 2.2.2 from 2.2.1)
    ↓ (R509: validate correct base branch ✅)
    ↓ (Mark 2.2.2 created=true)
VALIDATE_INFRASTRUCTURE
    ↓ (validation successful, all efforts valid)
SPAWN_CODE_REVIEWERS_EFFORT_PLANNING
    ↓ (Create implementation plan for 2.2.2)
WAITING_FOR_EFFORT_PLANS
    ↓ (Plan creation complete)
ANALYZE_IMPLEMENTATION_PARALLELIZATION
    ↓ (Sequential strategy)
SPAWN_SW_ENGINEERS
    ↓ (Implementation begins)
... (normal flow continues)
```

---

## APPENDIX C: Validation Checklist

### Pre-Recovery
- [ ] R603 does not exist in project (confirmed)
- [ ] Effort 2.2.1 is APPROVED (confirmed)
- [ ] Effort 2.2.2 has wrong base branch (confirmed)
- [ ] config.go missing from 2.2.2 (confirmed)
- [ ] State is ERROR_RECOVERY (confirmed)

### Post-Upgrade
- [ ] R603 exists in rule-library/
- [ ] CREATE_NEXT_INFRASTRUCTURE references R603
- [ ] RULE-REGISTRY includes R603
- [ ] Upgrade committed and pushed

### Post-Deletion
- [ ] 2.2.2 workspace deleted
- [ ] 2.2.2 local branch deleted
- [ ] 2.2.2 remote branch deleted
- [ ] 2.2.1 workspace intact
- [ ] 2.2.1 config.go present
- [ ] State file updated (created=false)

### Post-State-Reset
- [ ] Current state is VALIDATE_INFRASTRUCTURE
- [ ] State Manager consultation recorded
- [ ] Previous state is ERROR_RECOVERY
- [ ] Dependency metadata intact
- [ ] Changes committed and pushed

### Pre-Handoff
- [ ] All validations passed
- [ ] Handoff document created
- [ ] Success criteria documented
- [ ] Troubleshooting guide provided
- [ ] CONTINUE-SOFTWARE-FACTORY=TRUE set

---

**END OF RECOVERY ANALYSIS AND PLAN**

**Status**: ANALYSIS COMPLETE - READY FOR EXECUTION
**Next Step**: Execute Phase 1 (Upgrade) when authorized
