# Credential Management Implementation Plan

## Effort ID: E2.2.2-A

## Overview
- **Effort**: Implement real credential management for Gitea registry authentication
- **Phase**: 2, Wave: 2
- **Estimated Size**: 450-500 lines (excluding tests)
- **Implementation Time**: 4-6 hours
- **Dependencies**: E2.2.1 cli-commands (complete)

## Scope
Replace placeholder authentication with a comprehensive credential management system supporting:
1. CLI flags for username and token (--username, --token)
2. Environment variable credential loading (fallback)
3. Configuration file parsing (fallback)
4. System keyring integration for secure storage (fallback)
5. Proper credential retrieval functions
6. Removal of all credential-related TODOs

## File Structure

### New Files to Create
1. **pkg/gitea/credentials.go** (150 lines)
   - Core credential management logic
   - Credential provider interface
   - Credential resolution chain

2. **pkg/gitea/config.go** (70 lines)
   - Configuration file parsing
   - Config structure definitions
   - Config validation

3. **pkg/gitea/keyring.go** (80 lines)
   - System keyring integration
   - Secure credential storage/retrieval
   - Fallback handling

4. **pkg/gitea/credentials_test.go** (100 lines)
   - Unit tests for credential providers
   - Mock implementations
   - Test coverage for all scenarios

5. **pkg/gitea/config_test.go** (50 lines)
   - Config parsing tests
   - Validation tests
   - Error handling tests

6. **pkg/gitea/keyring_test.go** (50 lines)
   - Keyring integration tests
   - Mock keyring provider
   - Fallback scenario tests

### Files to Modify
1. **pkg/gitea/client.go** (~60 line changes)
   - Update getRegistryUsername() function
   - Update getRegistryPassword() function
   - Remove TODO comments at lines 30, 145, 151
   - Add credential manager initialization
   - Add methods to set CLI-provided credentials

2. **pkg/gitea/client_test.go** (~50 line changes)
   - Update tests to use real credential system
   - Add environment variable tests
   - Remove hardcoded credential expectations

3. **pkg/cmd/push.go** (~30 line changes)
   - Add --username and --token flags
   - Pass CLI credentials to client
   - Update help text and examples

4. **pkg/cmd/build.go** (~30 line changes)
   - Add --username and --token flags for registry operations
   - Pass CLI credentials to builder
   - Update help text

## Implementation Steps

