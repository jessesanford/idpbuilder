# 🔴🔴🔴 RULE R521: PRODUCT MANAGER PRD VALIDATION STATE

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
**Criticality:** BLOCKING
**Agent:** product-manager
**State:** PRD_VALIDATION
**Enforcement:** MANDATORY - All PRD validation must follow this protocol

---

## 🎯 STATE OBJECTIVE

Validate human-completed PRD to ensure all mandatory sections are populated and `[NEEDS INPUT]` markers have been resolved. Determine if PRD is complete enough to proceed to architecture phase, or if additional human input is required.

---

## 📋 MANDATORY EXECUTION SEQUENCE

### 1️⃣ LOAD PROJECT PARAMETERS

```bash
# Read from orchestrator state file
PROJECT_NAME=$(jq -r '.project_name' orchestrator-state-v3.json)
PRD_FILE="prd/${PROJECT_NAME}-prd.md"

echo "📋 PRD Validation Parameters:"
echo "   Project Name: $PROJECT_NAME"
echo "   PRD File: $PRD_FILE"
```

**Validation:**
- [ ] project_name is non-empty
- [ ] PRD file exists at expected path

---

### 2️⃣ READ HUMAN-EDITED PRD

```bash
# Verify PRD file exists
if [ ! -f "$PRD_FILE" ]; then
    echo "❌ ERROR: PRD file not found: $PRD_FILE"
    echo "Expected PRD at: $PRD_FILE"
    echo ""
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=PRD_FILE_NOT_FOUND"
    exit 1
fi

echo "✅ PRD file found: $PRD_FILE"

# Read PRD content
READ: $PRD_FILE

echo "📄 PRD loaded successfully"
```

---

### 3️⃣ CHECK FOR REMAINING GAPS

```bash
echo "🔍 Checking for [NEEDS INPUT] markers..."

# Count total [NEEDS INPUT] markers
TOTAL_GAPS=$(grep -c "\[NEEDS INPUT" "$PRD_FILE" || echo "0")

echo "   Total [NEEDS INPUT] markers found: $TOTAL_GAPS"

if [ "$TOTAL_GAPS" -eq 0 ]; then
    echo "   ✅ No gaps remaining - PRD appears complete"
    GAPS_REMAIN="no"
else
    echo "   ⚠️ Gaps still present - listing locations:"
    grep -n "\[NEEDS INPUT" "$PRD_FILE" | head -10
    GAPS_REMAIN="yes"
fi
```

---

### 4️⃣ VALIDATE MANDATORY SECTIONS

```bash
echo "📋 Validating mandatory sections..."

# Section 1: Project Overview
if grep -qA 5 "## 1. Project Overview" "$PRD_FILE" | grep -qv "\[NEEDS INPUT\]"; then
    echo "   ✅ Project Overview: Complete"
    OVERVIEW_OK="yes"
else
    echo "   ❌ Project Overview: Missing or incomplete"
    OVERVIEW_OK="no"
fi

# Section 2: Problem Statement
if grep -qA 5 "## 2. Problem Statement" "$PRD_FILE" | grep -qv "\[NEEDS INPUT\]"; then
    echo "   ✅ Problem Statement: Complete"
    PROBLEM_OK="yes"
else
    echo "   ❌ Problem Statement: Missing or incomplete"
    PROBLEM_OK="no"
fi

# Section 3: Core Features (must have ≥3 features)
FEATURE_COUNT=$(grep -c "### 3\.[0-9]" "$PRD_FILE" || echo "0")
if [ "$FEATURE_COUNT" -ge 3 ]; then
    # Check if features have content (not just [NEEDS INPUT])
    if grep -A 10 "## 3. Core Features" "$PRD_FILE" | grep -qv "\[NEEDS INPUT\]"; then
        echo "   ✅ Core Features: $FEATURE_COUNT features defined"
        FEATURES_OK="yes"
    else
        echo "   ❌ Core Features: Features listed but not detailed"
        FEATURES_OK="no"
    fi
else
    echo "   ❌ Core Features: Only $FEATURE_COUNT features (need ≥3)"
    FEATURES_OK="no"
fi

# Section 4: Success Criteria
if grep -qA 5 "## 4. Success Criteria" "$PRD_FILE" | grep -qv "\[NEEDS INPUT\]"; then
    echo "   ✅ Success Criteria: Defined"
    SUCCESS_OK="yes"
else
    echo "   ❌ Success Criteria: Missing or incomplete"
    SUCCESS_OK="no"
fi

# Section 5: Technical Requirements
if grep -qA 10 "## 5. Technical Requirements" "$PRD_FILE" | grep -qv "\[NEEDS INPUT\]"; then
    echo "   ✅ Technical Requirements: Specified"
    TECH_OK="yes"
else
    echo "   ❌ Technical Requirements: Missing or incomplete"
    TECH_OK="no"
fi
```

