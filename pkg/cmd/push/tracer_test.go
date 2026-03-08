package push

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewDebugLogger tests logger creation with different levels
func TestNewDebugLogger(t *testing.T) {
	tests := []struct {
		name  string
		level slog.Level
	}{
		{"debug", slog.LevelDebug},
		{"info", slog.LevelInfo},
		{"warn", slog.LevelWarn},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewDebugLogger(tt.level, "test-phase")
			assert.NotNil(t, logger)
		})
	}
}

// TestLogCredentialResolution verifies credentials are never logged
func TestLogCredentialResolution(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	LogCredentialResolution(logger, "flags", true, true, false)

	output := buf.String()
	assert.Contains(t, output, "source=flags")
	assert.Contains(t, output, "has_username=true")
	assert.Contains(t, output, "has_password=true")
	assert.Contains(t, output, "has_token=false")
}

// TestLogCredentialResolution_NoCredentials tests logging with no credentials
func TestLogCredentialResolution_NoCredentials(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	LogCredentialResolution(logger, "anonymous", false, false, false)

	output := buf.String()
	assert.Contains(t, output, "source=anonymous")
	assert.Contains(t, output, "has_username=false")
	assert.Contains(t, output, "has_password=false")
	assert.Contains(t, output, "has_token=false")
}

// mockRoundTripper for testing
type mockRoundTripper struct {
	response *http.Response
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.response, nil
}

// TestDebugTransport_RequestLogging verifies requests are logged with Authorization redaction
func TestDebugTransport_RequestLogging(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	mockRT := &mockRoundTripper{
		response: &http.Response{
			StatusCode: 200,
			Status:     "200 OK",
			Body:       io.NopCloser(bytes.NewBufferString("")),
		},
	}

	transport := &DebugTransport{
		Base:   mockRT,
		Logger: logger,
	}

	req, _ := http.NewRequest("GET", "https://registry.example.com/v2/", nil)
	req.Header.Set("Authorization", "Bearer secret-token")

	_, _ = transport.RoundTrip(req)

	output := buf.String()

	// Verify request was logged
	assert.Contains(t, output, "HTTP request")
	assert.Contains(t, output, "request_id=")

	// CRITICAL: Verify Authorization header was redacted
	assert.Contains(t, output, "[REDACTED]")
	assert.NotContains(t, output, "secret-token")
}

// TestDebugTransport_ResponseLogging verifies responses are logged
func TestDebugTransport_ResponseLogging(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	mockRT := &mockRoundTripper{
		response: &http.Response{
			StatusCode: 200,
			Status:     "200 OK",
			Body:       io.NopCloser(bytes.NewBufferString("")),
		},
	}

	transport := &DebugTransport{
		Base:   mockRT,
		Logger: logger,
	}

	req, _ := http.NewRequest("GET", "https://registry.example.com/v2/", nil)
	_, _ = transport.RoundTrip(req)

	output := buf.String()

	// Verify response was logged
	assert.Contains(t, output, "HTTP response")
	assert.Contains(t, output, "status_code=200")
	assert.Contains(t, output, "duration=")
}

// TestDebugTransport_RequestResponseCorrelation verifies same request_id in request and response
func TestDebugTransport_RequestResponseCorrelation(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	mockRT := &mockRoundTripper{
		response: &http.Response{
			StatusCode: 201,
			Status:     "201 Created",
			Body:       io.NopCloser(bytes.NewBufferString("")),
		},
	}

	transport := &DebugTransport{
		Base:   mockRT,
		Logger: logger,
	}

	req, _ := http.NewRequest("POST", "https://registry.example.com/v2/manifests", nil)
	_, _ = transport.RoundTrip(req)

	output := buf.String()
	lines := strings.Split(output, "\n")

	// Extract request_id from request log
	var requestID string
	for _, line := range lines {
		if strings.Contains(line, "HTTP request") {
			parts := strings.Split(line, "request_id=")
			if len(parts) > 1 {
				requestID = strings.Fields(parts[1])[0]
				break
			}
		}
	}

	// Verify request_id was found
	assert.NotEmpty(t, requestID, "request_id should be present in request log")

	// Verify same request_id appears in response log
	found := false
	for _, line := range lines {
		if strings.Contains(line, "HTTP response") && strings.Contains(line, requestID) {
			found = true
			break
		}
	}
	assert.True(t, found, "response should have same request_id as request")
}

// TestGenerateRequestID verifies uniqueness
func TestGenerateRequestID(t *testing.T) {
	id1 := generateRequestID()
	id2 := generateRequestID()

	assert.NotEqual(t, id1, id2, "request IDs should be unique")
	assert.True(t, strings.HasPrefix(id1, "req-"), "request ID should have req- prefix")
	assert.True(t, strings.HasPrefix(id2, "req-"), "request ID should have req- prefix")
}
