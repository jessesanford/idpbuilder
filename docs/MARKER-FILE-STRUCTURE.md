# Software Factory 2.0 - Marker File Structure Documentation

## Overview

This document describes the marker files and coordination files used in Software Factory 2.0, including their new locations in the `.software-factory/` directory structure and backward compatibility measures.

## File Categories

### 1. State Verification Markers

These temporary markers are created when agents read their state rules (R290 compliance).

**Location**: Project root (unchanged)
**Pattern**: `.state_rules_read_[agent]_[STATE]`
**Lifetime**: Temporary - deleted after state transition
**Purpose**: Verify agent has read state-specific rules before executing state work

Examples:
```
.state_rules_read_orchestrator_MONITOR
.state_rules_read_sw-engineer_IMPLEMENTATION
.state_rules_read_code-reviewer_CODE_REVIEW
.state_rules_read_architect_WAVE_REVIEW
```

### 2. Planning Files (Effort Coordination)

These files coordinate work between agents and persist throughout the effort lifecycle.

#### IMPLEMENTATION-PLAN Files

**Created by**: Code Reviewer in EFFORT_PLANNING state
**New Location**: `.software-factory/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}/plans/IMPLEMENTATION-PLAN-${TIMESTAMP}.md`
**Legacy Location**: `IMPLEMENTATION-PLAN.md` in effort root
**Backward Compatibility**: Symlink from root to new location

Example:
```bash
# New structure
.software-factory/phase1/wave2/effort3-api/plans/IMPLEMENTATION-PLAN-20250120-143000.md

# Backward compatible symlink
IMPLEMENTATION-PLAN.md -> .software-factory/phase1/wave2/effort3-api/plans/IMPLEMENTATION-PLAN-20250120-143000.md
```

#### SPLIT-PLAN Files

**Created by**: Code Reviewer in CREATE_SPLIT_PLAN state
**New Location**: `.software-factory/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}/plans/SPLIT-PLAN-${TIMESTAMP}.md`
**Legacy Location**: `SPLIT-PLAN.md` in effort root
**Backward Compatibility**: Symlink from root to new location

Example:
```bash
# For each split
.software-factory/phase1/wave2/effort3-api-split-001/plans/SPLIT-PLAN-20250120-145500.md
.software-factory/phase1/wave2/effort3-api-split-002/plans/SPLIT-PLAN-20250120-145600.md
```

### 3. Report Files

These files document outcomes and decisions.

#### CODE-REVIEW-REPORT Files

**Created by**: Code Reviewer in CODE_REVIEW state
**New Location**: `.software-factory/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}/reports/CODE-REVIEW-REPORT-${TIMESTAMP}.md`
**Legacy Location**: `CODE-REVIEW-REPORT.md` in effort root
**Backward Compatibility**: Symlink from root to new location

Example:
```bash
# New structure
.software-factory/phase1/wave2/effort3-api/reports/CODE-REVIEW-REPORT-20250120-150000.md

# Backward compatible symlink
CODE-REVIEW-REPORT.md -> .software-factory/phase1/wave2/effort3-api/reports/CODE-REVIEW-REPORT-20250120-150000.md
```

### 4. Work Logs

Track progress and metrics throughout implementation.

**Created by**: Code Reviewer initially, updated by SW Engineer
**New Location**: `.software-factory/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}/logs/work-log-${TIMESTAMP}.md`
**Legacy Location**: `work-log.md` in effort root
**Backward Compatibility**: Symlink from root to new location

## Directory Structure

```
effort-directory/
├── IMPLEMENTATION-PLAN.md (symlink)
├── CODE-REVIEW-REPORT.md (symlink)
├── work-log.md (symlink)
├── .software-factory/
│   └── phase${PHASE}/
│       └── wave${WAVE}/
│           └── ${EFFORT_NAME}/
│               ├── plans/
│               │   ├── IMPLEMENTATION-PLAN-20250120-143000.md
│               │   └── SPLIT-PLAN-20250120-145500.md
│               ├── reports/
│               │   └── CODE-REVIEW-REPORT-20250120-150000.md
│               └── logs/
│                   └── work-log-20250120-143000.md
└── [implementation files...]
```

## Detection Patterns

### Checking for Plan Files (New Pattern)

```bash
# Check new location first, then legacy
EFFORT_NAME="effort3-api"
PHASE=1
WAVE=2

# For IMPLEMENTATION-PLAN
if [ -n "$(ls .software-factory/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}/plans/IMPLEMENTATION-PLAN-*.md 2>/dev/null)" ]; then
    echo "Found plan in new location"
    PLAN=$(ls -t .software-factory/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}/plans/IMPLEMENTATION-PLAN-*.md | head -1)
elif [ -f "IMPLEMENTATION-PLAN.md" ]; then
    echo "Found plan in legacy location"
    PLAN="IMPLEMENTATION-PLAN.md"
else
    echo "No plan found"
fi
```

### Creating Files (New Pattern)

```bash
# Create directory structure
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
PLAN_DIR=".software-factory/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}/plans"
mkdir -p "$PLAN_DIR"

# Create timestamped file
PLAN_FILE="$PLAN_DIR/IMPLEMENTATION-PLAN-${TIMESTAMP}.md"
cat > "$PLAN_FILE" << 'EOF'
# Implementation Plan
...
EOF

# Create backward-compatible symlink
ln -sf "$PLAN_FILE" IMPLEMENTATION-PLAN.md
```

## Migration Strategy

### Phase 1: Dual Support (Current)
- New files created in `.software-factory/` structure with timestamps
- Symlinks created for backward compatibility
- Detection checks both locations

### Phase 2: Deprecation Warnings
- Warn when legacy files found without new structure
- Encourage migration to new structure

### Phase 3: Full Migration
- Remove symlink creation
- Only support new structure
- Provide migration tool for old projects

## Benefits of New Structure

1. **Timestamp Tracking**: Know when files were created/updated
2. **Version History**: Multiple versions can coexist
3. **Organization**: Clear separation of plans, reports, and logs
4. **Namespace Isolation**: Each effort has its own subdirectory
5. **Audit Trail**: Complete history of all coordination files

## Implementation Checklist

When updating agent states or rules:

- [ ] Update file creation to use new `.software-factory/` structure
- [ ] Add timestamp to filename
- [ ] Create backward-compatible symlink
- [ ] Update detection patterns to check new location first
- [ ] Fall back to legacy location if new not found
- [ ] Test both new and legacy patterns work
- [ ] Document any new marker files here

## Related Rules

- **R290**: State rule verification markers
- **R287**: TODO persistence (uses todos/ directory)
- **R295**: Implementation plan authority
- **R222**: Code review trigger requirements

---

Last Updated: 2025-01-20
Version: 2.0.0 (New .software-factory/ structure)