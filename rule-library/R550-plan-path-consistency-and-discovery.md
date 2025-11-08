# 🔴🔴🔴 RULE R550 - PLAN PATH CONSISTENCY AND DISCOVERY (SUPREME LAW)

**Criticality:** 🔴🔴🔴 SUPREME LAW - ABSOLUTE REQUIREMENT
**Category**: Planning / File Management / State Tracking
**Exit Code**: 550
**Grading Impact:** -20% to -100% for violations

---

## SUPREME LAW STATEMENT

**ALL PLANNING DOCUMENTS MUST FOLLOW STANDARDIZED NAMING CONVENTIONS AND BE TRACKED IN ORCHESTRATOR-STATE-V3.JSON FOR RELIABLE DISCOVERY. FILESYSTEM SEARCHING IS PROHIBITED - USE STATE FILE LOOKUPS.**

---

## 🚨🚨🚨 THE ABSOLUTE MANDATE 🚨🚨🚨

### Planning File Discovery IS:
```markdown
✅ DETERMINISTIC - Always know where plans are
✅ TRACKED - orchestrator-state-v3.json records all paths
✅ FAST - O(1) lookup from state file
✅ RELIABLE - No filesystem searching required
✅ VALIDATED - Schema enforcement prevents missing plans
```

### Planning File Discovery IS NOT:
```markdown
❌ Filesystem searching with `ls -t` or `find`
❌ Guessing paths from conventions
❌ Trial-and-error lookups
❌ Different paths for creation vs reading
❌ Untracked files in unknown locations
```

---

## 🔴🔴🔴 CANONICAL DIRECTORY STRUCTURE 🔴🔴🔴

### Standard Paths

#### Project Level
```
planning/project/
├── PROJECT-ARCHITECTURE-PLAN.md
├── PROJECT-IMPLEMENTATION-PLAN.md
└── PROJECT-TEST-PLAN.md
```

#### Phase Level
```
planning/phase{N}/
├── PHASE-ARCHITECTURE-PLAN.md
├── PHASE-TEST-PLAN.md
└── wave{W}/
    ├── WAVE-IMPLEMENTATION-PLAN.md
    ├── WAVE-TEST-PLAN.md
    └── effort-{id}-{name}.md
```

### Complete Example
```
planning/
├── project/
│   ├── PROJECT-ARCHITECTURE-PLAN.md
│   ├── PROJECT-IMPLEMENTATION-PLAN.md
│   └── PROJECT-TEST-PLAN.md
├── phase1/
│   ├── PHASE-ARCHITECTURE-PLAN.md
│   ├── PHASE-TEST-PLAN.md
│   └── wave1/
│       ├── WAVE-IMPLEMENTATION-PLAN.md
│       ├── WAVE-TEST-PLAN.md
│       └── effort-001-repository-initialization.md
│   └── wave2/
│       ├── WAVE-IMPLEMENTATION-PLAN.md
│       ├── WAVE-TEST-PLAN.md
│       └── effort-002-feature-implementation.md
└── phase2/
    └── [similar structure]
```

### Naming Rules

| Rule | Requirement | Example |
|------|-------------|---------|
| Root Directory | ALWAYS `planning/` | `planning/phase1/` ✅ `phase-plans/phase1/` ❌ |
| Phase Numbers | NO hyphen before number | `phase1` ✅ `phase-1` ❌ |
| Wave Numbers | NO hyphen before number | `wave2` ✅ `wave-2` ❌ |
| Plan Types | EXPLICIT descriptive names | `IMPLEMENTATION-PLAN` ✅ `PLAN` ❌ |
| Separators | Single hyphen | `WAVE-IMPLEMENTATION-PLAN` ✅ `WAVE--PLAN` ❌ |
| Timestamps | NEVER in planning/ filenames | `PHASE-ARCHITECTURE-PLAN.md` ✅ `PHASE-ARCHITECTURE-PLAN--20251031-142500.md` ❌ |
| Case | ALL CAPS for plan files | `PROJECT-ARCHITECTURE-PLAN.md` ✅ `project-architecture-plan.md` ❌ |

---

## 🔴🔴🔴 ORCHESTRATOR STATE TRACKING 🔴🔴🔴

