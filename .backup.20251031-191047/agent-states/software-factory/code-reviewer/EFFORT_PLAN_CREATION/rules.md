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

### 🔴🔴🔴 PRIMARY DIRECTIVE: DEMO PLANNING REQUIREMENTS (R330 - SUPREME LAW) 🔴🔴🔴

### R330: Demo Planning Requirements (SUPREME LAW - BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R330-demo-planning-requirements.md`
**Criticality**: BLOCKING - Automatic -100% failure for plans without demo requirements
**Penalty**: -25% to -50% for incomplete demo planning

**EVERY effort plan MUST include explicit demo requirements!**

#### MANDATORY DEMO SECTION IN EVERY PLAN:

```markdown
## 🎬 Demo Requirements (R330 MANDATORY)

### Demo Objectives (3-5 specific, verifiable objectives)
- [ ] Demonstrate [specific functionality] works with [specific input]
- [ ] Show proper error handling for [specific error case]
- [ ] Verify integration with [specific upstream service]
- [ ] Prove performance meets [specific requirement]
- [ ] Display proper [logging/monitoring/etc.]

**Success Criteria**: All objectives checked = demo passes

### Demo Scenarios (IMPLEMENT EXACTLY THESE - 2-4 scenarios)

#### Scenario 1: [Scenario Name]
- **Setup**: [Prerequisites and initial state]
- **Input**: [Exact input data/commands]
- **Action**: [Exact command or API call]
- **Expected Output**:
  ```
  [Exact expected response]
  ```
- **Verification**: [How to verify success]
- **Script Lines**: ~X lines

#### Scenario 2: [Error Handling Scenario]
- **Setup**: [Prerequisites]
- **Input**: [Invalid/error input]
- **Action**: [Command that triggers error]
- **Expected Output**:
  ```
  [Expected error response]
  ```
- **Verification**: [How to verify error handled correctly]
- **Script Lines**: ~X lines

**TOTAL SCENARIO LINES**: ~XX lines

### Demo Size Planning

#### Demo Artifacts (Excluded from line count per R007)
```
demo-features.sh:     XX lines  # Executable script
DEMO.md:             XX lines  # Documentation
test-data/:          XX lines  # Sample data files
integration-hook.sh: XX lines  # For wave integration
────────────────────────────────
TOTAL DEMO FILES:   XXX lines (NOT counted toward 800)
```

#### Effort Size Summary
```
Implementation:     XXX lines  # ← ONLY this counts toward 800
────────────────────────────────
Tests:             XXX lines  # Excluded per R007
Demos:             XXX lines  # Excluded per R007
────────────────────────────────
Implementation:    XXX/800 ✅ (within limit)
```

**NOTE**: While demos don't count toward the line limit, they MUST still be planned and implemented as specified!

### Demo Deliverables

Required Files:
- [ ] `demo-features.sh` - Main demo script (executable)
- [ ] `DEMO.md` - Demo documentation per template
- [ ] `test-data/valid.json` - Valid input examples (if applicable)
- [ ] `test-data/invalid.json` - Invalid input examples (if applicable)
- [ ] `.demo-config` - Demo environment settings (if applicable)

Integration Hooks:
- [ ] Export DEMO_READY=true when complete
- [ ] Provide integration point for wave demo
- [ ] Include cleanup function
```

#### VALIDATION BEFORE FINALIZING PLAN:

```bash
# MANDATORY: Verify demo requirements before completing plan
verify_demo_requirements() {
    local plan_file="$PLAN_FILE"

    # Check for demo section
    if ! grep -q "Demo Requirements" "$plan_file"; then
        echo "❌ R330 VIOLATION: Missing demo requirements section!"
        echo "🔴 CANNOT finalize plan without demo requirements"
        exit 330
    fi

    # Check for scenarios
    if ! grep -q "Demo Scenarios" "$plan_file"; then
        echo "❌ R330 VIOLATION: Missing demo scenarios!"
        echo "🔴 Plans MUST include 2-4 specific demo scenarios"
        exit 330
    fi

    # Check for size planning
    if ! grep -q "Demo Size Planning" "$plan_file"; then
        echo "❌ R330 VIOLATION: Missing demo size calculation!"
        echo "🔴 Demo artifacts must be sized (even if not counted)"
        exit 330
    fi

    # Check for deliverables
    if ! grep -q "Demo Deliverables" "$plan_file"; then
        echo "❌ R330 VIOLATION: Missing demo deliverables section!"
        echo "🔴 Must specify exact demo files to create"
        exit 330
    fi

    echo "✅ R330 COMPLIANT: Demo requirements complete"
}

# Run validation before plan finalization
verify_demo_requirements || exit 330
```

#### CRITICAL REMINDERS:

- **Demo objectives must be verifiable** - Not subjective
- **Demo scenarios must be complete** - Setup → Action → Verify
- **Demo size must be calculated** - Even though not counted toward 800
- **Demo deliverables must be listed** - Exact filenames required
- **Integration hooks must be specified** - How this contributes to wave demo

**FAILURE TO INCLUDE DEMO REQUIREMENTS = -100% BLOCKING VIOLATION**

**Related Rules:**
- R291: Integration Demo Requirement (demos must work at integration)
- R311: Scope Control (demos are part of scope)
- R007: Size Limits (demos excluded from count)

---

## 🔴🔴🔴 MANDATORY PRE-PLANNING RESEARCH (R374 - SUPREME LAW) 🔴🔴🔴

**CRITICAL: BEFORE creating ANY plan, you MUST research existing code!**

#### R374: Pre-Planning Research Protocol
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R374-pre-planning-research-protocol.md`
**Criticality**: BLOCKING - No research = -50% penalty

**MANDATORY RESEARCH SEQUENCE:**
```bash
# 1. Search current wave for existing interfaces/implementations
echo "=== R374: Searching current wave for existing code ==="
for branch in $(git branch -r | grep "phase${PHASE}-wave${WAVE}"); do
    echo "Checking $branch..."
    git checkout $branch 2>/dev/null || continue

    # Find interfaces
    grep -r "type.*interface" --include="*.go" . 2>/dev/null | head -10

    # Find key functions
    for method in Push Pull Upload Download Store Create Delete Get List; do
        grep -r "func.*${method}(" --include="*.go" . 2>/dev/null | head -3
    done
done

# 2. Search previous waves
echo "=== R374: Searching previous waves ==="
for wave in $(seq 1 $((CURRENT_WAVE - 1))); do
    for branch in $(git branch -r | grep "phase${PHASE}-wave${wave}"); do
        git checkout $branch 2>/dev/null || continue
        grep -r "type.*interface" --include="*.go" . 2>/dev/null | head -5
    done
done

# 3. Document ALL findings
echo "=== R374: Documenting found interfaces and implementations ==="
```

#### R373: Mandatory Code Reuse and Interface Compliance
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R373-mandatory-code-reuse-and-interface-compliance.md`
**Criticality**: BLOCKING - Duplicate implementation = -100% IMMEDIATE FAILURE

**YOUR PLAN MUST INCLUDE:**
```markdown
## Pre-Planning Research Results (R374 MANDATORY)

### Existing Interfaces Found
| Interface | Location | Signature | Must Implement |
|-----------|----------|-----------|----------------|
| [List ALL interfaces found] | | | |

### Existing Implementations to Reuse
| Component | Location | Purpose | How to Use |
|-----------|----------|---------|------------|
| [List ALL reusable code] | | | |

### APIs Already Defined
| API | Method | Signature | Notes |
|-----|--------|-----------|-------|
| [List ALL existing APIs] | | | |

### FORBIDDEN DUPLICATIONS (R373)
- DO NOT create [list interfaces that already exist]
- DO NOT reimplement [list existing functionality]
- DO NOT create alternative [list existing method signatures]

### REQUIRED INTEGRATE_WAVE_EFFORTSS (R373)
- MUST implement [interface] from [location] with EXACT signature
- MUST reuse [component] from [location]
- MUST import and use [package] for [functionality]
```

