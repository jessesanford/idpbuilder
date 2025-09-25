# Software Factory 2.0 - Initialization State Machine

## CRITICAL: Project Initialization from Idea to Ready-to-Run

**PURPOSE**: Transform an initial idea/description into a fully configured SF2.0 project ready for `/continue-orchestrating`

**ENTRY COMMAND**: `/init-software-factory`

**TIME TARGET**: < 30 minutes (vs 2-3 hours manual setup)

## Core Initialization States

### Phase 1: INIT_START
**Agent**: Orchestrator
**Purpose**: Capture initial project information
**Actions**:
1. Validate project name/prefix (alphanumeric, no spaces, no dashes)
2. Record initial idea/description
3. Create basic directory structure
4. Initialize state file for initialization process

**Inputs Required**:
- Project name/prefix (becomes $PROJECT_PREFIX)
- Initial idea or change description (1-3 sentences)

**Transition**: → INIT_LOAD_EXAMPLES

---

### Phase 2: INIT_LOAD_EXAMPLES
**Agent**: Architect
**Purpose**: Load templates and examples to understand requirements structure
**Actions**:
1. Read example IMPLEMENTATION-PLAN.md files
2. Read setup-config.yaml examples
3. Read target-repo-config.yaml examples
4. Understand required fields and structure
5. Prepare question framework

**Files to Read**:
- `/home/vscode/software-factory-template/templates/IMPLEMENTATION-PLAN-EXAMPLE.md`
- `/home/vscode/software-factory-template/setup-config-example.yaml`
- `/home/vscode/software-factory-template/target-repo-config.yaml`

**Transition**: → INIT_REQUIREMENTS_GATHERING

---

### Phase 3: INIT_REQUIREMENTS_GATHERING
**Agent**: Architect (Interactive Q&A)
**Purpose**: Gather all project requirements through structured questions
**Actions**:
1. Ask about target codebase (existing vs new)
2. Determine git repository configuration
3. Identify technology stack and frameworks
4. Understand architecture patterns
5. Collect deployment requirements
6. Document quality constraints

**Question Categories**:
1. **Target Codebase**
   - Existing codebase or new project?
   - Fork URL if existing
   - New repository name if creating

2. **Technology Stack**
   - Programming language(s)
   - Key libraries/frameworks
   - Build system
   - Testing framework

3. **Architecture**
   - Pattern (microservices, monolith, library, CLI)
   - Design patterns to follow
   - Compatibility requirements

4. **Deployment**
   - Target environment
   - Container requirements
   - CI/CD preferences

5. **Quality Requirements**
   - Performance targets
   - Security considerations
   - Compliance needs

**Output**: Requirements document with all answers

**Transition**: → INIT_REPO_DECISION

---

### Phase 4: INIT_REPO_DECISION
**Agent**: Orchestrator
**Purpose**: Determine repository setup path based on requirements
**Actions**:
1. Analyze requirements for repository type
2. Choose path:
   - UPSTREAM_FORK: Has existing codebase with fork
   - NEW_REPOSITORY: Creating new project
   - LIBRARY_ONLY: Documentation/library project

**Decision Logic**:
```
IF has_existing_codebase AND has_fork_url:
    → INIT_SETUP_UPSTREAM_FORK
ELIF wants_new_repository:
    → INIT_CREATE_NEW_REPO
ELSE:
    → INIT_LIBRARY_PROJECT
```

---

### Phase 5A: INIT_SETUP_UPSTREAM_FORK
**Agent**: Software Engineer
**Purpose**: Configure forked repository structure
**Actions**:
1. Clone upstream to `$CLAUDE_PROJECT_DIR/efforts/$PROJECT_PREFIX/target-repo/`
2. Set up remote tracking
3. Create development branches
4. Initialize .gitignore if needed

**Transition**: → INIT_GENERATE_CONFIGS

---

### Phase 5B: INIT_CREATE_NEW_REPO
**Agent**: Software Engineer
**Purpose**: Initialize new repository
**Actions**:
1. Create directory at `$CLAUDE_PROJECT_DIR/efforts/$PROJECT_PREFIX/project/`
2. Initialize git repository
3. Create initial .gitignore
4. Make initial commit
5. Set up branch structure

**Transition**: → INIT_GENERATE_CONFIGS

---

### Phase 5C: INIT_LIBRARY_PROJECT
**Agent**: Software Engineer
**Purpose**: Set up documentation/library project
**Actions**:
1. Create project directory structure
2. Initialize as library project (no target repo)
3. Set up documentation directories

**Transition**: → INIT_GENERATE_CONFIGS

---

