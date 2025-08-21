# Example Split: E3.1.1 Sync Engine

## Scenario
Effort E3.1.1 (Sync Engine) from Phase 3, Wave 1 exceeds 800 lines after implementation.

## Initial Review Result
```yaml
review_date: 2025-08-20-14:30
branch: /phase3/wave1/effort1-sync-engine  
line_count: 1456
quality_pass: true
exception_evaluation:
  granted: false  # Reviewer determined it CAN be split logically
  reason: "While complex, sync engine has clear API/implementation boundaries"
```

## Split Plan Created by @agent-kcp-kubernetes-code-reviewer

```yaml
split_plan:
  original_effort: E3.1.1
  total_lines: 1456
  split_count: 3
  
  splits:
    - number: 1
      description: "Sync Engine APIs and Interfaces"
      target_lines: 350
      content:
        - pkg/syncer/engine/types.go
        - pkg/syncer/engine/interfaces.go
        - apis/syncer/v1alpha1/types.go
      
    - number: 2  
      description: "Core Sync Implementation"
      target_lines: 650
      content:
        - pkg/syncer/engine/engine.go
        - pkg/syncer/engine/reconciler.go
        - pkg/syncer/engine/state.go
        
    - number: 3
      description: "Sync Engine Tests and Helpers"
      target_lines: 456
      content:
        - pkg/syncer/engine/engine_test.go
        - pkg/syncer/engine/helpers.go
        - test/e2e/syncer/basic_test.go
```

## Orchestrator Execution

### Step 1: Rename Original Branch
```bash
cd /workspaces/efforts/phase3/wave1/effort1-sync-engine
git branch -m /phase3/wave1/effort1-sync-engine /phase3/wave1/effort1-sync-engine-to-be-split
git push origin --delete /phase3/wave1/effort1-sync-engine
git push -u origin /phase3/wave1/effort1-sync-engine-to-be-split
```

### Step 2: Create Split Working Copies (NO STOPPING)
```bash
# Create split 1 directory
mkdir -p /workspaces/efforts/phase3/wave1/effort1-sync-engine-split1
cd /workspaces/efforts/phase3/wave1/effort1-sync-engine-split1
git clone --no-checkout https://github.com/jessesanford/kcp.git .
git sparse-checkout init --cone
git sparse-checkout set pkg apis cmd test hack
git checkout phase2-integration  # Base from Phase 2
git checkout -b /phase3/wave1/effort1-sync-engine-part1

# Create split 2 directory  
mkdir -p /workspaces/efforts/phase3/wave1/effort1-sync-engine-split2
cd /workspaces/efforts/phase3/wave1/effort1-sync-engine-split2
git clone --no-checkout https://github.com/jessesanford/kcp.git .
git sparse-checkout init --cone
git sparse-checkout set pkg apis cmd test hack
git checkout /phase3/wave1/effort1-sync-engine-part1  # Depends on part1
git checkout -b /phase3/wave1/effort1-sync-engine-part2

# Create split 3 directory
mkdir -p /workspaces/efforts/phase3/wave1/effort1-sync-engine-split3
cd /workspaces/efforts/phase3/wave1/effort1-sync-engine-split3
git clone --no-checkout https://github.com/jessesanford/kcp.git .
git sparse-checkout init --cone
git sparse-checkout set pkg apis cmd test hack
git checkout /phase3/wave1/effort1-sync-engine-part2  # Depends on part2
git checkout -b /phase3/wave1/effort1-sync-engine-part3
```

### Step 3: Create Tracking Documentation
```bash
mkdir -p /workspaces/agent-configs/tmc-orchestrator-impl-8-20-2025/phase3/wave1/effort1/splits

cat > /workspaces/agent-configs/tmc-orchestrator-impl-8-20-2025/phase3/wave1/effort1/splits/SPLIT-TRACKING-20250820-1445.md << 'EOF'
# E3.1.1 Sync Engine Split Tracking

## Original Effort (OVERSIZED)
- Branch: /phase3/wave1/effort1-sync-engine-to-be-split  
- Size: 1456 lines (656 over limit)
- Status: TO BE SPLIT - DO NOT MERGE
- Location: /workspaces/efforts/phase3/wave1/effort1-sync-engine/

## Split Implementation
| Split | Branch | Content | Target | Actual | Status |
|-------|--------|---------|--------|--------|--------|
| 1 | /phase3/wave1/effort1-sync-engine-part1 | APIs/Interfaces | 350 | TBD | In Progress |
| 2 | /phase3/wave1/effort1-sync-engine-part2 | Core Implementation | 650 | TBD | Pending |
| 3 | /phase3/wave1/effort1-sync-engine-part3 | Tests/Helpers | 456 | TBD | Pending |

## Continuous Execution Log
| Time | Action | Result |
|------|--------|--------|
| 14:45 | Renamed original branch | Success |
| 14:46 | Created split working copies | Success |
| 14:47 | Started agent on split 1 | In Progress |

## Dependencies
- Base: phase2-integration
- Required for: All Phase 3 Wave 2 efforts
EOF
```

