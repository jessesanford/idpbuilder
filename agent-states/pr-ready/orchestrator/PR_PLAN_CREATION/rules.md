# Orchestrator - PR_PLAN_CREATION State Rules

## 🔴🔴🔴 ABSOLUTE REQUIREMENT: STATE MANAGER CONSULTATION 🔴🔴🔴

**SUPREME LAW - R517 - UNIVERSAL STATE MANAGER CONSULTATION LAW**

**BEFORE EXITING THIS STATE, YOU MUST:**

1. **SPAWN STATE MANAGER FOR SHUTDOWN_CONSULTATION** (MANDATORY - NO EXCEPTIONS)
2. **NEVER UPDATE STATE FILES DIRECTLY** (orchestrator-state-v3.json, bug-tracking.json, etc.)
3. **NEVER COMMIT STATE FILES YOURSELF** (State Manager does this atomically)

**FAILURE TO CONSULT STATE MANAGER = IMMEDIATE SYSTEM HALT (-100% GRADE)**

### Enforcement Mechanism

If you attempt to exit this state without spawning State Manager:
- ❌ Pre-commit hooks will REJECT your commit
- ❌ Validation tools will FAIL the build
- ❌ Grading system will assign -100% penalty
- ❌ System will transition to ERROR_RECOVERY

### Required Pattern (COPY THIS EXACTLY)

```bash
# At end of state work, BEFORE any state file updates:

echo "🔴 MANDATORY: Spawning State Manager for SHUTDOWN_CONSULTATION"

# Spawn State Manager (REQUIRED - NOT OPTIONAL)
# Task: state-manager
# State: SHUTDOWN_CONSULTATION
# Current State: [YOUR_CURRENT_STATE]
# Proposed Next State: [YOUR_PROPOSED_NEXT_STATE]
# Work Summary: [SUMMARY_OF_WORK_COMPLETED]

# State Manager will:
# 1. Validate proposed transition against state machine
# 2. Update all 4 state files atomically
# 3. Commit with [R288] tag
# 4. Return REQUIRED next state (may differ from proposal)

# Wait for State Manager response
# Follow State Manager's directive (REQUIRED next state)
# DO NOT proceed until State Manager confirms
```

**YOU MUST NEVER:**
- ❌ Update orchestrator-state-v3.json yourself
- ❌ Update bug-tracking.json yourself
- ❌ Update integration-containers.json yourself
- ❌ Use `jq` to modify state files
- ❌ Use `sed/awk` to modify state files
- ❌ Set `validated_by: "orchestrator"` (must be "state-manager")
- ❌ Skip State Manager consultation "just this once"
- ❌ Think "I'll validate it manually"

**ONLY State Manager may update state files. This is NON-NEGOTIABLE.**

See: `rule-library/R517-universal-state-manager-consultation-law.md`
## 🛑🛑🛑 R322 MANDATORY STOP BEFORE STATE TRANSITIONS 🛑🛑🛑

**SUPREME LAW - VIOLATION = -100% IMMEDIATE FAILURE**

### YOU MUST STOP AFTER:
1. ✅ Completing all TODOs for this state
2. ✅ Updating orchestrator-state-v3.json with new state
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


## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED PR_PLAN_CREATION STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
mkdir -p markers/state-verification && touch "markers/state-verification/state_rules_read_orchestrator_PR_PLAN_CREATION-$(date +%Y%m%d-%H%M%S)"
echo "$(date +%s) - Rules read and acknowledged for PR_PLAN_CREATION" > "markers/state-verification/state_rules_read_orchestrator_PR_PLAN_CREATION-$(date +%Y%m%d-%H%M%S)"
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY PR PLANNING UNTIL RULES ARE READ:
- ❌ Start creating PR plan
- ❌ Generate PR messages
- ❌ Update state files
- ❌ Continue to next state
- ❌ Think about what to do in this state

### ✅ YOU MUST IMMEDIATELY:

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

## 📋 PRIMARY DIRECTIVES

### 🚨🚨🚨 RULE R006 - Orchestrator NEVER Writes Code [BLOCKING]
**Source:** `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`
**Criticality:** BLOCKING - Any code operation = -100% IMMEDIATE FAILURE

