package tls_test

import (
	"crypto/tls"
	"testing"

	tlspkg "github.com/cnoe-io/idpbuilder/pkg/tls"
	"github.com/stretchr/testify/assert"
)

// MockConfigProvider is a mock implementation of tls.ConfigProvider for testing
type MockConfigProvider struct {
	config   *tls.Config
	insecure bool
}

func (m *MockConfigProvider) GetTLSConfig() *tls.Config {
	if m.config == nil {
		return &tls.Config{
			InsecureSkipVerify: m.insecure,
		}
	}
	return m.config
}

func (m *MockConfigProvider) IsInsecure() bool {
	return m.insecure
}

func TestMockConfigProvider(t *testing.T) {
	t.Run("mock provider implements ConfigProvider interface", func(t *testing.T) {
		var _ tlspkg.ConfigProvider = &MockConfigProvider{}
	})

	t.Run("mock provider returns TLS config in secure mode", func(t *testing.T) {
		mock := &MockConfigProvider{
			insecure: false,
		}

		config := mock.GetTLSConfig()
		assert.NotNil(t, config)
		assert.False(t, config.InsecureSkipVerify)
		assert.False(t, mock.IsInsecure())
	})

	t.Run("mock provider returns TLS config in insecure mode", func(t *testing.T) {
		mock := &MockConfigProvider{
			insecure: true,
		}

		config := mock.GetTLSConfig()
		assert.NotNil(t, config)
		assert.True(t, config.InsecureSkipVerify)
		assert.True(t, mock.IsInsecure())
	})

	t.Run("mock provider returns custom TLS config", func(t *testing.T) {
		customConfig := &tls.Config{
			InsecureSkipVerify: false,
			MinVersion:         tls.VersionTLS12,
		}
		mock := &MockConfigProvider{
			config:   customConfig,
			insecure: false,
		}

		config := mock.GetTLSConfig()
		assert.Equal(t, customConfig, config)
		assert.Equal(t, uint16(tls.VersionTLS12), config.MinVersion)
	})
}
