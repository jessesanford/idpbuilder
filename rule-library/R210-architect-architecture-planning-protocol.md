# R210: Architect Architecture Planning Protocol

**Category:** Critical Rules  
**Agents:** Architect, Orchestrator, Code Reviewer  
**Criticality:** MISSION CRITICAL - Architecture drives all implementation  
**Priority:** HIGHEST - Must happen before any implementation planning

## 🚨 ARCHITECT OWNS ARCHITECTURAL VISION AND ALIGNMENT 🚨

The Architect is the guardian of the project's architectural integrity, ensuring all phases and waves align with the master plan while adapting based on implementation learnings.

## The Problem This Solves

Without architectural planning BEFORE implementation planning:
- Implementation drifts from original vision
- APIs and contracts emerge organically (inconsistently)
- Parallelization opportunities are missed
- Technical debt accumulates from phase to phase
- Integration becomes painful due to misaligned interfaces

## The Solution: Two-Level Architecture Planning

### Level 1: Phase Architecture Planning (After Previous Phase Passes)

When a PHASE_ASSESSMENT passes, the Architect creates the next phase's architecture:

```bash
# ARCHITECT RESPONSIBILITY: Create Phase Architecture Plan
create_phase_architecture_plan() {
    local PHASE="$1"
    local PHASE_PLAN="$CLAUDE_PROJECT_DIR/phase-plans/phase${PHASE}/PHASE-${PHASE}-ARCHITECTURE-PLAN--$(date +%Y%m%d-%H%M%S).md"

    echo "═══════════════════════════════════════════════════════"
    echo "🏗️ R210: Creating Phase ${PHASE} Architecture Plan"
    echo "═══════════════════════════════════════════════════════"

    # 0. Ensure planning directory exists
    mkdir -p "$CLAUDE_PROJECT_DIR/phase-plans/phase${PHASE}"

    # 1. Analyze what was built in previous phases
    echo "📊 Analyzing previous phase implementations..."
    analyze_previous_phases "$PHASE"

    # 2. Review master plan for this phase's goals
    echo "📋 Reviewing master plan requirements..."
    review_master_plan_for_phase "$PHASE"

    # 3. Create architecture plan
    cat > "$PHASE_PLAN" << 'EOF'
# Phase ${PHASE} Architecture Plan

## 🎯 Phase Vision Alignment
[How this phase aligns with master IMPLEMENTATION-PLAN.md]

## 📊 Analysis of Previous Phases
### What We've Built
- APIs: [List established APIs]
- Contracts: [List established contracts]
- Abstractions: [List key abstractions]
- Libraries: [List reusable libraries]

### What We've Learned
- [Key learnings that affect this phase]
- [Adjustments needed based on implementation reality]

## 🏗️ Phase ${PHASE} Architecture

### Core Architectural Decisions
1. **Decision**: [e.g., Use event-driven architecture]
   **Rationale**: [Why this decision]
   **Impact**: [How it affects implementation]

### APIs and Contracts (MUST BE FIRST!)
```go
// Example API Contract
type PhaseService interface {
    // Methods that MUST be implemented
    Initialize(ctx context.Context) error
    Process(req Request) (Response, error)
}

// Example Data Contract
type Request struct {
    ID       string                 `json:"id"`
    Metadata map[string]interface{} `json:"metadata"`
}
```

### Abstractions and Interfaces
```go
// Core abstraction for this phase
type DataProcessor interface {
    Validate(data []byte) error
    Transform(data []byte) ([]byte, error)
    Store(ctx context.Context, data []byte) error
}
```

### Shared Libraries (From Previous Phases)
- `pkg/common/auth`: Authentication from Phase 1
- `pkg/common/logger`: Logging from Phase 1
- Usage Example:
```go
import "pkg/common/auth"

func NewService() *Service {
    return &Service{
        auth: auth.NewValidator(), // Reuse from Phase 1
    }
}
```

### CRDs/Schemas (If Applicable)
```yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: services.example.com
spec:
  # Schema definition
