# Complete File Tree for Software Factory Template

```
software-factory-template/
│
├── README.md                                    # Main documentation
├── setup.sh                                     # Interactive setup script
├── orchestrator-state-example.yaml              # Example state file
├── PROJECT-IMPLEMENTATION-PLAN-TEMPLATE.md      # High-level plan template
├── HOW-TO-PLAN.md                               # Planning methodology guide
├── PLANNING-AGENT-ASSIGNMENTS.md                # Agent responsibility guide
├── TEMPLATE-CREATION-SUMMARY.md                 # Template creation details
├── CRITICAL-FILES-ADDED.md                      # Critical files documentation
├── FINAL-FILE-ORGANIZATION.md                   # File organization summary
├── CLAUDE-MD-FILE-VERIFICATION.md               # File reference verification
│
├── .claude/                                      # Claude Code configuration
│   ├── CLAUDE.md                                # Global rules and recovery procedures
│   ├── agents/                                  # Agent configurations
│   │   ├── orchestrator-task-master.md         # Orchestrator agent config
│   │   ├── code-reviewer.md                    # Code reviewer agent config
│   │   ├── architect-reviewer.md               # Architect agent config
│   │   └── sw-engineer-example-go.md           # SW engineer config (customize)
│   └── commands/                                # Custom commands
│       └── continue-orchestrating.md           # Main orchestration command
│
├── core/                                        # Core system files
│   ├── SOFTWARE-FACTORY-STATE-MACHINE.md       # State machine definition
│   └── ORCHESTRATOR-MASTER-OPERATIONS-GUIDE.md # Operations blueprint
│
├── protocols/                                   # Required protocol files (15 total)
│   ├── IMPERATIVE-LINE-COUNT-RULE.md           # Size limit enforcement
│   ├── EFFORT-SPLIT-CONTINUOUS-EXECUTION-PROTOCOL.md # Split handling
│   ├── SW-ENGINEER-STARTUP-REQUIREMENTS.md     # SW engineer startup
│   ├── SW-ENGINEER-EXPLICIT-INSTRUCTIONS.md    # SW engineer instructions
│   ├── ORCHESTRATOR-EFFORT-PLANNING-PROTOCOL.md # Planning coordination
│   ├── ORCHESTRATOR-TASKMASTER-EXECUTION-PLAN.md # Execution guide
│   ├── CODE-REVIEWER-EFFORT-PLANNING-INSTRUCTIONS.md # Plan creation
│   ├── CODE-REVIEWER-COMPREHENSIVE-GUIDE.md    # Review process
│   ├── WAVE-COMPLETION-ARCHITECT-REVIEW-PROTOCOL.md # Wave reviews
│   ├── ARCHITECT-REVIEWER-WAVE-INSTRUCTIONS.md # Architect instructions
│   ├── PHASE-START-ARCHITECT-REVIEW-PROTOCOL.md # Phase assessment
│   ├── TEST-DRIVEN-VALIDATION-REQUIREMENTS.md  # Testing requirements
│   ├── WORK-LOG-TEMPLATE.md                    # Progress tracking
│   └── TODO-STATE-MANAGEMENT-PROTOCOL.md       # TODO persistence
│
├── phase-plans/                                 # Phase planning templates
│   ├── README.md                               # How to use templates
│   ├── PHASEX-GENERIC-TEMPLATE.md             # Generic phase template
│   ├── PHASE1-TEMPLATE.md                     # API/Contract example
│   ├── PHASE2-TEMPLATE.md                     # Infrastructure example
│   └── PHASE3-TEMPLATE.md                     # Implementation example
│
├── efforts/                                     # Working directories
│   └── README.md                               # Effort isolation guide
│
├── todos/                                       # TODO state files
│   └── README.md                               # TODO recovery guide
│
├── tools/                                       # Utility scripts
│   └── line-counter.sh                         # Line counting tool
│
└── possibly-needed-but-not-sure/               # Optional enhancements
    ├── README.md                               # Optional files guide
    ├── WHEN-TO-USE-THESE-FILES.md             # Decision guide
    ├── FILES-NOT-INCLUDED.md                  # Omitted files explanation
    ├── CODE-REVIEWER-QUICK-REFERENCE.md       # Quick review guide
    ├── ORCHESTRATOR-QUICK-REFERENCE.md        # Quick orchestrator guide
    ├── ORCHESTRATOR-WORKFLOW-SUMMARY.md       # Visual workflow
    ├── CODE-REVIEW-EXAMPLES.md                # Example reviews
    ├── CODE-REVIEW-ENFORCEMENT-SUMMARY.md     # Enforcement points
    ├── ORCHESTRATOR-CODE-REVIEW-INTEGRATION.md # Review integration
    ├── ORCHESTRATOR-NEVER-WRITES-CODE-RULE.md # Redundant rule
    ├── PHASE-COMPLETION-FUNCTIONAL-TESTING.md # Testing protocols
    ├── SPLIT-EXAMPLE-E3.1.1-SYNC-ENGINE.md   # Split example
    └── SPLIT-REVIEW-LOOP-DIAGRAM.md          # Split process diagram
```

## File Count Summary

- **Root level**: 10 files
- **.claude/**: 1 CLAUDE.md + 4 agents + 1 command = 6 files
- **core/**: 2 files
- **protocols/**: 15 files
- **phase-plans/**: 5 files
- **efforts/**: 1 README
- **todos/**: 1 README
- **tools/**: 1 script
- **possibly-needed-but-not-sure/**: 13 files

**Total**: 54 files