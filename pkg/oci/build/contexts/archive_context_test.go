package contexts

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestArchiveFormat_String(t *testing.T) {
	tests := []struct {
		format   ArchiveFormat
		expected string
	}{
		{FormatUnknown, "FormatUnknown"},
		{FormatTar, "FormatTar"},
		{FormatTarGz, "FormatTarGz"},
		{FormatTarBz2, "FormatTarBz2"},
		{FormatZip, "FormatZip"},
	}

	for _, tt := range tests {
		// Note: ArchiveFormat doesn't have a String() method in the implementation
		// but we can test the enum values themselves
		if int(tt.format) < 0 || int(tt.format) > 4 {
			t.Errorf("Invalid ArchiveFormat value: %d", int(tt.format))
		}
	}
}

func TestArchiveContextImpl_Type(t *testing.T) {
	a := &ArchiveContextImpl{}
	if got := a.Type(); got != ArchiveContext {
		t.Errorf("Type() = %v, want %v", got, ArchiveContext)
	}
}

func TestArchiveContextImpl_Path(t *testing.T) {
	a := &ArchiveContextImpl{extractDir: "/test/extract/path"}
	if got := a.Path(); got != "/test/extract/path" {
		t.Errorf("Path() = %v, want %v", got, "/test/extract/path")
	}
}

func TestArchiveContextImpl_Cleanup(t *testing.T) {
	// Test cleanup with empty path
	a := &ArchiveContextImpl{}
	if err := a.Cleanup(); err != nil {
		t.Errorf("Cleanup() error = %v, want nil", err)
	}

	// Test cleanup with actual directory
	tempDir, err := os.MkdirTemp("", "archive_test_*")
	if err != nil {
		t.Fatal(err)
	}
	
	a.extractDir = tempDir
	if err := a.Cleanup(); err != nil {
		t.Errorf("Cleanup() error = %v, want nil", err)
	}
	
	// Verify directory was removed
	if _, err := os.Stat(tempDir); !os.IsNotExist(err) {
		t.Error("Cleanup() did not remove directory")
	}
}

func TestArchiveContextImpl_detectFormat(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "format_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name     string
		filename string
		content  []byte
		expected ArchiveFormat
	}{
		{
			name:     "tar.gz extension",
			filename: "test.tar.gz",
			content:  []byte("dummy"),
			expected: FormatTarGz,
		},
		{
			name:     "tgz extension",
			filename: "test.tgz", 
			content:  []byte("dummy"),
			expected: FormatTarGz,
		},
		{
			name:     "tar.bz2 extension",
			filename: "test.tar.bz2",
			content:  []byte("dummy"),
			expected: FormatTarBz2,
		},
		{
			name:     "tbz2 extension",
			filename: "test.tbz2",
			content:  []byte("dummy"),
			expected: FormatTarBz2,
		},
		{
			name:     "tar extension",
			filename: "test.tar",
			content:  []byte("dummy"),
			expected: FormatTar,
		},
		{
			name:     "zip extension",
			filename: "test.zip",
			content:  []byte("dummy"),
			expected: FormatZip,
		},
		{
			name:     "zip magic bytes",
			filename: "test.unknown",
			content:  []byte{0x50, 0x4B, 0x03, 0x04}, // ZIP magic
			expected: FormatZip,
		},
		{
			name:     "gzip magic bytes",
			filename: "test.unknown",
			content:  []byte{0x1F, 0x8B, 0x08, 0x00}, // GZIP magic
			expected: FormatTarGz,
		},
		{
			name:     "bzip2 magic bytes",
			filename: "test.unknown",
			content:  []byte{0x42, 0x5A, 0x68, 0x39}, // BZ2 magic
			expected: FormatTarBz2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tempDir, tt.filename)
			if err := os.WriteFile(testFile, tt.content, 0644); err != nil {
				t.Fatal(err)
			}
			defer os.Remove(testFile)

			a := &ArchiveContextImpl{archivePath: testFile}
			format, err := a.detectFormat()
			
			if err != nil && tt.expected != FormatUnknown {
				t.Errorf("detectFormat() error = %v, want nil", err)
				return
			}
			
			if format != tt.expected {
				t.Errorf("detectFormat() = %v, want %v", format, tt.expected)
			}
		})
	}
}

func TestArchiveContextImpl_detectFormat_UnknownFormat(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "unknown_format_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	testFile := filepath.Join(tempDir, "test.unknown")
	if err := os.WriteFile(testFile, []byte("invalid content"), 0644); err != nil {
		t.Fatal(err)
	}

	a := &ArchiveContextImpl{archivePath: testFile}
	format, err := a.detectFormat()
	
	if err == nil {
		t.Error("detectFormat() should return error for unknown format")
	}
	
	if format != FormatUnknown {
		t.Errorf("detectFormat() = %v, want %v", format, FormatUnknown)
	}
}

