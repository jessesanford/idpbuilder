# Code Reviewer State: DEMO_VALIDATION

## Purpose
Execute integration demos and validate they work successfully per R291 Gate 4.
This state is the ACTUAL EXECUTION POINT for demo validation.

**🚨 CRITICAL: All demos MUST pass R331 validation (no simulation allowed) before execution 🚨**

## Entry Conditions
- Spawned by orchestrator for demo validation
- Integration demos should exist in demos/ directory
- Integration code review already passed
- Integration build and tests already passed

## State Responsibilities

### 1. Read Task Instructions

Load the demo validation task file created by orchestrator:

```bash
# Read task file
TASK_FILE=".software-factory/demo-validation-task.json"

if [ ! -f "$TASK_FILE" ]; then
    echo "🔴 CRITICAL: Demo validation task file missing!"
    echo "   Expected: $TASK_FILE"
    exit 291
fi

# Extract task parameters
INTEGRATE_WAVE_EFFORTS_TYPE=$(jq -r '.integration_type' "$TASK_FILE")
DEMO_DIR=$(jq -r '.demo_directory' "$TASK_FILE")
REPORT_LOCATION=$(jq -r '.report_location' "$TASK_FILE")
LOG_LOCATION=$(jq -r '.log_location' "$TASK_FILE")

echo "📋 Demo Validation Task Parameters:"
echo "   Integration Type: $INTEGRATE_WAVE_EFFORTS_TYPE"
echo "   Demo Directory: $DEMO_DIR"
echo "   Report Location: $REPORT_LOCATION"
echo "   Log Location: $LOG_LOCATION"
```

### 2. Locate Demos (Using Pre-Planned Paths from R504 + R330)

**PRIORITY 1: Use pre-planned demo paths from orchestrator-state-v3.json**

