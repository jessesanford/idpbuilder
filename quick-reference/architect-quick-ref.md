# 🏗️ ARCHITECT QUICK REFERENCE

## 🚨 STARTUP CHECKLIST
```
□ Print: AGENT STARTUP: $(date '+%Y-%m-%d %H:%M:%S %Z')
□ Verify: Correct directory and git branch
□ Load: agent-states/architect/{STATE}/rules.md  
□ Read: orchestrator-state.yaml (check wave completion)
□ Read: WAVE-COMPLETION-ARCHITECT-REVIEW-PROTOCOL.md
```

## ⚡ CRITICAL METRICS (GRADING)
```
┌─────────────────────────────────────────────┐
│ 🚨 DECISION ACCURACY: No reversed decisions │
│ 🚨 ISSUE DETECTION: Catch critical problems │
│ 🚨 TRAJECTORY: Correct ON/OFF_TRACK assess  │
│ 🚨 ADDENDUM CLARITY: Next wave succeeds     │
└─────────────────────────────────────────────┘
```

## 🔄 STATE MACHINE FLOW
```
INIT → WAVE_REVIEW → [PROCEED|PROCEED_WITH_CHANGES|STOP]
  ↓
PHASE_ASSESSMENT → [ON_TRACK|NEEDS_CORRECTION|OFF_TRACK]
  ↓  
INTEGRATION_REVIEW → [CLEAN_MERGE|CONFLICTS|ARCHITECTURAL_ISSUES]
```

## 🎯 KEY DECISIONS BY STATE

| State | Decision Types | Critical Impact |
|-------|---------------|----------------|
| `WAVE_REVIEW` | PROCEED/STOP | Wave continuation |
| `PHASE_ASSESSMENT` | ON_TRACK/OFF_TRACK | Phase direction |
| `INTEGRATION_REVIEW` | MERGE/BLOCK | Code integration |

## 🛠️ ESSENTIAL COMMANDS

### Wave Completion Check
```bash
# Verify all splits under limit
for effort in $(grep "efforts_completed:" orchestrator-state.yaml -A 20); do
    if [[ $effort =~ branch:.*$ ]]; then
        branch=${effort#*: }
        lines=$(line-counter.sh -c $branch | grep "Total:" | awk '{print $2}')
        echo "Effort $branch: $lines/800 lines"
        if [ $lines -gt 800 ]; then
            echo "🚨 VIOLATION: $branch exceeds limit"
        fi
    fi
done
```

### Integration Branch Assessment
```bash
# Check merge conflicts
git checkout wave-integration-branch
git merge --no-commit --no-ff effort-branch-1
git merge --no-commit --no-ff effort-branch-2
# Assess conflicts and architectural coherence
```

## 🔍 WAVE REVIEW PROTOCOL

### Size Compliance Verification
```
FOR EACH completed effort:
  ✅ Measure: line-counter.sh -c ${EFFORT_BRANCH}
  ✅ Verify: <800 lines total  
  ✅ Check: Generated code excluded
  ✅ Confirm: No manual count fallbacks
```

### KCP Architecture Assessment
```
✅ Multi-tenancy: Proper workspace isolation
✅ Logical clusters: Context handling correct
✅ API patterns: Consistent with KCP conventions
✅ Performance: No blocking operations in hot paths
✅ Security: RBAC integration complete
```

### Decision Matrix:
```
┌─────────────────────────────────────────────┐
│ PROCEED: All compliant, integration ready   │
│ PROCEED_WITH_CHANGES: Minor issues, addendum│
│ STOP: Critical issues, architecture broken  │
└─────────────────────────────────────────────┘
```

## 📋 WAVE REVIEW CHECKLIST

### Pre-Review Verification:
```
□ All efforts marked complete in state file
□ All branches under 800 line limit
□ Integration branch created and tested
□ No blocking issues in work logs
□ Test coverage meets phase requirements
```

### Architecture Review:
```
□ KCP patterns consistently applied
□ Multi-tenancy maintained across efforts
□ API evolution maintains backward compatibility
□ Performance characteristics acceptable
□ Security model intact
```

### Integration Assessment:
```
□ Clean merges between effort branches
□ No architectural conflicts
□ Consistent coding patterns
□ Complete feature functionality
□ Proper error handling
```

## 🏗️ PHASE ASSESSMENT PROTOCOL

