package providers

import (
	"time"
)

// Artifact represents an OCI artifact with its manifest, layers, and metadata.
// This structure contains all the components needed for pushing or pulling
// artifacts from OCI registries.
type Artifact struct {
	// MediaType specifies the OCI media type of the artifact
	MediaType string

	// Manifest contains the raw manifest content as bytes
	Manifest []byte

	// Layers contains all the layers that make up this artifact
	Layers []Layer

	// Config contains the configuration blob for the artifact
	Config []byte

	// Annotations contains OCI annotations for additional metadata
	Annotations map[string]string
}

// Layer represents a single layer within an OCI artifact.
// Each layer has content, metadata, and integrity information.
type Layer struct {
	// MediaType specifies the media type of this layer
	MediaType string

	// Digest is the content digest for integrity verification
	Digest string

	// Size is the layer size in bytes
	Size int64

	// Data contains the actual layer content
	Data []byte
}

// ArtifactInfo provides metadata about an artifact in a registry.
// This is used for listing operations and contains summary information
// without the actual artifact content.
type ArtifactInfo struct {
	// Reference is the full reference (registry/repo:tag)
	Reference string

	// Digest is the manifest digest for the artifact
	Digest string

	// Tags contains all tags associated with this artifact
	Tags []string

	// Size is the total size of the artifact in bytes
	Size int64

	// Created is the creation timestamp of the artifact
	Created time.Time

	// Annotations contains metadata annotations for the artifact
	Annotations map[string]string
}

// ProviderConfig contains configuration settings for registry providers.
// This includes authentication, connection settings, and operational parameters.
type ProviderConfig struct {
	// Registry is the base URL of the registry
	Registry string

	// Auth contains authentication configuration
	Auth AuthConfig

	// Insecure allows connections to registries with invalid certificates
	Insecure bool

	// Timeout specifies the timeout for registry operations
	Timeout time.Duration
}

// AuthConfig contains authentication credentials and settings.
// Different authentication methods can be used depending on the registry.
type AuthConfig struct {
	// Username for basic authentication
	Username string

	// Password for basic authentication
	Password string

	// Token for bearer token authentication
	Token string

	// RegistryToken for registry-specific token authentication
	RegistryToken string
}

// RegistryCapabilities describes the capabilities and limits of a registry.
// This information helps clients understand what operations are supported
// and any size or feature limitations.
type RegistryCapabilities struct {
	// SupportsDelete indicates if the registry supports artifact deletion
	SupportsDelete bool

	// SupportsOCI indicates full OCI Distribution Spec compliance
	SupportsOCI bool

	// MaxLayerSize is the maximum size allowed for individual layers
	MaxLayerSize int64
}