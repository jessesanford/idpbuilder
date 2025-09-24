package push

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/cnoe-io/idpbuilder/pkg/cmd/helpers"
	"github.com/cnoe-io/idpbuilder/pkg/oci"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

const (
	imageUsage        = "Image reference to push (e.g., localhost:5000/myimage:latest)"
	usernameUsage     = "Registry username for authentication"
	passwordUsage     = "Registry password for authentication"
	insecureUsage     = "Allow connections to non-HTTPS registries"
	authModeUsage     = "Authentication mode: flags, env, docker-config, or auto (default: auto)"
	tagUsage          = "Additional tags to apply to the image before pushing"
	quietUsage        = "Suppress push progress output"
	forceUsage        = "Force push even if image already exists in registry"
)

var (
	// Command flags
	imageRef   string
	username   string
	password   string
	insecure   bool
	authMode   string
	tags       []string
	quiet      bool
	force      bool
)

// PushCmd represents the push command
var PushCmd = &cobra.Command{
	Use:          "push IMAGE",
	Short:        "Push a container image to an OCI registry",
	Long: `Push a container image to an OCI registry with authentication support.

The push command supports multiple authentication methods:
- CLI flags (--username, --password)
- Environment variables (OCI_USERNAME, OCI_PASSWORD)
- Docker config file (~/.docker/config.json)
- Automatic detection (tries all methods)

Examples:
  # Push with explicit authentication
  idpbuilder push localhost:5000/myapp:v1.0 --username admin --password secret

  # Push with environment variables
  export OCI_USERNAME=admin OCI_PASSWORD=secret
  idpbuilder push localhost:5000/myapp:v1.0

  # Push with additional tags
  idpbuilder push localhost:5000/myapp:v1.0 --tag latest --tag stable

  # Push to insecure registry
  idpbuilder push localhost:5000/myapp:v1.0 --insecure`,
	Args:         cobra.ExactArgs(1),
	RunE:         pushImage,
	PreRunE:      prePushE,
	SilenceUsage: true,
}

func init() {
	PushCmd.Flags().StringVarP(&username, "username", "u", "", usernameUsage)
	PushCmd.Flags().StringVarP(&password, "password", "p", "", passwordUsage)
	PushCmd.Flags().BoolVar(&insecure, "insecure", false, insecureUsage)
	PushCmd.Flags().StringVar(&authMode, "auth-mode", "auto", authModeUsage)
	PushCmd.Flags().StringSliceVarP(&tags, "tag", "t", []string{}, tagUsage)
	PushCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, quietUsage)
	PushCmd.Flags().BoolVar(&force, "force", false, forceUsage)
}

func prePushE(cmd *cobra.Command, args []string) error {
	imageRef = args[0]
	return helpers.SetLogger()
}

func pushImage(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	// Validate image reference
	if imageRef == "" {
		return fmt.Errorf("image reference is required")
	}

	if !quiet {
		fmt.Printf("Pushing image: %s\n", imageRef)
	}

	// Initialize Docker client
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %w", err)
	}
	defer dockerClient.Close()

	// Verify image exists locally
	if err := verifyLocalImage(ctx, dockerClient, imageRef); err != nil {
		return fmt.Errorf("image verification failed: %w", err)
	}

	// Apply additional tags if specified
	if len(tags) > 0 {
		if err := applyAdditionalTags(ctx, dockerClient, imageRef, tags); err != nil {
			return fmt.Errorf("failed to apply tags: %w", err)
		}
	}

	// Extract registry from image reference
	registry, err := extractRegistry(imageRef)
	if err != nil {
		return fmt.Errorf("failed to parse image reference: %w", err)
	}

	// Set up authentication
	auth, err := setupAuthentication(ctx, registry)
	if err != nil {
		return fmt.Errorf("authentication setup failed: %w", err)
	}

	// Perform the push
	if err := performPush(ctx, dockerClient, imageRef, auth); err != nil {
		return fmt.Errorf("push operation failed: %w", err)
	}

	// Push additional tags if any were applied
	for _, tag := range tags {
		taggedImage := replaceTag(imageRef, tag)
		if err := performPush(ctx, dockerClient, taggedImage, auth); err != nil {
			return fmt.Errorf("failed to push tag %s: %w", tag, err)
		}
	}

	if !quiet {
		fmt.Printf("Successfully pushed image: %s\n", imageRef)
		if len(tags) > 0 {
			fmt.Printf("Additional tags pushed: %v\n", tags)
		}
	}

	return nil
}

// verifyLocalImage checks if the image exists locally
func verifyLocalImage(ctx context.Context, dockerClient *client.Client, imageRef string) error {
	images, err := dockerClient.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list local images: %w", err)
	}

	for _, img := range images {
		for _, tag := range img.RepoTags {
			if tag == imageRef {
				return nil
			}
		}
		// Also check by ID if the reference looks like an ID
		if strings.HasPrefix(img.ID, "sha256:"+imageRef) || img.ID == imageRef {
			return nil
		}
	}

	return fmt.Errorf("image %s not found locally", imageRef)
}

