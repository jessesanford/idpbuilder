# R212: Phase Directory Isolation Protocol

**Category:** Critical Rules  
**Agents:** Orchestrator, Architect, Code Reviewer, All Agents  
**Criticality:** MISSION CRITICAL - Phase isolation ensures proper organization  
**Priority:** HIGHEST - All phase work must be properly isolated

## 🚨 ALL PHASE WORK MUST BE PROPERLY ISOLATED 🚨

Just as efforts need directory isolation (R209), phases need clear directory boundaries and metadata to ensure all agents work in the correct locations.

## The Problem This Solves

Without phase-level directory metadata:
- Agents don't know where phase work should happen
- Files get created in wrong locations
- Phase integration becomes chaotic
- Cross-phase contamination occurs
- Directory structure becomes inconsistent

## The Solution: Phase Metadata Injection

### Part 1: Orchestrator MUST Inject Phase Metadata

Similar to R209 for efforts, the orchestrator MUST inject metadata into ALL phase-level plans:

```bash
# ORCHESTRATOR RESPONSIBILITY: Add metadata to phase plans
inject_phase_metadata() {
    local PHASE="$1"
    local PHASE_DIR="$(pwd)/phase${PHASE}"
    local PHASE_ARCH_PLAN="phase-plans/PHASE-${PHASE}-ARCHITECTURE-PLAN.md"
    local PHASE_IMPL_PLAN="phase-plans/PHASE-${PHASE}-IMPLEMENTATION-PLAN.md"
    
    echo "═══════════════════════════════════════════════════════"
    echo "🚨 R212: Injecting Phase Metadata"
    echo "═══════════════════════════════════════════════════════"
    
    # Function to inject metadata into a plan
    inject_into_plan() {
        local PLAN="$1"
        
        if [ ! -f "$PLAN" ]; then 
            echo "⚠️ Plan not found: $PLAN"; 
            return 1; 
        fi
        
        # Check if metadata already exists
        if grep -q "PHASE INFRASTRUCTURE METADATA" "$PLAN"; then 
            echo "✅ Metadata already present in $PLAN"; 
            return 0; 
        fi
        
        # Create metadata header
        cat > /tmp/phase_metadata.md << EOF
<!-- ⚠️ PHASE INFRASTRUCTURE METADATA (Added by Orchestrator per R212) ⚠️ -->
<!-- ALL AGENTS: YOU MUST ACKNOWLEDGE THIS PHASE DIRECTORY! -->
**PHASE_NUMBER**: ${PHASE}
**WORKING_DIRECTORY**: $(pwd)
**PHASE_ROOT**: $(pwd)/phase${PHASE}
**EFFORTS_ROOT**: $(pwd)/efforts/phase${PHASE}
**PHASE_PLANS_DIR**: $(pwd)/phase-plans
**INTEGRATION_BRANCH**: phase${PHASE}-integration
**BASE_BRANCH**: main

⚠️ **CRITICAL PHASE ISOLATION RULES** ⚠️
1. ALL phase ${PHASE} work MUST happen under: $(pwd)/phase${PHASE}/
2. ALL efforts MUST go under: $(pwd)/efforts/phase${PHASE}/
3. ALL wave work MUST go under: $(pwd)/efforts/phase${PHASE}/wave{N}/
4. Phase plans stay in: $(pwd)/phase-plans/
5. NEVER mix files between phases!
6. This is the isolation boundary for Phase ${PHASE}

**PHASE DIRECTORY STRUCTURE**:
\`\`\`
$(pwd)/
├── phase-plans/                          # Phase-level plans
│   ├── PHASE-${PHASE}-ARCHITECTURE-PLAN.md
│   ├── PHASE-${PHASE}-IMPLEMENTATION-PLAN.md
│   └── PHASE-${PHASE}-WAVE-*-*.md
├── phase${PHASE}/                        # Phase-specific code
│   ├── api/                             # Phase APIs
│   ├── lib/                             # Phase libraries
│   └── integration/                     # Phase integration
└── efforts/phase${PHASE}/                # Phase efforts
    ├── wave1/
    │   ├── effort-1/
    │   ├── effort-2/
    │   └── ...
    └── wave2/
        └── ...
\`\`\`
<!-- END PHASE METADATA -->

EOF
        
        # Prepend metadata to plan
        cat /tmp/phase_metadata.md "$PLAN" > /tmp/updated_plan.md
        mv /tmp/updated_plan.md "$PLAN"
        echo "✅ Metadata injected into $PLAN"
    }
    
    # Inject into both architecture and implementation plans
    inject_into_plan "$PHASE_ARCH_PLAN"
    inject_into_plan "$PHASE_IMPL_PLAN"
    
    # Also inject into wave plans
    for wave_plan in phase-plans/PHASE-${PHASE}-WAVE-*.md; do
        if [ -f "$wave_plan" ]; then 
            inject_into_plan "$wave_plan"; 
        fi
    done
    
    echo "✅ Phase ${PHASE} metadata injection complete"
}
```

