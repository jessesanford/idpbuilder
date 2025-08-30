package get

import (
	"github.com/cnoe-io/idpbuilder/pkg/cmd/helpers"
	"github.com/spf13/cobra"
)

var (
	// Common flags for get commands
	namespace     string
	allNamespaces bool
	selector      string
	outputFormat  string
	showLabels    bool
	wide          bool
)

// GetCmd represents the get command
var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get IDP resources and display information",
	Long: `Get and display information about IDP resources including clusters, packages, and secrets.

This command allows you to retrieve information about various IDP components:
- Clusters and their status
- Packages and applications
- Secrets and configurations
- Custom resources

Examples:
  # Get all clusters
  idpbuilder get clusters

  # Get packages in a specific namespace
  idpbuilder get packages -n development

  # Get resources with custom output format
  idpbuilder get secrets -o yaml

  # Get resources by label selector
  idpbuilder get all -l app=myapp
`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Validate output format
		if outputFormat != "" {
			if _, err := helpers.ValidateOutputFormat(outputFormat); err != nil {
				helpers.PrintError("Invalid output format: %v", err)
				return
			}
		}
	},
}

func init() {
	// Output flags
	GetCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table", "Output format (table, json, yaml, wide)")
	GetCmd.PersistentFlags().BoolVar(&showLabels, "show-labels", false, "Show labels in output")
	GetCmd.PersistentFlags().BoolVar(&wide, "wide", false, "Show additional columns in output")

	// Selection flags
	GetCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "Namespace scope (default: current namespace)")
	GetCmd.PersistentFlags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "List resources across all namespaces")
	GetCmd.PersistentFlags().StringVarP(&selector, "selector", "l", "", "Selector (label query) to filter resources")

	// Add subcommands
	GetCmd.AddCommand(ClustersCmd)
	GetCmd.AddCommand(PackagesCmd)
	GetCmd.AddCommand(SecretsCmd)

	// Register flag completions
	GetCmd.RegisterFlagCompletionFunc("output", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"table", "json", "yaml", "wide"}, cobra.ShellCompDirectiveNoFileComp
	})
}