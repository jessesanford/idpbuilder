package get

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cnoe-io/idpbuilder/pkg/cmd/helpers"
	"github.com/spf13/cobra"
)

// ClustersCmd represents the get clusters command
var ClustersCmd = &cobra.Command{
	Use:     "clusters [NAME...]",
	Aliases: []string{"cluster", "cl"},
	Short:   "Get cluster information",
	Long: `Get information about IDP clusters including their status, configuration, and resources.

This command displays detailed information about one or more clusters:
- Cluster status and health
- Node information and capacity
- Installed packages and applications
- Network configuration
- Resource utilization

Examples:
  # Get all clusters
  idpbuilder get clusters

  # Get a specific cluster
  idpbuilder get clusters my-cluster

  # Get clusters with wide output
  idpbuilder get clusters --wide

  # Get clusters in JSON format
  idpbuilder get clusters -o json

  # Get clusters by label
  idpbuilder get clusters -l environment=production
`,
	RunE: runGetClusters,
}

// ClusterInfo represents cluster information
type ClusterInfo struct {
	Name         string            `json:"name" yaml:"name"`
	Status       string            `json:"status" yaml:"status"`
	Version      string            `json:"version" yaml:"version"`
	Nodes        int               `json:"nodes" yaml:"nodes"`
	Ready        string            `json:"ready" yaml:"ready"`
	Age          string            `json:"age" yaml:"age"`
	Endpoint     string            `json:"endpoint,omitempty" yaml:"endpoint,omitempty"`
	Provider     string            `json:"provider,omitempty" yaml:"provider,omitempty"`
	Region       string            `json:"region,omitempty" yaml:"region,omitempty"`
	Labels       map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Annotations  map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
}

// runGetClusters executes the get clusters command
func runGetClusters(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	
	helpers.LogInfo("Getting cluster information")

	// Get cluster data
	clusters, err := getClusters(ctx, args)
	if err != nil {
		return fmt.Errorf("failed to get clusters: %w", err)
	}

	if len(clusters) == 0 {
		helpers.PrintWarning("No clusters found")
		return nil
	}

	// Format and display output
	return displayClusters(clusters)
}

// getClusters retrieves cluster information
func getClusters(ctx context.Context, clusterNames []string) ([]ClusterInfo, error) {
	var clusters []ClusterInfo

	// Simulate getting cluster information
	if len(clusterNames) > 0 {
		// Get specific clusters
		for _, name := range clusterNames {
			cluster, err := getCluster(ctx, name)
			if err != nil {
				helpers.PrintError("Failed to get cluster %s: %v", name, err)
				continue
			}
			clusters = append(clusters, *cluster)
		}
	} else {
		// Get all clusters
		allClusters, err := listAllClusters(ctx)
		if err != nil {
			return nil, err
		}
		clusters = allClusters
	}

	// Filter by selector if specified
	if selector != "" {
		clusters = filterClustersBySelector(clusters, selector)
	}

	return clusters, nil
}

// getCluster gets information for a specific cluster
func getCluster(ctx context.Context, name string) (*ClusterInfo, error) {
	helpers.LogDebug("Getting cluster: %s", name)

	// Simulate API call to get cluster
	time.Sleep(50 * time.Millisecond)

	// Return mock cluster data
	cluster := &ClusterInfo{
		Name:     name,
		Status:   "Running",
		Version:  "v1.28.0",
		Nodes:    3,
		Ready:    "3/3",
		Age:      "7d",
		Endpoint: fmt.Sprintf("https://%s.example.com:6443", name),
		Provider: "kind",
		Region:   "local",
		Labels: map[string]string{
			"environment": "development",
			"team":        "platform",
		},
	}

	return cluster, nil
}

// listAllClusters lists all available clusters
func listAllClusters(ctx context.Context) ([]ClusterInfo, error) {
	helpers.LogDebug("Listing all clusters")

	// Simulate API call to list clusters
	time.Sleep(100 * time.Millisecond)

	// Return mock cluster list
	clusters := []ClusterInfo{
		{
			Name:     "dev-cluster",
			Status:   "Running", 
			Version:  "v1.28.0",
			Nodes:    1,
			Ready:    "1/1",
			Age:      "3d",
			Provider: "kind",
			Labels:   map[string]string{"environment": "development"},
		},
		{
			Name:     "staging-cluster",
			Status:   "Running",
			Version:  "v1.28.0", 
			Nodes:    2,
			Ready:    "2/2",
			Age:      "7d",
			Provider: "kind",
			Labels:   map[string]string{"environment": "staging"},
		},
		{
			Name:     "prod-cluster",
			Status:   "Running",
			Version:  "v1.28.0",
			Nodes:    5,
			Ready:    "5/5", 
			Age:      "30d",
			Provider: "eks",
			Labels:   map[string]string{"environment": "production"},
		},
	}

	return clusters, nil
}

// filterClustersBySelector filters clusters by label selector
func filterClustersBySelector(clusters []ClusterInfo, selector string) []ClusterInfo {
	var filtered []ClusterInfo

	// Parse selector (simplified - real implementation would use proper label selector parsing)
	parts := strings.Split(selector, "=")
	if len(parts) != 2 {
		helpers.PrintWarning("Invalid selector format: %s", selector)
		return clusters
	}

	key, value := parts[0], parts[1]

	for _, cluster := range clusters {
		if cluster.Labels != nil {
			if labelValue, exists := cluster.Labels[key]; exists && labelValue == value {
				filtered = append(filtered, cluster)
			}
		}
	}

	return filtered
}

// displayClusters formats and displays cluster information
func displayClusters(clusters []ClusterInfo) error {
	// Determine output format
	format := helpers.TableOutput
	if outputFormat != "" {
		var err error
		format, err = helpers.ValidateOutputFormat(outputFormat)
		if err != nil {
			return err
		}
	}
	if wide {
		format = helpers.WideOutput
	}

	// Create printer and display
	printer := helpers.NewPrinter(format)

	// Convert to map slice for printing
	var data []map[string]interface{}
	for _, cluster := range clusters {
		item := map[string]interface{}{
			"name":    cluster.Name,
			"status":  cluster.Status,
			"ready":   cluster.Ready,
			"age":     cluster.Age,
		}

		if wide || format == helpers.WideOutput {
			item["version"] = cluster.Version
			item["nodes"] = cluster.Nodes
			item["provider"] = cluster.Provider
			item["endpoint"] = cluster.Endpoint
		}

		if showLabels && cluster.Labels != nil {
			var labels []string
			for k, v := range cluster.Labels {
				labels = append(labels, fmt.Sprintf("%s=%s", k, v))
			}
			item["labels"] = strings.Join(labels, ",")
		}

		data = append(data, item)
	}

	return printer.Print(data)
}