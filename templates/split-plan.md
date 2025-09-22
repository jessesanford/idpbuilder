# SPLIT PLAN TEMPLATE - PREVENTING OVER-ENGINEERING

<!-- STORAGE LOCATION: This split plan should be saved in:
     .software-factory/phase[X]/wave[Y]/[effort-name]-split-[XXX]/SPLIT-PLAN-YYYYMMDD-HHMMSS.md
     within the too-large effort's working directory. This keeps split plans organized. -->

## 🔴🔴🔴 CRITICAL: ACTUAL 3.4X OVERRUN FROM TRANSCRIPT 🔴🔴🔴

**ACTUAL FAILURE: Split-003 exceeded plan by 3.4x (2,215 lines vs 650 planned)**

Why it happened:
- Plan said "selected command files" → Engineer built ALL commands
- Plan said "utilities" → Engineer built complete modules
- Plan said "handlers" → Engineer built comprehensive system
- Plan didn't mention tests → Engineer added 430 lines of tests
- Engineer had "complete the feature" mindset instead of "stay within budget"

## 🔴🔴🔴 CRITICAL: FILE COUNT ENFORCEMENT (R314) 🔴🔴🔴

**AFTER THE 2667% VIOLATION (80 FILES INSTEAD OF 3), FILE COUNT IS NOW MANDATORY**

### 📊 SCOPE METRICS (MANDATORY - FILL ALL WITH EXACT VALUES)
- **EXACT FILES TO IMPLEMENT: [NUMBER]** ← 🔴 LIST EACH FILE BY NAME BELOW
- **MAXIMUM ALLOWED FILES: [EXACT NUMBER x 1.5]** ← AUTOMATIC FAILURE IF EXCEEDED  
- **Functions to implement: [EXACT COUNT]** ← LIST EACH FUNCTION BY NAME BELOW
- **Methods to implement: [EXACT COUNT]** ← LIST EACH METHOD BY NAME BELOW
- **Tests to write: [EXACT COUNT OR ZERO]** ← BE EXPLICIT: "0" or exact number
- **Estimated total lines: [NUMBER]** (must be <600 for safety margin)

## 🔴🔴🔴 CRITICAL: EXPLICIT SCOPE BOUNDARIES 🔴🔴🔴

**THIS TEMPLATE PREVENTS 2667% VIOLATIONS BY REQUIRING EXPLICIT BOUNDARIES**

## Split Metadata
- **Split Number**: XXX
- **Parent Effort**: [effort-name]
- **Base Branch**: [MANDATORY per R337 - MUST be from orchestrator-state.json!]
  - For split-001: Use the same base as the oversized effort (from state file)
  - For split-N: Use split-(N-1) as the base (must be recorded in state file)
- **Base Branch Tracking** (R337 MANDATORY):
  - **Planned Base**: [What this plan specifies]
  - **Actual Base**: [What orchestrator records in state - MUST match planned]
  - **Base Commit**: [SHA from orchestrator when creating]
  - **Dependent Splits**: [Which splits depend on this one]
  - **Cascade Group**: [All splits that would need rebase if this changes]
- **Original Size**: [total lines that exceeded 800]
- **This Split Target**: [target lines, max 400-600 to leave buffer]
- **Created**: [timestamp]

## 🔴🔴🔴 CRITICAL: FILE PLACEMENT (R326) - NO SPLIT SUBDIRECTORIES! 🔴🔴🔴

**⚠️⚠️⚠️ CATASTROPHIC BUG PREVENTION ⚠️⚠️⚠️**

### ❌ WRONG - CAUSES MASSIVE MEASUREMENT ERRORS:
```
Working in: efforts/phase1/wave1/registry-SPLIT-001/
Creating: split-001/pkg/registry/auth.go  ❌❌❌ NEVER DO THIS!
```

### ✅ CORRECT - STANDARD PROJECT STRUCTURE:
```
Working in: efforts/phase1/wave1/registry-SPLIT-001/
Creating: pkg/registry/auth.go  ✅✅✅ FILES GO DIRECTLY HERE!
```

**CRITICAL**: Files MUST go in standard directories (pkg/, cmd/, tests/).
**NEVER** create a subdirectory matching the split name (split-001/, split-002/, etc.)!
**PENALTY**: Creating split subdirectories = -100% IMMEDIATE FAILURE

## 🚨 EXPLICIT SCOPE DEFINITION (MANDATORY)

### 📋 MINIMUM VIABLE IMPLEMENTATION (What MUST be done)
**This is your CONTRACT - implement EXACTLY these items:**

