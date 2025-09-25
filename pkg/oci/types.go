// Package oci provides OCI registry client functionality for idpbuilder-push.
// This implementation focuses on core types and authentication needed for
// pushing container images to OCI-compliant registries.
package oci

import (
	"context"
	"time"
)

// RegistryAuth represents authentication credentials for OCI registry access.
// It supports multiple authentication methods including basic auth and tokens.
type RegistryAuth struct {
	// Username for basic authentication
	Username string `json:"username,omitempty"`

	// Password for basic authentication
	Password string `json:"password,omitempty"`

	// Token for bearer token authentication
	Token string `json:"token,omitempty"`

	// ServerAddress is the registry server address
	ServerAddress string `json:"server_address,omitempty"`

	// Realm is the authentication realm
	Realm string `json:"realm,omitempty"`

	// Service is the authentication service
	Service string `json:"service,omitempty"`
}

// IsEmpty returns true if no authentication credentials are provided.
func (a *RegistryAuth) IsEmpty() bool {
	return a == nil || (a.Username == "" && a.Password == "" && a.Token == "")
}

// HasBasicAuth returns true if username and password are provided.
func (a *RegistryAuth) HasBasicAuth() bool {
	return a != nil && a.Username != "" && a.Password != ""
}

// HasTokenAuth returns true if a bearer token is provided.
func (a *RegistryAuth) HasTokenAuth() bool {
	return a != nil && a.Token != ""
}

// PushOptions represents configuration options for pushing images to a registry.
// These options control authentication, behavior, and validation settings.
type PushOptions struct {
	// ImageRef is the full image reference (registry/repo:tag)
	ImageRef string `json:"image_ref"`

	// Auth contains registry authentication credentials
	Auth *RegistryAuth `json:"auth,omitempty"`

	// Insecure allows connections to insecure (HTTP) registries
	Insecure bool `json:"insecure,omitempty"`

	// Context for request cancellation and timeouts
	Context context.Context `json:"-"`

	// Timeout for push operations
	Timeout time.Duration `json:"timeout,omitempty"`

	// Registry is the target registry hostname
	Registry string `json:"registry,omitempty"`

	// Repository is the target repository name
	Repository string `json:"repository,omitempty"`

	// Tag is the target image tag
	Tag string `json:"tag,omitempty"`
}

// Validate checks if the push options are valid and complete.
func (o *PushOptions) Validate() error {
	if o == nil {
		return NewRegistryError(400, "push options cannot be nil")
	}

	if o.ImageRef == "" {
		return NewRegistryError(400, "image reference is required")
	}

	return nil
}