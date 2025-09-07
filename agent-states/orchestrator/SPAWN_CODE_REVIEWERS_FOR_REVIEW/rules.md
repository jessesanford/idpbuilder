# Orchestrator - SPAWN_CODE_REVIEWERS_FOR_REVIEW State Rules

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

### STOP PROTOCOL:
```markdown
## 🛑 STATE TRANSITION CHECKPOINT: CURRENT_STATE → NEXT_STATE

### ✅ Current State Work Completed:
- [List completed work]

### 📊 Current Status:
- Current State: CURRENT_STATE
- Next State: NEXT_STATE
- TODOs Completed: X/Y
- State Files: Updated and committed ✅

### ⏸️ STOPPED - Awaiting User Continuation
Ready to transition to NEXT_STATE. Please use /continue-orchestrating.
```

**STOP MEANS STOP - Exit and wait for /continue-orchestrating**

---

## 🔴🔴🔴 R322 MANDATORY: STOP BEFORE STATE TRANSITIONS 🔴🔴🔴

**CRITICAL REQUIREMENT PER R322:**
After spawning ANY agents in this state, you MUST:
1. Record what was spawned in state file
2. Save TODOs per R287
3. Commit and push state changes
4. Display stop message with continuation instructions
5. EXIT immediately with code 0

**VIOLATION = AUTOMATIC -100% FAILURE**

See: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-after-spawn.md`

---

# Orchestrator - SPAWN_CODE_REVIEWERS_FOR_REVIEW State Rules

## 🔴🔴🔴 R290 ENFORCEMENT: READ THESE RULES FIRST! 🔴🔴🔴

**SUPREME LAW #3 (R290): STATE RULES MUST BE READ BEFORE STATE ACTIONS**

## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED SPAWN_CODE_REVIEWERS_FOR_REVIEW STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_SPAWN_CODE_REVIEWERS_FOR_REVIEW
echo "$(date +%s) - Rules read and acknowledged for SPAWN_CODE_REVIEWERS_FOR_REVIEW" > .state_rules_read_orchestrator_SPAWN_CODE_REVIEWERS_FOR_REVIEW
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY SPAWN_CODE_REVIEWERS_FOR_REVIEW WORK UNTIL RULES ARE READ:
- ❌ Start spawn code reviewer agents
- ❌ Start assign review work
- ❌ Start distribute review tasks
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
   ❌ WRONG: "I acknowledge all SPAWN_CODE_REVIEWERS_FOR_REVIEW rules"
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

### ✅ CORRECT PATTERN FOR SPAWN_CODE_REVIEWERS_FOR_REVIEW:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/[first-rule-file].md
2. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/[second-rule-file].md  
4. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. "Ready to execute SPAWN_CODE_REVIEWERS_FOR_REVIEW work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY SPAWN_CODE_REVIEWERS_FOR_REVIEW work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ You have stated readiness to execute SPAWN_CODE_REVIEWERS_FOR_REVIEW work
1. **READ** every rule file listed in PRIMARY DIRECTIVES below
2. **ACKNOWLEDGE** each rule individually with number and description
3. **ONLY THEN** proceed with SPAWN_CODE_REVIEWERS_FOR_REVIEW work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY SPAWN_CODE_REVIEWERS_FOR_REVIEW work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

**THE SYSTEM IS MONITORING YOUR READ TOOL CALLS!**

## 📋 PRIMARY DIRECTIVES FOR SPAWN_CODE_REVIEWERS_FOR_REVIEW

**YOU MUST READ EACH RULE LISTED HERE. YOUR READ TOOL CALLS ARE BEING MONITORED.**

## 🔴🔴🔴 MANDATORY MEASUREMENT RULE - R304 🔴🔴🔴

**ABSOLUTE REQUIREMENTS:**
- ✅ MUST use `$CLAUDE_PROJECT_DIR/tools/line-counter.sh` for ALL line counting
- ❌ NEVER use `wc -l` or any manual counting method
- ❌ NEVER count lines any other way - this is a -100% automatic failure
- ✅ MUST specify both -b (base branch) and -c (current branch) parameters
- ✅ Base branch MUST be phase integration branch (NOT "main")

**Failure to use the line counter tool = AUTOMATIC -100% GRADE**

### State-Specific Rules (NOT in orchestrator.md):

1. **🚨🚨🚨 R151** - Parallel Spawning Timestamp Requirement
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R151-parallel-agent-spawning-timing.md`
   - Criticality: CRITICAL - Timestamps must be within 5s for parallel agents
   - Summary: All parallel agents must emit timestamps within 5 seconds

