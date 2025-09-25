# Software Factory State File Management

## Overview

The Software Factory uses a **dual state file pattern** to maintain clean separation between project-level progress tracking and temporary fix/recovery operations. This pattern, formalized in R375, ensures the main state file remains focused and readable while providing detailed tracking for complex operations.

## The Two-File Pattern

### 1. Main State File: `orchestrator-state.json`

**Purpose**: Track overall project/effort progress through the Software Factory lifecycle

**Lifecycle**: Permanent - exists for entire project duration

**Contains**:
- Current phase/wave/effort status
- State machine position and transitions
- Integration branch tracking
- Effort completion records
- High-level project metrics
- Agent deployment status

**Protected From**:
- Temporary fix details
- Backport/forward-port specifics
- Transient error recovery data
- Fix-specific validation results

### 2. Fix State Files: `orchestrator-[fix-name]-state.json`

**Purpose**: Track specific fix cascades, hotfixes, and recovery operations

**Lifecycle**: Temporary - created at fix start, archived upon completion

**Contains**:
- Fix-specific progress tracking
- Backport/forward-port status per branch
- Detailed validation results
- Fix-specific error logs
- Recovery action history
- Branch-specific metrics

**Examples**:
- `orchestrator-gitea-api-fix-state.json`
- `orchestrator-auth-critical-fix-state.json`
- `orchestrator-pr367-backport-state.json`

## Separation of Concerns

### What Goes in Main State

```json
{
  "current_phase": 2,
  "current_wave": 3,
  "current_state": "ERROR_RECOVERY",
  "error_recovery": {
    "reason": "Critical API fix required",
    "entered_at": "2025-01-21T10:00:00Z",
    "fix_state_file": "orchestrator-api-fix-state.json"
  }
}
```

**Note**: Main state only references the fix state file, not the details

### What Goes in Fix State

```json
{
  "fix_identifier": "api-fix",
  "backports": {
    "release-1.0": {
      "status": "COMPLETED",
      "pr_number": 123,
      "validation": {
        "build": "PASSED",
        "tests": "PASSED"
      }
    }
  },
  "validation_results": {
    "unit_tests": "PASSED",
    "integration_tests": "RUNNING"
  }
}
```

**Note**: All fix-specific details isolated in fix state file

## Benefits of This Pattern

### 1. Clean Main State
- Remains readable and manageable
- Not polluted with temporary details
- Clear project progress visibility
- Easy to understand current position

### 2. Detailed Fix Tracking
- Complete audit trail for fixes
- Branch-specific status tracking
- Validation results per target
- Error and recovery history

### 3. Multiple Concurrent Operations
- Each fix gets its own state file
- No interference between fixes
- Independent lifecycles
- Parallel fix execution support

### 4. Clean Archival
- Fix states archived when complete
- Main state continues unaffected
- Historical record maintained
- Easy cleanup process

## Implementation Examples

### Starting a Fix Cascade

```bash
# 1. Create fix-specific state file
FIX_ID="security-hotfix"
cat > orchestrator-${FIX_ID}-state.json << 'EOF'
{
  "fix_identifier": "security-hotfix",
  "fix_type": "HOTFIX",
  "created_at": "2025-01-21T10:00:00Z",
  "status": "IN_PROGRESS",
  "source_branch": "main",
  "target_branches": ["release-2.0", "release-1.9"]
}
EOF

# 2. Update main state to reference fix
jq '.error_recovery.fix_state_file = "orchestrator-security-hotfix-state.json"' \
   orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json

# 3. Commit both files
git add orchestrator-*.json
git commit -m "state: initiate security-hotfix with separate tracking"
git push
```

### Updating Fix Progress

```bash
# Update ONLY the fix state file
jq '.backports."release-2.0" = {
  "status": "COMPLETED",
  "pr_number": 456,
  "validation": {"build": "PASSED", "tests": "PASSED"}
}' orchestrator-${FIX_ID}-state.json > tmp.json && \
mv tmp.json orchestrator-${FIX_ID}-state.json

git add orchestrator-${FIX_ID}-state.json
git commit -m "fix-state: release-2.0 backport completed"
git push
```

### Completing and Archiving

