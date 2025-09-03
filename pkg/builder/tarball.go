package builder

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
)

// TarballWriter handles creation of OCI image layers as tarballs.
type TarballWriter struct {
	timestampPolicy TimestampPolicy
	defaultMode     os.FileMode
	preserveOwner   bool
}

// TarballOptions configures tarball creation behavior.
type TarballOptions struct {
	TimestampPolicy TimestampPolicy `json:"timestamp_policy,omitempty"`
	PreserveOwner   bool            `json:"preserve_owner,omitempty"`
	DefaultMode     os.FileMode     `json:"default_mode,omitempty"`
	Compression     bool            `json:"compression,omitempty"`
}

// TarballEntry represents an entry in a tarball layer.
type TarballEntry struct {
	Header tar.Header
	Data   io.Reader
}

// NewTarballWriter creates a new TarballWriter with the given options.
func NewTarballWriter(opts TarballOptions) *TarballWriter {
	if opts.DefaultMode == 0 {
		opts.DefaultMode = 0644
	}

	return &TarballWriter{
		timestampPolicy: opts.TimestampPolicy,
		defaultMode:     opts.DefaultMode,
		preserveOwner:   opts.PreserveOwner,
	}
}

// CreateLayerFromFiles creates a layer from a list of file entries.
func (tw *TarballWriter) CreateLayerFromFiles(ctx context.Context, entries []FileEntry) (v1.Layer, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context cannot be nil")
	}

	if len(entries) == 0 {
		return nil, fmt.Errorf("no files specified for layer creation")
	}

	// Prepare and validate entries
	prepared, err := tw.prepareEntries(entries)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare entries: %w", err)
	}

	// Create tarball buffer
	var buf bytes.Buffer
	tarWriter := tar.NewWriter(&buf)
	defer tarWriter.Close()

	// Write all entries to tarball
	for _, entry := range prepared {
		if err := tw.writeEntryToTar(ctx, tarWriter, entry); err != nil {
			return nil, fmt.Errorf("failed to write entry %s: %w", entry.Path, err)
		}
	}

	// Close the tar writer to finalize
	if err := tarWriter.Close(); err != nil {
		return nil, fmt.Errorf("failed to close tar writer: %w", err)
	}

	// Create layer from tarball
	layer, err := tarball.LayerFromReader(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return nil, fmt.Errorf("failed to create layer from tarball: %w", err)
	}

	return layer, nil
}

// CreateLayerFromDirectory creates a layer from all files in a directory.
func (tw *TarballWriter) CreateLayerFromDirectory(ctx context.Context, sourceDir, targetDir string) (v1.Layer, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context cannot be nil")
	}

	if sourceDir == "" {
		return nil, fmt.Errorf("source directory cannot be empty")
	}

	if targetDir == "" {
		targetDir = "/"
	}

	// Collect all files from the directory
	entries, err := tw.collectDirectoryFiles(sourceDir, targetDir)
	if err != nil {
		return nil, fmt.Errorf("failed to collect directory files: %w", err)
	}

	return tw.CreateLayerFromFiles(ctx, entries)
}

// CreateLayerFromTarball creates a layer from an existing tarball reader.
func (tw *TarballWriter) CreateLayerFromTarball(ctx context.Context, reader io.Reader) (v1.Layer, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context cannot be nil")
	}

	if reader == nil {
		return nil, fmt.Errorf("reader cannot be nil")
	}

	// Create layer directly from the tarball reader
	layer, err := tarball.LayerFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to create layer from tarball reader: %w", err)
	}

	return layer, nil
}


// ExtractTarballToDirectory extracts a tarball to a directory.
func (tw *TarballWriter) ExtractTarballToDirectory(ctx context.Context, reader io.Reader, targetDir string) error {
	if ctx == nil {
		return fmt.Errorf("context cannot be nil")
	}
	if reader == nil {
		return fmt.Errorf("reader cannot be nil")
	}
	if targetDir == "" {
		return fmt.Errorf("target directory cannot be empty")
	}

	tarReader := tar.NewReader(reader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar header: %w", err)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		target := filepath.Join(targetDir, header.Name)
		if !strings.HasPrefix(target, filepath.Clean(targetDir)+string(os.PathSeparator)) &&
			target != filepath.Clean(targetDir) {
			return fmt.Errorf("path %s is outside target directory %s", header.Name, targetDir)
		}

		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return fmt.Errorf("failed to create parent directories for %s: %w", target, err)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, header.FileInfo().Mode()); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", target, err)
			}
		case tar.TypeReg:
			file, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, header.FileInfo().Mode())
			if err != nil {
				return fmt.Errorf("failed to create file %s: %w", target, err)
			}
			_, err = io.Copy(file, tarReader)
			file.Close()
			if err != nil {
				return fmt.Errorf("failed to write file %s: %w", target, err)
			}
		}

		if tw.timestampPolicy != TimestampPreserve {
			modTime := tw.applyTimestampPolicy(header.ModTime)
			os.Chtimes(target, modTime, modTime)
		}
	}
	return nil
}

