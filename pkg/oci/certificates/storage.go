// Package certificates provides certificate storage and management functionality
// for the idpbuilder OCI management system.
package certificates

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// EventType represents the type of certificate storage event.
type EventType int

const (
	// EventAdded indicates a certificate was added to storage.
	EventAdded EventType = iota
	// EventModified indicates a certificate was modified in storage.
	EventModified
	// EventDeleted indicates a certificate was deleted from storage.
	EventDeleted
)

// String returns the string representation of the EventType.
func (e EventType) String() string {
	switch e {
	case EventAdded:
		return "added"
	case EventModified:
		return "modified"
	case EventDeleted:
		return "deleted"
	default:
		return "unknown"
	}
}

// Event represents a storage event for certificate changes.
type Event struct {
	// Type is the type of event that occurred.
	Type EventType
	// ID is the unique identifier of the certificate.
	ID string
	// Timestamp is when the event occurred.
	Timestamp time.Time
}

// Certificate represents a certificate with its associated metadata.
type Certificate struct {
	// ID is the unique identifier for this certificate.
	ID string
	// Data contains the raw certificate data (PEM encoded).
	Data []byte
	// PrivateKey contains the private key data (PEM encoded), if available.
	PrivateKey []byte
	// X509 is the parsed X.509 certificate.
	X509 *x509.Certificate
	// CreatedAt is when the certificate was stored.
	CreatedAt time.Time
	// UpdatedAt is when the certificate was last modified.
	UpdatedAt time.Time
}

// CertificateValidator defines an interface for certificate validation.
type CertificateValidator interface {
	// Validate checks if a certificate is valid and safe to store.
	Validate(ctx context.Context, cert *Certificate) error
}

// CertificateStore defines the interface for certificate storage operations.
type CertificateStore interface {
	// Save stores a certificate with the given ID.
	Save(ctx context.Context, id string, cert *Certificate) error
	
	// Load retrieves a certificate by its ID.
	Load(ctx context.Context, id string) (*Certificate, error)
	
	// Delete removes a certificate by its ID.
	Delete(ctx context.Context, id string) error
	
	// List returns all certificate IDs currently in storage.
	List(ctx context.Context) ([]string, error)
	
	// Watch monitors the storage for changes and calls the provided function
	// when changes occur. The context can be used to stop watching.
	Watch(ctx context.Context, onChange func(Event)) error
	
	// DiscoverCertificates scans well-known locations for certificates
	// and imports them into the store.
	DiscoverCertificates(ctx context.Context) error
	
	// Close gracefully shuts down the store and any background operations.
	Close() error
}

// FilesystemStore implements CertificateStore using the filesystem for persistence.
type FilesystemStore struct {
	// basePath is the root directory for certificate storage.
	basePath string
	
	// mu protects concurrent access to the store.
	mu sync.RWMutex
	
	// watcher monitors filesystem changes for hot-reload functionality.
	watcher *fsnotify.Watcher
	
	// validators are used to validate certificates before storage.
	validators []CertificateValidator
	
	// watchChan is used to communicate watch events.
	watchChan chan Event
	
	// stopChan signals the watch goroutine to stop.
	stopChan chan struct{}
	
	// watchOnce ensures the watch goroutine is started only once.
	watchOnce sync.Once
}

// NewFilesystemStore creates a new filesystem-based certificate store.
func NewFilesystemStore(basePath string, validators ...CertificateValidator) (*FilesystemStore, error) {
	if basePath == "" {
		return nil, fmt.Errorf("basePath cannot be empty")
	}
	
	// Ensure the base directory exists with proper permissions
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create base directory: %w", err)
	}
	
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create filesystem watcher: %w", err)
	}
	
	store := &FilesystemStore{
		basePath:   basePath,
		watcher:    watcher,
		validators: validators,
		watchChan:  make(chan Event, 100), // Buffer to prevent blocking
		stopChan:   make(chan struct{}),
	}
	
	return store, nil
}

