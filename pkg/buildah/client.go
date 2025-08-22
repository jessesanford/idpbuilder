package buildah

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
)

// client implements the Client interface using the buildah binary
type client struct {
	// logger for structured logging
	logger logr.Logger

	// buildahPath is the path to the buildah binary
	buildahPath string

	// workDir is the working directory for builds
	workDir string

	// defaultRegistry is the default registry for operations
	defaultRegistry string

	// timeout is the default timeout for operations
	timeout time.Duration

	// logLevel controls buildah verbosity
	logLevel string
}

// NewClient creates a new buildah client with the given options
func NewClient(opts ClientOptions) (Client, error) {
	logger := ctrl.Log.WithName("buildah-client")

	// Determine buildah binary path
	buildahPath := opts.BuildahPath
	if buildahPath == "" {
		buildahPath = "buildah"
	}

	// Check if buildah binary is available
	if err := checkBuildahBinary(buildahPath); err != nil {
		return nil, fmt.Errorf("buildah binary not available: %w", err)
	}

	// Set up working directory
	workDir := opts.WorkDir
	if workDir == "" {
		var err error
		workDir, err = os.MkdirTemp("", "buildah-client-")
		if err != nil {
			return nil, fmt.Errorf("failed to create temp directory: %w", err)
		}
		logger.V(1).Info("Created temp working directory", "dir", workDir)
	}

	// Ensure working directory exists
	if err := os.MkdirAll(workDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create working directory: %w", err)
	}

	// Set default timeout
	timeout := opts.Timeout
	if timeout == 0 {
		timeout = 10 * time.Minute
	}

	// Set default log level
	logLevel := opts.LogLevel
	if logLevel == "" {
		logLevel = "warn"
	}

	c := &client{
		logger:          logger,
		buildahPath:     buildahPath,
		workDir:         workDir,
		defaultRegistry: opts.DefaultRegistry,
		timeout:         timeout,
		logLevel:        logLevel,
	}

	logger.Info("Buildah client created", 
		"buildah_path", buildahPath,
		"work_dir", workDir,
		"default_registry", opts.DefaultRegistry,
		"timeout", timeout)

	return c, nil
}

// Build builds a container image from a Dockerfile using buildah
func (c *client) Build(ctx context.Context, opts BuildOptions) (*BuildResult, error) {
	start := time.Now()
	
	c.logger.Info("Starting container build",
		"context_dir", opts.ContextDir,
		"dockerfile", opts.Dockerfile,
		"tag", opts.Tag)

	// Validate build options
	if err := c.validateBuildOptions(opts); err != nil {
		return nil, &BuildError{
			Op:  "validate_options",
			Err: err,
		}
	}

	// Build the command arguments
	args, err := c.buildBuildArgs(opts)
	if err != nil {
		return nil, &BuildError{
			Op:  "build_args",
			Err: err,
		}
	}

	// Execute the buildah build command
	output, err := c.runBuildahCommand(ctx, args, opts)
	if err != nil {
		return nil, err
	}

	// Parse the build result
	result, err := c.parseBuildResult(opts, output, time.Since(start))
	if err != nil {
		return nil, &BuildError{
			Op:     "parse_result",
			Err:    err,
			Output: string(output),
		}
	}

	c.logger.Info("Container build completed",
		"image_id", result.ImageID,
		"tag", result.Tag,
		"duration", result.Duration)

	return result, nil
}

// Push pushes a container image to a registry
func (c *client) Push(ctx context.Context, opts PushOptions) error {
	c.logger.Info("Starting image push",
		"image", opts.Image,
		"registry", opts.Registry)

	// Build the command arguments
	args := []string{"push"}
	
	// Add authentication if provided
	if opts.Username != "" && opts.Password != "" {
		args = append(args, "--creds", opts.Username+":"+opts.Password)
	}

	// Add the image to push
	args = append(args, opts.Image)

	// If registry is specified, include it in the destination
	dest := opts.Image
	if opts.Registry != "" && opts.Registry != c.defaultRegistry {
		// Parse image name and add registry prefix if needed
		dest = opts.Registry + "/" + opts.Image
	}
	
	args[len(args)-1] = dest

	// Execute the buildah push command
	_, err := c.runBuildahCommand(ctx, args, BuildOptions{Progress: opts.Progress})
	if err != nil {
		return err
	}

	c.logger.Info("Image push completed", "destination", dest)
	return nil
}

// Close cleans up client resources
func (c *client) Close() error {
	c.logger.Info("Closing buildah client")
	
	// Clean up temporary working directory if we created it
	if c.workDir != "" && filepath.Base(c.workDir) == filepath.Base(os.TempDir()) {
		if err := os.RemoveAll(c.workDir); err != nil {
			c.logger.Error(err, "Failed to clean up working directory", "dir", c.workDir)
			return fmt.Errorf("failed to clean up working directory: %w", err)
		}
		c.logger.V(1).Info("Cleaned up working directory", "dir", c.workDir)
	}

	return nil
}

// checkBuildahBinary verifies that the buildah binary is available and executable
func checkBuildahBinary(buildahPath string) error {
	cmd := exec.Command(buildahPath, "version")
	if err := cmd.Run(); err != nil {
		if pathErr, ok := err.(*exec.Error); ok && pathErr.Err == exec.ErrNotFound {
			return fmt.Errorf("buildah binary not found at %s", buildahPath)
		}
		return fmt.Errorf("buildah binary not working: %w", err)
	}
	return nil
}