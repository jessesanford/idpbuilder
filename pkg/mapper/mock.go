package mapper

import (
	"context"
	"errors"
)

// MockStackMapper provides a mock implementation of StackMapper for testing
type MockStackMapper struct {
	MapStackToContainersFunc func(ctx context.Context, config StackConfig) (MappingResult, error)
	ResolveReferencesFunc    func(ctx context.Context, refs []string) (map[string]ComponentRef, error)
	ValidateMappingFunc      func(mapping MappingResult) error

	// Configuration for mock behavior
	ShouldFailMapping    bool
	ShouldFailResolving  bool
	ShouldFailValidation bool
	ErrorToReturn        error
}

// NewMockStackMapper creates a new mock stack mapper with default behavior
func NewMockStackMapper() *MockStackMapper {
	return &MockStackMapper{
		ShouldFailMapping:    false,
		ShouldFailResolving:  false,
		ShouldFailValidation: false,
	}
}

// MapStackToContainers implements the StackMapper interface
func (m *MockStackMapper) MapStackToContainers(ctx context.Context, config StackConfig) (MappingResult, error) {
	if m.MapStackToContainersFunc != nil {
		return m.MapStackToContainersFunc(ctx, config)
	}

	if m.ShouldFailMapping {
		if m.ErrorToReturn != nil {
			return MappingResult{}, m.ErrorToReturn
		}
		return MappingResult{}, errors.New("mock mapping failure")
	}

	// Default mock behavior - create a simple mapping result
	containers := make([]ContainerSpec, len(config.Components))
	for i, comp := range config.Components {
		containers[i] = ContainerSpec{
			Name:      comp.Name,
			BaseImage: "alpine:latest",
			Env:       make(map[string]string),
			Labels:    map[string]string{"component.name": comp.Name},
		}
	}

	manifest := &PackageManifest{
		SchemaVersion: 2,
		MediaType:     "application/vnd.oci.image.manifest.v1+json",
		Config: Descriptor{
			MediaType: "application/vnd.oci.image.config.v1+json",
		},
		Layers: []LayerDescriptor{},
		Annotations: map[string]string{
			"org.opencontainers.image.title":   config.Name,
			"org.opencontainers.image.version": config.Version,
		},
	}

	return MappingResult{
		Containers: containers,
		Manifest:   manifest,
		Metadata:   map[string]string{"mock": "true"},
	}, nil
}

// ResolveReferences implements the StackMapper interface
func (m *MockStackMapper) ResolveReferences(ctx context.Context, refs []string) (map[string]ComponentRef, error) {
	if m.ResolveReferencesFunc != nil {
		return m.ResolveReferencesFunc(ctx, refs)
	}

	if m.ShouldFailResolving {
		if m.ErrorToReturn != nil {
			return nil, m.ErrorToReturn
		}
		return nil, errors.New("mock resolve failure")
	}

	// Default mock behavior - create simple resolved references
	resolved := make(map[string]ComponentRef, len(refs))
	for _, ref := range refs {
		resolved[ref] = ComponentRef{
			Registry:   "docker.io",
			Repository: "library/alpine",
			Tag:        "latest",
			Digest:     "",
		}
	}

	return resolved, nil
}

// ValidateMapping implements the StackMapper interface
func (m *MockStackMapper) ValidateMapping(mapping MappingResult) error {
	if m.ValidateMappingFunc != nil {
		return m.ValidateMappingFunc(mapping)
	}

	if m.ShouldFailValidation {
		if m.ErrorToReturn != nil {
			return m.ErrorToReturn
		}
		return errors.New("mock validation failure")
	}

	// Default mock behavior - always pass validation
	return nil
}

// WithMappingError configures the mock to return an error for mapping operations
func (m *MockStackMapper) WithMappingError(err error) *MockStackMapper {
	m.ShouldFailMapping = true
	m.ErrorToReturn = err
	return m
}

// WithResolvingError configures the mock to return an error for resolving operations
func (m *MockStackMapper) WithResolvingError(err error) *MockStackMapper {
	m.ShouldFailResolving = true
	m.ErrorToReturn = err
	return m
}

// WithValidationError configures the mock to return an error for validation operations
func (m *MockStackMapper) WithValidationError(err error) *MockStackMapper {
	m.ShouldFailValidation = true
	m.ErrorToReturn = err
	return m
}