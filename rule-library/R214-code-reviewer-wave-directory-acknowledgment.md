# R214: Code Reviewer Wave Directory Acknowledgment Protocol

**Category:** Critical Rules  
**Agents:** Code Reviewer (PRIMARY), Orchestrator  
**Criticality:** MISSION CRITICAL - Must use correct wave directory for effort planning  
**Priority:** HIGHEST - Prevents effort plans in wrong locations

## 🚨 CODE REVIEWERS MUST ACKNOWLEDGE WAVE DIRECTORY FROM METADATA 🚨

Code Reviewers MUST read wave implementation plan metadata (R213) and explicitly acknowledge they're using the correct directory structure when creating effort-specific plans.

## The Problem This Solves

Without explicit wave directory acknowledgment:
- Code reviewers might create effort plans in wrong directories
- Effort plans might not align with orchestrator's structure
- Directory paths in effort plans might contradict wave metadata
- SW Engineers get confused about where to work
- Integration becomes impossible due to misplaced files

## The Solution: Mandatory Wave Directory Acknowledgment

### Part 1: Code Reviewer MUST Read Wave Metadata First

Before creating ANY effort plan, the Code Reviewer MUST:
1. Read the wave implementation plan
2. Extract R213 wave metadata
3. Verify the metadata source is ORCHESTRATOR
4. Acknowledge the wave directory structure
5. Create effort plans using THOSE EXACT PATHS

```bash
# CODE REVIEWER MANDATORY ACKNOWLEDGMENT
code_reviewer_wave_acknowledgment() {
    local PHASE="$1"
    local WAVE="$2"
    local WAVE_IMPL_PLAN="phase-plans/PHASE-${PHASE}-WAVE-${WAVE}-IMPLEMENTATION-PLAN.md"
    
    echo "═══════════════════════════════════════════════════════"
    echo "🚨 R214: CODE REVIEWER WAVE DIRECTORY ACKNOWLEDGMENT 🚨"
    echo "═══════════════════════════════════════════════════════"
    
    # 1. MANDATORY: Wave implementation plan must exist
    if [ ! -f "$WAVE_IMPL_PLAN" ]; then 
        echo "❌ FATAL: No wave implementation plan found!"; 
        echo "❌ Cannot create effort plans without wave plan!"; 
        exit 1; 
    fi
    
    # 2. Extract R213 metadata
    if ! grep -q "WAVE INFRASTRUCTURE METADATA" "$WAVE_IMPL_PLAN"; then 
        echo "❌ FATAL: Wave plan missing R213 metadata!"; 
        echo "❌ Orchestrator must inject metadata first!"; 
        exit 1; 
    fi
    
    # 3. CRITICAL: Verify orchestrator is the source
    METADATA_SOURCE=$(grep "**METADATA_SOURCE**:" "$WAVE_IMPL_PLAN" | cut -d: -f2- | xargs)
    if [ "$METADATA_SOURCE" != "ORCHESTRATOR" ]; then 
        echo "❌ FATAL: Metadata not from orchestrator!"; 
        echo "   Source: $METADATA_SOURCE"; 
        echo "   Required: ORCHESTRATOR"; 
        echo "❌ Only orchestrator defines directory structure!"; 
        exit 1; 
    fi
    
    # 4. Extract wave directory structure
    WAVE_ROOT=$(grep "**WAVE_ROOT**:" "$WAVE_IMPL_PLAN" | cut -d: -f2- | xargs)
    EFFORT_COUNT=$(grep "**EFFORT_COUNT**:" "$WAVE_IMPL_PLAN" | cut -d: -f2- | xargs)
    PHASE_NUM=$(grep "**PHASE_NUMBER**:" "$WAVE_IMPL_PLAN" | cut -d: -f2- | xargs)
    WAVE_NUM=$(grep "**WAVE_NUMBER**:" "$WAVE_IMPL_PLAN" | cut -d: -f2- | xargs)
    
    # 5. EXPLICIT ACKNOWLEDGMENT
    echo ""
    echo "📋 R214 WAVE DIRECTORY ACKNOWLEDGMENT"
    echo "────────────────────────────────────────────────────────"
    echo "I, CODE REVIEWER, EXPLICITLY ACKNOWLEDGE:"
    echo ""
    echo "✅ WAVE METADATA SOURCE VERIFIED:"
    echo "   Source: ORCHESTRATOR ✓"
    echo "   Authority: Single Source of Truth ✓"
    echo ""
    echo "✅ WAVE STRUCTURE UNDERSTOOD:"
    echo "   Phase: ${PHASE_NUM}"
    echo "   Wave: ${WAVE_NUM}"
    echo "   Wave Root: ${WAVE_ROOT}"
    echo "   Effort Count: ${EFFORT_COUNT}"
    echo ""
    echo "✅ EFFORT DIRECTORIES I WILL USE:"
    for i in $(seq 1 $EFFORT_COUNT); do
        echo "   Effort $i: ${WAVE_ROOT}/effort-${i}/"
    done
    echo ""
    echo "📍 I UNDERSTAND AND AGREE TO:"
    echo "   1. ✓ Use ONLY orchestrator-defined paths"
    echo "   2. ✓ Create effort plans in: ${WAVE_ROOT}/effort-N/"
    echo "   3. ✓ Maintain consistency with wave metadata"
    echo "   4. ✓ Not create any new directory structures"
    echo "   5. ✓ Ensure all paths match R213 metadata"
    echo ""
    echo "🔒 WAVE ${WAVE_NUM} DIRECTORY STRUCTURE ACKNOWLEDGED"
    echo "═══════════════════════════════════════════════════════"
    
    # 6. Export for use during effort planning
    export WAVE_ROOT_DIR="$WAVE_ROOT"
    export WAVE_PHASE="$PHASE_NUM"
    export WAVE_NUMBER="$WAVE_NUM"
    export WAVE_EFFORT_COUNT="$EFFORT_COUNT"
    
    # 7. Create acknowledgment record
    echo "R214 Wave $WAVE_NUM Acknowledged at $(date '+%Y-%m-%d %H:%M:%S')" >> .r214-wave-acknowledged
    echo "Code Reviewer acknowledges wave root: $WAVE_ROOT" >> .r214-wave-acknowledged
    echo "Will create $EFFORT_COUNT effort plans" >> .r214-wave-acknowledged
    echo "" >> .r214-wave-acknowledged
    
    return 0
}
```

