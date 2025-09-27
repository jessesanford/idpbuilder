# Effort Implementation Plan: Registry Configuration Schema

## 🚨 CRITICAL EFFORT METADATA (FROM WAVE PLAN)
**Effort**: P1W1-E3 - Registry Configuration Schema
**Branch**: `phase1/wave1/P1W1-E3-registry-config`
**Base Branch**: `main`
**Base Branch Reason**: First wave of Phase 1 - no previous integration exists
**Can Parallelize**: Yes
**Parallel With**: P1W1-E1, P1W1-E2, P1W1-E4
**Size Estimate**: 180 lines (MUST be <800)
**Dependencies**: None
**Dependent Efforts**: P1W2-E1, P1W2-E2
**Atomic PR**: ✅ This effort = ONE PR to main (R220 REQUIREMENT)

## 📋 Source Information
**Wave Plan**: PROJECT-IMPLEMENTATION-PLAN.md
**Effort Section**: P1W1-E3
**Created By**: Code Reviewer Agent
**Date**: 2025-09-27
**Extracted**: 2025-09-27T14:36:45Z

## 🔴 BASE BRANCH VALIDATION (R337 MANDATORY)
**The orchestrator-state.json is the SOLE SOURCE OF TRUTH for base branches!**
- Base branch MUST be explicitly specified above: `main`
- Base branch MUST match what's in orchestrator-state.json
- Reason MUST explain why this base: First wave of implementation, no prior integration branches
- Orchestrator MUST record this in state file before creating infrastructure

## 🚀 Parallelization Context
**Can Parallelize**: Yes
**Parallel With**: P1W1-E1 (Provider Interface), P1W1-E2 (OCI Package Format), P1W1-E4 (CLI Contracts)
**Blocking Status**: N/A - This is a foundational effort
**Parallel Group**: Phase 1 Wave 1 Foundation Group
**Orchestrator Guidance**: Can spawn immediately with other P1W1 efforts

## 🚨 EXPLICIT SCOPE DEFINITION (R311 MANDATORY)

### IMPLEMENT EXACTLY (BE SPECIFIC!)

#### Types/Structs to Define (EXACTLY 5)
```go
// Type 1: Main registry configuration
type RegistryConfig struct {
    URL        string            // Registry URL
    Type       string            // Registry type (dockerhub, harbor, gitea, generic)
    Insecure   bool              // Allow insecure connections
    Auth       AuthConfig        // Authentication configuration
    Options    map[string]string // Registry-specific options
}

// Type 2: Authentication configuration
type AuthConfig struct {
    Type        string // auth type: basic, token, oauth2
    Username    string // username for basic auth
    Password    string // password for basic auth (from env/secret)
    Token       string // bearer token
    ConfigFile  string // path to docker config
}

// Type 3: Registry connection pool config
type ConnectionConfig struct {
    MaxIdleConns    int           // Maximum idle connections
    MaxOpenConns    int           // Maximum open connections
    ConnMaxLifetime time.Duration // Connection lifetime
    Timeout         time.Duration // Request timeout
}

// Type 4: TLS configuration
type TLSConfig struct {
    CACert     string // CA certificate path
    ClientCert string // Client certificate path
    ClientKey  string // Client key path
    SkipVerify bool   // Skip TLS verification (dev only)
}

// Type 5: Validation result
type ValidationResult struct {
    Valid   bool     // Whether config is valid
    Errors  []string // List of validation errors
    Warnings []string // List of warnings
}
```

#### Functions to Create (EXACTLY 6 - NO MORE)
```go
1. LoadRegistryConfig(path string) (*RegistryConfig, error)      // ~30 lines - Load from YAML/JSON
2. ValidateRegistryConfig(config *RegistryConfig) *ValidationResult // ~40 lines - Validate config
3. GetAuthConfig(config *RegistryConfig) (*AuthConfig, error)    // ~25 lines - Extract auth
4. SetDefaultValues(config *RegistryConfig)                      // ~20 lines - Apply defaults
5. MergeConfigs(base, override *RegistryConfig) *RegistryConfig  // ~25 lines - Merge configs
6. ToConnectionString(config *RegistryConfig) string             // ~15 lines - Build conn string
// STOP HERE - DO NOT ADD MORE FUNCTIONS
```

### 🛑 DO NOT IMPLEMENT (SCOPE BOUNDARIES)

**EXPLICITLY FORBIDDEN IN THIS EFFORT:**
- ❌ DO NOT implement actual registry connections (P1W2-E1)
- ❌ DO NOT implement authentication flows (P1W2-E2)
- ❌ DO NOT add registry health checks
- ❌ DO NOT implement config hot-reloading
- ❌ DO NOT add config migration or versioning
- ❌ DO NOT create CLI commands (P1W1-E4)
- ❌ DO NOT implement registry discovery
- ❌ DO NOT add comprehensive error types (basic only)
- ❌ DO NOT write integration tests
- ❌ DO NOT add config encryption/decryption

