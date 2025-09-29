package providers

import (
	"errors"
	"testing"
	"time"
)

// TestArtifactCreation tests the creation and basic properties of Artifact types.
func TestArtifactCreation(t *testing.T) {
	artifact := Artifact{
		MediaType: "application/vnd.oci.image.manifest.v1+json",
		Manifest:  []byte(`{"test": "manifest"}`),
		Layers: []Layer{
			{
				MediaType: "application/vnd.oci.image.layer.v1.tar+gzip",
				Digest:    "sha256:abc123",
				Size:      1024,
				Data:      []byte("layer data"),
			},
		},
		Config:      []byte(`{"test": "config"}`),
		Annotations: map[string]string{"org.opencontainers.image.title": "test"},
	}

	if artifact.MediaType != "application/vnd.oci.image.manifest.v1+json" {
		t.Errorf("Expected MediaType to be set correctly, got %s", artifact.MediaType)
	}

	if len(artifact.Layers) != 1 {
		t.Errorf("Expected 1 layer, got %d", len(artifact.Layers))
	}

	if artifact.Layers[0].Size != 1024 {
		t.Errorf("Expected layer size 1024, got %d", artifact.Layers[0].Size)
	}

	if artifact.Annotations["org.opencontainers.image.title"] != "test" {
		t.Errorf("Expected annotation value 'test', got %s", artifact.Annotations["org.opencontainers.image.title"])
	}
}

// TestProviderErrorFormat tests the error formatting of ProviderError.
func TestProviderErrorFormat(t *testing.T) {
	baseErr := errors.New("connection failed")

	// Test with reference
	errWithRef := ProviderError{
		Op:  "pull",
		Ref: "registry.example.com/repo:tag",
		Err: baseErr,
	}

	expected := "provider operation pull failed for registry.example.com/repo:tag: connection failed"
	if errWithRef.Error() != expected {
		t.Errorf("Expected error message %q, got %q", expected, errWithRef.Error())
	}

	// Test without reference
	errWithoutRef := ProviderError{
		Op:  "list",
		Err: baseErr,
	}

	expected = "provider operation list failed: connection failed"
	if errWithoutRef.Error() != expected {
		t.Errorf("Expected error message %q, got %q", expected, errWithoutRef.Error())
	}

	// Test error unwrapping
	if !errors.Is(errWithRef, baseErr) {
		t.Error("ProviderError should unwrap to the base error")
	}
}

// TestAuthConfigMasking tests that AuthConfig can be created with sensitive data.
// Note: This doesn't test actual masking since that would be in a later effort.
func TestAuthConfigMasking(t *testing.T) {
	auth := AuthConfig{
		Username:      "testuser",
		Password:      "secret123",
		Token:         "bearer-token-123",
		RegistryToken: "registry-specific-token",
	}

	if auth.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got %s", auth.Username)
	}

	if auth.Password != "secret123" {
		t.Errorf("Expected password to be set (not testing masking in this effort)")
	}

	if auth.Token == "" {
		t.Error("Expected token to be set")
	}

	if auth.RegistryToken == "" {
		t.Error("Expected registry token to be set")
	}

	// Test with ProviderConfig
	config := ProviderConfig{
		Registry: "registry.example.com",
		Auth:     auth,
		Insecure: false,
		Timeout:  30 * time.Second,
	}

	if config.Auth.Username != "testuser" {
		t.Errorf("Expected auth username in config, got %s", config.Auth.Username)
	}

	if config.Timeout != 30*time.Second {
		t.Errorf("Expected timeout 30s, got %v", config.Timeout)
	}
}