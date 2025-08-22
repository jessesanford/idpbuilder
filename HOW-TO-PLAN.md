# How to Plan a Software Project Using the Software Factory System

## Overview

This guide provides a step-by-step process for planning any software project using the Software Factory orchestration system. The process involves creating increasingly detailed plans that can be executed by specialized agents.

## Planning Phases

### Phase 1: High-Level Implementation Plan
**Target Agent**: @agent-architect-reviewer or @agent-code-reviewer (Project Maintainer)

#### Prompt Template:
```
Create a detailed step-by-step phased plan to implement [PROJECT/FEATURE REQUIREMENTS DOCUMENT] 
by an @agent-orchestrator-task-master who will task out each of the phases 1 by 1 to a group 
of @agent-sw-engineer agents.

Each Phase of your plan should have waves of steps that are logically grouped together around 
features/functionality that make sense to group with each other. The "waves" are like "sprints" 
where inside of each "wave" there are "efforts" which are like "user stories" but may or may not 
be "user oriented" - they could just be "maintenance tasks" or other types of tasks.

Requirements:
1. Each effort will only ever be worked on by 1 @agent-sw-engineer
2. EACH EFFORT WILL GET ITS OWN BRANCH: /phase{X}/wave{Y}/effort{Z}-{effort-name}
3. ALL WORK ON THE EFFORT MUST BE DONE IN THAT BRANCH
4. Specify whether efforts within a wave should be SERIAL or PARALLEL
5. CREATE A DEPENDENCY MAPPING TABLE AND GRAPH for each wave
6. SPECIFY WHICH BRANCHES should be based on which other branches
7. Order waves: APIs/INTERFACES/SCHEMAS → SHARED LIBRARIES → IMPLEMENTATIONS → TESTING → DOCUMENTATION
8. Each agent works in their own sparse clone with only needed branches
9. Include a comprehensive MERGE PLAN from efforts → waves → phases → main
```

#### Expected Output:
```markdown
# [Project] Implementation Plan

## Executive Summary
[Overview of phases, waves, and efforts]

## Phase Structure
### Phase 1: [Foundation & Contracts]
- Duration: X days
- Waves: Y
- Efforts: Z
- Dependencies: None

### Phase 2: [Core Infrastructure]
- Duration: X days
- Dependencies: Phase 1

[... continue for all phases ...]

## Wave Dependency Graphs
[Visual representations and tables]

## Merge Strategy
[Detailed merge plan]
```

### Phase 2: Phase-Specific Detailed Plans
**Target Agent**: @agent-architect-reviewer or @agent-code-reviewer (Project Maintainer)

#### Prompt Template:
```
Take the high-level plan and create MORE SPECIFIC plans for each phase that spell out the 
efforts within them in great detail.

For each effort, include:
1. Exactly which existing branches/code to use as fodder
2. Specific commits to cherry-pick (if reusing code)
3. PSEUDO CODE for implementation approach (without writing actual code)
4. TDD TEST REQUIREMENTS - write actual test cases
5. Clear, specific requirements for completion
6. Success criteria and validation steps

Focus on:
- PRIORITIZING CODE REUSE from existing work
- Being explicit enough that @agent-sw-engineer knows EXACTLY what to do
- Including validation commands after each step

Put each plan in: /idpbuilder/phase-plans/PHASE{X}-SPECIFIC-IMPL-PLAN.md
```

