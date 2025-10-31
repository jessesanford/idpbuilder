# SW Engineer - INIT_CUSTOMIZE_AGENTS State Rules

## Purpose
Add minimal language and framework expertise to agent configuration files.

## Entry Criteria
- IMPLEMENTATION-PLAN.md created
- Technology stack determined
- Agent files need customization

## Required Actions

### 1. Read Technology Requirements
From init-state-${PROJECT_PREFIX}.json:
- Primary language
- Key frameworks
- Testing tools
- Build system

### 2. Customize sw-engineer.md

**Location**: `$CLAUDE_PROJECT_DIR/.claude/agents/sw-engineer.md`

**Add After Line**: "## Core Responsibilities"

**Add Section** (MAXIMUM 10 lines):
```markdown
## Project-Specific Expertise

### Language: [LANGUAGE]
- Syntax and idioms for [LANGUAGE]
- Standard library usage
- Best practices and conventions

### Frameworks: [FRAMEWORK LIST]
- [Framework 1]: [1-line expertise]
- [Framework 2]: [1-line expertise]

### Build & Test: [BUILD_SYSTEM] / [TEST_FRAMEWORK]
- Build with [BUILD_SYSTEM] commands
- Test using [TEST_FRAMEWORK] patterns
```

### 3. Customize code-reviewer.md

**Location**: `$CLAUDE_PROJECT_DIR/.claude/agents/code-reviewer.md`

**Add After Line**: "## Review Criteria"

**Add Section** (MAXIMUM 10 lines):
```markdown
## Project-Specific Review Points

### [LANGUAGE] Standards
- Verify [LANGUAGE]-specific conventions
- Check for [LANGUAGE] anti-patterns
- Ensure proper error handling for [LANGUAGE]

### Framework Compliance
- [Framework]: Ensure proper usage patterns
- Check for framework best practices

### Testing Standards
- [TEST_FRAMEWORK] test coverage
- Test naming conventions for [LANGUAGE]
```

### 4. Optionally Customize architect.md

**Only if architecture pattern requires**:
```markdown
## Architecture Expertise

### Pattern: [PATTERN]
- Ensure [PATTERN] principles followed
- Validate component boundaries
```

### 5. Validation Rules

**CRITICAL CONSTRAINTS**:
- ❌ DO NOT modify any existing rules
- ❌ DO NOT remove any sections
- ❌ DO NOT exceed 10 lines per agent
- ✅ ONLY add expertise sections
- ✅ Keep additions minimal and focused
- ✅ Preserve all rule references

### 6. Update State File
Record customizations:
```json
"agent_customization": {
  "sw_engineer_customized": true,
  "code_reviewer_customized": true,
  "architect_customized": false,
  "language_added": "[LANGUAGE]",
  "frameworks_added": ["..."],
  "timestamp": "[ISO_TIME]"
}
```

## Examples by Language

### Go Customization
```markdown
## Project-Specific Expertise

### Language: Go
- Go idioms and effective patterns
- Proper error handling with error wrapping
- Context usage and cancellation

### Frameworks: cobra, gin
- cobra: CLI command structure and flags
- gin: HTTP routing and middleware

### Build & Test: make / go test
- Makefile targets for build/test/lint
- Table-driven tests with subtests
```

### Python Customization
```markdown
## Project-Specific Expertise

### Language: Python
- Python 3.10+ features and type hints
- Proper use of dataclasses and protocols
- Async/await patterns where applicable

### Frameworks: FastAPI, SQLAlchemy
- FastAPI: Dependency injection and Pydantic models
- SQLAlchemy: ORM patterns and session management

### Build & Test: poetry / pytest
- Poetry dependency management
- Pytest fixtures and parameterization
```

### TypeScript Customization
```markdown
## Project-Specific Expertise

### Language: TypeScript
- Strict TypeScript with no-any rule
- Proper type inference and generics
- Union types and discriminated unions

### Frameworks: React, Next.js
- React: Hooks and functional components
- Next.js: SSR/SSG patterns and API routes

### Build & Test: npm / jest
- NPM scripts for development workflow
- Jest with React Testing Library
```

## Exit Criteria
- Agent files have expertise sections
- Customizations are minimal (<10 lines)
- No existing rules modified
- Language expertise documented
- State file updated

## Transition
**MANDATORY**: → INIT_VALIDATION

## Common Mistakes to Avoid
- ❌ Adding too much detail
- ❌ Modifying existing sections
- ❌ Adding project-specific rules
- ❌ Exceeding line limits
- ✅ Keep it minimal and focused

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
3. **EXACTLY AS SHOWN**: Use exact format - no variations like CONTINUE-ORCHESTRATING
4. **ALWAYS REQUIRED**: Every single state must output this flag
5. **GREPPABLE**: Must be on its own line for automation parsing

### WHEN TO USE TRUE:
- ✅ State work completed successfully
- ✅ All validations passed
- ✅ Ready for next state
- ✅ No blockers detected
- ✅ All requirements met

### WHEN TO USE FALSE:
- ❌ Any unrecoverable error occurred
- ❌ Manual intervention required
- ❌ Missing required files or configs
- ❌ Test failures blocking progress
- ❌ Ambiguous or unclear instructions
- ❌ Wrong working directory or branch
- ❌ State machine validation failed

### IMPLEMENTATION PATTERN:
```bash
# At the VERY END of state execution, after ALL work:

# Determine success/failure
if [[ "$STATE_COMPLETED_PROJECT_DONEFULLY" == "true" ]]; then
    echo "✅ State work completed successfully"
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
else
    echo "❌ State encountered issues requiring intervention"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
fi

# DO NOT OUTPUT ANYTHING AFTER THE FLAG!
```

### GRADING IMPACT:
- **Missing flag**: -100% AUTOMATIC FAILURE
- **Text after flag**: -50% penalty
- **Wrong format**: -100% AUTOMATIC FAILURE
- **Multiple flags**: -50% penalty

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

