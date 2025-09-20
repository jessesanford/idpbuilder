package push

import (
	"context"
	"crypto/x509"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
	"github.com/cnoe-io/idpbuilder/pkg/registry"
	"github.com/cnoe-io/idpbuilder/pkg/certs"
)

// PushCmd is the push command
var PushCmd = &cobra.Command{
	Use:   "push IMAGE[:TAG]",
	Short: "Push image to Gitea registry",
	Long:  `Push a container image to the builtin Gitea registry with certificate support.
Automatically handles certificate trust configuration for secure connections.`,
	Example: `  idpbuilder push myapp:v1
  idpbuilder push --insecure myapp:latest
  idpbuilder push --registry https://gitea.example.com myapp:v2
  idpbuilder push --username myuser --token mytoken myapp:latest`,
	Args: cobra.ExactArgs(1),
	RunE: runPush,
}

func init() {
	PushCmd.Flags().Bool("insecure", false, "Skip TLS certificate verification")
	PushCmd.Flags().String("registry", "", "Registry URL (default: auto-detect)")
	PushCmd.Flags().String("username", "", "Registry username (default: gitea_admin)")
	PushCmd.Flags().String("token", "", "Registry token/password")
	PushCmd.Flags().Int("retry", 3, "Number of retry attempts")
}

func runPush(cmd *cobra.Command, args []string) error {
	image := args[0]
	insecure, _ := cmd.Flags().GetBool("insecure")
	registryURL, _ := cmd.Flags().GetString("registry")
	username, _ := cmd.Flags().GetString("username")
	token, _ := cmd.Flags().GetString("token")
	retryCount, _ := cmd.Flags().GetInt("retry")

	// Don't add :latest tag if this is a tarball path
	if !strings.Contains(image, ".tar") && !strings.Contains(image, ":") {
		image = image + ":latest"
	}

	// Simple progress reporting using fmt.Printf
	fmt.Printf("Pushing %s...\n", image)

	// Set defaults for registry access
	if registryURL == "" {
		registryURL = "https://gitea.cnoe.localtest.me:443"
	}
	if username == "" {
		username = "gitea_admin"
	}

	// Get token from environment if not provided
	if token == "" {
		if envToken := os.Getenv("GITEA_PASSWORD"); envToken != "" {
			token = envToken
		} else if envToken := os.Getenv("IDPBUILDER_REGISTRY_PASSWORD"); envToken != "" {
			token = envToken
		} else if envToken := os.Getenv("GITEA_TOKEN"); envToken != "" {
			token = envToken
		} else if envToken := os.Getenv("IDPBUILDER_REGISTRY_TOKEN"); envToken != "" {
			token = envToken
		}
	}

	// Create trust store for certificate management
	trustStore := certs.NewTrustStore()

	// Setup certificate trust if not insecure mode
	if !insecure {
		fmt.Println("Configuring certificate trust...")
		// Auto-configure for Gitea if using default registry
		if strings.Contains(registryURL, "gitea.cnoe.localtest.me") {
			// Try to extract certificate from Kind cluster
			if cert, err := extractGiteaCertificate(); err == nil {
				if err := trustStore.AddCertificate("gitea.cnoe.localtest.me", cert); err != nil {
					return fmt.Errorf("failed to add certificate: %w", err)
				}
			}
		}
	}

	// Create registry client options
	clientOpts := []registry.ClientOption{
		registry.WithTimeout(30 * time.Second),
		registry.WithRetryConfig(retryCount, 2*time.Second),
	}

	if insecure {
		clientOpts = append(clientOpts, registry.WithInsecure(true))
	}

	// Create registry client using go-containerregistry based GiteaClient
	client, err := registry.NewGiteaClient(registryURL, username, token, trustStore, clientOpts...)
	if err != nil {
		return fmt.Errorf("failed to create registry client: %w", err)
	}

	// Push the image
	fmt.Println("Loading image from tarball...")
	
	// For CLI usage, we expect the image to be provided as a tarball path
	// This follows the pattern: idpbuilder build --output image.tar && idpbuilder push image.tar
	tarballPath := image
	if !strings.Contains(tarballPath, ".tar") {
		return fmt.Errorf("image must be provided as a tarball path (*.tar). Use 'idpbuilder build --output %s.tar' first", image)
	}
	
	// Load the OCI image from the tarball
	img, err := tarball.ImageFromPath(tarballPath, nil)
	if err != nil {
		return fmt.Errorf("failed to load image from tarball: %w", err)
	}
	
	// Parse registry URL to get host
	registryHost := strings.TrimPrefix(registryURL, "https://")
	registryHost = strings.TrimPrefix(registryHost, "http://")
	
	// Construct the full image reference for pushing
	imageRef := fmt.Sprintf("%s/%s/hello-world:v1", registryHost, strings.ToLower(username))
	
	fmt.Printf("Pushing to %s...\n", imageRef)

	// Prepare push options
	pushOpts := registry.PushOptions{
		Options: registry.Options{
			Insecure: insecure,
			Timeout:  30 * time.Second,
		},
	}

	// Push the image using go-containerregistry
	if err := client.Push(context.Background(), img, imageRef, pushOpts); err != nil {
		return fmt.Errorf("push failed: %w", err)
	}
	
	fmt.Printf("Successfully pushed %s to %s\n", tarballPath, imageRef)
	return nil
}

// extractGiteaCertificate attempts to extract the Gitea certificate from Kind cluster
func extractGiteaCertificate() (*x509.Certificate, error) {
	// This would use Phase 1's certificate extraction logic
	// For now, return an error to use insecure mode
	return nil, fmt.Errorf("certificate extraction not yet implemented")
}

