// Package tls provides TLS configuration for registry connections.
// This is a Phase 1 stub interface for Phase 2 development.
package tls

import (
	"crypto/tls"
)

// ConfigProvider supplies TLS configuration
type ConfigProvider interface {
	// GetTLSConfig returns the TLS configuration
	GetTLSConfig() *tls.Config

	// IsInsecure returns whether TLS verification is disabled
	IsInsecure() bool
}

// NewConfigProvider creates a TLS config provider (stub for Phase 1)
func NewConfigProvider(insecure bool) ConfigProvider {
	return &configProvider{
		insecure: insecure,
	}
}

// configProvider is a minimal stub for planning purposes
type configProvider struct {
	insecure bool
}

func (p *configProvider) GetTLSConfig() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: p.insecure,
	}
}

func (p *configProvider) IsInsecure() bool {
	return p.insecure
}