---

### 5️⃣ MAKE VALIDATION DECISION

```bash
echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "VALIDATION DECISION"
echo "═══════════════════════════════════════════════════════════════"

# All mandatory sections must be complete
if [ "$OVERVIEW_OK" = "yes" ] && \
   [ "$PROBLEM_OK" = "yes" ] && \
   [ "$FEATURES_OK" = "yes" ] && \
   [ "$SUCCESS_OK" = "yes" ] && \
   [ "$TECH_OK" = "yes" ] && \
   [ "$GAPS_REMAIN" = "no" ]; then

    echo "✅ VALIDATION PASSED"
    echo ""
    echo "All mandatory sections complete:"
    echo "  ✅ Project Overview"
    echo "  ✅ Problem Statement"
    echo "  ✅ Core Features (≥3)"
    echo "  ✅ Success Criteria"
    echo "  ✅ Technical Requirements"
    echo "  ✅ No [NEEDS INPUT] markers remaining"
    echo ""
    echo "Decision: Proceed to architecture phase"
    VALIDATION_STATUS="PASS"
    CONTINUE_FLAG="TRUE"
else
    echo "❌ VALIDATION FAILED"
    echo ""
    echo "Missing or incomplete sections:"
    [ "$OVERVIEW_OK" = "no" ] && echo "  ❌ Project Overview"
    [ "$PROBLEM_OK" = "no" ] && echo "  ❌ Problem Statement"
    [ "$FEATURES_OK" = "no" ] && echo "  ❌ Core Features (need ≥3 detailed features)"
    [ "$SUCCESS_OK" = "no" ] && echo "  ❌ Success Criteria"
    [ "$TECH_OK" = "no" ] && echo "  ❌ Technical Requirements"
    [ "$GAPS_REMAIN" = "yes" ] && echo "  ❌ [NEEDS INPUT] markers still present ($TOTAL_GAPS)"
    echo ""
    echo "Decision: Human must complete remaining sections"
    VALIDATION_STATUS="NEEDS_REVISION"
    CONTINUE_FLAG="FALSE"
fi

echo "═══════════════════════════════════════════════════════════════"
```

---

### 6️⃣ GENERATE VALIDATION REPORT