#### Files to Create/Modify (EXACT LIST - NO ADDITIONS)
```
1. pkg/config/options.go (CREATE) - 150 lines MAX
2. pkg/config/config.go (MODIFY) - Add 50 lines MAX
3. pkg/config/options_test.go (CREATE) - 100 lines MAX
TOTAL FILES: 3 (implementing any other file = VIOLATION)
```

#### Functions to Implement (EXACT LIST - NO ADDITIONS)
**IMPLEMENT THESE [N] FUNCTIONS ONLY - NO MORE:**
```go
// List EVERY function by exact signature:
1. func WithImage(image string) Option         // 10-15 lines
2. func WithContext(ctx context.Context) Option // 10-15 lines  
3. func WithPlatform(platform string) Option   // 10-15 lines
// TOTAL: 3 functions (30-45 lines)
// STOP HERE - DO NOT ADD WithTimeout, WithRetry, etc.
```

#### Methods to Implement (EXACT LIST OR "NONE")
```go
// If NO methods: write "NONE - DO NOT ADD ANY METHODS"
// If methods needed, list exact signatures:
NONE - DO NOT ADD ANY METHODS IN THIS SPLIT
```

#### Structs/Types to Define (WITH FIELD COUNTS)
**CREATE EXACTLY [N] TYPES WITH THESE FIELDS ONLY:**
```go
// Example - BE THIS EXPLICIT:
type Config struct {  // EXACTLY 5 fields, NO methods yet
    Image    string
    Context  context.Context
    Platform string
    Debug    bool
    Timeout  time.Duration
}
// NO METHODS ON THIS STRUCT IN THIS SPLIT
```

#### Files to Create/Modify (EXACT COUNT - R314 ENFORCEMENT)

**🔴 IMPLEMENTING MORE FILES THAN LISTED = AUTOMATIC FAILURE 🔴**

### File: pkg/config/options.go (CREATE - 150 lines MAX)
  - ONLY the 3 option functions listed above
  - Basic Option type definition
  - NO validation, NO defaults, NO extra helpers

### File: pkg/config/config.go (CREATE - 200 lines MAX)
  - ONLY the Config struct
  - ONLY NewConfig() constructor
  - NO methods, NO validation, NO converters
```

### 📏 MAXIMUM ALLOWED SCOPE (Hard limits)
**Exceed these = AUTOMATIC FAILURE:**
- **Total lines:** 650 HARD MAXIMUM (target 400-500)
- **Total files:** [EXACT COUNT x 1.5] files MAXIMUM
- **Per file:** As specified above (DO NOT EXCEED)
- **Tests:** [0 or EXACT NUMBER] - if 0, write NO TESTS AT ALL

### 🛑 DO NOT IMPLEMENT LIST (FORBIDDEN ITEMS)

#### ⚠️⚠️⚠️ WARNING: THE 3.4X OVERRUN (FROM TRANSCRIPT) ⚠️⚠️⚠️
**ACTUAL FAILURE: Engineer saw "selected commands" and built ALL commands**
**Result: 2,215 lines instead of 650 = COMPLETE FAILURE**

**EXPLICITLY FORBIDDEN IN THIS SPLIT:**
- ❌ DO NOT implement ANY files not listed in "Files to Create/Modify" above
- ❌ DO NOT interpret "selected" as "all" or "complete"
- ❌ DO NOT interpret "utilities" as "comprehensive utility module"
- ❌ DO NOT add tests unless explicitly listed with exact count
- ❌ DO NOT add validation methods
- ❌ DO NOT add Clone() or Copy() methods
- ❌ DO NOT add conversion/transformation methods
- ❌ DO NOT implement error handling beyond basics
- ❌ DO NOT add comprehensive tests (basic tests only)
- ❌ DO NOT refactor existing code
- ❌ DO NOT add documentation beyond minimal comments
- ❌ DO NOT implement features "while you're there"
- ❌ DO NOT add helper/utility functions not listed
- ❌ DO NOT "make it complete" or "production ready"
- ❌ DO NOT exceed file count by >2x (automatic failure)

### 📝 Test Requirements (BE EXPLICIT ABOUT ZERO OR COUNT)

#### Option A: NO TESTS IN THIS SPLIT
```
TESTS: ZERO - DO NOT WRITE ANY TEST FILES
Reason: Tests allocated to Split-004
```

#### Option B: MINIMAL TESTS ONLY
```
File: options_test.go - 100 lines MAXIMUM
EXACT TEST COUNT: 3 tests
  1. TestWithImage() - 30 lines max
  2. TestWithContext() - 30 lines max
  3. TestWithPlatform() - 30 lines max
  
