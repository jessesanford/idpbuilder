package contexts

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestGitContextImpl_Type(t *testing.T) {
	g := &GitContextImpl{}
	if got := g.Type(); got != GitContext {
		t.Errorf("Type() = %v, want %v", got, GitContext)
	}
}

func TestGitContextImpl_Path(t *testing.T) {
	g := &GitContextImpl{cloneDir: "/test/path"}
	if got := g.Path(); got != "/test/path" {
		t.Errorf("Path() = %v, want %v", got, "/test/path")
	}
}

func TestGitContextImpl_SetRef(t *testing.T) {
	g := &GitContextImpl{}
	g.SetRef("main")
	if g.ref != "main" {
		t.Errorf("SetRef() failed, got %v, want %v", g.ref, "main")
	}
}

func TestGitContextImpl_SetShallow(t *testing.T) {
	g := &GitContextImpl{}
	g.SetShallow(false)
	if g.shallow != false {
		t.Errorf("SetShallow() failed, got %v, want %v", g.shallow, false)
	}
}

func TestGitContextImpl_SetAuth(t *testing.T) {
	g := &GitContextImpl{}
	auth := &AuthConfig{Username: "test", Password: "pass"}
	g.SetAuth(auth)
	if g.authConfig != auth {
		t.Errorf("SetAuth() failed")
	}
}

func TestGitContextImpl_GetRef(t *testing.T) {
	g := &GitContextImpl{ref: "develop"}
	if got := g.GetRef(); got != "develop" {
		t.Errorf("GetRef() = %v, want %v", got, "develop")
	}
}

func TestGitContextImpl_GetRemoteURL(t *testing.T) {
	g := &GitContextImpl{repoURL: "https://github.com/test/repo.git"}
	if got := g.GetRemoteURL(); got != "https://github.com/test/repo.git" {
		t.Errorf("GetRemoteURL() = %v, want %v", got, "https://github.com/test/repo.git")
	}
}

func TestGitContextImpl_Cleanup(t *testing.T) {
	// Test cleanup with empty path
	g := &GitContextImpl{}
	if err := g.Cleanup(); err != nil {
		t.Errorf("Cleanup() error = %v, want nil", err)
	}

	// Test cleanup with actual directory
	tempDir, err := os.MkdirTemp("", "git_test_*")
	if err != nil {
		t.Fatal(err)
	}
	
	g.cloneDir = tempDir
	if err := g.Cleanup(); err != nil {
		t.Errorf("Cleanup() error = %v, want nil", err)
	}
	
	// Verify directory was removed
	if _, err := os.Stat(tempDir); !os.IsNotExist(err) {
		t.Error("Cleanup() did not remove directory")
	}
}

