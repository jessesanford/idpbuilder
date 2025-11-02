# Code Reviewer - EFFORT_PLAN_CREATION State Rules

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
Create comprehensive implementation plans for individual efforts that will be assigned to SW Engineers. Plans MUST include thorough cross-effort awareness analysis to prevent integration failures.

## 🚨🚨🚨 MANDATORY PRE-FLIGHT: R420 Research Protocol 🚨🚨🚨

**BEFORE creating ANY effort plan, you MUST execute R420 Cross-Effort Planning Awareness Protocol:**

### Step 1: Discover Previous Work
```bash
# Read orchestrator state
ORCHESTRATOR_STATE="/path/to/orchestrator-state-v3.json"

# Discover all previous efforts
PREVIOUS_EFFORTS=$(jq -r '.efforts_completed[] | .effort_id' $ORCHESTRATOR_STATE)

echo "📋 R420 Discovery Phase:"
echo "Previous efforts to review: $PREVIOUS_EFFORTS"

# Count efforts to review
EFFORT_COUNT=$(echo "$PREVIOUS_EFFORTS" | wc -l)
if [ $EFFORT_COUNT -gt 0 ]; then
    echo "MANDATORY: Must review $EFFORT_COUNT previous effort(s)"
else
    echo "ℹ️ No previous efforts (first effort in wave/phase)"
fi
```

### Step 2: Read ALL Previous Implementations
```bash
# For EACH previous effort
echo "$PREVIOUS_EFFORTS" | while read effort_id; do
    echo "📖 Reading effort: $effort_id"

    # Get effort directory
    effort_dir=$(jq -r ".efforts_completed[] | select(.effort_id == \"$effort_id\") | .effort_directory" $ORCHESTRATOR_STATE)

    if [ -d "$effort_dir" ]; then
        echo "  Analyzing $effort_dir..."

        # Find exported functions
        echo "  → Exported functions:"
        grep -r "^func [A-Z]" --include="*.go" "$effort_dir" 2>/dev/null | head -5

        # Find interfaces
        echo "  → Interfaces:"
        grep -r "^type.*interface" --include="*.go" "$effort_dir" 2>/dev/null

        # Find structs
        echo "  → Structs:"
        grep -r "^type [A-Z].*struct" --include="*.go" "$effort_dir" 2>/dev/null | head -5

        # Find file structure
        echo "  → File structure:"
        find "$effort_dir" -name "*.go" -type f 2>/dev/null | head -10
    else
        echo "  ⚠️ Directory not found: $effort_dir"
    fi
done
```

### Step 3: Read Previous Plans (for unimplemented efforts)
```bash
# Find all implementation plans
PLAN_DIR="/path/to/.software-factory"
PREVIOUS_PLANS=$(find "$PLAN_DIR" -name "*IMPLEMENTATION-PLAN*.md" 2>/dev/null)

echo "📖 Reading previous plans..."
echo "$PREVIOUS_PLANS" | while read plan_file; do
    echo "  Plan: $plan_file"

    # Extract planned interfaces
    echo "  → Planned interfaces:"
    grep -A3 "interface" "$plan_file" 2>/dev/null | head -10

    # Extract planned packages
    echo "  → Planned packages:"
    grep "pkg/" "$plan_file" 2>/dev/null | head -5
done
```

### Step 4: Analyze for Conflicts
```bash
echo "🔍 R420 Conflict Analysis:"

# Check for duplicate file paths
echo "  Checking file structure conflicts..."
# (Compare planned files against existing files)

# Check for API mismatches
echo "  Checking API compatibility..."
# (Verify any assumed APIs actually exist)

# Check for method visibility errors
echo "  Checking method visibility..."
# (Ensure only exported methods are accessed)
```

## 🚨🚨🚨 MANDATORY PLAN SECTION: Prior Work Analysis

Every effort plan MUST include this section:

```markdown
## 🔍 PRIOR WORK ANALYSIS (R420 MANDATORY)

### Discovery Phase Results
- **Previous Efforts Reviewed**: [List effort IDs]
- **Previous Plans Reviewed**: [List plan files]
- **Research Timestamp**: [ISO timestamp]
- **Research Status**: [COMPLETE/INCOMPLETE]

### File Structure Findings
| File Path | Source Effort | Status | Action Required |
|-----------|---------------|--------|-----------------|
| [path] | [effort ID] | EXISTS/NEW | MUST (NOT) create |

### Interface/API Findings
| Interface/API | Source | Signature | Action Required |
|---------------|--------|-----------|-----------------|
| [name] | [effort ID] | [signature] | MUST implement/use |

### Type/Struct Findings
| Type | Source | Exported | Action Required |
|------|--------|----------|-----------------|
| [name] | [effort ID] | YES/NO | Can/Cannot use |

### Method Visibility Findings
| Method | Type | Visibility | Can Access? | Action Required |
|--------|------|------------|-------------|-----------------|
| [method] | [type] | EXPORTED/UNEXPORTED | YES/NO | [action] |

### Conflicts Detected
- ✅ NO duplicate file paths detected
- ✅ NO API mismatches detected
- ✅ NO method visibility violations detected

OR:

- ❌ [Describe detected conflict]
- ❌ [Describe detected conflict]

### Required Integrations
1. MUST [integration requirement]
2. MUST [integration requirement]

### Forbidden Actions
- ❌ DO NOT [forbidden action with reason]
- ❌ DO NOT [forbidden action with reason]
```

## 🚨🚨🚨 BLOCKING Validation Gate

Before plan can be approved, it MUST pass R420 validation:

```bash
# Run validation script
bash /path/to/tools/validate-R420-compliance.sh <plan-file>

# If validation fails, plan CANNOT be approved
# Must return to research phase and complete R420 protocol
```

## State-Specific Requirements

### 1. Load Current Effort Context
```bash
# From orchestrator state
CURRENT_EFFORT=$(jq -r '.efforts_in_progress[0]' orchestrator-state-v3.json)
EFFORT_ID=$(echo "$CURRENT_EFFORT" | jq -r '.effort_id')
EFFORT_DESC=$(echo "$CURRENT_EFFORT" | jq -r '.description')

echo "Creating plan for: $EFFORT_ID"
echo "Description: $EFFORT_DESC"
```

### 2. Load Wave/Phase Context
```bash
CURRENT_PHASE=$(jq -r '.current_phase' orchestrator-state-v3.json)
CURRENT_WAVE=$(jq -r '.current_wave' orchestrator-state-v3.json)

echo "Context: Phase $CURRENT_PHASE, Wave $CURRENT_WAVE"

# Read phase plan for context
PHASE_PLAN="/path/to/.software-factory/PHASE${CURRENT_PHASE}-PLAN.md"
if [ -f "$PHASE_PLAN" ]; then
    echo "Reading phase plan: $PHASE_PLAN"
    # Extract relevant requirements for this effort
fi

# Read wave plan for context
WAVE_PLAN="/path/to/.software-factory/PHASE${CURRENT_PHASE}-WAVE${CURRENT_WAVE}-PLAN.md"
if [ -f "$WAVE_PLAN" ]; then
    echo "Reading wave plan: $WAVE_PLAN"
    # Extract wave-specific requirements for this effort
fi
```

### 3. Execute R420 Research (MANDATORY)
See "MANDATORY PRE-FLIGHT: R420 Research Protocol" above.

**BLOCKING:** Cannot proceed without completing R420.

### 4. Create Implementation Plan
Using template: `/path/to/templates/IMPLEMENTATION-PLAN-TEMPLATE.md`

**MUST include:**
- Standard sections from template
- **R420 Prior Work Analysis section** (MANDATORY)
- Effort-specific requirements
- Size estimates (line count targets)
- Demo requirements (per R330)
- TDD test pseudo-code (per R400/R401)

### 5. Validate Plan
```bash
# Run R420 validation
bash tools/validate-R420-compliance.sh <plan-file> || exit 1

# Run general plan validation
bash tools/validate-implementation-plan.sh <plan-file> || exit 1

# Check size estimates (R535 - Code Reviewer enforcement at 900 lines)
ESTIMATED_LINES=$(grep "Estimated Lines" <plan-file> | grep -o "[0-9]*")
if [ $ESTIMATED_LINES -gt 900 ]; then
    echo "⚠️ WARNING: Estimated $ESTIMATED_LINES lines exceeds 900 line enforcement threshold (R535)"
    echo "   Consider splitting effort during planning phase"
elif [ $ESTIMATED_LINES -gt 800 ]; then
    echo "ℹ️ Note: $ESTIMATED_LINES lines in grace buffer (800-900) per R535"
    echo "   SW Engineers see 800 limit, Code Reviewers enforce at 900"
fi
```