// prepareEntries validates and prepares file entries for tarball creation.
func (tw *TarballWriter) prepareEntries(entries []FileEntry) ([]FileEntry, error) {
	if len(entries) == 0 {
		return nil, fmt.Errorf("no entries to prepare")
	}

	var prepared []FileEntry
	seenPaths := make(map[string]bool)

	for i, entry := range entries {
		// Validate required fields
		if entry.Path == "" {
			return nil, fmt.Errorf("entry %d: path cannot be empty", i)
		}

		// Normalize path
		entry.Path = filepath.Clean("/" + strings.TrimPrefix(entry.Path, "/"))

		// Check for duplicates
		if seenPaths[entry.Path] {
			return nil, fmt.Errorf("duplicate path: %s", entry.Path)
		}
		seenPaths[entry.Path] = true

		// Set defaults
		if entry.Mode == 0 {
			if entry.IsDir {
				entry.Mode = 0755
			} else {
				entry.Mode = tw.defaultMode
			}
		}

		// Apply timestamp policy
		entry.ModTime = tw.applyTimestampPolicy(entry.ModTime)

		// Validate source file exists for regular files
		if !entry.IsDir && entry.Source != "" {
			info, err := os.Stat(entry.Source)
			if err != nil {
				return nil, fmt.Errorf("entry %d: source file %s not accessible: %w", i, entry.Source, err)
			}
			entry.Size = info.Size()
			if entry.ModTime.IsZero() {
				entry.ModTime = tw.applyTimestampPolicy(info.ModTime())
			}
		}

		prepared = append(prepared, entry)
	}

	// Sort by path for deterministic output
	sort.Slice(prepared, func(i, j int) bool {
		return prepared[i].Path < prepared[j].Path
	})

	return prepared, nil
}

// writeEntryToTar writes a single entry to the tar writer.
func (tw *TarballWriter) writeEntryToTar(ctx context.Context, tarWriter *tar.Writer, entry FileEntry) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Create tar header
	header := &tar.Header{
		Name:     strings.TrimPrefix(entry.Path, "/"),
		Mode:     int64(entry.Mode),
		ModTime:  entry.ModTime,
		Typeflag: tar.TypeReg,
	}

	if entry.IsDir {
		header.Typeflag = tar.TypeDir
		header.Name = strings.TrimSuffix(header.Name, "/") + "/"
		header.Size = 0
	} else {
		header.Size = entry.Size
	}

	// Write header
	if err := tarWriter.WriteHeader(header); err != nil {
		return fmt.Errorf("failed to write header for %s: %w", entry.Path, err)
	}

	// Write file data for regular files
	if !entry.IsDir && entry.Source != "" {
		file, err := os.Open(entry.Source)
		if err != nil {
			return fmt.Errorf("failed to open source file %s: %w", entry.Source, err)
		}
		defer file.Close()

		_, err = io.Copy(tarWriter, file)
		if err != nil {
			return fmt.Errorf("failed to copy file data for %s: %w", entry.Source, err)
		}
	}

	return nil
}

// collectDirectoryFiles recursively collects files from a directory.
func (tw *TarballWriter) collectDirectoryFiles(sourceDir, targetDir string) ([]FileEntry, error) {
	var entries []FileEntry

	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if path == sourceDir {
			return nil
		}

		// Calculate relative path
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		// Create target path
		targetPath := filepath.Join(targetDir, relPath)

		entry := FileEntry{
			Path:    targetPath,
			Source:  path,
			Mode:    info.Mode(),
			ModTime: info.ModTime(),
			Size:    info.Size(),
			IsDir:   info.IsDir(),
		}

		entries = append(entries, entry)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory %s: %w", sourceDir, err)
	}

	return entries, nil
}

// applyTimestampPolicy applies the configured timestamp policy.
func (tw *TarballWriter) applyTimestampPolicy(original time.Time) time.Time {
	switch tw.timestampPolicy {
	case TimestampEpoch:
		return time.Unix(0, 0).UTC()
	case TimestampCurrent:
		return time.Now().UTC()
	case TimestampPreserve:
		fallthrough
	default:
		return original.UTC()
	}
}

// TarballInfo contains information about a tarball layer.
type TarballInfo struct {
	Digest v1.Hash `json:"digest"`
	Size   int64   `json:"size"`
}