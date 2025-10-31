# RULE R255 - POST-AGENT WORK VERIFICATION AND RECOVERY

## 🔴🔴🔴 SUPREME LAW: VERIFY EVERY AGENT'S WORK LOCATION 🔴🔴🔴

**Rule ID:** R255
**Criticality:** BLOCKING - Wrong location = DELETE AND RESTART
**Applies to:** Orchestrator verifying ALL agent work
**States:** MONITOR, WAITING_FOR_EFFORT_PLANS, post-agent-completion

## THE RULE

After EVERY agent (SW Engineer, Code Reviewer, Architect) reports completion, the orchestrator MUST:

1. **VERIFY** work is in correct `/efforts/phase{X}/wave{Y}/{effort}` directory
2. **VERIFY** work is committed to correct branch
3. **VERIFY** work is pushed to correct remote
4. **DELETE AND RESTART** if ANY verification fails

## MANDATORY VERIFICATION PROTOCOL

```bash
verify_agent_work() {
    local AGENT_TYPE="$1"
    local EFFORT_NAME="$2"
    local PHASE="$3"
    local WAVE="$4"
    
    echo "🔍 R255: POST-AGENT WORK VERIFICATION"
    
    # Expected values
    local EXPECTED_DIR="$(pwd)/efforts/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}"
    local EXPECTED_BRANCH="phase${PHASE}/wave${WAVE}/${EFFORT_NAME}"
    local EXPECTED_REMOTE="origin/${EXPECTED_BRANCH}"
    
    # Check 1: Verify directory exists and has work
    if [ ! -d "$EXPECTED_DIR" ]; then
        echo "❌ R255 VIOLATION: Directory does not exist: $EXPECTED_DIR"
        return 1
    fi
    
    cd "$EXPECTED_DIR"
    
    # Check 2: Verify we're on correct branch
    local CURRENT_BRANCH=$(git branch --show-current)
    if [ "$CURRENT_BRANCH" != "$EXPECTED_BRANCH" ]; then
        echo "❌ R255 VIOLATION: Wrong branch!"
        echo "   Expected: $EXPECTED_BRANCH"
        echo "   Found: $CURRENT_BRANCH"
        return 1
    fi
    
    # Check 3: Verify work exists (check for expected files)
    case "$AGENT_TYPE" in
        "code-reviewer")
            if [ ! -f "IMPLEMENTATION-PLAN.md" ]; then
                echo "❌ R255 VIOLATION: No IMPLEMENTATION-PLAN.md found!"
                return 1
            fi
            ;;
        "sw-engineer")
            if [ -z "$(find pkg/ -name "*.go" 2>/dev/null)" ]; then
                echo "❌ R255 VIOLATION: No code files in pkg/ directory!"
                return 1
            fi
            ;;
        "architect")
            if [ ! -f "ARCHITECTURE-REVIEW.md" ]; then
                echo "❌ R255 VIOLATION: No ARCHITECTURE-REVIEW.md found!"
                return 1
            fi
            ;;
    esac
    
    # Check 4: Verify commits exist
    local COMMITS=$(git log --oneline -n 1 2>/dev/null | wc -l)
    if [ "$COMMITS" -eq 0 ]; then
        echo "❌ R255 VIOLATION: No commits found!"
        return 1
    fi
    
    # Check 5: Verify pushed to remote
    local REMOTE_EXISTS=$(git ls-remote --heads origin "$EXPECTED_BRANCH" 2>/dev/null | wc -l)
    if [ "$REMOTE_EXISTS" -eq 0 ]; then
        echo "❌ R255 VIOLATION: Branch not pushed to remote!"
        return 1
    fi
    
    # Check 6: Verify local is up to date with remote
    git fetch origin "$EXPECTED_BRANCH" >/dev/null 2>&1
    local BEHIND=$(git rev-list HEAD..origin/"$EXPECTED_BRANCH" --count 2>/dev/null)
    local AHEAD=$(git rev-list origin/"$EXPECTED_BRANCH"..HEAD --count 2>/dev/null)
    
    if [ "$AHEAD" -gt 0 ]; then
        echo "❌ R255 VIOLATION: Local changes not pushed!"
        echo "   $AHEAD commits ahead of remote"
        return 1
    fi
    
    echo "✅ R255: All verifications passed!"
    echo "   Directory: $EXPECTED_DIR ✓"
    echo "   Branch: $EXPECTED_BRANCH ✓"
    echo "   Remote: $EXPECTED_REMOTE ✓"
    echo "   Work exists: ✓"
    echo "   Committed: ✓"
    echo "   Pushed: ✓"
    
    return 0
}
```

