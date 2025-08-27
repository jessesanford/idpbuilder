package contexts

import (
	"net/http"
	"testing"
	"time"
)

func TestContextType_String(t *testing.T) {
	tests := []struct {
		name     string
		ct       ContextType
		expected string
	}{
		{"Local context", LocalContext, "local"},
		{"URL context", URLContext, "url"},
		{"Git context", GitContext, "git"},
		{"Archive context", ArchiveContext, "archive"},
		{"Unknown context", ContextType(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ct.String(); got != tt.expected {
				t.Errorf("ContextType.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config == nil {
		t.Fatal("DefaultConfig() returned nil")
	}

	expectedMaxSize := int64(500 * 1024 * 1024) // 500MB
	if config.MaxSize != expectedMaxSize {
		t.Errorf("DefaultConfig().MaxSize = %v, want %v", config.MaxSize, expectedMaxSize)
	}

	if !config.CacheEnabled {
		t.Error("DefaultConfig().CacheEnabled = false, want true")
	}

	if config.TempDir != "/tmp" {
		t.Errorf("DefaultConfig().TempDir = %v, want %v", config.TempDir, "/tmp")
	}

	expectedTimeout := 30 * time.Second
	if config.HTTPTimeout != expectedTimeout {
		t.Errorf("DefaultConfig().HTTPTimeout = %v, want %v", config.HTTPTimeout, expectedTimeout)
	}
}

func TestContextConfig_DefaultValues(t *testing.T) {
	config := &ContextConfig{}

	// Test that zero values are reasonable
	if config.MaxSize != 0 {
		t.Errorf("Empty ContextConfig.MaxSize = %v, want 0", config.MaxSize)
	}

	if config.CacheEnabled {
		t.Error("Empty ContextConfig.CacheEnabled = true, want false")
	}

	if config.TempDir != "" {
		t.Errorf("Empty ContextConfig.TempDir = %v, want empty string", config.TempDir)
	}

	if config.HTTPTimeout != 0 {
		t.Errorf("Empty ContextConfig.HTTPTimeout = %v, want 0", config.HTTPTimeout)
	}
}

func TestContextConfig_CustomValues(t *testing.T) {
	customMaxSize := int64(100 * 1024 * 1024) // 100MB
	customTempDir := "/custom/temp"
	customTimeout := 60 * time.Second

	config := &ContextConfig{
		MaxSize:      customMaxSize,
		CacheEnabled: false,
		TempDir:      customTempDir,
		HTTPTimeout:  customTimeout,
	}

	if config.MaxSize != customMaxSize {
		t.Errorf("ContextConfig.MaxSize = %v, want %v", config.MaxSize, customMaxSize)
	}

	if config.CacheEnabled {
		t.Error("ContextConfig.CacheEnabled = true, want false")
	}

	if config.TempDir != customTempDir {
		t.Errorf("ContextConfig.TempDir = %v, want %v", config.TempDir, customTempDir)
	}

	if config.HTTPTimeout != customTimeout {
		t.Errorf("ContextConfig.HTTPTimeout = %v, want %v", config.HTTPTimeout, customTimeout)
	}
}

func TestAuthConfig(t *testing.T) {
	tests := []struct {
		name     string
		config   AuthConfig
		username string
		password string
		token    string
	}{
		{
			name:     "Empty auth config",
			config:   AuthConfig{},
			username: "",
			password: "",
			token:    "",
		},
		{
			name: "Username/password auth",
			config: AuthConfig{
				Username: "testuser",
				Password: "testpass",
			},
			username: "testuser",
			password: "testpass",
			token:    "",
		},
		{
			name: "Token auth",
			config: AuthConfig{
				Token: "test-token-123",
			},
			username: "",
			password: "",
			token:    "test-token-123",
		},
		{
			name: "Full auth config",
			config: AuthConfig{
				Username: "user",
				Password: "pass",
				Token:    "token",
			},
			username: "user",
			password: "pass",
			token:    "token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.config.Username != tt.username {
				t.Errorf("AuthConfig.Username = %v, want %v", tt.config.Username, tt.username)
			}
			if tt.config.Password != tt.password {
				t.Errorf("AuthConfig.Password = %v, want %v", tt.config.Password, tt.password)
			}
			if tt.config.Token != tt.token {
				t.Errorf("AuthConfig.Token = %v, want %v", tt.config.Token, tt.token)
			}
		})
	}
}

// Mock HTTPClient for testing
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if m.DoFunc != nil {
		return m.DoFunc(req)
	}
	return &http.Response{StatusCode: http.StatusOK}, nil
}

func TestHTTPClientInterface(t *testing.T) {
	// Test that our mock implements the interface
	var client HTTPClient = &MockHTTPClient{}
	
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("HTTPClient.Do() error = %v", err)
	}

	if resp == nil {
		t.Error("HTTPClient.Do() returned nil response")
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("HTTPClient.Do() status = %v, want %v", resp.StatusCode, http.StatusOK)
	}
}

func TestContextInterface(t *testing.T) {
	// Test that we can define the interface behavior expectations
	// This test documents the expected interface contract

	t.Run("Interface methods exist", func(t *testing.T) {
		// This is a compile-time test - if the interface methods don't exist,
		// this won't compile
		var ctx Context
		if ctx != nil {
			_ = ctx.Path()
			_ = ctx.Cleanup()
			_ = ctx.Type()
		}
	})
}