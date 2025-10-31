# R211: Code Reviewer Implementation Planning from Architecture

**Category:** Critical Rules  
**Agents:** Code Reviewer, Orchestrator, Architect  
**Criticality:** MISSION CRITICAL - Implementation must follow architecture  
**Priority:** HIGHEST - Required for proper planning hierarchy

## 🚨 CODE REVIEWER TRANSLATES ARCHITECTURE INTO IMPLEMENTATION 🚨

The Code Reviewer takes the Architect's vision and creates concrete, actionable implementation plans that SW Engineers can execute.

## The Problem This Solves

Without proper translation from architecture to implementation:
- Architectural vision gets lost in implementation details
- Engineers lack concrete guidance on what to build
- Parallelization opportunities aren't clearly communicated
- Code reuse isn't properly specified
- Implementation drifts from architectural intent

## The Solution: Two-Level Implementation Planning

### Level 1: Phase Implementation Planning (After Architect's Phase Architecture)

When the Architect completes a PHASE-X-ARCHITECTURE-PLAN.md, the Code Reviewer creates the implementation plan:

```bash
# CODE REVIEWER RESPONSIBILITY: Create Phase Implementation Plan
create_phase_implementation_plan() {
    local PHASE="$1"
    local ARCH_PLAN="$CLAUDE_PROJECT_DIR/planning/phase${PHASE}/PHASE-${PHASE}-ARCHITECTURE-PLAN.md"
    local IMPL_PLAN="$CLAUDE_PROJECT_DIR/planning/phase${PHASE}/PHASE-${PHASE}-IMPLEMENTATION-PLAN.md"
    
    echo "═══════════════════════════════════════════════════════"
    echo "📝 R211: Creating Phase ${PHASE} Implementation Plan"
    echo "═══════════════════════════════════════════════════════"
    
    # 1. MANDATORY: Read architect's plan first
    if [ ! -f "$ARCH_PLAN" ]; then
        echo "❌ FATAL: No architecture plan found!"
        echo "❌ Architect must create architecture plan first!"
        exit 1
    fi
    
    echo "📖 Reading architecture plan: $ARCH_PLAN"
    # Extract key architectural decisions
    CONTRACTS=$(grep -A 20 "APIs and Contracts" "$ARCH_PLAN")
    PARALLELIZATION=$(grep -A 20 "Wave Parallelization Strategy" "$ARCH_PLAN")
    REUSE=$(grep -A 20 "Shared Libraries" "$ARCH_PLAN")
    
    # 2. Ensure planning directory exists and use template
    mkdir -p "$CLAUDE_PROJECT_DIR/planning/phase${PHASE}"
    cp templates/PHASE-IMPLEMENTATION-PLAN.md "$IMPL_PLAN"
    
    # 3. Translate architecture into implementation details
    cat >> "$IMPL_PLAN" << 'EOF'
## 📋 Implementation Based on Architecture

### Architecture Compliance
This implementation plan follows PHASE-${PHASE}-ARCHITECTURE-PLAN.md:
- ✅ All contracts from architecture included
- ✅ Parallelization strategy preserved
- ✅ Code reuse requirements mapped
- ✅ MVP vs nice-to-have classification maintained

### Wave Implementation Strategy
Based on architectural parallelization analysis:

#### Sequential Waves (Must Complete in Order)
- Wave 1: Contract definitions (BLOCKS all other waves)
- Wave 2: Core implementations (Depends on Wave 1)

#### Parallel Wave Groups
- Waves 3-4: Can execute simultaneously (independent domains)
- Waves 5-6: Can execute simultaneously (separate features)

### Concrete File Mappings
[Translate architectural pseudo-code into real file paths]

### Branch Strategy
Based on parallelization opportunities:
- Sequential branches: phase${PHASE}/wave1, phase${PHASE}/wave2
- Parallel branches: Can spawn multiple phase${PHASE}/wave[3-4] simultaneously

EOF
    
    echo "✅ Phase implementation plan created: $IMPL_PLAN"
}
```

### Level 2: Wave Implementation Planning (After Architect's Wave Architecture)

When the Architect completes a WAVE-Y-ARCHITECTURE-PLAN.md, the Code Reviewer creates the implementation:

