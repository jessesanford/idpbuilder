package push

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
)

// MockRegistry provides a mock OCI registry for unit testing
type MockRegistry struct {
	Server     *httptest.Server  // Test server instance
	Images     map[string][]byte // Stored images keyed by name:tag
	AuthConfig *AuthConfig       // Optional auth configuration
}

// AuthConfig configures authentication requirements for the mock registry
type AuthConfig struct {
	Username string // Expected username
	Password string // Expected password
	Required bool   // Whether auth is required
}

// MockAuthTransport wraps an HTTP transport with authentication
type MockAuthTransport struct {
	Username string              // Auth username
	Password string              // Auth password
	Base     http.RoundTripper   // Base transport
}

// NewMockRegistry creates a new mock OCI registry for testing
func NewMockRegistry() *MockRegistry {
	reg := &MockRegistry{
		Images: make(map[string][]byte),
	}

	// Set up HTTP test server with registry API handlers
	mux := http.NewServeMux()

	// Handle manifest uploads (PUT /v2/{name}/manifests/{tag})
	mux.HandleFunc("/v2/", func(w http.ResponseWriter, r *http.Request) {
		// Check authentication if required
		if reg.AuthConfig != nil && reg.AuthConfig.Required {
			if !reg.checkAuth(r) {
				w.Header().Set("WWW-Authenticate", `Bearer realm="https://auth.docker.io/token"`)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		path := strings.TrimPrefix(r.URL.Path, "/v2/")

		switch r.Method {
		case "GET":
			if strings.HasSuffix(path, "/manifests/") {
				// Return stored manifest
				name := strings.TrimSuffix(path, "/manifests/")
				if data, exists := reg.Images[name]; exists {
					w.Header().Set("Content-Type", "application/vnd.docker.distribution.manifest.v2+json")
					w.Write(data)
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
			} else {
				// Registry version check
				w.Header().Set("Docker-Distribution-API-Version", "registry/2.0")
				w.WriteHeader(http.StatusOK)
			}
		case "PUT":
			if strings.Contains(path, "/manifests/") {
				// Store manifest
				parts := strings.Split(path, "/manifests/")
				if len(parts) == 2 {
					name := parts[0]
					tag := parts[1]
					key := fmt.Sprintf("%s:%s", name, tag)

					// Read manifest data
					buf := make([]byte, r.ContentLength)
					r.Body.Read(buf)
					reg.Images[key] = buf

					w.WriteHeader(http.StatusCreated)
				}
			}
		case "HEAD":
			// Check if manifest exists
			if strings.Contains(path, "/manifests/") {
				parts := strings.Split(path, "/manifests/")
				if len(parts) == 2 {
					name := parts[0]
					tag := parts[1]
					key := fmt.Sprintf("%s:%s", name, tag)
					if _, exists := reg.Images[key]; exists {
						w.WriteHeader(http.StatusOK)
					} else {
						w.WriteHeader(http.StatusNotFound)
					}
				}
			}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	reg.Server = httptest.NewServer(mux)
	return reg
}

// NewMockAuthTransport creates a transport that adds authentication headers
func NewMockAuthTransport(username, password string) *MockAuthTransport {
	return &MockAuthTransport{
		Username: username,
		Password: password,
		Base:     http.DefaultTransport,
	}
}

// RoundTrip implements http.RoundTripper interface with authentication
func (t *MockAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request to avoid modifying the original
	clonedReq := req.Clone(req.Context())

	// Add basic authentication
	if t.Username != "" && t.Password != "" {
		clonedReq.SetBasicAuth(t.Username, t.Password)
	}

	// Use base transport to execute the request
	if t.Base != nil {
		return t.Base.RoundTrip(clonedReq)
	}
	return http.DefaultTransport.RoundTrip(clonedReq)
}

// SetAuth configures authentication requirements for the mock registry
func (r *MockRegistry) SetAuth(username, password string, required bool) {
	r.AuthConfig = &AuthConfig{
		Username: username,
		Password: password,
		Required: required,
	}
}

// Cleanup shuts down the mock registry server
func (r *MockRegistry) Cleanup() {
	if r.Server != nil {
		r.Server.Close()
	}
}

// GetURL returns the base URL of the mock registry
func (r *MockRegistry) GetURL() string {
	if r.Server != nil {
		return r.Server.URL
	}
	return ""
}

// GetStoredImages returns a copy of stored images for testing verification
func (r *MockRegistry) GetStoredImages() map[string][]byte {
	images := make(map[string][]byte)
	for key, value := range r.Images {
		images[key] = value
	}
	return images
}

// checkAuth validates authentication from HTTP request
func (r *MockRegistry) checkAuth(req *http.Request) bool {
	if r.AuthConfig == nil || !r.AuthConfig.Required {
		return true
	}

	username, password, ok := req.BasicAuth()
	if !ok {
		return false
	}

	return username == r.AuthConfig.Username && password == r.AuthConfig.Password
}

// MockInsecureTransport creates a transport that ignores TLS certificate errors
func MockInsecureTransport() http.RoundTripper {
	return &mockInsecureTransport{
		base: http.DefaultTransport,
	}
}

// mockInsecureTransport implements http.RoundTripper for insecure connections
type mockInsecureTransport struct {
	base http.RoundTripper
}

// RoundTrip executes requests while ignoring certificate validation
func (t *mockInsecureTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// For testing purposes, we simply pass through to base transport
	// In real scenarios, this would configure TLS to skip verification
	return t.base.RoundTrip(req)
}

// CreateRegistryError creates a mock registry error response
func CreateRegistryError(status int, code, message string) []byte {
	errorResp := map[string]interface{}{
		"errors": []map[string]interface{}{
			{
				"code":    code,
				"message": message,
			},
		},
	}

	data, _ := json.Marshal(errorResp)
	return data
}