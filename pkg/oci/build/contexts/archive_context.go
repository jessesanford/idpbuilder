package contexts

import (
	"archive/tar"
	"archive/zip"
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ArchiveFormat represents supported archive formats
type ArchiveFormat int

const (
	FormatUnknown ArchiveFormat = iota
	FormatTar
	FormatTarGz
	FormatTarBz2
	FormatZip
)

// ArchiveContextImpl implements Context for archive file sources
type ArchiveContextImpl struct {
	archivePath string
	config      *ContextConfig
	extractDir  string
	format      ArchiveFormat
}

// Path returns the path where the archive was extracted
func (a *ArchiveContextImpl) Path() string {
	return a.extractDir
}

// Type returns ArchiveContext
func (a *ArchiveContextImpl) Type() ContextType {
	return ArchiveContext
}

// Cleanup removes the temporary extraction directory
func (a *ArchiveContextImpl) Cleanup() error {
	if a.extractDir != "" {
		return os.RemoveAll(a.extractDir)
	}
	return nil
}

// Extract extracts the archive and returns the extraction path
func (a *ArchiveContextImpl) Extract() (string, error) {
	// Detect archive format
	format, err := a.detectFormat()
	if err != nil {
		return "", fmt.Errorf("failed to detect archive format: %w", err)
	}
	a.format = format

	// Create temporary extraction directory
	tempDir, err := os.MkdirTemp(a.config.TempDir, "archive_context_*")
	if err != nil {
		return "", fmt.Errorf("failed to create extraction directory: %w", err)
	}
	a.extractDir = tempDir

	// Extract based on format
	switch format {
	case FormatTar:
		return tempDir, a.extractTar(nil)
	case FormatTarGz:
		return tempDir, a.extractTarGz()
	case FormatTarBz2:
		return tempDir, a.extractTarBz2()
	case FormatZip:
		return tempDir, a.extractZip()
	default:
		return "", fmt.Errorf("unsupported archive format")
	}
}

// detectFormat determines the archive format from file extension and magic bytes
func (a *ArchiveContextImpl) detectFormat() (ArchiveFormat, error) {
	// Check file extension first
	lower := strings.ToLower(a.archivePath)
	if strings.HasSuffix(lower, ".tar.gz") || strings.HasSuffix(lower, ".tgz") {
		return FormatTarGz, nil
	}
	if strings.HasSuffix(lower, ".tar.bz2") || strings.HasSuffix(lower, ".tbz2") {
		return FormatTarBz2, nil
	}
	if strings.HasSuffix(lower, ".tar") {
		return FormatTar, nil
	}
	if strings.HasSuffix(lower, ".zip") {
		return FormatZip, nil
	}

	// Check magic bytes
	file, err := os.Open(a.archivePath)
	if err != nil {
		return FormatUnknown, err
	}
	defer file.Close()

	magic := make([]byte, 4)
	if _, err := file.Read(magic); err != nil {
		return FormatUnknown, err
	}

	// ZIP magic
	if magic[0] == 0x50 && magic[1] == 0x4B {
		return FormatZip, nil
	}

	// GZIP magic
	if magic[0] == 0x1F && magic[1] == 0x8B {
		return FormatTarGz, nil
	}

	// BZ2 magic
	if magic[0] == 0x42 && magic[1] == 0x5A && magic[2] == 0x68 {
		return FormatTarBz2, nil
	}

	// Check for tar signature (ustar)
	file.Seek(257, io.SeekStart)
	ustarMagic := make([]byte, 5)
	if n, _ := file.Read(ustarMagic); n == 5 && string(ustarMagic) == "ustar" {
		return FormatTar, nil
	}

	return FormatUnknown, fmt.Errorf("unknown archive format")
}

// extractTar extracts a plain tar archive
func (a *ArchiveContextImpl) extractTar(reader io.Reader) error {
	var tarReader *tar.Reader
	
	if reader != nil {
		tarReader = tar.NewReader(reader)
	} else {
		file, err := os.Open(a.archivePath)
		if err != nil {
			return fmt.Errorf("failed to open archive: %w", err)
		}
		defer file.Close()
		tarReader = tar.NewReader(file)
	}

	return a.extractTarReader(tarReader)
}

// extractTarGz extracts a gzipped tar archive
func (a *ArchiveContextImpl) extractTarGz() error {
	file, err := os.Open(a.archivePath)
	if err != nil {
		return fmt.Errorf("failed to open archive: %w", err)
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzReader.Close()

	return a.extractTar(gzReader)
}

// extractTarBz2 extracts a bzip2 compressed tar archive
func (a *ArchiveContextImpl) extractTarBz2() error {
	file, err := os.Open(a.archivePath)
	if err != nil {
		return fmt.Errorf("failed to open archive: %w", err)
	}
	defer file.Close()

	bz2Reader := bzip2.NewReader(file)
	return a.extractTar(bz2Reader)
}

// extractTarReader performs the actual tar extraction
func (a *ArchiveContextImpl) extractTarReader(tarReader *tar.Reader) error {
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar header: %w", err)
		}

		// Validate and clean the path
		cleanPath, err := a.validatePath(header.Name)
		if err != nil {
			return fmt.Errorf("invalid path in archive: %w", err)
		}

		targetPath := filepath.Join(a.extractDir, cleanPath)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", targetPath, err)
			}

		case tar.TypeReg:
			// Create parent directories
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return fmt.Errorf("failed to create parent directory: %w", err)
			}

			// Extract file
			if err := a.extractTarFile(tarReader, targetPath, header.Mode); err != nil {
				return fmt.Errorf("failed to extract file %s: %w", cleanPath, err)
			}

		case tar.TypeSymlink, tar.TypeLink:
			// Skip symbolic and hard links for security
			continue
		}
	}

	return nil
}