```

### Plugin System Architecture (If Applicable)
```go
type Plugin interface {
    Name() string
    Version() string
    Execute(ctx context.Context, input interface{}) (interface{}, error)
}
```

## 🔄 Wave Parallelization Strategy

### Sequential Waves (Dependencies)
- Wave 1: MUST complete first (establishes contracts)
- Wave 2: Depends on Wave 1 contracts

### Parallel Waves (Independent)
- Waves 3 & 4: Can run in parallel (use same contracts)
- Waves 5 & 6: Can run in parallel (independent features)

### Parallelization Enablers
1. **Contract Definition**: Wave 1 defines all interfaces
2. **Mock Implementations**: Wave 1 provides mocks for testing
3. **Independent Domains**: Waves 3-4 work on separate domains

## 📈 MVP vs Nice-to-Have

### MVP Features (Early Waves)
- Core authentication
- Basic CRUD operations
- Essential API endpoints

### Nice-to-Have Features (Later Waves)
- Advanced analytics
- Performance optimizations
- Extended plugin support

## 🔧 Technical Specifications

### Performance Requirements
- Latency: <100ms for API calls
- Throughput: 1000 req/sec
- Memory: <500MB per service

### Security Requirements
- All APIs must use JWT authentication
- Data encryption at rest and in transit
- RBAC implementation required

### Integration Points
- Previous Phase APIs: [List what to integrate with]
- External Systems: [List external dependencies]
- Future Phase Hooks: [List extension points for next phase]

## ✅ Success Criteria
- All contracts defined before implementation
- Parallelization achieved where planned
- Previous phase code successfully reused
- Architecture remains aligned with master plan

## 🚨 Risk Mitigations
- Risk: API changes break previous phases
  Mitigation: Versioned APIs, backward compatibility
- Risk: Parallel waves conflict
  Mitigation: Clear domain boundaries, integration tests

EOF
    
    echo "✅ Phase architecture plan created: $PHASE_PLAN"
}
```

### Level 2: Wave Architecture Planning (After Wave Review)

When a REVIEW_WAVE_ARCHITECTURE passes, the Architect creates the next wave's architecture:

```bash
# ARCHITECT RESPONSIBILITY: Create Wave Architecture Plan
create_wave_architecture_plan() {
    local PHASE="$1"
    local WAVE="$2"
    local TIMESTAMP=$(date +%Y%m%d-%H%M%S)
    local WAVE_PLAN="$CLAUDE_PROJECT_DIR/phase-plans/phase${PHASE}/wave${WAVE}/WAVE-${PHASE}-${WAVE}-ARCHITECTURE-PLAN--${TIMESTAMP}.md"

    echo "═══════════════════════════════════════════════════════"
    echo "🌊 R210: Creating Wave ${WAVE} Architecture Plan"
    echo "═══════════════════════════════════════════════════════"

    # 0. Ensure planning directory exists
    mkdir -p "$CLAUDE_PROJECT_DIR/phase-plans/phase${PHASE}/wave${WAVE}"

    cat > "$WAVE_PLAN" << 'EOF'
# Phase ${PHASE} Wave ${WAVE} Architecture Plan

## 🎯 Wave Vision Alignment
[How this wave contributes to phase goals]

## 📊 Previous Wave Analysis
### Completed Implementations
- [What was built in previous waves]
- [APIs/contracts now available]

### Integration Points
- [How this wave connects to previous waves]
- [Shared libraries to use]

## 🏗️ Wave ${WAVE} Architecture

### Effort Parallelization Plan
```
Effort 1: API Contracts (MUST BE FIRST)
  └─→ Effort 2: Client Implementation (depends on 1)
  └─→ Effort 3: Server Implementation (depends on 1)
      
Effort 4: Independent Feature A (PARALLEL)
Effort 5: Independent Feature B (PARALLEL)
```

### API/Contract Definitions for This Wave
```go
// Wave-specific contracts
type WaveService interface {
    // Define all contracts UPFRONT
    StartProcess(ctx context.Context, config Config) error
    GetStatus(id string) (Status, error)
}
```

### Code Reuse Strategy
```go
// From Phase 1
import "pkg/phase1/common"

// From earlier waves in this phase
import "pkg/phase${PHASE}/wave1/contracts"

// Example usage
func NewWaveService() *Service {
    return &Service{
        validator: common.NewValidator(),     // Phase 1
        processor: contracts.NewProcessor(),  // Wave 1
    }
}
```

### Effort Dependencies
1. **Effort 1**: Define contracts
   - No dependencies
   - Creates: API interfaces, data models
   
2. **Effort 2**: Implement clients
   - Depends on: Effort 1 contracts
   - Creates: Client libraries
   
3. **Effort 3**: Implement servers
   - Depends on: Effort 1 contracts
   - Creates: Server implementations

## 📈 Implementation Priorities

### Order of Implementation
1. **APIs/Interfaces/Contracts** (Effort 1)
2. **Shared Libraries/Plugins** (Efforts 2-3)
3. **Private Implementations** (Efforts 4-5)
4. **Tests** (Effort 6)
5. **Documentation** (Effort 7)

### Critical Path
```
Contracts → Libraries → Implementations → Tests → Docs
    ↓           ↓              ↓            ↓       ↓
  Day 1      Day 2-3        Day 4-6      Day 7    Day 8