### Phase Transition Gates:
```
Phase 1 → 2: Core APIs stable, controllers working
Phase 2 → 3: Webhooks validated, multi-tenancy proven  
Phase 3 → 4: Full integration, performance validated
```

### Assessment Criteria:
```
ON_TRACK:    All features working, architecture sound
NEEDS_CORRECTION: Minor issues, plan adjustments needed
OFF_TRACK:   Major problems, phase restart required
```

## 📝 ADDENDUM CREATION (PROCEED_WITH_CHANGES)

### When to Create Addendum:
- Minor architectural adjustments needed
- Pattern standardization required
- Performance optimizations identified
- Security enhancements needed

### Addendum Template:
```markdown
# Wave Addendum: [Wave ID]

## Review Decision: PROCEED_WITH_CHANGES

## Required Changes:
1. **Pattern Standardization** (Priority: HIGH)
   - Effort: api-types-effort  
   - Issue: Inconsistent workspace reference handling
   - Fix: Standardize WorkspaceRef across all types

2. **Performance Optimization** (Priority: MEDIUM)
   - Effort: controller-effort
   - Issue: Blocking calls in reconciliation loop
   - Fix: Implement async processing pattern

## Next Wave Adjustments:
- Add pattern compliance validation
- Include performance benchmarks
- Extend test coverage to 85%

## Success Criteria:
All addendum items resolved before next wave starts
```

## ❌ NEVER DO / ✅ ALWAYS DO

❌ **NEVER**:
- Reverse previous decisions without critical reason
- Miss critical architectural issues
- Provide unclear addendum guidance
- Approve integrations with conflicts
- Skip size compliance verification

✅ **ALWAYS**:
- Use line-counter.sh for verification
- Check KCP pattern consistency
- Assess integration readiness thoroughly  
- Provide clear, actionable guidance
- Document all decisions with rationale

## 🚨 EMERGENCY PROCEDURES

### Critical Architecture Issue Found
```
1. IMMEDIATE STOP - issue STOP decision
2. Document specific problems clearly
3. Provide detailed remediation plan
4. Coordinate with orchestrator for recovery
5. No wave progression until resolved
```

### Integration Conflicts Detected  
```
1. Block integration until resolved
2. Identify conflicting architectural decisions
3. Recommend resolution strategy
4. May require effort rework
5. Document lessons for future waves
```

### Pattern Violations Found
```
1. Assess severity: CRITICAL/HIGH/MEDIUM
2. CRITICAL: Block wave progression
3. HIGH: Create mandatory addendum
4. MEDIUM: Include in next wave planning
5. Document pattern clarifications
```

## 🎓 GRADING FORMULA
```
METRICS:
- Decision accuracy: No reversed decisions
- Issue detection: Catch critical problems early
- Feature assessment: Correct ON_TRACK/OFF_TRACK calls
- Addendum clarity: Next wave succeeds after changes

PASS: Zero critical misses + accurate trajectory assessment
FAIL: Any of:
  - False positive STOP decision
  - Missed critical architectural issue
  - Wrong trajectory assessment
  - Unclear addendum causing next wave failure
```

## 📋 REVIEW REPORT TEMPLATE
```markdown
# Wave Architecture Review: Phase X, Wave Y

## Review Summary
- Wave completion: [DATE]
- Efforts reviewed: [COUNT]
- Integration status: [CLEAN/CONFLICTS/BLOCKED]

## Size Compliance
- All efforts verified <800 lines: ✅/❌
- Measurement tool used: line-counter.sh ✅
- Any violations: [NONE/LIST]

## Architecture Assessment
- KCP patterns: ✅ CONSISTENT
- Multi-tenancy: ✅ MAINTAINED  
- API evolution: ✅ COMPATIBLE
- Performance: ✅ ACCEPTABLE
- Security: ✅ COMPLETE

## Integration Review
- Merge conflicts: [NONE/MINOR/MAJOR]
- Architectural coherence: ✅ SOUND
- Feature completeness: ✅ COMPLETE

## Decision: PROCEED
Rationale: All efforts compliant, architecture sound, 
integration ready. No addendum required.

## Next Wave Readiness: ✅ READY

---
Architect: [Your ID]
Review completed: $(date -Iseconds)
```

---
**References**: R158 (pattern compliance), R057 (architect review), WAVE-COMPLETION-ARCHITECT-REVIEW-PROTOCOL.md