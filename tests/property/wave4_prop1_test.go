package property_test

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"pgregory.net/rapid"

	"github.com/cnoe-io/idpbuilder/pkg/cmd/push"
)

// genCredentials generates random credential combinations
func genCredentials(t *rapid.T) struct {
	username string
	password string
	token    string
} {
	return struct {
		username string
		password string
		token    string
	}{
		username: rapid.String().Draw(t, "username"),
		password: rapid.String().Draw(t, "password"),
		token:    rapid.String().Draw(t, "token"),
	}
}

func TestProperty_W1_4_1_NoCredentialLogging(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random credentials with distinct prefixes to avoid substring collisions
		usernameRandomPart := rapid.String().Draw(t, "username")
		passwordRandomPart := rapid.String().Draw(t, "password")
		tokenRandomPart := rapid.String().Draw(t, "token")

		username := "CRED_USER_" + usernameRandomPart + "_END"
		password := "CRED_PASS_" + passwordRandomPart + "_END"
		token := "CRED_TOKEN_" + tokenRandomPart + "_END"

		// Capture log output
		var buf bytes.Buffer
		logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))

		// Log credential resolution (should NOT log values)
		push.LogCredentialResolution(
			logger,
			"flags",
			len(username) > 0,
			len(password) > 0,
			len(token) > 0,
		)

		output := buf.String()

		// Property: No credential values in output
		// Check that the actual credential values don't appear
		assert.NotContains(t, output, username,
			"username value should NOT be logged")
		assert.NotContains(t, output, password,
			"password value should NOT be logged")
		assert.NotContains(t, output, token,
			"token value should NOT be logged")

		// Should only log presence/absence flags
		assert.Contains(t, output, "has_username=true")
		assert.Contains(t, output, "has_password=true")
		assert.Contains(t, output, "has_token=true")
	})
}
