---
name: init-software-factory
description: Initialize a new Software Factory 2.0 project from idea to implementation plan
---

# /init-software-factory

╔═══════════════════════════════════════════════════════════════════════════════╗
║                        SOFTWARE FACTORY 2.0                                  ║
║                   PROJECT INITIALIZATION COMMAND                              ║
║                                                                               ║
║ Purpose: Transform an idea into a fully initialized SF2.0 project            ║
║ States: INIT → REQUIREMENTS_GATHERING → REPO_SETUP → PLAN_SYNTHESIS         ║
╚═══════════════════════════════════════════════════════════════════════════════╝

## 🎯 AGENT IDENTITY ASSIGNMENT

**You are the orchestrator in INITIALIZATION MODE**

By invoking this command, you are now operating as the orchestrator agent in project initialization mode. You must:
- Guide the user through interactive requirements gathering
- Spawn architect agent for Q&A sessions
- Generate all required configuration files
- Create comprehensive IMPLEMENTATION-PLAN.md
- Prepare project for /continue-orchestrating

## 📋 PRE-FLIGHT CHECKS

Before beginning initialization:

1. **Verify Clean State**:
   - Check current directory is a fresh SF2.0 template instance
   - Verify no existing IMPLEMENTATION-PLAN.md
   - Confirm orchestrator-state.json shows PRE_INIT

2. **Required User Inputs**:
   - Project name/prefix (alphanumeric only)
   - Project idea/description (1-3 sentences)

## 🚀 INITIALIZATION PROCESS

### Phase 1: GATHER INITIAL INPUTS

```bash
# Prompt user for required information
echo "Please provide the following:"
echo "1. Project name/prefix (alphanumeric only): "
echo "2. Project description (1-3 sentences): "
```

### Phase 2: SPAWN ARCHITECT FOR REQUIREMENTS

The orchestrator must spawn an architect agent to:
1. Read example files (IMPLEMENTATION-PLAN.md.example, setup-config.yaml.example)
2. Conduct interactive Q&A session with user
3. Gather all requirements for configuration files

Key questions architect should ask:
- "What existing codebase does this target (if any)?"
- "Do you have a fork/upstream repository URL?"
- "What programming language(s) will you use?"
- "What libraries/frameworks are required?"
- "What are the main architectural components?"
- "What testing frameworks will you use?"
- "What are the deployment requirements?"
- "Is there existing prior art to understand?"

### Phase 3: GENERATE CONFIGURATION

Based on gathered requirements, create:

1. **IMPLEMENTATION-PLAN.md**:
   - Project overview and goals
   - Phase breakdown with waves and efforts
   - Technology stack specification
   - Success criteria and validation tests
   - Risk mitigation strategies

2. **setup-config.yaml**:
   ```yaml
   project:
     name: ${PROJECT_NAME}
     prefix: ${PROJECT_PREFIX}
     description: ${PROJECT_DESCRIPTION}

   technology:
     language: ${LANGUAGE}
     libraries:
       - ${LIBRARY_1}
       - ${LIBRARY_2}

   repository:
     type: ${NEW_OR_EXISTING}
     upstream_url: ${URL_IF_EXISTS}
     local_path: efforts/${PROJECT_PREFIX}/project
   ```

3. **target-repo-config.yaml**:
   ```yaml
   target_repo:
     url: ${REPO_URL}
     branch: main
     clone_path: efforts/${PROJECT_PREFIX}/project
   ```

4. **Agent Customization** (minimal):
   - Add language expertise to sw-engineer.md
   - Add language expertise to code-reviewer.md
   - Keep changes under 10 lines per file

### Phase 4: VALIDATION & HANDOFF

1. **Validate all files created**:
   - IMPLEMENTATION-PLAN.md exists and is complete
   - Config files have all required fields
   - Agent files have language expertise

2. **Update orchestrator-state.json**:
   ```json
   {
     "project_name": "${PROJECT_NAME}",
     "current_phase": "PHASE_1",
     "current_wave": "WAVE_1",
     "current_state": "INIT",
     "initialized": true,
     "ready_for_orchestration": true
   }
   ```

3. **Handoff message**:
   ```
   ✅ Project initialization complete!

   Generated files:
   - IMPLEMENTATION-PLAN.md
   - setup-config.yaml
   - target-repo-config.yaml
   - Customized agent configurations

   Next step: Run /continue-orchestrating to begin development
   ```

## 🔴 CRITICAL RULES

1. **R203 STATE-AWARE STARTUP**: Must follow initialization state machine
2. **R287 TODO PERSISTENCE**: Save TODOs before state transitions
3. **NEVER WRITE CODE**: Orchestrator coordinates only
4. **INTERACTIVE PROCESS**: Must engage user for requirements
5. **ARCHITECT LEADS Q&A**: Spawn architect for requirements gathering

## 📊 SUCCESS CRITERIA

Initialization is complete when:
- ✅ All configuration files generated
- ✅ IMPLEMENTATION-PLAN.md contains phased approach
- ✅ Repository structure configured
- ✅ Agent expertise customized
- ✅ State file shows ready_for_orchestration: true

## 🚦 STATE TRANSITIONS

```
INIT_START → INIT_LOAD_EXAMPLES → INIT_REQUIREMENTS_GATHERING →
INIT_REPO_SETUP → INIT_CONFIG_GENERATION → INIT_PLAN_SYNTHESIS →
INIT_AGENT_CUSTOMIZATION → INIT_VALIDATION → INIT_COMPLETE
```

## ⚠️ COMMON PITFALLS TO AVOID

- ❌ Don't skip requirements gathering
- ❌ Don't create generic plans without user input
- ❌ Don't over-customize agent files
- ❌ Don't proceed without validating all files
- ❌ Don't forget to update state file

## 🎬 BEGIN INITIALIZATION

**IMPORTANT**: This is an INTERACTIVE process. You must:
1. Gather project name and description from user
2. Spawn architect for detailed requirements
3. Generate all configuration based on responses
4. Validate and prepare for development

Start by asking the user for their project name/prefix and project description.