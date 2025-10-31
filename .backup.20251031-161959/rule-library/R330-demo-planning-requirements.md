# 🚨🚨🚨 RULE R330: Demo Planning Requirements

## Classification
- **Category**: Planning & Documentation
- **Criticality Level**: 🚨🚨🚨 BLOCKING
- **Enforcement**: MANDATORY for all effort plans
- **Penalty**: -25% to -50% for violations
- **Related Rules**: R291 (Integration Demo Requirement), R331 (Demo Validation Protocol), R311 (Scope Control)
- **Modified**: 2025-10-06 (Added R331 compliance requirement)

## The Rule

**EVERY effort plan AND EVERY merge plan created by Code Reviewer MUST include explicit demo requirements that specify:**
1. **WHAT** to demonstrate (objectives and features)
2. **HOW** to demonstrate it (exact scenarios)
3. **SIZE IMPACT** of demo artifacts on the 800-line limit (for efforts)
4. **VALIDATION GATES** for integration demos (for merge plans)
5. **R331 COMPLIANCE** - demos must be real, not simulated (BLOCKING)

**This ensures demos are properly planned BEFORE implementation/integration begins, preventing scope creep and ensuring consistent quality. All demos MUST pass R331 validation (no simulation, external verification required).**

### Scope of This Rule

**R330 applies to ALL Code Reviewer planning states:**

1. **EFFORT_PLAN_CREATION**: Demo requirements for individual efforts
2. **WAVE_MERGE_PLANNING**: Demo execution plan for wave integrations
3. **PHASE_MERGE_PLANNING**: Demo execution plan for phase integrations (MANDATORY)
4. **PROJECT_MERGE_PLANNING**: Demo execution plan for project integration (MANDATORY)

**If Code Reviewer is creating a plan, that plan MUST include demo requirements.**

## 🔴🔴🔴 SUPREME REQUIREMENT: DEMOS ARE SEPARATE FROM LINE LIMITS 🔴🔴🔴

**CRITICAL CLARIFICATION: Demo artifacts do NOT count toward the 800-line implementation limit!**

