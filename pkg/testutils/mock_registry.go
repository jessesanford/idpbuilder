package testutils

import (
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
)

// AuthConfig represents authentication configuration for the mock registry
type AuthConfig struct {
	Username string
	Password string
	Token    string
	Enabled  bool
}

// MockAuthTransport implements http.RoundTripper interface for testing
type MockAuthTransport struct {
	config     *AuthConfig
	transport  http.RoundTripper
	authHeader string
	mu         sync.RWMutex
}

// NewMockAuthTransport creates a new mock authentication transport
func NewMockAuthTransport(config *AuthConfig) *MockAuthTransport {
	return &MockAuthTransport{
		config:    config,
		transport: http.DefaultTransport,
	}
}

// RoundTrip implements http.RoundTripper interface
func (m *MockAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.config.Enabled {
		if m.config.Token != "" {
			req.Header.Set("Authorization", "Bearer "+m.config.Token)
		} else if m.config.Username != "" && m.config.Password != "" {
			req.SetBasicAuth(m.config.Username, m.config.Password)
		}
	}

	return m.transport.RoundTrip(req)
}

// GetAuthHeaders returns the authorization headers for the configured auth
func (m *MockAuthTransport) GetAuthHeaders() map[string]string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	headers := make(map[string]string)
	if !m.config.Enabled {
		return headers
	}

	if m.config.Token != "" {
		headers["Authorization"] = "Bearer " + m.config.Token
	} else if m.config.Username != "" && m.config.Password != "" {
		// For basic auth, we would normally use base64 encoding,
		// but for testing purposes, we'll return the credentials
		headers["X-Username"] = m.config.Username
		headers["X-Password"] = m.config.Password
	}

	return headers
}

// MockRegistry provides a mock OCI registry for testing
type MockRegistry struct {
	server     *httptest.Server
	authConfig *AuthConfig
	manifests  map[string][]byte
	blobs      map[string][]byte
	tags       map[string]string
	mu         sync.RWMutex
	reqCount   int
	lastReq    *http.Request
}

// NewMockRegistry creates a new mock registry instance
func NewMockRegistry(authConfig *AuthConfig) *MockRegistry {
	registry := &MockRegistry{
		authConfig: authConfig,
		manifests:  make(map[string][]byte),
		blobs:      make(map[string][]byte),
		tags:       make(map[string]string),
	}

	mux := http.NewServeMux()

	// OCI Distribution API v2 endpoints - single handler for all /v2/ routes
	mux.HandleFunc("/v2/", registry.handleV2Endpoint)

	registry.server = httptest.NewServer(mux)
	return registry
}

// URL returns the registry URL
func (m *MockRegistry) URL() string {
	return m.server.URL
}

// Host returns the registry host
func (m *MockRegistry) Host() string {
	u, _ := url.Parse(m.server.URL)
	return u.Host
}

// Close shuts down the mock registry
func (m *MockRegistry) Close() {
	m.server.Close()
}

// GetRequestCount returns the number of requests handled
func (m *MockRegistry) GetRequestCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.reqCount
}

// GetLastRequest returns the last request received
func (m *MockRegistry) GetLastRequest() *http.Request {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.lastReq
}

// AddImageManifest adds a test image manifest to the registry
func (m *MockRegistry) AddImageManifest(ref string, manifest []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Generate digest for the manifest
	digest := generateDigest(manifest)
	m.manifests[digest] = manifest

	// Store tag reference
	if strings.Contains(ref, ":") {
		parts := strings.Split(ref, ":")
		if len(parts) == 2 {
			m.tags[parts[1]] = digest
		}
	}

	return nil
}

// AddBlob adds a blob to the registry
func (m *MockRegistry) AddBlob(digest string, data []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.blobs[digest] = data
	return nil
}

