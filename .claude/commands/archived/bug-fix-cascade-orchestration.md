# /bug-fix-cascade-orchestration

╔═══════════════════════════════════════════════════════════════════════════════╗
║                        SOFTWARE FACTORY 2.0                                  ║
║                   BUG FIX CASCADE ORCHESTRATION                              ║
║                                                                               ║
║ Priority: SUPREME - Must complete before ANY other work                      ║
║ Round: 2 (Round 1 fixes already complete and archived)                       ║
╚═══════════════════════════════════════════════════════════════════════════════╝

## 🚨 CRITICAL PRIORITY OVERRIDE 🚨

**This command initiates the HIGHEST PRIORITY bug fix CASCADE.**
**No other work can proceed until these bugs are fixed.**

## 🎯 AGENT IDENTITY ASSIGNMENT

**You are @agent-orchestrator executing BUG FIX CASCADE ROUND 2**

By invoking this command, you MUST:
- Read and acknowledge BUG-FIX-CASCADE-ENFORCEMENT-RULES.md
- Execute fixes according to CRITICAL-BUGS-REQUIRING-FIX-CASCADE.md
- Track everything in orchestrator-state.json under bug_fix_cascade
- Verify with actual test execution
- CASCADE all fixes through every integration

## 🔴 MANDATORY PRE-FLIGHT PROTOCOL

### 1. Read and Acknowledge All Documents
```bash
echo "══════════════════════════════════════════════════"
echo "BUG FIX CASCADE ROUND 2 - INITIATION"
echo "══════════════════════════════════════════════════"
echo ""
echo "Reading enforcement rules..."
cat BUG-FIX-CASCADE-ENFORCEMENT-RULES.md | head -50
echo ""
echo "Reading bug inventory..."
cat CRITICAL-BUGS-REQUIRING-FIX-CASCADE.md | grep "^### BUG-"
echo ""
echo "I ACKNOWLEDGE:"
echo "✓ BUG-FIX-CASCADE-ENFORCEMENT-RULES BF-001 through BF-010"
echo "✓ Bug fix CASCADE takes priority over ALL other work"
echo "✓ No PR until ALL bugs fixed and verified"
echo "✓ Must fix at source, not in integrations"
echo "✓ Must CASCADE through all integrations after fixes"
echo "✓ Definition of complete: Binary builds and all tests pass"
```

### 2. Check Current Bug Status
```bash
# Check orchestrator state for bug tracking
cat orchestrator-state.json | python3 -c "
import json, sys
state = json.loads(sys.stdin.read())
bugs = state.get('bug_fix_cascade', {})
print(f'Bug Fix CASCADE Round: {bugs.get(\"round\", \"NOT FOUND\")}')
print(f'Status: {bugs.get(\"status\", \"NOT FOUND\")}')
print(f'Rules Acknowledged: {bugs.get(\"enforcement_rules_acknowledged\", False)}')
print()
print('Bugs to Fix:')
for bug in bugs.get('bugs_to_fix', []):
    status_icon = '❌' if bug['status'] == 'pending' else '✅'
    print(f'  {status_icon} {bug[\"id\"]}: {bug[\"description\"]} [{bug[\"status\"]}]')
"
```

### 3. Initialize Bug Fix Tracking
```bash
# Update state to acknowledge rules and begin
python3 -c "
import json
from datetime import datetime

with open('orchestrator-state.json', 'r') as f:
    state = json.loads(f.read())

state['bug_fix_cascade']['enforcement_rules_acknowledged'] = True
state['bug_fix_cascade']['status'] = 'in_progress'
state['bug_fix_cascade']['started_at'] = datetime.now().isoformat() + 'Z'
state['current_state'] = 'BUG_FIX_CASCADE'
state['previous_state'] = 'WAVE_COMPLETE'
state['transition_reason'] = 'Initiating Bug Fix CASCADE Round 2 - 5 critical bugs blocking build/test'

with open('orchestrator-state.json', 'w') as f:
    json.dump(state, f, indent=2)

print('✅ Bug Fix CASCADE tracking initialized')
print('State transitioned to: BUG_FIX_CASCADE')
"

git add orchestrator-state.json
git commit -m "state: Initialize Bug Fix CASCADE Round 2

Beginning fix cascade for 5 critical bugs:
- BUG-001: Duplicate ValidationMode
- BUG-002: Duplicate main() in tests
- BUG-003: Mock incompatibility
- BUG-004: Empty test data
- BUG-005: Missing fallback implementation

Priority: SUPREME - No other work until complete"
git push
```

## 📋 BUG FIX EXECUTION PROTOCOL

### For Each Bug (IN ORDER - dependencies matter):

#### Step 1: Spawn SW Engineer to Fix at Source
```bash
# For BUG-001 (Duplicate ValidationMode)
Task: software-engineer
Working Directory: efforts/phase1/wave2/cert-validation
Mission: Fix BUG-001 as documented in /CRITICAL-BUGS-REQUIRING-FIX-CASCADE.md
Instructions:
1. Read the bug description for BUG-001
2. Create pkg/certs/types.go with single ValidationMode definition
3. Remove ValidationMode from validator.go and chain_validator.go
4. Update imports in all affected files
5. Verify: go build ./pkg/certs/... succeeds
6. Commit with: "fix(cert-validation): BUG-001 - resolve duplicate ValidationMode"
7. Push to remote
```

