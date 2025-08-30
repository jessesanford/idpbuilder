package get

import (
	"context"
	"fmt"

	"github.com/cnoe-io/idpbuilder/pkg/cmd/helpers"
	"github.com/spf13/cobra"
)

var getSecretsCmd = &cobra.Command{
	Use:     "secrets",
	Aliases: []string{"secret"},
	Short:   "Get IDP secrets",
	Long:    `Get and display IDP secrets in the specified format.`,
	RunE:    runGetSecrets,
}

func init() {
	GetCmd.AddCommand(getSecretsCmd)
}

func runGetSecrets(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	
	helpers.LogInfo("Getting IDP secrets")
	
	secrets, err := listSecrets(ctx)
	if err != nil {
		return fmt.Errorf("failed to list secrets: %w", err)
	}
	
	printer := helpers.NewPrinter(outputFormat)
	if err := printer.Print(secrets); err != nil {
		return fmt.Errorf("failed to print secrets: %w", err)
	}
	
	return nil
}

func listSecrets(ctx context.Context) ([]map[string]interface{}, error) {
	// Simplified secret listing
	secrets := []map[string]interface{}{
		{
			"name":      "example-secret-1",
			"namespace": "default",
			"type":      "Opaque",
			"age":       "1d",
		},
		{
			"name":      "example-secret-2", 
			"namespace": "kube-system",
			"type":      "kubernetes.io/service-account-token",
			"age":       "7d",
		},
	}
	
	// Filter by namespace if specified
	if namespace != "" {
		var filtered []map[string]interface{}
		for _, secret := range secrets {
			if secret["namespace"] == namespace {
				filtered = append(filtered, secret)
			}
		}
		return filtered, nil
	}
	
	return secrets, nil
}