READ THIS RULE FILE: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`

**⚠️ R006 WARNING FOR PR_PLAN_CREATION STATE:**
- DO NOT edit any code in PR descriptions!
- DO NOT modify source files for PR preparation!
- DO NOT apply any last-minute fixes!
- You only create the PR PLAN document - NEVER touch code

### 🚨🚨🚨 RULE R370 - PR Plan Creation Requirements [BLOCKING]
**Source:** `$CLAUDE_PROJECT_DIR/rule-library/R370-pr-plan-creation-requirements.md`
**Criticality:** BLOCKING - Humans need clear PR instructions

**This state implements R370 to create the MASTER-PR-PLAN.md for humans to execute**

Key requirements per R370:
- Document all effort branches in exact merge order
- Provide complete PR templates (title + body)
- Include copy-paste ready commands
- Document rollback procedures
- Create progress tracking checklist

Software Factory creates the plan, humans execute the PRs. We NEVER push to main.

PRs must be ordered by dependencies to avoid conflicts and ensure clean merges.

## 🎯 STATE OBJECTIVES

In the PR_PLAN_CREATION state, you are responsible for:

1. **Creating MASTER-PR-PLAN.md**
   - List all branches to be PRed
   - Specify exact merge order
   - Provide PR title and body for each
   - Include verification steps

2. **Determining Dependency Order**
   - Core/infrastructure first
   - Features that depend on core next
   - UI/frontend typically last
   - Independent features can be parallel

3. **Generating PR Templates**
   - Title following project conventions
   - Comprehensive description
   - Testing notes
   - Breaking changes if any

4. **Creating Merge Instructions**
   - Step-by-step process for humans
   - Verification between PRs
   - Rollback instructions if needed

## 📝 REQUIRED ACTIONS

### Step 1: Gather All Effort Information
```bash
# List all effort branches
echo "📋 Gathering effort branches..."
EFFORT_BRANCHES=""

