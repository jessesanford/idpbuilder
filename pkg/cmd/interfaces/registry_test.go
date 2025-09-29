package interfaces

import (
	"context"
	"testing"
)

// Mock implementations for testing registry interface compliance

// mockRegistryCommand implements RegistryCommand for testing
type mockRegistryCommand struct {
	configured      bool
	credentialsValid bool
	connectionTested bool
	registryInfo    *RegistryInfo
	capabilities    []string
}

func (m *mockRegistryCommand) ConfigureRegistry(config RegistryConfig) error {
	m.configured = true
	return nil
}

func (m *mockRegistryCommand) ValidateCredentials(ctx context.Context) error {
	m.credentialsValid = true
	return nil
}

func (m *mockRegistryCommand) GetRegistryInfo(ctx context.Context) (*RegistryInfo, error) {
	return m.registryInfo, nil
}

func (m *mockRegistryCommand) TestConnection(ctx context.Context) error {
	m.connectionTested = true
	return nil
}

func (m *mockRegistryCommand) GetCapabilities(ctx context.Context) ([]string, error) {
	return m.capabilities, nil
}

// mockRegistryManager implements RegistryManager for testing
type mockRegistryManager struct {
	registries     map[string]RegistryConfig
	defaultRegistry string
	validated      map[string]bool
}

func newMockRegistryManager() *mockRegistryManager {
	return &mockRegistryManager{
		registries: make(map[string]RegistryConfig),
		validated:  make(map[string]bool),
	}
}

func (m *mockRegistryManager) AddRegistry(name string, config RegistryConfig) error {
	m.registries[name] = config
	return nil
}

func (m *mockRegistryManager) RemoveRegistry(name string) error {
	delete(m.registries, name)
	return nil
}

func (m *mockRegistryManager) UpdateRegistry(name string, config RegistryConfig) error {
	if _, exists := m.registries[name]; exists {
		m.registries[name] = config
	}
	return nil
}

func (m *mockRegistryManager) GetRegistry(name string) (*RegistryConfig, error) {
	config, exists := m.registries[name]
	if !exists {
		return nil, &CLIError{Type: ErrorTypeValidation, Message: "registry not found"}
	}
	return &config, nil
}

func (m *mockRegistryManager) ListRegistries() ([]RegistryConfig, error) {
	var configs []RegistryConfig
	for _, config := range m.registries {
		configs = append(configs, config)
	}
	return configs, nil
}

func (m *mockRegistryManager) SetDefault(name string) error {
	m.defaultRegistry = name
	return nil
}

func (m *mockRegistryManager) GetDefault() (*RegistryConfig, error) {
	if m.defaultRegistry == "" {
		return nil, &CLIError{Type: ErrorTypeValidation, Message: "no default registry set"}
	}
	return m.GetRegistry(m.defaultRegistry)
}

func (m *mockRegistryManager) ValidateRegistry(ctx context.Context, name string) error {
	m.validated[name] = true
	return nil
}

// mockRegistryAuthenticator implements RegistryAuthenticator for testing
type mockRegistryAuthenticator struct {
	authenticated bool
	refreshed     bool
	loggedOut     bool
	authToken     string
	supportedAuth []string
}

func (m *mockRegistryAuthenticator) Authenticate(ctx context.Context, config RegistryConfig) error {
	m.authenticated = true
	return nil
}

func (m *mockRegistryAuthenticator) RefreshCredentials(ctx context.Context) error {
	m.refreshed = true
	return nil
}

func (m *mockRegistryAuthenticator) GetAuthToken(ctx context.Context) (string, error) {
	return m.authToken, nil
}

func (m *mockRegistryAuthenticator) IsAuthenticated() bool {
	return m.authenticated
}

func (m *mockRegistryAuthenticator) GetSupportedAuthTypes() []string {
	return m.supportedAuth
}

func (m *mockRegistryAuthenticator) Logout() error {
	m.loggedOut = true
	m.authenticated = false
	return nil
}

