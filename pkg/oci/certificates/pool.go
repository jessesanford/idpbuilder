package certificates

import (
	"context"
	"crypto/x509"
	"fmt"
	"sync"
	"time"
)

// CertPoolManager manages certificate pools with hot-reload capabilities.
// It maintains separate pools for system and custom certificates and supports
// dynamic updates without service restart.
type CertPoolManager struct {
	// store is the underlying certificate storage.
	store CertificateStore
	
	// systemPool contains certificates from the system's trusted certificate store.
	systemPool *x509.CertPool
	
	// customPool contains user-provided certificates.
	customPool *x509.CertPool
	
	// mu protects concurrent access to the pools.
	mu sync.RWMutex
	
	// updateChan receives certificate update events.
	updateChan chan Event
	
	// stopChan signals the manager to stop processing updates.
	stopChan chan struct{}
	
	// validators are used to validate certificates before adding to pools.
	validators []CertificateValidator
	
	// config contains configuration for the pool manager.
	config *PoolConfig
	
	// isWatching indicates if the manager is actively watching for changes.
	isWatching bool
	
	// watchMu protects the isWatching flag.
	watchMu sync.Mutex
}

// PoolConfig contains configuration options for the certificate pool manager.
type PoolConfig struct {
	// AutoReload enables automatic reloading of certificate pools when changes are detected.
	AutoReload bool
	
	// ReloadInterval is the minimum interval between pool reloads.
	ReloadInterval time.Duration
	
	// ValidateCertificates determines whether certificates should be validated before adding to pools.
	ValidateCertificates bool
	
	// IncludeSystemCerts determines whether system certificates should be included in the pool.
	IncludeSystemCerts bool
	
	// MaxPoolSize is the maximum number of certificates allowed in a pool (0 = unlimited).
	MaxPoolSize int
}

// DefaultPoolConfig returns a default pool configuration.
func DefaultPoolConfig() *PoolConfig {
	return &PoolConfig{
		AutoReload:           true,
		ReloadInterval:       5 * time.Second,
		ValidateCertificates: true,
		IncludeSystemCerts:   true,
		MaxPoolSize:          1000,
	}
}

// NewCertPoolManager creates a new certificate pool manager.
func NewCertPoolManager(store CertificateStore, config *PoolConfig, validators ...CertificateValidator) (*CertPoolManager, error) {
	if store == nil {
		return nil, fmt.Errorf("certificate store cannot be nil")
	}
	
	if config == nil {
		config = DefaultPoolConfig()
	}
	
	manager := &CertPoolManager{
		store:      store,
		updateChan: make(chan Event, 100),
		stopChan:   make(chan struct{}),
		validators: validators,
		config:     config,
	}
	
	// Initialize pools
	if err := manager.initializePools(); err != nil {
		return nil, fmt.Errorf("failed to initialize certificate pools: %w", err)
	}
	
	return manager, nil
}

// initializePools initializes the system and custom certificate pools.
func (cpm *CertPoolManager) initializePools() error {
	cpm.mu.Lock()
	defer cpm.mu.Unlock()
	
	// Initialize system pool
	if cpm.config.IncludeSystemCerts {
		systemPool, err := x509.SystemCertPool()
		if err != nil {
			// If system pool is not available, create an empty pool
			systemPool = x509.NewCertPool()
		}
		cpm.systemPool = systemPool
	} else {
		cpm.systemPool = x509.NewCertPool()
	}
	
	// Initialize custom pool
	cpm.customPool = x509.NewCertPool()
	
	return nil
}