### Step 1: Create Credential Provider Interface (credentials.go)
```go
// Create pkg/gitea/credentials.go with:

package gitea

import (
    "fmt"
    "os"
    "path/filepath"
)

// CredentialProvider defines interface for credential sources
type CredentialProvider interface {
    GetUsername() (string, error)
    GetPassword() (string, error)
    IsAvailable() bool
    Priority() int
}

// CredentialManager manages multiple credential providers
type CredentialManager struct {
    providers []CredentialProvider
    cliUsername string
    cliPassword string
}

// NewCredentialManager creates a credential manager with default providers
func NewCredentialManager() *CredentialManager {
    return &CredentialManager{
        providers: []CredentialProvider{
            NewCLICredentialProvider(),
            NewEnvCredentialProvider(),
            NewConfigFileProvider(),
            NewKeyringProvider(),
        },
    }
}

// SetCLICredentials sets credentials provided via CLI flags
func (cm *CredentialManager) SetCLICredentials(username, password string) {
    if cliProvider, ok := cm.providers[0].(*CLICredentialProvider); ok {
        cliProvider.SetCredentials(username, password)
    }
}

// GetCredentials retrieves credentials from the first available provider
func (cm *CredentialManager) GetCredentials() (username, password string, err error) {
    for _, provider := range cm.providers {
        if provider.IsAvailable() {
            username, err = provider.GetUsername()
            if err != nil {
                continue
            }
            password, err = provider.GetPassword()
            if err != nil {
                continue
            }
            return username, password, nil
        }
    }
    return "", "", fmt.Errorf("no credentials available from any provider")
}

// EnvCredentialProvider reads from environment variables
type EnvCredentialProvider struct{}

func NewEnvCredentialProvider() *EnvCredentialProvider {
    return &EnvCredentialProvider{}
}

func (e *EnvCredentialProvider) GetUsername() (string, error) {
    username := os.Getenv("GITEA_USERNAME")
    if username == "" {
        return "", fmt.Errorf("GITEA_USERNAME not set")
    }
    return username, nil
}

func (e *EnvCredentialProvider) GetPassword() (string, error) {
    password := os.Getenv("GITEA_PASSWORD")
    if password == "" {
        return "", fmt.Errorf("GITEA_PASSWORD not set")
    }
    return password, nil
}

func (e *EnvCredentialProvider) IsAvailable() bool {
    return os.Getenv("GITEA_USERNAME") != "" && os.Getenv("GITEA_PASSWORD") != ""
}

func (e *EnvCredentialProvider) Priority() int {
    return 2 // Second priority after CLI
}

// CLICredentialProvider holds credentials from command-line flags
type CLICredentialProvider struct {
    username string
    password string
}

func NewCLICredentialProvider() *CLICredentialProvider {
    return &CLICredentialProvider{}
}

func (c *CLICredentialProvider) SetCredentials(username, password string) {
    c.username = username
    c.password = password
}

func (c *CLICredentialProvider) GetUsername() (string, error) {
    if c.username == "" {
        return "", fmt.Errorf("no CLI username provided")
    }
    return c.username, nil
}

func (c *CLICredentialProvider) GetPassword() (string, error) {
    if c.password == "" {
        return "", fmt.Errorf("no CLI password/token provided")
    }
    return c.password, nil
}

func (c *CLICredentialProvider) IsAvailable() bool {
    return c.username != "" && c.password != ""
}

func (c *CLICredentialProvider) Priority() int {
    return 1 // Highest priority
}
```

### Step 2: Create Configuration File Support (config.go)
```go
// Create pkg/gitea/config.go with:

package gitea

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
)

// Config represents the configuration file structure
type Config struct {
    Registries map[string]RegistryCredentials `json:"registries"`
}

// RegistryCredentials holds credentials for a specific registry
type RegistryCredentials struct {
    Username string `json:"username"`
    Password string `json:"password"`
    URL      string `json:"url"`
}

// ConfigFileProvider reads from ~/.idpbuilder/config
type ConfigFileProvider struct {
    configPath string
    config     *Config
}

func NewConfigFileProvider() *ConfigFileProvider {
    homeDir, _ := os.UserHomeDir()
    configPath := filepath.Join(homeDir, ".idpbuilder", "config")
    return &ConfigFileProvider{
        configPath: configPath,
    }
}

func (c *ConfigFileProvider) loadConfig() error {
    if c.config != nil {
        return nil
    }

    data, err := os.ReadFile(c.configPath)
    if err != nil {
        return err
    }

    var config Config
    if err := json.Unmarshal(data, &config); err != nil {
        return err
    }

    c.config = &config
    return nil
}

func (c *ConfigFileProvider) GetUsername() (string, error) {
    if err := c.loadConfig(); err != nil {
        return "", err
    }

    // Look for gitea registry config
    for name, creds := range c.config.Registries {
        if name == "gitea" || name == "default" {
            return creds.Username, nil
        }
    }

    return "", fmt.Errorf("no gitea credentials in config")
}

func (c *ConfigFileProvider) GetPassword() (string, error) {
    if err := c.loadConfig(); err != nil {
        return "", err
    }

    for name, creds := range c.config.Registries {
        if name == "gitea" || name == "default" {
            return creds.Password, nil
        }
    }

    return "", fmt.Errorf("no gitea credentials in config")
}

func (c *ConfigFileProvider) IsAvailable() bool {
    if _, err := os.Stat(c.configPath); err != nil {
        return false
    }
    return true
}

func (c *ConfigFileProvider) Priority() int {
    return 3 // Third priority
}
```

