package push

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// PushTestCase represents a test case for push operations
type PushTestCase struct {
	Name     string    // Test case name
	Image    v1.Image  // Image to push
	WantErr  bool      // Expected error
	ErrorMsg string    // Expected error message
}

// AssertPushSucceeds verifies that a push operation completes successfully
func AssertPushSucceeds(t *testing.T, fixtures *TestFixtures, img v1.Image, imageName string) {
	require.NotNil(t, fixtures, "Test fixtures must not be nil")
	require.NotNil(t, fixtures.Registry, "Mock registry must not be nil")
	require.NotNil(t, img, "Image must not be nil")

	// Get registry URL
	registryURL := fixtures.Registry.GetURL()
	require.NotEmpty(t, registryURL, "Registry URL must not be empty")

	// Construct manifest URL
	manifestURL := fmt.Sprintf("%s/v2/%s/manifests/latest", registryURL, imageName)

	// Get manifest from image
	manifest, err := img.RawManifest()
	require.NoError(t, err, "Failed to get image manifest")

	// Create PUT request to push manifest
	req, err := http.NewRequest("PUT", manifestURL, strings.NewReader(string(manifest)))
	require.NoError(t, err, "Failed to create PUT request")

	req.Header.Set("Content-Type", "application/vnd.docker.distribution.manifest.v2+json")

	// Execute push request
	resp, err := fixtures.Client.Do(req)
	require.NoError(t, err, "Push request failed")
	defer resp.Body.Close()

	// Assert successful push
	assert.Equal(t, http.StatusCreated, resp.StatusCode,
		"Push should return 201 Created, got %d", resp.StatusCode)

	// Verify image was stored in registry
	storedImages := fixtures.Registry.GetStoredImages()
	imageKey := fmt.Sprintf("%s:latest", imageName)
	assert.Contains(t, storedImages, imageKey,
		"Image should be stored in registry with key %s", imageKey)

	// Verify stored manifest matches original
	storedManifest := storedImages[imageKey]
	assert.Equal(t, manifest, storedManifest,
		"Stored manifest should match original")
}

// AssertAuthRequired verifies that authentication is required for the endpoint
func AssertAuthRequired(t *testing.T, fixtures *TestFixtures, endpoint string) {
	require.NotNil(t, fixtures, "Test fixtures must not be nil")
	require.NotNil(t, fixtures.Registry, "Mock registry must not be nil")
	require.NotEmpty(t, endpoint, "Endpoint must not be empty")

	// Ensure auth is required on the registry
	if fixtures.Registry.AuthConfig == nil {
		fixtures.Registry.SetAuth("testuser", "testpass", true)
	}

	registryURL := fixtures.Registry.GetURL()
	require.NotEmpty(t, registryURL, "Registry URL must not be empty")

	// Construct full URL
	fullURL := fmt.Sprintf("%s%s", registryURL, endpoint)

	// Create client without authentication
	client := &http.Client{
		Transport: http.DefaultTransport,
	}

	// Make request without auth
	req, err := http.NewRequest("GET", fullURL, nil)
	require.NoError(t, err, "Failed to create request")

	resp, err := client.Do(req)
	require.NoError(t, err, "Request failed")
	defer resp.Body.Close()

	// Assert unauthorized response
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode,
		"Request without auth should return 401 Unauthorized, got %d", resp.StatusCode)

	// Check for WWW-Authenticate header
	authHeader := resp.Header.Get("WWW-Authenticate")
	assert.NotEmpty(t, authHeader, "Response should include WWW-Authenticate header")
	assert.Contains(t, authHeader, "Bearer", "WWW-Authenticate should specify Bearer auth")
}

