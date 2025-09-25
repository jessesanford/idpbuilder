# Architect - PHASE_ARCHITECTURE_PLANNING State Rules

## State Context
This is the PHASE_ARCHITECTURE_PLANNING state for the architect.

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
# ALWAYS create in SF instance directory's phase-plans folder
cd "$SF_INSTANCE_DIR"
mkdir -p phase-plans

# Document naming convention is STRICT
DOC_NAME="phase-plans/PHASE-${PHASE}-ARCHITECTURE-PLAN.md"

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

## General Responsibilities
Follow all general architect rules and the Software Factory state machine.

## Next Steps
Proceed with the standard workflow for the PHASE_ARCHITECTURE_PLANNING state as defined in the state machine.
