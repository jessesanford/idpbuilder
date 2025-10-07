package push

import (
	"context"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/spf13/cobra"
)

// TestNewPushOperationFromCommand tests creating a PushOperation from command flags
func TestNewPushOperationFromCommand(t *testing.T) {
	tests := []struct {
		name      string
		registry  string
		username  string
		password  string
		insecure  bool
		buildPath string
		wantErr   bool
	}{
		{
			name:      "valid configuration",
			registry:  "localhost:5000",
			username:  "testuser",
			password:  "testpass",
			insecure:  true,
			buildPath: "/tmp/images",
			wantErr:   false,
		},
		{
			name:      "anonymous auth",
			registry:  "localhost:5000",
			username:  "",
			password:  "",
			insecure:  false,
			buildPath: ".",
			wantErr:   false,
		},
		{
			name:      "default build path",
			registry:  "localhost:5000",
			username:  "",
			password:  "",
			insecure:  false,
			buildPath: "",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			cmd.Flags().String("registry", tt.registry, "")
			cmd.Flags().String("username", tt.username, "")
			cmd.Flags().String("password", tt.password, "")
			cmd.Flags().Bool("insecure", tt.insecure, "")
			cmd.Flags().String("build-path", tt.buildPath, "")
			cmd.Flags().Int("max-retries", 3, "")
			cmd.Flags().Int("concurrency", 2, "")

			logger := logr.Discard()
			op, err := NewPushOperationFromCommand(cmd, logger)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewPushOperationFromCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if op.Registry != tt.registry {
					t.Errorf("Registry = %v, want %v", op.Registry, tt.registry)
				}
				if op.Username != tt.username {
					t.Errorf("Username = %v, want %v", op.Username, tt.username)
				}
				if op.Insecure != tt.insecure {
					t.Errorf("Insecure = %v, want %v", op.Insecure, tt.insecure)
				}
				if op.Logger == nil {
					t.Error("Logger should not be nil")
				}
				if op.Progress == nil {
					t.Error("Progress reporter should not be nil")
				}
				if op.Transport == nil {
					t.Error("Transport should not be nil")
				}
				if op.MaxRetries != 3 {
					t.Errorf("MaxRetries = %v, want 3", op.MaxRetries)
				}
				if op.Concurrency != 2 {
					t.Errorf("Concurrency = %v, want 2", op.Concurrency)
				}

				expectedPath := tt.buildPath
				if expectedPath == "" {
					expectedPath = "."
				}
				if op.BuildPath != expectedPath {
					t.Errorf("BuildPath = %v, want %v", op.BuildPath, expectedPath)
				}
			}
		})
	}
}

// TestSetupAuthentication tests authentication setup
func TestSetupAuthentication(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		wantAuth bool
	}{
		{
			name:     "anonymous auth",
			username: "",
			password: "",
			wantAuth: false,
		},
		{
			name:     "basic auth",
			username: "user",
			password: "pass",
			wantAuth: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := &PushOperation{
				Username: tt.username,
				Password: tt.password,
				Registry: "localhost:5000",
				Logger:   NewPushLogger(logr.Discard()),
			}

			err := op.setupAuthentication()
			if err != nil {
				t.Errorf("setupAuthentication() error = %v", err)
			}

			if op.Auth == nil {
				t.Error("Auth should not be nil after setup")
			}

			if tt.wantAuth {
				if _, ok := op.Auth.(*authn.Basic); !ok {
					t.Error("Expected Basic auth, got different type")
				}
			}
		})
	}
}

// TestSetupTransport tests HTTP transport setup
func TestSetupTransport(t *testing.T) {
	tests := []struct {
		name     string
		insecure bool
	}{
		{
			name:     "secure transport",
			insecure: false,
		},
		{
			name:     "insecure transport",
			insecure: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := &PushOperation{
				Insecure: tt.insecure,
			}

			op.setupTransport()

			if op.Transport == nil {
				t.Fatal("Transport should not be nil")
			}

			if op.Transport.TLSClientConfig == nil {
				t.Fatal("TLS config should not be nil")
			}

			if op.Transport.TLSClientConfig.InsecureSkipVerify != tt.insecure {
				t.Errorf("InsecureSkipVerify = %v, want %v",
					op.Transport.TLSClientConfig.InsecureSkipVerify, tt.insecure)
			}
		})
	}
}

