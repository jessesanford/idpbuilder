# WAVE_REVIEW State Rules

## 🔴🔴🔴 CRITICAL: REPORT LOCATION REQUIREMENTS 🔴🔴🔴

**YOU MUST CREATE THE WAVE REVIEW REPORT IN THE EXACT LOCATION:**
```bash
# MANDATORY LOCATION:
Directory: wave-reviews/phase{N}/wave{W}/
Filename:  PHASE-{N}-WAVE-{W}-REVIEW-REPORT.md
Full path: wave-reviews/phase{N}/wave{W}/PHASE-{N}-WAVE-{W}-REVIEW-REPORT.md

# EXAMPLE for Phase 1, Wave 2:
mkdir -p wave-reviews/phase1/wave2
Write wave-reviews/phase1/wave2/PHASE-1-WAVE-2-REVIEW-REPORT.md

# WRONG LOCATIONS (WILL FAIL):
❌ ~/WAVE-1-2-REVIEW-REPORT.md                      # Root directory
❌ PHASE-1-WAVE-2-REVIEW-REPORT.md                  # Current directory
❌ ./WAVE-REVIEW.md                                 # Current directory
❌ reports/PHASE-1-WAVE-2-REVIEW-REPORT.md          # Wrong directory
❌ wave2/PHASE-1-WAVE-2-REVIEW-REPORT.md            # Missing parent directories
❌ phase1/wave2/PHASE-1-WAVE-2-REVIEW-REPORT.md     # Missing wave-reviews parent

# CORRECT LOCATION (ONLY THIS WORKS):
✅ wave-reviews/phase1/wave2/PHASE-1-WAVE-2-REVIEW-REPORT.md
```

**VERIFICATION FUNCTION - USE THIS:**
```bash
verify_wave_report_location() {
    local PHASE=$1
    local WAVE=$2
    local EXPECTED="wave-reviews/phase${PHASE}/wave${WAVE}/PHASE-${PHASE}-WAVE-${WAVE}-REVIEW-REPORT.md"
    
    if [[ ! -f "$EXPECTED" ]]; then
        echo "❌ CRITICAL ERROR: Report not in correct location!"
        echo "❌ Expected: $EXPECTED"
        echo "❌ Orchestrator will NOT find your report!"
        exit 1
    fi
    echo "✅ Report in correct location: $EXPECTED"
}

# ALWAYS verify after creating:
verify_wave_report_location 1 2  # For Phase 1, Wave 2
```

**PENALTY FOR WRONG LOCATION: -50% grading penalty, orchestrator cannot proceed**

## Core Wave Assessment Rules

---
### 🚨🚨🚨 RULE R297 - Architect Split Detection Protocol
**Source:** rule-library/R297-architect-split-detection-protocol.md
**Criticality:** BLOCKING - Must check splits BEFORE measuring integration

MANDATE: Check split_count in orchestrator-state.yaml BEFORE measuring any effort.
If split_count > 0, the effort was already split and is COMPLIANT.
Integration branches merge all splits (will exceed limits - EXPECTED).
Measure ORIGINAL effort branches, NOT integration branches.
PRs come from effort branches, NOT integration.
---

---
### 🚨🚨 RULE R074.0.0 - Wave Completion Architectural Review
**Source:** rule-library/RULE-REGISTRY.md#R074
**Criticality:** MANDATORY - Required for approval

MANDATE: Architect must review completed waves before
integration into phase. Assess effort integration quality,
architectural consistency, and readiness for wave integration.

ASSESSMENT CRITERIA:
- Effort size compliance (<800 lines per effort)
- KCP pattern consistency across all efforts
- API coherence within wave scope
- Integration test coverage and pass rates
- Performance impact on existing components
---

---
### ℹ️ RULE R075.0.0 - Wave Integration Readiness
**Source:** rule-library/RULE-REGISTRY.md#R075
**Criticality:** INFO - Best practice

MANDATE: Before wave integration approval, verify all efforts
are properly tested, documented, and maintain system
stability when combined.

INTEGRATION REQUIREMENTS:
- All efforts merge cleanly with wave integration branch
- No architectural conflicts between efforts
- Combined wave maintains KCP multi-tenancy guarantees
- Performance degradation within acceptable limits
- Security boundaries preserved across all efforts
---

---
### 🚨 RULE R076.0.0 - Effort Size Compliance Verification
**Source:** rule-library/RULE-REGISTRY.md#R076
**Criticality:** CRITICAL - Major impact on grading