// mockRegistryDiscovery implements RegistryDiscovery for testing
type mockRegistryDiscovery struct {
	discoveredInfo   *RegistryInfo
	probedCaps       []string
	detectedType     string
	recommendedConfig *RegistryConfig
}

func (m *mockRegistryDiscovery) DiscoverRegistry(ctx context.Context, url string) (*RegistryInfo, error) {
	return m.discoveredInfo, nil
}

func (m *mockRegistryDiscovery) ProbeCapabilities(ctx context.Context, url string) ([]string, error) {
	return m.probedCaps, nil
}

func (m *mockRegistryDiscovery) DetectRegistryType(ctx context.Context, url string) (string, error) {
	return m.detectedType, nil
}

func (m *mockRegistryDiscovery) GetRecommendedConfig(ctx context.Context, url string) (*RegistryConfig, error) {
	return m.recommendedConfig, nil
}

// mockRegistryHealthChecker implements RegistryHealthChecker for testing
type mockRegistryHealthChecker struct {
	healthStatus     *HealthStatus
	connectivityOK   bool
	authenticationOK bool
	healthHistory    []HealthStatus
}

func (m *mockRegistryHealthChecker) CheckHealth(ctx context.Context, config RegistryConfig) (*HealthStatus, error) {
	return m.healthStatus, nil
}

func (m *mockRegistryHealthChecker) CheckConnectivity(ctx context.Context, url string) error {
	m.connectivityOK = true
	return nil
}

func (m *mockRegistryHealthChecker) CheckAuthentication(ctx context.Context, config RegistryConfig) error {
	m.authenticationOK = true
	return nil
}

func (m *mockRegistryHealthChecker) GetHealthHistory(registryName string) ([]HealthStatus, error) {
	return m.healthHistory, nil
}

// mockRegistryCatalogBrowser implements RegistryCatalogBrowser for testing
type mockRegistryCatalogBrowser struct {
	repositories []string
	tags         []string
	repoInfo     *RepositoryInfo
	searchResults []RepositoryInfo
	catalogInfo   *CatalogInfo
}

func (m *mockRegistryCatalogBrowser) ListRepositories(ctx context.Context, opts ListOptions) ([]string, error) {
	return m.repositories, nil
}

func (m *mockRegistryCatalogBrowser) ListTags(ctx context.Context, repository string, opts ListOptions) ([]string, error) {
	return m.tags, nil
}

func (m *mockRegistryCatalogBrowser) GetRepositoryInfo(ctx context.Context, repository string) (*RepositoryInfo, error) {
	return m.repoInfo, nil
}

func (m *mockRegistryCatalogBrowser) SearchRepositories(ctx context.Context, query string, opts ListOptions) ([]RepositoryInfo, error) {
	return m.searchResults, nil
}

func (m *mockRegistryCatalogBrowser) GetCatalogInfo(ctx context.Context) (*CatalogInfo, error) {
	return m.catalogInfo, nil
}

// Test RegistryCommand interface compliance
func TestRegistryCommandInterface(t *testing.T) {
	cmd := &mockRegistryCommand{
		registryInfo: &RegistryInfo{
			Name:    "test-registry",
			URL:     "https://test.registry.com",
			Version: "v2.0",
		},
		capabilities: []string{"push", "pull", "list"},
	}

	// Test interface compliance
	var _ RegistryCommand = cmd

	ctx := context.Background()
	config := RegistryConfig{
		Name: "test",
		URL:  "https://test.registry.com",
		Type: "generic",
	}

	// Test ConfigureRegistry
	err := cmd.ConfigureRegistry(config)
	if err != nil {
		t.Errorf("ConfigureRegistry failed: %v", err)
	}
	if !cmd.configured {
		t.Error("ConfigureRegistry was not called")
	}

	// Test ValidateCredentials
	err = cmd.ValidateCredentials(ctx)
	if err != nil {
		t.Errorf("ValidateCredentials failed: %v", err)
	}
	if !cmd.credentialsValid {
		t.Error("ValidateCredentials was not called")
	}

	// Test GetRegistryInfo
	info, err := cmd.GetRegistryInfo(ctx)
	if err != nil {
		t.Errorf("GetRegistryInfo failed: %v", err)
	}
	if info == nil || info.Name != "test-registry" {
		t.Error("GetRegistryInfo did not return expected info")
	}

	// Test TestConnection
	err = cmd.TestConnection(ctx)
	if err != nil {
		t.Errorf("TestConnection failed: %v", err)
	}
	if !cmd.connectionTested {
		t.Error("TestConnection was not called")
	}

	// Test GetCapabilities
	caps, err := cmd.GetCapabilities(ctx)
	if err != nil {
		t.Errorf("GetCapabilities failed: %v", err)
	}
	if len(caps) != 3 {
		t.Errorf("Expected 3 capabilities, got %d", len(caps))
	}
}

