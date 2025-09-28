package mapper

import (
	"context"
)

// StackMapper defines the interface for mapping stack definitions to containers
type StackMapper interface {
	// MapStackToContainers converts a stack configuration to container specifications
	MapStackToContainers(ctx context.Context, config StackConfig) (MappingResult, error)

	// ResolveReferences resolves all references in a stack configuration
	ResolveReferences(ctx context.Context, refs []string) (map[string]ComponentRef, error)

	// ValidateMapping validates a mapping result for consistency
	ValidateMapping(mapping MappingResult) error
}

// NewStackMapper creates a new StackMapper implementation
func NewStackMapper() StackMapper {
	return &mapperImpl{}
}

type mapperImpl struct{}