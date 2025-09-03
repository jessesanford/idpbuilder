# Rule R218: Orchestrator Parallel Code Reviewer Spawning for Effort Planning

## Rule Statement
The Orchestrator MUST:
1. **USE THE READ TOOL** to read the wave implementation plan
2. **EXPLICITLY ACKNOWLEDGE** reading the parallelization headers 
3. **STATE THE FILE PATH** of the wave plan being used
4. **SPAWN** Code Reviewers according to the parallelization metadata found

## Criticality Level
**MANDATORY** - Failure to acknowledge = -25% (Protocol violation)

## Enforcement Mechanism
- **Technical**: MUST use Read tool on wave plan file
- **Behavioral**: MUST output acknowledgment message with file path
- **Grading**: -25% if no acknowledgment, -50% if wrong spawning strategy
- **Audit**: Check for "as allowed by the parallelization headers in wave plan located at:" message

## Core Principle

```
Wave Plan Says Parallel → Spawn Code Reviewers Together
Wave Plan Says Sequential → Spawn Code Reviewers One at a Time
```

## Detailed Requirements

### 🚨 MANDATORY OUTPUT: Acknowledge Reading Wave Plan

**THE ORCHESTRATOR MUST OUTPUT ALL OF THE FOLLOWING:**

```
📖 R218: Using Read tool to read wave plan...
[Orchestrator MUST use Read tool on phase-plans/PHASE-X-WAVE-Y-IMPLEMENTATION-PLAN.md]

📖 R218: I have READ the parallelization headers from wave plan located at:
   phase-plans/PHASE-X-WAVE-Y-IMPLEMENTATION-PLAN.md
   
✅ Spawning Code Reviewers as allowed by the parallelization headers in wave plan located at:
   phase-plans/PHASE-X-WAVE-Y-IMPLEMENTATION-PLAN.md
   
Efforts marked "Can Parallelize: Yes": [list]
Efforts marked "Can Parallelize: No": [list]
```

**GRADING CHECKLIST:**
- [ ] Did orchestrator use Read tool on wave plan? (-25% if no)
- [ ] Did orchestrator say "as allowed by the parallelization headers in wave plan located at: [path]"? (-25% if no)
- [ ] Did orchestrator list which efforts are parallel vs sequential? (-10% if no)
- [ ] Did orchestrator spawn according to the metadata? (-50% if wrong)

### MANDATORY: Read Wave Plan Parallelization First

```bash
# ORCHESTRATOR MUST READ WAVE PLAN PARALLELIZATION
read_wave_parallelization_strategy() {
    local PHASE="$1"
    local WAVE="$2"
    local WAVE_PLAN="phase-plans/PHASE-${PHASE}-WAVE-${WAVE}-IMPLEMENTATION-PLAN.md"
    
    echo "═══════════════════════════════════════════════════════════════"
    echo "📖 R218: Reading Wave Plan Parallelization Strategy"
    echo "═══════════════════════════════════════════════════════════════"
    
    if [ ! -f "$WAVE_PLAN" ]; then
        echo "❌ FATAL: No wave implementation plan found!"
        echo "Cannot determine parallelization strategy!"
        exit 1
    fi
    
    # Extract parallelization groups from wave plan
    echo "Analyzing effort parallelization..."
    
    # Find blocking efforts (Can Parallelize: No)
    BLOCKING_EFFORTS=$(grep -B2 "Can Parallelize: No" "$WAVE_PLAN" | 
                       grep "### Effort" | 
                       sed 's/### Effort //')
    
    # Find parallel efforts (Can Parallelize: Yes)
    PARALLEL_EFFORTS=$(grep -B2 "Can Parallelize: Yes" "$WAVE_PLAN" | 
                       grep "### Effort" | 
                       sed 's/### Effort //')
    
    echo "BLOCKING EFFORTS (spawn sequentially):"
    echo "$BLOCKING_EFFORTS"
    echo ""
    echo "PARALLEL EFFORTS (spawn together):"
    echo "$PARALLEL_EFFORTS"
    
    # Store strategy for spawning
    cat > wave_spawn_strategy.yaml << EOF
wave_${WAVE}_spawn_strategy:
  blocking_efforts:
$(echo "$BLOCKING_EFFORTS" | sed 's/^/    - /')
  parallel_efforts:
$(echo "$PARALLEL_EFFORTS" | sed 's/^/    - /')
EOF
    
    echo "✅ Parallelization strategy loaded"
}
```

