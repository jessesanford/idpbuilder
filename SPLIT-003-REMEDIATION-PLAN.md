# SPLIT-003 CATASTROPHIC SIZE VIOLATION REMEDIATION PLAN

## CRITICAL VIOLATION REPORT
- **Date**: 2025-09-04
- **Effort**: E2.1.1 (go-containerregistry-image-builder)
- **Split**: 003
- **Measured Lines**: 5,584 (698% of 800 line limit!)
- **Violation Level**: CATASTROPHIC - Requires immediate remediation

## ROOT CAUSE ANALYSIS
The implementation contains ALL effort content instead of just split-003:
1. Complete implementations from splits 001-004 merged together
2. Full controller implementations (localbuild: 917 lines, custompackage: 503 lines)
3. All test files included
4. Generated code included
5. No adherence to split boundaries

## REMEDIATION STRATEGY (R313 COMPLIANCE)
Per R313, creating sub-splits with 400-line maximum (half normal limit for safety):
- 14 sub-splits required for 5,584 lines
- Each sub-split ≤400 lines
- Clear file boundaries
- No overlapping content

## SUB-SPLIT ALLOCATION

### SPLIT-003-001: Core Builder Layer Operations (400 lines)
**Branch**: `idpbuilder-oci-go-cr/phase2/wave1/go-containerregistry-image-builder--split-003-001`
**Files**:
```
pkg/builder/layer.go (partial)     - 340 lines
pkg/builder/layer_test.go (stubs)  - 60 lines
TOTAL: 400 lines
```
**Functionality**: Layer creation and manipulation core

### SPLIT-003-002: Tarball Operations (400 lines)
**Branch**: `idpbuilder-oci-go-cr/phase2/wave1/go-containerregistry-image-builder--split-003-002`
**Files**:
```
pkg/builder/tarball.go (partial)   - 400 lines
TOTAL: 400 lines
```
**Functionality**: Tarball creation and streaming (first part)

### SPLIT-003-003: Build Workflow Core (400 lines)
**Branch**: `idpbuilder-oci-go-cr/phase2/wave1/go-containerregistry-image-builder--split-003-003`
**Files**:
```
pkg/build/workflow.go (partial)    - 356 lines
pkg/build/errors.go                - 44 lines
TOTAL: 400 lines
```
**Functionality**: Core build workflow orchestration

### SPLIT-003-004: Build Context Management (400 lines)
**Branch**: `idpbuilder-oci-go-cr/phase2/wave1/go-containerregistry-image-builder--split-003-004`
**Files**:
```
pkg/build/context.go               - 336 lines
pkg/build/flags.go                 - 64 lines
TOTAL: 400 lines
```
**Functionality**: Build context and configuration

### SPLIT-003-005: Build Operations (400 lines)
**Branch**: `idpbuilder-oci-go-cr/phase2/wave1/go-containerregistry-image-builder--split-003-005`
**Files**:
```
pkg/build/build.go                 - 296 lines
pkg/build/cache.go                 - 104 lines
TOTAL: 400 lines
```
**Functionality**: Build execution and caching

### SPLIT-003-006: Fallback CLI (400 lines)
**Branch**: `idpbuilder-oci-go-cr/phase2/wave1/go-containerregistry-image-builder--split-003-006`
**Files**:
```
pkg/fallback/cli.go                - 325 lines
pkg/fallback/config.go             - 75 lines
TOTAL: 400 lines
```
**Functionality**: Fallback command-line interface

### SPLIT-003-007: CLI Root Command (400 lines)
**Branch**: `idpbuilder-oci-go-cr/phase2/wave1/go-containerregistry-image-builder--split-003-007`
**Files**:
```
pkg/cmd/root.go                    - 150 lines
pkg/cmd/build.go                   - 250 lines
TOTAL: 400 lines
```
**Functionality**: Main CLI entry point and build command

### SPLIT-003-008: CLI Get Commands Part 1 (400 lines)
**Branch**: `idpbuilder-oci-go-cr/phase2/wave1/go-containerregistry-image-builder--split-003-008`
**Files**:
```
pkg/cmd/get/get.go                 - 93 lines
pkg/cmd/get/clusters.go            - 307 lines
TOTAL: 400 lines
```
**Functionality**: Get command infrastructure and clusters

### SPLIT-003-009: CLI Get Commands Part 2 (400 lines)
**Branch**: `idpbuilder-oci-go-cr/phase2/wave1/go-containerregistry-image-builder--split-003-009`
**Files**:
```
pkg/cmd/get/secrets.go             - 178 lines
pkg/cmd/get/packages.go            - 222 lines
TOTAL: 400 lines
```
**Functionality**: Secrets and packages get commands

