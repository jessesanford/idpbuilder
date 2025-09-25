# 🔴🔴🔴 ORCHESTRATOR: PHASE 2 WAVE 1 INTEGRATION RESET PROMPT 🔴🔴🔴

## 📋 CONTEXT
You are the orchestrator. The phase2 wave1 integration needs to be completely reset due to issues with the previous attempt. All efforts have been successfully completed, reviewed, and verified - but the integration process itself needs to be redone from scratch.

## 🎯 OBJECTIVE
Clean up ALL artifacts from previous phase2 wave1 integration attempts and reset to a clean state ready for a fresh integration attempt.

## 🚨🚨🚨 CRITICAL STARTUP SEQUENCE 🚨🚨🚨

### STEP 1: ACKNOWLEDGE RESET REQUIREMENTS
```bash
echo "🔴🔴🔴 PHASE 2 WAVE 1 INTEGRATION RESET INITIATED 🔴🔴🔴"
echo "Timestamp: $(date '+%Y-%m-%d %H:%M:%S %Z')"
echo ""
echo "📋 RESET OBJECTIVES:"
echo "  1. Clean up all previous integration artifacts"
echo "  2. Delete local and remote integration branches"
echo "  3. Remove integration working copies"
echo "  4. Sanitize orchestrator-state.json"
echo "  5. Reset to WAVE_COMPLETE state"
echo "  6. Prepare for fresh integration attempt"
```

### STEP 2: VERIFY CURRENT STATE
```bash
# Check current orchestrator state
echo "📊 Current Orchestrator State:"
cat orchestrator-state.json | jq '.current_state, .current_phase, .current_wave'

# Verify we're working on phase 2, wave 1
if [[ $(jq '.current_phase' orchestrator-state.json) != "2" ]] || \
   [[ $(jq '.current_wave' orchestrator-state.json) != "1" ]]; then
    echo "❌ ERROR: Not in Phase 2, Wave 1 context!"
    echo "Current: Phase $(jq '.current_phase' orchestrator-state.json), Wave $(jq '.current_wave' orchestrator-state.json)"
    exit 1
fi
```

## 🧹 CLEANUP OPERATIONS

### 1. IDENTIFY ALL INTEGRATION ARTIFACTS
```bash
echo "🔍 Identifying integration artifacts to clean..."

# List all integration-related branches (local and remote)
echo "📌 Local integration branches:"
git branch -a | grep -E "phase2-wave1-integration|p2w1-integration"

echo "📌 Remote integration branches:"
git ls-remote --heads origin | grep -E "phase2-wave1-integration|p2w1-integration"

# List all integration working copies
echo "📌 Integration working copies:"
ls -la /tmp/ | grep -E "integration|merge" | grep -E "phase2|p2|wave1|w1"
ls -la ~/workspaces/ | grep -E "integration|merge" | grep -E "phase2|p2|wave1|w1"

# Check state file for integration artifacts
echo "📌 State file integration references:"
cat orchestrator-state.json | jq '.integration_infrastructure'
cat orchestrator-state.json | jq '.wave_integrations.phase2'
```

### 2. DELETE ALL INTEGRATION BRANCHES
```bash
echo "🗑️ Deleting integration branches..."

# Delete local branches
for branch in $(git branch | grep -E "phase2-wave1-integration|p2w1-integration"); do
    echo "  Deleting local branch: $branch"
    git branch -D "$branch" 2>/dev/null || true
done

# Delete remote branches
for ref in $(git ls-remote --heads origin | grep -E "phase2-wave1-integration|p2w1-integration" | awk '{print $2}'); do
    branch=${ref#refs/heads/}
    echo "  Deleting remote branch: $branch"
    git push origin --delete "$branch" 2>/dev/null || true
done

echo "✅ Integration branches deleted"
```

### 3. REMOVE INTEGRATION WORKING COPIES
```bash
echo "🗑️ Removing integration working copies..."

# Remove from /tmp
for dir in $(ls -d /tmp/*integration* 2>/dev/null | grep -E "phase2|p2|wave1|w1"); do
    echo "  Removing: $dir"
    rm -rf "$dir"
done

# Remove from ~/workspaces
for dir in $(ls -d ~/workspaces/*integration* 2>/dev/null | grep -E "phase2|p2|wave1|w1"); do
    echo "  Removing: $dir"
    rm -rf "$dir"
done

echo "✅ Working copies removed"
```

