package build

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockTrustManager implements TrustManager for testing
type MockTrustManager struct{}

func (m *MockTrustManager) ConfigureTLS() error {
	return nil
}

func (m *MockTrustManager) ValidateCertificate(cert []byte) error {
	return nil
}

func TestNewBuildahBuilder(t *testing.T) {
	mockTM := &MockTrustManager{}
	
	builder, err := NewBuildahBuilder(mockTM)
	
	require.NoError(t, err)
	assert.NotNil(t, builder)
	assert.Equal(t, mockTM, builder.trustManager)
}

func TestBuildImage_ValidationErrors(t *testing.T) {
	mockTM := &MockTrustManager{}
	builder, err := NewBuildahBuilder(mockTM)
	require.NoError(t, err)
	
	ctx := context.Background()
	
	// Test empty dockerfile path
	_, err = builder.BuildImage(ctx, BuildOptions{
		ContextDir: "/tmp",
		Tag:        "test:latest",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "dockerfile path cannot be empty")
	
	// Test empty context directory
	_, err = builder.BuildImage(ctx, BuildOptions{
		DockerfilePath: "/tmp/Dockerfile",
		Tag:            "test:latest",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context directory cannot be empty")
}

func TestGetRepository(t *testing.T) {
	tests := []struct {
		imageName  string
		expected   string
	}{
		{"", ""},
		{"nginx", "nginx"},
		{"nginx:latest", "nginx"},
		{"docker.io/library/nginx:1.21", "docker.io/library/nginx"},
		{"localhost:5000/myapp:v1.0.0", "localhost:5000/myapp"},
	}
	
	for _, tt := range tests {
		t.Run(tt.imageName, func(t *testing.T) {
			result := getRepository(tt.imageName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetTag(t *testing.T) {
	tests := []struct {
		imageName  string
		expected   string
	}{
		{"", ""},
		{"nginx", "latest"},
		{"nginx:v1.21", "v1.21"},
		{"docker.io/library/nginx:1.21", "1.21"},
		{"localhost:5000/myapp:v1.0.0", "v1.0.0"},
	}
	
	for _, tt := range tests {
		t.Run(tt.imageName, func(t *testing.T) {
			result := getTag(tt.imageName)
			assert.Equal(t, tt.expected, result)
		})
	}
}