### Step 3: Create Keyring Integration (keyring.go)
```go
// Create pkg/gitea/keyring.go with:

package gitea

import (
    "fmt"
    "github.com/zalando/go-keyring"
)

const (
    keyringService = "idpbuilder"
    keyringUser    = "gitea"
)

// KeyringProvider provides credentials from system keyring
type KeyringProvider struct {
    service string
    user    string
}

func NewKeyringProvider() *KeyringProvider {
    return &KeyringProvider{
        service: keyringService,
        user:    keyringUser,
    }
}

func (k *KeyringProvider) GetUsername() (string, error) {
    username, err := keyring.Get(k.service, k.user+"_username")
    if err != nil {
        return "", fmt.Errorf("failed to get username from keyring: %w", err)
    }
    return username, nil
}

func (k *KeyringProvider) GetPassword() (string, error) {
    password, err := keyring.Get(k.service, k.user+"_password")
    if err != nil {
        return "", fmt.Errorf("failed to get password from keyring: %w", err)
    }
    return password, nil
}

func (k *KeyringProvider) IsAvailable() bool {
    // Check if keyring is accessible
    _, err := keyring.Get(k.service, k.user+"_username")
    return err == nil || err == keyring.ErrNotFound
}

func (k *KeyringProvider) Priority() int {
    return 4 // Lowest priority
}

// SetCredentials stores credentials in the keyring
func (k *KeyringProvider) SetCredentials(username, password string) error {
    if err := keyring.Set(k.service, k.user+"_username", username); err != nil {
        return fmt.Errorf("failed to store username: %w", err)
    }
    if err := keyring.Set(k.service, k.user+"_password", password); err != nil {
        return fmt.Errorf("failed to store password: %w", err)
    }
    return nil
}

// DeleteCredentials removes credentials from the keyring
func (k *KeyringProvider) DeleteCredentials() error {
    keyring.Delete(k.service, k.user+"_username")
    keyring.Delete(k.service, k.user+"_password")
    return nil
}
```

### Step 4: Update push.go and build.go for CLI flags

#### Update pkg/cmd/push.go:
```go
// Add to var declarations:
var (
    pushUsername string
    pushToken    string
    // ... existing vars
)

// Update init function:
func init() {
    PushCmd.Flags().BoolVar(&pushInsecure, "insecure", false, "Skip certificate verification (not recommended)")
    PushCmd.Flags().StringVar(&pushRegistry, "registry", getDefaultRegistry(), "Target registry")
    PushCmd.Flags().StringVar(&pushUsername, "username", "", "Registry username")
    PushCmd.Flags().StringVar(&pushToken, "token", "", "Registry token/password")
}

// Update runPush function to pass credentials:
func runPush(cmd *cobra.Command, args []string) error {
    // ... existing code ...

    // After creating client, set CLI credentials if provided
    if pushUsername != "" && pushToken != "" {
        client.SetCredentials(pushUsername, pushToken)
    }

    // ... rest of function
}
```

#### Update pkg/cmd/build.go:
```go
// Add similar flags for build command if it needs registry access
var (
    buildUsername string
    buildToken    string
    // ... existing vars
)

func init() {
    // ... existing flags ...
    BuildCmd.Flags().StringVar(&buildUsername, "username", "", "Registry username for base image pulls")
    BuildCmd.Flags().StringVar(&buildToken, "token", "", "Registry token/password for base image pulls")
}
```

### Step 5: Update client.go
```go
// Modify pkg/gitea/client.go:

// Add at package level:
var credentialManager *CredentialManager

func init() {
    credentialManager = NewCredentialManager()
}

// SetCredentials sets credentials from CLI flags
func (c *Client) SetCredentials(username, password string) {
    credentialManager.SetCLICredentials(username, password)
}

// Update getRegistryUsername function (line ~144):
func getRegistryUsername() string {
    username, password, err := credentialManager.GetCredentials()
    if err != nil {
        // Log warning and return empty string for backward compatibility
        // In production, this should probably return an error
        return ""
    }
    return username
}

// Update getRegistryPassword function (line ~150):
func getRegistryPassword() string {
    username, password, err := credentialManager.GetCredentials()
    if err != nil {
        // Log warning and return empty string for backward compatibility
        // In production, this should probably return an error
        return ""
    }
    return password
}

// Remove TODO comments at lines 30, 145, 151
```

