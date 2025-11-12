package registry_test

import (
	"errors"
	"testing"

	"github.com/cnoe-io/idpbuilder/pkg/registry"
)

// T1.1.2-001: RegistryClient interface compiles
func TestRegistryClientInterfaceCompiles(t *testing.T) {
	var _ registry.RegistryClient = nil
}

// T1.1.2-002: LayerStatus enum compiles with String() method
func TestLayerStatus_StringMethod(t *testing.T) {
	tests := []struct {
		status   registry.LayerStatus
		expected string
	}{
		{registry.LayerWaiting, "Waiting"},
		{registry.LayerUploading, "Uploading"},
		{registry.LayerComplete, "Complete"},
		{registry.LayerFailed, "Failed"},
		{registry.LayerStatus(99), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.status.String(); got != tt.expected {
				t.Errorf("LayerStatus.String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

// T1.1.2-003: ProgressUpdate struct compiles
func TestProgressUpdate_StructValid(t *testing.T) {
	update := registry.ProgressUpdate{
		LayerDigest:   "sha256:abc123",
		LayerSize:     1000,
		BytesUploaded: 500,
		Status:        registry.LayerUploading,
	}

	if update.LayerDigest != "sha256:abc123" {
		t.Error("ProgressUpdate struct field assignment failed")
	}
}

// T1.1.2-004: ProgressCallback type signature valid
func TestProgressCallback_TypeValid(t *testing.T) {
	var callback registry.ProgressCallback = func(update registry.ProgressUpdate) {
		// Callback implementation
	}

	// Invoke callback to verify signature
	callback(registry.ProgressUpdate{
		LayerDigest:   "sha256:test",
		LayerSize:     100,
		BytesUploaded: 50,
		Status:        registry.LayerUploading,
	})
}

// T1.1.2-005: RegistryAuthError implements error with Unwrap
func TestRegistryAuthError_ImplementsError(t *testing.T) {
	cause := errors.New("401 Unauthorized")
	err := &registry.RegistryAuthError{
		Registry: "gitea.example.com:8443",
		Cause:    cause,
	}

	var _ error = err // Compile-time check

	if !errors.Is(err, cause) {
		t.Error("RegistryAuthError should unwrap to cause")
	}

	expectedMsg := "authentication failed for registry gitea.example.com:8443: 401 Unauthorized"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message %q, got %q", expectedMsg, err.Error())
	}
}

// T1.1.2-006: RegistryConnectionError implements error with Unwrap
func TestRegistryConnectionError_ImplementsError(t *testing.T) {
	cause := errors.New("connection refused")
	err := &registry.RegistryConnectionError{
		Registry: "gitea.example.com:8443",
		Cause:    cause,
	}

	var _ error = err // Compile-time check

	if !errors.Is(err, cause) {
		t.Error("RegistryConnectionError should unwrap to cause")
	}
}

// T1.1.2-007: LayerPushError implements error with Unwrap
func TestLayerPushError_ImplementsError(t *testing.T) {
	cause := errors.New("blob upload failed")
	err := &registry.LayerPushError{
		LayerDigest: "sha256:abc123...",
		Cause:       cause,
	}

	var _ error = err // Compile-time check

	if !errors.Is(err, cause) {
		t.Error("LayerPushError should unwrap to cause")
	}
}

// T1.1.2-008: NewRegistryClient constructor signature valid
func TestNewRegistryClient_SignatureValid(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			if r != "not implemented" {
				t.Errorf("Expected panic 'not implemented', got %v", r)
			}
		} else {
			t.Error("Expected NewRegistryClient to panic (not implemented)")
		}
	}()

	_, _ = registry.NewRegistryClient(nil, nil)
}
