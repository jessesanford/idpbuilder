# 🔴🔴🔴 SUPREME LAW: RULE R372 - Effort Theme Enforcement (BLOCKING)

## Rule Details
- **Rule ID**: R372
- **Rule Name**: Effort Theme Enforcement
- **Criticality**: 🚨🚨🚨 BLOCKING
- **Category**: SUPREME LAW - Theme Discipline
- **Created**: 2025-01-21
- **Updated**: 2025-01-21

## 🔴 SUPREME DIRECTIVE

**EACH EFFORT MUST HAVE ONE CLEAR THEME - NO KITCHEN SINK BRANCHES**

Every effort must embody a SINGLE, FOCUSED theme or spirit. Mixing unrelated work in effort branches is a catastrophic architectural violation that destroys traceability and creates unmergeable monsters.

## 🚨 THEME VIOLATION EXAMPLES

### ❌ CATASTROPHIC FAILURES DETECTED:
```yaml
registry-helpers branch:
  Intended Theme: "Registry helper functions"
  Actual Content:
    - ✅ Registry helpers (correct)
    - ❌ Kind cluster setup (DIFFERENT THEME)
    - ❌ Test infrastructure (DIFFERENT THEME)
    - ❌ Build system changes (DIFFERENT THEME)
  Result: 70 files of mixed themes (FAILURE)

gitea-client-split-001:
  Intended Theme: "Gitea API client"
  Actual Content:
    - ✅ Gitea client (correct)
    - ❌ DevContainer setup (DIFFERENT THEME)
    - ❌ CI/CD pipeline (DIFFERENT THEME)
    - ❌ Documentation overhaul (DIFFERENT THEME)
  Result: 124 files, multiple themes (CATASTROPHIC)
```

## ⚠️ MANDATORY THEME REQUIREMENTS

### 1. THEME DECLARATION
Every effort plan MUST declare its theme:
```markdown
## EFFORT THEME
**Primary Theme**: Registry API Integration
**Theme Boundary**: ONLY code that directly implements registry API calls
**Theme Spirit**: Clean, minimal API client with no side concerns

## THEME VIOLATIONS (FORBIDDEN):
- Infrastructure changes
- Build system modifications
- Test framework setup
- Documentation updates (unless API docs)
- Utility functions unrelated to API
```

### 2. THEME VALIDATION MATRIX
```markdown
| Change Type | Aligns with Theme? | Decision |
|-------------|-------------------|----------|
| API client code | ✅ YES | ALLOW |
| API types/models | ✅ YES | ALLOW |
| Unit tests for API | ✅ YES | ALLOW |
| Makefile changes | ❌ NO | REJECT |
| Docker configs | ❌ NO | REJECT |
| Logging framework | ❌ NO | REJECT |
```

### 3. THEME COHERENCE TEST
```bash
# Every file must answer YES to:
theme_coherence_test() {
    echo "1. Does this file directly support the theme?"
    echo "2. Would removing this file break the theme?"
    echo "3. Could this file be in a different effort?"

    # If ANY answer is NO, file violates theme
}
```

## 🔴 THEME ENFORCEMENT RULES

### PLANNING PHASE
```yaml
Orchestrator Requirements:
  - Define ONE clear theme per effort
  - Theme must be 1-2 sentences MAX
  - Theme must be specific and actionable
  - NO compound themes (X AND Y)
  - NO vague themes ("improvements", "refactoring")
```

### IMPLEMENTATION PHASE
```yaml
SW-Engineer Requirements:
  - Read theme before ANY code
  - Ask: "Does this support the theme?"
  - If unclear: DO NOT ADD
  - Re-read theme every 5 files
  - Stop if drifting from theme
```

### REVIEW PHASE
```yaml
Code-Reviewer Requirements:
  - First check: Theme coherence
  - Flag ANY file outside theme
  - Count theme violations
  - Score theme purity (must be >95%)
  - Reject if multiple themes detected
```

## 🚨 THEME CATEGORIES

### ALLOWED THEMES (Focused):
```markdown
✅ "Implement Gitea API client"
✅ "Add authentication layer"
✅ "Create database models"
✅ "Build CLI interface"
✅ "Add unit tests for package X"
```

