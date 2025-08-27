package contexts

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ArchiveContext represents a context from an archive file
type ArchiveContext struct {
	archivePath  string
	extractedDir string
	entries      []ContextEntry
	filter       *Filter
}

// NewArchiveContext creates a new archive context
func NewArchiveContext(archivePath string) (*ArchiveContext, error) {
	lower := strings.ToLower(archivePath)
	if !strings.HasSuffix(lower, ".tar") && !strings.HasSuffix(lower, ".tar.gz") && !strings.HasSuffix(lower, ".tgz") {
		return nil, fmt.Errorf("unsupported archive format: %s", archivePath)
	}
	return &ArchiveContext{archivePath: archivePath}, nil
}

// PrepareContext extracts the archive and prepares it
func (ac *ArchiveContext) PrepareContext() error {
	tempDir, err := os.MkdirTemp("", "buildah-archive-*")
	if err != nil {
		return err
	}
	ac.extractedDir = tempDir
	if err := ac.extractArchive(); err != nil {
		os.RemoveAll(tempDir)
		return err
	}
	return ac.scanFiles()
}

// extractArchive extracts the tar archive
func (ac *ArchiveContext) extractArchive() error {
	file, err := os.Open(ac.archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var reader io.Reader = file
	if strings.HasSuffix(ac.archivePath, ".gz") || strings.HasSuffix(ac.archivePath, ".tgz") {
		gzReader, err := gzip.NewReader(file)
		if err != nil {
			return err
		}
		defer gzReader.Close()
		reader = gzReader
	}

	tarReader := tar.NewReader(reader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(ac.extractedDir, header.Name)
		// Path security check
		if !strings.HasPrefix(target, ac.extractedDir) {
			continue
		}

		switch header.Typeflag {
		case tar.TypeDir:
			os.MkdirAll(target, os.FileMode(header.Mode))
		case tar.TypeReg:
			os.MkdirAll(filepath.Dir(target), 0755)
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				continue
			}
			io.Copy(f, tarReader)
			f.Close()
		}
	}
	return nil
}

// scanFiles builds the list of context entries
func (ac *ArchiveContext) scanFiles() error {
	ac.entries = []ContextEntry{}
	return filepath.Walk(ac.extractedDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || path == ac.extractedDir {
			return nil
		}
		relPath, _ := filepath.Rel(ac.extractedDir, path)
		ac.entries = append(ac.entries, ContextEntry{
			Path:    relPath,
			Size:    info.Size(),
			ModTime: info.ModTime(),
			IsDir:   info.IsDir(),
		})
		return nil
	})
}

// GetEntries returns filtered entries
func (ac *ArchiveContext) GetEntries() ([]ContextEntry, error) {
	if ac.entries == nil {
		return nil, fmt.Errorf("context not prepared")
	}
	if ac.filter != nil {
		return ac.filter.ApplyFilter(ac.entries), nil
	}
	return ac.entries, nil
}

// GetContextDir returns the extracted directory
func (ac *ArchiveContext) GetContextDir() string {
	return ac.extractedDir
}

// SetFilter sets a dockerignore filter
func (ac *ArchiveContext) SetFilter(filter *Filter) {
	ac.filter = filter
}

// Cleanup removes temporary directory
func (ac *ArchiveContext) Cleanup() error {
	if ac.extractedDir != "" {
		return os.RemoveAll(ac.extractedDir)
	}
	return nil
}