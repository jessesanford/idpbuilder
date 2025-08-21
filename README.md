# Software Factory Template

A comprehensive template for orchestrating large-scale software development using multiple specialized AI agents coordinated through a state machine workflow.

## 🎯 Overview

This template provides a complete system for:
- **Orchestrated Development**: Coordinate multiple AI agents through phases and waves
- **Quality Gates**: Enforce size limits, reviews, and architectural compliance  
- **State Management**: Track progress through a formal state machine
- **Continuous Integration**: Manage splits, fixes, and integration branches
- **Context Recovery**: Handle agent context loss and compaction events

## 🏗️ System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    OPERATOR (Human)                          │
│  - Sets up project structure                                 │
│  - Configures agents and limits                              │
│  - Invokes orchestrator with /continue-orchestrating         │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│              ORCHESTRATOR (orchestrator-task-master)         │
│  - Reads state machine and plans                             │
│  - Manages state file (orchestrator-state.yaml)              │
│  - Spawns and coordinates agents                             │
│  - NEVER writes code (coordination only)                     │
└─────────────────────────────────────────────────────────────┘
                            │
        ┌───────────────────┼───────────────────┐
        ▼                   ▼                   ▼
┌──────────────┐   ┌──────────────┐   ┌──────────────┐
│  ARCHITECT   │   │CODE REVIEWER │   │ SW ENGINEER  │
│   Reviews    │   │    Plans &   │   │  Implements  │
│ Architecture │   │   Reviews    │   │     Code     │
└──────────────┘   └──────────────┘   └──────────────┘
```

## 📁 Complete Directory Structure & File Reference

```
software-factory-template/
│
├── README.md                           # This file (main documentation)
├── setup.sh                           # Interactive setup script
├── orchestrator-state-example.yaml    # Example state file
├── PROJECT-IMPLEMENTATION-PLAN-TEMPLATE.md  # Template for high-level plan
├── HOW-TO-PLAN.md                     # Planning methodology guide
├── PLANNING-AGENT-ASSIGNMENTS.md      # Which agent does what in planning
├── TEMPLATE-CREATION-SUMMARY.md       # How this template was created
├── CRITICAL-FILES-ADDED.md           # Documentation of critical files
├── CRITICAL-SETTINGS-JSON.md         # 🔴 Why settings.json is essential
├── FINAL-FILE-ORGANIZATION.md        # File organization details
├── CLAUDE-MD-FILE-VERIFICATION.md    # Verification of all references
│
├── .claude/                           # Claude Code configuration
│   ├── settings.json                  # 🔴 CRITICAL: Compaction hooks for TODO preservation
│   ├── CLAUDE.md                      # 🔴 CRITICAL: Global rules, compaction recovery, TODO management
│   ├── agents/                        # Agent configurations
│   │   ├── orchestrator-task-master.md  # Orchestrator agent config
│   │   ├── architect-reviewer.md     # Architect agent config
│   │   ├── code-reviewer.md          # Code reviewer agent config
│   │   └── sw-engineer-example-go.md # Example SW engineer (customize for your stack)
│   └── commands/                      # Reusable command workflows
│       ├── continue-orchestrating.md # Main command to start/resume orchestration
│       ├── code-review.md            # Standard code review workflow
│       └── create-wave-impl-plan.md  # Wave planning workflow
│
├── core/                              # Core system files
│   ├── SOFTWARE-FACTORY-STATE-MACHINE.md  # ⚡ THE HEART - Defines all states and transitions
│   └── ORCHESTRATOR-MASTER-OPERATIONS-GUIDE.md  # Complete operational blueprint
│
├── protocols/                         # Critical execution protocols (18 files)
│   ├── IMPERATIVE-LINE-COUNT-RULE.md # 🚨 Size limit enforcement (referenced 6x in CLAUDE.md)
│   ├── ORCHESTRATOR-NEVER-WRITES-CODE-RULE.md # 🔴 Fundamental separation of concerns
│   ├── EFFORT-SPLIT-CONTINUOUS-EXECUTION-PROTOCOL.md  # How to handle >800 line efforts
│   ├── SW-ENGINEER-STARTUP-REQUIREMENTS.md  # SW Engineer initialization
│   ├── SW-ENGINEER-EXPLICIT-INSTRUCTIONS.md  # Detailed SW Engineer instructions
│   ├── ORCHESTRATOR-EFFORT-PLANNING-PROTOCOL.md  # How orchestrator manages planning
│   ├── ORCHESTRATOR-TASKMASTER-EXECUTION-PLAN.md  # Complete execution guide
│   ├── ORCHESTRATOR-CODE-REVIEW-INTEGRATION.md  # How orchestrator handles reviews
│   ├── CODE-REVIEWER-EFFORT-PLANNING-INSTRUCTIONS.md  # How to create effort plans
│   ├── CODE-REVIEWER-COMPREHENSIVE-GUIDE.md  # Complete review process
│   ├── CODE-REVIEW-ENFORCEMENT-SUMMARY.md  # Review enforcement rules
│   ├── WAVE-COMPLETION-ARCHITECT-REVIEW-PROTOCOL.md  # Wave review requirements
│   ├── ARCHITECT-REVIEWER-WAVE-INSTRUCTIONS.md  # Detailed architect instructions
│   ├── PHASE-START-ARCHITECT-REVIEW-PROTOCOL.md  # Phase assessment protocol
│   ├── PHASE-COMPLETION-FUNCTIONAL-TESTING.md  # Phase functional testing requirements
│   ├── TEST-DRIVEN-VALIDATION-REQUIREMENTS.md  # Testing coverage requirements
│   ├── WORK-LOG-TEMPLATE.md          # Template for effort work logs
│   └── TODO-STATE-MANAGEMENT-PROTOCOL.md  # TODO persistence (referenced in CLAUDE.md)
│
├── efforts/                           # Working directories for implementations
│   └── README.md                      # Explains effort isolation and structure
│
├── todos/                             # TODO state persistence
│   └── README.md                      # How TODO files work for context recovery
│
├── tools/                             # Utilities
│   └── line-counter.sh                # Configurable size measurement tool
│
├── phase-plans/                       # Templates for phase-specific planning
│   ├── README.md                      # How to create phase plans
│   ├── PHASEX-GENERIC-TEMPLATE.md    # Generic template for any phase
│   ├── PHASE1-TEMPLATE.md            # Example: API/Contract phase
│   ├── PHASE2-TEMPLATE.md            # Example: Infrastructure phase
│   └── PHASE3-TEMPLATE.md            # Example: Implementation phase
│
├── agent-instructions/                # Templates for spawning agents with tasks
│   ├── README.md                      # How to use instruction templates
│   ├── sw-engineer-implementation.md # Template for implementation tasks
│   ├── code-reviewer-planning.md     # Template for creating effort plans
│   ├── code-reviewer-review.md       # Template for reviewing code
│   └── architect-wave-review.md      # Template for wave architecture review
│
├── quick-reference/                   # Quick guides for agents
│   ├── ORCHESTRATOR-QUICK-REFERENCE.md  # Quick state transitions
│   ├── ORCHESTRATOR-WORKFLOW-SUMMARY.md # Visual workflow summary
│   └── CODE-REVIEWER-QUICK-REFERENCE.md # Quick decision trees
│
├── examples/                          # Real-world examples
│   ├── SPLIT-EXAMPLE-AUTHENTICATION-MODULE.md # Generic split example
│   ├── SPLIT-REVIEW-LOOP-DIAGRAM.md  # Visual split process
│   └── CODE-REVIEW-EXAMPLES.md       # Real review examples
│
└── possibly-needed-but-not-sure/      # Optional/advanced protocols
    ├── README.md                      # Guide to all optional files
    ├── WHEN-TO-USE-THESE-FILES.md    # Decision guide for activation
    ├── FILES-NOT-INCLUDED.md         # What wasn't included and why
    └── SPLIT-EXAMPLE-E3.1.1-SYNC-ENGINE.md  # TMC-specific split example
