# Upstream Bugs Identified During Integration

## Bug 1: Duplicate Type Declarations
**Affected Files**: 
- pkg/certs/validator.go (from E1.2.1)
- pkg/certs/types.go (from E1.2.2)

**Issue**: 
E1.2.2 created pkg/certs/types.go which duplicates type declarations that already exist in pkg/certs/validator.go from E1.2.1:
- CertValidator interface (line 40 in validator.go, line 11 in types.go)
- CertDiagnostics struct (line 56 in validator.go, line 26 in types.go)
- ValidationError type (line 69 in validator.go, line 39 in types.go)

**Impact**: Build failure with "redeclared in this block" errors

**Recommendation**: 
Remove duplicate declarations from pkg/certs/types.go since validator.go already contains these definitions. The types.go file appears to be unnecessary as E1.2.1 already provided these interfaces.

**Status**: NOT FIXED (upstream issue - documented for resolution by development team)