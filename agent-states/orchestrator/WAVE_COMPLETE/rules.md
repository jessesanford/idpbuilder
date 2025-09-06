# Orchestrator - WAVE_COMPLETE State Rules

## 🛑🛑🛑 R322 MANDATORY STOP BEFORE STATE TRANSITIONS 🛑🛑🛑

**SUPREME LAW - VIOLATION = -100% IMMEDIATE FAILURE**

### YOU MUST STOP AFTER:
1. ✅ Completing all TODOs for this state
2. ✅ Updating orchestrator-state.yaml with new state
3. ✅ Committing and pushing the state file  
4. ✅ Providing work summary

### YOU MUST NOT:
- ❌ Continue to the next state automatically
- ❌ Start work for the new state
- ❌ Spawn agents for the new state
- ❌ Assume permission to continue

### STOP PROTOCOL FOR WAVE_COMPLETE:
```markdown
## 🛑 STATE TRANSITION CHECKPOINT: WAVE_COMPLETE → INTEGRATION

### ✅ Current State Work Completed:
- All efforts verified complete with passed reviews
- Wave marked complete in orchestrator-state.yaml
- current_state updated to "INTEGRATION"
- State file committed and pushed

### 📊 Current Status:
- Current State: INTEGRATION (already updated in file)
- Previous State: WAVE_COMPLETE
- Wave Status: Complete and validated
- State Files: Updated to INTEGRATION and committed ✅

### ⏸️ STOPPED - Awaiting User Continuation
State file already updated to INTEGRATION. When you run /continue-orchestrating, 
integration work will begin immediately. Please use /continue-orchestrating.
```

**STOP MEANS STOP - Exit and wait for /continue-orchestrating**

---


## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED WAVE_COMPLETE STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_WAVE_COMPLETE
echo "$(date +%s) - Rules read and acknowledged for WAVE_COMPLETE" > .state_rules_read_orchestrator_WAVE_COMPLETE
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY WAVE_COMPLETE WORK UNTIL RULES ARE READ:
- ❌ Start finalize wave efforts
- ❌ Start collect implementation results
- ❌ Start prepare integration
- ❌ Update state files
- ❌ Continue to next state
- ❌ Think about what to do in this state

### ✅ YOU MUST IMMEDIATELY:

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

### ❌ ANTI-PATTERNS THAT WILL CAUSE FAILURE:

1. **Fake Acknowledgment Without Reading**:
   ```
   ❌ WRONG: "I acknowledge R151, R208, R053..."
   (No Read tool calls detected = AUTOMATIC FAILURE)
   ```

2. **Bulk Acknowledgment**:
   ```
   ❌ WRONG: "I acknowledge all WAVE_COMPLETE rules"
   (YOU Must READ AND ACKNOWLEDGE EACH rule individually)
   ```

3. **Silent Reading**:
   ```
   ❌ WRONG: [Reads rules but doesn't acknowledge]
   "Now I've read the rules, let me start work..."
   (MUST explicitly acknowledge EACH rule)
   ```

4. **Reading From Memory**:
   ```
   ❌ WRONG: "I know R208 requires CD before spawn..."
   (Must READ from file, not recall from memory)
   ```

5. **Skipping Rules in PRIMARY DIRECTIVES**:
   ```
   ❌ WRONG: Reading only some rules from the list
   (ALL rules in PRIMARY DIRECTIVES are MANDATORY)
   ```

### ✅ CORRECT PATTERN FOR WAVE_COMPLETE:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/[first-rule-file].md
2. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/[second-rule-file].md  
4. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. "Ready to execute WAVE_COMPLETE work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY WAVE_COMPLETE work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ You have stated readiness to execute WAVE_COMPLETE work
1. **READ** every rule file listed in PRIMARY DIRECTIVES below
2. **ACKNOWLEDGE** each rule individually with number and description
3. **ONLY THEN** proceed with WAVE_COMPLETE work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY WAVE_COMPLETE work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