```

## 📚 Critical Files - Who Uses Them and When

### 🔴 Always Active Files (Required for System to Function)

#### Core System Files
| File | Used By | When | Purpose |
|------|---------|------|---------|
| **SOFTWARE-FACTORY-STATE-MACHINE.md** | Orchestrator | FIRST at startup, continuously | Defines entire workflow, states, transitions |
| **ORCHESTRATOR-MASTER-OPERATIONS-GUIDE.md** | Orchestrator, Human | Setup and operations | Complete blueprint for running system |
| **CLAUDE.md** | ALL agents | Every startup | Global rules, recovery procedures, TODO management |

#### Protocol Files
| File | Used By | When | Purpose |
|------|---------|------|---------|
| **IMPERATIVE-LINE-COUNT-RULE.md** | ALL agents | Every effort | Absolute size limit enforcement |
| **ORCHESTRATOR-NEVER-WRITES-CODE-RULE.md** | Orchestrator | Always | Fundamental separation of concerns |
| **EFFORT-SPLIT-CONTINUOUS-EXECUTION-PROTOCOL.md** | Orchestrator, Code Reviewer | When effort >800 lines | Sequential split execution |
| **SW-ENGINEER-STARTUP-REQUIREMENTS.md** | SW Engineer | Every task startup | Environment verification, startup protocol |
| **SW-ENGINEER-EXPLICIT-INSTRUCTIONS.md** | SW Engineer | Every startup | Git commands, validation, build procedures |
| **ORCHESTRATOR-EFFORT-PLANNING-PROTOCOL.md** | Orchestrator | Before each effort | How to coordinate planning |
| **ORCHESTRATOR-TASKMASTER-EXECUTION-PLAN.md** | Orchestrator | Always | Complete execution guide |
| **ORCHESTRATOR-CODE-REVIEW-INTEGRATION.md** | Orchestrator | During reviews | How to handle review outcomes |
| **CODE-REVIEWER-EFFORT-PLANNING-INSTRUCTIONS.md** | Code Reviewer | When creating plans | How to create implementation plans |
| **CODE-REVIEWER-COMPREHENSIVE-GUIDE.md** | Code Reviewer | Every review | Complete review process and standards |
| **CODE-REVIEW-ENFORCEMENT-SUMMARY.md** | Code Reviewer | Every review | Review enforcement rules |
| **WAVE-COMPLETION-ARCHITECT-REVIEW-PROTOCOL.md** | Orchestrator, Architect | End of each wave | Mandatory wave review process |
| **ARCHITECT-REVIEWER-WAVE-INSTRUCTIONS.md** | Architect | Wave reviews | Detailed wave review instructions |
| **PHASE-START-ARCHITECT-REVIEW-PROTOCOL.md** | Architect | Phase boundaries | Phase assessment protocol |
| **PHASE-COMPLETION-FUNCTIONAL-TESTING.md** | Orchestrator, Code Reviewer | End of each phase | Functional testing before phase transition |
| **TEST-DRIVEN-VALIDATION-REQUIREMENTS.md** | SW Engineer, Code Reviewer | Every implementation/review | Testing coverage requirements |
| **WORK-LOG-TEMPLATE.md** | SW Engineer | Every effort | Progress tracking template |
| **TODO-STATE-MANAGEMENT-PROTOCOL.md** | ALL agents | State transitions | TODO file management procedures |

#### Agent Configurations
| File | Used By | When | Purpose |
|------|---------|------|---------|
| **orchestrator-task-master.md** | Claude Code | When spawning orchestrator | Defines orchestrator behavior |
| **code-reviewer.md** | Claude Code | When spawning reviewer | Defines reviewer behavior |
| **architect-reviewer.md** | Claude Code | When spawning architect | Defines architect behavior |
| **sw-engineer-example-go.md** | Claude Code | When spawning engineer | Example config (customize) |

#### Command Files
| File | Used By | When | Purpose |
|------|---------|------|---------|
| **continue-orchestrating.md** | Human | To start/resume | Main orchestration command |

### 🟡 Optional Files (In possibly-needed-but-not-sure/)

#### Quick References (Helpful for Learning)
| File | Best For | When to Activate | Purpose |
|------|----------|-----------------|---------|
| **CODE-REVIEWER-QUICK-REFERENCE.md** | New reviewers | Learning phase | Quick decision trees and checks |
| **ORCHESTRATOR-QUICK-REFERENCE.md** | New orchestrators | Learning phase | State transitions and commands |
| **ORCHESTRATOR-WORKFLOW-SUMMARY.md** | Everyone | Understanding flow | Visual workflow with examples |

#### Advanced Protocols (For Complex Projects)
| File | Best For | When to Activate | Purpose |
|------|----------|-----------------|---------|
| **PHASE-COMPLETION-FUNCTIONAL-TESTING.md** | Quality-critical projects | End of phases | Comprehensive testing |

#### Examples and References (For Understanding)
| File | Best For | When to Use | Purpose |
|------|----------|-------------|---------|
| **CODE-REVIEW-EXAMPLES.md** | Learning | Understanding reviews | Real review examples |
| **SPLIT-EXAMPLE-E3.1.1-SYNC-ENGINE.md** | Learning splits | First split needed | Real 2400→3 parts example |
| **SPLIT-REVIEW-LOOP-DIAGRAM.md** | Visual learners | Understanding splits | Diagram of split process |
| **CODE-REVIEW-ENFORCEMENT-SUMMARY.md** | Compliance | Setting up gates | All enforcement points |

## 🔄 State Machine Flow

The system follows a strict state machine defined in `SOFTWARE-FACTORY-STATE-MACHINE.md`:

### Major Loops
1. **Phase Loop**: Highest level, includes assessment gates
2. **Wave Loop**: Groups of related efforts with architecture review
3. **Effort Loop**: Individual implementation units
4. **Split Loop**: Handles oversized efforts sequentially
5. **Fix Loop**: Addresses review feedback

### Key States
- `INIT` → `PHASE_START_GATE` → `WAVE_START` → `EFFORT_SELECTION`
- `IMPLEMENTATION` → `MEASURE_SIZE` → `CODE_REVIEW`
- `WAVE_COMPLETE` → `ARCHITECT_REVIEW` → `NEXT_WAVE`

## 👥 Agent Roles and Their Files

### Orchestrator (orchestrator-task-master)
**Reads on Startup:**
- SOFTWARE-FACTORY-STATE-MACHINE.md
- ORCHESTRATOR-TASKMASTER-EXECUTION-PLAN.md
- orchestrator-state.yaml
- ORCHESTRATOR-EFFORT-PLANNING-PROTOCOL.md

**Responsibilities:**
- Coordinates all activities
- Manages state transitions
- Spawns other agents
- **NEVER writes code**

### Code Reviewer (code-reviewer)
**Reads on Startup:**
- CODE-REVIEWER-EFFORT-PLANNING-INSTRUCTIONS.md
- IMPERATIVE-LINE-COUNT-RULE.md
- TEST-DRIVEN-VALIDATION-REQUIREMENTS.md
- EFFORT-SPLIT-CONTINUOUS-EXECUTION-PROTOCOL.md

**Responsibilities:**
- Creates implementation plans
- Reviews code for quality
- Designs split strategies
- Ensures compliance

### SW Engineer (sw-engineer)
**Reads on Startup:**
- SW-ENGINEER-STARTUP-REQUIREMENTS.md
- IMPERATIVE-LINE-COUNT-RULE.md
- TEST-DRIVEN-VALIDATION-REQUIREMENTS.md
- IMPLEMENTATION-PLAN.md (in working directory)

**Responsibilities:**
- Implements code per plan
- Measures size continuously
- Fixes review issues
- Works in isolated environments

### Architect (architect-reviewer)
**Reads on Startup:**
- WAVE-COMPLETION-ARCHITECT-REVIEW-PROTOCOL.md
- orchestrator-state.yaml
- PHASE-START-ARCHITECT-REVIEW-PROTOCOL.md (if in possibly-needed)

**Responsibilities:**
- Reviews architectural consistency
- Assesses phase progress
- Approves wave completions
- Identifies technical debt

## 📏 Size Management

### Continuous Measurement
```bash
# After every logical change
./tools/line-counter.sh -c branch-name

