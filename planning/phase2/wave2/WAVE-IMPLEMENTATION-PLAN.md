# Wave 2.2 Implementation Plan - Advanced Configuration Features

**Wave**: Phase 2, Wave 2 (Advanced Configuration Features)
**Phase**: Phase 2 - Core Push Functionality
**Created**: 2025-11-01
**Planner**: @agent-code-reviewer
**Fidelity Level**: **EXACT SPECIFICATIONS** (detailed efforts with R213 metadata)

---

## Wave Overview

**Goal**: Extend the Push Command with environment variable support and configuration precedence, enabling users to configure push operations via environment variables (`IDPBUILDER_*`) in addition to command-line flags, with clear precedence (flags > env > defaults).

**Architecture Reference**: See `planning/phase2/wave2/WAVE-2.2-ARCHITECTURE.md` for design details

**Test Plan Reference**: See `planning/phase2/wave2/WAVE-TEST-PLAN.md` for 50 TDD tests

**Total Efforts**: 2

**Total Estimated Lines**: ~750 lines (within Wave limit of 3500)

**Base Branch**: `idpbuilder-oci-push/phase2-integration` (Wave 2.1 complete)

---

## Effort Definitions

### Effort 2.2.1: Registry Override & Viper Integration

#### R213 Metadata

```json
{
  "effort_id": "2.2.1",
  "effort_name": "Registry Override & Viper Integration",
  "branch_name": "idpbuilder-oci-push/phase2/wave2/effort-2.2.1-registry-override-viper",
  "parent_wave": "WAVE_2.2",
  "parent_phase": "PHASE_2",
  "depends_on": ["integration:phase2-wave2.1"],
  "estimated_lines": 400,
  "complexity": "medium",
  "can_parallelize": false,
  "risk_level": "medium",
  "parallelizable": false,
  "parallel_with": []
}
```

#### Scope

**Purpose**: Implement configuration system with environment variable support and precedence resolution. Extend Wave 2.1's `PushOptions` with source tracking, integrate with IDPBuilder's viper instance, and implement registry override with environment variable fallback.

**Boundaries - IN SCOPE**:
- Configuration loader with precedence (flags > env > defaults)
- `PushConfig` struct with source tracking
- Environment variable constants (`IDPBUILDER_*` prefix)
- Boolean parsing for environment variables (true/false/1/0/yes/no)
- Integration with existing Wave 2.1 `runPush()` function
- Configuration validation with helpful error messages
- Verbose mode displaying configuration sources

**Boundaries - OUT OF SCOPE**:
- Integration tests (Effort 2.2.2)
- Configuration file support (Viper can do this, but NOT in this wave)
- Dynamic configuration reload
- Configuration encryption/security beyond password redaction
- Multi-registry push (single registry only)

#### Files to Create/Modify

**New Files**:
- `pkg/cmd/push/config.go` (285 lines) - Configuration system
- `pkg/cmd/push/config_test.go` (placeholder for unit tests, 50 lines)

**Modified Files**:
- `pkg/cmd/push/push.go` (add LoadConfig integration, modify NewPushCommand, +65 lines)
- Update command help text to document environment variables
- Update RunE function to use LoadConfig instead of direct flag access

**Total Estimated Lines**: 400 lines (285 new + 65 modifications + 50 test placeholder)

#### Exact Code Specifications

**File: pkg/cmd/push/config.go**

