package delete

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cnoe-io/idpbuilder/pkg/cmd/helpers"
	"github.com/spf13/cobra"
)

var (
	// Command-specific flags
	force           bool
	wait            bool
	timeout         time.Duration
	namespace       string
	selector        string
	allNamespaces   bool
	gracePeriod     int
	cascade         string
)

// DeleteCmd represents the delete command
var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete IDP resources and configurations",
	Long: `Delete IDP resources and configurations by name, selector, or configuration file.

This command allows you to delete various IDP components including:
- Individual resources by name
- Multiple resources by label selector
- All resources in a namespace
- Resources defined in configuration files

Examples:
  # Delete a specific package
  idpbuilder delete package my-package

  # Delete all packages in a namespace
  idpbuilder delete packages --all -n development

  # Delete resources by label selector
  idpbuilder delete all -l app=myapp

  # Delete with confirmation bypass
  idpbuilder delete package my-package --force

  # Delete and wait for completion
  idpbuilder delete deployment my-app --wait --timeout=2m
`,
	Args:              cobra.MinimumNArgs(1),
	ValidArgsFunction: validResourceTypes,
	RunE:              runDelete,
}

var validResourceTypesList = []string{
	"package", "packages",
	"secret", "secrets", 
	"config", "configs",
	"application", "applications", "apps",
	"all",
}

func init() {
	// Target selection flags
	DeleteCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Target namespace (default: current namespace)")
	DeleteCmd.Flags().StringVarP(&selector, "selector", "l", "", "Selector (label query) to filter resources")
	DeleteCmd.Flags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "Delete resources across all namespaces")

	// Behavior flags
	DeleteCmd.Flags().BoolVar(&force, "force", false, "Skip confirmation prompts")
	DeleteCmd.Flags().BoolVar(&wait, "wait", false, "Wait for resources to be fully deleted")
	DeleteCmd.Flags().DurationVar(&timeout, "timeout", 2*time.Minute, "Timeout for wait operations")
	DeleteCmd.Flags().IntVar(&gracePeriod, "grace-period", -1, "Grace period in seconds for deletion")
	DeleteCmd.Flags().StringVar(&cascade, "cascade", "background", "Cascade deletion policy (background, foreground, orphan)")

	// Validation
	DeleteCmd.RegisterFlagCompletionFunc("cascade", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"background", "foreground", "orphan"}, cobra.ShellCompDirectiveNoFileComp
	})
}

// validResourceTypes provides completion for resource types
func validResourceTypes(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return validResourceTypesList, cobra.ShellCompDirectiveNoFileComp
	}
	return nil, cobra.ShellCompDirectiveNoFileComp
}

// runDelete executes the delete command
func runDelete(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	
	resourceType := args[0]
	resourceNames := args[1:]
	
	helpers.LogInfo("Starting IDP resource deletion")
	
	// Validate inputs
	if err := validateDeleteInputs(resourceType, resourceNames); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Build deletion request
	request := &DeleteRequest{
		ResourceType:  resourceType,
		ResourceNames: resourceNames,
		Namespace:     namespace,
		Selector:      selector,
		AllNamespaces: allNamespaces,
		Force:         force,
		Wait:          wait,
		Timeout:       timeout,
		GracePeriod:   gracePeriod,
		Cascade:       cascade,
	}

	// Confirm deletion unless force is specified
	if !force {
		confirmed, err := confirmDeletion(request)
		if err != nil {
			return fmt.Errorf("failed to get confirmation: %w", err)
		}
		if !confirmed {
			helpers.PrintWarning("Deletion cancelled by user")
			return nil
		}
	}

	// Execute deletion
	if err := executeDelete(ctx, request); err != nil {
		return fmt.Errorf("deletion failed: %w", err)
	}

	helpers.PrintSuccess("IDP resources deleted successfully")
	return nil
}

// DeleteRequest represents a deletion request
type DeleteRequest struct {
	ResourceType  string
	ResourceNames []string
	Namespace     string
	Selector      string
	AllNamespaces bool
	Force         bool
	Wait          bool
	Timeout       time.Duration
	GracePeriod   int
	Cascade       string
}

