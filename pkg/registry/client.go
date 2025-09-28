// Package registry provides a concrete registry client implementation that bridges
// the Provider interface with go-containerregistry for OCI registry operations.
//
// This client abstracts the complexity of different registry types and provides
// a unified interface for registry operations across different implementations.
package registry

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/types"

	"github.com/cnoe-io/idpbuilder/pkg/config"
	"github.com/cnoe-io/idpbuilder/pkg/providers"
)

// Client implements the Provider interface using go-containerregistry.
// It provides a concrete implementation for OCI registry operations
// with support for authentication, transport configuration, and error handling.
type Client struct {
	config     *config.RegistryConfig
	httpClient *http.Client
	auth       authn.Authenticator
	options    []remote.Option
}

// NewClient creates a new registry client with the provided configuration.
// It sets up authentication, transport options, and validates the configuration.
func NewClient(cfg *config.RegistryConfig) (*Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("registry configuration cannot be nil")
	}

	if cfg.URL == "" {
		return nil, fmt.Errorf("registry URL cannot be empty")
	}

	client := &Client{
		config:     cfg,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}

	// Configure authentication
	if err := client.configureAuth(); err != nil {
		return nil, fmt.Errorf("failed to configure authentication: %w", err)
	}

	// Configure transport options
	client.configureTransport()

	return client, nil
}

// configureAuth sets up authentication based on the registry configuration.
func (c *Client) configureAuth() error {
	if c.config.Auth.Username != "" && c.config.Auth.Password != "" {
		c.auth = &authn.Basic{
			Username: c.config.Auth.Username,
			Password: c.config.Auth.Password,
		}
	} else if c.config.Auth.Token != "" {
		c.auth = &authn.Bearer{Token: c.config.Auth.Token}
	} else {
		c.auth = authn.Anonymous
	}

	c.options = append(c.options, remote.WithAuth(c.auth))
	return nil
}

// configureTransport sets up HTTP transport options including TLS configuration.
func (c *Client) configureTransport() {
	if c.config.Insecure {
		c.options = append(c.options, remote.WithTransport(&http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}))
	}

	c.options = append(c.options, remote.WithContext(context.Background()))
}

// Push implements the Provider interface Push method.
// It uploads an OCI artifact to the registry at the specified reference.
func (c *Client) Push(ctx context.Context, ref string, artifact providers.Artifact) error {
	// Parse the reference
	repository, err := name.ParseReference(ref)
	if err != nil {
		return &providers.ProviderError{
			Op:  "push",
			Ref: ref,
			Err: fmt.Errorf("invalid reference: %w", err),
		}
	}

	// Convert artifact to go-containerregistry image
	img, err := c.artifactToImage(artifact)
	if err != nil {
		return &providers.ProviderError{
			Op:  "push",
			Ref: ref,
			Err: fmt.Errorf("failed to convert artifact: %w", err),
		}
	}

	// Push the image
	options := append(c.options, remote.WithContext(ctx))
	if err := remote.Write(repository, img, options...); err != nil {
		return &providers.ProviderError{
			Op:  "push",
			Ref: ref,
			Err: fmt.Errorf("push failed: %w", err),
		}
	}

	return nil
}

// Pull implements the Provider interface Pull method.
// It downloads an OCI artifact from the registry at the specified reference.
func (c *Client) Pull(ctx context.Context, ref string) (providers.Artifact, error) {
	// Parse the reference
	repository, err := name.ParseReference(ref)
	if err != nil {
		return providers.Artifact{}, &providers.ProviderError{
			Op:  "pull",
			Ref: ref,
			Err: fmt.Errorf("invalid reference: %w", err),
		}
	}

	// Pull the image
	options := append(c.options, remote.WithContext(ctx))
	img, err := remote.Image(repository, options...)
	if err != nil {
		return providers.Artifact{}, &providers.ProviderError{
			Op:  "pull",
			Ref: ref,
			Err: fmt.Errorf("pull failed: %w", err),
		}
	}

	// Convert image to artifact
	artifact, err := c.imageToArtifact(img)
	if err != nil {
		return providers.Artifact{}, &providers.ProviderError{
			Op:  "pull",
			Ref: ref,
			Err: fmt.Errorf("failed to convert image: %w", err),
		}
	}

	return artifact, nil
}

