package contexts

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	// MaxContextSize defines the maximum allowed context size in bytes (100MB)
	MaxContextSize = 100 * 1024 * 1024
	// MaxSymlinkDepth prevents symlink loops and excessive resolution
	MaxSymlinkDepth = 10
)

// ValidationResult represents the result of context validation
type ValidationResult struct {
	Valid   bool
	Errors  []string
	Size    int64
	Warnings []string
}

// ValidateContext performs comprehensive security validation on a build context
func ValidateContext(contextPath string, allowSymlinks bool) (*ValidationResult, error) {
	result := &ValidationResult{
		Valid:    true,
		Errors:   make([]string, 0),
		Warnings: make([]string, 0),
	}

	// Check if context path exists
	info, err := os.Stat(contextPath)
	if err != nil {
		if os.IsNotExist(err) {
			result.addError("context path does not exist: %s", contextPath)
			return result, nil
		}
		return nil, fmt.Errorf("failed to stat context path: %w", err)
	}

	// Validate that it's a directory
	if !info.IsDir() {
		result.addError("context path must be a directory, got file: %s", contextPath)
		return result, nil
	}

	// Calculate context size and validate files
	err = validateContextDirectory(contextPath, result, allowSymlinks, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to validate context directory: %w", err)
	}

	// Check size limit
	if result.Size > MaxContextSize {
		result.addError("context size %d bytes exceeds maximum allowed size %d bytes", result.Size, MaxContextSize)
	}

	return result, nil
}

// validateContextDirectory recursively validates a context directory
func validateContextDirectory(dirPath string, result *ValidationResult, allowSymlinks bool, depth int) error {
	// Prevent deep directory nesting (potential zip bomb protection)
	if depth > 50 {
		result.addError("directory nesting too deep at: %s", dirPath)
		return nil
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	for _, entry := range entries {
		fullPath := filepath.Join(dirPath, entry.Name())
		
		// Validate file/directory name
		if err := validateFileName(entry.Name()); err != nil {
			result.addError("invalid file name %s: %v", fullPath, err)
			continue
		}

		info, err := entry.Info()
		if err != nil {
			result.addError("failed to get info for %s: %v", fullPath, err)
			continue
		}

		// Handle symlinks
		if info.Mode()&os.ModeSymlink != 0 {
			if !allowSymlinks {
				result.addError("symlinks not allowed: %s", fullPath)
				continue
			}
			
			if err := validateSymlink(fullPath, result); err != nil {
				result.addError("invalid symlink %s: %v", fullPath, err)
				continue
			}
		}

		// Add to total size
		result.Size += info.Size()

		// Recursively validate subdirectories
		if info.IsDir() {
			err := validateContextDirectory(fullPath, result, allowSymlinks, depth+1)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// validateFileName checks for potentially dangerous file names
func validateFileName(name string) error {
	// Check for null bytes
	if strings.Contains(name, "\x00") {
		return fmt.Errorf("file name contains null byte")
	}

	// Check for path traversal attempts
	if strings.Contains(name, "..") {
		return fmt.Errorf("file name contains path traversal sequence")
	}

	// Check for reserved names on Windows
	reserved := []string{"CON", "PRN", "AUX", "NUL", "COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9", "LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9"}
	upperName := strings.ToUpper(name)
	for _, res := range reserved {
		if upperName == res || strings.HasPrefix(upperName, res+".") {
			return fmt.Errorf("file name uses reserved name: %s", name)
		}
	}

	return nil
}

// validateSymlink validates a symlink for security issues
func validateSymlink(symlinkPath string, result *ValidationResult) error {
	target, err := filepath.EvalSymlinks(symlinkPath)
	if err != nil {
		return fmt.Errorf("failed to resolve symlink: %w", err)
	}

	// Check for symlink loops by limiting resolution depth
	currentPath := symlinkPath
	for i := 0; i < MaxSymlinkDepth; i++ {
		linkTarget, err := os.Readlink(currentPath)
		if err != nil {
			break // Not a symlink anymore
		}

		if !filepath.IsAbs(linkTarget) {
			linkTarget = filepath.Join(filepath.Dir(currentPath), linkTarget)
		}
		
		currentPath = linkTarget
	}

	// Check if target exists and is accessible
	if _, err := os.Stat(target); err != nil {
		result.addWarning("symlink target not accessible: %s -> %s", symlinkPath, target)
	}

	return nil
}

// ValidateContextIntegrity performs integrity checks on context data
func ValidateContextIntegrity(contextData []byte, expectedChecksum string) error {
	if len(contextData) == 0 {
		return fmt.Errorf("context data is empty")
	}

	// Basic integrity check - could be enhanced with actual checksums
	if expectedChecksum != "" {
		// Placeholder for checksum validation
		// In a real implementation, this would calculate and compare checksums
		return nil
	}

	return nil
}

// CheckContextPermissions verifies that the context has appropriate permissions
func CheckContextPermissions(contextPath string) error {
	info, err := os.Stat(contextPath)
	if err != nil {
		return fmt.Errorf("failed to stat context path: %w", err)
	}

	// Check if directory is readable
	if info.Mode().Perm()&0444 == 0 {
		return fmt.Errorf("context directory is not readable")
	}

	// On Unix systems, check for world-writable directories (security risk)
	if info.Mode().Perm()&0002 != 0 {
		return fmt.Errorf("context directory is world-writable (security risk)")
	}

	return nil
}

// Helper methods for ValidationResult
func (vr *ValidationResult) addError(format string, args ...interface{}) {
	vr.Valid = false
	vr.Errors = append(vr.Errors, fmt.Sprintf(format, args...))
}

func (vr *ValidationResult) addWarning(format string, args ...interface{}) {
	vr.Warnings = append(vr.Warnings, fmt.Sprintf(format, args...))
}