## INTELLIGENT RECOVERY ASSESSMENT

Before deleting work, the orchestrator MUST assess if work is salvageable:

```bash
assess_salvageability() {
    local AGENT_TYPE="$1"
    local EFFORT_NAME="$2"
    local PHASE="$3"
    local WAVE="$4"
    
    echo "🔍 R255: Assessing if work can be salvaged..."
    
    # Expected values
    local EXPECTED_DIR="efforts/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}"
    local EXPECTED_BRANCH="phase${PHASE}/wave${WAVE}/${EFFORT_NAME}"
    local BASE_BRANCH="main"  # or from target-repo-config.yaml
    
    # SCENARIO 1: Work in wrong directory but correct content
    # SALVAGEABLE - Move files to correct location
    if [ -f "IMPLEMENTATION-PLAN.md" ] && [ "$(pwd)" != "$EXPECTED_DIR" ]; then
        echo "✅ SALVAGEABLE: Work in wrong directory, moving..."
        mkdir -p "$EXPECTED_DIR"
        cp IMPLEMENTATION-PLAN.md "$EXPECTED_DIR/"
        cd "$EXPECTED_DIR"
        git add IMPLEMENTATION-PLAN.md
        git commit -m "fix: salvage misplaced work [R255 recovery]"
        git push
        echo "✅ Work salvaged and moved to correct location"
        return 0
    fi
    
    # SCENARIO 2: Work on wrong branch but correct directory
    # SALVAGEABLE - Cherry-pick commits
    if [ -d "$EXPECTED_DIR" ]; then
        cd "$EXPECTED_DIR"
        CURRENT_BRANCH=$(git branch --show-current)
        
        if [ "$CURRENT_BRANCH" != "$EXPECTED_BRANCH" ]; then
            echo "⚠️ Work on wrong branch: $CURRENT_BRANCH"
            
            # Check if commits can be cherry-picked
            COMMITS_TO_SAVE=$(git log --oneline "$BASE_BRANCH".."$CURRENT_BRANCH" --grep="$EFFORT_NAME" | wc -l)
            
            if [ "$COMMITS_TO_SAVE" -gt 0 ]; then
                echo "✅ SALVAGEABLE: Can cherry-pick $COMMITS_TO_SAVE commits"
                
                # Create correct branch and cherry-pick
                git checkout -b "$EXPECTED_BRANCH" "$BASE_BRANCH"
                git cherry-pick "$BASE_BRANCH".."$CURRENT_BRANCH"
                git push -u origin "$EXPECTED_BRANCH"
                
                # Clean up wrong branch
                git branch -D "$CURRENT_BRANCH"
                echo "✅ Work salvaged via cherry-pick"
                return 0
            fi
        fi
    fi
    
    # SCENARIO 3: Work uncommitted but in correct location
    # SALVAGEABLE - Just commit and push
    if [ -d "$EXPECTED_DIR" ]; then
        cd "$EXPECTED_DIR"
        if [ -n "$(git status --porcelain)" ]; then
            echo "✅ SALVAGEABLE: Uncommitted work in correct location"
            git add -A
            git commit -m "fix: commit unsaved work [R255 recovery]"
            git push
            echo "✅ Work committed and pushed"
            return 0
        fi
    fi
    
    # SCENARIO 4: Work in completely wrong repository
    # NOT SALVAGEABLE - Wrong base, can't merge
    if [ -d "$EXPECTED_DIR/.git" ]; then
        cd "$EXPECTED_DIR"
        REMOTE_URL=$(git remote get-url origin 2>/dev/null)
        
        if [[ "$REMOTE_URL" != *"$TARGET_REPO_URL"* ]]; then
            echo "❌ NOT SALVAGEABLE: Wrong repository!"
            echo "   Expected repo: $TARGET_REPO_URL"
            echo "   Found repo: $REMOTE_URL"
            return 1
        fi
    fi
    
    # SCENARIO 5: Work has merge conflicts with correct branch
    # NOT SALVAGEABLE - Would break existing work
    if [ -d "$EXPECTED_DIR" ]; then
        cd "$EXPECTED_DIR"
        git fetch origin "$EXPECTED_BRANCH" 2>/dev/null
        
        if git merge --no-commit --no-ff origin/"$EXPECTED_BRANCH" 2>/dev/null; then
            git merge --abort
            echo "⚠️ MAYBE SALVAGEABLE: No conflicts detected"
        else
            git merge --abort
            echo "❌ NOT SALVAGEABLE: Merge conflicts with existing branch"
            return 1
        fi
    fi
    
    # SCENARIO 6: Work is for wrong effort entirely
    # NOT SALVAGEABLE - Code doesn't match effort
    if [ "$AGENT_TYPE" = "sw-engineer" ] && [ -d "$EXPECTED_DIR/pkg" ]; then
        cd "$EXPECTED_DIR"
        
        # Check if code mentions wrong effort
        WRONG_EFFORT=$(grep -r "package\|func\|type" pkg/ | grep -v "$EFFORT_NAME" | head -5)
        if [ -n "$WRONG_EFFORT" ]; then
            echo "❌ NOT SALVAGEABLE: Code is for wrong effort!"
            echo "$WRONG_EFFORT"
            return 1
        fi
    fi
    
    echo "❌ Could not determine salvageability - defaulting to NOT salvageable"
    return 1
}
```