// handleV2Endpoint handles all /v2/ endpoints by routing to appropriate handlers
func (m *MockRegistry) handleV2Endpoint(w http.ResponseWriter, r *http.Request) {
	m.recordRequest(r)

	if !m.checkAuth(w, r) {
		return
	}

	path := r.URL.Path

	// Root endpoint
	if path == "/v2/" {
		w.Header().Set("Docker-Distribution-API-Version", "registry/2.0")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Route to specific handlers based on path
	if strings.Contains(path, "/manifests/") {
		m.handleManifests(w, r)
		return
	}

	if strings.Contains(path, "/blobs/") {
		m.handleBlobs(w, r)
		return
	}

	if strings.Contains(path, "/tags/list") {
		m.handleTags(w, r)
		return
	}

	http.NotFound(w, r)
}

// handleManifests handles manifest-related endpoints
func (m *MockRegistry) handleManifests(w http.ResponseWriter, r *http.Request) {

	// Extract reference (tag or digest)
	parts := strings.Split(r.URL.Path, "/manifests/")
	if len(parts) != 2 {
		http.Error(w, "Invalid manifest path", http.StatusBadRequest)
		return
	}

	ref := parts[1]

	switch r.Method {
	case http.MethodGet:
		m.getManifest(w, r, ref)
	case http.MethodPut:
		m.putManifest(w, r, ref)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleBlobs handles blob-related endpoints
func (m *MockRegistry) handleBlobs(w http.ResponseWriter, r *http.Request) {

	// Extract digest
	parts := strings.Split(r.URL.Path, "/blobs/")
	if len(parts) != 2 {
		http.Error(w, "Invalid blob path", http.StatusBadRequest)
		return
	}

	digest := parts[1]

	switch r.Method {
	case http.MethodGet:
		m.getBlob(w, r, digest)
	case http.MethodHead:
		m.headBlob(w, r, digest)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleTags handles tag listing endpoints
func (m *MockRegistry) handleTags(w http.ResponseWriter, r *http.Request) {

	// Return list of tags
	tags := make([]string, 0, len(m.tags))
	for tag := range m.tags {
		tags = append(tags, tag)
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"name":"%s","tags":[`, extractRepoName(r.URL.Path))
	for i, tag := range tags {
		if i > 0 {
			w.Write([]byte(","))
		}
		fmt.Fprintf(w, `"%s"`, tag)
	}
	w.Write([]byte("]}"))
}

// recordRequest records request for testing verification
func (m *MockRegistry) recordRequest(r *http.Request) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.reqCount++
	m.lastReq = r
}

// checkAuth verifies authentication if enabled
func (m *MockRegistry) checkAuth(w http.ResponseWriter, r *http.Request) bool {
	if !m.authConfig.Enabled {
		return true
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		w.Header().Set("WWW-Authenticate", "Basic realm=\"registry\"")
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}

	// Verify auth header matches expected credentials
	if m.authConfig.Token != "" {
		expected := "Bearer " + m.authConfig.Token
		if authHeader != expected {
			w.WriteHeader(http.StatusForbidden)
			return false
		}
	} else if m.authConfig.Username != "" {
		// Basic auth verification would go here
		// For testing purposes, we accept any basic auth
	}

	return true
}

// getManifest handles GET requests for manifests
func (m *MockRegistry) getManifest(w http.ResponseWriter, r *http.Request, ref string) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var manifest []byte
	var found bool

	// Check if ref is a tag
	if digest, ok := m.tags[ref]; ok {
		manifest, found = m.manifests[digest]
	} else {
		// Assume ref is a digest
		manifest, found = m.manifests[ref]
	}

	if !found {
		http.Error(w, "Manifest not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.docker.distribution.manifest.v2+json")
	w.Header().Set("Docker-Content-Digest", ref)
	w.Write(manifest)
}

// putManifest handles PUT requests for manifests
func (m *MockRegistry) putManifest(w http.ResponseWriter, r *http.Request, ref string) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read manifest", http.StatusBadRequest)
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Generate digest for the manifest
	digest := generateDigest(body)
	m.manifests[digest] = body

	// If ref is a tag, store the tag mapping
	if !strings.HasPrefix(ref, "sha256:") {
		m.tags[ref] = digest
	}

	w.Header().Set("Location", fmt.Sprintf("/v2/%s/manifests/%s", extractRepoName(r.URL.Path), digest))
	w.WriteHeader(http.StatusCreated)
}

// getBlob handles GET requests for blobs
func (m *MockRegistry) getBlob(w http.ResponseWriter, r *http.Request, digest string) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	blob, found := m.blobs[digest]
	if !found {
		http.Error(w, "Blob not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(blob)
}

// headBlob handles HEAD requests for blobs
func (m *MockRegistry) headBlob(w http.ResponseWriter, r *http.Request, digest string) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	blob, found := m.blobs[digest]
	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(blob)))
	w.WriteHeader(http.StatusOK)
}

// Helper functions

func extractRepoName(path string) string {
	// Extract repository name from path like /v2/repo/name/manifests/tag
	parts := strings.Split(path, "/")
	if len(parts) >= 4 {
		return strings.Join(parts[2:len(parts)-2], "/")
	}
	return "test-repo"
}

func generateDigest(data []byte) string {
	// Simple digest generation for testing
	// In real implementation, this would use SHA256
	h := make([]byte, 8)
	rand.Read(h)
	return fmt.Sprintf("sha256:%x", h)
}