### Part 2: Agent Phase Acknowledgment

ALL agents MUST acknowledge when working with phase plans:

```bash
# MANDATORY FOR ALL AGENTS READING PHASE PLANS
acknowledge_phase_isolation() {
    local PLAN="$1"
    
    echo "═══════════════════════════════════════════════════════"
    echo "🚨 R212: PHASE DIRECTORY ISOLATION CHECK 🚨"
    echo "═══════════════════════════════════════════════════════"
    
    if [ ! -f "$PLAN" ]; then 
        echo "❌ FATAL: Phase plan not found: $PLAN"; 
        exit 1; 
    fi
    
    # Extract phase metadata
    PHASE_NUM=$(grep "**PHASE_NUMBER**:" "$PLAN" | cut -d: -f2- | xargs)
    WORKING_DIR=$(grep "**WORKING_DIRECTORY**:" "$PLAN" | cut -d: -f2- | xargs)
    PHASE_ROOT=$(grep "**PHASE_ROOT**:" "$PLAN" | cut -d: -f2- | xargs)
    EFFORTS_ROOT=$(grep "**EFFORTS_ROOT**:" "$PLAN" | cut -d: -f2- | xargs)
    PHASE_PLANS=$(grep "**PHASE_PLANS_DIR**:" "$PLAN" | cut -d: -f2- | xargs)
    
    if [ -z "$PHASE_NUM" ] || [ -z "$PHASE_ROOT" ]; then 
        echo "❌ FATAL: No R212 metadata in phase plan!"; 
        echo "Orchestrator must inject metadata first!"; 
        exit 1; 
    fi
    
    # Print acknowledgment
    echo ""
    echo "📋 R212 PHASE ISOLATION ACKNOWLEDGMENT"
    echo "────────────────────────────────────────────────────────"
    echo "I, $(whoami), EXPLICITLY ACKNOWLEDGE:"
    echo ""
    echo "✅ PHASE IDENTIFIED:"
    echo "   Phase Number: ${PHASE_NUM}"
    echo "   Phase Root: ${PHASE_ROOT}"
    echo ""
    echo "✅ DIRECTORY STRUCTURE CONFIRMED:"
    echo "   Working Directory: ${WORKING_DIR}"
    echo "   Phase Root: ${PHASE_ROOT}"
    echo "   Efforts Root: ${EFFORTS_ROOT}"
    echo "   Plans Directory: ${PHASE_PLANS}"
    echo ""
    echo "✅ CURRENT LOCATION:"
    echo "   Current: $(pwd)"
    
    # Check if in correct location for the work being done
    CURRENT=$(pwd)
    if [[ "$CURRENT" == "$PHASE_ROOT"* ]]; then 
        echo "   Status: ✓ In phase ${PHASE_NUM} directory"; 
    elif [[ "$CURRENT" == "$EFFORTS_ROOT"* ]]; then 
        echo "   Status: ✓ In phase ${PHASE_NUM} efforts directory"; 
    elif [[ "$CURRENT" == "$PHASE_PLANS"* ]]; then 
        echo "   Status: ✓ In phase plans directory"; 
    elif [[ "$CURRENT" == "$WORKING_DIR" ]]; then 
        echo "   Status: ✓ In project root"; 
    else 
        echo "   Status: ⚠️ Not in phase ${PHASE_NUM} directory"; 
        echo "   WARNING: May need to navigate to correct location!"; 
    fi
    
    echo ""
    echo "📍 I UNDERSTAND AND AGREE TO:"
    echo "   1. ✓ Phase ${PHASE_NUM} work happens in: ${PHASE_ROOT}"
    echo "   2. ✓ Efforts go in: ${EFFORTS_ROOT}"
    echo "   3. ✓ Plans stay in: ${PHASE_PLANS}"
    echo "   4. ✓ NEVER mix files between phases"
    echo "   5. ✓ Maintain phase isolation boundaries"
    echo ""
    echo "🔒 PHASE ${PHASE_NUM} ISOLATION ACKNOWLEDGED"
    echo "═══════════════════════════════════════════════════════"
    
    # Export for use
    export PHASE_NUMBER="$PHASE_NUM"
    export PHASE_ROOT_DIR="$PHASE_ROOT"
    export PHASE_EFFORTS_ROOT="$EFFORTS_ROOT"
    
    # Create acknowledgment file for audit
    echo "R212 Phase ${PHASE_NUM} Acknowledged at $(date '+%Y-%m-%d %H:%M:%S')" >> .r212-phase-acknowledged
    echo "Agent: $(whoami)" >> .r212-phase-acknowledged
    echo "Directory: $(pwd)" >> .r212-phase-acknowledged
    echo "" >> .r212-phase-acknowledged
}
```