### FORBIDDEN THEMES (Too Broad):
```markdown
❌ "Improve the system"
❌ "Add features and fix bugs"
❌ "Refactor and optimize"
❌ "Set up development environment"
❌ "General maintenance"
```

### FORBIDDEN THEME MIXING:
```markdown
❌ API client + Build system
❌ Business logic + Infrastructure
❌ Feature + Documentation overhaul
❌ Core code + Test framework setup
❌ Implementation + DevContainer configs
```

## 📊 THEME PURITY METRICS

### Measurement Formula:
```bash
theme_purity_score() {
    local theme_files=$(grep -c "supports theme" FILES.md)
    local total_files=$(wc -l < FILES.md)
    local purity=$((theme_files * 100 / total_files))

    if [ $purity -lt 95 ]; then
        echo "❌ THEME VIOLATION: Purity only $purity%"
        return 1
    fi
    echo "✅ Theme purity: $purity%"
}
```

### Required Thresholds:
- **95%+**: Minimum acceptable
- **98%+**: Good theme discipline
- **100%**: Perfect theme focus

## 🔴 KITCHEN SINK PREVENTION

### Signs of Kitchen Sink:
1. More than 3 different package directories modified
2. Both src/ and build/ changes
3. Both code and infrastructure changes
4. Multiple unrelated features
5. "While I'm here" additions

### Prevention Protocol:
```bash
prevent_kitchen_sink() {
    # Check package spread
    local package_count=$(git diff --name-only |
                         cut -d'/' -f1-2 |
                         sort -u |
                         wc -l)

    if [ $package_count -gt 3 ]; then
        echo "🚨 KITCHEN SINK ALERT: Too many packages!"
        echo "Split into focused efforts!"
        return 1
    fi
}
```

## 🚨 VIOLATION PENALTIES

### IMMEDIATE FAILURE (-100%):
- Effort with multiple themes
- Kitchen sink branch (>3 concerns)
- Theme purity <95%
- No theme declaration

### SEVERE PENALTIES (-50%):
- Vague theme definition
- Theme drift during implementation
- Missing theme validation

## 📝 REQUIRED ARTIFACTS

### EFFORT-THEME-DECLARATION.md:
```markdown
# Effort Theme Declaration

## PRIMARY THEME
"Implement Gitea API client for registry operations"

## THEME BOUNDARIES
- IN: API client code, request/response types
- OUT: Infrastructure, build system, documentation

## THEME VALIDATION CHECKLIST
- [ ] Single, focused concern
- [ ] Clear boundaries defined
- [ ] No mixing with other themes
- [ ] <20 files estimated
```

### THEME-COHERENCE-REPORT.md:
```markdown
# Theme Coherence Report

## Theme: "Gitea API Client"

## Files Analyzed: 15
- On-theme: 15 (100%)
- Off-theme: 0 (0%)

## Coherence Score: 100% ✅

## File Breakdown:
| File | Theme Alignment | Justification |
|------|----------------|---------------|
| client.go | ✅ ON-THEME | Core API client |
| types.go | ✅ ON-THEME | API types |
```

## 🔴 EMERGENCY PROTOCOL

If theme violation detected:
1. **STOP ALL WORK**
2. Identify theme drift point
3. Revert all off-theme changes
4. Create new effort for off-theme work
5. Document in THEME-VIOLATION-INCIDENT.md

## ✅ GOOD EXAMPLE

```markdown
Effort: "Implement Registry Authentication"
Theme: "Add OAuth2 authentication to registry client"
Files:
  - pkg/auth/oauth2.go
  - pkg/auth/token.go
  - pkg/auth/oauth2_test.go
Result: 3 files, 100% theme coherence
```

## ❌ BAD EXAMPLE

```markdown
Effort: "Improve Registry"
Theme: [VAGUE - VIOLATION]
Files:
  - pkg/auth/oauth2.go (authentication)
  - Makefile (build system)
  - .github/workflows/ci.yml (CI/CD)
  - docs/README.md (documentation)
  - pkg/cache/cache.go (caching)
Result: 5+ themes mixed (CATASTROPHIC)
```

---

**REMEMBER**: One effort, one theme, one purpose. Multi-theme branches are architectural cancer that metastasize into unmergeable disasters. MAINTAIN THEME DISCIPLINE!