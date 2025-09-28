package registry

import (
	"context"
)

// Artifact represents an OCI artifact for registry operations.
// This is a simplified interface for the provider types.
type Artifact struct {
	Repository string
	Tag        string
	Digest     string
	Layers     []Layer
	Config     []byte
	Manifest   []byte
}

// Layer represents an OCI layer within an artifact.
type Layer struct {
	Digest      string
	Size        int64
	MediaType   string
	Annotations map[string]string
}

// Client defines the interface for registry operations.
// This abstraction wraps go-containerregistry for testability and provides
// a clean interface for pushing, pulling, and managing OCI artifacts.
type Client interface {
	// Push uploads an artifact to the registry at the specified reference.
	// The reference should include registry, repository, and tag (e.g., "registry.example.com/repo:tag").
	// Returns an error if the push operation fails.
	Push(ctx context.Context, ref string, artifact *Artifact) error

	// Pull retrieves an artifact from the registry at the specified reference.
	// The reference should include registry, repository, and tag (e.g., "registry.example.com/repo:tag").
	// Returns the artifact and any error that occurred during the pull operation.
	Pull(ctx context.Context, ref string) (*Artifact, error)

	// Exists checks if an artifact exists in the registry at the specified reference.
	// The reference should include registry, repository, and tag (e.g., "registry.example.com/repo:tag").
	// Returns true if the artifact exists, false otherwise, and any error that occurred.
	Exists(ctx context.Context, ref string) (bool, error)

	// ListTags returns all tags for a repository in the registry.
	// The repository should include registry and repository name (e.g., "registry.example.com/repo").
	// Returns a slice of tag names and any error that occurred.
	ListTags(ctx context.Context, repository string) ([]string, error)
}