### Part 3: Agent-Specific Acknowledgments

#### Architect Acknowledgment
```bash
# When reviewing phase or creating phase architecture
architect_phase_acknowledgment() {
    local PHASE="$1"
    
    echo "═══════════════════════════════════════════════════════"
    echo "🏗️ ARCHITECT: Phase ${PHASE} Acknowledgment"
    echo "═══════════════════════════════════════════════════════"
    
    # Read phase plan and acknowledge
    if [ -f "phase-plans/PHASE-${PHASE}-ARCHITECTURE-PLAN.md" ]; then 
        acknowledge_phase_isolation "phase-plans/PHASE-${PHASE}-ARCHITECTURE-PLAN.md"; 
    elif [ -f "phase-plans/PHASE-${PHASE}-IMPLEMENTATION-PLAN.md" ]; then 
        acknowledge_phase_isolation "phase-plans/PHASE-${PHASE}-IMPLEMENTATION-PLAN.md"; 
    else 
        echo "❌ No phase plan found for Phase ${PHASE}"; 
        exit 1; 
    fi
    
    echo "ARCHITECT ACKNOWLEDGES:"
    echo "  - Will review all Phase ${PHASE} work"
    echo "  - Will ensure architectural consistency"
    echo "  - Will create plans in: phase-plans/"
    echo "  - Will review efforts in: efforts/phase${PHASE}/"
}
```

#### Code Reviewer Acknowledgment
```bash
# When creating implementation plans from architecture
code_reviewer_phase_acknowledgment() {
    local PHASE="$1"
    
    echo "═══════════════════════════════════════════════════════"
    echo "📝 CODE REVIEWER: Phase ${PHASE} Acknowledgment"
    echo "═══════════════════════════════════════════════════════"
    
    # Read architecture plan and acknowledge
    if [ -f "phase-plans/PHASE-${PHASE}-ARCHITECTURE-PLAN.md" ]; then 
        acknowledge_phase_isolation "phase-plans/PHASE-${PHASE}-ARCHITECTURE-PLAN.md"; 
    else
        echo "❌ No architecture plan for Phase ${PHASE}"
        exit 1
    fi
    
    echo "CODE REVIEWER ACKNOWLEDGES:"
    echo "  - Will create implementation from architecture"
    echo "  - Will place plans in: phase-plans/"
    echo "  - Will review efforts in: efforts/phase${PHASE}/"
    echo "  - Will ensure <800 lines per effort"
}
```

## Integration with Existing Rules

### Works with R209 (Effort Isolation)
- R212 handles phase-level isolation
- R209 handles effort-level isolation
- Both use similar metadata injection patterns

