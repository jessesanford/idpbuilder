# WAVE_REVIEW Grading Criteria

## Wave Architecture Quality Metrics

### Primary Grading Categories

| Category | Weight | Scoring Method | Pass Threshold | Critical Failure |
|----------|--------|----------------|----------------|------------------|
| Size Compliance | 25% | Binary Pass/Fail | 100% | Any effort >800 lines |
| KCP Pattern Consistency | 25% | Pattern Checklist | 95% | Pattern violations |
| Integration Stability | 20% | Test Results + Manual | 90% | Major conflicts |
| API Coherence | 15% | Compatibility Matrix | 85% | Breaking conflicts |
| Performance Impact | 10% | Benchmark Comparison | 90% baseline | >20% degradation |
| Documentation Quality | 5% | Coverage Analysis | 80% | Missing work logs |

## Detailed Grading Rubric

### Size Compliance (25% - Binary Scoring)

| Effort Size Status | Points | Action Required |
|-------------------|--------|-----------------|
| All efforts ≤800 lines | 100 | None - proceed |
| 1-2 efforts >800 lines | 0 | Split oversized efforts |
| 3+ efforts >800 lines | 0 | Major wave restructuring |

**Measurement Protocol**:
```bash
# For each effort in wave (tool auto-detects correct base)
$PROJECT_ROOT/tools/line-counter.sh effort-branch-name

# Document results
effort_1: 650 lines ✅
effort_2: 780 lines ✅ 
effort_3: 820 lines ❌ (requires split)
```

**Critical Failure**: ANY effort >800 lines = Automatic STOP decision

### KCP Pattern Consistency (25% - Pattern Checklist)

| Pattern Area | Max Points | Deduction Rules |
|-------------|-----------|-----------------|
| LogicalCluster Fields | 30 | -10 per missing field |
| Workspace Isolation | 25 | -25 if boundary violated |
| Controller Patterns | 20 | -5 per inconsistency |
| RBAC Consistency | 15 | -15 if roles conflict |
| Event Scoping | 10 | -2 per scoping error |

**Pattern Assessment Matrix**:

| Effort | LogicalCluster | Workspace Isolation | Controller Pattern | RBAC | Events | Score |
|--------|---------------|-------------------|------------------|------|--------|-------|
| Effort A | ✅ (30/30) | ✅ (25/25) | ✅ (20/20) | ✅ (15/15) | ✅ (10/10) | 100% |
| Effort B | ✅ (30/30) | ✅ (25/25) | ⚠️ (15/20) | ✅ (15/15) | ✅ (10/10) | 95% |
| Effort C | ❌ (20/30) | ✅ (25/25) | ✅ (20/20) | ⚠️ (10/15) | ✅ (10/10) | 85% |

**Wave Pattern Score**: Average of all effort pattern scores

### Integration Stability (20% - Test Results + Manual Review)

| Integration Aspect | Points | Measurement Method | Threshold |
|-------------------|--------|--------------------|-----------|
| Merge Compatibility | 50 | Git merge success | All efforts merge cleanly |
| Integration Tests | 30 | Test suite results | >95% pass rate |
| Cross-Effort APIs | 20 | Manual validation | No conflicts detected |

**Integration Scoring Examples**:

#### Perfect Integration (100/100)
- All effort branches merge without conflicts: **50/50 points**
- Integration test suite: 47/47 tests pass (100%): **30/30 points**  
- No API conflicts found in manual review: **20/20 points**

#### Good Integration (85/100)
- 2 minor merge conflicts resolved: **45/50 points** (-5 for conflicts)
- Integration test suite: 45/47 tests pass (96%): **30/30 points**
- One minor API naming inconsistency: **15/20 points** (-5 for inconsistency)

#### Poor Integration (60/100)  
- Multiple merge conflicts, some unresolved: **30/50 points** (-20 for major conflicts)
- Integration test suite: 42/47 tests pass (89%): **25/30 points** (-5 below threshold)
- API conflicts require design changes: **5/20 points** (-15 for major conflicts)

### API Coherence (15% - Compatibility Matrix)

| Coherence Factor | Weight | Assessment Method | Scoring |
|-----------------|--------|-------------------|---------|
| Naming Consistency | 40% | Pattern analysis | 0-100% based on violations |
| Interface Compatibility | 35% | Cross-effort API calls | 0-100% based on conflicts |
| Schema Alignment | 25% | CRD field analysis | 0-100% based on inconsistencies |

**API Coherence Calculation**:
```
Naming Score = 100 - (violations × 10)      # Max deduction: 100%
Interface Score = 100 - (conflicts × 25)    # Max deduction: 100% 
Schema Score = 100 - (misalignments × 15)   # Max deduction: 100%

API Coherence = (Naming × 0.40) + (Interface × 0.35) + (Schema × 0.25)
```

### Performance Impact (10% - Benchmark Comparison)

