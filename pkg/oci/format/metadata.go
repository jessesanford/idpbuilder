package format

import (
	"encoding/json"
	"errors"
)

// SerializeManifest converts a package manifest to JSON bytes
// Returns serialized manifest data or an error if serialization fails
func SerializeManifest(manifest *PackageManifest) ([]byte, error) {
	if manifest == nil {
		return nil, errors.New("manifest cannot be nil")
	}

	data, err := json.Marshal(manifest)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// ParseManifest parses JSON data into a package manifest structure
// Returns a PackageManifest or an error if parsing fails
func ParseManifest(data []byte) (*PackageManifest, error) {
	if len(data) == 0 {
		return nil, errors.New("manifest data cannot be empty")
	}

	var manifest PackageManifest
	err := json.Unmarshal(data, &manifest)
	if err != nil {
		return nil, err
	}

	return &manifest, nil
}

// NewLayerDescriptor creates a new layer descriptor with the specified properties
// Returns a LayerDescriptor with digest, media type, and size properly set
func NewLayerDescriptor(digest, mediaType string, size int64) *LayerDescriptor {
	if digest == "" || mediaType == "" || size <= 0 {
		return nil
	}

	return &LayerDescriptor{
		Descriptor: Descriptor{
			MediaType: mediaType,
			Digest:    digest,
			Size:      size,
		},
		Annotations: make(map[string]string),
	}
}

// GetMediaType maps an artifact type to its corresponding OCI media type
// Returns the appropriate media type string for the given artifact type
func GetMediaType(artifactType ArtifactType) string {
	switch artifactType {
	case "manifest":
		return MediaTypes.ManifestV2
	case "config":
		return MediaTypes.ConfigV1
	case "layer":
		return MediaTypes.LayerTarGzip
	default:
		return MediaTypes.ManifestV2
	}
}