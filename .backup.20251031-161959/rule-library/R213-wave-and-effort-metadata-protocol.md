# R213: Wave and Effort Directory Metadata Protocol

**Category:** Critical Rules  
**Agents:** Orchestrator (PRIMARY), All Agents (CONSUMERS)  
**Criticality:** MISSION CRITICAL - Orchestrator owns ALL directory structure  
**Priority:** HIGHEST - Must happen BEFORE any wave/effort work begins

## 🚨 ORCHESTRATOR IS THE MASTER OF ALL DIRECTORY STRUCTURES 🚨

The Orchestrator is the SINGLE SOURCE OF TRUTH for all directory structures. It MUST inject metadata into ALL plans (phase, wave, and effort) to ensure consistency.

## The Problem This Solves

Without orchestrator-controlled wave metadata:
- Directory structures become inconsistent
- Different agents create conflicting structures
- Wave boundaries get confused
- Efforts end up in wrong waves
- Integration becomes chaotic
- No single source of truth for structure

## The Solution: Orchestrator-Driven Wave Metadata

### Core Principle: Orchestrator Creates ALL Metadata

The orchestrator MUST:
1. **Create** all directory structures
2. **Define** all paths in metadata
3. **Inject** metadata into ALL plans
4. **Enforce** that agents follow the metadata

### Wave Metadata Injection Protocol

```bash
# ORCHESTRATOR RESPONSIBILITY: Master of wave directory structure
inject_wave_metadata() {
    local PHASE="$1"
    local WAVE="$2"
    local WAVE_ARCH_PLAN="phase-plans/PHASE-${PHASE}-WAVE-${WAVE}-ARCHITECTURE-PLAN.md"
    local WAVE_IMPL_PLAN="phase-plans/PHASE-${PHASE}-WAVE-${WAVE}-IMPLEMENTATION-PLAN.md"
    
    echo "═══════════════════════════════════════════════════════"
    echo "🌊 R213: Injecting Wave ${WAVE} Directory Metadata"
    echo "Orchestrator is MASTER of directory structure"
    echo "═══════════════════════════════════════════════════════"
    
    # Calculate all paths (ORCHESTRATOR DECIDES!)
    local PROJECT_ROOT="$(pwd)"
    local PHASE_ROOT="${PROJECT_ROOT}/phase${PHASE}"
    local WAVE_ROOT="${PROJECT_ROOT}/efforts/phase${PHASE}/wave${WAVE}"
    local INTEGRATE_WAVE_EFFORTS_BRANCH="phase${PHASE}/wave${WAVE}-integration"
    
    # Count efforts in this wave (orchestrator knows!)
    local EFFORT_COUNT=$(ls -d ${WAVE_ROOT}/effort-* 2>/dev/null | wc -l || echo "0")
    
    # Function to inject metadata
    inject_into_wave_plan() {
        local PLAN="$1"
        
        if [ ! -f "$PLAN" ]; then 
            echo "⚠️ Plan not found: $PLAN (will inject when created)"; 
            return 1; 
        fi
        
        # Check if metadata exists
        if grep -q "WAVE INFRASTRUCTURE METADATA" "$PLAN"; then 
            echo "✅ Metadata already present in $PLAN"; 
            return 0; 
        fi
        
        # Create AUTHORITATIVE metadata
        cat > /tmp/wave_metadata.md << EOF
<!-- ⚠️ WAVE INFRASTRUCTURE METADATA (R213 - ORCHESTRATOR DEFINED) ⚠️ -->
<!-- 🚨 THIS IS THE AUTHORITATIVE DIRECTORY STRUCTURE 🚨 -->
<!-- ALL AGENTS MUST USE THESE EXACT PATHS! -->

**METADATA_SOURCE**: ORCHESTRATOR (Single Source of Truth)
**METADATA_VERSION**: 1.0
**GENERATED_AT**: $(date -Iseconds)
**GENERATED_BY**: orchestrator

## Wave Identity
**PHASE_NUMBER**: ${PHASE}
**WAVE_NUMBER**: ${WAVE}
**WAVE_NAME**: Phase ${PHASE} Wave ${WAVE}
**INTEGRATE_WAVE_EFFORTS_BRANCH**: ${INTEGRATE_WAVE_EFFORTS_BRANCH}
**BASE_BRANCH**: main

## Directory Structure (AUTHORITATIVE)
**PROJECT_ROOT**: ${PROJECT_ROOT}
**PHASE_ROOT**: ${PHASE_ROOT}
**WAVE_ROOT**: ${WAVE_ROOT}
**PHASE_PLANS_DIR**: ${PROJECT_ROOT}/phase-plans
**EFFORT_COUNT**: ${EFFORT_COUNT}

## Wave Directory Layout (EXACTLY AS CREATED BY ORCHESTRATOR)
\`\`\`
${WAVE_ROOT}/
├── effort-1/                    # First effort (always contracts)
│   ├── IMPLEMENTATION-PLAN.md   # Has R209 metadata
│   ├── work-log.md
│   └── pkg/                     # Code goes here
├── effort-2/                    # Second effort
│   ├── IMPLEMENTATION-PLAN.md
│   ├── work-log.md
│   └── pkg/
├── effort-3/                    # Parallel effort
│   └── ...
├── effort-4/                    # Parallel effort
│   └── ...
└── wave-integration/            # Wave integration point
    └── README.md
\`\`\`

## Effort Parallelization (ORCHESTRATOR DEFINED)
**SEQUENTIAL_EFFORTS**: effort-1, effort-2
**PARALLEL_EFFORTS**: effort-3, effort-4, effort-5
**INTEGRATE_WAVE_EFFORTS_EFFORT**: effort-6

## ⚠️ CRITICAL WAVE ISOLATION RULES ⚠️
1. ALL Wave ${WAVE} work MUST happen under: ${WAVE_ROOT}
2. EACH effort gets its own subdirectory: ${WAVE_ROOT}/effort-N/
3. ALL code goes in: ${WAVE_ROOT}/effort-N/pkg/
4. NO files outside designated directories
5. NO cross-wave file sharing
6. Integration ONLY in: ${WAVE_ROOT}/wave-integration/

## Git Branch Strategy (ORCHESTRATOR ENFORCED)
- **Wave Integration Branch**: ${INTEGRATE_WAVE_EFFORTS_BRANCH}
- **Effort Branches**: phase${PHASE}/wave${WAVE}/effort-N
- **All branches** created by orchestrator BEFORE agent spawn

## Validation Commands
\`\`\`bash
# Verify you're in correct wave directory
if [[ "\$(pwd)" != "${WAVE_ROOT}"* ]]; then 
    echo "ERROR: Not in Wave ${WAVE} directory!"; 
    exit 1; 
fi

# Verify correct branch
EXPECTED_BRANCH_PREFIX="phase${PHASE}/wave${WAVE}"
if [[ "\$(git branch --show-current)" != "\$EXPECTED_BRANCH_PREFIX"* ]]; then 
    echo "ERROR: Not on Wave ${WAVE} branch!"; 
    exit 1; 
fi
\`\`\`

<!-- END WAVE METADATA -->

EOF
        
        # Prepend metadata to plan
        cat /tmp/wave_metadata.md "$PLAN" > /tmp/updated_plan.md
        mv /tmp/updated_plan.md "$PLAN"
        echo "✅ R213: Wave metadata injected into $PLAN"
        echo "   Wave root: $WAVE_ROOT"
        echo "   Efforts: $EFFORT_COUNT"
    }
    
    # Inject into both architecture and implementation plans
    inject_into_wave_plan "$WAVE_ARCH_PLAN"
    inject_into_wave_plan "$WAVE_IMPL_PLAN"
    
    echo "✅ Wave ${WAVE} metadata injection complete"
    echo "🔒 Directory structure is now LOCKED"
}
```