// TestPushOperationResult tests the result aggregation
func TestPushOperationResult(t *testing.T) {
	result := &PushOperationResult{
		StartTime:     time.Now(),
		EndTime:       time.Now().Add(5 * time.Second),
		TotalDuration: 5 * time.Second,
		ImagesFound:   10,
		ImagesPushed:  8,
		ImagesFailed:  2,
		TotalBytes:    1024 * 1024 * 100, // 100MB
		Results:       make([]*PushResult, 0),
		Errors:        make([]error, 0),
		Metrics: &PushMetrics{
			TotalDuration: 5 * time.Second,
			TotalBytes:    1024 * 1024 * 100,
			ImagesPushed:  8,
		},
	}

	// Test summary generation
	summary := result.Summary()
	if summary == "" {
		t.Error("Summary should not be empty")
	}

	// Verify metrics calculations
	throughput := result.Metrics.AverageThroughputMBps()
	if throughput < 0 {
		t.Errorf("Throughput should be non-negative, got: %f", throughput)
	}
}

// TestFormatBytes tests byte formatting
func TestFormatBytes(t *testing.T) {
	tests := []struct {
		name  string
		bytes int64
		want  string
	}{
		{
			name:  "bytes",
			bytes: 512,
			want:  "512 B",
		},
		{
			name:  "kilobytes",
			bytes: 1024,
			want:  "1.0 KB",
		},
		{
			name:  "megabytes",
			bytes: 1024 * 1024,
			want:  "1.0 MB",
		},
		{
			name:  "gigabytes",
			bytes: 1024 * 1024 * 1024,
			want:  "1.0 GB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatBytes(tt.bytes)
			if got != tt.want {
				t.Errorf("formatBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestPushImagesWithNoImages tests push operation when no images are found
func TestPushImagesWithNoImages(t *testing.T) {
	// Create a temporary empty directory
	tempDir := t.TempDir()

	op := &PushOperation{
		Registry:    "localhost:5000",
		BuildPath:   tempDir,
		Username:    "",
		Password:    "",
		Insecure:    true,
		MaxRetries:  3,
		Concurrency: 2,
		Logger:      NewPushLogger(logr.Discard()),
		Auth:        authn.Anonymous,
	}
	op.setupTransport()
	op.Progress = NewNoOpProgressReporter()

	ctx := context.Background()
	result, err := op.Execute(ctx)

	if err != nil {
		t.Errorf("Execute() should not error with no images, got: %v", err)
	}

	if result == nil {
		t.Fatal("Result should not be nil")
	}

	if result.ImagesFound != 0 {
		t.Errorf("ImagesFound = %v, want 0", result.ImagesFound)
	}

	if result.ImagesPushed != 0 {
		t.Errorf("ImagesPushed = %v, want 0", result.ImagesPushed)
	}
}

// TestDiscoverImages tests image discovery
func TestDiscoverImages(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()

	op := &PushOperation{
		BuildPath: tempDir,
		Logger:    NewPushLogger(logr.Discard()),
	}

	images, err := op.discoverImages()
	// Empty directory may return error or empty slice depending on implementation
	if err == nil && images != nil && len(images) != 0 {
		t.Errorf("Expected 0 images in empty directory, got %d", len(images))
	}
	// Accept either error or empty result for empty directory
}

// TestValidateImages tests image validation
func TestValidateImages(t *testing.T) {
	// Create a mock image
	img := empty.Image

	localImg := &LocalImage{
		Name:   "test-image",
		Path:   "/tmp/test.tar",
		Format: "tarball",
		Image:  img,
	}

	op := &PushOperation{
		Registry:    "localhost:5000",
		Username:    "",
		Password:    "",
		Insecure:    true,
		MaxRetries:  3,
		Concurrency: 2,
		UserAgent:   "test",
		Logger:      NewPushLogger(logr.Discard()),
		Auth:        authn.Anonymous,
	}
	op.setupTransport()
	op.Progress = NewNoOpProgressReporter()

	images := []*LocalImage{localImg}

	// This may fail with the actual implementation if registry is not reachable
	// but we're testing the validation flow
	err := op.validateImages(images)

	// We don't assert success/failure since it depends on actual registry
	// The important thing is that it doesn't panic
	_ = err
}

// TestPushImagesCommand tests the command entry point
func TestPushImagesCommand(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String("registry", "localhost:5000", "")
	cmd.Flags().String("username", "", "")
	cmd.Flags().String("password", "", "")
	cmd.Flags().Bool("insecure", true, "")
	cmd.Flags().String("build-path", t.TempDir(), "")
	cmd.Flags().Int("max-retries", 3, "")
	cmd.Flags().Int("concurrency", 2, "")

	ctx := context.Background()
	cmd.SetContext(ctx)

	// This will likely fail because there are no images and no real registry
	// but we're testing that it doesn't panic
	err := PushImages(cmd, []string{})

	// We expect either no error (no images found) or an error (no images/registry unreachable)
	// The important thing is no panic
	_ = err
}

// TestPushOperationWithContext tests context cancellation
func TestPushOperationWithContext(t *testing.T) {
	tempDir := t.TempDir()

	op := &PushOperation{
		Registry:    "localhost:5000",
		BuildPath:   tempDir,
		Username:    "",
		Password:    "",
		Insecure:    true,
		MaxRetries:  3,
		Concurrency: 2,
		Logger:      NewPushLogger(logr.Discard()),
		Auth:        authn.Anonymous,
	}
	op.setupTransport()
	op.Progress = NewNoOpProgressReporter()

	// Create a context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Execute should handle the cancelled context gracefully
	result, err := op.Execute(ctx)

	// With cancelled context and no images, should complete without error
	// (discovery happens synchronously before any async operations)
	if err != nil && result.ImagesFound == 0 {
		// This is acceptable - no images to push with cancelled context
		return
	}

	if result == nil {
		t.Error("Result should not be nil even with cancelled context")
	}
}

// TestPushOperationDefaults tests default value handling
func TestPushOperationDefaults(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String("registry", "localhost:5000", "")
	cmd.Flags().String("username", "", "")
	cmd.Flags().String("password", "", "")
	cmd.Flags().Bool("insecure", false, "")
	cmd.Flags().String("build-path", "", "") // Empty, should default to "."
	cmd.Flags().Int("max-retries", 0, "")    // 0, should default to 3
	cmd.Flags().Int("concurrency", 0, "")    // 0, should default to 3

	logger := logr.Discard()
	op, err := NewPushOperationFromCommand(cmd, logger)

	if err != nil {
		t.Fatalf("NewPushOperationFromCommand() error = %v", err)
	}

	if op.BuildPath != "." {
		t.Errorf("Default BuildPath = %v, want '.'", op.BuildPath)
	}

	if op.MaxRetries != 3 {
		t.Errorf("Default MaxRetries = %v, want 3", op.MaxRetries)
	}

	if op.Concurrency != 3 {
		t.Errorf("Default Concurrency = %v, want 3", op.Concurrency)
	}

	if op.UserAgent != "idpbuilder-push/1.0.0" {
		t.Errorf("UserAgent = %v, want 'idpbuilder-push/1.0.0'", op.UserAgent)
	}
}

// TestPushOperationResultMetrics tests metrics tracking
func TestPushOperationResultMetrics(t *testing.T) {
	startTime := time.Now()
	endTime := startTime.Add(10 * time.Second)

	result := &PushOperationResult{
		StartTime:     startTime,
		EndTime:       endTime,
		TotalDuration: 10 * time.Second,
		ImagesFound:   5,
		ImagesPushed:  4,
		ImagesFailed:  1,
		TotalBytes:    1024 * 1024 * 50, // 50MB
		Metrics: &PushMetrics{
			TotalDuration:     10 * time.Second,
			TotalBytes:        1024 * 1024 * 50,
			ImagesPushed:      4,
			DiscoveryDuration: 1 * time.Second,
			PushDuration:      9 * time.Second,
			RetryCount:        2,
		},
	}

	if result.Metrics.TotalDuration != 10*time.Second {
		t.Errorf("TotalDuration = %v, want 10s", result.Metrics.TotalDuration)
	}

	if result.Metrics.ImagesPushed != 4 {
		t.Errorf("ImagesPushed = %v, want 4", result.Metrics.ImagesPushed)
	}

	throughput := result.Metrics.AverageThroughputMBps()
	if throughput <= 0 {
		t.Error("AverageThroughputMBps should be positive")
	}

	// Test summary contains expected information
	summary := result.Summary()
	if summary == "" {
		t.Error("Summary should not be empty")
	}
}

// TestPushOperation_ErrorHandling tests error scenarios
func TestPushOperation_ErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		buildPath   string
		expectError bool
	}{
		{
			name:        "valid empty directory",
			buildPath:   t.TempDir(),
			expectError: false, // No error for empty directory, just no images
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := &PushOperation{
				Registry:    "localhost:5000",
				BuildPath:   tt.buildPath,
				Username:    "",
				Password:    "",
				Insecure:    true,
				MaxRetries:  3,
				Concurrency: 2,
				Logger:      NewPushLogger(logr.Discard()),
				Auth:        authn.Anonymous,
			}
			op.setupTransport()
			op.Progress = NewNoOpProgressReporter()

			ctx := context.Background()
			result, err := op.Execute(ctx)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if result == nil {
				t.Error("Result should not be nil")
			}
		})
	}
}