### Phase 6: INIT_GENERATE_CONFIGS
**Agent**: Code Reviewer
**Purpose**: Generate all configuration files
**Actions**:
1. Create `setup-config.yaml` with:
   - project_name
   - project_prefix
   - target_language
   - key_libraries
   - build_system
   - test_framework
   - deployment_target

2. Create `target-repo-config.yaml` with:
   - upstream_url (if applicable)
   - fork_url (if applicable)
   - main_branch
   - target_directories
   - excluded_paths

3. Create initial `.claude/CLAUDE.md` with:
   - Project-specific instructions
   - Technology guidelines
   - Coding standards

**Transition**: → INIT_SYNTHESIZE_PLAN

---

### Phase 7: INIT_SYNTHESIZE_PLAN
**Agent**: Architect
**Purpose**: Create comprehensive IMPLEMENTATION-PLAN.md
**Actions**:
1. Read all gathered requirements
2. Structure into phases:
   - Phase 1: Foundation/MVP
   - Phase 2: Core Features
   - Phase 3: Polish/Optimization
3. Define waves within each phase
4. Set success criteria
5. Document risk mitigation

**Plan Structure**:
```markdown
# IMPLEMENTATION PLAN: [Project Name]

## 1. Project Overview
[Expanded from initial idea]

## 2. Goals and Objectives
[Refined from requirements]

## 3. Technical Architecture
[From technology decisions]

## 4. Implementation Phases

### Phase 1: Foundation
#### Wave 1.1: Core Setup
- Effort 1.1.1: [Description]
- Effort 1.1.2: [Description]

### Phase 2: Core Features
[Similar structure]

### Phase 3: Enhancement
[Similar structure]

## 5. Success Criteria
[Measurable goals]

## 6. Risk Mitigation
[Identified risks and mitigation]
```

**Transition**: → INIT_CUSTOMIZE_AGENTS

---

### Phase 8: INIT_CUSTOMIZE_AGENTS
**Agent**: Software Engineer
**Purpose**: Add minimal language/tech expertise to agent files
**Actions**:
1. Edit `.claude/agents/sw-engineer.md`:
   - Add language expertise section (< 5 lines)
   - Add framework knowledge (< 5 lines)

2. Edit `.claude/agents/code-reviewer.md`:
   - Add language-specific review criteria (< 5 lines)
   - Add framework best practices (< 5 lines)

**Constraints**:
- MINIMAL modifications only
- Preserve ALL existing rules
- Only add expertise sections

**Transition**: → INIT_VALIDATION

---

### Phase 9: INIT_VALIDATION
**Agent**: Architect
**Purpose**: Comprehensive validation of initialization
**Actions**:
1. Verify all required files exist:
   - [ ] IMPLEMENTATION-PLAN.md
   - [ ] setup-config.yaml
   - [ ] target-repo-config.yaml
   - [ ] .claude/CLAUDE.md
   - [ ] Agent customizations

2. Validate content completeness:
   - [ ] All config fields populated
   - [ ] Plan has all phases defined
   - [ ] Git repositories accessible
   - [ ] Directory structure correct

3. Test readiness:
   - [ ] Can transition to normal operation
   - [ ] State file properly configured

**Transition**: → INIT_HANDOFF (if valid) or INIT_ERROR_RECOVERY (if issues)

---

### Phase 10: INIT_HANDOFF
**Agent**: Orchestrator
**Purpose**: Transition to normal SF2.0 operation
**Actions**:
1. Create production `orchestrator-state.json`:
   - Set current_state: "INIT"
   - Set current_phase: "1"
   - Set current_wave: "1"
   - Clear initialization markers

2. Generate summary report:
   - Project configuration summary
   - Files created
   - Next steps

3. Display instructions:
   ```
   ✅ Project initialization complete!

   Project: [name]
   Type: [upstream_fork|new_repo|library]
   Language: [language]

   Files created:
   - IMPLEMENTATION-PLAN.md
   - setup-config.yaml
   - target-repo-config.yaml
   - orchestrator-state.json

   Next step: Run /continue-orchestrating to begin Phase 1 implementation
   ```

**Terminal State**: Project ready for normal operation

---

### Error Recovery: INIT_ERROR_RECOVERY
**Agent**: Orchestrator
**Purpose**: Handle initialization failures
**Actions**:
1. Identify specific failure:
   - Missing information
   - Invalid git URL
   - File write failures
   - Validation failures

2. Determine recovery path:
   - Re-ask specific questions
   - Retry operations
   - Request manual intervention

3. Log issues and recovery attempts

**Transitions**:
- → INIT_REQUIREMENTS_GATHERING (for missing info)
- → INIT_REPO_DECISION (for repo issues)
- → INIT_VALIDATION (after fixes)

---

## State Transition Rules

