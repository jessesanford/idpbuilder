package contexts

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// ContextEntry represents a file or directory in the build context
type ContextEntry struct {
	Path    string
	Size    int64
	ModTime time.Time
	IsDir   bool
}

// Filter represents a dockerignore-style filter
type Filter struct {
	patterns []string
}

// NewFilter creates a new filter instance
func NewFilter() *Filter {
	return &Filter{patterns: []string{}}
}

// ParseDockerignore reads and parses a .dockerignore file
func ParseDockerignore(dockerignorePath string) (*Filter, error) {
	filter := NewFilter()
	file, err := os.Open(dockerignorePath)
	if err != nil {
		if os.IsNotExist(err) {
			return filter, nil
		}
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			filter.patterns = append(filter.patterns, line)
		}
	}
	return filter, scanner.Err()
}

// AddPattern adds a pattern to the filter
func (f *Filter) AddPattern(pattern string) {
	f.patterns = append(f.patterns, pattern)
}

// ApplyFilter filters entries based on dockerignore patterns
func (f *Filter) ApplyFilter(entries []ContextEntry) []ContextEntry {
	if len(f.patterns) == 0 {
		return entries
	}
	var filtered []ContextEntry
	for _, entry := range entries {
		if !f.shouldExclude(entry.Path) {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

// shouldExclude checks if a path should be excluded
func (f *Filter) shouldExclude(path string) bool {
	normalizedPath := filepath.ToSlash(path)
	excluded := false
	for _, pattern := range f.patterns {
		if strings.HasPrefix(pattern, "!") {
			if f.matchesPattern(normalizedPath, strings.TrimPrefix(pattern, "!")) {
				excluded = false
			}
		} else if f.matchesPattern(normalizedPath, pattern) {
			excluded = true
		}
	}
	return excluded
}

// matchesPattern checks if a path matches a dockerignore pattern
func (f *Filter) matchesPattern(path, pattern string) bool {
	escaped := regexp.QuoteMeta(pattern)
	escaped = strings.ReplaceAll(escaped, "\\*", ".*")
	escaped = "^" + escaped + "$"
	regex, err := regexp.Compile(escaped)
	return err == nil && regex.MatchString(path)
}