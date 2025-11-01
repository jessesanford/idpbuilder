# 🔴🔴🔴 RULE R520: PRODUCT MANAGER PRD CREATION STATE

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
**State:** PRD_CREATION
**Enforcement:** MANDATORY - All PRD generation must follow this protocol

---

## 🎯 STATE OBJECTIVE

Generate a comprehensive Product Requirements Document (PRD) from the provided project description, OR create a partial PRD with clearly marked gaps when insufficient information is provided, triggering human intervention via CONTINUE-SOFTWARE-FACTORY=FALSE.

---

## 📋 MANDATORY EXECUTION SEQUENCE

### 1️⃣ LOAD PROJECT PARAMETERS

```bash
# Read from orchestrator state file
PROJECT_NAME=$(jq -r '.project_name' orchestrator-state-v3.json)
PROJECT_DESC=$(jq -r '.project_description' orchestrator-state-v3.json)
PROJECT_TYPE=$(jq -r '.project_type // "service"' orchestrator-state-v3.json)

echo "📋 PRD Generation Parameters:"
echo "   Project Name: $PROJECT_NAME"
echo "   Project Type: $PROJECT_TYPE"
echo "   Description: $PROJECT_DESC"
```

**Validation:**
- [ ] project_name is non-empty
- [ ] project_description is non-empty
- [ ] project_type defaults to "service" if not specified

---

### 2️⃣ SELECT PRD TEMPLATE

```bash
# Template selection based on project type
case "$PROJECT_TYPE" in
    service)
        TEMPLATE="$CLAUDE_PROJECT_DIR/templates/prd-template-service.md"
        ;;
    cli)
        TEMPLATE="$CLAUDE_PROJECT_DIR/templates/prd-template-cli.md"
        ;;
    library)
        TEMPLATE="$CLAUDE_PROJECT_DIR/templates/prd-template-library.md"
        ;;
    web-app)
        TEMPLATE="$CLAUDE_PROJECT_DIR/templates/prd-template-web-app.md"
        ;;
    *)
        # Fallback to service template
        TEMPLATE="$CLAUDE_PROJECT_DIR/templates/prd-template-service.md"
        echo "⚠️ Unknown project type '$PROJECT_TYPE', using service template"
        ;;
esac

echo "📄 Selected template: $TEMPLATE"

# Verify template exists
if [ ! -f "$TEMPLATE" ]; then
    echo "❌ ERROR: Template not found: $TEMPLATE"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=TEMPLATE_NOT_FOUND"
    exit 1
fi

# Read template structure
READ: $TEMPLATE
```

**Template Files:**
- `templates/prd-template-service.md` - Backend services, APIs, microservices
- `templates/prd-template-cli.md` - Command-line tools
- `templates/prd-template-library.md` - Libraries, SDKs, packages
- `templates/prd-template-web-app.md` - Web applications, frontends

---

### 3️⃣ ANALYZE DESCRIPTION COMPLETENESS

**Check for Completeness Indicators:**

```bash
# Analyze the project description for key elements
echo "🔍 Analyzing description completeness..."

# Count specific feature mentions (need ≥3)
FEATURE_COUNT=$(echo "$PROJECT_DESC" | grep -oiE '(should|must|will|can) (do|handle|support|provide|expose|store|validate)' | wc -l)

# Check for technical constraints
HAS_TECH=$(echo "$PROJECT_DESC" | grep -qiE '(go|python|rust|typescript|javascript|java|c\+\+|postgresql|mysql|redis|docker|kubernetes)' && echo "yes" || echo "no")

# Check for scale/performance requirements
HAS_SCALE=$(echo "$PROJECT_DESC" | grep -qiE '([0-9,]+\s*(request|event|user|message|query)s?/(second|minute|hour|day)|[0-9]+%\s*uptime)' && echo "yes" || echo "no")

# Check for stakeholder mentions
HAS_USERS=$(echo "$PROJECT_DESC" | grep -qiE '(user|developer|customer|team|admin|operator)' && echo "yes" || echo "no")

echo "   Features mentioned: $FEATURE_COUNT"
echo "   Tech stack mentioned: $HAS_TECH"
echo "   Scale requirements: $HAS_SCALE"
echo "   Stakeholders identified: $HAS_USERS"
```

**Completeness Decision:**

