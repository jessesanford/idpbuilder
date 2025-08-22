// Package api defines the interface for container build operations
package api

import "context"

// Builder defines the interface for container build operations
type Builder interface {
	// BuildAndPush builds a container image and pushes to registry
	// This is the only method needed for MVP
	BuildAndPush(ctx context.Context, req BuildRequest) (*BuildResponse, error)
}

// BuilderConfig holds configuration for builder instances
type BuilderConfig struct {
	// Registry is the target OCI registry
	Registry string

	// Namespace is the registry namespace (hardcoded for MVP)
	Namespace string

	// InsecureSkipTLSVerify skips TLS verification (default true for MVP)
	InsecureSkipTLSVerify bool
}

// DefaultConfig returns MVP configuration
func DefaultConfig() BuilderConfig {
	return BuilderConfig{
		Registry:              "gitea.cnoe.localtest.me",
		Namespace:             "giteaadmin",
		InsecureSkipTLSVerify: true,
	}
}