# 🔴🔴🔴 RULE R344: Metadata Location Tracking (SUPREME LAW) 🔴🔴🔴

## Rule Identifier
**Rule ID**: R344
**Category**: State File Management
**Criticality**: SUPREME LAW - Violations cause wasted time and lost metadata
**Priority**: ABSOLUTE

## Description

ALL metadata file locations MUST be tracked in orchestrator-state-v3.json. Agents MUST report metadata locations immediately after creation. Agents MUST read metadata locations from state file ONLY - NO filesystem searching allowed.

## Rationale

Without centralized metadata tracking:
- Agents waste time searching filesystems
- Metadata gets lost or forgotten
- Recovery after failures is impossible
- Parallel agents cannot coordinate
- State becomes inconsistent with reality

## Requirements

### 1. METADATA TRACKING IN ORCHESTRATOR-STATE.JSON

**ALL metadata locations MUST be in state file:**

```json
{
  "metadata_locations": {
    "implementation_plans": {
      "buildah-builder-interface": {
        "file_path": "/efforts/phase1/wave2/buildah-builder-interface/.software-factory/IMPLEMENTATION-PLAN.md",
        "created_at": "2025-01-20T10:00:00Z",
        "created_by": "code-reviewer",
        "agent_id": "code-reviewer-wave2-001",
        "phase": 1,
        "wave": 2,
        "effort": "buildah-builder-interface"
      }
    },
    "code_review_reports": {
      "buildah-builder-interface-review-001": {
        "file_path": "/efforts/phase1/wave2/buildah-builder-interface/.software-factory/CODE-REVIEW-REPORT.md",
        "created_at": "2025-01-20T12:00:00Z",
        "created_by": "code-reviewer",
        "reviewed_branch": "phase1/wave2/buildah-builder-interface",
        "findings_count": 3
      }
    },
    "integration_reports": {
      "phase1_wave1": {
        "file_path": "/efforts/phase1/wave1/integration-workspace/.software-factory/INTEGRATE_WAVE_EFFORTS-REPORT.md",
        "created_at": "2025-01-20T14:00:00Z",
        "created_by": "integration",
        "integration_branch": "phase1-wave1-integration",
        "status": "PASSED"
      }
    },
    "work_logs": {
      "buildah-builder-interface": {
        "file_path": "/efforts/phase1/wave2/buildah-builder-interface/.software-factory/work-log.md",
        "created_at": "2025-01-20T10:30:00Z",
        "created_by": "sw-engineer",
        "last_updated": "2025-01-20T15:00:00Z"
      }
    },
    "split_plans": {
      "oci-types-split": {
        "file_path": "/efforts/phase1/wave1/oci-types/.software-factory/SPLIT-PLAN.md",
        "created_at": "2025-01-20T11:00:00Z",
        "created_by": "code-reviewer",
        "total_splits": 2,
        "split_branches": [
          "phase1/wave1/oci-types-split-001",
          "phase1/wave1/oci-types-split-002"
        ]
      }
    },
    "fix_plans": {
      "registry-auth-types-fix-001": {
        "file_path": "/efforts/phase1/wave1/registry-auth-types/.software-factory/FIX-PLAN.md",
        "created_at": "2025-01-20T13:00:00Z",
        "created_by": "code-reviewer",
        "fixes_required": 5
      }
    },
    "test_results": {
      "phase1_wave1": {
        "file_path": "/efforts/phase1/wave1/integration-workspace/.software-factory/test-results.md",
        "created_at": "2025-01-20T14:30:00Z",
        "created_by": "integration",
        "tests_passed": 45,
        "tests_failed": 0
      }
    },
    "validation_results": {
      "buildah-builder-interface": {
        "file_path": "/efforts/phase1/wave2/buildah-builder-interface/.software-factory/validation-results.md",
        "created_at": "2025-01-20T15:30:00Z",
        "created_by": "sw-engineer",
        "validation_status": "PASSED"
      }
    },
    "demo_results": {
      "phase1_wave1": {
        "file_path": "/efforts/phase1/wave1/integration-workspace/.software-factory/demo-results.md",
        "created_at": "2025-01-20T16:00:00Z",
        "created_by": "integration",
        "demo_status": "PROJECT_DONEFUL"
      }
    }
  }
}
```

### 2. AGENT RESPONSIBILITIES FOR REPORTING

#### Code Reviewer MUST Report:
```bash
# After creating implementation plan
yq -i ".metadata_locations.implementation_plans.\"$EFFORT_NAME\" = {
  \"file_path\": \"$PLAN_PATH\",
  \"created_at\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
  \"created_by\": \"code-reviewer\",
  \"agent_id\": \"$AGENT_ID\",
  \"phase\": $PHASE,
  \"wave\": $WAVE,
  \"effort\": \"$EFFORT_NAME\"
}" orchestrator-state-v3.json

# After creating review report
yq -i ".metadata_locations.code_review_reports.\"$REVIEW_ID\" = {
  \"file_path\": \"$REPORT_PATH\",
  \"created_at\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
  \"created_by\": \"code-reviewer\",
  \"reviewed_branch\": \"$BRANCH\",
  \"findings_count\": $FINDINGS
}" orchestrator-state-v3.json
```

#### SW Engineer MUST Report:
```bash
# After creating work log
yq -i ".metadata_locations.work_logs.\"$EFFORT_NAME\" = {
  \"file_path\": \"$LOG_PATH\",
  \"created_at\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
  \"created_by\": \"sw-engineer\",
  \"last_updated\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"
}" orchestrator-state-v3.json

# After validation results
yq -i ".metadata_locations.validation_results.\"$EFFORT_NAME\" = {
  \"file_path\": \"$VALIDATION_PATH\",
  \"created_at\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
  \"created_by\": \"sw-engineer\",
  \"validation_status\": \"$STATUS\"
}" orchestrator-state-v3.json
```

