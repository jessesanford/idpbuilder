package builder

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-containerregistry/pkg/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTarballWriter tests tarball export functionality
func TestTarballWriter(t *testing.T) {
	writer := NewTarballWriter()
	assert.NotNil(t, writer)
	
	// Create a simple test image
	tempDir, err := os.MkdirTemp("", "tarball-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	testFile := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	require.NoError(t, err)
	
	// Build an image for testing
	ctx := context.Background()
	builder, err := NewBuilder(DefaultBuildOptions())
	require.NoError(t, err)
	
	image, err := builder.Build(ctx, tempDir, DefaultBuildOptions())
	require.NoError(t, err)
	
	// Test tarball writing
	t.Run("write tarball", func(t *testing.T) {
		outputPath := filepath.Join(tempDir, "test.tar")
		
		err := writer.Write(image, outputPath, "localhost/test:latest")
		assert.NoError(t, err)
		
		// Verify file was created
		info, err := os.Stat(outputPath)
		assert.NoError(t, err)
		assert.Greater(t, info.Size(), int64(0))
	})
	
	t.Run("invalid parameters", func(t *testing.T) {
		outputPath := filepath.Join(tempDir, "test.tar")
		
		// Nil image
		err := writer.Write(nil, outputPath, "test:latest")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "image cannot be nil")
		
		// Empty output path
		err = writer.Write(image, "", "test:latest")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "output path cannot be empty")
		
		// Empty reference
		err = writer.Write(image, outputPath, "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "image reference cannot be empty")
	})
}

// TestTarballWriterWithOptions tests tarball writer configuration
func TestTarballWriterWithOptions(t *testing.T) {
	t.Run("with options", func(t *testing.T) {
		opts := TarballOptions{
			Compress:        true,
			IncludeManifest: false,
		}
		writer := NewTarballWriterWithOptions(opts)
		assert.NotNil(t, writer)
		assert.True(t, writer.options.Compress)
		assert.False(t, writer.options.IncludeManifest)
	})
	
	t.Run("with compression", func(t *testing.T) {
		writer := NewTarballWriter().WithCompression(true)
		assert.True(t, writer.options.Compress)
	})
}

// TestTarballWriterMultiple tests multi-image tarball export
func TestTarballWriterMultiple(t *testing.T) {
	writer := NewTarballWriter()
	
	t.Run("empty images map", func(t *testing.T) {
		outputPath := filepath.Join(os.TempDir(), "empty.tar")
		err := writer.WriteMultiple(map[string]v1.Image{}, outputPath)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no images provided for export")
	})
	
	t.Run("invalid reference", func(t *testing.T) {
		outputPath := filepath.Join(os.TempDir(), "invalid.tar")
		images := map[string]v1.Image{
			"invalid::reference": nil,
		}
		err := writer.WriteMultiple(images, outputPath)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid reference")
	})
}

// TestTarballInfo tests tarball information functions
func TestTarballInfo(t *testing.T) {
	t.Run("nonexistent file", func(t *testing.T) {
		_, err := GetTarballInfo("/nonexistent/file.tar")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tarball file not found")
	})
	
	t.Run("valid file", func(t *testing.T) {
		// Create a temporary file
		tempFile, err := os.CreateTemp("", "test*.tar")
		require.NoError(t, err)
		defer os.Remove(tempFile.Name())
		
		// Write some content
		_, err = tempFile.WriteString("test content")
		require.NoError(t, err)
		tempFile.Close()
		
		info, err := GetTarballInfo(tempFile.Name())
		assert.NoError(t, err)
		assert.Equal(t, tempFile.Name(), info.Path)
		assert.Greater(t, info.Size, int64(0))
	})
}

// TestValidateTarball tests tarball validation
func TestValidateTarball(t *testing.T) {
	t.Run("nonexistent file", func(t *testing.T) {
		err := ValidateTarball("/nonexistent/file.tar")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot open tarball file")
	})
	
	t.Run("empty file", func(t *testing.T) {
		tempFile, err := os.CreateTemp("", "empty*.tar")
		require.NoError(t, err)
		defer os.Remove(tempFile.Name())
		tempFile.Close()
		
		err = ValidateTarball(tempFile.Name())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tarball file is empty")
	})
	
	t.Run("valid file", func(t *testing.T) {
		tempFile, err := os.CreateTemp("", "valid*.tar")
		require.NoError(t, err)
		defer os.Remove(tempFile.Name())
		
		_, err = tempFile.WriteString("valid content")
		require.NoError(t, err)
		tempFile.Close()
		
		err = ValidateTarball(tempFile.Name())
		assert.NoError(t, err)
	})
}

// TestCompressTarball tests tarball compression
func TestCompressTarball(t *testing.T) {
	err := CompressTarball("input.tar", "output.tar.gz")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tarball compression not yet implemented")
}

// Benchmark tests for performance validation
func BenchmarkBuild(b *testing.B) {
	// Create test directory
	tempDir, err := os.MkdirTemp("", "bench-*")
	require.NoError(b, err)
	defer os.RemoveAll(tempDir)
	
	// Create some test files
	for i := 0; i < 10; i++ {
		err := os.WriteFile(filepath.Join(tempDir, "file"+string(rune(i))+".txt"), 
			[]byte("content "+string(rune(i))), 0644)
		require.NoError(b, err)
	}
	
	builder, err := NewBuilder(DefaultBuildOptions())
	require.NoError(b, err)
	
	ctx := context.Background()
	opts := DefaultBuildOptions()
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, err := builder.Build(ctx, tempDir, opts)
		if err != nil {
			b.Fatal(err)
		}
	}
}