```bash
# 🔴🔴🔴 R504 + R330: READ PRE-PLANNED DEMO PATHS FIRST 🔴🔴🔴
echo "🔍 Reading pre-planned demo information from orchestrator-state-v3.json..."

CURRENT_PHASE=$(jq -r '.current_phase' "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json" 2>/dev/null || echo "1")
CURRENT_WAVE=$(jq -r '.current_wave' "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json" 2>/dev/null || echo "1")

# Try to get pre-planned demo script file
if [ "$INTEGRATE_WAVE_EFFORTS_TYPE" = "wave" ]; then
    INTEGRATE_WAVE_EFFORTS_KEY="phase${CURRENT_PHASE}_wave${CURRENT_WAVE}"
    PRE_PLANNED_DEMO_SCRIPT=$(jq -r ".pre_planned_infrastructure.integrations.wave_integrations.\"${INTEGRATE_WAVE_EFFORTS_KEY}\".demo_script_file // empty" "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json" 2>/dev/null)
    PRE_PLANNED_DEMO_DESC=$(jq -r ".pre_planned_infrastructure.integrations.wave_integrations.\"${INTEGRATE_WAVE_EFFORTS_KEY}\".demo_description // empty" "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json" 2>/dev/null)
elif [ "$INTEGRATE_WAVE_EFFORTS_TYPE" = "phase" ]; then
    INTEGRATE_WAVE_EFFORTS_KEY="phase${CURRENT_PHASE}"
    PRE_PLANNED_DEMO_SCRIPT=$(jq -r ".pre_planned_infrastructure.integrations.phase_integrations.\"${INTEGRATE_WAVE_EFFORTS_KEY}\".demo_script_file // empty" "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json" 2>/dev/null)
    PRE_PLANNED_DEMO_DESC=$(jq -r ".pre_planned_infrastructure.integrations.phase_integrations.\"${INTEGRATE_WAVE_EFFORTS_KEY}\".demo_description // empty" "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json" 2>/dev/null)
elif [ "$INTEGRATE_WAVE_EFFORTS_TYPE" = "project" ]; then
    PRE_PLANNED_DEMO_SCRIPT=$(jq -r ".pre_planned_infrastructure.integrations.project_integration.demo_script_file // empty" "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json" 2>/dev/null)
    PRE_PLANNED_DEMO_DESC=$(jq -r ".pre_planned_infrastructure.integrations.project_integration.demo_description // empty" "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json" 2>/dev/null)
fi

if [ -n "$PRE_PLANNED_DEMO_SCRIPT" ] && [ "$PRE_PLANNED_DEMO_SCRIPT" != "null" ]; then
    echo "✅ Found pre-planned demo script: $PRE_PLANNED_DEMO_SCRIPT"
    echo "   Description: $PRE_PLANNED_DEMO_DESC"

    # Use pre-planned location
    if [ -f "$PRE_PLANNED_DEMO_SCRIPT" ]; then
        DEMO_SCRIPTS="$PRE_PLANNED_DEMO_SCRIPT"
        DEMO_DIR=$(dirname "$PRE_PLANNED_DEMO_SCRIPT")
        cd "$DEMO_DIR" || exit 291
        echo "✅ Using pre-planned demo location per R504 + R330"
    else
        echo "🔴 CRITICAL: Pre-planned demo script does not exist!"
        echo "   Expected at: $PRE_PLANNED_DEMO_SCRIPT"
        echo "   Integration agent MUST create demo at pre-planned location!"
        echo "   This is a R291 + R504 violation"
        exit 291
    fi
else
    echo "⚠️ No pre-planned demo script found in orchestrator-state-v3.json"
    echo "   Falling back to directory-based discovery..."

    # FALLBACK: Navigate to demos directory
    if [ -d "$DEMO_DIR" ]; then
        cd "$DEMO_DIR" || exit 291
        echo "✅ Found demo directory: $DEMO_DIR"
    else
        echo "⚠️ WARNING: Demo directory not found: $DEMO_DIR"
        echo "   This may indicate integration agent did not create demos"
        echo "   Searching for demos in alternate locations..."

        # Try alternate locations
        if [ -d "demos" ]; then
            cd demos || exit 291
            DEMO_DIR="demos"
            echo "   Found demos in: $DEMO_DIR"
        else
            echo "🔴 CRITICAL: No demo directory found anywhere!"
            echo "   Integration MUST have demos per R291"
            echo "   This is a R291 Gate 4 FAILURE"

            # Create failure report
            mkdir -p "$(dirname "$REPORT_LOCATION")"
            cat > "$REPORT_LOCATION" << EOF
# Demo Validation Report - ${INTEGRATE_WAVE_EFFORTS_TYPE} Integration

## Summary
- Integration Type: ${INTEGRATE_WAVE_EFFORTS_TYPE}
- Demo Directory: ${DEMO_DIR}
- Validation Date: $(date -Iseconds)

## Results
- Demos Passed: 0
- Demos Failed: 0
- Total Demos: 0

## Demo Validation Status: FAILED

## R291 Gate 4 Compliance
🔴 GATE 4: FAILED - No demo directory found

## Failure Reason
No demo directory exists. Integration agent did not create demos as
required by R291. This is a MANDATORY requirement.

## Recommendation
MUST enter ERROR_RECOVERY per R291. Demos are MANDATORY.
EOF
            exit 291
        fi
    fi

    # Find all demo scripts (fallback mode)
    DEMO_SCRIPTS=$(find . -maxdepth 2 -name "*.sh" -type f 2>/dev/null | sort)
fi

if [ -z "$DEMO_SCRIPTS" ]; then
    echo "⚠️ WARNING: No demo scripts found in $DEMO_DIR"
    echo "   This is a R291 Gate 4 FAILURE"
fi

echo "📝 Found demo scripts:"
echo "$DEMO_SCRIPTS" | nl
```

### 3. R331 Validation - Pre-Execution Compliance Check

**MANDATORY: Validate demos meet R331 requirements BEFORE execution**

