# Integration Agent - TESTING State Rules

## State Definition
The TESTING state validates the integrated code through build, test, and demo execution per R291 requirements.

## 🔴🔴🔴 CRITICAL: R291 DEMO GATE ENFORCEMENT 🔴🔴🔴

**Per R291, you MUST execute and verify demos as part of testing!**
- Build must succeed
- Tests must pass
- **DEMOS MUST RUN** (new requirement per R291)
- All results documented

## Required Actions

### 0. Find Test Harnesses from State (R341 + R340)
```bash
# Find test harnesses from orchestrator state per R340
echo "🔍 LOCATING TEST HARNESSES FROM STATE FILE..."

STATE_FILE="$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json"
if [ ! -f "$STATE_FILE" ]; then
    echo "❌ Cannot find orchestrator-state-v3.json!"
    exit 1
fi

# Get current phase and wave
PHASE_NUM=$(jq -r '.current_phase' "$STATE_FILE")
WAVE_NUM=$(jq -r '.current_wave' "$STATE_FILE")

# Find phase test harness
PHASE_KEY="phase${PHASE_NUM}"
PHASE_HARNESS=$(jq -r ".test_plans.phase[\"$PHASE_KEY\"].test_harness_path" "$STATE_FILE")
PHASE_DEMO=$(jq -r ".test_plans.phase[\"$PHASE_KEY\"].demo_plan_path" "$STATE_FILE")

# Find wave test harness
WAVE_KEY="phase${PHASE_NUM}_wave${WAVE_NUM}"
WAVE_HARNESS=$(jq -r ".test_plans.wave[\"$WAVE_KEY\"].test_harness_path" "$STATE_FILE")
WAVE_DEMO=$(jq -r ".test_plans.wave[\"$WAVE_KEY\"].demo_plan_path" "$STATE_FILE")

echo "📋 TEST HARNESSES LOCATED:"
echo "   Phase Harness: $PHASE_HARNESS"
echo "   Phase Demo: $PHASE_DEMO"
echo "   Wave Harness: $WAVE_HARNESS"
echo "   Wave Demo: $WAVE_DEMO"

# Verify harnesses exist
if [ -f "$WAVE_HARNESS" ]; then
    echo "✅ Wave test harness found"
    PRIMARY_HARNESS="$WAVE_HARNESS"
elif [ -f "$PHASE_HARNESS" ]; then
    echo "✅ Phase test harness found"
    PRIMARY_HARNESS="$PHASE_HARNESS"
else
    echo "⚠️ WARNING: No test harness found in tracked locations"
    echo "Falling back to searching for test scripts..."
fi
```

### 1. Build Attempt
```bash
# Attempt build - DO NOT FIX if it fails
echo "🏗️ R291 BUILD GATE: Starting build verification..."
make build || BUILD_STATUS="FAILED"
echo "Build Status: $BUILD_STATUS"

# Check for build artifacts (R291 ARTIFACT GATE)
if [ -d "dist" ] || [ -d "build" ] || [ -d "out" ] || [ -d "target" ]; then
    echo "✅ Build artifacts found"
    ARTIFACT_STATUS="PASSING"
else
    echo "❌ No build artifacts produced"
    ARTIFACT_STATUS="FAILED"
fi
```

### 2. Test Execution
```bash
# Run tests - DO NOT FIX failures
echo "🧪 R291 TEST GATE: Starting test execution..."

# Use test harness from state if found
if [ -n "$PRIMARY_HARNESS" ] && [ -f "$PRIMARY_HARNESS" ]; then
    echo "📋 Using test harness from state: $PRIMARY_HARNESS"
    bash "$PRIMARY_HARNESS" || TEST_STATUS="FAILED"
else
    # Fallback to standard test commands
    echo "⚠️ Using fallback test execution (no harness from state)"
    make test || TEST_STATUS="FAILED"
    go test ./... -v > test-output.txt 2>&1
fi

echo "Test Status: $TEST_STATUS"
```

