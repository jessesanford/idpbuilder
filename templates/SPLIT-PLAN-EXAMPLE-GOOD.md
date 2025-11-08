# EXAMPLE: GOOD Split Plan with Explicit Boundaries

## Split Metadata
- **Split Number**: 001
- **Parent Effort**: buildkit-options
- **Original Size**: 1,247 lines (exceeded 800)
- **This Split Target**: 400 lines (leaving buffer)
- **Created**: 2025-01-21 10:00:00

## 🚨 EXPLICIT SCOPE DEFINITION (MANDATORY)

### EXACTLY What to Implement (BE SPECIFIC!)

#### Functions to Create (EXACTLY 5 - NO MORE)
```go
1. WithImage(image string) Option              // ~15 lines - Sets container image
2. WithContext(ctx context.Context) Option     // ~15 lines - Sets build context
3. WithPlatform(platform string) Option        // ~20 lines - Sets target platform
4. WithCache(cacheDir string) Option           // ~20 lines - Sets cache directory
5. NewConfig(opts ...Option) (*Config, error)  // ~40 lines - Creates config with options
// STOP HERE - DO NOT ADD MORE OPTIONS
```

#### Structs to Define (EXACTLY 2)
```go
// Option function type
type Option func(*Config) error  // Simple function type, no methods

// Config struct - EXACTLY these 6 fields, NO methods in this split
type Config struct {
    Image    string           // Container image name
    Context  context.Context  // Build context
    Platform string          // Target platform
    Cache    string          // Cache directory path
    Debug    bool           // Debug mode flag
    created  time.Time      // Internal timestamp
}
// NO METHODS ON THESE TYPES IN THIS SPLIT
```

#### Files to Create/Modify
```
pkg/buildkit/options.go - CREATE - 110 lines MAX
  - Option type definition (5 lines)
  - 4 option functions listed above (80 lines total)
  - NewConfig constructor (25 lines)
  - NO validation methods, NO helpers

pkg/buildkit/config.go - CREATE - 40 lines MAX
  - Config struct definition only (40 lines)
  - NO methods, NO validation, NO constructors here

pkg/buildkit/options_test.go - CREATE - 150 lines MAX
  - Test each option function (4 tests × 25 lines = 100 lines)
  - Test NewConfig with options (50 lines)
  - NO benchmarks, NO examples, NO edge cases
```

### 🛑 STOP BOUNDARIES - DO NOT IMPLEMENT

**EXPLICITLY FORBIDDEN IN THIS SPLIT:**
- ❌ DO NOT add Validate() or IsValid() methods
- ❌ DO NOT add Clone() or Copy() methods
- ❌ DO NOT add String() or MarshalJSON() methods
- ❌ DO NOT implement WithTimeout(), WithRetry(), etc. (not listed above)
- ❌ DO NOT add comprehensive error handling (basic only)
- ❌ DO NOT write edge case tests or benchmarks
- ❌ DO NOT add logging or debugging utilities
- ❌ DO NOT refactor or "improve" the design
- ❌ DO NOT add comments beyond minimal documentation
- ❌ DO NOT create helper functions not explicitly listed

### Test Scope (MINIMAL)
```go
// options_test.go - EXACTLY these tests, no more:
func TestWithImage(t *testing.T)    // 25 lines - Basic functionality only
func TestWithContext(t *testing.T)   // 25 lines - Basic functionality only
func TestWithPlatform(t *testing.T)  // 25 lines - Basic functionality only
func TestWithCache(t *testing.T)     // 25 lines - Basic functionality only
func TestNewConfig(t *testing.T)     // 50 lines - Test with 2-3 options max
// STOP - No edge cases, no error scenarios beyond basics
```

## 📊 SIZE CALCULATION (REALISTIC)

```
Component Breakdown:
- Option type definition:        5 lines
- WithImage function:           15 lines
- WithContext function:         15 lines
- WithPlatform function:        20 lines
- WithCache function:           20 lines
- NewConfig function:           40 lines
- Config struct:                40 lines
- Basic tests (5 × 30):        150 lines

TOTAL: 305 lines (well under 400 target)
```

## 🔴 MANDATORY ADHERENCE CHECKPOINTS

### Before Starting:
```bash
echo "SPLIT 001 SCOPE LOCKED:"
echo "✓ Functions: EXACTLY 5 (WithImage, WithContext, WithPlatform, WithCache, NewConfig)"
echo "✓ Structs: EXACTLY 2 (Option type, Config struct)"
echo "✓ Methods: EXACTLY 0 (none in this split)"
echo "✓ Tests: EXACTLY 5 (one per function)"
echo "✗ Validation: NONE"
echo "✗ Clone/Copy: NONE"
echo "✗ Extra options: NONE"
```

### During Implementation:
```bash
# Check after each file
FUNC_COUNT=$(grep -c "^func With" options.go 2>/dev/null || echo 0)
if [ "$FUNC_COUNT" -gt 4 ]; then
    echo "⚠️ WARNING: Exceeding function count! Stop adding!"
fi
```

### After Implementation:
```bash
# Final validation (tool auto-detects base)
TOTAL_LINES=$(../../tools/line-counter.sh | grep Total | awk '{print $NF}')
if [ "$TOTAL_LINES" -gt 400 ]; then
    echo "❌ EXCEEDED TARGET: $TOTAL_LINES > 400"
    echo "Remove features to get under limit!"
fi
```

## Success Criteria
- [ ] Implemented EXACTLY 5 functions (no more, no less)
- [ ] Created EXACTLY 2 types (no more, no less)
- [ ] Wrote EXACTLY 5 tests (no more, no less)
- [ ] Total lines under 400
- [ ] NO validation logic added
- [ ] NO methods on structs
- [ ] NO extra "helpful" features

## Notes for SW Engineer

**IMPORTANT**: This split intentionally creates an incomplete implementation. The Config struct has no methods, validation is missing, and only 4 options are implemented. This is BY DESIGN. Other splits will add:
- Split-002: Validation methods and error handling
- Split-003: Additional options (WithTimeout, WithRetry, etc.)
- Split-004: Clone, String, and marshal methods

**DO NOT** try to make this "complete" or "production-ready". Implement EXACTLY what's specified above and stop.

## Integration Notes

This split creates the foundation. The next split (002) will branch from this and add validation. Each split builds on the previous sequentially.