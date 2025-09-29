package format

import (
	"errors"
	"fmt"
)

// NewPackageManifest creates a new package manifest with the specified name and version
// Returns a properly initialized PackageManifest with required fields set
func NewPackageManifest(name, version string) *PackageManifest {
	if name == "" {
		return nil
	}
	if version == "" {
		return nil
	}

	manifest := &PackageManifest{
		SchemaVersion: SchemaVersion,
		MediaType:     MediaTypes.ManifestV2,
		Layers:        make([]LayerDescriptor, 0),
		Annotations: map[string]string{
			"org.opencontainers.image.title":   name,
			"org.opencontainers.image.version": version,
		},
	}

	return manifest
}

// ValidateManifest validates the structure and content of a package manifest
// Returns an error if the manifest is invalid according to OCI specification
func ValidateManifest(manifest *PackageManifest) error {
	if manifest == nil {
		return errors.New("manifest cannot be nil")
	}

	if manifest.SchemaVersion != SchemaVersion {
		return fmt.Errorf("invalid schema version: expected %d, got %d", SchemaVersion, manifest.SchemaVersion)
	}

	if manifest.MediaType == "" {
		return errors.New("media type is required")
	}

	if manifest.Config.MediaType == "" {
		return errors.New("config media type is required")
	}

	if manifest.Config.Digest == "" {
		return errors.New("config digest is required")
	}

	if manifest.Config.Size <= 0 {
		return errors.New("config size must be positive")
	}

	return nil
}

// AddLayer adds a layer descriptor to the package manifest
// Returns an error if the layer is invalid or cannot be added
func AddLayer(manifest *PackageManifest, layer *LayerDescriptor) error {
	if manifest == nil {
		return errors.New("manifest cannot be nil")
	}

	if layer == nil {
		return errors.New("layer cannot be nil")
	}

	if layer.MediaType == "" {
		return errors.New("layer media type is required")
	}

	if layer.Digest == "" {
		return errors.New("layer digest is required")
	}

	if layer.Size <= 0 {
		return errors.New("layer size must be positive")
	}

	manifest.Layers = append(manifest.Layers, *layer)
	return nil
}

// GetConfigDescriptor returns the config descriptor from the package manifest
// Returns a pointer to the config descriptor for read-only access
func GetConfigDescriptor(manifest *PackageManifest) *Descriptor {
	if manifest == nil {
		return nil
	}

	return &manifest.Config
}