// List implements the Provider interface List method.
// It retrieves information about all artifacts in the specified repository.
func (c *Client) List(ctx context.Context, repository string) ([]providers.ArtifactInfo, error) {
	// Parse repository
	repo, err := name.NewRepository(repository)
	if err != nil {
		return nil, &providers.ProviderError{
			Op:  "list",
			Ref: repository,
			Err: fmt.Errorf("invalid repository: %w", err),
		}
	}

	// List tags
	options := append(c.options, remote.WithContext(ctx))
	tags, err := remote.List(repo, options...)
	if err != nil {
		return nil, &providers.ProviderError{
			Op:  "list",
			Ref: repository,
			Err: fmt.Errorf("list failed: %w", err),
		}
	}

	// Convert tags to artifact info
	var artifacts []providers.ArtifactInfo
	for _, tag := range tags {
		tagRef := repo.Tag(tag)

		// Get manifest to extract metadata
		manifest, err := remote.Get(tagRef, options...)
		if err != nil {
			// Continue with basic info if manifest fetch fails
			artifacts = append(artifacts, providers.ArtifactInfo{
				Reference: tagRef.String(),
				Tags:      []string{tag},
			})
			continue
		}

		artifacts = append(artifacts, providers.ArtifactInfo{
			Reference:   tagRef.String(),
			Digest:      manifest.Digest.String(),
			Tags:        []string{tag},
			Size:        manifest.Size,
			Created:     time.Now(), // go-containerregistry doesn't provide creation time easily
			Annotations: manifest.Annotations,
		})
	}

	return artifacts, nil
}

// Delete implements the Provider interface Delete method.
// It removes an artifact from the registry at the specified reference.
func (c *Client) Delete(ctx context.Context, ref string) error {
	// Parse the reference
	repository, err := name.ParseReference(ref)
	if err != nil {
		return &providers.ProviderError{
			Op:  "delete",
			Ref: ref,
			Err: fmt.Errorf("invalid reference: %w", err),
		}
	}

	// Delete the image
	options := append(c.options, remote.WithContext(ctx))
	if err := remote.Delete(repository, options...); err != nil {
		return &providers.ProviderError{
			Op:  "delete",
			Ref: ref,
			Err: fmt.Errorf("delete failed: %w", err),
		}
	}

	return nil
}

// artifactToImage converts a providers.Artifact to a v1.Image for go-containerregistry.
func (c *Client) artifactToImage(artifact providers.Artifact) (v1.Image, error) {
	// This is a simplified conversion - in a real implementation,
	// we would need to properly construct the image from layers and manifest
	return nil, fmt.Errorf("artifact to image conversion not implemented")
}

// imageToArtifact converts a v1.Image to a providers.Artifact.
func (c *Client) imageToArtifact(img v1.Image) (providers.Artifact, error) {
	// Get manifest
	manifest, err := img.RawManifest()
	if err != nil {
		return providers.Artifact{}, fmt.Errorf("failed to get manifest: %w", err)
	}

	// Get config
	_, err = img.ConfigFile()
	if err != nil {
		return providers.Artifact{}, fmt.Errorf("failed to get config: %w", err)
	}

	// Get layers
	layers, err := img.Layers()
	if err != nil {
		return providers.Artifact{}, fmt.Errorf("failed to get layers: %w", err)
	}

	var artifactLayers []providers.Layer
	for _, layer := range layers {
		digest, err := layer.Digest()
		if err != nil {
			continue
		}

		size, err := layer.Size()
		if err != nil {
			continue
		}

		mediaType, err := layer.MediaType()
		if err != nil {
			continue
		}

		// Get layer data - this is expensive, might want to do lazily
		compressed, err := layer.Compressed()
		if err != nil {
			continue
		}

		data, err := io.ReadAll(compressed)
		if err != nil {
			continue
		}

		artifactLayers = append(artifactLayers, providers.Layer{
			MediaType: string(mediaType),
			Digest:    digest.String(),
			Size:      size,
			Data:      data,
		})
	}

	return providers.Artifact{
		MediaType:   string(types.OCIManifestSchema1),
		Manifest:    manifest,
		Layers:      artifactLayers,
		Config:      []byte{}, // ConfigFile needs to be serialized
		Annotations: map[string]string{},
	}, nil
}