```
IF:
  - FEATURE_COUNT ≥ 3 AND
  - (HAS_TECH = "yes" OR HAS_SCALE = "yes") AND
  - HAS_USERS = "yes"
THEN:
  → Description is COMPLETE
  → Generate full PRD
  → Set CONTINUE=TRUE
ELSE:
  → Description is INCOMPLETE
  → Generate partial PRD with gaps
  → Create gap analysis report
  → Set CONTINUE=FALSE
```

---

### 4️⃣ GENERATE PRD

#### Option A: Complete Description → Full PRD

```bash
echo "✅ Description is comprehensive - generating complete PRD"

# Create PRD directory
mkdir -p prd

# Populate ALL template sections from description analysis
# Use AI/analysis to fill:
# - Project Overview (from description summary)
# - Problem Statement (from description context)
# - Core Features (from identified capabilities)
# - Success Criteria (from scale/performance requirements)
# - Technical Requirements (from mentioned tech stack)
# - Architecture (infer from project type and features)

# Write complete PRD
cat > "prd/${PROJECT_NAME}-prd.md" <<'EOF'
# ${PROJECT_NAME} - Product Requirements Document

[COMPLETE PRD CONTENT - ALL SECTIONS POPULATED]

## 1. Project Overview
${INFERRED_OVERVIEW}

## 2. Problem Statement
${INFERRED_PROBLEM}

## 3. Core Features
${LISTED_FEATURES}

## 4. Success Criteria
${DEFINED_METRICS}

## 5. Technical Requirements
${TECH_STACK}

... [ALL OTHER SECTIONS] ...
EOF

echo "✅ Complete PRD generated: prd/${PROJECT_NAME}-prd.md"
CONTINUE_FLAG="TRUE"
```

#### Option B: Incomplete Description → Partial PRD + Gap Report

```bash
echo "⚠️ Description is incomplete - generating partial PRD with gaps"

# Create PRD directory
mkdir -p prd

# Populate what's possible, mark gaps with [NEEDS INPUT]
cat > "prd/${PROJECT_NAME}-prd.md" <<'EOF'
# ${PROJECT_NAME} - Product Requirements Document

⚠️ **This PRD is INCOMPLETE** - Sections marked `[NEEDS INPUT]` require human input before proceeding to architecture phase.

## 1. Project Overview
${BASIC_OVERVIEW_FROM_DESCRIPTION}

## 2. Problem Statement
[NEEDS INPUT: What specific problem does this solve? Who experiences this problem?]

## 3. Core Features
[NEEDS INPUT: What are the 3-5 core capabilities this system must provide? Be specific about inputs, outputs, and operations.]

### 3.1 Feature 1: ${INFERRED_FEATURE_IF_ANY}
[NEEDS INPUT: Detailed requirements for this feature]

### 3.2 Feature 2
[NEEDS INPUT: Describe the second core feature]

... [MORE GAPS] ...

## 4. Success Criteria
[NEEDS INPUT: How will we measure success? What metrics indicate this is working?]

## 5. Technical Requirements

### 5.1 Technology Stack
[NEEDS INPUT: What programming language? What database? What frameworks?]

### 5.2 Performance Requirements
[NEEDS INPUT: What scale/throughput/latency requirements?]

... [OTHER SECTIONS WITH GAPS] ...
EOF

# Generate gap analysis report
cat > "prd/PRD-VALIDATION-REPORT.md" <<'EOF'
# PRD VALIDATION REPORT

**Project:** ${PROJECT_NAME}
**Status:** ⚠️ INCOMPLETE - HUMAN INPUT REQUIRED
**Generated:** $(date -Iseconds)

---

## ANALYSIS OF PROVIDED DESCRIPTION

**Description:** "${PROJECT_DESC}"

**What We Identified:**
- ✅ Project type: ${PROJECT_TYPE}
- ✅ General intent: ${BASIC_UNDERSTANDING}
${IDENTIFIED_ELEMENTS}

**What's Missing:**
${GAP_LIST}

---

## REQUIRED INFORMATION

### 1. Core Features (Section 3 of PRD)
**Problem:** Description doesn't specify WHAT the system should do in detail.

**What we need:**
- What specific operations/actions should the system perform?
- What inputs does it accept? What outputs does it produce?
- What are the 3-5 most important capabilities?
- What edge cases must be handled?

**Example:** Instead of "webhook service", specify:
- "Receive GitHub webhook events via HTTP POST"
- "Validate HMAC-SHA256 signatures per GitHub spec"
- "Store events in PostgreSQL with event type indexing"
- "Trigger downstream processing via message queue"

### 2. Technical Requirements (Section 5 of PRD)
**Problem:** No specific technology stack or constraints mentioned.

**What we need:**
- Programming language preference? (Go, Python, Rust, etc.)
- Database/storage requirements? (PostgreSQL, MySQL, Redis, etc.)
- External integrations? (APIs, message queues, etc.)
- Performance/scale requirements? (requests/sec, uptime SLA, etc.)

### 3. Success Criteria (Section 4 of PRD)
**Problem:** Unclear how to measure success of this project.

**What we need:**
- What metrics indicate this is working correctly?
- What does "done" look like?
- Any SLAs or performance targets?
- How will users validate the system meets their needs?

---

## WHAT WE GENERATED

We created a partial PRD with:
- ✅ Standard structure for ${PROJECT_TYPE} projects
- ✅ Common patterns and architectural best practices
- ⚠️ Sections marked `[NEEDS INPUT]` where information is missing

**PRD Location:** `prd/${PROJECT_NAME}-prd.md`

---

## NEXT STEPS

1. **Open the PRD:**
   \`\`\`bash
   cat prd/${PROJECT_NAME}-prd.md
   \`\`\`

2. **Find gaps:**
   \`\`\`bash
   grep -n "\[NEEDS INPUT\]" prd/${PROJECT_NAME}-prd.md
   \`\`\`

3. **Fill in missing details** based on your domain knowledge

4. **Continue orchestration:**
   \`\`\`bash
   /continue-orchestrating
   \`\`\`

   This will spawn the Product Manager again to validate your completed PRD.

---

## REASONING

**Why we stopped:** The description "${PROJECT_DESC}" provides general intent but lacks specific implementation requirements. Without detailed features, technical constraints, and success criteria, the Architect agent cannot create a complete implementation plan.

**What happens next:** After you complete the PRD sections marked `[NEEDS INPUT]`, the Product Manager will validate completeness and proceed to the architecture phase.

---

**Generated by:** Product Manager Agent (PRD_CREATION state)
**Timestamp:** $(date -Iseconds)
EOF

echo "⚠️ Partial PRD generated: prd/${PROJECT_NAME}-prd.md"
echo "📋 Gap report created: prd/PRD-VALIDATION-REPORT.md"
CONTINUE_FLAG="FALSE"
```

