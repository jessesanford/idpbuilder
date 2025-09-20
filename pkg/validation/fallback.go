package validation

import (
	"fmt"
	"sync"
)

// FallbackStrategy represents a validation fallback strategy
type FallbackStrategy interface {
	// Validate performs the fallback validation
	Validate(data interface{}) error
	// Name returns the strategy name
	Name() string
	// Priority returns the strategy priority (lower = higher priority)
	Priority() int
}

// FallbackRegistry manages validation fallback strategies
type FallbackRegistry struct {
	mutex      sync.RWMutex
	strategies map[string]FallbackStrategy
}

// NewFallbackRegistry creates a new fallback registry
func NewFallbackRegistry() *FallbackRegistry {
	return &FallbackRegistry{
		strategies: make(map[string]FallbackStrategy),
	}
}

// Register adds a fallback strategy to the registry
func (fr *FallbackRegistry) Register(strategy FallbackStrategy) error {
	if strategy == nil {
		return fmt.Errorf("strategy cannot be nil")
	}

	name := strategy.Name()
	if name == "" {
		return fmt.Errorf("strategy name cannot be empty")
	}

	fr.mutex.Lock()
	defer fr.mutex.Unlock()

	if _, exists := fr.strategies[name]; exists {
		return fmt.Errorf("strategy %q is already registered", name)
	}

	fr.strategies[name] = strategy
	return nil
}

// Get retrieves a fallback strategy by name
func (fr *FallbackRegistry) Get(name string) (FallbackStrategy, error) {
	if name == "" {
		return nil, fmt.Errorf("strategy name cannot be empty")
	}

	fr.mutex.RLock()
	defer fr.mutex.RUnlock()

	strategy, exists := fr.strategies[name]
	if !exists {
		return nil, fmt.Errorf("strategy %q not found", name)
	}

	return strategy, nil
}

// List returns all registered fallback strategy names
func (fr *FallbackRegistry) List() []string {
	fr.mutex.RLock()
	defer fr.mutex.RUnlock()

	names := make([]string, 0, len(fr.strategies))
	for name := range fr.strategies {
		names = append(names, name)
	}

	return names
}

// GetByPriority returns strategies sorted by priority (lowest first)
func (fr *FallbackRegistry) GetByPriority() []FallbackStrategy {
	fr.mutex.RLock()
	defer fr.mutex.RUnlock()

	strategies := make([]FallbackStrategy, 0, len(fr.strategies))
	for _, strategy := range fr.strategies {
		strategies = append(strategies, strategy)
	}

	// Sort by priority (simple bubble sort for small collections)
	for i := 0; i < len(strategies)-1; i++ {
		for j := 0; j < len(strategies)-i-1; j++ {
			if strategies[j].Priority() > strategies[j+1].Priority() {
				strategies[j], strategies[j+1] = strategies[j+1], strategies[j]
			}
		}
	}

	return strategies
}

// Remove removes a strategy from the registry
func (fr *FallbackRegistry) Remove(name string) error {
	if name == "" {
		return fmt.Errorf("strategy name cannot be empty")
	}

	fr.mutex.Lock()
	defer fr.mutex.Unlock()

	if _, exists := fr.strategies[name]; !exists {
		return fmt.Errorf("strategy %q not found", name)
	}

	delete(fr.strategies, name)
	return nil
}

// Count returns the number of registered strategies
func (fr *FallbackRegistry) Count() int {
	fr.mutex.RLock()
	defer fr.mutex.RUnlock()
	return len(fr.strategies)
}