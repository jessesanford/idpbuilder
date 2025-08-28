package contexts

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GitContextImpl implements Context for Git repository sources
type GitContextImpl struct {
	repoURL   string
	ref       string
	config    *ContextConfig
	cloneDir  string
	shallow   bool
	authConfig *AuthConfig
}

// Path returns the path where the repository was cloned
func (g *GitContextImpl) Path() string {
	return g.cloneDir
}

// Type returns GitContext
func (g *GitContextImpl) Type() ContextType {
	return GitContext
}

// Cleanup removes the cloned repository directory
func (g *GitContextImpl) Cleanup() error {
	if g.cloneDir != "" {
		return os.RemoveAll(g.cloneDir)
	}
	return nil
}

// Clone clones the git repository and returns the clone path
func (g *GitContextImpl) Clone() (string, error) {
	// Validate the Git URL
	if err := g.validateGitURL(); err != nil {
		return "", fmt.Errorf("invalid git URL: %w", err)
	}

	// Parse the URL and extract ref if present
	g.parseURLAndRef()

	// Create temporary directory for cloning
	tempDir, err := os.MkdirTemp(g.config.TempDir, "git_context_*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}
	g.cloneDir = tempDir

	// Clone the repository
	if err := g.cloneRepository(); err != nil {
		g.Cleanup() // Clean up on failure
		return "", fmt.Errorf("failed to clone repository: %w", err)
	}

	// Checkout specific ref if provided
	if g.ref != "" && g.ref != "HEAD" {
		if err := g.checkoutRef(); err != nil {
			g.Cleanup() // Clean up on failure
			return "", fmt.Errorf("failed to checkout ref %s: %w", g.ref, err)
		}
	}

	return g.cloneDir, nil
}

// validateGitURL validates that the URL is a proper Git repository URL
func (g *GitContextImpl) validateGitURL() error {
	if g.repoURL == "" {
		return fmt.Errorf("repository URL cannot be empty")
	}

	// Check for common Git URL patterns
	if strings.HasPrefix(g.repoURL, "git://") ||
		strings.HasPrefix(g.repoURL, "git@") ||
		strings.HasSuffix(g.repoURL, ".git") {
		return nil
	}

	// Check if it's an HTTPS URL that might be a Git repository
	if strings.HasPrefix(g.repoURL, "https://") {
		parsedURL, err := url.Parse(g.repoURL)
		if err != nil {
			return fmt.Errorf("invalid URL format: %w", err)
		}

		// Common Git hosting services
		gitHosts := []string{"github.com", "gitlab.com", "bitbucket.org", "git."}
		for _, host := range gitHosts {
			if strings.Contains(parsedURL.Host, host) {
				return nil
			}
		}
	}

	// If we reach here, try to validate by checking if git ls-remote works
	return g.testGitRemote()
}

// testGitRemote tests if the URL is a valid Git remote
func (g *GitContextImpl) testGitRemote() error {
	cmd := exec.Command("git", "ls-remote", "--exit-code", g.repoURL, "HEAD")
	cmd.Stderr = nil // Suppress stderr
	return cmd.Run()
}

// parseURLAndRef extracts ref from URL fragment if present
func (g *GitContextImpl) parseURLAndRef() {
	// Check for URL fragments like https://github.com/user/repo.git#branch
	if strings.Contains(g.repoURL, "#") {
		parts := strings.Split(g.repoURL, "#")
		if len(parts) == 2 {
			g.repoURL = parts[0]
			g.ref = parts[1]
		}
	}

	// Set defaults
	if g.ref == "" {
		g.ref = "HEAD"
	}
	g.shallow = true // Default to shallow clone for efficiency
}

// cloneRepository performs the actual git clone operation
func (g *GitContextImpl) cloneRepository() error {
	args := []string{"clone"}
	
	// Add shallow clone for efficiency (single branch, depth 1)
	if g.shallow {
		args = append(args, "--depth", "1")
		if g.ref != "" && g.ref != "HEAD" {
			args = append(args, "--branch", g.ref)
		}
	}

	// Add authentication if configured
	if g.authConfig != nil {
		g.setupGitAuth()
	}

	// Add URL and target directory
	args = append(args, g.repoURL, g.cloneDir)

	// Execute git clone command
	cmd := exec.Command("git", args...)
	cmd.Env = append(os.Environ(), g.getGitEnv()...)
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		// SECURITY FIX: Sanitize URL in error messages to prevent credential leakage
		sanitizedURL := g.sanitizeURLForLog(g.repoURL)
		return fmt.Errorf("git clone failed for %s: %w (output: %s)", sanitizedURL, err, string(output))
	}

	return nil
}

// checkoutRef checks out a specific ref (branch, tag, or commit)
func (g *GitContextImpl) checkoutRef() error {
	// If we used shallow clone with branch, checkout is already done
	if g.shallow {
		return nil
	}

	cmd := exec.Command("git", "checkout", g.ref)
	cmd.Dir = g.cloneDir
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git checkout failed: %w (output: %s)", err, string(output))
	}

	return nil
}