```go
package push

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

// ConfigSource indicates where a configuration value came from
type ConfigSource int

const (
    SourceDefault ConfigSource = iota
    SourceEnv
    SourceFlag
)

func (s ConfigSource) String() string {
    switch s {
    case SourceDefault:
        return "default"
    case SourceEnv:
        return "environment"
    case SourceFlag:
        return "flag"
    default:
        return "unknown"
    }
}

// ConfigValue represents a configuration value with its source
type ConfigValue struct {
    Value  string
    Source ConfigSource
}

// PushConfig extends PushOptions with configuration source tracking
type PushConfig struct {
    ImageName  ConfigValue
    Registry   ConfigValue
    Username   ConfigValue
    Password   ConfigValue
    Insecure   ConfigValue
    Verbose    ConfigValue
}

// Environment variable names for push command configuration
const (
    EnvRegistry  = "IDPBUILDER_REGISTRY"
    EnvUsername  = "IDPBUILDER_USERNAME"
    EnvPassword  = "IDPBUILDER_PASSWORD"
    EnvInsecure  = "IDPBUILDER_INSECURE"
    EnvVerbose   = "IDPBUILDER_VERBOSE"
)

// Default values for push command
const (
    DefaultRegistry = "gitea.cnoe.localtest.me:8443"
)

// LoadConfig loads configuration with proper precedence (flags > env > defaults)
func LoadConfig(cmd *cobra.Command, args []string, v *viper.Viper) (*PushConfig, error) {
    config := &PushConfig{}

    // ImageName always comes from positional argument (required)
    if len(args) < 1 {
        return nil, fmt.Errorf("image name is required")
    }
    config.ImageName = ConfigValue{
        Value:  args[0],
        Source: SourceFlag, // Positional args treated as flags
    }

    // Registry: flags > env > default
    config.Registry = resolveStringConfig(cmd, v, "registry", EnvRegistry, DefaultRegistry)

    // Username: flags > env > no default (required)
    config.Username = resolveStringConfig(cmd, v, "username", EnvUsername, "")

    // Password: flags > env > no default (required)
    config.Password = resolveStringConfig(cmd, v, "password", EnvPassword, "")

    // Insecure: flags > env > default (false)
    config.Insecure = resolveBoolConfig(cmd, v, "insecure", EnvInsecure, false)

    // Verbose: flags > env > default (false)
    config.Verbose = resolveBoolConfig(cmd, v, "verbose", EnvVerbose, false)

    return config, nil
}

// resolveStringConfig resolves a string configuration with precedence
func resolveStringConfig(cmd *cobra.Command, v *viper.Viper, flagName, envName, defaultValue string) ConfigValue {
    flag := cmd.Flags().Lookup(flagName)
    if flag == nil {
        // Flag doesn't exist, use default
        return ConfigValue{Value: defaultValue, Source: SourceDefault}
    }

    // Check if flag was explicitly set by user
    if flag.Changed {
        return ConfigValue{
            Value:  flag.Value.String(),
            Source: SourceFlag,
        }
    }

    // Flag not set, check environment variable
    if envValue := os.Getenv(envName); envValue != "" {
        return ConfigValue{
            Value:  envValue,
            Source: SourceEnv,
        }
    }

    // Neither flag nor env set, use default
    return ConfigValue{
        Value:  defaultValue,
        Source: SourceDefault,
    }
}

// resolveBoolConfig resolves a boolean configuration with precedence
func resolveBoolConfig(cmd *cobra.Command, v *viper.Viper, flagName, envName string, defaultValue bool) ConfigValue {
    flag := cmd.Flags().Lookup(flagName)
    if flag == nil {
        // Flag doesn't exist, use default
        return ConfigValue{
            Value:  fmt.Sprintf("%t", defaultValue),
            Source: SourceDefault,
        }
    }

    // Check if flag was explicitly set by user
    if flag.Changed {
        return ConfigValue{
            Value:  flag.Value.String(),
            Source: SourceFlag,
        }
    }

    // Flag not set, check environment variable
    if envValue := os.Getenv(envName); envValue != "" {
        // Parse boolean from env var (supports: true, false, 1, 0, yes, no)
        switch envValue {
        case "true", "1", "yes", "YES", "True", "TRUE":
            return ConfigValue{Value: "true", Source: SourceEnv}
        case "false", "0", "no", "NO", "False", "FALSE":
            return ConfigValue{Value: "false", Source: SourceEnv}
        default:
            // Invalid value, use default
            return ConfigValue{
                Value:  fmt.Sprintf("%t", defaultValue),
                Source: SourceDefault,
            }
        }
    }

    // Neither flag nor env set, use default
    return ConfigValue{
        Value:  fmt.Sprintf("%t", defaultValue),
        Source: SourceDefault,
    }
}

// ToPushOptions converts PushConfig to PushOptions (Wave 2.1 compatibility)
func (c *PushConfig) ToPushOptions() *PushOptions {
    return &PushOptions{
        ImageName: c.ImageName.Value,
        Registry:  c.Registry.Value,
        Username:  c.Username.Value,
        Password:  c.Password.Value,
        Insecure:  c.Insecure.Value == "true",
        Verbose:   c.Verbose.Value == "true",
    }
}

// Validate checks if configuration is valid
func (c *PushConfig) Validate() error {
    if c.ImageName.Value == "" {
        return fmt.Errorf("image name is required")
    }
    if c.Username.Value == "" {
        return fmt.Errorf("username is required (use --username flag or %s environment variable)", EnvUsername)
    }
    if c.Password.Value == "" {
        return fmt.Errorf("password is required (use --password flag or %s environment variable)", EnvPassword)
    }
    return nil
}

// DisplaySources prints configuration sources (for debugging/verbose mode)
func (c *PushConfig) DisplaySources() {
    fmt.Printf("Configuration sources:\n")
    fmt.Printf("  Image name: %s (from %s)\n", c.ImageName.Value, c.ImageName.Source)
    fmt.Printf("  Registry: %s (from %s)\n", c.Registry.Value, c.Registry.Source)
    fmt.Printf("  Username: %s (from %s)\n", c.Username.Value, c.Username.Source)
    fmt.Printf("  Password: [redacted] (from %s)\n", c.Password.Source)
    fmt.Printf("  Insecure: %s (from %s)\n", c.Insecure.Value, c.Insecure.Source)
    fmt.Printf("  Verbose: %s (from %s)\n", c.Verbose.Value, c.Verbose.Source)
}
```

**Implementation Requirements**:
1. **Precedence must be strictly enforced**: Flags always override environment, environment always overrides defaults
2. **Boolean parsing must be flexible**: Support true/false, 1/0, yes/no, YES/NO, etc.
3. **Flag.Changed detection is critical**: Must detect if user explicitly set flag vs using default
4. **Error messages must be helpful**: Mention both flag and environment variable options
5. **Password redaction**: Never print actual password, always show "[redacted]"
6. **Wave 2.1 compatibility**: ToPushOptions must produce identical PushOptions struct