### Spawning Code Reviewers Based on Strategy

```bash
# SPAWN CODE REVIEWERS FOR EFFORT PLANNING
spawn_code_reviewers_for_effort_planning() {
    local PHASE="$1"
    local WAVE="$2"
    
    echo "═══════════════════════════════════════════════════════════════"
    echo "🚀 R218: Spawning Code Reviewers for Effort Planning"
    echo "═══════════════════════════════════════════════════════════════"
    
    # 1. Read parallelization strategy
    read_wave_parallelization_strategy "$PHASE" "$WAVE"
    
    # MANDATORY OUTPUT: Acknowledge reading the wave plan
    echo ""
    echo "📖 R218: I have READ the parallelization headers from wave plan located at:"
    echo "   $WAVE_PLAN"
    echo ""
    echo "✅ Spawning Code Reviewers as allowed by the parallelization headers in wave plan located at:"
    echo "   $WAVE_PLAN"
    
    # 2. Process BLOCKING efforts first (sequentially)
    echo "Processing BLOCKING efforts (sequential spawn)..."
    for effort in $BLOCKING_EFFORTS; do
        echo "Spawning Code Reviewer for blocking effort: $effort"
        spawn_single_code_reviewer "$PHASE" "$WAVE" "$effort"
        echo "Waiting for completion before next effort..."
        wait_for_completion "$effort"
    done
    
    # 3. Process PARALLEL efforts together (R151 compliance!)
    if [ -n "$PARALLEL_EFFORTS" ]; then
        echo "Processing PARALLEL efforts (spawning ALL together for R151)..."
        echo "📋 As allowed by the parallelization headers in wave plan located at:"
        echo "   $WAVE_PLAN"
        echo "   These efforts are marked 'Can Parallelize: Yes'"
        spawn_parallel_code_reviewers "$PHASE" "$WAVE" "$PARALLEL_EFFORTS"
    fi
}

# SEQUENTIAL SPAWN for blocking efforts
spawn_single_code_reviewer() {
    local PHASE="$1"
    local WAVE="$2"
    local EFFORT="$3"
    
    echo "📝 Spawning single Code Reviewer for: $EFFORT"
    
    # Task single Code Reviewer
    Task: code-reviewer
    PURPOSE: Create effort implementation plan
    EFFORT: $EFFORT
    WAVE_PLAN: phase-plans/PHASE-${PHASE}-WAVE-${WAVE}-IMPLEMENTATION-PLAN.md
    WORKING_DIR: efforts/phase${PHASE}/wave${WAVE}/${EFFORT}
    
    REQUIREMENTS:
    - Extract effort section from wave plan
    - PRESERVE ALL HEADERS (R211)
    - Include parallelization metadata
    - Create IMPLEMENTATION-PLAN.md
    
    # Record spawn time
    echo "spawned_at: $(date -Iseconds)" >> orchestrator-state.yaml
}

# PARALLEL SPAWN for independent efforts (R151 CRITICAL!)
spawn_parallel_code_reviewers() {
    local PHASE="$1"
    local WAVE="$2"
    local EFFORTS="$3"
    
    echo "🚀🚀🚀 R151 CRITICAL: Spawning parallel Code Reviewers in ONE message!"
    echo "Efforts to spawn in parallel:"
    echo "$EFFORTS"
    
    # Record start time for R151 grading
    START_TIME=$(date +%s%N)
    
    # CRITICAL: ALL spawns in ONE message for R151 compliance!
    # Use multiple Task tool invocations in same message
    
    for effort in $EFFORTS; do
        Task: code-reviewer
        PURPOSE: Create effort implementation plan
        EFFORT: $effort
        WAVE_PLAN: phase-plans/PHASE-${PHASE}-WAVE-${WAVE}-IMPLEMENTATION-PLAN.md
        WORKING_DIR: efforts/phase${PHASE}/wave${WAVE}/${effort}
        
        REQUIREMENTS:
        - Extract effort section from wave plan
        - PRESERVE ALL HEADERS (R211)
        - Include parallelization metadata
        - Create IMPLEMENTATION-PLAN.md
    done
    
    # Calculate spawn deltas for R151 grading
    END_TIME=$(date +%s%N)
    DELTA=$((($END_TIME - $START_TIME) / 1000000000))
    
    echo "✅ Parallel spawn complete"
    echo "R151 Grade: Average delta = ${DELTA}s (must be <5s)"
}
```

