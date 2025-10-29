package docker

import (
	"context"
	"testing"

	v1 "github.com/google/go-containerregistry/pkg/v1"
)

// MockClient is a mock implementation of the Client interface for testing.
// This demonstrates that the Client interface is implementable.
type MockClient struct {
	ImageExistsFunc   func(ctx context.Context, imageName string) (bool, error)
	GetImageFunc      func(ctx context.Context, imageName string) (v1.Image, error)
	ValidateImageFunc func(imageName string) error
	CloseFunc         func() error
}

func (m *MockClient) ImageExists(ctx context.Context, imageName string) (bool, error) {
	if m.ImageExistsFunc != nil {
		return m.ImageExistsFunc(ctx, imageName)
	}
	return true, nil
}

func (m *MockClient) GetImage(ctx context.Context, imageName string) (v1.Image, error) {
	if m.GetImageFunc != nil {
		return m.GetImageFunc(ctx, imageName)
	}
	return nil, nil
}

func (m *MockClient) ValidateImageName(imageName string) error {
	if m.ValidateImageFunc != nil {
		return m.ValidateImageFunc(imageName)
	}
	return nil
}

func (m *MockClient) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

// TestMockClientImplementsInterface verifies MockClient satisfies Client interface.
// This is a compile-time verification test.
func TestMockClientImplementsInterface(t *testing.T) {
	var _ Client = &MockClient{}
	t.Log("✅ MockClient implements Client interface")
}

// TestMockClientMethodCalls verifies all methods are callable.
func TestMockClientMethodCalls(t *testing.T) {
	mock := &MockClient{
		ImageExistsFunc: func(ctx context.Context, imageName string) (bool, error) {
			return true, nil
		},
		GetImageFunc: func(ctx context.Context, imageName string) (v1.Image, error) {
			return nil, nil
		},
		ValidateImageFunc: func(imageName string) error {
			return nil
		},
		CloseFunc: func() error {
			return nil
		},
	}

	ctx := context.Background()

	// Test ImageExists
	exists, err := mock.ImageExists(ctx, "test:latest")
	if !exists || err != nil {
		t.Errorf("ImageExists failed: exists=%v, err=%v", exists, err)
	}

	// Test GetImage
	_, err = mock.GetImage(ctx, "test:latest")
	if err != nil {
		t.Errorf("GetImage failed: %v", err)
	}

	// Test ValidateImageName
	err = mock.ValidateImageName("test:latest")
	if err != nil {
		t.Errorf("ValidateImageName failed: %v", err)
	}

	// Test Close
	err = mock.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}

	t.Log("✅ All mock methods callable")
}
