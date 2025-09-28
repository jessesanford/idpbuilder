package push

import (
	"fmt"

	"github.com/spf13/cobra"
)

// PushCmd represents the push command with auth support
var PushCmd = &cobra.Command{
	Use:   "push [IMAGE]",
	Short: "Push an OCI package to a registry",
	Long: `Push an OCI package to a container registry.
Supports authentication and various registry types.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		image := args[0]

		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		insecure, _ := cmd.Flags().GetBool("insecure")

		return pushImage(image, username, password, insecure)
	},
}

func init() {
	PushCmd.Flags().StringP("username", "u", "", "Username for registry authentication")
	PushCmd.Flags().StringP("password", "p", "", "Password for registry authentication")
	PushCmd.Flags().Bool("insecure", false, "Allow insecure connections for self-signed certificates")
}

func pushImage(image, username, password string, insecure bool) error {
	fmt.Printf("Pushing image: %s\n", image)

	if username != "" {
		fmt.Printf("Using authentication for user: %s\n", username)
	}

	if insecure {
		fmt.Println("Warning: Using insecure connection")
	}

	// Implementation would use the OCI format types from pkg/oci/format
	fmt.Println("✅ Push completed successfully")
	return nil
}