## Example Output Required from Orchestrator

When spawning Code Reviewers for effort planning, the orchestrator MUST output:

```
═══════════════════════════════════════════════════════════════
📖 R218: Reading Wave Plan Parallelization Strategy
═══════════════════════════════════════════════════════════════
📖 R218: I MUST USE THE READ TOOL on wave plan at:
   phase-plans/PHASE-1-WAVE-1-IMPLEMENTATION-PLAN.md

🚨 USING READ TOOL NOW to extract parallelization metadata...

[Orchestrator MUST use Read tool here - NOT OPTIONAL]

Analyzing effort parallelization...
BLOCKING EFFORTS (spawn sequentially):
1: Contracts & Interfaces

PARALLEL EFFORTS (spawn together for R151 compliance):
2: Feature A
3: Feature B
4: Feature C

📖 R218: I have READ the parallelization headers from wave plan located at:
   phase-plans/PHASE-1-WAVE-1-IMPLEMENTATION-PLAN.md

✅ Spawning Code Reviewers as allowed by the parallelization headers in wave plan located at:
   phase-plans/PHASE-1-WAVE-1-IMPLEMENTATION-PLAN.md

Processing PARALLEL efforts (spawning ALL together for R151)...

📋 As allowed by the parallelization headers in wave plan located at:
   phase-plans/PHASE-1-WAVE-1-IMPLEMENTATION-PLAN.md
   These efforts are marked 'Can Parallelize: Yes'
```

## Example: Reading Wave Plan and Spawning

### Wave Plan Content
```markdown
### Effort 1: Contracts & Interfaces
**Can Parallelize**: No (blocks all other efforts)
**Parallel With**: None

### Effort 2: Shared Libraries  
**Can Parallelize**: No (blocks features)
**Parallel With**: None

### Effort 3: Feature A
**Can Parallelize**: Yes
**Parallel With**: Efforts 4, 5

### Effort 4: Feature B
**Can Parallelize**: Yes
**Parallel With**: Efforts 3, 5

### Effort 5: Feature C
**Can Parallelize**: Yes
**Parallel With**: Efforts 3, 4

### Effort 6: Integration
**Can Parallelize**: No
**Parallel With**: None
```

### Resulting Spawn Strategy
```yaml
spawn_sequence:
  step_1:
    type: sequential
    effort: 1 (Contracts)
    reason: "Can Parallelize: No"
    action: Spawn single Code Reviewer, wait for completion
    
  step_2:
    type: sequential
    effort: 2 (Libraries)
    reason: "Can Parallelize: No, depends on Effort 1"
    action: Spawn single Code Reviewer, wait for completion
    
  step_3:
    type: parallel
    efforts: [3, 4, 5] (Features A, B, C)
    reason: "All marked Can Parallelize: Yes"
    action: Spawn 3 Code Reviewers in ONE message (R151)
    
  step_4:
    type: sequential
    effort: 6 (Integration)
    reason: "Can Parallelize: No, depends on 3,4,5"
    action: Spawn single Code Reviewer after 3,4,5 complete
```

