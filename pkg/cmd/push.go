package cmd

import (
	"context"
	"fmt"

	"github.com/cnoe-io/idpbuilder/pkg/certs"
	"github.com/cnoe-io/idpbuilder/pkg/registry"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
	"github.com/spf13/cobra"
)

var (
	pushImage    string
	pushInsecure bool
	pushSkipTLS  bool
	pushUsername string
	pushPassword string

	PushCmd = &cobra.Command{
		Use:   "push",
		Short: "Push OCI image to registry",
		Long: `Push an OCI image to a container registry using go-containerregistry.

Examples:
  # Push to local Gitea registry
  idpbuilder push myimage.tar --insecure

  # Push with authentication
  idpbuilder push myimage.tar --username admin --password secret

  # Push with TLS verification skipped
  idpbuilder push myimage.tar --skip-tls-verify`,
		Args: cobra.ExactArgs(1),
		RunE: runPush,
	}
)

func init() {
	PushCmd.Flags().BoolVar(&pushInsecure, "insecure", false, "Allow insecure registry connections")
	PushCmd.Flags().BoolVar(&pushSkipTLS, "skip-tls-verify", false, "Skip TLS certificate verification")
	PushCmd.Flags().StringVarP(&pushUsername, "username", "u", "", "Registry username")
	PushCmd.Flags().StringVarP(&pushPassword, "password", "p", "", "Registry password")
}

func runPush(cmd *cobra.Command, args []string) error {
	imagePath := args[0]

	// Load the image from tarball
	img, err := tarball.ImageFromPath(imagePath, nil)
	if err != nil {
		return fmt.Errorf("failed to load image from %s: %w", imagePath, err)
	}

	// Parse the image reference
	ref, err := name.ParseReference(imagePath)
	if err != nil {
		// If parsing fails, use a default reference
		ref, err = name.ParseReference("localhost:5000/image:latest")
		if err != nil {
			return fmt.Errorf("failed to parse image reference: %w", err)
		}
	}

	// Create registry client
	baseURL := fmt.Sprintf("https://%s", ref.Context().Registry.Name())
	if pushInsecure {
		baseURL = fmt.Sprintf("http://%s", ref.Context().Registry.Name())
	}

	// Create a basic trust store manager (can be enhanced later)
	var trustStore certs.TrustStoreManager
	// trustStore can be nil for insecure connections

	client, err := registry.NewGiteaClient(baseURL, pushUsername, pushPassword, trustStore)
	if err != nil {
		return fmt.Errorf("failed to create registry client: %w", err)
	}

	// Push the image
	pushOpts := registry.PushOptions{
		Options: registry.Options{
			Insecure: pushSkipTLS,
		},
	}

	if err := client.Push(context.Background(), img, ref.String(), pushOpts); err != nil {
		return fmt.Errorf("failed to push image: %w", err)
	}

	fmt.Printf("Successfully pushed image to %s\n", ref.String())
	return nil
}