```bash
echo "🔍 R331: Demo Validation Protocol - Pre-Execution Check"
echo "=========================================================="

R331_VIOLATIONS=0
R331_VIOLATION_DETAILS=""

# R331 Check 1: Pre-demo implementation scan
echo "✓ R331 Check 1: Implementation completeness..."
if find ../src/ ../pkg/ ../internal/ -type f \( -name "*.go" -o -name "*.ts" -o -name "*.py" -o -name "*.rs" \) \
    -exec grep -l "TODO\|FIXME\|XXX\|STUB" {} \; 2>/dev/null | grep -q .; then
    echo "  ❌ FAILED: TODO/FIXME found in implementation"
    echo "  PENALTY: -100% (R331 violation - implementation incomplete)"
    R331_VIOLATIONS=$((R331_VIOLATIONS + 1))
    R331_VIOLATION_DETAILS="${R331_VIOLATION_DETAILS}\n- Implementation has TODO/FIXME/STUB in execution path"
else
    echo "  ✅ PASSED: No TODOs in implementation"
fi

# R331 Check 2: Demo scripts have error handling
echo "✓ R331 Check 2: Demo error handling..."
ERROR_HANDLING_MISSING=false
for demo_script in $DEMO_SCRIPTS; do
    if ! grep -q "set -e\|PIPESTATUS\|exit 1" "$demo_script"; then
        echo "  ❌ FAILED: $demo_script lacks error handling"
        ERROR_HANDLING_MISSING=true
    fi
done

if [ "$ERROR_HANDLING_MISSING" = true ]; then
    echo "  PENALTY: -50% (R331 violation - cannot detect failures)"
    R331_VIOLATIONS=$((R331_VIOLATIONS + 1))
    R331_VIOLATION_DETAILS="${R331_VIOLATION_DETAILS}\n- Demo scripts lack error handling (set -e/PIPESTATUS)"
else
    echo "  ✅ PASSED: All demo scripts have error handling"
fi

# R331 Check 3: No obvious simulation patterns
echo "✓ R331 Check 3: No simulation patterns..."
SIMULATION_DETECTED=false
for demo_script in $DEMO_SCRIPTS; do
    if grep -qi "echo.*✅.*without\|RESULT=.*PASS.*without\|MOCK.*PROJECT_DONE" "$demo_script"; then
        echo "  ❌ FAILED: $demo_script contains simulation patterns"
        SIMULATION_DETECTED=true
    fi
done

if [ "$SIMULATION_DETECTED" = true ]; then
    echo "  PENALTY: -100% (R331 violation - simulation detected)"
    R331_VIOLATIONS=$((R331_VIOLATIONS + 1))
    R331_VIOLATION_DETAILS="${R331_VIOLATION_DETAILS}\n- Simulation patterns detected (hardcoded success)"
else
    echo "  ✅ PASSED: No simulation patterns detected"
fi

# R331 Check 4: External verification present
echo "✓ R331 Check 4: External verification..."
EXTERNAL_VERIFY_MISSING=false
for demo_script in $DEMO_SCRIPTS; do
    if ! grep -q "docker images\|docker pull\|curl.*registry\|sqlite3\|ls -la.*file\|cat.*file" "$demo_script"; then
        echo "  ⚠️ WARNING: $demo_script may lack external verification"
        EXTERNAL_VERIFY_MISSING=true
    fi
done

if [ "$EXTERNAL_VERIFY_MISSING" = true ]; then
    echo "  PENALTY: -75% (R331 violation - no external verification)"
    R331_VIOLATIONS=$((R331_VIOLATIONS + 1))
    R331_VIOLATION_DETAILS="${R331_VIOLATION_DETAILS}\n- No external verification (registry, files, database)"
fi

echo "=========================================================="

# R331 FINAL VERDICT
if [ $R331_VIOLATIONS -gt 0 ]; then
    echo ""
    echo "🔴🔴🔴 R331 VALIDATION: FAILED 🔴🔴🔴"
    echo "Violations: $R331_VIOLATIONS"
    echo -e "Details:$R331_VIOLATION_DETAILS"
    echo ""
    echo "🚨 CRITICAL: Demos do not meet R331 requirements!"
    echo "   Demos are SIMULATED or INCOMPLETE"
    echo "   PENALTY: -100% for simulation, -50% to -75% for other violations"
    echo ""
    echo "REQUIRED ACTION: Fix demos to be R331 compliant before validation"
    echo ""

    # Document R331 failure in report
    cat > "$REPORT_LOCATION" << EOF
# Demo Validation Report - ${INTEGRATE_WAVE_EFFORTS_TYPE} Integration

## Summary
- Integration Type: ${INTEGRATE_WAVE_EFFORTS_TYPE}
- Demo Directory: ${DEMO_DIR}
- Validation Date: $(date -Iseconds)

## R331 Validation Status: FAILED

🔴🔴🔴 CRITICAL: R331 DEMO VALIDATION PROTOCOL VIOLATION 🔴🔴🔴

## R331 Violations Detected
Violation Count: $R331_VIOLATIONS

Details:
${R331_VIOLATION_DETAILS}

## Consequences
- Simulated demos: -100% IMMEDIATE FAILURE
- Missing external verification: -75%
- No error handling: -50%
- Incomplete implementation: -100%

## R291 Gate 4 Compliance
🔴 GATE 4: FAILED - Demos do not meet R331 requirements

## Recommendation
MUST enter ERROR_RECOVERY per R291 + R331.
Demos are not real working implementations - they are simulated or incomplete.

Fix demos to be R331 compliant:
1. Remove all TODO/FIXME from implementation execution path
2. Add error handling (set -e) to all demo scripts
3. Remove simulation patterns (hardcoded success)
4. Add external verification (registry, files, database, etc.)

---

**Generated by Code Reviewer DEMO_VALIDATION state**
**R331 Demo Validation Protocol Enforcement**
**Timestamp: $(date -Iseconds)**
EOF

    # Commit R331 failure
    git add "$REPORT_LOCATION"
    git commit -m "demo: R331 validation FAILED - demos are simulated/incomplete

${INTEGRATE_WAVE_EFFORTS_TYPE} integration demos violate R331 (Demo Validation Protocol).

Violations: $R331_VIOLATIONS
${R331_VIOLATION_DETAILS}

PENALTY: -100% (simulation) or -50% to -75% (other violations)

Report: $REPORT_LOCATION" || true
    git push || true

    echo "📝 R331 failure report created: $REPORT_LOCATION"
    echo "🔴 Exiting with R331 violation code"
    exit 331
else
    echo "✅✅✅ R331 VALIDATION: PASSED ✅✅✅"
    echo "Demos meet all R331 requirements"
    echo "  - Implementation complete (no TODOs)"
    echo "  - Error handling present"
    echo "  - No simulation patterns"
    echo "  - External verification detected"
    echo ""
    echo "Proceeding with demo execution..."
fi
```

