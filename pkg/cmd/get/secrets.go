// pkg/cmd/get/secrets.go
package get

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/cnoe-io/idpbuilder/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	argoCDInitialAdminSecretName = "argocd-initial-admin-secret"
	giteaAdminSecretName         = "gitea-admin-secret"
)

var packages = []string{"argocd", "gitea"}

// SecretData represents the structure for secret output
type SecretData struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Data      map[string]string `json:"data,omitempty"`
	Username  string            `json:"username,omitempty"`
	Password  string            `json:"password,omitempty"`
	IsCore    bool              `json:"isCore,omitempty"`
}

// printPackageSecrets prints secrets for specified packages
func printPackageSecrets(ctx context.Context, w io.Writer, kubeClient client.Client, outputFormat string) error {
	var allSecrets []SecretData

	// Handle packages based on the global packages variable
	for _, pkg := range packages {
		var secretName string
		var namespace string

		switch pkg {
		case "argocd":
			secretName = argoCDInitialAdminSecretName
			namespace = "argocd"
			secret := &v1.Secret{}
			key := client.ObjectKey{Name: secretName, Namespace: namespace}
			err := kubeClient.Get(ctx, key, secret, []client.GetOption{}...)
			if err != nil {
				continue // Skip if secret doesn't exist
			}

			secretData := SecretData{
				Name:      secret.Name,
				Namespace: secret.Namespace,
				IsCore:    true,
			}

			// Extract username and password if available
			if data, ok := secret.Data["username"]; ok {
				secretData.Username = string(data)
			}
			if data, ok := secret.Data["password"]; ok {
				secretData.Password = string(data)
			}

			allSecrets = append(allSecrets, secretData)

		case "gitea":
			secretName = giteaAdminSecretName
			namespace = "gitea"
			secret := &v1.Secret{}
			key := client.ObjectKey{Name: secretName, Namespace: namespace}
			err := kubeClient.Get(ctx, key, secret, []client.GetOption{}...)
			if err != nil {
				continue // Skip if secret doesn't exist
			}

			secretData := SecretData{
				Name:      secret.Name,
				Namespace: secret.Namespace,
				IsCore:    true,
			}

			// Extract username and password if available
			if data, ok := secret.Data["username"]; ok {
				secretData.Username = string(data)
			}
			if data, ok := secret.Data["password"]; ok {
				secretData.Password = string(data)
			}

			allSecrets = append(allSecrets, secretData)

		default:
			// Handle non-core packages - create label selector
			req, err := labels.NewRequirement(v1alpha1.CLISecretLabelKey, selection.Equals, []string{pkg})
			if err != nil {
				continue
			}
			selector := labels.NewSelector().Add(*req)

			secretList := &v1.SecretList{}
			listOpts := client.ListOptions{
				LabelSelector: selector,
			}

			err = kubeClient.List(ctx, secretList, &listOpts)
			if err != nil {
				continue
			}

			for _, secret := range secretList.Items {
				secretData := SecretData{
					Name:      secret.Name,
					Namespace: secret.Namespace,
					Data:      make(map[string]string),
				}

				// Convert byte data to strings
				for k, v := range secret.Data {
					secretData.Data[k] = string(v)
				}

				allSecrets = append(allSecrets, secretData)
			}
		}
	}

	return outputSecrets(w, allSecrets, outputFormat)
}

// printAllPackageSecrets prints all package secrets
func printAllPackageSecrets(ctx context.Context, w io.Writer, kubeClient client.Client, outputFormat string) error {
	var allSecrets []SecretData

	// Get core package secrets
	coreSecrets := map[string]string{
		"argocd": argoCDInitialAdminSecretName,
		"gitea":  giteaAdminSecretName,
	}

	for namespace, secretName := range coreSecrets {
		secret := &v1.Secret{}
		key := client.ObjectKey{Name: secretName, Namespace: namespace}
		err := kubeClient.Get(ctx, key, secret)
		if err != nil {
			continue // Skip if secret doesn't exist
		}

		secretData := SecretData{
			Name:      secret.Name,
			Namespace: secret.Namespace,
			IsCore:    true,
		}

		// Extract username and password if available
		if data, ok := secret.Data["username"]; ok {
			secretData.Username = string(data)
		}
		if data, ok := secret.Data["password"]; ok {
			secretData.Password = string(data)
		}

		allSecrets = append(allSecrets, secretData)
	}

	// Get all labeled secrets
	req, err := labels.NewRequirement(v1alpha1.CLISecretLabelKey, selection.Equals, []string{v1alpha1.CLISecretLabelValue})
	if err != nil {
		return err
	}
	selector := labels.NewSelector().Add(*req)

	secretList := &v1.SecretList{}
	listOpts := client.ListOptions{
		LabelSelector: selector,
	}

	err = kubeClient.List(ctx, secretList, &listOpts)
	if err != nil {
		return err
	}

	for _, secret := range secretList.Items {
		secretData := SecretData{
			Name:      secret.Name,
			Namespace: secret.Namespace,
			Data:      make(map[string]string),
		}

		// Convert byte data to strings
		for k, v := range secret.Data {
			secretData.Data[k] = string(v)
		}

		allSecrets = append(allSecrets, secretData)
	}

	return outputSecrets(w, allSecrets, outputFormat)
}

// outputSecrets formats and outputs the secrets data
func outputSecrets(w io.Writer, secrets []SecretData, outputFormat string) error {
	switch outputFormat {
	case "json":
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		return encoder.Encode(secrets)
	case "yaml":
		// For simplicity, using JSON format for YAML as well
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		return encoder.Encode(secrets)
	case "table":
		fallthrough
	default:
		return outputSecretsTable(w, secrets)
	}
}

// outputSecretsTable outputs secrets in table format
func outputSecretsTable(w io.Writer, secrets []SecretData) error {
	writer := tabwriter.NewWriter(w, 0, 8, 2, ' ', 0)
	defer writer.Flush()

	// Write header
	fmt.Fprintln(writer, "NAME\tNAMESPACE\tUSERNAME\tPASSWORD\tCORE")

	for _, secret := range secrets {
		username := secret.Username
		password := secret.Password

		// For non-core secrets, show data keys
		if !secret.IsCore && len(secret.Data) > 0 {
			username = ""
			password = ""
			for k, v := range secret.Data {
				if k == "username" {
					username = v
				} else if k == "password" {
					password = v
				}
			}
		}

		core := "false"
		if secret.IsCore {
			core = "true"
		}

		fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\n",
			secret.Name, secret.Namespace, username, password, core)
	}

	return nil
}