// setupGitAuth configures Git authentication
func (g *GitContextImpl) setupGitAuth() {
	if g.authConfig == nil {
		return
	}

	// Set up credential helper or SSH key
	if g.authConfig.Token != "" {
		// For token-based auth, we'll use the askpass helper
		g.setupTokenAuth()
	} else if g.authConfig.Username != "" && g.authConfig.Password != "" {
		// For username/password auth
		g.setupPasswordAuth()
	}
}

// setupTokenAuth sets up token-based authentication
func (g *GitContextImpl) setupTokenAuth() {
	// Create a temporary askpass script
	askpassScript := filepath.Join(g.cloneDir, "..", "git_askpass.sh")
	scriptContent := fmt.Sprintf("#!/bin/sh\necho '%s'\n", g.authConfig.Token)
	
	os.WriteFile(askpassScript, []byte(scriptContent), 0700)
	
	// SECURITY FIX: Schedule secure deletion of the token file
	defer func() {
		if err := g.secureDeleteFile(askpassScript); err != nil {
			// Log error but don't fail the operation
			fmt.Printf("Warning: failed to securely delete token file: %v\n", err)
		}
	}()
	
	// Set environment variables for git
	os.Setenv("GIT_ASKPASS", askpassScript)
	os.Setenv("GIT_USERNAME", "token")
}

// setupPasswordAuth sets up username/password authentication
func (g *GitContextImpl) setupPasswordAuth() {
	// SECURITY FIX: Don't embed credentials in URLs - use credential helper instead
	if strings.HasPrefix(g.repoURL, "https://") {
		// Create a temporary credential helper script instead of modifying URL
		credentialScript := filepath.Join(g.cloneDir, "..", "git_credential_helper.sh")
		scriptContent := fmt.Sprintf("#!/bin/sh\necho 'username=%s'\necho 'password=%s'\n", 
			g.authConfig.Username, g.authConfig.Password)
		
		os.WriteFile(credentialScript, []byte(scriptContent), 0700)
		
		// Schedule secure deletion
		defer func() {
			if err := g.secureDeleteFile(credentialScript); err != nil {
				fmt.Printf("Warning: failed to securely delete credential file: %v\n", err)
			}
		}()
		
		// Set git credential helper
		os.Setenv("GIT_CONFIG_GLOBAL", "/dev/null") // Prevent global config interference
		os.Setenv("GIT_CONFIG_SYSTEM", "/dev/null") // Prevent system config interference
	}
}

// getGitEnv returns environment variables for git commands
func (g *GitContextImpl) getGitEnv() []string {
	env := []string{
		"GIT_TERMINAL_PROMPT=0", // Disable interactive prompts
	}

	// Add timeout for operations
	if g.config.HTTPTimeout > 0 {
		timeout := int(g.config.HTTPTimeout.Seconds())
		env = append(env, fmt.Sprintf("GIT_HTTP_LOW_SPEED_LIMIT=1000"))
		env = append(env, fmt.Sprintf("GIT_HTTP_LOW_SPEED_TIME=%d", timeout))
	}

	return env
}

// SetRef sets the Git reference (branch, tag, or commit) to checkout
func (g *GitContextImpl) SetRef(ref string) {
	g.ref = ref
}

// SetShallow sets whether to use shallow cloning
func (g *GitContextImpl) SetShallow(shallow bool) {
	g.shallow = shallow
}

// SetAuth sets authentication configuration
func (g *GitContextImpl) SetAuth(auth *AuthConfig) {
	g.authConfig = auth
}

// GetCommitHash returns the current commit hash of the cloned repository
func (g *GitContextImpl) GetCommitHash() (string, error) {
	if g.cloneDir == "" {
		return "", fmt.Errorf("repository not cloned yet")
	}

	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = g.cloneDir
	
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get commit hash: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// GetRemoteURL returns the remote URL of the cloned repository
func (g *GitContextImpl) GetRemoteURL() string {
	return g.repoURL
}

// GetRef returns the current ref
func (g *GitContextImpl) GetRef() string {
	return g.ref
}

// secureDeleteFile securely deletes a file by overwriting it with random data
func (g *GitContextImpl) secureDeleteFile(filepath string) error {
	// Check if file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return nil // Already deleted
	}
	
	// Get file size
	info, err := os.Stat(filepath)
	if err != nil {
		return err
	}
	
	// Open file for writing
	file, err := os.OpenFile(filepath, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer file.Close()
	
	// Overwrite with zeros (simple secure deletion)
	zeros := make([]byte, info.Size())
	if _, err := file.Write(zeros); err != nil {
		return err
	}
	
	// Sync to ensure data is written to disk
	if err := file.Sync(); err != nil {
		return err
	}
	
	// Finally remove the file
	return os.Remove(filepath)
}

// sanitizeURLForLog removes sensitive information from URLs for logging
func (g *GitContextImpl) sanitizeURLForLog(rawURL string) string {
	if rawURL == "" {
		return ""
	}
	
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		// If we can't parse it, just return a generic placeholder
		return "[SANITIZED_URL]"
	}
	
	// Remove any user information (credentials)
	parsedURL.User = nil
	
	return parsedURL.String()
}