## SALVAGE EXAMPLES

### ✅ SALVAGEABLE SCENARIOS:

**Example 1: Files in orchestrator root instead of effort directory**
```bash
# Agent created files in / instead of /efforts/phase1/wave1/api-types/
# Found: /IMPLEMENTATION-PLAN.md
# Expected: /efforts/phase1/wave1/api-types/IMPLEMENTATION-PLAN.md

echo "✅ SALVAGING: Moving files to correct location"
mkdir -p efforts/phase1/wave1/api-types
mv IMPLEMENTATION-PLAN.md efforts/phase1/wave1/api-types/
cd efforts/phase1/wave1/api-types
git add IMPLEMENTATION-PLAN.md
git commit -m "fix: relocate misplaced plan [R255 salvage]"
git push origin phase1/wave1/api-types
```

**Example 2: Correct directory but wrong branch name**
```bash
# Work in: /efforts/phase1/wave1/api-types/
# Current branch: feature/api-types
# Expected branch: phase1/wave1/api-types

echo "✅ SALVAGING: Moving commits to correct branch"
git checkout -b phase1/wave1/api-types main
git cherry-pick feature/api-types~3..feature/api-types
git push -u origin phase1/wave1/api-types
git branch -D feature/api-types
```

**Example 3: Work not committed/pushed**
```bash
# Work in correct location but git status shows uncommitted
echo "✅ SALVAGING: Committing and pushing work"
git add -A
git commit -m "feat: implement api-types [R255 recovery]"
git push origin phase1/wave1/api-types
```

### ❌ NOT SALVAGEABLE SCENARIOS:

**Example 1: Wrong repository entirely**
```bash
# Agent cloned personal fork instead of upstream
# Remote: https://github.com/agent/kcp-fork
# Expected: https://github.com/kcp-dev/kcp

echo "❌ NOT SALVAGEABLE: Wrong repository base"
echo "Must delete and start fresh with correct repo"
rm -rf efforts/phase1/wave1/api-types
# Reclone from correct source
```

**Example 2: Work conflicts with existing branch**
```bash
# Agent based work on outdated main branch
# Conflicts with current phase1/wave1/api-types branch

echo "❌ NOT SALVAGEABLE: Merge conflicts detected"
echo "Work is incompatible with existing branch"
rm -rf efforts/phase1/wave1/api-types
# Start fresh to avoid breaking existing work
```

**Example 3: Wrong effort code**
```bash
# Agent implemented webhook code in api-types directory
# grep shows "package webhook" in api-types effort

echo "❌ NOT SALVAGEABLE: Code for wrong effort"
echo "Found webhook code in api-types directory"
rm -rf efforts/phase1/wave1/api-types
# Must restart with correct implementation
```

## RECOVERY PROTOCOL - DELETE AND RESTART

When verification fails, the orchestrator MUST FIRST try to salvage:

