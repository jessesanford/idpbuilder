package fallback

import (
	"net/http"
	"testing"
)

func TestNewInsecureConfig(t *testing.T) {
	config := NewInsecureConfig()

	if config == nil {
		t.Fatal("NewInsecureConfig returned nil")
	}

	if config.Enabled {
		t.Error("Expected insecure mode to be disabled by default")
	}

	if config.WarningShown {
		t.Error("Expected warning not shown by default")
	}

	if len(config.AuditLog) != 0 {
		t.Error("Expected empty audit log by default")
	}

	if config.logger == nil {
		t.Error("Expected logger to be initialized")
	}
}

func TestCreateInsecureTLSConfig(t *testing.T) {
	tlsConfig := CreateInsecureTLSConfig()

	if tlsConfig == nil {
		t.Fatal("CreateInsecureTLSConfig returned nil")
	}

	if !tlsConfig.InsecureSkipVerify {
		t.Error("Expected InsecureSkipVerify to be true")
	}
}

func TestApplyInsecureFlag_Enable(t *testing.T) {
	config := NewInsecureConfig()
	registryURL := "https://test-registry.com"

	err := config.ApplyInsecureFlag(true, registryURL)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !config.Enabled {
		t.Error("Expected insecure mode to be enabled")
	}

	if !config.WarningShown {
		t.Error("Expected warning to be shown")
	}

	if len(config.AuditLog) != 1 {
		t.Errorf("Expected 1 audit log entry, got: %d", len(config.AuditLog))
	}
}

func TestApplyInsecureFlag_Disable(t *testing.T) {
	config := NewInsecureConfig()

	// First enable it
	config.ApplyInsecureFlag(true, "test")

	// Then disable it
	err := config.ApplyInsecureFlag(false, "test")

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if config.Enabled {
		t.Error("Expected insecure mode to be disabled")
	}
}

func TestIsInsecureModeEnabled(t *testing.T) {
	config := NewInsecureConfig()

	if config.IsInsecureModeEnabled() {
		t.Error("Expected insecure mode to be disabled initially")
	}

	config.ApplyInsecureFlag(true, "test")

	if !config.IsInsecureModeEnabled() {
		t.Error("Expected insecure mode to be enabled after applying flag")
	}
}

func TestWrapTransportWithInsecure_NotInsecure(t *testing.T) {
	originalTransport := &http.Transport{}

	wrappedTransport := WrapTransportWithInsecure(originalTransport, false)

	if wrappedTransport != originalTransport {
		t.Error("Expected transport to be unchanged when not insecure")
	}
}

func TestWrapTransportWithInsecure_Insecure(t *testing.T) {
	originalTransport := &http.Transport{}

	wrappedTransport := WrapTransportWithInsecure(originalTransport, true)

	if wrappedTransport == originalTransport {
		t.Error("Expected transport to be modified when insecure")
	}

	// Check that the wrapped transport has insecure TLS config
	if httpTransport, ok := wrappedTransport.(*http.Transport); ok {
		if httpTransport.TLSClientConfig == nil {
			t.Error("Expected TLS client config to be set")
		}

		if !httpTransport.TLSClientConfig.InsecureSkipVerify {
			t.Error("Expected InsecureSkipVerify to be true")
		}
	} else {
		t.Error("Expected wrapped transport to be *http.Transport")
	}
}

func TestCreateInsecureHTTPClient(t *testing.T) {
	client := CreateInsecureHTTPClient()

	if client == nil {
		t.Fatal("CreateInsecureHTTPClient returned nil")
	}

	if client.Timeout == 0 {
		t.Error("Expected timeout to be set")
	}

	// Check that the transport has insecure TLS config
	if httpTransport, ok := client.Transport.(*http.Transport); ok {
		if httpTransport.TLSClientConfig == nil {
			t.Error("Expected TLS client config to be set")
		}

		if !httpTransport.TLSClientConfig.InsecureSkipVerify {
			t.Error("Expected InsecureSkipVerify to be true")
		}
	} else {
		t.Error("Expected client transport to be *http.Transport")
	}
}

func TestLogInsecureConnection(t *testing.T) {
	config := NewInsecureConfig()
	registryURL := "https://test-registry.com"

	config.LogInsecureConnection(registryURL)

	if len(config.AuditLog) != 1 {
		t.Errorf("Expected 1 audit log entry, got: %d", len(config.AuditLog))
	}

	logEntry := config.AuditLog[0]
	if !contains(logEntry, registryURL) {
		t.Error("Expected audit log to contain registry URL")
	}

	if !contains(logEntry, "Insecure connection established") {
		t.Error("Expected audit log to contain connection message")
	}
}

func TestValidateInsecureUsage_NotEnabled(t *testing.T) {
	config := NewInsecureConfig()

	err := config.ValidateInsecureUsage("https://production.com")

	if err != nil {
		t.Errorf("Expected no error when insecure mode not enabled, got: %v", err)
	}
}

func TestValidateInsecureUsage_EmptyURL(t *testing.T) {
	config := NewInsecureConfig()
	config.ApplyInsecureFlag(true, "test")

	err := config.ValidateInsecureUsage("")

	if err == nil {
		t.Error("Expected error for empty registry URL")
	}
}

func TestContainsProductionIndicators(t *testing.T) {
	tests := []struct {
		url      string
		expected bool
	}{
		{"https://production-registry.com", true},
		{"https://registry.prod.example.org", true},
		{"https://docker.io", true},
		{"https://gcr.io", true},
		{"https://localhost:5000", false},
		{"https://registry.example.dev", false},
		{"https://test-registry", false},
	}

	for _, test := range tests {
		result := containsProductionIndicators(test.url)
		if result != test.expected {
			t.Errorf("For URL %s, expected: %v, got: %v", test.url, test.expected, result)
		}
	}
}
