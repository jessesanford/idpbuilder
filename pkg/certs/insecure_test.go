package certs

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Test helper: creates a test server with self-signed certificate
func createTestTLSServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	if handler == nil {
		handler = func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("test response"))
		}
	}
	
	server := httptest.NewTLSServer(handler)
	t.Cleanup(func() {
		server.Close()
	})
	return server
}

// Test helper: creates basic test configuration
func createTestConfig(t *testing.T) *InsecureConfig {
	return &InsecureConfig{
		Enabled:           true,
		ShowWarnings:      false, // Disable warnings in tests
		AuditConnections:  true,
		AllowedRegistries: []string{},
		MaxAuditEntries:   10,
		auditEntries:      make([]AuditEntry, 0),
	}
}

// Test helper: creates test request
func createTestRequest(t *testing.T, method, url string) *http.Request {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		t.Fatalf("Failed to create test request: %v", err)
	}
	return req
}

// TestNewInsecureTransport tests the creation of a new insecure transport
func TestNewInsecureTransport(t *testing.T) {
	transport := NewInsecureTransport()
	
	if transport == nil {
		t.Fatal("NewInsecureTransport returned nil")
	}
	
	if transport.Base == nil {
		t.Error("Base transport should not be nil")
	}
	
	if transport.config == nil {
		t.Error("Config should not be nil")
	}
	
	if transport.auditLogger == nil {
		t.Error("Audit logger should not be nil")
	}
	
	// Test default configuration
	if transport.config.Enabled {
		t.Error("Insecure mode should be disabled by default")
	}
	
	if !transport.config.ShowWarnings {
		t.Error("ShowWarnings should be enabled by default")
	}
	
	if !transport.config.AuditConnections {
		t.Error("AuditConnections should be enabled by default")
	}
	
	if transport.config.MaxAuditEntries != 100 {
		t.Errorf("Expected MaxAuditEntries to be 100, got %d", transport.config.MaxAuditEntries)
	}
	
	if len(transport.config.AllowedRegistries) != 0 {
		t.Error("AllowedRegistries should be empty by default")
	}
	
	if len(transport.config.auditEntries) != 0 {
		t.Error("Audit entries should be empty by default")
	}
}

// TestInsecureRoundTrip_Success tests successful insecure round trip
func TestInsecureRoundTrip_Success(t *testing.T) {
	// Create a test TLS server
	server := createTestTLSServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success response"))
	})
	
	// Create insecure transport with test configuration
	transport := NewInsecureTransportWithConfig(createTestConfig(t))
	
	// Enable insecure mode
	transport.EnableInsecure([]string{})
	
	// Create test request
	req := createTestRequest(t, "GET", server.URL)
	
	// Make the request
	resp, err := transport.RoundTrip(req)
	if err != nil {
		t.Fatalf("RoundTrip failed: %v", err)
	}
	
	if resp == nil {
		t.Fatal("Response should not be nil")
	}
	
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	
	// Verify audit entry was created
	auditEntries := transport.GetAuditEntries()
	if len(auditEntries) == 0 {
		t.Error("Expected audit entries to be created")
	}
	
	// Check that insecure connection was audited
	foundAttempt := false
	foundSuccess := false
	for _, entry := range auditEntries {
		if strings.Contains(entry.Status, "ATTEMPTING") {
			foundAttempt = true
		}
		if strings.Contains(entry.Status, "HTTP_200") {
			foundSuccess = true
		}
	}
	
	if !foundAttempt {
		t.Error("Expected to find ATTEMPTING audit entry")
	}
	
	if !foundSuccess {
		t.Error("Expected to find success audit entry")
	}
	
	resp.Body.Close()
}

// TestInsecureRoundTrip_BasicError tests error handling in insecure round trip
func TestInsecureRoundTrip_BasicError(t *testing.T) {
	// Create insecure transport
	transport := NewInsecureTransport()
	
	// Enable insecure mode
	transport.EnableInsecure([]string{})
	
	// Create request to non-existent server
	req := createTestRequest(t, "GET", "https://non-existent-server.invalid:9999/test")
	
	// Make the request (should fail)
	resp, err := transport.RoundTrip(req)
	
	// Should get an error
	if err == nil {
		t.Error("Expected error for non-existent server")
		if resp != nil {
			resp.Body.Close()
		}
	}
	
	// Verify audit entry was created for the error
	auditEntries := transport.GetAuditEntries()
	if len(auditEntries) == 0 {
		t.Error("Expected audit entries even for errors")
	}
	
	// Check that error was audited
	foundError := false
	for _, entry := range auditEntries {
		if strings.Contains(entry.Status, "ERROR") {
			foundError = true
			break
		}
	}
	
	if !foundError {
		t.Error("Expected to find error audit entry")
	}
}