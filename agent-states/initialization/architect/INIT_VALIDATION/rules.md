# Architect - INIT_VALIDATION State Rules

## Purpose
Comprehensive validation that all initialization components are ready for production use.

## Entry Criteria
- All configuration files generated
- IMPLEMENTATION-PLAN.md created
- Agent files customized
- Repository set up

## Required Actions

### 1. File Existence Validation

**Check Required Files**:
```bash
# Core files that MUST exist
FILES_TO_CHECK=(
  "IMPLEMENTATION-PLAN.md"
  "efforts/$PROJECT_PREFIX/configs/setup-config.yaml"
  "efforts/$PROJECT_PREFIX/configs/target-repo-config.yaml"
  "orchestrator-state.json"  # Or will be created in handoff
)

for file in "${FILES_TO_CHECK[@]}"; do
  if [ ! -f "$CLAUDE_PROJECT_DIR/$file" ]; then
    VALIDATION_FAILED: "Missing required file: $file"
  fi
done
```

### 2. Content Completeness Validation

#### IMPLEMENTATION-PLAN.md Validation
```
Required Sections:
- [ ] Project Overview (non-empty)
- [ ] Goals and Objectives (at least 3)
- [ ] Technical Architecture (includes tech stack)
- [ ] Phase 1 defined (with waves)
- [ ] Phase 2 defined (with waves)
- [ ] Phase 3 defined (with waves)
- [ ] Success Criteria (measurable)
- [ ] Risk Mitigation (at least 2 risks)

Phase Structure:
- [ ] Each phase has 2-4 waves
- [ ] Each wave has 3-6 efforts
- [ ] Each effort has clear description
- [ ] Effort naming follows E#.#.# pattern
```

#### setup-config.yaml Validation
```
Required Fields:
- [ ] project.name (not empty)
- [ ] project.prefix (matches PROJECT_PREFIX)
- [ ] technology.primary_language (valid language)
- [ ] technology.build_system (not empty)
- [ ] technology.testing.framework (not empty)
- [ ] architecture.pattern (valid pattern)
- [ ] deployment.environment (not empty)

No Placeholders:
- [ ] No "[...]" values
- [ ] No "TODO" values
- [ ] No null values for required fields
```

#### target-repo-config.yaml Validation
```
For Upstream Fork:
- [ ] repository.type = "upstream_fork"
- [ ] upstream.url (valid git URL)
- [ ] fork.url (valid git URL)
- [ ] paths.target_repo (directory exists)

For New Project:
- [ ] repository.type = "new_project"
- [ ] project.path (directory exists)
- [ ] directories.source (defined)

Common:
- [ ] Valid YAML syntax
- [ ] Consistent paths
```

### 3. Repository Validation

**For Upstream Fork**:
```bash
cd $CLAUDE_PROJECT_DIR/efforts/$PROJECT_PREFIX/target-repo
# Check git status
git status || VALIDATION_FAILED: "Not a git repository"
# Check remotes
git remote -v | grep upstream || VALIDATION_FAILED: "No upstream remote"
# Check branches
git branch | grep -E "phase-1|development" || VALIDATION_FAILED: "Missing branches"
```

**For New Project**:
```bash
cd $CLAUDE_PROJECT_DIR/efforts/$PROJECT_PREFIX/project
# Check git status
git status || VALIDATION_FAILED: "Not a git repository"
# Check initial commit
git log --oneline | head -1 || VALIDATION_FAILED: "No initial commit"
```

### 4. Agent Customization Validation

```bash
# Check sw-engineer.md has expertise
grep -q "Project-Specific Expertise" $CLAUDE_PROJECT_DIR/.claude/agents/sw-engineer.md || \
  VALIDATION_FAILED: "SW Engineer not customized"

# Check code-reviewer.md has expertise
grep -q "Project-Specific Review Points" $CLAUDE_PROJECT_DIR/.claude/agents/code-reviewer.md || \
  VALIDATION_FAILED: "Code Reviewer not customized"

# Verify customizations are minimal
# Count lines in expertise sections, should be < 10
```

### 5. Create Validation Report

Generate INIT-VALIDATION-REPORT.md:
```markdown
# Initialization Validation Report

## Project: [PROJECT_NAME]
## Status: [PASSED|FAILED]

### File Checklist
- [✓] IMPLEMENTATION-PLAN.md exists
- [✓] setup-config.yaml exists
- [✓] target-repo-config.yaml exists
- [✓] Repository initialized

### Content Validation
- [✓] Plan has all required sections
- [✓] Config files fully populated
- [✓] No placeholder values
- [✓] Git repository functional

### Agent Customization
- [✓] SW Engineer has [LANGUAGE] expertise
- [✓] Code Reviewer has review criteria
- [✓] Customizations are minimal

### Readiness Assessment
- [✓] Ready for /continue-orchestrating
- [✓] Can begin Phase 1 implementation
- [✓] All dependencies resolved

### Issues Found
[None | List any issues]

### Validation Time: [TIMESTAMP]
```

### 6. Determine Next State

```
IF validation_passed:
  next_state = "INIT_HANDOFF"
  status = "ready_for_production"
ELSE:
  next_state = "INIT_ERROR_RECOVERY"
  status = "requires_fixes"
  document_specific_issues()
```

### 7. Update State File
```json
"validation": {
  "status": "[passed|failed]",
  "timestamp": "[ISO_TIME]",
  "files_checked": [list],
  "issues_found": [list],
  "ready_for_handoff": true/false
}
```

## Validation Criteria Summary

### MUST PASS (Blocking):
- All required files exist
- IMPLEMENTATION-PLAN.md has phases/waves/efforts
- Config files have no placeholders
- Git repository accessible
- Agent files customized

### SHOULD PASS (Warning):
- Proper effort sizing (<700 lines estimate)
- Reasonable phase distribution
- Complete risk mitigation
- Comprehensive success criteria

### NICE TO HAVE:
- Examples in documentation
- Additional reference materials
- Extended agent expertise

## Exit Criteria
- Validation complete
- Report generated
- All blocking issues resolved
- Ready for handoff or error recovery

## Transitions
**SUCCESS**: → INIT_HANDOFF
**FAILURE**: → INIT_ERROR_RECOVERY

## Time Guidance
- Validation checks: 2-3 minutes
- Should be thorough but quick
- Focus on blocking issues first

## Common Validation Failures

1. **Missing Sections in Plan**
   - Fix: Return to INIT_SYNTHESIZE_PLAN
   - Add missing sections

2. **Placeholder Values in Configs**
   - Fix: Return to INIT_GENERATE_CONFIGS
   - Fill in actual values

3. **Git Repository Issues**
   - Fix: Return to repository setup state
   - Reinitialize git

4. **Agent Files Not Customized**
   - Fix: Return to INIT_CUSTOMIZE_AGENTS
   - Add expertise sections