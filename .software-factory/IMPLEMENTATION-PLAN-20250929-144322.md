# E1.2.2 Registry Authentication - Implementation Plan

## 🚨 CRITICAL EFFORT METADATA
**Branch**: `phase1/wave2/registry-authentication`
**Can Parallelize**: Yes
**Parallel With**: E1.2.3 (image-push-operations)
**Size Estimate**: 650-700 lines
**Dependencies**: E1.2.1 (command-structure)

## Overview
- **Effort**: Implement registry authentication and credential handling for idpbuilder push command
- **Phase**: 1, Wave: 2
- **Estimated Size**: 650-700 lines
- **Implementation Time**: 4-6 hours

## Context from Master Plan
This effort implements registry authentication for the idpbuilder push command. Key requirements:
- Implement credential handling with precedence: flags > environment variables
- Add --insecure flag support for self-signed certificates
- Create comprehensive authentication error handling
- Add retry logic with exponential backoff (5-10 retries max)
- Use go-containerregistry library for OCI authentication

## File Structure
```
pkg/
├── push/
│   ├── auth/
│   │   ├── authenticator.go      # Main auth logic (120 lines)
│   │   ├── credentials.go        # Credential handling (80 lines)
│   │   ├── insecure.go          # Insecure mode support (60 lines)
│   │   └── auth_test.go         # Unit tests (150 lines)
│   ├── retry/
│   │   ├── backoff.go           # Exponential backoff (70 lines)
│   │   ├── retry.go             # Retry wrapper (90 lines)
│   │   └── retry_test.go        # Unit tests (100 lines)
│   └── errors/
│       └── auth_errors.go       # Auth error types (30 lines)
```

## Implementation Steps

### Step 1: Create Authentication Module Structure (50 lines)
```go
// pkg/push/auth/authenticator.go
package auth

import (
    "github.com/google/go-containerregistry/pkg/authn"
    "github.com/google/go-containerregistry/pkg/v1/remote"
)

type Authenticator struct {
    username   string
    password   string
    insecure   bool
    keychain   authn.Keychain
}

func NewAuthenticator(opts ...Option) (*Authenticator, error)
func (a *Authenticator) GetAuthOptions() []remote.Option
func (a *Authenticator) Validate() error
```

### Step 2: Implement Credential Handling (80 lines)
```go
// pkg/push/auth/credentials.go
package auth

// Credential precedence logic
type CredentialSource int
const (
    SourceNone CredentialSource = iota
    SourceEnv
    SourceFlags
)

type Credentials struct {
    Username string
    Password string
    Source   CredentialSource
}

// GetCredentials returns credentials with proper precedence:
// 1. Command line flags (highest priority)
// 2. Environment variables (fallback)
// 3. Docker config (if available)
func GetCredentials(flagUser, flagPass string) (*Credentials, error) {
    // Implementation will check flags first
    if flagUser != "" && flagPass != "" {
        return &Credentials{
            Username: flagUser,
            Password: flagPass,
            Source:   SourceFlags,
        }, nil
    }

    // Fall back to environment variables
    envUser := os.Getenv("IDPBUILDER_REGISTRY_USER")
    envPass := os.Getenv("IDPBUILDER_REGISTRY_PASSWORD")
    if envUser != "" && envPass != "" {
        return &Credentials{
            Username: envUser,
            Password: envPass,
            Source:   SourceEnv,
        }, nil
    }

    // Fall back to docker config (using go-containerregistry's DefaultKeychain)
    return nil, nil // Use default keychain
}
```

### Step 3: Implement Insecure Mode Support (60 lines)
```go
// pkg/push/auth/insecure.go
package auth

import (
    "crypto/tls"
    "net/http"
    "github.com/google/go-containerregistry/pkg/v1/remote"
)

// GetInsecureTransport returns an HTTP transport that skips TLS verification
func GetInsecureTransport() *http.Transport {
    return &http.Transport{
        TLSClientConfig: &tls.Config{
            InsecureSkipVerify: true,
        },
    }
}

// GetInsecureOption returns a remote.Option for insecure registries
func GetInsecureOption() remote.Option {
    return remote.WithTransport(GetInsecureTransport())
}
```