```bash
# CODE REVIEWER RESPONSIBILITY: Create Wave Implementation Plan
create_wave_implementation_plan() {
    local PHASE="$1"
    local WAVE="$2"
    local ARCH_PLAN="$CLAUDE_PROJECT_DIR/planning/phase${PHASE}/wave${WAVE}/WAVE-${PHASE}-${WAVE}-ARCHITECTURE-PLAN.md"
    local IMPL_PLAN="$CLAUDE_PROJECT_DIR/planning/phase${PHASE}/wave${WAVE}/WAVE-${PHASE}-${WAVE}-IMPLEMENTATION-PLAN.md"
    
    echo "═══════════════════════════════════════════════════════"
    echo "📝 R211: Creating Wave ${WAVE} Implementation Plan"
    echo "═══════════════════════════════════════════════════════"
    
    # 1. Read project prefix from target-repo-config.yaml (R191)
    local PROJECT_PREFIX=$(yq '.branch_naming.project_prefix' "$SF_ROOT/target-repo-config.yaml")
    local BRANCH_PREFIX=""
    if [ -n "$PROJECT_PREFIX" ] && [ "$PROJECT_PREFIX" != "null" ]; then
        BRANCH_PREFIX="${PROJECT_PREFIX}/"
        echo "📋 Using project prefix for branches: $PROJECT_PREFIX"
    fi
    
    # 2. MANDATORY: Verify architecture plan exists
    if [ ! -f "$ARCH_PLAN" ]; then
        echo "❌ FATAL: No wave architecture plan!"
        echo "❌ Architect must create wave architecture first!"
        exit 1
    fi
    
    # 3. Ensure planning directory exists and use template
    mkdir -p "$CLAUDE_PROJECT_DIR/planning/phase${PHASE}/wave${WAVE}"
    cp templates/WAVE-IMPLEMENTATION-PLAN.md "$IMPL_PLAN"
    
    # 4. Fill in implementation details from architecture
    # - Convert pseudo-code to actual implementation
    # - Specify exact files to create/edit
    # - Define concrete import paths
    # - Set specific line count targets
    # 🚨 NEW: MANDATORY PARALLELIZATION INFO FOR EACH EFFORT
    # - Can Parallelize: Yes/No for each effort
    # - Parallel With: List of efforts that can run simultaneously
    # 🚨 MANDATORY: Include project prefix in branch names!
    # - Branch format: ${BRANCH_PREFIX}phase${PHASE}/wave${WAVE}/effort-name
    
    echo "✅ Wave implementation plan created: $IMPL_PLAN"
}
```

### Level 3: Effort Implementation Planning (Based on Wave Implementation Plan)

**CRITICAL CHANGE**: Code Reviewer now uses WAVE implementation plans, not phase plans!

**🚨 MANDATORY: ALL EFFORT HEADERS MUST BE COPIED FROM WAVE PLAN! 🚨**

