# 🔴🔴🔴 RULE R532: Template Metadata Path Validation (BLOCKING) 🔴🔴🔴

## Rule Identifier
**Rule ID**: R532
**Category**: Template Validation & Code Quality
**Criticality**: BLOCKING - Prevents R383 violations at source
**Priority**: ABSOLUTE
**Penalty**: -100% for merged templates with violations
**Introduced**: 2025-10-06

## Description

ALL agent templates (configs and state rules) MUST be validated to ensure they do NOT contain hardcoded paths that violate R383. This rule PREVENTS the root cause of metadata file violations by catching template bugs BEFORE they reach production.

## Rationale

**The Problem We Solve:**
- Agents read templates literally
- If template says `cat > REPORT.md`, agent creates file in root
- Even if agents know R383, hardcoded template code beats general rules
- Result: 41+ metadata files in wrong locations (Phase 2 integration)

**Why Template Validation:**
- **Prevention > Detection**: Stop violations at the source
- **Defense in Depth**: R383 defines requirements, R532 enforces them
- **Automation**: Pre-commit hooks catch violations before merge
- **Documentation**: Templates become living examples of correct patterns

## ROOT CAUSE ANALYSIS

### Case Study: Integration Agent INTEGRATE_WAVE_EFFORTS-REPORT.md

**Bad Template (R383 Violation):**
```bash
# From .claude/agents/integration.md line 562
cat > INTEGRATE_WAVE_EFFORTS-REPORT.md << EOF
# Integration Report
...
EOF
```

**What Happened:**
1. Agent reads R383: "metadata goes in .software-factory/"
2. Agent reads template: `cat > INTEGRATE_WAVE_EFFORTS-REPORT.md`
3. Agent executes template code literally
4. **Hardcoded path beats general rule**
5. File created in project root ❌

**Good Template (R383 Compliant):**
```bash
# Use sf_metadata_path helper
REPORT_PATH=$(sf_metadata_path "$PHASE" "$WAVE" "integration" "INTEGRATE_WAVE_EFFORTS-REPORT" "md")
cat > "$REPORT_PATH" << EOF
# Integration Report
...
EOF
```

**Result:**
- File created in `.software-factory/phase2/wave1/integration/INTEGRATE_WAVE_EFFORTS-REPORT--20251006-143000.md` ✅
- Timestamp prevents merge conflicts ✅
- Proper structure for tracking ✅

## ABSOLUTE REQUIREMENTS

### 1. MANDATORY VALIDATION SCRIPT

**Location:** `/tools/validate-templates.sh`

**Purpose:** Scan ALL templates for R383 violations

**Must Check:**
```bash
# For every agent template file:
.claude/agents/*.md
agent-states/**/rules.md

# Find patterns like:
cat > METADATA-FILE.md        # ❌ REJECT
cat > REPORT.md                # ❌ REJECT
echo "..." > PLAN.md           # ❌ REJECT

# Accept patterns like:
sf_metadata_path ... "REPORT"  # ✅ ACCEPT
cat > .software-factory/phase  # ✅ ACCEPT (if structured)
cat > src/code.go              # ✅ ACCEPT (not metadata)
```

**Run Frequency:**
- ✅ Pre-commit hook (automatic)
- ✅ Manual before merging template changes
- ✅ CI/CD pipeline validation
- ✅ Weekly automated scans

### 2. METADATA FILE PATTERNS

**These patterns MUST trigger validation:**

```bash
# Integration Reports
INTEGRATE_WAVE_EFFORTS-REPORT*.md
PHASE-MERGE-PLAN*.md
WAVE-MERGE-PLAN*.md
PROJECT-MERGE-PLAN*.md

# Code Review Artifacts
CODE-REVIEW-REPORT*.md
SPLIT-INVENTORY*.md
SPLIT-PLAN*.md
FIX-PLAN*.md

# Planning Documents
IMPLEMENTATION-PLAN*.md
TEST-PLAN*.md
DEMO-PLAN*.md
DEMO-STATUS*.md

# PR-Ready Artifacts
MASTER-PR-PLAN*.md
PR-BODY*.md
PR-VALIDATION-REPORT*.md

# Build & Validation
BUILD-VALIDATION-REPORT*.md
TEST-REPORT*.md
VALIDATION-RESULTS*.md
```