```bash
echo "📋 Generating validation report..."

if [ "$VALIDATION_STATUS" = "PASS" ]; then
    # Validation passed - create success report
    cat > "prd/PRD-VALIDATION-REPORT.md" <<EOF
# PRD VALIDATION REPORT

**Project:** ${PROJECT_NAME}
**Status:** ✅ VALIDATION PASSED
**Validated:** $(date -Iseconds)
**Validated By:** Product Manager Agent (PRD_VALIDATION state)

---

## VALIDATION RESULTS

**All mandatory sections complete and validated:**

- ✅ **Project Overview**: Clear summary provided
- ✅ **Problem Statement**: Problem and rationale defined
- ✅ **Core Features**: ${FEATURE_COUNT} features specified with details
- ✅ **Success Criteria**: Measurable outcomes defined
- ✅ **Technical Requirements**: Technology stack and constraints specified

**Quality Checks:**
- ✅ No \`[NEEDS INPUT]\` markers remaining
- ✅ Features are specific and implementable
- ✅ Requirements are testable

---

## NEXT STEPS

**PRD is ready for architecture phase.**

The Architect agent will now:
1. Read this PRD as primary input
2. Design system architecture
3. Create PROJECT-ARCHITECTURE.md
4. Generate PROJECT-IMPLEMENTATION-PLAN.md

---

**Proceeding to:** SPAWN_ARCHITECT_MASTER_PLANNING

**Generated:** $(date -Iseconds)
EOF

else
    # Validation failed - create gap report
    cat > "prd/PRD-VALIDATION-REPORT.md" <<EOF
# PRD VALIDATION REPORT

**Project:** ${PROJECT_NAME}
**Status:** ❌ VALIDATION FAILED - ADDITIONAL INPUT REQUIRED
**Validated:** $(date -Iseconds)
**Validated By:** Product Manager Agent (PRD_VALIDATION state)

---

## VALIDATION RESULTS

**Incomplete sections detected:**

EOF

    # Add specific failures
    [ "$OVERVIEW_OK" = "no" ] && cat >> "prd/PRD-VALIDATION-REPORT.md" <<EOF
- ❌ **Project Overview** (Section 1)
  - **Problem**: Section is empty or contains \`[NEEDS INPUT]\`
  - **Required**: 2-3 sentence summary of the project
  - **Example**: "This project builds a GitHub webhook receiver that validates payloads, stores events, and triggers downstream processing for DevOps teams."

EOF

    [ "$PROBLEM_OK" = "no" ] && cat >> "prd/PRD-VALIDATION-REPORT.md" <<EOF
- ❌ **Problem Statement** (Section 2)
  - **Problem**: Section is empty or contains \`[NEEDS INPUT]\`
  - **Required**: Clear description of what problem this solves and why
  - **Example**: "Current CI/CD pipelines lack real-time GitHub event processing, causing delays in build triggers. This webhook receiver provides instant event capture and routing."

EOF

    [ "$FEATURES_OK" = "no" ] && cat >> "prd/PRD-VALIDATION-REPORT.md" <<EOF
- ❌ **Core Features** (Section 3)
  - **Problem**: Less than 3 features OR features contain \`[NEEDS INPUT]\`
  - **Found**: ${FEATURE_COUNT} feature(s)
  - **Required**: At least 3 specific, detailed features
  - **Example**:
    ### 3.1 Webhook Reception
    - Accept HTTP POST requests from GitHub
    - Validate HMAC-SHA256 signatures
    - Return appropriate HTTP status codes

    ### 3.2 Event Storage
    - Store events in PostgreSQL with event type indexing
    - Maintain event order and timestamps
    - Support queries by repository, event type, and time range

    ### 3.3 Metrics Exposure
    - Expose Prometheus metrics for event counts
    - Track processing latency (p50, p95, p99)
    - Monitor validation failures and storage errors

EOF

    [ "$SUCCESS_OK" = "no" ] && cat >> "prd/PRD-VALIDATION-REPORT.md" <<EOF
- ❌ **Success Criteria** (Section 4)
  - **Problem**: Section is empty or contains \`[NEEDS INPUT]\`
  - **Required**: Measurable outcomes that indicate success
  - **Example**:
    - Handle 1000 events/minute with <100ms p95 latency
    - 99.9% uptime over 30-day periods
    - Zero data loss (all events stored)
    - 100% signature validation accuracy

EOF

    [ "$TECH_OK" = "no" ] && cat >> "prd/PRD-VALIDATION-REPORT.md" <<EOF
- ❌ **Technical Requirements** (Section 5)
  - **Problem**: Section is empty or contains \`[NEEDS INPUT]\`
  - **Required**: Technology stack, frameworks, architecture patterns
  - **Example**:
    - Language: Go 1.21+
    - Database: PostgreSQL 15+
    - Metrics: Prometheus client library
    - Architecture: HTTP handler → validator → storage → metrics

EOF

    cat >> "prd/PRD-VALIDATION-REPORT.md" <<EOF
---

## REMAINING GAPS

Found ${TOTAL_GAPS} \`[NEEDS INPUT]\` marker(s) in PRD:

\`\`\`bash
$(grep -n "\[NEEDS INPUT" "$PRD_FILE" | head -10)
\`\`\`

---

## NEXT STEPS

1. **Open the PRD:**
   \`\`\`bash
   \$EDITOR prd/${PROJECT_NAME}-prd.md
   \`\`\`

2. **Complete the sections listed above** with specific, detailed information

3. **Remove all \`[NEEDS INPUT]\` markers** by replacing them with actual content

4. **Run validation again:**
   \`\`\`bash
   /continue-orchestrating
   \`\`\`

   This will re-run PRD validation and proceed if all sections are complete.

---

## TIPS FOR COMPLETING PRD

- **Be Specific**: "Fast" → "Handle 1000 requests/sec"
- **Be Testable**: Every feature should have acceptance criteria
- **Be Realistic**: Base on actual project requirements
- **Be Complete**: Don't leave placeholders or TODOs

---

**Status:** Returning to WAITING_FOR_PRD_VALIDATION state

**Generated:** $(date -Iseconds)
EOF

fi

echo "✅ Validation report created: prd/PRD-VALIDATION-REPORT.md"
```