```

## ✅ Wave Success Criteria
- Contracts defined in first effort
- Parallel efforts don't block each other
- All efforts under 800 lines
- Code reuse achieved as planned

EOF
    
    echo "✅ Wave architecture plan created: $WAVE_PLAN"
}
```

## Architect State Machine Integration

### New Architect States

1. **PHASE_ARCHITECTURE_PLANNING**
   - Triggered after successful PHASE_ASSESSMENT
   - Creates $CLAUDE_PROJECT_DIR/phase-plans/phaseX/PHASE-X-ARCHITECTURE-PLAN--TIMESTAMP.md
   - Transitions to: ARCHITECTURE_COMPLETE

2. **WAVE_ARCHITECTURE_PLANNING**
   - Triggered after successful REVIEW_WAVE_ARCHITECTURE
   - Creates $CLAUDE_PROJECT_DIR/phase-plans/phaseX/waveY/WAVE-X-Y-ARCHITECTURE-PLAN--TIMESTAMP.md
   - Transitions to: ARCHITECTURE_COMPLETE

### State Flow
```
PHASE_ASSESSMENT (PASS)
    ↓
PHASE_ARCHITECTURE_PLANNING
    ↓
Create $CLAUDE_PROJECT_DIR/planning/phaseX/PHASE-X-ARCHITECTURE-PLAN--TIMESTAMP.md
Create $CLAUDE_PROJECT_DIR/planning/phaseX/PHASE-X-PLAN--TIMESTAMP.md
    ↓
Signal Orchestrator
    ↓
Orchestrator spawns Code Reviewer
    ↓
Code Reviewer: PHASE_IMPLEMENTATION_PLANNING
```

## Architecture Plan Requirements

### MUST Include:
1. **Vision Alignment**: How it fits master plan
2. **Previous Work Analysis**: What to reuse
3. **Contracts/APIs**: Define FIRST, implement later
4. **Abstractions/Interfaces**: Clear boundaries
5. **Parallelization Strategy**: What can run concurrently
6. **Code Reuse Map**: What comes from where
7. **MVP vs Nice-to-Have**: Priority clarity

### MUST Consider:
1. **Lessons Learned**: Adapt based on reality
2. **Technical Debt**: Address or document
3. **Integration Points**: How pieces connect
4. **Future Extensibility**: Hooks for next phases

## Validation Checklist

### Before Creating Architecture Plan:
- [ ] Previous phase/wave implementation reviewed
- [ ] Master plan requirements understood
- [ ] Integration points identified
- [ ] Parallelization opportunities analyzed

### Architecture Plan Must Have:
- [ ] All contracts defined upfront
- [ ] Clear parallelization strategy
- [ ] Code reuse explicitly mapped
- [ ] MVP features separated from nice-to-haves
- [ ] Pseudo-code examples provided

### After Architecture Plan:
- [ ] Code Reviewer can create implementation plan
- [ ] Parallelization strategy is clear
- [ ] Contracts enable independent work
- [ ] Vision alignment maintained

## Examples

### ✅ Good Architecture Plan (in $CLAUDE_PROJECT_DIR/planning/phaseX/)
```markdown
## APIs and Contracts (DEFINED FIRST!)
```go
// All contracts defined upfront
type UserService interface {
    Create(ctx context.Context, user User) (*User, error)
    Get(ctx context.Context, id string) (*User, error)
    Update(ctx context.Context, id string, user User) (*User, error)
    Delete(ctx context.Context, id string) error
}
```

## Parallelization Strategy
- Effort 1: Define UserService interface (MUST BE FIRST)
- Efforts 2-3: Can parallelize after Effort 1:
  - Effort 2: PostgreSQL implementation
  - Effort 3: MongoDB implementation
```

### ❌ Bad Architecture Plan
```markdown
## Implementation Details
We'll figure out the APIs as we go...
Each effort can define its own interfaces...
```

## Summary

- **Architect** creates architecture plans AFTER reviews pass
- **Architecture plans** define contracts, APIs, and parallelization
- **Code Reviewer** uses architecture to create implementation plans
- **Contracts first**, implementations second
- **Early phases** focus on MVP, later phases add nice-to-haves
- **Parallelization** enabled by upfront contract definition