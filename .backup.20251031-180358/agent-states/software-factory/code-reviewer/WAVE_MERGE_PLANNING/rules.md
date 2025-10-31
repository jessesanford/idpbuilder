# Code Reviewer - WAVE_MERGE_PLANNING State Rules

## 🚨 CRITICAL: MERGE PLAN CREATION ONLY - NO EXECUTION! 🚨

### PRIMARY PURPOSE
Create a detailed MERGE PLAN document that the Integration Agent will follow.
DO NOT execute any merges - only plan them!

### 🔴🔴🔴 CRITICAL RULE: NO INTEGRATE_WAVE_EFFORTS BRANCHES IN MERGE PLAN! 🔴🔴🔴

**YOU MUST USE ONLY ORIGINAL EFFORT BRANCHES!**
- ✅ CORRECT: phase3/wave2/effort1-api-gateway
- ✅ CORRECT: phase3/wave2/effort2-controller-split1  
- ❌ WRONG: wave2-integration-20250827
- ❌ WRONG: phase3-integration-20250827

Integration branches are TARGETS, not SOURCES! Never merge from an integration branch.

## State Context
You are creating a comprehensive merge plan for wave integration. The orchestrator has already set up the integration infrastructure (directory and branch). Your role is to analyze what needs to be merged and create a detailed plan.

## 🔴🔴🔴 CRITICAL: FIRST ACTION - CD TO INTEGRATE_WAVE_EFFORTS DIRECTORY 🔴🔴🔴

**YOUR FIRST ACTION MUST BE:**
```bash
# The orchestrator will tell you the integration directory path
# You MUST CD to this directory before creating WAVE-MERGE-PLAN.md
cd [integration directory path from orchestrator]
pwd  # Verify you are in the integration workspace
```

**NEVER create WAVE-MERGE-PLAN.md in the root directory!**
**ALWAYS create it IN the integration workspace!**

## 🔴🔴🔴 WAVE SEQUENTIAL REBUILD MODEL (R009/R512) 🔴🔴🔴

### Critical Understanding

**WAVE INTEGRATE_WAVE_EFFORTS USES SEQUENTIAL REBUILD:**
- **Base**: FIRST effort of THIS wave (not previous wave integration!)
- **Merge**: Remaining efforts of THIS wave only
- **Purpose**: Test wave-internal sequential mergeability
- **Scope**: Wave-scoped (independent from other waves)

**Example:**
```bash
Wave 2 has efforts: E4, E5, E6
Base: E4 (first effort of wave 2)
Merge: E5, E6 (remaining efforts)
Result: Tests E4→E5→E6 sequential merge
```

**NOT:**
```bash
Base: wave1-integration (WRONG! This is phase integration model)
Merge: E4, E5, E6
```

## MERGE PLAN REQUIREMENTS

### 1. Determine Base Branch (First Effort of Wave)
```bash
#!/bin/bash
# 🔴 SEQUENTIAL REBUILD: Find base and source branches for wave integration
PHASE=$(jq -r '.state_machine.current_phase' "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json")
WAVE=$(jq -r '.state_machine.current_wave' "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json")

echo "📊 Wave Sequential Rebuild Analysis for Phase $PHASE Wave $WAVE..."

# Step 1: Get FIRST effort of this wave (BASE)
FIRST_EFFORT=$(jq -r --arg p "$PHASE" --arg w "$WAVE" '
    [.efforts_completed[], .efforts_in_progress[]] |
    map(select(.phase == ($p | tonumber) and .wave == ($w | tonumber))) |
    sort_by(.index // .branched_at) |
    .[0].branch
' "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json")

if [ -z "$FIRST_EFFORT" ] || [ "$FIRST_EFFORT" = "null" ]; then
    echo "❌ FATAL: Cannot find first effort of wave"
    exit 1
fi

echo "🔴 BASE BRANCH: $FIRST_EFFORT (first effort of wave)"
echo ""

# Step 2: Get REMAINING efforts of this wave (SOURCES)
echo "📋 SOURCE BRANCHES (remaining efforts to merge):"
jq -r --arg p "$PHASE" --arg w "$WAVE" --arg first "$FIRST_EFFORT" '
    [.efforts_completed[], .efforts_in_progress[]] |
    map(select(.phase == ($p | tonumber) and .wave == ($w | tonumber))) |
    map(select(.branch != $first)) |
    sort_by(.index // .branched_at) |
    .[] | .branch
' "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json" | while read branch; do
    # Check if deprecated/split
    DEPRECATED=$(jq -r --arg b "$branch" '
        .deprecated_branches[] | select(.original_branch == $b) | .deprecated
    ' "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json" 2>/dev/null)

    if [[ "$DEPRECATED" == "true" ]]; then
        echo "❌ EXCLUDE: $branch (deprecated, was split)"
        # Find splits to include instead
        jq -r --arg b "$branch" '
            .deprecated_branches[] |
            select(.original_branch == $b) |
            .split_branches[]
        ' "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json" 2>/dev/null | while read split; do
            echo "  ✅ INCLUDE: $split (split of $branch)"
        done
    else
        echo "✅ INCLUDE: $branch"
    fi
done
```

