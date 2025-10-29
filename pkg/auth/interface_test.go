package auth_test

import (
	"testing"

	"github.com/cnoe-io/idpbuilder/pkg/auth"
	"github.com/stretchr/testify/assert"
)

func TestCredentials(t *testing.T) {
	t.Run("creates credentials with username and password", func(t *testing.T) {
		creds := auth.Credentials{
			Username: "testuser",
			Password: "testpass",
		}

		assert.Equal(t, "testuser", creds.Username)
		assert.Equal(t, "testpass", creds.Password)
	})

	t.Run("supports special characters in password", func(t *testing.T) {
		creds := auth.Credentials{
			Username: "user",
			Password: `p@ss"w'o rd!#$%^&*()[]{}`,
		}

		assert.Equal(t, "user", creds.Username)
		assert.Contains(t, creds.Password, `"`)
		assert.Contains(t, creds.Password, `'`)
		assert.Contains(t, creds.Password, ` `)
	})
}