### Part 2: Code Reviewer MUST Use Acknowledged Paths

When creating effort plans, MUST use the paths from acknowledgment:

```bash
# CODE REVIEWER: Create effort plan with correct paths
create_effort_plan_with_acknowledgment() {
    local EFFORT_NUM="$1"
    local EFFORT_NAME="$2"
    
    # MANDATORY: Must have acknowledged wave first
    if [ -z "$WAVE_ROOT_DIR" ]; then 
        echo "❌ FATAL: Must acknowledge wave directory first (R214)!"; 
        exit 1; 
    fi
    
    local EFFORT_DIR="${WAVE_ROOT_DIR}/effort-${EFFORT_NUM}"
    local EFFORT_PLAN="${EFFORT_DIR}/IMPLEMENTATION-PLAN.md"
    
    echo "═══════════════════════════════════════════════════════"
    echo "📝 Creating Effort Plan with R214 Compliance"
    echo "═══════════════════════════════════════════════════════"
    
    # Verify directory exists (orchestrator should have created it)
    if [ ! -d "$EFFORT_DIR" ]; then 
        echo "❌ ERROR: Effort directory doesn't exist!"; 
        echo "   Expected: $EFFORT_DIR"; 
        echo "   Orchestrator must create directories first!"; 
        exit 1; 
    fi
    
    # Navigate to effort directory
    cd "$EFFORT_DIR"
    echo "✅ Navigated to: $(pwd)"
    
    # Verify we're in the right place
    if [[ "$(pwd)" != "$EFFORT_DIR" ]]; then 
        echo "❌ FATAL: Not in correct effort directory!"; 
        exit 1; 
    fi
    
    # Create effort plan using template
    cp templates/EFFORT-PLANNING-TEMPLATE.md "$EFFORT_PLAN"
    
    # Add R214 compliance note
    cat >> "$EFFORT_PLAN" << EOF

## R214 Compliance
This effort plan was created in the orchestrator-defined directory:
- Wave Root: ${WAVE_ROOT_DIR}
- Effort Directory: ${EFFORT_DIR}
- Metadata Source: ORCHESTRATOR (R213)
- Acknowledged by: Code Reviewer at $(date '+%Y-%m-%d %H:%M:%S')
EOF
    
    echo "✅ Effort plan created with R214 compliance"
    echo "   Location: $EFFORT_PLAN"
    echo "   Consistent with wave metadata: YES"
}
```

### Part 3: Validation During Effort Planning

```bash
# VALIDATE: Effort plan location matches wave metadata
validate_effort_plan_location() {
    local EFFORT_PLAN="$1"
    local WAVE_IMPL_PLAN="$2"
    
    echo "🔍 R214: Validating effort plan location..."
    
    # Get effort directory from plan location
    EFFORT_DIR=$(dirname "$EFFORT_PLAN")
    
    # Get expected wave root from metadata
    WAVE_ROOT=$(grep "**WAVE_ROOT**:" "$WAVE_IMPL_PLAN" | cut -d: -f2- | xargs)
    
    # Verify effort is under wave root
    if [[ "$EFFORT_DIR" == "$WAVE_ROOT/effort-"* ]]; then 
        echo "✅ Effort plan in correct wave directory"; 
        return 0; 
    else 
        echo "❌ VIOLATION: Effort plan in wrong directory!"; 
        echo "   Expected under: $WAVE_ROOT"; 
        echo "   Actually in: $EFFORT_DIR"; 
        return 1; 
    fi
}
```

## Integration with Existing Rules

### Works with R211 (Implementation from Architecture)
- Code Reviewer reads wave implementation plan (R211)
- Acknowledges wave directory (R214)
- Creates effort plans in correct locations

