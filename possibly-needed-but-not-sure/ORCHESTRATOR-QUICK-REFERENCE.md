# Orchestrator Quick Reference Card

## Phase Execution Order (STRICT)
```
Phase 1 → Phase 2 → Phase 3 → Phase 4 → Phase 5
```

## Parallel Execution Rules

### CAN Parallelize
- Phase 1: API types within same wave
- Phase 2: Controllers after base controller
- Phase 3: Resource syncers after sync engine
- Phase 4: Feature gaps after cross-workspace fix
- Phase 5: All testing efforts

### MUST Serialize
- Base controllers (E2.1.1)
- Sync engine (E3.1.1)
- Cross-workspace fix (E4.1.1)
- Final integration (E5.5.*)
- ALL BRANCH SPLITS

## Agent Commands

### Start Implementation Agent
```bash
# 1. Create working copy
mkdir -p /workspaces/efforts/phase${P}/wave${W}/effort${E}
cd /workspaces/efforts/phase${P}/wave${W}/effort${E}
git clone --no-checkout https://github.com/jessesanford/kcp.git .
git sparse-checkout init --cone
git sparse-checkout set pkg apis cmd test hack

# 2. Create branch
git checkout ${BASE_BRANCH}
git checkout -b /phase${P}/wave${W}/effort${E}-${NAME}

# 3. Seed work log
echo "# Work Log E${P}.${W}.${E}" > WORK-LOG-$(date +%Y%m%d-%H%M).md
git add WORK-LOG-*.md && git commit -m "Initialize ${NAME}" && git push -u origin

# 4. Start agent with instructions
```

### Start Review Agent
```bash
# Point to effort branch
BRANCH="/phase${P}/wave${W}/effort${E}-${NAME}"

# Run line counter
/workspaces/kcp-shared-tools/tmc-pr-line-counter.sh ${BRANCH}

# If > 800 lines, create split plan
# Otherwise, create review
```

## Critical Validation Commands

### After Every Cherry-pick
```bash
git status | grep -q conflict && echo "CONFLICT" || echo "OK"
go build ./... || echo "BUILD FAILED"
```

### After Every Effort
```bash
go test ./... -race -cover
golangci-lint run ./...
make generate && git diff --exit-code
```

### After Every Split (ITERATIVE)
```bash
# Check split size
/workspaces/kcp-shared-tools/tmc-pr-line-counter.sh ${SPLIT_BRANCH}
# If > 800, create new split plan and repeat
# If <= 800, proceed with quality checks
go test ./... -race -cover
golangci-lint run ./...
```

### After Every Phase
```bash
git checkout phase${N}-integration
for branch in ${EFFORT_BRANCHES}; do
    git merge --no-ff $branch || exit 1
    go build ./... || exit 1
done
```

## Branch Naming Convention
```
/phase{X}/wave{Y}/effort{Z}-{descriptive-name}
```
- X = Phase number (1-5)
- Y = Wave number
- Z = Effort number
- descriptive-name = e.g., "api-types-core"

## Coverage Requirements by Phase
- Phase 1: 80% (API types)
- Phase 2: 85% (Controllers)
- Phase 3: 90% (Syncer - CRITICAL)
- Phase 4: 85% (Features)
- Phase 5: 95% (Final testing)

## Line Count Exception Rules
- Default limit: 800 lines per PR
- Exception authority: @agent-kcp-kubernetes-code-reviewer ONLY
- Valid exception reasons:
  - Atomic transactions that cannot be split
  - Complex state machines
  - Tightly coupled interface/implementation
  - Generated code blocks
  - Critical bug fixes spanning multiple components
- Exception requires:
  - Documented justification
  - Risk mitigation plan
  - Extra validation/reviewers

## Key Commits to Cherry-pick

### Phase 1 - API Types
```bash
184b0a593  # NegotiatedAPIResource types
2cbbc56c3  # SyncTarget API types
b629c4831  # Split implementation
```

### Phase 2 - Controllers (USE tmc2-impl2!)
```bash
# From tmc2-impl2/00a1-controller-patterns (8 commits)
git cherry-pick $(git log --oneline origin/feature/tmc2-impl2/00a1-controller-patterns | head -8 | awk '{print $1}' | tac)
```

