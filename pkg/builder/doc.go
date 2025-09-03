// Package builder provides container image building capabilities for IDP Builder.
// Supports basic builds, multi-stage builds, layer caching, and registry integration.
//
// Basic usage:
//	config := &BuildConfig{BaseImage: "alpine:latest"}
//	builder, err := NewBuilder(config, registry, cache)
//	result, err := builder.Build(ctx)

package builder