## Integration with R151 (Parallel Spawn Timing)

When spawning parallel Code Reviewers:

```bash
# R151 COMPLIANCE FOR PARALLEL CODE REVIEWERS
ensure_r151_compliance() {
    echo "🚨 R151 CRITICAL: Parallel spawn timing check"
    
    # All parallel spawns MUST be in ONE message
    # Example with 3 parallel efforts:
    
    SPAWN_START=$(date +%s%N)
    
    # SINGLE MESSAGE with multiple Task invocations:
    Task: code-reviewer "Effort 3 planning"
    Task: code-reviewer "Effort 4 planning"
    Task: code-reviewer "Effort 5 planning"
    
    SPAWN_END=$(date +%s%N)
    
    # Calculate average delta
    NUM_AGENTS=3
    TOTAL_TIME=$((($SPAWN_END - $SPAWN_START) / 1000000))  # Convert to ms
    AVG_DELTA=$(($TOTAL_TIME / $NUM_AGENTS))
    
    if [ $AVG_DELTA -lt 5000 ]; then
        echo "✅ R151 PASS: Average delta ${AVG_DELTA}ms < 5000ms"
        GRADE="PASS"
    else
        echo "❌ R151 FAIL: Average delta ${AVG_DELTA}ms >= 5000ms"
        GRADE="FAIL"
    fi
    
    # Record in orchestrator-state.yaml
    cat >> orchestrator-state.yaml << EOF
parallel_code_reviewer_spawns:
  wave_${WAVE}_effort_planning:
    agents_spawned: $NUM_AGENTS
    total_time_ms: $TOTAL_TIME
    average_delta_ms: $AVG_DELTA
    r151_grade: $GRADE
EOF
}
```

## Common Violations to Avoid

### ❌ Spawning Parallel Efforts Sequentially
```bash
# WRONG - Violates parallelization strategy AND R151
spawn_code_reviewer "Effort 3"
wait_for_completion
spawn_code_reviewer "Effort 4"  # Should be parallel!
wait_for_completion
spawn_code_reviewer "Effort 5"  # Should be parallel!
```

### ❌ Spawning Blocking Efforts in Parallel
```bash
# WRONG - Violates dependencies
spawn_code_reviewer "Effort 1" &
spawn_code_reviewer "Effort 2" &  # Depends on Effort 1!
```

### ✅ Correct Implementation
```bash
# RIGHT - Respect wave plan parallelization
spawn_code_reviewer "Effort 1"  # Blocking
wait_for_completion

spawn_code_reviewer "Effort 2"  # Blocking
wait_for_completion

# Spawn 3, 4, 5 together (R151 compliance)
spawn_parallel_code_reviewers "3 4 5"  # All in ONE message!
```

## Grading Impact

- **Spawning parallel efforts sequentially**: -30% (Efficiency violation)
- **Spawning dependent efforts in parallel**: -40% (Dependency violation)  
- **R151 violation (>5s spawn delta)**: -50% (Critical performance failure)
- **Not reading wave plan parallelization**: -25% (Protocol violation)

## Integration with Other Rules

- **R211**: Code Reviewers must preserve parallelization headers
- **R151**: Parallel spawns must achieve <5s average delta
- **R053**: Parallelization decision guidelines
- **R052**: Agent spawning protocol

## Summary

**R218 ensures efficient Code Reviewer spawning by:**
1. Reading wave plan parallelization metadata
2. Spawning blocking efforts sequentially with dependencies
3. Spawning parallel efforts together in ONE message (R151)
4. Recording spawn timing for grading
5. Optimizing overall wave planning time