func TestArchiveContextImpl_validatePath(t *testing.T) {
	a := &ArchiveContextImpl{}

	tests := []struct {
		name        string
		archivePath string
		wantErr     bool
		errContains string
	}{
		{
			name:        "valid relative path",
			archivePath: "dir/file.txt",
			wantErr:     false,
		},
		{
			name:        "absolute path",
			archivePath: "/etc/passwd",
			wantErr:     true,
			errContains: "absolute paths not allowed",
		},
		{
			name:        "path traversal attempt",
			archivePath: "../../../etc/passwd",
			wantErr:     true,
			errContains: "path traversal attempt",
		},
		{
			name:        "path with null byte",
			archivePath: "file\x00.txt",
			wantErr:     true,
			errContains: "invalid characters in path",
		},
		{
			name:        "clean relative path with dots",
			archivePath: "./dir/file.txt",
			wantErr:     false,
		},
		{
			name:        "path with double dots in middle",
			archivePath: "dir/../file.txt",
			wantErr:     false, // filepath.Clean() resolves this to "file.txt"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanPath, err := a.validatePath(tt.archivePath)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("validatePath() expected error but got none")
					return
				}
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("validatePath() error = %v, want to contain %v", err, tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("validatePath() unexpected error = %v", err)
					return
				}
				if cleanPath == "" {
					t.Error("validatePath() should return clean path")
				}
			}
		})
	}
}

func TestArchiveContextImpl_Extract_InvalidFormat(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "extract_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create a file with unknown format
	testFile := filepath.Join(tempDir, "test.unknown")
	if err := os.WriteFile(testFile, []byte("invalid content"), 0644); err != nil {
		t.Fatal(err)
	}

	a := &ArchiveContextImpl{
		archivePath: testFile,
		config:      DefaultConfig(),
	}

	_, err = a.Extract()
	if err == nil {
		t.Error("Extract() should return error for unknown format")
	}
}

func TestArchiveContextImpl_extractTarFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "tar_file_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	a := &ArchiveContextImpl{
		config: &ContextConfig{MaxSize: 1024},
	}

	testContent := "Hello, World!"
	reader := strings.NewReader(testContent)
	targetPath := filepath.Join(tempDir, "test.txt")

	err = a.extractTarFile(reader, targetPath, 0644)
	if err != nil {
		t.Errorf("extractTarFile() error = %v, want nil", err)
		return
	}

	// Verify file was created with correct content
	content, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatal(err)
	}

	if string(content) != testContent {
		t.Errorf("extractTarFile() content = %v, want %v", string(content), testContent)
	}
}

// Helper function to create a simple tar archive for testing
func createTestTarArchive(t *testing.T, tempDir string) string {
	t.Helper()
	
	archivePath := filepath.Join(tempDir, "test.tar")
	file, err := os.Create(archivePath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	tw := tar.NewWriter(file)
	defer tw.Close()

	// Add a test file to the archive
	header := &tar.Header{
		Name: "testfile.txt",
		Mode: 0644,
		Size: int64(len("test content")),
		Typeflag: tar.TypeReg,
	}
	
	if err := tw.WriteHeader(header); err != nil {
		t.Fatal(err)
	}
	
	if _, err := tw.Write([]byte("test content")); err != nil {
		t.Fatal(err)
	}

	// Add a directory
	dirHeader := &tar.Header{
		Name: "testdir/",
		Mode: 0755,
		Typeflag: tar.TypeDir,
	}
	
	if err := tw.WriteHeader(dirHeader); err != nil {
		t.Fatal(err)
	}

	return archivePath
}

func TestArchiveContextImpl_extractTar_Success(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "tar_extract_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create test tar archive
	archivePath := createTestTarArchive(t, tempDir)

	a := &ArchiveContextImpl{
		archivePath: archivePath,
		config:      DefaultConfig(),
	}

	// Create extraction directory
	extractDir, err := os.MkdirTemp(tempDir, "extract_*")
	if err != nil {
		t.Fatal(err)
	}
	a.extractDir = extractDir

	err = a.extractTar(nil)
	if err != nil {
		t.Errorf("extractTar() error = %v, want nil", err)
		return
	}

	// Verify extracted files
	testFile := filepath.Join(extractDir, "testfile.txt")
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Error("extractTar() did not extract testfile.txt")
	}

	testDir := filepath.Join(extractDir, "testdir")
	if info, err := os.Stat(testDir); os.IsNotExist(err) {
		t.Error("extractTar() did not extract testdir")
	} else if !info.IsDir() {
		t.Error("extractTar() testdir should be a directory")
	}
}

