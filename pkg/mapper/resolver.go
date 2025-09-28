package mapper

import (
	"context"
	"fmt"
	"os"
	"strings"
)

// ResolveReferences resolves all references in a stack configuration
func (m *mapperImpl) ResolveReferences(ctx context.Context, refs []string) (map[string]ComponentRef, error) {
	if len(refs) == 0 {
		return make(map[string]ComponentRef), nil
	}

	resolved := make(map[string]ComponentRef, len(refs))

	for _, ref := range refs {
		componentRef, err := m.parseReference(ref)
		if err != nil {
			return nil, fmt.Errorf("parsing reference %s: %w", ref, err)
		}

		// Validate the parsed reference
		if err := m.validateReference(componentRef); err != nil {
			return nil, fmt.Errorf("validating reference %s: %w", ref, err)
		}

		resolved[ref] = componentRef
	}

	return resolved, nil
}

// parseReference parses a component source reference into its parts
func (m *mapperImpl) parseReference(ref string) (ComponentRef, error) {
	// Handle different reference formats:
	// 1. registry.example.com/repo/image:tag
	// 2. registry.example.com/repo/image@sha256:digest
	// 3. repo/image:tag (assumes default registry)
	// 4. image:tag (assumes default registry and namespace)

	componentRef := ComponentRef{}

	// Check for digest format first
	if strings.Contains(ref, "@") {
		parts := strings.SplitN(ref, "@", 2)
		if len(parts) != 2 {
			return componentRef, newMappingError(ErrInvalidReference,
				fmt.Sprintf("invalid digest format: %s", ref))
		}
		componentRef.Digest = parts[1]
		ref = parts[0] // Continue parsing the registry/repo part
	}

	// Handle tag format
	var repoPath string
	if strings.Contains(ref, ":") && !strings.Contains(ref, "://") {
		parts := strings.SplitN(ref, ":", 2)
		if len(parts) == 2 {
			repoPath = parts[0]
			componentRef.Tag = parts[1]
		} else {
			repoPath = ref
		}
	} else {
		repoPath = ref
		if componentRef.Tag == "" && componentRef.Digest == "" {
			componentRef.Tag = "latest" // Default tag
		}
	}

	// Parse registry and repository
	if err := m.parseRepositoryPath(repoPath, &componentRef); err != nil {
		return componentRef, err
	}

	return componentRef, nil
}

// parseRepositoryPath parses the repository path to extract registry and repository
func (m *mapperImpl) parseRepositoryPath(repoPath string, ref *ComponentRef) error {
	// Split by '/' to identify registry vs repository
	parts := strings.Split(repoPath, "/")

	switch len(parts) {
	case 1:
		// Just image name: image
		ref.Registry = m.getDefaultRegistry()
		ref.Repository = "library/" + parts[0] // Default namespace
	case 2:
		// Two parts could be:
		// 1. namespace/image
		// 2. registry.com/image
		if m.isRegistryHost(parts[0]) {
			ref.Registry = parts[0]
			ref.Repository = "library/" + parts[1] // Default namespace
		} else {
			ref.Registry = m.getDefaultRegistry()
			ref.Repository = strings.Join(parts, "/")
		}
	case 3:
		// Three parts: registry.com/namespace/image
		ref.Registry = parts[0]
		ref.Repository = strings.Join(parts[1:], "/")
	default:
		// More than 3 parts: registry.com/namespace/subpath/image
		ref.Registry = parts[0]
		ref.Repository = strings.Join(parts[1:], "/")
	}

	return nil
}

// isRegistryHost determines if a string looks like a registry hostname
func (m *mapperImpl) isRegistryHost(host string) bool {
	// Simple heuristic: contains a dot (domain) or is a known registry
	return strings.Contains(host, ".") ||
		   host == "localhost" ||
		   strings.Contains(host, ":")
}

// getDefaultRegistry returns the default registry from configuration
func (m *mapperImpl) getDefaultRegistry() string {
	// Try environment variable first
	if registry := strings.TrimSpace(os.Getenv("MAPPER_DEFAULT_REGISTRY")); registry != "" {
		return registry
	}
	// Fall back to Docker Hub
	return "docker.io"
}

// validateReference validates a parsed component reference
func (m *mapperImpl) validateReference(ref ComponentRef) error {
	if ref.Registry == "" {
		return newMappingError(ErrInvalidReference, "registry is required")
	}
	if ref.Repository == "" {
		return newMappingError(ErrInvalidReference, "repository is required")
	}
	if ref.Tag == "" && ref.Digest == "" {
		return newMappingError(ErrInvalidReference, "either tag or digest is required")
	}

	// Validate tag format
	if ref.Tag != "" {
		if err := m.validateTag(ref.Tag); err != nil {
			return err
		}
	}

	// Validate digest format
	if ref.Digest != "" {
		if err := m.validateDigest(ref.Digest); err != nil {
			return err
		}
	}

	return nil
}

// validateTag validates an image tag format
func (m *mapperImpl) validateTag(tag string) error {
	if tag == "" {
		return newMappingError(ErrInvalidReference, "tag cannot be empty")
	}
	if len(tag) > 128 {
		return newMappingError(ErrInvalidReference, "tag too long (max 128 characters)")
	}
	// Tag must start with alphanumeric and can contain .-_
	if !isValidTagCharacter(tag[0]) {
		return newMappingError(ErrInvalidReference, "tag must start with alphanumeric character")
	}
	for _, char := range tag {
		if !isValidTagCharacter(byte(char)) && char != '.' && char != '-' && char != '_' {
			return newMappingError(ErrInvalidReference,
				fmt.Sprintf("invalid character in tag: %c", char))
		}
	}
	return nil
}

// validateDigest validates a digest format (sha256:...)
func (m *mapperImpl) validateDigest(digest string) error {
	if !strings.HasPrefix(digest, "sha256:") {
		return newMappingError(ErrInvalidReference, "digest must start with 'sha256:'")
	}
	hashPart := strings.TrimPrefix(digest, "sha256:")
	if len(hashPart) != 64 {
		return newMappingError(ErrInvalidReference, "sha256 digest must be 64 characters")
	}
	// Validate hex characters
	for _, char := range hashPart {
		if !isHexCharacter(byte(char)) {
			return newMappingError(ErrInvalidReference,
				fmt.Sprintf("invalid hex character in digest: %c", char))
		}
	}
	return nil
}

// isValidTagCharacter checks if a character is valid in a tag
func isValidTagCharacter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}

// isHexCharacter checks if a character is a valid hex digit
func isHexCharacter(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}