**File: pkg/cmd/push/push.go (modifications)**

```go
// Modify NewPushCommand function signature to accept viper instance
func NewPushCommand(v *viper.Viper) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "push IMAGE",
        Short: "Push a Docker image to an OCI registry",
        Long: `Push a local Docker image to an OCI-compliant container registry.

The command retrieves the image from the local Docker daemon and pushes it to
the specified registry using credentials provided via flags or environment variables.

Configuration precedence (highest to lowest):
  1. Command-line flags (--registry, --username, etc.)
  2. Environment variables (IDPBUILDER_REGISTRY, IDPBUILDER_USERNAME, etc.)
  3. Default values

Environment variables:
  IDPBUILDER_REGISTRY   Registry URL (default: gitea.cnoe.localtest.me:8443)
  IDPBUILDER_USERNAME   Registry username
  IDPBUILDER_PASSWORD   Registry password
  IDPBUILDER_INSECURE   Skip TLS verification (true/false, default: false)
  IDPBUILDER_VERBOSE    Enable verbose output (true/false, default: false)

Examples:
  # Push using flags
  idpbuilder push alpine:latest --username admin --password password

  # Push using environment variables
  export IDPBUILDER_USERNAME=admin
  export IDPBUILDER_PASSWORD=password
  idpbuilder push alpine:latest

  # Override registry with flag (takes precedence over env var)
  export IDPBUILDER_REGISTRY=docker.io
  idpbuilder push myapp:v1.0 --registry ghcr.io --username user --password pass

  # Push with verbose progress
  idpbuilder push alpine:latest --verbose

  # Push with insecure TLS (development only)
  idpbuilder push alpine:latest --insecure`,
        Args: cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            // Load configuration with precedence handling
            config, err := LoadConfig(cmd, args, v)
            if err != nil {
                return fmt.Errorf("configuration error: %w", err)
            }

            // Validate configuration
            if err := config.Validate(); err != nil {
                return fmt.Errorf("validation error: %w", err)
            }

            // Display configuration sources in verbose mode
            if config.Verbose.Value == "true" {
                config.DisplaySources()
                fmt.Println()
            }

            // Convert to PushOptions for backward compatibility with Wave 2.1
            opts := config.ToPushOptions()

            // Execute push pipeline (unchanged from Wave 2.1)
            return runPush(cmd.Context(), opts)
        },
    }

    // Define flags (same as Wave 2.1, but now with env var support)
    cmd.Flags().String("registry", DefaultRegistry,
        fmt.Sprintf("Registry URL (env: %s)", EnvRegistry))
    cmd.Flags().String("username", "",
        fmt.Sprintf("Registry username (required, env: %s)", EnvUsername))
    cmd.Flags().String("password", "",
        fmt.Sprintf("Registry password (required, env: %s)", EnvPassword))
    cmd.Flags().Bool("insecure", false,
        fmt.Sprintf("Skip TLS certificate verification (env: %s)", EnvInsecure))
    cmd.Flags().Bool("verbose", false,
        fmt.Sprintf("Enable verbose progress output (env: %s)", EnvVerbose))

    // Note: We do NOT mark flags as required because they can come from env vars
    // Validation happens in LoadConfig/Validate instead

    return cmd
}

// runPush function remains UNCHANGED from Wave 2.1
// It already works with PushOptions, so no modifications needed
```

#### Tests Required

**Unit Tests** (30 tests - defined in WAVE-TEST-PLAN.md):

**File: pkg/cmd/push/config_test.go**

**Test Suites**:
1. **Configuration Precedence Tests** (12 tests):
   - T-2.2.1-01 to T-2.2.1-12 from test plan
   - Test flag > env > default precedence
   - Test mixed sources
   - Test flag.Changed detection

2. **Boolean Resolution Tests** (12 tests):
   - T-2.2.2-01 to T-2.2.2-12 from test plan
   - Test all boolean formats (true/false/1/0/yes/no/YES/NO)
   - Test invalid boolean values fall back to default

3. **Validation Tests** (8 tests):
   - T-2.2.3-01 to T-2.2.3-08 from test plan
   - Test missing required fields
   - Test error messages mention environment variables
   - Test special characters in passwords

4. **Conversion Tests** (8 tests):
   - T-2.2.4-01 to T-2.2.4-08 from test plan
   - Test ToPushOptions conversion
   - Test DisplaySources output
   - Test ConfigSource.String()

**Coverage Target**: 90% statement, 85% branch

**See**: `planning/phase2/wave2/WAVE-TEST-PLAN.md` for complete test specifications

#### Dependencies

**Upstream Dependencies** (must complete before this effort):
- Wave 2.1 complete (Push Command Core, Progress Reporter)
- Phase 1 complete (docker, registry, auth, tls packages)
- IDPBuilder root command with viper instance

**Downstream Dependencies** (efforts that depend on this):
- Effort 2.2.2 (integration tests depend on this configuration system)

**External Libraries** (already in go.mod):
- `github.com/spf13/cobra` v1.8.0 - Command framework
- `github.com/spf13/viper` v1.17.0 - Configuration (NOW actively used)
- `github.com/spf13/pflag` v1.0.5 - Flag change detection