**THE SYSTEM IS MONITORING YOUR READ TOOL CALLS!**

## 📋 PRIMARY DIRECTIVES FOR WAVE_COMPLETE

**YOU MUST READ EACH RULE LISTED HERE. YOUR READ TOOL CALLS ARE BEING MONITORED.**

### State-Specific Rules (NOT in orchestrator.md):
1. **R222** - Code Review Gate
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R222-code-review-gate.md`
   - Criticality: BLOCKING - ALL reviews must pass before proceeding

2. **R105** - Wave Completion Protocol
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R105-wave-completion-protocol.md`
   - Criticality: INFO - Best practices for wave completion

3. **R035** - Phase Completion Testing
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R035-phase-completion-testing.md`
   - Criticality: MANDATORY - Validation requirements

**Note**: R288 (state updates), R287 (TODO saves) are already in orchestrator.md.

## 📋 RULE SUMMARY FOR WAVE_COMPLETE STATE

### Rules Enforced in This State:
- R222: Code Review Gate [BLOCKING - ALL reviews must pass]
- R288: State File Update and Commit [SUPREME LAW - Update immediately, includes commit/push]
- R105: Wave Completion Protocol [INFO - Best practice]
- R035: Phase Completion Testing [MANDATORY - Validation required]
- R287: TODO Save Triggers [BLOCKING - Save on completion]

### Critical Requirements:
1. Verify ALL reviews passed (R222) - Penalty: -100%
2. Verify ALL size compliance (<800 lines) - Penalty: -100%
3. Update state file immediately - Penalty: -50%
4. Save and commit TODOs - Penalty: -20%
5. Prepare for INTEGRATION state transition - Penalty: -10%

### Success Criteria:
- ✅ ALL efforts have PASSED code reviews
- ✅ ALL efforts are size compliant
- ✅ State file updated with completion
- ✅ Ready to transition to INTEGRATION state
- ✅ TODOs saved and committed

### Failure Triggers:
- ❌ Enter with ANY failed review = R222 VIOLATION
- ❌ Enter with size violations = AUTOMATIC FAILURE
- ❌ Skip state file update = R288 VIOLATION
- ❌ Forget TODO saves = -20% penalty

## 🚨 WAVE_COMPLETE IS A VERB - START WAVE COMPLETION PROCESS IMMEDIATELY! 🚨

### IMMEDIATE ACTIONS UPON ENTERING WAVE_COMPLETE

**THE MOMENT YOU ENTER THIS STATE, YOU MUST:**
1. Verify all efforts are complete and passed reviews
2. Update orchestrator-state.yaml with completion
3. Check TodoWrite for pending items and process them
4. Prepare to transition to INTEGRATION state

**FORBIDDEN - AUTOMATIC FAILURE:**
- ❌ "STATE TRANSITION COMPLETE: Now in WAVE_COMPLETE" [stops]
- ❌ "Successfully entered WAVE_COMPLETE state" [waits]
- ❌ "Ready to start wave completion process" [pauses]
- ❌ "I'm in WAVE_COMPLETE state" [does nothing]

**REQUIRED - IMMEDIATE ACTION:**
- ✅ "Entering WAVE_COMPLETE, verifying all efforts complete..."
- ✅ "START WAVE COMPLETION PROCESS, update orchestrator-state.yaml with completion..."
- ✅ "WAVE_COMPLETE: Validating all reviews passed before proceeding..."

## State Context
You have completed all efforts in a wave and need to validate completion and prepare for next steps.

**NOTE: Integration branch creation happens in the INTEGRATION state, not here!**
- This state validates that all efforts are complete
- The INTEGRATION state handles creating integration branches
- See state machine: WAVE_COMPLETE → INTEGRATION

## 🔴🔴🔴 R222 CODE REVIEW GATE - ABSOLUTELY MANDATORY 🔴🔴🔴

**YOU CANNOT BE IN THIS STATE UNLESS:**
1. **ALL** Code Reviews have been run
2. **ALL** Code Reviews have PASSED
3. **ALL** Size compliance checks PASSED (<800 lines)
4. **NO** effort is in FIX_ISSUES state
5. **NO** effort has pending review issues

### MANDATORY VERIFICATION BEFORE PROCEEDING:
```bash
# R222 ENFORCEMENT - MUST CHECK EVERY EFFORT!
echo "🔍 R222: Verifying ALL reviews passed..."
ALL_PASSED=true

