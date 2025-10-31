# Phase N Implementation Plan

**Phase**: [Phase Number/Name]
**Created**: [Date]
**Planner**: [Code Reviewer Agent Name]
**Fidelity Level**: **WAVE LIST ONLY** (high-level descriptions, no detailed plans)

---

## Phase Overview

**Goal**: [1-2 sentence summary of what this phase achieves]

**Scope**: [Brief description of functional boundaries]

**Dependencies**:
- [Previous phase or external system]
- [Infrastructure requirements]

---

## Wave List

This section provides a **high-level roadmap** of waves in this phase. Detailed effort definitions will be created **just-in-time** during wave planning.

### Wave 1: [Wave Name]

**Description**: [1-2 sentences describing the main goal of this wave]

**Key Features**:
- [Feature 1]
- [Feature 2]
- [Feature 3]

**Efforts**: **TBD** (Will be defined during Wave 1 planning)

**Estimated Complexity**: [Low/Medium/High]

---

### Wave 2: [Wave Name]

**Description**: [1-2 sentences describing the main goal of this wave]

**Key Features**:
- [Feature 1]
- [Feature 2]
- [Feature 3]

**Efforts**: **TBD** (Will be defined during Wave 2 planning)

**Estimated Complexity**: [Low/Medium/High]

---

### Wave 3: [Wave Name]

**Description**: [1-2 sentences describing the main goal of this wave]

**Key Features**:
- [Feature 1]
- [Feature 2]
- [Feature 3]

**Efforts**: **TBD** (Will be defined during Wave 3 planning)

**Estimated Complexity**: [Low/Medium/High]

---

### Wave 4: [Wave Name] (if applicable)

**Description**: [1-2 sentences describing the main goal of this wave]

**Key Features**:
- [Feature 1]
- [Feature 2]

**Efforts**: **TBD** (Will be defined during Wave 4 planning)

**Estimated Complexity**: [Low/Medium/High]

---

## Wave Dependencies

```
Wave 1 (Foundation)
  ↓
Wave 2 (Core Features)
  ↓
Wave 3 (Advanced Features)
  ↓
Wave 4 (Refinement) [optional]
```

**Critical Path**: [Which waves must complete sequentially vs. which can parallelize]

---

## Risk Assessment

### High-Risk Waves

- **Wave [N]**: [Risk description and mitigation strategy]
- **Wave [N]**: [Risk description and mitigation strategy]

### Dependencies on External Systems

- [External system 1]: Required for [Wave N]
- [External system 2]: Required for [Wave N]

---

## Success Criteria

**Phase Complete When**:
- [ ] All waves integrated and tested
- [ ] Phase integration tests passing
- [ ] Documentation complete
- [ ] Architect phase assessment approved

---

## Progressive Planning Notes

### Why Wave List Only?

This phase implementation plan provides **only a wave list** because:

1. **Adaptive Planning**: Detailed effort plans must adapt based on what we learn during implementation
2. **Fidelity Gradient**: Effort definitions require **real code examples** from completed waves
3. **Just-In-Time Planning**: Creating detailed plans upfront wastes effort when requirements change
4. **Progressive Refinement**: Each wave's planning benefits from lessons learned in previous waves

### When Are Effort Plans Created?

**Effort plans** are created **just-in-time** at **wave start**:

1. Orchestrator transitions to **WAVE_START**
2. Architect creates **wave architecture plan** (with real code examples)
3. Code Reviewer creates **wave implementation plan** (with detailed effort definitions)
4. Efforts include:
   - Exact file lists
   - Real code specifications
   - R213 metadata
   - Dependencies

---

## Next Steps

1. **Phase Architecture Planning**: Architect creates `PHASE-N-ARCHITECTURE.md` with pseudocode patterns
2. **Wave 1 Start**: When phase approved, orchestrator proceeds to Wave 1
3. **Wave 1 Architecture**: Architect creates `WAVE-1-ARCHITECTURE.md` with real code
4. **Wave 1 Implementation**: Code Reviewer creates `WAVE-1-IMPLEMENTATION.md` with exact effort specs

**Note**: This document is intentionally **high-level**. Detailed planning happens progressively during wave execution.