#### Expected Output Structure:
```markdown
# Phase X: [Name] - Detailed Implementation Plan

## Phase Overview
- Duration: X days
- Critical Path: YES/NO
- Base Branch: [branch]
- Target Integration Branch: phase{X}-integration

## Wave X.1: [Wave Name]

### Effort X.1.1: [Effort Name]
**Branch**: /phase{X}/wave{1}/effort{1}-{name}
**Duration**: X hours
**Dependencies**: None

#### Source Material:
- Primary: origin/feature/[existing-branch]
- Secondary: origin/feature/[fallback-branch]

#### Specific Commits to Cherry-Pick:
```bash
# From existing work
git cherry-pick abc123def  # Feature implementation
git cherry-pick 456ghi789  # Bug fixes
```

#### Requirements:
1. MUST implement [specific feature]
2. MUST include [specific components]
3. MUST NOT [specific restrictions]

#### Test Requirements (TDD):
```[language]
// Test case 1: [Description]
func Test[Feature](t *testing.T) {
    // Given
    [setup code]
    
    // When
    [action code]
    
    // Then
    assert.Equal(t, expected, actual)
}
```

#### Pseudo-Code Implementation:
```
FUNCTION implement_feature():
    // Step 1: Setup base structure
    CREATE directory structure
    COPY templates from existing
    
    // Step 2: Core implementation
    FOR each component IN requirements:
        IF exists_in_source_branch:
            CHERRY_PICK component
        ELSE:
            IMPLEMENT minimal_version
    
    // Step 3: Integration
    WIRE components together
    VALIDATE interfaces
```

#### Validation Commands:
```bash
# After implementation
make build || exit 1
make test || exit 1
make lint || exit 1

# Verify line count
/tools/line-counter.sh -c [branch]
```

#### Success Criteria:
- [ ] All tests pass
- [ ] Build succeeds
- [ ] Line count < 800
- [ ] No lint errors
- [ ] Documentation updated
```

### Phase 3: Validation and Enhancement Review
**Target Agent**: @agent-architect-reviewer

#### Prompt Template:
```
Review all phase plans (PHASE1 through PHASEX) and evaluate:

1. Are they explicit enough for @agent-sw-engineer to know EXACTLY what to do?
2. Do they contain ENOUGH pseudocode and TDD tests?
3. Is there sufficient detail to ensure:
   - Clean, feature-complete output
   - Buildable code at end of every wave
   - All tests passing at end of every phase
4. What additional details would help ensure EXACTLY CORRECT implementation?

Create any needed enhancement documents for clarity.
```

#### Expected Enhancement Output:
```markdown
# Agent Explicit Instructions Enhancement

## Critical Instructions for ALL Agents

### Git Commands - EXACT SYNTAX
[Specific commands, not pseudo-code]

### Validation Gates
[Specific checks after each operation]

### Rollback Procedures
[Recovery steps if efforts fail]

### Success Criteria
[Measurable outcomes]
```

## Additional Planning Documents

### 1. Size Management and Split Planning
**Target Agent**: @agent-code-reviewer

Create instructions for handling efforts that exceed size limits:

```markdown
# Split Review Loop Process

## When effort > 800 lines:
1. Create split plan
2. Implement splits sequentially
3. Review each split
4. If still > limit, create sub-splits
5. Continue until all branches acceptable

## Exception Criteria:
Only @agent-code-reviewer can grant exceptions when:
- Atomic transactions cannot be split
- Complex state machines
- Tightly coupled interface/implementation
```

### 2. Test-Driven Validation Requirements
**Target Agent**: @agent-code-reviewer or @agent-architect-reviewer

```markdown
# Test Coverage Requirements by Phase

## Phase 1: APIs & Contracts (90% coverage)
## Phase 2: Core Infrastructure (85% coverage)
## Phase 3: Implementation (80% coverage)
## Phase 4: Features (85% coverage)
## Phase 5: Integration (95% coverage)
```

## Complete Planning Workflow

### Step 1: Initial Planning Session
1. **Input**: Requirements document, existing codebase analysis
2. **Agent**: @agent-architect-reviewer
3. **Output**: High-level implementation plan with phases/waves/efforts

### Step 2: Detailed Phase Planning
1. **Input**: High-level plan, existing code branches
2. **Agent**: @agent-code-reviewer (Project Maintainer)
3. **Output**: Phase-specific plans with explicit instructions

### Step 3: Validation Review
1. **Input**: All phase plans
2. **Agent**: @agent-architect-reviewer
3. **Output**: Enhancement documents, validation requirements

