# Architect - PHASE_ARCHITECTURE_PLANNING State Rules

## State Context
This is the PHASE_ARCHITECTURE_PLANNING state for the architect within SF 3.0.

## SF 3.0 Phase Planning Context

In this state, the Architect creates comprehensive phase-level architectural plans:
- Reads phase requirements and vision from `state_machine.current_state` in orchestrator-state-v3.json
- Reviews previous phase outcomes and architectural decisions from orchestrator-state-v3.json metadata
- Creates phase architecture plan that orchestrator records in `metadata_locations.phase_architecture_plans` per R340
- Defines wave breakdown and architectural approach for entire phase
- All planning artifacts stored with timestamps and reported for atomic orchestrator-state-v3.json updates per R288

## Acknowledgment Required
Thank you for reading the rules file for the PHASE_ARCHITECTURE_PLANNING state.

**IMPORTANT**: Please report that you have successfully read the PHASE_ARCHITECTURE_PLANNING rules file.

Say: "✅ Successfully read PHASE_ARCHITECTURE_PLANNING rules for architect"

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

### 🔴🔴🔴 INCREMENTAL BRANCHING ARCHITECTURE (R308 - CORE TENANT) 🔴🔴🔴

When designing phase architecture, you MUST account for incremental branching:

1. **Design for Wave-to-Wave Building**
   - Wave 2 efforts will branch from Wave 1 integration
   - Wave 3 efforts will branch from Wave 2 integration
   - Interfaces must support gradual additions

2. **Plan Incremental Dependencies**
   - Each wave builds on the previous wave's completed work
   - No wave can assume direct access to main branch
   - Features must layer incrementally

3. **Architecture Supporting Incremental Flow**
   - Define stable interfaces that won't break with additions
   - Plan for cumulative functionality across waves
   - Design APIs that extend rather than replace

4. **Document Incremental Strategy**
   ```yaml
   incremental_design:
     wave_1_foundation: "Core interfaces and base services"
     wave_2_builds_on: "Wave 1 interfaces, adds concrete implementations"
     wave_3_extends: "Wave 2 services, adds advanced features"
     no_breaking_changes: true
   ```

**CRITICAL**: Your architecture MUST assume each wave sees all previous waves' work, NOT just main branch.

### 🔴🔴🔴 ATOMIC PR ARCHITECTURE REQUIREMENTS (R220 - SUPREME LAW) 🔴🔴🔴

When designing phase architecture, you MUST ensure EVERY effort is designed for atomic PR mergeability:

1. **Split Work for Independent Mergeability**
   - Each effort = ONE atomic PR to upstream main
   - Each effort must merge independently without breaking build
   - No multi-effort dependencies that prevent individual merging

2. **Design Interfaces for Gradual Implementation**
   - Define clear interface contracts early
   - Allow concrete implementations to be added incrementally
   - Support swapping implementations without breaking other code

3. **Plan Feature Flags for Incomplete Functionality**
   - Every incomplete feature MUST be behind a flag
   - Document flag names and activation conditions
   - Ensure features can be tested in isolation

4. **Ensure Build Continuity**
   - Each effort PR must leave build working
   - No "temporary breakage" between PRs
   - All tests must pass after each PR

5. **Document Dependencies and Merge Order**
   - Clearly specify any merge order requirements
   - Identify which efforts can merge in parallel
   - Note any interfaces that must be defined first

### 🚨🚨🚨 DOCUMENT LOCATION REQUIREMENTS (R303) 🚨🚨🚨

**MANDATORY: Phase architecture plans MUST be created in the correct location:**

```bash
# ALWAYS create in SF instance directory's planning folder (R550)
cd "$SF_INSTANCE_DIR"
mkdir -p "planning/phase${PHASE}"

# Document naming convention is STRICT (R550 standardized)
DOC_NAME="planning/phase${PHASE}/PHASE-ARCHITECTURE-PLAN.md"

# Create and commit to main branch
git checkout main
cat > "$DOC_NAME" << 'EOF'
# Phase ${PHASE} Architecture Plan
[Architecture content]
EOF

git add "$DOC_NAME"
git commit -m "architect: create Phase ${PHASE} architecture plan"
git push origin main
```

**NEVER create phase documents in:**
- ❌ Effort directories
- ❌ Integration workspaces
- ❌ Target repository branches
- ❌ Random locations

### Architecture Plan Must Include