### Step 4: Implement Retry Logic with Backoff (160 lines)
```go
// pkg/push/retry/backoff.go
package retry

import (
    "math"
    "math/rand"
    "time"
)

type BackoffStrategy struct {
    InitialInterval time.Duration
    MaxInterval     time.Duration
    Multiplier      float64
    MaxRetries      int
    Jitter          bool
}

func DefaultBackoff() *BackoffStrategy {
    return &BackoffStrategy{
        InitialInterval: 1 * time.Second,
        MaxInterval:     30 * time.Second,
        Multiplier:      2.0,
        MaxRetries:      10, // 5-10 retries as specified
        Jitter:          true,
    }
}

func (b *BackoffStrategy) NextInterval(attempt int) time.Duration {
    if attempt >= b.MaxRetries {
        return 0
    }

    interval := float64(b.InitialInterval) * math.Pow(b.Multiplier, float64(attempt))
    if interval > float64(b.MaxInterval) {
        interval = float64(b.MaxInterval)
    }

    if b.Jitter {
        interval = interval * (0.5 + rand.Float64()*0.5)
    }

    return time.Duration(interval)
}
```

```go
// pkg/push/retry/retry.go
package retry

import (
    "context"
    "fmt"
    "time"
)

type RetryableFunc func() error

// IsRetryable determines if an error should trigger a retry
func IsRetryable(err error) bool {
    // Check for transient errors:
    // - Network timeouts
    // - 503 Service Unavailable
    // - 429 Too Many Requests
    // - Connection reset
    // Implementation details...
}

// WithRetry executes a function with retry logic
func WithRetry(ctx context.Context, fn RetryableFunc, strategy *BackoffStrategy) error {
    var lastErr error

    for attempt := 0; attempt <= strategy.MaxRetries; attempt++ {
        if attempt > 0 {
            interval := strategy.NextInterval(attempt - 1)
            select {
            case <-time.After(interval):
                // Continue with retry
            case <-ctx.Done():
                return fmt.Errorf("retry cancelled: %w", ctx.Err())
            }
        }

        err := fn()
        if err == nil {
            return nil
        }

        lastErr = err
        if !IsRetryable(err) {
            return err // Non-retryable error
        }

        // Log retry attempt
        fmt.Printf("Attempt %d failed, retrying in %v: %v\n",
            attempt+1, strategy.NextInterval(attempt), err)
    }

    return fmt.Errorf("max retries (%d) exceeded: %w", strategy.MaxRetries, lastErr)
}
```

### Step 5: Create Authentication Error Types (30 lines)
```go
// pkg/push/errors/auth_errors.go
package errors

import "fmt"

type AuthenticationError struct {
    Registry string
    Cause    error
}

func (e *AuthenticationError) Error() string {
    return fmt.Sprintf("authentication failed for registry %s: %v", e.Registry, e.Cause)
}

type InsecureRegistryError struct {
    Registry string
}

func (e *InsecureRegistryError) Error() string {
    return fmt.Sprintf("registry %s requires --insecure flag for self-signed certificates", e.Registry)
}
```

