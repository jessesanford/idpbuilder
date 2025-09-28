package mapper

import (
	"context"
	"fmt"
	"os"
	"strings"
)

// MapStackToContainers converts a stack configuration to container specifications
func (m *mapperImpl) MapStackToContainers(ctx context.Context, config StackConfig) (MappingResult, error) {
	// Validate input configuration
	if err := m.validateStackConfig(config); err != nil {
		return MappingResult{}, err
	}

	// Get configuration from context or environment
	registryURL := m.getRegistryURL(ctx)
	if registryURL == "" {
		return MappingResult{}, newMappingError(ErrInvalidConfig, "registry URL not configured")
	}

	// Initialize mapping result
	result := MappingResult{
		Containers: make([]ContainerSpec, 0, len(config.Components)),
		Metadata:   make(map[string]string),
	}

	// Copy stack metadata
	for k, v := range config.Metadata {
		result.Metadata[k] = v
	}
	result.Metadata["stack_name"] = config.Name
	result.Metadata["stack_version"] = config.Version

	// Process each component
	for _, comp := range config.Components {
		spec, err := m.processComponent(comp, registryURL)
		if err != nil {
			return MappingResult{}, fmt.Errorf("processing component %s: %w", comp.Name, err)
		}
		result.Containers = append(result.Containers, spec)
	}

	// Generate OCI manifest
	manifest, err := m.generateManifest(config, result.Containers)
	if err != nil {
		return MappingResult{}, fmt.Errorf("generating manifest: %w", err)
	}
	result.Manifest = manifest

	return result, nil
}

// validateStackConfig validates the basic structure of a stack configuration
func (m *mapperImpl) validateStackConfig(config StackConfig) error {
	if config.Name == "" {
		return newMappingError(ErrInvalidConfig, "stack name is required")
	}
	if config.Version == "" {
		return newMappingError(ErrInvalidConfig, "stack version is required")
	}
	if len(config.Components) == 0 {
		return newMappingError(ErrInvalidConfig, "stack must have at least one component")
	}

	// Validate each component
	for i, comp := range config.Components {
		if comp.Name == "" {
			return newMappingError(ErrInvalidConfig, fmt.Sprintf("component %d: name is required", i))
		}
		if comp.Type == "" {
			return newMappingError(ErrInvalidConfig, fmt.Sprintf("component %s: type is required", comp.Name))
		}
		if comp.Source == "" {
			return newMappingError(ErrInvalidConfig, fmt.Sprintf("component %s: source is required", comp.Name))
		}
	}

	return nil
}

// getRegistryURL retrieves the registry URL from context or environment
func (m *mapperImpl) getRegistryURL(ctx context.Context) string {
	// Try to get from context first
	if registryURL := ctx.Value("registry"); registryURL != nil {
		if url, ok := registryURL.(string); ok && url != "" {
			return url
		}
	}

	// Fall back to environment variable
	return os.Getenv("MAPPER_REGISTRY_URL")
}

// processComponent processes a single component into a container specification
func (m *mapperImpl) processComponent(comp Component, registryURL string) (ContainerSpec, error) {
	spec := ContainerSpec{
		Name:      comp.Name,
		BaseImage: m.resolveBaseImage(comp.Type),
		Layers:    make([]Layer, 0),
		Env:       make(map[string]string),
		Labels:    make(map[string]string),
	}

	// Set standard labels
	spec.Labels["component.name"] = comp.Name
	spec.Labels["component.type"] = comp.Type
	spec.Labels["component.source"] = comp.Source

	// Process component configuration
	if comp.Config != nil {
		m.processComponentConfig(comp.Config, &spec)
	}

	// Create layer for the component
	layer := Layer{
		MediaType: "application/vnd.oci.image.layer.v1.tar",
		Size:      0, // Size will be determined during build
		Digest:    "", // Digest will be calculated during build
	}
	spec.Layers = append(spec.Layers, layer)

	return spec, nil
}

// resolveBaseImage determines the base image for a component type
func (m *mapperImpl) resolveBaseImage(componentType string) string {
	switch strings.ToLower(componentType) {
	case "web", "frontend":
		return "nginx:alpine"
	case "api", "backend", "service":
		return "alpine:latest"
	case "database", "db":
		return "postgres:alpine"
	case "cache":
		return "redis:alpine"
	default:
		return "alpine:latest"
	}
}

