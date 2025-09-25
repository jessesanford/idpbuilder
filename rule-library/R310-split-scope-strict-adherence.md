# 🚨🚨🚨 BLOCKING: RULE R310 - Split Scope Strict Adherence Protocol 🚨🚨🚨

**Category:** Implementation Control  
**Agents:** SW Engineer (PRIMARY), Code Reviewer  
**Criticality:** BLOCKING - Over-engineering causes 3-5X size violations  
**Penalty:** -100% for exceeding split scope, -50% for unrequested features

## 🔴🔴🔴 SUPREME DIRECTIVE: IMPLEMENT EXACTLY WHAT'S SPECIFIED 🔴🔴🔴

**NO MORE, NO LESS - EVEN IF IT SEEMS INCOMPLETE!**

## The Critical Problem This Solves

**CATASTROPHIC VIOLATIONS FROM PRODUCTION:**

### 🚨 THE 2667% VIOLATION - WORST CASE EVER RECORDED 🚨
```
Split Plan Said: "Implement 3 files in pkg/builder" (~572 lines)
SW Engineer Did: Implemented 80 files across 26 packages (ENTIRE CODEBASE)
Result: 2667% scope violation - Complete project failure
Root Cause: Copy-pasted entire repo without reading plan
```

### 🚨 RECENT 3.4X VIOLATION FROM TRANSCRIPT:
```
Split-003 Plan Said: "Implement selected command files and utilities" (~650 lines)
SW Engineer Interpreted:
  - "selected" = all the important commands (not just listed ones)
  - "utilities" = complete utility module (not minimal helpers)
  - "handlers" = comprehensive handler system
  - No test count = add full test coverage (430 lines!)
Result: 2,215 lines instead of 650 (3.4X OVERRUN)
Root Cause: "Complete the feature" mindset vs "stay within budget"
Consequence: COMPLETE FAILURE - Split rejected, effort failed
```

### Previous Major Violation:
```
Split Plan Said: "Implement functional options pattern"
SW Engineer Did: Created 47 option functions, validation, Clone(), extensive tests
Result: 1,512 lines instead of 400 (3.8X OVERRUN)
Consequence: COMPLETE FAILURE - Split rejected, effort failed
```

## 🚨 MANDATORY ADHERENCE PROTOCOL

### 0. COMPLETE PRE-IMPLEMENTATION CHECKLIST (NEW - REQUIRED)

**BEFORE WRITING ANY CODE:**
```bash
# MANDATORY: Complete the pre-implementation checklist
cp $CLAUDE_PROJECT_DIR/templates/SW-ENGINEER-PRE-IMPLEMENTATION-CHECKLIST.md ./checklist-$(date +%s).md
echo "📋 Completing mandatory pre-implementation checklist..."
# Fill out the checklist completely
# This is NOT optional - it prevents 2667% violations
```

### 1. READ THE "DO NOT IMPLEMENT" SECTION FIRST

**Before writing ANY code:**
```bash
# MANDATORY: Extract and display what NOT to do
echo "═══════════════════════════════════════════════════════"
echo "🛑 DO NOT IMPLEMENT LIST FROM SPLIT PLAN:"
echo "═══════════════════════════════════════════════════════"
grep -A 20 "DO NOT IMPLEMENT\|STOP BOUNDARIES\|FORBIDDEN" SPLIT-PLAN-*.md

# Acknowledge the boundaries
echo "✅ I will NOT implement anything in the DO NOT list"
echo "✅ I will STOP at the specified boundaries"
echo "✅ I will RESIST adding 'helpful' extras"
```

### 2. COUNT BEFORE CODING - ENHANCED WITH FILE VALIDATION