### 6. Save Plan to Metadata Location
```bash
# Save to .software-factory directory
METADATA_DIR="/path/to/effort-dir/.software-factory"
mkdir -p "$METADATA_DIR"

PLAN_FILE="$METADATA_DIR/${EFFORT_ID}-IMPLEMENTATION-PLAN.md"
cp <plan-file> "$PLAN_FILE"

echo "Plan saved to: $PLAN_FILE"

# R340: Track effort plan in orchestrator state
EFFORT_ID=$(basename "$METADATA_DIR")
TIMESTAMP=$(date -u +%Y-%m-%dT%H:%M:%SZ)

jq \
   --arg effort_id "$EFFORT_ID" \
   --arg path "$PLAN_FILE" \
   --arg timestamp "$TIMESTAMP" \
   --arg creator "code-reviewer" \
   --arg branch "$EFFORT_BRANCH" \
   --argjson phase "$PHASE" \
   --argjson wave "$WAVE" \
   '
   # Initialize effort_repo_files if not exists
   .effort_repo_files //= {} |
   .effort_repo_files.effort_plans //= {} |

   # Record effort plan metadata
   .effort_repo_files.effort_plans[$effort_id] = {
     "file_path": $path,
     "created_at": $timestamp,
     "created_by": $creator,
     "target_branch": $branch,
     "effort_name": $effort_id,
     "phase": $phase,
     "wave": $wave,
     "status": "active",
     "deprecated": false
   }
   ' orchestrator-state-v3.json > orchestrator-state-v3.json.tmp

mv orchestrator-state-v3.json.tmp orchestrator-state-v3.json

echo "✅ R340: Effort plan tracked in orchestrator-state-v3.json"
echo "   Effort: $EFFORT_ID"
echo "   Path: $PLAN_FILE"

# Commit the state update
git add orchestrator-state-v3.json
git commit -m "feat: Track effort $EFFORT_ID plan per R340

- Recorded plan at: $PLAN_FILE
- Effort: $EFFORT_ID
- Phase: $PHASE, Wave: $WAVE
- R340 metadata tracking complete

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>"
git push
```

## Critical Rules for This State

### 🚨🚨🚨 R420 - Cross-Effort Planning Awareness Protocol (BLOCKING)
**Source:** rule-library/R420-cross-effort-planning-awareness-protocol.md

**MANDATORY in EFFORT_PLAN_CREATION:**
- Execute full R420 research protocol BEFORE planning
- Include R420 Prior Work Analysis section in plan
- Pass R420 validation before plan approval
- Document ALL conflicts found
- Define required integrations and forbidden actions

**Grading Penalty:** -50% for incomplete research, -100% if leads to duplicates

---

### 🚨🚨 R373 - Mandatory Code Reuse and Interface Compliance (BLOCKING)
**Source:** rule-library/R373-mandatory-code-reuse-and-interface-compliance.md

**MANDATORY in EFFORT_PLAN_CREATION:**
- R420 research discovers what interfaces/APIs exist
- Plan MUST specify using existing interfaces (not creating new ones)
- Plan MUST forbid duplicate implementations
- Plan MUST require exact signature implementation

**Grading Penalty:** -100% for creating duplicate interfaces

---

### 🚨🚨 R374 - Pre-Planning Research Protocol (BLOCKING)
**Source:** rule-library/R374-pre-planning-research-protocol.md

**MANDATORY in EFFORT_PLAN_CREATION:**
- Research ALL existing implementations across effort repository
- Document findings in plan
- List all code to reuse
- List forbidden duplications

**Grading Penalty:** -50% for no research, -30% for incomplete research

---

### 🚨 R330 - Demo Planning Requirements (BLOCKING)
**Source:** rule-library/R330-demo-planning-requirements.md

**MANDATORY in EFFORT_PLAN_CREATION:**
- Plan MUST include explicit demo requirements
- Demo MUST prove effort works as designed
- Demo implementation MUST be included in size estimates

**Grading Penalty:** -25% to -50% for missing demo requirements

---

### 🚨 R400/R401 - TDD Requirements (BLOCKING)
**Source:** rule-library/R400-tdd-mandatory.md, R401-tests-first-enforcement.md

**MANDATORY in EFFORT_PLAN_CREATION:**
- Plan MUST include test pseudo-code BEFORE implementation details
- Tests MUST be designed to fail first (RED phase)
- Implementation MUST be minimal code to pass tests (GREEN phase)
- Coverage requirements MUST be specified

**Grading Penalty:** -100% for planning implementation before tests

---