// processComponentConfig processes component configuration into container spec
func (m *mapperImpl) processComponentConfig(config map[string]any, spec *ContainerSpec) {
	// Process environment variables
	if env, ok := config["env"]; ok {
		if envMap, ok := env.(map[string]any); ok {
			for k, v := range envMap {
				if strVal, ok := v.(string); ok {
					spec.Env[k] = strVal
				}
			}
		}
	}

	// Process labels
	if labels, ok := config["labels"]; ok {
		if labelMap, ok := labels.(map[string]any); ok {
			for k, v := range labelMap {
				if strVal, ok := v.(string); ok {
					spec.Labels[k] = strVal
				}
			}
		}
	}

	// Process ports (as labels for now)
	if ports, ok := config["ports"]; ok {
		if portSlice, ok := ports.([]any); ok {
			var portStrs []string
			for _, port := range portSlice {
				if portStr, ok := port.(string); ok {
					portStrs = append(portStrs, portStr)
				}
			}
			if len(portStrs) > 0 {
				spec.Labels["component.ports"] = strings.Join(portStrs, ",")
			}
		}
	}
}

// generateManifest creates an OCI manifest for the mapped containers
func (m *mapperImpl) generateManifest(config StackConfig, containers []ContainerSpec) (*PackageManifest, error) {
	manifest := &PackageManifest{
		SchemaVersion: 2,
		MediaType:     "application/vnd.oci.image.manifest.v1+json",
		Config: Descriptor{
			MediaType: "application/vnd.oci.image.config.v1+json",
			Size:      0, // Will be calculated
			Digest:    "", // Will be calculated
		},
		Layers: make([]LayerDescriptor, 0, len(containers)),
	}

	// Add layer descriptors for each container
	for _, container := range containers {
		for _, layer := range container.Layers {
			layerDesc := LayerDescriptor{
				MediaType: layer.MediaType,
				Size:      layer.Size,
				Digest:    layer.Digest,
			}
			manifest.Layers = append(manifest.Layers, layerDesc)
		}
	}

	// Set annotations
	if manifest.Annotations == nil {
		manifest.Annotations = make(map[string]string)
	}
	manifest.Annotations["org.opencontainers.image.title"] = config.Name
	manifest.Annotations["org.opencontainers.image.version"] = config.Version
	manifest.Annotations["org.opencontainers.image.created"] = "" // Will be set during build

	return manifest, nil
}

// ValidateMapping validates a mapping result for consistency
func (m *mapperImpl) ValidateMapping(mapping MappingResult) error {
	if len(mapping.Containers) == 0 {
		return newMappingError(ErrValidationFailed, "mapping result must contain at least one container")
	}

	// Validate each container specification
	for i, container := range mapping.Containers {
		if err := m.validateContainerSpec(container, i); err != nil {
			return err
		}
	}

	// Validate manifest if present
	if mapping.Manifest != nil {
		if err := m.validateManifest(mapping.Manifest); err != nil {
			return fmt.Errorf("manifest validation failed: %w", err)
		}
	}

	return nil
}

// validateContainerSpec validates a single container specification
func (m *mapperImpl) validateContainerSpec(spec ContainerSpec, index int) error {
	if spec.Name == "" {
		return newMappingError(ErrValidationFailed,
			fmt.Sprintf("container %d: name is required", index))
	}
	if spec.BaseImage == "" {
		return newMappingError(ErrValidationFailed,
			fmt.Sprintf("container %s: base image is required", spec.Name))
	}

	// Validate environment variables
	for key, value := range spec.Env {
		if key == "" {
			return newMappingError(ErrValidationFailed,
				fmt.Sprintf("container %s: environment variable key cannot be empty", spec.Name))
		}
		if strings.Contains(key, "=") {
			return newMappingError(ErrValidationFailed,
				fmt.Sprintf("container %s: environment variable key cannot contain '='", spec.Name))
		}
		_ = value // Value can be empty
	}

	// Validate labels
	for key := range spec.Labels {
		if key == "" {
			return newMappingError(ErrValidationFailed,
				fmt.Sprintf("container %s: label key cannot be empty", spec.Name))
		}
	}

	return nil
}

// validateManifest validates an OCI manifest structure
func (m *mapperImpl) validateManifest(manifest *PackageManifest) error {
	if manifest.SchemaVersion == 0 {
		return newMappingError(ErrValidationFailed, "manifest schema version is required")
	}
	if manifest.MediaType == "" {
		return newMappingError(ErrValidationFailed, "manifest media type is required")
	}

	// Validate config descriptor
	if manifest.Config.MediaType == "" {
		return newMappingError(ErrValidationFailed, "config media type is required")
	}

	return nil
}