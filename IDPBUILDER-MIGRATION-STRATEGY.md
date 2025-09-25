# IDPBuilder SF 1.0 to 2.0 Migration Strategy

## 📊 Current State Analysis

Your SF 1.0 idpbuilder project has:
- **Progress:** Phase 1, Wave 1 complete (2/38 efforts = 5%)
- **Planning Artifacts:** High-level phase templates, no effort-level plans yet
- **Implementation:** 2 completed efforts (want to discard these)
- **Architecture:** Well-defined phase/wave structure

## 🎯 Migration Recommendation: **Hybrid Approach**

### Keep These SF 1.0 Planning Artifacts (High Value):

1. **Phase Structure** ✅ KEEP
   - 5 phases with clear goals
   - 15 waves with dependencies
   - 38 efforts identified
   - This decomposition is solid and reusable

2. **Project Architecture** ✅ KEEP
   ```yaml
   Phase 1: MVP Core (3 waves, 7 efforts)
   Phase 2: CLI Integration (3 waves, 7 efforts)  
   Phase 3: Production Ready (3 waves, 8 efforts)
   Phase 4: Enhanced Features (3 waves, 8 efforts)
   Phase 5: Advanced Features (3 waves, 8 efforts)
   ```

3. **Dependency Graph** ✅ KEEP
   - Wave dependencies (1.2 depends on 1.1)
   - Effort parallelization opportunities
   - Critical path identification

4. **Test Coverage Requirements** ✅ KEEP
   ```yaml
   Phase 1: 70%
   Phase 2: 80%
   Phase 3: 90%  # Production ready
   Phase 4: 85%
   Phase 5: 80%
   ```

### Regenerate with SF 2.0 Agents (Better Quality):

1. **Effort-Level Implementation Plans** 🔄 REGENERATE
   - SF 2.0 Code Reviewer creates better structured plans
   - Includes proper work-log.md templates
   - Has grading criteria built-in
   - Better test scaffolding

2. **Review Criteria** 🔄 REGENERATE
   - SF 2.0 has comprehensive grading rubrics
   - KCP pattern compliance checks
   - Security validation
   - Performance criteria

3. **Integration Strategies** 🔄 REGENERATE
   - SF 2.0 handles wave/phase integration better
   - Automatic conflict detection
   - Better merge strategies

## 📋 Migration Plan

### Step 1: Create Fresh SF 2.0 Project
```bash
cd /workspaces/software-factory-2.0-template
./setup.sh

# When prompted:
Project name: idpbuilder
Description: Container Build and Push Feature for IDP Builder
Language: Go
Technologies: Kubernetes/KCP, gRPC
Agents: All (Orchestrator, SW Engineer, Code Reviewer, Architect)
Existing plan: Yes (we'll import the phase structure)
```

### Step 2: Import Phase Structure
Create `/workspaces/idpbuilder-sf2/MASTER-IMPLEMENTATION-PLAN.md`:

```markdown
# IDPBuilder Implementation Plan (Migrated from SF 1.0)

## Project Overview
**Total Phases:** 5
**Total Waves:** 15  
**Total Efforts:** 38
**Target:** 6-7 weeks

## Phase 1: MVP Core - Minimal Working Build and Push

### Wave 1.1: Essential API Contracts (COMPLETED in 1.0 - REDO)
- E1.1.1: Minimal Build Types (196 lines)
- E1.1.2: Builder Interface (198 lines)

### Wave 1.2: Core Libraries
- E1.2.1: Buildah Client (est. 400 lines)
- E1.2.2: Registry Client (est. 300 lines)

### Wave 1.3: MVP Implementation
- E1.3.1: Build Controller (est. 600 lines)
- E1.3.2: Push Logic (est. 500 lines)
- E1.3.3: Self-Signed Cert Support (est. 300 lines)

[Continue with remaining phases...]
```

### Step 3: Configure SF 2.0 Specific Settings

Update `project-config.yaml`:
```yaml
project:
  name: "idpbuilder"
  migrated_from: "SF 1.0"
  restart_reason: "Want SF 2.0 quality and automation"
  
implementation:
  discard_sf1_code: true  # Start fresh
  keep_sf1_planning: true  # Preserve architecture
  
constraints:
  max_lines_per_effort: 800
  test_coverage_targets:
    phase_1: 70
    phase_2: 80
    phase_3: 90
    phase_4: 85
    phase_5: 80
```

### Step 4: Initialize Orchestrator State

Create initial `orchestrator-state.json`:
```yaml
current_phase: 1
current_wave: 1  # Restart from Wave 1
current_state: INIT

# Import high-level structure from SF 1.0
phases:
  - phase: 1
    name: "MVP Core - Minimal Working Build and Push"
    waves_total: 3
    efforts_total: 7
    # But mark as NOT_STARTED for fresh implementation
    status: "NOT_STARTED"
    waves_completed: 0

# SF 2.0 additions
grading_history:
  parallel_spawn_average: 0.0
  review_first_try_rate: 0.0
  integration_success_rate: 0.0

state_machine_version: "2.0"
```

### Step 5: Let SF 2.0 Generate Better Plans

Start the orchestrator and let it generate effort plans:

