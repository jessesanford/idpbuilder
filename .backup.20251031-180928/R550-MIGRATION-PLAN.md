# R550 Planning File Migration - Audit and Execution Plan

**Project**: idpbuilder-oci-push-planning
**Created**: 2025-10-31
**Purpose**: Migrate all planning files to R550 canonical naming and state tracking

---

## AUDIT RESULTS

### Current Planning Files Found

**Project-Level Plans (Non-Canonical Naming)**:
- `planning/PROJECT-ARCHITECTURE.md` → Should be: `planning/project/PROJECT-ARCHITECTURE-PLAN.md`
- `planning/PROJECT-TEST-PLAN.md` → Should be: `planning/project/PROJECT-TEST-PLAN.md`
- `planning/README.md` → Keep as-is (documentation, not a plan)

**Phase 2 Plans (Non-Canonical Paths)**:
- `planning/phase2/PHASE-2-ARCHITECTURE.md` → Should be: `planning/phase2/PHASE-ARCHITECTURE-PLAN.md`
- `planning/phase2/PHASE-2-IMPLEMENTATION.md` → Should be: `planning/phase2/PHASE-IMPLEMENTATION-PLAN.md`
- `planning/phase2/PHASE-2-TEST-PLAN.md` → Should be: `planning/phase2/PHASE-TEST-PLAN.md`

**Example Files (Keep as Documentation)**:
- `planning/phase1/PHASE-ARCHITECTURE-PLAN.example.md` - Template
- `planning/phase1/PHASE-TEST-PLAN.example.md` - Template
- `planning/phase1/wave1/EFFORT-PLANS-SUMMARY.example.md` - Template
- `planning/phase1/wave1/WAVE-IMPLEMENTATION-PLAN.example.md` - Template
- `planning/phase1/wave1/WAVE-TEST-PLAN.example.md` - Template
- `planning/phase1/wave1/effort-001-repository-initialization.example.md` - Template
- `planning/phase2/PHASE-ARCHITECTURE-PLAN.example.md` - Template
- `planning/project/PROJECT-ARCHITECTURE-PLAN.example.md` - Template
- `planning/project/PROJECT-IMPLEMENTATION-PLAN.example.md` - Template
- `planning/project/PROJECT-TEST-PLAN.example.md` - Template

### Current Directory Structure

```
planning/
├── PROJECT-ARCHITECTURE.md          ← RENAME
├── PROJECT-TEST-PLAN.md             ← RENAME
├── README.md                         ← KEEP
├── phase1/
│   ├── PHASE-ARCHITECTURE-PLAN.example.md
│   ├── PHASE-TEST-PLAN.example.md
│   └── wave1/
│       ├── EFFORT-PLANS-SUMMARY.example.md
│       ├── WAVE-IMPLEMENTATION-PLAN.example.md
│       ├── WAVE-TEST-PLAN.example.md
│       └── effort-001-repository-initialization.example.md
├── phase2/
│   ├── PHASE-2-ARCHITECTURE.md      ← RENAME
│   ├── PHASE-2-IMPLEMENTATION.md    ← RENAME
│   ├── PHASE-2-TEST-PLAN.md         ← RENAME
│   └── PHASE-ARCHITECTURE-PLAN.example.md
└── project/
    ├── PROJECT-ARCHITECTURE-PLAN.example.md
    ├── PROJECT-IMPLEMENTATION-PLAN.example.md
    └── PROJECT-TEST-PLAN.example.md
```

### Current State Context

From `orchestrator-state-v3.json`:
- Current Phase: 2
- Current Wave: 1
- Phases Completed: []
- Waves Completed: 0

---

## MIGRATION MAPPING TABLE

| Current Path | Current Name | R550 Canonical Path | R550 Canonical Name | Type | Action |
|--------------|--------------|---------------------|---------------------|------|--------|
| `planning/` | `PROJECT-ARCHITECTURE.md` | `planning/project/` | `PROJECT-ARCHITECTURE-PLAN.md` | Project Architecture | MOVE + RENAME |
| `planning/` | `PROJECT-TEST-PLAN.md` | `planning/project/` | `PROJECT-TEST-PLAN.md` | Project Test | MOVE |
| `planning/phase2/` | `PHASE-2-ARCHITECTURE.md` | `planning/phase2/` | `PHASE-ARCHITECTURE-PLAN.md` | Phase Architecture | RENAME |
| `planning/phase2/` | `PHASE-2-IMPLEMENTATION.md` | `planning/phase2/` | `PHASE-IMPLEMENTATION-PLAN.md` | RENAME |
| `planning/phase2/` | `PHASE-2-TEST-PLAN.md` | `planning/phase2/` | `PHASE-TEST-PLAN.md` | Phase Test | RENAME |
| `planning/README.md` | - | `planning/README.md` | - | Documentation | KEEP |
| `planning/phase1/*.example.md` | - | (same) | - | Templates | KEEP |
| `planning/phase2/*.example.md` | - | (same) | - | Templates | KEEP |
| `planning/project/*.example.md` | - | (same) | - | Templates | KEEP |

