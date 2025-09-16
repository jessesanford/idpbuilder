# PHASE 2 INTEGRATION MERGE PLAN (CORRECTED)

**Created**: 2025-09-16 02:35:00 UTC
**Purpose**: Correctly integrate Phase 2 by merging INTEGRATION branches, not individual efforts
**Critical**: This replaces the incorrect wave-level plan that was mistakenly used

## 🔴🔴🔴 CRITICAL CORRECTION 🔴🔴🔴

**PROBLEM IDENTIFIED**: The previous integration failed because it used a Wave-level merge plan that instructed merging individual effort branches instead of integration branches.

**THIS PLAN**: Correctly merges the two integration branches that contain ALL the work:
1. Phase 2 Wave 1 Integration (contains image-builder and gitea-client work)
2. Phase 2 Wave 2 Integration (contains cli-commands, credential-management, and image-operations)

## Executive Summary

This plan merges the COMPLETE Phase 2 work by using the proper integration branches that have already successfully integrated all individual efforts at the wave level.

## Integration Strategy

### What We're Merging (CORRECT APPROACH)

```
Phase 1 Integration (base)
    ├── Phase 2 Wave 1 Integration (idpbuilder-oci-build-push/phase2/wave1/integration-20250915-125755)
    │   ├── ✅ image-builder
    │   ├── ✅ gitea-client-split-001
    │   └── ✅ gitea-client-split-002
    │
    └── Phase 2 Wave 2 Integration (idpbuilder-oci-build-push/phase2/wave2/integration-20250916-002118)
        ├── ✅ cli-commands
        ├── ✅ credential-management (WITH --username/--token flags!)
        └── ✅ image-operations
```

### What NOT to Merge (INCORRECT APPROACH from old plan)
❌ DO NOT merge individual effort branches
❌ DO NOT merge phase2/wave2/cli-commands directly
❌ DO NOT merge phase2/wave2/credential-management directly
❌ DO NOT merge phase2/wave2/image-operations directly

## Pre-Integration Verification

### Step 1: Verify Integration Branches Exist
```bash
# Verify Wave 1 integration exists and has all efforts
cd /home/vscode/workspaces/idpbuilder-oci-build-push/worktrees/phase2-wave1-integration
git branch --show-current
# Should show: idpbuilder-oci-build-push/phase2/wave1/integration-20250915-125755

# Verify Wave 2 integration exists and has all efforts
cd /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase2/wave2/integration-workspace/repo
git branch --show-current
# Should show: idpbuilder-oci-build-push/phase2/wave2/integration-20250916-002118

# Verify credentials are present in Wave 2 integration
grep "pushUsername\|pushToken" pkg/cmd/push.go
# Should show the --username and --token flag definitions
```

### Step 2: Create New Phase 2 Integration Workspace
```bash
# Create clean workspace for Phase 2 integration
cd /home/vscode/workspaces/idpbuilder-oci-build-push
mkdir -p efforts/phase2/phase-integration-workspace-new
cd efforts/phase2/phase-integration-workspace-new

# Clone fresh copy
git clone https://github.com/jessesanford/idpbuilder.git repo
cd repo

# Start from Phase 1 integration as base
git checkout idpbuilder-oci-build-push/phase1/integration
git checkout -b idpbuilder-oci-build-push/phase2/integration-corrected-$(date +%Y%m%d-%H%M%S)
```

## Merge Sequence

### MERGE 1: Phase 2 Wave 1 Integration
```bash
# First merge Wave 1 integration (contains image-builder and gitea-client work)
git fetch origin idpbuilder-oci-build-push/phase2/wave1/integration-20250915-125755
git merge origin/idpbuilder-oci-build-push/phase2/wave1/integration-20250915-125755 \
    -m "integrate: Phase 2 Wave 1 integration (image-builder, gitea-client)"

# Verify build works
go build ./...
go test ./pkg/build/... ./pkg/gitea/...
```

### MERGE 2: Phase 2 Wave 2 Integration
```bash
# Then merge Wave 2 integration (contains cli-commands, credential-management, image-operations)
git fetch origin idpbuilder-oci-build-push/phase2/wave2/integration-20250916-002118
git merge origin/idpbuilder-oci-build-push/phase2/wave2/integration-20250916-002118 \
    -m "integrate: Phase 2 Wave 2 integration (cli-commands with credentials, image-operations)"

# Verify credentials were included
grep "pushUsername" pkg/cmd/push.go || echo "ERROR: Credentials missing!"
grep "pushToken" pkg/cmd/push.go || echo "ERROR: Credentials missing!"

# Verify all components build
go build ./...
```

## Critical Verification Points

### After Wave 1 Merge
- [ ] pkg/build/builder.go exists
- [ ] pkg/gitea/client.go exists (basic version)
- [ ] Tests pass for build and gitea packages

### After Wave 2 Merge
- [ ] pkg/cmd/build.go exists
- [ ] pkg/cmd/push.go exists WITH --username and --token flags
- [ ] pkg/gitea/credentials.go exists
- [ ] pkg/gitea/config.go exists
- [ ] pkg/gitea/keyring.go exists
- [ ] pkg/operations/ directory exists with real implementations
- [ ] ALL tests pass

## Post-Integration Testing

```bash
# Build the binary
go build -o idpbuilder .

# Verify credential flags exist
./idpbuilder push --help | grep -E "username|token"
# MUST show both flags

# Verify build command exists
./idpbuilder build --help
# MUST show build command with proper options

# Run all tests
go test ./...
# ALL must pass
```

## Integration Success Criteria

The integration is ONLY successful if:
1. ✅ Both Wave integration branches are merged (not individual efforts)
2. ✅ The binary has --username and --token flags on push command
3. ✅ All credential management files exist (credentials.go, config.go, keyring.go)
4. ✅ No hardcoded "admin"/"password" credentials remain
5. ✅ All tests pass
6. ✅ Binary builds without errors

## Rollback Plan

If integration fails:
```bash
# Delete the failed branch
git checkout main
git branch -D idpbuilder-oci-build-push/phase2/integration-corrected-*
git push origin --delete idpbuilder-oci-build-push/phase2/integration-corrected-*

# Start over with this plan
```

## Important Notes for Integration Agent

1. **USE INTEGRATION BRANCHES**: You are merging completed integration branches, NOT individual effort branches
2. **VERIFY CREDENTIALS**: After merging Wave 2, you MUST verify credential flags are present
3. **NO MANUAL SELECTION**: The branch names are exact - use them as specified
4. **CHECK WORKTREES**: The integration branches may be in worktrees, not the main repo
5. **FOLLOW ORDER**: Wave 1 first, then Wave 2 - this order is critical

## Why Previous Integration Failed

The previous integration failed because:
1. Wrong plan was used (Wave-level instead of Phase-level)
2. Plan specified individual effort branches instead of integration branches
3. Plan was created before Wave 2 integration was complete
4. Integration agent followed the incorrect instructions exactly

This plan corrects all these issues by:
1. Being a Phase-level integration plan
2. Specifying the correct integration branches
3. Being created after all integrations are complete
4. Providing clear verification steps

---
END OF MERGE PLAN - Integration Agent should execute EXACTLY as specified above.