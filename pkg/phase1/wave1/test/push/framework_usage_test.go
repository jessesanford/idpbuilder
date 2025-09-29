package push

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFrameworkUsageScenarios tests the framework in realistic usage scenarios
func TestFrameworkUsageScenarios(t *testing.T) {
	t.Run("BasicPushWorkflow", func(t *testing.T) {
		// Set up test environment
		fixtures := SetupTestFixtures(t)

		// Create a test image
		image := CreateTestImage("test-app", "v1.0.0")

		// Test registry health
		AssertRegistryHealth(t, fixtures)

		// Verify image doesn't exist initially
		AssertImageNotExists(t, fixtures, "test-app", "latest")

		// Push the image
		AssertPushSucceeds(t, fixtures, image, "test-app")

		// Verify image now exists
		AssertImageExists(t, fixtures, "test-app", "latest")

		// Test stored images retrieval
		stored := fixtures.Registry.GetStoredImages()
		assert.Contains(t, stored, "test-app:latest", "Image should be in stored images")
	})

	t.Run("AuthenticationWorkflow", func(t *testing.T) {
		// Set up authenticated environment
		fixtures := CreateAuthenticatedFixtures(t, "testuser", "testpass")

		// Test that auth is required
		AssertAuthRequired(t, fixtures, "/v2/")

		// Test that authenticated push works
		image := CreateTestImage("secure-app", "v1.0.0")
		AssertPushSucceeds(t, fixtures, image, "secure-app")

		// Verify image exists
		AssertImageExists(t, fixtures, "secure-app", "latest")
	})

	t.Run("TestImageInterfaceMethods", func(t *testing.T) {
		// Create test image
		image := CreateTestImage("test-methods", "v2.0.0")

		// Test various image methods for coverage
		digest, err := image.Digest()
		require.NoError(t, err)
		assert.NotEmpty(t, digest.Hex)

		rawManifest, err := image.RawManifest()
		require.NoError(t, err)
		assert.NotEmpty(t, rawManifest)

		size, err := image.Size()
		require.NoError(t, err)
		assert.Greater(t, size, int64(0))

		configName, err := image.ConfigName()
		require.NoError(t, err)
		assert.NotEmpty(t, configName.Hex)

		rawConfig, err := image.RawConfigFile()
		require.NoError(t, err)
		assert.NotEmpty(t, rawConfig)

		mediaType, err := image.MediaType()
		require.NoError(t, err)
		assert.NotEmpty(t, string(mediaType))

		// Test layer methods
		layers, err := image.Layers()
		require.NoError(t, err)
		require.Len(t, layers, 1)

		layer := layers[0]
		layerDigest, err := layer.Digest()
		require.NoError(t, err)
		assert.NotEmpty(t, layerDigest.Hex)

		diffID, err := layer.DiffID()
		require.NoError(t, err)
		assert.NotEmpty(t, diffID.Hex)

		compressed, err := layer.Compressed()
		require.NoError(t, err)
		assert.NotNil(t, compressed)
		compressed.Close()

		uncompressed, err := layer.Uncompressed()
		require.NoError(t, err)
		assert.NotNil(t, uncompressed)
		uncompressed.Close()

		layerSize, err := layer.Size()
		require.NoError(t, err)
		assert.Greater(t, layerSize, int64(0))

		layerMediaType, err := layer.MediaType()
		require.NoError(t, err)
		assert.NotEmpty(t, string(layerMediaType))

		// Test layer by digest lookup
		foundLayer, err := image.LayerByDigest(layerDigest)
		require.NoError(t, err)
		assert.NotNil(t, foundLayer)

		// Test layer by diff ID lookup
		foundLayer2, err := image.LayerByDiffID(diffID)
		require.NoError(t, err)
		assert.NotNil(t, foundLayer2)
	})

	t.Run("TestFileOperations", func(t *testing.T) {
		fixtures := SetupTestFixtures(t)

		// Test file write and read
		testContent := []byte("test file content")
		filename := "test-file.txt"

		filePath, err := WriteTestFile(fixtures, filename, testContent)
		require.NoError(t, err)
		assert.Contains(t, filePath, filename)

		readContent, err := ReadTestFile(fixtures, filename)
		require.NoError(t, err)
		assert.Equal(t, testContent, readContent)
	})

	t.Run("TestRegistryErrorGeneration", func(t *testing.T) {
		// Test registry error creation
		errorData := CreateRegistryError(404, "MANIFEST_UNKNOWN", "manifest not found")
		assert.NotEmpty(t, errorData)

		// This validates the error structure is properly formatted
		assert.Contains(t, string(errorData), "MANIFEST_UNKNOWN")
		assert.Contains(t, string(errorData), "manifest not found")
	})

	t.Run("TestPushTestCasesExecution", func(t *testing.T) {
		fixtures := SetupTestFixtures(t)

		// Create test cases - first test should succeed with no auth
		testCases := []PushTestCase{
			{
				Name:    "successful-push",
				Image:   CreateTestImage("success-app", "v1.0.0"),
				WantErr: false,
			},
		}

		// Run successful test case without auth
		RunPushTestCases(t, fixtures, testCases)

		// Now test with auth required to demonstrate failure scenario
		fixtures.Registry.SetAuth("admin", "secret", true)

		failureTestCases := []PushTestCase{
			{
				Name:     "auth-required-failure",
				Image:    CreateTestImage("fail-app", "v1.0.0"),
				WantErr:  true,
				ErrorMsg: "401", // Will fail with 401 status
			},
		}

		// Run failure test case with auth required
		RunPushTestCases(t, fixtures, failureTestCases)
	})
}