# Detailed breakdown
./tools/line-counter.sh -c branch-name -d
```

### Thresholds
- **Warning**: 700 lines - plan for completion
- **Error**: 800 lines - must split immediately

### Split Protocol
1. Detection → Stop implementation
2. Code Reviewer creates split plan
3. Sequential execution of splits
4. Each split reviewed independently
5. Continue until all under limit

## 🔒 Quality Gates

### Mandatory Gates
- ✅ Size compliance (every effort)
- ✅ Code review (every effort)
- ✅ Test coverage (per requirements)
- ✅ Architecture review (every wave)
- ✅ Phase assessment (phase boundaries)

### Gate Failures
- Size violation → Split protocol
- Review failure → Fix cycle
- Architecture issues → Addendum
- Phase off-track → STOP

## 💾 State Persistence

### TODO Management
Agents save TODO files during state transitions:
```
todos/orchestrator-WAVE_COMPLETE-20250121-143000.todo
todos/sw-eng-IMPLEMENTATION-20250121-145500.todo
```

### Recovery After Context Loss
1. Check compaction marker
2. Read latest TODO file
3. Load into TodoWrite tool
4. Resume from saved state

## 📋 Planning Your Project

### Planning Process Overview

Before starting orchestration, follow the structured planning process:

1. **Study Planning Guide**: Read `HOW-TO-PLAN.md` for the complete methodology
2. **Create High-Level Plan**: Structure your project into phases, waves, and efforts
3. **Detail Each Phase**: Use templates in `phase-plans/` for specific implementation plans
4. **Validate Plans**: Ensure every effort has explicit instructions, tests, and success criteria

### Key Planning Documents

- **HOW-TO-PLAN.md**: Step-by-step planning methodology with agent targeting
- **phase-plans/**: Templates and examples for detailed phase planning
- **PROJECT-IMPLEMENTATION-PLAN-TEMPLATE.md**: High-level orchestration template

## 🚀 Quick Start

### 1. Setup Your Project
```bash
# Copy this template to your project
cp -r /workspaces/software-factory-template /workspaces/your-project