```markdown
/continue-orchestrating

# Orchestrator will:
1. Load phase structure from MASTER-IMPLEMENTATION-PLAN.md
2. Spawn Code Reviewer to create effort plans for Wave 1
3. Plans will include:
   - Detailed IMPLEMENTATION-PLAN.md per effort
   - work-log.md templates
   - Test scaffolding
   - Grading criteria
```

## 🔄 Key Differences in SF 2.0 Generated Plans

### SF 1.0 Phase Plan (Template)
```markdown
### E1.1.1: Core API Types
**Requirements:**
1. Create base API group structure
2. Implement core types
3. Include standard patterns

**Test Requirements:**
[Basic test outline]
```

### SF 2.0 Effort Plan (Generated)
```markdown
# Implementation Plan for E1.1.1: Core API Types

## Grading Criteria
┌─────────────────────────────────────────────────────────────────┐
│ RULE R007.0.0 - Size Limit: MUST be <800 lines                 │
│ RULE R032.0.0 - Test Coverage: MUST be >70%                    │
│ RULE R037.0.0 - KCP Patterns: MUST follow multi-tenancy        │
└─────────────────────────────────────────────────────────────────┘

## Implementation Approach
1. Directory Structure:
   ```
   apis/
   └── v1alpha1/
       ├── types.go          # ~200 lines
       ├── defaults.go       # ~100 lines
       ├── validation.go     # ~150 lines
       ├── zz_generated.go   # (excluded from count)
       └── types_test.go     # ~200 lines
   ```

2. Test-Driven Development:
   - Write types_test.go FIRST
   - Implement types to pass tests
   - Add validation with tests
   - Verify 70% coverage

## Work Checkpoints
- [ ] CP1: Test scaffolding complete (30 min)
- [ ] CP2: Basic types implemented (1 hr)
- [ ] CP3: Validation added (30 min)
- [ ] CP4: DeepCopy generated (15 min)
- [ ] CP5: All tests passing (15 min)

## Size Monitoring
Check every checkpoint:
```bash
/workspaces/idpbuilder/tools/line-counter.sh -c $(git branch --show-current)
```

## Review Checklist
- [ ] Follows KCP patterns
- [ ] Multi-tenancy safe
- [ ] Proper validation
- [ ] Test coverage >70%
- [ ] Size <800 lines
```

## 🎯 Benefits of This Hybrid Approach

1. **Preserve Investment** - Keep your good architectural planning
2. **Better Quality** - SF 2.0 generates superior effort plans
3. **Automatic Enforcement** - Size limits, test coverage, grading
4. **Faster Development** - Parallel spawning, better recovery
5. **Clean Start** - No legacy code issues

## 📊 Comparison: Planning Artifacts

| Artifact | SF 1.0 | SF 2.0 | Recommendation |
|----------|--------|--------|----------------|
| Phase decomposition | ✅ Have | Use as-is | KEEP |
| Wave structure | ✅ Have | Use as-is | KEEP |
| Effort list | ✅ Have | Use as-is | KEEP |
| Dependencies | ✅ Have | Use as-is | KEEP |
| Effort plans | ❌ None yet | Generate better | REGENERATE |
| Work logs | ❌ None yet | Auto-template | REGENERATE |
| Review criteria | ❌ Basic | Comprehensive | REGENERATE |
| Test scaffolds | ❌ Basic | TDD templates | REGENERATE |

## 🚀 Recommended Migration Commands

```bash
# 1. Create new SF 2.0 project
cd /workspaces/software-factory-2.0-template
./setup.sh
# Answer: idpbuilder, Go, Kubernetes/KCP, etc.

# 2. Copy phase planning
cp /workspaces/idpbuidler-software-factory-attempt1/phase-plans/*.md \
   /workspaces/idpbuilder-sf2/planning/

# 3. Create master plan
cat > /workspaces/idpbuilder-sf2/MASTER-IMPLEMENTATION-PLAN.md << 'EOF'
[Paste consolidated phase structure]
EOF

# 4. Start orchestrator
cd /workspaces/idpbuilder-sf2
# In Claude: /continue-orchestrating

# 5. Let SF 2.0 generate effort plans
# Orchestrator spawns Code Reviewer
# Code Reviewer creates IMPLEMENTATION-PLAN.md for each effort
# Much better quality than SF 1.0 templates
```

## ⚠️ Important Notes

1. **Don't copy SF 1.0 code** - Start fresh with SF 2.0 quality
2. **Don't copy effort plans** - Let SF 2.0 generate better ones
3. **Do copy architecture** - Your phase/wave structure is good
4. **Do copy test requirements** - Your coverage targets are appropriate

## 📈 Expected Outcomes

After migration to SF 2.0:
- **Development Speed**: 2-3x faster with parallel spawning
- **Code Quality**: 85%+ first-review pass rate (vs 60% typical)
- **Size Compliance**: 100% automatic (vs manual checking)
- **Recovery Time**: 1-2 minutes (vs 15+ minutes)
- **Context Efficiency**: 75% less context needed

---

**Ready to migrate?** This hybrid approach gives you the best of both worlds:
- Your good architectural planning from SF 1.0
- Superior execution and automation from SF 2.0