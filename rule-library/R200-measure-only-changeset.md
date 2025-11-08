# Rule R200: Measure ONLY Effort Changeset - NOT Base Branch Files

## Rule Statement
Code reviewers and SW engineers MUST measure ONLY the files changed/added by the current effort implementation. Files from the base branch or elsewhere in the project that were NOT modified by this effort MUST NOT be included in line count measurements. Measuring base branch files is a CRITICAL ERROR requiring immediate stop.

## Criticality Level
**BLOCKING** - Measuring wrong files leads to incorrect splits and implementation failures

## Enforcement Mechanism
- **Technical**: Use `git diff` to identify changeset before measuring
- **Behavioral**: STOP IMMEDIATELY if measuring files outside changeset
- **Grading**: -50% for measuring base branch files (Major architectural failure)

## Core Principle

```
MEASURE ONLY WHAT YOU CHANGED
Base Branch Files = NOT YOUR CONCERN
Your Effort's Changes = Your ONLY Focus
```

## Detailed Requirements

### CODE REVIEWER: Correct Measurement Protocol

```bash
# ❌❌❌ WRONG - Measuring entire repository
measure_wrong() {
    cd efforts/phase1/wave1/api-types
    ./tools/line-counter.sh  # WRONG! This counts EVERYTHING
    # Result: 50,000 lines (includes entire base project!)
}

# ❌❌❌ WRONG - Including unmodified files
measure_wrong_2() {
    find . -name "*.go" | xargs wc -l  # WRONG! Counts all Go files
    # Includes thousands of lines from base branch
}

# ✅✅✅ CORRECT - Measure ONLY changeset
measure_correct() {
    local effort_name="api-types"
    local branch_name="phase1/wave1/api-types"
    
    echo "═══════════════════════════════════════════════════════"
    echo "MEASURING EFFORT CHANGESET ONLY"
    echo "Effort: $effort_name"
    echo "Branch: $branch_name"
    echo "═══════════════════════════════════════════════════════"
    
    # STEP 1: Identify what changed
    echo "Step 1: Identifying changeset..."
    git diff --name-only main..HEAD > changed_files.txt
    
    # STEP 2: Show what we're measuring
    echo "Files changed in this effort:"
    cat changed_files.txt
    
    # STEP 3: Measure ONLY changed files
    echo "Step 2: Measuring ONLY changed files..."
    local total_lines=0
    while read -r file; do
        if [ -f "$file" ]; then
            lines=$(wc -l < "$file")
            echo "  $file: $lines lines"
            total_lines=$((total_lines + lines))
        fi
    done < changed_files.txt
    
    echo "═══════════════════════════════════════════════════════"
    echo "TOTAL CHANGESET SIZE: $total_lines lines"
    echo "═══════════════════════════════════════════════════════"
    
    # STEP 4: Verify this is reasonable
    if [ $total_lines -gt 10000 ]; then
        echo "⚠️ WARNING: Changeset seems unusually large!"
        echo "Did you accidentally include base branch files?"
        echo "Please verify the changeset is correct!"
    fi
}
```

### SW ENGINEER: Implementation Measurement

```bash
# When implementing, track what you're adding/changing
track_implementation() {
    echo "Files I'm creating/modifying in this effort:"
    
    # Track new files
    git ls-files --others --exclude-standard > new_files.txt
    
    # Track modified files  
    git diff --name-only > modified_files.txt
    
    # Combine for measurement
    cat new_files.txt modified_files.txt | sort -u > my_changeset.txt
    
    echo "My changeset for this effort:"
    cat my_changeset.txt
    
    # Measure ONLY these files
    local total=0
    while read -r file; do
        if [ -f "$file" ]; then
            lines=$(wc -l < "$file")
            total=$((total + lines))
        fi
    done < my_changeset.txt
    
    echo "Total lines in MY changes: $total"
}
```

### CRITICAL: Stop Work Protocol

