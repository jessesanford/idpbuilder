# 🚨🚨🚨 BLOCKING RULE R053: Parallelization Decisions

## Rule Statement
The orchestrator MUST analyze effort dependencies and determine parallelization strategy before spawning agents. Independent efforts MUST be executed in parallel, while dependent efforts MUST be executed sequentially.

## Parallelization Requirements

### 1. Dependency Analysis
- **READ** all IMPLEMENTATION-PLAN.md files created by Code Reviewers
- **IDENTIFY** effort dependencies from the plans
- **CREATE** dependency graph showing which efforts can be parallel
- **DOCUMENT** parallelization strategy in orchestrator-state-v3.json

### 2. Parallel Execution Criteria
Efforts can be spawned in parallel when:
- ✅ No file overlap between efforts
- ✅ No API dependencies between efforts
- ✅ No shared database schema changes
- ✅ No sequential ordering specified in plans
- ✅ Different subsystems or modules

### 3. Sequential Execution Required
Efforts MUST be sequential when:
- ❌ Shared file modifications detected
- ❌ API dependencies exist (one needs the other's output)
- ❌ Database migrations that must be ordered
- ❌ Plan explicitly states sequential requirement
- ❌ Integration points require specific order

### 4. Parallelization Decision Documentation
```yaml
# In orchestrator-state-v3.json
parallelization_analysis:
  analyzed_at: "2025-01-20T10:30:00Z"
  parallel_groups:
    - group_1:
        efforts: [effort-1, effort-2, effort-3]
        reason: "Independent modules, no file overlap"
    - group_2:
        efforts: [effort-4]
        reason: "Depends on effort-1 completion"
  sequential_requirements:
    - effort-4 after effort-1: "API dependency"
    - effort-5 after all: "Integration effort"
```

## Enforcement

### Pre-Spawn Verification
```bash
# MANDATORY before spawning SW Engineers
verify_parallelization_analysis() {
    local state_file="orchestrator-state-v3.json"
    
    # Check analysis exists
    # Use text_editor tool with view command to check orchestrator-state-v3.json:
    # Look for parallelization_analysis field
    if [ ! parallelization_analysis field exists ]; then
        echo "❌ BLOCKING: No parallelization analysis found!"
        echo "Must analyze dependencies before spawning agents"
        exit 1
    fi
    
    # Check timestamp
    # Use text_editor tool with view command to read state file:
    # Find the parallelization_analysis.analyzed_at field
    local analyzed_at="<value from parallelization_analysis.analyzed_at>"
    if [[ -z "$analyzed_at" ]]; then
        echo "❌ BLOCKING: Parallelization analysis not timestamped!"
        exit 1
    fi
    
    echo "✅ Parallelization strategy verified"
}
```

### Spawning Pattern
```markdown
# For PARALLEL efforts (spawn together):
echo "🚀 Spawning parallel group 1: effort-1, effort-2, effort-3"
Task sw-engineer-1: Implement effort-1 [spawn immediately]
Task sw-engineer-2: Implement effort-2 [spawn immediately]  
Task sw-engineer-3: Implement effort-3 [spawn immediately]

# For SEQUENTIAL efforts (spawn after predecessor completes):
echo "⏳ Waiting for effort-1 completion before spawning effort-4"
[Monitor effort-1 until complete]
Task sw-engineer-4: Implement effort-4 [spawn only after effort-1 done]
```

## Common Violations

### ❌ VIOLATION: Spawning Without Analysis
```markdown
# WRONG - No dependency check
Task sw-engineer-1: Implement effort-1
Task sw-engineer-2: Implement effort-2
Task sw-engineer-3: Implement effort-3
```

### ❌ VIOLATION: Ignoring Dependencies
```markdown
# WRONG - effort-4 depends on effort-1 but spawned together
Task sw-engineer-1: Implement effort-1
Task sw-engineer-4: Implement effort-4  # Should wait!
```

### ✅ CORRECT: Proper Analysis and Spawning
```markdown
1. Read all IMPLEMENTATION-PLAN.md files
2. Analyze dependencies
3. Document in orchestrator-state-v3.json
4. Spawn parallel groups together
5. Monitor and spawn sequential efforts when ready
```

## Integration with Other Rules

- **R151**: Parallel agents must emit timestamps within 5s
- **R197**: One agent per effort (never multiple agents for same effort)
- **R208**: CD to correct directory before each spawn
- **R052**: Provide complete context to each spawned agent

## Grading Impact

- **-50%**: Spawning agents without parallelization analysis
- **-30%**: Spawning dependent efforts in parallel
- **-20%**: Not documenting parallelization decisions
- **-15%**: Sequential spawning when parallel was possible

## Example Analysis Output

```yaml
# After analyzing 5 efforts in wave 2:
parallelization_analysis:
  analyzed_at: "2025-01-20T14:45:00Z"
  total_efforts: 5
  parallel_groups:
    - group_1:
        efforts: ["auth-service", "logging-module", "metrics-collector"]
        reason: "Independent services with no shared files"
        spawn_together: true
    - group_2:  
        efforts: ["api-gateway"]
        reason: "Depends on auth-service API"
        spawn_after: ["auth-service"]
    - group_3:
        efforts: ["integration-tests"]
        reason: "Requires all services running"
        spawn_after: ["api-gateway"]
  execution_order:
    1: ["auth-service", "logging-module", "metrics-collector"]  # Parallel
    2: ["api-gateway"]  # After auth-service
    3: ["integration-tests"]  # After all
```

---

**REMEMBER**: Parallelization decisions MUST be made BEFORE spawning agents. The Code Reviewer's implementation plans contain the dependency information needed for this analysis.