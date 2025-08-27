package contexts

import (
	"io"
	"os"
	"time"
)

// ContextType represents the type of context
type ContextType int

const (
	LocalContextType ContextType = iota
	RemoteContextType
	ArchiveContextType
)

// Context interface defines the common operations for all context types
type Context interface {
	Type() ContextType
	PrepareContext(options *ContextOptions) error
	GetEntries() ([]ContextEntry, error)
	GetEntry(path string) (io.ReadCloser, error)
	Cleanup() error
}

// ContextOptions provides configuration for context preparation
type ContextOptions struct {
	// Additional options can be added as needed
	IgnorePatterns []string
	MaxDepth       int
}

// ContextEntry represents a file or directory entry in the context
type ContextEntry struct {
	Path    string
	Size    int64
	Mode    os.FileMode
	ModTime time.Time
	IsDir   bool
}