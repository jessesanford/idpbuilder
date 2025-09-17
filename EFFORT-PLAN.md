# Phase 1 Wave 3: Upstream Fixes Effort

## Overview
Fix pre-existing test failures in original idpbuilder codebase blocking R291 gates.

## Base
- Branch: idpbuilder-oci-build-push/phase1/wave3/upstream-fixes  
- Base: origin/idpbuilder-oci-build-push/phase1/wave2-integration

## Fixes Required
1. pkg/kind: Create missing NewCluster, IProvider
2. cmd/: Create missing main.go entry point
3. pkg/cmd/get, pkg/util, pkg/controllers: Fix test failures

## Success Criteria
- All tests pass (go test ./...)
- Total < 800 lines
