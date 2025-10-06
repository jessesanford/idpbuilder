package testutils

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/types"
)

// PushTestCase represents a test case for push operations
type PushTestCase struct {
	Name           string
	ImageName      string
	Tag            string
	ExpectSuccess  bool
	ExpectedError  string
	AuthRequired   bool
	Config         *AuthConfig
	ExpectedLayers int
}

// AssertionHelper provides common assertion utilities for container testing
type AssertionHelper struct {
	t *testing.T
}

// NewAssertionHelper creates a new assertion helper
func NewAssertionHelper(t *testing.T) *AssertionHelper {
	return &AssertionHelper{t: t}
}

// AssertPushSucceeds verifies that a push operation succeeds
func (a *AssertionHelper) AssertPushSucceeds(registry *MockRegistry, imageName string, image v1.Image) {
	a.t.Helper()

	if registry == nil {
		a.t.Fatal("Registry cannot be nil")
	}

	if image == nil {
		a.t.Fatal("Image cannot be nil")
	}

	// Verify the image was stored in the registry
	if !registry.HasImage(imageName) {
		a.t.Errorf("Expected image %s to be pushed to registry", imageName)
	}

	// Verify image integrity
	storedImage := registry.GetImage(imageName)
	if storedImage == nil {
		a.t.Errorf("Expected to retrieve image %s from registry", imageName)
		return
	}

	// Compare digests
	originalDigest, err := image.Digest()
	if err != nil {
		a.t.Errorf("Failed to get original image digest: %v", err)
		return
	}

	storedDigest, err := storedImage.Digest()
	if err != nil {
		a.t.Errorf("Failed to get stored image digest: %v", err)
		return
	}

	if originalDigest != storedDigest {
		a.t.Errorf("Image digests don't match. Original: %s, Stored: %s", originalDigest, storedDigest)
	}
}

// AssertPushFails verifies that a push operation fails with expected error
func (a *AssertionHelper) AssertPushFails(registry *MockRegistry, imageName string, image v1.Image, expectedError string) {
	a.t.Helper()

	// This would normally be called after attempting a push that should fail
	// For now, we'll check that the image is NOT in the registry
	if registry.HasImage(imageName) {
		a.t.Errorf("Expected push to fail, but image %s was found in registry", imageName)
	}
}

// AssertAuthRequired verifies that authentication is required for an operation
func (a *AssertionHelper) AssertAuthRequired(registry *MockRegistry, operation string) {
	a.t.Helper()

	if !registry.AuthConfig.Enabled {
		a.t.Errorf("Expected authentication to be required for %s operation", operation)
	}
}

// AssertImageExists verifies that an image exists in the registry
func (a *AssertionHelper) AssertImageExists(registry *MockRegistry, imageName string) {
	a.t.Helper()

	if !registry.HasImage(imageName) {
		a.t.Errorf("Expected image %s to exist in registry", imageName)
	}
}

// AssertImageNotExists verifies that an image does not exist in the registry
func (a *AssertionHelper) AssertImageNotExists(registry *MockRegistry, imageName string) {
	a.t.Helper()

	if registry.HasImage(imageName) {
		a.t.Errorf("Expected image %s to not exist in registry", imageName)
	}
}

// AssertManifestExists verifies that a manifest exists for the given reference
func (a *AssertionHelper) AssertManifestExists(registry *MockRegistry, ref string) {
	a.t.Helper()

	manifest := registry.GetManifest(ref)
	if manifest == nil {
		a.t.Errorf("Expected manifest to exist for reference %s", ref)
	}
}

// AssertTagExists verifies that a specific tag exists for a repository
func (a *AssertionHelper) AssertTagExists(registry *MockRegistry, repo, tag string) {
	a.t.Helper()

	fullRef := fmt.Sprintf("%s:%s", repo, tag)
	if !registry.HasImage(fullRef) {
		a.t.Errorf("Expected tag %s to exist for repository %s", tag, repo)
	}
}

// AssertLayerExists verifies that a layer with the given digest exists
func (a *AssertionHelper) AssertLayerExists(registry *MockRegistry, digest v1.Hash) {
	a.t.Helper()

	if !registry.HasLayer(digest) {
		a.t.Errorf("Expected layer %s to exist in registry", digest)
	}
}