FORBIDDEN:
  - NO edge case testing
  - NO error condition testing (unless listed)
  - NO benchmarks
  - NO example functions
  - NO table-driven tests (unless specified)
```

## 📊 SIZE CALCULATION (REALISTIC)

### How We Calculate (Go Example):
```
Function with validation: 30-50 lines average
Simple option function:   10-20 lines
Struct definition:        5-10 lines per 5 fields
Basic test:              20-30 lines
Comprehensive test:       50-100 lines

This Split:
- 3 option functions × 15 lines = 45 lines
- 1 Option type = 10 lines  
- 1 Config struct = 25 lines
- 1 NewConfig function = 30 lines
- 3 basic tests × 30 lines = 90 lines
TOTAL: ~200 lines (well under 400 limit)
```

## 🚨🚨🚨 R355 PRODUCTION READY CODE (SUPREME LAW) 🚨🚨🚨

**ALL CODE MUST BE PRODUCTION READY - ZERO TOLERANCE FOR VIOLATIONS**

### ❌ ABSOLUTELY FORBIDDEN (AUTOMATIC -100% FAILURE):
- NO stubs or placeholder implementations
- NO mocks except in test directories
- NO hardcoded credentials or secrets
- NO static configuration values
- NO TODO/FIXME markers in code
- NO returning nil for "later implementation"
- NO panic("not implemented")

### ✅ REQUIRED PATTERNS:
```go
// ❌ WRONG - Hardcoded values (AUTOMATIC FAILURE)
dbURL := "postgres://localhost/mydb"
apiKey := "sk-12345"
maxRetries := 3

// ✅ CORRECT - Configuration-driven
dbURL := os.Getenv("DATABASE_URL")
if dbURL == "" {
    dbURL = defaultDBURL // Default only, not hardcoded
}
apiKey := os.Getenv("API_KEY")
if apiKey == "" {
    return errors.New("API_KEY required")
}
maxRetries := config.GetInt("max_retries", 3)
```

### Implementation Example:
```go
// ❌ WRONG - Stub implementation (AUTOMATIC FAILURE)
func ProcessData(data []byte) error {
    // TODO: implement this
    return nil
}

// ✅ CORRECT - Complete minimal implementation
func ProcessData(data []byte) error {
    if len(data) == 0 {
        return errors.New("data cannot be empty")
    }
    // Minimal but complete implementation
    processed := bytes.ToUpper(data)
    return saveToMemory(processed)
}
```

## 🔴 MANDATORY ADHERENCE CHECKPOINTS

### Before Starting:
```bash
echo "SPLIT SCOPE CHECKPOINT:"
echo "Functions to implement: EXACTLY 3"
echo "Methods to implement: EXACTLY 0"
echo "Validation to add: NONE"
echo "Clone/Copy methods: NONE"
echo "Extra features: NONE"
echo "R355 Compliance: NO STUBS, NO HARDCODING"
```

### During Implementation:
```bash
# After each file, check scope
CURRENT_LINES=$(git diff --stat | tail -1 | awk '{print $4}')
echo "Current additions: $CURRENT_LINES lines"
echo "Remaining in scope: $(( 400 - CURRENT_LINES )) lines"

# If approaching limit with work remaining
if [ $CURRENT_LINES -gt 350 ]; then
    echo "⚠️ APPROACHING LIMIT - STOP ADDING FEATURES"
    echo "Complete ONLY what's explicitly listed"
fi
```

### After Implementation:
```bash
# Count what was actually implemented
FUNC_COUNT=$(grep -c "^func " *.go)
METHOD_COUNT=$(grep -c "^func (.*) " *.go)
TEST_COUNT=$(grep -c "^func Test" *_test.go)

echo "Functions implemented: $FUNC_COUNT (should be EXACTLY as planned)"
echo "Methods added: $METHOD_COUNT (should be 0 if not planned)"
echo "Tests written: $TEST_COUNT (should match plan)"
```

## 💡 ACTUAL EXAMPLES FROM FAILURES

### ❌ ACTUAL BAD EXAMPLE (Led to 3.4x overrun):
```markdown
## Split-003: CLI and Utilities
- Implement selected command files
- Add necessary utilities and handlers
- Create workflow orchestration components
```
Result: Engineer implemented ALL commands, ALL utilities, complete system = 2,215 lines

### ✅ CORRECTED GOOD EXAMPLE (Would prevent overrun):
```markdown
## Split-003: CLI Push Command Only (650 lines MAX)

