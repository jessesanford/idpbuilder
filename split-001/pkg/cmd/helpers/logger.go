package helpers

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"
)

var (
	// LogLevel controls the logging level
	LogLevel string = "info"
	// LogLevelMsg provides help text for log level flag
	LogLevelMsg = "Set log level (debug, info, warn, error)"
	// Verbose controls verbose output
	Verbose bool
	// VerboseMsg provides help text for verbose flag
	VerboseMsg = "Enable verbose output"
	// logger is the structured logger instance
	logger *slog.Logger
)

// LogLevelValue represents log level values
type LogLevelValue int

const (
	// DebugLevel for debug messages
	DebugLevel LogLevelValue = -4
	// InfoLevel for info messages
	InfoLevel LogLevelValue = 0
	// WarnLevel for warning messages
	WarnLevel LogLevelValue = 4
	// ErrorLevel for error messages
	ErrorLevel LogLevelValue = 8
)

// InitializeLogging sets up the logging system
func InitializeLogging() {
	level := parseLogLevel(LogLevel)
	
	// Create handler options
	opts := &slog.HandlerOptions{
		Level: slog.Level(level),
		AddSource: Verbose,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Customize timestamp format
			if a.Key == slog.TimeKey {
				return slog.String(slog.TimeKey, a.Value.Time().Format(time.RFC3339))
			}
			// Shorten source paths
			if a.Key == slog.SourceKey {
				source := a.Value.Any().(*slog.Source)
				parts := strings.Split(source.File, "/")
				if len(parts) > 2 {
					source.File = strings.Join(parts[len(parts)-2:], "/")
				}
			}
			return a
		},
	}

	var handler slog.Handler
	if ColoredOutput {
		handler = NewColorHandler(os.Stderr, opts)
	} else {
		handler = slog.NewTextHandler(os.Stderr, opts)
	}

	logger = slog.New(handler)
	slog.SetDefault(logger)
}

// LogDebug logs a debug message
func LogDebug(format string, args ...interface{}) {
	if logger != nil {
		logger.Debug(fmt.Sprintf(format, args...))
	}
}

// LogInfo logs an info message
func LogInfo(format string, args ...interface{}) {
	if logger != nil {
		logger.Info(fmt.Sprintf(format, args...))
	}
}

// LogWarn logs a warning message
func LogWarn(format string, args ...interface{}) {
	if logger != nil {
		logger.Warn(fmt.Sprintf(format, args...))
	}
}

// LogError logs an error message
func LogError(format string, args ...interface{}) {
	if logger != nil {
		logger.Error(fmt.Sprintf(format, args...))
	}
}

// LogWarning is an alias for LogWarn for consistency with output functions
func LogWarning(format string, args ...interface{}) {
	LogWarn(format, args...)
}

// parseLogLevel converts string log level to LogLevelValue
func parseLogLevel(level string) LogLevelValue {
	switch strings.ToLower(level) {
	case "debug":
		return DebugLevel
	case "info", "":
		return InfoLevel
	case "warn", "warning":
		return WarnLevel
	case "error":
		return ErrorLevel
	default:
		// Default to info for invalid levels
		LogWarning("Invalid log level '%s', using 'info'", level)
		return InfoLevel
	}
}

// ColorHandler is a custom handler that adds colors to log output
type ColorHandler struct {
	handler slog.Handler
	writer  io.Writer
}

// NewColorHandler creates a new colored log handler
func NewColorHandler(w io.Writer, opts *slog.HandlerOptions) *ColorHandler {
	return &ColorHandler{
		handler: slog.NewTextHandler(w, opts),
		writer:  w,
	}
}

// Enabled implements slog.Handler
func (h *ColorHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

// Handle implements slog.Handler with color support
func (h *ColorHandler) Handle(ctx context.Context, record slog.Record) error {
	// Get level color
	var color string
	switch {
	case record.Level >= slog.LevelError:
		color = ColorRed
	case record.Level >= slog.LevelWarn:
		color = ColorYellow
	case record.Level >= slog.LevelInfo:
		color = ColorGreen
	default:
		color = ColorCyan
	}

	// Format the message with colors
	levelStr := record.Level.String()
	if ColoredOutput {
		levelStr = color + strings.ToUpper(levelStr) + ColorReset
	}

	// Build the log message
	timestamp := record.Time.Format("15:04:05")
	message := fmt.Sprintf("%s[%s]%s %s %s",
		ColorCyan, timestamp, ColorReset,
		levelStr,
		record.Message)

	// Add attributes if any
	record.Attrs(func(a slog.Attr) bool {
		message += fmt.Sprintf(" %s=%v", a.Key, a.Value)
		return true
	})

	// Write to output
	fmt.Fprintln(h.writer, message)
	return nil
}

// WithAttrs implements slog.Handler
func (h *ColorHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ColorHandler{
		handler: h.handler.WithAttrs(attrs),
		writer:  h.writer,
	}
}

// WithGroup implements slog.Handler
func (h *ColorHandler) WithGroup(name string) slog.Handler {
	return &ColorHandler{
		handler: h.handler.WithGroup(name),
		writer:  h.writer,
	}
}

// SetLogLevel updates the global log level
func SetLogLevel(level string) {
	LogLevel = level
	InitializeLogging()
}

// GetLogger returns the configured logger instance
func GetLogger() *slog.Logger {
	return logger
}

// LogWithFields logs a message with structured fields
func LogWithFields(level LogLevelValue, message string, fields map[string]interface{}) {
	if logger == nil {
		return
	}

	var attrs []slog.Attr
	for key, value := range fields {
		attrs = append(attrs, slog.Any(key, value))
	}

	logger.LogAttrs(nil, slog.Level(level), message, attrs...)
}