// Package tls provides TLS configuration for registry connections.
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

// NewConfigProvider creates a TLS config provider
func NewConfigProvider(insecure bool) ConfigProvider {
	return &configProvider{
		insecure: insecure,
	}
}

// configProvider implements TLS configuration
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
