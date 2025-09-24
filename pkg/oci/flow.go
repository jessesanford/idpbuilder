package oci

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// AuthFlowConfig configures the authentication flow
type AuthFlowConfig struct {
	Username  string                // From flags
	Password  string                // From flags
	K8sClient kubernetes.Interface  // Kubernetes client for secrets
	Logger    logr.Logger          // Logger for debug messages
}

// AuthFlow manages credential precedence (flags override secrets)
type AuthFlow struct {
	config        *AuthFlowConfig
	authenticator Authenticator
	k8sClient     kubernetes.Interface
	logger        logr.Logger
	flagCreds     *Credentials
}

// FlowCredentials represents credentials with source information
type FlowCredentials struct {
	*Credentials
	Source string // Source description (flags, secrets, etc.)
	Valid  bool   // Whether credentials are valid
}

// NewAuthFlow creates a new authentication flow with the given configuration
func NewAuthFlow(config *AuthFlowConfig) *AuthFlow {
	flow := &AuthFlow{
		config:    config,
		k8sClient: config.K8sClient,
		logger:    config.Logger,
	}

	// Pre-populate flag credentials if provided
	if config.Username != "" && config.Password != "" {
		flow.flagCreds = &Credentials{
			Username: config.Username,
			Password: config.Password,
		}
	}

	return flow
}

// GetCredentials returns credentials with proper precedence (flags override secrets)
func (f *AuthFlow) GetCredentials(ctx context.Context, registry string) (*Credentials, error) {
	f.logger.Info("Getting credentials for registry", "registry", registry)

	// Try flags first (highest precedence)
	flagCreds, err := f.getFromFlags()
	if err == nil && flagCreds != nil {
		f.logger.Info("Using credentials from flags", "username", flagCreds.Username, "source", "command-line flags")
		flagCreds.Registry = registry
		return flagCreds, nil
	}

	// Fall back to Kubernetes secrets
	secretCreds, err := f.getFromSecrets(ctx, registry)
	if err == nil && secretCreds != nil {
		f.logger.Info("Using credentials from Kubernetes secret", "username", secretCreds.Username, "source", "kubernetes-secret")
		return secretCreds, nil
	}

	// No credentials available
	f.logger.Error(nil, "No credentials available", "registry", registry, "tried", "flags (empty), secrets (not found)")
	return nil, fmt.Errorf("no credentials available for registry %s", registry)
}

// getFromFlags extracts credentials from command-line flags
func (f *AuthFlow) getFromFlags() (*Credentials, error) {
	if f.flagCreds == nil {
		return nil, errors.New("no flag credentials available")
	}

	// Validate flag credentials
	if err := f.validateCredentials(f.flagCreds); err != nil {
		return nil, fmt.Errorf("invalid flag credentials: %w", err)
	}

	return f.flagCreds, nil
}

// getFromSecrets retrieves credentials from Kubernetes secrets
func (f *AuthFlow) getFromSecrets(ctx context.Context, registry string) (*Credentials, error) {
	if f.k8sClient == nil {
		return nil, errors.New("no kubernetes client available")
	}

	// Try to get the secret from the default namespace
	// This follows idpbuilder's secret naming convention
	secretName := "registry-credentials"
	namespace := "default"

	secret, err := f.k8sClient.CoreV1().Secrets(namespace).Get(ctx, secretName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get secret %s/%s: %w", namespace, secretName, err)
	}

	// Extract credentials from secret data
	username := string(secret.Data["username"])
	password := string(secret.Data["password"])

	if username == "" || password == "" {
		return nil, errors.New("secret missing username or password")
	}

	creds := &Credentials{
		Username: username,
		Password: password,
		Registry: registry,
	}

	// Validate secret credentials
	if err := f.validateCredentials(creds); err != nil {
		return nil, fmt.Errorf("invalid secret credentials: %w", err)
	}

	return creds, nil
}

// validateCredentials checks if credentials are minimally valid
func (f *AuthFlow) validateCredentials(creds *Credentials) error {
	if creds == nil {
		return errors.New("credentials are nil")
	}

	if creds.Username == "" {
		return errors.New("username is required")
	}

	if creds.Password == "" && creds.Token == "" {
		return errors.New("password or token is required")
	}

	return nil
}

// GetSource returns the source of credentials for display purposes
func (fc *FlowCredentials) GetSource() string {
	return fc.Source
}