# PHASE_ASSESSMENT State Rules

## 🔴🔴🔴 CRITICAL: ARCHITECT ROLE LIMITATIONS 🔴🔴🔴

**THE ARCHITECT IS AN ASSESSOR, NOT A DECIDER:**
- ✅ You assess phase quality and compliance
- ✅ You recommend proceeding or requesting changes
- ❌ You CANNOT end the project
- ❌ You CANNOT skip phases
- ❌ You CANNOT decide that "the MVP is complete"
- ❌ You CANNOT terminate Phase 2 or any other phase

**MANDATORY PRINCIPLE**: Every phase in the plan MUST be executed and assessed. Even if Phase 1 delivers a working system, Phase 2 and all subsequent phases MUST still be implemented.

## 🚨🚨🚨 MANDATORY PHASE ASSESSMENT REPORT - BLOCKING REQUIREMENT 🚨🚨🚨

### 🚨🚨🚨 RULE R257.0.0 - Mandatory Phase Assessment Report [BLOCKING]
**Source:** rule-library/R257-mandatory-phase-assessment-report.md
**Criticality:** BLOCKING - Phase cannot complete without this

**YOU MUST CREATE A PERMANENT ASSESSMENT REPORT FILE:**
- **File Name**: `PHASE-{N}-ASSESSMENT-REPORT.md`
- **Location**: `phase-assessments/phase{N}/`
- **When**: BEFORE signaling assessment complete
- **Verification**: Orchestrator will verify file exists

**CRITICAL STEPS:**
1. Complete your architectural assessment
2. CREATE the assessment report file with ALL mandatory sections
3. Include your DECISION in the report (PROCEED_NEXT_PHASE|CHANGES_REQUIRED|PHASE_FAILED)
4. Calculate and document all scoring metrics
5. Sign off with timestamp
6. Commit and push the report file
7. ONLY THEN signal assessment complete to orchestrator

**EXAMPLE:**
```bash
# After completing assessment
PHASE=1
DECISION="PROCEED_NEXT_PHASE"
REPORT_DIR="phase-assessments/phase${PHASE}"
REPORT_FILE="${REPORT_DIR}/PHASE-${PHASE}-ASSESSMENT-REPORT.md"

mkdir -p "$REPORT_DIR"
# Create report with all mandatory sections...
git add "$REPORT_FILE"
git commit -m "assessment: Phase ${PHASE} assessment report - ${DECISION}"
git push
```

**VIOLATIONS:**
- ❌ Providing verbal assessment without report = BLOCKING FAILURE
- ❌ Missing mandatory sections = INVALID ASSESSMENT
- ❌ Not committing report = LOST ASSESSMENT
- ❌ Orchestrator cannot proceed without verified report

---

## Core Architecture Assessment Rules

---
### ℹ️ RULE R071.0.0 - Phase Architectural Integrity Assessment
**Source:** rule-library/RULE-REGISTRY.md#R071
**Criticality:** INFO - Best practice

MANDATE: Architect must assess architectural integrity before
phase completion. Evaluate structural soundness, pattern
compliance, and integration readiness of all phase components.

CRITERIA:
- KCP multi-tenancy pattern consistency
- API design stability and backwards compatibility
- Resource hierarchy and namespace isolation
- Controller pattern adherence
- Performance implications at enterprise scale
---

---
### ℹ️ RULE R072.0.0 - KCP Pattern Compliance Verification
**Source:** rule-library/RULE-REGISTRY.md#R072
**Criticality:** INFO - Best practice

MANDATE: Verify all phase implementations adhere to KCP
architectural patterns including logical clusters, workspace
isolation, and multi-tenancy constraints.

VERIFICATION POINTS:
- LogicalCluster field presence in all CRDs
- Workspace-aware controllers and indexing
- ClusterRole vs Role usage patterns
- Cross-workspace reference handling
- Resource quota and limit enforcement
---

---
### ℹ️ RULE R073.0.0 - Phase Completion Prerequisites
**Source:** rule-library/RULE-REGISTRY.md#R073
**Criticality:** INFO - Best practice

MANDATE: Before approving phase completion, verify all waves
are integrated, tested, and meet architectural standards.

PREREQUISITES:
- All wave integration branches merged successfully
- No architectural debt or anti-patterns introduced
- API compatibility maintained with previous phases
- Performance benchmarks meet enterprise requirements
- Security posture maintained or improved
---