### 3. ALLOWED ROOT FILES (NOT METADATA)

**These are acceptable in project root:**

```bash
# Executable scripts
demo-wave.sh
demo-features.sh
integration-demo.sh
verify-prs.sh

# Source code files
*.go, *.py, *.js, *.ts, *.java
src/**, lib/**, pkg/**

# Project configuration
.gitignore
README.md
package.json
requirements.txt
go.mod
Makefile

# Temporary files
/tmp/**
```

### 4. VALIDATION SCRIPT IMPLEMENTATION

**Complete script:** `tools/validate-templates.sh`

**Key Logic:**
```bash
#!/bin/bash
# R532 Validation Script

# 1. Find all cat/echo > commands in templates
grep -rn "cat\|echo.*>" .claude/agents/ agent-states/

# 2. Extract target filenames
# 3. Check if target is metadata (matches patterns above)
# 4. If metadata, verify it uses:
#    - .software-factory/ prefix, OR
#    - sf_metadata_path() helper, OR
#    - Is a variable like $REPORT_PATH

# 5. Report violations with:
#    - File path
#    - Line number
#    - Violation details
#    - Fix suggestion

# 6. Exit 1 if any violations found
```

**Example Output:**
```
🔍 R532 TEMPLATE VALIDATION - Starting...

📋 Checking agent configs (.claude/agents/*.md)...
  Checking: integration.md
    ✓ Line 588: INTEGRATE_WAVE_EFFORTS-REPORT (uses sf_metadata_path)
    ✗ Line 686: INTEGRATE_WAVE_EFFORTS-REPORT.md ← VIOLATION!
       Content: cat > INTEGRATE_WAVE_EFFORTS-REPORT.md

❌ FAILED: 1 R383/R532 violation found

FIX REQUIRED:
1. Replace 'cat > METADATA-FILE.md' with sf_metadata_path()
2. Ensure all metadata goes to .software-factory/phase*/wave*/[agent]/
3. Add timestamps using --YYYYMMDD-HHMMSS format
```

### 5. PRE-COMMIT HOOK INTEGRATE_WAVE_EFFORTS

**Add to `.pre-commit-config.yaml`:**

```yaml
repos:
  - repo: local
    hooks:
      - id: validate-templates-r532
        name: R532 - Validate Template Metadata Paths
        entry: tools/validate-templates.sh
        language: script
        files: '\.(claude/agents/.*\.md|agent-states/.*/rules\.md)$'
        pass_filenames: false
        description: 'Ensures agent templates comply with R383 metadata requirements'
```

**Benefits:**
- Automatic validation on every commit
- Blocks commits with violations
- Immediate feedback to developers
- Zero-cost enforcement

### 6. TEMPLATE FIX PATTERNS

**Pattern 1: Simple Metadata File**

**Before (WRONG):**
```bash
cat > CODE-REVIEW-REPORT.md << EOF
# Code Review Report
...
EOF
```

**After (CORRECT):**
```bash
# Use sf_metadata_path helper (defined in R383)
REPORT_PATH=$(sf_metadata_path "$PHASE" "$WAVE" "$EFFORT" "CODE-REVIEW-REPORT" "md")
cat > "$REPORT_PATH" << EOF
# Code Review Report
...
EOF
```

**Pattern 2: Append to Existing File**

**Before (WRONG):**
```bash
cat >> INTEGRATE_WAVE_EFFORTS-REPORT.md << EOF
## New Section
...
EOF
```

**After (CORRECT):**
```bash
# Find latest timestamped report
LATEST_REPORT=$(ls -t .software-factory/phase*/wave*/integration/INTEGRATE_WAVE_EFFORTS-REPORT--*.md 2>/dev/null | head -1)
if [ -n "$LATEST_REPORT" ]; then
    cat >> "$LATEST_REPORT" << EOF
## New Section
...
EOF
fi
```

