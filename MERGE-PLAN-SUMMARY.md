# PROJECT MERGE PLAN SUMMARY

## Quick Reference for Integration Agent

**Plan Location**: `/home/vscode/workspaces/idpbuilder-oci-build-push/efforts/project/integration-workspace/PROJECT-MERGE-PLAN.md`

## Merge Sequence
1. **image-builder** first (no dependencies)
2. **gitea-client** second (depends on image-builder)

## Key Commands
```bash
# Step 1: Merge image-builder
git merge origin/idpbuilder-oci-build-push/phase2/wave1/image-builder --no-ff

# Step 2: Validate
go mod tidy && go test ./pkg/build/...

# Step 3: Merge gitea-client  
git merge origin/idpbuilder-oci-build-push/phase2/wave1/gitea-client --no-ff

# Step 4: Final validation
go mod tidy && go test ./...
```

## Critical Success Criteria
- ✅ All imports use github.com/cnoe-io/idpbuilder
- ✅ No references to jessesanford/idpbuilder
- ✅ All tests pass
- ✅ go.mod is clean

## High-Risk Areas
- **pkg/certs/**: 20+ overlapping files
- **go.mod**: Dependency conflicts
- **Resolution**: See full plan for detailed strategies

---
**Created**: 2025-01-09 19:17:00 UTC
**Ready for Execution**: YES