# Software Factory 1.0 to 2.0 Migration Guide

## 🔄 Overview

Software Factory 2.0 introduces significant improvements over 1.0:
- **State Machine Driven** - Predictable, recoverable workflows
- **Grading System** - Performance metrics and enforcement
- **Pre-compaction Hooks** - Automatic state preservation
- **Parallel Spawn Optimization** - <5s timing requirements
- **Enhanced Rule System** - 160+ rules with unique IDs
- **Better Context Management** - 75-85% reduction in context usage

## 📊 Key Differences

| Feature | SF 1.0 | SF 2.0 |
|---------|--------|--------|
| State Management | Basic YAML | Full state machine per agent |
| Rule Organization | Scattered files | Hierarchical with IDs (R001-R160) |
| Context Recovery | Manual | Automatic with hooks |
| Grading | None | Real-time performance metrics |
| Agent Spawning | Sequential | Parallel with timing requirements |
| Size Management | Manual checking | Automatic with splits |
| TODO Management | Basic | State-aware with preservation |

## 🚀 Migration Approaches

### Option 1: Clean Migration (Recommended)
Start fresh with SF 2.0 while preserving your implementation progress.

**Pros:**
- Clean state machine initialization
- Optimal grading baseline
- No legacy issues

**Cons:**
- Requires mapping existing work to new structure

### Option 2: In-Place Upgrade
Upgrade existing SF 1.0 project to 2.0 structure.

**Pros:**
- Preserves Git history
- Maintains existing branches

**Cons:**
- May have state inconsistencies
- Requires careful validation

## 📋 Pre-Migration Checklist

Before migrating, ensure you have:

- [ ] Current SF 1.0 state backed up
- [ ] All branches pushed to remote
- [ ] Documentation of current phase/wave
- [ ] List of completed efforts
- [ ] Any custom rules or modifications documented
- [ ] Test results and coverage data

## 🔧 Migration Steps

### Step 1: Backup Current State

```bash
# Create backup of current SF 1.0 project
cd /workspaces/your-sf1-project
tar -czf ../sf1-backup-$(date +%Y%m%d).tar.gz .

# Export current state
cp orchestrator-state.json ../sf1-state-backup.yaml
cp -r todos ../sf1-todos-backup/

# Document current progress
git log --oneline -20 > ../sf1-git-history.txt
git branch -a > ../sf1-branches.txt
```

### Step 2: Analyze Current Progress

```bash
# Check current phase/wave
grep -E "current_phase|current_wave" orchestrator-state.json

# List completed efforts
grep -A 100 "efforts_completed:" orchestrator-state.json

# Check integration branches
git branch -a | grep integration

# Measure current sizes
for branch in $(git branch -r | grep effort); do
  echo "$branch:"
  git diff main..$branch --stat
done
```

### Step 3: Create Migration Mapping

Create `migration-mapping.yaml`:

```yaml
# SF 1.0 to 2.0 Migration Mapping
migration_date: "2024-XX-XX"

sf1_state:
  current_phase: 1
  current_wave: 3
  completed_efforts:
    - phase1/wave1/effort1-api-types
    - phase1/wave1/effort2-controllers
    - phase1/wave2/effort1-webhooks
  
sf2_mapping:
  phase_1:
    wave_1:
      status: COMPLETED
      integration_branch: phase1/wave1-integration
      efforts:
        - id: effort1-api-types
          size: 650
          review_status: APPROVED
        - id: effort2-controllers
          size: 780
          review_status: APPROVED
    wave_2:
      status: IN_PROGRESS
      efforts:
        - id: effort1-webhooks
          size: 720
          review_status: PENDING
```

### Step 4: Run Migration Script

Use the migration script (see below) or perform manual migration:

```bash
# Run automated migration
cd /workspaces/software-factory-2.0-template
./migrate-from-1.0.sh /workspaces/your-sf1-project

# Or manual migration
cp -r /workspaces/software-factory-2.0-template /workspaces/your-project-sf2
cd /workspaces/your-project-sf2
# Manually update configurations
```

### Step 5: Update State Files

#### Transform orchestrator-state.json:

**SF 1.0 Format:**
```yaml
current_phase: 1
current_wave: 3
efforts_completed: [...]
```

**SF 2.0 Format:**
```yaml
current_phase: 1
current_wave: 3
current_state: WAVE_START  # Add state machine state

# Add grading metrics
grading_history:
  parallel_spawn_average: 0.0
  review_first_try_rate: 0.0
  integration_success_rate: 0.0
  size_compliance_rate: 1.0

# Add checkpoint data
last_checkpoint: "2024-XX-XX"
checkpoint_version: "2.0"

# Preserve existing efforts data
efforts_completed: [...]
```

### Step 6: Migrate Custom Rules

If you have custom rules in SF 1.0:

1. **Identify custom rules:**
```bash
# Find custom rules not in original template
diff -r /workspaces/software-factory-template /workspaces/your-sf1-project
```

2. **Add to Rule Library:**
```bash
# Add custom rules to rule-library/CUSTOM-RULES.md
cat > rule-library/CUSTOM-RULES.md << EOF
# Custom Rules from SF 1.0

## Rule RC01.0.0 - [Your Custom Rule Name]
**Source:** Migrated from SF 1.0
**Category:** Custom
**Content:** [Your rule content]
EOF
```

