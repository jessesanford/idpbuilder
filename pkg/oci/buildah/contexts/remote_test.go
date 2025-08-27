package contexts

import (
	"os"
	"testing"
)

func TestNewRemoteContext(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		ref      string
		wantErr  bool
	}{
		{
			name:    "valid https url",
			url:     "https://github.com/example/repo.git",
			ref:     "main",
			wantErr: false,
		},
		{
			name:    "valid ssh url",
			url:     "git@github.com:example/repo.git",
			ref:     "v1.0.0",
			wantErr: false,
		},
		{
			name:    "empty url",
			url:     "",
			ref:     "main",
			wantErr: true,
		},
		{
			name:    "invalid scheme",
			url:     "ftp://example.com/repo.git",
			ref:     "main",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, err := NewRemoteContext(tt.url, tt.ref)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRemoteContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if ctx.Type() != RemoteContextType {
					t.Errorf("Type() = %v, want %v", ctx.Type(), RemoteContextType)
				}
			}
		})
	}
}

func TestValidateGitURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"https URL", "https://github.com/user/repo.git", false},
		{"http URL", "http://example.com/repo.git", false},
		{"ssh URL", "ssh://git@github.com/user/repo.git", false},
		{"git protocol", "git://github.com/user/repo.git", false},
		{"SSH shorthand", "git@github.com:user/repo.git", false},
		{"empty URL", "", true},
		{"invalid scheme", "ftp://example.com/repo.git", true},
		{"no scheme and invalid format", "invalid-url", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateGitURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateGitURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRemoteContext_PrepareContextWithoutGit(t *testing.T) {
	// Test authentication token handling
	originalToken := os.Getenv("GIT_TOKEN")
	defer func() {
		if originalToken == "" {
			os.Unsetenv("GIT_TOKEN")
		} else {
			os.Setenv("GIT_TOKEN", originalToken)
		}
	}()

	os.Setenv("GIT_TOKEN", "test-token")

	ctx, err := NewRemoteContext("https://github.com/example/repo.git", "main")
	if err != nil {
		t.Fatalf("NewRemoteContext() error = %v", err)
	}

	if ctx.authToken != "test-token" {
		t.Errorf("authToken = %v, want %v", ctx.authToken, "test-token")
	}

	// Don't actually call PrepareContext as it requires git
	// In a real test environment, you would mock the git commands
}

func TestRemoteContext_Cleanup(t *testing.T) {
	ctx, err := NewRemoteContext("https://github.com/example/repo.git", "main")
	if err != nil {
		t.Fatalf("NewRemoteContext() error = %v", err)
	}

	// Test cleanup with no temp directory
	err = ctx.Cleanup()
	if err != nil {
		t.Errorf("Cleanup() error = %v", err)
	}

	// Set a temp directory and test cleanup
	ctx.tempDir = "/tmp/test-context"
	ctx.prepared = true
	
	err = ctx.Cleanup()
	if err != nil {
		t.Errorf("Cleanup() error = %v", err)
	}
	
	if ctx.prepared {
		t.Error("prepared should be false after cleanup")
	}
	
	if ctx.tempDir != "" {
		t.Error("tempDir should be empty after cleanup")
	}
}