### Step 6: Integration with Main Authenticator (120 lines)
```go
// pkg/push/auth/authenticator.go (full implementation)
package auth

import (
    "context"
    "fmt"
    "github.com/cnoe-io/idpbuilder/pkg/push/errors"
    "github.com/cnoe-io/idpbuilder/pkg/push/retry"
    "github.com/google/go-containerregistry/pkg/authn"
    "github.com/google/go-containerregistry/pkg/v1/remote"
)

type Option func(*Authenticator) error

func WithCredentials(username, password string) Option {
    return func(a *Authenticator) error {
        creds, err := GetCredentials(username, password)
        if err != nil {
            return err
        }
        if creds != nil {
            a.username = creds.Username
            a.password = creds.Password
        }
        return nil
    }
}

func WithInsecure(insecure bool) Option {
    return func(a *Authenticator) error {
        a.insecure = insecure
        return nil
    }
}

func NewAuthenticator(opts ...Option) (*Authenticator, error) {
    auth := &Authenticator{
        keychain: authn.DefaultKeychain,
    }

    for _, opt := range opts {
        if err := opt(auth); err != nil {
            return nil, err
        }
    }

    return auth, nil
}

func (a *Authenticator) GetAuthOptions() []remote.Option {
    var opts []remote.Option

    // Add authentication
    if a.username != "" && a.password != "" {
        auth := &authn.Basic{
            Username: a.username,
            Password: a.password,
        }
        opts = append(opts, remote.WithAuth(auth))
    } else {
        opts = append(opts, remote.WithAuthFromKeychain(a.keychain))
    }

    // Add insecure transport if needed
    if a.insecure {
        opts = append(opts, GetInsecureOption())
    }

    // Add retry logic wrapper
    opts = append(opts, remote.WithRetryBackoff(retry.DefaultBackoff()))

    return opts
}

// AuthenticateWithRetry performs authentication with retry logic
func (a *Authenticator) AuthenticateWithRetry(ctx context.Context, registry string) error {
    strategy := retry.DefaultBackoff()

    return retry.WithRetry(ctx, func() error {
        // Test authentication by attempting to list catalog or ping
        opts := a.GetAuthOptions()

        // Attempt to authenticate
        // This would be called by the actual push operation
        // For now, we're setting up the structure

        return nil
    }, strategy)
}
```

### Step 7: Create Comprehensive Unit Tests (250 lines)
```go
// pkg/push/auth/auth_test.go
package auth_test

import (
    "testing"
    "os"
    "github.com/cnoe-io/idpbuilder/pkg/push/auth"
)

func TestCredentialPrecedence(t *testing.T) {
    tests := []struct {
        name      string
        flagUser  string
        flagPass  string
        envUser   string
        envPass   string
        wantUser  string
        wantPass  string
        wantSource auth.CredentialSource
    }{
        {
            name:      "flags take precedence over env",
            flagUser:  "flag-user",
            flagPass:  "flag-pass",
            envUser:   "env-user",
            envPass:   "env-pass",
            wantUser:  "flag-user",
            wantPass:  "flag-pass",
            wantSource: auth.SourceFlags,
        },
        {
            name:      "env vars used when no flags",
            flagUser:  "",
            flagPass:  "",
            envUser:   "env-user",
            envPass:   "env-pass",
            wantUser:  "env-user",
            wantPass:  "env-pass",
            wantSource: auth.SourceEnv,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Set env vars
            os.Setenv("IDPBUILDER_REGISTRY_USER", tt.envUser)
            os.Setenv("IDPBUILDER_REGISTRY_PASSWORD", tt.envPass)
            defer os.Unsetenv("IDPBUILDER_REGISTRY_USER")
            defer os.Unsetenv("IDPBUILDER_REGISTRY_PASSWORD")

            creds, err := auth.GetCredentials(tt.flagUser, tt.flagPass)
            if err != nil {
                t.Fatalf("GetCredentials() error = %v", err)
            }

            if creds.Username != tt.wantUser {
                t.Errorf("Username = %v, want %v", creds.Username, tt.wantUser)
            }
            if creds.Password != tt.wantPass {
                t.Errorf("Password = %v, want %v", creds.Password, tt.wantPass)
            }
            if creds.Source != tt.wantSource {
                t.Errorf("Source = %v, want %v", creds.Source, tt.wantSource)
            }
        })
    }
}

func TestInsecureMode(t *testing.T) {
    // Test insecure transport configuration
    // Test TLS skip verify
}

func TestAuthenticatorOptions(t *testing.T) {
    // Test various authenticator configurations
    // Test option combinations
}
```

