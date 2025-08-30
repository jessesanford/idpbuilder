package get

import (
	"context"
	"fmt"

	"github.com/cnoe-io/idpbuilder/pkg/cmd/helpers"
	"github.com/spf13/cobra"
)

var getClustersCmd = &cobra.Command{
	Use:     "clusters",
	Aliases: []string{"cluster"},
	Short:   "Get IDP clusters",
	Long:    `Get and display IDP clusters in the specified format.`,
	RunE:    runGetClusters,
}

func init() {
	GetCmd.AddCommand(getClustersCmd)
}

func runGetClusters(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	
	helpers.LogInfo("Getting IDP clusters")
	
	clusters, err := listClusters(ctx)
	if err != nil {
		return fmt.Errorf("failed to list clusters: %w", err)
	}
	
	printer := helpers.NewPrinter(outputFormat)
	if err := printer.Print(clusters); err != nil {
		return fmt.Errorf("failed to print clusters: %w", err)
	}
	
	return nil
}

func listClusters(ctx context.Context) ([]map[string]interface{}, error) {
	// Simplified cluster listing
	clusters := []map[string]interface{}{
		{
			"name":    "local-cluster",
			"status":  "Ready",
			"version": "v1.28.0",
			"nodes":   "3",
			"age":     "5d",
		},
		{
			"name":    "dev-cluster",
			"status":  "NotReady",
			"version": "v1.27.3",
			"nodes":   "1",
			"age":     "2d",
		},
	}
	
	return clusters, nil
}