// validateDeleteInputs validates the delete command inputs
func validateDeleteInputs(resourceType string, resourceNames []string) error {
	// Validate resource type
	isValidType := false
	for _, validType := range validResourceTypesList {
		if resourceType == validType {
			isValidType = true
			break
		}
	}
	if !isValidType {
		return fmt.Errorf("invalid resource type: %s (valid types: %v)", resourceType, validResourceTypesList)
	}

	// Validate namespace if specified
	if namespace != "" {
		if err := helpers.ValidateNamespace(namespace); err != nil {
			return err
		}
	}

	// Validate resource names
	for _, name := range resourceNames {
		if err := helpers.ValidateName(name); err != nil {
			return fmt.Errorf("invalid resource name '%s': %w", name, err)
		}
	}

	// Validate cascade policy
	validCascadePolicies := []string{"background", "foreground", "orphan"}
	isValidCascade := false
	for _, policy := range validCascadePolicies {
		if cascade == policy {
			isValidCascade = true
			break
		}
	}
	if !isValidCascade {
		return fmt.Errorf("invalid cascade policy: %s (valid: %v)", cascade, validCascadePolicies)
	}

	// Validate timeout
	if timeout <= 0 {
		return fmt.Errorf("timeout must be greater than 0")
	}

	return nil
}

// confirmDeletion asks for user confirmation before deletion
func confirmDeletion(request *DeleteRequest) (bool, error) {
	// Build confirmation message
	var message strings.Builder
	message.WriteString(fmt.Sprintf("Are you sure you want to delete %s", request.ResourceType))
	
	if len(request.ResourceNames) > 0 {
		message.WriteString(fmt.Sprintf(" '%s'", strings.Join(request.ResourceNames, "', '")))
	}
	
	if request.Namespace != "" {
		message.WriteString(fmt.Sprintf(" in namespace '%s'", request.Namespace))
	} else if request.AllNamespaces {
		message.WriteString(" across all namespaces")
	}
	
	if request.Selector != "" {
		message.WriteString(fmt.Sprintf(" matching selector '%s'", request.Selector))
	}
	
	message.WriteString("? [y/N]: ")

	// Show warning for dangerous operations
	if request.ResourceType == "all" || request.AllNamespaces {
		helpers.PrintWarning("This will delete multiple resources!")
	}

	fmt.Print(message.String())
	
	var response string
	fmt.Scanln(&response)
	
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes", nil
}

// executeDelete performs the actual deletion
func executeDelete(ctx context.Context, request *DeleteRequest) error {
	helpers.PrintStep("DELETE", "Starting resource deletion")

	// Find resources to delete
	resources, err := findResources(ctx, request)
	if err != nil {
		return fmt.Errorf("failed to find resources: %w", err)
	}

	if len(resources) == 0 {
		helpers.PrintWarning("No resources found matching criteria")
		return nil
	}

	// Delete each resource
	for _, resource := range resources {
		if err := deleteResource(ctx, resource, request); err != nil {
			helpers.PrintError("Failed to delete %s/%s: %v", resource.Type, resource.Name, err)
			continue
		}
		helpers.LogInfo("Deleted %s: %s", resource.Type, resource.Name)
	}

	// Wait for deletion if requested
	if request.Wait {
		helpers.PrintStep("WAIT", "Waiting for resources to be fully deleted...")
		return waitForDeletion(ctx, resources, request.Timeout)
	}

	return nil
}

// ResourceRef represents a resource reference
type ResourceRef struct {
	Type      string
	Name      string
	Namespace string
}

// findResources finds resources matching the deletion criteria
func findResources(ctx context.Context, request *DeleteRequest) ([]ResourceRef, error) {
	var resources []ResourceRef
	
	// Simulate finding resources based on criteria
	if len(request.ResourceNames) > 0 {
		// Find specific named resources
		for _, name := range request.ResourceNames {
			resources = append(resources, ResourceRef{
				Type:      request.ResourceType,
				Name:      name,
				Namespace: request.Namespace,
			})
		}
	} else {
		// Find resources by selector or all resources
		// This would query the actual backend in a real implementation
		helpers.LogDebug("Finding resources by selector: %s", request.Selector)
	}

	return resources, nil
}

// deleteResource deletes a single resource
func deleteResource(ctx context.Context, resource ResourceRef, request *DeleteRequest) error {
	helpers.LogDebug("Deleting %s: %s in namespace %s", resource.Type, resource.Name, resource.Namespace)
	
	// Simulate resource deletion
	time.Sleep(100 * time.Millisecond)
	
	return nil
}

// waitForDeletion waits for resources to be fully deleted
func waitForDeletion(ctx context.Context, resources []ResourceRef, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Simulate waiting for deletion
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for resource deletion")
		case <-ticker.C:
			// Check if resources still exist (simulated)
			helpers.LogDebug("Checking deletion status...")
			// In real implementation, check actual resource existence
			return nil // Assume deleted for simulation
		}
	}
}