### Schema Structure

The `planning_files` section in `orchestrator-state-v3.json`:

```json
{
  "planning_files": {
    "project": {
      "architecture_plan": "planning/project/PROJECT-ARCHITECTURE-PLAN.md",
      "implementation_plan": "planning/project/PROJECT-IMPLEMENTATION-PLAN.md",
      "test_plan": "planning/project/PROJECT-TEST-PLAN.md"
    },
    "phases": {
      "phase1": {
        "architecture_plan": "planning/phase1/PHASE-ARCHITECTURE-PLAN.md",
        "test_plan": "planning/phase1/PHASE-TEST-PLAN.md",
        "waves": {
          "wave1": {
            "implementation_plan": "planning/phase1/wave1/WAVE-IMPLEMENTATION-PLAN.md",
            "test_plan": "planning/phase1/wave1/WAVE-TEST-PLAN.md",
            "efforts": {
              "effort-001-repository-initialization": "planning/phase1/wave1/effort-001-repository-initialization.md"
            }
          },
          "wave2": {
            "implementation_plan": "planning/phase1/wave2/WAVE-IMPLEMENTATION-PLAN.md",
            "test_plan": "planning/phase1/wave2/WAVE-TEST-PLAN.md",
            "efforts": {
              "effort-002-feature-implementation": "planning/phase1/wave2/effort-002-feature-implementation.md"
            }
          }
        }
      }
    }
  }
}
```

---

## 🔴🔴🔴 ORCHESTRATOR RESPONSIBILITIES 🔴🔴🔴

### 1. On Plan Receipt (WAITING_FOR_* States)

When receiving plans from agents, orchestrator MUST:

```bash
# 1. Verify Plan Location
PHASE=$(jq -r '.project_progression.current_phase.phase_number' orchestrator-state-v3.json)
EXPECTED_PATH="planning/phase${PHASE}/PHASE-ARCHITECTURE-PLAN.md"

if [ ! -f "$EXPECTED_PATH" ]; then
    echo "❌ Plan not at standard location: $EXPECTED_PATH"
    echo "Searching for plan..."

    # Search common wrong locations
    FOUND_PATH=$(find planning/ -name "PHASE*ARCHITECTURE*.md" -type f | head -1)

    if [ -n "$FOUND_PATH" ] && [ -f "$FOUND_PATH" ]; then
        echo "⚠️ Found plan at non-standard location: $FOUND_PATH"
        echo "Moving to standard location..."
        mkdir -p "$(dirname "$EXPECTED_PATH")"
        mv "$FOUND_PATH" "$EXPECTED_PATH"
    else
        echo "❌ CRITICAL: Cannot find phase ${PHASE} architecture plan anywhere!"
        exit 550
    fi
fi

# 2. Document in Orchestrator State
jq ".planning_files.phases.phase${PHASE}.architecture_plan = \"$EXPECTED_PATH\"" \
   orchestrator-state-v3.json > orchestrator-state-v3.json.tmp
mv orchestrator-state-v3.json.tmp orchestrator-state-v3.json

# 3. Verify Accessibility
if [ ! -f "$EXPECTED_PATH" ] || [ ! -r "$EXPECTED_PATH" ]; then
    echo "❌ Plan file not accessible: $EXPECTED_PATH"
    exit 550
fi

# Confirm tracking
TRACKED_PATH=$(jq -r ".planning_files.phases.phase${PHASE}.architecture_plan" orchestrator-state-v3.json)
if [ "$TRACKED_PATH" != "$EXPECTED_PATH" ]; then
    echo "❌ Path not tracked correctly in state file!"
    exit 550
fi

echo "✅ Plan verified and tracked: $EXPECTED_PATH"
```

### 2. On Plan Request (All States)

When needing to read a plan, orchestrator MUST:

