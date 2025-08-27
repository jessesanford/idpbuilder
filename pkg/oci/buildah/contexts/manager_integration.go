package contexts

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"
)

// ContextType represents different types of build contexts
type ContextType string

const (
	FileContext   ContextType = "file"
	GitContext    ContextType = "git"
	HTTPContext   ContextType = "http"
	TarContext    ContextType = "tar"
	DockerContext ContextType = "docker"
)

// ContextFactory creates context instances based on type and configuration
type ContextFactory struct {
	validators map[ContextType]ValidatorFunc
	mu         sync.RWMutex
}

// ValidatorFunc defines the signature for context validators
type ValidatorFunc func(source string, config *ContextConfig) error

// ContextConfig holds configuration for context creation
type ContextConfig struct {
	Type         ContextType
	Source       string
	AllowSymlinks bool
	Timeout      time.Duration
	Credentials  map[string]string
	Options      map[string]interface{}
}

// BuildContext represents a build context with validation and cleanup
type BuildContext interface {
	// GetData returns the context data as a reader
	GetData(ctx context.Context) (io.ReadCloser, error)
	
	// GetType returns the context type
	GetType() ContextType
	
	// Cleanup performs cleanup operations
	Cleanup() error
	
	// Validate performs context validation
	Validate() error
}

// ContextManager manages build contexts with caching and lifecycle
type ContextManager struct {
	factory    *ContextFactory
	cache      map[string]BuildContext
	cacheMu    sync.RWMutex
	maxCacheSize int
	defaultTimeout time.Duration
}

// NewContextFactory creates a new context factory with default validators
func NewContextFactory() *ContextFactory {
	factory := &ContextFactory{
		validators: make(map[ContextType]ValidatorFunc),
	}
	
	// Register default validators
	factory.RegisterValidator(FileContext, validateFileContext)
	factory.RegisterValidator(TarContext, validateTarContext)
	
	return factory
}

// RegisterValidator registers a validator for a specific context type
func (cf *ContextFactory) RegisterValidator(contextType ContextType, validator ValidatorFunc) {
	cf.mu.Lock()
	defer cf.mu.Unlock()
	cf.validators[contextType] = validator
}

// CreateContext creates a new build context based on configuration
func (cf *ContextFactory) CreateContext(config *ContextConfig) (BuildContext, error) {
	// Validate configuration
	if err := cf.validateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid context configuration: %w", err)
	}
	
	// Run type-specific validation
	cf.mu.RLock()
	validator, exists := cf.validators[config.Type]
	cf.mu.RUnlock()
	
	if exists {
		if err := validator(config.Source, config); err != nil {
			return nil, fmt.Errorf("context validation failed: %w", err)
		}
	}
	
	// Create context based on type
	switch config.Type {
	case FileContext:
		return newFileContext(config)
	case TarContext:
		return newTarContext(config)
	default:
		return nil, fmt.Errorf("unsupported context type: %s", config.Type)
	}
}

// NewContextManager creates a new context manager
func NewContextManager(opts ...ContextManagerOption) *ContextManager {
	manager := &ContextManager{
		factory:        NewContextFactory(),
		cache:          make(map[string]BuildContext),
		maxCacheSize:   10,
		defaultTimeout: 5 * time.Minute,
	}
	
	for _, opt := range opts {
		opt(manager)
	}
	
	return manager
}

// ContextManagerOption defines options for context manager
type ContextManagerOption func(*ContextManager)

// WithMaxCacheSize sets the maximum cache size
func WithMaxCacheSize(size int) ContextManagerOption {
	return func(cm *ContextManager) {
		cm.maxCacheSize = size
	}
}

// WithDefaultTimeout sets the default timeout for operations
func WithDefaultTimeout(timeout time.Duration) ContextManagerOption {
	return func(cm *ContextManager) {
		cm.defaultTimeout = timeout
	}
}

