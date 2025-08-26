package api

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// RegistryClient handles OCI registry operations including push, pull, and metadata
// management with support for authentication and advanced registry features.
type RegistryClient interface {
	// Push uploads an image to the registry with the specified authentication.
	Push(ctx context.Context, image string, auth AuthConfig) error

	// Pull downloads an image from the registry to the local system.
	Pull(ctx context.Context, image string, auth AuthConfig) (*Image, error)

	// GetManifest retrieves the manifest for a specific image.
	GetManifest(ctx context.Context, image string) (*Manifest, error)

	// ListTags returns all tags for a given repository.
	ListTags(ctx context.Context, repository string) ([]string, error)

	// Delete removes an image or tag from the registry.
	Delete(ctx context.Context, image string, auth AuthConfig) error

	// Ping checks registry connectivity and health.
	Ping(ctx context.Context) error

	// GetRegistryInfo returns detailed information about the registry.
	GetRegistryInfo(ctx context.Context) (*RegistryInfo, error)

	// ListRepositories returns all repositories accessible to the authenticated user.
	ListRepositories(ctx context.Context, auth AuthConfig) ([]string, error)

	// CopyImage copies an image between registries or repositories.
	CopyImage(ctx context.Context, source, destination string, auth AuthConfig) error

	// GetImageHistory retrieves the build history for an image.
	GetImageHistory(ctx context.Context, image string) ([]*HistoryEntry, error)
}

// AuthConfig contains authentication credentials for registry operations.
type AuthConfig struct {
	Username       string `json:"username,omitempty" yaml:"username,omitempty"`
	Password       string `json:"password,omitempty" yaml:"password,omitempty"`
	Auth           string `json:"auth,omitempty" yaml:"auth,omitempty"`
	ServerAddress  string `json:"serveraddress,omitempty" yaml:"serveraddress,omitempty"`
	IdentityToken  string `json:"identitytoken,omitempty" yaml:"identitytoken,omitempty"`
	RegistryToken  string `json:"registrytoken,omitempty" yaml:"registrytoken,omitempty"`
	Email          string `json:"email,omitempty" yaml:"email,omitempty"`
	Insecure       bool   `json:"insecure,omitempty" yaml:"insecure,omitempty"`
	CACertPath     string `json:"ca_cert_path,omitempty" yaml:"ca_cert_path,omitempty"`
	ClientCertPath string `json:"client_cert_path,omitempty" yaml:"client_cert_path,omitempty"`
	ClientKeyPath  string `json:"client_key_path,omitempty" yaml:"client_key_path,omitempty"`
}

// Validate checks if the authentication configuration has required fields.
func (a AuthConfig) Validate() error {
	if a.ServerAddress == "" {
		return fmt.Errorf("server address is required")
	}

	hasBasicAuth := a.Username != "" && a.Password != ""
	hasAuth := a.Auth != ""
	hasToken := a.IdentityToken != "" || a.RegistryToken != ""

	if !hasBasicAuth && !hasAuth && !hasToken {
		return fmt.Errorf("at least one authentication method must be provided")
	}

	if !strings.Contains(a.ServerAddress, "://") {
		return fmt.Errorf("server address must include protocol (http:// or https://)")
	}

	return nil
}

// IsEmpty returns true if the auth config has no authentication credentials.
func (a AuthConfig) IsEmpty() bool {
	return a.Username == "" && a.Password == "" && a.Auth == "" && 
		   a.IdentityToken == "" && a.RegistryToken == ""
}

// Image represents an OCI image with its metadata and layer information.
type Image struct {
	Name        string            `json:"name" yaml:"name"`
	Tag         string            `json:"tag" yaml:"tag"`
	Digest      string            `json:"digest" yaml:"digest"`
	MediaType   string            `json:"media_type" yaml:"media_type"`
	Size        int64             `json:"size" yaml:"size"`
	Created     time.Time         `json:"created" yaml:"created"`
	Config      *ImageConfig      `json:"config,omitempty" yaml:"config,omitempty"`
	Layers      []*LayerInfo      `json:"layers" yaml:"layers"`
	Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Platform    *Platform         `json:"platform,omitempty" yaml:"platform,omitempty"`
}

// ImageConfig contains the configuration data for an image.
type ImageConfig struct {
	Architecture string           `json:"architecture" yaml:"architecture"`
	OS           string           `json:"os" yaml:"os"`
	Config       *ContainerConfig `json:"config,omitempty" yaml:"config,omitempty"`
	RootFS       *RootFS          `json:"rootfs,omitempty" yaml:"rootfs,omitempty"`
	History      []*HistoryEntry  `json:"history,omitempty" yaml:"history,omitempty"`
}