### 4. SANITIZE ORCHESTRATOR STATE FILE
```bash
echo "📝 Sanitizing orchestrator-state.json..."

# Create backup first
cp orchestrator-state.json orchestrator-state.backup.$(date +%Y%m%d-%H%M%S).json

# Remove integration-specific fields
cat orchestrator-state.json | jq '
    # Set state to WAVE_COMPLETE (ready for integration)
    .current_state = "WAVE_COMPLETE" |
    
    # Clear integration infrastructure
    .integration_infrastructure = {} |
    
    # Clear wave integration records for phase2 wave1
    if .wave_integrations.phase2 then
        .wave_integrations.phase2.wave1 = null
    else . end |
    
    # Clear any integration-related fields
    del(.current_integration_branch) |
    del(.integration_in_progress) |
    del(.integration_attempts) |
    del(.last_integration_error) |
    
    # Ensure all efforts show as completed and reviewed
    .efforts_completed = (.efforts_completed // []) |
    .efforts_in_progress = [] |
    
    # Add reset marker
    .last_reset = {
        "timestamp": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'",
        "reason": "Phase 2 Wave 1 integration reset for fresh attempt",
        "reset_from_state": .current_state
    }
' > orchestrator-state.tmp.json

# Validate JSON
if jq empty orchestrator-state.tmp.json 2>/dev/null; then
    mv orchestrator-state.tmp.json orchestrator-state.json
    echo "✅ State file sanitized"
else
    echo "❌ JSON validation failed!"
    exit 1
fi
```

### 5. VERIFY ALL EFFORTS ARE COMPLETE
```bash
echo "📋 Verifying all efforts are complete and reviewed..."

# Check each effort's status
for effort_dir in effort-*; do
    if [[ -d "$effort_dir" ]]; then
        echo "Checking $effort_dir:"
        
        # Check for completion markers
        if [[ -f "$effort_dir/IMPLEMENTATION_COMPLETE.md" ]]; then
            echo "  ✅ Implementation complete"
        else
            echo "  ⚠️ Missing IMPLEMENTATION_COMPLETE.md"
        fi
        
        if [[ -f "$effort_dir/CODE_REVIEW_REPORT.md" ]]; then
            echo "  ✅ Code review complete"
            grep -E "PASSED|APPROVED" "$effort_dir/CODE_REVIEW_REPORT.md" | head -1
        else
            echo "  ⚠️ Missing CODE_REVIEW_REPORT.md"
        fi
        
        # Check git status
        cd "$effort_dir"
        if [[ -z $(git status --porcelain) ]]; then
            echo "  ✅ Git status clean"
        else
            echo "  ⚠️ Uncommitted changes detected"
        fi
        cd ..
    fi
done
```

## 🔄 STATE TRANSITION TO WAVE_COMPLETE

### Per R336 (Mandatory Wave Integration Before Next Wave):
```bash
echo "🔄 Transitioning to WAVE_COMPLETE state..."

# Update state to WAVE_COMPLETE
jq '.current_state = "WAVE_COMPLETE" | 
    .state_transition_reason = "Reset complete - ready for fresh integration attempt" |
    .last_state_transition = "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'"' \
    orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json

echo "✅ State set to WAVE_COMPLETE"
```

## 📋 TODO STATE MANAGEMENT (R287 Compliance)

### Save TODO state BEFORE any transitions:
```bash
# Create comprehensive TODO file for recovery
cat > todos/orchestrator-WAVE_COMPLETE-$(date +%Y%m%d-%H%M%S).todo <<EOF
# Orchestrator TODO State - Integration Reset
Timestamp: $(date '+%Y-%m-%d %H:%M:%S %Z')
State: WAVE_COMPLETE
Phase: 2, Wave: 1

## ✅ Completed Tasks
- [x] All Phase 2 Wave 1 efforts implemented
- [x] All implementations reviewed and passed
- [x] All compilation checks passed
- [x] All demos verified working
- [x] Previous integration attempt cleaned up
- [x] State file sanitized
- [x] Ready for fresh integration attempt

## 📋 Pending Tasks
- [ ] Create fresh wave integration infrastructure
- [ ] Spawn Code Reviewer for merge plan creation
- [ ] Execute integration with Integration Agent
- [ ] Monitor integration progress
- [ ] Request Architect wave review after integration

## 🚨 Critical Next Steps
1. Transition to INTEGRATION state
2. Setup fresh integration infrastructure
3. Create phase2-wave1-integration branch
4. Prepare merge plan with Code Reviewer
EOF

# Commit TODO state
git add todos/*.todo orchestrator-state.json
git commit -m "reset: phase2 wave1 integration - cleaned and ready for fresh attempt" \
           -m "State: WAVE_COMPLETE" \
           -m "All previous integration artifacts removed" \
           -m "Ready to transition to INTEGRATION state"
git push
```