### 📊 REALISTIC SIZE CALCULATION

```
Component Breakdown:
- Type definitions (5 structs):           40 lines
- LoadRegistryConfig function:            30 lines
- ValidateRegistryConfig function:        40 lines
- GetAuthConfig function:                 25 lines
- SetDefaultValues function:              20 lines
- MergeConfigs function:                  25 lines
- ToConnectionString function:            15 lines
- Basic unit tests (6 × 20):            120 lines
- Comments and imports:                   30 lines

TOTAL ESTIMATE: 345 lines (must be <800)
BUFFER: 455 lines for unforeseen needs
```

## 🔴🔴🔴 PRE-PLANNING RESEARCH RESULTS (R374 MANDATORY) 🔴🔴🔴

### Existing Interfaces Found
| Interface | Location | Signature | Must Implement |
|-----------|----------|-----------|----------------|
| N/A - First wave | N/A | N/A | N/A |

### Existing Implementations to Reuse
| Component | Location | Purpose | How to Use |
|-----------|----------|---------|------------|
| N/A - Foundational effort | N/A | N/A | N/A |

### APIs Already Defined
| API | Method | Signature | Notes |
|-----|--------|-----------|-------|
| N/A - No existing APIs | N/A | N/A | First implementation |

### FORBIDDEN DUPLICATIONS (R373)
- ❌ DO NOT create alternative config formats
- ❌ DO NOT reimplement YAML/JSON parsing (use standard library)
- ❌ DO NOT create custom validation framework

### REQUIRED INTEGRATIONS (R373)
- ✅ MUST use standard library for YAML/JSON parsing
- ✅ MUST follow Go conventions for configuration
- ✅ MUST use environment variables for secrets

## 📁 Files to Create

### Primary Implementation Files
```yaml
new_files:
  - path: pkg/config/registry.go
    lines: ~60 MAX
    purpose: Main registry configuration types and loading
    contains:
      - RegistryConfig struct
      - ConnectionConfig struct
      - LoadRegistryConfig function
      - ToConnectionString function

  - path: pkg/config/auth.go
    lines: ~50 MAX
    purpose: Authentication configuration
    contains:
      - AuthConfig struct
      - TLSConfig struct
      - GetAuthConfig function ONLY
      - NO helper functions

  - path: pkg/config/validation.go
    lines: ~70 MAX
    purpose: Configuration validation logic
    contains:
      - ValidationResult struct
      - ValidateRegistryConfig function
      - SetDefaultValues function
      - MergeConfigs function
```

### Test Files
```yaml
test_files:
  - path: pkg/config/registry_test.go
    lines: ~40 MAX
    coverage_target: 80%
    test_functions:
      - TestLoadRegistryConfig  # ~20 lines
      - TestToConnectionString   # ~20 lines

  - path: pkg/config/auth_test.go
    lines: ~40 MAX
    coverage_target: 80%
    test_functions:
      - TestGetAuthConfig        # ~20 lines
      - TestTLSConfig           # ~20 lines

  - path: pkg/config/validation_test.go
    lines: ~40 MAX
    coverage_target: 80%
    test_functions:
      - TestValidateRegistryConfig  # ~20 lines
      - TestMergeConfigs            # ~20 lines
```

## 📦 Files to Import/Reuse

### From Previous Efforts (This Wave)
```yaml
this_wave_imports:
  # None - this is a parallel foundational effort
```

### From Previous Waves/Phases
```yaml
previous_work_imports:
  # None - this is Phase 1 Wave 1
```

## 🔗 Dependencies

### Effort Dependencies
- **Must Complete First**: None - foundational effort
- **Can Run in Parallel With**: P1W1-E1, P1W1-E2, P1W1-E4
- **Blocks**: P1W2-E1 (Base OCI Registry Client), P1W2-E2 (Authentication Handler)

### Technical Dependencies
- Standard library: `encoding/json`, `gopkg.in/yaml.v3`
- No external dependencies in this effort

## 🔴 ATOMIC PR REQUIREMENTS (R220 - SUPREME LAW)

### 🔴🔴🔴 PARAMOUNT: Independent Mergeability (R307) 🔴🔴🔴
**This effort MUST be mergeable at ANY time, even YEARS later:**
- ✅ Must compile when merged alone to main
- ✅ Must NOT break any existing functionality
- ✅ Configuration structs are passive (no active features)
- ✅ Must work even if registry client never implemented
- ✅ Config validation works standalone

### Feature Flags for This Effort
```yaml
feature_flags:
  # Not needed - passive configuration types only
```

### 🚨🚨🚨 R355 PRODUCTION READY CODE (SUPREME LAW #5) 🚨🚨🚨

**ALL CODE MUST BE PRODUCTION READY - NO EXCEPTIONS**

#### ❌ ABSOLUTELY FORBIDDEN:
- NO placeholder validation logic
- NO hardcoded registry URLs
- NO embedded credentials
- NO TODO/FIXME markers
- NO panic for validation errors
- NO stub config loaders