**Extract exact counts from plan:**
```bash
# CRITICAL: Count FILES first (new requirement after 2667% violation)
FILE_COUNT=$(grep -c "^###.*File:" SPLIT-PLAN-*.md)
if [ -z "$FILE_COUNT" ] || [ "$FILE_COUNT" -eq 0 ]; then
    echo "🚨 FATAL: No file count found in split plan!"
    echo "Split plans MUST specify exact files to implement"
    exit 1
fi

# Set HARD limits based on file count
MAX_FILES=$((FILE_COUNT * 2))  # 200% is absolute maximum
STOP_FILES=$((FILE_COUNT + 2))  # Stop if 2 files over

echo "🔴🔴🔴 FILE COUNT ENFORCEMENT 🔴🔴🔴"
echo "  Files to implement: $FILE_COUNT (EXACTLY)"
echo "  Warning at: $STOP_FILES files"
echo "  AUTOMATIC FAILURE at: $MAX_FILES files"
echo "  Current violation record: 2667% (80 files instead of 3)"

# Parse the split plan for exact numbers
FUNCTION_COUNT=$(grep -i "EXACTLY.*functions\|IMPLEMENT.*functions" SPLIT-PLAN-*.md | grep -o '[0-9]+' | head -1)
METHOD_COUNT=$(grep -i "EXACTLY.*methods\|IMPLEMENT.*methods" SPLIT-PLAN-*.md | grep -o '[0-9]+' | head -1)
TEST_COUNT=$(grep -i "EXACTLY.*tests\|WRITE.*tests" SPLIT-PLAN-*.md | grep -o '[0-9]+' | head -1)

echo "📊 SCOPE LIMITS FOR THIS SPLIT:"
echo "  Files to implement: ${FILE_COUNT} (HARD LIMIT)"
echo "  Functions to implement: ${FUNCTION_COUNT:-0}"
echo "  Methods to implement: ${METHOD_COUNT:-0}"
echo "  Tests to write: ${TEST_COUNT:-0}"
echo "  ANYTHING ELSE: FORBIDDEN"

# Create a marker file to track start of implementation
touch .implementation-start-marker
```

### 3. IMPLEMENT ONLY WHAT'S EXPLICITLY NAMED

#### ❌ WRONG - Adding Unlisted Features:
```go
// Split plan says: "Implement WithImage, WithContext, WithPlatform"
// SW Engineer writes:
func WithImage(image string) Option { ... }      // ✅ Correct
func WithContext(ctx context.Context) Option { ... } // ✅ Correct  
func WithPlatform(platform string) Option { ... }    // ✅ Correct
func WithTimeout(timeout time.Duration) Option { ... } // ❌ NOT IN PLAN!
func WithRetry(retries int) Option { ... }           // ❌ NOT IN PLAN!
func WithDebug(debug bool) Option { ... }            // ❌ NOT IN PLAN!
// Result: 6 functions instead of 3 = VIOLATION
```

#### ✅ RIGHT - Exactly What's Listed:
```go
// Split plan says: "Implement WithImage, WithContext, WithPlatform"
// SW Engineer writes:
func WithImage(image string) Option { ... }      // ✅ Listed
func WithContext(ctx context.Context) Option { ... } // ✅ Listed
func WithPlatform(platform string) Option { ... }    // ✅ Listed
// STOP HERE - Even if it seems incomplete!
```

### 4. VALIDATION CHECKPOINTS - ENHANCED WITH FILE MONITORING

#### Every 100 Lines (NEW - MANDATORY):
```bash
# MANDATORY: Check file count every 100 lines to prevent 2667% violations
check_file_count() {
    PLANNED_FILES=$(grep -c "^###.*File:" SPLIT-PLAN-*.md)
    ACTUAL_FILES=$(find . -type f -name "*.go" -newer .implementation-start-marker 2>/dev/null | wc -l)
    
    echo "📊 FILE COUNT CHECK (Every 100 lines):"
    echo "  Planned: $PLANNED_FILES files"
    echo "  Current: $ACTUAL_FILES files"
    echo "  Percentage: $(($ACTUAL_FILES * 100 / $PLANNED_FILES))%"
    
    if [[ $ACTUAL_FILES -gt $((PLANNED_FILES + 2)) ]]; then
        echo "⚠️ WARNING: Exceeding file count by more than 2!"
        echo "STOP and review your implementation!"
    fi
    
    if [[ $ACTUAL_FILES -gt $((PLANNED_FILES * 2)) ]]; then
        echo "🚨🚨🚨 CATASTROPHIC FAILURE IMMINENT!"
        echo "You are at $(($ACTUAL_FILES * 100 / $PLANNED_FILES))% of planned files!"
        echo "This is approaching the 2667% violation record!"
        echo "STOP IMMEDIATELY!"
        exit 1
    fi
}

# Run this check after adding every ~100 lines
LINES_SINCE_CHECK=$(git diff --stat | tail -1 | awk '{print $4}')
if [ "$LINES_SINCE_CHECK" -gt 100 ]; then
    check_file_count
fi
```