---

### 7️⃣ COMMIT VALIDATION RESULTS

```bash
# Add validation report to git
git add prd/PRD-VALIDATION-REPORT.md

if [ "$VALIDATION_STATUS" = "PASS" ]; then
    git commit -m "prd: Validation PASSED for ${PROJECT_NAME}

All mandatory sections complete and validated.
Ready for architecture phase.

[PRODUCT_MANAGER_PRD_VALIDATION]"
else
    git commit -m "prd: Validation FAILED for ${PROJECT_NAME} - additional input required

Missing or incomplete sections identified.
See: prd/PRD-VALIDATION-REPORT.md

[PRODUCT_MANAGER_PRD_VALIDATION]"
fi

git push

echo "✅ Validation results committed and pushed"
```

---

### 8️⃣ STATE TRANSITION SIGNAL

```bash
# Output state completion report
echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "PRODUCT MANAGER PRD VALIDATION - COMPLETE"
echo "═══════════════════════════════════════════════════════════════"
echo ""
echo "Project: $PROJECT_NAME"
echo "PRD File: $PRD_FILE"
echo "Validation Report: prd/PRD-VALIDATION-REPORT.md"
echo ""

if [ "$VALIDATION_STATUS" = "PASS" ]; then
    echo "Status: ✅ VALIDATION PASSED"
    echo "Mandatory Sections: ALL COMPLETE"
    echo "Remaining Gaps: NONE"
    echo ""
    echo "Next State: SPAWN_ARCHITECT_MASTER_PLANNING"
    echo ""
    echo "🎉 PRD is ready for architecture phase!"
else
    echo "Status: ❌ VALIDATION FAILED"
    echo "Mandatory Sections: INCOMPLETE"
    echo "Remaining Gaps: $TOTAL_GAPS [NEEDS INPUT] markers"
    echo ""
    echo "Next State: WAITING_FOR_PRD_VALIDATION (loop)"
    echo ""
    echo "📋 ACTION REQUIRED:"
    echo "   1. Review: prd/PRD-VALIDATION-REPORT.md"
    echo "   2. Edit: prd/${PROJECT_NAME}-prd.md"
    echo "   3. Complete sections listed in validation report"
    echo "   4. Run: /continue-orchestrating (to re-validate)"
fi

echo ""
echo "═══════════════════════════════════════════════════════════════"
echo ""

# CRITICAL: Output continuation flag as LAST line (R405)
echo "CONTINUE-SOFTWARE-FACTORY=$CONTINUE_FLAG"
```

---

## 🚨 CRITICAL VALIDATION CHECKLIST

Before setting CONTINUE=TRUE (validation pass):
- [ ] PRD file exists and is readable
- [ ] NO `[NEEDS INPUT]` markers anywhere in file
- [ ] Project Overview section has content
- [ ] Problem Statement section has content
- [ ] ≥3 features listed in Core Features section
- [ ] Each feature has detailed description
- [ ] Success Criteria section has measurable outcomes
- [ ] Technical Requirements section specifies tech stack
- [ ] Validation report created (success)
- [ ] Git commit created and pushed
- [ ] CONTINUE flag output as LAST line

Before setting CONTINUE=FALSE (validation fail):
- [ ] PRD file exists and is readable
- [ ] Specific missing sections identified
- [ ] Validation report created with clear guidance
- [ ] Validation report lists ALL incomplete sections
- [ ] Validation report provides examples for each gap
- [ ] Git commit created and pushed
- [ ] CONTINUE flag output as LAST line