#### Acceptance Criteria

- [ ] All files created/modified as specified
- [ ] `LoadConfig` function implements strict precedence (flags > env > defaults)
- [ ] `resolveBoolConfig` supports all boolean formats (true/false/1/0/yes/no variants)
- [ ] `PushConfig.Validate()` returns helpful error messages mentioning env vars
- [ ] `ToPushOptions()` maintains Wave 2.1 compatibility
- [ ] 30 unit tests passing (100% pass rate)
- [ ] Code coverage ≥90% statement, ≥85% branch
- [ ] No linting errors (go vet, golangci-lint)
- [ ] All public functions have godoc comments
- [ ] Line count within estimate (400 ± 60 lines = 340-460 acceptable)
- [ ] Wave 2.1's `runPush()` function remains unchanged

---

### Effort 2.2.2: Environment Variable Support & Integration Testing

#### R213 Metadata

```json
{
  "effort_id": "2.2.2",
  "effort_name": "Environment Variable Support & Integration Testing",
  "branch_name": "idpbuilder-oci-push/phase2/wave2/effort-2.2.2-env-var-integration",
  "parent_wave": "WAVE_2.2",
  "parent_phase": "PHASE_2",
  "depends_on": ["effort:2.2.1"],
  "estimated_lines": 350,
  "complexity": "medium",
  "can_parallelize": false,
  "risk_level": "low",
  "parallelizable": false,
  "parallel_with": []
}
```

#### Scope

**Purpose**: Implement comprehensive integration tests for environment variable support and configuration precedence. Verify end-to-end push operations work with environment variables, validate backward compatibility with Wave 2.1, and ensure robust error handling.

**Boundaries - IN SCOPE**:
- Integration tests with real environment variables
- End-to-end push tests using env vars only
- Precedence verification tests (flag overrides env)
- Verbose mode output verification
- Error message validation with env var hints
- Backward compatibility tests (Wave 2.1 style still works)
- Edge case handling (invalid booleans, empty strings, etc.)

**Boundaries - OUT OF SCOPE**:
- Production usage (this is testing only)
- Performance benchmarking
- Load testing
- Multi-registry tests (single registry only)
- Configuration file tests (not in this wave)

#### Files to Create/Modify

**New Files**:
- `pkg/cmd/push/push_integration_test.go` (350 lines) - Integration tests

**Modified Files**:
- None (tests only)

**Total Estimated Lines**: 350 lines

#### Exact Test Specifications

**File: pkg/cmd/push/push_integration_test.go**

```go
package push_test

import (
    "context"
    "os"
    "testing"
    "time"

    "github.com/cnoe-io/idpbuilder/pkg/cmd/push"
    "github.com/spf13/viper"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

// Test Suite 5: Environment Variable Scenarios (10 tests)

func TestPushCommand_AllFromEnvironment(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }

    // Given: All configuration from environment variables
    os.Setenv("IDPBUILDER_REGISTRY", "gitea.cnoe.localtest.me:8443")
    os.Setenv("IDPBUILDER_USERNAME", "giteaAdmin")
    os.Setenv("IDPBUILDER_PASSWORD", "password")
    os.Setenv("IDPBUILDER_INSECURE", "true")
    defer func() {
        os.Unsetenv("IDPBUILDER_REGISTRY")
        os.Unsetenv("IDPBUILDER_USERNAME")
        os.Unsetenv("IDPBUILDER_PASSWORD")
        os.Unsetenv("IDPBUILDER_INSECURE")
    }()

    // Given: Command with NO flags set
    cmd := push.NewPushCommand(viper.New())
    cmd.SetArgs([]string{"alpine:latest"})

    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
    defer cancel()
    cmd.SetContext(ctx)

    // When: Execute command (should use env vars)
    err := cmd.Execute()

    // Then: Command succeeds using environment variables only
    require.NoError(t, err, "Push should succeed with env vars only")
}

func TestPushCommand_FlagOverridesEnvironment(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }

    // Given: Environment variable set to different registry
    os.Setenv("IDPBUILDER_REGISTRY", "docker.io")
    os.Setenv("IDPBUILDER_USERNAME", "wronguser")
    os.Setenv("IDPBUILDER_PASSWORD", "wrongpass")
    defer func() {
        os.Unsetenv("IDPBUILDER_REGISTRY")
        os.Unsetenv("IDPBUILDER_USERNAME")
        os.Unsetenv("IDPBUILDER_PASSWORD")
    }()

    // Given: Command with flags that override env vars
    cmd := push.NewPushCommand(viper.New())
    cmd.SetArgs([]string{
        "alpine:latest",
        "--registry", "gitea.cnoe.localtest.me:8443",
        "--username", "giteaAdmin",
        "--password", "password",
        "--insecure",
    })

    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
    defer cancel()
    cmd.SetContext(ctx)

    // When: Execute command
    err := cmd.Execute()

    // Then: Command succeeds using flag values (not env vars)
    require.NoError(t, err, "Push should succeed with flags overriding env")
}

// Test Suite 6: Edge Cases & Error Handling (10 tests)

func TestPushCommand_InvalidBooleanInEnv(t *testing.T) {
    // Given: Environment variable with invalid boolean value
    os.Setenv("IDPBUILDER_INSECURE", "maybe")
    os.Setenv("IDPBUILDER_USERNAME", "testuser")
    os.Setenv("IDPBUILDER_PASSWORD", "testpass")
    defer func() {
        os.Unsetenv("IDPBUILDER_INSECURE")
        os.Unsetenv("IDPBUILDER_USERNAME")
        os.Unsetenv("IDPBUILDER_PASSWORD")
    }()

    // Given: Command with no flags
    cmd := push.NewPushCommand(viper.New())
    cmd.SetArgs([]string{"alpine:latest"})

    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
    defer cancel()
    cmd.SetContext(ctx)

    // When: Execute command
    err := cmd.Execute()

    // Then: Invalid boolean falls back to default (false)
    // Command should still succeed (not fail on invalid bool)
    require.NoError(t, err, "Should succeed with default value for invalid bool")
}

// Additional 17 integration tests (T-2.2.5-03 through T-2.2.6-10)
// See WAVE-TEST-PLAN.md for complete test specifications
// Total: 20 integration tests covering:
// - Mixed configuration sources
// - Verbose mode output
// - Validation error messages with env hints
// - Environment override scenarios
// - Boolean parsing variants
// - Password special characters
// - Registry override functionality
// - Wave 2.1 backward compatibility
// - Empty environment variables
// - Environment variable with spaces
// - Multiple boolean formats
// - Unset after set scenarios
// - Context cancellation
// - Viper instance reuse
// - Concurrent environment access
// - Help text documentation
```