#### Step 2: Spawn Code Reviewer to Verify Fix
```bash
Task: code-reviewer
Working Directory: efforts/phase1/wave2/cert-validation
Mission: Verify BUG-001 is completely fixed
Instructions:
1. Verify ValidationMode exists only in types.go
2. Verify no duplicate declarations remain
3. Run: go build ./pkg/certs/...
4. Confirm build succeeds
5. Check no new issues introduced
6. Return: FIX_COMPLETE or FIX_INCOMPLETE
```

#### Step 3: Update Bug Status
```bash
python3 -c "
import json
with open('orchestrator-state.json', 'r') as f:
    state = json.loads(f.read())

for bug in state['bug_fix_cascade']['bugs_to_fix']:
    if bug['id'] == 'BUG-001':
        bug['status'] = 'fixed'
        bug['fixed_at'] = '$(date -Iseconds)'
        bug['commit'] = '$(git rev-parse HEAD)'
        break

with open('orchestrator-state.json', 'w') as f:
    json.dump(state, f, indent=2)

print('✅ BUG-001 marked as fixed')
"
```

### Repeat for BUG-002 through BUG-005

## 🔄 CASCADE EXECUTION AFTER ALL FIXES

Once ALL bugs are fixed at source:

### Phase 1: Delete All Integration Branches
```bash
for branch in \
  "phase1-wave1-integration" \
  "phase1-wave2-integration" \
  "phase1-integration" \
  "phase2-wave1-integration" \
  "phase2-wave2-integration" \
  "phase2-integration" \
  "project-integration"
do
  echo "Deleting old integration: $branch"
  git push origin --delete "idpbuilder-oci-build-push/$branch" 2>/dev/null || true
done
```

### Phase 2: Recreate Each Integration (IN ORDER)
```bash
# Spawn Integration Agent for each level
Task: integration
Mission: Recreate phase1-wave1-integration with fixed sources
Base: main
Merge: Fixed kind-cert-extraction effort
Verify: Build and tests pass

# Continue for each integration level...
```

## 🎯 SUCCESS VERIFICATION PROTOCOL

### After ALL integrations recreated:
```bash
# Clone final project integration
cd /tmp
rm -rf test-build
git clone https://github.com/jessesanford/idpbuilder.git test-build
cd test-build
git checkout idpbuilder-oci-build-push/project-integration

# Run comprehensive verification
echo "══════════════════════════════════════════════════"
echo "FINAL VERIFICATION"
echo "══════════════════════════════════════════════════"

# 1. Build test
echo "Testing build..."
if go build ./...; then
  echo "✅ BUILD PASSES"
else
  echo "❌ BUILD FAILS - CASCADE INCOMPLETE"
  exit 1
fi

# 2. Test suite
echo "Running tests..."
if go test ./... -count=1; then
  echo "✅ TESTS PASS"
else
  echo "❌ TESTS FAIL - CASCADE INCOMPLETE"
  exit 1
fi

# 3. Race detector
echo "Running race detector..."
if go test -race ./...; then
  echo "✅ NO RACE CONDITIONS"
else
  echo "❌ RACE CONDITIONS - CASCADE INCOMPLETE"
  exit 1
fi

# 4. Binary build
echo "Building binary..."
if go build -o idpbuilder; then
  echo "✅ BINARY BUILDS"
else
  echo "❌ BINARY BUILD FAILS - CASCADE INCOMPLETE"
  exit 1
fi

# 5. Binary runs
echo "Testing binary..."
if ./idpbuilder --help > /dev/null 2>&1; then
  echo "✅ BINARY RUNS"
else
  echo "❌ BINARY FAILS TO RUN - CASCADE INCOMPLETE"
  exit 1
fi

echo ""
echo "🎉 ALL VERIFICATION PASSED - BUG FIX CASCADE COMPLETE! 🎉"
```

## 📊 COMPLETION CRITERIA

The Bug Fix CASCADE is ONLY complete when:
- [ ] All 5 bugs marked "fixed" in orchestrator-state.json
- [ ] All source branches have fixes committed and pushed
- [ ] All 7 integrations recreated with fixed sources
- [ ] Project integration builds: `go build ./...` ✅
- [ ] All tests pass: `go test ./... -count=1` ✅
- [ ] No races: `go test -race ./...` ✅
- [ ] Binary builds: `go build -o idpbuilder` ✅
- [ ] Binary runs: `./idpbuilder --help` ✅
- [ ] BUG-FIX-VERIFICATION-CHECKLIST.md fully checked

## 🚨 ENFORCEMENT REMINDERS

- **BF-001**: This takes priority over EVERYTHING
- **BF-002**: Verify with actual execution, not assumptions
- **BF-003**: No partial fixes - completely fixed or not fixed
- **BF-004**: Fix at source, never in integrations
- **BF-005**: Must CASCADE after fixes
- **BF-009**: Not complete until binary works
- **BF-010**: You are personally responsible for completion

## 🛑 STOPPING POINTS

You MUST stop and report at these points:
1. After acknowledging rules (for confirmation)
2. After each bug fix (for verification)
3. After all source fixes before CASCADE (checkpoint)
4. After CASCADE before final verification
5. After final verification (completion report)

---

**INVOCATION**: When ready to begin, the orchestrator will start with rule acknowledgment and pre-flight checks.

**Remember**: This is THE highest priority. Nothing else matters until these bugs are fixed and verified.