# Navigate to your project
cd /workspaces/your-project

# Customize configurations
# - Edit .claude/agents/* for your tech stack
# - Update tools/line-counter.sh patterns for your language
# - Create your PROJECT-IMPLEMENTATION-PLAN.md
# - Define your phases and waves
```

### 2. Configure Size Limits
Edit `tools/line-counter.sh`:
```bash
# Configuration - CUSTOMIZE FOR YOUR PROJECT
MAX_LINES_WARNING=700
MAX_LINES_ERROR=800
```

### 3. Create Your Implementation Plans

#### Step 3a: High-Level Plan
Create `orchestrator/PROJECT-IMPLEMENTATION-PLAN.md`:
```markdown
# Project Implementation Plan

## Phase 1: Foundation
### Wave 1: Core Types
- E1.1.1: Basic data models
- E1.1.2: API interfaces
...
```

#### Step 3b: Detailed Phase Plans
For each phase, create a detailed plan using templates:
```bash
cp phase-plans/PHASEX-GENERIC-TEMPLATE.md phase-plans/PHASE1-SPECIFIC-IMPL-PLAN.md
# Edit with:
# - Exact source branches to reuse
# - Actual TDD test cases
# - Detailed pseudo-code
# - Specific validation commands
```

### 4. Initialize State
Create `orchestrator-state.yaml`:
```yaml
current_phase: 1
current_wave: 1
current_state: "INIT"
efforts_completed: []
efforts_in_progress: []
efforts_pending: [...]
```

### 5. Start Orchestration
```
/continue-orchestrating
```

## 🛠️ Customization Guide

### For Different Languages

#### Python Projects
Edit `.claude/agents/sw-engineer-python.md`:
```yaml
Focus Areas:
- Type hints and dataclasses
- Async/await patterns
- pytest and coverage
- Virtual environments
```

#### JavaScript/TypeScript
Edit `.claude/agents/sw-engineer-typescript.md`:
```yaml
Focus Areas:
- TypeScript strict mode
- React/Vue/Angular patterns
- Jest/Mocha testing
- npm/yarn workflows
```

### For Different Domains

#### Web Applications
- Add frontend/backend split protocols
- Include API documentation requirements
- Add E2E testing gates

#### Microservices
- Add service boundary protocols
- Include API contract testing
- Add deployment manifests

#### Data Pipelines
- Add data validation gates
- Include performance benchmarks
- Add data lineage tracking

## 🎯 Which Files to Activate When

### Minimum Viable Setup (Simple Project)
Use only the files in:
- `/core/`
- `/protocols/`
- `/.claude/`

### Medium Complexity (Multiple Phases)
Additionally activate from possibly-needed:
- ARCHITECT-REVIEWER-WAVE-INSTRUCTIONS.md
- PHASE-START-ARCHITECT-REVIEW-PROTOCOL.md
- Quick reference guides

### High Complexity (Enterprise)
Activate everything from possibly-needed:
- All architect protocols
- TODO-STATE-MANAGEMENT-PROTOCOL.md
- Comprehensive review guides
- Testing protocols

### Learning the System
Start with:
- Quick reference guides
- Workflow summary
- Split examples
- Review examples

## ⚠️ Critical Configuration

### 🔴 ESSENTIAL: settings.json

The `.claude/settings.json` file is **CRITICAL** for the system to function. It enables:
- **Compaction Recovery**: Preserves state when memory limits are reached
- **TODO Persistence**: Saves TODO lists before context compression
- **State Management**: Maintains working context across sessions

**WITHOUT settings.json, YOU WILL LOSE WORK!** See `CRITICAL-SETTINGS-JSON.md` for details.

## ⚠️ Critical Rules

### Never Do This ❌
1. Orchestrator writes code
2. Parallel splits
3. Skip reviews
4. Ignore size limits
5. Mix working directories

### Always Do This ✅
1. Measure continuously
2. Review everything
3. Follow state machine
4. Commit state file
5. Document progress

## 🔧 Troubleshooting

### Common Issues

| Problem | Solution | Relevant Files |
|---------|----------|----------------|
| "State file not found" | Create from implementation plan | orchestrator-state-example.yaml |
| "Effort over limit" | Implement split protocol | EFFORT-SPLIT-CONTINUOUS-EXECUTION-PROTOCOL.md |
| "Context lost" | Use TODO recovery process | CLAUDE.md sections 7-9 |
| "Wrong directory" | Never fix with cd, report error | SW-ENGINEER-STARTUP-REQUIREMENTS.md |
| "Review failed" | Fix and re-review, don't skip | CODE-REVIEWER-EFFORT-PLANNING-INSTRUCTIONS.md |

### Recovery Procedures
See `.claude/CLAUDE.md` sections:
- Section 7: Context Loss Recovery
- Section 8: TODO State Management  
- Section 9: Pre-Compaction Saving

## 📈 Success Metrics

Track these KPIs:
- **Size Compliance**: 100% of efforts under limit
- **Review Pass Rate**: >80% first-pass reviews
- **Architecture Drift**: 0 STOP decisions
- **Feature Coverage**: 100% requirements met

## 🤝 Contributing

This template is designed to be extended and customized. Key extension points:
- Agent configurations in `.claude/agents/`
- Language patterns in `tools/line-counter.sh`
- Protocols in `protocols/`
- State machine states in core files

## 📄 License

This template is provided as a reference implementation for AI-orchestrated software development workflows.

---

**Remember**: The state machine is law. Trust the process. Let the system work.