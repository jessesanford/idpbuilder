# Architect - INIT_DECOMPOSE_PRD State Rules

## 🔴🔴🔴 ABSOLUTE REQUIREMENT: STATE MANAGER CONSULTATION 🔴🔴🔴

**SUPREME LAW - R517 - UNIVERSAL STATE MANAGER CONSULTATION LAW**

**BEFORE EXITING THIS STATE, YOU MUST:**

1. **SPAWN STATE MANAGER FOR SHUTDOWN_CONSULTATION** (MANDATORY - NO EXCEPTIONS)
2. **NEVER UPDATE STATE FILES DIRECTLY** (orchestrator-state-v3.json, bug-tracking.json, etc.)
3. **NEVER COMMIT STATE FILES YOURSELF** (State Manager does this atomically)

**FAILURE TO CONSULT STATE MANAGER = IMMEDIATE SYSTEM HALT (-100% GRADE)**

### Enforcement Mechanism

If you attempt to exit this state without spawning State Manager:
- ❌ Pre-commit hooks will REJECT your commit
- ❌ Validation tools will FAIL the build
- ❌ Grading system will assign -100% penalty
- ❌ System will transition to ERROR_RECOVERY

### Required Pattern (COPY THIS EXACTLY)

```bash
# At end of state work, BEFORE any state file updates:

echo "🔴 MANDATORY: Spawning State Manager for SHUTDOWN_CONSULTATION"

# Spawn State Manager (REQUIRED - NOT OPTIONAL)
# Task: state-manager
# State: SHUTDOWN_CONSULTATION
# Current State: [YOUR_CURRENT_STATE]
# Proposed Next State: [YOUR_PROPOSED_NEXT_STATE]
# Work Summary: [SUMMARY_OF_WORK_COMPLETED]

# State Manager will:
# 1. Validate proposed transition against state machine
# 2. Update all 4 state files atomically
# 3. Commit with [R288] tag
# 4. Return REQUIRED next state (may differ from proposal)

# Wait for State Manager response
# Follow State Manager's directive (REQUIRED next state)
# DO NOT proceed until State Manager confirms
```

**YOU MUST NEVER:**
- ❌ Update orchestrator-state-v3.json yourself
- ❌ Update bug-tracking.json yourself
- ❌ Update integration-containers.json yourself
- ❌ Use `jq` to modify state files
- ❌ Use `sed/awk` to modify state files
- ❌ Set `validated_by: "orchestrator"` (must be "state-manager")
- ❌ Skip State Manager consultation "just this once"
- ❌ Think "I'll validate it manually"

**ONLY State Manager may update state files. This is NON-NEGOTIABLE.**

See: `rule-library/R517-universal-state-manager-consultation-law.md`
## Purpose
Decompose project requirements and PRD into concrete phases, waves, and efforts with sizing justifications.

## Entry Criteria
- Requirements gathered from user (INIT_REQUIREMENTS_GATHERING complete)
- Initial project description available
- Technical stack determined
- Ready to break down work structure

## Core Responsibility

**THIS IS THE CRITICAL STATE**: Transform a PRD/description into an actionable implementation plan structure.

The user should NEVER have to manually figure out phases/waves/efforts. That's YOUR job as the architect.

## Required Actions

### 1. Load All Requirements

Read from `init-state-temp.json`:
- Initial project description (the PRD)
- Requirements gathered from Q&A
- Technology stack decisions
- Project constraints

### 2. Identify Major Features

Analyze the PRD and extract all major features/capabilities:

**Example PRD**: "Implement idpbuilder push command to upload OCI images to Gitea registry"

**Feature Extraction**:
1. CLI command structure (`push` subcommand)
2. Authentication layer (username/password flags)
3. OCI image operations (go-containerregistry integration)
4. Registry communication (HTTPS, TLS handling)
5. Progress tracking and user feedback
6. Error handling and validation
7. Testing infrastructure

**Process**:
- List every distinct capability mentioned
- Include implied features (auth → credential validation, registry → connection handling)
- Don't forget non-functional requirements (testing, docs, error handling)

