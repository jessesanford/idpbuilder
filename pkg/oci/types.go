package oci

import (
	"fmt"
	"strings"
)

// OCIReference represents a parsed OCI image reference
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

// String returns the full image reference as a string
func (r *OCIReference) String() string {
	var parts []string
	
	// Add registry if not default
	if r.Registry != "" && r.Registry != DefaultRegistry {
		parts = append(parts, r.Registry)
	}
	
	// Add namespace/repository
	if r.Namespace != "" && r.Namespace != DefaultNamespace {
		parts = append(parts, r.Namespace+"/"+r.Repository)
	} else {
		parts = append(parts, r.Repository)
	}
	
	ref := strings.Join(parts, "/")
	
	// Add tag if present
	if r.Tag != "" {
		ref += ":" + r.Tag
	}
	
	// Add digest if present (takes precedence over tag)
	if r.Digest != "" {
		if r.Tag != "" {
			ref = strings.TrimSuffix(ref, ":"+r.Tag)
		}
		ref += "@" + r.Digest
	}
	
	return ref
}

// IsDigest returns true if the reference uses a digest instead of a tag
func (r *OCIReference) IsDigest() bool {
	return r.Digest != ""
}

// OCIImage represents an OCI container image
type OCIImage struct {
	// Reference is the parsed image reference
	Reference *OCIReference `json:"reference"`
	
	// MediaType is the manifest media type
	MediaType string `json:"mediaType"`
	
	// Size is the image size in bytes
	Size int64 `json:"size,omitempty"`
	
	// Digest is the image manifest digest
	Digest string `json:"digest,omitempty"`
	
	// Platform describes the target platform
	Platform *OCIPlatform `json:"platform,omitempty"`
	
	// Annotations contains arbitrary metadata
	Annotations map[string]string `json:"annotations,omitempty"`
}

// OCIPlatform represents the target platform for an image
type OCIPlatform struct {
	// Architecture is the CPU architecture (e.g., "amd64", "arm64")
	Architecture string `json:"architecture"`
	
	// OS is the operating system (e.g., "linux", "windows")
	OS string `json:"os"`
	
	// OSVersion is the OS version (optional)
	OSVersion string `json:"os.version,omitempty"`
	
	// OSFeatures are required OS features (optional)
	OSFeatures []string `json:"os.features,omitempty"`
	
	// Variant is the CPU variant (optional, e.g., "v7" for ARM)
	Variant string `json:"variant,omitempty"`
}

// String returns the platform as a string in the format "os/architecture[/variant]"
func (p *OCIPlatform) String() string {
	platform := p.OS + "/" + p.Architecture
	if p.Variant != "" {
		platform += "/" + p.Variant
	}
	return platform
}

// OCIDescriptor represents an OCI content descriptor
type OCIDescriptor struct {
	// MediaType is the media type of the referenced content
	MediaType string `json:"mediaType"`
	
	// Digest is the content digest
	Digest string `json:"digest"`
	
	// Size is the size in bytes of the referenced content
	Size int64 `json:"size"`
	
	// URLs are the URLs where the content can be retrieved
	URLs []string `json:"urls,omitempty"`
	
	// Annotations contains arbitrary metadata
	Annotations map[string]string `json:"annotations,omitempty"`
	
	// Platform describes the target platform (for manifest lists)
	Platform *OCIPlatform `json:"platform,omitempty"`
}

// OCIManifest represents an OCI image manifest
type OCIManifest struct {
	// SchemaVersion is the schema version (must be 2)
	SchemaVersion int `json:"schemaVersion"`
	
	// MediaType is the manifest media type
	MediaType string `json:"mediaType"`
	
	// Config is the image configuration descriptor
	Config OCIDescriptor `json:"config"`
	
	// Layers are the layer descriptors
	Layers []OCIDescriptor `json:"layers"`
	
	// Annotations contains arbitrary metadata
	Annotations map[string]string `json:"annotations,omitempty"`
}