### 2. Determine Merge Order
```markdown
## Merge Order Analysis

The correct merge order is critical to avoid conflicts and preserve intent.

### Checking Branch Bases:
- Branch A based on: main at commit abc123  
- Branch B based on: main at commit def456
- Branch C based on: Branch A at commit ghi789

### Therefore Order Must Be:
1. Branch A (base: main)
2. Branch C (base: Branch A) - depends on A
3. Branch B (base: main) - independent

### Dependency Graph:
```
main
├── Branch A
│   └── Branch C (depends on A)
└── Branch B (independent)
```
```

### 3. Handle Splits Correctly

#### 🔴🔴🔴 CRITICAL: Split Merge Ordering with Dependencies 🔴🔴🔴

**SUPREME LAW: Dependencies are at EFFORT level, not SPLIT level!**

When creating merge plans with splits and dependencies:
1. ALL splits of an effort must be listed sequentially
2. ALL splits must complete before dependent efforts
3. NEVER interleave dependent efforts between splits

```markdown
## Split Handling

### Effort E3.1.3 - Contexts Implementation
Original branch exceeded 800 lines and was split:

**EXCLUDE (too large):**
- phase3/wave1/effort3-contexts (1,234 lines)

**INCLUDE IN THIS ORDER (properly sized splits):**  
1. phase3/wave1/effort3-contexts-split-001 (423 lines) - base: main
2. phase3/wave1/effort3-contexts-split-002 (389 lines) - base: split-001
3. phase3/wave1/effort3-contexts-split-003 (401 lines) - base: split-002

Total after splits: 1,213 lines (compliant)

### Dependency Handling Example

If E3.1.4 depends on E3.1.3:
✅ CORRECT ORDER:
1. effort3-contexts-split-001
2. effort3-contexts-split-002
3. effort3-contexts-split-003
4. effort4-dependent-feature  # NOW has complete E3.1.3

❌ WRONG ORDER:
1. effort3-contexts-split-001
2. effort4-dependent-feature  # ERROR: E3.1.3 incomplete!
3. effort3-contexts-split-002
4. effort3-contexts-split-003
```

### 4. Create WAVE-MERGE-PLAN.md

**CRITICAL LOCATION REQUIREMENT:**
- **YOU MUST BE IN THE INTEGRATE_WAVE_EFFORTS DIRECTORY** (cd there first!)
- **Create the file IN the current directory** (./WAVE-MERGE-PLAN.md)
- **DO NOT create it in the root directory or anywhere else!**

```bash
# Verify you are in the integration workspace
pwd  # Should show: /efforts/phase${PHASE}/wave${WAVE}/integration-workspace

# Create the merge plan with proper metadata path
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
MERGE_PLAN=".software-factory/phase${PHASE}/wave${WAVE}/integration/WAVE-MERGE-PLAN--${TIMESTAMP}.md"
mkdir -p "$(dirname "$MERGE_PLAN")"
cat > "$MERGE_PLAN" << 'EOF'
[merge plan content]
EOF
```

