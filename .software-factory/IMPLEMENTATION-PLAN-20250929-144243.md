# E1.2.1 Command Structure Implementation Plan

## 🚨 CRITICAL EFFORT METADATA (FROM WAVE PLAN)
**Branch**: `phase1/wave2/command-structure`
**Can Parallelize**: Yes
**Parallel With**: [E1.2.2, E1.2.3]
**Size Estimate**: 450 lines
**Dependencies**: [E1.1.1, E1.1.2, E1.1.3]

## 🎯 Effort Infrastructure Metadata
**EFFORT_NAME**: E1.2.1-command-structure
**EFFORT_TYPE**: implementation
**EFFORT_PHASE**: 1
**EFFORT_WAVE**: 2
**BASE_BRANCH**: phase1-integration
**CASCADE_FROM**: phase1-integration

## Overview
- **Effort**: Implement complete push command structure with Cobra, including all flags and environment variable support
- **Phase**: 1, Wave: 2
- **Estimated Size**: 450 lines (well within 800 line limit)
- **Implementation Time**: 3 hours

## Dependencies from Wave 1
Based on analysis of Wave 1 efforts:
- **E1.1.1**: Provided idpbuilder structure analysis and patterns
  - Command registration in `pkg/cmd/root.go`
  - Helper utilities in `pkg/cmd/helpers/`
  - Established patterns for Cobra commands
- **E1.1.2**: Set up unit test framework
  - Mock registry client patterns
  - Test utilities in `pkg/testutils/`
- **E1.1.3**: Created integration test infrastructure
  - E2E test patterns
  - Gitea registry test setup

## File Structure
```
pkg/
├── cmd/
│   ├── push/
│   │   ├── push.go          (120 lines - main command implementation)
│   │   ├── flags.go         (80 lines - flag definitions and parsing)
│   │   ├── validation.go    (70 lines - input validation logic)
│   │   └── push_test.go     (150 lines - unit tests)
│   └── helpers/
│       └── env.go            (30 lines - environment variable utilities)
```

## Implementation Steps

### Step 1: Create Enhanced Push Command Structure (push.go)
```go
package push

import (
    "context"
    "fmt"
    "os"

    "github.com/cnoe-io/idpbuilder/pkg/cmd/helpers"
    "github.com/spf13/cobra"
)

// PushOptions holds all configuration for the push command
type PushOptions struct {
    // Registry configuration
    RegistryURL  string
    Repository   string
    Tag          string

    // Authentication
    Username     string
    Password     string

    // TLS configuration
    Insecure     bool

    // Image source
    ImagePath    string
    ImageRef     string

    // Behavior
    DryRun       bool
    Verbose      bool
}

var PushCmd = &cobra.Command{
    Use:   "push [IMAGE] [REGISTRY_URL]",
    Short: "Push OCI artifacts to a registry",
    Long: `Push OCI artifacts to any OCI-compliant registry.

The push command supports authentication via flags or environment variables:
  - Authentication: --username/--password or REGISTRY_USERNAME/REGISTRY_PASSWORD
  - TLS: --insecure flag for self-signed certificates

Examples:
  # Push with authentication via flags
  idpbuilder push myimage:latest registry.example.com/repo --username user --password pass

  # Push with environment variables
  export REGISTRY_USERNAME=user
  export REGISTRY_PASSWORD=pass
  idpbuilder push myimage:latest registry.example.com/repo

  # Push to insecure registry (self-signed cert)
  idpbuilder push myimage:latest registry.example.com/repo --insecure`,
    Args:         cobra.RangeArgs(1, 2),
    RunE:         runPush,
    SilenceUsage: true,
}

func runPush(cmd *cobra.Command, args []string) error {
    opts, err := buildPushOptions(cmd, args)
    if err != nil {
        return fmt.Errorf("invalid options: %w", err)
    }

    // Validate options
    if err := validatePushOptions(opts); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }

    // Log configuration
    helpers.CmdLogger.Info("Starting push operation",
        "image", opts.ImageRef,
        "registry", opts.RegistryURL,
        "insecure", opts.Insecure,
        "dry-run", opts.DryRun,
    )

    // Warning for insecure mode
    if opts.Insecure {
        fmt.Fprintf(os.Stderr, "⚠️  WARNING: TLS certificate verification disabled\n")
        fmt.Fprintf(os.Stderr, "   Only use with self-signed certificates in development\n\n")
    }

    // TODO: In E1.2.3, implement actual push logic here
    // For now, just log the configuration
    helpers.CmdLogger.Info("Push configuration validated successfully")

    if opts.DryRun {
        fmt.Printf("DRY RUN: Would push %s to %s\n", opts.ImageRef, opts.RegistryURL)
        return nil
    }

    fmt.Printf("✅ Push command structure ready (implementation pending E1.2.3)\n")
    return nil
}
```

