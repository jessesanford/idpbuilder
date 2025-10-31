# Architect - INIT_LOAD_EXAMPLES State Rules

## Purpose
Load and understand example files and templates to prepare for requirements gathering.

## Entry Criteria
- Spawned by orchestrator in initialization mode
- init-state-${PROJECT_PREFIX}.json exists
- Project idea/description available

## Required Actions

### 1. Read Example Files
**MANDATORY READS**:
```bash
# Read implementation plan examples
READ: $CLAUDE_PROJECT_DIR/templates/IMPLEMENTATION-PLAN-EXAMPLE.md
READ: $CLAUDE_PROJECT_DIR/templates/IMPLEMENTATION-PLAN-TEMPLATE.md

# Read configuration examples
READ: $CLAUDE_PROJECT_DIR/setup-config-example.yaml
READ: $CLAUDE_PROJECT_DIR/target-repo-config.yaml

# Read existing project examples if available
READ: $CLAUDE_PROJECT_DIR/examples/idpbuilder-oci/IMPLEMENTATION-PLAN.md
```

### 2. Analyze Structure Requirements
Identify from examples:
- Required sections in IMPLEMENTATION-PLAN.md
- Mandatory fields in setup-config.yaml
- Repository configuration options
- Phase/wave/effort structure patterns
- Success criteria formats

### 3. Prepare Question Framework
Based on examples, prepare questions for:
- Target codebase details
- Technology stack specifics
- Architecture patterns
- Deployment requirements
- Quality constraints
- Testing strategies

### 4. Create Question Templates
Structure questions to gather:
- All fields needed for setup-config.yaml
- All fields needed for target-repo-config.yaml
- Information for each section of IMPLEMENTATION-PLAN.md
- Agent expertise requirements

### 5. Update State File
Record in init-state-${PROJECT_PREFIX}.json:
- examples_loaded: true
- question_categories_prepared: [list]
- ready_for_requirements: true

## Exit Criteria
- All example files read and understood
- Question framework prepared
- Know all required fields to gather
- Ready to begin interactive Q&A

## Transition
**MANDATORY**: → INIT_REQUIREMENTS_GATHERING

## Best Practices
- Use examples as templates, not rigid structures
- Adapt questions to project type
- Prepare follow-up questions for each category
- Consider both new and existing codebases

## Time Limit
Maximum 5 minutes in this state

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