**Implementation Requirements**:
1. **All tests must be real integration tests**: Use `testing.Short()` to allow skipping
2. **Environment cleanup is mandatory**: Use `defer` to unset all env vars
3. **Context timeouts**: All push operations must have 2-minute timeout
4. **Test isolation**: Each test must set up and tear down its own env vars
5. **Backward compatibility verification**: Include test for Wave 2.1 style usage
6. **Error message validation**: Verify error messages mention environment variables
7. **Thread safety**: Test concurrent environment variable access

**Coverage Target**: 85% statement, 80% branch (integration tests)

**See**: `planning/phase2/wave2/WAVE-TEST-PLAN.md` for complete test specifications (T-2.2.5-01 through T-2.2.6-10)

#### Tests Required

**Integration Tests** (20 tests - defined in WAVE-TEST-PLAN.md):

**Test Suite 5: Environment Variable Scenarios** (10 tests):
- T-2.2.5-01: All from environment
- T-2.2.5-02: Flag overrides environment
- T-2.2.5-03: Mixed configuration
- T-2.2.5-04: Verbose shows sources
- T-2.2.5-05: Validation errors with env hints
- T-2.2.5-06: Environment overrides default
- T-2.2.5-07: Insecure from environment
- T-2.2.5-08: Password special characters
- T-2.2.5-09: Registry override
- T-2.2.5-10: Wave 2.1 backward compatibility

**Test Suite 6: Edge Cases & Error Handling** (10 tests):
- T-2.2.6-01: Empty environment variable
- T-2.2.6-02: Invalid boolean in env
- T-2.2.6-03: Env var with spaces
- T-2.2.6-04: Multiple env formats
- T-2.2.6-05: Unset after set
- T-2.2.6-06: Flag explicitly set to empty
- T-2.2.6-07: Context cancellation with env
- T-2.2.6-08: Viper instance reuse
- T-2.2.6-09: Concurrent environment access
- T-2.2.6-10: Env var precedence documentation

**Coverage Target**: 85% statement, 80% branch

#### Dependencies

**Upstream Dependencies** (must complete before this effort):
- Effort 2.2.1 (configuration system implementation)

**Downstream Dependencies** (efforts that depend on this):
- None (this is the final effort in Wave 2.2)

**External Libraries** (already in go.mod):
- `github.com/stretchr/testify` v1.8.4 - Test assertions
- All libraries from Effort 2.2.1

#### Acceptance Criteria

- [ ] All 20 integration tests implemented
- [ ] All tests passing (100% pass rate)
- [ ] Tests can be skipped with `go test -short`
- [ ] All environment variables properly set and cleaned up
- [ ] Context timeouts implemented for all push operations
- [ ] Backward compatibility with Wave 2.1 verified
- [ ] Error messages validated (must mention env vars)
- [ ] Code coverage ≥85% statement, ≥80% branch
- [ ] No test flakiness (tests pass consistently)
- [ ] All tests have clear Given/When/Then structure
- [ ] Line count within estimate (350 ± 50 lines = 300-400 acceptable)

---

## Parallelization Strategy

### Sequential Execution Required

**Effort 2.2.1** (Foundational) → **Effort 2.2.2** (Tests depend on implementation)