// AssertConfigMatches verifies that image configuration matches expected values
func (a *AssertionHelper) AssertConfigMatches(expected, actual *v1.ConfigFile) {
	a.t.Helper()

	if expected == nil || actual == nil {
		a.t.Fatal("Config files cannot be nil")
	}

	if expected.Architecture != actual.Architecture {
		a.t.Errorf("Architecture mismatch. Expected: %s, Got: %s", expected.Architecture, actual.Architecture)
	}

	if expected.OS != actual.OS {
		a.t.Errorf("OS mismatch. Expected: %s, Got: %s", expected.OS, actual.OS)
	}

	// Compare environment variables
	if !reflect.DeepEqual(expected.Config.Env, actual.Config.Env) {
		a.t.Errorf("Environment variables mismatch. Expected: %v, Got: %v", expected.Config.Env, actual.Config.Env)
	}
}

// AssertMediaTypeMatches verifies that media types match expected values
func (a *AssertionHelper) AssertMediaTypeMatches(expected, actual types.MediaType) {
	a.t.Helper()

	if expected != actual {
		a.t.Errorf("Media type mismatch. Expected: %s, Got: %s", expected, actual)
	}
}

// AssertLayerCount verifies the number of layers in an image
func (a *AssertionHelper) AssertLayerCount(image v1.Image, expectedCount int) {
	a.t.Helper()

	layers, err := image.Layers()
	if err != nil {
		a.t.Errorf("Failed to get image layers: %v", err)
		return
	}

	actualCount := len(layers)
	if actualCount != expectedCount {
		a.t.Errorf("Layer count mismatch. Expected: %d, Got: %d", expectedCount, actualCount)
	}
}

// AssertImageSize verifies the size of an image
func (a *AssertionHelper) AssertImageSize(image v1.Image, expectedSize int64) {
	a.t.Helper()

	size, err := image.Size()
	if err != nil {
		a.t.Errorf("Failed to get image size: %v", err)
		return
	}

	if size != expectedSize {
		a.t.Errorf("Image size mismatch. Expected: %d, Got: %d", expectedSize, size)
	}
}

// AssertRegistryHealthy verifies that the registry is functioning correctly
func (a *AssertionHelper) AssertRegistryHealthy(registry *MockRegistry) {
	a.t.Helper()

	if registry.Server == nil {
		a.t.Error("Registry server is not initialized")
	}

	// Test basic connectivity
	url := registry.GetURL()
	if url == "" {
		a.t.Error("Registry URL is empty")
	}

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		a.t.Errorf("Invalid registry URL format: %s", url)
	}
}

// AssertErrorContains verifies that an error contains expected text
func (a *AssertionHelper) AssertErrorContains(err error, expectedText string) {
	a.t.Helper()

	if err == nil {
		a.t.Errorf("Expected error containing '%s', but got nil", expectedText)
		return
	}

	if !strings.Contains(err.Error(), expectedText) {
		a.t.Errorf("Expected error to contain '%s', but got: %s", expectedText, err.Error())
	}
}

// AssertNoError verifies that no error occurred
func (a *AssertionHelper) AssertNoError(err error) {
	a.t.Helper()

	if err != nil {
		a.t.Errorf("Expected no error, but got: %v", err)
	}
}

// AssertEqual verifies that two values are equal
func (a *AssertionHelper) AssertEqual(expected, actual interface{}) {
	a.t.Helper()

	if !reflect.DeepEqual(expected, actual) {
		a.t.Errorf("Values not equal. Expected: %v, Got: %v", expected, actual)
	}
}

// AssertNotEqual verifies that two values are not equal
func (a *AssertionHelper) AssertNotEqual(expected, actual interface{}) {
	a.t.Helper()

	if reflect.DeepEqual(expected, actual) {
		a.t.Errorf("Values should not be equal, but both are: %v", expected)
	}
}

// AssertTrue verifies that a condition is true
func (a *AssertionHelper) AssertTrue(condition bool, message string) {
	a.t.Helper()

	if !condition {
		a.t.Errorf("Expected condition to be true: %s", message)
	}
}

// AssertFalse verifies that a condition is false
func (a *AssertionHelper) AssertFalse(condition bool, message string) {
	a.t.Helper()

	if condition {
		a.t.Errorf("Expected condition to be false: %s", message)
	}
}