**EXAMPLE VIOLATION TO PREVENT:**
```go
// WRONG - R373 VIOLATION - Creating competing interface
type MyRegistry interface {
    Push(image v1.Image, ref string) error  // Different signature!
}

// CORRECT - Implementing existing interface
import "existing/registry"
type MyClient struct{}
func (m *MyClient) Push(ctx context.Context, image string, content io.Reader) error {
    // Implements EXISTING interface exactly
}
```

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

### 🔴🔴🔴 PRODUCTION READY CODE REQUIREMENTS (R355 - SUPREME LAW) 🔴🔴🔴

**THIS IS SUPREME LAW #5 - ALL CODE MUST BE PRODUCTION READY FROM DAY ONE**

When creating implementation plans, you MUST explicitly forbid ALL non-production patterns:

#### ❌ EXPLICITLY FORBIDDEN IN IMPLEMENTATION:
```markdown
## 🚨🚨🚨 R355 PRODUCTION READINESS - ZERO TOLERANCE 🚨🚨🚨

This implementation MUST be production-ready from the first commit:
- ❌ NO STUBS or placeholder implementations
- ❌ NO MOCKS except in test directories
- ❌ NO hardcoded credentials or secrets
- ❌ NO static configuration values
- ❌ NO TODO/FIXME markers in code
- ❌ NO returning nil or empty for "later implementation"
- ❌ NO panic("not implemented") patterns
- ❌ NO fake or dummy data

VIOLATION = -100% AUTOMATIC FAILURE

## 🚨🚨🚨 R332 MANDATORY BUG FILING PROTOCOL (SUPREME LAW) 🚨🚨🚨

**File**: `$CLAUDE_PROJECT_DIR/rule-library/R332-mandatory-bug-filing-protocol.md`

**Integration with R355 and Effort Planning**:
1. Plans MUST include bug tracking system for effort
2. Plans MUST prevent "pre-existing" bug excuses
3. If deferring functionality to future effort, MUST provide exact evidence:
   - Exact plan file path for future effort
   - Exact line numbers showing planned implementation
   - Grep output as evidence
   - Phase/Wave comparison (current vs planned)
4. Vague "will be addressed later" = -100% FAILURE

**TODO Acceptance Criteria for Plans**:
- If plan includes TODO marker, MUST specify EXACT future effort ID
- Must document exact file and line where TODO will be placed
- Must document exact future effort that removes TODO
- Must be grep-verifiable in future effort plan

**See R332 for complete TODO acceptance criteria and bug filing protocol.**
```

#### ✅ REQUIRED PATTERNS IN YOUR PLAN:
Include configuration examples showing HOW to avoid violations:

```markdown
## Configuration Requirements (R355 Mandatory)

### WRONG - Will fail review:
```go
// ❌ VIOLATION - Hardcoded credential
password := "admin123"

// ❌ VIOLATION - Stub implementation
func ProcessPayment() error {
    // TODO: implement later
    return nil
}

// ❌ VIOLATION - Static configuration
apiEndpoint := "http://api.example.com"
```

### CORRECT - Production ready:
```go
// ✅ From environment variable
password := os.Getenv("DB_PASSWORD")
if password == "" {
    return errors.New("DB_PASSWORD not set")
}

// ✅ Full implementation required
func ProcessPayment(amount float64) error {
    client := payment.NewClient(config.PaymentKey)
    return client.Process(amount)
}

// ✅ Configurable endpoint
apiEndpoint := config.GetString("api.endpoint")
if apiEndpoint == "" {
    apiEndpoint = defaultEndpoint
}
```
```