#### Integration Agent MUST Report:
```bash
# After creating integration report
yq -i ".metadata_locations.integration_reports.\"${PHASE}_wave${WAVE}\" = {
  \"file_path\": \"$REPORT_PATH\",
  \"created_at\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
  \"created_by\": \"integration\",
  \"integration_branch\": \"$BRANCH\",
  \"status\": \"$STATUS\"
}" orchestrator-state-v3.json

# After test results
yq -i ".metadata_locations.test_results.\"${PHASE}_wave${WAVE}\" = {
  \"file_path\": \"$TEST_PATH\",
  \"created_at\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
  \"created_by\": \"integration\",
  \"tests_passed\": $PASSED,
  \"tests_failed\": $FAILED
}" orchestrator-state-v3.json
```

### 3. AGENT RESPONSIBILITIES FOR READING

#### Orchestrator MUST Read Locations:
```bash
# NEVER search filesystem - read from state
INTEGRATE_WAVE_EFFORTS_REPORT=$(yq ".metadata_locations.integration_reports.\"${PHASE}_wave${WAVE}\".file_path" orchestrator-state-v3.json)

if [ "$INTEGRATE_WAVE_EFFORTS_REPORT" = "null" ]; then
    echo "❌ R344 VIOLATION: Integration report location not in state file"
    exit 1
fi

# Read the report from known location
cat "$INTEGRATE_WAVE_EFFORTS_REPORT"
```

#### All Agents MUST:
```bash
# Read metadata location from state
PLAN_PATH=$(yq ".metadata_locations.implementation_plans.\"$EFFORT_NAME\".file_path" orchestrator-state-v3.json)

# FORBIDDEN - NEVER DO THIS:
# ❌ find /efforts -name "IMPLEMENTATION-PLAN*.md"
# ❌ ls /efforts/*/*/*/*PLAN*.md
# ❌ grep -r "IMPLEMENTATION" /efforts/

# REQUIRED - ALWAYS DO THIS:
# ✅ Read location from orchestrator-state-v3.json
# ✅ Use the exact path provided
# ✅ Report if location is missing
```

### 4. VALIDATION PROTOCOL

```bash
validate_metadata_tracking() {
    local metadata_type="$1"
    local identifier="$2"
    
    # Check if location is tracked
    local file_path=$(yq ".metadata_locations.${metadata_type}.\"${identifier}\".file_path" orchestrator-state-v3.json)
    
    if [ "$file_path" = "null" ] || [ -z "$file_path" ]; then
        echo "❌ R344 VIOLATION: ${metadata_type}/${identifier} not tracked in state"
        return 1
    fi
    
    # Check if file exists at tracked location
    if [ ! -f "$file_path" ]; then
        echo "❌ R344 VIOLATION: Tracked file missing at $file_path"
        return 1
    fi
    
    echo "✅ R344 COMPLIANT: ${metadata_type}/${identifier} properly tracked"
    return 0
}
```

### 5. RECOVERY PROTOCOL

**If metadata location is missing from state:**

```bash
# Step 1: Report violation
echo "🔴 R344 VIOLATION: Metadata location not in state file"

# Step 2: Agent that created it must re-report
if [ -f "$EXPECTED_PATH" ]; then
    # Re-add to state file
    yq -i ".metadata_locations.${TYPE}.\"${ID}\" = {
        \"file_path\": \"$EXPECTED_PATH\",
        \"created_at\": \"$(stat -c %y \"$EXPECTED_PATH\" | cut -d' ' -f1)T$(stat -c %y \"$EXPECTED_PATH\" | cut -d' ' -f2 | cut -d'.' -f1)Z\",
        \"created_by\": \"recovery\",
        \"recovered\": true
    }" orchestrator-state-v3.json
fi

# Step 3: Commit state update
git add orchestrator-state-v3.json
git commit -m "fix: recover metadata location for ${TYPE}/${ID}"
git push
```

## Common Violations

1. ❌ **Searching filesystem for metadata**
   - Wrong: `find /efforts -name "INTEGRATE_WAVE_EFFORTS-REPORT.md"`
   - Right: `yq ".metadata_locations.integration_reports.\"${ID}\".file_path" orchestrator-state-v3.json`

2. ❌ **Creating metadata without reporting location**
   - Wrong: Create file and forget to update state
   - Right: Create file AND immediately update state

3. ❌ **Hardcoding metadata paths**
   - Wrong: `cat /efforts/phase1/wave1/integration-workspace/.software-factory/INTEGRATE_WAVE_EFFORTS-REPORT.md`
   - Right: Read path from state file first

4. ❌ **Delayed metadata reporting**
   - Wrong: Report location after several operations
   - Right: Report IMMEDIATELY after creation

## Enforcement

- **Grading Impact**: -100% IMMEDIATE FAILURE for filesystem searching
- **Performance**: -50% for missing metadata tracking
- **Recovery**: -30% for delayed reporting

## Related Rules

- R343: Metadata directory standardization
- R340: Planning file metadata tracking  
- R288: Mandatory state file updates
- R206: State machine validation

## Implementation Notes

1. Metadata locations are ALWAYS in orchestrator-state-v3.json
2. Report location IMMEDIATELY after creating file
3. Read location BEFORE accessing file
4. NEVER search filesystem for metadata
5. Update state file atomically with git operations

## Key Principle

**"State file knows where everything is - no searching allowed"**

This ensures:
- Instant metadata discovery
- No wasted time searching
- Consistent agent behavior
- Recovery after failures
- Parallel agent coordination

---

**Remember**: If you're searching the filesystem for metadata, you're violating R344. The state file is the ONLY source of truth for metadata locations.