3. **Update Rule Registry:**
Add entries to `rule-library/RULE-REGISTRY.md`

### Step 7: Update Agent Configurations

Transform `.claude/commands/continue-orchestrating.md`:

**Add SF 2.0 sections:**
```markdown
## Pre-Flight Checks (New in 2.0)
┌─────────────────────────────────────────────────────────────────┐
│ RULE R001.0.0 - Pre-Flight Checks                              │
├─────────────────────────────────────────────────────────────────┤
│ Print startup acknowledgment                                    │
│ Verify environment                                             │
│ Load state machine                                             │
└─────────────────────────────────────────────────────────────────┘

## Grading Requirements (New in 2.0)
- Parallel spawn: <5s average
- Review success: >80% first-try
- Integration success: >90%
```

### Step 8: Set Up Hooks

```bash
# Enable pre-compaction hooks
chmod +x hooks/*.sh

# Test hook functionality
./hooks/pre-compact.sh
./hooks/post-compact.sh

# Configure Git hooks if needed
ln -s ../../hooks/pre-compact.sh .git/hooks/pre-commit
```

### Step 9: Validate Migration

Run validation checks:

```bash
# Check state consistency
./hooks/state-snapshot.sh

# Verify directory structure
tree -L 2 -d

# Test slash commands
grep -l "Pre-Flight Checks" .claude/commands/*.md

# Verify rule references
grep -r "RULE R[0-9]" --include="*.md" | head -20

# Check grading system
find . -name "grading.md" | wc -l  # Should be many
```

### Step 10: Test with Agents

```markdown
1. Start with status check:
   /check-status

2. Verify state machine:
   Check current_state in orchestrator-state.json

3. Test continuation:
   /continue-orchestrating

4. Verify grading:
   Check parallel_spawn_average after spawning agents
```

## 🔄 Post-Migration Tasks

### Update Working Practices

1. **Use State Machine States:**
   - Always check `current_state` not just phase/wave
   - Follow state transitions strictly

2. **Monitor Grading:**
   - Check grading metrics regularly
   - Address failures immediately

3. **Leverage Pre-compaction:**
   - Trust automatic TODO preservation
   - Use recovery assistant when needed

4. **Follow Rule IDs:**
   - Reference rules by ID (R001, R002, etc.)
   - Check rule-library for details

### Configure for Your Project

1. **Update project-config.yaml:**
```yaml
project:
  name: "your-project"
  migrated_from: "SF 1.0"
  migration_date: "2024-XX-XX"
```

2. **Set Constraints:**
```yaml
constraints:
  max_lines_per_effort: 800  # Enforce strictly in 2.0
  test_coverage_target: 80
  parallel_spawn_limit: 3
```

## ⚠️ Common Migration Issues

### Issue 1: State Machine Mismatch
**Problem:** Current work doesn't map to a state
**Solution:** Set state to closest match, use WAVE_START if unsure

### Issue 2: Missing Grading History
**Problem:** No baseline metrics
**Solution:** Initialize with conservative values:
```yaml
grading_history:
  parallel_spawn_average: 10.0  # Will improve
  review_first_try_rate: 0.5
  integration_success_rate: 0.7
```

### Issue 3: Rule Reference Errors
**Problem:** Old rules don't have IDs
**Solution:** Map to closest SF 2.0 rule or create custom rule with ID

### Issue 4: TODO State Loss
**Problem:** TODOs not in new format
**Solution:** Convert to new naming: `{agent}-{state}-{timestamp}.todo`

## 📊 Migration Success Metrics

Your migration is successful when:

- [ ] All slash commands work
- [ ] State machine transitions properly
- [ ] Grading metrics are tracked
- [ ] Pre-compaction hooks trigger
- [ ] Agents acknowledge rules on startup
- [ ] Size limits are enforced
- [ ] TODOs are preserved across sessions

## 🆘 Rollback Procedure

If migration fails:

```bash
# Restore from backup
cd /workspaces
tar -xzf sf1-backup-*.tar.gz
mv sf1-state-backup.yaml your-project/orchestrator-state.json
mv sf1-todos-backup/* your-project/todos/

# Return to SF 1.0
cd your-project
git checkout main  # Or your SF 1.0 branch
```

## 📚 Additional Resources

- **Quick Reference:** `quick-reference/` - New 2.0 features
- **State Machines:** `state-machines/` - Understand new flow
- **Rule Library:** `rule-library/` - All rules with IDs
- **Expertise Modules:** `expertise/` - Enhanced patterns

## 🎯 Benefits After Migration

Once migrated to SF 2.0, you'll have:

1. **Better Recovery:** Automatic context preservation
2. **Performance Tracking:** Know exactly how well agents perform
3. **Stricter Compliance:** Automatic enforcement of limits
4. **Clearer Workflow:** State machine removes ambiguity
5. **Faster Development:** Parallel spawning when possible
6. **Better Documentation:** All rules have IDs and locations

---

**Need help?** Check `quick-reference/emergency-procedures.md` or run `/check-status` for diagnostics.