package oci

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// ManifestReader provides functionality to read and parse OCI manifests
type ManifestReader struct {
	// MaxSize limits the maximum manifest size that can be read
	MaxSize int64
}

// NewManifestReader creates a new manifest reader with default settings
func NewManifestReader() *ManifestReader {
	return &ManifestReader{
		MaxSize: MaxManifestSize,
	}
}

// ReadManifest reads and parses an OCI manifest from a reader
func (mr *ManifestReader) ReadManifest(r io.Reader) (*OCIManifest, error) {
	// Limit the reader to prevent excessive memory usage
	limitedReader := io.LimitReader(r, mr.MaxSize)
	
	data, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest: %w", err)
	}
	
	if len(data) == 0 {
		return nil, fmt.Errorf("empty manifest")
	}
	
	var manifest OCIManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest JSON: %w", err)
	}
	
	// Validate the parsed manifest
	if err := manifest.Validate(); err != nil {
		return nil, fmt.Errorf("invalid manifest: %w", err)
	}
	
	return &manifest, nil
}

// ParseManifestJSON parses an OCI manifest from JSON bytes
func ParseManifestJSON(data []byte) (*OCIManifest, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty manifest data")
	}
	
	if len(data) > MaxManifestSize {
		return nil, fmt.Errorf("manifest size %d exceeds maximum %d", len(data), MaxManifestSize)
	}
	
	var manifest OCIManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest JSON: %w", err)
	}
	
	if err := manifest.Validate(); err != nil {
		return nil, fmt.Errorf("invalid manifest: %w", err)
	}
	
	return &manifest, nil
}

// ManifestWriter provides functionality to write OCI manifests
type ManifestWriter struct{}

// NewManifestWriter creates a new manifest writer
func NewManifestWriter() *ManifestWriter {
	return &ManifestWriter{}
}

// WriteManifest writes an OCI manifest to a writer as JSON
func (mw *ManifestWriter) WriteManifest(w io.Writer, manifest *OCIManifest) error {
	if manifest == nil {
		return fmt.Errorf("manifest cannot be nil")
	}
	
	// Validate before writing
	if err := manifest.Validate(); err != nil {
		return fmt.Errorf("invalid manifest: %w", err)
	}
	
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	
	if err := encoder.Encode(manifest); err != nil {
		return fmt.Errorf("failed to encode manifest: %w", err)
	}
	
	return nil
}

// MarshalManifest marshals an OCI manifest to JSON bytes
func MarshalManifest(manifest *OCIManifest) ([]byte, error) {
	if manifest == nil {
		return nil, fmt.Errorf("manifest cannot be nil")
	}
	
	if err := manifest.Validate(); err != nil {
		return nil, fmt.Errorf("invalid manifest: %w", err)
	}
	
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal manifest: %w", err)
	}
	
	return data, nil
}

// DetectManifestMediaType detects the media type of a manifest from its content
func DetectManifestMediaType(data []byte) (string, error) {
	if len(data) == 0 {
		return "", fmt.Errorf("empty data")
	}
	
	// Try to parse as JSON to get the mediaType field
	var temp struct {
		MediaType string `json:"mediaType"`
	}
	
	if err := json.Unmarshal(data, &temp); err != nil {
		return "", fmt.Errorf("not valid JSON: %w", err)
	}
	
	if temp.MediaType == "" {
		return "", fmt.Errorf("mediaType field is missing or empty")
	}
	
	// Validate that it's a known manifest media type
	switch temp.MediaType {
	case MediaTypeManifest, MediaTypeImageIndex:
		return temp.MediaType, nil
	case DockerMediaTypeManifest, DockerMediaTypeManifestList:
		return temp.MediaType, nil
	default:
		return "", fmt.Errorf("unknown manifest media type: %s", temp.MediaType)
	}
}

// IsOCIManifest checks if the provided data represents an OCI manifest
func IsOCIManifest(data []byte) bool {
	mediaType, err := DetectManifestMediaType(data)
	if err != nil {
		return false
	}
	
	return mediaType == MediaTypeManifest || mediaType == MediaTypeImageIndex
}

// IsDockerManifest checks if the provided data represents a Docker manifest
func IsDockerManifest(data []byte) bool {
	mediaType, err := DetectManifestMediaType(data)
	if err != nil {
		return false
	}
	
	return mediaType == DockerMediaTypeManifest || mediaType == DockerMediaTypeManifestList
}

// NormalizeManifestMediaType converts Docker manifest media types to OCI equivalents
func NormalizeManifestMediaType(mediaType string) string {
	switch mediaType {
	case DockerMediaTypeManifest:
		return MediaTypeManifest
	case DockerMediaTypeManifestList:
		return MediaTypeImageIndex
	case DockerMediaTypeConfig:
		return MediaTypeImageConfig
	case DockerMediaTypeLayer:
		return MediaTypeImageLayer
	default:
		return mediaType
	}
}

// ValidateDigest validates that a digest string follows the expected format
func ValidateDigest(digest string) error {
	if digest == "" {
		return fmt.Errorf("digest cannot be empty")
	}
	
	parts := strings.SplitN(digest, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("digest must contain algorithm and hex: %s", digest)
	}
	
	algorithm, hex := parts[0], parts[1]
	
	if algorithm == "" {
		return fmt.Errorf("digest algorithm cannot be empty")
	}
	
	if hex == "" {
		return fmt.Errorf("digest hex cannot be empty")
	}
	
	// Validate common algorithms
	switch algorithm {
	case "sha256":
		if len(hex) != 64 {
			return fmt.Errorf("sha256 digest must be 64 hex characters, got %d", len(hex))
		}
	case "sha512":
		if len(hex) != 128 {
			return fmt.Errorf("sha512 digest must be 128 hex characters, got %d", len(hex))
		}
	default:
		// Allow other algorithms but don't validate hex length
	}
	
	// Check that hex contains only valid characters
	for _, char := range hex {
		if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f') || (char >= 'A' && char <= 'F')) {
			return fmt.Errorf("digest hex contains invalid character: %c", char)
		}
	}
	
	return nil
}

// CreateManifest creates a new OCI manifest with the provided configuration and layers
func CreateManifest(config OCIDescriptor, layers []OCIDescriptor, annotations map[string]string) *OCIManifest {
	manifest := &OCIManifest{
		SchemaVersion: 2,
		MediaType:     MediaTypeManifest,
		Config:        config,
		Layers:        layers,
		Annotations:   annotations,
	}
	
	if manifest.Annotations == nil {
		manifest.Annotations = make(map[string]string)
	}
	
	return manifest
}