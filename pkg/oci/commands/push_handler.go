package commands

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/cnoe-io/idpbuilder/pkg/oci/certs"
)

// PushOptions configures the OCI image push operation
type PushOptions struct {
	Insecure bool
	Username string
	Password string
	Verbose  bool
}

// ExecutePush executes the OCI image push using buildah with certificate configuration
func ExecutePush(ctx context.Context, imageRef string, opts PushOptions) error {
	// Step 1: Validate image reference
	if err := validateImageReference(imageRef); err != nil {
		return fmt.Errorf("invalid image reference: %w", err)
	}

	// Step 2: Initialize trust environment (unless --insecure)
	var trustMgr *certs.TrustManager
	var err error
	
	if opts.Insecure {
		trustMgr, err = certs.InitializeFallbackEnvironment()
		if err != nil {
			return fmt.Errorf("failed to initialize fallback environment: %w", err)
		}
		if opts.Verbose {
			fmt.Printf("Running in insecure mode - skipping certificate verification\n")
		}
	} else {
		trustMgr, err = certs.InitializeTrustEnvironment(ctx)
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
	}

	// Step 3: Execute push with buildah
	if err := executeBuildahPush(ctx, imageRef, opts, trustMgr); err != nil {
		return fmt.Errorf("push failed: %w", err)
	}

	fmt.Printf("✅ Successfully pushed image: %s\n", imageRef)
	return nil
}

// validateImageReference validates the format of the image reference
func validateImageReference(imageRef string) error {
	if imageRef == "" {
		return fmt.Errorf("image reference cannot be empty")
	}

	// Basic validation - must contain a colon for tag
	if !strings.Contains(imageRef, ":") {
		return fmt.Errorf("image reference must include tag: %s", imageRef)
	}

	// Must contain at least one slash for registry/repository format
	parts := strings.Split(imageRef, "/")
	if len(parts) < 2 {
		return fmt.Errorf("image reference must be in format [registry/]repository:tag")
	}

	// Extract tag part
	lastPart := parts[len(parts)-1]
	tagParts := strings.Split(lastPart, ":")
	if len(tagParts) != 2 {
		return fmt.Errorf("invalid tag format in image reference")
	}

	return nil
}

// executeBuildahPush executes the actual buildah push command
func executeBuildahPush(ctx context.Context, imageRef string, opts PushOptions, trustMgr *certs.TrustManager) error {
	// Build the buildah command
	args := []string{"push"}

	// Add TLS verification settings
	if opts.Insecure {
		args = append(args, "--tls-verify=false")
	} else {
		args = append(args, "--tls-verify=true")
	}

	// Add authentication if provided
	if opts.Username != "" && opts.Password != "" {
		args = append(args, "--creds", fmt.Sprintf("%s:%s", opts.Username, opts.Password))
	}

	// Add image reference
	args = append(args, imageRef)

	// Create command
	cmd := exec.CommandContext(ctx, "buildah", args...)

	// Set up environment for certificate trust
	env := os.Environ()
	if !opts.Insecure && trustMgr.IsConfigured() && trustMgr.GetCertificatePath() != "" {
		env = append(env, "SSL_CERT_FILE="+trustMgr.GetCertificatePath())
		env = append(env, "CURL_CA_BUNDLE="+trustMgr.GetCertificatePath())
	}
	
	// Set insecure environment variables if needed
	if opts.Insecure {
		env = append(env, "BUILDAH_TLS_VERIFY=false")
	}
	
	cmd.Env = env

	// Configure output
	if opts.Verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		fmt.Printf("Executing: buildah %s\n", strings.Join(args, " "))
	} else {
		// Capture output to show progress
		cmd.Stdout = &pushProgressWriter{}
		cmd.Stderr = &pushProgressWriter{}
	}

	// Execute the push
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("buildah push command failed: %w", err)
	}

	return nil
}

// pushProgressWriter writes push progress with formatting
type pushProgressWriter struct{}

func (ppw *pushProgressWriter) Write(p []byte) (n int, err error) {
	fmt.Print(string(p))
	return len(p), nil
}

// CheckPodmanAvailable checks if podman/buildah is available for push operations
func CheckPodmanAvailable() error {
	// Check for buildah first (preferred)
	if _, err := exec.LookPath("buildah"); err == nil {
		return nil
	}
	
	// Fall back to podman
	if _, err := exec.LookPath("podman"); err == nil {
		return nil
	}
	
	return fmt.Errorf("neither buildah nor podman found in PATH - please install one to use push command")
}