### 4. Execute ALL Demo Scripts

Run each demo script and capture results:

```bash
# Initialize counters
DEMOS_PASSED=0
DEMOS_FAILED=0
TOTAL_DEMOS=0

# Create log directory
mkdir -p "$(dirname "$LOG_LOCATION")"
LOG_FILE="$LOG_LOCATION"

# Initialize log
cat > "$LOG_FILE" << EOF
# Demo Validation Execution Log
# Integration Type: ${INTEGRATE_WAVE_EFFORTS_TYPE}
# Validation Started: $(date -Iseconds)
# Demo Directory: ${DEMO_DIR}

EOF

echo "🎬 Executing demo scripts..."
echo "================================"

# Execute each demo script
for demo_script in $DEMO_SCRIPTS; do
    TOTAL_DEMOS=$((TOTAL_DEMOS + 1))

    echo "" | tee -a "$LOG_FILE"
    echo "────────────────────────────────────────" | tee -a "$LOG_FILE"
    echo "Executing: $demo_script" | tee -a "$LOG_FILE"
    echo "Started: $(date -Iseconds)" | tee -a "$LOG_FILE"
    echo "────────────────────────────────────────" | tee -a "$LOG_FILE"

    # Make script executable if it isn't already
    chmod +x "$demo_script"

    # Execute demo script and capture output
    if bash "$demo_script" >> "$LOG_FILE" 2>&1; then
        EXIT_CODE=$?
        echo "✅ $demo_script: PASSED (exit code: $EXIT_CODE)" | tee -a "$LOG_FILE"
        DEMOS_PASSED=$((DEMOS_PASSED + 1))
    else
        EXIT_CODE=$?
        echo "🔴 $demo_script: FAILED (exit code: $EXIT_CODE)" | tee -a "$LOG_FILE"
        DEMOS_FAILED=$((DEMOS_FAILED + 1))

        # Capture error details
        echo "   Error: Demo script returned non-zero exit code" | tee -a "$LOG_FILE"
    fi

    echo "Completed: $(date -Iseconds)" | tee -a "$LOG_FILE"
    echo "" | tee -a "$LOG_FILE"
done

echo "================================"
echo "📊 Demo Execution Summary:"
echo "   Total Demos: $TOTAL_DEMOS"
echo "   Passed: $DEMOS_PASSED"
echo "   Failed: $DEMOS_FAILED"

# Append summary to log
cat >> "$LOG_FILE" << EOF

────────────────────────────────────────
DEMO EXECUTION SUMMARY
────────────────────────────────────────
Total Demos: $TOTAL_DEMOS
Demos Passed: $DEMOS_PASSED
Demos Failed: $DEMOS_FAILED
Validation Completed: $(date -Iseconds)
EOF
```

