package mocks

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
)

// MockResponse represents a configured HTTP response for testing
type MockResponse struct {
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    string            `json:"body,omitempty"`
}

// CallRecord tracks HTTP requests made to the mock registry
type CallRecord struct {
	Method  string            `json:"method"`
	Path    string            `json:"path"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body,omitempty"`
}

// MockRegistry provides a test double for OCI registry servers
type MockRegistry struct {
	server       *httptest.Server
	responses    map[string]MockResponse
	sequences    map[string][]MockResponse
	sequenceIdx  map[string]int
	calls        []CallRecord
	mutex        sync.RWMutex
	defaultResponse MockResponse
}

// NewMockRegistry creates a new mock registry for testing
func NewMockRegistry() *MockRegistry {
	mock := &MockRegistry{
		responses:   make(map[string]MockResponse),
		sequences:   make(map[string][]MockResponse),
		sequenceIdx: make(map[string]int),
		calls:       make([]CallRecord, 0),
		defaultResponse: MockResponse{
			Status: 404,
			Body:   `{"error": "not found"}`,
		},
	}

	mock.server = httptest.NewServer(http.HandlerFunc(mock.ServeHTTP))
	return mock
}

// NewSecureMockRegistry creates a new HTTPS mock registry for testing
func NewSecureMockRegistry() *MockRegistry {
	mock := &MockRegistry{
		responses:   make(map[string]MockResponse),
		sequences:   make(map[string][]MockResponse),
		sequenceIdx: make(map[string]int),
		calls:       make([]CallRecord, 0),
		defaultResponse: MockResponse{
			Status: 404,
			Body:   `{"error": "not found"}`,
		},
	}

	mock.server = httptest.NewTLSServer(http.HandlerFunc(mock.ServeHTTP))
	return mock
}

// URL returns the mock registry URL
func (m *MockRegistry) URL() string {
	return m.server.URL
}

// Close shuts down the mock registry server
func (m *MockRegistry) Close() {
	m.server.Close()
}

// SetResponse configures a response for a specific path
func (m *MockRegistry) SetResponse(path string, response MockResponse) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.responses[path] = response
}

// SetResponseSequence configures a sequence of responses for a path
// Useful for testing retry logic - first call gets first response, etc.
func (m *MockRegistry) SetResponseSequence(path string, responses []MockResponse) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.sequences[path] = responses
	m.sequenceIdx[path] = 0
}

// SetDefaultResponse sets the fallback response for unmatched paths
func (m *MockRegistry) SetDefaultResponse(response MockResponse) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.defaultResponse = response
}

// ServeHTTP handles incoming HTTP requests to the mock registry
func (m *MockRegistry) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.mutex.Lock()

	// Record the call
	headers := make(map[string]string)
	for key, values := range r.Header {
		headers[key] = strings.Join(values, ", ")
	}

	call := CallRecord{
		Method:  r.Method,
		Path:    r.URL.Path,
		Headers: headers,
	}

	// Read body if present
	if r.Body != nil {
		defer r.Body.Close()
		bodyBytes := make([]byte, 1024)
		n, _ := r.Body.Read(bodyBytes)
		if n > 0 {
			call.Body = string(bodyBytes[:n])
		}
	}

	m.calls = append(m.calls, call)

	// Find response to return
	var response MockResponse
	var found bool

	// Check for sequence response first
	if sequence, hasSeq := m.sequences[r.URL.Path]; hasSeq {
		idx := m.sequenceIdx[r.URL.Path]
		if idx < len(sequence) {
			response = sequence[idx]
			m.sequenceIdx[r.URL.Path] = idx + 1
			found = true
		}
	}

	// Check for single response
	if !found {
		if resp, hasResp := m.responses[r.URL.Path]; hasResp {
			response = resp
			found = true
		}
	}

	// Use default if no specific response configured
	if !found {
		response = m.defaultResponse
	}

	m.mutex.Unlock()

	// Set response headers
	for key, value := range response.Headers {
		w.Header().Set(key, value)
	}

	// Set content type if not specified
	if response.Headers["Content-Type"] == "" {
		w.Header().Set("Content-Type", "application/json")
	}

	// Write status and body
	w.WriteHeader(response.Status)
	if response.Body != "" {
		w.Write([]byte(response.Body))
	}
}

// GetCalls returns all recorded HTTP calls
func (m *MockRegistry) GetCalls() []CallRecord {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Return a copy to avoid race conditions
	calls := make([]CallRecord, len(m.calls))
	copy(calls, m.calls)
	return calls
}

// GetRequestsTo returns all requests made to a specific path
func (m *MockRegistry) GetRequestsTo(path string) []CallRecord {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	requests := make([]CallRecord, 0)
	for _, call := range m.calls {
		if call.Path == path {
			requests = append(requests, call)
		}
	}
	return requests
}

// VerifyCall checks if a specific method and path was called
func (m *MockRegistry) VerifyCall(method, path string) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, call := range m.calls {
		if call.Method == method && call.Path == path {
			return true
		}
	}
	return false
}

// GetCallCount returns the number of times a path was called
func (m *MockRegistry) GetCallCount(path string) int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	count := 0
	for _, call := range m.calls {
		if call.Path == path {
			count++
		}
	}
	return count
}

// Reset clears all recorded calls and configured responses
func (m *MockRegistry) Reset() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.calls = make([]CallRecord, 0)
	m.responses = make(map[string]MockResponse)
	m.sequences = make(map[string][]MockResponse)
	m.sequenceIdx = make(map[string]int)
}

// SetupOCIRegistryResponses configures common OCI registry endpoints
func (m *MockRegistry) SetupOCIRegistryResponses() {
	// OCI API version check
	m.SetResponse("/v2/", MockResponse{
		Status: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: `{}`,
	})

	// Authentication challenge
	m.SetResponse("/v2/_catalog", MockResponse{
		Status: 401,
		Headers: map[string]string{
			"WWW-Authenticate": `Bearer realm="https://auth.docker.io/token",service="registry.docker.io",scope="registry:catalog:*"`,
		},
		Body: `{"errors":[{"code":"UNAUTHORIZED","message":"authentication required"}]}`,
	})
}

// SimulateRateLimiting sets up rate limiting responses
func (m *MockRegistry) SimulateRateLimiting(path string, retryAfterSeconds int) {
	m.SetResponse(path, MockResponse{
		Status: 429,
		Headers: map[string]string{
			"Retry-After": fmt.Sprintf("%d", retryAfterSeconds),
		},
		Body: `{"errors":[{"code":"TOOMANYREQUESTS","message":"rate limit exceeded"}]}`,
	})
}

// SimulateServerErrors sets up server error responses for retry testing
func (m *MockRegistry) SimulateServerErrors(path string, errorCount int) {
	responses := make([]MockResponse, errorCount+1)

	// First N responses are errors
	for i := 0; i < errorCount; i++ {
		responses[i] = MockResponse{
			Status: 500,
			Body:   `{"errors":[{"code":"INTERNAL","message":"internal server error"}]}`,
		}
	}

	// Final response is success
	responses[errorCount] = MockResponse{
		Status: 200,
		Body:   `{}`,
	}

	m.SetResponseSequence(path, responses)
}