// Validate validates the manifest structure and content
func (m *OCIManifest) Validate() error {
	if m.SchemaVersion != 2 {
		return fmt.Errorf("invalid schema version: %d, expected 2", m.SchemaVersion)
	}
	
	if m.MediaType == "" {
		return fmt.Errorf("mediaType is required")
	}
	
	if m.Config.MediaType == "" {
		return fmt.Errorf("config mediaType is required")
	}
	
	if m.Config.Digest == "" {
		return fmt.Errorf("config digest is required")
	}
	
	if m.Config.Size <= 0 {
		return fmt.Errorf("config size must be positive")
	}
	
	if len(m.Layers) == 0 {
		return fmt.Errorf("manifest must have at least one layer")
	}
	
	for i, layer := range m.Layers {
		if layer.MediaType == "" {
			return fmt.Errorf("layer %d mediaType is required", i)
		}
		if layer.Digest == "" {
			return fmt.Errorf("layer %d digest is required", i)
		}
		if layer.Size <= 0 {
			return fmt.Errorf("layer %d size must be positive", i)
		}
	}
	
	return nil
}

// ParseOCIReference parses an OCI image reference string
func ParseOCIReference(ref string) (*OCIReference, error) {
	if ref == "" {
		return nil, fmt.Errorf("empty reference")
	}
	
	// Split by @ first to separate digest
	var digestPart string
	if strings.Contains(ref, "@") {
		parts := strings.SplitN(ref, "@", 2)
		ref = parts[0]
		digestPart = parts[1]
	}
	
	// Split by : to separate tag  
	var tagPart string
	if strings.Contains(ref, ":") && !strings.Contains(ref, "://") {
		// Find the last colon (could be port or tag)
		lastColon := strings.LastIndex(ref, ":")
		beforeColon := ref[:lastColon]
		afterColon := ref[lastColon+1:]
		
		// Check if this looks like a port (registry:port) or tag
		// If there's no slash before the last colon, it could be registry:port
		// If there's a slash, it's more likely namespace/repo:tag
		lastSlash := strings.LastIndex(beforeColon, "/")
		
		if lastSlash == -1 {
			// No slash before colon, could be registry:port or just repo:tag
			// Check if afterColon looks like a port number (all digits) vs tag
			isPort := true
			for _, char := range afterColon {
				if char < '0' || char > '9' {
					isPort = false
					break
				}
			}
			if !isPort {
				// It's a tag
				tagPart = afterColon
				ref = beforeColon
			}
		} else {
			// There's a slash before the colon, so this is likely repo:tag
			tagPart = afterColon
			ref = beforeColon
		}
	}
	
	result := &OCIReference{}
	
	// Split the remaining ref by /
	parts := strings.Split(ref, "/")
	
	switch len(parts) {
	case 1:
		// Just repository name: nginx
		result.Registry = DefaultRegistry
		result.Namespace = DefaultNamespace
		result.Repository = parts[0]
	case 2:
		// namespace/repository: myorg/myapp OR registry/repository
		if strings.Contains(parts[0], ".") || strings.Contains(parts[0], ":") {
			// registry/repository: docker.io/nginx
			result.Registry = parts[0]
			result.Namespace = DefaultNamespace
			result.Repository = parts[1]
		} else {
			// namespace/repository: myorg/myapp
			result.Registry = DefaultRegistry
			result.Namespace = parts[0]
			result.Repository = parts[1]
		}
	case 3:
		// registry/namespace/repository: registry.example.com/myorg/myapp
		result.Registry = parts[0]
		result.Namespace = parts[1]
		result.Repository = parts[2]
	default:
		// More than 3 parts, assume last is repository, second-to-last is namespace
		if len(parts) >= 3 {
			result.Registry = strings.Join(parts[:len(parts)-2], "/")
			result.Namespace = parts[len(parts)-2]
			result.Repository = parts[len(parts)-1]
		} else {
			return nil, fmt.Errorf("invalid reference format: %s", ref)
		}
	}
	
	result.Tag = tagPart
	result.Digest = digestPart
	
	// Set default tag if no tag or digest specified
	if result.Tag == "" && result.Digest == "" {
		result.Tag = DefaultTag
	}
	
	return result, nil
}

// NewOCIImage creates a new OCIImage from a reference string
func NewOCIImage(ref string) (*OCIImage, error) {
	reference, err := ParseOCIReference(ref)
	if err != nil {
		return nil, fmt.Errorf("failed to parse reference: %w", err)
	}
	
	return &OCIImage{
		Reference:   reference,
		MediaType:   MediaTypeManifest,
		Annotations: make(map[string]string),
	}, nil
}