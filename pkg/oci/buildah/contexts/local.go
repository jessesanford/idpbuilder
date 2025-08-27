package contexts

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LocalContext implements Context interface for local directory contexts
type LocalContext struct {
	sourcePath string
	excludes   []string
	entries    []*ContextEntry
}

// PrepareContext prepares the local directory context
func (lc *LocalContext) PrepareContext(ctx context.Context, opts ContextOptions) (string, error) {
	if _, err := os.Stat(lc.sourcePath); err != nil {
		return "", fmt.Errorf("source path invalid: %w", err)
	}
	return lc.sourcePath, nil
}

// GetEntries returns all entries in the local context
func (lc *LocalContext) GetEntries() ([]*ContextEntry, error) {
	if lc.entries != nil {
		return lc.entries, nil
	}
	
	var entries []*ContextEntry
	err := filepath.Walk(lc.sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		relPath, err := filepath.Rel(lc.sourcePath, path)
		if err != nil {
			return err
		}
		
		if lc.isExcluded(relPath) {
			return nil
		}
		
		entries = append(entries, &ContextEntry{
			Path:    relPath,
			Mode:    uint32(info.Mode()),
			Size:    info.Size(),
			ModTime: info.ModTime(),
			IsDir:   info.IsDir(),
		})
		return nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to walk source directory: %w", err)
	}
	
	lc.entries = entries
	return entries, nil
}

// GetEntry returns a specific entry by path
func (lc *LocalContext) GetEntry(path string) (*ContextEntry, error) {
	fullPath := filepath.Join(lc.sourcePath, path)
	info, err := os.Stat(fullPath)
	if err != nil {
		return nil, fmt.Errorf("entry not found: %w", err)
	}
	
	entry := &ContextEntry{
		Path:    path,
		Mode:    uint32(info.Mode()),
		Size:    info.Size(),
		ModTime: info.ModTime(),
		IsDir:   info.IsDir(),
	}
	
	if !info.IsDir() {
		file, _ := os.Open(fullPath)
		entry.Content = file
	}
	
	return entry, nil
}

// isExcluded checks if a path should be excluded
func (lc *LocalContext) isExcluded(path string) bool {
	for _, exclude := range lc.excludes {
		if matched, _ := filepath.Match(exclude, path); matched || strings.HasPrefix(path, exclude) {
			return true
		}
	}
	return false
}

// Cleanup removes any temporary resources
func (lc *LocalContext) Cleanup() error {
	lc.entries = nil
	return nil
}

// Type returns the context type
func (lc *LocalContext) Type() ContextType {
	return Local
}