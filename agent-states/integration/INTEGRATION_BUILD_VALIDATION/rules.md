# INTEGRATION_BUILD_VALIDATION State Rules

## State Purpose
Validate that the integrated code builds successfully - the first major quality gate.

## Entry Conditions
- All branches successfully merged
- No unresolved conflicts
- Integration branch is complete
- Build validation is required (per requirements)

## Required Actions

### 1. Prepare Build Environment
```bash
# Clean any previous build artifacts
make clean 2>/dev/null || true
rm -rf build/ dist/ target/ 2>/dev/null || true

# Ensure all dependencies are available
if [[ -f "requirements.txt" ]]; then
    pip install -r requirements.txt
elif [[ -f "package.json" ]]; then
    npm install
elif [[ -f "go.mod" ]]; then
    go mod download
fi
```

### 2. Set Build Configuration
```bash
# Set build flags for integration testing
export BUILD_TYPE="integration"
export ENABLE_TESTS="true"
export VERBOSE_OUTPUT="true"

# Record build start time
BUILD_START=$(date +%s)
```

### 3. Capture Build Environment
```json
{
  "build_validation": {
    "started_at": "2025-01-21T10:30:00Z",
    "attempt": 2,
    "environment": {
      "build_type": "integration",
      "compiler": "gcc 9.3.0",
      "flags": "-O2 -Wall -Werror"
    }
  }
}
```

### 4. Transition Decision
- Proceed to INTEGRATION_RUN_BUILD to execute build

## Critical Validation Rules

### Build Must Be Clean
```bash
# No warnings allowed in integration
export CFLAGS="-Wall -Werror"
export GOFLAGS="-buildvcs=false"

# Treat warnings as errors
export CI=true
export INTEGRATION_BUILD=true
```

### Build Must Be Complete
- All modules must build
- All dependencies must resolve
- All artifacts must be generated

### Build Must Be Reproducible
```bash
# Clear ccache and other caches
ccache -C 2>/dev/null || true

# Use consistent timestamps
export SOURCE_DATE_EPOCH=$(date +%s)
```

## State Tracking Updates
```json
{
  "validation_status": {
    "merge": "COMPLETE",
    "build": "IN_PROGRESS",
    "unit_tests": "PENDING",
    "functional_tests": "PENDING"
  },
  "build_config": {
    "clean_build": true,
    "warnings_as_errors": true,
    "parallel_jobs": 4
  }
}
```

## Quality Gates

### Gate 1: Pre-Build Checks
- ✅ All merges complete
- ✅ No conflict markers in code
- ✅ Dependencies available
- ✅ Build system ready

### Gate 2: Build Configuration
- ✅ Appropriate flags set
- ✅ Error handling configured
- ✅ Output capture ready
- ✅ Timeout configured

### Gate 3: Environment Validation
- ✅ Correct branch checked out
- ✅ Clean workspace
- ✅ Tools available
- ✅ Disk space sufficient

## Error Handling

### Pre-Build Failures
- Missing dependencies → Attempt to install
- Dirty workspace → Clean and retry
- Wrong branch → Checkout correct branch

### Build Preparation Failures
- Cannot clean → Log warning, continue
- Cannot set environment → INTEGRATION_ERROR
- Missing tools → INTEGRATION_ABORT

## Logging Requirements
```bash
echo "[INTEGRATION_BUILD_VALIDATION] Preparing build validation"
echo "[INTEGRATION_BUILD_VALIDATION] Attempt: ${ATTEMPT_NUMBER}"
echo "[INTEGRATION_BUILD_VALIDATION] Build type: ${BUILD_TYPE}"
echo "[INTEGRATION_BUILD_VALIDATION] Merged branches: ${MERGED_COUNT}"

# Log build configuration
echo "Build Configuration:"
echo "  - Clean build: YES"
echo "  - Warnings as errors: YES"
echo "  - Parallel jobs: ${PARALLEL_JOBS:-4}"
echo "  - Timeout: ${BUILD_TIMEOUT:-30m}"
```

## Metrics to Track
- Build preparation time
- Configuration complexity
- Environment setup duration
- Cache hit rates

## Common Issues

### Issue 1: Dependency Conflicts
**Problem**: Merged branches have conflicting dependencies
**Solution**: Resolve in source branches, re-attempt integration

### Issue 2: Build Tool Version Mismatch
**Problem**: Different branches expect different tool versions
**Solution**: Use project-standard versions, fix in source branches

### Issue 3: Missing Build Configuration
**Problem**: No Makefile, package.json, or build config
**Solution**: Check if build is actually required, or fix in source

## Success Criteria
✅ Build environment prepared
✅ All dependencies available
✅ Build configuration set
✅ Ready to execute build

## Next State
- Always transitions to INTEGRATION_RUN_BUILD