FILES (EXACTLY 3):
1. cmd/push.go - Implement ONLY HandlePush() function (200 lines)
2. pkg/push/validator.go - Implement ONLY ValidatePushArgs() (150 lines)  
3. cmd/push_test.go - Write EXACTLY 2 tests (100 lines)

DO NOT IMPLEMENT:
- Other command files (build.go, tag.go, etc.)
- Workflow orchestration (separate split)
- Utility modules beyond validator
- Integration tests
```

### ❌ BAD (NO BOUNDARIES):
```markdown
## Scope
Complete implementation of options.go file
```

### ✅ GOOD (CLEAR BOUNDARIES):
```markdown
## Scope
Implement ONLY:
- Lines 1-150 of options.go (3 functions)
- NO validation logic (separate split)
- NO comprehensive tests (separate split)
STOP AT: 150 lines even if file incomplete
```

## 🚨 CONSEQUENCES OF IGNORING BOUNDARIES

### What Actually Happened (From Transcript):
```
Plan: "Implement selected command files"
SW Engineer Thought: "Selected must mean the important ones, I'll do them all"
Actual Implementation:
  - ALL command files (not selected)
  - Complete utility modules (not minimal)
  - Full test coverage (not requested)
  - Workflow orchestration (not in plan)
Result: 2,215 lines instead of 650 (3.4x overrun)
Consequence: COMPLETE FAILURE - Split rejected, effort failed
```

### What Happens With Explicit Boundaries:
```
Plan: "Implement EXACTLY WithImage(), WithContext(), WithPlatform() - STOP"
SW Engineer Thinks: "Clear scope, 3 functions only"
Result: Exactly 3 functions, no extras
Actual: 200 lines as planned
Consequence: SPLIT SUCCEEDS, MOVE TO NEXT
```

## 📋 SPLIT PLAN CHECKLIST

Before marking split plan complete, verify:

- [ ] Listed EXACT function/method counts
- [ ] Provided specific function/method names
- [ ] Included explicit line estimates per component
- [ ] Added "DO NOT IMPLEMENT" section
- [ ] Specified test scope limitations
- [ ] Included STOP boundaries
- [ ] Provided realistic size calculations
- [ ] Added scope checkpoint commands
- [ ] Included good/bad examples

## 🔴 FOR CODE REVIEWERS - PREVENTING THE 3.4X PROBLEM

When creating split plans:
1. **NEVER USE VAGUE WORDS**: "selected", "necessary", "utilities", "handlers"
2. **LIST EXACT FILES**: Not "command files", but "cmd/push.go, cmd/build.go"
3. **LIST EXACT FUNCTIONS**: Not "options", but "WithImage(), WithContext()"
4. **SPECIFY EXACT TEST COUNT**: Not "write tests", but "3 tests" or "0 tests"
5. **CALCULATE REALISTICALLY**: 
   - Simple function: 10-20 lines
   - With validation: 30-50 lines
   - Complex logic: 70-100 lines
6. **ADD EXPLICIT FORBIDDEN LIST**: Think what SW might add and forbid it
7. **USE NUMBERS NOT WORDS**: "3 functions" not "a few functions"
8. **SPECIFY LINE LIMITS PER FILE**: Not just total

## 🔴 FOR SW ENGINEERS - YOUR MINDSET MUST CHANGE

### ❌ WRONG MINDSET (Causes 3.4x overruns):
- "I should complete the feature"
- "This seems incomplete, I'll add what's missing"
- "While I'm here, I'll add related functionality"
- "Selected probably means the important ones"

### ✅ CORRECT MINDSET (Stays within budget):
- "I will implement EXACTLY what's listed"
- "Incomplete is intentional - other splits handle the rest"  
- "If it's not explicitly listed, I won't add it"
- "Selected means ONLY the ones named specifically"

### YOUR IMPLEMENTATION CHECKLIST:
1. **READ the DO NOT list** - This is your primary guide
2. **COUNT exact items** - Functions, files, tests
3. **STOP at limits** - Even mid-function if needed
4. **NEVER ASSUME** - Ask if unclear, don't interpret
5. **RESIST COMPLETENESS** - Partial is the goal
6. **MEASURE FREQUENTLY** - Every 100 lines

## Template Version: 2.0
## Purpose: Prevent 3-5X implementation overruns through explicit boundaries