```bash
recover_from_wrong_location() {
    local AGENT_TYPE="$1"
    local EFFORT_NAME="$2"
    local PHASE="$3"
    local WAVE="$4"
    local FAILURE_REASON="$5"
    
    echo "🚨🚨🚨 R255 RECOVERY PROTOCOL ACTIVATED 🚨🚨🚨"
    echo "Agent worked in WRONG location or didn't commit/push!"
    echo "Failure: $FAILURE_REASON"
    
    # STEP 1: Try to salvage the work first
    echo "🔧 Attempting to salvage work..."
    if assess_salvageability "$AGENT_TYPE" "$EFFORT_NAME" "$PHASE" "$WAVE"; then
        echo "✅ Work successfully salvaged! No restart needed."
        return 0
    fi
    
    echo "❌ Work cannot be salvaged - proceeding with DELETE AND RESTART"
    
    # Step 1: Document the violation
    cat >> "R255-violations.log" << EOF
$(date): $AGENT_TYPE failed verification for $EFFORT_NAME
Reason: $FAILURE_REASON
Phase: $PHASE, Wave: $WAVE
EOF
    
    # Step 2: Reset the effort state
    local EFFORT_DIR="efforts/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}"
    
    echo "🗑️ Deleting any incorrect work..."
    if [ -d "$EFFORT_DIR" ]; then
        # Save any salvageable work to temp
        cp -r "$EFFORT_DIR" "/tmp/R255-backup-$(date +%s)-${EFFORT_NAME}" 2>/dev/null
        
        # Reset the directory
        cd "$(pwd)"  # Return to orchestrator directory
        rm -rf "$EFFORT_DIR"
    fi
    
    # Step 3: Recreate clean infrastructure
    echo "🔧 Recreating clean effort infrastructure..."
    mkdir -p "$EFFORT_DIR"
    cd "$EFFORT_DIR"
    
    # Clone and setup branch
    git clone --sparse --filter=blob:none "$TARGET_REPO_URL" .
    git sparse-checkout init --cone
    git sparse-checkout set pkg/
    git checkout -b "phase${PHASE}/wave${WAVE}/${EFFORT_NAME}"
    git push -u origin "phase${PHASE}/wave${WAVE}/${EFFORT_NAME}"
    
    cd "$(pwd)"  # Return to orchestrator directory
    
    # Step 4: Update state to restart
    echo "📝 Updating orchestrator state to restart effort..."
    jq ".efforts_in_progress.${EFFORT_NAME}.status = \"RESTARTING_R255\"" orchestrator-state-v3.json
    jq ".efforts_in_progress.${EFFORT_NAME}.violations += [\"R255: Wrong location/no push\"]" orchestrator-state-v3.json
    
    # Step 5: Respawn agent with ENHANCED instructions
    respawn_with_enhanced_instructions "$AGENT_TYPE" "$EFFORT_NAME" "$PHASE" "$WAVE"
}

respawn_with_enhanced_instructions() {
    local AGENT_TYPE="$1"
    local EFFORT_NAME="$2"
    local PHASE="$3"
    local WAVE="$4"
    
    local EFFORT_DIR="$(pwd)/efforts/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}"
    
    echo "🔄 RESPAWNING $AGENT_TYPE with ENHANCED INSTRUCTIONS"
    
    # Create ultra-specific instructions
    local ENHANCED_PROMPT="⚠️⚠️⚠️ CRITICAL R255 RESTART - YOU PREVIOUSLY WORKED IN THE WRONG LOCATION! ⚠️⚠️⚠️

YOU MUST FOLLOW THESE EXACT STEPS OR FACE DELETION:

1. IMMEDIATELY CD TO THIS EXACT DIRECTORY (COPY AND PASTE):
   cd $EFFORT_DIR

2. VERIFY YOU ARE IN THE CORRECT DIRECTORY:
   pwd
   MUST OUTPUT EXACTLY: $EFFORT_DIR
   
3. VERIFY CORRECT GIT BRANCH:
   git branch --show-current
   MUST OUTPUT EXACTLY: phase${PHASE}/wave${WAVE}/${EFFORT_NAME}

4. DO ALL YOUR WORK IN THIS DIRECTORY ONLY!
   - NEVER cd to any other directory
   - NEVER create files outside this directory
   - ALL code goes in: $EFFORT_DIR/pkg/

5. AFTER COMPLETING WORK, YOU MUST:
   git add -A
   git commit -m \"feat: implement $EFFORT_NAME\"
   git push origin phase${PHASE}/wave${WAVE}/${EFFORT_NAME}

6. FINAL VERIFICATION BEFORE REPORTING COMPLETE:
   pwd  # Must show: $EFFORT_DIR
   git status  # Must show: nothing to commit, working tree clean
   git log --oneline -1  # Must show your commit
   git branch -vv  # Must show tracking origin

⚠️ FAILURE TO FOLLOW THESE STEPS = YOUR WORK WILL BE DELETED AND YOU'LL START OVER! ⚠️

Original task: [Include original task here]"

    # Spawn with enhanced prompt
    echo "Task: $AGENT_TYPE"
    echo "Prompt: $ENHANCED_PROMPT"
    echo "Effort: $EFFORT_NAME"
    
    # Record restart attempt
    echo "$(date): Restarted $AGENT_TYPE for $EFFORT_NAME with R255 enhancement" >> "R255-restarts.log"
}
```

