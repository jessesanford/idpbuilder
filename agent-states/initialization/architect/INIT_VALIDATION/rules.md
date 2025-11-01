# Architect - INIT_VALIDATION State Rules

## 🔴🔴🔴 ABSOLUTE REQUIREMENT: STATE MANAGER CONSULTATION 🔴🔴🔴

**SUPREME LAW - R517 - UNIVERSAL STATE MANAGER CONSULTATION LAW**

**BEFORE EXITING THIS STATE, YOU MUST:**

1. **SPAWN STATE MANAGER FOR SHUTDOWN_CONSULTATION** (MANDATORY - NO EXCEPTIONS)
2. **NEVER UPDATE STATE FILES DIRECTLY** (orchestrator-state-v3.json, bug-tracking.json, etc.)
3. **NEVER COMMIT STATE FILES YOURSELF** (State Manager does this atomically)

**FAILURE TO CONSULT STATE MANAGER = IMMEDIATE SYSTEM HALT (-100% GRADE)**

### Enforcement Mechanism

If you attempt to exit this state without spawning State Manager:
- ❌ Pre-commit hooks will REJECT your commit
- ❌ Validation tools will FAIL the build
- ❌ Grading system will assign -100% penalty
- ❌ System will transition to ERROR_RECOVERY

### Required Pattern (COPY THIS EXACTLY)

```bash
# At end of state work, BEFORE any state file updates:

echo "🔴 MANDATORY: Spawning State Manager for SHUTDOWN_CONSULTATION"

# Spawn State Manager (REQUIRED - NOT OPTIONAL)
# Task: state-manager
# State: SHUTDOWN_CONSULTATION
# Current State: [YOUR_CURRENT_STATE]
# Proposed Next State: [YOUR_PROPOSED_NEXT_STATE]
# Work Summary: [SUMMARY_OF_WORK_COMPLETED]

# State Manager will:
# 1. Validate proposed transition against state machine
# 2. Update all 4 state files atomically
# 3. Commit with [R288] tag
# 4. Return REQUIRED next state (may differ from proposal)

# Wait for State Manager response
# Follow State Manager's directive (REQUIRED next state)
# DO NOT proceed until State Manager confirms
```

**YOU MUST NEVER:**
- ❌ Update orchestrator-state-v3.json yourself
- ❌ Update bug-tracking.json yourself
- ❌ Update integration-containers.json yourself
- ❌ Use `jq` to modify state files
- ❌ Use `sed/awk` to modify state files
- ❌ Set `validated_by: "orchestrator"` (must be "state-manager")
- ❌ Skip State Manager consultation "just this once"
- ❌ Think "I'll validate it manually"

**ONLY State Manager may update state files. This is NON-NEGOTIABLE.**

See: `rule-library/R517-universal-state-manager-consultation-law.md`
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
  "orchestrator-state-v3.json"  # Or will be created in handoff
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
**PROJECT_DONE**: → INIT_HANDOFF
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

## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

### YOU MUST OUTPUT THE CONTINUATION FLAG AS YOUR LAST ACTION

**EVERY STATE COMPLETION MUST END WITH EXACTLY ONE OF:**
```bash
# If state completed successfully and factory should continue:
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"

# If error/block/manual review needed:
echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
```

### CRITICAL REQUIREMENTS:
1. **ABSOLUTE LAST OUTPUT**: This MUST be the very last line of output before state completion
2. **NO TEXT AFTER**: No explanations, summaries, or any text after the flag
3. **EXACTLY AS SHOWN**: Use exact format - no variations like CONTINUE-ORCHESTRATING
4. **ALWAYS REQUIRED**: Every single state must output this flag
5. **GREPPABLE**: Must be on its own line for automation parsing

### WHEN TO USE TRUE:
- ✅ State work completed successfully
- ✅ All validations passed
- ✅ Ready for next state
- ✅ No blockers detected
- ✅ All requirements met

### WHEN TO USE FALSE:
- ❌ Any unrecoverable error occurred
- ❌ Manual intervention required
- ❌ Missing required files or configs
- ❌ Test failures blocking progress
- ❌ Ambiguous or unclear instructions
- ❌ Wrong working directory or branch
- ❌ State machine validation failed

### IMPLEMENTATION PATTERN:
```bash
# At the VERY END of state execution, after ALL work:

# Determine success/failure
if [[ "$STATE_COMPLETED_PROJECT_DONEFULLY" == "true" ]]; then
    echo "✅ State work completed successfully"
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
else
    echo "❌ State encountered issues requiring intervention"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
fi

# DO NOT OUTPUT ANYTHING AFTER THE FLAG!
```

### GRADING IMPACT:
- **Missing flag**: -100% AUTOMATIC FAILURE
- **Text after flag**: -50% penalty
- **Wrong format**: -100% AUTOMATIC FAILURE
- **Multiple flags**: -50% penalty

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

