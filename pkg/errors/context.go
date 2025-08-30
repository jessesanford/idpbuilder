package errors

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// ErrorContext provides rich context information for certificate errors
type ErrorContext struct {
	Timestamp     time.Time         `json:"timestamp"`
	Component     string            `json:"component"`     // Which component failed
	Operation     string            `json:"operation"`     // What was being attempted
	Environment   map[string]string `json:"environment"`   // Relevant env vars
	SystemInfo    SystemInfo        `json:"system_info"`   // OS, arch, versions
	Diagnostics   map[string]string `json:"diagnostics"`   // Command outputs
	CorrelationID string            `json:"correlation_id"` // For tracking related errors
}

// SystemInfo contains information about the system environment
type SystemInfo struct {
	OS           string `json:"os"`
	Arch         string `json:"arch"`
	GoVersion    string `json:"go_version"`
	KindVersion  string `json:"kind_version,omitempty"`
	DockerInfo   string `json:"docker_info,omitempty"`
	KubeVersion  string `json:"kube_version,omitempty"`
}

// CaptureContext collects comprehensive error context
func CaptureContext(component, operation string) *ErrorContext {
	ctx := &ErrorContext{
		Timestamp:     time.Now(),
		Component:     component,
		Operation:     operation,
		Environment:   captureEnvironment(),
		SystemInfo:    captureSystemInfo(),
		Diagnostics:   captureDiagnostics(),
		CorrelationID: generateCorrelationID(),
	}
	return ctx
}

// captureEnvironment collects relevant environment variables
func captureEnvironment() map[string]string {
	env := make(map[string]string)
	
	// Certificate-related environment variables
	relevantEnvVars := []string{
		"KUBECONFIG",
		"DOCKER_HOST",
		"SSL_CERT_DIR",
		"SSL_CERT_FILE",
		"REQUESTS_CA_BUNDLE",
		"CURL_CA_BUNDLE",
		"HOME",
		"USER",
		"PWD",
		"PATH",
	}
	
	for _, varName := range relevantEnvVars {
		if value := os.Getenv(varName); value != "" {
			// Sanitize sensitive paths
			if strings.Contains(varName, "CERT") || varName == "KUBECONFIG" {
				env[varName] = sanitizePath(value)
			} else {
				env[varName] = value
			}
		}
	}
	
	return env
}

// captureSystemInfo gathers system information
func captureSystemInfo() SystemInfo {
	info := SystemInfo{
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		GoVersion: runtime.Version(),
	}
	
	// Try to get kind version
	if output, err := runCommand("kind", "version"); err == nil {
		info.KindVersion = strings.TrimSpace(output)
	}
	
	// Try to get docker info (just version)
	if output, err := runCommand("docker", "version", "--format", "{{.Server.Version}}"); err == nil {
		info.DockerInfo = strings.TrimSpace(output)
	}
	
	// Try to get kubernetes version
	if output, err := runCommand("kubectl", "version", "--client", "--short"); err == nil {
		info.KubeVersion = strings.TrimSpace(output)
	}
	
	return info
}

// captureDiagnostics runs diagnostic commands safely
func captureDiagnostics() map[string]string {
	diagnostics := make(map[string]string)
	
	// Define safe diagnostic commands
	commands := map[string][]string{
		"kubectl_config_view":     {"kubectl", "config", "view", "--minify"},
		"kubectl_cluster_info":    {"kubectl", "cluster-info"},
		"docker_ps":              {"docker", "ps", "--format", "table {{.Names}}\\t{{.Status}}"},
		"kind_get_clusters":      {"kind", "get", "clusters"},
		"openssl_version":        {"openssl", "version"},
		"ca_certificates_check":  {"ls", "/usr/local/share/ca-certificates"},
	}
	
	for name, cmd := range commands {
		if output, err := runCommandWithTimeout(cmd, 5*time.Second); err == nil {
			diagnostics[name] = strings.TrimSpace(output)
		} else {
			diagnostics[name] = fmt.Sprintf("Failed to run: %v", err)
		}
	}
	
	return diagnostics
}