## CONTINUOUS MONITORING_SWE_PROGRESS

The orchestrator must check after EVERY agent interaction:

```bash
# In MONITOR state
monitor_agent_completion() {
    local AGENT_ID="$1"
    
    # When agent reports completion
    if [ "$AGENT_STATUS" = "COMPLETE" ]; then
        echo "🔍 Agent $AGENT_ID reports complete - Running R255 verification..."
        
        if ! verify_agent_work "$AGENT_TYPE" "$EFFORT_NAME" "$PHASE" "$WAVE"; then
            echo "❌ R255 VERIFICATION FAILED!"
            recover_from_wrong_location "$AGENT_TYPE" "$EFFORT_NAME" "$PHASE" "$WAVE" "$FAILURE_REASON"
        else
            echo "✅ R255 VERIFICATION PASSED - Work accepted"
        fi
    fi
}
```

## GRADING IMPACT

- First violation: -20% (warning, restart)
- Second violation same effort: -40% (critical, restart)
- Third violation same effort: -100% (AUTOMATIC FAILURE)
- Not checking at all: -100% (AUTOMATIC FAILURE)

## DETECTION HELPERS

```bash
# Quick check if work exists in wrong places
find_misplaced_work() {
    echo "🔍 Searching for misplaced work..."
    
    # Look for work outside efforts directory
    find . -type f -name "*.go" ! -path "./efforts/*" -print 2>/dev/null
    find . -type f -name "IMPLEMENTATION-PLAN.md" ! -path "./efforts/*" -print 2>/dev/null
    
    # Check for uncommitted work
    for effort_dir in efforts/phase*/wave*/*/; do
        if [ -d "$effort_dir/.git" ]; then
            cd "$effort_dir"
            if [ -n "$(git status --porcelain)" ]; then
                echo "⚠️ Uncommitted work in: $effort_dir"
            fi
            cd - >/dev/null
        fi
    done
}

# Verify all efforts have correct remotes
verify_all_remotes() {
    for effort_dir in efforts/phase*/wave*/*/; do
        if [ -d "$effort_dir/.git" ]; then
            cd "$effort_dir"
            local BRANCH=$(git branch --show-current)
            local REMOTE=$(git config --get branch."$BRANCH".remote)
            
            if [ "$REMOTE" != "origin" ]; then
                echo "❌ Wrong/missing remote for: $effort_dir"
            fi
            cd - >/dev/null
        fi
    done
}
```

## MANDATORY ACKNOWLEDGMENT

The orchestrator MUST acknowledge:

"I understand R255: I MUST verify every agent's work is in the correct /efforts/ directory, committed, and pushed. If not, I MUST delete the work and restart the agent with enhanced instructions. Failure to verify = AUTOMATIC FAILURE."

## SALVAGE DECISION MATRIX