// Test RegistryManager interface compliance
func TestRegistryManagerInterface(t *testing.T) {
	manager := newMockRegistryManager()

	// Test interface compliance
	var _ RegistryManager = manager

	ctx := context.Background()
	config := RegistryConfig{
		Name: "test",
		URL:  "https://test.registry.com",
		Type: "dockerhub",
	}

	// Test AddRegistry
	err := manager.AddRegistry("test", config)
	if err != nil {
		t.Errorf("AddRegistry failed: %v", err)
	}

	// Test GetRegistry
	retrieved, err := manager.GetRegistry("test")
	if err != nil {
		t.Errorf("GetRegistry failed: %v", err)
	}
	if retrieved.Name != "test" {
		t.Error("GetRegistry did not return correct config")
	}

	// Test UpdateRegistry
	updatedConfig := config
	updatedConfig.URL = "https://updated.registry.com"
	err = manager.UpdateRegistry("test", updatedConfig)
	if err != nil {
		t.Errorf("UpdateRegistry failed: %v", err)
	}

	// Test SetDefault and GetDefault
	err = manager.SetDefault("test")
	if err != nil {
		t.Errorf("SetDefault failed: %v", err)
	}

	defaultConfig, err := manager.GetDefault()
	if err != nil {
		t.Errorf("GetDefault failed: %v", err)
	}
	if defaultConfig.Name != "test" {
		t.Error("GetDefault did not return correct config")
	}

	// Test ListRegistries
	registries, err := manager.ListRegistries()
	if err != nil {
		t.Errorf("ListRegistries failed: %v", err)
	}
	if len(registries) != 1 {
		t.Errorf("Expected 1 registry, got %d", len(registries))
	}

	// Test ValidateRegistry
	err = manager.ValidateRegistry(ctx, "test")
	if err != nil {
		t.Errorf("ValidateRegistry failed: %v", err)
	}
	if !manager.validated["test"] {
		t.Error("ValidateRegistry was not called")
	}

	// Test RemoveRegistry
	err = manager.RemoveRegistry("test")
	if err != nil {
		t.Errorf("RemoveRegistry failed: %v", err)
	}

	registries, err = manager.ListRegistries()
	if err != nil {
		t.Errorf("ListRegistries after remove failed: %v", err)
	}
	if len(registries) != 0 {
		t.Errorf("Expected 0 registries after remove, got %d", len(registries))
	}
}