MANDATE: Every effort in wave must be ≤800 lines as measured
by the line counter tool. Oversized efforts require split
before wave approval.

ENFORCEMENT:
- Use ONLY `$PROJECT_ROOT/tools/line-counter.sh` for measurement
- Navigate TO each effort directory before measuring
- Run with NO PARAMETERS (tool auto-detects branch)
- Exclude generated code (zz_generated*, *.pb.go, CRDs) - automatic
- Split efforts >800 lines before integration
- Verify all splits maintain functional coherence
- Document split rationale and dependencies
---

### 🚨🚨🚨 RULE R022 - Architect Size Verification Protocol
**Source:** rule-library/R022-architect-size-verification.md
**Criticality:** BLOCKING - Measurement errors = invalid review

**CRITICAL: How Architects MUST Measure Sizes:**
```bash
# For EACH effort in the wave:
cd /path/to/efforts/phase2/wave2/effort-name

# Find project root (where orchestrator-state.yaml lives):
PROJECT_ROOT=$(pwd)
while [ "$PROJECT_ROOT" != "/" ]; do
    [ -f "$PROJECT_ROOT/orchestrator-state.yaml" ] && break
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
done

# Run line counter with NO PARAMETERS:
$PROJECT_ROOT/tools/line-counter.sh  # NO PARAMETERS!
```

**❌ NEVER DO THIS:**
- `./tools/line-counter.sh phase2/wave2/effort-name` (NO parameters!)
- `$CLAUDE_PROJECT_DIR/tools/line-counter.sh` (often undefined!)
- Manual counting or alternative tools (ONLY use line-counter.sh per R304!)
---

---
### ℹ️ RULE R037.0.0 - KCP Resource Pattern Enforcement
**Source:** rule-library/RULE-REGISTRY.md#R037
**Criticality:** INFO - Best practice

MANDATE: All wave components must maintain KCP multi-tenant
patterns with workspace isolation and logical cluster support.

WAVE-LEVEL VALIDATION:
- Consistent LogicalCluster field usage across efforts
- Workspace-aware resource controllers
- Proper RBAC scoping for multi-tenancy
- Event and status propagation within workspace boundaries
- Cross-effort API compatibility
---

## State-Specific Context

### 🚨🚨🚨 RULE R258 - Mandatory Wave Review Report [BLOCKING]
**Source:** rule-library/R258-mandatory-wave-review-report.md
**Criticality:** BLOCKING - No wave progression without report

**MANDATE**: Architect MUST create a permanent wave review report file before signaling review complete.

**REQUIRED FILE**: `wave-reviews/phase{N}/wave{W}/PHASE-{N}-WAVE-{W}-REVIEW-REPORT.md`
⚠️⚠️⚠️ **SEE CRITICAL LOCATION SECTION AT TOP OF THIS FILE** ⚠️⚠️⚠️

**🔴 STEP-BY-STEP INSTRUCTIONS:**
```bash
# Step 1: Get phase and wave numbers
PHASE=1  # From orchestrator-state.yaml
WAVE=2   # From orchestrator-state.yaml

# Step 2: CREATE DIRECTORY STRUCTURE (MANDATORY!)
mkdir -p wave-reviews/phase${PHASE}/wave${WAVE}

# Step 3: Set exact report path
REPORT_FILE="wave-reviews/phase${PHASE}/wave${WAVE}/PHASE-${PHASE}-WAVE-${WAVE}-REVIEW-REPORT.md"

# Step 4: Create report with Write tool using EXACT path
# Example: Write wave-reviews/phase1/wave2/PHASE-1-WAVE-2-REVIEW-REPORT.md

# Step 5: Verify location (CRITICAL!)
ls -la wave-reviews/phase${PHASE}/wave${WAVE}/
[[ -f "$REPORT_FILE" ]] || { echo "❌ FAILED!"; exit 1; }

# Step 6: Commit and push
git add "$REPORT_FILE"
git commit -m "review: Phase ${PHASE} Wave ${WAVE} review - ${DECISION}"
git push
```

**DECISION TYPES**:
- PROCEED_NEXT_WAVE: Wave approved, start next wave
- PROCEED_PHASE_ASSESSMENT: Last wave complete, trigger phase assessment  
- CHANGES_REQUIRED: Fixes needed before progression
- WAVE_FAILED: Major issues, cannot proceed