// extractTarFile extracts a single file from tar
func (a *ArchiveContextImpl) extractTarFile(reader io.Reader, targetPath string, mode int64) error {
	file, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(mode))
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Copy with size limit
	_, err = io.CopyN(file, reader, a.config.MaxSize)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to write file content: %w", err)
	}

	return nil
}

// extractZip extracts a ZIP archive
func (a *ArchiveContextImpl) extractZip() error {
	reader, err := zip.OpenReader(a.archivePath)
	if err != nil {
		return fmt.Errorf("failed to open zip archive: %w", err)
	}
	defer reader.Close()

	for _, file := range reader.File {
		// Validate and clean the path
		cleanPath, err := a.validatePath(file.Name)
		if err != nil {
			return fmt.Errorf("invalid path in zip: %w", err)
		}

		targetPath := filepath.Join(a.extractDir, cleanPath)

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(targetPath, file.FileInfo().Mode()); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", targetPath, err)
			}
			continue
		}

		// Extract regular file
		if err := a.extractZipFile(file, targetPath); err != nil {
			return fmt.Errorf("failed to extract file %s: %w", cleanPath, err)
		}
	}

	return nil
}

// extractZipFile extracts a single file from ZIP
func (a *ArchiveContextImpl) extractZipFile(file *zip.File, targetPath string) error {
	// Create parent directories
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	reader, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open file in zip: %w", err)
	}
	defer reader.Close()

	targetFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, file.FileInfo().Mode())
	if err != nil {
		return fmt.Errorf("failed to create target file: %w", err)
	}
	defer targetFile.Close()

	// Copy with size limit
	_, err = io.CopyN(targetFile, reader, a.config.MaxSize)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to write file content: %w", err)
	}

	return nil
}

// validatePath validates and cleans an archive path to prevent directory traversal
func (a *ArchiveContextImpl) validatePath(archivePath string) (string, error) {
	// Clean the path and check for traversal attempts
	cleanPath := filepath.Clean(archivePath)
	
	// Reject absolute paths
	if filepath.IsAbs(cleanPath) {
		return "", fmt.Errorf("absolute paths not allowed: %s", archivePath)
	}
	
	// Reject paths that would escape the extraction directory
	if strings.Contains(cleanPath, "..") {
		return "", fmt.Errorf("path traversal attempt: %s", archivePath)
	}
	
	// Reject paths with null bytes or other problematic characters
	if strings.ContainsAny(cleanPath, "\x00") {
		return "", fmt.Errorf("invalid characters in path: %s", archivePath)
	}
	
	return cleanPath, nil
}