// runCommand executes a command and returns its output
func runCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// runCommandWithTimeout executes a command with a timeout
func runCommandWithTimeout(cmdArgs []string, timeout time.Duration) (string, error) {
	if len(cmdArgs) == 0 {
		return "", fmt.Errorf("empty command")
	}
	
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	
	// Create a channel to receive the result
	done := make(chan struct {
		output string
		err    error
	}, 1)
	
	go func() {
		output, err := cmd.CombinedOutput()
		done <- struct {
			output string
			err    error
		}{string(output), err}
	}()
	
	// Wait for either completion or timeout
	select {
	case result := <-done:
		return result.output, result.err
	case <-time.After(timeout):
		cmd.Process.Kill()
		return "", fmt.Errorf("command timeout after %v", timeout)
	}
}

// sanitizePath removes sensitive information from paths
func sanitizePath(path string) string {
	// Replace home directory with ~
	if home := os.Getenv("HOME"); home != "" && strings.HasPrefix(path, home) {
		return strings.Replace(path, home, "~", 1)
	}
	return path
}

// generateCorrelationID creates a unique ID for error correlation
func generateCorrelationID() string {
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), os.Getpid())
}

// ToJSON serializes the error context to JSON
func (ctx *ErrorContext) ToJSON() (string, error) {
	data, err := json.MarshalIndent(ctx, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal context to JSON: %w", err)
	}
	return string(data), nil
}

// ToString provides a human-readable representation of the context
func (ctx *ErrorContext) ToString() string {
	var sb strings.Builder
	
	sb.WriteString(fmt.Sprintf("Context ID: %s\n", ctx.CorrelationID))
	sb.WriteString(fmt.Sprintf("Timestamp: %s\n", ctx.Timestamp.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("Component: %s\n", ctx.Component))
	sb.WriteString(fmt.Sprintf("Operation: %s\n", ctx.Operation))
	
	sb.WriteString("\nSystem Information:\n")
	sb.WriteString(fmt.Sprintf("  OS: %s\n", ctx.SystemInfo.OS))
	sb.WriteString(fmt.Sprintf("  Architecture: %s\n", ctx.SystemInfo.Arch))
	sb.WriteString(fmt.Sprintf("  Go Version: %s\n", ctx.SystemInfo.GoVersion))
	if ctx.SystemInfo.KindVersion != "" {
		sb.WriteString(fmt.Sprintf("  Kind Version: %s\n", ctx.SystemInfo.KindVersion))
	}
	if ctx.SystemInfo.DockerInfo != "" {
		sb.WriteString(fmt.Sprintf("  Docker Version: %s\n", ctx.SystemInfo.DockerInfo))
	}
	if ctx.SystemInfo.KubeVersion != "" {
		sb.WriteString(fmt.Sprintf("  Kubectl Version: %s\n", ctx.SystemInfo.KubeVersion))
	}
	
	if len(ctx.Environment) > 0 {
		sb.WriteString("\nRelevant Environment:\n")
		for key, value := range ctx.Environment {
			sb.WriteString(fmt.Sprintf("  %s: %s\n", key, value))
		}
	}
	
	if len(ctx.Diagnostics) > 0 {
		sb.WriteString("\nDiagnostics:\n")
		for key, value := range ctx.Diagnostics {
			sb.WriteString(fmt.Sprintf("  %s: %s\n", key, value))
		}
	}
	
	return sb.String()
}

// EnrichError adds context to a CertificateError
func EnrichError(err *CertificateError, component, operation string) *CertificateError {
	ctx := CaptureContext(component, operation)
	
	// Add context information to error details
	if err.Details == nil {
		err.Details = make(map[string]string)
	}
	
	err.Details["correlation_id"] = ctx.CorrelationID
	err.Details["context_timestamp"] = ctx.Timestamp.Format(time.RFC3339)
	err.Component = ctx.Component
	err.Operation = ctx.Operation
	
	return err
}