```bash
# If you find yourself measuring wrong files
detect_measurement_error() {
    # Check if measuring base branch files
    if git diff --name-only main..HEAD | grep -q "^$"; then
        echo "❌❌❌ CRITICAL ERROR DETECTED ❌❌❌"
        echo "NO CHANGES FOUND - AM I MEASURING BASE BRANCH?"
        echo "STOPPING IMMEDIATELY!"
        exit 1
    fi
    
    # Check if changeset is suspiciously large
    CHANGESET_SIZE=$(git diff --name-only main..HEAD | xargs wc -l 2>/dev/null | tail -1 | awk '{print $1}')
    if [ "$CHANGESET_SIZE" -gt 5000 ]; then
        echo "❌❌❌ SUSPICIOUS SIZE DETECTED ❌❌❌"
        echo "Changeset is $CHANGESET_SIZE lines"
        echo "This seems too large for a single effort!"
        echo "Am I including files I shouldn't?"
        echo "STOPPING FOR VERIFICATION!"
        exit 1
    fi
}
```

## What to Measure vs What to Ignore

### ✅ MEASURE These Files:
- Files YOU created in this effort
- Files YOU modified in this effort  
- New test files YOU added
- New documentation YOU wrote
- Generated files that YOUR changes created

### ❌ NEVER MEASURE These Files:
- Files that existed in base branch unchanged
- Dependencies from vendor/node_modules
- Files from other efforts
- System files or IDE configs
- Files you're just reading/referencing

## Git Commands for Correct Measurement

```bash
# Show only files changed in this effort
git diff --name-only main..HEAD

# Show lines added/removed per file
git diff --stat main..HEAD

# Show detailed diff with line counts
git diff main..HEAD | diffstat

# Count lines in changed files only
git diff --name-only main..HEAD | xargs wc -l

# Show what's new (not in base branch)
git diff --diff-filter=A --name-only main..HEAD
```

## Common Mistakes to Avoid

### ❌ Using line-counter.sh on entire directory
```bash
cd efforts/phase1/wave1/my-effort
./tools/line-counter.sh  # WRONG if repo has other files!
```

### ❌ Measuring before implementation
```bash
# Clone repo
git clone --sparse $REPO efforts/phase1/wave1/my-effort
cd efforts/phase1/wave1/my-effort
./tools/line-counter.sh  # WRONG! No changes yet!
```

### ❌ Including test data or fixtures
```bash
# If you added 100 lines of code and 10,000 lines of test data
# Only the 100 lines of CODE count toward limit
```

## Correct Workflow

```
1. SW Engineer implements effort
2. SW Engineer commits changes
3. Code Reviewer measures ONLY the changeset
4. If changeset >800 lines, plan splits
5. Each split measures ONLY its subset of changes
6. NEVER include unchanged base files
```

## Emergency Stop Triggers

STOP IMMEDIATELY if:
1. Measurement includes files you didn't change
2. Measurement >5000 lines (suspicious)
3. Measurement includes vendor/dependencies  
4. Measurement includes entire directories unchanged
5. You're measuring before implementation
6. Git diff shows no changes but measurement shows lines

## Integration with Other Rules

- **R198**: Line counter usage (use correctly on changeset)
- **R199**: Single reviewer for splits (of changeset only)
- **R007**: Size limits (apply to changeset only)
- **R196**: Base branch selection (measure changes FROM base)

## Grading Impact

- **Measuring base branch files**: -50% (Critical failure)
- **Including unchanged files in splits**: -30% (Major error)
- **Proceeding after detecting error**: -40% (Ignored stop signal)
- **Creating splits of base files**: -45% (Architectural failure)

## Summary

**Remember**: 
- You're measuring YOUR WORK, not the entire project
- If you didn't change it, don't count it
- Base branch files are NOT your concern
- STOP if you detect measurement errors
- Use git diff to identify YOUR changeset
- Measure ONLY what YOU changed