// Test RegistryAuthenticator interface compliance
func TestRegistryAuthenticatorInterface(t *testing.T) {
	auth := &mockRegistryAuthenticator{
		authToken:     "test-token",
		supportedAuth: []string{"basic", "token"},
	}

	// Test interface compliance
	var _ RegistryAuthenticator = auth

	ctx := context.Background()
	config := RegistryConfig{AuthType: "basic"}

	// Test Authenticate
	err := auth.Authenticate(ctx, config)
	if err != nil {
		t.Errorf("Authenticate failed: %v", err)
	}
	if !auth.authenticated {
		t.Error("Authenticate was not called")
	}

	// Test IsAuthenticated
	if !auth.IsAuthenticated() {
		t.Error("IsAuthenticated should return true after authentication")
	}

	// Test GetAuthToken
	token, err := auth.GetAuthToken(ctx)
	if err != nil {
		t.Errorf("GetAuthToken failed: %v", err)
	}
	if token != "test-token" {
		t.Errorf("Expected 'test-token', got '%s'", token)
	}

	// Test RefreshCredentials
	err = auth.RefreshCredentials(ctx)
	if err != nil {
		t.Errorf("RefreshCredentials failed: %v", err)
	}
	if !auth.refreshed {
		t.Error("RefreshCredentials was not called")
	}

	// Test GetSupportedAuthTypes
	types := auth.GetSupportedAuthTypes()
	if len(types) != 2 {
		t.Errorf("Expected 2 auth types, got %d", len(types))
	}

	// Test Logout
	err = auth.Logout()
	if err != nil {
		t.Errorf("Logout failed: %v", err)
	}
	if !auth.loggedOut {
		t.Error("Logout was not called")
	}
	if auth.IsAuthenticated() {
		t.Error("IsAuthenticated should return false after logout")
	}
}

// Test RegistryDiscovery interface compliance
func TestRegistryDiscoveryInterface(t *testing.T) {
	discovery := &mockRegistryDiscovery{
		discoveredInfo: &RegistryInfo{Name: "discovered", Version: "v2.1"},
		probedCaps:     []string{"push", "pull"},
		detectedType:   "harbor",
		recommendedConfig: &RegistryConfig{
			Type:     "harbor",
			AuthType: "basic",
		},
	}

	// Test interface compliance
	var _ RegistryDiscovery = discovery

	ctx := context.Background()
	url := "https://registry.example.com"

	// Test DiscoverRegistry
	info, err := discovery.DiscoverRegistry(ctx, url)
	if err != nil {
		t.Errorf("DiscoverRegistry failed: %v", err)
	}
	if info.Name != "discovered" {
		t.Error("DiscoverRegistry did not return expected info")
	}

	// Test ProbeCapabilities
	caps, err := discovery.ProbeCapabilities(ctx, url)
	if err != nil {
		t.Errorf("ProbeCapabilities failed: %v", err)
	}
	if len(caps) != 2 {
		t.Errorf("Expected 2 capabilities, got %d", len(caps))
	}

	// Test DetectRegistryType
	regType, err := discovery.DetectRegistryType(ctx, url)
	if err != nil {
		t.Errorf("DetectRegistryType failed: %v", err)
	}
	if regType != "harbor" {
		t.Errorf("Expected 'harbor', got '%s'", regType)
	}

	// Test GetRecommendedConfig
	config, err := discovery.GetRecommendedConfig(ctx, url)
	if err != nil {
		t.Errorf("GetRecommendedConfig failed: %v", err)
	}
	if config.Type != "harbor" {
		t.Error("GetRecommendedConfig did not return expected config")
	}
}

// Test RegistryHealthChecker interface compliance
func TestRegistryHealthCheckerInterface(t *testing.T) {
	healthChecker := &mockRegistryHealthChecker{
		healthStatus: &HealthStatus{
			RegistryName: "test",
			Status:       "healthy",
			ResponseTime: 100,
		},
		healthHistory: []HealthStatus{
			{Status: "healthy", ResponseTime: 95},
			{Status: "healthy", ResponseTime: 105},
		},
	}

	// Test interface compliance
	var _ RegistryHealthChecker = healthChecker

	ctx := context.Background()
	config := RegistryConfig{Name: "test"}

	// Test CheckHealth
	health, err := healthChecker.CheckHealth(ctx, config)
	if err != nil {
		t.Errorf("CheckHealth failed: %v", err)
	}
	if health.Status != "healthy" {
		t.Error("CheckHealth did not return expected status")
	}

	// Test CheckConnectivity
	err = healthChecker.CheckConnectivity(ctx, "https://test.com")
	if err != nil {
		t.Errorf("CheckConnectivity failed: %v", err)
	}
	if !healthChecker.connectivityOK {
		t.Error("CheckConnectivity was not called")
	}

	// Test CheckAuthentication
	err = healthChecker.CheckAuthentication(ctx, config)
	if err != nil {
		t.Errorf("CheckAuthentication failed: %v", err)
	}
	if !healthChecker.authenticationOK {
		t.Error("CheckAuthentication was not called")
	}

	// Test GetHealthHistory
	history, err := healthChecker.GetHealthHistory("test")
	if err != nil {
		t.Errorf("GetHealthHistory failed: %v", err)
	}
	if len(history) != 2 {
		t.Errorf("Expected 2 history entries, got %d", len(history))
	}
}

