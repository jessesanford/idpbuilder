# Rule R219: Code Reviewer Dependency-Aware Effort Planning

## Rule Statement
When creating an effort implementation plan, the Code Reviewer MUST:
1. **READ** all dependency effort plans that the current effort depends on
2. **ANALYZE** how those dependencies influence the current effort
3. **THINK** about interfaces, contracts, and patterns established by dependencies
4. **PLAN** the current effort to properly integrate with and build upon dependencies

## Criticality Level
**MANDATORY** - Failure to understand dependencies leads to integration failures and rework

## Enforcement Mechanism
- **Technical**: Read dependency plans before creating current plan
- **Behavioral**: Document dependency analysis in the plan
- **Validation**: Verify imports match dependency exports

## Core Principle

```
Dependencies Define Context → Read Their Plans First → Think About Integration → Plan Accordingly
```

## Detailed Requirements

### MANDATORY: Read All Dependency Effort Plans First

```bash
# CODE REVIEWER MUST READ DEPENDENCY PLANS BEFORE PLANNING
read_dependency_effort_plans() {
    local PHASE="$1"
    local WAVE="$2"
    local CURRENT_EFFORT="$3"
    local WAVE_PLAN="phase-plans/PHASE-${PHASE}-WAVE-${WAVE}-IMPLEMENTATION-PLAN.md"
    
    echo "═══════════════════════════════════════════════════════════════"
    echo "📖 R219: Reading Dependency Effort Plans"
    echo "Current Effort: $CURRENT_EFFORT"
    echo "═══════════════════════════════════════════════════════════════"
    
    # 1. Extract dependencies from wave plan
    echo "Analyzing dependencies for $CURRENT_EFFORT..."
    DEPENDENCIES=$(grep -A 5 "### Effort.*${CURRENT_EFFORT}" "$WAVE_PLAN" | 
                   grep "**Dependencies**:" | 
                   sed 's/.*Dependencies**: //')
    
    if [ -z "$DEPENDENCIES" ] || [ "$DEPENDENCIES" = "None" ]; then
        echo "✅ No dependencies - can proceed with planning"
        return 0
    fi
    
    echo "Found dependencies: $DEPENDENCIES"
    echo ""
    
    # 2. Read each dependency's implementation plan
    for dep in $DEPENDENCIES; do
        local DEP_DIR="efforts/phase${PHASE}/wave${WAVE}/${dep}"
        local DEP_PLAN="${DEP_DIR}/IMPLEMENTATION-PLAN.md"
        
        if [ ! -f "$DEP_PLAN" ]; then
            echo "⚠️ WARNING: Dependency $dep plan not found!"
            echo "This effort should not have been spawned yet!"
            continue
        fi
        
        echo "📖 Reading dependency plan: $DEP_PLAN"
        
        # Extract key information from dependency
        echo "Analyzing what $dep provides:"
        echo "- Exported functions/types"
        echo "- API contracts"
        echo "- Shared libraries"
        echo "- Configuration patterns"
        echo ""
    done
    
    echo "✅ All dependency plans analyzed"
}
```

### Analyze Dependency Influence

```bash
# ANALYZE HOW DEPENDENCIES AFFECT CURRENT EFFORT
analyze_dependency_influence() {
    local CURRENT_EFFORT="$1"
    local DEPENDENCIES="$2"
    
    echo "═══════════════════════════════════════════════════════════════"
    echo "🧠 R219: THINKING About Dependency Influence"
    echo "═══════════════════════════════════════════════════════════════"
    
    cat << 'EOF'
DEPENDENCY ANALYSIS CHECKLIST:

1. **Contracts & Interfaces**
   - What interfaces must I implement from dependencies?
   - What contracts are already defined that I must follow?
   - What types/structs are available for reuse?

2. **Code Reuse Opportunities**
   - What libraries from dependencies can I import?
   - What utilities are already implemented?
   - What patterns should I follow for consistency?

3. **Integration Points**
   - How will my effort integrate with dependency outputs?
   - What are the expected data flows?
   - What are the API boundaries?

4. **Constraints & Limitations**
   - What design decisions from dependencies constrain me?
   - What naming conventions are established?
   - What error handling patterns are used?

5. **Testing Considerations**
   - What test utilities are available from dependencies?
   - What mocks/stubs can I reuse?
   - What integration test patterns are established?
EOF
    
    echo ""
    echo "🤔 THINKING: How do these dependencies shape my implementation?"
    echo "📝 Documenting influence in implementation plan..."
}
```

### Document Dependency Context in Plan

```bash
# CREATE EFFORT PLAN WITH DEPENDENCY CONTEXT
create_dependency_aware_effort_plan() {
    local PHASE="$1"
    local WAVE="$2"
    local EFFORT="$3"
    local DEPENDENCIES="$4"
    local EFFORT_DIR="efforts/phase${PHASE}/wave${WAVE}/${EFFORT}"
    local EFFORT_PLAN="${EFFORT_DIR}/IMPLEMENTATION-PLAN.md"
    
    echo "═══════════════════════════════════════════════════════════════"
    echo "📝 R219: Creating Dependency-Aware Implementation Plan"
    echo "═══════════════════════════════════════════════════════════════"
    
    # Add dependency context section to plan
    cat >> "$EFFORT_PLAN" << EOF

## 🔗 Dependency Context (R219 Required)

### Dependencies Analyzed
$(for dep in $DEPENDENCIES; do
    echo "- **$dep**: [What this provides to current effort]"
done)

### How Dependencies Influence This Implementation

#### Contracts to Implement
[List interfaces from dependencies that this effort must implement]

#### Libraries to Import
\`\`\`go
import (
    // From Effort 1 (contracts)
    "pkg/phase${PHASE}/wave${WAVE}/api/interfaces"
    
    // From Effort 2 (libraries)
    "pkg/phase${PHASE}/wave${WAVE}/lib/client"
)
\`\`\`

#### Patterns to Follow
[Document patterns established by dependencies]

#### Integration Strategy
[Explain how this effort integrates with dependency outputs]

### Dependency-Driven Design Decisions

Based on analyzing the dependency implementation plans:

1. **Decision**: [Specific design choice influenced by dependencies]
   **Rationale**: [Why this decision based on dependency analysis]

2. **Decision**: [Another design choice]
   **Rationale**: [Dependency-based reasoning]

### What This Effort Provides to Dependents

For efforts that depend on this one, we will provide:
- [Exported functions/types]
- [Libraries for reuse]
- [Established patterns]

EOF
    
    echo "✅ Dependency context documented in plan"
}
```

