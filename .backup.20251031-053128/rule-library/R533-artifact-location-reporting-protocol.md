# 🚨🚨🚨 RULE R533 - ARTIFACT LOCATION REPORTING PROTOCOL (BLOCKING)

## Metadata
- **Rule ID**: R533
- **Title**: Artifact Location Reporting Protocol
- **Category**: State Management / Artifact Tracking
- **Criticality**: BLOCKING - Violation = -20% per untracked artifact
- **Applies To**: Code-Reviewer, Architect (all artifact-creating agents)
- **Created**: 2025-10-29
- **Related Rules**: R340 (Planning File Metadata Tracking), R288 (State File Update Protocol), R510 (State Execution Checklist Compliance)

## Problem Statement

**Discovered Issue**: Code-Reviewer creates artifacts (fix plans, code reviews, split plans) but does NOT report their locations to Orchestrator. Result:

```
Code-Reviewer creates: .software-factory/phase1/FIX-PLAN-WAVE-1--20251029.md
Orchestrator searches: jq '.artifacts.fix_plans' orchestrator-state-v3.json
Result: null (artifact not tracked)
Fallback: Filesystem search (R340 VIOLATION!)
```

**Root Cause**: No rule requiring agents to record artifact metadata in state file.

## Rule Statement

**EVERY artifact created by ANY agent MUST be recorded in `orchestrator-state-v3.json` with complete metadata BEFORE the agent completes its current state.**

## Scope

### Artifacts Covered
1. **Fix Plans** (Code-Reviewer → CREATE_FIX_PLAN)
2. **Code Review Reports** (Code-Reviewer → PERFORM_CODE_REVIEW)
3. **Split Plans** (Code-Reviewer → CREATE_SPLIT_PLAN)
4. **Integration Reports** (Orchestrator → Integration States)
5. **Architecture Reviews** (Architect → Architecture Review States)

### Agents Affected
- Code-Reviewer (primary)
- Architect (architecture reviews)
- Orchestrator (integration reports)

## Requirements

### 1. Artifact Creation Pattern

When creating ANY artifact file:

```bash
# Step 1: Create the artifact file
ARTIFACT_FILE=".software-factory/phase1/FIX-PLAN-WAVE-1--20251029-120000.md"
cat > "$ARTIFACT_FILE" << 'EOF'
# Fix Plan for Wave 1
[... artifact content ...]
EOF

# Step 2: IMMEDIATELY record in state file (R533 REQUIREMENT)
ARTIFACT_ID="wave1-fix-001"  # Unique identifier
TIMESTAMP=$(date -u +%Y-%m-%dT%H:%M:%SZ)

jq --arg id "$ARTIFACT_ID" \
   --arg path "$ARTIFACT_FILE" \
   --arg timestamp "$TIMESTAMP" \
   --arg type "fix_plan" \
   --arg scope "wave" \
   --argjson bugs '["BUG-001-EXAMPLE", "BUG-002-EXAMPLE"]' \
   '.artifacts.fix_plans[$id] = {
     "file_path": $path,
     "created_at": $timestamp,
     "created_by": "code-reviewer",
     "artifact_type": $type,
     "scope": $scope,
     "status": "READY",
     "related_bugs": $bugs,
     "metadata": {}
   }' orchestrator-state-v3.json > tmp && mv tmp orchestrator-state-v3.json

# Step 3: Verify artifact tracked
jq -e ".artifacts.fix_plans[\"$ARTIFACT_ID\"]" orchestrator-state-v3.json > /dev/null || {
    echo "❌ CRITICAL: Artifact not tracked in state file (R533 violation)"
    exit 533
}

echo "✅ Artifact tracked: $ARTIFACT_ID at $ARTIFACT_FILE"
```

### 2. Required Metadata Fields

ALL artifacts MUST include these fields:

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| `file_path` | string | Path to artifact file (absolute or relative from project root) | ✅ |
| `created_at` | string (ISO 8601) | Timestamp when artifact was created | ✅ |
| `created_by` | enum | Agent that created artifact (code-reviewer/architect/orchestrator) | ✅ |
| `artifact_type` | enum | Type of artifact (fix_plan/code_review/split_plan/integration_report/architecture_review) | ✅ |
| `scope` | enum | Scope of artifact (effort/wave/phase/project) | ✅ |
| `status` | enum | Current status (CREATED/READY/IN_USE/COMPLETED/ARCHIVED) | ✅ |
| `related_bugs` | array\[string\] | Bug IDs this artifact addresses (empty array if none) | Optional |
| `metadata` | object | Additional artifact-specific metadata | Optional |

### 3. Artifact Tracking Location

Artifacts are tracked in `orchestrator-state-v3.json` under `.artifacts` section:

```json
{
  "state_machine": {...},
  "project_progression": {...},
  "references": {...},
  "artifacts": {
    "fix_plans": {
      "wave1-fix-001": {
        "file_path": ".software-factory/phase1/FIX-PLAN-WAVE-1--20251029.md",
        "created_at": "2025-10-29T12:00:00Z",
        "created_by": "code-reviewer",
        "artifact_type": "fix_plan",
        "scope": "wave",
        "status": "READY",
        "related_bugs": ["BUG-001-EXAMPLE"],
        "metadata": {}
      }
    },
    "code_reviews": {
      "wave1-effort1-review-001": {
        "file_path": "efforts/phase1/wave1/effort1/.software-factory/CODE-REVIEW-REPORT--20251029.md",
        "created_at": "2025-10-29T11:00:00Z",
        "created_by": "code-reviewer",
        "artifact_type": "code_review",
        "scope": "effort",
        "status": "COMPLETED",
        "related_bugs": [],
        "metadata": {"issues_found": 3, "severity": "MEDIUM"}
      }
    },
    "split_plans": {...},
    "integration_reports": {...},
    "architecture_reviews": {...}
  }
}
```

### 4. Artifact ID Naming Convention

Artifact IDs MUST follow this pattern:

```
{scope}-{descriptor}-{sequence}

Examples:
- wave1-fix-001 (first fix plan for wave 1)
- phase1-integration-001 (first integration report for phase 1)
- effort1-review-001 (first code review for effort 1)
- project-architecture-001 (first project-level architecture review)
```

## Enforcement

### State Checklist Requirement

ALL artifact-creating states MUST include this checklist item:

```markdown
- [ ] X. Record artifact location in orchestrator-state-v3.json per R533
  - Action: Update `.artifacts.[artifact_type]` with complete metadata
  - Required Fields: file_path, created_at, created_by, artifact_type, scope, status
  - Validation: `jq '.artifacts.[artifact_type].[artifact_id]' orchestrator-state-v3.json` returns metadata
  - **BLOCKING**: Orchestrator cannot discover artifacts without this (R340 compliance)
```

### Pre-Completion Validation

Before completing state, agent MUST verify:

```bash
# Verify all created artifacts are tracked
ARTIFACT_COUNT=$(find .software-factory -name "FIX-PLAN-*.md" | wc -l)
TRACKED_COUNT=$(jq '.artifacts.fix_plans | length // 0' orchestrator-state-v3.json)

if [ "$ARTIFACT_COUNT" -ne "$TRACKED_COUNT" ]; then
    echo "❌ CRITICAL: $ARTIFACT_COUNT artifacts created but only $TRACKED_COUNT tracked"
    echo "This is an R533 violation - all artifacts MUST be tracked"
    exit 533
fi

echo "✅ All $ARTIFACT_COUNT artifacts properly tracked in state file"
```

### Grading Impact

- **Missing metadata** = -20% per untracked artifact
- **Incorrect metadata** = -10% per artifact
- **Orchestrator cannot find artifacts** = SYSTEM FAILURE (-100%)

## Integration with Existing Rules

### R340: Planning File Metadata Tracking

R533 **implements** R340's requirement:

> **MUST** read artifact locations from orchestrator-state-v3.json
> **NEVER** search directories for artifacts
> **ALWAYS** use state file metadata sections

**Before R533**: R340 referenced non-existent `.effort_repo_files.fix_plans` field
**After R533**: R340 compliant using `.artifacts.fix_plans` field

### R288: State File Update Protocol

R533 **extends** R288 to include artifact metadata:

- Update state file with artifact metadata
- Validate state file after update
- Commit state file changes

### R510: State Execution Checklist Compliance

R533 **requires** checklist item in all artifact-creating states:

- CREATE_FIX_PLAN → artifact reporting checklist item
- PERFORM_CODE_REVIEW → artifact reporting checklist item
- CREATE_SPLIT_PLAN → artifact reporting checklist item

## Examples

### Example 1: Code-Reviewer Creating Fix Plan

```bash
# In CREATE_FIX_PLAN state
PHASE=1
WAVE=1
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
FIX_PLAN_FILE=".software-factory/phase${PHASE}/FIX-PLAN-WAVE-${WAVE}--${TIMESTAMP}.md"

# Create fix plan
cat > "$FIX_PLAN_FILE" << 'EOF'
# Fix Plan for Wave 1
## Bugs to Fix
- BUG-001: Duplicate push command
- BUG-002: Undefined variable
[... detailed fix instructions ...]
EOF

# R533: Record in state file
ARTIFACT_ID="wave${WAVE}-fix-001"
ISO_TIMESTAMP=$(date -u +%Y-%m-%dT%H:%M:%SZ)

jq --arg id "$ARTIFACT_ID" \
   --arg path "$FIX_PLAN_FILE" \
   --arg timestamp "$ISO_TIMESTAMP" \
   --argjson bugs '["BUG-001-DUPLICATE-PUSHCMD", "BUG-002-UNDEFINED-VAR"]' \
   '.artifacts.fix_plans[$id] = {
     "file_path": $path,
     "created_at": $timestamp,
     "created_by": "code-reviewer",
     "artifact_type": "fix_plan",
     "scope": "wave",
     "status": "READY",
     "related_bugs": $bugs,
     "metadata": {
       "bug_count": 2,
       "complexity": "MEDIUM"
     }
   }' orchestrator-state-v3.json > tmp && mv tmp orchestrator-state-v3.json

# Verify
jq -e ".artifacts.fix_plans[\"$ARTIFACT_ID\"]" orchestrator-state-v3.json
echo "✅ Fix plan artifact tracked per R533"
```