### Effort Metadata Injection Protocol

The orchestrator MUST also inject metadata into each effort's implementation plan:

```bash
# ORCHESTRATOR RESPONSIBILITY: Master of effort metadata
inject_effort_metadata() {
    local PHASE="$1"
    local WAVE="$2"
    local EFFORT_NAME="$3"
    local EFFORT_NUM="$4"
    local WORKING_DIR="$5"
    local BRANCH="$6"
    
    echo "═══════════════════════════════════════════════════════"
    echo "🔧 R213: Injecting Effort ${EFFORT_NUM} Metadata"
    echo "Orchestrator is MASTER of effort structure"
    echo "═══════════════════════════════════════════════════════"
    
    local EFFORT_PLAN="${WORKING_DIR}/IMPLEMENTATION-PLAN.md"
    
    # Check for placeholder
    if grep -q "<!-- ⚠️ ORCHESTRATOR METADATA PLACEHOLDER - DO NOT REMOVE ⚠️ -->" "$EFFORT_PLAN"; then
        echo "📝 Found metadata placeholder in effort plan"
        
        # Create the metadata to inject
        local METADATA="# 🔧 EFFORT INFRASTRUCTURE METADATA
**WORKING_DIRECTORY**: ${WORKING_DIR}
**BRANCH**: ${BRANCH}
**EFFORT_NAME**: ${EFFORT_NAME}"
        
        # Replace the placeholder with actual metadata
        sed -i "s|<!-- ⚠️ ORCHESTRATOR METADATA PLACEHOLDER - DO NOT REMOVE ⚠️ -->|${METADATA}|" "$EFFORT_PLAN"
        
        echo "✅ Effort metadata injected into $EFFORT_PLAN"
        echo "   Working dir: $WORKING_DIR"
        echo "   Branch: $BRANCH"
        echo "   Effort: $EFFORT_NAME"
    else
        echo "⚠️ No placeholder found in $EFFORT_PLAN"
        echo "   Cannot inject effort metadata!"
        return 1
    fi
}
```

