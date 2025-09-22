# 🚨🚨🚨 BLOCKING RULE R054: Implementation Plan Creation

## Rule Statement
Code Reviewers MUST create comprehensive IMPLEMENTATION-PLAN.md files for each effort before SW Engineers begin work. Plans MUST include technical architecture, implementation sequence, size management strategy, and validation requirements.

## Plan Requirements

### 1. Mandatory Plan Components
Every IMPLEMENTATION-PLAN.md MUST contain:
- ✅ **Effort Identification**: ID, name, phase, wave
- ✅ **Technical Architecture**: Design decisions and patterns
- ✅ **File Structure**: List of files to create/modify
- ✅ **Implementation Sequence**: Step-by-step order
- ✅ **Dependencies**: What this effort needs from others
- ✅ **Size Management**: Strategy to stay under 800 lines
- ✅ **Testing Strategy**: Unit, integration, E2E requirements
- ✅ **Validation Checkpoints**: When to verify correctness
- ✅ **Integration Points**: How this connects to other efforts

### 2. Plan Template Structure
```markdown
# IMPLEMENTATION PLAN: {effort-id} - {effort-name}

## Effort Metadata
- **Effort ID**: {effort-id}
- **Phase**: {X}
- **Wave**: {Y}
- **Created By**: code-reviewer
- **Created At**: {timestamp}
- **Estimated Lines**: {estimate}

## Technical Architecture

### Design Overview
[High-level design decisions and patterns to use]

### API Contracts
[Interfaces this effort exposes or consumes]

### Data Models
[Structures, schemas, or types to implement]

## File Structure
```
{effort-directory}/
├── src/
│   ├── {component1}.{ext}  # [New/Modified] - Description
│   ├── {component2}.{ext}  # [New/Modified] - Description
│   └── ...
├── tests/
│   ├── {test1}.{ext}       # Unit tests for component1
│   └── ...
└── docs/
    └── {documentation}.md   # API documentation
```

## Implementation Sequence

### Step 1: Foundation Setup
- [ ] Create directory structure
- [ ] Initialize configuration files
- [ ] Set up logging framework
- **Lines**: ~50

### Step 2: Core Implementation
- [ ] Implement {main-component}
- [ ] Add error handling
- [ ] Create utility functions
- **Lines**: ~300

### Step 3: API Layer
- [ ] Define API endpoints
- [ ] Implement request handlers
- [ ] Add validation middleware
- **Lines**: ~200

### Step 4: Testing
- [ ] Write unit tests (90% coverage)
- [ ] Add integration tests
- [ ] Create E2E test scenarios
- **Lines**: ~200

### Step 5: Documentation
- [ ] API documentation
- [ ] Code comments
- [ ] README updates
- **Lines**: ~50

**Total Estimated Lines**: ~800

## Dependencies

### Requires From Other Efforts
- **effort-1**: API client library
- **effort-2**: Shared data models
- **None**: This is independent

### Provides To Other Efforts
- **effort-4**: Authentication service API
- **effort-5**: User management endpoints

## Size Management Strategy

### Line Count Controls
1. **Monitor every 100 lines**: Run line counter
2. **Refactor at 600 lines**: Extract utilities
3. **Split trigger at 700 lines**: Stop and plan split
4. **Hard limit 800 lines**: Never exceed

### Potential Split Points
- After Step 2: Core can be separate effort
- After Step 3: API layer can be split
- Testing as separate effort if needed

## Testing Requirements

### Unit Tests (Required: 90% coverage)
```{language}
// Example test structure
describe('{ComponentName}', () => {
    test('should handle normal case', () => {
        // Test implementation
    });
    
    test('should handle error case', () => {
        // Error handling test
    });
});
```

### Integration Tests
- API endpoint testing
- Database interaction verification
- External service mocking

### E2E Tests
- Complete user workflow
- Cross-component interaction
- Performance validation

## Validation Checkpoints

### Checkpoint 1: After Foundation (Step 1)
- [ ] Directory structure correct
- [ ] Configuration loads properly
- [ ] Logging works

### Checkpoint 2: After Core (Step 2)
- [ ] Core logic functions correctly
- [ ] Error handling works
- [ ] Unit tests pass

### Checkpoint 3: After API (Step 3)
- [ ] API endpoints respond
- [ ] Validation works
- [ ] Integration tests pass

### Checkpoint 4: Final Validation
- [ ] All tests pass
- [ ] Line count under 800
- [ ] Documentation complete
- [ ] Code review ready

## Risk Assessment

### Technical Risks
1. **Risk**: API compatibility with effort-1
   **Mitigation**: Early integration testing
   
2. **Risk**: Size limit exceeded
   **Mitigation**: Pre-planned split points

### Schedule Risks
1. **Risk**: Dependency on effort-1 completion
   **Mitigation**: Mock interfaces for parallel work

## Integration Strategy

### How This Connects
1. **Upstream**: Consumes APIs from efforts 1,2
2. **Downstream**: Provides services to efforts 4,5
3. **Parallel**: No interaction with effort 3

### Integration Testing Plan
- Mock upstream dependencies initially
- Test with real services when available
- Validate downstream compatibility
```

## Code Reviewer Responsibilities