## Complete Effort Planning Flow with R219

```bash
# COMPLETE FLOW: DEPENDENCY-AWARE EFFORT PLANNING
plan_effort_with_dependencies() {
    local PHASE="$1"
    local WAVE="$2"
    local EFFORT="$3"
    
    echo "═══════════════════════════════════════════════════════════════"
    echo "🎯 R219: Dependency-Aware Effort Planning"
    echo "Phase: $PHASE | Wave: $WAVE | Effort: $EFFORT"
    echo "═══════════════════════════════════════════════════════════════"
    
    # Step 1: Read dependency plans (R219)
    read_dependency_effort_plans "$PHASE" "$WAVE" "$EFFORT"
    
    # Step 2: Analyze influence (R219)
    analyze_dependency_influence "$EFFORT" "$DEPENDENCIES"
    
    # Step 3: Extract effort from wave plan (R211)
    extract_effort_from_wave_plan "$PHASE" "$WAVE" "$EFFORT"
    
    # Step 4: Create plan with dependency context (R219)
    create_dependency_aware_effort_plan "$PHASE" "$WAVE" "$EFFORT" "$DEPENDENCIES"
    
    # Step 5: Preserve all headers (R211)
    preserve_effort_headers "$PHASE" "$WAVE" "$EFFORT"
    
    echo "✅ Dependency-aware effort plan complete!"
}
```

## Example: Planning Effort 3 (Feature A) with Dependencies

### Given Dependencies
- Effort 1: Contracts & Interfaces (completed)
- Effort 2: Shared Libraries (completed)

### Code Reviewer Actions

```bash
# 1. Read Effort 1's plan
echo "Reading efforts/phase1/wave1/contracts/IMPLEMENTATION-PLAN.md"
# Discovers: Service interface, Request/Response types

# 2. Read Effort 2's plan  
echo "Reading efforts/phase1/wave1/libraries/IMPLEMENTATION-PLAN.md"
# Discovers: HTTP client, Logger, Error handling utilities

# 3. THINK about integration
echo "🤔 THINKING: Feature A must:"
echo "- Implement Service interface from Effort 1"
echo "- Use HTTP client from Effort 2"
echo "- Follow error patterns from Effort 2"

# 4. Create plan incorporating dependencies
cat > efforts/phase1/wave1/feature-a/IMPLEMENTATION-PLAN.md << 'EOF'
# Effort Implementation Plan: Feature A

## 🔗 Dependency Context (R219 Required)

### Dependencies Analyzed
- **contracts**: Provides Service interface, Request/Response types
- **libraries**: Provides HTTP client, logger, error utilities

### How Dependencies Influence This Implementation

#### Contracts to Implement
```go
// From contracts effort
type Service interface {
    ProcessRequest(ctx context.Context, req *Request) (*Response, error)
}
```

#### Libraries to Import
```go
import (
    "pkg/phase1/wave1/api/interfaces"  // From contracts
    "pkg/phase1/wave1/lib/client"       // From libraries
    "pkg/phase1/wave1/lib/errors"       // From libraries
)
```

#### Integration Strategy
Feature A will implement the Service interface using the shared HTTP
client for external calls, following the error handling patterns
established in the libraries effort.

[Rest of implementation plan...]
EOF
```

## Common Violations to Avoid

### ❌ Planning Without Reading Dependencies
```bash
# WRONG - Creates plan without understanding dependencies
create_effort_plan "feature-a"  # No dependency analysis!
```

### ❌ Ignoring Dependency Patterns
```bash
# WRONG - Reinvents what dependencies provide
implement_custom_error_handling()  # Libraries already has this!
```

### ✅ Correct Implementation
```bash
# RIGHT - Read, analyze, think, then plan
read_dependency_effort_plans
analyze_dependency_influence
create_dependency_aware_effort_plan
```

## Grading Impact

- **Creating plan without reading dependencies**: -25% (Integration risk)
- **No dependency analysis documentation**: -20% (Traceability failure)
- **Mismatched imports from dependencies**: -30% (Will cause build failures)
- **Reinventing dependency functionality**: -25% (Code duplication)

## Integration with Other Rules

- **R211**: Still must preserve all headers from wave plan
- **R218**: Orchestrator uses parallelization info for spawning
- **R053**: Dependencies determine parallelization possibilities
- **R210**: Architecture defines high-level dependencies

## Summary

**R219 ensures intelligent effort planning by:**
1. Reading all dependency implementation plans first
2. Analyzing how dependencies influence the current effort
3. Thinking about integration points and constraints
4. Planning implementation to properly build on dependencies
5. Documenting dependency context for traceability

This prevents integration failures, reduces rework, and ensures efforts build cohesively on each other.