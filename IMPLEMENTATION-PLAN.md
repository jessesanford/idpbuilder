# Split-003A Implementation Requirements

This is split-003a of the resplit from oversized split-003 (1120 lines).
Split-002a and split-002b have already been completed.

## Scope: Core Builder Functionality (569 lines target)

### IMPLEMENT:
1. **builder.go** (174 lines) - Builder interface and SimpleBuilder
2. **options.go** (124 lines) - Build options
3. **Minimal test stubs** (~271 lines) - Basic tests

**Total Target**: 569 lines (well under 800 limit)

### Dependencies:
- Build on split-002b (which has tarball.go and builder_impl.go)
- Base branch: split-002b

### Implementation Notes:
- builder.go exists with stub - enhance with full SimpleBuilder
- options.go - implement build option structures
- Add minimal test coverage to ensure functionality

### DO NOT IMPLEMENT:
- TLS/certificate code (reserved for split-003b)
- Comprehensive test suites (defer to later)

### Size Requirements:
- Target: 569 lines
- Hard limit: 800 lines
- Use line-counter.sh to verify

### Implementation Order:
1. Enhance builder.go with SimpleBuilder
2. Implement options.go fully
3. Add basic test stubs
4. Verify size with line-counter.sh