2. **R108** - Code Review Protocol
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R108-code-review-protocol.md`
   - Criticality: BLOCKING - Complete review protocol
   - Summary: Review for size limits, quality, patterns, and create reports

3. **R222** - Code Review Gate
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R222-code-review-gate.md`
   - Criticality: BLOCKING - Must create standardized reports
   - Summary: Generate CODE-REVIEW-REPORT.md for all findings

4. **R255** - Post-Agent Work Verification
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R255-post-agent-work-verification.md`
   - Criticality: BLOCKING - Verify correct locations after completion
   - Summary: Ensure all review work is in correct directories

5. **🔴🔴🔴 R208** - CD Before Agent Spawn (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R208-cd-before-agent-spawn.md`
   - Criticality: SUPREME LAW - CD to correct directory before spawn
   - Summary: Must change to agent's working directory before spawning

**Note**: Additional TODO persistence rules (R287) apply from orchestrator.md.

## 🚨 SPAWN_CODE_REVIEWERS_FOR_REVIEW IS A VERB - START SPAWNING IMMEDIATELY! 🚨

**See R151 for immediate action requirements when entering this state.**

The SPAWN_CODE_REVIEWERS_FOR_REVIEW state requires IMMEDIATE spawning action - no pausing or waiting.

## State Context
You are spawning Code Reviewer agents to review:
- Completed implementation work from SW Engineers
- **INTEGRATION FIXES** when coming from MONITORING_FIX_PROGRESS
- Any code that needs review before proceeding

## 🔴🔴🔴 PREREQUISITES FOR SPAWN_CODE_REVIEWERS_FOR_REVIEW 🔴🔴🔴

**BEFORE ENTERING THIS STATE, YOU MUST ALREADY HAVE:**
1. ✅ All SW Engineers completed their implementation
2. ✅ All code committed and pushed by SW Engineers
3. ✅ Size measurements completed and within limits
4. ✅ All effort directories contain implemented code
5. ✅ **PARALLELIZATION ANALYSIS COMPLETE (ANALYZE_CODE_REVIEWER_PARALLELIZATION)**
6. ✅ **Code Reviewer parallelization plan in orchestrator-state.yaml**

**IF PARALLELIZATION NOT ANALYZED, GO BACK TO ANALYZE_CODE_REVIEWER_PARALLELIZATION!**

## Review Assignment Protocol

### For Each Code Reviewer to Spawn:
1. **CD to effort directory** (R208 SUPREME LAW)
2. **Spawn with clear review scope**:
   - Which efforts to review
   - What type of review (code quality, size, patterns)
   - Where to create reports
3. **Track spawn timestamps** (R151 requirement)
4. **Monitor completion** via orchestrator-state.yaml

### Parallel vs Sequential Spawning:
- **Parallel**: When reviewing independent efforts
- **Sequential**: When reviewing split efforts or dependencies
- **Decision made in**: ANALYZE_CODE_REVIEWER_PARALLELIZATION state

## Expected Deliverables

Each Code Reviewer must produce:
1. **CODE-REVIEW-REPORT.md** in effort directory
2. **Size compliance verification**
3. **Quality assessment**
4. **Recommendations for fixes if needed**

## 🚨🚨🚨 SPECIAL CASE: REVIEWING INTEGRATION FIXES 🚨🚨🚨

**When coming from MONITORING_FIX_PROGRESS:**
1. You are reviewing FIXES to integration issues
2. Focus review on:
   - Did the fixes resolve the integration problems?
   - Are the fixes properly implemented?
   - Do the fixes maintain code quality?
3. After successful review of fixes:
   - Transition to MONITOR state
   - Then to WAVE_COMPLETE
   - Then BACK TO INTEGRATION for full re-run
4. **NEVER skip directly to MONITORING_INTEGRATION**

## State Transitions

From SPAWN_CODE_REVIEWERS_FOR_REVIEW:
- **REVIEWERS_COMPLETE** → WAVE_COMPLETE (All reviews done)
- **REVIEWS_FAILED** → ERROR_RECOVERY (Critical issues found)
- **NEED_FIXES** → SPAWN_AGENTS (Re-spawn SW Engineers for fixes)
- **When reviewing fixes** → MONITOR → WAVE_COMPLETE → INTEGRATION (re-run)

## Critical Rule Enforcement Order

1. **R290**: Read these state rules FIRST and create verification marker
3. **R208**: CD to correct directory BEFORE spawning
4. **R151**: Ensure parallel timestamps within 5s
5. **R053**: Follow complete review protocol
6. **R054**: Generate standardized reports
7. **R255**: Verify work in correct locations

**Remember**: Code Reviewers check for size violations, quality issues, and create actionable reports for the orchestrator.

## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-orchestrating
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**