### 4. Create Demo Evaluation Report

Generate the structured report that orchestrator will read:

```bash
# Determine overall status
if [ $DEMOS_FAILED -eq 0 ] && [ $TOTAL_DEMOS -gt 0 ]; then
    DEMO_STATUS="PASSED"
    GATE_STATUS="✅ GATE 4: PASSED - All demos executed successfully"
    RECOMMENDATION="Integration may proceed to completion"
elif [ $TOTAL_DEMOS -eq 0 ]; then
    DEMO_STATUS="FAILED"
    GATE_STATUS="🔴 GATE 4: FAILED - No demos found"
    RECOMMENDATION="MUST enter ERROR_RECOVERY per R291 - demos are MANDATORY"
else
    DEMO_STATUS="FAILED"
    GATE_STATUS="🔴 GATE 4: FAILED - Demo execution failed"
    RECOMMENDATION="MUST enter ERROR_RECOVERY per R291 - demos are MANDATORY"
fi

# Create report directory
mkdir -p "$(dirname "$REPORT_LOCATION")"

# Generate report
cat > "$REPORT_LOCATION" << EOF
# Demo Validation Report - ${INTEGRATE_WAVE_EFFORTS_TYPE} Integration

## Summary
- Integration Type: ${INTEGRATE_WAVE_EFFORTS_TYPE}
- Demo Directory: ${DEMO_DIR}
- Validation Date: $(date -Iseconds)

## Results
- Demos Passed: ${DEMOS_PASSED}
- Demos Failed: ${DEMOS_FAILED}
- Total Demos: ${TOTAL_DEMOS}

Demo Validation Status: ${DEMO_STATUS}

## R291 Gate 4 Compliance
${GATE_STATUS}

## Individual Demo Results

EOF

# Add individual demo results
for demo_script in $DEMO_SCRIPTS; do
    # Extract result from log
    if grep -q "✅ $demo_script: PASSED" "$LOG_FILE"; then
        echo "- ✅ $demo_script: PASSED" >> "$REPORT_LOCATION"
    else
        echo "- 🔴 $demo_script: FAILED" >> "$REPORT_LOCATION"
    fi
done

# Add log reference and recommendation
cat >> "$REPORT_LOCATION" << EOF

## Demo Execution Log
See: $LOG_FILE

## Recommendation
${RECOMMENDATION}

---

**Generated by Code Reviewer DEMO_VALIDATION state**
**R291 Gate 4 Enforcement Mechanism**
**Timestamp: $(date -Iseconds)**
EOF

echo "✅ Demo evaluation report created: $REPORT_LOCATION"

# R340: Report validation report location
echo "📋 Validation Report: $REPORT_LOCATION"
echo "📋 Validation Log: $LOG_FILE"
echo "Validation Type: demo_validation"
echo "Integration Type: $INTEGRATE_WAVE_EFFORTS_TYPE"
echo "R340: Created validation report at: $REPORT_LOCATION"
```

