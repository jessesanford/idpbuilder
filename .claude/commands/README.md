# Software Factory 2.0 Slash Commands

This directory contains comprehensive slash commands for managing the Software Factory 2.0 development process. These commands are designed to be project-agnostic and follow all critical Software Factory protocols.

## Available Commands

### 🎯 Agent Continuation Commands

#### `/continue-orchestrating`
**Purpose**: Primary orchestrator agent continuation with full state machine support
**Use When**: 
- Starting fresh orchestration
- Resuming orchestration after interruption  
- Managing wave completions and integrations
- Spawning and coordinating other agents

**Key Features**:
- Mandatory pre-flight checks and agent acknowledgment
- State-driven execution (INIT, WAVE_START, WAVE_COMPLETE, etc.)
- TODO recovery and state persistence
- Line count compliance enforcement
- Integration gate management
- Never implements code directly - only coordinates

#### `/continue-implementing`
**Purpose**: Software Engineer agent continuation with implementation focus
**Use When**:
- Implementing features per implementation plans
- Working on effort branches
- Fixing issues after code review
- Working on effort splits

**Key Features**:
- Test-driven development protocols
- Incremental size monitoring (every 200 lines)
- Work log maintenance
- Branch management and Git hygiene
- Pre-review checklists
- Cannot exceed 800-line limit per effort

#### `/continue-reviewing` 
**Purpose**: Code Reviewer agent continuation with quality assurance
**Use When**:
- Creating implementation plans for efforts
- Reviewing completed implementations
- Planning effort splits when >800 lines
- Validating test coverage and compliance

**Key Features**:
- Line count assessment and enforcement
- Split planning when efforts exceed limits
- Test coverage validation
- Code quality and pattern compliance
- Documentation requirements checking
- Cannot approve non-compliant efforts

#### `/continue-architecting`
**Purpose**: Architect agent continuation with technical leadership
**Use When**:
- Reviewing completed waves
- Assessing phase readiness
- Evaluating integration branches
- Making PROCEED/CHANGES_REQUIRED/STOP decisions

**Key Features**:
- Comprehensive technical excellence review
- Integration conflict detection
- Scalability and security assessment
- Performance analysis
- Architectural pattern compliance
- Cannot implement code - only reviews and decides

### 🔧 System Management Commands

#### `/reset-state`
**Purpose**: Controlled state reset for corrupted or unrecoverable situations
**Use When**:
- State machine is corrupted beyond recovery
- TODO system is completely broken
- Multiple agents have conflicting state
- Fresh start is needed

**Safety Levels**:
- **Level 1**: TODO-only reset (safest)
- **Level 2**: State machine reset (moderate) 
- **Level 3**: Full reset (dangerous)

**Key Features**:
- Mandatory backup creation before reset
- Graduated reset levels with increasing scope
- Safety checks and abort options
- Recovery guidance after reset
- Preservation of work in backup branches

#### `/check-status`
**Purpose**: Comprehensive diagnostic analysis of current system state
**Use When**:
- Diagnosing system issues
- Understanding current progress
- Planning recovery actions
- Verifying system health

**Key Features**:
- Environment and configuration validation
- State machine consistency analysis
- TODO system status assessment
- Git branch and line count analysis
- Recovery recommendations
- System health scoring

## Command Usage Guidelines

### 🚨 Critical Requirements

All commands enforce these mandatory requirements:

1. **Agent Identity Verification**: Commands verify agent identity matches expected role
2. **Environment Validation**: Working directory, Git branch, and remote tracking verification
3. **Line Count Compliance**: Every effort must be ≤800 lines using project line counter
4. **TODO State Management**: Proper save/load protocols during state transitions
5. **Pre-flight Checks**: Comprehensive startup verification for all agents

### 📋 State Machine Integration

Commands are tightly integrated with the Software Factory 2.0 state machine:

- **State Detection**: Automatically detect current state from orchestrator-state.yaml
- **State Transitions**: Proper TODO saving before transitions
- **State Validation**: Ensure transitions follow valid state machine paths
- **Recovery Support**: Context recovery protocols for lost or corrupted state