### 3. Estimate Lines of Code per Feature

For each feature, estimate implementation size using **effort-sizing-guidelines.md**.

**Use These Heuristics**:

#### CLI Features
- Command skeleton + flags: 100-200 lines
- Command implementation logic: 200-500 lines
- Subcommand routing: 50-100 lines

#### Authentication
- Basic auth client: 150-250 lines
- Token-based auth: 300-500 lines
- OAuth flow: 500-800 lines
- Credential validation: 100-150 lines

#### API/Network Integration
- HTTP client setup: 100-200 lines
- API endpoint implementation: 200-400 lines per endpoint
- Request/response models: 50-100 lines per model
- Connection handling: 150-300 lines

#### Data Operations
- CRUD operations: 200-400 lines per resource
- Data validation: 100-200 lines
- Data transformation: 150-300 lines

#### Testing
- Unit test suite: 300-500 lines
- Integration tests: 400-600 lines
- Mock/fixture setup: 200-300 lines

**Output Format**:
```markdown
### Feature Sizing Analysis

| Feature | Components | Est. Lines | Notes |
|---------|-----------|-----------|-------|
| CLI push command | skeleton + flags + routing | 250 | cobra-style command |
| Auth client | basic auth + validation | 300 | username/password only |
| OCI push operations | go-containerregistry + upload | 600 | core push logic |
| Progress tracking | progress bars + status | 150 | user feedback |
| TLS handling | cert skip option | 100 | insecure flag |
| Error handling | validation + reporting | 200 | comprehensive errors |
| Unit tests | test suite + mocks | 400 | 70%+ coverage |
| Integration tests | e2e test scenarios | 400 | real registry tests |
```

### 4. Split Features into Efforts (≤800 lines each)

**HARD RULE**: No effort can exceed 800 lines of implementation code.

**Process**:
1. For each feature >700 lines: Split into multiple efforts
2. Find natural split boundaries:
   - Separate setup from logic
   - Separate read from write operations
   - Separate basic from advanced features
   - Separate implementation from testing

**Example Split**:
```
OCI Push Operations (600 lines) → Keep as single effort ✅

Auth Client + Token Refresh (850 lines) → Split:
  - Effort A: Basic auth client (300 lines)
  - Effort B: Token refresh + caching (350 lines)

Testing Suite (800 lines) → Split:
  - Effort A: Unit tests (400 lines)
  - Effort B: Integration tests (400 lines)
```

### 5. Organize Efforts into Waves

**Wave Rules**:
- Each wave contains 3-6 efforts
- Efforts in a wave should be thematically related
- Wave should have a clear deliverable goal
- Total wave size: 1000-3000 lines typically

**Wave Organization Patterns**:

#### Foundation Wave
- Project setup
- Core data models
- Basic infrastructure
- Initial test framework

#### Feature Implementation Wave
- Related feature implementations
- All share common dependencies
- Can be worked on by multiple engineers in parallel

#### Testing/Quality Wave
- Test suites
- Integration tests
- Documentation
- Performance optimization

**Example Wave Structure**:
```markdown
## Wave 1.1: CLI Foundation
**Goal**: Establish command structure and basic I/O
**Total Estimated Lines**: 800

- Effort 1.1.1: Push command skeleton (150 lines)
- Effort 1.1.2: Flag parsing and validation (200 lines)
- Effort 1.1.3: Command routing and help text (100 lines)
- Effort 1.1.4: Basic error handling framework (150 lines)
- Effort 1.1.5: Unit tests for CLI layer (200 lines)

## Wave 1.2: Authentication
**Goal**: Implement registry authentication
**Total Estimated Lines**: 650

- Effort 1.2.1: Basic auth client (300 lines)
- Effort 1.2.2: Credential validation (150 lines)
- Effort 1.2.3: Auth tests (200 lines)
```

### 6. Organize Waves into Phases

**Phase Rules**:
- Each phase contains 2-4 waves
- Phase should deliver significant value increment
- Phase 1 should be MVP/foundation
- Later phases build on earlier phases

**Typical Phase Structure**:

#### Phase 1: Foundation / MVP (30-40% of total work)
- Basic functionality working
- Core features implemented
- Can demonstrate value
- Tests passing

#### Phase 2: Feature Completion (40-50% of total work)
- All planned features
- Edge cases handled
- Performance optimized
- Integration complete

#### Phase 3: Production Readiness (20-30% of total work)
- Documentation complete
- Security hardened
- Performance tuned
- Deployment ready

**Example Phase Structure**:
```markdown
# Phase 1: MVP Push Command (3 waves, 2200 lines)

Goal: Working push command that can upload images to Gitea registry

## Wave 1.1: CLI Foundation (800 lines)
[efforts...]

## Wave 1.2: Authentication (650 lines)
[efforts...]

## Wave 1.3: Core Push Implementation (750 lines)
[efforts...]

# Phase 2: Production Features (2 waves, 1600 lines)

Goal: Production-ready with error handling, progress, and testing

## Wave 2.1: Error Handling & Progress (600 lines)
[efforts...]

## Wave 2.2: Comprehensive Testing (1000 lines)
[efforts...]
```

### 7. Map Dependencies Between Efforts

**Dependency Types**:
- **Sequential**: Effort B requires Effort A to be complete
- **Parallel**: Efforts can be worked on simultaneously
- **Optional**: Effort B is better with Effort A but not required

**Notation**:
```markdown
### Effort 1.2.1: Auth Client
**Dependencies**:
- Effort 1.1.1 (CLI skeleton) - REQUIRED
- Effort 1.1.2 (flag parsing) - REQUIRED

**Enables**:
- Effort 1.3.1 (OCI push) - This effort needs auth
```

**Dependency Mapping Rules**:
- CLI skeleton must come first (all other efforts depend on it)
- Auth must come before any authenticated operations
- Data models must come before operations using them
- Implementation must come before tests for that implementation
- Integration tests come after all implementation

### 8. Create Decomposition Output

Generate `init-decomposition.json`:

```json
{
  "project_name": "idpbuilder-push",
  "total_estimated_loc": 3800,
  "phases": [
    {
      "phase_number": 1,
      "phase_name": "MVP Push Command",
      "phase_goal": "Working push command for OCI images to Gitea",
      "estimated_lines": 2200,
      "waves": [
        {
          "wave_number": 1,
          "wave_name": "CLI Foundation",
          "wave_goal": "Command structure and I/O established",
          "estimated_lines": 800,
          "efforts": [
            {
              "effort_id": "E1.1.1",
              "name": "Push command skeleton",
              "description": "Add push subcommand to idpbuilder CLI with cobra, basic flag definitions, help text",
              "estimated_loc": 150,
              "sizing_justification": "Cobra command definition (50 lines) + flag setup (50 lines) + help text (50 lines)",
              "dependencies": [],
              "parallel_safe": false
            },
            {
              "effort_id": "E1.1.2",
              "name": "Flag parsing and validation",
              "description": "Parse -username, -password, -registry-url, -insecure flags, validate required flags present",
              "estimated_loc": 200,
              "sizing_justification": "Flag binding (50 lines) + validation functions (100 lines) + error messages (50 lines)",
              "dependencies": ["E1.1.1"],
              "parallel_safe": false
            }
          ]
        }
      ]
    }
  ],
  "dependency_graph": {
    "E1.1.1": [],
    "E1.1.2": ["E1.1.1"],
    "E1.2.1": ["E1.1.1", "E1.1.2"],
    "E1.3.1": ["E1.2.1"]
  }
}
```

### 9. Create Sizing Justification Document

Generate `init-sizing-justifications.md`:

```markdown
# Effort Sizing Justifications

## Methodology
All estimates based on:
- Historical data from similar Go projects
- go-containerregistry library complexity
- Cobra CLI framework patterns
- Standard Go testing practices

## Effort-by-Effort Breakdown

### E1.1.1: Push Command Skeleton (150 lines)
**Components**:
- Cobra command struct: 30 lines
- RunE function skeleton: 40 lines
- Flag definitions: 50 lines
- Help text and examples: 30 lines

**Similar Projects**:
- kubectl commands average 120-180 lines
- docker CLI commands average 150-200 lines

**Confidence**: HIGH (cobra pattern is well-established)

### E1.1.2: Flag Parsing and Validation (200 lines)
**Components**:
- Flag binding functions: 50 lines
- Validation logic: 100 lines (validate URLs, check required, format checks)
- Error message generation: 50 lines

**Similar Projects**:
- Similar validation in kubectl: 180-220 lines
- Similar validation in docker: 200-250 lines

**Confidence**: HIGH (validation is straightforward)

[... continue for all efforts ...]
```

## Validation Checklist

Before transitioning to INIT_SYNTHESIZE_PLAN:

- [ ] All major features identified from PRD
- [ ] Every feature has LOC estimate with justification
- [ ] All efforts are ≤800 lines (hard limit)
- [ ] Efforts organized into 3-6 per wave
- [ ] Waves organized into 2-4 per phase
- [ ] Dependencies mapped for all efforts
- [ ] Phase 1 delivers MVP (can demonstrate value)
- [ ] Total project scope is realistic
- [ ] init-decomposition.json created
- [ ] init-sizing-justifications.md created

## Common Sizing Mistakes to Avoid

**Mistake 1**: Forgetting testing
- Always include test efforts (usually 30-40% of implementation LOC)

**Mistake 2**: Underestimating integration
- Integration code often larger than unit implementations

**Mistake 3**: Ignoring error handling
- Comprehensive error handling adds 20-30% to estimates

**Mistake 4**: Missing infrastructure
- CLI setup, logging, config - don't forget the "boring" stuff

**Mistake 5**: Creating "catch-all" efforts
- "Implement everything" efforts are red flags - break them down

## Exit Criteria
- Decomposition complete with sizing
- All efforts validated ≤800 lines
- Dependencies mapped
- Ready for plan synthesis
- User has NOT had to manually structure anything

## Transition
**MANDATORY**: → INIT_SYNTHESIZE_PLAN (architect continues)

## Time Guidance
- Feature identification: 5-10 minutes
- LOC estimation: 10-15 minutes
- Effort organization: 10-15 minutes
- Dependency mapping: 5-10 minutes
- Total time: 30-50 minutes

## Error Handling
- Unclear scope → Ask user for clarification
- Too large → Recommend phased approach
- Too small → Suggest combining features
- Complex dependencies → Create dependency diagram

## Reference Files
- **READ**: `/home/vscode/software-factory-template/expertise/effort-sizing-guidelines.md`
- **READ**: `/home/vscode/software-factory-template/templates/IMPLEMENTATION-PLAN-TEMPLATE.md`
- **READ**: Examples in `/home/vscode/software-factory-template/templates/IMPLEMENTATION-PLAN-EXAMPLE.md`

## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

### YOU MUST OUTPUT THE CONTINUATION FLAG AS YOUR LAST ACTION

**EVERY STATE COMPLETION MUST END WITH EXACTLY ONE OF:**
```bash
# If state completed successfully and factory should continue:
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"

# If error/block/manual review needed:
echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
```

### CRITICAL REQUIREMENTS:
1. **ABSOLUTE LAST OUTPUT**: This MUST be the very last line of output before state completion
2. **NO TEXT AFTER**: No explanations, summaries, or any text after the flag
3. **EXACTLY AS SHOWN**: Use exact format - no variations
4. **ALWAYS REQUIRED**: Every single state must output this flag
5. **GREPPABLE**: Must be on its own line for automation parsing

### WHEN TO USE TRUE:
- ✅ Decomposition completed successfully
- ✅ All efforts ≤800 lines validated
- ✅ Dependencies mapped
- ✅ Ready for synthesis state
- ✅ Files created (init-decomposition.json, init-sizing-justifications.md)

### WHEN TO USE FALSE:
- ❌ Cannot determine scope from PRD
- ❌ User input needed for critical decisions
- ❌ Requirements insufficient for sizing
- ❌ Technical constraints unclear

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**