| Scenario | Salvageable? | Action | Example |
|----------|-------------|---------|---------|
| Files in wrong directory | ✅ YES | Move to correct location | Found in `/`, move to `/efforts/phase1/wave1/api-types/` |
| Correct directory, wrong branch | ✅ YES | Cherry-pick to correct branch | On `feature/api`, cherry-pick to `phase1/wave1/api-types` |
| Uncommitted work, correct location | ✅ YES | Commit and push | `git add -A && git commit && git push` |
| Wrong branch, no conflicts | ✅ YES | Rebase or cherry-pick | Rebase onto correct branch |
| Stale branch, can fast-forward | ✅ YES | Merge or rebase | `git pull --rebase origin main` |
| Work for correct effort | ✅ YES | Verify and relocate if needed | Code matches effort name |
| Wrong repository base | ❌ NO | Delete and restart | Cloned fork instead of upstream |
| Merge conflicts with target | ❌ NO | Delete to avoid breaking | Would break existing work |
| Wrong effort code | ❌ NO | Delete and restart | Webhook code in api-types directory |
| Corrupted git history | ❌ NO | Delete and restart | Rebase/merge broke history |
| Mixed effort code | ❌ NO | Delete and restart | Multiple efforts in one directory |
| Based on wrong base branch | ❌ NO | Delete and restart | Based on `develop` not `main` |

## SALVAGE FLOWCHART

```
Agent Reports Complete
        ↓
Run R255 Verification
        ↓
    Passed? ──YES──→ Accept Work ✅
        ↓
       NO
        ↓
Assess Salvageability
        ↓
Check Location ──WRONG──→ Can Move? ──YES──→ Move Files → Commit → Push → Done ✅
        ↓                      ↓
    CORRECT                   NO → Delete & Restart 🗑️
        ↓
Check Branch ──WRONG──→ Can Cherry-pick? ──YES──→ Cherry-pick → Push → Done ✅
        ↓                      ↓
    CORRECT                   NO → Delete & Restart 🗑️
        ↓
Check Commits ──MISSING──→ Can Commit? ──YES──→ Commit → Push → Done ✅
        ↓                      ↓
    EXISTS                    NO → Delete & Restart 🗑️
        ↓
Check Push ──NOT PUSHED──→ Can Push? ──YES──→ Push → Done ✅
        ↓                      ↓
     PUSHED                   NO → Delete & Restart 🗑️
        ↓
   All Good ✅
```

## ADDITIONAL SALVAGE SCENARIOS

### Partially Salvageable Work
```bash
# Some files good, others bad
echo "⚠️ PARTIAL SALVAGE: Saving good files only"

# Save good files
cp "$EXPECTED_DIR/IMPLEMENTATION-PLAN.md" /tmp/salvage/
cp "$EXPECTED_DIR/work-log.md" /tmp/salvage/

# Delete everything
rm -rf "$EXPECTED_DIR"

# Recreate and restore good files
mkdir -p "$EXPECTED_DIR"
cd "$EXPECTED_DIR"
git clone --sparse "$TARGET_REPO" .
git checkout -b "$EXPECTED_BRANCH"

# Restore salvaged files
cp /tmp/salvage/* .
git add -A
git commit -m "fix: restore salvaged work [R255 partial recovery]"
git push origin "$EXPECTED_BRANCH"
```

### Cross-Effort Contamination
```bash
# Agent mixed api-types and webhook code
echo "❌ NOT SALVAGEABLE: Cross-effort contamination"
echo "Found:"
echo "  - pkg/apis/... (correct for api-types)"
echo "  - pkg/webhook/... (WRONG - belongs to different effort)"
echo "Cannot untangle - must restart clean"
```

### Salvage with Manual Intervention
```bash
# Complex case needing human decision
echo "⚠️ NEEDS MANUAL REVIEW: Complex salvage scenario"
echo "Options:"
echo "  1. Interactive rebase to fix history"
echo "  2. Manual cherry-pick of specific commits"
echo "  3. Patch extraction and reapplication"

# Save patch for manual review
git format-patch main.."$CURRENT_BRANCH" -o /tmp/R255-patches/

echo "Patches saved to /tmp/R255-patches/ for manual review"
echo "Proceeding with DELETE AND RESTART"
```

## WHY THIS MATTERS

1. **Prevents lost work** - Work in wrong directories won't be found
2. **Ensures integration** - Only properly located work can be merged
3. **Maintains isolation** - Each effort must stay in its directory
4. **Enables tracking** - Can't track progress of misplaced work
5. **Avoids conflicts** - Wrong branches cause merge disasters

## VIOLATIONS TO DETECT

- Work created in orchestrator root directory
- Work created in /tmp or home directories  
- Files created but not committed
- Commits made but not pushed
- Work on wrong git branch
- No sparse checkout (full repo clone)
- Cross-effort contamination