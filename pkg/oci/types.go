// Package oci provides types and interfaces for OCI (Open Container Initiative) images and repositories.
package oci

import (
	"strings"
	"time"
)

// OCIImage represents an OCI container image with metadata.
type OCIImage struct {
	Name        string            `json:"name" yaml:"name"`
	Tag         string            `json:"tag" yaml:"tag"`
	Digest      string            `json:"digest,omitempty" yaml:"digest,omitempty"`
	Registry    string            `json:"registry" yaml:"registry"`
	Repository  string            `json:"repository" yaml:"repository"`
	Platform    OCIPlatform       `json:"platform" yaml:"platform"`
	Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	CreatedAt   time.Time         `json:"created_at" yaml:"created_at"`
	Size        int64             `json:"size" yaml:"size"`
}

// OCIPlatform describes the target platform for an OCI image.
type OCIPlatform struct {
	OS           string   `json:"os" yaml:"os"`
	Architecture string   `json:"architecture" yaml:"architecture"`
	Variant      string   `json:"variant,omitempty" yaml:"variant,omitempty"`
	OSVersion    string   `json:"os.version,omitempty" yaml:"os_version,omitempty"`
	OSFeatures   []string `json:"os.features,omitempty" yaml:"os_features,omitempty"`
}

// String returns platform as "os/arch" format.
func (p OCIPlatform) String() string {
	result := p.OS + "/" + p.Architecture
	if p.Variant != "" {
		result += "/" + p.Variant
	}
	return result
}

// OCIReference represents a complete reference to an OCI image.
type OCIReference struct {
	Registry   string `json:"registry" yaml:"registry"`
	Repository string `json:"repository" yaml:"repository"`
	Tag        string `json:"tag,omitempty" yaml:"tag,omitempty"`
	Digest     string `json:"digest,omitempty" yaml:"digest,omitempty"`
}

// String returns the canonical string representation: [registry/]repository[:tag][@digest]
func (r OCIReference) String() string {
	var parts []string
	if r.Registry != "" && r.Registry != DefaultRegistry {
		parts = append(parts, r.Registry)
	}
	parts = append(parts, r.Repository)
	reference := strings.Join(parts, "/")
	if r.Tag != "" {
		reference += ":" + r.Tag
	}
	if r.Digest != "" {
		reference += "@" + r.Digest
	}
	return reference
}

// IsValid checks if the OCI reference has required fields.
func (r OCIReference) IsValid() bool {
	return r.Repository != "" && (r.Tag != "" || r.Digest != "")
}

// OCILayer represents a single layer within an OCI image.
type OCILayer struct {
	Digest      string            `json:"digest" yaml:"digest"`
	Size        int64             `json:"size" yaml:"size"`
	MediaType   string            `json:"mediaType" yaml:"media_type"`
	URLs        []string          `json:"urls,omitempty" yaml:"urls,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	CreatedBy   string            `json:"created_by,omitempty" yaml:"created_by,omitempty"`
}

// OCIRepository defines the interface for interacting with OCI image repositories.
type OCIRepository interface {
	List(filter OCIListFilter) ([]OCIImage, error)
	Get(ref OCIReference) (*OCIImage, error)
	Push(image *OCIImage, options OCIPushOptions) error
	Pull(ref OCIReference, options OCIPullOptions) (*OCIImage, error)
	Delete(ref OCIReference) error
	Tags(repository string) ([]string, error)
}

// OCIListFilter provides filtering options when listing images.
type OCIListFilter struct {
	Repository       string       `json:"repository,omitempty"`
	Tag              string       `json:"tag,omitempty"`
	Platform         *OCIPlatform `json:"platform,omitempty"`
	MaxResults       int          `json:"max_results,omitempty"`
	IncludeManifests bool         `json:"include_manifests,omitempty"`
}

// OCIPushOptions configures options for pushing images.
type OCIPushOptions struct {
	Tags        []string          `json:"tags,omitempty"`
	Platform    *OCIPlatform      `json:"platform,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Force       bool              `json:"force,omitempty"`
}

// OCIPullOptions configures options for pulling images.
type OCIPullOptions struct {
	Platform         *OCIPlatform                  `json:"platform,omitempty"`
	ProgressCallback func(progress OCIPullProgress) `json:"-"`
	VerifySignature  bool                          `json:"verify_signature,omitempty"`
	MaxConcurrency   int                           `json:"max_concurrency,omitempty"`
}

// OCIPullProgress represents the progress of a pull operation.
type OCIPullProgress struct {
	Layer           OCILayer `json:"layer"`
	BytesDownloaded int64    `json:"bytes_downloaded"`
	TotalBytes      int64    `json:"total_bytes"`
	CompletedLayers int      `json:"completed_layers"`
	TotalLayers     int      `json:"total_layers"`
}