// Save implements the CertificateStore interface.
func (fs *FilesystemStore) Save(ctx context.Context, id string, cert *Certificate) error {
	if id == "" {
		return fmt.Errorf("certificate ID cannot be empty")
	}
	if cert == nil {
		return fmt.Errorf("certificate cannot be nil")
	}
	
	// Validate the certificate using all configured validators
	for _, validator := range fs.validators {
		if err := validator.Validate(ctx, cert); err != nil {
			return fmt.Errorf("certificate validation failed: %w", err)
		}
	}
	
	fs.mu.Lock()
	defer fs.mu.Unlock()
	
	certPath := filepath.Join(fs.basePath, id+".crt")
	keyPath := filepath.Join(fs.basePath, id+".key")
	
	// Check if certificate already exists for event type determination
	exists := fs.existsLocked(id)
	
	// Create a backup if the certificate already exists
	if exists {
		if err := fs.createBackupLocked(id); err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}
	}
	
	// Write certificate file atomically
	if err := fs.writeFileAtomic(certPath, cert.Data, 0644); err != nil {
		return fmt.Errorf("failed to write certificate file: %w", err)
	}
	
	// Write private key file if present (with restricted permissions)
	if len(cert.PrivateKey) > 0 {
		if err := fs.writeFileAtomic(keyPath, cert.PrivateKey, 0600); err != nil {
			// Cleanup certificate file on failure
			os.Remove(certPath)
			return fmt.Errorf("failed to write private key file: %w", err)
		}
	}
	
	// Emit event
	eventType := EventAdded
	if exists {
		eventType = EventModified
	}
	
	select {
	case fs.watchChan <- Event{Type: eventType, ID: id, Timestamp: time.Now()}:
	default:
		// Channel full, skip event (non-blocking)
	}
	
	return nil
}

// Load implements the CertificateStore interface.
func (fs *FilesystemStore) Load(ctx context.Context, id string) (*Certificate, error) {
	if id == "" {
		return nil, fmt.Errorf("certificate ID cannot be empty")
	}
	
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	
	certPath := filepath.Join(fs.basePath, id+".crt")
	keyPath := filepath.Join(fs.basePath, id+".key")
	
	// Read certificate file
	certData, err := os.ReadFile(certPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("certificate not found: %s", id)
		}
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}
	
	// Parse the certificate
	block, _ := pem.Decode(certData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM certificate")
	}
	
	x509Cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse X.509 certificate: %w", err)
	}
	
	cert := &Certificate{
		ID:   id,
		Data: certData,
		X509: x509Cert,
	}
	
	// Read private key if it exists
	if keyData, err := os.ReadFile(keyPath); err == nil {
		cert.PrivateKey = keyData
	}
	
	// Get file timestamps
	if info, err := os.Stat(certPath); err == nil {
		cert.UpdatedAt = info.ModTime()
		// For filesystem storage, we approximate CreatedAt as ModTime
		// In a real implementation, this might be stored as metadata
		cert.CreatedAt = info.ModTime()
	}
	
	return cert, nil
}

// Delete implements the CertificateStore interface.
func (fs *FilesystemStore) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("certificate ID cannot be empty")
	}
	
	fs.mu.Lock()
	defer fs.mu.Unlock()
	
	certPath := filepath.Join(fs.basePath, id+".crt")
	keyPath := filepath.Join(fs.basePath, id+".key")
	
	// Check if certificate exists
	if !fs.existsLocked(id) {
		return fmt.Errorf("certificate not found: %s", id)
	}
	
	// Create backup before deletion
	if err := fs.createBackupLocked(id); err != nil {
		return fmt.Errorf("failed to create backup before deletion: %w", err)
	}
	
	// Remove certificate file
	if err := os.Remove(certPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove certificate file: %w", err)
	}
	
	// Remove private key file if it exists
	if err := os.Remove(keyPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove private key file: %w", err)
	}
	
	// Emit event
	select {
	case fs.watchChan <- Event{Type: EventDeleted, ID: id, Timestamp: time.Now()}:
	default:
		// Channel full, skip event (non-blocking)
	}
	
	return nil
}

// List implements the CertificateStore interface.
func (fs *FilesystemStore) List(ctx context.Context) ([]string, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	
	entries, err := os.ReadDir(fs.basePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read storage directory: %w", err)
	}
	
	var certIDs []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".crt" {
			// Remove the .crt extension to get the certificate ID
			id := entry.Name()[:len(entry.Name())-4]
			certIDs = append(certIDs, id)
		}
	}
	
	return certIDs, nil
}

// Watch implements the CertificateStore interface.
func (fs *FilesystemStore) Watch(ctx context.Context, onChange func(Event)) error {
	if onChange == nil {
		return fmt.Errorf("onChange callback cannot be nil")
	}
	
	// Start the watch goroutine only once
	fs.watchOnce.Do(func() {
		go fs.watchLoop(ctx)
	})
	
	// Add the storage directory to the watcher
	if err := fs.watcher.Add(fs.basePath); err != nil {
		return fmt.Errorf("failed to watch directory: %w", err)
	}
	
	// Forward events to the callback
	go func() {
		for {
			select {
			case event := <-fs.watchChan:
				onChange(event)
			case <-ctx.Done():
				return
			}
		}
	}()
	
	return nil
}