#### After Each File:
```bash
# Count what you've actually implemented
ACTUAL_FILES=$(find . -type f -name "*.go" -newer .implementation-start-marker 2>/dev/null | wc -l)
ACTUAL_FUNCTIONS=$(grep -c "^func [A-Z]" *.go 2>/dev/null || echo 0)
ACTUAL_METHODS=$(grep -c "^func (.*) " *.go 2>/dev/null || echo 0)
ACTUAL_TESTS=$(grep -c "^func Test" *_test.go 2>/dev/null || echo 0)

echo "📈 CURRENT IMPLEMENTATION COUNT:"
echo "  Files: $ACTUAL_FILES / $FILE_COUNT 🔴 CRITICAL METRIC"
echo "  Functions: $ACTUAL_FUNCTIONS / $FUNCTION_COUNT"
echo "  Methods: $ACTUAL_METHODS / $METHOD_COUNT"  
echo "  Tests: $ACTUAL_TESTS / $TEST_COUNT"

# File violation check (HIGHEST PRIORITY)
if [ "$ACTUAL_FILES" -gt "$FILE_COUNT" ]; then
    echo "❌ VIOLATION: Too many files! You have $ACTUAL_FILES but plan specifies $FILE_COUNT!"
    echo "This is a $(($ACTUAL_FILES * 100 / $FILE_COUNT))% violation!"
    exit 1
fi

# Function violation check
if [ "$ACTUAL_FUNCTIONS" -gt "$FUNCTION_COUNT" ]; then
    echo "❌ VIOLATION: Too many functions! Remove extras!"
    exit 1
fi
```

#### Before Committing:
```bash
# Final scope validation
validate_split_scope() {
    echo "🔍 FINAL SCOPE VALIDATION"
    
    # Check for forbidden patterns
    if grep -q "func.*Clone\|func.*Copy" *.go; then
        echo "❌ FOUND Clone/Copy methods - were these requested?"
        [ "$ALLOW_CLONE" != "true" ] && exit 1
    fi
    
    if grep -q "func.*Validate\|func.*Valid" *.go; then
        echo "❌ FOUND validation methods - were these requested?"
        [ "$ALLOW_VALIDATION" != "true" ] && exit 1
    fi
    
    if grep -q "// TODO\|// FIXME" *.go; then
        echo "⚠️ Found TODOs - don't add work for later!"
    fi
    
    echo "✅ Scope validation passed"
}
```

## 🛑 WHEN THE PLAN SEEMS INSUFFICIENT

### DON'T ASSUME - ASK!

**If the split plan seems incomplete:**

```bash
# DO NOT DO THIS:
"The plan only lists 3 option functions but we need 15 for a complete implementation.
I'll add the other 12 to make it complete." # ❌ WRONG - 5X OVERRUN!

# DO THIS INSTEAD:
echo "📋 SCOPE CLARIFICATION NEEDED:"
echo "The split plan specifies 3 option functions: WithImage, WithContext, WithPlatform"
echo "This seems insufficient for a complete options pattern."
echo "Should I:"
echo "  A) Implement ONLY these 3 as specified (recommended)"
echo "  B) Wait for clarification on additional functions"
echo "Proceeding with Option A - implementing ONLY the 3 specified functions."
```

## 📏 SIZE CALCULATION REALITY CHECK

### Realistic Line Counts (Go Example):

```go
// Simple option function: 10-20 lines
func WithImage(image string) Option {
    return func(c *Config) error {
        if image == "" {
            return fmt.Errorf("image cannot be empty")
        }
        c.Image = image
        return nil
    }
}

// With validation and docs: 30-50 lines
// With complex logic: 50-100 lines
// Full CRUD operations: 200-400 lines per resource
```

### Common Over-Engineering Patterns:

| Plan Says | SW Does | Result |
|-----------|---------|---------|
| "Basic struct" | Adds 15 methods | 5X overrun |
| "3 functions" | Implements 47 | 15X overrun |
| "Simple tests" | Full coverage | 10X overrun |
| "Core types" | Adds validation | 3X overrun |

## 🔴 ENFORCEMENT MECHANISMS

### Pre-Implementation Gate:
```bash
# Must acknowledge scope before starting
ACKNOWLEDGED=$(grep -c "I will NOT implement" .scope-ack)
[ "$ACKNOWLEDGED" -eq 0 ] && {
    echo "❌ Must acknowledge scope boundaries first!"
    exit 1
}
```

