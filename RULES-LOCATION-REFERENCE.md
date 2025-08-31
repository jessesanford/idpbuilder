# SOFTWARE FACTORY 2.0 - RULES LOCATION REFERENCE

This document maps all rules by agent type and location for quick reference.

## 🔍 QUICK LOOKUP

### Rule Library (Numbered Rules: R001-R999)
**Location:** `/home/vscode/software-factory-template/rule-library/`

**Critical Supreme Laws:**
- R208: Orchestrator spawn directory protocol (SUPREME LAW #2)
- R234: Mandatory state traversal (SUPREME LAW)
- R235: Mandatory preflight verification (SUPREME LAW)

**TODO Persistence Rules (Critical for all agents):**
- R187: TODO save triggers
- R188: TODO save frequency  
- R189: TODO commit protocol
- R190: TODO recovery verification

**Core Operation Rules:**
- R002: Agent acknowledgment protocol
- R003: Performance grading system
- R006: Orchestrator never writes code
- R007: Size limit compliance
- R203: State-aware agent startup

---

## 📋 BY AGENT TYPE

### 🎼 ORCHESTRATOR AGENT

#### Main Configuration
- **File:** `/home/vscode/software-factory-template/.claude/agents/orchestrator.md`
- **Description:** Core orchestrator configuration and supreme laws

#### State-Specific Rules
**Location:** `/home/vscode/software-factory-template/agent-states/orchestrator/[STATE]/rules.md`

**Available States:**
- `INIT/rules.md` - Initialization requirements and immediate actions
- `PLANNING/rules.md` - Phase and wave planning protocols
- `SETUP_EFFORT_INFRASTRUCTURE/rules.md` - Workspace setup requirements
- `SPAWN_AGENTS/rules.md` - Agent spawning protocols
- `MONITOR/rules.md` - Progress monitoring and tracking
- `WAVE_COMPLETE/rules.md` - Wave completion verification
- `INTEGRATION/rules.md` - Integration management
- `ERROR_RECOVERY/rules.md` - Error handling and recovery
- `PHASE_INTEGRATION/rules.md` - Phase-level integration
- `SPAWN_ARCHITECT_PHASE_ASSESSMENT/rules.md` - Architect spawning for assessment
- `SPAWN_CODE_REVIEWERS_EFFORT_PLANNING/rules.md` - Code reviewer spawning
- `ANALYZE_IMPLEMENTATION_PARALLELIZATION/rules.md` - Parallelization analysis
- `ANALYZE_CODE_REVIEWER_PARALLELIZATION/rules.md` - Reviewer parallelization
- Plus 15+ additional specialized states

#### Special Rule Files
- `/home/vscode/software-factory-template/agent-states/orchestrator/R208-SUPREME-LAW-ENFORCEMENT.md` - Supreme law R208 details
- `/home/vscode/software-factory-template/agent-states/orchestrator/SUPREME_LAW_R217.md` - Supreme law R217 details
- `/home/vscode/software-factory-template/agent-states/orchestrator/CRITICAL_RULES.md` - Critical rule summary

---

### 👨‍💻 SOFTWARE ENGINEER AGENT

#### Main Configuration
- **File:** `/home/vscode/software-factory-template/.claude/agents/sw-engineer.md`
- **Description:** Software engineering agent configuration

#### State-Specific Rules
**Location:** `/home/vscode/software-factory-template/agent-states/sw-engineer/[STATE]/rules.md`

**Available States:**
- `INIT/rules.md` - Initialization and workspace setup
- `IMPLEMENTATION/rules.md` - Core implementation rules
- `SPLIT_IMPLEMENTATION/rules.md` - Split implementation handling
- `MEASURE_SIZE/rules.md` - Size measurement protocols
- `FIX_ISSUES/rules.md` - Issue resolution procedures
- `TEST_WRITING/rules.md` - Test creation requirements
- `BLOCKED/rules.md` - Blocking condition handling
- `COMPLETED/rules.md` - Completion verification
- `REQUEST_REVIEW/rules.md` - Review request procedures

---

### 🔍 CODE REVIEWER AGENT

#### Main Configuration
- **File:** `/home/vscode/software-factory-template/.claude/agents/code-reviewer.md`
- **Description:** Code review agent configuration

#### State-Specific Rules
**Location:** `/home/vscode/software-factory-template/agent-states/code-reviewer/[STATE]/rules.md`

**Available States:**
- `INIT/rules.md` - Initialization procedures
- `CODE_REVIEW/rules.md` - Code review protocols
- `CREATE_SPLIT_PLAN/rules.md` - Split planning creation
- `SPLIT_REVIEW/rules.md` - Split-specific review
- `EFFORT_PLANNING/rules.md` - Effort planning procedures
- `EFFORT_PLAN_CREATION/rules.md` - Plan creation protocols
- `VALIDATION/rules.md` - Validation procedures
- `SPLIT_PLANNING/rules.md` - Split planning protocols
- `PHASE_IMPLEMENTATION_PLANNING/rules.md` - Phase implementation planning
- `WAVE_IMPLEMENTATION_PLANNING/rules.md` - Wave implementation planning
- `PHASE_MERGE_PLANNING/rules.md` - Phase merge planning
- `WAVE_MERGE_PLANNING/rules.md` - Wave merge planning
- `WAVE_DIRECTORY_ACKNOWLEDGMENT/rules.md` - Directory acknowledgment
- `COMPLETED/rules.md` - Completion procedures

---

### 🏗️ ARCHITECT AGENT

#### Main Configuration
- **File:** `/home/vscode/software-factory-template/.claude/agents/architect.md`
- **Description:** Architecture agent configuration

#### State-Specific Rules
**Location:** `/home/vscode/software-factory-template/agent-states/architect/[STATE]/rules.md`

**Available States:**
- `INIT/rules.md` - Initialization procedures
- `WAVE_REVIEW/rules.md` - Wave review protocols
- `PHASE_ASSESSMENT/rules.md` - Phase assessment procedures
- `INTEGRATION_REVIEW/rules.md` - Integration review protocols
- `ARCHITECTURE_AUDIT/rules.md` - Architecture auditing
- `ARCHITECTURE_VALIDATION/rules.md` - Architecture validation
- `PHASE_ARCHITECTURE_PLANNING/rules.md` - Phase architecture planning
- `WAVE_ARCHITECTURE_PLANNING/rules.md` - Wave architecture planning
- `PHASE_DIRECTORY_ACKNOWLEDGMENT/rules.md` - Directory acknowledgment
- `DECISION/rules.md` - Decision making protocols

---

### 🔗 INTEGRATION AGENT

#### Main Configuration
- **File:** `/home/vscode/software-factory-template/.claude/agents/integration.md`
- **Description:** Integration agent configuration

#### State-Specific Rules
**Location:** `/home/vscode/software-factory-template/agent-states/integration/[STATE]/rules.md`

**Available States:**
- `INIT/rules.md` - Initialization procedures
- `PLANNING/rules.md` - Integration planning
- `MERGING/rules.md` - Merge operation protocols
- `TESTING/rules.md` - Integration testing
- `REPORTING/rules.md` - Integration reporting
- `COMPLETED/rules.md` - Completion procedures

#### Special Integration Rules
- `/home/vscode/software-factory-template/rule-library/agents/integration/core-rules.md` - Core integration rules

---

## 🗂️ RULE LIBRARY REFERENCE

### Universal Rules (All Agents)
**Location:** `/home/vscode/software-factory-template/rule-library/`

#### Supreme Laws & Critical Rules
- `R208-orchestrator-spawn-directory-protocol.md` - Directory spawn protocol (SUPREME LAW #2)
- `R234-mandatory-state-traversal-supreme-law.md` - State traversal (SUPREME LAW)  
- `R235-MANDATORY-PREFLIGHT-VERIFICATION-SUPREME-LAW.md` - Preflight verification (SUPREME LAW)
- `R251-REPOSITORY-SEPARATION-LAW.md` - Repository separation law

#### TODO Persistence (R187-R190)
- `R187-todo-save-triggers.md` - Save trigger requirements
- `R188-todo-save-frequency.md` - Save frequency requirements
- `R189-todo-commit-protocol.md` - Commit protocol requirements
- `R190-todo-recovery-verification.md` - Recovery verification requirements

#### Core Operation Rules
- `R002-agent-acknowledgment.md` - Agent acknowledgment protocol
- `R003-performance-grading.md` - Performance grading system
- `R006-orchestrator-never-writes-code.md` - Orchestrator role boundaries
- `R007-size-limit-compliance.md` - Size limit enforcement
- `R008-monitoring-frequency.md` - Monitoring requirements
- `R009-integration-branch-creation.md` - Branch creation protocols

#### Startup & State Management
- `R203-state-aware-agent-startup.md` - Agent startup protocol
- `R206-state-machine-transition-validation.md` - Transition validation
- `R217-post-transition-rule-reacknowledgment.md` - Post-transition acknowledgment
- `R252-MANDATORY-STATE-FILE-UPDATES.md` - State file update requirements
- `R253-MANDATORY-STATE-FILE-COMMIT-PUSH.md` - State file commit requirements

#### Size & Split Management
- `R198-line-counter-usage.md` - Line counter requirements
- `R199-single-reviewer-split-planning.md` - Split planning protocols
- `R200-measure-only-changeset.md` - Changeset measurement
- `R201-line-counter-location.md` - Line counter location requirements
- `R202-single-agent-per-split.md` - Split agent assignment
- `R204-orchestrator-split-infrastructure.md` - Split infrastructure
- `R205-sw-engineer-split-navigation.md` - Split navigation protocols
- `R207-split-boundary-validation.md` - Split boundary validation

#### Workspace & Repository Management
- `R181-orchestrator-workspace-setup.md` - Workspace setup protocols
- `R182-verify-git-repository.md` - Git repository verification  
- `R184-verify-git-branch.md` - Git branch verification
- `R191-target-repo-config.md` - Target repository configuration
- `R192-repo-separation.md` - Repository separation requirements
- `R193-effort-clone-protocol.md` - Clone protocol requirements
- `R194-remote-branch-tracking.md` - Branch tracking requirements
- `R195-branch-push-verification.md` - Push verification requirements
- `R196-base-branch-selection.md` - Base branch selection

#### Directory & Isolation Protocols
- `R209-effort-directory-isolation-protocol.md` - Effort directory isolation
- `R212-phase-directory-isolation-protocol.md` - Phase directory isolation
- `R213-wave-and-effort-metadata-protocol.md` - Wave and effort metadata
- `R214-code-reviewer-wave-directory-acknowledgment.md` - Directory acknowledgment

#### Integration Rules (R260-R280)
- `R260-integration-agent-core-requirements.md` - Core integration requirements
- `R261-integration-planning-requirements.md` - Integration planning
- `R262-merge-operation-protocols.md` - Merge operation protocols
- `R263-integration-documentation-requirements.md` - Documentation requirements
- `R264-work-log-tracking-requirements.md` - Work log tracking
- `R265-integration-testing-requirements.md` - Testing requirements
- `R266-upstream-bug-documentation.md` - Bug documentation
- `R267-integration-agent-grading-criteria.md` - Grading criteria
- `R268-orchestrator-integration-agent-spawn.md` - Integration agent spawning
- `R269-code-reviewer-merge-plan-no-execution.md` - Merge plan restrictions
- `R270-no-integration-branches-as-sources.md` - Branch source restrictions
- `R271-mandatory-production-ready-validation.md` - Production readiness validation
- `R272-integration-testing-branch.md` - Testing branch requirements
- `R273-runtime-specific-validation.md` - Runtime validation
- `R274-production-readiness-checklist.md` - Production readiness checklist
- `R275-deployment-verification.md` - Deployment verification
- `R276-runbook-requirement.md` - Runbook requirements
- `R277-continuous-build-verification.md` - Build verification
- `R278-external-user-validation.md` - External user validation
- `R279-master-pr-plan-requirement.md` - Master PR plan requirements
- `R280-main-branch-protection.md` - Main branch protection

#### Performance & Parallelization
- `R151-parallel-agent-spawning-timing.md` - Parallel spawning timing
- `R152-implementation-speed.md` - Implementation speed requirements
- `R153-review-turnaround.md` - Review turnaround requirements
- `R158-pattern-compliance-rate.md` - Pattern compliance rates
- `R218-orchestrator-parallel-code-reviewer-spawning.md` - Parallel reviewer spawning
- `R219-code-reviewer-dependency-aware-effort-planning.md` - Dependency-aware planning

#### Agent Specialization Rules
- `R210-architect-architecture-planning-protocol.md` - Architecture planning protocol
- `R211-code-reviewer-implementation-from-architecture.md` - Implementation from architecture
- `R215-orchestrator-state-ownership.md` - State ownership
- `R216-bash-execution-syntax.md` - Bash execution requirements
- `R022-architect-size-verification.md` - Architect size verification

#### Error Handling & Recovery
- `R171-precompact-hook.md` - Precompact hook requirements
- `R172-utility-scripts.md` - Utility script requirements
- `R173-state-preservation.md` - State preservation
- `R174-recovery-detection.md` - Recovery detection
- `R175-manual-utilities.md` - Manual utility requirements
- `R186-automatic-compaction-detection.md` - Compaction detection
- `R254-AGENT-ERROR-REPORTING.md` - Error reporting requirements
- `R255-POST-AGENT-WORK-VERIFICATION.md` - Post-work verification

#### Gate & Assessment Rules
- `R256-mandatory-phase-assessment-gate.md` - Phase assessment gate
- `R257-mandatory-phase-assessment-report.md` - Phase assessment report
- `R258-mandatory-wave-review-report.md` - Wave review report
- `R259-mandatory-phase-integration-after-fixes.md` - Phase integration after fixes

#### Enforcement & Compliance
- `R021-orchestrator-never-stops.md` - Continuous operation requirement
- `R014-branch-naming-convention.md` - Branch naming conventions
- `R197-one-agent-per-effort.md` - Agent assignment restrictions
- `R230-state-machine-visualization-requirement.md` - State machine visualization
- `R231-continuous-operation-through-transitions.md` - Continuous operation
- `R232-enforcement-examples.md` - Enforcement examples
- `R232-todowrite-pending-items-override.md` - TodoWrite overrides
- `R233-all-states-immediate-action.md` - Immediate action requirements

#### Additional Reference Files
- `DELIMITER-AND-CRITICALITY-SYSTEM.md` - Rule formatting and criticality system
- `HOOK-RULES-SUMMARY.md` - Hook system rule summary
- `RULE-CRITICALITY-FORMATTING-GUIDE.md` - Criticality formatting guide
- `RULE-REGISTRY.md` - Complete rule registry
- `STATE-TRANSITION-ANTI-PATTERNS.md` - Anti-patterns to avoid

---

## 🛠️ UTILITY LOCATIONS

### Tools & Scripts
- `/home/vscode/software-factory-template/tools/line-counter.sh` - Official line counter (MANDATORY for size measurement)
- `/home/vscode/software-factory-template/utilities/` - Recovery and validation utilities

### State Machine Definitions
- `/home/vscode/software-factory-template/SOFTWARE-FACTORY-STATE-MACHINE.md` - Master state machine
- `/home/vscode/software-factory-template/state-machines/` - Agent-specific state machines

### Templates
- `/home/vscode/software-factory-template/templates/` - Implementation and planning templates

---

## 🎯 QUICK ACCESS PATTERNS

### For Rule Lookup:
1. **Rule Number Known:** `/home/vscode/software-factory-template/rule-library/R###-*.md`
2. **Agent-Specific:** `/home/vscode/software-factory-template/.claude/agents/[agent].md`
3. **State-Specific:** `/home/vscode/software-factory-template/agent-states/[agent]/[STATE]/rules.md`

### For Emergency Recovery:
1. Check compaction: `/home/vscode/software-factory-template/utilities/check-compaction.sh`
2. TODO recovery: `/home/vscode/software-factory-template/utilities/todo-preservation.sh`
3. State verification: `/home/vscode/software-factory-template/utilities/validate-factory-compliance.sh`

### For Size Management:
1. Line counting: `/home/vscode/software-factory-template/tools/line-counter.sh`
2. Split validation: `/home/vscode/software-factory-template/utilities/validate-split-boundaries.sh`

---

*This reference document maps the complete rule system for Software Factory 2.0. All paths are absolute and verified current as of the template structure.*