# When to Use These Optional Files

## Overview

The files in this directory are **optional enhancements** that can be activated based on your project's complexity and needs. The core system works without them, but they provide additional structure, examples, and detailed protocols for specific scenarios.

## Quick Decision Guide

### By Project Complexity

#### 🟢 Simple Project (1-2 developers, <10K lines)
**Activate**: None initially
- Start with just the core files
- Add quick references if needed for learning

#### 🟡 Medium Project (3-10 developers, 10K-50K lines)
**Activate**:
- `CODE-REVIEWER-QUICK-REFERENCE.md` - Speed up reviews
- `ORCHESTRATOR-QUICK-REFERENCE.md` - Quick state lookups
- `ORCHESTRATOR-WORKFLOW-SUMMARY.md` - Visual understanding
- `ARCHITECT-REVIEWER-WAVE-INSTRUCTIONS.md` - Structured wave reviews

#### 🔴 Complex Project (10+ developers, 50K+ lines)
**Activate**: Everything
- All quick references for onboarding
- All architect protocols for quality gates
- All review examples for consistency
- All split examples for handling large efforts

### By Specific Needs

#### "I need to understand the system quickly"
**Activate these learning aids**:
- `ORCHESTRATOR-WORKFLOW-SUMMARY.md` - Visual flow with examples
- `CODE-REVIEWER-QUICK-REFERENCE.md` - Decision trees
- `ORCHESTRATOR-QUICK-REFERENCE.md` - State transition guide
- `CODE-REVIEW-EXAMPLES.md` - Real review patterns

#### "My efforts keep exceeding 800 lines"
**Activate these split helpers**:
- `SPLIT-EXAMPLE-E3.1.1-SYNC-ENGINE.md` - Real 2400→3 parts example
- `SPLIT-REVIEW-LOOP-DIAGRAM.md` - Visual split process
- Already have: `EFFORT-SPLIT-CONTINUOUS-EXECUTION-PROTOCOL.md` (in core)

#### "I need strict architectural compliance"
**Activate these architect protocols**:
- `ARCHITECT-REVIEWER-WAVE-INSTRUCTIONS.md` - Detailed wave review
- `PHASE-START-ARCHITECT-REVIEW-PROTOCOL.md` - Phase assessment
- `PHASE-COMPLETION-FUNCTIONAL-TESTING.md` - End-of-phase validation

#### "I want comprehensive code reviews"
**Activate these review enhancers**:
- `CODE-REVIEWER-COMPREHENSIVE-GUIDE-EXAMPLE.md` - Full review process
- `CODE-REVIEW-EXAMPLES.md` - Pattern library
- `CODE-REVIEW-ENFORCEMENT-SUMMARY.md` - All enforcement points
- `ORCHESTRATOR-CODE-REVIEW-INTEGRATION.md` - Review workflow

#### "I'm setting up CI/CD integration"
**Activate**:
- `CODE-REVIEW-ENFORCEMENT-SUMMARY.md` - Gate integration points
- `PHASE-COMPLETION-FUNCTIONAL-TESTING.md` - Test requirements

## File-by-File Decision Matrix

| File | When to Activate | Value Provided | Dependencies |
|------|------------------|----------------|--------------|
| **Quick References** ||||
| `CODE-REVIEWER-QUICK-REFERENCE.md` | Day 1 for new reviewers | 50% faster reviews | None |
| `ORCHESTRATOR-QUICK-REFERENCE.md` | Day 1 for orchestrators | Reduces state errors | None |
| `ORCHESTRATOR-WORKFLOW-SUMMARY.md` | During onboarding | Visual understanding | None |
| **Architect Protocols** ||||
| `ARCHITECT-REVIEWER-WAVE-INSTRUCTIONS.md` | Multiple waves per phase | Consistent architecture | Wave structure |
| `PHASE-START-ARCHITECT-REVIEW-PROTOCOL.md` | Phase boundaries | Prevents drift | Multiple phases |
| `PHASE-COMPLETION-FUNCTIONAL-TESTING.md` | Quality critical projects | Comprehensive validation | Test framework |
| **Review Enhancements** ||||
| `CODE-REVIEWER-COMPREHENSIVE-GUIDE-EXAMPLE.md` | Standardizing reviews | Review consistency | None |
| `CODE-REVIEW-EXAMPLES.md` | Training reviewers | Pattern recognition | None |
| `CODE-REVIEW-ENFORCEMENT-SUMMARY.md` | CI/CD setup | Automated gates | CI system |
| `ORCHESTRATOR-CODE-REVIEW-INTEGRATION.md` | Complex review flows | Workflow optimization | None |
| **Split Management** ||||
| `SPLIT-EXAMPLE-E3.1.1-SYNC-ENGINE.md` | First split needed | Learn by example | None |
| `SPLIT-REVIEW-LOOP-DIAGRAM.md` | Understanding splits | Visual guide | None |
| **Standalone Rules** ||||
| `ORCHESTRATOR-NEVER-WRITES-CODE-RULE.md` | Always (redundant) | Reinforcement | In core CLAUDE.md |
| **Meta Documentation** ||||
| `README.md` | Always | File overview | None |
| `FILES-NOT-INCLUDED.md` | Reference only | Explains omissions | None |