The line counter tool automatically excludes:
- Demo scripts (demo-*, demos/*, example-*)
- Demo documentation (DEMO.md)
- Test files (*_test.go, test/*, tests/*)
- All non-implementation files

```bash
# CORRECT SIZE CALCULATION - DEMOS TRACKED SEPARATELY
echo "📊 Effort Size Breakdown:"
echo "  Implementation Code: 500 lines"  # ← ONLY THIS COUNTS for 800 limit
echo "  -------------------------------"
echo "  Test Code: 150 lines"            # ← Does NOT count
echo "  Demo Script: 50 lines"           # ← Does NOT count
echo "  Demo Docs: 80 lines"             # ← Does NOT count
echo "  -------------------------------"
echo "  Implementation: 500/800 ✅"
echo "  (Tests/demos excluded per R007)"
```

**However, demos MUST still be planned and sized to ensure effort completeness!**

## Requirements

### 1. 🎯 Demo Objectives Section (MANDATORY)

**Every effort plan MUST define 3-5 specific demonstration objectives:**

```markdown
## 🎬 Demo Requirements (R330 + R331 Mandatory)

### Demo Objectives
- [ ] Demonstrate CREATE operation works with valid input (R331: real execution)
- [ ] Show proper error handling for invalid data (R331: must actually fail)
- [ ] Verify integration with upstream service (R331: external verification)
- [ ] Prove performance meets <100ms requirement (R331: actual measurement)
- [ ] Display proper logging and monitoring

**R331 Compliance**: Demo must execute real code, verify external state, and be capable of failing

**Success Criteria**: All objectives checked = demo passes
```

**Requirements:**
- Minimum 3, maximum 5 objectives
- Each must be verifiable (not subjective)
- Must cover core functionality
- Should include error scenarios
- Must align with effort requirements

### 2. 📋 Demo Scenarios Section (MANDATORY)

**Define EXACT scenarios that will be implemented:**

```markdown
### Demo Scenarios (IMPLEMENT EXACTLY THESE)

#### Scenario 1: Create Valid Resource
- **Setup**: Clean database, service running on port 8080
- **Input**: `{"name": "test", "value": 42}`
- **Action**: `curl -X POST localhost:8080/api/resource -d @valid.json`
- **Expected Output**: 
  ```json
  {"id": "uuid", "status": "created", "name": "test"}
  ```
- **Verification**: HTTP 201, resource appears in GET list
- **Script Lines**: ~20 lines

#### Scenario 2: Handle Invalid Input
- **Setup**: Service running
- **Input**: `{"name": ""}` (empty name)
- **Action**: `curl -X POST localhost:8080/api/resource -d @invalid.json`
- **Expected Output**:
  ```json
  {"error": "name is required", "code": "VALIDATION_ERROR"}
  ```
- **Verification**: HTTP 400, no resource created
- **Script Lines**: ~15 lines

#### Scenario 3: Integration Verification
- **Setup**: Upstream service mock running
- **Action**: Create resource that triggers upstream call
- **Expected**: Upstream receives correct request
- **Verification**: Check mock server logs
- **Script Lines**: ~25 lines

**TOTAL SCENARIO LINES**: ~60 lines
```

**Requirements:**
- 2-4 scenarios per effort
- Each scenario must be complete (setup → action → verify)
- Include exact commands or code
- Specify expected outputs precisely
- Count lines for size estimation

### 3. 📏 Demo Size Planning Section (MANDATORY)

**MUST plan demo artifacts but they DON'T count toward 800-line limit:**

```markdown
### Demo Size Planning

#### Demo Artifacts (Excluded from line count per R007)
```
demo-features.sh:     50 lines  # Executable script
DEMO.md:             80 lines  # Documentation
test-data/:          20 lines  # Sample data files
integration-hook.sh: 10 lines  # For wave integration
────────────────────────────────
TOTAL DEMO FILES:   160 lines (NOT counted)
```

#### Effort Size Summary
```
Implementation:     450 lines  # ← ONLY this counts toward 800
────────────────────────────────
Tests:             150 lines  # Excluded per R007
Demos:             160 lines  # Excluded per R007
────────────────────────────────
Implementation:    450/800 ✅ (within limit)
```

**NOTE**: While demos don't count toward the line limit, they MUST still be planned and implemented as specified!
```

### 4. 🎬 Demo Deliverables Section (MANDATORY)

**List exact files that MUST be created:**

```markdown
### Demo Deliverables

Required Files:
- [ ] `demo-features.sh` - Main demo script (executable)
- [ ] `DEMO.md` - Demo documentation per template
- [ ] `test-data/valid.json` - Valid input examples
- [ ] `test-data/invalid.json` - Invalid input examples
- [ ] `.demo-config` - Demo environment settings

Integration Hooks:
- [ ] Export DEMO_READY=true when complete
- [ ] Provide integration point for wave demo
- [ ] Include cleanup function
```

## 🚨 Enforcement Protocol

### Code Reviewer (EFFORT_PLAN_CREATION state)
```bash
# MANDATORY: Include demo section in every plan
verify_demo_requirements() {
    local plan_file="$1"
    
    # Check for demo section
    if ! grep -q "Demo Requirements" "$plan_file"; then
        echo "❌ R330 VIOLATION: Missing demo requirements section!"
        return 1
    fi
    
    # Check for scenarios
    if ! grep -q "Demo Scenarios" "$plan_file"; then
        echo "❌ R330 VIOLATION: Missing demo scenarios!"
        return 1
    fi
    
    # Check for size impact
    if ! grep -q "Demo Size Impact" "$plan_file"; then
        echo "❌ R330 VIOLATION: Missing demo size calculation!"
        return 1
    fi
    
    echo "✅ R330 COMPLIANT: Demo requirements complete"
}
```

### SW Engineer (IMPLEMENTATION state)
```bash
# MANDATORY: Follow demo plan exactly
implement_demo_per_plan() {
    echo "📋 Loading demo requirements from plan..."
    
    # Extract scenarios from plan
    grep -A 20 "Demo Scenarios" IMPLEMENTATION-PLAN.md
    
    # Implement ONLY specified scenarios
    echo "✅ Will implement EXACTLY the planned scenarios"
    echo "❌ Will NOT add extra demo features"
    echo "❌ Will NOT over-engineer the demo"
}
```

### Integration Agent (PLANNING/MERGING states)
```bash
# Use demo plans for integration strategy
plan_integration_demo() {
    echo "🔍 Reviewing effort demo plans..."
    
    for effort in "${EFFORTS[@]}"; do
        # Check effort has demo plan
        if grep -q "Demo Requirements" "$effort/IMPLEMENTATION-PLAN.md"; then
            echo "✅ $effort has demo plan - including in integration"
        else
            echo "❌ $effort missing demo plan - R330 VIOLATION!"
        fi
    done
    
    # Create integration demo strategy based on plans
    create_wave_demo_orchestration
}
```

## 🔴 Failure Conditions

### Critical Violations (-50% penalty)
- 🚨 No demo requirements in effort plan
- 🚨 Demo size not planned/estimated
- 🚨 No demo scenarios defined
- 🚨 Demo deliverables not specified

### Major Violations (-25% penalty)
- ⚠️ Fewer than 3 demo objectives
- ⚠️ Scenarios lack verification steps
- ⚠️ Demo deliverables not listed
- ⚠️ Integration hooks missing

### Minor Violations (-10% penalty)
- Demo documentation incomplete
- Scenarios missing exact commands
- Size estimates significantly wrong
- Demo doesn't match plan exactly

## Success Criteria

Before approving ANY effort plan:
- ✅ Demo objectives clearly defined (3-5 items)
- ✅ Demo scenarios fully specified (2-4 scenarios)
- ✅ Demo size calculated and included in total
- ✅ Total size including demos < 800 lines
- ✅ Demo deliverables listed explicitly
- ✅ Integration hooks identified

## Examples

### ✅ CORRECT: Complete Demo Planning
```markdown
## 🎬 Demo Requirements (R330)

### Demo Objectives
- [ ] Show user creation with valid data
- [ ] Demonstrate validation errors
- [ ] Verify database persistence
- [ ] Prove API returns correct status codes

### Demo Scenarios
#### Scenario 1: Valid User Creation
- Setup: Empty database
- Action: POST /users with valid JSON
- Expected: 201 Created, user ID returned
- Lines: ~25

#### Scenario 2: Validation Error
- Setup: Service running
- Action: POST /users with missing email
- Expected: 400 Bad Request, error message
- Lines: ~20

### Demo Size Impact
Demo script: 50 lines
Demo docs: 80 lines
Test data: 20 lines
TOTAL: 150 lines (included in 750 total)

### Demo Deliverables
- [ ] demo-features.sh (executable)
- [ ] DEMO.md (documentation)
- [ ] test-data/ (sample files)
```

### ❌ WRONG: Missing Demo Planning
```markdown
## Implementation Plan

[Requirements and code details...]

<!-- NO DEMO SECTION! R330 VIOLATION! -->
```

### ❌ WRONG: Demo Size Not Counted
```markdown
## Size Calculation
Production: 600 lines
Tests: 190 lines
Total: 790 lines  <!-- WRONG! Forgot demo artifacts! -->

## Demo Requirements
[Demo details that would add 150 lines...]
<!-- Now total is 940 lines - OVER LIMIT! -->
```

## 🔴🔴🔴 PRE-PLANNED DEMO INFRASTRUCTURE (R504 + R330 INTEGRATE_WAVE_EFFORTS) 🔴🔴🔴

### Integration Demo Pre-Planning

**For ALL integrations (wave, phase, project), demo information MUST be pre-calculated during planning and stored in pre_planned_infrastructure:**

```json
"pre_planned_infrastructure": {
  "integrations": {
    "wave_integrations": {
      "phase1_wave1": {
        ...
        "demo_script_file": "$CLAUDE_PROJECT_DIR/efforts/phase1/wave1/integration-phase1-wave1/.software-factory/phase1/wave1/demo/phase1/wave1/integration/demo-wave-integration.sh",
        "demo_description": "Showcases wave 1 integration validating OCI image push with local registry",
        "demo_plan_file": "$CLAUDE_PROJECT_DIR/efforts/phase1/wave1/integration-phase1-wave1/.software-factory/phase1/wave1/demo/demo-plan.md"
      }
    },
    "phase_integrations": {
      "phase1": {
        ...
        "demo_script_file": "$CLAUDE_PROJECT_DIR/efforts/phase1/integration-phase1/.software-factory/phase1/demo/phase1/integration/demo-phase-integration.sh",
        "demo_description": "Showcases complete phase 1 OCI workflow with authentication and registry push",
        "demo_plan_file": "$CLAUDE_PROJECT_DIR/efforts/phase1/integration-phase1/.software-factory/phase1/demo/demo-plan.md"
      }
    },
    "project_integration": {
      "demo_script_file": "$CLAUDE_PROJECT_DIR/efforts/project-integration/.software-factory/demo/project/demo-project-integration.sh",
      "demo_description": "Showcases complete project OCI solution across all phases",
      "demo_plan_file": "$CLAUDE_PROJECT_DIR/efforts/project-integration/.software-factory/demo/demo-plan.md"
    }
  }
}
```

**Architect/Planning Responsibilities:**
1. **During wave planning**: Pre-calculate wave integration demo info (script path, description, plan file)
2. **During phase planning**: Pre-calculate phase integration demo info
3. **During project init**: Pre-calculate project integration demo info
4. **Update demo_description**: Whenever demo plan changes, update description as summary

**Integration Agent Responsibilities:**
1. **READ pre-planned demo info**: Load demo_script_file, demo_description, demo_plan_file from orchestrator-state-v3.json
2. **USE pre-planned paths**: Create demo at exact pre-planned demo_script_file location
3. **FOLLOW demo plan**: Implement demo based on demo_plan_file content
4. **NO runtime decisions**: Never calculate demo paths or names dynamically

**Demo Validation Agent Responsibilities:**
1. **Find demos using pre-planned paths**: Read demo_script_file from orchestrator-state-v3.json
2. **Execute demos from pre-planned locations**: Use exact paths from pre_planned_infrastructure
3. **Verify against pre-planned description**: Check demo matches demo_description

### Pre-Planning Enforcement

```bash
# Architect MUST populate during planning
if ! jq -e '.pre_planned_infrastructure.integrations.wave_integrations.phase1_wave1.demo_script_file' orchestrator-state-v3.json > /dev/null 2>&1; then
    echo "❌ R330 + R504 VIOLATION: No demo_script_file pre-planned for wave integration!"
    exit 330
fi

# Integration agent MUST USE pre-planned info
DEMO_SCRIPT_FILE=$(jq -r '.pre_planned_infrastructure.integrations.wave_integrations.phase1_wave1.demo_script_file' orchestrator-state-v3.json)
if [ -z "$DEMO_SCRIPT_FILE" ]; then
    echo "❌ R504 VIOLATION: Demo path must come from pre_planned_infrastructure!"
    exit 504
fi

# Create demo at PRE-PLANNED location
mkdir -p "$(dirname "$DEMO_SCRIPT_FILE")"
cat > "$DEMO_SCRIPT_FILE" << 'EOF'
#!/bin/bash
# Demo content from demo plan...
EOF
```

## Related Rules
- **R291**: Integration Demo Requirement (demos must work)
- **R331**: Demo Validation Protocol (demos must be real, not simulated)
- **R504**: Pre-Infrastructure Planning (demos pre-planned with paths)
- **R311**: Scope Control (demos are part of scope)
- **R007**: Size Limits (demos excluded from 800 line count)
- **R304**: Line Counter Usage (measure excluding demos)

## Migration Guide

### For Existing Efforts Without Demo Plans
1. Review implementation to identify demo needs
2. Create demo requirements retroactively
3. Verify size still under 800 with demos
4. Update effort plan with demo section
5. Implement demos per new requirements

### For New Efforts
1. Code Reviewer MUST include demo section
2. Size demos during planning (typically 100-150 lines)
3. Ensure total with demos < 800 lines
4. SW Engineer implements exactly as specified
5. Integration Agent uses for orchestration

## Remember

**"Every integration needs a demo, no exceptions"**
**"Single-effort waves still integrate and demo"**
**"Integration demos prove the whole is greater than the parts"**

Integration demos are MANDATORY at wave, phase, and project levels. Even seemingly trivial integrations (single-effort waves, single-wave phases) MUST have demos to ensure the integration process works correctly and nothing breaks when components combine!