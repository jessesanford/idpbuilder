# Code-reviewer - EFFORT_PLAN_CREATION State Rules

## State Context
This is the EFFORT_PLAN_CREATION state for the code-reviewer.

## Acknowledgment Required
Thank you for reading the rules file for the EFFORT_PLAN_CREATION state.

**IMPORTANT**: Please report that you have successfully read the EFFORT_PLAN_CREATION rules file.

Say: "✅ Successfully read EFFORT_PLAN_CREATION rules for code-reviewer"

## 🔴🔴🔴 PARAMOUNT: Repository Separation Understanding (R251 & R309) 🔴🔴🔴

### R251: Universal Repository Separation Law
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R251-REPOSITORY-SEPARATION-LAW.md`
**Criticality**: PARAMOUNT - Automatic -100% failure for violation
**KEY UNDERSTANDING**: 
- Plans are created IN effort directories (TARGET repo clones)
- Plans are stored in `.software-factory/` subdirectory within the effort
- Implementation will happen IN TARGET repo clones
- NEVER in Software Factory repo

### R309: Never Create Efforts in SF Repo
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R309-never-create-efforts-in-sf-repo.md`
**Criticality**: PARAMOUNT - Automatic -100% failure for violation
**KEY UNDERSTANDING**:
- Your plan will be used by SW-Engineer IN TARGET repo clone
- Effort directory is under /efforts/ (not SF root)
- Plans are stored in `.software-factory/phaseX/waveY/effort-name/` subdirectory
- Plan must reference TARGET repo structure, not SF structure

**VERIFY YOUR UNDERSTANDING:**
```bash
echo "🔴 R251/R309: Understanding repository context..."
echo "I understand:"
echo "  ✅ I'm creating plan for TARGET repo implementation"
echo "  ✅ SW-Engineer will work in /efforts/ clone"
echo "  ✅ NOT in Software Factory planning repo"
echo "  ✅ Plan references pkg/, cmd/, etc. (TARGET structure)"
echo "  ✅ NOT .claude/, rule-library/ (SF structure)"
```

## 🔴🔴🔴 MANDATORY MEASUREMENT RULE - R304 🔴🔴🔴

### R304: Mandatory Line Counter Tool Enforcement
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R304-mandatory-line-counter-enforcement.md`
**Criticality**: BLOCKING - Manual counting = AUTOMATIC -100% FAILURE

**ABSOLUTE REQUIREMENTS:**
- ✅ MUST use `$CLAUDE_PROJECT_DIR/tools/line-counter.sh` for ALL line counting
- ❌ NEVER use `wc -l` or any manual counting method
- ❌ NEVER count lines any other way - this is a -100% automatic failure
- ✅ MUST specify both -b (base branch) and -c (current branch) parameters
- ✅ Base branch MUST be phase integration branch (NOT "main")

**Failure to use the line counter tool = AUTOMATIC -100% GRADE**

## State-Specific Rules

### 🔴🔴🔴 EXPLICIT SCOPE CONTROL (R311 - SUPREME LAW) 🔴🔴🔴

**CRITICAL**: Effort plans MUST include explicit scope boundaries to prevent over-engineering!

#### Mandatory Scope Definition Requirements:

1. **EXACT Function/Method Counts**
   - List EXACTLY how many functions to implement
   - Name each function explicitly where possible
   - Include line estimates for each function
   - State "NO MORE" after the list

2. **DO NOT IMPLEMENT Section (CRITICAL)**
   - MUST include explicit list of what NOT to build
   - Common exclusions: validation, caching, logging, extra CRUD ops
   - Be specific: "DO NOT add Update/Delete" not just "minimal scope"
   - This prevents 3-5X scope creep

3. **Realistic Size Calculations**
   ```
   Functions: 3 × 40 lines = 120 lines
   Types: 2 × 30 lines = 60 lines
   Tests: 5 × 30 lines = 150 lines
   TOTAL: 330 lines (well under 800)
   ```

4. **Scope Enforcement Checkpoints**
   - Before starting: Acknowledge boundaries
   - During implementation: Count functions/types
   - Before commit: Verify no extras added

#### Example Scope Definition:
```markdown
## EXPLICIT SCOPE (R311 MANDATORY)

### IMPLEMENT EXACTLY:
- Function: CreateUser(user User) error (~40 lines)
- Function: GetUser(id string) (*User, error) (~35 lines)
- Type: User struct with 5 fields, NO methods (~20 lines)
- Tests: 2 basic tests only (~60 lines)
TOTAL: ~155 lines

