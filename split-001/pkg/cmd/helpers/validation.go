package helpers

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

var (
	// ValidationMsg provides help text for validation flags
	ValidationMsg = "Enable strict validation of inputs"
	// NameRegex defines valid name patterns
	NameRegex = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)
)

// ValidateConfig validates configuration inputs
func ValidateConfig(configPath string) error {
	if configPath == "" {
		return fmt.Errorf("config path cannot be empty")
	}

	// Check if path exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("config file does not exist: %s", configPath)
	}

	// Validate file extension
	ext := filepath.Ext(configPath)
	validExts := []string{".yaml", ".yml", ".json"}
	isValid := false
	for _, validExt := range validExts {
		if ext == validExt {
			isValid = true
			break
		}
	}
	if !isValid {
		return fmt.Errorf("invalid config file extension %s, must be one of: %v", ext, validExts)
	}

	return nil
}

// ValidateName validates resource names according to Kubernetes naming conventions
func ValidateName(name string) error {
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	if len(name) > 63 {
		return fmt.Errorf("name cannot be longer than 63 characters")
	}

	if !NameRegex.MatchString(name) {
		return fmt.Errorf("name must contain only lowercase alphanumeric characters and hyphens, and must start and end with an alphanumeric character")
	}

	return nil
}

// ValidateNamespace validates Kubernetes namespace names
func ValidateNamespace(namespace string) error {
	if namespace == "" {
		return nil // empty namespace is valid (uses default)
	}

	if err := ValidateName(namespace); err != nil {
		return fmt.Errorf("invalid namespace: %w", err)
	}

	// Additional namespace-specific validations
	reserved := []string{"kube-system", "kube-public", "kube-node-lease", "default"}
	for _, res := range reserved {
		if namespace == res {
			LogWarning("Using reserved namespace: %s", namespace)
			break
		}
	}

	return nil
}

// ValidateDirectory validates directory paths
func ValidateDirectory(dir string) error {
	if dir == "" {
		return fmt.Errorf("directory path cannot be empty")
	}

	// Check if directory exists
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", dir)
	}
	if err != nil {
		return fmt.Errorf("error accessing directory %s: %w", dir, err)
	}

	// Check if it's actually a directory
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", dir)
	}

	// Check if directory is readable
	if !isReadable(dir) {
		return fmt.Errorf("directory is not readable: %s", dir)
	}

	return nil
}

// ValidatePort validates port numbers
func ValidatePort(port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535, got: %d", port)
	}

	// Check for common reserved ports
	reservedPorts := map[int]string{
		22: "SSH", 80: "HTTP", 443: "HTTPS", 
		3306: "MySQL", 5432: "PostgreSQL",
	}
	
	if service, isReserved := reservedPorts[port]; isReserved {
		LogWarning("Using reserved port %d (%s)", port, service)
	}

	return nil
}

// ValidateEmail validates email addresses (basic validation)
func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format: %s", email)
	}

	return nil
}

// ValidateVersion validates semantic version strings
func ValidateVersion(version string) error {
	if version == "" {
		return fmt.Errorf("version cannot be empty")
	}

	// Basic semantic version validation
	versionRegex := regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)(-[\w\.-]+)?(\+[\w\.-]+)?$`)
	if !versionRegex.MatchString(version) {
		return fmt.Errorf("invalid semantic version format: %s (expected format: v1.2.3 or 1.2.3)", version)
	}

	return nil
}

// SanitizeName sanitizes input to create valid names
func SanitizeName(input string) string {
	// Convert to lowercase
	name := strings.ToLower(input)
	
	// Replace invalid characters with hyphens
	var result strings.Builder
	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			result.WriteRune(r)
		} else {
			result.WriteRune('-')
		}
	}
	
	// Remove leading/trailing hyphens and collapse multiple hyphens
	name = strings.Trim(result.String(), "-")
	hyphenRegex := regexp.MustCompile(`-+`)
	name = hyphenRegex.ReplaceAllString(name, "-")
	
	// Truncate if too long
	if len(name) > 63 {
		name = name[:63]
		name = strings.TrimSuffix(name, "-")
	}
	
	return name
}

// isReadable checks if a directory is readable
func isReadable(dir string) bool {
	file, err := os.Open(dir)
	if err != nil {
		return false
	}
	file.Close()
	return true
}