package certificates

import (
	"context"
	"crypto/x509"
)

// CertificateStore defines the interface for certificate storage operations.
type CertificateStore interface {
	AddCertificate(ctx context.Context, cert *Certificate) error
	GetCertificate(ctx context.Context, id string) (*Certificate, error)
	UpdateCertificate(ctx context.Context, cert *Certificate) error
	DeleteCertificate(ctx context.Context, id string) error
	ListCertificates(ctx context.Context, filter *CertificateFilter) ([]*Certificate, error)
	SearchCertificates(ctx context.Context, query string) ([]*Certificate, error)
	GetCertificatesByStatus(ctx context.Context, status CertificateStatus) ([]*Certificate, error)
}

// CertificateValidator defines the interface for certificate validation operations.
type CertificateValidator interface {
	ValidateCertificate(ctx context.Context, cert *Certificate) (*ValidationResult, error)
	ValidatePEM(ctx context.Context, pemData []byte) (*ValidationResult, error)
	ValidateChain(ctx context.Context, chain []*Certificate) (*ValidationResult, error)
	AddValidationRule(ctx context.Context, rule *ValidationRule) error
	RemoveValidationRule(ctx context.Context, ruleName string) error
	ListValidationRules(ctx context.Context) ([]*ValidationRule, error)
}

// CertPoolManager defines the interface for managing certificate pools.
type CertPoolManager interface {
	CreatePool(ctx context.Context, name string) error
	DeletePool(ctx context.Context, name string) error
	AddCertificateToPool(ctx context.Context, poolName, certificateID string) error
	RemoveCertificateFromPool(ctx context.Context, poolName, certificateID string) error
	GetPool(ctx context.Context, name string) ([]*Certificate, error)
	GetPoolNames(ctx context.Context) ([]string, error)
	GetX509Pool(ctx context.Context, name string) (*x509.CertPool, error)
}

// ConfigLoader defines the interface for loading configuration data.
type ConfigLoader interface {
	LoadConfig(ctx context.Context, source string) (*Config, error)
	SaveConfig(ctx context.Context, config *Config, destination string) error
	ValidateConfig(ctx context.Context, config *Config) (*ValidationResult, error)
	GetDefaultConfig(ctx context.Context) (*Config, error)
}

// EventHandler defines the interface for handling certificate events.
type EventHandler interface {
	HandleEvent(ctx context.Context, event *Event) error
	SubscribeToEvents(ctx context.Context, eventTypes []EventType) (<-chan *Event, error)
	UnsubscribeFromEvents(ctx context.Context) error
}

// CertificateFilter provides filtering options for certificate queries.
type CertificateFilter struct {
	Status    []CertificateStatus   `json:"status,omitempty"`
	Tags      map[string]string     `json:"tags,omitempty"`
	Issuer    string                `json:"issuer,omitempty"`
	Subject   string                `json:"subject,omitempty"`
	ValidFrom *string               `json:"valid_from,omitempty"`
	ValidTo   *string               `json:"valid_to,omitempty"`
	KeyUsage  []string              `json:"key_usage,omitempty"`
	Limit     int                   `json:"limit,omitempty"`
	Offset    int                   `json:"offset,omitempty"`
}

// Config represents the configuration for the certificate management system.
type Config struct {
	StorageBackend  *StorageConfig    `json:"storage_backend"`
	ValidationRules []*ValidationRule `json:"validation_rules"`
	DefaultPools    []string          `json:"default_pools"`
	EventHandlers   *EventConfig      `json:"event_handlers"`
}

// StorageConfig defines storage backend configuration.
type StorageConfig struct {
	Type             string                 `json:"type"`
	ConnectionString string                 `json:"connection_string"`
	Options          map[string]interface{} `json:"options,omitempty"`
}

// EventConfig defines event handling configuration.
type EventConfig struct {
	Enabled    bool     `json:"enabled"`
	BufferSize int      `json:"buffer_size"`
	Handlers   []string `json:"handlers"`
}