# ⚠️ INTEGRATION BLOCKED - UPSTREAM BUG ⚠️

## Critical Issue
**Integration cannot proceed due to duplicate type declaration**

## Quick Fix Required
**File**: `pkg/certs/trust_store.go`
**Line**: 21-28
**Action**: DELETE the duplicate `CertificateInfo` struct

```go
// DELETE THESE LINES (21-28):
type CertificateInfo struct {
	Subject   string
	Issuer    string
	NotBefore time.Time
	NotAfter  time.Time
	IsCA      bool
	DNSNames  []string
}
```

## Location of Fix Needed
The fix must be applied in the E1.1.2 effort workspace:
`/home/vscode/workspaces/idpbuilder-oci-go-cr/efforts/phase1/wave1/registry-tls-trust-integration/pkg/certs/trust_store.go`

## Why This Happened
- E1.1.1 defines `CertificateInfo` in `types.go` (correct)
- E1.1.2 was supposed to import from E1.1.1 but still has duplicate in `trust_store.go`
- The TODO comment shows awareness but the duplicate wasn't removed

## Next Steps
1. SW Engineer fixes duplicate in E1.1.2 workspace
2. Commits the fix
3. Integration Agent re-runs merge (will be clean)
4. Build and tests should pass
5. Integration can complete

## Integration Agent Compliance
Per R266, I have documented but NOT fixed this upstream bug.
The fix must come from the responsible SW Engineer.