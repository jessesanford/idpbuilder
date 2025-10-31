# Orchestrator - INIT_HANDOFF State Rules

## Purpose
Complete initialization and transition project to normal Software Factory 2.0 operation.

## Entry Criteria
- All validation checks passed
- All required files generated
- Repository properly configured
- Agent files customized

## Required Actions

### 1. Create Production State File
Generate orchestrator-state-v3.json with:
```json
{
  "current_state": "INIT",
  "current_phase": "1",
  "current_wave": "1",
  "phase_1": {
    "status": "not_started",
    "waves": {}
  },
  "efforts_completed": [],
  "efforts_in_progress": [],
  "project_name": "[from init state]",
  "project_prefix": "[from init state]",
  "initialization_completed": "[timestamp]"
}
```

### 2. Archive Initialization State
- Move init-state-${PROJECT_PREFIX}.json to .completed/
- Add completion timestamp
- Mark as successful

### 3. Generate Summary Report
Create INITIALIZATION-SUMMARY.md with:
- Project name and type
- Technology stack chosen
- Repository configuration
- Files created (with paths)
- Customizations made
- Phase 1 overview from plan

### 4. Display Success Message
```
✅ SOFTWARE FACTORY 2.0 INITIALIZATION COMPLETE!

Project: [name]
Type: [upstream_fork|new_repo|library]
Language: [primary_language]
Framework: [primary_framework]

Created Files:
✓ IMPLEMENTATION-PLAN.md
✓ setup-config.yaml
✓ target-repo-config.yaml
✓ orchestrator-state-v3.json
✓ .claude/CLAUDE.md (customized)
✓ Agent files (with expertise)

Repository:
[Details based on type]

Next Steps:
1. Review IMPLEMENTATION-PLAN.md for Phase 1 details
2. Run: /continue-orchestrating
3. The orchestrator will begin Phase 1, Wave 1 implementation

Total initialization time: [duration]
```

## Exit Criteria
- Production state file created
- Initialization state archived
- Summary report generated
- User informed of next steps
- System ready for /continue-orchestrating

## Terminal State
This is a terminal state for initialization.
Next interaction uses normal state machine.

## Success Metrics
- All files present and valid
- < 30 minutes total time
- No manual intervention needed
- Seamless handoff to production

## Post-Handoff
User runs `/continue-orchestrating` which:
1. Reads orchestrator-state-v3.json
2. Sees state = "INIT"
3. Begins normal Phase 1 execution
4. Uses software-factory-3.0-state-machine.json

## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**


### 🚨 CRITICAL DISTINCTION: AGENT STOPS ≠ FACTORY STOPS 🚨

**TWO INDEPENDENT DECISIONS - DO NOT CONFUSE THEM:**

#### 1. Should Agent Stop Work? (R322 Technical Requirement)
- Agent completes current state
- Agent saves TODOs and commits state
- Agent exits with `exit 0` (preserves context)
- User runs /continue-orchestrating to resume
- **This is NORMAL at checkpoints**

#### 2. Should Factory Continue? (R405 Operational Status)
- Even though agent stopped, can automation proceed?
- TRUE = Healthy completion, automation can continue
- FALSE = Catastrophic failure, must halt everything
- **R322 checkpoints = TRUE (99.9% of cases)**

### THE PATTERN AT R322 CHECKPOINTS

```bash
# 1. Complete state work
echo "✅ State work complete"

# 2. Update state file
jq '.state_machine.current_state = "NEXT_STATE"' orchestrator-state-v3.json > tmp && mv tmp orchestrator-state-v3.json

# 3. Save TODOs
save_todos "R322_CHECKPOINT"

# 4. Factory continues (operational status)
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"

# 5. Agent stops (technical requirement)
exit 0
```

**Both happen together! Agent stops AND factory continues!**

### WHEN TO USE EACH FLAG VALUE

**TRUE (99.9%):**
- ✅ R322 checkpoint reached
- ✅ State work completed successfully
- ✅ Ready for /continue-orchestrating
- ✅ Waiting for user to continue (NORMAL)
- ✅ Plan ready for review (agent done, factory proceeds)

**FALSE (0.1%):**
- ❌ CATASTROPHIC unrecoverable error
- ❌ Data corruption spreading
- ❌ Critical security violation
- ❌ NOT for R322 checkpoints
- ❌ NOT for user review needs
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

