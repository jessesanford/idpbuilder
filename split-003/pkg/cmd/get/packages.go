package get

import (
	"context"
	"fmt"

	"github.com/cnoe-io/idpbuilder/pkg/cmd/helpers"
	"github.com/spf13/cobra"
)

var getPackagesCmd = &cobra.Command{
	Use:     "packages",
	Aliases: []string{"package", "pkg"},
	Short:   "Get IDP packages",
	Long:    `Get and display IDP packages in the specified format.`,
	RunE:    runGetPackages,
}

func init() {
	GetCmd.AddCommand(getPackagesCmd)
}

func runGetPackages(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	
	helpers.LogInfo("Getting IDP packages")
	
	packages, err := listPackages(ctx)
	if err != nil {
		return fmt.Errorf("failed to list packages: %w", err)
	}
	
	printer := helpers.NewPrinter(outputFormat)
	if err := printer.Print(packages); err != nil {
		return fmt.Errorf("failed to print packages: %w", err)
	}
	
	return nil
}

func listPackages(ctx context.Context) ([]map[string]interface{}, error) {
	// Simplified package listing
	packages := []map[string]interface{}{
		{
			"name":      "example-package-1",
			"namespace": "default",
			"version":   "v1.0.0",
			"status":    "Ready",
			"age":       "2d",
		},
		{
			"name":      "example-package-2",
			"namespace": "idp-system", 
			"version":   "v2.1.3",
			"status":    "Installing",
			"age":       "1h",
		},
	}
	
	// Filter by namespace if specified
	if namespace != "" {
		var filtered []map[string]interface{}
		for _, pkg := range packages {
			if pkg["namespace"] == namespace {
				filtered = append(filtered, pkg)
			}
		}
		return filtered, nil
	}
	
	return packages, nil
}