### 5. Commit Results

Commit the report and logs to the repository:

```bash
# Return to repository root
cd "$OLDPWD" || cd "$CLAUDE_PROJECT_DIR" || exit 291

# Add report and log
git add "$REPORT_LOCATION" "$LOG_FILE"

# Commit with appropriate message
if [ "$DEMO_STATUS" = "PASSED" ]; then
    git commit -m "demo: ${INTEGRATE_WAVE_EFFORTS_TYPE} demo validation PASSED (R291 Gate 4)

All ${TOTAL_DEMOS} demo(s) executed successfully.
Integration demos verified per R291.

Report: $REPORT_LOCATION
Log: $LOG_FILE"
else
    git commit -m "demo: ${INTEGRATE_WAVE_EFFORTS_TYPE} demo validation FAILED (R291 Gate 4)

${DEMOS_FAILED} demo(s) failed out of ${TOTAL_DEMOS}.
Integration MUST enter ERROR_RECOVERY per R291.

Report: $REPORT_LOCATION
Log: $LOG_FILE"
fi

git push

echo "✅ Demo validation results committed and pushed"
```

## Exit Conditions

### ✅ CONTINUE-SOFTWARE-FACTORY=TRUE
Conditions for successful completion:

**All demos executed** (whether passed or failed):
- Demo scripts found and executed
- Results captured in log
- Report created with structured data
- Results committed to repository
- Orchestrator can read report and enforce R291

**Both PASSED and FAILED demos result in TRUE** because:
- Code Reviewer's job is to RUN demos and REPORT results
- Orchestrator's job is to ENFORCE R291 based on those results
- Code Reviewer has completed its responsibility

### 🔴 CONTINUE-SOFTWARE-FACTORY=FALSE
Conditions requiring manual intervention:
- Demo directory completely missing (no recovery possible)
- Cannot create report (filesystem issues)
- Cannot commit results (git issues)
- Critical infrastructure failure

**Manual intervention needed** for systemic infrastructure issues,
not for demo failures (which are reported normally).

## Next State

**Agent completes** and returns control to orchestrator.

Orchestrator will:
1. Read the demo evaluation report
2. Enforce R291 Gate 4 in WAITING_FOR_DEMO_VALIDATION state
3. Transition to completion state (if passed) or ERROR_RECOVERY (if failed)

## R291 Compliance

This state IS the R291 Gate 4 execution mechanism.

**Responsibilities:**
- ✅ EXECUTE all demo scripts
- ✅ CAPTURE execution results
- ✅ CREATE structured report
- ✅ SAVE execution logs
- ✅ COMMIT results

**NOT Responsible For:**
- ❌ Deciding whether to proceed (orchestrator's job)
- ❌ Triggering ERROR_RECOVERY (orchestrator's job)
- ❌ Creating demos (integration agent's job)

**Per R291:** Demos MUST be executed, not just checked for existence.

This state ensures that requirement is met by ACTUALLY RUNNING the demos.

## Critical Rules Referenced
- **R291**: Integration Demo Requirement (BLOCKING - this state executes demos)
- **R331**: Demo Validation Protocol (BLOCKING - this state validates demos are real, not simulated)
- **R330**: Demo Planning Requirements (demos must be planned)
- **R263**: Integration Documentation Requirements
- **R265**: Integration Testing Requirements

## Remember

**You are the demo executor, not the demo enforcer.**

Your job:
1. Find demos
2. Run demos
3. Report results
4. Let orchestrator enforce

**Do NOT:**
- Skip demos that look "hard to run"
- Assume demos will pass without running them
- Create fake "PASSED" reports
- Try to enforce R291 yourself (orchestrator does that)

**Run ALL demos. Report ALL results. Let the system enforce.**

CONTINUE-SOFTWARE-FACTORY=TRUE