// ContainerConfig contains runtime configuration for containers.
type ContainerConfig struct {
	User         string                 `json:"user,omitempty" yaml:"user,omitempty"`
	Env          []string               `json:"env,omitempty" yaml:"env,omitempty"`
	Cmd          []string               `json:"cmd,omitempty" yaml:"cmd,omitempty"`
	Entrypoint   []string               `json:"entrypoint,omitempty" yaml:"entrypoint,omitempty"`
	WorkingDir   string                 `json:"working_dir,omitempty" yaml:"working_dir,omitempty"`
	ExposedPorts map[string]struct{}    `json:"exposed_ports,omitempty" yaml:"exposed_ports,omitempty"`
	Volumes      map[string]struct{}    `json:"volumes,omitempty" yaml:"volumes,omitempty"`
	Labels       map[string]string      `json:"labels,omitempty" yaml:"labels,omitempty"`
}

// RootFS describes the root filesystem of an image.
type RootFS struct {
	Type    string   `json:"type" yaml:"type"`
	DiffIDs []string `json:"diff_ids" yaml:"diff_ids"`
}

// LayerInfo contains metadata about an image layer.
type LayerInfo struct {
	Digest      string            `json:"digest" yaml:"digest"`
	Size        int64             `json:"size" yaml:"size"`
	MediaType   string            `json:"media_type" yaml:"media_type"`
	URLs        []string          `json:"urls,omitempty" yaml:"urls,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
}

// Platform specifies the target platform for multi-architecture images.
type Platform struct {
	Architecture string   `json:"architecture" yaml:"architecture"`
	OS           string   `json:"os" yaml:"os"`
	OSVersion    string   `json:"os.version,omitempty" yaml:"os.version,omitempty"`
	OSFeatures   []string `json:"os.features,omitempty" yaml:"os.features,omitempty"`
	Variant      string   `json:"variant,omitempty" yaml:"variant,omitempty"`
}

// Manifest represents an OCI image manifest.
type Manifest struct {
	SchemaVersion int                   `json:"schemaVersion" yaml:"schemaVersion"`
	MediaType     string                `json:"mediaType" yaml:"mediaType"`
	Config        *Descriptor           `json:"config" yaml:"config"`
	Layers        []*Descriptor         `json:"layers" yaml:"layers"`
	Annotations   map[string]string     `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Subject       *Descriptor           `json:"subject,omitempty" yaml:"subject,omitempty"`
}

// Descriptor represents a content descriptor in OCI manifests.
type Descriptor struct {
	MediaType   string            `json:"mediaType" yaml:"mediaType"`
	Digest      string            `json:"digest" yaml:"digest"`
	Size        int64             `json:"size" yaml:"size"`
	URLs        []string          `json:"urls,omitempty" yaml:"urls,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Platform    *Platform         `json:"platform,omitempty" yaml:"platform,omitempty"`
}

// HistoryEntry represents an entry in the image build history.
type HistoryEntry struct {
	Created    time.Time `json:"created" yaml:"created"`
	CreatedBy  string    `json:"created_by,omitempty" yaml:"created_by,omitempty"`
	Author     string    `json:"author,omitempty" yaml:"author,omitempty"`
	Comment    string    `json:"comment,omitempty" yaml:"comment,omitempty"`
	EmptyLayer bool      `json:"empty_layer,omitempty" yaml:"empty_layer,omitempty"`
}

// RegistryInfo contains information about a registry's capabilities.
type RegistryInfo struct {
	Name            string   `json:"name" yaml:"name"`
	Version         string   `json:"version" yaml:"version"`
	URL             string   `json:"url" yaml:"url"`
	Capabilities    []string `json:"capabilities,omitempty" yaml:"capabilities,omitempty"`
	AuthMethods     []string `json:"auth_methods,omitempty" yaml:"auth_methods,omitempty"`
	MaxManifestSize int64    `json:"max_manifest_size,omitempty" yaml:"max_manifest_size,omitempty"`
	MaxLayerSize    int64    `json:"max_layer_size,omitempty" yaml:"max_layer_size,omitempty"`
	SupportsDelete  bool     `json:"supports_delete" yaml:"supports_delete"`
	SupportsCatalog bool     `json:"supports_catalog" yaml:"supports_catalog"`
}