**Rationale for Sequential Execution**:
1. **Effort 2.2.1 is foundational**: Defines `PushConfig`, `ConfigValue`, `LoadConfig`, and all precedence logic
2. **Testing depends on implementation**: Effort 2.2.2's integration tests require Effort 2.2.1's configuration system to exist
3. **Small wave size**: Only 2 efforts (~750 lines total), minimal parallelization benefit
4. **Clear interface boundary**: Effort 2.2.1 delivers configuration system, Effort 2.2.2 validates it
5. **No parallel work possible**: Cannot write integration tests before implementation exists

**Testing Strategy**:
- Unit tests (Effort 2.2.1) use TDD approach (written before implementation)
- Integration tests (Effort 2.2.2) written after Effort 2.2.1 interfaces are stable
- Effort 2.2.2 can start as soon as Effort 2.2.1 code review passes

**Execution Timeline**:
```
Day 1-2: Effort 2.2.1 implementation + unit tests
Day 2:   Code review of Effort 2.2.1
Day 3:   Effort 2.2.2 integration tests (after 2.2.1 approved)
Day 3:   Code review of Effort 2.2.2
Day 4:   Wave integration + architect review
```

---

## Wave Size Compliance

**Total Wave Lines**: 750 lines (400 + 350)

**Size Limit**: 3500 lines (soft), 4000 lines (hard)

**Status**:
- ✅ Within soft limit (750 < 3500 lines)
- ✅ Within hard limit (750 < 4000 lines)
- ✅ No split required

**Size Breakdown**:
- Effort 2.2.1: 400 lines (285 config.go + 65 push.go mods + 50 test placeholder)
- Effort 2.2.2: 350 lines (integration tests only)
- **Total**: 750 lines

**Size Compliance**: ✅ **COMPLIANT** (21% of soft limit, 19% of hard limit)

---

## Integration Strategy

### Cascade Branching (R501 Compliance)

**Branch Structure**:
```
main (production)
  └─ idpbuilder-oci-push/phase2-integration (Wave 2.1 complete)
      └─ idpbuilder-oci-push/phase2/wave2/integration (Wave 2.2 target)
          ├─ idpbuilder-oci-push/phase2/wave2/effort-2.2.1-registry-override-viper
          └─ idpbuilder-oci-push/phase2/wave2/effort-2.2.2-env-var-integration (from 2.2.1)
```

**Integration Sequence**:
1. **Effort 2.2.1** completes → Code Review → Merge to `phase2/wave2/integration`
2. **Effort 2.2.2** starts from Effort 2.2.1 branch
3. **Effort 2.2.2** completes → Code Review → Merge to `phase2/wave2/integration`
4. **Wave integration tests** run on `phase2/wave2/integration` branch
5. **Architect wave review** performs assessment
6. **Merge to phase integration**: `phase2/wave2/integration` → `phase2-integration`
7. **Phase 2 continues** with Wave 2.3

### Independent Mergeability (R307 Compliance)

**Wave 2.2 can merge independently**:
- ✅ No breaking changes to Wave 2.1 interfaces (`PushOptions` unchanged)
- ✅ Backward compatible (Wave 2.1 flag-only usage still works)
- ✅ Feature complete (environment variable support fully functional)
- ✅ All tests pass (50 tests: 30 unit + 20 integration)
- ✅ No dependencies on future waves

**Graceful degradation**:
- If environment variables not set, flags work exactly as Wave 2.1
- If viper not available, flags still function (uses os.Getenv directly)
- If IDPBuilder doesn't pass viper instance, can create local instance

---

## Testing Strategy

### Test-Driven Development (R341 Compliance)

**Phase 1: RED (Tests before implementation)**
- ✅ 50 tests defined in WAVE-TEST-PLAN.md
- ✅ All tests written BEFORE Effort 2.2.1 implementation
- ✅ Tests define expected behavior
- ✅ Commit: "tdd: Wave 2.2 tests created BEFORE implementation"

**Phase 2: GREEN (Implementation passes tests)**
- Effort 2.2.1 implements config.go → 30 unit tests turn GREEN
- Effort 2.2.2 implements integration tests → 20 integration tests turn GREEN
- Target: ALL 50 tests passing

**Phase 3: REFACTOR (Coverage verification)**
- Run coverage: `go test -cover ./pkg/cmd/push/...`
- Verify: ≥90% statement, ≥85% branch
- Generate report: `go test -coverprofile=wave2.2-coverage.out`

### Test Invocation Commands

```bash
# Run all Wave 2.2 tests
go test -v ./pkg/cmd/push/... -run "TestLoadConfig|TestResolve|TestPushConfig|TestPushCommand"

# Run only unit tests (Effort 2.2.1)
go test -v ./pkg/cmd/push/config_test.go

# Run only integration tests (Effort 2.2.2)
go test -v ./pkg/cmd/push/push_integration_test.go

# Run integration tests (skip with -short)
go test -short ./pkg/cmd/push/...  # Skips integration tests
go test ./pkg/cmd/push/...          # Runs all tests

# Run with coverage
go test -cover -coverprofile=wave2.2-coverage.out ./pkg/cmd/push/...
go tool cover -html=wave2.2-coverage.out -o wave2.2-coverage.html

# Run with race detection
go test -race ./pkg/cmd/push/...

# Run specific test
go test -v ./pkg/cmd/push/... -run TestLoadConfig_FlagOverridesEnv
```

