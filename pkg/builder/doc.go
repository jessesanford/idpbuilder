// Package builder provides functionality for building OCI container images
// using the go-containerregistry library as the foundation.
//
// This package implements a builder pattern for creating container images
// with support for multiple architectures, custom configurations, and
// feature flags to control functionality availability.
//
// Basic Usage:
//
//	builder, err := NewBuilder(options...)
//	if err != nil {
//		return err
//	}
//
//	image, err := builder.Build(ctx, contextDir, BuildOptions{
//		Platform: &v1.Platform{
//			Architecture: "amd64",
//			OS:           "linux",
//		},
//		Tags: []string{"my-app:latest"},
//	})
//	if err != nil {
//		return err
//	}
//
// Feature Flags:
//
// This package uses feature flags to control availability of functionality
// during development. Some features may be disabled in certain splits:
//
//	const (
//		FeatureTarballExport = "tarball_export" // May be disabled in early splits
//		FeatureLayerCaching  = "layer_caching"   // May be disabled in early splits
//		FeatureMultiLayer    = "multi_layer"     // May be disabled in early splits
//	)
//
// Configuration:
//
// The package provides a ConfigFactory for generating OCI image configurations
// with platform-specific settings, labels, and environment variables.
//
// Platform Support:
//
// Multi-architecture builds are supported through platform configuration.
// The builder can generate images for different OS/architecture combinations.
package builder