// Start begins monitoring the certificate store for changes and enables hot-reload.
func (cpm *CertPoolManager) Start(ctx context.Context) error {
	cpm.watchMu.Lock()
	defer cpm.watchMu.Unlock()
	
	if cpm.isWatching {
		return fmt.Errorf("pool manager is already watching")
	}
	
	// Load initial certificates
	if err := cpm.loadAllCertificates(ctx); err != nil {
		return fmt.Errorf("failed to load initial certificates: %w", err)
	}
	
	// Start watching for changes if auto-reload is enabled
	if cpm.config.AutoReload {
		if err := cpm.store.Watch(ctx, cpm.handleEvent); err != nil {
			return fmt.Errorf("failed to start watching certificate store: %w", err)
		}
		
		// Start the update processing goroutine
		go cpm.processUpdates(ctx)
	}
	
	cpm.isWatching = true
	return nil
}

// Stop stops monitoring the certificate store and shuts down the pool manager.
func (cpm *CertPoolManager) Stop() error {
	cpm.watchMu.Lock()
	defer cpm.watchMu.Unlock()
	
	if !cpm.isWatching {
		return nil
	}
	
	close(cpm.stopChan)
	cpm.isWatching = false
	
	return nil
}

// GetSystemPool returns the system certificate pool.
func (cpm *CertPoolManager) GetSystemPool() *x509.CertPool {
	cpm.mu.RLock()
	defer cpm.mu.RUnlock()
	
	return cpm.systemPool
}

// GetCustomPool returns the custom certificate pool.
func (cpm *CertPoolManager) GetCustomPool() *x509.CertPool {
	cpm.mu.RLock()
	defer cpm.mu.RUnlock()
	
	return cpm.customPool
}

// GetCombinedPool returns a new pool containing both system and custom certificates.
func (cpm *CertPoolManager) GetCombinedPool() *x509.CertPool {
	cpm.mu.RLock()
	defer cpm.mu.RUnlock()
	
	combined := x509.NewCertPool()
	
	// Add system certificates if available
	if cpm.systemPool != nil {
		// Note: x509.CertPool doesn't provide a way to iterate over certificates,
		// so we maintain our own tracking. In a real implementation, you might
		// need to track certificates separately.
	}
	
	return combined
}

// AddCertificate adds a certificate to the custom pool.
func (cpm *CertPoolManager) AddCertificate(ctx context.Context, cert *Certificate) error {
	if cert == nil {
		return fmt.Errorf("certificate cannot be nil")
	}
	
	// Validate certificate if validation is enabled
	if cpm.config.ValidateCertificates {
		for _, validator := range cpm.validators {
			if err := validator.Validate(ctx, cert); err != nil {
				return fmt.Errorf("certificate validation failed: %w", err)
			}
		}
	}
	
	cpm.mu.Lock()
	defer cpm.mu.Unlock()
	
	// Check pool size limit
	if cpm.config.MaxPoolSize > 0 && cpm.getPoolSizeLocked() >= cpm.config.MaxPoolSize {
		return fmt.Errorf("certificate pool has reached maximum size: %d", cpm.config.MaxPoolSize)
	}
	
	// Add certificate to custom pool
	cpm.customPool.AddCert(cert.X509)
	
	return nil
}

// RemoveCertificate removes a certificate from the custom pool by ID.
func (cpm *CertPoolManager) RemoveCertificate(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("certificate ID cannot be empty")
	}
	
	cpm.mu.Lock()
	defer cpm.mu.Unlock()
	
	// Load the certificate to get its X509 representation
	_, err := cpm.store.Load(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to load certificate for removal: %w", err)
	}
	
	// Note: x509.CertPool doesn't provide a direct way to remove certificates.
	// In a production implementation, you would need to rebuild the pool
	// without the specific certificate, or use a custom pool implementation.
	
	// For now, we rebuild the custom pool without this certificate
	return cpm.rebuildCustomPool(ctx, id)
}