for effort in $WAVE_EFFORTS; do
    REVIEW_STATUS=$(check_effort_review_status "$effort")
    SIZE_STATUS=$(check_effort_size_compliance "$effort")
    
    if [ "$REVIEW_STATUS" != "PASSED" ]; then
        echo "❌ R222 VIOLATION: $effort review status: $REVIEW_STATUS"
        echo "🚫 CANNOT BE IN WAVE_COMPLETE STATE!"
        echo "🔄 Must return to MONITOR and fix issues"
        ALL_PASSED=false
    fi
    
    if [ "$SIZE_STATUS" != "COMPLIANT" ]; then
        echo "❌ SIZE VIOLATION: $effort exceeds 800 lines!"
        echo "🚫 CANNOT BE IN WAVE_COMPLETE STATE!"
        ALL_PASSED=false
    fi
done

if [ "$ALL_PASSED" = false ]; then
    echo "🔴🔴🔴 CRITICAL: INVALID STATE TRANSITION!"
    echo "Must return to MONITOR and execute review-fix loops"
    exit 222
fi

echo "✅ R222 VERIFIED: All reviews passed, proceeding with completion"
```

## 🔴🔴🔴 CRITICAL: MANDATORY STATE FILE UPDATE (R288) 🔴🔴🔴

### IMMEDIATELY upon entering WAVE_COMPLETE state:

```bash
# 1. Update state machine
update_orchestrator_state "WAVE_COMPLETE" "All efforts reviewed and passed"

# 2. Mark wave as complete in state file
mark_wave_complete "$PHASE" "$WAVE"

# Example state file update:
yq -i ".waves_completed.phase${PHASE}.wave${WAVE}.completed_at = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" orchestrator-state.yaml
yq -i ".waves_completed.phase${PHASE}.wave${WAVE}.status = \"COMPLETE\"" orchestrator-state.yaml
yq -i ".waves_completed.phase${PHASE}.wave${WAVE}.efforts_count = $EFFORT_COUNT" orchestrator-state.yaml
yq -i ".waves_completed.phase${PHASE}.wave${WAVE}.all_reviews_passed = true" orchestrator-state.yaml
yq -i ".waves_completed.phase${PHASE}.wave${WAVE}.size_compliant = true" orchestrator-state.yaml
yq -i ".waves_completed.phase${PHASE}.wave${WAVE}.integration_branch = \"$INTEGRATION_BRANCH\"" orchestrator-state.yaml

# 3. 🔴🔴🔴 CRITICAL: UPDATE STATE TO INTEGRATION BEFORE STOPPING! 🔴🔴🔴
# Per state machine: WAVE_COMPLETE → INTEGRATION is the REQUIRED transition
echo "📝 Updating current_state to INTEGRATION for next continuation..."
yq -i ".current_state = \"INTEGRATION\"" orchestrator-state.yaml
yq -i ".previous_state = \"WAVE_COMPLETE\"" orchestrator-state.yaml
yq -i ".transition_time = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" orchestrator-state.yaml
yq -i ".transition_reason = \"Wave $WAVE complete, all reviews passed, ready for integration\"" orchestrator-state.yaml
echo "✅ State updated to INTEGRATION - will execute integration work on next /continue-orchestrating"

# R301 MANDATORY: Update current_wave_integration
echo "📝 Updating current_wave_integration per R301..."