---

## 📏 MANDATORY SECTION REQUIREMENTS

### 1. Project Overview
- **Minimum**: 2-3 sentences
- **Content**: What the project is, what it does, who uses it
- **Test**: `grep -A 5 "## 1. Project Overview" | wc -w` should be ≥20 words

### 2. Problem Statement
- **Minimum**: 3-5 sentences
- **Content**: What problem this solves, why it matters, current gaps
- **Test**: Section should not contain `[NEEDS INPUT]`

### 3. Core Features
- **Minimum**: 3 features
- **Content**: Each feature with detailed description, inputs/outputs, edge cases
- **Test**: `grep -c "### 3\.[0-9]"` should be ≥3
- **Test**: Features should not contain `[NEEDS INPUT]`

### 4. Success Criteria
- **Minimum**: 1 measurable metric
- **Content**: How success is measured, target values, acceptance criteria
- **Test**: Should mention numbers, percentages, or measurable outcomes

### 5. Technical Requirements
- **Minimum**: Language OR framework mentioned
- **Content**: Tech stack, architecture pattern, key dependencies
- **Test**: Should mention specific technologies (Go, PostgreSQL, etc.)

---

## 🔴 CRITICAL ERROR CASES

### Case 1: PRD File Not Found
```bash
if [ ! -f "$PRD_FILE" ]; then
    echo "❌ ERROR: PRD file not found: $PRD_FILE"
    echo ""
    echo "Expected PRD at: prd/${PROJECT_NAME}-prd.md"
    echo "This file should have been created in PRD_CREATION state"
    echo "and edited by human in WAITING_FOR_PRD_VALIDATION state"
    echo ""
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=PRD_FILE_MISSING"
    exit 1
fi
```

### Case 2: PRD File Empty
```bash
if [ ! -s "$PRD_FILE" ]; then
    echo "❌ ERROR: PRD file is empty: $PRD_FILE"
    echo ""
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=PRD_FILE_EMPTY"
    exit 1
fi
```

### Case 3: Validation Loop Limit
```bash
# Check how many validation attempts
VALIDATION_COUNT=$(jq -r '.prd_validation_attempts // 0' orchestrator-state-v3.json)

if [ "$VALIDATION_COUNT" -ge 5 ]; then
    echo "⚠️ WARNING: PRD validation attempted $VALIDATION_COUNT times"
    echo "Consider reviewing requirements with stakeholders"
fi

# Update count
jq ".prd_validation_attempts = $((VALIDATION_COUNT + 1))" orchestrator-state-v3.json > tmp.json
mv tmp.json orchestrator-state-v3.json
```

---

## 💡 BEST PRACTICES

### DO:
✅ Provide specific, actionable feedback for each gap
✅ Include examples in validation report
✅ Be thorough in checking all mandatory sections
✅ Use quantitative checks (word counts, feature counts)
✅ Create clear next steps for human

### DON'T:
❌ Auto-pass PRDs with vague sections
❌ Accept "TBD" or "TODO" as valid content
❌ Skip validation report generation
❌ Fail to identify specific missing sections
❌ Output CONTINUE flag before committing files

---

## 📚 RELATED RULES

- **R405:** Continuation flag must be last output
- **R287:** Save TODOs before state transitions
- **R322:** Mandatory stop after spawning (Part A)
- **R520:** PRD creation protocol (previous state)

---

## 🎯 SUCCESS CRITERIA

**State Completes Successfully When:**
- ✅ PRD file read and analyzed
- ✅ All mandatory sections validated
- ✅ CONTINUE flag correctly set based on validation results
- ✅ Validation report created with detailed feedback
- ✅ Git push successful
- ✅ CONTINUE flag output as LAST line

**State Loops When:**
- 🔄 Validation fails (CONTINUE=FALSE)
- 🔄 Returns to WAITING_FOR_PRD_VALIDATION
- 🔄 Human completes remaining sections
- 🔄 Re-runs validation until pass

**State Fails When:**
- ❌ PRD file not found
- ❌ PRD file empty
- ❌ Validation report not created
- ❌ Git commit/push fails
- ❌ CONTINUE flag not output

---

**This rule is BLOCKING - failure to follow results in initialization failure.**