**🚨🚨🚨 FORBIDDEN DECISIONS - NEVER ALLOWED IN WAVE_REVIEW 🚨🚨🚨**:
- ❌ END_PHASE: The architect CANNOT end phases
- ❌ PROJECT_COMPLETE: The architect CANNOT end the project
- ❌ END_PROJECT: The architect CANNOT terminate the project
- ❌ SKIP_PHASE: The architect CANNOT skip phases

**CRITICAL REMINDERS**:
- You are reviewing a WAVE, not a PHASE
- You CANNOT make project-level decisions
- You CANNOT skip or end phases
- Only the ORCHESTRATOR decides when to end the project
- You CANNOT signal wave review complete without creating this report!
---

### WAVE_REVIEW State Purpose
This state is entered when all efforts in a wave are complete and need architectural review before wave integration. The architect validates effort cohesion, size compliance, and integration readiness.

### State Transitions

**ENTRY CONDITIONS:**
- All efforts in wave marked as `EFFORT_COMPLETE`
- Wave integration branch created by orchestrator
- No outstanding code review issues (all efforts approved)

**EXIT CONDITIONS:**
- **PROCEED_NEXT_WAVE**: Wave meets all architectural standards → Approve for next wave
- **PROCEED_PHASE_ASSESSMENT**: Last wave complete and approved → Trigger phase assessment
- **CHANGES_REQUIRED**: Issues identified requiring fixes → Return to orchestrator with specific actions
- **WAVE_FAILED**: Major architectural problems → Halt wave progression, require redesign

**⚠️⚠️⚠️ ARCHITECT AUTHORITY LIMITS ⚠️⚠️⚠️**:
- The architect reviews and recommends, but CANNOT control project flow
- The architect CANNOT end phases or the project
- The architect CANNOT skip any phases, even if they seem complete
- All phases in the plan MUST be executed - no exceptions
- Only the ORCHESTRATOR has authority to end the project

**STATE TRANSITION FLOW:**
```
WAVE_REVIEW → [Assessment Complete] → Decision State
├─ PROCEED → Wave Integration Approved
├─ CHANGES_REQUIRED → Fix Cycle (orchestrator coordinates)
└─ STOP → Wave Redesign Required
```

## Wave Assessment Criteria

### Primary Assessment Areas

| Area | Weight | Critical Threshold | Pass Criteria |
|------|--------|-------------------|---------------|
| Size Compliance | 25% | 100% efforts ≤800 lines | All pass |
| KCP Pattern Consistency | 25% | No pattern violations | All pass |
| Integration Stability | 20% | Clean merge + tests pass | 95% success |
| API Coherence | 15% | No conflicting interfaces | All resolved |
| Performance Impact | 10% | <10% degradation | Within limits |
| Documentation Quality | 5% | Complete work logs | 90% coverage |

### Size Compliance Assessment

**CRITICAL REQUIREMENT**: Every effort must be ≤800 lines (R297 MUST be checked FIRST!)
- **Split Detection**: Check split_count BEFORE measuring (R297)
- **Measurement Tool**: `$PROJECT_ROOT/tools/line-counter.sh` (find project root first!)
- **Measurement Target**: ORIGINAL effort branches, NOT integration branches
- **Exclusions**: Generated code (zz_generated*.go, *.pb.go, CRDs, SDK clients) - automatic
- **Action on Violation**: Immediate STOP, require effort splitting

**Size Assessment Process** (Per R297 and R022):
1. **CHECK SPLIT_COUNT FIRST** (R297):
   - Read `split_count` from orchestrator-state.yaml
   - If > 0: Mark as COMPLIANT, skip size measurement
   - If 0: Continue to measure original branch
2. Navigate TO each effort directory (cd /efforts/phase*/wave*/[effort-name])
3. Find project root: `while [ "$PWD" != "/" ]; do [ -f orchestrator-state.yaml ] && break; cd ..; done`
4. Run line counter with NO PARAMETERS: `$PROJECT_ROOT/tools/line-counter.sh`
5. Document results in assessment report  
6. Flag any efforts >800 lines for splitting
7. Verify splits maintain logical coherence
8. **NEVER measure integration branch for compliance** (integration merges all splits)

### KCP Pattern Consistency Assessment

**CONSISTENCY REQUIREMENTS**:
- Uniform LogicalCluster field usage across all efforts
- Consistent workspace isolation patterns
- Compatible RBAC models between efforts
- Aligned event handling and status reporting

