# Phase Planning Templates

## Overview

This directory contains THE template for creating detailed phase-specific implementation plans. These plans bridge the gap between high-level project planning and actual code implementation by providing explicit, actionable instructions for developer agents.

## THE Template

### PHASE-IMPL-PLAN-TEMPLATE.md
**This is the ONLY template to use for phase planning**  
**Purpose:** Creates comprehensive phase plans with library decisions, reuse enforcement, and critical code snippets  
**Created by:** Code Reviewer acting as Senior Maintainer  
**Contains:** Complete structure with maintainer guidance sections

## Example Phase Plans

### PHASE1-TEMPLATE.md  
**Example of:** Foundation phase (APIs, contracts, schemas)  
**Shows:** How to structure API-first development

### PHASE2-TEMPLATE.md
**Example of:** Infrastructure phase (controllers, frameworks, libraries)  
**Shows:** Building reusable components

### PHASE3-TEMPLATE.md
**Example of:** Implementation phase (business logic, features)  
**Shows:** Complex implementation with potential splits

## How to Create a Phase Plan

### Step 1: Orchestrator Tasks Code Reviewer as Senior Maintainer
```markdown
Task @agent-code-reviewer:
Act as SENIOR PROJECT MAINTAINER to create Phase [X] detailed implementation plan.
Use phase-plans/PHASE-IMPL-PLAN-TEMPLATE.md as your template.
```

### Step 2: Code Reviewer Creates Phase Plan
The Code Reviewer (as Senior Maintainer) will:

1. **Make Library Decisions**
   ```yaml
   core_libraries:
     - name: "github.com/gorilla/mux"
       version: "v1.8.0"
       reason: "Mature router with middleware support"
   ```

2. **Enforce Reuse from Previous Phases**
   ```yaml
   reused_from_previous:
     phase1:
       - "pkg/api/types.go: All API types"
       - "pkg/auth: Complete auth system"
   
   forbidden_duplications:
     - "DO NOT create new auth system"
     - "DO NOT reimplement validators"
   ```

3. **Provide Critical Code Snippets (10-30 lines ONLY)**
   ```go
   // MAINTAINER NOTE: Complex OAuth refresh logic
   func refreshToken(token *Token) (*Token, error) {
       // [15-20 lines of critical logic]
       // SW Engineer implements surrounding code
   }
   ```

4. **Define Interface Contracts**
   ```go
   // MUST implement this interface from Phase 1
   type DataProcessor interface {
       Process(data []byte) error  // From Phase 1
       Transform() error           // New in this phase
   }
   ```

### Step 3: Structure the Phase Plan

#### Phase Overview Section
- Duration and critical path
- Base branch (previous phase integration)
- Prerequisites from earlier phases
- Target integration branch

#### Library & Dependencies Section (MAINTAINER CRITICAL)
- Exact library versions with justification
- What to reuse from previous phases
- Forbidden duplications list
- Shared dependencies across phases

#### Wave Sections
For each wave:
- Overview and dependencies
- Whether efforts can parallelize
- Individual effort details

#### Effort Details
For each effort:
- Branch naming
- Duration and line estimates
- Source material (if reusing)
- Requirements (MUST/MUST NOT)
- Implementation guidance
- Test requirements
- Validation commands

### Step 4: Key Principles for Maintainer

#### DO Provide:
1. **Specific Library Versions**
   - Not "use a logging library"
   - But "use github.com/sirupsen/logrus v1.9.0"

2. **Critical Algorithm Snippets**
   ```go
   // ONLY the complex 20 lines of distributed locking
   // Not the entire service implementation
   ```

3. **Reuse Enforcement**
   ```markdown
   MUST reuse from Phase 1:
   - pkg/models: ALL database models
   - pkg/validators: ALL validation logic
   ```

4. **Clear Patterns**
   ```markdown
   Pattern: Repository pattern from Phase 2
   Location: pkg/repository/base.go
   Extension: Add caching layer only
   ```

#### DO NOT Provide:
1. **Complete Implementations** - SW Engineers write the code
2. **More than 30 lines per snippet** - Only critical complex parts
3. **Vague Guidance** - Be specific about versions and paths
4. **New Patterns when old work** - Enforce consistency

## Template Sections Explained

### Critical Libraries & Dependencies
**Owner:** Code Reviewer as Maintainer  
**Purpose:** Technical decisions that ensure consistency  
**Contains:** Library versions, reuse mandates, forbidden duplications

### Wave and Effort Structure
**Purpose:** Break down work into manageable chunks  
**Contains:** Dependencies, parallelization options, effort details

### Implementation Guidance
**Purpose:** Guide without over-implementing  
**Contains:** Patterns to follow, critical snippets, interface contracts

### Test Requirements
**Purpose:** Ensure quality and coverage  
**Contains:** TDD test structure, coverage targets

### Validation Commands
**Purpose:** Verify success  
**Contains:** Build, test, lint, line count commands

## Integration with Orchestrator Workflow

1. **Before Phase Start:**
   - Orchestrator spawns Code Reviewer as Senior Maintainer
   - Code Reviewer creates PHASE[X]-SPECIFIC-IMPL-PLAN.md

2. **During Phase Execution:**
   - Orchestrator uses plan to create effort directories
   - Tasks SW Engineers with specific efforts
   - Engineers follow maintainer guidance

3. **Reuse Enforcement:**
   - SW Engineers MUST use specified libraries
   - MUST reuse interfaces from previous phases
   - MUST NOT duplicate existing functionality

## Common Patterns

### API Development Pattern
```
Wave 1: Types and Interfaces (reuse enforced)
Wave 2: Implementation (using Phase 1 types)
Wave 3: Integration (using Phase 2 services)
```

### Feature Development Pattern
```
Wave 1: Core Logic (with maintainer snippets)
Wave 2: Integration (reusing Phase N components)
Wave 3: Optimization (following established patterns)
```

## Validation Checklist

Before considering a phase plan complete:

- [ ] All library versions specified exactly
- [ ] Reuse from previous phases explicitly mandated
- [ ] Forbidden duplications listed
- [ ] Critical code snippets provided (10-30 lines max)
- [ ] Interface contracts defined
- [ ] Every effort has clear requirements
- [ ] Test requirements specified
- [ ] Validation commands included
- [ ] Line count estimates provided
- [ ] Split strategies planned where needed

## Example: Good vs Bad Maintainer Guidance

### ❌ BAD (Too Vague):
```markdown
E2.1.1: Implement Authentication
Use some JWT library.
Implement OAuth flow.
```

### ✅ GOOD (Specific Guidance):
```markdown
E2.1.1: JWT Authentication with OAuth2

Libraries:
- github.com/golang-jwt/jwt/v5 v5.0.0 (JWT handling)
- golang.org/x/oauth2 v0.13.0 (OAuth2 client)

MUST reuse from Phase 1:
- pkg/models/user.go: User model
- pkg/crypto/keys.go: RSA key management

Critical Implementation (OAuth token exchange):
```go
// MAINTAINER: Complex token exchange with retry
func exchangeToken(code string) (*Token, error) {
    // [20 lines of critical retry/backoff logic]
    // SW Engineer implements the rest
}
```

Forbidden:
- DO NOT create new User type
- DO NOT implement new key management
```

## Tips for Success

1. **Maintainer decides, Engineers implement** - Clear separation
2. **Reuse is mandatory** - Not optional
3. **Snippets are for complexity** - Not whole implementations
4. **Versions are exact** - No ambiguity
5. **Patterns are consistent** - Across all phases

Each phase plan becomes the technical contract that ensures consistency, reuse, and quality across the entire project.