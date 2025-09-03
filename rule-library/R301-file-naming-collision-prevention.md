# 🚨🚨🚨 RULE R301: FILE NAMING COLLISION PREVENTION (BLOCKING)

## Rule Definition
**Criticality:** BLOCKING
**Category:** Universal - All Agents
**Priority:** HIGH
**Penalty:** -50% for violations causing data loss

## PROBLEM STATEMENT

Multiple agents working in parallel or sequentially can create files with identical names, causing:
- **Data Loss**: Files overwritten silently
- **Merge Conflicts**: Git cannot auto-merge identical filenames
- **Integration Failures**: Wrong version of file used
- **Audit Trail Loss**: Historical work logs destroyed

## MANDATORY NAMING REQUIREMENTS

### 1. HIGH-RISK FILES (MUST USE TIMESTAMPS)

These files MUST include timestamps to prevent collisions:

#### Integration Files
```bash
# Pattern: {BASE-NAME}-{OPERATION-ID}-{TIMESTAMP}.md
INTEGRATION-REPORT-wave1-20250120-143000.md
INTEGRATION-PLAN-phase1-20250120-145500.md
work-log-integration-20250120-150000.md
```

#### Review Reports
```bash
# Pattern: {REPORT-TYPE}-{EFFORT}-{TIMESTAMP}.md
CODE-REVIEW-REPORT-auth-system-20250120-143000.md
SPLIT-PLAN-user-mgmt-20250120-145500.md
ARCHITECTURE-REVIEW-wave1-20250120-160000.md
```

#### Wave/Phase Level Documents
```bash
# Pattern: {SCOPE}-{TYPE}-{IDENTIFIER}-{TIMESTAMP}.md
WAVE-IMPLEMENTATION-PLAN-wave1-20250120.md
PHASE-ARCHITECTURE-PLAN-phase1-20250120.md
```

### 2. TIMESTAMP FORMATS

Use ONE of these formats consistently:

```bash
# Full timestamp (when multiple per day possible)
TIMESTAMP=$(date +%Y%m%d-%H%M%S)  # 20250120-143000

# Daily timestamp (for once-per-day files)
TIMESTAMP=$(date +%Y%m%d)         # 20250120

# Epoch (for programmatic generation)
TIMESTAMP=$(date +%s)              # 1705759800
```

### 3. EFFORT-SCOPED FILES (PROTECTED BY DIRECTORY)

These files are naturally protected by directory structure:
- `efforts/phase1/wave1/auth-system/work-log.md` ✅
- `efforts/phase1/wave1/user-mgmt/work-log.md` ✅

**NO TIMESTAMP REQUIRED** when file is in unique effort directory.

## IMPLEMENTATION PATTERNS

### Pattern 1: Integration Report Creation
```bash
# ✅ CORRECT - Includes timestamp
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
REPORT_NAME="INTEGRATION-REPORT-${WAVE_ID}-${TIMESTAMP}.md"

cat > "$REPORT_NAME" << 'EOF'
# Integration Report
...
EOF

# ❌ WRONG - No timestamp, will overwrite
cat > "INTEGRATION-REPORT.md" << 'EOF'
...
EOF
```

### Pattern 2: Code Review Report
```bash
# ✅ CORRECT - Includes effort and timestamp
create_review_report() {
    local effort_name="$1"
    local timestamp=$(date +%Y%m%d-%H%M%S)
    local report_file="CODE-REVIEW-REPORT-${effort_name}-${timestamp}.md"
    
    echo "# Code Review Report" > "$report_file"
    echo "Effort: $effort_name" >> "$report_file"
    echo "Date: $(date)" >> "$report_file"
}

# ❌ WRONG - Generic name
echo "# Code Review" > "CODE-REVIEW-REPORT.md"
```

### Pattern 3: Work Log in Shared Space
```bash
# When creating work logs outside effort directories
# ✅ CORRECT - Unique identifier
WORK_LOG="work-log-${AGENT_TYPE}-${OPERATION}-$(date +%Y%m%d-%H%M%S).md"

# ❌ WRONG - Generic name in shared space
touch "work-log.md"  # COLLISION RISK!
```