### Step 2: Implement Flag Definitions (flags.go)
```go
package push

import (
    "os"
    "github.com/spf13/cobra"
)

var (
    // Flag variables
    username     string
    password     string
    insecure     bool
    dryRun       bool
    verbose      bool
    repository   string
    tag          string
)

const (
    // Flag usage strings
    usernameUsage    = "Registry username (env: REGISTRY_USERNAME)"
    passwordUsage    = "Registry password (env: REGISTRY_PASSWORD)"
    insecureUsage    = "Skip TLS certificate verification (use for self-signed certificates)"
    dryRunUsage      = "Perform a dry run without actually pushing"
    verboseUsage     = "Enable verbose output"
    repositoryUsage  = "Override target repository name"
    tagUsage         = "Override image tag"

    // Environment variable names
    EnvRegistryUsername = "REGISTRY_USERNAME"
    EnvRegistryPassword = "REGISTRY_PASSWORD"
    EnvRegistryInsecure = "REGISTRY_INSECURE"
)

func init() {
    // Authentication flags
    PushCmd.Flags().StringVarP(&username, "username", "u", "", usernameUsage)
    PushCmd.Flags().StringVarP(&password, "password", "p", "", passwordUsage)

    // TLS configuration
    PushCmd.Flags().BoolVar(&insecure, "insecure", false, insecureUsage)

    // Behavior flags
    PushCmd.Flags().BoolVar(&dryRun, "dry-run", false, dryRunUsage)
    PushCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, verboseUsage)

    // Repository overrides
    PushCmd.Flags().StringVar(&repository, "repository", "", repositoryUsage)
    PushCmd.Flags().StringVar(&tag, "tag", "", tagUsage)

    // Mark password as hidden in help
    PushCmd.Flags().MarkHidden("password")
}

// buildPushOptions constructs PushOptions from flags and environment
func buildPushOptions(cmd *cobra.Command, args []string) (*PushOptions, error) {
    opts := &PushOptions{
        DryRun:   dryRun,
        Verbose:  verbose,
        Insecure: insecure,
    }

    // Parse positional arguments
    if len(args) > 0 {
        opts.ImageRef = args[0]
    }
    if len(args) > 1 {
        opts.RegistryURL = args[1]
    }

    // Authentication: flags override environment variables
    opts.Username = getStringValue(username, os.Getenv(EnvRegistryUsername))
    opts.Password = getStringValue(password, os.Getenv(EnvRegistryPassword))

    // Check insecure from environment if not set via flag
    if !opts.Insecure && os.Getenv(EnvRegistryInsecure) == "true" {
        opts.Insecure = true
    }

    // Repository overrides
    if repository != "" {
        opts.Repository = repository
    }
    if tag != "" {
        opts.Tag = tag
    }

    return opts, nil
}

// getStringValue returns flag value if set, otherwise env value
func getStringValue(flagValue, envValue string) string {
    if flagValue != "" {
        return flagValue
    }
    return envValue
}
```

### Step 3: Implement Validation Logic (validation.go)
```go
package push

import (
    "fmt"
    "net/url"
    "strings"
)

// validatePushOptions validates all push command options
func validatePushOptions(opts *PushOptions) error {
    // Validate required fields
    if opts.ImageRef == "" {
        return fmt.Errorf("image reference is required")
    }

    if opts.RegistryURL == "" {
        return fmt.Errorf("registry URL is required")
    }

    // Validate registry URL format
    if err := validateRegistryURL(opts.RegistryURL); err != nil {
        return fmt.Errorf("invalid registry URL: %w", err)
    }

    // Validate image reference format
    if err := validateImageRef(opts.ImageRef); err != nil {
        return fmt.Errorf("invalid image reference: %w", err)
    }

    // Validate authentication if registry requires it
    if requiresAuth(opts.RegistryURL) {
        if opts.Username == "" || opts.Password == "" {
            return fmt.Errorf("authentication required: provide --username and --password or set REGISTRY_USERNAME and REGISTRY_PASSWORD")
        }
    }

    return nil
}

// validateRegistryURL checks if the registry URL is valid
func validateRegistryURL(registryURL string) error {
    // Add scheme if missing
    if !strings.Contains(registryURL, "://") {
        registryURL = "https://" + registryURL
    }

    u, err := url.Parse(registryURL)
    if err != nil {
        return err
    }

    if u.Host == "" {
        return fmt.Errorf("registry host cannot be empty")
    }

    // Check for valid schemes
    if u.Scheme != "http" && u.Scheme != "https" {
        return fmt.Errorf("unsupported scheme: %s (use http or https)", u.Scheme)
    }

    return nil
}

// validateImageRef validates the image reference format
func validateImageRef(imageRef string) error {
    if imageRef == "" {
        return fmt.Errorf("image reference cannot be empty")
    }

    // Basic validation - more comprehensive validation will be in E1.2.3
    // Check for invalid characters
    if strings.ContainsAny(imageRef, " \t\n") {
        return fmt.Errorf("image reference contains whitespace")
    }

    return nil
}

// requiresAuth checks if the registry requires authentication
func requiresAuth(registryURL string) bool {
    // Docker Hub and most registries require auth
    // Local registries (localhost, 127.0.0.1) typically don't
    if strings.Contains(registryURL, "localhost") ||
       strings.Contains(registryURL, "127.0.0.1") {
        return false
    }
    return true
}
```

