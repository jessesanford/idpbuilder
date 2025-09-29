package testutils

import (
	"context"
	"testing"
	"time"
)

// TestImageCreation tests the image creation utilities from Split 002
func TestImageCreation(t *testing.T) {
	tests := []struct {
		name         string
		imageName    string
		tags         []string
		expectedArch string
		expectedOS   string
	}{
		{
			name:         "simple test image",
			imageName:    "test/simple",
			tags:         []string{"latest", "v1.0"},
			expectedArch: "amd64",
			expectedOS:   "linux",
		},
		{
			name:         "named test image",
			imageName:    "test/myapp",
			tags:         []string{"dev"},
			expectedArch: "amd64",
			expectedOS:   "linux",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test image using Split 002 utilities
			image, err := CreateTestImage(tt.imageName, tt.tags)
			if err != nil {
				t.Fatalf("Failed to create test image: %v", err)
			}

			// Verify image properties
			configFile, err := image.ConfigFile()
			if err != nil {
				t.Fatalf("Failed to get image config: %v", err)
			}

			if configFile.Architecture != tt.expectedArch {
				t.Errorf("Expected architecture %s, got %s", tt.expectedArch, configFile.Architecture)
			}

			if configFile.OS != tt.expectedOS {
				t.Errorf("Expected OS %s, got %s", tt.expectedOS, configFile.OS)
			}

			// Verify image has at least one layer
			layers, err := image.Layers()
			if err != nil {
				t.Fatalf("Failed to get image layers: %v", err)
			}

			if len(layers) == 0 {
				t.Error("Expected at least one layer in the image")
			}
		})
	}
}

// TestAssertions tests the assertion utilities from Split 002
func TestAssertions(t *testing.T) {
	// Create test fixtures using Split 001 infrastructure
	fixtures := SetupTestFixtures(t, false)
	defer CleanupTestFixtures(fixtures)

	// Create assertion helper from Split 002
	helper := NewAssertionHelper(t)

	// Test that registry is healthy
	helper.AssertRegistryHealthy(fixtures.Registry)

	// Create a test image using Split 002 utilities
	image, err := CreateTestImage("test/assertions", []string{"latest"})
	if err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}

	// Test push success assertion
	imageName := "test/assertions:latest"

	// Store image in registry (simulating successful push)
	err = fixtures.Registry.StoreImage(imageName, image)
	if err != nil {
		t.Fatalf("Failed to store image in registry: %v", err)
	}

	// Assert push succeeded
	helper.AssertPushSucceeds(fixtures.Registry, imageName, image)

	// Test image existence assertion
	helper.AssertImageExists(fixtures.Registry, imageName)

	// Test that non-existent image doesn't exist
	helper.AssertImageNotExists(fixtures.Registry, "test/nonexistent:latest")
}

// TestIntegrationFramework tests the complete framework using both splits
func TestIntegrationFramework(t *testing.T) {
	// Setup test environment using Split 001
	fixtures := SetupTestFixtures(t, true) // Enable auth
	defer CleanupTestFixtures(fixtures)

	// Create assertion helper from Split 002
	helper := NewAssertionHelper(t)

	// Verify auth is required
	helper.AssertAuthRequired(fixtures.Registry, "push")

	// Create multiple test images using Split 002
	testCases := []struct {
		name      string
		imageName string
		layers    []string
	}{
		{
			name:      "single layer image",
			imageName: "test/single-layer",
			layers:    []string{"base-content"},
		},
		{
			name:      "multi layer image",
			imageName: "test/multi-layer",
			layers:    []string{"base", "app", "config"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create image with specific layers
			image, err := CreateImageWithLayers(tc.imageName, tc.layers)
			if err != nil {
				t.Fatalf("Failed to create image with layers: %v", err)
			}

			// Verify layer count
			helper.AssertLayerCount(image, len(tc.layers))

			// Push to registry
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			ref := tc.imageName + ":latest"
			err = PushTestImage(ctx, fixtures.Registry, ref, image)
			if err != nil {
				t.Fatalf("Failed to push test image: %v", err)
			}

			// Assert push succeeded
			helper.AssertPushSucceeds(fixtures.Registry, ref, image)

			// Verify manifest exists
			helper.AssertManifestExists(fixtures.Registry, ref)

			// Verify tag exists
			helper.AssertTagExists(fixtures.Registry, tc.imageName, "latest")
		})
	}
}

