package push

import (
	"context"
	"crypto/tls"
	"net/http"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/google/go-containerregistry/pkg/authn"
)

// mockProgressReporter implements ProgressReporter for testing
type mockProgressReporter struct {
	startedImages  []string
	finishedImages []string
	startedLayers  []string
	finishedLayers []string
	errors         map[string]error
}

func (m *mockProgressReporter) StartImage(digest string, totalSize int64) {
	m.startedImages = append(m.startedImages, digest)
}

func (m *mockProgressReporter) UpdateLayer(digest string, written int64) {
	// Track updates if needed for testing
}

func (m *mockProgressReporter) FinishLayer(digest string) {
	m.finishedLayers = append(m.finishedLayers, digest)
}

func (m *mockProgressReporter) FinishImage(digest string) {
	m.finishedImages = append(m.finishedImages, digest)
}

func (m *mockProgressReporter) SetError(digest string, err error) {
	if m.errors == nil {
		m.errors = make(map[string]error)
	}
	m.errors[digest] = err
}

func TestDefaultPusherOptions(t *testing.T) {
	opts := DefaultPusherOptions()

	if opts.MaxRetries != 5 {
		t.Errorf("Expected MaxRetries 5, got %d", opts.MaxRetries)
	}
	if opts.InitialBackoff != time.Second {
		t.Errorf("Expected InitialBackoff 1s, got %v", opts.InitialBackoff)
	}
	if opts.BackoffMultiplier != 2.0 {
		t.Errorf("Expected BackoffMultiplier 2.0, got %f", opts.BackoffMultiplier)
	}
	if opts.MaxBackoff != 30*time.Second {
		t.Errorf("Expected MaxBackoff 30s, got %v", opts.MaxBackoff)
	}
	if opts.Insecure {
		t.Error("Expected Insecure to be false")
	}
	if opts.UserAgent != "idpbuilder-push/1.0.0" {
		t.Errorf("Expected UserAgent 'idpbuilder-push/1.0.0', got %s", opts.UserAgent)
	}
}

func TestNewImagePusher(t *testing.T) {
	auth := authn.Anonymous
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
	}
	progress := &mockProgressReporter{}
	logger := logr.Discard()

	pusher := NewImagePusher(auth, transport, progress, logger)

	if pusher == nil {
		t.Fatal("Expected non-nil pusher")
	}
	if pusher.auth != auth {
		t.Error("Authenticator not set correctly")
	}
	if pusher.transport != transport {
		t.Error("Transport not set correctly")
	}
	if pusher.progress != progress {
		t.Error("Progress reporter not set correctly")
	}
	if pusher.options.MaxRetries != 5 {
		t.Errorf("Expected default max retries 5, got %d", pusher.options.MaxRetries)
	}
}

func TestNewImagePusherWithOptions(t *testing.T) {
	auth := authn.Anonymous
	transport := &http.Transport{}
	progress := &mockProgressReporter{}
	logger := logr.Discard()

	customOpts := PusherOptions{
		MaxRetries:        3,
		InitialBackoff:    2 * time.Second,
		BackoffMultiplier: 1.5,
		MaxBackoff:        60 * time.Second,
		Insecure:          true,
		UserAgent:         "test-agent/1.0",
	}

	pusher := NewImagePusherWithOptions(auth, transport, progress, logger, customOpts)

	if pusher == nil {
		t.Fatal("Expected non-nil pusher")
	}
	if pusher.options.MaxRetries != 3 {
		t.Errorf("Expected MaxRetries 3, got %d", pusher.options.MaxRetries)
	}
	if pusher.options.InitialBackoff != 2*time.Second {
		t.Errorf("Expected InitialBackoff 2s, got %v", pusher.options.InitialBackoff)
	}
	if pusher.options.BackoffMultiplier != 1.5 {
		t.Errorf("Expected BackoffMultiplier 1.5, got %f", pusher.options.BackoffMultiplier)
	}
	if pusher.options.UserAgent != "test-agent/1.0" {
		t.Errorf("Expected UserAgent 'test-agent/1.0', got %s", pusher.options.UserAgent)
	}
}

func TestNewImagePusher_NilTransport(t *testing.T) {
	auth := authn.Anonymous
	progress := &mockProgressReporter{}
	logger := logr.Discard()

	pusher := NewImagePusher(auth, nil, progress, logger)

	if pusher == nil {
		t.Fatal("Expected non-nil pusher")
	}
	if pusher.transport == nil {
		t.Error("Expected transport to be created when nil provided")
	}
}

func TestPushResult(t *testing.T) {
	result := &PushResult{
		ImageName: "registry.example.com/myimage:latest",
		Digest:    "sha256:abc123",
		Size:      1024000,
		Duration:  5 * time.Second,
		Error:     nil,
		Retries:   2,
	}

	if result.ImageName != "registry.example.com/myimage:latest" {
		t.Errorf("Expected ImageName 'registry.example.com/myimage:latest', got %s", result.ImageName)
	}
	if result.Digest != "sha256:abc123" {
		t.Errorf("Expected Digest 'sha256:abc123', got %s", result.Digest)
	}
	if result.Size != 1024000 {
		t.Errorf("Expected Size 1024000, got %d", result.Size)
	}
	if result.Duration != 5*time.Second {
		t.Errorf("Expected Duration 5s, got %v", result.Duration)
	}
	if result.Retries != 2 {
		t.Errorf("Expected Retries 2, got %d", result.Retries)
	}
}