---

## VIOLATIONS DETECTED

### R550 Violations

1. ❌ **Project-level plans in wrong directory**:
   - `planning/PROJECT-ARCHITECTURE.md` should be in `planning/project/`
   - `planning/PROJECT-TEST-PLAN.md` should be in `planning/project/`

2. ❌ **Phase naming with numeric prefix**:
   - `PHASE-2-ARCHITECTURE.md` should be `PHASE-ARCHITECTURE-PLAN.md`
   - `PHASE-2-IMPLEMENTATION.md` should be `PHASE-IMPLEMENTATION-PLAN.md`
   - `PHASE-2-TEST-PLAN.md` should be `PHASE-TEST-PLAN.md`
   - (Phase number comes from directory `phase2/`, not filename)

3. ❌ **Missing "-PLAN" suffix**:
   - `PROJECT-ARCHITECTURE.md` should be `PROJECT-ARCHITECTURE-PLAN.md`

4. ⚠️ **State tracking missing**:
   - `orchestrator-state-v3.json` has NO `planning_files` field

---

## STATE UPDATES REQUIRED

### Add planning_files Structure

```json
"planning_files": {
  "project": {
    "architecture_plan": "planning/project/PROJECT-ARCHITECTURE-PLAN.md",
    "test_plan": "planning/project/PROJECT-TEST-PLAN.md"
  },
  "current_phase": {
    "architecture_plan": "planning/phase2/PHASE-ARCHITECTURE-PLAN.md",
    "implementation_plan": "planning/phase2/PHASE-IMPLEMENTATION-PLAN.md",
    "test_plan": "planning/phase2/PHASE-TEST-PLAN.md"
  },
  "current_wave": {},
  "efforts": {}
}
```

---

## EXECUTION PLAN

### Step 1: Verify Directory Structure

```bash
# Check that planning/project/ exists
ls -la planning/project/

# Expected output: Should exist with example files
```

### Step 2: Move and Rename Files

```bash
# 1. Move project-level architecture plan
mv planning/PROJECT-ARCHITECTURE.md planning/project/PROJECT-ARCHITECTURE-PLAN.md

# 2. Move project-level test plan
mv planning/PROJECT-TEST-PLAN.md planning/project/PROJECT-TEST-PLAN.md

# 3. Rename phase2 architecture plan
mv planning/phase2/PHASE-2-ARCHITECTURE.md planning/phase2/PHASE-ARCHITECTURE-PLAN.md

# 4. Rename phase2 implementation plan
mv planning/phase2/PHASE-2-IMPLEMENTATION.md planning/phase2/PHASE-IMPLEMENTATION-PLAN.md

# 5. Rename phase2 test plan
mv planning/phase2/PHASE-2-TEST-PLAN.md planning/phase2/PHASE-TEST-PLAN.md
```

### Step 3: Update orchestrator-state-v3.json

Add `planning_files` structure at top level with all plan paths.

### Step 4: Verify All Files

```bash
# Check all canonical paths exist
ls -la planning/project/PROJECT-ARCHITECTURE-PLAN.md
ls -la planning/project/PROJECT-TEST-PLAN.md
ls -la planning/phase2/PHASE-ARCHITECTURE-PLAN.md
ls -la planning/phase2/PHASE-IMPLEMENTATION-PLAN.md
ls -la planning/phase2/PHASE-TEST-PLAN.md
```

### Step 5: Verify No Old Files Remain

```bash
# Check old paths are gone
! ls planning/PROJECT-ARCHITECTURE.md 2>/dev/null
! ls planning/PROJECT-TEST-PLAN.md 2>/dev/null
! ls planning/phase2/PHASE-2-ARCHITECTURE.md 2>/dev/null
! ls planning/phase2/PHASE-2-IMPLEMENTATION.md 2>/dev/null
! ls planning/phase2/PHASE-2-TEST-PLAN.md 2>/dev/null
```

---

## R550 COMPLIANCE VERIFICATION

After migration, verify:

✅ All planning files use canonical naming (no numeric prefixes, proper suffixes)
✅ All files in correct R550 directory structure
✅ `orchestrator-state-v3.json` tracks all planning file paths
✅ `planning_files` at top level of state
✅ All documented files exist at specified paths
✅ No phase-plans/ directory (legacy structure removed)
✅ No timestamped filenames (R550 prohibits timestamps)

---

## SUCCESS CRITERIA

- [ ] All 5 files moved/renamed
- [ ] orchestrator-state-v3.json updated with planning_files
- [ ] All canonical paths verified to exist
- [ ] All old paths verified removed
- [ ] Git commit created with all changes
- [ ] Git push completed
- [ ] Summary report generated

---

**Next**: Execute migration steps