### Step 6: Create Test Files

#### credentials_test.go
```go
package gitea

import (
    "os"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestEnvCredentialProvider(t *testing.T) {
    // Save and restore environment
    oldUsername := os.Getenv("GITEA_USERNAME")
    oldPassword := os.Getenv("GITEA_PASSWORD")
    defer func() {
        os.Setenv("GITEA_USERNAME", oldUsername)
        os.Setenv("GITEA_PASSWORD", oldPassword)
    }()

    // Test with environment variables set
    os.Setenv("GITEA_USERNAME", "testuser")
    os.Setenv("GITEA_PASSWORD", "testpass")

    provider := NewEnvCredentialProvider()
    assert.True(t, provider.IsAvailable())

    username, err := provider.GetUsername()
    require.NoError(t, err)
    assert.Equal(t, "testuser", username)

    password, err := provider.GetPassword()
    require.NoError(t, err)
    assert.Equal(t, "testpass", password)

    // Test with environment variables not set
    os.Unsetenv("GITEA_USERNAME")
    os.Unsetenv("GITEA_PASSWORD")

    assert.False(t, provider.IsAvailable())
}

func TestCredentialManager(t *testing.T) {
    // Test with mock providers
    manager := NewCredentialManager()

    // Set environment variables for testing
    os.Setenv("GITEA_USERNAME", "envuser")
    os.Setenv("GITEA_PASSWORD", "envpass")
    defer func() {
        os.Unsetenv("GITEA_USERNAME")
        os.Unsetenv("GITEA_PASSWORD")
    }()

    username, password, err := manager.GetCredentials()
    require.NoError(t, err)
    assert.Equal(t, "envuser", username)
    assert.Equal(t, "envpass", password)
}
```

#### config_test.go
```go
package gitea

import (
    "encoding/json"
    "os"
    "path/filepath"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestConfigFileProvider(t *testing.T) {
    // Create temporary config file
    tmpDir := t.TempDir()
    configPath := filepath.Join(tmpDir, ".idpbuilder", "config")
    os.MkdirAll(filepath.Dir(configPath), 0755)

    config := Config{
        Registries: map[string]RegistryCredentials{
            "gitea": {
                Username: "configuser",
                Password: "configpass",
                URL:      "https://gitea.example.com",
            },
        },
    }

    data, err := json.Marshal(config)
    require.NoError(t, err)

    err = os.WriteFile(configPath, data, 0600)
    require.NoError(t, err)

    // Test provider
    provider := &ConfigFileProvider{
        configPath: configPath,
    }

    assert.True(t, provider.IsAvailable())

    username, err := provider.GetUsername()
    require.NoError(t, err)
    assert.Equal(t, "configuser", username)

    password, err := provider.GetPassword()
    require.NoError(t, err)
    assert.Equal(t, "configpass", password)
}
```

#### keyring_test.go
```go
package gitea

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

// MockKeyringProvider for testing
type MockKeyringProvider struct {
    username string
    password string
    available bool
}

func (m *MockKeyringProvider) GetUsername() (string, error) {
    if m.username == "" {
        return "", fmt.Errorf("no username in keyring")
    }
    return m.username, nil
}

func (m *MockKeyringProvider) GetPassword() (string, error) {
    if m.password == "" {
        return "", fmt.Errorf("no password in keyring")
    }
    return m.password, nil
}

func (m *MockKeyringProvider) IsAvailable() bool {
    return m.available
}

func (m *MockKeyringProvider) Priority() int {
    return 3
}

func TestMockKeyringProvider(t *testing.T) {
    provider := &MockKeyringProvider{
        username:  "keyringuser",
        password:  "keyringpass",
        available: true,
    }

    assert.True(t, provider.IsAvailable())

    username, err := provider.GetUsername()
    assert.NoError(t, err)
    assert.Equal(t, "keyringuser", username)

    password, err := provider.GetPassword()
    assert.NoError(t, err)
    assert.Equal(t, "keyringpass", password)
}
```

