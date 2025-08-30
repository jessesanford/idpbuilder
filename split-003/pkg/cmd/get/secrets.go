package get

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cnoe-io/idpbuilder/pkg/cmd/helpers"
	"github.com/spf13/cobra"
)

// SecretsCmd represents the get secrets command
var SecretsCmd = &cobra.Command{
	Use:     "secrets [NAME...]",
	Aliases: []string{"secret", "sec"},
	Short:   "Get secret information",
	Long: `Get information about IDP secrets including their type, data keys, and metadata.

This command displays information about secrets without revealing sensitive data:
- Secret names and types
- Data keys (not values)
- Metadata and labels
- Creation and modification timestamps
- Associated resources

Examples:
  # Get all secrets
  idpbuilder get secrets

  # Get secrets in a specific namespace
  idpbuilder get secrets -n kube-system

  # Get a specific secret
  idpbuilder get secrets my-secret

  # Get secrets across all namespaces
  idpbuilder get secrets --all-namespaces

  # Get secrets by type
  idpbuilder get secrets -l type=kubernetes.io/tls
`,
	RunE: runGetSecrets,
}

// SecretInfo represents secret information (without sensitive data)
type SecretInfo struct {
	Name        string            `json:"name" yaml:"name"`
	Namespace   string            `json:"namespace" yaml:"namespace"`
	Type        string            `json:"type" yaml:"type"`
	DataKeys    []string          `json:"dataKeys" yaml:"dataKeys"`
	Age         string            `json:"age" yaml:"age"`
	Labels      map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
}

// runGetSecrets executes the get secrets command
func runGetSecrets(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	
	helpers.LogInfo("Getting secret information")

	// Get secret data
	secrets, err := getSecrets(ctx, args)
	if err != nil {
		return fmt.Errorf("failed to get secrets: %w", err)
	}

	if len(secrets) == 0 {
		if namespace != "" {
			helpers.PrintWarning("No secrets found in namespace: %s", namespace)
		} else {
			helpers.PrintWarning("No secrets found")
		}
		return nil
	}

	// Format and display output
	return displaySecrets(secrets)
}

// getSecrets retrieves secret information
func getSecrets(ctx context.Context, secretNames []string) ([]SecretInfo, error) {
	var secrets []SecretInfo

	// Simulate getting secret information
	if len(secretNames) > 0 {
		// Get specific secrets
		for _, name := range secretNames {
			secret, err := getSecret(ctx, name, namespace)
			if err != nil {
				helpers.PrintError("Failed to get secret %s: %v", name, err)
				continue
			}
			secrets = append(secrets, *secret)
		}
	} else {
		// Get all secrets
		allSecrets, err := listAllSecrets(ctx, namespace, allNamespaces)
		if err != nil {
			return nil, err
		}
		secrets = allSecrets
	}

	// Filter by selector if specified
	if selector != "" {
		secrets = filterSecretsBySelector(secrets, selector)
	}

	return secrets, nil
}

// getSecret gets information for a specific secret
func getSecret(ctx context.Context, name, ns string) (*SecretInfo, error) {
	helpers.LogDebug("Getting secret: %s in namespace: %s", name, ns)

	// Simulate API call
	time.Sleep(30 * time.Millisecond)

	if ns == "" {
		ns = "default"
	}

	// Return mock secret data (without actual secret values)
	secret := &SecretInfo{
		Name:      name,
		Namespace: ns,
		Type:      "Opaque",
		DataKeys:  []string{"username", "password"},
		Age:       "1h",
		Labels: map[string]string{
			"app": "myapp",
		},
		Annotations: map[string]string{
			"created-by": "idpbuilder",
		},
	}

	return secret, nil
}

// listAllSecrets lists all secrets
func listAllSecrets(ctx context.Context, ns string, allNS bool) ([]SecretInfo, error) {
	helpers.LogDebug("Listing secrets - namespace: %s, all-namespaces: %v", ns, allNS)

	// Simulate API call
	time.Sleep(100 * time.Millisecond)

	// Return mock secret list
	secrets := []SecretInfo{
		{
			Name:      "database-credentials",
			Namespace: "default",
			Type:      "Opaque",
			DataKeys:  []string{"username", "password", "host"},
			Age:       "7d",
			Labels:    map[string]string{"component": "database"},
		},
		{
			Name:      "api-tls-cert",
			Namespace: "default",
			Type:      "kubernetes.io/tls",
			DataKeys:  []string{"tls.crt", "tls.key"},
			Age:       "30d",
			Labels:    map[string]string{"component": "api", "type": "tls"},
		},
		{
			Name:      "docker-registry-secret",
			Namespace: "kube-system",
			Type:      "kubernetes.io/dockerconfigjson",
			DataKeys:  []string{".dockerconfigjson"},
			Age:       "14d",
			Labels:    map[string]string{"component": "registry"},
		},
		{
			Name:      "service-account-token",
			Namespace: "kube-system",
			Type:      "kubernetes.io/service-account-token",
			DataKeys:  []string{"ca.crt", "namespace", "token"},
			Age:       "60d",
			Labels:    map[string]string{"component": "auth"},
		},
	}

	// Filter by namespace if specified
	if ns != "" && !allNS {
		var filtered []SecretInfo
		for _, secret := range secrets {
			if secret.Namespace == ns {
				filtered = append(filtered, secret)
			}
		}
		secrets = filtered
	}

	return secrets, nil
}

// filterSecretsBySelector filters secrets by label selector
func filterSecretsBySelector(secrets []SecretInfo, selector string) []SecretInfo {
	var filtered []SecretInfo

	// Parse selector (simplified)
	parts := strings.Split(selector, "=")
	if len(parts) != 2 {
		helpers.PrintWarning("Invalid selector format: %s", selector)
		return secrets
	}

	key, value := parts[0], parts[1]

	for _, secret := range secrets {
		if secret.Labels != nil {
			if labelValue, exists := secret.Labels[key]; exists && labelValue == value {
				filtered = append(filtered, secret)
			}
		}
	}

	return filtered
}

// displaySecrets formats and displays secret information
func displaySecrets(secrets []SecretInfo) error {
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
	for _, secret := range secrets {
		item := map[string]interface{}{
			"name":      secret.Name,
			"namespace": secret.Namespace,
			"type":      secret.Type,
			"data":      fmt.Sprintf("%d keys", len(secret.DataKeys)),
			"age":       secret.Age,
		}

		if wide || format == helpers.WideOutput {
			item["dataKeys"] = strings.Join(secret.DataKeys, ",")
			if secret.Annotations != nil {
				item["annotations"] = fmt.Sprintf("%d", len(secret.Annotations))
			}
		}

		if showLabels && secret.Labels != nil {
			var labels []string
			for k, v := range secret.Labels {
				labels = append(labels, fmt.Sprintf("%s=%s", k, v))
			}
			item["labels"] = strings.Join(labels, ",")
		}

		data = append(data, item)
	}

	return printer.Print(data)
}