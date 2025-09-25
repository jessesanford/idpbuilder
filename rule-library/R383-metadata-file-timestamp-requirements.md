# 🔴🔴🔴 RULE R383: Software Factory Metadata File Organization (SUPREME LAW) 🔴🔴🔴

## Rule Identifier
**Rule ID**: R383
**Category**: Merge Conflict Prevention
**Criticality**: SUPREME LAW - Violations cause integration failures
**Priority**: ABSOLUTE
**Penalty**: -100% for ANY violation

## Description

ALL Software Factory metadata files MUST include unique timestamps in their filenames and be stored in standardized .software-factory directories. This PREVENTS merge conflicts during integration and ensures perfect file uniqueness across all agents and operations.

## Rationale

Without unique timestamps:
- Multiple agents create files with same names
- Integration causes merge conflicts
- Metadata gets overwritten/lost
- Parallel work becomes impossible
- Recovery after failures is blocked

## ABSOLUTE REQUIREMENTS

### 1. MANDATORY TIMESTAMP FORMAT

**EVERY metadata file MUST have timestamp suffix:**

```bash
# CORRECT FORMAT: filename--YYYYMMDD-HHMMSS.ext
IMPLEMENTATION-PLAN--20250121-143052.md
CODE-REVIEW-REPORT--20250121-153427.md
work-log--20250121-163915.log
integration-report--20250121-173245.md

# ❌ FORBIDDEN: No timestamp
IMPLEMENTATION-PLAN.md           # VIOLATION!
CODE-REVIEW-REPORT.md            # VIOLATION!
work-log.md                      # VIOLATION!

# ❌ FORBIDDEN: Incorrect timestamp format
IMPLEMENTATION-PLAN-2025-01-21.md    # WRONG FORMAT!
CODE-REVIEW-20250121.md              # MISSING TIME!
work-log_143052.md                   # WRONG SEPARATOR!
```

### 2. TIMESTAMP GENERATION PROTOCOL

**ALWAYS use this exact pattern:**

```bash
# CORRECT: Consistent timestamp generation
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
FILENAME="IMPLEMENTATION-PLAN--${TIMESTAMP}.md"

# ❌ WRONG: Various incorrect patterns
TIMESTAMP=$(date +%Y-%m-%d)              # Wrong format
FILENAME="PLAN-${date}.md"               # Missing time
FILENAME="PLAN-$(date +%s).md"           # Unix timestamp
```

### 3. METADATA DIRECTORY STRUCTURE (from R343)

**ALL metadata MUST be in .software-factory directories:**

```bash
# CORRECT: Proper structure with timestamps
.software-factory/
├── phase1/
│   └── wave1/
│       └── effort-name/
│           ├── IMPLEMENTATION-PLAN--20250121-143052.md
│           ├── CODE-REVIEW-REPORT--20250121-153427.md
│           ├── work-log--20250121-163915.log
│           └── validation-results--20250121-173245.md

# ❌ FORBIDDEN: Files without timestamps
.software-factory/
├── IMPLEMENTATION-PLAN.md         # NO TIMESTAMP!
├── CODE-REVIEW-REPORT.md          # NO TIMESTAMP!
└── work-log.md                    # NO TIMESTAMP!
```

### 4. COMPREHENSIVE METADATA FILE LIST

**ALL these files REQUIRE timestamps:**

- `IMPLEMENTATION-PLAN--*.md` - From Code Reviewer
- `CODE-REVIEW-REPORT--*.md` - Review results
- `SPLIT-PLAN--*.md` - Split strategies
- `FIX-PLAN--*.md` - Fix requirements
- `work-log--*.log` - Development logs
- `validation-results--*.md` - Test results
- `INTEGRATION-REPORT--*.md` - Integration results
- `WAVE-MERGE-PLAN--*.md` - Wave merge strategies
- `PHASE-MERGE-PLAN--*.md` - Phase merge strategies
- `PROJECT-ARCHITECTURE--*.md` - Architecture docs
- `WAVE-ARCHITECTURE--*.md` - Wave architecture
- `PHASE-ARCHITECTURE--*.md` - Phase architecture
- `TEST-PLAN--*.md` - Test strategies
- `DEMO-PLAN--*.md` - Demo scenarios
- `progress-marker--*.marker` - Progress tracking
- `status-update--*.status` - Status files
- `metadata--*.yaml` - Metadata files
- `metadata--*.json` - JSON metadata
- **ANY OTHER NON-CODE FILE** - Must have timestamp!

### 5. HELPER FUNCTION (MANDATORY USE)

**ALL agents MUST use this function:**

```bash
# MANDATORY: Include this in ALL agent scripts
sf_metadata_path() {
    local phase="$1"
    local wave="$2"
    local effort="$3"
    local filename="$4"
    local ext="$5"

    # Validate inputs
    if [[ -z "$phase" || -z "$wave" || -z "$effort" || -z "$filename" || -z "$ext" ]]; then
        echo "❌ R383 VIOLATION: Missing parameters to sf_metadata_path" >&2
        exit 1
    fi

    # Create directory structure
    local dir=".software-factory/phase${phase}/wave${wave}/${effort}"
    mkdir -p "$dir"

    # Generate unique timestamped filename
    local timestamp=$(date +%Y%m%d-%H%M%S)
    local full_path="${dir}/${filename}--${timestamp}.${ext}"

    echo "$full_path"
}

# USAGE EXAMPLES:
# Create implementation plan
PLAN_PATH=$(sf_metadata_path 1 2 "registry-client" "IMPLEMENTATION-PLAN" "md")
echo "# Implementation Plan" > "$PLAN_PATH"

# Create review report
REPORT_PATH=$(sf_metadata_path 1 2 "registry-client" "CODE-REVIEW-REPORT" "md")
echo "# Code Review Report" > "$REPORT_PATH"

# Create work log
LOG_PATH=$(sf_metadata_path 1 2 "registry-client" "work-log" "log")
echo "Starting work..." >> "$LOG_PATH"
```

