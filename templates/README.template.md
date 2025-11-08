# {{PROJECT_NAME}}

**Software Factory 3.0 Project**

This project was created using the [Software Factory 3.0](https://github.com/jessesanford/software-factory-template) template - an AI-native development system designed for building high-quality software with autonomous agent orchestration.

---

## 🎯 Project Overview

**Status**: Initialization Complete
**Created**: {{CREATION_DATE}}
**Software Factory Version**: 3.0

---

## 🚀 Getting Started

### What is Software Factory 3.0?

Software Factory 3.0 is a development framework that coordinates multiple AI agents to build, test, review, and integrate code systematically. It provides:

- **Orchestrator Agent**: Manages project phases, waves, and efforts
- **Software Engineer Agents**: Implement features in isolated branches
- **Code Reviewer Agents**: Perform systematic code reviews and track bugs
- **Architect Agent**: Validates architectural integrity across iterations

### Your Next Steps

1. **Review Project Plan** (if created)
   ```bash
   cat PROJECT-IMPLEMENTATION-PLAN.md
   ```
   The implementation plan defines all phases, waves, and efforts for your project.

2. **Initialize First Phase** (using Claude Code with `/init-software-factory` command)
   ```bash
   # The Orchestrator will:
   # - Parse your implementation plan
   # - Set up effort infrastructure (branches, working copies)
   # - Begin coordinating agent work
   ```

3. **Monitor Progress**
   ```bash
   # Check current state
   cat orchestrator-state-v3.json

   # View tracked bugs
   cat bug-tracking.json

   # Check active integrations
   cat integration-containers.json
   ```

4. **Continue Development** (using `/continue-software-factory` command)
   ```bash
   # The Orchestrator will:
   # - Resume from current state
   # - Spawn appropriate agents for next work
   # - Manage iteration convergence
   # - Track all state changes atomically
   ```

---

## 📁 Project Structure

```
{{PROJECT_NAME}}/
├── README.md                          # This file
├── PROJECT-IMPLEMENTATION-PLAN.md     # Your implementation plan
├── orchestrator-state-v3.json         # Current project state
├── bug-tracking.json                  # Bug tracking database
├── integration-containers.json        # Integration iteration tracking
├── fix-cascade-state.json             # Fix cascade tracking (created when needed)
├── state-machines/                    # State machine definitions
│   └── software-factory-3.0-state-machine.json
├── agent-states/                      # Agent-specific state rules
│   ├── software-factory/
│   │   ├── orchestrator/
│   │   ├── sw-engineer/
│   │   ├── code-reviewer/
│   │   └── architect/
│   └── state-manager/
├── rule-library/                      # Development rules (R###)
├── .claude/                           # Claude Code configuration
│   ├── agents/                        # Agent definitions
│   └── commands/                      # Slash commands
├── tools/                             # Development tools
│   ├── atomic-state-update.sh         # Atomic 4-file state updates
│   ├── validate-state-file.sh         # State file validation
│   └── line-counter.sh                # Code size measurement
└── utilities/                         # Helper scripts
    └── init-software-factory.sh       # This initialization script
```

---

## 🔐 Git Hooks & Quality Gates

### Pre-Commit Hook (R506)

This project enforces **R506: Absolute Prohibition on Pre-Commit Bypass**. The pre-commit hook:

- ✅ Validates all 4 state files before every commit
- ✅ Prevents invalid state from entering git history
- ✅ Ensures atomic 4-file updates (R288)
- ⛔ **NEVER use `git commit --no-verify`** - this is a critical violation

If your commit is blocked:
1. Fix the validation errors (don't bypass!)
2. Use `tools/validate-state-file.sh` to debug issues
3. The hook is your safety net - respect it

---

## 📊 State Files (SF 3.0)

### orchestrator-state-v3.json
The primary state file containing:
- `state_machine`: Current state, transitions, history
- `project_progression`: Phases, waves, efforts
- `references`: Links to other state files

### bug-tracking.json
Centralized bug tracking database:
- Bugs discovered during code reviews
- Cross-container bug relationships
- Fix cascade triggers

### integration-containers.json
Tracks integration iterations:
- Active integration containers (wave/phase/project level)
- Iteration counts and convergence metrics
- Upstream bug history

### fix-cascade-state.json
Created dynamically when cross-container bugs are detected:
- Cascade activation triggers
- Multi-container fix coordination
- Escalation thresholds

**All 4 files are updated atomically** using `tools/atomic-state-update.sh` per R288.

---

## 🛠️ Common Commands

### Check Current State
```bash
# View current orchestrator state
jq '.state_machine.current_state' orchestrator-state-v3.json

# Count open bugs
jq '.bugs | length' bug-tracking.json

# List active integrations
jq '.active_integrations[] | .container_id' integration-containers.json
```

### Validate State Files
```bash
# Validate all state files
bash tools/validate-state-file.sh orchestrator-state-v3.json
bash tools/validate-state-file.sh bug-tracking.json
bash tools/validate-state-file.sh integration-containers.json

# Validate against schemas
jq . orchestrator-state-v3.json >/dev/null && echo "Valid JSON"
```

### Measure Code Size
```bash
# Count lines in current branch (excludes generated code)
bash tools/line-counter.sh
```

---

## 🤖 Agent Coordination

### Orchestrator Agent
**Role**: Project manager and coordination hub

**Key States**:
- `INIT`: Initialize project structure
- `PLANNING`: Parse implementation plan
- `CREATE_NEXT_INFRASTRUCTURE`: Create branches and working copies
- `SPAWN_SW_ENGINEERS`: Deploy SWE agents for parallel work
- `MONITORING_SWE_PROGRESS`: Track agent work
- `INTEGRATE_WAVE_EFFORTS`: Merge completed efforts
- `REVIEW_WAVE_INTEGRATION`: Deploy Code Reviewer
- `FIX_WAVE_UPSTREAM_BUGS`: Address review findings

### Software Engineer Agent
**Role**: Implementation specialist

**Key States**:
- `IMPLEMENTATION`: Write code in isolated branch
- `MEASURE_SIZE`: Check if effort exceeds size limits
- `SPLIT_IMPLEMENTATION`: Break large efforts into smaller pieces
- `REQUEST_REVIEW`: Complete work and request review

### Code Reviewer Agent
**Role**: Quality assurance specialist

**Key States**:
- `EFFORT_PLAN_CREATION`: Create review checklist
- `CODE_REVIEW`: Systematic code review
- `CREATE_SPLIT_PLAN`: Plan effort splits when too large
- `CREATE_FIX_PLAN`: Organize bugs into actionable fixes

### Architect Agent
**Role**: Architectural integrity validator

**Key States**:
- `REVIEW_WAVE_ARCHITECTURE`: Validate wave-level architecture
- `PHASE_ASSESSMENT`: Validate phase-level architecture
- `INTEGRATE_WAVE_EFFORTS_REVIEW`: Review integration patterns

---

## 📚 Documentation

- **Architecture**: `docs/SOFTWARE-FACTORY-3.0-ARCHITECTURE.md`
- **Rule Migration Plan**: `docs/RULE-MIGRATION-PLAN-SF3.md`
- **Execution Checklist**: `docs/SF-3.0-COMPREHENSIVE-EXECUTION-CHECKLIST.md`
- **Rule Library**: `rule-library/R*.md` (475+ rules)

---

## ⚡ Quick Reference

| Task | Command |
|------|---------|
| Start new project | `/init-software-factory` |
| Continue work | `/continue-software-factory` |
| Check project state | `cat orchestrator-state-v3.json` |
| View bugs | `cat bug-tracking.json` |
| Measure code size | `bash tools/line-counter.sh` |
| Validate state | `bash tools/validate-state-file.sh <file>` |

---

## 🆘 Troubleshooting

### Commit Blocked by Pre-Commit Hook
**Problem**: Git commit fails with validation errors

**Solution**:
1. Don't use `--no-verify` (R506 violation!)
2. Run `bash tools/validate-state-file.sh <file>` to identify issues
3. Fix the validation errors
4. Commit again

### State Files Out of Sync
**Problem**: State files inconsistent after manual edits

**Solution**:
1. Use `tools/atomic-state-update.sh` for all state updates
2. Never manually edit state files outside of atomic updates
3. If corrupted, restore from git history

### Agent Stuck or Confused
**Problem**: Orchestrator or agent not progressing

**Solution**:
1. Check `orchestrator-state-v3.json` for current state
2. Verify state machine allows transition
3. Check `bug-tracking.json` for blocking bugs
4. Review agent state rules in `agent-states/`

---

## 📝 License

{{LICENSE_INFO}}

---

**Generated by Software Factory 3.0**
_Building better software through systematic AI agent orchestration_