```markdown
# Wave ${WAVE} Merge Plan

**Generated:** $(date -u +%Y-%m-%dT%H:%M:%SZ)
**Code Reviewer:** code-reviewer
**State:** WAVE_MERGE_PLANNING

## Target Integration Branch (Sequential Rebuild)
- **Branch Name:** phase${PHASE}-wave${WAVE}-integration-${timestamp}
- **Base:** ${FIRST_EFFORT_BRANCH} (first effort of wave ${WAVE})
- **Model:** Wave-scoped sequential rebuild per R009/R512
- **Location:** /efforts/phase${PHASE}/wave${WAVE}/integration-workspace

## Sequential Rebuild Architecture
**This wave integration uses wave-scoped sequential rebuild:**
- Base branch: First effort of THIS wave only
- Source branches: Remaining efforts of THIS wave only
- Purpose: Test wave-internal sequential mergeability
- Independent: Each wave tests its own mergeability
- Phase integration will test cross-wave mergeability separately

## 🎬 Demo Execution Plan (R330/R291 Compliance)

### Demo Requirements Overview
Per R330 and R291, ALL integrations MUST demonstrate working functionality.

### Demo Execution Sequence
1. **After Each Effort Merge:**
   - Run effort-specific demo script if exists
   - Capture output in `demo-results/effort-X-demo.log`
   - Continue even if individual demo fails (document for review)

2. **After All Merges Complete:**
   - Run integrated wave demo: `./wave-demo.sh`
   - Verify all effort features work together
   - Capture evidence in `demo-results/wave-integration-demo.log`

3. **Demo Validation Gates (R291):**
   - ✅ BUILD GATE: Code must compile
   - ✅ TEST GATE: All tests must pass
   - ✅ DEMO GATE: Demo scripts must execute
   - ✅ ARTIFACT GATE: Build outputs must exist

### Demo Files Expected
Based on effort plans (R330), these demos should exist:
- [ ] effort1/demo-features.sh (from effort plan)
- [ ] effort2/demo-features.sh (from effort plan)
- [ ] effort3/demo-features.sh (from effort plan)
- [ ] WAVE-DEMO.md (integration demo documentation)
- [ ] wave-demo.sh (integrated demo script)

### Demo Failure Protocol
If ANY demo fails during integration:
1. Document failure in INTEGRATE_WAVE_EFFORTS_REPORT.md
2. Mark Demo Status: FAILED
3. This will trigger ERROR_RECOVERY per R291
4. Fixes must be made in effort branches (R292)

## Branches to Merge (IN ORDER)

### 1. phase${PHASE}/wave${WAVE}/effort1-api-types
- **Type:** Original effort branch
- **Base:** main at abc123
- **Size:** 542 lines
- **Dependencies:** None
- **Conflicts Expected:** None
- **Merge Command:**
  ```bash
  git fetch origin phase${PHASE}/wave${WAVE}/effort1-api-types
  git merge origin/phase${PHASE}/wave${WAVE}/effort1-api-types --no-ff \
    -m "Integrate effort1-api-types into wave ${WAVE}"
  ```

### 2. phase${PHASE}/wave${WAVE}/effort2-controller-split1
- **Type:** Split branch (1 of 3)
- **Base:** main at abc123  
- **Size:** 398 lines
- **Dependencies:** None
- **Conflicts Expected:** Possible in controller.go
- **Merge Command:**
  ```bash
  git fetch origin phase${PHASE}/wave${WAVE}/effort2-controller-split1
  git merge origin/phase${PHASE}/wave${WAVE}/effort2-controller-split1 --no-ff \
    -m "Integrate effort2-controller-split1 into wave ${WAVE}"
  ```

### 3. phase${PHASE}/wave${WAVE}/effort2-controller-split2
- **Type:** Split branch (2 of 3)
- **Base:** effort2-controller-split1 at def456
- **Size:** 412 lines
- **Dependencies:** Must merge after split1
- **Conflicts Expected:** None (sequential splits)
- **Merge Command:**
  ```bash
  git fetch origin phase${PHASE}/wave${WAVE}/effort2-controller-split2
  git merge origin/phase${PHASE}/wave${WAVE}/effort2-controller-split2 --no-ff \
    -m "Integrate effort2-controller-split2 into wave ${WAVE}"
  ```

[Continue for all branches...]

## Excluded Branches (Too Large)
These original branches were split and should NOT be merged:
- phase${PHASE}/wave${WAVE}/effort2-controller (original, 1,456 lines)

## Merge Strategy
1. **Merge Type:** --no-ff (preserve branch history)
2. **Conflict Resolution:** Favor newer implementation when conflicts occur
3. **Testing:** Run unit tests after each merge
4. **Validation:** Check total size after all merges

## Expected Conflicts
Based on branch analysis, conflicts are likely in:
- `pkg/controller/controller.go` - Modified by both effort2 and effort3
- `api/v1/types.go` - Extended by multiple efforts

## Validation Steps
1. After each merge:
   ```bash
   # Test the code
   go test ./...
   
   # Run effort demo if exists (R330 compliance)
   EFFORT_NAME="[effort-name-from-merge]"
   if [ -f "${EFFORT_NAME}/demo-features.sh" ]; then
       echo "🎬 Running ${EFFORT_NAME} demo..."
       bash "${EFFORT_NAME}/demo-features.sh" > "demo-results/${EFFORT_NAME}-demo.log" 2>&1
       echo "Demo exit code: $?"
   fi
   ```
2. After all merges:
   ```bash
   # Run integration tests
   make test-integration
   
   # Check size compliance
   $PROJECT_ROOT/tools/line-counter.sh -c $(git branch --show-current)
   
   # Run integrated wave demo (R291 requirement)
   echo "🎬 Running integrated wave demo..."
   if [ -f "./wave-demo.sh" ]; then
       bash ./wave-demo.sh > demo-results/wave-integration-demo.log 2>&1
       DEMO_STATUS=$?
       echo "Wave demo status: $DEMO_STATUS"
   else
       echo "⚠️ WARNING: No wave-demo.sh found - creating basic demo"
       # Integration agent should create basic demo if missing
   fi
   ```
3. Final validation:
   - Verify all effort features are present
   - Confirm no effort was missed
   - Check combined size is reasonable
   - ✅ Verify all demos executed (R291)
   - ✅ Capture demo evidence for review

## Risk Assessment
- **Low Risk:** Sequential splits should merge cleanly
- **Medium Risk:** Controller modifications may conflict
- **Mitigation:** Test after each merge to catch issues early

## Integration Agent Instructions
1. CD to integration directory before starting
2. Execute merges in the EXACT order specified
3. Run tests after EACH merge
4. Document any conflicts encountered
5. Create work-log.md with all operations
6. Generate INTEGRATE_WAVE_EFFORTS-REPORT.md when complete
```