### 🔴🔴🔴 ATOMIC PR EFFORT REQUIREMENTS (R220 - SUPREME LAW) 🔴🔴🔴

When creating effort implementation plans, you MUST ensure the effort produces exactly ONE atomic PR:

1. **One Effort = One PR (ABSOLUTE)**
   - This effort must result in EXACTLY one PR to main
   - PR must merge independently of all other efforts
   - PR must not break the build when merged alone
   - NO EXCEPTIONS TO THIS RULE

2. **Feature Flags for Incomplete Features ONLY**
   - Define specific flags ONLY for features not ready for production
   - Document exact implementation location
   - Include flag initialization code
   - Plan tests with flag on/off
   - Specify cleanup conditions
   - NOTE: Feature flags are NOT an excuse for stubs (R355)

3. **Interface Contracts Instead of Stubs**
   - If depending on unimplemented services, use dependency injection
   - Define clear interface contracts
   - Implement minimal working version (not stub!)
   - Example: Use in-memory storage instead of "TODO: add database"
   - Document upgrade path to full implementation

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

  r355_production_ready_checklist:
    no_hardcoded_values: true
    all_config_from_env: true
    no_stub_implementations: true
    no_todo_markers: true
    all_functions_complete: true

  configuration_approach:
    - name: "Database URL"
      wrong: 'dbURL := "postgres://localhost/mydb"'
      correct: 'dbURL := os.Getenv("DATABASE_URL")'
    - name: "API Timeout"
      wrong: 'timeout := 30 * time.Second'
      correct: 'timeout := config.GetDuration("api.timeout", 30*time.Second)'

  feature_flags_needed:
    - flag: "EFFORT_X_FEATURE_Y"
      purpose: "Control feature rollout (NOT to hide stubs!)"
      default: false
      location: "config/features.yaml"
      activation: "When ready for production traffic"

  interface_implementations:
    - interface: "IServiceZ"
      implementation: "InMemoryServiceZ"
      production_ready: true
      notes: "Fully functional in-memory version, not a stub"
  
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

## 🔴🔴🔴 SUPREME LAW R359: SIZE LIMITS ARE FOR NEW CODE ONLY 🔴🔴🔴

**PENALTY: IMMEDIATE TERMINATION (-1000%)**

### CRITICAL CLARIFICATION IN YOUR PLANS:
When specifying size estimates, you MUST clarify:

```markdown
## Size Limit Clarification (R359):
- The 800-line limit applies to NEW CODE YOU ADD
- Repository will grow by ~800 lines (EXPECTED)
- NEVER delete existing code to meet size limits
- Example: If repo has 10,000 lines and you add 800, total will be 10,800

## Implementation Size Estimate:
- NEW code to be added: ~750 lines
- Existing codebase: [current size]
- Expected total after: [current + 750]
```

### YOUR PLANS MUST NEVER SUGGEST:
❌ "Remove unrelated packages to fit limit"
❌ "Keep only the essential files"
❌ "Extract this module to make it standalone"
❌ "Delete unused code to make room"

### YOUR PLANS MUST EMPHASIZE:
✅ "ADD 800 lines of new functionality"
✅ "BUILD ON TOP of existing code"
✅ "EXTEND the current implementation"
✅ "ENHANCE with new features"

## 🔴🔴🔴 CRITICAL: Plan Storage Location (R383 SUPREME LAW) 🔴🔴🔴

### Plans MUST use sf_metadata_path helper from R383