```bash
# CODE REVIEWER RESPONSIBILITY: Create Effort Implementation Plan
create_effort_implementation_plan() {
    local PHASE="$1"
    local WAVE="$2"
    local EFFORT="$3"
    local WAVE_IMPL_PLAN="$CLAUDE_PROJECT_DIR/planning/phase${PHASE}/wave${WAVE}/WAVE-${PHASE}-${WAVE}-IMPLEMENTATION-PLAN.md"
    local EFFORT_DIR="efforts/phase${PHASE}/wave${WAVE}/${EFFORT}"
    local TIMESTAMP=$(date +%Y%m%d-%H%M%S)
    local EFFORT_PLAN="${EFFORT_DIR}/.software-factory/phase${PHASE}/wave${WAVE}/${EFFORT}/EFFORT-IMPLEMENTATION-PLAN--${TIMESTAMP}.md"
    
    echo "═══════════════════════════════════════════════════════"
    echo "📝 R211: Creating Effort Implementation Plan"
    echo "Based on Wave Implementation Plan (NOT Phase Plan!)"
    echo "═══════════════════════════════════════════════════════"
    
    # 1. MANDATORY: Read WAVE implementation plan
    if [ ! -f "$WAVE_IMPL_PLAN" ]; then
        echo "❌ FATAL: No wave implementation plan!"
        echo "❌ Code Reviewer must create wave plan first!"
        exit 1
    fi
    
    echo "📖 Reading wave implementation plan: $WAVE_IMPL_PLAN"
    
    # 2. NEW (R219): Read dependency effort plans FIRST
    echo "🔗 R219: Checking for dependency effort plans..."
    DEPENDENCIES=$(grep -A 5 "### Effort.*${EFFORT}" "$WAVE_IMPL_PLAN" | 
                   grep "**Dependencies**:" | 
                   sed 's/.*Dependencies**: //')
    
    if [ -n "$DEPENDENCIES" ] && [ "$DEPENDENCIES" != "None" ]; then
        echo "📖 R219: Reading dependency plans: $DEPENDENCIES"
        echo "🧠 THINKING about how dependencies influence this effort..."
        for dep in $DEPENDENCIES; do
            DEP_PLAN="efforts/phase${PHASE}/wave${WAVE}/${dep}/.software-factory/phase${PHASE}/wave${WAVE}/${dep}/EFFORT-IMPLEMENTATION-PLAN--*.md"
            if [ -f "$DEP_PLAN" ]; then
                echo "✅ Analyzing dependency: $dep"
                # Code Reviewer must THINK about:
                # - What interfaces/contracts to implement
                # - What libraries to import
                # - What patterns to follow
                # - How to integrate with dependency outputs
            fi
        done
        echo "✅ R219: Dependency analysis complete - planning accordingly"
    fi
    
    # 3. Extract ALL effort headers (CRITICAL - MUST INCLUDE PARALLELIZATION!)
    # Find the effort section and extract complete header block
    EFFORT_NAME=$(grep -A 0 "### Effort.*${EFFORT}" "$WAVE_IMPL_PLAN" | head -1)
    BRANCH=$(grep -A 1 "### Effort.*${EFFORT}" "$WAVE_IMPL_PLAN" | grep "**Branch**:")
    CAN_PARALLELIZE=$(grep -A 2 "### Effort.*${EFFORT}" "$WAVE_IMPL_PLAN" | grep "**Can Parallelize**:")
    PARALLEL_WITH=$(grep -A 3 "### Effort.*${EFFORT}" "$WAVE_IMPL_PLAN" | grep "**Parallel With**:")
    SIZE_ESTIMATE=$(grep -A 4 "### Effort.*${EFFORT}" "$WAVE_IMPL_PLAN" | grep "**Size Estimate**:")
    DEPENDENCIES=$(grep -A 5 "### Effort.*${EFFORT}" "$WAVE_IMPL_PLAN" | grep "**Dependencies**:")
    
    # 3. Extract effort-specific content
    EFFORT_SECTION=$(grep -A 100 "### Effort.*${EFFORT}" "$WAVE_IMPL_PLAN")
    FILES_TO_CREATE=$(grep -A 30 "Files to Create" <<< "$EFFORT_SECTION")
    FILES_TO_REUSE=$(grep -A 20 "Files to Reuse" <<< "$EFFORT_SECTION")
    IMPLEMENTATION_INSTRUCTIONS=$(grep -A 30 "Implementation Instructions" <<< "$EFFORT_SECTION")
    TEST_REQUIREMENTS=$(grep -A 20 "Test Requirements" <<< "$EFFORT_SECTION")

    # 4. Create effort directory with .software-factory structure
    mkdir -p "${EFFORT_DIR}/.software-factory/phase${PHASE}/wave${WAVE}/${EFFORT}"

    # 5. Create detailed effort plan WITH ALL HEADERS
    cat > "$EFFORT_PLAN" << 'EOF'
# Effort Implementation Plan: ${EFFORT}

## 🚨 CRITICAL EFFORT METADATA (FROM WAVE PLAN)
${EFFORT_NAME}
${BRANCH}
${CAN_PARALLELIZE}
${PARALLEL_WITH}
${SIZE_ESTIMATE}
${DEPENDENCIES}

## 📋 Source Information
**Wave Plan**: $CLAUDE_PROJECT_DIR/planning/phase${PHASE}/wave${WAVE}/WAVE-${PHASE}-${WAVE}-IMPLEMENTATION-PLAN.md
**Effort Number**: ${EFFORT_NUMBER}
**Extracted**: $(date '+%Y-%m-%d %H:%M:%S')

## 🚀 Parallelization Context
${CAN_PARALLELIZE}
${PARALLEL_WITH}
**Blocking Status**: [If Can Parallelize = No, this effort blocks others]
**Parallel Group**: [If Can Parallelize = Yes, which efforts run simultaneously]

## Files to Create
${FILES_TO_CREATE}

## Files to Import/Reuse
${FILES_TO_REUSE}

## Dependencies
${DEPENDENCIES}

## Implementation Instructions
${IMPLEMENTATION_INSTRUCTIONS}

## Test Requirements
${TEST_REQUIREMENTS}

## Size Constraints
${SIZE_ESTIMATE}
- Maximum: 800 lines (HARD LIMIT - will trigger split if exceeded)
- Measurement Tool: Use line-counter.sh from project root

## Completion Criteria
- [ ] All files created as specified
- [ ] Size under 800 lines (verified by line-counter.sh)
- [ ] All tests passing
- [ ] Dependencies properly imported
- [ ] Code reviewed and approved
- [ ] Work log updated

EOF
    
    echo "✅ Effort implementation plan created with ALL headers: $EFFORT_PLAN"
    echo "✅ Parallelization info preserved: ${CAN_PARALLELIZE}"
}
```