```bash
# PREFERRED METHOD: Read from orchestrator state
PHASE=$(jq -r '.project_progression.current_phase.phase_number' orchestrator-state-v3.json)
PLAN_PATH=$(jq -r ".planning_files.phases.phase${PHASE}.architecture_plan" orchestrator-state-v3.json)

# FALLBACK: Standard path if not tracked yet
if [ -z "$PLAN_PATH" ] || [ "$PLAN_PATH" = "null" ] || [ ! -f "$PLAN_PATH" ]; then
    echo "⚠️ Plan path not in state file, using standard location"
    PLAN_PATH="planning/phase${PHASE}/PHASE-ARCHITECTURE-PLAN.md"
fi

# VERIFY: File must exist
if [ ! -f "$PLAN_PATH" ]; then
    echo "❌ CRITICAL: Cannot find phase ${PHASE} architecture plan"
    echo "Expected: planning/phase${PHASE}/PHASE-ARCHITECTURE-PLAN.md"
    echo "State file says: $(jq -r ".planning_files.phases.phase${PHASE}.architecture_plan" orchestrator-state-v3.json)"
    echo ""
    echo "Available planning files:"
    find planning/ -name "*.md" -type f
    exit 550
fi

# USE: Plan is now safe to read
cat "$PLAN_PATH"
```

### 3. On Plan Missing (Error Handling)

```bash
# If plan truly missing after all fallbacks:
echo "❌❌❌ R550 VIOLATION: Missing Planning Document ❌❌❌"
echo ""
echo "Expected at: planning/phase${PHASE}/PHASE-ARCHITECTURE-PLAN.md"
echo "State file tracking: $(jq -r ".planning_files.phases.phase${PHASE}.architecture_plan // \"NOT TRACKED\"" orchestrator-state-v3.json)"
echo ""
echo "This is a CRITICAL FAILURE - cannot proceed without planning documents!"
echo ""
echo "Transitioning to ERROR_RECOVERY per R550"
exit 550
```

---

## 🔴🔴🔴 AGENT RESPONSIBILITIES 🔴🔴🔴

### Plan Creation (Architect, Code Reviewer)

When creating plans, agents MUST use standard path construction:

```bash
# Phase Architecture Plan
PHASE=$(jq -r '.project_progression.current_phase.phase_number' orchestrator-state-v3.json)
PLAN_PATH="planning/phase${PHASE}/PHASE-ARCHITECTURE-PLAN.md"

# Create directory if needed
mkdir -p "$(dirname "$PLAN_PATH")"

# Write plan to standard location
cat > "$PLAN_PATH" << 'EOF'
# Phase ${PHASE} Architecture Plan

[Plan content here]
EOF

echo "✅ Created plan at standard location: $PLAN_PATH"
```

### Plan Reading (All Agents)

When reading plans, agents MUST:

```bash
# Check orchestrator state first
PLAN_PATH=$(jq -r ".planning_files.phases.phase${PHASE}.architecture_plan" orchestrator-state-v3.json)

# Fallback to standard path
if [ -z "$PLAN_PATH" ] || [ "$PLAN_PATH" = "null" ]; then
    PLAN_PATH="planning/phase${PHASE}/PHASE-ARCHITECTURE-PLAN.md"
fi

# Verify exists
if [ ! -f "$PLAN_PATH" ]; then
    echo "❌ Cannot find plan - reporting to orchestrator"
    exit 550
fi

# Read plan
cat "$PLAN_PATH"
```

---

## 🔴🔴🔴 ABSOLUTE PROHIBITIONS 🔴🔴🔴

### ❌ NEVER Do These

#### 1. Filesystem Searching
```bash
# ❌ WRONG - Violates R340 and R550
PLAN=$(ls -t planning/phase*/PHASE*ARCHITECTURE*.md | head -1)
PLAN=$(find planning/ -name "*PHASE*PLAN*.md" -type f | head -1)

# ✅ RIGHT - Use state file
PLAN=$(jq -r ".planning_files.phases.phase${PHASE}.architecture_plan" orchestrator-state-v3.json)
```

#### 2. Wrong Directory
```bash
# ❌ WRONG - phase-plans/ doesn't exist
PLAN="phase-plans/phase1/PHASE-1-ARCHITECTURE-PLAN.md"

# ✅ RIGHT - planning/ is canonical
PLAN="planning/phase1/PHASE-ARCHITECTURE-PLAN.md"
```

