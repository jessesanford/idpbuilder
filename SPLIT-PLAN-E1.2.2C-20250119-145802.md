# SPLIT-PLAN-E1.2.2C: fallback-security

## Split 003 of 003: Security and Insecure Mode
**Planner**: Code Reviewer Agent
**Parent Effort**: E1.2.2 fallback-strategies
**Target Size**: ~499 lines (CRITICAL: Original 833 lines requires optimization)
**Priority**: THIRD (can parallelize with E1.2.2B after E1.2.2A)

<!-- ORCHESTRATOR METADATA PLACEHOLDER - DO NOT REMOVE -->
<!-- The orchestrator will add infrastructure metadata below: -->
<!-- WORKING_DIRECTORY, BRANCH, REMOTE, BASE_BRANCH, etc. -->
<!-- SW Engineers MUST read this metadata to navigate to the correct directory -->
<!-- END PLACEHOLDER -->

## Boundaries
- **Previous Splits**:
  - Split 001 (fallback-core) of phase1/wave2/fallback-strategies
    - Path: efforts/phase1/wave2/fallback-core/
    - Branch: phase1/wave2/fallback-core
  - Split 002 (fallback-recommendations) of phase1/wave2/fallback-strategies
    - Path: efforts/phase1/wave2/fallback-recommendations/
    - Branch: phase1/wave2/fallback-recommendations
- **This Split**: Split 003 of phase1/wave2/fallback-strategies
  - Path: efforts/phase1/wave2/fallback-security/
  - Branch: phase1/wave2/fallback-security
- **Next Split**: None (final split)

## Critical Size Constraint
**WARNING**: Original files total 833 lines (security_log.go: 409, insecure.go: 424)
**REQUIREMENT**: Must optimize to <500 lines through:
1. Extracting common utilities
2. Consolidating redundant code
3. Focusing on essential functionality
4. Moving verbose logging to configuration

## Files in This Split (EXCLUSIVE - OPTIMIZED)
```
pkg/fallback/security/
├── log.go           # 250 lines - Security logging (reduced from 409)
├── insecure.go      # 200 lines - Insecure mode (reduced from 424)
├── security_test.go # ~49 lines - Unit tests
Total: ~499 lines (MUST NOT EXCEED 500)
```

## Functionality Scope

### Core Components (ESSENTIAL ONLY)
1. **Security Logger** (`log.go` - 250 lines MAX)
   - Security event recording
   - Audit trail generation
   - Compliance logging
   - Risk assessment
   - **REMOVED**: Verbose formatting, redundant checks

2. **Insecure Mode Manager** (`insecure.go` - 200 lines MAX)
   - Insecure fallback handling
   - Risk acceptance tracking
   - Warning generation
   - Override mechanisms
   - **REMOVED**: Duplicate validations, verbose documentation

## Optimization Strategy

### Required Refactoring from Original
1. **Extract Common Code** (saves ~150 lines)
   ```go
   // Original: Duplicate validation in both files
   // Optimized: Single shared validator
   type SecurityValidator struct {
       // Consolidated validation logic
   }
   ```

2. **Consolidate Logging** (saves ~80 lines)
   ```go
   // Original: Separate log formatting in each file
   // Optimized: Shared formatter
   func formatSecurityLog(event Event) string {
       // Unified formatting
   }
   ```

3. **Remove Verbosity** (saves ~100 lines)
   ```go
   // Original: Extensive inline documentation
   // Optimized: Concise implementation with external docs
   ```

## Implementation Instructions

### Step 1: Repository Setup
```bash
cd /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase1/wave2
mkdir -p fallback-security
cd fallback-security
git init
git remote add origin [same as parent]
git checkout -b phase1/wave2/fallback-security
```

### Step 2: Import Dependencies from E1.2.2A
```bash
mkdir -p pkg/fallback/security
cd pkg/fallback/security
```

### Step 3: Implement Optimized Security Components

