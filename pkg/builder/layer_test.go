package builder

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLayerFactory tests layer creation functionality
func TestLayerFactory(t *testing.T) {
	factory := NewLayerFactory()
	assert.NotNil(t, factory)
	
	// Create a temporary test directory
	tempDir, err := os.MkdirTemp("", "layer-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	// Create test content
	testFile := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	require.NoError(t, err)
	
	// Create subdirectory
	subDir := filepath.Join(tempDir, "subdir")
	err = os.Mkdir(subDir, 0755)
	require.NoError(t, err)
	
	subFile := filepath.Join(subDir, "sub.txt")
	err = os.WriteFile(subFile, []byte("sub content"), 0644)
	require.NoError(t, err)
	
	t.Run("create layer from directory", func(t *testing.T) {
		layer, err := factory.CreateLayer(tempDir)
		assert.NoError(t, err)
		assert.NotNil(t, layer)
		
		// Verify layer properties
		size, err := layer.Size()
		assert.NoError(t, err)
		assert.Greater(t, size, int64(0))
		
		diffID, err := layer.DiffID()
		assert.NoError(t, err)
		assert.NotEmpty(t, diffID.String())
		
		// Get layer info
		info, err := GetLayerInfo(layer)
		assert.NoError(t, err)
		assert.Equal(t, size, info.Size)
		assert.Equal(t, diffID, info.DiffID)
	})
	
	t.Run("empty context dir", func(t *testing.T) {
		_, err := factory.CreateLayer("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "context directory cannot be empty")
	})
	
	t.Run("nonexistent directory", func(t *testing.T) {
		_, err := factory.CreateLayer("/nonexistent")
		assert.Error(t, err)
	})
}

// TestLayerFactoryWithOptions tests layer factory with options
func TestLayerFactoryWithOptions(t *testing.T) {
	factory := NewLayerFactory().
		WithPermissionPreservation(false).
		WithTimestampPreservation(true)
	
	assert.NotNil(t, factory)
	assert.False(t, factory.preservePermissions)
	assert.True(t, factory.preserveTimestamps)
}