**MANDATORY STORAGE PATTERN PER R383:**
```bash
# MANDATORY: Include R383 helper function (from rule-library/R383-metadata-file-timestamp-requirements.md)
sf_metadata_path() {
    local phase="$1"
    local wave="$2"
    local effort="$3"
    local filename="$4"
    local ext="$5"

    # Validate inputs
    if [[ -z "$phase" || -z "$wave" || -z "$effort" || -z "$filename" || -z "$ext" ]]; then
        echo "❌ R383 VIOLATION: Missing parameters to sf_metadata_path" >&2
        exit 1
    fi

    # Create directory structure
    local dir=".software-factory/phase${phase}/wave${wave}/${effort}"
    mkdir -p "$dir"

    # Generate unique timestamped filename
    local timestamp=$(date +%Y%m%d-%H%M%S)
    local full_path="${dir}/${filename}--${timestamp}.${ext}"

    echo "$full_path"
}

# Determine phase, wave, and effort name from context
PHASE="1"  # Get from context
WAVE="1"   # Get from context
EFFORT_NAME="go-containerregistry-image-builder"  # Get from context

# Create plan using R383-compliant helper
PLAN_FILE=$(sf_metadata_path "$PHASE" "$WAVE" "$EFFORT_NAME" "IMPLEMENTATION-PLAN" "md")

echo "📁 R383 COMPLIANT: Using .software-factory structure"
echo "📝 Plan will be saved as: $PLAN_FILE"
```

**EXAMPLE PATHS (R383 COMPLIANT):**
- Effort plan: `.software-factory/phase1/wave1/effort-name/IMPLEMENTATION-PLAN--20250104-143000.md`
- Review report: `.software-factory/phase1/wave1/effort-name/CODE-REVIEW-REPORT--20250104-153000.md`
- Split plan: `.software-factory/phase1/wave1/effort-name/SPLIT-PLAN--20250104-143000.md`
- Fix plan: `.software-factory/phase1/wave1/effort-name/FIX-PLAN--20250104-143000.md`
- Work log: `.software-factory/phase1/wave1/effort-name/work-log--20250104-143000.log`

**WHY THIS STRUCTURE (R383 SUPREME LAW):**
- R383: EVERY metadata file MUST have timestamp suffix
- R343: All metadata in .software-factory directory structure
- Prevents merge conflicts during integration
- Maintains perfect file uniqueness across all agents
- Timestamps prevent collision per R383
- SW Engineers know exactly where to look

### Creating the Plan File
```bash
# Use the R383-compliant path we generated
cat > "$PLAN_FILE" << 'EOF'
# Implementation Plan for ${EFFORT_NAME}
Created: $(date -Iseconds)
Location: $(pwd)/$PLAN_FILE
Phase: ${PHASE}
Wave: ${WAVE}

## Effort Metadata
[Plan content here...]
EOF

# Validate the file was created with proper timestamp (R383)
if [[ ! "$PLAN_FILE" =~ --[0-9]{8}-[0-9]{6}\. ]]; then
    echo "❌ R383 VIOLATION: Plan file missing timestamp!"
    exit 383
fi

# Commit the plan
git add "$PLAN_FILE"
git commit -m "feat: add implementation plan for ${EFFORT_NAME}

Plan location: $PLAN_FILE
Phase $PHASE, Wave $WAVE
R383 compliant: timestamp included"
git push

echo "✅ R383 COMPLIANT: Plan created with timestamp: $PLAN_FILE"
```

### SW Engineer Will Read Plan From State File (R340)
```bash
# R340: SW Engineer MUST read plan location from orchestrator-state-v3.json
# They will NEVER search for plans!

# SW Engineer's R340-compliant plan discovery:
get_plan_from_state() {
    EFFORT_NAME="$1"
    STATE_FILE="/workspaces/software-factory-template/orchestrator-state-v3.json"
    
    if command -v jq &> /dev/null; then
        PLAN_PATH=$(jq -r ".effort_repo_files.effort_plans[\"${EFFORT_NAME}\"].file_path" "$STATE_FILE")
    elif command -v yq &> /dev/null; then
        PLAN_PATH=$(yq ".effort_repo_files.effort_plans[\"${EFFORT_NAME}\"].file_path" "$STATE_FILE")
    fi
    
    if [ "$PLAN_PATH" = "null" ] || [ -z "$PLAN_PATH" ]; then
        echo "❌ R340 VIOLATION: No plan tracked for effort '$EFFORT_NAME'"
        exit 340
    fi
    
    if [ ! -f "$PLAN_PATH" ]; then
        echo "❌ Tracked plan does not exist at: $PLAN_PATH"
        exit 1
    fi
    
    echo "$PLAN_PATH"
}

# YOU MUST REPORT THE PLAN TO ORCHESTRATOR FOR TRACKING!
```