#### 3. Timestamps in Planning Filenames
```bash
# ❌ WRONG - Timestamps make discovery unpredictable
PLAN="planning/phase1/PHASE-ARCHITECTURE-PLAN--20251031-142500.md"

# ✅ RIGHT - No timestamps in planning/ directory
PLAN="planning/phase1/PHASE-ARCHITECTURE-PLAN.md"

# NOTE: Timestamps ARE allowed in .software-factory/ metadata
METADATA=".software-factory/phase1/wave1/effort-001/IMPLEMENTATION-PLAN--20251031-142500.md"
```

#### 4. Generic "PLAN" Names
```bash
# ❌ WRONG - Ambiguous what type of plan
PLAN="planning/phase1/PHASE-1-PLAN.md"

# ✅ RIGHT - Explicit plan type
PLAN="planning/phase1/PHASE-ARCHITECTURE-PLAN.md"
PLAN="planning/phase1/PHASE-TEST-PLAN.md"
```

---

## 🔴🔴🔴 GRADING PENALTIES 🔴🔴🔴

### Violation Types and Penalties

```yaml
filesystem_searching:
  offense: Using ls -t, find, or glob to discover plans
  penalty: -30% per occurrence
  reason: Violates R340 (no filesystem searching) and R550

wrong_directory:
  offense: Using phase-plans/ instead of planning/
  penalty: -40% per occurrence
  reason: Directory doesn't exist, guaranteed failure

path_not_tracked:
  offense: Creating plan but not recording in state file
  penalty: -20% per plan
  reason: Future lookups will fail

wrong_naming:
  offense: Using non-canonical filenames (timestamps, hyphens, etc.)
  penalty: -25% per file
  reason: Discovery will fail, inconsistent with examples

plan_not_found:
  offense: Orchestrator can't find plan after creation
  penalty: -100% IMMEDIATE FAILURE
  reason: CRITICAL system failure, blocks all progress

missing_fallback:
  offense: No fallback logic when state file lookup fails
  penalty: -20%
  reason: Brittle system, recovery impossible
```

### Catastrophic Failures (Automatic -100%)

1. **Plan Created But Unfindable**
   - Agent creates plan at non-standard location
   - Orchestrator can't discover it
   - System deadlocked