### DO NOT IMPLEMENT:
- ❌ UpdateUser (future effort)
- ❌ DeleteUser (future effort)
- ❌ ListUsers (future effort)
- ❌ User validation methods
- ❌ Caching layer
- ❌ Comprehensive error handling
- ❌ Edge case tests
```

**FAILURE TO INCLUDE EXPLICIT SCOPE = -75% PENALTY**

### 🔴🔴🔴 ATOMIC PR EFFORT REQUIREMENTS (R220 - SUPREME LAW) 🔴🔴🔴

When creating effort implementation plans, you MUST ensure the effort produces exactly ONE atomic PR:

1. **One Effort = One PR (ABSOLUTE)**
   - This effort must result in EXACTLY one PR to main
   - PR must merge independently of all other efforts
   - PR must not break the build when merged alone
   - NO EXCEPTIONS TO THIS RULE

2. **Feature Flags for This Effort**
   - Define specific flags for incomplete features
   - Document exact implementation location
   - Include flag initialization code
   - Plan tests with flag on/off
   - Specify cleanup conditions

3. **Stubs for Dependencies**
   - Identify what this effort depends on
   - Create stubs for missing dependencies
   - Ensure stubs match interface contracts
   - Document when stubs get replaced
   - Test with both stubs and real implementations

4. **Interface Implementation**
   - If defining interface: complete specification
   - If implementing interface: match contract exactly
   - Support both current and future use cases
   - Maintain backward compatibility
   - Document any assumptions

5. **PR Completeness Checklist**
   - All code for effort in ONE PR
   - All tests pass independently
   - Feature flags control activation
   - Documentation included
   - No dependencies on unmerged PRs

### Effort Plan MUST Include

```yaml
effort_atomic_pr_design:
  pr_summary: "Single PR implementing [specific feature]"
  can_merge_to_main_alone: true  # MUST be true
  
  feature_flags_needed:
    - flag: "EFFORT_X_FEATURE_Y"
      purpose: "Hide incomplete feature Y"
      default: false
      location: "config/features.yaml"
      activation: "When all components ready"
  
  stubs_required:
    - stub: "MockServiceZ"
      replaces: "ServiceZ (from effort_5)"
      interface: "IServiceZ"
      behavior: "Returns default success response"
  
  interfaces_to_implement:
    - interface: "IDataProcessor"
      methods: ["process", "validate"]
      implementation: "Complete in this PR"
  
  pr_verification:
    tests_pass_alone: true
    build_remains_working: true
    flags_tested_both_ways: true
    no_external_dependencies: true
    backward_compatible: true
  
  example_pr_structure:
    files_added:
      - "src/feature_x.go"
      - "src/feature_x_test.go"
      - "config/features.yaml"
      - "stubs/mock_service_z.go"
    tests_included:
      - "Unit tests with flag off"
      - "Unit tests with flag on"
      - "Integration test with stubs"
    documentation:
      - "README update"
      - "API documentation"
```

### CRITICAL VALIDATION

Before completing effort plan, verify:
- ✅ This effort = ONE atomic PR to main
- ✅ PR can merge without any other effort
- ✅ Build stays green when PR merges
- ✅ Feature flags hide incomplete work
- ✅ All dependencies stubbed/mocked
- ✅ Tests pass in complete isolation

**FAILURE TO ENSURE ATOMIC PR = -100% IMMEDIATE FAILURE**

## 🔴🔴🔴 CRITICAL: Plan Storage Location (NEW REQUIREMENT) 🔴🔴🔴

### Plans MUST be Stored in .software-factory/ Subdirectory

**MANDATORY STORAGE PATTERN:**
```bash
# Determine phase, wave, and effort name from context
PHASE="1"  # Example
WAVE="1"   # Example
EFFORT_NAME="go-containerregistry-image-builder"  # Example

# Create the .software-factory directory structure
PLAN_DIR=".software-factory/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}"
mkdir -p "$PLAN_DIR"

# Create plan with timestamp
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
PLAN_FILE="$PLAN_DIR/IMPLEMENTATION-PLAN-${TIMESTAMP}.md"

echo "📁 Creating plan directory: $PLAN_DIR"
echo "📝 Plan will be saved as: $PLAN_FILE"
```

**EXAMPLE PATHS:**
- Effort plan: `.software-factory/phase2/wave1/go-containerregistry-image-builder/IMPLEMENTATION-PLAN-20250104-143000.md`
- Split plan: `.software-factory/phase2/wave1/go-containerregistry-image-builder-split-001/SPLIT-PLAN-20250104-145500.md`

**WHY THIS STRUCTURE?**
- Keeps plans organized within working copies
- Prevents clutter in root of effort directory
- Maintains clear hierarchy of phase/wave/effort
- Timestamps prevent collision per R301
- SW Engineers know exactly where to look

### Creating the Plan File
```bash
cat > "$PLAN_FILE" << 'EOF'
# Implementation Plan for ${EFFORT_NAME}
Created: $(date -Iseconds)
Location: ${PLAN_DIR}
Phase: ${PHASE}
Wave: ${WAVE}

## Effort Metadata
[Plan content here...]
EOF

# Commit the plan
git add "$PLAN_FILE"
git commit -m "feat: add implementation plan for ${EFFORT_NAME}

Plan location: $PLAN_FILE
Phase $PHASE, Wave $WAVE"
git push

echo "✅ Plan created and committed: $PLAN_FILE"
```

### SW Engineer Will Look For Plans Here
```bash
# SW Engineer's plan discovery logic:
find_latest_plan() {
    # First check new location
    if [ -d ".software-factory" ]; then
        LATEST_PLAN=$(find .software-factory -name "IMPLEMENTATION-PLAN-*.md" -type f | sort -r | head -1)
        if [ -n "$LATEST_PLAN" ]; then
            echo "✅ Found plan in new location: $LATEST_PLAN"
            return 0
        fi
    fi
    
    # Fallback to old location for backward compatibility
    if [ -f "IMPLEMENTATION-PLAN.md" ]; then
        echo "⚠️ Found legacy plan: IMPLEMENTATION-PLAN.md"
        echo "Consider migrating to .software-factory/ structure"
        LATEST_PLAN="IMPLEMENTATION-PLAN.md"
        return 0
    fi
    
    echo "❌ No implementation plan found!"
    return 1
}
```

## General Responsibilities
Follow all general code-reviewer rules and the Software Factory state machine.

## Next Steps
Proceed with the standard workflow for the EFFORT_PLAN_CREATION state as defined in the state machine.
