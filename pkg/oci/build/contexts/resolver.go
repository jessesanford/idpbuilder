package contexts

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// ContextResolver handles resolution of different context types
type ContextResolver struct {
	config     *ContextConfig
	httpClient HTTPClient
	cleanupFns []func() error
}

// NewContextResolver creates a new context resolver with given config
func NewContextResolver(config *ContextConfig) *ContextResolver {
	if config == nil {
		config = DefaultConfig()
	}

	return &ContextResolver{
		config: config,
		httpClient: &http.Client{
			Timeout: config.HTTPTimeout,
		},
		cleanupFns: make([]func() error, 0),
	}
}

// ResolveContext resolves a context source to an appropriate Context implementation
func (r *ContextResolver) ResolveContext(source string) (Context, error) {
	if source == "" {
		return nil, fmt.Errorf("context source cannot be empty")
	}

	contextType := r.detectContextType(source)
	
	switch contextType {
	case LocalContext:
		return r.createLocalContext(source)
	case URLContext:
		return r.createURLContext(source)
	case GitContext:
		return nil, fmt.Errorf("git context support requires additional components (split 002)")
	case ArchiveContext:
		return nil, fmt.Errorf("archive context support requires additional components (split 002)")
	default:
		return nil, fmt.Errorf("unknown context type for source: %s", source)
	}
}

// detectContextType determines the context type based on the source string
func (r *ContextResolver) detectContextType(source string) ContextType {
	// Check if it's a URL
	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		return URLContext
	}

	// Check if it's a git repository
	if strings.HasPrefix(source, "git://") || 
	   strings.HasPrefix(source, "git@") ||
	   (strings.HasPrefix(source, "https://") && strings.Contains(source, ".git")) ||
	   strings.HasSuffix(source, ".git") {
		return GitContext
	}

	// Check if it's an archive file
	if r.isArchiveFile(source) {
		return ArchiveContext
	}

	// Default to local context
	return LocalContext
}

// isArchiveFile checks if the source appears to be an archive file
func (r *ContextResolver) isArchiveFile(source string) bool {
	archiveExts := []string{".tar", ".tar.gz", ".tgz", ".tar.bz2", ".tbz2", ".zip"}
	lower := strings.ToLower(source)
	
	for _, ext := range archiveExts {
		if strings.HasSuffix(lower, ext) {
			return true
		}
	}
	return false
}

// createLocalContext creates a local filesystem context
func (r *ContextResolver) createLocalContext(source string) (Context, error) {
	// Convert to absolute path
	absPath, err := filepath.Abs(source)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path for %s: %w", source, err)
	}

	// Verify path exists and is a directory
	info, err := os.Stat(absPath)
	if err != nil {
		return nil, fmt.Errorf("context path does not exist: %s", absPath)
	}
	
	if !info.IsDir() {
		return nil, fmt.Errorf("context path is not a directory: %s", absPath)
	}

	return &LocalContextImpl{path: absPath}, nil
}

// createURLContext creates a URL-based context
func (r *ContextResolver) createURLContext(source string) (Context, error) {
	parsedURL, err := url.Parse(source)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %s", source)
	}

	ctx := &URLContextImpl{
		url:        parsedURL,
		config:     r.config,
		httpClient: r.httpClient,
	}

	// Register cleanup function
	r.addCleanupFn(ctx.Cleanup)
	
	return ctx, nil
}


// addCleanupFn registers a cleanup function to be called later
func (r *ContextResolver) addCleanupFn(fn func() error) {
	r.cleanupFns = append(r.cleanupFns, fn)
}

// Cleanup cleans up all managed contexts
func (r *ContextResolver) Cleanup() error {
	var errs []string
	
	for _, cleanupFn := range r.cleanupFns {
		if err := cleanupFn(); err != nil {
			errs = append(errs, err.Error())
		}
	}
	
	if len(errs) > 0 {
		return fmt.Errorf("cleanup errors: %s", strings.Join(errs, "; "))
	}
	
	return nil
}

// LocalContextImpl implements Context for local filesystem contexts
type LocalContextImpl struct {
	path string
}

func (l *LocalContextImpl) Path() string {
	return l.path
}

func (l *LocalContextImpl) Cleanup() error {
	// Local contexts don't need cleanup
	return nil
}

func (l *LocalContextImpl) Type() ContextType {
	return LocalContext
}