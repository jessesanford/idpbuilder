package property_test

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"pgregory.net/rapid"

	"github.com/cnoe-io/idpbuilder/pkg/registry"
)

// mockRoundTripper for testing
type mockRoundTripper struct {
	response *http.Response
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.response, nil
}

// genHTTPRequest generates random HTTP request parameters
func genHTTPRequest(t *rapid.T) struct {
	method     string
	statusCode int
} {
	return struct {
		method     string
		statusCode int
	}{
		method:     rapid.SampledFrom([]string{"GET", "POST", "PUT", "HEAD"}).Draw(t, "method"),
		statusCode: rapid.IntRange(200, 599).Draw(t, "status_code"),
	}
}

func TestProperty_W1_4_2_RequestResponseCorrelation(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random HTTP request
		reqParams := genHTTPRequest(t)

		var buf bytes.Buffer
		logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))

		// Create debug transport
		mockRT := &mockRoundTripper{
			response: &http.Response{
				StatusCode: reqParams.statusCode,
				Status:     fmt.Sprintf("%d Status", reqParams.statusCode),
				Body:       io.NopCloser(bytes.NewBufferString("")),
			},
		}

		transport := &registry.DebugTransport{
			Base:   mockRT,
			Logger: logger,
		}

		req, _ := http.NewRequest(reqParams.method, "https://registry.example.com/v2/", nil)
		_, _ = transport.RoundTrip(req)

		output := buf.String()
		lines := strings.Split(output, "\n")

		// Extract request_id from request log
		var requestID string
		for _, line := range lines {
			if strings.Contains(line, "HTTP request") && strings.Contains(line, "request_id=") {
				parts := strings.Split(line, "request_id=")
				if len(parts) > 1 {
					requestID = strings.Fields(parts[1])[0]
					break
				}
			}
		}

		// Property: Request ID must be present
		assert.NotEmpty(t, requestID, "request should have request_id")

		// Property: Response log has same request_id
		responseFound := false
		for _, line := range lines {
			if strings.Contains(line, "HTTP response") && strings.Contains(line, requestID) {
				responseFound = true
				break
			}
		}

		assert.True(t, responseFound,
			"response should be logged with same request_id as request")
	})
}