## 🔴🔴🔴 MANDATORY: Report Plan Location to Orchestrator (R340) 🔴🔴🔴

### R340: Planning File Metadata Tracking
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R340-planning-file-metadata-tracking.md`
**Criticality**: BLOCKING - Orchestrator must track all planning files

**AFTER CREATING ANY IMPLEMENTATION PLAN, YOU MUST:**

```markdown
## 📋 PLANNING FILE CREATED (R383 COMPLIANT)

**Type**: effort_plan
**Path**: /efforts/phase{X}/wave{Y}/{effort-name}/.software-factory/phase{X}/wave{Y}/{effort-name}/IMPLEMENTATION-PLAN--{YYYYMMDD-HHMMSS}.md
**Effort**: {effort-name}
**Phase**: {X}
**Wave**: {Y}
**Target Branch**: phase{X}/wave{Y}/{effort-name}
**Created At**: {ISO-8601-timestamp}
**R383 Compliance**: ✅ Timestamp included

ORCHESTRATOR: Please update effort_repo_files.effort_plans["{effort-name}"] in state file per R340
```

**EXAMPLE REPORT:**
```markdown
## 📋 PLANNING FILE CREATED (R383 COMPLIANT)

**Type**: effort_plan
**Path**: /efforts/phase1/wave2/buildah-builder-interface/.software-factory/phase1/wave2/buildah-builder-interface/IMPLEMENTATION-PLAN--20250120-100000.md
**Effort**: buildah-builder-interface
**Phase**: 1
**Wave**: 2
**Target Branch**: phase1/wave2/buildah-builder-interface
**Created At**: 2025-01-20T10:00:00Z
**R383 Compliance**: ✅ Timestamp included (--20250120-100000)

ORCHESTRATOR: Please update effort_repo_files.effort_plans["buildah-builder-interface"] in state file per R340
```

**VERIFICATION:**
```bash
# After reporting, verify orchestrator will be able to find it
echo "✅ Plan created at: $PLAN_PATH"
echo "📋 Orchestrator must update state file with this location"
echo "🔍 SW Engineer will read location from orchestrator-state-v3.json"
```

### R344: Report Metadata Location to State File

**MANDATORY: After creating IMPLEMENTATION-PLAN.md, MUST report location**

```bash
# After writing implementation plan, update state file with location (R344 MANDATORY)
PLAN_PATH="$(pwd)/$PLAN_FILE"  # Full path to the plan file
EFFORT_NAME="${EFFORT_NAME}"

# Update state file with plan location
yq -i ".metadata_locations.implementation_plans.\"$EFFORT_NAME\" = {
    \"file_path\": \"$PLAN_PATH\",
    \"created_at\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
    \"created_by\": \"code-reviewer\",
    \"agent_id\": \"$(hostname)-$$\",
    \"phase\": $PHASE,
    \"wave\": $WAVE,
    \"effort\": \"$EFFORT_NAME\"
}" "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json"

# Commit state update
cd "$CLAUDE_PROJECT_DIR"
git add orchestrator-state-v3.json
git commit -m "state: report implementation plan location per R344"
git push

echo "✅ R344 COMPLETE: Plan location reported to state file"
```

## General Responsibilities
Follow all general code-reviewer rules and the Software Factory state machine.

## Next Steps
Proceed with the standard workflow for the EFFORT_PLAN_CREATION state as defined in the state machine.


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

