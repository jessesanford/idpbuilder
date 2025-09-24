# Effort 2.2.2: Authentication Flow Implementation Plan

## Overview
- **Effort**: Implement Authentication Flow (TDD GREEN-REFACTOR Phase)
- **Phase**: 2 (Authentication & Credentials), Wave: 2
- **Estimated Size**: 300 lines
- **Implementation Time**: 2-3 hours
- **TDD Phase**: GREEN (make tests pass) then REFACTOR

## 🎯 Objectives

Implement the complete authentication flow with:
1. **Precedence Logic**: CLI flags → Environment → Docker config → Secrets → Anonymous
2. **Error Propagation**: Clear, helpful error messages at each level
3. **Integration**: Build on Wave 1 auth components
4. **Clean Code**: Refactor for maintainability after tests pass

## 📋 Dependencies from Wave 1

### Available Components
- **pkg/oci/types.go**: `Authenticator` interface, `Credentials` struct
- **pkg/oci/auth.go**: `DefaultAuthenticator` with:
  - `NewAuthenticatorFromFlags()` - CLI flag authentication
  - `NewAuthenticatorFromEnv()` - Environment variable authentication
  - `NewAuthenticatorFromSecrets()` - Kubernetes secrets authentication
  - Docker config loading methods

### Test Contract from Effort 2.2.1
- **pkg/oci/flow_test.go**: Defines the expected API and behavior
- Tests define precedence order and error handling requirements
- Mock structures show expected implementation approach

## 📁 File Structure

```
pkg/oci/
├── flow.go           # Main authentication flow implementation (NEW - 300 lines)
└── flow_test.go      # Existing tests from effort 2.2.1 (will turn GREEN)
```

## 🚀 Implementation Steps

### Step 1: Create AuthFlow Structure (50 lines)
```go
// pkg/oci/flow.go
package oci

type AuthFlow struct {
    // CLI flag values (highest priority)
    FlagUsername string
    FlagPassword string

    // Environment variables
    EnvUsername string
    EnvPassword string

    // Docker config path
    DockerConfigPath string

    // Default credentials from secrets
    DefaultSecretsData map[string][]byte

    // Allow anonymous access
    AllowAnonymous bool

    // Underlying authenticator
    authenticator Authenticator
}

type AuthFlowConfig struct {
    FlagUsername       string
    FlagPassword       string
    DockerConfigPath   string
    DefaultSecretsData map[string][]byte
    AllowAnonymous     bool
}
```

### Step 2: Implement NewAuthFlow Constructor (40 lines)
```go
func NewAuthFlow(config *AuthFlowConfig) (*AuthFlow, error) {
    // Validate config
    // Initialize AuthFlow with precedence awareness
    // Set up environment reading
    // Return configured flow
}
```

### Step 3: Implement Main Authenticate Method (100 lines)
```go
func (f *AuthFlow) Authenticate(ctx context.Context, registry string) (*Credentials, error) {
    // 1. Check CLI flags first (highest priority)
    if f.FlagUsername != "" && f.FlagPassword != "" {
        auth, err := NewAuthenticatorFromFlags(f.FlagUsername, f.FlagPassword)
        if err == nil {
            creds, authErr := auth.Authenticate(ctx, registry)
            if authErr == nil && creds != nil {
                return creds, nil
            }
        }
    }

    // 2. Check environment variables
    if envAuth, err := f.tryEnvironmentAuth(ctx, registry); err == nil {
        return envAuth, nil
    }

    // 3. Check Docker config
    if dockerAuth, err := f.tryDockerConfigAuth(ctx, registry); err == nil {
        return dockerAuth, nil
    }

    // 4. Check default secrets
    if secretsAuth, err := f.trySecretsAuth(ctx, registry); err == nil {
        return secretsAuth, nil
    }

    // 5. Check if anonymous is allowed
    if f.AllowAnonymous {
        return f.anonymousCredentials(registry), nil
    }

    // Return comprehensive error with all attempted methods
    return nil, f.buildAuthError(registry)
}
```

