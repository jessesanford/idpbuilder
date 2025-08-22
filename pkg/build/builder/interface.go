// Package builder defines the interface for container build operations
package builder

import (
	"context"
	
	"github.com/vscode/workspaces/idpbuilder/pkg/build/api"
)

// Builder defines the interface for container build operations
type Builder interface {
	// BuildAndPush builds a container image and pushes to registry
	// This is the only method needed for MVP
	BuildAndPush(ctx context.Context, req api.BuildRequest) (*api.BuildResponse, error)
}