// GetContext retrieves or creates a build context with caching
func (cm *ContextManager) GetContext(ctx context.Context, config *ContextConfig) (BuildContext, error) {
	cacheKey := cm.getCacheKey(config)
	
	// Check cache first
	cm.cacheMu.RLock()
	if cachedContext, exists := cm.cache[cacheKey]; exists {
		cm.cacheMu.RUnlock()
		return cachedContext, nil
	}
	cm.cacheMu.RUnlock()
	
	// Create new context
	buildContext, err := cm.factory.CreateContext(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create context: %w", err)
	}
	
	// Add to cache if there's space
	cm.addToCache(cacheKey, buildContext)
	
	return buildContext, nil
}

// CleanupCache removes expired contexts from cache
func (cm *ContextManager) CleanupCache() {
	cm.cacheMu.Lock()
	defer cm.cacheMu.Unlock()
	
	for key, context := range cm.cache {
		if err := context.Cleanup(); err == nil {
			delete(cm.cache, key)
		}
	}
}

// Shutdown performs cleanup of all contexts and resources
func (cm *ContextManager) Shutdown() error {
	cm.cacheMu.Lock()
	defer cm.cacheMu.Unlock()
	
	var lastErr error
	for _, context := range cm.cache {
		if err := context.Cleanup(); err != nil {
			lastErr = err
		}
	}
	
	cm.cache = make(map[string]BuildContext)
	return lastErr
}

// Helper methods

func (cf *ContextFactory) validateConfig(config *ContextConfig) error {
	if config == nil {
		return fmt.Errorf("configuration cannot be nil")
	}
	
	if config.Type == "" {
		return fmt.Errorf("context type is required")
	}
	
	if config.Source == "" {
		return fmt.Errorf("context source is required")
	}
	
	return nil
}

func (cm *ContextManager) getCacheKey(config *ContextConfig) string {
	return fmt.Sprintf("%s:%s", config.Type, config.Source)
}

func (cm *ContextManager) addToCache(key string, context BuildContext) {
	cm.cacheMu.Lock()
	defer cm.cacheMu.Unlock()
	
	// Remove oldest entry if cache is full
	if len(cm.cache) >= cm.maxCacheSize {
		for oldKey := range cm.cache {
			if oldContext := cm.cache[oldKey]; oldContext != nil {
				oldContext.Cleanup()
			}
			delete(cm.cache, oldKey)
			break
		}
	}
	
	cm.cache[key] = context
}

// Default validators

func validateFileContext(source string, config *ContextConfig) error {
	result, err := ValidateContext(source, config.AllowSymlinks)
	if err != nil {
		return err
	}
	
	if !result.Valid {
		return fmt.Errorf("context validation failed: %v", result.Errors)
	}
	
	return CheckContextPermissions(source)
}

func validateTarContext(source string, config *ContextConfig) error {
	// Placeholder for tar context validation
	return nil
}

// Simple context implementations for integration

type fileContext struct {
	config *ContextConfig
	path   string
}

func newFileContext(config *ContextConfig) (BuildContext, error) {
	return &fileContext{
		config: config,
		path:   config.Source,
	}, nil
}

func (fc *fileContext) GetData(ctx context.Context) (io.ReadCloser, error) {
	// Placeholder - would return compressed context data
	return nil, fmt.Errorf("not implemented")
}

func (fc *fileContext) GetType() ContextType {
	return FileContext
}

func (fc *fileContext) Cleanup() error {
	return nil
}

func (fc *fileContext) Validate() error {
	result, err := ValidateContext(fc.path, fc.config.AllowSymlinks)
	if err != nil {
		return err
	}
	
	if !result.Valid {
		return fmt.Errorf("validation failed: %v", result.Errors)
	}
	
	return nil
}

type tarContext struct {
	config *ContextConfig
}

func newTarContext(config *ContextConfig) (BuildContext, error) {
	return &tarContext{config: config}, nil
}

func (tc *tarContext) GetData(ctx context.Context) (io.ReadCloser, error) {
	return nil, fmt.Errorf("not implemented")
}

func (tc *tarContext) GetType() ContextType {
	return TarContext
}

func (tc *tarContext) Cleanup() error {
	return nil
}

func (tc *tarContext) Validate() error {
	return nil
}