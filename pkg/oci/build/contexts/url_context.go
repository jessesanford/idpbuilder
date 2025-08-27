package contexts

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

// URLContextImpl implements Context for HTTP/HTTPS URL sources
type URLContextImpl struct {
	url        *url.URL
	config     *ContextConfig
	httpClient HTTPClient
	localPath  string
}

// Path returns the local path where the URL content was downloaded
func (u *URLContextImpl) Path() string {
	return u.localPath
}

// Type returns URLContext
func (u *URLContextImpl) Type() ContextType {
	return URLContext
}

// Cleanup removes any temporary files created for this context
func (u *URLContextImpl) Cleanup() error {
	if u.localPath != "" {
		return os.RemoveAll(u.localPath)
	}
	return nil
}

// Fetch downloads the URL content and prepares it as a build context
func (u *URLContextImpl) Fetch() error {
	// Check cache first if enabled
	if u.config.CacheEnabled {
		cachedPath, err := u.getCachedPath()
		if err == nil {
			if _, err := os.Stat(cachedPath); err == nil {
				u.localPath = cachedPath
				return nil
			}
		}
	}

	// Download the content
	content, contentType, err := u.downloadWithRetry(3)
	if err != nil {
		return fmt.Errorf("failed to download %s: %w", u.url.String(), err)
	}

	// Determine if it's an archive or directory
	tempDir, err := u.createTempDir()
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}

	u.localPath = tempDir

	// Check if the content appears to be an archive
	if u.isArchiveContent(contentType, content) {
		return u.extractArchiveContent(content, tempDir)
	}

	// Treat as a single file context
	fileName := u.getFileName()
	filePath := filepath.Join(tempDir, fileName)
	
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		return fmt.Errorf("failed to write downloaded file: %w", err)
	}

	// Cache the result if enabled
	if u.config.CacheEnabled {
		u.cacheContent(content)
	}

	return nil
}

// downloadWithRetry attempts to download the URL with retries
func (u *URLContextImpl) downloadWithRetry(maxRetries int) ([]byte, string, error) {
	var lastErr error
	
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Duration(attempt) * time.Second)
		}

		req, err := http.NewRequest("GET", u.url.String(), nil)
		if err != nil {
			lastErr = fmt.Errorf("failed to create request: %w", err)
			continue
		}

		resp, err := u.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("request failed: %w", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("HTTP error: %d %s", resp.StatusCode, resp.Status)
			continue
		}

		// Check content length against max size
		if resp.ContentLength > 0 && resp.ContentLength > u.config.MaxSize {
			return nil, "", fmt.Errorf("content too large: %d bytes (max: %d)", 
				resp.ContentLength, u.config.MaxSize)
		}

		content, err := io.ReadAll(io.LimitReader(resp.Body, u.config.MaxSize+1))
		if err != nil {
			lastErr = fmt.Errorf("failed to read response: %w", err)
			continue
		}

		if int64(len(content)) > u.config.MaxSize {
			return nil, "", fmt.Errorf("content too large: %d bytes (max: %d)", 
				len(content), u.config.MaxSize)
		}

		contentType := resp.Header.Get("Content-Type")
		return content, contentType, nil
	}

	return nil, "", fmt.Errorf("all %d download attempts failed, last error: %w", 
		maxRetries, lastErr)
}

// isArchiveContent determines if the downloaded content is an archive
func (u *URLContextImpl) isArchiveContent(contentType string, content []byte) bool {
	// Check Content-Type header
	archiveTypes := []string{
		"application/x-tar",
		"application/x-gzip",
		"application/gzip",
		"application/zip",
		"application/x-bzip2",
	}
	
	for _, archiveType := range archiveTypes {
		if contentType == archiveType {
			return true
		}
	}

	// Check magic bytes for common archive formats
	if len(content) < 4 {
		return false
	}

	// ZIP magic number
	if content[0] == 0x50 && content[1] == 0x4B {
		return true
	}

	// TAR files (check for ustar magic)
	if len(content) >= 265 {
		ustarMagic := string(content[257:262])
		if ustarMagic == "ustar" {
			return true
		}
	}

	// GZIP magic number
	if content[0] == 0x1F && content[1] == 0x8B {
		return true
	}

	return false
}

// extractArchiveContent extracts archive content to the temp directory
func (u *URLContextImpl) extractArchiveContent(content []byte, tempDir string) error {
	// For Split 001, we'll treat archive content as a single file
	// Full archive extraction will be handled by ArchiveContextImpl in Split 002
	
	// Write archive to temp directory as-is
	fileName := u.getFileName()
	if fileName == "content" {
		fileName = "archive"
	}
	
	filePath := filepath.Join(tempDir, fileName)
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		return fmt.Errorf("failed to write archive file: %w", err)
	}

	return nil
}


// createTempDir creates a temporary directory for this context
func (u *URLContextImpl) createTempDir() (string, error) {
	return os.MkdirTemp(u.config.TempDir, "url_context_*")
}

// getFileName extracts a reasonable filename from the URL
func (u *URLContextImpl) getFileName() string {
	path := u.url.Path
	if path == "" || path == "/" {
		return "index.html"
	}
	
	fileName := filepath.Base(path)
	if fileName == "." || fileName == "/" {
		return "content"
	}
	
	return fileName
}

// getCachedPath returns the path where cached content would be stored
func (u *URLContextImpl) getCachedPath() (string, error) {
	hash := sha256.Sum256([]byte(u.url.String()))
	cacheKey := fmt.Sprintf("%x", hash[:8]) // Use first 8 bytes as cache key
	
	cacheDir := filepath.Join(u.config.TempDir, "url_cache")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", err
	}
	
	return filepath.Join(cacheDir, cacheKey), nil
}

// cacheContent stores content in cache for future use
func (u *URLContextImpl) cacheContent(content []byte) error {
	if !u.config.CacheEnabled {
		return nil
	}
	
	cachePath, err := u.getCachedPath()
	if err != nil {
		return err // Non-fatal error, just log it
	}
	
	return os.WriteFile(cachePath, content, 0644)
}