## COLLISION DETECTION REQUIREMENT

Before creating high-risk files, CHECK for existing files:

```bash
check_file_collision() {
    local base_name="$1"
    local directory="$2"
    
    # Check if files with this base name exist
    if ls "${directory}"/${base_name}* 2>/dev/null | grep -q .; then
        echo "⚠️ WARNING: Files with base name '${base_name}' already exist:"
        ls -la "${directory}"/${base_name}*
        echo "📝 Adding timestamp to prevent collision..."
        return 1
    fi
    return 0
}

# Usage
if ! check_file_collision "INTEGRATION-REPORT" "."; then
    REPORT_FILE="INTEGRATION-REPORT-$(date +%Y%m%d-%H%M%S).md"
else
    REPORT_FILE="INTEGRATION-REPORT.md"
fi
```

## EXEMPTIONS

These files are EXEMPT from timestamp requirements:

1. **Effort-scoped files** in unique directories
2. **Git-managed files** (README.md, .gitignore, etc.)
3. **Configuration files** with standard names
4. **Source code files** with unique paths

## ENFORCEMENT

```bash
# Validation function
validate_file_naming() {
    local file_path="$1"
    local file_name=$(basename "$file_path")
    
    # List of high-risk patterns
    local high_risk_patterns=(
        "INTEGRATION-REPORT"
        "CODE-REVIEW-REPORT"
        "SPLIT-PLAN"
        "WAVE-.*-PLAN"
        "PHASE-.*-PLAN"
        "work-log"  # When not in effort directory
    )
    
    # Check if file matches high-risk pattern
    for pattern in "${high_risk_patterns[@]}"; do
        if [[ "$file_name" =~ ^${pattern}\.md$ ]]; then
            echo "❌ VIOLATION: High-risk file '$file_name' lacks timestamp!"
            echo "   Required format: ${pattern}-[identifier]-[timestamp].md"
            return 1
        fi
    done
    
    # Check if file has timestamp
    if [[ "$file_name" =~ [0-9]{8}(-[0-9]{6})?\.md$ ]]; then
        echo "✅ File '$file_name' includes timestamp"
        return 0
    fi
}
```

## GRADING IMPACT

- **-50%** - File collision causing data loss
- **-30%** - Missing timestamps on high-risk files
- **-20%** - Inconsistent timestamp formats
- **-10%** - No collision detection before file creation

## MIGRATION FOR EXISTING FILES

When encountering existing files without timestamps:

```bash
# Archive with timestamp
mv INTEGRATION-REPORT.md INTEGRATION-REPORT-LEGACY-$(date +%Y%m%d-%H%M%S).md

# Create new with proper naming
touch INTEGRATION-REPORT-wave1-$(date +%Y%m%d-%H%M%S).md
```

## EXAMPLES

### Good Examples ✅
```bash
INTEGRATION-REPORT-wave1-20250120-143000.md
CODE-REVIEW-REPORT-auth-system-20250120-145500.md
SPLIT-PLAN-user-mgmt-split1-20250120-150000.md
work-log-integration-phase1-20250120-160000.md
WAVE-IMPLEMENTATION-PLAN-wave2-20250120.md
```

### Bad Examples ❌
```bash
INTEGRATION-REPORT.md           # No timestamp
CODE-REVIEW-REPORT.md          # No effort identifier
SPLIT-PLAN.md                  # No context
work-log.md                    # Too generic
WAVE-IMPLEMENTATION-PLAN.md    # No wave identifier
```

## RELATED RULES

- **R287** - TODO Persistence (uses timestamps)
- **R293** - Integration Report Distribution (archives with timestamps)
- **R264** - Work Log Tracking Requirements
- **R263** - Integration Documentation Requirements
- **R239** - Fix Plan Distribution

## RATIONALE

This rule prevents the **#3 cause of Software Factory failures**: file collisions during parallel agent operations. By requiring timestamps on high-risk files, we ensure:

1. **No data loss** from overwrites
2. **Clear audit trail** of all operations
3. **Successful integration** without conflicts
4. **Parallel agent safety** during concurrent work

Without this rule, integration attempts fail 40% of the time due to file conflicts.