## Code Reviewer State Machine Integration

### New Code Reviewer States

1. **PHASE_IMPLEMENTATION_PLANNING**
   - Triggered after Architect creates phase architecture
   - Reads: $CLAUDE_PROJECT_DIR/planning/phaseX/PHASE-X-ARCHITECTURE-PLAN.md
   - Creates: $CLAUDE_PROJECT_DIR/planning/phaseX/PHASE-X-IMPLEMENTATION-PLAN.md
   - Transitions to: PLANNING_COMPLETE

2. **WAVE_IMPLEMENTATION_PLANNING**
   - Triggered after Architect creates wave architecture
   - Reads: $CLAUDE_PROJECT_DIR/planning/phaseX/waveY/WAVE-X-Y-ARCHITECTURE-PLAN.md
   - Creates: $CLAUDE_PROJECT_DIR/planning/phaseX/waveY/WAVE-X-Y-IMPLEMENTATION-PLAN.md
   - Transitions to: PLANNING_COMPLETE

3. **EFFORT_PLANNING** (UPDATED!)
   - Now reads WAVE implementation plan (not phase plan)
   - **NEW (R219)**: Reads dependency effort plans first
   - Analyzes how dependencies influence current effort
   - Creates effort-specific implementation plans
   - Extracts relevant section from wave plan

### State Flow
```
Architect: PHASE_ARCHITECTURE_PLANNING
    ↓
    Creates $CLAUDE_PROJECT_DIR/planning/phaseX/PHASE-X-ARCHITECTURE-PLAN.md
    ↓
Orchestrator spawns Code Reviewer
    ↓
Code Reviewer: PHASE_IMPLEMENTATION_PLANNING
    ↓
    Creates $CLAUDE_PROJECT_DIR/planning/phaseX/PHASE-X-IMPLEMENTATION-PLAN.md
    ↓
=========================================
    ↓
Architect: WAVE_ARCHITECTURE_PLANNING
    ↓
    Creates $CLAUDE_PROJECT_DIR/planning/phaseX/waveY/WAVE-X-Y-ARCHITECTURE-PLAN.md
    ↓
Orchestrator spawns Code Reviewer
    ↓
Code Reviewer: WAVE_IMPLEMENTATION_PLANNING
    ↓
    Creates $CLAUDE_PROJECT_DIR/planning/phaseX/waveY/WAVE-X-Y-IMPLEMENTATION-PLAN.md
    ↓
=========================================
    ↓
Orchestrator spawns Code Reviewer for effort
    ↓
Code Reviewer: EFFORT_PLANNING
    ↓
    Reads $CLAUDE_PROJECT_DIR/planning/phaseX/waveY/WAVE-X-Y-IMPLEMENTATION-PLAN.md (NEW!)
    ↓
    Creates effort/.software-factory/phaseX/waveY/effort/EFFORT-IMPLEMENTATION-PLAN--TIMESTAMP.md
```