### ⚠️ R007 - Size Limit Compliance
**Source:** rule-library/R007-size-limit-compliance.md

**MANDATORY in EFFORT_PLAN_CREATION:**
- Estimate total line count for effort
- If estimate > 700 lines, consider pre-splitting during planning
- If estimate > 900 lines, MUST pre-split (R535 enforcement threshold)
- Document size estimates in plan

**Grading Penalty:** Exceeding 900 lines without split = automatic split required (R535)

## State Transition Criteria

### Can Transition to COMPLETED When:
1. ✅ R420 research protocol executed completely
2. ✅ R420 Prior Work Analysis section included in plan
3. ✅ Plan passes R420 validation script
4. ✅ Plan includes all required sections
5. ✅ Plan includes TDD test pseudo-code
6. ✅ Plan includes demo requirements
7. ✅ Plan saved to metadata location
8. ✅ Orchestrator state updated with plan location

### Blocked If:
- ❌ R420 research not completed
- ❌ R420 validation fails
- ❌ Conflicts detected but not resolved
- ❌ Plan missing required sections
- ❌ Size estimates exceed 900 lines without split consideration (R535 enforcement)

## Grading Criteria for This State

### Research Quality (30%)
- ✅ All previous efforts reviewed
- ✅ All previous plans reviewed
- ✅ All findings documented in tables
- ✅ Conflicts analyzed and documented
- ✅ Required integrations identified

### Plan Completeness (30%)
- ✅ All template sections included
- ✅ R420 Prior Work Analysis section complete
- ✅ TDD test pseudo-code included
- ✅ Demo requirements specified
- ✅ Size estimates provided

### Compliance Verification (20%)
- ✅ R420 validation passes
- ✅ No forbidden duplications planned
- ✅ Existing interfaces to be implemented
- ✅ Required integrations specified

### Documentation Quality (20%)
- ✅ Clear, actionable instructions for SW Engineer
- ✅ Sufficient detail for implementation
- ✅ Conflicts and forbidden actions clearly stated
- ✅ Integration requirements explicit

### Failure Conditions:
- **R420 research skipped:** -50%
- **R420 validation fails:** -40%
- **Planning duplicates:** -100% (via R373)
- **Missing TDD tests:** -100% (via R400)
- **Missing demo:** -50% (via R330)

## Exit Actions

### Before Transitioning to COMPLETED:
1. Save TODO state (R287)
2. Commit plan file to git
3. Push changes to remote
4. Update orchestrator state with plan metadata
5. Output continuation flag: `CONTINUE-SOFTWARE-FACTORY=TRUE`

## Common Pitfalls to Avoid

### ❌ WRONG: Skipping R420 Research
```markdown
## Implementation Plan
[Creates plan without reviewing previous efforts]
[Assumes APIs exist without verification]
[Plans file creation without checking for conflicts]
```

**Result:** Integration failures, duplicate declarations, API mismatches

### ✅ CORRECT: Full R420 Execution
```markdown
## 🔍 PRIOR WORK ANALYSIS (R420 MANDATORY)

### Discovery Phase Results
- Previous Efforts Reviewed: E1.1.0, E1.2.0
- Research Timestamp: 2025-10-05T14:30:00Z

### File Structure Findings
| File Path | Source | Status | Action |
|-----------|--------|--------|--------|
| pkg/cmd/push/root.go | E1.2.0 | EXISTS | MUST NOT duplicate |

### Interface/API Findings
| Interface | Source | Signature | Action |
|-----------|--------|-----------|--------|
| Registry | E1.1.0 | Push(ctx, image, content) | MUST implement exactly |

[Rest of plan follows...]
```

**Result:** Conflicts prevented, integration succeeds

---

## Summary Checklist

Before completing EFFORT_PLAN_CREATION state:

- [ ] R420 discovery phase executed
- [ ] All previous efforts reviewed
- [ ] All previous plans reviewed
- [ ] Conflict analysis completed
- [ ] R420 Prior Work Analysis section in plan
- [ ] File structure conflicts documented
- [ ] Interface/API findings documented
- [ ] Method visibility checked
- [ ] Required integrations listed
- [ ] Forbidden actions listed
- [ ] R420 validation script passes
- [ ] TDD test pseudo-code included
- [ ] Demo requirements specified
- [ ] Size estimates provided
- [ ] Plan saved to metadata location
- [ ] Orchestrator state updated
- [ ] Changes committed and pushed

**Only proceed when ALL boxes checked!**