// Helper function to create a simple ZIP archive for testing
func createTestZipArchive(t *testing.T, tempDir string) string {
	t.Helper()
	
	archivePath := filepath.Join(tempDir, "test.zip")
	file, err := os.Create(archivePath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	zw := zip.NewWriter(file)
	defer zw.Close()

	// Add a test file to the archive
	fw, err := zw.Create("testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	
	if _, err := fw.Write([]byte("test content")); err != nil {
		t.Fatal(err)
	}

	// Add a directory
	_, err = zw.Create("testdir/")
	if err != nil {
		t.Fatal(err)
	}

	return archivePath
}

func TestArchiveContextImpl_extractZip_Success(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "zip_extract_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create test zip archive
	archivePath := createTestZipArchive(t, tempDir)

	a := &ArchiveContextImpl{
		archivePath: archivePath,
		config:      DefaultConfig(),
	}

	// Create extraction directory
	extractDir, err := os.MkdirTemp(tempDir, "extract_*")
	if err != nil {
		t.Fatal(err)
	}
	a.extractDir = extractDir

	err = a.extractZip()
	if err != nil {
		t.Errorf("extractZip() error = %v, want nil", err)
		return
	}

	// Verify extracted files
	testFile := filepath.Join(extractDir, "testfile.txt")
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Error("extractZip() did not extract testfile.txt")
	}

	testDir := filepath.Join(extractDir, "testdir")
	if info, err := os.Stat(testDir); os.IsNotExist(err) {
		t.Error("extractZip() did not extract testdir")
	} else if !info.IsDir() {
		t.Error("extractZip() testdir should be a directory")
	}
}

func TestArchiveContextImpl_extractTarGz_Success(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "targz_extract_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create a gzipped tar archive
	archivePath := filepath.Join(tempDir, "test.tar.gz")
	file, err := os.Create(archivePath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	gzWriter := gzip.NewWriter(file)
	defer gzWriter.Close()

	tw := tar.NewWriter(gzWriter)
	defer tw.Close()

	// Add a test file
	header := &tar.Header{
		Name: "testfile.txt",
		Mode: 0644,
		Size: int64(len("test content")),
		Typeflag: tar.TypeReg,
	}
	
	if err := tw.WriteHeader(header); err != nil {
		t.Fatal(err)
	}
	
	if _, err := tw.Write([]byte("test content")); err != nil {
		t.Fatal(err)
	}

	// Close writers to flush content
	tw.Close()
	gzWriter.Close()
	file.Close()

	// Test extraction
	a := &ArchiveContextImpl{
		archivePath: archivePath,
		config:      DefaultConfig(),
	}

	extractDir, err := os.MkdirTemp(tempDir, "extract_*")
	if err != nil {
		t.Fatal(err)
	}
	a.extractDir = extractDir

	err = a.extractTarGz()
	if err != nil {
		t.Errorf("extractTarGz() error = %v, want nil", err)
		return
	}

	// Verify extracted file
	testFile := filepath.Join(extractDir, "testfile.txt")
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Error("extractTarGz() did not extract testfile.txt")
	}

	// Verify content
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}
	
	if string(content) != "test content" {
		t.Errorf("extractTarGz() content = %v, want %v", string(content), "test content")
	}
}

func TestArchiveContextImpl_Extract_Integration(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "extract_integration_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Test with tar.gz format
	archivePath := filepath.Join(tempDir, "test.tar.gz")
	file, err := os.Create(archivePath)
	if err != nil {
		t.Fatal(err)
	}

	gzWriter := gzip.NewWriter(file)
	tw := tar.NewWriter(gzWriter)

	// Add test content
	header := &tar.Header{
		Name: "integration/test.txt",
		Mode: 0644,
		Size: int64(len("integration test")),
		Typeflag: tar.TypeReg,
	}
	
	tw.WriteHeader(header)
	tw.Write([]byte("integration test"))
	tw.Close()
	gzWriter.Close()
	file.Close()

	// Test full Extract() workflow
	a := &ArchiveContextImpl{
		archivePath: archivePath,
		config:      DefaultConfig(),
	}

	extractedPath, err := a.Extract()
	if err != nil {
		t.Errorf("Extract() error = %v, want nil", err)
		return
	}

	if extractedPath == "" {
		t.Error("Extract() should return extraction path")
	}

	// Verify extracted content
	testFile := filepath.Join(extractedPath, "integration", "test.txt")
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	if string(content) != "integration test" {
		t.Errorf("Extract() content = %v, want %v", string(content), "integration test")
	}
}

