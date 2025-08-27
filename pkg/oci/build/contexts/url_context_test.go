package contexts

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Mock HTTP client for testing
type mockHTTPClient struct {
	responses map[string]*mockResponse
	doFunc    func(req *http.Request) (*http.Response, error)
}

type mockResponse struct {
	statusCode    int
	body          []byte
	contentType   string
	contentLength int64
	err           error
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if m.doFunc != nil {
		return m.doFunc(req)
	}
	
	if response, ok := m.responses[req.URL.String()]; ok {
		if response.err != nil {
			return nil, response.err
		}
		
		resp := &http.Response{
			StatusCode:    response.statusCode,
			Body:          io.NopCloser(bytes.NewReader(response.body)),
			Header:        make(http.Header),
			ContentLength: response.contentLength,
		}
		
		if response.contentType != "" {
			resp.Header.Set("Content-Type", response.contentType)
		}
		if response.contentLength == 0 && len(response.body) > 0 {
			resp.ContentLength = int64(len(response.body))
		}
		
		return resp, nil
	}
	
	// Default response
	return &http.Response{
		StatusCode: http.StatusNotFound,
		Body:       io.NopCloser(strings.NewReader("Not Found")),
	}, nil
}

func TestURLContextImpl_Path(t *testing.T) {
	testURL, _ := url.Parse("http://example.com/test")
	ctx := &URLContextImpl{
		url:       testURL,
		localPath: "/tmp/test-path",
	}
	
	if ctx.Path() != "/tmp/test-path" {
		t.Errorf("URLContextImpl.Path() = %v, want %v", ctx.Path(), "/tmp/test-path")
	}
}

func TestURLContextImpl_Type(t *testing.T) {
	ctx := &URLContextImpl{}
	if ctx.Type() != URLContext {
		t.Errorf("URLContextImpl.Type() = %v, want %v", ctx.Type(), URLContext)
	}
}

func TestURLContextImpl_Cleanup(t *testing.T) {
	// Create a temporary directory to test cleanup
	tempDir, err := os.MkdirTemp("", "url_context_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	
	ctx := &URLContextImpl{
		localPath: tempDir,
	}
	
	// Verify directory exists
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		t.Fatal("Temp directory should exist before cleanup")
	}
	
	// Cleanup should remove the directory
	err = ctx.Cleanup()
	if err != nil {
		t.Errorf("URLContextImpl.Cleanup() error = %v", err)
	}
	
	// Verify directory is removed
	if _, err := os.Stat(tempDir); !os.IsNotExist(err) {
		t.Error("Temp directory should be removed after cleanup")
	}
}

func TestURLContextImpl_Cleanup_EmptyPath(t *testing.T) {
	ctx := &URLContextImpl{}
	
	err := ctx.Cleanup()
	if err != nil {
		t.Errorf("URLContextImpl.Cleanup() with empty path error = %v", err)
	}
}

