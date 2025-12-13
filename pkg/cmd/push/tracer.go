// Package push implements OCI registry push functionality for idpbuilder.
package push

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

// NewDebugLogger creates logger configured for debug mode.
// Per REQ-005: When --log-level debug, enable verbose output.
// Per REQ-006: When --log-level info, log only high-level operational steps.
func NewDebugLogger(level slog.Level, phase string) *slog.Logger {
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	})
	return slog.New(handler).With(slog.String("phase", phase))
}

// LogCredentialResolution logs credential flow at debug level.
// Per REQ-020: NEVER log actual credential values.
// Only logs presence/absence flags for diagnosis.
func LogCredentialResolution(logger *slog.Logger, source string, hasUsername, hasPassword, hasToken bool) {
	logger.Debug("credential resolution",
		slog.String("source", source), // "flags", "env", "anonymous"
		slog.Bool("has_username", hasUsername),
		slog.Bool("has_password", hasPassword),
		slog.Bool("has_token", hasToken),
	)
}

// generateRequestID creates unique ID for request/response correlation.
// Format: "req-{UnixNano}" for uniqueness and sortability.
func generateRequestID() string {
	return fmt.Sprintf("req-%d", time.Now().UnixNano())
}

// DebugTransport wraps http.RoundTripper for request/response logging.
// Per REQ-025: Log full HTTP interactions in debug mode with correlation IDs.
// Per REQ-020: Redact Authorization header values in all logs.
type DebugTransport struct {
	// Base is the underlying HTTP RoundTripper
	Base http.RoundTripper
	// Logger is the structured logger for debug output
	Logger *slog.Logger
}

// RoundTrip implements http.RoundTripper with logging.
// Generates unique request ID and logs both request and response.
func (t *DebugTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	requestID := generateRequestID()

	// Log outgoing request (REQ-025)
	t.logRequest(req, requestID)

	start := time.Now()
	resp, err := t.Base.RoundTrip(req)
	duration := time.Since(start)

	if err != nil {
		t.Logger.Debug("HTTP request failed",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	// Log response (REQ-025 correlation)
	t.logResponse(resp, requestID, duration)

	return resp, nil
}

// logRequest dumps HTTP request with Authorization redaction.
// Per REQ-020: Redact Authorization header values.
// Clones request to avoid modifying the original.
func (t *DebugTransport) logRequest(req *http.Request, requestID string) {
	reqCopy := req.Clone(req.Context())

	// CRITICAL: Redact credentials (REQ-020)
	if reqCopy.Header.Get("Authorization") != "" {
		reqCopy.Header.Set("Authorization", "[REDACTED]")
	}

	dump, _ := httputil.DumpRequestOut(reqCopy, false)

	t.Logger.Debug("HTTP request",
		slog.String("request_id", requestID),
		slog.String("method", req.Method),
		slog.String("url", req.URL.String()),
		slog.String("dump", string(dump)),
	)
}

// logResponse dumps HTTP response with correlation.
// Per REQ-025: Same request_id as request for correlation.
func (t *DebugTransport) logResponse(resp *http.Response, requestID string, duration time.Duration) {
	dump, _ := httputil.DumpResponse(resp, false)

	t.Logger.Debug("HTTP response",
		slog.String("request_id", requestID),
		slog.Int("status_code", resp.StatusCode),
		slog.Duration("duration", duration),
		slog.String("dump", string(dump)),
	)
}