### Before Creating Plan
1. **READ** architecture documents
2. **ANALYZE** effort requirements
3. **IDENTIFY** dependencies with other efforts
4. **ESTIMATE** implementation size

### During Plan Creation
1. **DESIGN** technical architecture
2. **SEQUENCE** implementation steps
3. **PLAN** size management strategy
4. **DEFINE** test requirements
5. **IDENTIFY** split points

### After Plan Creation
1. **SAVE** in .software-factory subdirectory structure per R303 with MANDATORY timestamps per R383
   ```bash
   # 🔴🔴🔴 R383 MANDATORY: Use sf_metadata_path helper function
   source $CLAUDE_PROJECT_DIR/utilities/sf-metadata-path.sh

   # Determine phase, wave, and effort name from context
   PHASE="1"  # Example
   WAVE="2"   # Example
   EFFORT_NAME="buildah-builder-interface"  # Example

   # Create plan with MANDATORY timestamp (R383)
   PLAN_FILE=$(sf_metadata_path "$PHASE" "$WAVE" "$EFFORT_NAME" "IMPLEMENTATION-PLAN" "md")

   echo "📁 Creating plan at: $PLAN_FILE"
   # Write content to plan file
   cat > "$PLAN_FILE" << 'EOF'
   # Implementation Plan content here...
   EOF
   ```
2. **VALIDATE** plan completeness
3. **COMMIT** to git repository
4. **REPORT** completion to orchestrator with exact path
5. **UPDATE** orchestrator-state.json with plan location (if possible)

## Common Violations

### ❌ VIOLATION: Missing Size Strategy
```markdown
# WRONG - No size management plan
## Implementation Steps
1. Build everything
2. Test it
3. Ship it
```

### ❌ VIOLATION: No Dependencies Identified
```markdown
# WRONG - Dependencies section empty or missing
## Dependencies
TBD
```

### ❌ VIOLATION: Vague Implementation Sequence
```markdown
# WRONG - No clear steps
## Implementation
We'll figure it out as we go
```

### ✅ CORRECT: Complete Plan
- Clear implementation sequence with line estimates
- Identified dependencies and integration points
- Size management strategy with split points
- Comprehensive test requirements

## Enforcement

### Plan Validation
```bash
validate_implementation_plan() {
    local plan_file="$1"
    
    # Check required sections
    for section in "Effort Metadata" "Technical Architecture" "Implementation Sequence" \
                   "Dependencies" "Size Management" "Testing Requirements"; do
        if ! grep -q "## $section" "$plan_file"; then
            echo "❌ Missing required section: $section"
            return 1
        fi
    done
    
    # Check for line estimates
    if ! grep -q "Total Estimated Lines:" "$plan_file"; then
        echo "❌ Missing total line estimate"
        return 1
    fi
    
    echo "✅ Implementation plan validated"
}
```

## Reading Implementation Plans

When SW Engineers need to read an implementation plan:
```bash
# Find the most recent implementation plan
get_latest_implementation_plan() {
    # Check new .software-factory location first (per R303)
    if [ -d ".software-factory" ]; then
        LATEST_PLAN=$(find .software-factory -name "IMPLEMENTATION-PLAN-*.md" -type f | sort -r | head -1)
        if [ -n "$LATEST_PLAN" ]; then
            echo "✅ Found plan in .software-factory structure: $LATEST_PLAN"
            cat "$LATEST_PLAN"
            return 0
        fi
    fi
    
    # Fallback to root directory for backward compatibility
    LATEST_PLAN=$(ls -t IMPLEMENTATION-PLAN-*.md 2>/dev/null | head -n1)
    if [ -n "$LATEST_PLAN" ]; then
        echo "⚠️ Found plan in legacy location: $LATEST_PLAN"
        cat "$LATEST_PLAN"
        return 0
    fi
    
    # Final fallback to old format
    if [ -f "IMPLEMENTATION-PLAN.md" ]; then
        echo "⚠️ Using legacy plan format: IMPLEMENTATION-PLAN.md"
        cat "IMPLEMENTATION-PLAN.md"
        return 0
    fi
    
    echo "❌ ERROR: No implementation plan found!"
    echo "   Checked: .software-factory/phase*/wave*/*/IMPLEMENTATION-PLAN-*.md"
    echo "   Checked: ./IMPLEMENTATION-PLAN-*.md"
    echo "   Checked: ./IMPLEMENTATION-PLAN.md"
    exit 1
}
```

## Integration with Other Rules

- **R007**: Size limit compliance (800 lines)
- **R053**: Dependencies feed parallelization decisions
- **R197**: One agent implements one plan
- **R031**: Plan defines review criteria
- **R303**: Phase/Wave document location protocol
- **R343**: Metadata directory standardization
- **R383**: Software Factory metadata file organization (MANDATORY timestamps)

## Grading Impact

- **-50%**: No implementation plan created
- **-30%**: Plan missing critical sections
- **-20%**: No size management strategy
- **-15%**: Dependencies not identified
- **-10%**: Vague or incomplete sequence

---

**REMEMBER**: The implementation plan is the CONTRACT between Code Reviewer and SW Engineer. It MUST be complete, specific, and actionable.