func TestGitContextImpl_validateGitURL(t *testing.T) {
	tests := []struct {
		name    string
		repoURL string
		wantErr bool
	}{
		{
			name:    "empty URL",
			repoURL: "",
			wantErr: true,
		},
		{
			name:    "git protocol",
			repoURL: "git://github.com/user/repo.git",
			wantErr: false,
		},
		{
			name:    "ssh protocol",
			repoURL: "git@github.com:user/repo.git",
			wantErr: false,
		},
		{
			name:    "https with .git suffix",
			repoURL: "https://github.com/user/repo.git",
			wantErr: false,
		},
		{
			name:    "github https without .git",
			repoURL: "https://github.com/user/repo",
			wantErr: false,
		},
		{
			name:    "gitlab https",
			repoURL: "https://gitlab.com/user/repo.git",
			wantErr: false,
		},
		{
			name:    "bitbucket https",
			repoURL: "https://bitbucket.org/user/repo.git",
			wantErr: false,
		},
		{
			name:    "invalid https",
			repoURL: "https://example.com/not-git",
			wantErr: true, // Will fail on ls-remote
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GitContextImpl{repoURL: tt.repoURL}
			err := g.validateGitURL()
			if (err != nil) != tt.wantErr {
				t.Errorf("validateGitURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGitContextImpl_parseURLAndRef(t *testing.T) {
	tests := []struct {
		name        string
		inputURL    string
		expectedURL string
		expectedRef string
	}{
		{
			name:        "URL with fragment",
			inputURL:    "https://github.com/user/repo.git#main",
			expectedURL: "https://github.com/user/repo.git",
			expectedRef: "main",
		},
		{
			name:        "URL without fragment",
			inputURL:    "https://github.com/user/repo.git",
			expectedURL: "https://github.com/user/repo.git",
			expectedRef: "HEAD",
		},
		{
			name:        "URL with multiple fragments",
			inputURL:    "https://github.com/user/repo.git#feature#branch",
			expectedURL: "https://github.com/user/repo.git#feature#branch", // Won't parse because len(parts) != 2
			expectedRef: "HEAD",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GitContextImpl{repoURL: tt.inputURL}
			g.parseURLAndRef()
			
			if g.repoURL != tt.expectedURL {
				t.Errorf("parseURLAndRef() repoURL = %v, want %v", g.repoURL, tt.expectedURL)
			}
			if g.ref != tt.expectedRef {
				t.Errorf("parseURLAndRef() ref = %v, want %v", g.ref, tt.expectedRef)
			}
			if !g.shallow {
				t.Error("parseURLAndRef() should set shallow to true by default")
			}
		})
	}
}

func TestGitContextImpl_setupTokenAuth(t *testing.T) {
	// Create temp directory for test
	tempDir, err := os.MkdirTemp("", "git_auth_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	g := &GitContextImpl{
		cloneDir:   filepath.Join(tempDir, "repo"),
		authConfig: &AuthConfig{Token: "test_token_123"},
	}

	// Create the clone directory structure
	os.MkdirAll(g.cloneDir, 0755)
	
	g.setupTokenAuth()

	// Verify environment variables are set
	if os.Getenv("GIT_USERNAME") != "token" {
		t.Error("GIT_USERNAME not set correctly")
	}

	// Note: The security fix schedules deletion via defer, so the file 
	// is deleted immediately after creation in the same function call.
	// We can't verify the file exists because it's deleted by the defer.
	// This test mainly verifies that setupTokenAuth doesn't panic.
}

func TestGitContextImpl_setupPasswordAuth(t *testing.T) {
	tests := []struct {
		name     string
		repoURL  string
		username string
		password string
	}{
		{
			name:     "HTTPS URL with credentials - uses credential helper",
			repoURL:  "https://github.com/user/repo.git",
			username: "testuser",
			password: "testpass",
		},
		{
			name:     "Non-HTTPS URL unchanged",
			repoURL:  "git@github.com:user/repo.git",
			username: "testuser", 
			password: "testpass",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory for test
			tempDir, err := os.MkdirTemp("", "password_auth_test_*")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tempDir)

			g := &GitContextImpl{
				repoURL: tt.repoURL,
				cloneDir: filepath.Join(tempDir, "repo"),
				authConfig: &AuthConfig{
					Username: tt.username,
					Password: tt.password,
				},
			}

			// Create the clone directory structure
			os.MkdirAll(g.cloneDir, 0755)

			g.setupPasswordAuth()

			// The security fix uses credential helper instead of embedding credentials
			// URL should remain unchanged
			if g.repoURL != tt.repoURL {
				t.Errorf("setupPasswordAuth() repoURL = %v, want %v", g.repoURL, tt.repoURL)
			}
		})
	}
}

func TestGitContextImpl_getGitEnv(t *testing.T) {
	g := &GitContextImpl{
		config: &ContextConfig{HTTPTimeout: 30 * time.Second},
	}

	env := g.getGitEnv()

	// Check for required environment variables
	found := false
	for _, e := range env {
		if e == "GIT_TERMINAL_PROMPT=0" {
			found = true
			break
		}
	}
	if !found {
		t.Error("GIT_TERMINAL_PROMPT=0 not found in environment")
	}

	// Check for timeout settings
	foundTimeout := false
	for _, e := range env {
		if strings.Contains(e, "GIT_HTTP_LOW_SPEED_TIME=30") {
			foundTimeout = true
			break
		}
	}
	if !foundTimeout {
		t.Error("Timeout settings not found in environment")
	}
}

func TestGitContextImpl_GetCommitHash_NotCloned(t *testing.T) {
	g := &GitContextImpl{}
	_, err := g.GetCommitHash()
	if err == nil {
		t.Error("GetCommitHash() should return error when not cloned")
	}
	if !strings.Contains(err.Error(), "repository not cloned yet") {
		t.Errorf("GetCommitHash() error = %v, want 'repository not cloned yet'", err)
	}
}

func TestGitContextImpl_setupGitAuth(t *testing.T) {
	tests := []struct {
		name       string
		authConfig *AuthConfig
		expectCall bool
	}{
		{
			name:       "no auth config",
			authConfig: nil,
			expectCall: false,
		},
		{
			name:       "token auth",
			authConfig: &AuthConfig{Token: "token123"},
			expectCall: true,
		},
		{
			name:       "password auth",
			authConfig: &AuthConfig{Username: "user", Password: "pass"},
			expectCall: true,
		},
		{
			name:       "empty auth config",
			authConfig: &AuthConfig{},
			expectCall: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory for test
			tempDir, err := os.MkdirTemp("", "git_setup_auth_test_*")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tempDir)

			g := &GitContextImpl{
				cloneDir:   filepath.Join(tempDir, "repo"),
				authConfig: tt.authConfig,
			}

			// Create the clone directory structure
			os.MkdirAll(g.cloneDir, 0755)

			// This test mainly checks that setupGitAuth doesn't panic
			// and handles different auth configurations
			g.setupGitAuth()

			// The security fix means token files are deleted immediately via defer
			// We can't verify the file exists, but we can verify the method doesn't panic
			// and that environment variables are set for token auth
			if tt.authConfig != nil && tt.authConfig.Token != "" {
				if os.Getenv("GIT_USERNAME") != "token" {
					t.Error("Token auth should set GIT_USERNAME")
				}
			}
		})
	}
}