## 🚨 MANDATORY: Header Preservation Requirements

### ALL Effort Headers MUST Be Preserved in Effort Plans

When creating effort-specific IMPLEMENTATION-PLAN.md files from the wave plan, the Code Reviewer MUST:

1. **Copy ALL headers EXACTLY** from the wave plan effort section
2. **Preserve parallelization metadata** (Can Parallelize, Parallel With)  
3. **Maintain size estimates** as specified in wave plan
4. **Keep branch names** consistent with wave plan
5. **Retain dependencies** as listed in wave plan

**CRITICAL**: The SW Engineer and Orchestrator depend on these headers being IDENTICAL between wave plan and effort plan!

## 🚨 MANDATORY: Parallelization Information Requirements

### Every Effort MUST Include Parallelization Metadata

When creating WAVE-IMPLEMENTATION-PLAN.md, the Code Reviewer MUST specify for EACH effort:

```yaml
effort_parallelization_metadata:
  effort_1:
    can_parallelize: false  # Blocking effort
    parallel_with: []       # None - must complete first
    reason: "Defines contracts that all other efforts depend on"
    
  effort_2:
    can_parallelize: false  # Blocking effort
    parallel_with: []       # None - depends on effort_1
    reason: "Shared libraries needed by feature efforts"
    
  effort_3:
    can_parallelize: true   # Can run in parallel
    parallel_with: [4, 5]   # Can run with efforts 4 and 5
    reason: "Independent feature with no shared state"
    
  effort_4:
    can_parallelize: true   # Can run in parallel
    parallel_with: [3, 5]   # Can run with efforts 3 and 5
    reason: "Independent feature with no shared state"
    
  effort_5:
    can_parallelize: true   # Can run in parallel
    parallel_with: [3, 4]   # Can run with efforts 3 and 4
    reason: "Independent feature with no shared state"
    
  effort_6:
    can_parallelize: false  # Must wait for parallel group
    parallel_with: []       # None - integrates 3, 4, 5
    reason: "Requires all features to be complete for integration"
```

### Template Section for Each Effort

```markdown
### Effort N: [EFFORT_NAME]
**Branch**: `[PROJECT_PREFIX/]phase[PHASE]/wave[WAVE]/effort-[name]`  
**Can Parallelize**: [Yes/No] (MANDATORY FIELD)
**Parallel With**: [List effort numbers or "None"] (MANDATORY FIELD)
**Size Estimate**: [NUMBER] lines (MUST be <800)  
**Dependencies**: [List dependent efforts] (MANDATORY FIELD)
```

**NOTE**: Branch MUST include project prefix from target-repo-config.yaml if defined!
Example with prefix: `tmc-workspace/phase1/wave1/effort-api`
Example without: `phase1/wave1/effort-api`

### Orchestrator Uses This Information

The orchestrator will read the parallelization metadata to determine spawning strategy:

```bash
# Orchestrator reads parallelization info
read_effort_parallelization() {
    local WAVE_IMPL_PLAN="$1"
    
    # Extract effort parallelization
    CAN_PARALLELIZE=$(grep -A 1 "Effort $EFFORT_NUM" "$WAVE_IMPL_PLAN" | 
                      grep "Can Parallelize:" | 
                      cut -d: -f2 | 
                      xargs)
                      
    PARALLEL_WITH=$(grep -A 2 "Effort $EFFORT_NUM" "$WAVE_IMPL_PLAN" | 
                    grep "Parallel With:" | 
                    cut -d: -f2 | 
                    xargs)
    
    if [ "$CAN_PARALLELIZE" = "Yes" ]; then
        echo "✅ Effort $EFFORT_NUM can be parallelized with: $PARALLEL_WITH"
        # Orchestrator can spawn multiple engineers
    else
        echo "⚠️ Effort $EFFORT_NUM is blocking - must complete first"
        # Orchestrator must wait for completion
    fi
}
```

## Translation Requirements

### From Architecture to Implementation

The Code Reviewer MUST translate:

1. **Pseudo-code → Real Code**
   ```go
   // Architecture (pseudo-code)
   type Service interface {
       Method() error
   }
   
   // Implementation (specific)
   // File: pkg/phase2/wave1/api/service.go
   // Lines: 45-67
   type UserService interface {
       CreateUser(ctx context.Context, user *User) (*User, error)
       GetUser(ctx context.Context, id string) (*User, error)
       // ... specific methods
   }
   ```

2. **Abstract Paths → Concrete Paths**
   ```yaml
   # Architecture
   "Use authentication from Phase 1"
   
   # Implementation
   import "github.com/project/pkg/phase1/common/auth"
   authValidator := auth.NewValidator(config.AuthConfig)
   ```

3. **Concepts → Files**
   ```yaml
   # Architecture
   "Implement caching layer"
   
   # Implementation
   new_files:
     - pkg/phase2/wave3/cache/redis.go (200 lines)
     - pkg/phase2/wave3/cache/memory.go (150 lines)
     - pkg/phase2/wave3/cache/interface.go (50 lines)
   ```

4. **Parallelization Strategy → Concrete Assignments**
   ```yaml
   # Architecture
   "Waves 3-4 can run in parallel"
   
   # Implementation
   parallel_execution:
     engineer_1: phase2/wave3/effort1
     engineer_2: phase2/wave4/effort1
     start_time: same
   ```

## Validation Checklist

### Before Creating Implementation Plan:
- [ ] Architecture plan exists and is complete
- [ ] All architectural decisions understood
- [ ] Contracts and APIs identified
- [ ] Parallelization strategy clear

### Implementation Plan Must:
- [ ] Reference source architecture plan
- [ ] Include all contracts from architecture
- [ ] Preserve parallelization opportunities
- [ ] Specify exact files and line counts
- [ ] Map all reuse requirements
- [ ] Provide concrete import paths

### After Implementation Plan:
- [ ] SW Engineers have clear instructions
- [ ] File paths are specific
- [ ] Dependencies are explicit
- [ ] Size estimates are realistic

## Critical Rules

### RULE 1: Architecture First, Implementation Second
```bash
if [ ! -f "ARCHITECTURE-PLAN.md" ]; then
    echo "❌ CANNOT create implementation without architecture!"
    exit 1
fi
```

### RULE 2: Wave Plans Override Phase Plans for Efforts
```bash
# OLD (WRONG):
EFFORT_SOURCE="PHASE-X-IMPLEMENTATION-PLAN.md"

# NEW (CORRECT):
EFFORT_SOURCE="PHASE-X-WAVE-Y-IMPLEMENTATION-PLAN.md"
```

### RULE 3: Preserve Architectural Intent
- If architecture says "parallel", implementation MUST enable parallel execution
- If architecture says "reuse", implementation MUST import and reuse
- If architecture says "contract first", implementation MUST define contracts in Effort 1

## Examples

### ✅ Good Implementation Plan from Architecture
```markdown
## Based on Architecture Plan
Source: $CLAUDE_PROJECT_DIR/planning/phase2/wave1/WAVE-2-1-ARCHITECTURE-PLAN.md#contracts

## Effort 1: Contracts (From Architecture Section 3.1)
The architect specified these contracts (page 5, lines 123-145).
Converting to implementation:

Files to Create:
- pkg/phase2/wave1/api/service.go (150 lines)
  Implements Service interface from architecture line 127
- pkg/phase2/wave1/api/models.go (200 lines)  
  Implements data models from architecture lines 134-141
```

### ❌ Bad Implementation Plan
```markdown
## Implementation Plan
I'll figure out what to implement based on my own ideas...
[No reference to architecture plan]
```

## Summary

- **Code Reviewer** creates implementation plans FROM architecture plans
- **Phase implementation** comes from phase architecture
- **Wave implementation** comes from wave architecture
- **Effort implementation** comes from WAVE implementation (not phase!)
- **Translation** from abstract to concrete is key
- **Preservation** of architectural intent is mandatory

## Integration with R219 (Dependency-Aware Planning)

When creating effort implementation plans, Code Reviewers MUST also follow R219:
- **READ** dependency effort plans BEFORE creating current plan
- **THINK** about how dependencies influence the current effort
- **ANALYZE** available interfaces, libraries, and patterns from dependencies
- **DOCUMENT** dependency context in the implementation plan
- **PLAN** implementation to properly integrate with and build upon dependencies