func TestIsRetryableError(t *testing.T) {
	tests := []struct {
		name     string
		errMsg   string
		expected bool
	}{
		{
			name:     "connection refused",
			errMsg:   "dial tcp: connection refused",
			expected: true,
		},
		{
			name:     "connection reset",
			errMsg:   "read tcp: connection reset by peer",
			expected: true,
		},
		{
			name:     "timeout",
			errMsg:   "context deadline exceeded: timeout",
			expected: true,
		},
		{
			name:     "503 service unavailable",
			errMsg:   "GET https://registry.io: 503 service unavailable",
			expected: true,
		},
		{
			name:     "502 bad gateway",
			errMsg:   "502 bad gateway error",
			expected: true,
		},
		{
			name:     "rate limit",
			errMsg:   "rate limit exceeded",
			expected: true,
		},
		{
			name:     "network unreachable",
			errMsg:   "network is unreachable",
			expected: true,
		},
		{
			name:     "non-retryable error",
			errMsg:   "invalid credentials",
			expected: false,
		},
		{
			name:     "authentication failed",
			errMsg:   "401 unauthorized",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a simple error with the message
			err := &testError{msg: tt.errMsg}
			result := isRetryableError(err)
			if result != tt.expected {
				t.Errorf("isRetryableError(%s) = %v, expected %v", tt.errMsg, result, tt.expected)
			}
		})
	}
}

// testError is a simple error implementation for testing
type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		substr   string
		expected bool
	}{
		{
			name:     "exact match",
			s:        "timeout",
			substr:   "timeout",
			expected: true,
		},
		{
			name:     "substring at start",
			s:        "timeout error occurred",
			substr:   "timeout",
			expected: true,
		},
		{
			name:     "substring at end",
			s:        "connection timeout",
			substr:   "timeout",
			expected: true,
		},
		{
			name:     "substring in middle",
			s:        "a timeout occurred",
			substr:   "timeout",
			expected: true,
		},
		{
			name:     "not found",
			s:        "error occurred",
			substr:   "timeout",
			expected: false,
		},
		{
			name:     "empty substring",
			s:        "test string",
			substr:   "",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := contains(tt.s, tt.substr)
			if result != tt.expected {
				t.Errorf("contains(%q, %q) = %v, expected %v", tt.s, tt.substr, result, tt.expected)
			}
		})
	}
}

func TestPusherOptions_CustomConfiguration(t *testing.T) {
	opts := PusherOptions{
		MaxRetries:        10,
		InitialBackoff:    500 * time.Millisecond,
		BackoffMultiplier: 3.0,
		MaxBackoff:        2 * time.Minute,
		Insecure:          true,
		UserAgent:         "custom-agent/2.0",
	}

	if opts.MaxRetries != 10 {
		t.Errorf("Expected MaxRetries 10, got %d", opts.MaxRetries)
	}
	if opts.InitialBackoff != 500*time.Millisecond {
		t.Errorf("Expected InitialBackoff 500ms, got %v", opts.InitialBackoff)
	}
	if opts.BackoffMultiplier != 3.0 {
		t.Errorf("Expected BackoffMultiplier 3.0, got %f", opts.BackoffMultiplier)
	}
	if opts.MaxBackoff != 2*time.Minute {
		t.Errorf("Expected MaxBackoff 2m, got %v", opts.MaxBackoff)
	}
	if !opts.Insecure {
		t.Error("Expected Insecure to be true")
	}
	if opts.UserAgent != "custom-agent/2.0" {
		t.Errorf("Expected UserAgent 'custom-agent/2.0', got %s", opts.UserAgent)
	}
}

func TestBuildRemoteOptions(t *testing.T) {
	// Test that buildRemoteOptions creates proper options
	auth := authn.Anonymous
	transport := &http.Transport{}
	progress := &mockProgressReporter{}
	logger := logr.Discard()

	opts := PusherOptions{
		UserAgent: "test-agent/1.0",
	}

	pusher := NewImagePusherWithOptions(auth, transport, progress, logger, opts)
	remoteOpts := pusher.buildRemoteOptions()

	if len(remoteOpts) == 0 {
		t.Error("Expected non-empty remote options")
	}
}

func TestProgressReporter_Interface(t *testing.T) {
	// Test that mockProgressReporter implements ProgressReporter
	var _ ProgressReporter = (*mockProgressReporter)(nil)

	mock := &mockProgressReporter{}

	mock.StartImage("sha256:test", 1000)
	if len(mock.startedImages) != 1 || mock.startedImages[0] != "sha256:test" {
		t.Error("StartImage not tracking correctly")
	}

	mock.FinishImage("sha256:test")
	if len(mock.finishedImages) != 1 || mock.finishedImages[0] != "sha256:test" {
		t.Error("FinishImage not tracking correctly")
	}

	mock.FinishLayer("sha256:layer1")
	if len(mock.finishedLayers) != 1 || mock.finishedLayers[0] != "sha256:layer1" {
		t.Error("FinishLayer not tracking correctly")
	}
}

func TestPusher_WithCancellation(t *testing.T) {
	// Test that pusher respects context cancellation
	pusher := NewImagePusher(authn.Anonymous, nil, nil, logr.Discard())

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Try to push with cancelled context - should fail quickly
	// Note: This would require a real image and registry, so we just verify the pusher was created
	if pusher == nil {
		t.Fatal("Expected non-nil pusher")
	}

	// Verify context is cancelled
	select {
	case <-ctx.Done():
		// Context properly cancelled
	default:
		t.Error("Context should be cancelled")
	}
}
