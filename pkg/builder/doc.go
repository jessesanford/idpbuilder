// Package builder provides functionality for building OCI container images
// using the go-containerregistry library. It supports building images from
// Dockerfile-like configurations with multi-platform builds.
//
// Key Components:
//   - BuildConfig: Core configuration for image building
//   - ImageOptions: Options for image customization
//   - RegistryConfig: Registry authentication settings
//   - PlatformConfig: Platform-specific configurations
//
// Example:
//	config := &BuildConfig{
//		ContextPath: "./context",
//		Dockerfile:  "Dockerfile",
//		Tags:        []string{"myapp:latest"},
//	}
//	builder := NewBuilder(config, WithPlatform("linux/amd64"))
//	err := builder.Build(ctx)
package builder