### Phase 3 - Syncer (USE phase7!)
```bash
# From phase7-syncer-impl/p7w1-sync-engine
for commit in $(git rev-list --reverse origin/feature/phase7-syncer-impl/p7w1-sync-engine); do
    git cherry-pick $commit
done
```

### Phase 4 - Cross-workspace Fix
```bash
# MUST fix the bug from contrib-tmc
git cherry-pick $(git log --oneline origin/feature/tmc-phase4-19-crossworkspace-controller | awk '{print $1}')
# Then apply bugfix.go from instructions
```

## Recovery Procedures

### From Conflict
```bash
# Option 1: Take source version
git checkout --theirs . && git add -A && git cherry-pick --continue

# Option 2: Reset and try different approach
git cherry-pick --abort
git reset --hard HEAD
```

### From Failed Build
```bash
# Check what's broken
go build -v ./... 2>&1 | head -20

# Common fixes
go mod tidy
make generate
```

### From Interrupted Agent
```bash
# Check work log
cat WORK-LOG-*.md | tail -50

# Resume from last commit
git log --oneline -5
```

## State File Updates
```yaml
# Update after every event
/workspaces/orchestrator-state.yaml:
  current_phase: X
  current_wave: Y
  efforts_completed: [...]
  efforts_in_progress: [...]
```

## Final Success Checklist
- [ ] All 87 efforts complete
- [ ] All code reviews passed
- [ ] No branches > 800 lines
- [ ] Cross-workspace bug fixed
- [ ] 8 feature gaps filled
- [ ] Final validation passes
- [ ] Ready for main merge

## Continuous Execution Rules for Splits

### DO NOT STOP After:
- Creating split working copies
- Implementing first split
- Any individual split completion
- Setting up directories

### ONLY STOP When:
- ALL splits complete
- Error prevents continuation
- Split still >800 lines (needs recursive split)

### Split Directory Structure
```
/workspaces/efforts/phase{X}/wave{Y}/effort{Z}-{name}/          # Original
/workspaces/efforts/phase{X}/wave{Y}/effort{Z}-{name}-split1/   # Split 1
/workspaces/efforts/phase{X}/wave{Y}/effort{Z}-{name}-split2/   # Split 2
```

### Split Branch Naming
```
Original: /phase{X}/wave{Y}/effort{Z}-{name}
Renamed:  /phase{X}/wave{Y}/effort{Z}-{name}-to-be-split
Split 1:  /phase{X}/wave{Y}/effort{Z}-{name}-part1
Split 2:  /phase{X}/wave{Y}/effort{Z}-{name}-part2
```

### Split Tracking Docs
Create in: `/workspaces/agent-configs/tmc-orchestrator-impl-8-20-2025/phase{X}/wave{Y}/effort{Z}/splits/`
- SPLIT-TRACKING-{DATE}.md
- SPLIT-WORK-LOG-{DATE}.md
- SPLIT-COMPLETION-CHECKPOINT.md

## Phase Completion Testing

### After Each Phase
1. Create integration branch with all efforts merged
2. Build binaries: `make build-all`
3. Create test directory: `/workspaces/tests/phase{N}-functional/`
4. Task @agent-kcp-kubernetes-code-reviewer to create test harness
5. **USER MUST TEST** before proceeding to next phase

### Test Harness Requirements
- Interactive script showing TMC features
- Based on: `/workspaces/agent-configs/example-functional-test-scripts/tmc-multi-cluster-demo.sh`
- User can see features working
- Cleanup after completion

## Emergency Contacts
- Phase 1-2 Issues: Check controller patterns
- Phase 3 Issues: Verify using phase7 syncer
- Phase 4 Issues: Ensure cross-workspace fix applied
- Phase 5 Issues: Check test coverage requirements
- Split Issues: Never parallelize, always sequential
- Continuous Execution: See EFFORT-SPLIT-CONTINUOUS-EXECUTION-PROTOCOL.md
- Testing: See PHASE-COMPLETION-FUNCTIONAL-TESTING.md

## Remember
1. **Exact commands only** - No pseudo-code
2. **Validate everything** - Build/test after each step
3. **Keep agents working** - 24/7 until complete
4. **Document everything** - Work logs are critical
5. **Review all code** - Every effort needs review