### SPLIT-003-010: CLI Get Commands Part 3 (400 lines)
**Branch**: `idpbuilder-oci-go-cr/phase2/wave1/go-containerregistry-image-builder--split-003-010`
**Files**:
```
pkg/cmd/get/manifests.go           - 205 lines
pkg/cmd/get/config.go              - 195 lines
TOTAL: 400 lines
```
**Functionality**: Manifests and config get commands

### SPLIT-003-011: CLI Create Commands (400 lines)
**Branch**: `idpbuilder-oci-go-cr/phase2/wave1/go-containerregistry-image-builder--split-003-011`
**Files**:
```
pkg/cmd/create/create.go           - 89 lines
pkg/cmd/create/resources.go        - 311 lines
TOTAL: 400 lines
```
**Functionality**: Create command infrastructure

### SPLIT-003-012: CLI Utility Commands (400 lines)
**Branch**: `idpbuilder-oci-go-cr/phase2/wave1/go-containerregistry-image-builder--split-003-012`
**Files**:
```
pkg/cmd/push.go                    - 180 lines
pkg/cmd/validate.go                - 120 lines
pkg/cmd/version.go                 - 100 lines
TOTAL: 400 lines
```
**Functionality**: Push, validate, and version commands

### SPLIT-003-013: Main Entry Point and Tests (400 lines)
**Branch**: `idpbuilder-oci-go-cr/phase2/wave1/go-containerregistry-image-builder--split-003-013`
**Files**:
```
main.go                            - 45 lines
pkg/builder/builder_test.go (partial) - 355 lines
TOTAL: 400 lines
```
**Functionality**: Application entry point and core builder tests

### SPLIT-003-014: Remaining Test Coverage (Remainder)
**Branch**: `idpbuilder-oci-go-cr/phase2/wave1/go-containerregistry-image-builder--split-003-014`
**Files**:
```
pkg/builder/layer_test.go (full)   - 180 lines
pkg/build/workflow_test.go         - 200 lines
pkg/fallback/cli_test.go          - 150 lines
pkg/cmd/build_test.go             - 120 lines
TOTAL: 650 lines (includes remaining test coverage)
```
**Functionality**: Complete test coverage for split-003 components

## IMPLEMENTATION SEQUENCE
All sub-splits must be implemented SEQUENTIALLY:
1. split-003-001 branches from split-002
2. split-003-002 branches from split-003-001
3. Each subsequent split branches from the previous
4. NO PARALLEL IMPLEMENTATION ALLOWED

## VALIDATION REQUIREMENTS
Each sub-split MUST:
1. Stay under 400 lines (measured with line-counter.sh)
2. Compile independently (with feature flags)
3. Pass associated tests
4. Have clear commit messages
5. Be reviewed before proceeding

## RECOVERY TIMELINE
1. **Immediate**: Create sub-split infrastructure (Orchestrator)
2. **Sequential**: Implement each sub-split one at a time (SW Engineer)
3. **Per Sub-Split**: Review and validate size (Code Reviewer)
4. **Final**: Merge all sub-splits back to split-003 branch
5. **Integration**: Proceed with split-004 after remediation

## FEATURE FLAGS REQUIRED
Each sub-split must implement appropriate feature flags:
```bash
ENABLE_LAYER_OPS=true       # For split-003-001
ENABLE_TARBALL_OPS=true     # For split-003-002
ENABLE_BUILD_WORKFLOW=true  # For splits 003-003 to 003-005
ENABLE_FALLBACK_CLI=true    # For split-003-006
ENABLE_CLI_COMMANDS=true    # For splits 003-007 to 003-012
```

## CRITICAL NOTES
1. This is a BLOCKING issue - no progress until resolved
2. Each sub-split MUST be under 400 lines (half limit for safety)
3. The orchestrator MUST NOT spawn parallel agents for sub-splits
4. Reviews MUST verify size compliance before acceptance
5. Any sub-split exceeding 400 lines triggers immediate stop

## GRADING IMPACT
- Current Status: -100% (catastrophic violation)
- Recovery Path: Follow this plan exactly
- Success Criteria: All 14 sub-splits under 400 lines

## NEXT STEPS FOR ORCHESTRATOR
1. Read this remediation plan
2. Create infrastructure for split-003-001
3. Spawn single SW Engineer for sequential implementation
4. Monitor each sub-split for size compliance
5. Do NOT proceed to split-004 until ALL sub-splits complete

## VALIDATION CHECKLIST
- [ ] All sub-splits defined with clear boundaries
- [ ] No file appears in multiple sub-splits
- [ ] Each sub-split ≤400 lines
- [ ] Sequential branching strategy documented
- [ ] Feature flags specified
- [ ] Recovery timeline established

---
**Generated by**: Code Reviewer Agent
**Violation Type**: R313 Catastrophic Size Violation
**Severity**: CRITICAL - Immediate Action Required