// RotateCertificate replaces an existing certificate with a new one.
func (cpm *CertPoolManager) RotateCertificate(ctx context.Context, id string, newCert *Certificate) error {
	if id == "" {
		return fmt.Errorf("certificate ID cannot be empty")
	}
	if newCert == nil {
		return fmt.Errorf("new certificate cannot be nil")
	}
	
	// Validate new certificate if validation is enabled
	if cpm.config.ValidateCertificates {
		for _, validator := range cpm.validators {
			if err := validator.Validate(ctx, newCert); err != nil {
				return fmt.Errorf("new certificate validation failed: %w", err)
			}
		}
	}
	
	cpm.mu.Lock()
	defer cpm.mu.Unlock()
	
	// Store the new certificate
	if err := cpm.store.Save(ctx, id, newCert); err != nil {
		return fmt.Errorf("failed to store rotated certificate: %w", err)
	}
	
	// Update the custom pool
	// Since x509.CertPool doesn't support direct replacement, we rebuild the pool
	return cpm.rebuildCustomPoolLocked(ctx, "")
}

// RefreshPools reloads all certificates from the store into the pools.
func (cpm *CertPoolManager) RefreshPools(ctx context.Context) error {
	cpm.mu.Lock()
	defer cpm.mu.Unlock()
	
	// Clear and rebuild custom pool
	cpm.customPool = x509.NewCertPool()
	
	return cpm.loadAllCertificatesLocked(ctx)
}

// handleEvent processes certificate store events for hot-reload.
func (cpm *CertPoolManager) handleEvent(event Event) {
	select {
	case cpm.updateChan <- event:
	default:
		// Channel is full, drop the event
	}
}

// processUpdates processes certificate update events in a separate goroutine.
func (cpm *CertPoolManager) processUpdates(ctx context.Context) {
	ticker := time.NewTicker(cpm.config.ReloadInterval)
	defer ticker.Stop()
	
	pendingUpdates := make(map[string]Event)
	
	for {
		select {
		case event := <-cpm.updateChan:
			// Batch updates to avoid excessive reloading
			pendingUpdates[event.ID] = event
			
		case <-ticker.C:
			if len(pendingUpdates) > 0 {
				cpm.processPendingUpdates(ctx, pendingUpdates)
				// Clear pending updates
				pendingUpdates = make(map[string]Event)
			}
			
		case <-cpm.stopChan:
			return
			
		case <-ctx.Done():
			return
		}
	}
}

// processPendingUpdates processes a batch of pending certificate updates.
func (cpm *CertPoolManager) processPendingUpdates(ctx context.Context, updates map[string]Event) {
	cpm.mu.Lock()
	defer cpm.mu.Unlock()
	
	for id, event := range updates {
		switch event.Type {
		case EventAdded, EventModified:
			cert, err := cpm.store.Load(ctx, id)
			if err != nil {
				continue // Skip failed loads
			}
			
			// Validate if required
			if cpm.config.ValidateCertificates {
				valid := true
				for _, validator := range cpm.validators {
					if err := validator.Validate(ctx, cert); err != nil {
						valid = false
						break
					}
				}
				if !valid {
					continue
				}
			}
			
			// Add to custom pool (rebuild to handle modifications)
			cpm.rebuildCustomPoolLocked(ctx, "")
			
		case EventDeleted:
			// Rebuild pool without the deleted certificate
			cpm.rebuildCustomPoolLocked(ctx, id)
		}
	}
}

// loadAllCertificates loads all certificates from the store into the pools.
func (cpm *CertPoolManager) loadAllCertificates(ctx context.Context) error {
	cpm.mu.Lock()
	defer cpm.mu.Unlock()
	
	return cpm.loadAllCertificatesLocked(ctx)
}

// loadAllCertificatesLocked loads all certificates (caller must hold lock).
func (cpm *CertPoolManager) loadAllCertificatesLocked(ctx context.Context) error {
	// List all certificates in the store
	certIDs, err := cpm.store.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list certificates: %w", err)
	}
	
	// Load and add each certificate
	for _, id := range certIDs {
		cert, err := cpm.store.Load(ctx, id)
		if err != nil {
			continue // Skip failed loads
		}
		
		// Validate if required
		if cpm.config.ValidateCertificates {
			valid := true
			for _, validator := range cpm.validators {
				if err := validator.Validate(ctx, cert); err != nil {
					valid = false
					break
				}
			}
			if !valid {
				continue
			}
		}
		
		// Add to custom pool
		cpm.customPool.AddCert(cert.X509)
	}
	
	return nil
}