// AssertPushFails verifies that a push operation fails as expected
func AssertPushFails(t *testing.T, fixtures *TestFixtures, img v1.Image, imageName string, expectedError string) {
	require.NotNil(t, fixtures, "Test fixtures must not be nil")
	require.NotNil(t, fixtures.Registry, "Mock registry must not be nil")
	require.NotNil(t, img, "Image must not be nil")

	// Get registry URL
	registryURL := fixtures.Registry.GetURL()
	require.NotEmpty(t, registryURL, "Registry URL must not be empty")

	// Construct manifest URL
	manifestURL := fmt.Sprintf("%s/v2/%s/manifests/latest", registryURL, imageName)

	// Get manifest from image
	manifest, err := img.RawManifest()
	require.NoError(t, err, "Failed to get image manifest")

	// Create PUT request with potentially problematic data
	req, err := http.NewRequest("PUT", manifestURL, strings.NewReader(string(manifest)))
	require.NoError(t, err, "Failed to create PUT request")

	req.Header.Set("Content-Type", "application/vnd.docker.distribution.manifest.v2+json")

	// Execute push request
	resp, err := fixtures.Client.Do(req)

	// Assert that either the request failed or returned an error status
	if err != nil {
		assert.Contains(t, err.Error(), expectedError,
			"Error message should contain expected error: %s", expectedError)
	} else {
		defer resp.Body.Close()
		assert.True(t, resp.StatusCode >= 400,
			"Response should be an error status (>=400), got %d", resp.StatusCode)
	}
}

// AssertRegistryHealth verifies that the mock registry is responding correctly
func AssertRegistryHealth(t *testing.T, fixtures *TestFixtures) {
	require.NotNil(t, fixtures, "Test fixtures must not be nil")
	require.NotNil(t, fixtures.Registry, "Mock registry must not be nil")

	registryURL := fixtures.Registry.GetURL()
	require.NotEmpty(t, registryURL, "Registry URL must not be empty")

	// Check registry version endpoint
	versionURL := fmt.Sprintf("%s/v2/", registryURL)

	resp, err := fixtures.Client.Get(versionURL)
	require.NoError(t, err, "Registry health check failed")
	defer resp.Body.Close()

	// Assert registry responds with correct headers
	assert.Equal(t, http.StatusOK, resp.StatusCode,
		"Registry should return 200 OK for version check")

	apiVersion := resp.Header.Get("Docker-Distribution-API-Version")
	assert.Equal(t, "registry/2.0", apiVersion,
		"Registry should return correct API version header")
}

// AssertImageExists verifies that an image exists in the registry
func AssertImageExists(t *testing.T, fixtures *TestFixtures, imageName, tag string) {
	require.NotNil(t, fixtures, "Test fixtures must not be nil")
	require.NotNil(t, fixtures.Registry, "Mock registry must not be nil")

	registryURL := fixtures.Registry.GetURL()
	require.NotEmpty(t, registryURL, "Registry URL must not be empty")

	// Construct manifest URL for HEAD request
	manifestURL := fmt.Sprintf("%s/v2/%s/manifests/%s", registryURL, imageName, tag)

	// Create HEAD request to check if manifest exists
	req, err := http.NewRequest("HEAD", manifestURL, nil)
	require.NoError(t, err, "Failed to create HEAD request")

	resp, err := fixtures.Client.Do(req)
	require.NoError(t, err, "HEAD request failed")
	defer resp.Body.Close()

	// Assert image exists
	assert.Equal(t, http.StatusOK, resp.StatusCode,
		"Image %s:%s should exist in registry", imageName, tag)
}

// AssertImageNotExists verifies that an image does not exist in the registry
func AssertImageNotExists(t *testing.T, fixtures *TestFixtures, imageName, tag string) {
	require.NotNil(t, fixtures, "Test fixtures must not be nil")
	require.NotNil(t, fixtures.Registry, "Mock registry must not be nil")

	registryURL := fixtures.Registry.GetURL()
	require.NotEmpty(t, registryURL, "Registry URL must not be empty")

	// Construct manifest URL for HEAD request
	manifestURL := fmt.Sprintf("%s/v2/%s/manifests/%s", registryURL, imageName, tag)

	// Create HEAD request to check if manifest exists
	req, err := http.NewRequest("HEAD", manifestURL, nil)
	require.NoError(t, err, "Failed to create HEAD request")

	resp, err := fixtures.Client.Do(req)
	require.NoError(t, err, "HEAD request failed")
	defer resp.Body.Close()

	// Assert image does not exist
	assert.Equal(t, http.StatusNotFound, resp.StatusCode,
		"Image %s:%s should not exist in registry", imageName, tag)
}

// RunPushTestCases executes a series of push test cases
func RunPushTestCases(t *testing.T, fixtures *TestFixtures, testCases []PushTestCase) {
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			if tc.WantErr {
				AssertPushFails(t, fixtures, tc.Image, tc.Name, tc.ErrorMsg)
			} else {
				AssertPushSucceeds(t, fixtures, tc.Image, tc.Name)
			}
		})
	}
}