func TestArchiveContextImpl_extractZipFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "zip_file_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create a simple zip file
	zipPath := filepath.Join(tempDir, "test.zip")
	zipFile, err := os.Create(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	defer zipFile.Close()

	zw := zip.NewWriter(zipFile)
	fw, err := zw.Create("test.txt")
	if err != nil {
		t.Fatal(err)
	}
	fw.Write([]byte("zip file content"))
	zw.Close()
	zipFile.Close()

	// Open the zip for reading
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	if len(reader.File) == 0 {
		t.Fatal("Zip file should contain at least one file")
	}

	a := &ArchiveContextImpl{
		config: DefaultConfig(),
	}

	targetPath := filepath.Join(tempDir, "extracted.txt")
	err = a.extractZipFile(reader.File[0], targetPath)
	if err != nil {
		t.Errorf("extractZipFile() error = %v, want nil", err)
		return
	}

	// Verify extracted content
	content, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatal(err)
	}

	if string(content) != "zip file content" {
		t.Errorf("extractZipFile() content = %v, want %v", string(content), "zip file content")
	}
}

func TestArchiveContextImpl_extractTarBz2_Success(t *testing.T) {
	// This tests the bzip2 extraction path which may not be covered
	// Since we can't easily create a real bzip2 file, we'll test the method call path
	tempDir, err := os.MkdirTemp("", "bz2_extract_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create a mock bzip2 file (won't be valid, but tests method path)
	archivePath := filepath.Join(tempDir, "test.tar.bz2")
	if err := os.WriteFile(archivePath, []byte("BZ fake content"), 0644); err != nil {
		t.Fatal(err)
	}

	a := &ArchiveContextImpl{
		archivePath: archivePath,
		config:      DefaultConfig(),
	}

	// This will fail because it's not a real bz2 file, but covers the method
	err = a.extractTarBz2()
	if err == nil {
		t.Error("extractTarBz2() should fail with fake bz2 content")
	}
}

func TestArchiveContextImpl_FullWorkflow_ErrorCases(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "error_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Test Extract() with unsupported format
	a := &ArchiveContextImpl{
		archivePath: "/nonexistent/path",
		config:      DefaultConfig(),
	}

	_, err = a.Extract()
	if err == nil {
		t.Error("Extract() should fail with nonexistent file")
	}
}

func TestContextType_String(t *testing.T) {
	// Test the String() method of ContextType enum
	tests := []struct {
		ct       ContextType
		expected string
	}{
		{LocalContext, "local"},
		{URLContext, "url"},
		{GitContext, "git"},
		{ArchiveContext, "archive"},
		{ContextType(999), "unknown"}, // Invalid type
	}

	for _, tt := range tests {
		if got := tt.ct.String(); got != tt.expected {
			t.Errorf("ContextType.String() = %v, want %v", got, tt.expected)
		}
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	
	if config.MaxSize != 500*1024*1024 {
		t.Errorf("DefaultConfig().MaxSize = %v, want %v", config.MaxSize, 500*1024*1024)
	}
	
	if !config.CacheEnabled {
		t.Error("DefaultConfig().CacheEnabled should be true")
	}
	
	if config.TempDir != "/tmp" {
		t.Errorf("DefaultConfig().TempDir = %v, want %v", config.TempDir, "/tmp")
	}
	
	if config.HTTPTimeout != 30*time.Second {
		t.Errorf("DefaultConfig().HTTPTimeout = %v, want %v", config.HTTPTimeout, 30*time.Second)
	}
}

func TestArchiveContextImpl_extractTarReader_ErrorHandling(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "tar_reader_error_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	a := &ArchiveContextImpl{
		extractDir: tempDir,
		config:     DefaultConfig(),
	}

	// Create a tar reader with an invalid path entry
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	
	// Add an entry that will cause path validation to fail
	header := &tar.Header{
		Name:     "../../etc/passwd", // This should fail path validation
		Mode:     0644,
		Size:     5,
		Typeflag: tar.TypeReg,
	}
	
	tw.WriteHeader(header)
	tw.Write([]byte("test\n"))
	tw.Close()

	tr := tar.NewReader(&buf)
	
	// This should fail due to path traversal attempt
	err = a.extractTarReader(tr)
	if err == nil {
		t.Error("extractTarReader() should fail with path traversal attempt")
	}
	if !strings.Contains(err.Error(), "path traversal attempt") {
		t.Errorf("extractTarReader() error should mention path traversal, got: %v", err)
	}
}