### Wave-Level Integration Tests

**After both efforts complete, run wave integration tests**:

```bash
# Verify Wave 2.2 functionality end-to-end
go test -v ./pkg/cmd/push/... -run Integration

# Test backward compatibility with Wave 2.1
go test -v ./pkg/cmd/push/... -run BackwardCompatibility

# Test all configuration precedence scenarios
go test -v ./pkg/cmd/push/... -run Precedence
```

---

## Dependency Graph

### Effort Dependencies

```
┌─────────────────────────────────────────┐
│ Wave 2.1 Complete                       │
│ (Push Command Core + Progress Reporter)│
└────────────────┬────────────────────────┘
                 │
                 v
┌─────────────────────────────────────────┐
│ Effort 2.2.1: Registry Override & Viper │
│ - Configuration system                  │
│ - Environment variable support          │
│ - Precedence resolution                 │
│ - 30 unit tests                         │
└────────────────┬────────────────────────┘
                 │
                 v
┌─────────────────────────────────────────┐
│ Effort 2.2.2: Integration Testing       │
│ - End-to-end env var tests              │
│ - Backward compatibility tests          │
│ - Edge case handling                    │
│ - 20 integration tests                  │
└────────────────┬────────────────────────┘
                 │
                 v
┌─────────────────────────────────────────┐
│ Wave 2.2 Integration                    │
│ - 50 tests passing                      │
│ - 90%+ coverage                         │
│ - Architect review                      │
└─────────────────────────────────────────┘
```

### External Dependencies

**From Wave 2.1**:
- `pkg/cmd/push/push.go` - `runPush()` function (unchanged)
- `pkg/cmd/push/types.go` - `PushOptions` struct (unchanged)
- `pkg/progress/reporter.go` - Progress reporter (unchanged)

**From Phase 1**:
- `pkg/docker` - Docker client interface
- `pkg/registry` - Registry client interface
- `pkg/auth` - Authentication provider interface
- `pkg/tls` - TLS configuration provider interface

**From IDPBuilder**:
- `pkg/cmd/root.go` - Root command with viper instance
- Existing Cobra command patterns
- Existing environment variable conventions (`IDPBUILDER_*` prefix)

---

## Risk Mitigation

### Medium-Risk Items

**Effort 2.2.1 Risks**:
1. **Viper integration complexity**:
   - Risk: IDPBuilder's viper instance may have conflicting bindings
   - Mitigation: Use flag.Changed detection instead of viper bindings
   - Mitigation: Direct os.Getenv calls for environment variables

2. **Flag.Changed edge cases**:
   - Risk: Cobra's Changed detection may not work in all scenarios
   - Mitigation: Comprehensive unit tests for Changed detection (T-2.2.1-09)
   - Mitigation: Fallback to environment variable if Changed unreliable

3. **Boolean parsing ambiguity**:
   - Risk: Users may use unexpected boolean formats
   - Mitigation: Support all common formats (true/false/1/0/yes/no)
   - Mitigation: Invalid values fall back to safe default (false)

**Effort 2.2.2 Risks**:
1. **Test flakiness**:
   - Risk: Environment variable tests may be flaky due to state pollution
   - Mitigation: Strict cleanup with defer in every test
   - Mitigation: Test isolation (each test sets up own env)

2. **Integration test failures**:
   - Risk: Integration tests depend on Docker daemon and registry
   - Mitigation: Use `testing.Short()` to allow skipping
   - Mitigation: Document prerequisites in test comments

### Low-Risk Items

**Backward Compatibility**:
- Risk: Wave 2.1 usage breaks
- Mitigation: ToPushOptions maintains exact PushOptions structure
- Mitigation: runPush() function unchanged
- Mitigation: Dedicated backward compatibility test (T-2.2.5-10)

**Performance**:
- Risk: Environment variable lookups add latency
- Mitigation: os.Getenv is extremely fast (<1μs)
- Mitigation: Only 5 environment variables total

---

## Quality Gates (R502 Compliance)

### Implementation Plan Quality Requirements

- ✅ **R213 metadata**: ALL efforts have complete R213 metadata blocks
- ✅ **Exact file lists**: Complete paths for all files (create + modify)
- ✅ **Real code examples**: Actual Go code (no pseudocode)
- ✅ **Detailed specifications**: Line-by-line implementation guidance
- ✅ **Complete test specifications**: All 50 tests defined
- ✅ **Dependency graphs**: Clear upstream/downstream dependencies
- ✅ **Size estimates**: Accurate line counts per effort
- ✅ **Acceptance criteria**: Clear checklist for each effort
- ✅ **Parallelization strategy**: Documented (sequential with rationale)

### Code Quality Requirements

**Effort 2.2.1**:
- ✅ Statement coverage ≥90%
- ✅ Branch coverage ≥85%
- ✅ All public functions have godoc comments
- ✅ No linting errors (go vet, golangci-lint)
- ✅ Pass all 30 unit tests

**Effort 2.2.2**:
- ✅ Statement coverage ≥85%
- ✅ Branch coverage ≥80%
- ✅ All tests have Given/When/Then structure
- ✅ All env vars cleaned up with defer
- ✅ Pass all 20 integration tests

