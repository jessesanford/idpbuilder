# ORCHESTRATOR STATE: VERIFY_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE


## 🚨🚨🚨 STATE PURPOSE [BLOCKING]
Validate integration infrastructure configuration against authoritative sources AFTER creation and BEFORE spawning integration agents. This is a MANDATORY gate that prevents catastrophic failures from wrong repository targeting.

## SF 3.0 Infrastructure Validation Context

This validation state operates within SF 3.0 architecture:
- Reads integration configuration from `integration-containers.json` to verify container setup
- Validates all entries in orchestrator-state-v3.json match target repository configuration
- Updates `state_machine.current_state` in orchestrator-state-v3.json with validation results
- Records validation status atomically per R288 before any integration agents are spawned
- Prevents catastrophic failures where integration infrastructure points to wrong repository (SF planning repo vs target repo)

## 🔴🔴🔴 SUPREME LAW ENFORCEMENT
This state enforces:
- **R507**: Mandatory Infrastructure Validation [BLOCKING]
- **R508**: Target Repository Enforcement [SUPREME LAW]
- **R308**: Incremental Branching Strategy [SUPREME LAW]

ANY violation requires immediate transition to ERROR_RECOVERY with CONTINUE-SOFTWARE-FACTORY=FALSE.

## 🚨🚨🚨 CRITICAL: WRONG REPOSITORY = CATASTROPHIC FAILURE 🚨🚨🚨

**THIS STATE EXISTS TO PREVENT THE #1 CATASTROPHIC FAILURE:**
- Integration infrastructure on PLANNING repository instead of TARGET repository
- This is a SUPREME LAW violation (R508) = -100% IMMEDIATE FAILURE
- Exit code 911 = CATASTROPHIC repository mismatch

## 🚨🚨🚨 ENTRY CONDITIONS [BLOCKING]
MUST have:
1. ✅ Just completed SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE
2. ✅ Integration infrastructure directories created
3. ✅ Integration branch created and pushed
4. ✅ Need to verify EVERYTHING before spawning agents

## 🚨🚨🚨 REQUIRED ACTIONS [MANDATORY]

### 1. LOAD INTEGRATE_WAVE_EFFORTS METADATA
```bash
echo "🔍 Loading integration metadata from orchestrator-state-v3.json..."

INTEGRATE_WAVE_EFFORTS_TYPE=$(echo "✅ State file updated to: $NEXT_STATE"
```

---

### ✅ Step 4: Validate State File (R324)
```bash
# Validate state file before committing
"$CLAUDE_PROJECT_DIR/tools/validate-state.sh" orchestrator-state-v3.json || {
    echo "❌ State file validation failed!"
    exit 288
}
echo "✅ State file validated"
```

---

### ✅ Step 5: Commit State File (R288)
```bash
# Commit and push state file immediately
git add orchestrator-state-v3.json
git commit -m "state: VERIFY_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE → $NEXT_STATE - VERIFY_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE complete [R288]"
git push
echo "✅ State file committed and pushed"
```

---

### ✅ Step 6: Save TODOs (R287 - SUPREME LAW)
```bash
# Save TODO state before transition (R287 trigger)
save_todos "VERIFY_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo
git commit -m "todo: orchestrator - VERIFY_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE complete [R287]"
git push
echo "✅ TODOs saved and committed"
```

---

### ✅ Step 7: Output Continuation Flag (R405 - SUPREME LAW) ⚠️ MANDATORY
```bash
# Output continuation flag as LAST action (R405)
# Use TRUE for normal completion, FALSE only for catastrophic errors

echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
```

**⚠️ THIS MUST BE THE ABSOLUTE LAST LINE OF OUTPUT BEFORE EXIT! ⚠️**

---

### ✅ Step 8: Stop Processing (R322 - SUPREME LAW)
```bash
# Stop for context preservation (R322)
echo "🛑 Stopping for context preservation - use /continue-orchestrating to resume"
exit 0
```

---

## 🚨 CHECKLIST ENFORCEMENT 🚨

**Skipping ANY step in this checklist = FAILURE:**
- Missing Step 2: No next state = stuck forever
- Missing Step 3: No state update = state machine broken (R288 violation, -100%)
- Missing Step 4: Invalid state = corruption (R324 violation)
- Missing Step 5: No commit = state lost on compaction (R288 violation, -100%)
- Missing Step 6: No TODO save = work lost (R287 violation, -20% to -100%)
- Missing Step 7: No continuation flag = automation stops (R405 violation, -100%)
- Missing Step 8: No exit = R322 violation (-100%)

**ALL 8 STEPS ARE MANDATORY - NO EXCEPTIONS**

## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

**MUST output as LAST line before state completion:**
- `CONTINUE-SOFTWARE-FACTORY=TRUE` - State completed, automation can continue
- `CONTINUE-SOFTWARE-FACTORY=FALSE` - Error/block, manual intervention required

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-CONTINUATION-FLAG-MASTER-GUIDE.md**