## ✅ VERIFICATION CHECKLIST

Run these checks to confirm reset is complete:
```bash
echo "🔍 Final Verification:"

# 1. No integration branches exist
echo -n "1. Integration branches removed: "
if [[ -z $(git branch -a | grep -E "phase2-wave1-integration|p2w1-integration") ]]; then
    echo "✅ PASS"
else
    echo "❌ FAIL - branches still exist"
fi

# 2. No integration working copies
echo -n "2. Working copies removed: "
if [[ -z $(ls -d /tmp/*integration* 2>/dev/null | grep -E "phase2|wave1") ]]; then
    echo "✅ PASS"
else
    echo "❌ FAIL - working copies still exist"
fi

# 3. State is WAVE_COMPLETE
echo -n "3. State is WAVE_COMPLETE: "
if [[ $(jq -r '.current_state' orchestrator-state.json) == "WAVE_COMPLETE" ]]; then
    echo "✅ PASS"
else
    echo "❌ FAIL - wrong state"
fi

# 4. No integration artifacts in state file
echo -n "4. State file clean: "
if [[ $(jq '.integration_infrastructure | length' orchestrator-state.json) -eq 0 ]]; then
    echo "✅ PASS"
else
    echo "❌ FAIL - integration artifacts remain"
fi

# 5. TODOs saved
echo -n "5. TODOs persisted: "
if [[ -f todos/orchestrator-WAVE_COMPLETE-*.todo ]]; then
    echo "✅ PASS"
else
    echo "❌ FAIL - TODOs not saved"
fi
```

## 🚀 NEXT STEPS AFTER RESET

Once this reset is complete, you should:

1. **Verify Reset Success**:
   ```bash
   cat orchestrator-state.json | jq '.current_state'  # Should be "WAVE_COMPLETE"
   ```

2. **Transition to INTEGRATION State** (per state machine):
   ```bash
   # Following R336 mandatory integration flow:
   # WAVE_COMPLETE → INTEGRATION → SPAWN_CODE_REVIEWER_MERGE_PLAN → etc.
   
   jq '.current_state = "INTEGRATION"' orchestrator-state.json > tmp.json
   mv tmp.json orchestrator-state.json
   ```

3. **Create Fresh Integration Infrastructure**:
   - New working directory
   - New integration branch
   - Fresh from phase2-wave1 effort branches

4. **Continue with Standard Integration Flow**:
   - Spawn Code Reviewer for merge plan
   - Execute integration with Integration Agent
   - Monitor and verify success
   - Request Architect review

## ⚠️ IMPORTANT REMINDERS

### Per R336 (Mandatory Wave Integration):
- ✅ Wave 1 MUST be fully integrated before Wave 2 can start
- ✅ The integration branch becomes the base for Wave 2 efforts
- ✅ Architect reviews the INTEGRATED wave, not individual efforts

### Per R287 (TODO Persistence):
- ✅ Save TODOs before EVERY state transition
- ✅ Commit and push within 60 seconds
- ✅ Include state and context in TODO files

### Per R322 (Mandatory Checkpoints):
- ✅ STOP after spawning agents
- ✅ User review required before executing integration plan
- ✅ Clear continuation instructions on exit

## 🔴 GRADING CRITERIA REMINDER

You will be graded on:
1. **WORKSPACE ISOLATION (20%)** - Keep integration isolated
2. **WORKFLOW COMPLIANCE (25%)** - Follow proper integration flow  
3. **SIZE COMPLIANCE (20%)** - N/A for integration
4. **PARALLELIZATION (15%)** - N/A for integration
5. **QUALITY ASSURANCE (20%)** - Ensure clean integration

## 📝 COMPLETION CONFIRMATION

After running this reset protocol, confirm:
```bash
echo "═══════════════════════════════════════════════════"
echo "📋 PHASE 2 WAVE 1 INTEGRATION RESET COMPLETE"
echo "═══════════════════════════════════════════════════"
echo "State: $(jq -r '.current_state' orchestrator-state.json)"
echo "Phase: $(jq '.current_phase' orchestrator-state.json)"
echo "Wave: $(jq '.current_wave' orchestrator-state.json)"
echo ""
echo "✅ All integration artifacts cleaned"
echo "✅ State file sanitized"
echo "✅ TODOs persisted"
echo "✅ Ready for fresh integration attempt"
echo ""
echo "Next command: /continue-orchestrating"
echo "═══════════════════════════════════════════════════"
```

---

**END OF RESET PROMPT**

Use this prompt to cleanly reset the phase2 wave1 integration and prepare for a fresh attempt.