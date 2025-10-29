package registry

import (
	"context"
	"testing"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/stretchr/testify/assert"
)

// MockRegistryClient is a mock implementation of the Client interface for testing.
type MockRegistryClient struct {
	PushFunc                 func(ctx context.Context, image v1.Image, targetRef string, progressCallback ProgressCallback) error
	BuildImageReferenceFunc  func(registryURL, imageName string) (string, error)
	ValidateRegistryFunc     func(ctx context.Context, registryURL string) error
	ProgressCallbackInvoked  bool
	ProgressUpdatesReceived  []ProgressUpdate
}

// Push implements the Client interface.
func (m *MockRegistryClient) Push(ctx context.Context, image v1.Image, targetRef string, progressCallback ProgressCallback) error {
	if m.PushFunc != nil {
		return m.PushFunc(ctx, image, targetRef, progressCallback)
	}

	// Simulate progress updates if callback provided
	if progressCallback != nil {
		m.ProgressCallbackInvoked = true
		update := ProgressUpdate{
			LayerDigest: "sha256:mock123",
			LayerSize:   1024,
			BytesPushed: 1024,
			Status:      "complete",
		}
		m.ProgressUpdatesReceived = append(m.ProgressUpdatesReceived, update)
		progressCallback(update)
	}

	return nil
}

// BuildImageReference implements the Client interface.
func (m *MockRegistryClient) BuildImageReference(registryURL, imageName string) (string, error) {
	if m.BuildImageReferenceFunc != nil {
		return m.BuildImageReferenceFunc(registryURL, imageName)
	}
	return registryURL + "/" + imageName, nil
}

// ValidateRegistry implements the Client interface.
func (m *MockRegistryClient) ValidateRegistry(ctx context.Context, registryURL string) error {
	if m.ValidateRegistryFunc != nil {
		return m.ValidateRegistryFunc(ctx, registryURL)
	}
	return nil
}

// TestMockClientImplementsInterface verifies MockRegistryClient satisfies Client interface.
func TestMockClientImplementsInterface(t *testing.T) {
	var _ Client = &MockRegistryClient{}
	assert.True(t, true, "MockRegistryClient implements Client interface")
}

// TestMockClientPush verifies Push method is callable.
func TestMockClientPush(t *testing.T) {
	mock := &MockRegistryClient{}
	ctx := context.Background()

	err := mock.Push(ctx, nil, "registry.io/myapp:latest", nil)
	assert.NoError(t, err)
}

// TestMockClientPushWithCallback verifies callback invocation.
func TestMockClientPushWithCallback(t *testing.T) {
	mock := &MockRegistryClient{}
	ctx := context.Background()

	var receivedUpdate ProgressUpdate
	callback := func(update ProgressUpdate) {
		receivedUpdate = update
	}

	err := mock.Push(ctx, nil, "registry.io/myapp:latest", callback)
	assert.NoError(t, err)
	assert.True(t, mock.ProgressCallbackInvoked, "Callback should be invoked")
	assert.Equal(t, "sha256:mock123", receivedUpdate.LayerDigest)
	assert.Equal(t, int64(1024), receivedUpdate.LayerSize)
	assert.Equal(t, int64(1024), receivedUpdate.BytesPushed)
	assert.Equal(t, "complete", receivedUpdate.Status)
}

// TestMockClientBuildImageReference verifies BuildImageReference method.
func TestMockClientBuildImageReference(t *testing.T) {
	mock := &MockRegistryClient{}

	ref, err := mock.BuildImageReference("https://registry.io", "myapp:v1.0")
	assert.NoError(t, err)
	assert.Equal(t, "https://registry.io/myapp:v1.0", ref)
}

// TestMockClientValidateRegistry verifies ValidateRegistry method.
func TestMockClientValidateRegistry(t *testing.T) {
	mock := &MockRegistryClient{}
	ctx := context.Background()

	err := mock.ValidateRegistry(ctx, "https://registry.io")
	assert.NoError(t, err)
}

// TestMockClientCustomBehavior verifies custom function injection.
func TestMockClientCustomBehavior(t *testing.T) {
	mock := &MockRegistryClient{
		PushFunc: func(ctx context.Context, image v1.Image, targetRef string, progressCallback ProgressCallback) error {
			return &PushFailedError{TargetRef: targetRef, Cause: assert.AnError}
		},
	}

	ctx := context.Background()
	err := mock.Push(ctx, nil, "registry.io/myapp:latest", nil)
	assert.Error(t, err)
	assert.IsType(t, &PushFailedError{}, err)
}