### Step 4: Task Agent for Split 1 (CONTINUE WITHOUT STOPPING)

```markdown
Task for @agent-kcp-go-lang-sr-sw-eng:

You are implementing split 1 of 3 for effort E3.1.1 Sync Engine.

CRITICAL CONTINUOUS EXECUTION:
- DO NOT STOP after any step - keep working
- DO NOT ASK for confirmation  
- Measure after EVERY file addition
- Complete ALL work before reporting

Working directory: /workspaces/efforts/phase3/wave1/effort1-sync-engine-split1
Branch: /phase3/wave1/effort1-sync-engine-part1
Target: 350 lines (MAXIMUM 800)

Instructions:
1. Cherry-pick ONLY the API/interface commits from phase7:
   git cherry-pick 73d9e2a1  # Sync engine interfaces
   git cherry-pick 184b0a59  # API types

2. Verify structure:
   - pkg/syncer/engine/types.go
   - pkg/syncer/engine/interfaces.go  
   - apis/syncer/v1alpha1/types.go

3. Measure immediately:
   /workspaces/kcp-shared-tools/tmc-pr-line-counter.sh -c /phase3/wave1/effort1-sync-engine-part1

4. Create work log:
   echo "Split 1 complete: $(date)" > SPLIT-WORK-LOG.md
   git add . && git commit -m "feat: sync engine APIs and interfaces (split 1/3)"
   git push -u origin /phase3/wave1/effort1-sync-engine-part1

5. Create completion checkpoint:
   echo "READY FOR SPLIT 2" > SPLIT-COMPLETION-CHECKPOINT.md
```

### Step 5: Immediately Continue with Split 2 (NO PAUSE)

The orchestrator immediately tasks another agent (or same agent) with split 2:

```markdown
Task for @agent-kcp-go-lang-sr-sw-eng:

You are implementing split 2 of 3 for effort E3.1.1 Sync Engine.
Previous split 1 is complete.

[Similar instructions for implementation files]
```

### Step 6: Immediately Continue with Split 3 (NO PAUSE)

```markdown
Task for @agent-kcp-go-lang-sr-sw-eng:

You are implementing split 3 of 3 for effort E3.1.1 Sync Engine.
Previous splits 1-2 are complete.

[Similar instructions for test files]
```

### Step 7: Review All Splits (STILL NO STOPPING)

```python
# Orchestrator triggers reviews immediately
for split_branch in [part1, part2, part3]:
    review = start_code_review(split_branch)
    if review.line_count > 800:
        # Handle recursive split WITHOUT STOPPING
        create_recursive_split_plan(split_branch)
```

### Step 8: Final Report (ONLY NOW REPORT COMPLETION)

```markdown
## E3.1.1 Split Complete

All 3 splits implemented, reviewed, and ready:
- Split 1: 342 lines ✅ (APIs/Interfaces)
- Split 2: 634 lines ✅ (Core Implementation)  
- Split 3: 448 lines ✅ (Tests/Helpers)

Total: 1424 lines (original was 1456)
All splits build and test independently.

Ready for PR creation and merge sequence.
```

## PR Creation for Splits

### PR 1: APIs and Interfaces
```markdown
## Summary
Part 1 of 3 splits from E3.1.1 Sync Engine effort.

Implements sync engine APIs and interfaces required by all subsequent implementations.

- Original effort: 1456 lines (656 over limit)
- This split: 342 lines ✅

## Related PRs (merge in order)
1. This PR: APIs and Interfaces
2. #XXX: Core Implementation (depends on this)
3. #XXX: Tests and Helpers (depends on #2)
```

### PR 2: Core Implementation
```markdown
## Summary  
Part 2 of 3 splits from E3.1.1 Sync Engine effort.

Implements core synchronization logic using APIs from PR #1.

- Depends on: PR #1 (APIs/Interfaces)
- This split: 634 lines ✅
```

### PR 3: Tests and Helpers
```markdown
## Summary
Part 3 of 3 splits from E3.1.1 Sync Engine effort.

Complete test coverage and helper functions for sync engine.

- Depends on: PR #1 and #2
- This split: 448 lines ✅
- Coverage: 92% (exceeds Phase 3 requirement of 90%)
```

## Key Lessons from Example

1. **NEVER STOP** between creating directories and implementing splits
2. **Sequential dependency** - Each split builds on previous
3. **Clear documentation** at every step
4. **Immediate measurement** after file additions
5. **Automatic review trigger** without waiting
6. **Only report** when ALL splits complete
7. **PR dependencies** clearly documented

This example shows the complete continuous execution flow for a complex effort that needs splitting.