**Pattern 3: State File Detection**

**Before (WRONG):**
```bash
if [ -f "SPLIT-INVENTORY.md" ]; then
    CURRENT_STATE="SPLIT_WORK"
fi
```

**After (CORRECT):**
```bash
if ls .software-factory/phase*/wave*/*/SPLIT-INVENTORY--*.md >/dev/null 2>&1; then
    CURRENT_STATE="SPLIT_WORK"
fi
```

**Pattern 4: File Listing**

**Before (WRONG):**
```bash
TOTAL_SPLITS=$(ls SPLIT-PLAN-*.md | wc -l)
```

**After (CORRECT):**
```bash
TOTAL_SPLITS=$(ls .software-factory/phase*/wave*/*/SPLIT-PLAN--*.md 2>/dev/null | wc -l)
```

## ENFORCEMENT PROTOCOL

### Development Phase
```bash
# Before starting template changes
git checkout -b fix-template-r383-violations

# After each change
bash tools/validate-templates.sh

# If violations found
# → Fix immediately
# → Run validation again
# → Repeat until clean
```

### Pre-Merge Checklist
```markdown
- [ ] Ran `bash tools/validate-templates.sh` manually
- [ ] Zero violations reported
- [ ] Pre-commit hook is active
- [ ] All template changes use sf_metadata_path()
- [ ] Added comments referencing R383 in changed sections
- [ ] Tested templates in real agent execution
```

### CI/CD Integration
```yaml
# .github/workflows/template-validation.yml
name: R532 Template Validation
on: [pull_request]
jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Run R532 Validation
        run: bash tools/validate-templates.sh
      - name: Comment on PR if violations
        if: failure()
        uses: actions/github-script@v6
        with:
          script: |
            github.rest.issues.createComment({
              issue_number: context.issue.number,
              owner: context.repo.owner,
              repo: context.repo.repo,
              body: '❌ R532 Validation Failed: Template contains R383 violations. See logs for details.'
            })
```

## VIOLATION REMEDIATION

### When Validation Fails

**Step 1: Identify Violation**
```bash
bash tools/validate-templates.sh
# Output shows:
#   File: agent-states/integration/INTEGRATE_WAVE_EFFORTS_BUILD_VALIDATION/rules.md
#   Line: 469
#   Violation: cat > DEMO-STATUS.md
```

**Step 2: Locate in File**
```bash
# Go to file at exact line
vim +469 agent-states/integration/INTEGRATE_WAVE_EFFORTS_BUILD_VALIDATION/rules.md
```

**Step 3: Apply Fix Pattern**
```bash
# Replace hardcoded path with sf_metadata_path call
# OR use .software-factory/phase*/wave*/[agent]/ path with timestamp
```

**Step 4: Verify Fix**
```bash
bash tools/validate-templates.sh
# Should show: ✅ PASSED
```

**Step 5: Test in Agent**
```bash
# Spawn agent that uses the template
# Verify files created in correct location
ls -la .software-factory/phase*/wave*/[agent]/
```

### Mass Remediation Strategy

For projects with many violations:

1. **Run full scan**: `bash tools/validate-templates.sh > violations.txt`
2. **Group by file**: Sort violations by which template file
3. **Fix by priority**:
   - Critical templates (integration, code-reviewer) first
   - Frequently-used states next
   - Experimental/deprecated last
4. **Batch test**: Fix 5-10 at a time, test batch
5. **Commit incrementally**: Don't wait for all fixes

## COMMON VIOLATIONS & FIXES

### Violation 1: Integration Report in Root

**File:** `.claude/agents/integration.md`
**Line:** 562
**Violation:** `cat > INTEGRATE_WAVE_EFFORTS-REPORT.md`

**Fix:**
```bash
REPORT_PATH=$(sf_metadata_path "$PHASE" "$WAVE" "integration" "INTEGRATE_WAVE_EFFORTS-REPORT" "md")
cat > "$REPORT_PATH" << EOF
```