// Test RegistryCatalogBrowser interface compliance
func TestRegistryCatalogBrowserInterface(t *testing.T) {
	browser := &mockRegistryCatalogBrowser{
		repositories: []string{"repo1", "repo2"},
		tags:         []string{"v1.0", "v2.0"},
		repoInfo: &RepositoryInfo{
			Name:     "test-repo",
			TagCount: 5,
		},
		searchResults: []RepositoryInfo{
			{Name: "found1"},
			{Name: "found2"},
		},
		catalogInfo: &CatalogInfo{
			RepositoryCount: 10,
			TotalSize:       1024000,
		},
	}

	// Test interface compliance
	var _ RegistryCatalogBrowser = browser

	ctx := context.Background()
	opts := ListOptions{Limit: 10}

	// Test ListRepositories
	repos, err := browser.ListRepositories(ctx, opts)
	if err != nil {
		t.Errorf("ListRepositories failed: %v", err)
	}
	if len(repos) != 2 {
		t.Errorf("Expected 2 repositories, got %d", len(repos))
	}

	// Test ListTags
	tags, err := browser.ListTags(ctx, "test-repo", opts)
	if err != nil {
		t.Errorf("ListTags failed: %v", err)
	}
	if len(tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(tags))
	}

	// Test GetRepositoryInfo
	info, err := browser.GetRepositoryInfo(ctx, "test-repo")
	if err != nil {
		t.Errorf("GetRepositoryInfo failed: %v", err)
	}
	if info.Name != "test-repo" {
		t.Error("GetRepositoryInfo did not return expected info")
	}

	// Test SearchRepositories
	results, err := browser.SearchRepositories(ctx, "test", opts)
	if err != nil {
		t.Errorf("SearchRepositories failed: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("Expected 2 search results, got %d", len(results))
	}

	// Test GetCatalogInfo
	catalog, err := browser.GetCatalogInfo(ctx)
	if err != nil {
		t.Errorf("GetCatalogInfo failed: %v", err)
	}
	if catalog.RepositoryCount != 10 {
		t.Error("GetCatalogInfo did not return expected info")
	}
}

// Test CLIError type functionality
func TestCLIErrorType(t *testing.T) {
	// Test error creation
	err := &CLIError{
		Type:    ErrorTypeValidation,
		Message: "invalid input",
		Context: "command validation",
	}

	// Test Error() method
	errMsg := err.Error()
	expected := "command validation: invalid input"
	if errMsg != expected {
		t.Errorf("Expected '%s', got '%s'", expected, errMsg)
	}

	// Test error without context
	err2 := &CLIError{
		Type:    ErrorTypeNetwork,
		Message: "connection failed",
	}

	errMsg2 := err2.Error()
	if errMsg2 != "connection failed" {
		t.Errorf("Expected 'connection failed', got '%s'", errMsg2)
	}

	// Test Unwrap method
	underlyingErr := &CLIError{Type: ErrorTypeRegistry, Message: "registry error"}
	err3 := &CLIError{
		Type:  ErrorTypeNetwork,
		Cause: underlyingErr,
	}

	if err3.Unwrap() != underlyingErr {
		t.Error("Unwrap did not return the expected underlying error")
	}
}