### Integration with Effort Creation

When orchestrator creates efforts, it MUST ensure consistency:

```bash
# ORCHESTRATOR: Create effort with matching metadata
create_effort_with_metadata() {
    local PHASE="$1"
    local WAVE="$2"
    local EFFORT="$3"
    local EFFORT_NUM="$4"
    
    # Use SAME paths as in wave metadata!
    local WAVE_ROOT="$(pwd)/efforts/phase${PHASE}/wave${WAVE}"
    local EFFORT_DIR="${WAVE_ROOT}/effort-${EFFORT_NUM}"
    
    echo "Creating effort directory (R213 compliant)..."
    mkdir -p "$EFFORT_DIR/pkg"
    cd "$EFFORT_DIR"
    
    # Create branch (orchestrator controls naming!)
    git checkout -b "phase${PHASE}/wave${WAVE}/effort-${EFFORT_NUM}"
    
    # Create files
    touch IMPLEMENTATION-PLAN.md
    touch work-log.md
    
    # NOW inject R209 metadata that MATCHES wave metadata
    inject_r209_metadata "$EFFORT" "$PHASE" "$WAVE"
    
    echo "✅ Effort created with consistent metadata"
    echo "   Effort path: $EFFORT_DIR"
    echo "   Matches wave metadata: YES"
}
```

### Agent Validation of Wave Metadata

Agents MUST validate they're using orchestrator-defined structure:

```bash
# ALL AGENTS: Validate wave directory from metadata
validate_wave_directory() {
    local WAVE_PLAN="$1"
    
    echo "═══════════════════════════════════════════════════════"
    echo "🔍 R213: Validating Wave Directory from Orchestrator"
    echo "═══════════════════════════════════════════════════════"
    
    if [ ! -f "$WAVE_PLAN" ]; then 
        echo "❌ FATAL: No wave plan with orchestrator metadata!"; 
        exit 1; 
    fi
    
    # Extract orchestrator-defined paths
    WAVE_ROOT=$(grep "**WAVE_ROOT**:" "$WAVE_PLAN" | cut -d: -f2- | xargs)
    METADATA_SOURCE=$(grep "**METADATA_SOURCE**:" "$WAVE_PLAN" | cut -d: -f2- | xargs)
    
    if [ "$METADATA_SOURCE" != "ORCHESTRATOR" ]; then 
        echo "❌ FATAL: Metadata not from orchestrator!"; 
        echo "   Found: $METADATA_SOURCE"; 
        echo "   Required: ORCHESTRATOR"; 
        exit 1; 
    fi
    
    # Validate current directory matches orchestrator's definition
    CURRENT=$(pwd)
    if [[ "$CURRENT" == "$WAVE_ROOT"* ]]; then 
        echo "✅ In orchestrator-defined wave directory"; 
    else 
        echo "❌ NOT in orchestrator-defined directory!"; 
        echo "   Orchestrator says: $WAVE_ROOT"; 
        echo "   You are in: $CURRENT"; 
        exit 1; 
    fi
    
    echo "✅ Following orchestrator's directory structure"
}
```

## Metadata Hierarchy

The orchestrator creates a complete hierarchy of metadata:

```
ORCHESTRATOR (Master of Structure)
    ↓
R212: Phase Metadata (phase-plans/PHASE-X-*.md)
    ├── Phase root: /project/phase1
    └── Efforts root: /project/efforts/phase1
        ↓
R213: Wave Metadata (phase-plans/PHASE-X-WAVE-Y-*.md)
    ├── Wave root: /project/efforts/phase1/wave1
    └── Effort directories: effort-1/, effort-2/, ...
        ↓
R209: Effort Metadata (efforts/.../IMPLEMENTATION-PLAN.md)
    ├── Effort root: /project/efforts/phase1/wave1/effort-1
    └── Code location: .../pkg/
```

## Validation Rules

### Orchestrator MUST:
1. Create ALL directories before injecting metadata
2. Ensure metadata matches actual directory structure
3. Inject WAVE metadata BEFORE spawning any wave agents
4. Inject EFFORT metadata BEFORE spawning SW Engineers for that effort
5. Never change structure after metadata injection

### Agents MUST:
1. Read wave metadata from plans
2. Validate they're in orchestrator-defined directories
3. NEVER create directories that don't match metadata
4. Report errors if structure doesn't match metadata