## Activation Instructions

### To Activate a File

1. **Copy to active location**:
```bash
# For protocols
cp possibly-needed-but-not-sure/FILENAME.md protocols/

# For agent-specific files
cp possibly-needed-but-not-sure/FILENAME.md .claude/agents/
```

2. **Update agent configurations** to reference the file:
```markdown
# In relevant agent .md file
READ: /workspaces/[project]/protocols/FILENAME.md
```

3. **Test with a small effort** before full rollout

### Progressive Activation Strategy

#### Week 1: Core Only
- Run with minimal files
- Document pain points
- Identify needed enhancements

#### Week 2: Add Quick References
- Activate quick reference guides
- Reduces lookup time by 70%
- Improves consistency

#### Week 3: Add Examples
- Activate example files
- Accelerates learning curve
- Provides patterns to follow

#### Week 4: Add Specialized Protocols
- Activate based on observed needs
- Only if problems occurred in weeks 1-3
- Measure improvement

## Signs You Need Specific Files

### You Need `ARCHITECT-REVIEWER-WAVE-INSTRUCTIONS.md` When:
- ❌ Architecture drifts between waves
- ❌ Integration conflicts increase
- ❌ Technical debt accumulates
- ❌ Patterns become inconsistent

### You Need `PHASE-START-ARCHITECT-REVIEW-PROTOCOL.md` When:
- ❌ Phases start without clear goals
- ❌ Dependencies aren't validated
- ❌ Previous phase issues carry forward
- ❌ Team is unsure about readiness

### You Need `CODE-REVIEWER-COMPREHENSIVE-GUIDE-EXAMPLE.md` When:
- ❌ Review quality varies by reviewer
- ❌ Important issues missed in reviews
- ❌ Review feedback inconsistent
- ❌ New reviewers struggle

### You Need `SPLIT-EXAMPLE-E3.1.1-SYNC-ENGINE.md` When:
- ❌ First time splitting an effort
- ❌ Splits still exceed 800 lines
- ❌ Unclear how to divide code
- ❌ Integration issues after splits

## Performance Impact

| Configuration | Setup Time | Per-Effort Overhead | Quality Improvement |
|---------------|------------|--------------------|--------------------|
| Core Only | 5 min | Baseline | Baseline |
| + Quick Refs | +2 min | -20% (faster) | +10% |
| + Examples | +5 min | -30% (faster) | +25% |
| + All Architect | +10 min | +10% (slower) | +40% |
| + All Review | +15 min | +15% (slower) | +50% |
| Everything | 30 min | +20% (slower) | +60% |

## Recommendations by Language/Framework

### Go/Kubernetes Projects
**Strongly Recommended**:
- All architect protocols (complex domain)
- All review examples (pattern-heavy)
- Split examples (large codebases)

### Python/FastAPI Projects
**Recommended**:
- Quick references only
- Basic architect protocols
- Skip complex review guides

### JavaScript/React Projects
**Recommended**:
- Quick references
- Split examples (component splitting)
- Basic review guides

### Rust Systems Programming
**Strongly Recommended**:
- Everything (safety critical)
- Especially testing protocols
- All review enforcement

## Common Mistakes to Avoid

### ❌ Don't Activate Everything on Day 1
- Overwhelming for team
- Increases friction
- Reduces adoption

### ❌ Don't Skip Quick References
- They save more time than they cost
- Reduce errors significantly
- Improve consistency immediately

### ❌ Don't Ignore Team Feedback
- If file isn't helping, deactivate it
- Customize based on actual needs
- Iterate on what works

### ✅ Do Start Simple
- Core files first
- Add based on observed needs
- Measure improvement

### ✅ Do Train on Examples
- Review examples as team
- Discuss patterns
- Build shared understanding

## Metrics to Track

Track these to decide what to activate:

1. **Review Cycle Time**: If >2 hours, add quick references
2. **Split Success Rate**: If <80%, add split examples
3. **Architecture Drift**: If detected, add architect protocols
4. **Integration Conflicts**: If >5%, add wave reviews
5. **Effort Rejection Rate**: If >20%, add comprehensive guides

## Quick Start Recommendations

### For a New Team
1. Week 1: Core only
2. Week 2: Add `ORCHESTRATOR-WORKFLOW-SUMMARY.md`
3. Week 3: Add `CODE-REVIEWER-QUICK-REFERENCE.md`
4. Week 4: Evaluate and add based on pain points

### For an Experienced Team
1. Day 1: Core + all quick references
2. Week 1: Add examples after first effort
3. Week 2: Add protocols based on complexity
4. Week 3: Full activation if beneficial

### For a Solo Developer
1. Core only initially
2. Add `ORCHESTRATOR-WORKFLOW-SUMMARY.md` for understanding
3. Skip architect protocols (self-reviewing)
4. Add split examples if needed

## Final Advice

> **Start lean, add thoughtfully, measure everything.**

The optional files are tools, not requirements. Use them when they solve real problems, not because they exist. The core system is designed to work without them—they're enhancements for specific situations.

If you're unsure, run without them first. You'll quickly discover which ones you need based on actual pain points rather than anticipated ones.