package registry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestClientInterfaceCompilation verifies the Client interface compiles.
func TestClientInterfaceCompilation(t *testing.T) {
	// This test ensures the Client interface is properly defined.
	// If this compiles, the interface is valid.
	var _ Client
	assert.True(t, true, "Client interface compiles")
}

// TestProgressCallbackType verifies ProgressCallback type signature.
func TestProgressCallbackType(t *testing.T) {
	// Verify ProgressCallback is a function type
	var callback ProgressCallback
	assert.Nil(t, callback, "ProgressCallback can be nil")

	// Verify it accepts ProgressUpdate
	callback = func(update ProgressUpdate) {
		// Do nothing - just verify signature
	}
	assert.NotNil(t, callback, "ProgressCallback can be assigned")
}

// TestProgressUpdateStruct verifies ProgressUpdate struct fields.
func TestProgressUpdateStruct(t *testing.T) {
	update := ProgressUpdate{
		LayerDigest: "sha256:abc123",
		LayerSize:   1024,
		BytesPushed: 512,
		Status:      "uploading",
	}

	assert.Equal(t, "sha256:abc123", update.LayerDigest)
	assert.Equal(t, int64(1024), update.LayerSize)
	assert.Equal(t, int64(512), update.BytesPushed)
	assert.Equal(t, "uploading", update.Status)
}

// TestNewClientSignature verifies NewClient function signature.
func TestNewClientSignature(t *testing.T) {
	// This test verifies NewClient exists with correct signature.
	// We expect it to panic since it's not implemented in Wave 1.
	assert.Panics(t, func() {
		NewClient(nil, nil)
	}, "NewClient should panic in Wave 1")
}

// TestAuthProviderInterface verifies AuthProvider forward reference.
func TestAuthProviderInterface(t *testing.T) {
	// Verify AuthProvider interface compiles
	var _ AuthProvider
	assert.True(t, true, "AuthProvider interface compiles")
}

// TestTLSConfigProviderInterface verifies TLSConfigProvider forward reference.
func TestTLSConfigProviderInterface(t *testing.T) {
	// Verify TLSConfigProvider interface compiles
	var _ TLSConfigProvider
	assert.True(t, true, "TLSConfigProvider interface compiles")
}