### Works with R213 (Wave Metadata)
- Reads orchestrator-defined wave metadata
- Verifies ORCHESTRATOR is source
- Uses exact paths from metadata

### Works with R209 (Effort Isolation)
- Creates effort plans with correct paths
- R209 metadata will match R214 acknowledged paths
- SW Engineers get consistent directory information

## Code Reviewer Workflow with R214

```bash
# COMPLETE CODE REVIEWER WORKFLOW
code_reviewer_effort_planning_workflow() {
    local PHASE="$1"
    local WAVE="$2"
    
    echo "Starting Code Reviewer Effort Planning Workflow..."
    
    # Step 1: Read wave implementation plan (R211)
    echo "Step 1: Reading wave implementation plan..."
    READ: phase-plans/PHASE-${PHASE}-WAVE-${WAVE}-IMPLEMENTATION-PLAN.md
    
    # Step 2: ACKNOWLEDGE wave directory (R214) 
    echo "Step 2: R214 - Acknowledging wave directory..."
    code_reviewer_wave_acknowledgment "$PHASE" "$WAVE"
    
    # Step 3: For each effort, create plan in correct location
    echo "Step 3: Creating effort plans in acknowledged directories..."
    for i in $(seq 1 $WAVE_EFFORT_COUNT); do
        echo "Creating effort $i plan..."
        create_effort_plan_with_acknowledgment "$i" "effort-$i"
    done
    
    echo "✅ All effort plans created with R214 compliance"
}
```

## Validation Script Extension

```bash
#!/bin/bash
# validate-r214-compliance.sh

validate_code_reviewer_acknowledgments() {
    echo "🔍 R214: Code Reviewer Wave Acknowledgment Validation"
    
    # Check acknowledgment file
    if [ -f ".r214-wave-acknowledged" ]; then 
        echo "✅ R214 acknowledgment file found"; 
        echo "Recent acknowledgments:"; 
        tail -5 .r214-wave-acknowledged; 
    else 
        echo "⚠️  No R214 acknowledgments yet"; 
    fi
    
    # For each wave with implementation plan
    for wave_plan in phase-plans/PHASE-*-WAVE-*-IMPLEMENTATION-PLAN.md; do
        if [ -f "$wave_plan" ]; then 
            WAVE_ROOT=$(grep "**WAVE_ROOT**:" "$wave_plan" | cut -d: -f2- | xargs); 
            
            echo "Checking wave: $(basename $wave_plan)"; 
            echo "  Wave root: $WAVE_ROOT"; 
            
            # Check if effort plans exist in correct locations
            if [ -d "$WAVE_ROOT" ]; then 
                for effort_dir in ${WAVE_ROOT}/effort-*; do 
                    if [ -d "$effort_dir" ]; then 
                        EFFORT_PLAN="$effort_dir/IMPLEMENTATION-PLAN.md"; 
                        if [ -f "$EFFORT_PLAN" ]; then 
                            # Check for R214 compliance note
                            if grep -q "R214 Compliance" "$EFFORT_PLAN"; then 
                                echo "  ✅ $(basename $effort_dir): R214 compliant"; 
                            else 
                                echo "  ⚠️  $(basename $effort_dir): Missing R214 note"; 
                            fi; 
                        fi; 
                    fi; 
                done; 
            fi; 
        fi; 
    done
}
```

## Examples

### ✅ Correct: Code Reviewer Acknowledgment Output
```
═══════════════════════════════════════════════════════
🚨 R214: CODE REVIEWER WAVE DIRECTORY ACKNOWLEDGMENT 🚨
═══════════════════════════════════════════════════════

📋 R214 WAVE DIRECTORY ACKNOWLEDGMENT
────────────────────────────────────────────────────────
I, CODE REVIEWER, EXPLICITLY ACKNOWLEDGE:

✅ WAVE METADATA SOURCE VERIFIED:
   Source: ORCHESTRATOR ✓
   Authority: Single Source of Truth ✓

✅ WAVE STRUCTURE UNDERSTOOD:
   Phase: 2
   Wave: 1
   Wave Root: /workspaces/project/efforts/phase2/wave1
   Effort Count: 5

✅ EFFORT DIRECTORIES I WILL USE:
   Effort 1: /workspaces/project/efforts/phase2/wave1/effort-1/
   Effort 2: /workspaces/project/efforts/phase2/wave1/effort-2/
   [...]

🔒 WAVE 1 DIRECTORY STRUCTURE ACKNOWLEDGED
═══════════════════════════════════════════════════════
```

### ❌ Wrong: Creating Effort Plan Without Acknowledgment
```bash
# WRONG - No wave acknowledgment
cd efforts/phase2/wave1/effort-1
cp template.md IMPLEMENTATION-PLAN.md
```

## Summary

- **Code Reviewers** MUST acknowledge wave directory from metadata
- **Verification** that ORCHESTRATOR is the metadata source
- **Explicit listing** of all effort directories to be used
- **Audit trail** created for compliance tracking
- **Effort plans** created ONLY in acknowledged directories
- **Complete alignment** with orchestrator's structure