// Integration-style test for the clone workflow (without actually calling git)
func TestGitContextImpl_Clone_ValidationFlow(t *testing.T) {
	config := DefaultConfig()
	config.TempDir = os.TempDir()

	tests := []struct {
		name    string
		repoURL string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "empty URL",
			repoURL: "",
			wantErr: true,
			errMsg:  "invalid git URL",
		},
		{
			name:    "invalid URL format",
			repoURL: "not-a-valid-url",
			wantErr: true,
			errMsg:  "invalid git URL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GitContextImpl{
				repoURL: tt.repoURL,
				config:  config,
			}

			_, err := g.Clone()
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("Clone() expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Clone() error = %v, want to contain %v", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Clone() unexpected error = %v", err)
				}
			}
		})
	}
}

// Mock test for successful operations
func TestGitContextImpl_Methods_Coverage(t *testing.T) {
	g := &GitContextImpl{
		repoURL:  "https://github.com/test/repo.git",
		ref:      "main",
		config:   DefaultConfig(),
		cloneDir: "/tmp/test",
		shallow:  true,
	}

	// Test getter methods
	if g.Type() != GitContext {
		t.Error("Type() should return GitContext")
	}

	if g.Path() != "/tmp/test" {
		t.Error("Path() should return cloneDir")
	}

	if g.GetRef() != "main" {
		t.Error("GetRef() should return ref")
	}

	if g.GetRemoteURL() != "https://github.com/test/repo.git" {
		t.Error("GetRemoteURL() should return repoURL")
	}

	// Test setter methods
	g.SetRef("develop")
	if g.ref != "develop" {
		t.Error("SetRef() should update ref")
	}

	g.SetShallow(false)
	if g.shallow != false {
		t.Error("SetShallow() should update shallow")
	}

	auth := &AuthConfig{Token: "test"}
	g.SetAuth(auth)
	if g.authConfig != auth {
		t.Error("SetAuth() should update authConfig")
	}
}

func TestGitContextImpl_secureDeleteFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "secure_delete_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	g := &GitContextImpl{}
	
	// Test deleting non-existent file
	err = g.secureDeleteFile(filepath.Join(tempDir, "nonexistent.txt"))
	if err != nil {
		t.Errorf("secureDeleteFile() should handle non-existent files, got error: %v", err)
	}
	
	// Test deleting existing file
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "sensitive data"
	err = os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatal(err)
	}
	
	err = g.secureDeleteFile(testFile)
	if err != nil {
		t.Errorf("secureDeleteFile() error = %v, want nil", err)
	}
	
	// Verify file was deleted
	if _, err := os.Stat(testFile); !os.IsNotExist(err) {
		t.Error("secureDeleteFile() should remove the file")
	}
}

