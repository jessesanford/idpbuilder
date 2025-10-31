# 🔴🔴🔴 SUPREME LAW: RULE R371 - Effort Scope Immutability (BLOCKING)

## Rule Details
- **Rule ID**: R371
- **Rule Name**: Effort Scope Immutability
- **Criticality**: 🚨🚨🚨 BLOCKING
- **Category**: SUPREME LAW - Scope Control
- **Created**: 2025-01-21
- **Updated**: 2025-01-21

## 🔴 SUPREME DIRECTIVE

**EFFORT PLANS DEFINE ABSOLUTE SCOPE - NO ADDITIONS ALLOWED**

This is SUPREME LAW. Agents adding unrelated code to effort branches destroys the entire Software Factory architecture. An effort plan is a CONTRACT that defines EXACTLY what can and cannot be modified.

## 🚨 CATASTROPHIC VIOLATION EVIDENCE

**ACTUAL VIOLATIONS DETECTED:**
```
1. gitea-client-split-001: 124 files (MORE than unsplit with 89!)
2. registry-types: 45 files including unrelated SF artifacts
3. registry-helpers: 70 files with Kind contamination
4. Three fallback branches: All contaminated with unrelated files
5. Success rate: Only 18% of branches were clean
```

## ⚠️ MANDATORY REQUIREMENTS

### 1. EFFORT PLAN IS LAW
```yaml
IMMUTABLE RULES:
  - Effort plan defines EXACT scope
  - Every file change must trace to effort requirement
  - NO additions beyond plan allowed
  - NO "while I'm here" changes
  - NO "related" work outside scope
```

### 2. FILE CHANGE VALIDATION
```bash
# BEFORE ANY file modification:
validate_file_in_scope() {
    local file="$1"
    local effort_plan="$2"

    # File MUST be explicitly mentioned OR
    # Part of explicitly listed package/module
    if ! grep -q "$file" "$effort_plan"; then
        echo "🚨 SCOPE VIOLATION: $file not in effort plan!"
        return 1
    fi
}
```

### 3. SCOPE TRACEABILITY MATRIX
Every effort must maintain:
```markdown
## Scope Traceability
| File/Package | Effort Requirement | Justification |
|--------------|-------------------|---------------|
| pkg/gitea/client.go | "Create Gitea client" | Direct requirement |
| pkg/gitea/types.go | "Define request types" | Direct requirement |
```

### 4. OUT-OF-SCOPE DECLARATION
Effort plans MUST explicitly state what's NOT included:
```markdown
## OUT OF SCOPE (DO NOT MODIFY):
- Build system files (Makefile, go.mod)
- Test infrastructure (except unit tests for new code)
- DevContainer configurations
- Documentation files
- Unrelated packages
```

## 🔴 ENFORCEMENT MECHANISMS

### SW-ENGINEER ENFORCEMENT
```bash
# MANDATORY before ANY code writing:
1. READ effort plan completely
2. Extract allowed files/packages
3. Create scope validation checklist
4. REJECT any work not in plan
5. Re-read plan every 10 changes
```

### CODE-REVIEWER ENFORCEMENT
```bash
# MANDATORY scope validation:
1. Compare EVERY file against effort plan
2. Flag ANY file not traceable to requirements
3. Count files: splits must have FEWER than original
4. Create SCOPE-VIOLATION-REPORT.md
5. REJECT branches with ANY scope creep
```

### ORCHESTRATOR ENFORCEMENT
```bash
# MANDATORY planning discipline:
1. Create FOCUSED effort plans (10-20 files MAX)
2. One clear theme per effort
3. Explicit file/package lists
4. Clear OUT OF SCOPE section
5. Validate against wave/phase boundaries
```

## 🚨 VIOLATION PENALTIES

### IMMEDIATE FAILURES (-100%):
- Adding ANY file not in effort plan
- Split branch with MORE files than original
- Mixing multiple features in one branch
- "Kitchen sink" branches with unrelated work
- Ignoring OUT OF SCOPE declarations

### SEVERE PENALTIES (-50%):
- Vague effort plans without file lists
- Missing scope traceability
- No OUT OF SCOPE section
- Delayed scope validation

## 📊 VALIDATION CHECKLIST

```bash
# Run before EVERY commit:
validate_effort_scope() {
    echo "🔍 SCOPE VALIDATION CHECKLIST"

    # 1. Count files
    local file_count=$(git diff --name-only origin/main | wc -l)
    echo "Files changed: $file_count"

    # 2. Check each file against plan
    for file in $(git diff --name-only origin/main); do
        if ! grep -q "$file" EFFORT-PLAN.md; then
            echo "❌ VIOLATION: $file not in plan!"
            return 1
        fi
    done

    # 3. Verify no build system changes (unless planned)
    if git diff --name-only origin/main | grep -E "(Makefile|go.mod|package.json)"; then
        if ! grep -q "build system" EFFORT-PLAN.md; then
            echo "❌ VIOLATION: Unauthorized build system changes!"
            return 1
        fi
    fi

    echo "✅ All changes within scope"
}
```

## 🔴 EXAMPLE VIOLATIONS

### ❌ WRONG - Scope Creep:
```markdown
Effort: "Create Gitea client"
Actual changes:
- ✅ pkg/gitea/client.go (correct)
- ❌ Makefile (NOT in plan)
- ❌ .devcontainer/devcontainer.json (NOT in plan)
- ❌ tests/integration/setup.go (NOT in plan)
- ❌ docs/API.md (NOT in plan)
Result: 124 files changed (CATASTROPHIC FAILURE)
```

### ✅ CORRECT - Focused Scope:
```markdown
Effort: "Create Gitea client"
Actual changes:
- ✅ pkg/gitea/client.go (in plan)
- ✅ pkg/gitea/types.go (in plan)
- ✅ pkg/gitea/client_test.go (unit tests allowed)
Result: 3 files changed (PERFECT)
```

## 📝 MANDATORY ARTIFACTS

Each effort MUST produce:
1. `EFFORT-PLAN.md` with explicit file lists
2. `SCOPE-TRACEABILITY.md` mapping changes to requirements
3. `SCOPE-VALIDATION.log` showing all checks passed

## 🚨 EMERGENCY PROTOCOL

If scope violation detected:
1. **STOP ALL WORK IMMEDIATELY**
2. Revert ALL changes outside scope
3. Create `SCOPE-VIOLATION-INCIDENT.md`
4. Re-read and clarify effort plan
5. Get Orchestrator approval before continuing

---

**REMEMBER**: The effort plan is LAW. Not a suggestion, not a guideline - ABSOLUTE LAW. Every single file change must be traceable to an explicit requirement in the effort plan. NO EXCEPTIONS.