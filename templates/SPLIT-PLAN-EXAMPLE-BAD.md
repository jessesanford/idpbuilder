# EXAMPLE: BAD Split Plan (Causes 3-5X Overruns)

## ❌ THIS IS AN EXAMPLE OF WHAT NOT TO DO ❌

## Split Metadata
- **Split Number**: 001
- **Parent Effort**: buildkit-options
- **Original Size**: 1,247 lines
- **Target**: ~400 lines

## Implementation Scope (VAGUE - CAUSES OVERRUNS)

### What to Implement

Implement the functional options pattern for buildkit configuration. This should include the necessary option functions and configuration structure.

### Files to Create
- `pkg/buildkit/options.go` - Option functions
- `pkg/buildkit/config.go` - Configuration types
- `pkg/buildkit/options_test.go` - Tests

### Functionality
- Create option functions for configuration
- Define the Config struct
- Implement a constructor
- Add appropriate tests

## Technical Requirements

The options should follow Go best practices and be production-ready.

## Notes

Make sure the implementation is complete and follows standard patterns.

---

## 🚨 WHY THIS SPLIT PLAN FAILS 🚨

### Problem 1: No Explicit Counts
- **Bad**: "implement the necessary option functions"
- **Result**: SW implements 47 functions thinking they're all "necessary"
- **Should be**: "Implement EXACTLY 4 functions: WithImage, WithContext, WithPlatform, WithCache"

### Problem 2: No Boundaries
- **Bad**: "Add appropriate tests"
- **Result**: SW writes 1,000 lines of comprehensive tests
- **Should be**: "Write EXACTLY 5 tests, one per function, 25 lines each max"

### Problem 3: Vague Quality Standards
- **Bad**: "production-ready"
- **Result**: SW adds validation, error handling, logging, Clone(), marshaling
- **Should be**: "NO validation, NO Clone, NO extra methods in this split"

### Problem 4: "Complete" Language
- **Bad**: "Make sure the implementation is complete"
- **Result**: SW implements entire system to make it "complete"
- **Should be**: "This split is intentionally incomplete - implement ONLY what's listed"

### Problem 5: No DO NOT Section
- **Bad**: No forbidden items listed
- **Result**: SW adds everything they think is helpful
- **Should be**: Explicit "DO NOT IMPLEMENT" list

## ACTUAL RESULT OF THIS BAD PLAN

```
SW Engineer's interpretation:
"I need to make this production-ready and complete, so I'll add:"

Actual implementation:
- 47 option functions (instead of 4)
- Config with 25 methods (instead of 0)
- Validation logic (200 lines)
- Clone/Copy methods (100 lines)
- String/MarshalJSON (150 lines)
- Comprehensive tests (1,000 lines)
- Benchmarks (200 lines)
- Examples (150 lines)

Total: 2,847 lines (instead of 400)
Result: COMPLETE FAILURE - 7X overrun
```

## THE LESSON

**Vague instructions are interpreted as "make it complete"**

SW Engineers want to do good work. When instructions are vague, they fill in the gaps with what they think is best practice. This leads to massive overruns.

**BE EXPLICIT:**
- Count everything
- Name everything
- Forbid everything else
- Set hard boundaries
- Emphasize incompleteness is intentional