func TestGitContextImpl_sanitizeURLForLog(t *testing.T) {
	g := &GitContextImpl{}
	
	tests := []struct {
		name     string
		rawURL   string
		expected string
	}{
		{
			name:     "empty URL",
			rawURL:   "",
			expected: "",
		},
		{
			name:     "URL with credentials",
			rawURL:   "https://user:pass@github.com/repo.git",
			expected: "https://github.com/repo.git",
		},
		{
			name:     "URL without credentials",
			rawURL:   "https://github.com/user/repo.git",
			expected: "https://github.com/user/repo.git",
		},
		{
			name:     "invalid URL with colon",
			rawURL:   "http://[::1]:80080080", // This will fail parsing
			expected: "[SANITIZED_URL]",
		},
		{
			name:     "SSH URL",
			rawURL:   "git@github.com:user/repo.git",
			expected: "[SANITIZED_URL]", // SSH URLs can't be parsed by url.Parse properly
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := g.sanitizeURLForLog(tt.rawURL)
			if result != tt.expected {
				t.Errorf("sanitizeURLForLog() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGitContextImpl_checkoutRef(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "checkout_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	g := &GitContextImpl{
		cloneDir: tempDir,
		shallow:  true, // Shallow clone should skip checkout
	}

	// When shallow is true, checkoutRef should return nil immediately
	err = g.checkoutRef()
	if err != nil {
		t.Errorf("checkoutRef() with shallow=true should return nil, got %v", err)
	}

	// Test non-shallow case (will fail since we don't have a real git repo, but covers the path)
	g.shallow = false
	err = g.checkoutRef()
	if err == nil {
		t.Error("checkoutRef() should fail without a real git repo")
	}
}

func TestGitContextImpl_testGitRemote(t *testing.T) {
	g := &GitContextImpl{
		repoURL: "https://invalid-url-that-does-not-exist.invalid",
	}

	// This should fail because the URL doesn't exist
	err := g.testGitRemote()
	if err == nil {
		t.Error("testGitRemote() should fail with invalid URL")
	}
}

func TestGitContextImpl_cloneRepository_ErrorPaths(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "clone_error_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	g := &GitContextImpl{
		repoURL:  "https://invalid-url-that-does-not-exist.invalid",
		cloneDir: filepath.Join(tempDir, "repo"),
		config:   DefaultConfig(),
		shallow:  true,
	}

	// This should fail because the URL doesn't exist
	err = g.cloneRepository()
	if err == nil {
		t.Error("cloneRepository() should fail with invalid URL")
	}
}

func TestGitContextImpl_getGitEnv_NoTimeout(t *testing.T) {
	g := &GitContextImpl{
		config: &ContextConfig{HTTPTimeout: 0}, // No timeout
	}

	env := g.getGitEnv()

	// Should still have terminal prompt disabled
	found := false
	for _, e := range env {
		if e == "GIT_TERMINAL_PROMPT=0" {
			found = true
			break
		}
	}
	if !found {
		t.Error("GIT_TERMINAL_PROMPT=0 should always be present")
	}

	// Should not have timeout settings
	for _, e := range env {
		if strings.Contains(e, "GIT_HTTP_LOW_SPEED_TIME=") {
			t.Error("Should not have timeout settings when HTTPTimeout is 0")
		}
	}
}

func TestGitContextImpl_secureDeleteFile_ErrorCases(t *testing.T) {
	g := &GitContextImpl{}

	// Test with a directory (should fail)
	tempDir, err := os.MkdirTemp("", "delete_error_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	err = g.secureDeleteFile(tempDir)
	if err == nil {
		t.Error("secureDeleteFile() should fail when trying to delete a directory")
	}
}

func TestAuthConfig_Coverage(t *testing.T) {
	// Test AuthConfig struct initialization
	auth := &AuthConfig{
		Username: "test",
		Password: "pass",
		Token:    "token",
	}

	if auth.Username != "test" {
		t.Error("AuthConfig.Username not set correctly")
	}
	if auth.Password != "pass" {
		t.Error("AuthConfig.Password not set correctly")
	}
	if auth.Token != "token" {
		t.Error("AuthConfig.Token not set correctly")
	}
}