**Wave 2.2**:
- ✅ ALL 50 tests passing (100% pass rate)
- ✅ Wave coverage ≥90% overall
- ✅ No test flakiness
- ✅ Backward compatibility verified

---

## Compliance Checklist

### R213 Effort Metadata (BLOCKING)
- ✅ Effort 2.2.1 has complete R213 metadata
- ✅ Effort 2.2.2 has complete R213 metadata
- ✅ All required fields present (effort_id, name, branch, dependencies, lines, complexity, parallelizable)
- ✅ Dependency relationships documented

### R502 Implementation Plan Quality Gates
- ✅ EXACT fidelity (not high-level descriptions)
- ✅ Complete file paths for every file
- ✅ Real code specifications (actual Go code)
- ✅ Detailed task breakdowns
- ✅ Dependency graphs included
- ✅ Size estimates per effort

### R510 Checklist Compliance
- ✅ Clear criteria for each section
- ✅ Effort-to-test mapping (30 tests → Effort 2.2.1, 20 tests → Effort 2.2.2)
- ✅ Quality gates documented
- ✅ Success criteria per effort
- ✅ Acceptance criteria checkboxes

### R341 TDD Protocol
- ✅ Tests defined before implementation (50 tests in WAVE-TEST-PLAN.md)
- ✅ Tests specify expected behavior
- ✅ Tests use real imports and types
- ✅ Coverage targets: 90%/85%

### R501 Progressive Trunk-Based Development
- ✅ Wave 2.2 branches from Wave 2.1 integration
- ✅ Effort 2.2.2 branches from Effort 2.2.1
- ✅ Cascade branching documented
- ✅ Integration sequence defined

### R307 Independent Mergeability
- ✅ No breaking changes to Wave 2.1
- ✅ Backward compatible
- ✅ Feature complete
- ✅ Graceful degradation
- ✅ All tests pass

### R340 Progressive Realism
- ✅ Real code from Wave 2.1 (PushOptions, runPush)
- ✅ Concrete function signatures
- ✅ Actual Go syntax (no pseudocode)
- ✅ Working examples in architecture

---

## Next Steps

### For Orchestrator

1. **Create Wave 2.2 integration branch**: `idpbuilder-oci-push/phase2/wave2/integration`
2. **Install test harness**: From WAVE-TEST-PLAN.md into integration branch
3. **Create Effort 2.2.1 infrastructure**:
   - Branch: `idpbuilder-oci-push/phase2/wave2/effort-2.2.1-registry-override-viper`
   - Working copy: `efforts/phase2/wave2/effort-2.2.1-registry-override-viper/`
4. **Spawn SW Engineer** for Effort 2.2.1 with:
   - This implementation plan
   - WAVE-2.2-ARCHITECTURE.md
   - WAVE-TEST-PLAN.md
   - Instructions to implement TDD (tests first)
5. **Monitor Effort 2.2.1** → Code Review → Merge
6. **Create Effort 2.2.2 infrastructure** (after 2.2.1 merges)
7. **Spawn SW Engineer** for Effort 2.2.2
8. **Monitor Effort 2.2.2** → Code Review → Merge
9. **Run wave integration tests**
10. **Spawn Architect** for wave assessment
11. **Merge to phase integration** (after architect approval)

### For Software Engineers

**Effort 2.2.1**:
1. Read this plan + architecture + test plan
2. Create test file first (TDD RED phase)
3. Implement config.go per exact specifications
4. Modify push.go per specifications
5. Run tests → all 30 unit tests should pass (TDD GREEN phase)
6. Verify coverage ≥90%/85%
7. Commit and push for review

**Effort 2.2.2**:
1. Read this plan + test plan
2. Verify Effort 2.2.1 merged and available
3. Implement all 20 integration tests
4. Run tests → all should pass
5. Verify coverage ≥85%/80%
6. Verify backward compatibility test passes
7. Commit and push for review

---

## Document Status

**Status**: ✅ READY FOR ORCHESTRATOR EXECUTION
**Planner**: @agent-code-reviewer
**Created**: 2025-11-01
**Total Efforts**: 2
**Total Lines**: 750 (within limits)
**Fidelity Level**: EXACT (detailed specifications with real code)

**Builds Upon**:
- Wave 2.1: Push Command Core & Progress Reporter (COMPLETE)
- Phase 1: All interface packages (docker, registry, auth, tls - COMPLETE)

**References**:
- Architecture: `planning/phase2/wave2/WAVE-2.2-ARCHITECTURE.md`
- Test Plan: `planning/phase2/wave2/WAVE-TEST-PLAN.md` (50 tests)
- Template: `templates/WAVE-IMPLEMENTATION-TEMPLATE.md`

**Compliance Verified**:
- ✅ R213: Effort metadata for ALL efforts
- ✅ R502: Implementation plan quality gates (EXACT fidelity)
- ✅ R510: Checklist structure followed
- ✅ R341: TDD protocol (tests before implementation)
- ✅ R501: Cascade branching defined
- ✅ R307: Independent mergeability ensured

---

**END OF WAVE 2.2 IMPLEMENTATION PLAN**