| Metric | Weight | Baseline | Acceptable | Warning | Critical |
|--------|--------|----------|------------|---------|----------|
| API Response Time | 30% | 100ms | <110ms | 110-120ms | >120ms |
| Memory Usage | 25% | 512MB | <575MB | 575-640MB | >640MB |
| CPU Usage | 25% | 2 cores | <2.2 cores | 2.2-2.4 cores | >2.4 cores |
| Throughput | 20% | 1000 req/s | >900 req/s | 900-800 req/s | <800 req/s |

**Performance Scoring**:
- **Acceptable**: 100% score
- **Warning**: 75% score  
- **Critical**: 0% score (triggers STOP)

### Documentation Quality (5% - Coverage Analysis)

| Documentation Type | Weight | Requirement | Scoring |
|-------------------|--------|-------------|---------|
| Work Log Completeness | 60% | All efforts have work logs | 100% or 0% |
| Implementation Notes | 25% | Technical decisions documented | 0-100% coverage |
| Integration Guide | 15% | Wave integration documented | 0-100% completeness |

## Overall Wave Grade Calculation

### Final Score Formula
```
Wave Score = (Size_Compliance × 0.25) + (KCP_Patterns × 0.25) + 
             (Integration_Stability × 0.20) + (API_Coherence × 0.15) + 
             (Performance × 0.10) + (Documentation × 0.05)
```

### Decision Matrix

| Final Score | Decision | Next Action | Timeline |
|-------------|----------|-------------|----------|
| 95-100% | PROCEED | Approve wave integration | Immediate |
| 85-94% | CHANGES_REQUIRED | Fix specific issues | 1-3 days |
| 70-84% | CHANGES_REQUIRED | Major corrections needed | 5-10 days |
| <70% | STOP | Redesign wave architecture | 2-4 weeks |

### Grade Calculation Examples

#### Example 1: Excellent Wave (PROCEED)
```
Size Compliance: 100% (all efforts ≤800 lines)
KCP Patterns: 98% (minor controller inconsistency in one effort)
Integration: 95% (clean merges, all tests pass)
API Coherence: 92% (one naming inconsistency)
Performance: 95% (slight memory increase)
Documentation: 90% (complete work logs)

Final = (100×0.25) + (98×0.25) + (95×0.20) + (92×0.15) + (95×0.10) + (90×0.05)
Final = 25 + 24.5 + 19 + 13.8 + 9.5 + 4.5 = 96.3%
Decision = PROCEED
```

#### Example 2: Needs Correction
```
Size Compliance: 100% (all efforts ≤800 lines)  
KCP Patterns: 92% (LogicalCluster field missing in one CRD)
Integration: 88% (minor merge conflicts resolved)
API Coherence: 85% (interface compatibility issues)
Performance: 82% (performance regression needs optimization)
Documentation: 75% (incomplete implementation notes)

Final = (100×0.25) + (92×0.25) + (88×0.20) + (85×0.15) + (82×0.10) + (75×0.05)
Final = 25 + 23 + 17.6 + 12.75 + 8.2 + 3.75 = 90.3%
Decision = CHANGES_REQUIRED
```

#### Example 3: Stop Required
```
Size Compliance: 0% (2 efforts >800 lines - auto-fail)
KCP Patterns: 75% (multiple pattern violations)
Integration: 70% (major merge conflicts)
API Coherence: 60% (breaking API changes between efforts)
Performance: 45% (critical performance regression)
Documentation: 50% (missing work logs)

Final = (0×0.25) + (75×0.25) + (70×0.20) + (60×0.15) + (45×0.10) + (50×0.05)
Final = 0 + 18.75 + 14 + 9 + 4.5 + 2.5 = 48.75%
Decision = STOP
```

## Action Item Generation

### CHANGES_REQUIRED Action Items

Based on grade breakdown, generate specific, actionable items:

**Size Violations**:
- "Split effort-webhook-framework into 2 efforts: webhook-core (400 lines) and webhook-validation (350 lines)"
- "Estimated time: 8 hours development + 4 hours testing"

**KCP Pattern Issues**:
- "Add LogicalCluster field to ConfigMap CRD schema"
- "Update webhook controller to use workspace-scoped client"
- "Estimated time: 3 hours development + 2 hours testing"

**Integration Problems**:
- "Resolve API naming conflict between webhook and config efforts"
- "Fix integration test failure in cross-effort authentication flow"
- "Estimated time: 6 hours development + 4 hours testing"

**Performance Issues**:
- "Optimize webhook validation logic causing 15% response time increase"
- "Investigate memory leak in config controller affecting baseline"
- "Estimated time: 12 hours investigation + optimization"

### Quality Assurance Requirements

**Assessment Documentation**:
- All scoring decisions must be supported by evidence
- Benchmark results preserved for audit
- Integration test logs archived
- Manual review findings documented

**Review Reproducibility**:
- Assessment criteria applied consistently
- Measurement tools and versions documented
- Review timeline and effort distribution tracked
- Decision rationale preserved for future reference