// watchLoop runs the filesystem watching loop.
func (fs *FilesystemStore) watchLoop(ctx context.Context) {
	for {
		select {
		case event, ok := <-fs.watcher.Events:
			if !ok {
				return
			}
			
			// Only process certificate files
			if filepath.Ext(event.Name) == ".crt" {
				id := filepath.Base(event.Name)[:len(filepath.Base(event.Name))-4]
				
				var eventType EventType
				if event.Op&fsnotify.Create == fsnotify.Create {
					eventType = EventAdded
				} else if event.Op&fsnotify.Write == fsnotify.Write {
					eventType = EventModified
				} else if event.Op&fsnotify.Remove == fsnotify.Remove {
					eventType = EventDeleted
				} else {
					continue // Ignore other operations
				}
				
				select {
				case fs.watchChan <- Event{Type: eventType, ID: id, Timestamp: time.Now()}:
				default:
					// Channel full, skip event
				}
			}
			
		case err, ok := <-fs.watcher.Errors:
			if !ok {
				return
			}
			// In a production system, we would log this error
			_ = err
			
		case <-ctx.Done():
			return
		}
	}
}

// DiscoverCertificates implements the CertificateStore interface.
func (fs *FilesystemStore) DiscoverCertificates(ctx context.Context) error {
	wellKnownPaths := []string{
		"/etc/ssl/certs",
		"/etc/pki/tls/certs",
		"/usr/local/share/ca-certificates",
	}
	
	// Add user-specific paths if HOME is set
	if home := os.Getenv("HOME"); home != "" {
		wellKnownPaths = append(wellKnownPaths, filepath.Join(home, ".docker/certs.d"))
	}
	
	var discovered int
	for _, path := range wellKnownPaths {
		count, err := fs.discoverFromPath(ctx, path)
		if err != nil {
			// Continue with other paths on error
			continue
		}
		discovered += count
	}
	
	return nil
}

// discoverFromPath discovers certificates from a specific path.
func (fs *FilesystemStore) discoverFromPath(ctx context.Context, path string) (int, error) {
	var count int
	
	err := filepath.WalkDir(path, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // Continue walking on errors
		}
		
		if d.IsDir() {
			return nil
		}
		
		// Only process files with certificate extensions
		ext := filepath.Ext(filePath)
		if ext != ".crt" && ext != ".pem" && ext != ".cer" {
			return nil
		}
		
		// Read and validate the certificate
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil
		}
		
		block, _ := pem.Decode(data)
		if block == nil || block.Type != "CERTIFICATE" {
			return nil
		}
		
		x509Cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil
		}
		
		// Generate an ID based on the file path
		id := fmt.Sprintf("discovered_%s", filepath.Base(filePath))
		
		cert := &Certificate{
			ID:   id,
			Data: data,
			X509: x509Cert,
		}
		
		// Save the discovered certificate
		if err := fs.Save(ctx, id, cert); err != nil {
			return nil // Continue on save errors
		}
		
		count++
		return nil
	})
	
	return count, err
}

// Close implements the CertificateStore interface.
func (fs *FilesystemStore) Close() error {
	close(fs.stopChan)
	
	if fs.watcher != nil {
		return fs.watcher.Close()
	}
	
	return nil
}

// existsLocked checks if a certificate exists (caller must hold lock).
func (fs *FilesystemStore) existsLocked(id string) bool {
	certPath := filepath.Join(fs.basePath, id+".crt")
	_, err := os.Stat(certPath)
	return err == nil
}

// createBackupLocked creates a backup of the certificate (caller must hold lock).
func (fs *FilesystemStore) createBackupLocked(id string) error {
	certPath := filepath.Join(fs.basePath, id+".crt")
	keyPath := filepath.Join(fs.basePath, id+".key")
	
	timestamp := time.Now().Format("20060102-150405")
	
	// Backup certificate
	if _, err := os.Stat(certPath); err == nil {
		backupPath := filepath.Join(fs.basePath, fmt.Sprintf("%s.crt.backup.%s", id, timestamp))
		if err := fs.copyFile(certPath, backupPath); err != nil {
			return err
		}
	}
	
	// Backup private key if it exists
	if _, err := os.Stat(keyPath); err == nil {
		backupPath := filepath.Join(fs.basePath, fmt.Sprintf("%s.key.backup.%s", id, timestamp))
		if err := fs.copyFile(keyPath, backupPath); err != nil {
			return err
		}
	}
	
	return nil
}

// writeFileAtomic writes data to a file atomically by writing to a temporary file first.
func (fs *FilesystemStore) writeFileAtomic(filename string, data []byte, perm os.FileMode) error {
	tmpFile := filename + ".tmp"
	
	if err := os.WriteFile(tmpFile, data, perm); err != nil {
		return err
	}
	
	return os.Rename(tmpFile, filename)
}

// copyFile copies a file from src to dst.
func (fs *FilesystemStore) copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	
	return os.WriteFile(dst, data, info.Mode())
}