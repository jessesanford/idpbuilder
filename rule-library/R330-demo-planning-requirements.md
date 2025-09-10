# 🚨🚨🚨 RULE R330: Demo Planning Requirements

## Classification
- **Category**: Planning & Documentation
- **Criticality Level**: 🚨🚨🚨 BLOCKING
- **Enforcement**: MANDATORY for all effort plans
- **Penalty**: -25% to -50% for violations
- **Related Rules**: R291 (Integration Demo Requirement), R311 (Scope Control)

## The Rule

**EVERY effort plan created by Code Reviewer MUST include explicit demo requirements that specify:**
1. **WHAT** to demonstrate (objectives and features)
2. **HOW** to demonstrate it (exact scenarios)
3. **SIZE IMPACT** of demo artifacts on the 800-line limit

**This ensures demos are properly planned BEFORE implementation begins, preventing scope creep and ensuring consistent quality.**

## 🔴🔴🔴 SUPREME REQUIREMENT: DEMOS ARE PART OF SCOPE 🔴🔴🔴

**Demo artifacts COUNT toward the 800-line limit and MUST be included in size planning!**

```bash
# MANDATORY SIZE CALCULATION INCLUDING DEMOS
echo "📊 Effort Size Breakdown:"
echo "  Production Code: 500 lines"
echo "  Test Code: 150 lines"
echo "  Demo Script: 50 lines"    # ← MUST INCLUDE
echo "  Demo Docs: 80 lines"      # ← MUST INCLUDE
echo "  TOTAL: 780 lines (under 800 limit ✅)"
```

## Requirements

### 1. 🎯 Demo Objectives Section (MANDATORY)

**Every effort plan MUST define 3-5 specific demonstration objectives:**

```markdown
## 🎬 Demo Requirements (R330 Mandatory)

### Demo Objectives
- [ ] Demonstrate CREATE operation works with valid input
- [ ] Show proper error handling for invalid data
- [ ] Verify integration with upstream service
- [ ] Prove performance meets <100ms requirement
- [ ] Display proper logging and monitoring

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

### 3. 📏 Demo Size Impact Section (MANDATORY)

**MUST calculate and include demo overhead in total size:**

```markdown
### Demo Size Impact

#### Demo Artifacts Breakdown
```
demo-features.sh:     50 lines  # Executable script
DEMO.md:             80 lines  # Documentation
test-data/:          20 lines  # Sample data files
integration-hook.sh: 10 lines  # For wave integration
────────────────────────────────
TOTAL DEMO OVERHEAD: 160 lines
```

#### Impact on Effort Size
```
Production Code:     450 lines
Test Code:          150 lines
Demo Artifacts:     160 lines  # ← COUNTS TOWARD LIMIT!
────────────────────────────────
GRAND TOTAL:        760 lines (✅ under 800)
```

**WARNING**: Excluding demo size = automatic -50% penalty!
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
- 🚨 Demo size not included in total calculation
- 🚨 Demo causes effort to exceed 800 lines
- 🚨 No demo scenarios defined

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

## Related Rules
- **R291**: Integration Demo Requirement (demos must work)
- **R311**: Scope Control (demos are part of scope)
- **R007**: Size Limits (demos count toward 800 lines)
- **R304**: Line Counter Usage (measure including demos)

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

**"A demo without a plan is just hoping it works"**
**"Demo artifacts are CODE and count toward limits"**
**"Plan the demo, size the demo, implement the demo"**

Every effort needs a demo, every demo needs a plan, and every plan needs to account for demo size. This is how we ensure quality, consistency, and successful integrations!