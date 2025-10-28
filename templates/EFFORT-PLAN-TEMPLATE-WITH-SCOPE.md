# EFFORT PLAN - [EFFORT NAME]

## 🎯 EFFORT THEME DECLARATION (R372)

**PRIMARY THEME**: [One clear, focused theme - 1 sentence]
**THEME BOUNDARY**: [What's included in this theme]
**THEME SPIRIT**: [The core purpose/intent]

## 📋 SCOPE DEFINITION (R371)

### FILES IN SCOPE (ALLOWED)

**EXACT files to create/modify:**
```
pkg/[module]/[file1].go          # [Purpose]
pkg/[module]/[file2].go          # [Purpose]
pkg/[module]/[file1]_test.go     # [Unit tests]
```

**Maximum files**: [Number, should be <20]
**Estimated lines**: [Total new lines, must be <800]

### 🚫 OUT OF SCOPE (FORBIDDEN)

**DO NOT MODIFY:**
- Build system files (Makefile, go.mod, package.json)
- Infrastructure files (.devcontainer/, docker/)
- CI/CD files (.github/, .gitlab-ci.yml)
- Documentation (unless this IS a docs effort)
- Other packages not listed above
- Test framework/infrastructure
- Configuration files (unless listed above)

## 📊 SCOPE TRACEABILITY MATRIX

| File/Package | Requirement | Justification |
|--------------|-------------|---------------|
| pkg/gitea/client.go | "Create Gitea client" | Core implementation |
| pkg/gitea/types.go | "Define API types" | Required for client |
| pkg/gitea/client_test.go | "Unit tests" | Test coverage |

## ✅ ACCEPTANCE CRITERIA

1. [ ] All files in scope created/modified
2. [ ] NO files outside scope touched
3. [ ] Theme coherence maintained (>95%)
4. [ ] Unit tests pass
5. [ ] Line count <800
6. [ ] No stubs/mocks/TODOs (R355)

## 🔍 VALIDATION CHECKLIST

**BEFORE STARTING:**
- [ ] Theme is single and focused
- [ ] File list is explicit and complete
- [ ] OUT OF SCOPE section reviewed
- [ ] <20 files total
- [ ] No mixed concerns

**BEFORE COMMITTING:**
- [ ] Run: `tools/validate-effort-scope.sh`
- [ ] All changes trace to requirements
- [ ] No scope creep occurred
- [ ] Theme purity verified

## 📝 IMPLEMENTATION NOTES

[Any specific implementation guidance, patterns to follow, libraries to use]

## ⚠️ CRITICAL REMINDERS

- **R371**: The file list above is LAW - no additions allowed
- **R372**: Maintain single theme - no kitchen sinks
- **R355**: Production-ready code only - no stubs
- **R359**: Never delete existing code for size
- **R362**: Follow approved architecture exactly

---

**Generated**: [Date]
**Effort State**: PLANNED
**Review Status**: PENDING