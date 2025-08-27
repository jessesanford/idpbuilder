package contexts

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

// ContextManager manages build context creation and lifecycle
type ContextManager struct {
	workDir string
}

// NewContextManager creates a new context manager
func NewContextManager(workDir string) (*ContextManager, error) {
	if workDir == "" {
		workDir = "/tmp/buildah-contexts"
	}
	os.MkdirAll(workDir, 0755)
	return &ContextManager{workDir: workDir}, nil
}

// PrepareContext prepares a build context based on options
func (cm *ContextManager) PrepareContext(ctx context.Context, opts ContextOptions) (Context, error) {
	if opts.Type != Local {
		return nil, fmt.Errorf("unsupported context type: %s", opts.Type)
	}
	return cm.createLocalContext(opts)
}

// createLocalContext creates a local directory context
func (cm *ContextManager) createLocalContext(opts ContextOptions) (Context, error) {
	source, err := filepath.Abs(opts.Source)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve source path: %w", err)
	}
	
	if _, err := os.Stat(source); err != nil {
		return nil, fmt.Errorf("source path does not exist: %s", source)
	}
	
	return &LocalContext{
		sourcePath: source,
		excludes:   opts.Excludes,
	}, nil
}

// Cleanup removes all managed contexts and their resources
func (cm *ContextManager) Cleanup() error {
	// Simple cleanup for this split - no managed contexts to track
	return nil
}