---

### 5️⃣ COMMIT PRD ARTIFACTS

```bash
# Add PRD files to git
git add prd/

if [ "$CONTINUE_FLAG" = "TRUE" ]; then
    git commit -m "prd: Generate complete PRD for ${PROJECT_NAME}

Generated from comprehensive project description.
All sections populated and ready for architecture phase.

[PRODUCT_MANAGER_PRD_CREATION]"
else
    git commit -m "prd: Generate partial PRD for ${PROJECT_NAME} - human input required

Generated from minimal project description.
Sections marked [NEEDS INPUT] require human completion.

See: prd/PRD-VALIDATION-REPORT.md for gap analysis.

[PRODUCT_MANAGER_PRD_CREATION]"
fi

git push

echo "✅ PRD artifacts committed and pushed"
```

---

### 6️⃣ STATE TRANSITION SIGNAL

```bash
# Output state completion report
echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "PRODUCT MANAGER PRD CREATION - COMPLETE"
echo "═══════════════════════════════════════════════════════════════"
echo ""
echo "Project: $PROJECT_NAME"
echo "PRD Location: prd/${PROJECT_NAME}-prd.md"

if [ "$CONTINUE_FLAG" = "TRUE" ]; then
    echo "Status: ✅ COMPLETE"
    echo "Description Analysis: Comprehensive - all sections populated"
    echo "Next State: SPAWN_ARCHITECT_MASTER_PLANNING"
else
    echo "Status: ⚠️ INCOMPLETE - HUMAN INPUT REQUIRED"
    echo "Gap Report: prd/PRD-VALIDATION-REPORT.md"
    echo "Next State: WAITING_FOR_PRD_VALIDATION"
    echo ""
    echo "📋 ACTION REQUIRED:"
    echo "   1. Review: prd/PRD-VALIDATION-REPORT.md"
    echo "   2. Edit: prd/${PROJECT_NAME}-prd.md"
    echo "   3. Fill sections marked [NEEDS INPUT]"
    echo "   4. Run: /continue-orchestrating"
fi

echo ""
echo "═══════════════════════════════════════════════════════════════"
echo ""

# CRITICAL: Output continuation flag as LAST line (R405)
echo "CONTINUE-SOFTWARE-FACTORY=$CONTINUE_FLAG"
```

