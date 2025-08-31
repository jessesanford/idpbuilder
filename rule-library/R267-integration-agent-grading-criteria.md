# 🚨🚨🚨 RULE R267: Integration Agent Grading Criteria 🚨🚨🚨

## Rule Definition
**Criticality:** BLOCKING
**Category:** Agent-Specific
**Applies To:** integration-agent

## GRADING BREAKDOWN: 100% TOTAL

### 50% - Completeness of Integration

#### 20% - Successful Branch Merging
- All specified branches merged into integration branch
- Correct merge order based on lineage
- No branches skipped or forgotten
- Proper merge commit messages

#### 15% - Conflict Resolution
- ALL conflicts fully resolved
- No conflict markers remain
- Code compiles after resolution
- Intelligent resolution decisions documented

#### 10% - Branch Integrity
- Original branches remain unmodified
- No cherry-picking used
- Proper synthesis branches created when needed
- Git history preserved correctly

#### 5% - Final State Validation
- Integration branch contains all changes
- Can checkout and build integration branch
- All expected files present
- Commit history is clean and logical

### 50% - Meticulous Tracking and Documentation

#### 25% - Work Log Quality
- **Replayability** (10%): Can reproduce exact integration
- **Completeness** (10%): Every operation documented
- **Clarity** (5%): Clear, understandable entries

#### 25% - Integration Report Quality
- **Comprehensiveness** (10%): All sections complete
- **Bug Documentation** (5%): Upstream issues documented
- **Test Results** (5%): Build/test attempts documented
- **Recommendations** (5%): Actionable next steps provided

## MANDATORY ACKNOWLEDGMENT

The integration agent MUST acknowledge grading criteria at startup:

```markdown
## 🎯 GRADING ACKNOWLEDGMENT

I understand I will be graded on:
- **50%** - Completeness of Integration
  - 20% Successful branch merging
  - 15% Conflict resolution
  - 10% Branch integrity
  - 5% Final validation
- **50%** - Meticulous Tracking and Documentation
  - 25% Work log quality (replayable, complete)
  - 25% Integration report quality

**Key Requirements I Acknowledge:**
- ✅ NEVER modify original branches
- ✅ NEVER use cherry-pick
- ✅ Create comprehensive INTEGRATION-PLAN.md first
- ✅ Document EVERY operation in work-log.md
- ✅ Complete INTEGRATION-REPORT.md with all sections
- ✅ DO NOT fix upstream bugs - only document
- ✅ Produce completely merged working copy
```

## Automatic Failures (Instant 0% Grade)

### Critical Violations
1. **Modified Original Branches** - Instant failure
2. **Used Cherry-Pick** - Instant failure
3. **Fixed Upstream Bugs** - Instant failure
4. **No Documentation** - Instant failure
5. **Incomplete Integration** - Branches missing

## Grade Calculation Example

```bash
calculate_integration_grade() {
    local score=0
    
    # Check critical violations first
    if original_branches_modified; then
        echo "Grade: 0% - Critical violation: original branches modified"
        return 0
    fi
    
    if cherry_pick_detected; then
        echo "Grade: 0% - Critical violation: cherry-pick used"
        return 0
    fi
    
    # Integration completeness (50% max)
    branches_merged && score=$((score + 20))
    conflicts_resolved && score=$((score + 15))
    branch_integrity_maintained && score=$((score + 10))
    final_state_valid && score=$((score + 5))
    
    # Documentation quality (50% max)
    work_log_replayable && score=$((score + 10))
    work_log_complete && score=$((score + 10))
    work_log_clear && score=$((score + 5))
    
    report_comprehensive && score=$((score + 10))
    bugs_documented && score=$((score + 5))
    tests_documented && score=$((score + 5))
    recommendations_provided && score=$((score + 5))
    
    echo "Grade: ${score}%"
}
```

## Performance Levels

| Grade | Level | Description |
|-------|-------|-------------|
| 90-100% | EXCELLENT | Perfect integration with exemplary documentation |
| 80-89% | GOOD | Complete integration, minor documentation gaps |
| 70-79% | SATISFACTORY | Integration works, documentation adequate |
| 60-69% | NEEDS IMPROVEMENT | Issues in integration or documentation |
| 0-59% | FAILING | Major issues or violations |
| 0% | CRITICAL FAILURE | Violated core rules |

## Self-Assessment Checklist

Before completing, the agent should verify:

```markdown
## Pre-Completion Checklist
- [ ] All branches from plan merged
- [ ] All conflicts resolved
- [ ] Original branches unmodified
- [ ] No cherry-picks used
- [ ] INTEGRATION-PLAN.md created and followed
- [ ] work-log.md is complete and replayable
- [ ] INTEGRATION-REPORT.md all sections filled
- [ ] All upstream bugs documented (not fixed)
- [ ] Build/test results documented
- [ ] Integration branch pushed to remote
- [ ] Final `git status` is clean
```

## Related Rules
- R260 - Integration Agent Core Requirements
- R261 - Integration Planning Requirements
- R262 - Merge Operation Protocols
- R263 - Integration Documentation Requirements
- R264 - Work Log Tracking Requirements
- R265 - Integration Testing Requirements
- R266 - Upstream Bug Documentation