### Example 2: Orchestrator Waiting for Fix Plans

```bash
# In WAITING_FOR_FIX_PLANS state
# R533 + R340 compliant: Read from state file (NO filesystem search)

FIX_PLAN_COUNT=$(jq '.artifacts.fix_plans | length // 0' orchestrator-state-v3.json)

if [ "$FIX_PLAN_COUNT" -gt 0 ]; then
    echo "✅ Found $FIX_PLAN_COUNT fix plans tracked in state (R533/R340 compliant)"

    # Verify all tracked fix plans exist
    ALL_EXIST=true
    jq -r '.artifacts.fix_plans[] | @json' orchestrator-state-v3.json | while IFS= read -r artifact; do
        FILE_PATH=$(echo "$artifact" | jq -r '.file_path')
        ARTIFACT_ID=$(echo "$artifact" | jq -r '.artifact_id')

        if [ -f "$FILE_PATH" ]; then
            echo "✅ Fix plan exists: $FILE_PATH"
        else
            echo "❌ CRITICAL: Tracked plan missing: $FILE_PATH"
            ALL_EXIST=false
        fi
    done

    [ "$ALL_EXIST" = true ] && echo "✅ All fix plans verified - ready to proceed"
else
    echo "⏳ No fix plans tracked yet - waiting for Code-Reviewer to complete per R533"
fi
```

## Migration Notes

### For Existing States

States referencing `.effort_repo_files.fix_plans` (non-existent field) MUST be updated:

**OLD (broken)**:
```bash
FIX_PLANS=$(jq '.effort_repo_files.fix_plans' orchestrator-state-v3.json)
```

**NEW (R533 compliant)**:
```bash
FIX_PLANS=$(jq '.artifacts.fix_plans' orchestrator-state-v3.json)
```

### For State Rules

Add R533 to Primary Directives for all artifact-creating states:

```markdown
### State-Specific Rules:

X. **🚨🚨🚨 R533** - ARTIFACT LOCATION REPORTING PROTOCOL (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R533-artifact-location-reporting-protocol.md`
   - Criticality: BLOCKING
   - Summary: ALL artifacts MUST be tracked in orchestrator-state-v3.json with complete metadata
```

## Testing

### Unit Test: Artifact Tracking

```bash
# Create test artifact
TEST_ARTIFACT=".software-factory/test/FIX-PLAN-TEST--20251029.md"
mkdir -p .software-factory/test
echo "Test fix plan" > "$TEST_ARTIFACT"

# Track per R533
jq '.artifacts.fix_plans["test-fix-001"] = {
  "file_path": "'$TEST_ARTIFACT'",
  "created_at": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'",
  "created_by": "code-reviewer",
  "artifact_type": "fix_plan",
  "scope": "wave",
  "status": "READY",
  "related_bugs": [],
  "metadata": {}
}' orchestrator-state-v3.json > tmp && mv tmp orchestrator-state-v3.json

# Verify
jq -e '.artifacts.fix_plans["test-fix-001"]' orchestrator-state-v3.json && \
echo "✅ Test passed: Artifact tracked per R533"
```

### Integration Test: Runtime Test-03b

Runtime test-03b includes fix plan creation. After deploying R533:

```bash
bash tests/runtime-test-03b-review-fix-path.sh

# Verify artifacts tracked
TEST_STATE="/tmp/test-03b-*/orchestrator-state-v3.json"
ARTIFACT_COUNT=$(jq '.artifacts.fix_plans | length // 0' $TEST_STATE)

[ "$ARTIFACT_COUNT" -gt 0 ] && \
echo "✅ Runtime test passed: Fix plans tracked per R533" || \
echo "❌ Runtime test failed: No artifacts tracked (R533 violation)"
```

## Related Files

### Schema
- `/schemas/orchestrator-state-v3.schema.json` - Defines `.artifacts` structure

### State Rules (require R533 updates)
- `/agent-states/software-factory/code-reviewer/CREATE_FIX_PLAN/rules.md`
- `/agent-states/software-factory/code-reviewer/PERFORM_CODE_REVIEW/rules.md`
- `/agent-states/software-factory/code-reviewer/CREATE_SPLIT_PLAN/rules.md`
- `/agent-states/software-factory/orchestrator/WAITING_FOR_FIX_PLANS/rules.md`
- `/agent-states/software-factory/orchestrator/WAITING_FOR_INTEGRATION_CODE_REVIEW/rules.md`

### Workflows Affected
- Fix plan creation and distribution
- Code review creation and distribution
- Split plan creation and execution
- Integration report generation
- Architecture review flow

## Summary

**Rule**: ALL artifacts MUST be tracked in `orchestrator-state-v3.json` with complete metadata

**Purpose**: Enable artifact discovery via state file (R340 compliant), eliminate filesystem searches

**Enforcement**: -20% per untracked artifact, SYSTEM FAILURE if orchestrator can't find artifacts

**Impact**: Foundation for automated artifact management across entire Software Factory 3.0 system