#### 3.1 Create `log.go` (250 lines MAXIMUM)
```go
package security

import (
    "github.com/idpbuilder/idpbuilder/pkg/fallback" // from E1.2.2A
)

// SecurityLogger implements the SecurityLogger interface (40 lines)
type SecurityLogger struct {
    config    *Config
    writer    LogWriter
    validator *SecurityValidator // Shared
}

// Core sections (OPTIMIZED):
// 1. Event Recording (60 lines)
// 2. Audit Trail (50 lines)
// 3. Risk Assessment (40 lines)
// 4. Compliance Checks (30 lines)
// 5. Public API (30 lines)
// Total: 250 lines
```

#### 3.2 Create `insecure.go` (200 lines MAXIMUM)
```go
package security

// InsecureManager handles insecure fallback operations (30 lines)
type InsecureManager struct {
    logger    *SecurityLogger
    validator *SecurityValidator // Shared
    overrides map[string]Override
}

// Core sections (OPTIMIZED):
// 1. Mode Management (50 lines)
// 2. Risk Acceptance (40 lines)
// 3. Warning System (40 lines)
// 4. Override Handling (30 lines)
// 5. Public API (10 lines)
// Total: 200 lines
```

#### 3.3 Write Minimal Tests (49 lines)
Focus on critical paths only:
- Security event logging
- Insecure mode activation
- Risk assessment accuracy

### Step 4: Aggressive Size Validation
```bash
# Count lines BEFORE committing
find pkg -name "*.go" -exec wc -l {} \; | awk '{sum+=$1} END {print "Total:", sum}'

# If > 499 lines, MUST refactor more:
# - Remove comments
# - Consolidate functions
# - Extract to shared utilities
```

### Step 5: Final Validation
```bash
# Size check with official tool
PROJECT_ROOT=$(pwd)
while [ "$PROJECT_ROOT" != "/" ]; do
    [ -f "$PROJECT_ROOT/orchestrator-state.json" ] && break
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
done
$PROJECT_ROOT/tools/line-counter.sh

# MUST be < 500 lines or split fails review
```

## Dependencies

### Internal Dependencies
- E1.2.2A (fallback-core):
  - Import type definitions
  - Implement `SecurityLogger` interface
  - Use `SecurityEvent` types

### External Dependencies (go.mod)
```go
require (
    github.com/idpbuilder/idpbuilder v0.0.0 // local replace
    github.com/sirupsen/logrus v1.9.0
)

replace github.com/idpbuilder/idpbuilder => ../fallback-core
```

## Success Criteria

1. **SIZE COMPLIANCE**: Total lines < 500 (MANDATORY)
2. **Interface Implementation**: Fully implements `SecurityLogger`
3. **Functionality**: Core security features working
4. **Tests**: Critical paths covered
5. **No Stubs**: Complete implementation
6. **Performance**: Logging overhead < 5ms

## Fallback Plan (if size cannot be reduced)

If optimization cannot achieve <500 lines, create 4th split:
- E1.2.2C: fallback-security-log (409 lines as-is)
- E1.2.2D: fallback-insecure (424 lines as-is)

**NOTE**: 4-split approach requires orchestrator approval

## Risk Mitigation

- **Risk**: Cannot reduce to <500 lines
- **Mitigation**: 4-split fallback plan ready

- **Risk**: Lost functionality during optimization
- **Mitigation**: Document removed features for future enhancement

## Review Focus Areas

1. **SIZE LIMIT ADHERENCE** (Critical)
2. Security feature completeness
3. Interface compliance
4. No functionality regression
5. Audit trail integrity

## Parallel Development Notes

This split can be developed in parallel with E1.2.2B once E1.2.2A is complete:
- No dependencies on E1.2.2B
- Both depend only on E1.2.2A
- Integration testing after both complete

## Critical Implementation Notes

**SW ENGINEER MUST**:
1. Count lines frequently during implementation
2. Refactor aggressively to stay under 500 lines
3. Consider 4-split approach early if size problematic
4. Prioritize essential functionality only
5. Move verbose code to configuration/documentation