2. **Wrong Directory Used**
   - Creating plans in `phase-plans/` (doesn't exist)
   - Guaranteed filesystem errors

3. **No State Tracking**
   - Plans created but not recorded
   - No way to find them later
   - Violates R340

---

## 🔴🔴🔴 VALIDATION REQUIREMENTS 🔴🔴🔴

### Pre-Commit Validation

```bash
#!/bin/bash
# File: .git/hooks/pre-commit.d/validate-plan-paths.sh

# Validate no phase-plans/ references
if git diff --cached --name-only | xargs grep -l "phase-plans/" 2>/dev/null; then
    echo "❌ R550 VIOLATION: References to obsolete 'phase-plans/' directory found"
    echo ""
    echo "Use 'planning/' instead:"
    git diff --cached --name-only | xargs grep -n "phase-plans/" 2>/dev/null
    exit 550
fi

# Validate planning files use canonical names
if git diff --cached --name-only | grep "^planning/.*\.md$" | while read file; do
    # Check for timestamps in filename
    if echo "$file" | grep -q -- "--[0-9]\{8\}"; then
        echo "❌ R550 VIOLATION: Timestamp in planning filename: $file"
        echo "Remove timestamp - planning/ files should not have timestamps"
        exit 550
    fi

    # Check for generic "PLAN" names
    if echo "$file" | grep -qE "PLAN\.md$" && ! echo "$file" | grep -qE "(ARCHITECTURE|IMPLEMENTATION|TEST)-PLAN\.md$"; then
        echo "❌ R550 VIOLATION: Generic PLAN name: $file"
        echo "Use explicit type: ARCHITECTURE-PLAN, IMPLEMENTATION-PLAN, or TEST-PLAN"
        exit 550
    fi
done; then
    exit 550
fi

echo "✅ R550: Planning path validation passed"
```

### State File Validation

```bash
#!/bin/bash
# File: tools/validate-planning-paths.sh

validate_planning_paths() {
    echo "🔍 R550: Validating planning path tracking..."

    # Check planning_files exists in state
    if ! jq -e '.planning_files' orchestrator-state-v3.json > /dev/null 2>&1; then
        echo "❌ R550 VIOLATION: No planning_files in orchestrator state"
        return 1
    fi

    # Validate all tracked files exist
    local missing=0
    jq -r '.. | .architecture_plan? // .implementation_plan? // .test_plan? | select(. != null)' orchestrator-state-v3.json | while read path; do
        if [ ! -f "$path" ]; then
            echo "❌ Tracked plan missing: $path"
            ((missing++))
        fi
    done

    if [ $missing -gt 0 ]; then
        echo "❌ R550 VIOLATION: $missing tracked plans missing from filesystem"
        return 1
    fi

    # Validate all planning/ files are tracked
    find planning/ -name "*.md" -type f ! -name "*.example.md" | while read file; do
        if ! jq -r '.. | .architecture_plan? // .implementation_plan? // .test_plan? | select(. != null)' orchestrator-state-v3.json | grep -q "^$file$"; then
            echo "⚠️ WARNING: Untracked planning file: $file"
        fi
    done

    echo "✅ R550: Planning path validation passed"
}

validate_planning_paths
```

---

## 📋 EXAMPLES

### ✅ CORRECT USAGE

#### Creating Phase Architecture Plan
```bash
PHASE=$(jq -r '.project_progression.current_phase.phase_number' orchestrator-state-v3.json)
PLAN_PATH="planning/phase${PHASE}/PHASE-ARCHITECTURE-PLAN.md"

mkdir -p "$(dirname "$PLAN_PATH")"

cat > "$PLAN_PATH" << 'EOF'
# Phase 1 Architecture Plan

## Architecture Overview
[Content here]
EOF

# Track in state file
jq ".planning_files.phases.phase${PHASE}.architecture_plan = \"$PLAN_PATH\"" \
   orchestrator-state-v3.json > tmp && mv tmp orchestrator-state-v3.json

echo "✅ Created and tracked: $PLAN_PATH"
```

#### Reading Wave Implementation Plan
```bash
PHASE=$(jq -r '.project_progression.current_phase.phase_number' orchestrator-state-v3.json)
WAVE=$(jq -r '.project_progression.current_wave.wave_number' orchestrator-state-v3.json)

# Method 1: State file lookup
PLAN_PATH=$(jq -r ".planning_files.phases.phase${PHASE}.waves.wave${WAVE}.implementation_plan" orchestrator-state-v3.json)

# Method 2: Fallback to standard
if [ -z "$PLAN_PATH" ] || [ "$PLAN_PATH" = "null" ]; then
    PLAN_PATH="planning/phase${PHASE}/wave${WAVE}/WAVE-IMPLEMENTATION-PLAN.md"
fi

# Validate
if [ ! -f "$PLAN_PATH" ]; then
    echo "❌ Cannot find wave implementation plan"
    exit 550
fi

# Read
cat "$PLAN_PATH"
```

### ❌ INCORRECT USAGE

#### Wrong Directory
```bash
# ❌ WRONG - phase-plans/ doesn't exist
PLAN="phase-plans/phase1/PHASE-ARCHITECTURE-PLAN.md"
```

#### Filesystem Searching
```bash
# ❌ WRONG - Violates R340 and R550
PLAN=$(ls -t planning/phase*/PHASE*ARCHITECTURE*.md | head -1)
```

#### Not Tracking
```bash
# ❌ WRONG - Created but not tracked
cat > planning/phase1/PHASE-ARCHITECTURE-PLAN.md << 'EOF'
[Content]
EOF
# Missing: jq update to record path!
```

#### Timestamp in Filename
```bash
# ❌ WRONG - Timestamps make discovery unpredictable
PLAN="planning/phase1/PHASE-ARCHITECTURE-PLAN--20251031-142500.md"
```

---

## 🔴 INTEGRATION WITH OTHER RULES

### R340 - Metadata Location Tracking
- **Relationship**: R550 extends R340 to planning files
- **Enforcement**: All planning paths tracked in state file
- **Synergy**: NO filesystem searching, instant lookups

### R288 - State File Updates
- **Relationship**: Plan path recording is state update
- **Enforcement**: Must commit state file after tracking plans
- **Synergy**: Atomic plan creation + state update

### R405 - Automation Continuation Flag
- **Relationship**: Missing plans = CONTINUE-SOFTWARE-FACTORY=FALSE
- **Enforcement**: Cannot continue without discoverable plans
- **Synergy**: Clear blocking condition

### R502 - Mandatory Plan Validation Gates
- **Relationship**: R550 defines WHERE plans are, R502 defines WHEN needed
- **Enforcement**: Gates check both existence and tracking
- **Synergy**: Complete plan validation

---

## 📊 SUCCESS METRICS

### System Health Indicators

```yaml
plan_discoverability:
  target: 100%
  measurement: "Plans found / Plans requested"
  threshold: "100% - Any failure is critical"

state_tracking_coverage:
  target: 100%
  measurement: "Plans tracked / Plans created"
  threshold: "100% - All plans must be tracked"

filesystem_search_rate:
  target: 0%
  measurement: "grep 'ls -t|find' state rules"
  threshold: "0 - Absolutely forbidden"

path_consistency:
  target: 100%
  measurement: "planning/ refs / (planning/ + phase-plans/ refs)"
  threshold: "100% - No phase-plans/ references allowed"
```

---

## 🔄 MIGRATION GUIDANCE

### For Existing Projects

If your project is using the legacy `planning_artifacts.master_architecture_file` pattern:

**Legacy Pattern (DEPRECATED):**
```json
{
  "planning_artifacts": {
    "master_architecture_file": "planning/PROJECT-ARCHITECTURE.md",
    "master_architecture_status": "COMPLETE"
  }
}
```

**R550 Pattern (CURRENT STANDARD):**
```json
{
  "planning_files": {
    "project": {
      "architecture_plan": "planning/project/PROJECT-ARCHITECTURE-PLAN.md",
      "implementation_plan": "planning/project/PROJECT-IMPLEMENTATION-PLAN.md",
      "test_plan": "planning/project/PROJECT-TEST-PLAN.md"
    }
  }
}
```

**Migration Steps:**

1. **Add `planning_files` section** to orchestrator-state-v3.json
2. **Keep `planning_artifacts`** temporarily for backward compatibility
3. **Update all references** in code to use `planning_files` first with fallback:
   ```bash
   ARCH_PLAN=$(jq -r '.planning_files.project.architecture_plan // .planning_artifacts.master_architecture_file // "planning/project/PROJECT-ARCHITECTURE-PLAN.md"' state.json)
   ```
4. **Verify all agents** are using R550-compliant paths
5. **Remove `planning_artifacts`** in a future major version

### Schema Version 3.0.0+

As of SF 3.0, the schema marks `planning_artifacts` as **DEPRECATED**:
- Field remains in schema for backward compatibility
- New projects MUST use `planning_files`
- Old projects SHOULD migrate to `planning_files`
- `planning_artifacts` will be removed in SF 4.0

---

## THE PLAN PATH OATH

```
I, the Agent, swear by R550:

When I create a planning document:
- I WILL use the planning/ directory (never phase-plans/)
- I WILL use canonical filenames (no timestamps, explicit types)
- I WILL record the path in planning_files (NOT planning_artifacts)
- I WILL verify the file is accessible

When I read a planning document:
- I WILL check planning_files FIRST (orchestrator state file)
- I WILL fallback to planning_artifacts for legacy support
- I WILL fallback to standard path if not tracked
- I WILL report clear errors if not found
- I WILL NEVER use filesystem searching

My planning files are discoverable.
My paths are tracked and predictable.
My system is reliable and fast.

This is SUPREME LAW.
Violation means FAILURE.
I WILL FOLLOW R550 ABSOLUTELY.
```

---

**Remember:** Consistent planning file paths are CRITICAL for system reliability. R550 ensures plans are always findable, always tracked, and always predictable. NO filesystem searching, NO guessing, NO failures.

**See Also:**
- R340: Metadata Location Tracking (planning extends this)
- R502: Mandatory Plan Validation Gates (when plans needed)
- R288: State File Updates (recording plan paths)
- R405: Automation Continuation Flag (blocking on missing plans)
- Rule Library: All planning-related rules

**Migration Status:**
- Schema updated with deprecation notices (v3.0.0)
- Agent states updated to use R550 pattern
- Test fixtures updated with both fields
- Backward compatibility maintained