---

## 🚨 CRITICAL VALIDATION CHECKLIST

Before setting CONTINUE=TRUE, verify:
- [ ] PRD file created at `prd/${PROJECT_NAME}-prd.md`
- [ ] NO `[NEEDS INPUT]` markers in mandatory sections
- [ ] ≥3 specific features listed in Core Features section
- [ ] Success criteria defined with measurable outcomes
- [ ] Technical requirements specify language OR framework
- [ ] Git commit created and pushed
- [ ] CONTINUE flag output as LAST line

Before setting CONTINUE=FALSE, verify:
- [ ] PRD file created at `prd/${PROJECT_NAME}-prd.md`
- [ ] Sections marked with `[NEEDS INPUT: specific guidance]`
- [ ] Gap report created at `prd/PRD-VALIDATION-REPORT.md`
- [ ] Gap report clearly identifies what's missing
- [ ] Gap report provides examples of what to add
- [ ] Git commit created and pushed
- [ ] CONTINUE flag output as LAST line

---

## 📏 MANDATORY SECTIONS (Must Be Populated or Marked)

These sections MUST have content OR be marked `[NEEDS INPUT]`:

1. **Project Overview** - 2-3 sentence summary
2. **Problem Statement** - What problem this solves and why
3. **Core Features** - Minimum 3 specific capabilities
4. **Success Criteria** - At least 1 measurable outcome
5. **Technical Requirements** - Language, frameworks, or architecture pattern

If description provides info for a section → populate it
If description lacks info for a section → mark `[NEEDS INPUT: what's needed]`

---

## 🔴 CRITICAL ERROR CASES

### Case 1: Template Not Found
```bash
if [ ! -f "$TEMPLATE" ]; then
    echo "❌ ERROR: PRD template not found: $TEMPLATE"
    echo "Available templates:"
    ls $CLAUDE_PROJECT_DIR/templates/prd-template-*.md
    echo ""
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=TEMPLATE_NOT_FOUND"
    exit 1
fi
```

### Case 2: Empty Project Description
```bash
if [ -z "$PROJECT_DESC" ] || [ "$PROJECT_DESC" = "null" ]; then
    echo "❌ ERROR: No project description provided"
    echo "Cannot generate PRD without description"
    echo ""
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=NO_DESCRIPTION"
    exit 1
fi
```

### Case 3: PRD Directory Creation Fails
```bash
if ! mkdir -p prd; then
    echo "❌ ERROR: Cannot create prd/ directory"
    echo ""
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=FILESYSTEM_ERROR"
    exit 1
fi
```

---

## 💡 BEST PRACTICES

### DO:
✅ Be specific when populating PRD sections
✅ Provide clear guidance in `[NEEDS INPUT]` markers
✅ Create comprehensive gap reports
✅ Use examples in gap reports
✅ Commit PRD files before outputting CONTINUE flag

### DON'T:
❌ Make assumptions for missing critical information
❌ Auto-continue with vague requirements
❌ Write generic fluff ("should be fast and reliable")
❌ Skip gap report when setting CONTINUE=FALSE
❌ Output CONTINUE flag before committing files

---

## 📚 RELATED RULES

- **R405:** Continuation flag must be last output
- **R287:** Save TODOs before state transitions
- **R322:** Mandatory stop after spawning (Part A)
- **R521:** PRD validation protocol (next state)

---

## 🎯 SUCCESS CRITERIA

**State Completes Successfully When:**
- ✅ PRD file created and committed
- ✅ CONTINUE flag correctly set based on description completeness
- ✅ If CONTINUE=FALSE: Gap report created with clear guidance
- ✅ Git push successful
- ✅ CONTINUE flag output as LAST line

**State Fails When:**
- ❌ PRD file not created
- ❌ CONTINUE flag not output
- ❌ Git commit/push fails
- ❌ Template not found
- ❌ No project description provided

---

**This rule is BLOCKING - failure to follow results in initialization failure.**
