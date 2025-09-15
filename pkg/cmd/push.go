package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/cnoe-io/idpbuilder/pkg/certs"
	"github.com/cnoe-io/idpbuilder/pkg/cmd/helpers"
	"github.com/cnoe-io/idpbuilder/pkg/gitea"
	"github.com/spf13/cobra"
)

var (
	pushInsecure bool
	pushRegistry string
	pushUsername string
	pushToken    string

	PushCmd = &cobra.Command{
		Use:   "push IMAGE[:TAG]",
		Short: "Push image to Gitea registry",
		Long: `Push a container image to the builtin Gitea registry with certificate support.
The command automatically handles certificate extraction and configuration unless
the --insecure flag is specified.

Examples:
  # Push with automatic certificate handling
  idpbuilder push myapp:latest

  # Push to specific registry
  idpbuilder push --registry gitea.cnoe.localtest.me:8443 myapp:latest

  # Push with explicit credentials
  idpbuilder push --username admin --token mytoken myapp:latest

  # Push without certificate verification (not recommended)
  idpbuilder push --insecure myapp:latest`,
		Args: cobra.ExactArgs(1),
		RunE: runPush,
	}
)

func init() {
	PushCmd.Flags().BoolVar(&pushInsecure, "insecure", false, "Skip certificate verification (not recommended)")
	PushCmd.Flags().StringVar(&pushRegistry, "registry", "gitea.cnoe.localtest.me:8443", "Registry endpoint")
	PushCmd.Flags().StringVar(&pushUsername, "username", "", "Registry username (if not provided, will attempt auto-detection)")
	PushCmd.Flags().StringVar(&pushToken, "token", "", "Registry token/password (if not provided, will attempt auto-detection)")
}

func runPush(cmd *cobra.Command, args []string) error {
	imageName := args[0]
	
	// Ensure image has a tag
	if !strings.Contains(imageName, ":") {
		imageName = imageName + ":latest"
	}

	// Extract certificate unless insecure mode
	var certPath string
	if !pushInsecure {
		cert, err := certs.ExtractCertificate(cmd.Context())
		if err != nil {
			return fmt.Errorf("failed to extract certificate: %w", err)
		}
		
		certPath = cert.FilePath
		defer func() {
			// Clean up temp cert file
			if certPath != "" {
				os.Remove(certPath)
			}
		}()
		
		fmt.Printf("✓ Certificate extracted from ingress controller\n")
		fmt.Printf("✓ Registry: %s\n", pushRegistry)
	} else {
		fmt.Printf("⚠ Warning: Running in insecure mode - certificate verification disabled\n")
	}

	// Create Gitea client with appropriate configuration
	var client *gitea.Client
	var err error
	
	if pushUsername != "" && pushToken != "" {
		// Use explicit credentials
		client, err = gitea.NewClient(pushRegistry, pushUsername, pushToken, certPath)
	} else {
		// Attempt auto-detection of credentials
		client, err = gitea.NewClientAutoDetect(pushRegistry, certPath)
	}
	
	if err != nil {
		return fmt.Errorf("failed to create registry client: %w", err)
	}

	// Validate credentials before push
	if err := client.ValidateCredentials(cmd.Context()); err != nil {
		return fmt.Errorf("failed to validate credentials: %w", err)
	}
	fmt.Printf("✓ Credentials validated\n")

	// Push the image
	fmt.Printf("Pushing %s to %s...\n", imageName, pushRegistry)
	if err := client.PushImage(cmd.Context(), imageName); err != nil {
		return fmt.Errorf("failed to push image: %w", err)
	}

	fmt.Printf("✓ Image pushed successfully\n")
	return nil
}
