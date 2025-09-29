// Package format defines OCI package format structures and specifications
// following OCI Image Specification v1.0.2.
package format

// Constants for OCI specification compliance
const (
	SchemaVersion = 2 // OCI manifest schema version
)

// MediaTypes defines OCI media type constants following OCI spec v1.0.2
var MediaTypes = struct {
	ManifestV2   string
	ConfigV1     string
	LayerTarGzip string
}{
	ManifestV2:   "application/vnd.oci.image.manifest.v1+json",
	ConfigV1:     "application/vnd.oci.image.config.v1+json",
	LayerTarGzip: "application/vnd.oci.image.layer.v1.tar+gzip",
}

// ArtifactType represents the type of OCI artifact
type ArtifactType string

// Descriptor represents OCI descriptor for config and layers
// Contains media type, digest, and size information
type Descriptor struct {
	MediaType string `json:"mediaType"`
	Digest    string `json:"digest"`
	Size      int64  `json:"size"`
}

// LayerDescriptor extends Descriptor with additional metadata for layers
// Embeds Descriptor and adds optional annotations
type LayerDescriptor struct {
	Descriptor
	Annotations map[string]string `json:"annotations,omitempty"`
}

// PackageManifest represents the core OCI package manifest structure
// Follows OCI Image Specification v1.0.2 manifest format
type PackageManifest struct {
	SchemaVersion int               `json:"schemaVersion"`
	MediaType     string            `json:"mediaType"`
	Config        Descriptor        `json:"config"`
	Layers        []LayerDescriptor `json:"layers"`
	Annotations   map[string]string `json:"annotations,omitempty"`
}

// PackageMetadata contains metadata information for OCI packages
// Provides name, version, description and label information
type PackageMetadata struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Description string            `json:"description,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
}