## Validation Script Integration

```bash
#!/bin/bash
# validate-wave-metadata.sh

validate_wave_metadata_consistency() {
    local PHASE="$1"
    local WAVE="$2"
    
    echo "Validating Wave ${WAVE} metadata consistency..."
    
    # Check wave plans have metadata
    WAVE_ARCH="phase-plans/PHASE-${PHASE}-WAVE-${WAVE}-ARCHITECTURE-PLAN.md"
    WAVE_IMPL="phase-plans/PHASE-${PHASE}-WAVE-${WAVE}-IMPLEMENTATION-PLAN.md"
    
    for plan in "$WAVE_ARCH" "$WAVE_IMPL"; do
        if [ -f "$plan" ]; then 
            if grep -q "METADATA_SOURCE.*ORCHESTRATOR" "$plan"; then 
                echo "✅ $plan has orchestrator metadata"; 
                
                # Extract and verify paths
                WAVE_ROOT=$(grep "**WAVE_ROOT**:" "$plan" | cut -d: -f2- | xargs); 
                
                # Check directory exists
                if [ -d "$WAVE_ROOT" ]; then 
                    echo "✅ Wave directory exists: $WAVE_ROOT"; 
                    
                    # Count actual efforts vs metadata
                    ACTUAL_EFFORTS=$(ls -d ${WAVE_ROOT}/effort-* 2>/dev/null | wc -l); 
                    METADATA_EFFORTS=$(grep "**EFFORT_COUNT**:" "$plan" | cut -d: -f2- | xargs); 
                    
                    if [ "$ACTUAL_EFFORTS" = "$METADATA_EFFORTS" ]; then 
                        echo "✅ Effort count matches metadata"; 
                    else 
                        echo "⚠️  Effort count mismatch!"; 
                        echo "   Metadata says: $METADATA_EFFORTS"; 
                        echo "   Actually found: $ACTUAL_EFFORTS"; 
                    fi; 
                else 
                    echo "❌ Wave directory missing: $WAVE_ROOT"; 
                fi; 
            else 
                echo "❌ Missing orchestrator metadata in $plan"; 
            fi; 
        fi
    done
}
```

## Examples

### ✅ Correct Wave Plan with Orchestrator Metadata

```markdown
<!-- ⚠️ WAVE INFRASTRUCTURE METADATA (R213 - ORCHESTRATOR DEFINED) ⚠️ -->
<!-- 🚨 THIS IS THE AUTHORITATIVE DIRECTORY STRUCTURE 🚨 -->
**METADATA_SOURCE**: ORCHESTRATOR (Single Source of Truth)
**WAVE_ROOT**: /workspaces/project/efforts/phase2/wave1
**EFFORT_COUNT**: 5
[... complete metadata ...]
<!-- END WAVE METADATA -->

# Wave 1 Architecture Plan
[content...]
```

### ❌ Wrong: Wave Plan without Orchestrator Metadata

```markdown
# Wave 1 Architecture Plan
[No metadata - agents don't know structure!]
```

## Timing Requirements and Penalties

### Wave Metadata Injection:
- **WHEN**: During SETUP_WAVE_INFRASTRUCTURE state
- **BEFORE**: Spawning any wave-level agents
- **PENALTY**: -30% for missing wave metadata

### Effort Metadata Injection:
- **WHEN**: During ANALYZE_IMPLEMENTATION_PARALLELIZATION state
- **BEFORE**: Spawning SW Engineers for efforts
- **SPECIFICALLY**: After creating effort plans, before agent spawn
- **PENALTY**: -20% for missing effort metadata per effort

### Critical Enforcement:
```bash
# ORCHESTRATOR: Must inject effort metadata BEFORE spawning SWEs
prepare_effort_for_swe() {
    local EFFORT_NAME="$1"
    local WORKING_DIR="$2"
    local BRANCH="$3"
    
    # R213 MANDATE: Inject metadata first!
    inject_effort_metadata "$PHASE" "$WAVE" "$EFFORT_NAME" "$EFFORT_NUM" "$WORKING_DIR" "$BRANCH"
    
    # Only spawn after metadata injection
    spawn_sw_engineer "$EFFORT_NAME"
}
```

## Summary

- **Orchestrator** is the MASTER of all directory structures
- **R213** ensures wave AND effort-level directory consistency
- **Metadata** is injected into ALL wave and effort plans
- **Wave metadata** injected during SETUP_WAVE_INFRASTRUCTURE
- **Effort metadata** injected during ANALYZE_IMPLEMENTATION_PARALLELIZATION
- **Agents** MUST follow orchestrator-defined paths
- **No improvisation** - structure is centrally controlled
- **Complete hierarchy**: Phase (R212) → Wave (R213) → Effort (R213+R209)