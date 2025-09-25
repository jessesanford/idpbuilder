package mocks

import (
	"errors"
	"net/http"
	"sync"
	"time"
)

// MockRoundTripper implements http.RoundTripper for testing HTTP transport behavior
type MockRoundTripper struct {
	responses    []MockTransportResponse
	responseIdx  int
	requests     []*http.Request
	delays       []time.Duration
	errors       []error
	mutex        sync.RWMutex
	callCount    int
}

// MockTransportResponse represents a configured HTTP response for transport testing
type MockTransportResponse struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Delay      time.Duration
	Error      error
}

// NewMockTransport creates a new mock HTTP transport for testing
func NewMockTransport() *MockRoundTripper {
	return &MockRoundTripper{
		responses: make([]MockTransportResponse, 0),
		requests:  make([]*http.Request, 0),
		delays:    make([]time.Duration, 0),
		errors:    make([]error, 0),
	}
}

// RoundTrip implements the http.RoundTripper interface
func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Record the request
	m.requests = append(m.requests, req)
	m.callCount++

	// If we have a configured response, use it
	if m.responseIdx < len(m.responses) {
		response := m.responses[m.responseIdx]
		m.responseIdx++

		// Simulate delay if configured
		if response.Delay > 0 {
			m.mutex.Unlock()
			time.Sleep(response.Delay)
			m.mutex.Lock()
		}

		// Return error if configured
		if response.Error != nil {
			return nil, response.Error
		}

		// Create HTTP response
		resp := &http.Response{
			StatusCode:    response.StatusCode,
			Header:        response.Header,
			Body:          http.NoBody,
			ContentLength: int64(len(response.Body)),
			Request:       req,
		}

		if len(response.Body) > 0 {
			resp.Body = &mockResponseBody{data: response.Body}
		}

		return resp, nil
	}

	// No configured response - return default 200 OK
	resp := &http.Response{
		StatusCode:    200,
		Header:        make(http.Header),
		Body:          http.NoBody,
		ContentLength: 0,
		Request:       req,
	}

	return resp, nil
}

// mockResponseBody implements io.ReadCloser for response body
type mockResponseBody struct {
	data []byte
	pos  int
}

func (m *mockResponseBody) Read(p []byte) (int, error) {
	if m.pos >= len(m.data) {
		return 0, nil
	}

	n := copy(p, m.data[m.pos:])
	m.pos += n
	return n, nil
}

func (m *mockResponseBody) Close() error {
	return nil
}

// SetResponse configures a single response to be returned
func (m *MockRoundTripper) SetResponse(statusCode int, body []byte, headers http.Header) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	response := MockTransportResponse{
		StatusCode: statusCode,
		Body:       body,
		Header:     headers,
	}

	m.responses = []MockTransportResponse{response}
	m.responseIdx = 0
}

// AddResponse adds a response to the sequence of responses
func (m *MockRoundTripper) AddResponse(statusCode int, body []byte, headers http.Header) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	response := MockTransportResponse{
		StatusCode: statusCode,
		Body:       body,
		Header:     headers,
	}

	m.responses = append(m.responses, response)
}

// SetError configures the transport to return an error
func (m *MockRoundTripper) SetError(err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	response := MockTransportResponse{
		Error: err,
	}

	m.responses = []MockTransportResponse{response}
	m.responseIdx = 0
}

// SetDelay configures a delay before returning responses
func (m *MockRoundTripper) SetDelay(delay time.Duration) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for i := range m.responses {
		m.responses[i].Delay = delay
	}
}

// GetRequests returns all recorded HTTP requests
func (m *MockRoundTripper) GetRequests() []*http.Request {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Return a copy to avoid race conditions
	requests := make([]*http.Request, len(m.requests))
	copy(requests, m.requests)
	return requests
}

// GetCallCount returns the number of requests made
func (m *MockRoundTripper) GetCallCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.callCount
}

// Reset clears all recorded requests and responses
func (m *MockRoundTripper) Reset() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.responses = make([]MockTransportResponse, 0)
	m.requests = make([]*http.Request, 0)
	m.delays = make([]time.Duration, 0)
	m.errors = make([]error, 0)
	m.responseIdx = 0
	m.callCount = 0
}

// SimulateNetworkError creates a transport that simulates network failures
func SimulateNetworkError(errorType string) *MockRoundTripper {
	transport := NewMockTransport()

	var err error
	switch errorType {
	case "timeout":
		err = &mockTimeoutError{}
	case "dns":
		err = &mockDNSError{}
	case "connection_refused":
		err = &mockConnectionError{reason: "connection refused"}
	case "tls":
		err = &mockTLSError{}
	default:
		err = errors.New("network error")
	}

	transport.SetError(err)
	return transport
}

// Mock error types for realistic error simulation

type mockTimeoutError struct{}

func (e *mockTimeoutError) Error() string   { return "context deadline exceeded" }
func (e *mockTimeoutError) Timeout() bool   { return true }
func (e *mockTimeoutError) Temporary() bool { return true }

type mockDNSError struct{}

func (e *mockDNSError) Error() string   { return "no such host" }
func (e *mockDNSError) Timeout() bool   { return false }
func (e *mockDNSError) Temporary() bool { return false }

type mockConnectionError struct {
	reason string
}

func (e *mockConnectionError) Error() string   { return e.reason }
func (e *mockConnectionError) Timeout() bool   { return false }
func (e *mockConnectionError) Temporary() bool { return false }

type mockTLSError struct{}

func (e *mockTLSError) Error() string   { return "tls: bad certificate" }
func (e *mockTLSError) Timeout() bool   { return false }
func (e *mockTLSError) Temporary() bool { return false }

// SimulateRetryScenario creates responses for testing retry logic
func (m *MockRoundTripper) SimulateRetryScenario(failCount int, finalStatus int) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.responses = make([]MockTransportResponse, 0, failCount+1)

	// Add failure responses
	for i := 0; i < failCount; i++ {
		m.responses = append(m.responses, MockTransportResponse{
			StatusCode: 500,
			Body:       []byte(`{"error": "internal server error"}`),
			Header:     make(http.Header),
		})
	}

	// Add final success/failure response
	m.responses = append(m.responses, MockTransportResponse{
		StatusCode: finalStatus,
		Body:       []byte(`{"success": true}`),
		Header:     make(http.Header),
	})

	m.responseIdx = 0
}

// SimulateRateLimiting creates responses for testing rate limiting
func (m *MockRoundTripper) SimulateRateLimiting(retryAfter string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	headers := make(http.Header)
	headers.Set("Retry-After", retryAfter)

	m.responses = []MockTransportResponse{{
		StatusCode: 429,
		Body:       []byte(`{"error": "rate limit exceeded"}`),
		Header:     headers,
	}}

	m.responseIdx = 0
}

// MockInsecureTransport creates a transport that simulates insecure connections
type MockInsecureTransport struct {
	*MockRoundTripper
	insecure bool
}

// NewMockInsecureTransport creates a transport that tracks insecure flag usage
func NewMockInsecureTransport(insecure bool) *MockInsecureTransport {
	return &MockInsecureTransport{
		MockRoundTripper: NewMockTransport(),
		insecure:         insecure,
	}
}

// IsInsecure returns whether this transport allows insecure connections
func (m *MockInsecureTransport) IsInsecure() bool {
	return m.insecure
}

// SetInsecure updates the insecure flag
func (m *MockInsecureTransport) SetInsecure(insecure bool) {
	m.insecure = insecure
}