### 3. 🎬 Demo Execution (R291/R330 MANDATORY)
```bash
echo "🎬 R291 DEMO GATE: Starting demo execution..."

# Create demo results directory
mkdir -p demo-results

# Run individual effort demos (R330 compliance)
DEMO_STATUS="PASSING"
for effort_dir in effort*/; do
    if [ -f "${effort_dir}demo-features.sh" ]; then
        echo "Running demo for ${effort_dir}..."
        bash "${effort_dir}demo-features.sh" > "demo-results/${effort_dir}-demo.log" 2>&1
        if [ $? -ne 0 ]; then
            echo "❌ Demo failed for ${effort_dir}"
            DEMO_STATUS="FAILED"
        else
            echo "✅ Demo passed for ${effort_dir}"
        fi
    fi
done

# Run integrated wave demo (R291 requirement)
if [ -f "./wave-demo.sh" ]; then
    echo "🌊 Running integrated wave demo..."
    bash ./wave-demo.sh > demo-results/wave-integration-demo.log 2>&1
    if [ $? -ne 0 ]; then
        echo "❌ Wave demo failed"
        DEMO_STATUS="FAILED"
    else
        echo "✅ Wave demo passed"
    fi
elif [ -f "./demo-features.sh" ]; then
    # Fallback to generic demo script
    echo "🎯 Running generic demo script..."
    bash ./demo-features.sh > demo-results/integration-demo.log 2>&1
    if [ $? -ne 0 ]; then
        echo "❌ Demo script failed"
        DEMO_STATUS="FAILED"
    else
        echo "✅ Demo script passed"
    fi
else
    echo "⚠️ WARNING: No demo scripts found - creating basic validation"
    # Create basic demo to verify build works
    cat > basic-demo.sh << 'EOF'
#!/bin/bash
echo "🎬 Basic Integration Demo"
echo "Verifying build outputs exist..."
ls -la dist/ build/ out/ target/ 2>/dev/null
echo "Demo Status: Basic validation only"
EOF
    bash basic-demo.sh > demo-results/basic-demo.log 2>&1
fi

echo "Overall Demo Status: $DEMO_STATUS"
```

### 4. Coverage Analysis (if available)
```bash
go test ./... -cover -coverprofile=coverage.out
go tool cover -func=coverage.out
```

## CRITICAL RULE
- R266 - Upstream Bug Documentation
  - **NEVER FIX** any bugs found
  - **ONLY DOCUMENT** issues
  - Mark all bugs as "NOT FIXED - upstream"

## Documentation Requirements

### R344: Report Metadata Location to State File
```bash
# After creating INTEGRATE_WAVE_EFFORTS-REPORT.md, MUST report location
REPORT_PATH="$(pwd)/.software-factory/INTEGRATE_WAVE_EFFORTS-REPORT.md"
PHASE=$(jq -r '.current_phase' "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json")
WAVE=$(jq -r '.current_wave' "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json")

# Update state file with report location (R344 MANDATORY)
yq -i ".metadata_locations.integration_reports.\"phase${PHASE}_wave${WAVE}\" = {
    \"file_path\": \"$REPORT_PATH\",
    \"created_at\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
    \"created_by\": \"integration\",
    \"integration_branch\": \"$(git branch --show-current)\",
    \"status\": \"$OVERALL_STATUS\"
}" "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json"

# Also report test results location
if [ -f "test-output.txt" ]; then
    yq -i ".metadata_locations.test_results.\"phase${PHASE}_wave${WAVE}\" = {
        \"file_path\": \"$(pwd)/test-output.txt\",
        \"created_at\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
        \"created_by\": \"integration\",
        \"tests_passed\": $TESTS_PASSED,
        \"tests_failed\": $TESTS_FAILED
    }" "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json"
fi

# Report demo results location
if [ -d "demo-results" ]; then
    yq -i ".metadata_locations.demo_results.\"phase${PHASE}_wave${WAVE}\" = {
        \"file_path\": \"$(pwd)/demo-results\",
        \"created_at\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
        \"created_by\": \"integration\",
        \"demo_status\": \"$DEMO_STATUS\"
    }" "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json"
fi

# Commit state update
cd "$CLAUDE_PROJECT_DIR"
git add orchestrator-state-v3.json
git commit -m "state: report integration metadata locations per R344"
git push
```