### Works with R210 (Architecture Planning)
- Architect acknowledges phase directory before creating architecture
- Architecture plans include phase metadata

### Works with R211 (Implementation Planning)
- Code Reviewer acknowledges phase directory before creating implementation
- Implementation plans include phase metadata

## Validation Script

```bash
#!/bin/bash
# validate-phase-isolation.sh

validate_phase_metadata() {
    local PHASE="$1"
    local ERRORS=0
    
    echo "Validating Phase ${PHASE} metadata..."
    
    # Check architecture plan
    if [ -f "phase-plans/PHASE-${PHASE}-ARCHITECTURE-PLAN.md" ]; then 
        if grep -q "PHASE INFRASTRUCTURE METADATA" "phase-plans/PHASE-${PHASE}-ARCHITECTURE-PLAN.md"; then 
            echo "✅ Architecture plan has metadata"; 
        else 
            echo "❌ Architecture plan missing metadata"; 
            ((ERRORS++)); 
        fi
    fi
    
    # Check implementation plan  
    if [ -f "phase-plans/PHASE-${PHASE}-IMPLEMENTATION-PLAN.md" ]; then 
        if grep -q "PHASE INFRASTRUCTURE METADATA" "phase-plans/PHASE-${PHASE}-IMPLEMENTATION-PLAN.md"; then 
            echo "✅ Implementation plan has metadata"; 
        else 
            echo "❌ Implementation plan missing metadata"; 
            ((ERRORS++)); 
        fi
    fi
    
    # Check for acknowledgment files
    if grep -q "Phase ${PHASE}" .r212-phase-acknowledged 2>/dev/null; then 
        echo "✅ Phase acknowledgments found"; 
        grep "Phase ${PHASE}" .r212-phase-acknowledged | tail -3; 
    else 
        echo "⚠️ No phase acknowledgments yet"; 
    fi
    
    # Check directory structure
    if [ -d "phase${PHASE}" ]; then 
        echo "✅ Phase directory exists: phase${PHASE}/"; 
    else 
        echo "⚠️ Phase directory not created yet"; 
    fi
    
    if [ -d "efforts/phase${PHASE}" ]; then 
        echo "✅ Efforts directory exists: efforts/phase${PHASE}/"; 
    else 
        echo "⚠️ Efforts directory not created yet"; 
    fi
    
    return $ERRORS
}
```

## Examples

### ✅ Correct Phase Plan with Metadata

```markdown
<!-- ⚠️ PHASE INFRASTRUCTURE METADATA (Added by Orchestrator per R212) ⚠️ -->
<!-- ALL AGENTS: YOU MUST ACKNOWLEDGE THIS PHASE DIRECTORY! -->
**PHASE_NUMBER**: 2
**WORKING_DIRECTORY**: /workspaces/project
**PHASE_ROOT**: /workspaces/project/phase2
**EFFORTS_ROOT**: /workspaces/project/efforts/phase2
**PHASE_PLANS_DIR**: /workspaces/project/phase-plans
**INTEGRATION_BRANCH**: phase2-integration
**BASE_BRANCH**: main

⚠️ **CRITICAL PHASE ISOLATION RULES** ⚠️
1. ALL phase 2 work MUST happen under: /workspaces/project/phase2/
[...]
<!-- END PHASE METADATA -->

# Phase 2 Architecture Plan
[rest of plan...]
```

### ❌ Wrong: Phase Plan without Metadata

```markdown
# Phase 2 Architecture Plan
[plan content without metadata]
```

## Enforcement Workflow

1. **Orchestrator** creates phase plan (architecture or implementation)
2. **Orchestrator** injects R212 metadata into plan
3. **Agent** reads phase plan
4. **Agent** MUST acknowledge phase isolation
5. **Agent** creates audit trail (.r212-phase-acknowledged)
6. **Agent** works within phase boundaries

## Summary

- **Phase-level metadata** ensures proper directory isolation
- **All agents** must acknowledge phase boundaries
- **Orchestrator** injects metadata like it does for efforts (R209)
- **Audit trail** tracks acknowledgments
- **Prevents** cross-phase contamination
- **Ensures** consistent directory structure