package get

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cnoe-io/idpbuilder/pkg/cmd/helpers"
	"github.com/spf13/cobra"
)

// PackagesCmd represents the get packages command
var PackagesCmd = &cobra.Command{
	Use:     "packages [NAME...]",
	Aliases: []string{"package", "pkg", "pkgs"},
	Short:   "Get package information",
	Long: `Get information about IDP packages including their status, version, and configuration.

This command displays detailed information about installed packages:
- Package status and health
- Version and update information  
- Configuration and values
- Dependencies and requirements
- Resource usage

Examples:
  # Get all packages
  idpbuilder get packages

  # Get packages in a specific namespace
  idpbuilder get packages -n development

  # Get a specific package
  idpbuilder get packages nginx-ingress

  # Get packages with labels shown
  idpbuilder get packages --show-labels

  # Get packages in YAML format
  idpbuilder get packages -o yaml
`,
	RunE: runGetPackages,
}

// PackageInfo represents package information
type PackageInfo struct {
	Name        string            `json:"name" yaml:"name"`
	Namespace   string            `json:"namespace" yaml:"namespace"`
	Status      string            `json:"status" yaml:"status"`
	Version     string            `json:"version" yaml:"version"`
	Chart       string            `json:"chart,omitempty" yaml:"chart,omitempty"`
	AppVersion  string            `json:"appVersion,omitempty" yaml:"appVersion,omitempty"`
	Age         string            `json:"age" yaml:"age"`
	Labels      map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Description string            `json:"description,omitempty" yaml:"description,omitempty"`
}

// runGetPackages executes the get packages command
func runGetPackages(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	
	helpers.LogInfo("Getting package information")

	// Get package data
	packages, err := getPackages(ctx, args)
	if err != nil {
		return fmt.Errorf("failed to get packages: %w", err)
	}

	if len(packages) == 0 {
		if namespace != "" {
			helpers.PrintWarning("No packages found in namespace: %s", namespace)
		} else {
			helpers.PrintWarning("No packages found")
		}
		return nil
	}

	// Format and display output
	return displayPackages(packages)
}

// getPackages retrieves package information
func getPackages(ctx context.Context, packageNames []string) ([]PackageInfo, error) {
	var packages []PackageInfo

	// Simulate getting package information
	if len(packageNames) > 0 {
		// Get specific packages
		for _, name := range packageNames {
			pkg, err := getPackage(ctx, name, namespace)
			if err != nil {
				helpers.PrintError("Failed to get package %s: %v", name, err)
				continue
			}
			packages = append(packages, *pkg)
		}
	} else {
		// Get all packages
		allPackages, err := listAllPackages(ctx, namespace, allNamespaces)
		if err != nil {
			return nil, err
		}
		packages = allPackages
	}

	// Filter by selector if specified
	if selector != "" {
		packages = filterPackagesBySelector(packages, selector)
	}

	return packages, nil
}

// getPackage gets information for a specific package
func getPackage(ctx context.Context, name, ns string) (*PackageInfo, error) {
	helpers.LogDebug("Getting package: %s in namespace: %s", name, ns)

	// Simulate API call
	time.Sleep(30 * time.Millisecond)

	if ns == "" {
		ns = "default"
	}

	// Return mock package data
	pkg := &PackageInfo{
		Name:        name,
		Namespace:   ns,
		Status:      "Deployed",
		Version:     "1.0.0",
		Chart:       "stable/nginx",
		AppVersion:  "1.21.0",
		Age:         "2h",
		Description: fmt.Sprintf("Package %s deployed in %s namespace", name, ns),
		Labels: map[string]string{
			"app":     name,
			"version": "1.0.0",
		},
	}

	return pkg, nil
}

// listAllPackages lists all packages
func listAllPackages(ctx context.Context, ns string, allNS bool) ([]PackageInfo, error) {
	helpers.LogDebug("Listing packages - namespace: %s, all-namespaces: %v", ns, allNS)

	// Simulate API call
	time.Sleep(100 * time.Millisecond)

	// Return mock package list
	packages := []PackageInfo{
		{
			Name:       "nginx-ingress",
			Namespace:  "ingress-nginx",
			Status:     "Deployed", 
			Version:    "4.5.2",
			Chart:      "ingress-nginx/ingress-nginx",
			AppVersion: "1.7.0",
			Age:        "7d",
			Labels:     map[string]string{"component": "ingress"},
		},
		{
			Name:       "cert-manager",
			Namespace:  "cert-manager",
			Status:     "Deployed",
			Version:    "1.11.0",
			Chart:      "jetstack/cert-manager", 
			AppVersion: "v1.11.0",
			Age:        "5d",
			Labels:     map[string]string{"component": "security"},
		},
		{
			Name:       "prometheus",
			Namespace:  "monitoring",
			Status:     "Deployed",
			Version:    "19.7.2",
			Chart:      "prometheus-community/prometheus",
			AppVersion: "2.42.0",
			Age:        "3d",
			Labels:     map[string]string{"component": "monitoring"},
		},
	}

	// Filter by namespace if specified
	if ns != "" && !allNS {
		var filtered []PackageInfo
		for _, pkg := range packages {
			if pkg.Namespace == ns {
				filtered = append(filtered, pkg)
			}
		}
		packages = filtered
	}

	return packages, nil
}

// filterPackagesBySelector filters packages by label selector
func filterPackagesBySelector(packages []PackageInfo, selector string) []PackageInfo {
	var filtered []PackageInfo

	// Parse selector (simplified)
	parts := strings.Split(selector, "=")
	if len(parts) != 2 {
		helpers.PrintWarning("Invalid selector format: %s", selector)
		return packages
	}

	key, value := parts[0], parts[1]

	for _, pkg := range packages {
		if pkg.Labels != nil {
			if labelValue, exists := pkg.Labels[key]; exists && labelValue == value {
				filtered = append(filtered, pkg)
			}
		}
	}

	return filtered
}

// displayPackages formats and displays package information
func displayPackages(packages []PackageInfo) error {
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
	for _, pkg := range packages {
		item := map[string]interface{}{
			"name":      pkg.Name,
			"namespace": pkg.Namespace,
			"status":    pkg.Status,
			"version":   pkg.Version,
			"age":       pkg.Age,
		}

		if wide || format == helpers.WideOutput {
			item["chart"] = pkg.Chart
			item["appVersion"] = pkg.AppVersion
			item["description"] = pkg.Description
		}

		if showLabels && pkg.Labels != nil {
			var labels []string
			for k, v := range pkg.Labels {
				labels = append(labels, fmt.Sprintf("%s=%s", k, v))
			}
			item["labels"] = strings.Join(labels, ",")
		}

		data = append(data, item)
	}

	return printer.Print(data)
}