for effort_dir in /efforts/*/; do
    if [ -d "$effort_dir/.git" ]; then
        effort_name=$(basename "$effort_dir")
        cd "$effort_dir"
        branch=$(git branch --show-current)
        
        if [ "$branch" != "main" ] && [ "$branch" != "integration-testing" ]; then
            echo "Found effort: $effort_name on branch: $branch"
            EFFORT_BRANCHES="$EFFORT_BRANCHES$branch "
            
            # Get commit count
            commit_count=$(git rev-list --count origin/main..HEAD 2>/dev/null || echo "unknown")
            echo "  Commits ahead of main: $commit_count"
            
            # Get changed files count
            changed_files=$(git diff --name-only origin/main..HEAD 2>/dev/null | wc -l)
            echo "  Files changed: $changed_files"
        fi
    fi
done

echo "Total effort branches: $(echo $EFFORT_BRANCHES | wc -w)"
```

### Step 2: Analyze Dependencies
```bash
# Read implementation plan to understand dependencies
if [ -f ".software-factory/PROJECT-IMPLEMENTATION-PLAN--*.md" ]; then
    echo "📋 Reading implementation plan for dependencies..."
    grep -i "depend\|require\|before\|after" .software-factory/PROJECT-IMPLEMENTATION-PLAN--*.md || true
fi

# Check for explicit dependency documentation
for dep_file in DEPENDENCIES.md EFFORT-DEPENDENCIES.md effort-dependencies.txt; do
    if [ -f "$dep_file" ]; then
        echo "📋 Found dependency file: $dep_file"
        cat "$dep_file"
    fi
done
```

### Step 3: Create MASTER-PR-PLAN.md
```bash
cat > MASTER-PR-PLAN.md << 'EOF'
# MASTER PR PLAN
Generated: $(date)
By: Software Factory 2.0

## 🚨 CRITICAL INSTRUCTIONS FOR HUMANS 🚨

**SOFTWARE FACTORY HAS COMPLETED ALL IMPLEMENTATION AND TESTING**

This document contains the exact steps to create Pull Requests for merging all completed work into the main branch.

### Prerequisites
1. Ensure you have write access to the repository
2. Verify all branches listed below exist and are pushed
3. Have GitHub CLI (`gh`) installed and authenticated OR use GitHub web interface

## 📋 VERIFICATION CHECKLIST

Before starting PR creation:
- [ ] All effort branches are pushed to remote
- [ ] Integration testing branch shows successful build
- [ ] No uncommitted changes in any effort directory
- [ ] You understand the merge order below is MANDATORY

## 🔄 PR MERGE ORDER (CRITICAL - DO NOT DEVIATE)

**IMPORTANT**: PRs must be created and merged in this EXACT order to avoid conflicts.

### Phase 1: Core Infrastructure
These PRs establish the foundation and MUST be merged first.

#### PR #1: [Core Module Name]
**Branch**: `effort-core-module`
**Depends On**: None (merge first)

**PR Title**:
```
feat(core): implement core module infrastructure
```

**PR Body**:
```markdown
## Summary
Implements the core module that provides foundational functionality for the system.

## Changes
- Added core configuration system
- Implemented base interfaces
- Set up logging infrastructure

## Testing
- Unit tests: ✅ All passing
- Integration tests: ✅ Verified in integration-testing branch
- Build status: ✅ Successful

## Breaking Changes
None

## Related Issues
Closes #[issue-number]

## Verification
After merging, verify:
1. Build passes on main
2. No conflicts with subsequent PRs
```

**Merge Instructions**:
1. Create PR from `effort-core-module` to `main`
2. Wait for CI/CD checks to pass
3. Request code review if required
4. Merge using "Squash and merge" or project standard
5. Delete branch after merge

---

#### PR #2: [Dependent Module]
**Branch**: `effort-dependent-module`
**Depends On**: PR #1 (MUST merge after core)

**PR Title**:
```
feat(module): add dependent module functionality
```

**PR Body**:
```markdown
## Summary
Adds dependent module that builds on core infrastructure.

## Changes
- Implemented feature X using core module
- Added API endpoints
- Integrated with existing services

## Testing
- Unit tests: ✅ All passing
- Integration tests: ✅ Verified in integration-testing branch
- Build status: ✅ Successful

## Dependencies
- Requires PR #1 (core module) to be merged first

## Breaking Changes
None

## Verification
After merging, verify:
1. Build still passes
2. Integration with core module works
```

**Merge Instructions**:
1. WAIT for PR #1 to be fully merged
2. Rebase if necessary: `git rebase origin/main`
3. Create PR from `effort-dependent-module` to `main`
4. Follow standard merge process

---

### Phase 2: Feature Implementation
These can be merged in parallel after Phase 1 is complete.

#### PR #3: [Feature A]
**Branch**: `effort-feature-a`
**Depends On**: Phase 1 completion

[PR Template similar to above]

#### PR #4: [Feature B]
**Branch**: `effort-feature-b`  
**Depends On**: Phase 1 completion

[PR Template similar to above]

---

### Phase 3: UI/Frontend (if applicable)
Merge after all backend features are complete.

#### PR #5: [UI Components]
**Branch**: `effort-ui-components`
**Depends On**: All backend PRs

[PR Template similar to above]

---

## 📊 PR CREATION COMMANDS

### Using GitHub CLI (Recommended)
```bash
# PR #1: Core Module
gh pr create \
  --base main \
  --head effort-core-module \
  --title "feat(core): implement core module infrastructure" \
  --body-file PR-BODY-1.md

# PR #2: Dependent Module (after PR #1 merged)
gh pr create \
  --base main \
  --head effort-dependent-module \
  --title "feat(module): add dependent module functionality" \
  --body-file PR-BODY-2.md

# Continue for all PRs...
```

### Using GitHub Web Interface
1. Navigate to repository on GitHub
2. Click "Pull requests" → "New pull request"
3. Select base: `main`, compare: `effort-branch-name`
4. Copy title and body from templates above
5. Create pull request

## ⚠️ CRITICAL WARNINGS

### DO NOT:
- Skip the merge order - will cause conflicts
- Force merge with failing tests
- Merge multiple PRs simultaneously in wrong order
- Delete branches before confirming merge success

### DO:
- Follow the exact order specified
- Wait for each PR to fully merge before starting next
- Run verification steps after each merge
- Keep this document for reference

## 🔄 ROLLBACK PLAN

If issues arise after merging:

### For Single PR Issues:
1. Revert the problematic PR:
   ```bash
   gh pr revert [PR-NUMBER]
   ```
2. Fix the issue in a new branch
3. Create a new PR with fixes

### For Multiple PR Issues:
1. Revert in reverse order (last merged first)
2. Identify root cause
3. Fix and re-submit PRs in correct order

## 📈 PROGRESS TRACKING

Track your progress here:

### Merge Status
- [ ] PR #1: Core Module - Status: [CREATED/APPROVED/MERGED]
- [ ] PR #2: Dependent Module - Status: [WAITING/CREATED/APPROVED/MERGED]
- [ ] PR #3: Feature A - Status: [WAITING/CREATED/APPROVED/MERGED]
- [ ] PR #4: Feature B - Status: [WAITING/CREATED/APPROVED/MERGED]
- [ ] PR #5: UI Components - Status: [WAITING/CREATED/APPROVED/MERGED]

### Verification Checklist
- [ ] All PRs created
- [ ] All CI/CD checks passing
- [ ] All PRs reviewed (if required)
- [ ] All PRs merged
- [ ] Branches deleted
- [ ] Production deployment successful (if applicable)

## 📝 NOTES FOR REVIEWERS

When reviewing these PRs, please note:
1. All code has been integration tested together
2. The integration-testing branch contains the full working system
3. Each PR has been validated to work with others in sequence
4. Comprehensive testing has been completed

## 🏁 COMPLETION

Once all PRs are merged:
1. Verify main branch builds successfully
2. Run full test suite on main
3. Tag release if appropriate
4. Deploy to production if applicable
5. Archive this PR plan for future reference

---

## Appendix: Individual PR Bodies

### PR-BODY-1.md
[Full PR body for PR #1]

### PR-BODY-2.md
[Full PR body for PR #2]

[Continue for all PRs...]

---

END OF MASTER PR PLAN
EOF

echo "✅ MASTER-PR-PLAN.md created successfully"
```

### Step 4: Create Individual PR Body Files
```bash
# Create separate files for each PR body for easy copying
PR_NUMBER=1

for branch in $EFFORT_BRANCHES; do
    cat > "PR-BODY-${PR_NUMBER}.md" << EOF
## Summary
Implementation of [describe what this effort does]

## Changes
- [List key changes]
- [Include important modifications]
- [Note new features added]

## Files Modified
$(cd /efforts/effort-name && git diff --name-only origin/main..HEAD | head -20)

## Testing
- Unit Tests: ✅ All passing
- Integration Tests: ✅ Verified in integration-testing branch
- Build Status: ✅ Successful
- Code Coverage: [X%]

## Dependencies
[List any PR dependencies]

## Breaking Changes
[None OR describe breaking changes]

## Documentation
[Updated OR needs update]

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Comments added for complex code
- [ ] Documentation updated
- [ ] Tests added/updated
- [ ] All tests passing

## Screenshots (if applicable)
[Add any relevant screenshots]

## Additional Notes
[Any additional context for reviewers]
EOF
    
    echo "Created PR-BODY-${PR_NUMBER}.md"
    PR_NUMBER=$((PR_NUMBER + 1))
done
```

### Step 5: Create Verification Script
```bash
cat > verify-prs.sh << 'BASH'
#!/bin/bash
# PR Verification Script

echo "🔍 PR Readiness Verification"
echo "============================"

# Check all branches exist
echo "Checking branches..."
for branch in effort-1 effort-2 effort-3; do
    if git ls-remote --heads origin "$branch" > /dev/null 2>&1; then
        echo "✅ Branch exists: $branch"
    else
        echo "❌ Branch missing: $branch"
    fi
done

# Check for uncommitted changes
echo ""
echo "Checking for uncommitted changes..."
for effort in /efforts/*/; do
    if [ -d "$effort/.git" ]; then
        cd "$effort"
        if [ -n "$(git status --porcelain)" ]; then
            echo "⚠️ Uncommitted changes in: $(basename $effort)"
        else
            echo "✅ Clean: $(basename $effort)"
        fi
    fi