### Violation 2: Split Inventory in Root

**File:** `.claude/agents/code-reviewer.md`
**Line:** 1897
**Violation:** `cat > "SPLIT-INVENTORY.md"`

**Fix:**
```bash
INVENTORY_PATH=$(sf_metadata_path "$PHASE" "$WAVE" "$EFFORT" "SPLIT-INVENTORY" "md")
cat > "$INVENTORY_PATH" << EOF
```

### Violation 3: PR Plan in Root

**File:** `agent-states/pr-ready/orchestrator/PR_PLAN_CREATION/rules.md`
**Line:** 185
**Violation:** `cat > MASTER-PR-PLAN.md`

**Fix:**
```bash
PLAN_PATH=$(sf_metadata_path "X" "Y" "pr-ready" "MASTER-PR-PLAN" "md")
cat > "$PLAN_PATH" << EOF
```

### Violation 4: Demo Status in Root

**File:** `agent-states/integration/INTEGRATE_WAVE_EFFORTS_BUILD_VALIDATION/rules.md`
**Line:** 469
**Violation:** `cat > DEMO-STATUS.md`

**Fix:**
```bash
STATUS_PATH=$(sf_metadata_path "$PHASE" "$WAVE" "integration" "DEMO-STATUS" "md")
cat > "$STATUS_PATH" << EOF
```

## METRICS TO TRACK

**Template Quality Metrics:**
- Total templates scanned
- Violations found per scan
- Violation trend over time
- Most common violation types
- Time to fix violations
- Re-violation rate (same file/line)

**Deployment Metrics:**
- Commits blocked by pre-commit hook
- PRs failed by CI validation
- Manual scan frequency
- Agent execution failures due to missing files

## PROJECT_DONE CRITERIA

✅ **Zero Violations**: `bash tools/validate-templates.sh` reports 0 violations
✅ **Pre-Commit Active**: Hook blocks commits with violations
✅ **CI Passing**: All PRs must pass R532 validation
✅ **Agent Compliance**: No metadata files created in project root during agent execution
✅ **Documentation**: All templates reference R383 in comments

## PENALTIES

| Violation Type | Penalty | Severity |
|----------------|---------|----------|
| Merged template with violation | -100% | BLOCKING |
| Bypassed pre-commit check | -100% | BLOCKING |
| Failed to run validation before merge | -50% | CRITICAL |
| Multiple violations in one PR | -100% | BLOCKING |
| Re-violation of previously fixed template | -75% | SEVERE |

## RELATED RULES

- **R383**: Metadata File Timestamp Requirements (parent rule - defines what R532 validates)
- **R343**: Metadata Directory Standardization (directory structure requirements)
- **R506**: Absolute Prohibition on Pre-Commit Bypass (enforcement mechanism)
- **R054**: Implementation Plan Creation (specific metadata type)
- **R264**: Work Log Tracking Requirements (specific metadata type)

## KEY PRINCIPLE

**"Templates are truth. If templates violate R383, agents will too. Fix templates, fix system."**

This ensures:
- **Prevention at Source**: Stop violations before they happen
- **Consistent Behavior**: All agents follow same patterns
- **Easy Onboarding**: New developers see correct examples
- **Automated Quality**: No manual enforcement needed
- **Zero Merge Conflicts**: Timestamps + structure = clean integrations

## MAINTENANCE

**Weekly:**
- Run full validation scan
- Review any new violations
- Update METADATA_PATTERNS if new file types added

**Monthly:**
- Review validation script effectiveness
- Update allowed patterns if needed
- Check for false positives/negatives

**Per Release:**
- Ensure all templates pass validation
- Document any new template patterns
- Update this rule with learnings

---

**REMEMBER**: R532 is the enforcement arm of R383. Without R532, R383 is just documentation. Together, they ensure zero metadata conflicts forever.

**See Also:**
- `rule-library/R383-metadata-file-timestamp-requirements.md`
- `tools/validate-templates.sh` (implementation)
- `.pre-commit-config.yaml` (automation)
