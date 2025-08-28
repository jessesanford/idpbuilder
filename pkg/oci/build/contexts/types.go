package contexts

import (
	"net/http"
	"time"
)

// Context defines the interface for all build context types
type Context interface {
	Path() string
	Cleanup() error
	Type() ContextType
}

// ContextType represents the type of build context
type ContextType int

const (
	LocalContext ContextType = iota
	URLContext
	GitContext
	ArchiveContext
)

// String returns the string representation of ContextType
func (ct ContextType) String() string {
	switch ct {
	case LocalContext:
		return "local"
	case URLContext:
		return "url"
	case GitContext:
		return "git"
	case ArchiveContext:
		return "archive"
	default:
		return "unknown"
	}
}

// ContextConfig holds configuration for context resolution
type ContextConfig struct {
	MaxSize      int64
	CacheEnabled bool
	TempDir      string
	HTTPTimeout  time.Duration
}

// DefaultConfig returns a default configuration
func DefaultConfig() *ContextConfig {
	return &ContextConfig{
		MaxSize:      500 * 1024 * 1024, // 500MB
		CacheEnabled: true,
		TempDir:      "/tmp",
		HTTPTimeout:  30 * time.Second,
	}
}

// AuthConfig holds authentication configuration for git/http
type AuthConfig struct {
	Username string
	Password string
	Token    string
}

// HTTPClient interface for mocking
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}