### 6. AGENT-SPECIFIC REQUIREMENTS

#### Orchestrator
```bash
# Creating status update
STATUS_PATH=$(sf_metadata_path "$PHASE" "$WAVE" "orchestration" "status-update" "json")
echo '{"status": "in_progress"}' > "$STATUS_PATH"
```

#### Code Reviewer
```bash
# Creating implementation plan
PLAN_PATH=$(sf_metadata_path "$PHASE" "$WAVE" "$EFFORT" "IMPLEMENTATION-PLAN" "md")
cat > "$PLAN_PATH" << 'EOF'
# Implementation Plan
...
EOF
```

#### SW Engineer
```bash
# Creating work log
LOG_PATH=$(sf_metadata_path "$PHASE" "$WAVE" "$EFFORT" "work-log" "log")
echo "[$(date)] Starting implementation" >> "$LOG_PATH"
```

#### Integration Agent
```bash
# Creating integration report
REPORT_PATH=$(sf_metadata_path "$PHASE" "$WAVE" "integration" "INTEGRATION-REPORT" "md")
echo "# Integration Report" > "$REPORT_PATH"
```

### 7. VALIDATION REQUIREMENTS

```bash
# Validation function - MUST be called
validate_metadata_filename() {
    local filepath="$1"
    local filename=$(basename "$filepath")

    # Check for timestamp pattern
    if [[ ! "$filename" =~ --[0-9]{8}-[0-9]{6}\. ]]; then
        echo "❌ R383 VIOLATION: Missing timestamp in $filename"
        echo "   Required format: name--YYYYMMDD-HHMMSS.ext"
        exit 1
    fi

    # Check directory structure
    if [[ ! "$filepath" =~ \.software-factory/ ]]; then
        echo "❌ R383 VIOLATION: File not in .software-factory directory"
        echo "   File: $filepath"
        exit 1
    fi

    echo "✅ R383 COMPLIANT: $filename"
}
```

### 8. MIGRATION FROM OLD PATTERNS

**When encountering old files:**

```bash
# Migrate old files to new format
migrate_to_timestamped() {
    local old_file="$1"

    if [[ -f "$old_file" && ! "$old_file" =~ --[0-9]{8}-[0-9]{6}\. ]]; then
        # Extract base name and extension
        local base="${old_file%.*}"
        local ext="${old_file##*.}"

        # Create new timestamped name
        local timestamp=$(date +%Y%m%d-%H%M%S)
        local new_file="${base}--${timestamp}.${ext}"

        # Move file
        mv "$old_file" "$new_file"
        echo "✅ Migrated: $old_file -> $new_file"
    fi
}
```

## Common Violations

### ❌ VIOLATION 1: No Timestamp
```bash
# WRONG
echo "# Plan" > .software-factory/IMPLEMENTATION-PLAN.md

# RIGHT
PLAN_PATH=$(sf_metadata_path 1 1 "effort" "IMPLEMENTATION-PLAN" "md")
echo "# Plan" > "$PLAN_PATH"
```

### ❌ VIOLATION 2: Wrong Timestamp Format
```bash
# WRONG - Using different date format
echo "# Plan" > "PLAN-$(date +%Y-%m-%d).md"

# RIGHT - Using standard format
echo "# Plan" > "PLAN--$(date +%Y%m%d-%H%M%S).md"
```

### ❌ VIOLATION 3: Not Using Helper Function
```bash
# WRONG - Manual path construction
mkdir -p .software-factory
echo "# Plan" > ".software-factory/plan.md"

# RIGHT - Using helper function
PLAN_PATH=$(sf_metadata_path 1 1 "effort" "plan" "md")
echo "# Plan" > "$PLAN_PATH"
```

### ❌ VIOLATION 4: Metadata Outside .software-factory
```bash
# WRONG - Creating in project root
echo "# Work Log" > work-log.md

# RIGHT - In .software-factory with timestamp
LOG_PATH=$(sf_metadata_path 1 1 "effort" "work-log" "log")
echo "# Work Log" > "$LOG_PATH"
```

## Enforcement

- **Penalty**: -100% for ANY violation (IMMEDIATE FAILURE)
- **Detection**: Automated scanning for non-timestamped files
- **Prevention**: Helper function prevents violations
- **Recovery**: Migration function fixes old patterns

## Related Rules

- **R343**: Metadata directory standardization (parent rule)
- **R344**: Metadata location tracking
- **R054**: Implementation plan creation
- **R264**: Work log tracking requirements
- **R303**: Phase/Wave document location protocol

## Key Principle

**"Every metadata file is unique. No merge conflicts. Ever."**

This ensures:
- Zero merge conflicts from metadata
- Perfect parallel agent operation
- Complete audit trail with timestamps
- Easy recovery and debugging
- Clean integration every time

---

**REMEMBER**: If a file doesn't have a timestamp, it's WRONG. No exceptions. Use the helper function ALWAYS.