### Step 7: Update client_test.go
Update the existing tests to work with the new credential system:
- Remove hardcoded "admin"/"password" expectations
- Add tests for environment variable integration
- Add tests for missing credentials scenarios
- Update TestClientWithEnvironmentVariables to actually test env vars

## Dependencies

### External Libraries
Add to go.mod:
```
github.com/zalando/go-keyring v0.2.3
```

### Import Updates
Update imports in client.go to include the credential manager initialization.

## Size Management
- **Estimated Lines**: 490-540 (excluding tests)
- **Breakdown**:
  - credentials.go: 180 lines (added CLI provider)
  - config.go: 70 lines
  - keyring.go: 80 lines
  - client.go updates: 60 lines
  - push.go updates: 30 lines
  - build.go updates: 30 lines
  - Test files: 250 lines (not counted toward limit)
- **Measurement Tool**: ${PROJECT_ROOT}/tools/line-counter.sh
- **Check Frequency**: After each file completion
- **Split Threshold**: N/A (well under 800 line limit)

## Test Requirements
- **Unit Tests**: 90% coverage
- **Integration Tests**: Test all three credential providers
- **E2E Tests**: Test with actual environment variables
- **Test Files**:
  - credentials_test.go
  - config_test.go
  - keyring_test.go
  - Updated client_test.go

## Integration Points

### 1. CLI Flags (Highest Priority)
- --username - Registry username
- --token - Registry password/token
- Available on both push and build commands
- Override all other credential sources

### 2. Client Initialization
- NewClient() and NewInsecureClient() automatically use credential manager
- SetCredentials() method to set CLI-provided credentials
- Backward compatible - returns empty strings if no credentials found

### 3. Environment Variables (Second Priority)
- GITEA_USERNAME - Gitea registry username
- GITEA_PASSWORD - Gitea registry password
- Used when no CLI flags provided

### 4. Configuration File (Third Priority)
- Location: ~/.idpbuilder/config
- JSON format with registry credentials
- Used when no CLI flags or env vars

### 5. System Keyring (Fourth Priority)
- Service: "idpbuilder"
- Account: "gitea"
- Most secure but lowest priority

## Success Criteria
- ✅ All credential-related TODOs removed from client.go
- ✅ Environment variables properly read when set
- ✅ Config file parsed when present
- ✅ Keyring integration functional (with graceful fallback)
- ✅ All tests passing with >90% coverage
- ✅ No hardcoded credentials remaining
- ✅ Backward compatible with existing code

## Security Considerations
1. **Credential Storage**: Never log passwords in plaintext
2. **File Permissions**: Config file should be 0600 (user read/write only)
3. **Error Messages**: Don't expose credential values in errors
4. **Keyring Access**: Handle keyring unavailability gracefully

## Error Handling Strategy
1. **Cascade Through Providers**: Try each provider in order
2. **Graceful Degradation**: Return empty strings if no credentials found
3. **Logging**: Log warnings when credentials unavailable (not errors)
4. **Testing**: Ensure all error paths are tested

## Migration Path
1. **Phase 1**: Implement credential providers
2. **Phase 2**: Update client.go to use providers
3. **Phase 3**: Update tests
4. **Phase 4**: Documentation update (if needed)
5. **Phase 5**: Remove all TODOs

## Validation Checklist
- [ ] All TODOs removed from client.go
- [ ] Environment variables working
- [ ] Config file parsing working
- [ ] Keyring integration working
- [ ] All tests passing
- [ ] Size under 500 lines
- [ ] No security vulnerabilities
- [ ] Backward compatible

## Notes for Implementation
1. Start with the credential provider interface to establish the contract
2. Implement each provider independently for easy testing
3. Use dependency injection for testability
4. Keep the credential manager singleton for simplicity
5. Ensure thread-safety if concurrent access is expected
6. Consider adding credential caching to avoid repeated lookups
7. Add appropriate logging for debugging credential issues