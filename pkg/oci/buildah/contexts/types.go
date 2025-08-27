package contexts

import (
	"context"
	"io"
	"time"
)

// ContextType represents the type of build context
type ContextType string

const (
	Local   ContextType = "local"
	Remote  ContextType = "remote"
	Archive ContextType = "archive"
)

// ContextOptions holds configuration for build context preparation
type ContextOptions struct {
	Type     ContextType `json:"type"`
	Source   string      `json:"source"`
	Excludes []string    `json:"excludes,omitempty"`
}

// ContextEntry represents a single file or directory in the build context
type ContextEntry struct {
	Path    string    `json:"path"`
	Mode    uint32    `json:"mode"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"modTime"`
	IsDir   bool      `json:"isDir"`
	Content io.Reader `json:"-"`
}

// Context interface defines the contract for build contexts
type Context interface {
	PrepareContext(ctx context.Context, opts ContextOptions) (string, error)
	GetEntries() ([]*ContextEntry, error)
	GetEntry(path string) (*ContextEntry, error)
	Cleanup() error
	Type() ContextType
}