```bash
# 1. Mark fix as completed
jq '.status = "COMPLETED"' orchestrator-${FIX_ID}-state.json > tmp.json && \
mv tmp.json orchestrator-${FIX_ID}-state.json

# 2. Archive the fix state
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
mkdir -p archived-fixes/$(date +%Y)/$(date +%m)
mv orchestrator-${FIX_ID}-state.json \
   archived-fixes/$(date +%Y)/$(date +%m)/${FIX_ID}-${TIMESTAMP}.json

# 3. Clear reference from main state
jq 'del(.error_recovery.fix_state_file)' orchestrator-state.json > tmp.json && \
mv tmp.json orchestrator-state.json

# 4. Commit everything
git add orchestrator-state.json archived-fixes/
git commit -m "archive: ${FIX_ID} completed and archived"
git push
```

## State File Locations

### Planning Repository (Primary)
```
planning-repo/
├── orchestrator-state.json                    # Main state (permanent)
├── orchestrator-[fix1]-state.json            # Active fix state
├── orchestrator-[fix2]-state.json            # Another active fix
├── archived-fixes/                           # Completed fixes
│   └── 2025/
│       └── 01/
│           ├── api-fix-20250121-143000.json
│           └── auth-fix-20250122-091500.json
└── templates/
    └── fix-state-template.json              # Template for new fixes
```

### Never Store Fix States In:
- Target implementation repositories
- Working copies
- Agent-specific directories
- Temporary locations

## Common Patterns

### Pattern 1: Hotfix Cascade
```
Main State: Tracks ERROR_RECOVERY entry
Fix State: Tracks each branch backport
Result: Clean tracking without pollution
```

### Pattern 2: Multi-Branch Fix
```
Main State: References active fix file
Fix State: Detailed status per branch
Result: Parallel branch operations
```

### Pattern 3: Recovery from Failed Integration
```
Main State: Shows integration failure
Fix State: Tracks resolution attempts
Result: Complete recovery audit trail
```

## Anti-Patterns to Avoid

### ❌ DON'T: Mix fix details in main state
```json
// WRONG - Main state polluted
{
  "current_state": "ERROR_RECOVERY",
  "backports": {  // This belongs in fix state!
    "release-1.0": "completed"
  }
}
```

### ❌ DON'T: Delete fix states
```bash
# WRONG - Destroys audit trail
rm orchestrator-fix-state.json  # Never delete!
```

### ❌ DON'T: Store in implementation repo
```bash
# WRONG - Fix states belong in planning repo
cd implementation-repo
echo "{}" > orchestrator-fix-state.json  # Wrong location!
```

### ✅ DO: Keep clean separation
```json
// Main state - clean reference
{
  "error_recovery": {
    "fix_state_file": "orchestrator-api-fix-state.json"
  }
}

// Fix state - detailed tracking
{
  "backports": {
    "release-1.0": { /* details */ }
  }
}
```

## Monitoring and Reporting

### Check Active Fixes
```bash
# List all active fix states
ls orchestrator-*-state.json | grep -v "orchestrator-state.json"
```

### Generate Fix Report
```bash
# Summary of active fixes
for f in orchestrator-*-fix-state.json; do
  echo "Fix: $(jq -r '.fix_identifier' $f)"
  echo "Status: $(jq -r '.status' $f)"
  echo "Branches: $(jq -r '.target_branches | join(", ")' $f)"
  echo ""
done
```

### Audit Completed Fixes
```bash
# Monthly fix summary
find archived-fixes/2025/01 -name "*.json" -exec \
  jq '{fix: .fix_identifier, type: .fix_type, time: .metrics.total_time_minutes}' {} \;
```

## Rule Compliance

This dual state file pattern is mandated by:
- **R375**: Fix State File Management Protocol
- **R300**: Comprehensive Fix Management Protocol
- **R019**: Error Recovery Protocol

Violations include:
- Mixing fix details in main state: -30% penalty
- Not creating fix state files: -25% penalty
- Deleting instead of archiving: -20% penalty
- Wrong file locations: -15% penalty

## Summary

The dual state file pattern provides:
1. **Clean separation** of concerns
2. **Concurrent operation** support
3. **Complete audit trails** via archival
4. **Protected main state** from pollution
5. **Scalable tracking** for complex operations

Always remember:
- Main state = project journey
- Fix states = temporary operations
- Archive, don't delete
- Commit immediately
- Keep in planning repo

---

*Per R375 - Fix State File Management Protocol*