```yaml
phase_atomic_design:
  every_effort_independently_mergeable: true
  feature_flags_planned:
    - flag: "FEATURE_X_ENABLED"
      purpose: "Controls feature X activation"
  interfaces_defined:
    - interface: "IServiceContract"
      purpose: "Allows gradual service implementation"
  build_remains_working: true
  merge_strategy: "Each effort direct to main"
```

**VIOLATION = -100% IMMEDIATE FAILURE**

## 🔍 ARCHITECTURE AUDIT CHECKLIST (SF 2.0 Preservation)

**CRITICAL**: Before creating the phase architecture plan, perform a systematic architecture audit to ensure informed planning decisions.

### Pre-Planning Architecture Audit

**Purpose**: Identify existing patterns, technical debt, and architectural improvements needed before designing new phase.

**Audit Steps** (complete before writing architecture plan):

1. **Review Existing Architecture (if not Phase 1)**
   - [ ] Read previous phase architecture plans (planning/phase*/PHASE-ARCHITECTURE-PLAN.md per R550)
   - [ ] Identify architectural patterns currently in use
   - [ ] Document successful patterns to preserve
   - [ ] Identify problematic patterns to avoid

2. **Assess Current System State**
   - [ ] Review orchestrator-state-v3.json for system status
   - [ ] Check integration-containers.json for integration patterns
   - [ ] Review bug-tracking.json for recurring architectural issues
   - [ ] Examine completed efforts for implemented patterns

3. **Technical Debt Analysis**
   - [ ] Identify shortcuts or compromises in previous phases
   - [ ] Document areas requiring architectural cleanup
   - [ ] Assess impact on current phase design
   - [ ] Plan technical debt paydown if critical

4. **Architectural Pattern Inventory**
   - [ ] List design patterns in use (Factory, Observer, Strategy, etc.)
   - [ ] Document library/framework choices and rationale
   - [ ] Identify consistency issues across codebase
   - [ ] Note areas needing pattern standardization

5. **Quality and Performance Review**
   - [ ] Review test coverage patterns from previous phases
   - [ ] Identify performance bottlenecks or concerns
   - [ ] Document security considerations
   - [ ] Assess scalability of current architecture

6. **Improvement Opportunities**
   - [ ] Identify architectural refactorings needed
   - [ ] Document patterns to introduce in this phase
   - [ ] Plan interface improvements
   - [ ] Consider architectural evolution for future phases

### Audit Output

**Create**: `planning/phase${PHASE}/PHASE-ARCHITECTURE-AUDIT.md` (R550 compliant location)

**Contents**:
```markdown
# Phase ${PHASE} Architecture Audit

## Audit Date
[Date]

## Previous Architecture Review
[Summary of previous phase architecture - N/A if Phase 1]

## Current System State
[Findings from orchestrator-state-v3.json, integration-containers.json, bug-tracking.json]

## Technical Debt Inventory
- [Item 1: Description, impact, recommended action]
- [Item 2: ...]

## Architectural Patterns in Use
- [Pattern 1: Where used, effectiveness, continue/change recommendation]
- [Pattern 2: ...]

## Quality Assessment
- Test Coverage: [Current state, gaps]
- Performance: [Concerns, measurements]
- Security: [Current posture, improvements needed]
- Scalability: [Current capacity, growth plan]

## Improvement Recommendations for Phase ${PHASE}
1. [Recommendation 1: Rationale, priority, implementation approach]
2. [Recommendation 2: ...]

## Architecture Decisions for Phase ${PHASE}
[Key architectural decisions informed by this audit]
```

### Usage in Phase Architecture Planning

**Workflow**:
1. ✅ Complete architecture audit checklist above
2. ✅ Create PHASE-${PHASE}-ARCHITECTURE-AUDIT.md
3. ✅ Commit audit document to main branch
4. ✅ Use audit findings to inform PHASE-${PHASE}-ARCHITECTURE-PLAN.md creation
5. ✅ Reference audit document in architecture plan for traceability

**Why This Matters** (SF 2.0 Preservation):
- SF 2.0 had dedicated ARCHITECTURE_AUDIT state ensuring systematic review
- SF 3.0 merged this into PHASE_ARCHITECTURE_PLANNING for efficiency
- Checklist preserves SF 2.0's systematic audit rigor
- Prevents uninformed architectural decisions
- Ensures continuity and learning from previous phases

## General Responsibilities
Follow all general architect rules and the Software Factory state machine.

## Next Steps
Proceed with the standard workflow for the PHASE_ARCHITECTURE_PLANNING state as defined in the state machine.


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