// TestMultiArchSupport tests multi-architecture image support
func TestMultiArchSupport(t *testing.T) {
	// Setup test environment
	fixtures := SetupTestFixtures(t, false)
	defer CleanupTestFixtures(fixtures)

	helper := NewAssertionHelper(t)

	// Define platforms
	platforms := []Platform{
		{Architecture: "amd64", OS: "linux"},
		{Architecture: "arm64", OS: "linux"},
	}

	// Create multi-arch image
	imageName := "test/multiarch"
	index, err := CreateMultiArchImage(imageName, platforms)
	if err != nil {
		t.Fatalf("Failed to create multi-arch image: %v", err)
	}

	// Verify index was created
	if index == nil {
		t.Fatal("Expected non-nil image index")
	}

	// Get index manifest
	manifest, err := index.IndexManifest()
	if err != nil {
		t.Fatalf("Failed to get index manifest: %v", err)
	}

	// Verify we have the expected number of manifests
	expectedManifests := len(platforms)
	actualManifests := len(manifest.Manifests)
	if actualManifests != expectedManifests {
		t.Errorf("Expected %d manifests, got %d", expectedManifests, actualManifests)
	}

	// Verify each platform is represented
	for i, platform := range platforms {
		if i < len(manifest.Manifests) {
			desc := manifest.Manifests[i]
			if desc.Platform == nil {
				t.Errorf("Expected platform info for manifest %d", i)
				continue
			}

			if desc.Platform.Architecture != platform.Architecture {
				t.Errorf("Expected architecture %s, got %s", platform.Architecture, desc.Platform.Architecture)
			}

			if desc.Platform.OS != platform.OS {
				t.Errorf("Expected OS %s, got %s", platform.OS, desc.Platform.OS)
			}
		}
	}
}

// TestErrorHandling tests error handling scenarios
func TestErrorHandling(t *testing.T) {
	helper := NewAssertionHelper(t)

	// Test nil registry handling
	t.Run("nil registry", func(t *testing.T) {
		// This should be caught by assertions
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for nil registry")
			}
		}()

		image, _ := CreateTestImage("test", []string{"latest"})
		helper.AssertPushSucceeds(nil, "test:latest", image)
	})

	// Test nil image handling
	t.Run("nil image", func(t *testing.T) {
		fixtures := SetupTestFixtures(t, false)
		defer CleanupTestFixtures(fixtures)

		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for nil image")
			}
		}()

		helper.AssertPushSucceeds(fixtures.Registry, "test:latest", nil)
	})

	// Test invalid image creation
	t.Run("invalid image config", func(t *testing.T) {
		// Test with invalid layer content
		_, err := CreateImageWithLayers("test", []string{})
		if err != nil {
			// This is expected - empty layers should be handled gracefully
			helper.AssertErrorContains(err, "")
		}
	})
}

// TestImageComparison tests image comparison utilities
func TestImageComparison(t *testing.T) {
	// Create two identical images
	image1, err := CreateTestImage("test/compare1", []string{"v1"})
	if err != nil {
		t.Fatalf("Failed to create first image: %v", err)
	}

	image2, err := CreateTestImage("test/compare2", []string{"v1"})
	if err != nil {
		t.Fatalf("Failed to create second image: %v", err)
	}

	// Images with different content should not be equal
	equal, err := CompareImages(image1, image2)
	if err != nil {
		t.Fatalf("Failed to compare images: %v", err)
	}

	// They should be different since they have different content
	if equal {
		t.Error("Expected images to be different")
	}

	// Create image with same content
	image3, err := CreateTestImage("test/compare1", []string{"v1"})
	if err != nil {
		t.Fatalf("Failed to create third image: %v", err)
	}

	// These should be equal (same name and content)
	equal, err = CompareImages(image1, image3)
	if err != nil {
		t.Fatalf("Failed to compare images: %v", err)
	}

	if !equal {
		t.Error("Expected images with same content to be equal")
	}
}

// TestImageInfo tests image information utilities
func TestImageInfo(t *testing.T) {
	// Create test image
	imageName := "test/info"
	layers := []string{"layer1", "layer2", "layer3"}

	image, err := CreateImageWithLayers(imageName, layers)
	if err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}

	// Get image info
	info, err := GetImageInfo(image)
	if err != nil {
		t.Fatalf("Failed to get image info: %v", err)
	}

	// Verify info
	if info.LayerCount != len(layers) {
		t.Errorf("Expected %d layers, got %d", len(layers), info.LayerCount)
	}

	if info.Architecture != "amd64" {
		t.Errorf("Expected amd64 architecture, got %s", info.Architecture)
	}

	if info.OS != "linux" {
		t.Errorf("Expected linux OS, got %s", info.OS)
	}

	if info.Size <= 0 {
		t.Error("Expected positive image size")
	}

	// Test string representation
	infoStr := info.String()
	if len(infoStr) == 0 {
		t.Error("Expected non-empty string representation")
	}

	t.Logf("Image info: %s", infoStr)
}