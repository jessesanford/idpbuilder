package registry

import (
	"fmt"
	"net/http"
	"time"
)

type retryTransport struct {
	base       http.RoundTripper
	maxRetries int
	backoff    time.Duration
}

func newRetryTransport(base http.RoundTripper, maxRetries int, backoff time.Duration) *retryTransport {
	if base == nil {
		base = http.DefaultTransport
	}
	return &retryTransport{
		base:       base,
		maxRetries: maxRetries,
		backoff:    backoff,
	}
}

// RoundTrip implements the http.RoundTripper interface with retry logic
// Retries on network errors and 5xx server errors, but not on 4xx client errors
func (rt *retryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var lastErr error
	var lastResp *http.Response

	for attempt := 0; attempt <= rt.maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff with jitter
			delay := rt.backoff * time.Duration(1<<uint(attempt-1))
			time.Sleep(delay)
		}

		// Clone the request for retry attempts
		reqClone := rt.cloneRequest(req)
		
		resp, err := rt.base.RoundTrip(reqClone)
		
		// If successful or client error (4xx), don't retry
		if err == nil {
			if resp.StatusCode < 500 {
				return resp, nil
			}
			
			// Server error - close response body and retry
			resp.Body.Close()
			lastResp = resp
			lastErr = fmt.Errorf("server error: %d %s", resp.StatusCode, resp.Status)
			continue
		}

		// Network error - retry
		lastErr = err
	}

	// All retries exhausted
	if lastResp != nil {
		return lastResp, fmt.Errorf("max retries exceeded with server error: %w", lastErr)
	}
	
	return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}

// cloneRequest creates a shallow clone of the HTTP request for retry attempts
func (rt *retryTransport) cloneRequest(req *http.Request) *http.Request {
	clone := req.Clone(req.Context())
	
	// Reset body if it's seekable
	if req.Body != nil && req.GetBody != nil {
		body, err := req.GetBody()
		if err == nil {
			clone.Body = body
		}
	}
	
	return clone
}