## Validation Before Completion

```bash
#!/bin/bash
# Validate the merge plan is complete

validate_merge_plan() {
    local plan_file="$1"
    
    # Check all required sections exist
    for section in "Target Integration Branch" "Branches to Merge" "Merge Strategy" "Validation Steps"; do
        if ! grep -q "## $section" "$plan_file"; then
            echo "❌ Missing required section: $section"
            return 1
        fi
    done
    
    # Verify no integration branches as sources
    if grep -E "origin/.*-integration" "$plan_file" | grep -v "^##"; then
        echo "❌ ERROR: Integration branches found as merge sources!"
        return 1
    fi
    
    # Check merge commands are present
    merge_count=$(grep -c "git merge origin/" "$plan_file")
    if [[ $merge_count -lt 1 ]]; then
        echo "❌ No merge commands found in plan"
        return 1
    fi
    
    echo "✅ Merge plan validation passed"
    echo "📊 Total merges planned: $merge_count"
    return 0
}

# Run validation
validate_merge_plan "WAVE-MERGE-PLAN.md"
```

## State Transitions

From WAVE_MERGE_PLANNING state:
- **PLAN_COMPLETE** → Return to orchestrator
- **VALIDATION_FAILED** → Fix and re-validate

## Critical Success Criteria

1. ✅ CD'd to integration directory FIRST (verified with pwd)
2. ✅ WAVE-MERGE-PLAN.md created IN the integration directory (not root!)
3. ✅ All effort branches analyzed and categorized
4. ✅ Merge order determined based on dependencies
5. ✅ NO integration branches used as sources
6. ✅ Split branches handled correctly
7. ✅ Clear instructions for Integration Agent
8. ✅ Validation steps included
9. ✅ Risk assessment documented

## Common Mistakes to Avoid

1. **Using integration branches as sources**
   - ❌ WRONG: Merge from wave1-integration
   - ✅ RIGHT: Merge from effort branches only

