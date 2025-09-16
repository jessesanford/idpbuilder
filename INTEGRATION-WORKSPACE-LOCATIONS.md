# SOFTWARE FACTORY 2.0 - INTEGRATION WORKSPACE LOCATIONS

## ✅ CORRECT INTEGRATION WORKSPACES

### 1. PROJECT Integration (READY ✅)
- **Path**: `/home/vscode/workspaces/idpbuilder-oci-build-push/efforts/project/project-integration-workspace`
- **Branch**: `idpbuilder-oci-build-push/project-integration-20250916-152718`
- **Remote**: `https://github.com/jessesanford/idpbuilder.git` (TARGET REPO ✅)
- **Status**: READY FOR MERGE EXECUTION
- **Contains**: PROJECT-MERGE-PLAN.md

### 2. Phase 2 Integration (READY ✅)
- **Path**: `/home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase2/phase-integration-workspace-new/repo`
- **Branch**: `idpbuilder-oci-build-push/phase2-integration-20250916-033720`
- **Remote**: `https://github.com/jessesanford/idpbuilder.git` (TARGET REPO ✅)
- **Status**: COMPLETE - Contains both Wave 1 and Wave 2
- **Contains**: All Phase 2 efforts integrated

### 3. Phase 1 Integration (FIXED ✅)
- **Path**: `/home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase1/phase-integration-workspace`
- **Branch**: `idpbuilder-oci-build-push/phase1-integration`
- **Remote**: `https://github.com/jessesanford/idpbuilder.git` (TARGET REPO ✅)
- **Status**: READY - Pushed to target repo and workspace created
- **Contains**: All Phase 1 efforts integrated

## 🔧 TO RESTORE ALL INTEGRATIONS

Run this command to restore all integration workspaces from target repo:

```bash
# All integrations are now correctly set up:

# Phase 1 Integration (READY ✅)
# Located at: efforts/phase1/phase-integration-workspace

# Phase 2 Integration (READY ✅)
# Located at: efforts/phase2/phase-integration-workspace-new/repo

# Project Integration (READY ✅)
# Located at: efforts/project/project-integration-workspace
```

## 📂 DIRECTORY STRUCTURE

```
efforts/
├── phase1/
│   ├── phase-integration-workspace/         # ⚠️ NEEDS FIX - wrong repo
│   ├── wave1/
│   │   ├── cert-extraction/
│   │   └── integration-workspace/
│   └── wave2/
│       ├── cert-validation-split-001/
│       ├── cert-validation-split-002/
│       ├── fallback-strategies/
│       └── integration-workspace/
├── phase2/
│   ├── phase-integration-workspace-new/
│   │   └── repo/                           # ✅ CORRECT Phase 2 integration
│   ├── wave1/
│   │   ├── builder/
│   │   ├── client-split-001/
│   │   ├── client-split-002/
│   │   └── integration-workspace/
│   └── wave2/
│       ├── commands/
│       ├── credentials/
│       ├── operations/
│       └── integration-workspace/
└── project/
    └── project-integration-workspace/      # ✅ CORRECT Project integration

```

## 🚨 CRITICAL FOR AGENTS

When spawning integration agents:

1. **For PROJECT Integration**: Use `/home/vscode/workspaces/idpbuilder-oci-build-push/efforts/project/project-integration-workspace`
2. **DO NOT** create new workspaces
3. **DO NOT** use worktrees
4. **DO NOT** change remotes
5. **JUST** cd to the directory and execute the merge plan

## 📝 NOTES

- All integration branches exist on TARGET repo: `https://github.com/jessesanford/idpbuilder.git`
- Planning repo is separate: `https://github.com/jessesanford/idpbuilder-oci-build-push.git`
- Integration workspaces should ALWAYS point to TARGET repo
- Effort workspaces should ALWAYS point to TARGET repo