### Mid-Implementation Gate:
```bash
# Check every 100 lines added
LINES_ADDED=$(git diff --stat | tail -1 | awk '{print $4}')
if [ "$LINES_ADDED" -gt 100 ]; then
    validate_split_scope || exit 1
fi
```

### Pre-Commit Gate:
```bash
# Final validation before commit
EXCESS_FUNCTIONS=$(( $(grep -c "^func " *.go) - FUNCTION_COUNT ))
if [ "$EXCESS_FUNCTIONS" -gt 0 ]; then
    echo "❌ BLOCKING: $EXCESS_FUNCTIONS extra functions detected!"
    echo "Remove excess functions or split will be rejected!"
    exit 1
fi
```

## 🎯 Success Criteria

A split implementation is successful when:

1. ✅ Implements EXACTLY the named functions/methods
2. ✅ Stays within specified line count
3. ✅ Contains NO unrequested features
4. ✅ Includes NO "helpful" extras
5. ✅ Follows all "DO NOT" instructions
6. ✅ Stops at specified boundaries

## 📊 Grading Impact

- **FILE COUNT VIOLATIONS (MOST SEVERE):**
  - Exceeding by >200%: **-100% (AUTOMATIC FAILURE - NO APPEAL)**
  - Exceeding by >150%: **-100% (AUTOMATIC FAILURE)**
  - Exceeding by >100%: **-80%**
  - Exceeding by >50%: **-60%**
  
- **LINE COUNT VIOLATIONS:**
  - 3.4X overrun (like Split-003): **-100% (AUTOMATIC FAILURE)**
  - 2X-3X overrun: **-80%**
  - 1.5X-2X overrun: **-60%**
  
- **OTHER VIOLATIONS:**
  - Interpreting "selected" as "all": **-75%**
  - Adding unrequested features: **-50% per feature**
  - Ignoring "DO NOT" instructions: **-75%**
  - "Completing" partial implementations: **-50%**
  - Adding tests when not specified: **-40%**
  - Refactoring without request: **-25%**
  - Not completing pre-implementation checklist: **-40%**

## 💡 Mental Model

**Think of splits like LEGO instructions:**
- Step says "Add 3 red blocks" → Add EXACTLY 3 red blocks
- Don't add blue blocks because they'd look nice
- Don't add 5 red blocks to make it "stronger"
- Don't skip ahead to step 10
- Don't go back and "improve" step 1

## 🔑 Key Phrases to Internalize

- "If it's not listed, don't implement it"
- "Incomplete is better than oversized"
- "The plan is the contract"
- "Stop means stop"
- "Ask, don't assume"
- "Selected means ONLY the ones explicitly named"
- "Stay within budget, don't complete the feature"
- "Partial implementation is the goal"

## Examples of Proper Adherence

### Example 1: Options Implementation
```
Plan: Implement WithImage, WithContext, WithPlatform (3 functions)
Did: Implemented exactly those 3 functions
Size: 145 lines (well under 400 limit)
Result: SUCCESS ✅
```

### Example 2: Struct Definition
```
Plan: Define Config struct with 5 fields, no methods
Did: Created struct with 5 fields, added no methods
Size: 35 lines
Result: SUCCESS ✅
```

### Example 3: Test Implementation
```
Plan: Write 3 basic tests, one per function
Did: Wrote exactly 3 tests, no edge cases
Size: 90 lines
Result: SUCCESS ✅
```

## Summary

**The #1 cause of split failures is ambiguous language being interpreted as "make it complete".**

**RECENT FAILURES:**
- **3.4X OVERRUN**: "selected commands" → ALL commands (2,215 lines vs 650)
- **2667% VIOLATION**: 80 files instead of 3 (complete copy-paste of repo)

This rule ensures:
- SW Engineers implement EXACTLY what's specified (FILE AND LINE COUNT CRITICAL)
- No interpretation of vague words like "selected", "utilities", "handlers"
- Mandatory pre-implementation checklist completion
- Incremental monitoring every 100 lines
- Automatic failure for major overruns
- "Stay within budget" mindset, not "complete the feature"
- Clear communication when scope seems insufficient

**Remember: The engineers aren't trying to fail - they think they're being helpful by completing features. The solution is EXPLICIT, UNAMBIGUOUS plans that make it impossible to misinterpret.**

**NEVER FORGET: "Selected" doesn't mean "all important ones" - it means ONLY the ones explicitly listed by name!**