---
### ℹ️ RULE R037.0.0 - KCP Resource Pattern Enforcement
**Source:** rule-library/RULE-REGISTRY.md#R037
**Criticality:** INFO - Best practice

MANDATE: All Kubernetes resources must follow KCP multi-tenant
patterns with proper workspace isolation and logical cluster
field usage.

ENFORCEMENT AREAS:
- CRD schema includes LogicalCluster metadata
- Controllers use workspace-scoped clients
- RBAC follows least-privilege workspace boundaries
- Resource names avoid cross-workspace conflicts
- Event handling respects workspace isolation
---

## State-Specific Context

### PHASE_ASSESSMENT State Purpose
This state is entered when a complete phase (containing multiple waves) needs architectural review before transitioning to the next phase. The architect evaluates system-wide consistency, integration quality, and readiness for next phase work.

### State Transitions

**ENTRY CONDITIONS:**
- All waves in current phase marked as `WAVE_COMPLETE`
- Phase integration branch created and tested
- Orchestrator requests phase assessment

**EXIT CONDITIONS:**
- **PROCEED_NEXT_PHASE**: Phase meets all architectural standards → Continue to next phase
- **CHANGES_REQUIRED**: Minor issues identified → Return specific fixes to orchestrator
- **PHASE_FAILED**: Major architectural problems → Halt progression, require redesign

**🚨🚨🚨 FORBIDDEN DECISIONS - NEVER ALLOWED IN PHASE_ASSESSMENT 🚨🚨🚨**:
- ❌ END_PHASE: The architect CANNOT end phases
- ❌ PROJECT_COMPLETE: The architect CANNOT end the project
- ❌ END_PROJECT: The architect CANNOT terminate the project
- ❌ SKIP_PHASE: The architect CANNOT skip phases
- ❌ SKIP_TO_FINAL: The architect CANNOT skip to final integration

**⚠️⚠️⚠️ ARCHITECT AUTHORITY LIMITS ⚠️⚠️⚠️**:
- The architect assesses phases but CANNOT control project flow
- The architect CANNOT end the project or skip phases
- ALL phases in the plan MUST be executed and assessed
- Even if Phase 1 seems to "complete the MVP", Phase 2 MUST still be executed
- Only the ORCHESTRATOR decides when the project is complete

**STATE TRANSITION FLOW:**
```
PHASE_ASSESSMENT → [Assessment Complete] → Decision State
├─ PROCEED_NEXT_PHASE → Next Phase Planning
├─ CHANGES_REQUIRED → Fix Cycle (return to orchestrator)
└─ PHASE_FAILED → Architecture Redesign Required
```

## Assessment Criteria Matrix

| Area | Critical | Important | Nice-to-Have |
|------|----------|-----------|--------------|
| KCP Pattern Compliance | 100% | - | - |
| API Backwards Compatibility | 100% | - | - |
| Multi-tenancy Isolation | 100% | - | - |
| Performance at Scale | - | 95% | - |
| Code Quality Metrics | - | 90% | 100% |
| Documentation Coverage | - | 85% | 95% |

## Architecture Decision Requirements

### MUST ASSESS:
1. **Structural Integrity**: Are all components properly integrated?
2. **Pattern Consistency**: Do implementations follow KCP patterns uniformly?
3. **Scalability**: Will the phase handle enterprise-scale workloads?
4. **Security Posture**: Are multi-tenancy boundaries secure?
5. **API Stability**: Are APIs ready for production consumption?

### MUST DOCUMENT (IN THE MANDATORY ASSESSMENT REPORT FILE):
- **PRIMARY**: Create `PHASE-{N}-ASSESSMENT-REPORT.md` per R257
- Architecture decision record for phase approval/rejection
- Specific issues requiring correction (if CHANGES_REQUIRED)
- Performance benchmark results and analysis
- Security assessment findings
- Integration test results summary
- All documentation MUST be in the assessment report file
- Verbal/inline responses are NOT sufficient

### CRITICAL FAILURE CONDITIONS:
- **IMMEDIATE PHASE_FAILED**: Security vulnerabilities, data leakage between workspaces
- **IMMEDIATE PHASE_FAILED**: Breaking changes to public APIs without migration path