```go
// pkg/push/retry/retry_test.go
package retry_test

import (
    "context"
    "testing"
    "time"
    "errors"
    "github.com/cnoe-io/idpbuilder/pkg/push/retry"
)

func TestExponentialBackoff(t *testing.T) {
    strategy := retry.DefaultBackoff()

    // Test interval increases
    for i := 0; i < 5; i++ {
        interval := strategy.NextInterval(i)
        if interval == 0 {
            t.Errorf("Unexpected zero interval at attempt %d", i)
        }
        // Verify exponential growth with jitter
    }
}

func TestRetryLogic(t *testing.T) {
    tests := []struct {
        name        string
        attempts    int
        shouldRetry bool
        wantErr     bool
    }{
        {
            name:        "succeeds on first attempt",
            attempts:    1,
            shouldRetry: false,
            wantErr:     false,
        },
        {
            name:        "retries on transient error",
            attempts:    3,
            shouldRetry: true,
            wantErr:     false,
        },
        {
            name:        "fails after max retries",
            attempts:    11,
            shouldRetry: true,
            wantErr:     true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## Size Management
- **Estimated Lines**: 700 lines
- **Measurement Tool**: ${PROJECT_ROOT}/tools/line-counter.sh
- **Check Frequency**: After each major component
- **Split Threshold**: 700 lines (warning), 800 lines (stop)

## Test Requirements
- **Unit Tests**: 85% coverage minimum
- **Integration Tests**: Test against mock registry
- **E2E Tests**: Will be covered in Phase 2
- **Test Files**:
  - `pkg/push/auth/auth_test.go` - Authentication tests
  - `pkg/push/retry/retry_test.go` - Retry logic tests

## Pattern Compliance
- **idpbuilder Patterns**:
  - Use Cobra command patterns from existing commands
  - Follow existing error handling patterns
  - Maintain consistent logging approach
- **Security Requirements**:
  - Never log passwords in plain text
  - Use secure credential storage
  - Support secure TLS by default
- **Performance Targets**:
  - Retry delays should not exceed 30 seconds
  - Total retry time should not exceed 5 minutes

## Dependencies on Other Efforts
- **E1.2.1 (command-structure)**: Will integrate with the push command structure
  - Import command registration pattern
  - Use established flag handling
- **E1.2.3 (image-push-operations)**: This authentication will be used by push operations
  - Provides GetAuthOptions() for remote operations
  - Provides retry wrapper for push operations

## Integration Points
1. **Command Integration**: Authentication options will be passed to push command via flags
2. **Push Operation**: Authenticator will be created and used in push operations
3. **Error Handling**: Authentication errors will bubble up with proper context
4. **Logging**: Authentication attempts and retries will be logged appropriately

## Implementation Order
1. Create basic authenticator structure
2. Implement credential handling with precedence
3. Add insecure mode support
4. Implement retry logic with exponential backoff
5. Create comprehensive error types
6. Write unit tests for all components
7. Integrate with command structure (from E1.2.1)

## Known Limitations
- Rate limiting handling will be added in future phases (TODO marker)
- Registry-specific authentication methods not yet supported
- No OAuth/OIDC support in this phase

## Success Criteria
- ✅ Authentication works with username/password from flags
- ✅ Falls back to environment variables when flags not provided
- ✅ --insecure flag properly handles self-signed certificates
- ✅ Retry logic handles transient failures (5-10 retries)
- ✅ Clear error messages for authentication failures
- ✅ Unit test coverage >85%
- ✅ No credentials logged in plain text

## Notes for Software Engineer
1. Start with the authenticator.go as the main entry point
2. Ensure all error messages are user-friendly
3. Add debug logging for troubleshooting (use existing idpbuilder logging patterns)
4. The retry logic should be reusable for other operations
5. Consider making the backoff strategy configurable in future phases
6. Add TODO markers for future OAuth/OIDC support