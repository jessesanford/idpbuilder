package commands

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cnoe-io/idpbuilder/pkg/oci/certs"
)

// BuildOptions configures the OCI image build operation
type BuildOptions struct {
	Dockerfile string
	Tag        string
	Platform   string
	Verbose    bool
	Context    string
}

// ExecuteBuild executes the OCI image build using buildah with certificate configuration
func ExecuteBuild(ctx context.Context, contextPath string, opts BuildOptions) error {
	// Step 1: Initialize trust environment
	trustMgr, err := certs.InitializeTrustEnvironment(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize trust environment: %w", err)
	}

	if opts.Verbose {
		fmt.Printf("Trust environment initialized\n")
		if trustMgr.IsConfigured() {
			fmt.Printf("  Certificate path: %s\n", trustMgr.GetCertificatePath())
			fmt.Printf("  Registry URL: %s\n", trustMgr.GetGiteaURL())
		}
	}

	// Step 2: Validate build context
	contextDir, err := validateBuildContext(contextPath)
	if err != nil {
		return fmt.Errorf("invalid build context: %w", err)
	}

	// Step 3: Validate Dockerfile exists
	dockerfilePath := filepath.Join(contextDir, opts.Dockerfile)
	if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
		return fmt.Errorf("dockerfile not found: %s", dockerfilePath)
	}

	// Step 4: Execute build with buildah
	if err := executeBuildahBuild(ctx, contextDir, opts, trustMgr); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	fmt.Printf("✅ Successfully built image: %s\n", opts.Tag)
	return nil
}

// validateBuildContext validates and resolves the build context path
func validateBuildContext(contextPath string) (string, error) {
	if contextPath == "-" {
		return "", fmt.Errorf("stdin context not supported yet")
	}

	absPath, err := filepath.Abs(contextPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve context path: %w", err)
	}

	stat, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("context path does not exist: %s", absPath)
	}

	if !stat.IsDir() {
		return "", fmt.Errorf("context path must be a directory: %s", absPath)
	}

	return absPath, nil
}

// executeBuildahBuild executes the actual buildah build command
func executeBuildahBuild(ctx context.Context, contextDir string, opts BuildOptions, trustMgr *certs.TrustManager) error {
	// Build the buildah command
	args := []string{"build"}

	// Add dockerfile flag
	args = append(args, "--file", filepath.Join(contextDir, opts.Dockerfile))

	// Add tag flag
	args = append(args, "--tag", opts.Tag)

	// Add platform flag if specified
	if opts.Platform != "" {
		args = append(args, "--platform", opts.Platform)
	}

	// Add context directory
	args = append(args, contextDir)

	// Create command
	cmd := exec.CommandContext(ctx, "buildah", args...)

	// Set up environment for certificate trust
	env := os.Environ()
	if trustMgr.IsConfigured() && trustMgr.GetCertificatePath() != "" {
		env = append(env, "SSL_CERT_FILE="+trustMgr.GetCertificatePath())
		env = append(env, "CURL_CA_BUNDLE="+trustMgr.GetCertificatePath())
	}
	cmd.Env = env

	// Set working directory
	cmd.Dir = contextDir

	// Configure output
	if opts.Verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		fmt.Printf("Executing: buildah %s\n", strings.Join(args, " "))
	} else {
		// Capture output to show progress
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	// Execute the build
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("buildah command failed: %w", err)
	}

	return nil
}


// CheckBuildahAvailable checks if buildah is available on the system
func CheckBuildahAvailable() error {
	_, err := exec.LookPath("buildah")
	if err != nil {
		return fmt.Errorf("buildah not found in PATH - please install buildah to use build command")
	}
	return nil
}