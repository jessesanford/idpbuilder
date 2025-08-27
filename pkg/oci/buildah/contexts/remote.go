package contexts

import (
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// RemoteContext implements Context interface for remote git repositories
type RemoteContext struct {
	url           string
	ref           string
	tempDir       string
	authToken     string
	prepared      bool
	cloneTimeout  time.Duration
}

// NewRemoteContext creates a new remote git context
func NewRemoteContext(gitURL, ref string) (*RemoteContext, error) {
	if err := validateGitURL(gitURL); err != nil {
		return nil, fmt.Errorf("invalid git URL: %w", err)
	}

	return &RemoteContext{
		url:          gitURL,
		ref:          ref,
		authToken:    os.Getenv("GIT_TOKEN"),
		cloneTimeout: 5 * time.Minute,
	}, nil
}

// Type returns the context type
func (r *RemoteContext) Type() ContextType {
	return RemoteContextType
}

// PrepareContext clones the repository and checks out the specified ref
func (r *RemoteContext) PrepareContext(options *ContextOptions) error {
	if r.prepared {
		return nil
	}

	tempDir, err := os.MkdirTemp("", "buildah-context-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	r.tempDir = tempDir

	// Clone repository
	if err := r.cloneRepository(); err != nil {
		r.Cleanup()
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	// Checkout specific ref if provided
	if r.ref != "" {
		if err := r.checkoutRef(); err != nil {
			r.Cleanup()
			return fmt.Errorf("failed to checkout ref %s: %w", r.ref, err)
		}
	}

	r.prepared = true
	return nil
}

// GetEntries returns all entries in the cloned repository
func (r *RemoteContext) GetEntries() ([]ContextEntry, error) {
	if !r.prepared {
		return nil, fmt.Errorf("context not prepared")
	}

	var entries []ContextEntry

	err := filepath.WalkDir(r.tempDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip .git directory
		if d.IsDir() && d.Name() == ".git" {
			return filepath.SkipDir
		}

		relPath, err := filepath.Rel(r.tempDir, path)
		if err != nil {
			return err
		}

		// Skip root directory
		if relPath == "." {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		entries = append(entries, ContextEntry{
			Path:    relPath,
			Size:    info.Size(),
			Mode:    info.Mode(),
			ModTime: info.ModTime(),
			IsDir:   d.IsDir(),
		})

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	return entries, nil
}

// GetEntry returns a reader for the specified file
func (r *RemoteContext) GetEntry(path string) (io.ReadCloser, error) {
	if !r.prepared {
		return nil, fmt.Errorf("context not prepared")
	}

	fullPath := filepath.Join(r.tempDir, path)
	if !strings.HasPrefix(fullPath, r.tempDir) {
		return nil, fmt.Errorf("path outside context boundary: %s", path)
	}

	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", path, err)
	}

	return file, nil
}

// Cleanup removes the temporary directory
func (r *RemoteContext) Cleanup() error {
	if r.tempDir == "" {
		return nil
	}

	err := os.RemoveAll(r.tempDir)
	r.tempDir = ""
	r.prepared = false
	return err
}

// cloneRepository clones the git repository to temp directory
func (r *RemoteContext) cloneRepository() error {
	args := []string{"clone"}

	// Add authentication if token is available
	if r.authToken != "" && (strings.HasPrefix(r.url, "https://") || strings.HasPrefix(r.url, "http://")) {
		// Insert token into URL for HTTPS authentication
		authURL, err := r.addTokenToURL()
		if err != nil {
			return err
		}
		args = append(args, authURL)
	} else {
		args = append(args, r.url)
	}

	args = append(args, r.tempDir)

	cmd := exec.Command("git", args...)
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git clone failed: %w, output: %s", err, string(output))
	}

	return nil
}

// checkoutRef checks out the specified git reference
func (r *RemoteContext) checkoutRef() error {
	cmd := exec.Command("git", "checkout", r.ref)
	cmd.Dir = r.tempDir

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git checkout failed: %w, output: %s", err, string(output))
	}

	return nil
}

// addTokenToURL adds authentication token to HTTPS URLs
func (r *RemoteContext) addTokenToURL() (string, error) {
	u, err := url.Parse(r.url)
	if err != nil {
		return "", err
	}

	u.User = url.User(r.authToken)
	return u.String(), nil
}

// validateGitURL validates that the URL is a supported git URL
func validateGitURL(gitURL string) error {
	if gitURL == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	// Parse URL
	u, err := url.Parse(gitURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	// Check supported schemes
	switch u.Scheme {
	case "https", "http", "ssh", "git":
		return nil
	case "":
		// Check for SSH format like git@github.com:user/repo.git
		if strings.Contains(gitURL, "@") && strings.Contains(gitURL, ":") {
			return nil
		}
		return fmt.Errorf("missing scheme, supported schemes: https, http, ssh, git")
	default:
		return fmt.Errorf("unsupported scheme %s, supported schemes: https, http, ssh, git", u.Scheme)
	}
}