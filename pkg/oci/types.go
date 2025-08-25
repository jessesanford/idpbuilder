// Package oci provides minimal OCI types needed by the stack package
// This is a dependency stub for Split-002 that references types from Split-001
package oci

// OCIReference represents a parsed OCI image reference
// This is a minimal stub - full implementation is in Split-001
type OCIReference struct {
	// Registry is the hostname and optional port of the registry
	Registry string `json:"registry"`
	
	// Namespace is the namespace/organization (e.g., "library" for official images)
	Namespace string `json:"namespace"`
	
	// Repository is the repository name within the namespace
	Repository string `json:"repository"`
	
	// Tag is the image tag (e.g., "latest", "v1.0.0")
	Tag string `json:"tag,omitempty"`
	
	// Digest is the content digest for immutable references
	Digest string `json:"digest,omitempty"`
}

// OCIImage represents an OCI container image
// This is a minimal stub - full implementation is in Split-001
type OCIImage struct {
	// Reference is the parsed image reference
	Reference *OCIReference `json:"reference"`
	
	// MediaType is the manifest media type
	MediaType string `json:"mediaType"`
	
	// Size is the image size in bytes
	Size int64 `json:"size,omitempty"`
	
	// Digest is the image manifest digest
	Digest string `json:"digest,omitempty"`
	
	// Annotations contains arbitrary metadata
	Annotations map[string]string `json:"annotations,omitempty"`
}

// NOTE: This is a minimal dependency stub for Split-002
// The complete OCI types implementation is in Split-001
// In a real deployment, Split-001 would be merged first, making these types available