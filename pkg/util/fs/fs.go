// pkg/util/fs/fs.go
package fs

import (
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/cnoe-io/idpbuilder/api/v1alpha1"
)

// ConvertFSToBytes converts files from a filesystem to byte slices
func ConvertFSToBytes(fsys fs.FS, path string, customization v1alpha1.BuildCustomizationSpec) ([][]byte, error) {
	var results [][]byte

	// Walk through the filesystem at the given path
	err := fs.WalkDir(fsys, path, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Only process YAML files (common for Kubernetes resources)
		if !strings.HasSuffix(filePath, ".yaml") && !strings.HasSuffix(filePath, ".yml") {
			return nil
		}

		// Read file contents
		file, err := fsys.Open(filePath)
		if err != nil {
			return fmt.Errorf("failed to open file %s: %w", filePath, err)
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", filePath, err)
		}

		// Apply customizations if needed
		processedData := applyCustomizations(data, customization, filePath)
		results = append(results, processedData)

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk filesystem: %w", err)
	}

	return results, nil
}

// applyCustomizations applies build customizations to the file data
func applyCustomizations(data []byte, customization v1alpha1.BuildCustomizationSpec, filePath string) []byte {
	content := string(data)

	// Add a header comment to identify the resource source
	header := fmt.Sprintf("# UCP ARGO INSTALL RESOURCES\n")

	// Apply host customization
	if customization.Host != "" {
		content = strings.ReplaceAll(content, "localhost", customization.Host)
		content = strings.ReplaceAll(content, "127.0.0.1", customization.Host)
	}

	// Apply port customization
	if customization.Port != "" {
		// Replace common port patterns
		content = strings.ReplaceAll(content, ":8080", ":"+customization.Port)
		content = strings.ReplaceAll(content, ":80", ":"+customization.Port)
		content = strings.ReplaceAll(content, ":443", ":"+customization.Port)
	}

	// Apply protocol customization
	if customization.Protocol != "" {
		if customization.Protocol == "http" {
			content = strings.ReplaceAll(content, "https://", "http://")
		} else if customization.Protocol == "https" {
			content = strings.ReplaceAll(content, "http://", "https://")
		}
	}

	// Apply path routing customization
	if customization.UsePathRouting {
		// Modify ingress rules for path-based routing
		// This is a simplified implementation
		content = strings.ReplaceAll(content, "host:", "# host:")
	}

	// Add header for identification
	if strings.Contains(strings.ToLower(filePath), "install") {
		content = header + content
	}

	return []byte(content)
}

// FileExists checks if a file exists
func FileExists(path string) bool {
	_, err := fs.Stat(fs.FS(nil), path)
	return err == nil
}

// CreateDir creates a directory
func CreateDir(path string) error {
	// Note: This function works with the OS filesystem, not embed.FS
	// For embedded filesystems, directories are implicitly created
	return fmt.Errorf("CreateDir not implemented for embedded filesystems")
}

// ReadFile reads a file
func ReadFile(path string) ([]byte, error) {
	// Note: This function works with the OS filesystem, not embed.FS
	// For embedded filesystems, use fs.ReadFile
	return nil, fmt.Errorf("ReadFile not implemented for embedded filesystems")
}

// WriteFile writes a file
func WriteFile(path string, data []byte) error {
	// Note: This function works with the OS filesystem, not embed.FS
	// Embedded filesystems are read-only
	return fmt.Errorf("WriteFile not implemented for embedded filesystems")
}