**Pattern Validation Checklist**:
- [ ] All CRDs include LogicalCluster metadata field
- [ ] Controllers use workspace-scoped clients consistently
- [ ] Resource indexing includes workspace context
- [ ] Cross-effort APIs don't violate workspace boundaries
- [ ] Event propagation respects multi-tenancy

## Integration Readiness Validation

### Pre-Integration Requirements

1. **Clean Merge Validation**
   - All effort branches merge without conflicts
   - No duplicate resource definitions
   - Consistent dependency versions

2. **Integration Testing**
   - Cross-effort integration tests pass
   - Performance regression tests complete
   - End-to-end workflows functional

3. **API Compatibility**
   - No breaking changes between efforts
   - Consistent API patterns and naming
   - Backward compatibility maintained

### Wave Integration Approval Criteria

**PROCEED CONDITIONS** (All must be true):
- ✅ All efforts size-compliant (R297: check split_count first, then measure original branches ≤800 lines)
- ✅ All KCP patterns consistent across efforts
- ✅ Wave integration branch merges cleanly
- ✅ Integration tests pass with >95% success rate
- ✅ Performance impact <10% degradation
- ✅ No security boundary violations
- ✅ Split efforts recognized as compliant (even if integration exceeds limit)

**CHANGES_REQUIRED CONDITIONS**:
- Minor API inconsistencies between efforts
- Integration test failures <5% (fixable issues)
- Performance regression 10-20% (optimization needed)
- Documentation gaps in work logs

**STOP CONDITIONS**:
- Any effort >800 lines (splitting required)
- KCP pattern violations (multi-tenancy broken)
- Major integration conflicts (architectural mismatch)
- Performance regression >20% (design issue)
- Security vulnerabilities introduced

## Assessment Documentation Requirements

### 🚨🚨🚨 MANDATORY Wave Review Report (R258)

**YOU MUST CREATE THIS FILE BEFORE SIGNALING COMPLETE:**
- **File**: `wave-reviews/phase{N}/wave{W}/PHASE-{N}-WAVE-{W}-REVIEW-REPORT.md`
- **Example**: `wave-reviews/phase1/wave2/PHASE-1-WAVE-2-REVIEW-REPORT.md`
- **When**: BEFORE returning your review decision
- **Location**: ⚠️ MUST be in exact directory structure - see top of file ⚠️

**NEVER CREATE IN:**
- ❌ Root directory (~/WAVE-REVIEW.md)
- ❌ Current directory (./PHASE-1-WAVE-2-REVIEW-REPORT.md)
- ❌ Wrong structure (phase1/wave2/report.md)
- **Decision**: Must be one of:
  - PROCEED_NEXT_WAVE
  - PROCEED_PHASE_ASSESSMENT
  - CHANGES_REQUIRED
  - WAVE_FAILED

### Required Wave Review Deliverables

1. **Wave Review Report (MANDATORY per R258)**
   - Complete standardized report template
   - Size compliance verification for each effort
   - KCP pattern consistency evaluation
   - Integration stability assessment
   - Performance impact analysis
   - Architectural scoring (0-100 for each category)
   - Decision with rationale
   - Architect sign-off with timestamp

2. **Integration Readiness Certification**
   - Merge compatibility confirmation
   - Test result summary
   - API coherence validation
   - Security boundary verification

3. **Decision Record**
   - One of four valid decisions (per R258)
   - Rationale for decision
   - Specific action items (if CHANGES_REQUIRED)
   - Architecture approval signature with timestamp

### Quality Gates

**MANDATORY QUALITY GATES**:
- Size compliance: 100% efforts ≤800 lines
- Pattern compliance: 100% KCP pattern adherence
- Integration tests: >95% pass rate
- Security: Zero critical/high vulnerabilities

**PERFORMANCE GATES**:
- API response time: <10% increase from baseline
- Memory usage: <15% increase from baseline  
- Throughput: >90% of baseline performance
- Resource utilization: Within acceptable limits

### Review Timeline

- **Standard Wave Review**: 4-8 hours
- **Complex Wave Review**: 1-2 days
- **Problem Wave Review**: 2-5 days (includes fix cycles)

**Review Scope Factors**:
- Number of efforts in wave (typically 3-7)
- Complexity of inter-effort integration
- Performance testing requirements