# Deprecate existing wave integration if it exists
EXISTING_WAVE=$(yq ".current_wave_integration | select(.phase == $PHASE and .wave == $WAVE)" orchestrator-state.yaml)
if [ ! -z "$EXISTING_WAVE" ]; then
    yq -i ".deprecated_wave_integrations += (.current_wave_integration | select(.phase == $PHASE and .wave == $WAVE))" orchestrator-state.yaml
fi

# Set the NEW current wave integration
yq -i ".current_wave_integration = {
  \"phase\": $PHASE,
  \"wave\": $WAVE,
  \"branch\": \"$INTEGRATION_BRANCH\",
  \"status\": \"active\",
  \"created_at\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
  \"type\": \"initial\"
}" orchestrator-state.yaml

echo "✅ Current wave integration updated per R301"
```

**VIOLATION = AUTOMATIC FAILURE if wave completion not recorded in state file!**
**VIOLATION = AUTOMATIC FAILURE if current_wave_integration not updated per R301!**

### R287-R287 TODO PERSISTENCE ON COMPLETION
```bash
# R287: Major milestone trigger - save TODOs
echo "💾 R287: Wave complete - saving TODOs..."
save_todos "WAVE_COMPLETE - Phase $PHASE Wave $WAVE"

# R287: Commit within 60 seconds
cd $CLAUDE_PROJECT_DIR
git add todos/*.todo orchestrator-state.yaml
git commit -m "todo: wave $WAVE complete, all reviews passed"
git push

echo "✅ Wave completion and TODOs persisted"
```


## Wave Completion Validation


## Effort Completion Verification

```python
def verify_effort_completion(effort_id, effort_data):
    """Verify an individual effort is truly complete"""
    
    working_dir = effort_data['working_dir']
    branch = effort_data['branch']
    
    verification_results = {
        'effort_id': effort_id,
        'checks': {},
        'overall_complete': True,
        'issues': []
    }
    
    # 1. Implementation plan completion check
    plan_check = verify_implementation_plan_complete(working_dir)
    verification_results['checks']['plan_complete'] = plan_check
    if not plan_check['complete']:
        verification_results['overall_complete'] = False
        verification_results['issues'].extend(plan_check['missing_items'])
    
    # 2. Test coverage verification
    test_check = verify_test_coverage(working_dir)
    verification_results['checks']['test_coverage'] = test_check
    if not test_check['meets_requirements']:
        verification_results['overall_complete'] = False
        verification_results['issues'].append(f"Test coverage {test_check['coverage']}% < required {test_check['required']}%")
    
    # 3. Size compliance check (CRITICAL)
    size_check = verify_size_compliance(branch)
    verification_results['checks']['size_compliance'] = size_check
    if not size_check['compliant']:
        verification_results['overall_complete'] = False
        verification_results['issues'].append(f"Size {size_check['lines']} > limit {size_check['limit']}")
    
    # 4. Git status verification
    git_check = verify_git_status(working_dir, branch)
    verification_results['checks']['git_status'] = git_check
    if not git_check['clean']:
        verification_results['overall_complete'] = False
        verification_results['issues'].extend(git_check['issues'])
    
    # 5. Work log completeness
    worklog_check = verify_work_log_complete(working_dir)
    verification_results['checks']['work_log'] = worklog_check
    if not worklog_check['complete']:
        verification_results['overall_complete'] = False
        verification_results['issues'].append("Work log incomplete or missing final entries")
    
    return verification_results

def verify_size_from_review_report(effort_name):
    """🔴 R006: Orchestrator NEVER measures - read Code Reviewer reports"""
    
    # ORCHESTRATOR MUST NEVER RUN line-counter.sh!
    # Instead, read the review report created by Code Reviewer
    review_report = f"efforts/phase{PHASE}/wave{WAVE}/{effort_name}/CODE-REVIEW-REPORT.md"
    
    if not os.path.exists(review_report):
        return {
            'compliant': False,
            'error': "No Code Reviewer report found - spawn Code Reviewer first!",
            'lines': 'unknown',
            'limit': 800
        }
    
    # Parse the review report for size information
    with open(review_report, 'r') as f:
        content = f.read()
        # Look for size measurement in report
        if 'Size:' in content or 'lines' in content.lower():
            import re
            lines_match = re.search(r'(\d+)\s*lines', content, re.IGNORECASE)
            if lines_match:
                lines = int(lines_match.group(1))
                return {
                    'compliant': lines <= 800,
                    'lines': lines,
                    'limit': 800,
                    'source': 'Code Reviewer Report'
                }
    
    return {
        'compliant': False,
        'error': "Size not found in review report",
        'lines': 'unknown',
        'limit': 800
    }
```

## Wave-Level Integration Validation


```bash
#!/bin/bash
# Wave integration script

PHASE=$1
WAVE=$2
INTEGRATION_BRANCH="phase${PHASE}/wave${WAVE}-integration"

echo "🔗 Creating wave integration branch: $INTEGRATION_BRANCH"

# Create integration branch
git checkout -b "$INTEGRATION_BRANCH"

# Get all effort branches for this wave
EFFORT_BRANCHES=$(git branch -r | grep "phase${PHASE}/wave${WAVE}/effort" | sed 's/origin\///')

for EFFORT_BRANCH in $EFFORT_BRANCHES; do
    echo "Merging effort branch: $EFFORT_BRANCH"
    
    # Merge with no-fast-forward to maintain history
    git merge "origin/$EFFORT_BRANCH" --no-ff -m "integrate: $EFFORT_BRANCH into wave"
    
    if [ $? -ne 0 ]; then
        echo "❌ Conflict detected merging $EFFORT_BRANCH"
        echo "Manual resolution required before continuing"
        exit 1
    fi
    
    echo "✅ Successfully merged $EFFORT_BRANCH"
done

# 🔴 R006: ORCHESTRATOR NEVER MEASURES - Spawn Code Reviewer!
echo "📏 Wave integration needs size verification..."
echo "🚀 Spawning Code Reviewer to measure integrated wave..."
# Code Reviewer will use line-counter.sh to verify size
# Orchestrator will read the review report afterward
# DO NOT run line-counter.sh directly - that's Code Reviewer work!

# Run tests on integration
echo "🧪 Running tests on integrated wave..."
make test
if [ $? -ne 0 ]; then
    echo "❌ Tests failed on wave integration"
    exit 1
fi

echo "✅ Wave integration complete and validated"
```

## Architect Review Decision


```python
def should_request_architect_review(wave_data):
    """Determine if architect review is needed for this wave"""
    
    review_decision = {
        'review_required': False,
        'review_type': None,
        'reasons': [],
        'urgency': 'NORMAL'
    }
    
    # Check mandatory review triggers
    if wave_data.get('phase_end', False):
        review_decision['review_required'] = True
        review_decision['review_type'] = 'PHASE_COMPLETION'
        review_decision['reasons'].append('End of phase - mandatory review')
        review_decision['urgency'] = 'HIGH'
    
    # Check for size violations in any effort
    size_violations = []
    for effort in wave_data.get('efforts', []):
        if not effort.get('size_compliant', True):
            size_violations.append(effort['id'])
    
    if size_violations:
        review_decision['review_required'] = True
        review_decision['review_type'] = 'SIZE_VIOLATION_REVIEW'
        review_decision['reasons'].append(f'Size violations in efforts: {size_violations}')
        review_decision['urgency'] = 'CRITICAL'
    
    # Check for architecture violations
    arch_violations = wave_data.get('architecture_violations', [])
    if arch_violations:
        review_decision['review_required'] = True
        review_decision['review_type'] = 'ARCHITECTURE_VIOLATION'
        review_decision['reasons'].extend(arch_violations)
        review_decision['urgency'] = 'HIGH'
    
    # Check optional review triggers
    if not review_decision['review_required']:
        # Complex integration (>4 efforts)
        if len(wave_data.get('efforts', [])) > 4:
            review_decision['review_required'] = True
            review_decision['review_type'] = 'COMPLEXITY_REVIEW'
            review_decision['reasons'].append('Complex wave with many efforts')
            review_decision['urgency'] = 'NORMAL'
        
        # New patterns detected
        if wave_data.get('new_patterns_introduced', False):
            review_decision['review_required'] = True
            review_decision['review_type'] = 'PATTERN_REVIEW'
            review_decision['reasons'].append('New architectural patterns introduced')
            review_decision['urgency'] = 'NORMAL'
    
    return review_decision
```

## Next Wave Planning

```yaml
# Wave completion analysis for planning
wave_completion_analysis:
  completed_wave:
    phase: 1
    wave: 2
    completed_at: "2025-08-23T17:30:00Z"
    efforts_completed: 4
    total_lines_delivered: 2847
    
  performance_metrics:
    wave_duration_hours: 6.5
    average_effort_size: 711
    integration_conflicts: 1
    test_pass_rate: 100
    
  lessons_learned:
    - "API types effort completed ahead of schedule"
    - "Controller effort had minor integration conflict with webhooks"
    - "Size management worked well with 700-line target per effort"
    
  next_wave_recommendations:
    - "Continue with 4 effort pattern"
    - "Monitor controller-webhook dependencies closely"
    - "Consider pre-integration dependency analysis"
    
  readiness_for_next:
    dependencies_resolved: true
    blockers_identified: []
    resource_availability: "FULL"
    estimated_start: "2025-08-23T18:00:00Z"
```

## State Transition Decision Matrix

```python
def determine_next_state(wave_completion_data):
    """Determine next state after wave completion"""
    
    # Check if this completes the current phase
    if wave_completion_data.get('phase_complete', False):
        return {
            'next_state': 'WAVE_REVIEW',
            'reason': 'Phase complete - architect phase review required',
            'data': {
                'review_type': 'PHASE_COMPLETION',
                'phase': wave_completion_data['phase']
            }
        }
    
    # Check if architect review is required for other reasons
    review_decision = should_request_architect_review(wave_completion_data)
    if review_decision['review_required']:
        return {
            'next_state': 'WAVE_REVIEW',
            'reason': f'Architect review required: {review_decision["review_type"]}',
            'data': review_decision
        }
    
    # Check if integration issues need resolution
    if wave_completion_data.get('integration_issues', []):
        return {
            'next_state': 'INTEGRATION',
            'reason': 'Integration issues require resolution',
            'data': {
                'issues': wave_completion_data['integration_issues']
            }
        }
    
    # Check if any efforts need splits due to size
    efforts_needing_splits = [
        effort for effort in wave_completion_data.get('efforts', [])
        if not effort.get('size_compliant', True)
    ]
    
    if efforts_needing_splits:
        return {
            'next_state': 'SPAWN_AGENTS',
            'reason': 'Efforts need splitting due to size violations',
            'data': {
                'spawn_type': 'CODE_REVIEWER_SPLITS',
                'efforts_to_split': efforts_needing_splits
            }
        }
    
    # Check if ready for next wave
    next_wave_ready = wave_completion_data.get('next_wave_ready', True)
    if next_wave_ready:
        return {
            'next_state': 'WAVE_START',
            'reason': 'Ready to start next wave',
            'data': {
                'next_phase': wave_completion_data['phase'],
                'next_wave': wave_completion_data['wave'] + 1
            }
        }
    
    # Need to wait or investigate issues
    return {
        'next_state': 'MONITOR',
        'reason': 'Monitoring remaining completion tasks',
        'data': {
            'pending_tasks': wave_completion_data.get('pending_tasks', [])
        }
    }
```

## Wave Completion Reporting

```python
def generate_wave_completion_report(wave_data):
    """Generate comprehensive wave completion report"""
    
    report = {
        'wave_id': f"phase{wave_data['phase']}_wave{wave_data['wave']}",
        'completion_timestamp': datetime.now().isoformat(),
        'summary': generate_wave_summary(wave_data),
        'effort_details': generate_effort_reports(wave_data['efforts']),
        'integration_results': wave_data.get('integration_results', {}),
        'performance_metrics': calculate_wave_performance_metrics(wave_data),
        'quality_gates': validate_wave_quality_gates(wave_data),
        'recommendations': generate_next_wave_recommendations(wave_data)
    }
    
    print("📋 WAVE COMPLETION REPORT")
    print(f"Wave: {report['wave_id']}")
    print(f"Efforts Completed: {len(report['effort_details'])}")
    print(f"Total Lines Delivered: {report['performance_metrics']['total_lines']}")
    print(f"Quality Gates: {'✅ PASSED' if report['quality_gates']['all_passed'] else '❌ ISSUES'}")
    
    if not report['quality_gates']['all_passed']:
        print("Quality Gate Issues:")
        for issue in report['quality_gates']['failed_gates']:
            print(f"  - {issue}")
    
    return report

def validate_wave_quality_gates(wave_data):
    """Validate all quality gates for wave completion"""
    
    quality_gates = {
        'all_efforts_complete': True,
        'size_compliance': True,
        'test_coverage': True,
        'integration_clean': True,
        'performance_acceptable': True
    }
    
    failed_gates = []
    
    # Check each effort
    for effort in wave_data.get('efforts', []):
        if not effort.get('complete', False):
            quality_gates['all_efforts_complete'] = False
            failed_gates.append(f"Effort {effort['id']} not complete")
        
        if not effort.get('size_compliant', True):
            quality_gates['size_compliance'] = False
            failed_gates.append(f"Effort {effort['id']} exceeds size limits")
        
        if effort.get('test_coverage', 100) < effort.get('required_coverage', 80):
            quality_gates['test_coverage'] = False
            failed_gates.append(f"Effort {effort['id']} below test coverage requirements")
    
    # Check integration results
    integration = wave_data.get('integration_results', {})
    if integration.get('conflicts', 0) > 0:
        quality_gates['integration_clean'] = False
        failed_gates.append(f"Integration conflicts detected: {integration['conflicts']}")
    
    # Check performance metrics
    perf = wave_data.get('performance_metrics', {})
    if perf.get('build_time_minutes', 0) > 10:  # Build should complete in <10 min
        quality_gates['performance_acceptable'] = False
        failed_gates.append(f"Build time too long: {perf['build_time_minutes']} minutes")
    
    return {
        'quality_gates': quality_gates,
        'all_passed': all(quality_gates.values()),
        'failed_gates': failed_gates,
        'pass_count': sum(quality_gates.values()),
        'total_gates': len(quality_gates)
    }
```

## State Transitions

### 🔴🔴🔴 CRITICAL: DEFAULT TRANSITION IS TO INTEGRATION! 🔴🔴🔴

**UNLESS OTHERWISE DETERMINED, WAVE_COMPLETE ALWAYS TRANSITIONS TO INTEGRATION:**
1. Update current_state to "INTEGRATION" in orchestrator-state.yaml
2. Commit and push the state file
3. STOP per R322 and wait for /continue-orchestrating
4. When user continues, orchestrator will be in INTEGRATION state and execute integration work

From WAVE_COMPLETE state, the STANDARD transition is:
- **DEFAULT** → **INTEGRATION** (Always go here unless special conditions below)

Special condition transitions (RARE):
- **ARCHITECT_REVIEW_REQUIRED** → WAVE_REVIEW (Only if architect review explicitly needed)
- **SPLITS_REQUIRED** → SPAWN_AGENTS (Only if size violations detected)
- **NEXT_WAVE_READY** → WAVE_START (Only if skipping integration - VERY RARE)
- **PHASE_COMPLETE** → WAVE_REVIEW (Only at end of phase)

## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-orchestrating
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**
