// Package providers defines the contract for OCI registry providers.
//
// This package contains the core interfaces and types needed for implementing
// OCI registry operations. The Provider interface defines the fundamental
// operations that any registry implementation must support: pushing, pulling,
// listing, and deleting OCI artifacts.
//
// All implementations must handle context cancellation and return appropriate
// errors for different failure scenarios.
package providers

import (
	"context"
)

// Provider defines the contract for OCI registry providers.
// Implementations must support all core registry operations with proper
// context handling and error reporting.
type Provider interface {
	// Push uploads an OCI artifact to the registry at the specified reference.
	// The reference should be in the format registry/repository:tag or
	// registry/repository@digest.
	//
	// Returns ProviderError for registry-specific failures.
	Push(ctx context.Context, ref string, artifact Artifact) error

	// Pull downloads an OCI artifact from the registry at the specified reference.
	// The reference should be in the format registry/repository:tag or
	// registry/repository@digest.
	//
	// Returns the artifact and nil error on success, or zero value and
	// ProviderError on failure.
	Pull(ctx context.Context, ref string) (Artifact, error)

	// List retrieves information about all artifacts in the specified repository.
	// The repository should be in the format registry/repository.
	//
	// Returns a slice of artifact information or ProviderError on failure.
	List(ctx context.Context, repository string) ([]ArtifactInfo, error)

	// Delete removes an artifact from the registry at the specified reference.
	// The reference should be in the format registry/repository:tag or
	// registry/repository@digest.
	//
	// Returns ProviderError if the deletion fails or the artifact doesn't exist.
	Delete(ctx context.Context, ref string) error
}