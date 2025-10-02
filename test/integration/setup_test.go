package integration

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// TestEnvironment represents the integration test environment with idpbuilder
type TestEnvironment struct {
	ClusterName      string
	GiteaRegistry    string
	GiteaUsername    string
	GiteaPassword    string
	CertPath         string
	KeyPath          string
	CleanupFuncs     []func() error
	idpbuilderCmd    *exec.Cmd
	registryReady    bool
	clusterStartTime time.Time
}

// NewTestEnvironment creates a new test environment
func NewTestEnvironment(t *testing.T) *TestEnvironment {
	return &TestEnvironment{
		ClusterName:   fmt.Sprintf("idp-test-%d", time.Now().Unix()),
		GiteaUsername: "testuser",
		GiteaPassword: "testpass123",
		CleanupFuncs:  make([]func() error, 0),
	}
}

// SetupIDPBuilder initializes an idpbuilder cluster for testing
func (env *TestEnvironment) SetupIDPBuilder(t *testing.T) error {
	t.Log("Starting idpbuilder cluster setup...")
	env.clusterStartTime = time.Now()

	// Create temporary directory for cluster config
	tmpDir, err := os.MkdirTemp("", "idpbuilder-test-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Add cleanup for temp directory
	env.AddCleanup(func() error {
		return os.RemoveAll(tmpDir)
	})

	// Generate self-signed certificates for registry
	if err := env.generateSelfSignedCerts(tmpDir); err != nil {
		return fmt.Errorf("failed to generate certificates: %w", err)
	}

	// Create idpbuilder config
	configPath := filepath.Join(tmpDir, "idpbuilder-config.yaml")
	if err := env.createIDPBuilderConfig(configPath); err != nil {
		return fmt.Errorf("failed to create idpbuilder config: %w", err)
	}

	// Start idpbuilder cluster
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	env.idpbuilderCmd = exec.CommandContext(ctx, "idpbuilder", "create",
		"--name", env.ClusterName,
		"--config", configPath,
		"--wait",
	)

	env.idpbuilderCmd.Stdout = os.Stdout
	env.idpbuilderCmd.Stderr = os.Stderr

	t.Logf("Executing: %s", env.idpbuilderCmd.String())
	if err := env.idpbuilderCmd.Start(); err != nil {
		return fmt.Errorf("failed to start idpbuilder: %w", err)
	}

	// Add cleanup to destroy cluster
	env.AddCleanup(func() error {
		return env.destroyCluster(t)
	})

	// Wait for cluster to be ready
	if err := env.idpbuilderCmd.Wait(); err != nil {
		return fmt.Errorf("idpbuilder cluster creation failed: %w", err)
	}

	t.Logf("idpbuilder cluster created in %v", time.Since(env.clusterStartTime))

	// Get Gitea registry URL
	if err := env.discoverGiteaRegistry(t); err != nil {
		return fmt.Errorf("failed to discover Gitea registry: %w", err)
	}

	// Setup test user in Gitea
	if err := env.setupGiteaUser(t); err != nil {
		return fmt.Errorf("failed to setup Gitea user: %w", err)
	}

	env.registryReady = true
	t.Log("Test environment setup completed successfully")
	return nil
}

// generateSelfSignedCerts creates self-signed certificates for the registry
func (env *TestEnvironment) generateSelfSignedCerts(dir string) error {
	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %w", err)
	}

	// Create certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"IDPBuilder Test"},
			CommonName:   "localhost",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost", "gitea", "gitea.gitea.svc.cluster.local"},
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
	}

	// Create certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %w", err)
	}

	// Save certificate
	env.CertPath = filepath.Join(dir, "cert.pem")
	certFile, err := os.Create(env.CertPath)
	if err != nil {
		return fmt.Errorf("failed to create cert file: %w", err)
	}
	defer certFile.Close()

	if err := pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: certDER}); err != nil {
		return fmt.Errorf("failed to encode certificate: %w", err)
	}

	// Save private key
	env.KeyPath = filepath.Join(dir, "key.pem")
	keyFile, err := os.Create(env.KeyPath)
	if err != nil {
		return fmt.Errorf("failed to create key file: %w", err)
	}
	defer keyFile.Close()

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	if err := pem.Encode(keyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privateKeyBytes}); err != nil {
		return fmt.Errorf("failed to encode private key: %w", err)
	}

	return nil
}

// createIDPBuilderConfig creates the idpbuilder configuration file
func (env *TestEnvironment) createIDPBuilderConfig(path string) error {
	config := fmt.Sprintf(`apiVersion: idpbuilder.io/v1alpha1
kind: Config
metadata:
  name: integration-test-config
spec:
  gitea:
    enabled: true
    registry:
      enabled: true
      tls:
        enabled: true
        certPath: %s
        keyPath: %s
  argocd:
    enabled: false
  tekton:
    enabled: false
`, env.CertPath, env.KeyPath)

	return os.WriteFile(path, []byte(config), 0644)
}

// discoverGiteaRegistry discovers the Gitea registry URL from the cluster
func (env *TestEnvironment) discoverGiteaRegistry(t *testing.T) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "kubectl", "get", "service", "-n", "gitea", "gitea", "-o", "jsonpath={.spec.clusterIP}")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get Gitea service IP: %w", err)
	}

	env.GiteaRegistry = fmt.Sprintf("%s:5000", string(output))
	t.Logf("Discovered Gitea registry: %s", env.GiteaRegistry)
	return nil
}

// setupGiteaUser creates a test user in Gitea
func (env *TestEnvironment) setupGiteaUser(t *testing.T) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Execute command in Gitea pod to create user
	cmd := exec.CommandContext(ctx, "kubectl", "exec", "-n", "gitea", "deployment/gitea", "--",
		"gitea", "admin", "user", "create",
		"--username", env.GiteaUsername,
		"--password", env.GiteaPassword,
		"--email", "test@example.com",
		"--must-change-password=false",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("User creation output: %s", string(output))
		// User might already exist, which is okay
		if !contains(string(output), "already exists") {
			return fmt.Errorf("failed to create Gitea user: %w", err)
		}
	}

	t.Logf("Gitea user created: %s", env.GiteaUsername)
	return nil
}

// destroyCluster destroys the idpbuilder cluster
func (env *TestEnvironment) destroyCluster(t *testing.T) error {
	if env.ClusterName == "" {
		return nil
	}

	t.Logf("Destroying idpbuilder cluster: %s", env.ClusterName)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "idpbuilder", "delete", "--name", env.ClusterName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		t.Logf("Warning: failed to destroy cluster %s: %v", env.ClusterName, err)
		return err
	}

	t.Logf("Cluster %s destroyed successfully", env.ClusterName)
	return nil
}

// AddCleanup adds a cleanup function to be executed during teardown
func (env *TestEnvironment) AddCleanup(fn func() error) {
	env.CleanupFuncs = append(env.CleanupFuncs, fn)
}

// Cleanup executes all cleanup functions in reverse order
func (env *TestEnvironment) Cleanup(t *testing.T) {
	t.Log("Running test environment cleanup...")
	for i := len(env.CleanupFuncs) - 1; i >= 0; i-- {
		if err := env.CleanupFuncs[i](); err != nil {
			t.Errorf("Cleanup function failed: %v", err)
		}
	}
	t.Log("Test environment cleanup completed")
}

// IsReady returns whether the test environment is ready for use
func (env *TestEnvironment) IsReady() bool {
	return env.registryReady
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
