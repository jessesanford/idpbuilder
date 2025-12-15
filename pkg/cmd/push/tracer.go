// Package push implements OCI registry push functionality for idpbuilder.
package push

import (
	"log/slog"
	"os"
)

// NewDebugLogger creates logger configured for debug mode.
// Per REQ-005: When --log-level debug, enable verbose output.
// Per REQ-006: When --log-level info, log only high-level operational steps.
func NewDebugLogger(level slog.Level, phase string) *slog.Logger {
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	})
	return slog.New(handler).With(slog.String("phase", phase))
}

// LogCredentialResolution logs credential flow at debug level.
// Per REQ-020: NEVER log actual credential values.
// Only logs presence/absence flags for diagnosis.
func LogCredentialResolution(logger *slog.Logger, source string, hasUsername, hasPassword, hasToken bool) {
	logger.Debug("credential resolution",
		slog.String("source", source), // "flags", "env", "anonymous"
		slog.Bool("has_username", hasUsername),
		slog.Bool("has_password", hasPassword),
		slog.Bool("has_token", hasToken),
	)
}

// DebugTransport is imported from pkg/registry to avoid duplication
// See pkg/registry/debugtransport.go for implementation details