Document in INTEGRATE_WAVE_EFFORTS-REPORT.md:
- Build results (success/failure)
- Test results (passed/failed count)
- **Demo Status** (PASSING/FAILED) - R291 CRITICAL
- Demo execution logs location
- Failed test details
- Failed demo details
- Upstream bugs found
- Recommendations (but don't implement)

### R291 Gate Status Report Format
```markdown
## R291 Integration Gates Status
- BUILD GATE: [PASSING/FAILED]
- TEST GATE: [PASSING/FAILED]
- DEMO GATE: [PASSING/FAILED]
- ARTIFACT GATE: [PASSING/FAILED]

## Demo Execution Summary
- Individual Effort Demos: X/Y passed
- Wave Integration Demo: [PASSED/FAILED]
- Demo Logs: demo-results/
```

## Transition Rules
- Can transition to: REPORTING
- Always transitions regardless of test results
- Must document all findings

## Success Criteria
- Build attempted and documented
- Tests executed and results captured
- All bugs documented (not fixed)
- Results added to INTEGRATE_WAVE_EFFORTS-REPORT.md

## 🔴🔴🔴 MANDATORY MEASUREMENT RULE - R304 🔴🔴🔴

### R304: Mandatory Line Counter Tool Enforcement
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R304-mandatory-line-counter-enforcement.md`
**Criticality**: BLOCKING - Manual counting = AUTOMATIC -100% FAILURE

**ABSOLUTE REQUIREMENTS:**
- ✅ MUST use `$CLAUDE_PROJECT_DIR/tools/line-counter.sh` for ALL line counting
- ❌ NEVER use `wc -l` or any manual counting method
- ❌ NEVER count lines any other way - this is a -100% automatic failure
- ✅ MUST specify both -b (base branch) and -c (current branch) parameters
- ✅ Base branch MUST be phase integration branch (NOT "main")

**Failure to use the line counter tool = AUTOMATIC -100% GRADE**


## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

### YOU MUST OUTPUT THE CONTINUATION FLAG AS YOUR LAST ACTION

**EVERY STATE COMPLETION MUST END WITH EXACTLY ONE OF:**
```bash
# If state completed successfully and factory should continue:
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"

# If error/block/manual review needed:
echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
```

### CRITICAL REQUIREMENTS:
1. **ABSOLUTE LAST OUTPUT**: This MUST be the very last line of output before state completion
2. **NO TEXT AFTER**: No explanations, summaries, or any text after the flag
3. **EXACTLY AS SHOWN**: Use exact format - no variations like CONTINUE-ORCHESTRATING
4. **ALWAYS REQUIRED**: Every single state must output this flag
5. **GREPPABLE**: Must be on its own line for automation parsing

### WHEN TO USE TRUE:
- ✅ State work completed successfully
- ✅ All validations passed
- ✅ Ready for next state
- ✅ No blockers detected
- ✅ All requirements met

### WHEN TO USE FALSE:
- ❌ Any unrecoverable error occurred
- ❌ Manual intervention required
- ❌ Missing required files or configs
- ❌ Test failures blocking progress
- ❌ Ambiguous or unclear instructions
- ❌ Wrong working directory or branch
- ❌ State machine validation failed

### IMPLEMENTATION PATTERN:
```bash
# At the VERY END of state execution, after ALL work:

# Determine success/failure
if [[ "$STATE_COMPLETED_PROJECT_DONEFULLY" == "true" ]]; then
    echo "✅ State work completed successfully"
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
else
    echo "❌ State encountered issues requiring intervention"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
fi

# DO NOT OUTPUT ANYTHING AFTER THE FLAG!
```

### GRADING IMPACT:
- **Missing flag**: -100% AUTOMATIC FAILURE
- **Text after flag**: -50% penalty
- **Wrong format**: -100% AUTOMATIC FAILURE
- **Multiple flags**: -50% penalty

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

