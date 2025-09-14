# Fix Instructions: Duplicate TLSConfig Declarations

## Problem
Build failure due to duplicate struct declarations:
- TLSConfig struct appears in both pkg/certs/types.go and pkg/certs/utilities.go
- Syntax errors in pkg/certs/utilities.go at lines 131 and 159

## Required Fixes

### For cert-validation-split-001 branch
1. Clone target repository: https://github.com/jessesanford/idpbuilder.git
2. Checkout branch: idpbuilder-oci-build-push/phase1/wave2/cert-validation-split-001
3. Navigate to pkg/certs/
4. Verify TLSConfig exists in types.go (keep this one as the single source)
5. Remove any duplicate TLSConfig or DefaultTLSConfig structs from other files
6. Fix syntax errors if present
7. Test compilation: `go build ./pkg/certs`
8. Commit with message: "fix: remove duplicate TLSConfig declarations"
9. Push to remote

### For cert-validation-split-002 branch
1. Clone target repository: https://github.com/jessesanford/idpbuilder.git
2. Checkout branch: idpbuilder-oci-build-push/phase1/wave2/cert-validation-split-002
3. Check pkg/certs/utilities.go for:
   - Duplicate struct declarations
   - Syntax errors at lines 131 and 159
4. Fix any issues found
5. Test compilation
6. Commit and push if fixes made

### For cert-validation-split-003 branch
Same process as split-002

## Verification
- Ensure only ONE TLSConfig struct definition exists across all pkg/certs files
- All files must compile without errors
- Fixes MUST be in effort branches, NOT integration branches (R300)