2. **Including "too large" original branches**
   - ❌ WRONG: Include both original and splits
   - ✅ RIGHT: Exclude original, include only splits

3. **Wrong merge order**
   - ❌ WRONG: Random order
   - ✅ RIGHT: Dependency-aware ordering

4. **Executing merges**
   - ❌ WRONG: Actually running git merge
   - ✅ RIGHT: Only documenting what to merge

5. **Missing validation steps**
   - ❌ WRONG: No test commands
   - ✅ RIGHT: Test after each merge

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

CRITICAL: Only original effort branches in merge plans. Integration branches are TARGETS not SOURCES. This prevents recursive integration issues.
---

### 🔴🔴🔴 RULE R340: Report Merge Plan Location to Orchestrator 🔴🔴🔴

**File**: `$CLAUDE_PROJECT_DIR/rule-library/R340-planning-file-metadata-tracking.md`
**Criticality**: BLOCKING - All planning files must be tracked

**AFTER CREATING WAVE-MERGE-PLAN.md, YOU MUST REPORT:**

```markdown
## 📋 PLANNING FILE CREATED

**Type**: merge_plan (wave)
**Path**: /efforts/phase{X}/wave{Y}/integration-workspace/WAVE-MERGE-PLAN.md
**Phase**: {X}
**Wave**: {Y}
**Target Branch**: phase{X}-wave{Y}-integration
**Efforts Count**: {N}
**Created At**: {ISO-8601-timestamp}

ORCHESTRATOR: Please update planning_files.merge_plans.wave["phase{X}_wave{Y}"] in state file per R340
```

**EXAMPLE REPORT:**
```markdown
## 📋 PLANNING FILE CREATED

**Type**: merge_plan (wave)
**Path**: /efforts/phase1/wave2/integration-workspace/WAVE-MERGE-PLAN.md
**Phase**: 1
**Wave**: 2
**Target Branch**: phase1-wave2-integration
**Efforts Count**: 3
**Created At**: 2025-01-20T12:00:00Z

ORCHESTRATOR: Please update planning_files.merge_plans.wave["phase1_wave2"] in state file per R340
```

**R340 AUTOMATION: Update Orchestrator State Directly**

After creating the merge plan, Code Reviewer MUST update orchestrator-state-v3.json with R340 metadata:

```bash
# R340: Record merge plan location in orchestrator state
PHASE=$(jq -r '.state_machine.current_phase' orchestrator-state-v3.json)
WAVE=$(jq -r '.state_machine.current_wave' orchestrator-state-v3.json)
WAVE_ID="phase${PHASE}_wave${WAVE}"
TIMESTAMP=$(date -u +%Y-%m-%dT%H:%M:%SZ)

# Get absolute path to merge plan
MERGE_PLAN_ABS=$(realpath "$MERGE_PLAN")

# Update orchestrator state with R340 metadata
jq --arg wave_id "$WAVE_ID" \
   --arg path "$MERGE_PLAN_ABS" \
   --arg timestamp "$TIMESTAMP" \
   --arg creator "code-reviewer" \
   --argjson phase "$PHASE" \
   --argjson wave "$WAVE" \
   '
   # Initialize planning_repo_files if not exists
   .planning_repo_files //= {} |
   .planning_repo_files.merge_plans //= {} |
   .planning_repo_files.merge_plans.wave //= {} |

   # Record merge plan metadata
   .planning_repo_files.merge_plans.wave[$wave_id] = {
     "file_path": $path,
     "created_at": $timestamp,
     "created_by": $creator,
     "phase": $phase,
     "wave": $wave,
     "status": "active",
     "deprecated": false
   }
   ' orchestrator-state-v3.json > orchestrator-state-v3.json.tmp

mv orchestrator-state-v3.json.tmp orchestrator-state-v3.json

echo "✅ R340: Merge plan tracked in orchestrator-state-v3.json"
echo "   Wave ID: $WAVE_ID"
echo "   Path: $MERGE_PLAN_ABS"

# Commit the state update
git add orchestrator-state-v3.json
git commit -m "feat: Track wave ${WAVE} merge plan per R340

R340 metadata recorded in planning_repo_files.merge_plans.wave

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
3. **EXACTLY AS SHOWN**: Use exact format - no variations like CONTINUE-ORCHESTRATING
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

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