### 🔄 Context Recovery

All commands support context recovery after interruptions:

1. **Check for Compaction**: Automatic detection of Claude Code compaction
2. **TODO Recovery**: Load and merge saved TODO state files
3. **State Reconstruction**: Rebuild working context from persistent files
4. **Validation**: Verify recovered state is consistent and valid

### 📊 Project Agnostic Design

Commands use `[project]` placeholders for project-specific elements:

- `./agent-configs/[project]/` - Project configuration directory
- `./tools/[project]-line-counter.sh` - Project-specific line counting tool
- `[LANG]-sw-eng.md` - Language-specific agent configurations
- `PHASE{X}-SPECIFIC-IMPL-PLAN.md` - Phase-specific planning documents

## File Structure Requirements

For commands to work properly, projects should maintain this structure:

```
/workspaces/[project]/
├── agent-configs/[project]/          # Project-specific agent configs
│   ├── orchestrator-state.yaml      # State machine persistence
│   ├── todos/                       # TODO state files
│   ├── [PROJECT]-ORCHESTRATOR-IMPLEMENTATION-PLAN.md
│   ├── SOFTWARE-ENG-AGENT-STARTUP-REQUIREMENTS.md
│   ├── CODE-REVIEWER-COMPREHENSIVE-GUIDE.md
│   └── architect.md
├── agent-states/                    # Agent state definitions
│   ├── orchestrator/               # Orchestrator states
│   ├── sw-engineer/                # SW Engineer states
│   ├── code-reviewer/              # Code Reviewer states
│   └── architect/                  # Architect states
├── state-machines/                 # State machine definitions
├── tools/                         # Project-specific tools
│   └── [project]-line-counter.sh  # CRITICAL: Line counting tool
└── efforts/                       # Implementation work directories
```

## Integration with Software Factory 2.0

These commands are designed to work seamlessly with:

- **State Machine Framework**: Multi-agent state management
- **Grading System**: Automated quality assurance
- **Rule Library**: Comprehensive rule enforcement 
- **Expertise Modules**: Domain-specific knowledge
- **Line Count System**: Size compliance enforcement
- **TODO Management**: Task persistence across interruptions

## Best Practices

### Command Selection
- Use agent-specific continuation commands for normal operations
- Use `/check-status` first when diagnosing issues
- Use `/reset-state` only when recovery is impossible
- Always verify agent identity before using commands

### State Management
- Save TODO state before major transitions
- Commit and push TODO files to preserve state
- Use proper state machine transitions
- Verify environment before proceeding

### Recovery Protocols
- Check for compaction markers on startup
- Load TODO state using TodoWrite tool (not just reading)
- Deduplicate recovered TODOs with existing ones
- Validate recovered state before proceeding

### Quality Assurance
- Never bypass line count measurements
- Always use project-specific line counter tool
- Enforce test coverage requirements
- Maintain proper documentation and work logs

## Troubleshooting

### Common Issues

1. **Agent Identity Mismatch**: Verify your prompt contains correct `@agent-*` designation
2. **Wrong Directory**: Commands check working directory - don't try to cd to fix
3. **Missing State Files**: Use `/check-status` to diagnose, `/reset-state` to recover
4. **TODO State Loss**: Follow compaction recovery protocol in commands
5. **Line Count Failures**: Ensure project line counter tool exists and is executable

### Emergency Procedures

1. **System Unresponsive**: Use `/check-status` for comprehensive diagnosis
2. **State Corruption**: Use `/reset-state` with appropriate level
3. **Agent Confusion**: Verify identity and re-read agent-specific configurations
4. **Context Loss**: Follow compaction recovery protocols automatically triggered
5. **Progress Loss**: Check backup branches and TODO state files

## Contributing

When modifying these commands:

1. Maintain project-agnostic design with `[project]` placeholders
2. Preserve all critical safety checks and validations
3. Keep state machine integration intact
4. Test with multiple project types
5. Update this documentation

These commands form the operational backbone of Software Factory 2.0, ensuring consistent, high-quality software development processes across all projects.