done

# Verify integration branch
echo ""
echo "Checking integration branch..."
cd /efforts/integration-testing
if git log --oneline -1; then
    echo "✅ Integration branch ready"
else
    echo "❌ Integration branch issues"
fi

echo ""
echo "============================"
echo "Ready for PR creation: [YES/NO based on above]"
BASH

chmod +x verify-prs.sh
echo "✅ Created verification script: verify-prs.sh"
```

## ⚠️ CRITICAL REQUIREMENTS

### 🔴 R006 VIOLATION DETECTION 
If you find yourself about to:
- Edit code to prepare for PRs
- Apply last-minute fixes before PR creation
- Modify source files for better PR descriptions
- Clean up code before documenting PRs
- Make excuses like "just tidying up for the PR"

**STOP IMMEDIATELY - You are violating R006!**
You create documentation ONLY - NO code modifications allowed!

### PR Order Must Reflect Dependencies
- Core infrastructure ALWAYS first
- Dependencies resolved before dependents
- Independent features can be parallel
- UI/frontend typically last

### Each PR Must Be Complete
- Include all commits from effort branch
- Have comprehensive description
- Reference any related issues
- Include test results

### Never Actually Create PRs
- Software Factory creates the PLAN
- Humans execute the plan
- We prepare everything, they click buttons

## 🚫 FORBIDDEN ACTIONS

1. **NEVER edit any code files** - R006 VIOLATION = -100%
2. **NEVER apply fixes before PR creation** - R006 VIOLATION = -100%
3. **NEVER modify source for PR prep** - R006 VIOLATION = -100%
4. **NEVER use `gh pr create` yourself**
5. **NEVER push to main branch**
6. **NEVER merge PRs automatically**
7. **NEVER skip dependency analysis**
8. **NEVER create incomplete PR plans**

## ✅ PROJECT_DONE CRITERIA

Before transitioning to PROJECT_DONE state:
- [ ] MASTER-PR-PLAN.md created
- [ ] All effort branches documented
- [ ] Merge order determined
- [ ] PR templates generated
- [ ] Verification script created
- [ ] All files committed and pushed

## 🔄 STATE TRANSITIONS

### Success Path:
```
PR_PLAN_CREATION → PROJECT_DONE
```
- PR plan complete
- All documentation ready
- Software Factory work complete!

### Error Path:
```
PR_PLAN_CREATION → ERROR_RECOVERY
```
- Missing effort information
- Cannot determine dependencies
- Branches not ready

## 📊 VERIFICATION CHECKLIST

Before leaving this state:
```bash
# Verify master plan created
ls -la MASTER-PR-PLAN.md

