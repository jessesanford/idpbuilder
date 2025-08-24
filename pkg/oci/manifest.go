package oci

import (
	"encoding/json"
	"fmt"
	"time"
)

// OCIManifest represents an OCI Image Manifest following OCI spec.
type OCIManifest struct {
	SchemaVersion int                   `json:"schemaVersion" yaml:"schema_version"`
	MediaType     string                `json:"mediaType" yaml:"media_type"`
	Config        OCIDescriptor         `json:"config" yaml:"config"`
	Layers        []OCIDescriptor       `json:"layers" yaml:"layers"`
	Annotations   map[string]string     `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Subject       *OCIDescriptor        `json:"subject,omitempty" yaml:"subject,omitempty"`
}

// OCIDescriptor represents a content descriptor from OCI spec.
type OCIDescriptor struct {
	MediaType   string            `json:"mediaType" yaml:"media_type"`
	Digest      string            `json:"digest" yaml:"digest"`
	Size        int64             `json:"size" yaml:"size"`
	URLs        []string          `json:"urls,omitempty" yaml:"urls,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Data        []byte            `json:"data,omitempty" yaml:"data,omitempty"`
	Platform    *OCIPlatform      `json:"platform,omitempty" yaml:"platform,omitempty"`
}

// OCIManifestList represents an OCI Image Index for multi-platform images.
type OCIManifestList struct {
	SchemaVersion int                   `json:"schemaVersion" yaml:"schema_version"`
	MediaType     string                `json:"mediaType" yaml:"media_type"`
	Manifests     []OCIDescriptor       `json:"manifests" yaml:"manifests"`
	Annotations   map[string]string     `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Subject       *OCIDescriptor        `json:"subject,omitempty" yaml:"subject,omitempty"`
}

// OCIConfig represents image configuration from OCI Image Config Spec.
type OCIConfig struct {
	Created      *time.Time      `json:"created,omitempty" yaml:"created,omitempty"`
	Author       string          `json:"author,omitempty" yaml:"author,omitempty"`
	Architecture string          `json:"architecture" yaml:"architecture"`
	OS           string          `json:"os" yaml:"os"`
	OSVersion    string          `json:"os.version,omitempty" yaml:"os_version,omitempty"`
	OSFeatures   []string        `json:"os.features,omitempty" yaml:"os_features,omitempty"`
	Variant      string          `json:"variant,omitempty" yaml:"variant,omitempty"`
	Config       OCIImageConfig  `json:"config,omitempty" yaml:"config,omitempty"`
	RootFS       OCIRootFS       `json:"rootfs" yaml:"rootfs"`
	History      []OCIHistory    `json:"history,omitempty" yaml:"history,omitempty"`
}

// OCIImageConfig contains execution configuration for the container.
type OCIImageConfig struct {
	User         string              `json:"User,omitempty" yaml:"user,omitempty"`
	ExposedPorts map[string]struct{} `json:"ExposedPorts,omitempty" yaml:"exposed_ports,omitempty"`
	Env          []string            `json:"Env,omitempty" yaml:"env,omitempty"`
	Entrypoint   []string            `json:"Entrypoint,omitempty" yaml:"entrypoint,omitempty"`
	Cmd          []string            `json:"Cmd,omitempty" yaml:"cmd,omitempty"`
	Volumes      map[string]struct{} `json:"Volumes,omitempty" yaml:"volumes,omitempty"`
	WorkingDir   string              `json:"WorkingDir,omitempty" yaml:"working_dir,omitempty"`
	Labels       map[string]string   `json:"Labels,omitempty" yaml:"labels,omitempty"`
	StopSignal   string              `json:"StopSignal,omitempty" yaml:"stop_signal,omitempty"`
}

// OCIRootFS describes the root filesystem for the image.
type OCIRootFS struct {
	Type    string   `json:"type" yaml:"type"`
	DiffIDs []string `json:"diff_ids" yaml:"diff_ids"`
}

// OCIHistory represents one layer in the image's history.
type OCIHistory struct {
	Created    *time.Time `json:"created,omitempty" yaml:"created,omitempty"`
	CreatedBy  string     `json:"created_by,omitempty" yaml:"created_by,omitempty"`
	Author     string     `json:"author,omitempty" yaml:"author,omitempty"`
	Comment    string     `json:"comment,omitempty" yaml:"comment,omitempty"`
	EmptyLayer bool       `json:"empty_layer,omitempty" yaml:"empty_layer,omitempty"`
}

// IsValid validates manifest structure.
func (m OCIManifest) IsValid() error {
	if m.SchemaVersion != 2 {
		return fmt.Errorf("invalid schema version: expected 2, got %d", m.SchemaVersion)
	}
	if m.MediaType == "" {
		return fmt.Errorf("manifest media type is required")
	}
	if err := m.Config.IsValid(); err != nil {
		return fmt.Errorf("invalid config descriptor: %w", err)
	}
	if len(m.Layers) == 0 {
		return fmt.Errorf("manifest must contain at least one layer")
	}
	for i, layer := range m.Layers {
		if err := layer.IsValid(); err != nil {
			return fmt.Errorf("invalid layer %d: %w", i, err)
		}
	}
	return nil
}

// IsValid validates descriptor fields.
func (d OCIDescriptor) IsValid() error {
	if d.MediaType == "" {
		return fmt.Errorf("descriptor media type is required")
	}
	if d.Digest == "" {
		return fmt.Errorf("descriptor digest is required")
	}
	if d.Size < 0 {
		return fmt.Errorf("descriptor size cannot be negative")
	}
	return nil
}

// ToJSON serializes the manifest to JSON.
func (m OCIManifest) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

// FromJSON deserializes a manifest from JSON.
func (m *OCIManifest) FromJSON(data []byte) error {
	return json.Unmarshal(data, m)
}