### Valid Transitions
```
INIT_START → INIT_LOAD_EXAMPLES
INIT_LOAD_EXAMPLES → INIT_REQUIREMENTS_GATHERING
INIT_REQUIREMENTS_GATHERING → INIT_REPO_DECISION
INIT_REPO_DECISION → {INIT_SETUP_UPSTREAM_FORK, INIT_CREATE_NEW_REPO, INIT_LIBRARY_PROJECT}
INIT_SETUP_UPSTREAM_FORK → INIT_GENERATE_CONFIGS
INIT_CREATE_NEW_REPO → INIT_GENERATE_CONFIGS
INIT_LIBRARY_PROJECT → INIT_GENERATE_CONFIGS
INIT_GENERATE_CONFIGS → INIT_SYNTHESIZE_PLAN
INIT_SYNTHESIZE_PLAN → INIT_CUSTOMIZE_AGENTS
INIT_CUSTOMIZE_AGENTS → INIT_VALIDATION
INIT_VALIDATION → {INIT_HANDOFF, INIT_ERROR_RECOVERY}
INIT_ERROR_RECOVERY → {INIT_REQUIREMENTS_GATHERING, INIT_REPO_DECISION, INIT_VALIDATION}
INIT_HANDOFF → [END]
```

### Mandatory Stop Points
Per R313, stops required after:
- INIT_REQUIREMENTS_GATHERING (after Q&A session)
- INIT_REPO_DECISION (after path selection)
- INIT_VALIDATION (after validation complete)

## Critical Requirements

### 1. Interactive Q&A System
- Architect MUST read example files first
- Questions must be clear and specific
- Accept user responses and validate
- Build complete requirements document

### 2. File Generation
- All configs must have complete information
- No placeholder values
- Proper YAML structure
- Git-ready formatting

### 3. Minimal Agent Customization
- Maximum 10 lines per agent file
- Only add expertise sections
- Preserve all existing rules
- Use additive modifications only

### 4. Validation Completeness
- Every required file must exist
- All fields must be populated
- Git repositories must be accessible
- Must be ready for /continue-orchestrating

## Success Metrics

- **Time**: < 30 minutes from start to handoff
- **Completeness**: 100% of required files generated
- **Validity**: Passes all validation checks
- **Usability**: Can immediately run /continue-orchestrating
- **Quality**: Generated plan sufficient for Phase 1 start

## Error Handling

### Common Failures and Recovery

1. **Invalid Git URL**
   - Detection: Git clone fails
   - Recovery: Re-ask for correct URL
   - Validation: Test clone before proceeding

2. **Missing Requirements**
   - Detection: Validation finds empty fields
   - Recovery: Return to Q&A for specific items
   - Prevention: Validate answers immediately

3. **File Write Failures**
   - Detection: Write operations fail
   - Recovery: Check permissions, retry
   - Fallback: Request manual directory creation

4. **Network Failures**
   - Detection: Clone/fetch operations timeout
   - Recovery: Exponential backoff retry
   - Fallback: Allow manual repository setup

## Integration with Normal Operation

After INIT_HANDOFF completes:
1. System is in standard INIT state
2. orchestrator-state.json is production-ready
3. All configs are complete
4. User runs `/continue-orchestrating`
5. Normal state machine takes over from SOFTWARE-FACTORY-STATE-MACHINE.md

## State File Locations

Each initialization state has corresponding rule files:

```
agent-states/
└── initialization/
    ├── orchestrator/
    │   ├── INIT_START/rules.md
    │   ├── INIT_REPO_DECISION/rules.md
    │   ├── INIT_HANDOFF/rules.md
    │   └── INIT_ERROR_RECOVERY/rules.md
    ├── architect/
    │   ├── INIT_LOAD_EXAMPLES/rules.md
    │   ├── INIT_REQUIREMENTS_GATHERING/rules.md
    │   ├── INIT_SYNTHESIZE_PLAN/rules.md
    │   └── INIT_VALIDATION/rules.md
    ├── sw-engineer/
    │   ├── INIT_CREATE_NEW_REPO/rules.md
    │   ├── INIT_SETUP_UPSTREAM_FORK/rules.md
    │   ├── INIT_LIBRARY_PROJECT/rules.md
    │   └── INIT_CUSTOMIZE_AGENTS/rules.md
    └── code-reviewer/
        └── INIT_GENERATE_CONFIGS/rules.md
```

## Notes

- This state machine is SEPARATE from the main state machine
- Uses "INIT_" prefix to distinguish initialization states
- Focused on setup, not implementation
- Interactive and user-friendly
- Generates production-ready configuration
- State rules are in `agent-states/initialization/` to prevent conflicts