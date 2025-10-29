package tls_test

import (
	"testing"

	tlspkg "github.com/cnoe-io/idpbuilder/pkg/tls"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	t.Run("creates config with InsecureSkipVerify false", func(t *testing.T) {
		config := tlspkg.Config{
			InsecureSkipVerify: false,
		}

		assert.False(t, config.InsecureSkipVerify)
	})

	t.Run("creates config with InsecureSkipVerify true", func(t *testing.T) {
		config := tlspkg.Config{
			InsecureSkipVerify: true,
		}

		assert.True(t, config.InsecureSkipVerify)
	})
}
