package registry

import (
	"fmt"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
)

// Reference represents a parsed container image reference.
// This interface provides access to the components of an OCI artifact reference.
type Reference interface {
	// Parse parses a reference string and updates the internal state.
	// Returns an error if the reference format is invalid.
	Parse(ref string) error

	// Registry returns the registry hostname from the reference.
	Registry() string

	// Repository returns the repository path from the reference.
	Repository() string

	// Tag returns the tag from the reference, or "latest" if not specified.
	Tag() string

	// String returns the full reference string.
	String() string
}

// reference wraps go-containerregistry name.Reference to provide our interface.
type reference struct {
	ref name.Reference
}

// parseReference parses a reference string and returns a Reference interface.
// The reference string should be in the format: [registry/]repository[:tag|@digest]
func parseReference(ref string) (Reference, error) {
	if ref == "" {
		return nil, fmt.Errorf("reference cannot be empty")
	}

	// Use go-containerregistry to parse the reference
	parsedRef, err := name.ParseReference(ref)
	if err != nil {
		return nil, fmt.Errorf("failed to parse reference %q: %w", ref, err)
	}

	return &reference{ref: parsedRef}, nil
}

// Parse parses a reference string and updates the internal state.
// This allows reusing a Reference instance with different reference strings.
func (r *reference) Parse(ref string) error {
	if ref == "" {
		return fmt.Errorf("reference cannot be empty")
	}

	parsedRef, err := name.ParseReference(ref)
	if err != nil {
		return fmt.Errorf("failed to parse reference %q: %w", ref, err)
	}

	r.ref = parsedRef
	return nil
}

// Registry returns the registry hostname from the reference.
// If no registry is specified, returns the default registry.
func (r *reference) Registry() string {
	if r.ref == nil {
		return ""
	}
	return r.ref.Context().RegistryStr()
}

// Repository returns the repository path from the reference.
// This includes the namespace and repository name but excludes the registry.
func (r *reference) Repository() string {
	if r.ref == nil {
		return ""
	}
	return r.ref.Context().RepositoryStr()
}

// Tag returns the tag from the reference.
// If the reference uses a digest instead of a tag, returns empty string.
// If no tag is specified, returns "latest".
func (r *reference) Tag() string {
	if r.ref == nil {
		return ""
	}

	// Check if this is a tagged reference
	if tagged, ok := r.ref.(name.Tag); ok {
		return tagged.TagStr()
	}

	// If it's a digest reference, return empty string
	return ""
}

// String returns the full reference string in canonical format.
func (r *reference) String() string {
	if r.ref == nil {
		return ""
	}
	return r.ref.String()
}

// isValidRepository checks if a repository string is valid for listing tags.
// Repository should not include a tag or digest, just registry/repository.
func isValidRepository(repository string) error {
	if repository == "" {
		return fmt.Errorf("repository cannot be empty")
	}

	// Should not contain @ (digest) or : (tag) at the end
	if strings.Contains(repository, "@") {
		return fmt.Errorf("repository should not contain digest (@)")
	}

	// Check for tag (: followed by something that's not a port)
	parts := strings.Split(repository, ":")
	if len(parts) > 2 {
		return fmt.Errorf("repository should not contain tag (:)")
	}

	// If there's exactly one colon, it should be for the port
	if len(parts) == 2 {
		// Split by / to check if the colon is in the registry part (port) or repository part (tag)
		pathParts := strings.Split(repository, "/")
		if len(pathParts) > 1 {
			registryPart := pathParts[0]
			// If the colon is not in the first part (registry), it's likely a tag
			if !strings.Contains(registryPart, ":") {
				return fmt.Errorf("repository should not contain tag (:)")
			}
		}
	}

	return nil
}