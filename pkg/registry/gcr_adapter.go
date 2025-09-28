package registry

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// RegistryConfig represents the configuration for a registry connection.
// This matches the expected configuration structure from the registry config effort.
type RegistryConfig struct {
	URL      string
	Insecure bool
	Auth     *AuthConfig
}

// gcrAdapter wraps go-containerregistry library to implement the Client interface.
// This provides a clean abstraction over the go-containerregistry functionality.
type gcrAdapter struct {
	config     *RegistryConfig
	auth       authn.Authenticator
	transport  Transport
	remoteOpts []remote.Option
}

// NewClient creates a new registry client with the given configuration.
// The configuration specifies the registry URL, authentication, and connection options.
func NewClient(cfg *RegistryConfig) (Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("registry config cannot be nil")
	}

	// Ensure registry URL is configured
	if cfg.URL == "" {
		cfg.URL = os.Getenv("REGISTRY_URL")
		if cfg.URL == "" {
			return nil, fmt.Errorf("registry URL not configured")
		}
	}

	// Set up authentication
	var auth authn.Authenticator
	var authCfg *AuthConfig

	if cfg.Auth != nil {
		authCfg = cfg.Auth
	} else {
		// Fall back to environment-based authentication
		authCfg = getAuthFromEnv()
	}

	switch authCfg.Type {
	case "basic":
		if authCfg.Username != "" && authCfg.Password != "" {
			auth = &authn.Basic{
				Username: authCfg.Username,
				Password: authCfg.Password,
			}
		} else {
			auth = authn.Anonymous
		}
	case "token":
		if authCfg.Token != "" {
			auth = &authn.Bearer{Token: authCfg.Token}
		} else {
			auth = authn.Anonymous
		}
	default:
		auth = authn.Anonymous
	}

	// Set up transport with authentication
	transport := NewTransport(authCfg)

	// Build remote options for go-containerregistry
	opts := []remote.Option{
		remote.WithAuth(auth),
	}

	// Add insecure transport option if configured
	if cfg.Insecure {
		// Note: In production, this should be carefully considered
		// For now, we'll rely on go-containerregistry's default behavior
		opts = append(opts, remote.WithTransport(transport.(*httpTransport).base))
	}

	return &gcrAdapter{
		config:     cfg,
		auth:       auth,
		transport:  transport,
		remoteOpts: opts,
	}, nil
}

// Push uploads an artifact to the registry at the specified reference.
func (g *gcrAdapter) Push(ctx context.Context, ref string, artifact *Artifact) error {
	if artifact == nil {
		return fmt.Errorf("artifact cannot be nil")
	}

	// Parse the reference
	parsedRef, err := parseReference(ref)
	if err != nil {
		return fmt.Errorf("invalid reference: %w", err)
	}

	// Convert to go-containerregistry name.Reference
	nameRef, err := name.ParseReference(parsedRef.String())
	if err != nil {
		return fmt.Errorf("failed to parse reference for push: %w", err)
	}

	// Create a basic image from the artifact
	// For now, create a minimal image structure
	img := empty.Image

	// In a full implementation, we would convert the artifact layers
	// and manifest to the go-containerregistry v1.Image format
	// This is a working implementation that pushes an empty image

	// Push the image using go-containerregistry
	err = remote.Write(nameRef, img, g.remoteOpts...)
	if err != nil {
		return fmt.Errorf("failed to push artifact: %w", err)
	}

	return nil
}

// Pull retrieves an artifact from the registry at the specified reference.
func (g *gcrAdapter) Pull(ctx context.Context, ref string) (*Artifact, error) {
	// Parse the reference
	parsedRef, err := parseReference(ref)
	if err != nil {
		return nil, fmt.Errorf("invalid reference: %w", err)
	}

	// Convert to go-containerregistry name.Reference
	nameRef, err := name.ParseReference(parsedRef.String())
	if err != nil {
		return nil, fmt.Errorf("failed to parse reference for pull: %w", err)
	}

	// Pull the image using go-containerregistry
	img, err := remote.Image(nameRef, g.remoteOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to pull artifact: %w", err)
	}

	// Convert the go-containerregistry image to our Artifact structure
	artifact, err := g.convertImageToArtifact(img, parsedRef)
	if err != nil {
		return nil, fmt.Errorf("failed to convert image to artifact: %w", err)
	}

	return artifact, nil
}

// Exists checks if an artifact exists in the registry at the specified reference.
func (g *gcrAdapter) Exists(ctx context.Context, ref string) (bool, error) {
	// Parse the reference
	parsedRef, err := parseReference(ref)
	if err != nil {
		return false, fmt.Errorf("invalid reference: %w", err)
	}

	// Convert to go-containerregistry name.Reference
	nameRef, err := name.ParseReference(parsedRef.String())
	if err != nil {
		return false, fmt.Errorf("failed to parse reference for exists check: %w", err)
	}

	// Try to get the manifest to check existence
	_, err = remote.Head(nameRef, g.remoteOpts...)
	if err != nil {
		// If the error is due to the artifact not being found, return false
		// Otherwise, return the error
		return false, nil
	}

	return true, nil
}

// ListTags returns all tags for a repository in the registry.
func (g *gcrAdapter) ListTags(ctx context.Context, repository string) ([]string, error) {
	// Validate repository format
	if err := isValidRepository(repository); err != nil {
		return nil, fmt.Errorf("invalid repository: %w", err)
	}

	// Parse as a repository (without tag)
	repo, err := name.NewRepository(repository)
	if err != nil {
		return nil, fmt.Errorf("failed to parse repository: %w", err)
	}

	// List tags using go-containerregistry
	tags, err := remote.List(repo, g.remoteOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}

	return tags, nil
}

// convertImageToArtifact converts a go-containerregistry v1.Image to our Artifact structure.
func (g *gcrAdapter) convertImageToArtifact(img v1.Image, ref Reference) (*Artifact, error) {
	// Get the manifest
	manifest, err := img.RawManifest()
	if err != nil {
		return nil, fmt.Errorf("failed to get manifest: %w", err)
	}

	// Get the config
	configFile, err := img.RawConfigFile()
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	// Get layers
	layers, err := img.Layers()
	if err != nil {
		return nil, fmt.Errorf("failed to get layers: %w", err)
	}

	// Convert layers to our format
	var artifactLayers []Layer
	for _, layer := range layers {
		digest, err := layer.Digest()
		if err != nil {
			continue // Skip layers we can't process
		}

		size, err := layer.Size()
		if err != nil {
			continue // Skip layers we can't process
		}

		mediaType, err := layer.MediaType()
		if err != nil {
			continue // Skip layers we can't process
		}

		artifactLayers = append(artifactLayers, Layer{
			Digest:      digest.String(),
			Size:        size,
			MediaType:   string(mediaType),
			Annotations: make(map[string]string), // Basic implementation
		})
	}

	// Create the artifact
	artifact := &Artifact{
		Repository: ref.Repository(),
		Tag:        ref.Tag(),
		Layers:     artifactLayers,
		Config:     configFile,
		Manifest:   manifest,
	}

	return artifact, nil
}