### Step 4: Add Environment Variable Helper (helpers/env.go)
```go
package helpers

import (
    "os"
    "strconv"
)

// GetEnvString returns environment variable value or default
func GetEnvString(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

// GetEnvBool returns environment variable as bool or default
func GetEnvBool(key string, defaultValue bool) bool {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }

    b, err := strconv.ParseBool(value)
    if err != nil {
        return defaultValue
    }
    return b
}
```

### Step 5: Create Comprehensive Unit Tests (push_test.go)
```go
package push

import (
    "os"
    "testing"

    "github.com/spf13/cobra"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestBuildPushOptions(t *testing.T) {
    tests := []struct {
        name     string
        args     []string
        flags    map[string]string
        env      map[string]string
        expected *PushOptions
        wantErr  bool
    }{
        {
            name: "flags override environment",
            args: []string{"myimage:latest", "registry.example.com"},
            flags: map[string]string{
                "username": "flaguser",
                "password": "flagpass",
            },
            env: map[string]string{
                "REGISTRY_USERNAME": "envuser",
                "REGISTRY_PASSWORD": "envpass",
            },
            expected: &PushOptions{
                ImageRef:    "myimage:latest",
                RegistryURL: "registry.example.com",
                Username:    "flaguser",
                Password:    "flagpass",
            },
        },
        {
            name: "environment variables used when flags not set",
            args: []string{"myimage:latest", "registry.example.com"},
            env: map[string]string{
                "REGISTRY_USERNAME": "envuser",
                "REGISTRY_PASSWORD": "envpass",
                "REGISTRY_INSECURE": "true",
            },
            expected: &PushOptions{
                ImageRef:    "myimage:latest",
                RegistryURL: "registry.example.com",
                Username:    "envuser",
                Password:    "envpass",
                Insecure:    true,
            },
        },
        {
            name: "insecure flag",
            args: []string{"myimage:latest", "registry.example.com"},
            flags: map[string]string{
                "insecure": "true",
            },
            expected: &PushOptions{
                ImageRef:    "myimage:latest",
                RegistryURL: "registry.example.com",
                Insecure:    true,
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Clear environment
            os.Clearenv()

            // Set environment variables
            for k, v := range tt.env {
                os.Setenv(k, v)
            }

            // Reset flags
            username = ""
            password = ""
            insecure = false

            // Set flags
            for k, v := range tt.flags {
                switch k {
                case "username":
                    username = v
                case "password":
                    password = v
                case "insecure":
                    insecure = v == "true"
                }
            }

            // Build options
            cmd := &cobra.Command{}
            opts, err := buildPushOptions(cmd, tt.args)

            if tt.wantErr {
                require.Error(t, err)
                return
            }

            require.NoError(t, err)
            assert.Equal(t, tt.expected.ImageRef, opts.ImageRef)
            assert.Equal(t, tt.expected.RegistryURL, opts.RegistryURL)
            assert.Equal(t, tt.expected.Username, opts.Username)
            assert.Equal(t, tt.expected.Password, opts.Password)
            assert.Equal(t, tt.expected.Insecure, opts.Insecure)
        })
    }
}

func TestValidateRegistryURL(t *testing.T) {
    tests := []struct {
        name    string
        url     string
        wantErr bool
    }{
        {
            name:    "valid https URL",
            url:     "https://registry.example.com",
            wantErr: false,
        },
        {
            name:    "valid http URL",
            url:     "http://localhost:5000",
            wantErr: false,
        },
        {
            name:    "URL without scheme",
            url:     "registry.example.com",
            wantErr: false, // Should add https://
        },
        {
            name:    "invalid scheme",
            url:     "ftp://registry.example.com",
            wantErr: true,
        },
        {
            name:    "empty URL",
            url:     "",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validateRegistryURL(tt.url)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}

func TestValidatePushOptions(t *testing.T) {
    tests := []struct {
        name    string
        opts    *PushOptions
        wantErr bool
        errMsg  string
    }{
        {
            name: "valid options with auth",
            opts: &PushOptions{
                ImageRef:    "myimage:latest",
                RegistryURL: "registry.example.com",
                Username:    "user",
                Password:    "pass",
            },
            wantErr: false,
        },
        {
            name: "missing image reference",
            opts: &PushOptions{
                RegistryURL: "registry.example.com",
            },
            wantErr: true,
            errMsg:  "image reference is required",
        },
        {
            name: "missing registry URL",
            opts: &PushOptions{
                ImageRef: "myimage:latest",
            },
            wantErr: true,
            errMsg:  "registry URL is required",
        },
        {
            name: "localhost doesn't require auth",
            opts: &PushOptions{
                ImageRef:    "myimage:latest",
                RegistryURL: "localhost:5000",
            },
            wantErr: false,
        },
        {
            name: "remote registry requires auth",
            opts: &PushOptions{
                ImageRef:    "myimage:latest",
                RegistryURL: "registry.example.com",
            },
            wantErr: true,
            errMsg:  "authentication required",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validatePushOptions(tt.opts)
            if tt.wantErr {
                require.Error(t, err)
                if tt.errMsg != "" {
                    assert.Contains(t, err.Error(), tt.errMsg)
                }
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

## Size Management
- **Estimated Lines**: 450 lines total
  - push.go: 120 lines
  - flags.go: 80 lines
  - validation.go: 70 lines
  - push_test.go: 150 lines
  - helpers/env.go: 30 lines
- **Measurement Tool**: `/home/vscode/workspaces/idpbuilder-push-oci/tools/line-counter.sh`
- **Check Frequency**: After each file implementation
- **Split Threshold**: 700 lines (warning), 800 lines (stop)

## Test Requirements
- **Unit Tests**: 85% coverage minimum
  - Flag parsing and environment variable fallback
  - Validation logic for all inputs
  - Error handling paths
  - Dry-run mode behavior
- **Integration Tests**: Will be tested with E1.1.3 framework
- **Test Files**:
  - `pkg/cmd/push/push_test.go` - comprehensive unit tests
  - Integration tests will use existing E1.1.3 infrastructure

## Pattern Compliance
- **idpbuilder Patterns**:
  - Cobra command structure matching existing commands
  - Logger usage consistent with helpers package
  - Flag naming conventions following project standards
- **Go Best Practices**:
  - Options struct pattern for configuration
  - Clear separation of concerns (command, flags, validation)
  - Comprehensive error messages with context
- **Security Requirements**:
  - Password flag marked as hidden
  - Clear warnings for insecure mode
  - No credential logging

## Dependencies on Other Efforts
- **E1.2.2 (Registry Authentication)**: Will provide enhanced auth mechanisms
- **E1.2.3 (Image Push Operations)**: Will implement actual push logic using go-containerregistry
- This effort provides the command structure foundation that E1.2.2 and E1.2.3 will build upon

## Implementation Order
1. Create `pkg/cmd/push/` directory structure
2. Implement `push.go` with command structure and options
3. Implement `flags.go` with all flag definitions
4. Implement `validation.go` with input validation
5. Add `helpers/env.go` for environment utilities
6. Write comprehensive unit tests in `push_test.go`
7. Run tests and ensure >85% coverage
8. Verify integration with root command

## Success Criteria
- ✅ Push command properly registered with root command
- ✅ All flags work correctly (username, password, insecure, etc.)
- ✅ Environment variable fallback works when flags not provided
- ✅ Validation catches all invalid inputs with clear messages
- ✅ Dry-run mode prevents actual operations
- ✅ Insecure mode shows appropriate warnings
- ✅ Unit tests achieve >85% coverage
- ✅ Code follows idpbuilder patterns and conventions

## Notes for Software Engineer
1. **Important**: The basic push command already exists from Wave 1 with just the insecure flag. You need to enhance it with full functionality.
2. **Library Note**: go-containerregistry is already available in go.mod (v0.20.3)
3. **Logger**: Use `helpers.CmdLogger` for logging, it's already initialized by Wave 1
4. **Testing**: Use the test utilities from E1.1.2 in `pkg/testutils/` for mocking if needed
5. **TODO Markers**: Add TODO comments where E1.2.2 and E1.2.3 will integrate
6. **Branch Strategy**: This effort can be parallelized with E1.2.2 and E1.2.3

## Risk Mitigation
- **Risk**: Command structure doesn't align with E1.2.3 needs
  - **Mitigation**: Use extensible PushOptions struct that can be enhanced
- **Risk**: Flag conflicts with existing idpbuilder commands
  - **Mitigation**: Follow established patterns, use standard flag names