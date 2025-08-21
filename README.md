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

## 📁 Directory Structure

```
software-factory-template/
│
├── README.md                           # This file
│
├── .claude/                            # Claude Code configuration
│   ├── CLAUDE.md                       # Global rules for all agents
│   ├── agents/                         # Agent configurations
│   │   ├── orchestrator-task-master.md
│   │   ├── architect-reviewer.md
│   │   ├── code-reviewer.md
│   │   └── sw-engineer-example-go.md  # Example - customize for your stack
│   └── commands/
│       └── continue-orchestrating.md  # Main orchestration command
│
├── core/                               # Core system files
│   ├── SOFTWARE-FACTORY-STATE-MACHINE.md  # ⚡ THE HEART OF THE SYSTEM
│   └── ORCHESTRATOR-MASTER-OPERATIONS-GUIDE.md
│
├── protocols/                          # Execution protocols
│   ├── EFFORT-SPLIT-CONTINUOUS-EXECUTION-PROTOCOL.md
│   └── SW-ENGINEER-STARTUP-REQUIREMENTS.md
│
├── agent-instructions/                 # Agent-specific instructions
│   └── [Additional instruction files as needed]
│
├── efforts/                            # Working directories for implementations
│   └── README.md                       # Explains effort isolation
│
├── todos/                              # TODO state persistence
│   └── README.md                       # Explains TODO management
│
├── tools/                              # Utilities
│   └── line-counter.sh                 # Size measurement tool
│
└── possibly-needed-but-not-sure/       # Additional protocols
    └── [Various protocol files]
```

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

### 3. Create Your Implementation Plan

Create `orchestrator/PROJECT-IMPLEMENTATION-PLAN.md`:
```markdown
# Project Implementation Plan

## Phase 1: Foundation
### Wave 1: Core Types
- E1.1.1: Basic data models
- E1.1.2: API interfaces
...
```

### 4. Initialize State

Create `orchestrator/orchestrator-state.yaml`:
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

## 👥 Agent Roles

### Orchestrator (orchestrator-task-master)
- Coordinates all activities
- Manages state transitions
- Spawns other agents
- **NEVER writes code**

### Code Reviewer (code-reviewer)
- Creates implementation plans
- Reviews code for quality
- Designs split strategies
- Ensures compliance

### SW Engineer (sw-engineer)
- Implements code per plan
- Measures size continuously
- Fixes review issues
- Works in isolated environments

### Architect (architect-reviewer)
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

## 📚 Key Documents

### Must Read First
1. **SOFTWARE-FACTORY-STATE-MACHINE.md** - The core workflow
2. **ORCHESTRATOR-MASTER-OPERATIONS-GUIDE.md** - How to operate
3. **.claude/CLAUDE.md** - Agent rules and recovery

### Protocol Documents
- **EFFORT-SPLIT-CONTINUOUS-EXECUTION-PROTOCOL.md** - Handling large efforts
- **SW-ENGINEER-STARTUP-REQUIREMENTS.md** - Agent initialization

### Command Reference
- **/continue-orchestrating** - Start/resume orchestration

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

| Problem | Solution |
|---------|----------|
| "State file not found" | Create from implementation plan |
| "Effort over limit" | Implement split protocol |
| "Context lost" | Use TODO recovery process |
| "Wrong directory" | Never fix with cd, report error |
| "Review failed" | Fix and re-review, don't skip |

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