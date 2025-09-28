package mapper

// Layer represents an OCI layer (simplified Wave 1 equivalent)
type Layer struct {
	MediaType string
	Size      int64
	Digest    string
}

// PackageManifest represents an OCI manifest (simplified Wave 1 equivalent)
type PackageManifest struct {
	SchemaVersion int                 `json:"schemaVersion"`
	MediaType     string              `json:"mediaType"`
	Config        Descriptor          `json:"config"`
	Layers        []LayerDescriptor   `json:"layers"`
	Annotations   map[string]string   `json:"annotations,omitempty"`
}

// Descriptor represents a content descriptor
type Descriptor struct {
	MediaType string `json:"mediaType"`
	Size      int64  `json:"size"`
	Digest    string `json:"digest"`
}

// LayerDescriptor represents a layer descriptor
type LayerDescriptor struct {
	MediaType string `json:"mediaType"`
	Size      int64  `json:"size"`
	Digest    string `json:"digest"`
}

// StackConfig represents an IDPBuilder stack configuration
type StackConfig struct {
	Name       string            `json:"name"`
	Version    string            `json:"version"`
	Components []Component       `json:"components"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

// Component represents a single component in the stack
type Component struct {
	Name   string         `json:"name"`
	Type   string         `json:"type"`
	Source string         `json:"source"`
	Config map[string]any `json:"config,omitempty"`
}

// ContainerSpec defines the specification for building a container
type ContainerSpec struct {
	Name      string
	BaseImage string
	Layers    []Layer
	Env       map[string]string
	Labels    map[string]string
}

// MappingResult contains the result of stack-to-container mapping
type MappingResult struct {
	Containers []ContainerSpec
	Manifest   *PackageManifest
	Metadata   map[string]string
}

// ComponentRef represents a resolved component reference
type ComponentRef struct {
	Registry   string
	Repository string
	Tag        string
	Digest     string
}