### Step 4: Implement Helper Methods (80 lines)
```go
// Try authentication with environment variables
func (f *AuthFlow) tryEnvironmentAuth(ctx context.Context, registry string) (*Credentials, error) {
    // Check for OCI_USERNAME and OCI_PASSWORD
    // Use NewAuthenticatorFromEnv() from Wave 1
    // Return credentials or error with context
}

// Try authentication with Docker config
func (f *AuthFlow) tryDockerConfigAuth(ctx context.Context, registry string) (*Credentials, error) {
    // Create authenticator with Docker config source
    // Attempt authentication
    // Return credentials or error with context
}

// Try authentication with default secrets
func (f *AuthFlow) trySecretsAuth(ctx context.Context, registry string) (*Credentials, error) {
    // Check if secrets data is available
    // Use NewAuthenticatorFromSecrets() from Wave 1
    // Return credentials or error with context
}

// Create anonymous credentials
func (f *AuthFlow) anonymousCredentials(registry string) *Credentials {
    // Return minimal credentials for anonymous access
}

// Build comprehensive error message
func (f *AuthFlow) buildAuthError(registry string) error {
    // List all attempted methods
    // Provide helpful debugging information
    // Suggest possible solutions
}
```

### Step 5: Add Convenience Functions (30 lines)
```go
// CreateAuthFlowFromFlags creates an auth flow with CLI flag precedence
func CreateAuthFlowFromFlags(username, password string) (*AuthFlow, error) {
    config := &AuthFlowConfig{
        FlagUsername: username,
        FlagPassword: password,
    }
    return NewAuthFlow(config)
}

// CreateDefaultAuthFlow creates an auth flow with standard precedence
func CreateDefaultAuthFlow() (*AuthFlow, error) {
    // Read environment, detect Docker config, etc.
}
```

### Step 6: REFACTOR Phase (After Tests Pass)
1. **Extract Constants**: Define precedence levels as constants
2. **Improve Error Messages**: Add detailed context for each failure
3. **Add Logging**: Debug-level logging for auth attempts
4. **Optimize Performance**: Cache successful auth methods
5. **Documentation**: Add comprehensive godoc comments

## 🧪 Test Validation

### Tests to Make GREEN:
1. **TestAuthenticationPrecedence_FlagsOverrideSecrets** - Verify CLI flags have highest priority
2. **TestAuthenticationPrecedence_SecretsAsDefault** - Verify secrets used as fallback
3. **TestAuthenticationFailure_InvalidCredentials** - Proper error handling
4. **TestAuthenticationPrecedence_EnvOverrideSecrets** - Environment variables precedence
5. **TestAuthenticationFlow_DockerConfigFallback** - Docker config as fallback
6. **TestAuthenticationFlow_AnonymousAccess** - Anonymous when allowed
7. **TestAuthenticationFlow_AllSourcesFail** - Comprehensive error reporting

## 📏 Size Management

- **Target**: 300 lines
- **Breakdown**:
  - Core structures: 50 lines
  - Main authenticate: 100 lines
  - Helper methods: 80 lines
  - Constructor/convenience: 70 lines
- **Measurement Tool**: `$PROJECT_ROOT/tools/line-counter.sh`
- **Check Frequency**: After each major method implementation

## ✅ Success Criteria

1. **All tests pass** from effort 2.2.1 (flow_test.go)
2. **Clear precedence**: Flags → Env → Docker → Secrets → Anonymous
3. **Helpful errors**: Each failure explains what was tried
4. **Clean integration**: Uses Wave 1 components properly
5. **Maintainable code**: Well-structured after refactor
6. **Size compliant**: Under 300 lines

## 🔴 RED Phase Validation (Already Complete)
- Tests written in effort 2.2.1 define the contract
- All tests currently fail (expected in RED phase)
- Clear understanding of required behavior

## 🟢 GREEN Phase Implementation (This Effort)
- Implement minimal code to make tests pass
- Focus on correctness, not optimization
- Verify each test turns green

## 🔵 REFACTOR Phase (After GREEN)
- Improve code structure and readability
- Extract common patterns
- Add comprehensive error context
- Optimize for performance where needed

## 📝 Notes for SW Engineer

1. **Start with failing tests**: Run tests first to see RED state
2. **Implement incrementally**: Make one test pass at a time
3. **Use Wave 1 components**: Don't reimplement existing auth logic
4. **Keep it simple**: Minimal implementation for GREEN phase
5. **Document as you go**: Add comments for complex logic
6. **Refactor only after GREEN**: Get tests passing first

## 🚫 Out of Scope

- OAuth/OIDC flows (future enhancement)
- Token refresh logic (handled by authenticator)
- Registry-specific authentication (generic approach)
- Credential storage/persistence (handled by authenticator)

## 🎯 Command to Verify Implementation

```bash
# From effort directory
cd pkg/oci
go test -v -run "^TestAuthentication" ./...
```

All tests should pass after implementation is complete.