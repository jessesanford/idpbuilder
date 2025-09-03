package builder

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
)

// CompressionType represents the type of compression to use for tarballs.
type CompressionType int

const (
	// NoCompression indicates no compression
	NoCompression CompressionType = iota
	// GzipCompression indicates gzip compression
	GzipCompression
	// Bzip2Compression indicates bzip2 compression (future support)
	Bzip2Compression
	// XzCompression indicates xz compression (future support)
	XzCompression
)

// TarOption represents an option for tarball creation.
type TarOption func(*TarConfig)

// TarConfig holds configuration for tarball operations.
type TarConfig struct {
	Compression CompressionType
	OwnerUID    int
	OwnerGID    int
	ModTime     time.Time
	Excludes    []string
}

// TarProcessor processes each entry in a tarball stream.
type TarProcessor func(*tar.Header, io.Reader) error

// WithCompression sets the compression type for tarball operations.
func WithCompression(compression CompressionType) TarOption {
	return func(config *TarConfig) {
		config.Compression = compression
	}
}

// WithOwnership sets the owner UID and GID for tarball entries.
func WithOwnership(uid, gid int) TarOption {
	return func(config *TarConfig) {
		config.OwnerUID = uid
		config.OwnerGID = gid
	}
}

// WithModTime sets the modification time for tarball entries.
func WithModTime(modTime time.Time) TarOption {
	return func(config *TarConfig) {
		config.ModTime = modTime
	}
}

// WithExcludes sets patterns to exclude from the tarball.
func WithExcludes(excludes ...string) TarOption {
	return func(config *TarConfig) {
		config.Excludes = excludes
	}
}

// CreateTarball creates a tarball from a source directory to destination file.
func CreateTarball(src string, dest string, options ...TarOption) error {
	config := &TarConfig{
		Compression: NoCompression,
		OwnerUID:    0,
		OwnerGID:    0,
		ModTime:     time.Now(),
	}

	for _, option := range options {
		option(config)
	}

	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	var writer io.Writer = destFile

	// Apply compression if specified
	switch config.Compression {
	case GzipCompression:
		gzipWriter := gzip.NewWriter(destFile)
		defer gzipWriter.Close()
		writer = gzipWriter
	case Bzip2Compression, XzCompression:
		return fmt.Errorf("compression type %v not yet supported", config.Compression)
	}

	tarWriter := tar.NewWriter(writer)
	defer tarWriter.Close()

	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check exclusions
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		for _, exclude := range config.Excludes {
			if matched, _ := filepath.Match(exclude, relPath); matched {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}

		header.Name = relPath
		header.Uid = config.OwnerUID
		header.Gid = config.OwnerGID
		header.ModTime = config.ModTime

		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(tarWriter, file)
			return err
		}

		return nil
	})
}

// CreateTarballFromFiles creates a tarball from a list of files.
func CreateTarballFromFiles(files []string, dest string, options ...TarOption) error {
	config := &TarConfig{
		Compression: NoCompression,
		OwnerUID:    0,
		OwnerGID:    0,
		ModTime:     time.Now(),
	}

	for _, option := range options {
		option(config)
	}

	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	var writer io.Writer = destFile

	// Apply compression if specified
	switch config.Compression {
	case GzipCompression:
		gzipWriter := gzip.NewWriter(destFile)
		defer gzipWriter.Close()
		writer = gzipWriter
	}

	tarWriter := tar.NewWriter(writer)
	defer tarWriter.Close()

	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			return fmt.Errorf("failed to stat file %s: %w", file, err)
		}

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}

		header.Name = filepath.Base(file)
		header.Uid = config.OwnerUID
		header.Gid = config.OwnerGID
		header.ModTime = config.ModTime

		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		if !info.IsDir() {
			f, err := os.Open(file)
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(tarWriter, f)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ExtractTarball extracts a tarball to a destination directory.
func ExtractTarball(src string, dest string) error {
	file, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer file.Close()

	var reader io.Reader = file

	// Check if the file is gzipped
	if strings.HasSuffix(src, ".gz") || strings.HasSuffix(src, ".tgz") {
		gzipReader, err := gzip.NewReader(file)
		if err != nil {
			return fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer gzipReader.Close()
		reader = gzipReader
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

		target := filepath.Join(dest, header.Name)

		// Ensure target directory exists
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", target, err)
			}

		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("failed to create file %s: %w", target, err)
			}

			if _, err := io.Copy(f, tarReader); err != nil {
				f.Close()
				return fmt.Errorf("failed to extract file %s: %w", target, err)
			}
			f.Close()

		default:
			// Handle other types like symlinks, etc. if needed
			continue
		}
	}

	return nil
}

// StreamTarball streams a tarball and processes each entry.
func StreamTarball(reader io.Reader, processor TarProcessor) error {
	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar header: %w", err)
		}

		if err := processor(header, tarReader); err != nil {
			return fmt.Errorf("processor failed for entry %s: %w", header.Name, err)
		}
	}

	return nil
}

// CompressTarball compresses an existing tarball with the specified compression.
func CompressTarball(input string, output string, compression CompressionType) error {
	inputFile, err := os.Open(input)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	var writer io.Writer = outputFile

	switch compression {
	case GzipCompression:
		gzipWriter := gzip.NewWriter(outputFile)
		defer gzipWriter.Close()
		writer = gzipWriter
	case NoCompression:
		// No compression, just copy
	default:
		return fmt.Errorf("compression type %v not supported", compression)
	}

	_, err = io.Copy(writer, inputFile)
	return err
}

// CreateLayerFromTarball creates a v1.Layer from a tarball file.
func CreateLayerFromTarball(tarballPath string) (v1.Layer, error) {
	return tarball.LayerFromFile(tarballPath)
}

// CreateLayerFromTarballReader creates a v1.Layer from a tarball reader.
func CreateLayerFromTarballReader(reader io.ReadCloser) (v1.Layer, error) {
	return tarball.LayerFromReader(reader)
}