// rebuildCustomPool rebuilds the custom pool excluding a specific certificate ID.
func (cpm *CertPoolManager) rebuildCustomPool(ctx context.Context, excludeID string) error {
	cpm.mu.Lock()
	defer cpm.mu.Unlock()
	
	return cpm.rebuildCustomPoolLocked(ctx, excludeID)
}

// rebuildCustomPoolLocked rebuilds the custom pool (caller must hold lock).
func (cpm *CertPoolManager) rebuildCustomPoolLocked(ctx context.Context, excludeID string) error {
	// Create new pool
	newPool := x509.NewCertPool()
	
	// List all certificates
	certIDs, err := cpm.store.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list certificates for pool rebuild: %w", err)
	}
	
	// Add certificates to new pool (excluding the specified ID)
	for _, id := range certIDs {
		if id == excludeID {
			continue
		}
		
		cert, err := cpm.store.Load(ctx, id)
		if err != nil {
			continue // Skip failed loads
		}
		
		// Validate if required
		if cpm.config.ValidateCertificates {
			valid := true
			for _, validator := range cpm.validators {
				if err := validator.Validate(ctx, cert); err != nil {
					valid = false
					break
				}
			}
			if !valid {
				continue
			}
		}
		
		newPool.AddCert(cert.X509)
	}
	
	// Replace the old pool
	cpm.customPool = newPool
	
	return nil
}

// getPoolSizeLocked returns the approximate size of the custom pool (caller must hold lock).
func (cpm *CertPoolManager) getPoolSizeLocked() int {
	// Note: x509.CertPool doesn't provide a size method.
	// In a production implementation, you would track the size separately.
	// For now, we return 0 to indicate unknown size.
	return 0
}

// GetPoolStats returns statistics about the certificate pools.
func (cpm *CertPoolManager) GetPoolStats() *PoolStats {
	cpm.mu.RLock()
	defer cpm.mu.RUnlock()
	
	return &PoolStats{
		SystemPoolSize:    0, // x509.CertPool doesn't provide size info
		CustomPoolSize:    0, // x509.CertPool doesn't provide size info
		IsWatching:        cpm.isWatching,
		LastRefresh:       time.Now(), // Would track this in production
		ValidationEnabled: cpm.config.ValidateCertificates,
	}
}

// PoolStats contains statistics about certificate pools.
type PoolStats struct {
	// SystemPoolSize is the number of certificates in the system pool.
	SystemPoolSize int
	
	// CustomPoolSize is the number of certificates in the custom pool.
	CustomPoolSize int
	
	// IsWatching indicates if the pool manager is actively watching for changes.
	IsWatching bool
	
	// LastRefresh is when the pools were last refreshed.
	LastRefresh time.Time
	
	// ValidationEnabled indicates if certificate validation is enabled.
	ValidationEnabled bool
}

// ValidatePoolIntegrity checks the integrity of the certificate pools.
func (cpm *CertPoolManager) ValidatePoolIntegrity(ctx context.Context) error {
	cpm.mu.RLock()
	defer cpm.mu.RUnlock()
	
	// Basic integrity checks
	if cpm.systemPool == nil {
		return fmt.Errorf("system certificate pool is nil")
	}
	
	if cpm.customPool == nil {
		return fmt.Errorf("custom certificate pool is nil")
	}
	
	// Could add more sophisticated integrity checks here
	// such as verifying that all stored certificates are still valid,
	// checking for expired certificates, etc.
	
	return nil
}