func TestURLContextImpl_getFileName(t *testing.T) {
	tests := []struct {
		name     string
		urlStr   string
		expected string
	}{
		{"file with extension", "http://example.com/file.txt", "file.txt"},
		{"file without extension", "http://example.com/README", "README"},
		{"root path", "http://example.com/", "index.html"},
		{"empty path", "http://example.com", "index.html"},
		{"path with query", "http://example.com/file.zip?download=1", "file.zip"},
		{"path ending with slash", "http://example.com/dir/", "dir"},
		{"dot path", "http://example.com/.", "content"},
		{"complex path", "http://example.com/path/to/archive.tar.gz", "archive.tar.gz"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testURL, err := url.Parse(tt.urlStr)
			if err != nil {
				t.Fatalf("Failed to parse URL: %v", err)
			}
			
			ctx := &URLContextImpl{url: testURL}
			result := ctx.getFileName()
			
			if result != tt.expected {
				t.Errorf("getFileName() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestURLContextImpl_isArchiveContent(t *testing.T) {
	ctx := &URLContextImpl{}
	
	tests := []struct {
		name        string
		contentType string
		content     []byte
		expected    bool
	}{
		{
			name:        "zip content type",
			contentType: "application/zip",
			content:     []byte("test"),
			expected:    true,
		},
		{
			name:        "tar content type",
			contentType: "application/x-tar",
			content:     []byte("test"),
			expected:    true,
		},
		{
			name:        "gzip content type",
			contentType: "application/gzip",
			content:     []byte("test"),
			expected:    true,
		},
		{
			name:        "x-gzip content type",
			contentType: "application/x-gzip",
			content:     []byte("test"),
			expected:    true,
		},
		{
			name:        "bzip2 content type",
			contentType: "application/x-bzip2",
			content:     []byte("test"),
			expected:    true,
		},
		{
			name:        "text content type",
			contentType: "text/plain",
			content:     []byte("test"),
			expected:    false,
		},
		{
			name:        "zip magic bytes",
			contentType: "application/octet-stream",
			content:     []byte{0x50, 0x4B, 0x03, 0x04}, // ZIP magic
			expected:    true,
		},
		{
			name:        "gzip magic bytes",
			contentType: "",
			content:     []byte{0x1F, 0x8B, 0x08, 0x00}, // GZIP magic (need at least 4 bytes)
			expected:    true,
		},
		{
			name:        "tar ustar magic",
			contentType: "",
			content:     func() []byte {
				content := make([]byte, 265)
				copy(content[257:262], []byte("ustar"))
				return content
			}(), // TAR ustar at offset 257, need at least 265 bytes
			expected:    true,
		},
		{
			name:        "short content",
			contentType: "",
			content:     []byte("ab"), // Too short to check magic
			expected:    false,
		},
		{
			name:        "no magic match",
			contentType: "text/plain",
			content:     []byte("regular text content"),
			expected:    false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ctx.isArchiveContent(tt.contentType, tt.content)
			if result != tt.expected {
				t.Errorf("isArchiveContent() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestURLContextImpl_downloadWithRetry(t *testing.T) {
	config := &ContextConfig{
		MaxSize:     1024,
		HTTPTimeout: 5 * time.Second,
	}
	
	t.Run("successful download", func(t *testing.T) {
		testURL, _ := url.Parse("http://example.com/test.txt")
		testContent := []byte("test content")
		
		mockClient := &mockHTTPClient{
			responses: map[string]*mockResponse{
				"http://example.com/test.txt": {
					statusCode:    http.StatusOK,
					body:          testContent,
					contentType:   "text/plain",
					contentLength: int64(len(testContent)),
				},
			},
		}
		
		ctx := &URLContextImpl{
			url:        testURL,
			config:     config,
			httpClient: mockClient,
		}
		
		content, contentType, err := ctx.downloadWithRetry(3)
		if err != nil {
			t.Errorf("downloadWithRetry() error = %v", err)
		}
		
		if !bytes.Equal(content, testContent) {
			t.Errorf("downloadWithRetry() content = %v, want %v", content, testContent)
		}
		
		if contentType != "text/plain" {
			t.Errorf("downloadWithRetry() contentType = %v, want %v", contentType, "text/plain")
		}
	})
	
	t.Run("HTTP error status", func(t *testing.T) {
		testURL, _ := url.Parse("http://example.com/notfound.txt")
		
		mockClient := &mockHTTPClient{
			responses: map[string]*mockResponse{
				"http://example.com/notfound.txt": {
					statusCode: http.StatusNotFound,
					body:       []byte("Not Found"),
				},
			},
		}
		
		ctx := &URLContextImpl{
			url:        testURL,
			config:     config,
			httpClient: mockClient,
		}
		
		_, _, err := ctx.downloadWithRetry(1)
		if err == nil {
			t.Error("downloadWithRetry() should return error for 404")
		}
		
		if !strings.Contains(err.Error(), "404") {
			t.Errorf("downloadWithRetry() error = %v, should contain 404", err)
		}
	})
	
	t.Run("content too large via Content-Length", func(t *testing.T) {
		testURL, _ := url.Parse("http://example.com/large.txt")
		
		mockClient := &mockHTTPClient{
			responses: map[string]*mockResponse{
				"http://example.com/large.txt": {
					statusCode:    http.StatusOK,
					body:          []byte("small content"),
					contentType:   "text/plain",
					contentLength: 2048, // Exceeds MaxSize of 1024
				},
			},
		}
		
		ctx := &URLContextImpl{
			url:        testURL,
			config:     config,
			httpClient: mockClient,
		}
		
		_, _, err := ctx.downloadWithRetry(1)
		if err == nil {
			t.Error("downloadWithRetry() should return error for large content")
		}
		
		if !strings.Contains(err.Error(), "too large") {
			t.Errorf("downloadWithRetry() error = %v, should mention too large", err)
		}
	})
	
	t.Run("content too large via actual content", func(t *testing.T) {
		testURL, _ := url.Parse("http://example.com/large.txt")
		largeContent := make([]byte, 2048) // Exceeds MaxSize of 1024
		
		mockClient := &mockHTTPClient{
			responses: map[string]*mockResponse{
				"http://example.com/large.txt": {
					statusCode:    http.StatusOK,
					body:          largeContent,
					contentType:   "text/plain",
					contentLength: 0, // Unknown content length
				},
			},
		}
		
		ctx := &URLContextImpl{
			url:        testURL,
			config:     config,
			httpClient: mockClient,
		}
		
		_, _, err := ctx.downloadWithRetry(1)
		if err == nil {
			t.Error("downloadWithRetry() should return error for large content")
		}
		
		if !strings.Contains(err.Error(), "too large") {
			t.Errorf("downloadWithRetry() error = %v, should mention too large", err)
		}
	})
	
	t.Run("network error with retry", func(t *testing.T) {
		testURL, _ := url.Parse("http://example.com/retry.txt")
		
		attemptCount := 0
		mockClient := &mockHTTPClient{
			doFunc: func(req *http.Request) (*http.Response, error) {
				attemptCount++
				if attemptCount < 3 {
					return nil, fmt.Errorf("network error")
				}
				// Success on third attempt
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("success")),
					Header:     make(http.Header),
				}, nil
			},
		}
		
		ctx := &URLContextImpl{
			url:        testURL,
			config:     config,
			httpClient: mockClient,
		}
		
		content, _, err := ctx.downloadWithRetry(3)
		if err != nil {
			t.Errorf("downloadWithRetry() error = %v", err)
		}
		
		if string(content) != "success" {
			t.Errorf("downloadWithRetry() content = %v, want success", string(content))
		}
		
		if attemptCount != 3 {
			t.Errorf("downloadWithRetry() attempts = %v, want 3", attemptCount)
		}
	})
	
	t.Run("all retries fail", func(t *testing.T) {
		testURL, _ := url.Parse("http://example.com/fail.txt")
		
		mockClient := &mockHTTPClient{
			doFunc: func(req *http.Request) (*http.Response, error) {
				return nil, fmt.Errorf("network error")
			},
		}
		
		ctx := &URLContextImpl{
			url:        testURL,
			config:     config,
			httpClient: mockClient,
		}
		
		_, _, err := ctx.downloadWithRetry(2)
		if err == nil {
			t.Error("downloadWithRetry() should return error when all retries fail")
		}
		
		if !strings.Contains(err.Error(), "all 2 download attempts failed") {
			t.Errorf("downloadWithRetry() error = %v, should mention failed attempts", err)
		}
	})
}

func TestURLContextImpl_createTempDir(t *testing.T) {
	config := &ContextConfig{
		TempDir: "/tmp",
	}
	
	ctx := &URLContextImpl{
		config: config,
	}
	
	tempDir, err := ctx.createTempDir()
	if err != nil {
		t.Errorf("createTempDir() error = %v", err)
	}
	
	defer os.RemoveAll(tempDir)
	
	if !strings.HasPrefix(tempDir, "/tmp") {
		t.Errorf("createTempDir() = %v, should start with /tmp", tempDir)
	}
	
	if !strings.Contains(tempDir, "url_context_") {
		t.Errorf("createTempDir() = %v, should contain url_context_", tempDir)
	}
	
	// Verify directory was created
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		t.Error("createTempDir() should create the directory")
	}
}

func TestURLContextImpl_getCachedPath(t *testing.T) {
	config := &ContextConfig{
		TempDir: "/tmp",
	}
	
	testURL, _ := url.Parse("http://example.com/test.txt")
	ctx := &URLContextImpl{
		url:    testURL,
		config: config,
	}
	
	cachedPath, err := ctx.getCachedPath()
	if err != nil {
		t.Errorf("getCachedPath() error = %v", err)
	}
	
	expectedPrefix := "/tmp/url_cache/"
	if !strings.HasPrefix(cachedPath, expectedPrefix) {
		t.Errorf("getCachedPath() = %v, should start with %v", cachedPath, expectedPrefix)
	}
	
	// Should be consistent for same URL
	cachedPath2, err := ctx.getCachedPath()
	if err != nil {
		t.Errorf("getCachedPath() second call error = %v", err)
	}
	
	if cachedPath != cachedPath2 {
		t.Errorf("getCachedPath() inconsistent results: %v != %v", cachedPath, cachedPath2)
	}
}

func TestURLContextImpl_cacheContent(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "cache_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	config := &ContextConfig{
		TempDir:      tempDir,
		CacheEnabled: true,
	}
	
	testURL, _ := url.Parse("http://example.com/test.txt")
	ctx := &URLContextImpl{
		url:    testURL,
		config: config,
	}
	
	testContent := []byte("test content for caching")
	
	t.Run("cache enabled", func(t *testing.T) {
		err := ctx.cacheContent(testContent)
		if err != nil {
			t.Errorf("cacheContent() error = %v", err)
		}
		
		// Verify content was cached
		cachedPath, err := ctx.getCachedPath()
		if err != nil {
			t.Fatalf("getCachedPath() error = %v", err)
		}
		
		cachedContent, err := os.ReadFile(cachedPath)
		if err != nil {
			t.Errorf("Failed to read cached content: %v", err)
		}
		
		if !bytes.Equal(cachedContent, testContent) {
			t.Errorf("Cached content = %v, want %v", cachedContent, testContent)
		}
	})
	
	t.Run("cache disabled", func(t *testing.T) {
		config.CacheEnabled = false
		err := ctx.cacheContent(testContent)
		if err != nil {
			t.Errorf("cacheContent() with disabled cache error = %v", err)
		}
		// Should not error even when disabled
	})
}

func TestURLContextImpl_extractArchiveContent(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "extract_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	testURL, _ := url.Parse("http://example.com/archive.zip")
	ctx := &URLContextImpl{
		url: testURL,
	}
	
	testContent := []byte("archive content")
	
	err = ctx.extractArchiveContent(testContent, tempDir)
	if err != nil {
		t.Errorf("extractArchiveContent() error = %v", err)
	}
	
	// Verify file was created
	expectedFile := filepath.Join(tempDir, "archive.zip")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Error("extractArchiveContent() should create archive file")
	}
	
	// Verify content
	savedContent, err := os.ReadFile(expectedFile)
	if err != nil {
		t.Errorf("Failed to read saved content: %v", err)
	}
	
	if !bytes.Equal(savedContent, testContent) {
		t.Errorf("Saved content = %v, want %v", savedContent, testContent)
	}
}

func TestURLContextImpl_Fetch_Integration(t *testing.T) {
	config := &ContextConfig{
		MaxSize:      1024,
		CacheEnabled: false, // Disable cache for simpler testing
		TempDir:      "/tmp",
		HTTPTimeout:  5 * time.Second,
	}
	
	t.Run("fetch regular file", func(t *testing.T) {
		testURL, _ := url.Parse("http://example.com/test.txt")
		testContent := []byte("test file content")
		
		mockClient := &mockHTTPClient{
			responses: map[string]*mockResponse{
				"http://example.com/test.txt": {
					statusCode:  http.StatusOK,
					body:        testContent,
					contentType: "text/plain",
				},
			},
		}
		
		ctx := &URLContextImpl{
			url:        testURL,
			config:     config,
			httpClient: mockClient,
		}
		
		err := ctx.Fetch()
		if err != nil {
			t.Errorf("Fetch() error = %v", err)
		}
		
		if ctx.localPath == "" {
			t.Error("Fetch() should set localPath")
		}
		
		defer ctx.Cleanup()
		
		// Verify file was created
		expectedFile := filepath.Join(ctx.localPath, "test.txt")
		if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
			t.Error("Fetch() should create file")
		}
	})
	
	t.Run("fetch archive file", func(t *testing.T) {
		testURL, _ := url.Parse("http://example.com/archive.zip")
		testContent := []byte{0x50, 0x4B, 0x03, 0x04} // ZIP magic bytes
		
		mockClient := &mockHTTPClient{
			responses: map[string]*mockResponse{
				"http://example.com/archive.zip": {
					statusCode:  http.StatusOK,
					body:        testContent,
					contentType: "application/zip",
				},
			},
		}
		
		ctx := &URLContextImpl{
			url:        testURL,
			config:     config,
			httpClient: mockClient,
		}
		
		err := ctx.Fetch()
		if err != nil {
			t.Errorf("Fetch() error = %v", err)
		}
		
		defer ctx.Cleanup()
		
		// Should handle as archive
		expectedFile := filepath.Join(ctx.localPath, "archive.zip")
		if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
			t.Error("Fetch() should create archive file")
		}
	})
	
	t.Run("fetch with download failure", func(t *testing.T) {
		testURL, _ := url.Parse("http://example.com/fail.txt")
		
		mockClient := &mockHTTPClient{
			doFunc: func(req *http.Request) (*http.Response, error) {
				return nil, fmt.Errorf("network error")
			},
		}
		
		ctx := &URLContextImpl{
			url:        testURL,
			config:     config,
			httpClient: mockClient,
		}
		
		err := ctx.Fetch()
		if err == nil {
			t.Error("Fetch() should return error on download failure")
		}
		
		if !strings.Contains(err.Error(), "failed to download") {
			t.Errorf("Fetch() error = %v, should mention download failure", err)
		}
	})
}

// Test cache behavior in Fetch
func TestURLContextImpl_Fetch_WithCache(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "fetch_cache_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	config := &ContextConfig{
		MaxSize:      1024,
		CacheEnabled: true,
		TempDir:      tempDir,
		HTTPTimeout:  5 * time.Second,
	}
	
	testURL, _ := url.Parse("http://example.com/cached.txt")
	testContent := []byte("cached content")
	
	// Pre-populate cache
	ctx := &URLContextImpl{
		url:    testURL,
		config: config,
	}
	
	cachedPath, err := ctx.getCachedPath()
	if err != nil {
		t.Fatalf("getCachedPath() error = %v", err)
	}
	
	// Create cache directory and file
	cacheDir := filepath.Dir(cachedPath)
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		t.Fatalf("Failed to create cache dir: %v", err)
	}
	
	if err := os.WriteFile(cachedPath, testContent, 0644); err != nil {
		t.Fatalf("Failed to write cached file: %v", err)
	}
	
	// Mock client that should not be called
	mockClient := &mockHTTPClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			t.Error("HTTP client should not be called when cache hit occurs")
			return nil, fmt.Errorf("should not be called")
		},
	}
	
	ctx.httpClient = mockClient
	
	err = ctx.Fetch()
	if err != nil {
		t.Errorf("Fetch() with cache error = %v", err)
	}
	
	if ctx.localPath != cachedPath {
		t.Errorf("Fetch() should use cached path %v, got %v", cachedPath, ctx.localPath)
	}
}