### Step 4: Orchestration Preparation
1. **Input**: Complete plans and enhancements
2. **Agent**: @agent-orchestrator-task-master
3. **Output**: orchestrator-state.yaml, ready to execute

## Example Project Structure After Planning

```
/home/vscode/workspaces/idpbuilder/
├── orchestrator/
│   ├── PROJECT-IMPLEMENTATION-PLAN.md          # High-level plan
│   ├── orchestrator-state.yaml                 # Execution state
│   └── ORCHESTRATOR-TASKMASTER-EXECUTION-PLAN.md
├── phase-plans/
│   ├── PHASE1-SPECIFIC-IMPL-PLAN.md           # Detailed plans
│   ├── PHASE2-SPECIFIC-IMPL-PLAN.md
│   ├── PHASE3-SPECIFIC-IMPL-PLAN.md
│   └── ...
├── protocols/
│   ├── AGENT-EXPLICIT-INSTRUCTIONS.md         # Enhancement docs
│   ├── TEST-DRIVEN-VALIDATION-REQUIREMENTS.md
│   └── SPLIT-REVIEW-LOOP-DIAGRAM.md
└── efforts/                                    # Working directories
    └── (created during execution)
```

## Key Principles for Effective Planning

### 1. Prioritize Code Reuse
- Always check for existing implementations
- Prefer cherry-picking over rewriting
- Document source branches explicitly

### 2. Be Explicit, Not Abstract
- Use exact git commands, not pseudo-code
- Provide actual test cases, not descriptions
- Include specific file paths and structures

### 3. Plan for Iteration
- Expect split branches for large efforts
- Build in review loops
- Allow for exceptions with justification

### 4. Maintain Isolation
- Each effort in its own branch
- Each agent in its own working directory
- Clear merge strategy to prevent conflicts

### 5. Validate Continuously
- Tests after every effort
- Builds after every wave
- Integration tests after every phase

## Common Planning Patterns

### API-First Development
```
Phase 1: APIs and Contracts
├── Wave 1: Core Types
├── Wave 2: Interfaces
└── Wave 3: Schemas/Contracts

Phase 2: Implementation
├── Wave 1: Core Logic
├── Wave 2: Controllers
└── Wave 3: Integration
```

### Library Development
```
Phase 1: Public API
├── Wave 1: Interfaces
└── Wave 2: Types

Phase 2: Core Implementation
├── Wave 1: Basic Features
└── Wave 2: Advanced Features

Phase 3: Testing & Documentation
├── Wave 1: Unit Tests
├── Wave 2: Integration Tests
└── Wave 3: Documentation
```

### Microservice Development
```
Phase 1: Service Contracts
├── Wave 1: API Definitions
└── Wave 2: Message Formats

Phase 2: Service Implementation
├── Wave 1: Business Logic
├── Wave 2: Data Layer
└── Wave 3: API Layer

Phase 3: Integration
├── Wave 1: Service Communication
└── Wave 2: End-to-end Testing
```

## Validation Checklist

Before starting orchestration, ensure:

- [ ] All phases have detailed implementation plans
- [ ] Every effort has explicit success criteria
- [ ] Test requirements defined for each effort
- [ ] Dependency graphs are complete and accurate
- [ ] Merge strategy is documented
- [ ] Line count limits and exceptions documented
- [ ] Existing code reuse opportunities identified
- [ ] Validation commands specified for each step
- [ ] Rollback procedures defined
- [ ] Agent working directories planned

## Tips for Success

1. **Start with existing code analysis** - Know what you can reuse
2. **Front-load contracts** - APIs and interfaces first
3. **Plan for splits** - Assume some efforts will exceed limits
4. **Document exceptions** - Clear criteria for when rules can bend
5. **Test everything** - TDD approach with specific test cases
6. **Be specific** - Agents execute exactly what's written
7. **Plan the merge** - Know how everything comes together

This planning process ensures that complex software projects can be successfully orchestrated through the Software Factory system with multiple specialized agents working in coordination.