# Verify PR body files created
ls -la PR-BODY-*.md

# Run verification script
./verify-prs.sh

# Verify all committed
git status

# Final commit
git add -A
git commit -m "pr-plan: MASTER-PR-PLAN.md created - Software Factory complete!"
git push
```

## 🎓 GRADING CRITERIA

You will be evaluated on:
1. **Plan Completeness** (30%)
   - All branches included
   - All PR details provided
   
2. **Dependency Analysis** (30%)
   - Correct merge order
   - Dependencies identified
   
3. **Documentation Quality** (20%)
   - Clear instructions
   - Complete PR templates
   
4. **Usability** (20%)
   - Easy for humans to follow
   - Includes verification steps

## 💡 TIPS FOR PROJECT_DONE

1. **Think Like a Human**: Make instructions crystal clear
2. **Include Everything**: Commands, templates, warnings
3. **Order Matters**: Wrong order = merge conflicts
4. **Be Explicit**: Leave nothing to interpretation

## 🚨 COMMON PITFALLS TO AVOID

1. **Vague Instructions**: Humans need exact steps
2. **Missing Dependencies**: Causes merge conflicts
3. **Incomplete Templates**: PRs get rejected
4. **No Verification**: Problems discovered too late

## 🎯 THE END GOAL

**After this state:**
- Software Factory's work is COMPLETE
- Humans have everything needed to merge code
- Each PR will merge cleanly in order
- The project will be successfully integrated

Remember: This is the culmination of all Software Factory work!
## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-orchestrating
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**


## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**


### 🚨 CRITICAL DISTINCTION: AGENT STOPS ≠ FACTORY STOPS 🚨

**TWO INDEPENDENT DECISIONS - DO NOT CONFUSE THEM:**

#### 1. Should Agent Stop Work? (R322 Technical Requirement)
- Agent completes current state
- Agent saves TODOs and commits state
- Agent exits with `exit 0` (preserves context)
- User runs /continue-orchestrating to resume
- **This is NORMAL at checkpoints**

#### 2. Should Factory Continue? (R405 Operational Status)
- Even though agent stopped, can automation proceed?
- TRUE = Healthy completion, automation can continue
- FALSE = Catastrophic failure, must halt everything
- **R322 checkpoints = TRUE (99.9% of cases)**

### THE PATTERN AT R322 CHECKPOINTS

```bash
# 1. Complete state work
echo "✅ State work complete"

# 2. Update state file
jq '.state_machine.current_state = "NEXT_STATE"' orchestrator-state-v3.json > tmp && mv tmp orchestrator-state-v3.json

# 3. Save TODOs
save_todos "R322_CHECKPOINT"

# 4. Factory continues (operational status)
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"

# 5. Agent stops (technical requirement)
exit 0
```

**Both happen together! Agent stops AND factory continues!**

### WHEN TO USE EACH FLAG VALUE

**TRUE (99.9%):**
- ✅ R322 checkpoint reached
- ✅ State work completed successfully
- ✅ Ready for /continue-orchestrating
- ✅ Waiting for user to continue (NORMAL)
- ✅ Plan ready for review (agent done, factory proceeds)

**FALSE (0.1%):**
- ❌ CATASTROPHIC unrecoverable error
- ❌ Data corruption spreading
- ❌ Critical security violation
- ❌ NOT for R322 checkpoints
- ❌ NOT for user review needs
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