#### ✅ REQUIRED PATTERNS:
```go
// ❌ WRONG - Hardcoded
defaultURL := "docker.io"
defaultTimeout := 30

// ✅ CORRECT - Configuration-driven
defaultURL := os.Getenv("REGISTRY_URL")
if defaultURL == "" {
    defaultURL = "docker.io"
}
defaultTimeout := getIntEnv("REGISTRY_TIMEOUT", 30)
```

### PR Mergeability Checklist
- [ ] PR can merge to main independently
- [ ] Build passes with just this PR
- [ ] All tests pass in isolation
- [ ] No active features (config only)
- [ ] No breaking changes to existing code
- [ ] Backward compatible with main

## 🔴 MANDATORY ADHERENCE CHECKPOINTS (R311)

### Before Starting:
```bash
echo "EFFORT SCOPE LOCKED:"
echo "✓ Functions: EXACTLY 6 (Load, Validate, GetAuth, SetDefaults, Merge, ToString)"
echo "✓ Types: EXACTLY 5 (RegistryConfig, AuthConfig, ConnectionConfig, TLSConfig, ValidationResult)"
echo "✓ Endpoints: EXACTLY 0 (config only)"
echo "✓ Tests: EXACTLY 6 basic tests"
echo "✗ Validation: BASIC ONLY (no deep validation)"
echo "✗ Extra features: NONE"
echo "✗ Optimizations: NONE"
```

### During Implementation:
```bash
# Check scope adherence after each component
FUNC_COUNT=$(grep -c "^func [A-Z]" pkg/config/*.go 2>/dev/null || echo 0)
if [ "$FUNC_COUNT" -gt 6 ]; then
    echo "⚠️ WARNING: Exceeding function count! Stop adding!"
fi
```

## 📝 Implementation Instructions

### Step-by-Step Guide
1. **Create Directory Structure**
   ```bash
   mkdir -p pkg/config
   ```

2. **Implementation Order**
   - Start with `registry.go` - define RegistryConfig and ConnectionConfig structs
   - Create `auth.go` - define AuthConfig and TLSConfig structs
   - Implement `validation.go` - validation and utility functions
   - Write minimal unit tests for each file

3. **Key Implementation Details**
   ```go
   // registry.go - Main config structure
   type RegistryConfig struct {
       URL        string            `yaml:"url" json:"url"`
       Type       string            `yaml:"type" json:"type"`
       Insecure   bool              `yaml:"insecure" json:"insecure"`
       Auth       AuthConfig        `yaml:"auth" json:"auth"`
       Options    map[string]string `yaml:"options" json:"options"`
   }

   func LoadRegistryConfig(path string) (*RegistryConfig, error) {
       // Read file, unmarshal YAML/JSON
       // Apply environment variable overrides
       // Return config
   }
   ```

4. **Configuration Sources**
   - Primary: YAML/JSON config files
   - Override: Environment variables for sensitive data
   - Default: Built-in defaults for optional fields

## ✅ Test Requirements

### Coverage Requirements
- **Minimum Coverage**: 80%
- **Focus**: Config loading and validation
- **Scope**: Unit tests only (no integration tests)

### Test Categories
```yaml
required_tests:
  unit_tests:
    - Config loading from file
    - Validation of required fields
    - Default value application
    - Config merging logic
    - Connection string generation
    - Auth config extraction
```

## 📏 Size Constraints
**Target Size**: 180 lines (from wave plan)
**Maximum Size**: 800 lines (HARD LIMIT)
**Current Estimate**: 345 lines

### Size Monitoring Protocol
```bash
# After creating each file
cd /home/vscode/workspaces/idpbuilder-gitea-push/efforts/phase1/wave1/P1W1-E3-registry-config
find pkg -name "*.go" -not -name "*_test.go" | xargs wc -l

# If approaching 700 lines:
# 1. Alert Code Reviewer
# 2. Focus on core functionality only
# 3. Defer any nice-to-have features
```

## 🏁 Completion Criteria

### Implementation Checklist
- [ ] All 3 implementation files created
- [ ] All 5 struct types defined
- [ ] All 6 functions implemented
- [ ] Size verified under 800 lines
- [ ] No dependencies on other efforts

### Quality Checklist
- [ ] Test coverage ≥80%
- [ ] All tests passing
- [ ] No linting errors
- [ ] Basic validation complete
- [ ] Environment variable support

### Documentation Checklist
- [ ] Struct fields documented
- [ ] Function purposes clear
- [ ] Example usage in comments
- [ ] Configuration format documented

## 🎯 Success Metrics
- Configuration can be loaded from file
- Validation catches basic errors
- Auth config properly extracted
- Defaults applied correctly
- Configs can be merged
- Connection strings generated
- All without any registry connection code

CONTINUE-SOFTWARE-FACTORY=TRUE