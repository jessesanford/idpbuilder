// Package builder provides functionality for building OCI container images
// using the go-containerregistry library. It supports building images from
// Dockerfile-like configurations, managing build contexts, and handling
// various build options including multi-platform builds.
//
// The builder package is designed to be flexible and extensible, supporting
// various image build scenarios while maintaining compatibility with the
// OCI Image Specification and Docker Registry API.
//
// Key Components:
//   - BuildConfig: Core configuration for image building
//   - ImageOptions: Flexible options for image customization
//   - RegistryConfig: Registry authentication and connection settings
//   - PlatformConfig: Platform-specific build configurations
//
// Example usage:
//
//	config := &BuildConfig{
//		ContextPath: "./docker-context",
//		Dockerfile:  "Dockerfile",
//		Tags:        []string{"myapp:latest"},
//	}
//
//	opts := []BuildOption{
//		WithPlatform("linux/amd64"),
//		WithRegistry(registryConfig),
//	}
//
//	builder := NewBuilder(config, opts...)
//	if err := builder.Build(ctx); err != nil {
//		log.Fatal(err)
//	}
package builder