// applyAdditionalTags applies additional tags to the image
func applyAdditionalTags(ctx context.Context, dockerClient *client.Client, imageRef string, tags []string) error {
	for _, tag := range tags {
		newTag := replaceTag(imageRef, tag)
		if err := dockerClient.ImageTag(ctx, imageRef, newTag); err != nil {
			return fmt.Errorf("failed to tag image %s as %s: %w", imageRef, newTag, err)
		}
		if !quiet {
			fmt.Printf("Tagged image as: %s\n", newTag)
		}
	}
	return nil
}

// replaceTag replaces the tag portion of an image reference
func replaceTag(imageRef, newTag string) string {
	// Find the last colon to identify tag separator
	lastColonIndex := strings.LastIndex(imageRef, ":")

	// Check if the colon is part of a port number (registry:port/image)
	// If there's a slash after the last colon, it's part of a port
	if lastColonIndex != -1 {
		afterColon := imageRef[lastColonIndex+1:]
		if !strings.Contains(afterColon, "/") {
			// This colon is separating tag, not port
			return imageRef[:lastColonIndex] + ":" + newTag
		}
	}

	// No tag present, append it
	return imageRef + ":" + newTag
}

// extractRegistry extracts the registry hostname from an image reference
func extractRegistry(imageRef string) (string, error) {
	// Handle special cases
	if !strings.Contains(imageRef, "/") {
		return "docker.io", nil
	}

	parts := strings.Split(imageRef, "/")
	if len(parts) == 0 {
		return "", fmt.Errorf("invalid image reference: %s", imageRef)
	}

	// First part should be the registry if it contains a dot or port
	registry := parts[0]
	if strings.Contains(registry, ".") || strings.Contains(registry, ":") {
		return registry, nil
	}

	// Default to docker.io for simple names
	return "docker.io", nil
}

// setupAuthentication sets up authentication based on the specified mode
func setupAuthentication(ctx context.Context, registry string) (oci.Authenticator, error) {
	switch authMode {
	case "flags":
		if username == "" || password == "" {
			return nil, fmt.Errorf("username and password required for flags auth mode")
		}
		return oci.NewAuthenticatorFromFlags(username, password)

	case "env":
		return oci.NewAuthenticatorFromEnv()

	case "docker-config":
		return oci.NewAuthenticator(&oci.AuthConfig{
			Sources: []oci.CredentialSource{oci.SourceDockerConfig},
		}), nil

	case "auto":
		// Try flags first if provided
		if username != "" && password != "" {
			return oci.NewAuthenticatorFromFlags(username, password)
		}

		// Create authenticator with all sources
		return oci.NewAuthenticator(&oci.AuthConfig{
			Sources: []oci.CredentialSource{
				oci.SourceEnvironment,
				oci.SourceDockerConfig,
			},
		}), nil

	default:
		return nil, fmt.Errorf("unsupported auth mode: %s (supported: flags, env, docker-config, auto)", authMode)
	}
}

// performPush executes the actual push operation
func performPush(ctx context.Context, dockerClient *client.Client, imageRef string, auth oci.Authenticator) error {
	// Get registry from image reference for authentication
	registry, err := extractRegistry(imageRef)
	if err != nil {
		return fmt.Errorf("failed to extract registry: %w", err)
	}

	// Get credentials if authenticator is available
	var authConfig string
	if auth != nil {
		creds, err := auth.Authenticate(ctx, registry)
		if err != nil {
			if !quiet {
				fmt.Printf("Warning: authentication failed for %s: %v\n", registry, err)
				fmt.Println("Attempting push without authentication...")
			}
		} else if creds != nil && creds.IsValid() {
			// Create Docker auth config
			authConfig = createAuthConfig(creds.Username, creds.Password)
		}
	}

	// Set up push options
	pushOptions := types.ImagePushOptions{
		RegistryAuth: authConfig,
	}

	// Perform the push
	pushResponse, err := dockerClient.ImagePush(ctx, imageRef, pushOptions)
	if err != nil {
		return fmt.Errorf("docker push failed: %w", err)
	}
	defer pushResponse.Close()

	// Handle push response stream
	if !quiet {
		if err := streamPushResponse(pushResponse); err != nil {
			return fmt.Errorf("failed to read push response: %w", err)
		}
	} else {
		// Even if quiet, we need to consume the response to ensure push completes
		if err := consumePushResponse(pushResponse); err != nil {
			return fmt.Errorf("push failed: %w", err)
		}
	}

	return nil
}

// createAuthConfig creates a Docker auth config string
func createAuthConfig(username, password string) string {
	if username == "" && password == "" {
		return ""
	}

	// Create basic auth string (base64 encoded username:password)
	authData := map[string]string{
		"username": username,
		"password": password,
	}

	jsonData, _ := json.Marshal(authData)
	return base64.URLEncoding.EncodeToString(jsonData)
}

// streamPushResponse streams push progress to stdout
func streamPushResponse(response io.ReadCloser) error {
	// For now, just copy to stdout
	// In a production implementation, you would parse JSON progress messages
	_, err := io.Copy(os.Stdout, response)
	return err
}

// consumePushResponse consumes push response without output
func consumePushResponse(response io.ReadCloser) error {
	// Read and discard all response data to ensure push completes
	_, err := io.Copy(io.Discard, response)
	return err
}